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

  async function parseJD() {
    const text = pasteText.trim();
    if (text.length < 5) {
      parseError = 'Paste a job listing or URL first.';
      return;
    }
    parseError = '';
    parsing = true;
    try {
      const r = await api('/api/applications/parse', { method: 'POST', body: JSON.stringify({ text }) });
      // Overwrite form fields with whatever the model returned. Leave fields the
      // model didn't fill alone, so users can mix paste-then-correct.
      if (r.company)     form.company     = r.company;
      if (r.role)        form.role        = r.role;
      if (r.location)    form.location    = r.location;
      if (r.jd_url)      form.jd_url      = r.jd_url;
      if (r.source)      form.source      = r.source;
      if (r.salary_note) form.salary_note = r.salary_note;
      parsedHint = `Filled${r.seniority ? ` (level: ${r.seniority})` : ''} — review and Save.`;
      pasteText = '';
    } catch (e) {
      parseError = e.message || 'Parse failed.';
    } finally {
      parsing = false;
    }
  }

  function resetModal() {
    form = { company: '', role: '', status: 'applied', source: '', jd_url: '', cv_variant: '', location: '', salary_note: '' };
    pasteText = '';
    parseError = '';
    parsedHint = '';
    showNewModal = false;
  }

  onMount(refresh);

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
      Search applications, people…
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
    <div class="page-hd">
      <div>
        <div class="date">{todayString}</div>
        <h1>Today.</h1>
      </div>
      <div class="stats">
        {#each stats as s}
          <span>{s.lbl} <b>{s.n}</b></span>
        {/each}
      </div>
    </div>

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

{#if showNewModal}
  <div class="modal-overlay" onclick={resetModal} role="presentation">
    <form class="modal" onclick={(e) => e.stopPropagation()} onsubmit={createApp}>
      <h2>New application</h2>

      <!-- Paste-to-parse row -->
      <div class="paste-block">
        <div class="paste-label">
          <span class="ai-tag">AI</span>
          Paste a job posting, URL, or recruiter email — we'll fill the fields below.
        </div>
        <textarea
          bind:value={pasteText}
          placeholder="https://linkedin.com/jobs/… or the JD text…"
          rows="3"
          disabled={parsing}
        ></textarea>
        <div class="paste-row">
          <button type="button" class="btn" onclick={parseJD} disabled={parsing || !pasteText.trim()}>
            {parsing ? 'Parsing…' : 'Parse'}
          </button>
          {#if parsedHint}<span class="hint">{parsedHint}</span>{/if}
          {#if parseError}<span class="hint err">{parseError}</span>{/if}
        </div>
      </div>

      <div class="modal-divider"><span>or fill in by hand</span></div>

      <div class="fields">
        <label>Company <input bind:value={form.company} required /></label>
        <label>Role <input bind:value={form.role} required /></label>
        <label>Status
          <select bind:value={form.status}>
            <option value="wishlist">Wishlist</option>
            <option value="applied">Applied</option>
            <option value="screen">Screen</option>
            <option value="interview">Interview</option>
            <option value="offer">Offer</option>
            <option value="rejected">Rejected</option>
            <option value="withdrawn">Withdrawn</option>
          </select>
        </label>
        <label>Source <input bind:value={form.source} placeholder="LinkedIn / Referral / Cold email" /></label>
        <label>Location <input bind:value={form.location} placeholder="Remote / San Francisco" /></label>
        <label>CV variant <input bind:value={form.cv_variant} placeholder="v3-ai-focus" /></label>
        <label class="span-2">JD URL <input bind:value={form.jd_url} placeholder="https://…" /></label>
        <label class="span-2">Salary note <input bind:value={form.salary_note} placeholder="$220k-$280k base" /></label>
      </div>

      <div class="modal-actions">
        <button type="button" class="btn" onclick={resetModal}>Cancel</button>
        <button type="submit" class="btn btn-primary" disabled={saving || !form.company || !form.role}>
          {saving ? 'Saving…' : 'Add application'}
        </button>
      </div>
    </form>
  </div>
{/if}

<style>
  .modal-overlay {
    position: fixed; inset: 0;
    background: rgba(10,10,13,0.4);
    display: grid; place-items: center;
    z-index: 100;
    padding: 2rem;
  }
  .modal {
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 12px;
    padding: 1.5rem;
    width: 100%;
    max-width: 440px;
    display: flex; flex-direction: column; gap: .75rem;
    box-shadow: var(--sh-pop);
  }
  .modal {
    max-width: 560px !important;
  }
  .modal h2 {
    font-size: 18px; font-weight: 500;
    letter-spacing: -0.018em;
    margin: 0 0 .5rem;
  }
  .modal label {
    display: flex; flex-direction: column;
    font-size: 12px;
    color: var(--mute);
    gap: .35rem;
  }
  .modal input, .modal select, .modal textarea {
    font: inherit;
    color: var(--ink);
    background: var(--surface);
    border: 1px solid var(--rule);
    border-radius: 6px;
    padding: .45rem .6rem;
    font-size: 13.5px;
    outline: none;
    transition: border-color 100ms ease;
  }
  .modal textarea {
    resize: vertical;
    min-height: 70px;
    font-family: var(--sans);
    line-height: 1.4;
  }
  .modal input:focus, .modal select:focus, .modal textarea:focus {
    border-color: var(--accent);
  }
  .modal-actions {
    display: flex; justify-content: flex-end; gap: .5rem;
    margin-top: .75rem;
  }

  .paste-block { display: flex; flex-direction: column; gap: .4rem; }
  .paste-label {
    font-size: 12px; color: var(--ink-2);
    display: flex; align-items: center; gap: .4rem;
  }
  .ai-tag {
    font-weight: 500; font-size: 10px;
    color: var(--accent-text); background: var(--accent-tint);
    border-radius: 4px; padding: 1px 6px;
    letter-spacing: .04em;
  }
  .paste-row {
    display: flex; align-items: center; gap: .65rem;
  }
  .paste-row .hint {
    font-size: 12px; color: var(--positive-text);
  }
  .paste-row .hint.err { color: var(--danger-text); }

  .modal-divider {
    display: flex; align-items: center; gap: .65rem;
    color: var(--mute-2);
    font-size: 11px; letter-spacing: .04em; text-transform: uppercase;
    margin: .25rem 0 .25rem;
  }
  .modal-divider::before, .modal-divider::after {
    content: ''; flex: 1; height: 1px; background: var(--rule);
  }

  .fields {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: .65rem;
  }
  .fields .span-2 { grid-column: span 2; }
</style>
