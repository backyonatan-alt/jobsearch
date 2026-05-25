<script>
  // 3 funnel chart treatments to choose from. Same data, same dashboard
  // wrapper as funnel/b — only the funnel chart's colors and number
  // typography differ. Pick one.
  const counts = { applied: 30, screen: 12, interview: 5, offer: 2 };
  const stages = [
    { key: 'applied',   label: 'Applied',   n: counts.applied,   tone: 'mute'    },
    { key: 'screen',    label: 'Screen',    n: counts.screen,    tone: 'positive'},
    { key: 'interview', label: 'Interview', n: counts.interview, tone: 'accent'  },
    { key: 'offer',     label: 'Offer',     n: counts.offer,     tone: 'warm'    }
  ];
  function pct(n, base) { return base ? Math.round((n/base)*100) : null; }
  const overallPct = pct(counts.offer, counts.applied);
</script>

<svelte:head><title>Funnel · 3 chart options — Pursuit</title></svelte:head>

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
      <div class="crumb"><span class="here">Funnel · 3 chart options</span></div>
      <div class="right">
        <button class="btn">Last 90 days ▾</button>
        <img class="user-av" src="https://www.google.com/s2/favicons?sz=64&domain=google.com" alt="" />
      </div>
    </div>

    <div class="body">
      <div class="body-inner">
        <div class="hello">
          <div class="date">Pick one for the funnel chart</div>
          <h1>Funnel · 3 options.</h1>
          <p class="note">Same data in all three. Only the chart's colors and number styling change. Everything else around it (KPIs, source bars, time-in-stage) stays as variant B.</p>
        </div>

        <!-- ============= OPTION 1: Monochromatic accent ============= -->
        <div class="block">
          <div class="block-hd">
            <span class="opt-tag">Option 1</span>
            <h2>Monochromatic blue · numbers inside bars</h2>
          </div>
          <div class="funnel-bars opt-1">
            {#each stages as s, i}
              {@const width = (s.n / stages[0].n) * 100}
              {@const prev = i === 0 ? null : stages[i-1].n}
              {@const stagePct = prev !== null ? pct(s.n, prev) : null}
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

        <!-- ============= OPTION 2: Semantic saturated colors ============= -->
        <div class="block">
          <div class="block-hd">
            <span class="opt-tag">Option 2</span>
            <h2>Semantic colors · big confident numbers</h2>
          </div>
          <div class="funnel-bars opt-2">
            {#each stages as s, i}
              {@const width = (s.n / stages[0].n) * 100}
              {@const prev = i === 0 ? null : stages[i-1].n}
              {@const stagePct = prev !== null ? pct(s.n, prev) : null}
              {@const lost = prev !== null ? prev - s.n : null}
              <div class="fb-row">
                <div class="fb-lbl-col">
                  <span class={`fb-dot d-${s.tone}`}></span>
                  <span class="fb-lbl">{s.label}</span>
                </div>
                <div class="fb-track">
                  <div class={`fb-fill f-${s.tone}`} style={`width: ${width}%`}></div>
                  <span class="fb-n-big">{s.n}</span>
                </div>
                <div class="fb-pct">
                  {#if stagePct === null}
                    <span class="fb-base">100%</span>
                  {:else}
                    <span class={`fb-rate-num rate-${stagePct >= 50 ? 'good' : stagePct >= 25 ? 'mid' : 'low'}`}>{stagePct}%</span>
                    <span class="fb-of">−{lost} dropped</span>
                  {/if}
                </div>
              </div>
            {/each}
          </div>
        </div>

        <!-- ============= OPTION 3: Tinted pastel + huge outside numbers ============= -->
        <div class="block">
          <div class="block-hd">
            <span class="opt-tag">Option 3</span>
            <h2>Tinted pastels · huge numbers outside the bar</h2>
          </div>
          <div class="funnel-bars opt-3">
            {#each stages as s, i}
              {@const width = (s.n / stages[0].n) * 100}
              {@const prev = i === 0 ? null : stages[i-1].n}
              {@const stagePct = prev !== null ? pct(s.n, prev) : null}
              <div class="fb-row">
                <div class={`fb-huge-n t-${s.tone}`}>{s.n}</div>
                <div class="fb-mid">
                  <div class="fb-lbl-row">
                    <span class="fb-lbl">{s.label}</span>
                    {#if stagePct !== null}
                      <span class={`fb-rate-chip t-${s.tone}`}>
                        <svg width="9" height="9" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 4l4 4 4-4"/></svg>
                        {stagePct}%
                      </span>
                    {/if}
                  </div>
                  <div class="fb-track">
                    <div class={`fb-fill f-tint t-${s.tone}`} style={`width: ${width}%`}></div>
                  </div>
                </div>
              </div>
            {/each}
          </div>
        </div>

        <div class="block summary">
          <div class="block-hd">
            <h2>How to think about them</h2>
          </div>
          <ul class="cheat">
            <li><b>Option 1</b> — calm and unified. The funnel feels like a single object; no color story between stages.</li>
            <li><b>Option 2</b> — encodes meaning in color (green = progressing, blue = live, orange = decision). Highest information density.</li>
            <li><b>Option 3</b> — the numbers themselves are the visual. Most "designer-y", least dense.</li>
          </ul>
        </div>

        <p class="footer-link"><a href="/preview/redesign">← back to previews</a></p>
      </div>
    </div>
  </section>
</div>

<style>
  :global(html, body) { background: var(--surface); margin: 0; }
  .frame { display: grid; grid-template-columns: 220px 1fr; min-height: 100vh; font-family: var(--sans); color: var(--ink); }

  /* Sidebar (shared with B) */
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
  .hello .note { font-size: 13.5px; color: var(--mute); margin: 8px 0 0; line-height: 1.55; max-width: 70ch; }

  /* Block */
  .block { background: var(--card); border: 1px solid var(--rule); border-radius: 14px; padding: 22px 24px; margin-bottom: 14px; box-shadow: var(--sh-1); }
  .block-hd { display: flex; align-items: center; gap: 10px; margin-bottom: 18px; }
  .block-hd h2 { font-size: 17px; font-weight: 600; margin: 0; letter-spacing: -0.015em; }
  .opt-tag { font-size: 11px; font-weight: 700; color: white; background: var(--ink); padding: 3px 9px; border-radius: 99px; letter-spacing: 0.02em; }

  /* Shared funnel layout */
  .funnel-bars { display: flex; flex-direction: column; gap: 14px; }
  .fb-row { display: grid; grid-template-columns: 130px 1fr 170px; gap: 16px; align-items: center; }
  .fb-lbl-col { display: flex; align-items: center; gap: 8px; }
  .fb-lbl { font-size: 13.5px; font-weight: 600; }
  .fb-track { position: relative; height: 36px; background: var(--surface-2); border-radius: 10px; overflow: hidden; }
  .fb-fill { height: 100%; border-radius: 10px; transition: width 240ms ease; }
  .fb-pct { display: flex; align-items: baseline; gap: 8px; }

  /* OPTION 1 — monochromatic accent */
  .opt-1 .fb-n {
    position: absolute; top: 50%; left: 14px; transform: translateY(-50%);
    font-size: 14px; font-weight: 700; color: white;
    font-feature-settings: "tnum"; letter-spacing: -0.01em;
  }
  .opt-1 .fb-rate-pill {
    font-size: 14px; font-weight: 600; padding: 5px 12px; border-radius: 99px;
    background: var(--accent-tint); color: var(--accent-text);
    font-feature-settings: "tnum";
  }
  .opt-1 .fb-base { font-size: 12.5px; color: var(--mute); font-style: italic; }

  /* OPTION 2 — semantic colors + big numbers */
  .opt-2 .fb-dot { width: 10px; height: 10px; border-radius: 50%; }
  .opt-2 .fb-dot.d-mute     { background: var(--mute-2); }
  .opt-2 .fb-dot.d-positive { background: var(--positive); }
  .opt-2 .fb-dot.d-accent   { background: var(--accent); }
  .opt-2 .fb-dot.d-warm     { background: var(--warm); }
  .opt-2 .fb-fill.f-mute     { background: oklch(0.55 0.04 258); }
  .opt-2 .fb-fill.f-positive { background: var(--positive); }
  .opt-2 .fb-fill.f-accent   { background: var(--accent); }
  .opt-2 .fb-fill.f-warm     { background: var(--warm); }
  .opt-2 .fb-n-big {
    position: absolute; top: 50%; left: 16px; transform: translateY(-50%);
    font-size: 18px; font-weight: 700; color: white;
    font-feature-settings: "tnum"; letter-spacing: -0.02em;
    text-shadow: 0 1px 2px rgba(0,0,0,0.15);
  }
  .opt-2 .fb-rate-num { font-size: 22px; font-weight: 700; font-feature-settings: "tnum"; letter-spacing: -0.02em; }
  .opt-2 .rate-good { color: var(--positive-text); }
  .opt-2 .rate-mid  { color: var(--warm-text); }
  .opt-2 .rate-low  { color: var(--danger-text); }
  .opt-2 .fb-of { font-size: 11.5px; color: var(--mute); }
  .opt-2 .fb-base { font-size: 16px; font-weight: 600; color: var(--mute); }

  /* OPTION 3 — huge outside numbers + tinted bars */
  .opt-3 .fb-row { grid-template-columns: 96px 1fr; gap: 18px; align-items: center; }
  .opt-3 .fb-huge-n {
    font-size: 44px; font-weight: 700; letter-spacing: -0.04em;
    line-height: 1; text-align: right;
    font-feature-settings: "tnum";
  }
  .opt-3 .fb-huge-n.t-mute     { color: var(--mute); }
  .opt-3 .fb-huge-n.t-positive { color: var(--positive-text); }
  .opt-3 .fb-huge-n.t-accent   { color: var(--accent-text); }
  .opt-3 .fb-huge-n.t-warm     { color: var(--warm-text); }
  .opt-3 .fb-mid { display: flex; flex-direction: column; gap: 6px; }
  .opt-3 .fb-lbl-row { display: flex; align-items: center; gap: 10px; }
  .opt-3 .fb-track { height: 12px; }
  .opt-3 .fb-fill.f-tint { border-radius: 99px; }
  .opt-3 .fb-fill.f-tint.t-mute     { background: var(--mute-2); }
  .opt-3 .fb-fill.f-tint.t-positive { background: var(--positive); }
  .opt-3 .fb-fill.f-tint.t-accent   { background: var(--accent); }
  .opt-3 .fb-fill.f-tint.t-warm     { background: var(--warm); }
  .opt-3 .fb-rate-chip {
    font-size: 12px; font-weight: 600; padding: 3px 9px; border-radius: 99px;
    display: inline-flex; align-items: center; gap: 4px;
    font-feature-settings: "tnum";
  }
  .opt-3 .fb-rate-chip.t-mute     { background: var(--surface-2); color: var(--mute); }
  .opt-3 .fb-rate-chip.t-positive { background: var(--positive-tint); color: var(--positive-text); }
  .opt-3 .fb-rate-chip.t-accent   { background: var(--accent-tint); color: var(--accent-text); }
  .opt-3 .fb-rate-chip.t-warm     { background: var(--warm-tint); color: var(--warm-text); }

  .summary { background: var(--surface-2); box-shadow: none; }
  .cheat { margin: 0; padding-left: 18px; }
  .cheat li { font-size: 13.5px; color: var(--ink-2); line-height: 1.6; margin-bottom: 4px; }

  .footer-link { margin-top: 24px; font-size: 13px; }
  .footer-link a { color: var(--accent-text); text-decoration: none; }
</style>
