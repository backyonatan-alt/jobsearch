<script>
  // Design B — Soft (Airtable/Notion colorful). Warm whites, pastel pills,
  // generous whitespace.
  import { PREVIEW_APPS, FUNNEL, fmtDate } from '$lib/preview-data.js';

  const STATUS_STYLE = {
    wishlist:  { bg: '#ebe9e6', fg: '#5b574f' },
    applied:   { bg: '#d6e4f5', fg: '#1f4e8a' },
    screen:    { bg: '#fce8d8', fg: '#8a4a1a' },
    interview: { bg: '#e7d4f4', fg: '#5e2b8a' },
    offer:     { bg: '#d4f0d2', fg: '#1f6b1d' },
    rejected:  { bg: '#f5d8d8', fg: '#8a1f1f' },
    withdrawn: { bg: '#ebebeb', fg: '#5a5a5a' }
  };
</script>

<svelte:head>
  <title>Preview B — Soft</title>
</svelte:head>

<div class="root">
  <header class="topbar">
    <div class="brand">
      <span class="logo">P</span>
      Pursuit
    </div>
    <nav>
      <a href="#" class="navlink active">Applications</a>
      <a href="#" class="navlink">Funnel</a>
      <a href="#" class="navlink">Dossiers</a>
    </nav>
    <div class="me">
      <span class="avatar">B</span>
      back.yonatan
    </div>
  </header>

  <main>
    <div class="page-head">
      <div>
        <h1>Your applications 👋</h1>
        <p class="sub">{PREVIEW_APPS.length} in flight · {FUNNEL.interview + FUNNEL.offer} in late stage</p>
      </div>
      <button class="btn-primary">+ New application</button>
    </div>

    <div class="cards">
      {#each PREVIEW_APPS as a}
        <article class="card">
          <header>
            <div>
              <h3>{a.company}</h3>
              <p>{a.role}</p>
            </div>
            <span class="pill" style="background: {STATUS_STYLE[a.status].bg}; color: {STATUS_STYLE[a.status].fg}">
              {a.status}
            </span>
          </header>
          <footer>
            <span>📅 {fmtDate(a.applied_at)}</span>
            <span>📍 {a.location}</span>
            {#if a.cv_variant}<span>📄 {a.cv_variant}</span>{/if}
            <span class="src">via {a.source}</span>
          </footer>
        </article>
      {/each}
    </div>

    <a class="back" href="/preview">← back to directions</a>
  </main>
</div>

<style>
  :global(html, body) {
    background: #faf9f7;
    color: #37352f;
    font-family: -apple-system, 'Inter', system-ui, BlinkMacSystemFont, sans-serif;
    font-size: 15px;
    -webkit-font-smoothing: antialiased;
  }
  .root { min-height: 100vh; }

  .topbar {
    display: flex; align-items: center; gap: 1.5rem;
    padding: 1rem 2rem;
    background: #fff;
    border-bottom: 1px solid #ece9e2;
  }
  .brand {
    display: flex; align-items: center; gap: .65rem;
    font-weight: 600;
  }
  .logo {
    width: 28px; height: 28px;
    background: linear-gradient(135deg, #1f7ae0 0%, #5b6cff 100%);
    color: #fff;
    border-radius: 8px;
    display: grid; place-items: center;
    font-weight: 700;
    font-size: 14px;
  }
  nav { display: flex; gap: .25rem; margin-left: 1.5rem; }
  .navlink {
    padding: .45rem .85rem;
    color: #837e75; font-size: 14px;
    border-radius: 7px; text-decoration: none;
  }
  .navlink:hover { background: #f3f0eb; color: #37352f; }
  .navlink.active { color: #37352f; background: #f3f0eb; font-weight: 500; }
  .me {
    margin-left: auto;
    display: flex; align-items: center; gap: .5rem;
    color: #6b665d; font-size: 14px;
  }
  .avatar {
    width: 28px; height: 28px; border-radius: 999px;
    background: #e7d4f4; color: #5e2b8a;
    display: grid; place-items: center;
    font-weight: 600; font-size: 13px;
  }

  main { max-width: 1080px; margin: 0 auto; padding: 2.5rem 2rem; }

  .page-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 2rem; }
  .page-head h1 { font-size: 1.75rem; font-weight: 600; margin: 0; }
  .sub { color: #837e75; font-size: 14px; margin: .25rem 0 0; }
  .btn-primary {
    background: #1f7ae0; color: #fff;
    border: 0; border-radius: 10px;
    padding: .65rem 1.1rem;
    font-size: 14px; font-weight: 500;
    cursor: pointer;
    box-shadow: 0 1px 2px rgba(0,0,0,.05);
  }
  .btn-primary:hover { background: #1a68c1; }

  .cards { display: grid; gap: .85rem; }
  .card {
    background: #fff;
    border: 1px solid #ece9e2;
    border-radius: 14px;
    padding: 1.25rem 1.5rem;
    box-shadow: 0 1px 2px rgba(0,0,0,.03);
    transition: border-color .12s ease, transform .12s ease;
  }
  .card:hover { border-color: #d6d1c5; transform: translateY(-1px); }
  .card header {
    display: flex; align-items: flex-start; justify-content: space-between;
    gap: 1rem;
    margin-bottom: .75rem;
  }
  .card h3 { margin: 0; font-size: 1.05rem; font-weight: 600; }
  .card p { margin: .15rem 0 0; color: #6b665d; font-size: 14px; }
  .pill {
    padding: .3rem .7rem;
    border-radius: 999px;
    font-size: 12px;
    font-weight: 500;
    text-transform: capitalize;
    white-space: nowrap;
  }
  .card footer {
    display: flex; flex-wrap: wrap; gap: 1.25rem;
    color: #837e75; font-size: 13px;
  }
  .src { margin-left: auto; }

  .back {
    display: inline-block;
    margin-top: 2rem;
    color: #837e75;
    font-size: 13px;
    text-decoration: none;
  }
  .back:hover { color: #37352f; }
</style>
