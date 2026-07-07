# Pursuit — Strategy Investigation (Jul 7 2026)

> **What this is.** A one-time deep investigation into where Pursuit can go: scenarios,
> roadmap, monetization, distribution, messaging. Commissioned Jul 7 2026 with the frame:
> *"Can this become a product/business? Monetization in scope, but not super big. No
> funding. No timeline."* Built from six parallel research lanes (competitors, market,
> monetization, distribution, internal evidence, feature/scenario space) + a fresh
> aggregate prod-data snapshot + an adversarial red-team pass over the conclusions.
>
> **How to use it.** This is the reference doc, not a todo list. Re-read at decision
> points (the §12 kill-criteria table says when); log verdicts in §12; TODO.md stays the
> operational surface. Claims are marked **[observed]** (real data / real users) or
> **[hypothesis]** (reasoned but unproven) — don't let them blur.

---

## 1. The answer up front

**Yes, it can plausibly become a small, sustainable business — but today only one strategic
belief is actually earned, and everything else must be staged as experiments with
pre-committed kill criteria.**

The earned belief **[observed]**: *the playbook (grounded, cited, per-round interview prep)
is the wedge.* Three independent evidence sources agree: top engagement surface across four
consecutive data re-reads; two external testers praising it unprompted ("maybe the most
important value this tool gives"); 40% of a cold LinkedIn wave reaching it on day one.

Everything stacked on top — pricing, revenue, virality, the community-data moat, coaches —
is currently asserted, not proven. So the strategy is:

1. **Run the "Walk in cited" scenario (§6.2) as a priced experiment**, not a business plan:
   ship prep packs via Paddle, with success/kill thresholds written down *before* launch.
2. **Build trust infrastructure before wider distribution** — an automated grounding/citation
   eval harness and a published interviewer-privacy policy. "Trustworthy prep" backed by one
   manual QA pass is a slogan, not a mechanism.
3. **Keep the personal-power-tool scenario (§6.1) as the explicit floor.** Either verdict of
   the pricing experiment is a fine outcome. That's what "no funding, no timeline, not super
   big" buys: the freedom to let the data decide without existential stakes.

---

## 2. Where we are — the evidence base

### Funnel history [observed]

| Date | Invited | Signed in | Activated | Active ≤7d | Notes |
|---|---|---|---|---|---|
| Jun 17 | 36 | 22 | 13 | 7 | interview flow found broken (0/22), fixed |
| Jun 22 | 36 | 25 | 16 | ~6 | activation 64%; dossier confirmed as wedge |
| Jun 29 | 48 | 26 | 15 | 6 | +12 invites → ~0 sign-ins; "Playbook" rename |
| Jul 6 | 50 | 32 | 21 | 10 | reframe email moved +6/+6; open signup + LinkedIn launch |
| Jul 7 | open | 44 non-admin | — | 21 | +10 day-one signups, 10 with `src:li`; 4/10 to playbook day one |

### Depth of use (Jul 7 snapshot, non-admin) [observed]

- **19/44 (43%) ever generated a playbook**; distribution: 15 users×1, 3×3, 1×5 → the
  *average prep-generating user makes ~1–2.5 playbooks*. This distribution sets the paywall (§7).
- `dossier_open`: 86 opens / 21 users in 30d — top surface by far. `paste_parse` 31/7.
- **Only 8/33 signed-in users ever returned on a later calendar day.** But: only **5 users
  have interview rounds at all** — most users have had nothing to come back *for*. The honest
  retention metric is **return-per-interview-occasion**, currently uncomputable at n=5.
  (Red-team correction: calendar-day return is a social-app metric misapplied to an
  occasion-driven product. Don't panic on it; instrument the right one.)
- Real-user debriefs: **0** — but the debrief prompt was structurally unreachable until Jul 6,
  and the eligible denominator is ~5 users. Zero is currently *noise*, not a verdict (§5).
- Spend: 38 prep credits / 15 users; COGS ~$0.50–1.20 per full first-round prep.

### What users actually said [observed]

- *"This is excellent. Maybe the most important value this tool gives, in my eyes."* (Ayelet,
  first external tester, unprompted, about the playbook)
- The two things that nearly killed it were both **trust failures**: wrong same-named company
  researched (the 365scores collision — an *Israeli-company* collision, note for §9), and
  citations linking to homepages instead of sources ("the interviewer may ask where you know
  this from").
- Tracker-spine features are structurally underused (board/funnel/import barely touched;
  `interview_save` was a "vanity hole"). Every attempt to invest in tracker breadth has been
  a worse bet than deepening prep.

---

## 3. Market: the bet and the weather

Full lane report available in session history; the load-bearing points:

- **The application stage is a lottery; the interview is the scarce asset.** Application
  volume up >45% YoY (LinkedIn); ~2–3% of applications convert to interview (soft sources,
  directionally consistent); AI screening + ghost jobs make top-of-funnel near-random. Each
  landed human interview is therefore the highest-stakes moment in the funnel — Pursuit's
  core bet, directionally validated. **[observed-market]**
- **Prep is the least-saturated AI behavior**: only ~21% of seekers use AI to research the
  company (Resume Genius 2026) vs ~40–46% for resumes. *Red-team caveat: unsaturated ≠ open —
  79%-don't could also mean low demand. Treat as opportunity signal, not proof.* **[hypothesis]**
- **The "cheating copilot" backlash is a positioning gift.** Final Round AI / Cluely /
  live-answer tools are being actively hunted by employers (detection startups, in-person
  round reinstatement at Google/McKinsey). *Pre*-interview prep is the explicitly legitimate
  safe harbor. "We make you prepared, not fed" is the clean side of a line the market is
  drawing right now. **[observed-market]**
- **Platform posture is a quiet moat.** LinkedIn's 2026 scraping crackdown (suspensions,
  litigation, Proxycurl killed) makes Pursuit's screenshot-paste, user-as-conduit design
  durable while scraped-data competitors carry existential risk. **Sharpest named threat:
  LinkedIn's own Interview Prep AI** — they hold the JD, the interviewer profile, and the
  company page natively. Their consumer-product velocity is historically poor and their
  revenue is employer-side; this caps the ceiling rather than evicting us, but it's the
  competitor to watch. **[hypothesis]**
- **All three plausible 2027 futures** (agents-apply/humans-interview ~50%, AI-interviews-
  everything ~25%, trust-collapse/verification-regime ~25%) reward trust depth (grounding,
  citations, pipeline memory) and none reward tracker breadth. If AI screens eat early rounds,
  value shifts toward the company brief / "what this loop grades for" — already half the product.
- **Israel now**: the biggest layoff wave since 2022 (Wix −20%, Amdocs, Intuit IL), searches
  stretched to 45–60-day cycles. Demand for prep is at a local maximum; willingness-to-pay is
  under pressure. See §9 for why Israel is also our hardest grounding terrain.

## 4. Competitive map

| Segment | Who | Price | What they monetize | Verdict for us |
|---|---|---|---|---|
| Tracker-first | Teal (~$29/mo), Huntr ($40), Simplify ($40), Careerflow ($24–45) | $24–40/mo | AI resume/cover-letter + autofill (the *application* stage) | Don't compete. Tracker stays free scaffolding. |
| Live copilots | Final Round (~$148/mo), Cluely, LockedIn | $20–148/mo | In-interview answer feeding | Radioactive category, imploding on detection + billing scandals. Our anti-positioning. |
| Skill practice | Exponent (~$150/yr), interviewing.io ($225/mock), Prepfully | $120–300/session | Human coaching, FAANG-track | Different job. Their prices are our headroom anchor. |
| Agentic appliers | Simplify agent, Jobright, LazyApply, LoopCV | $30–100 | Application volume | Dying thesis; they *create* interviews they don't help win — a funnel into us, not competition. |
| **The real one** | ChatGPT / Claude / Gemini | $0–20/mo | — | Beats us on price, scope, model quality. Loses on grounding, citations, pipeline memory, feed-forward, and the 30–60 min prompting tax. We are "ChatGPT-grade research with the prompting, grounding, verification, and pipeline-memory done for you." |

**White space [hypothesis]:** nobody ships correct-company, cited, per-interviewer, per-round
briefs attached to a pipeline. Google retired Interview Warmup and points at Gemini — generic
prep gets absorbed into chatbots; only the grounded/stateful part is defensible.

**Explicitly not competing on:** resume/cover-letter generation, auto-apply volume, live
copilots, mock-practice-as-category (Gemini Live does it free; Yoodli fled to B2B), job
discovery, human-coaching marketplace.

## 5. The moat, honestly

Ranked by how real they are today:

1. **Trust infrastructure** (grounding, per-claim citations, identity block, re-ground
   control) — shipped, gate-passed once, *not yet systematized* (§9 makes it a workstream).
   This is the brand and the hardest thing for a casual competitor to bolt on credibly.
2. **Pipeline memory + within-application feed-forward** (round N's debrief informs round
   N+1's prep) — structurally something a stateless chat can't do without the user hand-
   managing context. Shipped; usage unproven (denominator ~5). The Jul 13 re-read is its
   first real data point. **[hypothesis with a date]**
3. **Compliance-clean data posture** — user-as-conduit (screenshots), no scraping. Gets
   *more* valuable as LinkedIn enforcement tightens. **[observed-market]**
4. ~~**Cross-user "what {company} actually asks" corpus**~~ — **cut from the plan** (red-team
   verdict accepted): needs per-company debrief *density* a bootstrap product never reaches
   (1k users spread over hundreds of companies = n=1–2 anecdotes per company), and the mature
   version already exists at Glassdoor and is mediocre/gamed — evidence the data structure
   itself has low defensive value. Within-app feed-forward stays; the network-effect chapter
   dies unless some future re-read shows organic debrief density we didn't earn on purpose.

**Commoditization ceiling, stated plainly:** ChatGPT with memory + browsing produces ~80% of
a playbook for a skilled prompter inside a subscription they already pay. Our defense is
workflow packaging at the panic moment — which genuinely sells (interview courses always sold
despite free libraries) — but it caps this at *niche product*, which matches the stated
ambition. It does not survive contact with "become big."

## 6. Scenarios

### 6.1 "Sharpest tool in my shed" — personal power tool (the floor)
Owner + friends; excellence over growth; costs <$50/mo; payoff is the owner's own search and
reputation. **Kills it:** owner lands a job, motivation decays. **This scenario is always
available and is a legitimate outcome.** Every investment below is chosen to also make this
version better, so "failure" of the business question wastes almost nothing.

### 6.2 "Walk in cited" — prep-first indie micro-SaaS ⭐ the priced experiment
Sell grounded, citable, per-round prep to individual seekers; free tracker scaffolding;
playbook = the metered unit; **credit packs, not subscriptions** — churn-on-success becomes
the model instead of its enemy (they leave happy, refer mid-search, return next search).
Positioning: *for candidates who refuse to walk into an interview on vibes — every claim has
a source you can name out loud.* **Kills it:** frontier chat ships grounded per-round prep
natively; or willingness-to-pay is ~$0. **Verdict mechanism: §7's pre-committed thresholds.**

### 6.3 "Interview-intelligence layer" — community corpus moat: **cut** (see §5.4)
Restaged: keep within-app feed-forward and user-visible prep-accuracy scoring. Revisit only
if debrief density materializes organically. No investment beyond that.

### 6.4 "Coach OS" — B2B2C: **declined now, door open**
Coaches are the persistent node in a transient market and would *force* debrief collection —
elegant on paper. But the binding constraint is founder-hours in sales/support, the one
resource Claude doesn't multiply, and it contradicts "not super big."
**Door open:** "buying for a client? email me" link on the pricing page + a few free coach
accounts as *distribution* (§8). Build seats only if ≥3 coaches ask unprompted.

### 6.5 "The prep library" — SEO/content engine: **deferred to a small citable library**
AI Overviews cannibalized informational SEO (CTR −15–47%, Chegg's collapse); a pSEO farm
also contradicts the trust brand. Later: 30–50 pages from real (consented, anonymized) brief
runs for companies users actually interview at, freshness-dated, with a "generate yours"
CTA. Measure 90 days before scaling. Marketing, not product.

## 7. Monetization: a priced experiment, not a forecast

**Red-team correction adopted: no revenue projections.** The observed usage distribution
(15 users×1 brief, 3×3, 1×5) means a free tier of 3 binds on roughly the top 20% of
prep-generating users — which is exactly where a paywall should sit: above the median,
binding only on demonstrated-value users.

**The setup:**
- **Free: 3 preps** (down from 10). One prep = one full round; company brief bundled free
  with the first prep per application (as today). Refresh of an existing brief = free —
  reinforces the trust loop, costs little.
- **Packs: $9 → 4 preps, $19 → 10 preps.** Never expire. Sold as "preps," never abstract
  credits (Cursor's credit-pool backlash is the cautionary tale).
- **"Interview Pass" $24 / 30 days** — unlimited-with-fair-use (~20). The subscription-shaped
  option with no cancellation machinery (one-time, auto-expires). A/B against packs after
  ~50 purchases.
- **Rails: Paddle** (merchant of record; supports Israeli sellers; handles global VAT/tax —
  worth the ~2% premium over raw Stripe, which doesn't directly support Israel anyway).
  One accountant conversation re מע"מ on payouts.
- **Abuse guards:** Google-OAuth-only (kills disposable emails) + per-account/IP rate limits
  + per-day cap. No fingerprinting at this scale.
- **Sequencing:** trust harness (§9) ships *before* the paywall — charging for a
  wrong-company brief converts a bug into a refund plus a reputation hit.

**Pre-committed thresholds (the scenario-6.2 verdict — write the numbers here before launch):**
- **Success signal:** ≥5 purchases from the first ~100 post-paywall users who generate a
  playbook, without discounting. (≈5% conversion of prep-activated users — the industry
  freemium prior applied to our funnel.)
- **Soft signal:** 1–4 purchases → price/packaging iteration (try the Pass, try $29), one
  more cohort, then re-verdict.
- **Kill signal:** 0 purchases from 100 prep-activated users → scenario 6.2 is falsified;
  fall back to 6.1 with heads held high. Also track **prep-accuracy** ("spot-on" ≥70% on
  ≥30 rated playbooks) as the quality gate — payment without accuracy is borrowed time.
- Unit sanity: ~80% gross margin on paid usage at these prices; the P&L is dominated by
  free-tier generosity, which is why free went from 10 → 3.

**Why charge early:** a purchase is the strongest trust signal we can instrument — it *is*
the "playbook good enough to walk into a real interview" activation metric the north star
asks for, expressed in currency.

## 8. Distribution

Red-team correction adopted: **job searches are secret.** Nobody shares their prep for
Company X while interviewing at Company X. The shareable moment is **the offer**. The
prep-share OG-card virality bet is demoted; private channels and the post-offer artifact are
promoted.

Ranked plan:

1. **LinkedIn (yours), systematized** — the only proven channel (1 post → 10 signups → 40%
   day-one activation). 2 posts/week alternating build-in-public and prep-artifact content;
   document/carousel format (top reach in 2026). Accept honestly: this is founder labor, it
   saturates, and it stops when the founder's attention does (§11).
2. **The post-offer artifact** *(replaces prep-share virality)* — at status→offer, generate
   "what I walked in knowing": a shareable one-pager + testimonial prompt. The only moment a
   job seeker *wants* to broadcast, and the only moment name-dropping the tool costs nothing.
3. **LayoffRadar → Pursuit handoff** — the ICP at peak need, via an owned channel. Tone rules
   (non-negotiable): standalone occasional message or pinned footer, never attached to a
   specific company's layoff post; agency framing ("when you land interviews, walk in more
   prepared than the panel"); context-aware landing page; `src=layoffradar` attribution
   (rails already proven with `src=li`).
4. **3–5 Israeli coaches / bootcamp career services, friend-scale** — free coach accounts +
   client invites, zero code (admin invites exist). This is Huntr's $1M-ARR channel at
   pilot scale, and it doubles as the §6.4 demand probe.
5. **Reddit/communities, value-first** — answer "interview at X next week" threads with the
   actual methodology; tool in profile only. Slow, compounding, free ICP research.

**Share links still ship** (read-only company brief, OG card, `?ref=`) — but as a *utility*
for private sends (mentor, friend, coach), not as the growth engine. Every shared brief is a
public demo, so §9 gates this too.

## 9. Trust infrastructure — a named workstream (gates everything)

The doc's biggest claim-vs-mechanism gap, per the red team: the north star is *trustworthy*
prep, but the grounding gate passed on **one** manual QA (Lusha) and there is no automated
evaluation. Two deliverables gate the next distribution wave:

1. **Grounding/citation eval harness.** Golden fixtures of production-shaped cases —
   *especially Israeli-company name collisions* (365scores was one; our beachhead concentrates
   our worst grounding terrain: thin English web presence, Hebrew sources, transliterated
   names). Assert: right company identified, citations resolve, cited passage substantiates
   the claim. Run on every deploy and on any model change (single-provider drift is otherwise
   silent). This is CLAUDE.md's own golden-output bar, finally applied to the core artifact.
2. **Interviewer-privacy policy, published.** The signature feature is AI-compiled research
   on a named person who never consented — one hostile screenshot ("creepy AI stalking tool")
   away from a reputation event that would travel through exactly the tight network we
   distribute in. Policy before the wave: public professional sources only; no personal-life
   inference; visible sourcing (already built for trust — doubles as defensibility); a stated
   line we don't cross; GDPR data-subject posture for EU-company interviews.

Supporting items: wrong-company kill-switch metric (alert on "Not them?" re-ground clicks),
per-claim confidence labels, freshness re-check, user-visible prep-accuracy score, 👍/👎
feedback-on-claims (micro-labels are cheaper debriefs and may succeed where the 20-sec
debrief is failing).

## 10. Feature almanac (effort for solo+Claude: S/M/L)

**Now-tier (cheap, high-leverage, all serve the wedge + trust):**
| Feature | Effort | Why |
|---|---|---|
| Mobile read-only playbook page | S/M | **Promoted by red team:** prep is consumed night-before, on a phone; we can't charge for an artifact unreadable at its moment of use. A read-only responsive brief view is a slice, not the full mobile pass. |
| Day-of cheat sheet (printable/phone one-pager) | S | The north star as an artifact — "walk in with and cite," literally. |
| Feedback-on-claims (👍/👎 per claim) | S | Trust data + cheaper debrief substitute. |
| Freshness indicator + pre-interview re-check | S | Drives the refresh trust signal. |
| Follow-up email draft from debrief | S | Rewards debriefing — attacks debrief=0 sideways. |
| Post-offer artifact (§8.2) | S/M | The viral moment that respects search secrecy. |
| Eval harness + privacy policy (§9) | M | Gates everything. |

**Next-tier (after the pricing verdict):**
mock-interview mode (text, seeded by the actual brief + debriefs — M, the natural paid-depth
feature), calendar-triggered prep (Google Calendar read scope — M, the *trigger* for the prep
moment), reverse-diligence/red-flag brief (M — LayoffRadar data bridge, same citation bar),
salary/negotiation brief at offer stage (M — least price-sensitive moment; cite or don't ship),
AI weekly review (M — gate on pipeline size; horoscope risk at n=6).

**Declined:** Gmail auto-status (CASA review + "it reads my email" chills trust-sensitive
users), voice mock mode (infra; Gemini Live does it free), browser extension (maintenance tax
now; revisit if capture friction shows up in data), white-label, pSEO farm, classic referral
program, donations.

## 11. Messaging & positioning

**Positioning statement:** *Pursuit is interview prep you can walk in with and cite. It
researches the right company and your actual interviewer, round by round — every claim
sourced — and remembers your pipeline so round 2's prep knows what round 1 asked. ChatGPT
guesses; Pursuit checks.*

**Message bank** (tested shapes, per audience):
- Seekers, sharp: "When the interviewer asks *'where did you read that?'* — have an answer."
- Seekers, differentiation: "ChatGPT prepped you for a company with the same name. Pursuit
  checked it's the right one."
- Seekers, feed-forward: "Round 2 prep that knows what Round 1 asked."
- Anti-copilot line (as backlash grows): "We make you prepared, not fed. Nothing to detect —
  because the work happened before you walked in."
- Laid-off audience (LayoffRadar, agency-framing): "When you land interviews, walk in more
  prepared than the panel."
- Coaches: "Stop spending Sunday nights googling your clients' interviewers."
- Build-in-public: "I built the interview-prep tool I wished existed. Here's what it got
  right in my last loop — with sources."

**Tone rules:** never attach promotion to a specific company's layoff; never imply the
interviewer was "profiled" — the vocabulary is *prepare, research, sources*, never *dossier,
track, intel* in public copy (internal name "dossier" already renamed "Playbook" for this
reason — the instinct was right, extend it to all public surfaces).

## 12. Decision log & kill criteria

| Bet | Verdict mechanism | Date/trigger | Outcome (fill in) |
|---|---|---|---|
| Scenario 6.2 pricing | ≥5 / 1–4 / 0 purchases per 100 prep-activated post-paywall users (§7) | ~6 weeks after Paddle ships | — |
| Debrief loop demand | `debrief_save` > 0 among eligible users with the Jul 6 surfacing fix live; then *density* (≥N per active user), not just >0 | Jul 13 re-read + one more cycle | — |
| Prep quality (trust) | prep-accuracy "spot-on" ≥70% on ≥30 rated playbooks | as ratings accumulate | — |
| Grounding holds live | eval harness green on deploys; Ayelet real-loop retry; "Not them?" click-rate ~0 | before next distribution wave | — |
| Coach demand (6.4 door) | ≥3 coaches actively run real clients; ≥1 pays unprompted | 4 weeks after coach pilot starts | — |
| SEO library | 20 pilot pages → ≥1 signup/week combined within 8 weeks of indexing | 90 days after publish | — |
| LinkedIn-native prep threat | LinkedIn ships per-interviewer candidate prep | monitor | — |

**Founder-exit paragraph (required honesty):** the founder is the content engine, the QA
harness, the coach channel, and — because he is job-searching — the lead dogfooder. A
successful search removes all four within weeks; the LinkedIn cadence (the only proven
channel) breaks first. Mitigation: the §9 eval harness replaces founder QA attention; Paddle
billing + credit packs idle gracefully; the product at 4 hours/week is scenario 6.1 with
revenue — which is an acceptable resting state, and the reason 6.1 is the designed floor
rather than a failure mode.

**Standing risks:** single-provider dependency (Anthropic pricing/web_search — credit model
reprices easily; eval harness catches drift); repo is public (no secrets in code; OAuth
secret rotation still pending — do it); brand "Pursuit" is a placeholder config string and
there is still no real domain — both become real work the moment money is charged (a payment
page on a raw-IP nip.io URL is a trust contradiction; **domain purchase is a prerequisite of
the §7 experiment**, the one new spend this plan requires, ~$10).

## 13. The sequenced roadmap (gated, no dates)

**Phase A — trust rails + the experiment's prerequisites** *(before charging or promoting)*
1. Grounding/citation eval harness on production-shaped fixtures incl. Israeli collisions (§9.1)
2. Interviewer-privacy policy published (§9.2) + OAuth secret rotation
3. Mobile read-only playbook view (the paid unit, consumable where it's consumed)
4. Real domain + brand check (config-string swap decision: keep "Pursuit" or rename once, now)
5. Day-of cheat sheet + feedback-on-claims (cheap trust surface)

**Phase B — the priced experiment**
6. Free tier 10 → 3; Paddle; $9/$19 packs (+$24 Pass behind a flag); thresholds of §7 logged
7. Post-offer artifact + share links (utility form)
8. Distribution wave 2: LinkedIn cadence + LayoffRadar handoff + 5 coach pilots (§8)

**Phase C — verdict-dependent**
- Pricing **success** → mock-interview (text), calendar-triggered prep, salary brief at offer,
  weekly review; scale channels; revisit coach door if it knocked (§6.4)
- Pricing **kill** → scenario 6.1: keep the tool excellent for the owner + friends, costs
  <$50/mo, write the postmortem in this file, feel zero regret
- Debrief data **appears** → user-visible prep-accuracy + feed-forward polish (never the
  cross-user corpus without density evidence)
- Jul 13 wave re-read → answers the five pre-committed questions; updates §12

---

*Compiled Jul 7 2026 by the six-lane research pass + red team. Sources for external claims
live in the lane reports (session transcript); key ones inline. The single most important
sentence for the next re-read: __prep-as-wedge is earned; everything else is an experiment
with a number attached.__*
