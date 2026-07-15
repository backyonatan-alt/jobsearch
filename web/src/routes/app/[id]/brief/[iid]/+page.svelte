<script>
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { api, isConnectionErr, pollForDossier } from '$lib/api.js';
  import { isPreview, mockApi } from '$lib/preview-mode.js';
  import { logEvent } from '$lib/analytics.js';
  import CompanyLogo from '$lib/CompanyLogo.svelte';
  import { toDisplayApp } from '$lib/app-helpers.js';

  const call = isPreview() ? mockApi : api;

  const id = $derived(page.params.id);
  const iid = $derived(page.params.iid);

  let app = $state(null);
  let dossier = $state(null);
  let round = $state(null);
  let loading = $state(true);
  let refreshing = $state(false);
  let refreshErr = $state('');
  let openLogged = null; // fire brief_page_open once per round viewed

  // Fake-door experiment (AB_TESTS.md "practice fake-door"): measures demand
  // for practice drills before building anything. One click per user counts.
  const PRACTICE_KEY = 'pursuit_practice_interest';
  let practiceClicked = $state(false);
  $effect(() => {
    try { practiceClicked = localStorage.getItem(PRACTICE_KEY) === '1'; } catch {}
  });
  function practiceInterest() {
    if (practiceClicked) return;
    practiceClicked = true;
    try { localStorage.setItem(PRACTICE_KEY, '1'); } catch {}
    logEvent('practice_interest', { app_id: Number(id) });
  }

  $effect(() => {
    void id; void iid;
    load();
  });

  async function load() {
    loading = true;
    try {
      const [a, d, ivs] = await Promise.all([
        call(`/api/applications/${id}`),
        call(`/api/applications/${id}/dossier?interview_id=${iid}`),
        call(`/api/applications/${id}/interviews`).catch(() => [])
      ]);
      if (!d) { goto(`/app/${id}`, { replaceState: true }); return; }
      app = toDisplayApp(a);
      dossier = d;
      round = (ivs || []).find(x => String(x.id) === String(iid)) || null;
      const key = `${id}:${iid}`;
      if (openLogged !== key) {
        openLogged = key;
        logEvent('brief_page_open', { app_id: Number(id) });
      }
    } catch (e) {
      if (e.message === 'unauthorized') return;
      // No brief for this round (404) or anything else broken → back to the detail page.
      goto(`/app/${id}`, { replaceState: true });
      return;
    } finally {
      loading = false;
    }
  }

  async function refresh() {
    if (refreshing) return;
    refreshing = true;
    refreshErr = '';
    const prevGeneratedAt = dossier?.generated_at ?? null;
    try {
      dossier = await call(`/api/applications/${id}/dossier/refresh`, {
        method: 'POST',
        body: JSON.stringify({ interview_id: Number(iid), interviewer_name: dossier?.interviewer_name || undefined })
      });
    } catch (e) {
      if (isConnectionErr(e)) {
        const d = await pollForDossier(`/api/applications/${id}/dossier?interview_id=${iid}`, prevGeneratedAt);
        if (d) { dossier = d; refreshing = false; return; }
      }
      refreshErr = /interview-prep limit/i.test(String(e.message))
        ? "You've used all your prep credits for the beta — email us for a top-up."
        : (e.message || 'Could not refresh.');
    } finally {
      refreshing = false;
    }
  }

  function print() { window.print(); }

  function initialsOf(name) {
    return (name || '').split(/\s+/).filter(Boolean).slice(0, 2).map(s => s[0]).join('').toUpperCase();
  }
  function domainOf(src) {
    if (!src) return '';
    try {
      return new URL(src.startsWith('http') ? src : `https://${src}`).hostname.replace(/^www\./, '');
    } catch { return src; }
  }

  const content = $derived(dossier?.content ?? null);
  const interviewer = $derived(content?.interviewer ?? null);
  const ivName = $derived(interviewer?.name ?? dossier?.interviewer_name ?? 'Your interviewer');
  const ivInitials = $derived(initialsOf(ivName) || '?');
  const ivMeta = $derived.by(() => {
    const bits = [];
    if (interviewer?.role) bits.push(interviewer.role);
    if (interviewer?.prior?.length) bits.push(interviewer.prior[0]);
    return bits.join(' · ');
  });
  const meeting = $derived(dossier?.meeting ?? null);

  const roundName = $derived((round?.summary || meeting?.panel || 'Round').trim() || 'Round');

  function whenStr() {
    const starts = round?.scheduled !== false ? (round?.starts_at || meeting?.starts_at) : null;
    if (!starts) return '';
    const d = new Date(starts);
    const startOfDay = (x) => new Date(x.getFullYear(), x.getMonth(), x.getDate());
    const days = Math.round((startOfDay(d) - startOfDay(new Date())) / 86400000);
    const time = d.toLocaleTimeString(undefined, { hour: 'numeric', minute: '2-digit' });
    if (days === 0) return `today, ${time}`;
    if (days === 1) return `tomorrow, ${time}`;
    if (days > 1 && days < 7) return `${d.toLocaleDateString(undefined, { weekday: 'long' })}, ${time}`;
    return `${d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })}, ${time}`;
  }
  const title = $derived(`${roundName} brief${whenStr() ? ` — ${whenStr()}` : ''}`);

  // "The one thing to remember" — snapshot first, else the strongest lands item.
  const oneThing = $derived(content?.snapshot || (content?.lands?.length ? content.lands[0] : ''));

  const signals = $derived(Array.isArray(content?.signals) ? content.signals : []);
  const lands = $derived(Array.isArray(content?.lands) ? content.lands : []);
  const avoid = $derived(Array.isArray(content?.avoid) ? content.avoid : []);
  const questions = $derived(Array.isArray(content?.questions) ? content.questions : []);
  const tells = $derived(Array.isArray(content?.style?.tells) ? content.style.tells : []);

  const linkedinHref = $derived.by(() => {
    const links = interviewer?.links || [];
    return links.find(l => /linkedin/i.test(l?.label || '') || /linkedin\.com/i.test(l?.href || ''))?.href || '';
  });

  // Every citation on the page in one chip row: explicit sources + signal deep
  // links + interviewer links, deduped by href.
  const allSources = $derived.by(() => {
    const out = [];
    const seen = new Set();
    const add = (label, href) => {
      if (!href || seen.has(href)) return;
      seen.add(href);
      out.push({ label: label || domainOf(href), href });
    };
    for (const s of (Array.isArray(content?.sources) ? content.sources : [])) add(s?.label, s?.href);
    for (const s of signals) add(s?.source || s?.kind, s?.source_url);
    for (const l of (interviewer?.links || [])) add(l?.label, l?.href);
    return out;
  });
</script>

<svelte:head>
  <title>{app?.co ? `${roundName} brief — ${app.co} — Pursuit` : 'Round brief — Pursuit'}</title>
</svelte:head>

<div class="wrap">
  {#if loading}
    <p class="loading">Loading…</p>
  {:else if app && dossier}

    <div class="crumb noprint">
      <a href="/app/applications">Applications</a> <span class="sep">/</span>
      <a href={`/app/${id}`}>{app.co}</a> <span class="sep">/</span>
      <span class="here">Round brief</span>
    </div>

    <div class="ttl">
      <CompanyLogo app={app} size={40} radius={10} />
      <h1>{title}</h1>
    </div>
    <div class="meta">
      {app.co} · {app.role}{#if dossier.generatedAgo}&nbsp;· generated {dossier.generatedAgo}{/if}
      <span class="noprint">
        · <button class="linkbtn" onclick={refresh} disabled={refreshing}>{refreshing ? 'refreshing…' : 'refresh'}</button>
        · <button class="linkbtn" onclick={print}>print</button>
      </span>
    </div>
    {#if refreshErr}<p class="err noprint">{refreshErr}</p>{/if}

    {#if oneThing}
      <div class="one">
        <div class="one-hd">✦ The one thing to remember</div>
        <div class="one-tx">{@html oneThing}</div>
      </div>
    {/if}

    <!-- Interviewer -->
    <div class="person">
      <div class="p-av">{ivInitials}</div>
      <div class="p-main">
        <div class="p-name">{ivName}
          {#if ivMeta}<span class="p-meta">· {ivMeta}</span>{/if}
          {#if linkedinHref}<a class="p-li noprint" href={linkedinHref} target="_blank" rel="noopener">LinkedIn ↗</a>{/if}
        </div>
        {#if content?.background}
          <div class="p-brief">{content.background}</div>
        {/if}
        {#if content?.style?.lead}
          <div class="p-line"><strong class="k">How they interview:</strong> {content.style.lead}</div>
        {/if}
        {#if tells.length}
          <div class="p-line tells">{#each tells as t, i (t.lbl)}{#if i > 0}<span class="tsep"> · </span>{/if}<span><strong class="k">{t.lbl}:</strong> {t.val}</span>{/each}</div>
        {/if}
        {#if lands.length}
          <div class="p-line"><strong class="k lands">What lands:</strong>
            <span class="p-list">{#each lands as l, i (i)}<span class="li"><span class="g pos">✓</span>{l}</span>{/each}</span>
          </div>
        {/if}
        {#if avoid.length}
          <div class="p-line"><strong class="k avoid">Avoid:</strong>
            <span class="p-list">{#each avoid as a, i (i)}<span class="li"><span class="g neg">✕</span>{a}</span>{/each}</span>
          </div>
        {/if}
        {#if interviewer?.links?.length}
          <div class="p-srcs">Sources: {#each interviewer.links as l, i (l.href)}{#if i > 0} · {/if}<a href={l.href} target="_blank" rel="noopener">{l.label || domainOf(l.href)}</a>{/each}</div>
        {/if}
      </div>
    </div>

    {#if signals.length}
      <div class="rule"></div>
      <div class="sec-hd">Recent signals</div>
      <div class="signals">
        {#each signals as s, i (i)}
          <div class="sig">
            <span class="sig-date">{s.date ?? ''}</span>
            <span class="sig-body">
              {#if s.kind}<span class="sig-kind">{s.kind}</span>{/if}{s.body}
              {#if s.source_url}
                <a class="sig-src" href={s.source_url} target="_blank" rel="noreferrer">{s.source || domainOf(s.source_url)}</a>
              {:else if s.source}
                <span class="sig-src plain">{s.source}</span>
              {/if}
            </span>
          </div>
        {/each}
      </div>
    {/if}

    {#if questions.length}
      <div class="qcard">
        <div class="q-hd">Questions worth asking them</div>
        <div class="q-list">
          {#each questions as item, i (i)}
            <div class="q-row">· “{item.q}”{#if item.why}&nbsp;<span class="q-why">({item.why})</span>{/if}</div>
          {/each}
        </div>
      </div>
    {/if}

    <!-- practice fake-door (see AB_TESTS.md) -->
    <div class="practice noprint">
      {#if practiceClicked}
        <span class="pr-spark">✓</span>
        <span class="pr-tx"><strong>Noted — practice drills are on our list.</strong> Your interest helps us decide what to build next.</span>
      {:else}
        <span class="pr-spark">✦</span>
        <span class="pr-tx"><strong>Practice this round.</strong> Turn the likely questions into a drill — answer out loud, get sharper before you walk in.</span>
        <button class="pr-btn" onclick={practiceInterest}>I'd use this →</button>
      {/if}
    </div>

    {#if allSources.length}
      <div class="srcs">
        <span class="srcs-lbl">All sources</span>
        {#each allSources as s (s.href)}
          <a class="src-chip" href={s.href} target="_blank" rel="noreferrer">{s.label}</a>
        {/each}
      </div>
    {/if}

    <div class="foot">
      Synthesised from public posts, talks, and papers — always verify before you walk in · <a href="/privacy" target="_blank" rel="noreferrer">how we research people</a>
    </div>
    <div class="foot hook">
      After the round: a 20-second debrief sharpens your next round's prep. <a href={`/app/${id}?debrief=${iid}`}>Debrief this round →</a>
    </div>

  {/if}
</div>

<style>
  .wrap {
    max-width: 760px; width: 100%; box-sizing: border-box;
    margin: 0 auto; padding: 28px 32px 90px;
    color: #16181c;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
  }
  .wrap :global(a) { color: #2463eb; text-decoration: none; }
  .loading { color: #8a9099; font-size: 13px; }
  .linkbtn { background: none; border: 0; cursor: pointer; padding: 0; font-family: inherit; font-size: inherit; color: #2463eb; }
  .linkbtn:hover:not(:disabled) { text-decoration: underline; }
  .linkbtn:disabled { color: #8a9099; cursor: default; }

  .crumb { font-size: 13px; color: #8a9099; margin-bottom: 24px; }
  .crumb a { color: #8a9099; }
  .crumb a:hover { color: #4b5158; }
  .crumb .sep { margin: 0 4px; }
  .crumb .here { color: #16181c; font-weight: 600; }

  .ttl { display: flex; align-items: center; gap: 14px; margin-bottom: 6px; }
  .ttl h1 { font-size: 28px; font-weight: 700; letter-spacing: -0.02em; margin: 0; }
  .meta { font-size: 13px; color: #8a9099; margin-bottom: 26px; }
  .err { color: #b3372a; font-size: 13px; margin: -16px 0 20px; }

  .one { background: #eef4ff; border: 1px solid #cdddfb; border-radius: 14px; padding: 18px 22px; margin-bottom: 34px; }
  .one-hd { font-size: 13px; font-weight: 700; color: #2463eb; margin-bottom: 6px; }
  .one-tx { font-size: 15px; line-height: 1.6; color: #1e3a6e; }
  .one-tx :global(em) { font-style: normal; font-weight: 600; }

  .person { display: flex; gap: 16px; margin-bottom: 10px; }
  .p-av { width: 44px; height: 44px; border-radius: 50%; background: #eef4ff; color: #2463eb;
    display: flex; align-items: center; justify-content: center; font-size: 14px; font-weight: 700; flex: none; }
  .p-main { flex: 1; min-width: 0; }
  .p-name { font-size: 16.5px; font-weight: 700; }
  .p-meta { font-size: 12.5px; font-weight: 400; color: #8a9099; }
  .p-li { font-size: 12px; margin-left: 6px; }
  .p-brief { font-size: 14px; line-height: 1.65; color: #4b5158; margin: 8px 0 12px; }
  .p-line { font-size: 13.5px; line-height: 1.65; color: #4b5158; margin-bottom: 8px; }
  .p-line .k { color: #16181c; }
  .p-line .k.lands { color: #1d7a4f; }
  .p-line .k.avoid { color: #b3372a; }
  .p-line.tells { color: #6f7680; font-size: 13px; }
  .p-list { display: block; margin-top: 4px; }
  .p-list .li { display: flex; gap: 8px; align-items: baseline; margin-bottom: 3px; }
  .p-list .g { font-size: 11px; flex: none; }
  .p-list .g.pos { color: #1d7a4f; }
  .p-list .g.neg { color: #b3372a; }
  .p-srcs { font-size: 12px; color: #8a9099; margin-top: 4px; }

  .rule { border-top: 1px solid #e2e2de; margin: 22px 0; }
  .sec-hd { font-size: 13px; font-weight: 700; margin-bottom: 12px; }
  .signals { display: flex; flex-direction: column; gap: 10px; }
  .sig { display: grid; grid-template-columns: 72px 1fr; gap: 14px; font-size: 13.5px; }
  .sig-date { font-size: 12px; font-weight: 500; color: #8a9099; padding-top: 2px; }
  .sig-body { color: #4b5158; line-height: 1.55; }
  .sig-kind { display: inline-block; font-size: 11px; font-weight: 500; color: #8a9099; margin-right: 8px;
    padding: 1px 7px; background: #f0f0ed; border-radius: 4px; vertical-align: 1px; }
  .sig-src { font-size: 12px; margin-left: 6px; }
  .sig-src.plain { color: #8a9099; }

  .qcard { background: #fff; border: 1px solid #e8e8e5; border-radius: 14px; padding: 20px 24px; margin-top: 34px; }
  .q-hd { font-size: 13px; font-weight: 700; margin-bottom: 10px; }
  .q-list { display: flex; flex-direction: column; gap: 8px; font-size: 13.5px; line-height: 1.6; color: #4b5158; }
  .q-why { color: #8a9099; }

  .srcs { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; font-size: 12px; color: #8a9099; margin-top: 18px; }
  .srcs-lbl { font-weight: 600; color: #6f7680; }
  .src-chip { border: 1px solid #e8e8e5; border-radius: 14px; padding: 3px 10px; background: #fff; color: #2463eb; }
  .src-chip:hover { border-color: #cdddfb; background: #eef4ff; }

  .practice {
    display: flex; align-items: center; gap: 12px;
    background: #fff7f1; border: 1px solid #f0d9c4; border-radius: 12px;
    padding: 14px 20px; margin-top: 18px; font-size: 13.5px; color: #4b5158;
  }
  .practice .pr-spark { color: #e0641f; flex: none; }
  .practice .pr-tx { flex: 1; min-width: 0; line-height: 1.5; }
  .practice .pr-tx strong { color: #16181c; }
  .practice .pr-btn {
    background: #fff; color: #c05310; border: 1px solid #f0d9c4; border-radius: 8px;
    padding: 7px 14px; font-size: 12.5px; font-weight: 600; cursor: pointer;
    flex: none; font-family: inherit; white-space: nowrap;
  }
  .practice .pr-btn:hover { border-color: #e0641f; }

  .foot { font-size: 12.5px; color: #8a9099; margin-top: 22px; border-top: 1px solid #e2e2de; padding-top: 14px; }
  .foot.hook { border-top: 0; padding-top: 0; margin-top: 10px; }

  @media (max-width: 700px) {
    .wrap { padding: 20px 16px 60px; }
    .ttl h1 { font-size: 22px; }
    .sig { grid-template-columns: 60px 1fr; gap: 10px; }
  }
  @media print {
    .noprint { display: none !important; }
    .wrap { padding: 0; max-width: none; }
  }
</style>
