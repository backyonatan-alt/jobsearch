<script>
  import { onMount } from 'svelte';
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api.js';
  import { isPreview } from '$lib/preview-mode.js';
  import { track } from '$lib/analytics.js';
  import GuidedTour from '$lib/GuidedTour.svelte';

  let { children } = $props();
  let me = $state(null);
  let applications = $state([]);
  let loading = $state(true);
  let previewMode = $state(false);
  let tourActive = $state(false);
  let tourDismissed = $state(false);
  $effect(() => { previewMode = isPreview(); });

  onMount(async () => {
    try {
      me = await api('/api/me');
      applications = await api('/api/applications');
      // GA4 recommended `login` event. Fired once per browser session on the
      // first authenticated app load — the closest client signal to a real
      // sign-in, since OAuth success is a server-side redirect. Skipped in
      // preview (no real auth).
      if (!isPreview() && typeof sessionStorage !== 'undefined'
          && !sessionStorage.getItem('pursuit_login_tracked')) {
        sessionStorage.setItem('pursuit_login_tracked', '1');
        track('login', { method: 'google' });
      }
    } catch (e) {
      if (e.message !== 'unauthorized') console.error(e);
    } finally {
      loading = false;
    }
  });

  // First-run guided tour: force via ?tour=1, otherwise first run (no onboarded_at).
  const forceTour = $derived(page.url.searchParams.get('tour') === '1');
  $effect(() => {
    if (tourDismissed) { tourActive = false; return; }
    // Latch on: once the tour starts it stays mounted until finished/skipped,
    // so it can't blink off if a navigation changes the trigger conditions.
    if (tourActive) return;
    tourActive = forceTour || (!loading && me != null && !me.onboarded_at);
  });

  async function finishTour() {
    tourDismissed = true;
    tourActive = false;
    try { await api('/api/me/onboarded', { method: 'POST' }); } catch {}
  }

  async function signOut() {
    await api('/api/auth/logout', { method: 'POST' });
    goto('/', { replaceState: true });
  }

  const path = $derived(page.url.pathname);
  function isCurrent(href, exact = false) {
    if (exact) return path === href;
    return path === href || path.startsWith(href + '/');
  }

  const userInitials = $derived(me?.email
    ? me.email.split('@')[0].slice(0, 2).toUpperCase()
    : '—');
</script>

<div class="app">
  <aside class="sidebar">
    <a class="brand" href="/app">
      <svg class="brand-mark" viewBox="0 0 24 24" width="22" height="22" fill="none" aria-hidden="true">
        <circle cx="12" cy="12" r="9.5" stroke="currentColor" stroke-width="1.4" opacity="0.65"/>
        <circle cx="12" cy="12" r="5.5" stroke="currentColor" stroke-width="1.4" opacity="0.9"/>
        <circle cx="17.5" cy="6.5" r="2.6" fill="currentColor"/>
      </svg>
      <span class="name">Pursuit</span>
    </a>

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
      <span>Insights</span>
      <span class="nav-count"></span>
    </a>

    <div class="sidebar-footer">
      <button class="profile" onclick={signOut} title="Sign out">
        {#if me?.picture_url}
          <img class="av av-img" src={me.picture_url} alt={me.email ?? ''} referrerpolicy="no-referrer" />
        {:else}
          <span class="av">{userInitials}</span>
        {/if}
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
    {#if previewMode}
      <div class="preview-banner">
        <span class="pb-dot"></span>
        <strong>Preview mode</strong>
        <span>· UI-only, no backend. Changes live in this tab and reset on reload.</span>
        <a class="pb-exit" href="?preview=0">Exit</a>
      </div>
    {/if}
    {@render children()}
  </section>
</div>

{#if tourActive}
  <GuidedTour onDone={finishTour} seedDemo={(me != null && !me.onboarded_at) || (page.url.searchParams.get('preview') === '1' && forceTour)} />
{/if}

<style>
  .preview-banner {
    background: var(--warm-tint);
    color: var(--warm-text);
    border-bottom: 1px solid var(--rule);
    padding: 6px 16px;
    font-size: 12.5px;
    display: flex; align-items: center; gap: 8px;
    font-weight: 500;
  }
  .preview-banner .pb-dot { width: 8px; height: 8px; border-radius: 50%; background: var(--warm); flex-shrink: 0; }
  .preview-banner strong { font-weight: 600; }
  .preview-banner .pb-exit { margin-left: auto; color: var(--warm-text); font-weight: 600; padding: 2px 10px; border-radius: 99px; border: 1px solid var(--warm); }
  .preview-banner .pb-exit:hover { background: var(--warm); color: white; }

  /* Brand mark: target-style SVG paired with the wordmark, swapping out the
     old square + dot. The SVG colors itself via currentColor. */
  :global(.sidebar .brand) {
    grid-template-columns: 22px 1fr;
    text-decoration: none;
  }
  :global(.sidebar .brand .brand-mark) { color: var(--accent); }
  :global(.sidebar .brand .name) { color: var(--ink); }

  /* When we have the Google profile picture, render it where the gradient
     initials square used to live. Same dimensions, just an <img>. */
  :global(.sidebar .profile .av.av-img) {
    background: var(--surface-2);
    object-fit: cover;
    padding: 0;
  }

  /* Disabled-looking nav items for the not-yet-built screens. */
  :global(.sidebar .nav-item[aria-disabled="true"]) {
    cursor: not-allowed;
    color: var(--mute);
  }
  :global(.sidebar .nav-item[aria-disabled="true"]:hover) {
    background: transparent;
  }
</style>
