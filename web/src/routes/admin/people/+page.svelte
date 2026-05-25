<script>
  import { onMount } from 'svelte';
  import { api } from '$lib/api.js';

  let invites = $state([]);
  let loading = $state(true);
  let email = $state('');
  let note = $state('');
  let saving = $state(false);
  let error = $state('');

  onMount(refresh);
  async function refresh() {
    try {
      invites = await api('/api/admin/invites');
    } catch (e) {
      if (e.message !== 'unauthorized') console.error(e);
    } finally {
      loading = false;
    }
  }

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
</style>
