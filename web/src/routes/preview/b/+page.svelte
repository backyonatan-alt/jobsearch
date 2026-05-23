<script>
  import { PREVIEW_APPS, fmtDate } from '$lib/preview-data.js';

  // SaaS-grade status: small colored dot + lowercase text.
  const STATUS_COLOR = {
    wishlist:  '#a39d92',
    applied:   '#3b82f6',
    screen:    '#f59e0b',
    interview: '#c45ba8',
    offer:     '#10b981',
    rejected:  '#ef4444',
    withdrawn: '#71717a'
  };

  // Generate distinctive monogram tones per company.
  const TONES = ['#ff8a5b', '#c45ba8', '#7e6cd6', '#5b8def', '#5bbb8a', '#d6b15b'];
  const monoTone = (s) => TONES[s.charCodeAt(0) % TONES.length];

  const metrics = [
    { label: 'Total applications', value: '24', delta: '+5 this week', spark: [3, 4, 4, 5, 5, 6, 8] },
    { label: 'Active loops',       value: '5',  delta: '+2 this week', spark: [2, 2, 3, 3, 4, 5, 5] },
    { label: 'Open offers',        value: '1',  delta: 'Linear · reply by Fri', spark: [0, 0, 0, 0, 0, 1, 1] },
    { label: 'Reply rate',         value: '62%', delta: '+8% vs. last month', spark: [40, 45, 50, 52, 55, 58, 62] }
  ];

  // Build sparkline points (0–100 → 0–24 vertical).
  function path(values) {
    const max = Math.max(...values, 1);
    const w = 80, h = 22;
    const step = w / (values.length - 1);
    return values
      .map((v, i) => `${i === 0 ? 'M' : 'L'} ${(i * step).toFixed(1)} ${(h - (v / max) * h).toFixed(1)}`)
      .join(' ');
  }
</script>

<svelte:head>
  <title>Applications — Pursuit</title>
</svelte:head>

<p class="breadcrumb"><a href="/preview/b">Workspace</a> <span>/</span> Applications</p>
<div class="page-head">
  <h1>Applications</h1>
  <div class="head-actions">
    <button class="btn-ghost">
      <svg viewBox="0 0 24 24" width="14" height="14"><path d="M4 6h16M7 12h10M10 18h4" stroke="currentColor" stroke-width="1.7" stroke-linecap="round"/></svg>
      Filter
    </button>
    <button class="btn-ghost">
      <svg viewBox="0 0 24 24" width="14" height="14"><path d="M7 4v16M3 8l4-4 4 4M17 20V4M21 16l-4 4-4-4" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"/></svg>
      Sort
    </button>
    <button class="btn-primary">+ New</button>
  </div>
</div>

<div class="metrics">
  {#each metrics as m}
    <article class="metric">
      <p class="m-label">{m.label}</p>
      <div class="m-row">
        <p class="m-value display">{m.value}</p>
        <svg class="spark" viewBox="0 0 80 22" preserveAspectRatio="none">
          <path d={path(m.spark)} stroke="url(#sparkGrad)" stroke-width="1.5" fill="none" stroke-linecap="round"/>
          <defs>
            <linearGradient id="sparkGrad" x1="0" x2="1">
              <stop offset="0%" stop-color="#ff8a5b"/>
              <stop offset="100%" stop-color="#c45ba8"/>
            </linearGradient>
          </defs>
        </svg>
      </div>
      <p class="m-delta">{m.delta}</p>
    </article>
  {/each}
</div>

<section class="card">
  <div class="card-head">
    <div class="tabs">
      <button class="tab active">All <span class="count">6</span></button>
      <button class="tab">Active <span class="count">3</span></button>
      <button class="tab">Offers <span class="count">1</span></button>
      <button class="tab">Closed <span class="count">1</span></button>
    </div>
  </div>

  <table>
    <thead>
      <tr>
        <th>Company</th>
        <th>Role</th>
        <th>Status</th>
        <th>Applied</th>
        <th>CV</th>
        <th>Source</th>
        <th></th>
      </tr>
    </thead>
    <tbody>
      {#each PREVIEW_APPS as a}
        <tr>
          <td>
            <a class="co" href={a.slug === 'anthropic' ? '/preview/b/anthropic' : '#'}>
              <span class="mono" style="background: {monoTone(a.company)}">{a.company[0]}</span>
              <span class="co-name">{a.company}</span>
            </a>
          </td>
          <td class="role">
            <p class="role-title">{a.role}</p>
            <p class="role-loc">{a.location}</p>
          </td>
          <td>
            <span class="status">
              <span class="dot" style="background: {STATUS_COLOR[a.status]}"></span>
              {a.status}
            </span>
          </td>
          <td class="t">{fmtDate(a.applied_at)}</td>
          <td class="t">{a.cv_variant ?? '—'}</td>
          <td class="t">{a.source}</td>
          <td><button class="row-action" aria-label="More">⋯</button></td>
        </tr>
      {/each}
    </tbody>
  </table>
</section>

<p class="hint">Tap <strong>Anthropic</strong> to see the dossier page — that's where this design takes you next.</p>

<style>
  .page-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 1.5rem; }
  .head-actions { display: flex; gap: .5rem; }
  .btn-ghost {
    display: inline-flex; align-items: center; gap: .35rem;
    padding: .4rem .7rem;
    background: #fbf9f4;
    border: 1px solid #ebe6dd;
    border-radius: 6px;
    font: inherit; font-size: 13px;
    color: #4a4842;
    cursor: pointer;
  }
  .btn-ghost:hover { background: #ece5d4; }
  .btn-ghost svg { color: #71717a; }
  .btn-primary {
    padding: .4rem .9rem;
    background: linear-gradient(135deg, #18181b 0%, #2c2a32 100%);
    color: #fbf9f4;
    border: 0; border-radius: 6px;
    font: inherit; font-size: 13px; font-weight: 500;
    cursor: pointer;
    box-shadow: 0 1px 2px rgba(0,0,0,.08);
  }
  .btn-primary:hover { transform: translateY(-1px); transition: transform .12s ease; }

  .metrics {
    display: grid; grid-template-columns: repeat(4, 1fr); gap: .75rem;
    margin-bottom: 2rem;
  }
  @media (max-width: 1000px) { .metrics { grid-template-columns: repeat(2, 1fr); } }
  .metric {
    background: #fbf9f4;
    border: 1px solid #ebe6dd;
    border-radius: 10px;
    padding: 1rem 1.1rem 1.1rem;
  }
  .m-label { margin: 0; font-size: 11px; letter-spacing: .04em; color: #71717a; font-weight: 500; text-transform: uppercase; }
  .m-row { display: flex; align-items: center; justify-content: space-between; gap: 1rem; margin: .35rem 0 .25rem; }
  .m-value { font-size: 1.65rem; margin: 0; color: #18181b; }
  .spark { width: 70px; height: 22px; opacity: .9; }
  .m-delta { margin: 0; font-size: 11px; color: #71717a; }

  .card {
    background: #fbf9f4;
    border: 1px solid #ebe6dd;
    border-radius: 12px;
    overflow: hidden;
  }
  .card-head { padding: .5rem .5rem 0; border-bottom: 1px solid #ebe6dd; }
  .tabs { display: flex; gap: 0; padding: 0 .25rem; }
  .tab {
    background: transparent; border: 0;
    padding: .65rem .9rem;
    color: #71717a; font: inherit; font-size: 13px;
    cursor: pointer;
    border-bottom: 2px solid transparent;
    margin-bottom: -1px;
    display: inline-flex; gap: .4rem; align-items: center;
  }
  .tab:hover { color: #18181b; }
  .tab.active { color: #18181b; border-bottom-color: #c45ba8; font-weight: 500; }
  .tab .count {
    font-family: 'JetBrains Mono', monospace;
    font-size: 10px;
    background: #ece5d4;
    color: #6f685c;
    padding: 1px 5px;
    border-radius: 4px;
  }

  table { width: 100%; border-collapse: collapse; }
  th {
    text-align: left;
    padding: .65rem 1rem;
    font-size: 11px;
    font-weight: 500;
    letter-spacing: .04em;
    color: #71717a;
    text-transform: uppercase;
    border-bottom: 1px solid #ebe6dd;
    background: #faf7f0;
  }
  td {
    padding: .85rem 1rem;
    border-bottom: 1px solid #f1ede3;
    font-size: 13px;
    vertical-align: middle;
  }
  tbody tr:hover { background: rgba(196, 91, 168, .03); }
  tbody tr:last-child td { border-bottom: 0; }

  .co { display: inline-flex; align-items: center; gap: .65rem; text-decoration: none; color: #18181b; }
  .co:hover .co-name { color: #c45ba8; }
  .mono {
    width: 28px; height: 28px; border-radius: 8px;
    color: #fff;
    display: grid; place-items: center;
    font-weight: 600; font-size: 13px;
    box-shadow: 0 1px 2px rgba(0,0,0,.06);
  }
  .co-name { font-weight: 500; transition: color .15s ease; }

  .role-title { margin: 0; color: #18181b; font-weight: 500; }
  .role-loc { margin: .15rem 0 0; color: #a39d92; font-size: 11px; }

  .status {
    display: inline-flex; align-items: center; gap: .4rem;
    font-size: 12px;
    color: #4a4842;
    text-transform: capitalize;
  }
  .dot { width: 7px; height: 7px; border-radius: 999px; }
  .t { color: #71717a; }

  .row-action {
    background: transparent; border: 0;
    color: #a39d92; cursor: pointer;
    font-size: 16px;
    width: 24px; height: 24px;
    border-radius: 4px;
  }
  .row-action:hover { background: rgba(0,0,0,.04); color: #18181b; }

  .hint {
    margin-top: 2rem; text-align: center;
    color: #a39d92; font-size: 12px;
  }
  .hint strong { color: #c45ba8; font-weight: 600; }
</style>
