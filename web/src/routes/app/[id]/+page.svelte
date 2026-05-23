<script>
  import { onMount } from 'svelte';
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api.js';
  import { toDisplayApp, STATUS_LABEL, STATUSES, fmtLongDate } from '$lib/app-helpers.js';
  import { SAMPLE_DOSSIER, buildTimelineFromApplication } from '$lib/dossier-sample.js';

  let app = $state(null);
  let loading = $state(true);
  let notFound = $state(false);
  let tab = $state('dossier');

  // Inline-action state
  let showStatusMenu = $state(false);
  let showEditModal = $state(false);
  let edit = $state({ company: '', role: '', source: '', location: '', cv_variant: '', jd_url: '', salary_note: '' });
  let saving = $state(false);

  const id = $derived(page.params.id);

  $effect(() => {
    void id; // re-fetch when the path param changes
    loadApp();
  });

  async function loadApp() {
    loading = true;
    notFound = false;
    try {
      const raw = await api(`/api/applications/${id}`);
      app = toDisplayApp(raw);
    } catch (e) {
      if (e.message === 'unauthorized') return;
      if (e.message.includes('not found') || e.message.includes('404')) notFound = true;
      else console.error(e);
    } finally {
      loading = false;
    }
  }

  async function setStatus(newStatus) {
    showStatusMenu = false;
    if (!app || newStatus === app.status) return;
    await api(`/api/applications/${id}`, { method: 'PATCH', body: JSON.stringify({ status: newStatus }) });
    await loadApp();
  }

  function openEdit() {
    if (!app) return;
    edit = {
      company:     app.raw.company ?? '',
      role:        app.raw.role ?? '',
      source:      app.raw.source ?? '',
      location:    app.raw.location ?? '',
      cv_variant:  app.raw.cv_variant ?? '',
      jd_url:      app.raw.jd_url ?? '',
      salary_note: app.raw.salary_note ?? ''
    };
    showEditModal = true;
  }

  async function saveEdit(e) {
    e.preventDefault();
    saving = true;
    try {
      // Send the full set so blank inputs clear stale values. The backend
      // PATCH uses COALESCE on NULLIF; an empty string keeps the existing
      // value, which is the wrong semantic here. Workaround: send fields
      // wrapped so the server can distinguish "unchanged" from "clear" —
      // for now we treat blank as "keep", matching the modal hint copy.
      const payload = { ...edit };
      for (const k of Object.keys(payload)) if (!payload[k]) delete payload[k];
      await api(`/api/applications/${id}`, { method: 'PATCH', body: JSON.stringify(payload) });
      showEditModal = false;
      await loadApp();
    } finally {
      saving = false;
    }
  }

  async function deleteApp() {
    if (!app) return;
    if (!confirm(`Delete the ${app.co} application? This can't be undone.`)) return;
    await api(`/api/applications/${id}`, { method: 'DELETE' });
    goto('/app', { replaceState: true });
  }

  function back() { goto('/app'); }

  // For the v0.1 demo, we render the rich SAMPLE_DOSSIER on any application
  // whose status indicates the interview loop is live. Otherwise show empty.
  const dossierAvailable = $derived(app && ['interview', 'screen', 'offer'].includes(app.status));
  const dossier = SAMPLE_DOSSIER;

  const timeline = $derived.by(() => {
    if (!app) return [];
    if (dossierAvailable && app.co === 'Anthropic') return dossier.timeline;
    return buildTimelineFromApplication(app.raw);
  });
</script>

<svelte:head>
  <title>{app?.co ? `${app.co} — Pursuit` : 'Pursuit'}</title>
</svelte:head>

<div class="topbar">
  <div class="crumb">
    <span class="root" onclick={back}>Applications</span>
    <span class="sep">/</span>
    <span class="here">{app?.co ?? '…'}</span>
  </div>
  <div class="right">
    <div class="search">
      <svg class="ico" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5">
        <circle cx="7" cy="7" r="4.5"/><path d="M10.5 10.5L13 13"/>
      </svg>
      Search applications, people…
      <span class="kbd">⌘K</span>
    </div>
    <div class="status-wrap">
      <button class="btn" onclick={() => (showStatusMenu = !showStatusMenu)} disabled={!app}>
        Update status
        <svg width="11" height="11" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.5" style="margin-left: 2px"><path d="M3 4.5l3 3 3-3" stroke-linecap="round" stroke-linejoin="round"/></svg>
      </button>
      {#if showStatusMenu}
        <div class="status-menu" role="menu" onclick={(e) => e.stopPropagation()}>
          {#each STATUSES as s}
            <button
              class="status-menu-item"
              class:current={app?.status === s}
              onclick={() => setStatus(s)}
            >
              <span class={`pill ${s}`} style="margin: 0"><span class="pdot"></span>{STATUS_LABEL[s]}</span>
              {#if app?.status === s}<span class="check">✓</span>{/if}
            </button>
          {/each}
        </div>
      {/if}
    </div>
    <button class="btn" onclick={openEdit} disabled={!app}>Edit</button>
    <button class="btn btn-danger" onclick={deleteApp} disabled={!app} title="Delete this application">Delete</button>
  </div>
</div>

<!-- click-out to close the status menu -->
{#if showStatusMenu}
  <div class="menu-scrim" onclick={() => (showStatusMenu = false)} role="presentation"></div>
{/if}

<div class="body">
  <div class="body-inner">
    {#if loading}
      <p style="color:var(--mute)">Loading…</p>
    {:else if notFound || !app}
      <div class="empty-tab">
        <h3>Application not found</h3>
        <p>It may have been deleted, or you might not have access. <a href="/app" style="color:var(--accent-text)">Back to Today →</a></p>
      </div>
    {:else}
      <!-- Application header -->
      <div class="dossier-hd">
        <div class="top">
          <span class={`logo-big ${app.logoCls}`}>{app.coShort}</span>
          <div>
            <h1>
              {app.co}
              <span class="role">{app.role}</span>
            </h1>
            <div class="sub">
              <span>Applied {app.applied}</span>
              <span class="dot">·</span>
              <span>via {app.source}</span>
              {#if app.cv && app.cv !== '—'}
                <span class="dot">·</span>
                <span>CV {app.cv}</span>
              {/if}
              {#if app.raw.jd_url}
                <span class="dot">·</span>
                <a href={app.raw.jd_url} target="_blank" rel="noopener">
                  JD
                  <svg width="11" height="11" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M4 8L8 4M5 4h3v3"/></svg>
                </a>
              {/if}
            </div>
          </div>
          <div class="actions">
            <span class={`pill ${app.status}`}><span class="pdot"></span>{STATUS_LABEL[app.status]}</span>
          </div>
        </div>
      </div>

      <!-- Tabs -->
      <div class="tabs">
        <div class={`tab ${tab === 'dossier' ? 'active' : ''}`} onclick={() => (tab = 'dossier')}>
          Dossier
          <span class="count">AI</span>
        </div>
        <div class={`tab ${tab === 'timeline' ? 'active' : ''}`} onclick={() => (tab = 'timeline')}>
          Timeline
          <span class="count">{timeline.length}</span>
        </div>
        <div class={`tab ${tab === 'notes' ? 'active' : ''}`} onclick={() => (tab = 'notes')}>Notes</div>
        <div class={`tab ${tab === 'files' ? 'active' : ''}`} onclick={() => (tab = 'files')}>Files</div>
      </div>

      {#if tab === 'dossier'}
        {#if dossierAvailable}
          <!-- Meeting hero -->
          <div class="meeting">
            <div class="lhs">
              <h3>
                {dossier.meeting.when}
                <span class="live-tag"><span class="pulse"></span> upcoming</span>
              </h3>
              <div class="when">
                <span>{dossier.meeting.duration}</span>
                <span class="dot">·</span>
                <span>{dossier.meeting.medium}</span>
                <span class="dot">·</span>
                <span>{dossier.meeting.panel}</span>
              </div>
            </div>
            <div class="rhs">
              <button class="btn">Join meeting →</button>
            </div>
          </div>

          <!-- Briefing meta -->
          <div class="brief-meta">
            <h2>Interviewer briefing</h2>
            <div class="gen">
              <span class="ai-tag">AI generated</span>
              <span>Refreshed {dossier.generatedAgo}</span>
              <a>
                <svg width="11" height="11" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M2 6a4 4 0 1 1 1.2 2.8M2 4v2h2"/></svg>
                Regenerate
              </a>
            </div>
          </div>

          <!-- Interviewer card -->
          <div class="interviewer">
            <div class="photo">{dossier.interviewer.initials}</div>
            <div class="who">
              <h4>{dossier.interviewer.name}</h4>
              <div class="role">{dossier.interviewer.role}</div>
              <div class="prior">
                <b>Prior:</b> {dossier.interviewer.prior.join(' · ')}
              </div>
            </div>
            <div class="links">
              {#each dossier.interviewer.links as l}
                <a href={l.href}>
                  <svg width="14" height="14" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5" style="opacity:.55">
                    <rect x="2" y="2" width="12" height="12" rx="3"/>
                  </svg>
                  {l.label}
                  <span class="ext">↗</span>
                </a>
              {/each}
            </div>
          </div>

          <!-- Snapshot -->
          <p class="snapshot">{@html dossier.snapshot}</p>

          <!-- Background + signals -->
          <div class="grid-2">
            <div class="section">
              <h3>Background</h3>
              <p>{dossier.background}</p>
            </div>
            <div class="section">
              <h3>Recent signals <span class="num">last 90 days</span></h3>
              <ul class="signals">
                {#each dossier.signals as s}
                  <li>
                    <span class="date">{s.date}</span>
                    <span class="body">
                      <span class="kind">{s.kind}</span>
                      {s.body}
                      <span class="source">{s.source} ↗</span>
                    </span>
                  </li>
                {/each}
              </ul>
            </div>
          </div>

          <!-- Style block -->
          <div class="style-block">
            <h3 style="font-size:13px; font-weight:500; color:var(--mute); margin:0 0 12px;">How they likely interview</h3>
            <p class="lead">{dossier.style.lead}</p>
            <div class="tells">
              {#each dossier.style.tells as t}
                <div class="tell">
                  <div class="lbl">{t.lbl}</div>
                  <div class="val">{t.val}</div>
                </div>
              {/each}
            </div>
          </div>

          <!-- Lands / Avoid -->
          <div class="la-grid">
            <div class="la-col lands">
              <h3><span class="glyph">+</span> What lands well</h3>
              <ul>
                {#each dossier.lands as line}
                  <li><span class="glyph">+</span><span>{line}</span></li>
                {/each}
              </ul>
            </div>
            <div class="la-col avoid">
              <h3><span class="glyph">−</span> What to avoid</h3>
              <ul>
                {#each dossier.avoid as line}
                  <li><span class="glyph">−</span><span>{line}</span></li>
                {/each}
              </ul>
            </div>
          </div>

          <!-- Questions -->
          <div class="questions">
            <div class="brief-meta" style="margin-bottom:0">
              <h2 style="font-size:16px; font-weight:500;">Questions worth asking</h2>
              <span class="gen"><span>Ranked by signal</span></span>
            </div>
            <p class="intro">Saving any of these adds them to your interview prep doc.</p>
            <ol>
              {#each dossier.questions as q}
                <li>
                  <span></span>
                  <div>
                    <div class="q">"{q.q}"</div>
                    <div class="why">{q.why}</div>
                  </div>
                  <span class="save" title="Save to prep doc">＋</span>
                </li>
              {/each}
            </ol>
          </div>

          <div class="disclaimer">
            Briefing synthesised from public posts, talks, and papers · always verify before you walk in · last refreshed {dossier.generatedAgo}
          </div>
        {:else}
          <div class="empty-tab">
            <h3>No dossier yet</h3>
            <p>Dossier is generated once an interview is scheduled. Move this application to <b>Screen</b>, <b>Interview</b>, or <b>Offer</b> to unlock it.</p>
          </div>
        {/if}
      {:else if tab === 'timeline'}
        {#if timeline.length > 0}
          <div class="timeline">
            {#each timeline as e}
              <div class={`timeline-event ${e.tag || ''}`}>
                <span class="date">{e.date}</span>
                <span class="axis"><span class="marker"></span></span>
                <div class="what">
                  <div class="label">{e.label}</div>
                  <div class="note">{e.note}</div>
                </div>
              </div>
            {/each}
          </div>
        {:else}
          <div class="empty-tab">
            <h3>Timeline is empty</h3>
            <p>Events appear here as the application progresses. Add notes from the Notes tab.</p>
          </div>
        {/if}
      {:else if tab === 'notes'}
        <div class="empty-tab">
          <h3>Notes</h3>
          <p>Free-form notes with autosave — out of scope for this round.</p>
        </div>
      {:else if tab === 'files'}
        <div class="empty-tab">
          <h3>Files</h3>
          <p>CV variants attached to this application — out of scope for this round.</p>
        </div>
      {/if}
    {/if}
  </div>
</div>

{#if showEditModal}
  <div class="modal-overlay" onclick={() => (showEditModal = false)} role="presentation">
    <form class="modal" onclick={(e) => e.stopPropagation()} onsubmit={saveEdit}>
      <h2>Edit application</h2>
      <p class="modal-hint">Empty fields keep their current value.</p>
      <div class="fields">
        <label>Company <input bind:value={edit.company} required /></label>
        <label>Role <input bind:value={edit.role} required /></label>
        <label>Source <input bind:value={edit.source} placeholder="LinkedIn / Referral / Cold email" /></label>
        <label>Location <input bind:value={edit.location} placeholder="Remote / San Francisco" /></label>
        <label>CV variant <input bind:value={edit.cv_variant} placeholder="v3-ai-focus" /></label>
        <label>Salary note <input bind:value={edit.salary_note} placeholder="$220k-$280k base" /></label>
        <label class="span-2">JD URL <input bind:value={edit.jd_url} placeholder="https://…" /></label>
      </div>
      <div class="modal-actions">
        <button type="button" class="btn" onclick={() => (showEditModal = false)}>Cancel</button>
        <button type="submit" class="btn btn-primary" disabled={saving || !edit.company || !edit.role}>
          {saving ? 'Saving…' : 'Save changes'}
        </button>
      </div>
    </form>
  </div>
{/if}

<style>
  /* Status menu (popover under the Update status button) */
  .status-wrap { position: relative; }
  .status-menu {
    position: absolute; top: calc(100% + 6px); right: 0;
    z-index: 50;
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 8px;
    box-shadow: var(--sh-pop);
    padding: 4px;
    min-width: 180px;
    display: flex; flex-direction: column;
    gap: 1px;
  }
  .status-menu-item {
    display: flex; align-items: center; gap: 8px;
    padding: 6px 8px;
    border-radius: 5px;
    background: transparent;
    font: inherit; font-size: 13px;
    color: var(--ink-2);
    cursor: pointer;
    text-align: left;
    width: 100%;
  }
  .status-menu-item:hover { background: var(--surface-2); }
  .status-menu-item.current { background: var(--surface-2); }
  .status-menu-item .check { margin-left: auto; color: var(--accent-text); font-weight: 600; font-size: 12px; }

  .menu-scrim {
    position: fixed; inset: 0;
    z-index: 49;
    background: transparent;
  }

  /* Danger button variant */
  .btn-danger {
    color: var(--danger-text);
    border-color: var(--rule);
  }
  .btn-danger:hover {
    background: var(--danger-tint);
    border-color: var(--danger-tint);
  }

  /* Edit modal */
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
    max-width: 560px;
    display: flex; flex-direction: column; gap: .75rem;
    box-shadow: var(--sh-pop);
  }
  .modal h2 {
    font-size: 18px; font-weight: 500;
    letter-spacing: -0.018em;
    margin: 0;
  }
  .modal-hint {
    font-size: 12px;
    color: var(--mute);
    margin: 0 0 .5rem;
  }
  .fields {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: .65rem;
  }
  .fields .span-2 { grid-column: span 2; }
  .modal label {
    display: flex; flex-direction: column;
    font-size: 12px;
    color: var(--mute);
    gap: .35rem;
  }
  .modal input {
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
  .modal input:focus { border-color: var(--accent); }
  .modal-actions {
    display: flex; justify-content: flex-end; gap: .5rem;
    margin-top: .75rem;
  }
</style>
