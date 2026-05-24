<script>
  import { onMount } from 'svelte';
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api.js';

  let { children } = $props();
  let me = $state(null);
  let applications = $state([]);
  let loading = $state(true);

  onMount(async () => {
    try {
      me = await api('/api/me');
      applications = await api('/api/applications');
    } catch (e) {
      if (e.message !== 'unauthorized') console.error(e);
    } finally {
      loading = false;
    }
  });

  async function signOut() {
    await api('/api/auth/logout', { method: 'POST' });
    goto('/', { replaceState: true });
  }

  const path = $derived(page.url.pathname);
  function isCurrent(href, exact = false) {
    if (exact) return path === href;
    return path === href || path.startsWith(href + '/');
  }

  const counts = $derived.by(() => {
    const c = { active: 0, interview: 0, offer: 0, wishlist: 0, screen: 0, closed: 0 };
    for (const a of applications) {
      if (a.status === 'wishlist') c.wishlist++;
      else if (a.status === 'interview') c.interview++;
      else if (a.status === 'screen') c.screen++;
      else if (a.status === 'offer') c.offer++;
      if (['rejected', 'withdrawn'].includes(a.status)) c.closed++;
      else c.active++;
    }
    return c;
  });

  function navWithFilter(filter) {
    goto(`/app?filter=${filter}`);
  }

  const userInitials = $derived(me?.email
    ? me.email.split('@')[0].slice(0, 2).toUpperCase()
    : '—');
</script>

<div class="app">
  <aside class="sidebar">
    <div class="brand">
      <span class="mark"></span>
      <span class="name">Pursuit</span>
      <span class="kbd">⌘K</span>
    </div>

    <a class="nav-item" class:active={isCurrent('/app', true)} href="/app">
      <span class="nav-icon">
        <svg viewBox="0 0 16 16" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round">
          <rect x="2" y="2.5" width="4" height="4" rx="1"/>
          <rect x="9" y="2.5" width="5" height="4" rx="1"/>
          <rect x="2" y="9" width="4" height="4.5" rx="1"/>
          <rect x="9" y="9" width="5" height="4.5" rx="1"/>
        </svg>
      </span>
      <span>Today</span>
      <span class="nav-count"></span>
    </a>
    <a class="nav-item" class:active={isCurrent('/app/board')} href="/app/board">
      <span class="nav-icon">
        <svg viewBox="0 0 16 16" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round">
          <rect x="2" y="3" width="3" height="10" rx="0.5"/>
          <rect x="7" y="3" width="3" height="7" rx="0.5"/>
          <rect x="12" y="3" width="2" height="5" rx="0.5"/>
        </svg>
      </span>
      <span>Board</span>
      <span class="nav-count">{applications.length || ''}</span>
    </a>
    <a class="nav-item" class:active={isCurrent('/app/funnel')} href="/app/funnel">
      <span class="nav-icon">
        <svg viewBox="0 0 16 16" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round">
          <path d="M2 3h12l-5 6v5l-2-1V9z"/>
        </svg>
      </span>
      <span>Funnel</span>
      <span class="nav-count"></span>
    </a>

    <div class="divider"></div>
    <div class="nav-section">Pipeline</div>
    <a class="nav-item" href="/app?filter=interview">
      <span class="nav-icon">
        <svg viewBox="0 0 16 16" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M3 8a5 5 0 0 1 9-3M13 8a5 5 0 0 1-9 3M11.5 4v2H9.5M4.5 12v-2h2"/>
        </svg>
      </span>
      <span>Interview loops</span>
      <span class="nav-count">{counts.interview}</span>
    </a>
    <a class="nav-item" href="/app?filter=offer">
      <span class="nav-icon">
        <svg viewBox="0 0 16 16" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M8 1.5l1.7 3.4 3.8.5-2.75 2.7.65 3.8L8 10l-3.4 1.9.65-3.8L2.5 5.4l3.8-.5z"/>
        </svg>
      </span>
      <span>Open offers</span>
      <span class="nav-count">{counts.offer}</span>
    </a>
    <a class="nav-item" href="/app?filter=wishlist">
      <span class="nav-icon">
        <svg viewBox="0 0 16 16" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M8 2.5l1.7 3.4 3.8.55-2.75 2.7.65 3.85L8 11.1l-3.4 1.9.65-3.85L2.5 6.45l3.8-.55z"/>
        </svg>
      </span>
      <span>Wishlist</span>
      <span class="nav-count">{counts.wishlist}</span>
    </a>
    <a class="nav-item" href="/app?filter=closed">
      <span class="nav-icon">
        <svg viewBox="0 0 16 16" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.5">
          <circle cx="8" cy="8" r="5.5"/>
          <path d="M5 8h6"/>
        </svg>
      </span>
      <span>Closed</span>
      <span class="nav-count">{counts.closed}</span>
    </a>

    <div class="sidebar-footer">
      <button class="profile" onclick={signOut} title="Sign out">
        <span class="av">{userInitials}</span>
        <span class="who">
          {me?.email?.split('@')[0] ?? 'Signed in'}
          <small>{me?.email ?? ''}</small>
        </span>
        <svg width="14" height="14" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5" style="opacity:0.5">
          <path d="M6 4l4 4-4 4" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      </button>
    </div>
  </aside>

  <section class="main">
    {@render children()}
  </section>
</div>

<style>
  /* Disabled-looking nav items for the not-yet-built screens. */
  :global(.sidebar .nav-item[aria-disabled="true"]) {
    cursor: not-allowed;
    color: var(--mute);
  }
  :global(.sidebar .nav-item[aria-disabled="true"]:hover) {
    background: transparent;
  }
</style>
