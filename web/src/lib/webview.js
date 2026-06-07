// Detects in-app / embedded browsers (LinkedIn, Instagram, Facebook, etc.).
// Google blocks "Sign in with Google" inside embedded webviews with
// `Error 403: disallowed_useragent`, so on these we steer users to open the
// page in a real browser before they hit the OAuth dead-end.

const IN_APP_TOKENS = [
  'FBAN', 'FBAV', 'FB_IAB',                    // Facebook / Messenger
  'Instagram',                                  // Instagram
  'LinkedInApp',                                // LinkedIn
  'Line/',                                       // LINE
  'Twitter',                                     // X / Twitter
  'Snapchat',                                    // Snapchat
  'TikTok', 'musical_ly', 'BytedanceWebview',   // TikTok
  'WhatsApp'                                     // WhatsApp
];

export function isInAppBrowser() {
  if (typeof navigator === 'undefined') return false;
  const ua = navigator.userAgent || '';
  if (IN_APP_TOKENS.some((t) => ua.includes(t))) return true;
  // Generic Android WebView (e.g. apps embedding a raw WebView).
  if (/Android/.test(ua) && /;\s*wv[);]/.test(ua)) return true;
  return false;
}
