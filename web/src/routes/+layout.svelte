<script>
  import '$lib/design-system.css';
  import { afterNavigate } from '$app/navigation';
  import { trackPageView } from '$lib/analytics.js';
  let { children } = $props();

  // Fires on the initial load (type 'enter') and every client-side navigation.
  afterNavigate((nav) => {
    trackPageView(nav?.to?.url?.pathname ?? (typeof location !== 'undefined' ? location.pathname : '/'));
  });
</script>

<svelte:head>
  <title>Pursuit</title>
</svelte:head>

{@render children()}
