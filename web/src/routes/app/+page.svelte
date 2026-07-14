<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api.js';
  import { STATUS_LABEL, toDisplayApp, daysSince, fmtRelativeDate } from '$lib/app-helpers.js';
  import { logEvent } from '$lib/analytics.js';
  import CompanyLogo from '$lib/CompanyLogo.svelte';

  const call = api;

  let apps = $state([]);
  let me = $state(null);
  let loading = $state(true);

  // Per-app interview events + debriefs, for the prep rows + debrief banner.
  let eventsByApp = $state({});
  let debriefsByApp = $state({});

  onMount(() => {
    refresh();
    const h = () => refresh();
    const onVis = () => { if (document.visibilityState === 'visible') refresh(); };
    window.addEventListener('pursuit:refresh', h);
    window.addEventListener('focus', h);
    document.addEventListener('visibilitychange', onVis);
    return () => {
      window.removeEventListener('pursuit:refresh', h);
      window.removeEventListener('focus', h);
      document.removeEventListener('visibilitychange', onVis);
    };
  });

  async function refresh() {
    try {
      const [meRes, raw] = await Promise.all([call('/api/me'), call('/api/applications')]);
      me = meRes;
      apps = raw.map(toDisplayApp);
      loading = false;
      await hydrate(raw);
    } catch (e) {
      if (e.message !== 'unauthorized') console.error(e);
      loading = false;
    }
  }

  async function hydrate(raw) {
    const live = raw.filter(a => ['screen', 'interview', 'offer'].includes(a.status));
    await Promise.all(live.map(async (a) => {
      const [ivs, dbs] = await Promise.all([
        call(`/api/applications/${a.id}/interviews`).catch(() => []),
        call(`/api/applications/${a.id}/debriefs`).catch(() => [])
      ]);
      if (Array.isArray(ivs) && ivs.length) eventsByApp[a.id] = ivs;
      if (Array.isArray(dbs) && dbs.length) debriefsByApp[a.id] = dbs;
    }));
    eventsByApp = { ...eventsByApp };
    debriefsByApp = { ...debriefsByApp };
  }

  // ── Upcoming interview per app ──────────────────────────────
  const upcomingByApp = $derived.by(() => {
    const now = Date.now();
    const out = {};
    for (const a of apps) {
      for (const ev of (eventsByApp[a.id] || [])) {
        if (!ev.starts_at) continue;
        const t = new Date(ev.starts_at).getTime();
        if (isNaN(t) || t < now) continue;
        if (!out[a.id] || t < new Date(out[a.id].starts_at).getTime()) out[a.id] = ev;
      }
    }
    return out;
  });

  function fmtWhen(iso) {
    const d = new Date(iso);
    const now = new Date();
    const time = d.toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit' });
    const days = Math.floor((d - new Date(now.getFullYear(), now.getMonth(), now.getDate())) / 86400000);
    if (days === 0) return `today ${time}`;
    if (days === 1) return `tomorrow ${time}`;
    if (days < 7) return `${d.toLocaleDateString(undefined, { weekday: 'short' })} ${time}`;
    return d.toLocaleDateString(undefined, { day: 'numeric', month: 'short' });
  }

  // ── Debrief loop (Phase 3b) — unchanged semantics ───────────
  const pendingDebrief = $derived.by(() => {
    const now = Date.now();
    const due = [];
    for (const a of apps) {
      if (!['screen', 'interview', 'offer'].includes(a.status)) continue;
      const done = new Set((debriefsByApp[a.id] || []).map(d => d.interview_id));
      for (const iv of (eventsByApp[a.id] || [])) {
        if (done.has(iv.id)) continue;
        const past = iv.scheduled === false || !iv.starts_at ||
          new Date(iv.starts_at).getTime() < now;
        if (past) due.push({ app: a, iv });
      }
    }
    if (!due.length) return null;
    due.sort((a, b) => new Date(b.iv.starts_at || 0) - new Date(a.iv.starts_at || 0));
    return due[0];
  });
  const debriefViewLogged = new Set();
  $effect(() => {
    const p = pendingDebrief;
    if (p && !debriefViewLogged.has(p.iv.id)) {
      debriefViewLogged.add(p.iv.id);
      logEvent('debrief_prompt_view', { app_id: Number(p.app.id), surface: 'today' });
    }
  });

  // ── The five pulse categories ───────────────────────────────
  const EXITS = ['rejected', 'withdrawn', 'closed'];
  const toPrep   = $derived(apps.filter(a => upcomingByApp[a.id])
    .sort((a, b) => new Date(upcomingByApp[a.id].starts_at) - new Date(upcomingByApp[b.id].starts_at)));
  const toDecide = $derived(apps.filter(a => a.status === 'offer' && !upcomingByApp[a.id]));
  const toNudge  = $derived(apps.filter(a => a.stale && !upcomingByApp[a.id] && a.status !== 'offer')
    .sort((a, b) => (daysSince(b.raw.applied_at) ?? 0) - (daysSince(a.raw.applied_at) ?? 0)));
  const waiting  = $derived(apps.filter(a =>
    !EXITS.includes(a.status) && !upcomingByApp[a.id] && a.status !== 'offer' && !a.stale
  ).sort((a, b) => (a.status === 'wishlist') - (b.status === 'wishlist')));
  const exits    = $derived(apps.filter(a => EXITS.includes(a.status)));

  const cats = $derived([
    { count: toPrep.length,   label: 'to prep',   color: '#e0641f', labelColor: '#16181c', note: toPrep[0] ? fmtWhen(upcomingByApp[toPrep[0].id].starts_at) : 'no interviews booked', group: toPrep },
    { count: toDecide.length, label: 'to decide', color: '#16a34a', labelColor: '#16181c', note: toDecide[0]?.raw.salary_note || (toDecide.length ? 'offer on the table' : 'no offers yet'), group: toDecide },
    { count: toNudge.length,  label: 'to nudge',  color: '#b3372a', labelColor: '#16181c', note: toNudge.length ? 'quiet over a week' : 'nothing stalled', group: toNudge },
    { count: waiting.length,  label: 'waiting',   color: '#16181c', labelColor: '#6f7680', note: 'nothing to do yet', group: waiting },
    { count: exits.length,    label: 'closed',    color: '#9aa1ab', labelColor: '#6f7680', note: 'kept for the record', group: exits, gray: true }
  ]);

  // ── One-line rows, in priority order ────────────────────────
  const rows = $derived.by(() => {
    const out = [];
    for (const a of toPrep) {
      const ev = upcomingByApp[a.id];
      out.push({ app: a, kind: 'act', meta: `${a.role} · ${(ev.summary || 'interview').trim()}`,
        hot: fmtWhen(ev.starts_at), hotColor: '#c05310', border: '#cdddfb',
        btn: 'Prep now →', btnBg: '#2463eb', btnColor: '#fff', btnBorder: '#2463eb',
        go: () => goto(`/app/${a.id}#interview-prep`) });
    }
    for (const a of toDecide) {
      out.push({ app: a, kind: 'act', meta: `${a.role} · offer`,
        hot: a.raw.salary_note || 'awaiting your decision', hotColor: '#1d7a4f', border: '#cfe5d2',
        btn: 'Decide →', btnBg: '#16a34a', btnColor: '#fff', btnBorder: '#16a34a',
        go: () => goto(`/app/${a.id}`) });
    }
    for (const a of toNudge) {
      out.push({ app: a, kind: 'act', meta: `${a.role} ·`,
        hot: `quiet ${daysSince(a.raw.applied_at) ?? 0} days`, hotColor: '#b3372a', border: '#f2d4cf',
        btn: 'Follow up', btnBg: '#fff', btnColor: '#2463eb', btnBorder: '#cdddfb',
        go: () => goto(`/app/${a.id}`) });
    }
    for (const a of waiting) {
      const wish = a.status === 'wishlist';
      out.push({ app: a, kind: 'quiet',
        meta: wish ? `${a.role} · saved ${fmtRelativeDate(a.raw.created_at)}`
          : `${a.role} ·${a.status !== 'applied' ? ` ${STATUS_LABEL[a.status].toLowerCase()} ·` : ''} applied ${a.appliedRel}${a.source && a.source !== '—' ? ` via ${a.source}` : ''}`,
        quiet: wish ? 'apply when ready' : ((daysSince(a.raw.applied_at) ?? 99) < 3 ? 'too early' : 'waiting'),
        go: () => goto(`/app/${a.id}`) });
    }
    for (const a of exits) {
      const label = a.status === 'closed' ? 'position closed' : a.status;
      out.push({ app: a, kind: 'exit',
        meta: `${a.role} · ${label} · ${daysSince(a.raw.updated_at || a.raw.applied_at) ?? 0}d`,
        go: () => goto(`/app/${a.id}`) });
    }
    return out;
  });

  async function reopen(e, a) {
    e.stopPropagation();
    await api(`/api/applications/${a.id}`, { method: 'PATCH', body: JSON.stringify({ status: 'applied' }) });
    logEvent('archive_reopened', { app_id: Number(a.id), from: a.status });
    try { window.dispatchEvent(new CustomEvent('pursuit:refresh')); } catch {}
    refresh();
  }

  // ── Noticing card — one honest observation from the data ────
  const noticing = $derived.by(() => {
    if (!apps.length) return null;
    const inPlay = apps.filter(a => !['wishlist', ...EXITS].includes(a.status));
    const referral = inPlay.filter(a => (a.source || '').toLowerCase().includes('referral'));
    const refProgressed = referral.filter(a => ['screen', 'interview', 'offer'].includes(a.status));
    if (referral.length >= 1 && refProgressed.length >= 1) {
      return { lead: 'Referrals are working:', rest: `${refProgressed.length} of your ${referral.length} referral ${referral.length === 1 ? 'application' : 'applications'} reached a screen or further. Worth asking for one more intro this week.` };
    }
    if (toNudge.length >= 2) {
      return { lead: `${toNudge[0].co} and ${toNudge[1].co} have gone quiet.`, rest: 'No reply in over a week — a direct nudge beats waiting.' };
    }
    if (toNudge.length === 1) {
      return { lead: `${toNudge[0].co} has gone quiet.`, rest: 'No reply in over a week — a direct nudge beats waiting.' };
    }
    const screens = inPlay.filter(a => ['screen', 'interview', 'offer'].includes(a.status)).length;
    if (inPlay.length >= 3 && screens > 0) {
      return { lead: `${Math.round((screens / inPlay.length) * 100)}% of your applications got a reply.`, rest: 'The bottleneck is getting replies, not passing rounds — keep the volume up.' };
    }
    return { lead: `${inPlay.length} ${inPlay.length === 1 ? 'application' : 'applications'} in play.`, rest: 'Keep the pipeline moving — add a few more this week.' };
  });

  // ── Greeting ────────────────────────────────────────────────
  const firstName = $derived.by(() => {
    if (!me?.email) return 'there';
    const local = me.email.split('@')[0].split(/[._]/).pop() || 'there';
    return local.charAt(0).toUpperCase() + local.slice(1);
  });
  const greeting = $derived.by(() => {
    const h = new Date().getHours();
    if (h < 12) return 'Good morning';
    if (h < 18) return 'Good afternoon';
    return 'Good evening';
  });
  const dateLine = $derived.by(() => {
    const d = new Date();
    const dow = d.toLocaleDateString('en-US', { weekday: 'long' });
    const rest = d.toLocaleDateString('en-US', { day: 'numeric', month: 'long', year: 'numeric' });
    return `${dow} · ${rest}`;
  });

  function scrollToList() {
    document.getElementById('list')?.scrollIntoView({ behavior: 'smooth' });
  }
  function openDebrief(p) { goto(`/app/${p.app.id}?debrief=${p.iv.id}`); }
  function openNewApp() { try { window.dispatchEvent(new CustomEvent('pursuit:new-app')); } catch {} }
</script>

<svelte:head><title>Home — Pursuit</title></svelte:head>

{#if loading}
  <div class="pg"><p class="loading">Loading your day…</p></div>
{:else if apps.length === 0}
  <!-- Empty state — first run, nothing tracked yet -->
  <div class="empty">
    <div class="e-date">{dateLine}</div>
    <h1>{greeting}, {firstName}.</h1>
    <p class="e-sub">Track your first application and Pursuit starts working — reply nudges, interview prep, the whole picture in one place.</p>
    <div class="e-cta">
      <button class="e-drop" onclick={openNewApp}>
        <span class="e-d1">Paste a job URL, drop a screenshot, or describe the role</span>
        <span class="e-d2">We fill in the company, role and source for you · ⌘V works right here</span>
      </button>
      <div class="e-btns">
        <button class="e-new" onclick={openNewApp}>New application <span class="nk">⌘N</span></button>
      </div>
    </div>
    <div class="e-legend">
      <span>Then this page becomes:</span>
      <span class="e-li"><span class="dot" style="background:#e0641f"></span>what to prep</span>
      <span class="e-li"><span class="dot" style="background:#16a34a"></span>what to decide</span>
      <span class="e-li"><span class="dot" style="background:#b3372a"></span>who to nudge</span>
    </div>
  </div>
{:else}
  <!-- 3a header band -->
  <div class="band">
    <div class="band-in">
      <div class="b-date">{dateLine}</div>
      <h1>{greeting}, {firstName}.</h1>
      <div class="cells" data-tour="stats">
        {#each cats as cat (cat.label)}
          <button class="cell" onclick={scrollToList}>
            <span class="c-top"><span class="c-n" style="color:{cat.color}">{cat.count}</span><span class="c-l" style="color:{cat.labelColor}">{cat.label}</span></span>
            <span class="c-logos">
              {#each cat.group.slice(0, 3) as a, i (a.id)}
                <span class="c-chip" style="margin-left:{i ? '-7px' : '0'}"><CompanyLogo app={a} size={24} gray={cat.gray || false} /></span>
              {/each}
              {#if cat.group.length > 3}<span class="c-extra">+{cat.group.length - 3}</span>{/if}
            </span>
            <span class="c-note">{cat.note}</span>
          </button>
        {/each}
      </div>
    </div>
  </div>

  <div class="pg">
    {#if pendingDebrief}
      <button type="button" class="db-banner" onclick={() => openDebrief(pendingDebrief)}>
        <span class="db-spark">✦</span>
        <span class="db-tx">How did the <b>{(pendingDebrief.iv.summary || 'last').trim() || 'last'}</b> round at <b>{pendingDebrief.app.co}</b> go? A 20-second debrief sharpens your next round's prep.</span>
        <span class="db-cta">Debrief →</span>
      </button>
    {/if}

    {#if noticing}
      <div class="noticing" data-tour="prep">
        <span class="n-spark">✦</span>
        <span class="n-tx"><strong>{noticing.lead}</strong> {noticing.rest}</span>
        <a href="/app/funnel" class="n-more">More in Insights →</a>
      </div>
    {/if}

    <div id="list" class="list-hd">
      <span class="lh-t">Everything · {apps.length}</span>
      <span class="lh-s">one line each · a button means it needs you · exits at the end</span>
    </div>

    <div class="rows">
      {#each rows as r (r.app.id)}
        {#if r.kind === 'act'}
          <div class="row act" style="border-color:{r.border}" onclick={r.go} role="button" tabindex="0">
            <CompanyLogo app={r.app} size={26} />
            <span class="r-tx"><strong>{r.app.co}</strong> <span class="r-meta">· {r.meta}</span> <strong style="color:{r.hotColor}">{r.hot}</strong></span>
            <button class="r-btn" style="background:{r.btnBg};color:{r.btnColor};border-color:{r.btnBorder}" onclick={(e) => { e.stopPropagation(); r.go(); }}>{r.btn}</button>
          </div>
        {:else if r.kind === 'quiet'}
          <div class="row quiet" onclick={r.go} role="button" tabindex="0">
            <CompanyLogo app={r.app} size={22} />
            <span class="r-tx sm"><strong class="q-co">{r.app.co}</strong> <span class="r-meta">· {r.meta}</span></span>
            <span class="r-quiet">{r.quiet}</span>
          </div>
        {:else}
          <div class="row exit" onclick={r.go} role="button" tabindex="0">
            <CompanyLogo app={r.app} size={22} gray />
            <span class="r-tx sm ex"><strong>{r.app.co}</strong> <span class="r-meta">· {r.meta}</span></span>
            <button class="r-btn reopen" onclick={(e) => reopen(e, r.app)}>Reopen</button>
          </div>
        {/if}
      {/each}
    </div>
  </div>
{/if}

<style>
  .pg { max-width: 1100px; margin: 0 auto; padding: 32px; width: 100%; box-sizing: border-box; }
  .loading { color: #8a9099; font-size: 13.5px; }

  /* header band */
  .band { background: #fff; border-bottom: 1px solid #e8e8e5; }
  .band-in { max-width: 1100px; margin: 0 auto; padding: 34px 32px 6px; }
  .b-date { font-size: 13px; color: #8a9099; margin-bottom: 4px; }
  .band h1 { font-size: 30px; font-weight: 700; letter-spacing: -0.02em; margin: 0 0 22px; }
  .cells { display: grid; grid-template-columns: repeat(5, 1fr); border-top: 1px solid #eeeeea; }
  .cell {
    display: block; padding: 18px 20px 16px; border-right: 1px solid #f0f0ed;
    background: none; border-top: 0; border-bottom: 0; border-left: 0; cursor: pointer;
    text-align: left; font-family: inherit;
  }
  .cell:last-child { border-right: 1px solid #f0f0ed; }
  .cell:hover { background: #fafaf8; }
  .c-top { display: flex; align-items: baseline; gap: 8px; }
  .c-n { font-size: 30px; font-weight: 700; letter-spacing: -0.03em; }
  .c-l { font-size: 12.5px; font-weight: 600; }
  .c-logos { display: flex; align-items: center; margin: 10px 0 5px; min-height: 24px; }
  .c-chip { display: inline-flex; border-radius: 8px; box-shadow: 0 0 0 2px #fff; }
  .c-extra { font-size: 11px; font-weight: 600; color: #8a9099; margin-left: 6px; }
  .c-note { display: block; font-size: 11px; color: #8a9099; }

  /* debrief banner (kept from Phase 3b) */
  .db-banner { width: 100%; display: flex; align-items: center; gap: 12px; text-align: left; cursor: pointer;
    background: #eef4ff; border: 1px solid #cdddfb; border-radius: 12px; padding: 13px 20px; font-family: inherit; margin: 0 0 14px; font-size: 13.5px; color: #4b5158; }
  .db-banner:hover { border-color: #2463eb; }
  .db-spark { color: #2463eb; flex: none; }
  .db-tx { flex: 1; min-width: 0; line-height: 1.4; }
  .db-tx b { color: #16181c; }
  .db-cta { flex: none; font-size: 13px; font-weight: 600; color: #2463eb; }

  /* noticing card */
  .noticing { display: flex; align-items: center; gap: 12px; border: 1px solid #e8e8e5; background: #fff; border-radius: 12px; padding: 13px 20px; margin-bottom: 32px; font-size: 13.5px; color: #4b5158; }
  .n-spark { color: #e0641f; flex: none; }
  .n-tx { min-width: 0; }
  .n-tx strong { color: #16181c; }
  .n-more { margin-left: auto; flex: none; font-size: 13px; font-weight: 500; color: #2463eb; text-decoration: none; }

  /* list */
  .list-hd { display: flex; align-items: baseline; gap: 10px; margin-bottom: 14px; }
  .lh-t { font-size: 11px; font-weight: 600; letter-spacing: .12em; text-transform: uppercase; color: #8a9099; }
  .lh-s { font-size: 12px; color: #8a9099; }
  .rows { display: flex; flex-direction: column; gap: 6px; }
  .row {
    display: flex; align-items: center; gap: 12px; background: #fff;
    border: 1px solid #eeeeea; border-radius: 10px; padding: 9px 16px; cursor: pointer;
  }
  .row:hover { border-color: #b9c6e8 !important; }
  .row.act { padding: 11px 16px; }
  .row.exit { background: #fbfbf9; opacity: .75; }
  .r-tx { flex: 1; min-width: 0; font-size: 13.5px; color: #16181c; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  .r-tx.sm { font-size: 13px; color: #4b5158; }
  .r-tx .q-co { color: #4b5158; }
  .r-tx.ex strong { color: #8a9099; }
  .r-tx.ex { color: #8a9099; }
  .r-meta { color: #8a9099; }
  .r-btn {
    border: 1px solid; border-radius: 8px; padding: 7px 14px; font-size: 12.5px;
    font-weight: 600; min-width: 126px; text-align: center; flex: none; cursor: pointer;
    font-family: inherit;
  }
  .r-btn.reopen { background: #fbfbf9; color: #8a9099; border-color: #e2e2de; }
  .r-btn.reopen:hover { color: #4b5158; border-color: #b8bdc4; }
  .r-quiet { font-size: 11.5px; color: #b8bdc4; min-width: 126px; text-align: center; flex: none; }

  /* empty state */
  .empty { max-width: 720px; margin: 0 auto; padding: 64px 32px 72px; text-align: center; }
  .e-date { font-size: 13px; color: #8a9099; margin-bottom: 6px; }
  .empty h1 { font-size: 30px; font-weight: 700; letter-spacing: -0.02em; margin: 0 0 10px; }
  .e-sub { font-size: 15px; line-height: 1.6; color: #4b5158; margin: 0 auto 28px; max-width: 440px; }
  .e-cta { display: inline-flex; flex-direction: column; gap: 10px; width: 460px; max-width: 100%; }
  .e-drop { border: 1.5px dashed #c8ccd2; background: #fff; border-radius: 14px; padding: 28px 24px; cursor: pointer; font-family: inherit; }
  .e-drop:hover { border-color: #2463eb; }
  .e-d1 { display: block; font-size: 14.5px; font-weight: 600; margin-bottom: 4px; color: #16181c; }
  .e-d2 { display: block; font-size: 12.5px; color: #8a9099; }
  .e-btns { display: flex; align-items: center; gap: 10px; justify-content: center; }
  .e-new { background: #2463eb; color: #fff; border: 0; border-radius: 9px; padding: 10px 20px; font-size: 13.5px; font-weight: 600; cursor: pointer; font-family: inherit; }
  .e-new:hover { background: #1a4fc4; }
  .e-new .nk { opacity: .65; font-weight: 400; }
  .e-legend { display: flex; align-items: center; justify-content: center; gap: 22px; margin-top: 44px; font-size: 12.5px; color: #b8bdc4; flex-wrap: wrap; }
  .e-li { display: flex; align-items: center; gap: 6px; }
  .e-li .dot { width: 7px; height: 7px; border-radius: 50%; }

  @media (max-width: 900px) {
    .band-in, .pg { padding-left: 16px; padding-right: 16px; }
    .cells { grid-template-columns: repeat(2, 1fr); }
    .cell { border-bottom: 1px solid #f0f0ed; }
    .r-btn, .r-quiet { min-width: 0; }
  }
</style>
