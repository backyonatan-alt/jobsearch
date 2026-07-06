<script>
  import { onMount } from 'svelte';
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api.js';
  import { isPreview } from '$lib/preview-mode.js';
  import { track, logEvent } from '$lib/analytics.js';
  import GuidedTour from '$lib/GuidedTour.svelte';
  import PrepFirstOnboarding from '$lib/PrepFirstOnboarding.svelte';

  // Beta feedback channel — a pre-addressed email so first users have an
  // obvious way to send notes (Michal's first ask).
  const FEEDBACK_EMAIL = 'back.yonatan@gmail.com';
  const feedbackHref = `mailto:${FEEDBACK_EMAIL}?subject=${encodeURIComponent('Pursuit feedback')}&body=${encodeURIComponent("What I was doing:\n\n\nWhat happened / what I'd expect instead:\n\n")}`;

  let { children } = $props();
  let me = $state(null);
  let applications = $state([]);
  let loading = $state(true);
  let previewMode = $state(false);
  let onboardMode = $state(null); // 'tour' | 'prep' | null
  let onboardDismissed = $state(false);
  let variantTracked = false;
  // Desktop-only notice for the beta (CSS hides it ≥820px). Dismissal sticks.
  let narrowDismissed = $state(true);
  function dismissNarrow() {
    narrowDismissed = true;
    try { localStorage.setItem('pursuit_narrow_dismissed', '1'); } catch {}
  }
  $effect(() => { previewMode = isPreview(); });

  // Keep the sidebar count (applications.length) fresh: any mutation anywhere
  // dispatches `pursuit:refresh`, and we also refetch on tab focus. Without this
  // the count was fetched once on mount and went stale after add/delete.
  async function refreshApplications() {
    try { applications = await api('/api/applications'); } catch (e) {
      if (e.message !== 'unauthorized') console.error(e);
    }
  }

  async function init() {
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
        let src = '';
        try { src = localStorage.getItem('pursuit_src') || ''; } catch {}
        track('login', src ? { method: 'google', src } : { method: 'google' });
        logEvent('login', src ? { src } : {});
      }
    } catch (e) {
      if (e.message !== 'unauthorized') console.error(e);
    } finally {
      loading = false;
    }
  }

  // onMount stays synchronous so the returned cleanup actually runs.
  onMount(() => {
    init();
    try { narrowDismissed = localStorage.getItem('pursuit_narrow_dismissed') === '1'; } catch { narrowDismissed = false; }
    const onRefresh = () => refreshApplications();
    const onVis = () => { if (document.visibilityState === 'visible') refreshApplications(); };
    window.addEventListener('pursuit:refresh', onRefresh);
    window.addEventListener('focus', onRefresh);
    document.addEventListener('visibilitychange', onVis);
    return () => {
      window.removeEventListener('pursuit:refresh', onRefresh);
      window.removeEventListener('focus', onRefresh);
      document.removeEventListener('visibilitychange', onVis);
    };
  });

  // First-run onboarding. The variant (prep-first cold start vs guided tour) comes
  // from me.onboarding_variant; ?tour=1 / ?onboard=prepfirst force one for QA.
  const forceTour = $derived(page.url.searchParams.get('tour') === '1');
  const forcePrep = $derived(page.url.searchParams.get('onboard') === 'prepfirst');
  $effect(() => {
    if (onboardDismissed) { onboardMode = null; return; }
    // Latch on: once onboarding starts it stays mounted until finished/skipped,
    // so it can't blink off if a navigation changes the trigger conditions.
    if (onboardMode) return;
    const needs = forceTour || forcePrep || (!loading && me != null && !me.onboarded_at);
    if (!needs) return;
    const variant = forcePrep ? 'prepfirst' : forceTour ? 'tour' : me?.onboarding_variant;
    onboardMode = variant === 'prepfirst' ? 'prep' : 'tour';
    if (!variantTracked && !isPreview()) {
      variantTracked = true;
      logEvent('onboard_variant_assigned', { variant: onboardMode === 'prep' ? 'prepfirst' : 'tour' });
    }
  });

  async function finishOnboarding() {
    onboardDismissed = true;
    onboardMode = null;
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
      <a class="nav-item feedback-link" href={feedbackHref} onclick={() => logEvent('feedback_click', { surface: 'sidebar' })}>
        <span class="nav-icon">
          <svg viewBox="0 0 16 16" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round">
            <path d="M3 3h10a1 1 0 0 1 1 1v6a1 1 0 0 1-1 1H7l-3 2.5V11H3a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1z"/>
          </svg>
        </span>
        <span>Send feedback</span>
        <span class="nav-count"></span>
      </a>
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
    {#if !narrowDismissed}
      <div class="narrow-note" role="note">
        <span><strong>Tip:</strong> the full board and editing experience is roomier on a laptop — but prep works great right here.</span>
        <button class="nn-x" onclick={dismissNarrow} aria-label="Dismiss">✕</button>
      </div>
    {/if}
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

{#if onboardMode === 'tour'}
  <GuidedTour onDone={finishOnboarding} seedDemo={(me != null && !me.onboarded_at) || (page.url.searchParams.get('preview') === '1' && forceTour)} />
{:else if onboardMode === 'prep'}
  <PrepFirstOnboarding onDone={finishOnboarding} />
{/if}

<style>
  /* Desktop-only beta notice — only shown on narrow viewports. */
  .narrow-note { display: none; }
  @media (max-width: 820px) {
    .narrow-note {
      display: flex; align-items: center; gap: 10px;
      background: var(--warm-tint); color: var(--warm-text);
      border-bottom: 1px solid var(--rule);
      padding: 9px 14px; font-size: 12.5px; line-height: 1.4;
    }
    .narrow-note strong { font-weight: 600; }
    .narrow-note .nn-x { margin-left: auto; flex-shrink: 0; background: none; border: none;
      color: var(--warm-text); font-size: 13px; cursor: pointer; padding: 2px 6px; border-radius: 5px; }
  }

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
