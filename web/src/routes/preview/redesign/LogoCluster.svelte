<script>
  import { logo } from './fixtures.js';
  let { domains = [], extra = 0, gray = false } = $props();
</script>

<div style="display:flex;align-items:center">
  {#each domains.slice(0, 3) as d, i (d)}
    <div style="width:24px;height:24px;border-radius:7px;border:2px solid #fff;background-color:#fff;background-image:url({logo(d)});background-size:contain;background-repeat:no-repeat;background-position:center;margin-left:{i ? '-7px' : '0'};box-shadow:0 0 0 1px #eeeeea;filter:{gray ? 'grayscale(1)' : 'none'}"></div>
  {/each}
  {#if extra}<span style="font-size:11px;font-weight:600;color:#8a9099;margin-left:6px">+{extra}</span>{/if}
</div>
