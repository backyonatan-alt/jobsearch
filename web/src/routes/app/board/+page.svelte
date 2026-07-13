<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api.js';
  import { isPreview, mockApi } from '$lib/preview-mode.js';
  import { logEvent } from '$lib/analytics.js';
  import { STATUS_LABEL, toDisplayApp, fmtLongDate, isStale } from '$lib/app-helpers.js';
  import ImportApplications from '$lib/ImportApplications.svelte';
  import AddApplication from '$lib/AddApplication.svelte';

  const call = isPreview() ? mockApi : api;

  let apps = $state([]);
  let loading = $state(true);
  let showImport = $state(false);
  let showAdd = $state(false);
  let addStatus = $state('applied');
  let dragOver = $state(null);   // column key being hovered during drag
  let dragging = $state(null);   // id of card being dragged

  // Board shows only the 5 active pipeline columns. Rejected/withdrawn are
  // hidden — they clutter the board and are visible in the list view.
  const COLS = [
    { k: 'wishlist',   lbl: 'Wishlist'   },
    { k: 'applied',    lbl: 'Applied'    },
    { k: 'screen',     lbl: 'Screen'     },
    { k: 'interview',  lbl: 'Interview'  },
    { k: 'offer',      lbl: 'Offer'      },
  ];

  onMount(() => {
    refresh();
    const h = () => refresh();
    window.addEventListener('pursuit:refresh', h);
    return () => window.removeEventListener('pursuit:refresh', h);
  });
  async function refresh() {
    try { apps = (await call('/api/applications')).map(toDisplayApp); }
    catch (e) { if (e.message !== 'unauthorized') console.error(e); }
    finally { loading = false; }
  }

  // Group by status. Only the 5 active columns — rejected/withdrawn are excluded.
  const byStatus = $derived.by(() => {
    const g = Object.fromEntries(COLS.map(c => [c.k, []]));
    for (const a of apps) {
      if (g[a.status]) g[a.status].push(a);
    }
    return g;
  });

  const inFlight = $derived(apps.filter(a => !['rejected', 'withdrawn', 'closed'].includes(a.status)).length);

  const dateLong = fmtLongDate(new Date());

  function open(id) { goto(`/app/${id}`); }

  function openAdd(status) {
    addStatus = status;
    showAdd = true;
  }

  // ── Drag handlers ──────────────────────────────────────────

  function onDragStart(e, id) {
    dragging = id;
    e.dataTransfer.effectAllowed = 'move';
    e.dataTransfer.setData('text/plain', String(id));
  }

  function onDragEnd() {
    dragging = null;
    dragOver = null;
  }

  function onDragOver(e, colKey) {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
    dragOver = colKey;
  }

  function onDragLeave(e, colKey) {
    // Only clear if we're leaving the column itself, not a child element.
    if (!e.currentTarget.contains(e.relatedTarget)) {
      if (dragOver === colKey) dragOver = null;
    }
  }

  async function onDrop(e, newStatus) {
    e.preventDefault();
    const id = Number(e.dataTransfer.getData('text/plain'));
    dragOver = null;
    dragging = null;
    if (!id) return;
    await moveCard(id, newStatus, 'board');
  }

  // Shared by drop (desktop drag) and the per-card Move select (touch — HTML5
  // drag events never fire on mobile, so the select is the only path there).
  async function moveCard(id, newStatus, surface) {
    const app = apps.find(a => a.id === id);
    if (!app || app.status === newStatus) return;
    const fromStatus = app.status;

    // Optimistic update: move card immediately + mark as "just moved" so
    // its time-ago reads "just now" and the stale border clears.
    const prevApps = apps.slice();
    apps = apps.map(a =>
      a.id === id
        ? { ...a, status: newStatus, appliedRel: 'just now', stale: false, raw: { ...a.raw, status: newStatus } }
        : a
    );

    try {
      await call(`/api/applications/${id}`, {
        method: 'PATCH',
        body: JSON.stringify({ status: newStatus })
      });
      // Confirmed-success only — never on the optimistic update or a revert.
      logEvent('status_change', { from: fromStatus, to: newStatus, surface });
    } catch (err) {
      console.error('status update failed', err);
      apps = prevApps;
    }
  }
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
      <div class="bhd-main">
        <div class="bdate">{dateLong}</div>
        <h1>Board.</h1>
        <div class="bsub">
          <b>{inFlight}</b> in flight · <span class="drag-hint">drag a card across columns to move its status.</span><span class="move-hint">use "Move to" on a card to change its status.</span>
          <span class="legend"><span class="rd"></span>red border = no movement in 7+ days</span>
        </div>
      </div>
      <div class="hd-actions">
        <button class="import-btn" onclick={() => (showImport = true)}>
          <svg width="14" height="14" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round"><path d="M8 2v8M5 7l3 3 3-3M3 12v1.5h10V12"/></svg>
          Import from spreadsheet
        </button>
        <button class="btn btn-primary" onclick={() => openAdd('applied')}>New application</button>
      </div>
    </div>

    {#if loading}
      <p class="loading-msg">Loading…</p>
    {:else if apps.length === 0}
      <div class="empty-state">
        <h3>No applications yet</h3>
        <p>Add your first one with "New application" above — they'll appear here grouped by status.</p>
      </div>
    {:else}
      <div class="bcols" data-tour="board">
        {#each COLS as col (col.k)}
          <section
            class="bcol"
            class:over={dragOver === col.k}
            ondragover={(e) => onDragOver(e, col.k)}
            ondragleave={(e) => onDragLeave(e, col.k)}
            ondrop={(e) => onDrop(e, col.k)}
          >
            <header class="bcol-h">
              <span class="bcol-tag {col.k}">
                <span class="dot"></span>{col.lbl}
              </span>
              <span class="bcol-ct">{byStatus[col.k].length}</span>
              <button class="bcol-add" title="Add to {col.lbl}" onclick={() => openAdd(col.k)}>+</button>
            </header>

            <div class="bcol-list">
              {#each byStatus[col.k] as a (a.id)}
                <div
                  role="button"
                  tabindex="0"
                  class="bcard"
                  class:stale={a.stale}
                  class:dragging={dragging === a.id}
                  draggable="true"
                  ondragstart={(e) => onDragStart(e, a.id)}
                  ondragend={onDragEnd}
                  onclick={() => open(a.id)}
                  onkeydown={(e) => (e.key === 'Enter' || e.key === ' ') && (e.preventDefault(), open(a.id))}
                >
                  <div class="bcard-top">
                    {#if a.logoSrc}
                      <img class="logo" src={a.logoSrc} alt="" />
                    {:else}
                      <span class="logo letter {a.logoCls}">{a.coShort}</span>
                    {/if}
                    <span class="co">{a.co}</span>
                    {#if a.stale}
                      <span class="stale-dot" title="No movement for over a week"></span>
                    {/if}
                  </div>

                  <div class="bcard-role">{a.role}</div>

                  <div class="bcard-ft">
                    <span class="src">{a.source}</span>
                    <span class="ago" class:red={a.stale}>{a.appliedRel}</span>
                  </div>

                  <!-- Touch path: HTML5 drag never fires on mobile, so this
                       select is the only way to move a card from a phone. -->
                  <label class="bcard-move" onclick={(e) => e.stopPropagation()} onkeydown={(e) => e.stopPropagation()}>
                    <span>Move to</span>
                    <select
                      value={a.status}
                      onchange={(e) => moveCard(a.id, e.currentTarget.value, 'board_move')}
                    >
                      {#each COLS as c (c.k)}
                        <option value={c.k}>{c.lbl}</option>
                      {/each}
                    </select>
                  </label>
                </div>
              {/each}

              {#if byStatus[col.k].length === 0}
                <div class="bcol-empty">Drop here</div>
              {/if}
            </div>
          </section>
        {/each}
      </div>
    {/if}

  </div>
</div>

<ImportApplications bind:open={showImport} onImported={refresh} />
<AddApplication bind:open={showAdd} initialStatus={addStatus} onCreated={refresh} />

<style>
  /* ── Layout ──────────────────────────────────────────────── */
  .body { padding: 30px 36px 60px; }
  .board-page { max-width: none; margin: 0; }

  /* ── Header ──────────────────────────────────────────────── */
  .board-hd { margin-bottom: 26px; display: flex; align-items: flex-start; justify-content: space-between; gap: 16px; }
  .hd-actions { display: flex; align-items: center; gap: 10px; flex-shrink: 0; }
  .import-btn { display: inline-flex; align-items: center; gap: 7px; flex-shrink: 0; background: var(--card); border: 1px solid var(--rule); border-radius: 9px; padding: 9px 14px; font-size: 13px; font-weight: 500; color: var(--ink-2); cursor: pointer; font-family: inherit; }
  .import-btn:hover { border-color: var(--rule-strong); color: var(--ink); }
  .import-btn svg { color: var(--mute); }
  .bdate { font-size: 13px; color: var(--mute); margin-bottom: 6px; letter-spacing: -0.003em; }
  .board-hd h1 {
    font-size: 34px; font-weight: 500; letter-spacing: -0.035em; margin: 0 0 11px; color: var(--ink);
  }
  .bsub {
    font-size: 13px; color: var(--mute); display: flex; align-items: center; gap: 7px; flex-wrap: wrap;
  }
  .bsub b { color: var(--ink-2); font-weight: 600; }
  .bsub .legend {
    display: inline-flex; align-items: center; gap: 7px;
    margin-left: 8px; padding-left: 15px; border-left: 1px solid var(--rule);
  }
  .bsub .legend .rd {
    width: 8px; height: 8px; border-radius: 50%; background: var(--danger); flex-shrink: 0;
  }

  /* ── Columns grid ────────────────────────────────────────── */
  .bcols {
    display: grid;
    grid-template-columns: repeat(5, minmax(0, 1fr));
    gap: 14px;
    align-items: start;
    overflow-x: auto;
    padding-bottom: 8px;
  }
  .bcol {
    background: oklch(0.975 0.003 255);
    border: 1px solid var(--rule);
    border-radius: 14px;
    padding: 6px 8px 10px;
    min-height: 140px;
    transition: background 120ms ease, border-color 120ms ease;
  }
  .bcol.over {
    background: var(--accent-tint);
    border-color: oklch(0.62 0.19 258 / 0.4);
  }

  /* ── Column header ───────────────────────────────────────── */
  .bcol-h {
    display: flex; align-items: center; gap: 9px; padding: 8px 6px 16px;
  }
  .bcol-tag {
    display: inline-flex; align-items: center; gap: 7px;
    font-size: 12.5px; font-weight: 500; padding: 4px 11px 4px 9px;
    border-radius: 999px; color: var(--ink-2);
  }
  .bcol-tag .dot { width: 7px; height: 7px; border-radius: 50%; }

  .bcol-tag.wishlist  { background: var(--surface-2); }
  .bcol-tag.wishlist .dot { background: var(--mute-2); }

  .bcol-tag.applied   { background: var(--surface-2); }
  .bcol-tag.applied .dot  { background: var(--mute); }

  .bcol-tag.screen    { background: var(--accent-tint); color: var(--accent-text); }
  .bcol-tag.screen .dot   { background: var(--accent); }

  .bcol-tag.interview { background: var(--warm-tint); color: var(--warm-text); }
  .bcol-tag.interview .dot { background: var(--warm); }

  .bcol-tag.offer     { background: var(--positive-tint); color: var(--positive-text); }
  .bcol-tag.offer .dot    { background: var(--positive); }

  .bcol-ct {
    font-size: 12px; color: var(--mute-2); font-variant-numeric: tabular-nums;
  }
  .bcol-add {
    margin-left: auto; width: 24px; height: 24px; border-radius: 7px; border: none;
    background: none; color: var(--mute-2); font-size: 18px; line-height: 1; cursor: pointer;
    display: flex; align-items: center; justify-content: center;
    transition: background 120ms, color 120ms;
  }
  .bcol-add:hover { background: var(--surface-2); color: var(--ink-2); }

  /* ── Cards list ──────────────────────────────────────────── */
  .bcol-list { display: flex; flex-direction: column; gap: 10px; }

  /* ── Card ────────────────────────────────────────────────── */
  .bcard {
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 12px;
    padding: 13px 14px;
    box-shadow: var(--sh-1);
    cursor: grab;
    text-align: left;
    display: flex; flex-direction: column; gap: 0;
    width: 100%;
    overflow: hidden;
    transition: box-shadow 140ms ease, border-color 140ms ease, transform 140ms ease;
  }
  .bcard:hover { box-shadow: var(--sh-pop); }
  .bcard:active { cursor: grabbing; }
  .bcard.dragging { opacity: 0.45; transform: scale(0.98) rotate(-0.5deg); }

  /* Stale: red border + inner glow */
  .bcard.stale {
    border-color: oklch(0.68 0.19 22 / 0.7);
    box-shadow: 0 0 0 1px oklch(0.68 0.19 22 / 0.18);
  }
  .bcard.stale:hover {
    box-shadow: 0 0 0 1px oklch(0.68 0.19 22 / 0.18), var(--sh-pop);
  }

  /* Card top row: logo + company + optional stale dot */
  .bcard-top {
    display: flex; align-items: center; gap: 9px; margin-bottom: 9px; min-width: 0;
  }
  .bcard-top .logo {
    width: 20px; height: 20px; border-radius: 6px;
    background: var(--surface-2); object-fit: contain; flex-shrink: 0;
  }
  .bcard-top .logo.letter {
    display: grid; place-items: center; padding: 0;
    color: var(--ink-2); font-size: 10px; font-weight: 600;
  }
  .bcard-top .co {
    font-size: 14px; font-weight: 600; letter-spacing: -0.01em;
    white-space: nowrap; overflow: hidden; text-overflow: ellipsis; min-width: 0; flex: 1;
  }
  .bcard-top .stale-dot {
    width: 7px; height: 7px; border-radius: 50%; background: var(--danger);
    flex-shrink: 0; box-shadow: 0 0 0 3px oklch(0.68 0.19 22 / 0.15);
  }

  /* Role */
  .bcard-role {
    font-size: 12.5px; color: var(--mute); line-height: 1.35; margin-bottom: 12px;
  }

  /* Footer: source + time-ago */
  .bcard-ft {
    display: flex; align-items: center; justify-content: space-between; gap: 8px;
    padding-top: 11px; border-top: 1px solid var(--rule);
    font-size: 12px; color: var(--mute);
  }
  .bcard-ft .src {
    overflow: hidden; text-overflow: ellipsis; white-space: nowrap; min-width: 0;
  }
  .bcard-ft .ago {
    font-variant-numeric: tabular-nums; flex-shrink: 0; white-space: nowrap;
  }
  .bcard-ft .ago.red { color: var(--danger-text); font-weight: 500; }

  /* Empty lane placeholder */
  .bcol-empty {
    font-size: 12px; color: var(--mute-2); text-align: center;
    padding: 16px 0; border: 1px dashed var(--rule); border-radius: 10px;
  }

  /* Loading / empty-state ────────────────────────────────── */
  .loading-msg { color: var(--mute); padding: 1rem 0; font-size: 13.5px; }
  .empty-state {
    border: 1px dashed var(--rule); border-radius: 12px;
    padding: 32px; text-align: center; background: var(--card);
  }
  .empty-state h3 { margin: 0 0 .5rem; font-size: 16px; font-weight: 500; color: var(--ink); }
  .empty-state p  { color: var(--mute); margin: 0; font-size: 13.5px; }

  /* Topbar (shared chrome override — board needs full-bleed width) */
  .topbar {
    display: flex; justify-content: space-between; align-items: center;
    padding: 12px 28px; border-bottom: 1px solid var(--rule); background: var(--surface);
  }
  .crumb .here { font-weight: 600; font-size: 14px; }
  .right { display: flex; align-items: center; gap: 8px; }
  .search {
    display: flex; align-items: center; gap: 6px; background: var(--card);
    border: 1px solid var(--rule); border-radius: 7px; padding: 5px 10px;
    font-size: 13px; color: var(--mute); min-width: 260px;
  }
  .search .ico { width: 14px; height: 14px; flex-shrink: 0; }
  .search .kbd {
    margin-left: auto; font-size: 11px; color: var(--mute-2); background: var(--surface-2);
    border: 1px solid var(--rule); border-radius: 4px; padding: 1px 5px;
  }

  /* Touch move control — the only status-change path on phones. */
  .bcard-move { display: none; }
  .move-hint { display: none; }

  /* Tablet: horizontal scroll, narrower cards */
  @media (max-width: 900px) {
    .body { padding: 18px 14px 40px; }
    .board-hd { flex-wrap: wrap; }
    .board-hd h1 { font-size: 26px; }
    .bcols { grid-template-columns: repeat(5, minmax(200px, 200px)); }
    .bcol { padding: 6px 8px 10px; }
    .bcard { padding: 11px 12px; }
  }

  /* Phone: stacked full-width lanes — no horizontal scroll, no drag. */
  @media (max-width: 720px) {
    /* design-system's global .board-page padding (28/32) eats 64px of a phone. */
    .board-page { padding: 0 0 40px; }
    .bcols { display: flex; flex-direction: column; align-items: stretch; gap: 12px; overflow-x: visible; }
    .bcol { min-height: 0; }
    .bcol-h { padding: 8px 6px 10px; }
    .bcol-list { gap: 8px; }
    .bcol-empty { padding: 10px 0; font-size: 11.5px; }
    .bcard { cursor: pointer; }
    .drag-hint { display: none; }
    .move-hint { display: inline; }
    .bsub .legend { margin-left: 0; padding-left: 0; border-left: none; width: 100%; }
    .bcard-move {
      display: flex; align-items: center; gap: 8px; margin-top: 10px;
      font-size: 12px; color: var(--mute);
    }
    .bcard-move select {
      flex: 1; min-height: 38px; font-size: 13px; font-family: inherit; color: var(--ink);
      background: var(--surface-2); border: 1px solid var(--rule); border-radius: 8px; padding: 6px 10px;
    }
  }
</style>
