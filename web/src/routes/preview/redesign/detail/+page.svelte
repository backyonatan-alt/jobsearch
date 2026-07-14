<script>
  import PreviewBanner from '../PreviewBanner.svelte';
  import TopNav from '../TopNav.svelte';
  import { logo, logoLg } from '../fixtures.js';

  let tab = $state('panel');
  let identity = $state('bar'); // bar | confirm | fix | verified

  const TABS = [
    { id: 'company', label: '▦ Company', done: false },
    { id: 'recruiter', label: 'Recruiter screen', done: true },
    { id: 'hm', label: 'HM interview', done: true },
    { id: 'panel', label: 'Panel · tomorrow', done: false }
  ];
  const tabStyle = (t) => {
    const active = tab === t.id;
    const bg = active ? (t.id === 'panel' ? '#fff7f1' : '#eef4ff') : '#fff';
    const border = active ? (t.id === 'panel' ? '#f0d9c4' : '#cdddfb') : '#e8e8e5';
    const color = active ? (t.id === 'panel' ? '#c05310' : '#2463eb') : (t.done ? '#1d7a4f' : '#4b5158');
    return `background:${bg};border:1px solid ${border};color:${color};font-weight:${active ? 700 : 500}`;
  };
  const PANELISTS = [
    { init: 'RB', bg: '#eef4ff', color: '#2463eb', name: 'Rotem Bar', role: 'VP Engineering' },
    { init: 'NG', bg: '#eef7ef', color: '#1d7a4f', name: 'Nadav Gil', role: 'Head of Design' },
    { init: 'TA', bg: '#fdf3ec', color: '#c05310', name: 'Tal Amir', role: 'Group PM · peer' }
  ];
</script>

<svelte:head><title>Redesign preview — Application detail</title></svelte:head>

<PreviewBanner tag="3.2 v4" title="Application detail — live tabs + “Not them?” flow" note="Tabs switch content. “Not them?” opens the wrong-company check." />

<div class="page">
  <TopNav active="applications" primaryCta={false} />

  <div style="max-width:1160px;margin:0 auto;padding:26px 32px 80px">

    <div style="font-size:13px;color:#8a9099;margin-bottom:20px"><a href="/preview/redesign/applications" style="color:#8a9099">Applications</a> <span style="margin:0 4px">/</span> <span style="color:#16181c;font-weight:600">Kayma</span></div>

    <div style="display:flex;align-items:flex-start;gap:18px;margin-bottom:22px">
      <div title="Kayma" style="width:52px;height:52px;border-radius:13px;border:1px solid #e8e8e5;background:#fff url({logoLg('kayma.com')}) center/contain no-repeat;flex:none"></div>
      <div style="flex:1;min-width:0">
        <div style="display:flex;align-items:baseline;gap:12px">
          <h1 style="font-size:28px;font-weight:700;letter-spacing:-0.02em;margin:0">Kayma</h1>
          <span style="display:flex;align-items:center;gap:6px;font-size:13px;font-weight:600;color:#c05310"><span style="width:8px;height:8px;border-radius:50%;background:#e0641f"></span>Interview</span>
        </div>
        <div style="font-size:14.5px;color:#4b5158;margin:3px 0 6px">VP Product</div>
        <div style="font-size:12.5px;color:#8a9099">Applied June 22 · Referral (Noa Adler) · Tel Aviv · <a href="#top">Open job post ↗</a></div>
      </div>
      <div style="display:flex;align-items:center;gap:14px;flex:none;font-size:13px;padding-top:6px">
        <a href="#top" style="color:#4b5158">Update status ▾</a>
        <a href="#top" style="color:#8a9099">⋯</a>
      </div>
    </div>

    <!-- tomorrow banner -->
    <div style="display:flex;align-items:center;gap:16px;background:#fff7f1;border:1px solid #f0d9c4;border-radius:14px;padding:16px 22px;margin-bottom:28px">
      <div style="flex:none;width:44px;border:1px solid #f0d9c4;border-radius:9px;overflow:hidden;text-align:center;background:#fff"><div style="background:#e0641f;color:#fff;font-size:9px;font-weight:700;letter-spacing:.08em;padding:2px 0">JUL</div><div style="font-size:17px;font-weight:700;color:#c05310;padding:2px 0 3px">15</div></div>
      <div style="flex:1;min-width:0">
        <div style="font-size:15px;font-weight:700">Panel tomorrow at 10:00 — 3 interviewers, 90 minutes.</div>
        <div style="font-size:13px;color:#6f7680">Your briefs are ready below. 20 focused minutes tonight is enough.</div>
      </div>
      <a href="/preview/redesign/brief" style="background:#2463eb;color:#fff;border-radius:9px;padding:10px 18px;font-size:13.5px;font-weight:600;cursor:pointer;flex:none;display:block">Review panel prep →</a>
    </div>

    <div style="display:grid;grid-template-columns:1.9fr 1fr;gap:40px;align-items:start">

      <!-- playbook -->
      <div>
        <div style="background:#fff;border:1px solid #e8e8e5;border-radius:16px;overflow:hidden;box-shadow:0 1px 3px rgba(22,24,28,.04)">

          {#if identity === 'bar'}
            <div style="display:flex;align-items:center;gap:12px;background:#fbfbf9;border-bottom:1px solid #eeeeea;padding:11px 20px">
              <div title="Kayma" style="width:24px;height:24px;border-radius:6px;border:1px solid #eeeeea;background:#fff url({logo('kayma.com')}) center/contain no-repeat;flex:none"></div>
              <span style="font-size:12.5px;color:#4b5158;min-width:0;white-space:nowrap;overflow:hidden;text-overflow:ellipsis">Researched for <strong>Kayma Ltd</strong> · <a href="#top">kayma.com</a> <span style="color:#8a9099">— Tel Aviv data & AI product studio, ~120 people</span></span>
              <button class="linkbtn" style="margin-left:auto;flex:none;font-size:12px;color:#8a9099" onclick={() => (identity = 'confirm')}>Not them?</button>
            </div>
          {:else if identity === 'confirm'}
            <div style="background:#fdf6ef;border-bottom:1px solid #f0d9c4;padding:18px 22px">
              <div style="font-size:14px;font-weight:700;margin-bottom:12px">Is this the company you’re interviewing with?</div>
              <div style="display:flex;align-items:center;gap:14px;background:#fff;border:1px solid #eeeeea;border-radius:12px;padding:14px 18px;margin-bottom:12px">
                <div title="Kayma" style="width:38px;height:38px;border-radius:10px;border:1px solid #eeeeea;background:#fff url({logoLg('kayma.com')}) center/contain no-repeat;flex:none"></div>
                <div style="flex:1;min-width:0;font-size:13px;line-height:1.5"><strong>Kayma Ltd</strong> · <a href="#top">kayma.com</a><br><span style="color:#6f7680">Tel Aviv data & AI product studio, ~120 people, founded 2018. Clients incl. government & enterprise.</span></div>
              </div>
              <div style="display:flex;align-items:center;gap:10px">
                <button style="background:#2463eb;color:#fff;border:0;border-radius:8px;padding:8px 16px;font-size:13px;font-weight:600;cursor:pointer;font-family:inherit;white-space:nowrap" onclick={() => (identity = 'verified')}>Yes — that’s them</button>
                <button style="border:1px solid #e2d4c4;background:#fff;color:#4b5158;border-radius:8px;padding:8px 16px;font-size:13px;font-weight:600;cursor:pointer;font-family:inherit;white-space:nowrap" onclick={() => (identity = 'fix')}>No — wrong company</button>
                <span style="font-size:12px;color:#8a9099">Same-name companies are the #1 cause of a wasted brief.</span>
              </div>
            </div>
          {:else if identity === 'fix'}
            <div style="background:#fdf6ef;border-bottom:1px solid #f0d9c4;padding:18px 22px">
              <div style="font-size:14px;font-weight:700;margin-bottom:4px">Point us at the right one.</div>
              <div style="font-size:12.5px;color:#6f7680;margin-bottom:12px">Paste their website or LinkedIn page — we’ll re-research and rebuild the briefs. Your credit for this playbook is returned.</div>
              <div style="display:flex;align-items:center;gap:10px">
                <input placeholder="https://" style="flex:1;border:1px solid #e2d4c4;border-radius:8px;padding:9px 13px;font-size:13px;background:#fff;color:#16181c;outline:none;font-family:inherit">
                <button style="background:#2463eb;color:#fff;border:0;border-radius:8px;padding:9px 16px;font-size:13px;font-weight:600;cursor:pointer;flex:none;font-family:inherit" onclick={() => (identity = 'verified')}>Re-research →</button>
                <button class="linkbtn" style="font-size:12.5px;color:#8a9099;flex:none" onclick={() => (identity = 'bar')}>Cancel</button>
              </div>
            </div>
          {:else}
            <div style="display:flex;align-items:center;gap:10px;background:#f3faf4;border-bottom:1px solid #cfe5d2;padding:10px 20px">
              <span style="color:#1d7a4f;font-size:13px;font-weight:600">✓ Verified — Kayma Ltd · kayma.com</span>
              <button class="linkbtn" style="margin-left:auto;font-size:12px;color:#8a9099" onclick={() => (identity = 'confirm')}>change</button>
            </div>
          {/if}

          <div style="padding:18px 22px 0">
            <div style="display:flex;align-items:center;gap:10px;margin-bottom:14px">
              <div style="font-size:19px;font-weight:700;letter-spacing:-0.01em">✦ Interview playbook</div>
              <div style="font-size:12px;color:#8a9099">refreshed 19h ago</div>
            </div>
            <div style="display:flex;align-items:center;gap:7px;padding-bottom:14px;border-bottom:1px solid #f0f0ed;flex-wrap:wrap">
              {#each TABS as t (t.id)}
                <button class="tab" style="display:flex;align-items:center;gap:6px;border-radius:9px;padding:8px 15px;font-size:13px;cursor:pointer;font-family:inherit;{tabStyle(t)}" onclick={() => (tab = t.id)}>{t.done && tab !== t.id ? '✓ ' + t.label : t.label}</button>
              {/each}
              <div style="border:1px dashed #e2e2de;border-radius:9px;padding:8px 15px;font-size:13px;color:#b8bdc4;cursor:pointer">+ Add round</div>
            </div>
          </div>

          <div style="padding:20px 22px 24px">
            {#if tab === 'company'}
              <div style="font-size:15px;font-weight:700;margin-bottom:12px">Company brief <span style="font-size:12px;font-weight:400;color:#8a9099">· shared across every round</span></div>
              <div style="background:#eef4ff;border:1px solid #cdddfb;border-radius:12px;padding:16px 20px;margin-bottom:14px">
                <div style="font-size:13px;font-weight:700;color:#2463eb;margin-bottom:8px">✦ What this team grades for</div>
                <div style="display:flex;flex-direction:column;gap:7px;font-size:13.5px;line-height:1.55;color:#1e3a6e">
                  <div>· Data-first product judgment — every case comes back to “how would you measure it”.</div>
                  <div>· Comfort selling to enterprise & government buyers, their core market.</div>
                  <div>· Founders still interview finalists; expect a “why us, why now” conversation.</div>
                </div>
              </div>
              <div style="display:flex;align-items:center;gap:8px;flex-wrap:wrap;font-size:12px;color:#8a9099"><span style="font-weight:600;color:#6f7680">Sources</span><a href="#top" style="border:1px solid #e8e8e5;border-radius:14px;padding:3px 10px">kayma.com/about</a><a href="#top" style="border:1px solid #e8e8e5;border-radius:14px;padding:3px 10px">Calcalist profile · 2026</a><a href="#top" style="border:1px solid #e8e8e5;border-radius:14px;padding:3px 10px">+ 5 more</a></div>
            {:else if tab === 'recruiter'}
              <div style="display:flex;align-items:center;gap:10px;margin-bottom:10px"><span style="color:#1d7a4f;font-weight:700;font-size:15px">✓ Recruiter screen — happened Jun 30</span><a href="#top" style="margin-left:auto;font-size:12.5px">View archived brief →</a></div>
              <div style="font-size:13.5px;line-height:1.6;color:#4b5158">Your debrief: comp range confirmed (asked for ₪75–85k), process is 4 rounds, they’re moving fast because the incumbent left. Maya flagged the panel as “the real bar”.</div>
            {:else if tab === 'hm'}
              <div style="display:flex;align-items:center;gap:10px;margin-bottom:10px"><span style="color:#1d7a4f;font-weight:700;font-size:15px">✓ HM interview — happened Jul 8</span><a href="#top" style="margin-left:auto;font-size:12.5px">View archived brief →</a></div>
              <div style="font-size:13.5px;line-height:1.6;color:#4b5158">Your debrief: Dan cares most about platform-as-product thinking; he pushed on stakeholder conflict twice. He mentioned Rotem will “want the technical story” — that’s tomorrow’s panel.</div>
            {:else}
              <div style="display:flex;align-items:center;gap:10px;margin-bottom:16px">
                <div style="font-size:15px;font-weight:700">Panel — 3 interviewers · 90 min</div>
              </div>
              <div style="display:grid;grid-template-columns:1fr 1fr 1fr;gap:10px;margin-bottom:16px">
                {#each PANELISTS as p (p.init)}
                  <a href="/preview/redesign/brief" class="pcard" style="border:1px solid #eeeeea;border-radius:12px;padding:14px;text-align:center;color:#16181c;display:block">
                    <div style="width:38px;height:38px;border-radius:50%;background:{p.bg};color:{p.color};display:flex;align-items:center;justify-content:center;font-size:13px;font-weight:700;margin:0 auto 8px">{p.init}</div>
                    <div style="font-size:13.5px;font-weight:700">{p.name}</div>
                    <div style="font-size:11.5px;color:#8a9099">{p.role}</div>
                  </a>
                {/each}
              </div>
              <a href="/preview/redesign/brief" class="doorbar" style="display:flex;align-items:center;gap:14px;background:#2463eb;border-radius:12px;padding:16px 20px;cursor:pointer;color:#fff">
                <div style="flex:1">
                  <div style="font-size:14.5px;font-weight:700">Open the panel brief →</div>
                  <div style="font-size:12px;opacity:.75">Who grades what, questions to expect, your angle per interviewer, sources · 5 min read</div>
                </div>
                <span style="font-size:20px">→</span>
              </a>
            {/if}
          </div>
        </div>
        <div style="display:flex;align-items:center;gap:8px;margin-top:10px;font-size:12px;color:#8a9099">
          <a href="#top">Refresh company brief</a><span style="color:#d8dade">·</span><span>1 credit per round brief · <strong style="color:#16181c">7 left</strong></span>
        </div>
      </div>

      <!-- right column -->
      <div style="font-size:13.5px">
        <div style="margin-bottom:24px">
          <div style="display:flex;align-items:baseline;margin-bottom:8px">
            <span style="font-size:11px;font-weight:600;letter-spacing:.12em;text-transform:uppercase;color:#8a9099">The role, in short</span>
            <a href="#top" style="font-size:12px;margin-left:auto">Full JD →</a>
          </div>
          <div style="font-size:13px;line-height:1.65;color:#4b5158">Own the data & AI platform group — 3 squads, ~14 people, reports to the CPO. Wants B2B SaaS scaled past $50M ARR. Hybrid TLV, 3 days on-site.</div>
        </div>

        <div style="border-top:1px solid #e2e2de;padding-top:18px;margin-bottom:24px">
          <div style="font-size:11px;font-weight:600;letter-spacing:.12em;text-transform:uppercase;color:#8a9099;margin-bottom:12px">Process</div>
          <div style="display:flex;flex-direction:column;gap:9px;font-size:13px">
            <div style="display:flex;align-items:center;gap:10px;color:#1d7a4f"><span style="width:20px;height:20px;border-radius:50%;background:#eef7ef;display:flex;align-items:center;justify-content:center;font-size:11px;flex:none">✓</span>Recruiter screen <span style="color:#b8bdc4;margin-left:auto">Jun 30</span></div>
            <div style="display:flex;align-items:center;gap:10px;color:#1d7a4f"><span style="width:20px;height:20px;border-radius:50%;background:#eef7ef;display:flex;align-items:center;justify-content:center;font-size:11px;flex:none">✓</span>HM interview <span style="color:#b8bdc4;margin-left:auto">Jul 8</span></div>
            <div style="display:flex;align-items:center;gap:10px;color:#c05310;font-weight:700"><span style="width:20px;height:20px;border-radius:50%;background:#fff7f1;border:1.5px solid #e0641f;display:flex;align-items:center;justify-content:center;flex:none"><span style="width:7px;height:7px;border-radius:50%;background:#e0641f"></span></span>Panel <span style="margin-left:auto">tomorrow</span></div>
            <div style="display:flex;align-items:center;gap:10px;color:#9aa1ab"><span style="width:20px;height:20px;border-radius:50%;border:1.5px dashed #d8dade;flex:none"></span>Final / CEO</div>
          </div>
          <div style="font-size:11.5px;color:#b8bdc4;margin-top:10px">as you mapped it from the recruiter call · <a href="#top" style="color:#8a9099">edit</a></div>
        </div>

        <div style="border-top:1px solid #e2e2de;padding-top:18px;margin-bottom:24px">
          <div style="font-size:11px;font-weight:600;letter-spacing:.12em;text-transform:uppercase;color:#8a9099;margin-bottom:12px">People</div>
          <div style="display:flex;flex-direction:column;gap:10px;font-size:13px">
            <div style="display:flex;align-items:center;gap:10px"><span style="width:26px;height:26px;border-radius:50%;background:#eef4ff;color:#2463eb;display:flex;align-items:center;justify-content:center;font-size:11px;font-weight:700;flex:none">MS</span><span style="flex:1;min-width:0"><strong>Maya Sharon</strong> <span style="color:#8a9099">· recruiter</span></span><a href="#top" style="font-size:12px;flex:none">Email →</a></div>
            <div style="display:flex;align-items:center;gap:10px"><span style="width:26px;height:26px;border-radius:50%;background:#fdf3ec;color:#c05310;display:flex;align-items:center;justify-content:center;font-size:11px;font-weight:700;flex:none">DL</span><span style="flex:1;min-width:0"><strong>Dan Lev</strong> <span style="color:#8a9099">· CPO, hiring manager</span></span><a href="#top" style="font-size:12px;flex:none">LinkedIn →</a></div>
            <a href="#top" style="font-size:12px;color:#8a9099">+ Add a contact</a>
          </div>
        </div>

        <div style="border-top:1px solid #e2e2de;padding-top:18px">
          <div style="display:flex;align-items:baseline;margin-bottom:12px">
            <span style="font-size:11px;font-weight:600;letter-spacing:.12em;text-transform:uppercase;color:#8a9099">Activity</span>
            <span style="margin-left:auto;display:flex;gap:12px;font-size:12px"><a href="#top">+ Follow-up</a><a href="#top">+ Note</a></span>
          </div>
          <div style="display:flex;flex-direction:column;gap:9px;font-size:12.5px;color:#4b5158">
            <div style="display:flex;gap:10px"><span style="color:#b8bdc4;width:40px;flex:none">Jul 13</span><span>Panel confirmed — 3 interviewers</span></div>
            <div style="display:flex;gap:10px"><span style="color:#b8bdc4;width:40px;flex:none">Jul 8</span><span>HM interview — debrief logged</span></div>
            <div style="display:flex;gap:10px"><span style="color:#b8bdc4;width:40px;flex:none">Jun 30</span><span>Recruiter screen — debrief logged</span></div>
            <div style="display:flex;gap:10px"><span style="color:#b8bdc4;width:40px;flex:none">Jun 22</span><span>Applied via referral (Noa Adler)</span></div>
          </div>
        </div>
      </div>

    </div>
  </div>
</div>

<style>
  .page { min-height: 100vh; background: #f6f6f3; color: #16181c; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; -webkit-font-smoothing: antialiased; }
  .page :global(a) { color: #2463eb; text-decoration: none; }
  .linkbtn { background: none; border: 0; cursor: pointer; padding: 0; font-family: inherit; }
  .tab { border: 1px solid #e8e8e5; background: #fff; }
  .pcard { color: #16181c !important; }
  .pcard:hover { border-color: #b9c6e8 !important; }
  .doorbar { color: #fff !important; }
  .doorbar:hover { background: #1a4fc4 !important; }
</style>
