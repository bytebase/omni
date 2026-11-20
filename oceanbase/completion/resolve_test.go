package completion

import (
	"testing"

	"github.com/bytebase/omni/oceanbase/catalog"
)

// TestRemoteTableOffersNoLocalColumns pins that completing a query over a
// catalog-qualified table lists no columns of a same-named local table, and no
// fallback to every local column.
func TestRemoteTableOffersNoLocalColumns(t *testing.T) {
	cat := catalog.New()
	cat.SetCurrentDatabase("testdb")
	mustExec(t, cat, "CREATE DATABASE testdb")
	mustExec(t, cat, "CREATE TABLE orders (id INT, secret INT)")

	for _, sql := range []string{
		"SELECT  FROM ext_catalog.db.orders",
		"SELECT o. FROM ext_catalog.db.orders AS o",
	} {
		for _, c := range Complete(sql, len(sql), cat) {
			if c.Type == CandidateColumn {
				t.Errorf("%q offered local column %q", sql, c.Text)
			}
		}
	}
}
