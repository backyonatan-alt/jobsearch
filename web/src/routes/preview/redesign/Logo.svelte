<script>
  import { logo } from './fixtures.js';
  let { domain, size = 24, radius = 7, gray = false, title = '' } = $props();
</script>

<div {title} style="width:{size}px;height:{size}px;border-radius:{radius}px;border:1px solid #eeeeea;background-color:#fff;background-size:contain;background-repeat:no-repeat;background-position:center;background-image:url({logo(domain)});flex:none;filter:{gray ? 'grayscale(1)' : 'none'}"></div>
