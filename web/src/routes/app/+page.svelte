<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/state';
  import { api } from '$lib/api.js';
  import {
    STATUS_LABEL, toDisplayApp, countsByStatus, daysSince, faviconUrl
  } from '$lib/app-helpers.js';
  import Onboarding from '$lib/Onboarding.svelte';

  let apps = $state([]);
  let me = $state(null);
  let loading = $state(true);
  let filter = $state('active');
  let showNewModal = $state(false);
  let showOnboarding = $state(false);

  // new-application form state
  let form = $state({ company: '', role: '', status: 'applied', source: '', jd_url: '', cv_variant: '', location: '', salary_note: '' });
  let saving = $state(false);

  // paste-to-parse state
  let pasteText = $state('');
  let parsing = $state(false);
  let parseError = $state('');
  let parsedHint = $state('');
  let attachedImage = $state(null);
  let isDraggingFile = $state(false);
  let parsingStage = $state('');

  const ALLOWED_IMG = ['image/png', 'image/jpeg', 'image/gif', 'image/webp'];
  const MAX_IMG_BYTES = 6 * 1024 * 1024;

  function setAttachedFile(f) {
    if (!f) return;
    if (!ALLOWED_IMG.includes(f.type)) {
      parseError = 'Only PNG, JPEG, GIF, or WebP screenshots are supported.';
      return;
    }
    if (f.size > MAX_IMG_BYTES) {
      parseError = 'Screenshot is too large (6 MB max — try a tighter crop).';
      return;
    }
    parseError = '';
    attachedImage = { name: f.name || 'pasted-image.png', size: f.size, mediaType: f.type, file: f };
  }

  function fileToBase64(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onerror = () => reject(new Error('Could not read screenshot.'));
      reader.onload = () => {
        const result = reader.result || '';
        const i = String(result).indexOf(',');
        resolve(i === -1 ? '' : String(result).slice(i + 1));
      };
      reader.readAsDataURL(file);
    });
  }

  async function parseJD() {
    const text = pasteText.trim();
    if (!attachedImage && text.length < 5) {
      parseError = 'Paste a job listing or URL, or drop a screenshot first.';
      return;
    }
    parseError = '';
    parsedHint = '';
    parsing = true;
    parsingStage = attachedImage ? 'Reading your screenshot…' : 'Reading the text…';
    const stageTimers = [
      setTimeout(() => { if (parsing) parsingStage = 'Pulling out company and role…'; }, 1800),
      setTimeout(() => { if (parsing) parsingStage = 'Almost there…'; }, 4200)
    ];
    try {
      const payload = { text };
      if (attachedImage) {
        const data = await fileToBase64(attachedImage.file);
        payload.image = { media_type: attachedImage.mediaType, data };
      }
      const r = await api('/api/applications/parse', { method: 'POST', body: JSON.stringify(payload) });
      if (r.company)     form.company     = r.company;
      if (r.role)        form.role        = r.role;
      if (r.location)    form.location    = r.location;
      if (r.jd_url)      form.jd_url      = r.jd_url;
      if (r.source)      form.source      = r.source;
      if (r.salary_note) form.salary_note = r.salary_note;
      const filled = [];
      if (r.company) filled.push('Company');
      if (r.role)    filled.push('Role');
      if (r.source)  filled.push('Source');
      parsedHint = filled.length
        ? `Filled ${filled.join(', ')}${r.seniority ? ` · level: ${r.seniority}` : ''} — review below.`
        : 'Parsed — review below.';
      pasteText = '';
      attachedImage = null;
    } catch (e) {
      parseError = e.message || 'Parse failed.';
    } finally {
      stageTimers.forEach(clearTimeout);
      parsing = false;
      parsingStage = '';
    }
  }

  function onDropZoneDragOver(e) {
    if (e.dataTransfer?.types?.includes('Files')) {
      e.preventDefault();
      isDraggingFile = true;
    }
  }
  function onDropZoneDragLeave() { isDraggingFile = false; }
  function onDropZoneDrop(e) {
    e.preventDefault();
    isDraggingFile = false;
    setAttachedFile(e.dataTransfer?.files?.[0]);
  }
  function onTextareaPaste(e) {
    const item = [...(e.clipboardData?.items || [])].find(i => i.type.startsWith('image/'));
    if (item) setAttachedFile(item.getAsFile());
  }

  function onModalKeydown(e) { if (e.key === 'Escape') resetModal(); }
  function onModalPaste(e) {
    const item = [...(e.clipboardData?.items || [])].find(i => i.type.startsWith('image/'));
    if (item) setAttachedFile(item.getAsFile());
  }

  function resetModal() {
    form = { company: '', role: '', status: 'applied', source: '', jd_url: '', cv_variant: '', location: '', salary_note: '' };
    pasteText = '';
    parseError = '';
    parsedHint = '';
    attachedImage = null;
    isDraggingFile = false;
    showNewModal = false;
  }

  onMount(refresh);

  $effect(() => {
    const f = page.url.searchParams.get('filter');
    if (f && filterMap[f]) filter = f;
  });

  async function refresh() {
    try {
      const [meRes, raw] = await Promise.all([
        api('/api/me'),
        api('/api/applications')
      ]);
      me = meRes;
      apps = raw.map(toDisplayApp);
      maybeShowOnboarding();
    } catch (e) {
      if (e.message !== 'unauthorized') console.error(e);
    } finally {
      loading = false;
    }
  }

  function maybeShowOnboarding() {
    const forced = page.url.searchParams.get('onboarding') === '1';
    const fresh = !me?.onboarded_at && apps.length === 0;
    showOnboarding = forced || fresh;
  }

  function finishOnboarding() { showOnboarding = false; refresh(); }

  async function createApp(e) {
    e.preventDefault();
    saving = true;
    try {
      const payload = { ...form };
      for (const k of ['jd_url', 'cv_variant', 'source', 'location', 'salary_note']) if (!payload[k]) delete payload[k];
      await api('/api/applications', { method: 'POST', body: JSON.stringify(payload) });
      resetModal();
      await refresh();
    } finally {
      saving = false;
    }
  }

  function open(id) { goto(`/app/${id}`); }

  // ── Demo data ────────────────────────────────────────────
  // Friends signing in for the first time land on a blank /app, which is a
  // terrible first impression. Offer a one-click "seed me with 15 realistic
  // apps" so they can see the full product, with a matching clear action.
  let seeding = $state(false);
  let seedError = $state('');
  async function seedDemo() {
    if (seeding) return;
    seeding = true;
    seedError = '';
    try {
      await api('/api/me/demo-seed', { method: 'POST' });
      await refresh();
    } catch (e) {
      seedError = e.message || 'Could not seed demo data.';
    } finally {
      seeding = false;
    }
  }
  async function clearDemo() {
    if (!confirm('Remove all demo applications? Your real ones are kept.')) return;
    try {
      await api('/api/me/demo-seed', { method: 'DELETE' });
      await refresh();
    } catch (e) {
      console.error('clear demo', e);
    }
  }
  const hasDemoRows = $derived(apps.some(a => (a.raw?.notes || '').startsWith('[demo] ')));

  // ── Derived state ────────────────────────────────────────
  const counts = $derived(countsByStatus(apps));
  const active = $derived(apps.filter(a => !['rejected', 'withdrawn'].includes(a.status)));

  // Count cards — interviews, offers, applied & waiting (no movement yet), wishlist.
  const stale = $derived(apps.filter(a => a.stale));
  const countCards = $derived.by(() => {
    const interview = apps.find(a => a.status === 'interview');
    const offer = apps.find(a => a.status === 'offer');
    const oldestWish = apps
      .filter(a => a.status === 'wishlist')
      .sort((a, b) => (a.appliedDate || '').localeCompare(b.appliedDate || ''))[0];
    const wishDays = oldestWish ? daysSince(oldestWish.raw.created_at) : null;
    return [
      {
        lbl: 'Interviews',
        n: counts.interview || 0,
        tone: 'accent',
        sub: interview ? `${interview.co} · ${STATUS_LABEL[interview.status]}` : 'None scheduled',
        hint: 'An interview loop = a series of rounds with one company'
      },
      {
        lbl: 'Open offers',
        n: counts.offer || 0,
        tone: 'warm',
        sub: offer ? `${offer.co} — keep the conversation moving` : 'None yet'
      },
      {
        lbl: 'Applied & waiting',
        n: (counts.applied || 0) + (counts.screen || 0),
        tone: 'positive',
        sub: stale.length > 0 ? `${stale.length} worth a nudge this week` : 'All recent'
      },
      {
        lbl: 'Wishlist',
        n: counts.wishlist || 0,
        tone: 'mute',
        sub: oldestWish && wishDays !== null ? `Oldest sat ${wishDays} days` : 'Nothing saved yet'
      }
    ];
  });

  // AI-suggested actions — derived from the real pipeline.
  const actions = $derived.by(() => {
    const out = [];
    const interview = apps.find(a => a.status === 'interview');
    if (interview) {
      out.push({
        urgency: 'Up next', tone: 'accent', appId: interview.id,
        title: `Prep for your ${interview.co} loop`,
        sub: `${interview.role} — open the brief to see the interviewer briefing.`,
        cta: 'Open brief', logoSrc: interview.logoSrc
      });
    }
    const offer = apps.find(a => a.status === 'offer');
    if (offer) {
      out.push({
        urgency: 'Decide soon', tone: 'warm', appId: offer.id,
        title: `Decide on the ${offer.co} offer`,
        sub: `${offer.role} — use any active loops as leverage.`,
        cta: 'Open offer', logoSrc: offer.logoSrc
      });
    }
    const oldestStale = stale.sort((a, b) => (a.appliedDate || '').localeCompare(b.appliedDate || ''))[0];
    if (oldestStale) {
      out.push({
        urgency: `Quiet ${daysSince(oldestStale.raw.applied_at)} days`, tone: 'mute', appId: oldestStale.id,
        title: `Nudge ${oldestStale.co}`,
        sub: `${oldestStale.source !== '—' ? `Via ${oldestStale.source}` : 'Applied'} — no response yet.`,
        cta: 'Draft follow-up', logoSrc: oldestStale.logoSrc
      });
    }
    const screen = apps.find(a => a.status === 'screen');
    if (screen && screen.id !== interview?.id) {
      out.push({
        urgency: 'Recently moved', tone: 'accent', appId: screen.id,
        title: `Learn about ${screen.co}`,
        sub: `Moved to screen — generate the company brief.`,
        cta: 'Generate', logoSrc: screen.logoSrc
      });
    }
    return out.slice(0, 4);
  });

  // "What we're noticing" — rule-based for v0.1; AI weekly review lands in v0.3.
  const insights = $derived.by(() => {
    const out = [];
    const total = apps.filter(a => ['applied','screen','interview','offer','rejected','withdrawn'].includes(a.status)).length;
    if (total >= 3) {
      const referrals = apps.filter(a => /referral/i.test(a.source || ''));
      const refAdvanced = referrals.filter(a => ['screen','interview','offer'].includes(a.status)).length;
      const cold = apps.filter(a => !/referral/i.test(a.source || ''));
      const coldAdvanced = cold.filter(a => ['screen','interview','offer'].includes(a.status)).length;
      if (referrals.length >= 2 && cold.length >= 2) {
        const refRate = refAdvanced / referrals.length;
        const coldRate = coldAdvanced / cold.length || 0.0001;
        const mult = (refRate / coldRate).toFixed(1);
        if (refRate > coldRate) {
          out.push({
            icon: 'people', tone: 'positive',
            text: `Referrals convert at <b>${mult}×</b> the rate of cold apps.`,
            detail: `${refAdvanced} of ${referrals.length} referrals advanced vs. ${coldAdvanced} of ${cold.length} cold`
          });
        }
      }
    }
    if (apps.length > 0) {
      const newest = apps
        .filter(a => a.appliedDate)
        .sort((a, b) => (b.appliedDate || '').localeCompare(a.appliedDate || ''))[0];
      const d = newest ? daysSince(newest.appliedDate) : null;
      if (d !== null && d >= 4) {
        out.push({
          icon: 'pause', tone: 'warm',
          text: `You haven't applied in <b>${d} days</b>.`,
          detail: `Your last application was to ${newest.co}.`
        });
      }
    }
    if (stale.length >= 2) {
      const cos = stale.slice(0, 3).map(s => s.co).join(', ');
      out.push({
        icon: 'moon', tone: 'accent',
        text: `<b>${stale.length} ${stale.length === 1 ? 'loop has' : 'loops have'}</b> gone quiet for over a week.`,
        detail: `${cos} — worth a nudge.`
      });
    }
    return out;
  });

  const firstName = $derived.by(() => {
    if (!me?.email) return 'there';
    const local = me.email.split('@')[0].split(/[._]/).pop() || 'there';
    return local.charAt(0).toUpperCase() + local.slice(1);
  });
  const timeOfDayGreeting = $derived.by(() => {
    const h = new Date().getHours();
    if (h < 12) return 'Good morning';
    if (h < 18) return 'Good afternoon';
    return 'Good evening';
  });
  const today = new Date();
  const dateLong = today.toLocaleDateString('en-US', { weekday: 'long', day: 'numeric', month: 'long' });

  const filterMap = {
    active:    a => !['rejected', 'withdrawn'].includes(a.status),
    all:       () => true,
    wishlist:  a => a.status === 'wishlist',
    applied:   a => a.status === 'applied',
    screen:    a => a.status === 'screen',
    interview: a => a.status === 'interview',
    offer:     a => a.status === 'offer',
    closed:    a => ['rejected', 'withdrawn'].includes(a.status)
  };
  const visible = $derived(apps.filter(filterMap[filter] ?? filterMap.active));

  const chips = $derived([
    { k: 'active',    lbl: 'Active' },
    { k: 'interview', lbl: 'Interview' },
    { k: 'offer',     lbl: 'Offer' },
    { k: 'wishlist',  lbl: 'Wishlist' },
    { k: 'all',       lbl: 'All' }
  ]);
</script>

<svelte:head>
  <title>Today — Pursuit</title>
</svelte:head>

{#if showOnboarding}
  <Onboarding onDone={finishOnboarding} />
{/if}

<div class="topbar">
  <div class="crumb"><span class="here">Today</span></div>
  <div class="right">
    <div class="search">
      <svg class="ico" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5">
        <circle cx="7" cy="7" r="4.5"/><path d="M10.5 10.5L13 13"/>
      </svg>
      <span>Search applications, people…</span>
      <span class="kbd">⌘K</span>
    </div>
    <button class="btn">Import</button>
    <button class="btn btn-primary" onclick={() => (showNewModal = true)}>
      New application <span class="kbd">⌘N</span>
    </button>
  </div>
</div>

<div class="body">
  <div class="body-inner">

    <div class="hello">
      <div class="date">{dateLong}</div>
      <h1>{timeOfDayGreeting}, {firstName}.</h1>
    </div>

    {#if !loading && apps.length === 0}
      <div class="welcome-card">
        <div class="welcome-glyph">
          <svg viewBox="0 0 24 24" width="28" height="28" fill="none">
            <circle cx="12" cy="12" r="9.5" stroke="currentColor" stroke-width="1.4" opacity="0.65"/>
            <circle cx="12" cy="12" r="5.5" stroke="currentColor" stroke-width="1.4" opacity="0.9"/>
            <circle cx="17.5" cy="6.5" r="2.6" fill="currentColor"/>
          </svg>
        </div>
        <h2>Welcome to Pursuit.</h2>
        <p>
          Pursuit tracks your job search end-to-end — applications, interviews,
          and AI briefings on the people you're about to meet. Add your first
          application to start, or seed your account with 15 realistic
          applications so you can see every surface populated.
        </p>
        <div class="welcome-actions">
          <button class="btn btn-primary" onclick={() => (showNewModal = true)}>
            Add your first application
          </button>
          <button class="btn" onclick={seedDemo} disabled={seeding}>
            {seeding ? 'Seeding…' : 'Try with demo data'}
          </button>
        </div>
        {#if seedError}<p class="welcome-err">{seedError}</p>{/if}
        <p class="welcome-note">Demo rows are tagged so you can clear them in one click later.</p>
      </div>
    {:else}
    <div class="counts">
      {#each countCards as c}
        <div class={`count-cell tone-${c.tone}`}>
          <span class="ribbon"></span>
          <div class="cell-top">
            <span class="lbl">
              {c.lbl}
              {#if c.hint}<span class="hint" title={c.hint}>i</span>{/if}
            </span>
          </div>
          <div class="n">{c.n}</div>
          <div class="sub">{c.sub}</div>
        </div>
      {/each}
    </div>

    {#if actions.length > 0}
      <div class="section-hd">
        <h2>What you can do today</h2>
        <span class="ai-tag">
          <svg width="11" height="11" viewBox="0 0 12 12" fill="currentColor"><path d="M6 0l1.2 3.6L11 5l-3.8 1.4L6 10 4.8 6.4 1 5l3.8-1.4z"/></svg>
          <span>AI suggested</span>
        </span>
      </div>
      <div class="action-grid">
        {#each actions as a}
          <div class={`action-card ${a.tone}`} onclick={() => open(a.appId)} role="button" tabindex="0">
            <div class="action-top">
              {#if a.logoSrc}<img class="action-logo" src={a.logoSrc} alt="" />{/if}
              <span class={`urgency u-${a.tone}`}>{a.urgency}</span>
            </div>
            <h3>{a.title}</h3>
            <p>{a.sub}</p>
            <button class="action-cta">{a.cta} <span class="arrow">→</span></button>
          </div>
        {/each}
      </div>
    {/if}

    {#if insights.length > 0}
      <div class="section-hd">
        <h2>What we're noticing</h2>
        <span class="ai-tag">
          <svg width="11" height="11" viewBox="0 0 12 12" fill="currentColor"><path d="M6 0l1.2 3.6L11 5l-3.8 1.4L6 10 4.8 6.4 1 5l3.8-1.4z"/></svg>
          <span>This week</span>
        </span>
      </div>
      <div class="insight-list">
        {#each insights as ins}
          <div class="insight">
            <span class={`ins-icon t-${ins.tone}`}>
              {#if ins.icon === 'people'}
                <svg viewBox="0 0 24 24" width="22" height="22" fill="currentColor">
                  <circle cx="8.5" cy="8" r="3.5"/>
                  <circle cx="16.5" cy="9" r="2.8"/>
                  <path d="M2 19c0-3.3 2.9-5.5 6.5-5.5s6.5 2.2 6.5 5.5v1H2zM15 19.5c0-2.4 2-4 4.5-4S22 17 22 19.5V20h-7z"/>
                </svg>
              {:else if ins.icon === 'pause'}
                <svg viewBox="0 0 24 24" width="22" height="22" fill="currentColor">
                  <rect x="6" y="5" width="4" height="14" rx="1.5"/>
                  <rect x="14" y="5" width="4" height="14" rx="1.5"/>
                </svg>
              {:else if ins.icon === 'moon'}
                <svg viewBox="0 0 24 24" width="22" height="22" fill="currentColor">
                  <path d="M14 3a9 9 0 1 0 7 14 7 7 0 0 1-7-14z"/>
                </svg>
              {/if}
            </span>
            <div class="ins-body">
              <div class="ins-line">{@html ins.text}</div>
              <div class="ins-detail">{ins.detail}</div>
            </div>
          </div>
        {/each}
      </div>
    {/if}

    <div class="section-hd">
      <h2>Applications <span class="count-pill">{apps.length}</span></h2>
      <div class="filters">
        {#if hasDemoRows}
          <button class="chip chip-warn" onclick={clearDemo} title="Remove demo data — real applications are kept">
            Clear demo data
          </button>
        {/if}
        {#each chips as c}
          <button class={`chip ${filter === c.k ? 'active' : ''}`} onclick={() => (filter = c.k)}>{c.lbl}</button>
        {/each}
      </div>
    </div>
    <div class="table">
      <div class="tr head">
        <span>Company</span><span>Role</span><span>Status</span><span>Applied</span><span></span>
      </div>
      {#if loading}
        <div class="tr"><span style="color:var(--mute); grid-column: 1 / -1; text-align:center; padding:1rem 0;">Loading…</span></div>
      {:else if visible.length === 0}
        <div class="tr"><span style="color:var(--mute); grid-column: 1 / -1; text-align:center; padding:2rem 0;">
          {filter === 'all' && apps.length === 0
            ? 'No applications yet. Press ⌘N to add your first one.'
            : `No applications in "${filter}."`}
        </span></div>
      {:else}
        {#each visible as a (a.id)}
          <div class={`tr ${a.stale ? 'stale' : ''}`} onclick={() => open(a.id)}>
            <span class="co">
              {#if a.logoSrc}
                <img class="logo" src={a.logoSrc} alt="" />
              {:else}
                <span class={`logo letter ${a.logoCls}`}>{a.coShort}</span>
              {/if}
              <span>{a.co}</span>
            </span>
            <span class="role">{a.role}</span>
            <span class={`pill ${a.status}`}><span class="pdot"></span>{STATUS_LABEL[a.status]}</span>
            <span class="applied">
              {a.appliedRel}
              {#if a.stale}<span class="stale-tag">stale</span>{/if}
            </span>
            <span class="arrow">→</span>
          </div>
        {/each}
      {/if}
    </div>
    {/if}
  </div>
</div>

<svelte:window
  onkeydown={(e) => { if (showNewModal) onModalKeydown(e); }}
  onpaste={(e) => { if (showNewModal) onModalPaste(e); }}
/>

{#if showNewModal}
  <div class="modal-overlay" onclick={resetModal} role="presentation">
    <form class="modal" onclick={(e) => e.stopPropagation()} onsubmit={createApp}>
      <header class="m-head">
        <div>
          <h2>New application</h2>
          <p class="sub">Drop a JD, paste a URL, or describe the role — we'll fill the fields. Or do it by hand below.</p>
        </div>
        <button type="button" class="x-close" onclick={resetModal} aria-label="Close">
          <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="1.6"><path d="M3 3l8 8M11 3l-8 8" stroke-linecap="round"/></svg>
        </button>
      </header>

      <div class="quick-add">
        <div class="qa-head">
          <span class="qa-label">Quick add</span>
          <span class="ai-tag-sm">AI</span>
        </div>

        <div
          class={`drop ${isDraggingFile ? 'drag' : ''} ${attachedImage ? 'has-image' : ''} ${parsing ? 'loading' : ''}`}
          ondragover={onDropZoneDragOver}
          ondragleave={onDropZoneDragLeave}
          ondrop={onDropZoneDrop}
        >
          {#if attachedImage}
            <div class="attached">
              <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="2" y="3" width="12" height="10" rx="1.5"/><circle cx="5.5" cy="6.5" r="1"/><path d="M2 11l3.5-3.5 3 3 2-2L14 11"/></svg>
              <span class="att-name">{attachedImage.name}</span>
              <span class="att-size">{Math.round(attachedImage.size / 1024)} KB</span>
              <button type="button" class="att-x" onclick={() => (attachedImage = null)} aria-label="Remove" disabled={parsing}>×</button>
            </div>
          {:else}
            <textarea
              bind:value={pasteText}
              onpaste={onTextareaPaste}
              placeholder="Drop a screenshot here, paste a job URL, or describe the role you're applying to…"
              rows="3"
              disabled={parsing}
            ></textarea>
          {/if}

          <div class="drop-footer">
            <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="2" y="3" width="12" height="10" rx="1.5"/><circle cx="5.5" cy="6.5" r="1"/><path d="M2 11l3.5-3.5 3 3 2-2L14 11"/></svg>
            <span>Drag &amp; drop a screenshot, or paste anything.</span>
          </div>

          {#if parsing}
            <div class="drop-loading" role="status" aria-live="polite">
              <span class="spinner" aria-hidden="true"></span>
              <span class="loading-text">{parsingStage}</span>
            </div>
          {/if}
        </div>

        <div class="qa-actions">
          <span class="kb-hints">
            <kbd>⌘V</kbd> here · <kbd>⌘⇧4</kbd> to screenshot
          </span>
          <button type="button" class="btn btn-accent" onclick={parseJD} disabled={parsing || (!pasteText.trim() && !attachedImage)}>
            {#if parsing}
              <span class="spinner spinner-sm" aria-hidden="true"></span>
            {:else}
              <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M8 2l1.4 3.6L13 7l-3.6 1.4L8 12l-1.4-3.6L3 7l3.6-1.4L8 2z"/></svg>
            {/if}
            {parsing ? 'Parsing…' : 'Parse with AI'}
          </button>
        </div>

        {#if parsedHint}<div class="qa-msg ok"><span class="ok-check">✓</span> {parsedHint}</div>{/if}
        {#if parseError}<div class="qa-msg err">{parseError}</div>{/if}
      </div>

      <div class="modal-divider"><span>Or enter by hand</span></div>

      <div class="fields">
        <label>
          <span class="lbl">Company <em class="req">*</em></span>
          <input bind:value={form.company} placeholder="Anthropic" required />
        </label>
        <label>
          <span class="lbl">Role <em class="req">*</em></span>
          <input bind:value={form.role} placeholder="Senior Software Engineer" required />
        </label>
        <label>
          <span class="lbl">Status</span>
          <div class={`status-select ${form.status}`}>
            <span class="sdot"></span>
            <select bind:value={form.status}>
              <option value="wishlist">Wishlist</option>
              <option value="applied">Applied</option>
              <option value="screen">Screen</option>
              <option value="interview">Interview</option>
              <option value="offer">Offer</option>
              <option value="rejected">Rejected</option>
              <option value="withdrawn">Withdrawn</option>
            </select>
            <svg class="chev" width="10" height="10" viewBox="0 0 10 10" fill="none" stroke="currentColor" stroke-width="1.6"><path d="M2 4l3 3 3-3" stroke-linecap="round" stroke-linejoin="round"/></svg>
          </div>
        </label>
        <label>
          <span class="lbl">Source <span class="opt">— optional</span></span>
          <input bind:value={form.source} placeholder="LinkedIn / Referral / Cold email" />
        </label>
      </div>

      <footer class="m-foot">
        <span class="esc-hint"><kbd>Esc</kbd> to cancel</span>
        <div class="foot-actions">
          <button type="button" class="btn" onclick={resetModal}>Cancel</button>
          <button type="submit" class="btn btn-primary" disabled={saving || !form.company || !form.role}>
            {saving ? 'Saving…' : 'Add application'} <kbd class="dark-kbd">↵</kbd>
          </button>
        </div>
      </footer>
    </form>
  </div>
{/if}

<style>
  /* ── Greeting ─────────────────────────────────────────── */
  .hello { margin-bottom: 28px; }
  .hello .date { font-size: 13.5px; color: var(--mute); margin-bottom: 6px; font-weight: 400; }
  .hello h1 {
    font-size: 30px; font-weight: 600;
    margin: 0; letter-spacing: -0.025em;
    color: var(--ink);
  }

  /* ── Count cards ──────────────────────────────────────── */
  .counts {
    display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px;
    margin-bottom: 40px;
  }
  .count-cell {
    position: relative;
    background: var(--card); border: 1px solid var(--rule); border-radius: 14px;
    padding: 18px 20px 18px;
    box-shadow: var(--sh-1);
    overflow: hidden;
    transition: transform 140ms ease, box-shadow 140ms ease, border-color 140ms ease;
  }
  .count-cell:hover { transform: translateY(-2px); box-shadow: var(--sh-pop); border-color: var(--rule-strong); }
  .ribbon { position: absolute; top: 0; left: 0; right: 0; height: 3px; }
  .count-cell.tone-accent   .ribbon { background: var(--accent); }
  .count-cell.tone-warm     .ribbon { background: var(--warm); }
  .count-cell.tone-positive .ribbon { background: var(--positive); }
  .count-cell.tone-mute     .ribbon { background: var(--rule-strong); }
  .cell-top { display: flex; justify-content: space-between; align-items: center; margin-bottom: 4px; }
  .count-cell .lbl { font-size: 13px; color: var(--mute); display: flex; align-items: center; gap: 6px; font-weight: 500; }
  .count-cell .n {
    font-size: 42px; font-weight: 600; letter-spacing: -0.035em;
    line-height: 1.05; color: var(--ink); font-feature-settings: "tnum";
    margin-top: 2px;
  }
  .count-cell.tone-accent   .n { color: var(--accent-text); }
  .count-cell.tone-warm     .n { color: var(--warm-text); }
  .count-cell.tone-positive .n { color: var(--positive-text); }
  .count-cell .sub {
    margin-top: 8px; font-size: 12.5px; color: var(--mute); line-height: 1.4;
    padding-top: 8px; border-top: 1px dashed var(--rule);
  }
  .hint { width: 14px; height: 14px; display: inline-grid; place-items: center; border-radius: 50%; background: var(--surface-2); color: var(--mute); font-size: 10px; font-style: italic; font-weight: 600; cursor: help; }

  /* ── Section heading ──────────────────────────────────── */
  .section-hd { display: flex; align-items: center; justify-content: space-between; margin: 32px 0 14px; }
  .section-hd h2 {
    font-size: 18px; font-weight: 600; margin: 0; letter-spacing: -0.02em;
  }
  .section-hd h2 .count-pill { font-size: 14px; color: var(--mute); margin-left: 8px; font-weight: 400; }
  .ai-tag {
    display: inline-flex; align-items: center; gap: 5px;
    font-size: 12.5px; color: var(--accent-text);
    background: var(--accent-tint); padding: 3px 10px; border-radius: 99px;
    font-weight: 500;
  }

  /* ── Action grid ──────────────────────────────────────── */
  .action-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(260px, 1fr)); gap: 12px; }
  .action-card { background: var(--card); border: 1px solid var(--rule); border-radius: 14px; padding: 18px; display: flex; flex-direction: column; gap: 8px; transition: transform 140ms ease, border-color 140ms ease, box-shadow 140ms ease; cursor: pointer; }
  .action-card:hover { transform: translateY(-2px); border-color: var(--accent); box-shadow: var(--sh-pop); }
  .action-top { display: flex; align-items: center; justify-content: space-between; }
  .action-logo { width: 30px; height: 30px; border-radius: 7px; background: var(--surface-2); object-fit: contain; padding: 3px; }
  .urgency { font-size: 12px; padding: 3px 9px; border-radius: 99px; font-weight: 500; }
  .urgency.u-accent { background: var(--accent-tint); color: var(--accent-text); }
  .urgency.u-warm { background: var(--warm-tint); color: var(--warm-text); }
  .urgency.u-mute { background: var(--surface-2); color: var(--mute); }
  .action-card h3 { font-size: 15.5px; font-weight: 500; margin: 4px 0 0; letter-spacing: -0.01em; }
  .action-card p { font-size: 13px; color: var(--mute); margin: 0; line-height: 1.5; }
  .action-cta { background: transparent; border: 0; color: var(--accent-text); font-size: 13px; font-weight: 500; text-align: left; padding: 4px 0 0; cursor: pointer; align-self: flex-start; }
  .action-cta .arrow { transition: transform 140ms ease; display: inline-block; }
  .action-card:hover .action-cta .arrow { transform: translateX(2px); }

  /* ── Insights ─────────────────────────────────────────── */
  .insight-list { display: flex; flex-direction: column; gap: 1px; background: var(--rule); border: 1px solid var(--rule); border-radius: 14px; overflow: hidden; }
  .insight { background: var(--card); padding: 14px 20px; display: grid; grid-template-columns: 44px 1fr; gap: 16px; align-items: center; }
  .ins-icon {
    width: 44px; height: 44px; border-radius: 12px;
    display: grid; place-items: center;
  }
  .ins-icon.t-positive { background: var(--positive-tint); color: var(--positive-text); }
  .ins-icon.t-warm     { background: var(--warm-tint);     color: var(--warm-text); }
  .ins-icon.t-accent   { background: var(--accent-tint);   color: var(--accent-text); }
  .ins-line { font-size: 14px; color: var(--ink); line-height: 1.45; }
  .ins-detail { font-size: 12.5px; color: var(--mute); margin-top: 3px; }

  /* ── Applications table ───────────────────────────────── */
  .filters { display: flex; gap: 4px; }
  .chip { background: transparent; border: 1px solid var(--rule); border-radius: 99px; padding: 4px 12px; font-size: 12.5px; color: var(--mute); cursor: pointer; }
  .chip.active { background: var(--ink); border-color: var(--ink); color: white; }
  .chip.chip-warn { color: var(--danger-text); border-color: var(--danger-tint); }
  .chip.chip-warn:hover { background: var(--danger-tint); }

  /* Welcome card — shown only when the signed-in user has zero applications. */
  .welcome-card {
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 18px;
    padding: 36px 40px;
    text-align: center;
    box-shadow: var(--sh-1);
    max-width: 640px;
    margin: 16px auto 0;
  }
  .welcome-glyph {
    width: 56px; height: 56px;
    margin: 0 auto 14px;
    border-radius: 14px;
    background: var(--accent-tint);
    color: var(--accent);
    display: grid; place-items: center;
  }
  .welcome-card h2 { font-size: 22px; font-weight: 600; margin: 0 0 .5rem; letter-spacing: -0.02em; }
  .welcome-card > p { font-size: 14px; color: var(--mute); line-height: 1.55; margin: 0 auto 1.5rem; max-width: 52ch; }
  .welcome-actions { display: flex; justify-content: center; gap: 10px; flex-wrap: wrap; }
  .welcome-err { color: var(--danger-text); font-size: 12.5px; margin: 12px 0 0; }
  .welcome-note { font-size: 12px; color: var(--mute-2); margin: 12px 0 0; }
  .table { background: var(--card); border: 1px solid var(--rule); border-radius: 14px; overflow: hidden; }
  .tr { display: grid; grid-template-columns: 220px 1fr 130px 160px 20px; align-items: center; padding: 12px 20px; border-bottom: 1px solid var(--rule); font-size: 13.5px; cursor: pointer; }
  .tr:last-child { border-bottom: 0; }
  .tr:hover { background: var(--surface-2); }
  .tr.head { background: var(--surface-2); font-size: 12px; color: var(--mute); cursor: default; padding: 10px 20px; font-weight: 500; }
  .tr.head:hover { background: var(--surface-2); }
  .tr.stale .applied { color: var(--danger-text); }
  .co { display: flex; align-items: center; gap: 10px; font-weight: 500; }
  .co .logo { width: 22px; height: 22px; border-radius: 5px; background: var(--surface-2); object-fit: contain; padding: 2px; }
  .co .logo.letter {
    display: grid; place-items: center;
    padding: 0;
    color: var(--ink-2); font-size: 11px; font-weight: 600;
  }
  .role { color: var(--mute); }
  .pill { display: inline-flex; align-items: center; gap: 5px; padding: 3px 9px; border-radius: 99px; font-size: 12px; font-weight: 500; background: var(--surface-2); color: var(--ink-2); width: max-content; }
  .pill .pdot { width: 5px; height: 5px; border-radius: 50%; background: var(--mute-2); }
  .pill.interview { background: var(--accent-tint); color: var(--accent-text); }
  .pill.interview .pdot { background: var(--accent); }
  .pill.offer { background: var(--warm-tint); color: var(--warm-text); }
  .pill.offer .pdot { background: var(--warm); }
  .pill.screen { background: var(--positive-tint); color: var(--positive-text); }
  .pill.screen .pdot { background: var(--positive); }
  .applied { color: var(--mute); }
  .stale-tag { font-size: 11px; background: var(--danger-tint); color: var(--danger-text); padding: 1px 7px; border-radius: 99px; margin-left: 4px; font-weight: 500; }
  .arrow { color: var(--mute-2); }

  /* ── Modal styles (unchanged) ─────────────────────────── */
  .modal-overlay {
    position: fixed; inset: 0;
    background: rgba(10,10,13,0.45);
    backdrop-filter: blur(2px);
    display: grid; place-items: center;
    z-index: 100;
    padding: 2rem;
  }
  .modal {
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 14px;
    width: 100%;
    max-width: 580px;
    display: flex; flex-direction: column;
    box-shadow: var(--sh-pop);
    overflow: hidden;
  }
  .m-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 16px; padding: 20px 22px 14px; }
  .m-head h2 { font-size: 19px; font-weight: 500; letter-spacing: -0.018em; margin: 0; color: var(--ink); }
  .m-head .sub { margin: 4px 0 0; font-size: 13px; color: var(--mute); line-height: 1.5; max-width: 460px; }
  .x-close {
    background: transparent; border: 0; color: var(--mute);
    width: 28px; height: 28px; border-radius: 6px;
    display: grid; place-items: center; cursor: pointer;
    flex-shrink: 0;
    transition: background 100ms ease, color 100ms ease;
  }
  .x-close:hover { background: var(--surface-2); color: var(--ink); }
  .quick-add { padding: 0 22px 16px; display: flex; flex-direction: column; gap: 10px; }
  .qa-head { display: flex; align-items: center; justify-content: space-between; }
  .qa-label { font-size: 12px; font-weight: 500; color: var(--ink-2); letter-spacing: 0; }
  .ai-tag-sm { font-weight: 500; font-size: 10px; color: var(--accent-text); background: var(--accent-tint); border-radius: 4px; padding: 2px 7px; letter-spacing: .04em; }

  .drop { position: relative; background: var(--surface); border: 1.5px dashed var(--rule-strong); border-radius: 10px; transition: border-color 120ms ease, background 120ms ease; overflow: hidden; }
  .drop.drag { border-color: var(--accent); background: var(--accent-tint); }
  .drop.loading { border-color: var(--accent); border-style: solid; }
  .drop.loading::before {
    content: ''; position: absolute; inset: 0 0 auto 0; height: 2px;
    background: linear-gradient(90deg, transparent, var(--accent), transparent);
    background-size: 40% 100%; background-repeat: no-repeat;
    animation: drop-sweep 1.4s linear infinite;
  }
  @keyframes drop-sweep {
    0%   { background-position: -40% 0; }
    100% { background-position: 140% 0; }
  }
  .drop-loading { display: flex; align-items: center; gap: 10px; padding: 10px 14px; background: var(--accent-tint); border-top: 1px solid var(--accent-tint-2); font-size: 12.5px; color: var(--accent-text); font-weight: 500; }
  .loading-text { letter-spacing: 0; }
  .spinner { display: inline-block; width: 14px; height: 14px; border: 1.8px solid var(--accent-tint-2); border-top-color: var(--accent); border-radius: 50%; animation: spinner-rot 0.7s linear infinite; flex-shrink: 0; }
  .spinner-sm { width: 12px; height: 12px; border-color: rgba(255,255,255,.35); border-top-color: white; }
  @keyframes spinner-rot { to { transform: rotate(360deg); } }
  .ok-check { display: inline-block; width: 14px; height: 14px; line-height: 14px; text-align: center; border-radius: 50%; background: var(--positive-tint); color: var(--positive-text); font-size: 10px; font-weight: 600; margin-right: 4px; }
  .drop textarea { width: 100%; font: inherit; font-family: var(--sans); font-size: 13.5px; color: var(--ink); background: transparent; border: 0; padding: 12px 14px 8px; outline: none; resize: none; min-height: 78px; line-height: 1.5; display: block; }
  .drop textarea::placeholder { color: var(--mute-2); }
  .attached { display: flex; align-items: center; gap: 10px; padding: 12px 14px 8px; font-size: 13px; color: var(--ink-2); }
  .attached svg { color: var(--accent-text); flex-shrink: 0; }
  .attached .att-name { font-weight: 500; color: var(--ink); }
  .attached .att-size { color: var(--mute); font-family: var(--mono); font-size: 11px; }
  .attached .att-x { margin-left: auto; background: transparent; border: 0; color: var(--mute); font-size: 18px; line-height: 1; cursor: pointer; padding: 0 4px; }
  .attached .att-x:hover { color: var(--ink); }
  .drop-footer { display: flex; align-items: center; gap: 8px; padding: 8px 14px 10px; font-size: 11.5px; color: var(--mute-2); border-top: 1px dashed var(--rule); background: transparent; }
  .drop-footer svg { color: var(--mute-2); flex-shrink: 0; }
  .qa-actions { display: flex; align-items: center; justify-content: space-between; gap: 12px; }
  .kb-hints { font-size: 11.5px; color: var(--mute); }
  .kb-hints kbd { font-family: var(--mono); font-size: 10.5px; background: var(--surface); border: 1px solid var(--rule); border-bottom-width: 2px; border-radius: 3px; padding: 0 4px; color: var(--ink-2); }
  .btn-accent { display: inline-flex; align-items: center; gap: 6px; background: var(--accent); color: white; border: 1px solid var(--accent); border-radius: 7px; padding: 7px 12px; font-size: 13px; font-weight: 500; cursor: pointer; transition: background 100ms ease, border-color 100ms ease; }
  .btn-accent:hover:not(:disabled) { background: var(--accent-strong); border-color: var(--accent-strong); }
  .btn-accent:disabled { opacity: .55; cursor: not-allowed; }
  .btn-accent svg { color: white; }
  .qa-msg { font-size: 12px; }
  .qa-msg.ok  { color: var(--positive-text); }
  .qa-msg.err { color: var(--danger-text); }
  .modal-divider { display: flex; align-items: center; gap: 12px; color: var(--mute-2); font-size: 10.5px; letter-spacing: .06em; text-transform: uppercase; padding: 0 22px; margin: 4px 0 14px; }
  .modal-divider::before, .modal-divider::after { content: ''; flex: 1; height: 1px; background: var(--rule); }
  .fields { padding: 0 22px 18px; display: grid; grid-template-columns: 1fr 1fr; gap: 12px 14px; }
  .fields label { display: flex; flex-direction: column; gap: 5px; }
  .fields .lbl { font-size: 11.5px; color: var(--mute); font-weight: 500; letter-spacing: 0; }
  .fields .req { color: var(--accent-text); font-style: normal; margin-left: 1px; }
  .fields .opt { color: var(--mute-2); font-weight: 400; }
  .fields input { font: inherit; font-size: 13.5px; color: var(--ink); background: var(--surface); border: 1px solid var(--rule); border-radius: 8px; padding: 0 11px; height: 36px; outline: 0; transition: border-color 100ms ease, box-shadow 100ms ease, background 100ms ease; }
  .fields input:hover { border-color: var(--rule-strong); }
  .fields input:focus { border-color: var(--accent); box-shadow: 0 0 0 3px var(--accent-tint); }
  .fields input::placeholder { color: var(--mute-2); }
  .status-select { position: relative; display: flex; align-items: center; background: var(--surface); border: 1px solid var(--rule); border-radius: 8px; height: 36px; padding: 0 11px; transition: border-color 100ms ease, box-shadow 100ms ease; }
  .status-select:hover { border-color: var(--rule-strong); }
  .status-select:focus-within { border-color: var(--accent); box-shadow: 0 0 0 3px var(--accent-tint); }
  .status-select .sdot { width: 7px; height: 7px; border-radius: 50%; background: var(--mute); margin-right: 8px; flex-shrink: 0; }
  .status-select.wishlist  .sdot { background: var(--mute); }
  .status-select.applied   .sdot { background: var(--ink-2); }
  .status-select.screen    .sdot { background: var(--positive); }
  .status-select.interview .sdot { background: var(--accent); }
  .status-select.offer     .sdot { background: var(--warm); }
  .status-select.rejected  .sdot { background: var(--mute-2); }
  .status-select.withdrawn .sdot { background: var(--mute-2); }
  .status-select select { flex: 1; appearance: none; -webkit-appearance: none; background: transparent; border: 0; outline: 0; font: inherit; font-size: 13.5px; color: var(--ink); padding: 0 18px 0 0; cursor: pointer; }
  .status-select .chev { position: absolute; right: 11px; top: 50%; transform: translateY(-50%); color: var(--mute); pointer-events: none; }
  .m-foot { display: flex; align-items: center; justify-content: space-between; padding: 12px 22px; background: var(--surface); border-top: 1px solid var(--rule); }
  .esc-hint { font-size: 11.5px; color: var(--mute); }
  .esc-hint kbd { font-family: var(--mono); font-size: 10.5px; background: var(--card); border: 1px solid var(--rule); border-bottom-width: 2px; border-radius: 3px; padding: 0 5px; color: var(--ink-2); }
  .foot-actions { display: flex; gap: 8px; }
  .foot-actions .btn-primary { display: inline-flex; align-items: center; gap: 6px; }
  .dark-kbd { font-family: var(--mono); font-size: 10.5px; background: rgba(255,255,255,.18); border: 1px solid rgba(255,255,255,.22); border-radius: 3px; padding: 0 5px; color: rgba(255,255,255,.9); }
</style>
