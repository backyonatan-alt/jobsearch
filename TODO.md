# TODO

Time-bound items. Cross off as completed. Things that don't have a date go in `CLAUDE.md` roadmap, not here.

## ✅ Jul 13 2026 — Phase A trust rails (§9) — first two gates shipped

The Jul 12 re-read's fork said: activation holds, so the next lever is
distribution — but STRATEGY §9 gates the next wave on trust infrastructure.
Built both gating deliverables before any wider push:

- [x] **Grounding/citation eval harness (§9.1).** `internal/grounding` — the
      accuracy checker. `CheckCompanyBrief` asserts right-company
      (identity.domain == expected) + citations are real deep links (not
      homepages). Two layers: deterministic `go test` over golden fixtures
      (every deploy, no API spend, incl. Israeli-collision cases + negative
      fixtures that must fail) and `cmd/groundingeval` + the `grounding-eval`
      workflow (live: generates for real, checks citation reachability over
      HTTP, refreshes fixtures). PR #38.
- [x] **Interviewer-privacy policy published (§9.2).** `/privacy` page — public
      professional sources only, no personal-life inference, visible sourcing,
      a stated line, GDPR posture, one-email removal. Linked from the homepage
      + the playbook's interviewer disclaimer. PR #38.

**Still open in Phase A before distribution wave 2:**
- [ ] **Add the `ANTHROPIC_API_KEY` repo secret** so the `grounding-eval`
      workflow can run live. The deterministic test guards regressions for
      free; the live gate needs the secret (same value as the VM). Run it green
      once — that's the actual "grounding holds" evidence STRATEGY §12 wants —
      *before* pouring in a new cohort.
- [ ] OAuth secret rotation (the other half of §9.2 item 2).
- [ ] Phase A #3–5: mobile read-only playbook view, day-of cheat sheet,
      feedback-on-claims. Then Phase B (pricing experiment).

> Distribution stays gated until the live eval runs green. Growth also needs a
> *repeated* cadence, not another one-off post (funnel died 48h after Jul 6).

## 🚀 Jul 6 2026 — PUBLIC LAUNCH (LinkedIn post live)

Open signup went live (`OPEN_SIGNUP=true` on the VM) and the user posted the
public LinkedIn announcement with the tagged link (`/?src=li`). Same-day ships
that made it possible: mobile pass on the funnel path, open-signup gate +
open-beta homepage, `src` attribution on the `login` event, nightly encrypted
backups, grounding gate passed live QA. Cost guardrails: 10 prep credits/user
(~$3–6 worst case each), auto-reload enabled on the Anthropic account.

**✅ Jul 12 — wave re-read DONE** (ran a day early, user call; the Jul 13
trigger is now moot). Full numbers in FREE_RUN_NOTES (Jul 12). Answers to the
five pre-committed questions:
1. **Wave:** 17 signups Jul 5–8 (peak 8 on Jul 7), 14/15 logging-in users
   tagged `src:li` — then **zero signups Jul 9–12**. One post ≈ 3-day reach.
2. **Activation:** 9/17 (53%) signup → playbook, up from day-1 40%. Prepfirst:
   17 view → 9 submit → 7 generate_ok.
3. **Debrief:** `debrief_prompt_view` fires on all 3 surfaces (6 users) →
   **first real `debrief_save` (1 user)**. The 3b gate is OPEN.
4. **Spend:** 39 credits / 17 users, sane. One user at the 10-credit limit.
5. **Errors:** only the two known pre-fix Jul 7 connection drops; zero since
   the `context.WithoutCancel` fix. Clean.

**Decisions this implies (per the pre-committed tree):**
- **Build Phase 3b** (debriefs > 0): proactive Today prompt + admin
  prep-accuracy stat.
- **Scale distribution** (activation holds, errors clean) — but as a repeated
  cadence, not another one-off post: the funnel was dry within 48h.
- **Watch retention**: 3/17 returned on a later day (all to the playbook).
  The debrief loop is the retention mechanic — 3b serves this too.
- Check the at-limit prep-credit UX before pushing more users in.

Still open (soft): privacy paragraph on the homepage (what's stored, who sees
it) — ship early in the wave window.

## ✅ Jul 6 2026 — grounding gate PASSED → cleared to go wider

Claude-for-Chrome QA on prod, all 8 steps green ("Lusha" test row, cleaned up):

- [x] **Right company:** "Researched: Lusha Systems · lusha.com" — correct
      B2B sales-intelligence Lusha, not a lookalike. The 365scores failure
      class is closed on a live generation.
- [x] **Real citations:** Sources were specific deep links (TechCrunch,
      Calcalist, Lusha blog, Glassdoor…); sampled TechCrunch link genuinely
      substantiated the brief's $1.5B-valuation claim.
- [x] **Debrief surfacing (shipped earlier today):** card on round tab, cross-tab
      banner → form jump, save → "Debriefed" badge — all verified on prod.

**→ The gate is open: start the wider outreach** (pending list is empty, so
wider = LinkedIn/LayoffRadar push + new invites).

- [ ] After going wider: watch `debrief_prompt_view` (`surface` prop) →
      `debrief_save` convert on the new surfacing. When `debrief_save` > 0 for
      real users → build 3b (proactive Today prompt + admin prep-accuracy stat).

## ✅ Jul 6 2026 — re-engagement cohort re-read done → amber-green, one gate left

Pulled adoption/invite-funnel/users via the admin console. Funnel:
**50 invited → 32 signed in → 21 activated → 10 active** (Jun 29 baseline:
48→26→15→6). +6 sign-ins, +6 activated, active up 6→10 — the email + reframe
moved people (11 invitees logged in since it went out; clear Jul 5 cluster).

The four pre-committed questions:
1. **Came back?** Yes, modestly — 11 logins since the email, 6 brand-new sign-ins.
2. **Reframe convert?** Inconclusive at n=6 prepfirst assignments (4 submit,
   3 generate-ok, 3 skip). The two `prepfirst_generate_error`s were Jun 30 +
   Jul 1 only (launch-day hiccups, empty props — no reason logged); Jul 5 cohort
   generated clean. Not a live bug.
3. **Grounding hold?** Unproven either way — no wrong-company signal, but no
   positive confirmation. Hence the gate above.
4. **Anyone debrief?** **Zero — and `debrief_prompt_view` was zero too**, i.e.
   the prompt was never even SEEN. Root cause found (not a demand signal): the
   debrief card only rendered on a past round's tab, but the page defaulted to
   `nextRound ?? company` — exactly when all rounds are past (the debrief case),
   it opened on Company and the prompt was structurally unreachable. Plus the
   event fired on click, not render. **Fixed + shipped Jul 6:** default tab
   prefers an un-debriefed past round when nothing's upcoming; slim cross-tab
   banner ("How did the {round} go? Debrief →") when the pending round isn't the
   selected tab; `debrief_prompt_view` now fires on render (once per round per
   visit, `surface: round_tab|banner|stage_done`); prepfirst error event now
   carries `step`+`reason`, and prepfirst retry no longer creates duplicate
   applications. Verified rendered: Playwright against local stack — both
   surfaces + banner→form jump + save → events land with right props.

**Decision:** go wider once the grounding gate above passes. Debrief loop was a
surfacing bug, not a retention verdict — re-judge it with the new data.

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
- [x] Jul 6 re-read: still 0, and prompt-views were 0 — the prompt itself was
      unreachable (default-tab logic). Fixed Jul 6, see the top entry.

*Unblock — one-tap rounds — ✅ SHIPPED Jul (deploy #79, migration 0021):*
- [x] `interviews.scheduled` flag — a round no longer needs a date.
- [x] "+ Add round" chip picker + stage-done hook ("Just did the {stage} round?
      Debrief it →") so rounds get created from behavior users already do.
- [ ] Watch: rounds created per activated user, and `debrief_save` finally > 0.
      If debriefs now happen → build 3b (proactive Today prompt + admin prep-accuracy).

*3b — proactive + metrics — ✅ SHIPPED Jul 12 (PR #36, deploy #92 green; user
verified the new mobile board on a real iPhone against prod — "looks perfect"):*
- [x] proactive Today prompt: past/undated un-debriefed round → banner "How did
      the {round} at {company} go? → Debrief", deep-links to `?debrief=<iid>`
      with the form open. Fires `debrief_prompt_view {surface:'today'}`.
- [x] admin Adoption: **prep-accuracy stat** (% spot-on headline + spot_on/partly/off
      bars) — first real trust metric.
- [x] "what we learned" on Today: right-rail block with the 2 most recent debriefs
      (feel · accuracy · topics), links to the app. (The old vague "Pursuit debrief"
      card no longer existed — this is its real-meaning replacement.)
- Also shipped alongside: **mobile pass** (board = stacked full-width lanes ≤720px
  + per-card "Move to" select — HTML5 drag never fires on touch, so phones had NO
  status-change path; edit-modal capped scroll box so Save stays reachable with
  the keyboard up; sticky pipeline-editor Save; datalist dropped on touch;
  `viewport-fit=cover`) and **prep-credit cap UX** (`/api/me` exposes credits;
  cap note with mailto replaces Generate at limit; "N left" hint at ≤2;
  `feedback_click {surface:'credit_cap'}`). All verified rendered via local
  Playwright (20/20 green, desktop + 390px viewports).

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
- [x] Backup story for Postgres — nightly workflow (Jul 6 2026): `pg_dump` on the
      VM (last 7 rotated locally) + encrypted off-VM copy as a GH Actions artifact
      (30-day retention; repo is public so plaintext never leaves the VM). Restore
      runbook in `deploy/README.md`. Optional hardening: set a `BACKUP_PASSPHRASE`
      secret (else key derives from `DEPLOY_SSH_KEY`).
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
