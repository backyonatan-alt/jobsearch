<script>
  import { onMount } from 'svelte';
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api.js';
  import { toDisplayApp, STATUS_LABEL, STATUSES, fmtLongDate } from '$lib/app-helpers.js';
  import { buildTimelineFromApplication } from '$lib/dossier-sample.js';

  let app = $state(null);
  let loading = $state(true);
  let notFound = $state(false);
  let tab = $state('dossier');

  // Interviews tab state
  let interviews = $state([]);
  let interviewsLoading = $state(false);
  let icsText = $state('');
  let icsParsing = $state(false);
  let icsParseError = $state('');
  let icsPreview = $state([]); // parsed events awaiting save
  let icsSaving = $state(false);

  // Inline-action state
  let showStatusMenu = $state(false);
  let showEditModal = $state(false);
  let edit = $state({ company: '', role: '', source: '', location: '', cv_variant: '', jd_url: '', salary_note: '' });
  let saving = $state(false);

  // Dossier state
  let dossier = $state(null);          // { content, meeting, generatedAgo, interviewer_name, ... } | null
  let dossierLoading = $state(false);
  let dossierError = $state('');
  let interviewerInput = $state('');
  let generating = $state(false);

  const id = $derived(page.params.id);

  $effect(() => {
    void id; // re-fetch when the path param changes
    loadApp();
    loadDossier();
    loadInterviews();
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

  async function loadDossier() {
    dossierLoading = true;
    dossierError = '';
    try {
      const d = await api(`/api/applications/${id}/dossier`);
      dossier = d;
      interviewerInput = d.interviewer_name ?? '';
    } catch (e) {
      // 404 = no dossier yet, that's normal — keep dossier = null
      if (!String(e.message).toLowerCase().includes('no dossier') && e.message !== 'unauthorized') {
        console.error(e);
      }
      dossier = null;
    } finally {
      dossierLoading = false;
    }
  }

  async function generateDossier() {
    if (generating) return;
    generating = true;
    dossierError = '';
    try {
      const d = await api(`/api/applications/${id}/dossier/refresh`, {
        method: 'POST',
        body: JSON.stringify({ interviewer_name: interviewerInput.trim() })
      });
      dossier = d;
      interviewerInput = d.interviewer_name ?? '';
    } catch (e) {
      dossierError = e.message || 'Could not generate dossier.';
    } finally {
      generating = false;
    }
  }

  async function loadInterviews() {
    interviewsLoading = true;
    try {
      interviews = await api(`/api/applications/${id}/interviews`);
    } catch (e) {
      if (e.message !== 'unauthorized') console.error(e);
      interviews = [];
    } finally {
      interviewsLoading = false;
    }
  }

  async function onIcsFile(e) {
    const f = e.target.files?.[0];
    e.target.value = '';
    if (!f) return;
    if (f.size > 256 * 1024) {
      icsParseError = 'File too large (256 KB max).';
      return;
    }
    icsText = await f.text();
    await parseIcs();
  }

  async function parseIcs() {
    if (icsParsing) return;
    icsParseError = '';
    icsPreview = [];
    const body = icsText.trim();
    if (!body) {
      icsParseError = 'Paste the .ics contents or pick a file.';
      return;
    }
    icsParsing = true;
    try {
      const r = await api(`/api/applications/${id}/interviews/parse`, {
        method: 'POST',
        body: JSON.stringify({ ics: body })
      });
      icsPreview = r.events ?? [];
      if (icsPreview.length === 0) icsParseError = 'No events found in that file.';
    } catch (e) {
      icsParseError = e.message || 'Could not parse calendar.';
    } finally {
      icsParsing = false;
    }
  }

  async function saveParsedEvents() {
    if (icsSaving || icsPreview.length === 0) return;
    icsSaving = true;
    try {
      for (const ev of icsPreview) {
        await api(`/api/applications/${id}/interviews`, {
          method: 'POST',
          body: JSON.stringify(ev)
        });
      }
      icsText = '';
      icsPreview = [];
      await loadInterviews();
      await loadDossier(); // refresh meeting block in dossier tab
    } catch (e) {
      icsParseError = e.message || 'Could not save events.';
    } finally {
      icsSaving = false;
    }
  }

  async function deleteInterview(iv) {
    if (!confirm(`Delete "${iv.summary}"?`)) return;
    await api(`/api/applications/${id}/interviews/${iv.id}`, { method: 'DELETE' });
    await loadInterviews();
    await loadDossier();
  }

  function fmtEventWhen(ev) {
    const d = new Date(ev.starts_at);
    if (ev.all_day) {
      return d.toLocaleDateString(undefined, { weekday: 'short', month: 'short', day: 'numeric' }) + ' · all day';
    }
    const date = d.toLocaleDateString(undefined, { weekday: 'short', month: 'short', day: 'numeric' });
    const time = d.toLocaleTimeString(undefined, { hour: 'numeric', minute: '2-digit' });
    let suffix = '';
    if (ev.ends_at) {
      const end = new Date(ev.ends_at);
      const mins = Math.round((end - d) / 60000);
      if (mins > 0) suffix = ` · ${mins >= 60 && mins % 60 === 0 ? `${mins / 60}h` : `${mins} min`}`;
    }
    return `${date}, ${time}${suffix}`;
  }

  function isPast(ev) {
    return new Date(ev.starts_at) < new Date();
  }

  // Meeting hero formatters — used when the dossier has a real interview
  // linked. starts_at/ends_at are ISO timestamps from the server; format in
  // the viewer's TZ so the hero matches the Scheduled list below.
  function meetingWhen(meeting) {
    if (!meeting?.starts_at) return meeting?.when ?? 'Upcoming · time TBD';
    const d = new Date(meeting.starts_at);
    if (meeting.all_day) {
      return d.toLocaleDateString(undefined, { weekday: 'short', month: 'short', day: 'numeric' }) + ' · all day';
    }
    const now = new Date();
    const startOfDay = (x) => new Date(x.getFullYear(), x.getMonth(), x.getDate());
    const days = Math.round((startOfDay(d) - startOfDay(now)) / 86400000);
    const time = d.toLocaleTimeString(undefined, { hour: 'numeric', minute: '2-digit' });
    if (days === 0) return `Today · ${time}`;
    if (days === 1) return `Tomorrow · ${time}`;
    if (days > 1 && days < 7) return `${d.toLocaleDateString(undefined, { weekday: 'long' })} · ${time}`;
    if (days < 0) return `${d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })} · ${time} (past)`;
    return `${d.toLocaleDateString(undefined, { weekday: 'short', month: 'short', day: 'numeric' })} · ${time}`;
  }

  function meetingDuration(meeting) {
    if (!meeting?.starts_at || !meeting?.ends_at) return meeting?.duration ?? '—';
    const mins = Math.round((new Date(meeting.ends_at) - new Date(meeting.starts_at)) / 60000);
    if (mins <= 0) return meeting.duration ?? '—';
    return mins >= 60 && mins % 60 === 0 ? `${mins / 60}h` : `${mins} min`;
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

  // The dossier feature only makes sense for live interview-loop apps.
  const dossierEligible = $derived(app && ['screen', 'interview', 'offer'].includes(app.status));
  const dossierAvailable = $derived(!!dossier);

  const timeline = $derived.by(() => {
    if (!app) return [];
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
      <span>Search applications, people…</span>
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
        <div class={`tab ${tab === 'interviews' ? 'active' : ''}`} onclick={() => (tab = 'interviews')}>
          Interviews
          <span class="count">{interviews.length}</span>
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
                {meetingWhen(dossier.meeting)}
                <span class="live-tag"><span class="pulse"></span> upcoming</span>
              </h3>
              <div class="when">
                <span>{meetingDuration(dossier.meeting)}</span>
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
              <a class="regen" onclick={generateDossier} class:busy={generating}>
                <svg width="11" height="11" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M2 6a4 4 0 1 1 1.2 2.8M2 4v2h2"/></svg>
                {generating ? 'Regenerating…' : 'Regenerate'}
              </a>
            </div>
          </div>

          {#if dossierError}
            <p class="dossier-err">{dossierError}</p>
          {/if}

          <!-- Interviewer card -->
          {#if dossier.content.interviewer}
            <div class="interviewer">
              <div class="photo">{dossier.content.interviewer.initials ?? ''}</div>
              <div class="who">
                <h4>{dossier.content.interviewer.name ?? ''}</h4>
                <div class="role">{dossier.content.interviewer.role ?? ''}</div>
                {#if dossier.content.interviewer.prior?.length}
                  <div class="prior">
                    <b>Prior:</b> {dossier.content.interviewer.prior.join(' · ')}
                  </div>
                {/if}
              </div>
              <div class="links">
                {#each dossier.content.interviewer.links ?? [] as l}
                  <a href={l.href} target="_blank" rel="noopener">
                    <svg width="14" height="14" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5" style="opacity:.55">
                      <rect x="2" y="2" width="12" height="12" rx="3"/>
                    </svg>
                    {l.label}
                    <span class="ext">↗</span>
                  </a>
                {/each}
              </div>
            </div>
          {/if}

          <!-- Snapshot -->
          {#if dossier.content.snapshot}
            <p class="snapshot">{@html dossier.content.snapshot}</p>
          {/if}

          <!-- Background + signals -->
          <div class="grid-2">
            {#if dossier.content.background}
              <div class="section">
                <h3>Background</h3>
                <p>{dossier.content.background}</p>
              </div>
            {/if}
            {#if dossier.content.signals?.length}
              <div class="section">
                <h3>Recent signals <span class="num">last 90 days</span></h3>
                <ul class="signals">
                  {#each dossier.content.signals as s}
                    <li>
                      <span class="date">{s.date}</span>
                      <span class="body">
                        {#if s.kind}<span class="kind">{s.kind}</span>{/if}
                        {s.body}
                        {#if s.source}<span class="source">{s.source} ↗</span>{/if}
                      </span>
                    </li>
                  {/each}
                </ul>
              </div>
            {/if}
          </div>

          <!-- Style block -->
          {#if dossier.content.style?.lead || dossier.content.style?.tells?.length}
            <div class="style-block">
              <h3 style="font-size:13px; font-weight:500; color:var(--mute); margin:0 0 12px;">How they likely interview</h3>
              {#if dossier.content.style.lead}<p class="lead">{dossier.content.style.lead}</p>{/if}
              {#if dossier.content.style.tells?.length}
                <div class="tells">
                  {#each dossier.content.style.tells as t}
                    <div class="tell">
                      <div class="lbl">{t.lbl}</div>
                      <div class="val">{t.val}</div>
                    </div>
                  {/each}
                </div>
              {/if}
            </div>
          {/if}

          <!-- Lands / Avoid -->
          {#if dossier.content.lands?.length || dossier.content.avoid?.length}
            <div class="la-grid">
              <div class="la-col lands">
                <h3><span class="glyph">+</span> What lands well</h3>
                <ul>
                  {#each dossier.content.lands ?? [] as line}
                    <li><span class="glyph">+</span><span>{line}</span></li>
                  {/each}
                </ul>
              </div>
              <div class="la-col avoid">
                <h3><span class="glyph">−</span> What to avoid</h3>
                <ul>
                  {#each dossier.content.avoid ?? [] as line}
                    <li><span class="glyph">−</span><span>{line}</span></li>
                  {/each}
                </ul>
              </div>
            </div>
          {/if}

          <!-- Questions -->
          {#if dossier.content.questions?.length}
            <div class="questions">
              <div class="brief-meta" style="margin-bottom:0">
                <h2 style="font-size:16px; font-weight:500;">Questions worth asking</h2>
                <span class="gen"><span>Ranked by signal</span></span>
              </div>
              <p class="intro">Saving any of these adds them to your interview prep doc.</p>
              <ol>
                {#each dossier.content.questions as q}
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
          {/if}

          <div class="disclaimer">
            Briefing synthesised from public posts, talks, and papers · always verify before you walk in · last refreshed {dossier.generatedAgo}
          </div>
        {:else if generating}
          <div class="generating-card">
            <div class="spinner"></div>
            <h3>Researching {app.co}{interviewerInput ? ` & ${interviewerInput}` : ''}…</h3>
            <p>Claude is searching the web for recent posts, talks, and the company's current direction. Typically 30–60 seconds.</p>
          </div>
        {:else}
          <div class="generate-card">
            <h3>Generate the dossier</h3>
            <p>
              Claude reads the web — recent essays, talks, papers, company news —
              and writes you a focused briefing for this interview. Optional: name
              the interviewer for a person-specific brief.
            </p>
            <div class="generate-row">
              <input
                type="text"
                placeholder="Interviewer name (optional) — e.g. Dario Amodei"
                bind:value={interviewerInput}
                disabled={generating}
              />
              <button class="btn btn-primary" onclick={generateDossier} disabled={generating || !dossierEligible}>
                Generate
              </button>
            </div>
            {#if !dossierEligible}
              <p class="muted-note">Move this application to <b>Screen</b>, <b>Interview</b>, or <b>Offer</b> to enable the dossier.</p>
            {/if}
            {#if dossierError}
              <p class="dossier-err">{dossierError}</p>
            {/if}
          </div>
        {/if}
      {:else if tab === 'interviews'}
        <div class="ics-card">
          <h3>Add an interview from a calendar invite</h3>
          <p class="ics-hint">
            Drop or pick the <code>.ics</code> file your interviewer sent — or paste the raw text below.
            We parse it locally (no AI round-trip) and link the event to this application.
          </p>
          <div class="ics-row">
            <label class="btn ics-file">
              Choose .ics file…
              <input type="file" accept=".ics,text/calendar" onchange={onIcsFile} hidden />
            </label>
            <span class="ics-sep">or paste below</span>
          </div>
          <textarea
            class="ics-textarea"
            placeholder={"BEGIN:VCALENDAR\nVERSION:2.0\nBEGIN:VEVENT\n…"}
            bind:value={icsText}
            rows="6"
          ></textarea>
          <div class="ics-actions">
            <button class="btn btn-primary" onclick={parseIcs} disabled={icsParsing || !icsText.trim()}>
              {icsParsing ? 'Parsing…' : 'Parse'}
            </button>
            {#if icsText || icsPreview.length}
              <button class="btn" onclick={() => { icsText = ''; icsPreview = []; icsParseError = ''; }}>Clear</button>
            {/if}
          </div>
          {#if icsParseError}
            <p class="dossier-err" style="margin-top:12px">{icsParseError}</p>
          {/if}

          {#if icsPreview.length > 0}
            <div class="ics-preview">
              <h4>Preview — {icsPreview.length} event{icsPreview.length === 1 ? '' : 's'}</h4>
              {#each icsPreview as ev}
                <div class="ics-preview-row">
                  <div class="iv-when">{fmtEventWhen(ev)}</div>
                  <div class="iv-summary">{ev.summary || 'Interview'}</div>
                  {#if ev.location}<div class="iv-loc">{ev.location}</div>{/if}
                </div>
              {/each}
              <div class="ics-actions" style="margin-top:12px">
                <button class="btn btn-primary" onclick={saveParsedEvents} disabled={icsSaving}>
                  {icsSaving ? 'Saving…' : `Save ${icsPreview.length} event${icsPreview.length === 1 ? '' : 's'}`}
                </button>
              </div>
            </div>
          {/if}
        </div>

        <div class="iv-list">
          <h3>Scheduled</h3>
          {#if interviewsLoading}
            <p style="color:var(--mute)">Loading…</p>
          {:else if interviews.length === 0}
            <p style="color:var(--mute); font-size:13px;">No interviews scheduled yet.</p>
          {:else}
            {#each interviews as iv}
              <div class={`iv-card ${isPast(iv) ? 'past' : ''}`}>
                <div class="iv-card-main">
                  <div class="iv-when">{fmtEventWhen(iv)}{isPast(iv) ? ' · past' : ''}</div>
                  <div class="iv-summary">{iv.summary}</div>
                  {#if iv.location}<div class="iv-loc">{iv.location}</div>{/if}
                  {#if iv.attendees && iv.attendees.length > 0}
                    <div class="iv-att">
                      with {iv.attendees.map(a => a.name || a.email).filter(Boolean).join(', ')}
                    </div>
                  {/if}
                </div>
                <button class="btn btn-danger" onclick={() => deleteInterview(iv)} title="Delete this interview">Delete</button>
              </div>
            {/each}
          {/if}
        </div>
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

  /* Dossier — generate / regenerate states */
  .generate-card {
    border: 1px dashed var(--rule-strong);
    border-radius: 12px;
    padding: 28px 32px;
    background: var(--card);
    text-align: center;
  }
  .generate-card h3 {
    font-size: 17px; font-weight: 500;
    letter-spacing: -0.012em;
    margin: 0 0 .35rem;
    color: var(--ink);
  }
  .generate-card > p {
    margin: 0 0 1.25rem;
    color: var(--mute);
    font-size: 13.5px;
    line-height: 1.55;
    max-width: 56ch;
    margin-left: auto;
    margin-right: auto;
  }
  .generate-row {
    display: flex;
    gap: 8px;
    max-width: 520px;
    margin: 0 auto;
  }
  .generate-row input {
    flex: 1;
    font: inherit;
    font-size: 13.5px;
    color: var(--ink);
    background: var(--surface);
    border: 1px solid var(--rule);
    border-radius: 6px;
    padding: 6px 10px;
    outline: none;
  }
  .generate-row input:focus { border-color: var(--accent); }
  .muted-note { color: var(--mute); font-size: 12px; margin: .75rem 0 0; }
  .muted-note b { color: var(--ink-2); font-weight: 500; }
  .dossier-err {
    color: var(--danger-text);
    background: var(--danger-tint);
    border: 1px solid var(--danger-tint);
    border-radius: 8px;
    padding: 8px 12px;
    font-size: 13px;
    margin: 0 0 16px;
  }

  .generating-card {
    border: 1px solid var(--accent-tint-2);
    background: linear-gradient(180deg, var(--accent-tint), var(--card) 80%);
    border-radius: 12px;
    padding: 32px;
    text-align: center;
  }
  .generating-card h3 {
    font-size: 17px; font-weight: 500;
    letter-spacing: -0.012em;
    margin: 16px 0 .5rem;
    color: var(--ink);
  }
  .generating-card p {
    color: var(--mute);
    font-size: 13.5px;
    margin: 0;
    max-width: 56ch;
    margin-left: auto;
    margin-right: auto;
  }
  .spinner {
    width: 24px; height: 24px;
    margin: 0 auto;
    border: 2px solid var(--accent-tint-2);
    border-top-color: var(--accent);
    border-radius: 50%;
    animation: spin 800ms linear infinite;
  }
  @keyframes spin { to { transform: rotate(360deg); } }

  .regen { cursor: pointer; user-select: none; }
  .regen.busy { opacity: 0.6; pointer-events: none; }

  /* Interviews tab */
  .ics-card {
    border: 1px dashed var(--rule-strong);
    border-radius: 12px;
    padding: 20px 24px;
    background: var(--card);
    margin-bottom: 24px;
  }
  .ics-card h3 {
    font-size: 15px; font-weight: 500;
    letter-spacing: -0.012em;
    margin: 0 0 .35rem;
    color: var(--ink);
  }
  .ics-hint {
    font-size: 13px;
    color: var(--mute);
    margin: 0 0 1rem;
    line-height: 1.5;
  }
  .ics-hint code {
    background: var(--surface-2);
    padding: 1px 5px;
    border-radius: 3px;
    font-size: 12px;
  }
  .ics-row {
    display: flex; align-items: center; gap: 12px;
    margin-bottom: .75rem;
  }
  .ics-file {
    cursor: pointer;
  }
  .ics-sep {
    color: var(--mute);
    font-size: 12px;
  }
  .ics-textarea {
    width: 100%;
    font: inherit;
    font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
    font-size: 12px;
    color: var(--ink);
    background: var(--surface);
    border: 1px solid var(--rule);
    border-radius: 6px;
    padding: 8px 10px;
    outline: none;
    resize: vertical;
    box-sizing: border-box;
  }
  .ics-textarea:focus { border-color: var(--accent); }
  .ics-actions {
    display: flex; gap: 8px;
    margin-top: 10px;
  }
  .ics-preview {
    margin-top: 16px;
    padding-top: 16px;
    border-top: 1px solid var(--rule);
  }
  .ics-preview h4 {
    font-size: 12px;
    font-weight: 500;
    color: var(--mute);
    text-transform: uppercase;
    letter-spacing: 0.04em;
    margin: 0 0 8px;
  }
  .ics-preview-row {
    padding: 8px 0;
    border-bottom: 1px solid var(--rule);
  }
  .ics-preview-row:last-of-type { border-bottom: none; }

  .iv-list h3 {
    font-size: 13px;
    font-weight: 500;
    color: var(--mute);
    text-transform: uppercase;
    letter-spacing: 0.04em;
    margin: 0 0 12px;
  }
  .iv-card {
    display: flex; align-items: flex-start; justify-content: space-between;
    gap: 16px;
    padding: 14px 16px;
    border: 1px solid var(--rule);
    border-radius: 8px;
    background: var(--card);
    margin-bottom: 8px;
  }
  .iv-card.past { opacity: 0.6; }
  .iv-card-main { flex: 1; min-width: 0; }
  .iv-when {
    font-size: 12px;
    color: var(--accent-text);
    font-weight: 500;
    margin-bottom: 2px;
  }
  .iv-summary {
    font-size: 14px;
    color: var(--ink);
    margin-bottom: 2px;
  }
  .iv-loc, .iv-att {
    font-size: 12px;
    color: var(--mute);
    margin-top: 2px;
  }
</style>
