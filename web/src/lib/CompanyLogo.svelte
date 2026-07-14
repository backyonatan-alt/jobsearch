<script>
  // Redesign logo chip: square, hairline border, real favicon with a letter
  // fallback while it loads / when it 404s. Grayscale for archived rows.
  let { app, size = 24, radius = 7, gray = false } = $props();
  let failed = $state(false);
</script>

<span class="chip" class:gray style="width:{size}px;height:{size}px;border-radius:{radius}px" title={app.co}>
  {#if app.logoSrc && !failed}
    <img src={app.logoSrc} alt="" onerror={() => (failed = true)} />
  {:else}
    <span class="letter" style="font-size:{Math.max(9, Math.round(size * 0.42))}px">{app.coShort}</span>
  {/if}
</span>

<style>
  .chip {
    display: inline-flex; align-items: center; justify-content: center;
    border: 1px solid #eeeeea; background: #fff; flex: none; overflow: hidden;
  }
  .chip img { width: 100%; height: 100%; object-fit: contain; padding: 2px; box-sizing: border-box; }
  .chip .letter { color: #4b5158; font-weight: 600; }
  .chip.gray { filter: grayscale(1); }
</style>
