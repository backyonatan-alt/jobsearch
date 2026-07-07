// Thin fetch wrapper for /api/* endpoints. Throws on non-2xx.
// When ?preview=1 is on the URL (see $lib/preview-mode.js), we short-circuit
// to an in-memory mock so the UI is reviewable without a running backend.

import { isPreview, mockApi } from './preview-mode.js';

export async function api(path, opts = {}) {
  if (isPreview()) return mockApi(path, opts);
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

// A dropped connection (phone locks, tab hidden, network blip — the browser's
// "Failed to fetch"/"Load failed") kills the response, not the generation:
// the server detaches and persists the brief anyway. 502/504 mean the proxy
// gave up the same way.
export function isConnectionErr(e) {
  return e instanceof TypeError || /failed to fetch|load failed|networkerror|http 50[24]/i.test(String(e?.message || ''));
}

// Poll a dossier GET until it returns a brief newer than prevGeneratedAt
// (pass null when any brief counts). Resolves null on timeout.
export async function pollForDossier(path, prevGeneratedAt = null, { timeoutMs = 150000, everyMs = 6000 } = {}) {
  const deadline = Date.now() + timeoutMs;
  while (Date.now() < deadline) {
    await new Promise((res) => setTimeout(res, everyMs));
    try {
      const d = await api(path);
      if (d && d.generated_at !== prevGeneratedAt) return d;
    } catch { /* 404 → not ready yet */ }
  }
  return null;
}

export const STATUSES = [
  'wishlist',
  'applied',
  'screen',
  'interview',
  'offer',
  'rejected',
  'withdrawn',
  'closed'
];
