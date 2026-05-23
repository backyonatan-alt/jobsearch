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
</style>
