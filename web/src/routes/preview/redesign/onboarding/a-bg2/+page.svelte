<script>
  // Variant A v2 — centered card, X close, fixed-height stage body,
  // larger nav buttons. Backdrop: dim translucent ink + blur (clean modal).
  let step = $state(0);
  const total = 5;
  let pickerMode = $state('pick');
  function next() { step = Math.min(total - 1, step + 1); }
  function prev() {
    if (step === total - 1 && pickerMode !== 'pick') { pickerMode = 'pick'; return; }
    step = Math.max(0, step - 1);
  }
  function close() { alert('Static mockup — in the real flow this closes the overlay and lands on /app with a tooltip on New application.'); }
</script>

<svelte:head><title>Onboarding · A v2 (mesh) — Pursuit</title></svelte:head>

<div class="ob-overlay bg-mesh">
  <div class="ob-card">
    <button class="x-close" onclick={close} aria-label="Close">
      <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="1.6">
        <path d="M3 3l8 8M11 3l-8 8" stroke-linecap="round"/>
      </svg>
    </button>

    <div class="dots">
      {#each Array(total) as _, i}
        <span class="dot" class:active={i === step} class:done={i < step}></span>
      {/each}
    </div>

    <div class="stage">
      <div class="art">
        {#if step === 0}
          <svg viewBox="0 0 240 200" width="240" height="200" fill="none">
            <circle cx="120" cy="105" r="78" stroke="var(--accent)" stroke-width="1" opacity="0.16"/>
            <circle cx="120" cy="105" r="58" stroke="var(--accent)" stroke-width="1.2" opacity="0.32"/>
            <circle cx="120" cy="105" r="38" stroke="var(--accent)" stroke-width="1.6" opacity="0.7"/>
            <circle cx="120" cy="105" r="20" fill="var(--accent-tint)" stroke="var(--accent)" stroke-width="1.6"/>
            <circle cx="150" cy="78" r="9" fill="var(--accent)"/>
          </svg>
        {:else if step === 1}
          <svg viewBox="0 0 320 200" width="320" height="200" fill="none">
            <rect x="14" y="34" width="124" height="132" rx="10" fill="var(--surface-2)" stroke="var(--rule)" stroke-width="1"/>
            <rect x="26" y="46" width="68" height="6" rx="3" fill="var(--mute-2)"/>
            <rect x="26" y="60" width="100" height="4" rx="2" fill="var(--rule-strong)"/>
            <rect x="26" y="70" width="92" height="4" rx="2" fill="var(--rule-strong)"/>
            <rect x="26" y="80" width="100" height="4" rx="2" fill="var(--rule-strong)"/>
            <rect x="26" y="90" width="78" height="4" rx="2" fill="var(--rule-strong)"/>
            <rect x="26" y="104" width="88" height="4" rx="2" fill="var(--rule-strong)"/>
            <rect x="26" y="114" width="100" height="4" rx="2" fill="var(--rule-strong)"/>
            <rect x="26" y="124" width="64" height="4" rx="2" fill="var(--rule-strong)"/>
            <path d="M148 100h22m-6-6 6 6-6 6" stroke="var(--accent)" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"/>
            <rect x="180" y="60" width="126" height="80" rx="10" fill="white" stroke="var(--accent)" stroke-width="1.2"/>
            <circle cx="196" cy="76" r="8" fill="var(--accent-tint)"/>
            <rect x="210" y="72" width="50" height="5" rx="2.5" fill="var(--ink)"/>
            <rect x="210" y="82" width="76" height="4" rx="2" fill="var(--mute-2)"/>
            <rect x="192" y="100" width="52" height="14" rx="7" fill="var(--accent-tint)"/>
            <text x="218" y="111" font-family="system-ui" font-size="9" fill="var(--accent-text)" font-weight="600">Applied</text>
          </svg>
        {:else if step === 2}
          <svg viewBox="0 0 320 200" width="320" height="200" fill="none">
            <rect x="20" y="40" width="120" height="130" rx="10" fill="var(--card)" stroke="var(--rule)" stroke-width="1"/>
            <rect x="20" y="40" width="120" height="16" rx="10" fill="var(--surface-2)"/>
            <circle cx="30" cy="48" r="2.5" fill="#ff6058"/><circle cx="38" cy="48" r="2.5" fill="#ffbd2e"/><circle cx="46" cy="48" r="2.5" fill="#28ca42"/>
            <rect x="30" y="66" width="100" height="38" rx="6" fill="var(--accent-tint)"/>
            <circle cx="46" cy="85" r="8" fill="white"/>
            <rect x="58" y="78" width="60" height="4" rx="2" fill="var(--accent-text)"/>
            <rect x="58" y="86" width="40" height="3" rx="1.5" fill="var(--accent)"/>
            <path d="M150 108h24m-6-5 6 5-6 5" stroke="var(--accent)" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"/>
            <rect x="184" y="56" width="124" height="98" rx="10" fill="white" stroke="var(--accent)" stroke-width="1.2"/>
            <circle cx="200" cy="74" r="9" fill="var(--accent-tint)"/>
            <rect x="214" y="70" width="54" height="5" rx="2.5" fill="var(--ink)"/>
            <rect x="194" y="100" width="60" height="14" rx="7" fill="var(--positive-tint)"/>
            <text x="222" y="111" font-family="system-ui" font-size="9" fill="var(--positive-text)" font-weight="600">Screen</text>
          </svg>
        {:else if step === 3}
          <svg viewBox="0 0 320 200" width="320" height="200" fill="none">
            <rect x="92" y="32" width="136" height="136" rx="14" fill="white" stroke="var(--rule)" stroke-width="1"/>
            <circle cx="160" cy="78" r="22" fill="var(--accent-tint)"/>
            <text x="160" y="85" text-anchor="middle" font-family="system-ui" font-weight="600" font-size="16" fill="var(--accent-text)">SC</text>
            <rect x="116" y="110" width="88" height="5" rx="2.5" fill="var(--ink)"/>
            <rect x="124" y="120" width="72" height="4" rx="2" fill="var(--mute-2)"/>
            <rect x="108" y="138" width="40" height="14" rx="7" fill="var(--accent-tint)"/>
            <text x="128" y="148.5" text-anchor="middle" font-family="system-ui" font-size="9" fill="var(--accent-text)" font-weight="600">Talk</text>
            <rect x="154" y="138" width="40" height="14" rx="7" fill="var(--positive-tint)"/>
            <text x="174" y="148.5" text-anchor="middle" font-family="system-ui" font-size="9" fill="var(--positive-text)" font-weight="600">Post</text>
            <rect x="18" y="46" width="60" height="34" rx="8" fill="white" stroke="var(--rule)" stroke-width="1"/>
            <rect x="26" y="54" width="20" height="3" rx="1.5" fill="var(--mute-2)"/>
            <rect x="26" y="62" width="40" height="3" rx="1.5" fill="var(--ink-2)"/>
            <rect x="242" y="58" width="60" height="34" rx="8" fill="white" stroke="var(--rule)" stroke-width="1"/>
            <rect x="250" y="66" width="22" height="3" rx="1.5" fill="var(--mute-2)"/>
            <rect x="250" y="74" width="40" height="3" rx="1.5" fill="var(--ink-2)"/>
            <rect x="34" y="118" width="60" height="34" rx="8" fill="white" stroke="var(--rule)" stroke-width="1"/>
            <rect x="42" y="126" width="22" height="3" rx="1.5" fill="var(--mute-2)"/>
            <rect x="42" y="134" width="40" height="3" rx="1.5" fill="var(--ink-2)"/>
          </svg>
        {:else}
          <svg viewBox="0 0 320 200" width="320" height="200" fill="none">
            <rect x="14" y="48" width="92" height="108" rx="12" fill="var(--surface)" stroke="var(--rule)" stroke-width="1"/>
            <circle cx="60" cy="84" r="20" fill="var(--accent-tint)"/>
            <path d="M52 80h16M52 84h12M52 88h14" stroke="var(--accent-text)" stroke-width="1.6" stroke-linecap="round"/>
            <rect x="32" y="118" width="56" height="4" rx="2" fill="var(--ink)"/>
            <rect x="114" y="44" width="92" height="112" rx="12" fill="var(--accent-tint)" stroke="var(--accent)" stroke-width="1.2"/>
            <circle cx="160" cy="80" r="20" fill="white"/>
            <path d="M150 76h20M150 80h20M150 84h14" stroke="var(--accent)" stroke-width="1.6" stroke-linecap="round"/>
            <path d="M172 70l.9 1.7 1.7.9-1.7.9-.9 1.7-.9-1.7-1.7-.9 1.7-.9z" fill="var(--accent)"/>
            <rect x="132" y="116" width="56" height="4" rx="2" fill="var(--accent-text)"/>
            <text x="160" y="148" text-anchor="middle" font-family="system-ui" font-size="9" font-weight="600" fill="var(--accent-text)">RECOMMENDED</text>
            <rect x="214" y="48" width="92" height="108" rx="12" fill="var(--surface)" stroke="var(--rule)" stroke-width="1"/>
            <circle cx="260" cy="84" r="20" fill="var(--warm-tint)"/>
            <circle cx="260" cy="84" r="10" stroke="var(--warm-text)" stroke-width="1.6" fill="none"/>
            <path d="M268 92l5 5" stroke="var(--warm-text)" stroke-width="1.6" stroke-linecap="round"/>
            <rect x="232" y="118" width="56" height="4" rx="2" fill="var(--ink)"/>
          </svg>
        {/if}
      </div>

      <div class="text">
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
          {#if pickerMode === 'pick'}
            <h1>Pick where to start.</h1>
            <p class="picker-help">All three paths land you on the dashboard. You can clear demo data any time.</p>
            <div class="picker">
              <button class="pick" onclick={() => (pickerMode = 'paste')}>
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
          {:else}
            <h1>Paste your applications.</h1>
            <p class="picker-help">Separate entries with a blank line.</p>
            <textarea rows="5" placeholder={"https://job-boards.greenhouse.io/anthropic/jobs/4020693008\n\nStripe — Staff Backend Engineer · referred by Mia"}></textarea>
          {/if}
        {/if}
      </div>
    </div>

    <footer>
      {#if step > 0 && (step < total - 1 || pickerMode !== 'pick')}
        <button class="btn ghost" onclick={prev}>← Back</button>
      {/if}
      <span style="flex:1"></span>
      {#if step < total - 1}
        <button class="btn primary" onclick={next}>Next →</button>
      {:else if pickerMode === 'paste'}
        <button class="btn primary">Parse + open dashboard →</button>
      {/if}
    </footer>
  </div>
</div>

<p class="back-link"><a href="/preview/redesign">← back to previews</a></p>

<style>
  :global(html, body) { background: var(--surface); margin: 0; font-family: var(--sans); color: var(--ink); }
  .ob-overlay {
    min-height: 100vh;
    display: grid; place-items: center;
    padding: 24px;
  }
  /* BACKDROP — soft mesh corners (warm light, no center glow). */
  .bg-mesh {
    background:
      radial-gradient(60vw 60vw at 0% 0%,   var(--accent-tint)   0%, transparent 55%),
      radial-gradient(60vw 60vw at 100% 100%, var(--warm-tint)   0%, transparent 55%),
      radial-gradient(50vw 50vw at 100% 0%, var(--positive-tint) 0%, transparent 55%),
      var(--surface);
  }
  .ob-card {
    position: relative;
    width: 100%; max-width: 560px;
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 18px;
    padding: 28px 30px 22px;
    box-shadow: 0 24px 80px -8px rgba(10,10,13,0.30), var(--sh-1);
    display: flex; flex-direction: column;
  }
  .x-close {
    position: absolute; top: 14px; right: 14px;
    width: 28px; height: 28px; border-radius: 8px;
    background: transparent; border: 0;
    display: grid; place-items: center;
    color: var(--mute); cursor: pointer;
    transition: background 100ms, color 100ms;
  }
  .x-close:hover { background: var(--surface-2); color: var(--ink); }
  .dots { display: flex; justify-content: center; gap: 6px; margin-bottom: 22px; }
  .dot { width: 6px; height: 6px; border-radius: 50%; background: var(--rule-strong); transition: background 200ms, width 200ms; }
  .dot.done { background: var(--accent); }
  .dot.active { background: var(--accent); width: 24px; border-radius: 99px; }

  /* Fixed-height stage so the card doesn't jiggle between steps. */
  .stage { flex: 1; min-height: 380px; display: flex; flex-direction: column; }
  .art { display: grid; place-items: center; height: 200px; margin-bottom: 18px; flex-shrink: 0; }

  .text h1 { font-size: 22px; font-weight: 600; letter-spacing: -0.022em; margin: 0 0 8px; text-align: center; }
  .text > p { font-size: 14px; color: var(--mute); line-height: 1.55; margin: 0 auto 12px; text-align: center; max-width: 44ch; }
  .text kbd { font-family: var(--mono); font-size: 11px; background: var(--surface-2); border: 1px solid var(--rule); border-bottom-width: 2px; border-radius: 3px; padding: 0 4px; color: var(--ink-2); }
  .picker-help { margin-bottom: 16px !important; }

  .picker { display: flex; flex-direction: column; gap: 8px; margin: 4px 0 8px; }
  .pick { position: relative; text-align: left; background: var(--surface); border: 1px solid var(--rule); border-radius: 10px; padding: 12px 14px; cursor: pointer; transition: border-color 120ms, transform 120ms; font: inherit; color: inherit; width: 100%; }
  .pick:hover { border-color: var(--accent); transform: translateY(-1px); }
  .pick.rec { background: var(--accent-tint); border-color: var(--accent); }
  .pick-tag { position: absolute; top: 10px; right: 12px; font-size: 10px; font-weight: 600; letter-spacing: 0.06em; color: var(--accent-text); text-transform: uppercase; }
  .pick-h { font-size: 14px; font-weight: 600; color: var(--ink); }
  .pick.rec .pick-h { color: var(--accent-text); }
  .pick-s { font-size: 12.5px; color: var(--mute); margin-top: 2px; }
  .pick.rec .pick-s { color: var(--accent-text); opacity: 0.8; }

  textarea { width: 100%; font: inherit; font-size: 13.5px; line-height: 1.5; color: var(--ink); background: var(--surface); border: 1px solid var(--rule); border-radius: 10px; padding: 12px 14px; outline: none; resize: vertical; font-family: var(--sans); }
  textarea:focus { border-color: var(--accent); }

  footer { display: flex; align-items: center; gap: 8px; margin-top: 18px; border-top: 1px solid var(--rule); padding-top: 16px; flex-shrink: 0; }
  .btn { height: 36px; padding: 0 16px; border-radius: 8px; font: inherit; font-size: 13.5px; font-weight: 500; border: 1px solid var(--rule); background: var(--card); color: var(--ink-2); cursor: pointer; }
  .btn.ghost { background: transparent; border-color: transparent; color: var(--mute); padding: 0 12px; }
  .btn.ghost:hover { color: var(--ink); background: var(--surface-2); }
  .btn.primary { background: var(--accent); color: white; border-color: var(--accent); font-weight: 600; padding: 0 18px; }
  .btn.primary:hover { background: var(--accent-strong); }

  .back-link { position: fixed; bottom: 10px; left: 12px; font-size: 12.5px; }
  .back-link a { color: var(--accent-text); text-decoration: none; background: var(--card); padding: 4px 10px; border-radius: 999px; border: 1px solid var(--rule); }
</style>
