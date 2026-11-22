package analysis

import "testing"

func TestWalkThroughStatementKindString(t *testing.T) {
	tests := []struct {
		kind WalkThroughStatementKind
		want string
	}{
		{WalkThroughStatementOther, "OTHER"},
		{WalkThroughStatementCreateDatabase, "CREATE_DATABASE"},
		{WalkThroughStatementAlterDatabase, "ALTER_DATABASE"},
		{WalkThroughStatementDropDatabase, "DROP_DATABASE"},
		{WalkThroughStatementCreateTable, "CREATE_TABLE"},
		{WalkThroughStatementDropTable, "DROP_TABLE"},
		{WalkThroughStatementTableRefCheck, "TABLE_REF_CHECK"},
		{WalkThroughStatementRenameTable, "RENAME_TABLE"},
		{WalkThroughStatementKind(99), "OTHER"},
	}

	for _, tt := range tests {
		if got := tt.kind.String(); got != tt.want {
			t.Errorf("WalkThroughStatementKind(%d).String() = %q, want %q", tt.kind, got, tt.want)
		}
	}
}
