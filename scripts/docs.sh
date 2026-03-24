#!/usr/bin/env bash
set -euo pipefail

OUT_DIR="docs/generated"
GOMARKDOC="go run github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest"

rm -rf "$OUT_DIR"
mkdir -p "$OUT_DIR"

# Find all packages that have a doc.go and generate markdown for each.
for docfile in $(find internal -name doc.go -not -path '*/docs/*'); do
  pkg=$(dirname "$docfile")
  # Derive a slug from the package path: internal/exercise/handlers -> exercise-handlers
  slug=$(echo "$pkg" | sed 's|^internal/||; s|/|-|g; s|^user-||')
  $GOMARKDOC -o "$OUT_DIR/$slug.md" "./$pkg/"
done

echo "Generated:"
ls "$OUT_DIR"/*.md
