// Mapping helpers that translate backend Application objects into the shape
// expected by the Crystalline Console design system (logo class, short label,
// formatted date, etc.). Keeps presentation logic out of the components.

export const STATUSES = [
  'wishlist', 'applied', 'screen', 'interview', 'offer', 'rejected', 'withdrawn'
];

export const STATUS_LABEL = {
  wishlist:  'Wishlist',
  applied:   'Applied',
  screen:    'Screen',
  interview: 'Interview',
  offer:     'Offer',
  rejected:  'Rejected',
  withdrawn: 'Withdrawn'
};

// Logo CSS class lookup. The design system defines .lg-A, .lg-S, .lg-V, etc.
// We fall back to .lg-A when the first letter isn't represented.
const LOGO_CLASSES = new Set(['A', 'M', 'S', 'V', 'L', 'F', 'N', 'R', 'C', 'P']);

export function toDisplayApp(a) {
  const first = (a.company || '?').trim().charAt(0).toUpperCase();
  const logoLetter = LOGO_CLASSES.has(first) ? first : 'A';
  return {
    id:       a.id,
    co:       a.company,
    coShort:  first,
    logoCls:  `lg-${logoLetter}`,
    role:     a.role,
    status:   a.status,
    source:   a.source ?? '—',
    applied:  fmtShortDate(a.applied_at),
    appliedDate: a.applied_at,
    cv:       a.cv_variant ?? '—',
    location: a.location ?? '',
    raw:      a
  };
}

export function fmtShortDate(d) {
  if (!d) return '—';
  return new Date(d).toLocaleDateString('en-US', { day: 'numeric', month: 'short' });
}

export function fmtLongDate(d) {
  if (!d) return '—';
  return new Date(d).toLocaleDateString('en-US', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' });
}

export function fmtWeekday(d) {
  if (!d) return '';
  return new Date(d).toLocaleDateString('en-US', { weekday: 'short' });
}

export function countsByStatus(apps) {
  return apps.reduce((acc, a) => {
    acc[a.status] = (acc[a.status] || 0) + 1;
    return acc;
  }, {});
}
