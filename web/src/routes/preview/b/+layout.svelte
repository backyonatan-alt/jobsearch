<script>
  import { page } from '$app/state';
  let { children } = $props();

  const sections = [
    {
      title: 'Workspace',
      items: [
        { href: '/preview/b',           label: 'Applications', icon: 'apps' },
        { href: '/preview/b/board',     label: 'Board',        icon: 'board' },
        { href: '/preview/b/funnel',    label: 'Funnel',       icon: 'funnel' }
      ]
    },
    {
      title: 'Recent dossiers',
      items: [
        { href: '/preview/b/anthropic', label: 'Anthropic',    icon: 'dot', tag: 'INT' },
        { href: '#',                    label: 'Stripe',       icon: 'dot', tag: 'APP' },
        { href: '#',                    label: 'Vercel',       icon: 'dot', tag: 'SCR' }
      ]
    }
  ];

  function isCurrent(href) {
    if (href === '/preview/b') return page.url.pathname === '/preview/b';
    return page.url.pathname === href;
  }

  // Tiny inline SVG icon set, sized to match Inter at 14px.
  const ICONS = {
    apps:   'M3 3h7v7H3zM14 3h7v7h-7zM3 14h7v7H3zM14 14h7v7h-7z',
    board:  'M4 5h4v14H4zM10 5h4v9h-4zM16 5h4v6h-4z',
    funnel: 'M3 4h18l-7 9v7l-4-2v-5z',
    dot:    'M12 12m-3 0a3 3 0 1 0 6 0a3 3 0 1 0-6 0'
  };
</script>

<svelte:head>
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous" />
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=Fraunces:opsz,wght@9..144,400;9..144,500;9..144,600&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet" />
</svelte:head>

<div class="shell">
  <aside class="sidebar">
    <div class="ws">
      <div class="ws-mark">
        <span></span>
      </div>
      <div class="ws-info">
        <p class="ws-name">Pursuit</p>
        <p class="ws-sub">back.yonatan@gmail.com</p>
      </div>
      <button class="ws-switch" aria-label="Switch workspace">⌃</button>
    </div>

    <button class="cta">
      + New application
      <span class="kbd">N</span>
    </button>

    {#each sections as section}
      <nav class="group">
        <p class="group-title">{section.title}</p>
        {#each section.items as item}
          <a href={item.href} class="navlink" class:current={isCurrent(item.href)}>
            <svg viewBox="0 0 24 24" width="14" height="14" aria-hidden="true">
              <path d={ICONS[item.icon]} fill="none" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round" stroke-linecap="round"/>
            </svg>
            <span>{item.label}</span>
            {#if item.tag}<span class="tag">{item.tag}</span>{/if}
          </a>
        {/each}
      </nav>
    {/each}

    <div class="spacer"></div>

    <a class="leave" href="/preview">‹ leave preview</a>
  </aside>

  <main class="main">
    <header class="topbar">
      <div class="search">
        <svg viewBox="0 0 24 24" width="14" height="14" aria-hidden="true">
          <circle cx="11" cy="11" r="6" fill="none" stroke="currentColor" stroke-width="1.6"/>
          <path d="M16 16l4 4" stroke="currentColor" stroke-width="1.6" stroke-linecap="round"/>
        </svg>
        <input placeholder="Search applications, companies, notes…" />
        <span class="kbd">⌘ K</span>
      </div>
      <button class="icon-btn" aria-label="Notifications">
        <svg viewBox="0 0 24 24" width="16" height="16"><path d="M6 16V11a6 6 0 1 1 12 0v5l1 2H5l1-2zM9 19a3 3 0 0 0 6 0" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round" stroke-linecap="round"/></svg>
      </button>
      <div class="me">
        <div class="avatar">B</div>
      </div>
    </header>

    <div class="content">
      {@render children()}
    </div>
  </main>
</div>

<style>
  :global(html, body) {
    background: #f7f5f0;
    color: #18181b;
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
    font-size: 14px;
    line-height: 1.45;
    -webkit-font-smoothing: antialiased;
  }

  .shell { display: grid; grid-template-columns: 244px 1fr; min-height: 100vh; }

  /* SIDEBAR */
  .sidebar {
    background: #fbf9f4;
    border-right: 1px solid #ebe6dd;
    padding: 1rem .75rem;
    display: flex; flex-direction: column;
    gap: 1rem;
  }
  .ws {
    display: flex; align-items: center; gap: .65rem;
    padding: .55rem .65rem;
    border-radius: 8px;
    transition: background .15s ease;
  }
  .ws:hover { background: rgba(0,0,0,.03); }
  .ws-mark {
    width: 28px; height: 28px; border-radius: 8px;
    background: linear-gradient(135deg, #ff8a5b 0%, #c45ba8 100%);
    display: grid; place-items: center;
    box-shadow: 0 1px 2px rgba(196, 91, 168, .25);
  }
  .ws-mark span {
    width: 10px; height: 10px; border-radius: 999px;
    background: #fbf9f4;
  }
  .ws-info { flex: 1; min-width: 0; }
  .ws-name {
    margin: 0; font-weight: 600; font-size: 13px;
    color: #18181b;
    white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
  }
  .ws-sub {
    margin: 0; font-size: 11px; color: #8a8378;
    white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
  }
  .ws-switch {
    background: transparent; border: 0; color: #8a8378;
    cursor: pointer; padding: 4px; border-radius: 4px;
  }
  .ws-switch:hover { background: rgba(0,0,0,.05); color: #18181b; }

  .cta {
    display: flex; align-items: center; justify-content: space-between;
    padding: .55rem .85rem;
    background: linear-gradient(135deg, #18181b 0%, #2c2a32 100%);
    color: #fbf9f4;
    border: 0; border-radius: 8px;
    font-size: 13px; font-weight: 500;
    cursor: pointer;
    box-shadow: 0 1px 2px rgba(0,0,0,.06);
  }
  .cta:hover { transform: translateY(-1px); transition: transform .12s ease; }
  .cta .kbd {
    font-family: 'JetBrains Mono', monospace;
    font-size: 10px;
    background: rgba(255,255,255,.12);
    border-radius: 4px;
    padding: 2px 6px;
    color: #d4d2cc;
  }

  .group { display: flex; flex-direction: column; gap: 1px; }
  .group-title {
    margin: 0 0 .25rem .65rem;
    font-size: 10px; font-weight: 600; letter-spacing: .08em; text-transform: uppercase;
    color: #a39d92;
  }
  .navlink {
    display: flex; align-items: center; gap: .55rem;
    padding: .4rem .65rem;
    border-radius: 6px;
    color: #4a4842;
    text-decoration: none;
    font-size: 13px;
    transition: background .12s ease, color .12s ease;
  }
  .navlink:hover { background: rgba(0,0,0,.04); color: #18181b; }
  .navlink.current { background: #ece5d4; color: #18181b; font-weight: 500; }
  .navlink.current svg { color: #c45ba8; }
  .navlink svg { color: #a39d92; flex-shrink: 0; }
  .navlink .tag {
    margin-left: auto;
    font-family: 'JetBrains Mono', monospace;
    font-size: 9px; font-weight: 500;
    padding: 1px 5px; border-radius: 4px;
    color: #8a8378;
    background: rgba(0,0,0,.04);
    letter-spacing: .05em;
  }

  .spacer { flex: 1; }
  .leave {
    color: #a39d92; font-size: 12px; text-decoration: none;
    padding: .5rem .65rem;
  }
  .leave:hover { color: #c45ba8; }

  /* MAIN */
  .main { display: flex; flex-direction: column; min-width: 0; }
  .topbar {
    display: flex; align-items: center; gap: .65rem;
    padding: .65rem 1.5rem;
    border-bottom: 1px solid #ebe6dd;
    background: rgba(247, 245, 240, .85);
    backdrop-filter: blur(8px);
    position: sticky; top: 0; z-index: 10;
  }
  .search {
    flex: 1; max-width: 480px;
    display: flex; align-items: center; gap: .5rem;
    padding: .45rem .75rem;
    background: #fbf9f4;
    border: 1px solid #ebe6dd;
    border-radius: 8px;
    color: #a39d92;
  }
  .search input {
    flex: 1; border: 0; background: transparent;
    font-family: inherit; font-size: 13px;
    color: #18181b; outline: 0;
  }
  .search input::placeholder { color: #a39d92; }
  .search .kbd {
    font-family: 'JetBrains Mono', monospace;
    font-size: 10px;
    background: #ece5d4; color: #6f685c;
    padding: 2px 6px; border-radius: 4px;
    letter-spacing: .02em;
  }
  .icon-btn {
    background: transparent; border: 0;
    width: 32px; height: 32px; border-radius: 6px;
    color: #4a4842;
    display: grid; place-items: center;
    cursor: pointer;
  }
  .icon-btn:hover { background: rgba(0,0,0,.04); }
  .me { margin-left: .25rem; }
  .avatar {
    width: 28px; height: 28px; border-radius: 999px;
    background: linear-gradient(135deg, #ff8a5b 0%, #c45ba8 100%);
    color: #fbf9f4;
    display: grid; place-items: center;
    font-weight: 600; font-size: 12px;
    cursor: pointer;
  }

  .content { padding: 2rem 2.5rem 4rem; max-width: 1200px; width: 100%; }

  /* Shared content helpers */
  :global(.main h1) {
    font-family: 'Inter', sans-serif;
    font-weight: 600;
    font-size: 1.6rem;
    letter-spacing: -.015em;
    color: #18181b;
    margin: 0;
  }
  :global(.main h2) {
    font-family: 'Inter', sans-serif;
    font-weight: 600;
    font-size: 1.1rem;
    letter-spacing: -.005em;
    color: #18181b;
    margin: 0;
  }
  :global(.main .breadcrumb) {
    font-size: 12px; color: #a39d92;
    margin: 0 0 .35rem;
    display: flex; align-items: center; gap: .35rem;
  }
  :global(.main .breadcrumb a) { color: inherit; text-decoration: none; }
  :global(.main .breadcrumb a:hover) { color: #18181b; }
  :global(.main .display) {
    font-family: 'Fraunces', Georgia, serif;
    font-weight: 500;
    letter-spacing: -.015em;
  }

  @media (max-width: 800px) {
    .shell { grid-template-columns: 1fr; }
    .sidebar { display: none; }
  }
</style>
