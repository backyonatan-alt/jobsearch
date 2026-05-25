<script>
  import { api } from '$lib/api.js';

  let { onDone } = $props();

  let text = $state('');
  let working = $state(false);
  let progress = $state([]); // [{idx, status: 'pending'|'parsing'|'creating'|'done'|'error', label, error?}]
  let globalError = $state('');
  let finished = $state(false);
  const added = $derived(progress.filter(p => p.status === 'done').length);

  function splitEntries(s) {
    const parts = s.split(/\n{2,}/).map(p => p.trim()).filter(p => p.length > 4);
    if (parts.length > 1) return parts;
    // Fall back to per-line splitting if everything looks URL-shaped.
    const lines = s.split(/\n/).map(p => p.trim()).filter(p => p.length > 4);
    if (lines.length > 1 && lines.every(l => /^https?:\/\//.test(l))) return lines;
    return parts.length ? parts : (s.trim() ? [s.trim()] : []);
  }

  async function start() {
    const entries = splitEntries(text);
    if (entries.length === 0) {
      // No applications — just mark onboarded and let them in.
      await markOnboarded();
      return;
    }
    working = true;
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

        // Map only the fields the create endpoint accepts — don't spread the
        // parser response blindly. `seniority`, for example, is parsed but not
        // stored. Backend's DisallowUnknownFields would reject extras.
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

    finished = true;
    working = false;
    // If at least one was added, mark onboarded automatically and let the
    // user click Done. If none, they can still click Done to skip.
  }

  async function markOnboarded() {
    try {
      await api('/api/me/onboarded', { method: 'POST' });
    } catch (e) {
      console.error('mark onboarded', e);
    }
    onDone?.();
  }

  function skipForNow() {
    markOnboarded();
  }
</script>

<div class="ob-overlay" role="dialog" aria-modal="true">
  <div class="ob-card" onclick={(e) => e.stopPropagation()}>
    <header>
      <h1>Welcome to Pursuit.</h1>
      <p class="lede">A job-search tracker with two AI moments most trackers don't have. Sixty seconds to get oriented.</p>
    </header>

    {#if !finished}
      <ul class="feature-pills">
        <li>
          <div class="fp-ic">
            <svg width="18" height="18" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="3" y="2" width="10" height="12" rx="1.5"/><path d="M6 5h4M6 8h4M6 11h3"/></svg>
          </div>
          <div>
            <h3>Paste any job</h3>
            <p>JD body, URL, recruiter email, or a one-line note. Claude parses it into a row.</p>
          </div>
        </li>
        <li>
          <div class="fp-ic">
            <svg width="18" height="18" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="2" y="3" width="12" height="10" rx="1.5"/><circle cx="5.5" cy="6.5" r="1"/><path d="M2 11l3.5-3.5 3 3 2-2L14 11"/></svg>
          </div>
          <div>
            <h3>Drop a screenshot</h3>
            <p>LinkedIn blocks scrapers — so screenshot the page (<kbd>⌘⌃⇧4</kbd>) and paste. Vision reads it.</p>
          </div>
        </li>
        <li>
          <div class="fp-ic">
            <svg width="18" height="18" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="8" cy="6" r="2.5"/><path d="M3 14c0-2.8 2.2-5 5-5s5 2.2 5 5"/></svg>
          </div>
          <div>
            <h3>AI interviewer dossier</h3>
            <p>Before an interview, get a 60-sec brief on the person — recent posts, style, watch-fors. Live web search.</p>
          </div>
        </li>
      </ul>

      <section class="step">
        <div class="step-head">
          <span class="step-num">1</span>
          <h2>Seed your search now</h2>
        </div>
        <p class="step-help">
          Paste any jobs you've already applied to — URLs, JD body text, recruiter
          emails, or your own one-liners. <b>Separate entries with a blank line.</b>
          Or skip and add them one-by-one with <kbd class="hk">⌘N</kbd> later.
        </p>
        <textarea
          bind:value={text}
          rows="6"
          disabled={working}
          placeholder={"https://job-boards.greenhouse.io/anthropic/jobs/4020693008\n\nStripe — Staff Backend Engineer · referred by Mia · applied 14 May\n\nJD body text pasted from any company's careers page"}
        ></textarea>

        {#if progress.length > 0}
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
        {/if}
      </section>

      <section class="step">
        <div class="step-head">
          <span class="step-num">2</span>
          <h2>The shortcuts that matter</h2>
        </div>
        <ul class="tips">
          <li><kbd class="hk">⌘N</kbd> anywhere opens the New-application modal. Drop a screenshot or paste a URL; it parses in ~3s.</li>
          <li>Open any application → <b>Refresh dossier</b> generates the AI briefing (interview-status apps benefit the most).</li>
          <li>Sidebar pipeline links (Interview loops / Open offers / Wishlist) jump straight to filtered views.</li>
          <li><b>Today</b> surfaces the live loop and recent offers; <b>Board</b> is Kanban; <b>Funnel</b> shows your conversion rates.</li>
        </ul>
      </section>

      <footer class="actions">
        <button class="btn" onclick={skipForNow} disabled={working}>Skip — just explore</button>
        <button class="btn btn-primary" onclick={start} disabled={working}>
          {working ? 'Adding…' : (text.trim() ? 'Parse + open dashboard →' : 'Open the dashboard →')}
        </button>
      </footer>
    {:else}
      <section class="done">
        {#if added > 0}
          <h2>You're set.</h2>
          <p>
            Added {added} of {progress.length}
            {progress.length === 1 ? 'application' : 'applications'}.
          </p>
        {:else}
          <h2>Nothing added.</h2>
          <p>
            None of the {progress.length} {progress.length === 1 ? 'entry' : 'entries'} could be parsed.
            See below — you can add them by hand from the dashboard.
          </p>
        {/if}

        {#if progress.some(p => p.status === 'done')}
          <ul class="result-list ok">
            {#each progress.filter(p => p.status === 'done') as p}
              <li><span class="r-ic">✓</span><span class="r-lbl">{p.label}</span></li>
            {/each}
          </ul>
        {/if}

        {#if progress.some(p => p.status === 'error')}
          <p class="failed-head">Didn't parse:</p>
          <ul class="result-list err">
            {#each progress.filter(p => p.status === 'error') as p}
              <li>
                <span class="r-ic">!</span>
                <span class="r-lbl">{p.label}</span>
                <span class="r-err">— {p.error}</span>
              </li>
            {/each}
          </ul>
        {/if}

        <div class="done-actions">
          <button class="btn btn-primary" onclick={markOnboarded}>Open the dashboard →</button>
        </div>
      </section>
    {/if}
  </div>
</div>

<style>
  .ob-overlay {
    position: fixed; inset: 0;
    background: rgba(10,10,13,0.55);
    display: grid; place-items: center;
    z-index: 200;
    padding: 24px;
    overflow-y: auto;
  }
  .ob-card {
    width: 100%;
    max-width: 640px;
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 14px;
    padding: 32px;
    box-shadow: var(--sh-pop);
    margin: 24px auto;
  }
  header h1 { font-size: 26px; font-weight: 500; letter-spacing: -0.022em; margin: 0 0 .4rem; }
  .lede { color: var(--mute); margin: 0 0 24px; font-size: 14px; line-height: 1.5; }

  .feature-pills {
    list-style: none; padding: 0; margin: 0 0 28px;
    display: grid; grid-template-columns: 1fr; gap: 8px;
  }
  .feature-pills li {
    display: grid; grid-template-columns: 36px 1fr; gap: 12px;
    align-items: start;
    padding: 12px 14px;
    background: var(--surface);
    border: 1px solid var(--rule);
    border-radius: 10px;
  }
  .feature-pills .fp-ic {
    width: 32px; height: 32px;
    border-radius: 7px;
    background: var(--accent-tint);
    color: var(--accent-text);
    display: grid; place-items: center;
  }
  .feature-pills h3 { font-size: 13.5px; font-weight: 500; margin: 0 0 2px; color: var(--ink); letter-spacing: -0.005em; }
  .feature-pills p { font-size: 12.5px; color: var(--mute); margin: 0; line-height: 1.5; }
  .feature-pills kbd {
    font-family: var(--mono); font-size: 10.5px;
    background: var(--card); border: 1px solid var(--rule);
    border-bottom-width: 2px; border-radius: 3px;
    padding: 0 4px; color: var(--ink-2);
  }

  .hk {
    font-family: var(--mono); font-size: 11.5px;
    background: var(--surface-2); border: 1px solid var(--rule);
    border-bottom-width: 2px; border-radius: 4px;
    padding: 0 5px; color: var(--ink-2);
  }
  .tips kbd { font-family: var(--mono); font-size: 11.5px; background: var(--card); border: 1px solid var(--rule); border-bottom-width: 2px; border-radius: 4px; padding: 0 5px; color: var(--ink-2); }
  .step-help kbd { font-family: var(--mono); font-size: 11.5px; background: var(--card); border: 1px solid var(--rule); border-bottom-width: 2px; border-radius: 4px; padding: 0 5px; color: var(--ink-2); }

  .step { margin-bottom: 24px; }
  .step-head {
    display: flex; align-items: center; gap: 10px;
    margin-bottom: 10px;
  }
  .step-num {
    width: 20px; height: 20px;
    border-radius: 50%;
    background: var(--ink); color: white;
    font-family: var(--mono); font-size: 11px;
    display: grid; place-items: center;
    font-weight: 500;
  }
  .step h2 { font-size: 15px; font-weight: 500; letter-spacing: -0.01em; margin: 0; }
  .step-help { font-size: 13.5px; color: var(--mute); margin: 0 0 12px; line-height: 1.5; }
  .step-help b { color: var(--ink-2); font-weight: 500; }

  textarea {
    width: 100%;
    font: inherit;
    font-size: 13px;
    color: var(--ink);
    background: var(--surface);
    border: 1px solid var(--rule);
    border-radius: 8px;
    padding: 10px 12px;
    outline: none;
    resize: vertical;
    line-height: 1.5;
    font-family: var(--sans);
  }
  textarea:focus { border-color: var(--accent); }
  textarea::placeholder { color: var(--mute-2); }

  .progress {
    list-style: none; padding: 0; margin: 12px 0 0;
    display: flex; flex-direction: column; gap: 4px;
    background: var(--surface);
    border: 1px solid var(--rule);
    border-radius: 8px;
    padding: 10px 12px;
    max-height: 200px;
    overflow-y: auto;
  }
  .pr {
    display: grid;
    grid-template-columns: 16px 1fr auto;
    gap: 8px;
    font-size: 12.5px;
    color: var(--ink-2);
    align-items: baseline;
  }
  .pi {
    font-family: var(--mono);
    color: var(--mute-2);
    text-align: center;
  }
  .pr.done .pi { color: var(--positive-text); }
  .pr.error .pi { color: var(--danger-text); }
  .pr.parsing .pi, .pr.creating .pi { color: var(--accent-text); animation: spin 1s linear infinite; }
  @keyframes spin { to { transform: rotate(360deg); } }
  .pl { color: var(--ink-2); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .pe { color: var(--danger-text); font-size: 11.5px; }

  .tips {
    list-style: none; padding: 0; margin: 0;
    display: flex; flex-direction: column; gap: 8px;
  }
  .tips li {
    padding: 10px 14px;
    background: var(--surface);
    border: 1px solid var(--rule);
    border-radius: 8px;
    font-size: 13px;
    line-height: 1.5;
    color: var(--ink-2);
  }
  .tips b { color: var(--ink); font-weight: 500; }
  .tips em { font-style: italic; color: var(--accent-text); }

  .actions {
    display: flex; justify-content: flex-end; gap: 8px;
    margin-top: 28px;
  }
  .btn {
    height: 32px;
    padding: 0 14px;
    border-radius: 6px;
    font: inherit;
    font-size: 13px;
    font-weight: 500;
    border: 1px solid var(--rule);
    background: var(--card);
    color: var(--ink-2);
    cursor: pointer;
  }
  .btn:hover { background: var(--surface-2); }
  .btn-primary {
    background: var(--ink);
    color: white;
    border-color: var(--ink);
  }
  .btn-primary:hover { background: #1a1a1f; }
  .btn:disabled { opacity: 0.55; cursor: not-allowed; }

  .done { padding: 8px 0 0; }
  .done h2 { font-size: 22px; font-weight: 500; letter-spacing: -0.018em; margin: 0 0 .4rem; text-align: center; }
  .done > p { color: var(--mute); margin: 0 0 20px; font-size: 14px; line-height: 1.55; text-align: center; }

  .failed-head {
    margin: 16px 0 6px;
    font-size: 12px;
    font-weight: 500;
    letter-spacing: .04em;
    text-transform: uppercase;
    color: var(--danger-text);
  }
  .result-list {
    list-style: none; padding: 0; margin: 0 0 12px;
    background: var(--surface);
    border: 1px solid var(--rule);
    border-radius: 8px;
    overflow: hidden;
  }
  .result-list.err { border-color: var(--danger-tint); background: var(--danger-tint); }
  .result-list li {
    display: grid;
    grid-template-columns: 18px 1fr;
    gap: 6px;
    padding: 8px 12px;
    font-size: 13px;
    color: var(--ink-2);
    border-top: 1px solid rgba(0,0,0,0.04);
  }
  .result-list li:first-child { border-top: none; }
  .result-list .r-ic {
    text-align: center;
    font-family: var(--mono);
    font-size: 12px;
  }
  .result-list.ok .r-ic { color: var(--positive-text); }
  .result-list.err .r-ic { color: var(--danger-text); }
  .result-list .r-lbl { overflow: hidden; text-overflow: ellipsis; }
  .result-list .r-err {
    grid-column: 2;
    color: var(--danger-text);
    font-size: 12px;
    line-height: 1.4;
  }

  .done-actions { text-align: center; margin-top: 20px; }
</style>
