<script>
  import { onMount } from 'svelte';
  import { api } from '$lib/api.js';

  let invites = $state([]);
  let interest = $state([]);
  let loading = $state(true);
  let email = $state('');
  let note = $state('');
  let saving = $state(false);
  let error = $state('');
  let promoting = $state(''); // email currently being promoted

  onMount(refresh);
  async function refresh() {
    try {
      const [inv, intr] = await Promise.all([
        api('/api/admin/invites'),
        api('/api/admin/beta-interest')
      ]);
      invites = inv;
      interest = intr;
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
      <div class="lr">
        <span class="email">{i.email}</span>
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
    <p class="empty">No access requests yet. Share <code>https://178.105.213.124.nip.io</code> with people you want to invite — they can drop their email from the homepage.</p>
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

  .list { border: 1px solid var(--rule); border-radius: 10px; background: var(--card); overflow: hidden; }
  .lh, .lr { display: grid; grid-template-columns: 1.6fr 2fr 1.6fr 32px; gap: 12px; align-items: center; padding: 10px 16px; }
  .lh { font-size: 12px; color: var(--mute); background: var(--surface); border-bottom: 1px solid var(--rule); font-weight: 500; }
  .lr { border-top: 1px solid var(--rule); font-size: 13.5px; }
  .lr:first-of-type { border-top: none; }
  .email { font-weight: 500; }
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
</style>
