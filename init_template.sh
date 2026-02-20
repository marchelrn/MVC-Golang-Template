#!/usr/bin/env bash
set -euo pipefail

OLD_MODULE="mini_jira"

if [[ $# -ne 1 ]]; then
  echo "Usage: $0 <new-go-module>"
  echo "Example: $0 github.com/username/my_project"
  exit 1
fi

NEW_MODULE="$1"

if ! command -v rg >/dev/null 2>&1; then
  echo "Error: ripgrep (rg) is required."
  exit 1
fi

if ! command -v go >/dev/null 2>&1; then
  echo "Error: go is required."
  exit 1
fi

go mod edit -module "$NEW_MODULE"

while IFS= read -r file; do
  sed -i "s|\"${OLD_MODULE}/|\"${NEW_MODULE}/|g" "$file"
done < <(rg -l "\"${OLD_MODULE}/" \
  --glob '*.go' \
  --glob '!mvc_template/**')

echo "Module initialized:"
echo "- old: ${OLD_MODULE}"
echo "- new: ${NEW_MODULE}"
echo "Next steps:"
echo "1) cp .env.example .env"
echo "2) go mod tidy"
echo "3) go run main.go"
