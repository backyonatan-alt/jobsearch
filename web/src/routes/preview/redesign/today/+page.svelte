<script>
  // Static mockup — fully hard-coded data to demonstrate the design.
  // No backend, no real apps used. The fake companies were chosen so
  // Clearbit's logo CDN returns a real image for each one.
  const me = {
    name: 'Yonatan',
    email: 'back.yonatan@gmail.com',
    picture: 'https://lh3.googleusercontent.com/a/default-user=s96-c'
  };

  const counts = { interviews: 2, offers: 1, applied: 8, wishlist: 3 };

  // "What you can do today" — proactive cards driven by pipeline state.
  const actions = [
    {
      kind: 'prep',
      urgency: 'today',
      title: 'Prep for your Stripe loop',
      sub: 'Technical screen with Sarah Chen · tomorrow 2:00 PM',
      cta: 'Open brief',
      logo: 'stripe.com',
      color: 'accent'
    },
    {
      kind: 'decide',
      urgency: '2 days left',
      title: 'Decide on the Vercel offer',
      sub: 'They asked for an answer by Friday — last touch was 3 days ago.',
      cta: 'Open offer',
      logo: 'vercel.com',
      color: 'warm'
    },
    {
      kind: 'follow-up',
      urgency: '7d quiet',
      title: 'Nudge Linear',
      sub: 'Applied 7 days ago via Referral — no response yet.',
      cta: 'Draft follow-up',
      logo: 'linear.app',
      color: 'mute'
    },
    {
      kind: 'research',
      urgency: 'new',
      title: 'Learn about Anthropic',
      sub: 'Moved to Screen yesterday — generate the company brief.',
      cta: 'Generate',
      logo: 'anthropic.com',
      color: 'accent'
    }
  ];

  // "What we're noticing" — insights moved from /funnel onto Today.
  const insights = [
    { glyph: '↗', text: 'Referrals convert at <b>3×</b> the rate of cold apps.', detail: '4 / 5 referrals → screen vs. 3 / 10 cold' },
    { glyph: '◷', text: 'You haven\'t applied in <b>5 days</b>.', detail: 'Your usual pace is ~2 per day' },
    { glyph: '⚠', text: '<b>3 loops</b> have gone quiet for over a week.', detail: 'Linear, Notion, Figma — worth a nudge' }
  ];

  // Pipeline preview list
  const apps = [
    { co: 'Stripe',    role: 'Staff Eng, Payments',     status: 'interview', applied: '4d ago', domain: 'stripe.com',   stale: false },
    { co: 'Vercel',    role: 'Senior PM, Edge',         status: 'offer',     applied: '18d ago', domain: 'vercel.com',   stale: false },
    { co: 'Anthropic', role: 'Research Engineer',       status: 'screen',    applied: '6d ago', domain: 'anthropic.com', stale: false },
    { co: 'Linear',    role: 'Senior Frontend Eng',     status: 'applied',   applied: '7d ago', domain: 'linear.app',    stale: true  },
    { co: 'Notion',    role: 'Eng Manager, Editor',     status: 'applied',   applied: '9d ago', domain: 'notion.so',     stale: true  },
    { co: 'Supabase',  role: 'Developer Advocate',      status: 'applied',   applied: '3d ago', domain: 'supabase.com',  stale: false },
    { co: 'Figma',     role: 'Design Eng',              status: 'applied',   applied: '11d ago', domain: 'figma.com',    stale: true  }
  ];

  const STATUS_LABEL = { wishlist:'Wishlist', applied:'Applied', screen:'Screen', interview:'Interview', offer:'Offer' };

  const today = new Date();
  const dow  = today.toLocaleDateString('en-US', { weekday: 'long' });
  const dnum = today.toLocaleDateString('en-US', { day: 'numeric', month: 'long', year: 'numeric' });
</script>

<svelte:head><title>Today (redesign preview) — Pursuit</title></svelte:head>

<div class="frame">
  <!-- Compact sidebar so the mockup feels framed in-app -->
  <aside class="sidebar">
    <div class="brand"><span class="mark"></span><span class="name">Pursuit</span></div>
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
    <!-- Top bar with avatar swap -->
    <div class="topbar">
      <div class="crumb"><span class="here">Today</span></div>
      <div class="right">
        <div class="search">
          <svg width="14" height="14" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="7" cy="7" r="4.5"/><path d="M10.5 10.5L13 13"/></svg>
          <span>Search applications, people…</span>
          <span class="kbd">⌘K</span>
        </div>
        <button class="btn">Import</button>
        <button class="btn btn-primary">New application <span class="kbd">⌘N</span></button>
        <!-- Google avatar replacing the "BA" initials -->
        <img class="user-av" src="https://www.gravatar.com/avatar/?d=mp&s=64" alt={me.name} title={me.email} />
      </div>
    </div>

    <div class="body">
      <div class="body-inner">

        <!-- Greeting row -->
        <div class="hello">
          <div>
            <div class="date">{dow} <span class="sep">·</span> {dnum}</div>
            <h1>Good afternoon, {me.name}.</h1>
          </div>
        </div>

        <!-- BIG counts header -->
        <div class="counts">
          <div class="count-cell live">
            <div class="n">{counts.interviews}</div>
            <div class="lbl">Interviews <span class="hint" title="An interview loop = a series of rounds with one company">i</span></div>
          </div>
          <div class="count-cell warm">
            <div class="n">{counts.offers}</div>
            <div class="lbl">Open offers</div>
          </div>
          <div class="count-cell">
            <div class="n">{counts.applied}</div>
            <div class="lbl">Applied · waiting</div>
          </div>
          <div class="count-cell">
            <div class="n">{counts.wishlist}</div>
            <div class="lbl">Wishlist</div>
          </div>
        </div>

        <!-- What you can do today -->
        <div class="section-hd">
          <h2>What you can do today</h2>
          <span class="ai-tag">AI suggested</span>
        </div>
        <div class="action-grid">
          {#each actions as a}
            <div class={`action-card ${a.color}`}>
              <div class="action-top">
                <img class="action-logo" src={`https://www.google.com/s2/favicons?sz=128&domain=${a.logo}`} alt="" />
                <span class={`urgency u-${a.color}`}>{a.urgency}</span>
              </div>
              <h3>{a.title}</h3>
              <p>{a.sub}</p>
              <button class="action-cta">{a.cta} <span class="arrow">→</span></button>
            </div>
          {/each}
        </div>

        <!-- What we're noticing (moved from funnel) -->
        <div class="section-hd">
          <h2>What we're noticing</h2>
          <span class="ai-tag">AI · this week</span>
        </div>
        <div class="insight-list">
          {#each insights as ins}
            <div class="insight">
              <span class="ins-glyph">{ins.glyph}</span>
              <div class="ins-body">
                <div class="ins-line">{@html ins.text}</div>
                <div class="ins-detail">{ins.detail}</div>
              </div>
              <button class="ins-act">View →</button>
            </div>
          {/each}
        </div>

        <!-- Pipeline -->
        <div class="section-hd">
          <h2>Applications <span class="count">{apps.length}</span></h2>
          <div class="filters">
            <button class="chip active">Active <span class="cn">{apps.length}</span></button>
            <button class="chip">Interview <span class="cn">{counts.interviews}</span></button>
            <button class="chip">Offer <span class="cn">{counts.offers}</span></button>
            <button class="chip">Wishlist <span class="cn">{counts.wishlist}</span></button>
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
                <img class="logo" src={`https://www.google.com/s2/favicons?sz=128&domain=${a.domain}`} alt="" />
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

  /* Sidebar (compact mockup version) */
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
  .crumb .here { font-weight: 500; font-size: 13.5px; }
  .right { display: flex; align-items: center; gap: 8px; }
  .search { display: flex; align-items: center; gap: 6px; background: var(--card); border: 1px solid var(--rule); border-radius: 7px; padding: 5px 10px; font-size: 12.5px; color: var(--mute); min-width: 280px; }
  .kbd { font-family: var(--mono); font-size: 10.5px; color: var(--mute); background: var(--surface-2); padding: 1px 5px; border-radius: 3px; margin-left: auto; }
  .btn { background: var(--card); border: 1px solid var(--rule); border-radius: 7px; padding: 6px 11px; font-size: 12.5px; font-weight: 500; color: var(--ink); cursor: pointer; }
  .btn-primary { background: var(--accent); border-color: var(--accent-strong); color: white; }
  .user-av { width: 30px; height: 30px; border-radius: 50%; object-fit: cover; background: var(--rule-strong); cursor: pointer; margin-left: 8px; border: 1px solid var(--rule); }

  /* Body */
  .body { padding: 28px; }
  .body-inner { max-width: 1080px; margin: 0 auto; }

  /* Greeting */
  .hello { margin-bottom: 28px; }
  .hello .date { font-family: var(--mono); font-size: 11.5px; color: var(--mute); letter-spacing: 0.02em; text-transform: uppercase; margin-bottom: 4px; }
  .hello .date .sep { opacity: 0.5; margin: 0 4px; }
  .hello h1 { font-size: 26px; font-weight: 500; margin: 0; letter-spacing: -0.02em; }

  /* BIG counts */
  .counts { display: grid; grid-template-columns: repeat(4, 1fr); gap: 1px; background: var(--rule); border: 1px solid var(--rule); border-radius: 12px; overflow: hidden; margin-bottom: 36px; box-shadow: var(--sh-1); }
  .count-cell { background: var(--card); padding: 22px 24px 20px; }
  .count-cell .n { font-size: 44px; font-weight: 500; letter-spacing: -0.04em; line-height: 1; color: var(--ink); font-feature-settings: "tnum"; }
  .count-cell .lbl { font-size: 12px; color: var(--mute); text-transform: uppercase; letter-spacing: 0.06em; margin-top: 10px; display: flex; align-items: center; gap: 6px; }
  .count-cell.live .n { color: var(--accent-text); }
  .count-cell.warm .n { color: var(--warm-text); }
  .hint { width: 13px; height: 13px; display: inline-grid; place-items: center; border-radius: 50%; background: var(--surface-2); color: var(--mute); font-size: 9.5px; font-style: italic; font-weight: 600; cursor: help; }

  /* Section heading */
  .section-hd { display: flex; align-items: center; justify-content: space-between; margin: 30px 0 12px; }
  .section-hd h2 { font-size: 14px; font-weight: 600; margin: 0; letter-spacing: -0.01em; }
  .section-hd h2 .count { font-family: var(--mono); font-size: 11px; color: var(--mute); margin-left: 6px; font-weight: 400; }
  .ai-tag { font-family: var(--mono); font-size: 10px; background: var(--accent-tint); color: var(--accent-text); padding: 2px 7px; border-radius: 4px; letter-spacing: 0.04em; }

  /* Action grid */
  .action-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(260px, 1fr)); gap: 12px; }
  .action-card { background: var(--card); border: 1px solid var(--rule); border-radius: 12px; padding: 16px; display: flex; flex-direction: column; gap: 8px; transition: transform 120ms ease, border-color 120ms ease, box-shadow 120ms ease; cursor: pointer; }
  .action-card:hover { transform: translateY(-2px); border-color: var(--accent); box-shadow: var(--sh-pop); }
  .action-top { display: flex; align-items: center; justify-content: space-between; }
  .action-logo { width: 28px; height: 28px; border-radius: 6px; background: var(--surface-2); object-fit: contain; padding: 3px; }
  .urgency { font-family: var(--mono); font-size: 10px; padding: 2px 7px; border-radius: 4px; letter-spacing: 0.04em; }
  .urgency.u-accent { background: var(--accent-tint); color: var(--accent-text); }
  .urgency.u-warm { background: var(--warm-tint); color: var(--warm-text); }
  .urgency.u-mute { background: var(--surface-2); color: var(--mute); }
  .action-card h3 { font-size: 14.5px; font-weight: 500; margin: 4px 0 0; letter-spacing: -0.01em; }
  .action-card p { font-size: 12.5px; color: var(--mute); margin: 0; line-height: 1.45; }
  .action-cta { background: transparent; border: 0; color: var(--accent-text); font-size: 12.5px; font-weight: 500; text-align: left; padding: 4px 0 0; cursor: pointer; align-self: flex-start; }
  .action-cta .arrow { transition: transform 120ms ease; display: inline-block; }
  .action-card:hover .action-cta .arrow { transform: translateX(2px); }

  /* Insights */
  .insight-list { display: flex; flex-direction: column; gap: 1px; background: var(--rule); border: 1px solid var(--rule); border-radius: 12px; overflow: hidden; }
  .insight { background: var(--card); padding: 14px 18px; display: grid; grid-template-columns: 28px 1fr auto; gap: 12px; align-items: center; }
  .ins-glyph { width: 28px; height: 28px; border-radius: 7px; background: var(--accent-tint); color: var(--accent-text); display: grid; place-items: center; font-size: 14px; }
  .ins-line { font-size: 13.5px; color: var(--ink); }
  .ins-detail { font-size: 11.5px; color: var(--mute); margin-top: 2px; }
  .ins-act { background: transparent; border: 0; color: var(--accent-text); font-size: 12.5px; font-weight: 500; cursor: pointer; }

  /* Pipeline table */
  .filters { display: flex; gap: 4px; }
  .chip { background: transparent; border: 1px solid var(--rule); border-radius: 6px; padding: 4px 9px; font-size: 12px; color: var(--mute); cursor: pointer; }
  .chip.active { background: var(--ink); border-color: var(--ink); color: white; }
  .chip .cn { font-family: var(--mono); font-size: 10.5px; margin-left: 4px; opacity: 0.7; }
  .table { background: var(--card); border: 1px solid var(--rule); border-radius: 12px; overflow: hidden; }
  .tr { display: grid; grid-template-columns: 220px 1fr 130px 120px 20px; align-items: center; padding: 11px 18px; border-bottom: 1px solid var(--rule); font-size: 13px; cursor: pointer; }
  .tr:last-child { border-bottom: 0; }
  .tr:hover { background: var(--surface-2); }
  .tr.head { background: var(--surface-2); font-size: 11px; color: var(--mute); text-transform: uppercase; letter-spacing: 0.06em; cursor: default; padding: 9px 18px; font-weight: 500; }
  .tr.head:hover { background: var(--surface-2); }
  .tr.stale .applied { color: var(--danger-text); }
  .co { display: flex; align-items: center; gap: 10px; font-weight: 500; }
  .co .logo { width: 22px; height: 22px; border-radius: 5px; background: var(--surface-2); object-fit: contain; padding: 2px; }
  .role { color: var(--mute); }
  .pill { display: inline-flex; align-items: center; gap: 5px; padding: 3px 8px; border-radius: 99px; font-size: 11.5px; font-weight: 500; background: var(--surface-2); color: var(--ink-2); width: max-content; }
  .pill .pdot { width: 5px; height: 5px; border-radius: 50%; background: var(--mute-2); }
  .pill.interview { background: var(--accent-tint); color: var(--accent-text); }
  .pill.interview .pdot { background: var(--accent); }
  .pill.offer { background: var(--warm-tint); color: var(--warm-text); }
  .pill.offer .pdot { background: var(--warm); }
  .pill.screen { background: var(--positive-tint); color: var(--positive-text); }
  .pill.screen .pdot { background: var(--positive); }
  .applied { color: var(--mute); font-feature-settings: "tnum"; }
  .stale-tag { font-family: var(--mono); font-size: 10px; background: var(--danger-tint); color: var(--danger-text); padding: 1px 6px; border-radius: 3px; margin-left: 4px; }
  .arrow { color: var(--mute-2); }

  .footer-link { margin-top: 36px; font-size: 12.5px; }
  .footer-link a { color: var(--accent-text); text-decoration: none; }
</style>
