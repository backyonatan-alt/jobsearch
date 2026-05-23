<script>
  import '$lib/design-system.css';
  let { children } = $props();
</script>

<svelte:head>
  <title>Pursuit</title>
</svelte:head>

{@render children()}
