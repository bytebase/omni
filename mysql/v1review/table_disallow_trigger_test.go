package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestTableDisallowTrigger(t *testing.T) {
	for _, tt := range []struct {
		name string
		sql  string
		want []finding
	}{
		{
			name: "create trigger",
			sql:  "CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;",
			want: []finding{{v1review.TableTriggerCreated, "", "CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1"}},
		},
		{name: "drop trigger not judged", sql: "DROP TRIGGER tr;"},
		{name: "create table not judged", sql: "CREATE TABLE t (a INT);"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, tableFindings(t, tt.sql, TableDisallowTrigger)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
