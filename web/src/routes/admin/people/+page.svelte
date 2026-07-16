<script>
  import { onMount } from 'svelte';
  import { api } from '$lib/api.js';

  let invites = $state([]);
  let interest = $state([]);
  let users = $state([]);
  let granting = $state(0); // user id currently being granted
  let loading = $state(true);
  let email = $state('');
  let note = $state('');
  let saving = $state(false);
  let error = $state('');
  let promoting = $state(''); // email currently being promoted

  onMount(refresh);
  async function refresh() {
    try {
      const [inv, intr, usr] = await Promise.all([
        api('/api/admin/invites'),
        api('/api/admin/beta-interest'),
        api('/api/admin/users').catch(() => [])
      ]);
      invites = inv;
      interest = intr;
      users = usr;
    } catch (e) {
      if (e.message !== 'unauthorized') console.error(e);
    } finally {
      loading = false;
    }
  }

  async function inviteFromInterest(e) {
    if (!confirm(`Invite ${e}? They'll be able to sign in with this Gmail.`)) return;
    promoting = e;
    try {
      await api(`/api/admin/beta-interest/${encodeURIComponent(e)}/invite`, { method: 'POST' });
      await refresh();
    } catch (err) {
      alert(err.message);
    } finally {
      promoting = '';
    }
  }

  const pendingInterest = $derived(interest.filter(i => !i.invited_at));
  const invitedInterest = $derived(interest.filter(i => i.invited_at));

  async function add(e) {
    e.preventDefault();
    error = '';
    saving = true;
    try {
      await api('/api/admin/invites', { method: 'POST', body: JSON.stringify({ email: email.trim(), note: note.trim() }) });
      email = ''; note = '';
      await refresh();
    } catch (e) {
      error = e.message;
    } finally {
      saving = false;
    }
  }

  async function remove(emailToRemove) {
    if (!confirm(`Remove invite for ${emailToRemove}?`)) return;
    await api(`/api/admin/invites/${encodeURIComponent(emailToRemove)}`, { method: 'DELETE' });
    await refresh();
  }

  function fmtDate(d) {
    return new Date(d).toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: 'numeric' });
  }

  // "3h ago" / "2d ago" / "just now" — for last-seen recency at a glance.
  function fmtAgo(d) {
    if (!d) return 'never';
    const secs = Math.max(0, (Date.now() - new Date(d).getTime()) / 1000);
    if (secs < 90) return 'just now';
    const mins = secs / 60;
    if (mins < 60) return `${Math.round(mins)}m ago`;
    const hrs = mins / 60;
    if (hrs < 24) return `${Math.round(hrs)}h ago`;
    const days = hrs / 24;
    if (days < 30) return `${Math.round(days)}d ago`;
    return fmtDate(d);
  }

  // Emails that have actually signed in (activated their account), lowercased
  // for case-insensitive matching against the invite list.
  const activatedEmails = $derived(new Set(users.map(u => u.email.toLowerCase())));
  function hasSignedIn(email) {
    return activatedEmails.has(email.toLowerCase());
  }
  // Pilot funnel: of everyone invited, how many have signed in at least once.
  const invitedCount = $derived(invites.length);
  const activatedInvited = $derived(invites.filter(i => hasSignedIn(i.email)).length);

  // Grant a user more AI interview-prep generations.
  async function grantPrep(u, add) {
    granting = u.id;
    try {
      await api(`/api/admin/users/${u.id}/prep-credits`, { method: 'POST', body: JSON.stringify({ add }) });
      await refresh();
    } catch (e) {
      alert(e.message);
    } finally {
      granting = 0;
    }
  }

  // Demo-data seeding — admin-only, scoped to the calling admin's own
  // applications. Tagged via [demo] prefix on notes for clean teardown.
  let seedBusy = $state(false);
  let seedMsg = $state('');
  let seedErr = $state('');

  async function seedDemo() {
    if (!confirm('Add 15 demo applications to YOUR account? (Tagged so you can clear them with one click.)')) return;
    seedBusy = true; seedMsg = ''; seedErr = '';
    try {
      const r = await api('/api/admin/demo-seed', { method: 'POST' });
      seedMsg = `Seeded ${r.inserted} demo applications. Reload the dashboard to see them.`;
    } catch (e) {
      seedErr = e.message;
    } finally {
      seedBusy = false;
    }
  }
  async function clearDemo() {
    if (!confirm('Remove all demo applications from your account? (Will NOT touch real applications.)')) return;
    seedBusy = true; seedMsg = ''; seedErr = '';
    try {
      const r = await api('/api/admin/demo-seed', { method: 'DELETE' });
      seedMsg = `Removed ${r.deleted} demo applications.`;
    } catch (e) {
      seedErr = e.message;
    } finally {
      seedBusy = false;
    }
  }
</script>

<header class="page">
  <div>
    <div class="kicker">Closed beta</div>
    <h1>People</h1>
    <p class="lede">Emails on this list can sign in. Adding someone here is the only thing you need to do — once they hit the homepage and click Continue with Google, they're in.</p>
  </div>
</header>

<form class="invite-form" onsubmit={add}>
  <div class="row">
    <input type="email" placeholder="friend@gmail.com" bind:value={email} required />
    <input type="text" placeholder="Note (optional, e.g. 'Sarah from work')" bind:value={note} />
    <button type="submit" class="btn btn-primary" disabled={saving}>{saving ? 'Adding…' : 'Add invite'}</button>
  </div>
  {#if error}<p class="error">{error}</p>{/if}
</form>

{#if !loading && invites.length > 0}
  <p class="funnel">
    <b>{activatedInvited}</b> of <b>{invitedCount}</b> invited have signed in.
    {#if activatedInvited < invitedCount}<span class="funnel-note">The rest haven't activated yet — tagged below.</span>{/if}
  </p>
{/if}

<section class="list">
  <header class="lh">
    <span>Email</span>
    <span>Note</span>
    <span>Invited</span>
    <span></span>
  </header>
  {#if loading}
    <p style="color: var(--mute); padding: 1rem">Loading…</p>
  {:else if invites.length === 0}
    <p style="color: var(--mute); padding: 1rem">Nobody on the list yet.</p>
  {:else}
    {#each invites as i (i.email)}
      {@const signedIn = hasSignedIn(i.email)}
      <div class="lr">
        <span class="email">
          {i.email}
          {#if signedIn}<span class="status on" title="Signed in at least once">signed in</span>
          {:else}<span class="status off" title="Invited but hasn't signed in yet">not signed in</span>{/if}
        </span>
        <span class="note">{i.note || '—'}</span>
        <span class="when">{fmtDate(i.invited_at)}{i.invited_by_email ? ` · by ${i.invited_by_email}` : ''}</span>
        <button class="x" onclick={() => remove(i.email)} title="Remove">×</button>
      </div>
    {/each}
  {/if}
</section>

<section class="interest">
  <header class="dh">
    <h2>Access requests <span class="count">{pendingInterest.length}</span></h2>
    <p>People who submitted the form on the homepage. Click <b>Invite</b> to promote them to the invite list above — once invited, they can sign in with that Gmail.</p>
  </header>
  {#if pendingInterest.length === 0 && invitedInterest.length === 0}
    <p class="empty">No access requests yet. Share <code>https://pursuit-playbook.com</code> with people you want to invite — they can drop their email from the homepage.</p>
  {:else}
    {#if pendingInterest.length > 0}
      <ul class="ir-list">
        {#each pendingInterest as i (i.email)}
          <li class="ir">
            <div class="ir-main">
              <div class="ir-email">{i.email}</div>
              {#if i.note}<div class="ir-note">{i.note}</div>{/if}
              <div class="ir-meta">requested {fmtDate(i.created_at)}{i.source ? ` · src: ${i.source}` : ''}</div>
            </div>
            <button class="btn btn-primary" onclick={() => inviteFromInterest(i.email)} disabled={promoting === i.email}>
              {promoting === i.email ? 'Inviting…' : 'Invite'}
            </button>
          </li>
        {/each}
      </ul>
    {/if}
    {#if invitedInterest.length > 0}
      <details class="invited-fold">
        <summary>Already invited from this list ({invitedInterest.length})</summary>
        <ul class="ir-list muted">
          {#each invitedInterest as i (i.email)}
            <li class="ir">
              <div class="ir-main">
                <div class="ir-email">{i.email}</div>
                {#if i.note}<div class="ir-note">{i.note}</div>{/if}
                <div class="ir-meta">invited {fmtDate(i.invited_at)}</div>
              </div>
            </li>
          {/each}
        </ul>
      </details>
    {/if}
  {/if}
</section>

<section class="users">
  <header class="dh">
    <h2>Users <span class="count">{users.length}</span></h2>
    <p>Everyone who's signed in, most-recently-active first — your pilot contact list. <b>Used</b> = real applications they've created (demo rows excluded), plus interviews and dossiers. <b>AI prep</b> shows interview-prep generations used vs the cap (new users start at 10); grant more when someone runs out — each generation costs you Claude credits.</p>
  </header>
  {#if users.length === 0}
    <p class="empty">No users yet.</p>
  {:else}
    <ul class="ulist">
      {#each users as u (u.id)}
        {@const capped = u.prep_credits_used >= u.prep_credits_limit}
        {@const tried = u.app_count > 0 || u.interview_count > 0 || u.dossier_count > 0}
        <li>
          <div class="ux">
            <span class="uemail">{u.email}{#if u.is_admin}<span class="utag">admin</span>{/if}{#if !tried}<span class="utag idle" title="Signed in but hasn't created anything yet">signed in, not tried</span>{/if}</span>
            <span class="umeta">
              Last seen <b>{fmtAgo(u.last_login_at)}</b> · joined {fmtDate(u.created_at)}{#if !u.onboarded_at} · not onboarded{/if}
            </span>
            <span class="umeta usage">
              {u.app_count} app{u.app_count === 1 ? '' : 's'} · {u.interview_count} interview{u.interview_count === 1 ? '' : 's'} · {u.dossier_count} dossier{u.dossier_count === 1 ? '' : 's'} · AI prep <b class:warn={capped}>{u.prep_credits_used} / {u.prep_credits_limit}</b>
            </span>
          </div>
          <div class="uactions">
            <button class="btn sm" onclick={() => grantPrep(u, 5)} disabled={granting === u.id}>+5</button>
            <button class="btn sm" onclick={() => grantPrep(u, 10)} disabled={granting === u.id}>+10</button>
          </div>
        </li>
      {/each}
    </ul>
  {/if}
</section>

<section class="demo">
  <header class="dh">
    <h2>Demo data</h2>
    <p>Fill your own account with 15 realistic applications spanning every status — for screenshots, demos, or just exploring the dashboard with something to look at. Only touches rows tagged <code>[demo]</code>; your real applications are never affected.</p>
  </header>
  <div class="demo-actions">
    <button class="btn btn-primary" onclick={seedDemo} disabled={seedBusy}>{seedBusy ? 'Working…' : 'Seed demo data'}</button>
    <button class="btn" onclick={clearDemo} disabled={seedBusy}>Clear demo data</button>
  </div>
  {#if seedMsg}<p class="seed-ok">{seedMsg}</p>{/if}
  {#if seedErr}<p class="error">{seedErr}</p>{/if}
</section>

<style>
  .page { margin-bottom: 24px; }
  .kicker { font-size: 11px; font-weight: 500; color: var(--mute); letter-spacing: .04em; text-transform: uppercase; margin-bottom: 6px; }
  h1 { font-size: 24px; font-weight: 500; letter-spacing: -0.02em; margin: 0 0 8px; }
  .lede { color: var(--mute); max-width: 60ch; margin: 0; }

  .invite-form { background: var(--card); border: 1px solid var(--rule); border-radius: 10px; padding: 16px; margin-bottom: 24px; }
  .row { display: grid; grid-template-columns: 1fr 1fr auto; gap: 10px; }
  .row input { padding: 6px 10px; border: 1px solid var(--rule); border-radius: 6px; background: var(--surface); font: inherit; font-size: 13.5px; color: var(--ink); outline: none; }
  .row input:focus { border-color: var(--accent); }
  .error { color: var(--danger-text); font-size: 12px; margin: 8px 0 0; }

  .funnel { font-size: 13px; color: var(--ink-2); margin: 0 0 10px; }
  .funnel b { color: var(--ink); font-weight: 600; font-variant-numeric: tabular-nums; }
  .funnel-note { color: var(--mute); margin-left: 4px; }

  .list { border: 1px solid var(--rule); border-radius: 10px; background: var(--card); overflow: hidden; }
  .lh, .lr { display: grid; grid-template-columns: 1.6fr 2fr 1.6fr 32px; gap: 12px; align-items: center; padding: 10px 16px; }
  .lh { font-size: 12px; color: var(--mute); background: var(--surface); border-bottom: 1px solid var(--rule); font-weight: 500; }
  .lr { border-top: 1px solid var(--rule); font-size: 13.5px; }
  .lr:first-of-type { border-top: none; }
  .email { font-weight: 500; display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
  .status { font-size: 10.5px; font-weight: 600; letter-spacing: 0.03em; text-transform: uppercase; border-radius: 4px; padding: 1px 6px; white-space: nowrap; }
  .status.on { color: var(--positive-text); background: var(--positive-tint, rgba(34,160,90,0.12)); }
  .status.off { color: var(--mute); background: var(--surface-2); }
  .note { color: var(--mute); }
  .when { color: var(--mute); font-size: 12px; font-variant-numeric: tabular-nums; }
  .x { background: transparent; border: 0; color: var(--mute-2); font-size: 18px; cursor: pointer; width: 24px; height: 24px; border-radius: 4px; }
  .x:hover { background: var(--surface-2); color: var(--danger-text); }

  .demo { margin-top: 32px; border: 1px solid var(--rule); border-radius: 10px; background: var(--card); padding: 18px 20px; }
  .demo .dh h2 { font-size: 16px; font-weight: 500; letter-spacing: -0.015em; margin: 0 0 4px; }
  .demo .dh p { color: var(--mute); font-size: 13px; margin: 0 0 14px; max-width: 70ch; line-height: 1.5; }
  .demo .dh code { font-family: var(--mono); font-size: 11.5px; background: var(--surface-2); padding: 1px 5px; border-radius: 3px; color: var(--ink-2); }
  .demo-actions { display: flex; gap: 8px; }
  .seed-ok { color: var(--positive-text); font-size: 12.5px; margin: 10px 0 0; }

  .interest { margin-top: 32px; border: 1px solid var(--rule); border-radius: 10px; background: var(--card); padding: 18px 20px; }
  .interest .dh h2 { font-size: 16px; font-weight: 500; letter-spacing: -0.015em; margin: 0 0 4px; display: flex; align-items: center; gap: 8px; }
  .interest .dh .count { font-family: var(--mono); font-size: 11px; background: var(--accent-tint); color: var(--accent-text); padding: 2px 7px; border-radius: 4px; }
  .interest .dh p { color: var(--mute); font-size: 13px; margin: 0 0 14px; max-width: 70ch; line-height: 1.5; }
  .interest .empty { color: var(--mute); font-size: 13px; margin: 0; padding: 14px 0; }
  .interest .empty code { font-family: var(--mono); font-size: 12px; background: var(--surface-2); padding: 1px 5px; border-radius: 3px; color: var(--ink-2); }

  .ir-list { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 8px; }
  .ir { display: flex; align-items: center; gap: 12px; padding: 12px 14px; background: var(--surface); border: 1px solid var(--rule); border-radius: 8px; }
  .ir-main { flex: 1; min-width: 0; }
  .ir-email { font-weight: 500; font-size: 13.5px; color: var(--ink); }
  .ir-note { font-size: 12.5px; color: var(--ink-2); margin-top: 2px; line-height: 1.4; }
  .ir-meta { font-size: 11.5px; color: var(--mute); margin-top: 4px; }
  .ir-list.muted .ir { background: var(--surface-2); opacity: .75; }
  .invited-fold { margin-top: 12px; }
  .invited-fold summary { cursor: pointer; font-size: 12.5px; color: var(--mute); padding: 6px 0; }
  .invited-fold summary:hover { color: var(--ink-2); }
  .invited-fold[open] summary { margin-bottom: 8px; }

  /* Users + AI-prep credits */
  .users { margin-top: 36px; }
  .users .ulist { list-style: none; margin: 14px 0 0; padding: 0; display: flex; flex-direction: column; }
  .users .ulist li { display: flex; align-items: center; justify-content: space-between; gap: 12px; padding: 12px 2px; border-top: 1px solid var(--rule); }
  .users .ulist li:first-child { border-top: 0; }
  .users .ux { min-width: 0; }
  .users .uemail { font-size: 13.5px; font-weight: 500; color: var(--ink); display: flex; align-items: center; gap: 8px; }
  .users .utag { font-size: 10.5px; font-weight: 600; letter-spacing: 0.04em; text-transform: uppercase; color: var(--accent-text); background: var(--accent-tint); border-radius: 4px; padding: 1px 6px; }
  .users .utag.idle { color: var(--mute); background: var(--surface-2); }
  .users .umeta { font-size: 12.5px; color: var(--mute); margin-top: 2px; }
  .users .umeta.usage { font-variant-numeric: tabular-nums; }
  .users .umeta b { color: var(--ink-2); font-weight: 600; font-variant-numeric: tabular-nums; }
  .users .umeta b.warn { color: var(--warm-text); }
  .users .uactions { display: flex; gap: 6px; flex-shrink: 0; }
  .users .btn.sm { height: 30px; padding: 0 12px; font-size: 12.5px; }
  .users .empty { color: var(--mute); font-size: 13px; margin-top: 10px; }
</style>
