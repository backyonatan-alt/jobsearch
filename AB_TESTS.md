# A/B Tests

Active and completed experiments. Each entry: hypothesis, variants, success metric, status, result.

## Active

### prep-first cold start (started Jun 29 2026)

**Hypothesis.** New users who hit an instant interview Playbook *before* being
asked to log applications activate at a higher rate — and form a clearer "this
is what Pursuit is for" model — than users who land in the tracker-first guided
tour. (Driven by the Jun 29 re-read: the Playbook is the wedge; the funnel is
stalled at the top; we're testing how prep-first the product should be without
ripping out the tracker spine.)

**Variants.**
- **control (`tour`):** post-OAuth → GuidedTour (welcome → 6 highlights → "add
  your first application" modal).
- **treatment (`prepfirst`):** post-OAuth → full-screen "Who are you interviewing
  with?" prompt (Company autofocused; optional Role) → CTA "Build my playbook" →
  generate-loading state → land in `/app/[id]` with the **company Playbook**
  rendered + a "first playbook" welcome banner. Secondary "I'm just exploring —
  skip →". (Interviewer is collected later on the detail page for round-by-round
  prep — we don't collect input the cold-start brief would ignore.)

**Mechanics (reuses existing endpoints — no new model).** On submit: create
application `{company, role?, status:'screen'}` → `POST /dossier/refresh` (company
scope) → redirect to `/app/{id}` → set `onboarded_at`. The prep question creates
the first tracked application as a byproduct, so the spine stays intact.

**Instrumentation.** `onboard_variant_assigned{variant}`, `prepfirst_prompt_view`,
`prepfirst_submit{has_role,has_interviewer}`, `prepfirst_generate_ok|error`,
`prepfirst_skip` + existing `application_create`/`dossier_open`.

**Success metric.** Primary: signup → first playbook generated (activation),
treatment vs baseline. Guardrail: do prep-first users still add a 2nd+ application
(the spine is the retention engine)? Secondary: D1/D7 return.

**Assignment (pragmatic, given N).** At ~6 active users a 50/50 split is badly
underpowered (would take months to hit the ~30-through-step rule). So treatment
ships to **100% of new signups**, compared against the historical tracker-first
cohort (25 signed-in / 16 activated baseline). Quasi-experimental on purpose; add
a holdback only if traffic picks up. The `users.onboarding_variant` flag exists so
we *can* split later.

**Status.** Building behind the `onboarding_variant` flag (Jun 29).

## Backlog of ideas (not yet running)

- **Onboarding empty state:** "Paste your first LinkedIn job URL" vs "Add an application manually" — which gets the user past zero faster?
- **Dossier teaser placement:** show the dossier CTA on the application card vs only when an interview event exists.
- **Magic-link copy:** "Sign in to Pursuit" vs "Open your job search dashboard".
- **Funnel view vs list view as the default landing screen** once a user has ≥10 applications.

## Completed

(none yet)

## Methodology notes

- Carry the Aegis discipline: every measurable surface gets a tracked event before the test starts, not bolted on after.
- Sample size: don't call a result before the smaller variant has at least ~30 users through the funnel step being measured. Closed beta means real tests will move slowly — that's fine.
