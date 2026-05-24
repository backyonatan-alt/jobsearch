<script>
  const opts = [
    { slug: 'a', name: 'Unified zone, quiet hint',
      blurb: 'One paste/drop area, single subtle hint line. Lowest friction; the screenshot trick is easy to miss the first time.' },
    { slug: 'b', name: 'Three explicit chips',
      blurb: 'Mode picker at top — Paste text / Screenshot / Type it. Each mode swaps the input + help text. Zero ambiguity; one extra click for power users.' },
    { slug: 'c', name: 'Three visual cards',
      blurb: 'Big chooser at the top: 📋 Paste text · 📸 Screenshot · ✍️ Type fields. Most discoverable, friendliest for first-timers. Heaviest visually.' },
    { slug: 'd', name: 'Zone + expandable how-to',
      blurb: 'Single zone like A, but with an expandable “How to add from LinkedIn” walkthrough one click away. Clean default, real help in reach.' }
  ];
</script>

<svelte:head><title>Add-modal previews — Pursuit</title></svelte:head>

<main class="wrap">
  <header>
    <h1>Add application — UI options</h1>
    <p>Four versions of the “+ New application” modal, rendered in the live Crystalline design system. Click each, then come back and pick one (or a mix).</p>
  </header>

  <ul class="grid">
    {#each opts as o}
      <li>
        <a href="/preview/add-modal/{o.slug}" class="card">
          <div class="badge">Option {o.slug.toUpperCase()}</div>
          <h2>{o.name}</h2>
          <p>{o.blurb}</p>
          <span class="cta">Open preview →</span>
        </a>
      </li>
    {/each}
  </ul>

  <footer>
    These are mockups — Parse buttons and form Save don't actually do anything.
    <a href="/app">← back to app</a>
  </footer>
</main>

<style>
  :global(html, body) { background: var(--surface); }
  .wrap { max-width: 920px; margin: 6vh auto; padding: 0 2rem; }
  header h1 { font-size: 1.6rem; font-weight: 600; margin: 0 0 .5rem; letter-spacing: -0.015em; }
  header p { color: var(--mute); max-width: 60ch; margin: 0 0 2rem; }
  .grid { list-style: none; padding: 0; margin: 0; display: grid; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); gap: 12px; }
  .card { display: block; background: var(--card); border: 1px solid var(--rule); border-radius: 12px; padding: 1.25rem; box-shadow: var(--sh-1); text-decoration: none; color: var(--ink); transition: border-color 120ms ease, transform 120ms ease; }
  .card:hover { border-color: var(--accent); transform: translateY(-1px); }
  .badge { font-family: var(--mono); font-size: 10.5px; color: var(--accent-text); background: var(--accent-tint); border-radius: 4px; padding: 2px 7px; display: inline-block; margin-bottom: 8px; letter-spacing: .04em; }
  .card h2 { font-size: 16px; font-weight: 500; margin: 0 0 .35rem; letter-spacing: -0.01em; }
  .card p { font-size: 13px; color: var(--mute); margin: 0 0 1rem; line-height: 1.55; }
  .cta { color: var(--accent-text); font-size: 13px; font-weight: 500; }
  footer { margin-top: 3rem; color: var(--mute); font-size: 12.5px; }
  footer a { color: var(--accent-text); margin-left: 8px; }
</style>
