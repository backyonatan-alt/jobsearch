// Mapping helpers that translate backend Application objects into the shape
// expected by the Crystalline Console design system (logo class, short label,
// formatted date, etc.). Keeps presentation logic out of the components.

export const STATUSES = [
  'wishlist', 'applied', 'screen', 'interview', 'offer', 'rejected', 'withdrawn'
];

// Common application sources, offered as a dropdown (via <datalist>) while
// keeping the field free-text so anything not listed still works.
export const SOURCE_SUGGESTIONS = [
  'LinkedIn',
  'Company website',
  'Referral',
  'Recruiter reached out',
  'Cold outreach',
  'Job board',
  'Other'
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
    last_follow_up_at: a.last_follow_up_at ?? null,
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

// Hosts that are job boards / ATS systems, not the hiring company. When the
// JD URL points to one of these, the URL host is useless for branding the
// row (we'd render LinkedIn's logo for every LinkedIn-posted job), so we
// skip the URL and fall back to the company-name slug.
const JOB_BOARD_HOSTS = new Set([
  'linkedin.com', 'indeed.com', 'glassdoor.com', 'monster.com', 'ziprecruiter.com',
  'wellfound.com', 'angel.co', 'builtin.com', 'simplyhired.com', 'dice.com',
  'greenhouse.io', 'lever.co', 'ashbyhq.com', 'ashby.com',
  'smartrecruiters.com', 'workable.com', 'breezy.hr', 'recruitee.com',
  'bamboohr.com', 'jobvite.com', 'icims.com', 'taleo.net', 'successfactors.com',
  'workday.com', 'myworkdayjobs.com', 'oraclecloud.com',
  'jobboards.greenhouse.io', 'job-boards.greenhouse.io', 'boards.greenhouse.io',
  'careers.google.com'
]);

function isJobBoardHost(host) {
  if (JOB_BOARD_HOSTS.has(host)) return true;
  // myworkdayjobs.com and similar use customer subdomains — match the suffix.
  for (const suffix of ['.myworkdayjobs.com', '.greenhouse.io', '.lever.co', '.ashbyhq.com', '.workable.com', '.icims.com', '.taleo.net']) {
    if (host.endsWith(suffix)) return true;
  }
  return false;
}

// Best-effort domain for a company so we can fetch a real favicon. Prefers the
// JD url host (drops "www." and any "jobs."/"careers." subdomain) WHEN the
// host is a real company site, not a job board. Otherwise — and as the
// fallback — uses a lowercased company-name slug + ".com".
export function companyDomain(company, jdUrl) {
  if (jdUrl) {
    try {
      const host = new URL(jdUrl).hostname.toLowerCase().replace(/^www\./, '');
      if (!isJobBoardHost(host)) {
        return host.replace(/^(jobs|careers|boards|apply)\./, '');
      }
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

// "Stale" = no activity for over a week on an early-pipeline app. "Activity"
// is the LATER of applied_at and the most recent logged follow-up — logging a
// follow-up resets the clock even though the displayed "Applied X ago" stays
// pinned to applied_at.
export function isStale(a) {
  if (!a) return false;
  if (!['applied', 'screen'].includes(a.status)) return false;
  const appliedAt = a.applied_at ?? a.appliedDate ?? a.raw?.applied_at ?? null;
  const followUpAt = a.last_follow_up_at ?? a.raw?.last_follow_up_at ?? null;
  const basis = mostRecentIso(appliedAt, followUpAt);
  const d = daysSince(basis);
  return d !== null && d >= 7;
}

// Returns the later of two ISO date strings (ignoring nulls). Null if both null.
function mostRecentIso(a, b) {
  if (!a) return b ?? null;
  if (!b) return a;
  return new Date(a).getTime() >= new Date(b).getTime() ? a : b;
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
