package completion

import (
	"context"
	"slices"
	"strings"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/stretchr/testify/require"
)

func catchCaretUTF16(t *testing.T, input string) (string, int, int) {
	t.Helper()
	line := 1
	col := 0
	for i, r := range input {
		if r == '|' {
			return input[:i] + input[i+1:], line, col
		}
		if r == '\n' {
			line++
			col = 0
			continue
		}
		if r > 0xFFFF {
			col += 2
		} else {
			col++
		}
	}
	t.Fatalf("missing caret marker in %q", input)
	return "", 0, 0
}

func hasSessionCandidate(cands []Candidate, text string, typ CandidateType) bool {
	for _, c := range cands {
		if c.Text == text && c.Type == typ {
			return true
		}
	}
	return false
}

func mssqlTestDatabases() map[string]*metadatapb.DatabaseSchemaMetadata {
	return map[string]*metadatapb.DatabaseSchemaMetadata{
		"Company": {
			Name: "Company",
			Schemas: []*metadatapb.SchemaMetadata{
				{
					Name: "dbo",
					Tables: []*metadatapb.TableMetadata{
						{
							Name: "Employees",
							Columns: []*metadatapb.ColumnMetadata{
								{Name: "Id", Type: "int"},
								{Name: "Name", Type: "varchar"},
							},
						},
						{
							Name: "Address",
							Columns: []*metadatapb.ColumnMetadata{
								{Name: "EmployeeId", Type: "int"},
								{Name: "Street", Type: "varchar"},
							},
						},
						{
							Name: "Order Items",
							Columns: []*metadatapb.ColumnMetadata{
								{Name: "Item ID", Type: "int"},
								{Name: "Select", Type: "varchar"},
							},
						},
						{
							Name: "a]b",
							Columns: []*metadatapb.ColumnMetadata{
								{Name: "Id", Type: "int"},
							},
						},
						{
							Name: `a"b`,
							Columns: []*metadatapb.ColumnMetadata{
								{Name: "Id", Type: "int"},
							},
						},
					},
					Functions: []*metadatapb.FunctionMetadata{
						{Name: "Calc Total"},
					},
					Procedures: []*metadatapb.ProcedureMetadata{
						{Name: "SyncEmployees"},
					},
				},
				{
					Name: "MySchema",
					Tables: []*metadatapb.TableMetadata{
						{
							Name: "SalaryLevel",
							Columns: []*metadatapb.ColumnMetadata{
								{Name: "SalaryUpBound", Type: "int"},
							},
						},
					},
				},
			},
		},
		"School": {
			Name: "School",
			Schemas: []*metadatapb.SchemaMetadata{
				{
					Name: "dbo",
					Tables: []*metadatapb.TableMetadata{
						{
							Name: "Student",
							Columns: []*metadatapb.ColumnMetadata{
								{Name: "Id", Type: "int"},
								{Name: "ParentName", Type: "varchar"},
							},
						},
					},
				},
			},
		},
		"Warehouse": {
			Name: "Warehouse",
			Schemas: []*metadatapb.SchemaMetadata{
				{
					Name: "dbo",
					Tables: []*metadatapb.TableMetadata{
						{
							Name: "Item",
							Columns: []*metadatapb.ColumnMetadata{
								{Name: "Sku", Type: "varchar"},
							},
						},
					},
				},
			},
		},
	}
}

func TestLineOffsetToBytePos(t *testing.T) {
	input := "SELECT '😀', | FROM Employees"
	stmt, line, col := catchCaretUTF16(t, input)
	require.Equal(t, strings.Index(input, "|"), LineOffsetToBytePos(stmt, line, col))
	require.Equal(t, 3, LineOffsetToBytePos("abc\ndef", 1, 99))
	require.Equal(t, 7, LineOffsetToBytePos("abc\ndef", 9, 0))
}

func TestCaretTokenIsQuoted(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "unquoted identifier", input: "SELECT * FROM Emp|", want: false},
		{name: "open bracket identifier", input: "SELECT * FROM [Emp|", want: true},
		{name: "caret inside closed brackets", input: "SELECT * FROM [Emp|]", want: true},
		{name: "caret inside empty brackets", input: "SELECT * FROM [|]", want: true},
		{name: "caret after closed brackets", input: "SELECT * FROM [Emp]|", want: false},
		{name: "escaped closing bracket inside open bracket", input: "SELECT * FROM [a]]b|", want: true},
		{name: "open double quote identifier", input: `SELECT * FROM "Emp|`, want: true},
		{name: "caret inside closed double quotes", input: `SELECT * FROM "Emp|"`, want: true},
		{name: "caret after closed double quotes", input: `SELECT * FROM "Emp"|`, want: false},
		{name: "space before bracketed identifier is not quoted", input: "SELECT * FROM | [Order Items]", want: false},
		{name: "bracket inside string literal", input: "SELECT '[abc', | FROM Employees", want: false},
		{name: "bracket inside block comment", input: "SELECT /* [ */ | FROM Employees", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stmt, line, col := catchCaretUTF16(t, tc.input)
			pos := LineOffsetToBytePos(stmt, line, col)
			require.Equal(t, tc.want, CaretTokenIsQuoted(stmt, pos))
		})
	}
}

func TestQuoteIdentifierIfNeeded(t *testing.T) {
	require.Equal(t, "Employees", QuoteIdentifierIfNeeded("Employees", false))
	require.Equal(t, "_temp", QuoteIdentifierIfNeeded("_temp", false))
	require.Equal(t, "#temp", QuoteIdentifierIfNeeded("#temp", false))
	require.Equal(t, "@var", QuoteIdentifierIfNeeded("@var", false))
	require.Equal(t, "[Order Items]", QuoteIdentifierIfNeeded("Order Items", false))
	require.Equal(t, "[Select]", QuoteIdentifierIfNeeded("Select", false))
	require.Equal(t, "[123col]", QuoteIdentifierIfNeeded("123col", false))
	require.Equal(t, "[a]]b]", QuoteIdentifierIfNeeded("a]b", false))
	require.Equal(t, "[]", QuoteIdentifierIfNeeded("", false))

	require.Equal(t, "Order Items", QuoteIdentifierIfNeeded("Order Items", true))
	require.Equal(t, "a]]b", QuoteIdentifierIfNeeded("a]b", true))
}

func TestCompleteSession_QuotingFunctionsAndPseudoColumn(t *testing.T) {
	dbs := mssqlTestDatabases()
	cfg := SessionConfig{
		DefaultDatabase: "Company",
		ListDatabaseNames: func(context.Context) ([]string, error) {
			return []string{"Company", "School", "Warehouse"}, nil
		},
		LoadDatabase: func(_ context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error) {
			return dbs[name], nil
		},
	}

	// Regular identifiers stay unquoted; identifiers with spaces, brackets, or keywords get bracket-quoted.
	stmt, line, col := catchCaretUTF16(t, "SELECT * FROM |")
	got := CompleteSession(context.Background(), cfg, stmt, line, col)
	require.True(t, hasSessionCandidate(got, "Employees", CandidateTable), "got %v", got)
	require.True(t, hasSessionCandidate(got, "[Order Items]", CandidateTable), "got %v", got)
	require.True(t, hasSessionCandidate(got, "[a]]b]", CandidateTable), "got %v", got)
	require.True(t, slices.IsSortedFunc(got, compareSessionCandidates))

	// Inside bracketed or double-quoted identifiers, omit outer delimiters while escaping the active delimiter.
	stmtBrk, lineB, colB := catchCaretUTF16(t, "SELECT * FROM [a|]")
	gotBrk := CompleteSession(context.Background(), cfg, stmtBrk, lineB, colB)
	require.True(t, hasSessionCandidate(gotBrk, "a]]b", CandidateTable), "got %v", gotBrk)
	require.False(t, hasSessionCandidate(gotBrk, "a]b", CandidateTable), "got %v", gotBrk)

	stmtDQ, lineDQ, colDQ := catchCaretUTF16(t, `SELECT * FROM "a|"`)
	gotDQ := CompleteSession(context.Background(), cfg, stmtDQ, lineDQ, colDQ)
	require.True(t, hasSessionCandidate(gotDQ, `a""b`, CandidateTable), "got %v", gotDQ)
	require.False(t, hasSessionCandidate(gotDQ, `a"b`, CandidateTable), "got %v", gotDQ)

	// Columns with spaces / reserved keywords are bracket-quoted.
	stmtC, lineC, colC := catchCaretUTF16(t, "SELECT | FROM [Order Items]")
	gotC := CompleteSession(context.Background(), cfg, stmtC, lineC, colC)
	require.True(t, hasSessionCandidate(gotC, "[Item ID]", CandidateColumn), "got %v", gotC)
	require.True(t, hasSessionCandidate(gotC, "[Select]", CandidateColumn), "got %v", gotC)

	// Built-in functions stay unquoted with "()"; non-regular user functions are bracket-quoted before "()".
	stmtF, lineF, colF := catchCaretUTF16(t, "SELECT COU| FROM Employees")
	gotF := CompleteSession(context.Background(), cfg, stmtF, lineF, colF)
	require.True(t, hasSessionCandidate(gotF, "COUNT()", CandidateFunction), "got %v", gotF)

	stmtCoalesce, lineCo, colCo := catchCaretUTF16(t, "SELECT COAL| FROM Employees")
	gotCo := CompleteSession(context.Background(), cfg, stmtCoalesce, lineCo, colCo)
	require.True(t, hasSessionCandidate(gotCo, "COALESCE()", CandidateFunction), "got %v", gotCo)

	require.Equal(t, "[Calc Total]()", formatSessionCandidateText(Candidate{Text: "Calc Total", Type: CandidateFunction}, 0))
	require.Equal(t, "Calc Total", formatSessionCandidateText(Candidate{Text: "Calc Total", Type: CandidateFunction}, '['))

	// $action pseudo-column is NOT bracket-quoted.
	stmtA, lineA, colA := catchCaretUTF16(t, "MERGE INTO Employees AS d USING Address AS s ON d.Id = s.EmployeeId WHEN MATCHED THEN UPDATE SET d.Name = s.Street OUTPUT |")
	gotA := CompleteSession(context.Background(), cfg, stmtA, lineA, colA)
	require.True(t, hasSessionCandidate(gotA, "$action", CandidateColumn), "got %v", gotA)
	require.False(t, hasSessionCandidate(gotA, "[$action]", CandidateColumn), "got %v", gotA)
}

func TestCompleteSession_StatementScopedDatabaseLoading(t *testing.T) {
	dbs := mssqlTestDatabases()
	loads := map[string]int{}
	cfg := SessionConfig{
		DefaultDatabase: "company", // case-folded; canonicalizes to "Company"
		ListDatabaseNames: func(context.Context) ([]string, error) {
			return []string{"Company", "School", "Warehouse"}, nil
		},
		LoadDatabase: func(_ context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error) {
			loads[name]++
			return dbs[name], nil
		},
	}

	// Default-DB query loads only Company, not School or Warehouse.
	stmt, line, col := catchCaretUTF16(t, "SELECT | FROM Employees")
	got := CompleteSession(context.Background(), cfg, stmt, line, col)
	require.True(t, hasSessionCandidate(got, "Id", CandidateColumn), "got %v", got)
	require.Equal(t, 1, loads["Company"])
	require.Zero(t, loads["School"])
	require.Zero(t, loads["Warehouse"])

	// USE lists all databases without loading School or Warehouse metadata.
	clear(loads)
	stmtUse, lineUse, colUse := catchCaretUTF16(t, "USE |")
	gotUse := CompleteSession(context.Background(), cfg, stmtUse, lineUse, colUse)
	require.True(t, hasSessionCandidate(gotUse, "School", CandidateDatabase), "got %v", gotUse)
	require.True(t, hasSessionCandidate(gotUse, "Warehouse", CandidateDatabase), "got %v", gotUse)
	require.Zero(t, loads["School"])
	require.Zero(t, loads["Warehouse"])

	// Cross-database reference loads School, still skips Warehouse.
	clear(loads)
	stmtCross, lineCross, colCross := catchCaretUTF16(t, "SELECT School.dbo.Student.|")
	gotCross := CompleteSession(context.Background(), cfg, stmtCross, lineCross, colCross)
	require.True(t, hasSessionCandidate(gotCross, "ParentName", CandidateColumn), "got %v", gotCross)
	require.Equal(t, 1, loads["School"])
	require.Zero(t, loads["Warehouse"])

	// Database referenced only in a preceding statement is isolated away and not loaded.
	clear(loads)
	stmtMulti, lineMulti, colMulti := catchCaretUTF16(t, "SELECT * FROM School.dbo.Student; SELECT | FROM Employees")
	gotMulti := CompleteSession(context.Background(), cfg, stmtMulti, lineMulti, colMulti)
	require.True(t, hasSessionCandidate(gotMulti, "Name", CandidateColumn), "got %v", gotMulti)
	require.Zero(t, loads["School"])
}

func TestCompleteSession_StatementIsolationAndUnsyncedDB(t *testing.T) {
	dbs := mssqlTestDatabases()
	cfg := SessionConfig{
		DefaultDatabase: "Company",
		DefaultSchema:   "MySchema",
		LoadDatabase: func(_ context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error) {
			return dbs[name], nil
		},
	}

	// DefaultSchema resolves SalaryLevel unqualified.
	stmtS, lineS, colS := catchCaretUTF16(t, "SELECT | FROM SalaryLevel")
	gotS := CompleteSession(context.Background(), cfg, stmtS, lineS, colS)
	require.True(t, hasSessionCandidate(gotS, "SalaryUpBound", CandidateColumn), "got %v", gotS)

	// Following statement after semicolon does not leak alias scope.
	stmtFollow, lineF, colF := catchCaretUTF16(t, "SELECT a.|; SELECT * FROM dbo.Address AS a")
	gotF := CompleteSession(context.Background(), cfg, stmtFollow, lineF, colF)
	require.False(t, hasSessionCandidate(gotF, "Street", CandidateColumn), "got %v", gotF)

	// GO batch separator isolates statements.
	stmtGo, lineGo, colGo := catchCaretUTF16(t, "SELECT 1\nGO\nSELECT * FROM |")
	gotGo := CompleteSession(context.Background(), cfg, stmtGo, lineGo, colGo)
	require.True(t, hasSessionCandidate(gotGo, "SalaryLevel", CandidateTable), "got %v", gotGo)

	// Unsynced database returns no schema/table/column candidates.
	unsyncedCfg := SessionConfig{
		DefaultDatabase: "Company",
		ListDatabaseNames: func(context.Context) ([]string, error) {
			return []string{"Company"}, nil
		},
		LoadDatabase: func(context.Context, string) (*metadatapb.DatabaseSchemaMetadata, error) {
			return nil, nil
		},
	}
	stmtU, lineU, colU := catchCaretUTF16(t, "SELECT * FROM |")
	for _, c := range CompleteSession(context.Background(), unsyncedCfg, stmtU, lineU, colU) {
		require.NotContains(t, []CandidateType{CandidateSchema, CandidateTable, CandidateColumn}, c.Type)
	}
}
