package catalog

import "testing"

// TestOracle_OceanBaseTTL holds checkTTL to the engine: a TTL's columns exist and are
// DATETIME or TIMESTAMP, entry by entry.
func TestOracle_OceanBaseTTL(t *testing.T) {
	runExecSteps(t, "oceanbase_catalog_ttl", []struct {
		name  string
		steps []engineStep
	}{
		{"columns", []engineStep{
			{"CREATE TABLE a (b INT PRIMARY KEY, c INT, dt DATETIME, ts TIMESTAMP, d DATE)", taken},
			{"ALTER TABLE a TTL (c + INTERVAL 1 DAY)", 10503},
			{"ALTER TABLE a TTL (d + INTERVAL 1 DAY)", 10503},
			{"ALTER TABLE a TTL (zz + INTERVAL 1 DAY)", 10502},
			{"ALTER TABLE a TTL (c + INTERVAL 1 DAY, zz + INTERVAL 1 DAY)", 10503},
			{"ALTER TABLE a TTL (zz + INTERVAL 1 DAY, c + INTERVAL 1 DAY)", 10502},
			{"ALTER TABLE a TTL (dt + INTERVAL 1 DAY, ts + INTERVAL 1 SECOND)", taken},
			{"CREATE TABLE t2 (b INT PRIMARY KEY, c INT) TTL (c + INTERVAL 1 DAY)", 10503},
			{"CREATE TABLE t2 (b INT PRIMARY KEY, c DATETIME) TTL (zz + INTERVAL 1 DAY)", 10502},
			{"CREATE TABLE t2 (b INT PRIMARY KEY, c DATETIME) TTL (c + INTERVAL 1 DAY)", taken},
		}},
	})
}
