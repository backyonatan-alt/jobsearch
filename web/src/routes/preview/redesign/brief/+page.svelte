<script>
  // Static mockup. Hard-coded data, nothing wired.
  const app = {
    co: 'Stripe',
    domain: 'stripe.com',
    role: 'Staff Engineer, Payments Infrastructure',
    location: 'Remote · US',
    status: 'interview',
    applied: 'May 18, 2026',
    appliedAgo: '7 days ago',
    source: 'Referral · Maya Patel',
    cv: 'v3 — infra-leaning',
    jdUrl: 'https://stripe.com/jobs/listing/staff-engineer-payments/abc123',
    salary: '$280k–$340k base + equity'
  };

  // AI company blurb
  const companyBlurb = `Stripe is the payments and financial-infrastructure layer behind a large share of internet businesses, from Shopify to OpenAI. Recently shipped Stripe Tax for marketplaces and is reportedly investing heavily in agentic-commerce APIs. Engineering culture leans toward written docs (the Atlas memo), strong on-call discipline, and a long-tenured senior bench.`;

  // JD highlights
  const jdHighlights = {
    summary: `Own the next generation of Stripe's payment-routing engine. Lead a small team designing for latency under 50ms p99 at multi-region scale. Heavy collaboration with risk, treasury, and the partner-bank integrations team.`,
    musts: ['10+ yrs distributed systems', 'Production Go or Rust at scale', 'Has owned a Tier-0 service'],
    nice:  ['Payments / fintech background', 'Multi-region active-active', 'Open-source maintainership']
  };

  // Hiring manager (LinkedIn)
  const hiringManager = {
    name: 'Devon Marquez',
    role: 'Director of Engineering, Payments Infra',
    linkedin: 'https://www.linkedin.com/in/example-devon-marquez',
    initials: 'DM'
  };

  // Interviewer block
  const interviewer = {
    name: 'Sarah Chen',
    role: 'Staff Eng · Payments Routing',
    initials: 'SC',
    linkedin: 'https://www.linkedin.com/in/example-sarah-chen',
    prior: ['Cloudflare (5y)', 'Two Sigma (3y)', 'CMU MS'],
    snapshot: `Sarah has been at Stripe for 2 years on the payments-routing team. She writes infrequently but well — her 2024 talk on Stripe's regional failover (linked below) is the canonical reference. Likely to spend most of the hour on system-design depth, especially failure modes.`,
    signals: [
      { date: 'Apr 26', kind: 'Talk', body: 'QCon SF — "Regional failover at p99.99"', source: 'qconsf.com' },
      { date: 'Mar 04', kind: 'Post', body: 'Twitter thread on idempotency-key collisions at scale', source: 'twitter.com' },
      { date: 'Jan 12', kind: 'Doc',  body: 'Cited in Stripe\'s engineering blog post on multi-region writes', source: 'stripe.com' }
    ],
    lands: [
      'Concrete failure-mode reasoning, not happy-path designs',
      'Numbers — latency budgets, RPS, error rates',
      'Operational empathy: rollback, observability, on-call'
    ],
    avoid: [
      'Hand-waving past consistency tradeoffs',
      'Big-bang rewrites without rollout plan',
      'Skipping straight to the database choice before constraints'
    ],
    questions: [
      { q: "What's the single hardest production incident the payments-routing team has shipped through this year?", why: 'Forces a real story; reveals how much she still owns vs. delegates.' },
      { q: 'How does the team decide between hot-patching the legacy router vs. building forward on v2?', why: 'Tests whether the org actually pays down debt or just talks about it.' },
      { q: 'What does a great first 90 days look like for this hire?', why: 'Standard but high-signal — if she fumbles it, the role is loosely defined.' }
    ]
  };

  let tab = $state('brief');
</script>

<svelte:head><title>{app.co} (redesign preview) — Pursuit</title></svelte:head>

<div class="frame">
  <!-- Compact sidebar -->
  <aside class="sidebar">
    <div class="brand"><span class="mark"></span><span class="name">Pursuit</span></div>
    <a class="nav-item"><span class="dot"></span>Today</a>
    <a class="nav-item active"><span class="dot"></span>Board</a>
    <a class="nav-item"><span class="dot"></span>Funnel</a>
    <div class="sidebar-footer">
      <div class="profile">
        <img class="av" src="https://www.gravatar.com/avatar/?d=mp&s=64" alt="" />
        <div class="who">Yonatan<small>back.yonatan@gmail.com</small></div>
      </div>
    </div>
  </aside>

  <section class="main">
    <!-- Top bar -->
    <div class="topbar">
      <div class="crumb">
        <span class="root">Applications</span>
        <span class="sep">/</span>
        <span class="here">{app.co}</span>
      </div>
      <div class="right">
        <button class="btn">Update status</button>
        <button class="btn">Edit</button>
        <button class="btn btn-danger">Delete</button>
        <img class="user-av" src="https://www.gravatar.com/avatar/?d=mp&s=64" alt="" />
      </div>
    </div>

    <div class="body">
      <div class="body-inner">

        <!-- Application header — now with real company logo -->
        <div class="app-hd">
          <img class="logo-big" src={`https://logo.clearbit.com/${app.domain}`} alt={app.co} />
          <div class="hd-text">
            <div class="co-line">
              <h1>{app.co}</h1>
              <span class="pill interview"><span class="pdot"></span>Interview</span>
            </div>
            <div class="role-line">{app.role}</div>
            <div class="sub">
              <span>{app.location}</span>
              <span class="dot">·</span>
              <span>Applied {app.applied}</span>
              <span class="dot">·</span>
              <span>via {app.source}</span>
              <span class="dot">·</span>
              <span>CV {app.cv}</span>
            </div>
          </div>
          <div class="hd-cta">
            <!-- The link to the source posting, prominent -->
            <a class="src-link" href={app.jdUrl} target="_blank" rel="noopener">
              <svg width="14" height="14" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M5 11l6-6M6 5h5v5"/></svg>
              View posting on stripe.com
            </a>
          </div>
        </div>

        <!-- Tabs -->
        <div class="tabs">
          <button class={`tab ${tab === 'brief' ? 'active' : ''}`} onclick={() => (tab = 'brief')}>
            Brief <span class="count">AI</span>
          </button>
          <button class={`tab ${tab === 'timeline' ? 'active' : ''}`} onclick={() => (tab = 'timeline')}>
            Timeline <span class="count">4</span>
          </button>
          <button class={`tab ${tab === 'notes' ? 'active' : ''}`} onclick={() => (tab = 'notes')}>Notes</button>
          <button class={`tab ${tab === 'files' ? 'active' : ''}`} onclick={() => (tab = 'files')}>Files</button>
        </div>

        {#if tab === 'brief'}
        <!-- ============ THE ROLE ============ -->
        <div class="block">
          <div class="block-hd">
            <h2>The role</h2>
            <span class="ai-tag">AI · from JD</span>
          </div>
          <div class="role-grid">
            <div class="role-summary">
              <p>{jdHighlights.summary}</p>
              <dl class="role-meta">
                <div><dt>Salary</dt><dd>{app.salary}</dd></div>
                <div><dt>Location</dt><dd>{app.location}</dd></div>
                <div><dt>Posted</dt><dd>{app.applied}</dd></div>
              </dl>
            </div>
            <div class="role-reqs">
              <div class="req-col">
                <div class="req-hd">Must-have</div>
                <ul>{#each jdHighlights.musts as m}<li>{m}</li>{/each}</ul>
              </div>
              <div class="req-col">
                <div class="req-hd">Nice-to-have</div>
                <ul>{#each jdHighlights.nice as n}<li>{n}</li>{/each}</ul>
              </div>
            </div>
          </div>
        </div>

        <!-- ============ THE COMPANY ============ -->
        <div class="block">
          <div class="block-hd">
            <h2>About {app.co}</h2>
            <span class="ai-tag">AI · refreshed today</span>
          </div>
          <div class="company-row">
            <img class="company-logo" src={`https://logo.clearbit.com/${app.domain}`} alt="" />
            <div>
              <p class="company-blurb">{companyBlurb}</p>
              <div class="company-meta">
                <span>~8,000 employees</span>
                <span class="dot">·</span>
                <span>San Francisco · Remote-friendly</span>
                <span class="dot">·</span>
                <a href={`https://${app.domain}`} target="_blank" rel="noopener">{app.domain} ↗</a>
              </div>
            </div>
          </div>
        </div>

        <!-- ============ HIRING MANAGER ============ -->
        <div class="block">
          <div class="block-hd">
            <h2>Hiring manager</h2>
            <span class="ai-tag">from posting</span>
          </div>
          <div class="hm-card">
            <div class="hm-av">{hiringManager.initials}</div>
            <div class="hm-info">
              <h4>{hiringManager.name}</h4>
              <div class="role">{hiringManager.role}</div>
            </div>
            <a class="hm-li" href={hiringManager.linkedin} target="_blank" rel="noopener">
              <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M3.5 6h2v6h-2zM4.5 3a1 1 0 1 1 0 2 1 1 0 0 1 0-2zM7 6h2v.9c.3-.5.9-1 1.8-1 1.6 0 2.2 1 2.2 2.6V12h-2V9c0-.9-.3-1.4-1.1-1.4-.6 0-1 .4-1 1.2V12H7z"/></svg>
              View on LinkedIn ↗
            </a>
          </div>
        </div>

        <!-- ============ INTERVIEWER BRIEF ============ -->
        <div class="block">
          <div class="block-hd">
            <h2>Your interviewer</h2>
            <span class="ai-tag">AI · web research</span>
            <button class="regen">↻ Regenerate</button>
          </div>

          <div class="interviewer">
            <div class="iv-photo">{interviewer.initials}</div>
            <div class="iv-who">
              <h4>{interviewer.name}</h4>
              <div class="role">{interviewer.role}</div>
              <div class="prior"><b>Prior:</b> {interviewer.prior.join(' · ')}</div>
            </div>
            <a class="iv-li" href={interviewer.linkedin} target="_blank" rel="noopener">
              <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M3.5 6h2v6h-2zM4.5 3a1 1 0 1 1 0 2 1 1 0 0 1 0-2zM7 6h2v.9c.3-.5.9-1 1.8-1 1.6 0 2.2 1 2.2 2.6V12h-2V9c0-.9-.3-1.4-1.1-1.4-.6 0-1 .4-1 1.2V12H7z"/></svg>
              LinkedIn ↗
            </a>
          </div>

          <p class="snapshot">{interviewer.snapshot}</p>

          <div class="signals">
            <div class="sig-hd">Recent signals <span class="num">last 90 days</span></div>
            <ul>
              {#each interviewer.signals as s}
                <li>
                  <span class="date">{s.date}</span>
                  <span class="body">
                    <span class="kind">{s.kind}</span>
                    {s.body}
                    <span class="source">{s.source} ↗</span>
                  </span>
                </li>
              {/each}
            </ul>
          </div>

          <div class="la-grid">
            <div class="la-col lands">
              <h3><span class="glyph">+</span> What lands well</h3>
              <ul>{#each interviewer.lands as l}<li><span class="glyph">+</span><span>{l}</span></li>{/each}</ul>
            </div>
            <div class="la-col avoid">
              <h3><span class="glyph">−</span> What to avoid</h3>
              <ul>{#each interviewer.avoid as a}<li><span class="glyph">−</span><span>{a}</span></li>{/each}</ul>
            </div>
          </div>

          <div class="questions">
            <div class="block-hd" style="margin: 12px 0 8px;">
              <h2 style="font-size:15px">Questions worth asking</h2>
              <span class="ai-tag">ranked by signal</span>
            </div>
            <ol>
              {#each interviewer.questions as q, i}
                <li>
                  <span class="qn">{i + 1}</span>
                  <div>
                    <div class="q">"{q.q}"</div>
                    <div class="why">{q.why}</div>
                  </div>
                  <button class="save" title="Save to prep doc">＋</button>
                </li>
              {/each}
            </ol>
          </div>

          <div class="disclaimer">
            Briefing synthesised from public posts, talks, and papers · always verify before you walk in · last refreshed 12 min ago
          </div>
        </div>
        {/if}

        <p class="footer-link"><a href="/preview/redesign">← back to previews</a></p>
      </div>
    </div>
  </section>
</div>

<style>
  :global(html, body) { background: var(--surface); margin: 0; }
  .frame { display: grid; grid-template-columns: 220px 1fr; min-height: 100vh; font-family: var(--sans); color: var(--ink); }

  /* Sidebar */
  .sidebar { background: var(--surface-2); border-right: 1px solid var(--rule); padding: 18px 14px; display: flex; flex-direction: column; gap: 4px; }
  .brand { display: flex; align-items: center; gap: 8px; padding: 4px 8px 18px; font-weight: 600; font-size: 14.5px; letter-spacing: -0.01em; }
  .brand .mark { width: 16px; height: 16px; border-radius: 4px; background: linear-gradient(135deg, var(--accent), var(--accent-strong)); }
  .nav-item { display: flex; align-items: center; gap: 10px; padding: 7px 10px; border-radius: 6px; font-size: 13px; color: var(--ink-2); cursor: pointer; }
  .nav-item .dot { width: 14px; height: 14px; border-radius: 3px; background: var(--rule-strong); }
  .nav-item.active { background: var(--card); color: var(--ink); box-shadow: var(--sh-1); }
  .sidebar-footer { margin-top: auto; padding-top: 16px; }
  .profile { display: flex; align-items: center; gap: 10px; padding: 8px; }
  .profile .av { width: 28px; height: 28px; border-radius: 50%; object-fit: cover; background: var(--rule-strong); }
  .profile .who { font-size: 13px; line-height: 1.2; }
  .profile .who small { display: block; font-size: 11px; color: var(--mute); }

  /* Topbar */
  .topbar { display: flex; justify-content: space-between; align-items: center; padding: 12px 28px; border-bottom: 1px solid var(--rule); background: var(--surface); }
  .crumb .root { color: var(--mute); cursor: pointer; }
  .crumb .sep { color: var(--mute-2); margin: 0 6px; }
  .crumb .here { font-weight: 500; font-size: 13.5px; }
  .right { display: flex; align-items: center; gap: 8px; }
  .btn { background: var(--card); border: 1px solid var(--rule); border-radius: 7px; padding: 6px 11px; font-size: 12.5px; font-weight: 500; color: var(--ink); cursor: pointer; }
  .btn-danger { color: var(--danger-text); }
  .user-av { width: 30px; height: 30px; border-radius: 50%; object-fit: cover; background: var(--rule-strong); cursor: pointer; margin-left: 8px; border: 1px solid var(--rule); }

  /* Body */
  .body { padding: 28px; }
  .body-inner { max-width: 980px; margin: 0 auto; }

  /* App header with logo */
  .app-hd { display: grid; grid-template-columns: 64px 1fr auto; gap: 18px; align-items: center; margin-bottom: 24px; }
  .logo-big { width: 64px; height: 64px; border-radius: 12px; background: var(--card); object-fit: contain; padding: 8px; border: 1px solid var(--rule); }
  .co-line { display: flex; align-items: center; gap: 12px; }
  .co-line h1 { font-size: 24px; font-weight: 500; margin: 0; letter-spacing: -0.02em; }
  .role-line { font-size: 16px; color: var(--ink-2); margin-top: 2px; }
  .sub { font-size: 12.5px; color: var(--mute); margin-top: 6px; display: flex; flex-wrap: wrap; gap: 0 6px; }
  .sub .dot { color: var(--mute-2); }
  .src-link { display: inline-flex; align-items: center; gap: 6px; background: var(--card); border: 1px solid var(--rule); border-radius: 8px; padding: 8px 12px; font-size: 12.5px; font-weight: 500; color: var(--ink); text-decoration: none; transition: border-color 120ms; }
  .src-link:hover { border-color: var(--accent); color: var(--accent-text); }

  /* Pills */
  .pill { display: inline-flex; align-items: center; gap: 5px; padding: 3px 8px; border-radius: 99px; font-size: 11.5px; font-weight: 500; background: var(--surface-2); color: var(--ink-2); }
  .pill .pdot { width: 5px; height: 5px; border-radius: 50%; background: var(--mute-2); }
  .pill.interview { background: var(--accent-tint); color: var(--accent-text); }
  .pill.interview .pdot { background: var(--accent); }

  /* Tabs */
  .tabs { display: flex; gap: 4px; border-bottom: 1px solid var(--rule); margin-bottom: 24px; }
  .tab { background: transparent; border: 0; padding: 10px 14px; font-size: 13px; color: var(--mute); cursor: pointer; border-bottom: 2px solid transparent; margin-bottom: -1px; font-weight: 500; }
  .tab.active { color: var(--ink); border-bottom-color: var(--ink); }
  .tab .count { font-family: var(--mono); font-size: 10px; background: var(--accent-tint); color: var(--accent-text); padding: 1px 5px; border-radius: 3px; margin-left: 4px; }

  /* Block */
  .block { background: var(--card); border: 1px solid var(--rule); border-radius: 12px; padding: 22px 24px; margin-bottom: 16px; box-shadow: var(--sh-1); }
  .block-hd { display: flex; align-items: center; gap: 10px; margin-bottom: 14px; }
  .block-hd h2 { font-size: 16px; font-weight: 600; margin: 0; letter-spacing: -0.01em; }
  .ai-tag { font-family: var(--mono); font-size: 10px; background: var(--accent-tint); color: var(--accent-text); padding: 2px 7px; border-radius: 4px; letter-spacing: 0.04em; }
  .regen { background: transparent; border: 0; color: var(--accent-text); font-size: 12px; cursor: pointer; margin-left: auto; font-weight: 500; }

  /* Role */
  .role-grid { display: grid; grid-template-columns: 2fr 1fr; gap: 24px; }
  .role-summary p { font-size: 13.5px; line-height: 1.6; color: var(--ink-2); margin: 0 0 14px; }
  .role-meta { display: flex; flex-direction: column; gap: 6px; margin: 0; }
  .role-meta > div { display: grid; grid-template-columns: 80px 1fr; font-size: 12.5px; }
  .role-meta dt { color: var(--mute); margin: 0; }
  .role-meta dd { margin: 0; color: var(--ink); font-feature-settings: "tnum"; }
  .role-reqs { display: flex; flex-direction: column; gap: 14px; }
  .req-hd { font-size: 11px; color: var(--mute); text-transform: uppercase; letter-spacing: 0.06em; margin-bottom: 6px; font-weight: 500; }
  .req-col ul { margin: 0; padding: 0; list-style: none; }
  .req-col li { font-size: 12.5px; padding: 4px 0; color: var(--ink-2); position: relative; padding-left: 14px; }
  .req-col li::before { content: '·'; position: absolute; left: 4px; color: var(--mute-2); }

  /* Company */
  .company-row { display: grid; grid-template-columns: 56px 1fr; gap: 18px; align-items: flex-start; }
  .company-logo { width: 56px; height: 56px; border-radius: 10px; background: var(--surface-2); object-fit: contain; padding: 6px; border: 1px solid var(--rule); }
  .company-blurb { font-size: 13.5px; line-height: 1.6; color: var(--ink-2); margin: 0 0 8px; }
  .company-meta { font-size: 12.5px; color: var(--mute); display: flex; gap: 6px; flex-wrap: wrap; }
  .company-meta a { color: var(--accent-text); text-decoration: none; }
  .dot { color: var(--mute-2); }

  /* Hiring manager */
  .hm-card { display: grid; grid-template-columns: 48px 1fr auto; gap: 14px; align-items: center; padding: 12px; background: var(--surface-2); border-radius: 10px; }
  .hm-av { width: 48px; height: 48px; border-radius: 50%; background: var(--accent-tint); color: var(--accent-text); display: grid; place-items: center; font-weight: 600; font-size: 16px; }
  .hm-info h4 { margin: 0; font-size: 14.5px; font-weight: 500; }
  .hm-info .role { font-size: 12.5px; color: var(--mute); margin-top: 2px; }
  .hm-li { display: inline-flex; align-items: center; gap: 6px; background: var(--card); border: 1px solid var(--rule); border-radius: 7px; padding: 6px 10px; font-size: 12px; font-weight: 500; color: var(--ink); text-decoration: none; }
  .hm-li svg { color: #0a66c2; }

  /* Interviewer */
  .interviewer { display: grid; grid-template-columns: 48px 1fr auto; gap: 14px; align-items: center; margin-bottom: 16px; }
  .iv-photo { width: 48px; height: 48px; border-radius: 50%; background: var(--positive-tint); color: var(--positive-text); display: grid; place-items: center; font-weight: 600; font-size: 16px; }
  .iv-who h4 { margin: 0; font-size: 14.5px; font-weight: 500; }
  .iv-who .role { font-size: 12.5px; color: var(--mute); margin-top: 2px; }
  .iv-who .prior { font-size: 11.5px; color: var(--mute); margin-top: 4px; }
  .iv-li { display: inline-flex; align-items: center; gap: 6px; background: var(--card); border: 1px solid var(--rule); border-radius: 7px; padding: 6px 10px; font-size: 12px; font-weight: 500; color: var(--ink); text-decoration: none; }
  .iv-li svg { color: #0a66c2; }

  .snapshot { font-size: 13.5px; line-height: 1.6; color: var(--ink-2); margin: 0 0 18px; padding: 12px 14px; background: var(--accent-tint); border-radius: 8px; }

  /* Signals */
  .signals { margin-bottom: 20px; }
  .sig-hd { font-size: 13px; font-weight: 500; margin-bottom: 8px; }
  .sig-hd .num { font-family: var(--mono); font-size: 11px; color: var(--mute); margin-left: 6px; font-weight: 400; }
  .signals ul { margin: 0; padding: 0; list-style: none; }
  .signals li { display: grid; grid-template-columns: 60px 1fr; padding: 6px 0; font-size: 12.5px; border-top: 1px dashed var(--rule); }
  .signals li:first-child { border-top: 0; }
  .signals .date { color: var(--mute); font-family: var(--mono); font-size: 11px; padding-top: 2px; }
  .signals .kind { font-family: var(--mono); font-size: 10px; background: var(--surface-2); color: var(--mute); padding: 1px 5px; border-radius: 3px; margin-right: 6px; }
  .signals .source { color: var(--accent-text); margin-left: 4px; font-size: 11.5px; }

  /* Lands / avoid */
  .la-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; margin: 20px 0; }
  .la-col { padding: 14px; border-radius: 10px; }
  .la-col.lands { background: var(--positive-tint); }
  .la-col.avoid { background: var(--danger-tint); }
  .la-col h3 { font-size: 12.5px; font-weight: 600; margin: 0 0 10px; display: flex; align-items: center; gap: 6px; }
  .la-col.lands h3 { color: var(--positive-text); }
  .la-col.avoid h3 { color: var(--danger-text); }
  .la-col h3 .glyph { width: 16px; height: 16px; border-radius: 50%; background: rgba(255,255,255,0.6); display: grid; place-items: center; font-size: 11px; }
  .la-col ul { margin: 0; padding: 0; list-style: none; }
  .la-col li { display: grid; grid-template-columns: 16px 1fr; gap: 8px; font-size: 12.5px; padding: 3px 0; line-height: 1.45; }
  .la-col li .glyph { color: var(--mute-2); font-weight: 600; }
  .la-col.lands li .glyph { color: var(--positive-text); }
  .la-col.avoid li .glyph { color: var(--danger-text); }

  /* Questions */
  .questions { margin-top: 8px; }
  .questions ol { margin: 0; padding: 0; list-style: none; display: flex; flex-direction: column; gap: 1px; background: var(--rule); border: 1px solid var(--rule); border-radius: 10px; overflow: hidden; }
  .questions li { background: var(--card); padding: 12px 14px; display: grid; grid-template-columns: 24px 1fr 28px; gap: 12px; align-items: center; }
  .qn { font-family: var(--mono); font-size: 11px; color: var(--mute); }
  .q { font-size: 13px; font-style: italic; color: var(--ink); }
  .why { font-size: 11.5px; color: var(--mute); margin-top: 3px; }
  .save { background: var(--surface-2); border: 0; width: 24px; height: 24px; border-radius: 6px; cursor: pointer; color: var(--mute); font-size: 14px; }
  .save:hover { background: var(--accent-tint); color: var(--accent-text); }

  .disclaimer { margin-top: 18px; font-size: 11.5px; color: var(--mute); padding-top: 14px; border-top: 1px dashed var(--rule); }

  .footer-link { margin-top: 36px; font-size: 12.5px; }
  .footer-link a { color: var(--accent-text); text-decoration: none; }
</style>
