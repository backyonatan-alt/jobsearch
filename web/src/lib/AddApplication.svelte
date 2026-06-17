<script>
  // New-application modal — extracted verbatim from the old Today screen so the
  // ⌘N / paste-to-parse / screenshot-vision flow stays intact and reusable.
  // Usage: <AddApplication bind:open onCreated={refresh} />  (open is bindable)
  import { api } from '$lib/api.js';
  import { logEvent } from '$lib/analytics.js';
  import { SOURCE_SUGGESTIONS } from '$lib/app-helpers.js';

  let { open = $bindable(false), onCreated } = $props();

  // Intent signal: fire once each time the modal opens, so we can see who opens
  // the add flow but never completes it (application_create). The drop between
  // the two is the onboarding leak we're trying to make visible.
  let wasOpen = false;
  $effect(() => {
    if (open && !wasOpen) logEvent('addmodal_open');
    wasOpen = open;
  });

  // Whether AI parse was used to fill the form (drives application_create `via`).
  let parseUsed = $state(false);

  // new-application form state
  let form = $state({ company: '', role: '', status: 'applied', source: '', jd_url: '', jd_text: '', cv_variant: '', location: '', salary_note: '' });
  let saving = $state(false);

  // paste-to-parse state
  let pasteText = $state('');
  let parsing = $state(false);
  let parseError = $state('');
  let parsedHint = $state('');
  let attachedImage = $state(null);
  let isDraggingFile = $state(false);
  let parsingStage = $state('');

  const ALLOWED_IMG = ['image/png', 'image/jpeg', 'image/gif', 'image/webp'];
  const MAX_IMG_BYTES = 6 * 1024 * 1024;

  const isBareUrl = (s) => /^https?:\/\/\S+$/.test(s.trim());

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
      parseUsed = true;
      if (r.company)     form.company     = r.company;
      if (r.role)        form.role        = r.role;
      if (r.location)    form.location    = r.location;
      if (r.jd_url)      form.jd_url      = r.jd_url;
      if (r.source)      form.source      = r.source;
      if (r.salary_note) form.salary_note = r.salary_note;
      // Persist the pasted JD body so it survives the posting being taken down.
      // Only a typed-out description counts — a bare URL or a screenshot has no
      // body text worth keeping (the URL goes to jd_url, the image isn't stored).
      if (!attachedImage && text && !isBareUrl(text) && text.length >= 40) {
        form.jd_text = text;
      }
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

  function onModalKeydown(e) { if (e.key === 'Escape') resetModal(); }
  function onModalPaste(e) {
    const item = [...(e.clipboardData?.items || [])].find(i => i.type.startsWith('image/'));
    if (item) setAttachedFile(item.getAsFile());
  }

  function resetModal() {
    form = { company: '', role: '', status: 'applied', source: '', jd_url: '', jd_text: '', cv_variant: '', location: '', salary_note: '' };
    pasteText = '';
    parseError = '';
    parsedHint = '';
    attachedImage = null;
    isDraggingFile = false;
    parseUsed = false;
    open = false;
  }

  async function createApp(e) {
    e.preventDefault();
    saving = true;
    try {
      const payload = { ...form };
      for (const k of ['jd_url', 'jd_text', 'cv_variant', 'source', 'location', 'salary_note']) if (!payload[k]) delete payload[k];
      await api('/api/applications', { method: 'POST', body: JSON.stringify(payload) });
      // Confirmed-success only (after the await), before resetModal clears state.
      logEvent('application_create', {
        via: parseUsed ? 'parse' : 'manual',
        source: form.source || 'none',
        status: form.status
      });
      resetModal();
      onCreated?.();
    } finally {
      saving = false;
    }
  }
</script>

<svelte:window
  onkeydown={(e) => { if (open) onModalKeydown(e); }}
  onpaste={(e) => { if (open) onModalPaste(e); }}
/>

{#if open}
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
          <span class="ai-tag-sm">AI</span>
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
          <input bind:value={form.source} list="source-suggestions" placeholder="LinkedIn / Referral / Cold email" />
          <datalist id="source-suggestions">
            {#each SOURCE_SUGGESTIONS as s}<option value={s}></option>{/each}
          </datalist>
        </label>
        <label class="span-2">
          <span class="lbl">Job description <span class="opt">— optional, kept even if the posting comes down</span></span>
          <textarea class="jd-area" bind:value={form.jd_text} rows="3" placeholder="Paste the full JD text here so it's saved with the application."></textarea>
        </label>
        <p class="privacy-note span-2">
          <svg width="12" height="12" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5" aria-hidden="true"><rect x="3" y="7" width="10" height="6.5" rx="1.5"/><path d="M5.5 7V5a2.5 2.5 0 0 1 5 0v2"/></svg>
          Private to your account — never shared with companies or anyone else.
        </p>
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
  .m-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 16px; padding: 20px 22px 14px; }
  .m-head h2 { font-size: 19px; font-weight: 500; letter-spacing: -0.018em; margin: 0; color: var(--ink); }
  .m-head .sub { margin: 4px 0 0; font-size: 13px; color: var(--mute); line-height: 1.5; max-width: 460px; }
  .x-close {
    background: transparent; border: 0; color: var(--mute);
    width: 28px; height: 28px; border-radius: 6px;
    display: grid; place-items: center; cursor: pointer;
    flex-shrink: 0;
    transition: background 100ms ease, color 100ms ease;
  }
  .x-close:hover { background: var(--surface-2); color: var(--ink); }
  .quick-add { padding: 0 22px 16px; display: flex; flex-direction: column; gap: 10px; }
  .qa-head { display: flex; align-items: center; justify-content: space-between; }
  .qa-label { font-size: 12px; font-weight: 500; color: var(--ink-2); letter-spacing: 0; }
  .ai-tag-sm { font-weight: 500; font-size: 10px; color: var(--accent-text); background: var(--accent-tint); border-radius: 4px; padding: 2px 7px; letter-spacing: .04em; }

  .drop { position: relative; background: var(--surface); border: 1.5px dashed var(--rule-strong); border-radius: 10px; transition: border-color 120ms ease, background 120ms ease; overflow: hidden; }
  .drop.drag { border-color: var(--accent); background: var(--accent-tint); }
  .drop.loading { border-color: var(--accent); border-style: solid; }
  .drop.loading::before {
    content: ''; position: absolute; inset: 0 0 auto 0; height: 2px;
    background: linear-gradient(90deg, transparent, var(--accent), transparent);
    background-size: 40% 100%; background-repeat: no-repeat;
    animation: drop-sweep 1.4s linear infinite;
  }
  @keyframes drop-sweep {
    0%   { background-position: -40% 0; }
    100% { background-position: 140% 0; }
  }
  .drop-loading { display: flex; align-items: center; gap: 10px; padding: 10px 14px; background: var(--accent-tint); border-top: 1px solid var(--accent-tint-2); font-size: 12.5px; color: var(--accent-text); font-weight: 500; }
  .loading-text { letter-spacing: 0; }
  .spinner { display: inline-block; width: 14px; height: 14px; border: 1.8px solid var(--accent-tint-2); border-top-color: var(--accent); border-radius: 50%; animation: spinner-rot 0.7s linear infinite; flex-shrink: 0; }
  .spinner-sm { width: 12px; height: 12px; border-color: rgba(255,255,255,.35); border-top-color: white; }
  @keyframes spinner-rot { to { transform: rotate(360deg); } }
  .ok-check { display: inline-block; width: 14px; height: 14px; line-height: 14px; text-align: center; border-radius: 50%; background: var(--positive-tint); color: var(--positive-text); font-size: 10px; font-weight: 600; margin-right: 4px; }
  .drop textarea { width: 100%; font: inherit; font-family: var(--sans); font-size: 13.5px; color: var(--ink); background: transparent; border: 0; padding: 12px 14px 8px; outline: none; resize: none; min-height: 78px; line-height: 1.5; display: block; }
  .drop textarea::placeholder { color: var(--mute-2); }
  .attached { display: flex; align-items: center; gap: 10px; padding: 12px 14px 8px; font-size: 13px; color: var(--ink-2); }
  .attached svg { color: var(--accent-text); flex-shrink: 0; }
  .attached .att-name { font-weight: 500; color: var(--ink); }
  .attached .att-size { color: var(--mute); font-family: var(--mono); font-size: 11px; }
  .attached .att-x { margin-left: auto; background: transparent; border: 0; color: var(--mute); font-size: 18px; line-height: 1; cursor: pointer; padding: 0 4px; }
  .attached .att-x:hover { color: var(--ink); }
  .drop-footer { display: flex; align-items: center; gap: 8px; padding: 8px 14px 10px; font-size: 11.5px; color: var(--mute-2); border-top: 1px dashed var(--rule); background: transparent; }
  .drop-footer svg { color: var(--mute-2); flex-shrink: 0; }
  .qa-actions { display: flex; align-items: center; justify-content: space-between; gap: 12px; }
  .kb-hints { font-size: 11.5px; color: var(--mute); }
  .kb-hints kbd { font-family: var(--mono); font-size: 10.5px; background: var(--surface); border: 1px solid var(--rule); border-bottom-width: 2px; border-radius: 3px; padding: 0 4px; color: var(--ink-2); }
  .btn-accent { display: inline-flex; align-items: center; gap: 6px; background: var(--accent); color: white; border: 1px solid var(--accent); border-radius: 7px; padding: 7px 12px; font-size: 13px; font-weight: 500; cursor: pointer; transition: background 100ms ease, border-color 100ms ease; }
  .btn-accent:hover:not(:disabled) { background: var(--accent-strong); border-color: var(--accent-strong); }
  .btn-accent:disabled { opacity: .55; cursor: not-allowed; }
  .btn-accent svg { color: white; }
  .qa-msg { font-size: 12px; }
  .qa-msg.ok  { color: var(--positive-text); }
  .qa-msg.err { color: var(--danger-text); }
  .modal-divider { display: flex; align-items: center; gap: 12px; color: var(--mute-2); font-size: 10.5px; letter-spacing: .06em; text-transform: uppercase; padding: 0 22px; margin: 4px 0 14px; }
  .modal-divider::before, .modal-divider::after { content: ''; flex: 1; height: 1px; background: var(--rule); }
  .fields { padding: 0 22px 18px; display: grid; grid-template-columns: 1fr 1fr; gap: 12px 14px; }
  .fields label { display: flex; flex-direction: column; gap: 5px; }
  .fields .span-2 { grid-column: 1 / -1; }
  .privacy-note { display: flex; align-items: center; gap: 7px; font-size: 11.5px; color: var(--mute); margin: 2px 0 0; line-height: 1.4; }
  .privacy-note svg { color: var(--mute-2); flex-shrink: 0; }
  .fields .jd-area { font: inherit; font-family: var(--sans); font-size: 13.5px; color: var(--ink); background: var(--surface); border: 1px solid var(--rule); border-radius: 8px; padding: 9px 11px; outline: 0; resize: vertical; min-height: 64px; line-height: 1.5; transition: border-color 100ms ease, box-shadow 100ms ease; }
  .fields .jd-area:hover { border-color: var(--rule-strong); }
  .fields .jd-area:focus { border-color: var(--accent); box-shadow: 0 0 0 3px var(--accent-tint); }
  .fields .jd-area::placeholder { color: var(--mute-2); }
  .fields .lbl { font-size: 11.5px; color: var(--mute); font-weight: 500; letter-spacing: 0; }
  .fields .req { color: var(--accent-text); font-style: normal; margin-left: 1px; }
  .fields .opt { color: var(--mute-2); font-weight: 400; }
  .fields input { font: inherit; font-size: 13.5px; color: var(--ink); background: var(--surface); border: 1px solid var(--rule); border-radius: 8px; padding: 0 11px; height: 36px; outline: 0; transition: border-color 100ms ease, box-shadow 100ms ease, background 100ms ease; }
  .fields input:hover { border-color: var(--rule-strong); }
  .fields input:focus { border-color: var(--accent); box-shadow: 0 0 0 3px var(--accent-tint); }
  .fields input::placeholder { color: var(--mute-2); }
  .status-select { position: relative; display: flex; align-items: center; background: var(--surface); border: 1px solid var(--rule); border-radius: 8px; height: 36px; padding: 0 11px; transition: border-color 100ms ease, box-shadow 100ms ease; }
  .status-select:hover { border-color: var(--rule-strong); }
  .status-select:focus-within { border-color: var(--accent); box-shadow: 0 0 0 3px var(--accent-tint); }
  .status-select .sdot { width: 7px; height: 7px; border-radius: 50%; background: var(--mute); margin-right: 8px; flex-shrink: 0; }
  .status-select.wishlist  .sdot { background: var(--mute); }
  .status-select.applied   .sdot { background: var(--ink-2); }
  .status-select.screen    .sdot { background: var(--positive); }
  .status-select.interview .sdot { background: var(--accent); }
  .status-select.offer     .sdot { background: var(--warm); }
  .status-select.rejected  .sdot { background: var(--mute-2); }
  .status-select.withdrawn .sdot { background: var(--mute-2); }
  .status-select select { flex: 1; appearance: none; -webkit-appearance: none; background: transparent; border: 0; outline: 0; font: inherit; font-size: 13.5px; color: var(--ink); padding: 0 18px 0 0; cursor: pointer; }
  .status-select .chev { position: absolute; right: 11px; top: 50%; transform: translateY(-50%); color: var(--mute); pointer-events: none; }
  .m-foot { display: flex; align-items: center; justify-content: space-between; padding: 12px 22px; background: var(--surface); border-top: 1px solid var(--rule); }
  .esc-hint { font-size: 11.5px; color: var(--mute); }
  .esc-hint kbd { font-family: var(--mono); font-size: 10.5px; background: var(--card); border: 1px solid var(--rule); border-bottom-width: 2px; border-radius: 3px; padding: 0 5px; color: var(--ink-2); }
  .foot-actions { display: flex; gap: 8px; }
  .foot-actions .btn-primary { display: inline-flex; align-items: center; gap: 6px; }
  .dark-kbd { font-family: var(--mono); font-size: 10.5px; background: rgba(255,255,255,.18); border: 1px solid rgba(255,255,255,.22); border-radius: 3px; padding: 0 5px; color: rgba(255,255,255,.9); }

  /* Buttons reused from the shell. */
  .btn { background: var(--card); border: 1px solid var(--rule); border-radius: 7px; padding: 6px 11px; font-size: 13px; font-weight: 500; color: var(--ink); cursor: pointer; }
  .btn:hover { border-color: var(--rule-strong); }
  .btn-primary { background: var(--accent); border-color: var(--accent-strong); color: white; }
  .btn-primary:hover { background: var(--accent-strong); }
  .btn-primary:disabled { opacity: .55; cursor: not-allowed; }

  @media (max-width: 720px) {
    .modal { max-width: 100%; }
    .fields { grid-template-columns: 1fr; }
  }
</style>
