# TODO

Time-bound items. Cross off as completed. Things that don't have a date go in `CLAUDE.md` roadmap, not here.

## ✅ Jun 30 2026 — prep-first cold start shipped (first A/B test)

Strategic call (user, Jun 30): don't gut the tracker, but pivot the *story* to
interview prep and test how prep-first the cold start should be. Shipped to prod
(deploy #74, green; migration 0019). See `AB_TESTS.md` "prep-first cold start".

- New signups (`users.onboarding_variant='prepfirst'`, 100%) land on "Who are you
  interviewing with?" → instant **company Playbook**, instead of the guided tour.
  The prep question creates the first tracked application as a byproduct — spine
  intact. Control = the old tour (NULL/`tour` variant).
- Funnel instrumented: `onboard_variant_assigned`, `prepfirst_prompt_view`,
  `prepfirst_submit`, `prepfirst_generate_ok|error`, `prepfirst_skip`.
- **Measure (give it a cohort):** signup → first playbook (activation) for the
  new prepfirst cohort vs the historical tracker-first baseline (16/25 activated).
  Guardrail: do prepfirst users add a 2nd+ application (the spine/retention)?
- [ ] Live-QA the flow on prod (throwaway company via `?onboard=prepfirst`).
- [ ] After ~1 cohort, read the prepfirst funnel + decide: deepen the Playbook
      (debrief feed-forward loop) vs. iterate the cold start.

## ✅ Jun 29 2026 — activation re-read done → reframe shipped

Re-read the data (adoption/invite-funnel/users console pull). Funnel:
**48 invited → 26 signed in → 15 activated → 6 active.** vs Jun 22 (36→25→16→6):
+12 invites brought ~0 signins/activations, active flat at 6. Manual nudges moved
no one. Answers to the four pre-committed questions:

1. **Interview fix:** not broken, barely used. `interview_save` now fires (2 users/3
   saves) and `addevent_open`→save converts ~75% — but only 2 users touch it. Demote
   confirmed.
2. **Dossier still the wedge:** yes, by a wide margin — `dossier_open` 12 users / 79
   total / 16 in 7d (next surface is 10). Thesis holds on a 2nd week.
3. **Open→complete:** interview side now measurable; **application side still blind** —
   `addmodal_open` (3 users) << `application_create` (8) << `first_application` (13).
   The Jun-22 instrumentation task was only half-done (interview_save shipped,
   addmodal_open under-firing never root-caused). **Still open.**
4. **Nudges:** no measurable movement (signed-in/activated flat).

**Decision:** dossier-thesis branch fired → Pursuit is "AI interview prep + tracker
spine." New data added a louder leak the tree didn't predict: **invite→sign-in (22/48
cold)**, so the reframe (lead the funnel entry points with the wedge) outranks
deepening the dossier.

**Shipped Jun 29 (this session):** the reframe + renamed the feature **"dossier" →
"Playbook"** (UI only; `dossier_open`/`dossier_refresh` event names kept for time-series
continuity). Surfaces: landing tagline + request-access copy, onboarding tour (playbook
promoted to lead step 1/6), in-app prep section + Today CTAs.

**Still open after this ship (next candidates):**
- [ ] Root-cause `addmodal_open` under-firing — application open→complete is still
      unmeasurable (the half-done Jun-22 task).
- [ ] Invite→sign-in leak: 22/48 invited never signed in (gate already open, pending=0).
- [ ] Deepen the Playbook (company + JD summary, hiring-manager link, logo) — for the
      already-hooked cohort, lower priority than the top-of-funnel leak.
- [ ] Watch whether the reframe lifts invite→sign-in on the next cohort before
      pouring more invites in.

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

## Product analytics (GA4 — funnel-first)

- [ ] Create GA4 property + web data stream; put the Measurement ID in a `GA4_MEASUREMENT_ID` env var and inject the gtag snippet into the SvelteKit shell (skip in dev / on localhost)
- [ ] SPA pageview tracking on Svelte client-side navigation (GA4 won't auto-fire on route changes)
- [ ] Instrument the AI-native moments — one event each: `paste_parse`, `screenshot_parse`, `dossier_open`, `dossier_refresh`, `interview_parse`
- [ ] Instrument the application funnel — `application_create` + a `status_change` event carrying old→new status (so applied→screen→onsite→offer is reconstructable)
- [ ] Instrument onboarding + activation — `onboarding_view`, `onboarding_dismiss`, `demo_seed_click`, first-application milestone
- [ ] Saved GA4 Explorations: application funnel (applied→offer) and AI-feature adoption funnel
- [ ] Verify events land in GA4 Realtime/DebugView from prod before calling it done (rendered-surface bar: confirm in the GA4 UI, not just that gtag fired)

## Parked decisions

- Domain name (pursuit.app etc.) — deferred until after beta validates demand. Running on `<ip>.nip.io` with a real Lets-Encrypt cert until then.

## Parked (post-beta)

- Post-interview recording analysis — privacy/consent decision required first
- CV A/B with per-variant tracking — needs volume to show signal, build in v2
- Payments / public signup
- Custom domain swap (DNS + nginx server_name + cert reissue)
