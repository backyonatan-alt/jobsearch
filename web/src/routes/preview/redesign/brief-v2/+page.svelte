<script>
  // Brief v2 — adds an "About [Company]" card above the interviewer block.
  // Static mockup. Generated content is hard-coded; real flow will fire one
  // Claude call (web_search on) that returns both interviewer + company.
  const app = {
    co: 'Stripe',
    domain: 'stripe.com',
    role: 'Staff Engineer, Payments Infrastructure',
    location: 'Remote · US',
    status: 'interview',
    applied: 'May 18',
    source: 'Referral · Maya Patel',
    cv: 'v3',
    jdUrl: 'https://stripe.com/jobs/listing/staff-engineer-payments/abc123'
  };
  const company = {
    blurb: 'Payments + financial infrastructure used by a meaningful share of the internet.',
    direction: 'Recently shipping agentic-commerce APIs and pushing into AI-mediated transactions. Last 12 months: agent toolkit, deeper Issuing rollouts, expanded Atlas.',
    hq: 'San Francisco',
    employees: '~8,000',
    stage: 'Late stage · $50B valuation',
    founded: '2010',
    process: [
      { kind: 'Recruiter screen', detail: '30 min · culture + role fit' },
      { kind: 'Hiring-manager call', detail: '45 min · architecture, prior wins' },
      { kind: 'Technical screen',  detail: '60 min · live coding (Go or Python)' },
      { kind: 'Onsite loop',       detail: '4 panels · sys design, IC code, behavioral, bar-raiser' },
      { kind: 'Team match',        detail: 'Two 30-min chats with prospective teams' }
    ],
    watchFors: [
      'Sys design is graded on rollout + observability, not just the diagram.',
      'Be explicit about failure modes — they probe consistency tradeoffs.',
      'Stripe interviewers value "what would you measure?" over hand-waving.',
      'Bar raiser is structured — STAR format pays off here.',
      'Have a take on agentic commerce; they\'ll ask "where do you think this goes?"'
    ]
  };
  let generating = $state(false);
  let dossierAvailable = $state(true);
</script>

<svelte:head><title>Brief v2 — Pursuit</title></svelte:head>

<main class="wrap">
  <div class="topbar-stub">
    <div class="crumb"><span class="root">Applications</span><span class="sep">/</span><span class="here">{app.co}</span></div>
  </div>

  <div class="body-inner">
    <!-- HERO STRIP (unchanged) -->
    <div class="hero">
      <div class="hero-top">
        <img class="logo-big" src={`https://www.google.com/s2/favicons?sz=128&domain=${app.domain}`} alt={app.co} />
        <div class="hero-text">
          <div class="co-row">
            <h1>{app.co}</h1>
            <span class="pill interview"><span class="pdot"></span>Interview</span>
          </div>
          <div class="role-line">{app.role}</div>
        </div>
        <a class="src-link" href={app.jdUrl} target="_blank" rel="noopener">
          <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.6"><path d="M5 11l6-6M6 5h5v5"/></svg>
          View on {app.domain}
        </a>
      </div>
      <div class="facts">
        <span class="fact"><span class="fdot d-loc"></span>{app.location}</span>
        <span class="fact"><span class="fdot d-app"></span>Applied {app.applied}</span>
        <span class="fact"><span class="fdot d-src"></span>{app.source}</span>
        <span class="fact"><span class="fdot d-cv"></span>CV {app.cv}</span>
      </div>
    </div>

    <!-- STATS (3 cards, unchanged) -->
    <div class="stats">
      <div class="stat tone-accent">
        <span class="ribbon"></span>
        <div class="stat-lbl">Days in pipeline</div>
        <div class="stat-n">7</div>
        <div class="stat-sub">Applied {app.applied}</div>
      </div>
      <div class="stat tone-positive">
        <span class="ribbon"></span>
        <div class="stat-lbl">Current stage</div>
        <div class="stat-n">3<span class="of">/ 4</span></div>
        <div class="stat-sub">Interview</div>
      </div>
      <div class="stat tone-warm">
        <span class="ribbon"></span>
        <div class="stat-lbl">Match score</div>
        <div class="stat-n">A−</div>
        <div class="stat-sub">CV {app.cv} vs. JD</div>
      </div>
    </div>

    <!-- TABS -->
    <div class="tabs">
      <button class="tab active">Brief <span class="t-tag">AI</span></button>
      <button class="tab">Interviews <span class="t-tag">2</span></button>
      <button class="tab">Timeline <span class="t-tag">3</span></button>
      <button class="tab">Notes</button>
      <button class="tab">Files</button>
    </div>

    <!-- NEW: ABOUT COMPANY CARD (compact, sits above interviewer) -->
    <div class="block company-block">
      <div class="block-hd">
        <h2>About {app.co}</h2>
        <span class="ai-tag">AI · refreshed today</span>
      </div>
      <p class="company-blurb">{company.blurb}</p>
      <p class="company-direction">{company.direction}</p>
      <div class="facts-grid">
        <div class="f-cell"><div class="f-lbl">HQ</div><div class="f-val">{company.hq}</div></div>
        <div class="f-cell"><div class="f-lbl">Employees</div><div class="f-val">{company.employees}</div></div>
        <div class="f-cell"><div class="f-lbl">Stage</div><div class="f-val">{company.stage}</div></div>
        <div class="f-cell"><div class="f-lbl">Founded</div><div class="f-val">{company.founded}</div></div>
      </div>

      <div class="sub-hd">
        <h3>Typical interview process</h3>
        <span class="ai-tag-mute">from Glassdoor / Blind / levels.fyi</span>
      </div>
      <ol class="process-list">
        {#each company.process as p, i}
          <li>
            <span class="step-n">{i + 1}</span>
            <div>
              <div class="step-kind">{p.kind}</div>
              <div class="step-detail">{p.detail}</div>
            </div>
          </li>
        {/each}
      </ol>

      <div class="sub-hd">
        <h3>Watch-fors for this loop</h3>
        <span class="ai-tag-mute">AI · synthesized</span>
      </div>
      <ul class="watch-list">
        {#each company.watchFors as w}
          <li>
            <span class="w-dot">▸</span>
            <span>{w}</span>
          </li>
        {/each}
      </ul>
    </div>

    <!-- INTERVIEWER BLOCK (existing, unchanged in layout) -->
    <div class="block person-block">
      <div class="block-hd">
        <h2>Your interviewer</h2>
        <span class="ai-tag">AI · refreshed today</span>
      </div>
      <div class="person">
        <div class="p-av t-accent">SC</div>
        <div class="p-info">
          <h4>Sarah Chen</h4>
          <div class="p-role">Staff Engineer · Payments Routing</div>
        </div>
        <a class="p-li" href="#" target="_blank" rel="noopener">
          <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M3.5 6h2v6h-2zM4.5 3a1 1 0 1 1 0 2 1 1 0 0 1 0-2zM7 6h2v.9c.3-.5.9-1 1.8-1 1.6 0 2.2 1 2.2 2.6V12h-2V9c0-.9-.3-1.4-1.1-1.4-.6 0-1 .4-1 1.2V12H7z"/></svg>
          LinkedIn
        </a>
      </div>
      <div class="prior-row">
        <span class="p-lbl">Prior</span>
        <span class="prior-chip">Cloudflare</span>
        <span class="prior-chip">Two Sigma</span>
        <span class="prior-chip">CMU MS</span>
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

  .body-inner { max-width: 1080px; margin: 0 auto; }

  /* HERO */
  .hero { background: var(--card); border: 1px solid var(--rule); border-radius: 18px; padding: 22px 24px; margin-bottom: 14px; box-shadow: var(--sh-1); }
  .hero-top { display: grid; grid-template-columns: 64px 1fr auto; gap: 18px; align-items: center; }
  .logo-big { width: 64px; height: 64px; border-radius: 14px; background: var(--card); object-fit: contain; padding: 8px; border: 1px solid var(--rule); }
  .co-row { display: flex; align-items: center; gap: 12px; }
  .co-row h1 { font-size: 28px; font-weight: 600; margin: 0; letter-spacing: -0.025em; }
  .role-line { font-size: 15px; color: var(--ink-2); margin-top: 2px; font-weight: 500; }
  .src-link { display: inline-flex; align-items: center; gap: 6px; background: var(--card); border: 1px solid var(--rule); border-radius: 99px; padding: 8px 14px; font-size: 12.5px; font-weight: 500; color: var(--ink); text-decoration: none; }
  .facts { display: flex; flex-wrap: wrap; gap: 6px; margin-top: 14px; padding-top: 14px; border-top: 1px dashed var(--rule); }
  .fact { display: inline-flex; align-items: center; gap: 6px; font-size: 12.5px; color: var(--ink-2); background: var(--card); border: 1px solid var(--rule); padding: 4px 10px; border-radius: 99px; font-weight: 500; }
  .fdot { width: 6px; height: 6px; border-radius: 50%; }
  .d-loc { background: var(--accent); } .d-app { background: var(--positive); } .d-src { background: var(--warm); } .d-cv { background: var(--mute-2); }

  /* PILL */
  .pill { display: inline-flex; align-items: center; gap: 5px; padding: 4px 10px; border-radius: 99px; font-size: 12px; font-weight: 500; background: var(--surface-2); color: var(--ink-2); }
  .pill .pdot { width: 6px; height: 6px; border-radius: 50%; background: var(--mute-2); }
  .pill.interview { background: var(--accent-tint); color: var(--accent-text); }
  .pill.interview .pdot { background: var(--accent); }

  /* STATS */
  .stats { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; margin-bottom: 24px; }
  .stat { position: relative; background: var(--card); border: 1px solid var(--rule); border-radius: 14px; padding: 16px 18px; overflow: hidden; box-shadow: var(--sh-1); }
  .stat .ribbon { position: absolute; top: 0; left: 0; right: 0; height: 3px; }
  .stat.tone-accent .ribbon { background: var(--accent); }
  .stat.tone-positive .ribbon { background: var(--positive); }
  .stat.tone-warm .ribbon { background: var(--warm); }
  .stat-lbl { font-size: 12.5px; color: var(--mute); margin-bottom: 4px; font-weight: 500; }
  .stat-n { font-size: 32px; font-weight: 600; letter-spacing: -0.035em; line-height: 1.1; }
  .stat.tone-accent .stat-n { color: var(--accent-text); }
  .stat.tone-positive .stat-n { color: var(--positive-text); }
  .stat.tone-warm .stat-n { color: var(--warm-text); }
  .of { font-size: 16px; color: var(--mute); margin-left: 6px; font-weight: 500; }
  .stat-sub { font-size: 12px; color: var(--mute); margin-top: 6px; padding-top: 6px; border-top: 1px dashed var(--rule); }

  /* TABS */
  .tabs { display: flex; gap: 4px; border-bottom: 1px solid var(--rule); margin-bottom: 18px; }
  .tab { background: transparent; border: 0; padding: 10px 14px; font-size: 13.5px; color: var(--mute); cursor: pointer; border-bottom: 2px solid transparent; margin-bottom: -1px; font-weight: 600; }
  .tab.active { color: var(--ink); border-bottom-color: var(--ink); }
  .t-tag { font-size: 11px; background: var(--accent-tint); color: var(--accent-text); padding: 1px 7px; border-radius: 99px; margin-left: 4px; font-weight: 500; }

  /* BLOCK */
  .block { background: var(--card); border: 1px solid var(--rule); border-radius: 14px; padding: 20px 22px; margin-bottom: 14px; box-shadow: var(--sh-1); }
  .block-hd { display: flex; align-items: center; gap: 10px; margin-bottom: 14px; }
  .block-hd h2 { font-size: 16px; font-weight: 600; margin: 0; letter-spacing: -0.015em; }
  .ai-tag { display: inline-flex; align-items: center; gap: 5px; font-size: 12px; background: var(--accent-tint); color: var(--accent-text); padding: 3px 10px; border-radius: 99px; font-weight: 500; }
  .ai-tag-mute { display: inline-flex; align-items: center; gap: 5px; font-size: 11.5px; background: var(--surface-2); color: var(--mute); padding: 3px 9px; border-radius: 99px; font-weight: 500; }

  /* COMPANY CARD */
  .company-blurb { font-size: 14.5px; line-height: 1.55; color: var(--ink); margin: 0 0 6px; font-weight: 500; }
  .company-direction { font-size: 13.5px; line-height: 1.6; color: var(--ink-2); margin: 0 0 14px; }
  .facts-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 8px; margin-bottom: 18px; }
  .f-cell { background: var(--surface-2); border-radius: 10px; padding: 10px 12px; }
  .f-lbl { font-size: 11px; color: var(--mute); font-weight: 500; }
  .f-val { font-size: 13px; color: var(--ink); font-weight: 600; margin-top: 2px; }

  .sub-hd { display: flex; align-items: center; gap: 10px; margin: 18px 0 10px; padding-top: 14px; border-top: 1px dashed var(--rule); }
  .sub-hd h3 { font-size: 13.5px; font-weight: 600; margin: 0; letter-spacing: -0.005em; color: var(--ink); }

  /* PROCESS LIST */
  .process-list { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 6px; }
  .process-list li { display: grid; grid-template-columns: 24px 1fr; gap: 10px; align-items: baseline; padding: 8px 0; border-bottom: 1px dashed var(--rule); }
  .process-list li:last-child { border-bottom: 0; }
  .step-n { width: 22px; height: 22px; border-radius: 50%; background: var(--accent-tint); color: var(--accent-text); display: grid; place-items: center; font-size: 11px; font-weight: 600; font-feature-settings: "tnum"; }
  .step-kind { font-size: 13.5px; color: var(--ink); font-weight: 600; }
  .step-detail { font-size: 12.5px; color: var(--mute); margin-top: 2px; }

  /* WATCH-FORS LIST */
  .watch-list { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 6px; }
  .watch-list li { display: grid; grid-template-columns: 18px 1fr; gap: 8px; align-items: baseline; padding: 6px 0; font-size: 13.5px; color: var(--ink-2); line-height: 1.5; }
  .w-dot { color: var(--accent); font-weight: 700; }

  /* PERSON */
  .person-block { padding: 18px 20px; }
  .person { display: grid; grid-template-columns: 46px 1fr auto; gap: 12px; align-items: center; }
  .p-av { width: 46px; height: 46px; border-radius: 50%; display: grid; place-items: center; font-weight: 600; font-size: 16px; background: var(--accent-tint); color: var(--accent-text); }
  .p-info h4 { margin: 0; font-size: 14.5px; font-weight: 600; letter-spacing: -0.01em; }
  .p-role { font-size: 12.5px; color: var(--mute); margin-top: 2px; }
  .p-li { display: inline-flex; align-items: center; gap: 6px; background: var(--surface-2); border: 1px solid var(--rule); border-radius: 99px; padding: 6px 12px; font-size: 12px; font-weight: 600; color: var(--ink); text-decoration: none; }
  .p-li svg { color: #0a66c2; }
  .prior-row { display: flex; align-items: center; gap: 8px; margin-top: 14px; padding-top: 12px; border-top: 1px dashed var(--rule); flex-wrap: wrap; }
  .p-lbl { font-size: 11.5px; color: var(--mute); font-weight: 600; }
  .prior-chip { font-size: 11.5px; background: var(--surface-2); color: var(--ink-2); padding: 3px 9px; border-radius: 99px; font-weight: 500; }

  .footer-link { margin-top: 30px; font-size: 13px; }
  .footer-link a { color: var(--accent-text); text-decoration: none; }
</style>
