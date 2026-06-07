// Thin GA4 wrapper. No-ops when gtag is absent — which is the case in dev and
// on any deploy without GA4_MEASUREMENT_ID set (the snippet is injected
// server-side only when the ID is configured). So callers never need to guard.

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
