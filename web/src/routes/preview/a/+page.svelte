<script>
  // Atelier — editorial / magazine. Serif body, oxblood ink, cream paper.
  // Layout is asymmetric: each application is a numbered entry with the
  // company set in display serif, role as deck, status as small caps.
  import { PREVIEW_APPS, FUNNEL, fmtDate } from '$lib/preview-data.js';

  const STATUS_TONE = {
    wishlist:  '#8a8275',
    applied:   '#1d1a14',
    screen:    '#5a3d1a',
    interview: '#7a1f2f',
    offer:     '#2e5b3a',
    rejected:  '#6b6058',
    withdrawn: '#a39d92'
  };

  const issue = String(new Date().getFullYear() - 2000).padStart(2, '0');
  const month = new Date().toLocaleString('en-US', { month: 'long' });
</script>

<svelte:head>
  <title>Atelier — Pursuit</title>
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous" />
  <link href="https://fonts.googleapis.com/css2?family=Cormorant+Garamond:ital,wght@0,400;0,500;0,600;1,400&family=Inter:wght@400;500;600&display=swap" rel="stylesheet" />
</svelte:head>

<div class="paper">
  <header class="masthead">
    <div class="row">
      <span class="vol">Vol. {issue} · {month}</span>
      <span class="me">back.yonatan@gmail.com</span>
    </div>
    <h1>Pursuit</h1>
    <p class="subhead">A quarterly index of applications, conversations, and outcomes</p>
    <nav>
      <a href="#" class="navlink current">Applications</a>
      <span class="sep">·</span>
      <a href="#" class="navlink">Funnel</a>
      <span class="sep">·</span>
      <a href="#" class="navlink">Dossiers</a>
      <span class="sep">·</span>
      <a href="#" class="navlink">Archive</a>
    </nav>
  </header>

  <section class="lede">
    <p>
      <span class="dropcap">S</span>ix entries this season. <em>One offer outstanding</em>; one
      interview loop in progress; three further conversations pending review.
      A note on the screening process appears overleaf.
    </p>
  </section>

  <ol class="entries">
    {#each PREVIEW_APPS as a, i}
      <li class="entry">
        <div class="folio">№ {String(i + 1).padStart(2, '0')}</div>
        <div class="body">
          <h2 class="company">{a.company}</h2>
          <p class="deck">{a.role} <span class="loc">· {a.location}</span></p>
          <div class="meta">
            <span class="status" style="color: {STATUS_TONE[a.status]}">{a.status}</span>
            <span class="rule"></span>
            <span class="when">{fmtDate(a.applied_at)}</span>
            {#if a.cv_variant}
              <span class="rule"></span>
              <span class="cv">CV / {a.cv_variant}</span>
            {/if}
            <span class="rule"></span>
            <span class="src">via {a.source}</span>
          </div>
        </div>
      </li>
    {/each}
  </ol>

  <footer class="colophon">
    <p>
      Set in <span class="oldstyle">Cormorant Garamond</span> and Inter. Composed in private.
      <a href="/preview">‹ return to directions</a>
    </p>
  </footer>
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

  .masthead { text-align: center; border-bottom: 1px solid #1d1a14; padding-bottom: 1.5rem; }
  .row {
    display: flex; justify-content: space-between;
    font-family: 'Inter', sans-serif; font-size: 11px;
    letter-spacing: .12em; text-transform: uppercase;
    color: #6b6358;
    margin-bottom: 2rem;
  }
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
  nav { display: flex; justify-content: center; align-items: center; gap: .75rem; font-family: 'Inter', sans-serif; font-size: 12px; letter-spacing: .1em; text-transform: uppercase; }
  .navlink { color: #6b6358; text-decoration: none; }
  .navlink.current { color: #7a1f2f; font-weight: 600; }
  .navlink:hover { color: #1d1a14; }
  .sep { color: #cdc4b3; }

  .lede {
    margin: 2.5rem 0 3rem;
    font-size: 1.15rem;
    font-style: italic;
    color: #2c2820;
    border-bottom: 1px solid #d9cfb9;
    padding-bottom: 2rem;
  }
  .lede em { color: #7a1f2f; font-style: italic; }
  .dropcap {
    float: left;
    font-size: 3.5rem;
    line-height: .85;
    font-weight: 600;
    margin: .15rem .35rem 0 0;
    color: #7a1f2f;
    font-style: normal;
  }

  .entries { list-style: none; padding: 0; margin: 0; }
  .entry {
    display: grid;
    grid-template-columns: 4rem 1fr;
    gap: 1.25rem;
    padding: 1.5rem 0;
    border-bottom: 1px solid #d9cfb9;
  }
  .entry:last-child { border-bottom: 0; }
  .folio {
    font-family: 'Inter', sans-serif;
    font-size: 11px;
    letter-spacing: .12em;
    color: #a39d92;
    padding-top: .5rem;
  }
  .company {
    font-size: 1.8rem;
    font-weight: 500;
    margin: 0;
    line-height: 1.15;
    letter-spacing: -.01em;
  }
  .deck {
    margin: .35rem 0 .85rem;
    font-style: italic;
    color: #3d3830;
    font-size: 1.05rem;
  }
  .deck .loc { color: #8a8275; font-style: normal; }
  .meta {
    display: flex; align-items: center; gap: .65rem; flex-wrap: wrap;
    font-family: 'Inter', sans-serif;
    font-size: 11px;
    letter-spacing: .08em;
    text-transform: uppercase;
  }
  .meta .status { font-weight: 600; }
  .meta .when, .meta .cv, .meta .src { color: #6b6358; }
  .rule { width: 14px; height: 1px; background: #c9bfa8; display: inline-block; }

  .colophon {
    margin-top: 4rem;
    padding-top: 1.5rem;
    border-top: 1px solid #1d1a14;
    text-align: center;
    font-family: 'Inter', sans-serif;
    font-size: 11px;
    letter-spacing: .12em;
    text-transform: uppercase;
    color: #6b6358;
  }
  .colophon a { color: #7a1f2f; text-decoration: none; margin-left: .5rem; }
  .colophon a:hover { text-decoration: underline; }
  .oldstyle { font-family: 'Cormorant Garamond', serif; font-style: italic; text-transform: none; letter-spacing: 0; font-size: 13px; }
</style>
