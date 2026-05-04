#!/usr/bin/env bash
# Clears application data in the local MySQL used by docker compose (default DB: jobcrawler).
# Usage:
#   ./scripts/reset-jobcrawler-db.sh              # TRUNCATE/DELETE while stack is running
#   ./scripts/reset-jobcrawler-db.sh --down-volumes   # remove all compose volumes (fresh DB files)

set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

if [[ "${1:-}" == "--down-volumes" ]]; then
  docker compose down -v
  echo "Compose volumes removed. Run 'docker compose up -d' for a completely empty MySQL data directory."
  exit 0
fi

DB="${MYSQL_DATABASE:-jobcrawler}"
ROOTPW="${MYSQL_ROOT_PASSWORD:-root}"

docker compose exec -T mysql mysql -uroot -p"${ROOTPW}" "${DB}" <<'SQL'
SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE feed_decisions;
TRUNCATE TABLE notification_deliveries;
TRUNCATE TABLE notification_settings;
TRUNCATE TABLE jobs;
DELETE FROM users;
SET FOREIGN_KEY_CHECKS = 1;
SQL

echo "Cleared users, jobs, decisions, and notification tables in database '${DB}'."
