<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/state';
  import { api } from '$lib/api.js';
  import {
    STATUS_LABEL, toDisplayApp, fmtWeekday, countsByStatus
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
  let attachedImage = $state(null); // { name, size, mediaType, file }
  let isDraggingFile = $state(false);
  let parsingStage = $state(''); // user-facing status text while parsing

  const ALLOWED_IMG = ['image/png', 'image/jpeg', 'image/gif', 'image/webp'];
  const MAX_IMG_BYTES = 6 * 1024 * 1024; // 6 MB raw — Anthropic caps at 5MB and we leave slack

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
        // result is "data:image/png;base64,XXXXX..." — strip the prefix.
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
    // Cycle through honest progress messages so the user feels something is
    // happening. Vision calls take 3-6s; pure-text calls 2-4s.
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

  function onModalKeydown(e) {
    if (e.key === 'Escape') resetModal();
  }
  function onModalPaste(e) {
    // Catch image paste anywhere in the modal, not just inside the textarea.
    // Text paste falls through to whatever input has focus.
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

  // Sidebar Pipeline links pass ?filter=X to jump straight to a filtered table.
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

  function finishOnboarding() {
    showOnboarding = false;
    // Re-fetch in case applications were created during onboarding.
    refresh();
  }

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

  function open(id) {
    goto(`/app/${id}`);
  }

  const counts  = $derived(countsByStatus(apps));
  const active  = $derived(apps.filter(a => !['rejected', 'withdrawn'].includes(a.status)));
  const upcoming = $derived(apps.filter(a => a.status === 'interview' || a.status === 'screen'));

  // Today section: prefer the interview-status app with the most recent applied date.
  const todayItem = $derived(apps.find(a => a.status === 'interview') ?? null);
  // Other near-term items: screens & offers, excluding the one already in Today.
  const weekItems = $derived(apps
    .filter(a => (a.status === 'screen' || a.status === 'offer' || a.status === 'interview') && a !== todayItem)
    .slice(0, 5));

  const today = new Date();
  const todayString = today.toLocaleDateString('en-US', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' });

  const stats = $derived([
    { lbl: 'Active',   n: active.length },
    { lbl: 'Loops',    n: counts.interview || 0 },
    { lbl: 'Offers',   n: counts.offer || 0 },
    { lbl: 'Wishlist', n: counts.wishlist || 0 }
  ]);

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
    { k: 'active',    lbl: 'Active',    n: active.length },
    { k: 'interview', lbl: 'Interview', n: counts.interview || 0 },
    { k: 'offer',     lbl: 'Offer',     n: counts.offer || 0 },
    { k: 'wishlist',  lbl: 'Wishlist',  n: counts.wishlist || 0 },
    { k: 'closed',    lbl: 'Closed',    n: (counts.rejected || 0) + (counts.withdrawn || 0) },
    { k: 'all',       lbl: 'All',       n: apps.length }
  ]);

  // ── Rich page header pieces ────────────────────────────────
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
  const dow  = $derived(today.toLocaleDateString('en-US', { weekday: 'long' }));
  const dnum = $derived(today.toLocaleDateString('en-US', { day: 'numeric', month: 'long', year: 'numeric' }));

  // ── Briefer copy derived from real pipeline state (no AI yet) ─
  const briefer = $derived.by(() => {
    if (apps.length === 0) return null;
    const offer     = apps.find(a => a.status === 'offer');
    const interview = apps.find(a => a.status === 'interview');
    const screen    = apps.find(a => a.status === 'screen');
    if (offer && interview) {
      return {
        msg: `Two threads worth your attention — the <em>${offer.co} offer</em> needs a reply and the <em>${interview.co} loop</em> is live.`,
        primary: { label: `Open ${offer.co}`, id: offer.id }
      };
    }
    if (offer) {
      return {
        msg: `One open offer from <em>${offer.co}</em> — use it as leverage on anything else live.`,
        primary: { label: `Open ${offer.co}`, id: offer.id }
      };
    }
    if (interview) {
      return {
        msg: `<em>${interview.co}</em> is your most active loop. Refresh the dossier before you talk.`,
        primary: { label: `Open ${interview.co}`, id: interview.id }
      };
    }
    if (screen) {
      return {
        msg: `<em>${screen.co}</em> just moved to <em>Screen</em>. Worth thinking about what you want from the conversation.`,
        primary: { label: `Open ${screen.co}`, id: screen.id }
      };
    }
    // No live loop — fall back to the most recent active application so
    // there's still a real company name in the message and a real CTA.
    const applied = apps.filter(a => a.status === 'applied');
    if (applied.length > 0) {
      const recent = applied[0]; // /api/applications is sorted applied_at DESC
      return {
        msg: `<em>${recent.co}</em> is your most recent application. Nothing has moved yet — worth a gentle follow-up if it has been a few days.`,
        primary: { label: `Open ${recent.co}`, id: recent.id }
      };
    }
    const wishlist = apps.find(a => a.status === 'wishlist');
    if (wishlist) {
      return {
        msg: `<em>${wishlist.co}</em> is on your wishlist. Move it to <em>Applied</em> when you actually send.`,
        primary: { label: `Open ${wishlist.co}`, id: wishlist.id }
      };
    }
    // True last resort: everything closed.
    return {
      msg: `${apps.length} ${apps.length === 1 ? 'application' : 'applications'} on file, all closed. Add a new one with <b>⌘N</b>.`,
      primary: null
    };
  });
  let briefDismissed = $state(false);

  // ── Sub-greeting line ─────────────────────────────────────
  const subGreeting = $derived.by(() => {
    if (apps.length === 0) {
      return `Add your first applications with <b>⌘N</b> to start seeing the patterns.`;
    }
    const parts = [];
    if (counts.interview > 0) parts.push(`<b>${counts.interview}</b> interview ${counts.interview === 1 ? 'loop' : 'loops'}`);
    if (counts.offer > 0)     parts.push(`<b>${counts.offer}</b> open ${counts.offer === 1 ? 'offer' : 'offers'}`);
    if (counts.screen > 0)    parts.push(`<b>${counts.screen}</b> in ${counts.screen === 1 ? 'screen' : 'screens'}`);
    if (parts.length === 0) {
      const n = active.length;
      return `<b>${n}</b> ${n === 1 ? 'application' : 'applications'} in flight — keep going.`;
    }
    return parts.join(' &middot; ');
  });
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
    <div class="page-hd page-hd-rich">
      <div>
        <div class="date">
          <span class="dow">{dow}</span>
          <span class="sep">/</span>
          <span class="dnum">{dnum}</span>
        </div>
        <h1>{timeOfDayGreeting}, {firstName}.</h1>
        <div class="sub-greeting">{@html subGreeting}</div>
      </div>
      <div class="stats stats-rich">
        {#each stats as s}
          <div class="stat-cell">
            <div class="lbl">{s.lbl}</div>
            <div class="n-row">
              <span class="n">{s.n}</span>
            </div>
          </div>
        {/each}
      </div>
    </div>

    {#if briefer && !briefDismissed}
      <div class="briefer">
        <span class="glyph">P</span>
        <div>
          <div class="who">
            <span class="label">Pursuit · briefed</span>
            <span class="tag">derived from your pipeline</span>
          </div>
          <p class="msg">{@html briefer.msg}</p>
          {#if briefer.primary}
            <div class="chips">
              <a class="chip-c primary" href="/app/{briefer.primary.id}">
                {briefer.primary.label} <span class="arrow">→</span>
              </a>
            </div>
          {/if}
        </div>
        <button class="dismiss" onclick={() => (briefDismissed = true)} title="Dismiss for now">
          <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M3 3l6 6M9 3l-6 6"/></svg>
        </button>
      </div>
    {/if}

    <!-- Today section -->
    <div class="section-hd">
      <h2>Today <span class="count">{todayItem ? 1 : 0}</span></h2>
      <div class="actions">
        <button class="btn" style="height:26px; font-size:12.5px">Sync calendar</button>
      </div>
    </div>
    {#if todayItem}
      <div class="today">
        <div class="today-item urgent" onclick={() => open(todayItem.id)}>
          <div class="when">
            <span class="time">Soon</span>
            <span class="tag">{STATUS_LABEL[todayItem.status]}</span>
          </div>
          <div class="summary">
            <div class="title">{todayItem.co} · {todayItem.role}</div>
            <div class="meta">
              Open the dossier to see the interviewer briefing
              <span class="dot">·</span>
              Dossier refreshed 12 min ago
            </div>
          </div>
          <button class="cta">
            Open dossier <span class="kbd">↵</span>
          </button>
        </div>
      </div>
    {:else}
      <div class="empty-tab">
        <h3>No interviews scheduled today</h3>
        <p>When you move an application to <b>Interview</b>, it appears here with the AI dossier ready to open.</p>
      </div>
    {/if}

    <!-- This week -->
    {#if weekItems.length > 0}
      <div class="section-hd">
        <h2>This week <span class="count">{weekItems.length}</span></h2>
      </div>
      <div class="weeklist">
        {#each weekItems as a (a.id)}
          <div class="week-item" onclick={() => open(a.id)}>
            <div class="when"><b>{fmtWeekday(a.appliedDate)}</b> {a.applied}</div>
            <div class="what">{a.co} · <span class="role">{a.role}</span></div>
            <span class={`pill ${a.status}`}><span class="pdot"></span>{STATUS_LABEL[a.status]}</span>
          </div>
        {/each}
      </div>
    {/if}

    <!-- Applications table -->
    <div class="section-hd">
      <h2>Applications <span class="count">{apps.length}</span></h2>
      <div class="actions">
        <div class="filters">
          {#each chips as c}
            <button class={`chip ${filter === c.k ? 'active' : ''}`} onclick={() => (filter = c.k)}>
              {c.lbl}<span class="count">{c.n}</span>
            </button>
          {/each}
        </div>
      </div>
    </div>
    <div class="table">
      <div class="tr head">
        <span>Company</span>
        <span>Role</span>
        <span>Status</span>
        <span>Source</span>
        <span style="text-align:right">Applied</span>
        <span></span>
      </div>
      {#if loading}
        <div class="tr"><span style="color:var(--mute); grid-column: 1 / -1; text-align:center; padding:1rem 0;">Loading…</span></div>
      {:else if visible.length === 0}
        <div class="tr"><span style="color:var(--mute); grid-column: 1 / -1; text-align:center; padding:2rem 0;">
          {filter === 'all' && apps.length === 0
            ? 'No applications yet. Press ⌘N to add your first one.'
            : `No applications in “${filter}.”`}
        </span></div>
      {:else}
        {#each visible as a (a.id)}
          <div class="tr" onclick={() => open(a.id)}>
            <span class="co">
              <span class={`logo ${a.logoCls}`}>{a.coShort}</span>
              <span>{a.co}</span>
            </span>
            <span class="role">{a.role}</span>
            <span class={`pill ${a.status}`}><span class="pdot"></span>{STATUS_LABEL[a.status]}</span>
            <span class="source">{a.source}</span>
            <span class="applied">{a.applied}</span>
            <span class="arrow">→</span>
          </div>
        {/each}
      {/if}
    </div>

    <p class="disclaimer">
      Showing {visible.length} of {apps.length} · sorted by applied date · click any row to open
    </p>
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
          <span class="ai-tag">AI</span>
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

  /* Header */
  .m-head {
    display: flex; align-items: flex-start; justify-content: space-between;
    gap: 16px;
    padding: 20px 22px 14px;
  }
  .m-head h2 {
    font-size: 19px; font-weight: 500;
    letter-spacing: -0.018em;
    margin: 0;
    color: var(--ink);
  }
  .m-head .sub {
    margin: 4px 0 0;
    font-size: 13px; color: var(--mute);
    line-height: 1.5;
    max-width: 460px;
  }
  .x-close {
    background: transparent; border: 0; color: var(--mute);
    width: 28px; height: 28px; border-radius: 6px;
    display: grid; place-items: center; cursor: pointer;
    flex-shrink: 0;
    transition: background 100ms ease, color 100ms ease;
  }
  .x-close:hover { background: var(--surface-2); color: var(--ink); }

  /* Quick add zone */
  .quick-add {
    padding: 0 22px 16px;
    display: flex; flex-direction: column; gap: 10px;
  }
  .qa-head {
    display: flex; align-items: center; justify-content: space-between;
  }
  .qa-label {
    font-size: 12px; font-weight: 500;
    color: var(--ink-2);
    letter-spacing: 0;
  }
  .ai-tag {
    font-weight: 500; font-size: 10px;
    color: var(--accent-text); background: var(--accent-tint);
    border-radius: 4px; padding: 2px 7px;
    letter-spacing: .04em;
  }

  .drop {
    position: relative;
    background: var(--surface);
    border: 1.5px dashed var(--rule-strong);
    border-radius: 10px;
    transition: border-color 120ms ease, background 120ms ease;
    overflow: hidden;
  }
  .drop.drag {
    border-color: var(--accent);
    background: var(--accent-tint);
  }
  .drop.loading {
    border-color: var(--accent);
    border-style: solid;
  }
  .drop.loading::before {
    /* Sweeping accent bar across the top edge — purely cosmetic, gives the
       zone a "thinking" feel that's distinct from a static disabled state. */
    content: '';
    position: absolute;
    inset: 0 0 auto 0;
    height: 2px;
    background: linear-gradient(90deg, transparent, var(--accent), transparent);
    background-size: 40% 100%;
    background-repeat: no-repeat;
    animation: drop-sweep 1.4s linear infinite;
  }
  @keyframes drop-sweep {
    0%   { background-position: -40% 0; }
    100% { background-position: 140% 0; }
  }

  .drop-loading {
    display: flex; align-items: center; gap: 10px;
    padding: 10px 14px;
    background: var(--accent-tint);
    border-top: 1px solid var(--accent-tint-2);
    font-size: 12.5px;
    color: var(--accent-text);
    font-weight: 500;
  }
  .loading-text { letter-spacing: 0; }

  .spinner {
    display: inline-block;
    width: 14px; height: 14px;
    border: 1.8px solid var(--accent-tint-2);
    border-top-color: var(--accent);
    border-radius: 50%;
    animation: spinner-rot 0.7s linear infinite;
    flex-shrink: 0;
  }
  .spinner-sm {
    width: 12px; height: 12px;
    border-color: rgba(255,255,255,.35);
    border-top-color: white;
  }
  @keyframes spinner-rot {
    to { transform: rotate(360deg); }
  }

  .ok-check {
    display: inline-block;
    width: 14px; height: 14px; line-height: 14px;
    text-align: center;
    border-radius: 50%;
    background: var(--positive-tint);
    color: var(--positive-text);
    font-size: 10px; font-weight: 600;
    margin-right: 4px;
  }
  .drop textarea {
    width: 100%;
    font: inherit; font-family: var(--sans);
    font-size: 13.5px; color: var(--ink);
    background: transparent;
    border: 0;
    padding: 12px 14px 8px;
    outline: none; resize: none;
    min-height: 78px;
    line-height: 1.5;
    display: block;
  }
  .drop textarea::placeholder { color: var(--mute-2); }

  .attached {
    display: flex; align-items: center; gap: 10px;
    padding: 12px 14px 8px;
    font-size: 13px;
    color: var(--ink-2);
  }
  .attached svg { color: var(--accent-text); flex-shrink: 0; }
  .attached .att-name { font-weight: 500; color: var(--ink); }
  .attached .att-size { color: var(--mute); font-family: var(--mono); font-size: 11px; }
  .attached .att-x {
    margin-left: auto;
    background: transparent; border: 0; color: var(--mute);
    font-size: 18px; line-height: 1; cursor: pointer;
    padding: 0 4px;
  }
  .attached .att-x:hover { color: var(--ink); }

  .drop-footer {
    display: flex; align-items: center; gap: 8px;
    padding: 8px 14px 10px;
    font-size: 11.5px; color: var(--mute-2);
    border-top: 1px dashed var(--rule);
    background: transparent;
  }
  .drop-footer svg { color: var(--mute-2); flex-shrink: 0; }

  .qa-actions {
    display: flex; align-items: center; justify-content: space-between;
    gap: 12px;
  }
  .kb-hints {
    font-size: 11.5px; color: var(--mute);
  }
  .kb-hints kbd {
    font-family: var(--mono); font-size: 10.5px;
    background: var(--surface); border: 1px solid var(--rule);
    border-bottom-width: 2px; border-radius: 3px;
    padding: 0 4px; color: var(--ink-2);
  }

  .btn-accent {
    display: inline-flex; align-items: center; gap: 6px;
    background: var(--accent); color: white;
    border: 1px solid var(--accent);
    border-radius: 7px;
    padding: 7px 12px;
    font-size: 13px; font-weight: 500;
    cursor: pointer;
    transition: background 100ms ease, border-color 100ms ease;
  }
  .btn-accent:hover:not(:disabled) { background: var(--accent-strong); border-color: var(--accent-strong); }
  .btn-accent:disabled { opacity: .55; cursor: not-allowed; }
  .btn-accent svg { color: white; }

  .qa-msg {
    font-size: 12px;
  }
  .qa-msg.ok  { color: var(--positive-text); }
  .qa-msg.err { color: var(--danger-text); }

  /* Divider */
  .modal-divider {
    display: flex; align-items: center; gap: 12px;
    color: var(--mute-2);
    font-size: 10.5px; letter-spacing: .06em; text-transform: uppercase;
    padding: 0 22px;
    margin: 4px 0 14px;
  }
  .modal-divider::before, .modal-divider::after {
    content: ''; flex: 1; height: 1px; background: var(--rule);
  }

  /* Fields */
  .fields {
    padding: 0 22px 18px;
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 12px 14px;
  }
  .fields label {
    display: flex; flex-direction: column;
    gap: 5px;
  }
  .fields .lbl {
    font-size: 11.5px; color: var(--mute);
    font-weight: 500;
    letter-spacing: 0;
  }
  .fields .req { color: var(--accent-text); font-style: normal; margin-left: 1px; }
  .fields .opt { color: var(--mute-2); font-weight: 400; }

  .fields input {
    font: inherit; font-size: 13.5px;
    color: var(--ink); background: var(--surface);
    border: 1px solid var(--rule);
    border-radius: 8px;
    padding: 0 11px;
    height: 36px;
    outline: 0;
    transition: border-color 100ms ease, box-shadow 100ms ease, background 100ms ease;
  }
  .fields input:hover { border-color: var(--rule-strong); }
  .fields input:focus {
    border-color: var(--accent);
    box-shadow: 0 0 0 3px var(--accent-tint);
  }
  .fields input::placeholder { color: var(--mute-2); }

  .status-select {
    position: relative;
    display: flex; align-items: center;
    background: var(--surface);
    border: 1px solid var(--rule);
    border-radius: 8px;
    height: 36px;
    padding: 0 11px;
    transition: border-color 100ms ease, box-shadow 100ms ease;
  }
  .status-select:hover { border-color: var(--rule-strong); }
  .status-select:focus-within {
    border-color: var(--accent);
    box-shadow: 0 0 0 3px var(--accent-tint);
  }
  .status-select .sdot {
    width: 7px; height: 7px; border-radius: 50%;
    background: var(--mute);
    margin-right: 8px;
    flex-shrink: 0;
  }
  .status-select.wishlist  .sdot { background: var(--mute); }
  .status-select.applied   .sdot { background: var(--ink-2); }
  .status-select.screen    .sdot { background: var(--accent); }
  .status-select.interview .sdot { background: var(--warm, oklch(0.7 0.15 60)); }
  .status-select.offer     .sdot { background: var(--positive); }
  .status-select.rejected  .sdot { background: var(--mute-2); }
  .status-select.withdrawn .sdot { background: var(--mute-2); }

  .status-select select {
    flex: 1;
    appearance: none; -webkit-appearance: none;
    background: transparent;
    border: 0; outline: 0;
    font: inherit; font-size: 13.5px; color: var(--ink);
    padding: 0 18px 0 0;
    cursor: pointer;
  }
  .status-select .chev {
    position: absolute; right: 11px; top: 50%;
    transform: translateY(-50%);
    color: var(--mute);
    pointer-events: none;
  }

  /* Footer zone */
  .m-foot {
    display: flex; align-items: center; justify-content: space-between;
    padding: 12px 22px;
    background: var(--surface);
    border-top: 1px solid var(--rule);
  }
  .esc-hint {
    font-size: 11.5px; color: var(--mute);
  }
  .esc-hint kbd {
    font-family: var(--mono); font-size: 10.5px;
    background: var(--card); border: 1px solid var(--rule);
    border-bottom-width: 2px; border-radius: 3px;
    padding: 0 5px; color: var(--ink-2);
  }
  .foot-actions { display: flex; gap: 8px; }
  .foot-actions .btn-primary { display: inline-flex; align-items: center; gap: 6px; }
  .dark-kbd {
    font-family: var(--mono); font-size: 10.5px;
    background: rgba(255,255,255,.18);
    border: 1px solid rgba(255,255,255,.22);
    border-radius: 3px;
    padding: 0 5px;
    color: rgba(255,255,255,.9);
  }
</style>
