<script>
  import { onMount } from 'svelte';
  import { api } from '$lib/api.js';

  let invitees = $state([]);
  let pendingCount = $state(0);
  let loading = $state(true);
  let error = $state('');
  let filter = $state('all'); // all | invited | signed_in | activated | active | dormant

  onMount(async () => {
    try {
      const r = await api('/api/admin/invite-funnel');
      invitees = r.invitees || [];
      pendingCount = r.pending_count || 0;
    } catch (e) {
      if (e.message !== 'unauthorized') error = e.message;
    } finally {
      loading = false;
    }
  });

  function fmtDate(d) {
    if (!d) return '—';
    return new Date(d).toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: 'numeric' });
  }
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

  // Funnel stages, in linear order. Each invitee sits at exactly one stage.
  const STAGES = {
    invited:   { label: 'Invited',     cls: 'st-invited',   blurb: "hasn't signed in" },
    signed_in: { label: 'Signed in',   cls: 'st-signedin',  blurb: 'no activity yet' },
    activated: { label: 'Activated',   cls: 'st-activated', blurb: 'created something' },
    active:    { label: 'Active',      cls: 'st-active',    blurb: 'seen this week' },
    dormant:   { label: 'Dormant',    cls: 'st-dormant',   blurb: 'quiet 21d+' }
  };

  const count = (s) => invitees.filter(i => i.stage === s).length;
  const signedInTotal = $derived(invitees.filter(i => i.stage !== 'invited').length);
  const activatedTotal = $derived(invitees.filter(i => ['activated', 'active', 'dormant'].includes(i.stage)).length);
  const activeTotal = $derived(count('active'));

  // Sort: most-recently-active first within the view; never-seen sink to the bottom.
  const shown = $derived(
    invitees
      .filter(i => filter === 'all' || i.stage === filter)
      .slice()
      .sort((a, b) => {
        const ta = a.last_activity_at ? new Date(a.last_activity_at).getTime() : 0;
        const tb = b.last_activity_at ? new Date(b.last_activity_at).getTime() : 0;
        return tb - ta;
      })
  );

  // People who signed in but never created anything — the onboarding-leak bucket.
  const stalled = $derived(invitees.filter(i => i.stage === 'signed_in').length);
  // Invited a while ago, still no-show.
  const noShow = $derived(invitees.filter(i =>
    i.stage === 'invited' && (Date.now() - new Date(i.invited_at).getTime()) > 7 * 864e5).length);
</script>

<header class="page">
  <div class="kicker">Closed beta</div>
  <h1>Invite funnel</h1>
  <p class="lede">Every invited email, matched to the account it became — so you can see who actually showed up and what they did. Matched on email; demo rows excluded from counts.</p>
</header>

{#if loading}
  <p class="muted">Loading…</p>
{:else if error}
  <p class="error">{error}</p>
{:else}
  <!-- Funnel summary -->
  <section class="funnel">
    <div class="step"><b>{invitees.length}</b><span>Invited</span></div>
    <div class="arrow">→</div>
    <div class="step"><b>{signedInTotal}</b><span>Signed in</span></div>
    <div class="arrow">→</div>
    <div class="step"><b>{activatedTotal}</b><span>Activated</span></div>
    <div class="arrow">→</div>
    <div class="step accent"><b>{activeTotal}</b><span>Active this week</span></div>
  </section>

  {#if pendingCount > 0 || stalled > 0 || noShow > 0}
    <p class="nudges">
      {#if pendingCount > 0}<span><b>{pendingCount}</b> access request{pendingCount === 1 ? '' : 's'} waiting — <a href="/admin/people">review on People</a>.</span>{/if}
      {#if noShow > 0}<span><b>{noShow}</b> invited 7d+ ago, never signed in.</span>{/if}
      {#if stalled > 0}<span><b>{stalled}</b> signed in but created nothing — onboarding leak.</span>{/if}
    </p>
  {/if}

  <!-- Stage filter -->
  <div class="filters">
    <button class:on={filter === 'all'} onclick={() => filter = 'all'}>All <span>{invitees.length}</span></button>
    {#each Object.entries(STAGES) as [key, s]}
      <button class:on={filter === key} onclick={() => filter = key}>{s.label} <span>{count(key)}</span></button>
    {/each}
  </div>

  {#if shown.length === 0}
    <p class="muted">{invitees.length === 0 ? 'Nobody invited yet.' : 'None in this stage.'}</p>
  {:else}
    <ul class="list">
      {#each shown as i (i.email)}
        {@const s = STAGES[i.stage]}
        <li>
          <div class="main">
            <div class="top">
              <span class="email">{i.email}</span>
              <span class="badge {s.cls}" title={s.blurb}>{s.label}</span>
            </div>
            <div class="meta">
              Invited {fmtDate(i.invited_at)}{i.invited_by_email ? ` · by ${i.invited_by_email}` : ''}
              {#if i.signed_in_at} · joined {fmtDate(i.signed_in_at)}{:else} · {s.blurb}{/if}
              {#if i.signed_in_at && !i.onboarded_at} · not onboarded{/if}
            </div>
            {#if i.signed_in_at}
              <div class="meta usage">
                {i.app_count} app{i.app_count === 1 ? '' : 's'} ·
                {i.interview_count} interview{i.interview_count === 1 ? '' : 's'} ·
                {i.dossier_count} dossier{i.dossier_count === 1 ? '' : 's'} ·
                {i.event_count} event{i.event_count === 1 ? '' : 's'}
              </div>
            {/if}
          </div>
          <div class="seen">
            {#if i.signed_in_at}
              <span class="ago">{fmtAgo(i.last_activity_at)}</span>
              <span class="ago-label">last active</span>
            {:else}
              <span class="ago muted">—</span>
            {/if}
          </div>
        </li>
      {/each}
    </ul>
  {/if}
{/if}

<style>
  .page { margin-bottom: 20px; }
  .kicker { font-size: 11px; font-weight: 500; color: var(--mute); letter-spacing: .04em; text-transform: uppercase; margin-bottom: 6px; }
  h1 { font-size: 24px; font-weight: 500; letter-spacing: -0.02em; margin: 0 0 8px; }
  .lede { color: var(--mute); max-width: 64ch; margin: 0; line-height: 1.5; }
  .muted { color: var(--mute); padding: 0.5rem 0; }
  .error { color: var(--danger-text); }

  .funnel {
    display: flex; align-items: stretch; gap: 4px;
    background: var(--card); border: 1px solid var(--rule); border-radius: 10px;
    padding: 16px; margin-bottom: 14px;
  }
  .funnel .step { flex: 1; display: flex; flex-direction: column; align-items: center; gap: 3px; padding: 4px 0; }
  .funnel .step b { font-size: 26px; font-weight: 600; letter-spacing: -0.02em; color: var(--ink); font-variant-numeric: tabular-nums; }
  .funnel .step span { font-size: 11.5px; color: var(--mute); text-align: center; }
  .funnel .step.accent b { color: var(--accent-text); }
  .funnel .arrow { display: flex; align-items: center; color: var(--mute-2); font-size: 14px; }

  .nudges { display: flex; flex-wrap: wrap; gap: 6px 16px; font-size: 12.5px; color: var(--ink-2); margin: 0 0 18px; }
  .nudges b { color: var(--ink); font-weight: 600; font-variant-numeric: tabular-nums; }
  .nudges a { color: var(--accent-text); }

  .filters { display: flex; flex-wrap: wrap; gap: 6px; margin-bottom: 14px; }
  .filters button {
    font: inherit; font-size: 12.5px; cursor: pointer;
    background: var(--surface); border: 1px solid var(--rule); color: var(--ink-2);
    padding: 4px 10px; border-radius: 999px; display: inline-flex; align-items: center; gap: 6px;
  }
  .filters button:hover { border-color: var(--mute-2); }
  .filters button.on { background: var(--ink); color: var(--surface); border-color: var(--ink); }
  .filters button span { font-variant-numeric: tabular-nums; opacity: .7; }

  .list { list-style: none; margin: 0; padding: 0; border: 1px solid var(--rule); border-radius: 10px; background: var(--card); overflow: hidden; }
  .list li { display: flex; align-items: center; justify-content: space-between; gap: 12px; padding: 13px 16px; border-top: 1px solid var(--rule); }
  .list li:first-child { border-top: none; }
  .main { min-width: 0; }
  .top { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
  .email { font-size: 13.5px; font-weight: 500; color: var(--ink); }
  .meta { font-size: 12px; color: var(--mute); margin-top: 3px; }
  .meta.usage { font-variant-numeric: tabular-nums; color: var(--ink-2); }
  .seen { text-align: right; flex-shrink: 0; }
  .seen .ago { display: block; font-size: 12.5px; font-weight: 600; color: var(--ink-2); font-variant-numeric: tabular-nums; }
  .seen .ago.muted { color: var(--mute-2); font-weight: 400; }
  .seen .ago-label { font-size: 10.5px; color: var(--mute); }

  .badge { font-size: 10.5px; font-weight: 600; letter-spacing: 0.03em; text-transform: uppercase; border-radius: 4px; padding: 1px 6px; white-space: nowrap; }
  .st-invited   { color: var(--mute); background: var(--surface-2); }
  .st-signedin  { color: var(--warm-text, #9a6b00); background: var(--warm-tint, rgba(200,140,0,0.12)); }
  .st-activated { color: var(--accent-text); background: var(--accent-tint); }
  .st-active    { color: var(--positive-text); background: var(--positive-tint, rgba(34,160,90,0.12)); }
  .st-dormant   { color: var(--mute); background: var(--surface-2); border: 1px solid var(--rule); }
</style>
