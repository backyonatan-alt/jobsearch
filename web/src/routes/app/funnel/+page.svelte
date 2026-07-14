<script>
  import { onMount } from 'svelte';
  import { api } from '$lib/api.js';
  import { isPreview, mockApi } from '$lib/preview-mode.js';
  import { toDisplayApp, daysSince, fmtShortDate } from '$lib/app-helpers.js';
  import CompanyLogo from '$lib/CompanyLogo.svelte';

  const call = isPreview() ? mockApi : api;

  let apps = $state([]);
  let loading = $state(true);

  async function load() {
    try { apps = (await call('/api/applications')).map(toDisplayApp); }
    catch (e) { if (e.message !== 'unauthorized') console.error(e); }
    finally { loading = false; }
  }
  onMount(() => {
    load();
    const h = () => load();
    window.addEventListener('pursuit:refresh', h);
    return () => window.removeEventListener('pursuit:refresh', h);
  });

  // ── Funnel counts (cumulative reach — an app at interview passed screen) ──
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

  const CHART_H = 180;
  const stages = $derived.by(() => {
    const c = counts;
    const defs = [
      { key: 'applied',   label: 'Applied',   n: c.applied,   sub: 'applications' },
      { key: 'screen',    label: 'Screen',    n: c.screen },
      { key: 'interview', label: 'Interview', n: c.interview },
      { key: 'offer',     label: 'Offer',     n: c.offer }
    ];
    const max = Math.max(1, c.applied);
    return defs.map((d, i) => {
      const h = Math.max(d.n > 0 ? 12 : 4, Math.round((d.n / max) * CHART_H));
      const pct = Math.round((d.n / max) * 100);
      return { ...d, h, sub: d.sub || `${pct}% of all`, cluster: apps.filter(a => a.status === d.key).slice(0, 3), clusterN: apps.filter(a => a.status === d.key).length };
    });
  });
  // Connector between stage i and i+1: trapezoid sloping from the top of bar i
  // down to the top of bar i+1, drawn inside a container as tall as bar i.
  const drops = $derived.by(() => {
    const out = [];
    for (let i = 0; i < stages.length - 1; i++) {
      const a = stages[i], b = stages[i + 1];
      const dropPct = a.n > 0 ? Math.round((1 - b.n / a.n) * 100) : 0;
      const topPct = a.h > 0 ? Math.round((1 - b.h / a.h) * 100) : 0;
      out.push({ h: a.h, topPct, dropPct, bad: dropPct >= 60 });
    }
    return out;
  });
  const overallPct = $derived(counts.applied ? Math.round((counts.offer / counts.applied) * 100) : 0);

  // ── Headline stats ──────────────────────────────────────────
  const replyRate = $derived(counts.applied ? Math.round((counts.screen / counts.applied) * 100) : null);
  const avgReplyDays = $derived.by(() => {
    const replied = apps.filter(a => ['screen', 'interview', 'offer'].includes(a.status) && a.raw.applied_at);
    if (!replied.length) return null;
    return Math.round(replied.reduce((s, a) => s + (daysSince(a.raw.applied_at) ?? 0), 0) / replied.length);
  });
  const furthest = $derived.by(() => {
    for (const s of ['offer', 'interview', 'screen', 'applied']) {
      const m = apps.filter(a => a.status === s);
      if (m.length) return { label: s.charAt(0).toUpperCase() + s.slice(1), app: m[0], green: s === 'offer' };
    }
    return null;
  });

  // ── Exits ───────────────────────────────────────────────────
  const exits = $derived.by(() => {
    let rejected = 0, withdrawn = 0, closed = 0;
    for (const a of apps) {
      if (a.status === 'rejected') rejected++;
      else if (a.status === 'withdrawn') withdrawn++;
      else if (a.status === 'closed') closed++;
    }
    return { rejected, withdrawn, closed, any: rejected + withdrawn + closed > 0 };
  });

  // ── Sources with conversion ─────────────────────────────────
  const SRC_COLORS = { referral: '#16a34a', linkedin: '#2463eb' };
  const sources = $derived.by(() => {
    const inPlay = apps.filter(a => a.status !== 'wishlist');
    const map = {};
    for (const a of inPlay) {
      const src = (a.source && a.source !== '—') ? a.source : 'Direct / other';
      (map[src] ||= []).push(a);
    }
    const maxN = Math.max(1, ...Object.values(map).map(v => v.length));
    return Object.entries(map).map(([src, list]) => {
      const past = { screen: 0, interview: 0, offer: 0 };
      for (const a of list) if (past[a.status] !== undefined) past[a.status]++;
      const conv = [];
      if (past.screen) conv.push(`${past.screen} screen`);
      if (past.interview) conv.push(`${past.interview} interview`);
      if (past.offer) conv.push(`${past.offer} offer`);
      const key = src.toLowerCase();
      const color = SRC_COLORS[Object.keys(SRC_COLORS).find(k => key.includes(k))] || '#8fabef';
      return { src, n: list.length, conv: conv.join(', '), w: Math.round((list.length / maxN) * 100), color };
    }).sort((a, b) => b.n - a.n);
  });
  const sourceInsight = $derived.by(() => {
    const ref = sources.find(s => s.src.toLowerCase().includes('referral'));
    if (ref?.conv && sources.length > 1) return { lead: 'Referrals convert best for you', rest: ` — ${ref.n} sent, ${ref.conv}.` };
    if (sources[0]) return { lead: `${sources[0].src} carries your volume`, rest: ` — ${sources[0].n} of ${apps.filter(a => a.status !== 'wishlist').length} applications.` };
    return null;
  });

  // ── Noticing — live observations from the same heuristics as Home ──
  const noticed = $derived.by(() => {
    const out = [];
    const inPlay = apps.filter(a => !['wishlist', 'rejected', 'withdrawn', 'closed'].includes(a.status));
    const quiet = apps.filter(a => a.stale);
    const referral = inPlay.filter(a => (a.source || '').toLowerCase().includes('referral'));
    const refProg = referral.filter(a => ['screen', 'interview', 'offer'].includes(a.status));
    if (refProg.length) out.push({ lead: 'Referrals are working.', rest: `${refProg.length} of ${referral.length} reached a screen or further — ask for one more intro this week.` });
    if (quiet.length >= 2) out.push({ lead: `${quiet[0].co} and ${quiet[1].co} have gone quiet.`, rest: 'No reply in over a week — a direct nudge beats waiting.' });
    else if (quiet.length === 1) out.push({ lead: `${quiet[0].co} has gone quiet.`, rest: 'No reply in over a week — a direct nudge beats waiting.' });
    if (counts.screen > 0 && counts.interview > 0) {
      const r = Math.round((counts.interview / counts.screen) * 100);
      if (r >= 50) out.push({ lead: `Your screen → interview rate is strong (${r}%).`, rest: 'The bottleneck is getting replies, not passing rounds.' });
    }
    return out.slice(0, 4);
  });

  // Date range line — from the earliest applied_at to today.
  const rangeLine = $derived.by(() => {
    const dates = apps.map(a => a.raw.applied_at).filter(Boolean).map(d => new Date(d));
    if (!dates.length) return '';
    const min = new Date(Math.min(...dates));
    return `${fmtShortDate(min)} – ${fmtShortDate(new Date())}`;
  });
</script>

<svelte:head><title>Insights — Pursuit</title></svelte:head>

<div class="pg">
  <div class="head">
    <div>
      <h1>Insights.</h1>
      <div class="sub">How your search is actually going{rangeLine ? ` · ${rangeLine}` : ''}</div>
    </div>
  </div>

  {#if loading}
    <p class="loading">Loading…</p>
  {:else if counts.applied === 0}
    <div class="empty">
      <h3>No data to chart yet</h3>
      <p>Once you add an application or two, the funnel and sources appear here.</p>
    </div>
  {:else}
    <!-- headline numbers -->
    <div class="stats">
      <div class="st">
        <div class="st-v">{replyRate ?? '—'}{#if replyRate !== null}%{/if}</div>
        <div class="st-l">reply rate</div>
        <div class="st-s">{counts.screen} of {counts.applied} got a reply</div>
      </div>
      <div class="st">
        <div class="st-v">{avgReplyDays !== null ? `${avgReplyDays}d` : '—'}</div>
        <div class="st-l">to first reply</div>
        <div class="st-s">proxy: days since applied, replied apps</div>
      </div>
      <div class="st">
        <div class="st-v" style={furthest?.green ? 'color:#16a34a' : ''}>{furthest?.label ?? '—'}</div>
        <div class="st-l">furthest stage</div>
        {#if furthest?.app}
          <div class="st-s chip-line"><CompanyLogo app={furthest.app} size={20} radius={6} />{furthest.app.co}{furthest.app.raw.salary_note ? ` · ${furthest.app.raw.salary_note}` : ''}</div>
        {:else}
          <div class="st-s">keep applying</div>
        {/if}
      </div>
    </div>

    <!-- funnel -->
    <div class="card" data-tour="funnel">
      <div class="card-hd">
        <span class="ch-t">Pipeline funnel</span>
        <span class="ch-s">every application, by furthest stage reached</span>
        <span class="ch-ov"><b>{overallPct}%</b> overall · applied → offer</span>
      </div>

      <div class="chart">
        <div class="grid" style="height:{CHART_H}px">
          {#each stages as s, i (s.key)}
            <div class="bar-col">
              <div class="bar-n">{s.n} <span class="bar-sub">{s.sub}</span></div>
              <div class="bar" style="height:{s.h}px"></div>
            </div>
            {#if i < stages.length - 1}
              <div class="conn" style="height:{drops[i].h}px">
                <div class="slope" style="clip-path:polygon(0 0, 100% {drops[i].topPct}%, 100% 100%, 0 100%)"></div>
                <span class="drop" class:ok={!drops[i].bad} style="top:{Math.min(drops[i].h - 24, Math.max(-14, Math.round(drops[i].h * drops[i].topPct / 100 / 2)))}px">↓ {drops[i].dropPct}%</span>
              </div>
            {/if}
          {/each}
        </div>
        <div class="labels">
          {#each stages as s, i (s.key)}
            <div class="lab">
              <div class="lab-t">{s.label}</div>
              <div class="lab-logos">
                {#each s.cluster as a, j (a.id)}
                  <span class="lab-chip" style="margin-left:{j ? '-6px' : '0'}"><CompanyLogo app={a} size={20} radius={6} /></span>
                {/each}
                {#if s.clusterN > 3}<span class="lab-extra">+{s.clusterN - 3}</span>{/if}
              </div>
            </div>
            {#if i < stages.length - 1}<div></div>{/if}
          {/each}
        </div>
      </div>

      {#if counts.applied > counts.screen}
        <div class="card-ft">
          <span>The applied → screen step is where you lose the most — <strong>{Math.round((1 - counts.screen / counts.applied) * 100)}% never reply</strong>.</span>
          <a href="/app/applications">See the quiet ones →</a>
        </div>
      {/if}
    </div>

    <!-- exits strip -->
    {#if exits.any}
      <div class="exits">
        <span><strong class="red">{exits.rejected}</strong> rejected</span>
        <span><strong>{exits.closed}</strong> position closed</span>
        <span><strong>{exits.withdrawn}</strong> withdrawn</span>
        <span class="ex-note">exits live in <a href="/app/applications">Applications → No longer in play</a></span>
      </div>
    {/if}

    <!-- sources + noticing -->
    <div class="two">
      <div>
        <div class="sec-t">Where they come from</div>
        <div class="srcs">
          {#each sources as s (s.src)}
            <div class="src">
              <div class="src-hd"><span class="src-n">{s.src}</span><span class="src-m">{s.n}{s.conv ? ` · ${s.conv}` : ''}</span></div>
              <div class="track"><i style="width:{s.w}%;background:{s.color}"></i></div>
            </div>
          {/each}
        </div>
        {#if sourceInsight}
          <div class="src-insight"><strong>{sourceInsight.lead}</strong>{sourceInsight.rest}</div>
        {/if}
      </div>
      <div>
        <div class="sec-t orange">✦ What we've noticed <span class="sec-s">computed live from your pipeline · the freshest one shows on Home</span></div>
        <div class="notes">
          {#if noticed.length === 0}
            <div class="n-empty">Nothing notable yet — observations appear as your pipeline moves.</div>
          {:else}
            {#each noticed as n, i (i)}
              <div class="note"><strong>{n.lead}</strong> {n.rest}</div>
            {/each}
          {/if}
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .pg { max-width: 1000px; margin: 0 auto; padding: 36px 32px 80px; width: 100%; box-sizing: border-box; }
  .head { display: flex; align-items: baseline; justify-content: space-between; margin-bottom: 26px; }
  .head h1 { font-size: 30px; font-weight: 700; letter-spacing: -0.02em; margin: 0 0 6px; }
  .sub { font-size: 13.5px; color: #6f7680; }
  .loading { color: #8a9099; font-size: 13.5px; }

  .stats { display: grid; grid-template-columns: 1fr 1fr 1fr; border-top: 1px solid #e2e2de; margin-bottom: 44px; }
  .st { padding: 24px 26px 0; border-right: 1px solid #eeeeea; }
  .st:last-child { border-right: 0; }
  .st-v { font-size: 40px; font-weight: 700; letter-spacing: -0.03em; line-height: 1; }
  .st-l { font-size: 13px; font-weight: 600; margin: 9px 0 4px; }
  .st-s { font-size: 12px; color: #8a9099; }
  .chip-line { display: flex; align-items: center; gap: 7px; }

  .card { background: #fff; border: 1px solid #e8e8e5; border-radius: 16px; padding: 26px 30px; margin-bottom: 20px; box-shadow: 0 1px 3px rgba(22,24,28,.04); }
  .card-hd { display: flex; align-items: baseline; gap: 10px; margin-bottom: 26px; flex-wrap: wrap; }
  .ch-t { font-size: 16px; font-weight: 700; }
  .ch-s { font-size: 12.5px; color: #8a9099; }
  .ch-ov { margin-left: auto; font-size: 12.5px; font-weight: 600; color: #6f7680; }
  .ch-ov b { font-size: 24px; font-weight: 700; letter-spacing: -0.02em; color: #16a34a; margin-right: 4px; }

  .chart { position: relative; }
  .grid {
    display: grid; grid-template-columns: 1fr 56px 1fr 56px 1fr 56px 1fr;
    align-items: end; position: relative;
    background-image: repeating-linear-gradient(to top, transparent 0, transparent calc(50% - 1px), #f0f0ed calc(50% - 1px), #f0f0ed 50%);
  }
  .bar-col { display: flex; flex-direction: column; justify-content: flex-end; }
  .bar-n { font-size: 24px; font-weight: 700; letter-spacing: -0.02em; margin-bottom: 6px; }
  .bar-sub { font-size: 12px; font-weight: 600; color: #8a9099; }
  .bar { background: linear-gradient(180deg, #2463eb, #4d7bee); border-radius: 9px 9px 0 0; }
  .conn { position: relative; align-self: end; }
  .slope { position: absolute; inset: 0; background: #eef4ff; }
  .drop {
    position: absolute; left: 50%; transform: translateX(-50%);
    font-size: 11px; font-weight: 700; color: #b3372a; background: #fff;
    border: 1px solid #f2d4cf; border-radius: 5px; padding: 2px 6px; white-space: nowrap;
    box-shadow: 0 1px 2px rgba(22,24,28,.06);
  }
  .drop.ok { color: #1d7a4f; border-color: #cfe5d2; }
  .labels { display: grid; grid-template-columns: 1fr 56px 1fr 56px 1fr 56px 1fr; border-top: 2px solid #e2e2de; padding-top: 10px; }
  .lab-t { font-size: 14.5px; font-weight: 700; margin-bottom: 8px; }
  .lab-logos { display: flex; align-items: center; min-height: 20px; }
  .lab-chip { display: inline-flex; border-radius: 7px; box-shadow: 0 0 0 2px #fff; }
  .lab-extra { font-size: 11px; font-weight: 600; color: #8a9099; margin-left: 5px; }
  .card-ft { display: flex; align-items: center; gap: 16px; margin-top: 16px; font-size: 12.5px; color: #6f7680; }
  .card-ft strong { color: #16181c; }
  .card-ft a { margin-left: auto; flex: none; color: #2463eb; text-decoration: none; }

  .exits { display: flex; align-items: center; gap: 20px; border: 1px solid #e8e8e5; background: #fff; border-radius: 12px; padding: 12px 20px; margin-bottom: 40px; font-size: 13px; color: #4b5158; flex-wrap: wrap; }
  .exits .red { color: #b3372a; }
  .ex-note { color: #8a9099; margin-left: auto; }
  .ex-note a { color: #2463eb; text-decoration: none; }

  .two { display: grid; grid-template-columns: 1fr 1fr; gap: 40px; align-items: start; }
  .sec-t { font-size: 11px; font-weight: 600; letter-spacing: .12em; text-transform: uppercase; color: #8a9099; margin-bottom: 14px; }
  .sec-t.orange { color: #e0641f; }
  .sec-s { font-size: 12px; color: #8a9099; letter-spacing: 0; text-transform: none; font-weight: 400; margin-left: 8px; }
  .srcs { display: flex; flex-direction: column; gap: 12px; font-size: 13px; }
  .src-hd { display: flex; align-items: baseline; margin-bottom: 5px; }
  .src-n { font-weight: 600; }
  .src-m { margin-left: auto; color: #8a9099; }
  .track { height: 6px; background: #f0f0ed; border-radius: 3px; }
  .track i { display: block; height: 6px; border-radius: 3px; }
  .src-insight { font-size: 12.5px; color: #6f7680; margin-top: 14px; line-height: 1.55; }
  .src-insight strong { color: #16181c; }
  .notes { display: flex; flex-direction: column; gap: 12px; font-size: 13px; color: #4b5158; }
  .note { line-height: 1.55; }
  .note strong { color: #16181c; }
  .n-empty { font-size: 12.5px; color: #8a9099; }

  .empty { border: 1px dashed #e2e2de; border-radius: 12px; padding: 32px; text-align: center; background: #fff; }
  .empty h3 { margin: 0 0 .5rem; font-size: 16px; font-weight: 500; }
  .empty p { color: #8a9099; margin: 0; font-size: 13.5px; }

  @media (max-width: 900px) {
    .pg { padding: 24px 16px 60px; }
    .stats { grid-template-columns: 1fr; }
    .st { border-right: 0; border-bottom: 1px solid #eeeeea; padding-bottom: 18px; }
    .two { grid-template-columns: 1fr; }
  }
</style>
