#!/bin/bash
# SessionStart hook for Pursuit (claude-code-on-the-web).
#
# Two jobs:
#  1. Make sure Go + pnpm deps are installed so the agent can build, test,
#     and lint without a cold-start round trip.
#  2. Print the user's free-run notes + a one-line repo state summary so
#     the next agent walks in already knowing what's open and where we are.
#
# Stays synchronous on the first iteration per the skill guidance — guarantees
# the agent never races a missing dependency.
set -euo pipefail

# Only do install work inside the remote/web container; locally the user has
# their own dev env and this would just churn.
if [ "${CLAUDE_CODE_REMOTE:-}" = "true" ]; then
  if [ -f "$CLAUDE_PROJECT_DIR/go.mod" ]; then
    (cd "$CLAUDE_PROJECT_DIR" && go mod download) >/dev/null 2>&1 || true
  fi
  if [ -f "$CLAUDE_PROJECT_DIR/web/pnpm-lock.yaml" ]; then
    (cd "$CLAUDE_PROJECT_DIR/web" && pnpm install --frozen-lockfile) >/dev/null 2>&1 || true
  fi
fi

# Stdout from this hook becomes context for the agent. Print the notes file
# first (the headline reason this hook exists), then a tight git summary.
cd "$CLAUDE_PROJECT_DIR"

echo "## Pursuit — session start context"
echo

if [ -f FREE_RUN_NOTES.md ]; then
  echo "### Free run notes (FREE_RUN_NOTES.md)"
  echo
  cat FREE_RUN_NOTES.md
  echo
else
  echo "_(FREE_RUN_NOTES.md not found — create one to log free-run observations.)_"
  echo
fi

echo "### Repo state"
echo
echo "Branch: $(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo 'detached')"
echo "Latest on main:"
git log --oneline -5 origin/main 2>/dev/null || git log --oneline -5 main 2>/dev/null || echo "  (no main ref yet)"
echo

# Status only if there are uncommitted changes — a clean tree shouldn't take
# up space in the agent's context.
if [ -n "$(git status --porcelain 2>/dev/null)" ]; then
  echo "Uncommitted changes:"
  git status --short
fi
