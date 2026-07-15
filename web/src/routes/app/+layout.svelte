<script>
  import { onMount } from 'svelte';
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api.js';
  import { isPreview } from '$lib/preview-mode.js';
  import { track, logEvent } from '$lib/analytics.js';
  import GuidedTour from '$lib/GuidedTour.svelte';
  import PrepFirstOnboarding from '$lib/PrepFirstOnboarding.svelte';
  import AddApplication from '$lib/AddApplication.svelte';

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
  let showNewModal = $state(false);
  // Desktop-only notice for the beta (CSS hides it ≥820px). Dismissal sticks.
  let narrowDismissed = $state(true);
  function dismissNarrow() {
    narrowDismissed = true;
    try { localStorage.setItem('pursuit_narrow_dismissed', '1'); } catch {}
  }
  $effect(() => { previewMode = isPreview(); });

  // Keep the nav count (applications.length) fresh: any mutation anywhere
  // dispatches `pursuit:refresh`, and we also refetch on tab focus.
  async function refreshApplications() {
    try { applications = await api('/api/applications'); } catch (e) {
      if (e.message !== 'unauthorized') console.error(e);
    }
  }
  async function refreshMe() {
    try { me = await api('/api/me'); } catch (e) {
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
    const onRefresh = () => { refreshApplications(); refreshMe(); };
    const onVis = () => { if (document.visibilityState === 'visible') onRefresh(); };
    const onNewApp = () => { showNewModal = true; };
    window.addEventListener('pursuit:refresh', onRefresh);
    window.addEventListener('pursuit:new-app', onNewApp);
    window.addEventListener('focus', onRefresh);
    document.addEventListener('visibilitychange', onVis);
    return () => {
      window.removeEventListener('pursuit:refresh', onRefresh);
      window.removeEventListener('pursuit:new-app', onNewApp);
      window.removeEventListener('focus', onRefresh);
      document.removeEventListener('visibilitychange', onVis);
    };
  });

  // ⌘N / Ctrl+N opens the new-application modal from anywhere in the app.
  function onKeydown(e) {
    if ((e.metaKey || e.ctrlKey) && (e.key === 'n' || e.key === 'N') && !showNewModal) {
      e.preventDefault();
      showNewModal = true;
    }
  }

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
  // Detail pages carry their own primary CTA — New application drops to the
  // outline variant there (design CTA rule: one primary per page).
  const onDetail = $derived(/^\/app\/\d+/.test(path));

  const userInitials = $derived(me?.email
    ? me.email.split('@')[0].slice(0, 2).toUpperCase()
    : '—');
  const creditsUsed = $derived(me?.prep_credits_used ?? 0);
  const creditsLimit = $derived(me?.prep_credits_limit ?? 10);
  const creditsLeft = $derived(Math.max(0, creditsLimit - creditsUsed));
</script>

<svelte:window onkeydown={onKeydown} />

<div class="shell">
  <header class="topnav">
    <div class="tn-in">
      <a class="brand" href="/app">
        <svg class="brand-mark" viewBox="0 0 24 24" width="24" height="24" fill="none" aria-hidden="true">
          <circle cx="12" cy="12" r="10" stroke="#16181c" stroke-width="2.5"/>
          <circle cx="12" cy="12" r="4.5" fill="#2463eb"/>
        </svg>
        <span class="name">Pursuit</span>
      </a>

      <nav class="pills">
        <a class="pill-link" class:active={isCurrent('/app', true)} href="/app">Home</a>
        <a class="pill-link" class:active={isCurrent('/app/applications') || onDetail} href="/app/applications">
          Applications {#if applications.length}<span class="ct">{applications.length}</span>{/if}
        </a>
        <a class="pill-link" class:active={isCurrent('/app/funnel')} href="/app/funnel">Insights</a>
      </nav>

      <div class="tn-right">
        <div class="search" role="presentation">
          <span class="ico">⌕</span>Search…
          <span class="kbd">⌘K</span>
        </div>
        <button class="new-app" class:outline={onDetail} data-tour="new-app" onclick={() => (showNewModal = true)}>
          New application <span class="nk">⌘N</span>
        </button>
        {#if creditsLimit <= 1000}
          <div class="credits" title="Prep credits — 1 per generated round brief">
            ✦ <strong>{creditsLeft}</strong><span class="of">/{creditsLimit}</span>
          </div>
        {/if}
        <a class="fb" href={feedbackHref} title="Send feedback" onclick={() => logEvent('feedback_click', { surface: 'topnav' })}>
          <svg viewBox="0 0 16 16" width="15" height="15" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round">
            <path d="M3 3h10a1 1 0 0 1 1 1v6a1 1 0 0 1-1 1H7l-3 2.5V11H3a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1z"/>
          </svg>
        </a>
        <button class="avatar" onclick={signOut} title={me?.email ? `${me.email} — sign out` : 'Sign out'}>
          {#if me?.picture_url}
            <img src={me.picture_url} alt={me.email ?? ''} referrerpolicy="no-referrer" />
          {:else}
            {userInitials}
          {/if}
        </button>
      </div>
    </div>
  </header>

  <section class="main">
    {#if !narrowDismissed}
      <div class="narrow-note" role="note">
        <span><strong>Tip:</strong> the full editing experience is roomier on a laptop — but prep works great right here.</span>
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

<AddApplication bind:open={showNewModal} onCreated={refreshApplications} />

{#if onboardMode === 'tour'}
  <GuidedTour onDone={finishOnboarding} seedDemo={(me != null && !me.onboarded_at) || (page.url.searchParams.get('preview') === '1' && forceTour)} />
{:else if onboardMode === 'prep'}
  <PrepFirstOnboarding onDone={finishOnboarding} />
{/if}

<style>
  .shell {
    min-height: 100vh; display: flex; flex-direction: column;
    background: #f6f6f3; color: #16181c;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
  }
  .main { flex: 1; display: flex; flex-direction: column; min-height: 0; }

  .topnav { background: #fff; border-bottom: 1px solid #e8e8e5; flex-shrink: 0; }
  .tn-in {
    max-width: 1160px; margin: 0 auto; display: flex; align-items: center;
    gap: 20px; padding: 12px 32px;
  }
  .brand { display: flex; align-items: center; gap: 9px; flex: none; text-decoration: none; color: #16181c; }
  .brand .name { font-size: 16px; font-weight: 700; letter-spacing: -0.01em; }

  .pills { display: flex; align-items: center; gap: 4px; }
  .pill-link {
    display: flex; align-items: center; gap: 6px; border-radius: 8px;
    padding: 7px 14px; font-size: 13.5px; color: #4b5158; text-decoration: none;
    white-space: nowrap;
  }
  .pill-link:hover { background: #f0f0ed; }
  .pill-link.active { background: #16181c; color: #fff; font-weight: 600; }
  .pill-link .ct { font-size: 11.5px; color: #8a9099; }
  .pill-link.active .ct { color: #fff; opacity: .6; }

  .tn-right { display: flex; align-items: center; gap: 10px; margin-left: auto; min-width: 0; }
  .search {
    display: flex; align-items: center; gap: 8px; border: 1px solid #e8e8e5;
    border-radius: 8px; padding: 7px 12px; width: 210px; color: #8a9099; font-size: 13px;
  }
  .search .kbd { margin-left: auto; font-size: 11px; border: 1px solid #e8e8e5; border-radius: 4px; padding: 1px 5px; }
  .new-app {
    background: #2463eb; color: #fff; border: 1px solid #2463eb; border-radius: 8px;
    padding: 8px 15px; font-size: 13px; font-weight: 600; cursor: pointer;
    flex: none; font-family: inherit; white-space: nowrap;
  }
  .new-app:hover { background: #1a4fc4; }
  .new-app .nk { opacity: .65; font-weight: 400; margin-left: 3px; }
  .new-app.outline { background: #fff; color: #4b5158; border-color: #e8e8e5; }
  .new-app.outline:hover { border-color: #b9c6e8; color: #2463eb; background: #fff; }
  .new-app.outline .nk { opacity: .55; }
  .credits { border: 1px solid #e8e8e5; border-radius: 8px; padding: 7px 11px; font-size: 12px; color: #6f7680; flex: none; white-space: nowrap; }
  .credits strong { color: #16181c; }
  .credits .of { color: #8a9099; }
  .fb { color: #8a9099; display: flex; padding: 6px; border-radius: 6px; }
  .fb:hover { color: #2463eb; background: #f0f0ed; }
  .avatar {
    width: 30px; height: 30px; border-radius: 50%; background: #e0641f; color: #fff;
    display: flex; align-items: center; justify-content: center; font-size: 12px;
    font-weight: 700; flex: none; border: 0; cursor: pointer; padding: 0; overflow: hidden;
    font-family: inherit;
  }
  .avatar img { width: 100%; height: 100%; object-fit: cover; }

  /* Desktop-only beta notice — only shown on narrow viewports. */
  .narrow-note { display: none; }
  @media (max-width: 820px) {
    .narrow-note {
      display: flex; align-items: center; gap: 10px;
      background: #fff7f1; color: #c05310;
      border-bottom: 1px solid #f0d9c4;
      padding: 9px 14px; font-size: 12.5px; line-height: 1.4;
    }
    .narrow-note strong { font-weight: 600; }
    .narrow-note .nn-x { margin-left: auto; flex-shrink: 0; background: none; border: none;
      color: #c05310; font-size: 13px; cursor: pointer; padding: 2px 6px; border-radius: 5px; }
  }

  .preview-banner {
    background: #fff7f1; color: #c05310; border-bottom: 1px solid #f0d9c4;
    padding: 6px 16px; font-size: 12.5px;
    display: flex; align-items: center; gap: 8px; font-weight: 500;
  }
  .preview-banner .pb-dot { width: 8px; height: 8px; border-radius: 50%; background: #e0641f; flex-shrink: 0; }
  .preview-banner strong { font-weight: 600; }
  .preview-banner .pb-exit { margin-left: auto; color: #c05310; font-weight: 600; padding: 2px 10px; border-radius: 99px; border: 1px solid #e0641f; text-decoration: none; }
  .preview-banner .pb-exit:hover { background: #e0641f; color: white; }

  /* Narrow screens: hide the search stub, tighten paddings. */
  @media (max-width: 900px) {
    .tn-in { padding: 10px 14px; gap: 10px; flex-wrap: wrap; }
    .search { display: none; }
    .credits { display: none; }
    .brand .name { display: none; }
  }
</style>
