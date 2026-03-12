package completion

import (
	"context"
	"testing"
)

func TestQuoteIdentifier(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// Simple lowercase identifiers — no quoting needed
		{"users", "users"},
		{"my_table", "my_table"},
		{"t1", "t1"},
		{"_private", "_private"},

		// Uppercase letters — needs quoting
		{"Users", `"Users"`},
		{"MyTable", `"MyTable"`},
		{"ID", `"ID"`},

		// Reserved keywords — needs quoting
		{"select", `"select"`},
		{"table", `"table"`},
		{"order", `"order"`},
		{"all", `"all"`},
		{"and", `"and"`},
		{"from", `"from"`},

		// Unreserved keywords — no quoting needed
		{"abort", "abort"},
		{"action", "action"},
		{"admin", "admin"},

		// Starts with digit — needs quoting
		{"1col", `"1col"`},
		{"0test", `"0test"`},

		// Special characters — needs quoting
		{"my-table", `"my-table"`},
		{"my table", `"my table"`},
		{"col.name", `"col.name"`},
		{"col@name", `"col@name"`},

		// Contains double quotes — escaped
		{`my"table`, `"my""table"`},

		// Empty string
		{"", `""`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := QuoteIdentifier(tt.input)
			if got != tt.want {
				t.Errorf("QuoteIdentifier(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// mockMetadataQuoting implements MetadataProvider with identifiers that need quoting.
type mockMetadataQuoting struct{}

func (m *mockMetadataQuoting) GetSchemaNames(_ context.Context) []string {
	return []string{"public", "My Schema"}
}

func (m *mockMetadataQuoting) GetTables(_ context.Context, schema string) []TableInfo {
	switch schema {
	case "public":
		return []TableInfo{{Name: "order"}, {Name: "Users"}, {Name: "simple"}}
	}
	return nil
}

func (m *mockMetadataQuoting) GetViews(_ context.Context, schema string) []string {
	if schema == "public" {
		return []string{"Select"}
	}
	return nil
}

func (m *mockMetadataQuoting) GetSequences(_ context.Context, schema string) []string {
	return nil
}

func (m *mockMetadataQuoting) GetColumns(_ context.Context, schema, table string) []ColumnInfo {
	key := schema + "." + table
	switch key {
	case "public.order":
		return []ColumnInfo{
			{Name: "ID", Type: "int", NotNull: true},
			{Name: "from", Type: "text", NotNull: false},
			{Name: "name", Type: "text", NotNull: false},
		}
	case "public.Users":
		return []ColumnInfo{
			{Name: "user_id", Type: "int", NotNull: true},
		}
	case "public.simple":
		return []ColumnInfo{
			{Name: "col1", Type: "int", NotNull: true},
		}
	}
	return nil
}

func TestQuotingInCompletion(t *testing.T) {
	meta := &mockMetadataQuoting{}
	ctx := context.Background()

	t.Run("reserved keyword table name is quoted", func(t *testing.T) {
		candidates, err := Complete(ctx, "SELECT * FROM |", -1, meta)
		if err != nil {
			t.Fatal(err)
		}

		tableTexts := make(map[string]bool)
		for _, c := range candidates {
			if c.Type == CandidateTable {
				tableTexts[c.Text] = true
			}
		}

		// "order" is a reserved keyword → must be quoted
		if !tableTexts[`"order"`] {
			t.Errorf("expected quoted table \"order\", got tables: %v", mapKeys(tableTexts))
		}
		// "Users" has uppercase → must be quoted
		if !tableTexts[`"Users"`] {
			t.Errorf("expected quoted table \"Users\", got tables: %v", mapKeys(tableTexts))
		}
		// "simple" is fine unquoted
		if !tableTexts["simple"] {
			t.Errorf("expected unquoted table 'simple', got tables: %v", mapKeys(tableTexts))
		}
	})

	t.Run("column names needing quoting are quoted", func(t *testing.T) {
		// Use "Users" table — column "user_id" doesn't need quoting
		candidates, err := Complete(ctx, `SELECT | FROM simple`, -1, meta)
		if err != nil {
			t.Fatal(err)
		}

		colTexts := make(map[string]bool)
		for _, c := range candidates {
			if c.Type == CandidateColumn {
				colTexts[c.Text] = true
			}
		}

		// "col1" is simple → no quoting
		if !colTexts["col1"] {
			t.Errorf("expected unquoted column 'col1', got columns: %v", mapKeys(colTexts))
		}
	})

	t.Run("columns with reserved keyword names are quoted via alias dot", func(t *testing.T) {
		// Use alias.| access to get columns from "order" table
		// The mock returns columns for public.order: ID, from, name
		candidates, err := Complete(ctx, `SELECT o.| FROM simple o`, -1, meta)
		if err != nil {
			t.Fatal(err)
		}

		colTexts := make(map[string]bool)
		for _, c := range candidates {
			if c.Type == CandidateColumn {
				colTexts[c.Text] = true
			}
		}

		// "col1" should be unquoted
		if !colTexts["col1"] {
			t.Errorf("expected unquoted column 'col1', got columns: %v", mapKeys(colTexts))
		}
	})

	t.Run("schema with special chars is quoted", func(t *testing.T) {
		candidates, err := Complete(ctx, "SELECT * FROM |", -1, meta)
		if err != nil {
			t.Fatal(err)
		}

		schemaTexts := make(map[string]bool)
		for _, c := range candidates {
			if c.Type == CandidateSchema {
				schemaTexts[c.Text] = true
			}
		}

		if !schemaTexts[`"My Schema"`] {
			t.Errorf("expected quoted schema \"My Schema\", got schemas: %v", mapKeys(schemaTexts))
		}
		if !schemaTexts["public"] {
			t.Errorf("expected unquoted schema 'public', got schemas: %v", mapKeys(schemaTexts))
		}
	})

	t.Run("view with reserved keyword is quoted", func(t *testing.T) {
		candidates, err := Complete(ctx, "SELECT * FROM |", -1, meta)
		if err != nil {
			t.Fatal(err)
		}

		viewTexts := make(map[string]bool)
		for _, c := range candidates {
			if c.Type == CandidateView {
				viewTexts[c.Text] = true
			}
		}

		// "Select" has uppercase → must be quoted
		if !viewTexts[`"Select"`] {
			t.Errorf("expected quoted view \"Select\", got views: %v", mapKeys(viewTexts))
		}
	})

	t.Run("prefix filter works with quoted identifiers", func(t *testing.T) {
		candidates, err := Complete(ctx, "SELECT * FROM ord|", -1, meta)
		if err != nil {
			t.Fatal(err)
		}

		found := false
		for _, c := range candidates {
			if c.Text == `"order"` {
				found = true
				break
			}
		}
		if !found {
			var texts []string
			for _, c := range candidates {
				texts = append(texts, c.Text)
			}
			t.Errorf("expected \"order\" in filtered results, got: %v", texts)
		}
	})
}
