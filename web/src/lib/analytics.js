// Two analytics sinks:
//  - track()/trackPageView() → GA4, for the public homepage / acquisition.
//  - logEvent() → first-party `events` table (POST /api/events), for in-app
//    product behaviour. First-party is the source of truth at this small N
//    (GA4 thresholds low-volume data); GA4 stays marketing-only.
// PII rule for logEvent props: enums, booleans, counts, durations only — never
// email/company/role/JD/dossier text.

import { isPreview } from './preview-mode.js';

function gtagReady() {
  return typeof window !== 'undefined' && typeof window.gtag === 'function';
}

export function track(name, params = {}) {
  if (!gtagReady()) return;
  window.gtag('event', name, params);
}

// GA4 is configured with send_page_view:false (SPA — it can't see client-side
// route changes), so we fire page_view ourselves on every navigation.
export function trackPageView(path) {
  if (!gtagReady()) return;
  window.gtag('event', 'page_view', { page_path: path });
}

// First-party in-app event → POST /api/events. Fire-and-forget; never throws,
// never blocks the UI action, and no-ops in preview (mock) mode.
export function logEvent(name, props = {}) {
  if (typeof window === 'undefined' || isPreview()) return;
  try {
    fetch('/api/events', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name, props }),
      keepalive: true
    }).catch(() => {});
  } catch {
    /* analytics must never break the app */
  }
}
