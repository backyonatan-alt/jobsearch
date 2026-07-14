<script>
  import PreviewBanner from '../PreviewBanner.svelte';
  import TopNav from '../TopNav.svelte';
  import Logo from '../Logo.svelte';
  import LogoCluster from '../LogoCluster.svelte';
  import { CATS, HOME_ROWS } from '../fixtures.js';
</script>

<svelte:head><title>Redesign preview — Home</title></svelte:head>

<PreviewBanner tag="D v3" title="Home — no sidebar, top nav" note="3a header band + 3d one-line list. Approved direction." />

<div class="page">
  <TopNav active="home" />

  <!-- 3a header band -->
  <div style="background:#fff;border-bottom:1px solid #e8e8e5">
    <div style="max-width:1100px;margin:0 auto;padding:34px 32px 6px">
      <div style="font-size:13px;color:#8a9099;margin-bottom:4px">Tuesday · July 14, 2026</div>
      <h1 style="font-size:30px;font-weight:700;letter-spacing:-0.02em;margin:0 0 22px">Good morning, Yonatan.</h1>
      <div style="display:grid;grid-template-columns:1fr 1fr 1fr 1fr 1fr;border-top:1px solid #eeeeea">
        {#each CATS as cat (cat.label)}
          <a href="#list" class="cat" style="display:block;padding:18px 20px 16px;border-right:1px solid #f0f0ed;color:#16181c">
            <div style="display:flex;align-items:baseline;gap:8px"><span style="font-size:30px;font-weight:700;letter-spacing:-0.03em;color:{cat.color}">{cat.count}</span><span style="font-size:12.5px;font-weight:600;color:{cat.labelColor}">{cat.label}</span></div>
            <div style="display:flex;align-items:center;margin:10px 0 5px"><LogoCluster domains={cat.domains} extra={cat.extra || 0} gray={cat.gray || false} /></div>
            <div style="font-size:11px;color:#8a9099">{cat.note}</div>
          </a>
        {/each}
      </div>
    </div>
  </div>

  <div style="max-width:1100px;margin:0 auto;padding:32px 32px 80px">

    <!-- noticing card -->
    <div style="display:flex;align-items:center;gap:12px;border:1px solid #e8e8e5;background:#fff;border-radius:12px;padding:13px 20px;margin-bottom:32px;font-size:13.5px;color:#4b5158">
      <span style="color:#e0641f;flex:none">✦</span>
      <span><strong style="color:#16181c">Referrals are working:</strong> your one referral reached interview; your 5 LinkedIn applies have a 20% reply rate. Worth asking for one more intro this week.</span>
      <a href="/preview/redesign/insights" style="margin-left:auto;flex:none;font-size:13px;font-weight:500">More in Insights →</a>
    </div>

    <!-- 3d list -->
    <div id="list" style="display:flex;align-items:baseline;gap:10px;margin-bottom:14px">
      <div style="font-size:11px;font-weight:600;letter-spacing:.12em;text-transform:uppercase;color:#8a9099">Everything · 16</div>
      <div style="font-size:12px;color:#8a9099">one line each · a button means it needs you · exits at the end</div>
    </div>

    <div style="display:flex;flex-direction:column;gap:6px">
      {#each HOME_ROWS as app (app.company)}
        {#if app.kind === 'act'}
          <div class="row" style="display:flex;align-items:center;gap:12px;background:#fff;border:1px solid {app.border};border-radius:10px;padding:11px 16px">
            <Logo domain={app.domain} size={26} title={app.company} />
            <span style="flex:1;min-width:0;font-size:13.5px;color:#16181c;white-space:nowrap;overflow:hidden;text-overflow:ellipsis"><strong style="color:#16181c">{app.company}</strong> <span style="color:#8a9099">· {app.meta}</span> <strong style="color:{app.hotColor}">{app.hot}</strong></span>
            <div style="background:{app.btnBg};color:{app.btnColor};border:1px solid {app.btnBorder};border-radius:8px;padding:7px 14px;font-size:12.5px;font-weight:600;min-width:126px;text-align:center;flex:none;cursor:pointer">{app.btn}</div>
          </div>
        {:else if app.kind === 'quiet'}
          <div class="row" style="display:flex;align-items:center;gap:12px;background:#fff;border:1px solid #eeeeea;border-radius:10px;padding:9px 16px">
            <Logo domain={app.domain} size={22} title={app.company} />
            <span style="flex:1;min-width:0;font-size:13px;color:#4b5158;white-space:nowrap;overflow:hidden;text-overflow:ellipsis"><strong style="color:#4b5158">{app.company}</strong> <span style="color:#8a9099">· {app.meta}</span></span>
            <span style="font-size:11.5px;color:#b8bdc4;min-width:126px;text-align:center;flex:none">{app.quiet}</span>
          </div>
        {:else}
          <div class="row" style="display:flex;align-items:center;gap:12px;background:#fbfbf9;border:1px solid #eeeeea;border-radius:10px;padding:9px 16px;opacity:.75">
            <Logo domain={app.domain} size={22} gray title={app.company} />
            <span style="flex:1;min-width:0;font-size:13px;color:#8a9099;white-space:nowrap;overflow:hidden;text-overflow:ellipsis"><strong style="color:#8a9099">{app.company}</strong> <span style="color:#8a9099">· {app.meta}</span></span>
            <div style="background:#fbfbf9;color:#8a9099;border:1px solid #e2e2de;border-radius:8px;padding:7px 14px;font-size:12.5px;font-weight:600;min-width:126px;text-align:center;flex:none;cursor:pointer">Reopen</div>
          </div>
        {/if}
      {/each}
    </div>

  </div>
</div>

<style>
  .page { min-height: 100vh; background: #f6f6f3; color: #16181c; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; -webkit-font-smoothing: antialiased; }
  .page :global(a) { color: #2463eb; text-decoration: none; }
  .cat:hover { background: #fafaf8; }
  .cat { color: #16181c !important; }
  .row:hover { border-color: #b9c6e8 !important; }
</style>
