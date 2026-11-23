package analysis_test

import (
	"testing"

	"github.com/bytebase/omni/analysis"
)

func TestStatementTypeStringAndUniqueness(t *testing.T) {
	cases := []struct {
		typ  analysis.StatementType
		want string
	}{
		{analysis.StatementTypeUnspecified, "STATEMENT_TYPE_UNSPECIFIED"},
		{analysis.CreateDatabase, "CREATE_DATABASE"},
		{analysis.CreateTable, "CREATE_TABLE"},
		{analysis.CreateView, "CREATE_VIEW"},
		{analysis.CreateIndex, "CREATE_INDEX"},
		{analysis.CreateSequence, "CREATE_SEQUENCE"},
		{analysis.CreateSchema, "CREATE_SCHEMA"},
		{analysis.CreateFunction, "CREATE_FUNCTION"},
		{analysis.CreateTrigger, "CREATE_TRIGGER"},
		{analysis.CreateProcedure, "CREATE_PROCEDURE"},
		{analysis.CreateEvent, "CREATE_EVENT"},
		{analysis.CreateExtension, "CREATE_EXTENSION"},
		{analysis.CreateType, "CREATE_TYPE"},
		{analysis.DropDatabase, "DROP_DATABASE"},
		{analysis.DropTable, "DROP_TABLE"},
		{analysis.DropView, "DROP_VIEW"},
		{analysis.DropIndex, "DROP_INDEX"},
		{analysis.DropSequence, "DROP_SEQUENCE"},
		{analysis.DropSchema, "DROP_SCHEMA"},
		{analysis.DropFunction, "DROP_FUNCTION"},
		{analysis.DropTrigger, "DROP_TRIGGER"},
		{analysis.DropProcedure, "DROP_PROCEDURE"},
		{analysis.DropEvent, "DROP_EVENT"},
		{analysis.DropExtension, "DROP_EXTENSION"},
		{analysis.DropType, "DROP_TYPE"},
		{analysis.AlterDatabase, "ALTER_DATABASE"},
		{analysis.AlterTable, "ALTER_TABLE"},
		{analysis.AlterView, "ALTER_VIEW"},
		{analysis.AlterSequence, "ALTER_SEQUENCE"},
		{analysis.AlterEvent, "ALTER_EVENT"},
		{analysis.AlterType, "ALTER_TYPE"},
		{analysis.AlterIndex, "ALTER_INDEX"},
		{analysis.Truncate, "TRUNCATE"},
		{analysis.Rename, "RENAME"},
		{analysis.RenameIndex, "RENAME_INDEX"},
		{analysis.RenameSchema, "RENAME_SCHEMA"},
		{analysis.RenameSequence, "RENAME_SEQUENCE"},
		{analysis.Comment, "COMMENT"},
		{analysis.Insert, "INSERT"},
		{analysis.Update, "UPDATE"},
		{analysis.Delete, "DELETE"},
		{analysis.Merge, "MERGE"},
		{analysis.Execute, "EXECUTE"},
	}

	if analysis.StatementTypeUnspecified != 0 {
		t.Fatalf("StatementTypeUnspecified = %d, want 0", int(analysis.StatementTypeUnspecified))
	}

	seenTypes := make(map[analysis.StatementType]bool, len(cases))
	seenNames := make(map[string]bool, len(cases))
	for _, tc := range cases {
		if seenTypes[tc.typ] {
			t.Errorf("duplicate StatementType numeric value %d (%s)", int(tc.typ), tc.want)
		}
		seenTypes[tc.typ] = true

		got := tc.typ.String()
		if got != tc.want {
			t.Errorf("StatementType(%d).String() = %q, want %q", int(tc.typ), got, tc.want)
		}
		if seenNames[got] {
			t.Errorf("duplicate StatementType name %q", got)
		}
		seenNames[got] = true
	}

	if got := analysis.StatementType(999).String(); got != "STATEMENT_TYPE_UNSPECIFIED" {
		t.Errorf("StatementType(999).String() = %q, want STATEMENT_TYPE_UNSPECIFIED", got)
	}

	type aliased = analysis.StatementType
	var st analysis.StatementType = aliased(analysis.CreateTable)
	if st != analysis.CreateTable {
		t.Fatalf("aliased StatementType = %v, want %v", st, analysis.CreateTable)
	}
}
