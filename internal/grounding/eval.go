// Package grounding is the automatic accuracy checker for company-brief
// generation — the "trustworthy prep" north star, finally made testable.
//
// The wedge is prep you can walk in with and cite. Two ways it silently breaks:
//   - it researches the WRONG same-named company (the "365scores" collision
//     class — our Israeli beachhead is the worst terrain for this), or
//   - its citations point at homepages, not the page that backs the claim
//     (so "where do you know this from?" has no answer).
//
// CheckCompanyBrief asserts both, deterministically, over a captured brief.
// It runs two ways: as a Go test over golden fixtures (every deploy, no API
// spend — see eval_test.go) and live over real generations (cmd/groundingeval,
// which also refreshes the fixtures).
package grounding

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// Case is one production-shaped generation to check: the inputs a user would
// give, plus the canonical domain the RIGHT company must resolve to.
type Case struct {
	Slug         string `json:"slug"`
	Company      string `json:"company"`
	Role         string `json:"role"`
	Location     string `json:"location,omitempty"`
	JDURL        string `json:"jd_url,omitempty"`
	CompanyURL   string `json:"company_url,omitempty"`
	ExpectDomain string `json:"expect_domain"`
	Note         string `json:"note,omitempty"`
}

// Check is one assertion's outcome. Critical checks failing means the brief is
// not trustworthy; non-critical ones are warnings worth surfacing.
type Check struct {
	Name     string `json:"name"`
	Pass     bool   `json:"pass"`
	Critical bool   `json:"critical"`
	Detail   string `json:"detail"`
}

// Result is the verdict for one case. Pass is true only when every critical
// check passed.
type Result struct {
	Slug   string  `json:"slug"`
	Pass   bool    `json:"pass"`
	Checks []Check `json:"checks"`
}

// FailedCritical returns the names of the critical checks that failed.
func (r Result) FailedCritical() []string {
	var out []string
	for _, c := range r.Checks {
		if c.Critical && !c.Pass {
			out = append(out, c.Name)
		}
	}
	return out
}

// companyBrief mirrors the fields of the company-brief JSON the checker reads.
// It intentionally ignores everything else the model emits.
type companyBrief struct {
	Identity struct {
		Name    string `json:"name"`
		Domain  string `json:"domain"`
		Summary string `json:"summary"`
	} `json:"identity"`
	Company struct {
		Blurb string `json:"blurb"`
	} `json:"company"`
	Sources []struct {
		Label string `json:"label"`
		Href  string `json:"href"`
	} `json:"sources"`
	Error string `json:"error"`
}

// CheckCompanyBrief runs every grounding assertion for one case against a
// captured brief. It never touches the network — citation *reachability* is a
// live-only concern (see CheckCitationsReachable).
func CheckCompanyBrief(c Case, raw json.RawMessage) Result {
	res := Result{Slug: c.Slug}
	add := func(name string, critical, pass bool, detail string) {
		res.Checks = append(res.Checks, Check{Name: name, Pass: pass, Critical: critical, Detail: detail})
	}

	var b companyBrief
	if err := json.Unmarshal(raw, &b); err != nil {
		add("parses", true, false, fmt.Sprintf("brief is not valid JSON: %v", err))
		res.Pass = false
		return res
	}
	if b.Error != "" {
		add("generated", true, false, fmt.Sprintf("model returned an error instead of a brief: %q", b.Error))
		res.Pass = false
		return res
	}

	// The core anti-collision check: did we research the company the user meant?
	got := NormalizeDomain(b.Identity.Domain)
	want := NormalizeDomain(c.ExpectDomain)
	add("right_company", true, got != "" && got == want,
		fmt.Sprintf("identity.domain=%q (want %q)", b.Identity.Domain, c.ExpectDomain))

	// The candidate must be able to confirm the company from the brief itself.
	add("identity_complete", true, b.Identity.Name != "" && b.Identity.Summary != "",
		fmt.Sprintf("name=%q summary=%q", b.Identity.Name, truncate(b.Identity.Summary, 60)))

	add("blurb_present", true, strings.TrimSpace(b.Company.Blurb) != "",
		fmt.Sprintf("blurb=%q", truncate(b.Company.Blurb, 60)))

	add("sources_present", true, len(b.Sources) >= 2,
		fmt.Sprintf("%d sources (need >= 2)", len(b.Sources)))

	// Every citation must be a deep link to the page that backs the claim —
	// never a bare homepage. This is Ayelet's "clicking opens the main page" bug.
	var homepages []string
	for _, s := range b.Sources {
		if !isDeepLink(s.Href) {
			homepages = append(homepages, s.Href)
		}
	}
	add("citations_are_deeplinks", true, len(homepages) == 0,
		fmt.Sprintf("homepage/invalid citations: %s", strings.Join(homepages, ", ")))

	// Duplicate citations are a quality smell, not a trust failure.
	seen := map[string]bool{}
	dupes := false
	for _, s := range b.Sources {
		k := NormalizeDomain(s.Href) + strings.TrimRight(pathOf(s.Href), "/")
		if seen[k] {
			dupes = true
		}
		seen[k] = true
	}
	add("no_duplicate_sources", false, !dupes, "")

	res.Pass = len(res.FailedCritical()) == 0
	return res
}

// CheckCitationsReachable fetches each citation and records whether it resolves.
// Live-only: pass an http getter (nil detail on success). This is where a dead
// or moved source link gets caught — it needs the network, so it never runs in
// the deterministic unit test.
func CheckCitationsReachable(raw json.RawMessage, statusOf func(url string) (int, error)) []Check {
	var b companyBrief
	if err := json.Unmarshal(raw, &b); err != nil {
		return []Check{{Name: "citations_reachable", Critical: true, Pass: false, Detail: "unparseable brief"}}
	}
	var checks []Check
	for _, s := range b.Sources {
		code, err := statusOf(s.Href)
		switch {
		case err != nil:
			checks = append(checks, Check{Name: "reachable:" + s.Href, Critical: true, Pass: false, Detail: err.Error()})
		case code >= 400:
			checks = append(checks, Check{Name: "reachable:" + s.Href, Critical: true, Pass: false, Detail: fmt.Sprintf("HTTP %d", code)})
		default:
			checks = append(checks, Check{Name: "reachable:" + s.Href, Critical: true, Pass: true, Detail: fmt.Sprintf("HTTP %d", code)})
		}
	}
	return checks
}

// NormalizeDomain reduces a domain or URL to a bare comparable host:
// lowercased, scheme/path/port stripped, leading "www." removed.
func NormalizeDomain(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return ""
	}
	if i := strings.Index(s, "://"); i != -1 {
		s = s[i+3:]
	}
	s = strings.TrimPrefix(s, "www.")
	if i := strings.IndexAny(s, "/?#"); i != -1 {
		s = s[:i]
	}
	if i := strings.Index(s, ":"); i != -1 {
		s = s[:i]
	}
	return strings.TrimSuffix(s, ".")
}

// isDeepLink is true when href is an absolute http(s) URL that points past the
// homepage — it has a real path segment or a query string.
func isDeepLink(href string) bool {
	u, err := url.Parse(strings.TrimSpace(href))
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	if u.Host == "" {
		return false
	}
	return strings.Trim(u.Path, "/") != "" || u.RawQuery != ""
}

func pathOf(href string) string {
	u, err := url.Parse(strings.TrimSpace(href))
	if err != nil {
		return ""
	}
	return u.Path
}

func truncate(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
