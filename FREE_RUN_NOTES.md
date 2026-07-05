# Free run notes — Pursuit

> ⏳ **PENDING — re-read the re-engagement cohort (~Jul 7 2026).** After shipping
> P0–P3a + one-tap rounds, the user sent a Hebrew re-engagement email to all
> invited users (~Jul 1). Run the dated re-read in `TODO.md` (top entry): pull
> adoption/invite-funnel/users/events via the admin console snippet and answer the
> pre-committed questions — did people come back, did the reframe convert, did
> anyone finally **debrief** (`debrief_save > 0`), does `prep_accuracy` trend
> spot-on? Then decide: open the gate wider vs. fix what surfaces. **Do this before
> starting new feature work.**

A scratchpad for observations while actually using the app. Capture
the moment you notice something; triage later.

## Legend

- `[bug]`   something is broken
- `[ux]`    works but feels wrong
- `[gap]`   feature is missing where I expected it
- `[idea]`  speculative — maybe worth building, maybe not
- `[wow]`   something worked better than I expected (keep doing this)

## How to use this file

1. While exploring the app, jot a one-line note per observation. Tag it.
2. Don't fix or design here — that's for the next coding session.
3. At session start, the next agent reads this file first, you pick what
   to tackle, and the items you've shipped get moved to the bottom under
   "Shipped".

---

## Jun 30 2026 — first external beta tester (Ayelet) — live feedback

> Real beta user, testing on **mobile first then desktop**. Has a real interview
> tomorrow, said she'll put real details in and see how much it helps. Offered
> her phone number for direct contact — high-signal, engaged tester. (More
> screenshots may follow.)

**[wow] The playbook is the wedge — confirmed by a stranger.** "This is excellent.
Maybe the most important value this tool gives, in my eyes." Validates the
dossier→Playbook pivot externally, unprompted.

**[bug][CRITICAL] Wrong-company research.** Gave it a LinkedIn job link; the brief
researched a *different company with the same name* → the whole playbook was
irrelevant. The company brief keys off the bare company name with web_search and
has nothing to disambiguate. Must ground generation in JD URL / company domain /
location. Kills the core value when it misfires. (Case: "365scores" name collision.)
→ also: confirm the entry path — how did a LinkedIn URL reach generation given the
LinkedIn-fetch rejection? (paste prefill → name only?)

**[bug][high] Citations link to the homepage, not the source.** "I'd add the source
the data is based on, not just the general site (the interviewer may ask where you
know this from). Clicking the link opens the site's main page, not the
company-specific source." Same for the "Recent signals" sources — great, "but it
must have the context." Credibility of the wedge depends on this.

**[bug] Counts don't refresh after a mutation.** Added a 2nd application; the
left-side count stayed "1". Same stale-count class as the delete-confirm two-click
bug (now seen twice) → real fix: invalidate/refetch counts after create + delete.

**[bug][mobile] Source dropdown won't open** ("clicking doesn't open List") — the
Source select in the add/edit form. → ✅ FIXED (deploy #80): tappable source
preset chips under the input, Add + Edit.
**[bug][mobile] Save button disappears while editing** job/contact details — every
tap on a field hides the bottom Save button (keyboard / sticky footer). → ✅ FIXED
(deploy #80): modals sized to 100dvh + scroll, action footer made sticky.
**[bug][mobile] Pipeline/stage section cut off**, can't scroll to read or Save.
→ ✅ FIXED (deploy #80): modals 100vh→100dvh; inline pipeline verified full +
scrollable at 390px. (Full responsive/PWA pass = 3c, still deferred.)
**[ux] "Only fits Mac 😅"** — the app is desktop-sized; mobile is rough across the
board. → needs a scope decision: do a real mobile pass now, or tell testers
desktop-only for the beta?

**[ux][desktop] Right rail too wide / sparse.** "Where things stand" takes a big
chunk of screen with little info; let sections resize or rebalance the columns.

**[gap] Need a "position closed / not relevant" status.** Twice, a process she was
doing well in got the req cancelled — that's not a rejection or a withdrawal. Add a
distinct status.

**[idea] Pipeline stages: drag-to-reorder**, not just up/down arrows.

## Jun 30 2026 — prep-first cold start live QA (Claude-for-Chrome)

> Prep-first cold start shipped + QA'd green end-to-end on prod (all 8 steps).
> Real, company-specific playbook in ~40s. Two observations:

- `[bug]` **Delete-confirm needed two clicks once.** Deleting a throwaway app from
  /app/[id]: first "Delete application" confirm appeared not to take (row persisted,
  Today count briefly read 11), second click worked. Most likely the delete *did*
  succeed server-side and Today re-rendered a **stale count** for a beat before
  settling (the dialog/runConfirm wiring looks correct). Not reproduced — watch for
  it; if it recurs, repro then fix (likely a refetch/`goto` race on delete, not the
  confirm button). Low priority.
- `[gap]` Failed prep-first build (pre-fix, empty-role 400) left a throwaway
  "Anthropic Test" row behind. Already handled in code: the error card no longer
  claims "saved" when create itself failed. Noting in case orphaned rows show up
  from other partial-failure paths.

## May 24 2026 — first free run (post screenshot-parse + demo seed)

> Demo data is seeded (15 apps). Logged in as back.yonatan@gmail.com.
> Goal: click through every surface as if I were a new beta user
> walking the product for the first time.

### Today dashboard (/app)

- `[ux]` The "Pursuit debrief — two threads worth your attention" card is
  not interesting. What are the threads? What is an offer? Card copy /
  concept needs a rethink.
- `[ux]` Top stat row "14 / 2 / 1 — offers / wishlist / loops" should
  be **more prominent** (big counts at the top).
- `[ux]` The word **"loop"** is unclear — what does it mean for something
  to be a loop? Rename or explain.
- `[gap]` Want a "what can you do today" section with proactive offers:
  *"get ready for this interview", "learn about this company"* etc.
  → **wants a design preview before we build it.**
- `[idea]` "What we're noticing" (currently on funnel) belongs on the
  main page — that's the whole product. Move it here.

### New application modal (⌘N)

- 

### Application detail (/app/[id]) + dossier

- `[ux]` The word **"dossier"** / "Open dossier" is unclear. Find a
  different term.
- `[gap]` Add a link back to **where the application was taken from**
  (the original job posting URL on the source site).
- `[gap]` Show **information about the job** itself (JD summary).
- `[gap]` Add a short **AI-generated company summary**.
- `[gap]` Show the **company logo**.
- `[gap]` If we have it, link to the **hiring manager's LinkedIn**.

### Board (/app/board)

- `[ux]` **No spacing between cards** — they're crammed together.
- `[ux]` Instead of the first-letter square avatar, **fetch the actual
  company logo**.
- `[ux]` **Dragging looks really bad** — needs a micro-animation /
  better drag affordance.
- `[idea]` Turn a card **red when it's been a long time** since the
  last activity (needs-your-attention signal).
- `[bug]` The **list ↔ board toggle takes you back to the main page**
  and doesn't scroll inside the view. Remove it (or fix it).

### Funnel (/app/funnel)

- `[ux]` **Don't like the chart.** Rethink it — make it look like a
  real product analytics chart (Mixpanel-style funnel viz).
- `[ux]` The page looks **really plain**, needs a design upgrade —
  more modern.
- `[idea]` Move **"what we're noticing"** to the main page.
- `[gap]` **Why do we need the Pipeline view?** It's just the
  applications table filtered by status — not sure it's needed.
  Consider removing or justifying it.

### Admin / People (/admin/people)

- 

### Onboarding (sign out, sign back in, or ?onboarding=1)

- 

### Anything else

- `[bug]` Top-right user avatar shows initials **"BA"** — should show
  the **Google profile picture** from OAuth.

---

## Jun 9 2026 — first external beta feedback (Michal)

> First real user. Walked the product mid-job-search (already has live
> processes), so she hit the "I'm already deep in my search" edges a
> brand-new user wouldn't. Raw notes, lightly triaged.

### Add / edit application — JD ingest

- `[gap]` Wants a **free-text JD paste** path on top of URL + screenshot.
  Her case: applied, waiting for interview, the posting is already down
  from both the company site and LinkedIn, but she has the JD text saved
  (from a chat). No way to get that text in today.
- `[gap]` **We don't store the JD body — only `jd_url`.** Postings get
  taken down mid-process, so the saved URL rots and the description is
  lost. Need a `jd_text` column + persist parsed/pasted JD text.
- `[ux]` Adding a JD to an **existing** application is hard to find, and
  when found it's **URL-only**. Make "add/replace JD" (URL / text /
  screenshot) a first-class action on the detail page.

### Source field

- `[ux]/[gap]` Source is a free-text input; she effectively saw one value
  ("Referral"). Wants a **dropdown of common sources** (LinkedIn,
  Referral, Cold email, Company site, Recruiter reached out, Other) —
  keep free-text as fallback.

### Contacts

- `[gap]` Wants a **recruiter / point-of-contact** field per application —
  the person scheduling interviews and relaying answers. Distinct from
  `hiring_manager_*` (the role's manager) which is all we have today.

### Pipeline / funnel

- `[idea]` **Per-application customizable pipeline.** After a recruiter
  call she knows the exact stages (direct-manager interview → take-home →
  HRBP → manager-X → offer). Wants to define/edit those per role, and
  have them mutate mid-process. Hard; needs a design.
- `[gap]` **Where do rejected apps go in the funnel?** No clear home for
  "got a no" in the funnel view. Need an explicit rejected/closed lane
  and conversion accounting.

### Bulk import

- `[idea]` She already runs an **Excel** of live + past roles. Wants a
  **bulk import** (CSV/paste) instead of one-by-one entry. Irrelevant for
  users who start in-app, valuable for switchers — likely an activation
  lever for the beta cohort (all mid-search).

### Calendar event parse → save

- `[bug]` **Save fails on AI-parsed events.** Parse (text/screenshot via
  Haiku) succeeded, but "Save 1 event" returned an error. Root cause:
  the AI path tags events `source:"ai"`, and `handleInterviewCreate`
  only allowed `ics`/`manual` → 400 "source must be 'ics' or 'manual'".
  **FIXED Jun 9 2026** (allow `ai`). See Shipped.
- `[bug]` Haiku got the **day-of-week wrong** on first parse ("Tue Jun 10"
  for a Wednesday Jun 10 invite); a second parse corrected it. Worth
  pinning a "today" anchor / weekday cross-check in the parse prompt so
  it's deterministic, and surfacing the weekday in the preview so the
  user can sanity-check before saving.

### Trust / latency

- `[gap]` Wants a **privacy line** — what's stored, who sees it — flagged
  salary specifically ("curious who'll actually enter salary info").
- `[ux]` **Interview-prep felt slow** (her gut: >2 min). Sonnet +
  web_search. Add an explicit progress/expected-latency state; consider
  measuring p50/p95 from the `interview_prep` events we already log.

### Process

- `[idea]` Tell first users **how to send feedback** in the invite email.

---

## Shipped (move items here once fixed)

### Per-round prep — shared company brief + per-interviewer briefs (Jun 28 2026)

> Triggered live by the user: a 2nd interview at the same company (new round,
> new interviewer) kept showing the **1st** round's prep, and the Today page
> didn't refresh after adding it. Root cause: the dossier was `UNIQUE` per
> application — structurally only one brief could exist. PR #24, merged to prod,
> full Claude-for-Chrome QA passed (10/10), test row cleaned up.

- `[bug]` **Second round inherited the first round's prep.** Now each interview
  row holds its own interviewer brief; the soonest round is the default tab and
  a round with no prep shows its own generate state — never another round's.
- `[bug]` **Today didn't refresh after adding an interview.** Detail page now
  fires `pursuit:refresh` on add/delete and Today refetches on tab focus.
- `[ux→model]` **"What's true for all rounds vs this interviewer?"** (user's
  call). Split into a **shared company brief** (researched once — what they do,
  direction, the loop, what the team grades for) shown on a **Company** tab
  first, + **per-round interviewer briefs**. Company is generated **once**;
  generating the first round builds it alongside in parallel (one prep credit).
- **How.** `migration 0018` adds `dossiers.interview_id` + two partial uniques
  (one brief per interview, one app-level company slot); non-destructive — old
  dossiers become company-tab content. LLM split into `GenerateCompanyBrief` +
  `GenerateInterviewerBrief`. Dossier GET serves `?scope=company`,
  `?interview_id=N`, and a no-param composed view (next round + company block)
  so Today needed no change. Verified against live Postgres + on prod.
- `[gap]` **Still parked:** the debrief feed-forward ("how did the last round
  go?" → feed into the next round's prep) the user originally described. The
  schema now supports it cleanly; build on this shape when picked up.

### Interview flow — three bugs found via adoption data + QA (Jun 17 2026)

> Triggered by the new admin Adoption view: **0 / 22** signed-in users had
> ever saved an interview, despite the dossier (interview prep) being the
> #1 surface. Investigation + Claude-for-Chrome QA surfaced three distinct
> bugs, all on the path to a saved interview. Fixed and verified end-to-end
> (parse → save → reload-persists → delete).

- `[ux]` **Discoverability dead-end.** The empty "Next interview" rail card
  had no way to add one, and the only entry point was a vaguely-named
  "+ Log an event" buried in the Activity actions. Added a "+ Add interview"
  CTA to the empty card and renamed the Activity button to match. Explains
  why only 2/22 ever fired `interview_parse`.
- `[bug]` **Time parsed in the wrong timezone.** A bare "2:30pm" was read as
  US Eastern (prompt default) → 9:30 PM for an Israel-tz user. Frontend now
  sends the browser IANA timezone; `ParseEvent` anchors "now" in it and
  interprets untimezoned times in it (US Eastern only as last resort).
  Regression test guards offset survival through `parseFlexibleTime`.
- `[bug]` **Interviews could be created but never deleted.** Backend delete +
  a `deleteInterview()` JS fn existed, but the timeline only rendered a
  delete button for follow-ups. Now renders for interview rows too (reuses
  the existing confirm dialog).

### Adjacent: admin analytics shipped alongside (Jun 17 2026)

- **Invite funnel** (`/admin/invites`) and **Adoption** (`/admin/adoption`)
  views — see CLAUDE.md "Adjacent infra shipped". The Adoption view is what
  surfaced the interview-flow hole in the first place.


### Michal feedback — Chunk 1: application capture (Jun 9 2026)

- Migration 0016 adds `jd_text`, `recruiter_name/email/linkedin` to `applications`; backend list/get/create/update plumb all of them.
- `jd_text`: free-text JD field on the Add modal; the AI parse path also captures the pasted body (not a bare URL/screenshot). Detail page shows a collapsible "Saved job description" so the body survives the posting coming down.
- Source: free-text input backed by a `<datalist>` of 7 common values (LinkedIn / Company website / Referral / Recruiter reached out / Cold outreach / Job board / Other) on Add + Edit. Free text still works. Shared list in `app-helpers.js` (`SOURCE_SUGGESTIONS`).
- Recruiter / point-of-contact (name + email + LinkedIn), distinct from hiring manager. Edit modal groups both; detail Contacts card renders recruiter (warm avatar) above hiring manager.

### Michal feedback — Chunk 2: calendar weekday (Jun 9 2026)

- `[bug]` Haiku put the wrong weekday on a year-less invite ("Tue, Jun 10" for a Wed Jun 10 — June 10 is a Tuesday in 2025, a Wednesday in 2026, and the model had no "today" to pick the year). Fixed: `ParseEvent` now passes the current date in the user message and the prompt requires weekday-consistency (a named weekday must match the resolved date) and resolves year-less/relative dates relative to today.
- Preview now shows the full weekday + year computed from the parsed date (never the model's prose) — e.g. "Wednesday, June 10, 2026 · 11:00 AM" — with a "Double-check the day and time before saving" hint.
- `[ux]` "Add an interview" modal was too dense (two side-by-side zones, each with a drop area + textarea + button). Redesigned to one unified box that auto-routes .ics / screenshot / pasted-text, plus a single "Find the event" button — matches the New Application modal.

### Michal feedback — Chunk 3: funnel outcomes (Jun 9 2026)

- `[gap]` "Where do the rejected apps go in the funnel?" — they were folded silently into the cumulative Applied count. Added an Outcomes strip under the funnel (Offer / Rejected / Withdrawn, colored dots) + a one-liner that names where a "no" lands, so rejections have a visible home.

### Michal feedback — Chunk 4: privacy + prep latency (Jun 9 2026)

- `[gap]` "A few words on privacy?" (she flagged salary specifically) — added a privacy microcopy line to the New Application modal and the detail Edit modal ("Private to your account — never shared…"), naming notes/salary explicitly.
- `[ux]` Interview-prep felt slow / looked stuck (>2 min, static "30–60s" copy). The generating state now cycles honest stages (searching → reading → spotting signals → writing) and sets an accurate "1–2 minutes, keep working" estimate.

### Michal feedback — Chunk 5: in-app feedback link (Jun 9 2026)

- `[idea]` "Tell first users how to send feedback." Added a "Send feedback" link in the app sidebar that opens a pre-addressed email (subject + a small template). Fires a first-party `feedback_click` event.

### Michal feedback — Chunk 6: per-app pipeline (Jun 9 2026)

- `[idea]` Per-role customizable interview pipeline. Migration 0017 adds a `pipeline` JSONB array (`[{name, done}]`) to applications; `PUT /api/applications/{id}/pipeline` replaces it (sanitizes: trims, drops empties, caps 30 stages / 80 chars). Detail page rail gets a "Pipeline" card: a clean checkable vertical stepper (tap a node to mark done) with "N of M done", an Edit mode (rename / reorder ↑↓ / remove / add), and an empty-state "Start from a typical loop" seed. Mutable mid-process.

### Michal feedback — Chunk 7: bulk import (Jun 9 2026)

- `[idea]` Load an existing spreadsheet of past/current roles in bulk. New `POST /api/applications/import` (caps 500 rows, skips rows missing company+role, normalizes unknown status → applied, sets nothing it isn't given). `ImportApplications.svelte` modal off the Board: paste TSV/CSV → auto-map columns by header (with manual override dropdowns + sample values) → preview + count → import → "Added N applications". Emits `bulk_import` event.

### Bug fixes

- `[bug]` AI-parsed calendar events (text/screenshot → Haiku) could never be saved: the parse path tags them `source:"ai"` but `handleInterviewCreate` rejected anything but `ics`/`manual` with 400 "source must be 'ics' or 'manual'". Allowed `ai` as a valid source. (Jun 9 2026)
- `[bug]` dossier meeting hero rendered start time in the **server's** TZ while the Scheduled list rendered it in the **browser's** TZ — same event showed two different wall-clock times. Fixed by sending raw `starts_at`/`ends_at` from `meetingDTO` and letting the Svelte component format. (May 25 2026)
- `[bug]` Today's Applications table briefly showed **"-1 days ago"** for just-created rows because `Math.floor` on a negative ms diff rounds toward −∞. Fixed `fmtRelativeDate` to treat `d ≤ 0` as "today". (May 26 2026)

### Mobile pass (May 26 2026)

- Layout collapses to icons-only horizontal nav strip on screens ≤ 720px (sidebar → top bar). Brand on the left, nav items (icons only) in the middle, profile avatar on the right.
- Topbar wraps; search box shrinks to icon + flex placeholder; kbd hints hidden.
- Per-page media queries stack multi-column grids: count cards (Today) → 2-col, action grid → 1-col, applications table sheds role/applied/arrow columns (relative date moves under company); funnel KPIs → 2-col, two-col blocks → 1-col, time-in-stage → 2-col; brief hero stacks, stats → 1-col, signals/approach → 1-col, edit modal goes full-screen.
- Board stays horizontal-scroll on mobile (six columns at 220px each, swipe to see all).

### Hiring manager on the Brief (May 26 2026)

- Migration 0012 adds `hiring_manager_name` and `hiring_manager_linkedin` to `applications`.
- Backend list/get/create/update plumb both fields. PATCH treats blank as keep.
- Edit modal on the Brief gets two new inputs.
- Brief tab renders a dedicated "Hiring manager" card (warm-tinted initials avatar + LinkedIn button) above the dossier section, always visible when the data is set — independent of whether a dossier has been generated.

### Demo data for new users (May 26 2026)

- New `POST /api/me/demo-seed` + `DELETE /api/me/demo-seed` routes (same handlers as the admin ones, just open to any signed-in user via `requireUser`).
- On Today, when the signed-in user has zero applications, the whole body collapses to a welcome card: brand glyph, one paragraph about Pursuit, "Add your first application" + "Try with demo data" buttons.
- Once seeded, a small "Clear demo data" chip appears in the Applications-table filter row; it deletes only rows whose notes start with `[demo] `, so real applications are kept.
- The auto-deploy for PR #5 silently failed to trigger; merged a no-op PR to force a fresh push-to-main event and bring the welcome-card code online.

### Design decisions locked (May 25 2026 review session)

- **Logo source**: Google favicon service
  (`https://www.google.com/s2/favicons?sz=128&domain=<domain>`).
  Clearbit's free API was deprecated. Fallback = coloured letter square
  when the favicon returns nothing.
- **Today page**: variant A locked
  (`/preview/redesign/today/a`). Modern sans, narrative count cards
  with colored top ribbon + per-metric subtitle, action grid, "What
  we're noticing" with people / pause / moon icons in tinted squares.
- **Board page**: variant A locked
  (`/preview/redesign/board/a`). Real logos, generous card gaps,
  drag micro-animation, stale = red dot + red applied date. Broken
  list/board segmented toggle removed (sidebar handles nav).
- **Application Brief** (was "Dossier"): current pass at
  `/preview/redesign/brief`. Flat (no gradients), no salary block,
  neutral chips with colored group dots, lands/avoid replaced by
  "How to approach this interview" with check/cross markers.
- **Funnel page**: variant B dashboard locked, chart Option 1 locked
  (`/preview/redesign/funnel/b`). 4 KPI cards on top (overall, best
  CV, avg time-to-offer, in flight), stepped funnel in monochromatic
  blue with numbers inside bars and soft accent-pill conversion %,
  source + CV bar charts side-by-side, time-in-stage cells below.
- **Brand mark**: target-style SVG (concentric circles + offset accent
  dot) paired with sentence-case "Pursuit" wordmark in Geist.

### Wired into the real app (May 25 2026)

- **Sidebar Pipeline section killed** — it was just filtered Today views.
  Layout now shows Today / Board / Funnel only, plus the target-style
  brand mark.
- **Today (/app)**: locked Today A applied. Greeting + 4 count cards
  with colored ribbons + per-metric subtitles derived from real apps,
  AI-suggested action grid (prep for next interview, decide on open
  offer, nudge oldest stale, learn about latest screen), insights row
  with people/pause/moon icons in tinted squares (referral lift, days
  since last apply, stale loops), applications table with real
  favicons + relative dates + stale tag.
- **Board (/app/board)**: locked Board A applied. Real favicons, 10px
  card gaps, drag micro-animation (`scale(0.99) rotate(-0.5deg)` on
  active), stale = red border + red stale-dot + red applied date,
  six-column layout (rejected + withdrawn collapsed into a single
  Closed column), broken list/board toggle removed.
- **Funnel (/app/funnel)**: locked Funnel B applied. 4 KPI cards
  (overall conversion, best CV variant, avg time to offer, in flight)
  derived from real apps, monochromatic-blue stepped funnel with
  numbers inside bars and soft accent-pill conversion %, source
  breakdown (Referral / LinkedIn / Cold / Other) + CV variant bars
  side-by-side, time-in-stage cells below. "What we're noticing"
  removed (now lives on Today).
- **Application detail (/app/[id])**: locked Brief applied. Hero strip
  with real logo + status pill + JD link, Up-next card (when dossier
  has a meeting), 3 at-a-glance stats (days in pipeline / current
  stage of 4 / match score placeholder), tabs, dossier-driven content
  rendered through the new Brief layout (interviewer card, snapshot,
  background, recent posts & talks signals with per-source favicons,
  "How to approach this interview" with check/cross markers,
  questions). Generate/regenerate flow preserved.

### Helpers added

- `companyDomain(co, jdUrl)`, `faviconUrl(...)`, `daysSince(iso)`,
  `isStale(a)`, `fmtRelativeDate(iso)` in `web/src/lib/app-helpers.js`.
  `toDisplayApp` now also exposes `domain`, `logoSrc`, `appliedRel`,
  `stale`.

### Google profile picture on user avatar (May 25 2026)

- Added `users.picture_url` (migration 0011).
- Google OAuth callback now reads the `picture` claim from the ID token
  and persists it on upsert. Refreshed on every sign-in via
  `COALESCE(NULLIF($3,''), users.picture_url)` so a missing claim
  doesn't clobber an existing value.
- `/api/me` returns `picture_url`. Sidebar avatar renders the `<img>`
  when present, falls back to the initials square when null. Existing
  users will see the picture on their next sign-in.

### Still to design / decide

- **New application modal**: no notes yet, may not need a redesign.
