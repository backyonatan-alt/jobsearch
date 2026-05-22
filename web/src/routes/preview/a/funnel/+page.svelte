<script>
  import { FUNNEL_VIEW } from '$lib/preview-data.js';
</script>

<svelte:head>
  <title>Atelier — Funnel</title>
</svelte:head>

<section>
  <h2 class="section-title">A study of the funnel</h2>
  <p class="lede-text">
    Five applications progressed beyond submission. Three reached the screen,
    two the interview, one the offer. Observations follow.
  </p>

  <table class="stages">
    <thead>
      <tr>
        <th>Stage</th>
        <th>Count</th>
        <th>Conversion</th>
        <th></th>
      </tr>
    </thead>
    <tbody>
      {#each FUNNEL_VIEW.stages as s, i}
        <tr>
          <td class="stage-name">{s.label}</td>
          <td class="num">{s.count}</td>
          <td class="num pct">
            {#if i === 0}—{:else}{s.pct}%{/if}
          </td>
          <td class="bar-cell">
            <div class="bar" style="width: {s.pct}%"></div>
          </td>
        </tr>
      {/each}
    </tbody>
  </table>

  <h3 class="annotations-title">Notes from the editor</h3>
  <ol class="annotations">
    {#each FUNNEL_VIEW.insights as note}
      <li>
        <h4>{note.title}</h4>
        <p>{note.body}</p>
      </li>
    {/each}
  </ol>
</section>

<style>
  table { width: 100%; border-collapse: collapse; margin: 0 0 3rem; }
  th {
    text-align: left;
    font-family: 'Inter', sans-serif;
    font-size: 11px; letter-spacing: .12em; text-transform: uppercase;
    color: #6b6358; font-weight: 500;
    padding: .5rem .75rem;
    border-bottom: 1px solid #d9cfb9;
  }
  td { padding: .85rem .75rem; border-bottom: 1px solid #ece4d0; vertical-align: middle; }
  .stage-name { font-style: italic; font-size: 1.15rem; font-weight: 500; }
  .num { font-family: 'Inter', sans-serif; font-size: 14px; color: #1d1a14; }
  .pct { color: #7a1f2f; font-weight: 600; }
  .bar-cell { width: 50%; }
  .bar {
    height: 4px;
    background: linear-gradient(to right, #7a1f2f, #c97f8a);
    border-radius: 999px;
    transition: width .3s ease;
  }

  .annotations-title {
    font-family: 'Cormorant Garamond', serif;
    font-style: italic; font-weight: 500;
    font-size: 1.55rem;
    margin: 2rem 0 1rem;
    border-top: 1px solid #1d1a14;
    padding-top: 2rem;
  }
  .annotations { list-style: decimal; padding-left: 1.5rem; margin: 0; }
  .annotations li + li { margin-top: 1.25rem; }
  .annotations h4 {
    font-style: italic; font-weight: 500;
    font-size: 1.1rem;
    margin: 0 0 .35rem;
    color: #1d1a14;
  }
  .annotations p { margin: 0; color: #3d3830; }
</style>
