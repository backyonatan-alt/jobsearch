<script>
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { api, isConnectionErr, pollForDossier } from '$lib/api.js';
  import { isPreview, mockApi } from '$lib/preview-mode.js';
  import { logEvent } from '$lib/analytics.js';
  import ConfirmDialog from '$lib/ConfirmDialog.svelte';
  import StatusPill from '$lib/StatusPill.svelte';
  import CompanyLogo from '$lib/CompanyLogo.svelte';
  import {
    toDisplayApp, STATUS_LABEL, SOURCE_SUGGESTIONS,
    fmtLongDate, fmtRelativeDate, daysSince
  } from '$lib/app-helpers.js';

  const call = isPreview() ? mockApi : api;

  // On touch, <datalist> is a dead affordance (tap doesn't open it) — the
  // source chips below the input are the mobile path, so drop the list there.
  const coarsePointer = typeof window !== 'undefined' && !!window.matchMedia?.('(pointer: coarse)').matches;

  // First playbook from the prep-first cold start lands here with ?welcome=1.
  let welcomeDismissed = $state(false);
  const showWelcome = $derived(!welcomeDismissed && page.url.searchParams.get('welcome') === '1');

  // Fire-once guard for dossier_open (one event per app viewed, not per render).
  let lastDossierOpenId = null;

  let app = $state(null);
  let loading = $state(true);
  let notFound = $state(false);

  // Interview prep (dossier) — the selected tab's brief + the shared company brief.
  let dossier = $state(null);
  let dossierLoading = $state(true);
  let companyDossier = $state(null);
  let generating = $state(false);
  let regrounding = $state(false); // identity "Re-research" in flight → show the build state
  let prepStage = $state('');
  let genError = $state('');
  let interviewerInput = $state('');
  let companyUrlInput = $state('');
  // Identity strip ("Not them?") state: bar | confirm | fix | verified. Session-local.
  let identity = $state('bar');
  // Which tab is showing: 'company' (shared brief) or an interview id (that
  // round's interviewer brief).
  let selectedTab = $state('company');
  let prepReady = $state(false); // interviews+debriefs loaded and the default tab chosen
  // Rounds we know have (or lack) a generated brief — filled as tabs load.
  let hasBriefByIv = $state({});

  // Interviews
  let interviews = $state([]);
  let interviewsLoading = $state(false);

  // Prep credits — surfaced so the cap is visible BEFORE a user burns their
  // last generate on a 429. limit 0 = unknown (older /api/me) → don't gate.
  let prepCredits = $state(null); // { used, limit }
  const prepLeft = $derived(
    prepCredits && prepCredits.limit > 0 ? Math.max(0, prepCredits.limit - prepCredits.used) : null
  );
  const atPrepLimit = $derived(prepLeft === 0);
  async function loadCredits() {
    try {
      const me = await call('/api/me');
      prepCredits = { used: me.prep_credits_used || 0, limit: me.prep_credits_limit || 0 };
    } catch { prepCredits = null; }
  }
  const capMailto = `mailto:back.yonatan@gmail.com?subject=${encodeURIComponent('Pursuit — more prep credits please')}`;

  // Debriefs — keyed by interview_id. The post-round "how did it go / was the
  // prep right?" that feeds the next round's playbook.
  let debriefsByIv = $state({});
  let debriefDraft = $state({ feel: '', prep_accuracy: '', topics: '', notes: '' });
  let debriefSaving = $state(false);
  let debriefEditing = $state(false);

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

  // Edit modal
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
    loadCredits();
    initPrep();
    if (id && id !== lastDossierOpenId) {
      lastDossierOpenId = id;
      logEvent('dossier_open', { app_id: Number(id) });
    }
  });

  // Load interviews first so we can default the prep to the next upcoming round.
  // With no upcoming round, prefer a past round still waiting on its debrief
  // (that's when the prompt matters) before falling back to the Company tab.
  // ?debrief=<interview_id> (Today's proactive prompt) deep-links straight into
  // that round with the debrief form open.
  async function initPrep() {
    prepReady = false;
    await loadInterviews();
    await loadDebriefs();
    loadCompanyBrief();
    const want = Number(page.url.searchParams.get('debrief'));
    const deepLinked = want && (interviews || []).some(iv => iv.id === want);
    selectedTab = deepLinked ? want : (nextRoundId ?? pendingDebriefRoundId ?? 'company');
    prepReady = true;
    if (deepLinked && !debriefsByIv[want]) startDebrief();
    await loadDossier(selectedTab);
    if (deepLinked) {
      try { document.getElementById('interview-prep')?.scrollIntoView({ behavior: 'smooth', block: 'start' }); } catch {}
    }
  }

  async function loadCompanyBrief() {
    try {
      companyDossier = (await call(`/api/applications/${id}/dossier?scope=company`)) || null;
    } catch { companyDossier = null; }
  }

  async function loadDebriefs() {
    try {
      const list = await call(`/api/applications/${id}/debriefs`);
      const m = {};
      for (const d of (list || [])) m[d.interview_id] = d;
      debriefsByIv = m;
    } catch { debriefsByIv = {}; }
  }

  const onCompany = $derived(selectedTab === 'company');
  const selectedRound = $derived(onCompany ? null : (interviews || []).find(iv => iv.id === selectedTab) || null);

  // The round's own debrief (if any), and whether the round is in the past (so we
  // prompt for a debrief). Undated rounds are treated as debrief-able.
  const roundDebrief = $derived(onCompany ? null : (debriefsByIv[selectedTab] ?? null));
  const roundIsPast = $derived.by(() => {
    if (onCompany || !selectedRound) return false;
    if (selectedRound.scheduled === false) return true; // one-tap round → always debrief-able
    if (!selectedRound.starts_at) return true;
    return new Date(selectedRound.starts_at).getTime() < Date.now();
  });
  // Most recent debrief-able round (past, or undated one-tap) with no debrief yet.
  const pendingDebriefRoundId = $derived.by(() => {
    const now = Date.now();
    const due = (interviews || []).filter(iv => {
      if (debriefsByIv[iv.id]) return false;
      if (iv.scheduled === false || !iv.starts_at) return true;
      return new Date(iv.starts_at).getTime() < now;
    });
    if (!due.length) return null;
    due.sort((a, b) => new Date(b.starts_at || 0) - new Date(a.starts_at || 0));
    return due[0].id;
  });

  // Did an EARLIER round's debrief feed this round's playbook?
  const informedByDebrief = $derived.by(() => {
    if (onCompany || !selectedRound) return false;
    const mine = selectedRound.starts_at ? new Date(selectedRound.starts_at).getTime() : Infinity;
    return (interviews || []).some(iv =>
      iv.id !== selectedTab && iv.starts_at &&
      new Date(iv.starts_at).getTime() < mine && debriefsByIv[iv.id]);
  });

  function feelLabel(f) { return { strong: 'strong', mixed: 'mixed', rough: 'rough' }[f] || f; }
  function accLabel(a) { return { spot_on: 'spot on', partly: 'partly right', off: 'off' }[a] || a; }

  // One-tap "Add round" (no date). Creates an unscheduled round and opens it.
  const ROUND_PRESETS = ['Recruiter screen', 'Phone screen', 'Hiring manager', 'Technical screen', 'System design', 'Onsite', 'Team match', 'Behavioral'];
  let showAddRound = $state(false);
  let addRoundText = $state('');
  async function createRound(summary) {
    const label = (summary ?? '').trim();
    if (!label) return null;
    // Reuse an existing round with the same label rather than duplicating.
    const existing = (interviews || []).find(x => (x.summary || '').trim().toLowerCase() === label.toLowerCase());
    let iv = existing;
    if (!iv) {
      try {
        iv = await call(`/api/applications/${id}/interviews`, { method: 'POST', body: JSON.stringify({ summary: label, source: 'manual' }) });
      } catch (e) { if (e.message !== 'unauthorized') console.error(e); return null; }
      await loadInterviews();
      try { window.dispatchEvent(new CustomEvent('pursuit:refresh')); } catch {}
    }
    showAddRound = false; addRoundText = '';
    selectedTab = iv.id;
    debriefEditing = false;
    await loadDossier(iv.id);
    return iv;
  }

  // Stage-done hook: marking a pipeline stage done offers to debrief that round.
  let stageDebriefPrompt = $state(null); // { name }
  async function debriefFromStage() {
    const name = stageDebriefPrompt?.name;
    stageDebriefPrompt = null;
    if (!name) return;
    const iv = await createRound(name);
    if (iv) { startDebrief(); try { document.getElementById('interview-prep')?.scrollIntoView({ behavior: 'smooth', block: 'start' }); } catch {} }
  }

  function startDebrief() {
    const d = roundDebrief;
    debriefDraft = d
      ? { feel: d.feel, prep_accuracy: d.prep_accuracy, topics: d.topics ?? '', notes: d.notes ?? '' }
      : { feel: '', prep_accuracy: '', topics: '', notes: '' };
    debriefEditing = true;
  }

  // Log the prompt being SEEN (rendered), not clicked — once per round per visit.
  const debriefViewLogged = new Set();
  function logDebriefView(key, surface) {
    if (debriefViewLogged.has(key)) return;
    debriefViewLogged.add(key);
    logEvent('debrief_prompt_view', { app_id: Number(id), surface });
  }
  $effect(() => {
    if (prepReady && !onCompany && roundIsPast && !roundDebrief && !debriefEditing) logDebriefView(`round:${selectedTab}`, 'round_tab');
  });
  $effect(() => {
    if (prepReady && pendingDebriefRoundId && pendingDebriefRoundId !== selectedTab) logDebriefView(`banner:${pendingDebriefRoundId}`, 'banner');
  });
  $effect(() => {
    if (stageDebriefPrompt) logDebriefView(`stage:${stageDebriefPrompt.name}`, 'stage_done');
  });

  // Jump from the banner to the pending round with the debrief form open.
  async function debriefPendingRound() {
    const iv = pendingDebriefRoundId;
    if (!iv) return;
    await selectTab(iv);
    startDebrief();
    try { document.getElementById('interview-prep')?.scrollIntoView({ behavior: 'smooth', block: 'start' }); } catch {}
  }
  async function saveDebrief() {
    if (debriefSaving || !debriefDraft.feel || !debriefDraft.prep_accuracy) return;
    debriefSaving = true;
    try {
      const saved = await call(`/api/applications/${id}/interviews/${selectedTab}/debrief`, {
        method: 'POST', body: JSON.stringify(debriefDraft)
      });
      debriefsByIv = { ...debriefsByIv, [selectedTab]: saved };
      debriefEditing = false;
    } catch (e) {
      if (e.message !== 'unauthorized') console.error(e);
    } finally { debriefSaving = false; }
  }

  // The soonest upcoming interview — the round we prep for by default.
  const nextRoundId = $derived.by(() => {
    const now = Date.now();
    const future = (interviews || [])
      .filter(iv => iv?.scheduled !== false && iv?.starts_at && new Date(iv.starts_at).getTime() >= now)
      .sort((a, b) => new Date(a.starts_at) - new Date(b.starts_at));
    return future[0]?.id ?? null;
  });

  function roundLabel(iv) {
    // Only dated (scheduled) rounds show a relative date; one-tap rounds show the label alone.
    const d = (iv?.scheduled !== false && iv?.starts_at) ? fmtRelativeDate(iv.starts_at) : '';
    const s = (iv?.summary || '').trim();
    return s ? (d ? `${d} · ${s}` : s) : (d || 'Interview');
  }
  // Tab caption: name first, near-future suffix ("Panel · tomorrow").
  function tabLabel(iv) {
    const s = (iv?.summary || '').trim() || 'Interview';
    if (iv?.scheduled !== false && iv?.starts_at && new Date(iv.starts_at).getTime() >= Date.now()) {
      const rel = evRelative(iv);
      if (['today', 'tomorrow'].includes(rel) || rel.startsWith('in ')) return `${s} · ${rel}`;
    }
    return s;
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
    debriefEditing = false; // don't carry an open debrief editor across rounds
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
    const wasDone = !!pipeline[i]?.done;
    const next = pipeline.map((s, idx) => idx === i ? { ...s, done: !s.done } : s);
    pipeline = next;
    savePipeline(next);
    // Just completed a stage → offer to debrief that round (unless already debriefed).
    if (!wasDone && next[i].done && next[i].name) {
      const already = (interviews || []).find(x => (x.summary || '').trim().toLowerCase() === next[i].name.trim().toLowerCase());
      if (!already || !debriefsByIv[already.id]) stageDebriefPrompt = { name: next[i].name };
    }
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
      if (tab === 'company') companyDossier = dossier;
      else hasBriefByIv = { ...hasBriefByIv, [tab]: !!dossier };
    } catch (e) {
      // 404 / empty → this tab hasn't been generated yet
      dossier = null;
      if (tab !== 'company') hasBriefByIv = { ...hasBriefByIv, [tab]: false };
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
    const prevGeneratedAt = dossier?.generated_at ?? null;
    try {
      const companyUrl = companyUrlInput.trim() || undefined;
      const body = onCompany
        ? { company_url: companyUrl }
        : { interview_id: selectedTab, interviewer_name: interviewerInput.trim() || undefined, company_url: companyUrl };
      const d = await call(`/api/applications/${id}/dossier/refresh`, {
        method: 'POST',
        body: JSON.stringify(body)
      });
      applyGenerated(d);
    } catch (e) {
      // A dropped connection doesn't kill the build — the server finishes and
      // saves. Poll for the fresh brief before showing an error.
      if (isConnectionErr(e)) {
        timers.forEach(clearTimeout);
        prepStage = 'Connection hiccup — still building on our side, watching for it…';
        const q = onCompany ? '?scope=company' : `?interview_id=${selectedTab}`;
        const d = await pollForDossier(`/api/applications/${id}/dossier${q}`, prevGeneratedAt);
        if (d) {
          applyGenerated(d);
          generating = false;
          prepStage = '';
          return;
        }
      }
      genError = friendlyGenErr(e.message);
    } finally {
      timers.forEach(clearTimeout);
      generating = false;
      prepStage = '';
    }
  }
  function applyGenerated(d) {
    dossier = d;
    interviewerInput = d.interviewer_name ?? interviewerInput;
    companyUrlInput = '';
    if (onCompany) companyDossier = d;
    else {
      hasBriefByIv = { ...hasBriefByIv, [selectedTab]: true };
      // Round generation may have built the shared company brief alongside it.
      loadCompanyBrief();
    }
    loadCredits();
  }

  function friendlyGenErr(msg) {
    const m = String(msg || '');
    if (/interview-prep limit/i.test(m)) {
      loadCredits(); // refresh so the cap note replaces the button
      return "You've used all your prep credits for the beta — email me and I'll top you up.";
    }
    if (/failed to fetch|load failed|networkerror/i.test(m))
      return 'The connection dropped mid-build. If it finished on our side it will appear when you reload — otherwise try again.';
    if (m.includes('rate_limit_error') || m.includes('429'))
      return 'AI usage limit hit — wait a minute and try again.';
    if (m.includes('http 504') || /\btimeout\b/i.test(m))
      return 'Web search timed out — try again.';
    if (m.includes('http 5') || m.includes('not configured'))
      return 'Something went wrong — try again in a moment.';
    return m || 'Could not build the playbook.';
  }

  // ── Identity strip ("Not them?") — wrong-company disambiguation flow ──
  function identityOpen() {
    identity = 'confirm';
    logEvent('dossier_identity_open', { app_id: Number(id) });
  }
  function identityConfirm() {
    identity = 'verified';
    logEvent('identity_confirmed', { app_id: Number(id) });
  }
  function identityReject() {
    identity = 'fix';
    companyUrlInput = '';
    logEvent('identity_rejected', { app_id: Number(id) });
  }
  async function identityReresearch() {
    if (!companyUrlInput.trim() || generating) return;
    logEvent('identity_reresearch_submitted', { app_id: Number(id) });
    identity = 'bar';
    regrounding = true;
    // Re-grounding rebuilds the COMPANY brief — show the build on the company tab.
    if (selectedTab !== 'company') await selectTab('company');
    try { await generateDossier(); } finally { regrounding = false; }
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

  // ── Edit / delete (preserved) ─────────────────────────────────
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

  // Selected-tab brief content (round or company, depending on the tab).
  const dosContent = $derived(dossier?.content ?? null);
  const dosInterviewer = $derived(dosContent?.interviewer ?? null);
  const dosIvInitials = $derived(
    initialsOf(dosInterviewer?.name ?? dossier?.interviewer_name ?? '')
  );
  const dosIvName = $derived(
    dosInterviewer?.name ?? dossier?.interviewer_name ?? ''
  );
  const dosGeneratedAgo = $derived(dossier?.generatedAgo ?? '');
  const headerAgo = $derived(dossier?.generatedAgo ?? companyDossier?.generatedAgo ?? '');

  // Shared company brief content — used for the identity strip + Company tab.
  const companyContent = $derived(companyDossier?.content ?? null);
  const dosIdentity = $derived(companyContent?.identity ?? null);
  const dosCompany = $derived(companyContent?.company ?? null);
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
  const companySources = $derived(Array.isArray(companyContent?.sources) ? companyContent.sources.filter(s => s?.href) : []);
  const tips = $derived((dosCompany?.watch_fors ?? []).slice(0, 5));

  function dosSigDomain(src) {
    if (!src) return '';
    try {
      return new URL(src.startsWith('http') ? src : `https://${src}`).hostname.replace(/^www\./, '');
    } catch { return src; }
  }
  const identityDomainHref = $derived(
    dosIdentity?.domain ? `https://${String(dosIdentity.domain).replace(/^https?:\/\//, '')}` : ''
  );

  const appliedLong = $derived(app ? fmtLongDate(app.raw.applied_at) : '');

  // Next event = the soonest scheduled interview dated now-or-later.
  const upcoming = $derived.by(() => {
    const now = Date.now();
    const future = (interviews || [])
      .filter(iv => iv?.scheduled !== false && iv?.starts_at && new Date(iv.starts_at).getTime() >= now)
      .sort((a, b) => new Date(a.starts_at) - new Date(b.starts_at));
    return future[0] || null;
  });

  const closedOut = $derived(app && ['rejected', 'withdrawn', 'closed'].includes(app.status));

  // Tomorrow-banner: the next round when it's within a week.
  const soonRound = $derived.by(() => {
    if (!upcoming || closedOut) return null;
    const days = Math.round((startOfDay(new Date(upcoming.starts_at)) - startOfDay(new Date())) / 86400000);
    return days <= 7 ? upcoming : null;
  });
  function startOfDay(x) { return new Date(x.getFullYear(), x.getMonth(), x.getDate()); }
  const soonCal = $derived.by(() => {
    if (!soonRound) return null;
    const d = new Date(soonRound.starts_at);
    return {
      mon: d.toLocaleDateString('en-US', { month: 'short' }).toUpperCase(),
      day: d.getDate()
    };
  });
  const soonLine = $derived.by(() => {
    if (!soonRound) return '';
    const d = new Date(soonRound.starts_at);
    const rel = evRelative(soonRound);
    const time = soonRound.all_day ? '' : ` at ${d.toLocaleTimeString(undefined, { hour: 'numeric', minute: '2-digit' })}`;
    const mins = evDuration(soonRound);
    return `${(soonRound.summary || 'Interview').trim()} ${rel}${time}${mins ? ` — ${mins} minutes.` : '.'}`;
  });
  function reviewPrep() {
    if (!soonRound) return;
    if (hasBriefByIv[soonRound.id]) { goto(`/app/${id}/brief/${soonRound.id}`); return; }
    selectTab(soonRound.id);
    try { document.getElementById('interview-prep')?.scrollIntoView({ behavior: 'smooth', block: 'start' }); } catch {}
  }

  function evWhen(iv) {
    if (!iv?.starts_at) return '';
    const d = new Date(iv.starts_at);
    if (iv.all_day) return d.toLocaleDateString(undefined, { weekday: 'short', month: 'short', day: 'numeric' }) + ' · all day';
    const now = new Date();
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
  function isPastEvent(ev) { return new Date(ev.starts_at) < new Date(); }

  const monoDate = (iso) => iso ? new Date(iso).toLocaleDateString('en-US', { day: 'numeric', month: 'short' }) : '';

  // The round tile's interviewer name (from the brief, else the round's attendee).
  const roundTileName = $derived.by(() => {
    if (onCompany) return '';
    return dosIvName || (selectedRound ? attendeeName(selectedRound) : '') || '';
  });

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
        note: isPastEvent(iv) ? 'Past event' : 'Scheduled',
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

  // "The role, in short" — honest one-liner from the fields we actually have.
  const roleShort = $derived.by(() => {
    if (!app) return '';
    const bits = [];
    if (app.raw.location) bits.push(app.raw.location);
    if (app.raw.salary_note) bits.push(app.raw.salary_note);
    if (app.cv && app.cv !== '—') bits.push(`CV ${app.cv}`);
    return bits.join(' · ');
  });
  const jdSnippet = $derived.by(() => {
    const t = (app?.raw?.jd_text || '').trim().replace(/\s+/g, ' ');
    if (!t) return '';
    return t.length > 220 ? t.slice(0, 220).replace(/\s+\S*$/, '') + '…' : t;
  });
</script>

<svelte:head>
  <title>{app?.co ? `${app.co} — Pursuit` : 'Pursuit'}</title>
</svelte:head>

<svelte:window onkeydown={onWindowKeydown} />

<div class="wrap">
  {#if loading}
    <p class="loading">Loading…</p>
  {:else if notFound || !app}
    <div class="empty-tab">
      <h3>Application not found</h3>
      <p>It may have been deleted, or you might not have access. <a href="/app">Back to Home →</a></p>
    </div>
  {:else}

    <div class="crumb"><a href="/app/applications">Applications</a> <span class="sep">/</span> <span class="here">{app.co}</span></div>

    {#if showWelcome}
      <div class="welcome-banner">
        <span class="wb-spark">✦</span>
        <span class="wb-tx">Here's your first playbook. Add who's interviewing you for round-by-round prep, or <a href="/app">track another application</a>.</span>
        <button class="wb-x" onclick={() => (welcomeDismissed = true)} aria-label="Dismiss">✕</button>
      </div>
    {/if}

    <!-- HEADER -->
    <div class="hd">
      <div class="hd-logo"><CompanyLogo app={app} size={52} radius={13} /></div>
      <div class="hd-main">
        <div class="hd-row">
          <h1>{app.co}</h1>
          <StatusPill id={Number(id)} status={app.status} surface="detail_pill" align="left" onchanged={() => loadApp()} />
        </div>
        <div class="hd-role">{app.role}</div>
        <div class="hd-sub">
          {#if app.raw.applied_at}<span>Applied {appliedLong}</span>{/if}
          {#if app.source && app.source !== '—'}<span class="dot">·</span><span>{app.source}</span>{/if}
          {#if app.raw.location}<span class="dot">·</span><span>{app.raw.location}</span>{/if}
          {#if app.raw.jd_url}<span class="dot">·</span><a href={app.raw.jd_url} target="_blank" rel="noopener">Open job post ↗</a>{/if}
        </div>
      </div>
      <div class="hd-actions">
        <button class="linkbtn" onclick={openEdit}>Edit</button>
        <button class="linkbtn danger" onclick={deleteApp}>Delete</button>
      </div>
    </div>

    <!-- Soon banner -->
    {#if soonRound && soonCal}
      <div class="soon">
        <div class="soon-cal"><div class="sc-mon">{soonCal.mon}</div><div class="sc-day">{soonCal.day}</div></div>
        <div class="soon-tx">
          <div class="soon-t">{soonLine}</div>
          <div class="soon-s">{hasBriefByIv[soonRound.id] ? 'Your brief is ready below. 20 focused minutes is enough.' : 'Build your brief below — 20 focused minutes is enough.'}</div>
        </div>
        <button class="soon-cta" onclick={reviewPrep}>Review prep →</button>
      </div>
    {/if}

    <div class="grid">
      <!-- LEFT — playbook -->
      <div class="col-main">
        <div id="interview-prep" class="pb-card">

          <!-- Identity strip + "Not them?" wrong-company flow (from the company brief) -->
          {#if dosIdentity}
            {#if identity === 'bar'}
              <div class="idbar">
                <CompanyLogo app={app} size={24} radius={6} />
                <span class="id-tx">Researched for <strong>{dosIdentity.name || app.co}</strong>{#if dosIdentity.domain}&nbsp;· <a href={identityDomainHref} target="_blank" rel="noreferrer">{dosIdentity.domain}</a>{/if}{#if dosIdentity.summary}&nbsp;<span class="id-sum">— {dosIdentity.summary}</span>{/if}</span>
                <button class="linkbtn id-not" onclick={identityOpen}>Not them?</button>
              </div>
            {:else if identity === 'confirm'}
              <div class="idbox">
                <div class="id-q">Is this the company you're interviewing with?</div>
                <div class="id-card">
                  <CompanyLogo app={app} size={38} radius={10} />
                  <div class="id-card-tx"><strong>{dosIdentity.name || app.co}</strong>{#if dosIdentity.domain}&nbsp;· <a href={identityDomainHref} target="_blank" rel="noreferrer">{dosIdentity.domain}</a>{/if}<br><span class="id-sum">{dosIdentity.summary || ''}</span></div>
                </div>
                <div class="id-actions">
                  <button class="btn-blue" onclick={identityConfirm}>Yes — that's them</button>
                  <button class="btn-ghost-warm" onclick={identityReject}>No — wrong company</button>
                  <span class="id-hint">Same-name companies are the #1 cause of a wasted brief.</span>
                </div>
              </div>
            {:else if identity === 'fix'}
              <div class="idbox">
                <div class="id-q">Point us at the right one.</div>
                <div class="id-fix-sub">Paste their website or LinkedIn page — we'll re-research and rebuild the company brief.</div>
                <div class="id-fix-row">
                  <input class="id-input" placeholder="https://" bind:value={companyUrlInput} disabled={generating}
                    onkeydown={(e) => e.key === 'Enter' && identityReresearch()} />
                  <button class="btn-blue" onclick={identityReresearch} disabled={generating || !companyUrlInput.trim()}>Re-research →</button>
                  <button class="linkbtn id-cancel" onclick={() => (identity = 'bar')}>Cancel</button>
                </div>
              </div>
            {:else}
              <div class="idok">
                <span>✓ Verified — {dosIdentity.name || app.co}{#if dosIdentity.domain}&nbsp;· {dosIdentity.domain}{/if}</span>
                <button class="linkbtn id-change" onclick={identityOpen}>change</button>
              </div>
            {/if}
          {/if}

          <div class="pb-hd">
            <div class="pb-title-row">
              <div class="pb-title">✦ Interview playbook</div>
              {#if headerAgo}<div class="pb-ago">refreshed {headerAgo}</div>{/if}
              {#if informedByDebrief && dossier}
                <span class="informed-chip" title="This round's prep was tailored using your debrief of an earlier round">
                  <span class="ic-dot"></span>Informed by your last round
                </span>
              {/if}
            </div>
            <div class="tabs" role="tablist" aria-label="Interview round">
              <button type="button" role="tab" aria-selected={onCompany}
                class="tab company" class:active={onCompany}
                onclick={() => selectTab('company')}>▦ Company</button>
              {#each interviews as iv (iv.id)}
                <button type="button" role="tab" aria-selected={selectedTab === iv.id}
                  class="tab round" class:active={selectedTab === iv.id} class:done={!!debriefsByIv[iv.id]}
                  onclick={() => selectTab(iv.id)}>{debriefsByIv[iv.id] && selectedTab !== iv.id ? '✓ ' : ''}{tabLabel(iv)}</button>
              {/each}
              <button type="button" class="tab add" onclick={() => (showAddRound = !showAddRound)} title="Add a round">+ Add round</button>
            </div>
          </div>

          <div class="pb-body">
            {#if showAddRound}
              <div class="add-round">
                <span class="ar-lbl">Add a round you did or have coming up:</span>
                <div class="ar-chips">
                  {#each ROUND_PRESETS as p (p)}
                    <button type="button" class="ar-chip" onclick={() => createRound(p)}>{p}</button>
                  {/each}
                </div>
                <div class="ar-custom">
                  <input class="ar-input" placeholder="Or type a round name…" bind:value={addRoundText}
                    onkeydown={(e) => e.key === 'Enter' && addRoundText.trim() && createRound(addRoundText)} />
                  <button type="button" class="btn" onclick={() => createRound(addRoundText)} disabled={!addRoundText.trim()}>Add</button>
                </div>
              </div>
            {/if}

            <!-- A past round on ANOTHER tab still needs its debrief → slim jump-in banner. -->
            {#if prepReady && pendingDebriefRoundId && pendingDebriefRoundId !== selectedTab}
              {@const pendingIv = (interviews || []).find(x => x.id === pendingDebriefRoundId)}
              <button type="button" class="db-banner" onclick={debriefPendingRound}>
                <span class="db-spark">✦</span>
                <span class="db-banner-tx">How did the <b>{(pendingIv?.summary || 'last').trim() || 'last'}</b> round go? A 20-second debrief sharpens this round's prep.</span>
                <span class="db-prompt-cta">Debrief →</span>
              </button>
            {/if}

            <!-- Debrief: capture how a past round went → feeds next round. -->
            {#if roundIsPast}
              <section class="debrief">
                {#if roundDebrief && !debriefEditing}
                  <div class="db-summary">
                    <div class="db-sum-main">
                      <span class="db-badge">Debriefed</span>
                      <span class="db-sum-line">Went <b>{feelLabel(roundDebrief.feel)}</b> · prep was <b>{accLabel(roundDebrief.prep_accuracy)}</b></span>
                    </div>
                    <button class="db-edit" type="button" onclick={startDebrief}>Edit</button>
                  </div>
                  {#if roundDebrief.topics}<p class="db-topics">What came up: {roundDebrief.topics}</p>{/if}
                {:else if debriefEditing}
                  <div class="db-form">
                    <div class="db-hd"><h3>How did this round go?</h3><span class="db-hint">20 seconds — it sharpens your next round's prep</span></div>
                    <div class="db-q">
                      <span class="db-q-lbl">How it felt</span>
                      <div class="db-opts">
                        {#each [['strong','Strong'],['mixed','Mixed'],['rough','Rough']] as opt (opt[0])}
                          <button type="button" class="db-opt" class:sel={debriefDraft.feel === opt[0]} onclick={() => (debriefDraft.feel = opt[0])}>{opt[1]}</button>
                        {/each}
                      </div>
                    </div>
                    <div class="db-q">
                      <span class="db-q-lbl">Was our prep right?</span>
                      <div class="db-opts">
                        {#each [['spot_on','Spot on'],['partly','Partly'],['off','Off']] as opt (opt[0])}
                          <button type="button" class="db-opt" class:sel={debriefDraft.prep_accuracy === opt[0]} onclick={() => (debriefDraft.prep_accuracy = opt[0])}>{opt[1]}</button>
                        {/each}
                      </div>
                    </div>
                    <input class="db-input" placeholder="What actually came up? (optional)" bind:value={debriefDraft.topics} />
                    <div class="db-actions">
                      <button class="btn" type="button" onclick={() => (debriefEditing = false)}>Cancel</button>
                      <button class="btn btn-primary" type="button" onclick={saveDebrief} disabled={debriefSaving || !debriefDraft.feel || !debriefDraft.prep_accuracy}>{debriefSaving ? 'Saving…' : 'Save debrief'}</button>
                    </div>
                  </div>
                {:else}
                  <button type="button" class="db-prompt" onclick={startDebrief}>
                    <span class="db-spark">✦</span>
                    <span class="db-prompt-tx"><b>How did this round go?</b><small>A 20-second debrief sharpens your next round's prep.</small></span>
                    <span class="db-prompt-cta">Debrief →</span>
                  </button>
                {/if}
              </section>
            {/if}

            {#if generating}
              <!-- Build state (generation or identity re-research) -->
              <div class="genbox">
                <h3>Researching {app.co}{!onCompany && interviewerInput ? ` & ${interviewerInput}` : ''}…</h3>
                <p class="gen-sub" aria-live="polite">{prepStage || PREP_STAGES[0]}</p>
                <div class="big-spinner"></div>
                <p class="gen-eta">This usually takes 1–2 minutes — you can keep working, it'll be here when it's done.</p>
              </div>

            {:else if dossierLoading}
              <p class="loading">Loading…</p>

            {:else if !dossier}
              <!-- Generate / empty state -->
              {#if atPrepLimit}
                <div class="genbox">
                  <h3>You've used all {prepCredits.limit} prep credits for the beta.</h3>
                  <p class="gen-sub">Everything already generated stays yours. <a href={capMailto} onclick={() => logEvent('feedback_click', { surface: 'credit_cap' })}>Write to us</a> — we won't leave you hanging the night before.</p>
                  <p class="gen-eta">Tracker features — statuses, follow-ups, archive, the home page — are never limited by credits.</p>
                </div>
              {:else if onCompany}
                <div class="genbox">
                  <h3>Generate the company brief</h3>
                  <p class="gen-sub">A shared brief on {app.co} — what they do, where they're headed, the typical loop, and what this team grades for. Researched once and used across every round.</p>
                  <button class="btn-generate" onclick={generateDossier} disabled={generating}>Generate company brief</button>
                  {#if prepLeft !== null && prepLeft <= 2}<p class="credit-left">{prepLeft} prep credit{prepLeft === 1 ? '' : 's'} left in your beta allowance.</p>{/if}
                  {#if genError}<p class="gen-err">{genError}</p>{/if}
                </div>
              {:else}
                <div class="genbox">
                  <h3>Brief for {selectedRound ? (selectedRound.summary || 'this round') : 'this round'}</h3>
                  <p class="gen-sub">We'll research the person interviewing you in this round — their background, how they tend to interview, what lands, and smart questions to ask. The shared {app.co} company brief is generated alongside it if you don't have one yet, so you only wait once.</p>
                  <div class="gen-row">
                    <input class="gen-input" type="text"
                      placeholder="Interviewer name (optional) — e.g. Sarah Chen"
                      bind:value={interviewerInput} disabled={generating}
                      onkeydown={(e) => e.key === 'Enter' && generateDossier()} />
                    <button class="btn-generate" onclick={generateDossier} disabled={generating}>Build the brief</button>
                  </div>
                  {#if prepLeft !== null && prepLeft <= 2}<p class="credit-left">{prepLeft} prep credit{prepLeft === 1 ? '' : 's'} left in your beta allowance.</p>{/if}
                  {#if genError}<p class="gen-err">{genError}</p>{/if}
                </div>
              {/if}

            {:else if onCompany}
              <!-- COMPANY BRIEF — rendered inline -->
              <div class="cb-hd">Company brief <span class="cb-shared">· shared across every round</span></div>
              {#if companyBlurb}<p class="cb-blurb">{companyBlurb}</p>{/if}
              {#if companyAbout}<p class="cb-about">{companyAbout}</p>{/if}
              {#if companyFacts.length}
                <div class="cb-facts">
                  {#each companyFacts as f (f.lbl)}
                    <div class="cb-cell"><div class="f-lbl">{f.lbl}</div><div class="f-val">{f.val}</div></div>
                  {/each}
                </div>
              {/if}
              {#if companyProcess.length}
                <div class="cb-proc">
                  <div class="cb-lbl">The loop, as reported</div>
                  <div class="cb-chips">
                    {#each companyProcess as step, i (i)}
                      <span class="chip">{step}</span>{#if i < companyProcess.length - 1}<span class="chip-arrow">→</span>{/if}
                    {/each}
                  </div>
                </div>
              {/if}
              {#if tips.length}
                <div class="grades">
                  <div class="grades-hd">✦ What this team grades for</div>
                  <div class="grades-list">
                    {#each tips as t (t)}<div>· {t}</div>{/each}
                  </div>
                </div>
              {/if}
              {#if companySources.length}
                <div class="srcs">
                  <span class="srcs-lbl">Sources</span>
                  {#each companySources as s (s.href)}
                    <a class="src-chip" href={s.href} target="_blank" rel="noreferrer">{s.label || dosSigDomain(s.href)}</a>
                  {/each}
                </div>
              {/if}
              {#if genError}<p class="gen-err">{genError}</p>{/if}

            {:else}
              <!-- ROUND — the tab is a door into the reading page -->
              <div class="rd-hd">{(selectedRound?.summary || 'Interview').trim()}{#if selectedRound?.starts_at && selectedRound?.scheduled !== false}&nbsp;<span class="rd-when">· {evWhen(selectedRound)}</span>{/if}</div>
              {#if roundTileName}
                <div class="pcards">
                  <a class="pcard" href={`/app/${id}/brief/${selectedTab}`}>
                    <div class="pc-av">{initialsOf(roundTileName) || '?'}</div>
                    <div class="pc-name">{roundTileName}</div>
                    {#if dosInterviewer?.role}<div class="pc-role">{dosInterviewer.role}</div>{/if}
                  </a>
                </div>
              {/if}
              <a class="doorbar" href={`/app/${id}/brief/${selectedTab}`}>
                <div class="db-main">
                  <div class="db-t">Open the {(selectedRound?.summary || 'round').trim().toLowerCase()} brief →</div>
                  <div class="db-s">What they grade for, likely questions, your angle, sources · 5 min read</div>
                </div>
                <span class="db-arrow">→</span>
              </a>
              {#if genError}<p class="gen-err">{genError}</p>{/if}
            {/if}
          </div>
        </div>

        <!-- credits footnote -->
        <div class="pb-foot">
          {#if dossier && !generating}
            <button class="linkbtn foot-refresh" onclick={generateDossier} disabled={generating}>
              {onCompany ? 'Refresh company brief' : "Refresh this round's brief"}
            </button>
            <span class="foot-sep">·</span>
          {/if}
          <span>1 credit per round brief{#if prepLeft !== null}&nbsp;· <strong>{prepLeft} left</strong>{/if}</span>
          {#if !onCompany && dossier}
            <span class="foot-sep">·</span>
            <a href="/privacy" target="_blank" rel="noreferrer">how we research people</a>
          {/if}
          {#if selectedRound && selectedRound.scheduled === false}
            <button class="linkbtn foot-remove" onclick={() => deleteInterview(selectedRound)}>Remove round</button>
          {/if}
        </div>
      </div>

      <!-- RIGHT column -->
      <aside class="col-side">
        <div class="side-sec first">
          <div class="side-hd-row">
            <span class="side-lbl">The role, in short</span>
            {#if app.raw.jd_url}<a class="side-act" href={app.raw.jd_url} target="_blank" rel="noopener">Full JD →</a>{/if}
          </div>
          <div class="side-prose">
            {#if jdSnippet}{jdSnippet}{:else}{app.role}{#if roleShort}&nbsp;· {roleShort}{/if}{/if}
          </div>
          {#if jdSnippet && roleShort}<div class="side-note">{roleShort}</div>{/if}
          {#if app.raw.jd_text}
            <details class="jd-saved">
              <summary>Saved job description</summary>
              <p class="jd-body">{app.raw.jd_text}</p>
            </details>
          {/if}
          <div class="side-note">Last activity {fmtRelativeDate(app.raw.updated_at ?? app.raw.applied_at)} · <button class="linkbtn side-edit" onclick={openEdit}>edit details</button></div>
        </div>

        <div class="side-sec">
          <div class="side-hd-row">
            <span class="side-lbl">Process</span>
            {#if pipeline.length && !pipelineEditing}
              <button class="side-act linkbtn" onclick={startEditPipeline}>Edit</button>
            {/if}
          </div>
          {#if pipelineEditing}
            <div class="pipe-edit">
              {#each pipelineDraft as st, i (i)}
                <div class="pe-row" class:pe-dragging={pipeDragIdx === i} ondragover={(e) => onStageDragOver(e, i)} role="listitem">
                  <span class="pe-grip" draggable="true" ondragstart={() => onStageDragStart(i)} ondragend={onStageDragEnd} title="Drag to reorder" aria-label="Drag to reorder" role="button" tabindex="-1">⠿</span>
                  <input class="pe-input" bind:value={st.name} placeholder="Stage name" />
                  <button class="pe-btn" onclick={() => moveDraft(i, -1)} disabled={i === 0} aria-label="Move up">↑</button>
                  <button class="pe-btn" onclick={() => moveDraft(i, 1)} disabled={i === pipelineDraft.length - 1} aria-label="Move down">↓</button>
                  <button class="pe-btn pe-x" onclick={() => removeDraftStage(i)} aria-label="Remove stage">×</button>
                </div>
              {/each}
              <button class="add-line" onclick={addDraftStage}>+ Add stage</button>
              <div class="pe-actions">
                <button class="btn" onclick={() => (pipelineEditing = false)}>Cancel</button>
                <button class="btn btn-primary" onclick={saveEditPipeline} disabled={pipelineSaving}>{pipelineSaving ? 'Saving…' : 'Save'}</button>
              </div>
            </div>
          {:else if pipeline.length}
            <div class="proc-list">
              {#each pipeline as st, i (i)}
                <div class="proc-row" class:done={st.done}>
                  <button class="proc-dot" class:done={st.done} onclick={() => toggleStage(i)} aria-label={st.done ? 'Mark not done' : 'Mark done'}>
                    {#if st.done}✓{/if}
                  </button>
                  <span class="proc-name">{st.name}</span>
                </div>
              {/each}
            </div>
            <div class="side-note">{pipelineDone} of {pipeline.length} done · as you mapped it</div>
            {#if stageDebriefPrompt}
              <div class="stage-debrief">
                <span class="sd-tx">Just did the <b>{stageDebriefPrompt.name}</b> round?</span>
                <div class="sd-actions">
                  <button type="button" class="sd-x" onclick={() => (stageDebriefPrompt = null)}>Not yet</button>
                  <button type="button" class="sd-go" onclick={debriefFromStage}>Debrief it →</button>
                </div>
              </div>
            {/if}
          {:else}
            <p class="side-empty">No stages yet — map the steps the recruiter described.</p>
            <button class="add-line" onclick={seedTypicalLoop}>+ Start from a typical loop</button>
          {/if}
        </div>

        <div class="side-sec">
          <div class="side-lbl">People</div>
          <div class="people">
            {#if app.raw.recruiter_name}
              <div class="pp-row">
                <span class="pp-av warm">{recruiterInitials || '—'}</span>
                <span class="pp-main"><strong>{app.raw.recruiter_name}</strong> <span class="pp-role">· recruiter</span></span>
                {#if app.raw.recruiter_email}
                  <a class="pp-link" href={`mailto:${app.raw.recruiter_email}`}>Email →</a>
                {:else if app.raw.recruiter_linkedin}
                  <a class="pp-link" href={app.raw.recruiter_linkedin} target="_blank" rel="noopener">LinkedIn →</a>
                {/if}
              </div>
              {#if app.raw.recruiter_email && app.raw.recruiter_linkedin}
                <div class="pp-extra"><a class="pp-link" href={app.raw.recruiter_linkedin} target="_blank" rel="noopener">LinkedIn →</a></div>
              {/if}
            {/if}
            {#if app.raw.hiring_manager_name}
              <div class="pp-row">
                <span class="pp-av">{hiringManagerInitials || '—'}</span>
                <span class="pp-main"><strong>{app.raw.hiring_manager_name}</strong> <span class="pp-role">· hiring manager</span></span>
                {#if app.raw.hiring_manager_linkedin}
                  <a class="pp-link" href={app.raw.hiring_manager_linkedin} target="_blank" rel="noopener">LinkedIn →</a>
                {/if}
              </div>
            {/if}
            {#if !app.raw.recruiter_name && !app.raw.hiring_manager_name}
              <p class="side-empty">No contacts yet — add the recruiter or hiring manager.</p>
            {/if}
            <button class="linkbtn add-contact" onclick={openEdit}>+ Add a contact</button>
          </div>
        </div>

        <div class="side-sec">
          <div class="side-hd-row">
            <span class="side-lbl">Activity</span>
            <span class="side-acts">
              <button class="side-act linkbtn" onclick={openFollowUp}>+ Follow-up</button>
              <button class="side-act linkbtn" onclick={openEdit}>+ Note</button>
              <button class="side-act linkbtn" onclick={openAddEvent}>+ Interview</button>
            </span>
          </div>
          {#if timeline.length > 0}
            <div class="tl">
              {#each timeline as e (e.ts + e.title)}
                <div class="tl-row">
                  <span class="tl-date">{e.date}</span>
                  <span class="tl-body">
                    <span class="tl-title">{e.title}{#if e.note}<span class="tl-note"> — {e.note}</span>{/if}</span>
                    {#if e.followUp}
                      <button class="tl-del" title="Delete follow-up" aria-label="Delete follow-up" onclick={() => deleteFollowUp(e.followUp)}>✕</button>
                    {:else if e.interview}
                      <button class="tl-del" title="Delete interview" aria-label="Delete interview" onclick={() => deleteInterview(e.interview)}>✕</button>
                    {/if}
                  </span>
                </div>
              {/each}
            </div>
          {:else}
            <p class="side-empty">No activity yet.</p>
          {/if}
        </div>
      </aside>
    </div>
  {/if}
</div>

{#if toast}
  <div class="toast"><span class="ok">✓</span>{toast}</div>
{/if}

{#if showEditModal}
  <div class="modal-overlay" onclick={() => (showEditModal = false)} role="presentation">
    <form class="modal" onclick={(e) => e.stopPropagation()} onsubmit={saveEdit}>
      <h2>Edit application</h2>
      <p class="modal-hint">Empty fields keep their current value.</p>
      <div class="fields">
        <label>Company <input bind:value={edit.company} required /></label>
        <label>Role <input bind:value={edit.role} required /></label>
        <label class="span-2">Source
          <input bind:value={edit.source} list={coarsePointer ? undefined : 'edit-source-suggestions'} placeholder="LinkedIn / Referral / Cold email" />
          <div class="src-chips">
            {#each SOURCE_SUGGESTIONS as s (s)}<button type="button" class="src-pick" onclick={() => (edit.source = s)}>{s}</button>{/each}
          </div>
        </label>
        <datalist id="edit-source-suggestions">
          {#each SOURCE_SUGGESTIONS as s (s)}<option value={s}></option>{/each}
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
      <p class="privacy-note">Private to your account. Your notes and salary info are never shared or shown to anyone else.</p>
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
  <div class="modal-overlay" onclick={closeEventModal} role="presentation">
    <div class="modal ev-card" onclick={(e) => e.stopPropagation()} role="dialog" aria-modal="true" aria-label="Add an interview">
      <button class="x-close" onclick={closeEventModal} aria-label="Close">✕</button>
      <div class="add-hd">
        <h3>Add an interview</h3>
        <p>Paste the invite, drop a screenshot or .ics file, or just type the details — we'll pull out the event.</p>
      </div>

      <div class="ev-input" class:drag={evDragOver} class:loading={icsParsing}
        ondragover={onEvDragOver} ondragleave={() => (evDragOver = false)} ondrop={onEvDrop} role="presentation">
        {#if evAttach}
          <div class="ev-attached">
            <span class="ev-att-name">{evAttach.name}</span>
            <span class="ev-att-kind">{evAttach.kind === 'image' ? 'screenshot' : 'calendar file'} · {Math.round(evAttach.size / 1024)} KB</span>
            <button type="button" class="ev-att-x" onclick={() => (evAttach = null)} aria-label="Remove" disabled={icsParsing}>×</button>
          </div>
        {:else}
          <textarea class="ev-ta" rows="3" bind:value={evText} onpaste={onEvPaste}
            placeholder={"Paste an invite, or type it — e.g. “Interview Wed Jun 10, 11:00, Google Meet”"}
            disabled={icsParsing}></textarea>
        {/if}
        <div class="ev-foot">
          <label class="ev-browse">
            <span>↓ Drop a screenshot or .ics, paste, or <u>browse</u></span>
            <input type="file" accept=".ics,text/calendar,image/png,image/jpeg,image/gif,image/webp" onchange={onEvFileInput} hidden />
          </label>
          {#if icsParsing}<span class="ev-loading"><span class="ev-spin" aria-hidden="true"></span> Reading the event…</span>{/if}
        </div>
      </div>

      <div class="ev-actions">
        <span class="ev-hint"><kbd>⌘V</kbd> pastes a screenshot</span>
        <button class="btn btn-primary" onclick={parseEvent} disabled={icsParsing || (!evText.trim() && !evAttach)}>
          {icsParsing ? 'Reading…' : 'Find the event'}
        </button>
      </div>

      {#if icsParseError}<p class="parse-err">{icsParseError}</p>{/if}

      {#if icsPreview.length > 0}
        <div class="ics-preview">
          <h4>Preview</h4>
          <p class="prev-check">Double-check the day and time before saving.</p>
          {#each icsPreview as ev, i (i)}
            <div class="prev-row">
              <div class="prev-summary">{ev.summary || 'Untitled event'}</div>
              <div class="prev-when"><strong>{fmtEventDay(ev)}</strong>{fmtEventTimeSuffix(ev)}</div>
              {#if ev.location}<div class="prev-loc">{ev.location}</div>{/if}
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
  <div class="modal-overlay" onclick={closeFollowUp} role="presentation">
    <form class="modal fu-card" onclick={(e) => e.stopPropagation()} onsubmit={saveFollowUp} role="dialog" aria-modal="true" aria-label="Log a follow-up">
      <button type="button" class="x-close" onclick={closeFollowUp} aria-label="Close">✕</button>
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
              {#each FU_CHANNELS as c (c)}<option value={c}>{c}</option>{/each}
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
  .wrap {
    max-width: 1160px; width: 100%; box-sizing: border-box;
    margin: 0 auto; padding: 26px 32px 80px;
    color: #16181c;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
  }
  .wrap :global(a) { color: #2463eb; text-decoration: none; }
  .loading { color: #8a9099; font-size: 13px; }
  .linkbtn { background: none; border: 0; cursor: pointer; padding: 0; font-family: inherit; font-size: 13px; color: #4b5158; }
  .linkbtn.danger { color: #b3372a; }
  .linkbtn:hover { color: #16181c; }
  .linkbtn.danger:hover { color: #b3372a; text-decoration: underline; }

  .crumb { font-size: 13px; color: #8a9099; margin-bottom: 20px; }
  .crumb a { color: #8a9099; }
  .crumb a:hover { color: #4b5158; }
  .crumb .sep { margin: 0 4px; }
  .crumb .here { color: #16181c; font-weight: 600; }

  .welcome-banner { display: flex; align-items: center; gap: 11px; margin-bottom: 22px;
    padding: 12px 14px; border-radius: 12px; background: #eef4ff; border: 1px solid #cdddfb; }
  .welcome-banner .wb-spark { color: #2463eb; font-size: 14px; flex-shrink: 0; }
  .welcome-banner .wb-tx { font-size: 13.5px; color: #4b5158; line-height: 1.5; }
  .welcome-banner .wb-x { margin-left: auto; flex-shrink: 0; background: none; border: none; color: #8a9099;
    font-size: 13px; cursor: pointer; padding: 4px 6px; border-radius: 6px; }
  .welcome-banner .wb-x:hover { background: #fff; color: #4b5158; }

  /* HEADER */
  .hd { display: flex; align-items: flex-start; gap: 18px; margin-bottom: 22px; }
  .hd-logo { flex: none; }
  .hd-main { flex: 1; min-width: 0; }
  .hd-row { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; }
  .hd-row h1 { font-size: 28px; font-weight: 700; letter-spacing: -0.02em; margin: 0; }
  .hd-role { font-size: 14.5px; color: #4b5158; margin: 3px 0 6px; }
  .hd-sub { font-size: 12.5px; color: #8a9099; display: flex; align-items: baseline; gap: 6px; flex-wrap: wrap; }
  .hd-sub .dot { color: #b8bdc4; }
  .hd-actions { display: flex; align-items: center; gap: 14px; flex: none; padding-top: 8px; }

  /* SOON BANNER */
  .soon { display: flex; align-items: center; gap: 16px; background: #fff7f1; border: 1px solid #f0d9c4;
    border-radius: 14px; padding: 16px 22px; margin-bottom: 28px; }
  .soon-cal { flex: none; width: 44px; border: 1px solid #f0d9c4; border-radius: 9px; overflow: hidden; text-align: center; background: #fff; }
  .sc-mon { background: #e0641f; color: #fff; font-size: 9px; font-weight: 700; letter-spacing: .08em; padding: 2px 0; }
  .sc-day { font-size: 17px; font-weight: 700; color: #c05310; padding: 2px 0 3px; }
  .soon-tx { flex: 1; min-width: 0; }
  .soon-t { font-size: 15px; font-weight: 700; }
  .soon-s { font-size: 13px; color: #6f7680; }
  .soon-cta { background: #2463eb; color: #fff; border: 0; border-radius: 9px; padding: 10px 18px;
    font-size: 13.5px; font-weight: 600; cursor: pointer; flex: none; font-family: inherit; }
  .soon-cta:hover { background: #1a4fc4; }

  /* GRID */
  .grid { display: grid; grid-template-columns: 1.9fr 1fr; gap: 40px; align-items: start; }
  .col-main { min-width: 0; }

  /* PLAYBOOK CARD */
  .pb-card { background: #fff; border: 1px solid #e8e8e5; border-radius: 16px; overflow: hidden;
    box-shadow: 0 1px 3px rgba(22,24,28,.04); }

  /* Identity strip */
  .idbar { display: flex; align-items: center; gap: 12px; background: #fbfbf9; border-bottom: 1px solid #eeeeea; padding: 11px 20px; }
  .id-tx { font-size: 12.5px; color: #4b5158; min-width: 0; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  .id-sum { color: #8a9099; }
  .id-not { margin-left: auto; flex: none; font-size: 12px; color: #8a9099; }
  .idbox { background: #fdf6ef; border-bottom: 1px solid #f0d9c4; padding: 18px 22px; }
  .id-q { font-size: 14px; font-weight: 700; margin-bottom: 12px; }
  .id-card { display: flex; align-items: center; gap: 14px; background: #fff; border: 1px solid #eeeeea; border-radius: 12px; padding: 14px 18px; margin-bottom: 12px; }
  .id-card-tx { flex: 1; min-width: 0; font-size: 13px; line-height: 1.5; }
  .id-actions { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
  .id-hint { font-size: 12px; color: #8a9099; }
  .btn-blue { background: #2463eb; color: #fff; border: 0; border-radius: 8px; padding: 8px 16px;
    font-size: 13px; font-weight: 600; cursor: pointer; font-family: inherit; white-space: nowrap; }
  .btn-blue:hover:not(:disabled) { background: #1a4fc4; }
  .btn-blue:disabled { opacity: .5; cursor: default; }
  .btn-ghost-warm { border: 1px solid #e2d4c4; background: #fff; color: #4b5158; border-radius: 8px; padding: 8px 16px;
    font-size: 13px; font-weight: 600; cursor: pointer; font-family: inherit; white-space: nowrap; }
  .id-fix-sub { font-size: 12.5px; color: #6f7680; margin: -8px 0 12px; }
  .id-fix-row { display: flex; align-items: center; gap: 10px; }
  .id-input { flex: 1; min-width: 0; border: 1px solid #e2d4c4; border-radius: 8px; padding: 9px 13px;
    font-size: 13px; background: #fff; color: #16181c; outline: none; font-family: inherit; }
  .id-cancel { font-size: 12.5px; color: #8a9099; flex: none; }
  .idok { display: flex; align-items: center; gap: 10px; background: #f3faf4; border-bottom: 1px solid #cfe5d2; padding: 10px 20px; }
  .idok span { color: #1d7a4f; font-size: 13px; font-weight: 600; }
  .id-change { margin-left: auto; font-size: 12px; color: #8a9099; }

  /* Playbook header + tabs */
  .pb-hd { padding: 18px 22px 0; }
  .pb-title-row { display: flex; align-items: center; gap: 10px; margin-bottom: 14px; flex-wrap: wrap; }
  .pb-title { font-size: 19px; font-weight: 700; letter-spacing: -0.01em; }
  .pb-ago { font-size: 12px; color: #8a9099; }
  .informed-chip { display: inline-flex; align-items: center; gap: 6px; font-size: 11.5px; font-weight: 500;
    color: #1d7a4f; background: #f3faf4; border: 1px solid #cfe5d2; border-radius: 99px; padding: 3px 10px; }
  .informed-chip .ic-dot { width: 5px; height: 5px; border-radius: 50%; background: #16a34a; }
  .tabs { display: flex; align-items: center; gap: 7px; padding-bottom: 14px; border-bottom: 1px solid #f0f0ed; flex-wrap: wrap; }
  .tab { display: flex; align-items: center; gap: 6px; border-radius: 9px; padding: 8px 15px; font-size: 13px;
    cursor: pointer; font-family: inherit; background: #fff; border: 1px solid #e8e8e5; color: #4b5158; font-weight: 500; }
  .tab.done { color: #1d7a4f; }
  .tab.company.active { background: #eef4ff; border-color: #cdddfb; color: #2463eb; font-weight: 700; }
  .tab.round.active { background: #fff7f1; border-color: #f0d9c4; color: #c05310; font-weight: 700; }
  .tab.add { border: 1px dashed #e2e2de; color: #b8bdc4; }
  .tab.add:hover { color: #4b5158; border-color: #b8bdc4; }

  .pb-body { padding: 20px 22px 24px; }

  /* One-tap add round */
  .add-round { margin-bottom: 18px; padding: 13px 15px; border: 1px solid #e8e8e5; border-radius: 12px;
    background: #fbfbf9; display: flex; flex-direction: column; gap: 11px; }
  .ar-lbl { font-size: 12.5px; color: #4b5158; font-weight: 500; }
  .ar-chips { display: flex; flex-wrap: wrap; gap: 7px; }
  .ar-chip { font: inherit; font-size: 12.5px; color: #4b5158; background: #fff; border: 1px solid #e8e8e5;
    border-radius: 8px; padding: 6px 12px; cursor: pointer; }
  .ar-chip:hover { border-color: #cdddfb; color: #2463eb; background: #eef4ff; }
  .ar-custom { display: flex; gap: 8px; }
  .ar-input { flex: 1; min-width: 0; font: inherit; font-size: 13px; color: #16181c; background: #fff;
    border: 1px solid #e8e8e5; border-radius: 8px; padding: 7px 10px; outline: none; }
  .ar-input:focus { border-color: #2463eb; }

  /* Debrief */
  .db-banner { width: 100%; display: flex; align-items: center; gap: 10px; text-align: left; cursor: pointer;
    background: #eef4ff; border: 1px solid #cdddfb; border-radius: 10px; padding: 9px 13px; font-family: inherit; margin-bottom: 14px; }
  .db-banner:hover { border-color: #2463eb; }
  .db-spark { width: 24px; height: 24px; border-radius: 7px; background: #fff; color: #2463eb;
    display: grid; place-items: center; flex-shrink: 0; font-size: 12px; }
  .db-banner-tx { flex: 1; min-width: 0; font-size: 12.5px; color: #4b5158; }
  .db-banner-tx b { color: #16181c; }
  .db-prompt-cta { flex-shrink: 0; font-size: 13px; font-weight: 600; color: #2463eb; }
  .debrief { margin-bottom: 18px; }
  .db-prompt { width: 100%; display: flex; align-items: center; gap: 12px; text-align: left; cursor: pointer;
    background: #eef4ff; border: 1px solid #cdddfb; border-radius: 12px; padding: 13px 15px; font-family: inherit; }
  .db-prompt:hover { border-color: #2463eb; }
  .db-prompt-tx { flex: 1; min-width: 0; }
  .db-prompt-tx b { display: block; font-size: 13.5px; color: #16181c; }
  .db-prompt-tx small { display: block; font-size: 12px; color: #8a9099; margin-top: 1px; }
  .db-form { border: 1px solid #e8e8e5; border-radius: 12px; padding: 15px 16px; background: #fbfbf9;
    display: flex; flex-direction: column; gap: 13px; }
  .db-hd { display: flex; align-items: baseline; gap: 10px; flex-wrap: wrap; }
  .db-hd h3 { margin: 0; font-size: 15px; font-weight: 600; color: #16181c; }
  .db-hint { font-size: 12px; color: #8a9099; }
  .db-q { display: flex; flex-direction: column; gap: 7px; }
  .db-q-lbl { font-size: 12px; font-weight: 500; color: #4b5158; }
  .db-opts { display: flex; gap: 7px; flex-wrap: wrap; }
  .db-opt { font: inherit; font-size: 13px; font-weight: 500; color: #4b5158; background: #fff;
    border: 1px solid #e8e8e5; border-radius: 8px; padding: 7px 14px; cursor: pointer; }
  .db-opt:hover { border-color: #b8bdc4; }
  .db-opt.sel { background: #2463eb; border-color: #1a4fc4; color: #fff; }
  .db-input { font: inherit; font-size: 13.5px; color: #16181c; background: #fff; border: 1px solid #e8e8e5;
    border-radius: 8px; padding: 9px 11px; outline: none; }
  .db-input:focus { border-color: #2463eb; }
  .db-actions { display: flex; justify-content: flex-end; gap: 8px; }
  .db-summary { display: flex; align-items: center; gap: 12px; border: 1px solid #e8e8e5; border-radius: 11px;
    padding: 11px 14px; background: #fbfbf9; }
  .db-badge { font-size: 10.5px; font-weight: 700; letter-spacing: 0.05em; text-transform: uppercase;
    color: #1d7a4f; background: #eef7ef; padding: 3px 8px; border-radius: 6px; }
  .db-sum-main { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; min-width: 0; }
  .db-sum-line { font-size: 13px; color: #4b5158; }
  .db-sum-line b { color: #16181c; font-weight: 600; }
  .db-edit { margin-left: auto; flex-shrink: 0; background: none; border: 1px solid #e8e8e5; color: #4b5158;
    font: 500 12px/1 inherit; font-family: inherit; padding: 6px 11px; border-radius: 7px; cursor: pointer; }
  .db-edit:hover { border-color: #b8bdc4; }
  .db-topics { font-size: 12.5px; color: #8a9099; margin: 8px 0 0; line-height: 1.5; }

  /* Generate / build / gate box */
  .genbox { text-align: center; padding: 26px 18px 22px; }
  .genbox h3 { font-size: 18px; font-weight: 600; letter-spacing: -0.015em; margin: 0 0 10px; }
  .gen-sub { font-size: 13.5px; color: #6f7680; line-height: 1.6; margin: 0 auto 20px; max-width: 46ch; }
  .gen-row { display: flex; gap: 10px; flex-direction: column; max-width: 420px; margin: 0 auto; }
  .gen-input { width: 100%; padding: 11px 14px; font-size: 13.5px; font-family: inherit; border: 1px solid #e8e8e5;
    border-radius: 9px; background: #fbfbf9; color: #16181c; outline: none; box-sizing: border-box; }
  .gen-input:focus { border-color: #2463eb; background: #fff; }
  .gen-input::placeholder { color: #b8bdc4; }
  .btn-generate { background: #2463eb; color: #fff; border: none; border-radius: 9px; padding: 12px 20px;
    font-size: 14px; font-weight: 600; font-family: inherit; cursor: pointer; }
  .btn-generate:hover:not(:disabled) { background: #1a4fc4; }
  .btn-generate:disabled { opacity: 0.5; cursor: default; }
  .genbox > .btn-generate { min-width: 240px; }
  .gen-err { color: #b3372a; font-size: 13px; margin: 14px 0 0; }
  .credit-left { color: #8a9099; font-size: 12.5px; margin: 10px 0 0; }
  .big-spinner { width: 36px; height: 36px; border: 2.5px solid #e2e2de; border-top-color: #2463eb;
    border-radius: 50%; animation: prep-spin 0.75s linear infinite; margin: 24px auto 0; }
  .gen-eta { font-size: 12px; color: #8a9099; margin: 16px auto 0; max-width: 44ch; }
  @keyframes prep-spin { to { transform: rotate(360deg); } }

  /* Company brief content */
  .cb-hd { font-size: 15px; font-weight: 700; margin-bottom: 12px; }
  .cb-shared { font-size: 12px; font-weight: 400; color: #8a9099; }
  .cb-blurb { font-size: 14px; font-weight: 600; line-height: 1.55; margin: 0 0 8px; }
  .cb-about { font-size: 13.5px; line-height: 1.6; color: #4b5158; margin: 0 0 16px; }
  .cb-facts { display: grid; grid-template-columns: repeat(4, 1fr); border: 1px solid #eeeeea; border-radius: 10px;
    overflow: hidden; margin-bottom: 16px; }
  .cb-cell { padding: 11px 13px; }
  .cb-cell + .cb-cell { border-left: 1px solid #eeeeea; }
  .f-lbl { font-size: 11px; color: #8a9099; font-weight: 500; }
  .f-val { font-size: 13px; color: #16181c; font-weight: 600; margin-top: 3px; line-height: 1.3; }
  .cb-proc { margin-bottom: 16px; }
  .cb-lbl { font-size: 11px; font-weight: 600; letter-spacing: .12em; text-transform: uppercase; color: #8a9099; margin-bottom: 9px; }
  .cb-chips { display: flex; align-items: center; flex-wrap: wrap; gap: 7px; }
  .chip { font-size: 12px; font-weight: 500; color: #4b5158; background: #fbfbf9; border: 1px solid #e8e8e5;
    border-radius: 7px; padding: 4px 10px; }
  .chip-arrow { color: #b8bdc4; font-size: 12px; }
  .grades { background: #eef4ff; border: 1px solid #cdddfb; border-radius: 12px; padding: 16px 20px; margin-bottom: 14px; }
  .grades-hd { font-size: 13px; font-weight: 700; color: #2463eb; margin-bottom: 8px; }
  .grades-list { display: flex; flex-direction: column; gap: 7px; font-size: 13.5px; line-height: 1.55; color: #1e3a6e; }
  .srcs { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; font-size: 12px; color: #8a9099; }
  .srcs-lbl { font-weight: 600; color: #6f7680; }
  .src-chip { border: 1px solid #e8e8e5; border-radius: 14px; padding: 3px 10px; color: #2463eb; }
  .src-chip:hover { border-color: #cdddfb; background: #eef4ff; }

  /* Round tab — door */
  .rd-hd { font-size: 15px; font-weight: 700; margin-bottom: 14px; }
  .rd-when { font-size: 12.5px; font-weight: 400; color: #8a9099; }
  .pcards { display: grid; grid-template-columns: repeat(auto-fill, minmax(160px, 200px)); gap: 10px; margin-bottom: 16px; }
  .pcard { border: 1px solid #eeeeea; border-radius: 12px; padding: 14px; text-align: center; color: #16181c !important; display: block; }
  .pcard:hover { border-color: #b9c6e8; }
  .pc-av { width: 38px; height: 38px; border-radius: 50%; background: #eef4ff; color: #2463eb;
    display: flex; align-items: center; justify-content: center; font-size: 13px; font-weight: 700; margin: 0 auto 8px; }
  .pc-name { font-size: 13.5px; font-weight: 700; }
  .pc-role { font-size: 11.5px; color: #8a9099; }
  .doorbar { display: flex; align-items: center; gap: 14px; background: #2463eb; border-radius: 12px;
    padding: 16px 20px; cursor: pointer; color: #fff !important; }
  .doorbar:hover { background: #1a4fc4; }
  .db-main { flex: 1; min-width: 0; }
  .db-t { font-size: 14.5px; font-weight: 700; }
  .db-s { font-size: 12px; opacity: .75; }
  .db-arrow { font-size: 20px; }

  /* Footnote */
  .pb-foot { display: flex; align-items: center; gap: 8px; margin-top: 10px; font-size: 12px; color: #8a9099; flex-wrap: wrap; }
  .pb-foot strong { color: #16181c; }
  .foot-refresh { font-size: 12px; color: #2463eb; }
  .foot-refresh:hover { text-decoration: underline; color: #1a4fc4; }
  .foot-refresh:disabled { opacity: .5; cursor: default; }
  .foot-sep { color: #d8dade; }
  .foot-remove { margin-left: auto; font-size: 12px; color: #b8bdc4; }
  .foot-remove:hover { color: #b3372a; }

  /* RIGHT COLUMN */
  .col-side { font-size: 13.5px; min-width: 0; }
  .side-sec { border-top: 1px solid #e2e2de; padding-top: 18px; margin-bottom: 24px; }
  .side-sec.first { border-top: 0; padding-top: 0; }
  .side-lbl { font-size: 11px; font-weight: 600; letter-spacing: .12em; text-transform: uppercase; color: #8a9099; display: block; margin-bottom: 10px; }
  .side-hd-row { display: flex; align-items: baseline; margin-bottom: 10px; gap: 10px; }
  .side-hd-row .side-lbl { margin-bottom: 0; }
  .side-act { font-size: 12px; color: #2463eb; margin-left: auto; }
  .side-acts { margin-left: auto; display: flex; gap: 12px; }
  .side-acts .side-act { margin-left: 0; }
  .side-act:hover { text-decoration: underline; color: #1a4fc4; }
  .side-prose { font-size: 13px; line-height: 1.65; color: #4b5158; }
  .side-note { font-size: 11.5px; color: #b8bdc4; margin-top: 10px; }
  .side-edit { font-size: 11.5px; color: #8a9099; }
  .side-edit:hover { color: #4b5158; text-decoration: underline; }
  .side-empty { font-size: 12.5px; color: #8a9099; margin: 0 0 10px; line-height: 1.5; }
  .jd-saved { margin-top: 10px; }
  .jd-saved summary { font-size: 12.5px; font-weight: 500; color: #2463eb; cursor: pointer; list-style: none; }
  .jd-saved summary::-webkit-details-marker { display: none; }
  .jd-saved summary:hover { text-decoration: underline; }
  .jd-body { margin: 10px 0 0; font-size: 12.5px; line-height: 1.55; color: #4b5158; white-space: pre-wrap;
    max-height: 280px; overflow-y: auto; padding: 10px 12px; background: #fff; border: 1px solid #e8e8e5; border-radius: 8px; }

  /* Process */
  .proc-list { display: flex; flex-direction: column; gap: 9px; font-size: 13px; }
  .proc-row { display: flex; align-items: center; gap: 10px; color: #4b5158; }
  .proc-row.done { color: #1d7a4f; }
  .proc-row.done .proc-name { text-decoration: line-through; color: #8a9099; }
  .proc-dot { width: 20px; height: 20px; border-radius: 50%; border: 1.5px solid #d8dade; background: #fff;
    display: flex; align-items: center; justify-content: center; font-size: 11px; flex: none; cursor: pointer;
    color: #1d7a4f; padding: 0; font-family: inherit; }
  .proc-dot:hover { border-color: #2463eb; }
  .proc-dot.done { background: #eef7ef; border-color: #cfe5d2; }
  .pipe-edit { display: flex; flex-direction: column; gap: 8px; }
  .pe-row { display: grid; grid-template-columns: auto 1fr auto auto auto; gap: 4px; align-items: center; border-radius: 8px; }
  .pe-row.pe-dragging { opacity: 0.5; background: #eef4ff; }
  .pe-grip { display: grid; place-items: center; width: 20px; height: 30px; color: #b8bdc4; cursor: grab; font-size: 14px; user-select: none; }
  .pe-grip:active { cursor: grabbing; }
  .pe-grip:hover { color: #4b5158; }
  .pe-input { font: inherit; font-size: 13px; color: #16181c; background: #fff; border: 1px solid #e8e8e5;
    border-radius: 7px; padding: 6px 9px; outline: none; min-width: 0; }
  .pe-input:focus { border-color: #2463eb; }
  .pe-btn { width: 26px; height: 30px; display: grid; place-items: center; background: #fbfbf9; border: 1px solid #e8e8e5;
    border-radius: 7px; color: #8a9099; font-size: 13px; cursor: pointer; font-family: inherit; }
  .pe-btn:hover:not(:disabled) { border-color: #b8bdc4; color: #16181c; }
  .pe-btn:disabled { opacity: 0.4; cursor: default; }
  .pe-x:hover:not(:disabled) { color: #b3372a; border-color: #f2d4cf; }
  .pe-actions { display: flex; justify-content: flex-end; gap: 8px; margin-top: 4px;
    position: sticky; bottom: 0; z-index: 2; background: #f6f6f3; padding: 8px 0 2px; }
  .add-line { background: none; border: 0; font: inherit; font-size: 12px; color: #8a9099; cursor: pointer;
    padding: 0; text-align: left; }
  .add-line:hover { color: #2463eb; }
  .stage-debrief { margin-top: 12px; padding: 10px 12px; border-radius: 10px; background: #eef4ff;
    border: 1px solid #cdddfb; display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
  .stage-debrief .sd-tx { font-size: 12.5px; color: #4b5158; }
  .stage-debrief .sd-tx b { color: #16181c; }
  .stage-debrief .sd-actions { margin-left: auto; display: flex; gap: 6px; flex-shrink: 0; }
  .stage-debrief .sd-x { background: none; border: none; color: #8a9099; font-family: inherit;
    font-size: 12px; font-weight: 500; padding: 5px 8px; cursor: pointer; }
  .stage-debrief .sd-go { background: #2463eb; color: #fff; border: none; border-radius: 7px; font-family: inherit;
    font-size: 12px; font-weight: 600; padding: 6px 11px; cursor: pointer; }
  .stage-debrief .sd-go:hover { background: #1a4fc4; }

  /* People */
  .people { display: flex; flex-direction: column; gap: 10px; font-size: 13px; }
  .pp-row { display: flex; align-items: center; gap: 10px; }
  .pp-av { width: 26px; height: 26px; border-radius: 50%; background: #eef4ff; color: #2463eb;
    display: flex; align-items: center; justify-content: center; font-size: 11px; font-weight: 700; flex: none; }
  .pp-av.warm { background: #fdf3ec; color: #c05310; }
  .pp-main { flex: 1; min-width: 0; }
  .pp-role { color: #8a9099; }
  .pp-link { font-size: 12px; flex: none; }
  .pp-extra { padding-left: 36px; margin-top: -4px; }
  .add-contact { font-size: 12px; color: #8a9099; text-align: left; }
  .add-contact:hover { color: #2463eb; }

  /* Activity */
  .tl { display: flex; flex-direction: column; gap: 9px; font-size: 12.5px; color: #4b5158; }
  .tl-row { display: flex; gap: 10px; align-items: baseline; }
  .tl-date { color: #b8bdc4; width: 44px; flex: none; font-variant-numeric: tabular-nums; }
  .tl-body { display: flex; align-items: baseline; gap: 6px; min-width: 0; flex: 1; }
  .tl-title { line-height: 1.45; }
  .tl-note { color: #8a9099; }
  .tl-del { flex: none; border: 0; background: transparent; color: #b8bdc4; border-radius: 5px; font-size: 11px;
    cursor: pointer; opacity: 0; padding: 1px 4px; font-family: inherit; }
  .tl-row:hover .tl-del { opacity: 1; }
  .tl-del:hover { background: #fdf1ef; color: #b3372a; }

  /* Buttons */
  .btn { font: inherit; font-size: 13px; font-weight: 500; color: #16181c; background: #fff;
    border: 1px solid #e8e8e5; border-radius: 8px; padding: 7px 14px; cursor: pointer; }
  .btn:hover:not(:disabled) { border-color: #b8bdc4; }
  .btn:disabled { opacity: 0.5; cursor: default; }
  .btn-primary { background: #2463eb; border-color: #2463eb; color: #fff; font-weight: 600; }
  .btn-primary:hover:not(:disabled) { background: #1a4fc4; border-color: #1a4fc4; }

  /* TOAST */
  .toast { position: fixed; bottom: 26px; left: 50%; transform: translateX(-50%); background: #16181c; color: #fff;
    font-size: 13px; font-weight: 500; padding: 11px 18px; border-radius: 10px; z-index: 200;
    display: inline-flex; align-items: center; gap: 9px; box-shadow: 0 16px 36px -12px rgba(20,20,50,0.5); animation: rise .2s ease; }
  .toast .ok { width: 16px; height: 16px; border-radius: 50%; background: #16a34a; color: #fff; font-size: 10px;
    display: inline-flex; align-items: center; justify-content: center; }
  @keyframes rise { from { transform: translate(-50%, 12px); opacity: 0; } }

  /* MODALS */
  .modal-overlay { position: fixed; inset: 0; background: rgba(10,10,13,0.45); backdrop-filter: blur(6px);
    -webkit-backdrop-filter: blur(6px); display: grid; place-items: center; z-index: 200; padding: 24px; overflow-y: auto; }
  .modal { position: relative; background: #fff; border: 1px solid #e8e8e5; border-radius: 16px; padding: 24px 26px 22px;
    width: 100%; max-width: 560px; max-height: calc(100dvh - 4rem); overflow-y: auto; display: flex; flex-direction: column;
    gap: .75rem; box-shadow: 0 24px 80px -8px rgba(10,10,13,0.30); box-sizing: border-box;
    color: #16181c; font-family: inherit; }
  .modal h2 { font-size: 18px; font-weight: 700; letter-spacing: -0.018em; margin: 0; }
  .modal-hint { font-size: 12px; color: #8a9099; margin: 0 0 .5rem; }
  .fields { display: grid; grid-template-columns: 1fr 1fr; gap: .65rem; }
  .fields .span-2 { grid-column: span 2; }
  .field-group { grid-column: span 2; font-size: 11px; font-weight: 600; letter-spacing: .04em; text-transform: uppercase;
    color: #4b5158; margin-top: .5rem; padding-top: .65rem; border-top: 1px solid #eeeeea; }
  .field-group .fg-sub { font-weight: 400; text-transform: none; letter-spacing: 0; color: #b8bdc4; }
  .modal label { display: flex; flex-direction: column; font-size: 12px; color: #8a9099; gap: .35rem; }
  .modal input, .modal select { font: inherit; color: #16181c; background: #fbfbf9; border: 1px solid #e8e8e5;
    border-radius: 6px; padding: .45rem .6rem; font-size: 13.5px; outline: none; }
  .modal input:focus, .modal select:focus { border-color: #2463eb; }
  .modal .jd-area { font: inherit; color: #16181c; background: #fbfbf9; border: 1px solid #e8e8e5; border-radius: 6px;
    padding: .45rem .6rem; font-size: 13.5px; line-height: 1.5; outline: none; resize: vertical; min-height: 72px; }
  .modal .jd-area:focus { border-color: #2463eb; }
  .src-chips { display: flex; flex-wrap: wrap; gap: 6px; margin-top: 7px; }
  .src-pick { font: inherit; font-size: 11.5px; color: #4b5158; background: #fff; border: 1px solid #e8e8e5;
    border-radius: 7px; padding: 4px 9px; cursor: pointer; }
  .src-pick:hover { border-color: #cdddfb; color: #2463eb; background: #eef4ff; }
  /* Sticky footer so Cancel/Save stay reachable when the mobile keyboard is up. */
  .modal-actions { display: flex; justify-content: flex-end; gap: .5rem;
    position: sticky; bottom: 0; background: #fff; padding-top: .7rem; margin-top: .5rem; }
  .privacy-note { font-size: 11.5px; color: #8a9099; margin: 12px 0 0; line-height: 1.4; }
  .x-close { position: absolute; top: 14px; right: 14px; width: 28px; height: 28px; border-radius: 8px;
    background: transparent; border: 0; display: grid; place-items: center; color: #8a9099; cursor: pointer;
    font-size: 13px; font-family: inherit; }
  .x-close:hover { background: #f0f0ed; color: #16181c; }
  .add-hd { margin-bottom: 14px; padding-right: 26px; }
  .add-hd h3 { font-size: 15px; font-weight: 700; margin: 0 0 4px; letter-spacing: -0.015em; }
  .add-hd p { font-size: 13px; color: #8a9099; margin: 0; line-height: 1.5; }

  /* Add-event modal internals */
  .ev-card { max-width: 540px; }
  .ev-input { position: relative; background: #fbfbf9; border: 1.5px dashed #d8dade; border-radius: 12px; overflow: hidden; }
  .ev-input.drag { border-color: #2463eb; background: #eef4ff; }
  .ev-input.loading { border-style: solid; border-color: #2463eb; }
  .ev-ta { width: 100%; font: inherit; font-size: 13.5px; line-height: 1.55; color: #16181c; background: transparent;
    border: 0; padding: 13px 15px 8px; outline: none; resize: none; min-height: 78px; box-sizing: border-box; display: block; }
  .ev-ta::placeholder { color: #b8bdc4; }
  .ev-attached { display: flex; align-items: center; gap: 10px; padding: 14px 15px 10px; font-size: 13px; }
  .ev-att-name { font-weight: 500; color: #16181c; }
  .ev-att-kind { color: #8a9099; font-size: 11.5px; }
  .ev-att-x { margin-left: auto; background: transparent; border: 0; color: #8a9099; font-size: 18px; line-height: 1;
    cursor: pointer; padding: 0 4px; }
  .ev-att-x:hover { color: #16181c; }
  .ev-foot { display: flex; align-items: center; justify-content: space-between; gap: 10px; padding: 8px 14px 10px;
    border-top: 1px dashed #e8e8e5; }
  .ev-browse { display: inline-flex; align-items: center; gap: 7px; font-size: 11.5px; color: #8a9099; cursor: pointer; }
  .ev-browse u { color: #2463eb; text-decoration: none; font-weight: 500; }
  .ev-browse:hover u { text-decoration: underline; }
  .ev-loading { display: inline-flex; align-items: center; gap: 7px; font-size: 12px; font-weight: 500; color: #2463eb; }
  .ev-spin { width: 13px; height: 13px; border: 1.8px solid #cdddfb; border-top-color: #2463eb; border-radius: 50%;
    animation: ev-spin 0.7s linear infinite; flex-shrink: 0; }
  @keyframes ev-spin { to { transform: rotate(360deg); } }
  .ev-actions { display: flex; align-items: center; justify-content: space-between; gap: 12px; margin-top: 12px; }
  .ev-hint { font-size: 11.5px; color: #8a9099; }
  .ev-hint kbd { font-family: ui-monospace, monospace; font-size: 10.5px; background: #fbfbf9; border: 1px solid #e8e8e5;
    border-bottom-width: 2px; border-radius: 3px; padding: 0 4px; color: #4b5158; }
  .parse-err { color: #b3372a; background: #fdf1ef; border: 1px solid #f2d4cf; border-radius: 8px;
    padding: 8px 12px; font-size: 13px; margin: .75rem 0 0; }
  .ics-preview { margin-top: 18px; padding-top: 16px; border-top: 1px solid #eeeeea; }
  .ics-preview h4 { font-size: 11.5px; font-weight: 600; color: #8a9099; text-transform: uppercase;
    letter-spacing: 0.04em; margin: 0 0 6px; }
  .prev-check { font-size: 12px; color: #8a9099; margin: 0 0 10px; }
  .prev-when strong { font-weight: 600; }
  .prev-row { background: #eef4ff; border: 1px solid #cdddfb; border-radius: 10px; padding: 12px 14px; margin-bottom: 10px; }
  .prev-summary { font-size: 13.5px; font-weight: 600; color: #16181c; }
  .prev-when { font-size: 12.5px; color: #2463eb; margin-top: 3px; font-weight: 500; }
  .prev-loc { font-size: 12px; color: #8a9099; margin-top: 4px; }

  /* Follow-up modal */
  .fu-card { max-width: 460px; }
  .fu-fields { display: flex; flex-direction: column; gap: 14px; }
  .fu-label { display: flex; flex-direction: column; gap: 6px; font-size: 12px; color: #8a9099; }
  .fu-row { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
  .fu-ta, .fu-input { font: inherit; color: #16181c; background: #fbfbf9; border: 1px solid #e8e8e5;
    border-radius: 8px; padding: 9px 11px; font-size: 13.5px; outline: none; box-sizing: border-box; width: 100%; }
  .fu-ta { resize: vertical; line-height: 1.5; }
  .fu-ta:focus, .fu-input:focus { border-color: #2463eb; }

  .empty-tab { border: 1px dashed #e2e2de; border-radius: 12px; padding: 32px; text-align: center; background: #fff; }
  .empty-tab h3 { margin: 0 0 .5rem; font-size: 16px; font-weight: 600; color: #16181c; }
  .empty-tab p { color: #8a9099; margin: 0; font-size: 13.5px; }

  /* MOBILE */
  @media (max-width: 900px) {
    .wrap { padding: 18px 14px 60px; }
    .grid { grid-template-columns: 1fr; gap: 28px; }
    .hd-row h1 { font-size: 23px; }
    .soon { flex-wrap: wrap; }
    .cb-facts { grid-template-columns: repeat(2, 1fr); }
    .cb-cell:nth-child(3) { border-left: none; }
    .cb-cell:nth-child(n+3) { border-top: 1px solid #eeeeea; }
    .modal-overlay { padding: 0; }
    /* Capped scroll box (NOT min-height) — with min-height the modal grows past
       the viewport, the sticky Save footer resolves against the page, and the
       keyboard pushes it off-screen (Ayelet's vanishing-Save bug). */
    .modal { max-width: 100%; border-radius: 0; height: 100dvh; max-height: 100dvh; padding: 1rem;
      padding-bottom: calc(1rem + env(safe-area-inset-bottom)); }
    .fields { grid-template-columns: 1fr; }
    .fields .span-2 { grid-column: auto; }
    .fu-row { grid-template-columns: 1fr; }
  }
</style>
