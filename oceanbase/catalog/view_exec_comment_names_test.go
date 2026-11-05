package catalog

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

// viewExecCommentNameCases are the columns OceanBase 4.3.5.6 gives a view
// whose unaliased select items hold executable comments, read from
// information_schema.COLUMNS on 2026-11-05. OceanBase names an item after its
// text from its first token to its last exactly as written, the markers of an
// executable comment included, where MySQL 8.0 drops the markers (see
// deparse's sourceAutoName). A statement OceanBase reads again with its
// executable comments as plain ones (oceanbase/parser/exec_fallback.go) keeps
// them in the name as the comments they became.
var viewExecCommentNameCases = []struct {
	sql     string
	columns []string
}{
	{"SELECT a/*!50100 +2*/ FROM n1", []string{"a/*!50100 +2"}},
	{"SELECT a/*!+3*/ + 1 FROM n1", []string{"a/*!+3*/ + 1"}},
	{"SELECT /*!1234 */ + 1 FROM n1", []string{"1234 */ + 1"}},
	// A version is five digits, so a sixth, or a run shorter than five, is the
	// body's first token, and the item's text runs to it.
	{"SELECT 2 * /*!100005 */ FROM n1", []string{"2 * /*!100005"}},
	{"SELECT 1 + /*!1234567 */ FROM n1", []string{"1 + /*!1234567"}},
	{"SELECT 2 * /*!5 */ FROM n1", []string{"2 * /*!5"}},
	{"SELECT 1 + /*!12345 2 */ FROM n1", []string{"1 + /*!12345 2"}},
	{"SELECT 1 + /*!99999 2 */ FROM n1", []string{"1 + /*!99999 2"}},
	// Read again with the comments as plain ones.
	{"SELECT 1 + 1 /*! from */", []string{"1 + 1"}},
	{"SELECT 1 /*! + 1 */ + 1 /*! from */", []string{"1 /*! + 1 */ + 1"}},
	{"SELECT 1 /*! + 10 */ + 100 /*! , 2 */ /*! from */", []string{"1 /*! + 10 */ + 100"}},
	{"SELECT a + /*! 1 */ 2, b FROM n1 /*! where */", []string{"a + /*! 1 */ 2", "b"}},
}

// TestViewExecCommentNames pins the catalog to viewExecCommentNameCases.
func TestViewExecCommentNames(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE test; USE test; CREATE TABLE n1 (a INT, b INT);", nil); err != nil {
		t.Fatal(err)
	}
	for i, tc := range viewExecCommentNameCases {
		name := fmt.Sprintf("v%d", i)
		results, err := c.Exec("CREATE VIEW "+name+" AS "+tc.sql+";", nil)
		if err != nil {
			t.Fatalf("%s: %v", tc.sql, err)
		}
		for _, r := range results {
			if r.Error != nil {
				t.Fatalf("%s: %v", tc.sql, r.Error)
			}
		}
		if got := c.GetDatabase("test").Views[name].Columns; !reflect.DeepEqual(got, tc.columns) {
			t.Errorf("%s: columns %q, want %q", tc.sql, got, tc.columns)
		}
	}
}

// TestViewExecCommentNames_OceanBase creates every viewExecCommentNameCases
// view on a live OceanBase 4.3.5 (DSN-gated, see oceanBaseOracle) and checks
// the columns it lists in information_schema against the case.
func TestViewExecCommentNames_OceanBase(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx := context.Background()
	for _, stmt := range []string{
		"DROP DATABASE IF EXISTS ob_view_exec_names",
		"CREATE DATABASE ob_view_exec_names",
		"USE ob_view_exec_names",
		"CREATE TABLE n1 (a INT, b INT)",
	} {
		if _, err := o.db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	defer o.db.ExecContext(ctx, "DROP DATABASE IF EXISTS ob_view_exec_names")
	for i, tc := range viewExecCommentNameCases {
		name := fmt.Sprintf("v%d", i)
		if _, err := o.db.ExecContext(ctx, "CREATE VIEW "+name+" AS "+tc.sql); err != nil {
			t.Errorf("%s: OceanBase refuses it: %v", tc.sql, err)
			continue
		}
		rows, err := o.db.QueryContext(ctx, "SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'ob_view_exec_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION", name)
		if err != nil {
			t.Fatal(err)
		}
		var got []string
		for rows.Next() {
			var col string
			if err := rows.Scan(&col); err != nil {
				t.Fatal(err)
			}
			got = append(got, col)
		}
		if err := rows.Close(); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, tc.columns) {
			t.Errorf("%s: OceanBase names the columns %q, the case says %q", tc.sql, got, tc.columns)
		}
	}
}
