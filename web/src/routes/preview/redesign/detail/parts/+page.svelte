<script>
  // Static mockup — parts to pick for the detail page.
  // Three "Company brief" treatments + one redesigned Activity log.
  // All data hard-coded; no backend, no fetching.
  const company = {
    name: 'Astelia',
    headline: 'Astelia builds infrastructure for climate-tech data.',
    about: 'Series B climate-data infra company, ~120 people, founded 2021. Research-heavy product culture that values crisp problem framing over feature lists.',
    facts: [
      { lbl: 'Stage', val: 'Series B' },
      { lbl: 'Size', val: '~120 people' },
      { lbl: 'Founded', val: '2021' },
      { lbl: 'HQ', val: 'Remote-first' }
    ],
    process: ['Recruiter screen', 'PM craft', 'Product sense', 'Exec/values'],
    values: [
      'Research-heavy discovery',
      'Crisp problem framing over feature lists'
    ]
  };

  const timeline = [
    { kind: 'positive', date: 'Jun 3', label: 'Application submitted', note: 'via LinkedIn' },
    { kind: 'accent', date: 'Jun 3', label: 'Status', note: 'Applied → Interview' },
    { kind: 'warm', date: 'Jun 5', label: 'Follow-up', note: 'Emailed Sarah to check in' }
  ];
</script>

<svelte:head><title>Detail parts — Pursuit</title></svelte:head>

<main class="wrap">
  <header class="pg-hd">
    <h1>Detail page — parts to pick</h1>
    <p class="footer-link"><a href="/preview/redesign">← back to previews</a></p>
  </header>

  <!-- ============ PART 1 ============ -->
  <section class="part">
    <div class="part-hd"><span class="part-no">Part 1</span><h2>Company brief — 3 options</h2></div>

    <!-- OPTION 1 — Editorial -->
    <div class="opt-label">Option 1 — Editorial</div>
    <section class="card brief-1">
      <div class="b1-hd">
        <div class="logo">A</div>
        <div class="b1-titles">
          <div class="b1-name-row">
            <h3>{company.name}</h3>
            <span class="ai-pill"><span class="spark">✦</span> AI</span>
          </div>
          <div class="b1-headline">{company.headline}</div>
        </div>
      </div>
      <p class="about">{company.about}</p>
      <div class="b1-facts">
        {#each company.facts as f, i}
          <div class="b1-cell">
            <div class="f-lbl">{f.lbl}</div>
            <div class="f-val">{f.val}</div>
          </div>
        {/each}
      </div>
      <div class="b1-process">
        <div class="proc-lbl">Interview process</div>
        <div class="proc-chips">
          {#each company.process as step, i}
            <span class="chip">{step}</span>{#if i < company.process.length - 1}<span class="proc-arrow">→</span>{/if}
          {/each}
        </div>
      </div>
    </section>

    <!-- OPTION 2 — Stat-forward -->
    <div class="opt-label">Option 2 — Stat-forward</div>
    <section class="card brief-2">
      <div class="b2-topbar">
        <span class="b2-co">{company.name}</span>
        <span class="ai-pill"><span class="spark">✦</span> AI</span>
      </div>
      <div class="b2-stats">
        {#each company.facts as f}
          <div class="b2-stat">
            <div class="b2-val">{f.val}</div>
            <div class="b2-lbl">{f.lbl}</div>
          </div>
        {/each}
      </div>
      <p class="about">{company.about}</p>
      <div class="b2-process">
        <div class="proc-lbl">Interview process</div>
        <ol class="stepper">
          {#each company.process as step, i}
            <li><span class="step-no">{i + 1}</span><span class="step-txt">{step}</span></li>
          {/each}
        </ol>
      </div>
    </section>

    <!-- OPTION 3 — Narrative + sidebar -->
    <div class="opt-label">Option 3 — Narrative + sidebar facts</div>
    <section class="card brief-3">
      <div class="b3-hd">
        <h3>{company.name}</h3>
        <span class="ai-pill"><span class="spark">✦</span> AI</span>
      </div>
      <div class="b3-cols">
        <div class="b3-main">
          <p class="about">{company.about}</p>
          <div class="b3-values">
            <div class="proc-lbl">What they value</div>
            <ul class="value-list">
              {#each company.values as v}
                <li><span class="v-dot"></span><span>{v}</span></li>
              {/each}
            </ul>
          </div>
        </div>
        <aside class="b3-side">
          <dl class="b3-facts">
            {#each company.facts as f}
              <dt>{f.lbl}</dt><dd>{f.val}</dd>
            {/each}
          </dl>
          <div class="b3-proc">
            <div class="proc-lbl">Process</div>
            <ol class="b3-steps">
              {#each company.process as step, i}
                <li><span class="step-no">{i + 1}</span>{step}</li>
              {/each}
            </ol>
          </div>
        </aside>
      </div>
    </section>
  </section>

  <!-- ============ PART 2 ============ -->
  <section class="part">
    <div class="part-hd"><span class="part-no">Part 2</span><h2>Activity — redesigned</h2></div>

    <section class="card activity">
      <div class="act-hd">
        <h3>Activity</h3>
        <div class="act-actions">
          <button class="ghost-btn">+ Log a follow-up</button>
          <button class="ghost-btn">+ Add a note</button>
          <button class="ghost-btn">+ Log an event</button>
        </div>
      </div>
      <ul class="timeline">
        {#each timeline as e}
          <li class="tl-row {e.kind}">
            <span class="tl-rail"><span class="tl-dot"></span></span>
            <span class="tl-date">{e.date}</span>
            <span class="tl-body">
              <span class="tl-title">{e.label}</span>
              <span class="tl-note">{e.note}</span>
            </span>
          </li>
        {/each}
      </ul>
    </section>
  </section>
</main>

<style>
  :global(html, body) { background: var(--surface); margin: 0; font-family: var(--sans); color: var(--ink); }
  .wrap { max-width: 760px; margin: 0 auto; padding: 0 24px 80px; }

  /* PAGE HEADER */
  .pg-hd { padding: 32px 0 24px; border-bottom: 1px solid var(--rule); margin-bottom: 32px; }
  .pg-hd h1 { font-size: 24px; font-weight: 600; letter-spacing: -0.025em; margin: 0; }
  .footer-link { margin: 10px 0 0; font-size: 13px; }
  .footer-link a { color: var(--accent-text); text-decoration: none; }

  /* PART SCAFFOLD */
  .part { margin-bottom: 56px; }
  .part-hd { display: flex; align-items: baseline; gap: 12px; margin-bottom: 20px; }
  .part-no { font-family: var(--mono); font-size: 11px; letter-spacing: 0.04em; text-transform: uppercase; color: var(--mute); background: var(--surface-2); border: 1px solid var(--rule); border-radius: 999px; padding: 3px 10px; }
  .part-hd h2 { font-size: 19px; font-weight: 600; letter-spacing: -0.02em; margin: 0; }
  .opt-label { font-size: 12px; font-weight: 600; color: var(--mute); letter-spacing: 0.01em; margin: 24px 0 10px; }
  .opt-label:first-of-type { margin-top: 0; }

  /* SHARED */
  .card { background: var(--card); border: 1px solid var(--rule); border-radius: 14px; padding: 22px 24px; box-shadow: var(--sh-1); }
  .about { font-size: 14.5px; line-height: 1.6; color: var(--ink-2); margin: 0; }
  .ai-pill { display: inline-flex; align-items: center; gap: 5px; font-size: 11.5px; font-weight: 600; color: var(--accent-text); background: var(--accent-tint); border: 1px solid var(--accent-tint-2); border-radius: 999px; padding: 3px 9px; letter-spacing: 0.01em; white-space: nowrap; }
  .ai-pill .spark { font-size: 11px; }
  .proc-lbl { font-size: 11.5px; font-weight: 600; color: var(--mute); letter-spacing: 0.01em; margin-bottom: 10px; }
  h3 { font-size: 17px; font-weight: 600; letter-spacing: -0.02em; margin: 0; }

  /* ---- OPTION 1 — Editorial ---- */
  .b1-hd { display: grid; grid-template-columns: 44px 1fr; gap: 14px; align-items: start; margin-bottom: 16px; }
  .logo { width: 44px; height: 44px; border-radius: 11px; background: var(--ink); color: #fff; display: grid; place-items: center; font-size: 19px; font-weight: 600; }
  .b1-name-row { display: flex; align-items: center; gap: 10px; }
  .b1-headline { font-size: 13.5px; color: var(--mute); margin-top: 4px; line-height: 1.4; }
  .brief-1 .about { margin-bottom: 18px; }
  .b1-facts { display: grid; grid-template-columns: repeat(4, 1fr); border: 1px solid var(--rule); border-radius: 10px; overflow: hidden; margin-bottom: 18px; }
  .b1-cell { padding: 11px 13px; }
  .b1-cell + .b1-cell { border-left: 1px solid var(--rule); }
  .f-lbl { font-size: 11px; color: var(--mute); font-weight: 500; }
  .f-val { font-size: 13px; color: var(--ink); font-weight: 600; margin-top: 3px; line-height: 1.3; }
  .b1-process .proc-chips { display: flex; align-items: center; flex-wrap: wrap; gap: 7px; }
  .chip { font-size: 12px; font-weight: 500; color: var(--ink-2); background: var(--surface-2); border: 1px solid var(--rule); border-radius: 7px; padding: 4px 10px; }
  .proc-arrow { color: var(--mute-2); font-size: 12px; }

  /* ---- OPTION 2 — Stat-forward ---- */
  .b2-topbar { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
  .b2-co { font-size: 17px; font-weight: 600; letter-spacing: -0.02em; }
  .b2-stats { display: grid; grid-template-columns: repeat(4, 1fr); gap: 0; border: 1px solid var(--rule); border-radius: 12px; overflow: hidden; background: var(--surface-2); margin-bottom: 18px; }
  .b2-stat { padding: 16px 14px; background: var(--card); }
  .b2-stat + .b2-stat { border-left: 1px solid var(--rule); }
  .b2-val { font-size: 18px; font-weight: 600; letter-spacing: -0.02em; color: var(--ink); line-height: 1.2; }
  .b2-lbl { font-size: 11.5px; color: var(--mute); font-weight: 500; margin-top: 5px; }
  .brief-2 .about { margin-bottom: 20px; }
  .stepper { list-style: none; margin: 0; padding: 0; display: flex; flex-wrap: wrap; gap: 8px 16px; }
  .stepper li { display: inline-flex; align-items: center; gap: 8px; font-size: 13px; color: var(--ink-2); font-weight: 500; }
  .step-no { width: 19px; height: 19px; border-radius: 50%; background: var(--accent-tint); color: var(--accent-text); font-size: 11px; font-weight: 600; display: grid; place-items: center; flex-shrink: 0; }

  /* ---- OPTION 3 — Narrative + sidebar ---- */
  .b3-hd { display: flex; align-items: center; gap: 10px; margin-bottom: 16px; }
  .b3-cols { display: grid; grid-template-columns: 1fr 200px; gap: 28px; align-items: start; }
  .b3-main .about { margin-bottom: 18px; }
  .value-list { list-style: none; margin: 0; padding: 0; display: flex; flex-direction: column; gap: 8px; }
  .value-list li { display: grid; grid-template-columns: 14px 1fr; gap: 8px; align-items: start; font-size: 13.5px; line-height: 1.5; color: var(--ink-2); }
  .v-dot { width: 5px; height: 5px; border-radius: 50%; background: var(--accent); margin-top: 8px; }
  .b3-side { background: var(--surface-2); border: 1px solid var(--rule); border-radius: 12px; padding: 16px; }
  .b3-facts { margin: 0; display: grid; grid-template-columns: auto 1fr; gap: 9px 12px; }
  .b3-facts dt { font-size: 12px; color: var(--mute); }
  .b3-facts dd { margin: 0; font-size: 12.5px; color: var(--ink); font-weight: 600; text-align: right; }
  .b3-proc { border-top: 1px solid var(--rule); margin-top: 14px; padding-top: 14px; }
  .b3-steps { list-style: none; margin: 0; padding: 0; display: flex; flex-direction: column; gap: 8px; }
  .b3-steps li { display: flex; align-items: center; gap: 9px; font-size: 12.5px; color: var(--ink-2); }

  /* ============ ACTIVITY ============ */
  .activity { padding: 22px 24px; }
  .act-hd { display: flex; align-items: center; justify-content: space-between; gap: 14px; flex-wrap: wrap; margin-bottom: 18px; }
  .act-hd h3 { font-size: 16px; }
  .act-actions { display: flex; gap: 8px; flex-wrap: wrap; }
  .ghost-btn { font-size: 12.5px; font-weight: 500; color: var(--ink-2); background: var(--surface-2); border: 1px solid var(--rule); border-radius: 8px; padding: 6px 11px; cursor: pointer; }
  .ghost-btn:hover { border-color: var(--rule-strong); color: var(--ink); }

  /* TIMELINE — connector lives in a per-row rail so it lines up dot-to-dot */
  .timeline { list-style: none; margin: 0; padding: 0; }
  .tl-row {
    display: grid; grid-template-columns: 16px 50px 1fr; gap: 14px;
    align-items: start; padding: 14px 0;
  }
  /* rail spans the full row height; the line is centered in it */
  .tl-rail { position: relative; align-self: stretch; }
  .tl-rail::before {
    content: ""; position: absolute; left: 50%; transform: translateX(-50%);
    top: 0; bottom: 0; width: 1px; background: var(--rule);
  }
  /* trim the line so it begins/ends exactly at the first/last dot center */
  .tl-row:first-child .tl-rail::before { top: 9px; }
  .tl-row:last-child .tl-rail::before { bottom: auto; height: 9px; }
  .tl-dot {
    position: absolute; left: 50%; top: 9px; transform: translate(-50%, -50%);
    width: 8px; height: 8px; border-radius: 50%;
    background: var(--mute-2); z-index: 1;
  }
  .tl-row.positive .tl-dot { background: var(--positive); box-shadow: 0 0 0 3px var(--positive-tint); }
  .tl-row.accent .tl-dot { background: var(--accent); box-shadow: 0 0 0 3px var(--accent-tint); }
  .tl-row.warm .tl-dot { background: var(--warm); box-shadow: 0 0 0 3px var(--warm-tint); }
  .tl-date { font-size: 12.5px; color: var(--mute); line-height: 18px; font-variant-numeric: tabular-nums; }
  .tl-body { display: flex; flex-direction: column; gap: 2px; min-width: 0; }
  .tl-title { line-height: 18px; }
  .tl-title { font-size: 13.5px; font-weight: 600; color: var(--ink); }
  .tl-note { font-size: 12.5px; color: var(--mute); line-height: 1.45; }

  /* MOBILE */
  @media (max-width: 680px) {
    .b1-facts, .b2-stats { grid-template-columns: repeat(2, 1fr); }
    .b1-cell:nth-child(3) { border-left: none; }
    .b1-cell:nth-child(n+3) { border-top: 1px solid var(--rule); }
    .b2-stat:nth-child(3) { border-left: none; }
    .b2-stat:nth-child(n+3) { border-top: 1px solid var(--rule); }
    .b3-cols { grid-template-columns: 1fr; gap: 20px; }
  }
</style>
