// Sample data shared by the three design previews. Six applications spanning
// every status so the visual treatment for each status reads clearly.

export const PREVIEW_APPS = [
  {
    id: 1,
    company: 'Anthropic',
    role: 'Senior Software Engineer',
    status: 'interview',
    applied_at: '2026-05-10',
    cv_variant: 'v3-ai-focus',
    source: 'Referral',
    location: 'San Francisco'
  },
  {
    id: 2,
    company: 'Stripe',
    role: 'Staff Backend Engineer',
    status: 'applied',
    applied_at: '2026-05-12',
    cv_variant: 'v2-payments',
    source: 'LinkedIn',
    location: 'Remote'
  },
  {
    id: 3,
    company: 'Vercel',
    role: 'Founding Engineer',
    status: 'screen',
    applied_at: '2026-05-05',
    cv_variant: 'v3-ai-focus',
    source: 'Twitter',
    location: 'New York'
  },
  {
    id: 4,
    company: 'Linear',
    role: 'Senior Product Engineer',
    status: 'offer',
    applied_at: '2026-04-22',
    cv_variant: 'v2-payments',
    source: 'Cold email',
    location: 'San Francisco'
  },
  {
    id: 5,
    company: 'Figma',
    role: 'Tech Lead, Platform',
    status: 'rejected',
    applied_at: '2026-05-01',
    cv_variant: 'v3-ai-focus',
    source: 'LinkedIn',
    location: 'Remote'
  },
  {
    id: 6,
    company: 'Notion',
    role: 'Senior Engineer',
    status: 'wishlist',
    applied_at: null,
    cv_variant: null,
    source: '—',
    location: 'New York'
  }
];

export const FUNNEL = {
  wishlist: 1,
  applied: 1,
  screen: 1,
  interview: 1,
  offer: 1,
  rejected: 1,
  withdrawn: 0
};

export function fmtDate(d) {
  if (!d) return '—';
  return new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
}
