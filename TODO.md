# TODO

Time-bound items. Cross off as completed. Things that don't have a date go in `CLAUDE.md` roadmap, not here.

## This week (v0.1 — spine)

- [x] Repo scaffold, Go module, Postgres migrations runner
- [x] `users`, `sessions`, `magic_links`, `applications` schema
- [x] Magic-link auth endpoints (`/api/auth/request`, `/api/auth/verify`)
- [x] Applications CRUD endpoints behind auth middleware
- [x] Static frontend shell (login → app shell → application list/kanban)
- [x] systemd unit + nginx config samples in `deploy/`
- [x] GH Actions CI (build/test) + deploy stub for Hetzner
- [x] Provision Hetzner VM (CX22), run `deploy/bootstrap.sh` (uses `<ip>.nip.io` as the hostname — no domain needed for the beta)
- [x] Add GH Actions secrets (`DEPLOY_HOST`, `DEPLOY_USER`, `DEPLOY_SSH_KEY`) + var `DEPLOY_ENABLED=true`, push to main → first deploy
- [x] Google Cloud OAuth project + consent screen + Web client; `GOOGLE_CLIENT_ID` / `GOOGLE_CLIENT_SECRET` set on the VM
- [ ] Replace magic-link auth with Google OAuth in the backend + frontend (in progress)
- [ ] First end-to-end with Google sign-in: open the URL, click Continue with Google, land on /app, add 3 real applications
- [ ] Frontend design pass: SvelteKit migration + 3 design directions deployed at `/preview/a|b|c`, pick one
- [ ] Rotate Google OAuth Client Secret (current one was exposed in chat transcript on May 22 2026)
- [ ] Send invite links to 3 friends for the closed beta

## Next (v0.2 — ingest + dossier)

- [ ] LinkedIn job URL paste → backend fetch + parse → prefilled new-application form
- [ ] `.ics` paste/upload → interview event linked to application
- [ ] Anthropic API integration (Claude) with prompt caching on the system prompt
- [ ] Interviewer dossier endpoint + frontend card

## Pre-launch (before sending the first invite)

- [ ] Postmark account + `MAIL_DRIVER=postmark` wired in
- [ ] Privacy note: closed beta, what we store, retention policy. Even one paragraph.
- [ ] Backup story for Postgres: nightly `pg_dump` → off-VM (S3 or Hetzner Storage Box)
- [ ] A single Plausible (or GA4) property wired into the frontend

## Product analytics

Architecture (locked Jun 7 2026): **GA4 = public homepage / acquisition only**
(it thresholds low-volume beta data); **first-party `events` table = in-app
product behaviour** (source of truth, read per-user journeys). Application-stage
funnel stays the existing `/app/funnel` DB query — not reconstructed in GA4.

### Shipped

- [x] **PR 1 (#13)** — GA4 property `G-XJFYSBRVEW` (Google Signals off); gtag injected server-side only when `GA4_MEASUREMENT_ID` set (dev/local emit nothing); SPA `page_view` on Svelte nav; homepage events `beta_interest_submit`, `login`. Verified live via curl + Realtime.
- [x] **PR 1.5 (#14)** — in-app-browser (LinkedIn/IG/FB) sign-in nudge → fixes Google `disallowed_useragent` 403; `signin_webview_nudge` event.
- [x] **PR 2 (#15)** — `events` table + `POST /api/events` + `GET /api/admin/events` (per-user timeline). Server-side AI-moment events `paste_parse` / `screenshot_parse` / `dossier_refresh` / `interview_parse` (outcome + `duration_ms`). Client events `application_create` {via}, `status_change` {from,to,surface} (confirmed-success), `dossier_open` (fire-once). `first_application` milestone (demo rows excluded). PII rule enforced at call sites.

### PR 3 — parked until beta data rolls in (Jun 7 2026)

- [ ] Onboarding/activation events: `tutorial_begin`, `tutorial_complete`, `onboarding_dismiss`, `demo_seed_click`
- [ ] Internal-traffic filter — exclude admin/own user_id so my own clicks don't skew the beta numbers (GA4 internal-traffic + an exclusion in the events queries)
- [ ] `ANALYTICS.md` tracking plan — North Star, activation/retention definitions, event dictionary, the PII rule, which store answers which question
- [ ] Saved analyses: GA4 acquisition Exploration (homepage→`beta_interest_submit`→`login`); Postgres activation funnel (real-data-only); **AI-moment → week-2 retention cohort** (the core-thesis test)
- [ ] Quality/voice layer (from the PM review): `ai_feedback` thumbs-up/down on dossier+parse; `outcome_selfreport` at offer; churn ping; beta-interview script

### Launch ops (Jun 7 2026 — LinkedIn beta post is live)

- [ ] Rotate the **Anthropic API key** — exposed in chat transcript (`tail` of `/opt/jobsearch/.env`). console.anthropic.com → revoke + new key → update VM `.env` → `systemctl restart jobsearch`.
- [ ] Invite discipline: requests land on `/admin/people` tagged `source=linkedin` — invite a **~10 cohort, waitlist the rest** so per-user interviews stay doable.

## Parked decisions

- Domain name (pursuit.app etc.) — deferred until after beta validates demand. Running on `<ip>.nip.io` with a real Lets-Encrypt cert until then.

## Parked (post-beta)

- Post-interview recording analysis — privacy/consent decision required first
- CV A/B with per-variant tracking — needs volume to show signal, build in v2
- Payments / public signup
- Custom domain swap (DNS + nginx server_name + cert reissue)
