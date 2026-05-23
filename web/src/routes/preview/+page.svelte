<script>
  const directions = [
    {
      slug: 'attio',
      name: 'Attio / Plain',
      tagline: 'Pure white, sharp, fast',
      blurb:
        'Hairline borders, Inter at small sizes, indigo accent, dense data. The 2026 default for serious SaaS. Cool and professional.',
      bg: '#ffffff',
      fg: '#0f0f10',
      accent: '#5b5cff'
    },
    {
      slug: 'cal',
      name: 'Cal.com / Linear Light',
      tagline: 'Friendly, warm, approachable',
      blurb:
        'Off-white background, generous rounded corners, soft shadows, indigo accent. Friendlier than Attio but still crisp and modern.',
      bg: '#fafafa',
      fg: '#0a0a0a',
      accent: '#4f46e5'
    },
    {
      slug: 'mercury',
      name: 'Mercury — editorial',
      tagline: 'Cream + serif headlines, considered',
      blurb:
        'Cream-tinted off-white, Inter for UI + Tiempos-style serif for display numbers. Subtle gradients in summary cards. More personality, slightly editorial.',
      bg: '#faf9f7',
      fg: '#1a1814',
      accent: '#7c3aed'
    },
    {
      slug: 'vercel',
      name: 'Vercel / Geist',
      tagline: 'Minimal, geometric, dev-tool',
      blurb:
        'Pure white, Geist sans + mono throughout, monochrome with one bright accent. Dense, sharp, very current. The "shipped by engineers" look.',
      bg: '#ffffff',
      fg: '#000000',
      accent: '#ff5400'
    }
  ];
</script>

<svelte:head>
  <title>Pursuit — design directions</title>
</svelte:head>

<main class="wrap">
  <header>
    <h1>Pick a direction</h1>
    <p>
      Four flavors of modern light SaaS, each a single-page sketch of the
      dashboard. Click through, then tell me which one feels right. I'll build
      the rest of the app in the chosen direction.
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
          <span class="cta">Open sketch →</span>
        </a>
      </li>
    {/each}
  </ul>

  <footer>
    Each sketch only renders the dashboard — once you pick, I build out Board,
    Funnel, and Dossier in the same system. The real app at <a href="/">/</a>
    still works.
  </footer>
</main>

<style>
  :global(html, body) {
    background: #0a0a0a;
    color: #e6e6e6;
    font-family: 'Inter', system-ui, sans-serif;
  }

  .wrap { max-width: 1100px; margin: 6vh auto; padding: 0 2rem; }
  header h1 { font-size: 2rem; margin: 0 0 .5rem; font-weight: 600; }
  header p { color: #a0a0a0; max-width: 60ch; margin: 0 0 2.5rem; line-height: 1.55; }
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
    border: 1px solid #1f1f23;
    border-radius: 14px;
    padding: 1.5rem;
    min-height: 240px;
    transition: transform .12s ease, border-color .12s ease;
  }
  .card:hover { transform: translateY(-2px); border-color: var(--card-accent); }
  .swatch { display: flex; align-items: center; gap: .5rem; font-size: .75rem; font-weight: 600; letter-spacing: .04em; text-transform: uppercase; opacity: .7; }
  .dot { width: 10px; height: 10px; border-radius: 999px; }
  .card h2 { font-size: 1.1rem; margin: 1rem 0 .5rem; font-weight: 600; }
  .card p { font-size: .9rem; opacity: .8; margin: 0 0 1.5rem; line-height: 1.5; }
  .cta { color: var(--card-accent); font-weight: 500; font-size: .9rem; }
  footer { margin-top: 3rem; color: #707070; font-size: .85rem; line-height: 1.55; }
  footer a { color: #5b5cff; }
</style>
