#!/usr/bin/env bash
# Starts the pinned TiDB oracle for adjudication. Stop with: docker rm -f tidb-conformance
set -euo pipefail
docker rm -f tidb-conformance 2>/dev/null || true
# Loopback bind: the oracle runs passwordless root, so it must not listen on
# the LAN. The DSN below already targets 127.0.0.1.
docker run -d --name tidb-conformance -p 127.0.0.1:14001:4000 pingcap/tidb:v8.5.5
echo "waiting for TiDB to listen..."
for i in $(seq 1 60); do
  if nc -z 127.0.0.1 14001 2>/dev/null; then break; fi
  sleep 1
done
nc -z 127.0.0.1 14001 || { echo "TiDB did not listen within 60s" >&2; exit 1; }
DIGEST=$(docker inspect --format '{{index .RepoDigests 0}}' pingcap/tidb:v8.5.5 2>/dev/null || echo "")
echo "export TIDB_CONTAINER_DIGEST=$DIGEST"
# No default schema in the DSN: a default schema is droppable by adjudicated
# DDL, after which every fresh connection fails its handshake with 1049.
echo "export TIDB_DSN='root@tcp(127.0.0.1:14001)/?multiStatements=true&timeout=5s&readTimeout=10s&writeTimeout=10s'"
