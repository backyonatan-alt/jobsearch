<script>
  // Variant A — "Faithful funnel"
  // True trapezoid funnel chart (the Mixpanel canonical view).
  // Below: source-breakdown table so you can see WHICH source is converting.
  const counts = { applied: 30, screen: 12, interview: 5, offer: 2 };
  const stages = [
    { key: 'applied',   label: 'Applied',   n: counts.applied,   tone: 'mute'    },
    { key: 'screen',    label: 'Screen',    n: counts.screen,    tone: 'positive'},
    { key: 'interview', label: 'Interview', n: counts.interview, tone: 'accent'  },
    { key: 'offer',     label: 'Offer',     n: counts.offer,     tone: 'warm'    }
  ];

  const sourceRows = [
    { source: 'Referral',  domain: 'linkedin.com',    applied: 8,  screen: 5, interview: 3, offer: 1, tone: 'positive' },
    { source: 'LinkedIn',  domain: 'linkedin.com',    applied: 8,  screen: 2, interview: 0, offer: 0, tone: 'accent'   },
    { source: 'Cold',      domain: 'mail.google.com', applied: 14, screen: 5, interview: 2, offer: 1, tone: 'warm'     }
  ];

  function pct(n, base) { return base ? Math.round((n/base)*100) : null; }

  const overallPct = pct(counts.offer, counts.applied);

  // Trapezoid geometry — each stage's polygon points, scaled to a 800×400 viewBox.
  const W = 800, ROW_H = 80, GAP = 6;
  const base = counts.applied || 1;
  function trap(i) {
    const top    = (1 - stages[i].n   / base) * W / 2;
    const bottom = (1 - (stages[i+1]?.n ?? stages[i].n) / base) * W / 2;
    const y = i * (ROW_H + GAP);
    const yB = y + ROW_H;
    return `${top},${y} ${W - top},${y} ${W - bottom},${yB} ${bottom},${yB}`;
  }
  const toneFill = {
    mute: 'oklch(0.92 0.02 258)',
    positive: 'oklch(0.86 0.1 152)',
    accent: 'oklch(0.78 0.14 258)',
    warm: 'oklch(0.82 0.13 50)'
  };
</script>

<svelte:head><title>Funnel · variant A — Pursuit</title></svelte:head>

<div class="frame">
  <aside class="sidebar">
    <div class="brand">
      <svg viewBox="0 0 24 24" width="24" height="24" fill="none" class="brand-mark">
        <circle cx="12" cy="12" r="9.5" stroke="currentColor" stroke-width="1.4" opacity="0.65"/>
        <circle cx="12" cy="12" r="5.5" stroke="currentColor" stroke-width="1.4" opacity="0.9"/>
        <circle cx="17.5" cy="6.5" r="2.6" fill="currentColor"/>
      </svg>
      <span class="brand-name">Pursuit</span>
    </div>
    <a class="nav-item"><span class="dot"></span>Today</a>
    <a class="nav-item"><span class="dot"></span>Board</a>
    <a class="nav-item active"><span class="dot"></span>Funnel</a>
    <div class="sidebar-footer">
      <div class="profile">
        <img class="av" src="https://www.google.com/s2/favicons?sz=64&domain=google.com" alt="" />
        <div class="who">Yonatan<small>back.yonatan@gmail.com</small></div>
      </div>
    </div>
  </aside>

  <section class="main">
    <div class="topbar">
      <div class="crumb"><span class="here">Funnel</span></div>
      <div class="right">
        <button class="btn">Last 90 days ▾</button>
        <img class="user-av" src="https://www.google.com/s2/favicons?sz=64&domain=google.com" alt="" />
      </div>
    </div>

    <div class="body">
      <div class="body-inner">
        <div class="hello">
          <div class="date">Conversion across the pipeline</div>
          <h1>Funnel.</h1>
        </div>

        <!-- KPI ROW: 3 narrative count cards (same component as Today A) -->
        <div class="kpis">
          <div class="kpi tone-accent">
            <span class="ribbon"></span>
            <div class="kpi-lbl">Overall conversion</div>
            <div class="kpi-n">{overallPct}<span class="kpi-unit">%</span></div>
            <div class="kpi-sub">Applied → Offer · {counts.offer} of {counts.applied}</div>
          </div>
          <div class="kpi tone-positive">
            <span class="ribbon"></span>
            <div class="kpi-lbl">Best source</div>
            <div class="kpi-n">Referral</div>
            <div class="kpi-sub">3× the rate of cold applies</div>
          </div>
          <div class="kpi tone-warm">
            <span class="ribbon"></span>
            <div class="kpi-lbl">Where you're stalling</div>
            <div class="kpi-n">Screen</div>
            <div class="kpi-sub">42% pass · usually 50%+</div>
          </div>
        </div>

        <!-- FUNNEL VISUALIZATION — true trapezoid stack -->
        <div class="block">
          <div class="block-hd">
            <h2>Funnel</h2>
            <span class="ai-tag">last 90 days · {counts.applied} apps</span>
          </div>
          <div class="funnel-wrap">
            <svg viewBox="0 0 {W} {stages.length * (ROW_H + GAP)}" class="funnel-svg" preserveAspectRatio="xMidYMid meet">
              {#each stages as s, i}
                <polygon points={trap(i)} fill={toneFill[s.tone]} />
              {/each}
            </svg>
            <!-- Stage labels overlaid as an absolutely-positioned column -->
            <div class="funnel-labels">
              {#each stages as s, i}
                {@const prev = i === 0 ? null : stages[i-1].n}
                {@const stagePct = prev !== null ? pct(s.n, prev) : null}
                <div class="fl-row">
                  <div class="fl-stage">
                    <div class="fl-lbl">{s.label}</div>
                    <div class="fl-meta">{s.n} {s.n === 1 ? 'app' : 'apps'}</div>
                  </div>
                  <div class="fl-conv">
                    {#if i === 0}
                      <span class="fl-base">baseline</span>
                    {:else}
                      <span class="fl-pct">{stagePct}%</span>
                      <span class="fl-of">from {stages[i-1].label}</span>
                    {/if}
                  </div>
                </div>
              {/each}
            </div>
          </div>
        </div>

        <!-- SOURCE BREAKDOWN — same funnel, sliced by source -->
        <div class="block">
          <div class="block-hd">
            <h2>By source</h2>
            <span class="ai-tag">where your apps come from</span>
          </div>
          <div class="src-table">
            <div class="src-head">
              <span>Source</span>
              <span class="num">Applied</span>
              <span class="num">Screen</span>
              <span class="num">Interview</span>
              <span class="num">Offer</span>
              <span class="num">Rate</span>
            </div>
            {#each sourceRows as r}
              {@const rate = pct(r.offer, r.applied)}
              <div class="src-row">
                <span class="src-name">
                  <img class="src-logo" src={`https://www.google.com/s2/favicons?sz=64&domain=${r.domain}`} alt="" />
                  {r.source}
                </span>
                <span class="num">{r.applied}</span>
                <span class="num">{r.screen}</span>
                <span class="num">{r.interview}</span>
                <span class="num">{r.offer}</span>
                <span class={`num rate t-${r.tone}`}>{rate ?? '—'}%</span>
              </div>
            {/each}
          </div>
        </div>

        <p class="footer-link"><a href="/preview/redesign">← back to previews</a></p>
      </div>
    </div>
  </section>
</div>

<style>
  :global(html, body) { background: var(--surface); margin: 0; }
  .frame { display: grid; grid-template-columns: 220px 1fr; min-height: 100vh; font-family: var(--sans); color: var(--ink); }

  /* Sidebar */
  .sidebar { background: var(--surface-2); border-right: 1px solid var(--rule); padding: 18px 14px; display: flex; flex-direction: column; gap: 4px; }
  .brand { display: flex; align-items: center; gap: 10px; padding: 4px 8px 18px; color: var(--accent); }
  .brand-name { font-size: 18px; font-weight: 600; letter-spacing: -0.02em; color: var(--ink); }
  .nav-item { display: flex; align-items: center; gap: 10px; padding: 7px 10px; border-radius: 6px; font-size: 13.5px; color: var(--ink-2); cursor: pointer; }
  .nav-item .dot { width: 14px; height: 14px; border-radius: 3px; background: var(--rule-strong); }
  .nav-item.active { background: var(--card); color: var(--ink); box-shadow: var(--sh-1); }
  .sidebar-footer { margin-top: auto; padding-top: 16px; }
  .profile { display: flex; align-items: center; gap: 10px; padding: 8px; }
  .profile .av { width: 28px; height: 28px; border-radius: 50%; object-fit: cover; background: var(--rule-strong); }
  .profile .who { font-size: 13px; line-height: 1.2; }
  .profile .who small { display: block; font-size: 11.5px; color: var(--mute); }

  /* Topbar */
  .topbar { display: flex; justify-content: space-between; align-items: center; padding: 12px 28px; border-bottom: 1px solid var(--rule); background: var(--surface); }
  .crumb .here { font-weight: 600; font-size: 14px; }
  .right { display: flex; align-items: center; gap: 8px; }
  .btn { background: var(--card); border: 1px solid var(--rule); border-radius: 99px; padding: 6px 13px; font-size: 13px; font-weight: 500; color: var(--ink); cursor: pointer; }
  .user-av { width: 30px; height: 30px; border-radius: 50%; cursor: pointer; margin-left: 8px; border: 1px solid var(--rule); object-fit: cover; }

  .body { padding: 28px; }
  .body-inner { max-width: 1080px; margin: 0 auto; }

  .hello { margin-bottom: 24px; }
  .hello .date { font-size: 13.5px; color: var(--mute); margin-bottom: 6px; }
  .hello h1 { font-size: 30px; font-weight: 600; letter-spacing: -0.025em; margin: 0; }

  /* KPIs */
  .kpis { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; margin-bottom: 24px; }
  .kpi { position: relative; background: var(--card); border: 1px solid var(--rule); border-radius: 14px; padding: 18px 20px; overflow: hidden; box-shadow: var(--sh-1); }
  .ribbon { position: absolute; top: 0; left: 0; right: 0; height: 3px; }
  .kpi.tone-accent   .ribbon { background: var(--accent); }
  .kpi.tone-positive .ribbon { background: var(--positive); }
  .kpi.tone-warm     .ribbon { background: var(--warm); }
  .kpi-lbl { font-size: 13px; color: var(--mute); margin-bottom: 4px; font-weight: 500; }
  .kpi-n   { font-size: 34px; font-weight: 600; letter-spacing: -0.035em; line-height: 1.1; font-feature-settings: "tnum"; }
  .kpi.tone-accent   .kpi-n { color: var(--accent-text); }
  .kpi.tone-positive .kpi-n { color: var(--positive-text); }
  .kpi.tone-warm     .kpi-n { color: var(--warm-text); }
  .kpi-unit { font-size: 18px; margin-left: 2px; opacity: 0.6; }
  .kpi-sub { font-size: 12px; color: var(--mute); margin-top: 8px; padding-top: 8px; border-top: 1px dashed var(--rule); }

  /* Block */
  .block { background: var(--card); border: 1px solid var(--rule); border-radius: 14px; padding: 22px 24px; margin-bottom: 14px; box-shadow: var(--sh-1); }
  .block-hd { display: flex; align-items: center; gap: 10px; margin-bottom: 18px; }
  .block-hd h2 { font-size: 17px; font-weight: 600; margin: 0; letter-spacing: -0.015em; }
  .ai-tag { display: inline-flex; align-items: center; gap: 5px; font-size: 12px; background: var(--accent-tint); color: var(--accent-text); padding: 3px 10px; border-radius: 99px; font-weight: 500; }

  /* Funnel — trapezoid SVG + labels overlay */
  .funnel-wrap { display: grid; grid-template-columns: 1fr 240px; gap: 24px; align-items: stretch; }
  .funnel-svg { width: 100%; height: auto; display: block; }
  .funnel-labels { display: flex; flex-direction: column; gap: 6px; padding-top: 0; }
  .fl-row { flex: 1; display: flex; flex-direction: column; justify-content: center; border-left: 1px solid var(--rule); padding: 0 0 0 16px; }
  .fl-stage { display: flex; align-items: baseline; gap: 10px; }
  .fl-lbl { font-size: 14px; font-weight: 600; }
  .fl-meta { font-size: 12px; color: var(--mute); font-feature-settings: "tnum"; }
  .fl-conv { font-size: 12.5px; color: var(--mute); margin-top: 2px; display: flex; align-items: baseline; gap: 6px; }
  .fl-base { font-style: italic; opacity: 0.7; }
  .fl-pct { font-size: 18px; font-weight: 600; color: var(--ink); font-feature-settings: "tnum"; }

  /* Source table */
  .src-table { background: var(--surface-2); border-radius: 12px; overflow: hidden; }
  .src-head, .src-row {
    display: grid;
    grid-template-columns: 1.6fr repeat(5, 1fr);
    gap: 12px; align-items: center; padding: 12px 18px;
    font-size: 13px;
  }
  .src-head { font-size: 11.5px; color: var(--mute); font-weight: 600; background: var(--surface-2); }
  .src-row { background: var(--card); border-top: 1px solid var(--rule); }
  .src-row:first-of-type { border-top: 0; }
  .num { text-align: right; font-feature-settings: "tnum"; color: var(--ink-2); }
  .src-name { display: flex; align-items: center; gap: 10px; font-weight: 500; }
  .src-logo { width: 18px; height: 18px; border-radius: 4px; background: var(--surface-2); padding: 1px; }
  .rate { font-weight: 600; font-size: 14px; }
  .rate.t-positive { color: var(--positive-text); }
  .rate.t-accent   { color: var(--accent-text); }
  .rate.t-warm     { color: var(--warm-text); }

  .footer-link { margin-top: 30px; font-size: 13px; }
  .footer-link a { color: var(--accent-text); text-decoration: none; }
</style>
