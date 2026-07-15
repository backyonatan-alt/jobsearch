<script>
  // Redesign logo chip: square, hairline border, real favicon with a letter
  // fallback. The letter renders immediately and the favicon swaps in only
  // once it has actually loaded — no blank square while fetching or on 404.
  let { app, size = 24, radius = 7, gray = false } = $props();
  let loaded = $state(false);
  let failed = $state(false);
</script>

<span class="chip" class:gray style="width:{size}px;height:{size}px;border-radius:{radius}px" title={app.co}>
  <span class="letter" class:hide={loaded} style="font-size:{Math.max(9, Math.round(size * 0.42))}px">{app.coShort}</span>
  {#if app.logoSrc && !failed}
    <img src={app.logoSrc} alt="" class:show={loaded}
      onload={() => (loaded = true)} onerror={() => (failed = true)} />
  {/if}
</span>

<style>
  .chip {
    position: relative;
    display: inline-flex; align-items: center; justify-content: center;
    border: 1px solid #eeeeea; background: #fff; flex: none; overflow: hidden;
  }
  .chip .letter { color: #4b5158; font-weight: 600; }
  .chip .letter.hide { visibility: hidden; }
  .chip img {
    position: absolute; inset: 0; width: 100%; height: 100%;
    object-fit: contain; padding: 2px; box-sizing: border-box;
    opacity: 0;
  }
  .chip img.show { opacity: 1; }
  .chip.gray { filter: grayscale(1); }
</style>
