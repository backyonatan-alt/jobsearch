<script>
  // Static mockup — Layout B: two-pane application detail.
  // AI prep is the dominant left column (hero); meta rail on the right.
  // No backend wiring; all data is hard-coded. The person here is the
  // HIRING MANAGER — no interview is scheduled yet, stated explicitly.
  const app = {
    co: 'Astelia',
    role: 'Senior Product Manager',
    status: 'Applied',
    appliedLong: 'June 3 2026',
    source: 'LinkedIn',
    location: 'Remote (US)',
    salary: '$190k–$230k',
    lastActivity: 'today'
  };

  const company = {
    about: 'Astelia builds infrastructure for climate-tech data. Series B, ~120 people, founded 2021. Research-heavy product culture — they value crisp problem framing over feature lists.',
    facts: [
      { lbl: 'Stage', val: 'Series B' },
      { lbl: 'People', val: '~120' },
      { lbl: 'Founded', val: '2021' },
      { lbl: 'Setup', val: 'Remote-first' },
      { lbl: 'Typical process', val: 'Recruiter → HM → panel' }
    ]
  };

  const hm = {
    name: 'Sarah Chen',
    role: 'VP Product',
    initials: 'SC',
    snapshot: 'Led product at two infra startups; writes about discovery-driven roadmaps. Probes how you de-risk decisions, not just outcomes.',
    tells: [
      { lbl: 'Format', val: 'Conversational · no slides' },
      { lbl: 'Tempo', val: 'Digs into one decision deeply' },
      { lbl: 'Length', val: '~45 min' }
    ]
  };

  const lands = [
    'Crisp problem statements before solutions',
    'Metrics you actually moved + how you measured',
    'A clear stance on prioritization tradeoffs'
  ];
  const avoid = [
    'Feature laundry lists',
    "Vague 'we aligned stakeholders' with no specifics"
  ];

  const signals = [
    { date: 'Apr 2026', kind: 'Talk', body: "'Discovery debt' talk at ProductCon — teams ship too early.", source: 'productcon.com' },
    { date: 'Feb 2026', kind: 'Post', body: 'Essay on killing features as a PM skill.', source: 'sarahchen.substack.com' }
  ];

  const questions = [
    { q: 'How do you decide what NOT to build this quarter?', why: "Mirrors her 'discovery debt' thesis." },
    { q: 'What does great PM craft look like here in 6 months?', why: 'Outcomes + ramp thinking.' }
  ];

  const timeline = [
    { kind: 'accent', label: 'Application submitted', note: 'via LinkedIn' },
    { kind: 'flat', label: 'Status', note: 'Applied' }
  ];
</script>

<svelte:head><title>Detail B — Pursuit</title></svelte:head>

<main class="wrap">
  <!-- Header bar -->
  <header class="hdr">
    <div class="hdr-left">
      <div class="mark"></div>
      <div class="hdr-titles">
        <div class="co-row">
          <h1>{app.co}</h1>
          <span class="pill applied"><span class="pdot"></span>{app.status}</span>
        </div>
        <div class="role-line">{app.role}</div>
        <div class="meta-line">
          <span>Applied {app.appliedLong}</span>
          <span class="dot">·</span>
          <span>{app.source}</span>
        </div>
      </div>
    </div>
    <span class="mock-badge">Static mockup</span>
  </header>

  <!-- Two-pane grid -->
  <div class="grid">
    <!-- LEFT: AI prep (hero) -->
    <div class="left">
      <div class="prep-lead">
        <span class="ai-pill"><span class="spark">✦</span> AI</span>
        <h2>Interview prep</h2>
        <p class="lead-sub">No interview booked yet — prep's ready for when {app.co} gets back to you.</p>
      </div>

      <!-- Company brief -->
      <section class="card">
        <div class="card-hd"><h3>Company brief</h3><span class="ai-tag">AI · refreshed today</span></div>
        <p class="about">{company.about}</p>
        <div class="facts-grid">
          {#each company.facts as f}
            <div class="f-cell"><div class="f-lbl">{f.lbl}</div><div class="f-val">{f.val}</div></div>
          {/each}
        </div>
      </section>

      <!-- Hiring manager -->
      <section class="card">
        <div class="card-hd"><h3>Hiring manager</h3><span class="ai-tag">AI · refreshed today</span></div>
        <div class="person">
          <div class="p-av">{hm.initials}</div>
          <div class="p-info">
            <div class="p-name-row">
              <h4>{hm.name}</h4>
              <span class="role-tag">Hiring manager</span>
            </div>
            <div class="p-role">{hm.role} · {app.co}</div>
          </div>
        </div>
        <p class="snapshot">{hm.snapshot}</p>
        <div class="tells-block">
          <div class="tells-hd">How she interviews</div>
          <div class="tells">
            {#each hm.tells as t}
              <div class="tell"><div class="t-lbl">{t.lbl}</div><div class="t-val">{t.val}</div></div>
            {/each}
          </div>
        </div>
      </section>

      <!-- Lands / avoid -->
      <section class="la-grid">
        <div class="la-col lands">
          <h3><span class="glyph">✓</span> What lands</h3>
          <ul>
            {#each lands as l}<li><span class="g">✓</span><span>{l}</span></li>{/each}
          </ul>
        </div>
        <div class="la-col avoid">
          <h3><span class="glyph">✕</span> What to avoid</h3>
          <ul>
            {#each avoid as a}<li><span class="g">✕</span><span>{a}</span></li>{/each}
          </ul>
        </div>
      </section>

      <!-- Recent signals -->
      <section class="card">
        <div class="card-hd"><h3>Recent signals</h3><span class="ai-tag">AI · web search</span></div>
        <ul class="signals">
          {#each signals as s}
            <li>
              <span class="s-date">{s.date}</span>
              <span class="s-body">
                <span class="s-kind">{s.kind}</span>{s.body}
                <a class="s-source" href="#">{s.source}</a>
              </span>
            </li>
          {/each}
        </ul>
      </section>

      <!-- Questions -->
      <section class="card">
        <div class="card-hd"><h3>Questions worth asking</h3></div>
        <ol class="questions">
          {#each questions as item}
            <li>
              <div class="q">{item.q}</div>
              <div class="why">{item.why}</div>
            </li>
          {/each}
        </ol>
      </section>
    </div>

    <!-- RIGHT: meta rail -->
    <aside class="rail">
      <div class="status-card">
        <div class="sc-top"><span class="pill applied"><span class="pdot"></span>Applied</span></div>
        <p class="sc-line">No interview scheduled. Log one when it's booked.</p>
        <button class="sc-btn">+ Log an interview</button>
      </div>

      <div class="rail-card">
        <div class="rc-hd">Details</div>
        <dl class="kv">
          <dt>Status</dt><dd>{app.status}</dd>
          <dt>Source</dt><dd>{app.source}</dd>
          <dt>Location</dt><dd>{app.location}</dd>
          <dt>Salary</dt><dd>{app.salary}</dd>
          <dt>Last activity</dt><dd>{app.lastActivity}</dd>
        </dl>
      </div>

      <div class="rail-card">
        <div class="rc-hd">Contact</div>
        <div class="contact">
          <div class="c-av">{hm.initials}</div>
          <div class="c-info">
            <div class="c-name">{hm.name}</div>
            <div class="c-role">{hm.role}</div>
          </div>
        </div>
        <a class="c-li" href="#">
          <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M3.5 6h2v6h-2zM4.5 3a1 1 0 1 1 0 2 1 1 0 0 1 0-2zM7 6h2v.9c.3-.5.9-1 1.8-1 1.6 0 2.2 1 2.2 2.6V12h-2V9c0-.9-.3-1.4-1.1-1.4-.6 0-1 .4-1 1.2V12H7z"/></svg>
          LinkedIn
        </a>
      </div>

      <div class="rail-card">
        <div class="rc-hd">Activity</div>
        <div class="timeline">
          {#each timeline as e}
            <div class="tl-event {e.kind}">
              <div class="axis"><span class="marker"></span></div>
              <div class="what"><div class="label">{e.label}</div><div class="note">{e.note}</div></div>
            </div>
          {/each}
        </div>
      </div>
    </aside>
  </div>

  <p class="footer-link"><a href="/preview/redesign">← back to previews</a></p>
</main>

<style>
  :global(html, body) { background: var(--surface); margin: 0; font-family: var(--sans); color: var(--ink); }
  .wrap { max-width: 1180px; margin: 0 auto; padding: 0 28px 64px; }

  /* HEADER BAR */
  .hdr {
    display: flex; align-items: flex-start; justify-content: space-between;
    gap: 24px; padding: 24px 0 20px; border-bottom: 1px solid var(--rule); margin-bottom: 28px;
  }
  .hdr-left { display: grid; grid-template-columns: 48px 1fr; gap: 16px; align-items: start; }
  .mark { width: 48px; height: 48px; border-radius: 12px; background: var(--ink); position: relative; }
  .mark::after { content: ""; position: absolute; inset: 11px 11px auto auto; width: 9px; height: 9px; border-radius: 50%; background: var(--accent); }
  .co-row { display: flex; align-items: center; gap: 12px; }
  .co-row h1 { font-size: 24px; font-weight: 600; margin: 0; letter-spacing: -0.025em; }
  .role-line { font-size: 15px; color: var(--ink-2); margin-top: 3px; font-weight: 500; }
  .meta-line { font-size: 13px; color: var(--mute); margin-top: 6px; font-variant-numeric: tabular-nums; }
  .meta-line .dot { margin: 0 8px; opacity: 0.5; }

  .mock-badge {
    font-size: 11px; font-weight: 500; color: var(--mute);
    background: var(--surface-2); border: 1px solid var(--rule);
    border-radius: 999px; padding: 3px 10px; white-space: nowrap; flex-shrink: 0;
  }

  /* PILL */
  .pill { display: inline-flex; align-items: center; gap: 6px; padding: 3px 10px; border-radius: 999px; font-size: 12px; font-weight: 500; background: var(--surface-2); color: var(--ink-2); white-space: nowrap; }
  .pill .pdot { width: 6px; height: 6px; border-radius: 50%; background: var(--mute-2); }
  .pill.applied { background: var(--surface-2); color: var(--ink-2); }

  /* GRID */
  .grid { display: grid; grid-template-columns: 1fr 320px; gap: 24px; align-items: start; }

  /* LEFT — prep hero */
  .left { display: flex; flex-direction: column; gap: 16px; min-width: 0; }
  .prep-lead { margin-bottom: 2px; }
  .ai-pill {
    display: inline-flex; align-items: center; gap: 6px; font-size: 12px; font-weight: 600;
    color: var(--accent-text); background: var(--accent-tint); border-radius: 999px;
    padding: 4px 11px; letter-spacing: 0.01em;
  }
  .ai-pill .spark { font-size: 12px; }
  .prep-lead h2 { font-size: 26px; font-weight: 600; letter-spacing: -0.028em; margin: 12px 0 6px; }
  .lead-sub { font-size: 15px; color: var(--mute); margin: 0; line-height: 1.5; max-width: 60ch; }

  /* CARD */
  .card { background: var(--card); border: 1px solid var(--rule); border-radius: 14px; padding: 20px 22px; box-shadow: var(--sh-1); }
  .card-hd { display: flex; align-items: center; gap: 10px; margin-bottom: 14px; }
  .card-hd h3 { font-size: 16px; font-weight: 600; margin: 0; letter-spacing: -0.015em; }
  .ai-tag { display: inline-flex; align-items: center; gap: 5px; font-size: 11.5px; background: var(--accent-tint); color: var(--accent-text); padding: 3px 9px; border-radius: 999px; font-weight: 500; margin-left: auto; }

  /* COMPANY */
  .about { font-size: 14.5px; line-height: 1.6; color: var(--ink-2); margin: 0 0 16px; }
  .facts-grid { display: grid; grid-template-columns: repeat(5, 1fr); gap: 8px; }
  .f-cell { background: var(--surface-2); border-radius: 10px; padding: 10px 12px; }
  .f-lbl { font-size: 11px; color: var(--mute); font-weight: 500; }
  .f-val { font-size: 12.5px; color: var(--ink); font-weight: 600; margin-top: 3px; line-height: 1.35; }

  /* PERSON */
  .person { display: grid; grid-template-columns: 52px 1fr; gap: 14px; align-items: center; margin-bottom: 14px; }
  .p-av { width: 52px; height: 52px; border-radius: 50%; display: grid; place-items: center; font-weight: 600; font-size: 18px; background: var(--accent-tint); color: var(--accent-text); }
  .p-name-row { display: flex; align-items: center; gap: 10px; }
  .p-info h4 { margin: 0; font-size: 17px; font-weight: 600; letter-spacing: -0.015em; }
  .role-tag { font-size: 10.5px; font-weight: 600; letter-spacing: 0.04em; text-transform: uppercase; color: var(--warm-text); background: var(--warm-tint); padding: 2px 8px; border-radius: 5px; }
  .p-role { font-size: 13px; color: var(--mute); margin-top: 3px; }
  .snapshot { font-size: 15px; line-height: 1.55; letter-spacing: -0.008em; color: var(--ink); margin: 0 0 18px; padding-left: 14px; border-left: 2px solid var(--accent); }

  .tells-block { border-top: 1px solid var(--rule); padding-top: 16px; }
  .tells-hd { font-size: 12.5px; font-weight: 600; color: var(--mute); margin-bottom: 12px; }
  .tells { display: grid; grid-template-columns: repeat(3, 1fr); gap: 14px; }
  .t-lbl { font-size: 11.5px; font-weight: 500; color: var(--mute); margin-bottom: 4px; }
  .t-val { font-size: 13.5px; color: var(--ink); line-height: 1.4; }

  /* LANDS / AVOID */
  .la-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 0; border: 1px solid var(--rule); border-radius: 14px; background: var(--card); overflow: hidden; box-shadow: var(--sh-1); }
  .la-col { padding: 20px 22px; }
  .la-col + .la-col { border-left: 1px solid var(--rule); }
  .la-col h3 { font-size: 13px; font-weight: 600; margin: 0 0 14px; display: flex; align-items: center; gap: 8px; }
  .la-col.lands h3 { color: var(--positive-text); }
  .la-col.avoid h3 { color: var(--danger-text); }
  .la-col h3 .glyph { width: 18px; height: 18px; border-radius: 5px; display: inline-flex; align-items: center; justify-content: center; font-size: 11px; font-weight: 600; }
  .la-col.lands h3 .glyph { background: var(--positive-tint); color: var(--positive-text); }
  .la-col.avoid h3 .glyph { background: var(--danger-tint); color: var(--danger-text); }
  .la-col ul { list-style: none; margin: 0; padding: 0; }
  .la-col li { font-size: 13.5px; line-height: 1.5; color: var(--ink-2); padding: 9px 0; border-top: 1px solid var(--rule); display: grid; grid-template-columns: 16px 1fr; gap: 8px; align-items: start; }
  .la-col li:first-child { border-top: none; padding-top: 0; }
  .la-col li .g { font-size: 11px; margin-top: 2px; }
  .la-col.lands li .g { color: var(--positive-text); }
  .la-col.avoid li .g { color: var(--danger-text); }

  /* SIGNALS */
  .signals { list-style: none; padding: 0; margin: 0; }
  .signals li { display: grid; grid-template-columns: 72px 1fr; gap: 14px; padding: 12px 0; border-top: 1px solid var(--rule); font-size: 13.5px; }
  .signals li:first-child { border-top: none; padding-top: 0; }
  .s-date { font-size: 12px; font-weight: 500; color: var(--mute); padding-top: 2px; }
  .s-body { color: var(--ink-2); line-height: 1.5; }
  .s-kind { display: inline-block; font-size: 11px; font-weight: 500; color: var(--mute); margin-right: 8px; padding: 1px 7px; background: var(--surface-2); border-radius: 4px; vertical-align: 1px; }
  .s-source { font-size: 12px; color: var(--accent-text); margin-left: 6px; text-decoration: none; }

  /* QUESTIONS */
  .questions { margin: 0; padding: 0; list-style: none; counter-reset: q; }
  .questions li { counter-increment: q; display: grid; grid-template-columns: 28px 1fr; gap: 12px; padding: 13px 0; border-top: 1px solid var(--rule); align-items: start; }
  .questions li:first-child { border-top: none; padding-top: 0; }
  .questions li::before { content: counter(q, decimal-leading-zero); font-size: 12px; font-weight: 500; color: var(--mute-2); font-variant-numeric: tabular-nums; padding-top: 2px; }
  .questions .q { font-size: 15px; line-height: 1.45; color: var(--ink); letter-spacing: -0.005em; }
  .questions .why { font-size: 12.5px; color: var(--mute); margin-top: 4px; line-height: 1.45; }

  /* RIGHT RAIL */
  .rail { display: flex; flex-direction: column; gap: 14px; position: sticky; top: 20px; }
  .status-card { background: var(--surface-2); border: 1px solid var(--rule); border-radius: 12px; padding: 16px; }
  .sc-top { margin-bottom: 10px; }
  .sc-line { font-size: 13px; color: var(--ink-2); line-height: 1.5; margin: 0 0 14px; }
  .sc-btn { width: 100%; height: 32px; border-radius: 8px; background: var(--ink); color: white; font-size: 13px; font-weight: 500; cursor: pointer; }

  .rail-card { background: var(--card); border: 1px solid var(--rule); border-radius: 12px; padding: 16px; box-shadow: var(--sh-1); }
  .rc-hd { font-size: 12.5px; font-weight: 600; color: var(--mute); margin-bottom: 12px; }

  .kv { margin: 0; display: grid; grid-template-columns: auto 1fr; gap: 9px 12px; }
  .kv dt { font-size: 12.5px; color: var(--mute); }
  .kv dd { margin: 0; font-size: 12.5px; color: var(--ink); font-weight: 500; text-align: right; font-variant-numeric: tabular-nums; }

  .contact { display: grid; grid-template-columns: 36px 1fr; gap: 10px; align-items: center; margin-bottom: 12px; }
  .c-av { width: 36px; height: 36px; border-radius: 50%; display: grid; place-items: center; font-weight: 600; font-size: 13px; background: var(--accent-tint); color: var(--accent-text); }
  .c-name { font-size: 13.5px; font-weight: 600; }
  .c-role { font-size: 12px; color: var(--mute); margin-top: 1px; }
  .c-li { display: inline-flex; align-items: center; gap: 6px; width: 100%; justify-content: center; background: var(--surface-2); border: 1px solid var(--rule); border-radius: 8px; padding: 7px 12px; font-size: 12.5px; font-weight: 500; color: var(--ink); text-decoration: none; }
  .c-li svg { color: #0a66c2; }

  /* TIMELINE */
  .timeline { padding: 2px 0; }
  .tl-event { display: grid; grid-template-columns: 16px 1fr; gap: 12px; padding: 10px 0; align-items: start; }
  .tl-event:first-child { padding-top: 0; }
  .tl-event .axis { position: relative; display: flex; justify-content: center; padding-top: 3px; }
  .tl-event .axis::before { content: ""; position: absolute; top: 14px; bottom: -16px; width: 1px; background: var(--rule); }
  .tl-event:last-child .axis::before { display: none; }
  .tl-event .axis .marker { width: 10px; height: 10px; border-radius: 50%; background: var(--card); border: 2px solid var(--mute-2); z-index: 1; }
  .tl-event.accent .axis .marker { border-color: var(--accent); }
  .tl-event .what .label { font-size: 13px; font-weight: 500; color: var(--ink); margin-bottom: 1px; }
  .tl-event .what .note { font-size: 12px; color: var(--mute); line-height: 1.45; }

  .footer-link { margin-top: 32px; font-size: 13px; }
  .footer-link a { color: var(--accent-text); text-decoration: none; }

  /* MOBILE */
  @media (max-width: 860px) {
    .grid { grid-template-columns: 1fr; }
    .rail { position: static; }
    .facts-grid { grid-template-columns: repeat(2, 1fr); }
    .tells { grid-template-columns: 1fr; }
    .la-grid { grid-template-columns: 1fr; }
    .la-col + .la-col { border-left: none; border-top: 1px solid var(--rule); }
  }
</style>
