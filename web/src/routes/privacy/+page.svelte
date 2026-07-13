<script>
  // Published interviewer-privacy policy (STRATEGY §9.2). Pursuit compiles prep
  // on a named person who never consented; this page states, publicly, the line
  // we don't cross. Linked from the homepage and the playbook's interviewer card.
  const CONTACT = 'back.yonatan@gmail.com';
  const removeHref = `mailto:${CONTACT}?subject=${encodeURIComponent('Pursuit — remove me from interviewer research')}`;
</script>

<svelte:head>
  <title>How Pursuit researches people — Pursuit</title>
</svelte:head>

<main class="doc">
  <article>
    <a class="back" href="/">← Pursuit</a>
    <h1>How Pursuit researches people</h1>
    <p class="lede">
      Pursuit prepares you for interviews — and part of that is a briefing on the
      person interviewing you. That means researching a named individual who
      didn't ask to be researched, so we hold ourselves to a clear line. This is
      it.
    </p>

    <section>
      <h2>Public professional sources only</h2>
      <p>
        Interviewer briefings are built from what a person has published, or what
        has been published about their <strong>professional</strong> work:
        conference talks, essays, podcasts, papers, company bios, public
        engineering/design writing, LinkedIn, and interview-review sites. Nothing
        behind a login, a paywall, or a private database. No scraping of private
        profiles.
      </p>
    </section>

    <section>
      <h2>What we never do</h2>
      <ul>
        <li>No personal-life inference — family, health, religion, politics, home address, relationships. If it isn't about their professional work, it isn't in the brief.</li>
        <li>No guessing. If we can't find and source a claim, we leave it out rather than invent it.</li>
        <li>No private or purchased data — no people-search databases, no leaked datasets, no private contact scraping.</li>
        <li>We don't notify the interviewer, and we don't share your prep with them or anyone else.</li>
      </ul>
    </section>

    <section>
      <h2>Every claim is sourced</h2>
      <p>
        Each briefing links to the specific page a claim came from — a deep link
        to the essay, talk, or review, not a homepage. You can check it before you
        walk in, and if a claim has no source behind it, it doesn't belong in the
        brief. We run an automated check that a brief identifies the right company
        and that its citations resolve before we ship changes to how briefs are
        generated.
      </p>
    </section>

    <section>
      <h2>It's prep, not a profile</h2>
      <p>
        A briefing is scoped to one thing: helping you have a better, more
        informed conversation in a specific interview. It's about their public
        professional work and how they're known to interview — not a standing
        dossier on a private citizen.
      </p>
    </section>

    <section>
      <h2>Your data</h2>
      <p>
        Everything you enter — applications, notes, salary, prep briefs — is
        private to your account; no other user sees it. We store your Google email
        and profile picture for sign-in. Briefs are generated with Anthropic's
        API, which doesn't train on your data. Want your account deleted? One email
        and we remove everything.
      </p>
    </section>

    <section>
      <h2>If you're an interviewer and want out</h2>
      <p>
        Email <a href={removeHref}>{CONTACT}</a> and we'll remove any briefing that
        names you and suppress your name from future research. If you're in the EU,
        we honor data-subject requests (access, correction, erasure) under GDPR.
      </p>
    </section>

    <p class="foot">Questions: <a href={`mailto:${CONTACT}`}>{CONTACT}</a></p>
  </article>
</main>

<style>
  .doc {
    min-height: 100vh;
    background: var(--surface);
    padding: 48px 20px 80px;
    display: flex;
    justify-content: center;
  }
  article {
    width: 100%;
    max-width: 640px;
    color: var(--ink-2);
    font-size: 15px;
    line-height: 1.65;
  }
  .back {
    display: inline-block;
    font-size: 13px;
    color: var(--mute);
    text-decoration: none;
    margin-bottom: 28px;
  }
  .back:hover { color: var(--ink-2); }
  h1 {
    font-size: 30px;
    font-weight: 500;
    letter-spacing: -0.025em;
    color: var(--ink);
    margin: 0 0 14px;
  }
  .lede { font-size: 16px; color: var(--ink-2); margin: 0 0 8px; }
  section { margin-top: 30px; }
  h2 {
    font-size: 15px;
    font-weight: 600;
    color: var(--ink);
    margin: 0 0 8px;
  }
  p { margin: 0 0 10px; }
  ul { margin: 0; padding-left: 20px; }
  li { margin-bottom: 7px; }
  a { color: var(--accent-text); }
  strong { color: var(--ink); font-weight: 600; }
  .foot {
    margin-top: 40px;
    padding-top: 20px;
    border-top: 1px solid var(--rule);
    font-size: 13px;
    color: var(--mute);
  }
</style>
