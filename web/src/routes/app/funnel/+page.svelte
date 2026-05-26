<script>
  import { onMount } from 'svelte';
  import { api } from '$lib/api.js';
  import { toDisplayApp, daysSince } from '$lib/app-helpers.js';

  let apps = $state([]);
  let loading = $state(true);

  onMount(async () => {
    try {
      apps = (await api('/api/applications')).map(toDisplayApp);
    } catch (e) {
      if (e.message !== 'unauthorized') console.error(e);
    } finally {
      loading = false;
    }
  });

  function pct(n, base) {
    if (!base) return null;
    return Math.round((n / base) * 100);
  }

  // Funnel: each downstream stage is a subset since an app can't reach
  // "interview" without going through "applied" + "screen". We don't store
  // status history yet, so we infer from current status.
  const counts = $derived.by(() => {
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

  const stages = $derived([
    { key: 'applied',   label: 'Applied',   n: counts.applied   },
    { key: 'screen',    label: 'Screen',    n: counts.screen    },
    { key: 'interview', label: 'Interview', n: counts.interview },
    { key: 'offer',     label: 'Offer',     n: counts.offer     }
  ]);
  const overallPct = $derived(pct(counts.offer, counts.applied));

  // Source breakdown: bucket free-text sources into a few classes, then
  // compute Applied → Offer conversion.
  function classifySource(s) {
    const v = (s || '').toLowerCase();
    if (!v || v === '—') return 'Other';
    if (v.includes('referral')) return 'Referral';
    if (v.includes('linkedin')) return 'LinkedIn';
    if (v.includes('cold')) return 'Cold';
    return 'Other';
  }
  const sourceBars = $derived.by(() => {
    const buckets = {};
    for (const a of apps) {
      const c = classifySource(a.source);
      if (!buckets[c]) buckets[c] = { name: c, n: 0, offers: 0 };
      buckets[c].n++;
      if (a.status === 'offer') buckets[c].offers++;
    }
    const TONE = { Referral: 'positive', LinkedIn: 'accent', Cold: 'warm', Other: 'mute' };
    return Object.values(buckets)
      .filter(b => b.n > 0)
      .map(b => ({ ...b, rate: pct(b.offers, b.n) ?? 0, tone: TONE[b.name] || 'mute' }))
      .sort((a, b) => b.rate - a.rate);
  });
  const maxSource = $derived(Math.max(1, ...sourceBars.map(s => s.rate)));

  // CV variant: Applied → Screen rate (earliest signal that the CV worked).
  const cvBars = $derived.by(() => {
    const buckets = {};
    for (const a of apps) {
      const v = a.raw.cv_variant || '—';
      if (v === '—') continue;
      if (!buckets[v]) buckets[v] = { name: v, n: 0, advanced: 0 };
      buckets[v].n++;
      if (['screen', 'interview', 'offer'].includes(a.status)) buckets[v].advanced++;
    }
    const tones = ['positive', 'accent', 'warm', 'mute'];
    return Object.values(buckets)
      .map((b, i) => ({ ...b, rate: pct(b.advanced, b.n) ?? 0, tone: tones[i % tones.length] }))
      .sort((a, b) => b.rate - a.rate);
  });
  const maxCv = $derived(Math.max(1, ...cvBars.map(c => c.rate)));

  // Time-in-stage: best-effort proxy — avg days since applied_at for apps
  // currently sitting in each stage. Real "time spent in stage" needs a
  // stage_history table (v0.3).
  const timeInStage = $derived.by(() => {
    function avgDays(predicate) {
      const matched = apps.filter(predicate).map(a => daysSince(a.raw.applied_at)).filter(d => d !== null);
      if (matched.length === 0) return null;
      return Math.round(matched.reduce((s, d) => s + d, 0) / matched.length);
    }
    return [
      { stage: 'Applied',   days: avgDays(a => a.status === 'applied'),                                   tone: 'mute' },
      { stage: 'Screen',    days: avgDays(a => a.status === 'screen'),                                    tone: 'positive' },
      { stage: 'Interview', days: avgDays(a => a.status === 'interview'),                                 tone: 'accent' },
      { stage: 'Offer',     days: avgDays(a => a.status === 'offer'),                                     tone: 'warm' }
    ];
  });
  const maxTime = $derived(Math.max(1, ...timeInStage.map(t => t.days ?? 0)));

  // Headline KPIs.
  const bestCv = $derived(cvBars[0] ?? null);
  // Avg time to offer = avg days from applied_at to now for offers (proxy).
  const avgToOffer = $derived.by(() => {
    const offers = apps.filter(a => a.status === 'offer').map(a => daysSince(a.raw.applied_at)).filter(d => d !== null);
    if (offers.length === 0) return null;
    return Math.round(offers.reduce((s, d) => s + d, 0) / offers.length);
  });
  const inFlight = $derived(apps.filter(a => !['rejected','withdrawn','offer'].includes(a.status)).length);
</script>

<svelte:head><title>Funnel — Pursuit</title></svelte:head>

<div class="topbar">
  <div class="crumb"><span class="here">Funnel</span></div>
  <div class="right">
    <button class="btn">Last 90 days</button>
  </div>
</div>

<div class="body">
  <div class="body-inner">
    <div class="hello">
      <div class="date">How the pipeline is performing</div>
      <h1>Funnel.</h1>
    </div>

    {#if loading}
      <p style="color:var(--mute)">Loading…</p>
    {:else if counts.applied === 0}
      <div class="empty-tab">
        <h3>No data to chart yet</h3>
        <p>Once you add an application or two, the funnel and conversion rates appear here.</p>
      </div>
    {:else}
      <!-- KPI ROW -->
      <div class="kpis">
        <div class="kpi tone-accent">
          <span class="ribbon"></span>
          <div class="kpi-lbl">Overall conversion</div>
          <div class="kpi-n">{overallPct ?? '—'}<span class="kpi-unit">%</span></div>
          <div class="kpi-sub">Applied → Offer</div>
        </div>
        <div class="kpi tone-positive">
          <span class="ribbon"></span>
          <div class="kpi-lbl">Best CV variant</div>
          <div class="kpi-n">{bestCv ? bestCv.name : '—'}</div>
          <div class="kpi-sub">{bestCv ? `${bestCv.rate}% reach screen` : 'Tag a CV variant to compare'}</div>
        </div>
        <div class="kpi tone-warm">
          <span class="ribbon"></span>
          <div class="kpi-lbl">Avg time to offer</div>
          <div class="kpi-n">{avgToOffer ?? '—'}<span class="kpi-unit">d</span></div>
          <div class="kpi-sub">Applied through Offer</div>
        </div>
        <div class="kpi tone-mute">
          <span class="ribbon"></span>
          <div class="kpi-lbl">In flight</div>
          <div class="kpi-n">{inFlight}</div>
          <div class="kpi-sub">Active applications</div>
        </div>
      </div>

      <!-- FUNNEL — stepped bars, monochromatic blue -->
      <div class="block">
        <div class="block-hd">
          <h2>Conversion funnel</h2>
          <span class="ai-tag">all time · {counts.applied} apps</span>
        </div>
        <div class="funnel-bars">
          {#each stages as s, i}
            {@const width = stages[0].n ? (s.n / stages[0].n) * 100 : 0}
            {@const stagePct = i === 0 ? null : pct(s.n, stages[i-1].n)}
            <div class="fb-row">
              <div class="fb-lbl-col">
                <span class="fb-lbl">{s.label}</span>
              </div>
              <div class="fb-track">
                <div class="fb-fill" style={`width: ${width}%; background: oklch(${0.78 - i*0.1} 0.16 258);`}></div>
                <span class="fb-n">{s.n}</span>
              </div>
              <div class="fb-pct">
                {#if stagePct === null}
                  <span class="fb-base">baseline</span>
                {:else}
                  <span class="fb-rate-pill">{stagePct}%</span>
                {/if}
              </div>
            </div>
          {/each}
        </div>
      </div>

      <!-- SOURCE + CV -->
      <div class="two-col">
        <div class="block">
          <div class="block-hd">
            <h2>Conversion by source</h2>
            <span class="ai-tag">Applied → Offer</span>
          </div>
          {#if sourceBars.length === 0}
            <p class="block-empty">Tag a source on your applications to see breakdowns.</p>
          {:else}
            <div class="bar-list">
              {#each sourceBars as s}
                <div class="bar-row">
                  <div class="bar-lbl">
                    <span>{s.name}</span>
                    <span class="bar-n">({s.n})</span>
                  </div>
                  <div class="bar-track">
                    <div class={`bar-fill f-${s.tone}`} style="width: {(s.rate / maxSource) * 100}%"></div>
                  </div>
                  <div class="bar-pct">{s.rate}%</div>
                </div>
              {/each}
            </div>
          {/if}
        </div>

        <div class="block">
          <div class="block-hd">
            <h2>Conversion by CV variant</h2>
            <span class="ai-tag">Applied → Screen</span>
          </div>
          {#if cvBars.length === 0}
            <p class="block-empty">Tag a CV variant on your applications to compare.</p>
          {:else}
            <div class="bar-list">
              {#each cvBars as c}
                <div class="bar-row">
                  <div class="bar-lbl">
                    <span>{c.name}</span>
                    <span class="bar-n">({c.n})</span>
                  </div>
                  <div class="bar-track">
                    <div class={`bar-fill f-${c.tone}`} style="width: {(c.rate / maxCv) * 100}%"></div>
                  </div>
                  <div class="bar-pct">{c.rate}%</div>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      </div>

      <!-- TIME IN STAGE -->
      <div class="block">
        <div class="block-hd">
          <h2>Average days in each stage</h2>
          <span class="ai-tag">how long things sit</span>
        </div>
        <div class="tis-row">
          {#each timeInStage as t}
            <div class={`tis-cell t-${t.tone}`}>
              <div class="tis-stage">{t.stage}</div>
              <div class="tis-n">{t.days ?? '—'}<span class="tis-unit">d</span></div>
              <div class="tis-bar"><div class="tis-fill" style="width: {((t.days ?? 0)/maxTime)*100}%"></div></div>
            </div>
          {/each}
        </div>
      </div>
    {/if}
  </div>
</div>

<style>
  .body { padding: 28px; }
  .body-inner { max-width: 1080px; margin: 0 auto; }

  .hello { margin-bottom: 24px; }
  .hello .date { font-size: 13.5px; color: var(--mute); margin-bottom: 6px; }
  .hello h1 { font-size: 30px; font-weight: 600; letter-spacing: -0.025em; margin: 0; }

  /* KPIs */
  .kpis { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; margin-bottom: 24px; }
  .kpi { position: relative; background: var(--card); border: 1px solid var(--rule); border-radius: 14px; padding: 18px 20px; overflow: hidden; box-shadow: var(--sh-1); }
  .ribbon { position: absolute; top: 0; left: 0; right: 0; height: 3px; }
  .kpi.tone-accent   .ribbon { background: var(--accent); }
  .kpi.tone-positive .ribbon { background: var(--positive); }
  .kpi.tone-warm     .ribbon { background: var(--warm); }
  .kpi.tone-mute     .ribbon { background: var(--rule-strong); }
  .kpi-lbl { font-size: 12.5px; color: var(--mute); margin-bottom: 4px; font-weight: 500; }
  .kpi-n   { font-size: 30px; font-weight: 600; letter-spacing: -0.035em; line-height: 1.1; font-feature-settings: "tnum"; }
  .kpi.tone-accent   .kpi-n { color: var(--accent-text); }
  .kpi.tone-positive .kpi-n { color: var(--positive-text); }
  .kpi.tone-warm     .kpi-n { color: var(--warm-text); }
  .kpi-unit { font-size: 16px; margin-left: 2px; opacity: 0.6; }
  .kpi-sub { font-size: 11.5px; color: var(--mute); margin-top: 8px; padding-top: 8px; border-top: 1px dashed var(--rule); }

  /* Block */
  .block { background: var(--card); border: 1px solid var(--rule); border-radius: 14px; padding: 22px 24px; margin-bottom: 14px; box-shadow: var(--sh-1); }
  .block-hd { display: flex; align-items: center; gap: 10px; margin-bottom: 18px; }
  .block-hd h2 { font-size: 17px; font-weight: 600; margin: 0; letter-spacing: -0.015em; }
  .block-empty { font-size: 13px; color: var(--mute); margin: 0; }
  .ai-tag { display: inline-flex; align-items: center; gap: 5px; font-size: 12px; background: var(--accent-tint); color: var(--accent-text); padding: 3px 10px; border-radius: 99px; font-weight: 500; }

  /* Funnel bars — monochromatic blue */
  .funnel-bars { display: flex; flex-direction: column; gap: 14px; }
  .fb-row { display: grid; grid-template-columns: 130px 1fr 140px; gap: 16px; align-items: center; }
  .fb-lbl-col { display: flex; align-items: center; gap: 8px; }
  .fb-lbl { font-size: 13.5px; font-weight: 600; }
  .fb-track { position: relative; height: 36px; background: var(--surface-2); border-radius: 10px; overflow: hidden; }
  .fb-fill { height: 100%; border-radius: 10px; transition: width 240ms ease; }
  .fb-n {
    position: absolute; top: 50%; left: 14px; transform: translateY(-50%);
    font-size: 14px; font-weight: 700; color: white;
    font-feature-settings: "tnum"; letter-spacing: -0.01em;
  }
  .fb-pct { display: flex; align-items: baseline; gap: 6px; }
  .fb-rate-pill {
    font-size: 14px; font-weight: 600; padding: 5px 12px; border-radius: 99px;
    background: var(--accent-tint); color: var(--accent-text);
    font-feature-settings: "tnum";
  }
  .fb-base { font-size: 12.5px; color: var(--mute); font-style: italic; }

  /* Two-column */
  .two-col { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; margin-bottom: 14px; }

  /* Generic bar list */
  .bar-list { display: flex; flex-direction: column; gap: 12px; }
  .bar-row { display: grid; grid-template-columns: 1fr 1fr 50px; gap: 12px; align-items: center; }
  .bar-lbl { display: flex; align-items: center; gap: 8px; font-size: 13px; font-weight: 500; }
  .bar-n { font-size: 11.5px; color: var(--mute); font-feature-settings: "tnum"; }
  .bar-track { height: 14px; background: var(--surface-2); border-radius: 999px; overflow: hidden; }
  .bar-fill { height: 100%; border-radius: 999px; }
  .bar-fill.f-positive { background: var(--positive); }
  .bar-fill.f-accent   { background: var(--accent); }
  .bar-fill.f-warm     { background: var(--warm); }
  .bar-fill.f-mute     { background: var(--mute-2); }
  .bar-pct { font-size: 14px; font-weight: 600; text-align: right; font-feature-settings: "tnum"; color: var(--ink); }

  /* Time-in-stage */
  .tis-row { display: grid; grid-template-columns: repeat(4, 1fr); gap: 10px; }
  .tis-cell { background: var(--surface-2); border-radius: 10px; padding: 14px 16px; }
  .tis-stage { font-size: 12px; color: var(--mute); font-weight: 500; }
  .tis-n { font-size: 24px; font-weight: 600; letter-spacing: -0.025em; margin: 4px 0 8px; font-feature-settings: "tnum"; }
  .tis-cell.t-positive .tis-n { color: var(--positive-text); }
  .tis-cell.t-accent   .tis-n { color: var(--accent-text); }
  .tis-cell.t-warm     .tis-n { color: var(--warm-text); }
  .tis-unit { font-size: 14px; opacity: 0.6; margin-left: 1px; }
  .tis-bar { height: 6px; background: var(--card); border-radius: 99px; overflow: hidden; }
  .tis-fill { height: 100%; background: var(--mute-2); border-radius: 99px; }
  .tis-cell.t-positive .tis-fill { background: var(--positive); }
  .tis-cell.t-accent   .tis-fill { background: var(--accent); }
  .tis-cell.t-warm     .tis-fill { background: var(--warm); }

  .empty-tab {
    border: 1px dashed var(--rule);
    border-radius: 12px;
    padding: 32px;
    text-align: center;
    background: var(--card);
  }
  .empty-tab h3 { margin: 0 0 .5rem; font-size: 16px; font-weight: 500; color: var(--ink); }
  .empty-tab p { color: var(--mute); margin: 0; font-size: 13.5px; }

  /* Mobile — stack KPIs and breakdowns; keep funnel bars readable. */
  @media (max-width: 720px) {
    .body { padding: 18px 14px; }
    .hello h1 { font-size: 22px; }
    .kpis { grid-template-columns: 1fr 1fr; gap: 8px; }
    .kpi { padding: 14px 16px; }
    .kpi-n { font-size: 24px; }
    .two-col { grid-template-columns: 1fr; gap: 10px; }
    .block { padding: 16px 16px; }
    .block-hd { flex-wrap: wrap; }
    .fb-row { grid-template-columns: 90px 1fr 70px; gap: 10px; }
    .fb-track { height: 30px; }
    .fb-rate-pill { font-size: 13px; padding: 4px 10px; }
    .tis-row { grid-template-columns: 1fr 1fr; gap: 8px; }
    .bar-row { grid-template-columns: 1fr 1fr 40px; gap: 8px; }
  }
</style>
