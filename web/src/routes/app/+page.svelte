<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/state';
  import { api } from '$lib/api.js';
  import {
    STATUS_LABEL, toDisplayApp, daysSince
  } from '$lib/app-helpers.js';
  import Onboarding from '$lib/Onboarding.svelte';
  import AddApplication from '$lib/AddApplication.svelte';

  // `api` self-routes to the in-memory mock when ?preview=1 is on (see
  // $lib/api.js → preview-mode.js), so the Brief renders with no backend.
  const call = api;

  let apps = $state([]);          // display apps
  let me = $state(null);
  let loading = $state(true);
  let showNewModal = $state(false);
  let showOnboarding = $state(false);

  // Per-app interview events + dossier, keyed by app id. Loaded lazily for the
  // "live" statuses so the Brief can find the soonest upcoming interview and a
  // sourced framing without N round-trips for the whole pipeline.
  let eventsByApp = $state({});   // id -> [{starts_at, ends_at, summary, location, medium, panel, interviewer}]
  let dossierByApp = $state({});  // id -> dossier DTO (content + meeting + interviewer_name)

  onMount(refresh);

  async function refresh() {
    try {
      const [meRes, raw] = await Promise.all([
        call('/api/me'),
        call('/api/applications')
      ]);
      me = meRes;
      apps = raw.map(toDisplayApp);
      maybeShowOnboarding();
      loading = false;
      await hydrateBrief(raw);
    } catch (e) {
      if (e.message !== 'unauthorized') console.error(e);
      loading = false;
    }
  }

  // Pull interviews + dossier for the live-pipeline apps so the Brief can name
  // today's interview and surface a sourced inference. Soft-fails per app.
  async function hydrateBrief(raw) {
    const live = raw.filter(a => ['screen', 'interview', 'offer'].includes(a.status));
    await Promise.all(live.map(async (a) => {
      const [ivs, doss] = await Promise.all([
        call(`/api/applications/${a.id}/interviews`).catch(() => []),
        call(`/api/applications/${a.id}/dossier`).catch(() => null)
      ]);
      if (Array.isArray(ivs) && ivs.length) eventsByApp[a.id] = ivs;
      if (doss) dossierByApp[a.id] = doss;
    }));
    // reassign to trigger reactivity
    eventsByApp = { ...eventsByApp };
    dossierByApp = { ...dossierByApp };
    seedTasks();
  }

  function maybeShowOnboarding() {
    const forced = page.url.searchParams.get('onboarding') === '1';
    const fresh = !me?.onboarded_at && apps.length === 0;
    showOnboarding = forced || fresh;
  }
  function finishOnboarding() { showOnboarding = false; refresh(); }

  // ⌘N / Ctrl+N opens the new-application modal from anywhere on Today.
  function onKeydown(e) {
    if ((e.metaKey || e.ctrlKey) && (e.key === 'n' || e.key === 'N') && !showNewModal) {
      e.preventDefault();
      showNewModal = true;
    }
  }

  // ── Brief: find the soonest upcoming interview ──────────────
  // An app's interview events come from /interviews; the dossier may also carry
  // a linked `meeting` (set when an .ics is attached). We merge both, keep only
  // future-dated ones, and pick the soonest.
  const upcomingEvents = $derived.by(() => {
    const now = Date.now();
    const out = [];
    for (const a of apps) {
      if (!['screen', 'interview', 'offer'].includes(a.status)) continue;
      const doss = dossierByApp[a.id];
      const ivs = eventsByApp[a.id] || [];
      for (const ev of ivs) {
        if (!ev.starts_at) continue;
        const t = new Date(ev.starts_at).getTime();
        if (isNaN(t) || t < now) continue;
        out.push({
          app: a, startsAt: ev.starts_at, endsAt: ev.ends_at,
          summary: ev.summary || `${a.co} interview`,
          medium: ev.location || doss?.meeting?.medium || '',
          interviewer: pickInterviewer(ev.attendees) || doss?.interviewer_name || '',
          panel: doss?.meeting?.panel || ''
        });
      }
      // Dossier-linked meeting (if it has a real timestamp).
      const m = doss?.meeting;
      if (m?.starts_at) {
        const t = new Date(m.starts_at).getTime();
        if (!isNaN(t) && t >= now) {
          out.push({
            app: a, startsAt: m.starts_at, endsAt: m.ends_at,
            summary: doss.interviewer_name ? `${a.co} · ${m.panel || 'interview'}` : `${a.co} interview`,
            medium: m.medium || '', interviewer: doss.interviewer_name || '', panel: m.panel || ''
          });
        }
      }
    }
    // De-dupe identical (app, startsAt) pairs — a linked meeting often mirrors
    // an interview event — then sort soonest-first.
    const seen = new Set();
    const deduped = out.filter(e => {
      const k = `${e.app.id}|${e.startsAt}`;
      if (seen.has(k)) return false;
      seen.add(k);
      return true;
    });
    deduped.sort((a, b) => new Date(a.startsAt) - new Date(b.startsAt));
    return deduped;
  });

  function pickInterviewer(attendees) {
    if (!attendees) return '';
    const arr = Array.isArray(attendees) ? attendees : [attendees];
    const named = arr.find(x => x && typeof x === 'object' && x.name && !/recruit/i.test(x.name || ''));
    return named?.name || '';
  }

  const nextInterview = $derived(upcomingEvents[0] || null);
  // Dossier for the soonest interview — the source for the insight + tips.
  const nextDossier = $derived(nextInterview ? dossierByApp[nextInterview.app.id] : null);
  const nextContent = $derived(nextDossier?.content || null);

  // One-line sourced framing for the insight box (from dossier snapshot).
  const insightText = $derived.by(() => {
    const c = nextContent;
    if (!c) return '';
    return c.snapshot || c.style?.lead || c.background || '';
  });

  // "Worth reviewing" — the things that land with this person (3 max).
  const worthReviewing = $derived.by(() => {
    const c = nextContent;
    if (!c || !Array.isArray(c.lands) || c.lands.length === 0) return [];
    return c.lands.slice(0, 3);
  });

  // Two quick tips — prefer company watch_fors (grounded loop advice), else
  // fall back to the remaining "lands". Empty → block omitted.
  const quickTips = $derived.by(() => {
    const c = nextContent;
    if (!c) return [];
    const wf = c.company?.watch_fors;
    if (Array.isArray(wf) && wf.length >= 2) return wf.slice(0, 2);
    if (Array.isArray(c.lands) && c.lands.length >= 5) return c.lands.slice(3, 5);
    return [];
  });

  // Meta line: "60 min · Google Meet · Final round".
  const metaBits = $derived.by(() => {
    const e = nextInterview;
    if (!e) return [];
    const bits = [];
    if (e.startsAt && e.endsAt) {
      const mins = Math.round((new Date(e.endsAt) - new Date(e.startsAt)) / 60000);
      if (mins > 0) bits.push(mins >= 60 && mins % 60 === 0 ? `${mins / 60} hr` : `${mins} min`);
    }
    if (e.medium) bits.push(e.medium);
    if (e.panel) bits.push(e.panel);
    else if (e.summary && !e.panel) bits.push(e.summary);
    return bits;
  });

  // Agenda — the upcoming events AFTER the lede one (so "later this week" reads
  // as the rest of the schedule), up to 3. With no lede interview we show the
  // soonest 3 directly.
  const agenda = $derived(
    nextInterview ? upcomingEvents.slice(1, 4) : upcomingEvents.slice(0, 3)
  );

  function fmtAgendaWhen(iso) {
    const d = new Date(iso);
    const day = d.toLocaleDateString(undefined, { weekday: 'short' });
    const time = d.toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit' });
    return { day, time };
  }

  // ── Pulse (right pane) ──────────────────────────────────────
  const activeApps = $derived(apps.filter(a => ['applied', 'screen', 'interview', 'offer'].includes(a.status)));
  const quietApps  = $derived(apps.filter(a => a.stale));
  // Awaiting = applied/screen with no upcoming event linked.
  const awaitingApps = $derived(apps.filter(a =>
    ['applied', 'screen'].includes(a.status) &&
    !upcomingEvents.some(e => e.app.id === a.id)
  ));

  const activeCount   = $derived(activeApps.length);
  const awaitingCount = $derived(awaitingApps.length);
  const quietCount    = $derived(quietApps.length);

  // Waiting list — longest first (by days since applied).
  const waiting = $derived(
    [...awaitingApps].sort((a, b) =>
      (daysSince(b.raw.applied_at) ?? 0) - (daysSince(a.raw.applied_at) ?? 0)
    )
  );
  function waitDays(a) { return daysSince(a.raw.applied_at) ?? 0; }

  // ── "Your move" tasks — local only, persisted to localStorage ──
  const TASKS_KEY = 'pursuit_tasks';
  let tasks = $state([]);
  let tasksLoaded = false;

  function seedTasks() {
    if (tasksLoaded) return;
    tasksLoaded = true;
    let stored = null;
    try { stored = JSON.parse(localStorage.getItem(TASKS_KEY) || 'null'); } catch {}
    if (Array.isArray(stored) && stored.length) { tasks = stored; return; }
    // Seed from upcoming interviews + the most pressing pipeline facts.
    const seeded = [];
    const e = upcomingEvents[0];
    if (e) {
      const who = e.interviewer ? e.interviewer.split(' ')[0] : e.app.co;
      seeded.push({ id: 't-prep', b: `Prep 3 questions for ${who}`, s: `${e.summary} · ${e.app.co}`, due: 'Today', hot: true, done: false });
    }
    const offer = apps.find(a => a.status === 'offer');
    if (offer) seeded.push({ id: 't-offer', b: `Decide on the ${offer.co} offer`, s: offer.raw.salary_note || offer.role, due: 'Soon', hot: true, done: false });
    const q = quietApps[0];
    if (q) seeded.push({ id: 't-quiet', b: `Follow up on ${q.co}`, s: `Quiet ${waitDays(q)} days · log it once you reach out`, due: `${waitDays(q)}d`, hot: false, done: false });
    if (seeded.length === 0) {
      seeded.push({ id: 't-empty', b: 'Add your next application', s: 'Press ⌘N to log a role you just applied to', due: '', hot: false, done: false });
    }
    tasks = seeded;
    persistTasks();
  }
  function persistTasks() {
    try { localStorage.setItem(TASKS_KEY, JSON.stringify(tasks)); } catch {}
  }
  function toggleTask(id) {
    tasks = tasks.map(t => t.id === id ? { ...t, done: !t.done } : t);
    persistTasks();
  }
  const openTaskCount = $derived(tasks.filter(t => !t.done).length);

  // Advisory footer — only when something has gone quiet.
  const advisoryNames = $derived(quietApps.slice(0, 2).map(a => a.co));
  const advisoryLabel = $derived(
    advisoryNames.length === 2 ? `${advisoryNames[0]} and ${advisoryNames[1]}`
    : advisoryNames.length === 1 ? advisoryNames[0] : ''
  );

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
  // "Wednesday · 20 May 2026" — sans, never mono.
  const dateLine = $derived.by(() => {
    const d = new Date();
    const dow = d.toLocaleDateString('en-US', { weekday: 'long' });
    const rest = d.toLocaleDateString('en-US', { day: 'numeric', month: 'long', year: 'numeric' });
    return `${dow} · ${rest}`;
  });

  function openDetail(id) { goto(`/app/${id}`); }
  function openBoard() { goto('/app/board'); }
  function openPlaybook(id) { goto(`/app/${id}/playbook`); }
</script>

<svelte:head>
  <title>Today — Pursuit</title>
</svelte:head>

<svelte:window onkeydown={onKeydown} />

{#if showOnboarding}
  <Onboarding onDone={finishOnboarding} />
{/if}

<div class="topbar">
  <div class="crumb"><span class="here">Today</span></div>
  <div class="right">
    <div class="search">
      <svg class="ico" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5">
        <circle cx="7" cy="7" r="4.5"/><path d="M10.5 10.5L13 13"/>
      </svg>
      <span>Search applications, people…</span>
      <span class="kbd">⌘K</span>
    </div>
    <button class="btn btn-primary" onclick={() => (showNewModal = true)}>
      New application <span class="kbd">⌘N</span>
    </button>
  </div>
</div>

<div class="ob ob-alt ob-swap">
  <!-- ══ LEFT — The Brief ══════════════════════════════════ -->
  <div class="brief">
    <div class="brief-in">
      <div class="brief-date">{dateLine}</div>
      <div class="brief-head">
        <h1>{greeting},<br /><b>{firstName}.</b></h1>
        {#if !loading}
          <div class="brief-stats">
            <button class="bstat" onclick={openBoard} title="Applications still active — applied through offer">
              <span class="bstat-n">{activeCount}</span>
              <span class="bstat-l">In progress</span>
            </button>
            <button class="bstat" onclick={openBoard} title="Applied and waiting to hear back">
              <span class="bstat-n">{awaitingCount}</span>
              <span class="bstat-l">Awaiting reply</span>
            </button>
            <button class="bstat" class:warn={quietCount > 0} onclick={openBoard} title="No reply in over a week">
              <span class="bstat-n">{quietCount}</span>
              <span class="bstat-l">Gone quiet</span>
            </button>
          </div>
        {/if}
      </div>

      {#if loading}
        <p class="lede">Loading your day…</p>
      {:else if nextInterview}
        <p class="lede">
          Today it's your
          <span class="hot">interview at {nextInterview.app.co}</span>
          — the {nextInterview.app.role} role{#if nextInterview.interviewer}, one-on-one with {nextInterview.interviewer}{/if}.
          {#if upcomingEvents.length > 1}{upcomingEvents.length - 1} more {upcomingEvents.length - 1 === 1 ? 'event follows' : 'events follow'} later this week.{/if}
        </p>

        <div class="kick">{@render Spark()}&nbsp;Prep for today</div>

        {#if insightText}
          <div class="insight">
            <span class="ic">{@render Spark(15)}</span>
            <span class="tx">{@html insightText}</span>
          </div>
        {/if}

        {#if worthReviewing.length}
          <div class="brief-sub">Worth reviewing</div>
          <ul class="brief-review">
            {#each worthReviewing as item}
              <li>{@html item}</li>
            {/each}
          </ul>
        {/if}

        {#if quickTips.length}
          <div class="brief-sub">Two quick tips</div>
          {#each quickTips as tip}
            <div class="brief-tip"><span class="sp">{@render Spark()}</span><span>{@html tip}</span></div>
          {/each}
        {/if}

        {#if metaBits.length}
          <div class="meta">
            {#each metaBits as b, i}
              {#if i === 0}<b>{b}</b>{:else}<span class="dot"></span>{b}{/if}
            {/each}
          </div>
        {/if}

        <button class="cta" onclick={() => openPlaybook(nextInterview.app.id)}>Open interview prep {@render Arrow()}</button>
      {:else}
        <p class="lede">Nothing's on the calendar today — here's where your search stands.</p>
      {/if}

      {#if !loading && agenda.length}
        <div class="agenda">
          <div class="kick">Later this week</div>
          {#each agenda as e}
            {@const w = fmtAgendaWhen(e.startsAt)}
            <div class="ag-row" onclick={() => openDetail(e.app.id)} role="button" tabindex="0">
              <span class="when"><b>{w.day}</b> {w.time}</span>
              <span><span class="co">{e.app.co}</span> <span class="role">· {e.app.role}</span></span>
              <span class={`pill ${e.app.status}`}><span class="pdot"></span>{STATUS_LABEL[e.app.status]}</span>
            </div>
          {/each}
        </div>
      {/if}

      {#if !loading}
        <div class="foot">
          <button class="foot-link" onclick={openBoard}>View all {apps.length} {apps.length === 1 ? 'application' : 'applications'} {@render Arrow()}</button>
        </div>
      {/if}
    </div>
  </div>

  <!-- ══ RIGHT — Where things stand (pulse) ════════════════ -->
  <div class="pulse-stage">
    <div class="pulse-tag"><span class="d"></span>Where things stand</div>


    <div class="tasks">
      <div class="pulse-sec">
        <span class="t">Your move</span>
        <span class="c">{openTaskCount} to do</span>
      </div>
      {#each tasks as t (t.id)}
        <div class={`task ${t.done ? 'done' : ''}`} onclick={() => toggleTask(t.id)} role="button" tabindex="0">
          <span class="box"></span>
          <span class="tx"><b>{t.b}</b><small>{t.s}</small></span>
          {#if t.due}<span class={`due ${t.hot && !t.done ? 'hot' : ''}`}>{t.due}</span>{/if}
        </div>
      {/each}
      <div class="tasks-note">Personal checklist · stays on this device.</div>
    </div>

    {#if advisoryLabel}
      <div class="pulse-foot">
        <span class="fic">{@render Spark(15)}</span>
        <span class="ftx"><b>{advisoryLabel} {quietApps.length > 1 ? 'have' : 'has'} gone quiet</b><small>No reply in over a week — it might be a good time to reach out to them directly.</small></span>
        <button class="pulse-link" onclick={openBoard}>{quietApps.length > 1 ? 'See both' : 'See it'} {@render Arrow()}</button>
      </div>
    {/if}

    <div class="pulse-sec waiting-sec">
      <span class="t">Waiting to hear back</span>
      <span class="c">longest first</span>
    </div>
    <div class="pulse-list">
      {#if waiting.length === 0}
        <div class="pulse-empty">Nothing waiting — every open thread has a next step.</div>
      {:else}
        {#each waiting as w (w.id)}
          <div class={`pulse-row ${w.stale ? 'quiet' : ''}`} onclick={() => openDetail(w.id)} role="button" tabindex="0">
            {#if w.logoSrc}
              <img class="row-logo" src={w.logoSrc} alt="" />
            {:else}
              <span class={`row-logo letter ${w.logoCls}`}>{w.coShort}</span>
            {/if}
            <span class="wx"><b>{w.co}</b><small>{STATUS_LABEL[w.status]}</small></span>
            <span class="days">{waitDays(w)}d</span>
            <span class="ok"><span class={`okdot ${w.stale ? 'warn' : ''}`}></span></span>
          </div>
        {/each}
      {/if}
    </div>
  </div>
</div>

<AddApplication bind:open={showNewModal} onCreated={refresh} />

{#snippet Spark(s)}
  <svg width={s ?? 13} height={s ?? 13} viewBox="0 0 16 16" fill="currentColor" aria-hidden="true"><path d="M8 1l1.5 4.2L14 7l-4.5 1.8L8 13l-1.5-4.2L2 7l4.5-1.8z"/></svg>
{/snippet}
{#snippet Arrow()}
  <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.6" aria-hidden="true"><path d="M3 8h9M8 4l4 4-4 4" stroke-linecap="round" stroke-linejoin="round"/></svg>
{/snippet}

<style>
  /* ════ Two-pane Today (Option B, swapped) ════════════════ */
  /* The page is a flex child of .main (column, overflow hidden). We fill the
     remaining height and let each pane scroll independently. */
  .ob {
    flex: 1; min-height: 0;
    display: grid; grid-template-columns: 0.92fr 1.08fr;
    font-family: var(--sans); color: var(--ink);
  }
  .ob.ob-swap { grid-template-columns: 1.08fr 0.92fr; }
  .ob.ob-swap .brief { border-right: 1px solid var(--rule); }

  /* ── Topbar (matches shell convention) ── */
  .topbar { display: flex; justify-content: space-between; align-items: center; padding: 0 20px; height: 48px; border-bottom: 1px solid var(--rule); background: var(--surface); flex-shrink: 0; }
  .crumb .here { font-weight: 500; font-size: 14px; }
  .right { display: flex; align-items: center; gap: 8px; }
  .search { display: flex; align-items: center; gap: 6px; background: var(--card); border: 1px solid var(--rule); border-radius: 7px; padding: 5px 10px; font-size: 13px; color: var(--mute); min-width: 280px; }
  .search .ico { width: 14px; height: 14px; }
  .search .kbd { margin-left: auto; font-family: var(--mono); font-size: 11px; color: var(--mute-2); }
  .btn { background: var(--card); border: 1px solid var(--rule); border-radius: 7px; padding: 6px 11px; font-size: 13px; font-weight: 500; color: var(--ink); cursor: pointer; }
  .btn-primary { background: var(--accent); border-color: var(--accent-strong); color: white; display: inline-flex; align-items: center; gap: 8px; }
  .btn-primary .kbd { font-family: var(--mono); font-size: 11px; color: rgba(255,255,255,.75); }

  /* ── LEFT: editorial brief ── */
  .brief { overflow-y: auto; }
  .brief-in { max-width: 560px; padding: 44px 48px 56px; }
  .brief-date { font-size: 13px; color: var(--mute); margin-bottom: 16px; letter-spacing: -0.003em; }
  .brief h1 { font-size: 30px; font-weight: 300; letter-spacing: -0.03em; line-height: 1.12; margin: 0 0 12px; }
  .brief h1 b { font-weight: 500; }
  .lede { font-size: 14.5px; color: var(--ink-2); line-height: 1.6; margin: 0 0 30px; max-width: 50ch; }
  .lede .hot { color: var(--warm-text); font-weight: 500; }

  .brief-head { display: flex; align-items: center; justify-content: space-between; gap: 32px; margin: 0 0 30px; flex-wrap: wrap; }
  .brief-head h1 { margin: 0; }
  .brief-stats { display: flex; align-items: center; }
  .brief-stats .bstat { display: flex; flex-direction: column; align-items: flex-start; gap: 3px; cursor: pointer; padding: 0 22px; transition: opacity .12s; }
  .brief-stats .bstat:first-child { padding-left: 0; }
  .brief-stats .bstat:last-child { padding-right: 0; }
  .brief-stats .bstat + .bstat { border-left: 1px solid var(--rule); }
  .brief-stats .bstat:hover { opacity: 0.65; }
  .brief-stats .bstat-n { font-size: 30px; font-weight: 500; line-height: 1; letter-spacing: -0.022em; color: var(--ink); font-variant-numeric: tabular-nums; }
  .brief-stats .bstat-l { font-size: 12px; color: var(--mute); letter-spacing: -0.003em; }
  .brief-stats .bstat.warn .bstat-n { color: var(--warm-text); }

  .kick { font-size: 11.5px; font-weight: 600; letter-spacing: 0.07em; text-transform: uppercase; color: var(--mute-2); margin-bottom: 14px; display: flex; align-items: center; gap: 10px; }
  .kick::after { content: ""; flex: 1; height: 1px; background: var(--rule); }

  .insight { display: flex; gap: 12px; padding: 15px 16px; background: var(--accent-tint); border-radius: 13px; margin-bottom: 14px; }
  .insight .ic { color: var(--accent-text); flex-shrink: 0; margin-top: 1px; }
  .insight .tx { font-size: 13.5px; line-height: 1.55; color: var(--accent-text); }
  .insight .tx :global(b), .insight .tx :global(em) { font-weight: 600; font-style: normal; }

  .brief-sub { font-size: 11.5px; font-weight: 600; letter-spacing: 0.06em; text-transform: uppercase; color: var(--mute-2); margin: 24px 0 13px; }
  .brief-review { list-style: none; margin: 0; padding: 0; display: flex; flex-direction: column; gap: 11px; }
  .brief-review li { font-size: 13.5px; line-height: 1.5; color: var(--ink-2); padding-left: 25px; position: relative; }
  .brief-review li :global(b) { font-weight: 600; }
  .brief-review li::before {
    content: ""; position: absolute; left: 4px; top: 4px; width: 6px; height: 10px;
    border: solid var(--accent-text); border-width: 0 1.7px 1.7px 0; transform: rotate(42deg);
  }
  .brief-tip { display: flex; gap: 11px; font-size: 13px; line-height: 1.55; color: var(--mute); margin-top: 11px; }
  .brief-tip .sp { color: var(--warm-text); flex-shrink: 0; margin-top: 2px; }
  .brief-tip :global(b) { color: var(--ink-2); font-weight: 500; }

  .meta { font-size: 13px; color: var(--mute); margin: 24px 0 16px; display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
  .meta b { color: var(--ink-2); font-weight: 500; }
  .meta .dot { width: 3px; height: 3px; border-radius: 50%; background: var(--mute-2); }

  .cta { display: inline-flex; align-items: center; justify-content: center; gap: 8px; background: var(--ink); color: white; border: none; font-size: 13.5px; font-weight: 500; padding: 11px 18px; border-radius: 10px; cursor: pointer; transition: background .12s; }
  .cta:hover { background: var(--ink-2); }

  .agenda { margin-top: 36px; }
  .ag-row { display: grid; grid-template-columns: 70px 1fr auto; gap: 14px; align-items: center; padding: 13px 2px; border-top: 1px solid var(--rule); font-size: 13.5px; cursor: pointer; transition: background .12s; }
  .ag-row:hover { background: var(--surface-2); }
  .ag-row .when { font-variant-numeric: tabular-nums; color: var(--mute); }
  .ag-row .when b { color: var(--ink-2); font-weight: 500; }
  .ag-row .co { font-weight: 500; }
  .ag-row .role { color: var(--mute); font-size: 12.5px; }

  .foot { margin-top: 28px; display: flex; justify-content: flex-end; }
  .foot-link { background: none; border: none; padding: 4px 0; font-family: inherit; font-size: 12.5px; color: var(--mute); display: inline-flex; align-items: center; gap: 6px; cursor: pointer; transition: color .12s; }
  .foot-link:hover { color: var(--accent-text); }

  /* status pills (agenda) */
  .pill { display: inline-flex; align-items: center; gap: 5px; padding: 3px 9px; border-radius: 99px; font-size: 12px; font-weight: 500; background: var(--surface-2); color: var(--ink-2); width: max-content; }
  .pill .pdot { width: 5px; height: 5px; border-radius: 50%; background: var(--mute-2); }
  .pill.screen { background: var(--accent-tint); color: var(--accent-text); }
  .pill.screen .pdot { background: var(--accent); }
  .pill.interview { background: var(--warm-tint); color: var(--warm-text); }
  .pill.interview .pdot { background: var(--warm); }
  .pill.offer { background: var(--positive-tint); color: var(--positive-text); }
  .pill.offer .pdot { background: var(--positive); }

  /* ── RIGHT: pulse ── */
  .pulse-stage {
    overflow-y: auto;
    background:
      radial-gradient(72% 50% at 50% 0%, oklch(0.97 0.028 258), transparent 70%),
      var(--surface);
    display: flex; flex-direction: column; justify-content: flex-start;
    padding: 40px 40px 46px;
  }
  .pulse-tag { font-size: 11px; font-weight: 600; letter-spacing: 0.08em; text-transform: uppercase; color: var(--warm-text); display: inline-flex; align-items: center; gap: 7px; margin-bottom: 20px; }
  .pulse-tag .d { width: 6px; height: 6px; border-radius: 50%; background: var(--warm); box-shadow: 0 0 0 3px var(--warm-tint); }

  .pulse-stats { display: grid; grid-template-columns: repeat(3, 1fr); background: var(--card); border: 1px solid var(--rule); border-radius: 14px; box-shadow: 0 12px 30px -16px rgba(40,40,90,0.22); margin-bottom: 30px; overflow: hidden; }
  .pulse-stats .st { padding: 18px 16px 16px; text-align: center; position: relative; cursor: pointer; transition: background .12s; }
  .pulse-stats .st:hover { background: var(--surface-2); }
  .pulse-stats .st + .st::before { content: ""; position: absolute; left: 0; top: 16px; bottom: 16px; width: 1px; background: var(--rule); }
  .pulse-stats .num { display: block; font-size: 34px; font-weight: 400; letter-spacing: -0.03em; line-height: 1; font-variant-numeric: tabular-nums; }
  .pulse-stats .st.warn .num { color: var(--warm-text); }
  .pulse-stats .lbl { display: block; font-size: 11.5px; color: var(--mute); margin-top: 7px; }

  .pulse-sec { display: flex; align-items: baseline; justify-content: space-between; margin-bottom: 8px; }
  .pulse-sec .t { font-size: 11.5px; font-weight: 600; letter-spacing: 0.06em; text-transform: uppercase; color: var(--mute-2); }
  .pulse-sec .c { font-size: 11px; color: var(--mute-2); }

  .pulse-list { display: flex; flex-direction: column; }
  .pulse-empty { font-size: 12.5px; color: var(--mute); padding: 14px 4px; border-top: 1px solid var(--rule); }
  .pulse-row { display: grid; grid-template-columns: 30px 1fr auto auto; gap: 12px; align-items: center; padding: 12px 4px; border-top: 1px solid var(--rule); cursor: pointer; border-radius: 8px; transition: background .12s; }
  .pulse-row:hover { background: var(--surface-2); }
  .row-logo { width: 30px; height: 30px; border-radius: 8px; background: var(--surface-2); object-fit: contain; padding: 4px; }
  .row-logo.letter { display: grid; place-items: center; padding: 0; color: var(--ink-2); font-size: 12px; font-weight: 600; }
  .pulse-row .wx { line-height: 1.3; min-width: 0; }
  .pulse-row .wx b { font-size: 13.5px; font-weight: 500; }
  .pulse-row .wx small { display: block; font-size: 12px; color: var(--mute); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  .pulse-row .days { font-family: var(--mono); font-size: 12.5px; color: var(--mute); font-variant-numeric: tabular-nums; }
  .pulse-row.quiet .days { color: var(--warm-text); font-weight: 500; }
  .pulse-row .ok { width: 28px; display: flex; justify-content: center; }
  .pulse-row .okdot { width: 7px; height: 7px; border-radius: 50%; background: var(--positive); box-shadow: 0 0 0 3px var(--positive-tint); }
  .pulse-row .okdot.warn { background: var(--warm); box-shadow: 0 0 0 3px var(--warm-tint); }

  .tasks { margin-top: 0; }
  .tasks .pulse-sec { margin-bottom: 10px; }
  .waiting-sec { margin-top: 32px; }
  .task { display: grid; grid-template-columns: 22px 1fr auto; gap: 12px; align-items: center; padding: 11px 4px; border-top: 1px solid var(--rule); cursor: pointer; border-radius: 8px; transition: background .12s; }
  .task:hover { background: var(--surface-2); }
  .task .box { width: 18px; height: 18px; border-radius: 6px; border: 1.5px solid var(--rule-strong); background: var(--card); flex-shrink: 0; position: relative; }
  .task.done .box { background: var(--positive); border-color: var(--positive); }
  .task.done .box::after { content: ""; position: absolute; left: 5px; top: 2px; width: 4px; height: 8px; border: solid #fff; border-width: 0 1.6px 1.6px 0; transform: rotate(42deg); }
  .task .tx { line-height: 1.35; min-width: 0; }
  .task .tx b { font-size: 13.5px; font-weight: 500; }
  .task.done .tx b { color: var(--mute); text-decoration: line-through; }
  .task .tx small { display: block; font-size: 12px; color: var(--mute); }
  .task .due { font-family: var(--mono); font-size: 11.5px; color: var(--mute); white-space: nowrap; }
  .task .due.hot { color: var(--warm-text); }
  .tasks-note { font-size: 11px; color: var(--mute-2); margin-top: 12px; padding-left: 4px; }

  .pulse-foot { display: flex; align-items: center; gap: 13px; margin-top: 24px; padding: 13px 14px; background: var(--card); border: 1px solid var(--rule); border-radius: 12px; box-shadow: var(--sh-1); }
  .pulse-foot .fic { flex-shrink: 0; width: 30px; height: 30px; border-radius: 8px; display: flex; align-items: center; justify-content: center; background: var(--warm-tint); color: var(--warm-text); }
  .pulse-foot .ftx { flex: 1; font-size: 12.5px; color: var(--ink-2); line-height: 1.4; }
  .pulse-foot .ftx b { font-weight: 600; }
  .pulse-foot .ftx small { display: block; color: var(--mute); font-size: 11.5px; margin-top: 1px; }
  .pulse-foot .pulse-link { flex-shrink: 0; display: inline-flex; align-items: center; gap: 6px; white-space: nowrap; font-size: 12.5px; font-weight: 500; color: var(--ink-2); background: none; border: none; padding: 6px 4px; cursor: pointer; }
  .pulse-foot .pulse-link:hover { color: var(--ink); }

  /* Mobile — stack the two panes. */
  @media (max-width: 860px) {
    .ob, .ob.ob-swap { grid-template-columns: 1fr; }
    .ob.ob-swap .brief { border-right: 0; border-bottom: 1px solid var(--rule); }
    .search { display: none; }
    .brief-in { padding: 28px 22px 40px; }
    .pulse-stage { padding: 28px 22px 40px; }
  }
</style>
