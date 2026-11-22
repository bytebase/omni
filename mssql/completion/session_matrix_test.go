package completion

import (
	"context"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/stretchr/testify/require"
)

type matrixCandidateSpec struct {
	text string
	typ  CandidateType
}

func mssqlMatrixSessionConfig() SessionConfig {
	dbs := map[string]*metadatapb.DatabaseSchemaMetadata{
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
					Sequences: []*metadatapb.SequenceMetadata{
						{Name: "EmployeeIdSeq"},
						{Name: "OrderSeq"},
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
								{Name: "Id", Type: "int"},
								{Name: "SalaryUpBound", Type: "int"},
							},
						},
					},
					Views: []*metadatapb.ViewMetadata{
						{Name: "SalaryView"},
					},
					Sequences: []*metadatapb.SequenceMetadata{
						{Name: "SalarySeq"},
					},
					Procedures: []*metadatapb.ProcedureMetadata{
						{Name: "SyncSalary"},
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
					Sequences: []*metadatapb.SequenceMetadata{
						{Name: "StudentSeq"},
					},
					Procedures: []*metadatapb.ProcedureMetadata{
						{Name: "SyncStudents"},
					},
				},
			},
		},
	}

	return SessionConfig{
		DefaultDatabase: "Company",
		ListDatabaseNames: func(context.Context) ([]string, error) {
			return []string{"Company", "School"}, nil
		},
		LoadDatabase: func(_ context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error) {
			return dbs[name], nil
		},
	}
}

func TestCompleteSession_CoverageMatrix(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		want    []matrixCandidateSpec
		notWant []matrixCandidateSpec
	}{
		{
			name: "select list uses table columns",
			sql:  "SELECT | FROM Employees",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "select list uses schema qualified table columns",
			sql:  "SELECT | FROM dbo.Employees",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "select list uses implicit table alias columns",
			sql:  "SELECT e.| FROM Employees e",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "select list uses explicit table alias columns",
			sql:  "SELECT e.| FROM Employees AS e",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "select list prefix keeps matching column available",
			sql:  "SELECT Na| FROM Employees",
			want: matrixColumns("Name"),
		},
		{
			name: "qualified column prefix keeps matching column available",
			sql:  "SELECT e.Na| FROM Employees AS e",
			want: matrixColumns("Name"),
		},
		{
			name:    "cross database qualified column reference",
			sql:     "SELECT School.dbo.Student.|",
			want:    matrixColumns("Id", "ParentName"),
			notWant: matrixColumns("Name"),
		},
		{
			name: "where clause uses current table columns",
			sql:  "SELECT * FROM Employees WHERE |",
			want: matrixColumns("Id", "Name"),
		},
		{
			name:    "where clause excludes select alias",
			sql:     "SELECT Id AS IdAlias FROM Employees WHERE |",
			want:    matrixColumns("Id", "Name"),
			notWant: matrixColumns("IdAlias"),
		},
		{
			name:    "where clause uses qualified alias columns",
			sql:     "SELECT * FROM Employees e WHERE e.|",
			want:    matrixColumns("Id", "Name"),
			notWant: matrixColumns("EmployeeId", "Street"),
		},
		{
			name: "group by uses current table columns",
			sql:  "SELECT * FROM Employees GROUP BY |",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "having clause uses current table columns",
			sql:  "SELECT Name, COUNT(*) FROM Employees GROUP BY Name HAVING |",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "between expression uses current table columns",
			sql:  "SELECT * FROM Employees WHERE Id BETWEEN | AND 10",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "like expression uses current table columns",
			sql:  "SELECT * FROM Employees WHERE Name LIKE |",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "in expression uses current table columns",
			sql:  "SELECT * FROM Employees WHERE Id IN (|)",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "in expression after comma uses current table columns",
			sql:  "SELECT * FROM Employees WHERE Id IN (1, |)",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "or predicate uses current table columns",
			sql:  "SELECT * FROM Employees WHERE Id = 1 OR |",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "and predicate after is null uses current table columns",
			sql:  "SELECT * FROM Employees WHERE Id IS NULL AND |",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "not predicate uses current table columns",
			sql:  "SELECT * FROM Employees WHERE NOT |",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "comparison expression uses current table columns",
			sql:  "SELECT * FROM Employees WHERE Id > |",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "qualified comparison expression uses alias columns",
			sql:  "SELECT * FROM Employees e WHERE e.Id = |",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "arithmetic expression uses current table columns",
			sql:  "SELECT * FROM Employees WHERE Id + | > 0",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "between upper bound uses current table columns",
			sql:  "SELECT * FROM Employees WHERE Id BETWEEN 1 AND |",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "parenthesized predicate uses current table columns",
			sql:  "SELECT * FROM Employees WHERE (|)",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "case predicate uses current table columns",
			sql:  "SELECT CASE WHEN | THEN 1 ELSE 0 END FROM Employees",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "case then expression uses current table columns",
			sql:  "SELECT CASE WHEN Id > 0 THEN | ELSE 0 END FROM Employees",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "case else expression uses current table columns",
			sql:  "SELECT CASE WHEN Id > 0 THEN 1 ELSE | END FROM Employees",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "function argument uses current table columns",
			sql:  "SELECT COUNT(|) FROM Employees",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "scalar function argument uses current table columns",
			sql:  "SELECT ABS(|) FROM Employees",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "nested function argument uses current table columns",
			sql:  "SELECT ABS(ROUND(|, 0)) FROM Employees",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "isnull argument uses current table columns",
			sql:  "SELECT ISNULL(|, 0) FROM Employees",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "coalesce argument uses current table columns",
			sql:  "SELECT COALESCE(NULL, |) FROM Employees",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "nullif argument uses current table columns",
			sql:  "SELECT NULLIF(|, 0) FROM Employees",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "cast argument uses current table columns",
			sql:  "SELECT CAST(| AS INT) FROM Employees",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "convert argument uses current table columns",
			sql:  "SELECT CONVERT(INT, |) FROM Employees",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "iif predicate uses current table columns",
			sql:  "SELECT IIF(|, 1, 0) FROM Employees",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "order by includes select alias",
			sql:  "SELECT Id AS IdAlias, Name FROM Employees ORDER BY |",
			want: matrixColumns("Id", "Name", "IdAlias"),
		},
		{
			name: "order by includes tsql equals alias",
			sql:  "SELECT IdAlias = Id, Name FROM Employees ORDER BY |",
			want: matrixColumns("Id", "Name", "IdAlias"),
		},
		{
			name: "order by alias prefix keeps select alias available",
			sql:  "SELECT Id AS IdAlias, Name FROM Employees ORDER BY IdA|",
			want: matrixColumns("IdAlias"),
		},
		{
			name: "order by column prefix keeps matching column available",
			sql:  "SELECT * FROM Employees ORDER BY Na|",
			want: matrixColumns("Name"),
		},
		{
			name: "top select list uses columns",
			sql:  "SELECT TOP 10 | FROM Employees",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "distinct select list uses columns",
			sql:  "SELECT DISTINCT | FROM Employees",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "window partition by uses columns",
			sql:  "SELECT ROW_NUMBER() OVER (PARTITION BY |) FROM Employees",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "window partition by uses qualified alias columns",
			sql:  "SELECT ROW_NUMBER() OVER (PARTITION BY e.| ORDER BY e.Id) FROM Employees e",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "window order by uses columns",
			sql:  "SELECT ROW_NUMBER() OVER (ORDER BY |) FROM Employees",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "window order by uses qualified alias columns",
			sql:  "SELECT ROW_NUMBER() OVER (PARTITION BY e.Id ORDER BY e.|) FROM Employees e",
			want: matrixColumns("Id", "Name"),
		},
		{
			name:    "window order by excludes select alias",
			sql:     "SELECT Id AS IdAlias, ROW_NUMBER() OVER (ORDER BY |) FROM Employees",
			want:    matrixColumns("Id", "Name"),
			notWant: matrixColumns("IdAlias"),
		},
		{
			name:    "bracket quoted alias columns",
			sql:     "SELECT [e].| FROM Employees AS [e]",
			want:    matrixColumns("Id", "Name"),
			notWant: matrixColumns("EmployeeId", "Street"),
		},
		{
			name: "double quoted alias columns",
			sql:  `SELECT "e".| FROM Employees AS "e"`,
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "order by alias columns before offset",
			sql:  "SELECT * FROM Employees e ORDER BY e.| OFFSET 10 ROWS",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "offset fetch order by keeps alias columns",
			sql:  "SELECT * FROM Employees e ORDER BY e.| OFFSET 10 ROWS FETCH NEXT 5 ROWS ONLY",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "join table reference",
			sql:  "SELECT * FROM Employees JOIN |",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "left join table reference",
			sql:  "SELECT * FROM Employees LEFT JOIN |",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "right join table reference",
			sql:  "SELECT * FROM Employees RIGHT JOIN |",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "full outer join table reference",
			sql:  "SELECT * FROM Employees FULL OUTER JOIN |",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "cross join table reference",
			sql:  "SELECT * FROM Employees CROSS JOIN |",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "comma join table reference",
			sql:  "SELECT * FROM Employees, |",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "cross apply table reference",
			sql:  "SELECT * FROM Employees CROSS APPLY |",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "outer apply table reference",
			sql:  "SELECT * FROM Employees OUTER APPLY |",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "join on uses left alias columns",
			sql:  "SELECT * FROM Employees e JOIN Address a ON e.|",
			want: matrixColumns("Id", "Name"),
		},
		{
			name:    "join on excludes select alias",
			sql:     "SELECT e.Id AS IdAlias FROM Employees e JOIN Address a ON |",
			want:    matrixColumns("Id", "Name", "EmployeeId", "Street"),
			notWant: matrixColumns("IdAlias"),
		},
		{
			name:    "join on uses right alias columns",
			sql:     "SELECT * FROM Employees e JOIN Address a ON e.Id = a.|",
			want:    matrixColumns("EmployeeId", "Street"),
			notWant: matrixColumns("Name"),
		},
		{
			name: "multi join where sees all joined table columns",
			sql:  "SELECT * FROM Employees e JOIN Address a ON e.Id = a.EmployeeId JOIN MySchema.SalaryLevel s ON e.Id = s.Id WHERE |",
			want: matrixColumns("Id", "Name", "EmployeeId", "Street", "SalaryUpBound"),
		},
		{
			name: "cross join where uses right alias columns",
			sql:  "SELECT * FROM Employees e CROSS JOIN Address a WHERE a.|",
			want: matrixColumns("EmployeeId", "Street"),
		},
		{
			name: "schema qualified join on uses right alias columns",
			sql:  "SELECT * FROM Employees e LEFT JOIN MySchema.SalaryLevel s ON s.|",
			want: matrixColumns("Id", "SalaryUpBound"),
		},
		{
			name:    "joined order by uses right alias columns",
			sql:     "SELECT * FROM Employees e JOIN Address a ON e.Id = a.EmployeeId ORDER BY a.|",
			want:    matrixColumns("EmployeeId", "Street"),
			notWant: matrixColumns("Name"),
		},
		{
			name: "joined group by uses right alias columns",
			sql:  "SELECT a.EmployeeId, COUNT(*) FROM Employees e JOIN Address a ON e.Id = a.EmployeeId GROUP BY a.|",
			want: matrixColumns("EmployeeId", "Street"),
		},
		{
			name: "default schema table reference",
			sql:  "SELECT * FROM |",
			want: append(append(matrixTables("Address", "Employees"), matrixSchemas("dbo", "MySchema")...), matrixSequences("EmployeeIdSeq", "OrderSeq")...),
		},
		{
			name:    "database qualified schema reference",
			sql:     "SELECT * FROM Company.|",
			want:    matrixSchemas("dbo", "MySchema"),
			notWant: append(append(matrixDatabases("Company", "School"), matrixTables("Address", "Employees")...), matrixSequences("EmployeeIdSeq", "OrderSeq")...),
		},
		{
			name:    "database default schema table reference",
			sql:     "SELECT * FROM School..|",
			want:    matrixTables("Student"),
			notWant: matrixTables("Employees"),
		},
		{
			name:    "linked server qualified table reference does not use local database",
			sql:     "SELECT * FROM Linked.Company.dbo.|",
			notWant: matrixTables("Address", "Employees"),
		},
		{
			name: "dbo table reference",
			sql:  "SELECT * FROM dbo.|",
			want: matrixTables("Address", "Employees"),
		},
		{
			name:    "alternate schema table reference",
			sql:     "SELECT * FROM MySchema.|",
			want:    matrixTables("SalaryLevel"),
			notWant: matrixTables("Employees"),
		},
		{
			name: "alternate schema view reference",
			sql:  "SELECT * FROM MySchema.|",
			want: matrixViews("SalaryView"),
		},
		{
			name:    "cross database table reference",
			sql:     "SELECT * FROM School.dbo.|",
			want:    matrixTables("Student"),
			notWant: matrixTables("Employees"),
		},
		{
			name: "database qualified dbo table reference",
			sql:  "SELECT * FROM Company.dbo.|",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "database qualified alternate schema table reference",
			sql:  "SELECT * FROM Company.MySchema.|",
			want: matrixTables("SalaryLevel"),
		},
		{
			name: "bracket quoted schema table reference",
			sql:  "SELECT * FROM [dbo].|",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "bracket quoted table alias columns",
			sql:  "SELECT * FROM [dbo].[Employees] AS [emp] WHERE [emp].|",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "schema qualified table prefix keeps matching table available",
			sql:  "SELECT * FROM dbo.Emp|",
			want: matrixTables("Employees"),
		},
		{
			name: "cross database table prefix keeps matching table available",
			sql:  "SELECT * FROM School.dbo.St|",
			want: matrixTables("Student"),
		},
		{
			name: "table prefix keeps matching table available",
			sql:  "SELECT * FROM Emp|",
			want: matrixTables("Employees"),
		},
		{
			name: "schema prefix keeps matching schema available",
			sql:  "SELECT * FROM My|",
			want: matrixSchemas("MySchema"),
		},
		{
			name: "database qualified schema prefix keeps matching schema available",
			sql:  "SELECT * FROM Company.My|",
			want: matrixSchemas("MySchema"),
		},
		{
			name: "cte explicit columns are available",
			sql:  "WITH cte(EmpId, EmpName) AS (SELECT Id, Name FROM Employees) SELECT cte.| FROM cte",
			want: matrixColumns("EmpId", "EmpName"),
		},
		{
			name: "cte table is available in from",
			sql:  "WITH cte AS (SELECT Id, Name FROM Employees) SELECT * FROM |",
			want: matrixTables("cte"),
		},
		{
			name:    "cte table is excluded from schema qualified from",
			sql:     "WITH cte AS (SELECT Id, Name FROM Employees) SELECT * FROM dbo.|",
			want:    matrixTables("Address", "Employees"),
			notWant: matrixTables("cte"),
		},
		{
			name: "cte projected columns are available",
			sql:  "WITH cte AS (SELECT Id, Name FROM Employees) SELECT | FROM cte",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "cte qualified projected columns are available",
			sql:  "WITH cte AS (SELECT Id, Name FROM Employees) SELECT cte.| FROM cte",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "cte where uses projected columns",
			sql:  "WITH cte AS (SELECT Id, Name FROM Employees) SELECT * FROM cte WHERE |",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "cte insert select uses projected columns",
			sql:  "WITH cte AS (SELECT Id, Name FROM Employees) INSERT INTO Employees SELECT | FROM cte",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "chained cte table is available",
			sql:  "WITH a AS (SELECT Id FROM Employees), b AS (SELECT Id FROM a) SELECT * FROM |",
			want: matrixTables("a", "b"),
		},
		{
			name: "subquery select list uses inner table columns",
			sql:  "SELECT * FROM Employees WHERE EXISTS (SELECT | FROM Address)",
			want: matrixColumns("EmployeeId", "Street"),
		},
		{
			name: "subquery where uses inner table columns",
			sql:  "SELECT * FROM Employees WHERE EXISTS (SELECT 1 FROM Address WHERE |)",
			want: matrixColumns("EmployeeId", "Street"),
		},
		{
			name: "in subquery select list uses inner table columns",
			sql:  "SELECT * FROM Employees WHERE Id IN (SELECT | FROM Address)",
			want: matrixColumns("EmployeeId", "Street"),
		},
		{
			name: "apply subquery select list uses inner table columns",
			sql:  "SELECT * FROM Employees e CROSS APPLY (SELECT | FROM Address) x",
			want: matrixColumns("EmployeeId", "Street"),
		},
		{
			name: "scalar subquery select list uses inner table columns",
			sql:  "SELECT (SELECT | FROM Employees) FROM Address",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "aggregate scalar subquery uses inner table columns",
			sql:  "SELECT * FROM Employees WHERE Id = (SELECT MAX(|) FROM Address)",
			want: matrixColumns("EmployeeId", "Street"),
		},
		{
			name: "derived table alias columns",
			sql:  "SELECT d.| FROM (SELECT Id, Name FROM Employees) d",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "nested derived table alias columns",
			sql:  "SELECT x.| FROM (SELECT d.Id FROM (SELECT Id FROM Employees) d) x",
			want: matrixColumns("Id"),
		},
		{
			name: "values derived table alias columns",
			sql:  "SELECT v.| FROM (VALUES (1, 'a')) AS v(Id, ValueLabel)",
			want: matrixColumns("Id", "ValueLabel"),
		},
		{
			name: "correlated subquery can complete outer alias",
			sql:  "SELECT * FROM Employees e WHERE EXISTS (SELECT 1 FROM Address a WHERE a.EmployeeId = e.|)",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "correlated subquery can complete inner alias",
			sql:  "SELECT * FROM Employees e WHERE EXISTS (SELECT 1 FROM Address a WHERE a.| = e.Id)",
			want: matrixColumns("EmployeeId", "Street"),
		},
		{
			name:    "inner alias shadows outer alias",
			sql:     "SELECT * FROM Employees o WHERE EXISTS (SELECT 1 FROM Address o WHERE o.| = 1)",
			want:    matrixColumns("EmployeeId", "Street"),
			notWant: matrixColumns("Name"),
		},
		{
			name: "any subquery select list uses inner table columns",
			sql:  "SELECT * FROM Employees WHERE Id = ANY (SELECT | FROM Address)",
			want: matrixColumns("EmployeeId", "Street"),
		},
		{
			name: "all subquery select list uses inner table columns",
			sql:  "SELECT * FROM Employees WHERE Id > ALL (SELECT | FROM Address)",
			want: matrixColumns("EmployeeId", "Street"),
		},
		{
			name: "union right arm uses right table columns",
			sql:  "SELECT Id FROM Employees UNION SELECT | FROM Address",
			want: matrixColumns("EmployeeId", "Street"),
		},
		{
			name: "union all right arm uses right table columns",
			sql:  "SELECT Id FROM Employees UNION ALL SELECT | FROM Address",
			want: matrixColumns("EmployeeId", "Street"),
		},
		{
			name:    "chained union final arm uses final table columns",
			sql:     "SELECT Id FROM Employees UNION SELECT EmployeeId FROM Address UNION SELECT | FROM MySchema.SalaryLevel",
			want:    matrixColumns("Id", "SalaryUpBound"),
			notWant: matrixColumns("EmployeeId", "Street"),
		},
		{
			name: "except right arm uses right table columns",
			sql:  "SELECT Id FROM Employees EXCEPT SELECT | FROM Address",
			want: matrixColumns("EmployeeId", "Street"),
		},
		{
			name: "intersect right arm uses right table columns",
			sql:  "SELECT Id FROM Employees INTERSECT SELECT | FROM Address",
			want: matrixColumns("EmployeeId", "Street"),
		},
		{
			name: "insert target table reference",
			sql:  "INSERT INTO |",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "insert target schema table reference",
			sql:  "INSERT INTO dbo.|",
			want: matrixTables("Address", "Employees"),
		},
		{
			name:    "insert target cross database table reference",
			sql:     "INSERT INTO School.dbo.|",
			want:    matrixTables("Student"),
			notWant: matrixTables("Employees"),
		},
		{
			name: "insert target alternate schema table reference",
			sql:  "INSERT INTO MySchema.|",
			want: matrixTables("SalaryLevel"),
		},
		{
			name: "insert column list uses target columns",
			sql:  "INSERT INTO Employees(|) VALUES (1)",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "insert schema qualified column list uses target columns",
			sql:  "INSERT INTO dbo.Employees(|) VALUES (1)",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "insert cross database column list uses target columns",
			sql:  "INSERT INTO School.dbo.Student(|) VALUES (1)",
			want: matrixColumns("Id", "ParentName"),
		},
		{
			name:    "insert column list after comma uses remaining target columns",
			sql:     "INSERT INTO Employees(Id, |) VALUES (1, 'a')",
			want:    matrixColumns("Name"),
			notWant: matrixColumns("EmployeeId", "Id", "SalaryUpBound", "Street"),
		},
		{
			name: "insert select uses source columns",
			sql:  "INSERT INTO Employees SELECT | FROM Address",
			want: matrixColumns("EmployeeId", "Street"),
		},
		{
			name: "insert output inserted columns",
			sql:  "INSERT INTO Employees(Name) OUTPUT INSERTED.| VALUES ('a')",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "update target table reference",
			sql:  "UPDATE | SET Name = 'a'",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "update target schema table reference",
			sql:  "UPDATE dbo.| SET Name = 'a'",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "update target cross database table reference",
			sql:  "UPDATE School.dbo.| SET ParentName = 'a'",
			want: matrixTables("Student"),
		},
		{
			name: "update set column uses target columns",
			sql:  "UPDATE Employees SET |",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "update cross database set column uses target columns",
			sql:  "UPDATE School.dbo.Student SET |",
			want: matrixColumns("Id", "ParentName"),
		},
		{
			name: "update where uses target columns",
			sql:  "UPDATE Employees SET Name = 'a' WHERE |",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "update set value uses target columns",
			sql:  "UPDATE Employees SET Name = |",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "update top target table reference",
			sql:  "UPDATE TOP (10) | SET Name = 'a'",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "update alias set column uses target columns",
			sql:  "UPDATE e SET | FROM Employees e",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "update alias set value uses target columns",
			sql:  "UPDATE e SET Name = | FROM Employees e",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "update joined set value uses joined alias columns",
			sql:  "UPDATE e SET Name = a.| FROM Employees e JOIN Address a ON e.Id = a.EmployeeId",
			want: matrixColumns("EmployeeId", "Street"),
		},
		{
			name: "update from joined alias columns",
			sql:  "UPDATE e SET Name = 'x' FROM Employees e JOIN Address a ON e.Id = a.EmployeeId WHERE a.|",
			want: matrixColumns("EmployeeId", "Street"),
		},
		{
			name: "update output inserted columns",
			sql:  "UPDATE Employees SET Name = 'x' OUTPUT INSERTED.|",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "update output deleted columns",
			sql:  "UPDATE Employees SET Name = 'x' OUTPUT DELETED.|",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "delete target table reference",
			sql:  "DELETE FROM |",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "delete target schema table reference",
			sql:  "DELETE FROM dbo.|",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "delete target cross database table reference",
			sql:  "DELETE FROM School.dbo.|",
			want: matrixTables("Student"),
		},
		{
			name: "delete top target table reference",
			sql:  "DELETE TOP (10) FROM |",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "delete where uses target columns",
			sql:  "DELETE FROM Employees WHERE |",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "delete cross database where uses target columns",
			sql:  "DELETE FROM School.dbo.Student WHERE |",
			want: matrixColumns("Id", "ParentName"),
		},
		{
			name: "delete alias where uses target columns",
			sql:  "DELETE e FROM Employees e WHERE e.|",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "delete joined alias where uses joined columns",
			sql:  "DELETE e FROM Employees e JOIN Address a ON e.Id = a.EmployeeId WHERE a.|",
			want: matrixColumns("EmployeeId", "Street"),
		},
		{
			name: "delete output deleted columns",
			sql:  "DELETE FROM Employees OUTPUT DELETED.|",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "create index target table reference",
			sql:  "CREATE INDEX ix ON |",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "create index target schema table reference",
			sql:  "CREATE INDEX ix ON dbo.|",
			want: matrixTables("Address", "Employees"),
		},
		{
			name:    "create index column list uses target columns",
			sql:     "CREATE INDEX ix ON Employees(|)",
			want:    matrixColumns("Id", "Name"),
			notWant: matrixColumns("EmployeeId", "SalaryUpBound", "Street"),
		},
		{
			name: "create index schema qualified column list uses target columns",
			sql:  "CREATE INDEX ix ON dbo.Employees(|)",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "create unique clustered index column list uses target columns",
			sql:  "CREATE UNIQUE CLUSTERED INDEX ix ON Employees(|)",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "create index include column list uses target columns",
			sql:  "CREATE INDEX ix ON Employees(Id) INCLUDE (|)",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "create view body select uses source columns",
			sql:  "CREATE VIEW v AS SELECT | FROM Employees",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "create view with schema body select uses source columns",
			sql:  "CREATE VIEW MySchema.v AS SELECT | FROM Employees",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "create procedure body select uses source columns",
			sql:  "CREATE PROCEDURE p AS SELECT | FROM Employees",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "create procedure parameter type name",
			sql:  "CREATE PROCEDURE p @Name | AS SELECT 1",
			want: matrixTypes("INT", "NVARCHAR"),
		},
		{
			name: "create table references table",
			sql:  "CREATE TABLE NewTable (EmployeeId INT REFERENCES |)",
			want: matrixTables("Address", "Employees"),
		},
		{
			name:    "create table references cross database table",
			sql:     "CREATE TABLE NewTable (StudentId INT REFERENCES School.dbo.|)",
			want:    matrixTables("Student"),
			notWant: matrixTables("Employees"),
		},
		{
			name:    "create table foreign key source column list",
			sql:     "CREATE TABLE NewTable (EmployeeId INT, FOREIGN KEY (|) REFERENCES Employees(Id))",
			want:    matrixColumns("EmployeeId"),
			notWant: matrixColumns("Id", "Name", "Street"),
		},
		{
			name: "create table named foreign key source column list",
			sql:  "CREATE TABLE NewTable (EmployeeId INT, CONSTRAINT fk FOREIGN KEY (|) REFERENCES Employees(Id))",
			want: matrixColumns("EmployeeId"),
		},
		{
			name:    "create table references column list",
			sql:     "CREATE TABLE NewTable (EmployeeId INT REFERENCES Employees(|))",
			want:    matrixColumns("Id", "Name"),
			notWant: matrixColumns("EmployeeId", "SalaryUpBound", "Street"),
		},
		{
			name: "create table schema qualified references column list",
			sql:  "CREATE TABLE NewTable (EmployeeId INT REFERENCES dbo.Employees(|))",
			want: matrixColumns("Id", "Name"),
		},
		{
			name:    "create table cross database references column list",
			sql:     "CREATE TABLE NewTable (StudentId INT REFERENCES School.dbo.Student(|))",
			want:    matrixColumns("Id", "ParentName"),
			notWant: matrixColumns("Name"),
		},
		{
			name: "create table type name",
			sql:  "CREATE TABLE NewTable (Name |)",
			want: matrixTypes("INT", "NVARCHAR"),
		},
		{
			name: "create table type name after comma",
			sql:  "CREATE TABLE NewTable (Id INT, Name |)",
			want: matrixTypes("INT", "NVARCHAR"),
		},
		{
			name: "declare variable type name",
			sql:  "DECLARE @Name |",
			want: matrixTypes("INT", "NVARCHAR"),
		},
		{
			name: "alter table target reference",
			sql:  "ALTER TABLE |",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "alter table target schema reference",
			sql:  "ALTER TABLE dbo.|",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "alter table cross database target reference",
			sql:  "ALTER TABLE School.dbo.|",
			want: matrixTables("Student"),
		},
		{
			name: "alter table add column type name",
			sql:  "ALTER TABLE Employees ADD Name |",
			want: matrixTypes("INT", "NVARCHAR"),
		},
		{
			name: "alter table references table",
			sql:  "ALTER TABLE Employees ADD CONSTRAINT fk FOREIGN KEY (Id) REFERENCES |",
			want: matrixTables("Address", "Employees"),
		},
		{
			name:    "alter table references column list",
			sql:     "ALTER TABLE Employees ADD CONSTRAINT fk FOREIGN KEY (Id) REFERENCES Address(|)",
			want:    matrixColumns("EmployeeId", "Street"),
			notWant: matrixColumns("Name"),
		},
		{
			name: "drop table reference",
			sql:  "DROP TABLE |",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "drop table prefix keeps matching table available",
			sql:  "DROP TABLE Emp|",
			want: matrixTables("Employees"),
		},
		{
			name: "drop table if exists reference",
			sql:  "DROP TABLE IF EXISTS |",
			want: matrixTables("Address", "Employees"),
		},
		{
			name:    "drop table if exists database qualified schema reference",
			sql:     "DROP TABLE IF EXISTS School.|",
			want:    matrixSchemas("dbo"),
			notWant: matrixTables("Address", "Employees"),
		},
		{
			name: "drop schema qualified table reference",
			sql:  "DROP TABLE dbo.|",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "drop cross database table reference",
			sql:  "DROP TABLE School.dbo.|",
			want: matrixTables("Student"),
		},
		{
			name: "drop view schema reference",
			sql:  "DROP VIEW MySchema.|",
			want: matrixViews("SalaryView"),
		},
		{
			name: "drop procedure schema reference",
			sql:  "DROP PROCEDURE dbo.|",
			want: matrixProcedures("SyncEmployees"),
		},
		{
			name:    "drop procedure database qualified schema reference",
			sql:     "DROP PROCEDURE Company.|",
			want:    matrixSchemas("dbo", "MySchema"),
			notWant: matrixProcedures("SyncEmployees", "SyncSalary"),
		},
		{
			name:    "drop procedure if exists database qualified schema reference",
			sql:     "DROP PROCEDURE IF EXISTS Company.|",
			want:    matrixSchemas("dbo", "MySchema"),
			notWant: matrixProcedures("SyncEmployees", "SyncSalary"),
		},
		{
			name: "drop procedure prefix keeps matching routine available",
			sql:  "DROP PROCEDURE Sync|",
			want: matrixProcedures("SyncEmployees"),
		},
		{
			name: "drop sequence schema reference",
			sql:  "DROP SEQUENCE dbo.|",
			want: matrixSequences("EmployeeIdSeq", "OrderSeq"),
		},
		{
			name: "drop sequence prefix keeps matching sequence available",
			sql:  "DROP SEQUENCE Employee|",
			want: matrixSequences("EmployeeIdSeq"),
		},
		{
			name: "drop database reference",
			sql:  "DROP DATABASE |",
			want: matrixDatabases("Company", "School"),
		},
		{
			name: "drop database prefix keeps matching database available",
			sql:  "DROP DATABASE Comp|",
			want: matrixDatabases("Company"),
		},
		{
			name: "use database reference",
			sql:  "USE |",
			want: matrixDatabases("Company", "School"),
		},
		{
			name: "use database prefix keeps matching database available",
			sql:  "USE Sch|",
			want: matrixDatabases("School"),
		},
		{
			name: "next value for sequence reference",
			sql:  "SELECT NEXT VALUE FOR |",
			want: matrixSequences("EmployeeIdSeq", "OrderSeq"),
		},
		{
			name:    "next value for sequence prefix",
			sql:     "SELECT NEXT VALUE FOR Employee|",
			want:    matrixSequences("EmployeeIdSeq"),
			notWant: matrixSequences("OrderSeq"),
		},
		{
			name: "schema qualified next value for sequence reference",
			sql:  "SELECT NEXT VALUE FOR dbo.|",
			want: matrixSequences("EmployeeIdSeq", "OrderSeq"),
		},
		{
			name:    "alternate schema next value for sequence reference",
			sql:     "SELECT NEXT VALUE FOR MySchema.|",
			want:    matrixSequences("SalarySeq"),
			notWant: matrixSequences("EmployeeIdSeq"),
		},
		{
			name:    "cross database next value for sequence reference",
			sql:     "SELECT NEXT VALUE FOR School.dbo.|",
			want:    matrixSequences("StudentSeq"),
			notWant: matrixSequences("EmployeeIdSeq"),
		},
		{
			name:    "database qualified next value for schema reference",
			sql:     "SELECT NEXT VALUE FOR Company.|",
			want:    matrixSchemas("dbo", "MySchema"),
			notWant: matrixSequences("EmployeeIdSeq", "OrderSeq"),
		},
		{
			name: "truncate table reference",
			sql:  "TRUNCATE TABLE |",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "truncate cross database table reference",
			sql:  "TRUNCATE TABLE School.dbo.|",
			want: matrixTables("Student"),
		},
		{
			name: "merge target table reference",
			sql:  "MERGE INTO |",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "merge source table reference",
			sql:  "MERGE INTO Employees USING |",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "merge on uses target and source columns",
			sql:  "MERGE INTO Employees USING Address ON |",
			want: matrixColumns("Id", "Name", "EmployeeId", "Street"),
		},
		{
			name:    "merge on uses aliased target and source columns",
			sql:     "MERGE INTO Employees AS target USING Address AS source ON source.|",
			want:    matrixColumns("EmployeeId", "Street"),
			notWant: matrixColumns("Name"),
		},
		{
			name: "merge on uses cte source columns",
			sql:  "WITH source AS (SELECT EmployeeId, Street FROM Address) MERGE INTO Employees AS target USING source ON |",
			want: matrixColumns("Id", "Name", "EmployeeId", "Street"),
		},
		{
			name:    "merge on uses subquery source columns",
			sql:     "MERGE INTO Employees AS target USING (SELECT EmployeeId, Street FROM Address) AS source ON source.|",
			want:    matrixColumns("EmployeeId", "Street"),
			notWant: matrixColumns("Name"),
		},
		{
			name: "merge when keyword",
			sql:  "MERGE INTO Employees USING Address ON Employees.Id = Address.EmployeeId WHEN |",
			want: matrixKeywords("MATCHED", "NOT"),
		},
		{
			name: "merge matched update set uses target columns",
			sql:  "MERGE INTO Employees USING Address ON Employees.Id = Address.EmployeeId WHEN MATCHED THEN UPDATE SET |",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "merge matched update value uses source columns",
			sql:  "MERGE INTO Employees USING Address ON Employees.Id = Address.EmployeeId WHEN MATCHED THEN UPDATE SET Name = Address.|",
			want: matrixColumns("EmployeeId", "Street"),
		},
		{
			name: "merge output inserted columns",
			sql:  "MERGE INTO Employees USING Address ON Employees.Id = Address.EmployeeId WHEN MATCHED THEN UPDATE SET Name = Address.Street OUTPUT INSERTED.|",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "merge output deleted columns",
			sql:  "MERGE INTO Employees USING Address ON Employees.Id = Address.EmployeeId WHEN MATCHED THEN UPDATE SET Name = Address.Street OUTPUT DELETED.|",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "table hint keyword",
			sql:  "SELECT * FROM Employees WITH (|)",
			want: matrixKeywords("NOLOCK", "UPDLOCK"),
		},
		{
			name: "table hint keyword after comma",
			sql:  "SELECT * FROM Employees WITH (NOLOCK, |)",
			want: matrixKeywords("UPDLOCK", "TABLOCK"),
		},
		{
			name: "table hint keyword after multiple hints",
			sql:  "SELECT * FROM Employees WITH (NOLOCK, ROWLOCK, |)",
			want: matrixKeywords("HOLDLOCK", "TABLOCKX"),
		},
		{
			name: "join table hint keeps joined alias columns",
			sql:  "SELECT * FROM Employees e JOIN Address a WITH (NOLOCK) ON a.|",
			want: matrixColumns("EmployeeId", "Street"),
		},
		{
			name: "for xml mode",
			sql:  "SELECT * FROM Employees FOR XML |",
			want: matrixKeywords("PATH", "RAW", "AUTO", "EXPLICIT"),
		},
		{
			name: "for json mode",
			sql:  "SELECT * FROM Employees FOR JSON |",
			want: matrixKeywords("PATH", "AUTO"),
		},
		{
			name: "openjson with column type name",
			sql:  "SELECT * FROM OPENJSON(@json) WITH (Name |)",
			want: matrixTypes("INT", "NVARCHAR"),
		},
		{
			name: "openjson with second column type name",
			sql:  "SELECT * FROM OPENJSON(@json) WITH (Name NVARCHAR(100), Age |)",
			want: matrixTypes("INT", "NVARCHAR"),
		},
		{
			name: "pivot in list uses source columns",
			sql:  "SELECT * FROM Employees PIVOT (COUNT(Id) FOR Name IN (|)) p",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "unpivot in list uses source columns",
			sql:  "SELECT * FROM Employees UNPIVOT (Value FOR Field IN (|)) u",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "execute procedure reference",
			sql:  "EXEC |",
			want: matrixProcedures("SyncEmployees"),
		},
		{
			name: "execute procedure prefix",
			sql:  "EXEC Sync|",
			want: matrixProcedures("SyncEmployees"),
		},
		{
			name:    "execute database qualified schema reference",
			sql:     "EXEC Company.|",
			want:    matrixSchemas("dbo", "MySchema"),
			notWant: matrixProcedures("SyncEmployees", "SyncSalary"),
		},
		{
			name:    "alternate schema execute procedure reference",
			sql:     "EXEC MySchema.|",
			want:    matrixProcedures("SyncSalary"),
			notWant: matrixProcedures("SyncEmployees"),
		},
		{
			name:    "cross database execute procedure reference",
			sql:     "EXEC School.dbo.|",
			want:    matrixProcedures("SyncStudents"),
			notWant: matrixProcedures("SyncEmployees"),
		},
		{
			name: "option query hint",
			sql:  "SELECT * FROM Employees OPTION (|)",
			want: matrixKeywords("RECOMPILE", "MAXDOP"),
		},
		{
			name: "option query hint after comma",
			sql:  "SELECT * FROM Employees OPTION (HASH JOIN, |)",
			want: matrixKeywords("RECOMPILE", "MAXDOP"),
		},
		{
			name: "multi statement skips earlier invalid sql",
			sql:  "SELECT FROM broken; SELECT | FROM Employees",
			want: matrixColumns("Id", "Name"),
		},
		{
			name: "multi statement update after select",
			sql:  "SELECT 1; UPDATE | SET Name = 'a'",
			want: matrixTables("Address", "Employees"),
		},
		{
			name: "go batch separator starts new table reference",
			sql:  "SELECT 1\nGO\nSELECT * FROM |",
			want: matrixTables("Address", "Employees"),
		},
	}

	cfg := mssqlMatrixSessionConfig()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stmt, line, col := catchCaretUTF16(t, tc.sql)
			got := CompleteSession(context.Background(), cfg, stmt, line, col)
			for _, want := range tc.want {
				require.Truef(t, hasSessionCandidate(got, want.text, want.typ), "missing %v candidate %q in %#v", want.typ, want.text, got)
			}
			for _, notWant := range tc.notWant {
				require.Falsef(t, hasSessionCandidate(got, notWant.text, notWant.typ), "unexpected %v candidate %q in %#v", notWant.typ, notWant.text, got)
			}
		})
	}

	// ALTER TABLE column completion scopes to target table and preserves column Definition.
	stmtAlter, lineA, colA := catchCaretUTF16(t, "ALTER TABLE Employees ADD CONSTRAINT pk PRIMARY KEY (|")
	gotAlter := CompleteSession(context.Background(), cfg, stmtAlter, lineA, colA)
	require.Contains(t, gotAlter, Candidate{
		Text:       "Name",
		Type:       CandidateColumn,
		Definition: "Company.dbo.Employees | varchar, NOT NULL",
	})
	require.False(t, hasSessionCandidate(gotAlter, "Street", CandidateColumn))
}

func matrixColumns(texts ...string) []matrixCandidateSpec {
	return matrixSpecs(CandidateColumn, texts...)
}

func matrixTables(texts ...string) []matrixCandidateSpec {
	return matrixSpecs(CandidateTable, texts...)
}

func matrixViews(texts ...string) []matrixCandidateSpec {
	return matrixSpecs(CandidateView, texts...)
}

func matrixSchemas(texts ...string) []matrixCandidateSpec {
	return matrixSpecs(CandidateSchema, texts...)
}

func matrixDatabases(texts ...string) []matrixCandidateSpec {
	return matrixSpecs(CandidateDatabase, texts...)
}

func matrixSequences(texts ...string) []matrixCandidateSpec {
	return matrixSpecs(CandidateSequence, texts...)
}

func matrixKeywords(texts ...string) []matrixCandidateSpec {
	return matrixSpecs(CandidateKeyword, texts...)
}

func matrixTypes(texts ...string) []matrixCandidateSpec {
	return matrixSpecs(CandidateType_, texts...)
}

func matrixProcedures(texts ...string) []matrixCandidateSpec {
	return matrixSpecs(CandidateProcedure, texts...)
}

func matrixSpecs(typ CandidateType, texts ...string) []matrixCandidateSpec {
	specs := make([]matrixCandidateSpec, 0, len(texts))
	for _, text := range texts {
		specs = append(specs, matrixCandidateSpec{text: text, typ: typ})
	}
	return specs
}
