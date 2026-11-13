package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mariadb/parser"
	"github.com/bytebase/omni/v1review"
)

func TestEngineMySQLUseInnoDB(t *testing.T) {
	for _, tt := range []struct {
		name string
		sql  string
		want []v1review.Finding // Start and End are checked to be the statement's
	}{
		{name: "create innodb in any case", sql: "CREATE TABLE t (a INT) ENGINE = InnoDb;"},
		{name: "create without engine", sql: "CREATE TABLE t (a INT);"},
		{
			name: "create with another engine carries the search text",
			sql:  "CREATE TABLE t (a INT) ENGINE = MyISAM;",
			want: []v1review.Finding{{Code: v1review.EngineNotInnoDB, Name: "MyISAM", Search: "MyISAM"}},
		},
		{name: "create judges its first engine option only", sql: "CREATE TABLE t (a INT) ENGINE = InnoDB ENGINE = CSV;"},
		{
			name: "alter with another engine",
			sql:  "ALTER TABLE t ADD COLUMN b INT, ENGINE = CSV, ENGINE = MEMORY;",
			want: []v1review.Finding{{Code: v1review.EngineNotInnoDB, Name: "CSV"}},
		},
		{name: "alter innodb", sql: "ALTER TABLE t ENGINE = INNODB;"},
		{
			name: "set a bare word",
			sql:  "SET default_storage_engine = CSV;",
			want: []v1review.Finding{{Code: v1review.EngineNotInnoDB, Name: "CSV"}},
		},
		{
			name: "set a string, any scope word and letter case",
			sql:  "SET SESSION DEFAULT_STORAGE_ENGINE = 'MyISAM';",
			want: []v1review.Finding{{Code: v1review.EngineNotInnoDB, Name: "MyISAM"}},
		},
		{
			name: "set reports once",
			sql:  "SET default_storage_engine = INNODB, default_storage_engine = CSV, GLOBAL default_storage_engine = MEMORY;",
			want: []v1review.Finding{{Code: v1review.EngineNotInnoDB, Name: "CSV"}},
		},
		{name: "set innodb", sql: "SET default_storage_engine = 'innodb';"},
		{name: "set by the @@ spelling is not judged", sql: "SET @@default_storage_engine = CSV, @@session.default_storage_engine = 'CSV';"},
		{name: "set another variable", sql: "SET default_tmp_storage_engine = MyISAM;"},
		{name: "a routine body is not judged", sql: "CREATE PROCEDURE p() BEGIN CREATE TABLE t (a INT) ENGINE = CSV; END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			list, err := parser.Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			got := EngineMySQLUseInnoDB(list.Items)
			for i := range got {
				if text := tt.sql[got[i].Start:got[i].End]; text+";" != tt.sql {
					t.Errorf("finding %d anchors %q, want the statement", i, text)
				}
				got[i].Start, got[i].End = 0, 0
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
