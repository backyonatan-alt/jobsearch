<script>
  import { PREVIEW_APPS, STATUSES } from '$lib/preview-data.js';

  const STATUS_INK = {
    wishlist: '#8a7d5e', applied: '#1c3a72', screen: '#7a5a1f',
    interview: '#9a3a1f', offer: '#2e5b3a', rejected: '#6b6058', withdrawn: '#a39d92'
  };

  // Slight per-card rotation for pinned-paper feel.
  const rotations = [-1.2, 0.8, -0.6, 1.4, -1.1, 0.5, -0.7, 1.2];
  const grouped = STATUSES.reduce((acc, s) => { acc[s] = PREVIEW_APPS.filter(a => a.status === s); return acc; }, {});
</script>

<svelte:head>
  <title>Letter — Pinboard</title>
</svelte:head>

<article class="letter no-margin">
  <header>
    <h1>The pinboard</h1>
    <p class="caveat sub">cards on a wall, by column</p>
  </header>

  <div class="cols">
    {#each STATUSES as s}
      <section class="col">
        <h2><span class="ink" style="color: {STATUS_INK[s]}">●</span> {s}<span class="count">{grouped[s].length}</span></h2>
        {#each grouped[s] as a, i}
          <div class="pin" style="transform: rotate({rotations[(a.id + i) % rotations.length]}deg);">
            <span class="thumb"></span>
            <p class="pin-co">{a.company}</p>
            <p class="pin-role">{a.role}</p>
          </div>
        {/each}
        {#if grouped[s].length === 0}
          <p class="empty caveat">empty</p>
        {/if}
      </section>
    {/each}
  </div>
</article>

<style>
  header { text-align: center; margin-bottom: 2rem; }
  header h1 { font-family: 'Spectral', serif; font-weight: 500; font-style: italic; font-size: 2.5rem; margin: 0; color: #1c3a72; line-height: 1; }
  .sub { font-size: 1.3rem; margin: .35rem 0 0; }

  .cols { display: grid; grid-template-columns: repeat(2, 1fr); gap: 1.5rem; }
  @media (max-width: 600px) { .cols { grid-template-columns: 1fr; } }

  .col h2 {
    font-family: 'Spectral', serif; font-weight: 500; font-style: italic;
    font-size: 1.2rem;
    margin: 0 0 1rem;
    text-transform: lowercase;
    border-bottom: 1px dashed rgba(31, 28, 23, .15);
    padding-bottom: .35rem;
    display: flex; align-items: center; gap: .5rem;
  }
  .ink { font-style: normal; }
  .count { margin-left: auto; font-family: 'JetBrains Mono', monospace; font-size: 11px; color: #8a7d5e; }

  .pin {
    background: #fdfaf1;
    padding: .85rem 1rem;
    margin-bottom: .85rem;
    box-shadow: 0 1px 2px rgba(0,0,0,.04), 0 4px 10px rgba(60, 50, 40, .08);
    position: relative;
    border: 1px solid rgba(0,0,0,.04);
  }
  .thumb {
    position: absolute;
    top: -6px; left: 50%;
    transform: translateX(-50%);
    width: 12px; height: 12px;
    background: radial-gradient(circle at 30% 30%, #d04646, #7a1f1f);
    border-radius: 999px;
    box-shadow: 0 1px 2px rgba(0,0,0,.25);
  }
  .pin-co { font-family: 'Spectral', serif; font-weight: 500; font-style: italic; font-size: 1.05rem; margin: 0; color: #1f1c17; }
  .pin-role { font-family: 'JetBrains Mono', monospace; font-size: 11px; color: #6b6358; margin: .35rem 0 0; }

  .empty { font-size: 1.1rem; color: #b8a888; text-align: center; padding: 1rem 0; margin: 0; }
</style>
