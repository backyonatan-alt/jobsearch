<script>
  import { ANTHROPIC_DETAIL, fmtLongDate } from '$lib/preview-data.js';
  const a = ANTHROPIC_DETAIL;
</script>

<svelte:head>
  <title>Garden — Anthropic</title>
</svelte:head>

<section class="hero">
  <p class="eyebrow">Currently flowering · {a.location}</p>
  <h1>{a.company}</h1>
  <p class="role">{a.role}</p>
  <div class="next-card">
    <div>
      <p class="next-label">Next step</p>
      <p class="next-body">{a.next_step}</p>
    </div>
    <button class="btn-pill">Prep notes ›</button>
  </div>
</section>

<section class="dossier">
  <header class="dossier-head">
    <div class="avatar">D</div>
    <div>
      <p class="who">{a.dossier.name}</p>
      <p class="who-title">{a.dossier.title}</p>
    </div>
    <span class="ai-pill">AI dossier</span>
  </header>
  <p class="summary">{a.dossier.summary}</p>

  <div class="dossier-grid">
    <div class="block">
      <h3>Recent</h3>
      <ul class="recent">
        {#each a.dossier.recent as t}
          <li>
            <span class="when">{t.date}</span>
            <span>{t.text}</span>
          </li>
        {/each}
      </ul>
    </div>
    <div class="block">
      <h3>Likely style</h3>
      <ul class="bullets">
        {#each a.dossier.style as s}<li>{s}</li>{/each}
      </ul>
      <h3 style="margin-top: 1.5rem">Watch for</h3>
      <ul class="bullets warn">
        {#each a.dossier.watchfor as w}<li>{w}</li>{/each}
      </ul>
    </div>
  </div>
</section>

<section class="timeline-block">
  <h2>Conversation so far</h2>
  <ol class="timeline">
    {#each a.timeline as t}
      <li>
        <span class="dot" data-kind={t.kind}></span>
        <div>
          <p class="t-meta">
            <span class="t-when">{fmtLongDate(t.date)}</span>
            <span class="t-kind">{t.kind}</span>
          </p>
          <p class="t-text">{t.text}</p>
        </div>
      </li>
    {/each}
  </ol>
</section>

<section class="notes-block">
  <h2>Notes to self</h2>
  <ul class="notes">
    {#each a.notes as n}
      <li>{n}</li>
    {/each}
  </ul>
</section>

<p class="back-row"><a href="/preview/b">‹ back to the garden</a></p>

<style>
  .hero { margin: 1rem 0 2.5rem; }
  .hero h1 { font-size: 3rem; margin: .25rem 0 .25rem; line-height: 1; }
  .role { color: #5a5550; font-size: 1.2rem; margin: 0 0 1.5rem; }

  .next-card {
    background: linear-gradient(135deg, #f8c89d 0%, #f4b577 100%);
    border-radius: 18px;
    padding: 1.25rem 1.5rem;
    display: flex; align-items: center; justify-content: space-between; gap: 1.5rem;
    color: #4a2614;
  }
  .next-label { font-size: 11px; letter-spacing: .14em; text-transform: uppercase; font-weight: 700; margin: 0; opacity: .8; }
  .next-body { margin: .35rem 0 0; font-weight: 500; }
  .btn-pill { background: #4a2614; color: #fff; border: 0; border-radius: 999px; padding: .65rem 1.25rem; font-weight: 600; font-size: 14px; cursor: pointer; white-space: nowrap; }

  .dossier {
    background: rgba(255,255,255,.5);
    border-radius: 18px;
    padding: 1.5rem 1.75rem;
    margin-bottom: 2.5rem;
    border: 1px solid rgba(107, 61, 95, .12);
  }
  .dossier-head { display: flex; align-items: center; gap: 1rem; margin-bottom: 1rem; }
  .avatar {
    width: 48px; height: 48px;
    background: linear-gradient(135deg, #e4d0db, #c2a8b8);
    color: #6b3d5f;
    border-radius: 14px;
    display: grid; place-items: center;
    font-family: 'Fraunces', serif; font-weight: 600; font-size: 22px;
  }
  .who { margin: 0; font-weight: 600; }
  .who-title { color: #837e75; font-size: 13px; margin: .15rem 0 0; }
  .ai-pill {
    margin-left: auto;
    background: #6b3d5f; color: #fff;
    padding: .25rem .7rem; border-radius: 999px;
    font-size: 11px; letter-spacing: .12em; text-transform: uppercase; font-weight: 700;
  }
  .summary { color: #3d3830; line-height: 1.6; margin: 0 0 1.5rem; }

  .dossier-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 2rem; }
  @media (max-width: 700px) { .dossier-grid { grid-template-columns: 1fr; } }
  .block h3 { font-size: 1.1rem; margin: 0 0 .65rem; }
  .recent { list-style: none; padding: 0; margin: 0; }
  .recent li { padding: .5rem 0; border-bottom: 1px dashed rgba(0,0,0,.08); display: grid; grid-template-columns: 4rem 1fr; gap: .65rem; font-size: 14px; }
  .recent li:last-child { border-bottom: 0; }
  .recent .when { font-weight: 600; color: #6b3d5f; font-size: 11px; letter-spacing: .08em; text-transform: uppercase; padding-top: .15rem; }

  .bullets { list-style: none; padding: 0; margin: 0; }
  .bullets li { padding: .35rem 0 .35rem 1.25rem; position: relative; color: #3d3830; font-size: 14px; }
  .bullets li::before { content: '🌿'; position: absolute; left: 0; }
  .bullets.warn li::before { content: '⚠️'; }

  .timeline-block, .notes-block { margin-bottom: 2.5rem; }
  .timeline-block h2, .notes-block h2 { font-size: 1.6rem; margin: 0 0 1rem; }
  .timeline { list-style: none; padding: 0; margin: 0; }
  .timeline li { display: grid; grid-template-columns: 1rem 1fr; gap: 1rem; padding: .75rem 0; border-bottom: 1px dashed rgba(0,0,0,.08); }
  .timeline li:last-child { border-bottom: 0; }
  .timeline .dot { width: 12px; height: 12px; border-radius: 999px; background: #c2dbb8; margin-top: .35rem; }
  .timeline .dot[data-kind="applied"] { background: #e2d5f0; }
  .timeline .dot[data-kind="screen"] { background: #f5d8d8; }
  .timeline .dot[data-kind="interview"] { background: #f8c89d; }
  .timeline .dot[data-kind="reply"] { background: #cad8c0; }
  .t-meta { margin: 0; font-size: 11px; letter-spacing: .08em; text-transform: uppercase; color: #837e75; font-weight: 600; }
  .t-kind { margin-left: .65rem; color: #6b3d5f; }
  .t-text { margin: .35rem 0 0; color: #3d3830; }

  .notes { list-style: none; padding: 0; margin: 0; }
  .notes li {
    background: rgba(255,255,255,.55);
    border-radius: 14px;
    padding: 1rem 1.25rem;
    margin-bottom: .65rem;
    color: #3d3830;
    border-left: 3px solid #c2dbb8;
  }

  .back-row { text-align: center; margin: 2rem 0; }
  .back-row a { color: #6b3d5f; text-decoration: none; font-weight: 500; }
  .back-row a:hover { text-decoration: underline; }
</style>
