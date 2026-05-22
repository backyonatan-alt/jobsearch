<script>
  import { page } from '$app/state';
  let { children } = $props();

  const links = [
    { href: '/preview/c',           label: 'The letter' },
    { href: '/preview/c/board',     label: 'Pinboard' },
    { href: '/preview/c/funnel',    label: 'Ledger' },
    { href: '/preview/c/anthropic', label: 'Envelope · Anthropic' }
  ];

  function isCurrent(href) {
    if (href === '/preview/c') return page.url.pathname === '/preview/c';
    return page.url.pathname.startsWith(href);
  }
</script>

<svelte:head>
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous" />
  <link href="https://fonts.googleapis.com/css2?family=Spectral:ital,wght@0,400;0,500;0,600;1,400;1,500&family=JetBrains+Mono:wght@400;500&family=Caveat:wght@500;700&display=swap" rel="stylesheet" />
</svelte:head>

<div class="desk">
  <nav class="desk-nav">
    <div class="brand">
      <span class="seal">P</span>
      Pursuit
    </div>
    <div class="links">
      {#each links as link}
        <a href={link.href} class:current={isCurrent(link.href)}>{link.label}</a>
      {/each}
    </div>
    <a class="leave" href="/preview">leave preview ›</a>
  </nav>

  {@render children()}
</div>

<style>
  :global(html, body) {
    background: radial-gradient(ellipse at top, #efe6d0, #e6dcc5 60%, #ddd2b8);
    background-attachment: fixed;
    color: #1f1c17;
    font-family: 'Spectral', 'Iowan Old Style', Georgia, serif;
    font-size: 18px;
    line-height: 1.65;
    -webkit-font-smoothing: antialiased;
  }
  .desk { min-height: 100vh; padding: 2rem 1rem 5rem; }

  .desk-nav {
    max-width: 720px;
    margin: 0 auto 1.5rem;
    display: flex; align-items: center; gap: 1.5rem;
    padding: .65rem 1rem;
    background: rgba(241, 233, 212, .9);
    border-radius: 10px;
    border: 1px solid #d8cbac;
    font-family: 'JetBrains Mono', Menlo, monospace;
    font-size: 11px;
    letter-spacing: .04em;
    color: #6b6358;
  }
  .brand {
    display: flex; align-items: center; gap: .5rem;
    font-family: 'Spectral', serif;
    font-style: italic;
    font-weight: 600;
    font-size: 1.05rem;
    color: #1c3a72;
    text-transform: none;
    letter-spacing: 0;
  }
  .seal {
    width: 22px; height: 22px; border-radius: 999px;
    background: #1c3a72; color: #efe6d0;
    display: grid; place-items: center;
    font-family: 'Spectral', serif; font-style: italic; font-weight: 600;
    font-size: 14px;
  }
  .links { display: flex; gap: .15rem; flex-wrap: wrap; }
  .links a {
    padding: .35rem .65rem;
    color: #6b6358;
    text-decoration: none;
    border-radius: 6px;
  }
  .links a:hover { background: #e6dcc5; color: #1f1c17; }
  .links a.current { background: #1c3a72; color: #efe6d0; }
  .leave { margin-left: auto; color: #aa2828; text-decoration: none; }
  .leave:hover { text-decoration: underline; }

  /* Shared paper styles for child pages */
  :global(.desk .letter) {
    width: 100%;
    max-width: 680px;
    margin: 0 auto;
    padding: 3rem 4rem;
    background:
      repeating-linear-gradient(
        to bottom,
        transparent 0,
        transparent 31px,
        rgba(31, 28, 23, .06) 31px,
        rgba(31, 28, 23, .06) 32px
      ),
      linear-gradient(to right, transparent 0, transparent 70px, rgba(170, 40, 40, .25) 70px, rgba(170, 40, 40, .25) 72px, transparent 72px),
      #f1e9d4;
    box-shadow: 0 1px 1px rgba(0,0,0,.05), 0 8px 30px rgba(60, 50, 40, .12), 0 30px 60px rgba(60, 50, 40, .12);
    position: relative;
  }
  :global(.desk .letter.no-margin) {
    background:
      repeating-linear-gradient(to bottom, transparent 0, transparent 31px, rgba(31, 28, 23, .06) 31px, rgba(31, 28, 23, .06) 32px),
      #f1e9d4;
    padding: 3rem 3rem;
  }
  :global(.desk .letter::before) {
    content: ''; position: absolute; top: 0; right: 0;
    width: 32px; height: 32px;
    background: linear-gradient(225deg, #d8cbac 0 50%, transparent 50%);
    box-shadow: -1px 1px 2px rgba(0,0,0,.08);
  }
  :global(.desk em) { color: #1c3a72; font-style: italic; font-weight: 500; }
  :global(.desk .caveat) {
    font-family: 'Caveat', 'Snell Roundhand', cursive;
    font-weight: 500;
    color: #1c3a72;
  }
</style>
