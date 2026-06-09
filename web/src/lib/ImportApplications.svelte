<script>
  // Bulk import — paste rows copied from Excel / Google Sheets, map the
  // columns, preview, and create them all at once.
  import { api } from '$lib/api.js';
  import { logEvent } from '$lib/analytics.js';

  let { open = $bindable(false), onImported } = $props();

  let step = $state('paste'); // 'paste' | 'map'
  let raw = $state('');
  let hasHeader = $state(true);
  let headers = $state([]);    // column labels
  let dataLines = $state([]);  // array of cell arrays
  let mapping = $state([]);    // field key per column ('' = ignore)
  let importing = $state(false);
  let error = $state('');
  let result = $state(null);

  const FIELDS = [
    { k: '',            lbl: '— Ignore —' },
    { k: 'company',     lbl: 'Company' },
    { k: 'role',        lbl: 'Role' },
    { k: 'status',      lbl: 'Status' },
    { k: 'source',      lbl: 'Source' },
    { k: 'location',    lbl: 'Location' },
    { k: 'salary_note', lbl: 'Salary' },
    { k: 'applied_at',  lbl: 'Applied date' },
    { k: 'notes',       lbl: 'Notes' }
  ];

  function guessField(h) {
    const s = (h || '').toLowerCase();
    if (/company|employer|organi/.test(s)) return 'company';
    if (/role|title|position|\bjob\b/.test(s)) return 'role';
    if (/status|stage/.test(s)) return 'status';
    if (/source|via|channel|referr/.test(s)) return 'source';
    if (/location|city|where|remote/.test(s)) return 'location';
    if (/salary|comp|pay/.test(s)) return 'salary_note';
    if (/date|applied|when/.test(s)) return 'applied_at';
    if (/note|comment/.test(s)) return 'notes';
    return '';
  }

  function splitLine(line, delim) {
    return line.split(delim).map(c => c.trim());
  }

  function parseTable() {
    error = '';
    const lines = raw.replace(/\r/g, '').split('\n').filter(l => l.trim() !== '');
    if (lines.length === 0) { error = 'Paste some rows first.'; return; }
    const delim = lines[0].includes('\t') ? '\t' : ',';
    const rows = lines.map(l => splitLine(l, delim));
    const width = Math.max(...rows.map(r => r.length));

    if (hasHeader) {
      headers = rows[0];
      while (headers.length < width) headers.push(`Column ${headers.length + 1}`);
      dataLines = rows.slice(1);
    } else {
      headers = Array.from({ length: width }, (_, i) => `Column ${i + 1}`);
      dataLines = rows;
    }
    mapping = headers.map(h => hasHeader ? guessField(h) : '');
    if (dataLines.length === 0) { error = 'No data rows found below the header.'; return; }
    step = 'map';
  }

  function parseDate(v) {
    if (!v) return null;
    const d = new Date(v);
    return isNaN(d.getTime()) ? null : d.toISOString();
  }

  const builtRows = $derived.by(() =>
    dataLines.map(cells => {
      const o = {};
      mapping.forEach((field, i) => {
        if (!field) return;
        const v = (cells[i] ?? '').trim();
        if (field === 'applied_at') o.applied_at = parseDate(v);
        else o[field] = v;
      });
      return o;
    })
  );
  const validRows = $derived(builtRows.filter(r => r.company && r.role));
  const skippedCount = $derived(builtRows.length - validRows.length);
  const hasCompany = $derived(mapping.includes('company'));
  const hasRole = $derived(mapping.includes('role'));

  async function doImport() {
    if (importing || !validRows.length) return;
    importing = true;
    error = '';
    try {
      const r = await api('/api/applications/import', {
        method: 'POST',
        body: JSON.stringify({ applications: validRows })
      });
      result = r;
      logEvent('bulk_import', { created: r.created ?? 0, skipped: r.skipped ?? 0 });
      onImported?.();
    } catch (e) {
      error = e.message || 'Import failed.';
    } finally {
      importing = false;
    }
  }

  function reset() {
    step = 'paste'; raw = ''; headers = []; dataLines = []; mapping = [];
    error = ''; result = null; importing = false;
  }
  function close() { reset(); open = false; }
</script>

{#if open}
  <div class="ovl" onclick={close} role="presentation">
    <div class="card" onclick={(e) => e.stopPropagation()} role="dialog" aria-modal="true" aria-label="Import applications">
      <button class="x" onclick={close} aria-label="Close">
        <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="1.6"><path d="M3 3l8 8M11 3l-8 8" stroke-linecap="round"/></svg>
      </button>

      {#if result}
        <div class="hd"><h2>Import complete</h2></div>
        <p class="done-line">
          <span class="done-ic">✓</span>
          Added <strong>{result.created}</strong> application{result.created === 1 ? '' : 's'}{result.skipped ? ` · skipped ${result.skipped} incomplete row${result.skipped === 1 ? '' : 's'}` : ''}.
        </p>
        <div class="foot">
          <span></span>
          <button class="btn btn-primary" onclick={close}>Done</button>
        </div>

      {:else if step === 'paste'}
        <div class="hd">
          <h2>Import from a spreadsheet</h2>
          <p class="sub">Copy the rows from Excel or Google Sheets (including the header row) and paste them here. We'll match the columns next.</p>
        </div>
        <textarea
          class="paste"
          bind:value={raw}
          rows="9"
          placeholder={"Company\tRole\tStatus\tSource\tApplied\nAcme\tProduct Manager\tInterview\tLinkedIn\t2026-05-12\nGlobex\tSenior PM\tApplied\tReferral\t2026-05-20"}
        ></textarea>
        <label class="chk">
          <input type="checkbox" bind:checked={hasHeader} />
          First row is a header
        </label>
        {#if error}<p class="err">{error}</p>{/if}
        <div class="foot">
          <button class="btn" onclick={close}>Cancel</button>
          <button class="btn btn-primary" onclick={parseTable} disabled={!raw.trim()}>Continue</button>
        </div>

      {:else}
        <div class="hd">
          <h2>Match your columns</h2>
          <p class="sub">Tell us what each column is. Company and Role are required — rows missing either are skipped.</p>
        </div>

        <div class="map-grid">
          {#each headers as h, i}
            <div class="map-row">
              <div class="map-col">
                <div class="map-h">{h}</div>
                <div class="map-sample">{dataLines[0]?.[i] ?? ''}</div>
              </div>
              <div class="map-arrow">→</div>
              <select bind:value={mapping[i]} class="map-sel">
                {#each FIELDS as f}<option value={f.k}>{f.lbl}</option>{/each}
              </select>
            </div>
          {/each}
        </div>

        {#if !hasCompany || !hasRole}
          <p class="warn">Map a column to <strong>Company</strong> and <strong>Role</strong> to continue.</p>
        {:else}
          <div class="preview-note">
            Ready to import <strong>{validRows.length}</strong> application{validRows.length === 1 ? '' : 's'}{skippedCount ? ` · ${skippedCount} row${skippedCount === 1 ? '' : 's'} skipped (missing company or role)` : ''}.
          </div>
          <div class="prev-wrap">
            <table class="prev">
              <thead><tr><th>Company</th><th>Role</th><th>Status</th><th>Source</th></tr></thead>
              <tbody>
                {#each validRows.slice(0, 5) as r}
                  <tr><td>{r.company}</td><td>{r.role}</td><td>{r.status || 'applied'}</td><td>{r.source || '—'}</td></tr>
                {/each}
              </tbody>
            </table>
            {#if validRows.length > 5}<div class="prev-more">+ {validRows.length - 5} more</div>{/if}
          </div>
        {/if}

        {#if error}<p class="err">{error}</p>{/if}
        <div class="foot">
          <button class="btn" onclick={() => (step = 'paste')}>Back</button>
          <button class="btn btn-primary" onclick={doImport} disabled={importing || !validRows.length || !hasCompany || !hasRole}>
            {importing ? 'Importing…' : `Import ${validRows.length} application${validRows.length === 1 ? '' : 's'}`}
          </button>
        </div>
      {/if}
    </div>
  </div>
{/if}

<style>
  .ovl { position: fixed; inset: 0; background: rgba(10,10,13,0.45); backdrop-filter: blur(2px); display: grid; place-items: center; z-index: 100; padding: 2rem; }
  .card { position: relative; background: var(--card); border: 1px solid var(--rule); border-radius: 14px; width: 100%; max-width: 600px; max-height: calc(100vh - 4rem); overflow-y: auto; padding: 22px; box-shadow: var(--sh-pop); }
  .x { position: absolute; top: 16px; right: 16px; background: transparent; border: 0; color: var(--mute); width: 28px; height: 28px; border-radius: 6px; display: grid; place-items: center; cursor: pointer; }
  .x:hover { background: var(--surface-2); color: var(--ink); }
  .hd { margin-bottom: 14px; padding-right: 28px; }
  .hd h2 { font-size: 19px; font-weight: 500; letter-spacing: -0.018em; margin: 0; }
  .sub { font-size: 13px; color: var(--mute); margin: 5px 0 0; line-height: 1.5; }
  .paste { width: 100%; box-sizing: border-box; font-family: var(--mono, ui-monospace, monospace); font-size: 12.5px; line-height: 1.6; color: var(--ink); background: var(--surface); border: 1px solid var(--rule); border-radius: 10px; padding: 12px; outline: none; resize: vertical; }
  .paste:focus { border-color: var(--accent); box-shadow: 0 0 0 3px var(--accent-tint); }
  .chk { display: flex; align-items: center; gap: 8px; font-size: 13px; color: var(--ink-2); margin-top: 12px; cursor: pointer; }
  .map-grid { display: flex; flex-direction: column; gap: 8px; }
  .map-row { display: grid; grid-template-columns: 1fr 18px 180px; gap: 10px; align-items: center; }
  .map-col { min-width: 0; background: var(--surface); border: 1px solid var(--rule); border-radius: 8px; padding: 7px 11px; }
  .map-h { font-size: 13px; font-weight: 500; color: var(--ink); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  .map-sample { font-size: 11.5px; color: var(--mute); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  .map-arrow { color: var(--mute-2); text-align: center; }
  .map-sel { font: inherit; font-size: 13px; color: var(--ink); background: var(--surface); border: 1px solid var(--rule); border-radius: 8px; padding: 8px 10px; outline: none; cursor: pointer; }
  .map-sel:focus { border-color: var(--accent); }
  .preview-note { font-size: 12.5px; color: var(--ink-2); margin: 16px 0 10px; }
  .warn { font-size: 12.5px; color: var(--warm-text, var(--accent-text)); background: var(--warm-tint, var(--accent-tint)); border-radius: 8px; padding: 9px 12px; margin: 16px 0 0; }
  .prev-wrap { border: 1px solid var(--rule); border-radius: 10px; overflow: hidden; }
  .prev { width: 100%; border-collapse: collapse; font-size: 12.5px; }
  .prev th { text-align: left; font-weight: 500; color: var(--mute); padding: 8px 12px; border-bottom: 1px solid var(--rule); background: var(--surface); }
  .prev td { padding: 8px 12px; border-bottom: 1px solid var(--rule); color: var(--ink-2); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 160px; }
  .prev tr:last-child td { border-bottom: 0; }
  .prev-more { font-size: 12px; color: var(--mute); padding: 8px 12px; background: var(--surface); }
  .done-line { display: flex; align-items: center; gap: 9px; font-size: 14px; color: var(--ink); margin: 6px 0 0; }
  .done-ic { display: inline-grid; place-items: center; width: 22px; height: 22px; border-radius: 50%; background: var(--positive-tint); color: var(--positive-text); font-size: 12px; font-weight: 700; }
  .err { font-size: 12.5px; color: var(--danger-text); margin: 12px 0 0; }
  .foot { display: flex; align-items: center; justify-content: space-between; gap: 10px; margin-top: 18px; }
  .btn { background: var(--card); border: 1px solid var(--rule); border-radius: 8px; padding: 8px 14px; font-size: 13px; font-weight: 500; color: var(--ink); cursor: pointer; }
  .btn:hover { border-color: var(--rule-strong); }
  .btn-primary { background: var(--accent); border-color: var(--accent-strong); color: white; margin-left: auto; }
  .btn-primary:hover { background: var(--accent-strong); }
  .btn-primary:disabled { opacity: 0.55; cursor: not-allowed; }

  @media (max-width: 640px) {
    .card { max-width: 100%; border-radius: 0; min-height: 100vh; }
    .map-row { grid-template-columns: 1fr 14px 130px; }
  }
</style>
