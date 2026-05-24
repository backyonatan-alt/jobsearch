<script>
  let chosen = $state(null); // null | 'paste' | 'screenshot' | 'manual'
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

<svelte:head><title>Option C — Pursuit</title></svelte:head>

<header class="prev-bar">
  <span class="badge">Option C</span>
  <span class="title">Three visual cards</span>
  <a href="/preview/add-modal">← back to options</a>
</header>

<div class="stage">
  <div class="modal" role="dialog" aria-modal="true">
    <header class="m-head">
      <h2>Add a job</h2>
      <p class="sub">Pick how you want to give us the data.</p>
    </header>

    <div class="picker">
      <button class={`pick ${chosen === 'paste' ? 'active' : ''}`} onclick={() => chosen = 'paste'}>
        <div class="ic">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="5" y="3" width="14" height="18" rx="2"/><path d="M9 7h6M9 11h6M9 15h4"/></svg>
        </div>
        <h3>Paste text</h3>
        <p>From a job page, recruiter email, or your own note.<br/><strong>LinkedIn:</strong> ⌘A then ⌘V.</p>
      </button>

      <button class={`pick ${chosen === 'screenshot' ? 'active' : ''}`} onclick={() => chosen = 'screenshot'}>
        <div class="ic">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="3" y="5" width="18" height="14" rx="2"/><circle cx="8.5" cy="10" r="1.5"/><path d="M3 17l5-5 4 4 3-3 6 6"/></svg>
        </div>
        <h3>Screenshot</h3>
        <p>Take a picture of the job page (<strong>⌘⇧4</strong>), drop or paste it here.</p>
      </button>

      <button class={`pick ${chosen === 'manual' ? 'active' : ''}`} onclick={() => chosen = 'manual'}>
        <div class="ic">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M4 17l9-9 3 3-9 9H4v-3z"/><path d="M14 6l3 3"/></svg>
        </div>
        <h3>Type fields</h3>
        <p>Fill in company, role, status, source yourself.</p>
      </button>
    </div>

    {#if chosen === 'paste'}
      <div class="expander">
        <textarea bind:value={text} onpaste={onPaste} placeholder="Paste the JD body, URL, or your note. From LinkedIn: ⌘A on the job page, then ⌘V here." rows="6"></textarea>
        <div class="m-actions">
          <button class="btn">Cancel</button>
          <button class="btn btn-primary" disabled={!text.trim()}>Parse</button>
        </div>
      </div>
    {:else if chosen === 'screenshot'}
      <div class="expander">
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
              <svg width="36" height="36" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.3" style="opacity:.45"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4M17 8l-5-5-5 5M12 3v12"/></svg>
              <p>Drop screenshot, or paste with <kbd>⌘V</kbd></p>
            </div>
          {/if}
        </div>
        <div class="m-actions">
          <button class="btn">Cancel</button>
          <button class="btn btn-primary" disabled={!attachedImage}>Parse screenshot</button>
        </div>
      </div>
    {:else if chosen === 'manual'}
      <div class="expander">
        <div class="fields">
          <label>Company <input placeholder="Anthropic" /></label>
          <label>Role <input placeholder="Senior Software Engineer" /></label>
          <label>Status
            <select><option>Applied</option><option>Wishlist</option><option>Screen</option><option>Interview</option><option>Offer</option></select>
          </label>
          <label>Source <input placeholder="Referral" /></label>
        </div>
        <div class="m-actions">
          <button class="btn">Cancel</button>
          <button class="btn btn-primary">Add application</button>
        </div>
      </div>
    {:else}
      <div class="m-actions">
        <button class="btn">Cancel</button>
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
  .modal { width: 100%; max-width: 640px; background: var(--card); border: 1px solid var(--rule); border-radius: 12px; padding: 24px; box-shadow: var(--sh-pop); display: flex; flex-direction: column; gap: 18px; }
  .m-head h2 { font-size: 19px; font-weight: 500; letter-spacing: -0.018em; margin: 0; }
  .m-head .sub { color: var(--mute); font-size: 13px; margin: 2px 0 0; }

  .picker { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; }
  .pick {
    display: flex; flex-direction: column; align-items: flex-start;
    gap: 8px; padding: 16px; text-align: left;
    background: var(--card); border: 1px solid var(--rule); border-radius: 10px;
    cursor: pointer; transition: border-color 120ms ease, transform 120ms ease;
    color: var(--ink-2); font: inherit;
  }
  .pick:hover { border-color: var(--rule-strong); transform: translateY(-1px); }
  .pick.active { border-color: var(--accent); background: var(--accent-tint); }
  .ic {
    width: 36px; height: 36px; border-radius: 8px;
    background: var(--surface-2); color: var(--accent-text);
    display: grid; place-items: center;
  }
  .pick.active .ic { background: white; }
  .pick h3 { font-size: 14px; font-weight: 500; margin: 0; color: var(--ink); }
  .pick p { font-size: 12px; color: var(--mute); margin: 0; line-height: 1.5; }
  .pick p strong { color: var(--ink-2); font-weight: 500; }

  .expander { display: flex; flex-direction: column; gap: 12px; }

  textarea { width: 100%; font: inherit; font-family: var(--sans); font-size: 13.5px; color: var(--ink); background: var(--surface); border: 1px solid var(--rule); border-radius: 8px; padding: 12px; outline: none; resize: vertical; min-height: 130px; line-height: 1.5; }
  textarea:focus { border-color: var(--accent); }

  .drop.big { border: 2px dashed var(--rule-strong); border-radius: 10px; padding: 36px; background: var(--surface); }
  .drop-inner { display: flex; flex-direction: column; align-items: center; gap: 12px; color: var(--mute); }
  .drop-inner p { margin: 0; font-size: 13.5px; }
  .drop-inner kbd { font-family: var(--mono); font-size: 11px; background: var(--card); border: 1px solid var(--rule); border-bottom-width: 2px; border-radius: 3px; padding: 0 4px; color: var(--ink-2); }

  .attached { display: flex; align-items: center; gap: 10px; background: var(--accent-tint); border: 1px solid var(--accent-tint-2); border-radius: 8px; padding: 10px 12px; font-size: 13px; }
  .attached svg { color: var(--accent-text); }
  .attached .size { color: var(--mute); font-family: var(--mono); font-size: 11px; }
  .attached .x { margin-left: auto; background: transparent; border: 0; color: var(--mute); font-size: 18px; cursor: pointer; }

  .fields { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
  .fields label { display: flex; flex-direction: column; gap: 4px; font-size: 12px; color: var(--mute); }
  .fields input, .fields select { font: inherit; font-size: 13.5px; color: var(--ink); background: var(--surface); border: 1px solid var(--rule); border-radius: 6px; padding: 7px 9px; outline: 0; }

  .m-actions { display: flex; justify-content: flex-end; gap: 8px; }
</style>
