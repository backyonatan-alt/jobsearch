<script>
  import PreviewBanner from '../PreviewBanner.svelte';
  import { logoLg } from '../fixtures.js';

  const INTERVIEWERS = [
    {
      init: 'RB', bg: '#eef4ff', color: '#2463eb',
      name: 'Rotem Bar', meta: '· VP Engineering · 6y at Kayma',
      brief: 'Owns the platform org your squads would build on. In past loops he runs a live prioritization case — expect to trade off two roadmaps under a deadline. He grades how you use engineers’ time, not your framework vocabulary.',
      asks: '“Walk me through a platform bet that failed.” · “How do you decide build vs. buy?”',
      angle: 'the data-pipeline consolidation story — engineers freed up, with the numbers.',
      sources: ['Kayma eng blog · 2026', 'LinkedIn', 'conference talk']
    },
    {
      init: 'NG', bg: '#eef7ef', color: '#1d7a4f',
      name: 'Nadav Gil', meta: '· Head of Design · 3y at Kayma',
      brief: 'Cares about discovery craft and hates polished-case theater — bring one real research-to-ship story with the messy middle intact. He tends to ask what you’d cut from your own proposal.',
      asks: '“Tell me about a time research changed your roadmap.” · “What would you cut from this plan?”',
      angle: 'the enterprise-onboarding study — where the first design died and why that was right.',
      sources: ['design team blog', 'LinkedIn']
    },
    {
      init: 'TA', bg: '#fdf3ec', color: '#c05310',
      name: 'Tal Amir', meta: '· Group PM · peer interview',
      brief: 'Your would-be peer. She published on outcome roadmaps — mirror that language and she’ll advocate for you. Peer interviews at Kayma carry veto weight; treat this as the round to win, not survive.',
      asks: '“How would we split the platform surface between our groups?”',
      angle: 'come with one genuine question about her squad’s roadmap — peers grade curiosity.',
      sources: ['“Outcome roadmaps” talk · 2025', 'LinkedIn']
    }
  ];
</script>

<svelte:head><title>Redesign preview — Panel brief</title></svelte:head>

<PreviewBanner tag="5c" title="Panel brief — the full reading page" note="Where all the depth lives; the detail page is just the door." />

<div class="page">
  <!-- compact top bar (no search / new-application on the reading page) -->
  <div style="background:#fff;border-bottom:1px solid #e8e8e5">
    <div style="max-width:1100px;margin:0 auto;display:flex;align-items:center;gap:24px;padding:12px 32px">
      <div style="display:flex;align-items:center;gap:9px;flex:none">
        <div style="width:24px;height:24px;border-radius:50%;border:2.5px solid #16181c;display:flex;align-items:center;justify-content:center"><div style="width:9px;height:9px;border-radius:50%;background:#2463eb"></div></div>
        <div style="font-size:16px;font-weight:700;letter-spacing:-0.01em">Pursuit</div>
      </div>
      <div style="display:flex;align-items:center;gap:4px">
        <a href="/preview/redesign/home" class="navlink" style="border-radius:8px;padding:7px 14px;font-size:13.5px;color:#4b5158">Home</a>
        <a href="/preview/redesign/applications" style="display:flex;align-items:center;gap:6px;background:#16181c;color:#fff;border-radius:8px;padding:7px 14px;font-size:13.5px;font-weight:600">Applications <span style="font-size:11.5px;opacity:.6">16</span></a>
        <a href="/preview/redesign/insights" class="navlink" style="border-radius:8px;padding:7px 14px;font-size:13.5px;color:#4b5158">Insights</a>
      </div>
      <div style="display:flex;align-items:center;gap:12px;margin-left:auto">
        <div title="Prep credits" style="border:1px solid #e8e8e5;border-radius:8px;padding:7px 11px;font-size:12px;color:#6f7680;flex:none">✦ <strong style="color:#16181c">7</strong><span style="color:#8a9099">/10</span></div>
        <div title="back.yonatan@gmail.com" style="width:30px;height:30px;border-radius:50%;background:#e0641f;color:#fff;display:flex;align-items:center;justify-content:center;font-size:12px;font-weight:700;flex:none">Y</div>
      </div>
    </div>
  </div>

  <div style="max-width:760px;margin:0 auto;padding:28px 32px 90px">

    <div style="font-size:13px;color:#8a9099;margin-bottom:24px"><a href="/preview/redesign/applications" style="color:#8a9099">Applications</a> <span style="margin:0 4px">/</span> <a href="/preview/redesign/detail" style="color:#8a9099">Kayma</a> <span style="margin:0 4px">/</span> <span style="color:#16181c;font-weight:600">Panel brief</span></div>

    <div style="display:flex;align-items:center;gap:14px;margin-bottom:6px">
      <div title="Kayma" style="width:40px;height:40px;border-radius:10px;border:1px solid #e8e8e5;background:#fff url({logoLg('kayma.com')}) center/contain no-repeat;flex:none"></div>
      <h1 style="font-size:28px;font-weight:700;letter-spacing:-0.02em;margin:0">Panel brief — tomorrow, 10:00</h1>
    </div>
    <div style="font-size:13px;color:#8a9099;margin-bottom:26px">Kayma · VP Product · 3 interviewers · 90 min · generated 19h ago · <a href="#top">refresh</a> · <a href="#top">print</a></div>

    <div style="background:#eef4ff;border:1px solid #cdddfb;border-radius:14px;padding:18px 22px;margin-bottom:34px">
      <div style="font-size:13px;font-weight:700;color:#2463eb;margin-bottom:6px">✦ The one thing to remember</div>
      <div style="font-size:15px;line-height:1.6;color:#1e3a6e">This panel is built to disagree with you on purpose. Ask discovery questions before proposing — it’s their known filter, flagged in both of your earlier debriefs.</div>
    </div>

    {#each INTERVIEWERS as p, i (p.init)}
      {#if i > 0}<div style="border-top:1px solid #e2e2de;margin:22px 0"></div>{/if}
      <div style="display:flex;gap:16px;margin-bottom:10px">
        <div style="width:44px;height:44px;border-radius:50%;background:{p.bg};color:{p.color};display:flex;align-items:center;justify-content:center;font-size:14px;font-weight:700;flex:none">{p.init}</div>
        <div style="flex:1;min-width:0">
          <div style="font-size:16.5px;font-weight:700">{p.name} <span style="font-size:12.5px;font-weight:400;color:#8a9099">{p.meta}</span> <a href="#top" style="font-size:12px;margin-left:6px">LinkedIn ↗</a></div>
          <div style="font-size:14px;line-height:1.65;color:#4b5158;margin:8px 0 12px">{p.brief}</div>
          <div style="font-size:13.5px;line-height:1.65;color:#4b5158;margin-bottom:8px"><strong style="color:#16181c">Likely asks:</strong> {p.asks}</div>
          <div style="font-size:13.5px;line-height:1.65;color:#4b5158;margin-bottom:8px"><strong style="color:#1d7a4f">Your angle:</strong> {p.angle}</div>
          <div style="font-size:12px;color:#8a9099">Sources: {#each p.sources as s, j (s)}{#if j > 0} · {/if}<a href="#top">{s}</a>{/each}</div>
        </div>
      </div>
    {/each}

    <div style="background:#fff;border:1px solid #e8e8e5;border-radius:14px;padding:20px 24px;margin-top:34px">
      <div style="font-size:13px;font-weight:700;margin-bottom:10px">Questions worth asking them</div>
      <div style="display:flex;flex-direction:column;gap:8px;font-size:13.5px;line-height:1.6;color:#4b5158">
        <div>· “What broke the last time the platform group and a product squad disagreed?” <span style="color:#8a9099">(Rotem)</span></div>
        <div>· “Where does discovery get skipped today, and who feels it?” <span style="color:#8a9099">(Nadav)</span></div>
        <div>· “What would make the first 90 days a win from your seat?” <span style="color:#8a9099">(Tal)</span></div>
      </div>
    </div>

    <div style="display:flex;align-items:center;gap:8px;flex-wrap:wrap;font-size:12px;color:#8a9099;margin-top:18px">
      <span style="font-weight:600;color:#6f7680">All sources</span>
      <a href="#top" style="border:1px solid #e8e8e5;border-radius:14px;padding:3px 10px;background:#fff">Kayma eng blog</a>
      <a href="#top" style="border:1px solid #e8e8e5;border-radius:14px;padding:3px 10px;background:#fff">Glassdoor · VP loop (TLV)</a>
      <a href="#top" style="border:1px solid #e8e8e5;border-radius:14px;padding:3px 10px;background:#fff">Your HM debrief · Jul 8</a>
      <a href="#top" style="border:1px solid #e8e8e5;border-radius:14px;padding:3px 10px;background:#fff">LinkedIn ×3</a>
      <a href="#top" style="border:1px solid #e8e8e5;border-radius:14px;padding:3px 10px;background:#fff">“Outcome roadmaps” · 2025</a>
    </div>
    <div style="font-size:12.5px;color:#8a9099;margin-top:22px;border-top:1px solid #e2e2de;padding-top:14px">After the round: a 20-second debrief sharpens your Final/CEO prep. We’ll ask you tomorrow afternoon.</div>

  </div>
</div>

<style>
  .page { min-height: 100vh; background: #f6f6f3; color: #16181c; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; -webkit-font-smoothing: antialiased; }
  .page :global(a) { color: #2463eb; text-decoration: none; }
  .navlink { color: #4b5158 !important; }
  .navlink:hover { background: #f0f0ed; }
</style>
