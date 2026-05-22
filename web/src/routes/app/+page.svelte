<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { api, STATUSES } from '$lib/api.js';

  let me = $state(null);
  let apps = $state([]);
  let view = $state('list');
  let loading = $state(true);

  // New-application form state.
  let form = $state({
    company: '',
    role: '',
    status: 'applied',
    jd_url: '',
    cv_variant: ''
  });

  const byStatus = $derived.by(() => {
    const grouped = Object.fromEntries(STATUSES.map((s) => [s, []]));
    for (const a of apps) (grouped[a.status] ||= []).push(a);
    return grouped;
  });

  onMount(refresh);

  async function refresh() {
    try {
      me = await api('/api/me');
      apps = await api('/api/applications');
    } catch (e) {
      // 401 from api() already redirects to /
      if (e.message !== 'unauthorized') console.error(e);
    } finally {
      loading = false;
    }
  }

  async function addApp(event) {
    event.preventDefault();
    const payload = { ...form };
    for (const k of ['jd_url', 'cv_variant']) {
      if (payload[k] === '') delete payload[k];
    }
    await api('/api/applications', { method: 'POST', body: JSON.stringify(payload) });
    form = { company: '', role: '', status: 'applied', jd_url: '', cv_variant: '' };
    await refresh();
  }

  async function setStatus(id, status) {
    await api(`/api/applications/${id}`, {
      method: 'PATCH',
      body: JSON.stringify({ status })
    });
    await refresh();
  }

  async function remove(id) {
    if (!confirm('Delete this application?')) return;
    await api(`/api/applications/${id}`, { method: 'DELETE' });
    await refresh();
  }

  async function signOut() {
    await api('/api/auth/logout', { method: 'POST' });
    goto('/', { replaceState: true });
  }

  function fmtDate(d) {
    return d ? new Date(d).toLocaleDateString() : '—';
  }
</script>

<svelte:head>
  <title>Pursuit — applications</title>
</svelte:head>

<header class="app-header">
  <h1>Pursuit</h1>
  <nav>
    <button class:active={view === 'list'} onclick={() => (view = 'list')}>List</button>
    <button class:active={view === 'kanban'} onclick={() => (view = 'kanban')}>Board</button>
    <span class="who">{me?.email ?? ''}</span>
    <button class="link" onclick={signOut}>Sign out</button>
  </nav>
</header>

<section class="app-body">
  <form class="row" onsubmit={addApp}>
    <input bind:value={form.company} placeholder="Company" required />
    <input bind:value={form.role} placeholder="Role" required />
    <select bind:value={form.status}>
      {#each STATUSES as s}
        <option value={s}>{s}</option>
      {/each}
    </select>
    <input bind:value={form.jd_url} placeholder="JD URL (optional)" />
    <input bind:value={form.cv_variant} placeholder="CV variant (optional)" />
    <button type="submit">Add</button>
  </form>

  {#if loading}
    <p class="empty">Loading…</p>
  {:else if view === 'list'}
    <div class="view">
      <table>
        <thead>
          <tr>
            <th>Company</th><th>Role</th><th>Status</th><th>Applied</th><th>CV</th><th></th>
          </tr>
        </thead>
        <tbody>
          {#if apps.length === 0}
            <tr><td colspan="6" class="empty">No applications yet. Add your first one above.</td></tr>
          {:else}
            {#each apps as a (a.id)}
              <tr>
                <td>{a.company}</td>
                <td>{a.role}</td>
                <td>
                  <select class="status-select" value={a.status} onchange={(e) => setStatus(a.id, e.currentTarget.value)}>
                    {#each STATUSES as s}
                      <option value={s}>{s}</option>
                    {/each}
                  </select>
                </td>
                <td>{fmtDate(a.applied_at)}</td>
                <td>{a.cv_variant ?? ''}</td>
                <td><button class="link" onclick={() => remove(a.id)}>delete</button></td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
  {:else}
    <div class="view">
      <div class="kanban">
        {#each STATUSES as s}
          <div class="kcol">
            <h3>{s} <span class="count">{byStatus[s].length}</span></h3>
            {#each byStatus[s] as a (a.id)}
              <div class="kcard">
                <strong>{a.company}</strong><br />
                <span>{a.role}</span>
              </div>
            {/each}
          </div>
        {/each}
      </div>
    </div>
  {/if}
</section>
