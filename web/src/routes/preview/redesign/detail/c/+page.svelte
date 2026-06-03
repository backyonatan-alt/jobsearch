<script>
  // Static mockup — Direction C: "Spotlight cards + tabbed prep".
  // No backend; sample data hard-coded inline. Tab switching is local $state.
  const app = {
    co: 'Astelia',
    domain: 'astelia.com',
    role: 'Senior Product Manager',
    status: 'Applied',
    appliedFull: 'Applied June 3 2026',
    source: 'LinkedIn',
    location: 'Remote (US)',
    salary: '$190k–$230k'
  };

  const company = {
    oneLiner: 'Infrastructure for climate-tech data.',
    blurb: 'Astelia builds infrastructure for climate-tech data. Series B, ~120 people, founded 2021. Research-heavy product culture — values crisp problem framing over feature lists.',
    facts: ['Series B', '~120 people', 'Founded 2021', 'Remote-first'],
    process: 'screen → PM craft → product sense → exec'
  };

  const person = {
    name: 'Sarah Chen',
    initials: 'SC',
    title: 'VP Product',
    who: 'Led product at two infra startups; writes about discovery-driven roadmaps.',
    snapshotLead: 'Sarah probes how you de-risk decisions, not just outcomes. She led product at two infra startups and writes about discovery-driven roadmaps.',
    tells: [
      { k: 'Format', v: 'Conversational, no slides' },
      { k: 'Tempo', v: 'Digs into one decision in depth' },
      { k: 'Length', v: '~45 min' }
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
    { q: 'What does great PM craft look like here in 6 months?', why: 'Outcomes + ramp.' }
  ];

  const activity = [
    { title: 'Application submitted', detail: 'via LinkedIn', date: 'Jun 3' },
    { title: 'Status', detail: 'Applied', date: 'Jun 3' }
  ];

  const details = [
    { k: 'Status', v: 'Applied' },
    { k: 'Source', v: 'LinkedIn' },
    { k: 'Location', v: 'Remote (US)' },
    { k: 'Salary', v: '$190k–$230k' },
    { k: 'Last activity', v: 'Today' }
  ];

  let tab = $state('prep');
  const tabs = [
    { id: 'prep', label: 'Prep' },
    { id: 'activity', label: 'Activity' },
    { id: 'details', label: 'Details' }
  ];
</script>

<svelte:head><title>Detail · C — Spotlight + tabs — Pursuit</title></svelte:head>

<main class="wrap">
  <div class="topbar">
    <a class="back" href="/preview/redesign">
      <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.7"><path d="M10 3l-5 5 5 5"/></svg>
      All mockups
    </a>
    <span class="badge">Static mockup</span>
  </div>

  <div class="body">
    <!-- HEADER -->
    <header class="head">
      <img class="logo" src={`https://www.google.com/s2/favicons?sz=128&domain=${app.domain}`} alt={app.co} />
      <div class="head-text">
        <div class="co-row">
          <h1>{app.co}</h1>
          <span class="pill applied"><span class="pdot"></span>{app.status}</span>
        </div>
        <div class="role">{app.role}</div>
        <div class="meta">
          <span>{app.appliedFull}</span>
          <span class="dot">·</span>
          <span>{app.source}</span>
          <span class="dot">·</span>
          <span>{app.location}</span>
          <span class="dot">·</span>
          <span>{app.salary}</span>
        </div>
      </div>
    </header>

    <!-- SPOTLIGHT ROW -->
    <section class="spotlight">
      <article class="spot">
        <div class="spot-top">
          <img class="spot-logo" src={`https://www.google.com/s2/favicons?sz=128&domain=${app.domain}`} alt={app.co} />
          <div class="spot-head">
            <div class="spot-kicker">Company</div>
            <h2>{app.co}</h2>
          </div>
          <span class="ai-pill">✦ AI</span>
        </div>
        <p class="spot-blurb">{company.blurb}</p>
        <div class="chips">
          {#each company.facts as f}<span class="chip">{f}</span>{/each}
        </div>
        <div class="process">
          <span class="proc-lbl">Process</span>
          <span class="proc-val">{company.process}</span>
        </div>
      </article>

      <article class="spot">
        <div class="spot-top">
          <div class="avatar">{person.initials}</div>
          <div class="spot-head">
            <div class="spot-kicker hm">Hiring manager</div>
            <h2>{person.name}</h2>
          </div>
          <span class="ai-pill">✦ AI</span>
        </div>
        <div class="person-title">{person.title}</div>
        <p class="spot-blurb">{person.who} Probes how you de-risk decisions, not just outcomes.</p>
      </article>
    </section>
    <p class="spot-caption">No interview scheduled yet — this prep is ready for when one is.</p>

    <!-- TABS -->
    <div class="segmented" role="tablist">
      {#each tabs as t}
        <button
          role="tab"
          class="seg"
          class:active={tab === t.id}
          aria-selected={tab === t.id}
          onclick={() => (tab = t.id)}
        >{t.label}</button>
      {/each}
    </div>

    <!-- TAB PANELS -->
    {#if tab === 'prep'}
      <div class="panel">
        <section class="block">
          <h3>How {person.name.split(' ')[0]} interviews</h3>
          <p class="lead">{person.snapshotLead}</p>
          <div class="tells">
            {#each person.tells as t}
              <div class="tell"><span class="tell-k">{t.k}</span><span class="tell-v">{t.v}</span></div>
            {/each}
          </div>
        </section>

        <section class="block two-col">
          <div class="col lands">
            <h4><span class="glyph">+</span>What lands</h4>
            <ul>
              {#each lands as l}<li>{l}</li>{/each}
            </ul>
          </div>
          <div class="col avoid">
            <h4><span class="glyph">−</span>What to avoid</h4>
            <ul>
              {#each avoid as a}<li>{a}</li>{/each}
            </ul>
          </div>
        </section>

        <section class="block">
          <h3>Recent signals</h3>
          <div class="signals">
            {#each signals as s}
              <div class="signal">
                <div class="signal-meta">
                  <span class="sig-kind">{s.kind}</span>
                  <span class="sig-date">{s.date}</span>
                </div>
                <p class="sig-body">{s.body}</p>
                <span class="sig-src">{s.source}</span>
              </div>
            {/each}
          </div>
        </section>

        <section class="block">
          <h3>Questions worth asking</h3>
          <div class="questions">
            {#each questions as item}
              <div class="q">
                <p class="q-text">{item.q}</p>
                <p class="q-why"><span class="why-lbl">Why</span>{item.why}</p>
              </div>
            {/each}
          </div>
        </section>
      </div>
    {/if}

    {#if tab === 'activity'}
      <div class="panel">
        <section class="block">
          <div class="timeline">
            {#each activity as ev}
              <div class="event">
                <span class="ev-node"></span>
                <div class="ev-body">
                  <div class="ev-title">{ev.title}</div>
                  <div class="ev-detail">{ev.detail}</div>
                </div>
                <span class="ev-date">{ev.date}</span>
              </div>
            {/each}
          </div>
        </section>
      </div>
    {/if}

    {#if tab === 'details'}
      <div class="panel">
        <section class="block">
          <dl class="details">
            {#each details as d}
              <div class="drow"><dt>{d.k}</dt><dd>{d.v}</dd></div>
            {/each}
            <div class="drow"><dt>Contact</dt><dd>{person.name} <span class="sep">·</span> <a class="link" href="#">LinkedIn</a></dd></div>
          </dl>
        </section>
      </div>
    {/if}
  </div>
</main>

<style>
  .wrap {
    min-height: 100vh;
    background: var(--surface);
    color: var(--ink);
    font-family: var(--sans);
  }

  .topbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    max-width: 880px;
    margin: 0 auto;
    padding: 18px 24px 0;
  }
  .back {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    color: var(--mute);
    font-size: 13px;
    font-weight: 500;
  }
  .back:hover { color: var(--ink-2); }
  .badge {
    font-family: var(--mono);
    font-size: 11px;
    color: var(--mute);
    background: var(--surface-2);
    border: 1px solid var(--rule);
    border-radius: 999px;
    padding: 3px 9px;
  }

  .body {
    max-width: 880px;
    margin: 0 auto;
    padding: 22px 24px 72px;
  }

  /* HEADER */
  .head {
    display: grid;
    grid-template-columns: 52px 1fr;
    gap: 16px;
    align-items: start;
    padding-bottom: 22px;
  }
  .logo {
    width: 52px; height: 52px;
    border-radius: 12px;
    border: 1px solid var(--rule);
    background: var(--card);
    object-fit: contain;
    padding: 6px;
  }
  .co-row {
    display: flex;
    align-items: center;
    gap: 12px;
  }
  .co-row h1 {
    margin: 0;
    font-size: 24px;
    font-weight: 650;
    letter-spacing: -0.02em;
    line-height: 1.1;
  }
  .role {
    margin-top: 3px;
    font-size: 15px;
    color: var(--ink-2);
    font-weight: 500;
  }
  .meta {
    margin-top: 8px;
    font-size: 13px;
    color: var(--mute);
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 7px;
  }
  .meta .dot { color: var(--mute-2); }

  .pill {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    font-size: 12px;
    font-weight: 550;
    padding: 3px 10px;
    border-radius: 999px;
  }
  .pill .pdot { width: 6px; height: 6px; border-radius: 50%; background: currentColor; }
  .pill.applied { color: var(--accent-text); background: var(--accent-tint); }

  /* SPOTLIGHT */
  .spotlight {
    display: grid;
    grid-template-columns: 1.15fr 1fr;
    gap: 16px;
  }
  .spot {
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 14px;
    padding: 18px 18px 20px;
    box-shadow: var(--sh-1);
    display: flex;
    flex-direction: column;
  }
  .spot-top {
    display: grid;
    grid-template-columns: 40px 1fr auto;
    gap: 12px;
    align-items: center;
  }
  .spot-logo {
    width: 40px; height: 40px;
    border-radius: 10px;
    border: 1px solid var(--rule);
    background: var(--surface);
    object-fit: contain;
    padding: 5px;
  }
  .avatar {
    width: 40px; height: 40px;
    border-radius: 10px;
    background: var(--accent-tint);
    color: var(--accent-text);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    font-weight: 650;
    letter-spacing: 0.02em;
  }
  .spot-kicker {
    font-size: 11px;
    text-transform: uppercase;
    letter-spacing: 0.07em;
    font-weight: 600;
    color: var(--mute);
  }
  .spot-kicker.hm { color: var(--warm-text); }
  .spot-head h2 {
    margin: 2px 0 0;
    font-size: 17px;
    font-weight: 620;
    letter-spacing: -0.01em;
  }
  .ai-pill {
    align-self: start;
    font-size: 11px;
    font-weight: 600;
    color: var(--accent-text);
    background: var(--accent-tint);
    border: 1px solid var(--accent-tint-2);
    border-radius: 999px;
    padding: 3px 9px;
    letter-spacing: 0.01em;
    white-space: nowrap;
  }
  .person-title {
    margin-top: 14px;
    font-size: 13px;
    font-weight: 550;
    color: var(--ink-2);
  }
  .spot-blurb {
    margin: 12px 0 0;
    font-size: 13.5px;
    line-height: 1.55;
    color: var(--ink-2);
  }
  .person-title + .spot-blurb { margin-top: 6px; }
  .chips {
    margin-top: 14px;
    display: flex;
    flex-wrap: wrap;
    gap: 7px;
  }
  .chip {
    font-size: 12px;
    color: var(--ink-2);
    background: var(--surface-2);
    border: 1px solid var(--rule);
    border-radius: 7px;
    padding: 4px 9px;
    font-weight: 500;
  }
  .process {
    margin-top: 14px;
    padding-top: 14px;
    border-top: 1px solid var(--rule);
    display: flex;
    flex-direction: column;
    gap: 3px;
  }
  .proc-lbl {
    font-size: 11px;
    text-transform: uppercase;
    letter-spacing: 0.07em;
    font-weight: 600;
    color: var(--mute);
  }
  .proc-val {
    font-size: 13px;
    color: var(--ink-2);
    font-weight: 500;
  }
  .spot-caption {
    margin: 14px 2px 0;
    font-size: 12.5px;
    color: var(--mute);
    font-style: italic;
  }

  /* SEGMENTED CONTROL */
  .segmented {
    margin-top: 26px;
    display: inline-flex;
    gap: 2px;
    padding: 3px;
    background: var(--surface-2);
    border: 1px solid var(--rule);
    border-radius: 10px;
  }
  .seg {
    font-size: 13px;
    font-weight: 530;
    color: var(--mute);
    padding: 6px 16px;
    border-radius: 7px;
    transition: background 0.12s, color 0.12s;
  }
  .seg:hover { color: var(--ink-2); }
  .seg.active {
    background: var(--card);
    color: var(--ink);
    box-shadow: var(--sh-1), 0 0 0 1px var(--rule);
  }

  /* PANELS */
  .panel {
    margin-top: 18px;
    display: flex;
    flex-direction: column;
    gap: 14px;
  }
  .block {
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 14px;
    padding: 18px 20px 20px;
    box-shadow: var(--sh-1);
  }
  .block h3 {
    margin: 0 0 12px;
    font-size: 14px;
    font-weight: 620;
    letter-spacing: -0.01em;
  }
  .lead {
    margin: 0;
    font-size: 13.5px;
    line-height: 1.6;
    color: var(--ink-2);
  }
  .tells {
    margin-top: 14px;
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }
  .tell {
    display: flex;
    flex-direction: column;
    gap: 2px;
    background: var(--surface-2);
    border: 1px solid var(--rule);
    border-radius: 9px;
    padding: 8px 12px;
    min-width: 120px;
  }
  .tell-k {
    font-size: 11px;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--mute);
    font-weight: 600;
  }
  .tell-v {
    font-size: 12.5px;
    color: var(--ink-2);
    font-weight: 500;
  }

  /* lands / avoid */
  .two-col {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 24px;
  }
  .col h4 {
    margin: 0 0 10px;
    font-size: 12px;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    font-weight: 650;
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .col .glyph {
    width: 18px; height: 18px;
    border-radius: 6px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    font-size: 13px;
    font-weight: 700;
    line-height: 1;
  }
  .col.lands h4 { color: var(--positive-text); }
  .col.lands .glyph { background: var(--positive-tint); color: var(--positive-text); }
  .col.avoid h4 { color: var(--warm-text); }
  .col.avoid .glyph { background: var(--warm-tint); color: var(--warm-text); }
  .col ul {
    margin: 0;
    padding: 0;
    list-style: none;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
  .col li {
    font-size: 13px;
    line-height: 1.5;
    color: var(--ink-2);
    padding-left: 14px;
    position: relative;
  }
  .col li::before {
    content: "";
    position: absolute;
    left: 0; top: 8px;
    width: 5px; height: 5px;
    border-radius: 50%;
  }
  .col.lands li::before { background: var(--positive); }
  .col.avoid li::before { background: var(--warm); }

  /* signals */
  .signals {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
  .signal {
    border-left: 2px solid var(--rule-strong);
    padding-left: 14px;
  }
  .signal-meta {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .sig-kind {
    font-size: 11px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--accent-text);
    background: var(--accent-tint);
    border-radius: 5px;
    padding: 2px 7px;
  }
  .sig-date {
    font-family: var(--mono);
    font-size: 11.5px;
    color: var(--mute);
  }
  .sig-body {
    margin: 6px 0 4px;
    font-size: 13px;
    line-height: 1.5;
    color: var(--ink-2);
  }
  .sig-src {
    font-family: var(--mono);
    font-size: 11px;
    color: var(--mute-2);
  }

  /* questions */
  .questions {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
  .q {
    background: var(--surface-2);
    border: 1px solid var(--rule);
    border-radius: 10px;
    padding: 12px 14px;
  }
  .q-text {
    margin: 0;
    font-size: 13.5px;
    font-weight: 550;
    color: var(--ink);
    line-height: 1.45;
  }
  .q-why {
    margin: 6px 0 0;
    font-size: 12.5px;
    color: var(--mute);
    line-height: 1.45;
  }
  .why-lbl {
    font-size: 10.5px;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    font-weight: 600;
    color: var(--mute-2);
    margin-right: 7px;
  }

  /* timeline */
  .timeline {
    display: flex;
    flex-direction: column;
  }
  .event {
    display: grid;
    grid-template-columns: 18px 1fr auto;
    gap: 12px;
    align-items: start;
    padding: 10px 0;
    position: relative;
  }
  .event:not(:last-child) { border-bottom: 1px solid var(--rule); }
  .ev-node {
    width: 9px; height: 9px;
    border-radius: 50%;
    background: var(--accent);
    margin-top: 4px;
    margin-left: 4px;
  }
  .ev-title { font-size: 13.5px; font-weight: 550; color: var(--ink); }
  .ev-detail { font-size: 12.5px; color: var(--mute); margin-top: 2px; }
  .ev-date {
    font-family: var(--mono);
    font-size: 11.5px;
    color: var(--mute);
    padding-top: 1px;
  }

  /* details */
  .details { margin: 0; }
  .drow {
    display: grid;
    grid-template-columns: 130px 1fr;
    gap: 12px;
    padding: 11px 0;
  }
  .drow:not(:last-child) { border-bottom: 1px solid var(--rule); }
  .drow dt {
    font-size: 12.5px;
    color: var(--mute);
    font-weight: 500;
  }
  .drow dd {
    margin: 0;
    font-size: 13.5px;
    color: var(--ink-2);
    font-weight: 500;
  }
  .drow .sep { color: var(--mute-2); }
  .link { color: var(--accent-text); font-weight: 550; }
  .link:hover { text-decoration: underline; }

  @media (max-width: 720px) {
    .spotlight { grid-template-columns: 1fr; }
    .two-col { grid-template-columns: 1fr; gap: 18px; }
    .head { grid-template-columns: 44px 1fr; }
  }
</style>
