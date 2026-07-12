<script>
  // Clickable status pill: the pill showing an application's status is also
  // the control for changing it. Shared across detail header and Today rows.
  import { api } from '$lib/api.js';
  import { logEvent } from '$lib/analytics.js';
  import { STATUSES, STATUS_LABEL } from '$lib/app-helpers.js';

  let { id, status, surface, onchanged = null, align = 'right' } = $props();

  let open = $state(false);
  let busy = $state(false);

  function toggle(e) {
    e.stopPropagation();
    open = !open;
  }
  function close(e) {
    e?.stopPropagation();
    open = false;
  }
  async function pick(e, s) {
    e.stopPropagation();
    open = false;
    if (busy || s === status) return;
    busy = true;
    const from = status;
    try {
      await api(`/api/applications/${id}`, { method: 'PATCH', body: JSON.stringify({ status: s }) });
      logEvent('status_change', { from, to: s, surface });
      try { window.dispatchEvent(new CustomEvent('pursuit:refresh')); } catch {}
      onchanged?.(s);
    } finally {
      busy = false;
    }
  }
  function onKey(e) {
    if (e.key === 'Escape' && open) { open = false; }
  }
</script>

<svelte:window onkeydown={onKey} />

<span class="sp-wrap" role="presentation" onclick={(e) => e.stopPropagation()}>
  <button
    type="button"
    class={`pill sp-pill ${status}`}
    disabled={busy}
    aria-haspopup="menu"
    aria-expanded={open}
    title="Change status"
    onclick={toggle}
  >
    <span class="pdot"></span>{STATUS_LABEL[status]}
    <svg class="sp-chev" width="10" height="10" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M3 4.5l3 3 3-3" stroke-linecap="round" stroke-linejoin="round"/></svg>
  </button>

  {#if open}
    <div class="sp-menu" class:left={align === 'left'} role="menu">
      {#each STATUSES as s}
        <button type="button" class="sp-item" class:current={status === s} role="menuitem" onclick={(e) => pick(e, s)}>
          <span class={`pill ${s}`} style="margin: 0"><span class="pdot"></span>{STATUS_LABEL[s]}</span>
          {#if status === s}<span class="check">✓</span>{/if}
        </button>
      {/each}
    </div>
    <div class="sp-scrim" onclick={close} role="presentation"></div>
  {/if}
</span>

<style>
  .sp-wrap { position: relative; display: inline-flex; }
  .sp-pill {
    border: 0; cursor: pointer; font: inherit;
    font-size: 12px; font-weight: 500;
    transition: filter .12s, box-shadow .12s;
  }
  .sp-pill:hover:not(:disabled) { filter: brightness(0.96); box-shadow: 0 0 0 1px var(--rule) inset; }
  .sp-pill:disabled { opacity: 0.6; cursor: default; }
  .sp-chev { opacity: 0.55; margin-left: 1px; flex: none; }
  .sp-menu {
    position: absolute; top: calc(100% + 6px); right: 0; z-index: 60;
    background: var(--card); border: 1px solid var(--rule); border-radius: 8px;
    box-shadow: var(--sh-pop); padding: 4px; min-width: 180px;
    display: flex; flex-direction: column; gap: 1px;
  }
  .sp-menu.left { right: auto; left: 0; }
  .sp-item {
    display: flex; align-items: center; gap: 8px; padding: 6px 8px;
    border-radius: 5px; background: transparent; font: inherit; font-size: 13px;
    color: var(--ink-2); cursor: pointer; text-align: left; width: 100%; border: 0;
  }
  .sp-item:hover { background: var(--surface-2); }
  .sp-item.current { background: var(--surface-2); }
  .sp-item .check { margin-left: auto; color: var(--accent-text); font-weight: 600; font-size: 12px; }
  .sp-scrim { position: fixed; inset: 0; z-index: 55; }
</style>
