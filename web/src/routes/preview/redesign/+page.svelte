<script>
  const pages = [
    { slug: 'ia', name: 'IA rationale', blurb: 'One home, one list, exits included — the before/after map and every decision with the feedback that drove it.' },
    { slug: 'home', name: 'Home (D v3 — top nav)', blurb: 'Sidebar gone. Pipeline pulse with logos, one noticing card, the one-line everything list with exits at the end.' },
    { slug: 'applications', name: 'Applications', blurb: 'The full list grouped by stage — no board, no chips. "No longer in play" archive at the end with Reopen.' },
    { slug: 'detail', name: 'Application detail', blurb: 'Live playbook tabs + the "Not them?" wrong-company flow (click it — all 3 states work). Process, people, activity in the right column.' },
    { slug: 'insights', name: 'Insights', blurb: 'Mixpanel-grade funnel with drop-off pills, sources with conversion, noticing history.' },
    { slug: 'states', name: 'States', blurb: 'Empty home (first run) + at-credit-limit gate.' }
  ];
</script>

<svelte:head><title>Redesign previews — Pursuit</title></svelte:head>

<main class="wrap">
  <header>
    <h1>Tracker redesign — preview</h1>
    <p>Pixel port of the approved Claude Design handoff (Jul 14). Static fixtures, no backend needed. The dark banner on each page is preview chrome — it won't ship.</p>
  </header>

  <ul class="grid">
    {#each pages as p (p.slug)}
      <li>
        <a href="/preview/redesign/{p.slug}" class="card">
          <h2>{p.name}</h2>
          <p>{p.blurb}</p>
          <span class="cta">Open preview →</span>
        </a>
      </li>
    {/each}
  </ul>

  <footer>
    These are mockups — buttons don't persist anything. <a href="/app">← back to app</a>
  </footer>
</main>

<style>
  .wrap { max-width: 860px; margin: 0 auto; padding: 48px 24px 80px; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; color: #16181c; }
  header h1 { font-size: 28px; font-weight: 700; letter-spacing: -0.02em; margin: 0 0 8px; }
  header p { color: #4b5158; font-size: 14px; margin: 0 0 32px; max-width: 560px; line-height: 1.55; }
  .grid { list-style: none; margin: 0; padding: 0; display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }
  .card { display: block; border: 1px solid #e8e8e5; border-radius: 12px; padding: 18px 20px; color: inherit; text-decoration: none; background: #fff; }
  .card:hover { border-color: #b9c6e8; }
  .card h2 { font-size: 16px; font-weight: 700; margin: 0 0 6px; }
  .card p { font-size: 13px; color: #4b5158; line-height: 1.5; margin: 0 0 10px; }
  .card .cta { font-size: 12.5px; font-weight: 600; color: #2463eb; }
  footer { margin-top: 36px; font-size: 13px; color: #8a9099; }
  footer a { color: #2463eb; text-decoration: none; }
</style>
