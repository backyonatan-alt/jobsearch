<script>
  import { onMount } from 'svelte';
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api.js';
  import { isPreview, mockApi } from '$lib/preview-mode.js';
  import { faviconUrl, companyDomain, fmtRelativeDate, STATUS_LABEL } from '$lib/app-helpers.js';

  const call = isPreview() ? mockApi : api;

  let app      = $state(null);
  let dossier  = $state(null);
  let loading  = $state(true);
  let generating = $state(false);
  let genError   = $state('');
  let interviewerInput = $state('');

  const id = $derived(page.params.id);

  $effect(() => {
    void id;
    load();
  });

  async function load() {
    loading = true;
    genError = '';
    try {
      const [a, d] = await Promise.allSettled([
        call(`/api/applications/${id}`),
        call(`/api/applications/${id}/dossier`)
      ]);
      app = a.status === 'fulfilled' ? a.value : null;
      dossier = d.status === 'fulfilled' ? d.value : null;
      if (dossier?.interviewer_name) interviewerInput = dossier.interviewer_name;
      else if (app?.hiring_manager_name) interviewerInput = app.hiring_manager_name;
    } finally {
      loading = false;
    }
  }

  async function generate() {
    if (generating) return;
    generating = true;
    genError = '';
    try {
      const d = await call(`/api/applications/${id}/dossier/refresh`, {
        method: 'POST',
        body: JSON.stringify({ interviewer_name: interviewerInput.trim() || undefined })
      });
      dossier = d;
      interviewerInput = d.interviewer_name ?? interviewerInput;
    } catch (e) {
      genError = friendlyErr(e.message);
    } finally {
      generating = false;
    }
  }

  function friendlyErr(msg) {
    const m = String(msg || '');
    if (m.includes('rate_limit_error') || m.includes('429'))
      return 'AI usage limit hit — wait a minute and try again.';
    if (m.includes('http 504') || /\btimeout\b/i.test(m))
      return 'Web search timed out — try again.';
    if (m.includes('http 5') || m.includes('not configured'))
      return 'Something went wrong — try again in a moment.';
    return m || 'Could not generate the interview prep.';
  }

  function back() {
    if (history.length > 1) history.back();
    else goto(`/app/${id}`);
  }

  // Derived helpers
  const c = $derived(dossier?.content ?? null);
  const interviewer = $derived(c?.interviewer ?? null);

  function initialsOf(name) {
    return (name || '').split(/\s+/).filter(Boolean).slice(0, 2).map(s => s[0]).join('').toUpperCase();
  }

  const ivInitials = $derived(
    initialsOf(interviewer?.name ?? dossier?.interviewer_name ?? app?.hiring_manager_name ?? '')
  );
  const ivName = $derived(
    interviewer?.name ?? dossier?.interviewer_name ?? app?.hiring_manager_name ?? 'Your interviewer'
  );

  // Meeting from dossier.meeting (the shape: { starts_at?, ends_at?, when?, duration?, medium?, panel? })
  const meeting = $derived(dossier?.meeting ?? null);

  function fmtWhen(m) {
    if (!m) return '—';
    if (m.starts_at) {
      const d = new Date(m.starts_at);
      const now = new Date();
      const startOfDay = x => new Date(x.getFullYear(), x.getMonth(), x.getDate());
      const days = Math.round((startOfDay(d) - startOfDay(now)) / 86400000);
      const time = d.toLocaleTimeString(undefined, { hour: 'numeric', minute: '2-digit' });
      if (days === 0) return `Today · ${time}`;
      if (days === 1) return `Tomorrow · ${time}`;
      if (days > 1 && days < 7) return `${d.toLocaleDateString(undefined, { weekday: 'long' })} · ${time}`;
      if (days < 0) return `${d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })} · ${time}`;
      return `${d.toLocaleDateString(undefined, { weekday: 'short', month: 'short', day: 'numeric' })} · ${time}`;
    }
    return m.when ?? '—';
  }

  function fmtDuration(m) {
    if (!m) return '—';
    if (m.starts_at && m.ends_at) {
      const mins = Math.round((new Date(m.ends_at) - new Date(m.starts_at)) / 60000);
      if (mins <= 0) return m.duration ?? '—';
      return mins >= 60 && mins % 60 === 0 ? `${mins / 60}h` : `${mins} min`;
    }
    return m.duration ?? '—';
  }

  // Fact-row helpers — show '—' for unknowns
  const factWhen     = $derived(fmtWhen(meeting));
  const factDuration = $derived(fmtDuration(meeting));
  const factMedium   = $derived(meeting?.medium ?? '—');
  const factPanel    = $derived(meeting?.panel ?? '—');

  // Domain for signal source favicons
  function sigDomain(src) {
    if (!src) return '';
    try {
      return new URL(src.startsWith('http') ? src : `https://${src}`).hostname.replace(/^www\./, '');
    } catch { return src; }
  }

  // Generated timestamp
  const generatedAgo = $derived(dossier?.generatedAgo ?? '');
</script>

<svelte:head>
  <title>{app?.company ? `Interview prep — ${app.company}` : 'Interview prep'} — Pursuit</title>
</svelte:head>

<!-- ── In-page topbar crumb ── -->
<div class="pb-topbar">
  <button class="back" onclick={back}>
    <svg width="13" height="13" viewBox="0 0 13 13" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round">
      <path d="M8 2.5L3.5 6.5 8 10.5"/>
    </svg>
    Back
  </button>
  {#if app}
    <span class="crumb-sep">/</span>
    <span class="crumb-here">{app.company}</span>
  {/if}
</div>

<div class="dos">

  {#if loading}
    <div class="loading-state">
      <div class="spinner"></div>
      <p>Loading…</p>
    </div>

  {:else if !app}
    <div class="empty-state">
      <h3>Application not found</h3>
      <p>It may have been deleted. <a href="/app">Back to Today →</a></p>
    </div>

  {:else}
    <!-- ── Header rule ── -->
    <div class="dos-top">
      <span class="ttl">Interview prep</span>
      {#if generatedAgo}
        <span class="gen">
          <span class="sp" aria-hidden="true">
            <!-- Spark mark -->
            <svg width="13" height="13" viewBox="0 0 13 13" fill="none">
              <path d="M6.5 1.5C6.5 4.5 4.5 6.5 1.5 6.5C4.5 6.5 6.5 8.5 6.5 11.5C6.5 8.5 8.5 6.5 11.5 6.5C8.5 6.5 6.5 4.5 6.5 1.5Z" fill="currentColor"/>
            </svg>
          </span>
          Generated by Pursuit · {generatedAgo}
        </span>
      {/if}
    </div>

    {#if !dossier}
      <!-- ── Empty / generate state ── -->
      <div class="generate-wrap">
        <div class="generate-card">
          <div class="gen-icon" aria-hidden="true">
            <svg width="28" height="28" viewBox="0 0 24 24" fill="none">
              <path d="M12 3C12 7.97 8.97 11 4 11C8.97 11 12 14.03 12 19C12 14.03 15.03 11 20 11C15.03 11 12 7.97 12 3Z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/>
            </svg>
          </div>
          {#if generating}
            <h3>Researching {app.company}{interviewerInput ? ` & ${interviewerInput}` : ''}…</h3>
            <p class="gen-sub">Claude is searching the web for recent posts, talks, and the company's current direction. This typically takes 30–60 seconds.</p>
            <div class="big-spinner"></div>
          {:else}
            <h3>Generate interview prep</h3>
            <p class="gen-sub">
              We'll build an AI brief on the person interviewing you — their background, how they tend to interview, what lands well, and smart questions to ask — pulled from public posts, talks, papers, and company news. Add a name below to make it about a specific interviewer.
            </p>
            <div class="gen-row">
              <input
                class="gen-input"
                type="text"
                placeholder="Interviewer name (optional) — e.g. Sarah Chen"
                bind:value={interviewerInput}
                disabled={generating}
                onkeydown={(e) => e.key === 'Enter' && generate()}
              />
              <button class="btn-generate" onclick={generate} disabled={generating}>
                Generate interview prep
              </button>
            </div>
            {#if genError}
              <p class="gen-err">{genError}</p>
            {/if}
          {/if}
        </div>
      </div>

    {:else}
      <!-- ── Two-column grid ── -->
      <div class="dos-grid">

        <!-- ── LEFT RAIL ── -->
        <div class="dos-rail">

          <!-- Person card -->
          <div class="dos-person">
            <span class="iv-av">{ivInitials || '?'}</span>
            <div class="nm">{ivName}</div>
            {#if interviewer?.role}
              <div class="ro">{interviewer.role}</div>
            {/if}
            {#if interviewer?.prior?.length}
              <div class="prior">
                {#each interviewer.prior as p}
                  <span>{p}</span>
                {/each}
              </div>
            {/if}
            {#if interviewer?.links?.length}
              <div class="links">
                {#each interviewer.links as l}
                  <a href={l.href} target="_blank" rel="noopener">{l.label}</a>
                {/each}
              </div>
            {/if}
          </div>

          <!-- Facts card -->
          <div class="dos-facts">
            <div class="f"><span class="l">Company</span><span class="v">{app.company}</span></div>
            <div class="f"><span class="l">Role</span><span class="v">{app.role}</span></div>
            <div class="f"><span class="l">When</span><span class="v">{factWhen}</span></div>
            <div class="f"><span class="l">Duration</span><span class="v">{factDuration}</span></div>
            <div class="f"><span class="l">Where</span><span class="v">{factMedium}</span></div>
            <div class="f"><span class="l">Round</span><span class="v">{factPanel}</span></div>
          </div>

          <!-- Refresh -->
          <button class="dos-refresh" type="button" onclick={generate} disabled={generating}>
            <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M2 6a4 4 0 1 1 1.2 2.8M2 4v2h2"/></svg>
            {generating ? 'Refreshing…' : 'Refresh prep'}
          </button>
        </div>

        <!-- ── MAIN COLUMN ── -->
        <div class="dos-main">

          <h1>Before you meet <b>{ivName}.</b></h1>
          <p class="dek">An AI brief on the person interviewing you — their background, how they tend to interview, what lands well, and smart questions to ask. Generated from public sources.</p>

          <!-- Snapshot -->
          {#if c?.snapshot}
            <div class="dos-sec">
              <div class="kick">
                <span>Snapshot</span>
              </div>
              <p class="lead">{@html c.snapshot}</p>
            </div>
          {/if}

          <!-- How they interview (style) -->
          {#if c?.style}
            <div class="dos-sec">
              <div class="kick">
                <span>How {ivName.split(' ')[0]} interviews</span>
              </div>
              {#if c.style.lead}
                <p class="lead">{c.style.lead}</p>
              {/if}
              {#if c.style.tells?.length}
                <div class="dos-tells">
                  {#each c.style.tells as t}
                    <div class="t">
                      <div class="l">{t.lbl}</div>
                      <div class="v">{t.val}</div>
                    </div>
                  {/each}
                </div>
              {/if}
            </div>
          {/if}

          <!-- Lands & lands flat -->
          {#if c?.lands?.length || c?.avoid?.length}
            <div class="dos-sec">
              <div class="kick">
                <span>Lands &amp; lands flat</span>
              </div>
              <div class="dos-two">
                {#if c.lands?.length}
                  <div class="dos-list good">
                    <div class="h"><span class="dot"></span>What lands</div>
                    <ul>
                      {#each c.lands as item}
                        <li>{item}</li>
                      {/each}
                    </ul>
                  </div>
                {/if}
                {#if c.avoid?.length}
                  <div class="dos-list bad">
                    <div class="h"><span class="dot"></span>What to avoid</div>
                    <ul>
                      {#each c.avoid as item}
                        <li>{item}</li>
                      {/each}
                    </ul>
                  </div>
                {/if}
              </div>
            </div>
          {/if}

          <!-- Recent signals -->
          {#if c?.signals?.length}
            <div class="dos-sec">
              <div class="kick">
                <span>Recent signals</span>
              </div>
              <div class="dos-signals">
                {#each c.signals as s}
                  <div class="dos-sig">
                    <div class="when">
                      {s.date ?? ''}
                      {#if s.kind}<span class="tg">{s.kind}</span>{/if}
                    </div>
                    <div>
                      <div class="body">{s.body}</div>
                      {#if s.source}
                        <div class="src">
                          {#if sigDomain(s.source)}
                            <img
                              class="sig-favicon"
                              src={`https://www.google.com/s2/favicons?sz=32&domain=${sigDomain(s.source)}`}
                              alt=""
                              width="12"
                              height="12"
                            />
                          {/if}
                          {s.source}
                        </div>
                      {/if}
                    </div>
                  </div>
                {/each}
              </div>
            </div>
          {/if}

          <!-- Questions worth asking -->
          {#if c?.questions?.length}
            <div class="dos-sec">
              <div class="kick">
                <span>Questions worth asking</span>
              </div>
              {#each c.questions as q}
                <div class="dos-q">
                  <div class="q">"{q.q}"</div>
                  {#if q.why}
                    <div class="why">
                      <span class="sp" aria-hidden="true">
                        <svg width="12" height="12" viewBox="0 0 13 13" fill="none">
                          <path d="M6.5 1.5C6.5 4.5 4.5 6.5 1.5 6.5C4.5 6.5 6.5 8.5 6.5 11.5C6.5 8.5 8.5 6.5 11.5 6.5C8.5 6.5 6.5 4.5 6.5 1.5Z" fill="currentColor"/>
                        </svg>
                      </span>
                      {q.why}
                    </div>
                  {/if}
                </div>
              {/each}
            </div>
          {/if}

          <div class="disclaimer">
            Synthesised from public posts, talks, and papers · {generatedAgo ? `refreshed ${generatedAgo}` : 'just generated'} · always verify before you walk in
          </div>

          {#if genError}
            <p class="gen-err" style="margin-top: 16px">{genError}</p>
          {/if}
        </div>
      </div>
    {/if}
  {/if}
</div>

<style>
  /* ── Layout ─────────────────────────────────────────────── */
  .pb-topbar {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 14px 40px;
    border-bottom: 1px solid var(--rule);
    background: var(--surface);
  }
  .back {
    display: inline-flex;
    align-items: center;
    gap: 5px;
    font-size: 13px;
    color: var(--mute);
    background: none;
    border: none;
    padding: 0;
    cursor: pointer;
    font-family: inherit;
  }
  .back:hover { color: var(--ink); }
  .crumb-sep { color: var(--mute-2); font-size: 13px; }
  .crumb-here { font-size: 13px; font-weight: 600; color: var(--ink); }

  .dos {
    max-width: 1120px;
    margin: 0 auto;
    padding: 0 40px 72px;
  }

  /* ── Header rule ─────────────────────────────────────────── */
  .dos-top {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 22px 0 18px;
    border-bottom: 1.5px solid var(--ink);
    margin-bottom: 30px;
  }
  .dos-top .ttl {
    font-size: 13px;
    font-weight: 600;
    display: inline-flex;
    align-items: center;
    gap: 9px;
  }
  .dos-top .gen {
    font-family: inherit;
    font-size: 12px;
    color: var(--mute);
    display: inline-flex;
    align-items: center;
    gap: 7px;
  }
  .dos-top .gen .sp {
    color: var(--accent);
    display: inline-flex;
    align-items: center;
  }

  /* ── Two-column grid ─────────────────────────────────────── */
  .dos-grid {
    display: grid;
    grid-template-columns: 340px 1fr;
    gap: 44px;
    align-items: start;
  }

  /* ── Left rail ───────────────────────────────────────────── */
  .dos-rail {
    position: sticky;
    top: 18px;
  }

  .dos-person {
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 16px;
    padding: 24px;
    box-shadow: var(--sh-pop);
    text-align: center;
    margin-bottom: 16px;
  }
  .dos-person .iv-av {
    width: 76px;
    height: 76px;
    border-radius: 20px;
    font-size: 26px;
    margin: 0 auto 16px;
    background: linear-gradient(155deg, oklch(0.6 0.16 30), oklch(0.46 0.17 32));
    color: #fff;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 600;
    box-shadow: 0 10px 24px -10px oklch(0.5 0.16 30 / 0.6);
  }
  .dos-person .nm {
    font-size: 21px;
    font-weight: 500;
    letter-spacing: -0.02em;
  }
  .dos-person .ro {
    font-size: 13.5px;
    color: var(--mute);
    margin-top: 4px;
  }
  .dos-person .prior {
    display: flex;
    flex-direction: column;
    gap: 6px;
    margin-top: 16px;
  }
  .dos-person .prior span {
    font-size: 12px;
    color: var(--mute);
  }
  .dos-person .links {
    display: flex;
    flex-wrap: wrap;
    gap: 7px;
    justify-content: center;
    margin-top: 18px;
  }
  .dos-person .links a {
    font-size: 11.5px;
    color: var(--accent-text);
    text-decoration: none;
    border: 1px solid var(--rule);
    border-radius: 999px;
    padding: 5px 11px;
    cursor: pointer;
    transition: background 120ms, border-color 120ms;
  }
  .dos-person .links a:hover {
    background: var(--accent-tint);
    border-color: var(--accent-tint-2);
  }

  .dos-facts {
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 14px;
    padding: 6px 16px;
    box-shadow: var(--sh-1);
  }
  .dos-facts .f {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 0;
    border-top: 1px solid var(--rule);
    font-size: 13px;
  }
  .dos-facts .f:first-child { border-top: none; }
  .dos-facts .f .l { color: var(--mute); }
  .dos-facts .f .v { font-weight: 500; max-width: 180px; text-align: right; }

  .dos-cta {
    display: flex;
    width: 100%;
    align-items: center;
    justify-content: center;
    gap: 8px;
    background: var(--ink);
    color: #fff;
    border: none;
    border-radius: 11px;
    padding: 13px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    margin-top: 16px;
    font-family: inherit;
    transition: background 120ms;
  }
  .dos-cta:hover { background: var(--ink-2); }

  .dos-refresh {
    display: flex;
    width: 100%;
    align-items: center;
    justify-content: center;
    gap: 7px;
    background: none;
    color: var(--mute);
    border: 1px solid var(--rule);
    border-radius: 11px;
    padding: 10px;
    font-size: 12.5px;
    font-weight: 500;
    cursor: pointer;
    margin-top: 8px;
    font-family: inherit;
    transition: color 120ms, border-color 120ms, background 120ms;
  }
  .dos-refresh:hover:not(:disabled) {
    color: var(--ink);
    border-color: var(--rule-strong);
    background: var(--surface-2);
  }
  .dos-refresh:disabled { opacity: 0.5; cursor: default; }

  /* ── Main column ─────────────────────────────────────────── */
  .dos-main h1 {
    font-size: 30px;
    font-weight: 300;
    letter-spacing: -0.03em;
    line-height: 1.1;
    margin: 0 0 6px;
  }
  .dos-main h1 b { font-weight: 500; }
  .dos-main .dek {
    font-size: 14.5px;
    color: var(--mute);
    margin: 0 0 30px;
  }

  /* ── Sections ────────────────────────────────────────────── */
  .dos-sec { margin-bottom: 34px; }
  .dos-sec > .kick {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 14px;
    font-size: 11.5px;
    font-weight: 600;
    letter-spacing: 0.07em;
    text-transform: uppercase;
    color: var(--mute-2);
  }
  .dos-sec > .kick::after {
    content: "";
    flex: 1;
    height: 1px;
    background: var(--rule);
  }
  .dos-sec .lead {
    font-size: 15px;
    line-height: 1.65;
    color: var(--ink-2);
    margin: 0;
    max-width: 64ch;
  }

  /* Style tells */
  .dos-tells {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 12px;
    margin-top: 18px;
  }
  .dos-tells .t {
    background: var(--surface-2);
    border-radius: 12px;
    padding: 15px 16px;
  }
  .dos-tells .t .l {
    font-size: 11px;
    font-weight: 600;
    letter-spacing: 0.05em;
    text-transform: uppercase;
    color: var(--mute-2);
    margin-bottom: 7px;
  }
  .dos-tells .t .v {
    font-size: 12.5px;
    line-height: 1.5;
    color: var(--ink-2);
  }

  /* Lands & avoid */
  .dos-two {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 28px;
  }
  .dos-list .h {
    font-size: 13px;
    font-weight: 600;
    display: inline-flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 12px;
  }
  .dos-list .h .dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
  }
  .dos-list.good .h .dot { background: var(--positive); }
  .dos-list.bad  .h .dot { background: var(--danger); }
  .dos-list ul {
    margin: 0;
    padding: 0;
    list-style: none;
    display: flex;
    flex-direction: column;
    gap: 11px;
  }
  .dos-list li {
    font-size: 13.5px;
    line-height: 1.5;
    color: var(--ink-2);
    padding-left: 20px;
    position: relative;
  }
  .dos-list.good li::before {
    content: "";
    position: absolute;
    left: 2px;
    top: 6px;
    width: 6px;
    height: 10px;
    border: solid var(--positive-text);
    border-width: 0 1.8px 1.8px 0;
    transform: rotate(42deg);
  }
  .dos-list.bad li::before {
    content: "×";
    position: absolute;
    left: 2px;
    top: -1px;
    color: var(--danger-text);
    font-size: 15px;
  }

  /* Recent signals */
  .dos-signals {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
  .dos-sig {
    display: grid;
    grid-template-columns: 92px 1fr;
    gap: 16px;
    padding: 16px;
    border: 1px solid var(--rule);
    border-radius: 13px;
    background: var(--card);
    box-shadow: var(--sh-1);
  }
  .dos-sig .when {
    font-family: var(--mono);
    font-size: 11.5px;
    color: var(--mute);
  }
  .dos-sig .when .tg {
    display: inline-block;
    margin-top: 8px;
    font-size: 10px;
    font-weight: 600;
    letter-spacing: 0.05em;
    text-transform: uppercase;
    color: var(--accent-text);
    background: var(--accent-tint);
    border-radius: 6px;
    padding: 3px 8px;
  }
  .dos-sig .body {
    font-size: 13.5px;
    line-height: 1.55;
    color: var(--ink-2);
  }
  .dos-sig .src {
    font-family: var(--mono);
    font-size: 11px;
    color: var(--mute-2);
    margin-top: 8px;
    display: flex;
    align-items: center;
    gap: 5px;
  }
  .sig-favicon {
    border-radius: 2px;
    opacity: 0.7;
  }

  /* Questions */
  .dos-q {
    border: 1px solid var(--rule);
    border-radius: 13px;
    background: var(--card);
    padding: 18px 20px;
    margin-bottom: 12px;
    box-shadow: var(--sh-1);
  }
  .dos-q .q {
    font-size: 14.5px;
    font-weight: 500;
    line-height: 1.45;
    margin-bottom: 8px;
    letter-spacing: -0.01em;
    color: var(--ink);
  }
  .dos-q .why {
    font-size: 12.5px;
    color: var(--mute);
    line-height: 1.5;
    display: flex;
    gap: 8px;
  }
  .dos-q .why .sp {
    color: var(--accent);
    flex-shrink: 0;
    margin-top: 1px;
    display: inline-flex;
  }

  /* Disclaimer */
  .disclaimer {
    margin-top: 18px;
    font-size: 11.5px;
    color: var(--mute);
    padding-top: 14px;
    border-top: 1px dashed var(--rule);
  }

  /* ── Generate / empty state ──────────────────────────────── */
  .generate-wrap {
    display: flex;
    align-items: flex-start;
    justify-content: center;
    padding-top: 40px;
  }
  .generate-card {
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 18px;
    padding: 36px 40px;
    max-width: 560px;
    width: 100%;
    box-shadow: var(--sh-2);
    text-align: center;
  }
  .gen-icon {
    width: 56px;
    height: 56px;
    border-radius: 16px;
    background: var(--accent-tint);
    color: var(--accent-text);
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 20px;
  }
  .generate-card h3 {
    font-size: 20px;
    font-weight: 500;
    letter-spacing: -0.02em;
    margin: 0 0 10px;
  }
  .gen-sub {
    font-size: 14px;
    color: var(--mute);
    line-height: 1.6;
    margin: 0 0 24px;
    max-width: 42ch;
    margin-left: auto;
    margin-right: auto;
  }
  .gen-row {
    display: flex;
    gap: 10px;
    flex-direction: column;
  }
  .gen-input {
    width: 100%;
    padding: 11px 14px;
    font-size: 13.5px;
    font-family: inherit;
    border: 1px solid var(--rule);
    border-radius: 9px;
    background: var(--surface-2);
    color: var(--ink);
    outline: none;
    transition: border-color 120ms;
    box-sizing: border-box;
  }
  .gen-input:focus { border-color: var(--accent); background: var(--card); }
  .gen-input::placeholder { color: var(--mute-2); }
  .btn-generate {
    background: var(--ink);
    color: #fff;
    border: none;
    border-radius: 9px;
    padding: 12px 20px;
    font-size: 14px;
    font-weight: 500;
    font-family: inherit;
    cursor: pointer;
    transition: background 120ms;
    width: 100%;
  }
  .btn-generate:hover:not(:disabled) { background: var(--ink-2); }
  .btn-generate:disabled { opacity: 0.5; cursor: default; }

  .gen-err {
    color: var(--danger-text);
    font-size: 13px;
    margin: 14px 0 0;
    text-align: left;
  }

  /* ── Loading / error states ──────────────────────────────── */
  .loading-state, .empty-state {
    text-align: center;
    padding: 80px 40px;
    color: var(--mute);
  }
  .spinner, .big-spinner {
    width: 24px;
    height: 24px;
    border: 2px solid var(--rule-strong);
    border-top-color: var(--accent);
    border-radius: 50%;
    animation: spin 0.75s linear infinite;
    margin: 0 auto 12px;
  }
  .big-spinner {
    width: 36px;
    height: 36px;
    border-width: 2.5px;
    margin: 24px auto 0;
  }
  @keyframes spin { to { transform: rotate(360deg); } }

  .empty-state h3 { font-size: 18px; font-weight: 500; margin: 0 0 8px; color: var(--ink); }
  .empty-state a { color: var(--accent-text); text-decoration: none; }
</style>
