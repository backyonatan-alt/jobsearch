<script>
  import { onMount } from 'svelte';
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api.js';
  import { toDisplayApp, STATUS_LABEL, STATUSES, fmtLongDate, daysSince } from '$lib/app-helpers.js';
  import { buildTimelineFromApplication } from '$lib/dossier-sample.js';

  let app = $state(null);
  let loading = $state(true);
  let notFound = $state(false);
  let tab = $state('brief');

  // Interviews tab state
  let interviews = $state([]);
  let interviewsLoading = $state(false);
  // Left zone — .ics
  let icsText = $state('');
  // Right zone — AI (screenshot or email body text)
  let aiText = $state('');
  let aiImage = $state(null); // { name, mediaType, size, file }
  let aiDragOver = $state(false);
  // Shared parse state
  let icsParsing = $state(false);
  let icsParseError = $state('');
  let icsPreview = $state([]);
  let icsSaving = $state(false);
  const AI_ALLOWED_IMG = ['image/png', 'image/jpeg', 'image/gif', 'image/webp'];
  const AI_MAX_BYTES = 6 * 1024 * 1024;
  function setAiImage(f) {
    if (!f) return;
    if (!AI_ALLOWED_IMG.includes(f.type)) { icsParseError = 'Only PNG / JPEG / GIF / WebP screenshots are supported.'; return; }
    if (f.size > AI_MAX_BYTES) { icsParseError = 'Screenshot too large (6 MB max).'; return; }
    icsParseError = '';
    aiImage = { name: f.name || 'pasted.png', mediaType: f.type, size: f.size, file: f };
  }
  function onAiDrop(e) {
    e.preventDefault();
    aiDragOver = false;
    setAiImage(e.dataTransfer?.files?.[0]);
  }
  function onAiDragOver(e) {
    if (e.dataTransfer?.types?.includes('Files')) { e.preventDefault(); aiDragOver = true; }
  }
  function onAiPaste(e) {
    const item = [...(e.clipboardData?.items || [])].find(i => i.type.startsWith('image/'));
    if (item) setAiImage(item.getAsFile());
  }
  function fileToBase64(f) {
    return new Promise((resolve, reject) => {
      const r = new FileReader();
      r.onerror = () => reject(new Error('Could not read file.'));
      r.onload = () => {
        const s = String(r.result || '');
        const i = s.indexOf(',');
        resolve(i === -1 ? '' : s.slice(i + 1));
      };
      r.readAsDataURL(f);
    });
  }

  // Inline-action state
  let showStatusMenu = $state(false);
  let showEditModal = $state(false);
  let edit = $state({ company: '', role: '', source: '', location: '', cv_variant: '', jd_url: '', salary_note: '', hiring_manager_name: '', hiring_manager_linkedin: '' });
  let saving = $state(false);

  // Dossier state
  let dossier = $state(null);
  let dossierLoading = $state(false);
  let dossierError = $state('');
  let interviewerInput = $state('');
  let generating = $state(false);

  const id = $derived(page.params.id);

  $effect(() => {
    void id;
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

  // Auto-trigger dossier on first load. Fires once per page-mount only when:
  // the dossier is genuinely absent (load returned no_dossier), the app is in
  // a stage that benefits from a brief (screen/interview/offer), and we
  // haven't already kicked one off this session for this app id.
  let autoTriggered = $state(false);
  $effect(() => {
    if (
      !autoTriggered &&
      !dossierLoading &&
      !dossierAvailable &&
      !generating &&
      dossierEligible &&
      app
    ) {
      autoTriggered = true;
      generateDossier();
    }
  });

  async function loadDossier() {
    dossierLoading = true;
    dossierError = '';
    try {
      const d = await api(`/api/applications/${id}/dossier`);
      dossier = d;
      interviewerInput = d.interviewer_name ?? '';
    } catch (e) {
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

  // Parse the AI zone — sends image and/or text to the unified parse endpoint.
  async function parseAi() {
    if (icsParsing) return;
    icsParseError = '';
    icsPreview = [];
    if (!aiImage && !aiText.trim()) {
      icsParseError = 'Drop a screenshot or paste the event text first.';
      return;
    }
    icsParsing = true;
    try {
      const payload = {};
      if (aiText.trim()) payload.text = aiText.trim();
      if (aiImage) {
        const data = await fileToBase64(aiImage.file);
        payload.image = { media_type: aiImage.mediaType, data };
      }
      const r = await api(`/api/applications/${id}/interviews/parse`, {
        method: 'POST',
        body: JSON.stringify(payload)
      });
      icsPreview = r.events ?? [];
      if (icsPreview.length === 0) icsParseError = "Couldn't extract an event.";
    } catch (e) {
      icsParseError = e.message || 'Could not parse.';
    } finally {
      icsParsing = false;
    }
  }
  async function onAiFile(e) {
    const f = e.target.files?.[0];
    e.target.value = '';
    setAiImage(f);
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
      await loadDossier();
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
      company:                 app.raw.company ?? '',
      role:                    app.raw.role ?? '',
      source:                  app.raw.source ?? '',
      location:                app.raw.location ?? '',
      cv_variant:              app.raw.cv_variant ?? '',
      jd_url:                  app.raw.jd_url ?? '',
      salary_note:             app.raw.salary_note ?? '',
      hiring_manager_name:     app.raw.hiring_manager_name ?? '',
      hiring_manager_linkedin: app.raw.hiring_manager_linkedin ?? ''
    };
    showEditModal = true;
  }

  async function saveEdit(e) {
    e.preventDefault();
    saving = true;
    try {
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

  const dossierEligible = $derived(app && ['screen', 'interview', 'offer'].includes(app.status));
  const dossierAvailable = $derived(!!dossier);
  function initialsOf(name) {
    return (name || '').split(/\s+/).filter(Boolean).slice(0, 2).map(s => s[0]).join('').toUpperCase();
  }
  const interviewerInitials = $derived(initialsOf(dossier?.content?.interviewer?.name));
  const hiringManagerInitials = $derived(initialsOf(app?.raw?.hiring_manager_name));

  const timeline = $derived.by(() => app ? buildTimelineFromApplication(app.raw) : []);

  // Stage index for the "Current stage" stat card (1 of 4).
  const STAGE_ORDER = ['applied', 'screen', 'interview', 'offer'];
  const stageIdx = $derived(app ? STAGE_ORDER.indexOf(app.status) : -1);

  // Days in pipeline: days since applied_at.
  const daysInPipe = $derived(app ? daysSince(app.raw.applied_at) : null);

  const appliedLong = $derived(app ? fmtLongDate(app.raw.applied_at) : '');
  const jdHost = $derived.by(() => {
    if (!app?.raw?.jd_url) return '';
    try { return new URL(app.raw.jd_url).hostname.replace(/^www\./, ''); } catch { return ''; }
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

      <!-- HERO STRIP -->
      <div class="hero">
        <div class="hero-top">
          {#if app.logoSrc}
            <img class="logo-big" src={app.logoSrc} alt={app.co} />
          {:else}
            <span class={`logo-big letter ${app.logoCls}`}>{app.coShort}</span>
          {/if}
          <div class="hero-text">
            <div class="co-row">
              <h1>{app.co}</h1>
              <span class={`pill ${app.status}`}><span class="pdot"></span>{STATUS_LABEL[app.status]}</span>
            </div>
            <div class="role-line">{app.role}</div>
          </div>
          {#if app.raw.jd_url}
            <a class="src-link" href={app.raw.jd_url} target="_blank" rel="noopener">
              <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.6"><path d="M5 11l6-6M6 5h5v5"/></svg>
              View on {jdHost || 'source'}
            </a>
          {/if}
        </div>
        <div class="facts">
          {#if app.raw.location}
            <span class="fact"><span class="fdot d-loc"></span>{app.raw.location}</span>
          {/if}
          {#if app.raw.applied_at}
            <span class="fact"><span class="fdot d-app"></span>Applied {appliedLong}</span>
          {/if}
          {#if app.source && app.source !== '—'}
            <span class="fact"><span class="fdot d-src"></span>{app.source}</span>
          {/if}
          {#if app.cv && app.cv !== '—'}
            <span class="fact"><span class="fdot d-cv"></span>CV {app.cv}</span>
          {/if}
          {#if app.raw.salary_note}
            <span class="fact"><span class="fdot d-sal"></span>{app.raw.salary_note}</span>
          {/if}
        </div>
      </div>

      <!-- UP NEXT (only when we have a dossier with a meeting) -->
      {#if dossier?.meeting && (dossier.meeting.starts_at || dossier.meeting.when)}
        <div class="upnext">
          <div class="up-left">
            <span class="up-tag"><span class="up-pulse"></span>Up next</span>
            <h3>{meetingWhen(dossier.meeting)}</h3>
            <div class="up-meta">
              <span>{meetingDuration(dossier.meeting)}</span>
              {#if dossier.meeting.medium}<span class="dot">·</span><span>{dossier.meeting.medium}</span>{/if}
              {#if dossier.meeting.panel}<span class="dot">·</span><span>{dossier.meeting.panel}</span>{/if}
            </div>
          </div>
          <div class="up-right">
            <button class="btn-prep">Open prep ↓</button>
          </div>
        </div>
      {/if}

      <!-- STATS -->
      <div class="stats">
        <div class="stat tone-accent">
          <span class="ribbon"></span>
          <div class="stat-lbl">Days in pipeline</div>
          <div class="stat-n">{daysInPipe ?? '—'}</div>
          <div class="stat-sub">{app.raw.applied_at ? `Applied ${appliedLong}` : 'Not yet applied'}</div>
        </div>
        <div class="stat tone-positive">
          <span class="ribbon"></span>
          <div class="stat-lbl">Current stage</div>
          <div class="stat-n">
            {stageIdx >= 0 ? stageIdx + 1 : '—'}<span class="of">/ 4</span>
          </div>
          <div class="stat-sub">{STATUS_LABEL[app.status]}</div>
        </div>
        <div class="stat tone-warm">
          <span class="ribbon"></span>
          <div class="stat-lbl">Match score</div>
          <div class="stat-n">—</div>
          <div class="stat-sub">{app.cv && app.cv !== '—' ? `CV ${app.cv} vs. JD` : 'Tag a CV variant to score'}</div>
        </div>
      </div>

      <!-- TABS -->
      <div class="tabs">
        <button class={`tab ${tab === 'brief' ? 'active' : ''}`} onclick={() => (tab = 'brief')}>
          Brief <span class="t-tag">AI</span>
        </button>
        <button class={`tab ${tab === 'interviews' ? 'active' : ''}`} onclick={() => (tab = 'interviews')}>
          Interviews <span class="t-tag">{interviews.length}</span>
        </button>
        <button class={`tab ${tab === 'timeline' ? 'active' : ''}`} onclick={() => (tab = 'timeline')}>
          Timeline <span class="t-tag">{timeline.length}</span>
        </button>
        <button class={`tab ${tab === 'notes' ? 'active' : ''}`} onclick={() => (tab = 'notes')}>Notes</button>
        <button class={`tab ${tab === 'files' ? 'active' : ''}`} onclick={() => (tab = 'files')}>Files</button>
      </div>

      {#if tab === 'brief'}
        <!-- HIRING MANAGER — always rendered (when set), regardless of dossier state. -->
        {#if app.raw.hiring_manager_name || app.raw.hiring_manager_linkedin}
          <div class="block person-block">
            <div class="block-hd">
              <h2>Hiring manager</h2>
              <span class="ai-tag-mute">from posting</span>
            </div>
            <div class="person">
              <div class="p-av t-warm">{hiringManagerInitials || '—'}</div>
              <div class="p-info">
                <h4>{app.raw.hiring_manager_name || 'Name not set'}</h4>
                {#if app.raw.role}<div class="p-role">{app.role}</div>{/if}
              </div>
              {#if app.raw.hiring_manager_linkedin}
                <a class="p-li" href={app.raw.hiring_manager_linkedin} target="_blank" rel="noopener">
                  <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M3.5 6h2v6h-2zM4.5 3a1 1 0 1 1 0 2 1 1 0 0 1 0-2zM7 6h2v.9c.3-.5.9-1 1.8-1 1.6 0 2.2 1 2.2 2.6V12h-2V9c0-.9-.3-1.4-1.1-1.4-.6 0-1 .4-1 1.2V12H7z"/></svg>
                  LinkedIn
                </a>
              {/if}
            </div>
          </div>
        {/if}
        {#if dossierAvailable}
          <!-- ABOUT COMPANY — compact card above the interviewer block. -->
          {#if dossier.content.company}
            {@const c = dossier.content.company}
            <div class="block company-block">
              <div class="block-hd">
                <h2>About {app.co}</h2>
                <span class="ai-tag">AI · refreshed {dossier.generatedAgo}</span>
              </div>
              {#if c.blurb}<p class="company-blurb">{c.blurb}</p>{/if}
              {#if c.direction}<p class="company-direction">{c.direction}</p>{/if}
              {#if c.hq || c.employees || c.stage || c.founded}
                <div class="facts-grid">
                  {#if c.hq}       <div class="f-cell"><div class="f-lbl">HQ</div><div class="f-val">{c.hq}</div></div>{/if}
                  {#if c.employees}<div class="f-cell"><div class="f-lbl">Employees</div><div class="f-val">{c.employees}</div></div>{/if}
                  {#if c.stage}    <div class="f-cell"><div class="f-lbl">Stage</div><div class="f-val">{c.stage}</div></div>{/if}
                  {#if c.founded}  <div class="f-cell"><div class="f-lbl">Founded</div><div class="f-val">{c.founded}</div></div>{/if}
                </div>
              {/if}
              {#if c.process?.length}
                <div class="sub-hd">
                  <h3>Typical interview process</h3>
                  <span class="ai-tag-mute">from Glassdoor / Blind / levels.fyi</span>
                </div>
                <ol class="process-list">
                  {#each c.process as p, i}
                    <li>
                      <span class="step-n">{i + 1}</span>
                      <div>
                        <div class="step-kind">{p.kind}</div>
                        {#if p.detail}<div class="step-detail">{p.detail}</div>{/if}
                      </div>
                    </li>
                  {/each}
                </ol>
              {/if}
              {#if c.watch_fors?.length}
                <div class="sub-hd">
                  <h3>Watch-fors for this loop</h3>
                  <span class="ai-tag-mute">AI · synthesized</span>
                </div>
                <ul class="watch-list">
                  {#each c.watch_fors as w}
                    <li><span class="w-dot">▸</span><span>{w}</span></li>
                  {/each}
                </ul>
              {/if}
            </div>
          {/if}

          <!-- INTERVIEWER -->
          {#if dossier.content.interviewer?.name}
            <div class="block person-block">
              <div class="block-hd">
                <h2>Your interviewer</h2>
                <span class="ai-tag">AI · refreshed {dossier.generatedAgo}</span>
                <button class="regen" onclick={generateDossier} disabled={generating} title="Regenerate">
                  <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M2 6a4 4 0 1 1 1.2 2.8M2 4v2h2"/></svg>
                  {generating ? '…' : 'Refresh'}
                </button>
              </div>
              <div class="person">
                <div class="p-av t-accent">{interviewerInitials}</div>
                <div class="p-info">
                  <h4>{dossier.content.interviewer.name}</h4>
                  {#if dossier.content.interviewer.role}
                    <div class="p-role">{dossier.content.interviewer.role}</div>
                  {/if}
                </div>
                {#each dossier.content.interviewer.links ?? [] as l}
                  <a class="p-li" href={l.href} target="_blank" rel="noopener">
                    <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M3.5 6h2v6h-2zM4.5 3a1 1 0 1 1 0 2 1 1 0 0 1 0-2zM7 6h2v.9c.3-.5.9-1 1.8-1 1.6 0 2.2 1 2.2 2.6V12h-2V9c0-.9-.3-1.4-1.1-1.4-.6 0-1 .4-1 1.2V12H7z"/></svg>
                    {l.label}
                  </a>
                {/each}
              </div>
              {#if dossier.content.interviewer.prior?.length}
                <div class="prior-row">
                  <span class="p-lbl">Prior</span>
                  {#each dossier.content.interviewer.prior as p}<span class="prior-chip">{p}</span>{/each}
                </div>
              {/if}
            </div>
          {/if}

          <!-- SNAPSHOT -->
          {#if dossier.content.snapshot}
            <div class="snapshot-card">
              <div class="snap-lbl">In one line</div>
              <p>{@html dossier.content.snapshot}</p>
            </div>
          {/if}

          <!-- BACKGROUND -->
          {#if dossier.content.background}
            <div class="block">
              <div class="block-hd">
                <h2>Background</h2>
                <span class="ai-tag">AI · web research</span>
              </div>
              <p class="prose">{dossier.content.background}</p>
            </div>
          {/if}

          <!-- RECENT POSTS & TALKS -->
          {#if dossier.content.signals?.length}
            <div class="block">
              <div class="block-hd">
                <h2>Recent posts &amp; talks</h2>
                <span class="ai-tag">last 90 days</span>
              </div>
              <div class="signals-row">
                {#each dossier.content.signals as s}
                  <div class="signal">
                    {#if s.source}
                      {@const sHost = (() => { try { return new URL(s.source.startsWith('http') ? s.source : `https://${s.source}`).hostname.replace(/^www\./,''); } catch { return s.source; } })()}
                      <img class="sig-logo" src={`https://www.google.com/s2/favicons?sz=128&domain=${sHost}`} alt="" />
                    {/if}
                    <div class="sig-meta">
                      {#if s.kind}<span class="sig-kind">{s.kind}</span>{/if}
                      {#if s.date}<span class="sig-date">{s.date}</span>{/if}
                    </div>
                    <div class="sig-body">{s.body}</div>
                    {#if s.source}<div class="sig-src">{s.source}</div>{/if}
                  </div>
                {/each}
              </div>
            </div>
          {/if}

          <!-- HOW TO APPROACH -->
          {#if dossier.content.lands?.length || dossier.content.avoid?.length}
            <div class="block">
              <div class="block-hd">
                <h2>How to approach this interview</h2>
                <span class="ai-tag">AI · interviewer-specific</span>
              </div>
              <div class="approach-grid">
                <div class="approach-col">
                  <div class="approach-hd">
                    <span class="approach-glyph ok">
                      <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2.5 6.5l2.5 2.5 4.5-5"/></svg>
                    </span>
                    Lead with
                  </div>
                  <ul class="approach-list">{#each dossier.content.lands ?? [] as l}
                    <li>
                      <span class="approach-marker ok">
                        <svg width="10" height="10" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2.5 6.5l2.5 2.5 4.5-5"/></svg>
                      </span>
                      <span>{l}</span>
                    </li>
                  {/each}</ul>
                </div>
                <div class="approach-col">
                  <div class="approach-hd">
                    <span class="approach-glyph no">
                      <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M3 3l6 6M9 3l-6 6"/></svg>
                    </span>
                    Steer clear of
                  </div>
                  <ul class="approach-list">{#each dossier.content.avoid ?? [] as a}
                    <li>
                      <span class="approach-marker no">
                        <svg width="10" height="10" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M3 3l6 6M9 3l-6 6"/></svg>
                      </span>
                      <span>{a}</span>
                    </li>
                  {/each}</ul>
                </div>
              </div>
            </div>
          {/if}

          <!-- QUESTIONS -->
          {#if dossier.content.questions?.length}
            <div class="block">
              <div class="block-hd">
                <h2>Questions worth asking</h2>
                <span class="ai-tag">ranked by signal</span>
              </div>
              <ol class="q-list">
                {#each dossier.content.questions as q, i}
                  <li>
                    <span class="qn">{i + 1}</span>
                    <div>
                      <div class="q">"{q.q}"</div>
                      {#if q.why}<div class="why">{q.why}</div>{/if}
                    </div>
                  </li>
                {/each}
              </ol>
            </div>
          {/if}

          <div class="disclaimer">
            Synthesised from public posts, talks, and papers · refreshed {dossier.generatedAgo} · always verify before you walk in
          </div>
        {:else if generating}
          <div class="generating-card">
            <div class="big-spinner"></div>
            <h3>Researching {app.co}{interviewerInput ? ` & ${interviewerInput}` : ''}…</h3>
            <p>Claude is searching the web for recent posts, talks, and the company's current direction. Typically 30–60 seconds.</p>
          </div>
        {:else}
          <div class="generate-card">
            <h3>Generate the brief</h3>
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
              <p class="muted-note">Move this application to <b>Screen</b>, <b>Interview</b>, or <b>Offer</b> to enable the brief.</p>
            {/if}
            {#if dossierError}
              <p class="dossier-err">{dossierError}</p>
            {/if}
          </div>
        {/if}
      {:else if tab === 'interviews'}
        <div class="add-card">
          <div class="add-hd">
            <h3>Add an interview</h3>
            <p>Drop a calendar file, paste a screenshot, or just paste the email body — we'll extract the event.</p>
          </div>
          <div class="zones">
            <!-- LEFT — .ics file + raw text -->
            <div class="zone">
              <div class="zone-hd">
                <span class="zone-ic">
                  <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="2" y="3" width="12" height="11" rx="1.5"/><path d="M2 6h12M6 2v2M10 2v2"/></svg>
                </span>
                <div>
                  <div class="zone-title">Calendar file</div>
                  <div class="zone-sub">.ics from Google / Outlook / Apple Calendar</div>
                </div>
              </div>
              <label class="drop drop-file">
                <svg width="20" height="20" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M8 11V3M5 6l3-3 3 3M3 11v2h10v-2"/></svg>
                <span class="drop-l1">Drop .ics file</span>
                <span class="drop-l2">or click to browse</span>
                <input type="file" accept=".ics,text/calendar" onchange={onIcsFile} style="display:none" />
              </label>
              <div class="or">or paste raw .ics text below</div>
              <textarea class="ai-ta" rows="3" placeholder={"BEGIN:VCALENDAR&#10;VERSION:2.0&#10;BEGIN:VEVENT…"} bind:value={icsText}></textarea>
              <button class="btn btn-primary zone-parse" onclick={parseIcs} disabled={icsParsing || !icsText.trim()}>
                {icsParsing ? 'Parsing…' : 'Parse .ics'}
              </button>
            </div>

            <!-- RIGHT — screenshot + email text -->
            <div class="zone">
              <div class="zone-hd">
                <span class="zone-ic accent">
                  <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="2" y="3" width="12" height="10" rx="1.5"/><circle cx="5.5" cy="6.5" r="1"/><path d="M2 11l3.5-3.5 3 3 2-2L14 11"/></svg>
                </span>
                <div>
                  <div class="zone-title">Screenshot or email text<span class="ai-pill">AI</span></div>
                  <div class="zone-sub">Gmail invite, Calendar screenshot, anything readable</div>
                </div>
              </div>
              <label
                class="drop drop-image"
                class:drag={aiDragOver}
                ondragover={onAiDragOver}
                ondragleave={() => (aiDragOver = false)}
                ondrop={onAiDrop}
              >
                {#if aiImage}
                  <svg width="20" height="20" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="2" y="3" width="12" height="10" rx="1.5"/><circle cx="5.5" cy="6.5" r="1"/><path d="M2 11l3.5-3.5 3 3 2-2L14 11"/></svg>
                  <span class="drop-l1">{aiImage.name}</span>
                  <span class="drop-l2">{Math.round(aiImage.size / 1024)} KB · click to replace</span>
                {:else}
                  <svg width="20" height="20" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="2" y="3" width="12" height="10" rx="1.5"/><circle cx="5.5" cy="6.5" r="1"/><path d="M2 11l3.5-3.5 3 3 2-2L14 11"/></svg>
                  <span class="drop-l1">Drop a screenshot</span>
                  <span class="drop-l2"><kbd>⌘V</kbd> works too</span>
                {/if}
                <input type="file" accept="image/png,image/jpeg,image/gif,image/webp" onchange={onAiFile} style="display:none" />
              </label>
              <div class="or">or paste the email body below</div>
              <textarea
                class="ai-ta"
                rows="3"
                placeholder={"You're invited to: Stripe — Technical screen&#10;When: Tue, May 28, 2:00 PM EDT&#10;Where: Google Meet"}
                bind:value={aiText}
                onpaste={onAiPaste}
              ></textarea>
              <button class="btn btn-primary zone-parse" onclick={parseAi} disabled={icsParsing || (!aiImage && !aiText.trim())}>
                {icsParsing ? 'Parsing…' : 'Parse with AI'}
              </button>
            </div>
          </div>

          {#if icsParseError}<p class="dossier-err" style="margin-top: 14px">{icsParseError}</p>{/if}

          {#if icsPreview.length > 0}
            <div class="ics-preview">
              <h4>Preview</h4>
              {#each icsPreview as ev}
                <div class="prev-row">
                  <div class="prev-summary">{ev.summary || 'Untitled event'}</div>
                  <div class="prev-when">{fmtEventWhen(ev)}</div>
                  {#if ev.location}<div class="prev-loc">📍 {ev.location}</div>{/if}
                </div>
              {/each}
              <button class="btn btn-primary" onclick={saveParsedEvents} disabled={icsSaving}>
                {icsSaving ? 'Saving…' : `Save ${icsPreview.length} event${icsPreview.length === 1 ? '' : 's'}`}
              </button>
            </div>
          {/if}
        </div>

        <div class="iv-list">
          <h3>Scheduled</h3>
          {#if interviewsLoading}
            <p style="color:var(--mute); font-size:13px;">Loading…</p>
          {:else if interviews.length === 0}
            <div class="empty-tab">
              <h3>No interviews on file</h3>
              <p>Once you add one above, it appears here and the brief picks it up.</p>
            </div>
          {:else}
            {#each interviews as iv (iv.id)}
              <div class="iv-card" class:past={isPast(iv)}>
                <div class="iv-card-main">
                  <div class="iv-when">{fmtEventWhen(iv)}</div>
                  <div class="iv-summary">{iv.summary || 'Untitled event'}</div>
                  {#if iv.location}<div class="iv-loc">📍 {iv.location}</div>{/if}
                  {#if iv.attendees}<div class="iv-att">👥 {iv.attendees}</div>{/if}
                </div>
                <button class="btn btn-danger" onclick={() => deleteInterview(iv)} title="Delete">Delete</button>
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
            <p>Events appear here as the application progresses.</p>
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
        <label>Hiring manager <input bind:value={edit.hiring_manager_name} placeholder="Jane Doe" /></label>
        <label>Hiring manager LinkedIn <input bind:value={edit.hiring_manager_linkedin} placeholder="https://linkedin.com/in/…" /></label>
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
  .body { padding: 28px; }
  .body-inner { max-width: 1100px; margin: 0 auto; }

  /* HERO */
  .hero {
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 18px;
    padding: 22px 24px;
    margin-bottom: 14px;
    box-shadow: var(--sh-1);
  }
  .hero-top { display: grid; grid-template-columns: 64px 1fr auto; gap: 18px; align-items: center; }
  .logo-big { width: 64px; height: 64px; border-radius: 14px; background: var(--card); object-fit: contain; padding: 8px; border: 1px solid var(--rule); }
  .logo-big.letter {
    display: grid; place-items: center;
    padding: 0;
    color: var(--ink); font-size: 24px; font-weight: 600;
    background: var(--surface-2);
  }
  .co-row { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; }
  .co-row h1 { font-size: 28px; font-weight: 600; margin: 0; letter-spacing: -0.025em; }
  .role-line { font-size: 15px; color: var(--ink-2); margin-top: 2px; font-weight: 500; }
  .src-link { display: inline-flex; align-items: center; gap: 6px; background: var(--card); border: 1px solid var(--rule); border-radius: 99px; padding: 8px 14px; font-size: 12.5px; font-weight: 500; color: var(--ink); text-decoration: none; transition: border-color 120ms, color 120ms; }
  .src-link:hover { border-color: var(--accent); color: var(--accent-text); }
  .facts { display: flex; flex-wrap: wrap; gap: 6px; margin-top: 14px; padding-top: 14px; border-top: 1px dashed var(--rule); }
  .fact { display: inline-flex; align-items: center; gap: 6px; font-size: 12.5px; color: var(--ink-2); background: var(--card); border: 1px solid var(--rule); padding: 4px 10px; border-radius: 99px; font-weight: 500; }
  .fdot { width: 6px; height: 6px; border-radius: 50%; background: var(--mute-2); }
  .d-loc { background: var(--accent); }
  .d-app { background: var(--positive); }
  .d-src { background: var(--warm); }
  .d-cv  { background: var(--mute-2); }
  .d-sal { background: var(--ink); }

  /* Pills (page-scoped so we don't fight design-system) */
  .pill { display: inline-flex; align-items: center; gap: 5px; padding: 4px 10px; border-radius: 99px; font-size: 12px; font-weight: 500; background: var(--surface-2); color: var(--ink-2); width: max-content; }
  .pill .pdot { width: 6px; height: 6px; border-radius: 50%; background: var(--mute-2); }
  .pill.wishlist { background: var(--surface-2); color: var(--mute); }
  .pill.applied { background: var(--surface-2); color: var(--ink-2); }
  .pill.screen { background: var(--positive-tint); color: var(--positive-text); }
  .pill.screen .pdot { background: var(--positive); }
  .pill.interview { background: var(--accent-tint); color: var(--accent-text); }
  .pill.interview .pdot { background: var(--accent); }
  .pill.offer { background: var(--warm-tint); color: var(--warm-text); }
  .pill.offer .pdot { background: var(--warm); }
  .pill.rejected, .pill.withdrawn { background: var(--surface-2); color: var(--mute); }

  /* UP NEXT */
  .upnext {
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 14px;
    padding: 18px 22px;
    margin-bottom: 14px;
    display: grid; grid-template-columns: 1fr auto; gap: 16px; align-items: center;
    box-shadow: var(--sh-1);
    position: relative;
    overflow: hidden;
  }
  .upnext::before { content: ''; position: absolute; left: 0; top: 0; bottom: 0; width: 4px; background: var(--accent); }
  .up-tag { display: inline-flex; align-items: center; gap: 6px; font-size: 12px; color: var(--accent-text); background: var(--accent-tint); padding: 4px 10px; border-radius: 99px; font-weight: 500; margin-bottom: 8px; }
  @keyframes up-pulse {
    0%, 100% { box-shadow: 0 0 0 0 var(--accent); }
    50%      { box-shadow: 0 0 0 5px transparent; }
  }
  .up-pulse { width: 6px; height: 6px; border-radius: 50%; background: var(--accent); animation: up-pulse 1.6s ease-in-out infinite; }
  .upnext h3 { font-size: 18px; font-weight: 600; margin: 0; letter-spacing: -0.015em; }
  .up-meta { margin-top: 4px; font-size: 13px; color: var(--mute); display: flex; flex-wrap: wrap; gap: 0 6px; }
  .up-meta .dot { color: var(--mute-2); }
  .btn-prep {
    background: var(--accent); color: white; border: 0; border-radius: 99px;
    padding: 10px 18px; font-size: 13.5px; font-weight: 600; cursor: pointer;
    transition: transform 120ms ease;
  }
  .btn-prep:hover { transform: translateY(-1px); }

  /* STATS */
  .stats { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; margin-bottom: 24px; }
  .stat { position: relative; background: var(--card); border: 1px solid var(--rule); border-radius: 14px; padding: 16px 18px; overflow: hidden; box-shadow: var(--sh-1); }
  .stat .ribbon { position: absolute; top: 0; left: 0; right: 0; height: 3px; }
  .stat.tone-accent   .ribbon { background: var(--accent); }
  .stat.tone-positive .ribbon { background: var(--positive); }
  .stat.tone-warm     .ribbon { background: var(--warm); }
  .stat-lbl { font-size: 12.5px; color: var(--mute); margin-bottom: 4px; font-weight: 500; }
  .stat-n { font-size: 32px; font-weight: 600; letter-spacing: -0.035em; line-height: 1.1; font-feature-settings: "tnum"; }
  .stat.tone-accent   .stat-n { color: var(--accent-text); }
  .stat.tone-positive .stat-n { color: var(--positive-text); }
  .stat.tone-warm     .stat-n { color: var(--warm-text); }
  .of { font-size: 16px; color: var(--mute); margin-left: 6px; font-weight: 500; }
  .stat-sub { font-size: 12px; color: var(--mute); margin-top: 6px; padding-top: 6px; border-top: 1px dashed var(--rule); }

  /* TABS */
  .tabs { display: flex; gap: 4px; border-bottom: 1px solid var(--rule); margin-bottom: 18px; }
  .tab { background: transparent; border: 0; padding: 10px 14px; font-size: 13.5px; color: var(--mute); cursor: pointer; border-bottom: 2px solid transparent; margin-bottom: -1px; font-weight: 600; }
  .tab.active { color: var(--ink); border-bottom-color: var(--ink); }
  .t-tag { font-size: 11px; background: var(--accent-tint); color: var(--accent-text); padding: 1px 7px; border-radius: 99px; margin-left: 4px; font-weight: 500; }

  /* BLOCK */
  .block { background: var(--card); border: 1px solid var(--rule); border-radius: 14px; padding: 20px 22px; margin-bottom: 14px; box-shadow: var(--sh-1); }
  .block-hd { display: flex; align-items: center; gap: 10px; margin-bottom: 14px; }
  .block-hd h2 { font-size: 16px; font-weight: 600; margin: 0; letter-spacing: -0.015em; }
  .ai-tag { display: inline-flex; align-items: center; gap: 5px; font-size: 12px; background: var(--accent-tint); color: var(--accent-text); padding: 3px 10px; border-radius: 99px; font-weight: 500; }
  .ai-tag-mute { display: inline-flex; align-items: center; gap: 5px; font-size: 12px; background: var(--surface-2); color: var(--mute); padding: 3px 10px; border-radius: 99px; font-weight: 500; }
  .regen { margin-left: auto; display: inline-flex; align-items: center; gap: 6px; background: transparent; border: 1px solid var(--rule); border-radius: 99px; padding: 4px 10px; font-size: 12px; color: var(--ink-2); cursor: pointer; }
  .regen:hover:not(:disabled) { background: var(--surface-2); }
  .regen:disabled { opacity: 0.55; cursor: wait; }
  .prose { margin: 0; font-size: 13.5px; line-height: 1.55; color: var(--ink-2); }

  /* PEOPLE */
  .person-block { padding: 18px 20px; }
  .person { display: grid; grid-template-columns: 46px 1fr auto; gap: 12px; align-items: center; }
  .p-av { width: 46px; height: 46px; border-radius: 50%; display: grid; place-items: center; font-weight: 600; font-size: 16px; }
  .p-av.t-accent { background: var(--accent-tint); color: var(--accent-text); }
  .p-av.t-warm { background: var(--warm-tint); color: var(--warm-text); }
  .p-info h4 { margin: 0; font-size: 14.5px; font-weight: 600; letter-spacing: -0.01em; }
  .p-info .p-role { font-size: 12.5px; color: var(--mute); margin-top: 2px; }
  .p-li { display: inline-flex; align-items: center; gap: 6px; background: var(--surface-2); border: 1px solid var(--rule); border-radius: 99px; padding: 6px 12px; font-size: 12px; font-weight: 600; color: var(--ink); text-decoration: none; }
  .p-li svg { color: #0a66c2; }
  .prior-row { display: flex; align-items: center; gap: 8px; margin-top: 14px; padding-top: 12px; border-top: 1px dashed var(--rule); flex-wrap: wrap; }
  .p-lbl { font-size: 11.5px; color: var(--mute); font-weight: 600; }
  .prior-chip { font-size: 11.5px; background: var(--surface-2); color: var(--ink-2); padding: 3px 9px; border-radius: 99px; font-weight: 500; }

  /* SNAPSHOT */
  .snapshot-card {
    background: var(--accent-tint);
    border: 1px solid var(--rule);
    border-radius: 14px;
    padding: 16px 20px;
    margin: 14px 0;
  }
  .snap-lbl { font-size: 11.5px; color: var(--accent-text); font-weight: 600; margin-bottom: 4px; }

  /* COMPANY CARD — sits above the interviewer block on the Brief tab. */
  .company-block .company-blurb { font-size: 14.5px; line-height: 1.55; color: var(--ink); margin: 0 0 6px; font-weight: 500; }
  .company-block .company-direction { font-size: 13.5px; line-height: 1.6; color: var(--ink-2); margin: 0 0 14px; }
  .company-block .facts-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 8px; margin-bottom: 18px; }
  .company-block .f-cell { background: var(--surface-2); border-radius: 10px; padding: 10px 12px; }
  .company-block .f-lbl { font-size: 11px; color: var(--mute); font-weight: 500; }
  .company-block .f-val { font-size: 13px; color: var(--ink); font-weight: 600; margin-top: 2px; }
  .company-block .sub-hd { display: flex; align-items: center; gap: 10px; margin: 18px 0 10px; padding-top: 14px; border-top: 1px dashed var(--rule); }
  .company-block .sub-hd h3 { font-size: 13.5px; font-weight: 600; margin: 0; letter-spacing: -0.005em; color: var(--ink); }
  .company-block .process-list { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 6px; }
  .company-block .process-list li { display: grid; grid-template-columns: 24px 1fr; gap: 10px; align-items: baseline; padding: 8px 0; border-bottom: 1px dashed var(--rule); }
  .company-block .process-list li:last-child { border-bottom: 0; }
  .company-block .step-n { width: 22px; height: 22px; border-radius: 50%; background: var(--accent-tint); color: var(--accent-text); display: grid; place-items: center; font-size: 11px; font-weight: 600; font-feature-settings: "tnum"; }
  .company-block .step-kind { font-size: 13.5px; color: var(--ink); font-weight: 600; }
  .company-block .step-detail { font-size: 12.5px; color: var(--mute); margin-top: 2px; }
  .company-block .watch-list { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 6px; }
  .company-block .watch-list li { display: grid; grid-template-columns: 18px 1fr; gap: 8px; align-items: baseline; padding: 6px 0; font-size: 13.5px; color: var(--ink-2); line-height: 1.5; }
  .company-block .w-dot { color: var(--accent); font-weight: 700; }

  @media (max-width: 720px) {
    .company-block .facts-grid { grid-template-columns: repeat(2, 1fr); }
  }
  .snapshot-card p { margin: 0; font-size: 14.5px; line-height: 1.55; color: var(--ink); }

  /* SIGNALS */
  .signals-row { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 10px; }
  .signal {
    background: var(--surface-2);
    border: 1px solid var(--rule);
    border-radius: 12px;
    padding: 12px 14px;
    display: grid;
    grid-template-columns: 28px 1fr;
    grid-template-rows: auto auto auto;
    grid-template-areas:
      "logo meta"
      "body body"
      "src  src";
    gap: 6px 10px;
  }
  .sig-logo { grid-area: logo; width: 28px; height: 28px; border-radius: 7px; background: var(--card); object-fit: contain; padding: 3px; border: 1px solid var(--rule); }
  .sig-meta { grid-area: meta; display: flex; align-items: center; gap: 8px; }
  .sig-kind { font-size: 11px; color: var(--mute); font-weight: 600; }
  .sig-date { font-size: 11px; color: var(--mute-2); }
  .sig-body { grid-area: body; font-size: 13px; color: var(--ink); line-height: 1.4; }
  .sig-src  { grid-area: src; font-size: 11.5px; color: var(--accent-text); }

  /* APPROACH */
  .approach-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 22px; }
  .approach-hd { display: flex; align-items: center; gap: 8px; font-size: 13px; font-weight: 600; margin-bottom: 10px; color: var(--ink); }
  .approach-glyph { width: 22px; height: 22px; border-radius: 50%; display: grid; place-items: center; }
  .approach-glyph.ok { background: var(--positive-tint); color: var(--positive-text); }
  .approach-glyph.no { background: var(--danger-tint); color: var(--danger-text); }
  .approach-list { margin: 0; padding: 0; list-style: none; display: flex; flex-direction: column; gap: 8px; }
  .approach-list li { display: grid; grid-template-columns: 18px 1fr; gap: 10px; align-items: flex-start; font-size: 13px; color: var(--ink-2); line-height: 1.5; }
  .approach-marker { width: 18px; height: 18px; border-radius: 50%; display: grid; place-items: center; margin-top: 1px; }
  .approach-marker.ok { background: var(--positive-tint); color: var(--positive-text); }
  .approach-marker.no { background: var(--danger-tint); color: var(--danger-text); }

  /* QUESTIONS */
  .q-list { margin: 0; padding: 0; list-style: none; display: flex; flex-direction: column; gap: 8px; }
  .q-list li { background: var(--surface-2); border-radius: 12px; padding: 12px 14px; display: grid; grid-template-columns: 28px 1fr; gap: 12px; align-items: center; }
  .qn { width: 24px; height: 24px; border-radius: 50%; background: var(--accent-tint); color: var(--accent-text); display: grid; place-items: center; font-size: 12px; font-weight: 600; font-feature-settings: "tnum"; }
  .q { font-size: 13.5px; color: var(--ink); font-weight: 500; }
  .why { font-size: 11.5px; color: var(--mute); margin-top: 3px; }

  .disclaimer { margin: 18px 0 0; font-size: 11.5px; color: var(--mute); padding-top: 14px; border-top: 1px dashed var(--rule); }

  /* GENERATE / GENERATING */
  .generate-card {
    border: 1px dashed var(--rule-strong);
    border-radius: 14px;
    padding: 28px 32px;
    background: var(--card);
    text-align: center;
  }
  .generate-card h3 { font-size: 17px; font-weight: 600; letter-spacing: -0.012em; margin: 0 0 .35rem; color: var(--ink); }
  .generate-card > p { margin: 0 0 1.25rem; color: var(--mute); font-size: 13.5px; line-height: 1.55; max-width: 56ch; margin-left: auto; margin-right: auto; }
  .generate-row { display: flex; gap: 8px; max-width: 520px; margin: 0 auto; }
  .generate-row input { flex: 1; font: inherit; font-size: 13.5px; color: var(--ink); background: var(--surface); border: 1px solid var(--rule); border-radius: 8px; padding: 6px 10px; outline: none; }
  .generate-row input:focus { border-color: var(--accent); }
  .muted-note { color: var(--mute); font-size: 12px; margin: .75rem 0 0; }
  .muted-note b { color: var(--ink-2); font-weight: 600; }
  .dossier-err { color: var(--danger-text); background: var(--danger-tint); border: 1px solid var(--danger-tint); border-radius: 8px; padding: 8px 12px; font-size: 13px; margin: .75rem 0 0; }
  .generating-card {
    border: 1px solid var(--accent-tint-2);
    background: var(--accent-tint);
    border-radius: 14px;
    padding: 32px;
    text-align: center;
  }
  .generating-card h3 { font-size: 17px; font-weight: 600; letter-spacing: -0.012em; margin: 16px 0 .5rem; color: var(--ink); }
  .generating-card p { color: var(--mute); font-size: 13.5px; margin: 0; max-width: 56ch; margin-left: auto; margin-right: auto; }
  .big-spinner { width: 24px; height: 24px; margin: 0 auto; border: 2px solid var(--accent-tint-2); border-top-color: var(--accent); border-radius: 50%; animation: spin 800ms linear infinite; }
  @keyframes spin { to { transform: rotate(360deg); } }

  /* TIMELINE (kept from old) */
  .timeline { display: flex; flex-direction: column; gap: 0; }
  .timeline-event { display: grid; grid-template-columns: 80px 16px 1fr; gap: 14px; padding: 12px 0; align-items: flex-start; border-bottom: 1px solid var(--rule); }
  .timeline-event .date { font-size: 12px; color: var(--mute); padding-top: 2px; }
  .timeline-event .axis { position: relative; height: 100%; display: flex; justify-content: center; }
  .timeline-event .marker { width: 8px; height: 8px; border-radius: 50%; background: var(--accent); margin-top: 6px; }
  .timeline-event .label { font-size: 13.5px; font-weight: 500; color: var(--ink); }
  .timeline-event .note { font-size: 12.5px; color: var(--mute); margin-top: 2px; }

  /* STATUS MENU */
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
  .status-menu-item { display: flex; align-items: center; gap: 8px; padding: 6px 8px; border-radius: 5px; background: transparent; font: inherit; font-size: 13px; color: var(--ink-2); cursor: pointer; text-align: left; width: 100%; border: 0; }
  .status-menu-item:hover { background: var(--surface-2); }
  .status-menu-item.current { background: var(--surface-2); }
  .status-menu-item .check { margin-left: auto; color: var(--accent-text); font-weight: 600; font-size: 12px; }
  .menu-scrim { position: fixed; inset: 0; z-index: 49; background: transparent; }

  .btn-danger { color: var(--danger-text); border-color: var(--rule); }
  .btn-danger:hover { background: var(--danger-tint); border-color: var(--danger-tint); }

  /* EDIT MODAL */
  .modal-overlay { position: fixed; inset: 0; background: rgba(10,10,13,0.4); display: grid; place-items: center; z-index: 100; padding: 2rem; }
  .modal { background: var(--card); border: 1px solid var(--rule); border-radius: 12px; padding: 1.5rem; width: 100%; max-width: 560px; display: flex; flex-direction: column; gap: .75rem; box-shadow: var(--sh-pop); }
  .modal h2 { font-size: 18px; font-weight: 600; letter-spacing: -0.018em; margin: 0; }
  .modal-hint { font-size: 12px; color: var(--mute); margin: 0 0 .5rem; }
  .fields { display: grid; grid-template-columns: 1fr 1fr; gap: .65rem; }
  .fields .span-2 { grid-column: span 2; }
  .modal label { display: flex; flex-direction: column; font-size: 12px; color: var(--mute); gap: .35rem; }
  .modal input { font: inherit; color: var(--ink); background: var(--surface); border: 1px solid var(--rule); border-radius: 6px; padding: .45rem .6rem; font-size: 13.5px; outline: none; transition: border-color 100ms ease; }
  .modal input:focus { border-color: var(--accent); }
  .modal-actions { display: flex; justify-content: flex-end; gap: .5rem; margin-top: .75rem; }

  .empty-tab {
    border: 1px dashed var(--rule);
    border-radius: 12px;
    padding: 32px;
    text-align: center;
    background: var(--card);
  }
  .empty-tab h3 { margin: 0 0 .5rem; font-size: 16px; font-weight: 500; color: var(--ink); }
  .empty-tab p { color: var(--mute); margin: 0; font-size: 13.5px; }

  /* Interviews tab — two side-by-side zones (ics + AI). */
  .add-card { background: var(--card); border: 1px solid var(--rule); border-radius: 14px; padding: 22px 24px; margin-bottom: 18px; box-shadow: var(--sh-1); }
  .add-hd { margin-bottom: 16px; }
  .add-hd h3 { font-size: 16px; font-weight: 600; margin: 0 0 4px; letter-spacing: -0.015em; }
  .add-hd p { font-size: 13.5px; color: var(--mute); margin: 0; line-height: 1.5; }
  .zones { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
  .zone { background: var(--surface); border: 1px solid var(--rule); border-radius: 12px; padding: 14px 16px; display: flex; flex-direction: column; gap: 10px; }
  .zone-hd { display: grid; grid-template-columns: 32px 1fr; gap: 10px; align-items: center; }
  .zone-ic { width: 32px; height: 32px; border-radius: 8px; background: var(--surface-2); color: var(--ink-2); display: grid; place-items: center; }
  .zone-ic.accent { background: var(--accent-tint); color: var(--accent-text); }
  .zone-title { font-size: 13.5px; font-weight: 600; color: var(--ink); display: inline-flex; align-items: center; gap: 6px; }
  .zone-sub { font-size: 12px; color: var(--mute); margin-top: 1px; }
  .ai-pill { font-size: 10px; font-weight: 600; color: var(--accent-text); background: var(--accent-tint); border-radius: 4px; padding: 1px 5px; letter-spacing: .04em; }
  .drop { background: var(--card); border: 1.5px dashed var(--rule-strong); border-radius: 10px; padding: 16px 12px; display: flex; flex-direction: column; align-items: center; gap: 3px; color: var(--mute); transition: border-color 120ms, background 120ms; cursor: pointer; }
  .drop:hover, .drop.drag { border-color: var(--accent); background: var(--accent-tint); color: var(--accent-text); }
  .drop-l1 { font-size: 12.5px; font-weight: 500; color: var(--ink-2); }
  .drop:hover .drop-l1, .drop.drag .drop-l1 { color: var(--accent-text); }
  .drop-l2 { font-size: 11.5px; color: var(--mute-2); }
  .drop kbd { font-family: var(--mono); font-size: 10px; background: var(--surface-2); border: 1px solid var(--rule); border-bottom-width: 2px; border-radius: 3px; padding: 0 4px; color: var(--ink-2); }
  .or { font-size: 11.5px; color: var(--mute-2); text-align: center; }
  .ai-ta { width: 100%; font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 11.5px; line-height: 1.5; color: var(--ink); background: var(--card); border: 1px solid var(--rule); border-radius: 8px; padding: 8px 10px; outline: none; resize: vertical; box-sizing: border-box; }
  .ai-ta:focus { border-color: var(--accent); }
  .zone-parse { margin-top: 4px; align-self: flex-start; }
  .ics-preview { margin-top: 18px; padding-top: 16px; border-top: 1px solid var(--rule); }
  .ics-preview h4 { font-size: 11.5px; font-weight: 600; color: var(--mute); text-transform: uppercase; letter-spacing: 0.04em; margin: 0 0 10px; }
  .prev-row { background: var(--accent-tint); border: 1px solid var(--accent); border-radius: 10px; padding: 12px 14px; margin-bottom: 10px; }
  .prev-summary { font-size: 13.5px; font-weight: 600; color: var(--ink); }
  .prev-when { font-size: 12.5px; color: var(--accent-text); margin-top: 3px; font-weight: 500; }
  .prev-loc { font-size: 12px; color: var(--mute); margin-top: 4px; }

  .iv-list h3 { font-size: 13px; font-weight: 500; color: var(--mute); text-transform: uppercase; letter-spacing: 0.04em; margin: 0 0 12px; }
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
  .iv-when { font-size: 12px; color: var(--accent-text); font-weight: 500; margin-bottom: 2px; }
  .iv-summary { font-size: 14px; color: var(--ink); margin-bottom: 2px; }
  .iv-loc, .iv-att { font-size: 12px; color: var(--mute); margin-top: 2px; }

  /* Mobile — stack hero and grids; keep tabs scrollable. */
  @media (max-width: 720px) {
    .body { padding: 18px 14px; }
    .hero { padding: 18px 16px; border-radius: 14px; }
    .hero-top { grid-template-columns: 48px 1fr; gap: 12px; }
    .logo-big { width: 48px; height: 48px; border-radius: 12px; padding: 6px; }
    .logo-big.letter { font-size: 18px; }
    .co-row h1 { font-size: 22px; }
    .role-line { font-size: 13.5px; }
    .src-link { grid-column: 1 / -1; justify-self: start; margin-top: 4px; padding: 6px 10px; font-size: 12px; }
    .facts { gap: 4px; }
    .fact { font-size: 11.5px; padding: 3px 8px; }
    .upnext { grid-template-columns: 1fr; padding: 14px 16px; }
    .upnext h3 { font-size: 16px; }
    .stats { grid-template-columns: 1fr; gap: 8px; margin-bottom: 18px; }
    .stat { padding: 14px 16px; }
    .stat-n { font-size: 26px; }
    .tabs { overflow-x: auto; flex-wrap: nowrap; }
    .tab { white-space: nowrap; padding: 10px 12px; font-size: 13px; }
    .block { padding: 16px 16px; }
    .block-hd { flex-wrap: wrap; }
    .person { grid-template-columns: 40px 1fr; row-gap: 8px; }
    .p-li { grid-column: 1 / -1; justify-self: start; }
    .signals-row { grid-template-columns: 1fr; }
    .approach-grid { grid-template-columns: 1fr; gap: 18px; }
    .q-list li { grid-template-columns: 24px 1fr; }
    .modal-overlay { padding: 0; }
    .modal { max-width: 100%; border-radius: 0; min-height: 100vh; padding: 1rem; }
    .fields { grid-template-columns: 1fr; }
    .fields .span-2 { grid-column: auto; }
    .add-card { padding: 16px 16px; }
    .zones { grid-template-columns: 1fr; gap: 10px; }
  }
</style>
