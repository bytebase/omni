package driver

import (
	"bytes"
	"context"
	"slices"
	"testing"

	"github.com/bytebase/omni/driver"
)

// columnKindsSetup holds one row of every column kind the client names and the
// decode path decides: the BLOB family, an unsigned integer of each width, and
// the binary strings.
const columnKindsSetup = `CREATE TABLE qkinds (
  id BIGINT PRIMARY KEY,
  body BLOB,
  tiny TINYBLOB,
  medium MEDIUMBLOB,
  large LONGBLOB,
  raw VARBINARY(8),
  ubig BIGINT UNSIGNED,
  uintc INT UNSIGNED,
  umedium MEDIUMINT UNSIGNED,
  usmall SMALLINT UNSIGNED,
  utiny TINYINT UNSIGNED
);
INSERT INTO qkinds VALUES
  (1, 0x00ff, 0x01, 0x02, 0x03, 0x10, 18446744073709551615, 4294967295, 16777215, 65535, 255),
  (2, '', '', '', '', '', 0, 0, 0, 0, 0),
  (3, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);`

// TestQueryColumnKinds runs the decode path over the column kinds OceanBase
// names for a binary payload and an unsigned integer. Each keeps the kind the
// server states instead of decoding as a string, and an empty []byte crosses
// as Bytes rather than NULL.
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
	wantTypes := []string{"BIGINT", "BLOB", "BLOB", "BLOB", "BLOB", "VARBINARY",
		"UNSIGNED BIGINT", "UNSIGNED INT", "UNSIGNED MEDIUMINT", "UNSIGNED SMALLINT", "UNSIGNED TINYINT"}
	if !slices.Equal(types, wantTypes) {
		t.Fatalf("qkinds types = %v, want %v", types, wantTypes)
	}
	if len(kinds[0].Rows) != 3 {
		t.Fatalf("qkinds rows = %d, want 3", len(kinds[0].Rows))
	}

	filled, blank, nulls := kinds[0].Rows[0], kinds[0].Rows[1], kinds[0].Rows[2]
	for i, want := range [][]byte{{0x00, 0xff}, {0x01}, {0x02}, {0x03}, {0x10}} {
		if got := filled[i+1]; got.Kind != driver.Bytes || !bytes.Equal(got.Bytes, want) {
			t.Errorf("filled %q = %+v, want bytes %x", columns[i+1].Name, got, want)
		}
	}
	for i, want := range []uint64{18446744073709551615, 4294967295, 16777215, 65535, 255} {
		if got := filled[i+6]; got.Kind != driver.Uint64 || got.Uint64 != want {
			t.Errorf("filled %q = %+v, want uint64 %d", columns[i+6].Name, got, want)
		}
	}
	for i := 1; i < len(columns); i++ {
		if got := blank[i]; i < 6 {
			if got.Kind != driver.Bytes || len(got.Bytes) != 0 {
				t.Errorf("blank %q = %+v, want empty bytes", columns[i].Name, got)
			}
		} else if got.Kind != driver.Uint64 || got.Uint64 != 0 {
			t.Errorf("blank %q = %+v, want uint64 0", columns[i].Name, got)
		}
		if got := nulls[i]; got.Kind != driver.Null {
			t.Errorf("null %q = %+v, want NULL", columns[i].Name, got)
		}
	}
}
