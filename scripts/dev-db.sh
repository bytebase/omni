#!/usr/bin/env bash
# Keeps one long-lived server per engine on this machine, for every session
# that wants to ask an engine something, so none of them pays for a boot. The
# first call starts the server; later calls, from any checkout, find it.
#
#   scripts/dev-db.sh sql mysql 'PREPARE p FROM "SELECT 1"'   # run SQL
#   scripts/dev-db.sh sql mariadb < probe.sql                 # or from stdin
#   scripts/dev-db.sh dsn postgres17                          # print a DSN
#   scripts/dev-db.sh up tidb                                 # start, print nothing
#   scripts/dev-db.sh list                                    # what is running
#
# Engines: mysql (8.0), mariadb (11.8.8), tidb (v8.5.5), postgres17, postgres18.
# OceanBase needs 8 GB and stays with scripts/oceanbase-oracle.sh.
#
# `sql` and `dsn` work in a database named after this checkout, created on first
# use, so sessions do not see each other's tables. The account is the server's
# root, and every session on the machine shares the server: a probe that changes
# something server-wide (SET GLOBAL, an account, a role, SHUTDOWN) or needs an
# init-time setting belongs on a server of its own.
#
# There is deliberately no `down`: the server may be serving another session. A
# server nobody needs is removed by hand, after checking who is connected:
#   docker rm -f omni-dev-<engine>
#
# internal/sharedserver finds these servers on the same ports, so `go test` on a
# package that uses it needs no container once the engine's server is up.
set -euo pipefail

usage() {
	sed -n 2,/^set/p "$0" | sed -e '/^set/d' -e 's/^# \{0,1\}//' >&2
	exit 2
}

engine_config() {
	case "$1" in
	mysql) image=mysql:8.0 port=13380 inner=3306 family=mysql env=(-e MYSQL_ROOT_PASSWORD=test) args=(--skip-log-bin --max-connections=1000) ;;
	mariadb) image=mariadb:11.8.8 port=13311 inner=3306 family=mariadb env=(-e MARIADB_ROOT_PASSWORD=test) args=(--max-connections=1000) ;;
	tidb) image=pingcap/tidb:v8.5.5 port=14055 inner=4000 family=tidb env=() args=() ;;
	postgres17) image=postgres:17-alpine port=15417 inner=5432 family=postgres env=(-e POSTGRES_PASSWORD=test) args=(-c max_connections=500) ;;
	postgres18) image=postgres:18-alpine port=15418 inner=5432 family=postgres env=(-e POSTGRES_PASSWORD=test) args=(-c max_connections=500) ;;
	*) echo "unknown engine: $1" >&2; usage ;;
	esac
	name=omni-dev-$1
}

# client runs the engine's client on its server, in a container already running.
# tidb's image has no client, so it borrows the mysql server's.
client() { # engine database [client args...]
	local engine=$1 database=$2
	shift 2
	engine_config "$engine"
	case "$family" in
	mysql) docker exec -i "$name" mysql -uroot -ptest --protocol=tcp --silent ${database:+-D "$database"} "$@" 2> >(grep -v 'Using a password' >&2) ;;
	mariadb) docker exec -i "$name" mariadb -uroot -ptest --protocol=tcp --silent ${database:+-D "$database"} "$@" ;;
	tidb)
		local ip
		ip=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' "$name")
		docker exec -i omni-dev-mysql mysql -h"$ip" -P4000 -uroot --protocol=tcp --silent ${database:+-D "$database"} "$@"
		;;
	postgres) docker exec -i "$name" psql -U postgres -h 127.0.0.1 -v ON_ERROR_STOP=1 -X ${database:+-d "$database"} "$@" ;;
	esac
}

ready() { # engine
	case "$1" in
	postgres*) engine_config "$1"; docker exec "$name" pg_isready -q -h 127.0.0.1 -U postgres ;;
	*) client "$1" "" -e 'SELECT 1' >/dev/null 2>&1 ;;
	esac
}

up() { # engine
	engine_config "$1"
	local engine=$1 have
	[ "$engine" != tidb ] || up mysql
	engine_config "$engine"
	have=$(docker inspect -f '{{.Config.Image}} {{.State.Running}}' "$name" 2>/dev/null || true)
	case "$have" in
	"") # Two sessions can race here; docker run lets one win and the other finds it.
		local out
		out=$(docker run -d --name "$name" --restart=no "${env[@]}" -p "127.0.0.1:$port:$inner" \
			--label omni.dev-db=true "$image" "${args[@]}" 2>&1) || case "$out" in
		*"is already in use"*) ;;
		*)
			echo "$out" >&2
			exit 1
			;;
		esac
		;;
	"$image true") ;;
	"$image false") docker start "$name" >/dev/null 2>&1 || true ;;
	*)
		echo "$name runs ${have% *}, and this script pins $image; compare them and remove it by hand" >&2
		exit 1
		;;
	esac
	local deadline=$((SECONDS + 180))
	until ready "$engine"; do
		if [ "$SECONDS" -ge "$deadline" ]; then
			echo "$name did not become ready; see docker logs $name" >&2
			exit 1
		fi
		sleep 1
	done
}

# scratch names this checkout's database, from the directory it is in.
scratch() {
	local top
	top=$(git rev-parse --show-toplevel 2>/dev/null || pwd)
	top=$(basename "$top" | tr 'A-Z' 'a-z' | tr -c 'a-z0-9\n' '_')
	echo "scratch_${top:0:48}"
}

ensure_scratch() { # engine
	engine_config "$1"
	if [ "$family" = postgres ]; then
		client "$1" "" -tAc "SELECT 1 FROM pg_database WHERE datname = '$(scratch)'" | grep -q 1 ||
			client "$1" "" -qc "CREATE DATABASE $(scratch)" 2>/dev/null || true
	else
		client "$1" "" -e "CREATE DATABASE IF NOT EXISTS $(scratch)"
	fi
}

[ $# -ge 1 ] || usage
command=$1
shift
case "$command" in
up)
	[ $# -eq 1 ] || usage
	up "$1"
	;;
dsn)
	[ $# -eq 1 ] || usage
	up "$1"
	ensure_scratch "$1"
	engine_config "$1"
	case "$family" in
	mysql | mariadb) echo "root:test@tcp(127.0.0.1:$port)/$(scratch)" ;;
	tidb) echo "root@tcp(127.0.0.1:$port)/$(scratch)" ;;
	postgres) echo "postgres://postgres:test@127.0.0.1:$port/$(scratch)?sslmode=disable" ;;
	esac
	;;
sql)
	[ $# -ge 1 ] || usage
	engine=$1
	shift
	up "$engine"
	ensure_scratch "$engine"
	engine_config "$engine"
	if [ $# -eq 0 ]; then
		client "$engine" "$(scratch)"
	elif [ "$family" = postgres ]; then
		client "$engine" "$(scratch)" -c "$*"
	else
		client "$engine" "$(scratch)" -e "$*"
	fi
	;;
list)
	docker ps -a --filter label=omni.dev-db=true --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}'
	;;
*) usage ;;
esac
