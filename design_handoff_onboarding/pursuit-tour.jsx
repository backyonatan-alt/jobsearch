/* ═══════════════════════════════════════════════════════════
   PURSUIT — First-run onboarding (final)
   A no-dim Coachmark tour: a friendly speech bubble points at each
   thing Pursuit does, then opens the Add-application modal so the
   new user creates their first application. Ends on a confirmation.
   Exposes window.TourDesigns.{ CoachBubble, Harness }
   ═══════════════════════════════════════════════════════════ */
(function () {
  const { useState, useEffect, useRef, useCallback } = React;

  /* ── icons ── */
  const Spk = ({ s = 14 }) => (
    <svg width={s} height={s} viewBox="0 0 16 16" fill="currentColor"><path d="M8 1l1.5 4.2L14 7l-4.5 1.8L8 13l-1.5-4.2L2 7l4.5-1.8z" /></svg>
  );
  const Arr = () => (
    <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.7"><path d="M3 8h9M8 4l4 4-4 4" strokeLinecap="round" strokeLinejoin="round" /></svg>
  );
  const X = ({ s = 13 }) => (
    <svg width={s} height={s} viewBox="0 0 14 14" fill="none" stroke="currentColor" strokeWidth="1.6"><path d="M3 3l8 8M11 3l-8 8" strokeLinecap="round" /></svg>
  );
  const Check = ({ s = 14 }) => (
    <svg width={s} height={s} viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="2"><path d="M3 8.5l3.2 3.2L13 5" strokeLinecap="round" strokeLinejoin="round" /></svg>
  );
  const ICONS = {
    today: <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinejoin="round"><rect x="2" y="2.5" width="4" height="4" rx="1" /><rect x="9" y="2.5" width="5" height="4" rx="1" /><rect x="2" y="9" width="4" height="4.5" rx="1" /><rect x="9" y="9" width="5" height="4.5" rx="1" /></svg>,
    board: <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinejoin="round"><rect x="2" y="3" width="3" height="10" rx="0.5" /><rect x="7" y="3" width="3" height="7" rx="0.5" /><rect x="12" y="3" width="2" height="5" rx="0.5" /></svg>,
    insights: <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinejoin="round"><path d="M2 3h12l-5 6v5l-2-1V9z" /></svg>,
    prep: <Spk s={15} />,
    add: <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.7"><path d="M8 3v10M3 8h10" strokeLinecap="round" /></svg>,
  };

  /* ── tour script: welcome → 5 highlights → add-application modal ── */
  const STEPS = [
    { kind: "welcome" },
    { sel: '[data-tour="stats"]', view: "today", icon: ICONS.today, kick: "Today", title: "Your morning briefing", body: "Every day opens here. Pursuit reads your pipeline and tells you what's moving, what's waiting, and what's gone quiet — before you ask.", place: "bottom" },
    { sel: '[data-tour="prep"]', view: "today", icon: ICONS.prep, kick: "AI prep", title: "It preps you for interviews", body: "Ahead of every interview, Pursuit researches who you'll meet and what they care about — so you walk in ready, not guessing.", place: "bottom" },
    { sel: '[data-tour="board"]', view: "board", icon: ICONS.board, kick: "Board", title: "Your whole search, one board", body: "Each role is a card moving across stages. Drag a card to update its status — and it turns red when it's stalled for a week.", place: "bottom" },
    { sel: '[data-tour="funnel"]', view: "insights", icon: ICONS.insights, kick: "Insights", title: "See what's actually working", body: "Your reply rate, how far you get, and which sources convert — so you spend effort where it pays off.", place: "right" },
    { sel: '[data-tour="new-app"]', view: "today", icon: ICONS.add, kick: "Your turn", title: "Add your first application", body: "This is where every role starts. Let's add your first one — it takes seconds.", place: "bottom" },
    { modal: true, sel: '[data-tour="paste"]', view: "today", icon: ICONS.add, kick: "Add an application", title: "Paste a link — Pursuit does the rest", body: "Drop a job URL or a screenshot here and Pursuit reads the company, role, and location for you. Prefer to type it? The fields are right below.", place: "right", last: true },
  ];
  const N = STEPS.length - 1; // number of real (non-welcome) stops

  /* place a card of size (cw×ch) relative to a target rect */
  function place(r, side, cw, ch) {
    const m = 16, vw = innerWidth, vh = innerHeight;
    let top, left;
    if (side === "right") { left = r.left + r.width + m; top = r.top + r.height / 2 - ch / 2; }
    else if (side === "left") { left = r.left - cw - m; top = r.top + r.height / 2 - ch / 2; }
    else if (side === "top") { left = r.left + r.width / 2 - cw / 2; top = r.top - ch - m; }
    else { left = r.left + r.width / 2 - cw / 2; top = r.top + r.height + m; }
    let actual = side;
    if (side === "bottom" && top + ch > vh - 14) { top = r.top - ch - m; actual = "top"; }
    if (side === "right" && left + cw > vw - 14) { left = r.left - cw - m; actual = "left"; }
    left = Math.max(14, Math.min(left, vw - cw - 14));
    top = Math.max(14, Math.min(top, vh - ch - 14));
    return { top, left, actual };
  }
  const ringBox = (t, pad = 6) => ({ top: t.top - pad, left: t.left - pad, width: t.width + pad * 2, height: t.height + pad * 2 });
  const tailOffset = (ring, pos, cw, ch) => ring ? (pos.actual === "bottom" || pos.actual === "top"
    ? Math.max(22, Math.min(cw - 22, ring.left + ring.width / 2 - pos.left))
    : Math.max(22, Math.min(ch - 22, ring.top + ring.height / 2 - pos.top))) : cw / 2;

  /* engine: step index, keeps app on the right view, measures the target.
     One persistent listener reads the current step via a ref (no stale duplicates),
     plus rAF + timed re-measures on each change — robust across view swaps + modal mounts. */
  function useCoach(view, setView) {
    const [i, setI] = useState(0);
    const [target, setTarget] = useState(null);
    const step = STEPS[i];
    const stepRef = useRef(step);
    stepRef.current = step;

    useEffect(() => { if (step.view && step.view !== view) setView(step.view); }, [i]);

    const measure = useCallback(() => {
      const s = stepRef.current;
      if (s.kind === "welcome") { setTarget(null); return; }
      const el = document.querySelector(s.sel);
      if (!el) return; // keep prior box while the view / modal finishes painting
      const r = el.getBoundingClientRect();
      setTarget({ top: r.top, left: r.left, width: r.width, height: r.height });
    }, []);

    useEffect(() => {
      addEventListener("resize", measure); addEventListener("scroll", measure, true);
      return () => { removeEventListener("resize", measure); removeEventListener("scroll", measure, true); };
    }, [measure]);

    useEffect(() => {
      if (step.kind === "welcome") { setTarget(null); return; }
      const raf = requestAnimationFrame(() => requestAnimationFrame(measure));
      const ts = [40, 130, 260, 420, 640].map((d) => setTimeout(measure, d));
      return () => { cancelAnimationFrame(raf); ts.forEach(clearTimeout); };
    }, [i, view, measure]);

    return { i, setI, step, target };
  }

  /* ── Add-application modal (closing stage). shift nudges it left to clear the bubble ── */
  function AddAppModal({ onClose, onAdd }) {
    return (
      <div className="cm-modalback shift" onMouseDown={(e) => { if (e.target === e.currentTarget) onClose(); }}>
        <div className="nux-modal">
          <header className="nm-hd">
            <div><h3>New application</h3><p>Paste a link or screenshot and Pursuit fills it in — or type it by hand.</p></div>
            <button className="nm-x" onClick={onClose}><X /></button>
          </header>
          <div data-tour="paste" className="nm-quick">
            <span className="qi"><Spk s={15} /></span>
            <span className="qt">Drop a screenshot or paste a job URL<small>Pursuit parses the company, role, and location for you.</small></span>
          </div>
          <div className="nm-or"><span>or enter by hand</span></div>
          <div className="nm-form">
            <div className="nm-field"><label>Company <span className="req">*</span></label><input className="nm-input" defaultValue="Anthropic" /></div>
            <div className="nm-field"><label>Role <span className="req">*</span></label><input className="nm-input" defaultValue="Member of Technical Staff" /></div>
            <div className="nm-field"><label>Status</label><select className="nm-select" defaultValue="Applied"><option>Wishlist</option><option>Applied</option><option>Screen</option><option>Interview</option><option>Offer</option></select></div>
            <div className="nm-field"><label>Source</label><input className="nm-input" defaultValue="Referral" /></div>
          </div>
          <footer className="nm-ft">
            <button className="nm-btn" onClick={onClose}>Cancel</button>
            <button className="nm-btn primary" onClick={onAdd}>Add application</button>
          </footer>
        </div>
      </div>
    );
  }

  const Dots = ({ i }) => (
    <div className="t-dots">{Array.from({ length: N }).map((_, k) => <i key={k} className={k === i - 1 ? "on" : k < i - 1 ? "done" : ""} />)}</div>
  );

  /* ════════════════════════════════════════════════════════
     COACHMARK — friendly speech bubble + tail + breathing beacon
     ════════════════════════════════════════════════════════ */
  function CoachBubble({ view, setView, onClose }) {
    const { i, setI, step, target } = useCoach(view, setView);
    const [done, setDone] = useState(false);

    const advance = () => (i >= STEPS.length - 1 ? setDone(true) : setI(i + 1));
    const back = () => setI(Math.max(0, i - 1));

    // confirmation auto-dismisses
    useEffect(() => {
      if (!done) return;
      const t = setTimeout(onClose, 3400);
      return () => clearTimeout(t);
    }, [done]);

    if (done) {
      return (
        <div className="cmA-toastwrap">
          <div className="cmA-toast">
            <span className="cmA-toast-ic"><Check s={15} /></span>
            <div className="cmA-toast-tx">
              <b>You're all set.</b>
              <span>Your applications live here — add another anytime from <em>New application</em></span>
            </div>
            <button className="cmA-toast-x" onClick={onClose}><X s={12} /></button>
          </div>
        </div>
      );
    }

    if (step.kind === "welcome") {
      return (
        <div className="cmA-introwrap">
          <div className="cmA-intro">
            <div className="cmA-intro-row">
              <div className="cmA-av"><Spk s={17} /></div>
              <div className="cmA-intro-tx">
                <b>Welcome to Pursuit 👋</b>
                <p>Give me a minute and I'll show you the handful of things that make your search easier — then help you add your first application.</p>
              </div>
            </div>
            <div className="cmA-intro-foot">
              <button className="t-ghost sm" onClick={onClose}>No thanks</button>
              <button className="t-cta sm" onClick={() => setI(1)}>Show me <Arr /></button>
            </div>
          </div>
        </div>
      );
    }

    const cw = 314, ch = 188, ring = target;
    const pos = target ? place(target, step.place, cw, ch) : { top: 120, left: innerWidth - cw - 40, actual: step.place };
    const tail = tailOffset(ring, pos, cw, ch);

    return (
      <>
        {step.modal && <AddAppModal onClose={back} onAdd={advance} />}
        {ring && <div className="cmA-ring" style={ringBox(ring, 6)} />}
        {ring && <div className="cmA-beacon" style={{ top: ring.top - 9, left: ring.left + ring.width - 3 }}><i /><i /></div>}
        <div className={`t-card cmA-bubble tail-${pos.actual}`} style={{ top: pos.top, left: pos.left, width: cw, "--tail": tail + "px" }}>
          <div className="cmA-head">
            <span className="cmA-num">{i} <span>/ {N}</span></span>
            <button className="cmA-skip" onClick={onClose}>Skip tour</button>
          </div>
          <h4><span className="cmA-h-ic">{step.icon}</span>{step.title}</h4>
          <p>{step.body}</p>
          <div className="t-foot">
            <Dots i={i} />
            <div className="t-btns">
              <button className="t-btn ghost" onClick={back}>Back</button>
              <button className="t-btn primary" onClick={advance}>{step.last ? "Add it" : "Next"}{!step.last && <Arr />}</button>
            </div>
          </div>
        </div>
      </>
    );
  }

  /* ── demo control (not part of the product): replay as a new user / dismiss ── */
  function Harness({ active, replay, dismiss }) {
    return (
      <div className="harness">
        <span className="hz-lbl">First-run onboarding</span>
        <button className="hz-replay" onClick={replay} title="Replay as a brand-new user">
          <svg width="13" height="13" viewBox="0 0 14 14" fill="none" stroke="currentColor" strokeWidth="1.6"><path d="M2.5 7a4.5 4.5 0 1 1 1.4 3.2M2.5 4.5v2.5h2.5" strokeLinecap="round" strokeLinejoin="round" /></svg>
          Replay
        </button>
        <span className="hz-div" />
        <button className="hz-replay" onClick={dismiss} disabled={!active} style={{ opacity: active ? 1 : 0.4 }}>Dismiss</button>
      </div>
    );
  }

  window.TourDesigns = { CoachBubble, Harness };
})();
