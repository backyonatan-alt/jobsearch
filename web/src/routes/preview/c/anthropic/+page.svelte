<script>
  import { ANTHROPIC_DETAIL, fmtLongDate } from '$lib/preview-data.js';
  const a = ANTHROPIC_DETAIL;
</script>

<svelte:head>
  <title>Letter — Anthropic envelope</title>
</svelte:head>

<article class="letter">
  <header class="envelope-head">
    <div class="stamp">
      <p class="stamp-kind">meeting</p>
      <p class="stamp-date">{fmtLongDate(a.applied_at)}</p>
    </div>
    <p class="to">
      <span class="caveat">To —</span> <em>{a.company}</em>
    </p>
    <h1>{a.company}</h1>
    <p class="role">{a.role}<span class="loc"> · {a.location}</span></p>
    <div class="next">
      <p class="next-label">Next</p>
      <p>{a.next_step}</p>
    </div>
  </header>

  <hr class="divider" />

  <section class="block dossier">
    <h2 class="block-title">A note on your interviewer</h2>
    <p class="hand caveat">— from the briefing —</p>
    <p class="who"><em>{a.dossier.name}</em>, {a.dossier.title}</p>
    <p>{a.dossier.summary}</p>

    <h3 class="sub">Recent traces</h3>
    <ul class="recent">
      {#each a.dossier.recent as t}
        <li><span class="when">{t.date}</span> — {t.text}</li>
      {/each}
    </ul>

    <h3 class="sub">Likely style</h3>
    <ul class="bullets">
      {#each a.dossier.style as s}<li>{s}</li>{/each}
    </ul>

    <h3 class="sub">Watch for</h3>
    <ul class="bullets warn">
      {#each a.dossier.watchfor as w}<li>{w}</li>{/each}
    </ul>
  </section>

  <hr class="divider" />

  <section class="block">
    <h2 class="block-title">Conversation so far</h2>
    <ol class="timeline">
      {#each a.timeline as t}
        <li>
          <span class="t-date">{fmtLongDate(t.date)}</span>
          <span class="t-kind">{t.kind}</span>
          <p>{t.text}</p>
        </li>
      {/each}
    </ol>
  </section>

  <hr class="divider" />

  <section class="block notes">
    <h2 class="block-title caveat">Notes to self</h2>
    {#each a.notes as n}
      <p class="note-line">— {n}</p>
    {/each}
  </section>

  <footer class="sign">
    <p class="caveat sign-off">Best of luck —</p>
    <p class="signature caveat">Y.</p>
  </footer>
</article>

<style>
  .envelope-head { position: relative; margin-bottom: 1rem; }
  .stamp {
    position: absolute; top: -1rem; right: 0;
    width: 110px; padding: .5rem;
    background: #f9e6c4;
    border: 1px dashed #aa2828;
    transform: rotate(4deg);
    text-align: center;
  }
  .stamp-kind { font-family: 'Spectral', serif; font-weight: 600; font-style: italic; color: #aa2828; font-size: 1.05rem; margin: 0; text-transform: uppercase; letter-spacing: .04em; }
  .stamp-date { font-family: 'JetBrains Mono', monospace; font-size: 10px; color: #6b3a3a; margin: .25rem 0 0; }

  .to { font-size: 1.1rem; color: #6b6358; margin: 0; }
  .to .caveat { font-size: 1.4rem; margin-right: .35rem; }
  h1 { font-family: 'Spectral', serif; font-weight: 500; font-style: italic; font-size: 3rem; line-height: 1; margin: .5rem 0; color: #1c3a72; letter-spacing: -.02em; }
  .role { margin: 0 0 1.5rem; color: #3d3830; }
  .loc { color: #6b6358; font-style: italic; }

  .next {
    border-left: 3px solid #aa2828;
    padding: .65rem 1rem;
    background: rgba(170, 40, 40, .05);
  }
  .next-label { font-family: 'JetBrains Mono', monospace; font-size: 10px; letter-spacing: .12em; text-transform: uppercase; color: #aa2828; margin: 0 0 .25rem; }
  .next p { margin: 0; }

  .divider { border: 0; text-align: center; margin: 2rem 0; }
  .divider::before { content: '✦   ·   ✦'; letter-spacing: .5em; font-size: .9rem; color: #8a7d5e; }

  .block-title { font-family: 'Spectral', serif; font-weight: 500; font-style: italic; font-size: 1.55rem; margin: 0 0 .75rem; color: #1c3a72; }
  .hand { font-size: 1.2rem; margin: 0 0 .75rem; color: #aa2828; }
  .who { margin: 0 0 .75rem; }
  .who em { font-weight: 600; font-style: italic; color: #1c3a72; }

  .sub { font-family: 'Caveat', cursive; font-weight: 700; font-size: 1.35rem; margin: 1.5rem 0 .5rem; color: #1c3a72; }

  .recent { list-style: none; padding: 0; margin: 0; }
  .recent li { padding: .35rem 0; border-bottom: 1px dotted rgba(31, 28, 23, .12); }
  .recent li:last-child { border-bottom: 0; }
  .recent .when { font-family: 'JetBrains Mono', monospace; font-size: 11px; color: #1c3a72; }

  .bullets { list-style: none; padding: 0; margin: 0; }
  .bullets li { padding: .35rem 0 .35rem 1.5rem; position: relative; color: #3d3830; }
  .bullets li::before { content: '✦'; position: absolute; left: 0; color: #1c3a72; font-size: .85rem; top: .55rem; }
  .bullets.warn li::before { content: '⚠'; color: #aa2828; }

  .timeline { list-style: none; padding: 0; margin: 0; }
  .timeline li { padding: .65rem 0; border-bottom: 1px dashed rgba(31, 28, 23, .1); }
  .timeline li:last-child { border-bottom: 0; }
  .t-date { font-family: 'JetBrains Mono', monospace; font-size: 11px; color: #6b6358; margin-right: .65rem; }
  .t-kind { font-style: italic; color: #1c3a72; font-size: .95rem; }
  .timeline p { margin: .35rem 0 0; color: #2c2820; }

  .notes .note-line {
    margin: .5rem 0;
    color: #3d3830;
    font-style: italic;
  }

  .sign { margin-top: 2rem; }
  .sign-off { font-size: 1.4rem; margin: 0; }
  .signature { font-size: 3rem; margin: 0; line-height: 1; }
</style>
