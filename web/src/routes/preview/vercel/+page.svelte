<script>
  import { PREVIEW_APPS, fmtDate } from '$lib/preview-data.js';

  const STATUS_DOT = {
    wishlist:  '#a1a1a1',
    applied:   '#0070f3',
    screen:    '#f5a623',
    interview: '#7928ca',
    offer:     '#0070f3',
    rejected:  '#ee0000',
    withdrawn: '#a1a1a1'
  };
</script>

<svelte:head>
  <title>Vercel sketch — Pursuit</title>
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous" />
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500;600&display=swap" rel="stylesheet" />
</svelte:head>

<div class="shell">
  <header class="top">
    <div class="left">
      <span class="mark">▲</span>
      <span class="div">/</span>
      <span class="ws">back-yonatan</span>
      <span class="div">/</span>
      <span class="proj">pursuit</span>
    </div>
    <nav class="mid">
      <a class="t cur" href="#">Overview</a>
      <a class="t" href="#">Applications</a>
      <a class="t" href="#">Activity</a>
      <a class="t" href="#">Settings</a>
    </nav>
    <div class="right">
      <button class="ic" aria-label="Search">⌘ K</button>
      <button class="ic" aria-label="Help">?</button>
      <div class="av">b</div>
    </div>
  </header>

  <header class="sub">
    <div class="sub-inner">
      <nav class="subnav">
        <a class="sn cur" href="#">Overview</a>
        <a class="sn" href="#">Applications</a>
        <a class="sn" href="#">Board</a>
        <a class="sn" href="#">Funnel</a>
        <a class="sn" href="#">Dossiers</a>
      </nav>
      <button class="primary">+ New application</button>
    </div>
  </header>

  <main>
    <div class="page">
      <h1>Overview</h1>

      <div class="stats">
        <div class="stat">
          <p class="sl">Total applications</p>
          <p class="sv">24</p>
          <p class="sd">↑ 5 this week</p>
        </div>
        <div class="stat">
          <p class="sl">Active loops</p>
          <p class="sv">5</p>
          <p class="sd">↑ 2</p>
        </div>
        <div class="stat">
          <p class="sl">Open offers</p>
          <p class="sv">1</p>
          <p class="sd">Linear · Fri</p>
        </div>
        <div class="stat">
          <p class="sl">Reply rate</p>
          <p class="sv">62<span class="pct">%</span></p>
          <p class="sd">↑ 8%</p>
        </div>
      </div>

      <section class="block">
        <header class="bh">
          <h2>Applications</h2>
          <div class="bhr">
            <input class="search" placeholder="search by company, role…" />
            <select class="select"><option>All statuses</option></select>
          </div>
        </header>
        <table>
          <thead>
            <tr>
              <th>Company</th>
              <th>Role</th>
              <th>Status</th>
              <th>Applied</th>
              <th>CV</th>
              <th>Source</th>
            </tr>
          </thead>
          <tbody>
            {#each PREVIEW_APPS as a}
              <tr>
                <td><span class="co">{a.company}</span></td>
                <td><span class="role">{a.role}</span> <span class="loc">— {a.location}</span></td>
                <td>
                  <span class="status">
                    <span class="dot" style="background: {STATUS_DOT[a.status]}"></span>
                    {a.status}
                  </span>
                </td>
                <td class="mono">{fmtDate(a.applied_at)}</td>
                <td class="mono">{a.cv_variant ?? '—'}</td>
                <td class="mono">{a.source}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </section>

      <p class="leave-row"><a href="/preview">‹ leave preview</a></p>
    </div>
  </main>
</div>

<style>
  :global(html, body) {
    background: #ffffff; color: #000;
    font-family: 'Inter', system-ui, sans-serif;
    font-size: 13px; line-height: 1.5;
    -webkit-font-smoothing: antialiased;
    font-feature-settings: 'ss01', 'cv11';
  }

  .top { display: flex; align-items: center; gap: .5rem; padding: .65rem 1.5rem; border-bottom: 1px solid #eaeaea; }
  .left { display: flex; align-items: center; gap: .5rem; font-size: 14px; }
  .mark { color: #000; font-size: 16px; }
  .div { color: #d1d1d1; }
  .ws { color: #666; }
  .proj { font-weight: 600; }
  .mid { display: flex; gap: 0; margin-left: 1.5rem; }
  .t { padding: .5rem .75rem; color: #666; text-decoration: none; font-size: 13px; border-radius: 4px; }
  .t:hover { background: #fafafa; color: #000; }
  .t.cur { color: #000; font-weight: 500; }
  .right { margin-left: auto; display: flex; align-items: center; gap: .5rem; }
  .ic { background: #fafafa; border: 1px solid #eaeaea; padding: .25rem .5rem; border-radius: 4px; font: inherit; font-size: 11px; color: #666; cursor: pointer; font-family: 'JetBrains Mono', monospace; }
  .av { width: 28px; height: 28px; border-radius: 999px; background: #000; color: #fff; display: grid; place-items: center; font-weight: 600; font-size: 12px; font-family: 'JetBrains Mono', monospace; }

  .sub { border-bottom: 1px solid #eaeaea; }
  .sub-inner { display: flex; justify-content: space-between; align-items: center; padding: 0 1.5rem; max-width: 1280px; margin: 0 auto; }
  .subnav { display: flex; gap: 0; }
  .sn { padding: .75rem 1rem; color: #666; text-decoration: none; font-size: 13px; border-bottom: 2px solid transparent; margin-bottom: -1px; }
  .sn:hover { color: #000; }
  .sn.cur { color: #000; border-bottom-color: #000; }
  .primary { padding: .35rem .85rem; background: #000; color: #fff; border: 1px solid #000; border-radius: 6px; font: inherit; font-size: 13px; font-weight: 500; cursor: pointer; }
  .primary:hover { background: #fff; color: #000; }

  .page { max-width: 1280px; margin: 0 auto; padding: 2rem 1.5rem 4rem; }
  h1 { font-size: 32px; font-weight: 700; letter-spacing: -.02em; margin: 0 0 1.5rem; }

  .stats { display: grid; grid-template-columns: repeat(4, 1fr); gap: 1rem; margin-bottom: 2rem; }
  @media (max-width: 800px) { .stats { grid-template-columns: repeat(2, 1fr); } }
  .stat { border: 1px solid #eaeaea; border-radius: 8px; padding: 1.25rem; }
  .sl { margin: 0; font-size: 12px; color: #666; }
  .sv { margin: .25rem 0; font-size: 36px; font-weight: 700; letter-spacing: -.03em; color: #000; line-height: 1; }
  .sv .pct { font-size: 20px; color: #666; margin-left: 2px; }
  .sd { margin: 0; font-size: 11px; color: #0070f3; font-family: 'JetBrains Mono', monospace; }

  .block { border: 1px solid #eaeaea; border-radius: 8px; overflow: hidden; }
  .bh { display: flex; align-items: center; justify-content: space-between; padding: 1rem 1.25rem; border-bottom: 1px solid #eaeaea; }
  .bh h2 { margin: 0; font-size: 16px; font-weight: 600; }
  .bhr { display: flex; gap: .5rem; }
  .search, .select { padding: .35rem .65rem; border: 1px solid #eaeaea; border-radius: 6px; font: inherit; font-size: 12px; background: #fff; }
  .search { width: 220px; }
  .search::placeholder { color: #999; }

  table { width: 100%; border-collapse: collapse; }
  th { text-align: left; padding: .65rem 1.25rem; font-size: 11px; font-weight: 500; color: #666; text-transform: uppercase; letter-spacing: .04em; border-bottom: 1px solid #eaeaea; background: #fafafa; }
  td { padding: .85rem 1.25rem; border-bottom: 1px solid #f5f5f5; font-size: 13px; vertical-align: middle; }
  tbody tr:hover { background: #fafafa; }
  tbody tr:last-child td { border-bottom: 0; }
  .co { font-weight: 600; color: #000; }
  .role { color: #000; }
  .loc { color: #666; }
  .status { display: inline-flex; align-items: center; gap: .4rem; color: #000; text-transform: capitalize; }
  .dot { width: 7px; height: 7px; border-radius: 999px; }
  .mono { font-family: 'JetBrains Mono', monospace; font-size: 12px; color: #666; }

  .leave-row { margin-top: 2rem; text-align: center; }
  .leave-row a { color: #666; text-decoration: none; font-size: 12px; }
  .leave-row a:hover { color: #000; }
</style>
