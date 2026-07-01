// Preview mode — lets the user click through /app/* without a backend.
//
// Triggered by appending ?preview=1 to any /app/* URL. Once set, it sticks
// across navigation via sessionStorage so deep links inside /app/[id] stay
// in preview without re-adding the query param.
//
// The `mockApi` function returns the same shape the real Go handlers do.
// Mutations (PATCH/POST/DELETE) apply to an in-memory copy of the apps so
// the UI feels live; closing the tab discards them.

const PREVIEW_KEY = 'pursuit_preview_mode';

export function isPreview() {
  if (typeof window === 'undefined') return false;
  try {
    const url = new URL(window.location.href);
    if (url.searchParams.get('preview') === '1') {
      sessionStorage.setItem(PREVIEW_KEY, '1');
      return true;
    }
    if (url.searchParams.get('preview') === '0') {
      sessionStorage.removeItem(PREVIEW_KEY);
      return false;
    }
    return sessionStorage.getItem(PREVIEW_KEY) === '1';
  } catch {
    return false;
  }
}

// ── Fixtures ───────────────────────────────────────────────
// Synthesized to populate every surface: count cards, action grid,
// insights ("3 stale", "referrals convert better", "haven't applied
// in N days"), funnel with conversion rates, board with stale cards,
// brief with hiring manager + dossier on one app.

const PREVIEW_USER = {
  id: 1,
  email: 'yonatan@pursuit.app',
  is_admin: false,
  onboarded_at: '2026-01-01T00:00:00Z',
  picture_url: 'https://www.google.com/s2/favicons?sz=128&domain=pursuit.app'
};

function daysAgoIso(n) {
  const d = new Date();
  d.setDate(d.getDate() - n);
  return d.toISOString();
}

const INITIAL_APPS = [
  { id: 101, company: 'Anthropic',      role: 'Member of Technical Staff, Applied AI', status: 'offer',     source: 'Referral',    location: 'San Francisco / Remote', salary_note: '$340k base + equity', cv_variant: 'v3-ai-focus',   jd_url: 'https://www.anthropic.com/jobs', notes: 'Loop went well. Verbal offer Tue.', hiring_manager_name: 'Devon Marquez', hiring_manager_linkedin: 'https://www.linkedin.com/in/example-devon-marquez', applied_at: daysAgoIso(38) },
  { id: 102, company: 'Stripe',         role: 'Staff Engineer, Platform',               status: 'interview', source: 'Referral',    location: 'Remote (US)',           salary_note: '$280-320k base',     cv_variant: 'v2-infra',       jd_url: null,                          notes: 'Onsite Thursday — 4 panels.',          hiring_manager_name: null,            hiring_manager_linkedin: null,                                              applied_at: daysAgoIso(24) },
  { id: 103, company: 'Linear',         role: 'Senior Product Engineer',               status: 'interview', source: 'LinkedIn',    location: 'Remote (EU)',           salary_note: '$210k base',         cv_variant: 'v3-ai-focus',    jd_url: null,                          notes: null,                                  hiring_manager_name: null,            hiring_manager_linkedin: null,                                              applied_at: daysAgoIso(17) },
  { id: 104, company: 'Vercel',         role: 'Senior Software Engineer, Edge',        status: 'screen',    source: 'Greenhouse',  location: 'Remote',                salary_note: null,                 cv_variant: 'v2-infra',       jd_url: null,                          notes: null,                                  hiring_manager_name: null,            hiring_manager_linkedin: null,                                              applied_at: daysAgoIso(11) },
  { id: 105, company: 'Figma',          role: 'Senior Engineer, Multiplayer',          status: 'screen',    source: 'Referral',    location: 'San Francisco',         salary_note: null,                 cv_variant: 'v3-ai-focus',    jd_url: null,                          notes: null,                                  hiring_manager_name: null,            hiring_manager_linkedin: null,                                              applied_at: daysAgoIso(9) },
  { id: 106, company: 'Notion',         role: 'Staff Engineer, AI',                    status: 'screen',    source: 'Cold email',  location: 'San Francisco',         salary_note: null,                 cv_variant: null,             jd_url: null,                          notes: null,                                  hiring_manager_name: null,            hiring_manager_linkedin: null,                                              applied_at: daysAgoIso(6) },
  { id: 107, company: 'Cursor',         role: 'Founding Engineer, Agents',             status: 'applied',   source: 'X DM',        location: 'San Francisco',         salary_note: null,                 cv_variant: null,             jd_url: null,                          notes: null,                                  hiring_manager_name: null,            hiring_manager_linkedin: null,                                              applied_at: daysAgoIso(4) },
  { id: 108, company: 'Mistral AI',     role: 'Senior Research Engineer',              status: 'applied',   source: 'Company site',location: 'Paris / Remote',        salary_note: null,                 cv_variant: 'v3-ai-focus',    jd_url: null,                          notes: null,                                  hiring_manager_name: null,            hiring_manager_linkedin: null,                                              applied_at: daysAgoIso(3) },
  { id: 109, company: 'Perplexity',     role: 'Staff Engineer, Search',                status: 'applied',   source: 'LinkedIn',    location: 'San Francisco',         salary_note: null,                 cv_variant: null,             jd_url: 'https://www.linkedin.com/jobs/view/4012345678', notes: null,                        hiring_manager_name: null,            hiring_manager_linkedin: null,                                              applied_at: daysAgoIso(2) },
  { id: 110, company: 'Granola',        role: 'Founding Backend Engineer',             status: 'applied',   source: 'Referral',    location: 'London',                salary_note: null,                 cv_variant: null,             jd_url: null,                          notes: null,                                  hiring_manager_name: null,            hiring_manager_linkedin: null,                                              applied_at: daysAgoIso(2) },
  { id: 111, company: 'OpenAI',         role: 'Member of Technical Staff',             status: 'applied',   source: 'Company site',location: 'San Francisco',         salary_note: null,                 cv_variant: 'v3-ai-focus',    jd_url: null,                          notes: null,                                  hiring_manager_name: null,            hiring_manager_linkedin: null,                                              applied_at: daysAgoIso(1) },
  // Stale (>= 7 days in applied/screen) — drives red border + stale insights.
  { id: 112, company: 'Replicate',      role: 'Platform Engineer',                     status: 'applied',   source: 'Cold',        location: 'Remote',                salary_note: null,                 cv_variant: null,             jd_url: null,                          notes: null,                                  hiring_manager_name: null,            hiring_manager_linkedin: null,                                              applied_at: daysAgoIso(12) },
  { id: 113, company: 'Modal',          role: 'Infra Engineer',                        status: 'closed',    source: 'LinkedIn',    location: 'Remote',                salary_note: null,                 cv_variant: 'v2-infra',       jd_url: null,                          notes: 'Req cancelled mid-process.',          hiring_manager_name: null,            hiring_manager_linkedin: null,                                              applied_at: daysAgoIso(15) },
  { id: 114, company: 'Arc (The Browser Company)', role: 'Senior Engineer, AI Browser', status: 'wishlist', source: 'LinkedIn',   location: 'Remote',                salary_note: null,                 cv_variant: null,             jd_url: null,                          notes: null,                                  hiring_manager_name: null,            hiring_manager_linkedin: null,                                              applied_at: null },
  { id: 115, company: 'Plaid',          role: 'Senior Engineer, Identity',             status: 'rejected',  source: 'LinkedIn',    location: 'Remote (US)',           salary_note: null,                 cv_variant: 'v1-generalist',  jd_url: null,                          notes: 'Failed at sys design.',               hiring_manager_name: null,            hiring_manager_linkedin: null,                                              applied_at: daysAgoIso(52) }
];

// In-memory state for the session. Mutations stick until reload.
// Demo rows are tagged `_demo` so the guided tour can seed/clear them in
// preview exactly like the real backend (which tags via a [demo] notes prefix).
const demoSet = () => INITIAL_APPS.map(a => ({
  ...a, _demo: true, last_follow_up_at: null,
  created_at: a.applied_at || new Date().toISOString(),
  updated_at: new Date().toISOString()
}));
let apps = demoSet();
let interviewsByApp = {}; // appId -> []
let dossiersByApp = {};   // appId -> { content, meeting, generatedAgo, interviewer_name }
let followUpsByApp = {};  // appId -> [] (newest first)

// Seed one quiet app (Modal #113) with an existing follow-up so the timeline
// shows real content out of the box.
{
  const seededAt = daysAgoIso(3);
  followUpsByApp[113] = [{
    id: 1,
    application_id: 113,
    note: 'Emailed the recruiter to check in on timeline',
    channel: 'Email',
    occurred_at: seededAt,
    created_at: seededAt
  }];
  const modal = apps.find(a => a.id === 113);
  if (modal) modal.last_follow_up_at = seededAt;
}

function findApp(id) { return apps.find(a => a.id === Number(id)); }

// One pre-baked dossier so the Brief tab renders fully on Stripe.
dossiersByApp[102] = {
  interviewer_name: 'Sarah Chen',
  generatedAgo: '12 min ago',
  meeting: {
    starts_at: daysAgoIso(-1).replace(/T.*/, 'T18:00:00Z'),
    ends_at:   daysAgoIso(-1).replace(/T.*/, 'T19:00:00Z'),
    medium: 'Google Meet',
    panel: '1:1 with Sarah Chen'
  },
  content: {
    identity: {
      name: 'Stripe',
      domain: 'stripe.com',
      summary: 'Payments & financial infrastructure · San Francisco'
    },
    sources: [
      { label: 'TechCrunch — Stripe ships agentic-commerce APIs (May 2026)', href: 'https://techcrunch.com/2026/05/stripe-agentic-commerce/' },
      { label: 'Glassdoor — Stripe interview reviews', href: 'https://www.glassdoor.com/Interview/Stripe-Interview-Questions-E671954.htm' }
    ],
    interviewer: {
      name: 'Sarah Chen',
      role: 'Staff Engineer · Payments Routing',
      prior: ['Cloudflare', 'Two Sigma', 'CMU MS'],
      links: [{ href: 'https://www.linkedin.com/in/example', label: 'LinkedIn' }]
    },
    snapshot: '2 years at Stripe on payments-routing. Known for her QCon SF 2024 talk on regional failover — likely to drill on failure modes more than happy paths.',
    background: 'Started on the data infrastructure team before moving to payments. Tends to ask deep questions about consistency tradeoffs and rollout strategy. Active on tech Twitter about idempotency-key collisions.',
    signals: [
      { date: 'Apr 26', kind: 'Talk', body: '"Regional failover at p99.99"',           source: 'qconsf.com',  source_url: 'https://qconsf.com/talk/regional-failover-p9999' },
      { date: 'Mar 04', kind: 'Post', body: 'Thread on idempotency-key collisions',     source: 'twitter.com', source_url: 'https://twitter.com/example/status/1764000000000000000' },
      { date: 'Jan 12', kind: 'Doc',  body: 'Cited in Stripe multi-region writes',      source: 'stripe.com',  source_url: 'https://stripe.com/blog/multi-region-writes' }
    ],
    lands: [
      'Concrete failure-mode reasoning',
      'Latency, RPS, error-rate numbers',
      'Rollback + observability'
    ],
    avoid: [
      'Hand-waving consistency tradeoffs',
      'Big-bang rewrites with no rollout',
      'Database choice before constraints'
    ],
    questions: [
      { q: "What's the hardest production incident this year?",  why: 'Real story; reveals what she owns vs. delegates.' },
      { q: 'Hot-patch legacy router vs. forward on v2?',          why: 'Tests whether the org pays down debt.' },
      { q: 'What does a great first 90 days look like?',          why: 'Standard, high-signal — loose answer = loose role.' }
    ],
    company: {
      blurb: 'Payments + financial infrastructure used by a meaningful share of the internet.',
      direction: 'Recently shipping agentic-commerce APIs and pushing into AI-mediated transactions. Last 12 months: agent toolkit, deeper Issuing rollouts, expanded Atlas.',
      hq: 'San Francisco',
      employees: '~8,000',
      stage: 'Late stage · $50B valuation',
      founded: '2010',
      process: [
        { kind: 'Recruiter screen',    detail: '30 min · culture + role fit' },
        { kind: 'Hiring-manager call', detail: '45 min · architecture, prior wins' },
        { kind: 'Technical screen',    detail: '60 min · live coding (Go or Python)' },
        { kind: 'Onsite loop',         detail: '4 panels · sys design, IC code, behavioral, bar-raiser' },
        { kind: 'Team match',          detail: 'Two 30-min chats with prospective teams' }
      ],
      watch_fors: [
        'Sys design is graded on rollout + observability, not just the diagram.',
        'Be explicit about failure modes — they probe consistency tradeoffs.',
        'Stripe interviewers value "what would you measure?" over hand-waving.',
        'Bar raiser is structured — STAR format pays off here.',
        "Have a take on agentic commerce; they'll ask where you think this goes."
      ]
    }
  }
};

// ── Mock router ────────────────────────────────────────────

function ok(body) { return Promise.resolve(body); }
function notFound(msg) { return Promise.reject(new Error(msg || 'not found')); }

export async function mockApi(path, opts = {}) {
  const method = (opts.method || 'GET').toUpperCase();
  // Strip a trailing query string.
  const cleanPath = path.split('?')[0];

  // GET /api/me
  if (method === 'GET' && cleanPath === '/api/me') return ok(PREVIEW_USER);

  // POST /api/auth/logout — no-op in preview.
  if (method === 'POST' && cleanPath === '/api/auth/logout') return ok({ status: 'ok' });

  // POST /api/me/onboarded — flips the flag on the fake user.
  if (method === 'POST' && cleanPath === '/api/me/onboarded') {
    PREVIEW_USER.onboarded_at = new Date().toISOString();
    return ok({ status: 'ok' });
  }

  // GET /api/applications
  if (method === 'GET' && cleanPath === '/api/applications') return ok(apps.slice());

  // POST /api/applications  — create
  if (method === 'POST' && cleanPath === '/api/applications') {
    const body = JSON.parse(opts.body || '{}');
    const id = Math.max(0, ...apps.map(a => a.id)) + 1;
    const now = new Date().toISOString();
    const created = {
      id,
      company: body.company, role: body.role, status: body.status || 'applied',
      source: body.source ?? null, jd_url: body.jd_url ?? null, location: body.location ?? null,
      salary_note: body.salary_note ?? null, cv_variant: body.cv_variant ?? null,
      notes: body.notes ?? null,
      hiring_manager_name: body.hiring_manager_name ?? null,
      hiring_manager_linkedin: body.hiring_manager_linkedin ?? null,
      applied_at: body.status === 'wishlist' ? null : (body.applied_at ?? now),
      last_follow_up_at: null,
      created_at: now, updated_at: now
    };
    apps = [created, ...apps];
    return ok(created);
  }

  // POST /api/applications/parse — fake parse result.
  if (method === 'POST' && cleanPath === '/api/applications/parse') {
    return ok({ company: 'Preview Co', role: 'Senior Engineer', source: 'LinkedIn', location: 'Remote' });
  }

  // POST/DELETE /api/me/demo-seed — mirror the real backend so the guided tour
  // demonstrates seed→clear in preview: POST (re)seeds the demo rows, DELETE
  // removes only demo rows (user-added apps, with no _demo flag, survive).
  if (cleanPath === '/api/me/demo-seed') {
    if (method === 'POST') {
      apps = [...demoSet(), ...apps.filter(a => !a._demo)];
      return ok({ inserted: INITIAL_APPS.length });
    }
    if (method === 'DELETE') {
      const before = apps.length;
      apps = apps.filter(a => !a._demo);
      return ok({ deleted: before - apps.length });
    }
  }

  // /api/applications/:id …
  const appMatch = cleanPath.match(/^\/api\/applications\/(\d+)(?:\/(.+))?$/);
  if (appMatch) {
    const id = Number(appMatch[1]);
    const sub = appMatch[2] || '';
    const app = findApp(id);
    if (!app && method === 'GET' && !sub) return notFound();

    // GET /api/applications/:id
    if (method === 'GET' && !sub) return ok(app);

    // PATCH /api/applications/:id — partial update.
    if (method === 'PATCH' && !sub) {
      if (!app) return notFound();
      const body = JSON.parse(opts.body || '{}');
      Object.assign(app, body);
      app.updated_at = new Date().toISOString();
      return ok(app);
    }

    // DELETE /api/applications/:id
    if (method === 'DELETE' && !sub) {
      apps = apps.filter(a => a.id !== id);
      return ok(null);
    }

    // GET /api/applications/:id/dossier
    if (method === 'GET' && sub === 'dossier') {
      const d = dossiersByApp[id];
      if (!d) return Promise.reject(new Error('no dossier'));
      return ok(d);
    }

    // POST /api/applications/:id/dossier/refresh — generate a stub.
    if (method === 'POST' && sub === 'dossier/refresh') {
      // Honor a short delay so the spinner is visible.
      await new Promise(r => setTimeout(r, 800));
      const body = JSON.parse(opts.body || '{}');
      const baseline = dossiersByApp[102];
      dossiersByApp[id] = { ...baseline, interviewer_name: body.interviewer_name || baseline.interviewer_name, generatedAgo: 'just now' };
      return ok(dossiersByApp[id]);
    }

    // GET /api/applications/:id/interviews
    if (method === 'GET' && sub === 'interviews') {
      return ok(interviewsByApp[id] || []);
    }

    // POST /api/applications/:id/interviews
    if (method === 'POST' && sub === 'interviews') {
      const body = JSON.parse(opts.body || '{}');
      const list = interviewsByApp[id] || [];
      const iid = Math.max(0, ...list.map(x => x.id)) + 1;
      const ev = { id: iid, ...body };
      interviewsByApp[id] = [...list, ev];
      return ok(ev);
    }

    // POST /api/applications/:id/interviews/parse — supports {ics}, {text}, {image}.
    // Returns the same {events:[...]} shape regardless so the preview UI can
    // exercise both zones.
    if (method === 'POST' && sub === 'interviews/parse') {
      const body = JSON.parse(opts.body || '{}');
      const source = body.image ? 'ai' : body.text ? 'ai' : 'ics';
      const summary = body.image
        ? 'Stripe — Technical screen (parsed from screenshot)'
        : body.text
          ? 'Stripe — Technical screen (parsed from email)'
          : 'Onsite — Stripe';
      // Two business days from now at 14:00 local — keeps the preview always future-dated.
      const d = new Date(); d.setDate(d.getDate() + 2); d.setHours(14, 0, 0, 0);
      const e = new Date(d); e.setHours(15, 0, 0, 0);
      return ok({ events: [{
        source, summary,
        starts_at: d.toISOString(),
        ends_at: e.toISOString(),
        location: 'Google Meet',
        attendees: [{ name: 'Sarah Chen' }, { email: 'recruiter@stripe.com' }]
      }] });
    }

    // DELETE /api/applications/:id/interviews/:iid
    const ivMatch = sub.match(/^interviews\/(\d+)$/);
    if (method === 'DELETE' && ivMatch) {
      const iid = Number(ivMatch[1]);
      interviewsByApp[id] = (interviewsByApp[id] || []).filter(x => x.id !== iid);
      return ok(null);
    }

    // GET /api/applications/:id/follow-ups — newest first.
    if (method === 'GET' && sub === 'follow-ups') {
      return ok((followUpsByApp[id] || []).slice());
    }

    // POST /api/applications/:id/follow-ups — record a follow-up the user did.
    if (method === 'POST' && sub === 'follow-ups') {
      const body = JSON.parse(opts.body || '{}');
      const now = new Date().toISOString();
      const occurredAt = body.occurred_at || now;
      const list = followUpsByApp[id] || [];
      const fid = Math.max(0, ...list.map(x => x.id)) + 1;
      const f = {
        id: fid,
        application_id: id,
        note: body.note ?? '',
        channel: body.channel ?? '',
        occurred_at: occurredAt,
        created_at: now
      };
      followUpsByApp[id] = [f, ...list];
      // Reset the quiet clock to the latest follow-up.
      if (app) {
        const latest = followUpsByApp[id]
          .map(x => x.occurred_at)
          .sort((a, b) => new Date(b) - new Date(a))[0];
        app.last_follow_up_at = latest;
      }
      return ok(f);
    }

    // DELETE /api/applications/:id/follow-ups/:fid
    const fuMatch = sub.match(/^follow-ups\/(\d+)$/);
    if (method === 'DELETE' && fuMatch) {
      const fid = Number(fuMatch[1]);
      followUpsByApp[id] = (followUpsByApp[id] || []).filter(x => x.id !== fid);
      if (app) {
        const remaining = followUpsByApp[id];
        app.last_follow_up_at = remaining.length
          ? remaining.map(x => x.occurred_at).sort((a, b) => new Date(b) - new Date(a))[0]
          : null;
      }
      return ok(null);
    }
  }

  return Promise.reject(new Error(`preview: no mock for ${method} ${cleanPath}`));
}
