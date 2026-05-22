<script>
  // Design C — Terminal (Bloomberg style). Pure black, monospace data, sharp
  // green accent. Treats the funnel like a trading dashboard.
  import { PREVIEW_APPS, FUNNEL, fmtDate } from '$lib/preview-data.js';

  const STATUS_COLOR = {
    wishlist:  '#6b7280',
    applied:   '#d9d9d9',
    screen:    '#22d3ee',
    interview: '#ffb000',
    offer:     '#00ff66',
    rejected:  '#ef4444',
    withdrawn: '#6b7280'
  };

  const stats = [
    { key: 'applied',   label: 'APP', count: FUNNEL.applied + FUNNEL.screen + FUNNEL.interview + FUNNEL.offer + FUNNEL.rejected },
    { key: 'screen',    label: 'SCR', count: FUNNEL.screen + FUNNEL.interview + FUNNEL.offer },
    { key: 'interview', label: 'INT', count: FUNNEL.interview + FUNNEL.offer },
    { key: 'offer',     label: 'OFF', count: FUNNEL.offer }
  ];
  const conversion = stats.length > 1 ? (stats[1].count / Math.max(1, stats[0].count) * 100).toFixed(0) : 0;
  const ts = new Date().toLocaleString('en-US', { hour12: false });
</script>

<svelte:head>
  <title>Preview C — Terminal</title>
</svelte:head>

<div class="root">
  <header class="topbar">
    <div class="brand">
      <span class="caret">▮</span>
      PURSUIT
      <span class="ver">v0.1</span>
    </div>
    <div class="ticker">
      {#each stats as s}
        <span class="stat">
          <span class="stat-label">{s.label}</span>
          <span class="stat-num">{String(s.count).padStart(2, '0')}</span>
        </span>
      {/each}
      <span class="stat">
        <span class="stat-label">CONV</span>
        <span class="stat-num accent">{conversion}%</span>
      </span>
    </div>
    <div class="me">
      <span class="status-led"></span>
      back.yonatan@gmail.com · {ts}
    </div>
  </header>

  <main>
    <div class="cmdbar">
      <span class="prompt">›</span>
      <span class="cmd">applications --sort=applied_desc</span>
      <span class="hint">[N]ew · [F]unnel · [/]search · [⌘K] commands</span>
    </div>

    <table>
      <thead>
        <tr>
          <th class="num">#</th>
          <th>COMPANY</th>
          <th>ROLE</th>
          <th>STATUS</th>
          <th>APPLIED</th>
          <th>CV</th>
          <th>SOURCE</th>
        </tr>
      </thead>
      <tbody>
        {#each PREVIEW_APPS as a, i}
          <tr>
            <td class="num">{String(i + 1).padStart(2, '0')}</td>
            <td class="company">{a.company}</td>
            <td class="role">{a.role}</td>
            <td>
              <span class="status" style="color: {STATUS_COLOR[a.status]}">
                ● {a.status.toUpperCase()}
              </span>
            </td>
            <td class="mono">{fmtDate(a.applied_at)}</td>
            <td class="mono">{a.cv_variant ?? '-'}</td>
            <td class="muted">{a.source}</td>
          </tr>
        {/each}
      </tbody>
    </table>

    <a class="back" href="/preview">‹ back to directions</a>
  </main>
</div>

<style>
  :global(html, body) {
    background: #000;
    color: #d9d9d9;
    font-family: 'JetBrains Mono', Menlo, Consolas, 'Liberation Mono', monospace;
    font-size: 13px;
    -webkit-font-smoothing: antialiased;
  }
  .root { min-height: 100vh; }

  .topbar {
    display: flex; align-items: center; gap: 2rem;
    padding: .5rem 1rem;
    border-bottom: 1px solid #1f1f1f;
    background: #050505;
  }
  .brand {
    display: flex; align-items: center; gap: .5rem;
    font-weight: 700;
    color: #00ff66;
    letter-spacing: .08em;
  }
  .caret {
    color: #00ff66;
    animation: blink 1s steps(1) infinite;
  }
  @keyframes blink { 50% { opacity: 0; } }
  .ver { color: #4a4a4a; font-weight: 400; font-size: 11px; }

  .ticker { display: flex; gap: 1.25rem; }
  .stat { display: inline-flex; gap: .45rem; align-items: baseline; }
  .stat-label { color: #6b6b6b; font-size: 11px; letter-spacing: .06em; }
  .stat-num { font-weight: 700; }
  .stat-num.accent { color: #00ff66; }

  .me {
    margin-left: auto;
    display: flex; align-items: center; gap: .5rem;
    color: #6b6b6b; font-size: 11px; letter-spacing: .04em;
  }
  .status-led {
    width: 7px; height: 7px; border-radius: 999px;
    background: #00ff66;
    box-shadow: 0 0 6px #00ff66;
  }

  main { padding: 1rem; }

  .cmdbar {
    display: flex; align-items: center; gap: .65rem;
    padding: .5rem .75rem;
    margin-bottom: .75rem;
    background: #0a0a0a;
    border: 1px solid #1f1f1f;
  }
  .prompt { color: #00ff66; font-weight: 700; }
  .cmd { color: #d9d9d9; }
  .hint { margin-left: auto; color: #4a4a4a; font-size: 11px; letter-spacing: .04em; }

  table {
    width: 100%;
    border-collapse: collapse;
    border: 1px solid #1f1f1f;
    background: #050505;
  }
  th, td {
    padding: .35rem .75rem;
    text-align: left;
    border-bottom: 1px solid #131313;
    white-space: nowrap;
  }
  th {
    color: #6b6b6b;
    font-weight: 600;
    font-size: 10px;
    letter-spacing: .06em;
    border-bottom: 1px solid #1f1f1f;
    background: #0a0a0a;
  }
  tbody tr:hover { background: #0a1a0a; }
  td.num, th.num { color: #4a4a4a; width: 3ch; text-align: right; }
  td.company { color: #fff; font-weight: 600; }
  td.role { color: #d9d9d9; }
  td.mono { color: #d9d9d9; font-variant-numeric: tabular-nums; }
  td.muted { color: #6b6b6b; }
  .status { font-weight: 700; font-size: 11px; letter-spacing: .04em; }

  .back {
    display: inline-block;
    margin-top: 1.5rem;
    color: #6b6b6b;
    font-size: 12px;
    text-decoration: none;
  }
  .back:hover { color: #00ff66; }
</style>
