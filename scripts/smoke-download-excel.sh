#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PORT="${1:-3333}"
HOST="${2:-127.0.0.1}"
TMP_XLSX="/tmp/quox.smoke.download.xlsx"
TMP_HDR="/tmp/quox.smoke.download.headers.txt"

cd "$ROOT"
rm -f "$TMP_XLSX" "$TMP_HDR"

curl -sS -D "$TMP_HDR" -o "$TMP_XLSX" \
	-X POST \
	-d 'DownloadExcel=slim=false' \
	"http://${HOST}:${PORT}/download-excel"

if ! grep -Eiq '^HTTP/[0-9.]+ 200' "$TMP_HDR"; then
	echo "smoke failed: non-200 response"
	sed -n '1,20p' "$TMP_HDR"
	exit 1
fi

if ! grep -Eiq '^Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' "$TMP_HDR"; then
	echo "smoke failed: missing xlsx content type"
	sed -n '1,20p' "$TMP_HDR"
	exit 1
fi

FILE_NAME="$(sed -n 's/^Content-Disposition:[[:space:]]*attachment; filename=\(.*\)\r\{0,1\}$/\1/ip' "$TMP_HDR" | tail -n 1)"
FILE_NAME="$(printf '%s' "$FILE_NAME" | tr -d '\r' | sed 's/^"//; s/"$//')"
if [[ -z "$FILE_NAME" ]]; then
	echo "smoke failed: missing content-disposition filename"
	sed -n '1,20p' "$TMP_HDR"
	exit 1
fi

if [[ ! -f "$TMP_XLSX" ]]; then
	echo "smoke failed: download file missing"
	exit 1
fi

if [[ ! -s "$TMP_XLSX" ]]; then
	echo "smoke failed: download file empty"
	exit 1
fi

if [[ ! -f "assets/work/$FILE_NAME" ]]; then
	echo "smoke failed: expected generated file not found in assets/work: $FILE_NAME"
	ls -lt assets/work | sed -n '1,20p'
	exit 1
fi

echo "smoke ok"
echo "file=$FILE_NAME"
ls -l "$TMP_XLSX" "assets/work/$FILE_NAME"
