// Thin fetch wrapper for /api/* endpoints. Throws on non-2xx.

export async function api(path, opts = {}) {
  const r = await fetch(path, {
    headers: { 'Content-Type': 'application/json', ...(opts.headers || {}) },
    ...opts
  });
  if (r.status === 401) {
    if (typeof window !== 'undefined') window.location.href = '/';
    throw new Error('unauthorized');
  }
  if (r.status === 204) return null;
  const body = await r.json().catch(() => ({}));
  if (!r.ok) throw new Error(body.error || `http ${r.status}`);
  return body;
}

export const STATUSES = [
  'wishlist',
  'applied',
  'screen',
  'interview',
  'offer',
  'rejected',
  'withdrawn'
];
