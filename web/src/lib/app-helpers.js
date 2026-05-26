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
  const domain = companyDomain(a.company, a.jd_url);
  return {
    id:       a.id,
    co:       a.company,
    coShort:  first,
    logoCls:  `lg-${logoLetter}`,
    role:     a.role,
    status:   a.status,
    source:   a.source ?? '—',
    applied:  fmtShortDate(a.applied_at),
    appliedRel: fmtRelativeDate(a.applied_at),
    appliedDate: a.applied_at,
    cv:       a.cv_variant ?? '—',
    location: a.location ?? '',
    domain,
    logoSrc:  domain ? `https://www.google.com/s2/favicons?sz=128&domain=${domain}` : '',
    stale:    isStale(a),
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

// Best-effort domain for a company so we can fetch a real favicon. Prefers the
// JD url host (drops "www." and any "jobs."/"careers." subdomain), else falls
// back to a lowercased company-name slug + ".com".
export function companyDomain(company, jdUrl) {
  if (jdUrl) {
    try {
      const h = new URL(jdUrl).hostname.toLowerCase();
      return h.replace(/^www\./, '').replace(/^(jobs|careers|boards|apply)\./, '');
    } catch {}
  }
  const slug = String(company || '').toLowerCase()
    .replace(/[^a-z0-9]+/g, '')
    .replace(/^(the|a|an)/, '');
  return slug ? `${slug}.com` : '';
}

export function faviconUrl(company, jdUrl, size = 128) {
  const d = companyDomain(company, jdUrl);
  return d ? `https://www.google.com/s2/favicons?sz=${size}&domain=${d}` : '';
}

// Days since an ISO date string. Null when missing.
export function daysSince(iso) {
  if (!iso) return null;
  const ms = Date.now() - new Date(iso).getTime();
  return Math.floor(ms / (1000 * 60 * 60 * 24));
}

// "Stale" = no activity for over a week on an early-pipeline app.
export function isStale(a) {
  if (!a) return false;
  if (!['applied', 'screen'].includes(a.status)) return false;
  const d = daysSince(a.applied_at ?? a.appliedDate ?? a.raw?.applied_at);
  return d !== null && d >= 7;
}

// "X days ago" style relative date, matching the locked design copy.
export function fmtRelativeDate(iso) {
  const d = daysSince(iso);
  if (d === null) return '—';
  // Server-side applied_at can land a few ms after the client clock, which
  // turns Math.floor into -1. Treat anything ≤ 0 days as "today".
  if (d <= 0) return 'today';
  if (d === 1) return 'yesterday';
  if (d < 30) return `${d} days ago`;
  const months = Math.floor(d / 30);
  return `${months} ${months === 1 ? 'month' : 'months'} ago`;
}
