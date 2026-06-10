<script>
  // Small in-app confirmation dialog — replaces native window.confirm() for
  // destructive actions so it matches the app's design and is testable.
  let {
    open = false,
    title = 'Are you sure?',
    message = '',
    confirmLabel = 'Delete',
    cancelLabel = 'Cancel',
    busy = false,
    onConfirm,
    onCancel
  } = $props();

  function cancel() { if (!busy) onCancel?.(); }
  function onKeydown(e) {
    if (!open) return;
    if (e.key === 'Escape') cancel();
    else if (e.key === 'Enter') onConfirm?.();
  }
</script>

<svelte:window onkeydown={onKeydown} />

{#if open}
  <div class="cfm-ovl" onclick={cancel} role="presentation">
    <div class="cfm" onclick={(e) => e.stopPropagation()} role="alertdialog" aria-modal="true" aria-label={title}>
      <h3 class="cfm-title">{title}</h3>
      {#if message}<p class="cfm-msg">{message}</p>{/if}
      <div class="cfm-actions">
        <button class="cfm-btn" onclick={cancel} disabled={busy}>{cancelLabel}</button>
        <button class="cfm-btn cfm-danger" onclick={() => onConfirm?.()} disabled={busy}>
          {busy ? 'Working…' : confirmLabel}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .cfm-ovl { position: fixed; inset: 0; background: rgba(10,10,13,0.5); backdrop-filter: blur(3px); -webkit-backdrop-filter: blur(3px); display: grid; place-items: center; z-index: 300; padding: 24px; }
  .cfm { width: 100%; max-width: 380px; background: var(--card); border: 1px solid var(--rule); border-radius: 14px; padding: 22px 22px 18px; box-shadow: 0 24px 80px -8px rgba(10,10,13,0.32), var(--sh-1); }
  .cfm-title { font-size: 16px; font-weight: 600; letter-spacing: -0.015em; margin: 0 0 6px; color: var(--ink); }
  .cfm-msg { font-size: 13.5px; line-height: 1.5; color: var(--mute); margin: 0; }
  .cfm-actions { display: flex; justify-content: flex-end; gap: 9px; margin-top: 20px; }
  .cfm-btn { font: inherit; font-size: 13px; font-weight: 500; padding: 8px 15px; border-radius: 9px; border: 1px solid var(--rule); background: var(--card); color: var(--ink); cursor: pointer; }
  .cfm-btn:hover:not(:disabled) { border-color: var(--rule-strong); }
  .cfm-btn:disabled { opacity: 0.6; cursor: default; }
  .cfm-danger { background: var(--danger, oklch(0.58 0.2 25)); border-color: var(--danger, oklch(0.58 0.2 25)); color: white; }
  .cfm-danger:hover:not(:disabled) { background: var(--danger-strong, oklch(0.52 0.21 25)); border-color: var(--danger-strong, oklch(0.52 0.21 25)); }
</style>
