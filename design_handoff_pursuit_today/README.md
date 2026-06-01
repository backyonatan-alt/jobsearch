# Handoff: Pursuit — Today, Board, Playbook, Detail & Insights

## Overview
Pursuit is a calm job-search tracker. This bundle covers the **activated "Option B" direction** — a
two-pane **Today** screen (an editorial morning *Brief* on the left, a *pipeline pulse* on the right)
plus the four screens it links to: a drag-and-drop **Board**, an **Application detail** view, the full
**Interview playbook**, and an **Insights** dashboard.

The product premise: Pursuit *advises* the user (what to prep, who's gone quiet) but **never sends
messages on their behalf**. There is intentionally no in-app "nudge/compose/send" feature — quiet
applications surface a gentle "it might be a good time to reach out directly" prompt and a
"Log a follow-up" action that only records that the user reached out themselves.

## About the Design Files
The files in this bundle are **design references built in HTML/React (via in-browser Babel)** — a
working prototype that shows the intended look and behavior. They are **not production code to copy
directly**. The task is to **recreate these designs inside your existing project**, using its
established framework, component library, routing, state, and design-token system. If the project has
no front-end environment yet, pick the most appropriate framework and implement the designs there.

The prototype splits logic across files purely for the Babel sandbox; in a real codebase you'd model
these as proper components/routes.

## Fidelity
**High-fidelity.** Final colors, typography, spacing, and interactions. Recreate the UI faithfully
using your codebase's libraries and patterns. Exact token values are listed under **Design Tokens**.

---

## Design Tokens

### Color (light theme)
| Token | Value | Use |
|---|---|---|
| `--surface` | `#fbfbfd` | app background |
| `--surface-2` | `#f4f4f6` | hover fills, chips, board lane base |
| `--card` | `#ffffff` | cards, panels |
| `--ink` | `#0a0a0d` | primary text, dark buttons |
| `--ink-2` | `#2a2a2e` | secondary text |
| `--mute` | `#71717a` | muted text |
| `--mute-2` | `#a1a1ab` | faint labels |
| `--rule` | `#e7e7eb` | hairline borders |
| `--rule-strong` | `#d4d4d8` | stronger borders |
| `--accent` | `oklch(0.62 0.19 258)` | blue — Screen status, links, AI |
| `--accent-strong` | `oklch(0.55 0.21 258)` | |
| `--accent-text` | `oklch(0.42 0.19 258)` | accent text on tint |
| `--accent-tint` | `oklch(0.96 0.03 258)` | accent backgrounds |
| `--accent-tint-2` | `oklch(0.93 0.05 258)` | |
| `--positive` | `oklch(0.65 0.14 152)` | green — Offer, "on track" dot |
| `--positive-text` | `oklch(0.4 0.14 152)` | |
| `--positive-tint` | `oklch(0.96 0.05 152)` | |
| `--warm` | `oklch(0.7 0.16 50)` | orange — Interview, "needs attention" |
| `--warm-text` | `oklch(0.5 0.18 50)` | |
| `--warm-tint` | `oklch(0.96 0.06 50)` | |
| `--danger` | `oklch(0.6 0.18 22)` | red — stale (7+ days), rejected |
| `--danger-text` | `oklch(0.5 0.18 22)` | |

**Status → color mapping (used app-wide, keep consistent):**
Wishlist = mute-2/grey · Applied = mute/grey · **Screen = accent (blue)** ·
**Interview = warm (orange)** · **Offer = positive (green)** · Rejected/closed = danger (red).

### Typography
- **Sans:** `"Geist"`, fallback `-apple-system, system-ui, sans-serif`. Weights 300/400/500/600/700.
- **Mono:** `"Geist Mono"`, fallback `ui-monospace, "SF Mono", monospace`. Used sparingly: card
  "time ago", small numeric metadata.
- Global: `-webkit-font-smoothing: antialiased; font-feature-settings: "ss01","cv11";
  letter-spacing: -0.003em;`
- Headings use negative tracking (`-0.03em` to `-0.035em`) and weight 300–500 (light, editorial).
- **Never** use monospace for human-readable dates (a prior reviewer note) — dates render in sans.

### Shadows / radius / spacing
- `--sh-1: 0 1px 0 rgba(10,10,13,0.025)` (subtle lift)
- `--sh-2: 0 1px 2px rgba(10,10,13,0.04), 0 0 0 1px var(--rule)`
- `--sh-pop: 0 8px 24px -8px rgba(10,10,13,0.12), 0 2px 4px rgba(10,10,13,0.04)` (hover/elevated)
- Radii: cards 12px · panels 14–16px · pills 999px · buttons 7–11px.
- Card padding ~13–24px; screen padding ~30–40px; column gap 14px; card gap 10px.

---

## App Shell (persistent chrome)
Grid: `216px sidebar | 1fr content`. Full viewport height; content column scrolls.

- **Sidebar (216px):** brand mark + "Pursuit" wordmark; nav items **Today / Board (count 12) /
  Insights**; active item has a solid `--ink` background, white text. Footer: user avatar (initials
  "YB"), name, email, chevron. Nav items are clickable and switch the top-level view.
- **Topbar (50px):** left = current view title; for drill-in views (detail, playbook) it becomes
  `‹ Back / <crumb>`. Right = a search field placeholder (`⌘K`) and a primary **New application**
  button (`⌘N`). Search/New are visual only in the prototype.
- The brand mark: 22px rounded square, 1.5px accent border, a 7px accent dot centered + a 5px accent
  dot top-right (mirrors an "orbit" motif).

---

## Screens / Views

### 1. Today  (`TodayScreen`)
Two-column grid filling the content area: **left = Brief (1.08fr, scrolls), right = Pulse (0.92fr,
scrolls)**, divided by a 1px rule. (This is the *swapped* layout — the warm Brief leads on the left;
the pipeline status supports on the right.)

**LEFT — "The Brief" (editorial):**
- Quiet sans date line: `Wednesday · 20 May 2026`.
- Big light greeting: `Good morning,` / `Yonatan.` (h1, weight 300 with the name in 500).
- Lede naming the day's interview: *"Today it's your **final round at Anthropic** — the ML Engineer
  role, one-on-one with Dario Amodei. Two more screens follow later this week."* ("final round at
  Anthropic" is warm-colored).
- Section kicker **BEFORE THE ROOM** (uppercase, tracked, with a trailing hairline).
- **Insight box** (accent-tint background, 13px radius, spark icon): a *sourced* inference —
  *"Going on his essays and recent interviews, Dario reasons from first principles and seems more
  interested in **how you think** than what you've shipped — so narrate your reasoning, not just the
  result."* — Framing matters: claims must read as derived from public material, never insider intel.
- **Worth reviewing** — a 3-item checklist (accent check bullets).
- **Two quick tips** — each a spark icon + one line of grounded advice.
- Meta line: `60 min · Google Meet · Final round` (dot separators).
- Primary CTA `Open the full playbook →` (solid ink) → navigates to **Interview playbook**.
- **Later this week** agenda: 3 rows `[ Day time | Company · role | status pill ]`, each clickable →
  Application detail.
- Footer link: `12 applications tracked · open the board →` → Board.

**RIGHT — "Where things stand" (pulse):**
- Kicker **WHERE THINGS STAND** with a warm dot.
- **Stat card** — one bordered card split into 3 columns by hairlines: `12 Active loops`,
  `4 Awaiting reply`, `2 Gone quiet` (last value warm). Each cell hover-fills; clicking opens the Board.
- **Waiting to hear back** list (`longest first`): rows `[ logo | company + stage | "Nd" | dot ]`.
  Active rows show a **green** dot; quiet rows (7+ days) show a **warm** dot. Row click → detail.
  *(No "Nudge" button — removed by design.)*
- **Your move** — a personal task checklist (label + count "N to do"). Rows: checkbox + title +
  subtitle + due chip (warm if urgent). Toggling a row checks it off (local state); checked rows show
  strikethrough + green check. One task starts pre-completed.
- **Footer advisory** (card, spark icon): *"Figma and Notion have gone quiet — No reply in over a
  week — it might be a good time to reach out to them directly."* + a quiet `See both →` link → Board.

### 2. Board  (`BoardScreen`) — drag-and-drop kanban
- Header: sans date; big `Board.` title (34px, -0.035em); subline `12 in flight · drag a card across
  columns to move its status.` + a legend `● red border = no movement in 7+ days` (red dot).
- **5 columns**, `grid-template-columns: repeat(5, minmax(0, 1fr))`, gap 14px. Each column is a
  **filled lane**: background `oklch(0.975 0.003 255)`, 1px `--rule` border, 14px radius, min-height
  140px. Cards (white) sit inside and pop against the lane.
- **Column header:** a tinted **pill** (`dot + label`) colored by status, then the count, then a `+`
  button pushed to the right.
- **Card:** white, 12px radius, `--sh-1`, `cursor: grab`. Top row = 20px logo + bold company name
  (ellipsis on overflow). Role line (muted). Optional **schedule pill** colored by status showing the
  next event time only (e.g. `Thu · 11:00`, `Today · 14:00`, `Fri · 17:00`). Footer (top hairline):
  source on the left (ellipsis), **time-ago** on the right.
- **Stale state:** cards with `days >= 7` (and not wishlist) get a **red border + red time-ago**.
- Card click → Application detail.
- **Drag to move:** cards are `draggable`. On `dragstart`, set `dataTransfer` `text/plain` to the
  app id (`effectAllowed = "move"`). Columns are drop targets: `onDragOver` calls `preventDefault()`
  and highlights the lane (accent-tint + accent border); `onDrop` reads the id and reassigns that
  app's status to the column. A moved card is marked "moved" → its time-ago becomes `just now` and the
  stale border clears. (In the prototype, moves live in local component state — wire to your real
  mutation/API.)

### 3. Application detail  (`DetailScreen`)
- Header: 56px logo + company (24px) + role + meta row (`Applied`, `Source`, `CV`), status pill on
  the right.
- Two-column grid `1fr | 320px`.
- **Main:**
  - **Next-step card** — for apps with an upcoming event: a solid **ink** card, warm-dot kicker
    (`Next step · <relative>`), title `<label> · <when>`, meta (`who`, `medium`, `mins`), and actions.
    For the Anthropic interview the primary action is `Open the playbook →`; otherwise
    `Add to calendar` + `Reschedule` (ghost).
  - For apps **waiting/quiet** (no next event): a muted `--surface-2` card instead — kicker
    `Gone quiet`/`Waiting to hear back`, the stage, `No reply in N days — it might be a good time to
    reach out to them directly.`, and a single **`Log a follow-up →`** action (records a follow-up,
    fires the toast; **does not** send anything).
  - **Activity** timeline — vertical line with dots; each row = mono date + title + note. Dot color by
    tag (accent/positive/offer/danger).
- **Side cards:** Contact (recruiter avatar/name/title, or "No contact yet"); Details (status, last
  activity, source, résumé); Actions (`Log a follow-up` for quiet apps, `Add a note`, `Log an event`,
  `Open job post` — list of left-aligned ghost buttons with leading icons).

### 4. Interview playbook  (`DossierScreen`)
The deep prep doc the Brief teases. Header rule: `Interview playbook` (left) · `✦ Generated by
Pursuit · 12 min ago` (right). Two-column `340px | 1fr`.
- **Left rail (sticky):** person card (76px gradient avatar with initials, name, role, prior roles,
  link chips) · facts card (Company / Role / When / Duration / Where / Round) · `Add to calendar` CTA.
- **Main:** light h1 `Before you meet <Name>.` + dek · **Snapshot** paragraph · **How he interviews**
  (lead paragraph + 3 "tells" cards: Pace / Energy / Depth) · **Lands & lands flat** (two lists: green
  checks vs red ×) · **Recent signals** (sourced cards: date + kind tag + body + source domain) ·
  **Questions worth asking** (cards: the question + a spark-prefixed "why").
- Note: the "How he interviews" copy is currently asserted; if you want it as defensible as the Brief,
  source it to the public signals too.

### 5. Insights  (`InsightsScreen`)
- Title `Insights` + subline.
- **3 KPI cards:** Reply rate (`67%`, up delta), Avg. time to first reply (`4.2 days`), Furthest stage
  (`Offer ×1`).
- Grid `1.3fr | 1fr`:
  - Left panel: **Pipeline funnel** (horizontal bars per stage, colored by status) + **Application
    activity** (12 weekly bars, last bar accent) with month ticks.
  - Right panel: **Where they come from** (source rows with mini progress tracks + counts) + an
    accent-tint insight callout.

---

## Interactions & Behavior
- **Navigation:** sidebar switches top-level views (Today/Board/Insights) and clears the back stack;
  in-content links push onto a stack so the topbar `‹ Back` returns. Scroll resets to top on nav.
- **Today:** stat cells, pulse rows, agenda rows, footer links all navigate. Task rows toggle complete
  (local). The playbook CTA → playbook.
- **Board:** full HTML5 drag-and-drop between columns (see screen 2). Hover lifts cards to `--sh-pop`.
  Lane highlights while dragging over. `+` and search are placeholders.
- **Follow-up:** "Log a follow-up" (detail page, side actions) shows a confirmation **toast**
  (`Follow-up logged — we've reset the clock`, bottom-center, auto-dismiss ~2.6s). No message is sent.
- **Transitions:** subtle (120–160ms) on hover fills, shadows, lane highlight. Toast/scrim use small
  fade/rise keyframes.
- No loading/error/empty states beyond a board column "Drop here" placeholder when a lane is empty.

## State Management
- `view` `{ name, params }` + a history `stack` for back; `active` for sidebar highlight.
- Board: `moves` map (`appId → newStatus`) applied over the base data; `over` (column being dragged
  over). In production these become real status mutations.
- Today: `tasks` array with `done` flags (local toggle).
- App-level: `toast` string (transient).
- **Data shape** (`today-data.js`): `USER`, `NOW`, `STATUS_LABEL`, `APPS[]` (id, co, short, logo class,
  role, status, stage, `days` since last activity, `quiet`, `hot`, source, applied, cv, recruiter,
  `next` {kind,label,when,in,who,medium,mins}, `timeline[]`), `DOSSIER` (the Anthropic playbook),
  `INSIGHTS` (funnel/activity/sources). Use this as the schema for your models.

## Assets
- **Fonts:** Geist + Geist Mono (Google Fonts) — swap for your app's font system if different.
- **Logos:** company logos are **placeholder colored squares with initials** (`.lg-*` classes in
  `options.css` / `today.css`). Replace with real logos/favicons in production.
- **Icons:** inline SVGs defined in `today-shell.jsx` (`ICONS`) + a `Spark` mark for AI moments. Swap
  for your icon set.
- No raster images or brand assets are required.

## Files (in this bundle)
- `Pursuit Today.html` — entry; loads React 18 + Babel, then the files below. Open in a browser to view.
- `options.css` — **design tokens**, app shell (sidebar/topbar), pills, logo colors.
- `left-options.css` — pulse-stage styles (the right pane of Today).
- `orbit-brief.css` — Brief styles (`.ob`, `.brief`, agenda, insight).
- `today.css` — everything for the new screens: Today tweaks, Board kanban, Detail, Playbook,
  Insights, toast.
- `today-data.js` — the data model (see State Management).
- `today-shell.jsx` — `ICONS`, `Spark`, `Logo`, `Pill`, `AppShell`.
- `today-screens.jsx` — `TodayScreen`, `BriefRight`, `BoardScreen`, `InsightsScreen`.
- `today-detail.jsx` — `DossierScreen` (playbook), `DetailScreen`.
- `today-app.jsx` — router/state + mount.

> Recreate these as idiomatic components/routes in your stack; map the tokens above to your theme, the
> data shape to your models, and the HTML5 drag handlers to whatever DnD library your project uses.
