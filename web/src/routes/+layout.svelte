<script>
  import '../app.css';
  let { children } = $props();
</script>

<svelte:head>
  <title>Pursuit</title>
</svelte:head>

{@render children()}
