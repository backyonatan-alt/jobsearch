<script>
  import { onMount } from 'svelte';
  import { api } from '$lib/api.js';

  let total = $state(0);
  let milestones = $state([]);
  let events = $state([]);
  let loading = $state(true);
  let error = $state('');

  onMount(async () => {
    try {
      const r = await api('/api/admin/adoption');
      total = r.total_users || 0;
      milestones = r.milestones || [];
      events = r.events || [];
    } catch (e) {
      if (e.message !== 'unauthorized') error = e.message;
    } finally {
      loading = false;
    }
  });

  const pct = (n) => (total > 0 ? Math.round((n / total) * 100) : 0);

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
    return `${Math.round(days / 30)}mo ago`;
  }

  // A surface fired by nobody in the last 7 days is worth a second look.
  const stale = (e) => e.recent === 0;
</script>

<header class="page">
  <div class="kicker">Closed beta</div>
  <h1>Adoption</h1>
  <p class="lede">Of the <b>{total}</b> people who've signed in, how many reached each moment that matters — and which surfaces actually get used. Demo rows excluded from the milestones.</p>
</header>

{#if loading}
  <p class="muted">Loading…</p>
{:else if error}
  <p class="error">{error}</p>
{:else}
  <!-- Feature reach -->
  <section class="block">
    <h2>Feature reach</h2>
    <p class="sub">Distinct signed-in users who reached each milestone. The drop from one bar to the next is where you're losing people.</p>
    <ul class="reach">
      {#each milestones as m (m.key)}
        <li>
          <div class="rlabel">{m.label}</div>
          <div class="rbar">
            <div class="rfill" style="width: {Math.max(pct(m.users), m.users > 0 ? 3 : 0)}%"></div>
          </div>
          <div class="rnum"><b>{m.users}</b> <span>/ {total} · {pct(m.users)}%</span></div>
        </li>
      {/each}
    </ul>
  </section>

  <!-- Event activity -->
  <section class="block">
    <h2>Event activity</h2>
    <p class="sub">Every first-party event logged, most-used first. <b>Users</b> = distinct people who fired it; <b>7d</b> = times fired this week. A surface with 0 this week is greyed.</p>
    {#if events.length === 0}
      <p class="muted">No events logged yet.</p>
    {:else}
      <div class="etable">
        <div class="eh">
          <span>Event</span>
          <span class="n">Users</span>
          <span class="n">Total</span>
          <span class="n">7d</span>
          <span class="n">Last</span>
        </div>
        {#each events as e (e.name)}
          <div class="er" class:dim={stale(e)}>
            <span class="ename">{e.name}</span>
            <span class="n"><b>{e.users}</b></span>
            <span class="n">{e.total}</span>
            <span class="n">{e.recent || '—'}</span>
            <span class="n ago">{fmtAgo(e.last_at)}</span>
          </div>
        {/each}
      </div>
    {/if}
  </section>

  <p class="note">
    Blind spot: we don't yet log modal-opens or failed parses, so users who <i>tried</i> and bounced look identical to those who never tried. Worth adding a couple of events before the next cohort.
  </p>
{/if}

<style>
  .page { margin-bottom: 22px; }
  .kicker { font-size: 11px; font-weight: 500; color: var(--mute); letter-spacing: .04em; text-transform: uppercase; margin-bottom: 6px; }
  h1 { font-size: 24px; font-weight: 500; letter-spacing: -0.02em; margin: 0 0 8px; }
  .lede { color: var(--mute); max-width: 64ch; margin: 0; line-height: 1.5; }
  .lede b { color: var(--ink); font-weight: 600; }
  .muted { color: var(--mute); padding: 0.5rem 0; }
  .error { color: var(--danger-text); }

  .block { background: var(--card); border: 1px solid var(--rule); border-radius: 10px; padding: 18px 20px; margin-bottom: 16px; }
  .block h2 { font-size: 16px; font-weight: 500; letter-spacing: -0.015em; margin: 0 0 4px; }
  .block .sub { color: var(--mute); font-size: 12.5px; margin: 0 0 16px; max-width: 70ch; line-height: 1.5; }
  .block .sub b { color: var(--ink-2); font-weight: 600; }

  .reach { list-style: none; margin: 0; padding: 0; display: flex; flex-direction: column; gap: 11px; }
  .reach li { display: grid; grid-template-columns: 1.4fr 2.4fr auto; gap: 14px; align-items: center; }
  .rlabel { font-size: 13px; color: var(--ink); }
  .rbar { height: 12px; background: var(--surface-2); border-radius: 6px; overflow: hidden; }
  .rfill { height: 100%; background: var(--accent); border-radius: 6px; transition: width .3s ease; }
  .rnum { font-size: 12.5px; color: var(--mute); white-space: nowrap; font-variant-numeric: tabular-nums; }
  .rnum b { color: var(--ink); font-weight: 600; font-size: 14px; }

  .etable { font-size: 13px; }
  .eh, .er { display: grid; grid-template-columns: 2fr 0.8fr 0.8fr 0.7fr 1fr; gap: 10px; align-items: center; padding: 8px 0; }
  .eh { font-size: 11.5px; color: var(--mute); border-bottom: 1px solid var(--rule); font-weight: 500; }
  .er { border-bottom: 1px solid var(--rule); }
  .er:last-child { border-bottom: none; }
  .er.dim { opacity: .5; }
  .ename { font-family: var(--mono); font-size: 12px; color: var(--ink); }
  .n { text-align: right; font-variant-numeric: tabular-nums; color: var(--ink-2); }
  .n b { color: var(--ink); font-weight: 600; }
  .er .ago { color: var(--mute); font-size: 12px; }

  .note { font-size: 12px; color: var(--mute); line-height: 1.5; margin: 4px 2px 0; max-width: 70ch; }
  .note i { font-style: italic; }
</style>
