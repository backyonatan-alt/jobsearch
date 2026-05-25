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

- 

### New application modal (⌘N)

- 

### Application detail (/app/[id]) + dossier

- 

### Board (/app/board)

- 

### Funnel (/app/funnel)

- 

### Admin / People (/admin/people)

- 

### Onboarding (sign out, sign back in, or ?onboarding=1)

- 

### Anything else

- 

---

## Shipped (move items here once fixed)

- `[bug]` dossier meeting hero rendered start time in the **server's** TZ while the Scheduled list rendered it in the **browser's** TZ — same event showed two different wall-clock times. Fixed by sending raw `starts_at`/`ends_at` from `meetingDTO` and letting the Svelte component format. (May 25 2026)
