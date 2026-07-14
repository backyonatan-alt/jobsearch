<script>
  import PreviewBanner from '../PreviewBanner.svelte';
  const DECISIONS = [
    { what: 'Right rail', how: 'Deleted. “Your move” and “What you can do today” merge into one prioritized action list in the center — the page’s single spine. “Waiting to hear back” becomes a count in the pipeline pulse that filters the Applications list.', why: 'Feedback #2 — “no reason to compete with more things”; rail “too wide and sparse”.', strong: 'Deleted.' },
    { what: 'Board', how: 'Demoted from nav; flattened into the list. Stage becomes group labels + an inline status control on every row — same one-gesture status change the board offered, without a page.', why: 'Feedback #1 — “maybe you don’t need a board at all — just one list with a few kinds of action buttons.”' },
    { what: 'Exits (rejected · withdrawn · closed)', how: 'Archive section at the end of Applications — “the list of everything at the end.” Rows keep logo, what happened, and when; one click reopens if the req comes back. Exits stop being invisible.', why: 'Feedback #4 — closed applications currently have no surface at all.' },
    { what: 'Top numbers', how: 'Promoted to a pipeline pulse: big counts (to prep · to decide · to nudge · waiting · closed) rendered with company logos, each one clickable into the filtered list.', why: 'Feedback #3 + #5 — “numbers should carry more… with company logos, less list-like.”' },
    { what: '“What we’re noticing”', how: 'One card on Today — a single observation per day, never a list. Insights keeps the full history and the funnel (rebuilt as a proper stepped funnel).', why: 'Feedback #5 — “belongs on the main page, not buried in Insights.”' },
    { what: 'Prep & credits', how: 'Prep surfaces untouched. On home, prep is the top action when an interview is near, not a separate block competing with tasks. Wrong-company confirmation is designed at the playbook frame.', why: 'Locked strategy — prep is the wedge; the tracker just stays out of its way.' }
  ];
</script>

<svelte:head><title>Redesign preview — IA proposal</title></svelte:head>

<PreviewBanner tag="IA" title="IA rationale — one home, one list, exits included" />

<div class="page">
  <div style="max-width:920px;margin:0 auto;padding:64px 40px 96px">

    <div style="font-size:11px;font-weight:600;letter-spacing:.12em;text-transform:uppercase;color:#e0641f;margin-bottom:12px">Pursuit · IA proposal</div>
    <h1 style="font-size:40px;line-height:1.1;font-weight:700;letter-spacing:-0.02em;margin:0 0 16px">One home, one list, exits included.</h1>
    <p style="font-size:17px;line-height:1.6;color:#4b5158;margin:0 0 48px;max-width:680px">The tracker spine collapses from four competing surfaces (Today center, right rail, Board, Insights) into <strong>two</strong>: a task-first <strong>Home</strong> and a complete <strong>Applications</strong> list where every application — including rejected, withdrawn, and closed — has a permanent home. Prep stays exactly where it is.</p>

    <div style="display:grid;grid-template-columns:1fr 1fr;gap:20px;margin-bottom:20px">
      <div style="background:#fff;border:1px solid #e8e8e5;border-radius:12px;padding:24px">
        <div style="font-size:11px;font-weight:600;letter-spacing:.1em;text-transform:uppercase;color:#8a9099;margin-bottom:16px">Current — 4 surfaces, 6 lists</div>
        <div style="display:flex;flex-direction:column;gap:8px">
          <div style="border:1px solid #e8e8e5;border-radius:8px;padding:10px 14px;font-size:13px"><strong>Today</strong> — stats · prep block · what-you-can-do · recently added</div>
          <div style="border:1px solid #f0c9a8;background:#fdf6ef;border-radius:8px;padding:10px 14px;font-size:13px"><strong>Right rail</strong> — “Your move” + “Waiting to hear back” <span style="color:#c05310">← duplicates the center</span></div>
          <div style="border:1px solid #f0c9a8;background:#fdf6ef;border-radius:8px;padding:10px 14px;font-size:13px"><strong>Board</strong> — 5 columns, mostly empty <span style="color:#c05310">← 8 of 9 cards sit in one column</span></div>
          <div style="border:1px solid #e8e8e5;border-radius:8px;padding:10px 14px;font-size:13px"><strong>Insights</strong> — funnel · sources · “what we’re noticing”</div>
          <div style="border:1px dashed #d9958a;background:#fdf1ef;border-radius:8px;padding:10px 14px;font-size:13px;color:#a13a2c"><strong>Nowhere</strong> — rejected / withdrawn / closed are invisible</div>
        </div>
      </div>
      <div style="background:#fff;border:1px solid #e8e8e5;border-radius:12px;padding:24px">
        <div style="font-size:11px;font-weight:600;letter-spacing:.1em;text-transform:uppercase;color:#1d7a4f;margin-bottom:16px">Proposed — 2 tracker surfaces + prep</div>
        <div style="display:flex;flex-direction:column;gap:8px">
          <div style="border:1px solid #cfe0cd;background:#f4faf4;border-radius:8px;padding:10px 14px;font-size:13px"><strong>Home</strong> — only what needs you: pipeline pulse up top, prioritized actions, one “noticing” card. No right rail.</div>
          <div style="border:1px solid #cfe0cd;background:#f4faf4;border-radius:8px;padding:10px 14px;font-size:13px"><strong>Applications</strong> — the full list, grouped by stage. Archive section at the end for every exit.</div>
          <div style="border:1px solid #e8e8e5;border-radius:8px;padding:10px 14px;font-size:13px"><strong>Insights</strong> — kept, slimmer: real funnel + sources. “Noticing” observations move to Home.</div>
          <div style="border:1px solid #e8e8e5;border-radius:8px;padding:10px 14px;font-size:13px;color:#6f7680"><strong>Prep / Playbook</strong> — untouched. Detail page gains logo, JD summary, posting link, recruiter contact, wrong-company check.</div>
        </div>
      </div>
    </div>

    <div style="background:#fff;border:1px solid #e8e8e5;border-radius:12px;padding:28px;margin-bottom:20px">
      <h2 style="font-size:20px;font-weight:700;letter-spacing:-0.01em;margin:0 0 20px">Where everything goes — and why</h2>
      <div style="display:flex;flex-direction:column">
        {#each DECISIONS as d (d.what)}
          <div style="display:grid;grid-template-columns:180px 1fr 220px;gap:20px;padding:14px 0;border-top:1px solid #f0f0ed;font-size:13.5px;line-height:1.55">
            <div style="font-weight:600">{d.what}</div>
            <div style="color:#4b5158">{d.how}</div>
            <div style="color:#8a9099;font-size:12.5px">{d.why}</div>
          </div>
        {/each}
      </div>
    </div>

    <div style="background:#fff;border:1px solid #e8e8e5;border-radius:12px;padding:28px">
      <h2 style="font-size:20px;font-weight:700;letter-spacing:-0.01em;margin:0 0 6px">Resulting navigation</h2>
      <p style="font-size:13.5px;color:#6f7680;margin:0 0 20px">Three items, unchanged count — but each now has one unambiguous job.</p>
      <div style="display:grid;grid-template-columns:1fr 1fr 1fr;gap:16px">
        <div style="border:1px solid #e8e8e5;border-radius:10px;padding:18px">
          <div style="font-size:15px;font-weight:700;margin-bottom:6px">Home</div>
          <div style="font-size:13px;line-height:1.55;color:#4b5158">“What do I do right now?” Pipeline pulse · prioritized actions · one noticing card. Nothing else.</div>
        </div>
        <div style="border:1px solid #e8e8e5;border-radius:10px;padding:18px">
          <div style="font-size:15px;font-weight:700;margin-bottom:6px">Applications <span style="font-weight:400;color:#8a9099;font-size:12px">(was Board)</span></div>
          <div style="font-size:13px;line-height:1.55;color:#4b5158">“Everything, findable.” Full list grouped by stage — archive of exits at the end.</div>
        </div>
        <div style="border:1px solid #e8e8e5;border-radius:10px;padding:18px">
          <div style="font-size:15px;font-weight:700;margin-bottom:6px">Insights</div>
          <div style="font-size:13px;line-height:1.55;color:#4b5158">“How is the search going?” Real funnel + sources + reply stats. Observations surface on Home; deep-dive lives here.</div>
        </div>
      </div>
    </div>

  </div>
</div>

<style>
  .page { min-height: 100vh; background: #f6f6f3; color: #16181c; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; -webkit-font-smoothing: antialiased; }
  .page :global(a) { color: #2463eb; text-decoration: none; }
</style>
