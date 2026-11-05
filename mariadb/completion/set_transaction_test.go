package completion

import (
	"slices"
	"testing"

	"github.com/bytebase/omni/mariadb/catalog"
)

// TestSetTransactionCompletion checks what SET TRANSACTION offers at each word. After a
// comma it offers only the kind the list does not hold yet.
func TestSetTransactionCompletion(t *testing.T) {
	cat := catalog.New()
	for _, tc := range []struct {
		sql  string
		want []string
	}{
		{"set transaction <|>", []string{"ISOLATION", "READ"}},
		{"set session transaction <|>", []string{"ISOLATION", "READ"}},
		{"set local transaction <|>", []string{"ISOLATION", "READ"}},
		{"set transaction read only, <|>", []string{"ISOLATION"}},
		{"set transaction isolation level serializable, <|>", []string{"READ"}},
		{"set transaction read <|>", []string{"ONLY", "WRITE"}},
		{"set transaction isolation level <|>", []string{"READ", "REPEATABLE", "SERIALIZABLE"}},
		{"set transaction isolation level read <|>", []string{"COMMITTED", "UNCOMMITTED"}},
		{"set transaction read only, isolation level read <|>", []string{"COMMITTED", "UNCOMMITTED"}},
		{"set transaction isolation level serializable, read <|>", []string{"ONLY", "WRITE"}},
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
