<script>
  import { page } from '$app/state';
  let { children } = $props();

  const links = [
    { href: '/preview/b',           label: 'Garden' },
    { href: '/preview/b/board',     label: 'Board' },
    { href: '/preview/b/funnel',    label: 'Funnel' },
    { href: '/preview/b/anthropic', label: 'Anthropic' }
  ];

  function isCurrent(href) {
    if (href === '/preview/b') return page.url.pathname === '/preview/b';
    return page.url.pathname.startsWith(href);
  }
</script>

<svelte:head>
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous" />
  <link href="https://fonts.googleapis.com/css2?family=Fraunces:opsz,wght@9..144,300;9..144,400;9..144,500;9..144,600&family=Nunito+Sans:wght@400;500;600;700&display=swap" rel="stylesheet" />
</svelte:head>

<div class="root">
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
      {#each links as link}
        <a href={link.href} class="navlink" class:current={isCurrent(link.href)}>{link.label}</a>
      {/each}
    </nav>
    <div class="me">
      <a href="/preview" class="leave">leave preview</a>
      <span class="avatar">B</span>
    </div>
  </header>

  <main>
    {@render children()}
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

  .topbar { position: relative; z-index: 2; display: flex; align-items: center; gap: 1.5rem; padding: 1.5rem 2.5rem; }
  .brand { display: flex; align-items: center; gap: .65rem; }
  .wordmark { font-family: 'Fraunces', Georgia, serif; font-weight: 500; font-size: 1.4rem; letter-spacing: -.01em; color: #6b3d5f; }
  nav { display: flex; gap: .15rem; margin-left: 1rem; flex-wrap: wrap; }
  .navlink { padding: .5rem 1rem; color: #837e75; font-size: 14px; border-radius: 999px; text-decoration: none; font-weight: 500; transition: background .15s ease, color .15s ease; }
  .navlink:hover { background: rgba(107, 61, 95, .08); color: #6b3d5f; }
  .navlink.current { background: rgba(107, 61, 95, .14); color: #6b3d5f; }
  .me { margin-left: auto; display: flex; align-items: center; gap: 1rem; }
  .leave { color: #837e75; font-size: 12px; text-decoration: none; }
  .leave:hover { color: #6b3d5f; }
  .avatar { width: 36px; height: 36px; border-radius: 999px; background: linear-gradient(135deg, #e4d0db, #cab4c3); color: #6b3d5f; display: grid; place-items: center; font-family: 'Fraunces', serif; font-weight: 500; font-size: 16px; }

  main { position: relative; z-index: 1; max-width: 980px; margin: 0 auto; padding: 1rem 2.5rem 4rem; }

  :global(.root h1) {
    font-family: 'Fraunces', Georgia, serif;
    font-weight: 400;
    letter-spacing: -.02em;
    color: #2d2a27;
  }
  :global(.root h2), :global(.root h3) {
    font-family: 'Fraunces', Georgia, serif;
    font-weight: 500;
    letter-spacing: -.01em;
    color: #2d2a27;
  }
  :global(.root .eyebrow) {
    font-family: 'Nunito Sans', sans-serif;
    font-size: 12px; letter-spacing: .14em; text-transform: uppercase;
    color: #9a7a90; margin: 0; font-weight: 600;
  }
</style>
