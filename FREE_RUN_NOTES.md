# Free run notes — Pursuit

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

## Shipped (move items here once fixed)

### Bug fixes

- `[bug]` dossier meeting hero rendered start time in the **server's** TZ while the Scheduled list rendered it in the **browser's** TZ — same event showed two different wall-clock times. Fixed by sending raw `starts_at`/`ends_at` from `meetingDTO` and letting the Svelte component format. (May 25 2026)
- `[bug]` Today's Applications table briefly showed **"-1 days ago"** for just-created rows because `Math.floor` on a negative ms diff rounds toward −∞. Fixed `fmtRelativeDate` to treat `d ≤ 0` as "today". (May 26 2026)

### Demo data for new users (May 26 2026)

- New `POST /api/me/demo-seed` + `DELETE /api/me/demo-seed` routes (same handlers as the admin ones, just open to any signed-in user via `requireUser`).
- On Today, when the signed-in user has zero applications, the whole body collapses to a welcome card: brand glyph, one paragraph about Pursuit, "Add your first application" + "Try with demo data" buttons.
- Once seeded, a small "Clear demo data" chip appears in the Applications-table filter row; it deletes only rows whose notes start with `[demo] `, so real applications are kept.

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
