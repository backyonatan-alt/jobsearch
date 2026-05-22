<script>
  // Letter — handwritten correspondence. Ruled cream paper, fountain-pen blue,
  // wax-seal status stamps. Serif body, monospace for meta. Quiet & intimate.
  import { PREVIEW_APPS, FUNNEL, fmtDate } from '$lib/preview-data.js';

  const STATUS_SEAL = {
    wishlist:  { color: '#8a7d5e', label: 'noted' },
    applied:   { color: '#1c3a72', label: 'sent' },
    screen:    { color: '#7a5a1f', label: 'reply' },
    interview: { color: '#9a3a1f', label: 'meeting' },
    offer:     { color: '#2e5b3a', label: 'offered' },
    rejected:  { color: '#6b6058', label: 'closed' },
    withdrawn: { color: '#a39d92', label: 'withdrawn' }
  };

  const today = new Date().toLocaleString('en-US', { weekday: 'long', month: 'long', day: 'numeric', year: 'numeric' });
</script>

<svelte:head>
  <title>Letter — Pursuit</title>
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous" />
  <link href="https://fonts.googleapis.com/css2?family=Spectral:ital,wght@0,400;0,500;0,600;1,400&family=JetBrains+Mono:wght@400;500&family=Caveat:wght@500&display=swap" rel="stylesheet" />
</svelte:head>

<div class="desk">
  <article class="letter">
    <header class="head">
      <div class="head-meta">
        <span class="date">{today}</span>
        <span class="addr">back.yonatan@gmail.com</span>
      </div>
      <h1>Pursuit</h1>
      <p class="salutation"><span class="caveat">Dear reader,</span></p>
    </header>

    <section class="intro">
      <p>
        A short letter on this week's pursuits. Six in flight; one offer arrived
        from <em>Linear</em> and awaits reply; the conversation with
        <em>Anthropic</em> deepens. The remainder are quiet — a gentle nudge to
        Stripe on Friday should do.
      </p>
    </section>

    <hr class="divider" />

    <ol class="entries">
      {#each PREVIEW_APPS as a, i}
        <li class="entry">
          <div class="margin">
            <span class="num">{String(i + 1).padStart(2, '0')}</span>
          </div>
          <div class="body">
            <div class="line line-head">
              <h2>{a.company}</h2>
              <span class="seal" style="--seal: {STATUS_SEAL[a.status].color}">
                {STATUS_SEAL[a.status].label}
              </span>
            </div>
            <p class="role">{a.role} — <span class="loc">{a.location}</span></p>
            <p class="note">
              <span class="key">re:</span> applied {fmtDate(a.applied_at)}
              {#if a.cv_variant} · cv <code>{a.cv_variant}</code>{/if}
              · via <em>{a.source}</em>
            </p>
          </div>
        </li>
      {/each}
    </ol>

    <hr class="divider" />

    <footer class="sign">
      <p class="caveat sign-off">Yours, in pursuit —</p>
      <p class="signature caveat">Y.</p>
      <p class="back-row">
        <a href="/preview" class="back">‹ back to directions</a>
      </p>
    </footer>
  </article>
</div>

<style>
  :global(html, body) {
    background:
      radial-gradient(ellipse at top, #efe6d0, #e6dcc5 60%, #ddd2b8);
    background-attachment: fixed;
    color: #1f1c17;
    font-family: 'Spectral', 'Iowan Old Style', Georgia, serif;
    font-size: 18px;
    line-height: 1.65;
    -webkit-font-smoothing: antialiased;
  }

  .desk { min-height: 100vh; padding: 3rem 1rem 5rem; display: flex; justify-content: center; }

  .letter {
    width: 100%;
    max-width: 680px;
    padding: 3.5rem 4rem;
    background:
      /* horizontal ruled lines */
      repeating-linear-gradient(
        to bottom,
        transparent 0,
        transparent 31px,
        rgba(31, 28, 23, .06) 31px,
        rgba(31, 28, 23, .06) 32px
      ),
      /* vertical margin line */
      linear-gradient(to right, transparent 0, transparent 70px, rgba(170, 40, 40, .25) 70px, rgba(170, 40, 40, .25) 72px, transparent 72px),
      /* paper */
      #f1e9d4;
    box-shadow:
      0 1px 1px rgba(0,0,0,.05),
      0 8px 30px rgba(60, 50, 40, .12),
      0 30px 60px rgba(60, 50, 40, .12);
    position: relative;
  }
  /* tiny corner crease */
  .letter::before {
    content: '';
    position: absolute;
    top: 0; right: 0;
    width: 32px; height: 32px;
    background: linear-gradient(225deg, #d8cbac 0 50%, transparent 50%);
    box-shadow: -1px 1px 2px rgba(0,0,0,.08);
  }

  .head { padding-left: 0; }
  .head-meta {
    display: flex; justify-content: space-between;
    font-family: 'JetBrains Mono', Menlo, monospace;
    font-size: 11px;
    letter-spacing: .04em;
    color: #6b6358;
    margin-bottom: 1.5rem;
  }
  .head h1 {
    font-family: 'Spectral', serif;
    font-weight: 500;
    font-style: italic;
    font-size: 3.5rem;
    line-height: 1;
    margin: 0;
    color: #1c3a72;
    letter-spacing: -.02em;
  }
  .salutation { margin: 1.5rem 0 0; }
  .caveat {
    font-family: 'Caveat', 'Snell Roundhand', cursive;
    font-weight: 500;
    color: #1c3a72;
  }
  .salutation .caveat { font-size: 1.6rem; }

  .intro { margin: 1.5rem 0; }
  .intro p { margin: 0; }
  em { color: #1c3a72; font-style: italic; font-weight: 500; }

  .divider {
    border: 0;
    text-align: center;
    margin: 2rem 0;
    color: #8a7d5e;
  }
  .divider::before {
    content: '✦   ·   ✦';
    letter-spacing: .5em;
    font-size: .9rem;
    color: #8a7d5e;
  }

  .entries { list-style: none; padding: 0; margin: 0; }
  .entry {
    display: grid;
    grid-template-columns: 3rem 1fr;
    gap: 1rem;
    padding: 1rem 0;
  }
  .entry + .entry { border-top: 1px dashed rgba(31, 28, 23, .12); }
  .margin { display: flex; justify-content: flex-end; }
  .num {
    font-family: 'JetBrains Mono', Menlo, monospace;
    font-size: 11px;
    color: #b8a888;
    padding-top: .5rem;
  }
  .line-head {
    display: flex; align-items: baseline; justify-content: space-between;
    gap: 1rem;
  }
  .body h2 {
    font-family: 'Spectral', serif;
    font-weight: 500;
    font-style: italic;
    font-size: 1.55rem;
    margin: 0;
    line-height: 1.1;
    color: #1f1c17;
  }
  .role { margin: .3rem 0; color: #3d3830; font-size: .98rem; }
  .role .loc { color: #6b6358; font-style: italic; }
  .note {
    margin: .35rem 0 0;
    font-family: 'JetBrains Mono', Menlo, monospace;
    font-size: 11px;
    color: #6b6358;
    letter-spacing: .02em;
  }
  .note .key { color: #1c3a72; font-weight: 500; }
  .note code { font-family: inherit; color: #3d3830; }

  .seal {
    flex-shrink: 0;
    font-family: 'Spectral', serif;
    font-style: italic;
    font-size: .85rem;
    letter-spacing: .04em;
    text-transform: lowercase;
    color: var(--seal);
    border: 1px solid var(--seal);
    border-radius: 999px;
    padding: .15rem .65rem;
    white-space: nowrap;
    opacity: .85;
  }

  .sign { margin-top: 2.5rem; }
  .sign-off { font-size: 1.4rem; margin: 0; }
  .signature { font-size: 3rem; margin: 0; line-height: 1; }
  .back-row {
    font-family: 'JetBrains Mono', monospace;
    font-size: 11px;
    letter-spacing: .04em;
    color: #6b6358;
    margin-top: 2rem;
  }
  .back { color: #1c3a72; text-decoration: none; }
  .back:hover { text-decoration: underline; }

  @media (max-width: 600px) {
    .letter { padding: 2.5rem 2rem 2.5rem 3rem; }
    .letter { background-position: 0 0; }
    .entry { grid-template-columns: 2rem 1fr; gap: .65rem; }
    .head h1 { font-size: 2.5rem; }
  }
</style>
