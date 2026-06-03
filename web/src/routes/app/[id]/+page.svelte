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

  // Interview prep (dossier) — inline section
  let dossier = $state(null);
  let dossierLoading = $state(true);
  let generating = $state(false);
  let genError = $state('');
  let interviewerInput = $state('');

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
    loadDossier();
  });

  async function loadApp() {
    loading = true;
    notFound = false;
    try {
      const raw = await call(`/api/applications/${id}`);
      app = toDisplayApp(raw);
      if (!interviewerInput) {
        if (dossier?.interviewer_name) interviewerInput = dossier.interviewer_name;
        else if (app?.raw?.hiring_manager_name) interviewerInput = app.raw.hiring_manager_name;
      }
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
    genError = '';
    try {
      const d = await call(`/api/applications/${id}/dossier`);
      dossier = d || null;
      if (dossier?.interviewer_name) interviewerInput = dossier.interviewer_name;
    } catch (e) {
      // 404 / empty / "no dossier" → not generated yet
      dossier = null;
    } finally {
      dossierLoading = false;
    }
  }

  async function generateDossier() {
    if (generating) return;
    generating = true;
    genError = '';
    try {
      const d = await call(`/api/applications/${id}/dossier/refresh`, {
        method: 'POST',
        body: JSON.stringify({ interviewer_name: interviewerInput.trim() || undefined })
      });
      dossier = d;
      interviewerInput = d.interviewer_name ?? interviewerInput;
    } catch (e) {
      genError = friendlyGenErr(e.message);
    } finally {
      generating = false;
    }
  }

  function friendlyGenErr(msg) {
    const m = String(msg || '');
    if (m.includes('rate_limit_error') || m.includes('429'))
      return 'AI usage limit hit — wait a minute and try again.';
    if (m.includes('http 504') || /\btimeout\b/i.test(m))
      return 'Web search timed out — try again.';
    if (m.includes('http 5') || m.includes('not configured'))
      return 'Something went wrong — try again in a moment.';
    return m || 'Could not generate the interview prep.';
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

  // ── Interview-prep (dossier) derived view data ───────────────
  const dosContent = $derived(dossier?.content ?? null);
  const dosInterviewer = $derived(dosContent?.interviewer ?? null);
  const dosIvInitials = $derived(
    initialsOf(dosInterviewer?.name ?? dossier?.interviewer_name ?? app?.raw?.hiring_manager_name ?? '')
  );
  const dosIvName = $derived(
    dosInterviewer?.name ?? dossier?.interviewer_name ?? app?.raw?.hiring_manager_name ?? 'Your interviewer'
  );
  const dosMeeting = $derived(dossier?.meeting ?? null);

  function dosFmtWhen(m) {
    if (!m) return '—';
    if (m.starts_at) {
      const d = new Date(m.starts_at);
      const now = new Date();
      const startOfDay = x => new Date(x.getFullYear(), x.getMonth(), x.getDate());
      const days = Math.round((startOfDay(d) - startOfDay(now)) / 86400000);
      const time = d.toLocaleTimeString(undefined, { hour: 'numeric', minute: '2-digit' });
      if (days === 0) return `Today · ${time}`;
      if (days === 1) return `Tomorrow · ${time}`;
      if (days > 1 && days < 7) return `${d.toLocaleDateString(undefined, { weekday: 'long' })} · ${time}`;
      if (days < 0) return `${d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })} · ${time}`;
      return `${d.toLocaleDateString(undefined, { weekday: 'short', month: 'short', day: 'numeric' })} · ${time}`;
    }
    return m.when ?? '—';
  }
  function dosFmtDuration(m) {
    if (!m) return '—';
    if (m.starts_at && m.ends_at) {
      const mins = Math.round((new Date(m.ends_at) - new Date(m.starts_at)) / 60000);
      if (mins <= 0) return m.duration ?? '—';
      return mins >= 60 && mins % 60 === 0 ? `${mins / 60}h` : `${mins} min`;
    }
    return m.duration ?? '—';
  }
  const dosFactWhen     = $derived(dosFmtWhen(dosMeeting));
  const dosFactDuration = $derived(dosFmtDuration(dosMeeting));
  const dosFactMedium   = $derived(dosMeeting?.medium ?? '—');
  const dosFactPanel    = $derived(dosMeeting?.panel ?? '—');

  // Company brief (editorial) — from dossier content.company.*
  const dosCompany = $derived(dosContent?.company ?? null);
  const companyBlurb = $derived(dosCompany?.blurb ?? '');
  const companyAbout = $derived(dosCompany?.direction ?? dosCompany?.about ?? '');
  const companyFacts = $derived.by(() => {
    const c = dosCompany;
    if (!c) return [];
    const out = [];
    if (c.stage)     out.push({ lbl: 'Stage', val: c.stage });
    if (c.employees) out.push({ lbl: 'Size', val: c.employees });
    if (c.founded)   out.push({ lbl: 'Founded', val: c.founded });
    if (c.hq)        out.push({ lbl: 'HQ', val: c.hq });
    return out;
  });
  const companyProcess = $derived.by(() => {
    const p = dosCompany?.process;
    if (!Array.isArray(p)) return [];
    return p.map(s => s?.kind || s?.detail || '').filter(Boolean);
  });

  // AI tips box — up to 3 company watch-fors.
  const tips = $derived((dosCompany?.watch_fors ?? []).slice(0, 3));
  function dosSigDomain(src) {
    if (!src) return '';
    try {
      return new URL(src.startsWith('http') ? src : `https://${src}`).hostname.replace(/^www\./, '');
    } catch { return src; }
  }
  const dosGeneratedAgo = $derived(dossier?.generatedAgo ?? '');

  const appliedLong = $derived(app ? fmtLongDate(app.raw.applied_at) : '');

  // Next event = the soonest interview dated now-or-later.
  const upcoming = $derived.by(() => {
    const now = Date.now();
    const future = (interviews || [])
      .filter(iv => iv?.starts_at && new Date(iv.starts_at).getTime() >= now)
      .sort((a, b) => new Date(a.starts_at) - new Date(b.starts_at));
    return future[0] || null;
  });

  // Next interview = soonest future interview, else the dossier meeting.
  const nextWhen = $derived.by(() => {
    if (upcoming) return evWhen(upcoming);
    if (dosMeeting?.starts_at) return dosFmtWhen(dosMeeting);
    return '';
  });
  const nextTitle = $derived.by(() => {
    if (upcoming) return `${upcoming.summary || 'Interview'} · ${evWhen(upcoming)}`;
    if (dosMeeting?.starts_at) return `${dosMeeting.panel || 'Interview'} · ${dosFmtWhen(dosMeeting)}`;
    return '';
  });
  const hasNext = $derived(!!(upcoming || dosMeeting?.starts_at));
  const nextRows = $derived.by(() => {
    const rows = [];
    if (upcoming) {
      const who = evWho(upcoming) || app?.raw?.hiring_manager_name;
      if (who) rows.push(`${who} (Hiring manager)`);
      if (upcoming.location) rows.push(upcoming.location);
      const mins = evDuration(upcoming);
      if (mins) rows.push(`${mins} min`);
    } else if (dosMeeting?.starts_at) {
      const who = dosInterviewer?.name || dossier?.interviewer_name || app?.raw?.hiring_manager_name;
      if (who) rows.push(`${who} (Hiring manager)`);
      if (dosFactMedium !== '—') rows.push(dosFactMedium);
      if (dosFactDuration !== '—') rows.push(dosFactDuration);
    }
    return rows;
  });
  // Label the prep person "Hiring manager" only when nothing's scheduled.
  const personLabel = $derived(hasNext ? 'Likely interviewer' : 'Hiring manager');

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
        <!-- LEFT — AI prep hero + activity -->
        <div class="left">

          <!-- ── INTERVIEW PREP (inline dossier) ── -->
          <div id="interview-prep" class="prep-lead">
            <span class="ai-pill"><span class="spark">✦</span> AI</span>
            <h2>Interview prep</h2>
            {#if dosGeneratedAgo && dossier}
              <p class="prep-gen">
                <span class="sp" aria-hidden="true">
                  <svg width="13" height="13" viewBox="0 0 13 13" fill="none">
                    <path d="M6.5 1.5C6.5 4.5 4.5 6.5 1.5 6.5C4.5 6.5 6.5 8.5 6.5 11.5C6.5 8.5 8.5 6.5 11.5 6.5C8.5 6.5 6.5 4.5 6.5 1.5Z" fill="currentColor"/>
                  </svg>
                </span>
                Generated by Pursuit · {dosGeneratedAgo}
              </p>
            {/if}
          </div>

          {#if dossierLoading}
            <p style="color:var(--mute); font-size:13px;">Loading…</p>

          {:else if !dossier}
            <!-- Generate / empty state -->
            <div class="generate-card">
              <div class="gen-icon" aria-hidden="true">
                <svg width="28" height="28" viewBox="0 0 24 24" fill="none">
                  <path d="M12 3C12 7.97 8.97 11 4 11C8.97 11 12 14.03 12 19C12 14.03 15.03 11 20 11C15.03 11 12 7.97 12 3Z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/>
                </svg>
              </div>
              {#if generating}
                <h3>Researching {app.co}{interviewerInput ? ` & ${interviewerInput}` : ''}…</h3>
                <p class="gen-sub">Claude is searching the web for recent posts, talks, and the company's current direction. This typically takes 30–60 seconds.</p>
                <div class="big-spinner"></div>
              {:else}
                <h3>Generate interview prep</h3>
                <p class="gen-sub">
                  We'll build an AI brief on the person interviewing you — their background, how they tend to interview, what lands well, and smart questions to ask — pulled from public posts, talks, papers, and company news. Add a name below to make it about a specific interviewer.
                </p>
                <div class="gen-row">
                  <input
                    class="gen-input"
                    type="text"
                    placeholder="Interviewer name (optional) — e.g. Sarah Chen"
                    bind:value={interviewerInput}
                    disabled={generating}
                    onkeydown={(e) => e.key === 'Enter' && generateDossier()}
                  />
                  <button class="btn-generate" onclick={generateDossier} disabled={generating}>
                    Generate interview prep
                  </button>
                </div>
                {#if genError}
                  <p class="gen-err">{genError}</p>
                {/if}
              {/if}
            </div>

          {:else}
            <!-- Full brief -->

            <!-- AI tips box (baby blue) -->
            {#if tips.length}
              <section class="tips">
                <div class="tips-hd">
                  <span class="tips-spark">✦</span>
                  <h3>Tips for this one</h3>
                  <span class="tips-ai">AI</span>
                </div>
                <ul class="tips-list">
                  {#each tips as t}
                    <li><span class="tip-dot"></span><span>{t}</span></li>
                  {/each}
                </ul>
              </section>
            {/if}

            <!-- Company brief — Editorial -->
            {#if dosCompany && (companyBlurb || companyAbout || companyFacts.length || companyProcess.length)}
              <section class="card brief-1">
                <div class="b1-hd">
                  {#if app.logoSrc}
                    <img class="logo" src={app.logoSrc} alt={app.co} />
                  {:else}
                    <div class="logo letter">{app.coShort}</div>
                  {/if}
                  <div class="b1-titles">
                    <div class="b1-name-row">
                      <h3>{app.co}</h3>
                      <span class="ai-pill"><span class="spark">✦</span> AI</span>
                    </div>
                    {#if companyBlurb}<div class="b1-headline">{companyBlurb}</div>{/if}
                  </div>
                </div>
                {#if companyAbout}<p class="about">{companyAbout}</p>{/if}
                {#if companyFacts.length}
                  <div class="b1-facts">
                    {#each companyFacts as f}
                      <div class="b1-cell">
                        <div class="f-lbl">{f.lbl}</div>
                        <div class="f-val">{f.val}</div>
                      </div>
                    {/each}
                  </div>
                {/if}
                {#if companyProcess.length}
                  <div class="b1-process">
                    <div class="proc-lbl">Interview process</div>
                    <div class="proc-chips">
                      {#each companyProcess as step, i}
                        <span class="chip">{step}</span>{#if i < companyProcess.length - 1}<span class="proc-arrow">→</span>{/if}
                      {/each}
                    </div>
                  </div>
                {/if}
              </section>
            {/if}

            <!-- Hiring manager / interviewer -->
            {#if dosInterviewer || dosContent?.snapshot || app.raw.hiring_manager_name}
              <section class="card">
                <div class="card-hd"><h3>{personLabel}</h3><span class="ai-tag">AI{dosGeneratedAgo ? ` · ${dosGeneratedAgo}` : ''}</span></div>
                <div class="person">
                  <div class="p-av">{dosIvInitials || '?'}</div>
                  <div class="p-info">
                    <div class="p-name-row">
                      <h4>{dosIvName}</h4>
                      <span class="role-tag">{personLabel}</span>
                    </div>
                    {#if dosInterviewer?.role}<div class="p-role">{dosInterviewer.role}</div>{/if}
                    {#if dosInterviewer?.prior?.length}
                      <div class="p-prior">{dosInterviewer.prior.join(' · ')}</div>
                    {/if}
                  </div>
                  {#if dosInterviewer?.links?.length}
                    <div class="p-links">
                      {#each dosInterviewer.links as l}
                        <a href={l.href} target="_blank" rel="noopener">{l.label}</a>
                      {/each}
                    </div>
                  {/if}
                </div>
                {#if dosContent?.snapshot}
                  <p class="snapshot">{@html dosContent.snapshot}</p>
                {/if}
                {#if dosContent?.style?.lead || dosContent?.style?.tells?.length}
                  <div class="tells-block">
                    {#if dosContent.style.lead}
                      <div class="tells-hd">How {dosIvName.split(' ')[0]} interviews</div>
                      <p class="tells-lead">{dosContent.style.lead}</p>
                    {/if}
                    {#if dosContent.style.tells?.length}
                      <div class="tells">
                        {#each dosContent.style.tells as t}
                          <div class="tell"><div class="t-lbl">{t.lbl}</div><div class="t-val">{t.val}</div></div>
                        {/each}
                      </div>
                    {/if}
                  </div>
                {/if}
              </section>
            {/if}

            <!-- Lands / avoid -->
            {#if dosContent?.lands?.length || dosContent?.avoid?.length}
              <section class="la-grid">
                {#if dosContent.lands?.length}
                  <div class="la-col lands">
                    <h3><span class="glyph">✓</span> What lands</h3>
                    <ul>{#each dosContent.lands as l}<li><span class="g">✓</span><span>{l}</span></li>{/each}</ul>
                  </div>
                {/if}
                {#if dosContent.avoid?.length}
                  <div class="la-col avoid">
                    <h3><span class="glyph">✕</span> What to avoid</h3>
                    <ul>{#each dosContent.avoid as a}<li><span class="g">✕</span><span>{a}</span></li>{/each}</ul>
                  </div>
                {/if}
              </section>
            {/if}

            <!-- Recent signals -->
            {#if dosContent?.signals?.length}
              <section class="card">
                <div class="card-hd"><h3>Recent signals</h3><span class="ai-tag">AI · web search</span></div>
                <ul class="signals">
                  {#each dosContent.signals as s}
                    <li>
                      <span class="s-date">{s.date ?? ''}</span>
                      <span class="s-body">
                        {#if s.kind}<span class="s-kind">{s.kind}</span>{/if}{s.body}
                        {#if s.source}<span class="s-source">
                          {#if dosSigDomain(s.source)}
                            <img class="sig-favicon" src={`https://www.google.com/s2/favicons?sz=32&domain=${dosSigDomain(s.source)}`} alt="" width="12" height="12" />
                          {/if}
                          {s.source}
                        </span>{/if}
                      </span>
                    </li>
                  {/each}
                </ul>
              </section>
            {/if}

            <!-- Questions worth asking (A's card design) -->
            {#if dosContent?.questions?.length}
              <section class="card">
                <div class="card-hd"><h3>Questions worth asking</h3></div>
                <div class="questions">
                  {#each dosContent.questions as item}
                    <div class="prep-q">
                      <div class="q">"{item.q}"</div>
                      {#if item.why}
                        <div class="why">
                          <span class="why-spark" aria-hidden="true">
                            <svg width="11" height="11" viewBox="0 0 13 13" fill="none">
                              <path d="M6.5 1.5C6.5 4.5 4.5 6.5 1.5 6.5C4.5 6.5 6.5 8.5 6.5 11.5C6.5 8.5 8.5 6.5 11.5 6.5C8.5 6.5 6.5 4.5 6.5 1.5Z" fill="currentColor"/>
                            </svg>
                          </span>
                          {item.why}
                        </div>
                      {/if}
                    </div>
                  {/each}
                </div>
              </section>
            {/if}

            <div class="prep-disclaimer">
              Synthesised from public posts, talks, and papers · {dosGeneratedAgo ? `refreshed ${dosGeneratedAgo}` : 'just generated'} · always verify before you walk in
              <button class="prep-refresh" type="button" onclick={generateDossier} disabled={generating}>
                <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M2 6a4 4 0 1 1 1.2 2.8M2 4v2h2"/></svg>
                {generating ? 'Refreshing…' : 'Refresh prep'}
              </button>
            </div>

            {#if genError}
              <p class="gen-err" style="margin-top: 16px">{genError}</p>
            {/if}
          {/if}

          <!-- Activity (full width, redesigned) -->
          <section class="card activity">
            <div class="act-hd">
              <h3>Activity</h3>
              <div class="act-actions">
                <button class="ghost-btn" onclick={openFollowUp}>+ Log a follow-up</button>
                <button class="ghost-btn" onclick={openEdit}>+ Add a note</button>
                <button class="ghost-btn" onclick={openAddEvent}>+ Log an event</button>
              </div>
            </div>
            {#if timeline.length > 0}
              <ul class="timeline">
                {#each timeline as e}
                  <li class="tl-row {e.tag}">
                    <span class="tl-rail"><span class="tl-dot"></span></span>
                    <span class="tl-date">{e.date}</span>
                    <span class="tl-body">
                      <span class="tl-title">
                        {e.title}
                        {#if e.followUp}
                          <button class="tl-del" title="Delete follow-up" aria-label="Delete follow-up" onclick={() => deleteFollowUp(e.followUp)}>
                            <svg width="12" height="12" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M3 3l8 8M11 3l-8 8" stroke-linecap="round"/></svg>
                          </button>
                        {/if}
                      </span>
                      {#if e.note}<span class="tl-note">{e.note}</span>{/if}
                    </span>
                  </li>
                {/each}
              </ul>
            {:else}
              <p style="color:var(--mute); font-size:13px; margin:0;">No activity yet.</p>
            {/if}
          </section>
        </div>

        <!-- RIGHT — meta rail -->
        <aside class="rail">
          <!-- Next interview -->
          {#if hasNext}
            <div class="next-card">
              <div class="nc-kicker"><span class="nc-dot"></span>NEXT INTERVIEW</div>
              <div class="nc-title">{nextTitle}</div>
              {#if nextRows.length}
                <div class="nc-rows">
                  {#each nextRows as r}
                    <div class="nc-row">{r}</div>
                  {/each}
                </div>
              {/if}
            </div>
          {:else}
            <div class="rail-card next-empty">
              <div class="rc-hd">Next interview</div>
              <p class="next-empty-line">No interview scheduled yet — your prep's ready for when one is.</p>
            </div>
          {/if}

          <!-- Details -->
          <div class="rail-card">
            <div class="rc-hd">Details</div>
            <dl class="kv">
              <dt>Status</dt><dd>{STATUS_LABEL[app.status]}</dd>
              <dt>Source</dt><dd>{app.source}</dd>
              {#if app.raw.location}<dt>Location</dt><dd>{app.raw.location}</dd>{/if}
              {#if app.raw.salary_note}<dt>Salary</dt><dd>{app.raw.salary_note}</dd>{/if}
              <dt>Last activity</dt><dd>{fmtRelativeDate(app.raw.updated_at ?? app.raw.applied_at)}</dd>
            </dl>
            {#if app.raw.jd_url}
              <a class="jd-link" href={app.raw.jd_url} target="_blank" rel="noopener">
                <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M5 11l6-6M6 5h5v5"/></svg>
                Open job post
              </a>
            {/if}
          </div>

          <!-- Contact -->
          <div class="rail-card">
            <div class="rc-hd">Contact</div>
            {#if app.raw.hiring_manager_name}
              <div class="contact">
                <div class="c-av">{hiringManagerInitials || '—'}</div>
                <div class="c-info">
                  <div class="c-name">{app.raw.hiring_manager_name}</div>
                  <div class="c-role">Hiring manager</div>
                </div>
              </div>
              {#if app.raw.hiring_manager_linkedin}
                <a class="c-li" href={app.raw.hiring_manager_linkedin} target="_blank" rel="noopener">
                  <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M3.5 6h2v6h-2zM4.5 3a1 1 0 1 1 0 2 1 1 0 0 1 0-2zM7 6h2v.9c.3-.5.9-1 1.8-1 1.6 0 2.2 1 2.2 2.6V12h-2V9c0-.9-.3-1.4-1.1-1.4-.6 0-1 .4-1 1.2V12H7z"/></svg>
                  LinkedIn
                </a>
              {/if}
            {:else}
              <p class="contact-empty">No hiring manager yet.</p>
              <button class="add-hm" onclick={openEdit}>
                <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.6"><path d="M8 3v10M3 8h10" stroke-linecap="round"/></svg>
                Add hiring manager
              </button>
            {/if}
          </div>
        </aside>
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
  .det { max-width: 1080px; margin: 0 auto; }

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
  .det-grid { display: grid; grid-template-columns: 1fr 320px; gap: 24px; align-items: start; }
  .left { display: flex; flex-direction: column; gap: 16px; min-width: 0; }

  /* PREP LEAD */
  .prep-lead { margin-bottom: 2px; }
  .prep-lead .ai-pill { display: inline-flex; align-items: center; gap: 6px; font-size: 12px; font-weight: 600; color: var(--accent-text); background: var(--accent-tint); border-radius: 999px; padding: 4px 11px; letter-spacing: 0.01em; }
  .prep-lead .ai-pill .spark { font-size: 12px; }
  .prep-lead h2 { font-size: 26px; font-weight: 600; letter-spacing: -0.028em; margin: 12px 0 6px; }
  .prep-gen { font-family: inherit; font-size: 13px; color: var(--mute); display: inline-flex; align-items: center; gap: 7px; margin: 0; }
  .prep-gen .sp { color: var(--accent); display: inline-flex; align-items: center; }

  /* CARD */
  .card { background: var(--card); border: 1px solid var(--rule); border-radius: 14px; padding: 20px 22px; box-shadow: var(--sh-1); }
  .card-hd { display: flex; align-items: center; gap: 10px; margin-bottom: 14px; }
  .card-hd h3 { font-size: 16px; font-weight: 600; margin: 0; letter-spacing: -0.015em; }
  .ai-tag { display: inline-flex; align-items: center; gap: 5px; font-size: 11.5px; background: var(--accent-tint); color: var(--accent-text); padding: 3px 9px; border-radius: 999px; font-weight: 500; margin-left: auto; }

  /* AI TIPS BOX (baby blue) */
  .tips { background: var(--accent-tint); border: 1px solid var(--accent-tint-2); border-radius: 13px; padding: 18px 20px; }
  .tips-hd { display: flex; align-items: center; gap: 8px; margin-bottom: 12px; }
  .tips-spark { font-size: 14px; color: var(--accent-text); }
  .tips-hd h3 { font-size: 15px; font-weight: 600; margin: 0; color: var(--accent-text); letter-spacing: -0.01em; }
  .tips-ai { margin-left: auto; font-size: 10.5px; font-weight: 600; letter-spacing: 0.04em; color: var(--accent-text); background: var(--card); border: 1px solid var(--accent-tint-2); padding: 2px 8px; border-radius: 999px; }
  .tips-list { list-style: none; margin: 0; padding: 0; display: flex; flex-direction: column; gap: 9px; }
  .tips-list li { display: grid; grid-template-columns: 14px 1fr; gap: 9px; align-items: start; font-size: 13.5px; line-height: 1.5; color: var(--accent-text); }
  .tip-dot { width: 5px; height: 5px; border-radius: 50%; background: var(--accent); margin-top: 8px; }

  /* COMPANY BRIEF — Editorial */
  .b1-hd { display: grid; grid-template-columns: 44px 1fr; gap: 14px; align-items: start; margin-bottom: 16px; }
  .b1-hd .logo { width: 44px; height: 44px; border-radius: 11px; background: var(--ink); color: #fff; display: grid; place-items: center; font-size: 19px; font-weight: 600; object-fit: contain; }
  .b1-hd .logo.letter { background: var(--surface-2); color: var(--ink); }
  .b1-name-row { display: flex; align-items: center; gap: 10px; }
  .b1-name-row h3 { font-size: 17px; font-weight: 600; letter-spacing: -0.02em; margin: 0; }
  .b1-headline { font-size: 13.5px; color: var(--mute); margin-top: 4px; line-height: 1.4; }
  .ai-pill { display: inline-flex; align-items: center; gap: 5px; font-size: 11.5px; font-weight: 600; color: var(--accent-text); background: var(--accent-tint); border: 1px solid var(--accent-tint-2); border-radius: 999px; padding: 3px 9px; letter-spacing: 0.01em; white-space: nowrap; }
  .ai-pill .spark { font-size: 11px; }
  .about { font-size: 14.5px; line-height: 1.6; color: var(--ink-2); margin: 0 0 18px; }
  .b1-facts { display: grid; grid-template-columns: repeat(4, 1fr); border: 1px solid var(--rule); border-radius: 10px; overflow: hidden; margin-bottom: 18px; }
  .b1-cell { padding: 11px 13px; }
  .b1-cell + .b1-cell { border-left: 1px solid var(--rule); }
  .f-lbl { font-size: 11px; color: var(--mute); font-weight: 500; }
  .f-val { font-size: 13px; color: var(--ink); font-weight: 600; margin-top: 3px; line-height: 1.3; }
  .proc-lbl { font-size: 11.5px; font-weight: 600; color: var(--mute); letter-spacing: 0.01em; margin-bottom: 10px; }
  .b1-process { margin-bottom: 0; }
  .b1-process .proc-chips { display: flex; align-items: center; flex-wrap: wrap; gap: 7px; }
  .chip { font-size: 12px; font-weight: 500; color: var(--ink-2); background: var(--surface-2); border: 1px solid var(--rule); border-radius: 7px; padding: 4px 10px; }
  .proc-arrow { color: var(--mute-2); font-size: 12px; }

  /* PERSON (hiring manager / interviewer) */
  .person { display: grid; grid-template-columns: 52px 1fr auto; gap: 14px; align-items: center; margin-bottom: 14px; }
  .p-av { width: 52px; height: 52px; border-radius: 50%; display: grid; place-items: center; font-weight: 600; font-size: 18px; background: var(--accent-tint); color: var(--accent-text); }
  .p-info { min-width: 0; }
  .p-name-row { display: flex; align-items: center; gap: 10px; }
  .p-info h4 { margin: 0; font-size: 17px; font-weight: 600; letter-spacing: -0.015em; }
  .role-tag { font-size: 10.5px; font-weight: 600; letter-spacing: 0.04em; text-transform: uppercase; color: var(--warm-text); background: var(--warm-tint); padding: 2px 8px; border-radius: 5px; white-space: nowrap; }
  .p-role { font-size: 13px; color: var(--mute); margin-top: 3px; }
  .p-prior { font-size: 12px; color: var(--mute-2); margin-top: 2px; }
  .p-links { display: flex; flex-wrap: wrap; gap: 6px; justify-content: flex-end; }
  .p-links a { font-size: 11.5px; color: var(--accent-text); text-decoration: none; border: 1px solid var(--rule); border-radius: 999px; padding: 4px 10px; }
  .p-links a:hover { background: var(--accent-tint); border-color: var(--accent-tint-2); }
  .snapshot { font-size: 15px; line-height: 1.55; letter-spacing: -0.008em; color: var(--ink); margin: 0 0 18px; padding-left: 14px; border-left: 2px solid var(--accent); }
  .snapshot :global(em) { font-style: normal; font-weight: 600; }
  .tells-block { border-top: 1px solid var(--rule); padding-top: 16px; }
  .tells-hd { font-size: 12.5px; font-weight: 600; color: var(--mute); margin-bottom: 10px; }
  .tells-lead { font-size: 13.5px; line-height: 1.55; color: var(--ink-2); margin: 0 0 14px; }
  .tells { display: grid; grid-template-columns: repeat(3, 1fr); gap: 14px; }
  .t-lbl { font-size: 11.5px; font-weight: 500; color: var(--mute); margin-bottom: 4px; }
  .t-val { font-size: 13.5px; color: var(--ink); line-height: 1.4; }

  /* LANDS / AVOID */
  .la-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 0; border: 1px solid var(--rule); border-radius: 14px; background: var(--card); overflow: hidden; box-shadow: var(--sh-1); }
  .la-col { padding: 20px 22px; }
  .la-col + .la-col { border-left: 1px solid var(--rule); }
  .la-col h3 { font-size: 13px; font-weight: 600; margin: 0 0 14px; display: flex; align-items: center; gap: 8px; }
  .la-col.lands h3 { color: var(--positive-text); }
  .la-col.avoid h3 { color: var(--danger-text); }
  .la-col h3 .glyph { width: 18px; height: 18px; border-radius: 5px; display: inline-flex; align-items: center; justify-content: center; font-size: 11px; font-weight: 600; }
  .la-col.lands h3 .glyph { background: var(--positive-tint); color: var(--positive-text); }
  .la-col.avoid h3 .glyph { background: var(--danger-tint); color: var(--danger-text); }
  .la-col ul { list-style: none; margin: 0; padding: 0; }
  .la-col li { font-size: 13.5px; line-height: 1.5; color: var(--ink-2); padding: 9px 0; border-top: 1px solid var(--rule); display: grid; grid-template-columns: 16px 1fr; gap: 8px; align-items: start; }
  .la-col li:first-child { border-top: none; padding-top: 0; }
  .la-col li .g { font-size: 11px; margin-top: 2px; }
  .la-col.lands li .g { color: var(--positive-text); }
  .la-col.avoid li .g { color: var(--danger-text); }

  /* SIGNALS */
  .signals { list-style: none; padding: 0; margin: 0; }
  .signals li { display: grid; grid-template-columns: 72px 1fr; gap: 14px; padding: 12px 0; border-top: 1px solid var(--rule); font-size: 13.5px; }
  .signals li:first-child { border-top: none; padding-top: 0; }
  .s-date { font-size: 12px; font-weight: 500; color: var(--mute); padding-top: 2px; }
  .s-body { color: var(--ink-2); line-height: 1.5; }
  .s-kind { display: inline-block; font-size: 11px; font-weight: 500; color: var(--mute); margin-right: 8px; padding: 1px 7px; background: var(--surface-2); border-radius: 4px; vertical-align: 1px; }
  .s-source { font-size: 12px; color: var(--accent-text); margin-left: 6px; text-decoration: none; display: inline-flex; align-items: center; gap: 4px; }
  .sig-favicon { border-radius: 2px; opacity: 0.7; }

  /* QUESTIONS — A card design */
  .questions { display: flex; flex-direction: column; gap: 12px; }
  .prep-q { background: var(--surface); border: 1px solid var(--rule); border-radius: 11px; padding: 13px 15px; }
  .prep-q .q { font-size: 14px; font-weight: 500; color: var(--ink); line-height: 1.45; }
  .prep-q .why { font-size: 12.5px; color: var(--mute); margin-top: 5px; display: flex; align-items: baseline; gap: 6px; }
  .why-spark { color: var(--accent-text); font-size: 11px; display: inline-flex; flex-shrink: 0; }

  /* Disclaimer + refresh */
  .prep-disclaimer { font-size: 11.5px; color: var(--mute); padding-top: 4px; display: flex; align-items: center; justify-content: space-between; gap: 12px; flex-wrap: wrap; }
  .prep-refresh { display: inline-flex; align-items: center; gap: 7px; background: none; color: var(--mute); border: 1px solid var(--rule); border-radius: 999px; padding: 6px 12px; font-size: 12px; font-weight: 500; cursor: pointer; font-family: inherit; transition: color 120ms, border-color 120ms, background 120ms; }
  .prep-refresh:hover:not(:disabled) { color: var(--ink); border-color: var(--rule-strong); background: var(--surface-2); }
  .prep-refresh:disabled { opacity: 0.5; cursor: default; }

  /* ACTIVITY (full width) */
  .activity { padding: 20px 22px; }
  .act-hd { display: flex; align-items: center; justify-content: space-between; gap: 14px; flex-wrap: wrap; margin-bottom: 16px; }
  .act-hd h3 { font-size: 16px; font-weight: 600; margin: 0; letter-spacing: -0.015em; }
  .act-actions { display: flex; gap: 8px; flex-wrap: wrap; }
  .ghost-btn { font-size: 12.5px; font-weight: 500; color: var(--ink-2); background: var(--surface-2); border: 1px solid var(--rule); border-radius: 8px; padding: 6px 11px; cursor: pointer; }
  .ghost-btn:hover { border-color: var(--rule-strong); color: var(--ink); }

  /* TIMELINE — connector in a per-row rail (parts design) */
  .timeline { list-style: none; margin: 0; padding: 0; }
  .tl-row { display: grid; grid-template-columns: 16px 50px 1fr; gap: 14px; align-items: start; padding: 14px 0; }
  .tl-rail { position: relative; align-self: stretch; }
  .tl-rail::before { content: ""; position: absolute; left: 50%; transform: translateX(-50%); top: 0; bottom: 0; width: 1px; background: var(--rule); }
  .tl-row:first-child .tl-rail::before { top: 9px; }
  .tl-row:last-child .tl-rail::before { bottom: auto; height: 9px; }
  .tl-dot { position: absolute; left: 50%; top: 9px; transform: translate(-50%, -50%); width: 8px; height: 8px; border-radius: 50%; background: var(--mute-2); z-index: 1; }
  .tl-row.positive .tl-dot { background: var(--positive); box-shadow: 0 0 0 3px var(--positive-tint); }
  .tl-row.accent .tl-dot { background: var(--accent); box-shadow: 0 0 0 3px var(--accent-tint); }
  .tl-row.offer .tl-dot { background: var(--positive); box-shadow: 0 0 0 3px var(--positive-tint); }
  .tl-row.danger .tl-dot { background: var(--danger); box-shadow: 0 0 0 3px var(--danger-tint); }
  .tl-date { font-size: 12.5px; color: var(--mute); line-height: 18px; font-variant-numeric: tabular-nums; }
  .tl-body { display: flex; flex-direction: column; gap: 2px; min-width: 0; }
  .tl-title { font-size: 13.5px; font-weight: 600; color: var(--ink); line-height: 18px; display: flex; align-items: center; gap: 6px; }
  .tl-note { font-size: 12.5px; color: var(--mute); line-height: 1.45; }
  .tl-del { width: 20px; height: 20px; flex-shrink: 0; border: 0; background: transparent; color: var(--mute-2); border-radius: 5px; display: inline-grid; place-items: center; cursor: pointer; opacity: 0; transition: opacity 100ms ease, background 100ms ease, color 100ms ease; }
  .tl-row:hover .tl-del { opacity: 1; }
  .tl-del:hover { background: var(--danger-tint); color: var(--danger-text); }

  /* RIGHT RAIL */
  .rail { display: flex; flex-direction: column; gap: 14px; position: sticky; top: 20px; }
  .next-card { background: var(--ink); border-radius: 14px; padding: 18px; box-shadow: var(--sh-1); }
  .nc-kicker { display: flex; align-items: center; gap: 7px; font-size: 11px; font-weight: 600; letter-spacing: 0.06em; color: rgba(255,255,255,0.7); }
  .nc-dot { width: 7px; height: 7px; border-radius: 50%; background: var(--warm); box-shadow: 0 0 0 3px rgba(255,255,255,0.12); }
  .nc-title { font-size: 16px; font-weight: 600; color: #fff; line-height: 1.35; letter-spacing: -0.01em; margin: 12px 0 14px; }
  .nc-rows { display: flex; flex-direction: column; gap: 7px; border-top: 1px solid rgba(255,255,255,0.14); padding-top: 12px; }
  .nc-row { font-size: 13px; color: rgba(255,255,255,0.85); }
  .next-empty .next-empty-line { font-size: 12.5px; color: var(--mute); line-height: 1.5; margin: 0; }

  .rail-card { background: var(--card); border: 1px solid var(--rule); border-radius: 12px; padding: 16px; box-shadow: var(--sh-1); }
  .rc-hd { font-size: 12.5px; font-weight: 600; color: var(--mute); margin-bottom: 12px; }
  .kv { margin: 0; display: grid; grid-template-columns: auto 1fr; gap: 9px 12px; }
  .kv dt { font-size: 12.5px; color: var(--mute); }
  .kv dd { margin: 0; font-size: 12.5px; color: var(--ink); font-weight: 500; text-align: right; font-variant-numeric: tabular-nums; }
  .jd-link { display: inline-flex; align-items: center; gap: 6px; margin-top: 14px; font-size: 12.5px; font-weight: 500; color: var(--accent-text); text-decoration: none; }
  .jd-link:hover { text-decoration: underline; }

  .contact { display: grid; grid-template-columns: 36px 1fr; gap: 10px; align-items: center; margin-bottom: 12px; }
  .c-av { width: 36px; height: 36px; border-radius: 50%; display: grid; place-items: center; font-weight: 600; font-size: 13px; background: var(--accent-tint); color: var(--accent-text); }
  .c-name { font-size: 13.5px; font-weight: 600; }
  .c-role { font-size: 12px; color: var(--mute); margin-top: 1px; }
  .c-li { display: inline-flex; align-items: center; gap: 6px; width: 100%; justify-content: center; background: var(--surface-2); border: 1px solid var(--rule); border-radius: 8px; padding: 7px 12px; font-size: 12.5px; font-weight: 500; color: var(--ink); text-decoration: none; box-sizing: border-box; }
  .c-li svg { color: #0a66c2; }
  .contact-empty { font-size: 12.5px; color: var(--mute); margin: 0 0 10px; }
  .add-hm { display: inline-flex; align-items: center; gap: 6px; width: 100%; justify-content: center; background: var(--surface-2); border: 1px solid var(--rule); border-radius: 8px; padding: 7px 12px; font-size: 12.5px; font-weight: 500; color: var(--ink-2); cursor: pointer; font-family: inherit; }
  .add-hm:hover { border-color: var(--rule-strong); color: var(--ink); }

  /* FOLLOW-UP MODAL */
  .fu-card { max-width: 460px; }
  .fu-fields { display: flex; flex-direction: column; gap: 14px; }
  .fu-label { display: flex; flex-direction: column; gap: 6px; font-size: 12px; color: var(--mute); }
  .fu-row { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
  .fu-ta, .fu-input { font: inherit; color: var(--ink); background: var(--surface); border: 1px solid var(--rule); border-radius: 8px; padding: 9px 11px; font-size: 13.5px; outline: none; transition: border-color 100ms ease; box-sizing: border-box; width: 100%; }
  .fu-ta { resize: vertical; line-height: 1.5; }
  .fu-ta:focus, .fu-input:focus { border-color: var(--accent); }

  /* Generate / empty state */
  .generate-card { background: var(--card); border: 1px solid var(--rule); border-radius: 18px; padding: 32px 36px; max-width: 560px; box-shadow: var(--sh-2); text-align: center; }
  .gen-icon { width: 56px; height: 56px; border-radius: 16px; background: var(--accent-tint); color: var(--accent-text); display: flex; align-items: center; justify-content: center; margin: 0 auto 20px; }
  .generate-card h3 { font-size: 20px; font-weight: 500; letter-spacing: -0.02em; margin: 0 0 10px; }
  .gen-sub { font-size: 14px; color: var(--mute); line-height: 1.6; margin: 0 auto 24px; max-width: 42ch; }
  .gen-row { display: flex; gap: 10px; flex-direction: column; }
  .gen-input { width: 100%; padding: 11px 14px; font-size: 13.5px; font-family: inherit; border: 1px solid var(--rule); border-radius: 9px; background: var(--surface-2); color: var(--ink); outline: none; transition: border-color 120ms; box-sizing: border-box; }
  .gen-input:focus { border-color: var(--accent); background: var(--card); }
  .gen-input::placeholder { color: var(--mute-2); }
  .btn-generate { background: var(--ink); color: #fff; border: none; border-radius: 9px; padding: 12px 20px; font-size: 14px; font-weight: 500; font-family: inherit; cursor: pointer; transition: background 120ms; width: 100%; }
  .btn-generate:hover:not(:disabled) { background: var(--ink-2); }
  .btn-generate:disabled { opacity: 0.5; cursor: default; }
  .gen-err { color: var(--danger-text); font-size: 13px; margin: 14px 0 0; text-align: left; }
  .big-spinner { width: 36px; height: 36px; border: 2.5px solid var(--rule-strong); border-top-color: var(--accent); border-radius: 50%; animation: prep-spin 0.75s linear infinite; margin: 24px auto 0; }
  @keyframes prep-spin { to { transform: rotate(360deg); } }

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
  .zone-title .ai-pill { font-size: 10px; padding: 1px 5px; border-radius: 4px; }
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
    .rail { position: static; }
    .b1-facts { grid-template-columns: repeat(2, 1fr); }
    .b1-cell:nth-child(3) { border-left: none; }
    .b1-cell:nth-child(n+3) { border-top: 1px solid var(--rule); }
    .tells { grid-template-columns: 1fr; }
    .la-grid { grid-template-columns: 1fr; }
    .la-col + .la-col { border-left: none; border-top: 1px solid var(--rule); }
  }
</style>
