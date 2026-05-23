// Hardcoded sample dossier content used until the AI generation endpoint
// lands in v0.2. The shape mirrors what /api/applications/:id/dossier will
// return once Claude is wired in.

export const SAMPLE_DOSSIER = {
  meeting: {
    when:     'Wed 20 May · 14:00 PT',
    duration: '60 min',
    medium:   'Google Meet',
    panel:    'Dario Amodei (CEO)'
  },
  generatedAgo: '12 min ago',
  interviewer: {
    initials: 'DA',
    name:     'Dario Amodei',
    role:     'CEO & Co-founder, Anthropic',
    prior: [
      'VP Research, OpenAI',
      'Senior Research Director, Google Brain',
      'Princeton PhD (computational neuroscience)'
    ],
    links: [
      { label: 'LinkedIn',                  href: '#' },
      { label: 'Essay: Machines of Loving Grace', href: '#' },
      { label: 'Stratechery interview',     href: '#' }
    ]
  },
  snapshot:
    "Dario is a long-form thinker who values written argument; expect open-ended " +
    "<em>why</em> questions, not gotcha trivia. Strong AI-safety conviction. " +
    "Comfortable with discomfort — he asks <em>what would change your mind?</em> and means it.",
  background:
    "Co-founded Anthropic in 2021 after leaving OpenAI with seven colleagues, primarily " +
    "around AI-safety convictions. Background in computational neuroscience (Princeton PhD). " +
    "Has spent the last 18 months as a frequent public communicator — essays, Senate testimony, " +
    "podcast appearances. Writes more like a scientist than a CEO.",
  signals: [
    {
      date: 'Apr 2026',
      kind: 'Essay',
      body: 'Updated "Machines of Loving Grace" — argues for compressed AGI timelines through public-health and economic-growth analogies.',
      source: 'darioamodei.com'
    },
    {
      date: 'Mar 2026',
      kind: 'Senate',
      body: 'Joint testimony on frontier model regulation — emphasised compute reporting and pre-deployment evaluations.',
      source: 'senate.gov'
    },
    {
      date: 'Feb 2026',
      kind: 'Podcast',
      body: 'Stratechery interview — went deep on the trade-off between safety research and competitive dynamics with frontier labs.',
      source: 'stratechery.com'
    }
  ],
  style: {
    lead:
      'Reasoning-first. He cares more about how you think than what you already know. ' +
      'Expect to be asked to defend a position out loud, then asked the opposite — he is testing whether ' +
      'you change your mind in real time when given new information.',
    tells: [
      { lbl: 'Format', val: 'Open-ended discussion, no whiteboard. Comfortable with silence.' },
      { lbl: 'Tempo',  val: 'Slower than most. Take a breath before answering.' },
      { lbl: 'Length', val: '~60 minutes; usually runs over by 10–15.' }
    ]
  },
  lands: [
    'Boring, scalable system-design choices. He prefers obvious correctness over clever tricks.',
    "A clear 2-sentence personal take on AI safety. Doesn't need to match his — needs to be considered.",
    'Real examples from your work where you changed your mind based on data.',
    "Comfort with \"I don't know\" followed by how you'd find out."
  ],
  avoid: [
    'Over-optimising the system-design answer for cleverness. He will read it as immaturity.',
    "Hand-wavy safety takes (\"it's important\") without a concrete position.",
    'Comp-first energy. He has heard it; he is not the right audience.',
    "Citing his own essays back to him verbatim — engage with them, don't repeat them."
  ],
  questions: [
    { q: 'What changed your view on safety research vs. competitive racing in the last 12 months?',
      why: 'Direct, current, references a real tension he writes about. Shows you have actually read him.' },
    { q: 'How are you thinking about the trade-off between RLHF data quality and model scale at the current frontier?',
      why: 'Specific to his stated technical priorities. Shows you can talk shop.' },
    { q: 'If Claude-5 is meaningfully more capable than Claude-4, what is the single thing about your release process that has to change?',
      why: 'Forces a concrete answer about something he is actively working on. Hard to dodge.' },
    { q: "What is one thing you believed about Anthropic's hiring 18 months ago that you no longer believe?",
      why: 'Shows you care about how the team is built. Tests whether he is reflective in public the way he is in essays.' }
  ],
  timeline: [
    { date: '28 Apr', tag: '',         label: 'Noted',           note: 'Heard about the role through Sarah at the AI demo night.' },
    { date: '10 May', tag: 'accent',   label: 'Applied',         note: 'Sent CV v3-ai-focus + cover. Referred by Sarah.' },
    { date: '13 May', tag: '',         label: 'Recruiter reply', note: 'Maya Chen scheduled the screen for May 15.' },
    { date: '15 May', tag: 'accent',   label: 'Phone screen',    note: 'Maya. Strong signal. Moving to technical loop.' },
    { date: '22 May', tag: 'positive', label: 'Round 1',         note: 'System design with Dario. Discussed scaling RLHF data pipelines. Felt good — open follow-ups on cost modeling.' }
  ]
};

export function buildTimelineFromApplication(app) {
  // Stub a minimal timeline from the application data alone — used for
  // applications other than the demo Anthropic one until we track real events.
  const events = [];
  if (app.applied_at) {
    events.push({
      date: shortDate(app.applied_at),
      tag: 'accent',
      label: 'Applied',
      note: `Sent ${app.cv_variant ? `CV ${app.cv_variant}` : 'application'}${app.source ? ` via ${app.source}` : ''}.`
    });
  }
  if (app.status === 'screen' || app.status === 'interview' || app.status === 'offer') {
    events.push({
      date: shortDate(app.updated_at ?? app.applied_at),
      tag: app.status === 'offer' ? 'positive' : 'accent',
      label: capitalize(app.status),
      note: `Moved to ${app.status}.`
    });
  }
  return events;
}

function shortDate(d) {
  if (!d) return '—';
  return new Date(d).toLocaleDateString('en-US', { day: 'numeric', month: 'short' });
}
function capitalize(s) {
  return s ? s[0].toUpperCase() + s.slice(1) : s;
}
