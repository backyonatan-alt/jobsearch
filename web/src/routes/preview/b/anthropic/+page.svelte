<script>
  import { ANTHROPIC_DETAIL, fmtLongDate } from '$lib/preview-data.js';
  const a = ANTHROPIC_DETAIL;
  let tab = $state('dossier');
</script>

<svelte:head>
  <title>Anthropic — Pursuit</title>
</svelte:head>

<p class="breadcrumb">
  <a href="/preview/b">Workspace</a>
  <span>/</span>
  <a href="/preview/b">Applications</a>
  <span>/</span>
  Anthropic
</p>

<header class="profile-head">
  <div class="left">
    <div class="mono">A</div>
    <div>
      <h1>{a.company}</h1>
      <p class="role">{a.role} <span class="dot">·</span> {a.location}</p>
      <div class="meta-row">
        <span class="status">
          <span class="status-dot"></span>
          interview · round 1 complete
        </span>
        <span class="sep">·</span>
        <span class="muted">Applied {fmtLongDate(a.applied_at)}</span>
        <span class="sep">·</span>
        <a class="link" href={a.jd_url}>Job description ↗</a>
      </div>
    </div>
  </div>
  <div class="right">
    <button class="btn-ghost">Add note</button>
    <button class="btn-primary">Update status</button>
  </div>
</header>

<div class="next-banner">
  <div>
    <p class="next-kicker">Up next</p>
    <p class="next-body">{a.next_step}</p>
  </div>
  <button class="btn-primary sm">Open prep notes →</button>
</div>

<nav class="tabs">
  <button class="tab" class:active={tab === 'dossier'}  onclick={() => (tab = 'dossier')}>Dossier <span class="ai-pill">AI</span></button>
  <button class="tab" class:active={tab === 'timeline'} onclick={() => (tab = 'timeline')}>Timeline</button>
  <button class="tab" class:active={tab === 'notes'}    onclick={() => (tab = 'notes')}>Notes</button>
  <button class="tab" class:active={tab === 'files'}    onclick={() => (tab = 'files')}>Files</button>
</nav>

{#if tab === 'dossier'}
  <div class="dossier-grid">
    <article class="panel main">
      <header class="panel-head">
        <div class="who">
          <div class="who-avatar">D</div>
          <div>
            <p class="who-name">{a.dossier.name}</p>
            <p class="who-title">{a.dossier.title}</p>
          </div>
        </div>
        <span class="generated">Updated 12m ago</span>
      </header>
      <p class="summary">{a.dossier.summary}</p>

      <h3 class="sub-title">Recent</h3>
      <ul class="recent">
        {#each a.dossier.recent as t}
          <li>
            <span class="when">{t.date}</span>
            <p>{t.text}</p>
          </li>
        {/each}
      </ul>
    </article>

    <aside class="side">
      <article class="panel">
        <h3 class="panel-title">Likely style</h3>
        <ul class="bullets">
          {#each a.dossier.style as s}<li>{s}</li>{/each}
        </ul>
      </article>
      <article class="panel warn">
        <h3 class="panel-title">Watch for</h3>
        <ul class="bullets">
          {#each a.dossier.watchfor as w}<li>{w}</li>{/each}
        </ul>
      </article>
    </aside>
  </div>
{:else if tab === 'timeline'}
  <article class="panel">
    <ol class="timeline">
      {#each a.timeline as t}
        <li>
          <span class="t-marker"></span>
          <div class="t-body">
            <div class="t-row">
              <span class="t-kind">{t.kind}</span>
              <span class="t-when">{fmtLongDate(t.date)}</span>
            </div>
            <p>{t.text}</p>
          </div>
        </li>
      {/each}
    </ol>
  </article>
{:else if tab === 'notes'}
  <div class="notes">
    {#each a.notes as n, i}
      <article class="note">
        <p>{n}</p>
        <p class="note-meta">Note {i + 1} · 4 days ago</p>
      </article>
    {/each}
    <button class="add-note">+ Add a note</button>
  </div>
{:else}
  <div class="empty-files">
    <p>No files attached yet.</p>
    <button class="btn-ghost">+ Attach CV variant</button>
  </div>
{/if}

<style>
  .profile-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 1rem; margin-bottom: 1.5rem; }
  .left { display: flex; gap: 1rem; }
  .mono {
    width: 56px; height: 56px; border-radius: 14px;
    background: linear-gradient(135deg, #ff8a5b 0%, #c45ba8 100%);
    color: #fff;
    display: grid; place-items: center;
    font-weight: 600; font-size: 24px;
    box-shadow: 0 2px 6px rgba(196, 91, 168, .25);
  }
  .role { color: #4a4842; margin: .25rem 0 .65rem; font-size: 14px; }
  .role .dot { color: #c4bda9; margin: 0 .35rem; }
  .meta-row { display: flex; align-items: center; gap: .65rem; font-size: 12px; }
  .status {
    display: inline-flex; align-items: center; gap: .4rem;
    color: #c45ba8; font-weight: 500;
  }
  .status-dot { width: 7px; height: 7px; border-radius: 999px; background: #c45ba8; box-shadow: 0 0 0 3px rgba(196, 91, 168, .15); }
  .sep { color: #c4bda9; }
  .muted { color: #71717a; }
  .link { color: #c45ba8; text-decoration: none; }
  .link:hover { text-decoration: underline; }

  .right { display: flex; gap: .5rem; }
  .btn-ghost {
    padding: .45rem .85rem;
    background: #fbf9f4;
    border: 1px solid #ebe6dd; border-radius: 6px;
    font: inherit; font-size: 13px; color: #4a4842; cursor: pointer;
  }
  .btn-ghost:hover { background: #ece5d4; }
  .btn-primary {
    padding: .45rem 1rem;
    background: linear-gradient(135deg, #18181b 0%, #2c2a32 100%);
    color: #fbf9f4;
    border: 0; border-radius: 6px;
    font: inherit; font-size: 13px; font-weight: 500;
    cursor: pointer;
    box-shadow: 0 1px 2px rgba(0,0,0,.06);
  }
  .btn-primary.sm { padding: .35rem .85rem; font-size: 12px; }

  .next-banner {
    display: flex; justify-content: space-between; align-items: center; gap: 1rem;
    background: linear-gradient(135deg, rgba(255, 138, 91, .12) 0%, rgba(196, 91, 168, .12) 100%);
    border: 1px solid rgba(196, 91, 168, .2);
    border-radius: 10px;
    padding: 1rem 1.25rem;
    margin-bottom: 1.5rem;
  }
  .next-kicker { margin: 0; font-size: 10px; font-weight: 600; color: #c45ba8; letter-spacing: .12em; text-transform: uppercase; }
  .next-body { margin: .25rem 0 0; color: #18181b; font-weight: 500; }

  .tabs {
    display: flex; gap: 0;
    border-bottom: 1px solid #ebe6dd;
    margin-bottom: 1.5rem;
  }
  .tab {
    background: transparent; border: 0;
    padding: .65rem 1rem;
    color: #71717a;
    font: inherit; font-size: 13px;
    cursor: pointer;
    border-bottom: 2px solid transparent;
    margin-bottom: -1px;
    display: inline-flex; align-items: center; gap: .4rem;
  }
  .tab:hover { color: #18181b; }
  .tab.active { color: #18181b; border-bottom-color: #c45ba8; font-weight: 500; }
  .ai-pill {
    background: linear-gradient(135deg, #ff8a5b, #c45ba8);
    color: #fff;
    font-size: 9px; font-weight: 700;
    padding: 1px 5px; border-radius: 4px;
    letter-spacing: .04em;
  }

  .dossier-grid { display: grid; grid-template-columns: 1fr 320px; gap: 1rem; }
  @media (max-width: 900px) { .dossier-grid { grid-template-columns: 1fr; } }

  .panel {
    background: #fbf9f4;
    border: 1px solid #ebe6dd;
    border-radius: 10px;
    padding: 1.25rem 1.5rem;
  }
  .panel.warn { background: linear-gradient(135deg, rgba(245, 158, 11, .06), rgba(245, 158, 11, .02)); border-color: rgba(245, 158, 11, .25); }
  .panel-head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem; }
  .who { display: flex; align-items: center; gap: .65rem; }
  .who-avatar {
    width: 36px; height: 36px; border-radius: 10px;
    background: linear-gradient(135deg, #7e6cd6, #c45ba8);
    color: #fff;
    display: grid; place-items: center;
    font-weight: 600; font-size: 14px;
  }
  .who-name { margin: 0; font-weight: 600; font-size: 14px; color: #18181b; }
  .who-title { margin: .15rem 0 0; font-size: 12px; color: #71717a; }
  .generated {
    font-family: 'JetBrains Mono', monospace;
    font-size: 10px; color: #a39d92;
    background: #fff; padding: 2px 6px; border-radius: 4px;
  }
  .summary { margin: 0 0 1.25rem; color: #4a4842; line-height: 1.6; font-size: 13px; }

  .sub-title, .panel-title {
    font-size: 11px; font-weight: 600;
    letter-spacing: .06em; text-transform: uppercase;
    color: #71717a;
    margin: 1rem 0 .5rem;
  }
  .panel-title { margin-top: 0; margin-bottom: .65rem; }
  .panel.warn .panel-title { color: #b45309; }

  .recent { list-style: none; padding: 0; margin: 0; }
  .recent li {
    display: grid; grid-template-columns: 4rem 1fr; gap: .85rem;
    padding: .65rem 0;
    border-top: 1px dashed #ebe6dd;
  }
  .recent li:first-of-type { border-top: 0; padding-top: 0; }
  .recent .when {
    font-family: 'JetBrains Mono', monospace;
    font-size: 11px; color: #c45ba8;
    padding-top: .15rem;
  }
  .recent p { margin: 0; font-size: 13px; color: #4a4842; line-height: 1.5; }

  .side { display: flex; flex-direction: column; gap: 1rem; }
  .bullets { list-style: none; padding: 0; margin: 0; }
  .bullets li {
    position: relative;
    padding: .35rem 0 .35rem 1rem;
    font-size: 12px; color: #4a4842;
    line-height: 1.5;
  }
  .bullets li::before {
    content: ''; position: absolute; left: 0; top: .7rem;
    width: 5px; height: 5px; border-radius: 999px;
    background: #c45ba8;
  }
  .panel.warn .bullets li::before { background: #f59e0b; }

  .timeline { list-style: none; padding: 0; margin: 0; }
  .timeline li { display: grid; grid-template-columns: 16px 1fr; gap: 1rem; padding: .65rem 0; position: relative; }
  .timeline li::before {
    content: ''; position: absolute;
    left: 7px; top: 1.6rem; bottom: -0.2rem;
    width: 2px; background: #ebe6dd;
  }
  .timeline li:last-child::before { display: none; }
  .t-marker {
    width: 16px; height: 16px; border-radius: 999px;
    background: #fff; border: 2px solid #c45ba8;
    margin-top: .25rem;
  }
  .t-row { display: flex; align-items: baseline; gap: .65rem; margin-bottom: .25rem; }
  .t-kind { font-weight: 500; color: #18181b; font-size: 13px; text-transform: capitalize; }
  .t-when {
    font-family: 'JetBrains Mono', monospace;
    font-size: 11px; color: #71717a;
  }
  .t-body p { margin: 0; color: #4a4842; line-height: 1.55; font-size: 13px; }

  .notes { display: flex; flex-direction: column; gap: .65rem; }
  .note {
    background: #fbf9f4;
    border: 1px solid #ebe6dd;
    border-left: 3px solid #c45ba8;
    border-radius: 8px;
    padding: 1rem 1.25rem;
  }
  .note p { margin: 0; color: #18181b; line-height: 1.55; font-size: 13px; }
  .note-meta { margin-top: .5rem !important; font-size: 11px; color: #a39d92 !important; }
  .add-note {
    background: transparent;
    border: 1px dashed #ebe6dd;
    color: #a39d92;
    padding: .65rem;
    border-radius: 8px;
    font: inherit;
    cursor: pointer;
  }
  .add-note:hover { background: #fbf9f4; color: #c45ba8; border-color: #c45ba8; }

  .empty-files {
    background: #fbf9f4;
    border: 1px solid #ebe6dd;
    border-radius: 10px;
    padding: 3rem 2rem;
    text-align: center;
    color: #71717a;
  }
  .empty-files p { margin: 0 0 1rem; }
</style>
