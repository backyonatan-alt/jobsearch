<script>
  // Interviews v2 — side-by-side calendar / screenshot+text drop zones.
  // Static mockup. Real flow: same /api/applications/{id}/interviews/parse
  // endpoint, extended to also accept {image} and {text} payloads.
  let parsed = $state([]);
  function fakeParse() {
    parsed = [
      { id: 1, summary: 'Stripe — Technical screen', when: 'Tue, May 28 · 2:00 PM (60 min)', location: 'Google Meet' }
    ];
  }
  function clearPreview() { parsed = []; }
</script>

<svelte:head><title>Interviews v2 — Pursuit</title></svelte:head>

<main class="wrap">
  <div class="topbar-stub">
    <div class="crumb"><span class="root">Applications</span><span class="sep">/</span><span class="here">Stripe</span></div>
  </div>

  <div class="body-inner">
    <!-- TABS stub -->
    <div class="tabs">
      <button class="tab">Brief <span class="t-tag">AI</span></button>
      <button class="tab active">Interviews <span class="t-tag">0</span></button>
      <button class="tab">Timeline <span class="t-tag">3</span></button>
      <button class="tab">Notes</button>
      <button class="tab">Files</button>
    </div>

    <!-- NEW: two side-by-side drop zones -->
    <div class="add-card">
      <div class="add-hd">
        <h3>Add an interview</h3>
        <p>Drop a calendar file, paste a screenshot, or just paste the email body — we'll extract the event.</p>
      </div>

      <div class="zones">
        <!-- LEFT zone — .ics file or text -->
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
          <div class="drop drop-file">
            <svg width="20" height="20" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M8 11V3M5 6l3-3 3 3M3 11v2h10v-2"/></svg>
            <span class="drop-l1">Drop .ics file</span>
            <span class="drop-l2">or click to browse</span>
          </div>
          <div class="or">or paste raw .ics text below</div>
          <textarea rows="3" placeholder={"BEGIN:VCALENDAR&#10;VERSION:2.0&#10;BEGIN:VEVENT…"}></textarea>
        </div>

        <!-- RIGHT zone — screenshot or email body -->
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
          <div class="drop drop-image">
            <svg width="20" height="20" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="2" y="3" width="12" height="10" rx="1.5"/><circle cx="5.5" cy="6.5" r="1"/><path d="M2 11l3.5-3.5 3 3 2-2L14 11"/></svg>
            <span class="drop-l1">Drop a screenshot</span>
            <span class="drop-l2"><kbd>⌘V</kbd> works too</span>
          </div>
          <div class="or">or paste the email body below</div>
          <textarea rows="3" placeholder={"You're invited to: Stripe — Technical screen&#10;When: Tue, May 28, 2:00 PM EDT&#10;Where: Google Meet"}></textarea>
        </div>
      </div>

      <div class="card-actions">
        <button class="btn primary" onclick={fakeParse}>Parse</button>
        {#if parsed.length > 0}
          <button class="btn ghost" onclick={clearPreview}>Clear preview</button>
        {/if}
      </div>

      {#if parsed.length > 0}
        <div class="preview-wrap">
          <h4>Preview</h4>
          {#each parsed as ev}
            <div class="prev-row">
              <div>
                <div class="prev-summary">{ev.summary}</div>
                <div class="prev-when">{ev.when}</div>
                <div class="prev-loc">📍 {ev.location}</div>
              </div>
            </div>
          {/each}
          <button class="btn primary save">Save {parsed.length} event{parsed.length === 1 ? '' : 's'}</button>
        </div>
      {/if}
    </div>

    <div class="scheduled">
      <h3>Scheduled</h3>
      <div class="empty">
        <p>No interviews on file. Once you add one above, it appears here and the Brief picks it up.</p>
      </div>
    </div>

    <p class="footer-link"><a href="/preview/redesign">← back to previews</a></p>
  </div>
</main>

<style>
  :global(html, body) { background: var(--surface); margin: 0; font-family: var(--sans); color: var(--ink); }
  .wrap { max-width: 1080px; margin: 0 auto; padding: 0 28px 60px; }
  .topbar-stub { padding: 18px 0 14px; font-size: 13.5px; color: var(--mute); border-bottom: 1px solid var(--rule); margin-bottom: 24px; }
  .topbar-stub .root { color: var(--mute); cursor: pointer; }
  .topbar-stub .sep { color: var(--mute-2); margin: 0 6px; }
  .topbar-stub .here { color: var(--ink); font-weight: 500; }

  /* TABS */
  .tabs { display: flex; gap: 4px; border-bottom: 1px solid var(--rule); margin-bottom: 18px; }
  .tab { background: transparent; border: 0; padding: 10px 14px; font-size: 13.5px; color: var(--mute); cursor: pointer; border-bottom: 2px solid transparent; margin-bottom: -1px; font-weight: 600; }
  .tab.active { color: var(--ink); border-bottom-color: var(--ink); }
  .t-tag { font-size: 11px; background: var(--accent-tint); color: var(--accent-text); padding: 1px 7px; border-radius: 99px; margin-left: 4px; font-weight: 500; }

  /* ADD CARD */
  .add-card { background: var(--card); border: 1px solid var(--rule); border-radius: 14px; padding: 22px 24px; margin-bottom: 18px; box-shadow: var(--sh-1); }
  .add-hd { margin-bottom: 16px; }
  .add-hd h3 { font-size: 16px; font-weight: 600; margin: 0 0 4px; letter-spacing: -0.015em; }
  .add-hd p { font-size: 13.5px; color: var(--mute); margin: 0; line-height: 1.5; }

  /* TWO ZONES SIDE-BY-SIDE */
  .zones { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
  .zone { background: var(--surface); border: 1px solid var(--rule); border-radius: 12px; padding: 14px 16px; display: flex; flex-direction: column; gap: 10px; }
  .zone-hd { display: grid; grid-template-columns: 32px 1fr; gap: 10px; align-items: center; }
  .zone-ic { width: 32px; height: 32px; border-radius: 8px; background: var(--surface-2); color: var(--ink-2); display: grid; place-items: center; }
  .zone-ic.accent { background: var(--accent-tint); color: var(--accent-text); }
  .zone-title { font-size: 13.5px; font-weight: 600; color: var(--ink); display: inline-flex; align-items: center; gap: 6px; }
  .zone-sub { font-size: 12px; color: var(--mute); margin-top: 1px; }
  .ai-pill { font-size: 10px; font-weight: 600; color: var(--accent-text); background: var(--accent-tint); border-radius: 4px; padding: 1px 5px; letter-spacing: .04em; }

  .drop { background: var(--card); border: 1.5px dashed var(--rule-strong); border-radius: 10px; padding: 16px 12px; display: flex; flex-direction: column; align-items: center; gap: 3px; color: var(--mute); transition: border-color 120ms, background 120ms; cursor: pointer; }
  .drop:hover { border-color: var(--accent); background: var(--accent-tint); color: var(--accent-text); }
  .drop-l1 { font-size: 12.5px; font-weight: 500; color: var(--ink-2); }
  .drop:hover .drop-l1 { color: var(--accent-text); }
  .drop-l2 { font-size: 11.5px; color: var(--mute-2); }
  .drop kbd { font-family: var(--mono); font-size: 10px; background: var(--surface-2); border: 1px solid var(--rule); border-bottom-width: 2px; border-radius: 3px; padding: 0 4px; color: var(--ink-2); }

  .or { font-size: 11.5px; color: var(--mute-2); text-align: center; }

  textarea { width: 100%; font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 11.5px; line-height: 1.5; color: var(--ink); background: var(--card); border: 1px solid var(--rule); border-radius: 8px; padding: 8px 10px; outline: none; resize: vertical; box-sizing: border-box; }
  textarea:focus { border-color: var(--accent); }
  textarea::placeholder { color: var(--mute-2); }

  .card-actions { display: flex; gap: 8px; margin-top: 14px; }
  .btn { height: 32px; padding: 0 14px; border-radius: 7px; font: inherit; font-size: 13px; font-weight: 500; border: 1px solid var(--rule); background: var(--card); color: var(--ink-2); cursor: pointer; }
  .btn.primary { background: var(--accent); color: white; border-color: var(--accent); font-weight: 600; }
  .btn.primary:hover { background: var(--accent-strong); }
  .btn.ghost { background: transparent; border-color: transparent; color: var(--mute); }
  .btn.ghost:hover { background: var(--surface-2); color: var(--ink); }

  /* PREVIEW */
  .preview-wrap { margin-top: 18px; padding-top: 16px; border-top: 1px solid var(--rule); }
  .preview-wrap h4 { font-size: 11.5px; font-weight: 600; color: var(--mute); text-transform: uppercase; letter-spacing: .04em; margin: 0 0 10px; }
  .prev-row { background: var(--accent-tint); border: 1px solid var(--accent); border-radius: 10px; padding: 12px 14px; margin-bottom: 10px; }
  .prev-summary { font-size: 13.5px; font-weight: 600; color: var(--ink); }
  .prev-when { font-size: 12.5px; color: var(--accent-text); margin-top: 3px; font-weight: 500; }
  .prev-loc { font-size: 12px; color: var(--mute); margin-top: 4px; }
  .save { margin-top: 4px; }

  /* SCHEDULED */
  .scheduled { background: var(--card); border: 1px solid var(--rule); border-radius: 14px; padding: 22px 24px; }
  .scheduled h3 { font-size: 11.5px; font-weight: 600; color: var(--mute); text-transform: uppercase; letter-spacing: .04em; margin: 0 0 12px; }
  .empty { border: 1px dashed var(--rule); border-radius: 10px; padding: 20px; text-align: center; }
  .empty p { font-size: 13px; color: var(--mute); margin: 0; line-height: 1.5; }

  .footer-link { margin-top: 30px; font-size: 13px; }
  .footer-link a { color: var(--accent-text); text-decoration: none; }

  /* Mobile */
  @media (max-width: 720px) {
    .zones { grid-template-columns: 1fr; gap: 10px; }
    .wrap { padding: 0 18px 60px; }
  }
</style>
