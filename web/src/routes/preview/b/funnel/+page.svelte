<script>
  import { FUNNEL_VIEW } from '$lib/preview-data.js';

  const stages = FUNNEL_VIEW.stages;
  const drops = stages.slice(1).map((s, i) => ({
    from: stages[i].label,
    to:   s.label,
    rate: ((s.count / stages[i].count) * 100).toFixed(0)
  }));
</script>

<svelte:head>
  <title>Funnel — Pursuit</title>
</svelte:head>

<p class="breadcrumb"><a href="/preview/b">Workspace</a> <span>/</span> Funnel</p>
<div class="page-head">
  <h1>Funnel</h1>
  <div class="head-actions">
    <button class="btn-ghost">Last 30 days ▾</button>
    <button class="btn-ghost">Export</button>
  </div>
</div>

<div class="kpis">
  <div class="kpi">
    <p class="k-label">Applied → Offer</p>
    <p class="k-val display">20%</p>
    <p class="k-sub">+5% vs. last month</p>
  </div>
  <div class="kpi">
    <p class="k-label">Time to first reply</p>
    <p class="k-val display">3.2d</p>
    <p class="k-sub">median</p>
  </div>
  <div class="kpi">
    <p class="k-label">Active loops</p>
    <p class="k-val display">2</p>
    <p class="k-sub">interview · offer</p>
  </div>
  <div class="kpi accent">
    <p class="k-label">Best CV variant</p>
    <p class="k-val display">v3‑ai‑focus</p>
    <p class="k-sub">2× the reply rate</p>
  </div>
</div>

<section class="funnel-card">
  <h2>Stage conversion</h2>
  <div class="bars">
    {#each stages as s, i}
      <div class="bar-row">
        <div class="bar-label">
          <span class="b-name">{s.label}</span>
          <span class="b-count">{s.count}</span>
        </div>
        <div class="bar-track">
          <div class="bar-fill" style="width: {s.pct}%"></div>
        </div>
        <div class="bar-pct">
          {#if i === 0}
            <span class="muted">baseline</span>
          {:else}
            <span class="pct">{s.pct}%</span>
            <span class="drop">↓ {(100 - s.pct)}%</span>
          {/if}
        </div>
      </div>
    {/each}
  </div>

  <div class="drops">
    {#each drops as d}
      <div class="drop-card">
        <p class="d-from">{d.from} <span class="arrow">→</span> {d.to}</p>
        <p class="d-rate">{d.rate}% advanced</p>
      </div>
    {/each}
  </div>
</section>

<section class="insights">
  <h2>What we're noticing</h2>
  <div class="cards">
    {#each FUNNEL_VIEW.insights as note, i}
      <article class="insight" data-tone={i === 0 ? 'good' : i === 1 ? 'warn' : 'go'}>
        <span class="tone-dot"></span>
        <div>
          <h3>{note.title}</h3>
          <p>{note.body}</p>
        </div>
      </article>
    {/each}
  </div>
</section>

<style>
  .page-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 1.5rem; }
  .head-actions { display: flex; gap: .5rem; }
  .btn-ghost {
    padding: .4rem .7rem; background: #fbf9f4;
    border: 1px solid #ebe6dd; border-radius: 6px;
    font: inherit; font-size: 13px; color: #4a4842; cursor: pointer;
  }

  .kpis { display: grid; grid-template-columns: repeat(4, 1fr); gap: .75rem; margin-bottom: 2rem; }
  @media (max-width: 800px) { .kpis { grid-template-columns: repeat(2, 1fr); } }
  .kpi {
    background: #fbf9f4;
    border: 1px solid #ebe6dd;
    border-radius: 10px;
    padding: 1rem 1.1rem;
  }
  .kpi.accent {
    background: linear-gradient(135deg, rgba(255,138,91,.12), rgba(196,91,168,.12));
    border-color: rgba(196,91,168,.3);
  }
  .k-label { margin: 0; font-size: 11px; letter-spacing: .04em; color: #71717a; font-weight: 500; text-transform: uppercase; }
  .k-val { margin: .35rem 0 .25rem; font-size: 1.65rem; color: #18181b; }
  .k-sub { margin: 0; font-size: 11px; color: #a39d92; }

  .funnel-card, .insights {
    background: #fbf9f4;
    border: 1px solid #ebe6dd;
    border-radius: 12px;
    padding: 1.5rem 1.75rem;
    margin-bottom: 1rem;
  }
  .funnel-card h2, .insights h2 { margin: 0 0 1.25rem; }

  .bars { display: flex; flex-direction: column; gap: .75rem; }
  .bar-row {
    display: grid;
    grid-template-columns: 120px 1fr 140px;
    align-items: center;
    gap: 1rem;
  }
  .bar-label { display: flex; align-items: baseline; gap: .5rem; }
  .b-name { font-size: 13px; font-weight: 500; color: #18181b; }
  .b-count {
    font-family: 'JetBrains Mono', monospace;
    font-size: 11px;
    color: #71717a;
  }
  .bar-track {
    background: #f1ede3;
    height: 26px;
    border-radius: 6px;
    overflow: hidden;
  }
  .bar-fill {
    height: 100%;
    background: linear-gradient(90deg, #ff8a5b 0%, #c45ba8 100%);
    border-radius: 6px;
    transition: width .3s ease;
  }
  .bar-pct {
    display: flex; align-items: baseline; gap: .5rem;
    font-family: 'JetBrains Mono', monospace;
    font-size: 12px;
  }
  .bar-pct .pct { font-weight: 600; color: #18181b; }
  .bar-pct .drop { color: #ef4444; font-size: 10px; }
  .bar-pct .muted { color: #a39d92; font-size: 11px; }

  .drops {
    display: grid; grid-template-columns: repeat(3, 1fr);
    gap: .65rem;
    margin-top: 1.5rem;
    padding-top: 1.5rem;
    border-top: 1px solid #ebe6dd;
  }
  @media (max-width: 800px) { .drops { grid-template-columns: 1fr; } }
  .drop-card {
    background: #fff;
    border: 1px solid #ebe6dd;
    border-radius: 8px;
    padding: .65rem .85rem;
  }
  .d-from { margin: 0; font-size: 11px; color: #71717a; }
  .d-from .arrow { color: #c45ba8; margin: 0 .25rem; }
  .d-rate { margin: .25rem 0 0; font-weight: 600; color: #18181b; font-size: 14px; }

  .insights .cards { display: flex; flex-direction: column; gap: .75rem; }
  .insight {
    display: flex; gap: .85rem;
    background: #fff;
    border: 1px solid #ebe6dd;
    border-left-width: 3px;
    border-radius: 8px;
    padding: 1rem 1.25rem;
  }
  .insight[data-tone="good"] { border-left-color: #10b981; }
  .insight[data-tone="warn"] { border-left-color: #f59e0b; }
  .insight[data-tone="go"]   { border-left-color: #c45ba8; }
  .tone-dot {
    width: 8px; height: 8px; border-radius: 999px;
    margin-top: 7px; flex-shrink: 0;
  }
  .insight[data-tone="good"] .tone-dot { background: #10b981; }
  .insight[data-tone="warn"] .tone-dot { background: #f59e0b; }
  .insight[data-tone="go"]   .tone-dot { background: #c45ba8; }
  .insight h3 { margin: 0 0 .35rem; font-size: 14px; font-weight: 600; }
  .insight p { margin: 0; color: #4a4842; line-height: 1.55; font-size: 13px; }
</style>
