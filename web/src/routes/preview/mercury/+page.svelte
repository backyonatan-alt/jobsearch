<script>
  import { PREVIEW_APPS, fmtDate } from '$lib/preview-data.js';

  const STATUS = {
    wishlist:  { bg: '#f4f1ea', fg: '#75695a' },
    applied:   { bg: '#e8e8fe', fg: '#3b3bb5' },
    screen:    { bg: '#fdebd2', fg: '#8a4f0e' },
    interview: { bg: '#f1e1f6', fg: '#7c3aed' },
    offer:     { bg: '#daf1de', fg: '#1f6b3e' },
    rejected:  { bg: '#f7dedf', fg: '#9b2c2c' },
    withdrawn: { bg: '#ebe7df', fg: '#6b6358' }
  };

  const TONES = ['#7c3aed', '#dc2626', '#0891b2', '#16a34a', '#c2410c', '#a16207'];
  const tone = (s) => TONES[s.charCodeAt(0) % TONES.length];
</script>

<svelte:head>
  <title>Mercury sketch — Pursuit</title>
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous" />
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=Newsreader:opsz,wght@6..72,400;6..72,500;6..72,600&display=swap" rel="stylesheet" />
</svelte:head>

<div class="shell">
  <aside class="side">
    <div class="brand">
      <div class="bm">P</div>
      <span class="bn">Pursuit</span>
    </div>

    <nav>
      <a class="nl cur" href="#">Applications<span class="ct">6</span></a>
      <a class="nl" href="#">Board</a>
      <a class="nl" href="#">Funnel</a>
    </nav>

    <p class="g">Companies</p>
    <nav>
      <a class="nl" href="#"><span class="dot" style="background: {tone('A')}"></span>Anthropic</a>
      <a class="nl" href="#"><span class="dot" style="background: {tone('S')}"></span>Stripe</a>
      <a class="nl" href="#"><span class="dot" style="background: {tone('V')}"></span>Vercel</a>
      <a class="nl" href="#"><span class="dot" style="background: {tone('L')}"></span>Linear</a>
    </nav>

    <div class="flex"></div>
    <a class="leave" href="/preview">‹ leave preview</a>
  </aside>

  <main>
    <header class="top">
      <div class="crumb"><span>Workspace</span> · Applications</div>
      <div class="t-right">
        <button class="ghost">⌘ K</button>
        <div class="av">B</div>
      </div>
    </header>

    <div class="page">
      <div class="kicker">May 22, 2026</div>
      <h1 class="display">Good afternoon, Yonatan.</h1>
      <p class="lede">
        <span class="hl">24</span> applications in flight.
        <span class="hl">5</span> active loops.
        <span class="hl accent">1</span> open offer from Linear, replying by Friday.
      </p>

      <div class="kpis">
        <article class="kpi accent">
          <p class="kl">Reply rate</p>
          <p class="kv display">62<span class="pct">%</span></p>
          <p class="kd">↑ 8 points · last 30 days</p>
        </article>
        <article class="kpi">
          <p class="kl">Time to first reply</p>
          <p class="kv display">3.2<span class="pct">d</span></p>
          <p class="kd">median across applied</p>
        </article>
        <article class="kpi">
          <p class="kl">Active loops</p>
          <p class="kv display">5</p>
          <p class="kd">2 advancing this week</p>
        </article>
      </div>

      <section class="block">
        <header class="block-head">
          <h2>Applications</h2>
          <div class="b-actions">
            <button class="ghost">Filter</button>
            <button class="ghost">Sort</button>
            <button class="primary">+ New</button>
          </div>
        </header>

        <table>
          <thead>
            <tr>
              <th>Company</th>
              <th>Role</th>
              <th>Status</th>
              <th>Applied</th>
              <th>CV</th>
              <th>Source</th>
            </tr>
          </thead>
          <tbody>
            {#each PREVIEW_APPS as a}
              <tr>
                <td>
                  <span class="co">
                    <span class="mono" style="background: {tone(a.company)}">{a.company[0]}</span>
                    {a.company}
                  </span>
                </td>
                <td><div class="role">{a.role}</div><div class="loc">{a.location}</div></td>
                <td><span class="pill" style="background: {STATUS[a.status].bg}; color: {STATUS[a.status].fg}">{a.status}</span></td>
                <td class="muted">{fmtDate(a.applied_at)}</td>
                <td class="muted">{a.cv_variant ?? '—'}</td>
                <td class="muted">{a.source}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </section>
    </div>
  </main>
</div>

<style>
  :global(html, body) {
    background: #faf9f7; color: #1a1814;
    font-family: 'Inter', system-ui, sans-serif;
    font-size: 14px; line-height: 1.5;
    -webkit-font-smoothing: antialiased;
  }
  .shell { display: grid; grid-template-columns: 220px 1fr; min-height: 100vh; }

  .side { background: #f4f1ea; border-right: 1px solid #e6dfd0; padding: 1.5rem 1rem; display: flex; flex-direction: column; gap: 1rem; }
  .brand { display: flex; align-items: center; gap: .65rem; padding: 0 .35rem; }
  .bm { width: 30px; height: 30px; border-radius: 6px; background: #1a1814; color: #f4f1ea; display: grid; place-items: center; font-family: 'Newsreader', serif; font-weight: 600; }
  .bn { font-family: 'Newsreader', serif; font-weight: 500; font-size: 18px; color: #1a1814; letter-spacing: -.01em; }

  nav { display: flex; flex-direction: column; gap: 1px; }
  .g { margin: .25rem 0 .15rem .55rem; font-size: 11px; color: #75695a; font-weight: 500; letter-spacing: .04em; text-transform: uppercase; }
  .nl { display: flex; align-items: center; gap: .55rem; padding: .4rem .55rem; border-radius: 5px; color: #44402d; text-decoration: none; font-size: 13px; }
  .nl:hover { background: #eae3d0; color: #1a1814; }
  .nl.cur { background: #1a1814; color: #f4f1ea; font-weight: 500; }
  .nl .ct { margin-left: auto; font-size: 11px; color: #c1b8a3; }
  .nl.cur .ct { color: #c1b8a3; }
  .nl .dot { width: 8px; height: 8px; border-radius: 999px; }

  .flex { flex: 1; }
  .leave { color: #75695a; font-size: 11px; text-decoration: none; padding: .35rem .55rem; }

  main { display: flex; flex-direction: column; min-width: 0; }
  .top { display: flex; align-items: center; justify-content: space-between; padding: 1rem 2rem; border-bottom: 1px solid #e6dfd0; }
  .crumb { font-size: 13px; color: #75695a; }
  .crumb span { color: #1a1814; font-weight: 500; }
  .t-right { display: flex; align-items: center; gap: .65rem; }
  .ghost { padding: .35rem .75rem; background: #fff; border: 1px solid #e6dfd0; border-radius: 6px; font: inherit; font-size: 12px; color: #44402d; cursor: pointer; }
  .ghost:hover { background: #faf9f7; }
  .av { width: 30px; height: 30px; border-radius: 999px; background: linear-gradient(135deg, #7c3aed, #dc2626); color: #fff; display: grid; place-items: center; font-weight: 600; font-size: 12px; }

  .page { padding: 3rem 2.5rem 4rem; max-width: 1200px; width: 100%; }
  .kicker { font-size: 11px; letter-spacing: .08em; text-transform: uppercase; color: #75695a; margin-bottom: .5rem; font-weight: 600; }
  h1.display { font-family: 'Newsreader', Georgia, serif; font-weight: 500; font-size: 38px; line-height: 1.1; margin: 0 0 .85rem; letter-spacing: -.02em; color: #1a1814; }
  .lede { font-size: 17px; color: #44402d; margin: 0 0 2.5rem; max-width: 60ch; line-height: 1.55; }
  .lede .hl { font-family: 'Newsreader', serif; font-weight: 600; font-size: 19px; color: #1a1814; }
  .lede .hl.accent { color: #7c3aed; }

  .kpis { display: grid; grid-template-columns: repeat(3, 1fr); gap: 1rem; margin-bottom: 3rem; }
  @media (max-width: 800px) { .kpis { grid-template-columns: 1fr; } }
  .kpi { background: #fff; border: 1px solid #e6dfd0; border-radius: 10px; padding: 1.25rem 1.35rem; }
  .kpi.accent { background: linear-gradient(135deg, #ffffff 0%, #f4eafe 100%); border-color: #d8c3f3; }
  .kl { margin: 0; font-size: 12px; color: #75695a; font-weight: 500; }
  .kv { margin: .35rem 0 .35rem; font-family: 'Newsreader', serif; font-weight: 500; font-size: 44px; line-height: 1; color: #1a1814; letter-spacing: -.02em; }
  .kv .pct { font-size: 26px; color: #75695a; margin-left: 2px; }
  .kpi.accent .kv { color: #7c3aed; }
  .kd { margin: 0; font-size: 12px; color: #75695a; }

  .block { background: #fff; border: 1px solid #e6dfd0; border-radius: 12px; overflow: hidden; }
  .block-head { display: flex; align-items: center; justify-content: space-between; padding: 1rem 1.5rem; border-bottom: 1px solid #e6dfd0; }
  .block-head h2 { font-family: 'Newsreader', serif; font-weight: 500; font-size: 22px; margin: 0; letter-spacing: -.01em; color: #1a1814; }
  .b-actions { display: flex; gap: .35rem; }
  .primary { padding: .35rem .85rem; background: #1a1814; color: #f4f1ea; border: 0; border-radius: 6px; font: inherit; font-size: 12px; font-weight: 500; cursor: pointer; }

  table { width: 100%; border-collapse: collapse; }
  th { text-align: left; padding: .65rem 1rem; font-size: 11px; font-weight: 600; color: #75695a; text-transform: uppercase; letter-spacing: .04em; border-bottom: 1px solid #e6dfd0; }
  td { padding: .85rem 1rem; border-bottom: 1px solid #f3eee2; font-size: 13px; vertical-align: middle; color: #1a1814; }
  tbody tr:hover { background: #faf9f7; }
  tbody tr:last-child td { border-bottom: 0; }
  .co { display: inline-flex; align-items: center; gap: .65rem; font-weight: 500; }
  .mono { width: 28px; height: 28px; border-radius: 6px; color: #fff; display: grid; place-items: center; font-weight: 600; font-size: 13px; font-family: 'Newsreader', serif; }
  .role { color: #1a1814; }
  .loc { color: #a39d92; font-size: 11px; margin-top: 1px; }
  .pill { padding: 2px 9px; border-radius: 4px; font-size: 11px; font-weight: 500; text-transform: capitalize; }
  .muted { color: #75695a; }
</style>
