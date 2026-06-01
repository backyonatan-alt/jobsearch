/* ═══════════════════════════════════════════════════════════
   PURSUIT · Today (Option B, activated) — prototype data
   One coherent model behind every screen: Today, Board,
   Application detail, Dossier, Insights.
   ═══════════════════════════════════════════════════════════ */

window.USER = {
  first: "Yonatan",
  last: "B.",
  initials: "YB",
  email: "back.yonatan@gmail.com",
  role: "Product Designer / Design Engineer",
};

window.NOW = {
  dow: "Wednesday",
  short: "Wed",
  full: "20 May 2026",
  stamp: "Wed · 20 May 2026",
};

/* ── Applications ────────────────────────────────────────────
   status: wishlist · applied · screen · interview · offer · closed
   days   = days since last activity (drives "gone quiet")
   ──────────────────────────────────────────────────────────── */
window.APPS = [
  {
    id: "anthropic", co: "Anthropic", short: "A", cls: "lg-A",
    role: "ML Engineer, Pretraining", status: "interview",
    stage: "Final round today", days: 0, hot: true,
    source: "Founder intro", applied: "19 Apr", cv: "v3.2",
    recruiter: { name: "Priya R.", title: "Technical Recruiter" },
    next: { kind: "Interview", label: "Final round", when: "Today · 14:00", in: "in 3h 22m", who: "Dario Amodei", medium: "Google Meet", mins: 60 },
    timeline: [
      { d: "Today", t: "Interview · 14:00 with Dario Amodei", n: "Final round · technical depth", tag: "accent" },
      { d: "12 May", t: "Recruiter call — Priya R.", n: "30 min · pipeline screen passed", tag: "" },
      { d: "10 May", t: "Intro made by Ben K. (advisor)", n: "Warm founder introduction", tag: "" },
      { d: "19 Apr", t: "Application submitted", n: "CV v3.2 · custom cover note", tag: "positive" },
    ],
  },
  {
    id: "stripe", co: "Stripe", short: "S", cls: "lg-S",
    role: "Product Designer, Dashboard", status: "screen",
    stage: "Screen scheduled", days: 1,
    source: "Referral", applied: "13 May", cv: "v3.0",
    recruiter: { name: "Marcus L.", title: "Design Recruiter" },
    next: { kind: "Screen", label: "Recruiter screen", when: "Thu · 11:00", in: "tomorrow", who: "Marcus L.", medium: "Zoom", mins: 30 },
    timeline: [
      { d: "Thu", t: "Screen call · 11:00 with Marcus L.", n: "30 min · intro + role fit", tag: "accent" },
      { d: "15 May", t: "Referral submitted by Dana", n: "Internal referral logged", tag: "" },
      { d: "13 May", t: "Application submitted", n: "CV v3.0", tag: "positive" },
    ],
  },
  {
    id: "plain", co: "Plain", short: "P", cls: "lg-P",
    role: "Senior Designer", status: "screen",
    stage: "Take-home due", days: 1,
    source: "Cold app", applied: "05 May", cv: "v3.0",
    recruiter: { name: "Sofia M.", title: "Founder" },
    next: { kind: "Take-home", label: "Take-home review", when: "Thu · 16:30", in: "tomorrow", who: "Sofia M.", medium: "Google Meet", mins: 45 },
    timeline: [
      { d: "Thu", t: "Take-home review · 16:30", n: "Walk through the submitted concept", tag: "accent" },
      { d: "16 May", t: "Take-home submitted", n: "Onboarding redesign · Figma link", tag: "positive" },
      { d: "08 May", t: "Intro call with Sofia", n: "20 min · liked the portfolio", tag: "" },
      { d: "05 May", t: "Application submitted", n: "CV v3.0", tag: "positive" },
    ],
  },
  {
    id: "vercel", co: "Vercel", short: "V", cls: "lg-V",
    role: "Staff Frontend Engineer", status: "offer",
    stage: "Offer · decide by Fri", days: 2, hot: true,
    source: "Referral", applied: "06 May", cv: "v3.1",
    recruiter: { name: "Ален K.", title: "Eng Recruiter" },
    next: { kind: "Decide", label: "Offer decision", when: "Fri · 17:00", in: "in 2 days", who: "—", medium: "—", mins: 0 },
    timeline: [
      { d: "Fri", t: "Offer decision due · 17:00", n: "Verbal offer extended Mon", tag: "offer" },
      { d: "18 May", t: "Offer received", n: "$210k base · 0.08% · remote", tag: "positive" },
      { d: "14 May", t: "Final panel · 4 rounds", n: "System design + craft", tag: "" },
      { d: "06 May", t: "Application submitted", n: "CV v3.1 · referral", tag: "positive" },
    ],
  },
  {
    id: "figma", co: "Figma", short: "F", cls: "lg-F",
    role: "Staff Product Designer", status: "applied",
    stage: "Application sent", days: 14, quiet: true,
    source: "Job board", applied: "06 May", cv: "v3.2",
    recruiter: { name: "—", title: "No contact yet" },
    next: null,
    timeline: [
      { d: "06 May", t: "Application submitted", n: "CV v3.2 · no reply in 14 days", tag: "" },
    ],
  },
  {
    id: "notion", co: "Notion", short: "N", cls: "lg-N",
    role: "Senior Product Designer", status: "screen",
    stage: "Recruiter chat", days: 9, quiet: true,
    source: "Job board", applied: "09 May", cv: "v3.0",
    recruiter: { name: "Jordan P.", title: "Recruiter" },
    next: null,
    timeline: [
      { d: "11 May", t: "Recruiter chat — Jordan P.", n: "Said they'd follow up — 9 days ago", tag: "" },
      { d: "09 May", t: "Application submitted", n: "CV v3.0", tag: "positive" },
    ],
  },
  {
    id: "ramp", co: "Ramp", short: "R", cls: "lg-R",
    role: "Design Engineer", status: "screen",
    stage: "Take-home submitted", days: 4,
    source: "Cold app", applied: "08 May", cv: "v3.2",
    recruiter: { name: "Casey T.", title: "Recruiter" },
    next: null,
    timeline: [
      { d: "16 May", t: "Take-home submitted", n: "Dashboard prototype · awaiting review", tag: "positive" },
      { d: "10 May", t: "Recruiter screen passed", n: "Casey T. · 25 min", tag: "" },
      { d: "08 May", t: "Application submitted", n: "CV v3.2", tag: "positive" },
    ],
  },
  {
    id: "mistral", co: "Mistral", short: "M", cls: "lg-M",
    role: "Founding Design Engineer", status: "applied",
    stage: "Referral intro", days: 2,
    source: "Referral", applied: "18 May", cv: "v3.2",
    recruiter: { name: "Léa V.", title: "Talent" },
    next: null,
    timeline: [
      { d: "18 May", t: "Referral intro sent", n: "Via Tom — awaiting recruiter reply", tag: "" },
    ],
  },
  {
    id: "linear", co: "Linear", short: "L", cls: "lg-L",
    role: "Design Engineer", status: "applied",
    stage: "Application sent", days: 9,
    source: "Cold app", applied: "11 May", cv: "v3.2",
    recruiter: { name: "—", title: "No contact yet" },
    next: null,
    timeline: [{ d: "11 May", t: "Application submitted", n: "CV v3.2", tag: "" }],
  },
  {
    id: "mercury", co: "Mercury", short: "M", cls: "lg-MC",
    role: "Senior Product Designer", status: "screen",
    stage: "Screen next week", days: 3,
    source: "Job board", applied: "14 May", cv: "v3.0",
    recruiter: { name: "Robin H.", title: "Recruiter" },
    next: { kind: "Screen", label: "Recruiter screen", when: "Tue · 16:00", in: "next week", who: "Robin H.", medium: "Zoom", mins: 30 },
    timeline: [
      { d: "Tue", t: "Screen · 16:00 with Robin H.", n: "Scheduled for next week", tag: "accent" },
      { d: "14 May", t: "Application submitted", n: "CV v3.0", tag: "positive" },
    ],
  },
  {
    id: "cursor", co: "Cursor", short: "C", cls: "lg-CU",
    role: "Product Designer", status: "wishlist",
    stage: "Saved", days: 0,
    source: "—", applied: "—", cv: "—",
    recruiter: { name: "—", title: "Not applied yet" },
    next: null,
    timeline: [{ d: "07 May", t: "Saved to wishlist", n: "Referral path via Sam", tag: "" }],
  },
  {
    id: "replit", co: "Replit", short: "R", cls: "lg-RP",
    role: "Design Engineer", status: "wishlist",
    stage: "Saved", days: 0,
    source: "—", applied: "—", cv: "—",
    recruiter: { name: "—", title: "Not applied yet" },
    next: null,
    timeline: [{ d: "08 May", t: "Saved to wishlist", n: "Cold app planned", tag: "" }],
  },
  {
    id: "calcom", co: "Cal.com", short: "C", cls: "lg-CA",
    role: "Product Designer", status: "closed", outcome: "Rejected",
    stage: "Rejected", days: 11,
    source: "Job board", applied: "30 Apr", cv: "v2.8",
    recruiter: { name: "—", title: "Closed" },
    next: null,
    timeline: [
      { d: "20 May", t: "Rejected after screen", n: "Role put on hold", tag: "danger" },
      { d: "30 Apr", t: "Application submitted", n: "CV v2.8", tag: "" },
    ],
  },
];

window.STATUS_LABEL = {
  wishlist: "Wishlist", applied: "Applied", screen: "Screen",
  interview: "Interview", offer: "Offer", closed: "Closed",
};

/* ── The Anthropic prep dossier — the centrepiece ──────────── */
window.DOSSIER = {
  appId: "anthropic",
  generatedAgo: "12 min ago",
  meeting: { when: "Today · 14:00", duration: "60 min", medium: "Google Meet", panel: "Final · 1 of 1" },
  interviewer: {
    name: "Dario Amodei", initials: "DA",
    role: "Co-founder & CEO, Anthropic",
    prior: ["VP of Research, OpenAI", "PhD Physics, Princeton"],
    links: ["Personal site", "Twitter / X", "Anthropic profile"],
  },
  snapshot:
    "Co-founder and CEO of Anthropic. Trained as a physicist; led early scaling work at OpenAI before leaving in 2021 to start Anthropic with his sister Daniela. Comfortable with first-principles reasoning and long silences. Will go deeper than you expect on whatever you claim to have built — and is allergic to bravado.",
  style: {
    lead: "Conceptual, not behavioural. Expect \"walk me through how you'd think about X\" rather than STAR-style stories. He's comfortable with long pauses while you think — don't fill the silence. He'll probe one or two technical claims deeply rather than survey-style, so come ready to defend whatever you put on your CV.",
    tells: [
      { lbl: "Pace", val: "Slow, deliberate. Often goes quiet for 5–10 seconds after your answer." },
      { lbl: "Energy", val: "Calm, direct. Warm in a low-key way; not performative." },
      { lbl: "Depth", val: "Will pick one or two threads and go all the way down. Won't survey." },
    ],
  },
  lands: [
    "Concrete reasoning about scaling intuitions — \"what surprised me on my last model was…\"",
    "Genuine engagement with mechanistic interpretability if you have real taste there",
    "Comfort saying \"I don't know\" out loud, with a hypothesis attached",
    "Arriving with a sharp question about a recent Anthropic paper",
  ],
  avoid: [
    "Hyping capabilities without a safety lens — reads as missing the point",
    "Vague answers when pushed. Better to admit a gap than to handwave",
    "Credentials-heavy talk. He cares about taste and reasoning, not pedigree",
    "Filling silences. The pause is his thinking, not yours to rescue",
  ],
  signals: [
    { date: "Apr 2026", kind: "Essay", body: "Follow-up to \"Machines of Loving Grace\" on policy implications of capability gains in the next 18 months. Tone shifted toward urgency on governance.", src: "anthropic.com/essays" },
    { date: "Mar 2026", kind: "Interview", body: "~2h Stratechery conversation. Spent most of it on research bets, not competitive positioning. Light on capability claims, heavy on methodology.", src: "stratechery.com" },
    { date: "Feb 2026", kind: "Paper", body: "Co-authored interpretability paper on circuit-level analysis in Claude 4. Suggests he still personally reviews technical work.", src: "transformer-circuits.pub" },
  ],
  questions: [
    { q: "What would surprise you, in either direction, in the next 12 months of pretraining research?", why: "Invites him to think out loud about open problems — territory he enjoys." },
    { q: "How do you personally evaluate research taste outside published work?", why: "Signals you care about the same axis he optimises for in hiring." },
    { q: "Where do you think the field's collective intuitions are most wrong right now?", why: "Open-ended, conceptual — the kind of question he poses to himself." },
  ],
};

/* ── Insights numbers ────────────────────────────────────────*/
window.INSIGHTS = {
  funnel: [
    { stage: "Applied", n: 12 },
    { stage: "Screen", n: 8 },
    { stage: "Interview", n: 4 },
    { stage: "Offer", n: 1 },
  ],
  replyRate: 67,
  avgFirstReply: "4.2 days",
  medianStage: "Screen",
  activity: [3, 5, 2, 6, 4, 7, 5, 9, 6, 8, 4, 7],
  sources: [
    { src: "Referral", n: 4, pct: 33 },
    { src: "Job board", n: 4, pct: 33 },
    { src: "Cold app", n: 3, pct: 25 },
    { src: "Founder intro", n: 1, pct: 9 },
  ],
  bestSource: "Referral",
};
