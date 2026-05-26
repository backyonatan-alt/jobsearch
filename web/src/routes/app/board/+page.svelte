<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api.js';
  import { STATUS_LABEL, STATUSES, toDisplayApp } from '$lib/app-helpers.js';

  let apps = $state([]);
  let loading = $state(true);
  let dragOver = $state(null);
  let dragging = $state(null);

  // "Closed" — collapse rejected + withdrawn into one column to keep the
  // board scannable. Six columns reads cleanly across a desktop width.
  const COL_ORDER = ['wishlist', 'applied', 'screen', 'interview', 'offer', 'rejected'];
  const COL_LABEL = { ...STATUS_LABEL, rejected: 'Closed' };
  const MUTED = new Set(['rejected']);

  onMount(refresh);
  async function refresh() {
    try { apps = (await api('/api/applications')).map(toDisplayApp); }
    catch (e) { if (e.message !== 'unauthorized') console.error(e); }
    finally { loading = false; }
  }

  const byStatus = $derived.by(() => {
    const g = Object.fromEntries(COL_ORDER.map(s => [s, []]));
    for (const a of apps) {
      // Bucket withdrawn under the same column as rejected so it shows up.
      const key = a.status === 'withdrawn' ? 'rejected' : a.status;
      (g[key] ??= []).push(a);
    }
    return g;
  });

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
  const dateLong = today.toLocaleDateString('en-US', { weekday: 'long', day: 'numeric', month: 'long' });
  const inFlight = $derived(apps.filter(a => !['rejected','withdrawn'].includes(a.status)).length);
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
        <div class="date">{dateLong}</div>
        <h1>Board.</h1>
        <p class="sub"><b>{inFlight}</b> in flight · drag a card across columns to move its status.</p>
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
        {#each COL_ORDER as s (s)}
          <section
            class="board-col"
            class:muted={MUTED.has(s)}
            class:drag-over={dragOver === s}
            ondragover={(e) => onDragOver(e, s)}
            ondragleave={() => onDragLeave(s)}
            ondrop={(e) => onDrop(e, s)}
          >
            <header class="board-col-hd">
              <span class={`pill ${s}`}><span class="pdot"></span>{COL_LABEL[s]}</span>
              <span class="count">{byStatus[s].length}</span>
              <button class="add" title="Add to {COL_LABEL[s]}" onclick={() => goto('/app')}>+</button>
            </header>
            <div class="cards">
              {#each byStatus[s] as a (a.id)}
                <button
                  type="button"
                  class="bcard"
                  class:dragging={dragging === a.id}
                  class:stale={a.stale}
                  draggable="true"
                  ondragstart={(e) => onDragStart(e, a.id)}
                  ondragend={onDragEnd}
                  onclick={() => open(a.id)}
                >
                  <div class="top">
                    {#if a.logoSrc}
                      <img class="logo" src={a.logoSrc} alt="" />
                    {:else}
                      <span class={`logo letter ${a.logoCls}`}>{a.coShort}</span>
                    {/if}
                    <span class="co">{a.co}</span>
                    {#if a.stale}<span class="stale-dot" title="No movement for over a week"></span>{/if}
                  </div>
                  <p class="role">{a.role}</p>
                  {#if a.status === 'interview'}
                    <div class="next"><span class="next-dot"></span>Interview · time TBD</div>
                  {/if}
                  <div class="foot">
                    <span class="source">{a.source}</span>
                    <span class="applied" class:stale-text={a.stale}>{a.appliedRel}</span>
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
    {/if}
  </div>
</div>

<style>
  .body { padding: 28px; }
  .board-page { max-width: 1400px; margin: 0 auto; }

  .board-hd { margin-bottom: 24px; }
  .board-hd .date { font-size: 13.5px; color: var(--mute); margin-bottom: 6px; }
  .board-hd h1 { font-size: 28px; font-weight: 600; letter-spacing: -0.025em; margin: 0; }
  .board-hd .sub { font-size: 13.5px; color: var(--mute); margin: 6px 0 0; }

  .board-cols {
    display: grid;
    grid-template-columns: repeat(6, minmax(240px, 1fr));
    gap: 14px;
    overflow-x: auto;
    padding-bottom: 8px;
  }
  .board-col {
    background: var(--surface-2);
    border-radius: 14px;
    padding: 12px;
    min-height: 60vh;
    transition: background 120ms ease;
  }
  .board-col.muted { opacity: 0.65; }
  .board-col.drag-over { background: var(--accent-tint); }
  .board-col-hd { display: flex; align-items: center; gap: 8px; margin-bottom: 12px; padding: 0 4px; }
  .board-col-hd .count { font-size: 12.5px; color: var(--mute); font-weight: 500; }
  .board-col-hd .add {
    margin-left: auto; width: 22px; height: 22px; border-radius: 6px;
    background: transparent; border: 1px solid var(--rule); cursor: pointer;
    display: grid; place-items: center; color: var(--mute); font-size: 16px; line-height: 1;
    transition: background 120ms, color 120ms;
  }
  .board-col-hd .add:hover { background: var(--card); color: var(--ink); }

  .cards { display: flex; flex-direction: column; gap: 10px; }
  .bcard {
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 11px;
    padding: 12px 13px;
    text-align: left;
    cursor: grab;
    box-shadow: var(--sh-1);
    transition: transform 140ms ease, box-shadow 140ms ease, border-color 140ms ease;
    display: flex; flex-direction: column; gap: 6px;
    width: 100%;
  }
  .bcard:hover { transform: translateY(-2px); box-shadow: var(--sh-pop); border-color: var(--rule-strong); }
  .bcard:active { cursor: grabbing; transform: scale(0.99) rotate(-0.5deg); }
  .bcard.dragging { opacity: 0.45; }
  .bcard.stale {
    border-color: var(--danger);
    box-shadow: 0 0 0 1px var(--danger-tint), var(--sh-1);
  }
  .bcard.stale:hover { box-shadow: 0 0 0 1px var(--danger-tint), var(--sh-pop); }

  .top { display: flex; align-items: center; gap: 8px; }
  .top .logo { width: 22px; height: 22px; border-radius: 5px; background: var(--surface-2); object-fit: contain; padding: 2px; flex-shrink: 0; }
  .top .logo.letter {
    display: grid; place-items: center;
    padding: 0;
    color: var(--ink-2); font-size: 11px; font-weight: 600;
  }
  .top .co { font-size: 13.5px; font-weight: 600; letter-spacing: -0.01em; }
  .top .stale-dot { width: 7px; height: 7px; border-radius: 50%; background: var(--danger); margin-left: auto; box-shadow: 0 0 0 3px var(--danger-tint); }

  .role { font-size: 12.5px; color: var(--mute); margin: 0; line-height: 1.4; }

  .next {
    font-size: 11.5px; color: var(--accent-text);
    background: var(--accent-tint); padding: 4px 8px; border-radius: 6px;
    display: inline-flex; align-items: center; gap: 6px; align-self: flex-start;
    font-weight: 500;
  }
  .next-dot { width: 5px; height: 5px; border-radius: 50%; background: var(--accent); }

  .foot { display: flex; justify-content: space-between; font-size: 11.5px; color: var(--mute); padding-top: 4px; border-top: 1px dashed var(--rule); margin-top: 2px; }
  .foot .stale-text { color: var(--danger-text); font-weight: 500; }

  .board-empty {
    border: 1.5px dashed var(--rule-strong);
    border-radius: 10px;
    text-align: center;
    color: var(--mute-2);
    font-size: 12px;
    padding: 18px;
  }

  /* Pills (overrides for board styling) */
  .pill { display: inline-flex; align-items: center; gap: 5px; padding: 3px 9px; border-radius: 99px; font-size: 12px; font-weight: 500; background: var(--card); color: var(--ink-2); }
  .pill .pdot { width: 5px; height: 5px; border-radius: 50%; background: var(--mute-2); }
  .pill.wishlist { background: var(--surface-2); color: var(--mute); }
  .pill.applied { background: var(--card); color: var(--ink-2); }
  .pill.applied .pdot { background: var(--mute-2); }
  .pill.screen { background: var(--positive-tint); color: var(--positive-text); }
  .pill.screen .pdot { background: var(--positive); }
  .pill.interview { background: var(--accent-tint); color: var(--accent-text); }
  .pill.interview .pdot { background: var(--accent); }
  .pill.offer { background: var(--warm-tint); color: var(--warm-text); }
  .pill.offer .pdot { background: var(--warm); }
  .pill.rejected { background: var(--surface-2); color: var(--mute); }

  .empty-tab {
    border: 1px dashed var(--rule);
    border-radius: 12px;
    padding: 32px;
    text-align: center;
    background: var(--card);
  }
  .empty-tab h3 { margin: 0 0 .5rem; font-size: 16px; font-weight: 500; color: var(--ink); }
  .empty-tab p { color: var(--mute); margin: 0; font-size: 13.5px; }

  /* Mobile — tighten the board so each column is closer to phone width.
     Horizontal scroll stays (already overflow-x: auto on .board-cols). */
  @media (max-width: 720px) {
    .body { padding: 18px 14px; }
    .board-hd h1 { font-size: 22px; }
    .board-cols { grid-template-columns: repeat(6, minmax(220px, 220px)); gap: 10px; }
    .board-col { padding: 10px; min-height: 50vh; }
    .bcard { padding: 10px 11px; }
  }
</style>
