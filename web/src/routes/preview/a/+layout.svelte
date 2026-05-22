<script>
  import { page } from '$app/state';
  let { children } = $props();

  const links = [
    { href: '/preview/a',          label: 'Applications' },
    { href: '/preview/a/board',    label: 'Board' },
    { href: '/preview/a/funnel',   label: 'Funnel' },
    { href: '/preview/a/anthropic', label: 'Dossier · Anthropic' }
  ];

  const issue = String(new Date().getFullYear() - 2000).padStart(2, '0');
  const month = new Date().toLocaleString('en-US', { month: 'long' });

  function isCurrent(href) {
    if (href === '/preview/a') return page.url.pathname === '/preview/a';
    return page.url.pathname.startsWith(href);
  }
</script>

<svelte:head>
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous" />
  <link href="https://fonts.googleapis.com/css2?family=Cormorant+Garamond:ital,wght@0,400;0,500;0,600;1,400;1,500&family=Inter:wght@400;500;600&display=swap" rel="stylesheet" />
</svelte:head>

<div class="paper">
  <header class="masthead">
    <div class="row">
      <span class="vol">Vol. {issue} · {month}</span>
      <span class="me">back.yonatan@gmail.com · <a href="/preview" class="leave">leave preview</a></span>
    </div>
    <h1>Pursuit</h1>
    <p class="subhead">A quarterly index of applications, conversations, and outcomes</p>
    <nav>
      {#each links as link, i}
        {#if i > 0}<span class="sep">·</span>{/if}
        <a href={link.href} class="navlink" class:current={isCurrent(link.href)}>{link.label}</a>
      {/each}
    </nav>
  </header>

  {@render children()}
</div>

<style>
  :global(html, body) {
    background: #efe6d4;
    color: #1d1a14;
    font-family: 'Cormorant Garamond', 'Iowan Old Style', 'Palatino Linotype', Georgia, serif;
    font-size: 17px;
    line-height: 1.55;
    -webkit-font-smoothing: antialiased;
  }

  .paper {
    max-width: 720px;
    margin: 0 auto;
    padding: 4rem 2rem 6rem;
    background:
      radial-gradient(circle at 20% 10%, rgba(122,31,47,.02), transparent 40%),
      radial-gradient(circle at 80% 90%, rgba(28,58,114,.02), transparent 40%),
      #f6f0e4;
    min-height: 100vh;
    box-shadow: 0 0 80px rgba(0,0,0,.06);
  }

  .masthead { text-align: center; border-bottom: 1px solid #1d1a14; padding-bottom: 1.5rem; margin-bottom: 2.5rem; }
  .row {
    display: flex; justify-content: space-between;
    font-family: 'Inter', sans-serif; font-size: 11px;
    letter-spacing: .12em; text-transform: uppercase;
    color: #6b6358;
    margin-bottom: 2rem;
  }
  .me .leave { color: #7a1f2f; text-decoration: none; margin-left: .5rem; }
  .me .leave:hover { text-decoration: underline; }
  .masthead h1 {
    font-size: 4rem;
    font-weight: 500;
    font-style: italic;
    margin: 0;
    line-height: 1;
    letter-spacing: -.02em;
  }
  .subhead {
    font-size: 1rem;
    color: #6b6358;
    margin: .75rem 0 1.5rem;
    font-style: italic;
  }
  :global(.paper nav) {
    display: flex; justify-content: center; align-items: center; gap: .75rem;
    flex-wrap: wrap;
    font-family: 'Inter', sans-serif; font-size: 12px;
    letter-spacing: .1em; text-transform: uppercase;
  }
  :global(.paper .navlink) { color: #6b6358; text-decoration: none; }
  :global(.paper .navlink.current) { color: #7a1f2f; font-weight: 600; }
  :global(.paper .navlink:hover) { color: #1d1a14; }
  :global(.paper .sep) { color: #cdc4b3; }

  /* Shared content helpers for child pages */
  :global(.paper h2.section-title) {
    font-size: 2rem; font-weight: 500; font-style: italic;
    margin: 0 0 .5rem; letter-spacing: -.01em;
  }
  :global(.paper .lede-text) {
    font-size: 1.1rem; font-style: italic; color: #2c2820;
    margin: 0 0 2rem;
  }
  :global(.paper em) { color: #7a1f2f; font-style: italic; }
</style>
