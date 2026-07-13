<script>
  import { page } from '$app/state';
  import { track } from '$lib/analytics.js';
  import { isInAppBrowser } from '$lib/webview.js';

  // In-app browsers (LinkedIn/IG/FB) can't complete Google OAuth — show a
  // nudge to open in a real browser instead of letting users hit the 403.
  let inApp = $state(false);
  let copied = $state(false);
  let nudgeTracked = false;
  $effect(() => {
    // ?webview=1 forces the nudge on for preview/QA in a normal browser.
    inApp = isInAppBrowser() || page.url.searchParams.get('webview') === '1';
    if (inApp && !nudgeTracked) {
      nudgeTracked = true;
      track('signin_webview_nudge', { source: source || 'direct' });
    }
  });

  async function copyLink() {
    try {
      await navigator.clipboard.writeText(window.location.href);
      copied = true;
      setTimeout(() => (copied = false), 2000);
    } catch {
      copied = false;
    }
  }

  const messages = {
    oauth_denied:  'Sign-in was canceled.',
    oauth_state:   'Your sign-in session expired. Try again.',
    oauth_no_code: "Sign-in didn't return a code. Try again.",
    oauth_failed:  "We couldn't verify your Google account. Try again.",
    not_invited:   "That Google account isn't on the beta invite list yet — drop your email below and the admin will see it.",
    internal:      'Something went wrong on our end. Try again.'
  };

  const err = $derived(page.url.searchParams.get('err'));
  const statusMessage = $derived(err ? messages[err] || null : null);
  // If they bounced off "not_invited", auto-expand the request-access form
  // so the affordance is right where their eyes already are.
  let showInterest = $state(false);
  $effect(() => { if (err === 'not_invited') showInterest = true; });

  // Acquisition attribution: ?src=li etc. survives the OAuth round-trip via
  // localStorage and lands on the first-party `login` event.
  $effect(() => {
    const src = page.url.searchParams.get('src');
    if (src && typeof localStorage !== 'undefined') {
      try { localStorage.setItem('pursuit_src', src.slice(0, 40)); } catch {}
    }
  });

  let email = $state('');
  let note = $state('');
  let submitting = $state(false);
  let interestState = $state(''); // '' | 'sent' | 'already_invited'
  let interestError = $state('');
  const source = $derived(page.url.searchParams.get('src') || '');

  async function submitInterest(e) {
    e.preventDefault();
    interestError = '';
    submitting = true;
    try {
      const r = await fetch('/api/beta-interest', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: email.trim(), note: note.trim(), source })
      });
      const body = await r.json().catch(() => ({}));
      if (!r.ok) throw new Error(body.error || 'Something went wrong.');
      interestState = body.status === 'already_invited' ? 'already_invited' : 'sent';
      track('beta_interest_submit', { source: source || 'direct', outcome: interestState });
    } catch (e) {
      interestError = e.message;
      track('beta_interest_submit', { source: source || 'direct', outcome: 'error' });
    } finally {
      submitting = false;
    }
  }
</script>

<svelte:head>
  <title>Pursuit — sign in</title>
</svelte:head>

<main class="auth-shell">
  <div class="auth-card">
    <div class="mark"></div>
    <h1>Pursuit</h1>
    <p class="tagline">Walk into every interview with a playbook.</p>
    <p class="subtagline">Pursuit researches your interviewer and the company, then builds you an AI prep playbook — plus tracks every application.</p>

    {#if inApp}
      <div class="webview-note" role="note">
        <strong>To sign in with Google, open this page in your browser.</strong>
        <span>Tap the <b>•••</b> (or share) menu at the top, choose <b>“Open in Safari / Chrome”</b>, then sign in there.</span>
        <button type="button" class="wv-copy" onclick={copyLink}>{copied ? 'Link copied ✓' : 'Copy link'}</button>
      </div>
    {/if}

    <a class="google" href="/auth/google/start" data-sveltekit-reload>
      <svg viewBox="0 0 18 18" width="18" height="18" aria-hidden="true">
        <path d="M17.64 9.2c0-.637-.057-1.251-.164-1.84H9v3.481h4.844a4.14 4.14 0 0 1-1.796 2.716v2.258h2.908c1.702-1.567 2.684-3.875 2.684-6.615z" fill="#4285F4"/>
        <path d="M9 18c2.43 0 4.467-.806 5.956-2.18l-2.908-2.259c-.806.54-1.837.86-3.048.86-2.344 0-4.328-1.584-5.036-3.711H.957v2.332A8.997 8.997 0 0 0 9 18z" fill="#34A853"/>
        <path d="M3.964 10.71A5.41 5.41 0 0 1 3.682 9c0-.593.102-1.17.282-1.71V4.958H.957A8.996 8.996 0 0 0 0 9c0 1.452.348 2.827.957 4.042l3.007-2.332z" fill="#FBBC05"/>
        <path d="M9 3.58c1.321 0 2.508.454 3.44 1.345l2.582-2.58C13.463.891 11.426 0 9 0A8.997 8.997 0 0 0 .957 4.958L3.964 7.29C4.672 5.163 6.656 3.58 9 3.58z" fill="#EA4335"/>
      </svg>
      <span>Continue with Google</span>
    </a>

    {#if statusMessage}
      <p class="status error" role="status">{statusMessage}</p>
    {:else}
      <p class="status">&nbsp;</p>
    {/if}

    <div class="interest-block">
      {#if !showInterest && interestState === ''}
        <!-- Open beta: Google sign-in is the door; the interest form stays only as
             the fallback surface for the not_invited error (OPEN_SIGNUP kill switch). -->
      {:else if interestState === 'sent'}
        <div class="interest-done">
          <h3>Thanks — you're on the list.</h3>
          <p>The admin sees your email on their dashboard. When you're invited you'll be able to sign in with this Google account.</p>
        </div>
      {:else if interestState === 'already_invited'}
        <div class="interest-done">
          <h3>You're already on the invite list.</h3>
          <p>Click <b>Continue with Google</b> above with that email.</p>
        </div>
      {:else}
        <form class="interest-form" onsubmit={submitInterest}>
          <h3>Request access</h3>
          <p class="interest-help">Pursuit is in closed beta — AI interview playbooks + a tracker for your whole search. Drop your email and one line on what you're job-searching for; the admin reviews and invites manually.</p>
          <label>
            <span class="lbl">Your Gmail</span>
            <input type="email" bind:value={email} placeholder="you@gmail.com" required />
          </label>
          <label>
            <span class="lbl">One line about your search <span class="opt">— optional</span></span>
            <input type="text" bind:value={note} placeholder="Looking for staff roles in AI infra, US remote" maxlength="500" />
          </label>
          {#if interestError}<p class="ifail">{interestError}</p>{/if}
          <div class="iactions">
            <button type="button" class="ghost" onclick={() => (showInterest = false)}>Cancel</button>
            <button type="submit" class="primary" disabled={submitting || !email.trim()}>
              {submitting ? 'Sending…' : 'Request access'}
            </button>
          </div>
        </form>
      {/if}
    </div>

    <p class="footnote">Free while in beta. Sign in with Google — no waitlist.</p>
    <p class="privacy">Everything you enter — applications, notes, salary, prep briefs — is private to your account; no other user sees it. We store your Google email and profile picture for sign-in. Briefs are generated with Anthropic's API, which doesn't train on your data. Want out? One email and we delete everything. <a href="/privacy">How we research interviewers →</a></p>
  </div>
</main>

<style>
  .auth-shell {
    min-height: 100vh;
    display: grid; place-items: center;
    background: var(--surface);
    padding: 32px 16px;
  }
  .auth-card {
    width: 100%; max-width: 380px;
    background: var(--card);
    border: 1px solid var(--rule);
    border-radius: 14px;
    padding: 32px 28px 24px;
    box-shadow: var(--sh-pop);
    text-align: center;
  }
  .mark {
    width: 36px; height: 36px;
    background: var(--ink);
    border-radius: 8px;
    margin: 0 auto 16px;
  }
  h1 { font-size: 28px; font-weight: 500; letter-spacing: -0.025em; margin: 0 0 6px; }
  .tagline { color: var(--ink); margin: 0 0 8px; font-size: 16px; font-weight: 500; letter-spacing: -0.01em; }
  .subtagline { color: var(--mute); margin: 0 0 22px; font-size: 13px; line-height: 1.5; }

  .google {
    display: flex; align-items: center; justify-content: center;
    gap: 10px;
    height: 42px;
    background: var(--ink); color: white;
    border: 1px solid var(--ink);
    border-radius: 8px;
    text-decoration: none;
    font-size: 14px; font-weight: 500;
    transition: background 100ms ease;
  }
  .google:hover { background: #1a1a1f; }
  .google svg { background: white; border-radius: 3px; padding: 2px; box-sizing: content-box; }

  .status { font-size: 12.5px; min-height: 18px; margin: 12px 0 0; color: var(--mute); }
  .status.error { color: var(--danger-text); }

  .webview-note {
    text-align: left;
    background: var(--warm-tint);
    border: 1px solid var(--warm);
    border-radius: 10px;
    padding: 12px 14px;
    margin: 0 0 14px;
    display: flex; flex-direction: column; gap: 6px;
  }
  .webview-note strong { font-size: 13px; color: var(--warm-text); font-weight: 600; }
  .webview-note span { font-size: 12.5px; color: var(--ink-2); line-height: 1.5; }
  .webview-note b { font-weight: 600; }
  .wv-copy {
    align-self: flex-start; margin-top: 2px;
    font: inherit; font-size: 12px; font-weight: 500;
    height: 28px; padding: 0 12px;
    background: var(--card); color: var(--ink);
    border: 1px solid var(--rule); border-radius: 6px;
    cursor: pointer;
  }
  .wv-copy:hover { background: var(--surface); }

  .interest-block { margin-top: 18px; text-align: left; }
  .link-btn {
    background: transparent; border: 0; padding: 0;
    font: inherit; font-size: 13px; color: var(--mute);
    cursor: pointer;
    text-align: center; width: 100%;
  }
  .link-btn .link { color: var(--accent-text); font-weight: 500; }
  .link-btn:hover .link { text-decoration: underline; }

  .interest-form {
    background: var(--surface);
    border: 1px solid var(--rule);
    border-radius: 10px;
    padding: 16px;
    display: flex; flex-direction: column; gap: 10px;
  }
  .interest-form h3 { font-size: 14px; font-weight: 500; margin: 0; letter-spacing: -0.01em; }
  .interest-help { font-size: 12px; color: var(--mute); margin: 0 0 4px; line-height: 1.5; }
  .interest-form label { display: flex; flex-direction: column; gap: 4px; }
  .interest-form .lbl { font-size: 11.5px; color: var(--mute); font-weight: 500; }
  .interest-form .opt { color: var(--mute-2); font-weight: 400; }
  .interest-form input {
    font: inherit; font-size: 13px;
    height: 34px; padding: 0 10px;
    border: 1px solid var(--rule);
    border-radius: 7px;
    background: var(--card);
    color: var(--ink);
    outline: none;
  }
  .interest-form input:focus {
    border-color: var(--accent);
    box-shadow: 0 0 0 3px var(--accent-tint);
  }
  .ifail { color: var(--danger-text); font-size: 12px; margin: 0; }

  .iactions { display: flex; justify-content: flex-end; gap: 6px; margin-top: 4px; }
  .iactions button {
    font: inherit; font-size: 12.5px; font-weight: 500;
    height: 30px; padding: 0 12px;
    border-radius: 6px; cursor: pointer;
  }
  .iactions .ghost { background: transparent; border: 1px solid var(--rule); color: var(--ink-2); }
  .iactions .primary { background: var(--ink); color: white; border: 1px solid var(--ink); }
  .iactions .primary:disabled { opacity: .55; cursor: not-allowed; }

  .interest-done {
    background: var(--positive-tint);
    border: 1px solid var(--positive-tint);
    border-radius: 10px; padding: 14px 16px;
    text-align: left;
  }
  .interest-done h3 { font-size: 14px; font-weight: 500; margin: 0 0 4px; color: var(--positive-text); }
  .interest-done p { font-size: 12.5px; color: var(--ink-2); margin: 0; line-height: 1.5; }

  .footnote { font-size: 11.5px; color: var(--mute-2); margin: 18px 0 0; }
  .privacy { font-size: 11px; line-height: 1.5; color: var(--mute-2); margin: 10px 0 0; text-wrap: pretty; }
  .privacy a { color: var(--accent-text); white-space: nowrap; }
</style>
