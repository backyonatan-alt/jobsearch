<script>
  // Design A — Minimal (Linear/Notion style). Light background, tiny status
  // dots, dense table, almost no chrome.
  import { PREVIEW_APPS, FUNNEL, fmtDate } from '$lib/preview-data.js';

  const STATUS_DOT = {
    wishlist: '#94a3b8',
    applied: '#3b82f6',
    screen: '#0ea5e9',
    interview: '#a855f7',
    offer: '#10b981',
    rejected: '#ef4444',
    withdrawn: '#64748b'
  };
</script>

<svelte:head>
  <title>Preview A — Minimal</title>
</svelte:head>

<div class="root">
  <header class="topbar">
    <div class="brand">
      <span class="logo"></span>
      Pursuit
    </div>
    <nav>
      <button class="navlink active">Applications</button>
      <button class="navlink">Funnel</button>
      <button class="navlink">Dossiers</button>
    </nav>
    <div class="me">back.yonatan@gmail.com</div>
  </header>

  <main>
    <div class="page-head">
      <div>
        <h1>Applications</h1>
        <p class="sub">{PREVIEW_APPS.length} tracked · {FUNNEL.interview + FUNNEL.offer} in active loop</p>
      </div>
      <div class="actions">
        <button class="btn-ghost">⌘N New</button>
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
        </tr>
      </thead>
      <tbody>
        {#each PREVIEW_APPS as a}
          <tr>
            <td class="company">{a.company}</td>
            <td>{a.role}</td>
            <td>
              <span class="status">
                <span class="dot" style="background: {STATUS_DOT[a.status]}"></span>
                {a.status}
              </span>
            </td>
            <td class="muted">{fmtDate(a.applied_at)}</td>
            <td class="muted">{a.cv_variant ?? '—'}</td>
            <td class="muted">{a.source}</td>
          </tr>
        {/each}
      </tbody>
    </table>

    <a class="back" href="/preview">← back to directions</a>
  </main>
</div>

<style>
  :global(html, body) {
    background: #fafafa;
    color: #1a1a1a;
    font-family: 'Inter', ui-sans-serif, system-ui, -apple-system, sans-serif;
    font-size: 13px;
    -webkit-font-smoothing: antialiased;
  }
  .root { min-height: 100vh; }

  .topbar {
    display: flex; align-items: center; gap: 1.5rem;
    padding: .65rem 1.5rem;
    border-bottom: 1px solid #ececec;
    background: #fff;
  }
  .brand {
    display: flex; align-items: center; gap: .5rem;
    font-weight: 600; font-size: 13px;
  }
  .logo {
    width: 14px; height: 14px;
    background: #1a1a1a;
    border-radius: 3px;
  }
  nav { display: flex; gap: .25rem; margin-left: 1rem; }
  .navlink {
    background: transparent; border: 0;
    padding: .35rem .65rem;
    color: #6b7280; font-size: 13px;
    border-radius: 5px; cursor: pointer;
  }
  .navlink:hover { background: #f3f3f3; color: #1a1a1a; }
  .navlink.active { color: #1a1a1a; background: #f3f3f3; }
  .me {
    margin-left: auto;
    color: #6b7280; font-size: 12px;
  }

  main { max-width: 1080px; margin: 0 auto; padding: 1.5rem; }

  .page-head { display: flex; align-items: flex-end; justify-content: space-between; margin-bottom: 1rem; }
  .page-head h1 { font-size: 1.15rem; font-weight: 600; margin: 0; }
  .sub { color: #6b7280; font-size: 12px; margin: .15rem 0 0; }
  .btn-ghost {
    background: #fff; border: 1px solid #e5e5e5;
    padding: .35rem .65rem; border-radius: 5px;
    color: #1a1a1a; font-size: 12px; cursor: pointer;
  }
  .btn-ghost:hover { background: #f6f6f6; }

  table {
    width: 100%;
    border-collapse: collapse;
    background: #fff;
    border: 1px solid #ececec;
    border-radius: 6px;
    overflow: hidden;
  }
  th, td {
    padding: .55rem .85rem;
    text-align: left;
    border-bottom: 1px solid #f3f3f3;
  }
  th {
    font-weight: 500;
    color: #6b7280;
    font-size: 11px;
    text-transform: uppercase;
    letter-spacing: .04em;
    background: #fbfbfb;
  }
  tbody tr:hover { background: #fbfbfb; }
  tbody tr:last-child td { border-bottom: 0; }
  td.company { font-weight: 500; }
  td.muted { color: #6b7280; }

  .status { display: inline-flex; align-items: center; gap: .4rem; }
  .dot { width: 7px; height: 7px; border-radius: 999px; }

  .back {
    display: inline-block;
    margin-top: 1.5rem;
    color: #6b7280;
    font-size: 12px;
    text-decoration: none;
  }
  .back:hover { color: #1a1a1a; }
</style>
