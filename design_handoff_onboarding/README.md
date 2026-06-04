# Handoff: Pursuit — First-Run Onboarding (Coachmark Tour)

## Overview
A first-run onboarding for Pursuit (a job-search command center). The first time a new
user opens the app, a friendly **coachmark tour** introduces what Pursuit does and ends by
having the user create their first application. It is a guided, *non-blocking* overlay: the
app stays visible and bright underneath; a small speech bubble points at one element at a
time and the user clicks **Next** to advance.

Flow (7 stages):
1. **Welcome** — centered card with a light scrim ("Welcome to Pursuit 👋").
2. **Today / morning briefing** — highlights the stats cluster.
3. **AI prep** — highlights the interview-prep insight card.
4. **Board** — switches to the Board view, highlights the pipeline columns.
5. **Insights** — switches to the Insights view, highlights the funnel.
6. **Your turn** — back to Today, highlights the "New application" button.
7. **Add an application** — opens the New Application modal and highlights the paste field; the primary button is "Add it".
On finish, a "You're all set" confirmation toast appears and auto-dismisses.

## About the Design Files
The files in this bundle are **design references created in HTML/React (via in-browser
Babel)** — a working prototype showing the intended look and behavior. They are **not meant
to be shipped as-is.** The task is to **recreate this onboarding in Pursuit's real codebase**
using its existing framework, component library, router, and design tokens. If the target app
has an onboarding/tour mechanism already, prefer it; otherwise the structure below maps cleanly
to a small custom component or a library like Shepherd.js / Driver.js / Intro.js.

The prototype is self-contained:
- `Pursuit Guided Tour.html` — the mock Pursuit app (sidebar, Today, Board, Insights, Detail) **plus** the onboarding mounted on top. The app shell here is only a stand-in so the tour has something to point at; **don't port the shell** — point the tour at your real app.
- `pursuit-tour.jsx` — the onboarding logic + components (this is the part to recreate).
- `pursuit-tour.css` — the onboarding styles.

## Fidelity
**High-fidelity.** Colors, typography, spacing, radii, shadows, copy, and interactions are all
final. Recreate the bubble, ring, modal, and confirmation pixel-accurately using your codebase's
primitives. The underlying app screens in the HTML are placeholders — ignore their styling.

---

## How it works (architecture)

The tour is a **single overlay component** that:
1. Holds a **step index** (`i`). Step 0 is the welcome; steps 1…N are real stops.
2. On each step, optionally **navigates the app to the right view** (Today / Board / Insights) so the target element exists.
3. **Measures the target element** in the live DOM via `getBoundingClientRect()` and positions a "ring" highlight + a speech "bubble" relative to it.
4. On the last step, **renders the New Application modal** and points at the paste field inside it.
5. On finish, shows a confirmation toast, then unmounts.

### Target anchoring — `data-tour` attributes
The tour locates elements by a `data-tour` attribute (not by fragile class selectors). Your app
must stamp these on the real elements. The prototype uses exactly these:

| `data-tour` value | Element it marks | View |
|---|---|---|
| `stats` | the in-progress / awaiting / gone-quiet stat cluster on Today | `today` |
| `prep` | the "Prep for today" AI insight card on Today | `today` |
| `board` | the pipeline columns container on Board | `board` |
| `funnel` | the pipeline-funnel chart on Insights | `insights` |
| `new-app` | the "New application" button (top bar) | `today` |
| `paste` | the dashed "paste a job URL / drop a screenshot" dropzone **inside the New Application modal** | modal |

### Step model
Each stop is a plain object (see `STEPS` in `pursuit-tour.jsx`):
```
{ sel: '[data-tour="stats"]',   // CSS selector for the target
  view: 'today',                // app view that must be active for the target to exist
  icon: <…/>,                   // 16px line icon shown in the bubble heading chip
  kick: 'Today',                // not rendered in final bubble (kept for reference)
  title: 'Your morning briefing',
  body:  '…',
  place: 'bottom',              // preferred side: 'bottom' | 'top' | 'left' | 'right'
  last:  true }                 // only on the final (modal) stop
```
The last stop also carries `modal: true`, which tells the component to render the New
Application modal and anchor to `[data-tour="paste"]`.

### Positioning logic (recreate faithfully)
- `place(rect, side, cardW, cardH)` computes the bubble's top/left for the preferred side, then **flips** (bottom→top, right→left) and **clamps** to the viewport (14px margin) so the bubble never goes off-screen. It returns the *actual* side used, which selects the tail direction.
- The **ring** is the target rect inflated by 6px on each side.
- The **beacon** dot sits at the ring's top-right corner.
- The **tail** (little diamond on the bubble) is offset along the bubble edge to point at the target's center (`tailOffset`).
- **Important:** the ring uses **no CSS position transition** — it snaps to each measured box. (An earlier version animated position and could appear a step behind if the browser tab was backgrounded mid-transition. Keep it snapping, or only animate if you can guarantee completion.)

### Measurement timing (why there's a small dance)
When a step changes the view (e.g. Today→Board) or opens the modal, the target isn't in the DOM
on the same frame. The prototype handles this robustly:
- A **single** persistent `resize`/`scroll` listener calls a `measure()` that reads the **current** step from a ref (so there are never stale per-step listeners fighting over the value).
- On each step change it re-measures via `requestAnimationFrame` **and** a few timed passes (`40, 130, 260, 420, 640 ms`) to catch the view/modal finishing paint.
- `measure()` **does not clear** the highlight when the target is briefly missing — it keeps the previous box until the new one appears, avoiding flicker.

In a production app with a real router, prefer awaiting navigation/layout (e.g. `useLayoutEffect`
after route change, or a ResizeObserver on the target) instead of timed passes.

---

## Components & exact specs

All colors below are Pursuit's existing tokens (see **Design Tokens**). The bubble/modal use the
app's surface/ink/rule/accent tokens — reuse your own equivalents.

### 1. Welcome card (`.cmA-intro`, inside `.cmA-introwrap`)
- **Placement:** centered in the viewport. The wrapper is a full-screen overlay with a light scrim `rgba(10,10,13,0.30)`, `display:grid; place-items:center;`.
- **Card:** width `392px` (max 100%), background `--card` (#ffffff), border `1px solid --rule`, radius `18px`, padding `22px`, shadow `0 30px 70px -20px rgba(10,10,13,.40)`. Entrance: translateY(10px)+scale(.985)→none, .3s.
- **Row:** 38px circular avatar (`--accent` bg, white spark icon) + text block (13px gap).
  - Title: `b`, 14.5px / 600, color --ink — "Welcome to Pursuit 👋"
  - Body: 13px / line-height 1.55, color --ink-2 — "Give me a minute and I'll show you the handful of things that make your search easier — then help you add your first application."
- **Footer (right-aligned):** ghost text button "No thanks" (dismisses) + dark pill CTA "Show me →" (`--ink` bg, white).

### 2. Step bubble (`.cmA-bubble` — a `.t-card`)
- **Card:** width `314px`, background --card, border `1px solid --rule`, radius `14px`, padding `17px 18px 15px`, shadow `0 18px 50px -14px rgba(10,10,13,.30), 0 4px 12px rgba(10,10,13,.06)`. Entrance: translateY(10px)+scale(.985)→none, .24s. Fixed-positioned at the computed top/left; `z-index:120`.
- **Tail:** 14px square rotated 45°, same fill/border as card, on the side facing the target (`tail-top/bottom/left/right`), offset by `--tail`.
- **Header row (space-between):**
  - Left: step counter — accent-colored current number + muted "/ N" (e.g. **1** / 6), 11.5px / 600, tabular-nums.
  - Right: "Skip tour" text button — 12px / 500, color --mute-2, hover --ink-2 on --surface-2, `white-space:nowrap`. Dismisses the tour.
- **Title (`h4`):** 16px / 600, letter-spacing -0.02em, color --ink. Preceded by a 26px rounded chip (`--accent-tint` bg, `--accent-text` icon) holding the step icon.
- **Body (`p`):** 13.3px / line-height 1.58, color --ink-2, `text-wrap: pretty`.
- **Footer (space-between, margin-top 15px):**
  - Left: **progress dots** — N dots, 6px; done = `--accent-tint-2`; current = `--accent`, widened to 17px pill; upcoming = `--rule-strong`.
  - Right: **Back** (ghost: --ink-2 text, 1px --rule border, white bg) shown on **every** stop — on step 1 it returns to the welcome — and **Next** (primary: --accent bg, white) with a → arrow. On the final stop the primary reads **"Add it"** (no arrow).
- Buttons: 12.8px / 500, padding 7px 13px, radius 8px, `white-space:nowrap`.

### 3. Highlight ring + beacon
- **Ring (`.cmA-ring`):** fixed box at target rect + 6px padding, radius 11px, `pointer-events:none`, `z-index:108`. Visual = layered box-shadows that **breathe** (2.4s ease-in-out infinite): from `0 0 0 3px --accent-tint-2, 0 0 0 1.5px --accent` to `0 0 0 6px --accent-tint, 0 0 0 1.5px --accent`. No position transition (snaps).
- **Beacon (`.cmA-beacon`):** 14px, at ring top-right; two stacked accent dots, the top one pulses outward (`0→9px` expanding shadow, fading, 2s infinite), `z-index:109`.

### 4. New Application modal (final stop)
- **Backdrop (`.cm-modalback.shift`):** full-screen, `rgba(10,10,13,0.42)`, `display:grid; place-items:center`, `padding:24px`. The `.shift` modifier adds `padding-right: min(372px, 38vw)` so the centered modal slides **left**, leaving room for the bubble on its right. `z-index:100` (below ring 108 / bubble 120, so the highlight + explanation sit on top).
- **Modal:** width `540px`, --card bg, border 1px --rule, radius 16px, shadow `0 30px 70px -20px rgba(10,10,13,.42)`. Entrance .26s rise.
- **Header:** title "New application" (20px/600) + subtitle "Paste a link or screenshot and Pursuit fills it in — or type it by hand." (13.5px, --mute) + close ✕.
- **Dropzone (`.nm-quick`, this is `[data-tour="paste"]`):** dashed 1px `--rule-strong` border, radius 11px, `--surface` bg, padding 13px 15px. 30px accent-tint chip + spark icon, then "Drop a screenshot or paste a job URL" / sub "Pursuit parses the company, role, and location for you."
- **"or enter by hand" divider**, then a 2-col form: Company* (Anthropic), Role* (Member of Technical Staff), Status (select, Applied), Source (Referral). Inputs 38px, radius 9px, focus ring `0 0 0 3px --accent-tint` + accent border.
- **Footer:** "Cancel" (closes → returns to previous stop) + "Add application" primary (advances → finishes tour). The bubble's "Add it" and this button do the same thing.

### 5. Completion toast (`.cmA-toast`, in `.cmA-toastwrap`)
- **Placement:** bottom-center, `bottom:84px`, centered, `z-index:120`. Entrance: slide up 14px, .34s.
- **Toast:** width `min(440px, 100vw-32px)`, --card bg, 1px --rule, radius 14px, padding ~14px, shadow `0 20px 50px -14px rgba(10,10,13,.32)`. 30px rounded chip (`--positive-tint` bg, `--positive-text` check icon) + text + close ✕.
- **Copy:** bold "You're all set." / "Your applications live here — add another anytime from *New application*". Auto-dismisses after **3.4s** (also closeable).

---

## Interactions & Behavior
- **Advance:** Next/Add it → `i+1`; on the last step → show confirmation toast.
- **Back:** present on every stop; `i-1`, clamped so step 1 → welcome (step 0).
- **Skip tour / No thanks / modal Cancel/✕ / toast ✕:** dismiss the onboarding.
- **View driving:** entering a stop sets the app view (today/board/insights) so the target exists; the tour does the navigating, the user just clicks Next.
- **Modal stop:** opens the New Application modal; the ring lands on the dashed paste field; the modal is shifted left to make room for the bubble at its right.
- **Confirmation:** appears on finish, auto-dismisses after 3.4s.
- **No backdrop dim on the steps** (only the welcome and the modal have a scrim). The app stays usable-looking underneath.
- **Reduced motion:** all entrance animations, the breathe, and the beacon pulse are disabled under `prefers-reduced-motion: reduce`; resting states are fully visible.

## State Management
Minimal:
- `i` — current step index (0 = welcome). In the prototype this lives in the tour component; `setI` drives everything.
- `target` — the measured `{top,left,width,height}` of the current step's element.
- `done` — whether the confirmation toast is showing.
- App-level: `active` (is onboarding showing) and a `view` the tour can set. In production, gate `active` on a "has the user completed/seen onboarding?" flag (e.g. persisted to user settings / localStorage) so it only runs on first run; expose a "Replay" entry point if desired.

> The prototype also has a small **demo control** pinned bottom-center ("FIRST-RUN ONBOARDING · Replay · Dismiss"). That is a harness for reviewing the prototype — **do not ship it.**

## Design Tokens
Pulled from Pursuit's existing token set (defined at the top of the HTML; use your app's equivalents):
- **Surfaces:** `--surface #fbfbfd`, `--surface-2 #f4f4f6`, `--card #ffffff`
- **Ink/text:** `--ink #0a0a0d`, `--ink-2 #2a2a2e`, `--mute #71717a`, `--mute-2 #a1a1ab`
- **Rules:** `--rule #e7e7eb`, `--rule-strong #d4d4d8`
- **Accent (indigo):** `--accent oklch(0.62 0.19 258)`, `--accent-strong oklch(0.55 0.21 258)`, `--accent-text oklch(0.42 0.19 258)`, `--accent-tint oklch(0.96 0.03 258)`, `--accent-tint-2 oklch(0.93 0.05 258)`
- **Positive (success, used by toast):** `--positive`, `--positive-text`, `--positive-tint` (green, hue 152)
- **Type:** Geist (sans) / Geist Mono. Body 14px.
- **Radii:** bubble/modal 14–18px, chips 7–10px, dots/pills 99px.
- **Shadows:** see each component above.
- **Spacing:** card padding 17–22px; viewport clamp margin 14px; bubble gap to target 16px; ring inflation 6px.

## Assets
No raster assets. All icons are inline SVGs (spark/sparkle for the brand mark + AI prep, plus
today/board/insights/add/check/arrow/×). Recreate with your icon set; the spark is Pursuit's
brand mark and appears as the welcome avatar, the AI-prep icon, and inside accent chips.

## Files
- `Pursuit Guided Tour.html` — full prototype (mock app shell + onboarding). Run it to see the exact behavior. The shell is a placeholder; only the onboarding is in scope.
- `pursuit-tour.jsx` — the onboarding component to recreate (`CoachBubble`, the `useCoach` engine, `AddAppModal`, step model, positioning helpers). Ignore the `Harness` (demo control).
- `pursuit-tour.css` — onboarding styles (classes prefixed `cmA-` for the bubble/ring/beacon/intro/toast, `nm-`/`nux-modal`/`cm-modalback` for the modal, `t-` for shared button/dot primitives).
