# TODO

Time-bound items. Cross off as completed. Things that don't have a date go in `CLAUDE.md` roadmap, not here.

## ⏳ ~Jun 25 2026 — activation data re-read (do this first when back)

Context: PR #22 shipped `interview_save` instrumentation + locked the Jun-22
decisions (see CLAUDE.md "Data re-read (Jun 22 2026)"). We held the invite gate
2–3 days to let the new events accrue on the current cohort. Now re-read.

**Pull the data** — paste in the DevTools console at `https://178.105.213.124.nip.io/admin`
(signed in as admin):

```js
(async () => {
  const get = u => fetch(u, {credentials:'include'}).then(r => r.json());
  const [ad, fn, users] = await Promise.all([
    get('/api/admin/adoption'), get('/api/admin/invite-funnel'), get('/api/admin/users')]);
  const stages = {}; for (const i of fn.invitees) stages[i.stage]=(stages[i.stage]||0)+1;
  const ev = {}; for (const e of ad.events) ev[e.name]=e;
  console.log('milestones:', ad.milestones.map(m=>`${m.label}: ${m.users}`).join(' | '));
  console.log('stages:', JSON.stringify(stages), 'pending:', fn.pending_count);
  console.table(ad.events.map(e=>({name:e.name,users:e.users,total:e.total,recent7d:e.recent})));
  window.__pursuit={ad,fn,users};
})();
```

**Pre-committed questions (answer these, don't re-debate):**
1. Did the interview fix work? `addevent_open` → `interview_save:ok`. If saves ≈0
   AND `interview_save:error` ≈0 → nobody tries → feature is **dead, not broken**.
2. Is the dossier still the wedge? `dossier_open`/`dossier_refresh` 7-day trend.
3. Real open→complete drop: `addmodal_open`→`application_create`,
   `addevent_open`→`interview_save` as conversion rates.
4. Did the manual nudges move anyone signed-in → activated? (compare stages vs Jun 22)

**Decision tree (the strategic fork):**
- Dossier thesis holds → Pursuit is **"AI interview prep + a tracker spine"**, not a
  tracker with AI. Reframe activation metric / Today page / onboarding / invite copy
  around it; next build = deepen the dossier (company + JD summary, hiring-manager
  link, logo — see FREE_RUN_NOTES). **Then open the gate** (promote pending beta-interest).
- Activation actually leaks → build the **nudge-email system** → forces the parked
  **mail decision** (deferred since OAuth replaced magic-link).
- Activates but feels thin → UX-polish backlog (rename "dossier"/"loop", proactive
  "what can you do today" section).

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
