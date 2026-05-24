<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api.js';
  import { STATUS_LABEL, STATUSES, toDisplayApp, fmtShortDate } from '$lib/app-helpers.js';

  let apps = $state([]);
  let loading = $state(true);
  let dragOver = $state(null); // status currently being hovered during a drag

  // For status moves while dragging
  let dragging = $state(null); // app id being dragged

  onMount(refresh);

  async function refresh() {
    try {
      apps = (await api('/api/applications')).map(toDisplayApp);
    } catch (e) {
      if (e.message !== 'unauthorized') console.error(e);
    } finally {
      loading = false;
    }
  }

  const byStatus = $derived.by(() => {
    const g = Object.fromEntries(STATUSES.map(s => [s, []]));
    for (const a of apps) (g[a.status] ??= []).push(a);
    return g;
  });

  function open(id) {
    goto(`/app/${id}`);
  }

  function onDragStart(e, id) {
    dragging = id;
    e.dataTransfer.effectAllowed = 'move';
    e.dataTransfer.setData('text/plain', String(id));
  }

  function onDragEnd() {
    dragging = null;
    dragOver = null;
  }

  function onDragOver(e, status) {
    e.preventDefault();
    dragOver = status;
  }

  function onDragLeave(status) {
    if (dragOver === status) dragOver = null;
  }

  async function onDrop(e, newStatus) {
    e.preventDefault();
    const id = Number(e.dataTransfer.getData('text/plain'));
    dragOver = null;
    dragging = null;
    if (!id) return;
    const app = apps.find(a => a.id === id);
    if (!app || app.status === newStatus) return;
    // Optimistic update
    apps = apps.map(a => a.id === id ? { ...a, status: newStatus, raw: { ...a.raw, status: newStatus } } : a);
    try {
      await api(`/api/applications/${id}`, {
        method: 'PATCH',
        body: JSON.stringify({ status: newStatus })
      });
    } catch (err) {
      console.error('status update failed', err);
      await refresh(); // roll back to server state
    }
  }
</script>

<svelte:head>
  <title>Board — Pursuit</title>
</svelte:head>

<div class="topbar">
  <div class="crumb"><span class="here">Board</span></div>
  <div class="right">
    <div class="search">
      <svg class="ico" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5">
        <circle cx="7" cy="7" r="4.5"/><path d="M10.5 10.5L13 13"/>
      </svg>
      Search applications, people…
      <span class="kbd">⌘K</span>
    </div>
    <a class="btn" href="/app">List view</a>
  </div>
</div>

<div class="body">
  <div class="board-inner">
    {#if loading}
      <p style="color:var(--mute); padding: 1rem 0;">Loading…</p>
    {:else if apps.length === 0}
      <div class="empty-tab">
        <h3>No applications yet</h3>
        <p>Add your first one from the Today page (⌘N) — they'll appear here grouped by status.</p>
      </div>
    {:else}
      <div class="cols">
        {#each STATUSES as s (s)}
          <section
            class="col {dragOver === s ? 'drag-over' : ''}"
            ondragover={(e) => onDragOver(e, s)}
            ondragleave={() => onDragLeave(s)}
            ondrop={(e) => onDrop(e, s)}
          >
            <header>
              <span class={`pill ${s}`}><span class="pdot"></span>{STATUS_LABEL[s]}</span>
              <span class="ccount">{byStatus[s].length}</span>
            </header>
            <div class="cards">
              {#each byStatus[s] as a (a.id)}
                <button
                  type="button"
                  class="card"
                  class:dragging={dragging === a.id}
                  draggable="true"
                  ondragstart={(e) => onDragStart(e, a.id)}
                  ondragend={onDragEnd}
                  onclick={() => open(a.id)}
                >
                  <div class="card-head">
                    <span class={`logo ${a.logoCls}`}>{a.coShort}</span>
                    <span class="co">{a.co}</span>
                  </div>
                  <div class="role" title={a.role}>{a.role}</div>
                  <div class="meta">
                    <span class="applied">{a.applied}</span>
                    {#if a.cv && a.cv !== '—'}
                      <span class="cv">{a.cv}</span>
                    {/if}
                  </div>
                </button>
              {/each}
              {#if byStatus[s].length === 0}
                <p class="empty-col">drop here</p>
              {/if}
            </div>
          </section>
        {/each}
      </div>
      <p class="hint">Drag a card between columns to change its status, or click to open the dossier.</p>
    {/if}
  </div>
</div>

<style>
  .board-inner {
    max-width: 100%;
    margin: 0;
  }
  .cols {
    display: grid;
    grid-template-columns: repeat(7, minmax(180px, 1fr));
    gap: 12px;
    align-items: start;
  }
  .col {
    background: var(--surface);
    border: 1px solid var(--rule);
    border-radius: 10px;
    padding: 10px;
    min-height: 320px;
    display: flex;
    flex-direction: column;
    gap: 8px;
    transition: background 80ms ease, border-color 80ms ease;
  }
  .col.drag-over {
    background: var(--accent-tint);
    border-color: var(--accent);
  }
  .col header {
    display: flex; align-items: center; justify-content: space-between;
    padding: 4px 2px;
  }
  .ccount {
    font-family: var(--mono);
    font-size: 11px;
    color: var(--mute);
    font-variant-numeric: tabular-nums;
  }

  .cards { display: flex; flex-direction: column; gap: 6px; }

  .card {
    text-align: left;
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 8px;
    padding: 10px 12px;
    box-shadow: var(--sh-1);
    cursor: grab;
    color: var(--ink-2);
    transition: transform 80ms ease, box-shadow 80ms ease, border-color 80ms ease;
    width: 100%;
  }
  .card:hover {
    border-color: var(--rule-strong);
    transform: translateY(-1px);
    box-shadow: 0 2px 4px rgba(10,10,13,0.04), 0 1px 0 var(--rule);
  }
  .card.dragging { opacity: 0.5; cursor: grabbing; }

  .card-head {
    display: flex; align-items: center; gap: 8px;
    margin-bottom: 4px;
  }
  .card .logo {
    width: 20px; height: 20px;
    border-radius: 5px;
    display: inline-flex; align-items: center; justify-content: center;
    color: white;
    font-size: 10px; font-weight: 600;
    flex-shrink: 0;
  }
  .card .co {
    font-weight: 500;
    font-size: 13px;
    color: var(--ink);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .card .role {
    font-size: 12.5px;
    color: var(--mute);
    line-height: 1.3;
    margin-bottom: 6px;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
  }
  .card .meta {
    display: flex; gap: 6px; align-items: baseline;
    font-size: 11px;
    color: var(--mute-2);
    font-variant-numeric: tabular-nums;
  }
  .card .meta .cv { color: var(--accent-text); font-family: var(--mono); font-size: 10.5px; }

  .empty-col {
    color: var(--mute-2);
    font-size: 11.5px;
    text-align: center;
    padding: 24px 0;
    border: 1px dashed transparent;
    border-radius: 6px;
    margin: 0;
  }
  .col.drag-over .empty-col {
    color: var(--accent-text);
    border-color: var(--accent);
  }

  .hint {
    margin: 16px 0 0;
    text-align: center;
    font-size: 12px;
    color: var(--mute);
  }
</style>
