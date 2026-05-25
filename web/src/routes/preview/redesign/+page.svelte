<script>
  const previews = [
    {
      slug: 'today/a',
      name: 'Today · variant A (locked)',
      blurb: 'The Today direction we picked. Big narrative count cards, action grid, "What we\'re noticing" with people / pause / moon icons.'
    },
    {
      slug: 'board/a',
      name: 'Board · variant A — Refined classic',
      blurb: 'Tight kanban with real company logos, generous gaps between cards, subtle drag lift, and stale cards marked with a red dot + red applied date. The broken list/board toggle is removed.'
    },
    {
      slug: 'board/b',
      name: 'Board · variant B — Cinematic',
      blurb: 'Larger cards with a colored status accent strip on top, circular logos, salary/next-interview chips, and stale cards with a pulsing red dot + "needs nudge" tag. Drag tilts the card slightly. Higher visual personality.'
    },
    {
      slug: 'brief',
      name: 'Application Brief (was Dossier)',
      blurb: 'Renamed "Brief". Company logo, JD link back to source, AI company summary, hiring-manager LinkedIn, plus a tightened layout.'
    },
    {
      slug: 'today/b',
      name: 'Today · variant B (rejected)',
      blurb: 'Kept here for reference only.'
    },
    {
      slug: 'today',
      name: 'Today · v1 (rejected)',
      blurb: 'Kept here for reference only.'
    }
  ];
</script>

<svelte:head><title>Redesign previews — Pursuit</title></svelte:head>

<main class="wrap">
  <header>
    <h1>Redesign previews</h1>
    <p>Static mockups built from the free-run feedback (May 25 2026). Nothing here is wired — Open / Save / Generate buttons don't do anything. Pick what you want to keep and we'll wire it in.</p>
  </header>

  <ul class="grid">
    {#each previews as p}
      <li>
        <a href="/preview/redesign/{p.slug}" class="card">
          <div class="badge">Preview</div>
          <h2>{p.name}</h2>
          <p>{p.blurb}</p>
          <span class="cta">Open preview →</span>
        </a>
      </li>
    {/each}
  </ul>

  <footer>
    <a href="/app">← back to app</a>
  </footer>
</main>

<style>
  :global(html, body) { background: var(--surface); }
  .wrap { max-width: 920px; margin: 6vh auto; padding: 0 2rem; }
  header h1 { font-size: 1.6rem; font-weight: 600; margin: 0 0 .5rem; letter-spacing: -0.015em; }
  header p { color: var(--mute); max-width: 64ch; margin: 0 0 2rem; line-height: 1.55; }
  .grid { list-style: none; padding: 0; margin: 0; display: grid; grid-template-columns: repeat(auto-fit, minmax(320px, 1fr)); gap: 12px; }
  .card { display: block; background: var(--card); border: 1px solid var(--rule); border-radius: 12px; padding: 1.5rem; box-shadow: var(--sh-1); text-decoration: none; color: var(--ink); transition: border-color 120ms ease, transform 120ms ease; }
  .card:hover { border-color: var(--accent); transform: translateY(-1px); }
  .badge { font-family: var(--mono); font-size: 10.5px; color: var(--accent-text); background: var(--accent-tint); border-radius: 4px; padding: 2px 7px; display: inline-block; margin-bottom: 10px; letter-spacing: .04em; }
  .card h2 { font-size: 17px; font-weight: 500; margin: 0 0 .4rem; letter-spacing: -0.01em; }
  .card p { font-size: 13px; color: var(--mute); margin: 0 0 1rem; line-height: 1.55; }
  .cta { color: var(--accent-text); font-size: 13px; font-weight: 500; }
  footer { margin-top: 3rem; color: var(--mute); font-size: 12.5px; }
  footer a { color: var(--accent-text); }
</style>
