<script>
  // Redesign top nav (replaces the sidebar). `active`: home | applications | insights.
  // `primaryCta`: whether "New application" renders filled (Home/Applications) or outline (Detail).
  let { active = 'home', primaryCta = true } = $props();
  const items = [
    { k: 'home', label: 'Home', href: '/preview/redesign/home' },
    { k: 'applications', label: 'Applications', count: 16, href: '/preview/redesign/applications' },
    { k: 'insights', label: 'Insights', href: '/preview/redesign/insights' }
  ];
</script>

<div style="background:#fff;border-bottom:1px solid #e8e8e5">
  <div style="max-width:1100px;margin:0 auto;display:flex;align-items:center;gap:24px;padding:12px 32px">
    <div style="display:flex;align-items:center;gap:9px;flex:none">
      <div style="width:24px;height:24px;border-radius:50%;border:2.5px solid #16181c;display:flex;align-items:center;justify-content:center"><div style="width:9px;height:9px;border-radius:50%;background:#2463eb"></div></div>
      <div style="font-size:16px;font-weight:700;letter-spacing:-0.01em">Pursuit</div>
    </div>
    <div style="display:flex;align-items:center;gap:4px">
      {#each items as it (it.k)}
        {#if active === it.k}
          <div style="display:flex;align-items:center;gap:6px;background:#16181c;color:#fff;border-radius:8px;padding:7px 14px;font-size:13.5px;font-weight:600">{it.label}{#if it.count}<span style="font-size:11.5px;opacity:.6">{it.count}</span>{/if}</div>
        {:else}
          <a href={it.href} class="navlink" style="display:flex;align-items:center;gap:6px;border-radius:8px;padding:7px 14px;font-size:13.5px;color:#4b5158">{it.label}{#if it.count}<span style="font-size:11.5px;color:#8a9099">{it.count}</span>{/if}</a>
        {/if}
      {/each}
    </div>
    <div style="display:flex;align-items:center;gap:12px;margin-left:auto">
      <div style="display:flex;align-items:center;gap:8px;border:1px solid #e8e8e5;border-radius:8px;padding:7px 12px;width:240px;color:#8a9099;font-size:13px"><span>⌕</span>Search… <span style="margin-left:auto;font-size:11px;border:1px solid #e8e8e5;border-radius:4px;padding:1px 5px">⌘K</span></div>
      {#if primaryCta}
        <div style="background:#2463eb;color:#fff;border-radius:8px;padding:8px 15px;font-size:13px;font-weight:600;cursor:pointer;flex:none">New application <span style="opacity:.65;font-weight:400;margin-left:3px">⌘N</span></div>
      {:else}
        <div style="border:1px solid #e8e8e5;background:#fff;color:#4b5158;border-radius:8px;padding:8px 15px;font-size:13px;font-weight:600;cursor:pointer;flex:none">New application <span style="opacity:.55;font-weight:400;margin-left:3px">⌘N</span></div>
      {/if}
      <div title="Prep credits" style="border:1px solid #e8e8e5;border-radius:8px;padding:7px 11px;font-size:12px;color:#6f7680;flex:none">✦ <strong style="color:#16181c">7</strong><span style="color:#8a9099">/10</span></div>
      <div title="back.yonatan@gmail.com" style="width:30px;height:30px;border-radius:50%;background:#e0641f;color:#fff;display:flex;align-items:center;justify-content:center;font-size:12px;font-weight:700;flex:none">Y</div>
    </div>
  </div>
</div>

<style>
  .navlink:hover { background: #f0f0ed; }
  a { color: inherit; text-decoration: none; }
</style>
