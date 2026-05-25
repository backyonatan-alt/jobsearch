<script>
  // Variant B — "Warm playful".
  // Sans throughout but rounded + softer weights, no UPPERCASE, no mono.
  // Icons: filled custom glyphs sitting in tinted circles, each in its own hue.
  // Brand: dynamic upward path + dot mark with gradient.
  const me = { name: 'Yonatan', email: 'back.yonatan@gmail.com' };
  const counts = { interviews: 2, offers: 1, applied: 8, wishlist: 3 };

  const actions = [
    { urgency: 'Today',         tone: 'accent', title: 'Prep for your Stripe loop',     sub: 'Technical screen with Sarah Chen — tomorrow 2:00 PM',     cta: 'Open brief',     logo: 'stripe.com' },
    { urgency: '2 days left',   tone: 'warm',   title: 'Decide on the Vercel offer',    sub: 'They asked for an answer by Friday — last touch 3 days ago.', cta: 'Open offer',  logo: 'vercel.com' },
    { urgency: 'Quiet 7 days',  tone: 'mute',   title: 'Nudge Linear',                  sub: 'Applied via referral — no response yet.',                  cta: 'Draft follow-up', logo: 'linear.app' },
    { urgency: 'Just moved',    tone: 'positive', title: 'Learn about Anthropic',         sub: 'Moved to screen yesterday — generate the company brief.',  cta: 'Generate',       logo: 'anthropic.com' }
  ];

  const insights = [
    { kind: 'trend', tone: 'positive', text: 'Referrals convert at <b>3×</b> the rate of cold apps.', detail: '4 of 5 referrals reached screen vs. 3 of 10 cold' },
    { kind: 'pace',  tone: 'warm',     text: "You haven't applied in <b>5 days</b>.",                  detail: 'Your usual pace is about two per day' },
    { kind: 'bell',  tone: 'accent',   text: '<b>3 loops</b> have gone quiet for over a week.',        detail: 'Linear, Notion, Figma — worth a nudge' }
  ];

  const apps = [
    { co: 'Stripe',    role: 'Staff Eng, Payments',     status: 'interview', applied: '4 days ago',  domain: 'stripe.com',    stale: false },
    { co: 'Vercel',    role: 'Senior PM, Edge',         status: 'offer',     applied: '18 days ago', domain: 'vercel.com',    stale: false },
    { co: 'Anthropic', role: 'Research Engineer',       status: 'screen',    applied: '6 days ago',  domain: 'anthropic.com', stale: false },
    { co: 'Linear',    role: 'Senior Frontend Eng',     status: 'applied',   applied: '7 days ago',  domain: 'linear.app',    stale: true  },
    { co: 'Notion',    role: 'Eng Manager, Editor',     status: 'applied',   applied: '9 days ago',  domain: 'notion.so',     stale: true  },
    { co: 'Supabase',  role: 'Developer Advocate',      status: 'applied',   applied: '3 days ago',  domain: 'supabase.com',  stale: false },
    { co: 'Figma',     role: 'Design Eng',              status: 'applied',   applied: '11 days ago', domain: 'figma.com',     stale: true  }
  ];

  const STATUS_LABEL = { wishlist:'Wishlist', applied:'Applied', screen:'Screen', interview:'Interview', offer:'Offer' };

  const today = new Date();
  const dateLong = today.toLocaleDateString('en-US', { weekday: 'long', day: 'numeric', month: 'long' });
</script>

<svelte:head><title>Today · variant B — Pursuit</title></svelte:head>

<div class="frame">
  <aside class="sidebar">
    <div class="brand">
      <svg class="brand-mark" viewBox="0 0 28 28" width="26" height="26" fill="none">
        <defs>
          <linearGradient id="pursuit-grad-b" x1="0" y1="28" x2="28" y2="0">
            <stop offset="0" stop-color="oklch(0.62 0.19 258)"/>
            <stop offset="1" stop-color="oklch(0.7 0.18 318)"/>
          </linearGradient>
        </defs>
        <path d="M4 19 L11 12 L16 16 L24 6"
              stroke="url(#pursuit-grad-b)"
              stroke-width="2.4"
              stroke-linecap="round"
              stroke-linejoin="round"/>
        <circle cx="24" cy="6" r="3" fill="url(#pursuit-grad-b)"/>
      </svg>
      <span class="brand-name">Pursuit</span>
    </div>
    <a class="nav-item active"><span class="dot"></span>Today</a>
    <a class="nav-item"><span class="dot"></span>Board</a>
    <a class="nav-item"><span class="dot"></span>Funnel</a>
    <div class="sidebar-footer">
      <div class="profile">
        <img class="av" src="https://www.gravatar.com/avatar/?d=mp&s=64" alt="" />
        <div class="who">{me.name}<small>{me.email}</small></div>
      </div>
    </div>
  </aside>

  <section class="main">
    <div class="topbar">
      <div class="crumb"><span class="here">Today</span></div>
      <div class="right">
        <div class="search">
          <svg width="14" height="14" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="7" cy="7" r="4.5"/><path d="M10.5 10.5L13 13"/></svg>
          <span>Search applications, people</span>
        </div>
        <button class="btn">Import</button>
        <button class="btn btn-primary">New application</button>
        <img class="user-av" src="https://www.gravatar.com/avatar/?d=mp&s=64" alt={me.name} title={me.email} />
      </div>
    </div>

    <div class="body">
      <div class="body-inner">

        <!-- Greeting block with soft tinted backdrop -->
        <div class="hello-card">
          <div class="hello-text">
            <div class="date"><span class="ddot"></span>{dateLong}</div>
            <h1>Good afternoon, {me.name}.</h1>
            <p class="quip">Two interviews on your plate this week — and a Vercel offer waiting on you.</p>
          </div>
        </div>

        <!-- Big counts, each with its own colored accent ribbon -->
        <div class="counts">
          <div class="count-cell live">
            <span class="ribbon"></span>
            <div class="n">{counts.interviews}</div>
            <div class="lbl">Interviews <span class="hint" title="An interview loop = a series of rounds with one company">?</span></div>
          </div>
          <div class="count-cell warm">
            <span class="ribbon"></span>
            <div class="n">{counts.offers}</div>
            <div class="lbl">Open offers</div>
          </div>
          <div class="count-cell positive">
            <span class="ribbon"></span>
            <div class="n">{counts.applied}</div>
            <div class="lbl">Applied &amp; waiting</div>
          </div>
          <div class="count-cell mute">
            <span class="ribbon"></span>
            <div class="n">{counts.wishlist}</div>
            <div class="lbl">Wishlist</div>
          </div>
        </div>

        <!-- What you can do today -->
        <div class="section-hd">
          <h2>What you can do today</h2>
          <span class="ai-tag">
            <svg width="11" height="11" viewBox="0 0 12 12" fill="currentColor"><path d="M6 0l1.2 3.6L11 5l-3.8 1.4L6 10 4.8 6.4 1 5l3.8-1.4z"/></svg>
            <span>AI suggested</span>
          </span>
        </div>
        <div class="action-grid">
          {#each actions as a}
            <div class={`action-card ${a.tone}`}>
              <div class="action-top">
                <img class="action-logo" src={`https://logo.clearbit.com/${a.logo}`} alt="" />
                <span class={`urgency u-${a.tone}`}>{a.urgency}</span>
              </div>
              <h3>{a.title}</h3>
              <p>{a.sub}</p>
              <button class="action-cta">{a.cta} <span class="arrow">→</span></button>
            </div>
          {/each}
        </div>

        <!-- What we're noticing — each row in its own hue -->
        <div class="section-hd">
          <h2>What we're noticing</h2>
          <span class="ai-tag">
            <svg width="11" height="11" viewBox="0 0 12 12" fill="currentColor"><path d="M6 0l1.2 3.6L11 5l-3.8 1.4L6 10 4.8 6.4 1 5l3.8-1.4z"/></svg>
            <span>This week</span>
          </span>
        </div>
        <div class="insight-list">
          {#each insights as ins}
            <div class={`insight ${ins.tone}`}>
              <span class={`ins-icon i-${ins.tone}`}>
                {#if ins.kind === 'trend'}
                  <svg width="22" height="22" viewBox="0 0 24 24" fill="none">
                    <path d="M3 18l5-5 4 3 8-10" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
                    <circle cx="20" cy="6" r="3.2" fill="currentColor"/>
                  </svg>
                {:else if ins.kind === 'pace'}
                  <svg width="22" height="22" viewBox="0 0 24 24" fill="none">
                    <circle cx="12" cy="13" r="8" fill="currentColor" opacity="0.18"/>
                    <circle cx="12" cy="13" r="8" stroke="currentColor" stroke-width="2" fill="none"/>
                    <path d="M12 8v5l3 2" stroke="currentColor" stroke-width="2" stroke-linecap="round" fill="none"/>
                    <path d="M9 3h6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                  </svg>
                {:else if ins.kind === 'bell'}
                  <svg width="22" height="22" viewBox="0 0 24 24" fill="none">
                    <path d="M6 16V11a6 6 0 0 1 12 0v5l1.5 2H4.5z" fill="currentColor" opacity="0.18"/>
                    <path d="M6 16V11a6 6 0 0 1 12 0v5l1.5 2H4.5z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
                    <path d="M10 20a2 2 0 0 0 4 0" stroke="currentColor" stroke-width="2" stroke-linecap="round" fill="none"/>
                    <circle cx="18" cy="6" r="3" fill="currentColor"/>
                  </svg>
                {/if}
              </span>
              <div class="ins-body">
                <div class="ins-line">{@html ins.text}</div>
                <div class="ins-detail">{ins.detail}</div>
              </div>
              <button class="ins-act">View →</button>
            </div>
          {/each}
        </div>

        <!-- Applications -->
        <div class="section-hd">
          <h2>Applications <span class="count">{apps.length}</span></h2>
          <div class="filters">
            <button class="chip active">Active</button>
            <button class="chip">Interview</button>
            <button class="chip">Offer</button>
            <button class="chip">Wishlist</button>
            <button class="chip">All</button>
          </div>
        </div>
        <div class="table">
          <div class="tr head">
            <span>Company</span><span>Role</span><span>Status</span><span>Applied</span><span></span>
          </div>
          {#each apps as a}
            <div class={`tr ${a.stale ? 'stale' : ''}`}>
              <span class="co">
                <img class="logo" src={`https://logo.clearbit.com/${a.domain}`} alt="" />
                <span>{a.co}</span>
              </span>
              <span class="role">{a.role}</span>
              <span class={`pill ${a.status}`}><span class="pdot"></span>{STATUS_LABEL[a.status]}</span>
              <span class="applied">{a.applied} {#if a.stale}<span class="stale-tag">stale</span>{/if}</span>
              <span class="arrow">→</span>
            </div>
          {/each}
        </div>

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
  .brand { display: flex; align-items: center; gap: 10px; padding: 4px 8px 18px; }
  .brand-name { font-size: 18px; font-weight: 600; letter-spacing: -0.02em; color: var(--ink); }
  .nav-item { display: flex; align-items: center; gap: 10px; padding: 8px 10px; border-radius: 8px; font-size: 13.5px; color: var(--ink-2); cursor: pointer; }
  .nav-item .dot { width: 14px; height: 14px; border-radius: 4px; background: var(--rule-strong); }
  .nav-item.active { background: var(--card); color: var(--ink); box-shadow: var(--sh-1); }
  .sidebar-footer { margin-top: auto; padding-top: 16px; }
  .profile { display: flex; align-items: center; gap: 10px; padding: 8px; }
  .profile .av { width: 28px; height: 28px; border-radius: 50%; object-fit: cover; background: var(--rule-strong); }
  .profile .who { font-size: 13px; line-height: 1.2; }
  .profile .who small { display: block; font-size: 11.5px; color: var(--mute); }

  /* Topbar */
  .topbar { display: flex; justify-content: space-between; align-items: center; padding: 12px 28px; border-bottom: 1px solid var(--rule); background: var(--surface); }
  .crumb .here { font-weight: 600; font-size: 14px; }
  .right { display: flex; align-items: center; gap: 8px; }
  .search { display: flex; align-items: center; gap: 6px; background: var(--card); border: 1px solid var(--rule); border-radius: 99px; padding: 6px 12px; font-size: 13px; color: var(--mute); min-width: 280px; }
  .btn { background: var(--card); border: 1px solid var(--rule); border-radius: 99px; padding: 7px 13px; font-size: 13px; font-weight: 500; color: var(--ink); cursor: pointer; }
  .btn-primary {
    background: linear-gradient(135deg, oklch(0.62 0.19 258), oklch(0.7 0.18 318));
    border-color: transparent; color: white;
  }
  .user-av { width: 30px; height: 30px; border-radius: 50%; object-fit: cover; cursor: pointer; margin-left: 8px; border: 1px solid var(--rule); }

  .body { padding: 28px; }
  .body-inner { max-width: 1080px; margin: 0 auto; }

  /* Greeting card with very soft tinted gradient */
  .hello-card {
    background: linear-gradient(135deg, var(--accent-tint) 0%, oklch(0.97 0.04 318) 100%);
    border: 1px solid var(--rule);
    border-radius: 18px;
    padding: 26px 28px;
    margin-bottom: 24px;
  }
  .hello-text .date {
    display: inline-flex; align-items: center; gap: 8px;
    font-size: 13px; color: var(--accent-text); margin-bottom: 6px; font-weight: 500;
  }
  .hello-text .ddot { width: 7px; height: 7px; border-radius: 50%; background: var(--accent); display: inline-block; }
  .hello-text h1 {
    font-size: 30px; font-weight: 600; margin: 0; letter-spacing: -0.025em;
    color: var(--ink);
  }
  .quip { margin: 8px 0 0; font-size: 14px; color: var(--ink-2); line-height: 1.5; }

  /* Big counts with colored ribbon on top */
  .counts { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; margin-bottom: 36px; }
  .count-cell {
    position: relative;
    background: var(--card); border: 1px solid var(--rule); border-radius: 16px;
    padding: 22px 22px 18px;
    box-shadow: var(--sh-1);
    overflow: hidden;
  }
  .ribbon { position: absolute; top: 0; left: 0; right: 0; height: 4px; }
  .count-cell.live    .ribbon { background: var(--accent); }
  .count-cell.warm    .ribbon { background: var(--warm); }
  .count-cell.positive .ribbon { background: var(--positive); }
  .count-cell.mute    .ribbon { background: var(--rule-strong); }
  .count-cell .n { font-size: 44px; font-weight: 600; letter-spacing: -0.04em; line-height: 1; color: var(--ink); font-feature-settings: "tnum"; }
  .count-cell .lbl { font-size: 13.5px; color: var(--mute); margin-top: 10px; display: flex; align-items: center; gap: 6px; font-weight: 500; }
  .hint { width: 14px; height: 14px; display: inline-grid; place-items: center; border-radius: 50%; background: var(--surface-2); color: var(--mute); font-size: 10px; font-weight: 600; cursor: help; }

  /* Section heading */
  .section-hd { display: flex; align-items: center; justify-content: space-between; margin: 32px 0 14px; gap: 12px; }
  .section-hd h2 { font-size: 17px; font-weight: 600; margin: 0; letter-spacing: -0.015em; }
  .section-hd h2 .count { font-size: 14px; color: var(--mute); margin-left: 6px; font-weight: 400; }
  .ai-tag {
    display: inline-flex; align-items: center; gap: 6px;
    font-size: 12px; color: var(--accent-text);
    background: var(--accent-tint); padding: 4px 11px; border-radius: 99px;
    font-weight: 500;
  }

  /* Action grid */
  .action-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(260px, 1fr)); gap: 14px; }
  .action-card { background: var(--card); border: 1px solid var(--rule); border-radius: 16px; padding: 18px; display: flex; flex-direction: column; gap: 8px; transition: transform 160ms ease, border-color 160ms ease, box-shadow 160ms ease; cursor: pointer; }
  .action-card:hover { transform: translateY(-3px); border-color: var(--accent); box-shadow: var(--sh-pop); }
  .action-top { display: flex; align-items: center; justify-content: space-between; }
  .action-logo { width: 32px; height: 32px; border-radius: 9px; background: var(--surface-2); object-fit: contain; padding: 3px; }
  .urgency { font-size: 12px; padding: 4px 11px; border-radius: 99px; font-weight: 500; }
  .urgency.u-accent   { background: var(--accent-tint);   color: var(--accent-text); }
  .urgency.u-warm     { background: var(--warm-tint);     color: var(--warm-text); }
  .urgency.u-positive { background: var(--positive-tint); color: var(--positive-text); }
  .urgency.u-mute     { background: var(--surface-2);     color: var(--mute); }
  .action-card h3 { font-size: 15.5px; font-weight: 600; margin: 4px 0 0; letter-spacing: -0.01em; }
  .action-card p { font-size: 13px; color: var(--mute); margin: 0; line-height: 1.5; }
  .action-cta { background: transparent; border: 0; color: var(--accent-text); font-size: 13px; font-weight: 600; text-align: left; padding: 4px 0 0; cursor: pointer; align-self: flex-start; }
  .action-cta .arrow { transition: transform 160ms ease; display: inline-block; }
  .action-card:hover .action-cta .arrow { transform: translateX(3px); }

  /* Insights — each row tinted with its own hue */
  .insight-list { display: flex; flex-direction: column; gap: 10px; }
  .insight { background: var(--card); border: 1px solid var(--rule); border-radius: 14px; padding: 14px 18px; display: grid; grid-template-columns: 44px 1fr auto; gap: 14px; align-items: center; box-shadow: var(--sh-1); }
  .ins-icon { width: 44px; height: 44px; border-radius: 50%; display: grid; place-items: center; }
  .ins-icon.i-positive { background: var(--positive-tint); color: var(--positive-text); }
  .ins-icon.i-warm     { background: var(--warm-tint);     color: var(--warm-text); }
  .ins-icon.i-accent   { background: var(--accent-tint);   color: var(--accent-text); }
  .ins-line { font-size: 14.5px; color: var(--ink); line-height: 1.45; }
  .ins-detail { font-size: 12.5px; color: var(--mute); margin-top: 3px; }
  .ins-act { background: transparent; border: 0; color: var(--accent-text); font-size: 13px; font-weight: 600; cursor: pointer; }

  /* Pipeline table */
  .filters { display: flex; gap: 6px; }
  .chip { background: transparent; border: 1px solid var(--rule); border-radius: 99px; padding: 5px 13px; font-size: 12.5px; color: var(--mute); cursor: pointer; font-weight: 500; }
  .chip.active { background: var(--ink); border-color: var(--ink); color: white; }
  .table { background: var(--card); border: 1px solid var(--rule); border-radius: 16px; overflow: hidden; box-shadow: var(--sh-1); }
  .tr { display: grid; grid-template-columns: 220px 1fr 130px 140px 20px; align-items: center; padding: 12px 20px; border-bottom: 1px solid var(--rule); font-size: 13.5px; cursor: pointer; }
  .tr:last-child { border-bottom: 0; }
  .tr:hover { background: var(--surface-2); }
  .tr.head { background: var(--surface-2); font-size: 12px; color: var(--mute); cursor: default; padding: 10px 20px; font-weight: 500; }
  .tr.head:hover { background: var(--surface-2); }
  .tr.stale .applied { color: var(--danger-text); }
  .co { display: flex; align-items: center; gap: 10px; font-weight: 500; }
  .co .logo { width: 24px; height: 24px; border-radius: 6px; background: var(--surface-2); object-fit: contain; padding: 2px; }
  .role { color: var(--mute); }
  .pill { display: inline-flex; align-items: center; gap: 5px; padding: 3px 10px; border-radius: 99px; font-size: 12px; font-weight: 500; background: var(--surface-2); color: var(--ink-2); width: max-content; }
  .pill .pdot { width: 5px; height: 5px; border-radius: 50%; background: var(--mute-2); }
  .pill.interview { background: var(--accent-tint); color: var(--accent-text); }
  .pill.interview .pdot { background: var(--accent); }
  .pill.offer { background: var(--warm-tint); color: var(--warm-text); }
  .pill.offer .pdot { background: var(--warm); }
  .pill.screen { background: var(--positive-tint); color: var(--positive-text); }
  .pill.screen .pdot { background: var(--positive); }
  .applied { color: var(--mute); }
  .stale-tag { font-size: 11px; background: var(--danger-tint); color: var(--danger-text); padding: 1px 8px; border-radius: 99px; margin-left: 4px; font-weight: 500; }
  .arrow { color: var(--mute-2); }

  .footer-link { margin-top: 36px; font-size: 13px; }
  .footer-link a { color: var(--accent-text); text-decoration: none; }
</style>
