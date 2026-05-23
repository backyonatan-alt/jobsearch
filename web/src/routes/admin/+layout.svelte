<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api.js';

  let { children } = $props();
  let ready = $state(false);
  let allowed = $state(false);
  let me = $state(null);

  onMount(async () => {
    try {
      me = await api('/api/me');
      allowed = !!me?.is_admin;
    } catch (e) {
      if (e.message === 'unauthorized') return;
      console.error(e);
    } finally {
      ready = true;
      if (!allowed) {
        // Quiet redirect for non-admins. Nothing here for you.
        setTimeout(() => goto('/app', { replaceState: true }), 50);
      }
    }
  });
</script>

<svelte:head><title>Admin — Pursuit</title></svelte:head>

<div class="admin-shell">
  <header class="admin-top">
    <div class="brand">
      <span class="mark"></span>
      <span class="name">Pursuit</span>
      <span style="color: var(--mute-2); font-size: 11px; margin-left: 6px;">admin</span>
    </div>
    <div class="right">
      <a class="btn" href="/app">‹ back to app</a>
    </div>
  </header>

  <main class="admin-body">
    {#if !ready}
      <p style="color:var(--mute)">Loading…</p>
    {:else if !allowed}
      <p style="color:var(--mute)">Not authorized. Redirecting…</p>
    {:else}
      {@render children()}
    {/if}
  </main>
</div>

<style>
  .admin-shell { min-height: 100vh; display: flex; flex-direction: column; background: var(--surface); }
  .admin-top {
    height: 48px; padding: 0 24px;
    display: flex; align-items: center; justify-content: space-between;
    border-bottom: 1px solid var(--rule);
    background: var(--surface);
  }
  .brand { display: flex; align-items: center; gap: 10px; }
  .brand .mark {
    width: 20px; height: 20px; border-radius: 5px;
    background: var(--ink); position: relative;
  }
  .brand .mark::after {
    content: ""; position: absolute; inset: 5px 5px auto auto;
    width: 6px; height: 6px; border-radius: 50%; background: var(--accent);
  }
  .brand .name { font-weight: 600; font-size: 14px; letter-spacing: -0.015em; }
  .admin-body { padding: 32px; max-width: 720px; width: 100%; margin: 0 auto; }
</style>
