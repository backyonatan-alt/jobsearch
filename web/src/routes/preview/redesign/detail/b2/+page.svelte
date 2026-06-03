<script>
  // Static mockup — Layout B2: refinement of Option B's two-pane detail.
  // AI prep dominates the LEFT column; meta in the RIGHT rail.
  // This variant shows a POPULATED next-interview card and adds a baby-blue
  // AI tips box, a redesigned full-width Activity log, and A's question cards.
  // All data hard-coded; no backend, no fetching.
  const app = {
    co: 'Astelia',
    role: 'Senior Product Manager',
    status: 'Interview',
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
      { lbl: 'Typical process', val: 'screen → PM craft → product sense → exec' }
    ]
  };

  const tips = [
    "Map a concrete story to each of the JD's top 3 responsibilities.",
    "Re-read your own application — they'll reference what you wrote.",
    "Have 2 questions ready that show you get Astelia's stage."
  ];

  const nextInterview = {
    title: 'PM craft round · Thu Jun 5, 11:00',
    rows: ['Sarah Chen (Hiring manager)', 'Google Meet', '45 min']
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
    { kind: 'accent', date: 'Jun 3', label: 'Application submitted', note: 'via LinkedIn' },
    { kind: 'warm', date: 'Jun 3', label: 'Status', note: 'Applied → Interview' }
  ];
</script>

<svelte:head><title>Detail B2 — Pursuit</title></svelte:head>

<main class="wrap">
  <!-- Header bar -->
  <header class="hdr">
    <div class="hdr-left">
      <div class="mark"></div>
      <div class="hdr-titles">
        <div class="co-row">
          <h1>{app.co}</h1>
          <span class="pill interview"><span class="pdot"></span>{app.status}</span>
        </div>
        <div class="role-line">{app.role}</div>
        <div class="meta-line">
          <span>Applied {app.appliedLong}</span>
          <span class="dot">·</span>
          <span>{app.source}</span>
          <span class="dot">·</span>
          <span>{app.location}</span>
          <span class="dot">·</span>
          <span>{app.salary}</span>
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
        <p class="lead-sub">Your PM craft round with Sarah is Thursday — here's everything to walk in ready.</p>
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

      <!-- AI tips box (baby blue) -->
      <section class="tips">
        <div class="tips-hd">
          <span class="tips-spark">✦</span>
          <h3>Tips for this one</h3>
          <span class="tips-ai">AI</span>
        </div>
        <ul class="tips-list">
          {#each tips as t}
            <li><span class="tip-dot"></span><span>{t}</span></li>
          {/each}
        </ul>
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
          <div class="tells-hd">How Sarah interviews</div>
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

      <!-- Questions worth asking (Option A card design) -->
      <section class="card">
        <div class="card-hd"><h3>Questions worth asking</h3></div>
        <div class="questions">
          {#each questions as item}
            <div class="prep-q">
              <div class="q">{item.q}</div>
              <div class="why"><span class="why-spark">✦</span>{item.why}</div>
            </div>
          {/each}
        </div>
      </section>

      <!-- Activity log (full width, with add affordances) -->
      <section class="card activity">
        <div class="act-hd">
          <h3>Activity</h3>
          <div class="act-actions">
            <button class="ghost-btn">+ Log a follow-up</button>
            <button class="ghost-btn">+ Add a note</button>
            <button class="ghost-btn">+ Log an event</button>
          </div>
        </div>
        <div class="timeline">
          {#each timeline as e}
            <div class="tl-event {e.kind}">
              <div class="axis"><span class="marker"></span></div>
              <div class="what">
                <div class="tl-row"><span class="tl-date">{e.date}</span><span class="label">{e.label}</span></div>
                <div class="note">{e.note}</div>
              </div>
            </div>
          {/each}
        </div>
      </section>
    </div>

    <!-- RIGHT: meta rail -->
    <aside class="rail">
      <!-- Next interview (populated, prominent) -->
      <div class="next-card">
        <div class="nc-kicker"><span class="nc-dot"></span>NEXT INTERVIEW</div>
        <div class="nc-title">{nextInterview.title}</div>
        <div class="nc-rows">
          {#each nextInterview.rows as r}
            <div class="nc-row">{r}</div>
          {/each}
        </div>
      </div>
      <p class="next-caption">When nothing's booked, this shows “No interview scheduled yet”.</p>

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
            <div class="c-role">{hm.role} · Hiring manager</div>
          </div>
        </div>
        <a class="c-li" href="#">
          <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M3.5 6h2v6h-2zM4.5 3a1 1 0 1 1 0 2 1 1 0 0 1 0-2zM7 6h2v.9c.3-.5.9-1 1.8-1 1.6 0 2.2 1 2.2 2.6V12h-2V9c0-.9-.3-1.4-1.1-1.4-.6 0-1 .4-1 1.2V12H7z"/></svg>
          LinkedIn
        </a>
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
    font-family: var(--mono); font-size: 11px; letter-spacing: 0.02em; color: var(--mute);
    background: var(--surface-2); border: 1px solid var(--rule);
    border-radius: 999px; padding: 3px 10px; white-space: nowrap; flex-shrink: 0;
  }

  /* PILL */
  .pill { display: inline-flex; align-items: center; gap: 6px; padding: 3px 10px; border-radius: 999px; font-size: 12px; font-weight: 500; background: var(--surface-2); color: var(--ink-2); white-space: nowrap; }
  .pill .pdot { width: 6px; height: 6px; border-radius: 50%; background: var(--mute-2); }
  .pill.interview { background: var(--warm-tint); color: var(--warm-text); }
  .pill.interview .pdot { background: var(--warm); }

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

  /* AI TIPS BOX (baby blue) */
  .tips { background: var(--accent-tint); border: 1px solid var(--accent-tint-2); border-radius: 13px; padding: 18px 20px; }
  .tips-hd { display: flex; align-items: center; gap: 8px; margin-bottom: 12px; }
  .tips-spark { font-size: 14px; color: var(--accent-text); }
  .tips-hd h3 { font-size: 15px; font-weight: 600; margin: 0; color: var(--accent-text); letter-spacing: -0.01em; }
  .tips-ai { margin-left: auto; font-size: 10.5px; font-weight: 600; letter-spacing: 0.04em; color: var(--accent-text); background: var(--card); border: 1px solid var(--accent-tint-2); padding: 2px 8px; border-radius: 999px; }
  .tips-list { list-style: none; margin: 0; padding: 0; display: flex; flex-direction: column; gap: 9px; }
  .tips-list li { display: grid; grid-template-columns: 14px 1fr; gap: 9px; align-items: start; font-size: 13.5px; line-height: 1.5; color: var(--accent-text); }
  .tip-dot { width: 5px; height: 5px; border-radius: 50%; background: var(--accent); margin-top: 8px; }

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

  /* QUESTIONS — Option A card design */
  .questions { display: flex; flex-direction: column; gap: 12px; }
  .prep-q { background: var(--surface); border: 1px solid var(--rule); border-radius: 11px; padding: 13px 15px; }
  .prep-q .q { font-size: 14px; font-weight: 500; color: var(--ink); line-height: 1.45; }
  .prep-q .why { font-size: 12.5px; color: var(--mute); margin-top: 5px; display: flex; align-items: baseline; gap: 6px; }
  .why-spark { color: var(--accent-text); font-size: 11px; }

  /* ACTIVITY (full width) */
  .activity { padding: 20px 22px; }
  .act-hd { display: flex; align-items: center; justify-content: space-between; gap: 14px; flex-wrap: wrap; margin-bottom: 16px; }
  .act-hd h3 { font-size: 16px; font-weight: 600; margin: 0; letter-spacing: -0.015em; }
  .act-actions { display: flex; gap: 8px; flex-wrap: wrap; }
  .ghost-btn { font-size: 12.5px; font-weight: 500; color: var(--ink-2); background: var(--surface-2); border: 1px solid var(--rule); border-radius: 8px; padding: 6px 11px; cursor: pointer; }
  .ghost-btn:hover { border-color: var(--rule-strong); color: var(--ink); }

  /* TIMELINE */
  .timeline { padding: 2px 0; }
  .tl-event { display: grid; grid-template-columns: 16px 1fr; gap: 14px; padding: 11px 0; align-items: start; }
  .tl-event:first-child { padding-top: 0; }
  .tl-event .axis { position: relative; display: flex; justify-content: center; padding-top: 4px; }
  .tl-event .axis::before { content: ""; position: absolute; top: 15px; bottom: -17px; width: 1px; background: var(--rule); }
  .tl-event:last-child .axis::before { display: none; }
  .tl-event .axis .marker { width: 10px; height: 10px; border-radius: 50%; background: var(--card); border: 2px solid var(--mute-2); z-index: 1; }
  .tl-event.accent .axis .marker { border-color: var(--accent); }
  .tl-event.warm .axis .marker { border-color: var(--warm); }
  .tl-event .tl-row { display: flex; align-items: baseline; gap: 10px; }
  .tl-event .tl-date { font-family: var(--mono); font-size: 11.5px; color: var(--mute); min-width: 40px; }
  .tl-event .label { font-size: 13.5px; font-weight: 600; color: var(--ink); }
  .tl-event .note { font-size: 12.5px; color: var(--mute); line-height: 1.45; margin-top: 2px; margin-left: 50px; }

  /* RIGHT RAIL */
  .rail { display: flex; flex-direction: column; gap: 14px; position: sticky; top: 20px; }

  /* NEXT INTERVIEW — prominent solid ink card */
  .next-card { background: var(--ink); border-radius: 14px; padding: 18px; box-shadow: var(--sh-1); }
  .nc-kicker { display: flex; align-items: center; gap: 7px; font-size: 11px; font-weight: 600; letter-spacing: 0.06em; color: rgba(255,255,255,0.7); }
  .nc-dot { width: 7px; height: 7px; border-radius: 50%; background: var(--warm); box-shadow: 0 0 0 3px rgba(255,255,255,0.12); }
  .nc-title { font-size: 16px; font-weight: 600; color: #fff; line-height: 1.35; letter-spacing: -0.01em; margin: 12px 0 14px; }
  .nc-rows { display: flex; flex-direction: column; gap: 7px; border-top: 1px solid rgba(255,255,255,0.14); padding-top: 12px; }
  .nc-row { font-size: 13px; color: rgba(255,255,255,0.85); }
  .next-caption { font-size: 11.5px; color: var(--mute); line-height: 1.45; margin: -4px 2px 2px; }

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
