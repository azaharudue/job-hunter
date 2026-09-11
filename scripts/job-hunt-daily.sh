#!/usr/bin/env bash
# Daily job-hunt run. Invoked by cron. Usage: job-hunt-daily.sh
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LOG="$ROOT/joblists/job-hunt-run.log"
mkdir -p "$ROOT/joblists"
echo "[$(date '+%Y-%m-%d %H:%M:%S')] start" >> "$LOG"
cd "$ROOT"
opencode run --agent job-hunter --auto "Run your daily job-hunt routine now." >> "$LOG" 2>&1
echo "[$(date '+%Y-%m-%d %H:%M:%S')] exit=$?" >> "$LOG"