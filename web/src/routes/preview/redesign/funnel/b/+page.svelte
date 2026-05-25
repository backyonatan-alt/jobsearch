<script>
  // Variant B — "Analytics dashboard"
  // Funnel is one card among many — KPIs, stage bars, source comparison,
  // CV-variant breakdown, time-in-stage. Closer to a real PM dashboard.
  const counts = { applied: 30, screen: 12, interview: 5, offer: 2 };
  const stages = [
    { key: 'applied',   label: 'Applied',   n: counts.applied,   tone: 'mute'    },
    { key: 'screen',    label: 'Screen',    n: counts.screen,    tone: 'positive'},
    { key: 'interview', label: 'Interview', n: counts.interview, tone: 'accent'  },
    { key: 'offer',     label: 'Offer',     n: counts.offer,     tone: 'warm'    }
  ];
  const sourceBars = [
    { name: 'Referral',  domain: 'linkedin.com', n: 8, rate: 38, tone: 'positive' },
    { name: 'LinkedIn',  domain: 'linkedin.com', n: 8, rate: 12, tone: 'accent'   },
    { name: 'Cold',      domain: 'mail.google.com', n: 14, rate: 7,  tone: 'warm' }
  ];
  const cvBars = [
    { name: 'CV v3 — infra',     n: 14, rate: 28, tone: 'positive' },
    { name: 'CV v2 — generalist', n: 10, rate: 18, tone: 'accent'  },
    { name: 'CV v1 — old',        n: 6,  rate: 11, tone: 'warm'    }
  ];
  const timeInStage = [
    { stage: 'Applied',   days: 4,  tone: 'mute'     },
    { stage: 'Screen',    days: 6,  tone: 'positive' },
    { stage: 'Interview', days: 8,  tone: 'accent'   },
    { stage: 'Offer',     days: 3,  tone: 'warm'     }
  ];

  function pct(n, base) { return base ? Math.round((n/base)*100) : null; }
  const overallPct = pct(counts.offer, counts.applied);
  const maxSource = Math.max(...sourceBars.map(s => s.rate));
  const maxCv     = Math.max(...cvBars.map(c => c.rate));
  const maxTime   = Math.max(...timeInStage.map(t => t.days));
</script>

<svelte:head><title>Funnel · variant B — Pursuit</title></svelte:head>

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
          <div class="date">How the pipeline is performing</div>
          <h1>Funnel.</h1>
        </div>

        <!-- KPI ROW -->
        <div class="kpis">
          <div class="kpi tone-accent">
            <span class="ribbon"></span>
            <div class="kpi-lbl">Overall conversion</div>
            <div class="kpi-n">{overallPct}<span class="kpi-unit">%</span></div>
            <div class="kpi-sub">Applied → Offer</div>
          </div>
          <div class="kpi tone-positive">
            <span class="ribbon"></span>
            <div class="kpi-lbl">Best CV variant</div>
            <div class="kpi-n">v3</div>
            <div class="kpi-sub">28% reach screen</div>
          </div>
          <div class="kpi tone-warm">
            <span class="ribbon"></span>
            <div class="kpi-lbl">Avg time to offer</div>
            <div class="kpi-n">21<span class="kpi-unit">d</span></div>
            <div class="kpi-sub">Applied through Offer</div>
          </div>
          <div class="kpi tone-mute">
            <span class="ribbon"></span>
            <div class="kpi-lbl">In flight</div>
            <div class="kpi-n">{counts.applied - counts.offer}</div>
            <div class="kpi-sub">Active applications</div>
          </div>
        </div>

        <!-- FUNNEL — stepped bars, big numbers -->
        <div class="block">
          <div class="block-hd">
            <h2>Conversion funnel</h2>
            <span class="ai-tag">last 90 days · {counts.applied} apps</span>
          </div>
          <div class="funnel-bars">
            {#each stages as s, i}
              {@const width = (s.n / stages[0].n) * 100}
              {@const prev = i === 0 ? null : stages[i-1].n}
              {@const stagePct = prev !== null ? pct(s.n, prev) : null}
              <div class="fb-row">
                <div class="fb-lbl-col">
                  <span class={`fb-dot d-${s.tone}`}></span>
                  <span class="fb-lbl">{s.label}</span>
                </div>
                <div class="fb-track">
                  <div class={`fb-fill f-${s.tone}`} style="width: {width}%"></div>
                  <span class="fb-n">{s.n}</span>
                </div>
                <div class="fb-pct">
                  {#if stagePct === null}
                    <span class="fb-base">100%</span>
                  {:else}
                    <span class={`fb-rate rate-${stagePct >= 50 ? 'good' : stagePct >= 25 ? 'mid' : 'low'}`}>{stagePct}%</span>
                    <span class="fb-of">of {stages[i-1].label}</span>
                  {/if}
                </div>
              </div>
            {/each}
          </div>
        </div>

        <!-- TWO COLUMN: Source + CV -->
        <div class="two-col">
          <div class="block">
            <div class="block-hd">
              <h2>Conversion by source</h2>
              <span class="ai-tag">Applied → Offer</span>
            </div>
            <div class="bar-list">
              {#each sourceBars as s}
                <div class="bar-row">
                  <div class="bar-lbl">
                    <img class="bar-logo" src={`https://www.google.com/s2/favicons?sz=64&domain=${s.domain}`} alt="" />
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
          </div>

          <div class="block">
            <div class="block-hd">
              <h2>Conversion by CV variant</h2>
              <span class="ai-tag">Applied → Screen</span>
            </div>
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
                <div class="tis-n">{t.days}<span class="tis-unit">d</span></div>
                <div class="tis-bar"><div class="tis-fill" style="width: {(t.days/maxTime)*100}%"></div></div>
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
  .ai-tag { display: inline-flex; align-items: center; gap: 5px; font-size: 12px; background: var(--accent-tint); color: var(--accent-text); padding: 3px 10px; border-radius: 99px; font-weight: 500; }

  /* Funnel bars */
  .funnel-bars { display: flex; flex-direction: column; gap: 14px; }
  .fb-row { display: grid; grid-template-columns: 130px 1fr 140px; gap: 16px; align-items: center; }
  .fb-lbl-col { display: flex; align-items: center; gap: 8px; }
  .fb-dot { width: 9px; height: 9px; border-radius: 50%; background: var(--mute-2); }
  .fb-dot.d-mute     { background: var(--mute-2); }
  .fb-dot.d-positive { background: var(--positive); }
  .fb-dot.d-accent   { background: var(--accent); }
  .fb-dot.d-warm     { background: var(--warm); }
  .fb-lbl { font-size: 13.5px; font-weight: 600; }
  .fb-track { position: relative; height: 32px; background: var(--surface-2); border-radius: 8px; overflow: hidden; }
  .fb-fill { height: 100%; border-radius: 8px; transition: width 200ms ease; }
  .fb-fill.f-mute     { background: oklch(0.82 0.02 258); }
  .fb-fill.f-positive { background: var(--positive); }
  .fb-fill.f-accent   { background: var(--accent); }
  .fb-fill.f-warm     { background: var(--warm); }
  .fb-n { position: absolute; top: 50%; left: 12px; transform: translateY(-50%); font-size: 13px; font-weight: 600; color: white; mix-blend-mode: difference; font-feature-settings: "tnum"; }
  .fb-pct { display: flex; align-items: baseline; gap: 6px; }
  .fb-rate { font-size: 20px; font-weight: 600; font-feature-settings: "tnum"; }
  .rate-good { color: var(--positive-text); }
  .rate-mid  { color: var(--warm-text); }
  .rate-low  { color: var(--danger-text); }
  .fb-base { font-size: 14px; color: var(--mute); font-weight: 500; }
  .fb-of { font-size: 12px; color: var(--mute); }

  /* Two-column */
  .two-col { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; margin-bottom: 14px; }

  /* Generic bar list */
  .bar-list { display: flex; flex-direction: column; gap: 12px; }
  .bar-row { display: grid; grid-template-columns: 1fr 1fr 50px; gap: 12px; align-items: center; }
  .bar-lbl { display: flex; align-items: center; gap: 8px; font-size: 13px; font-weight: 500; }
  .bar-logo { width: 16px; height: 16px; border-radius: 4px; background: var(--surface-2); padding: 1px; }
  .bar-n { font-size: 11.5px; color: var(--mute); font-feature-settings: "tnum"; }
  .bar-track { height: 14px; background: var(--surface-2); border-radius: 999px; overflow: hidden; }
  .bar-fill { height: 100%; border-radius: 999px; }
  .bar-fill.f-positive { background: var(--positive); }
  .bar-fill.f-accent   { background: var(--accent); }
  .bar-fill.f-warm     { background: var(--warm); }
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

  .footer-link { margin-top: 30px; font-size: 13px; }
  .footer-link a { color: var(--accent-text); text-decoration: none; }
</style>
