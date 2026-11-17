package driver

import (
	"bytes"
	"context"
	"database/sql"
	"slices"
	"testing"
	"time"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/tidb/internal/testcontainer"
)

// columnKindsSetup holds one row of every column kind the client names and the
// decode path decides: the BLOB family, the two binary strings, and an
// unsigned integer of each width.
const columnKindsSetup = `CREATE TABLE qkinds (
  id BIGINT PRIMARY KEY,
  body BLOB,
  tiny TINYBLOB,
  medium MEDIUMBLOB,
  large LONGBLOB,
  raw VARBINARY(8),
  fixedb BINARY(4),
  ubig BIGINT UNSIGNED,
  uintc INT UNSIGNED,
  umedium MEDIUMINT UNSIGNED,
  usmall SMALLINT UNSIGNED,
  utiny TINYINT UNSIGNED
);
INSERT INTO qkinds VALUES
  (1, 0x00ff, 0x01, 0x02, 0x03, 0x10, 0xdeadbeef, 18446744073709551615, 4294967295, 16777215, 65535, 255),
  (2, '', '', '', '', '', '', 0, 0, 0, 0, 0),
  (3, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);`

// TestQueryColumnKinds runs the decode path over the column kinds TiDB names
// for a binary payload and an unsigned integer. Each keeps the kind the server
// states instead of decoding as a string, and an empty []byte crosses as Bytes
// rather than NULL. TiDB has no spatial types, so BINARY and VARBINARY are the
// binary strings it names.
func TestQueryColumnKinds(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t, "omni_qe_kinds")
	sess, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sess.Close() }()

	if _, err := sess.Query(ctx, columnKindsSetup, driver.QueryOptions{}); err != nil {
		t.Fatalf("create qkinds: %v", err)
	}
	kinds, err := sess.Query(ctx, "SELECT * FROM qkinds ORDER BY id;", driver.QueryOptions{})
	if err != nil || len(kinds) != 1 || kinds[0].Err != nil {
		t.Fatalf("qkinds query = %+v, %v", kinds, err)
	}
	columns := kinds[0].Columns
	var types []string
	for _, c := range columns {
		types = append(types, c.TypeName)
	}
	wantTypes := []string{"BIGINT", "BLOB", "BLOB", "BLOB", "BLOB", "VARBINARY", "BINARY",
		"UNSIGNED BIGINT", "UNSIGNED INT", "UNSIGNED MEDIUMINT", "UNSIGNED SMALLINT", "UNSIGNED TINYINT"}
	if !slices.Equal(types, wantTypes) {
		t.Fatalf("qkinds types = %v, want %v", types, wantTypes)
	}
	if len(kinds[0].Rows) != 3 {
		t.Fatalf("qkinds rows = %d, want 3", len(kinds[0].Rows))
	}

	filled, blank, nulls := kinds[0].Rows[0], kinds[0].Rows[1], kinds[0].Rows[2]
	for i, want := range [][]byte{{0x00, 0xff}, {0x01}, {0x02}, {0x03}, {0x10}, {0xde, 0xad, 0xbe, 0xef}} {
		if got := filled[i+1]; got.Kind != driver.Bytes || !bytes.Equal(got.Bytes, want) {
			t.Errorf("filled %q = %+v, want bytes %x", columns[i+1].Name, got, want)
		}
	}
	for i, want := range []uint64{18446744073709551615, 4294967295, 16777215, 65535, 255} {
		if got := filled[i+7]; got.Kind != driver.Uint64 || got.Uint64 != want {
			t.Errorf("filled %q = %+v, want uint64 %d", columns[i+7].Name, got, want)
		}
	}
	for i := 1; i < len(columns); i++ {
		switch got := blank[i]; {
		case i == 6:
			// BINARY(4) pads a short value with zero bytes.
			if got.Kind != driver.Bytes || !bytes.Equal(got.Bytes, []byte{0, 0, 0, 0}) {
				t.Errorf("blank %q = %+v, want four zero bytes", columns[i].Name, got)
			}
		case i < 6:
			if got.Kind != driver.Bytes || len(got.Bytes) != 0 {
				t.Errorf("blank %q = %+v, want empty bytes", columns[i].Name, got)
			}
		default:
			if got.Kind != driver.Uint64 || got.Uint64 != 0 {
				t.Errorf("blank %q = %+v, want uint64 0", columns[i].Name, got)
			}
		}
		if got := nulls[i]; got.Kind != driver.Null {
			t.Errorf("null %q = %+v, want NULL", columns[i].Name, got)
		}
	}
}

// execScript runs script on sess and fails on the first statement the engine
// refused: Query returns a statement's error on its result, not from itself.
func execScript(t *testing.T, sess driver.Session, script string) {
	t.Helper()
	results, err := sess.Query(context.Background(), script, driver.QueryOptions{})
	if err != nil {
		t.Fatalf("script: %v", err)
	}
	for _, res := range results {
		if res.Err != nil {
			t.Fatalf("%s: %v", res.Statement, res.Err)
		}
	}
}

// temporalKindsSetup holds the temporal shapes a TiDB column may declare:
// a TIMESTAMP and a DATETIME at the highest precision, one at the default
// precision, the DATE the decode path reads as text, and the NULL of each.
const temporalKindsSetup = `DROP TABLE IF EXISTS qtimes;
CREATE TABLE qtimes (
  id BIGINT PRIMARY KEY,
  ts TIMESTAMP(6) NULL,
  dt DATETIME(6) NULL,
  dt0 DATETIME NULL,
  d DATE NULL
);
INSERT INTO qtimes VALUES
  (1, '2026-03-04 05:06:07.123456', '2026-03-04 05:06:07.123456', '2026-03-04 05:06:07', '2026-03-04'),
  (2, NULL, NULL, NULL, NULL);`

// TestQueryTemporalKinds runs the decode path over those shapes on a
// connection whose client parses time, in the default location and eight hours
// east: a TIMESTAMP or DATETIME answers as Timestamp with the column's
// precision and the wall clock the column holds under either location, as
// mysql/driver answers it.
func TestQueryTemporalKinds(t *testing.T) {
	ctx := context.Background()
	const dbName = "omni_qe_temporal"
	testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
		for _, stmt := range []string{
			"DROP DATABASE IF EXISTS " + quoteBacktick(dbName),
			"CREATE DATABASE " + quoteBacktick(dbName),
		} {
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				return err
			}
		}
		return nil
	})
	for _, tc := range []struct {
		name string
		loc  *time.Location
	}{
		{"default-loc", time.UTC},
		{"loc-plus8", time.FixedZone("UTC+8", 8*3600)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db := testcontainer.OpenWith(t, testcontainer.Options{
				Database:  dbName,
				ParseTime: true,
				Loc:       tc.loc,
			})
			t.Cleanup(func() { _ = db.Close() })
			d := &Driver{db: db, databaseName: dbName}
			sess, err := d.Session(ctx)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = sess.Close() }()

			execScript(t, sess, temporalKindsSetup)
			times, err := sess.Query(ctx, "SELECT * FROM qtimes ORDER BY id;", driver.QueryOptions{})
			if err != nil || len(times) != 1 || times[0].Err != nil {
				t.Fatalf("qtimes query = %+v, %v", times, err)
			}
			columns := times[0].Columns
			var types []string
			for _, c := range columns {
				types = append(types, c.TypeName)
			}
			wantTypes := []string{"BIGINT", "TIMESTAMP", "DATETIME", "DATETIME", "DATE"}
			if !slices.Equal(types, wantTypes) {
				t.Fatalf("qtimes types = %v, want %v", types, wantTypes)
			}

			filled, nulls := times[0].Rows[0], times[0].Rows[1]
			for i, want := range []struct {
				precision int32
				nano      int
			}{
				{6, 123456000},
				{6, 123456000},
				{0, 0},
			} {
				wantTime := time.Date(2026, 3, 4, 5, 6, 7, want.nano, time.UTC)
				if got := filled[i+1]; got.Kind != driver.Timestamp || got.Precision != want.precision || !got.Time.Equal(wantTime) {
					t.Errorf("filled %q = %+v, want Timestamp %s at precision %d", columns[i+1].Name, got, wantTime, want.precision)
				}
			}
			wantDate := time.Date(2026, 3, 4, 0, 0, 0, 0, tc.loc).Format(time.RFC3339Nano)
			if got := filled[4]; got.Kind != driver.String || got.String != wantDate {
				t.Errorf("filled %q = %+v, want String %q", columns[4].Name, got, wantDate)
			}
			for i := 1; i < len(columns); i++ {
				if got := nulls[i]; got.Kind != driver.Null {
					t.Errorf("null %q = %+v, want NULL", columns[i].Name, got)
				}
			}
		})
	}
}
