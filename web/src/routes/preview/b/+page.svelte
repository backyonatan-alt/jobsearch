<script>
  import { PREVIEW_APPS, fmtDate } from '$lib/preview-data.js';

  const STATUS_STYLE = {
    wishlist:  { bg: 'linear-gradient(135deg, #efe7d8, #e6dcc9)', fg: '#6b5b3d', tag: 'seedling' },
    applied:   { bg: 'linear-gradient(135deg, #e2d5f0, #d4c5e8)', fg: '#4a2f6e', tag: 'sown' },
    screen:    { bg: 'linear-gradient(135deg, #f5d8d8, #f0c7c7)', fg: '#8a3a3a', tag: 'budding' },
    interview: { bg: 'linear-gradient(135deg, #f8c89d, #f4b577)', fg: '#6e3b14', tag: 'flowering' },
    offer:     { bg: 'linear-gradient(135deg, #c2dbb8, #a8c9a0)', fg: '#2e5b3a', tag: 'fruited' },
    rejected:  { bg: 'linear-gradient(135deg, #d6cfc4, #c5beb1)', fg: '#5a544a', tag: 'fallow' },
    withdrawn: { bg: 'linear-gradient(135deg, #e8e3da, #d8d2c5)', fg: '#6b6358', tag: 'pruned' }
  };
</script>

<svelte:head>
  <title>Garden — Pursuit</title>
</svelte:head>

<section class="hero">
  <p class="eyebrow">Tending {PREVIEW_APPS.length} seedlings</p>
  <h1>Your garden, this season</h1>
  <p class="lede">
    One flowering, one fruited, three sown. A quiet week — tend to the screening
    round with <em>Vercel</em>, and gently follow up with <em>Stripe</em>.
  </p>
</section>

<ul class="grid">
  {#each PREVIEW_APPS as a}
    <li>
      <a class="card-link" href={a.slug === 'anthropic' ? '/preview/b/anthropic' : '#'}>
        <article class="card" style="--card-bg: {STATUS_STYLE[a.status].bg};">
          <div class="card-inner">
            <div class="tag" style="color: {STATUS_STYLE[a.status].fg}">
              <span class="tag-dot" style="background: {STATUS_STYLE[a.status].fg}"></span>
              {STATUS_STYLE[a.status].tag}
            </div>
            <h3>{a.company}</h3>
            <p class="role">{a.role}</p>
            <div class="meta">
              <span>{a.location}</span>
              <span class="d">·</span>
              <span>{fmtDate(a.applied_at)}</span>
              {#if a.cv_variant}
                <span class="d">·</span>
                <span>{a.cv_variant}</span>
              {/if}
            </div>
          </div>
        </article>
      </a>
    </li>
  {/each}
</ul>

<p class="hint">Tap on <em>Anthropic</em> to open the dossier — that page shows where this design takes you.</p>

<style>
  .hero { margin: 2rem 0 3rem; }
  .hero h1 { font-size: 2.75rem; line-height: 1.1; margin: .5rem 0 1rem; }
  .lede { color: #5a5550; font-size: 1.1rem; line-height: 1.55; margin: 0; max-width: 56ch; }
  .lede em { color: #6b3d5f; font-style: italic; font-weight: 500; }

  .grid { list-style: none; padding: 0; margin: 0; display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 1rem; }
  .card-link { text-decoration: none; color: inherit; }
  .card {
    background: var(--card-bg);
    border-radius: 20px;
    padding: 1.5rem;
    box-shadow: 0 1px 2px rgba(60, 50, 50, .04), 0 8px 24px rgba(60, 50, 50, .04);
    transition: transform .15s ease, box-shadow .15s ease;
  }
  .card:hover {
    transform: translateY(-3px);
    box-shadow: 0 2px 4px rgba(60, 50, 50, .06), 0 14px 32px rgba(60, 50, 50, .08);
  }
  .card-inner { background: rgba(255, 255, 255, .55); border-radius: 14px; padding: 1.25rem; backdrop-filter: blur(4px); }
  .tag { display: inline-flex; align-items: center; gap: .4rem; font-size: 11px; letter-spacing: .12em; text-transform: uppercase; font-weight: 600; margin-bottom: .75rem; }
  .tag-dot { width: 6px; height: 6px; border-radius: 999px; }
  .card h3 { font-size: 1.35rem; margin: 0; }
  .role { color: #5a5550; margin: .25rem 0 1rem; font-size: 14px; }
  .meta { display: flex; flex-wrap: wrap; gap: .35rem; font-size: 12px; color: #7e7770; }
  .meta .d { color: #b0a8a0; }

  .hint {
    margin-top: 2.5rem; text-align: center;
    color: #837e75; font-size: 13px;
  }
  .hint em { color: #6b3d5f; font-style: italic; font-weight: 600; }
</style>
