#!/usr/bin/env bash
# Starts the pinned TiDB oracle for adjudication. Stop with: docker rm -f tidb-conformance
set -euo pipefail
docker rm -f tidb-conformance 2>/dev/null || true
docker run -d --name tidb-conformance -p 14001:4000 pingcap/tidb:v8.5.5
echo "waiting for TiDB to listen..."
for i in $(seq 1 60); do
  if nc -z 127.0.0.1 14001 2>/dev/null; then break; fi
  sleep 1
done
nc -z 127.0.0.1 14001 || { echo "TiDB did not listen within 60s" >&2; exit 1; }
DIGEST=$(docker inspect --format '{{index .RepoDigests 0}}' pingcap/tidb:v8.5.5 2>/dev/null || echo "")
echo "export TIDB_CONTAINER_DIGEST=$DIGEST"
echo "export TIDB_DSN='root@tcp(127.0.0.1:14001)/test?multiStatements=true&timeout=5s&readTimeout=10s&writeTimeout=10s'"
