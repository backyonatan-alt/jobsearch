<script>
  import { onMount } from 'svelte';
  import { api } from '$lib/api.js';

  let me = $state(null);
  let cvText = $state('');
  let savedChars = $state(0);
  let cvLoading = $state(true);
  let saving = $state(false);
  let savedFlash = $state(false);
  let extracting = $state(false);
  let cvError = $state('');
  let reminders = $state(true);
  let remindersSaving = $state(false);

  onMount(async () => {
    try {
      me = await api('/api/me');
      reminders = me.email_reminders !== false;
      const cv = await api('/api/me/cv');
      cvText = cv.cv_text || '';
      savedChars = cvText.length;
    } catch (e) { if (e.message !== 'unauthorized') console.error(e); }
    cvLoading = false;
  });

  async function saveCV() {
    if (saving) return;
    saving = true; cvError = '';
    try {
      const r = await api('/api/me/cv', { method: 'PUT', body: JSON.stringify({ cv_text: cvText }) });
      savedChars = r.cv_chars;
      savedFlash = true; setTimeout(() => (savedFlash = false), 2500);
    } catch (e) { cvError = e.message; }
    saving = false;
  }

  function fileToBase64(file) {
    return new Promise((res, rej) => {
      const rd = new FileReader();
      rd.onload = () => res(String(rd.result).split(',')[1] ?? '');
      rd.onerror = rej;
      rd.readAsDataURL(file);
    });
  }

  async function onFile(e) {
    const file = e.target.files?.[0];
    e.target.value = '';
    if (!file) return;
    cvError = '';
    if (file.type === 'text/plain' || /\.txt$/i.test(file.name)) {
      cvText = await file.text();
      return;
    }
    const ok = file.type === 'application/pdf' || file.type.startsWith('image/');
    if (!ok) { cvError = 'PDF, image, or .txt — those are the formats I can read.'; return; }
    if (file.size > 8 * 1024 * 1024) { cvError = 'That file is over 8 MB — export a smaller PDF.'; return; }
    extracting = true;
    try {
      const data = await fileToBase64(file);
      const r = await api('/api/me/cv/extract', {
        method: 'POST',
        body: JSON.stringify({ media_type: file.type || 'application/pdf', data })
      });
      cvText = r.cv_text;
    } catch (err) { cvError = err.message; }
    extracting = false;
  }

  async function toggleReminders() {
    if (remindersSaving) return;
    remindersSaving = true;
    const next = !reminders;
    try {
      await api('/api/me/reminders', { method: 'PUT', body: JSON.stringify({ enabled: next }) });
      reminders = next;
    } catch (e) { if (e.message !== 'unauthorized') console.error(e); }
    remindersSaving = false;
  }
</script>

<svelte:head><title>Profile — Pursuit</title></svelte:head>

<div class="page">
  <header class="hd">
    <h1>Profile</h1>
    {#if me}<p class="sub">{me.email}</p>{/if}
  </header>

  <section class="card" data-testid="cv-card">
    <div class="card-hd">
      <h2>Your CV</h2>
      {#if savedChars > 0}<span class="chip on">On file — powering your prep</span>
      {:else}<span class="chip">Not added yet</span>{/if}
    </div>
    <p class="lead">Paste it or upload a PDF. Every round brief you generate gains a
      <strong>"Make sure they hear"</strong> section — the things from <em>your</em> background
      to land in that specific interview.</p>

    {#if cvLoading}
      <p class="muted">Loading…</p>
    {:else}
      <div class="upload-row">
        <label class="btn upload">
          {extracting ? 'Reading your CV…' : 'Upload PDF / image / .txt'}
          <input type="file" accept=".pdf,.txt,image/*,application/pdf,text/plain" onchange={onFile} disabled={extracting} />
        </label>
        <span class="or">or paste below</span>
      </div>
      <textarea class="cv-box" rows="14" placeholder="Paste your CV text here…" bind:value={cvText} disabled={extracting}></textarea>
      {#if cvError}<p class="err">{cvError}</p>{/if}
      <div class="actions">
        <span class="privacy">Private to your account — used only to personalize your prep, never shared.</span>
        <button class="btn btn-primary" onclick={saveCV} disabled={saving || extracting}>
          {saving ? 'Saving…' : savedFlash ? 'Saved ✓' : 'Save CV'}
        </button>
      </div>
    {/if}
  </section>

  <section class="card">
    <div class="card-hd"><h2>Interview reminders</h2></div>
    <label class="toggle-row">
      <input type="checkbox" checked={reminders} onchange={toggleReminders} disabled={remindersSaving} />
      <span>Email me the playbook the day before a scheduled round</span>
    </label>
    <p class="muted small">One email per round, only when a date is set. Every email has a one-click stop link too.</p>
  </section>
</div>

<style>
  .page { max-width: 720px; margin: 0 auto; padding: 26px 20px 60px; }
  .hd h1 { margin: 0; font-size: 22px; font-weight: 700; color: #16181c; }
  .sub { margin: 3px 0 0; font-size: 13px; color: #8a9099; }
  .card { background: #fff; border: 1px solid #e8e8e5; border-radius: 14px; padding: 18px 20px; margin-top: 18px; }
  .card-hd { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
  .card-hd h2 { margin: 0; font-size: 15.5px; font-weight: 650; color: #16181c; }
  .chip { font-size: 11.5px; color: #6b7280; background: #f4f4f2; border-radius: 999px; padding: 3px 9px; }
  .chip.on { color: #166534; background: #ecfdf3; }
  .lead { font-size: 13px; color: #4b5158; margin: 8px 0 14px; }
  .upload-row { display: flex; align-items: center; gap: 10px; margin-bottom: 10px; }
  .btn { font-family: inherit; font-size: 13px; border: 1px solid #d9d9d5; background: #fff; color: #16181c;
    border-radius: 9px; padding: 7px 13px; cursor: pointer; }
  .btn:hover { border-color: #b9b9b4; }
  .upload { display: inline-block; }
  .upload input { display: none; }
  .or { font-size: 12px; color: #8a9099; }
  .cv-box { width: 100%; box-sizing: border-box; font-family: inherit; font-size: 12.5px; line-height: 1.5;
    color: #16181c; border: 1px solid #e0e0dc; border-radius: 10px; padding: 10px 12px; resize: vertical; background: #fbfbf9; }
  .cv-box:focus { outline: none; border-color: #2463eb; background: #fff; }
  .err { font-size: 12.5px; color: #b42318; margin: 8px 0 0; }
  .actions { display: flex; align-items: center; justify-content: space-between; gap: 14px; margin-top: 12px; flex-wrap: wrap; }
  .privacy { font-size: 11.5px; color: #8a9099; flex: 1; min-width: 200px; }
  .btn-primary { background: #2463eb; border-color: #2463eb; color: #fff; }
  .btn-primary:hover { background: #1d53c9; }
  .btn-primary:disabled { opacity: 0.6; cursor: default; }
  .toggle-row { display: flex; align-items: center; gap: 10px; font-size: 13.5px; color: #16181c; margin-top: 10px; cursor: pointer; }
  .muted { color: #8a9099; font-size: 13px; }
  .muted.small { font-size: 12px; margin: 8px 0 0; }
</style>
