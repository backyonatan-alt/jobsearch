<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api.js';
  import { STATUS_LABEL, STATUSES, toDisplayApp } from '$lib/app-helpers.js';

  let apps = $state([]);
  let loading = $state(true);
  let dragOver = $state(null);
  let dragging = $state(null);

  onMount(refresh);
  async function refresh() {
    try { apps = (await api('/api/applications')).map(toDisplayApp); }
    catch (e) { if (e.message !== 'unauthorized') console.error(e); }
    finally { loading = false; }
  }

  const byStatus = $derived.by(() => {
    const g = Object.fromEntries(STATUSES.map(s => [s, []]));
    for (const a of apps) (g[a.status] ??= []).push(a);
    return g;
  });

  // Mute closed columns so they don't fight for attention.
  const MUTED = new Set(['rejected', 'withdrawn']);

  function open(id) { goto(`/app/${id}`); }

  function onDragStart(e, id) {
    dragging = id;
    e.dataTransfer.effectAllowed = 'move';
    e.dataTransfer.setData('text/plain', String(id));
  }
  function onDragEnd() { dragging = null; dragOver = null; }
  function onDragOver(e, status) { e.preventDefault(); dragOver = status; }
  function onDragLeave(status) { if (dragOver === status) dragOver = null; }
  async function onDrop(e, newStatus) {
    e.preventDefault();
    const id = Number(e.dataTransfer.getData('text/plain'));
    dragOver = null; dragging = null;
    if (!id) return;
    const app = apps.find(a => a.id === id);
    if (!app || app.status === newStatus) return;
    apps = apps.map(a => a.id === id ? { ...a, status: newStatus, raw: { ...a.raw, status: newStatus } } : a);
    try {
      await api(`/api/applications/${id}`, { method: 'PATCH', body: JSON.stringify({ status: newStatus }) });
    } catch (err) {
      console.error('status update failed', err);
      await refresh();
    }
  }

  const today = new Date();
  const dow  = today.toLocaleDateString('en-US', { weekday: 'long' });
  const dnum = today.toLocaleDateString('en-US', { day: 'numeric', month: 'long', year: 'numeric' });

  const active = $derived(apps.filter(a => !MUTED.has(a.status)).length);
  const interview = $derived(apps.filter(a => a.status === 'interview').length);
</script>

<svelte:head><title>Board — Pursuit</title></svelte:head>

<div class="topbar">
  <div class="crumb"><span class="here">Board</span></div>
  <div class="right">
    <div class="search">
      <svg class="ico" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5">
        <circle cx="7" cy="7" r="4.5"/><path d="M10.5 10.5L13 13"/>
      </svg>
      <span>Search applications, people…</span>
      <span class="kbd">⌘K</span>
    </div>
  </div>
</div>

<div class="body">
  <div class="board-page">
    <div class="board-hd">
      <div>
        <div class="date">
          <span class="dow">{dow}</span>
          <span class="sep">/</span>
          <span class="dnum">{dnum}</span>
        </div>
        <h1>Board.</h1>
        <p class="sub">
          <b>{active}</b> in flight · <b>{interview}</b> in interview ·
          drag a card across columns to move its status.
        </p>
      </div>
      <div class="view-switch" role="tablist">
        <a class="seg" href="/app" role="tab">List</a>
        <a class="seg active" href="/app/board" role="tab" aria-selected="true">Board</a>
      </div>
    </div>

    {#if loading}
      <p style="color:var(--mute); padding: 1rem 0;">Loading…</p>
    {:else if apps.length === 0}
      <div class="empty-tab">
        <h3>No applications yet</h3>
        <p>Add your first one from the Today page (⌘N) — they'll appear here grouped by status.</p>
      </div>
    {:else}
      <div class="board-cols">
        {#each STATUSES as s (s)}
          <section
            class="board-col"
            class:muted={MUTED.has(s)}
            class:drag-over={dragOver === s}
            ondragover={(e) => onDragOver(e, s)}
            ondragleave={() => onDragLeave(s)}
            ondrop={(e) => onDrop(e, s)}
          >
            <header class="board-col-hd">
              <span class={`pill ${s}`}><span class="pdot"></span>{STATUS_LABEL[s]}</span>
              <span class="count">{byStatus[s].length}</span>
              <button class="add" title="Add to {STATUS_LABEL[s]}" onclick={() => goto('/app')}>+</button>
            </header>
            <div class="cards">
              {#each byStatus[s] as a (a.id)}
                <button
                  type="button"
                  class="bcard"
                  class:dragging={dragging === a.id}
                  class:closed={MUTED.has(a.status)}
                  draggable="true"
                  ondragstart={(e) => onDragStart(e, a.id)}
                  ondragend={onDragEnd}
                  onclick={() => open(a.id)}
                >
                  <div class="top">
                    <span class={`logo ${a.logoCls}`}>{a.coShort}</span>
                    <span class="co">{a.co}</span>
                    {#if a.cv && a.cv !== '—'}<span class="cv">{a.cv}</span>{/if}
                  </div>
                  <p class="role">{a.role}</p>
                  {#if a.status === 'interview'}
                    <div class="next">
                      <span class="when">Time TBD</span>
                      <span class="kind">interview</span>
                    </div>
                  {/if}
                  <div class="foot">
                    <span class="source">
                      <span class="sdot"></span>
                      {a.source}
                    </span>
                    <span class="applied">{a.applied}</span>
                  </div>
                </button>
              {/each}
              {#if byStatus[s].length === 0}
                <div class="board-empty">drop here</div>
              {/if}
            </div>
          </section>
        {/each}
      </div>

      <div class="board-foot">
        <span><kbd class="kbd">←</kbd> <kbd class="kbd">→</kbd> scroll columns</span>
        <span><kbd class="kbd">drag</kbd> a card to change status</span>
        <span>click any card to open the dossier</span>
      </div>
    {/if}
  </div>
</div>
