<script>
  import PreviewBanner from '../PreviewBanner.svelte';
  import TopNav from '../TopNav.svelte';
  import Logo from '../Logo.svelte';
  import { APP_GROUPS } from '../fixtures.js';
</script>

<svelte:head><title>Redesign preview — Applications</title></svelte:head>

<PreviewBanner tag="3.1 v2" title="Applications — grouped by stage (4c)" note="No chips, no controls, no board; stage labels do the work. Exits at the end." />

<div class="page">
  <TopNav active="applications" />

  <div style="max-width:900px;margin:0 auto;padding:36px 32px 80px">

    <div style="display:flex;align-items:baseline;justify-content:space-between;margin-bottom:30px">
      <div>
        <h1 style="font-size:30px;font-weight:700;letter-spacing:-0.02em;margin:0 0 6px">Applications.</h1>
        <div style="font-size:13.5px;color:#6f7680">16 total · 13 in play · 3 closed</div>
      </div>
      <div style="font-size:13px;color:#8a9099"><a href="#top">↓ Import from spreadsheet</a></div>
    </div>

    {#each APP_GROUPS as g (g.name)}
      <div style="margin-bottom:24px">
        <div style="display:flex;align-items:baseline;gap:8px;margin-bottom:8px">
          <span style="font-size:11px;font-weight:600;letter-spacing:.12em;text-transform:uppercase;color:{g.labelColor}">{g.name} · {g.count}</span>
          {#if g.note}<span style="font-size:12px;color:{g.noteColor}">{g.note}</span>{/if}
        </div>
        <div style="background:{g.bg};border:1px solid #eeeeea;border-radius:10px">
          {#each g.rows as app (app.company)}
            <div class="approw" style="display:flex;align-items:center;gap:12px;padding:10px 16px;border-bottom:1px solid #f4f4f1;opacity:{app.kind === 'exit' ? 0.8 : 1}">
              <Logo domain={app.domain} size={24} gray={app.kind === 'exit'} title={app.company} />
              <span style="flex:1;min-width:0;font-size:13.5px;color:{app.kind === 'exit' ? '#8a9099' : '#4b5158'};white-space:nowrap;overflow:hidden;text-overflow:ellipsis"><strong style="color:{app.kind === 'exit' ? '#8a9099' : '#16181c'}">{app.company}</strong> <span style="color:#8a9099">· {app.meta}</span> {#if app.hot}<strong style="color:{app.hotColor}">{app.hot}</strong>{/if}</span>
              {#if app.action}
                <a href="/preview/redesign/detail" style="font-size:12.5px;font-weight:600;flex:none">{app.action}</a>
              {/if}
            </div>
          {/each}
        </div>
      </div>
    {/each}

  </div>
</div>

<style>
  .page { min-height: 100vh; background: #f6f6f3; color: #16181c; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; -webkit-font-smoothing: antialiased; }
  .page :global(a) { color: #2463eb; text-decoration: none; }
  .approw:hover { background: #fafaf8; }
</style>
