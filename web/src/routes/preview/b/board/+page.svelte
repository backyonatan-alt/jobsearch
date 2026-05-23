<script>
  import { PREVIEW_APPS, STATUSES, fmtDate } from '$lib/preview-data.js';

  const STATUS_COLOR = {
    wishlist:  '#a39d92',
    applied:   '#3b82f6',
    screen:    '#f59e0b',
    interview: '#c45ba8',
    offer:     '#10b981',
    rejected:  '#ef4444',
    withdrawn: '#71717a'
  };

  const TONES = ['#ff8a5b', '#c45ba8', '#7e6cd6', '#5b8def', '#5bbb8a', '#d6b15b'];
  const monoTone = (s) => TONES[s.charCodeAt(0) % TONES.length];

  const grouped = STATUSES.reduce((acc, s) => {
    acc[s] = PREVIEW_APPS.filter((a) => a.status === s);
    return acc;
  }, {});
</script>

<svelte:head>
  <title>Board — Pursuit</title>
</svelte:head>

<p class="breadcrumb"><a href="/preview/b">Workspace</a> <span>/</span> Board</p>
<div class="page-head">
  <h1>Board</h1>
  <div class="head-actions">
    <button class="btn-ghost">Group: status ▾</button>
    <button class="btn-primary">+ New</button>
  </div>
</div>

<div class="board">
  {#each STATUSES as s}
    <section class="col">
      <header>
        <span class="col-title">
          <span class="col-dot" style="background: {STATUS_COLOR[s]}"></span>
          {s}
        </span>
        <span class="col-count">{grouped[s].length}</span>
      </header>

      <div class="cards">
        {#each grouped[s] as a}
          <a href={a.slug === 'anthropic' ? '/preview/b/anthropic' : '#'} class="card">
            <div class="card-row">
              <span class="mono" style="background: {monoTone(a.company)}">{a.company[0]}</span>
              <p class="co">{a.company}</p>
            </div>
            <p class="role">{a.role}</p>
            <div class="card-foot">
              <span class="when">{fmtDate(a.applied_at)}</span>
              {#if a.cv_variant}
                <span class="cv">{a.cv_variant}</span>
              {/if}
            </div>
          </a>
        {/each}
        <button class="add">+ Add</button>
      </div>
    </section>
  {/each}
</div>

<style>
  .page-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 1.5rem; }
  .head-actions { display: flex; gap: .5rem; }
  .btn-ghost {
    padding: .4rem .7rem; background: #fbf9f4;
    border: 1px solid #ebe6dd; border-radius: 6px;
    font: inherit; font-size: 13px; color: #4a4842; cursor: pointer;
  }
  .btn-ghost:hover { background: #ece5d4; }
  .btn-primary {
    padding: .4rem .9rem;
    background: linear-gradient(135deg, #18181b 0%, #2c2a32 100%);
    color: #fbf9f4; border: 0; border-radius: 6px;
    font: inherit; font-size: 13px; font-weight: 500; cursor: pointer;
  }

  .board {
    display: grid;
    grid-template-columns: repeat(7, minmax(180px, 1fr));
    gap: .65rem;
    overflow-x: auto;
    padding-bottom: 1rem;
  }
  .col {
    background: #fbf9f4;
    border: 1px solid #ebe6dd;
    border-radius: 10px;
    padding: .65rem;
    min-height: 320px;
    display: flex; flex-direction: column;
  }
  .col header {
    display: flex; align-items: center; justify-content: space-between;
    padding: 0 .25rem .65rem;
    margin-bottom: .5rem;
    border-bottom: 1px solid #ebe6dd;
  }
  .col-title {
    display: inline-flex; align-items: center; gap: .4rem;
    font-size: 12px; color: #18181b;
    text-transform: capitalize; font-weight: 500;
  }
  .col-dot { width: 7px; height: 7px; border-radius: 999px; }
  .col-count {
    font-family: 'JetBrains Mono', monospace;
    font-size: 10px; color: #71717a;
    background: #ece5d4;
    padding: 1px 6px; border-radius: 4px;
  }
  .cards { display: flex; flex-direction: column; gap: .5rem; }

  .card {
    background: #fff;
    border: 1px solid #ebe6dd;
    border-radius: 8px;
    padding: .65rem .75rem;
    text-decoration: none;
    color: inherit;
    transition: transform .12s ease, border-color .12s ease, box-shadow .12s ease;
  }
  .card:hover {
    transform: translateY(-2px);
    border-color: #c45ba8;
    box-shadow: 0 4px 10px rgba(196, 91, 168, .12);
  }
  .card-row { display: flex; align-items: center; gap: .5rem; margin-bottom: .35rem; }
  .mono {
    width: 22px; height: 22px; border-radius: 6px;
    color: #fff;
    display: grid; place-items: center;
    font-weight: 600; font-size: 11px;
  }
  .co { margin: 0; font-weight: 500; font-size: 13px; color: #18181b; }
  .role { margin: 0; font-size: 11px; color: #71717a; line-height: 1.35; }
  .card-foot {
    display: flex; gap: .5rem;
    margin-top: .5rem;
    padding-top: .5rem;
    border-top: 1px dashed #ebe6dd;
    font-size: 10px; color: #a39d92;
    font-family: 'JetBrains Mono', monospace;
  }
  .cv { background: #ece5d4; color: #6f685c; padding: 0 4px; border-radius: 3px; }

  .add {
    background: transparent; border: 1px dashed #ebe6dd;
    color: #a39d92; cursor: pointer;
    padding: .55rem; border-radius: 8px;
    font: inherit; font-size: 12px;
    margin-top: .25rem;
  }
  .add:hover { background: #fbf9f4; color: #c45ba8; border-color: #c45ba8; }
</style>
