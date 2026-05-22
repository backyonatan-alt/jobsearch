<script>
  import { PREVIEW_APPS, STATUSES } from '$lib/preview-data.js';

  const STATUS_STYLE = {
    wishlist:  { bg: '#efe7d8', fg: '#6b5b3d', label: 'seedlings' },
    applied:   { bg: '#e2d5f0', fg: '#4a2f6e', label: 'sown' },
    screen:    { bg: '#f5d8d8', fg: '#8a3a3a', label: 'budding' },
    interview: { bg: '#f8c89d', fg: '#6e3b14', label: 'flowering' },
    offer:     { bg: '#c2dbb8', fg: '#2e5b3a', label: 'fruited' },
    rejected:  { bg: '#d6cfc4', fg: '#5a544a', label: 'fallow' },
    withdrawn: { bg: '#e8e3da', fg: '#6b6358', label: 'pruned' }
  };

  const grouped = STATUSES.reduce((acc, s) => {
    acc[s] = PREVIEW_APPS.filter((a) => a.status === s);
    return acc;
  }, {});
</script>

<svelte:head>
  <title>Garden — Board</title>
</svelte:head>

<section class="hero">
  <p class="eyebrow">By bed</p>
  <h1>Your beds, this season</h1>
</section>

<div class="beds">
  {#each STATUSES as s}
    <article class="bed" style="--bed-bg: {STATUS_STYLE[s].bg}; --bed-fg: {STATUS_STYLE[s].fg};">
      <header>
        <span class="badge">{STATUS_STYLE[s].label}</span>
        <span class="count">{grouped[s].length}</span>
      </header>
      <div class="bed-body">
        {#if grouped[s].length === 0}
          <p class="empty">empty bed</p>
        {:else}
          {#each grouped[s] as a}
            <div class="seed">
              <p class="seed-name">{a.company}</p>
              <p class="seed-role">{a.role}</p>
            </div>
          {/each}
        {/if}
      </div>
    </article>
  {/each}
</div>

<style>
  .hero { margin: 1rem 0 2rem; }
  .hero h1 { font-size: 2.25rem; margin: .35rem 0 0; }
  .beds { display: grid; grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); gap: 1rem; }
  .bed {
    background: var(--bed-bg);
    border-radius: 18px;
    padding: 1rem 1rem 1.25rem;
    min-height: 220px;
  }
  .bed header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem; }
  .badge {
    background: rgba(255,255,255,.5); color: var(--bed-fg);
    padding: .2rem .6rem; border-radius: 999px;
    font-size: 11px; letter-spacing: .12em; text-transform: uppercase; font-weight: 700;
  }
  .count { color: var(--bed-fg); font-weight: 700; font-family: 'Fraunces', serif; }
  .bed-body { display: flex; flex-direction: column; gap: .5rem; }
  .seed {
    background: rgba(255,255,255,.6);
    border-radius: 12px;
    padding: .65rem .85rem;
    backdrop-filter: blur(4px);
  }
  .seed-name { font-weight: 600; margin: 0; font-size: 14px; color: #2d2a27; }
  .seed-role { color: #6b6358; font-size: 12px; margin: .15rem 0 0; }
  .empty { color: var(--bed-fg); opacity: .5; font-style: italic; font-size: 13px; margin: 0; text-align: center; padding-top: 1.5rem; }
</style>
