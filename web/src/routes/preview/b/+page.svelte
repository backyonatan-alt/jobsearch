<script>
  // Garden — botanical / dreamy. Blush + sage + plum palette, soft gradients,
  // rounded type, SVG leaf-curl ornaments at the page corners.
  import { PREVIEW_APPS, FUNNEL, fmtDate } from '$lib/preview-data.js';

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
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous" />
  <link href="https://fonts.googleapis.com/css2?family=Fraunces:opsz,wght@9..144,300;9..144,400;9..144,500;9..144,600&family=Nunito+Sans:wght@400;500;600&display=swap" rel="stylesheet" />
</svelte:head>

<div class="root">
  <!-- Corner leaf ornaments -->
  <svg class="leaf leaf-tl" viewBox="0 0 120 120" aria-hidden="true">
    <path d="M0 0 Q 80 20, 90 80 Q 70 100, 30 95 Q 20 60, 0 30 Z" fill="#cad8c0" opacity=".55"/>
    <path d="M10 15 Q 60 35, 70 75" stroke="#7e9477" stroke-width="1.5" fill="none" opacity=".6"/>
  </svg>
  <svg class="leaf leaf-br" viewBox="0 0 120 120" aria-hidden="true">
    <path d="M120 120 Q 40 100, 30 40 Q 50 20, 90 25 Q 100 60, 120 90 Z" fill="#e4d0db" opacity=".55"/>
    <path d="M110 105 Q 60 85, 50 45" stroke="#9a7a90" stroke-width="1.5" fill="none" opacity=".6"/>
  </svg>

  <header class="topbar">
    <div class="brand">
      <span class="mark">
        <svg viewBox="0 0 24 24" width="22" height="22">
          <path d="M12 2 Q 18 10, 14 18 Q 10 14, 12 2 Z" fill="#6b3d5f"/>
          <path d="M12 22 Q 6 14, 10 6 Q 14 10, 12 22 Z" fill="#9a7a90"/>
        </svg>
      </span>
      <span class="wordmark">Pursuit</span>
    </div>
    <nav>
      <a href="#" class="navlink current">Garden</a>
      <a href="#" class="navlink">Funnel</a>
      <a href="#" class="navlink">Dossiers</a>
    </nav>
    <div class="me">
      <span class="avatar">B</span>
    </div>
  </header>

  <main>
    <section class="hero">
      <p class="eyebrow">Tending {PREVIEW_APPS.length} seedlings</p>
      <h1>Your garden, this season</h1>
      <p class="lede">
        One flowering, one fruited, three sown. A quiet week — tend to the screening
        round on Vercel, and gently follow up with Stripe.
      </p>
    </section>

    <ul class="grid">
      {#each PREVIEW_APPS as a}
        <li class="card" style="--card-bg: {STATUS_STYLE[a.status].bg};">
          <div class="card-inner">
            <div class="tag" style="color: {STATUS_STYLE[a.status].fg}">
              <span class="tag-dot" style="background: {STATUS_STYLE[a.status].fg}"></span>
              {STATUS_STYLE[a.status].tag}
            </div>
            <h3>{a.company}</h3>
            <p class="role">{a.role}</p>
            <div class="meta">
              <span>{a.location}</span>
              <span class="dot">·</span>
              <span>{fmtDate(a.applied_at)}</span>
              {#if a.cv_variant}
                <span class="dot">·</span>
                <span>{a.cv_variant}</span>
              {/if}
            </div>
          </div>
        </li>
      {/each}
    </ul>

    <a class="back" href="/preview">‹ back to directions</a>
  </main>
</div>

<style>
  :global(html, body) {
    background:
      radial-gradient(circle at 15% 5%, #efe1ea 0%, transparent 35%),
      radial-gradient(circle at 90% 100%, #dee5d5 0%, transparent 40%),
      #f4ece2;
    background-attachment: fixed;
    color: #33302d;
    font-family: 'Nunito Sans', -apple-system, BlinkMacSystemFont, sans-serif;
    font-size: 15px;
    -webkit-font-smoothing: antialiased;
  }

  .root { position: relative; min-height: 100vh; }

  .leaf { position: fixed; width: 240px; height: 240px; pointer-events: none; z-index: 0; }
  .leaf-tl { top: -40px; left: -40px; }
  .leaf-br { bottom: -40px; right: -40px; transform: scaleX(-1); }

  .topbar {
    position: relative; z-index: 2;
    display: flex; align-items: center; gap: 2rem;
    padding: 1.5rem 2.5rem;
  }
  .brand { display: flex; align-items: center; gap: .65rem; }
  .mark { display: inline-flex; }
  .wordmark {
    font-family: 'Fraunces', Georgia, serif;
    font-weight: 500;
    font-size: 1.4rem;
    letter-spacing: -.01em;
    color: #6b3d5f;
  }
  nav { display: flex; gap: .25rem; margin-left: 1.5rem; }
  .navlink {
    padding: .5rem 1rem;
    color: #837e75; font-size: 14px;
    border-radius: 999px;
    text-decoration: none;
    font-weight: 500;
  }
  .navlink:hover { background: rgba(107, 61, 95, .08); color: #6b3d5f; }
  .navlink.current { background: rgba(107, 61, 95, .12); color: #6b3d5f; }
  .me { margin-left: auto; }
  .avatar {
    width: 36px; height: 36px; border-radius: 999px;
    background: linear-gradient(135deg, #e4d0db, #cab4c3);
    color: #6b3d5f;
    display: grid; place-items: center;
    font-family: 'Fraunces', serif;
    font-weight: 500; font-size: 16px;
  }

  main { position: relative; z-index: 1; max-width: 980px; margin: 0 auto; padding: 1.5rem 2.5rem 4rem; }

  .hero { margin: 2rem 0 3rem; }
  .eyebrow {
    font-family: 'Nunito Sans', sans-serif;
    font-size: 12px;
    letter-spacing: .14em;
    text-transform: uppercase;
    color: #9a7a90;
    margin: 0;
    font-weight: 600;
  }
  .hero h1 {
    font-family: 'Fraunces', Georgia, serif;
    font-weight: 400;
    font-size: 2.75rem;
    line-height: 1.1;
    margin: .5rem 0 1rem;
    letter-spacing: -.02em;
    color: #2d2a27;
  }
  .lede {
    color: #5a5550;
    font-size: 1.1rem;
    line-height: 1.55;
    margin: 0;
    max-width: 56ch;
  }

  .grid {
    list-style: none; padding: 0; margin: 0;
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 1rem;
  }
  .card {
    background: var(--card-bg);
    border-radius: 20px;
    padding: 1.5rem;
    box-shadow:
      0 1px 2px rgba(60, 50, 50, .04),
      0 8px 24px rgba(60, 50, 50, .04);
    transition: transform .15s ease, box-shadow .15s ease;
  }
  .card:hover {
    transform: translateY(-3px);
    box-shadow:
      0 2px 4px rgba(60, 50, 50, .06),
      0 14px 32px rgba(60, 50, 50, .08);
  }
  .card-inner { background: rgba(255, 255, 255, .55); border-radius: 14px; padding: 1.25rem; backdrop-filter: blur(4px); }
  .tag {
    display: inline-flex; align-items: center; gap: .4rem;
    font-size: 11px;
    letter-spacing: .12em;
    text-transform: uppercase;
    font-weight: 600;
    margin-bottom: .75rem;
  }
  .tag-dot { width: 6px; height: 6px; border-radius: 999px; }
  .card h3 {
    font-family: 'Fraunces', Georgia, serif;
    font-weight: 500;
    font-size: 1.35rem;
    letter-spacing: -.01em;
    margin: 0;
    color: #2d2a27;
  }
  .role { color: #5a5550; margin: .25rem 0 1rem; font-size: 14px; }
  .meta { display: flex; flex-wrap: wrap; gap: .35rem; font-size: 12px; color: #7e7770; }
  .meta .dot { color: #b0a8a0; }

  .back {
    display: inline-block;
    margin-top: 2.5rem;
    color: #6b3d5f;
    font-size: 13px;
    text-decoration: none;
    letter-spacing: .04em;
  }
  .back:hover { color: #4a2940; }
</style>
