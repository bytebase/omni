#!/usr/bin/env bash
# Boots OceanBase CE in MySQL mode, runs one command against it with
# OCEANBASE_WALK_DSN set, and removes the container, however the command
# ends. The WALK_THROUGH oracle in mysql/review skips without that DSN;
# CI has no room for the server (docs/specs/2026-10-22-walk-through-engines-design.md, D4).
#
#   scripts/oceanbase-oracle.sh go test ./mysql/review/ -run 'OceanBase' -count=1
#
# The server needs about 8 GB, so the script shares a devbox politely:
#   - it holds the heavy-engine lock, which the other heavy oracles (MSSQL,
#     Oracle) take too, so one of them runs at a time;
#   - inside the lock it waits until MemAvailable is at least
#     OB_GATE_MEM_GB (12) and at most OB_GATE_RUNNERS (1) CI jobs run,
#     polling every minute, and gives up with status 75 after
#     OB_GATE_TIMEOUT seconds (2 hours);
#   - the container is capped at 8 GB, swap included.
#
# OB_IMAGE overrides the image; the default is the LTS line Bytebase users
# run, 4.3.5. OB_LOG_DIR, when set, receives the observer's logs before the
# container goes.
set -euo pipefail

image=${OB_IMAGE:-oceanbase/oceanbase-ce:4.3.5-lts}
name=${OB_CONTAINER:-omni-walk-oceanbase}
lock=${OMNI_HEAVY_LOCK:-$(git -C "$(dirname "$0")" rev-parse --path-format=absolute --git-common-dir)/omni-heavy-engine.lock}
mem_gb=${OB_GATE_MEM_GB:-12}
runners=${OB_GATE_RUNNERS:-1}
timeout=${OB_GATE_TIMEOUT:-7200}

if [ $# -eq 0 ]; then
	echo "usage: $0 <command> [args...]" >&2
	exit 2
fi

if [ -z "${OB_ORACLE_LOCKED:-}" ]; then
	echo "$(date +%T) waiting for $lock"
	OB_ORACLE_LOCKED=1 exec flock "$lock" "$0" "$@"
fi

gate_open() {
	local avail jobs
	avail=$(awk '/^MemAvailable:/ {print $2}' /proc/meminfo)
	jobs=$(pgrep -fc Runner.Worker || true)
	echo "$(date +%T) MemAvailable $((avail / 1024 / 1024)) GB, $jobs CI jobs"
	[ "$avail" -ge $((mem_gb * 1024 * 1024)) ] && [ "$jobs" -le "$runners" ]
}

start=$(date +%s)
until gate_open; do
	if [ $(($(date +%s) - start)) -ge "$timeout" ]; then
		echo "the memory gate did not open within $timeout s" >&2
		exit 75
	fi
	sleep 60
done

docker rm -f "$name" >/dev/null 2>&1 || true
cleanup() {
	if [ -n "${OB_LOG_DIR:-}" ]; then
		mkdir -p "$OB_LOG_DIR"
		docker cp "$name:/root/ob/observer/log" "$OB_LOG_DIR" >/dev/null 2>&1 || true
	fi
	docker rm -f "$name" >/dev/null 2>&1 || true
}
trap cleanup EXIT
# The deployer refuses to start the observer below 20000 open files.
docker run -d --name "$name" -e MODE=mini --memory=8g --memory-swap=8g --ulimit nofile=65536:65536 \
	-p 127.0.0.1::2881 "$image" >/dev/null

# The image prints "boot success!" once the tenant takes connections.
echo "$(date +%T) booting $image"
for _ in $(seq 1 180); do
	logs=$(docker logs "$name" 2>&1 || true)
	case $logs in
	*"boot success!"*) break ;;
	*"boot failed!"*)
		echo "$logs" | tail -30 >&2
		exit 1
		;;
	esac
	sleep 5
done
case $(docker logs "$name" 2>&1) in
*"boot success!"*) ;;
*)
	echo "OceanBase did not boot in 15 minutes" >&2
	exit 1
	;;
esac
port=$(docker port "$name" 2881/tcp | head -1 | sed 's/.*://')
# External tables read files only under secure_file_priv, which the server
# lets a client set only over its unix socket.
timeout 30 docker exec "$name" timeout 20 obclient -S /root/ob/observer/run/sql.sock -uroot@test \
	-e "SET GLOBAL secure_file_priv = '/'" || echo "secure_file_priv not set: external table cases will fail" >&2
echo "$(date +%T) $(docker exec "$name" obclient -h127.1 -uroot@test -P2881 -N -e 'SELECT VERSION()' 2>/dev/null) on port $port"

# The tenant is test, and its root user has no password.
export OCEANBASE_WALK_DSN="root@test:@tcp(127.0.0.1:$port)/"
"$@"
