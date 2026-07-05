# TODO

Time-bound items. Cross off as completed. Things that don't have a date go in `CLAUDE.md` roadmap, not here.

## ⏳ ~Jul 7 2026 — read the re-engagement cohort (do this next)

Context: after shipping P0 (grounding+citations), P1 (stale counts), P2 (closed
status), P3a (debrief loop), the user sent a Hebrew re-engagement email to **all
invited** users (~Jun 30/Jul 1) announcing the pivot to an interview-prep system.
This is the "validate before scaling" gate — watch the cohort for ~a week, then
read. Pull via the console snippet (adoption/invite-funnel/users) at `/admin`.

**Pre-committed questions:**
1. Did people come back? `last_login_at` spike since the email; signed-in count
   vs the ~26 baseline.
2. Did the reframe convert? prep-first activation vs the tracker-first baseline
   (16/25 activated); `onboard_variant_assigned` prepfirst vs tour.
3. Did grounding hold? any wrong-company complaints / `dossier_refresh` with a
   company_url re-ground ("Not them?"); no repeat of the 365scores class.
4. Did anyone debrief? `debrief_save` count + whether `prep_accuracy` trends
   spot-on (our first real trust metric).

**Then decide:** if it reads green → staged batch-promote of pending beta-interest
(open wider). If bugs/leaks surface → fix those first. Don't open the gate wide
until the wow-moment is proven to survive real use.

## ⏳ Jun 30 2026 — Ayelet feedback → "trustworthy interview-ready prep" plan

First external tester (Ayelet) validated the wedge ("maybe the most important
value") but exposed that its real risk is **trust**, not features. North star
sharpened to: *prep you can walk in with and cite* (see CLAUDE.md roadmap "v0.4").
Full feedback in FREE_RUN_NOTES (Jun 30). Sequenced fixes:

**Phase 0 — Trust the wedge (priority) — ✅ SHIPPED Jun 30 (deploy #75):**
- [x] **Company disambiguation/grounding.** Generation now takes location + JD-URL
      (linkedin.com dropped) + optional confirmed company website; prompts treat
      them as authoritative, never default to a famous same-named company. Brief
      returns an `identity` block → UI shows "Researched: <name> · <domain>" with a
      "Not them? →" re-ground control (regenerates against a corrected website).
- [x] **Real citations.** Interviewer signals carry `source_url` deep links (render
      as links, not homepages); company brief returns a `sources` deep-link list,
      rendered as a Sources section.
- [ ] Watch live: does grounding actually fix the 365scores-class collision? Ask
      Ayelet to retry her real interview prep on prod.

**Phase 1 — Spine reliability — ✅ SHIPPED Jun 30 (deploy #76):**
- [x] **Stale counts** fixed at the root: the sidebar count was fetched once on
      mount and never refreshed. Layout now refetches on `pursuit:refresh` + focus +
      visibility; every mutation (create/delete/status/edit/import) dispatches it.
      Verified live: Board count 15→16 on create, no reload. Should also resolve the
      delete "two clicks" (count now updates instantly → no "did it work?" re-click).
- [x] **Desktop-only** beta notice ("Best on desktop") on narrow viewports (≤820px),
      dismissal sticks. Full mobile pass stays Phase 3.
- [ ] Confirm live the delete-confirm no longer needs two clicks; if it does, repro
      + fix ConfirmDialog separately.

**Phase 2 — UX & gaps — ✅ SHIPPED Jun 30 (deploy #77):**
- [x] **"Closed" status** (req cancelled mid-process) — neutral terminal, threaded
      through backend + dropdowns + muted pill + Board; in the funnel it's a neutral
      outcome kept OUT of reach/reply-rate/active so it never drags conversion.
- [x] **Today rebalanced** — right rail 1.08/0.92 → 1.32/0.68 so the brief gets room.
- [x] **Drag-to-reorder pipeline stages** (grip handle; up/down arrows kept as fallback).

**Phase 3 — Debrief feed-forward loop (planned Jun 30; decisions locked)** — the
core mechanic: retention + trust-closing + grounding that compounds. Round N+1
prep knowing what round N actually asked is something ChatGPT structurally can't
do. Locked decisions: **~20-sec debrief** (2 taps + 1 optional line); **enrich if
present, never block** generation; **build 3a first**.

*3a — core loop — ✅ SHIPPED Jun 30 (deploy #78, migration 0020):*
- [x] migration `0020_debriefs` (one per round: feel, prep_accuracy, topics, notes).
- [x] endpoints: GET `/applications/{id}/debriefs` + POST
      `/applications/{id}/interviews/{iid}/debrief` (upsert).
- [x] feed-forward: `GenerateInterviewerBrief(priorDebriefs)`; handler assembles
      earlier rounds' debriefs (starts_at < this round). Enrich-only, never blocks.
- [x] UI: 20-sec debrief card on a past round (prompt → 2-tap form → summary) +
      "Informed by your last round" chip when an earlier debrief fed the playbook.
- [x] events: `debrief_prompt_view`, `debrief_save {feel, prep_accuracy}`.
- [x] Watch live (Jul re-read): debrief=0 — because it required a calendar interview
      almost nobody creates. **Fixed:** one-tap rounds (deploy #79, migration 0021).

*Unblock — one-tap rounds — ✅ SHIPPED Jul (deploy #79, migration 0021):*
- [x] `interviews.scheduled` flag — a round no longer needs a date.
- [x] "+ Add round" chip picker + stage-done hook ("Just did the {stage} round?
      Debrief it →") so rounds get created from behavior users already do.
- [ ] Watch: rounds created per activated user, and `debrief_save` finally > 0.
      If debriefs now happen → build 3b (proactive Today prompt + admin prep-accuracy).

*3b — proactive + metrics (after 3a):*
- [ ] proactive Today prompt: once an interview's `starts_at` has passed and it's
      un-debriefed → "How did the {company} {round} go? → Debrief".
- [ ] admin Adoption: **prep-accuracy stat** (% "spot-on") — first real trust metric.
- [ ] give the old vague "Pursuit debrief" Today card real meaning ("what we learned").

*3c — mobile/PWA pass (separate track, deferred):*
- [ ] responsive fixes to the desktop-only surfaces Ayelet hit — vanishing Save
      button, Source dropdown, pipeline editor, board horizontal scroll, playbook
      readability — + a PWA manifest/install. Big + independent; after 3a/3b.

**Measurement:** add trust signals (playbook refresh/keep rate; post-interview
"was it accurate?"). **Ayelet:** desktop-only for now; after her real interview
tomorrow ask "did it help / was it right?" — gold-standard trust data point.

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
- [x] Live-QA the flow on prod (Jun 30, Claude-for-Chrome, all 8 steps green —
      real company-specific playbook in ~40s). Caught + fixed an empty-role 400
      (POST /applications requires role → Role now required in the prompt).
      Open `[bug]` to watch: delete-confirm needed two clicks once (see FREE_RUN_NOTES).
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
