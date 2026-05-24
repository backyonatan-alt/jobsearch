<script>
  let text = $state('');
  let isDraggingFile = $state(false);
  let attachedImage = $state(null);
  let howToOpen = $state(false);

  function onDragOver(e) {
    if (e.dataTransfer?.types?.includes('Files')) {
      e.preventDefault();
      isDraggingFile = true;
    }
  }
  function onDragLeave() { isDraggingFile = false; }
  function onDrop(e) {
    e.preventDefault();
    isDraggingFile = false;
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

<svelte:head><title>Option D — Pursuit</title></svelte:head>

<header class="prev-bar">
  <span class="badge">Option D</span>
  <span class="title">Zone + expandable how-to</span>
  <a href="/preview/add-modal">← back to options</a>
</header>

<div class="stage">
  <div class="modal" role="dialog" aria-modal="true">
    <header class="m-head">
      <h2>New application</h2>
    </header>

    <div class="paste-block">
      <div class="paste-label">
        <span class="ai-tag">AI</span>
        Paste a job — we'll fill the fields below.
      </div>

      <div
        class={`drop ${isDraggingFile ? 'drag' : ''}`}
        ondragover={onDragOver}
        ondragleave={onDragLeave}
        ondrop={onDrop}
      >
        {#if attachedImage}
          <div class="attached">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="2" y="3" width="12" height="10" rx="1.5"/><circle cx="5.5" cy="6.5" r="1"/><path d="M2 11l3.5-3.5 3 3 2-2L14 11"/></svg>
            <span>{attachedImage.name}</span>
            <span class="size">{Math.round(attachedImage.size / 1024)} KB</span>
            <button class="x" onclick={() => attachedImage = null}>×</button>
          </div>
        {:else}
          <textarea
            bind:value={text}
            onpaste={onPaste}
            placeholder="Drop a screenshot here, or paste text / URL"
            rows="4"
          ></textarea>
        {/if}
      </div>

      <div class="paste-row">
        <button class="btn btn-primary" disabled={!text.trim() && !attachedImage}>Parse</button>
      </div>

      <details class="howto" bind:open={howToOpen}>
        <summary>
          <svg class="caret" width="10" height="10" viewBox="0 0 10 10" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M3 1l4 4-4 4" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          How to add from LinkedIn
        </summary>
        <ol class="steps">
          <li>
            <span class="num">1</span>
            <div>
              <p class="t">Open the job on LinkedIn.</p>
            </div>
          </li>
          <li>
            <span class="num">2</span>
            <div>
              <p class="t">Either:</p>
              <ul class="alts">
                <li><kbd>⌘A</kbd> to select everything on the page, then <kbd>⌘C</kbd> to copy, then <kbd>⌘V</kbd> into the box above</li>
                <li>— or — <kbd>⌘⇧4</kbd> to take a screenshot, then drop the file in (or <kbd>⌘V</kbd> with the image on your clipboard)</li>
              </ul>
            </div>
          </li>
          <li>
            <span class="num">3</span>
            <div>
              <p class="t">Click <strong>Parse</strong>. Review what Claude pulled out, then <strong>Add application</strong>.</p>
            </div>
          </li>
        </ol>
        <p class="howto-note">LinkedIn URLs alone don't work because they block scrapers; selecting the page text or screenshotting gets around that.</p>
      </details>
    </div>

    <div class="modal-divider"><span>or fill in by hand</span></div>

    <div class="fields">
      <label>Company <input placeholder="Anthropic" /></label>
      <label>Role <input placeholder="Senior Software Engineer" /></label>
      <label>Status
        <select><option>Applied</option><option>Wishlist</option><option>Screen</option><option>Interview</option><option>Offer</option></select>
      </label>
      <label>Source <input placeholder="Referral / LinkedIn" /></label>
    </div>

    <div class="m-actions">
      <button class="btn">Cancel</button>
      <button class="btn btn-primary">Add application</button>
    </div>
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

  .paste-block { display: flex; flex-direction: column; gap: 10px; }
  .paste-label { font-size: 13px; color: var(--ink-2); display: flex; align-items: center; gap: 8px; }
  .ai-tag { font-weight: 500; font-size: 10px; color: var(--accent-text); background: var(--accent-tint); border-radius: 4px; padding: 1px 6px; letter-spacing: .04em; }

  .drop { position: relative; border-radius: 8px; outline: 0 dashed transparent; outline-offset: -1px; }
  .drop.drag { outline: 2px dashed var(--accent); background: var(--accent-tint); }
  .drop textarea { width: 100%; font: inherit; font-family: var(--sans); font-size: 13.5px; color: var(--ink); background: var(--surface); border: 1px solid var(--rule); border-radius: 8px; padding: 12px; outline: none; resize: vertical; min-height: 88px; line-height: 1.5; }
  .drop textarea:focus { border-color: var(--accent); }

  .attached { display: flex; align-items: center; gap: 10px; background: var(--accent-tint); border: 1px solid var(--accent-tint-2); border-radius: 8px; padding: 10px 12px; font-size: 13px; }
  .attached svg { color: var(--accent-text); }
  .attached .size { color: var(--mute); font-family: var(--mono); font-size: 11px; }
  .attached .x { margin-left: auto; background: transparent; border: 0; color: var(--mute); font-size: 18px; cursor: pointer; }

  .paste-row { display: flex; gap: 8px; align-items: center; }

  /* Expandable how-to */
  .howto {
    border: 1px solid var(--rule); border-radius: 8px;
    background: var(--surface);
    transition: background 120ms ease;
  }
  .howto[open] { background: var(--card); }
  .howto summary {
    list-style: none;
    cursor: pointer;
    padding: 10px 12px;
    font-size: 13px; font-weight: 500; color: var(--ink-2);
    display: flex; align-items: center; gap: 8px;
  }
  .howto summary::-webkit-details-marker { display: none; }
  .caret { color: var(--mute); transition: transform 120ms ease; }
  .howto[open] .caret { transform: rotate(90deg); color: var(--accent-text); }

  .steps { list-style: none; padding: 0 16px 12px; margin: 0; display: flex; flex-direction: column; gap: 10px; }
  .steps li { display: grid; grid-template-columns: 22px 1fr; gap: 10px; align-items: start; }
  .num { width: 20px; height: 20px; border-radius: 50%; background: var(--ink); color: white; font-family: var(--mono); font-size: 11px; display: grid; place-items: center; }
  .steps .t { margin: 0; font-size: 13px; color: var(--ink); line-height: 1.5; }
  .alts { list-style: none; padding: 0; margin: 6px 0 0; display: flex; flex-direction: column; gap: 6px; }
  .alts li { font-size: 12.5px; color: var(--ink-2); line-height: 1.55; padding-left: 0; }
  kbd { font-family: var(--mono); font-size: 11px; background: var(--surface); border: 1px solid var(--rule); border-bottom-width: 2px; border-radius: 3px; padding: 0 4px; color: var(--ink-2); }
  .howto[open] kbd { background: var(--surface-2); }
  .howto-note { margin: 0 16px 14px; font-size: 12px; color: var(--mute); padding-top: 8px; border-top: 1px dashed var(--rule); }

  .modal-divider { display: flex; align-items: center; gap: 10px; color: var(--mute-2); font-size: 11px; letter-spacing: .04em; text-transform: uppercase; }
  .modal-divider::before, .modal-divider::after { content: ''; flex: 1; height: 1px; background: var(--rule); }

  .fields { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
  .fields label { display: flex; flex-direction: column; gap: 4px; font-size: 12px; color: var(--mute); }
  .fields input, .fields select { font: inherit; font-size: 13.5px; color: var(--ink); background: var(--surface); border: 1px solid var(--rule); border-radius: 6px; padding: 7px 9px; outline: 0; }

  .m-actions { display: flex; justify-content: flex-end; gap: 8px; margin-top: 4px; }
</style>
