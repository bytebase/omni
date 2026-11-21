package completion

import (
	"context"
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

func trinoTestCatalogs() map[string]*metadatapb.DatabaseSchemaMetadata {
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
					Name: "sch",
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
					Name: "wh",
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

func TestCatalogNeeded(t *testing.T) {
	require.False(t, CatalogNeeded("", "company", "select * from school.sch.student"))
	require.True(t, CatalogNeeded("company", "company", "select * from "))
	require.True(t, CatalogNeeded("school", "company", "select * from school."))
	require.False(t, CatalogNeeded("warehouse", "company", "select * from school."))
	// Catalog names with embedded double quotes are always loaded.
	require.True(t, CatalogNeeded(`weird"cat`, "company", "select * from "))
}

func TestCompleteSession_ThreeLevelResolutionAndLazyCatalogLoad(t *testing.T) {
	catalogs := trinoTestCatalogs()
	loads := map[string]int{}
	cfg := SessionConfig{
		DefaultCatalog: "Company",
		DefaultSchema:  "dbo",
		ListCatalogNames: func(context.Context) ([]string, error) {
			return []string{"Company", "School", "Warehouse"}, nil
		},
		LoadCatalog: func(_ context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error) {
			loads[name]++
			return catalogs[name], nil
		},
	}

	// FROM lists all catalogs, current-catalog schemas, and current-schema tables,
	// while loading full metadata only for DefaultCatalog.
	stmt, line, col := catchCaretUTF16(t, "SELECT * FROM |")
	got := CompleteSession(context.Background(), cfg, stmt, line, col)
	require.True(t, hasSessionCandidate(got, "company", CandidateCatalog), "got %v", got)
	require.True(t, hasSessionCandidate(got, "school", CandidateCatalog), "got %v", got)
	require.True(t, hasSessionCandidate(got, "warehouse", CandidateCatalog), "got %v", got)
	require.True(t, hasSessionCandidate(got, "dbo", CandidateSchema), "got %v", got)
	require.True(t, hasSessionCandidate(got, "myschema", CandidateSchema), "got %v", got)
	require.True(t, hasSessionCandidate(got, "employees", CandidateTable), "got %v", got)
	require.Equal(t, 1, loads["Company"])
	require.Zero(t, loads["School"])
	require.Zero(t, loads["Warehouse"])

	// Naming School in the statement loads School and offers its schemas, still not Warehouse.
	clear(loads)
	stmtSchool, lineS, colS := catchCaretUTF16(t, "SELECT * FROM School.|")
	gotSchool := CompleteSession(context.Background(), cfg, stmtSchool, lineS, colS)
	require.True(t, hasSessionCandidate(gotSchool, "sch", CandidateSchema), "got %v", gotSchool)
	require.Equal(t, 1, loads["School"])
	require.Zero(t, loads["Warehouse"])

	// Three-level catalog.schema.table completion.
	clear(loads)
	stmtThree, line3, col3 := catchCaretUTF16(t, "SELECT * FROM School.sch.|")
	gotThree := CompleteSession(context.Background(), cfg, stmtThree, line3, col3)
	require.True(t, hasSessionCandidate(gotThree, "student", CandidateTable), "got %v", gotThree)
	require.Equal(t, 1, loads["School"])
	require.Zero(t, loads["Warehouse"])

	// Catalog referenced only in an earlier statement is not loaded.
	clear(loads)
	stmtMulti, lineM, colM := catchCaretUTF16(t, "SELECT * FROM School.sch.Student; SELECT | FROM Employees")
	gotMulti := CompleteSession(context.Background(), cfg, stmtMulti, lineM, colM)
	require.True(t, hasSessionCandidate(gotMulti, "id", CandidateColumn), "got %v", gotMulti)
	require.True(t, hasSessionCandidate(gotMulti, "name", CandidateColumn), "got %v", gotMulti)
	require.Zero(t, loads["School"])
	require.Zero(t, loads["Warehouse"])

	// Trailing semicolon after a cross-catalog query does not load that catalog.
	clear(loads)
	stmtAfterSemi, lineAS, colAS := catchCaretUTF16(t, "SELECT * FROM School.sch.Student; |")
	gotAS := CompleteSession(context.Background(), cfg, stmtAfterSemi, lineAS, colAS)
	require.True(t, hasSessionCandidate(gotAS, "SELECT", CandidateKeyword), "got %v", gotAS)
	require.Zero(t, loads["School"])
	require.Zero(t, loads["Warehouse"])

	// UTF-16 caret after emoji resolves in-scope columns.
	stmtEmoji, lineE, colE := catchCaretUTF16(t, "SELECT '😀', | FROM Employees")
	gotEmoji := CompleteSession(context.Background(), cfg, stmtEmoji, lineE, colE)
	require.True(t, hasSessionCandidate(gotEmoji, "id", CandidateColumn), "got %v", gotEmoji)
	require.True(t, hasSessionCandidate(gotEmoji, "name", CandidateColumn), "got %v", gotEmoji)
}

func TestCompleteSession_NilAndDefaultOnlyConfig(t *testing.T) {
	// Nil callbacks return no object candidates and never panic.
	stmt, line, col := catchCaretUTF16(t, "SELECT * FROM |")
	for _, c := range CompleteSession(context.Background(), SessionConfig{}, stmt, line, col) {
		require.NotContains(t, []CandidateType{
			CandidateCatalog, CandidateSchema, CandidateTable, CandidateView, CandidateColumn,
		}, c.Type)
	}

	// DefaultCatalog + LoadCatalog works even when ListCatalogNames is nil.
	catalogs := trinoTestCatalogs()
	cfg := SessionConfig{
		DefaultCatalog: "Company",
		DefaultSchema:  "dbo",
		LoadCatalog: func(_ context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error) {
			return catalogs[name], nil
		},
	}
	got := CompleteSession(context.Background(), cfg, stmt, line, col)
	require.True(t, hasSessionCandidate(got, "employees", CandidateTable), "got %v", got)
}
