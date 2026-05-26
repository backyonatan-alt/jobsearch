<script>
  // Variant C — Two-column split. Big illustration on the left over a
  // colored panel; text + actions on the right. Step pill top-right of the
  // right column.
  let step = $state(0);
  const total = 5;
  function next() { step = Math.min(total - 1, step + 1); }
  function prev() { step = Math.max(0, step - 1); }
  function close() { alert('Static mockup — final wiring lands the user on /app with a tooltip on the New-application button.'); }
</script>

<svelte:head><title>Onboarding · C — Pursuit</title></svelte:head>

<div class="ob-overlay">
  <div class="ob-card">
    <!-- LEFT: illustration panel -->
    <aside class="art-panel" data-step={step}>
      <div class="art-bg"></div>
      <div class="art-inner">
        {#if step === 0}
          <svg viewBox="0 0 320 320" width="280" height="280" fill="none">
            <circle cx="160" cy="160" r="140" stroke="var(--accent)" stroke-width="1" opacity="0.18"/>
            <circle cx="160" cy="160" r="106" stroke="var(--accent)" stroke-width="1.4" opacity="0.32"/>
            <circle cx="160" cy="160" r="74" stroke="var(--accent)" stroke-width="1.6" opacity="0.5"/>
            <circle cx="160" cy="160" r="46" stroke="var(--accent)" stroke-width="1.8" opacity="0.85"/>
            <circle cx="160" cy="160" r="22" fill="var(--accent-tint)" stroke="var(--accent)" stroke-width="2.5"/>
            <circle cx="218" cy="108" r="16" fill="var(--accent)"/>
            <path d="M210 108l5 5 11-11" stroke="white" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        {:else if step === 1}
          <svg viewBox="0 0 320 320" width="280" height="280" fill="none">
            <g transform="translate(14 60) rotate(-4 70 80)">
              <rect width="150" height="200" rx="14" fill="white" stroke="var(--rule)" stroke-width="1"/>
              <rect x="14" y="14" width="68" height="6" rx="3" fill="var(--ink)"/>
              <rect x="14" y="30" width="120" height="4" rx="2" fill="var(--rule-strong)"/>
              <rect x="14" y="40" width="106" height="4" rx="2" fill="var(--rule-strong)"/>
              <rect x="14" y="50" width="114" height="4" rx="2" fill="var(--rule-strong)"/>
              <rect x="14" y="60" width="80" height="4" rx="2" fill="var(--rule-strong)"/>
              <rect x="14" y="78" width="120" height="4" rx="2" fill="var(--rule-strong)"/>
              <rect x="14" y="88" width="100" height="4" rx="2" fill="var(--rule-strong)"/>
              <rect x="14" y="98" width="110" height="4" rx="2" fill="var(--rule-strong)"/>
              <rect x="14" y="108" width="64" height="4" rx="2" fill="var(--rule-strong)"/>
              <rect x="14" y="124" width="120" height="4" rx="2" fill="var(--rule-strong)"/>
              <rect x="14" y="134" width="84" height="4" rx="2" fill="var(--rule-strong)"/>
              <rect x="14" y="158" width="64" height="16" rx="8" fill="var(--ink)"/>
            </g>
            <g transform="translate(170 56) rotate(5 70 80)">
              <rect width="150" height="190" rx="14" fill="white" stroke="var(--accent)" stroke-width="2"/>
              <circle cx="26" cy="28" r="12" fill="var(--accent-tint)"/>
              <rect x="44" y="22" width="62" height="6" rx="3" fill="var(--ink)"/>
              <rect x="44" y="34" width="86" height="4" rx="2" fill="var(--mute-2)"/>
              <rect x="14" y="60" width="70" height="18" rx="9" fill="var(--positive-tint)"/>
              <text x="49" y="72" text-anchor="middle" font-family="system-ui" font-size="11" fill="var(--positive-text)" font-weight="600">Applied</text>
              <rect x="14" y="92" width="60" height="4" rx="2" fill="var(--mute-2)"/>
              <rect x="14" y="102" width="50" height="4" rx="2" fill="var(--mute-2)"/>
              <rect x="14" y="126" width="120" height="3" rx="1.5" fill="var(--rule)"/>
              <rect x="14" y="134" width="100" height="3" rx="1.5" fill="var(--rule)"/>
              <rect x="14" y="142" width="110" height="3" rx="1.5" fill="var(--rule)"/>
              <rect x="14" y="150" width="80" height="3" rx="1.5" fill="var(--rule)"/>
            </g>
            <g transform="translate(150 144)">
              <path d="M0 0l3.4 7.6 7.6 3.4-7.6 3.4-3.4 7.6-3.4-7.6L-11 11l7.6-3.4z" fill="var(--accent)"/>
            </g>
          </svg>
        {:else if step === 2}
          <svg viewBox="0 0 320 320" width="280" height="280" fill="none">
            <g transform="translate(36 50)">
              <rect width="120" height="220" rx="18" fill="var(--card)" stroke="var(--rule)" stroke-width="2"/>
              <rect x="10" y="10" width="100" height="18" rx="9" fill="var(--surface-2)"/>
              <rect x="22" y="16" width="40" height="5" rx="2.5" fill="var(--mute)"/>
              <rect x="10" y="40" width="100" height="58" rx="10" fill="var(--accent-tint)"/>
              <circle cx="26" cy="68" r="10" fill="white"/>
              <rect x="40" y="60" width="58" height="5" rx="2.5" fill="var(--accent-text)"/>
              <rect x="40" y="70" width="46" height="4" rx="2" fill="var(--accent)"/>
              <rect x="40" y="78" width="44" height="4" rx="2" fill="var(--accent)" opacity="0.6"/>
              <rect x="10" y="110" width="100" height="4" rx="2" fill="var(--rule-strong)"/>
              <rect x="10" y="120" width="92" height="4" rx="2" fill="var(--rule-strong)"/>
              <rect x="10" y="130" width="96" height="4" rx="2" fill="var(--rule-strong)"/>
              <rect x="10" y="140" width="70" height="4" rx="2" fill="var(--rule-strong)"/>
              <rect x="10" y="156" width="100" height="4" rx="2" fill="var(--rule-strong)"/>
              <rect x="10" y="166" width="80" height="4" rx="2" fill="var(--rule-strong)"/>
              <rect x="10" y="188" width="60" height="14" rx="7" fill="var(--ink)"/>
            </g>
            <g transform="translate(180 130)">
              <path d="M0 0l4 9 9 4-9 4-4 9-4-9-9-4 9-4z" fill="var(--accent)"/>
              <path d="M30 -28l1.6 3.6 3.6 1.6-3.6 1.6L30 -17l-1.6-3.6-3.6-1.6 3.6-1.6z" fill="var(--accent)" opacity="0.8"/>
              <path d="M-8 36l1.2 2.6 2.6 1.2-2.6 1.2L-8 43.6l-1.2-2.6-2.6-1.2 2.6-1.2z" fill="var(--accent)" opacity="0.6"/>
            </g>
            <g transform="translate(196 80)">
              <rect width="92" height="160" rx="14" fill="white" stroke="var(--accent)" stroke-width="2"/>
              <circle cx="20" cy="22" r="10" fill="var(--accent-tint)"/>
              <rect x="36" y="18" width="44" height="6" rx="3" fill="var(--ink)"/>
              <rect x="36" y="30" width="50" height="4" rx="2" fill="var(--mute-2)"/>
              <rect x="12" y="50" width="60" height="14" rx="7" fill="var(--positive-tint)"/>
              <text x="42" y="60" text-anchor="middle" font-family="system-ui" font-size="9" fill="var(--positive-text)" font-weight="600">Screen</text>
              <rect x="12" y="78" width="70" height="3.5" rx="1.75" fill="var(--mute-2)"/>
              <rect x="12" y="86" width="56" height="3.5" rx="1.75" fill="var(--mute-2)"/>
              <rect x="12" y="104" width="70" height="3" rx="1.5" fill="var(--rule)"/>
              <rect x="12" y="112" width="60" height="3" rx="1.5" fill="var(--rule)"/>
              <rect x="12" y="120" width="68" height="3" rx="1.5" fill="var(--rule)"/>
              <rect x="12" y="128" width="44" height="3" rx="1.5" fill="var(--rule)"/>
              <rect x="12" y="142" width="50" height="10" rx="5" fill="var(--ink)"/>
            </g>
          </svg>
        {:else if step === 3}
          <svg viewBox="0 0 320 320" width="280" height="280" fill="none">
            <path d="M76 110 L 158 158" stroke="var(--rule-strong)" stroke-width="1" stroke-dasharray="2 3"/>
            <path d="M76 230 L 158 174" stroke="var(--rule-strong)" stroke-width="1" stroke-dasharray="2 3"/>
            <path d="M260 110 L 200 158" stroke="var(--rule-strong)" stroke-width="1" stroke-dasharray="2 3"/>
            <path d="M260 224 L 200 172" stroke="var(--rule-strong)" stroke-width="1" stroke-dasharray="2 3"/>
            <!-- signals -->
            <g transform="translate(14 80)">
              <rect width="76" height="44" rx="10" fill="white" stroke="var(--rule)" stroke-width="1"/>
              <rect x="10" y="10" width="22" height="3.5" rx="1.75" fill="var(--accent)"/>
              <rect x="10" y="20" width="56" height="3.5" rx="1.75" fill="var(--ink-2)"/>
              <rect x="10" y="28" width="42" height="3.5" rx="1.75" fill="var(--ink-2)"/>
            </g>
            <g transform="translate(14 196)">
              <rect width="76" height="44" rx="10" fill="white" stroke="var(--rule)" stroke-width="1"/>
              <rect x="10" y="10" width="22" height="3.5" rx="1.75" fill="var(--accent)"/>
              <rect x="10" y="20" width="56" height="3.5" rx="1.75" fill="var(--ink-2)"/>
              <rect x="10" y="28" width="40" height="3.5" rx="1.75" fill="var(--ink-2)"/>
            </g>
            <g transform="translate(230 80)">
              <rect width="76" height="44" rx="10" fill="white" stroke="var(--rule)" stroke-width="1"/>
              <rect x="10" y="10" width="22" height="3.5" rx="1.75" fill="var(--accent)"/>
              <rect x="10" y="20" width="56" height="3.5" rx="1.75" fill="var(--ink-2)"/>
              <rect x="10" y="28" width="48" height="3.5" rx="1.75" fill="var(--ink-2)"/>
            </g>
            <g transform="translate(230 196)">
              <rect width="76" height="44" rx="10" fill="white" stroke="var(--rule)" stroke-width="1"/>
              <rect x="10" y="10" width="22" height="3.5" rx="1.75" fill="var(--accent)"/>
              <rect x="10" y="20" width="56" height="3.5" rx="1.75" fill="var(--ink-2)"/>
              <rect x="10" y="28" width="38" height="3.5" rx="1.75" fill="var(--ink-2)"/>
            </g>
            <!-- interviewer card -->
            <rect x="98" y="96" width="124" height="138" rx="16" fill="white" stroke="var(--accent)" stroke-width="2"/>
            <circle cx="160" cy="138" r="26" fill="var(--accent-tint)"/>
            <text x="160" y="146" text-anchor="middle" font-family="system-ui" font-weight="700" font-size="18" fill="var(--accent-text)">SC</text>
            <rect x="118" y="174" width="84" height="6" rx="3" fill="var(--ink)"/>
            <rect x="126" y="186" width="68" height="4" rx="2" fill="var(--mute)"/>
            <rect x="108" y="204" width="48" height="14" rx="7" fill="var(--accent-tint)"/>
            <text x="132" y="214" text-anchor="middle" font-family="system-ui" font-size="9" fill="var(--accent-text)" font-weight="600">Talk</text>
            <rect x="166" y="204" width="48" height="14" rx="7" fill="var(--positive-tint)"/>
            <text x="190" y="214" text-anchor="middle" font-family="system-ui" font-size="9" fill="var(--positive-text)" font-weight="600">Post</text>
          </svg>
        {:else}
          <svg viewBox="0 0 320 320" width="280" height="280" fill="none">
            <g transform="translate(20 60) rotate(-4 50 80)">
              <rect width="100" height="180" rx="16" fill="white" stroke="var(--rule)" stroke-width="1"/>
              <circle cx="50" cy="62" r="26" fill="var(--accent-tint)"/>
              <path d="M38 58h24M38 62h18M38 66h22" stroke="var(--accent-text)" stroke-width="2" stroke-linecap="round"/>
              <rect x="20" y="106" width="60" height="5" rx="2.5" fill="var(--ink)"/>
              <rect x="26" y="120" width="48" height="4" rx="2" fill="var(--mute)"/>
              <rect x="26" y="130" width="44" height="4" rx="2" fill="var(--mute)"/>
            </g>
            <g transform="translate(110 38)">
              <rect width="100" height="218" rx="18" fill="var(--accent-tint)" stroke="var(--accent)" stroke-width="2.5"/>
              <circle cx="50" cy="74" r="30" fill="white"/>
              <path d="M36 70h28M36 74h28M36 78h22" stroke="var(--accent)" stroke-width="2" stroke-linecap="round"/>
              <path d="M70 50l1.4 3 3 1.4-3 1.4-1.4 3-1.4-3-3-1.4 3-1.4z" fill="var(--accent)"/>
              <rect x="18" y="124" width="64" height="6" rx="3" fill="var(--accent-text)"/>
              <rect x="24" y="140" width="52" height="4" rx="2" fill="var(--accent-text)" opacity="0.7"/>
              <rect x="24" y="150" width="44" height="4" rx="2" fill="var(--accent-text)" opacity="0.7"/>
              <text x="50" y="180" text-anchor="middle" font-family="system-ui" font-size="9" font-weight="700" fill="var(--accent-text)" letter-spacing="0.05em">RECOMMENDED</text>
            </g>
            <g transform="translate(200 60) rotate(4 50 80)">
              <rect width="100" height="180" rx="16" fill="white" stroke="var(--rule)" stroke-width="1"/>
              <circle cx="50" cy="62" r="26" fill="var(--warm-tint)"/>
              <circle cx="50" cy="62" r="14" stroke="var(--warm-text)" stroke-width="2.5" fill="none"/>
              <path d="M60 72l8 8" stroke="var(--warm-text)" stroke-width="2.5" stroke-linecap="round"/>
              <rect x="20" y="106" width="60" height="5" rx="2.5" fill="var(--ink)"/>
              <rect x="26" y="120" width="48" height="4" rx="2" fill="var(--mute)"/>
              <rect x="26" y="130" width="40" height="4" rx="2" fill="var(--mute)"/>
            </g>
          </svg>
        {/if}
      </div>
    </aside>

    <!-- RIGHT: text + actions -->
    <section class="text-panel">
      <div class="step-pill">{step + 1} of {total}</div>

      {#if step === 0}
        <h1>Welcome to Pursuit.</h1>
        <p>Track your job search end-to-end — applications, interviews, and AI briefings on the people you're about to meet.</p>
      {:else if step === 1}
        <h1>Paste any job, get a row.</h1>
        <p>Drop a JD URL, a body of text, or a recruiter email. Claude parses it into a clean application row in about three seconds.</p>
      {:else if step === 2}
        <h1>Screenshot, no scraping.</h1>
        <p>LinkedIn blocks bots — so just <kbd>⌘⌃⇧4</kbd> the page and paste. Vision reads the screenshot and fills the fields for you.</p>
      {:else if step === 3}
        <h1>A briefing on your interviewer.</h1>
        <p>Before any conversation, generate a 60-second brief: recent talks and posts, what tends to land, what to avoid, three questions worth asking.</p>
      {:else}
        <h1>Pick where to start.</h1>
        <p class="picker-help">All three paths get you to the dashboard. You can clear demo data any time.</p>
        <div class="picker">
          <button class="pick">
            <div class="pick-h">Paste your applications</div>
            <div class="pick-s">URLs, text, emails — Claude parses them in.</div>
          </button>
          <button class="pick rec">
            <span class="pick-tag">Recommended</span>
            <div class="pick-h">Try with demo data</div>
            <div class="pick-s">15 realistic apps so you can see every surface populated.</div>
          </button>
          <button class="pick">
            <div class="pick-h">Just explore</div>
            <div class="pick-s">Land on an empty dashboard. Add by hand later.</div>
          </button>
        </div>
      {/if}

      <footer>
        {#if step > 0 && step < total - 1}
          <button class="btn ghost" onclick={prev}>← Back</button>
        {/if}
        <span style="flex:1"></span>
        {#if step < total - 1}
          <button class="btn ghost" onclick={close}>Skip</button>
          <button class="btn primary" onclick={next}>Next →</button>
        {/if}
      </footer>
    </section>
  </div>
</div>

<p class="back-link"><a href="/preview/redesign">← back to previews</a></p>

<style>
  :global(html, body) { background: var(--surface); margin: 0; font-family: var(--sans); color: var(--ink); }
  .ob-overlay {
    min-height: 100vh;
    display: grid; place-items: center;
    padding: 32px;
    background: var(--surface);
  }
  .ob-card {
    width: 100%;
    max-width: 920px;
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 20px;
    overflow: hidden;
    display: grid;
    grid-template-columns: 1fr 1fr;
    min-height: 520px;
    box-shadow: 0 24px 80px -32px rgba(10,10,13,0.20), var(--sh-1);
  }

  /* Left art panel */
  .art-panel {
    position: relative;
    overflow: hidden;
    display: grid; place-items: center;
    border-right: 1px solid var(--rule);
  }
  .art-bg {
    position: absolute; inset: 0;
    background: var(--accent-tint);
    transition: background 240ms ease;
  }
  .art-panel[data-step="0"] .art-bg { background: var(--accent-tint); }
  .art-panel[data-step="1"] .art-bg { background: var(--surface-2); }
  .art-panel[data-step="2"] .art-bg { background: var(--accent-tint); }
  .art-panel[data-step="3"] .art-bg { background: var(--positive-tint); opacity: 0.5; }
  .art-panel[data-step="4"] .art-bg { background: var(--warm-tint); opacity: 0.45; }
  .art-inner { position: relative; padding: 24px; }

  /* Right text panel */
  .text-panel {
    position: relative;
    padding: 40px 36px 24px;
    display: flex; flex-direction: column;
  }
  .step-pill {
    position: absolute; top: 22px; right: 28px;
    font-size: 11px; font-weight: 500;
    color: var(--mute);
    font-family: var(--mono);
    padding: 3px 9px;
    background: var(--surface-2);
    border-radius: 99px;
    letter-spacing: 0.04em;
  }
  .text-panel h1 {
    font-size: 24px; font-weight: 600;
    letter-spacing: -0.022em;
    margin: 0 0 10px;
  }
  .text-panel > p {
    font-size: 14.5px; color: var(--mute);
    line-height: 1.6;
    margin: 0 0 16px;
    max-width: 38ch;
  }
  .text-panel kbd {
    font-family: var(--mono); font-size: 11px;
    background: var(--surface-2); border: 1px solid var(--rule);
    border-bottom-width: 2px; border-radius: 3px;
    padding: 0 4px; color: var(--ink-2);
  }
  .picker-help { margin-bottom: 14px !important; }

  .picker { display: flex; flex-direction: column; gap: 8px; margin: 4px 0 8px; }
  .pick {
    position: relative;
    text-align: left;
    background: var(--surface);
    border: 1px solid var(--rule);
    border-radius: 10px;
    padding: 12px 14px;
    cursor: pointer;
    transition: border-color 120ms, transform 120ms;
  }
  .pick:hover { border-color: var(--accent); transform: translateY(-1px); }
  .pick.rec { background: var(--accent-tint); border-color: var(--accent); }
  .pick-tag {
    position: absolute; top: 10px; right: 12px;
    font-size: 10px; font-weight: 600; letter-spacing: 0.06em;
    color: var(--accent-text); text-transform: uppercase;
  }
  .pick-h { font-size: 14px; font-weight: 600; color: var(--ink); }
  .pick.rec .pick-h { color: var(--accent-text); }
  .pick-s { font-size: 12.5px; color: var(--mute); margin-top: 2px; }
  .pick.rec .pick-s { color: var(--accent-text); opacity: 0.8; }

  footer {
    margin-top: auto;
    padding-top: 18px;
    display: flex; align-items: center; gap: 8px;
    border-top: 1px solid var(--rule);
  }
  .btn {
    height: 32px;
    padding: 0 14px;
    border-radius: 7px;
    font: inherit; font-size: 13px; font-weight: 500;
    border: 1px solid var(--rule);
    background: var(--card);
    color: var(--ink-2);
    cursor: pointer;
  }
  .btn.ghost { background: transparent; border-color: transparent; color: var(--mute); }
  .btn.ghost:hover { color: var(--ink); background: var(--surface-2); }
  .btn.primary { background: var(--accent); color: white; border-color: var(--accent); }
  .btn.primary:hover { background: var(--accent-strong); }

  .back-link { text-align: center; margin: 14px 0 30px; font-size: 12.5px; }
  .back-link a { color: var(--accent-text); text-decoration: none; }

  @media (max-width: 720px) {
    .ob-card { grid-template-columns: 1fr; min-height: 0; }
    .art-panel { border-right: 0; border-bottom: 1px solid var(--rule); min-height: 280px; }
    .text-panel { padding: 28px 22px 18px; }
    .step-pill { top: 18px; right: 18px; }
  }
</style>
