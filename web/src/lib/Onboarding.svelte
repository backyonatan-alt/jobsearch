<script>
  // Onboarding — Variant A (centered card, dot progress, illustrative).
  // Five sequential stages sell the magic before the user has to pick a
  // path. The final stage has three choice tiles (paste / demo / explore)
  // that wire into the same backend flows the rest of the app uses.
  import { api } from '$lib/api.js';

  let { onDone } = $props();

  let step = $state(0);
  const total = 5;

  // Stage-4 sub-state — when the user picks "Paste your applications" we
  // flip to a textarea+parse view, then to a results view when finished.
  let pickerMode = $state('pick'); // 'pick' | 'paste' | 'parsing' | 'done'
  let text = $state('');
  let progress = $state([]); // [{idx, status: 'pending'|'parsing'|'creating'|'done'|'error', label, error?}]
  let globalError = $state('');
  let seeding = $state(false);
  const added = $derived(progress.filter(p => p.status === 'done').length);

  function next() { step = Math.min(total - 1, step + 1); }
  function prev() {
    if (step === total - 1 && pickerMode !== 'pick') {
      // Back from a sub-state returns to the picker.
      pickerMode = 'pick';
      progress = [];
      globalError = '';
      return;
    }
    step = Math.max(0, step - 1);
  }

  function splitEntries(s) {
    const parts = s.split(/\n{2,}/).map(p => p.trim()).filter(p => p.length > 4);
    if (parts.length > 1) return parts;
    const lines = s.split(/\n/).map(p => p.trim()).filter(p => p.length > 4);
    if (lines.length > 1 && lines.every(l => /^https?:\/\//.test(l))) return lines;
    return parts.length ? parts : (s.trim() ? [s.trim()] : []);
  }

  async function startPaste() {
    const entries = splitEntries(text);
    if (entries.length === 0) {
      // Nothing to parse — just finish.
      await markOnboarded();
      return;
    }
    pickerMode = 'parsing';
    globalError = '';
    progress = entries.map((e, idx) => ({
      idx,
      status: 'pending',
      label: e.slice(0, 60) + (e.length > 60 ? '…' : '')
    }));
    for (let i = 0; i < entries.length; i++) {
      progress[i] = { ...progress[i], status: 'parsing' };
      try {
        const parsed = await api('/api/applications/parse', {
          method: 'POST',
          body: JSON.stringify({ text: entries[i] })
        });
        progress[i] = { ...progress[i], status: 'creating', label: `${parsed.company || '?'} · ${parsed.role || '?'}` };
        const payload = {
          company:     parsed.company,
          role:        parsed.role,
          status:      'applied',
          source:      parsed.source,
          jd_url:      parsed.jd_url,
          location:    parsed.location,
          salary_note: parsed.salary_note
        };
        for (const k of Object.keys(payload)) if (!payload[k]) delete payload[k];
        await api('/api/applications', { method: 'POST', body: JSON.stringify(payload) });
        progress[i] = { ...progress[i], status: 'done' };
      } catch (e) {
        progress[i] = { ...progress[i], status: 'error', error: e.message };
      }
    }
    pickerMode = 'done';
  }

  async function seedDemo() {
    if (seeding) return;
    seeding = true;
    try {
      await api('/api/me/demo-seed', { method: 'POST' });
      await markOnboarded();
    } catch (e) {
      globalError = e.message || 'Could not seed demo data.';
      seeding = false;
    }
  }

  async function markOnboarded() {
    try {
      await api('/api/me/onboarded', { method: 'POST' });
    } catch (e) {
      console.error('mark onboarded', e);
    }
    // Drop a hint so the Today page can show a one-time tooltip on the
    // New-application button after the overlay closes.
    try { sessionStorage.setItem('pursuit_first_open', '1'); } catch {}
    onDone?.();
  }

  function skipForNow() { markOnboarded(); }
</script>

<div class="ob-overlay" role="dialog" aria-modal="true">
  <div class="ob-card">
    <!-- Progress dots — active dot expands -->
    <div class="dots">
      {#each Array(total) as _, i}
        <span class="dot" class:active={i === step} class:done={i < step}></span>
      {/each}
    </div>

    <!-- Illustration well -->
    <div class="art">
      {#if step === 0}
        <svg viewBox="0 0 240 200" width="240" height="200" fill="none" aria-hidden="true">
          <circle cx="120" cy="105" r="78" stroke="var(--accent)" stroke-width="1" opacity="0.16"/>
          <circle cx="120" cy="105" r="58" stroke="var(--accent)" stroke-width="1.2" opacity="0.32"/>
          <circle cx="120" cy="105" r="38" stroke="var(--accent)" stroke-width="1.6" opacity="0.7"/>
          <circle cx="120" cy="105" r="20" fill="var(--accent-tint)" stroke="var(--accent)" stroke-width="1.6"/>
          <circle cx="150" cy="78" r="9" fill="var(--accent)"/>
        </svg>
      {:else if step === 1}
        <svg viewBox="0 0 320 200" width="320" height="200" fill="none" aria-hidden="true">
          <rect x="14" y="34" width="124" height="132" rx="10" fill="var(--surface-2)" stroke="var(--rule)" stroke-width="1"/>
          <rect x="26" y="46" width="68" height="6" rx="3" fill="var(--mute-2)"/>
          <rect x="26" y="60" width="100" height="4" rx="2" fill="var(--rule-strong)"/>
          <rect x="26" y="70" width="92" height="4" rx="2" fill="var(--rule-strong)"/>
          <rect x="26" y="80" width="100" height="4" rx="2" fill="var(--rule-strong)"/>
          <rect x="26" y="90" width="78" height="4" rx="2" fill="var(--rule-strong)"/>
          <rect x="26" y="104" width="88" height="4" rx="2" fill="var(--rule-strong)"/>
          <rect x="26" y="114" width="100" height="4" rx="2" fill="var(--rule-strong)"/>
          <rect x="26" y="124" width="64" height="4" rx="2" fill="var(--rule-strong)"/>
          <rect x="26" y="138" width="92" height="4" rx="2" fill="var(--rule-strong)"/>
          <rect x="26" y="148" width="74" height="4" rx="2" fill="var(--rule-strong)"/>
          <path d="M148 100h22m-6-6 6 6-6 6" stroke="var(--accent)" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M154 84l1.4 2.6L158 88l-2.6 1.4L154 92l-1.4-2.6L150 88l2.6-1.4z" fill="var(--accent)"/>
          <rect x="180" y="60" width="126" height="80" rx="10" fill="white" stroke="var(--accent)" stroke-width="1.2"/>
          <circle cx="196" cy="76" r="8" fill="var(--accent-tint)"/>
          <rect x="210" y="72" width="50" height="5" rx="2.5" fill="var(--ink)"/>
          <rect x="210" y="82" width="76" height="4" rx="2" fill="var(--mute-2)"/>
          <rect x="192" y="100" width="52" height="14" rx="7" fill="var(--accent-tint)"/>
          <text x="218" y="111" font-family="system-ui" font-size="9" fill="var(--accent-text)" font-weight="600">Applied</text>
          <rect x="192" y="120" width="100" height="3" rx="1.5" fill="var(--rule)"/>
          <rect x="192" y="128" width="80" height="3" rx="1.5" fill="var(--rule)"/>
        </svg>
      {:else if step === 2}
        <svg viewBox="0 0 320 200" width="320" height="200" fill="none" aria-hidden="true">
          <rect x="20" y="40" width="120" height="130" rx="10" fill="var(--card)" stroke="var(--rule)" stroke-width="1"/>
          <rect x="20" y="40" width="120" height="16" rx="10" fill="var(--surface-2)"/>
          <circle cx="30" cy="48" r="2.5" fill="#ff6058"/>
          <circle cx="38" cy="48" r="2.5" fill="#ffbd2e"/>
          <circle cx="46" cy="48" r="2.5" fill="#28ca42"/>
          <rect x="30" y="66" width="100" height="38" rx="6" fill="var(--accent-tint)"/>
          <circle cx="46" cy="85" r="8" fill="white"/>
          <rect x="58" y="78" width="60" height="4" rx="2" fill="var(--accent-text)"/>
          <rect x="58" y="86" width="40" height="3" rx="1.5" fill="var(--accent)"/>
          <rect x="30" y="112" width="100" height="3" rx="1.5" fill="var(--rule-strong)"/>
          <rect x="30" y="120" width="80" height="3" rx="1.5" fill="var(--rule-strong)"/>
          <rect x="30" y="128" width="92" height="3" rx="1.5" fill="var(--rule-strong)"/>
          <rect x="30" y="136" width="74" height="3" rx="1.5" fill="var(--rule-strong)"/>
          <rect x="30" y="148" width="60" height="6" rx="3" fill="var(--ink)"/>
          <path d="M150 96l1.4 2.6L154 100l-2.6 1.4L150 104l-1.4-2.6L146 100l2.6-1.4z" fill="var(--accent)"/>
          <path d="M155 86l.7 1.3 1.3.7-1.3.7-.7 1.3-.7-1.3-1.3-.7 1.3-.7z" fill="var(--accent)"/>
          <path d="M150 108h24m-6-5 6 5-6 5" stroke="var(--accent)" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"/>
          <rect x="184" y="56" width="124" height="98" rx="10" fill="white" stroke="var(--accent)" stroke-width="1.2"/>
          <circle cx="200" cy="74" r="9" fill="var(--accent-tint)"/>
          <rect x="214" y="70" width="54" height="5" rx="2.5" fill="var(--ink)"/>
          <rect x="214" y="80" width="78" height="3.5" rx="1.75" fill="var(--mute-2)"/>
          <rect x="194" y="100" width="60" height="14" rx="7" fill="var(--positive-tint)"/>
          <text x="222" y="111" font-family="system-ui" font-size="9" fill="var(--positive-text)" font-weight="600">Screen</text>
          <rect x="194" y="122" width="100" height="3" rx="1.5" fill="var(--rule)"/>
          <rect x="194" y="130" width="70" height="3" rx="1.5" fill="var(--rule)"/>
          <rect x="194" y="138" width="86" height="3" rx="1.5" fill="var(--rule)"/>
        </svg>
      {:else if step === 3}
        <svg viewBox="0 0 320 200" width="320" height="200" fill="none" aria-hidden="true">
          <rect x="92" y="32" width="136" height="136" rx="14" fill="white" stroke="var(--rule)" stroke-width="1"/>
          <circle cx="160" cy="78" r="22" fill="var(--accent-tint)"/>
          <text x="160" y="85" text-anchor="middle" font-family="system-ui" font-weight="600" font-size="16" fill="var(--accent-text)">SC</text>
          <rect x="116" y="110" width="88" height="5" rx="2.5" fill="var(--ink)"/>
          <rect x="124" y="120" width="72" height="4" rx="2" fill="var(--mute-2)"/>
          <rect x="108" y="138" width="40" height="14" rx="7" fill="var(--accent-tint)"/>
          <text x="128" y="148.5" text-anchor="middle" font-family="system-ui" font-size="9" fill="var(--accent-text)" font-weight="600">Talk</text>
          <rect x="154" y="138" width="40" height="14" rx="7" fill="var(--positive-tint)"/>
          <text x="174" y="148.5" text-anchor="middle" font-family="system-ui" font-size="9" fill="var(--positive-text)" font-weight="600">Post</text>
          <rect x="18" y="46" width="60" height="34" rx="8" fill="white" stroke="var(--rule)" stroke-width="1"/>
          <rect x="26" y="54" width="20" height="3" rx="1.5" fill="var(--mute-2)"/>
          <rect x="26" y="62" width="40" height="3" rx="1.5" fill="var(--ink-2)"/>
          <rect x="26" y="70" width="32" height="3" rx="1.5" fill="var(--ink-2)"/>
          <rect x="242" y="58" width="60" height="34" rx="8" fill="white" stroke="var(--rule)" stroke-width="1"/>
          <rect x="250" y="66" width="22" height="3" rx="1.5" fill="var(--mute-2)"/>
          <rect x="250" y="74" width="40" height="3" rx="1.5" fill="var(--ink-2)"/>
          <rect x="250" y="82" width="30" height="3" rx="1.5" fill="var(--ink-2)"/>
          <rect x="34" y="118" width="60" height="34" rx="8" fill="white" stroke="var(--rule)" stroke-width="1"/>
          <rect x="42" y="126" width="22" height="3" rx="1.5" fill="var(--mute-2)"/>
          <rect x="42" y="134" width="40" height="3" rx="1.5" fill="var(--ink-2)"/>
          <rect x="42" y="142" width="34" height="3" rx="1.5" fill="var(--ink-2)"/>
          <path d="M232 38l1.4 2.6L236 42l-2.6 1.4L232 46l-1.4-2.6L228 42l2.6-1.4z" fill="var(--accent)"/>
          <path d="M104 162l1.4 2.6L108 166l-2.6 1.4L104 170l-1.4-2.6L100 166l2.6-1.4z" fill="var(--accent)"/>
        </svg>
      {:else}
        <svg viewBox="0 0 320 200" width="320" height="200" fill="none" aria-hidden="true">
          <rect x="14" y="48" width="92" height="108" rx="12" fill="var(--surface)" stroke="var(--rule)" stroke-width="1"/>
          <circle cx="60" cy="84" r="20" fill="var(--accent-tint)"/>
          <path d="M52 80h16M52 84h12M52 88h14" stroke="var(--accent-text)" stroke-width="1.6" stroke-linecap="round"/>
          <rect x="32" y="118" width="56" height="4" rx="2" fill="var(--ink)"/>
          <rect x="38" y="130" width="44" height="3" rx="1.5" fill="var(--mute)"/>
          <rect x="114" y="44" width="92" height="112" rx="12" fill="var(--accent-tint)" stroke="var(--accent)" stroke-width="1.2"/>
          <circle cx="160" cy="80" r="20" fill="white"/>
          <path d="M150 76h20M150 80h20M150 84h14" stroke="var(--accent)" stroke-width="1.6" stroke-linecap="round"/>
          <path d="M172 70l.9 1.7 1.7.9-1.7.9-.9 1.7-.9-1.7-1.7-.9 1.7-.9z" fill="var(--accent)"/>
          <rect x="132" y="116" width="56" height="4" rx="2" fill="var(--accent-text)"/>
          <rect x="138" y="128" width="44" height="3" rx="1.5" fill="var(--accent-text)" opacity="0.7"/>
          <text x="160" y="148" text-anchor="middle" font-family="system-ui" font-size="9" font-weight="600" fill="var(--accent-text)">RECOMMENDED</text>
          <rect x="214" y="48" width="92" height="108" rx="12" fill="var(--surface)" stroke="var(--rule)" stroke-width="1"/>
          <circle cx="260" cy="84" r="20" fill="var(--warm-tint)"/>
          <circle cx="260" cy="84" r="10" stroke="var(--warm-text)" stroke-width="1.6" fill="none"/>
          <path d="M268 92l5 5" stroke="var(--warm-text)" stroke-width="1.6" stroke-linecap="round"/>
          <rect x="232" y="118" width="56" height="4" rx="2" fill="var(--ink)"/>
          <rect x="238" y="130" width="44" height="3" rx="1.5" fill="var(--mute)"/>
        </svg>
      {/if}
    </div>

    <!-- Stage text -->
    <div class="text">
      {#if step === 0}
        <h1>Welcome to Pursuit.</h1>
        <p>Track your job search end-to-end — applications, interviews, and AI briefings on the people you're about to meet.</p>
      {:else if step === 1}
        <h1>Paste any job, get a row.</h1>
        <p>Drop a JD URL, a body of text, or a recruiter email. Claude parses it into a clean application row in about three seconds.</p>
      {:else if step === 2}
        <h1>Screenshot, no scraping.</h1>
        <p>LinkedIn blocks bots — so just <kbd>⌘⌃⇧4</kbd> the page and paste. Vision reads the screenshot and fills the fields for you.</p>
      {:else if step === 3}
        <h1>A briefing on your interviewer.</h1>
        <p>Before any conversation, generate a 60-second brief: recent talks and posts, what tends to land, what to avoid, three questions worth asking.</p>
      {:else}
        <!-- Stage 4 — picker / paste / parsing / done -->
        {#if pickerMode === 'pick'}
          <h1>Pick where to start.</h1>
          <p class="picker-help">All three paths land you on the dashboard. You can clear demo data any time.</p>
          <div class="picker">
            <button class="pick" onclick={() => (pickerMode = 'paste')}>
              <div class="pick-h">Paste your applications</div>
              <div class="pick-s">URLs, text, emails — Claude parses them in.</div>
            </button>
            <button class="pick rec" onclick={seedDemo} disabled={seeding}>
              <span class="pick-tag">Recommended</span>
              <div class="pick-h">{seeding ? 'Seeding…' : 'Try with demo data'}</div>
              <div class="pick-s">15 realistic apps so you can see every surface populated.</div>
            </button>
            <button class="pick" onclick={skipForNow}>
              <div class="pick-h">Just explore</div>
              <div class="pick-s">Land on an empty dashboard. Add by hand later.</div>
            </button>
          </div>
          {#if globalError}<p class="err">{globalError}</p>{/if}
        {:else if pickerMode === 'paste'}
          <h1>Paste your applications.</h1>
          <p class="picker-help">URLs, JD body text, recruiter emails, or one-liners. <b>Separate entries with a blank line.</b></p>
          <textarea
            bind:value={text}
            rows="5"
            placeholder={"https://job-boards.greenhouse.io/anthropic/jobs/4020693008\n\nStripe — Staff Backend Engineer · referred by Mia · applied 14 May\n\nJD body text pasted from any company's careers page"}
          ></textarea>
          {#if globalError}<p class="err">{globalError}</p>{/if}
        {:else if pickerMode === 'parsing'}
          <h1>Adding your applications…</h1>
          <ul class="progress">
            {#each progress as p}
              <li class={`pr ${p.status}`}>
                <span class="pi">
                  {#if p.status === 'pending'}·{:else if p.status === 'parsing' || p.status === 'creating'}⟳{:else if p.status === 'done'}✓{:else}!{/if}
                </span>
                <span class="pl">{p.label}</span>
                {#if p.error}<span class="pe">— {p.error}</span>{/if}
              </li>
            {/each}
          </ul>
        {:else}
          <!-- done -->
          {#if added > 0}
            <h1>You're set.</h1>
            <p class="picker-help">Added {added} of {progress.length} {progress.length === 1 ? 'application' : 'applications'}.</p>
          {:else}
            <h1>Nothing added.</h1>
            <p class="picker-help">None of the entries could be parsed — you can add them by hand from the dashboard.</p>
          {/if}
          {#if progress.some(p => p.status === 'done')}
            <ul class="result-list ok">
              {#each progress.filter(p => p.status === 'done') as p}
                <li><span class="r-ic">✓</span><span class="r-lbl">{p.label}</span></li>
              {/each}
            </ul>
          {/if}
          {#if progress.some(p => p.status === 'error')}
            <ul class="result-list err-list">
              {#each progress.filter(p => p.status === 'error') as p}
                <li><span class="r-ic">!</span><span class="r-lbl">{p.label}</span>{#if p.error}<span class="r-err">— {p.error}</span>{/if}</li>
              {/each}
            </ul>
          {/if}
        {/if}
      {/if}
    </div>

    <!-- Footer -->
    <footer>
      {#if step > 0 && (step < total - 1 || (step === total - 1 && pickerMode !== 'pick' && pickerMode !== 'parsing'))}
        <button class="btn ghost" onclick={prev}>← Back</button>
      {/if}
      <span style="flex:1"></span>
      {#if step < total - 1}
        <button class="btn ghost" onclick={skipForNow}>Skip</button>
        <button class="btn primary" onclick={next}>Next →</button>
      {:else if pickerMode === 'paste'}
        <button class="btn primary" onclick={startPaste} disabled={!text.trim()}>
          {text.trim() ? 'Parse + open dashboard →' : 'Add some text first'}
        </button>
      {:else if pickerMode === 'done'}
        <button class="btn primary" onclick={markOnboarded}>Open the dashboard →</button>
      {/if}
    </footer>
  </div>
</div>

<style>
  .ob-overlay {
    position: fixed; inset: 0;
    background: radial-gradient(circle at 50% 30%, var(--accent-tint), rgba(10,10,13,0.55) 75%);
    display: grid; place-items: center;
    z-index: 200;
    padding: 24px;
    overflow-y: auto;
  }
  .ob-card {
    width: 100%;
    max-width: 560px;
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 18px;
    padding: 28px 30px 22px;
    box-shadow: 0 24px 80px -32px rgba(10,10,13,0.18), var(--sh-1);
    margin: 24px auto;
  }
  .dots {
    display: flex; justify-content: center; gap: 6px;
    margin-bottom: 22px;
  }
  .dot {
    width: 6px; height: 6px; border-radius: 50%;
    background: var(--rule-strong);
    transition: background 200ms, width 200ms;
  }
  .dot.done { background: var(--accent); }
  .dot.active { background: var(--accent); width: 24px; border-radius: 99px; }

  .art {
    display: grid; place-items: center;
    height: 200px;
    margin-bottom: 18px;
  }

  .text h1 {
    font-size: 22px; font-weight: 600;
    letter-spacing: -0.022em;
    margin: 0 0 8px;
    text-align: center;
  }
  .text > p {
    font-size: 14px; color: var(--mute);
    line-height: 1.55;
    margin: 0 auto 12px;
    text-align: center;
    max-width: 44ch;
  }
  .text kbd {
    font-family: var(--mono); font-size: 11px;
    background: var(--surface-2); border: 1px solid var(--rule);
    border-bottom-width: 2px; border-radius: 3px;
    padding: 0 4px; color: var(--ink-2);
  }
  .picker-help { margin-bottom: 16px !important; }

  .picker { display: flex; flex-direction: column; gap: 8px; margin: 4px 0 8px; }
  .pick {
    position: relative;
    text-align: left;
    background: var(--surface);
    border: 1px solid var(--rule);
    border-radius: 10px;
    padding: 12px 14px;
    cursor: pointer;
    transition: border-color 120ms, transform 120ms;
    font: inherit;
    color: inherit;
    width: 100%;
  }
  .pick:hover:not(:disabled) { border-color: var(--accent); transform: translateY(-1px); }
  .pick:disabled { opacity: 0.6; cursor: progress; }
  .pick.rec { background: var(--accent-tint); border-color: var(--accent); }
  .pick-tag {
    position: absolute; top: 10px; right: 12px;
    font-size: 10px; font-weight: 600; letter-spacing: 0.06em;
    color: var(--accent-text); text-transform: uppercase;
  }
  .pick-h { font-size: 14px; font-weight: 600; color: var(--ink); }
  .pick.rec .pick-h { color: var(--accent-text); }
  .pick-s { font-size: 12.5px; color: var(--mute); margin-top: 2px; }
  .pick.rec .pick-s { color: var(--accent-text); opacity: 0.8; }

  textarea {
    width: 100%;
    font: inherit;
    font-size: 13.5px; line-height: 1.5;
    color: var(--ink); background: var(--surface);
    border: 1px solid var(--rule);
    border-radius: 10px;
    padding: 12px 14px;
    outline: none; resize: vertical;
    font-family: var(--sans);
  }
  textarea:focus { border-color: var(--accent); }
  textarea::placeholder { color: var(--mute-2); }

  .progress {
    list-style: none; padding: 12px 14px; margin: 8px 0 0;
    display: flex; flex-direction: column; gap: 6px;
    background: var(--surface); border: 1px solid var(--rule);
    border-radius: 10px;
    max-height: 220px; overflow-y: auto;
  }
  .pr {
    display: grid; grid-template-columns: 16px 1fr auto;
    gap: 8px; font-size: 12.5px; color: var(--ink-2);
    align-items: baseline;
  }
  .pi { font-family: var(--mono); color: var(--mute-2); text-align: center; }
  .pr.done .pi { color: var(--positive-text); }
  .pr.error .pi { color: var(--danger-text); }
  .pr.parsing .pi, .pr.creating .pi { color: var(--accent-text); animation: spin 1s linear infinite; }
  @keyframes spin { to { transform: rotate(360deg); } }
  .pl { color: var(--ink-2); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .pe { color: var(--danger-text); font-size: 11.5px; }
  .err { color: var(--danger-text); font-size: 12.5px; margin-top: 8px; }

  .result-list {
    list-style: none; padding: 8px 12px; margin: 8px 0;
    background: var(--surface); border: 1px solid var(--rule);
    border-radius: 10px;
    display: flex; flex-direction: column; gap: 4px;
  }
  .result-list.err-list { background: var(--danger-tint); border-color: var(--danger-tint); }
  .result-list li {
    display: grid; grid-template-columns: 18px 1fr;
    gap: 6px; font-size: 13px; color: var(--ink-2);
  }
  .result-list .r-ic { text-align: center; font-family: var(--mono); font-size: 12px; }
  .result-list.ok .r-ic { color: var(--positive-text); }
  .result-list.err-list .r-ic { color: var(--danger-text); }
  .result-list .r-lbl { overflow: hidden; text-overflow: ellipsis; }
  .result-list .r-err { grid-column: 2; color: var(--danger-text); font-size: 12px; }

  footer {
    display: flex; align-items: center; gap: 8px;
    margin-top: 18px;
    border-top: 1px solid var(--rule);
    padding-top: 14px;
  }
  .btn {
    height: 32px;
    padding: 0 14px;
    border-radius: 7px;
    font: inherit; font-size: 13px; font-weight: 500;
    border: 1px solid var(--rule);
    background: var(--card);
    color: var(--ink-2);
    cursor: pointer;
  }
  .btn:disabled { opacity: 0.55; cursor: not-allowed; }
  .btn.ghost { background: transparent; border-color: transparent; color: var(--mute); }
  .btn.ghost:hover { color: var(--ink); background: var(--surface-2); }
  .btn.primary { background: var(--accent); color: white; border-color: var(--accent); }
  .btn.primary:hover:not(:disabled) { background: var(--accent-strong); }

  @media (max-width: 720px) {
    .ob-card { padding: 22px 18px 18px; }
    .text h1 { font-size: 20px; }
  }
</style>
