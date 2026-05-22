// Sample data shared by the three design previews.

export const PREVIEW_APPS = [
  {
    id: 1, slug: 'anthropic',
    company: 'Anthropic',
    role: 'Senior Software Engineer',
    status: 'interview',
    applied_at: '2026-05-10',
    cv_variant: 'v3-ai-focus',
    source: 'Referral',
    location: 'San Francisco'
  },
  {
    id: 2, slug: 'stripe',
    company: 'Stripe',
    role: 'Staff Backend Engineer',
    status: 'applied',
    applied_at: '2026-05-12',
    cv_variant: 'v2-payments',
    source: 'LinkedIn',
    location: 'Remote'
  },
  {
    id: 3, slug: 'vercel',
    company: 'Vercel',
    role: 'Founding Engineer',
    status: 'screen',
    applied_at: '2026-05-05',
    cv_variant: 'v3-ai-focus',
    source: 'Twitter',
    location: 'New York'
  },
  {
    id: 4, slug: 'linear',
    company: 'Linear',
    role: 'Senior Product Engineer',
    status: 'offer',
    applied_at: '2026-04-22',
    cv_variant: 'v2-payments',
    source: 'Cold email',
    location: 'San Francisco'
  },
  {
    id: 5, slug: 'figma',
    company: 'Figma',
    role: 'Tech Lead, Platform',
    status: 'rejected',
    applied_at: '2026-05-01',
    cv_variant: 'v3-ai-focus',
    source: 'LinkedIn',
    location: 'Remote'
  },
  {
    id: 6, slug: 'notion',
    company: 'Notion',
    role: 'Senior Engineer',
    status: 'wishlist',
    applied_at: null,
    cv_variant: null,
    source: '—',
    location: 'New York'
  }
];

export const STATUSES = ['wishlist', 'applied', 'screen', 'interview', 'offer', 'rejected', 'withdrawn'];

export const FUNNEL = {
  wishlist: 1,
  applied: 1,
  screen: 1,
  interview: 1,
  offer: 1,
  rejected: 1,
  withdrawn: 0
};

// Detailed view for /preview/{a,b,c}/anthropic
export const ANTHROPIC_DETAIL = {
  ...PREVIEW_APPS[0],
  jd_url: 'https://www.anthropic.com/careers',
  next_step: 'Technical interview · May 24, 10:00 AM PT · with Dario Amodei',
  timeline: [
    { date: '2026-04-28', kind: 'noted',     text: 'Heard about the role through Sarah at the AI demo night.' },
    { date: '2026-05-10', kind: 'applied',   text: 'Sent CV v3-ai-focus + cover. Referred by Sarah.' },
    { date: '2026-05-13', kind: 'reply',     text: 'Recruiter screen scheduled with Maya Chen for May 15.' },
    { date: '2026-05-15', kind: 'screen',    text: 'Phone screen with Maya. Strong signal. Moving to technical loop.' },
    { date: '2026-05-22', kind: 'interview', text: 'Round 1: system design with Dario. Discussed scaling RLHF data pipelines. Felt good — open follow-ups on cost modeling.' }
  ],
  notes: [
    "Dario mentioned they're rethinking the data layer for Claude 5 training. Worth referencing the Stripe payments-pipeline work in the next round.",
    "Maya asked about why I want to leave my current role — keep the answer about mission and scale of impact, not comp.",
    "Glassdoor average for this band: $310k base + equity. My ask: $335k base. Anchor first."
  ],
  dossier: {
    name: 'Dario Amodei',
    title: 'CEO & Co-founder, Anthropic',
    summary:
      "Former VP of Research at OpenAI; left in 2021 with seven colleagues to start Anthropic around AI safety. Has a background in computational neuroscience (Princeton PhD). Public communicator — comfortable with long-form interviews and essays.",
    recent: [
      { date: '2026-04', text: '"Machines of Loving Grace" essay update — argues for compressed timeline to powerful AI, frames it through public-health and economic-growth analogies.' },
      { date: '2026-03', text: 'Joint Senate testimony on frontier model regulation — emphasised compute reporting and pre-deployment evals.' },
      { date: '2026-02', text: 'Stratechery interview — went deep on the trade-off between safety research and racing dynamics with frontier labs.' }
    ],
    style: [
      'Long-form thinker — prepare for open-ended "why" questions, not gotcha trivia.',
      'Values written communication — expect a take-home or written exercise downstream.',
      'Likely to probe on AI safety convictions; have a 2-sentence personal take ready.'
    ],
    watchfor: [
      'Don\'t over-optimise the system-design answer for clever — he prefers boring, scalable choices.',
      'He will ask "what would change your mind on X?" — have a real answer.'
    ]
  }
};

// Per-design funnel breakdown.
export const FUNNEL_VIEW = {
  stages: [
    { key: 'applied',   label: 'Applied',   count: 5, pct: 100 },
    { key: 'screen',    label: 'Screen',    count: 3, pct: 60 },
    { key: 'interview', label: 'Interview', count: 2, pct: 40 },
    { key: 'offer',     label: 'Offer',     count: 1, pct: 20 }
  ],
  insights: [
    {
      title: 'Your application → screen rate (60%) is well above average.',
      body: 'Industry benchmark for senior engineering roles is ~25%. Whatever CV variant you\'re sending is working — keep v3-ai-focus as your default.'
    },
    {
      title: "You're stalling at the screen → interview step.",
      body: 'Only 2 of 3 screens have advanced. The two stalled ones (Stripe, Vercel) are both > 7 days old. Consider a gentle nudge to the recruiters this week.'
    },
    {
      title: 'One live offer (Linear). Reply window closes Friday.',
      body: 'Set the negotiation anchor at $335k base; you have leverage from Anthropic\'s active interview loop.'
    }
  ]
};

export function fmtDate(d) {
  if (!d) return '—';
  return new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
}

export function fmtLongDate(d) {
  if (!d) return '—';
  return new Date(d).toLocaleDateString('en-US', { weekday: 'short', month: 'long', day: 'numeric' });
}
