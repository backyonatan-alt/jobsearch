# Pursuit — current design snapshot

`pursuit-current.html` is a **single, self-contained** snapshot of Pursuit's
current production UI. It renders by just opening the file in a browser — no
build, no backend, no `pnpm dev`. React 18 + Babel standalone load from a CDN;
all design tokens, screen CSS, and mock data are inline.

## What's in it

A left-nav switches between the five production surfaces, all populated with the
real preview-mode fixtures (Anthropic / Stripe / Linear / Vercel / Figma /
Notion / etc. and the baked Stripe dossier for Sarah Chen):

- **Today (interview)** — editorial Brief (date, greeting + stat strip, "PREP FOR
  TODAY", baby-blue insight box, Worth reviewing, Two quick tips, meta, CTA,
  Later this week) + the right Pulse pane.
- **Today (no interview)** — lede + "What you can do today" hairline rows +
  Recently added + same Pulse pane.
- **Board** — the 5-column kanban with status pills, stale red-border cards,
  empty "Drop here" lanes.
- **Detail** — the B v2 detail page: interview-prep section (AI tips, company
  brief, likely interviewer, lands/avoid, signals, questions, activity log) +
  the right rail (Next interview, Details, Contact).
- **Insights** — KPIs, pipeline funnel with conversion %, weekly activity bars,
  "Where they come from", insight callout.

Buttons/links are visual only — this is a design artifact, not the live app.

## How to use it in claude.ai

1. Open a new chat at [claude.ai](https://claude.ai).
2. **Upload `pursuit-current.html`** (paperclip / drag-drop).
3. Ask Claude to **render it as an Artifact** and iterate, e.g.:
   > "Render this HTML as an Artifact. This is the current design of my app
   > Pursuit. Keep all five views working via the left nav. I want to try
   > [your change] — show me the new version in the Artifact."
4. Iterate on copy, spacing, color, layout, new states, etc. inside the
   Artifact until it looks right.
5. **Bring changes back**: copy the relevant CSS/markup deltas out of the
   Artifact and wire them into the real Svelte app:
   - Design tokens live in `web/src/lib/design-system.css` (`:root` + shell).
   - Per-screen markup/CSS lives in the matching route:
     - Today → `web/src/routes/app/+page.svelte`
     - Board → `web/src/routes/app/board/+page.svelte`
     - Detail → `web/src/routes/app/[id]/+page.svelte`
     - Insights → `web/src/routes/app/funnel/+page.svelte`
     - Shell → `web/src/routes/app/+layout.svelte`
   - Then follow the normal ship playbook (local preview before deploy — see
     CLAUDE.md).

## Fidelity notes

- This is a faithful **visual** snapshot of what's in production now. It's the
  starting point for iteration, not a second source of truth — the Svelte
  routes remain authoritative.
- The HTML mirrors the original handoff format in
  `design_handoff_pursuit_today/` (React 18 + Babel via CDN, one `text/babel`
  script, one `<style>` block), collapsed into a single file for Artifact use.
