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

  # Convert escaped "OpenAPI: /api/docs\#/operations/op\-id" lines into
  # clickable markdown links: [OpenAPI docs](/api/docs#/operations/op-id)
  perl -pi -e 's{OpenAPI: /api/docs\\#/operations/([\w\\-]+)}{
    my $id = $1; $id =~ s/\\-/-/g; "[OpenAPI docs](/api/docs#/operations/$id)"
  }ge' "$OUT_DIR/$slug.md"
done

echo "Generated:"
ls "$OUT_DIR"/*.md
