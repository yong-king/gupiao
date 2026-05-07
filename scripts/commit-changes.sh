#!/usr/bin/env sh
set -eu

CONFIG_PATH="${JIJIN_REPOSITORY_CONFIG:-config/repository.local.json}"
if [ ! -f "$CONFIG_PATH" ]; then
  CONFIG_PATH="config/repository.example.json"
fi
REMOTE_URL="${REPOSITORY_REMOTE_URL:-}"
BRANCH="${REPOSITORY_BRANCH:-}"
PREFIX="${COMMIT_MESSAGE_PREFIX:-}"

if command -v python3 >/dev/null 2>&1 && [ -f "$CONFIG_PATH" ]; then
  REMOTE_URL="${REMOTE_URL:-$(python3 -c 'import json,sys; print(json.load(open(sys.argv[1])).get("remote_url",""))' "$CONFIG_PATH")}"
  BRANCH="${BRANCH:-$(python3 -c 'import json,sys; print(json.load(open(sys.argv[1])).get("branch","main"))' "$CONFIG_PATH")}"
  PREFIX="${PREFIX:-$(python3 -c 'import json,sys; print(json.load(open(sys.argv[1])).get("commit_message_prefix","chore: update stock agent project"))' "$CONFIG_PATH")}"
fi

BRANCH="${BRANCH:-main}"
PREFIX="${PREFIX:-chore: update stock agent project}"

echo "This will stage current changes and commit them."
printf "Commit now? [y/N] "
read answer
case "$answer" in
  y|Y|yes|YES) ;;
  *) echo "commit skipped"; exit 0 ;;
esac

if [ ! -d .git ]; then
  git init
  git checkout -b "$BRANCH"
elif ! git rev-parse --abbrev-ref HEAD >/dev/null 2>&1; then
  git checkout -b "$BRANCH"
fi

if [ -n "$REMOTE_URL" ]; then
  if git remote get-url origin >/dev/null 2>&1; then
    git remote set-url origin "$REMOTE_URL"
  else
    git remote add origin "$REMOTE_URL"
  fi
fi

if [ -z "$(git status --short)" ]; then
  echo "no changes to commit"
  exit 0
fi

git add .
git commit -m "$PREFIX"

if [ -n "$REMOTE_URL" ]; then
  printf "Push to origin/%s now? [y/N] " "$BRANCH"
  read push_answer
  case "$push_answer" in
    y|Y|yes|YES) git push -u origin "$BRANCH" ;;
    *) echo "push skipped";;
  esac
fi
