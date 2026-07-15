<script>
  import PreviewBanner from '../PreviewBanner.svelte';
  import TopNav from '../TopNav.svelte';
  import { logo } from '../fixtures.js';

  const chip = (d) => `width:20px;height:20px;border-radius:6px;border:2px solid #fff;box-shadow:0 0 0 1px #eeeeea;background:#fff url(${logo(d)}) center/contain no-repeat`;
  const SOURCES = [
    { name: 'LinkedIn', meta: '8 · 1 screen', w: 50, color: '#2463eb' },
    { name: 'Referral', meta: '3 · 1 interview, 1 offer', w: 19, color: '#16a34a' },
    { name: 'Company site', meta: '2', w: 13, color: '#8fabef' },
    { name: 'Direct / other', meta: '3', w: 19, color: '#c8ccd2' }
  ];
  const NOTICED = [
    { d: 'Jul 14', lead: 'Referrals are working.', rest: 'Your one referral reached interview; ask for one more intro this week.' },
    { d: 'Jul 11', lead: 'HiBob and Rylo have gone quiet.', rest: 'No reply in 3+ weeks — a direct nudge beats waiting.' },
    { d: 'Jul 8', lead: 'Your screen → interview rate is strong (50%).', rest: 'The bottleneck is getting replies, not passing rounds.' },
    { d: 'Jul 2', lead: 'Applications sent Mon–Tue got replies 2× more often', rest: 'than end-of-week ones in your history.' }
  ];
</script>

<svelte:head><title>Redesign preview — Insights</title></svelte:head>

<PreviewBanner tag="3.4" title="Insights — real funnel" note="Mixpanel-grade funnel with drop-off; “noticing” lives on Home, its history lives here." />

<div class="page">
  <TopNav active="insights" primaryCta={false} />

  <div style="max-width:1000px;margin:0 auto;padding:36px 32px 80px">

    <div style="display:flex;align-items:baseline;justify-content:space-between;margin-bottom:26px">
      <div>
        <h1 style="font-size:30px;font-weight:700;letter-spacing:-0.02em;margin:0 0 6px">Insights.</h1>
        <div style="font-size:13.5px;color:#6f7680">How your search is actually going · last 12 weeks</div>
      </div>
      <div style="font-size:13px;color:#8a9099">Last 12 weeks ▾</div>
    </div>

    <!-- headline numbers -->
    <div style="display:grid;grid-template-columns:1fr 1fr 1fr;border-top:1px solid #e2e2de;margin-bottom:44px">
      <div style="padding:24px 26px 0;border-right:1px solid #eeeeea">
        <div style="font-size:40px;font-weight:700;letter-spacing:-0.03em;line-height:1">25%</div>
        <div style="font-size:13px;font-weight:600;margin:9px 0 4px">reply rate</div>
        <div style="font-size:12px;color:#8a9099">4 of 16 got a reply</div>
      </div>
      <div style="padding:24px 26px 0;border-right:1px solid #eeeeea">
        <div style="font-size:40px;font-weight:700;letter-spacing:-0.03em;line-height:1">22d</div>
        <div style="font-size:13px;font-weight:600;margin:9px 0 4px">to first reply</div>
        <div style="font-size:12px;color:#8a9099">median, applied → screen</div>
      </div>
      <div style="padding:24px 26px 0">
        <div style="font-size:40px;font-weight:700;letter-spacing:-0.03em;line-height:1;color:#16a34a">Offer</div>
        <div style="font-size:13px;font-weight:600;margin:9px 0 4px">furthest stage</div>
        <div style="display:flex;align-items:center;gap:7px;font-size:12px;color:#8a9099"><span title="Linear" style={chip('linear.app')}></span>Linear · $210k base</div>
      </div>
    </div>

    <!-- funnel -->
    <div style="background:#fff;border:1px solid #e8e8e5;border-radius:16px;padding:26px 30px;margin-bottom:20px;box-shadow:0 1px 3px rgba(22,24,28,.04)">
      <div style="display:flex;align-items:baseline;gap:10px;margin-bottom:26px">
        <div style="font-size:16px;font-weight:700">Pipeline funnel</div>
        <div style="font-size:12.5px;color:#8a9099">every application, by furthest stage reached · Apr 21 – Jul 14</div>
        <div style="margin-left:auto;display:flex;align-items:baseline;gap:7px"><span style="font-size:24px;font-weight:700;letter-spacing:-0.02em;color:#16a34a">6%</span><span style="font-size:12.5px;font-weight:600;color:#6f7680">overall · applied → offer</span></div>
      </div>

      <div style="display:flex;gap:14px">
        <div style="position:relative;width:20px;height:212px;flex:none">
          <span style="position:absolute;right:0;top:24px;font-size:10.5px;color:#b8bdc4;transform:translateY(-50%)">16</span>
          <span style="position:absolute;right:0;top:114px;font-size:10.5px;color:#b8bdc4;transform:translateY(-50%)">8</span>
          <span style="position:absolute;right:0;top:159px;font-size:10.5px;color:#b8bdc4;transform:translateY(-50%)">4</span>
          <span style="position:absolute;right:0;top:204px;font-size:10.5px;color:#b8bdc4;transform:translateY(-50%)">0</span>
        </div>
        <div style="flex:1;position:relative">
          <div style="position:absolute;left:0;right:0;top:24px;height:1px;background:#f0f0ed"></div>
          <div style="position:absolute;left:0;right:0;top:114px;height:1px;background:#f0f0ed"></div>
          <div style="position:absolute;left:0;right:0;top:159px;height:1px;background:#f0f0ed"></div>
          <div style="display:grid;grid-template-columns:1fr 56px 1fr 56px 1fr 56px 1fr;align-items:end;height:204px;position:relative">
            <div style="display:flex;flex-direction:column;justify-content:flex-end">
              <div style="font-size:24px;font-weight:700;letter-spacing:-0.02em;margin-bottom:6px">16 <span style="font-size:12px;font-weight:600;color:#8a9099">applications</span></div>
              <div style="height:180px;background:linear-gradient(180deg,#2463eb,#4d7bee);border-radius:9px 9px 0 0"></div>
            </div>
            <div style="position:relative;height:180px">
              <div style="position:absolute;inset:0;background:#eef4ff;clip-path:polygon(0 0, 100% 75%, 100% 100%, 0 100%)"></div>
              <span style="position:absolute;left:50%;top:96px;transform:translateX(-50%);font-size:11px;font-weight:700;color:#b3372a;background:#fff;border:1px solid #f2d4cf;border-radius:5px;padding:2px 6px;white-space:nowrap;box-shadow:0 1px 2px rgba(22,24,28,.06)">↓ 75%</span>
            </div>
            <div style="display:flex;flex-direction:column;justify-content:flex-end">
              <div style="font-size:24px;font-weight:700;letter-spacing:-0.02em;margin-bottom:6px">4 <span style="font-size:12px;font-weight:600;color:#8a9099">25% of all</span></div>
              <div style="height:45px;background:linear-gradient(180deg,#2463eb,#4d7bee);border-radius:9px 9px 0 0"></div>
            </div>
            <div style="position:relative;height:45px">
              <div style="position:absolute;inset:0;background:#eef4ff;clip-path:polygon(0 0, 100% 49%, 100% 100%, 0 100%)"></div>
              <span style="position:absolute;left:50%;top:-14px;transform:translateX(-50%);font-size:11px;font-weight:700;color:#b3372a;background:#fff;border:1px solid #f2d4cf;border-radius:5px;padding:2px 6px;white-space:nowrap;box-shadow:0 1px 2px rgba(22,24,28,.06)">↓ 50%</span>
            </div>
            <div style="display:flex;flex-direction:column;justify-content:flex-end">
              <div style="font-size:24px;font-weight:700;letter-spacing:-0.02em;margin-bottom:6px">2 <span style="font-size:12px;font-weight:600;color:#8a9099">13% of all</span></div>
              <div style="height:23px;background:linear-gradient(180deg,#2463eb,#4d7bee);border-radius:9px 9px 0 0"></div>
            </div>
            <div style="position:relative;height:23px">
              <div style="position:absolute;inset:0;background:#eef4ff;clip-path:polygon(0 0, 100% 48%, 100% 100%, 0 100%)"></div>
              <span style="position:absolute;left:50%;top:-14px;transform:translateX(-50%);font-size:11px;font-weight:700;color:#1d7a4f;background:#fff;border:1px solid #cfe5d2;border-radius:5px;padding:2px 6px;white-space:nowrap;box-shadow:0 1px 2px rgba(22,24,28,.06)">↓ 50%</span>
            </div>
            <div style="display:flex;flex-direction:column;justify-content:flex-end">
              <div style="font-size:24px;font-weight:700;letter-spacing:-0.02em;margin-bottom:6px">1 <span style="font-size:12px;font-weight:600;color:#8a9099">6% of all</span></div>
              <div style="height:12px;background:linear-gradient(180deg,#2463eb,#4d7bee);border-radius:6px 6px 0 0"></div>
            </div>
          </div>
          <div style="display:grid;grid-template-columns:1fr 56px 1fr 56px 1fr 56px 1fr;border-top:2px solid #e2e2de;padding-top:10px">
            <div>
              <div style="font-size:14.5px;font-weight:700;color:#16181c;margin-bottom:8px">Applied</div>
              <div style="display:flex;align-items:center"><span style={chip('fiverr.com')}></span><span style="{chip('gong.io')};margin-left:-6px"></span><span style="{chip('hibob.com')};margin-left:-6px"></span><span style="font-size:11px;font-weight:600;color:#8a9099;margin-left:5px">+13</span></div>
            </div>
            <div></div>
            <div>
              <div style="font-size:14.5px;font-weight:700;color:#16181c;margin-bottom:8px">Screen</div>
              <div style="display:flex;align-items:center"><span style={chip('wiz.io')}></span><span style="{chip('meliopayments.com')};margin-left:-6px"></span><span style="{chip('wix.com')};margin-left:-6px"></span><span style="font-size:11px;font-weight:600;color:#8a9099;margin-left:5px">+1</span></div>
            </div>
            <div></div>
            <div>
              <div style="font-size:14.5px;font-weight:700;color:#16181c;margin-bottom:8px">Interview</div>
              <div style="display:flex;align-items:center"><span style={chip('kayma.com')}></span><span style="{chip('linear.app')};margin-left:-6px"></span></div>
            </div>
            <div></div>
            <div>
              <div style="font-size:14.5px;font-weight:700;color:#16181c;margin-bottom:8px">Offer</div>
              <div style="display:flex;align-items:center"><span style={chip('linear.app')}></span></div>
            </div>
          </div>
        </div>
      </div>
      <div style="display:flex;align-items:center;gap:16px;margin-top:16px;font-size:12.5px;color:#6f7680">
        <span>The applied → screen step is where you lose the most — <strong style="color:#16181c">75% never reply</strong>.</span>
        <a href="/preview/redesign/applications" style="margin-left:auto;flex:none">See the 12 quiet ones →</a>
      </div>
    </div>

    <!-- exits strip -->
    <div style="display:flex;align-items:center;gap:20px;border:1px solid #e8e8e5;background:#fff;border-radius:12px;padding:12px 20px;margin-bottom:40px;font-size:13px;color:#4b5158">
      <span><strong style="color:#b3372a">1</strong> rejected</span>
      <span><strong>1</strong> position closed</span>
      <span><strong>1</strong> withdrawn</span>
      <span style="color:#8a9099;margin-left:auto">exits live in <a href="/preview/redesign/applications">Applications → No longer in play</a></span>
    </div>

    <!-- sources + noticing history -->
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:40px;align-items:start">
      <div>
        <div style="font-size:11px;font-weight:600;letter-spacing:.12em;text-transform:uppercase;color:#8a9099;margin-bottom:14px">Where they come from</div>
        <div style="display:flex;flex-direction:column;gap:12px;font-size:13px">
          {#each SOURCES as s (s.name)}
            <div>
              <div style="display:flex;align-items:baseline;margin-bottom:5px"><span style="font-weight:600">{s.name}</span><span style="margin-left:auto;color:#8a9099">{s.meta}</span></div>
              <div style="height:6px;background:#f0f0ed;border-radius:3px"><div style="width:{s.w}%;height:6px;background:{s.color};border-radius:3px"></div></div>
            </div>
          {/each}
        </div>
        <div style="font-size:12.5px;color:#6f7680;margin-top:14px;line-height:1.55">Referrals are <strong style="color:#16181c">your only source that converts past screen</strong> — half your volume goes to LinkedIn, which has produced one screen in 12 weeks.</div>
      </div>
      <div>
        <div style="display:flex;align-items:baseline;margin-bottom:14px">
          <span style="font-size:11px;font-weight:600;letter-spacing:.12em;text-transform:uppercase;color:#e0641f">✦ What we’ve noticed</span>
          <span style="font-size:12px;color:#8a9099;margin-left:10px">one per day, shown on Home — full history here</span>
        </div>
        <div style="display:flex;flex-direction:column;gap:12px;font-size:13px;color:#4b5158">
          {#each NOTICED as n (n.d)}
            <div style="display:flex;gap:12px"><span style="color:#b8bdc4;width:42px;flex:none;font-size:12px;padding-top:1px">{n.d}</span><span style="line-height:1.55"><strong>{n.lead}</strong> {n.rest}</span></div>
          {/each}
        </div>
      </div>
    </div>

  </div>
</div>

<style>
  .page { min-height: 100vh; background: #f6f6f3; color: #16181c; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; -webkit-font-smoothing: antialiased; }
  .page :global(a) { color: #2463eb; text-decoration: none; }
</style>
