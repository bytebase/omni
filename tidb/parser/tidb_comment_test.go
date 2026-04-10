package parser

import (
	"testing"

	"github.com/bytebase/omni/tidb/ast"
)

func TestTiDBCommentBasic(t *testing.T) {
	// /*T! ... */ — always inject inner SQL
	sql := "CREATE TABLE t (id BIGINT /*T! PRIMARY KEY */)"
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("failed to parse TiDB comment: %v", err)
	}
	if list.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", list.Len())
	}
	// Verify PRIMARY KEY was injected (column has PK constraint)
	stmt := list.Items[0].(*ast.CreateTableStmt)
	if len(stmt.Columns) == 0 {
		t.Fatal("no columns parsed")
	}
	hasPK := false
	for _, c := range stmt.Columns[0].Constraints {
		if c.Type == ast.ColConstrPrimaryKey {
			hasPK = true
		}
	}
	if !hasPK {
		t.Error("PRIMARY KEY from /*T! */ comment was not injected")
	}
}

func TestTiDBCommentFeatureSupported(t *testing.T) {
	// /*T![auto_rand] ... */ — inject because auto_rand is supported in v8.5
	sql := "SELECT /*T![auto_rand] 1 AS */ col FROM t"
	_, err := Parse(sql)
	if err != nil {
		t.Fatalf("failed to parse supported feature comment: %v", err)
	}
}

func TestTiDBCommentFeatureUnsupported(t *testing.T) {
	// /*T![unsupported_feature_xyz] ... */ — skip as comment
	sql := "SELECT /*T![unsupported_feature_xyz] WEIRD_STUFF */ 1"
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("failed to parse unsupported feature comment: %v", err)
	}
	if list.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", list.Len())
	}
}

func TestTiDBCommentMultiFeature(t *testing.T) {
	// /*T![ttl] ... */ — single supported feature
	sql := "SELECT /*T![ttl] 1 AS */ col FROM t"
	_, err := Parse(sql)
	if err != nil {
		t.Fatalf("failed to parse single-feature comment: %v", err)
	}
}

func TestTiDBCommentPreservesMySQL(t *testing.T) {
	// Standard MySQL conditional comments still work
	sql := "SELECT /*!50000 1 */ + 1"
	_, err := Parse(sql)
	if err != nil {
		t.Fatalf("MySQL conditional comment broke: %v", err)
	}
}

func TestTiDBCommentEmpty(t *testing.T) {
	// A segment with only a TiDB comment is not empty
	sql := "/*T! SELECT 1 */"
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("TiDB comment-only segment failed: %v", err)
	}
	if list.Len() != 1 {
		t.Fatalf("expected 1 statement from TiDB comment, got %d", list.Len())
	}
}

func TestAllTiDBFeaturesSupported(t *testing.T) {
	tests := []struct {
		features string
		want     bool
	}{
		{"auto_rand", true},
		{"auto_id_cache", true},
		{"clustered_index", true},
		{"placement", true},
		{"ttl", true},
		{"auto_rand_base", true},
		{"unsupported_xyz", false},
		{"auto_rand,clustered_index", true},
		{"auto_rand,unsupported", false},
		{"", true}, // empty = no features required
	}
	for _, tt := range tests {
		t.Run(tt.features, func(t *testing.T) {
			got := allTiDBFeaturesSupported(tt.features)
			if got != tt.want {
				t.Errorf("allTiDBFeaturesSupported(%q) = %v, want %v", tt.features, got, tt.want)
			}
		})
	}
}
