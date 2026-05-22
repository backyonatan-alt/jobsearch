<script>
  import { FUNNEL_VIEW } from '$lib/preview-data.js';
</script>

<svelte:head>
  <title>Garden — Funnel</title>
</svelte:head>

<section class="hero">
  <p class="eyebrow">Yield report</p>
  <h1>The funnel, gently traced</h1>
  <p class="lede">From sown to fruited — where energy is flowing, where it isn't.</p>
</section>

<div class="stages">
  {#each FUNNEL_VIEW.stages as s, i}
    <article class="stage" style="--w: {s.pct}%;">
      <div class="bar">
        <div class="fill"></div>
        <div class="label">
          <span class="stage-name">{s.label}</span>
          <span class="stage-num">{s.count}</span>
        </div>
      </div>
      {#if i > 0}
        <span class="pct">{s.pct}% carried through</span>
      {:else}
        <span class="pct seed">— starting point</span>
      {/if}
    </article>
  {/each}
</div>

<section class="notes">
  <h2>Reading the garden</h2>
  <div class="cards">
    {#each FUNNEL_VIEW.insights as note, i}
      <article class="note" data-tone={i === 0 ? 'good' : i === 1 ? 'warn' : 'go'}>
        <h3>{note.title}</h3>
        <p>{note.body}</p>
      </article>
    {/each}
  </div>
</section>

<style>
  .hero { margin: 1rem 0 2.5rem; }
  .hero h1 { font-size: 2.5rem; margin: .35rem 0 .75rem; }
  .hero .lede { color: #5a5550; margin: 0; }

  .stages { display: flex; flex-direction: column; gap: 1rem; margin-bottom: 3rem; }
  .stage { display: flex; align-items: center; gap: 1.25rem; }
  .bar {
    flex: 1;
    background: rgba(255,255,255,.5);
    border-radius: 18px;
    height: 56px;
    position: relative;
    overflow: hidden;
  }
  .fill {
    position: absolute; inset: 0;
    width: var(--w);
    background: linear-gradient(90deg, #c2dbb8 0%, #9ec496 80%, #6b3d5f 100%);
    border-radius: 18px;
    transition: width .4s ease;
  }
  .label {
    position: relative; z-index: 1;
    display: flex; align-items: center; justify-content: space-between;
    height: 100%;
    padding: 0 1.5rem;
    color: #2d2a27;
  }
  .stage-name { font-family: 'Fraunces', serif; font-weight: 500; font-size: 1.2rem; }
  .stage-num { font-family: 'Fraunces', serif; font-weight: 600; font-size: 1.4rem; color: #6b3d5f; }
  .pct { font-size: 12px; color: #837e75; min-width: 12rem; }
  .pct.seed { color: #b0a8a0; font-style: italic; }

  .notes h2 { font-size: 1.65rem; margin: 0 0 1rem; }
  .cards { display: grid; gap: 1rem; }
  .note {
    background: rgba(255,255,255,.55);
    border-radius: 16px;
    padding: 1.25rem 1.5rem;
    border-left: 4px solid #c2dbb8;
  }
  .note[data-tone="warn"] { border-left-color: #f8c89d; }
  .note[data-tone="go"] { border-left-color: #6b3d5f; }
  .note h3 { font-size: 1.1rem; margin: 0 0 .35rem; font-weight: 500; }
  .note p { margin: 0; color: #5a5550; line-height: 1.55; }
</style>
