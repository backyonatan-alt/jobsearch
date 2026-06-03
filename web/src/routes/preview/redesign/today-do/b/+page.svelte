<script>
  // Static mockup — Today, NO interview today.
  // Direction B: "Recent-first with smart prompts".
  // Left column: suggested-for-today chip strip + recent-applications hero
  // (richer cards, each with a contextual inline action chip).
  // Right column reuses the real /app Pulse pane content as-is.
  const firstName = 'Yonatan';
  const greeting = 'Good evening';
  const dateLine = 'Wednesday · 3 June 2026';

  const totalApps = 21;

  // ── Suggested for today (compact glanceable strip) ──
  const suggestions = [
    { label: 'Prep for Vercel screen', spark: true },
    { label: 'Follow up on Eleos Health (7d)', spark: false },
    { label: 'Decide on Linear offer', spark: false }
  ];

  // ── Recent applications (hero of the left column) ──
  const STATUS_LABEL = { applied: 'Applied', screen: 'Screen', interview: 'Interview', offer: 'Offer' };
  const recent = [
    { co: 'Vercel',     role: 'Senior SWE, Edge',          status: 'screen',  added: '19d', chip: 'Generate prep',     cls: 'k1' },
    { co: 'Figma',      role: 'Senior Engineer, Multiplayer', status: 'screen',  added: '17d', chip: 'Add hiring manager', cls: 'k2' },
    { co: 'Granola',    role: 'Founding Backend Engineer',  status: 'applied', added: '2d',  chip: null,                 cls: 'k3' },
    { co: 'Perplexity', role: 'Staff Engineer, Search',     status: 'applied', added: '2d',  chip: null,                 cls: 'k4' },
    { co: 'OpenAI',     role: 'Member of Technical Staff',  status: 'applied', added: '1d',  chip: null,                 cls: 'k5' }
  ];

  // ── Pulse (right pane) — reused as-is from the real Today page ──
  const tasks = [
    { id: 't-prep',  b: 'Prep 3 questions for the Vercel screen', s: 'Senior SWE, Edge · Vercel',                              due: 'Today', hot: true,  done: false },
    { id: 't-offer', b: 'Decide on the Linear offer',             s: 'Staff Engineer · they asked by Friday',                  due: 'Soon',  hot: true,  done: false },
    { id: 't-quiet', b: 'Follow up on Eleos Health',              s: 'Quiet 7 days · log it once you reach out',               due: '7d',    hot: false, done: false }
  ];
  const openTaskCount = tasks.filter(t => !t.done).length;

  const advisoryLabel = 'Eleos Health and Notion';
  const advisoryPlural = true;

  const waiting = [
    { co: 'Eleos Health', short: 'EH', status: 'screen',  days: 7,  stale: true,  cls: 'k6' },
    { co: 'Notion',       short: 'N',  status: 'applied', days: 9,  stale: true,  cls: 'k7' },
    { co: 'Figma',        short: 'F',  status: 'screen',  days: 17, stale: false, cls: 'k2' },
    { co: 'Granola',      short: 'G',  status: 'applied', days: 2,  stale: false, cls: 'k3' },
    { co: 'Perplexity',   short: 'P',  status: 'applied', days: 2,  stale: false, cls: 'k4' }
  ];
</script>

<svelte:head><title>Today (no interview) · variant B — Pursuit</title></svelte:head>

<div class="mockwrap">
  <div class="mocktop">
    <a class="back" href="/preview/redesign">← back to previews</a>
    <span class="mock-badge">Static mockup</span>
  </div>

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
    <!-- ══ LEFT — Recent-first with smart prompts ══════════ -->
    <div class="brief">
      <div class="brief-in">
        <div class="brief-date">{dateLine}</div>
        <div class="brief-head">
          <h1>{greeting}, <b>{firstName}.</b></h1>
          <div class="brief-stats">
            <div class="bstat"><span class="bstat-n">17</span><span class="bstat-l">In progress</span></div>
            <div class="bstat"><span class="bstat-n">14</span><span class="bstat-l">Awaiting reply</span></div>
            <div class="bstat warn"><span class="bstat-n">8</span><span class="bstat-l">Gone quiet</span></div>
          </div>
        </div>

        <p class="lede">Nothing on the calendar today — pick up where you left off.</p>

        <!-- Suggested for today: compact chip strip -->
        <div class="sg-strip">
          <span class="sg-eyebrow">Suggested for today</span>
          <div class="sg-row">
            {#each suggestions as s}
              <button class="sg-chip">
                {#if s.spark}<span class="sp">{@render Spark(12)}</span>{/if}
                {s.label}
              </button>
            {/each}
          </div>
        </div>

        <!-- Recent applications: the hero -->
        <div class="kick">Recent applications</div>
        <div class="recent">
          {#each recent as a}
            <div class="rec-card" role="button" tabindex="0">
              <span class={`rec-logo ${a.cls}`}>{a.co.charAt(0)}</span>
              <div class="rec-text">
                <div class="rec-top">
                  <span class="rec-co">{a.co}</span>
                  <span class={`pill ${a.status}`}><span class="pdot"></span>{STATUS_LABEL[a.status]}</span>
                </div>
                <div class="rec-role">{a.role}</div>
                <div class="rec-added">added {a.added} ago</div>
              </div>
              {#if a.chip}
                <button class="rec-action">{@render Spark(12)} {a.chip}</button>
              {/if}
            </div>
          {/each}
        </div>

        <div class="foot">
          <button class="foot-link">View all {totalApps} applications {@render Arrow()}</button>
        </div>
      </div>
    </div>

    <!-- ══ RIGHT — Where things stand (pulse, reused) ══════ -->
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

      {#if advisoryLabel}
        <div class="pulse-foot">
          <span class="fic">{@render Spark(15)}</span>
          <span class="ftx"><b>{advisoryLabel} {advisoryPlural ? 'have' : 'has'} gone quiet</b><small>No reply in over a week — it might be a good time to reach out to them directly.</small></span>
          <button class="pulse-link">{advisoryPlural ? 'See both' : 'See it'} {@render Arrow()}</button>
        </div>
      {/if}

      <div class="pulse-sec waiting-sec">
        <span class="t">Waiting to hear back</span>
        <span class="c">longest first</span>
      </div>
      <div class="pulse-list">
        {#each waiting as w}
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
  /* Standalone mockup shell — fills the viewport, no app chrome. */
  .mockwrap {
    position: fixed; inset: 0; display: flex; flex-direction: column;
    background: var(--surface); font-family: var(--sans); color: var(--ink);
  }
  .mocktop {
    display: flex; align-items: center; justify-content: space-between;
    padding: 8px 20px; border-bottom: 1px solid var(--rule); background: var(--surface-2);
    flex-shrink: 0;
  }
  .mocktop .back { font-size: 12.5px; color: var(--mute); text-decoration: none; }
  .mocktop .back:hover { color: var(--ink-2); }
  .mock-badge {
    font-size: 11px; font-weight: 600; letter-spacing: 0.04em; color: var(--accent-text);
    background: var(--accent-tint); border-radius: 99px; padding: 3px 10px;
  }

  /* ════ Two-pane Today (Option B, swapped) ════════════════ */
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

  /* ── LEFT ── */
  .brief { overflow-y: auto; }
  .brief-in { max-width: 660px; padding: 44px 40px 56px; }
  .brief-date { font-size: 13px; color: var(--mute); margin-bottom: 16px; letter-spacing: -0.003em; }
  .brief h1 { font-size: 26px; font-weight: 300; letter-spacing: -0.03em; line-height: 1.12; margin: 0; white-space: nowrap; }
  .brief h1 b { font-weight: 500; }
  .lede { font-size: 14.5px; color: var(--ink-2); line-height: 1.6; margin: 0 0 26px; max-width: 50ch; }

  .brief-head { display: flex; align-items: center; justify-content: space-between; gap: 20px; margin: 0 0 22px; flex-wrap: nowrap; }
  .brief-stats { display: flex; align-items: center; flex-shrink: 0; }
  .brief-stats .bstat { display: flex; flex-direction: column; align-items: flex-start; gap: 2px; padding: 0 14px; }
  .brief-stats .bstat:first-child { padding-left: 0; }
  .brief-stats .bstat:last-child { padding-right: 0; }
  .brief-stats .bstat + .bstat { border-left: 1px solid var(--rule); }
  .brief-stats .bstat-n { font-size: 23px; font-weight: 500; line-height: 1; letter-spacing: -0.022em; color: var(--ink); font-variant-numeric: tabular-nums; }
  .brief-stats .bstat-l { font-size: 11px; color: var(--mute); letter-spacing: -0.003em; white-space: nowrap; }
  .brief-stats .bstat.warn .bstat-n { color: var(--warm-text); }

  /* Suggested-for-today chip strip */
  .sg-strip { margin: 0 0 34px; }
  .sg-eyebrow { display: block; font-size: 11.5px; font-weight: 600; letter-spacing: 0.06em; text-transform: uppercase; color: var(--mute-2); margin-bottom: 11px; }
  .sg-row { display: flex; flex-wrap: wrap; gap: 8px; }
  .sg-chip {
    display: inline-flex; align-items: center; gap: 6px;
    font-family: inherit; font-size: 12.5px; font-weight: 500; color: var(--accent-text);
    background: var(--accent-tint); border: 1px solid transparent; border-radius: 99px;
    padding: 6px 13px; cursor: pointer; transition: border-color .12s, background .12s;
  }
  .sg-chip:hover { border-color: var(--accent); }
  .sg-chip .sp { display: inline-flex; }

  /* Recent applications — hero list */
  .kick { font-size: 11.5px; font-weight: 600; letter-spacing: 0.07em; text-transform: uppercase; color: var(--mute-2); margin-bottom: 14px; display: flex; align-items: center; gap: 10px; }
  .kick::after { content: ""; flex: 1; height: 1px; background: var(--rule); }

  .recent { display: flex; flex-direction: column; gap: 10px; }
  .rec-card {
    display: grid; grid-template-columns: 42px 1fr auto; gap: 14px; align-items: center;
    background: var(--card); border: 1px solid var(--rule); border-radius: 13px;
    padding: 15px 16px; cursor: pointer; transition: border-color .12s, box-shadow .12s;
  }
  .rec-card:hover { border-color: var(--rule-strong); box-shadow: var(--sh-1); }
  .rec-logo { width: 42px; height: 42px; border-radius: 11px; display: grid; place-items: center; font-size: 17px; font-weight: 600; color: var(--accent-text); background: var(--accent-tint); }
  .rec-logo.k1 { background: oklch(0.94 0.04 258); color: oklch(0.42 0.12 258); }
  .rec-logo.k2 { background: oklch(0.94 0.045 320); color: oklch(0.42 0.13 320); }
  .rec-logo.k3 { background: oklch(0.94 0.05 95);  color: oklch(0.44 0.11 95); }
  .rec-logo.k4 { background: oklch(0.93 0.05 165); color: oklch(0.42 0.11 165); }
  .rec-logo.k5 { background: oklch(0.94 0.04 30);  color: oklch(0.45 0.13 30); }
  .rec-text { min-width: 0; }
  .rec-top { display: flex; align-items: center; gap: 9px; }
  .rec-co { font-size: 14.5px; font-weight: 600; letter-spacing: -0.01em; }
  .rec-role { font-size: 13px; color: var(--ink-2); margin-top: 2px; }
  .rec-added { font-size: 11.5px; color: var(--mute); margin-top: 4px; }
  .rec-action {
    flex-shrink: 0; display: inline-flex; align-items: center; gap: 6px;
    font-family: inherit; font-size: 12.5px; font-weight: 500; color: var(--ink-2);
    background: var(--surface-2); border: 1px solid var(--rule); border-radius: 8px;
    padding: 7px 11px; cursor: pointer; white-space: nowrap; transition: background .12s, color .12s, border-color .12s;
  }
  .rec-action:hover { background: var(--accent-tint); color: var(--accent-text); border-color: transparent; }

  .foot { margin-top: 26px; display: flex; justify-content: flex-end; }
  .foot-link { background: none; border: none; padding: 4px 0; font-family: inherit; font-size: 12.5px; color: var(--mute); display: inline-flex; align-items: center; gap: 6px; cursor: pointer; transition: color .12s; }
  .foot-link:hover { color: var(--accent-text); }

  /* status pills */
  .pill { display: inline-flex; align-items: center; gap: 5px; padding: 3px 9px; border-radius: 99px; font-size: 12px; font-weight: 500; background: var(--surface-2); color: var(--ink-2); width: max-content; }
  .pill .pdot { width: 5px; height: 5px; border-radius: 50%; background: var(--mute-2); }
  .pill.screen { background: var(--accent-tint); color: var(--accent-text); }
  .pill.screen .pdot { background: var(--accent); }
  .pill.interview { background: var(--warm-tint); color: var(--warm-text); }
  .pill.interview .pdot { background: var(--warm); }
  .pill.offer { background: var(--positive-tint); color: var(--positive-text); }
  .pill.offer .pdot { background: var(--positive); }

  /* ── RIGHT: pulse (reused) ── */
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
  .row-logo { width: 30px; height: 30px; border-radius: 8px; background: var(--surface-2); }
  .row-logo.letter { display: grid; place-items: center; padding: 0; color: var(--ink-2); font-size: 12px; font-weight: 600; }
  .row-logo.k2 { background: oklch(0.94 0.045 320); color: oklch(0.42 0.13 320); }
  .row-logo.k3 { background: oklch(0.94 0.05 95);  color: oklch(0.44 0.11 95); }
  .row-logo.k4 { background: oklch(0.93 0.05 165); color: oklch(0.42 0.11 165); }
  .row-logo.k6 { background: var(--warm-tint); color: var(--warm-text); }
  .row-logo.k7 { background: oklch(0.92 0.01 258); color: var(--ink-2); }
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
    .ob.ob-swap { grid-template-columns: 1fr; }
    .ob.ob-swap .brief { border-right: 0; border-bottom: 1px solid var(--rule); }
    .search { display: none; }
    .brief-in { padding: 28px 22px 40px; }
    .pulse-stage { padding: 28px 22px 40px; }
  }
</style>
