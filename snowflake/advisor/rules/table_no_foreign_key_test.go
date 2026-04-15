package rules_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/advisor/rules"
)

func TestTableNoForeignKey(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantCount int
	}{
		{
			name:      "no FK — no findings",
			sql:       "CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(100))",
			wantCount: 0,
		},
		{
			name:      "inline FK — one finding",
			sql:       "CREATE TABLE orders (id INT NOT NULL PRIMARY KEY, user_id INT FOREIGN KEY REFERENCES users(id))",
			wantCount: 1,
		},
		{
			name:      "table-level FK — one finding",
			sql:       "CREATE TABLE orders (id INT NOT NULL, user_id INT NOT NULL, PRIMARY KEY (id), FOREIGN KEY (user_id) REFERENCES users(id))",
			wantCount: 1,
		},
		{
			name:      "multiple inline FKs — multiple findings",
			sql:       "CREATE TABLE orders (id INT NOT NULL, user_id INT FOREIGN KEY REFERENCES users(id), product_id INT FOREIGN KEY REFERENCES products(id))",
			wantCount: 2,
		},
		{
			name:      "unique constraint — no finding",
			sql:       "CREATE TABLE t (id INT NOT NULL, email VARCHAR(100) UNIQUE)",
			wantCount: 0,
		},
		{
			name:      "SELECT statement — no findings",
			sql:       "SELECT a FROM t",
			wantCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := runRule(t, tt.sql, rules.TableNoForeignKeyRule{})
			if len(findings) != tt.wantCount {
				t.Errorf("got %d findings, want %d", len(findings), tt.wantCount)
				for _, f := range findings {
					t.Logf("  finding: %s", f.Message)
				}
			}
		})
	}
}

func TestTableNoForeignKey_Metadata(t *testing.T) {
	r := rules.TableNoForeignKeyRule{}
	if r.ID() != "snowflake.table.no-foreign-key" {
		t.Errorf("ID() = %q", r.ID())
	}
	if r.Severity() != advisor.SeverityError {
		t.Errorf("Severity() = %v, want ERROR", r.Severity())
	}
}
