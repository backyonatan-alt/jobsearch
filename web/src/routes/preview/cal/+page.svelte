<script>
  import { PREVIEW_APPS, fmtDate } from '$lib/preview-data.js';

  const STATUS = {
    wishlist:  { bg: '#f3f4f6', fg: '#475569' },
    applied:   { bg: '#dbeafe', fg: '#1e40af' },
    screen:    { bg: '#fef3c7', fg: '#92400e' },
    interview: { bg: '#ede9fe', fg: '#5b21b6' },
    offer:     { bg: '#dcfce7', fg: '#166534' },
    rejected:  { bg: '#fee2e2', fg: '#991b1b' },
    withdrawn: { bg: '#f3f4f6', fg: '#475569' }
  };

  const TONES = ['#6366f1', '#ec4899', '#06b6d4', '#10b981', '#f59e0b', '#8b5cf6'];
  const tone = (s) => TONES[s.charCodeAt(0) % TONES.length];
</script>

<svelte:head>
  <title>Cal.com sketch — Pursuit</title>
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous" />
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap" rel="stylesheet" />
</svelte:head>

<div class="shell">
  <aside class="side">
    <div class="brand">
      <div class="bm">P</div>
      <div>
        <p class="bn">Pursuit</p>
        <p class="bs">back.yonatan@gmail.com</p>
      </div>
    </div>

    <button class="cta">+ New application</button>

    <nav>
      <p class="g">Workspace</p>
      <a class="nl cur" href="#">
        <svg viewBox="0 0 24 24" width="16" height="16"><path d="M3 3h7v7H3zM14 3h7v7h-7zM3 14h7v7H3zM14 14h7v7h-7z" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round"/></svg>
        Applications
        <span class="ct">6</span>
      </a>
      <a class="nl" href="#">
        <svg viewBox="0 0 24 24" width="16" height="16"><path d="M4 5h4v14H4zM10 5h4v9h-4zM16 5h4v6h-4z" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round"/></svg>
        Board
      </a>
      <a class="nl" href="#">
        <svg viewBox="0 0 24 24" width="16" height="16"><path d="M3 4h18l-7 9v7l-4-2v-5z" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round"/></svg>
        Funnel
      </a>
    </nav>

    <div class="flex"></div>
    <a class="leave" href="/preview">‹ leave preview</a>
  </aside>

  <main>
    <header class="top">
      <div class="search">
        <svg viewBox="0 0 24 24" width="16" height="16"><circle cx="11" cy="11" r="6" fill="none" stroke="currentColor" stroke-width="1.7"/><path d="M16 16l4 4" stroke="currentColor" stroke-width="1.7" stroke-linecap="round"/></svg>
        <input placeholder="Search…" />
      </div>
      <button class="ic" aria-label="Notifications">
        <svg viewBox="0 0 24 24" width="18" height="18"><path d="M6 16V11a6 6 0 1 1 12 0v5l1 2H5l1-2zM9 19a3 3 0 0 0 6 0" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linejoin="round" stroke-linecap="round"/></svg>
      </button>
      <div class="av">B</div>
    </header>

    <div class="page">
      <div class="head">
        <div>
          <h1>Applications</h1>
          <p class="sub">A clean overview of every conversation you're in.</p>
        </div>
        <button class="primary">+ New application</button>
      </div>

      <div class="kpis">
        <div class="kpi">
          <p class="kl">Total</p>
          <p class="kv">24</p>
          <p class="kd"><span class="up">↑ 5</span> this week</p>
        </div>
        <div class="kpi">
          <p class="kl">Active loops</p>
          <p class="kv">5</p>
          <p class="kd"><span class="up">↑ 2</span> from last week</p>
        </div>
        <div class="kpi">
          <p class="kl">Open offers</p>
          <p class="kv">1</p>
          <p class="kd">Linear · reply by Fri</p>
        </div>
        <div class="kpi">
          <p class="kl">Reply rate</p>
          <p class="kv">62%</p>
          <p class="kd"><span class="up">↑ 8%</span> vs last month</p>
        </div>
      </div>

      <div class="card">
        <div class="card-head">
          <div class="tabs">
            <button class="tab cur">All</button>
            <button class="tab">Active</button>
            <button class="tab">Offers</button>
            <button class="tab">Closed</button>
          </div>
          <button class="ghost">Filter</button>
        </div>

        <ul class="apps">
          {#each PREVIEW_APPS as a}
            <li>
              <div class="left">
                <span class="mono" style="background: {tone(a.company)}">{a.company[0]}</span>
                <div>
                  <p class="co">{a.company}</p>
                  <p class="ro">{a.role} · {a.location}</p>
                </div>
              </div>
              <div class="right">
                <span class="pill" style="background: {STATUS[a.status].bg}; color: {STATUS[a.status].fg}">{a.status}</span>
                <span class="when">{fmtDate(a.applied_at)}</span>
                <button class="dots">⋯</button>
              </div>
            </li>
          {/each}
        </ul>
      </div>
    </div>
  </main>
</div>

<style>
  :global(html, body) {
    background: #fafafa; color: #0a0a0a;
    font-family: 'Inter', system-ui, sans-serif;
    font-size: 14px; line-height: 1.5;
    -webkit-font-smoothing: antialiased;
  }
  .shell { display: grid; grid-template-columns: 260px 1fr; min-height: 100vh; }

  .side { background: #fff; border-right: 1px solid #e5e7eb; padding: 1.25rem 1rem; display: flex; flex-direction: column; gap: 1.25rem; }
  .brand { display: flex; align-items: center; gap: .65rem; padding: .25rem .35rem; }
  .bm { width: 36px; height: 36px; border-radius: 10px; background: linear-gradient(135deg, #4f46e5, #7c3aed); color: #fff; font-weight: 700; display: grid; place-items: center; box-shadow: 0 2px 8px rgba(79, 70, 229, .25); }
  .bn { margin: 0; font-weight: 600; font-size: 14px; }
  .bs { margin: 0; font-size: 12px; color: #6b7280; }

  .cta { padding: .65rem 1rem; background: #4f46e5; color: #fff; border: 0; border-radius: 10px; font: inherit; font-size: 13px; font-weight: 600; cursor: pointer; box-shadow: 0 1px 2px rgba(79, 70, 229, .25); }
  .cta:hover { background: #4338ca; }

  nav { display: flex; flex-direction: column; gap: 2px; }
  .g { margin: 0 0 .35rem .75rem; font-size: 11px; font-weight: 600; color: #9ca3af; text-transform: uppercase; letter-spacing: .04em; }
  .nl { display: flex; align-items: center; gap: .65rem; padding: .55rem .75rem; border-radius: 8px; color: #374151; text-decoration: none; font-size: 14px; font-weight: 500; }
  .nl:hover { background: #f3f4f6; color: #111827; }
  .nl.cur { background: #eef2ff; color: #4f46e5; }
  .nl.cur svg { color: #4f46e5; }
  .nl svg { color: #9ca3af; }
  .nl .ct { margin-left: auto; background: #fff; border: 1px solid #e5e7eb; padding: 1px 7px; border-radius: 999px; font-size: 11px; color: #6b7280; font-weight: 500; }

  .flex { flex: 1; }
  .leave { color: #9ca3af; font-size: 12px; text-decoration: none; padding: .5rem .75rem; }
  .leave:hover { color: #4f46e5; }

  main { display: flex; flex-direction: column; min-width: 0; }
  .top { display: flex; align-items: center; gap: .65rem; padding: .85rem 1.5rem; border-bottom: 1px solid #e5e7eb; background: rgba(250, 250, 250, .85); backdrop-filter: blur(8px); position: sticky; top: 0; z-index: 10; }
  .search { flex: 1; max-width: 480px; display: flex; align-items: center; gap: .55rem; padding: .55rem .85rem; background: #fff; border: 1px solid #e5e7eb; border-radius: 10px; color: #9ca3af; }
  .search input { flex: 1; border: 0; background: transparent; outline: 0; font: inherit; color: #0a0a0a; }
  .ic { width: 36px; height: 36px; border-radius: 10px; background: #fff; border: 1px solid #e5e7eb; color: #4b5563; cursor: pointer; display: grid; place-items: center; }
  .ic:hover { background: #f9fafb; }
  .av { width: 32px; height: 32px; border-radius: 999px; background: linear-gradient(135deg, #4f46e5, #ec4899); color: #fff; display: grid; place-items: center; font-weight: 600; font-size: 12px; }

  .page { padding: 2rem 2.5rem 4rem; max-width: 1200px; width: 100%; }
  .head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
  h1 { font-size: 28px; font-weight: 700; letter-spacing: -.02em; margin: 0 0 .35rem; }
  .sub { color: #6b7280; margin: 0; }
  .primary { padding: .65rem 1.25rem; background: #111827; color: #fff; border: 0; border-radius: 10px; font: inherit; font-weight: 500; font-size: 14px; cursor: pointer; box-shadow: 0 1px 3px rgba(0,0,0,.08); }
  .primary:hover { background: #1f2937; }

  .kpis { display: grid; grid-template-columns: repeat(4, 1fr); gap: 1rem; margin-bottom: 2rem; }
  @media (max-width: 1000px) { .kpis { grid-template-columns: repeat(2, 1fr); } }
  .kpi { background: #fff; border: 1px solid #e5e7eb; border-radius: 14px; padding: 1.25rem 1.35rem; box-shadow: 0 1px 2px rgba(0,0,0,.03); }
  .kl { margin: 0; font-size: 12px; color: #6b7280; font-weight: 500; }
  .kv { margin: .5rem 0 .35rem; font-size: 32px; font-weight: 700; letter-spacing: -.02em; color: #111827; }
  .kd { margin: 0; font-size: 12px; color: #6b7280; }
  .up { color: #10b981; font-weight: 600; }

  .card { background: #fff; border: 1px solid #e5e7eb; border-radius: 14px; overflow: hidden; box-shadow: 0 1px 2px rgba(0,0,0,.03); }
  .card-head { display: flex; align-items: center; justify-content: space-between; padding: 1rem 1.5rem; border-bottom: 1px solid #e5e7eb; }
  .tabs { display: flex; gap: .25rem; }
  .tab { background: transparent; border: 0; padding: .35rem .85rem; border-radius: 8px; color: #6b7280; font: inherit; font-size: 13px; font-weight: 500; cursor: pointer; }
  .tab:hover { background: #f3f4f6; color: #111827; }
  .tab.cur { background: #111827; color: #fff; }
  .ghost { background: transparent; border: 1px solid #e5e7eb; padding: .35rem .85rem; border-radius: 8px; font: inherit; font-size: 13px; color: #4b5563; cursor: pointer; }
  .ghost:hover { background: #f9fafb; }

  .apps { list-style: none; padding: 0; margin: 0; }
  .apps li { display: flex; align-items: center; justify-content: space-between; padding: 1rem 1.5rem; border-bottom: 1px solid #f3f4f6; transition: background .12s ease; }
  .apps li:hover { background: #fafbfd; }
  .apps li:last-child { border-bottom: 0; }
  .left { display: flex; align-items: center; gap: 1rem; }
  .mono { width: 40px; height: 40px; border-radius: 10px; color: #fff; display: grid; place-items: center; font-weight: 700; font-size: 16px; box-shadow: 0 1px 2px rgba(0,0,0,.06); }
  .co { margin: 0; font-weight: 600; color: #111827; }
  .ro { margin: 1px 0 0; font-size: 13px; color: #6b7280; }
  .right { display: flex; align-items: center; gap: 1rem; }
  .pill { padding: 4px 10px; border-radius: 999px; font-size: 12px; font-weight: 500; text-transform: capitalize; }
  .when { color: #9ca3af; font-size: 13px; min-width: 4rem; text-align: right; }
  .dots { background: transparent; border: 0; color: #9ca3af; cursor: pointer; width: 28px; height: 28px; border-radius: 6px; font-size: 16px; }
  .dots:hover { background: #f3f4f6; color: #111827; }
</style>
