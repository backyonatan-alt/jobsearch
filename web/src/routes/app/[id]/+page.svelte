<script>
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api.js';
  import { isPreview, mockApi } from '$lib/preview-mode.js';
  import {
    toDisplayApp, STATUS_LABEL, STATUSES,
    fmtLongDate, fmtRelativeDate, daysSince, isStale
  } from '$lib/app-helpers.js';

  const call = isPreview() ? mockApi : api;

  let app = $state(null);
  let loading = $state(true);
  let notFound = $state(false);

  // Interviews
  let interviews = $state([]);
  let interviewsLoading = $state(false);

  // Add-event flow (popup modal)
  let showEventModal = $state(false);
  let icsText = $state('');
  let aiText = $state('');
  let aiImage = $state(null); // { name, mediaType, size, file }
  let aiDragOver = $state(false);
  let icsParsing = $state(false);
  let icsParseError = $state('');
  let icsPreview = $state([]);
  let icsSaving = $state(false);
  const AI_ALLOWED_IMG = ['image/png', 'image/jpeg', 'image/gif', 'image/webp'];
  const AI_MAX_BYTES = 6 * 1024 * 1024;

  // Inline-action state
  let showStatusMenu = $state(false);
  let showEditModal = $state(false);
  let edit = $state({ company: '', role: '', source: '', location: '', cv_variant: '', jd_url: '', salary_note: '', hiring_manager_name: '', hiring_manager_linkedin: '' });
  let saving = $state(false);

  // Follow-up (records what the user did — sends nothing). Loaded from backend.
  let followUps = $state([]);
  let followUpsLoading = $state(false);
  let toast = $state('');
  let toastTimer = null;

  // Follow-up popup modal
  let showFollowUpModal = $state(false);
  let fuNote = $state('');
  let fuChannel = $state('');
  let fuDate = $state('');
  let fuSaving = $state(false);
  const FU_CHANNELS = ['Email', 'LinkedIn', 'Phone', 'In person', 'Other'];

  const id = $derived(page.params.id);

  $effect(() => {
    void id;
    loadApp();
    loadInterviews();
    loadFollowUps();
  });

  async function loadApp() {
    loading = true;
    notFound = false;
    try {
      const raw = await call(`/api/applications/${id}`);
      app = toDisplayApp(raw);
    } catch (e) {
      if (e.message === 'unauthorized') return;
      if (e.message.includes('not found') || e.message.includes('404')) notFound = true;
      else console.error(e);
    } finally {
      loading = false;
    }
  }

  async function loadInterviews() {
    interviewsLoading = true;
    try {
      interviews = await call(`/api/applications/${id}/interviews`);
    } catch (e) {
      if (e.message !== 'unauthorized') console.error(e);
      interviews = [];
    } finally {
      interviewsLoading = false;
    }
  }

  async function loadFollowUps() {
    followUpsLoading = true;
    try {
      followUps = await call(`/api/applications/${id}/follow-ups`);
    } catch (e) {
      if (e.message !== 'unauthorized') console.error(e);
      followUps = [];
    } finally {
      followUpsLoading = false;
    }
  }

  // ── Add-event parse/save flow (preserved) ────────────────────
  function setAiImage(f) {
    if (!f) return;
    if (!AI_ALLOWED_IMG.includes(f.type)) { icsParseError = 'Only PNG / JPEG / GIF / WebP screenshots are supported.'; return; }
    if (f.size > AI_MAX_BYTES) { icsParseError = 'Screenshot too large (6 MB max).'; return; }
    icsParseError = '';
    aiImage = { name: f.name || 'pasted.png', mediaType: f.type, size: f.size, file: f };
  }
  function onAiDrop(e) { e.preventDefault(); aiDragOver = false; setAiImage(e.dataTransfer?.files?.[0]); }
  function onAiDragOver(e) { if (e.dataTransfer?.types?.includes('Files')) { e.preventDefault(); aiDragOver = true; } }
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
  async function onIcsFile(e) {
    const f = e.target.files?.[0];
    e.target.value = '';
    if (!f) return;
    if (f.size > 256 * 1024) { icsParseError = 'File too large (256 KB max).'; return; }
    icsText = await f.text();
    await parseIcs();
  }
  async function parseIcs() {
    if (icsParsing) return;
    icsParseError = '';
    icsPreview = [];
    const body = icsText.trim();
    if (!body) { icsParseError = 'Paste the .ics contents or pick a file.'; return; }
    icsParsing = true;
    try {
      const r = await call(`/api/applications/${id}/interviews/parse`, { method: 'POST', body: JSON.stringify({ ics: body }) });
      icsPreview = r.events ?? [];
      if (icsPreview.length === 0) icsParseError = 'No events found in that file.';
    } catch (e) {
      icsParseError = e.message || 'Could not parse calendar.';
    } finally {
      icsParsing = false;
    }
  }
  async function parseAi() {
    if (icsParsing) return;
    icsParseError = '';
    icsPreview = [];
    if (!aiImage && !aiText.trim()) { icsParseError = 'Drop a screenshot or paste the event text first.'; return; }
    icsParsing = true;
    try {
      const payload = {};
      if (aiText.trim()) payload.text = aiText.trim();
      if (aiImage) {
        const data = await fileToBase64(aiImage.file);
        payload.image = { media_type: aiImage.mediaType, data };
      }
      const r = await call(`/api/applications/${id}/interviews/parse`, { method: 'POST', body: JSON.stringify(payload) });
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
        await call(`/api/applications/${id}/interviews`, { method: 'POST', body: JSON.stringify(ev) });
      }
      icsText = ''; aiText = ''; aiImage = null;
      icsPreview = [];
      showEventModal = false;
      await loadInterviews();
    } catch (e) {
      icsParseError = e.message || 'Could not save events.';
    } finally {
      icsSaving = false;
    }
  }
  async function deleteInterview(iv) {
    if (!confirm(`Delete "${iv.summary}"?`)) return;
    await call(`/api/applications/${id}/interviews/${iv.id}`, { method: 'DELETE' });
    await loadInterviews();
  }

  function openAddEvent() {
    showEventModal = true;
    icsParseError = '';
    icsPreview = [];
  }
  function closeEventModal() {
    showEventModal = false;
  }
  function onWindowKeydown(e) {
    if (e.key !== 'Escape') return;
    if (showEventModal) closeEventModal();
    if (showFollowUpModal) closeFollowUp();
  }

  // ── Status / edit / delete (preserved) ───────────────────────
  async function setStatus(newStatus) {
    showStatusMenu = false;
    if (!app || newStatus === app.status) return;
    await call(`/api/applications/${id}`, { method: 'PATCH', body: JSON.stringify({ status: newStatus }) });
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
      await call(`/api/applications/${id}`, { method: 'PATCH', body: JSON.stringify(payload) });
      showEditModal = false;
      await loadApp();
    } finally {
      saving = false;
    }
  }
  async function deleteApp() {
    if (!app) return;
    if (!confirm(`Delete the ${app.co} application? This can't be undone.`)) return;
    await call(`/api/applications/${id}`, { method: 'DELETE' });
    goto('/app', { replaceState: true });
  }
  function back() { goto('/app'); }

  function openPlaybook() { goto(`/app/${id}/playbook`); }

  // ── Follow-up (records what the user did — sends nothing) ─────
  function todayInputValue() {
    const d = new Date();
    const off = d.getTimezoneOffset();
    return new Date(d.getTime() - off * 60000).toISOString().slice(0, 10);
  }
  function openFollowUp() {
    fuNote = '';
    fuChannel = '';
    fuDate = todayInputValue();
    showFollowUpModal = true;
  }
  function closeFollowUp() { showFollowUpModal = false; }
  async function saveFollowUp(e) {
    e.preventDefault();
    if (fuSaving) return;
    fuSaving = true;
    try {
      // Date input is a local YYYY-MM-DD; turn it into an RFC3339 timestamp.
      // Blank → now.
      const occurred_at = fuDate ? new Date(`${fuDate}T12:00:00`).toISOString() : new Date().toISOString();
      const payload = { note: fuNote.trim(), channel: fuChannel, occurred_at };
      await call(`/api/applications/${id}/follow-ups`, { method: 'POST', body: JSON.stringify(payload) });
      showFollowUpModal = false;
      await loadFollowUps();
      await loadApp();
      showToast("Follow-up logged — we've reset the clock");
    } finally {
      fuSaving = false;
    }
  }
  async function deleteFollowUp(f) {
    await call(`/api/applications/${id}/follow-ups/${f.id}`, { method: 'DELETE' });
    await loadFollowUps();
    await loadApp();
  }
  function showToast(msg) {
    toast = msg;
    if (toastTimer) clearTimeout(toastTimer);
    toastTimer = setTimeout(() => { toast = ''; }, 2600);
  }

  // ── Derived view data ────────────────────────────────────────
  function initialsOf(name) {
    return (name || '').split(/\s+/).filter(Boolean).slice(0, 2).map(s => s[0]).join('').toUpperCase();
  }
  const hiringManagerInitials = $derived(initialsOf(app?.raw?.hiring_manager_name));

  const appliedLong = $derived(app ? fmtLongDate(app.raw.applied_at) : '');

  // Next event = the soonest interview dated now-or-later.
  const upcoming = $derived.by(() => {
    const now = Date.now();
    const future = (interviews || [])
      .filter(iv => iv?.starts_at && new Date(iv.starts_at).getTime() >= now)
      .sort((a, b) => new Date(a.starts_at) - new Date(b.starts_at));
    return future[0] || null;
  });

  const awaiting = $derived(app && ['applied', 'screen'].includes(app.status));
  const quiet = $derived(app ? isStale(app.raw) : false);
  const lastActivityDays = $derived(app ? daysSince(app.raw.updated_at ?? app.raw.applied_at) : null);
  const waitDays = $derived(app ? daysSince(app.raw.applied_at) : null);

  function evWhen(iv) {
    if (!iv?.starts_at) return '';
    const d = new Date(iv.starts_at);
    if (iv.all_day) return d.toLocaleDateString(undefined, { weekday: 'short', month: 'short', day: 'numeric' }) + ' · all day';
    const now = new Date();
    const startOfDay = (x) => new Date(x.getFullYear(), x.getMonth(), x.getDate());
    const days = Math.round((startOfDay(d) - startOfDay(now)) / 86400000);
    const time = d.toLocaleTimeString(undefined, { hour: 'numeric', minute: '2-digit' });
    if (days === 0) return `Today · ${time}`;
    if (days === 1) return `Tomorrow · ${time}`;
    if (days > 1 && days < 7) return `${d.toLocaleDateString(undefined, { weekday: 'long' })} · ${time}`;
    return `${d.toLocaleDateString(undefined, { weekday: 'short', month: 'short', day: 'numeric' })} · ${time}`;
  }
  function evRelative(iv) {
    if (!iv?.starts_at) return 'upcoming';
    const now = new Date();
    const d = new Date(iv.starts_at);
    const startOfDay = (x) => new Date(x.getFullYear(), x.getMonth(), x.getDate());
    const days = Math.round((startOfDay(d) - startOfDay(now)) / 86400000);
    if (days <= 0) return 'today';
    if (days === 1) return 'tomorrow';
    if (days < 7) return `in ${days} days`;
    if (days < 14) return 'next week';
    return `in ${days} days`;
  }
  function evDuration(iv) {
    if (!iv?.starts_at || !iv?.ends_at) return null;
    const mins = Math.round((new Date(iv.ends_at) - new Date(iv.starts_at)) / 60000);
    if (mins <= 0) return null;
    return mins;
  }
  function evWho(iv) {
    const att = iv?.attendees;
    if (!att) return '';
    if (typeof att === 'string') return att;
    if (Array.isArray(att)) {
      const first = att.find(a => a?.name) || att[0];
      return first?.name || first?.email || '';
    }
    return '';
  }

  function fmtEventWhen(ev) {
    const d = new Date(ev.starts_at);
    if (ev.all_day) return d.toLocaleDateString(undefined, { weekday: 'short', month: 'short', day: 'numeric' }) + ' · all day';
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
  function isPast(ev) { return new Date(ev.starts_at) < new Date(); }

  const monoDate = (iso) => iso ? new Date(iso).toLocaleDateString('en-US', { day: 'numeric', month: 'short' }) : '';

  // Activity timeline — derived honestly from real data. We DON'T fabricate
  // events: just the application being submitted, each interview on file, the
  // logged follow-ups, and the current status as the head of the line.
  const timeline = $derived.by(() => {
    if (!app) return [];
    const rows = [];
    const submittedAt = app.raw.applied_at || app.raw.created_at;
    if (submittedAt) {
      rows.push({
        ts: new Date(submittedAt).getTime(),
        date: monoDate(submittedAt),
        title: 'Application submitted',
        note: app.cv && app.cv !== '—' ? `CV ${app.cv}` : (app.source && app.source !== '—' ? `via ${app.source}` : ''),
        tag: 'positive'
      });
    }
    for (const iv of (interviews || [])) {
      if (!iv?.starts_at) continue;
      rows.push({
        ts: new Date(iv.starts_at).getTime(),
        date: monoDate(iv.starts_at),
        title: iv.summary || 'Interview',
        note: isPast(iv) ? 'Past event' : 'Scheduled',
        tag: 'accent'
      });
    }
    for (const f of (followUps || [])) {
      const when = f.occurred_at || f.created_at;
      rows.push({
        ts: when ? new Date(when).getTime() : Date.now(),
        date: monoDate(when),
        title: (f.note && f.note.trim()) ? f.note.trim() : 'Follow-up',
        note: f.channel ? `Follow-up · ${f.channel}` : 'Follow-up · you reached out directly',
        tag: '',
        followUp: f
      });
    }
    rows.push({
      ts: app.raw.updated_at ? new Date(app.raw.updated_at).getTime() : Date.now(),
      date: monoDate(app.raw.updated_at) || 'Now',
      title: `Status · ${STATUS_LABEL[app.status]}`,
      note: 'Current stage',
      tag: app.status === 'offer' ? 'offer' : app.status === 'interview' ? 'accent'
        : (app.status === 'rejected' || app.status === 'withdrawn') ? 'danger' : 'positive'
    });
    return rows.sort((a, b) => b.ts - a.ts);
  });
</script>

<svelte:head>
  <title>{app?.co ? `${app.co} — Pursuit` : 'Pursuit'}</title>
</svelte:head>

<svelte:window onkeydown={onWindowKeydown} />

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
            <button class="status-menu-item" class:current={app?.status === s} onclick={() => setStatus(s)}>
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
  <div class="det">
    {#if loading}
      <p style="color:var(--mute)">Loading…</p>
    {:else if notFound || !app}
      <div class="empty-tab">
        <h3>Application not found</h3>
        <p>It may have been deleted, or you might not have access. <a href="/app" style="color:var(--accent-text)">Back to Today →</a></p>
      </div>
    {:else}

      <!-- HEADER -->
      <div class="det-hd">
        {#if app.logoSrc}
          <img class="logo-big" src={app.logoSrc} alt={app.co} />
        {:else}
          <span class={`logo-big letter ${app.logoCls}`}>{app.coShort}</span>
        {/if}
        <div class="meta">
          <div class="co">{app.co}</div>
          <div class="role">{app.role}</div>
          <div class="sub">
            {#if app.raw.applied_at}<span>Applied <b>{appliedLong}</b></span>{/if}
            {#if app.source && app.source !== '—'}<span>Source <b>{app.source}</b></span>{/if}
            {#if app.cv && app.cv !== '—'}<span>CV <b>{app.cv}</b></span>{/if}
          </div>
        </div>
        <span class={`pill ${app.status}`}><span class="pdot"></span>{STATUS_LABEL[app.status]}</span>
      </div>

      <!-- TWO-COLUMN GRID -->
      <div class="det-grid">
        <!-- MAIN -->
        <div>
          {#if upcoming}
            <div class="det-next">
              <div class="k"><span class="d"></span>Next step · {evRelative(upcoming)}</div>
              <div class="ttl">{upcoming.summary || 'Interview'} · {evWhen(upcoming)}</div>
              <div class="mt">
                {#if evWho(upcoming)}<span><b>{evWho(upcoming)}</b></span>{/if}
                {#if upcoming.location}<span>{upcoming.location}</span>{/if}
                {#if evDuration(upcoming)}<span>{evDuration(upcoming)} min</span>{/if}
              </div>
              <div class="row">
                <button class="cta" onclick={openPlaybook}>Open interview prep
                  <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.7"><path d="M3 8h9M8 4l4 4-4 4" stroke-linecap="round" stroke-linejoin="round"/></svg>
                </button>
              </div>
              <p class="next-prep-note">AI brief on your interviewer — their style, what lands, and questions to ask.</p>
            </div>
          {:else if awaiting}
            <div class="det-next muted">
              <div class="k"><span class="d" style={`background:${quiet ? 'var(--warm)' : 'var(--mute-2)'}`}></span>{quiet ? 'Gone quiet' : 'Waiting to hear back'}</div>
              <div class="ttl">{STATUS_LABEL[app.status]}</div>
              <div class="mt">
                {#if quiet && waitDays != null}
                  <span>No reply in <b class="warn">{waitDays} days</b> — it might be a good time to reach out to them directly.</span>
                {:else}
                  <span>Still in the pipeline. We'll surface a nudge if it goes quiet.</span>
                {/if}
              </div>
              <div class="row">
                <button class="cta dark" onclick={openFollowUp}>Log a follow-up
                  <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.7"><path d="M3 8h9M8 4l4 4-4 4" stroke-linecap="round" stroke-linejoin="round"/></svg>
                </button>
              </div>
            </div>
          {:else}
            <div class="det-next muted">
              <div class="k"><span class="d" style="background:var(--mute-2)"></span>No upcoming step</div>
              <div class="ttl">{STATUS_LABEL[app.status]}</div>
              <div class="mt"><span>Nothing scheduled right now.</span></div>
              <div class="row">
                <button class="cta dark" onclick={openAddEvent}>Log an event
                  <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.7"><path d="M3 8h9M8 4l4 4-4 4" stroke-linecap="round" stroke-linejoin="round"/></svg>
                </button>
              </div>
            </div>
          {/if}

          <div class="sec-lbl">Activity</div>
          {#if timeline.length > 0}
            <div class="tl">
              {#each timeline as e}
                <div class={`tlrow ${e.tag}`}>
                  <span class="pt"></span>
                  <div class="d">{e.date}</div>
                  <div class="t">
                    {e.title}
                    {#if e.followUp}
                      <button class="tl-del" title="Delete follow-up" aria-label="Delete follow-up" onclick={() => deleteFollowUp(e.followUp)}>
                        <svg width="12" height="12" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M3 3l8 8M11 3l-8 8" stroke-linecap="round"/></svg>
                      </button>
                    {/if}
                  </div>
                  {#if e.note}<div class="n">{e.note}</div>{/if}
                </div>
              {/each}
            </div>
          {:else}
            <p style="color:var(--mute); font-size:13px;">No activity yet.</p>
          {/if}

          <!-- Interviews on file + add flow -->
          <div class="iv-section">
            <div class="iv-hd">
              <div class="sec-lbl" style="margin:0">Interviews <span class="iv-count">{interviews.length}</span></div>
              <button class="iv-add-btn" onclick={openAddEvent}>
                <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.6"><path d="M8 3v10M3 8h10" stroke-linecap="round"/></svg>
                Add event
              </button>
            </div>

            {#if interviewsLoading}
              <p style="color:var(--mute); font-size:13px;">Loading…</p>
            {:else if interviews.length > 0}
              {#each interviews as iv (iv.id)}
                <div class="iv-card" class:past={isPast(iv)}>
                  <div class="iv-card-main">
                    <div class="iv-when">{fmtEventWhen(iv)}</div>
                    <div class="iv-summary">{iv.summary || 'Untitled event'}</div>
                    {#if iv.location}<div class="iv-loc">📍 {iv.location}</div>{/if}
                  </div>
                  <button class="btn btn-danger" onclick={() => deleteInterview(iv)} title="Delete">Delete</button>
                </div>
              {/each}
            {:else}
              <p style="color:var(--mute); font-size:13px;">No interviews on file yet.</p>
            {/if}
          </div>
        </div>

        <!-- SIDE -->
        <div>
          <!-- Contact -->
          <div class="side-card">
            <div class="ttl">Contact</div>
            {#if app.raw.hiring_manager_name}
              <div class="person">
                <span class="iv-av sm">{hiringManagerInitials || '—'}</span>
                <div>
                  <div class="nm">{app.raw.hiring_manager_name}</div>
                  <div class="ro">Hiring manager</div>
                </div>
                {#if app.raw.hiring_manager_linkedin}
                  <a class="p-li" href={app.raw.hiring_manager_linkedin} target="_blank" rel="noopener" title="LinkedIn">
                    <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M3.5 6h2v6h-2zM4.5 3a1 1 0 1 1 0 2 1 1 0 0 1 0-2zM7 6h2v.9c.3-.5.9-1 1.8-1 1.6 0 2.2 1 2.2 2.6V12h-2V9c0-.9-.3-1.4-1.1-1.4-.6 0-1 .4-1 1.2V12H7z"/></svg>
                  </a>
                {/if}
              </div>
            {:else}
              <div class="person"><div><div class="nm" style="color:var(--mute)">No contact yet</div></div></div>
            {/if}
          </div>

          <!-- Details -->
          <div class="side-card">
            <div class="ttl">Details</div>
            <div class="kv"><span class="l">Status</span><span class="v">{STATUS_LABEL[app.status]}</span></div>
            <div class="kv"><span class="l">Last activity</span><span class="v">{fmtRelativeDate(app.raw.updated_at ?? app.raw.applied_at)}</span></div>
            <div class="kv"><span class="l">Source</span><span class="v">{app.source}</span></div>
            <div class="kv"><span class="l">Résumé</span><span class="v">{app.cv}</span></div>
          </div>

          <!-- Actions -->
          <div class="side-card">
            <div class="ttl">Actions</div>
            <div class="side-act">
              {#if awaiting}
                <button class="warn" onclick={openFollowUp}>
                  <span class="ic"><svg width="15" height="15" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M2 4h12v8H2zM2 4l6 5 6-5"/></svg></span>
                  Log a follow-up
                </button>
              {/if}
              <button onclick={openEdit}>
                <span class="ic"><svg width="15" height="15" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M11 2l3 3-8 8H3v-3z"/></svg></span>
                Add a note
              </button>
              <button onclick={openAddEvent}>
                <span class="ic"><svg width="15" height="15" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="2" y="3" width="12" height="11" rx="1.5"/><path d="M2 6h12M6 2v2M10 2v2"/></svg></span>
                Log an event
              </button>
              {#if app.raw.jd_url}
                <a class="act-link" href={app.raw.jd_url} target="_blank" rel="noopener">
                  <span class="ic"><svg width="15" height="15" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M5 11l6-6M6 5h5v5"/></svg></span>
                  Open job post
                </a>
              {/if}
            </div>
          </div>

          <!-- Interview prep link -->
          <div class="side-card playbook-card">
            <div class="ttl">Interview prep<span class="ai-pill">AI</span></div>
            <p class="playbook-blurb">AI brief on your interviewer — their style, what lands, and questions to ask.</p>
            <button class="playbook-link" onclick={openPlaybook}>
              Open interview prep
              <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.7"><path d="M3 8h9M8 4l4 4-4 4" stroke-linecap="round" stroke-linejoin="round"/></svg>
            </button>
          </div>
        </div>
      </div>
    {/if}
  </div>
</div>

{#if toast}
  <div class="toast">
    <span class="ok"><svg width="10" height="10" viewBox="0 0 12 12" fill="none" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2.5 6.5l2.5 2.5 4.5-5"/></svg></span>
    {toast}
  </div>
{/if}

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

{#if showEventModal}
  <div class="ev-overlay" onclick={closeEventModal} role="presentation">
    <div class="ev-card" onclick={(e) => e.stopPropagation()} role="dialog" aria-modal="true" aria-label="Add an interview">
      <button class="x-close" onclick={closeEventModal} aria-label="Close">
        <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="1.6"><path d="M3 3l8 8M11 3l-8 8" stroke-linecap="round"/></svg>
      </button>
      <div class="add-hd">
        <h3>Add an interview</h3>
        <p>Drop a calendar file, paste a screenshot, or just paste the email body — we'll extract the event.</p>
      </div>
      <div class="zones">
        <!-- LEFT — .ics -->
        <div class="zone">
          <div class="zone-hd">
            <span class="zone-ic"><svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="2" y="3" width="12" height="11" rx="1.5"/><path d="M2 6h12M6 2v2M10 2v2"/></svg></span>
            <div>
              <div class="zone-title">Calendar file</div>
              <div class="zone-sub">.ics from Google / Outlook / Apple</div>
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

        <!-- RIGHT — screenshot/text AI -->
        <div class="zone">
          <div class="zone-hd">
            <span class="zone-ic accent"><svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="2" y="3" width="12" height="10" rx="1.5"/><circle cx="5.5" cy="6.5" r="1"/><path d="M2 11l3.5-3.5 3 3 2-2L14 11"/></svg></span>
            <div>
              <div class="zone-title">Screenshot or email text<span class="ai-pill">AI</span></div>
              <div class="zone-sub">Gmail invite, Calendar screenshot, anything readable</div>
            </div>
          </div>
          <label class="drop drop-image" class:drag={aiDragOver} ondragover={onAiDragOver} ondragleave={() => (aiDragOver = false)} ondrop={onAiDrop}>
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
          <textarea class="ai-ta" rows="3" placeholder={"You're invited to: Stripe — Technical screen&#10;When: Tue, May 28, 2:00 PM EDT&#10;Where: Google Meet"} bind:value={aiText} onpaste={onAiPaste}></textarea>
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
  </div>
{/if}

{#if showFollowUpModal}
  <div class="ev-overlay" onclick={closeFollowUp} role="presentation">
    <form class="ev-card fu-card" onclick={(e) => e.stopPropagation()} onsubmit={saveFollowUp} role="dialog" aria-modal="true" aria-label="Log a follow-up">
      <button type="button" class="x-close" onclick={closeFollowUp} aria-label="Close">
        <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="1.6"><path d="M3 3l8 8M11 3l-8 8" stroke-linecap="round"/></svg>
      </button>
      <div class="add-hd">
        <h3>Log a follow-up</h3>
        <p>Record something you did yourself — Pursuit doesn't send anything. We'll reset the quiet clock.</p>
      </div>
      <div class="fu-fields">
        <label class="fu-label">What did you do?
          <textarea class="fu-ta" rows="3" placeholder="e.g. Emailed Sarah to check in" bind:value={fuNote}></textarea>
        </label>
        <div class="fu-row">
          <label class="fu-label">Channel
            <select class="fu-input" bind:value={fuChannel}>
              <option value=""></option>
              {#each FU_CHANNELS as c}<option value={c}>{c}</option>{/each}
            </select>
          </label>
          <label class="fu-label">Date
            <input class="fu-input" type="date" bind:value={fuDate} />
          </label>
        </div>
      </div>
      <div class="modal-actions">
        <button type="button" class="btn" onclick={closeFollowUp}>Cancel</button>
        <button type="submit" class="btn btn-primary" disabled={fuSaving}>
          {fuSaving ? 'Saving…' : 'Save follow-up'}
        </button>
      </div>
    </form>
  </div>
{/if}

<style>
  .body { padding: 28px; }
  .det { max-width: 980px; margin: 0 auto; }

  /* HEADER */
  .det-hd { display: flex; align-items: flex-start; gap: 18px; margin-bottom: 30px; }
  .logo-big { width: 56px; height: 56px; border-radius: 15px; background: var(--card); object-fit: contain; padding: 8px; border: 1px solid var(--rule); flex-shrink: 0; }
  .logo-big.letter { display: grid; place-items: center; padding: 0; color: var(--ink); font-size: 22px; font-weight: 600; background: var(--surface-2); }
  .det-hd .meta { flex: 1; min-width: 0; }
  .det-hd .co { font-size: 24px; font-weight: 600; letter-spacing: -0.025em; }
  .det-hd .role { font-size: 15px; color: var(--mute); margin-top: 3px; }
  .det-hd .sub { font-size: 12.5px; color: var(--mute-2); margin-top: 10px; display: flex; gap: 16px; flex-wrap: wrap; }
  .det-hd .sub b { color: var(--ink-2); font-weight: 500; }

  /* Pills (page-scoped) */
  .pill { display: inline-flex; align-items: center; gap: 5px; padding: 4px 10px; border-radius: 99px; font-size: 12px; font-weight: 500; background: var(--surface-2); color: var(--ink-2); width: max-content; flex-shrink: 0; }
  .pill .pdot { width: 6px; height: 6px; border-radius: 50%; background: var(--mute-2); }
  .pill.wishlist { background: var(--surface-2); color: var(--mute); }
  .pill.applied { background: var(--surface-2); color: var(--ink-2); }
  .pill.screen { background: var(--accent-tint); color: var(--accent-text); }
  .pill.screen .pdot { background: var(--accent); }
  .pill.interview { background: var(--warm-tint); color: var(--warm-text); }
  .pill.interview .pdot { background: var(--warm); }
  .pill.offer { background: var(--positive-tint); color: var(--positive-text); }
  .pill.offer .pdot { background: var(--positive); }
  .pill.rejected, .pill.withdrawn { background: var(--danger-tint); color: var(--danger-text); }
  .pill.rejected .pdot, .pill.withdrawn .pdot { background: var(--danger); }

  /* GRID */
  .det-grid { display: grid; grid-template-columns: 1fr 320px; gap: 36px; align-items: start; }

  /* NEXT-STEP */
  .det-next { background: var(--ink); color: #fff; border-radius: 16px; padding: 24px; margin-bottom: 26px; position: relative; overflow: hidden; }
  .det-next.muted { background: var(--surface-2); color: var(--ink); }
  .det-next .k { font-size: 11px; font-weight: 600; letter-spacing: 0.08em; text-transform: uppercase; color: rgba(255,255,255,0.6); margin-bottom: 14px; display: inline-flex; align-items: center; gap: 8px; }
  .det-next.muted .k { color: var(--mute); }
  .det-next .k .d { width: 6px; height: 6px; border-radius: 50%; background: var(--warm); }
  .det-next .ttl { font-size: 20px; font-weight: 500; letter-spacing: -0.02em; margin-bottom: 6px; }
  .det-next.muted .ttl { color: var(--ink); }
  .det-next .mt { font-size: 13px; color: rgba(255,255,255,0.7); margin-bottom: 20px; display: flex; gap: 14px; flex-wrap: wrap; }
  .det-next.muted .mt { color: var(--mute); }
  .det-next .mt b { color: #fff; font-weight: 500; }
  .det-next.muted .mt b.warn { color: var(--warm-text); }
  .det-next .row { display: flex; gap: 10px; }
  .det-next .cta { background: #fff; color: var(--ink); border: none; border-radius: 9px; padding: 11px 16px; font-size: 13.5px; font-weight: 600; cursor: pointer; display: inline-flex; align-items: center; gap: 7px; }
  .det-next .cta.dark { background: var(--ink); color: #fff; }
  .det-next .ghost { background: rgba(255,255,255,0.1); color: #fff; border: 1px solid rgba(255,255,255,0.18); border-radius: 9px; padding: 11px 16px; font-size: 13.5px; font-weight: 500; cursor: pointer; }
  .det-next .ghost:hover { background: rgba(255,255,255,0.18); }

  .sec-lbl { font-size: 11.5px; font-weight: 600; letter-spacing: 0.06em; text-transform: uppercase; color: var(--mute-2); margin-bottom: 16px; }

  /* TIMELINE */
  .tl { position: relative; padding-left: 26px; }
  .tl::before { content: ""; position: absolute; left: 4px; top: 8px; bottom: 8px; width: 1.5px; background: var(--rule); }
  .tlrow { position: relative; padding-bottom: 22px; }
  .tlrow:last-child { padding-bottom: 0; }
  .tlrow .pt { position: absolute; left: -26px; top: 3px; width: 9px; height: 9px; border-radius: 50%; background: var(--card); border: 2px solid var(--rule-strong); }
  .tlrow.accent .pt { border-color: var(--warm); background: var(--warm); box-shadow: 0 0 0 4px var(--warm-tint); }
  .tlrow.positive .pt { border-color: var(--positive); background: var(--positive); }
  .tlrow.offer .pt { border-color: var(--positive); background: var(--positive); box-shadow: 0 0 0 4px var(--positive-tint); }
  .tlrow.danger .pt { border-color: var(--danger); background: var(--danger); }
  .tlrow .d { font-family: var(--mono, ui-monospace, monospace); font-size: 11px; color: var(--mute); margin-bottom: 3px; }
  .tlrow .t { font-size: 13.5px; font-weight: 500; display: flex; align-items: center; gap: 6px; }
  .tlrow .n { font-size: 12.5px; color: var(--mute); margin-top: 2px; }
  .tl-del { width: 20px; height: 20px; flex-shrink: 0; border: 0; background: transparent; color: var(--mute-2); border-radius: 5px; display: inline-grid; place-items: center; cursor: pointer; opacity: 0; transition: opacity 100ms ease, background 100ms ease, color 100ms ease; }
  .tlrow:hover .tl-del { opacity: 1; }
  .tl-del:hover { background: var(--danger-tint); color: var(--danger-text); }

  /* FOLLOW-UP MODAL */
  .fu-card { max-width: 460px; }
  .fu-fields { display: flex; flex-direction: column; gap: 14px; }
  .fu-label { display: flex; flex-direction: column; gap: 6px; font-size: 12px; color: var(--mute); }
  .fu-row { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
  .fu-ta, .fu-input { font: inherit; color: var(--ink); background: var(--surface); border: 1px solid var(--rule); border-radius: 8px; padding: 9px 11px; font-size: 13.5px; outline: none; transition: border-color 100ms ease; box-sizing: border-box; width: 100%; }
  .fu-ta { resize: vertical; line-height: 1.5; }
  .fu-ta:focus, .fu-input:focus { border-color: var(--accent); }

  /* SIDE CARDS */
  .side-card { background: var(--card); border: 1px solid var(--rule); border-radius: 14px; padding: 18px; box-shadow: var(--sh-1); margin-bottom: 16px; }
  .side-card .ttl { font-size: 11.5px; font-weight: 600; letter-spacing: 0.06em; text-transform: uppercase; color: var(--mute-2); margin-bottom: 14px; display: inline-flex; align-items: center; gap: 8px; }
  .side-card .person { display: flex; align-items: center; gap: 12px; }
  .side-card .person .nm { font-size: 13.5px; font-weight: 500; }
  .side-card .person .ro { font-size: 12px; color: var(--mute); margin-top: 1px; }
  .iv-av { background: linear-gradient(155deg, oklch(0.6 0.16 30), oklch(0.46 0.17 32)); color: #fff; border-radius: 11px; display: inline-flex; align-items: center; justify-content: center; font-weight: 600; flex-shrink: 0; }
  .iv-av.sm { width: 38px; height: 38px; font-size: 13px; }
  .p-li { margin-left: auto; display: inline-flex; align-items: center; justify-content: center; width: 32px; height: 32px; background: var(--surface-2); border: 1px solid var(--rule); border-radius: 8px; color: #0a66c2; text-decoration: none; flex-shrink: 0; }
  .side-card .kv { display: flex; justify-content: space-between; font-size: 13px; padding: 9px 0; border-top: 1px solid var(--rule); }
  .side-card .kv:first-of-type { border-top: none; padding-top: 0; }
  .side-card .kv .l { color: var(--mute); }
  .side-card .kv .v { font-weight: 500; }
  .side-act { display: flex; flex-direction: column; gap: 8px; }
  .side-act button, .side-act .act-link { display: flex; align-items: center; gap: 10px; width: 100%; text-align: left; font: inherit; font-size: 13px; font-weight: 500; color: var(--ink-2); background: var(--card); border: 1px solid var(--rule); border-radius: 9px; padding: 11px 13px; cursor: pointer; text-decoration: none; box-sizing: border-box; }
  .side-act button:hover, .side-act .act-link:hover { background: var(--surface-2); border-color: var(--rule-strong); }
  .side-act .ic { color: var(--mute); display: inline-flex; }
  .side-act button.warn { color: var(--warm-text); }
  .side-act button.warn .ic { color: var(--warm-text); }

  .playbook-card .ai-pill { margin-left: 6px; }
  .playbook-blurb { font-size: 12.5px; color: var(--mute); line-height: 1.5; margin: 0 0 14px; }
  .playbook-link { display: inline-flex; align-items: center; gap: 7px; width: 100%; justify-content: center; font: inherit; font-size: 13px; font-weight: 600; color: #fff; background: var(--ink); border: none; border-radius: 9px; padding: 11px 13px; cursor: pointer; }
  .playbook-link:hover { background: var(--ink-2); }
  .ai-pill { font-size: 10px; font-weight: 600; color: var(--accent-text); background: var(--accent-tint); border-radius: 4px; padding: 1px 5px; letter-spacing: .04em; }

  /* INTERVIEWS SECTION */
  .iv-section { margin-top: 30px; padding-top: 24px; border-top: 1px solid var(--rule); }
  .iv-hd { display: flex; align-items: center; justify-content: space-between; margin-bottom: 14px; }
  .iv-count { font-size: 11px; background: var(--accent-tint); color: var(--accent-text); padding: 1px 7px; border-radius: 99px; margin-left: 4px; font-weight: 500; }
  .iv-add-btn { display: inline-flex; align-items: center; gap: 6px; background: var(--card); border: 1px solid var(--rule); border-radius: 8px; padding: 6px 11px; font-size: 12.5px; font-weight: 600; color: var(--ink-2); cursor: pointer; }
  .iv-add-btn:hover { background: var(--surface-2); border-color: var(--rule-strong); }
  .iv-card { display: flex; align-items: flex-start; justify-content: space-between; gap: 16px; padding: 14px 16px; border: 1px solid var(--rule); border-radius: 10px; background: var(--card); margin-bottom: 8px; box-shadow: var(--sh-1); }
  .iv-card.past { opacity: 0.6; }
  .iv-card-main { flex: 1; min-width: 0; }
  .iv-when { font-size: 12px; color: var(--accent-text); font-weight: 500; margin-bottom: 2px; }
  .iv-summary { font-size: 14px; color: var(--ink); margin-bottom: 2px; }
  .iv-loc { font-size: 12px; color: var(--mute); margin-top: 2px; }

  /* ADD-EVENT MODAL */
  .ev-overlay { position: fixed; inset: 0; background: rgba(10,10,13,0.55); backdrop-filter: blur(8px); -webkit-backdrop-filter: blur(8px); display: grid; place-items: center; z-index: 200; padding: 24px; overflow-y: auto; }
  .ev-card { position: relative; width: 100%; max-width: 540px; background: var(--card); border: 1px solid var(--rule); border-radius: 18px; padding: 26px 28px 24px; box-shadow: 0 24px 80px -8px rgba(10,10,13,0.30), var(--sh-1); margin: 24px auto; display: flex; flex-direction: column; }
  .x-close { position: absolute; top: 14px; right: 14px; width: 28px; height: 28px; border-radius: 8px; background: transparent; border: 0; display: grid; place-items: center; color: var(--mute); cursor: pointer; transition: background 100ms ease, color 100ms ease; }
  .x-close:hover { background: var(--surface-2); color: var(--ink); }
  .add-hd { margin-bottom: 16px; padding-right: 26px; }
  .add-hd h3 { font-size: 15px; font-weight: 600; margin: 0 0 4px; letter-spacing: -0.015em; }
  .add-hd p { font-size: 13px; color: var(--mute); margin: 0; line-height: 1.5; }
  .zones { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
  .zone { background: var(--surface); border: 1px solid var(--rule); border-radius: 12px; padding: 14px 16px; display: flex; flex-direction: column; gap: 10px; }
  .zone-hd { display: grid; grid-template-columns: 32px 1fr; gap: 10px; align-items: center; }
  .zone-ic { width: 32px; height: 32px; border-radius: 8px; background: var(--surface-2); color: var(--ink-2); display: grid; place-items: center; }
  .zone-ic.accent { background: var(--accent-tint); color: var(--accent-text); }
  .zone-title { font-size: 13.5px; font-weight: 600; color: var(--ink); display: inline-flex; align-items: center; gap: 6px; }
  .zone-sub { font-size: 12px; color: var(--mute); margin-top: 1px; }
  .drop { background: var(--card); border: 1.5px dashed var(--rule-strong); border-radius: 10px; padding: 16px 12px; display: flex; flex-direction: column; align-items: center; gap: 3px; color: var(--mute); transition: border-color 120ms, background 120ms; cursor: pointer; }
  .drop:hover, .drop.drag { border-color: var(--accent); background: var(--accent-tint); color: var(--accent-text); }
  .drop-l1 { font-size: 12.5px; font-weight: 500; color: var(--ink-2); }
  .drop:hover .drop-l1, .drop.drag .drop-l1 { color: var(--accent-text); }
  .drop-l2 { font-size: 11.5px; color: var(--mute-2); }
  .drop kbd { font-family: var(--mono, ui-monospace, monospace); font-size: 10px; background: var(--surface-2); border: 1px solid var(--rule); border-bottom-width: 2px; border-radius: 3px; padding: 0 4px; color: var(--ink-2); }
  .or { font-size: 11.5px; color: var(--mute-2); text-align: center; }
  .ai-ta { width: 100%; font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 11.5px; line-height: 1.5; color: var(--ink); background: var(--card); border: 1px solid var(--rule); border-radius: 8px; padding: 8px 10px; outline: none; resize: vertical; box-sizing: border-box; }
  .ai-ta:focus { border-color: var(--accent); }
  .ai-pill { font-size: 10px; }
  .zone-parse { margin-top: 4px; align-self: flex-start; }
  .ics-preview { margin-top: 18px; padding-top: 16px; border-top: 1px solid var(--rule); }
  .ics-preview h4 { font-size: 11.5px; font-weight: 600; color: var(--mute); text-transform: uppercase; letter-spacing: 0.04em; margin: 0 0 10px; }
  .prev-row { background: var(--accent-tint); border: 1px solid var(--accent); border-radius: 10px; padding: 12px 14px; margin-bottom: 10px; }
  .prev-summary { font-size: 13.5px; font-weight: 600; color: var(--ink); }
  .prev-when { font-size: 12.5px; color: var(--accent-text); margin-top: 3px; font-weight: 500; }
  .prev-loc { font-size: 12px; color: var(--mute); margin-top: 4px; }
  .dossier-err { color: var(--danger-text); background: var(--danger-tint); border: 1px solid var(--danger-tint); border-radius: 8px; padding: 8px 12px; font-size: 13px; margin: .75rem 0 0; }

  /* TOAST */
  .toast { position: fixed; bottom: 26px; left: 50%; transform: translateX(-50%); background: var(--ink); color: #fff; font-size: 13px; font-weight: 500; padding: 11px 18px; border-radius: 10px; z-index: 200; display: inline-flex; align-items: center; gap: 9px; box-shadow: 0 16px 36px -12px rgba(20,20,50,0.5); animation: rise .2s ease; }
  .toast .ok { width: 16px; height: 16px; border-radius: 50%; background: var(--positive); display: inline-flex; align-items: center; justify-content: center; }
  @keyframes rise { from { transform: translate(-50%, 12px); opacity: 0; } }

  /* STATUS MENU */
  .status-wrap { position: relative; }
  .status-menu { position: absolute; top: calc(100% + 6px); right: 0; z-index: 50; background: var(--card); border: 1px solid var(--rule); border-radius: 8px; box-shadow: var(--sh-pop); padding: 4px; min-width: 180px; display: flex; flex-direction: column; gap: 1px; }
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

  .empty-tab { border: 1px dashed var(--rule); border-radius: 12px; padding: 32px; text-align: center; background: var(--card); }
  .empty-tab h3 { margin: 0 0 .5rem; font-size: 16px; font-weight: 500; color: var(--ink); }
  .empty-tab p { color: var(--mute); margin: 0; font-size: 13.5px; }

  /* MOBILE */
  @media (max-width: 820px) {
    .body { padding: 18px 14px; }
    .det-grid { grid-template-columns: 1fr; gap: 24px; }
    .det-hd { gap: 12px; }
    .logo-big { width: 48px; height: 48px; }
    .det-hd .co { font-size: 21px; }
    .zones { grid-template-columns: 1fr; gap: 10px; }
    .modal-overlay { padding: 0; }
    .modal { max-width: 100%; border-radius: 0; min-height: 100vh; padding: 1rem; }
    .ev-overlay { padding: 0; }
    .ev-card { max-width: 100%; border-radius: 0; min-height: 100vh; margin: 0; padding: 20px 16px; }
    .fields { grid-template-columns: 1fr; }
    .fields .span-2 { grid-column: auto; }
    .fu-row { grid-template-columns: 1fr; }
  }
</style>
