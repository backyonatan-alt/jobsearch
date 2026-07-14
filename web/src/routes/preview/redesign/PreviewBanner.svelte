<script>
  // Preview-only chrome (NOT part of the design — dropped at implementation).
  let { tag, title, note = '' } = $props();
</script>

<div style="display:flex;align-items:center;gap:12px;background:#16181c;color:#fff;padding:10px 20px;font-size:12.5px">
  <span style="background:#2463eb;border-radius:5px;padding:2px 8px;font-weight:700;font-size:11px;letter-spacing:.05em">{tag}</span>
  <span style="font-weight:600">{title}</span>
  {#if note}<span style="color:#9aa1ab">{note}</span>{/if}
  <span style="margin-left:auto;color:#9aa1ab">
    <a href="/preview/redesign" style="color:#8ab0ff">Index</a> ·
    <a href="/preview/redesign/home" style="color:#8ab0ff">Home</a> ·
    <a href="/preview/redesign/applications" style="color:#8ab0ff">Applications</a> ·
    <a href="/preview/redesign/detail" style="color:#8ab0ff">Detail</a> ·
    <a href="/preview/redesign/insights" style="color:#8ab0ff">Insights</a> ·
    <a href="/preview/redesign/states" style="color:#8ab0ff">States</a> ·
    <a href="/preview/redesign/ia" style="color:#8ab0ff">IA</a>
  </span>
</div>
