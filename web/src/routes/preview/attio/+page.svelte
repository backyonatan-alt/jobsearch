<script>
  import { PREVIEW_APPS, fmtDate } from '$lib/preview-data.js';

  const STATUS = {
    wishlist:  { bg: '#f4f4f5', fg: '#52525b' },
    applied:   { bg: '#eef0ff', fg: '#3939ad' },
    screen:    { bg: '#fef3c7', fg: '#92400e' },
    interview: { bg: '#f3e8ff', fg: '#6b21a8' },
    offer:     { bg: '#d1fae5', fg: '#065f46' },
    rejected:  { bg: '#fee2e2', fg: '#991b1b' },
    withdrawn: { bg: '#f4f4f5', fg: '#52525b' }
  };

  const TONES = ['#fb923c', '#a855f7', '#3b82f6', '#10b981', '#ec4899', '#f59e0b'];
  const tone = (s) => TONES[s.charCodeAt(0) % TONES.length];
</script>

<svelte:head>
  <title>Attio sketch — Pursuit</title>
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous" />
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet" />
</svelte:head>

<div class="shell">
  <aside class="side">
    <div class="ws">
      <div class="ws-mark">P</div>
      <div class="ws-name">Pursuit</div>
      <button class="ws-x" aria-label="Switch">▾</button>
    </div>

    <button class="new">+ New application <span class="k">N</span></button>

    <div class="group">
      <p class="g-title">Workspace</p>
      <a class="nl current" href="#"><svg viewBox="0 0 24 24" width="14" height="14"><path d="M3 3h7v7H3zM14 3h7v7h-7zM3 14h7v7H3zM14 14h7v7h-7z" fill="none" stroke="currentColor" stroke-width="1.5"/></svg>Applications<span class="ct">6</span></a>
      <a class="nl" href="#"><svg viewBox="0 0 24 24" width="14" height="14"><path d="M4 5h4v14H4zM10 5h4v9h-4zM16 5h4v6h-4z" fill="none" stroke="currentColor" stroke-width="1.5"/></svg>Board</a>
      <a class="nl" href="#"><svg viewBox="0 0 24 24" width="14" height="14"><path d="M3 4h18l-7 9v7l-4-2v-5z" fill="none" stroke="currentColor" stroke-width="1.5"/></svg>Funnel</a>
    </div>

    <div class="group">
      <p class="g-title">Companies</p>
      <a class="nl" href="#"><span class="cdot" style="background: {tone('A')}">A</span>Anthropic</a>
      <a class="nl" href="#"><span class="cdot" style="background: {tone('S')}">S</span>Stripe</a>
      <a class="nl" href="#"><span class="cdot" style="background: {tone('V')}">V</span>Vercel</a>
      <a class="nl" href="#"><span class="cdot" style="background: {tone('L')}">L</span>Linear</a>
    </div>

    <div class="flex"></div>
    <a class="leave" href="/preview">‹ leave preview</a>
  </aside>

  <main>
    <header class="top">
      <div class="search">
        <svg viewBox="0 0 24 24" width="14" height="14"><circle cx="11" cy="11" r="6" fill="none" stroke="currentColor" stroke-width="1.5"/><path d="M16 16l4 4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/></svg>
        <input placeholder="Search" />
        <span class="kbd">⌘K</span>
      </div>
      <button class="ic" aria-label="Notifications"><svg viewBox="0 0 24 24" width="15" height="15"><path d="M6 16V11a6 6 0 1 1 12 0v5l1 2H5l1-2zM9 19a3 3 0 0 0 6 0" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round" stroke-linecap="round"/></svg></button>
      <div class="av">B</div>
    </header>

    <div class="page">
      <p class="crumb">Workspace · <span>Applications</span></p>
      <div class="head">
        <h1>Applications</h1>
        <div class="actions">
          <button class="btn">+ View</button>
          <button class="btn">Filter</button>
          <button class="btn">Sort</button>
          <button class="btn primary">+ New</button>
        </div>
      </div>

      <div class="kpis">
        <div class="kpi"><p class="kl">Total</p><p class="kv">24</p><p class="kd">+5 this week</p></div>
        <div class="kpi"><p class="kl">Active loops</p><p class="kv">5</p><p class="kd up">+2</p></div>
        <div class="kpi"><p class="kl">Open offers</p><p class="kv">1</p><p class="kd">Linear · reply by Fri</p></div>
        <div class="kpi"><p class="kl">Reply rate</p><p class="kv">62%</p><p class="kd up">+8% vs. last month</p></div>
      </div>

      <div class="tabs">
        <button class="tab cur">All <span>6</span></button>
        <button class="tab">Active <span>3</span></button>
        <button class="tab">Offers <span>1</span></button>
        <button class="tab">Closed <span>1</span></button>
      </div>

      <table>
        <thead>
          <tr><th>Company</th><th>Role</th><th>Status</th><th>Applied</th><th>CV</th><th>Source</th><th></th></tr>
        </thead>
        <tbody>
          {#each PREVIEW_APPS as a}
            <tr>
              <td><span class="co"><span class="mono" style="background: {tone(a.company)}">{a.company[0]}</span><span>{a.company}</span></span></td>
              <td><div class="role">{a.role}</div><div class="loc">{a.location}</div></td>
              <td><span class="pill" style="background: {STATUS[a.status].bg}; color: {STATUS[a.status].fg}">{a.status}</span></td>
              <td class="muted">{fmtDate(a.applied_at)}</td>
              <td class="muted">{a.cv_variant ?? '—'}</td>
              <td class="muted">{a.source}</td>
              <td><button class="row-x">⋯</button></td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </main>
</div>

<style>
  :global(html, body) {
    background: #ffffff; color: #0f0f10;
    font-family: 'Inter', system-ui, sans-serif;
    font-size: 13px; line-height: 1.45;
    -webkit-font-smoothing: antialiased;
    font-feature-settings: 'ss01', 'cv11';
  }

  .shell { display: grid; grid-template-columns: 220px 1fr; min-height: 100vh; }

  .side { background: #fafafa; border-right: 1px solid #ececef; padding: .65rem; display: flex; flex-direction: column; gap: .85rem; }
  .ws { display: flex; align-items: center; gap: .55rem; padding: .35rem .5rem; border-radius: 6px; }
  .ws:hover { background: #f0f0f3; }
  .ws-mark { width: 22px; height: 22px; border-radius: 5px; background: linear-gradient(135deg, #5b5cff, #a855f7); color: #fff; font-weight: 700; font-size: 11px; display: grid; place-items: center; }
  .ws-name { font-weight: 600; font-size: 13px; flex: 1; }
  .ws-x { background: transparent; border: 0; color: #71717a; cursor: pointer; font-size: 10px; }

  .new { display: flex; align-items: center; justify-content: space-between; padding: .4rem .55rem; background: #0f0f10; color: #fff; border: 0; border-radius: 6px; font: inherit; font-size: 12px; font-weight: 500; cursor: pointer; }
  .new:hover { background: #2c2c2f; }
  .new .k { background: rgba(255,255,255,.15); color: #d4d4d8; padding: 1px 5px; border-radius: 3px; font-family: 'JetBrains Mono', monospace; font-size: 10px; }

  .group { display: flex; flex-direction: column; gap: 1px; }
  .g-title { margin: .25rem 0 .15rem .55rem; font-size: 10px; font-weight: 600; color: #a1a1aa; letter-spacing: .04em; text-transform: uppercase; }
  .nl { display: flex; align-items: center; gap: .5rem; padding: .35rem .55rem; border-radius: 5px; color: #3f3f46; text-decoration: none; font-size: 13px; }
  .nl:hover { background: #f0f0f3; color: #0f0f10; }
  .nl.current { background: #ececef; color: #0f0f10; font-weight: 500; }
  .nl.current svg { color: #5b5cff; }
  .nl svg { color: #71717a; }
  .nl .ct { margin-left: auto; font-family: 'JetBrains Mono', monospace; font-size: 10px; color: #a1a1aa; }
  .nl .cdot { width: 16px; height: 16px; border-radius: 4px; color: #fff; display: grid; place-items: center; font-size: 9px; font-weight: 600; }

  .flex { flex: 1; }
  .leave { padding: .35rem .55rem; color: #a1a1aa; font-size: 11px; text-decoration: none; }

  main { display: flex; flex-direction: column; min-width: 0; }
  .top { display: flex; align-items: center; gap: .5rem; padding: .5rem .85rem; border-bottom: 1px solid #ececef; background: rgba(255,255,255,.85); backdrop-filter: blur(8px); position: sticky; top: 0; z-index: 10; }
  .search { flex: 1; max-width: 400px; display: flex; align-items: center; gap: .45rem; padding: .35rem .55rem; background: #fafafa; border: 1px solid #ececef; border-radius: 5px; color: #a1a1aa; }
  .search input { flex: 1; border: 0; background: transparent; outline: 0; font: inherit; color: #0f0f10; }
  .search .kbd { font-family: 'JetBrains Mono', monospace; font-size: 10px; background: #ececef; color: #71717a; padding: 1px 5px; border-radius: 3px; }
  .ic { width: 28px; height: 28px; border-radius: 5px; background: transparent; border: 0; color: #52525b; cursor: pointer; display: grid; place-items: center; }
  .ic:hover { background: #f0f0f3; }
  .av { width: 24px; height: 24px; border-radius: 999px; background: linear-gradient(135deg, #5b5cff, #a855f7); color: #fff; display: grid; place-items: center; font-weight: 600; font-size: 11px; }

  .page { padding: 1.5rem 2rem 3rem; max-width: 1200px; width: 100%; }
  .crumb { font-size: 12px; color: #a1a1aa; margin: 0 0 .25rem; }
  .crumb span { color: #0f0f10; }
  .head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.5rem; }
  h1 { font-size: 22px; font-weight: 600; letter-spacing: -.01em; margin: 0; }
  .actions { display: flex; gap: .35rem; }
  .btn { padding: .35rem .7rem; background: #fff; border: 1px solid #ececef; border-radius: 5px; font: inherit; font-size: 12px; color: #3f3f46; cursor: pointer; }
  .btn:hover { background: #fafafa; color: #0f0f10; }
  .btn.primary { background: #5b5cff; color: #fff; border-color: #5b5cff; font-weight: 500; }
  .btn.primary:hover { background: #4848e0; }

  .kpis { display: grid; grid-template-columns: repeat(4, 1fr); gap: .65rem; margin-bottom: 1.5rem; }
  @media (max-width: 1000px) { .kpis { grid-template-columns: repeat(2, 1fr); } }
  .kpi { background: #fff; border: 1px solid #ececef; border-radius: 8px; padding: .85rem 1rem; }
  .kl { margin: 0; font-size: 11px; color: #71717a; font-weight: 500; }
  .kv { margin: .25rem 0; font-size: 22px; font-weight: 600; letter-spacing: -.01em; color: #0f0f10; }
  .kd { margin: 0; font-size: 11px; color: #71717a; }
  .kd.up { color: #15803d; }

  .tabs { display: flex; gap: .25rem; border-bottom: 1px solid #ececef; margin-bottom: 0; }
  .tab { background: transparent; border: 0; padding: .55rem .85rem; color: #71717a; font: inherit; font-size: 12px; cursor: pointer; border-bottom: 2px solid transparent; margin-bottom: -1px; display: inline-flex; align-items: center; gap: .35rem; }
  .tab:hover { color: #0f0f10; }
  .tab.cur { color: #0f0f10; border-bottom-color: #5b5cff; font-weight: 500; }
  .tab span { font-family: 'JetBrains Mono', monospace; font-size: 10px; background: #ececef; color: #71717a; padding: 1px 5px; border-radius: 3px; }

  table { width: 100%; border-collapse: collapse; background: #fff; }
  th { text-align: left; padding: .5rem .85rem; font-size: 11px; font-weight: 500; color: #71717a; border-bottom: 1px solid #ececef; }
  td { padding: .6rem .85rem; border-bottom: 1px solid #f4f4f5; font-size: 13px; vertical-align: middle; color: #0f0f10; }
  tbody tr:hover { background: #fafafa; }
  tbody tr:last-child td { border-bottom: 0; }
  .co { display: inline-flex; align-items: center; gap: .5rem; }
  .mono { width: 22px; height: 22px; border-radius: 5px; color: #fff; display: grid; place-items: center; font-weight: 600; font-size: 11px; }
  .role { color: #0f0f10; }
  .loc { color: #a1a1aa; font-size: 11px; margin-top: 1px; }
  .pill { display: inline-block; padding: 2px 8px; border-radius: 4px; font-size: 11px; font-weight: 500; text-transform: capitalize; }
  .muted { color: #71717a; }
  .row-x { background: transparent; border: 0; color: #a1a1aa; cursor: pointer; width: 22px; height: 22px; border-radius: 4px; }
  .row-x:hover { background: #f0f0f3; color: #0f0f10; }
</style>
