<script>
  // The Board was removed in the Jul 14 redesign — the full list (with archive)
  // lives at /app/applications. Keep the route so old links don't 404.
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  onMount(() => goto('/app/applications', { replaceState: true }));
</script>
