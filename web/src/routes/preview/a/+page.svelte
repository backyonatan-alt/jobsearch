<script>
  import { PREVIEW_APPS, fmtDate } from '$lib/preview-data.js';

  const STATUS_TONE = {
    wishlist:  '#8a8275',
    applied:   '#1d1a14',
    screen:    '#5a3d1a',
    interview: '#7a1f2f',
    offer:     '#2e5b3a',
    rejected:  '#6b6058',
    withdrawn: '#a39d92'
  };
</script>

<svelte:head>
  <title>Atelier — Applications</title>
</svelte:head>

<section class="lede">
  <p>
    <span class="dropcap">S</span>ix entries this season. <em>One offer outstanding</em> from Linear;
    one interview loop in progress with Anthropic; three further conversations pending review.
    A note on screening cadence appears in the funnel.
  </p>
</section>

<ol class="entries">
  {#each PREVIEW_APPS as a, i}
    <li class="entry">
      <div class="folio">№ {String(i + 1).padStart(2, '0')}</div>
      <div class="body">
        <a class="company-link" href={a.slug === 'anthropic' ? '/preview/a/anthropic' : '#'}>
          <h2 class="company">{a.company}</h2>
        </a>
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

<p class="hint">Tap “Anthropic” to read the dossier — that is where this design comes alive.</p>

<style>
  .lede {
    margin: 0 0 2.5rem;
    font-size: 1.15rem; font-style: italic; color: #2c2820;
    border-bottom: 1px solid #d9cfb9;
    padding-bottom: 2rem;
  }
  .dropcap {
    float: left;
    font-size: 3.5rem; line-height: .85; font-weight: 600;
    margin: .15rem .35rem 0 0;
    color: #7a1f2f; font-style: normal;
  }

  .entries { list-style: none; padding: 0; margin: 0; }
  .entry {
    display: grid; grid-template-columns: 4rem 1fr; gap: 1.25rem;
    padding: 1.5rem 0; border-bottom: 1px solid #d9cfb9;
  }
  .entry:last-child { border-bottom: 0; }
  .folio { font-family: 'Inter', sans-serif; font-size: 11px; letter-spacing: .12em; color: #a39d92; padding-top: .5rem; }
  .company-link { text-decoration: none; color: inherit; }
  .company-link:hover .company { color: #7a1f2f; }
  .company { font-size: 1.8rem; font-weight: 500; margin: 0; line-height: 1.15; letter-spacing: -.01em; transition: color .15s ease; }
  .deck { margin: .35rem 0 .85rem; font-style: italic; color: #3d3830; font-size: 1.05rem; }
  .deck .loc { color: #8a8275; font-style: normal; }
  .meta { display: flex; align-items: center; gap: .65rem; flex-wrap: wrap; font-family: 'Inter', sans-serif; font-size: 11px; letter-spacing: .08em; text-transform: uppercase; }
  .meta .status { font-weight: 600; }
  .meta .when, .meta .cv, .meta .src { color: #6b6358; }
  .rule { width: 14px; height: 1px; background: #c9bfa8; display: inline-block; }

  .hint {
    margin: 3rem 0 0;
    text-align: center;
    font-family: 'Inter', sans-serif;
    font-size: 11px; letter-spacing: .12em; text-transform: uppercase;
    color: #a39d92;
  }
</style>
