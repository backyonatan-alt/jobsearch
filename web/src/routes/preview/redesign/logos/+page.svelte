<script>
  // Standalone logo gallery — every company we use across the demos,
  // pulled from logo.clearbit.com. This is the stable logo library we'll
  // ship with. Letter-square is the fallback when Clearbit has no hit.
  const companies = [
    { name: 'Stripe',     domain: 'stripe.com' },
    { name: 'Vercel',     domain: 'vercel.com' },
    { name: 'Anthropic',  domain: 'anthropic.com' },
    { name: 'Linear',     domain: 'linear.app' },
    { name: 'Notion',     domain: 'notion.so' },
    { name: 'Supabase',   domain: 'supabase.com' },
    { name: 'Figma',      domain: 'figma.com' },
    { name: 'Cursor',     domain: 'cursor.com' },
    { name: 'Modal',      domain: 'modal.com' },
    { name: 'Replicate',  domain: 'replicate.com' },
    { name: 'Perplexity', domain: 'perplexity.ai' },
    { name: 'Mistral',    domain: 'mistral.ai' },
    { name: 'Pinecone',   domain: 'pinecone.io' },
    { name: 'OpenAI',     domain: 'openai.com' },
    { name: 'Google',     domain: 'google.com' },
    { name: 'Meta',       domain: 'meta.com' },
    { name: 'Airbnb',     domain: 'airbnb.com' },
    { name: 'Shopify',    domain: 'shopify.com' },
    { name: 'Cloudflare', domain: 'cloudflare.com' },
    { name: 'GitHub',     domain: 'github.com' },
    { name: 'Datadog',    domain: 'datadoghq.com' },
    { name: 'Asana',      domain: 'asana.com' }
  ];

  // The three sizes the design system uses.
  const sizes = [
    { lbl: 'Table row / chip', px: 22, radius: 5 },
    { lbl: 'Card / list',      px: 32, radius: 8 },
    { lbl: 'Detail header',    px: 64, radius: 14 }
  ];
</script>

<svelte:head><title>Logo gallery — Pursuit preview</title></svelte:head>

<main class="wrap">
  <header>
    <div class="kicker">Stable library · logo.clearbit.com</div>
    <h1>Company logos for the demo</h1>
    <p>
      Every logo on every Pursuit surface comes from
      <code>https://logo.clearbit.com/&lt;domain&gt;</code> — free, public, no API key.
      When Clearbit has no hit we fall back to a coloured letter square (last row below).
      This page is just so you can eyeball the library before we wire it in.
    </p>
  </header>

  <!-- The three sizes side by side, using Stripe as the example -->
  <section>
    <h2>Three sizes the design system uses</h2>
    <div class="sizes">
      {#each sizes as s}
        <div class="size-cell">
          <img
            src="https://logo.clearbit.com/stripe.com"
            alt="Stripe"
            class="size-img"
            style={`width:${s.px}px;height:${s.px}px;border-radius:${s.radius}px;`}
          />
          <div class="size-meta">
            <div class="size-lbl">{s.lbl}</div>
            <div class="size-sub">{s.px}×{s.px} · radius {s.radius}px</div>
          </div>
        </div>
      {/each}
    </div>
  </section>

  <!-- Grid of all the companies in our demo data -->
  <section>
    <h2>Every company in the demo data</h2>
    <div class="grid">
      {#each companies as c}
        <div class="logo-card">
          <img
            class="logo-img"
            src={`https://logo.clearbit.com/${c.domain}`}
            alt={c.name}
            loading="lazy"
          />
          <div class="logo-meta">
            <div class="co-name">{c.name}</div>
            <div class="co-domain">{c.domain}</div>
          </div>
        </div>
      {/each}
    </div>
  </section>

  <!-- Fallback rendering for unknown companies -->
  <section>
    <h2>Fallback for companies Clearbit doesn't have</h2>
    <p class="hint">When the CDN returns nothing (private companies, tiny startups, mistyped domains), we render a colored letter square keyed off the company name.</p>
    <div class="fallback-row">
      {#each ['Acme Stealth Co', 'NewCo', 'XYZ Robotics', 'Foo Bar Inc'] as name}
        <div class="logo-card">
          <span class="letter-sq" style={`background: hsl(${(name.charCodeAt(0) * 13) % 360} 65% 92%); color: hsl(${(name.charCodeAt(0) * 13) % 360} 55% 35%);`}>
            {name.split(' ').map(w => w[0]).join('').slice(0,2).toUpperCase()}
          </span>
          <div class="logo-meta">
            <div class="co-name">{name}</div>
            <div class="co-domain">no clearbit hit</div>
          </div>
        </div>
      {/each}
    </div>
  </section>

  <footer>
    <a href="/preview/redesign">← back to previews</a>
  </footer>
</main>

<style>
  :global(html, body) { background: var(--surface); margin: 0; }
  .wrap { max-width: 1020px; margin: 5vh auto; padding: 0 2rem 6rem; font-family: var(--sans); color: var(--ink); }

  header { margin-bottom: 36px; }
  .kicker { font-size: 12px; color: var(--accent-text); background: var(--accent-tint); display: inline-block; padding: 3px 10px; border-radius: 99px; font-weight: 500; margin-bottom: 12px; }
  header h1 { font-size: 28px; font-weight: 600; letter-spacing: -0.025em; margin: 0 0 .5rem; }
  header p { color: var(--mute); max-width: 68ch; margin: 0; line-height: 1.6; font-size: 14px; }
  header code { background: var(--surface-2); padding: 2px 6px; border-radius: 4px; font-size: 12.5px; color: var(--ink); }

  section { margin-bottom: 40px; }
  section h2 { font-size: 16px; font-weight: 600; letter-spacing: -0.015em; margin: 0 0 14px; }
  .hint { font-size: 13px; color: var(--mute); margin: 0 0 14px; line-height: 1.55; max-width: 70ch; }

  /* Three-sizes preview */
  .sizes {
    display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px;
    background: var(--card); border: 1px solid var(--rule); border-radius: 14px;
    padding: 22px; box-shadow: var(--sh-1);
  }
  .size-cell { display: flex; align-items: center; gap: 14px; }
  .size-img { background: var(--surface-2); object-fit: contain; padding: 3px; border: 1px solid var(--rule); }
  .size-lbl { font-size: 13px; font-weight: 600; }
  .size-sub { font-size: 11.5px; color: var(--mute); margin-top: 2px; }

  /* Logo grid */
  .grid {
    display: grid; grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
    gap: 10px;
  }
  .logo-card {
    background: var(--card); border: 1px solid var(--rule); border-radius: 12px;
    padding: 14px 16px;
    display: flex; align-items: center; gap: 12px;
    box-shadow: var(--sh-1);
    transition: border-color 120ms ease, transform 120ms ease;
  }
  .logo-card:hover { border-color: var(--rule-strong); transform: translateY(-1px); }
  .logo-img {
    width: 36px; height: 36px; border-radius: 8px;
    background: var(--surface-2); object-fit: contain; padding: 4px;
    border: 1px solid var(--rule); flex-shrink: 0;
  }
  .logo-meta { min-width: 0; }
  .co-name { font-size: 13.5px; font-weight: 600; letter-spacing: -0.01em; }
  .co-domain { font-size: 11.5px; color: var(--mute); margin-top: 1px; }

  /* Fallback row */
  .fallback-row { display: grid; grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); gap: 10px; }
  .letter-sq {
    width: 36px; height: 36px; border-radius: 8px;
    display: grid; place-items: center;
    font-size: 13px; font-weight: 700; letter-spacing: -0.01em;
    flex-shrink: 0;
    border: 1px solid var(--rule);
  }

  footer { margin-top: 40px; font-size: 13px; }
  footer a { color: var(--accent-text); text-decoration: none; }
</style>
