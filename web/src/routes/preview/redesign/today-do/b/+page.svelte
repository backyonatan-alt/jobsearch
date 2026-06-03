<script>
  // Static mockup — Today, NO interview today.
  // Direction B: "Recent-first with smart prompts".
  // Left column: suggested-for-today chip strip + recent-applications hero
  // (richer cards, each with a contextual inline action chip).
  // Right column reuses the real /app Pulse pane content as-is.
  const firstName = 'Yonatan';
  const greeting = 'Good evening';
  const dateLine = 'Wednesday · 3 June 2026';

  const STATUS_LABEL = { applied: 'Applied', screen: 'Screen', interview: 'Interview', offer: 'Offer' };

  // ── LEFT: compact suggested-for-today chips ──
  const chips = [
    { spark: true,  label: 'Prep for Vercel screen' },
    { spark: false, label: 'Follow up on Eleos Health (7d)' },
    { spark: false, label: 'Decide on Linear offer' }
  ];

  // ── LEFT (hero): recent applications — richer cards, contextual chip ──
  const recent = [
    { co: 'Vercel',     role: 'Senior SWE, Edge',          status: 'screen',  short: 'Ve', cls: 'c1', added: '19d', chip: 'Generate prep',     chipSpark: true  },
    { co: 'Figma',      role: 'Senior Engineer, Multiplayer', status: 'screen',  short: 'Fi', cls: 'c3', added: '17d', chip: 'Add hiring manager', chipSpark: false },
    { co: 'Granola',    role: 'Founding Backend Engineer', status: 'applied', short: 'Gr', cls: 'c4', added: '2d',  chip: null },
    { co: 'Perplexity', role: 'Staff Engineer, Search',    status: 'applied', short: 'Pe', cls: 'c2', added: '2d',  chip: null },
    { co: 'OpenAI',     role: 'Member of Technical Staff', status: 'applied', short: 'Op', cls: 'c3', added: '1d',  chip: null }
  ];

  // ── RIGHT: pulse pane (hardcoded sample, mirrors /app) ──
  const tasks = [
    { id: 't-prep',  b: 'Prep 3 questions for the Vercel team', s: 'Senior SWE, Edge · Vercel', due: 'Today', hot: true,  done: false },
    { id: 't-offer', b: 'Decide on the Linear offer',           s: '$210k base · waiting on you', due: 'Soon', hot: true,  done: false },
    { id: 't-quiet', b: 'Follow up on Eleos Health',            s: 'Quiet 7 days · log it once you reach out', due: '7d', hot: false, done: false }
  ];
  const openTaskCount = tasks.filter(t => !t.done).length;

  const advisoryLabel = 'Eleos Health and Ramp';

  const waiting = [
    { id: 'w1', co: 'Notion',       status: 'applied', short: 'No', cls: 'c2', days: 12, stale: true  },
    { id: 'w2', co: 'Ramp',         status: 'applied', short: 'Ra', cls: 'c1', days: 9,  stale: true  },
    { id: 'w3', co: 'Eleos Health', status: 'applied', short: 'El', cls: 'c3', days: 7,  stale: true  },
    { id: 'w4', co: 'Granola',      status: 'applied', short: 'Gr', cls: 'c4', days: 2,  stale: false },
    { id: 'w5', co: 'Perplexity',   status: 'screen',  short: 'Pe', cls: 'c2', days: 2,  stale: false }
  ];

  const activeCount = 17, awaitingCount = 14, quietCount = 8, totalApps = 21;
</script>

<svelte:head><title>Today (no interview) · variant B — Pursuit</title></svelte:head>

<div class="mock-badge">Static mockup</div>
<a class="mock-back" href="/preview/redesign">← previews</a>

<div class="shell">
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
      <button class="btn btn-primary">New application <span class="kbd">⌘N</span></button>
    </div>
  </div>

  <div class="ob ob-swap">
    <!-- ══ LEFT — Direction B: recent-first with smart prompts ══ -->
    <div class="brief">
      <div class="brief-in">
        <div class="brief-date">{dateLine}</div>
        <div class="brief-head">
          <h1>{greeting}, <b>{firstName}.</b></h1>
          <div class="brief-stats">
            <button class="bstat" title="Applications still active — applied through offer">
              <span class="bstat-n">{activeCount}</span>
              <span class="bstat-l">In progress</span>
            </button>
            <button class="bstat" title="Applied and waiting to hear back">
              <span class="bstat-n">{awaitingCount}</span>
              <span class="bstat-l">Awaiting reply</span>
            </button>
            <button class="bstat warn" title="No reply in over a week">
              <span class="bstat-n">{quietCount}</span>
              <span class="bstat-l">Gone quiet</span>
            </button>
          </div>
        </div>

        <p class="lede">Nothing on the calendar today — pick up where you left off.</p>

        <!-- Compact "Suggested for today" chip strip -->
        <div class="suggest-strip">
          <span class="strip-label">Suggested for today</span>
          <div class="chips">
            {#each chips as c}
              <button class="chip">
                {#if c.spark}<span class="chip-ic">{@render Spark(12)}</span>{/if}
                {c.label}
              </button>
            {/each}
          </div>
        </div>

        <!-- Recent applications — HERO of the left column -->
        <div class="kick recent-kick">Recent applications</div>
        <div class="recent">
          {#each recent as r}
            <div class="rcard" role="button" tabindex="0">
              <span class={`row-logo letter ${r.cls}`}>{r.short}</span>
              <span class="rx">
                <b>{r.co}</b>
                <small>{r.role}</small>
              </span>
              <span class="rmeta">
                <span class={`pill ${r.status}`}><span class="pdot"></span>{STATUS_LABEL[r.status]}</span>
                <span class="ago">added {r.added} ago</span>
              </span>
              {#if r.chip}
                <button class="rchip">
                  {#if r.chipSpark}<span class="rchip-ic">{@render Spark(12)}</span>{/if}
                  {r.chip}
                </button>
              {/if}
            </div>
          {/each}
        </div>

        <div class="foot">
          <button class="foot-link">View all {totalApps} applications {@render Arrow()}</button>
        </div>
      </div>
    </div>

    <!-- ══ RIGHT — Where things stand (pulse) — reused as-is ══ -->
    <div class="pulse-stage">
      <div class="pulse-tag"><span class="d"></span>Where things stand</div>

      <div class="tasks">
        <div class="pulse-sec">
          <span class="t">Your move</span>
          <span class="c">{openTaskCount} to do</span>
        </div>
        {#each tasks as t (t.id)}
          <div class={`task ${t.done ? 'done' : ''}`} role="button" tabindex="0">
            <span class="box"></span>
            <span class="tx"><b>{t.b}</b><small>{t.s}</small></span>
            {#if t.due}<span class={`due ${t.hot && !t.done ? 'hot' : ''}`}>{t.due}</span>{/if}
          </div>
        {/each}
        <div class="tasks-note">Personal checklist · stays on this device.</div>
      </div>

      <div class="pulse-foot">
        <span class="fic">{@render Spark(15)}</span>
        <span class="ftx"><b>{advisoryLabel} have gone quiet</b><small>No reply in over a week — it might be a good time to reach out to them directly.</small></span>
        <button class="pulse-link">See both {@render Arrow()}</button>
      </div>

      <div class="pulse-sec waiting-sec">
        <span class="t">Waiting to hear back</span>
        <span class="c">longest first</span>
      </div>
      <div class="pulse-list">
        {#each waiting as w (w.id)}
          <div class={`pulse-row ${w.stale ? 'quiet' : ''}`} role="button" tabindex="0">
            <span class={`row-logo letter ${w.cls}`}>{w.short}</span>
            <span class="wx"><b>{w.co}</b><small>{STATUS_LABEL[w.status]}</small></span>
            <span class="days">{w.days}d</span>
            <span class="ok"><span class={`okdot ${w.stale ? 'warn' : ''}`}></span></span>
          </div>
        {/each}
      </div>
    </div>
  </div>
</div>

{#snippet Spark(s)}
  <svg width={s ?? 13} height={s ?? 13} viewBox="0 0 16 16" fill="currentColor" aria-hidden="true"><path d="M8 1l1.5 4.2L14 7l-4.5 1.8L8 13l-1.5-4.2L2 7l4.5-1.8z"/></svg>
{/snippet}
{#snippet Arrow()}
  <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.6" aria-hidden="true"><path d="M3 8h9M8 4l4 4-4 4" stroke-linecap="round" stroke-linejoin="round"/></svg>
{/snippet}

<style>
  :global(html, body) { background: var(--surface); margin: 0; }

  .shell {
    height: 100vh; display: flex; flex-direction: column; overflow: hidden;
    font-family: var(--sans); color: var(--ink);
  }

  /* Static-mockup chrome */
  .mock-badge {
    position: fixed; top: 10px; right: 12px; z-index: 50;
    font-size: 11px; font-weight: 600; letter-spacing: 0.04em;
    background: var(--ink); color: #fff; padding: 4px 10px; border-radius: 99px;
  }
  .mock-back {
    position: fixed; top: 10px; left: 12px; z-index: 50;
    font-size: 12px; color: var(--mute); text-decoration: none;
    background: var(--card); border: 1px solid var(--rule);
    padding: 4px 10px; border-radius: 99px;
  }
  .mock-back:hover { color: var(--ink); }

  /* ── Two-pane Today (swapped: Brief left, Pulse right) ── */
  .ob {
    flex: 1; min-height: 0;
    display: grid; grid-template-columns: 1.08fr 0.92fr;
    font-family: var(--sans); color: var(--ink);
  }
  .ob.ob-swap .brief { border-right: 1px solid var(--rule); }

  /* ── Topbar ── */
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
  .brief-in { max-width: 660px; padding: 44px 40px 56px; }
  .brief-date { font-size: 13px; color: var(--mute); margin-bottom: 16px; letter-spacing: -0.003em; }
  .brief h1 { font-size: 26px; font-weight: 300; letter-spacing: -0.03em; line-height: 1.12; margin: 0; }
  .brief h1 b { font-weight: 500; }
  .lede { font-size: 14.5px; color: var(--ink-2); line-height: 1.6; margin: 0 0 24px; max-width: 50ch; }

  .brief-head { display: flex; align-items: center; justify-content: space-between; gap: 20px; margin: 0 0 22px; flex-wrap: nowrap; }
  .brief-head h1 { margin: 0; white-space: nowrap; }
  .brief-stats { display: flex; align-items: center; flex-shrink: 0; }
  .brief-stats .bstat { display: flex; flex-direction: column; align-items: flex-start; gap: 2px; cursor: pointer; padding: 0 14px; transition: opacity .12s; background: none; border: 0; }
  .brief-stats .bstat:first-child { padding-left: 0; }
  .brief-stats .bstat:last-child { padding-right: 0; }
  .brief-stats .bstat + .bstat { border-left: 1px solid var(--rule); }
  .brief-stats .bstat:hover { opacity: 0.65; }
  .brief-stats .bstat-n { font-size: 23px; font-weight: 500; line-height: 1; letter-spacing: -0.022em; color: var(--ink); font-variant-numeric: tabular-nums; }
  .brief-stats .bstat-l { font-size: 11px; color: var(--mute); letter-spacing: -0.003em; white-space: nowrap; }
  .brief-stats .bstat.warn .bstat-n { color: var(--warm-text); }

  /* Compact "Suggested for today" chip strip */
  .suggest-strip { margin-bottom: 34px; }
  .strip-label { display: block; font-size: 11.5px; font-weight: 600; letter-spacing: 0.06em; text-transform: uppercase; color: var(--mute-2); margin-bottom: 11px; }
  .chips { display: flex; flex-wrap: wrap; gap: 9px; }
  .chip {
    display: inline-flex; align-items: center; gap: 7px; cursor: pointer;
    background: var(--accent-tint); border: 1px solid transparent; border-radius: 99px;
    padding: 7px 14px; font-size: 12.5px; font-weight: 500; color: var(--accent-text);
    font-family: inherit; transition: border-color .12s, transform .12s;
  }
  .chip:hover { border-color: var(--accent); transform: translateY(-1px); }
  .chip-ic { display: inline-flex; color: var(--accent-text); }

  .kick { font-size: 11.5px; font-weight: 600; letter-spacing: 0.07em; text-transform: uppercase; color: var(--mute-2); margin-bottom: 14px; display: flex; align-items: center; gap: 10px; }
  .kick::after { content: ""; flex: 1; height: 1px; background: var(--rule); }
  .recent-kick { margin-top: 0; }

  /* Recent applications — richer hero cards */
  .recent { display: flex; flex-direction: column; gap: 10px; }
  .rcard {
    display: grid; grid-template-columns: 38px 1fr auto auto; gap: 14px; align-items: center;
    background: var(--card); border: 1px solid var(--rule); border-radius: 13px;
    padding: 15px 17px; cursor: pointer;
    transition: border-color .12s, box-shadow .12s, transform .12s;
  }
  .rcard:hover { border-color: var(--rule-strong); box-shadow: var(--sh-pop); transform: translateY(-1px); }
  .rcard .row-logo { width: 38px; height: 38px; border-radius: 10px; font-size: 13px; }
  .rx { line-height: 1.35; min-width: 0; }
  .rx b { font-size: 14.5px; font-weight: 600; color: var(--ink); }
  .rx small { display: block; font-size: 12.5px; color: var(--mute); margin-top: 2px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  .rmeta { display: flex; flex-direction: column; align-items: flex-end; gap: 5px; white-space: nowrap; }
  .rmeta .ago { font-size: 11.5px; color: var(--mute-2); }
  .rchip {
    flex-shrink: 0; white-space: nowrap; cursor: pointer; font-family: inherit;
    display: inline-flex; align-items: center; gap: 6px;
    background: var(--surface-2); border: 1px solid var(--rule); border-radius: 8px;
    padding: 7px 13px; font-size: 12.5px; font-weight: 500; color: var(--ink-2);
    transition: background .12s, border-color .12s, color .12s;
  }
  .rchip:hover { background: var(--accent-tint); border-color: var(--accent); color: var(--accent-text); }
  .rchip-ic { display: inline-flex; color: var(--accent-text); }

  /* shared letter logo */
  .row-logo { width: 30px; height: 30px; border-radius: 8px; background: var(--surface-2); object-fit: contain; padding: 4px; }
  .row-logo.letter { display: grid; place-items: center; padding: 0; color: var(--ink-2); font-size: 12px; font-weight: 600; }
  .row-logo.letter.c1 { background: oklch(0.94 0.04 258); color: var(--accent-text); }
  .row-logo.letter.c2 { background: oklch(0.95 0.05 152); color: var(--positive-text); }
  .row-logo.letter.c3 { background: oklch(0.95 0.06 50);  color: var(--warm-text); }
  .row-logo.letter.c4 { background: var(--surface-2); color: var(--ink-2); }

  /* status pills */
  .pill { display: inline-flex; align-items: center; gap: 5px; padding: 3px 9px; border-radius: 99px; font-size: 12px; font-weight: 500; background: var(--surface-2); color: var(--ink-2); width: max-content; }
  .pill .pdot { width: 5px; height: 5px; border-radius: 50%; background: var(--mute-2); }
  .pill.screen { background: var(--accent-tint); color: var(--accent-text); }
  .pill.screen .pdot { background: var(--accent); }
  .pill.interview { background: var(--warm-tint); color: var(--warm-text); }
  .pill.interview .pdot { background: var(--warm); }
  .pill.offer { background: var(--positive-tint); color: var(--positive-text); }
  .pill.offer .pdot { background: var(--positive); }

  .foot { margin-top: 26px; display: flex; justify-content: flex-end; }
  .foot-link { background: none; border: none; padding: 4px 0; font-family: inherit; font-size: 12.5px; color: var(--mute); display: inline-flex; align-items: center; gap: 6px; cursor: pointer; transition: color .12s; }
  .foot-link:hover { color: var(--accent-text); }

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

  .pulse-sec { display: flex; align-items: baseline; justify-content: space-between; margin-bottom: 8px; }
  .pulse-sec .t { font-size: 11.5px; font-weight: 600; letter-spacing: 0.06em; text-transform: uppercase; color: var(--mute-2); }
  .pulse-sec .c { font-size: 11px; color: var(--mute-2); }

  .pulse-list { display: flex; flex-direction: column; }
  .pulse-row { display: grid; grid-template-columns: 30px 1fr auto auto; gap: 12px; align-items: center; padding: 12px 4px; border-top: 1px solid var(--rule); cursor: pointer; border-radius: 8px; transition: background .12s; }
  .pulse-row:hover { background: var(--surface-2); }
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
    .rcard { grid-template-columns: 38px 1fr auto; }
    .rchip { grid-column: 2 / -1; justify-self: start; }
  }
</style>
