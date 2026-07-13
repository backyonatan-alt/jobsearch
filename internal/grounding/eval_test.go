package grounding

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func loadCases(t *testing.T) []Case {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", "cases.json"))
	if err != nil {
		t.Fatalf("read cases.json: %v", err)
	}
	var cases []Case
	if err := json.Unmarshal(b, &cases); err != nil {
		t.Fatalf("parse cases.json: %v", err)
	}
	if len(cases) == 0 {
		t.Fatal("no cases defined")
	}
	return cases
}

func loadBrief(t *testing.T, path string) json.RawMessage {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return json.RawMessage(b)
}

// Every golden brief must pass every critical check. This is the deploy-time
// regression guard: it runs on every deploy with no API spend, and it documents
// the expected right-company answer for each case.
func TestGoldenBriefsPass(t *testing.T) {
	for _, c := range loadCases(t) {
		c := c
		t.Run(c.Slug, func(t *testing.T) {
			raw := loadBrief(t, filepath.Join("testdata", "briefs", c.Slug+".json"))
			res := CheckCompanyBrief(c, raw)
			if !res.Pass {
				t.Errorf("golden brief %q failed critical checks: %v", c.Slug, res.FailedCritical())
				for _, ch := range res.Checks {
					t.Logf("  [%v critical=%v] %s — %s", ch.Pass, ch.Critical, ch.Name, ch.Detail)
				}
			}
		})
	}
}

// The checker is only useful if it actually catches the two failure modes it
// exists for. These fixtures are deliberately broken and MUST fail.
func TestCheckerCatchesFailures(t *testing.T) {
	t.Run("wrong_company", func(t *testing.T) {
		c := Case{Slug: "365scores", ExpectDomain: "365scores.com"}
		raw := loadBrief(t, filepath.Join("testdata", "negative", "wrong_company.json"))
		res := CheckCompanyBrief(c, raw)
		if res.Pass {
			t.Fatal("wrong-company brief passed — the anti-collision check is not working")
		}
		if !failed(res, "right_company") {
			t.Errorf("expected right_company to fail; got %v", res.FailedCritical())
		}
	})

	t.Run("homepage_citations", func(t *testing.T) {
		c := Case{Slug: "lusha", ExpectDomain: "lusha.com"}
		raw := loadBrief(t, filepath.Join("testdata", "negative", "homepage_citations.json"))
		res := CheckCompanyBrief(c, raw)
		if res.Pass {
			t.Fatal("homepage-only citations passed — the deep-link check is not working")
		}
		if !failed(res, "citations_are_deeplinks") {
			t.Errorf("expected citations_are_deeplinks to fail; got %v", res.FailedCritical())
		}
	})

	t.Run("model_error", func(t *testing.T) {
		c := Case{Slug: "x", ExpectDomain: "x.com"}
		res := CheckCompanyBrief(c, json.RawMessage(`{"error":"could not find enough public information about this company"}`))
		if res.Pass {
			t.Fatal("an error response passed as a valid brief")
		}
	})
}

func TestNormalizeDomain(t *testing.T) {
	cases := map[string]string{
		"https://www.Lusha.com/careers": "lusha.com",
		"WWW.365scores.com":             "365scores.com",
		"lusha.com/":                    "lusha.com",
		"http://wix.com:8080/x":         "wix.com",
		"":                              "",
	}
	for in, want := range cases {
		if got := NormalizeDomain(in); got != want {
			t.Errorf("NormalizeDomain(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestIsDeepLink(t *testing.T) {
	deep := []string{
		"https://techcrunch.com/2026/03/lusha-raises/",
		"https://glassdoor.com/Interview/x.htm",
		"https://example.com/?q=1",
	}
	shallow := []string{
		"https://techcrunch.com",
		"https://www.glassdoor.com/",
		"ftp://example.com/x",
		"not a url",
		"",
	}
	for _, u := range deep {
		if !isDeepLink(u) {
			t.Errorf("isDeepLink(%q) = false, want true", u)
		}
	}
	for _, u := range shallow {
		if isDeepLink(u) {
			t.Errorf("isDeepLink(%q) = true, want false", u)
		}
	}
}

func failed(res Result, name string) bool {
	for _, c := range res.Checks {
		if c.Name == name && !c.Pass {
			return true
		}
	}
	return false
}
