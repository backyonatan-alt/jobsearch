# Grounding eval fixtures

Golden company briefs + the cases that produce them, for the accuracy checker
(`internal/grounding`). Two layers:

- **`go test ./internal/grounding/...`** — runs on every deploy, no API spend.
  Replays the briefs in `briefs/` against `cases.json` and asserts each still
  passes (right company, real deep-link citations). `negative/` holds
  deliberately-broken briefs that MUST fail, so the checker can't silently rot.
- **`cmd/groundingeval` / the `grounding-eval` workflow** — the live gate.
  Generates each case for real, checks it (incl. citation reachability over
  HTTP), and with `-update` rewrites the briefs here.

The briefs in `briefs/` are hand-seeded baselines (Lusha is the Jul 6
grounding-gate case; 365Scores is the Jun 30 wrong-company failure class).
Refresh them from real generations before a distribution wave:

    ANTHROPIC_API_KEY=… go run ./cmd/groundingeval -update

Add cases to `cases.json` — especially Israeli name-collisions (thin English
web, Hebrew sources, transliterated names): that's the terrain where grounding
breaks, and every case there is a regression that can't come back.
