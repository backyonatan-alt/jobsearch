<script>
  import { onMount } from 'svelte';
  import { api } from '$lib/api.js';
  import { isPreview, mockApi } from '$lib/preview-mode.js';
  import { toDisplayApp, daysSince } from '$lib/app-helpers.js';

  const call = isPreview() ? mockApi : api;

  let apps = $state([]);
  let loading = $state(true);

  onMount(async () => {
    try {
      apps = (await call('/api/applications')).map(toDisplayApp);
    } catch (e) {
      if (e.message !== 'unauthorized') console.error(e);
    } finally {
      loading = false;
    }
  });

  // ── Funnel counts (cumulative reach) ──────────────────────
  // An app at "interview" has passed through applied+screen, so counts are
  // additive: applied = everyone who got past wishlist.
  const funnelCounts = $derived.by(() => {
    let applied = 0, screen = 0, interview = 0, offer = 0;
    for (const a of apps) {
      const s = a.status;
      if (['applied', 'screen', 'interview', 'offer', 'rejected', 'withdrawn'].includes(s)) applied++;
      if (['screen', 'interview', 'offer'].includes(s)) screen++;
      if (['interview', 'offer'].includes(s)) interview++;
      if (s === 'offer') offer++;
    }
    return { applied, screen, interview, offer };
  });

  const funnelStages = $derived([
    { key: 'applied',   label: 'Applied',   n: funnelCounts.applied,   color: 'oklch(0.62 0.19 258)' },
    { key: 'screen',    label: 'Screen',    n: funnelCounts.screen,    color: 'oklch(0.62 0.19 258 / 0.78)' },
    { key: 'interview', label: 'Interview', n: funnelCounts.interview, color: 'oklch(0.7 0.16 50)' },
    { key: 'offer',     label: 'Offer',     n: funnelCounts.offer,     color: 'oklch(0.65 0.14 152)' }
  ]);

  const activeCount = $derived(
    apps.filter(a => !['wishlist', 'rejected', 'withdrawn'].includes(a.status)).length
  );

  // ── KPI 1: Reply rate ──────────────────────────────────────
  // Numerator: apps that reached screen or further (screen, interview, offer, rejected after screen)
  // We only have current status, so: screen + interview + offer counts as "got a reply".
  // Denominator: all apps that were at least "applied" (exclude wishlist only).
  const replyRate = $derived.by(() => {
    const denom = funnelCounts.applied; // applied+screen+interview+offer+rejected+withdrawn
    const numer = funnelCounts.screen;  // screen+interview+offer (cumulative)
    if (!denom) return null;
    return Math.round((numer / denom) * 100);
  });

  // ── KPI 2: Avg. time to first reply ───────────────────────
  // Proxy: avg days from applied_at to now for apps that reached screen or further.
  // Not a true "days to reply" but the best honest estimate from available data.
  const avgReplyDays = $derived.by(() => {
    const replied = apps.filter(a =>
      ['screen', 'interview', 'offer'].includes(a.status) && a.raw.applied_at
    );
    if (!replied.length) return null;
    const sum = replied.reduce((s, a) => s + (daysSince(a.raw.applied_at) ?? 0), 0);
    return Math.round(sum / replied.length);
  });

  // ── KPI 3: Furthest stage reached ──────────────────────────
  const stageRank = { offer: 4, interview: 3, screen: 2, applied: 1, wishlist: 0, rejected: 0, withdrawn: 0 };
  const furthestStage = $derived.by(() => {
    const ORDER = ['offer', 'interview', 'screen', 'applied'];
    for (const s of ORDER) {
      const matching = apps.filter(a => a.status === s);
      if (matching.length) {
        return {
          label: s.charAt(0).toUpperCase() + s.slice(1),
          count: matching.length,
          company: matching.length === 1 ? matching[0].co : null
        };
      }
    }
    return null;
  });

  // ── Activity histogram: last 12 weeks ─────────────────────
  // Buckets applied_at (or created_at) into ISO weeks.
  const activityWeeks = $derived.by(() => {
    const now = new Date();
    const buckets = Array.from({ length: 12 }, (_, i) => {
      // week 0 = current week, week 11 = 11 weeks ago
      const start = new Date(now);
      start.setDate(start.getDate() - (11 - i) * 7 - start.getDay());
      start.setHours(0, 0, 0, 0);
      const end = new Date(start);
      end.setDate(end.getDate() + 7);
      return { start, end, count: 0, label: '' };
    });

    // Label only when the month changes relative to the previous bucket
    let lastMonth = -1;
    for (const b of buckets) {
      const m = b.start.getMonth();
      if (m !== lastMonth) {
        b.label = b.start.toLocaleDateString('en-US', { month: 'short' });
        lastMonth = m;
      }
    }

    for (const a of apps) {
      const d = a.raw.applied_at ? new Date(a.raw.applied_at) : null;
      if (!d) continue;
      for (const b of buckets) {
        if (d >= b.start && d < b.end) { b.count++; break; }
      }
    }
    return buckets;
  });

  const activityMax = $derived(Math.max(1, ...activityWeeks.map(b => b.count)));

  // ── Sources breakdown ──────────────────────────────────────
  // Group by raw source field; null/"—" → "Direct". Show count + % of total.
  // Sorted by count desc.
  const sourceRows = $derived.by(() => {
    const total = apps.filter(a => a.status !== 'wishlist').length;
    const map = {};
    for (const a of apps) {
      if (a.status === 'wishlist') continue;
      const src = (a.source && a.source !== '—') ? a.source : 'Direct';
      map[src] = (map[src] || 0) + 1;
    }
    const maxN = Math.max(1, ...Object.values(map));
    return Object.entries(map)
      .map(([src, n]) => ({
        src,
        n,
        pct: total ? Math.round((n / total) * 100) : 0,
        barWidth: Math.round((n / maxN) * 100)
      }))
      .sort((a, b) => b.n - a.n);
  });

  // ── Insight callout ────────────────────────────────────────
  // Pick the most interesting honest one-liner from the data.
  const insightLine = $derived.by(() => {
    // Referral conversion: how many referral apps reached screen or further
    const referralApps = apps.filter(a =>
      (a.source || '').toLowerCase().includes('referral') && a.status !== 'wishlist'
    );
    const referralProgressed = referralApps.filter(a =>
      ['screen', 'interview', 'offer'].includes(a.status)
    );
    if (referralApps.length >= 2 && referralProgressed.length >= 1) {
      return `Your ${referralProgressed.length} referral application${referralProgressed.length > 1 ? 's' : ''} reached a screen or further. Cold apps stall at the inbox — lean on intros.`;
    }
    // Most active source
    if (sourceRows.length >= 2) {
      const top = sourceRows[0];
      return `${top.pct}% of your applications came via ${top.src} — ${top.n} in total.`;
    }
    // Fallback: funnel note
    if (funnelCounts.applied > 0) {
      return `${funnelCounts.applied} applications tracked. Keep the pipeline moving — add a few more this week.`;
    }
    return 'Add your first application to see insights here.';
  });
</script>

<svelte:head><title>Insights — Pursuit</title></svelte:head>

<div class="topbar">
  <div class="crumb"><span class="here">Insights</span></div>
</div>

<div class="body">
  <div class="body-inner">
    <div class="ins-header">
      <h1>Insights</h1>
      <div class="sub">How your search is actually going · last 12 weeks</div>
    </div>

    {#if loading}
      <p class="loading-msg">Loading…</p>
    {:else if funnelCounts.applied === 0}
      <div class="empty-tab">
        <h3>No data to chart yet</h3>
        <p>Once you add an application or two, the funnel and activity appear here.</p>
      </div>
    {:else}
      <!-- KPI ROW: 3 cards -->
      <div class="ins-kpis">
        <!-- Reply rate -->
        <div class="kpi">
          <div class="kpi-l">Reply rate</div>
          <div class="kpi-v">{replyRate ?? '—'}{#if replyRate !== null}<span class="kpi-unit">%</span>{/if}</div>
          <div class="kpi-delta flat">based on your {funnelCounts.applied} application{funnelCounts.applied !== 1 ? 's' : ''}</div>
        </div>

        <!-- Avg time to first reply -->
        <div class="kpi">
          <div class="kpi-l">Avg. time to first reply</div>
          {#if avgReplyDays !== null}
            <div class="kpi-v">{avgReplyDays}<span class="kpi-unit">d</span></div>
            <div class="kpi-delta flat">proxy: days from applied to screen</div>
          {:else}
            <div class="kpi-v">—</div>
            <div class="kpi-delta flat">not enough data yet</div>
          {/if}
        </div>

        <!-- Furthest stage reached -->
        <div class="kpi">
          <div class="kpi-l">Furthest stage reached</div>
          {#if furthestStage}
            <div class="kpi-v kpi-v-stage">{furthestStage.label} <span class="kpi-count">×{furthestStage.count}</span></div>
            {#if furthestStage.company}
              <div class="kpi-delta flat">{furthestStage.company}</div>
            {:else}
              <div class="kpi-delta flat">{furthestStage.count} application{furthestStage.count !== 1 ? 's' : ''} at this stage</div>
            {/if}
          {:else}
            <div class="kpi-v">—</div>
            <div class="kpi-delta flat">keep applying</div>
          {/if}
        </div>
      </div>

      <!-- MAIN GRID: 1.3fr | 1fr -->
      <div class="ins-grid">
        <!-- LEFT PANEL -->
        <div class="panel">
          <!-- Pipeline funnel -->
          <div class="ph">Pipeline funnel</div>
          <div class="psub">Where your {activeCount} active application{activeCount !== 1 ? 's' : ''} sit today</div>
          <div class="funnel">
            {#each funnelStages as s}
              {@const width = funnelCounts.applied ? Math.max(4, (s.n / funnelCounts.applied) * 100) : 0}
              {@const pct = funnelCounts.applied ? Math.round((s.n / funnelCounts.applied) * 100) : 0}
              <div class="fn">
                <span class="fn-l">{s.label}</span>
                <span class="fn-bar" style="width: {width}%; background: {s.color};"></span>
                <span class="fn-n">{s.n}</span>
                <span class="fn-pct">{pct}%</span>
              </div>
            {/each}
          </div>

          <div class="divider"></div>

          <!-- Application activity -->
          <div class="ph">Application activity</div>
          <div class="psub">New applications per week</div>
          <div class="act-bars">
            {#each activityWeeks as week, i}
              {@const h = Math.max(6, (week.count / activityMax) * 100)}
              <div
                class="act-col"
                title="{week.count} application{week.count !== 1 ? 's' : ''}"
              >
                {#if week.count > 0}
                  <span class="act-val">{week.count}</span>
                {/if}
                <div
                  class="act-bar"
                  class:act-bar-current={i === activityWeeks.length - 1}
                  style="height: {h}%;"
                ></div>
              </div>
            {/each}
          </div>
          <div class="act-x">
            {#each activityWeeks as week}
              <span>{week.label}</span>
            {/each}
          </div>
        </div>

        <!-- RIGHT PANEL -->
        <div class="panel">
          <div class="ph">Where they come from</div>
          <div class="psub">Sorted by volume</div>
          <div class="src-list">
            {#each sourceRows as row}
              <div class="src-row">
                <div class="src-meta">
                  <span class="src-nm">{row.src}</span>
                  <span class="src-track"><i style="width: {row.barWidth}%;"></i></span>
                </div>
                <span class="src-cnt">{row.n}</span>
                <span class="src-pct">{row.pct}%</span>
              </div>
            {/each}
          </div>

          <!-- Insight callout -->
          <div class="insight-callout">
            <span class="spark-icon">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                <path d="M12 2L14.5 9.5L22 12L14.5 14.5L12 22L9.5 14.5L2 12L9.5 9.5L12 2Z"
                  fill="currentColor" opacity="0.9"/>
              </svg>
            </span>
            <span class="insight-txt">{insightLine}</span>
          </div>
        </div>
      </div>
    {/if}
  </div>
</div>

<style>
  .body { padding: 28px; }
  .body-inner { max-width: 1080px; margin: 0 auto; }

  /* Header */
  .ins-header { margin-bottom: 28px; }
  .ins-header h1 { font-size: 26px; font-weight: 500; letter-spacing: -0.03em; margin: 0 0 6px; }
  .sub { font-size: 13.5px; color: var(--mute); }

  .loading-msg { color: var(--mute); font-size: 13.5px; }

  /* KPI row */
  .ins-kpis { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; margin-bottom: 28px; }
  .kpi { background: var(--card); border: 1px solid var(--rule); border-radius: 14px; padding: 20px 22px; box-shadow: var(--sh-1); }
  .kpi-l { font-size: 12px; color: var(--mute); margin-bottom: 10px; }
  .kpi-v { font-size: 34px; font-weight: 400; letter-spacing: -0.03em; line-height: 1; font-variant-numeric: tabular-nums; }
  .kpi-v-stage { font-size: 26px; }
  .kpi-count { font-size: 20px; opacity: 0.7; }
  .kpi-unit { font-size: 18px; opacity: 0.5; margin-left: 2px; }
  .kpi-delta { font-size: 12px; margin-top: 9px; display: inline-flex; align-items: center; gap: 5px; }
  .kpi-delta.flat { color: var(--mute); }
  .kpi-delta.up { color: var(--positive-text); }

  /* Main grid */
  .ins-grid { display: grid; grid-template-columns: 1.3fr 1fr; gap: 16px; }

  /* Panel */
  .panel { background: var(--card); border: 1px solid var(--rule); border-radius: 16px; padding: 22px 24px; box-shadow: var(--sh-1); }
  .ph { font-size: 13.5px; font-weight: 500; margin-bottom: 4px; }
  .psub { font-size: 12px; color: var(--mute); margin-bottom: 20px; }

  /* Funnel bars */
  .funnel { display: flex; flex-direction: column; gap: 12px; }
  .fn { display: grid; grid-template-columns: 78px 1fr 28px 44px; gap: 12px; align-items: center; }
  .fn-l { font-size: 12.5px; color: var(--ink-2); }
  .fn-bar { height: 30px; border-radius: 8px; display: flex; align-items: center; min-width: 30px; transition: width 280ms ease; }
  .fn-n { font-size: 13px; color: var(--ink-2); font-variant-numeric: tabular-nums; font-family: var(--mono, ui-monospace, monospace); text-align: right; }
  .fn-pct { font-size: 12.5px; color: var(--mute); font-variant-numeric: tabular-nums; text-align: right; }

  /* Divider between funnel and activity */
  .divider { height: 1px; background: var(--rule); margin: 24px 0; }

  /* Activity bars */
  .act-bars {
    display: flex; align-items: flex-end; gap: 6px; height: 120px;
    background-image: repeating-linear-gradient(to top, transparent 0, transparent calc(25% - 1px), var(--rule) calc(25% - 1px), var(--rule) 25%);
  }
  .act-col { flex: 1; display: flex; flex-direction: column; justify-content: flex-end; align-items: stretch; height: 100%; min-width: 0; }
  .act-val { font-family: var(--mono, ui-monospace, monospace); font-size: 10px; line-height: 1; color: var(--mute-2); font-variant-numeric: tabular-nums; text-align: center; margin-bottom: 4px; }
  .act-bar { width: 100%; background: var(--accent-tint-2); border-radius: 5px 5px 2px 2px; min-height: 6px; transition: background .12s; }
  .act-col:hover .act-bar { background: var(--accent); }
  .act-bar-current { background: var(--accent) !important; }
  .act-x { display: flex; justify-content: space-between; margin-top: 8px; font-size: 11px; color: var(--mute-2); }

  /* Sources */
  .src-list { margin-bottom: 0; }
  .src-row { display: grid; grid-template-columns: 1fr auto auto; gap: 10px; align-items: center; padding: 11px 0; border-top: 1px solid var(--rule); }
  .src-row:first-child { border-top: none; }
  .src-meta { display: flex; flex-direction: column; gap: 6px; }
  .src-nm { font-size: 13px; font-weight: 400; }
  .src-track { height: 6px; border-radius: 3px; background: var(--surface-2); overflow: hidden; width: 160px; display: block; }
  .src-track i { display: block; height: 100%; background: var(--accent); border-radius: 3px; }
  .src-cnt { font-size: 12.5px; color: var(--mute); font-variant-numeric: tabular-nums; font-family: var(--mono, ui-monospace, monospace); text-align: right; }
  .src-pct { font-size: 12.5px; color: var(--mute-2); font-variant-numeric: tabular-nums; min-width: 32px; text-align: right; }

  /* Insight callout */
  .insight-callout { margin-top: 18px; padding: 14px 16px; background: var(--accent-tint); border-radius: 12px; display: flex; gap: 10px; align-items: flex-start; }
  .spark-icon { color: var(--accent-text); flex-shrink: 0; margin-top: 1px; }
  .insight-txt { font-size: 12.5px; line-height: 1.5; color: var(--accent-text); }

  /* Empty state */
  .empty-tab { border: 1px dashed var(--rule); border-radius: 12px; padding: 32px; text-align: center; background: var(--card); }
  .empty-tab h3 { margin: 0 0 .5rem; font-size: 16px; font-weight: 500; color: var(--ink); }
  .empty-tab p { color: var(--mute); margin: 0; font-size: 13.5px; }

  /* Mobile */
  @media (max-width: 768px) {
    .body { padding: 18px 14px; }
    .ins-kpis { grid-template-columns: 1fr 1fr; gap: 10px; }
    .ins-grid { grid-template-columns: 1fr; }
    .panel { padding: 16px 18px; }
    .kpi-v { font-size: 26px; }
  }
  @media (max-width: 480px) {
    .ins-kpis { grid-template-columns: 1fr; }
  }
</style>
