#!/bin/sh
set -e

DB_PATH="${DATABASE_PATH:-/app/db/gesitr.db}"

# Seed compendium data on first run
if [ ! -f "$DB_PATH" ]; then
  echo "First run detected — seeding compendium data..."
  ./seed
fi

exec ./gesitr
