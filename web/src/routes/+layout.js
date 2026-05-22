// SPA mode: all rendering happens in the browser. Static adapter emits a
// single index.html fallback; the Svelte router takes it from there.
export const ssr = false;
export const prerender = false;
export const trailingSlash = 'never';
