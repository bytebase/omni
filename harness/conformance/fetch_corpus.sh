#!/usr/bin/env bash
# Fetches upstream engine test corpora at pinned tags into harness/conformance/corpus/
# (gitignored — fetch-don't-vendor: mtr is GPLv2; tidb/starrocks are Apache-2.0, but we
# vendor none of them). Sparse + blob-filtered: only the test dirs are materialized.
set -euo pipefail
cd "$(dirname "$0")"
mkdir -p corpus

TIDB_TAG="v8.5.5"

if [ "$(git -C corpus/tidb describe --tags 2>/dev/null || true)" != "$TIDB_TAG" ] ||
   [ ! -f corpus/tidb/pkg/parser/parser_test.go ]; then
  rm -rf corpus/tidb
  git clone --depth 1 --branch "$TIDB_TAG" --filter=blob:none --sparse \
    https://github.com/pingcap/tidb.git corpus/tidb
  git -C corpus/tidb sparse-checkout set pkg/parser
else
  echo "corpus/tidb already present ($(git -C corpus/tidb describe --tags))"
fi
echo "corpus/tidb at $(git -C corpus/tidb rev-parse HEAD)"

echo "corpus ready:"
ls corpus/tidb/pkg/parser/*_test.go
