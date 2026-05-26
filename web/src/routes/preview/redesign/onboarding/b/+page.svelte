<script>
  // Variant B — Card with two "peek" cards behind hinting at progress.
  // Counter "1 / 5" sits subtly top-right. Bigger illustration well.
  let step = $state(0);
  const total = 5;
  function next() { step = Math.min(total - 1, step + 1); }
  function prev() { step = Math.max(0, step - 1); }
  function close() { alert('Static mockup — final wiring lands the user on /app with a tooltip on the New-application button.'); }
</script>

<svelte:head><title>Onboarding · B — Pursuit</title></svelte:head>

<div class="ob-overlay">
  <div class="stack">
    <!-- Peek cards behind -->
    <div class="peek peek-2"></div>
    <div class="peek peek-1"></div>

    <div class="ob-card">
      <div class="counter">{step + 1} / {total}</div>

      <div class="art">
        {#if step === 0}
          <!-- Welcome — large target with concentric rings -->
          <svg viewBox="0 0 320 240" width="320" height="240" fill="none">
            <circle cx="160" cy="120" r="100" stroke="var(--accent)" stroke-width="1" opacity="0.12"/>
            <circle cx="160" cy="120" r="76" stroke="var(--accent)" stroke-width="1.2" opacity="0.22"/>
            <circle cx="160" cy="120" r="54" stroke="var(--accent)" stroke-width="1.4" opacity="0.4"/>
            <circle cx="160" cy="120" r="34" stroke="var(--accent)" stroke-width="1.6" opacity="0.8"/>
            <circle cx="160" cy="120" r="16" fill="var(--accent-tint)" stroke="var(--accent)" stroke-width="2"/>
            <circle cx="200" cy="84" r="12" fill="var(--accent)"/>
            <path d="M194 84l4 4 8-8" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        {:else if step === 1}
          <!-- Paste & parse — overlapping cards -->
          <svg viewBox="0 0 320 240" width="320" height="240" fill="none">
            <!-- JD source -->
            <g transform="translate(20 36) rotate(-3 70 84)">
              <rect width="140" height="170" rx="12" fill="var(--surface-2)" stroke="var(--rule)" stroke-width="1"/>
              <rect x="14" y="14" width="64" height="6" rx="3" fill="var(--ink-2)"/>
              <rect x="14" y="28" width="112" height="3.5" rx="1.75" fill="var(--rule-strong)"/>
              <rect x="14" y="36" width="100" height="3.5" rx="1.75" fill="var(--rule-strong)"/>
              <rect x="14" y="44" width="108" height="3.5" rx="1.75" fill="var(--rule-strong)"/>
              <rect x="14" y="58" width="80" height="3.5" rx="1.75" fill="var(--rule-strong)"/>
              <rect x="14" y="72" width="108" height="3.5" rx="1.75" fill="var(--rule-strong)"/>
              <rect x="14" y="80" width="100" height="3.5" rx="1.75" fill="var(--rule-strong)"/>
              <rect x="14" y="88" width="64" height="3.5" rx="1.75" fill="var(--rule-strong)"/>
              <rect x="14" y="102" width="112" height="3.5" rx="1.75" fill="var(--rule-strong)"/>
              <rect x="14" y="110" width="88" height="3.5" rx="1.75" fill="var(--rule-strong)"/>
              <rect x="14" y="118" width="100" height="3.5" rx="1.75" fill="var(--rule-strong)"/>
              <rect x="14" y="138" width="60" height="14" rx="7" fill="var(--ink)"/>
              <text x="44" y="148.5" text-anchor="middle" font-family="system-ui" font-size="9" fill="white" font-weight="600">Apply</text>
            </g>
            <!-- Sparkle pile -->
            <g transform="translate(150 110)">
              <path d="M0 0l2.4 5.6 5.6 2.4-5.6 2.4-2.4 5.6-2.4-5.6L-8 8l5.6-2.4z" fill="var(--accent)"/>
              <path d="M22 -28l1.4 3.2 3.2 1.4-3.2 1.4L22 -19l-1.4-3.2-3.2-1.4 3.2-1.4z" fill="var(--accent)" opacity="0.7"/>
              <path d="M20 36l1 2.4 2.4 1-2.4 1-1 2.4-1-2.4-2.4-1 2.4-1z" fill="var(--accent)" opacity="0.5"/>
            </g>
            <!-- Parsed row card -->
            <g transform="translate(174 50) rotate(4 70 84)">
              <rect width="140" height="170" rx="12" fill="white" stroke="var(--accent)" stroke-width="1.6"/>
              <circle cx="22" cy="24" r="10" fill="var(--accent-tint)"/>
              <rect x="38" y="20" width="56" height="6" rx="3" fill="var(--ink)"/>
              <rect x="38" y="32" width="78" height="4" rx="2" fill="var(--mute-2)"/>
              <rect x="14" y="58" width="64" height="16" rx="8" fill="var(--positive-tint)"/>
              <text x="46" y="69" text-anchor="middle" font-family="system-ui" font-size="10" fill="var(--positive-text)" font-weight="600">Applied</text>
              <rect x="14" y="86" width="50" height="3.5" rx="1.75" fill="var(--mute-2)"/>
              <rect x="14" y="94" width="44" height="3.5" rx="1.75" fill="var(--mute-2)"/>
              <rect x="14" y="116" width="112" height="3" rx="1.5" fill="var(--rule)"/>
              <rect x="14" y="124" width="98" height="3" rx="1.5" fill="var(--rule)"/>
              <rect x="14" y="132" width="104" height="3" rx="1.5" fill="var(--rule)"/>
              <rect x="14" y="146" width="80" height="3" rx="1.5" fill="var(--rule)"/>
            </g>
          </svg>
        {:else if step === 2}
          <!-- Screenshot → row -->
          <svg viewBox="0 0 320 240" width="320" height="240" fill="none">
            <!-- Phone-style screenshot -->
            <g transform="translate(46 30)">
              <rect width="100" height="180" rx="14" fill="var(--card)" stroke="var(--rule)" stroke-width="1.5"/>
              <rect x="8" y="8" width="84" height="14" rx="7" fill="var(--surface-2)"/>
              <rect x="18" y="13" width="34" height="4" rx="2" fill="var(--mute)"/>
              <rect x="8" y="32" width="84" height="50" rx="8" fill="var(--accent-tint)"/>
              <circle cx="22" cy="58" r="9" fill="white"/>
              <rect x="34" y="50" width="50" height="4" rx="2" fill="var(--accent-text)"/>
              <rect x="34" y="58" width="40" height="3" rx="1.5" fill="var(--accent)"/>
              <rect x="34" y="66" width="38" height="3" rx="1.5" fill="var(--accent)" opacity="0.6"/>
              <rect x="8" y="92" width="84" height="3" rx="1.5" fill="var(--rule-strong)"/>
              <rect x="8" y="100" width="76" height="3" rx="1.5" fill="var(--rule-strong)"/>
              <rect x="8" y="108" width="80" height="3" rx="1.5" fill="var(--rule-strong)"/>
              <rect x="8" y="116" width="60" height="3" rx="1.5" fill="var(--rule-strong)"/>
              <rect x="8" y="130" width="84" height="3" rx="1.5" fill="var(--rule-strong)"/>
              <rect x="8" y="138" width="64" height="3" rx="1.5" fill="var(--rule-strong)"/>
              <rect x="8" y="156" width="50" height="10" rx="5" fill="var(--ink)"/>
            </g>
            <!-- Vision wand -->
            <g transform="translate(160 110)">
              <path d="M0 0l3 7 7 3-7 3-3 7-3-7-7-3 7-3z" fill="var(--accent)"/>
              <path d="M20 -22l1.4 3.2 3.2 1.4-3.2 1.4L20 -13l-1.4-3.2-3.2-1.4 3.2-1.4z" fill="var(--accent)" opacity="0.8"/>
              <path d="M-6 30l.9 2 2 .9-2 .9-.9 2-.9-2-2-.9 2-.9z" fill="var(--accent)" opacity="0.6"/>
            </g>
            <!-- Application card -->
            <g transform="translate(190 50)">
              <rect width="116" height="156" rx="12" fill="white" stroke="var(--accent)" stroke-width="1.5"/>
              <circle cx="22" cy="22" r="10" fill="var(--accent-tint)"/>
              <rect x="38" y="18" width="56" height="6" rx="3" fill="var(--ink)"/>
              <rect x="38" y="30" width="64" height="4" rx="2" fill="var(--mute-2)"/>
              <rect x="14" y="50" width="64" height="14" rx="7" fill="var(--positive-tint)"/>
              <text x="46" y="60" text-anchor="middle" font-family="system-ui" font-size="9" fill="var(--positive-text)" font-weight="600">Screen</text>
              <rect x="14" y="76" width="90" height="3.5" rx="1.75" fill="var(--mute-2)"/>
              <rect x="14" y="84" width="70" height="3.5" rx="1.75" fill="var(--mute-2)"/>
              <rect x="14" y="100" width="90" height="3" rx="1.5" fill="var(--rule)"/>
              <rect x="14" y="108" width="86" height="3" rx="1.5" fill="var(--rule)"/>
              <rect x="14" y="116" width="90" height="3" rx="1.5" fill="var(--rule)"/>
              <rect x="14" y="124" width="60" height="3" rx="1.5" fill="var(--rule)"/>
            </g>
          </svg>
        {:else if step === 3}
          <!-- Brief — central interviewer with orbiting signal cards -->
          <svg viewBox="0 0 320 240" width="320" height="240" fill="none">
            <!-- Connector lines -->
            <path d="M70 78 L 138 118" stroke="var(--rule-strong)" stroke-width="1" stroke-dasharray="2 3"/>
            <path d="M70 178 L 138 130" stroke="var(--rule-strong)" stroke-width="1" stroke-dasharray="2 3"/>
            <path d="M252 80 L 188 116" stroke="var(--rule-strong)" stroke-width="1" stroke-dasharray="2 3"/>
            <path d="M252 178 L 188 132" stroke="var(--rule-strong)" stroke-width="1" stroke-dasharray="2 3"/>
            <!-- Signals -->
            <g transform="translate(14 50)">
              <rect width="60" height="36" rx="8" fill="white" stroke="var(--rule)" stroke-width="1"/>
              <rect x="8" y="8" width="20" height="3" rx="1.5" fill="var(--accent)"/>
              <rect x="8" y="16" width="44" height="3" rx="1.5" fill="var(--ink-2)"/>
              <rect x="8" y="24" width="36" height="3" rx="1.5" fill="var(--ink-2)"/>
            </g>
            <g transform="translate(14 150)">
              <rect width="60" height="36" rx="8" fill="white" stroke="var(--rule)" stroke-width="1"/>
              <rect x="8" y="8" width="20" height="3" rx="1.5" fill="var(--accent)"/>
              <rect x="8" y="16" width="44" height="3" rx="1.5" fill="var(--ink-2)"/>
              <rect x="8" y="24" width="32" height="3" rx="1.5" fill="var(--ink-2)"/>
            </g>
            <g transform="translate(244 52)">
              <rect width="60" height="36" rx="8" fill="white" stroke="var(--rule)" stroke-width="1"/>
              <rect x="8" y="8" width="20" height="3" rx="1.5" fill="var(--accent)"/>
              <rect x="8" y="16" width="44" height="3" rx="1.5" fill="var(--ink-2)"/>
              <rect x="8" y="24" width="36" height="3" rx="1.5" fill="var(--ink-2)"/>
            </g>
            <g transform="translate(244 150)">
              <rect width="60" height="36" rx="8" fill="white" stroke="var(--rule)" stroke-width="1"/>
              <rect x="8" y="8" width="20" height="3" rx="1.5" fill="var(--accent)"/>
              <rect x="8" y="16" width="44" height="3" rx="1.5" fill="var(--ink-2)"/>
              <rect x="8" y="24" width="32" height="3" rx="1.5" fill="var(--ink-2)"/>
            </g>
            <!-- Center interviewer card -->
            <rect x="98" y="64" width="124" height="124" rx="14" fill="white" stroke="var(--accent)" stroke-width="1.6"/>
            <circle cx="160" cy="106" r="22" fill="var(--accent-tint)"/>
            <text x="160" y="113" text-anchor="middle" font-family="system-ui" font-weight="700" font-size="16" fill="var(--accent-text)">SC</text>
            <rect x="120" y="140" width="80" height="5" rx="2.5" fill="var(--ink)"/>
            <rect x="128" y="150" width="64" height="3.5" rx="1.75" fill="var(--mute)"/>
            <rect x="110" y="166" width="44" height="13" rx="6.5" fill="var(--accent-tint)"/>
            <text x="132" y="175" text-anchor="middle" font-family="system-ui" font-size="9" fill="var(--accent-text)" font-weight="600">Talk</text>
            <rect x="166" y="166" width="44" height="13" rx="6.5" fill="var(--positive-tint)"/>
            <text x="188" y="175" text-anchor="middle" font-family="system-ui" font-size="9" fill="var(--positive-text)" font-weight="600">Post</text>
          </svg>
        {:else}
          <!-- Get started — three tilted tiles -->
          <svg viewBox="0 0 320 240" width="320" height="240" fill="none">
            <g transform="translate(14 56) rotate(-3 50 70)">
              <rect width="100" height="140" rx="14" fill="var(--surface)" stroke="var(--rule)" stroke-width="1"/>
              <circle cx="50" cy="50" r="24" fill="var(--accent-tint)"/>
              <path d="M38 46h22M38 50h16M38 54h18" stroke="var(--accent-text)" stroke-width="1.8" stroke-linecap="round"/>
              <rect x="20" y="90" width="60" height="5" rx="2.5" fill="var(--ink)"/>
              <rect x="26" y="104" width="48" height="3.5" rx="1.75" fill="var(--mute)"/>
            </g>
            <g transform="translate(110 38)">
              <rect width="100" height="156" rx="14" fill="var(--accent-tint)" stroke="var(--accent)" stroke-width="2"/>
              <circle cx="50" cy="58" r="26" fill="white"/>
              <path d="M37 54h26M37 58h26M37 62h20" stroke="var(--accent)" stroke-width="1.8" stroke-linecap="round"/>
              <path d="M68 38l1.2 2.6 2.6 1.2-2.6 1.2L68 45.6l-1.2-2.6-2.6-1.2 2.6-1.2z" fill="var(--accent)"/>
              <rect x="20" y="102" width="60" height="5" rx="2.5" fill="var(--accent-text)"/>
              <rect x="26" y="116" width="48" height="3.5" rx="1.75" fill="var(--accent-text)" opacity="0.7"/>
              <text x="50" y="138" text-anchor="middle" font-family="system-ui" font-size="9" font-weight="700" fill="var(--accent-text)" letter-spacing="0.05em">RECOMMENDED</text>
            </g>
            <g transform="translate(206 56) rotate(3 50 70)">
              <rect width="100" height="140" rx="14" fill="var(--surface)" stroke="var(--rule)" stroke-width="1"/>
              <circle cx="50" cy="50" r="24" fill="var(--warm-tint)"/>
              <circle cx="50" cy="50" r="12" stroke="var(--warm-text)" stroke-width="2" fill="none"/>
              <path d="M59 59l7 7" stroke="var(--warm-text)" stroke-width="2" stroke-linecap="round"/>
              <rect x="20" y="90" width="60" height="5" rx="2.5" fill="var(--ink)"/>
              <rect x="26" y="104" width="48" height="3.5" rx="1.75" fill="var(--mute)"/>
            </g>
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
      </div>

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
    </div>
  </div>
</div>

<p class="back-link"><a href="/preview/redesign">← back to previews</a></p>

<style>
  :global(html, body) { background: var(--surface); margin: 0; font-family: var(--sans); color: var(--ink); }
  .ob-overlay {
    min-height: 100vh;
    background:
      radial-gradient(circle at 20% 20%, var(--accent-tint), transparent 40%),
      radial-gradient(circle at 80% 80%, var(--warm-tint), transparent 40%),
      var(--surface);
    display: grid; place-items: center;
    padding: 32px;
  }
  /* Stack with peek-cards behind */
  .stack {
    position: relative;
    width: 100%;
    max-width: 600px;
  }
  .peek {
    position: absolute;
    inset: 0;
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 20px;
    box-shadow: var(--sh-1);
  }
  .peek-1 { transform: translate(8px, 8px) rotate(1.2deg); opacity: 0.85; }
  .peek-2 { transform: translate(16px, 16px) rotate(2.4deg); opacity: 0.55; }

  .ob-card {
    position: relative;
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 20px;
    padding: 30px 32px 22px;
    box-shadow: 0 24px 80px -32px rgba(10,10,13,0.18), var(--sh-1);
  }

  .counter {
    position: absolute;
    top: 18px; right: 22px;
    font-size: 11px; color: var(--mute);
    font-family: var(--mono);
    letter-spacing: 0.06em;
  }

  /* Illustration */
  .art {
    display: grid; place-items: center;
    height: 240px;
    margin-bottom: 18px;
  }

  /* Text */
  .text h1 {
    font-size: 24px; font-weight: 600;
    letter-spacing: -0.022em;
    margin: 0 0 8px;
    text-align: center;
  }
  .text > p {
    font-size: 14px; color: var(--mute);
    line-height: 1.55;
    margin: 0 auto 12px;
    text-align: center;
    max-width: 46ch;
  }
  .text kbd {
    font-family: var(--mono); font-size: 11px;
    background: var(--surface-2); border: 1px solid var(--rule);
    border-bottom-width: 2px; border-radius: 3px;
    padding: 0 4px; color: var(--ink-2);
  }
  .picker-help { margin-bottom: 16px !important; }

  .picker { display: flex; flex-direction: column; gap: 8px; margin: 4px 0 8px; }
  .pick {
    position: relative;
    text-align: left;
    background: var(--surface);
    border: 1px solid var(--rule);
    border-radius: 12px;
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
    display: flex; align-items: center; gap: 8px;
    margin-top: 18px;
    border-top: 1px solid var(--rule);
    padding-top: 14px;
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

  .back-link { text-align: center; margin: 10px 0 30px; font-size: 12.5px; }
  .back-link a { color: var(--accent-text); text-decoration: none; }
</style>
