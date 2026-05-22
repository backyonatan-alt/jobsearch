<script>
  import { PREVIEW_APPS, STATUSES } from '$lib/preview-data.js';

  const STATUS_LABELS = {
    wishlist: 'I.', applied: 'II.', screen: 'III.',
    interview: 'IV.', offer: 'V.', rejected: 'VI.', withdrawn: 'VII.'
  };

  const grouped = STATUSES.reduce((acc, s) => {
    acc[s] = PREVIEW_APPS.filter((a) => a.status === s);
    return acc;
  }, {});
</script>

<svelte:head>
  <title>Atelier — Board</title>
</svelte:head>

<section>
  <h2 class="section-title">The board, in seven movements</h2>
  <p class="lede-text">An arrangement by stage, from <em>wishlist</em> through <em>withdrawn</em> — each application read as a movement in a longer composition.</p>

  <div class="movements">
    {#each STATUSES as s}
      <article class="movement">
        <header>
          <span class="numeral">{STATUS_LABELS[s]}</span>
          <h3>{s}</h3>
          <span class="count">{grouped[s].length}</span>
        </header>
        {#if grouped[s].length === 0}
          <p class="silent">— silent —</p>
        {:else}
          <ul>
            {#each grouped[s] as a}
              <li>
                <span class="company">{a.company}</span>
                <span class="role">— {a.role}</span>
              </li>
            {/each}
          </ul>
        {/if}
      </article>
    {/each}
  </div>
</section>

<style>
  .movements { display: grid; gap: 1.5rem; }
  .movement {
    border-top: 1px solid #d9cfb9;
    padding-top: 1rem;
  }
  .movement header {
    display: flex; align-items: baseline; gap: .75rem;
    margin-bottom: .65rem;
  }
  .numeral {
    font-family: 'Inter', sans-serif;
    font-size: 11px; letter-spacing: .12em;
    color: #a39d92;
    width: 2rem;
  }
  .movement h3 {
    font-style: italic; font-weight: 500;
    font-size: 1.4rem;
    margin: 0;
    text-transform: lowercase;
    letter-spacing: -.01em;
  }
  .count {
    margin-left: auto;
    font-family: 'Inter', sans-serif;
    font-size: 11px; letter-spacing: .12em; text-transform: uppercase;
    color: #6b6358;
  }
  ul { list-style: none; padding: 0 0 0 2.75rem; margin: 0; }
  li { padding: .3rem 0; }
  .company { font-weight: 500; font-size: 1.1rem; }
  .role { color: #6b6358; font-style: italic; font-size: 1rem; }
  .silent {
    padding-left: 2.75rem; margin: 0;
    color: #a39d92; font-style: italic; font-size: .95rem;
  }
</style>
