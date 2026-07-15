<script>
  import PreviewBanner from '../PreviewBanner.svelte';
  import { logo } from '../fixtures.js';
</script>

<svelte:head><title>Redesign preview — States</title></svelte:head>

<PreviewBanner tag="3.3" title="States — empty home · at credit limit" />

<div class="page">

  <!-- STATE 1: EMPTY HOME -->
  <div style="max-width:1100px;margin:0 auto;padding:36px 32px 12px">
    <div style="font-size:11px;font-weight:600;letter-spacing:.12em;text-transform:uppercase;color:#8a9099;margin-bottom:12px">State 1 · First run — nothing tracked yet</div>
  </div>
  <div style="max-width:1100px;margin:0 auto 56px;border:1px solid #e2e2de;border-radius:16px;overflow:hidden;background:#f6f6f3">

    <div style="background:#fff;border-bottom:1px solid #e8e8e5">
      <div style="display:flex;align-items:center;gap:24px;padding:12px 28px">
        <div style="display:flex;align-items:center;gap:9px;flex:none">
          <div style="width:24px;height:24px;border-radius:50%;border:2.5px solid #16181c;display:flex;align-items:center;justify-content:center"><div style="width:9px;height:9px;border-radius:50%;background:#2463eb"></div></div>
          <div style="font-size:16px;font-weight:700;letter-spacing:-0.01em">Pursuit</div>
        </div>
        <div style="display:flex;align-items:center;gap:4px">
          <div style="background:#16181c;color:#fff;border-radius:8px;padding:7px 14px;font-size:13.5px;font-weight:600">Home</div>
          <div style="border-radius:8px;padding:7px 14px;font-size:13.5px;color:#b8bdc4">Applications</div>
          <div style="border-radius:8px;padding:7px 14px;font-size:13.5px;color:#b8bdc4">Insights</div>
        </div>
        <div style="display:flex;align-items:center;gap:12px;margin-left:auto">
          <div title="Prep credits" style="border:1px solid #e8e8e5;border-radius:8px;padding:7px 11px;font-size:12px;color:#6f7680;flex:none">✦ <strong style="color:#16181c">10</strong><span style="color:#8a9099">/10</span></div>
          <div style="width:30px;height:30px;border-radius:50%;background:#e0641f;color:#fff;display:flex;align-items:center;justify-content:center;font-size:12px;font-weight:700;flex:none">Y</div>
        </div>
      </div>
    </div>

    <div style="max-width:720px;margin:0 auto;padding:64px 32px 72px;text-align:center">
      <div style="font-size:13px;color:#8a9099;margin-bottom:6px">Tuesday · July 14, 2026</div>
      <h1 style="font-size:30px;font-weight:700;letter-spacing:-0.02em;margin:0 0 10px">Good morning, Yonatan.</h1>
      <p style="font-size:15px;line-height:1.6;color:#4b5158;margin:0 auto 28px;max-width:440px">Track your first application and Pursuit starts working — reply nudges, interview prep, the whole picture in one place.</p>

      <div style="display:inline-flex;flex-direction:column;gap:10px;width:460px;max-width:100%">
        <div class="drop" style="border:1.5px dashed #c8ccd2;background:#fff;border-radius:14px;padding:28px 24px;cursor:pointer">
          <div style="font-size:14.5px;font-weight:600;margin-bottom:4px">Paste a job URL, drop a screenshot, or describe the role</div>
          <div style="font-size:12.5px;color:#8a9099">We fill in the company, role and source for you · ⌘V works right here</div>
        </div>
        <div style="display:flex;align-items:center;gap:10px;justify-content:center">
          <div style="background:#2463eb;color:#fff;border-radius:9px;padding:10px 20px;font-size:13.5px;font-weight:600;cursor:pointer">New application <span style="opacity:.65;font-weight:400">⌘N</span></div>
          <a href="#top" style="font-size:13px">or import from a spreadsheet</a>
        </div>
      </div>

      <div style="display:flex;align-items:center;justify-content:center;gap:22px;margin-top:44px;font-size:12.5px;color:#b8bdc4">
        <span>Then this page becomes:</span>
        <span style="display:flex;align-items:center;gap:6px"><span style="width:7px;height:7px;border-radius:50%;background:#e0641f"></span>what to prep</span>
        <span style="display:flex;align-items:center;gap:6px"><span style="width:7px;height:7px;border-radius:50%;background:#16a34a"></span>what to decide</span>
        <span style="display:flex;align-items:center;gap:6px"><span style="width:7px;height:7px;border-radius:50%;background:#b3372a"></span>who to nudge</span>
      </div>
    </div>
  </div>

  <!-- STATE 2: AT CREDIT LIMIT -->
  <div style="max-width:1100px;margin:0 auto;padding:0 32px 12px">
    <div style="font-size:11px;font-weight:600;letter-spacing:.12em;text-transform:uppercase;color:#8a9099;margin-bottom:12px">State 2 · At the credit limit — generating a new round brief with 0 / 10 left</div>
  </div>
  <div style="max-width:1100px;margin:0 auto 64px;border:1px solid #e2e2de;border-radius:16px;overflow:hidden;background:#f6f6f3;padding:36px 40px">

    <div style="max-width:700px;margin:0 auto">
      <div style="background:#fff;border:1px solid #e8e8e5;border-radius:16px;overflow:hidden;box-shadow:0 1px 3px rgba(22,24,28,.04)">
        <div style="display:flex;align-items:center;gap:12px;background:#fbfbf9;border-bottom:1px solid #eeeeea;padding:11px 20px">
          <div title="Wiz" style="width:24px;height:24px;border-radius:6px;border:1px solid #eeeeea;background:#fff url({logo('wiz.io')}) center/contain no-repeat;flex:none"></div>
          <span style="font-size:12.5px;color:#4b5158">Researched for <strong>Wiz Inc.</strong> · <a href="#top">wiz.io</a></span>
          <span title="Prep credits" style="margin-left:auto;border:1px solid #f2d4cf;background:#fdf1ef;border-radius:8px;padding:5px 10px;font-size:12px;color:#b3372a;flex:none">✦ <strong>0</strong>/10 credits</span>
        </div>
        <div style="padding:18px 22px 0">
          <div style="display:flex;align-items:baseline;gap:10px;margin-bottom:14px">
            <div style="font-size:19px;font-weight:700;letter-spacing:-0.01em">✦ Interview playbook</div>
          </div>
          <div style="display:flex;align-items:center;gap:7px;padding-bottom:14px;border-bottom:1px solid #f0f0ed">
            <div style="background:#eef4ff;border:1px solid #cdddfb;border-radius:9px;padding:8px 15px;font-size:13px;font-weight:700;color:#2463eb">▦ Company</div>
            <div style="border:1px solid #e8e8e5;border-radius:9px;padding:8px 15px;font-size:13px;color:#4b5158">Recruiter screen</div>
            <div style="border:1px dashed #e2e2de;border-radius:9px;padding:8px 15px;font-size:13px;color:#b8bdc4">+ Add round</div>
          </div>
        </div>
        <div style="padding:28px 22px 30px;text-align:center">
          <div style="font-size:15px;font-weight:700;margin-bottom:6px">You’ve used all 10 prep credits.</div>
          <div style="font-size:13.5px;line-height:1.6;color:#4b5158;max-width:420px;margin:0 auto 18px">Everything already generated stays yours — every brief, every citation. New round briefs unlock when your credits renew on <strong>August 1</strong>.</div>
          <div style="display:flex;align-items:center;justify-content:center;gap:10px">
            <div style="background:#2463eb;color:#fff;border-radius:9px;padding:9px 18px;font-size:13px;font-weight:600;cursor:pointer">Get more credits</div>
            <a href="#top" style="font-size:13px">See what’s included</a>
          </div>
          <div style="font-size:12px;color:#8a9099;margin-top:16px">Tight on time before an interview? <a href="#top">Write to us</a> — we won’t leave you hanging the night before.</div>
        </div>
      </div>
      <div style="font-size:12px;color:#8a9099;margin-top:10px">Tracker features — statuses, follow-ups, archive, the home page — are never limited by credits.</div>
    </div>
  </div>

</div>

<style>
  .page { min-height: 100vh; background: #f6f6f3; color: #16181c; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; -webkit-font-smoothing: antialiased; }
  .page :global(a) { color: #2463eb; text-decoration: none; }
  .drop:hover { border-color: #2463eb !important; }
</style>
