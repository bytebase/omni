package redshift

import "testing"

func TestValidateSQLForEditor(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantQuery bool
		wantPlan  bool
	}{
		{
			name:      "select",
			sql:       "SELECT * FROM users;",
			wantQuery: true,
		},
		{
			name:      "plain explain select",
			sql:       "EXPLAIN SELECT * FROM users;",
			wantQuery: true,
			wantPlan:  true,
		},
		{
			name:      "show",
			sql:       "SHOW DATABASES;",
			wantQuery: true,
		},
		{
			name: "dml",
			sql:  "INSERT INTO users SELECT 1;",
		},
		{
			name: "ddl",
			sql:  "CREATE TABLE users(id INT);",
		},
		{
			name: "explain analyze",
			sql:  "EXPLAIN ANALYZE SELECT * FROM users;",
		},
		{
			name: "select into",
			sql:  "SELECT * INTO copied_users FROM users;",
		},
		{
			name:      "set then select",
			sql:       "SET search_path TO public; SELECT 1;",
			wantQuery: true,
		},
		{
			name: "mixed readonly and dml",
			sql:  "SELECT 1; DELETE FROM users WHERE id = 1;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotPlan, err := ValidateSQLForEditor(tt.sql)
			if err != nil {
				t.Fatalf("ValidateSQLForEditor returned error: %v", err)
			}
			if gotQuery != tt.wantQuery {
				t.Fatalf("query = %v, want %v", gotQuery, tt.wantQuery)
			}
			if gotPlan != tt.wantPlan {
				t.Fatalf("plan = %v, want %v", gotPlan, tt.wantPlan)
			}
		})
	}
}
