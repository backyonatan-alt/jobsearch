<script>
  // Variant A — "Refined classic kanban"
  // Tight cards with generous gaps, real logos, subtle drag lift,
  // stale = red dot + red applied date. Modern sans throughout.
  const STATUSES = ['wishlist', 'applied', 'screen', 'interview', 'offer', 'rejected'];
  const STATUS_LABEL = { wishlist:'Wishlist', applied:'Applied', screen:'Screen', interview:'Interview', offer:'Offer', rejected:'Closed' };

  const apps = [
    { co:'Vercel',    domain:'vercel.com',    role:'Senior PM, Edge',           status:'offer',     source:'Referral', applied:'18d ago', stale:false, cv:'v3' },
    { co:'Stripe',    domain:'stripe.com',    role:'Staff Eng, Payments',       status:'interview', source:'Referral', applied:'4d ago',  stale:false, cv:'v3', next:'Tomorrow · 2:00 PM' },
    { co:'Anthropic', domain:'anthropic.com', role:'Research Engineer',         status:'screen',    source:'LinkedIn', applied:'6d ago',  stale:false, cv:'v2' },
    { co:'Linear',    domain:'linear.app',    role:'Senior Frontend Eng',       status:'applied',   source:'Referral', applied:'7d ago',  stale:true,  cv:'v3' },
    { co:'Notion',    domain:'notion.so',     role:'Eng Manager, Editor',       status:'applied',   source:'LinkedIn', applied:'9d ago',  stale:true,  cv:'v2' },
    { co:'Supabase',  domain:'supabase.com',  role:'Developer Advocate',        status:'applied',   source:'Cold',     applied:'3d ago',  stale:false, cv:'v1' },
    { co:'Figma',     domain:'figma.com',     role:'Design Eng',                status:'applied',   source:'LinkedIn', applied:'11d ago', stale:true,  cv:'v3' },
    { co:'Replicate', domain:'replicate.com', role:'Platform Eng',              status:'applied',   source:'Cold',     applied:'2d ago',  stale:false, cv:'v2' },
    { co:'Modal',     domain:'modal.com',     role:'Infra Eng',                 status:'screen',    source:'Referral', applied:'5d ago',  stale:false, cv:'v3' },
    { co:'Cursor',    domain:'cursor.com',    role:'Engineering Lead',          status:'interview', source:'Referral', applied:'2d ago',  stale:false, cv:'v3', next:'Thu · panel' },
    { co:'Perplexity',domain:'perplexity.ai', role:'Search Engineer',           status:'wishlist',  source:'—',        applied:'—',       stale:false, cv:'—' },
    { co:'Mistral',   domain:'mistral.ai',    role:'Research Engineer',         status:'wishlist',  source:'—',        applied:'—',       stale:false, cv:'—' },
    { co:'Pinecone',  domain:'pinecone.io',   role:'PM, Platform',              status:'rejected',  source:'Cold',     applied:'30d ago', stale:false, cv:'v1' }
  ];

  const me = { name: 'Yonatan', email: 'back.yonatan@gmail.com' };

  const today = new Date();
  const dateLong = today.toLocaleDateString('en-US', { weekday: 'long', day: 'numeric', month: 'long' });

  const byStatus = Object.fromEntries(STATUSES.map(s => [s, apps.filter(a => a.status === s)]));
  const MUTED = new Set(['rejected']);
</script>

<svelte:head><title>Board · variant A — Pursuit</title></svelte:head>

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
          <!-- List/Board segmented toggle removed — the sidebar already
               navigates between Today (list) and Board, and the old
               segment caused a full page reload + scroll reset. -->
        </div>

        <div class="board-cols">
          {#each STATUSES as s (s)}
            <section class="board-col" class:muted={MUTED.has(s)}>
              <header class="board-col-hd">
                <span class={`pill ${s}`}><span class="pdot"></span>{STATUS_LABEL[s]}</span>
                <span class="count">{byStatus[s].length}</span>
                <button class="add" title="Add to {STATUS_LABEL[s]}">+</button>
              </header>
              <div class="cards">
                {#each byStatus[s] as a}
                  <button type="button" class="bcard" class:stale={a.stale} draggable="true">
                    <div class="top">
                      <img class="logo" src={`https://www.google.com/s2/favicons?sz=128&domain=${a.domain}`} alt="" />
                      <span class="co">{a.co}</span>
                      {#if a.stale}<span class="stale-dot" title="No movement for over a week"></span>{/if}
                    </div>
                    <p class="role">{a.role}</p>
                    {#if a.next}
                      <div class="next"><span class="next-dot"></span>{a.next}</div>
                    {/if}
                    <div class="foot">
                      <span class="source">{a.source}</span>
                      <span class="applied" class:stale-text={a.stale}>{a.applied}</span>
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
  .board-page { max-width: 1400px; margin: 0 auto; }

  /* Header */
  .board-hd { margin-bottom: 24px; }
  .board-hd .date { font-size: 13.5px; color: var(--mute); margin-bottom: 6px; }
  .board-hd h1 { font-size: 28px; font-weight: 600; letter-spacing: -0.025em; margin: 0; }
  .board-hd .sub { font-size: 13.5px; color: var(--mute); margin: 6px 0 0; }

  /* Columns */
  .board-cols {
    display: grid;
    grid-template-columns: repeat(6, minmax(240px, 1fr));
    gap: 14px;
    overflow-x: auto;
    padding-bottom: 8px;
  }
  .board-col {
    background: var(--surface-2);
    border-radius: 14px;
    padding: 12px;
    min-height: 60vh;
  }
  .board-col.muted { opacity: 0.65; }
  .board-col-hd { display: flex; align-items: center; gap: 8px; margin-bottom: 12px; padding: 0 4px; }
  .board-col-hd .count { font-size: 12.5px; color: var(--mute); font-weight: 500; }
  .board-col-hd .add {
    margin-left: auto; width: 22px; height: 22px; border-radius: 6px;
    background: transparent; border: 1px solid var(--rule); cursor: pointer;
    display: grid; place-items: center; color: var(--mute); font-size: 16px; line-height: 1;
    transition: background 120ms, color 120ms;
  }
  .board-col-hd .add:hover { background: var(--card); color: var(--ink); }

  /* Cards — gap between them, fixed */
  .cards { display: flex; flex-direction: column; gap: 10px; }
  .bcard {
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 11px;
    padding: 12px 13px;
    text-align: left;
    cursor: grab;
    box-shadow: var(--sh-1);
    transition: transform 140ms ease, box-shadow 140ms ease, border-color 140ms ease;
    display: flex; flex-direction: column; gap: 6px;
  }
  .bcard:hover { transform: translateY(-2px); box-shadow: var(--sh-pop); border-color: var(--rule-strong); }
  .bcard:active { cursor: grabbing; transform: scale(0.99) rotate(-0.5deg); }
  .bcard.stale {
    border-color: var(--danger);
    box-shadow: 0 0 0 1px var(--danger-tint), var(--sh-1);
  }
  .bcard.stale:hover {
    box-shadow: 0 0 0 1px var(--danger-tint), var(--sh-pop);
  }

  .top { display: flex; align-items: center; gap: 8px; }
  .top .logo { width: 22px; height: 22px; border-radius: 5px; background: var(--surface-2); object-fit: contain; padding: 2px; flex-shrink: 0; }
  .top .co { font-size: 13.5px; font-weight: 600; letter-spacing: -0.01em; }
  .top .stale-dot { width: 7px; height: 7px; border-radius: 50%; background: var(--danger); margin-left: auto; box-shadow: 0 0 0 3px var(--danger-tint); }

  .role { font-size: 12.5px; color: var(--mute); margin: 0; line-height: 1.4; }

  .next {
    font-size: 11.5px; color: var(--accent-text);
    background: var(--accent-tint); padding: 4px 8px; border-radius: 6px;
    display: inline-flex; align-items: center; gap: 6px; align-self: flex-start;
    font-weight: 500;
  }
  .next-dot { width: 5px; height: 5px; border-radius: 50%; background: var(--accent); }

  .foot { display: flex; justify-content: space-between; font-size: 11.5px; color: var(--mute); padding-top: 4px; border-top: 1px dashed var(--rule); margin-top: 2px; }
  .foot .stale-text { color: var(--danger-text); font-weight: 500; }

  .board-empty {
    border: 1.5px dashed var(--rule-strong);
    border-radius: 10px;
    text-align: center;
    color: var(--mute-2);
    font-size: 12px;
    padding: 18px;
  }

  /* Pills */
  .pill { display: inline-flex; align-items: center; gap: 5px; padding: 3px 9px; border-radius: 99px; font-size: 12px; font-weight: 500; background: var(--card); color: var(--ink-2); }
  .pill .pdot { width: 5px; height: 5px; border-radius: 50%; background: var(--mute-2); }
  .pill.wishlist { background: var(--surface-2); color: var(--mute); }
  .pill.applied { background: var(--card); color: var(--ink-2); }
  .pill.applied .pdot { background: var(--mute-2); }
  .pill.screen { background: var(--positive-tint); color: var(--positive-text); }
  .pill.screen .pdot { background: var(--positive); }
  .pill.interview { background: var(--accent-tint); color: var(--accent-text); }
  .pill.interview .pdot { background: var(--accent); }
  .pill.offer { background: var(--warm-tint); color: var(--warm-text); }
  .pill.offer .pdot { background: var(--warm); }
  .pill.rejected { background: var(--surface-2); color: var(--mute); }

  .footer-link { margin-top: 30px; font-size: 13px; }
  .footer-link a { color: var(--accent-text); text-decoration: none; }
</style>
