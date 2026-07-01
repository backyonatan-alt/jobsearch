<script>
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api.js';
  import { isPreview, mockApi } from '$lib/preview-mode.js';
  import { logEvent } from '$lib/analytics.js';
  import ConfirmDialog from '$lib/ConfirmDialog.svelte';
  import {
    toDisplayApp, STATUS_LABEL, STATUSES, SOURCE_SUGGESTIONS,
    fmtLongDate, fmtRelativeDate, daysSince, isStale
  } from '$lib/app-helpers.js';

  const call = isPreview() ? mockApi : api;

  // First playbook from the prep-first cold start lands here with ?welcome=1.
  let welcomeDismissed = $state(false);
  const showWelcome = $derived(!welcomeDismissed && page.url.searchParams.get('welcome') === '1');

  // Fire-once guard for dossier_open (one event per app viewed, not per render).
  let lastDossierOpenId = null;

  let app = $state(null);
  let loading = $state(true);
  let notFound = $state(false);

  // Interview prep (dossier) — inline section
  let dossier = $state(null);
  let dossierLoading = $state(true);
  let generating = $state(false);
  let prepStage = $state('');
  let genError = $state('');
  let interviewerInput = $state('');
  // "Not them? →" re-ground: user-confirmed company website that overrides
  // same-named-company drift on regeneration.
  let showReground = $state(false);
  let companyUrlInput = $state('');
  // Which tab is showing: 'company' (shared brief) or an interview id (that
  // round's interviewer brief).
  let selectedTab = $state('company');

  // Interviews
  let interviews = $state([]);
  let interviewsLoading = $state(false);

  // Add-event flow (popup modal) — one unified input: paste / drop / type.
  let showEventModal = $state(false);
  let evText = $state('');
  let evAttach = $state(null); // { kind:'image'|'ics', name, size, mediaType?, file?, content? }
  let evDragOver = $state(false);
  let icsParsing = $state(false);
  let icsParseError = $state('');
  let icsPreview = $state([]);
  let icsSaving = $state(false);
  const EV_IMG_TYPES = ['image/png', 'image/jpeg', 'image/gif', 'image/webp'];
  const EV_IMG_MAX = 6 * 1024 * 1024;
  const looksLikeIcs = (t) => /BEGIN:VCALENDAR/i.test(t);

  // Inline-action state
  let showStatusMenu = $state(false);
  let showEditModal = $state(false);
  let edit = $state({ company: '', role: '', source: '', location: '', cv_variant: '', jd_url: '', jd_text: '', salary_note: '', hiring_manager_name: '', hiring_manager_linkedin: '', recruiter_name: '', recruiter_email: '', recruiter_linkedin: '' });
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

  // Per-app interview pipeline (ordered, user-defined stages).
  let pipeline = $state([]);
  let pipelineEditing = $state(false);
  let pipelineDraft = $state([]);
  let pipelineSaving = $state(false);
  const TYPICAL_LOOP = ['Recruiter call', 'Hiring manager', 'Take-home', 'Team interview', 'Offer'];
  const pipelineDone = $derived(pipeline.filter(s => s.done).length);

  const id = $derived(page.params.id);

  $effect(() => {
    void id;
    loadApp();
    loadFollowUps();
    initPrep();
    if (id && id !== lastDossierOpenId) {
      lastDossierOpenId = id;
      logEvent('dossier_open', { app_id: Number(id) });
    }
  });

  // Load interviews first so we can default the prep to the next upcoming round
  // (falling back to the shared Company tab when nothing's scheduled).
  async function initPrep() {
    await loadInterviews();
    selectedTab = nextRoundId ?? 'company';
    await loadDossier(selectedTab);
  }

  const onCompany = $derived(selectedTab === 'company');
  const selectedRound = $derived(onCompany ? null : (interviews || []).find(iv => iv.id === selectedTab) || null);

  // The soonest upcoming interview — the round we prep for by default.
  const nextRoundId = $derived.by(() => {
    const now = Date.now();
    const future = (interviews || [])
      .filter(iv => iv?.starts_at && new Date(iv.starts_at).getTime() >= now)
      .sort((a, b) => new Date(a.starts_at) - new Date(b.starts_at));
    return future[0]?.id ?? null;
  });

  function roundLabel(iv) {
    const d = iv?.starts_at ? fmtRelativeDate(iv.starts_at) : '';
    const s = (iv?.summary || '').trim();
    return s ? (d ? `${d} · ${s}` : s) : (d || 'Interview');
  }

  function attendeeName(iv) {
    const a = iv?.attendees;
    const arr = Array.isArray(a) ? a : (a ? [a] : []);
    const named = arr.find(x => x && typeof x === 'object' && x.name && !/recruit/i.test(x.name || ''));
    return named?.name || '';
  }

  // Switch tabs. On a round with no brief yet, pre-fill the interviewer name
  // from that round's attendees so generating is one click away.
  async function selectTab(tab) {
    if (selectedTab === tab) return;
    selectedTab = tab;
    interviewerInput = '';
    await loadDossier(tab);
    if (!dossier && tab !== 'company') {
      const iv = (interviews || []).find(x => x.id === tab);
      interviewerInput = (iv ? attendeeName(iv) : '') || app?.raw?.hiring_manager_name || '';
    }
  }

  async function loadApp() {
    loading = true;
    notFound = false;
    try {
      const raw = await call(`/api/applications/${id}`);
      app = toDisplayApp(raw);
      pipeline = Array.isArray(raw.pipeline) ? raw.pipeline : [];
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

  async function savePipeline(stages) {
    pipelineSaving = true;
    try {
      const r = await call(`/api/applications/${id}/pipeline`, { method: 'PUT', body: JSON.stringify({ stages }) });
      pipeline = Array.isArray(r.pipeline) ? r.pipeline : stages;
    } catch (e) {
      if (e.message !== 'unauthorized') console.error(e);
    } finally {
      pipelineSaving = false;
    }
  }
  function toggleStage(i) {
    const next = pipeline.map((s, idx) => idx === i ? { ...s, done: !s.done } : s);
    pipeline = next;
    savePipeline(next);
  }
  function startEditPipeline() {
    pipelineDraft = pipeline.length ? pipeline.map(s => ({ ...s })) : [{ name: '', done: false }];
    pipelineEditing = true;
  }
  function seedTypicalLoop() {
    pipelineDraft = TYPICAL_LOOP.map(name => ({ name, done: false }));
    pipelineEditing = true;
  }
  function addDraftStage() { pipelineDraft = [...pipelineDraft, { name: '', done: false }]; }
  function removeDraftStage(i) { pipelineDraft = pipelineDraft.filter((_, idx) => idx !== i); }
  function moveDraft(i, dir) {
    const j = i + dir;
    if (j < 0 || j >= pipelineDraft.length) return;
    const next = [...pipelineDraft];
    [next[i], next[j]] = [next[j], next[i]];
    pipelineDraft = next;
  }
  // Drag-to-reorder stages (arrows kept as a fallback). Reorders live on hover.
  let pipeDragIdx = $state(null);
  function onStageDragStart(i) { pipeDragIdx = i; }
  function onStageDragEnd() { pipeDragIdx = null; }
  function onStageDragOver(e, i) {
    e.preventDefault();
    if (pipeDragIdx === null || pipeDragIdx === i) return;
    const next = [...pipelineDraft];
    const [moved] = next.splice(pipeDragIdx, 1);
    next.splice(i, 0, moved);
    pipelineDraft = next;
    pipeDragIdx = i;
  }
  async function saveEditPipeline() {
    const cleaned = pipelineDraft.map(s => ({ name: s.name.trim(), done: !!s.done })).filter(s => s.name);
    await savePipeline(cleaned);
    pipelineEditing = false;
  }

  async function loadDossier(tab) {
    dossierLoading = true;
    genError = '';
    try {
      const q = (tab === 'company') ? '?scope=company' : `?interview_id=${tab}`;
      const d = await call(`/api/applications/${id}/dossier${q}`);
      dossier = d || null;
      if (dossier?.interviewer_name) interviewerInput = dossier.interviewer_name;
    } catch (e) {
      // 404 / empty → this tab hasn't been generated yet
      dossier = null;
    } finally {
      dossierLoading = false;
    }
  }

  const PREP_STAGES = [
    'Searching the web for recent posts, talks, and company news…',
    'Reading through what we found…',
    'Spotting how they tend to interview…',
    'Writing your brief…'
  ];
  async function generateDossier() {
    if (generating) return;
    generating = true;
    genError = '';
    prepStage = PREP_STAGES[0];
    // Web search + Sonnet can run ~1–2 min; cycle stages so it never looks stuck.
    const timers = [
      setTimeout(() => { if (generating) prepStage = PREP_STAGES[1]; }, 8000),
      setTimeout(() => { if (generating) prepStage = PREP_STAGES[2]; }, 25000),
      setTimeout(() => { if (generating) prepStage = PREP_STAGES[3]; }, 50000)
    ];
    try {
      const companyUrl = companyUrlInput.trim() || undefined;
      const body = onCompany
        ? { company_url: companyUrl }
        : { interview_id: selectedTab, interviewer_name: interviewerInput.trim() || undefined, company_url: companyUrl };
      const d = await call(`/api/applications/${id}/dossier/refresh`, {
        method: 'POST',
        body: JSON.stringify(body)
      });
      dossier = d;
      interviewerInput = d.interviewer_name ?? interviewerInput;
      showReground = false;
    } catch (e) {
      genError = friendlyGenErr(e.message);
    } finally {
      timers.forEach(clearTimeout);
      generating = false;
      prepStage = '';
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
    return m || 'Could not build the playbook.';
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

  // ── Add-event flow — one box that takes an .ics file, a screenshot, or text ──
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
  async function setEvFile(f) {
    if (!f) return;
    icsParseError = '';
    const name = (f.name || '').toLowerCase();
    if (EV_IMG_TYPES.includes(f.type)) {
      if (f.size > EV_IMG_MAX) { icsParseError = 'Screenshot too large (6 MB max).'; return; }
      evAttach = { kind: 'image', name: f.name || 'pasted.png', size: f.size, mediaType: f.type, file: f };
      return;
    }
    if (name.endsWith('.ics') || f.type === 'text/calendar') {
      if (f.size > 256 * 1024) { icsParseError = 'Calendar file too large (256 KB max).'; return; }
      evAttach = { kind: 'ics', name: f.name || 'invite.ics', size: f.size, content: await f.text() };
      return;
    }
    icsParseError = 'Drop a screenshot (PNG/JPEG) or an .ics calendar file.';
  }
  function onEvDrop(e) { e.preventDefault(); evDragOver = false; setEvFile(e.dataTransfer?.files?.[0]); }
  function onEvDragOver(e) { if (e.dataTransfer?.types?.includes('Files')) { e.preventDefault(); evDragOver = true; } }
  function onEvPaste(e) {
    const item = [...(e.clipboardData?.items || [])].find(i => i.type.startsWith('image/'));
    if (item) { e.preventDefault(); setEvFile(item.getAsFile()); }
  }
  async function onEvFileInput(e) { const f = e.target.files?.[0]; e.target.value = ''; await setEvFile(f); }

  async function parseEvent() {
    if (icsParsing) return;
    icsParseError = '';
    icsPreview = [];
    const text = evText.trim();
    if (!evAttach && !text) { icsParseError = 'Drop a file or screenshot, or paste the invite text first.'; return; }
    icsParsing = true;
    try {
      let payload;
      if (evAttach?.kind === 'ics') {
        payload = { ics: evAttach.content };
      } else if (evAttach?.kind === 'image') {
        payload = { image: { media_type: evAttach.mediaType, data: await fileToBase64(evAttach.file) } };
        if (text) payload.text = text;
      } else if (looksLikeIcs(text)) {
        payload = { ics: text };
      } else {
        payload = { text };
      }
      // So a bare "2:30pm" resolves to the user's wall clock, not US Eastern.
      try { payload.tz = Intl.DateTimeFormat().resolvedOptions().timeZone; } catch {}
      const r = await call(`/api/applications/${id}/interviews/parse`, { method: 'POST', body: JSON.stringify(payload) });
      icsPreview = r.events ?? [];
      if (icsPreview.length === 0) icsParseError = "Couldn't find an event in that — try a screenshot or the email body.";
    } catch (e) {
      icsParseError = e.message || 'Could not parse.';
    } finally {
      icsParsing = false;
    }
  }
  async function saveParsedEvents() {
    if (icsSaving || icsPreview.length === 0) return;
    icsSaving = true;
    try {
      for (const ev of icsPreview) {
        await call(`/api/applications/${id}/interviews`, { method: 'POST', body: JSON.stringify(ev) });
      }
      evText = ''; evAttach = null;
      icsPreview = [];
      showEventModal = false;
      await loadInterviews();
      // A new round shifts the default prep — point at it and load its (empty)
      // prep so the user lands on "generate for this round". Tell Today to refetch.
      if (nextRoundId) await selectTab(nextRoundId);
      try { window.dispatchEvent(new CustomEvent('pursuit:refresh')); } catch {}
    } catch (e) {
      icsParseError = e.message || 'Could not save events.';
    } finally {
      icsSaving = false;
    }
  }
  // In-app confirmation for destructive actions (replaces window.confirm).
  let confirmDlg = $state({ open: false, title: '', message: '', confirmLabel: 'Delete', busy: false, action: null });
  function askConfirm(opts) {
    confirmDlg = { open: true, busy: false, confirmLabel: 'Delete', ...opts };
  }
  async function runConfirm() {
    if (!confirmDlg.action || confirmDlg.busy) return;
    confirmDlg.busy = true;
    try {
      await confirmDlg.action();
      confirmDlg.open = false;
    } catch (e) {
      if (e.message !== 'unauthorized') console.error(e);
    } finally {
      confirmDlg.busy = false;
    }
  }

  function deleteInterview(iv) {
    askConfirm({
      title: 'Delete this event?',
      message: `"${iv.summary}" will be removed from this application.`,
      confirmLabel: 'Delete event',
      action: async () => {
        await call(`/api/applications/${id}/interviews/${iv.id}`, { method: 'DELETE' });
        await loadInterviews();
        if (selectedTab === iv.id) { selectedTab = nextRoundId ?? 'company'; await loadDossier(selectedTab); }
        try { window.dispatchEvent(new CustomEvent('pursuit:refresh')); } catch {}
      }
    });
  }

  function openAddEvent() {
    showEventModal = true;
    icsParseError = '';
    icsPreview = [];
    evText = '';
    evAttach = null;
    // Intent signal: opened the add-interview flow. The drop to a saved
    // interview tells us if the flow we just fixed is actually converting.
    logEvent('addevent_open', { app_id: Number(id) });
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
    const fromStatus = app.status;
    await call(`/api/applications/${id}`, { method: 'PATCH', body: JSON.stringify({ status: newStatus }) });
    logEvent('status_change', { from: fromStatus, to: newStatus, surface: 'detail' });
    await loadApp();
    try { window.dispatchEvent(new CustomEvent('pursuit:refresh')); } catch {}
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
      jd_text:                 app.raw.jd_text ?? '',
      salary_note:             app.raw.salary_note ?? '',
      hiring_manager_name:     app.raw.hiring_manager_name ?? '',
      hiring_manager_linkedin: app.raw.hiring_manager_linkedin ?? '',
      recruiter_name:          app.raw.recruiter_name ?? '',
      recruiter_email:         app.raw.recruiter_email ?? '',
      recruiter_linkedin:      app.raw.recruiter_linkedin ?? ''
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
      try { window.dispatchEvent(new CustomEvent('pursuit:refresh')); } catch {}
    } finally {
      saving = false;
    }
  }
  function deleteApp() {
    if (!app) return;
    askConfirm({
      title: `Delete the ${app.co} application?`,
      message: "This removes the application and everything attached to it — events, contacts, pipeline, and prep. This can't be undone.",
      confirmLabel: 'Delete application',
      action: async () => {
        await call(`/api/applications/${id}`, { method: 'DELETE' });
        try { window.dispatchEvent(new CustomEvent('pursuit:refresh')); } catch {}
        goto('/app', { replaceState: true });
      }
    });
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
  const recruiterInitials = $derived(initialsOf(app?.raw?.recruiter_name));

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

  // Company identity (disambiguation) + sources (citations) — from the brief.
  const dosIdentity = $derived(dosContent?.identity ?? null);
  const dosSources = $derived(Array.isArray(dosContent?.sources) ? dosContent.sources.filter(s => s?.href) : []);

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

  // Preview formatters: the weekday + year are always computed from the parsed
  // date here, never taken from the model's prose, so a wrong day is visible.
  function fmtEventDay(ev) {
    const d = new Date(ev.starts_at);
    return d.toLocaleDateString(undefined, { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' });
  }
  function fmtEventTimeSuffix(ev) {
    const d = new Date(ev.starts_at);
    if (ev.all_day) return ' · all day';
    let s = ' · ' + d.toLocaleTimeString(undefined, { hour: 'numeric', minute: '2-digit' });
    if (ev.ends_at) {
      const mins = Math.round((new Date(ev.ends_at) - d) / 60000);
      if (mins > 0) s += ` · ${mins >= 60 && mins % 60 === 0 ? `${mins / 60}h` : `${mins} min`}`;
    }
    return s;
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
        tag: 'accent',
        interview: iv
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

      {#if showWelcome}
        <div class="welcome-banner">
          <span class="wb-spark">✦</span>
          <span class="wb-tx">Here's your first playbook. Add who's interviewing you for round-by-round prep, or <a href="/app">track another application</a>.</span>
          <button class="wb-x" onclick={() => (welcomeDismissed = true)} aria-label="Dismiss">✕</button>
        </div>
      {/if}

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
            <h2>Interview playbook</h2>
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

          <div class="round-tabs" role="tablist" aria-label="Interview round">
            <button
              type="button" role="tab"
              aria-selected={onCompany}
              class="round-tab company" class:active={onCompany}
              onclick={() => selectTab('company')}
            >
              <svg width="12" height="12" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5" aria-hidden="true"><rect x="2.5" y="2" width="11" height="12" rx="1"/><path d="M5.5 5h2M5.5 8h2M5.5 11h2M9.5 5h1.5M9.5 8h1.5"/></svg>
              Company
            </button>
            {#each interviews as iv}
              <button
                type="button" role="tab"
                aria-selected={selectedTab === iv.id}
                class="round-tab" class:active={selectedTab === iv.id}
                onclick={() => selectTab(iv.id)}
              >{roundLabel(iv)}</button>
            {/each}
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
                <h3>Researching {app.co}{!onCompany && interviewerInput ? ` & ${interviewerInput}` : ''}…</h3>
                <p class="gen-sub" aria-live="polite">{prepStage || PREP_STAGES[0]}</p>
                <div class="big-spinner"></div>
                <p class="gen-eta">This usually takes 1–2 minutes — you can keep working, it'll be here when it's done.</p>
              {:else if onCompany}
                <h3>Generate company brief</h3>
                <p class="gen-sub">
                  A shared brief on {app.co} — what they do, where they're headed, the typical interview loop, and what this team grades for. Researched once and used across every round.
                </p>
                <div class="gen-row">
                  <button class="btn-generate" onclick={generateDossier} disabled={generating}>
                    Generate company brief
                  </button>
                </div>
                {#if genError}<p class="gen-err">{genError}</p>{/if}
              {:else}
                <h3>Playbook for {selectedRound ? roundLabel(selectedRound) : 'this round'}</h3>
                <p class="gen-sub">
                  We'll research the person interviewing you in this round — their background, how they tend to interview, what lands, and smart questions to ask. The shared {app.co} company brief is generated alongside it if you don't have one yet, so you only wait once.
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
                    Build playbook
                  </button>
                </div>
                {#if genError}<p class="gen-err">{genError}</p>{/if}
              {/if}
            </div>

          {:else}
            <!-- Full brief -->

            <!-- Researched-the-right-company assurance + re-ground (disambiguation) -->
            {#if dosIdentity}
              <div class="researched">
                <div class="rs-main">
                  <span class="rs-lbl">Researched</span>
                  <strong class="rs-name">{dosIdentity.name || app.co}</strong>
                  {#if dosIdentity.domain}
                    <a class="rs-dom" href={`https://${String(dosIdentity.domain).replace(/^https?:\/\//, '')}`} target="_blank" rel="noreferrer">{dosIdentity.domain}</a>
                  {/if}
                  {#if dosIdentity.summary}<span class="rs-sum">— {dosIdentity.summary}</span>{/if}
                </div>
                <button class="rs-not" type="button" onclick={() => { showReground = !showReground; companyUrlInput = ''; }}>
                  {showReground ? 'Cancel' : 'Not them?'}
                </button>
              </div>
              {#if showReground}
                <div class="reground">
                  <p>Paste the company's website so we research the right one — this regenerates the playbook.</p>
                  <div class="rg-row">
                    <input class="gen-input" type="url" placeholder="https://company.com" bind:value={companyUrlInput}
                      onkeydown={(e) => e.key === 'Enter' && companyUrlInput.trim() && generateDossier()} disabled={generating} />
                    <button class="btn-generate" onclick={generateDossier} disabled={generating || !companyUrlInput.trim()}>
                      {generating ? 'Re-grounding…' : 'Re-ground'}
                    </button>
                  </div>
                </div>
              {/if}
            {/if}

            {#if onCompany}
            <!-- What this team grades for (company watch-fors) -->
            {#if tips.length}
              <section class="tips">
                <div class="tips-hd">
                  <span class="tips-spark">✦</span>
                  <h3>What this team grades for</h3>
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
                {#if dosSources.length}
                  <div class="b1-sources">
                    <div class="proc-lbl">Sources</div>
                    <ul>
                      {#each dosSources as s}
                        <li><a href={s.href} target="_blank" rel="noreferrer">
                          <img class="src-favicon" src={`https://www.google.com/s2/favicons?sz=32&domain=${dosSigDomain(s.href)}`} alt="" width="12" height="12" />
                          {s.label || dosSigDomain(s.href)}
                        </a></li>
                      {/each}
                    </ul>
                  </div>
                {/if}
              </section>
            {/if}

            {:else}

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
                </div>
                {#if dosInterviewer?.links?.length}
                  <div class="p-links">
                    {#each dosInterviewer.links as l}
                      <a href={l.href} target="_blank" rel="noopener">{l.label}</a>
                    {/each}
                  </div>
                {/if}
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
                        {#if s.source || s.source_url}
                          {#if s.source_url}
                            <a class="s-source" href={s.source_url} target="_blank" rel="noreferrer">
                              <img class="sig-favicon" src={`https://www.google.com/s2/favicons?sz=32&domain=${dosSigDomain(s.source_url || s.source)}`} alt="" width="12" height="12" />
                              {s.source || dosSigDomain(s.source_url)}
                            </a>
                          {:else}
                            <span class="s-source">
                              {#if dosSigDomain(s.source)}
                                <img class="sig-favicon" src={`https://www.google.com/s2/favicons?sz=32&domain=${dosSigDomain(s.source)}`} alt="" width="12" height="12" />
                              {/if}
                              {s.source}
                            </span>
                          {/if}
                        {/if}
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
            {/if}

            <div class="prep-disclaimer">
              {#if onCompany}
                Shared across every round · {dosGeneratedAgo ? `refreshed ${dosGeneratedAgo}` : 'just generated'}
              {:else}
                Synthesised from public posts, talks, and papers · {dosGeneratedAgo ? `refreshed ${dosGeneratedAgo}` : 'just generated'} · always verify before you walk in
              {/if}
              <button class="prep-refresh" type="button" onclick={generateDossier} disabled={generating}>
                <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M2 6a4 4 0 1 1 1.2 2.8M2 4v2h2"/></svg>
                {generating ? 'Refreshing…' : (onCompany ? 'Refresh company brief' : 'Refresh playbook')}
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
                <button class="ghost-btn" onclick={openAddEvent}>+ Add interview</button>
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
                        {:else if e.interview}
                          <button class="tl-del" title="Delete interview" aria-label="Delete interview" onclick={() => deleteInterview(e.interview)}>
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
              <button class="next-add-btn" onclick={openAddEvent}>+ Add interview</button>
            </div>
          {/if}

          <!-- Pipeline -->
          <div class="rail-card">
            <div class="rc-hd rc-hd-row">
              <span>Pipeline</span>
              {#if pipeline.length && !pipelineEditing}
                <button class="rc-edit" onclick={startEditPipeline}>Edit</button>
              {/if}
            </div>

            {#if pipelineEditing}
              <div class="pipe-edit">
                {#each pipelineDraft as st, i (i)}
                  <div class="pe-row" class:pe-dragging={pipeDragIdx === i} ondragover={(e) => onStageDragOver(e, i)}>
                    <span class="pe-grip" draggable="true" ondragstart={() => onStageDragStart(i)} ondragend={onStageDragEnd} title="Drag to reorder" aria-label="Drag to reorder">⠿</span>
                    <input class="pe-input" bind:value={st.name} placeholder="Stage name" />
                    <button class="pe-btn" onclick={() => moveDraft(i, -1)} disabled={i === 0} aria-label="Move up">↑</button>
                    <button class="pe-btn" onclick={() => moveDraft(i, 1)} disabled={i === pipelineDraft.length - 1} aria-label="Move down">↓</button>
                    <button class="pe-btn pe-x" onclick={() => removeDraftStage(i)} aria-label="Remove stage">×</button>
                  </div>
                {/each}
                <button class="add-hm" onclick={addDraftStage}>
                  <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.6"><path d="M8 3v10M3 8h10" stroke-linecap="round"/></svg>
                  Add stage
                </button>
                <div class="pe-actions">
                  <button class="btn" onclick={() => (pipelineEditing = false)}>Cancel</button>
                  <button class="btn btn-primary" onclick={saveEditPipeline} disabled={pipelineSaving}>{pipelineSaving ? 'Saving…' : 'Save'}</button>
                </div>
              </div>
            {:else if pipeline.length}
              <div class="pl-prog">{pipelineDone} of {pipeline.length} done</div>
              <ol class="pipe">
                {#each pipeline as st, i}
                  <li class="pipe-step" class:done={st.done}>
                    <button class="pipe-dot" onclick={() => toggleStage(i)} aria-label={st.done ? 'Mark not done' : 'Mark done'}>
                      {#if st.done}<svg width="10" height="10" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M2.5 6.5l2.5 2.5 4.5-5"/></svg>{/if}
                    </button>
                    <span class="pipe-name">{st.name}</span>
                  </li>
                {/each}
              </ol>
            {:else}
              <p class="contact-empty">No stages yet — map the steps the recruiter described.</p>
              <button class="add-hm" onclick={seedTypicalLoop}>
                <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.6"><path d="M8 3v10M3 8h10" stroke-linecap="round"/></svg>
                Start from a typical loop
              </button>
            {/if}
          </div>

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
            {#if app.raw.jd_text}
              <details class="jd-saved">
                <summary>
                  <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M4 2h6l3 3v9H4z"/><path d="M9 2v4h4"/></svg>
                  Saved job description
                </summary>
                <p class="jd-body">{app.raw.jd_text}</p>
              </details>
            {/if}
          </div>

          <!-- Contact -->
          <div class="rail-card">
            <div class="rc-hd">Contacts</div>
            {#if app.raw.recruiter_name}
              <div class="contact">
                <div class="c-av c-av-warm">{recruiterInitials || '—'}</div>
                <div class="c-info">
                  <div class="c-name">{app.raw.recruiter_name}</div>
                  <div class="c-role">Recruiter / contact</div>
                </div>
              </div>
              <div class="c-links">
                {#if app.raw.recruiter_email}
                  <a class="c-li" href={`mailto:${app.raw.recruiter_email}`}>
                    <svg width="14" height="14" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.4"><rect x="2" y="3.5" width="12" height="9" rx="1.5"/><path d="M2.5 4.5L8 9l5.5-4.5"/></svg>
                    Email
                  </a>
                {/if}
                {#if app.raw.recruiter_linkedin}
                  <a class="c-li" href={app.raw.recruiter_linkedin} target="_blank" rel="noopener">
                    <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M3.5 6h2v6h-2zM4.5 3a1 1 0 1 1 0 2 1 1 0 0 1 0-2zM7 6h2v.9c.3-.5.9-1 1.8-1 1.6 0 2.2 1 2.2 2.6V12h-2V9c0-.9-.3-1.4-1.1-1.4-.6 0-1 .4-1 1.2V12H7z"/></svg>
                    LinkedIn
                  </a>
                {/if}
              </div>
            {/if}
            {#if app.raw.hiring_manager_name}
              <div class="contact" class:contact-stacked={app.raw.recruiter_name}>
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
            {/if}
            {#if !app.raw.recruiter_name && !app.raw.hiring_manager_name}
              <p class="contact-empty">No contacts yet — add the recruiter or hiring manager.</p>
              <button class="add-hm" onclick={openEdit}>
                <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.6"><path d="M8 3v10M3 8h10" stroke-linecap="round"/></svg>
                Add a contact
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
        <label>Source <input bind:value={edit.source} list="edit-source-suggestions" placeholder="LinkedIn / Referral / Cold email" /></label>
        <datalist id="edit-source-suggestions">
          {#each SOURCE_SUGGESTIONS as s}<option value={s}></option>{/each}
        </datalist>
        <label>Location <input bind:value={edit.location} placeholder="Remote / San Francisco" /></label>
        <label>CV variant <input bind:value={edit.cv_variant} placeholder="v3-ai-focus" /></label>
        <label>Salary note <input bind:value={edit.salary_note} placeholder="$220k-$280k base" /></label>
        <label class="span-2">JD URL <input bind:value={edit.jd_url} placeholder="https://…" /></label>
        <label class="span-2">Job description
          <textarea class="jd-area" bind:value={edit.jd_text} rows="4" placeholder="Paste the full JD text so it's kept even if the posting comes down."></textarea>
        </label>
        <div class="field-group span-2">Hiring manager <span class="fg-sub">— who the role reports to</span></div>
        <label>Name <input bind:value={edit.hiring_manager_name} placeholder="Jane Doe" /></label>
        <label>LinkedIn <input bind:value={edit.hiring_manager_linkedin} placeholder="https://linkedin.com/in/…" /></label>
        <div class="field-group span-2">Recruiter / contact <span class="fg-sub">— who's running your process</span></div>
        <label>Name <input bind:value={edit.recruiter_name} placeholder="Sam Levi" /></label>
        <label>Email <input bind:value={edit.recruiter_email} type="email" placeholder="sam@company.com" /></label>
        <label class="span-2">LinkedIn <input bind:value={edit.recruiter_linkedin} placeholder="https://linkedin.com/in/…" /></label>
      </div>
      <p class="privacy-note">
        <svg width="12" height="12" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5" aria-hidden="true"><rect x="3" y="7" width="10" height="6.5" rx="1.5"/><path d="M5.5 7V5a2.5 2.5 0 0 1 5 0v2"/></svg>
        Private to your account. Your notes and salary info are never shared or shown to anyone else.
      </p>
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
        <p>Paste the invite, drop a screenshot or .ics file, or just type the details — we'll pull out the event.</p>
      </div>

      <div
        class="ev-input"
        class:drag={evDragOver}
        class:loading={icsParsing}
        ondragover={onEvDragOver}
        ondragleave={() => (evDragOver = false)}
        ondrop={onEvDrop}
        role="presentation"
      >
        {#if evAttach}
          <div class="ev-attached">
            <span class="ev-att-ic">
              {#if evAttach.kind === 'image'}
                <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="2" y="3" width="12" height="10" rx="1.5"/><circle cx="5.5" cy="6.5" r="1"/><path d="M2 11l3.5-3.5 3 3 2-2L14 11"/></svg>
              {:else}
                <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="2" y="3" width="12" height="11" rx="1.5"/><path d="M2 6h12M6 2v2M10 2v2"/></svg>
              {/if}
            </span>
            <span class="ev-att-name">{evAttach.name}</span>
            <span class="ev-att-kind">{evAttach.kind === 'image' ? 'screenshot' : 'calendar file'} · {Math.round(evAttach.size / 1024)} KB</span>
            <button type="button" class="ev-att-x" onclick={() => (evAttach = null)} aria-label="Remove" disabled={icsParsing}>×</button>
          </div>
        {:else}
          <textarea
            class="ev-ta"
            rows="3"
            bind:value={evText}
            onpaste={onEvPaste}
            placeholder={"Paste an invite, or type it — e.g. “Interview Wed Jun 10, 11:00, Google Meet”"}
            disabled={icsParsing}
          ></textarea>
        {/if}

        <div class="ev-foot">
          <label class="ev-browse">
            <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M8 11V3M5 6l3-3 3 3M3 11v2h10v-2"/></svg>
            <span>Drop a screenshot or .ics, paste, or <u>browse</u></span>
            <input type="file" accept=".ics,text/calendar,image/png,image/jpeg,image/gif,image/webp" onchange={onEvFileInput} hidden />
          </label>
          {#if icsParsing}
            <span class="ev-loading"><span class="ev-spin" aria-hidden="true"></span> Reading the event…</span>
          {/if}
        </div>
      </div>

      <div class="ev-actions">
        <span class="ev-hint"><kbd>⌘V</kbd> pastes a screenshot</span>
        <button class="btn btn-primary" onclick={parseEvent} disabled={icsParsing || (!evText.trim() && !evAttach)}>
          {icsParsing ? 'Reading…' : 'Find the event'}
        </button>
      </div>

      {#if icsParseError}<p class="dossier-err" style="margin-top: 14px">{icsParseError}</p>{/if}

      {#if icsPreview.length > 0}
        <div class="ics-preview">
          <h4>Preview</h4>
          <p class="prev-check">Double-check the day and time before saving.</p>
          {#each icsPreview as ev}
            <div class="prev-row">
              <div class="prev-summary">{ev.summary || 'Untitled event'}</div>
              <div class="prev-when"><strong>{fmtEventDay(ev)}</strong>{fmtEventTimeSuffix(ev)}</div>
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

<ConfirmDialog
  open={confirmDlg.open}
  title={confirmDlg.title}
  message={confirmDlg.message}
  confirmLabel={confirmDlg.confirmLabel}
  busy={confirmDlg.busy}
  onConfirm={runConfirm}
  onCancel={() => (confirmDlg.open = false)}
/>

<style>
  .body { padding: 28px; }
  .det { max-width: 1080px; margin: 0 auto; }

  .welcome-banner { display: flex; align-items: center; gap: 11px; margin-bottom: 22px;
    padding: 12px 14px; border-radius: 12px; background: var(--accent-tint); border: 1px solid var(--accent-tint-2); }
  .welcome-banner .wb-spark { color: var(--accent-text); font-size: 14px; flex-shrink: 0; }
  .welcome-banner .wb-tx { font-size: 13.5px; color: var(--ink-2); line-height: 1.5; }
  .welcome-banner .wb-tx a { color: var(--accent-text); font-weight: 500; }
  .welcome-banner .wb-x { margin-left: auto; flex-shrink: 0; background: none; border: none; color: var(--mute);
    font-size: 13px; cursor: pointer; padding: 4px 6px; border-radius: 6px; }
  .welcome-banner .wb-x:hover { background: var(--card); color: var(--ink-2); }

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
  /* closed = neutral terminal (req cancelled), not a rejection — muted, not red. */
  .pill.closed { background: var(--surface-2); color: var(--mute); }
  .pill.closed .pdot { background: var(--mute-2); }

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

  .round-tabs { display: flex; flex-wrap: wrap; gap: 8px; margin: 14px 0 18px; }
  .round-tab { display: inline-flex; align-items: center; gap: 6px; font-family: inherit; font-size: 12.5px; font-weight: 500; color: var(--mute); background: var(--card); border: 1px solid var(--rule); border-radius: 999px; padding: 6px 13px; cursor: pointer; white-space: nowrap; transition: background .12s, color .12s, border-color .12s; }
  .round-tab:hover { color: var(--ink); border-color: var(--accent-tint-2); }
  .round-tab.active { color: var(--accent-text); background: var(--accent-tint); border-color: var(--accent-tint-2); font-weight: 600; }
  .round-tab.company { border-style: dashed; }
  .round-tab.company.active { border-style: solid; }

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

  /* Researched-the-right-company assurance + re-ground */
  .researched { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; padding: 9px 13px; margin-bottom: 14px;
    background: var(--surface-2); border: 1px solid var(--rule); border-radius: 10px; }
  .researched .rs-main { display: flex; align-items: baseline; gap: 7px; flex-wrap: wrap; min-width: 0; font-size: 12.5px; color: var(--ink-2); }
  .researched .rs-lbl { font-size: 10.5px; font-weight: 700; letter-spacing: 0.06em; text-transform: uppercase; color: var(--mute-2); }
  .researched .rs-name { color: var(--ink); font-weight: 600; }
  .researched .rs-dom { color: var(--accent-text); text-decoration: none; }
  .researched .rs-dom:hover { text-decoration: underline; }
  .researched .rs-sum { color: var(--mute); }
  .researched .rs-not { margin-left: auto; flex-shrink: 0; background: none; border: 1px solid var(--rule); color: var(--ink-2);
    font: 500 12px/1 var(--sans); padding: 6px 11px; border-radius: 7px; cursor: pointer; }
  .researched .rs-not:hover { border-color: var(--rule-strong); background: var(--card); }
  .reground { margin: -6px 0 16px; padding: 12px 13px; background: var(--accent-tint); border: 1px solid var(--accent-tint-2); border-radius: 10px; }
  .reground p { margin: 0 0 9px; font-size: 12.5px; color: var(--ink-2); }
  .reground .rg-row { display: flex; gap: 8px; align-items: stretch; }
  .reground .rg-row .gen-input { flex: 1 1 auto; min-width: 0; }
  .reground .rg-row .btn-generate { flex: 0 0 auto; width: auto; white-space: nowrap; }

  /* Company sources (citations) */
  .b1-sources { margin-top: 16px; }
  .b1-sources ul { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 6px; }
  .b1-sources a { display: inline-flex; align-items: center; gap: 6px; font-size: 12.5px; color: var(--accent-text); text-decoration: none; }
  .b1-sources a:hover { text-decoration: underline; }
  .src-favicon { border-radius: 3px; flex-shrink: 0; }

  /* PERSON (hiring manager / interviewer) */
  .person { display: grid; grid-template-columns: 52px 1fr; gap: 14px; align-items: center; margin-bottom: 12px; }
  .p-av { width: 52px; height: 52px; border-radius: 50%; display: grid; place-items: center; font-weight: 600; font-size: 18px; background: var(--accent-tint); color: var(--accent-text); }
  .p-info { min-width: 0; }
  .p-name-row { display: flex; align-items: center; gap: 10px; }
  .p-info h4 { margin: 0; font-size: 17px; font-weight: 600; letter-spacing: -0.015em; }
  .role-tag { font-size: 10.5px; font-weight: 600; letter-spacing: 0.04em; text-transform: uppercase; color: var(--warm-text); background: var(--warm-tint); padding: 2px 8px; border-radius: 5px; white-space: nowrap; }
  .p-role { font-size: 13px; color: var(--mute); margin-top: 3px; }
  .p-prior { font-size: 12px; color: var(--mute-2); margin-top: 2px; }
  .p-links { display: flex; flex-wrap: wrap; gap: 6px; justify-content: flex-start; margin-bottom: 14px; }
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
  .next-add-btn {
    margin-top: 12px; width: 100%; padding: 8px 12px;
    font: inherit; font-size: 13px; font-weight: 500; cursor: pointer;
    color: var(--accent-text); background: var(--accent-tint);
    border: 1px solid transparent; border-radius: 8px;
  }
  .next-add-btn:hover { border-color: var(--accent); }

  .rail-card { background: var(--card); border: 1px solid var(--rule); border-radius: 12px; padding: 16px; box-shadow: var(--sh-1); }
  .rc-hd { font-size: 12.5px; font-weight: 600; color: var(--mute); margin-bottom: 12px; }
  .rc-hd-row { display: flex; align-items: center; justify-content: space-between; }
  .rc-edit { background: transparent; border: 0; color: var(--accent-text); font-size: 12px; font-weight: 500; cursor: pointer; padding: 0; font-family: inherit; }
  .rc-edit:hover { text-decoration: underline; }

  /* Pipeline stepper */
  .pl-prog { font-size: 12px; color: var(--mute); margin: -4px 0 12px; }
  .pipe { list-style: none; margin: 0; padding: 0; }
  .pipe-step { position: relative; display: grid; grid-template-columns: 22px 1fr; gap: 10px; align-items: center; padding-bottom: 14px; }
  .pipe-step:not(:last-child)::before { content: ''; position: absolute; left: 10px; top: 22px; bottom: 0; width: 1.5px; background: var(--rule); }
  .pipe-step.done:not(:last-child)::before { background: var(--positive, oklch(0.65 0.14 152)); }
  .pipe-dot { position: relative; z-index: 1; width: 22px; height: 22px; border-radius: 50%; border: 1.5px solid var(--rule-strong); background: var(--card); display: grid; place-items: center; cursor: pointer; color: white; padding: 0; transition: background 100ms ease, border-color 100ms ease; }
  .pipe-dot:hover { border-color: var(--accent); }
  .pipe-step.done .pipe-dot { background: var(--positive, oklch(0.65 0.14 152)); border-color: var(--positive, oklch(0.65 0.14 152)); }
  .pipe-name { font-size: 13px; color: var(--ink-2); }
  .pipe-step.done .pipe-name { color: var(--mute); text-decoration: line-through; }

  /* Pipeline edit mode */
  .pipe-edit { display: flex; flex-direction: column; gap: 8px; }
  .pe-row { display: grid; grid-template-columns: auto 1fr auto auto auto; gap: 4px; align-items: center; border-radius: 8px; }
  .pe-row.pe-dragging { opacity: 0.5; background: var(--accent-tint); }
  .pe-grip { display: grid; place-items: center; width: 20px; height: 30px; color: var(--mute-2); cursor: grab; font-size: 14px; user-select: none; }
  .pe-grip:active { cursor: grabbing; }
  .pe-grip:hover { color: var(--ink-2); }
  .pe-input { font: inherit; font-size: 13px; color: var(--ink); background: var(--surface); border: 1px solid var(--rule); border-radius: 7px; padding: 6px 9px; outline: none; min-width: 0; }
  .pe-input:focus { border-color: var(--accent); }
  .pe-btn { width: 26px; height: 30px; display: grid; place-items: center; background: var(--surface-2); border: 1px solid var(--rule); border-radius: 7px; color: var(--mute); font-size: 13px; cursor: pointer; font-family: inherit; }
  .pe-btn:hover:not(:disabled) { border-color: var(--rule-strong); color: var(--ink); }
  .pe-btn:disabled { opacity: 0.4; cursor: default; }
  .pe-x:hover:not(:disabled) { color: var(--danger-text); border-color: var(--danger-tint); }
  .pe-actions { display: flex; justify-content: flex-end; gap: 8px; margin-top: 4px; }
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
  .c-av-warm { background: var(--warm-tint, var(--accent-tint)); color: var(--warm-text, var(--accent-text)); }
  .c-links { display: flex; gap: 8px; margin-bottom: 4px; }
  .c-links .c-li { width: auto; flex: 1; }
  .contact-stacked { padding-top: 12px; border-top: 1px solid var(--rule); }

  .jd-saved { margin-top: 12px; }
  .jd-saved summary { display: inline-flex; align-items: center; gap: 6px; font-size: 12.5px; font-weight: 500; color: var(--accent-text); cursor: pointer; list-style: none; }
  .jd-saved summary::-webkit-details-marker { display: none; }
  .jd-saved summary:hover { text-decoration: underline; }
  .jd-body { margin: 10px 0 0; font-size: 12.5px; line-height: 1.55; color: var(--ink-2); white-space: pre-wrap; max-height: 280px; overflow-y: auto; padding: 10px 12px; background: var(--surface); border: 1px solid var(--rule); border-radius: 8px; }

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
  .gen-eta { font-size: 12px; color: var(--mute-2); margin: 16px auto 0; max-width: 40ch; }
  @keyframes prep-spin { to { transform: rotate(360deg); } }

  /* ADD-EVENT MODAL */
  .ev-overlay { position: fixed; inset: 0; background: rgba(10,10,13,0.55); backdrop-filter: blur(8px); -webkit-backdrop-filter: blur(8px); display: grid; place-items: center; z-index: 200; padding: 24px; overflow-y: auto; }
  .ev-card { position: relative; width: 100%; max-width: 540px; background: var(--card); border: 1px solid var(--rule); border-radius: 18px; padding: 26px 28px 24px; box-shadow: 0 24px 80px -8px rgba(10,10,13,0.30), var(--sh-1); margin: 24px auto; display: flex; flex-direction: column; }
  .x-close { position: absolute; top: 14px; right: 14px; width: 28px; height: 28px; border-radius: 8px; background: transparent; border: 0; display: grid; place-items: center; color: var(--mute); cursor: pointer; transition: background 100ms ease, color 100ms ease; }
  .x-close:hover { background: var(--surface-2); color: var(--ink); }
  .add-hd { margin-bottom: 16px; padding-right: 26px; }
  .add-hd h3 { font-size: 15px; font-weight: 600; margin: 0 0 4px; letter-spacing: -0.015em; }
  .add-hd p { font-size: 13px; color: var(--mute); margin: 0; line-height: 1.5; }
  /* Unified add-event input: one box that takes a file, screenshot, or text. */
  .ev-input { position: relative; background: var(--surface); border: 1.5px dashed var(--rule-strong); border-radius: 12px; overflow: hidden; transition: border-color 120ms ease, background 120ms ease; }
  .ev-input.drag { border-color: var(--accent); background: var(--accent-tint); }
  .ev-input.loading { border-style: solid; border-color: var(--accent); }
  .ev-ta { width: 100%; font: inherit; font-family: var(--sans); font-size: 13.5px; line-height: 1.55; color: var(--ink); background: transparent; border: 0; padding: 13px 15px 8px; outline: none; resize: none; min-height: 78px; box-sizing: border-box; display: block; }
  .ev-ta::placeholder { color: var(--mute-2); }
  .ev-attached { display: flex; align-items: center; gap: 10px; padding: 14px 15px 10px; font-size: 13px; }
  .ev-att-ic { color: var(--accent-text); display: inline-flex; flex-shrink: 0; }
  .ev-att-name { font-weight: 500; color: var(--ink); }
  .ev-att-kind { color: var(--mute); font-size: 11.5px; }
  .ev-att-x { margin-left: auto; background: transparent; border: 0; color: var(--mute); font-size: 18px; line-height: 1; cursor: pointer; padding: 0 4px; }
  .ev-att-x:hover { color: var(--ink); }
  .ev-foot { display: flex; align-items: center; justify-content: space-between; gap: 10px; padding: 8px 14px 10px; border-top: 1px dashed var(--rule); }
  .ev-browse { display: inline-flex; align-items: center; gap: 7px; font-size: 11.5px; color: var(--mute); cursor: pointer; }
  .ev-browse svg { color: var(--mute-2); flex-shrink: 0; }
  .ev-browse u { color: var(--accent-text); text-decoration: none; font-weight: 500; }
  .ev-browse:hover u { text-decoration: underline; }
  .ev-loading { display: inline-flex; align-items: center; gap: 7px; font-size: 12px; font-weight: 500; color: var(--accent-text); }
  .ev-spin { width: 13px; height: 13px; border: 1.8px solid var(--accent-tint-2); border-top-color: var(--accent); border-radius: 50%; animation: ev-spin 0.7s linear infinite; flex-shrink: 0; }
  @keyframes ev-spin { to { transform: rotate(360deg); } }
  .ev-actions { display: flex; align-items: center; justify-content: space-between; gap: 12px; margin-top: 12px; }
  .ev-hint { font-size: 11.5px; color: var(--mute); }
  .ev-hint kbd { font-family: var(--mono, ui-monospace, monospace); font-size: 10.5px; background: var(--surface); border: 1px solid var(--rule); border-bottom-width: 2px; border-radius: 3px; padding: 0 4px; color: var(--ink-2); }
  .ics-preview { margin-top: 18px; padding-top: 16px; border-top: 1px solid var(--rule); }
  .ics-preview h4 { font-size: 11.5px; font-weight: 600; color: var(--mute); text-transform: uppercase; letter-spacing: 0.04em; margin: 0 0 6px; }
  .prev-check { font-size: 12px; color: var(--mute); margin: 0 0 10px; }
  .prev-when strong { font-weight: 600; }
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
  .modal { background: var(--card); border: 1px solid var(--rule); border-radius: 12px; padding: 1.5rem; width: 100%; max-width: 560px; max-height: calc(100vh - 4rem); overflow-y: auto; display: flex; flex-direction: column; gap: .75rem; box-shadow: var(--sh-pop); }
  .modal h2 { font-size: 18px; font-weight: 600; letter-spacing: -0.018em; margin: 0; }
  .modal-hint { font-size: 12px; color: var(--mute); margin: 0 0 .5rem; }
  .fields { display: grid; grid-template-columns: 1fr 1fr; gap: .65rem; }
  .fields .span-2 { grid-column: span 2; }
  .field-group { grid-column: span 2; font-size: 11px; font-weight: 600; letter-spacing: .04em; text-transform: uppercase; color: var(--ink-2); margin-top: .5rem; padding-top: .65rem; border-top: 1px solid var(--rule); }
  .field-group .fg-sub { font-weight: 400; text-transform: none; letter-spacing: 0; color: var(--mute-2); }
  .modal label { display: flex; flex-direction: column; font-size: 12px; color: var(--mute); gap: .35rem; }
  .modal input { font: inherit; color: var(--ink); background: var(--surface); border: 1px solid var(--rule); border-radius: 6px; padding: .45rem .6rem; font-size: 13.5px; outline: none; transition: border-color 100ms ease; }
  .modal input:focus { border-color: var(--accent); }
  .modal .jd-area { font: inherit; font-family: var(--sans); color: var(--ink); background: var(--surface); border: 1px solid var(--rule); border-radius: 6px; padding: .45rem .6rem; font-size: 13.5px; line-height: 1.5; outline: none; resize: vertical; min-height: 72px; transition: border-color 100ms ease; }
  .modal .jd-area:focus { border-color: var(--accent); }
  .modal-actions { display: flex; justify-content: flex-end; gap: .5rem; margin-top: .75rem; }
  .privacy-note { display: flex; align-items: center; gap: 7px; font-size: 11.5px; color: var(--mute); margin: 12px 0 0; line-height: 1.4; }
  .privacy-note svg { color: var(--mute-2); flex-shrink: 0; }

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
