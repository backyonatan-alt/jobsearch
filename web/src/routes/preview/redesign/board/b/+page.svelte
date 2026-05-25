<script>
  // Variant B — "Cinematic kanban"
  // Larger cards, colored status accent strip on top, big circular logos,
  // stale = soft red glow + pulse, drag = tilt animation.
  const STATUSES = ['wishlist', 'applied', 'screen', 'interview', 'offer', 'rejected'];
  const STATUS_LABEL = { wishlist:'Wishlist', applied:'Applied', screen:'Screen', interview:'Interview', offer:'Offer', rejected:'Closed' };

  const apps = [
    { co:'Vercel',    domain:'vercel.com',    role:'Senior PM, Edge',           status:'offer',     source:'Referral', applied:'18d ago', stale:false, salary:'$280k base' },
    { co:'Stripe',    domain:'stripe.com',    role:'Staff Eng, Payments',       status:'interview', source:'Referral', applied:'4d ago',  stale:false, next:'Tomorrow · 2:00 PM' },
    { co:'Cursor',    domain:'cursor.com',    role:'Engineering Lead',          status:'interview', source:'Referral', applied:'2d ago',  stale:false, next:'Thu · panel' },
    { co:'Anthropic', domain:'anthropic.com', role:'Research Engineer',         status:'screen',    source:'LinkedIn', applied:'6d ago',  stale:false },
    { co:'Modal',     domain:'modal.com',     role:'Infra Eng',                 status:'screen',    source:'Referral', applied:'5d ago',  stale:false },
    { co:'Linear',    domain:'linear.app',    role:'Senior Frontend Eng',       status:'applied',   source:'Referral', applied:'7d ago',  stale:true  },
    { co:'Notion',    domain:'notion.so',     role:'Eng Manager, Editor',       status:'applied',   source:'LinkedIn', applied:'9d ago',  stale:true  },
    { co:'Supabase',  domain:'supabase.com',  role:'Developer Advocate',        status:'applied',   source:'Cold',     applied:'3d ago',  stale:false },
    { co:'Figma',     domain:'figma.com',     role:'Design Eng',                status:'applied',   source:'LinkedIn', applied:'11d ago', stale:true  },
    { co:'Replicate', domain:'replicate.com', role:'Platform Eng',              status:'applied',   source:'Cold',     applied:'2d ago',  stale:false },
    { co:'Perplexity',domain:'perplexity.ai', role:'Search Engineer',           status:'wishlist',  source:'—',        applied:'—',       stale:false },
    { co:'Mistral',   domain:'mistral.ai',    role:'Research Engineer',         status:'wishlist',  source:'—',        applied:'—',       stale:false },
    { co:'Pinecone',  domain:'pinecone.io',   role:'PM, Platform',              status:'rejected',  source:'Cold',     applied:'30d ago', stale:false }
  ];

  const me = { name: 'Yonatan', email: 'back.yonatan@gmail.com' };
  const today = new Date();
  const dateLong = today.toLocaleDateString('en-US', { weekday: 'long', day: 'numeric', month: 'long' });
  const byStatus = Object.fromEntries(STATUSES.map(s => [s, apps.filter(a => a.status === s)]));
  const MUTED = new Set(['rejected']);
</script>

<svelte:head><title>Board · variant B — Pursuit</title></svelte:head>

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
        <div class="who">{me.name}<small>{me.email}</small></div>
      </div>
    </div>
  </aside>

  <section class="main">
    <div class="topbar">
      <div class="crumb"><span class="here">Board</span></div>
      <div class="right">
        <div class="search">
          <svg width="14" height="14" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="7" cy="7" r="4.5"/><path d="M10.5 10.5L13 13"/></svg>
          <span>Search applications, people</span>
        </div>
        <button class="btn btn-primary">New application</button>
        <img class="user-av" src="https://www.gravatar.com/avatar/?d=mp&s=64" alt="" />
      </div>
    </div>

    <div class="body">
      <div class="board-page">
        <div class="board-hd">
          <div>
            <div class="date">{dateLong}</div>
            <h1>Board.</h1>
            <p class="sub"><b>{apps.filter(a => !MUTED.has(a.status)).length}</b> in flight · drag a card across columns to move its status.</p>
          </div>
        </div>

        <div class="board-cols">
          {#each STATUSES as s (s)}
            <section class={`board-col col-${s}`} class:muted={MUTED.has(s)}>
              <header class="board-col-hd">
                <span class={`col-dot d-${s}`}></span>
                <span class="col-title">{STATUS_LABEL[s]}</span>
                <span class="count">{byStatus[s].length}</span>
                <button class="add" title="Add to {STATUS_LABEL[s]}">+</button>
              </header>
              <div class="cards">
                {#each byStatus[s] as a}
                  <button type="button" class={`bcard tone-${s}`} class:stale={a.stale} draggable="true">
                    <span class="card-accent"></span>
                    <div class="top">
                      <img class="logo" src={`https://logo.clearbit.com/${a.domain}`} alt="" />
                      <div class="top-text">
                        <span class="co">{a.co}</span>
                        <span class="role">{a.role}</span>
                      </div>
                      {#if a.stale}<span class="stale-pulse" title="No movement for over a week"></span>{/if}
                    </div>
                    {#if a.next}
                      <div class="next">
                        <svg width="11" height="11" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.6"><circle cx="6" cy="6" r="4.5"/><path d="M6 3.5V6l1.5 1"/></svg>
                        {a.next}
                      </div>
                    {/if}
                    {#if a.salary}
                      <div class="salary">{a.salary}</div>
                    {/if}
                    <div class="foot">
                      <span class="source">{a.source}</span>
                      <span class="applied" class:stale-text={a.stale}>
                        {#if a.stale}<span class="needs-nudge">needs nudge</span>{/if}
                        {a.applied}
                      </span>
                    </div>
                  </button>
                {/each}
                {#if byStatus[s].length === 0}
                  <div class="board-empty">drop here</div>
                {/if}
              </div>
            </section>
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
  .brand { display: flex; align-items: center; gap: 10px; padding: 4px 8px 18px; color: var(--accent-text); }
  .brand-name { font-size: 18px; font-weight: 600; letter-spacing: -0.02em; color: var(--ink); }
  .brand-mark { color: var(--accent); }
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
  .crumb .here { font-weight: 600; font-size: 14px; }
  .right { display: flex; align-items: center; gap: 8px; }
  .search { display: flex; align-items: center; gap: 6px; background: var(--card); border: 1px solid var(--rule); border-radius: 7px; padding: 5px 10px; font-size: 13px; color: var(--mute); min-width: 280px; }
  .btn { background: var(--card); border: 1px solid var(--rule); border-radius: 7px; padding: 6px 11px; font-size: 13px; font-weight: 500; color: var(--ink); cursor: pointer; }
  .btn-primary { background: var(--accent); border-color: var(--accent-strong); color: white; }
  .user-av { width: 30px; height: 30px; border-radius: 50%; object-fit: cover; cursor: pointer; margin-left: 8px; border: 1px solid var(--rule); }

  .body { padding: 28px; }
  .board-page { max-width: 1500px; margin: 0 auto; }

  /* Header */
  .board-hd { margin-bottom: 28px; }
  .board-hd .date { font-size: 13.5px; color: var(--mute); margin-bottom: 6px; }
  .board-hd h1 { font-size: 30px; font-weight: 600; letter-spacing: -0.025em; margin: 0; }
  .board-hd .sub { font-size: 13.5px; color: var(--mute); margin: 6px 0 0; }

  /* Columns */
  .board-cols {
    display: grid;
    grid-template-columns: repeat(6, minmax(260px, 1fr));
    gap: 16px;
    overflow-x: auto;
    padding-bottom: 8px;
  }
  .board-col {
    background: var(--surface-2);
    border-radius: 16px;
    padding: 14px;
    min-height: 60vh;
  }
  .board-col.muted { opacity: 0.55; }
  .board-col-hd { display: flex; align-items: center; gap: 8px; margin-bottom: 14px; padding: 0 4px; }
  .col-dot { width: 9px; height: 9px; border-radius: 50%; background: var(--mute-2); }
  .d-wishlist { background: var(--mute-2); }
  .d-applied  { background: var(--ink-2); }
  .d-screen   { background: var(--positive); }
  .d-interview{ background: var(--accent); }
  .d-offer    { background: var(--warm); }
  .d-rejected { background: var(--mute-2); }
  .col-title { font-size: 13.5px; font-weight: 600; letter-spacing: -0.01em; }
  .board-col-hd .count {
    font-size: 12px; color: var(--mute); background: var(--card);
    padding: 2px 8px; border-radius: 99px; font-weight: 500;
  }
  .board-col-hd .add {
    margin-left: auto; width: 24px; height: 24px; border-radius: 7px;
    background: var(--card); border: 1px solid var(--rule); cursor: pointer;
    display: grid; place-items: center; color: var(--mute); font-size: 16px; line-height: 1;
  }
  .board-col-hd .add:hover { color: var(--ink); }

  /* Cards */
  .cards { display: flex; flex-direction: column; gap: 12px; }
  .bcard {
    position: relative;
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 13px;
    padding: 14px 14px 12px;
    text-align: left;
    cursor: grab;
    box-shadow: var(--sh-1);
    transition: transform 180ms cubic-bezier(0.4, 0, 0.2, 1), box-shadow 180ms ease, border-color 180ms ease;
    display: flex; flex-direction: column; gap: 8px;
    overflow: hidden;
  }
  /* Top accent strip — colored by status */
  .card-accent { position: absolute; top: 0; left: 0; right: 0; height: 3px; background: var(--rule-strong); }
  .tone-wishlist .card-accent { background: var(--mute-2); }
  .tone-applied  .card-accent { background: var(--ink-2); }
  .tone-screen   .card-accent { background: var(--positive); }
  .tone-interview.card-accent,
  .tone-interview .card-accent { background: var(--accent); }
  .tone-offer    .card-accent { background: var(--warm); }
  .tone-rejected .card-accent { background: var(--mute-2); }

  .bcard:hover { transform: translateY(-3px); box-shadow: var(--sh-pop); border-color: var(--rule-strong); }
  .bcard:active { cursor: grabbing; transform: scale(0.98) rotate(-1deg); transition-duration: 60ms; box-shadow: var(--sh-pop); }

  /* Stale: red glow + pulsing dot */
  .bcard.stale {
    border-color: var(--danger);
    box-shadow: 0 0 0 3px var(--danger-tint), var(--sh-1);
  }
  .bcard.stale .card-accent { background: var(--danger); }
  @keyframes stale-pulse {
    0%, 100% { box-shadow: 0 0 0 0 var(--danger); opacity: 1; }
    50%      { box-shadow: 0 0 0 5px transparent; opacity: 0.55; }
  }
  .stale-pulse {
    width: 8px; height: 8px; border-radius: 50%; background: var(--danger);
    margin-left: auto; animation: stale-pulse 1.6s ease-in-out infinite;
  }

  .top { display: flex; align-items: center; gap: 10px; }
  .top .logo {
    width: 32px; height: 32px; border-radius: 50%;
    background: var(--surface-2); object-fit: contain; padding: 4px;
    border: 1px solid var(--rule); flex-shrink: 0;
  }
  .top-text { display: flex; flex-direction: column; gap: 1px; min-width: 0; flex: 1; }
  .top-text .co { font-size: 14px; font-weight: 600; letter-spacing: -0.01em; }
  .top-text .role { font-size: 12px; color: var(--mute); line-height: 1.3; }

  .next {
    font-size: 12px; color: var(--accent-text);
    background: var(--accent-tint); padding: 5px 9px; border-radius: 8px;
    display: inline-flex; align-items: center; gap: 6px; align-self: flex-start;
    font-weight: 500;
  }

  .salary {
    font-size: 12px; color: var(--warm-text);
    background: var(--warm-tint); padding: 4px 9px; border-radius: 8px;
    align-self: flex-start; font-weight: 500;
  }

  .foot { display: flex; justify-content: space-between; align-items: center; font-size: 11.5px; color: var(--mute); padding-top: 8px; border-top: 1px dashed var(--rule); margin-top: 2px; }
  .foot .stale-text { color: var(--danger-text); font-weight: 500; display: inline-flex; align-items: center; gap: 6px; }
  .needs-nudge {
    font-size: 10px; background: var(--danger-tint); color: var(--danger-text);
    padding: 1px 6px; border-radius: 99px; font-weight: 600;
  }

  .board-empty {
    border: 1.5px dashed var(--rule-strong);
    border-radius: 12px;
    text-align: center;
    color: var(--mute-2);
    font-size: 12px;
    padding: 22px;
  }

  .footer-link { margin-top: 30px; font-size: 13px; }
  .footer-link a { color: var(--accent-text); text-decoration: none; }
</style>
