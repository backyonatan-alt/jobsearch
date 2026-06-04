<script>
  import { onMount, onDestroy } from 'svelte';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api.js';

  let { onDone = () => {}, seedDemo = false } = $props();

  // First-run only: clear any stale demo rows, seed fresh ones so the tour's
  // screens look alive, then tell the views to refetch. Cleared again on finish.
  let seeded = false;
  async function seedDemoData() {
    if (!seedDemo || seeded) return;
    seeded = true;
    try {
      await api('/api/me/demo-seed', { method: 'DELETE' }).catch(() => {});
      await api('/api/me/demo-seed', { method: 'POST' });
      window.dispatchEvent(new Event('pursuit:refresh'));
    } catch {}
  }
  async function clearDemoData() {
    if (!seedDemo) return;
    try { await api('/api/me/demo-seed', { method: 'DELETE' }); } catch {}
    window.dispatchEvent(new Event('pursuit:refresh'));
  }
  // finish = clean up demo rows (keeps the user's own added app), then hand off.
  let finishing = false;
  async function finish() {
    if (finishing) return;
    finishing = true;
    await clearDemoData();
    onDone();
  }

  // ── tour script: welcome → 5 highlights → add-application modal ──
  const STEPS = [
    { kind: 'welcome' },
    { sel: '[data-tour="stats"]', view: 'today', icon: 'today', title: 'Your morning briefing', body: "Every day opens here. Pursuit reads your pipeline and tells you what's moving, what's waiting, and what's gone quiet — before you ask.", place: 'bottom' },
    { sel: '[data-tour="prep"]', view: 'today', icon: 'spark', title: 'It preps you for interviews', body: "Ahead of every interview, Pursuit researches who you'll meet and what they care about — so you walk in ready, not guessing.", place: 'bottom' },
    { sel: '[data-tour="board"]', view: 'board', icon: 'board', title: 'Your whole search, one board', body: "Each role is a card moving across stages. Drag a card to update its status — and it turns red when it's stalled for a week.", place: 'bottom' },
    { sel: '[data-tour="funnel"]', view: 'insights', icon: 'insights', title: "See what's actually working", body: 'Your reply rate, how far you get, and which sources convert — so you spend effort where it pays off.', place: 'right' },
    { sel: '[data-tour="new-app"]', view: 'today', icon: 'add', title: 'Add your first application', body: "This is where every role starts. Let's add your first one — it takes seconds.", place: 'bottom' },
    { modal: true, sel: '[data-tour="paste"]', view: 'today', icon: 'add', title: 'Paste a link — Pursuit does the rest', body: 'Drop a job URL or a screenshot here and Pursuit reads the company, role, and location for you. Prefer to type it? The fields are right below.', place: 'right', last: true }
  ];
  const N = STEPS.length - 1; // number of real (non-welcome) stops

  const VIEW_PATH = { today: '/app', board: '/app/board', insights: '/app/funnel' };

  let i = $state(0);
  let target = $state(null);
  let done = $state(false);
  let saving = $state(false);

  const step = $derived(STEPS[i]);

  // Modal form (final stop) — prefilled per the handoff.
  let form = $state({ company: 'Anthropic', role: 'Member of Technical Staff', status: 'Applied', source: 'Referral' });

  // ── positioning helpers (ported from pursuit-tour.jsx) ──
  function place(r, side, cw, ch) {
    const m = 16, vw = window.innerWidth, vh = window.innerHeight;
    let top, left;
    if (side === 'right') { left = r.left + r.width + m; top = r.top + r.height / 2 - ch / 2; }
    else if (side === 'left') { left = r.left - cw - m; top = r.top + r.height / 2 - ch / 2; }
    else if (side === 'top') { left = r.left + r.width / 2 - cw / 2; top = r.top - ch - m; }
    else { left = r.left + r.width / 2 - cw / 2; top = r.top + r.height + m; }
    let actual = side;
    if (side === 'bottom' && top + ch > vh - 14) { top = r.top - ch - m; actual = 'top'; }
    if (side === 'right' && left + cw > vw - 14) { left = r.left - cw - m; actual = 'left'; }
    left = Math.max(14, Math.min(left, vw - cw - 14));
    top = Math.max(14, Math.min(top, vh - ch - 14));
    return { top, left, actual };
  }
  const ringBox = (t, pad = 6) => ({ top: t.top - pad, left: t.left - pad, width: t.width + pad * 2, height: t.height + pad * 2 });
  const tailOffset = (ring, pos, cw, ch) => ring ? (pos.actual === 'bottom' || pos.actual === 'top'
    ? Math.max(22, Math.min(cw - 22, ring.left + ring.width / 2 - pos.left))
    : Math.max(22, Math.min(ch - 22, ring.top + ring.height / 2 - pos.top))) : cw / 2;

  // ── measurement engine ──
  // A single persistent listener reads the CURRENT step via the closure so we
  // never fight stale per-step listeners. Keeps the prior box if the target is
  // briefly missing (no flicker).
  function measure() {
    const s = STEPS[i];
    if (!s || s.kind === 'welcome' || done) { target = null; return; }
    const el = document.querySelector(s.sel);
    if (!el) return; // keep prior box while the view / modal finishes painting
    const r = el.getBoundingClientRect();
    target = { top: r.top, left: r.left, width: r.width, height: r.height };
  }

  let timers = [];
  let raf = 0;
  function scheduleMeasure() {
    cancelAnimationFrame(raf);
    timers.forEach(clearTimeout);
    timers = [];
    if (done) return;
    raf = requestAnimationFrame(() => requestAnimationFrame(measure));
    timers = [40, 130, 260, 420, 640].map((d) => setTimeout(measure, d));
  }

  // On step change: drive the view, then re-measure with timed passes.
  $effect(() => {
    const s = STEPS[i];
    if (done) return;
    // Navigate to the step's view only when the path actually differs — and
    // PRESERVE the current query string (?preview=1&tour=1) so the tour (and
    // preview data) survive the navigation.
    if (s && s.view) {
      const dest = VIEW_PATH[s.view];
      if (location.pathname !== dest) goto(dest + location.search, { keepFocus: true, noScroll: true });
    }
    if (!s || s.kind === 'welcome') { target = null; return; }
    scheduleMeasure();
  });

  onMount(() => {
    seedDemoData();
    window.addEventListener('resize', measure);
    window.addEventListener('scroll', measure, true);
  });
  onDestroy(() => {
    window.removeEventListener('resize', measure);
    window.removeEventListener('scroll', measure, true);
    cancelAnimationFrame(raf);
    timers.forEach(clearTimeout);
  });

  // confirmation auto-dismisses after 3.4s
  $effect(() => {
    if (!done) return;
    const t = setTimeout(() => finish(), 3400);
    return () => clearTimeout(t);
  });

  function advance() { if (i >= STEPS.length - 1) addApp(); else i += 1; }
  function back() { i = Math.max(0, i - 1); }

  async function addApp() {
    if (saving) return;
    saving = true;
    try {
      await api('/api/applications', {
        method: 'POST',
        body: JSON.stringify({
          company: form.company,
          role: form.role,
          status: (form.status || 'applied').toLowerCase(),
          source: form.source
        })
      });
    } catch (e) {
      console.error(e);
    } finally {
      saving = false;
      done = true;
    }
  }

  // bubble geometry
  const CW = 314, CH = 188;
  const pos = $derived(target ? place(target, step.place, CW, CH) : { top: 120, left: (typeof window !== 'undefined' ? window.innerWidth : 1000) - CW - 40, actual: step.place });
  const tail = $derived(tailOffset(target, pos, CW, CH));
  const ring = $derived(target ? ringBox(target, 6) : null);
</script>

{#snippet Icon(name, s)}
  {#if name === 'spark'}
    <svg width={s ?? 14} height={s ?? 14} viewBox="0 0 16 16" fill="currentColor"><path d="M8 1l1.5 4.2L14 7l-4.5 1.8L8 13l-1.5-4.2L2 7l4.5-1.8z"/></svg>
  {:else if name === 'arrow'}
    <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.7"><path d="M3 8h9M8 4l4 4-4 4" stroke-linecap="round" stroke-linejoin="round"/></svg>
  {:else if name === 'x'}
    <svg width={s ?? 13} height={s ?? 13} viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="1.6"><path d="M3 3l8 8M11 3l-8 8" stroke-linecap="round"/></svg>
  {:else if name === 'check'}
    <svg width={s ?? 14} height={s ?? 14} viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 8.5l3.2 3.2L13 5" stroke-linecap="round" stroke-linejoin="round"/></svg>
  {:else if name === 'today'}
    <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round"><rect x="2" y="2.5" width="4" height="4" rx="1"/><rect x="9" y="2.5" width="5" height="4" rx="1"/><rect x="2" y="9" width="4" height="4.5" rx="1"/><rect x="9" y="9" width="5" height="4.5" rx="1"/></svg>
  {:else if name === 'board'}
    <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round"><rect x="2" y="3" width="3" height="10" rx="0.5"/><rect x="7" y="3" width="3" height="7" rx="0.5"/><rect x="12" y="3" width="2" height="5" rx="0.5"/></svg>
  {:else if name === 'insights'}
    <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round"><path d="M2 3h12l-5 6v5l-2-1V9z"/></svg>
  {:else if name === 'add'}
    <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.7"><path d="M8 3v10M3 8h10" stroke-linecap="round"/></svg>
  {/if}
{/snippet}

{#if done}
  <div class="cmA-toastwrap">
    <div class="cmA-toast">
      <span class="cmA-toast-ic">{@render Icon('check', 15)}</span>
      <div class="cmA-toast-tx">
        <b>You're all set.</b>
        <span>Your applications live here — add another anytime from <em>New application</em></span>
      </div>
      <button class="cmA-toast-x" onclick={() => finish()}>{@render Icon('x', 12)}</button>
    </div>
  </div>
{:else if step.kind === 'welcome'}
  <div class="cmA-introwrap">
    <div class="cmA-intro">
      <div class="cmA-intro-row">
        <div class="cmA-av">{@render Icon('spark', 17)}</div>
        <div class="cmA-intro-tx">
          <b>Welcome to Pursuit 👋</b>
          <p>Give me a minute and I'll show you the handful of things that make your search easier — then help you add your first application.</p>
        </div>
      </div>
      <div class="cmA-intro-foot">
        <button class="t-ghost sm" onclick={() => finish()}>No thanks</button>
        <button class="t-cta sm" onclick={() => (i = 1)}>Show me {@render Icon('arrow')}</button>
      </div>
    </div>
  </div>
{:else}
  {#if step.modal}
    <div class="cm-modalback shift" onmousedown={(e) => { if (e.target === e.currentTarget) back(); }} role="presentation">
      <div class="nux-modal">
        <header class="nm-hd">
          <div><h3>New application</h3><p>Paste a link or screenshot and Pursuit fills it in — or type it by hand.</p></div>
          <button class="nm-x" onclick={back}>{@render Icon('x')}</button>
        </header>
        <div data-tour="paste" class="nm-quick">
          <span class="qi">{@render Icon('spark', 15)}</span>
          <span class="qt">Drop a screenshot or paste a job URL<small>Pursuit parses the company, role, and location for you.</small></span>
        </div>
        <div class="nm-or"><span>or enter by hand</span></div>
        <div class="nm-form">
          <div class="nm-field"><label>Company <span class="req">*</span></label><input class="nm-input" bind:value={form.company} /></div>
          <div class="nm-field"><label>Role <span class="req">*</span></label><input class="nm-input" bind:value={form.role} /></div>
          <div class="nm-field"><label>Status</label>
            <select class="nm-select" bind:value={form.status}>
              <option>Wishlist</option><option>Applied</option><option>Screen</option><option>Interview</option><option>Offer</option>
            </select>
          </div>
          <div class="nm-field"><label>Source</label><input class="nm-input" bind:value={form.source} /></div>
        </div>
        <footer class="nm-ft">
          <button class="nm-btn" onclick={back}>Cancel</button>
          <button class="nm-btn primary" onclick={addApp} disabled={saving}>Add application</button>
        </footer>
      </div>
    </div>
  {/if}

  {#if ring}
    <div class="cmA-ring" style="top:{ring.top}px;left:{ring.left}px;width:{ring.width}px;height:{ring.height}px"></div>
    <div class="cmA-beacon" style="top:{target.top - 9}px;left:{target.left + target.width - 3}px"><i></i><i></i></div>
  {/if}

  <div class="t-card cmA-bubble tail-{pos.actual}" style="top:{pos.top}px;left:{pos.left}px;width:{CW}px;--tail:{tail}px">
    <div class="cmA-head">
      <span class="cmA-num">{i} <span>/ {N}</span></span>
      <button class="cmA-skip" onclick={() => finish()}>Skip tour</button>
    </div>
    <h4><span class="cmA-h-ic">{@render Icon(step.icon, 15)}</span>{step.title}</h4>
    <p>{step.body}</p>
    <div class="t-foot">
      <div class="t-dots">
        {#each Array.from({ length: N }) as _, k}
          <i class={k === i - 1 ? 'on' : k < i - 1 ? 'done' : ''}></i>
        {/each}
      </div>
      <div class="t-btns">
        <button class="t-btn ghost" onclick={back}>Back</button>
        <button class="t-btn primary" onclick={advance}>{step.last ? 'Add it' : 'Next'}{#if !step.last}{@render Icon('arrow')}{/if}</button>
      </div>
    </div>
  </div>
{/if}

<style>
  @keyframes t-rise { from { transform: translateY(10px) scale(.985); } to { transform: none; } }
  @keyframes t-up { from { transform: translateX(-50%) translateY(14px); } to { transform: translateX(-50%) translateY(0); } }
  @keyframes t-pulse { 0% { box-shadow: 0 0 0 0 var(--accent); opacity: .9; } 70% { box-shadow: 0 0 0 9px transparent; opacity: 0; } 100% { opacity: 0; } }
  @keyframes t-breathe { 0%,100% { box-shadow: 0 0 0 3px var(--accent-tint-2), 0 0 0 1.5px var(--accent); } 50% { box-shadow: 0 0 0 6px var(--accent-tint), 0 0 0 1.5px var(--accent); } }

  /* ── shared primitives ── */
  .t-card { position: fixed; z-index: 120; background: var(--card); border: 1px solid var(--rule);
    border-radius: 14px; box-shadow: 0 18px 50px -14px rgba(10,10,13,.30), 0 4px 12px rgba(10,10,13,.06);
    padding: 17px 18px 15px; }
  .t-card h4 { margin: 0 0 7px; font-size: 16px; font-weight: 600; letter-spacing: -0.02em; color: var(--ink); line-height: 1.25; }
  .t-card p { margin: 0; font-size: 13.3px; line-height: 1.58; color: var(--ink-2); text-wrap: pretty; }

  .t-foot { display: flex; align-items: center; justify-content: space-between; margin-top: 15px; gap: 12px; }
  .t-dots { display: flex; gap: 5px; }
  .t-dots i { width: 6px; height: 6px; border-radius: 50%; background: var(--rule-strong); transition: all .25s ease; }
  .t-dots i.done { background: var(--accent-tint-2); }
  .t-dots i.on { background: var(--accent); width: 17px; border-radius: 3px; }

  .t-btns { display: flex; gap: 7px; }
  .t-btn { display: inline-flex; align-items: center; gap: 5px; font-size: 12.8px; font-weight: 500; white-space: nowrap;
    padding: 7px 13px; border-radius: 8px; cursor: pointer; transition: background .12s ease, border-color .12s ease; border: 1px solid transparent; font-family: inherit; }
  .t-btn.primary { background: var(--accent); color: #fff; }
  .t-btn.primary:hover { background: var(--accent-strong); }
  .t-btn.ghost { color: var(--ink-2); border: 1px solid var(--rule); background: var(--card); }
  .t-btn.ghost:hover { background: var(--surface-2); }

  .t-cta { display: inline-flex; align-items: center; gap: 8px; background: var(--ink); color: #fff; white-space: nowrap;
    font-size: 13.5px; font-weight: 500; padding: 11px 18px; border-radius: 9px; cursor: pointer; border: none; font-family: inherit; }
  .t-cta:hover { background: #000; }
  .t-cta.sm { padding: 8px 14px; font-size: 12.8px; border-radius: 8px; }
  .t-ghost { font-size: 13px; font-weight: 500; color: var(--mute); padding: 11px 8px; cursor: pointer; white-space: nowrap; background: none; border: none; font-family: inherit; }
  .t-ghost:hover { color: var(--ink-2); }
  .t-ghost.sm { padding: 8px 10px; font-size: 12.5px; }

  /* ════════════════ BUBBLE ════════════════ */
  .cmA-ring { position: fixed; z-index: 108; border-radius: 11px; pointer-events: none;
    animation: t-breathe 2.4s ease-in-out infinite; }
  .cmA-beacon { position: fixed; z-index: 109; width: 14px; height: 14px; pointer-events: none; }
  .cmA-beacon i { position: absolute; inset: 0; border-radius: 50%; background: var(--accent); }
  .cmA-beacon i:last-child { animation: t-pulse 2s ease-out infinite; }

  .cmA-bubble { animation: t-rise .24s cubic-bezier(.2,.7,.3,1) both; }
  .cmA-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 8px; }
  .cmA-num { font-size: 11.5px; font-weight: 600; color: var(--accent-text); font-variant-numeric: tabular-nums; }
  .cmA-num span { color: var(--mute-2); font-weight: 500; }
  .cmA-skip { font-size: 12px; font-weight: 500; color: var(--mute-2); padding: 3px 4px; border-radius: 5px; cursor: pointer; white-space: nowrap; background: none; border: none; font-family: inherit; }
  .cmA-skip:hover { color: var(--ink-2); background: var(--surface-2); }
  .cmA-bubble h4 { display: flex; align-items: center; gap: 8px; }
  .cmA-h-ic { display: inline-flex; width: 26px; height: 26px; border-radius: 7px; background: var(--accent-tint);
    color: var(--accent-text); align-items: center; justify-content: center; flex-shrink: 0; }

  .cmA-bubble::before { content: ""; position: absolute; width: 14px; height: 14px; background: var(--card);
    border: 1px solid var(--rule); transform: rotate(45deg); }
  /* tail-* classes are built dynamically; keep them global so Svelte's
     static-selector scoping doesn't tree-shake them away. */
  :global(.cmA-bubble.tail-bottom)::before { top: -8px; left: var(--tail); margin-left: -7px; border-right: none; border-bottom: none; }
  :global(.cmA-bubble.tail-top)::before { bottom: -8px; left: var(--tail); margin-left: -7px; border-left: none; border-top: none; }
  :global(.cmA-bubble.tail-right)::before { left: -8px; top: var(--tail); margin-top: -7px; border-right: none; border-top: none; }
  :global(.cmA-bubble.tail-left)::before { right: -8px; top: var(--tail); margin-top: -7px; border-left: none; border-bottom: none; }

  .cmA-introwrap { position: fixed; inset: 0; z-index: 120; display: grid; place-items: center; padding: 24px;
    background: rgba(10,10,13,.30); }
  .cmA-intro { width: 392px; max-width: 100%; background: var(--card); border: 1px solid var(--rule); border-radius: 18px;
    padding: 22px; box-shadow: 0 30px 70px -20px rgba(10,10,13,.40); animation: t-rise .3s cubic-bezier(.2,.7,.3,1) both; }
  .cmA-intro-row { display: flex; gap: 13px; }
  .cmA-av { width: 38px; height: 38px; border-radius: 50%; background: var(--accent); color: #fff;
    display: grid; place-items: center; flex-shrink: 0; }
  .cmA-intro-tx b { font-size: 14.5px; font-weight: 600; color: var(--ink); }
  .cmA-intro-tx p { margin: 5px 0 0; font-size: 13px; line-height: 1.55; color: var(--ink-2); }
  .cmA-intro-foot { display: flex; align-items: center; justify-content: flex-end; gap: 4px; margin-top: 15px; }

  /* completion confirmation */
  .cmA-toastwrap { position: fixed; left: 50%; bottom: 84px; transform: translateX(-50%); z-index: 120;
    animation: t-up .34s cubic-bezier(.2,.7,.3,1) both; }
  .cmA-toast { display: flex; align-items: flex-start; gap: 12px; width: min(440px, calc(100vw - 32px));
    background: var(--card); border: 1px solid var(--rule); border-radius: 14px; padding: 14px 14px 14px 15px;
    box-shadow: 0 20px 50px -14px rgba(10,10,13,.32); }
  .cmA-toast-ic { width: 30px; height: 30px; border-radius: 8px; background: var(--positive-tint); color: var(--positive-text);
    display: grid; place-items: center; flex-shrink: 0; margin-top: 1px; }
  .cmA-toast-tx { flex: 1; min-width: 0; }
  .cmA-toast-tx b { display: block; font-size: 14px; font-weight: 600; color: var(--ink); letter-spacing: -0.01em; }
  .cmA-toast-tx span { display: block; font-size: 12.8px; line-height: 1.5; color: var(--mute); margin-top: 2px; }
  .cmA-toast-tx em { font-style: normal; font-weight: 500; color: var(--ink-2); }
  .cmA-toast-x { width: 22px; height: 22px; display: grid; place-items: center; border-radius: 6px; color: var(--mute-2); flex-shrink: 0; background: none; border: none; cursor: pointer; }
  .cmA-toast-x:hover { background: var(--surface-2); color: var(--ink-2); }

  /* ════════════════ Add-application modal ════════════════ */
  .cm-modalback { position: fixed; inset: 0; z-index: 100; background: rgba(10,10,13,.42);
    display: grid; place-items: center; padding: 24px; }
  .cm-modalback.shift { padding-right: min(372px, 38vw); }
  .nux-modal { width: 540px; max-width: 100%; background: var(--card); border: 1px solid var(--rule); border-radius: 16px;
    box-shadow: 0 30px 70px -20px rgba(10,10,13,0.42); animation: t-rise .26s cubic-bezier(.2,.7,.3,1) both; }
  .nux-modal .nm-hd { padding: 22px 24px 16px; display: grid; grid-template-columns: 1fr 26px; gap: 12px; align-items: start; }
  .nux-modal .nm-hd h3 { margin: 0; font-size: 20px; font-weight: 600; letter-spacing: -0.02em; }
  .nux-modal .nm-hd p { margin: 6px 0 0; font-size: 13.5px; color: var(--mute); line-height: 1.5; }
  .nux-modal .nm-x { width: 28px; height: 28px; border-radius: 7px; color: var(--mute); display: grid; place-items: center; background: none; border: none; cursor: pointer; }
  .nux-modal .nm-x:hover { background: var(--surface-2); color: var(--ink-2); }
  .nm-quick { margin: 0 24px; border: 1px dashed var(--rule-strong); border-radius: 11px; background: var(--surface);
    padding: 13px 15px; display: grid; grid-template-columns: 30px 1fr; gap: 12px; align-items: center; }
  .nm-quick .qi { width: 30px; height: 30px; border-radius: 8px; background: var(--accent-tint); color: var(--accent-text); display: grid; place-items: center; }
  .nm-quick .qt { font-size: 13px; color: var(--ink-2); } .nm-quick .qt small { display: block; color: var(--mute); font-size: 12px; margin-top: 1px; }
  .nm-or { display: grid; grid-template-columns: 1fr auto 1fr; align-items: center; gap: 14px; margin: 16px 24px 14px; font-size: 12px; color: var(--mute-2); }
  .nm-or::before, .nm-or::after { content: ""; height: 1px; background: var(--rule); }
  .nm-form { padding: 0 24px; display: grid; grid-template-columns: 1fr 1fr; gap: 13px 15px; }
  .nm-field { display: flex; flex-direction: column; gap: 6px; min-width: 0; }
  .nm-field label { font-size: 12.5px; font-weight: 500; color: var(--ink-2); }
  .nm-field label .req { color: var(--accent); }
  .nm-input, .nm-select { height: 38px; border: 1px solid var(--rule); border-radius: 9px; background: var(--card);
    color: var(--ink); font: 400 14px/1.2 var(--sans); padding: 0 12px; width: 100%; outline: none; appearance: none; -webkit-appearance: none; }
  .nm-input::placeholder { color: var(--mute-2); }
  .nm-input:focus, .nm-select:focus { border-color: var(--accent); box-shadow: 0 0 0 3px var(--accent-tint); }
  .nm-select { cursor: pointer; }
  .nm-ft { margin-top: 20px; padding: 16px 24px 18px; border-top: 1px solid var(--rule); background: var(--surface);
    border-radius: 0 0 16px 16px; display: flex; align-items: center; justify-content: flex-end; gap: 9px; }
  .nm-btn { height: 38px; padding: 0 16px; border-radius: 9px; font-size: 13.5px; font-weight: 500; border: 1px solid var(--rule);
    background: var(--card); color: var(--ink-2); cursor: pointer; font-family: inherit; transition: background .12s, border-color .12s; }
  .nm-btn:hover { background: var(--surface-2); border-color: var(--rule-strong); }
  .nm-btn.primary { background: var(--accent); color: #fff; border-color: var(--accent-strong); }
  .nm-btn.primary:hover { background: var(--accent-strong); }
  .nm-btn:disabled { opacity: .6; cursor: default; }

  /* ── reduced motion ── */
  @media (prefers-reduced-motion: reduce) {
    .cmA-bubble, .cmA-introwrap, .cmA-toastwrap, .nux-modal { animation: none !important; }
    .cmA-ring, .cmA-beacon i:last-child { animation: none !important; }
    .cmA-ring, .t-dots i { transition: none !important; }
  }
</style>
