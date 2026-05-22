<script>
  const directions = [
    {
      slug: 'a',
      name: 'Minimal',
      tagline: 'Linear / Notion style — gets out of the way',
      blurb:
        'Monochrome, dense, subtle borders. Almost no chrome. The data is the product; the app is just a window into it.',
      bg: '#fafafa',
      fg: '#1a1a1a',
      accent: '#2563eb'
    },
    {
      slug: 'b',
      name: 'Soft',
      tagline: 'Airtable / Notion colorful — pleasant to open daily',
      blurb:
        'Warm whites, pastel status tags, generous whitespace, rounded corners. Friendly and inviting. Makes the job search feel less grim.',
      bg: '#faf9f7',
      fg: '#37352f',
      accent: '#1f7ae0'
    },
    {
      slug: 'c',
      name: 'Terminal',
      tagline: 'Bloomberg style — for people who like dense data',
      blurb:
        'Pure black, monospace data, bright accent. Treats the funnel like a trading dashboard. Distinctive and a little intense.',
      bg: '#000',
      fg: '#d9d9d9',
      accent: '#00ff66'
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
      Three takes on Pursuit. Click each, walk through it for 30 seconds, then come back and pick one.
      The chosen direction becomes the design system.
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
    These are mockups with fake data — no sign-in needed.
    The real app stays where it is at <a href="/">/</a>.
  </footer>
</main>

<style>
  :global(html, body) {
    background: #0f1115;
    color: #e6e8ee;
  }

  .wrap {
    max-width: 1100px;
    margin: 6vh auto;
    padding: 0 2rem;
  }
  header h1 {
    font-size: 2rem;
    margin: 0 0 .5rem;
  }
  header p {
    color: #8b91a0;
    max-width: 60ch;
    margin: 0 0 2.5rem;
  }
  .grid {
    list-style: none;
    padding: 0;
    margin: 0;
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
    display: flex;
    align-items: center;
    gap: .5rem;
    font-size: .8rem;
    font-weight: 600;
    letter-spacing: .02em;
    text-transform: uppercase;
    opacity: .7;
  }
  .dot {
    width: 10px;
    height: 10px;
    border-radius: 999px;
  }
  .card h2 {
    font-size: 1.1rem;
    margin: 1rem 0 .5rem;
  }
  .card p {
    font-size: .9rem;
    opacity: .8;
    margin: 0 0 1.5rem;
  }
  .cta {
    color: var(--card-accent);
    font-weight: 500;
    font-size: .9rem;
  }
  footer {
    margin-top: 3rem;
    color: #8b91a0;
    font-size: .85rem;
  }
  footer a { color: #4f8cff; }
</style>
