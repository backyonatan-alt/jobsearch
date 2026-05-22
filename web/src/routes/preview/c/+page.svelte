<script>
  import { PREVIEW_APPS, fmtDate } from '$lib/preview-data.js';

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
</svelte:head>

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
        <div class="margin"><span class="num">{String(i + 1).padStart(2, '0')}</span></div>
        <div class="body">
          <div class="line line-head">
            <a class="entry-link" href={a.slug === 'anthropic' ? '/preview/c/anthropic' : '#'}>
              <h2>{a.company}</h2>
            </a>
            <span class="seal" style="--seal: {STATUS_SEAL[a.status].color}">{STATUS_SEAL[a.status].label}</span>
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
    <p class="ps">
      <em>P.S.</em> Open the envelope on Anthropic for the full dossier.
    </p>
  </footer>
</article>

<style>
  .head-meta { display: flex; justify-content: space-between; font-family: 'JetBrains Mono', Menlo, monospace; font-size: 11px; letter-spacing: .04em; color: #6b6358; margin-bottom: 1.5rem; }
  .head h1 { font-family: 'Spectral', serif; font-weight: 500; font-style: italic; font-size: 3.5rem; line-height: 1; margin: 0; color: #1c3a72; letter-spacing: -.02em; }
  .salutation { margin: 1.5rem 0 0; }
  .salutation .caveat { font-size: 1.6rem; }

  .intro { margin: 1.5rem 0; }
  .intro p { margin: 0; }
  .divider { border: 0; text-align: center; margin: 2rem 0; }
  .divider::before { content: '✦   ·   ✦'; letter-spacing: .5em; font-size: .9rem; color: #8a7d5e; }

  .entries { list-style: none; padding: 0; margin: 0; }
  .entry { display: grid; grid-template-columns: 3rem 1fr; gap: 1rem; padding: 1rem 0; }
  .entry + .entry { border-top: 1px dashed rgba(31, 28, 23, .12); }
  .margin { display: flex; justify-content: flex-end; }
  .num { font-family: 'JetBrains Mono', Menlo, monospace; font-size: 11px; color: #b8a888; padding-top: .5rem; }
  .line-head { display: flex; align-items: baseline; justify-content: space-between; gap: 1rem; }
  .entry-link { text-decoration: none; color: inherit; }
  .entry-link:hover h2 { color: #1c3a72; }
  .body h2 { font-family: 'Spectral', serif; font-weight: 500; font-style: italic; font-size: 1.55rem; margin: 0; line-height: 1.1; color: #1f1c17; transition: color .15s ease; }
  .role { margin: .3rem 0; color: #3d3830; font-size: .98rem; }
  .role .loc { color: #6b6358; font-style: italic; }
  .note { margin: .35rem 0 0; font-family: 'JetBrains Mono', Menlo, monospace; font-size: 11px; color: #6b6358; letter-spacing: .02em; }
  .note .key { color: #1c3a72; font-weight: 500; }
  .note code { font-family: inherit; color: #3d3830; }

  .seal {
    flex-shrink: 0;
    font-family: 'Spectral', serif; font-style: italic; font-size: .85rem;
    letter-spacing: .04em; text-transform: lowercase; color: var(--seal);
    border: 1px solid var(--seal); border-radius: 999px;
    padding: .15rem .65rem; white-space: nowrap; opacity: .85;
  }

  .sign { margin-top: 2.5rem; }
  .sign-off { font-size: 1.4rem; margin: 0; }
  .signature { font-size: 3rem; margin: 0; line-height: 1; }
  .ps { font-size: .95rem; color: #3d3830; margin-top: 1.5rem; }
  .ps em { font-style: italic; font-weight: 500; }
</style>
