/* ───────────────────────────────────────────────────────────
   PURSUIT · Today — primary screens
   TodayScreen (Option B, activated + widened) · Board · Insights
   ─────────────────────────────────────────────────────────── */
const { useState } = React;

/* ═══ TODAY ══════════════════════════════════════════════════ */
function TodayScreen({ nav }) {
  const waitingIds = ["figma", "notion", "ramp", "mistral"];
  const waiting = waitingIds.map(window.getApp);

  const [tasks, setTasks] = useState([
    { id: "t1", b: "Prep 3 questions for Dario", s: "Final round · Anthropic", due: "Today", hot: true, done: false },
    { id: "t2", b: "Decide on the Vercel offer", s: "$210k · respond by Friday", due: "Fri", hot: true, done: false },
    { id: "t3", b: "Send Plain your take-home link", s: "Walkthrough tomorrow 16:30", due: "Thu", hot: false, done: false },
    { id: "t4", b: "Confirm Anthropic call time", s: "Accepted the calendar invite", due: "Done", hot: false, done: true },
  ]);
  const toggle = (id) => setTasks((t) => t.map((x) => (x.id === id ? { ...x, done: !x.done } : x)));
  const openTasks = tasks.filter((t) => !t.done).length;

  return (
    <div className="ob ob-alt ob-swap">
      {/* ── LEFT: the brief ── */}
      <BriefRight nav={nav} />

      {/* ── RIGHT: pipeline pulse + tasks ── */}
      <div className="pulse-stage">
        <div className="pulse-tag"><span className="d"></span>Where things stand</div>

        <div className="pulse-stats">
          <div className="st" onClick={() => nav("board")}><span className="num">12</span><span className="lbl">Active loops</span></div>
          <div className="st" onClick={() => nav("board")}><span className="num">4</span><span className="lbl">Awaiting reply</span></div>
          <div className="st warn" onClick={() => nav("board")}><span className="num">2</span><span className="lbl">Gone quiet</span></div>
        </div>

        <div className="pulse-sec">
          <span className="t">Waiting to hear back</span>
          <span className="c">longest first</span>
        </div>
        <div className="pulse-list">
          {waiting.map((w) => (
            <div className={"pulse-row" + (w.quiet ? " quiet" : "")} key={w.id} onClick={() => nav("detail", { appId: w.id })}>
              <Logo co={w.short} cls={w.cls} size={30} />
              <span className="wx"><b>{w.co}</b><small>{w.stage}</small></span>
              <span className="days">{w.days}d</span>
              <span className="ok"><span className={"okdot" + (w.quiet ? " warn" : "")}></span></span>
            </div>
          ))}
        </div>

        {/* widened: a small "needs a move" task list */}
        <div className="tasks">
          <div className="pulse-sec">
            <span className="t">Your move</span>
            <span className="c">{openTasks} to do</span>
          </div>
          {tasks.map((t) => (
            <div className={"task" + (t.done ? " done" : "")} key={t.id} onClick={() => toggle(t.id)}>
              <span className="box"></span>
              <span className="tx"><b>{t.b}</b><small>{t.s}</small></span>
              <span className={"due" + (t.hot && !t.done ? " hot" : "")}>{t.due}</span>
            </div>
          ))}
        </div>

        <div className="pulse-foot">
          <span className="fic"><Spark s={15} /></span>
          <span className="ftx"><b>Figma and Notion have gone quiet</b><small>No reply in over a week — it might be a good time to reach out to them directly.</small></span>
          <button className="pulse-link" onClick={() => nav("board")}>See both {ICONS.arrow}</button>
        </div>
      </div>

    </div>
  );
}

/* The editorial brief (right column) — clickable */
function BriefRight({ nav }) {
  return (
    <div className="brief">
      <div className="brief-in">
        <div className="brief-date">{window.NOW.dow} · {window.NOW.full}</div>

        <h1>Good morning,<br /><b>{window.USER.first}.</b></h1>
        <p className="lede">Today it's your <span className="hot">final round at Anthropic</span> — the ML Engineer role, one-on-one with Dario Amodei. Two more screens follow later this week.</p>

        <div className="kick"><Spark s={13} />&nbsp;Before the room</div>
        <div className="insight">
          <span className="ic"><Spark s={15} /></span>
          <span className="tx">Going on his essays and recent interviews, Dario reasons from first principles and seems more interested in <b>how you think</b> than what you've shipped — so narrate your reasoning, not just the result.</span>
        </div>

        <div className="brief-sub">Worth reviewing</div>
        <ul className="brief-review">
          <li>Your reasoning on scaling intuitions — lead with the <b>why</b>, not just the result.</li>
          <li>One recent Anthropic paper or essay, with a genuine question ready to ask about it.</li>
          <li>Two tradeoffs you weighed on your last project, and what you'd do differently.</li>
        </ul>

        <div className="brief-sub">Two quick tips</div>
        <div className="brief-tip"><span className="sp"><Spark s={13} /></span><span>Referencing a <b>specific piece of his writing</b> signals you actually engage with the work — better than generic enthusiasm.</span></div>
        <div className="brief-tip"><span className="sp"><Spark s={13} /></span><span>He talks about safety constantly in public — <b>anchor your interest to a concrete problem</b>, not the hype.</span></div>

        <div className="meta"><b>60 min</b><span className="dot"></span>Google Meet<span className="dot"></span>Final round</div>
        <button className="cta" onClick={() => nav("dossier", { appId: "anthropic" })}>Open the full playbook {ICONS.arrow}</button>

        <div className="agenda">
          <div className="kick">Later this week</div>
          <div className="ag-row" onClick={() => nav("detail", { appId: "stripe" })}><span className="when"><b>Thu</b> 11:00</span><span><span className="co">Stripe</span> <span className="role">· Product Designer screen</span></span><Pill status="screen" label="Screen" /></div>
          <div className="ag-row" onClick={() => nav("detail", { appId: "plain" })}><span className="when"><b>Thu</b> 16:30</span><span><span className="co">Plain</span> <span className="role">· Senior Designer take-home</span></span><Pill status="screen" label="Screen" /></div>
          <div className="ag-row" onClick={() => nav("detail", { appId: "vercel" })}><span className="when"><b>Fri</b> 17:00</span><span><span className="co">Vercel</span> <span className="role">· offer decision due</span></span><Pill status="offer" label="Offer" /></div>
        </div>

        <div className="foot" onClick={() => nav("board")}><b>12 applications</b> tracked · open the board {ICONS.arrow}</div>
      </div>
    </div>
  );
}

/* ═══ BOARD ══════════════════════════════════════════════════ */
const COLS = [
  { k: "wishlist", lbl: "Wishlist" },
  { k: "applied", lbl: "Applied" },
  { k: "screen", lbl: "Screen" },
  { k: "interview", lbl: "Interview" },
  { k: "offer", lbl: "Offer" },
];

function BoardScreen({ nav }) {
  const [moves, setMoves] = useState({});  // id -> new status after a drag
  const [over, setOver] = useState(null);  // column being dragged over

  const apps = window.APPS
    .filter((a) => a.status !== "closed")
    .map((a) => ({ ...a, status: moves[a.id] || a.status, moved: !!moves[a.id] }));

  const drop = (id, status) => {
    if (id) setMoves((m) => ({ ...m, [id]: status }));
    setOver(null);
  };

  const agoText = (a) => {
    if (a.moved) return "just now";
    if (a.status === "wishlist") return "—";
    if (a.days === 0) return "today";
    if (a.days === 1) return "yesterday";
    return a.days + " days ago";
  };

  return (
    <div className="board2">
      <div className="bhead">
        <div className="bdate">{window.NOW.dow}, {window.NOW.full}</div>
        <h1>Board.</h1>
        <div className="bsub">
          <b>{apps.length}</b> in flight · drag a card across columns to move its status.
          <span className="legend"><span className="rd"></span>red border = no movement in 7+ days</span>
        </div>
      </div>

      <div className="bcols">
        {COLS.map((c) => {
          const items = apps.filter((a) => a.status === c.k);
          return (
            <div
              key={c.k}
              className={"bcol" + (over === c.k ? " over" : "")}
              onDragOver={(e) => { e.preventDefault(); if (over !== c.k) setOver(c.k); }}
              onDragLeave={(e) => { if (e.currentTarget === e.target) setOver(null); }}
              onDrop={(e) => { e.preventDefault(); drop(e.dataTransfer.getData("text/plain"), c.k); }}
            >
              <div className="bcol-h">
                <span className={"bcol-tag " + c.k}><span className="dot"></span>{c.lbl}</span>
                <span className="bcol-ct">{items.length}</span>
                <button className="bcol-add" title="Add application">+</button>
              </div>
              <div className="bcol-list">
                {items.map((a) => {
                  const stale = !a.moved && a.status !== "wishlist" && a.days >= 7;
                  return (
                    <div
                      key={a.id}
                      className={"bcard" + (stale ? " stale" : "")}
                      draggable
                      onDragStart={(e) => { e.dataTransfer.setData("text/plain", a.id); e.dataTransfer.effectAllowed = "move"; }}
                      onDragEnd={() => setOver(null)}
                      onClick={() => nav("detail", { appId: a.id })}
                    >
                      <div className="bcard-top">
                        <Logo co={a.short} cls={a.cls} size={20} />
                        <span className="co">{a.co}</span>
                      </div>
                      <div className="bcard-role">{a.role}</div>
                      {a.next && (
                        <span className={"bcard-when " + a.status}>
                          <span className="d"></span>{a.next.when}
                        </span>
                      )}
                      <div className="bcard-ft">
                        <span className="src">{a.source === "—" ? "Not applied" : a.source}</span>
                        <span className={"ago" + (stale ? " red" : "")}>{agoText(a)}</span>
                      </div>
                    </div>
                  );
                })}
                {items.length === 0 && <div className="bcol-empty">Drop here</div>}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}

/* ═══ INSIGHTS ═══════════════════════════════════════════════ */
function InsightsScreen() {
  const D = window.INSIGHTS;
  const max = Math.max(...D.funnel.map((f) => f.n));
  const amax = Math.max(...D.activity);
  const weeks = ["", "", "Mar", "", "", "Apr", "", "", "", "May", "", ""];
  return (
    <div className="page ins">
      <h1>Insights</h1>
      <div className="sub">How your search is actually going · last 12 weeks</div>

      <div className="ins-kpis">
        <div className="kpi"><div className="l">Reply rate</div><div className="v">{D.replyRate}%</div><div className="delta up">{ICONS.arrow} up 9% vs last month</div></div>
        <div className="kpi"><div className="l">Avg. time to first reply</div><div className="v">{D.avgFirstReply}</div><div className="delta flat">about the same</div></div>
        <div className="kpi"><div className="l">Furthest stage reached</div><div className="v" style={{ fontSize: 26 }}>Offer ×1</div><div className="delta up">Vercel · decide Friday</div></div>
      </div>

      <div className="ins-grid">
        <div className="panel">
          <div className="ph">Pipeline funnel</div>
          <div className="psub">Where your 12 active applications sit today</div>
          <div className="funnel">
            {D.funnel.map((s) => (
              <div className="fn" key={s.stage}>
                <span className="l">{s.stage}</span>
                <span className="bar" style={{ width: (s.n / max * 100) + "%" }}></span>
                <span className="n">{s.n}</span>
              </div>
            ))}
          </div>
          <div style={{ height: 28 }}></div>
          <div className="ph">Application activity</div>
          <div className="psub">New touchpoints per week</div>
          <div className="bars">
            {D.activity.map((v, i) => <div className="b" key={i} style={{ height: (v / amax * 100) + "%" }}></div>)}
          </div>
          <div className="bars-x">{weeks.map((w, i) => <span key={i}>{w}</span>)}</div>
        </div>

        <div className="panel">
          <div className="ph">Where they come from</div>
          <div className="psub">Referrals are converting best</div>
          {D.sources.map((s) => (
            <div className="src-row" key={s.src}>
              <div className="meta">
                <span className="nm">{s.src}</span>
                <span className="track"><i style={{ width: s.pct + "%" }}></i></span>
              </div>
              <span className="n">{s.n}</span>
            </div>
          ))}
          <div style={{ marginTop: 20, padding: 16, background: "var(--accent-tint)", borderRadius: 12, display: "flex", gap: 10 }}>
            <span style={{ color: "var(--accent-text)", flexShrink: 0, marginTop: 1 }}><Spark s={15} /></span>
            <span style={{ fontSize: 12.5, lineHeight: 1.5, color: "var(--accent-text)" }}>Your <b>4 referral applications</b> reached a screen or further. Cold apps stall at the inbox — lean on intros.</span>
          </div>
        </div>
      </div>
    </div>
  );
}

Object.assign(window, { TodayScreen, BoardScreen, InsightsScreen });
