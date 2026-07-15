<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api.js';
  import { isPreview, mockApi } from '$lib/preview-mode.js';
  import { logEvent } from '$lib/analytics.js';
  import { toDisplayApp, daysSince, fmtRelativeDate } from '$lib/app-helpers.js';
  import CompanyLogo from '$lib/CompanyLogo.svelte';
  import ImportApplications from '$lib/ImportApplications.svelte';

  const call = isPreview() ? mockApi : api;

  let apps = $state([]);
  let loading = $state(true);
  let showImport = $state(false);

  onMount(() => {
    refresh();
    const h = () => refresh();
    window.addEventListener('pursuit:refresh', h);
    return () => window.removeEventListener('pursuit:refresh', h);
  });
  async function refresh() {
    try { apps = (await call('/api/applications')).map(toDisplayApp); }
    catch (e) { if (e.message !== 'unauthorized') console.error(e); }
    finally { loading = false; }
  }

  const EXITS = ['rejected', 'withdrawn', 'closed'];
  const quietDays = (a) => daysSince(a.raw.applied_at) ?? 0;

  function rowMeta(a) {
    const bits = [a.role];
    if (a.status === 'wishlist') bits.push(`saved ${fmtRelativeDate(a.raw.created_at)}`);
    else if (EXITS.includes(a.status)) {
      const label = a.status === 'closed' ? 'position closed' : a.status;
      bits.push(`${label} · ${daysSince(a.raw.updated_at || a.raw.applied_at) ?? 0}d`);
    } else {
      bits.push(`applied ${a.appliedRel}${a.source && a.source !== '—' ? ` via ${a.source}` : ''}`);
    }
    return bits.join(' · ');
  }

  const groups = $derived.by(() => {
    const g = (name, labelColor, statuses, opts = {}) => {
      const rows = apps.filter(a => statuses.includes(a.status));
      return { name, labelColor, rows, ...opts };
    };
    const applied = g('Applied', '#8a9099', ['applied']);
    const nQuiet = applied.rows.filter(a => a.stale).length;
    if (nQuiet) { applied.note = `— ${nQuiet} quiet too long`; applied.noteColor = '#b3372a'; }
    const out = [
      g('Interview', '#e0641f', ['interview']),
      g('Offer', '#1d7a4f', ['offer']),
      g('Screen', '#0e9f6e', ['screen']),
      applied,
      g('Wishlist', '#8a9099', ['wishlist']),
      g('No longer in play', '#b8bdc4', EXITS, {
        note: '— kept forever, reopen if the req comes back', noteColor: '#b8bdc4', exit: true
      })
    ];
    return out.filter(x => x.rows.length > 0 || x.exit);
  });

  const inPlay = $derived(apps.filter(a => !EXITS.includes(a.status)).length);
  const closed = $derived(apps.length - inPlay);

  function actionFor(a) {
    if (a.status === 'interview') return { label: 'Prep →', go: () => goto(`/app/${a.id}#interview-prep`) };
    if (a.status === 'offer') return { label: 'Decide →', go: () => goto(`/app/${a.id}`) };
    if (a.status === 'wishlist') return { label: 'Apply →', go: () => goto(`/app/${a.id}`) };
    if (a.status === 'applied' && a.stale) return { label: 'Follow up', go: () => goto(`/app/${a.id}`) };
    return null;
  }

  async function reopen(e, a) {
    e.stopPropagation();
    await api(`/api/applications/${a.id}`, { method: 'PATCH', body: JSON.stringify({ status: 'applied' }) });
    logEvent('archive_reopened', { app_id: Number(a.id), from: a.status });
    try { window.dispatchEvent(new CustomEvent('pursuit:refresh')); } catch {}
    refresh();
  }

  function open(id) { goto(`/app/${id}`); }
</script>

<svelte:head><title>Applications — Pursuit</title></svelte:head>

<div class="pg">
  <div class="head">
    <div>
      <h1>Applications.</h1>
      <div class="sub">{apps.length} total · {inPlay} in play · {closed} closed</div>
    </div>
    <button class="imp" onclick={() => (showImport = true)}>↓ Import from spreadsheet</button>
  </div>

  {#if loading}
    <p class="loading">Loading…</p>
  {:else if apps.length === 0}
    <div class="empty">
      <h3>No applications yet</h3>
      <p>Add your first one with "New application" above — they'll appear here grouped by stage.</p>
    </div>
  {:else}
    <div data-tour="board">
      {#each groups as g (g.name)}
        {#if g.rows.length}
          <div class="grp">
            <div class="grp-hd">
              <span class="grp-t" style="color:{g.labelColor}">{g.name} · {g.rows.length}</span>
              {#if g.note}<span class="grp-n" style="color:{g.noteColor}">{g.note}</span>{/if}
            </div>
            <div class="box" class:dim={g.exit}>
              {#each g.rows as a (a.id)}
                {@const act = g.exit ? null : actionFor(a)}
                <div class="row" class:ex={g.exit} onclick={() => open(a.id)} role="button" tabindex="0">
                  <CompanyLogo app={a} size={24} gray={g.exit} />
                  <span class="tx" class:extx={g.exit}>
                    <strong>{a.co}</strong> <span class="meta">· {rowMeta(a)}</span>
                    {#if !g.exit && a.stale && a.status === 'applied'}<strong class="hot">quiet {quietDays(a)}d</strong>{/if}
                  </span>
                  {#if g.exit}
                    <button class="lnk reopen" onclick={(e) => reopen(e, a)}>Reopen</button>
                  {:else if act}
                    <button class="lnk" onclick={(e) => { e.stopPropagation(); act.go(); }}>{act.label}</button>
                  {/if}
                </div>
              {/each}
            </div>
          </div>
        {/if}
      {/each}
    </div>
  {/if}
</div>

<ImportApplications bind:open={showImport} onImported={refresh} />

<style>
  .pg { max-width: 900px; margin: 0 auto; padding: 36px 32px 80px; width: 100%; box-sizing: border-box; }
  .head { display: flex; align-items: baseline; justify-content: space-between; margin-bottom: 30px; }
  .head h1 { font-size: 30px; font-weight: 700; letter-spacing: -0.02em; margin: 0 0 6px; }
  .sub { font-size: 13.5px; color: #6f7680; }
  .imp { background: none; border: 0; font-size: 13px; color: #2463eb; cursor: pointer; font-family: inherit; padding: 0; }
  .imp:hover { color: #1a4fc4; }
  .loading { color: #8a9099; font-size: 13.5px; }

  .grp { margin-bottom: 24px; }
  .grp-hd { display: flex; align-items: baseline; gap: 8px; margin-bottom: 8px; }
  .grp-t { font-size: 11px; font-weight: 600; letter-spacing: .12em; text-transform: uppercase; }
  .grp-n { font-size: 12px; }
  .box { background: #fff; border: 1px solid #eeeeea; border-radius: 10px; }
  .box.dim { background: #fbfbf9; }
  .row {
    display: flex; align-items: center; gap: 12px; padding: 10px 16px;
    border-bottom: 1px solid #f4f4f1; cursor: pointer;
  }
  .row:last-child { border-bottom: 0; }
  .row:hover { background: #fafaf8; }
  .row.ex { opacity: .8; }
  .tx { flex: 1; min-width: 0; font-size: 13.5px; color: #4b5158; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  .tx strong { color: #16181c; }
  .tx.extx, .tx.extx strong { color: #8a9099; }
  .meta { color: #8a9099; }
  .hot { color: #b3372a; }
  .lnk { background: none; border: 0; font-size: 12.5px; font-weight: 600; color: #2463eb; cursor: pointer; flex: none; font-family: inherit; }
  .lnk:hover { color: #1a4fc4; }
  .lnk.reopen { color: #8a9099; }
  .lnk.reopen:hover { color: #4b5158; }

  .empty { border: 1px dashed #e2e2de; border-radius: 12px; padding: 32px; text-align: center; background: #fff; }
  .empty h3 { margin: 0 0 .5rem; font-size: 16px; font-weight: 500; }
  .empty p { color: #8a9099; margin: 0; font-size: 13.5px; }

  @media (max-width: 900px) { .pg { padding: 24px 16px 60px; } }
</style>
