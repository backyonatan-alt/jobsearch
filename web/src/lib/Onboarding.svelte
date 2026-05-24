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

    let added = 0;
    for (let i = 0; i < entries.length; i++) {
      progress[i] = { ...progress[i], status: 'parsing' };
      try {
        const parsed = await api('/api/applications/parse', {
          method: 'POST',
          body: JSON.stringify({ text: entries[i] })
        });
        progress[i] = { ...progress[i], status: 'creating', label: `${parsed.company || '?'} · ${parsed.role || '?'}` };

        const payload = { ...parsed, status: 'applied' };
        for (const k of Object.keys(payload)) if (!payload[k]) delete payload[k];
        await api('/api/applications', { method: 'POST', body: JSON.stringify(payload) });

        progress[i] = { ...progress[i], status: 'done' };
        added++;
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
      <p class="lede">A working tool for your job search. Two minutes to get oriented.</p>
    </header>

    {#if !finished}
      <section class="step">
        <div class="step-head">
          <span class="step-num">1</span>
          <h2>Add the jobs you've applied to</h2>
        </div>
        <p class="step-help">
          Paste any job postings, URLs, recruiter emails, or notes in your own
          words. <b>Separate entries with a blank line.</b> Claude will parse
          each one into a row.
        </p>
        <textarea
          bind:value={text}
          rows="8"
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
          <h2>What you can do next</h2>
        </div>
        <ul class="tips">
          <li><b>Update status</b> on any application from its detail page — the pipeline runs wishlist → applied → screen → interview → offer.</li>
          <li><b>AI dossier</b>: once an application is in Screen / Interview / Offer, open it and click <em>Generate</em>. Claude researches the interviewer or company live and writes you a 60-second briefing.</li>
          <li><b>Today</b> in the sidebar surfaces what needs your attention now — open offers, upcoming loops.</li>
        </ul>
      </section>

      <footer class="actions">
        <button class="btn" onclick={skipForNow} disabled={working}>Skip for now</button>
        <button class="btn btn-primary" onclick={start} disabled={working}>
          {working ? 'Adding…' : (text.trim() ? 'Add applications + start' : 'Start with no data')}
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
  .lede { color: var(--mute); margin: 0 0 28px; font-size: 14px; }

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
