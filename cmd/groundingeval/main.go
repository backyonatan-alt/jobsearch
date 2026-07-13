// Command groundingeval runs the company-brief accuracy checker against LIVE
// generations — the trust gate from STRATEGY §9.1. It's the paid, networked
// counterpart to the deterministic go test: it actually generates each case,
// checks the result (right company + real, reachable citations), and can
// refresh the golden fixtures the deploy-time test replays.
//
// Usage:
//
//	ANTHROPIC_API_KEY=… go run ./cmd/groundingeval            # check, exit 1 on any failure
//	ANTHROPIC_API_KEY=… go run ./cmd/groundingeval -update    # also rewrite golden fixtures
//	go run ./cmd/groundingeval -offline                       # re-check existing fixtures, no API
//
// Meant for the grounding-eval workflow (manual dispatch / pre-distribution),
// never the per-push CI — it costs tokens.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/backyonatan-alt/jobsearch/internal/grounding"
	"github.com/backyonatan-alt/jobsearch/internal/llm"
)

func main() {
	update := flag.Bool("update", false, "rewrite golden fixtures with the fresh generations")
	offline := flag.Bool("offline", false, "re-check existing golden fixtures without calling the API")
	checkLinks := flag.Bool("check-links", true, "fetch each citation to confirm it resolves (live only)")
	dir := flag.String("dir", "internal/grounding/testdata", "testdata directory")
	flag.Parse()

	cases, err := loadCases(filepath.Join(*dir, "cases.json"))
	if err != nil {
		fatal(err)
	}

	var client *llm.Client
	if !*offline {
		key := os.Getenv("ANTHROPIC_API_KEY")
		if key == "" {
			fatal(fmt.Errorf("ANTHROPIC_API_KEY not set (use -offline to re-check existing fixtures)"))
		}
		client = llm.New(key)
	}

	httpClient := &http.Client{Timeout: 15 * time.Second}
	statusOf := func(url string) (int, error) {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return 0, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; PursuitGroundingEval/1.0)")
		resp, err := httpClient.Do(req)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()
		return resp.StatusCode, nil
	}

	failures := 0
	for _, c := range cases {
		var raw json.RawMessage
		if *offline {
			raw, err = os.ReadFile(filepath.Join(*dir, "briefs", c.Slug+".json"))
			if err != nil {
				fmt.Printf("✗ %s — no fixture: %v\n", c.Slug, err)
				failures++
				continue
			}
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), 170*time.Second)
			raw, err = client.GenerateCompanyBrief(ctx, c.Company, c.Role, c.Location, c.JDURL, c.CompanyURL)
			cancel()
			if err != nil {
				fmt.Printf("✗ %s — generation failed: %v\n", c.Slug, err)
				failures++
				continue
			}
		}

		res := grounding.CheckCompanyBrief(c, raw)
		printResult(c, res)

		if !*offline && *checkLinks {
			for _, ch := range grounding.CheckCitationsReachable(raw, statusOf) {
				mark := "✓"
				if !ch.Pass {
					mark = "✗"
					if ch.Critical {
						res.Pass = false
					}
				}
				fmt.Printf("    %s %s — %s\n", mark, ch.Name, ch.Detail)
			}
		}

		if !res.Pass {
			failures++
		}
		if *update && !*offline {
			out := filepath.Join(*dir, "briefs", c.Slug+".json")
			pretty, _ := json.MarshalIndent(json.RawMessage(raw), "", "  ")
			if err := os.WriteFile(out, append(pretty, '\n'), 0o644); err != nil {
				fatal(fmt.Errorf("write fixture %s: %w", out, err))
			}
			fmt.Printf("    ↳ updated %s\n", out)
		}
	}

	fmt.Printf("\n%d/%d cases passed\n", len(cases)-failures, len(cases))
	if failures > 0 {
		os.Exit(1)
	}
}

func printResult(c grounding.Case, res grounding.Result) {
	mark := "✓"
	if !res.Pass {
		mark = "✗"
	}
	fmt.Printf("%s %s (%s)\n", mark, c.Slug, c.Company)
	for _, ch := range res.Checks {
		cm := "✓"
		if !ch.Pass {
			cm = "✗"
		}
		crit := ""
		if ch.Critical {
			crit = " [critical]"
		}
		fmt.Printf("    %s %s%s — %s\n", cm, ch.Name, crit, ch.Detail)
	}
}

func loadCases(path string) ([]grounding.Case, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cases []grounding.Case
	if err := json.Unmarshal(b, &cases); err != nil {
		return nil, err
	}
	return cases, nil
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(2)
}
