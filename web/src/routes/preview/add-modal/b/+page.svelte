<script>
  let mode = $state('paste'); // 'paste' | 'screenshot' | 'manual'
  let text = $state('');
  let attachedImage = $state(null);

  function onDrop(e) {
    e.preventDefault();
    const f = e.dataTransfer?.files?.[0];
    if (f && f.type.startsWith('image/')) attachedImage = { name: f.name, size: f.size };
  }
  function onPaste(e) {
    const item = [...(e.clipboardData?.items || [])].find(i => i.type.startsWith('image/'));
    if (item) {
      const f = item.getAsFile();
      if (f) attachedImage = { name: f.name || 'pasted-image.png', size: f.size };
    }
  }
</script>

<svelte:head><title>Option B — Pursuit</title></svelte:head>

<header class="prev-bar">
  <span class="badge">Option B</span>
  <span class="title">Three explicit chips</span>
  <a href="/preview/add-modal">← back to options</a>
</header>

<div class="stage">
  <div class="modal" role="dialog" aria-modal="true">
    <header class="m-head">
      <h2>New application</h2>
    </header>

    <div class="mode-row">
      <span class="mode-label">How?</span>
      <button class={`chip ${mode === 'paste' ? 'active' : ''}`} onclick={() => mode = 'paste'}>
        <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="3" y="2" width="10" height="12" rx="1.5"/><path d="M6 4h4M6 7h4M6 10h4"/></svg>
        Paste text
      </button>
      <button class={`chip ${mode === 'screenshot' ? 'active' : ''}`} onclick={() => mode = 'screenshot'}>
        <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="2" y="3" width="12" height="10" rx="1.5"/><circle cx="5.5" cy="6.5" r="1"/><path d="M2 11l3.5-3.5 3 3 2-2L14 11"/></svg>
        Screenshot
      </button>
      <button class={`chip ${mode === 'manual' ? 'active' : ''}`} onclick={() => mode = 'manual'}>
        <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M3 12l8-8 2 2-8 8H3v-2z"/></svg>
        Type it
      </button>
    </div>

    {#if mode === 'paste'}
      <p class="help">
        Copy a job page, recruiter email, or your own note and paste it. From LinkedIn: open the job → <kbd>⌘A</kbd> → <kbd>⌘V</kbd> here.
      </p>
      <textarea bind:value={text} onpaste={onPaste} placeholder="Paste anything — JD, URL, recruiter email, a short note" rows="6"></textarea>
      <div class="m-actions">
        <button class="btn">Cancel</button>
        <button class="btn btn-primary" disabled={!text.trim()}>Parse</button>
      </div>
    {:else if mode === 'screenshot'}
      <p class="help">
        Take a screenshot of the job page (<kbd>⌘⇧4</kbd> on Mac, <kbd>Win+Shift+S</kbd> on Windows), then paste it below or drag the file in.
      </p>
      <div class="drop big" ondragover={(e) => e.preventDefault()} ondrop={onDrop}>
        {#if attachedImage}
          <div class="attached">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="2" y="3" width="12" height="10" rx="1.5"/><circle cx="5.5" cy="6.5" r="1"/><path d="M2 11l3.5-3.5 3 3 2-2L14 11"/></svg>
            <span>{attachedImage.name}</span>
            <span class="size">{Math.round(attachedImage.size / 1024)} KB</span>
            <button class="x" onclick={() => attachedImage = null}>×</button>
          </div>
        {:else}
          <div class="drop-inner">
            <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.4" style="opacity:.45"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4M17 8l-5-5-5 5M12 3v12"/></svg>
            <p>Drop screenshot here, or paste with <kbd>⌘V</kbd></p>
          </div>
        {/if}
      </div>
      <div class="m-actions">
        <button class="btn">Cancel</button>
        <button class="btn btn-primary" disabled={!attachedImage}>Parse screenshot</button>
      </div>
    {:else}
      <p class="help">Fill in what you know. Status defaults to Applied.</p>
      <div class="fields">
        <label>Company <input placeholder="Anthropic" /></label>
        <label>Role <input placeholder="Senior Software Engineer" /></label>
        <label>Status
          <select><option>Applied</option><option>Wishlist</option><option>Screen</option><option>Interview</option><option>Offer</option></select>
        </label>
        <label>Source <input placeholder="Referral / LinkedIn" /></label>
        <label>Location <input placeholder="Remote / SF" /></label>
        <label>CV variant <input placeholder="v3-ai-focus" /></label>
        <label class="span-2">JD URL <input placeholder="https://…" /></label>
      </div>
      <div class="m-actions">
        <button class="btn">Cancel</button>
        <button class="btn btn-primary">Add application</button>
      </div>
    {/if}
  </div>
</div>

<style>
  :global(html, body) { background: var(--surface); }
  .prev-bar { position: sticky; top: 0; z-index: 50; display: flex; align-items: center; gap: 12px; padding: 10px 16px; background: var(--ink); color: white; font-size: 12px; }
  .prev-bar .badge { font-family: var(--mono); background: rgba(255,255,255,.15); border-radius: 4px; padding: 2px 7px; letter-spacing: .04em; }
  .prev-bar .title { font-weight: 500; }
  .prev-bar a { margin-left: auto; color: rgba(255,255,255,.7); text-decoration: none; }

  .stage { min-height: calc(100vh - 36px); background: rgba(10,10,13,.4); display: grid; place-items: start center; padding: 64px 24px; }
  .modal { width: 100%; max-width: 560px; background: var(--card); border: 1px solid var(--rule); border-radius: 12px; padding: 24px; box-shadow: var(--sh-pop); display: flex; flex-direction: column; gap: 16px; }
  .m-head h2 { font-size: 18px; font-weight: 500; letter-spacing: -0.018em; margin: 0; }

  .mode-row { display: flex; align-items: center; gap: 6px; flex-wrap: wrap; }
  .mode-label { font-size: 12px; color: var(--mute); margin-right: 4px; }
  .chip {
    display: inline-flex; align-items: center; gap: 6px;
    font-size: 13px; font-weight: 500;
    color: var(--ink-2); background: var(--card);
    border: 1px solid var(--rule); border-radius: 6px;
    padding: 5px 10px; cursor: pointer;
  }
  .chip:hover { background: var(--surface-2); }
  .chip.active { background: var(--ink); color: white; border-color: var(--ink); }

  .help { margin: 0; font-size: 13px; color: var(--mute); line-height: 1.55; }
  kbd { font-family: var(--mono); font-size: 11px; background: var(--surface); border: 1px solid var(--rule); border-bottom-width: 2px; border-radius: 3px; padding: 0 4px; color: var(--ink-2); }

  textarea { width: 100%; font: inherit; font-family: var(--sans); font-size: 13.5px; color: var(--ink); background: var(--surface); border: 1px solid var(--rule); border-radius: 8px; padding: 12px; outline: none; resize: vertical; min-height: 120px; line-height: 1.5; }
  textarea:focus { border-color: var(--accent); }

  .drop.big { border: 2px dashed var(--rule-strong); border-radius: 10px; padding: 32px; background: var(--surface); }
  .drop-inner { display: flex; flex-direction: column; align-items: center; gap: 12px; color: var(--mute); }
  .drop-inner p { margin: 0; font-size: 13.5px; }

  .attached { display: flex; align-items: center; gap: 10px; background: var(--accent-tint); border: 1px solid var(--accent-tint-2); border-radius: 8px; padding: 10px 12px; font-size: 13px; }
  .attached svg { color: var(--accent-text); }
  .attached .size { color: var(--mute); font-family: var(--mono); font-size: 11px; }
  .attached .x { margin-left: auto; background: transparent; border: 0; color: var(--mute); font-size: 18px; cursor: pointer; }

  .fields { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
  .fields .span-2 { grid-column: span 2; }
  .fields label { display: flex; flex-direction: column; gap: 4px; font-size: 12px; color: var(--mute); }
  .fields input, .fields select { font: inherit; font-size: 13.5px; color: var(--ink); background: var(--surface); border: 1px solid var(--rule); border-radius: 6px; padding: 7px 9px; outline: 0; }

  .m-actions { display: flex; justify-content: flex-end; gap: 8px; margin-top: 4px; }
</style>
