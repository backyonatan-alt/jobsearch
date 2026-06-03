<script>
  // Static mockup — Layout A: "Prep hero up top, two spotlight cards".
  // All data hard-coded inline; no backend, no fetching.
  const app = {
    co: 'Astelia',
    domain: 'astelia.com',
    role: 'Senior Product Manager',
    status: 'Applied',
    applied: 'June 3 2026',
    source: 'LinkedIn',
    location: 'Remote (US)',
    salary: '$190k–$230k'
  };

  const company = {
    about:
      'Astelia builds infrastructure for climate-tech data. Series B, ~120 people, founded 2021. Product culture is research-heavy — they value crisp problem framing over feature lists.',
    facts: ['Stage Series B', '~120 people', 'Founded 2021', 'Remote-first'],
    process: 'Recruiter screen → PM craft → product sense → exec/values'
  };

  const hm = {
    name: 'Sarah Chen',
    title: 'VP Product',
    initials: 'SC',
    snapshot:
      'Led product at two infra startups; writes about discovery-driven roadmaps. Probes how you de-risk decisions, not just outcomes.',
    tells: [
      { k: 'Format', v: 'conversational / no slides' },
      { k: 'Tempo', v: 'digs into one decision' },
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
    {
      date: 'Apr 2026',
      kind: 'Talk',
      body: "'Discovery debt' talk at ProductCon — argued teams ship too early.",
      source: 'productcon.com'
    },
    {
      date: 'Feb 2026',
      kind: 'Post',
      body: 'Essay on killing features as a PM skill.',
      source: 'sarahchen.substack.com'
    }
  ];
  const questions = [
    {
      q: 'How do you decide what NOT to build this quarter?',
      why: "Mirrors her 'discovery debt' thesis."
    },
    {
      q: 'What does great PM craft look like here in 6 months?',
      why: 'Shows you think in outcomes + ramp.'
    }
  ];

  const timeline = [
    { title: 'Application submitted', detail: 'via LinkedIn', date: 'June 3 2026' },
    { title: 'Status', detail: 'Applied', date: 'June 3 2026' }
  ];
</script>

<svelte:head><title>Detail (redesign A) — Pursuit</title></svelte:head>

<main class="wrap">
  <div class="topbar">
    <a class="back" href="/preview/redesign">← back to previews</a>
    <span class="mock-badge">Static mockup</span>
  </div>

  <!-- HEADER -->
  <header class="head">
    <div class="head-main">
      <div class="logo">{app.co.charAt(0)}</div>
      <div class="head-text">
        <div class="co-row">
          <h1>{app.co}</h1>
          <span class="pill applied"><span class="pdot"></span>{app.status}</span>
        </div>
        <div class="role">{app.role}</div>
        <div class="meta">
          Applied {app.applied} <span class="dot">·</span> Source {app.source}
          <span class="dot">·</span> {app.location} <span class="dot">·</span> {app.salary}
        </div>
      </div>
    </div>
    <div class="status-strip">
      No interview scheduled yet — here's your prep for when you hear back.
    </div>
  </header>

  <!-- AI PREP HERO -->
  <section class="hero">
    <div class="hero-head">
      <span class="ai-pill">✦ AI</span>
      <h2>Prep for {app.co}</h2>
    </div>

    <div class="spotlights">
      <!-- COMPANY CARD -->
      <article class="spot">
        <div class="spot-head">
          <div class="logo sm">{app.co.charAt(0)}</div>
          <div>
            <div class="spot-eyebrow">Company</div>
            <div class="spot-title">{app.co}</div>
          </div>
        </div>
        <p class="spot-body">{company.about}</p>
        <div class="facts">
          {#each company.facts as f}<span class="fact">{f}</span>{/each}
        </div>
        <div class="process">
          <span class="process-label">Typical process</span>
          <span class="process-val">{company.process}</span>
        </div>
      </article>

      <!-- HIRING MANAGER CARD -->
      <article class="spot">
        <div class="spot-head">
          <div class="avatar">{hm.initials}</div>
          <div>
            <div class="spot-eyebrow">Hiring manager</div>
            <div class="spot-title">{hm.name} <span class="spot-sub">· {hm.title}</span></div>
          </div>
        </div>
        <p class="spot-body">{hm.snapshot}</p>
        <div class="tells">
          {#each hm.tells as t}
            <div class="tell"><span class="tell-k">{t.k}</span><span class="tell-v">{t.v}</span></div>
          {/each}
        </div>
      </article>
    </div>

    <!-- PREP BODY -->
    <div class="prep-body">
      <div class="two-col">
        <div class="col">
          <h3 class="col-head">What lands</h3>
          <ul class="checklist">
            {#each lands as item}
              <li><span class="mark good">✓</span><span>{item}</span></li>
            {/each}
          </ul>
        </div>
        <div class="col">
          <h3 class="col-head">What to avoid</h3>
          <ul class="checklist">
            {#each avoid as item}
              <li><span class="mark bad">×</span><span>{item}</span></li>
            {/each}
          </ul>
        </div>
      </div>

      <div class="block">
        <h3 class="block-head">Recent signals</h3>
        <div class="signals">
          {#each signals as s}
            <div class="signal">
              <div class="signal-meta">
                <span class="signal-date">{s.date}</span>
                <span class="signal-kind">{s.kind}</span>
              </div>
              <div class="signal-body">{s.body}</div>
              <div class="signal-source">{s.source}</div>
            </div>
          {/each}
        </div>
      </div>

      <div class="block">
        <h3 class="block-head">Questions worth asking</h3>
        <div class="questions">
          {#each questions as item}
            <div class="question">
              <div class="q">{item.q}</div>
              <div class="why">{item.why}</div>
            </div>
          {/each}
        </div>
      </div>
    </div>
  </section>

  <!-- SECONDARY -->
  <section class="secondary">
    <div class="sec-card">
      <h3 class="sec-head">Details</h3>
      <dl class="kv">
        <dt>Status</dt><dd>{app.status}</dd>
        <dt>Source</dt><dd>{app.source}</dd>
        <dt>Location</dt><dd>{app.location}</dd>
        <dt>Salary</dt><dd>{app.salary}</dd>
        <dt>Last activity</dt><dd>{app.applied}</dd>
      </dl>
    </div>

    <div class="sec-card">
      <h3 class="sec-head">Contact</h3>
      <div class="contact">
        <div class="avatar sm">{hm.initials}</div>
        <div class="contact-text">
          <div class="contact-name">{hm.name}</div>
          <div class="contact-title">{hm.title}</div>
        </div>
      </div>
      <button class="add-link">+ Add LinkedIn</button>
    </div>

    <div class="sec-card">
      <h3 class="sec-head">Activity</h3>
      <ol class="timeline">
        {#each timeline as t}
          <li>
            <span class="tl-dot"></span>
            <div class="tl-text">
              <div class="tl-title">{t.title} <span class="tl-detail">{t.detail}</span></div>
              <div class="tl-date">{t.date}</div>
            </div>
          </li>
        {/each}
      </ol>
    </div>
  </section>
</main>

<style>
  :global(html, body) { background: var(--surface); margin: 0; font-family: var(--sans); color: var(--ink); }

  .wrap { max-width: 980px; margin: 0 auto; padding: 24px 24px 80px; }

  .topbar { display: flex; align-items: center; justify-content: space-between; margin-bottom: 18px; }
  .back { font-size: 13px; color: var(--mute); font-weight: 500; }
  .back:hover { color: var(--ink); }
  .mock-badge {
    font-family: var(--mono); font-size: 11px; letter-spacing: 0.02em;
    color: var(--mute); background: var(--surface-2); border: 1px solid var(--rule);
    padding: 3px 8px; border-radius: 6px;
  }

  /* HEADER */
  .head {
    background: var(--card); border: 1px solid var(--rule); border-radius: 16px;
    padding: 20px 22px; box-shadow: var(--sh-1); margin-bottom: 14px;
  }
  .head-main { display: flex; gap: 16px; align-items: flex-start; }
  .logo {
    width: 56px; height: 56px; border-radius: 13px; flex-shrink: 0;
    display: grid; place-items: center; font-size: 24px; font-weight: 600;
    color: var(--accent-text); background: var(--accent-tint); border: 1px solid var(--rule);
  }
  .logo.sm { width: 36px; height: 36px; font-size: 16px; border-radius: 9px; }
  .head-text { flex: 1; min-width: 0; }
  .co-row { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
  h1 { font-size: 24px; font-weight: 600; margin: 0; letter-spacing: -0.01em; }
  .role { font-size: 15px; color: var(--ink-2); margin-top: 2px; }
  .meta { font-size: 12.5px; color: var(--mute); margin-top: 8px; }
  .meta .dot { color: var(--mute-2); margin: 0 5px; }

  .pill { display: inline-flex; align-items: center; gap: 6px; padding: 4px 10px; border-radius: 99px; font-size: 12px; font-weight: 500; background: var(--surface-2); color: var(--ink-2); }
  .pill .pdot { width: 6px; height: 6px; border-radius: 50%; background: var(--mute-2); }
  .pill.applied { background: var(--accent-tint); color: var(--accent-text); }
  .pill.applied .pdot { background: var(--accent); }

  .status-strip {
    margin-top: 16px; padding-top: 14px; border-top: 1px solid var(--rule);
    font-size: 13px; color: var(--mute);
  }

  /* HERO */
  .hero {
    background: var(--card); border: 1px solid var(--rule); border-radius: 16px;
    padding: 22px 24px; box-shadow: var(--sh-1); margin-bottom: 14px;
  }
  .hero-head { display: flex; align-items: center; gap: 10px; margin-bottom: 18px; }
  .ai-pill {
    display: inline-flex; align-items: center; font-size: 12px; font-weight: 600;
    color: var(--accent-text); background: var(--accent-tint); border: 1px solid var(--accent-tint-2);
    padding: 3px 9px; border-radius: 99px; letter-spacing: 0.02em;
  }
  .hero-head h2 { font-size: 18px; font-weight: 600; margin: 0; letter-spacing: -0.01em; }

  .spotlights { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }
  .spot {
    background: var(--surface); border: 1px solid var(--rule); border-radius: 13px; padding: 18px;
  }
  .spot-head { display: flex; align-items: center; gap: 12px; margin-bottom: 12px; }
  .spot-eyebrow { font-size: 11px; font-weight: 600; letter-spacing: 0.04em; text-transform: uppercase; color: var(--mute); }
  .spot-title { font-size: 15.5px; font-weight: 600; }
  .spot-sub { font-weight: 400; color: var(--mute); font-size: 13.5px; }
  .spot-body { font-size: 13.5px; line-height: 1.55; color: var(--ink-2); margin: 0 0 14px; }

  .avatar {
    width: 36px; height: 36px; border-radius: 50%; flex-shrink: 0;
    display: grid; place-items: center; font-size: 13px; font-weight: 600;
    color: var(--warm-text); background: var(--warm-tint); border: 1px solid var(--rule);
  }
  .avatar.sm { width: 34px; height: 34px; }

  .facts { display: flex; flex-wrap: wrap; gap: 6px; margin-bottom: 14px; }
  .fact { font-size: 12px; color: var(--ink-2); background: var(--card); border: 1px solid var(--rule); padding: 3px 9px; border-radius: 99px; font-weight: 500; }

  .process { border-top: 1px solid var(--rule); padding-top: 12px; }
  .process-label { display: block; font-size: 11px; font-weight: 600; letter-spacing: 0.04em; text-transform: uppercase; color: var(--mute); margin-bottom: 4px; }
  .process-val { font-size: 13px; color: var(--ink-2); line-height: 1.5; }

  .tells { display: flex; flex-direction: column; gap: 8px; border-top: 1px solid var(--rule); padding-top: 12px; }
  .tell { display: flex; gap: 10px; font-size: 13px; }
  .tell-k { color: var(--mute); min-width: 56px; font-weight: 500; }
  .tell-v { color: var(--ink-2); }

  /* PREP BODY */
  .prep-body { margin-top: 22px; padding-top: 22px; border-top: 1px solid var(--rule); }
  .two-col { display: grid; grid-template-columns: 1fr 1fr; gap: 28px; margin-bottom: 24px; }
  .col-head, .block-head { font-size: 13px; font-weight: 600; margin: 0 0 12px; letter-spacing: -0.005em; }
  .checklist { list-style: none; margin: 0; padding: 0; display: flex; flex-direction: column; gap: 10px; }
  .checklist li { display: flex; gap: 10px; font-size: 13.5px; line-height: 1.5; color: var(--ink-2); }
  .mark { flex-shrink: 0; width: 18px; height: 18px; border-radius: 50%; display: grid; place-items: center; font-size: 12px; font-weight: 700; margin-top: 1px; }
  .mark.good { color: var(--positive-text); background: var(--positive-tint); }
  .mark.bad { color: var(--danger-text); background: var(--danger-tint); }

  .block { margin-top: 22px; }
  .signals { display: flex; flex-direction: column; gap: 12px; }
  .signal { background: var(--surface); border: 1px solid var(--rule); border-radius: 11px; padding: 13px 15px; }
  .signal-meta { display: flex; align-items: center; gap: 8px; margin-bottom: 5px; }
  .signal-date { font-family: var(--mono); font-size: 11.5px; color: var(--mute); }
  .signal-kind { font-size: 11px; font-weight: 600; color: var(--accent-text); background: var(--accent-tint); padding: 1px 7px; border-radius: 99px; }
  .signal-body { font-size: 13.5px; color: var(--ink-2); line-height: 1.5; }
  .signal-source { font-family: var(--mono); font-size: 11.5px; color: var(--mute-2); margin-top: 5px; }

  .questions { display: flex; flex-direction: column; gap: 12px; }
  .question { background: var(--surface); border: 1px solid var(--rule); border-radius: 11px; padding: 13px 15px; }
  .q { font-size: 13.5px; font-weight: 500; color: var(--ink); line-height: 1.45; }
  .why { font-size: 12.5px; color: var(--mute); margin-top: 4px; }

  /* SECONDARY */
  .secondary { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 14px; }
  .sec-card { background: var(--card); border: 1px solid var(--rule); border-radius: 13px; padding: 16px 18px; box-shadow: var(--sh-1); }
  .sec-head { font-size: 11px; font-weight: 600; letter-spacing: 0.04em; text-transform: uppercase; color: var(--mute); margin: 0 0 12px; }

  .kv { display: grid; grid-template-columns: auto 1fr; gap: 8px 14px; margin: 0; }
  .kv dt { font-size: 12.5px; color: var(--mute); }
  .kv dd { font-size: 12.5px; color: var(--ink-2); margin: 0; text-align: right; }

  .contact { display: flex; align-items: center; gap: 11px; margin-bottom: 12px; }
  .contact-name { font-size: 13.5px; font-weight: 600; }
  .contact-title { font-size: 12px; color: var(--mute); }
  .add-link { font-size: 12.5px; font-weight: 500; color: var(--accent-text); }
  .add-link:hover { text-decoration: underline; }

  .timeline { list-style: none; margin: 0; padding: 0; display: flex; flex-direction: column; gap: 14px; }
  .timeline li { display: flex; gap: 10px; position: relative; }
  .tl-dot { width: 8px; height: 8px; border-radius: 50%; background: var(--accent); margin-top: 5px; flex-shrink: 0; }
  .tl-title { font-size: 13px; color: var(--ink); }
  .tl-detail { color: var(--mute); font-weight: 400; }
  .tl-date { font-family: var(--mono); font-size: 11px; color: var(--mute-2); margin-top: 2px; }

  @media (max-width: 720px) {
    .spotlights, .two-col, .secondary { grid-template-columns: 1fr; }
  }
</style>
