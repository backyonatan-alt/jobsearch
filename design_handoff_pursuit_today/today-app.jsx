/* ───────────────────────────────────────────────────────────
   PURSUIT · Today — app router
   Screen stack + back, sidebar nav, follow-up toast.
   ─────────────────────────────────────────────────────────── */

function App() {
  const [view, setView] = React.useState({ name: "today", params: {} });
  const [stack, setStack] = React.useState([]);
  const [active, setActive] = React.useState("today");
  const [toast, setToast] = React.useState(null);

  const nav = (name, params = {}) => {
    setStack((s) => [...s, view]);
    setView({ name, params });
    if (name === "today" || name === "board" || name === "insights") setActive(name);
    document.querySelector(".scroll")?.scrollTo(0, 0);
  };
  const onNavSide = (k) => { setStack([]); setView({ name: k, params: {} }); setActive(k); };
  const back = () => {
    setStack((s) => {
      if (!s.length) return s;
      const prev = s[s.length - 1];
      setView(prev);
      if (["today", "board", "insights"].includes(prev.name)) setActive(prev.name);
      return s.slice(0, -1);
    });
  };

  const logFollowup = () => {
    setToast("Follow-up logged — we've reset the clock");
    setTimeout(() => setToast(null), 2600);
  };

  /* topbar content */
  const back_btn = <button className="back" onClick={back}>{ICONS.arrow} Back</button>;
  let topbar;
  if (view.name === "today") topbar = "Today";
  else if (view.name === "board") topbar = "Board";
  else if (view.name === "insights") topbar = "Insights";
  else if (view.name === "dossier") topbar = <>{back_btn}<span className="crumb-sep">/</span> Playbook</>;
  else if (view.name === "detail") topbar = <>{back_btn}<span className="crumb-sep">/</span> {window.getApp(view.params.appId).co}</>;

  let screen;
  if (view.name === "today") screen = <TodayScreen nav={nav} />;
  else if (view.name === "board") screen = <BoardScreen nav={nav} filter={view.params.filter} key={view.params.filter || "all"} />;
  else if (view.name === "insights") screen = <InsightsScreen />;
  else if (view.name === "dossier") screen = <DossierScreen />;
  else if (view.name === "detail") screen = <DetailScreen appId={view.params.appId} nav={nav} onLog={logFollowup} key={view.params.appId} />;

  return (
    <div className="app-root">
      <AppShell active={active} topbar={topbar} onNav={onNavSide}>
        {screen}
      </AppShell>
      {toast && <div className="toast"><span className="ok"><svg width="10" height="10" viewBox="0 0 16 16" fill="none" stroke="#fff" strokeWidth="2.4"><path d="M3 8.5l3.5 3.5L13 4"/></svg></span>{toast}</div>}
    </div>
  );
}

ReactDOM.createRoot(document.getElementById("root")).render(<App />);
