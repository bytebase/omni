package driver

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/bytebase/omni/clickhouse/deparse"
	"github.com/bytebase/omni/clickhouse/internal/testcontainer"
	"github.com/bytebase/omni/driver"
)

// TestDeparseRoundTrip creates ClickHouse tables and views on the pinned server,
// syncs them via SyncDatabase, renders the snapshot with clickhouse/deparse,
// drops and recreates the database, executes the deparsed DDL, and requires
// SyncDatabase on the restored database to match the source snapshot.
func TestDeparseRoundTrip(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, `
CREATE TABLE events (
  id UInt64,
  name String COMMENT 'the ''event'' name\\v1',
  at DateTime64(3, 'UTC'),
  amount Nullable(Decimal(12, 2)),
  tags Array(String),
  attrs Map(String, String),
  level Enum8('info' = 1, 'warn' = 2) DEFAULT 'info',
  seen Date DEFAULT '2026-01-02',
  retries UInt32 DEFAULT 3
)
ENGINE = MergeTree
ORDER BY (id, name)
PRIMARY KEY id
COMMENT 'events ''table''';

CREATE TABLE `+"`odd-name space`"+` (
  `+"`select`"+` Int32 DEFAULT 0,
  `+"`note col`"+` Nullable(String)
)
ENGINE = Memory
COMMENT 'line1\nline2';

CREATE TABLE scratch (
  k Int32
)
ENGINE = MergeTree
ORDER BY tuple();

CREATE VIEW event_names AS SELECT id, name FROM events COMMENT 'active ''view''';
`)

	source, err := d.SyncDatabase(ctx)
	require.NoError(t, err)

	ddl, err := deparse.DatabaseDefinition(source)
	require.NoError(t, err)

	recreateDatabase(t, d, testcontainer.Database(t))
	_, err = d.Execute(ctx, ddl, driver.ExecuteOptions{})
	require.NoError(t, err, "executing deparsed DDL:\n%s", ddl)

	restored, err := d.SyncDatabase(ctx)
	require.NoError(t, err)

	if diff := cmp.Diff(source, restored, protocmp.Transform()); diff != "" {
		t.Errorf("restored snapshot differs from source (-source +restored):\n%s\nDDL:\n%s", diff, ddl)
	}
}
