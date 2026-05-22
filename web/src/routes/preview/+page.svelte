<script>
  const directions = [
    {
      slug: 'a',
      name: 'Atelier',
      tagline: 'Editorial — your job search as a curated index',
      blurb:
        'Cream paper, serif headlines, oxblood ink. Each application reads like an entry in a hand-set publication. Slow, considered, design-forward.',
      bg: '#f6f0e4',
      fg: '#1d1a14',
      accent: '#7a1f2f'
    },
    {
      slug: 'b',
      name: 'Garden',
      tagline: 'Botanical — soft, dreamy, plant-energy',
      blurb:
        'Blush + sage + plum, soft gradient washes, rounded type, leaf-curl ornaments at the page edges. Makes opening the app feel like a Sunday morning.',
      bg: '#f4ece2',
      fg: '#33302d',
      accent: '#6b3d5f'
    },
    {
      slug: 'c',
      name: 'Letter',
      tagline: 'Correspondence — like a hand-written journal',
      blurb:
        'Ruled cream paper, fountain-pen blue ink, status as wax seals. Serif body for the content, monospace for meta. Quiet, intimate, almost analogue.',
      bg: '#efe7d4',
      fg: '#1f1c17',
      accent: '#1c3a72'
    }
  ];
</script>

<svelte:head>
  <title>Pursuit — design previews</title>
</svelte:head>

<main class="wrap">
  <header>
    <h1>Pick a direction</h1>
    <p>
      Round three — each direction is now a full four-screen experience:
      a dashboard, a board view, a funnel/stats page, and a deep dossier
      page for the Anthropic application. Click around inside each one
      before deciding.
    </p>
  </header>

  <ul class="grid">
    {#each directions as d}
      <li>
        <a href="/preview/{d.slug}" class="card" style="--card-bg: {d.bg}; --card-fg: {d.fg}; --card-accent: {d.accent}">
          <div class="swatch">
            <span class="dot" style="background: {d.accent}"></span>
            <span class="label">{d.name}</span>
          </div>
          <h2>{d.tagline}</h2>
          <p>{d.blurb}</p>
          <span class="cta">Open preview →</span>
        </a>
      </li>
    {/each}
  </ul>

  <footer>
    Mockups with the same six fake applications — only the styling changes.
    The real app stays at <a href="/">/</a>.
  </footer>
</main>

<style>
  :global(html, body) {
    background: #0f1115;
    color: #e6e8ee;
  }

  .wrap { max-width: 1100px; margin: 6vh auto; padding: 0 2rem; }
  header h1 { font-size: 2rem; margin: 0 0 .5rem; }
  header p { color: #8b91a0; max-width: 60ch; margin: 0 0 2.5rem; }
  .grid {
    list-style: none; padding: 0; margin: 0;
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 1rem;
  }
  .card {
    display: block;
    text-decoration: none;
    color: var(--card-fg);
    background: var(--card-bg);
    border: 1px solid #262a33;
    border-radius: 14px;
    padding: 1.5rem;
    min-height: 240px;
    transition: transform .12s ease, border-color .12s ease;
  }
  .card:hover {
    transform: translateY(-2px);
    border-color: var(--card-accent);
  }
  .swatch {
    display: flex; align-items: center; gap: .5rem;
    font-size: .8rem; font-weight: 600; letter-spacing: .04em;
    text-transform: uppercase;
    opacity: .75;
  }
  .dot { width: 10px; height: 10px; border-radius: 999px; }
  .card h2 { font-size: 1.15rem; margin: 1rem 0 .5rem; font-weight: 600; }
  .card p { font-size: .9rem; opacity: .85; margin: 0 0 1.5rem; line-height: 1.5; }
  .cta { color: var(--card-accent); font-weight: 500; font-size: .9rem; }
  footer { margin-top: 3rem; color: #8b91a0; font-size: .85rem; }
  footer a { color: #4f8cff; }
</style>
