#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

watch=0
args=()
for arg in "$@"; do
	case "$arg" in
	-watch|--watch) watch=1 ;;
	*) args+=("$arg") ;;
	esac
done

run_once() {
	./scripts/check-guardrails.sh
	exec env GOCACHE=/tmp/go-build go run . "${args[@]}"
}

if [ "$watch" -eq 0 ]; then
	run_once
fi

snapshot() {
	rg --files -g '*.go' -g 'go.mod' -g 'go.sum' |
		sort |
		xargs stat -c '%n:%Y:%s' |
		sha1sum |
		cut -d' ' -f1
}

pid=""
bin="/tmp/quo2-dev-watch"
if [ "${#args[@]}" -gt 0 ]; then
	key="$(printf '%s\0' "${args[@]}" | sha1sum | cut -d' ' -f1)"
	bin="/tmp/quo2-dev-watch-${key}"
fi

start_app() {
	echo "[run-dev] starting app"
	(
		./scripts/check-guardrails.sh
		env GOCACHE=/tmp/go-build go build -o "$bin" .
		exec "$bin" "${args[@]}"
	) &
	pid=$!
	echo "[run-dev] app pid=$pid"
}

stop_app() {
	if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
		echo "[run-dev] stopping pid=$pid"
		kill "$pid" 2>/dev/null || true
		for _ in $(seq 1 30); do
			if ! kill -0 "$pid" 2>/dev/null; then break; fi
			sleep 0.1
		done
		if kill -0 "$pid" 2>/dev/null; then
			kill -9 "$pid" 2>/dev/null || true
		fi
		wait "$pid" 2>/dev/null || true
	fi
	pid=""
}

shutdown() {
	echo "[run-dev] shutting down"
	stop_app
	exit 0
}

trap shutdown INT TERM

last="$(snapshot)"
echo "[run-dev] watch mode: go files only (ignores css)"
start_app

while true; do
	sleep 1
	now="$(snapshot)"
	if [ "$now" != "$last" ]; then
		echo "[run-dev] go change detected, rebuilding"
		last="$now"
		stop_app
		start_app
	fi
done
