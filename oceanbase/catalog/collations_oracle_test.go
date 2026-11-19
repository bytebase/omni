package catalog

import (
	"testing"

	"github.com/bytebase/omni/oceanbase/internal/testcontainer"
)

// TestCollationsMatchOceanBase holds oceanBaseCollations, which judgeCollationName
// reads, to the pinned server's information_schema.COLLATIONS. An entry the server
// lacks would accept a collation it refuses (1273), and a missing entry would refuse
// one it takes.
func TestCollationsMatchOceanBase(t *testing.T) {
	rows, err := testcontainer.Open(t).Query("SELECT COLLATION_NAME, CHARACTER_SET_NAME FROM information_schema.COLLATIONS")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	reported := map[string]string{}
	for rows.Next() {
		var collation, charset string
		if err := rows.Scan(&collation, &charset); err != nil {
			t.Fatal(err)
		}
		reported[collation] = charset
		if got, ok := oceanBaseCollations[collation]; !ok {
			t.Errorf("OceanBase has collation %s (%s), the catalog does not", collation, charset)
		} else if got != charset {
			t.Errorf("collation %s: catalog has charset %q, OceanBase %q", collation, got, charset)
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	for collation, charset := range oceanBaseCollations {
		if _, ok := reported[collation]; !ok {
			t.Errorf("catalog has collation %s (%s), which OceanBase does not", collation, charset)
		}
	}
}
