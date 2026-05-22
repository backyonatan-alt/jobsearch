<script>
  import { ANTHROPIC_DETAIL, fmtLongDate } from '$lib/preview-data.js';
  const a = ANTHROPIC_DETAIL;
</script>

<svelte:head>
  <title>Atelier — Anthropic</title>
</svelte:head>

<article class="entry-detail">
  <header class="dossier-head">
    <div class="kicker">Dossier · Entry № 01</div>
    <h2 class="company">{a.company}</h2>
    <p class="deck">{a.role}<span class="loc"> · {a.location}</span></p>
    <p class="next">
      <span class="next-label">Next:</span> {a.next_step}
    </p>
  </header>

  <section class="block">
    <h3>The interviewer</h3>
    <div class="who">
      <div>
        <p class="name">{a.dossier.name}</p>
        <p class="title">{a.dossier.title}</p>
      </div>
    </div>
    <p class="summary">{a.dossier.summary}</p>
  </section>

  <section class="block">
    <h3>Recent traces</h3>
    <ol class="traces">
      {#each a.dossier.recent as t}
        <li>
          <span class="when">{t.date}</span>
          <span class="text">{t.text}</span>
        </li>
      {/each}
    </ol>
  </section>

  <section class="block split">
    <div>
      <h3>Likely style</h3>
      <ul>
        {#each a.dossier.style as s}
          <li>{s}</li>
        {/each}
      </ul>
    </div>
    <div>
      <h3>Watch for</h3>
      <ul class="watch">
        {#each a.dossier.watchfor as w}
          <li>{w}</li>
        {/each}
      </ul>
    </div>
  </section>

  <section class="block">
    <h3>Timeline</h3>
    <ol class="timeline">
      {#each a.timeline as t}
        <li>
          <span class="t-date">{fmtLongDate(t.date)}</span>
          <span class="t-kind">{t.kind}</span>
          <p class="t-text">{t.text}</p>
        </li>
      {/each}
    </ol>
  </section>

  <section class="block">
    <h3>Notes to self</h3>
    <ul class="notes">
      {#each a.notes as n}
        <li>{n}</li>
      {/each}
    </ul>
  </section>

  <p class="back-row">
    <a href="/preview/a">‹ back to the index</a>
  </p>
</article>

<style>
  .dossier-head { margin-bottom: 2.5rem; border-bottom: 1px solid #1d1a14; padding-bottom: 1.5rem; }
  .kicker {
    font-family: 'Inter', sans-serif;
    font-size: 11px; letter-spacing: .14em; text-transform: uppercase;
    color: #a39d92;
    margin-bottom: .5rem;
  }
  .company { font-size: 3rem; font-weight: 500; font-style: italic; margin: 0; line-height: 1; letter-spacing: -.02em; }
  .deck { font-style: italic; color: #3d3830; font-size: 1.2rem; margin: .5rem 0 1rem; }
  .loc { color: #8a8275; font-style: normal; }
  .next {
    font-family: 'Inter', sans-serif; font-size: 13px;
    background: rgba(122, 31, 47, .08);
    padding: .75rem 1rem;
    border-left: 3px solid #7a1f2f;
    margin: 1rem 0 0;
    color: #2c2820;
  }
  .next-label { color: #7a1f2f; font-weight: 600; letter-spacing: .04em; text-transform: uppercase; font-size: 11px; margin-right: .5rem; }

  .block { margin: 0 0 2.5rem; }
  .block h3 {
    font-style: italic; font-weight: 500; font-size: 1.4rem;
    margin: 0 0 1rem;
    border-bottom: 1px solid #d9cfb9;
    padding-bottom: .35rem;
  }

  .who { display: flex; align-items: center; gap: 1rem; margin-bottom: .75rem; }
  .name { font-weight: 600; font-size: 1.15rem; margin: 0; }
  .title { color: #6b6358; font-style: italic; margin: 0; font-size: 1rem; }
  .summary { margin: 0; color: #2c2820; }

  .traces { list-style: none; padding: 0; margin: 0; }
  .traces li {
    display: grid;
    grid-template-columns: 5rem 1fr;
    gap: 1rem;
    padding: .75rem 0;
    border-bottom: 1px dashed #d9cfb9;
  }
  .traces li:last-child { border-bottom: 0; }
  .traces .when { font-family: 'Inter', sans-serif; font-size: 11px; color: #a39d92; letter-spacing: .08em; text-transform: uppercase; padding-top: .25rem; }
  .traces .text { color: #2c2820; }

  .split { display: grid; grid-template-columns: 1fr 1fr; gap: 2rem; }
  @media (max-width: 600px) { .split { grid-template-columns: 1fr; } }
  .split ul { list-style: none; padding: 0; margin: 0; }
  .split li {
    padding: .55rem 0;
    border-bottom: 1px dotted #d9cfb9;
    font-style: italic;
    color: #2c2820;
  }
  .split li:last-child { border-bottom: 0; }
  .split ul.watch li::before { content: '⚠ '; color: #7a1f2f; font-style: normal; margin-right: .25rem; }

  .timeline { list-style: none; padding: 0; margin: 0; }
  .timeline li {
    padding: .85rem 0;
    border-bottom: 1px solid #ece4d0;
  }
  .timeline li:last-child { border-bottom: 0; }
  .t-date { font-family: 'Inter', sans-serif; font-size: 11px; letter-spacing: .12em; text-transform: uppercase; color: #6b6358; }
  .t-kind { font-style: italic; color: #7a1f2f; margin-left: .65rem; font-size: .95rem; }
  .t-text { margin: .35rem 0 0; color: #2c2820; }

  .notes { list-style: none; padding: 0; margin: 0; }
  .notes li {
    padding: .85rem 0;
    border-bottom: 1px dashed #d9cfb9;
    font-style: italic;
    color: #2c2820;
  }
  .notes li:last-child { border-bottom: 0; }
  .notes li::before {
    content: '“ ';
    color: #7a1f2f;
    font-size: 1.6rem;
    font-style: normal;
    line-height: 0;
    vertical-align: -.35rem;
  }

  .back-row {
    margin-top: 3rem;
    text-align: center;
    font-family: 'Inter', sans-serif;
    font-size: 11px; letter-spacing: .12em; text-transform: uppercase;
  }
  .back-row a { color: #7a1f2f; text-decoration: none; }
  .back-row a:hover { text-decoration: underline; }
</style>
