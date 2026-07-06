<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api.js';
  import { logEvent } from '$lib/analytics.js';

  // onDone marks the user onboarded + unmounts the overlay (shared with the tour).
  let { onDone = () => {} } = $props();

  let phase = $state('prompt'); // 'prompt' | 'generating' | 'error'
  let company = $state('');
  let role = $state('');
  let createdId = $state(null); // set once the application exists, so error/skip can route there
  let createdFor = ''; // company|role the created app belongs to (guards retry reuse)
  let errMsg = $state('');
  const ready = $derived(!!company.trim() && !!role.trim());

  onMount(() => logEvent('prepfirst_prompt_view'));

  async function build() {
    if (phase === 'generating' || !ready) return;
    logEvent('prepfirst_submit', { has_role: !!role.trim() });
    phase = 'generating';
    errMsg = '';
    try {
      // The prep question creates the first tracked application as a byproduct —
      // the tracker spine stays intact; we just lead with the wedge. Company AND
      // role are both required by POST /applications (prep is role-specific).
      // On retry after a generate failure, reuse the app we already created —
      // unless they edited the company/role, then start fresh.
      const key = `${company.trim()}|${role.trim()}`;
      if (createdId && key !== createdFor) createdId = null;
      if (!createdId) {
        createdFor = key;
        const a = await api('/api/applications', {
          method: 'POST',
          body: JSON.stringify({ company: company.trim(), role: role.trim(), status: 'screen' })
        });
        createdId = a.id;
        logEvent('application_create', { source: 'prepfirst' });
      }
      // No interview_id → the shared company brief (interviewer-optional). This is
      // the cold-start wow; round-by-round prep is added later on the detail page.
      await api(`/api/applications/${createdId}/dossier/refresh`, { method: 'POST', body: '{}' });
      logEvent('prepfirst_generate_ok');
      await onDone();
      goto(`/app/${createdId}?welcome=1`);
    } catch (e) {
      console.error(e);
      errMsg = e?.message || '';
      // step: did the create or the brief generation fail? reason: capped error string.
      logEvent('prepfirst_generate_error', { step: createdId ? 'generate' : 'create', reason: (e?.message || 'unknown').slice(0, 120) });
      phase = 'error';
    }
  }

  function backToPrompt() { phase = 'prompt'; }

  async function skip() {
    logEvent('prepfirst_skip');
    await onDone();
    if (createdId) goto(`/app/${createdId}`);
  }

  async function retry() {
    // The application already exists — just send them to it to generate there.
    await onDone();
    goto(createdId ? `/app/${createdId}` : '/app');
  }
</script>

<div class="pf-back">
  <div class="pf-card">
    {#if phase === 'prompt'}
      <div class="pf-badge"><span class="spark">✦</span> AI interview prep</div>
      <h1>Who are you interviewing with?</h1>
      <p class="pf-sub">Name a company you're interviewing with and Pursuit builds your playbook — what they do, where they're headed, the typical loop, and what the team grades for.</p>
      <form onsubmit={(e) => { e.preventDefault(); build(); }}>
        <label class="pf-field">
          <span>Company</span>
          <input class="pf-input" bind:value={company} placeholder="e.g. Stripe" autofocus />
        </label>
        <label class="pf-field">
          <span>Role you're interviewing for</span>
          <input class="pf-input" bind:value={role} placeholder="e.g. Staff Software Engineer" />
        </label>
        <button class="pf-cta" type="submit" disabled={!ready}>Build my playbook</button>
      </form>
      <button class="pf-skip" type="button" onclick={skip}>I'm just exploring — skip →</button>

    {:else if phase === 'generating'}
      <div class="pf-badge"><span class="spark">✦</span> AI interview prep</div>
      <h1>Researching {company}…</h1>
      <div class="pf-spinner"></div>
      <p class="pf-sub">Building your company playbook from public sources. This usually takes 1–2 minutes — hang tight, it'll open as soon as it's ready.</p>

    {:else}
      <div class="pf-badge err">Couldn't build it</div>
      <h1>That didn't go through</h1>
      <p class="pf-sub">{errMsg || 'We hit a snag building your playbook.'}{#if createdId} Your {company} application was saved — you can try again from its page.{/if}</p>
      {#if createdId}
        <button class="pf-cta" type="button" onclick={retry}>Go to {company}</button>
      {:else}
        <button class="pf-cta" type="button" onclick={backToPrompt}>Try again</button>
      {/if}
      <button class="pf-skip" type="button" onclick={skip}>Skip to my dashboard →</button>
    {/if}
  </div>
</div>

<style>
  .pf-back { position: fixed; inset: 0; z-index: 130; display: grid; place-items: center;
    background: var(--surface); padding: 24px; }
  .pf-card { width: 440px; max-width: 100%; background: var(--card); border: 1px solid var(--rule);
    border-radius: 18px; padding: 34px 32px 26px; box-shadow: var(--sh-pop);
    text-align: center; animation: pf-rise .28s cubic-bezier(.2,.7,.3,1) both; }
  @keyframes pf-rise { from { transform: translateY(10px); opacity: 0; } to { transform: none; opacity: 1; } }

  .pf-badge { display: inline-flex; align-items: center; gap: 6px; font-size: 12px; font-weight: 600;
    color: var(--accent-text); background: var(--accent-tint); padding: 5px 11px; border-radius: 99px; margin-bottom: 16px; }
  .pf-badge.err { color: var(--danger-text); background: var(--danger-tint); }
  .pf-badge .spark { font-size: 12px; }

  h1 { font-size: 23px; font-weight: 600; letter-spacing: -0.025em; margin: 0 0 10px; color: var(--ink); line-height: 1.2; }
  .pf-sub { font-size: 13.5px; line-height: 1.55; color: var(--mute); margin: 0 0 22px; text-wrap: pretty; }

  form { display: flex; flex-direction: column; gap: 13px; text-align: left; }
  .pf-field { display: flex; flex-direction: column; gap: 5px; }
  .pf-field span { font-size: 12px; font-weight: 500; color: var(--ink-2); }
  .pf-field em { color: var(--mute-2); font-style: normal; font-weight: 400; }
  .pf-input { height: 42px; padding: 0 13px; border: 1px solid var(--rule); border-radius: 10px;
    background: var(--card); color: var(--ink); font: 400 14.5px/1.2 var(--sans); outline: none; }
  .pf-input::placeholder { color: var(--mute-2); }
  .pf-input:focus { border-color: var(--accent); box-shadow: 0 0 0 3px var(--accent-tint); }

  .pf-cta { height: 44px; margin-top: 6px; background: var(--accent); color: #fff; border: none;
    border-radius: 10px; font: 600 14.5px/1 var(--sans); cursor: pointer; transition: background .12s; }
  .pf-cta:hover { background: var(--accent-strong); }
  .pf-cta:disabled { opacity: .5; cursor: not-allowed; }

  .pf-skip { margin-top: 16px; background: none; border: none; color: var(--mute); font: 500 13px/1 var(--sans);
    cursor: pointer; padding: 6px; }
  .pf-skip:hover { color: var(--ink-2); }

  .pf-spinner { width: 30px; height: 30px; margin: 6px auto 18px; border-radius: 50%;
    border: 3px solid var(--accent-tint); border-top-color: var(--accent); animation: pf-spin .8s linear infinite; }
  @keyframes pf-spin { to { transform: rotate(360deg); } }

  @media (prefers-reduced-motion: reduce) {
    .pf-card { animation: none; }
    .pf-spinner { animation-duration: 1.6s; }
  }
</style>
