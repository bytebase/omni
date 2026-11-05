package completion

import (
	"slices"
	"testing"

	"github.com/bytebase/omni/tidb/catalog"
)

// TestSetTransactionCompletion checks what SET TRANSACTION offers at each word.
func TestSetTransactionCompletion(t *testing.T) {
	cat := catalog.New()
	for _, tc := range []struct {
		sql  string
		want []string
	}{
		{"set transaction <|>", []string{"ISOLATION", "READ"}},
		{"set session transaction <|>", []string{"ISOLATION", "READ"}},
		{"set transaction read only, <|>", []string{"ISOLATION", "READ"}},
		{"set transaction read <|>", []string{"ONLY", "WRITE"}},
		{"set transaction isolation level <|>", []string{"READ", "REPEATABLE", "SERIALIZABLE"}},
		{"set transaction isolation level read <|>", []string{"COMMITTED", "UNCOMMITTED"}},
		{"set transaction read only as of <|>", []string{"TIMESTAMP"}},
		{"set transaction read only as of timestamp now(), <|>", []string{"ISOLATION", "READ"}},
		// TRANSACTION is a variable name only once `=`, `:=` or `.` is written after it, so
		// text past the cursor does not take these away.
		{"set global transaction <|>", []string{"ISOLATION", "READ"}},
		{"set transaction <|> = 1", []string{"ISOLATION", "READ"}},
		{"set session transaction <|> := 1", []string{"ISOLATION", "READ"}},
		{"set transaction <|>.a = 1", []string{"ISOLATION", "READ"}},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			sql, cursor := splitCursor(t, tc.sql)
			var got []string
			for _, c := range Complete(sql, cursor, cat) {
				got = append(got, c.Text)
			}
			slices.Sort(got)
			if !slices.Equal(got, tc.want) {
				t.Errorf("candidates = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestSetListCompletion checks that a later assignment offers what the first one does:
// the list holds the same assignments at every place.
func TestSetListCompletion(t *testing.T) {
	cat := catalog.New()
	for _, tc := range []struct {
		sql      string
		want     []string
		wantNone []string
	}{
		{"set <|>", []string{"GLOBAL", "SESSION", "CHARACTER", "@@autocommit"}, []string{"SELECT"}},
		{"set a = 1, <|>", []string{"GLOBAL", "SESSION", "CHARACTER", "@@autocommit"}, []string{"SELECT"}},
		{"set a = 1, global <|>", []string{"@@autocommit"}, []string{"GLOBAL", "SELECT"}},
		{"set global a = 1, <|>", []string{"SESSION", "@@autocommit"}, []string{"SELECT"}},
		{"set a = 1, names <|>", []string{"utf8mb4"}, []string{"SELECT"}},
		{"set names <|>", []string{"utf8mb4"}, []string{"SELECT"}},
		{"set a = 1, charset <|>", []string{"utf8mb4"}, []string{"SELECT"}},
		{"set a = 1, character set <|>", []string{"utf8mb4"}, []string{"SELECT"}},
		{"set char set <|>", []string{"utf8mb4"}, []string{"SELECT"}},
		{"set names utf8, <|>", []string{"GLOBAL", "@@autocommit"}, []string{"SELECT"}},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			sql, cursor := splitCursor(t, tc.sql)
			var got []string
			for _, c := range Complete(sql, cursor, cat) {
				got = append(got, c.Text)
			}
			for _, w := range tc.want {
				if !slices.Contains(got, w) {
					t.Errorf("candidates %v lack %q", got, w)
				}
			}
			for _, w := range tc.wantNone {
				if slices.Contains(got, w) {
					t.Errorf("candidates %v hold %q", got, w)
				}
			}
		})
	}
}
