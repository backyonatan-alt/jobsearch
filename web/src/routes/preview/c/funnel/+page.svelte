<script>
  import { FUNNEL_VIEW } from '$lib/preview-data.js';
</script>

<svelte:head>
  <title>Letter — Ledger</title>
</svelte:head>

<article class="letter">
  <header>
    <h1>The ledger</h1>
    <p class="caveat sub">a quiet accounting</p>
  </header>

  <table class="book">
    <thead>
      <tr>
        <th>stage</th>
        <th class="num">count</th>
        <th class="num">carry</th>
      </tr>
    </thead>
    <tbody>
      {#each FUNNEL_VIEW.stages as s, i}
        <tr>
          <td><em>{s.label}</em></td>
          <td class="num">{s.count}</td>
          <td class="num pct">{#if i === 0}—{:else}{s.pct}%{/if}</td>
        </tr>
      {/each}
    </tbody>
  </table>

  <hr class="divider" />

  <section class="margin-notes">
    <h2 class="caveat">margin notes</h2>
    {#each FUNNEL_VIEW.insights as note, i}
      <article class="note">
        <p class="hand caveat">{i + 1}.</p>
        <div>
          <h3>{note.title}</h3>
          <p>{note.body}</p>
        </div>
      </article>
    {/each}
  </section>
</article>

<style>
  header { text-align: center; margin-bottom: 1.5rem; }
  header h1 { font-family: 'Spectral', serif; font-weight: 500; font-style: italic; font-size: 2.5rem; margin: 0; color: #1c3a72; line-height: 1; }
  .sub { font-size: 1.3rem; margin: .35rem 0 0; }

  .book { width: 100%; border-collapse: collapse; font-family: 'JetBrains Mono', monospace; font-size: 13px; }
  .book th { text-align: left; padding: .5rem .75rem; font-weight: 500; color: #6b6358; border-bottom: 2px solid #1f1c17; text-transform: uppercase; letter-spacing: .08em; font-size: 10px; }
  .book td { padding: .55rem .75rem; border-bottom: 1px dotted rgba(31, 28, 23, .2); }
  .book .num { text-align: right; }
  .book em { font-family: 'Spectral', serif; font-style: italic; font-size: 1.05rem; color: #1f1c17; font-weight: 500; }
  .book .pct { color: #1c3a72; font-weight: 500; }

  .divider { border: 0; text-align: center; margin: 2rem 0; }
  .divider::before { content: '✦   ·   ✦'; letter-spacing: .5em; font-size: .9rem; color: #8a7d5e; }

  .margin-notes h2 { font-size: 1.65rem; margin: 0 0 1rem; }
  .note { display: grid; grid-template-columns: 2.5rem 1fr; gap: 1rem; padding: .85rem 0; border-bottom: 1px dashed rgba(31, 28, 23, .12); }
  .note:last-child { border-bottom: 0; }
  .hand { font-size: 1.8rem; line-height: 1; padding-top: .25rem; text-align: right; }
  .note h3 { font-family: 'Spectral', serif; font-weight: 500; font-style: italic; font-size: 1.15rem; margin: 0 0 .35rem; color: #1f1c17; }
  .note p { margin: 0; color: #3d3830; font-size: .98rem; }
</style>
