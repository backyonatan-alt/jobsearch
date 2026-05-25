<script>
  // Static mockup, fully hard-coded. No backend wiring.
  const app = {
    co: 'Stripe',
    domain: 'stripe.com',
    role: 'Staff Engineer, Payments Infrastructure',
    location: 'Remote · US',
    status: 'interview',
    applied: 'May 18',
    daysInPipeline: 7,
    source: 'Referral · Maya Patel',
    cv: 'v3',
    jdUrl: 'https://stripe.com/jobs/listing/staff-engineer-payments/abc123'
  };

  const upNext = {
    when: 'Tomorrow · 2:00 PM',
    countdown: '~18 hours',
    kind: 'Technical screen',
    duration: '60 min',
    medium: 'Google Meet',
    panel: '1:1 with Sarah Chen'
  };

  const stats = [
    { lbl: 'Days in pipeline', n: 7,    sub: 'Applied May 18',    tone: 'accent'   },
    { lbl: 'Current stage',    n: 1,    sub: 'of 4 expected',     tone: 'positive', isFraction: true, dn: 4 },
    { lbl: 'Match score',      n: 'A−', sub: 'CV v3 vs. JD',      tone: 'warm'     }
  ];

  // Compact tags for the role section (single neutral chip style, with
  // a colored dot indicator on the group label).
  const mustHave = ['10+ yrs distributed systems', 'Go or Rust at scale', 'Owned a Tier-0 service'];
  const niceHave = ['Payments / fintech', 'Multi-region active-active', 'OSS maintainership'];
  const techStack = ['Go', 'Rust', 'PostgreSQL', 'Kafka', 'AWS'];

  // Trimmed company copy.
  const company = {
    blurb: 'Payments + financial infrastructure behind a large share of the internet. Recently shipping agentic-commerce APIs.',
    employees: '~8,000',
    hq: 'San Francisco',
    founded: '2010',
    funding: 'Late stage · last round $6.5B at $50B valuation'
  };

  const hiringManager = {
    name: 'Devon Marquez',
    role: 'Director of Engineering, Payments Infra',
    linkedin: 'https://www.linkedin.com/in/example-devon-marquez',
    initials: 'DM'
  };

  // Interviewer brief — shorter prose, more visual chips.
  const interviewer = {
    name: 'Sarah Chen',
    role: 'Staff Engineer · Payments Routing',
    initials: 'SC',
    linkedin: 'https://www.linkedin.com/in/example-sarah-chen',
    prior: ['Cloudflare', 'Two Sigma', 'CMU MS'],
    snapshot: '2 years at Stripe on payments-routing. Known for her QCon SF 2024 talk on regional failover — likely to drill on failure modes more than happy paths.',
    signals: [
      { date: 'Apr 26', kind: 'Talk', body: '"Regional failover at p99.99"',     source: 'qconsf.com',  domain: 'qconsf.com' },
      { date: 'Mar 04', kind: 'Post', body: 'Thread on idempotency-key collisions', source: 'twitter.com', domain: 'twitter.com' },
      { date: 'Jan 12', kind: 'Doc',  body: 'Cited in Stripe multi-region writes', source: 'stripe.com',  domain: 'stripe.com' }
    ],
    lands: [
      'Concrete failure-mode reasoning',
      'Latency, RPS, error-rate numbers',
      'Rollback + observability'
    ],
    avoid: [
      'Hand-waving consistency tradeoffs',
      'Big-bang rewrites with no rollout',
      'Database choice before constraints'
    ],
    questions: [
      { q: "What's the hardest production incident this year?",                 why: 'Real story; reveals what she owns vs. delegates.' },
      { q: 'Hot-patch legacy router vs. forward on v2?',                         why: 'Tests whether the org pays down debt.' },
      { q: 'What does a great first 90 days look like?',                         why: 'Standard, high-signal — loose answer = loose role.' }
    ]
  };

  let tab = $state('brief');
</script>

<svelte:head><title>{app.co} brief — Pursuit</title></svelte:head>

<div class="frame">
  <aside class="sidebar">
    <div class="brand">
      <svg viewBox="0 0 24 24" width="24" height="24" fill="none" class="brand-mark">
        <circle cx="12" cy="12" r="9.5" stroke="currentColor" stroke-width="1.4" opacity="0.65"/>
        <circle cx="12" cy="12" r="5.5" stroke="currentColor" stroke-width="1.4" opacity="0.9"/>
        <circle cx="17.5" cy="6.5" r="2.6" fill="currentColor"/>
      </svg>
      <span class="brand-name">Pursuit</span>
    </div>
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
    <div class="topbar">
      <div class="crumb">
        <span class="root">Applications</span>
        <span class="sep">/</span>
        <span class="here">{app.co}</span>
      </div>
      <div class="right">
        <button class="btn">Update status</button>
        <button class="btn">Edit</button>
        <img class="user-av" src="https://www.gravatar.com/avatar/?d=mp&s=64" alt="" />
      </div>
    </div>

    <div class="body">
      <div class="body-inner">

        <!-- HERO STRIP: company id + status + JD link + quick facts -->
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
          <!-- Quick-facts chip row replaces the old "Applied · via · CV · JD" subline -->
          <div class="facts">
            <span class="fact"><span class="fdot d-loc"></span>{app.location}</span>
            <span class="fact"><span class="fdot d-app"></span>Applied {app.applied}</span>
            <span class="fact"><span class="fdot d-src"></span>{app.source}</span>
            <span class="fact"><span class="fdot d-cv"></span>CV {app.cv}</span>
          </div>
        </div>

        <!-- UP NEXT — prominent next-action card -->
        <div class="upnext">
          <div class="up-left">
            <span class="up-tag"><span class="up-pulse"></span>Up next · {upNext.countdown}</span>
            <h3>{upNext.kind}</h3>
            <div class="up-meta">
              <span>{upNext.when}</span>
              <span class="dot">·</span>
              <span>{upNext.duration}</span>
              <span class="dot">·</span>
              <span>{upNext.medium}</span>
              <span class="dot">·</span>
              <span>{upNext.panel}</span>
            </div>
          </div>
          <div class="up-right">
            <button class="btn-prep">Open prep ↓</button>
          </div>
        </div>

        <!-- AT-A-GLANCE STATS -->
        <div class="stats">
          {#each stats as s}
            <div class={`stat tone-${s.tone}`}>
              <span class="ribbon"></span>
              <div class="stat-lbl">{s.lbl}</div>
              <div class="stat-n">
                {s.n}{#if s.isFraction}<span class="of">/ {s.dn}</span>{/if}
              </div>
              <div class="stat-sub">{s.sub}</div>
            </div>
          {/each}
        </div>

        <!-- TABS -->
        <div class="tabs">
          <button class={`tab ${tab === 'brief' ? 'active' : ''}`} onclick={() => (tab = 'brief')}>
            Brief <span class="t-tag">AI</span>
          </button>
          <button class={`tab ${tab === 'timeline' ? 'active' : ''}`} onclick={() => (tab = 'timeline')}>
            Timeline <span class="t-tag">4</span>
          </button>
          <button class={`tab ${tab === 'notes' ? 'active' : ''}`} onclick={() => (tab = 'notes')}>Notes</button>
          <button class={`tab ${tab === 'files' ? 'active' : ''}`} onclick={() => (tab = 'files')}>Files</button>
        </div>

        {#if tab === 'brief'}
        <!-- TWO-COLUMN GRID: role + company side by side, more efficient -->
        <div class="two-col">
          <!-- ROLE BLOCK -->
          <div class="block">
            <div class="block-hd">
              <h2>The role</h2>
              <span class="ai-tag">AI · from JD</span>
            </div>
            <div class="req-group">
              <div class="req-hd">
                <span class="req-dot d-must"></span>Must-have
              </div>
              <div class="chip-row">
                {#each mustHave as m}<span class="chip">{m}</span>{/each}
              </div>
            </div>
            <div class="req-group">
              <div class="req-hd">
                <span class="req-dot d-nice"></span>Nice-to-have
              </div>
              <div class="chip-row">
                {#each niceHave as n}<span class="chip">{n}</span>{/each}
              </div>
            </div>
            <div class="req-group">
              <div class="req-hd">
                <span class="req-dot d-tech"></span>Tech in the stack
              </div>
              <div class="chip-row">
                {#each techStack as t}<span class="chip">{t}</span>{/each}
              </div>
            </div>
          </div>

          <!-- COMPANY BLOCK -->
          <div class="block">
            <div class="block-hd">
              <h2>About {app.co}</h2>
              <span class="ai-tag">AI · refreshed today</span>
            </div>
            <div class="company-id">
              <img class="company-logo" src={`https://www.google.com/s2/favicons?sz=128&domain=${app.domain}`} alt="" />
              <div class="company-meta-mini">
                <a href={`https://${app.domain}`} target="_blank" rel="noopener">{app.domain} ↗</a>
                <span class="dot">·</span>
                <span>{company.hq}</span>
              </div>
            </div>
            <p class="company-blurb">{company.blurb}</p>
            <div class="meta-grid">
              <div class="meta-cell">
                <div class="m-lbl">Employees</div>
                <div class="m-val">{company.employees}</div>
              </div>
              <div class="meta-cell">
                <div class="m-lbl">Founded</div>
                <div class="m-val">{company.founded}</div>
              </div>
              <div class="meta-cell wide">
                <div class="m-lbl">Funding</div>
                <div class="m-val">{company.funding}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- HIRING MANAGER — small card -->
        <div class="people-row">
          <div class="block person-block">
            <div class="block-hd">
              <h2>Hiring manager</h2>
              <span class="ai-tag">from posting</span>
            </div>
            <div class="person">
              <div class="p-av t-warm">{hiringManager.initials}</div>
              <div class="p-info">
                <h4>{hiringManager.name}</h4>
                <div class="p-role">{hiringManager.role}</div>
              </div>
              <a class="p-li" href={hiringManager.linkedin} target="_blank" rel="noopener">
                <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M3.5 6h2v6h-2zM4.5 3a1 1 0 1 1 0 2 1 1 0 0 1 0-2zM7 6h2v.9c.3-.5.9-1 1.8-1 1.6 0 2.2 1 2.2 2.6V12h-2V9c0-.9-.3-1.4-1.1-1.4-.6 0-1 .4-1 1.2V12H7z"/></svg>
                LinkedIn
              </a>
            </div>
          </div>

          <div class="block person-block">
            <div class="block-hd">
              <h2>Your interviewer</h2>
              <span class="ai-tag">AI · web research</span>
            </div>
            <div class="person">
              <div class="p-av t-accent">{interviewer.initials}</div>
              <div class="p-info">
                <h4>{interviewer.name}</h4>
                <div class="p-role">{interviewer.role}</div>
              </div>
              <a class="p-li" href={interviewer.linkedin} target="_blank" rel="noopener">
                <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M3.5 6h2v6h-2zM4.5 3a1 1 0 1 1 0 2 1 1 0 0 1 0-2zM7 6h2v.9c.3-.5.9-1 1.8-1 1.6 0 2.2 1 2.2 2.6V12h-2V9c0-.9-.3-1.4-1.1-1.4-.6 0-1 .4-1 1.2V12H7z"/></svg>
                LinkedIn
              </a>
            </div>
            <div class="prior-row">
              <span class="p-lbl">Prior</span>
              {#each interviewer.prior as p}<span class="prior-chip">{p}</span>{/each}
            </div>
          </div>
        </div>

        <!-- INTERVIEWER SNAPSHOT (tinted hero) -->
        <div class="snapshot-card">
          <div class="snap-lbl">In one line</div>
          <p>{interviewer.snapshot}</p>
        </div>

        <!-- RECENT POSTS & TALKS — public things the interviewer has said/done
             that you can reference or read before the meeting -->
        <div class="block">
          <div class="block-hd">
            <h2>Recent posts &amp; talks</h2>
            <span class="ai-tag">last 90 days</span>
          </div>
          <div class="signals-row">
            {#each interviewer.signals as s}
              <a class="signal" href="#">
                <img class="sig-logo" src={`https://www.google.com/s2/favicons?sz=128&domain=${s.domain}`} alt="" />
                <div class="sig-meta">
                  <span class="sig-kind">{s.kind}</span>
                  <span class="sig-date">{s.date}</span>
                </div>
                <div class="sig-body">{s.body}</div>
                <div class="sig-src">{s.source} ↗</div>
              </a>
            {/each}
          </div>
        </div>

        <!-- HOW TO APPROACH — was "What lands / What to avoid".
             Now a single block with two clean lists and check / cross icons. -->
        <div class="block">
          <div class="block-hd">
            <h2>How to approach this interview</h2>
            <span class="ai-tag">AI · interviewer-specific</span>
          </div>
          <div class="approach-grid">
            <div class="approach-col">
              <div class="approach-hd ok">
                <span class="approach-glyph ok">
                  <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2.5 6.5l2.5 2.5 4.5-5"/></svg>
                </span>
                Lead with
              </div>
              <ul class="approach-list">{#each interviewer.lands as l}
                <li>
                  <span class="approach-marker ok">
                    <svg width="10" height="10" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2.5 6.5l2.5 2.5 4.5-5"/></svg>
                  </span>
                  <span>{l}</span>
                </li>
              {/each}</ul>
            </div>
            <div class="approach-col">
              <div class="approach-hd no">
                <span class="approach-glyph no">
                  <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M3 3l6 6M9 3l-6 6"/></svg>
                </span>
                Steer clear of
              </div>
              <ul class="approach-list">{#each interviewer.avoid as a}
                <li>
                  <span class="approach-marker no">
                    <svg width="10" height="10" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M3 3l6 6M9 3l-6 6"/></svg>
                  </span>
                  <span>{a}</span>
                </li>
              {/each}</ul>
            </div>
          </div>
        </div>

        <!-- QUESTIONS -->
        <div class="block">
          <div class="block-hd">
            <h2>Questions worth asking</h2>
            <span class="ai-tag">ranked by signal</span>
          </div>
          <ol class="q-list">
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
          Synthesised from public posts, talks, and papers · refreshed 12 min ago · always verify before you walk in
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
  .brand { display: flex; align-items: center; gap: 10px; padding: 4px 8px 18px; color: var(--accent); }
  .brand-name { font-size: 18px; font-weight: 600; letter-spacing: -0.02em; color: var(--ink); }
  .nav-item { display: flex; align-items: center; gap: 10px; padding: 7px 10px; border-radius: 6px; font-size: 13.5px; color: var(--ink-2); cursor: pointer; }
  .nav-item .dot { width: 14px; height: 14px; border-radius: 3px; background: var(--rule-strong); }
  .nav-item.active { background: var(--card); color: var(--ink); box-shadow: var(--sh-1); }
  .sidebar-footer { margin-top: auto; padding-top: 16px; }
  .profile { display: flex; align-items: center; gap: 10px; padding: 8px; }
  .profile .av { width: 28px; height: 28px; border-radius: 50%; object-fit: cover; background: var(--rule-strong); }
  .profile .who { font-size: 13px; line-height: 1.2; }
  .profile .who small { display: block; font-size: 11.5px; color: var(--mute); }

  /* Topbar */
  .topbar { display: flex; justify-content: space-between; align-items: center; padding: 12px 28px; border-bottom: 1px solid var(--rule); background: var(--surface); }
  .crumb .root { color: var(--mute); cursor: pointer; }
  .crumb .sep { color: var(--mute-2); margin: 0 6px; }
  .crumb .here { font-weight: 600; font-size: 14px; }
  .right { display: flex; align-items: center; gap: 8px; }
  .btn { background: var(--card); border: 1px solid var(--rule); border-radius: 7px; padding: 6px 11px; font-size: 13px; font-weight: 500; color: var(--ink); cursor: pointer; }
  .user-av { width: 30px; height: 30px; border-radius: 50%; object-fit: cover; cursor: pointer; margin-left: 8px; border: 1px solid var(--rule); }

  /* Body */
  .body { padding: 28px; }
  .body-inner { max-width: 1100px; margin: 0 auto; }

  /* HERO */
  .hero {
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 18px;
    padding: 22px 24px;
    margin-bottom: 14px;
    box-shadow: var(--sh-1);
  }
  .hero-top { display: grid; grid-template-columns: 64px 1fr auto; gap: 18px; align-items: center; }
  .logo-big { width: 64px; height: 64px; border-radius: 14px; background: var(--card); object-fit: contain; padding: 8px; border: 1px solid var(--rule); }
  .co-row { display: flex; align-items: center; gap: 12px; }
  .co-row h1 { font-size: 28px; font-weight: 600; margin: 0; letter-spacing: -0.025em; }
  .role-line { font-size: 15px; color: var(--ink-2); margin-top: 2px; font-weight: 500; }
  .src-link { display: inline-flex; align-items: center; gap: 6px; background: var(--card); border: 1px solid var(--rule); border-radius: 99px; padding: 8px 14px; font-size: 12.5px; font-weight: 500; color: var(--ink); text-decoration: none; transition: border-color 120ms, color 120ms; }
  .src-link:hover { border-color: var(--accent); color: var(--accent-text); }
  /* Quick-facts chip row */
  .facts { display: flex; flex-wrap: wrap; gap: 6px; margin-top: 14px; padding-top: 14px; border-top: 1px dashed var(--rule); }
  .fact { display: inline-flex; align-items: center; gap: 6px; font-size: 12.5px; color: var(--ink-2); background: var(--card); border: 1px solid var(--rule); padding: 4px 10px; border-radius: 99px; font-weight: 500; }
  .fdot { width: 6px; height: 6px; border-radius: 50%; background: var(--mute-2); }
  .d-loc { background: var(--accent); }
  .d-app { background: var(--positive); }
  .d-src { background: var(--warm); }
  .d-cv  { background: var(--mute-2); }

  /* Pills */
  .pill { display: inline-flex; align-items: center; gap: 5px; padding: 4px 10px; border-radius: 99px; font-size: 12px; font-weight: 500; background: var(--surface-2); color: var(--ink-2); }
  .pill .pdot { width: 6px; height: 6px; border-radius: 50%; background: var(--mute-2); }
  .pill.interview { background: var(--accent-tint); color: var(--accent-text); }
  .pill.interview .pdot { background: var(--accent); }

  /* UP NEXT */
  .upnext {
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 14px;
    padding: 18px 22px;
    margin-bottom: 14px;
    display: grid; grid-template-columns: 1fr auto; gap: 16px; align-items: center;
    box-shadow: var(--sh-1);
    position: relative;
    overflow: hidden;
  }
  .upnext::before { content: ''; position: absolute; left: 0; top: 0; bottom: 0; width: 4px; background: var(--accent); }
  .up-tag { display: inline-flex; align-items: center; gap: 6px; font-size: 12px; color: var(--accent-text); background: var(--accent-tint); padding: 4px 10px; border-radius: 99px; font-weight: 500; margin-bottom: 8px; }
  @keyframes up-pulse {
    0%, 100% { box-shadow: 0 0 0 0 var(--accent); }
    50%      { box-shadow: 0 0 0 5px transparent; }
  }
  .up-pulse { width: 6px; height: 6px; border-radius: 50%; background: var(--accent); animation: up-pulse 1.6s ease-in-out infinite; }
  .upnext h3 { font-size: 18px; font-weight: 600; margin: 0; letter-spacing: -0.015em; }
  .up-meta { margin-top: 4px; font-size: 13px; color: var(--mute); display: flex; flex-wrap: wrap; gap: 0 6px; }
  .up-meta .dot { color: var(--mute-2); }
  .btn-prep {
    background: var(--accent); color: white; border: 0; border-radius: 99px;
    padding: 10px 18px; font-size: 13.5px; font-weight: 600; cursor: pointer;
    transition: transform 120ms ease;
  }
  .btn-prep:hover { transform: translateY(-1px); }

  /* AT-A-GLANCE STATS */
  .stats { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; margin-bottom: 24px; }
  .stat { position: relative; background: var(--card); border: 1px solid var(--rule); border-radius: 14px; padding: 16px 18px; overflow: hidden; box-shadow: var(--sh-1); }
  .stat .ribbon { position: absolute; top: 0; left: 0; right: 0; height: 3px; }
  .stat.tone-accent   .ribbon { background: var(--accent); }
  .stat.tone-positive .ribbon { background: var(--positive); }
  .stat.tone-warm     .ribbon { background: var(--warm); }
  .stat-lbl { font-size: 12.5px; color: var(--mute); margin-bottom: 4px; font-weight: 500; }
  .stat-n { font-size: 32px; font-weight: 600; letter-spacing: -0.035em; line-height: 1.1; font-feature-settings: "tnum"; }
  .stat.tone-accent   .stat-n { color: var(--accent-text); }
  .stat.tone-positive .stat-n { color: var(--positive-text); }
  .stat.tone-warm     .stat-n { color: var(--warm-text); }
  .of { font-size: 16px; color: var(--mute); margin-left: 6px; font-weight: 500; }
  .stat-sub { font-size: 12px; color: var(--mute); margin-top: 6px; padding-top: 6px; border-top: 1px dashed var(--rule); }

  /* TABS */
  .tabs { display: flex; gap: 4px; border-bottom: 1px solid var(--rule); margin-bottom: 18px; }
  .tab { background: transparent; border: 0; padding: 10px 14px; font-size: 13.5px; color: var(--mute); cursor: pointer; border-bottom: 2px solid transparent; margin-bottom: -1px; font-weight: 600; }
  .tab.active { color: var(--ink); border-bottom-color: var(--ink); }
  .t-tag { font-size: 11px; background: var(--accent-tint); color: var(--accent-text); padding: 1px 7px; border-radius: 99px; margin-left: 4px; font-weight: 500; }

  /* TWO-COLUMN */
  .two-col { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; margin-bottom: 14px; }

  /* BLOCK */
  .block { background: var(--card); border: 1px solid var(--rule); border-radius: 14px; padding: 20px 22px; margin-bottom: 14px; box-shadow: var(--sh-1); }
  .block-hd { display: flex; align-items: center; gap: 10px; margin-bottom: 14px; }
  .block-hd h2 { font-size: 16px; font-weight: 600; margin: 0; letter-spacing: -0.015em; }
  .ai-tag { display: inline-flex; align-items: center; gap: 5px; font-size: 12px; background: var(--accent-tint); color: var(--accent-text); padding: 3px 10px; border-radius: 99px; font-weight: 500; }

  /* ROLE — single neutral chip style; group label carries the only color. */
  .req-group { margin-bottom: 16px; }
  .req-group:last-child { margin-bottom: 0; }
  .req-hd { display: flex; align-items: center; gap: 8px; font-size: 12.5px; color: var(--ink-2); margin-bottom: 8px; font-weight: 600; }
  .req-dot { width: 8px; height: 8px; border-radius: 50%; background: var(--mute-2); }
  .d-must { background: var(--positive); }
  .d-nice { background: var(--accent); }
  .d-tech { background: var(--warm); }
  .chip-row { display: flex; flex-wrap: wrap; gap: 6px; }
  .chip { font-size: 12.5px; padding: 5px 11px; border-radius: 99px; font-weight: 500; background: var(--surface-2); color: var(--ink-2); border: 1px solid var(--rule); }

  /* COMPANY */
  .company-id { display: flex; align-items: center; gap: 12px; margin-bottom: 12px; }
  .company-logo { width: 36px; height: 36px; border-radius: 8px; background: var(--surface-2); object-fit: contain; padding: 4px; border: 1px solid var(--rule); }
  .company-meta-mini { font-size: 12.5px; color: var(--mute); display: flex; gap: 6px; }
  .company-meta-mini a { color: var(--accent-text); text-decoration: none; }
  .dot { color: var(--mute-2); }
  .company-blurb { font-size: 13.5px; line-height: 1.55; color: var(--ink-2); margin: 0 0 14px; }
  .meta-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; }
  .meta-cell { background: var(--surface-2); border-radius: 10px; padding: 10px 12px; }
  .meta-cell.wide { grid-column: span 2; }
  .m-lbl { font-size: 11px; color: var(--mute); font-weight: 500; }
  .m-val { font-size: 13px; color: var(--ink); font-weight: 600; margin-top: 2px; }

  /* PEOPLE CARDS */
  .people-row { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }
  .person-block { padding: 18px 20px; margin-bottom: 0; }
  .person { display: grid; grid-template-columns: 46px 1fr auto; gap: 12px; align-items: center; }
  .p-av { width: 46px; height: 46px; border-radius: 50%; display: grid; place-items: center; font-weight: 600; font-size: 16px; }
  .p-av.t-warm { background: var(--warm-tint); color: var(--warm-text); }
  .p-av.t-accent { background: var(--accent-tint); color: var(--accent-text); }
  .p-info h4 { margin: 0; font-size: 14.5px; font-weight: 600; letter-spacing: -0.01em; }
  .p-info .p-role { font-size: 12.5px; color: var(--mute); margin-top: 2px; }
  .p-li { display: inline-flex; align-items: center; gap: 6px; background: var(--surface-2); border: 1px solid var(--rule); border-radius: 99px; padding: 6px 12px; font-size: 12px; font-weight: 600; color: var(--ink); text-decoration: none; }
  .p-li svg { color: #0a66c2; }
  .prior-row { display: flex; align-items: center; gap: 8px; margin-top: 14px; padding-top: 12px; border-top: 1px dashed var(--rule); flex-wrap: wrap; }
  .p-lbl { font-size: 11.5px; color: var(--mute); font-weight: 600; }
  .prior-chip { font-size: 11.5px; background: var(--surface-2); color: var(--ink-2); padding: 3px 9px; border-radius: 99px; font-weight: 500; }

  /* SNAPSHOT — flat accent-tint card, no gradient */
  .snapshot-card {
    background: var(--accent-tint);
    border: 1px solid var(--rule);
    border-radius: 14px;
    padding: 16px 20px;
    margin: 14px 0;
  }
  .snap-lbl { font-size: 11.5px; color: var(--accent-text); font-weight: 600; margin-bottom: 4px; }
  .snapshot-card p { margin: 0; font-size: 14.5px; line-height: 1.55; color: var(--ink); }

  /* RECENT POSTS & TALKS — three small cards with the source's logo */
  .signals-row { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; }
  .signal {
    background: var(--surface-2);
    border: 1px solid var(--rule);
    border-radius: 12px;
    padding: 12px 14px;
    display: grid;
    grid-template-columns: 28px 1fr;
    grid-template-rows: auto auto auto;
    grid-template-areas:
      "logo meta"
      "body body"
      "src  src";
    gap: 6px 10px;
    text-decoration: none; color: inherit;
    transition: border-color 120ms ease, transform 120ms ease;
  }
  .signal:hover { border-color: var(--rule-strong); transform: translateY(-1px); }
  .sig-logo { grid-area: logo; width: 28px; height: 28px; border-radius: 7px; background: var(--card); object-fit: contain; padding: 3px; border: 1px solid var(--rule); }
  .sig-meta { grid-area: meta; display: flex; align-items: center; gap: 8px; }
  .sig-kind { font-size: 11px; color: var(--mute); font-weight: 600; }
  .sig-date { font-size: 11px; color: var(--mute-2); }
  .sig-body { grid-area: body; font-size: 13px; color: var(--ink); line-height: 1.4; }
  .sig-src  { grid-area: src; font-size: 11.5px; color: var(--accent-text); }

  /* HOW TO APPROACH — two simple lists, no heavy color blocks. */
  .approach-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 22px; }
  .approach-hd { display: flex; align-items: center; gap: 8px; font-size: 13px; font-weight: 600; margin-bottom: 10px; color: var(--ink); }
  .approach-glyph { width: 22px; height: 22px; border-radius: 50%; display: grid; place-items: center; }
  .approach-glyph.ok { background: var(--positive-tint); color: var(--positive-text); }
  .approach-glyph.no { background: var(--danger-tint); color: var(--danger-text); }
  .approach-list { margin: 0; padding: 0; list-style: none; display: flex; flex-direction: column; gap: 8px; }
  .approach-list li { display: grid; grid-template-columns: 18px 1fr; gap: 10px; align-items: flex-start; font-size: 13px; color: var(--ink-2); line-height: 1.5; }
  .approach-marker { width: 18px; height: 18px; border-radius: 50%; display: grid; place-items: center; margin-top: 1px; }
  .approach-marker.ok { background: var(--positive-tint); color: var(--positive-text); }
  .approach-marker.no { background: var(--danger-tint); color: var(--danger-text); }

  /* QUESTIONS */
  .q-list { margin: 0; padding: 0; list-style: none; display: flex; flex-direction: column; gap: 8px; }
  .q-list li { background: var(--surface-2); border-radius: 12px; padding: 12px 14px; display: grid; grid-template-columns: 28px 1fr 28px; gap: 12px; align-items: center; }
  .qn { width: 24px; height: 24px; border-radius: 50%; background: var(--accent-tint); color: var(--accent-text); display: grid; place-items: center; font-size: 12px; font-weight: 600; font-feature-settings: "tnum"; }
  .q { font-size: 13.5px; color: var(--ink); font-weight: 500; }
  .why { font-size: 11.5px; color: var(--mute); margin-top: 3px; }
  .save { background: var(--card); border: 1px solid var(--rule); width: 26px; height: 26px; border-radius: 8px; cursor: pointer; color: var(--mute); font-size: 14px; }
  .save:hover { background: var(--accent-tint); color: var(--accent-text); border-color: var(--accent-tint); }

  .disclaimer { margin-top: 18px; font-size: 11.5px; color: var(--mute); padding-top: 14px; border-top: 1px dashed var(--rule); }

  .footer-link { margin-top: 30px; font-size: 13px; }
  .footer-link a { color: var(--accent-text); text-decoration: none; }
</style>
