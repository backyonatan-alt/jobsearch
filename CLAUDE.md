# Job Search Management — Handoff Notes for the New Project

> **How to use this file:** Codebase root. Carries handoff context from Aegis / StrikeRadar / LayoffRadar, plus the live roadmap. Re-confirm and check items off as they complete.

---

## Decisions locked at kickoff (May 22 2026)

- **Audience:** closed beta, invite-only friends. Multi-tenant data model from day 1, but no public signup.
- **Stack:** Go + Postgres + static frontend on a user-owned Hetzner VM. systemd + nginx + GitHub Actions deploy. Same shape as Aegis but without the black-box constraint — `ssh` and `psql` are fair game on this project.
- **LLM provider:** Anthropic (Claude).
- **MVP wedge:** application tracker + LinkedIn/calendar ingest (spine) + AI interviewer dossier (the AI-native moment).
- **Working brand name:** **Pursuit** — repo stays `jobsearch`, brand is a single config string (`BRAND_NAME` env var) so we can swap it before public launch.
- **Domain:** none for the beta. Running on `<public-ip>.nip.io` with a real Let's-Encrypt cert. Real domain decision deferred until after beta validates demand.
- **Auth:** Google OAuth only (decision changed from magic-link on May 22 2026). All closed-beta users have Gmail anyway, one-click is better UX, and it removes the "where's my email" friction. Magic-link code removed.
- **Mail (post-OAuth):** no longer required for auth. Deferred until we need outbound notifications (reminders, weekly review).
- **Deferred (v1.5+):** post-interview recording analysis (privacy/consent story first), CV A/B testing (needs application volume to show signal).

## Roadmap

### v0.1 — Spine (current)
1. Repo scaffold: Go backend, Postgres, static frontend, GH Actions deploy, nginx + systemd
2. Google OAuth sign-in (was magic-link in initial scope; switched May 22 2026)
3. Applications CRUD: company, role, source, status, applied_at, JD url, notes, CV variant ref
4. Per-user data isolation (multi-tenant schema, not bolted on)
5. Kanban + list views on the frontend (placeholder UI — full design pass scheduled before v0.2 features)

### v0.2 — Ingest + Dossier
6. LinkedIn job paste → parse title/company/location/JD → prefill new application
7. Calendar `.ics` upload or paste → create interview event linked to application
8. AI interviewer dossier: name + company (or LinkedIn URL) → Claude brief on background, recent posts/talks, likely style, watch-fors

### v0.3 — Insight
9. Funnel view: applied → screen → onsite → offer, with conversion rates
10. AI weekly review: "what's working, where you're stalling"
11. Reminders (follow-ups, prep tasks) — email + browser

### v1.5+ — Parked, needs decisions
- Post-interview review + recording analysis (privacy/consent first)
- CV A/B with per-variant tracking link or unique submission email
- Public signup, payments

---

## Project vision (initial sketch from the user, May 21 2026)

A personal system for managing a job search. Core features as stated:

- **Application tracker** — who I applied to, when, status.
- **CV version control + A/B testing** — track which CV variant got which response rate, by role/company type.
- **Interview-type learning loop** — recommendations and prep tailored to interview format (behavioral, system design, coding, take-home, etc.). Learn from what went well and badly.
- **Interviewer dossier** — pre-interview briefing on the interviewer (background, prior talks, signals to watch for).
- **Post-interview review** — points to improve, structured debrief.
- **Recording analysis** — transcribe and analyze interview recordings for self-review (delivery, pacing, gaps).
- **Reminders** — follow-ups, prep tasks, deadlines.

Treat this list the same way the Aegis roadmap is treated: pin it, re-confirm and update as items complete, and add a `TODO.md` for time-bound items and an `AB_TESTS.md` for experiments.

---

## Who you're working with — preferences carried from Aegis

### Product philosophy

- **Ship fast, measure with real data.** Default to A/B testing real surfaces (modals, banners, copy variants). Keep an `AB_TESTS.md` docket: active and completed experiments, methodology, lessons.
- **Funnel-first analytics.** Every meaningful surface gets a tracked event (e.g. on Aegis: `lr_modal_view`, `lr_modal_click_open`, `lr_modal_dismiss`). Set up saved GA4 Explorations for cross-surface funnels.
- **Persistent roadmap in `CLAUDE.md`.** Numbered list, current week vs next vs nice-to-have. Re-confirm and check items off as they complete — the user explicitly asks for this.
- **Cross-product growth thinking.** On Aegis: reused the StrikeRadar Telegram list to launch LayoffRadar. Apply the same instinct here — every new surface should feed another.
- **OG cards / social-share artifacts are first-class.** Don't treat them as an afterthought.

### Quality bar — verify at the *rendered* surface, not at the function level

This is the single most important lesson from Aegis. The May 19 2026 incident (two pins with the identical headline shipped because verification stopped at `events.json` instead of the rendered chart) drove a permanent checklist update. Carry the same bar here:

1. **State acceptance criteria as "what the rendered UI must look like" before coding** — and put it in the PR description as a checklist.
2. **Golden-output tests on production-shaped fixtures.** Snapshot the real prod output, commit as a fixture, assert the *full* output. "A field exists" is not the bar.
3. **Branch dry-run before merge.** If the job is idempotent / `workflow_dispatch`-able, run it on the feature branch against real data and inspect the diff before merging. Preview-on-branch beats observe-prod-and-hotfix.
4. **2-minute adversarial pass: "how could this look dumb?"** Enumerate failure modes against the rendered surface — duplicates, off-by-ones, empty states, too-many-items, identical-looking-things — *before* coding.
5. **Post-deploy spot check** for anything touching the data pipeline.
6. **Incidents drive checklist updates.** When something ships broken, update this file so the next change of that shape can't ship the same mistake.

For a job-search app, "rendered surface" means: the application list view, the CV diff view, the interviewer dossier card, the post-interview review form, the reminder notifications. Not just the API response.

### Communication / working style

- **Direct.** Don't over-ask clarifying questions — make a call and proceed. Push back when you disagree; the user prefers being challenged to being humored.
- **Hebrew + English both fine** in conversation; output code/docs in English.
- **Persist context in repo docs**, not just chat. Re-confirm lists when they grow.
- **Short replies.** End-of-turn summary in 1–2 sentences.
- **Don't narrate deliberation** — state results and decisions.
- **Don't suggest infra paths the user doesn't have access to** (see below).

---

## Technical profile

### Comfortable with (use these by default)

- **Go** — primary backend language on Aegis.
- **PostgreSQL** — main store. Comfortable with JSONB columns, migrations, normalized time-series tables.
- **Static frontend** — plain HTML/CSS/JS, no heavy framework on Aegis. Direct DOM, `fetch` API. If a job-search app justifies a framework, pitch it explicitly.
- **Bash / `jq` / `curl`** — daily driver for ops and verification.
- **GitHub Actions** — primary CI/CD and ops automation surface. Used for deploys, diagnostics, scheduled jobs.
- **systemd** services on Linux VMs.
- **nginx** reverse-proxy config (retry logic, restart-tolerance during deploys).
- **Telegram bot broadcasting** — proven distribution channel.
- **GA4 + LinkedIn + X** for distribution and measurement.

### Less familiar / don't assume

- **Heavy SPA frameworks** (React/Vue/Svelte/Next/Nuxt) — not part of the Aegis stack. If introducing one for this project, justify it (e.g., the CV diff view is genuinely too dynamic for vanilla JS).
- **Mobile native** (Swift/Kotlin). If mobile is needed, prefer a PWA.
- **Kubernetes / heavy DevOps.** The stack should fit on a single VM + GitHub Actions; don't over-engineer.
- **Auth flows beyond a simple admin token.** If this app needs multi-user auth, that's a deliberate choice to discuss before building.

### Infrastructure constraint (the big one)

On Aegis, the user does **not** own the GCP project (`strike-radar`). All access is *indirect* through GitHub Actions — no `gcloud`, no Cloud Shell, no SSH, no `psql` from the user's machine. This shaped the toolset:

- **Paste-the-curl workflow** when the sandbox can't reach a host: the agent gives the user a one-liner with `jq` pre-trim, the user runs it locally and pastes back.
- **Admin-token-protected backend endpoints** for ops queries (e.g. `/api/admin/history`).
- **`workflow_dispatch`** for any prod action.

For a *new* user-owned project this constraint may not apply — but **the habit of building around restricted prod access is good practice anyway**. Default to endpoints + GH Actions over assuming SSH; design so the user can debug from a browser + a terminal with `curl`, not from a GCP console.

---

## Default ship playbook for this project

1. New feature → write acceptance criteria as "what the rendered UI must look like" in the PR.
2. Build with a golden-output test or a branch-deployed preview.
3. Branch dry-run on real-shaped data.
4. 2-minute adversarial pass: list the failure modes that would look dumb.
5. Merge → post-deploy spot check.
6. Add a GA4 (or equivalent) event for any user interaction worth measuring.
7. If something ships broken, update this file's checklist so the bug class can't recur.

---

## Open questions to resolve at kickoff

Ask these before scoping the first sprint:

1. **Single-user (just you) or multi-user (other candidates too)?** Drives auth, multi-tenant data model, hosting cost shape, privacy concerns around recordings.
2. **Stack:** stick with Go + Postgres + static frontend, or try something new for this one? (e.g., is the CV diff/A-B view dynamic enough to warrant React/Svelte?)
3. **Hosting:** same hanan-GCP black-box constraint, or a fresh user-owned project (Fly.io, Railway, Render, Vercel, a personal GCP project)? Recording storage and transcription cost shape this.
4. **Recording analysis pipeline:** local Whisper, OpenAI/Anthropic API, or a hosted transcription service? Privacy of interview recordings is a real concern — clarify storage and retention before building.
5. **Interviewer dossier source:** manual entry, scraped LinkedIn, or LLM-summarized from a pasted URL/profile? Each has very different legal/ToS implications.
6. **CV A/B test mechanic:** is the "response rate" signal trackable (e.g., per-CV tracking link / unique submission email) or self-reported?

---

## What *not* to do (lessons from things that went sideways on Aegis)

- **Don't verify only at the JSON/function level for a user-facing change.** The May 19 2026 duplicate-headline incident shipped a green test suite. Render it and look.
- **Don't suggest `gcloud`, Cloud Shell, `psql`, or SSH paths unless the user has confirmed access for this new project.** Default to "what can you do from a browser + a terminal with `curl`?"
- **Don't migrate infrastructure as a side-effect of another task.** The user has explicitly chosen to live with constraints rather than spin up new projects mid-flight. Same instinct should apply here — pick the hosting answer up front, then stop touching it.
- **Don't add features beyond what the task requires.** A bug fix is not a refactor invitation.
- **Don't write multi-paragraph docstrings or comment blocks.** One short line max; default to no comments unless the *why* is non-obvious.
- **Don't create docs unless asked.** This file is the exception because the user explicitly requested a handoff.
