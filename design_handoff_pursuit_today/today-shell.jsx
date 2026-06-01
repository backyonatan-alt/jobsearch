/* ───────────────────────────────────────────────────────────
   PURSUIT · Today — shell + shared primitives
   ICONS · Spark · Logo · Pill · AppShell (sidebar + topbar)
   Exposed on window for the screen files.
   ─────────────────────────────────────────────────────────── */

const ICONS = {
  today: <svg width="15" height="15" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5"><rect x="2.5" y="3" width="11" height="10" rx="2"/><path d="M2.5 6h11M6 1.5v3M10 1.5v3"/></svg>,
  board: <svg width="15" height="15" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5"><rect x="2.5" y="2.5" width="3" height="11" rx="1"/><rect x="6.5" y="2.5" width="3" height="8" rx="1"/><rect x="10.5" y="2.5" width="3" height="5" rx="1"/></svg>,
  insights: <svg width="15" height="15" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5"><path d="M2 13h12M4 11V7M8 11V3M12 11V8"/></svg>,
  search: <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5"><circle cx="7" cy="7" r="4.5"/><path d="M10.5 10.5L13 13"/></svg>,
  chev: <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5" style={{opacity:0.5}}><path d="M6 4l4 4-4 4"/></svg>,
  arrow: <svg width="14" height="14" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.6"><path d="M3 8h9M9 5l3 3-3 3"/></svg>,
  bell: <svg width="15" height="15" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5"><path d="M8 2a4 4 0 0 0-4 4c0 3-1.2 4-1.2 4h10.4S12 9 12 6a4 4 0 0 0-4-4zM6.5 13a1.6 1.6 0 0 0 3 0"/></svg>,
  note: <svg width="15" height="15" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5"><path d="M3 2.5h7l3 3V13a.5.5 0 0 1-.5.5h-9A.5.5 0 0 1 3 13V2.5z"/><path d="M9.5 2.5V6h3.5"/></svg>,
  cal: <svg width="15" height="15" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5"><rect x="2.5" y="3" width="11" height="10" rx="2"/><path d="M2.5 6h11M6 1.5v3M10 1.5v3"/></svg>,
  file: <svg width="15" height="15" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5"><path d="M4 2.5h5l3 3V13a.5.5 0 0 1-.5.5h-7A.5.5 0 0 1 4 13V2.5z"/><path d="M8.5 2.5V6H12"/></svg>,
  send: <svg width="14" height="14" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5"><path d="M14 2L7 9M14 2l-4.5 12-2.5-5-5-2.5L14 2z"/></svg>,
  refresh: <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5"><path d="M13 8a5 5 0 1 1-1.5-3.6M13 2v3h-3"/></svg>,
  link: <svg width="14" height="14" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5"><path d="M6.5 9.5l3-3M7 4.5l1-1a2.5 2.5 0 0 1 3.5 3.5l-1 1M9 11.5l-1 1A2.5 2.5 0 0 1 4.5 9l1-1"/></svg>,
  x: <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.6"><path d="M4 4l8 8M12 4l-8 8"/></svg>,
};

const Spark = ({ s = 15 }) => (
  <svg width={s} height={s} viewBox="0 0 16 16" fill="currentColor">
    <path d="M8 1l1.2 4.1L13.5 6 9.8 7.4 8 11.5 6.2 7.4 2.5 6l4.3-.9L8 1z" opacity="0.95"/>
    <circle cx="13" cy="12.5" r="1.3"/>
  </svg>
);

function Logo({ co, cls, size = 24 }) {
  return <span className={"logo " + cls} style={{ width: size, height: size, fontSize: Math.round(size * 0.42), borderRadius: Math.max(6, Math.round(size * 0.27)) }}>{co}</span>;
}

function Pill({ status, label }) {
  const s = status === "closed" ? "applied" : status;
  return <span className={"pill " + s}><span className="pdot"></span>{label}</span>;
}

const getApp = (id) => window.APPS.find((a) => a.id === id);

/* App shell: clickable sidebar, dynamic topbar, content slot */
function AppShell({ active, topbar, children, onNav }) {
  const nav = [
    { k: "today", ic: ICONS.today, lbl: "Today" },
    { k: "board", ic: ICONS.board, lbl: "Board", ct: "12" },
    { k: "insights", ic: ICONS.insights, lbl: "Insights" },
  ];
  return (
    <div className="fr">
      <aside className="sb">
        <div className="brand">
          <span className="mark"></span>
          <span className="name">Pursuit</span>
        </div>
        {nav.map((n) => (
          <div key={n.k} className={"nav " + (active === n.k ? "on" : "")} onClick={() => onNav(n.k)}>
            <span className="ic">{n.ic}</span>
            <span>{n.lbl}</span>
            <span className="ct">{n.ct || ""}</span>
          </div>
        ))}
        <div className="foot">
          <div className="prof">
            <span className="av">{window.USER.initials}</span>
            <span className="who">{window.USER.first} {window.USER.last}<small>{window.USER.email}</small></span>
            {ICONS.chev}
          </div>
        </div>
      </aside>
      <div className="col">
        <div className="tb">
          <span className="here">{topbar}</span>
          <div className="right">
            <div className="srch">{ICONS.search}<span>Search…</span><span className="kbd">⌘K</span></div>
            <button className="btn pri">New application<span className="kbd">⌘N</span></button>
          </div>
        </div>
        <div className="scroll">{children}</div>
      </div>
    </div>
  );
}

Object.assign(window, { ICONS, Spark, Logo, Pill, getApp, AppShell });
