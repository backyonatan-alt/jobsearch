<script>
  import { onMount } from 'svelte';
  import { api } from '$lib/api.js';
  import { toDisplayApp } from '$lib/app-helpers.js';

  let apps = $state([]);
  let loading = $state(true);

  onMount(async () => {
    try {
      apps = (await api('/api/applications')).map(toDisplayApp);
    } catch (e) {
      if (e.message !== 'unauthorized') console.error(e);
    } finally {
      loading = false;
    }
  });

  // Funnel logic: each downstream stage is a subset of the previous, since
  // an application can't reach "interview" without having gone through "applied"
  // and "screen". We don't have status history, so we infer from current
  // status — anything currently at Screen-or-beyond is "Applied" too.
  // wishlist and rejected/withdrawn are not counted in the funnel itself.
  const counts = $derived.by(() => {
    let applied = 0, screen = 0, interview = 0, offer = 0;
    for (const a of apps) {
      const s = a.status;
      // anything that's progressed at all (not just wishlist) counts as "applied"
      if (['applied', 'screen', 'interview', 'offer', 'rejected', 'withdrawn'].includes(s)) applied++;
      if (['screen', 'interview', 'offer'].includes(s)) screen++;
      if (['interview', 'offer'].includes(s)) interview++;
      if (['offer'].includes(s)) offer++;
    }
    return { applied, screen, interview, offer };
  });

  const stages = $derived([
    { key: 'applied',   label: 'Applied',   n: counts.applied   },
    { key: 'screen',    label: 'Screen',    n: counts.screen    },
    { key: 'interview', label: 'Interview', n: counts.interview },
    { key: 'offer',     label: 'Offer',     n: counts.offer     }
  ]);

  function pct(n, base) {
    if (base === 0) return null;
    return Math.round((n / base) * 100);
  }

  // Conversion between adjacent stages.
  const transitions = $derived(stages.slice(1).map((s, i) => ({
    from:    stages[i].label,
    to:      s.label,
    fromN:   stages[i].n,
    toN:     s.n,
    rate:    pct(s.n, stages[i].n)
  })));

  // Simple rule-based insights — replace with AI weekly review in v0.3.
  const insights = $derived.by(() => {
    const out = [];
    if (counts.applied === 0) {
      out.push({ tone: 'info', title: 'Nothing in flight yet.', body: 'Add your first applications from the Today page and check back as they move.' });
      return out;
    }
    const appliedToScreen = pct(counts.screen, counts.applied);
    if (appliedToScreen !== null) {
      if (appliedToScreen >= 25) {
        out.push({ tone: 'good', title: `Strong reply rate (${appliedToScreen}%).`, body: 'Industry baseline for senior eng applies is around 20%. Whatever CV variant you\'re sending is landing — keep it.' });
      } else if (appliedToScreen < 10 && counts.applied >= 5) {
        out.push({ tone: 'warn', title: `Reply rate is low (${appliedToScreen}%).`, body: 'Try a different CV variant on the next batch, or look at the source — referrals usually convert higher than cold applications.' });
      }
    }
    const screenToInterview = pct(counts.interview, counts.screen);
    if (screenToInterview !== null && counts.screen >= 3 && screenToInterview < 50) {
      out.push({ tone: 'warn', title: `Screens aren't converting to interviews (${screenToInterview}%).`, body: 'Review what happens in screens — the recruiter framing of your background may not match what their loops are looking for.' });
    }
    if (counts.offer > 0) {
      out.push({ tone: 'go', title: `${counts.offer} open offer${counts.offer === 1 ? '' : 's'}.`, body: 'Use them as leverage. The active interview loops are stronger negotiation anchors than any number on a salary site.' });
    }
    if (out.length === 0) {
      out.push({ tone: 'info', title: 'Healthy funnel.', body: 'Nothing jumps out as off. Keep applying and the stage-by-stage rates will get more meaningful as the sample grows.' });
    }
    return out;
  });

  const totalActive = $derived(apps.filter(a => !['rejected', 'withdrawn'].includes(a.status)).length);
  const totalClosed = $derived(apps.filter(a => ['rejected', 'withdrawn'].includes(a.status)).length);
</script>

<svelte:head>
  <title>Funnel — Pursuit</title>
</svelte:head>

<div class="topbar">
  <div class="crumb"><span class="here">Funnel</span></div>
  <div class="right">
    <div class="search">
      <svg class="ico" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5">
        <circle cx="7" cy="7" r="4.5"/><path d="M10.5 10.5L13 13"/>
      </svg>
      <span>Search applications, people…</span>
      <span class="kbd">⌘K</span>
    </div>
  </div>
</div>

<div class="body">
  <div class="body-inner">
    <div class="page-hd">
      <div>
        <div class="date">Conversion across the pipeline</div>
        <h1>Funnel.</h1>
      </div>
      <div class="stats">
        <span>Active <b>{totalActive}</b></span>
        <span>Closed <b>{totalClosed}</b></span>
      </div>
    </div>

    {#if loading}
      <p style="color:var(--mute)">Loading…</p>
    {:else if counts.applied === 0}
      <div class="empty-tab">
        <h3>No data to chart yet</h3>
        <p>Once you add an application or two, the funnel and conversion rates appear here.</p>
      </div>
    {:else}
      <!-- Stage bars -->
      <section class="bars">
        {#each stages as s, i}
          {@const baseN = stages[0].n}
          {@const width = baseN ? (s.n / baseN) * 100 : 0}
          <div class="bar-row">
            <div class="b-label">
              <span class={`pill ${s.key}`}><span class="pdot"></span>{s.label}</span>
            </div>
            <div class="b-track">
              <div class="b-fill" style="width: {width}%"></div>
              <div class="b-num">{s.n}</div>
            </div>
            <div class="b-pct">
              {#if i === 0}
                <span class="baseline">baseline</span>
              {:else}
                {@const rate = pct(s.n, stages[i - 1].n)}
                <span class="rate {rate >= 50 ? 'ok' : rate >= 25 ? 'mid' : 'low'}">
                  {rate ?? '—'}%
                </span>
                <span class="vs">of {stages[i - 1].label}</span>
              {/if}
            </div>
          </div>
        {/each}
      </section>

      <!-- Transition cards -->
      <section class="transitions">
        {#each transitions as t}
          <div class="trans">
            <div class="trans-label">{t.from} <span class="arrow">→</span> {t.to}</div>
            <div class="trans-num">{t.toN}/{t.fromN}</div>
            <div class="trans-rate {t.rate >= 50 ? 'ok' : t.rate >= 25 ? 'mid' : 'low'}">
              {t.rate ?? '—'}%
            </div>
          </div>
        {/each}
      </section>

      <!-- Insights -->
      <section class="insights">
        <h2>What we're noticing</h2>
        <ul>
          {#each insights as ins}
            <li class="ins" data-tone={ins.tone}>
              <span class="ins-dot"></span>
              <div>
                <h3>{ins.title}</h3>
                <p>{ins.body}</p>
              </div>
            </li>
          {/each}
        </ul>
        <p class="ai-coming">
          Rule-based for v0.1. AI weekly review with Claude lands in v0.3 — it'll
          spot patterns these rules miss.
        </p>
      </section>
    {/if}
  </div>
</div>

<style>
  .bars {
    display: flex; flex-direction: column;
    gap: 12px;
    margin: 16px 0 32px;
    border: 1px solid var(--rule);
    background: var(--card);
    border-radius: 12px;
    padding: 20px 24px;
    box-shadow: var(--sh-1);
  }
  .bar-row {
    display: grid;
    grid-template-columns: 110px 1fr 170px;
    gap: 16px;
    align-items: center;
  }
  .b-label { display: flex; }
  .b-track {
    position: relative;
    height: 26px;
    background: var(--surface-2);
    border-radius: 6px;
    overflow: hidden;
  }
  .b-fill {
    height: 100%;
    background: linear-gradient(90deg, var(--accent-tint-2), var(--accent));
    border-radius: 6px;
    transition: width 200ms ease;
  }
  .b-num {
    position: absolute;
    inset: 0;
    display: flex; align-items: center; justify-content: flex-end;
    padding-right: 10px;
    font-family: var(--mono);
    font-size: 12px;
    color: var(--ink);
    font-variant-numeric: tabular-nums;
  }
  .b-pct {
    display: flex; align-items: baseline; gap: 6px;
    font-size: 13px;
  }
  .b-pct .baseline { color: var(--mute-2); font-size: 12px; }
  .b-pct .rate {
    font-weight: 600;
    font-variant-numeric: tabular-nums;
  }
  .b-pct .rate.ok { color: var(--positive-text); }
  .b-pct .rate.mid { color: var(--warm-text); }
  .b-pct .rate.low { color: var(--danger-text); }
  .b-pct .vs { color: var(--mute); font-size: 12px; }

  .transitions {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 8px;
    margin-bottom: 32px;
  }
  .trans {
    border: 1px solid var(--rule);
    background: var(--card);
    border-radius: 10px;
    padding: 12px 16px;
    box-shadow: var(--sh-1);
  }
  .trans-label {
    font-size: 12px;
    color: var(--mute);
    margin-bottom: 6px;
  }
  .trans-label .arrow { color: var(--accent-text); margin: 0 4px; }
  .trans-num {
    font-family: var(--mono);
    font-size: 12px;
    color: var(--mute-2);
    font-variant-numeric: tabular-nums;
  }
  .trans-rate {
    font-size: 22px;
    font-weight: 500;
    letter-spacing: -0.02em;
    font-variant-numeric: tabular-nums;
    margin-top: 2px;
  }
  .trans-rate.ok { color: var(--positive-text); }
  .trans-rate.mid { color: var(--warm-text); }
  .trans-rate.low { color: var(--danger-text); }

  .insights h2 {
    font-size: 15px;
    font-weight: 500;
    letter-spacing: -0.01em;
    margin: 0 0 12px;
  }
  .insights ul { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 8px; }
  .ins {
    display: grid;
    grid-template-columns: 8px 1fr;
    gap: 14px;
    padding: 14px 16px;
    background: var(--card);
    border: 1px solid var(--rule);
    border-left-width: 3px;
    border-radius: 8px;
  }
  .ins[data-tone="good"] { border-left-color: var(--positive); }
  .ins[data-tone="warn"] { border-left-color: var(--warm); }
  .ins[data-tone="go"]   { border-left-color: var(--accent); }
  .ins[data-tone="info"] { border-left-color: var(--rule-strong); }
  .ins-dot {
    width: 8px; height: 8px; border-radius: 999px;
    margin-top: 6px;
  }
  .ins[data-tone="good"] .ins-dot { background: var(--positive); }
  .ins[data-tone="warn"] .ins-dot { background: var(--warm); }
  .ins[data-tone="go"]   .ins-dot { background: var(--accent); }
  .ins[data-tone="info"] .ins-dot { background: var(--rule-strong); }
  .ins h3 {
    font-size: 14px;
    font-weight: 500;
    letter-spacing: -0.005em;
    margin: 0 0 4px;
    color: var(--ink);
  }
  .ins p { margin: 0; font-size: 13px; color: var(--ink-2); line-height: 1.5; }

  .ai-coming {
    margin: 24px 0 0;
    text-align: center;
    font-size: 11.5px;
    color: var(--mute-2);
  }

  @media (max-width: 800px) {
    .bar-row { grid-template-columns: 1fr; gap: 6px; }
    .transitions { grid-template-columns: 1fr; }
  }
</style>
