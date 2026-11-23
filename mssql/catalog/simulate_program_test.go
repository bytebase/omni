package catalog

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql/internal/testcontainer"
)

// TestProgramSimulation holds the catalog to testdata/programs/simulate.tsv:
// stored procedures, scalar, inline and multi-statement table-valued
// functions, and DML and DDL triggers with real bodies, in programFeatures,
// run on programs@objects, with what SQL Server 2022 did and what Exec does.
// testdata/simulate.tsv runs the grammar's program statements, whose bodies
// name nothing; these hold the catalog to how SQL Server binds a body, or
// defers binding it, when it creates the program.
//
// SQL Server's side is on the test's tape, so the test starts no server.
// OMNI_RECORD=1 runs every row on the pinned SQL Server again; -update-compat
// commits the rows a catalog fix moved; and -simulate-show=<text> prints both
// schemas of the rows that hold the text.
func TestProgramSimulation(t *testing.T) {
	rows := make([]simulate.Generated, 0, len(programFeatures))
	seen := map[string]bool{}
	for _, f := range programFeatures {
		if seen[f.statement] {
			t.Fatalf("programFeatures holds %q twice", f.statement)
		}
		seen[f.statement] = true
		rows = append(rows, simulate.Generated{Fixture: "programs@objects", Rule: "program:" + f.kind, Statement: f.statement})
	}
	simulate.Check(t, simulate.Suite{
		Path:     filepath.Join("testdata", "programs", "simulate.tsv"),
		Pins:     testcontainer.Pins,
		Fixtures: filepath.Join("testdata", "fixtures"),
		Rows:     rows,
		MadeUp:   testcontainer.Codec{}.MadeUp(),
		Omni:     simulationOmni,
		Engine:   simulationEngineWith(programResetSQL),
	})
}

// programResetSQL drops the schema-bound functions and views a row created,
// which would keep _sim_reset from dropping the tables they read.
const programResetSQL = `DECLARE @sql NVARCHAR(MAX) = N'';
SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
IF @sql <> N'' EXEC sp_executesql @sql;`

// programFeatures are the programs the suite creates, alters and drops,
// labelled program:<kind>. Some are refused by SQL Server on purpose.
var programFeatures = []struct{ kind, statement string }{
	// Procedures.
	{"procedure", "CREATE PROCEDURE p AS SELECT id, name FROM customers"},
	{"procedure", "CREATE PROCEDURE p @id INT AS SELECT name FROM customers WHERE id = @id"},
	{"procedure", "CREATE PROCEDURE p @id INT = NULL, @name NVARCHAR(100) OUTPUT AS SELECT @name = name FROM customers WHERE id = @id"},
	{"procedure", "CREATE PROCEDURE p AS SELECT nope FROM customers"},
	{"procedure", "CREATE PROCEDURE p AS SELECT nope FROM missing_t"},
	{"procedure", "CREATE PROCEDURE p AS SELECT c.nope FROM customers c JOIN missing_t m ON m.id = c.id"},
	{"procedure", "CREATE PROCEDURE p AS BEGIN SELECT nope FROM missing_t; SELECT nope FROM customers; END"},
	{"procedure", "CREATE PROCEDURE p AS BEGIN SET NOCOUNT ON; INSERT INTO orders (customer_id, amount) VALUES (1, 2.5); SELECT SCOPE_IDENTITY() AS id; END"},
	{"procedure", "CREATE PROCEDURE p AS INSERT INTO orders (nope) VALUES (1)"},
	{"procedure", "CREATE PROCEDURE p AS INSERT INTO orders (customer_id) VALUES (1, 2)"},
	{"procedure", "CREATE PROCEDURE p AS INSERT INTO orders (customer_id, amount) SELECT id FROM customers"},
	{"procedure", "CREATE PROCEDURE p AS INSERT INTO missing_t (x) VALUES (1)"},
	{"procedure", "CREATE PROCEDURE p AS UPDATE orders SET status = 'x' WHERE id = 1"},
	{"procedure", "CREATE PROCEDURE p AS UPDATE orders SET nope = 1"},
	{"procedure", "CREATE PROCEDURE p AS UPDATE o SET status = c.region FROM orders o JOIN customers c ON c.id = o.customer_id"},
	{"procedure", "CREATE PROCEDURE p AS DELETE FROM orders WHERE nope = 1"},
	{"procedure", "CREATE PROCEDURE p AS DELETE o FROM orders o JOIN customers c ON c.id = o.customer_id WHERE c.region = 'EU'"},
	{"procedure", "CREATE PROCEDURE p AS MERGE orders AS t USING customers AS s ON t.customer_id = s.id WHEN MATCHED THEN UPDATE SET t.status = 'x' WHEN NOT MATCHED THEN INSERT (customer_id) VALUES (s.id);"},
	{"procedure", "CREATE PROCEDURE p AS MERGE orders AS t USING customers AS s ON t.customer_id = s.nope WHEN MATCHED THEN DELETE;"},
	{"procedure", "CREATE PROCEDURE p AS BEGIN DECLARE @n INT; SET @n = (SELECT COUNT(*) FROM orders); IF @n > 10 SELECT 'many' AS r ELSE SELECT 'few' AS r; END"},
	{"procedure", "CREATE PROCEDURE p AS BEGIN DECLARE @i INT = 0; WHILE @i < 10 BEGIN SET @i += 1; END; RETURN @i; END"},
	{"procedure", "CREATE PROCEDURE p AS SELECT @undeclared AS v"},
	{"procedure", "CREATE PROCEDURE p AS BEGIN SET @x = 1; END"},
	{"procedure", "CREATE PROCEDURE p AS BEGIN SELECT @n AS v; DECLARE @n INT; END"},
	{"procedure", "CREATE PROCEDURE p AS BEGIN TRY SELECT 1 / 0 AS x; END TRY BEGIN CATCH SELECT ERROR_MESSAGE() AS m; THROW; END CATCH"},
	{"procedure", "CREATE PROCEDURE p AS BEGIN CREATE TABLE #t (x INT); INSERT INTO #t VALUES (1); SELECT x FROM #t; END"},
	{"procedure", "CREATE PROCEDURE p AS BEGIN CREATE TABLE #t (x INT); SELECT y FROM #t; END"},
	{"procedure", "CREATE PROCEDURE p AS BEGIN DECLARE @t TABLE (x INT); INSERT INTO @t VALUES (1); SELECT x FROM @t; END"},
	{"procedure", "CREATE PROCEDURE p AS BEGIN DECLARE @t TABLE (x INT); SELECT y FROM @t; END"},
	{"procedure", "CREATE PROCEDURE p AS SELECT x FROM @t"},
	{"procedure", "CREATE PROCEDURE p AS EXEC p_existing 1"},
	{"procedure", "CREATE PROCEDURE p AS EXEC nope_proc"},
	{"procedure", "CREATE PROCEDURE p AS SELECT dbo.f_double(id) AS d FROM customers"},
	{"procedure", "CREATE PROCEDURE p AS SELECT dbo.nope(id) AS d FROM customers"},
	{"procedure", "CREATE PROCEDURE p AS SELECT dbo.f_double(id, 1) AS d FROM customers"},
	{"procedure", "CREATE PROCEDURE p AS SELECT id FROM dbo.f_orders(1)"},
	{"procedure", "CREATE PROCEDURE p AS SELECT id FROM dbo.f_orders(1, 2)"},
	{"procedure", "CREATE PROCEDURE p AS SELECT * FROM v_customers WHERE nope = 1"},
	{"procedure", "CREATE PROCEDURE p @id INT AS SELECT name FROM customers WHERE id = @id ORDER BY nope"},
	{"procedure", "CREATE PROCEDURE p AS SELECT id INTO #copy FROM customers"},
	{"procedure", "CREATE PROCEDURE p AS SELECT customer_id, status, COUNT(*) AS n FROM orders GROUP BY customer_id"},
	{"procedure", "CREATE PROCEDURE p AS SELECT id FROM customers UNION SELECT id, status FROM orders"},
	{"procedure", "CREATE PROCEDURE p AS SELECT id FROM customers c JOIN orders o ON o.customer_id = c.id"},
	{"procedure", "CREATE PROCEDURE p @c CURSOR VARYING OUTPUT AS BEGIN SET @c = CURSOR FOR SELECT id FROM customers; OPEN @c; END"},
	{"procedure", "CREATE PROCEDURE p AS BEGIN DECLARE c CURSOR FOR SELECT nope FROM customers; OPEN c; END"},
	{"procedure", "CREATE PROCEDURE p AS WITH x AS (SELECT id FROM customers) SELECT id FROM x"},
	{"procedure", "CREATE PROCEDURE p AS TRUNCATE TABLE audit"},
	{"procedure", "CREATE PROCEDURE p AS BEGIN BEGIN TRANSACTION; UPDATE orders SET amount = amount * 2; COMMIT; END"},
	{"procedure", "CREATE PROCEDURE p @n INT AS SELECT TOP (@n) id FROM orders ORDER BY id"},
	{"procedure", "CREATE PROCEDURE p AS SELECT id FROM customers WHERE id = @@ROWCOUNT"},
	{"procedure", "CREATE PROCEDURE p AS RAISERROR('x', 16, 1)"},
	{"procedure", "CREATE PROCEDURE p AS BEGIN DECLARE @sql NVARCHAR(MAX) = N'SELECT nope FROM nowhere'; EXEC sp_executesql @sql; END"},
	{"procedure", "CREATE PROCEDURE p @id INT AS IF EXISTS (SELECT 1 FROM customers WHERE id = @id) RETURN 1 ELSE RETURN 0"},
	{"procedure", "CREATE PROCEDURE p WITH RECOMPILE AS SELECT id FROM customers"},
	{"procedure", "CREATE PROCEDURE p WITH EXECUTE AS OWNER AS SELECT id FROM customers"},
	{"procedure", "CREATE PROCEDURE s.p AS SELECT id FROM customers"},
	{"procedure", "CREATE PROCEDURE nope.p AS SELECT 1"},
	{"procedure", "CREATE PROCEDURE p_existing AS SELECT 1"},
	{"procedure", "CREATE PROCEDURE p @id NOPE AS SELECT 1"},
	{"procedure", "ALTER PROCEDURE p_existing @id INT, @extra INT = 0 AS SELECT name FROM customers WHERE id = @id"},
	{"procedure", "ALTER PROCEDURE p_existing AS SELECT nope FROM customers"},
	{"procedure", "ALTER PROCEDURE nope AS SELECT 1"},
	{"procedure", "CREATE OR ALTER PROCEDURE p_existing AS SELECT id FROM customers"},
	{"procedure", "CREATE OR ALTER PROCEDURE p AS SELECT id FROM customers"},
	{"procedure", "DROP PROCEDURE p_existing"},
	{"procedure", "DROP PROCEDURE IF EXISTS p_existing, nope"},
	// Scalar functions.
	{"scalar-function", "CREATE FUNCTION f(@a INT, @b INT) RETURNS INT AS BEGIN RETURN @a + @b END"},
	{"scalar-function", "CREATE FUNCTION f(@c INT) RETURNS DECIMAL(10,2) AS BEGIN DECLARE @s DECIMAL(10,2); SELECT @s = SUM(amount) FROM orders WHERE customer_id = @c; RETURN @s; END"},
	{"scalar-function", "CREATE FUNCTION f(@c INT) RETURNS INT AS BEGIN RETURN (SELECT COUNT(*) FROM orders WHERE nope = @c) END"},
	{"scalar-function", "CREATE FUNCTION f() RETURNS INT AS BEGIN RETURN (SELECT COUNT(*) FROM missing_t) END"},
	{"scalar-function", "CREATE FUNCTION f() RETURNS INT AS BEGIN INSERT INTO orders (customer_id) VALUES (1); RETURN 1 END"},
	{"scalar-function", "CREATE FUNCTION f() RETURNS UNIQUEIDENTIFIER AS BEGIN RETURN NEWID() END"},
	{"scalar-function", "CREATE FUNCTION f() RETURNS INT AS BEGIN EXEC p_existing 1; RETURN 1 END"},
	{"scalar-function", "CREATE FUNCTION f() RETURNS INT AS BEGIN RETURN @x END"},
	{"scalar-function", "CREATE FUNCTION f() RETURNS INT AS BEGIN SELECT id FROM customers; RETURN 1 END"},
	{"scalar-function", "CREATE FUNCTION f() RETURNS INT AS BEGIN DECLARE @t TABLE (x INT); INSERT INTO @t VALUES (1); RETURN (SELECT MAX(x) FROM @t) END"},
	{"scalar-function", "CREATE FUNCTION f() RETURNS INT WITH SCHEMABINDING AS BEGIN RETURN (SELECT COUNT(*) FROM dbo.orders) END"},
	{"scalar-function", "CREATE FUNCTION f() RETURNS INT WITH SCHEMABINDING AS BEGIN RETURN (SELECT COUNT(*) FROM orders) END"},
	{"scalar-function", "CREATE FUNCTION f(@x INT) RETURNS INT AS BEGIN IF @x > 0 RETURN 1; RETURN 0; END"},
	{"scalar-function", "CREATE FUNCTION f(@x INT) RETURNS INT AS BEGIN IF @x > 0 RETURN 1; END"},
	{"scalar-function", "CREATE FUNCTION f() RETURNS INT AS BEGIN RETURN dbo.f_double(2) END"},
	{"scalar-function", "CREATE FUNCTION f() RETURNS INT AS BEGIN RETURN dbo.nope(2) END"},
	{"scalar-function", "CREATE FUNCTION f() RETURNS INT AS BEGIN RETURN (SELECT TOP 1 id FROM dbo.f_orders(1)) END"},
	{"scalar-function", "CREATE FUNCTION f() RETURNS NOPE AS BEGIN RETURN 1 END"},
	{"scalar-function", "CREATE FUNCTION f(@x INT) RETURNS INT AS BEGIN DECLARE @y INT; SET @y = @x * 2; RETURN @y END"},
	{"scalar-function", "CREATE FUNCTION f(@x INT) RETURNS INT WITH RETURNS NULL ON NULL INPUT AS BEGIN RETURN @x END"},
	{"scalar-function", "CREATE FUNCTION f_double(@x INT) RETURNS INT AS BEGIN RETURN @x END"},
	{"scalar-function", "ALTER FUNCTION f_double(@x INT) RETURNS INT AS BEGIN RETURN @x * 3 END"},
	{"scalar-function", "ALTER FUNCTION f_double(@x INT) RETURNS TABLE AS RETURN SELECT @x AS x"},
	{"scalar-function", "CREATE OR ALTER FUNCTION f_double(@x INT) RETURNS INT AS BEGIN RETURN @x + 1 END"},
	{"scalar-function", "DROP FUNCTION f_double"},
	// Inline table-valued functions.
	{"inline-function", "CREATE FUNCTION f(@c INT) RETURNS TABLE AS RETURN SELECT id, amount FROM orders WHERE customer_id = @c"},
	{"inline-function", "CREATE FUNCTION f(@c INT) RETURNS TABLE AS RETURN (SELECT id, amount FROM orders WHERE customer_id = @c)"},
	{"inline-function", "CREATE FUNCTION f() RETURNS TABLE AS RETURN SELECT nope FROM orders"},
	{"inline-function", "CREATE FUNCTION f() RETURNS TABLE AS RETURN SELECT id FROM missing_t"},
	{"inline-function", "CREATE FUNCTION f() RETURNS TABLE AS RETURN SELECT COUNT(*) FROM orders"},
	{"inline-function", "CREATE FUNCTION f() RETURNS TABLE AS RETURN SELECT id, id FROM orders"},
	{"inline-function", "CREATE FUNCTION f() RETURNS TABLE AS RETURN SELECT id FROM orders ORDER BY id"},
	{"inline-function", "CREATE FUNCTION f() RETURNS TABLE AS RETURN SELECT TOP 5 id FROM orders ORDER BY id"},
	{"inline-function", "CREATE FUNCTION f() RETURNS TABLE AS RETURN WITH x AS (SELECT id FROM orders) SELECT id FROM x"},
	{"inline-function", "CREATE FUNCTION f() RETURNS TABLE AS RETURN SELECT @x AS v"},
	{"inline-function", "CREATE FUNCTION f(@c INT) RETURNS TABLE WITH SCHEMABINDING AS RETURN SELECT id FROM dbo.orders WHERE customer_id = @c"},
	{"inline-function", "CREATE FUNCTION f(@c INT) RETURNS TABLE WITH SCHEMABINDING AS RETURN SELECT id FROM orders WHERE customer_id = @c"},
	{"inline-function", "CREATE FUNCTION f(@c INT) RETURNS TABLE AS RETURN SELECT o.id, x.id AS xid FROM orders o CROSS APPLY dbo.f_orders(@c) x"},
	{"inline-function", "CREATE FUNCTION f() RETURNS TABLE AS RETURN SELECT id FROM v_customers"},
	{"inline-function", "ALTER FUNCTION f_orders(@c INT) RETURNS TABLE AS RETURN SELECT id FROM orders WHERE customer_id = @c"},
	{"inline-function", "ALTER FUNCTION f_orders(@c INT) RETURNS TABLE AS RETURN SELECT nope FROM orders"},
	{"inline-function", "DROP FUNCTION f_orders, f_double"},
	// Multi-statement table-valued functions.
	{"multi-statement-function", "CREATE FUNCTION f(@c INT) RETURNS @r TABLE (id INT, amount DECIMAL(10,2)) AS BEGIN INSERT INTO @r SELECT id, amount FROM orders WHERE customer_id = @c; RETURN; END"},
	{"multi-statement-function", "CREATE FUNCTION f() RETURNS @r TABLE (id INT) AS BEGIN INSERT INTO @r SELECT nope FROM orders; RETURN; END"},
	{"multi-statement-function", "CREATE FUNCTION f() RETURNS @r TABLE (id INT) AS BEGIN INSERT INTO @r (nope) VALUES (1); RETURN; END"},
	{"multi-statement-function", "CREATE FUNCTION f() RETURNS @r TABLE (id INT, id INT) AS BEGIN RETURN; END"},
	{"multi-statement-function", "CREATE FUNCTION f() RETURNS @r TABLE (id INT PRIMARY KEY, v NVARCHAR(10) DEFAULT 'x') AS BEGIN INSERT INTO @r (id) VALUES (1); RETURN; END"},
	{"multi-statement-function", "CREATE FUNCTION f() RETURNS @r TABLE (id INT) AS BEGIN INSERT INTO @r VALUES (1); RETURN 1; END"},
	{"multi-statement-function", "CREATE FUNCTION f() RETURNS @r TABLE (id INT) AS BEGIN INSERT INTO @r SELECT id FROM missing_t; RETURN; END"},
	{"multi-statement-function", "CREATE FUNCTION f() RETURNS @r TABLE (id NOPE) AS BEGIN RETURN; END"},
	{"multi-statement-function", "CREATE FUNCTION f() RETURNS @r TABLE (id INT) AS BEGIN UPDATE @r SET id = id + 1; DELETE FROM @r WHERE id > 5; RETURN; END"},
	{"multi-statement-function", "CREATE FUNCTION f() RETURNS @r TABLE (id INT) AS BEGIN UPDATE orders SET amount = 0; RETURN; END"},
	// DML triggers.
	{"trigger", "CREATE TRIGGER t ON orders AFTER INSERT AS INSERT INTO audit (what) SELECT 'order ' + CAST(id AS NVARCHAR(10)) FROM inserted"},
	{"trigger", "CREATE TRIGGER t ON orders AFTER UPDATE AS IF UPDATE(amount) INSERT INTO audit (what) SELECT 'amount' FROM inserted i JOIN deleted d ON d.id = i.id"},
	{"trigger", "CREATE TRIGGER t ON orders AFTER DELETE AS SELECT nope FROM deleted"},
	{"trigger", "CREATE TRIGGER t ON orders AFTER INSERT AS SELECT nope FROM inserted"},
	{"trigger", "CREATE TRIGGER t ON orders AFTER INSERT AS IF UPDATE(nope) RETURN"},
	{"trigger", "CREATE TRIGGER t ON orders INSTEAD OF DELETE AS UPDATE o SET status = 'deleted' FROM orders o JOIN deleted d ON d.id = o.id"},
	{"trigger", "CREATE TRIGGER t ON v_customers INSTEAD OF INSERT AS INSERT INTO customers (id, name) SELECT id, name FROM inserted"},
	{"trigger", "CREATE TRIGGER t ON v_customers INSTEAD OF INSERT AS INSERT INTO customers (id, name) SELECT id, region FROM inserted"},
	{"trigger", "CREATE TRIGGER t ON v_customers AFTER INSERT AS SET NOCOUNT ON"},
	{"trigger", "CREATE TRIGGER t ON orders FOR INSERT, UPDATE AS SET NOCOUNT ON"},
	{"trigger", "CREATE TRIGGER t ON orders AFTER INSERT AS BEGIN DECLARE @n INT = (SELECT COUNT(*) FROM inserted); IF @n > 1 THROW 50000, 'one at a time', 1; END"},
	{"trigger", "CREATE TRIGGER t ON orders AFTER INSERT AS SELECT COLUMNS_UPDATED() AS c"},
	{"trigger", "CREATE TRIGGER t ON orders AFTER INSERT AS SELECT id FROM missing_t"},
	{"trigger", "CREATE TRIGGER t ON orders AFTER INSERT AS SELECT @v AS v"},
	{"trigger", "CREATE TRIGGER t ON orders WITH EXECUTE AS OWNER AFTER INSERT AS SET NOCOUNT ON"},
	{"trigger", "CREATE TRIGGER t ON orders AFTER INSERT NOT FOR REPLICATION AS SET NOCOUNT ON"},
	{"trigger", "CREATE TRIGGER t ON orders AFTER INSERT AS ROLLBACK TRANSACTION"},
	{"trigger", "CREATE TRIGGER t ON missing_t AFTER INSERT AS SET NOCOUNT ON"},
	{"trigger", "CREATE TRIGGER s.t ON orders AFTER INSERT AS SET NOCOUNT ON"},
	{"trigger", "CREATE TRIGGER trg_existing ON orders AFTER INSERT AS SET NOCOUNT ON"},
	{"trigger", "CREATE TRIGGER trg_existing ON customers AFTER INSERT AS SET NOCOUNT ON"},
	{"trigger", "ALTER TRIGGER trg_existing ON orders AFTER UPDATE AS SET NOCOUNT ON"},
	{"trigger", "ALTER TRIGGER trg_existing ON customers AFTER UPDATE AS SET NOCOUNT ON"},
	{"trigger", "CREATE OR ALTER TRIGGER trg_existing ON orders AFTER DELETE AS SET NOCOUNT ON"},
	{"trigger", "DROP TRIGGER trg_existing"},
	{"trigger", "DROP TRIGGER IF EXISTS nope"},
	// DDL triggers.
	{"ddl-trigger", "CREATE TRIGGER t ON DATABASE FOR CREATE_TABLE AS PRINT 'created'"},
	{"ddl-trigger", "CREATE TRIGGER t ON DATABASE FOR DDL_TABLE_EVENTS AS SELECT EVENTDATA().value('(/EVENT_INSTANCE/ObjectName)[1]', 'nvarchar(128)') AS n"},
	{"ddl-trigger", "CREATE TRIGGER t ON DATABASE FOR NOPE_EVENT AS PRINT 'x'"},
	{"ddl-trigger", "CREATE TRIGGER t ON DATABASE AFTER DROP_TABLE AS ROLLBACK"},
	{"ddl-trigger", "CREATE TRIGGER t ON DATABASE FOR CREATE_TABLE AS SELECT nope FROM customers"},
	{"ddl-trigger", "CREATE TRIGGER t ON DATABASE FOR CREATE_TABLE AS INSERT INTO audit (what) VALUES (EVENTDATA().value('(/EVENT_INSTANCE/ObjectName)[1]', 'nvarchar(100)'))"},
	{"ddl-trigger", "DROP TRIGGER nope ON DATABASE"},
}

// TestProgramBinding pins module bindings beyond the simulated rows.
func TestProgramBinding(t *testing.T) {
	cases := []struct {
		name  string
		setup []string
		sql   string
		code  int // 0 when SQL Server creates the module
	}{
		{"CTE as an UPDATE's target", nil, "CREATE PROCEDURE p AS WITH x AS (SELECT id FROM t) UPDATE x SET id = 1", 0},
		{"CTE as an INSERT's source", nil, "CREATE PROCEDURE p AS WITH x AS (SELECT id FROM t) INSERT INTO t (id) SELECT id FROM x", 0},
		{"INSERT skips a computed column", nil, "CREATE PROCEDURE p AS INSERT INTO c VALUES (1)", 0},
		{"INSERT with a value too many", nil, "CREATE PROCEDURE p AS INSERT INTO c VALUES (1, 2)", 110},
		{"temporary table defers its statement", nil, "CREATE PROCEDURE p AS SELECT nope FROM #t", 0},
		{"existing table binds its statement", nil, "CREATE PROCEDURE p AS SELECT nope FROM t", ErrInvalidColumnName},
		{"variable declared in a block is declared after it", nil, "CREATE PROCEDURE p AS BEGIN IF 1 = 1 BEGIN DECLARE @a INT; END; SET @a = 1; END", 0},
		{"UPDATE target beside its FROM sources", nil, "CREATE PROCEDURE p AS UPDATE t SET id = t.id + 1 FROM u WHERE t.id = u.n", 0},
		{"function updates its table variable by alias", nil, "CREATE FUNCTION f() RETURNS @r TABLE (id INT) AS BEGIN INSERT INTO @r VALUES (1); UPDATE r SET id = 2 FROM @r AS r; RETURN; END", 0},
		{"function updates a table by alias", nil, "CREATE FUNCTION f() RETURNS @r TABLE (id INT) AS BEGIN UPDATE r SET id = 2 FROM t AS r; RETURN; END", 443},
		{"qualified UPDATE target in its FROM", nil, "CREATE PROCEDURE p AS UPDATE dbo.t SET id = id + 1 FROM dbo.t", 0},
		{"INSERT skips a table variable's computed column", nil, "CREATE PROCEDURE p AS BEGIN DECLARE @t TABLE (id INT, doubled AS id * 2); INSERT INTO @t VALUES (1); END", 0},
		{"method receiver is bound", nil, "CREATE VIEW v AS SELECT CAST(nope AS xml).value('(/x)[1]', 'int') AS x FROM t", ErrInvalidColumnName},
		{"recursive call binds the new signature", []string{"CREATE FUNCTION dbo.f(@x INT) RETURNS INT AS BEGIN RETURN @x END"}, "ALTER FUNCTION dbo.f(@x INT, @y INT) RETURNS INT AS BEGIN RETURN CASE WHEN @x <= 0 THEN @y ELSE dbo.f(@x - 1, @y) END; END", 0},
		{"recursive call with the old signature", []string{"CREATE FUNCTION dbo.f(@x INT) RETURNS INT AS BEGIN RETURN @x END"}, "ALTER FUNCTION dbo.f(@x INT, @y INT) RETURNS INT AS BEGIN RETURN dbo.f(@x); END", 313},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := LoadSQL("CREATE TABLE t (id INT NOT NULL);\nCREATE TABLE u (n INT NOT NULL);\nCREATE TABLE c (id INT NOT NULL, d AS id * 2);")
			if err != nil {
				t.Fatal(err)
			}
			for _, stmt := range append(append([]string(nil), tc.setup...), tc.sql) {
				results, err := c.Exec(stmt, nil)
				if err != nil {
					t.Fatal(err)
				}
				var got error
				for _, r := range results {
					if r.Error != nil {
						got = r.Error
					}
				}
				if stmt != tc.sql {
					if got != nil {
						t.Fatalf("%s: %v", stmt, got)
					}
					continue
				}
				var refusal *Error
				switch {
				case tc.code == 0 && got != nil:
					t.Fatalf("%s: %v", stmt, got)
				case tc.code != 0 && (!errors.As(got, &refusal) || refusal.Code != tc.code):
					t.Fatalf("%s: got %v, want Msg %d", stmt, got, tc.code)
				}
			}
		})
	}
	// A desired-state script may define a module before the tables it reads.
	if _, err := LoadSQL("CREATE PROCEDURE p AS SELECT id FROM t;\nGO\nCREATE TABLE t (id INT NOT NULL);"); err != nil {
		t.Fatalf("LoadSQL: %v", err)
	}
}

// TestAlterTriggerPrefersExactName pins ALTER TRIGGER on a case-sensitive
// snapshot holding trg and TRG: it replaces the one its name matches exactly.
func TestAlterTriggerPrefersExactName(t *testing.T) {
	c, err := LoadSQL("CREATE TABLE t (id INT NOT NULL);")
	if err != nil {
		t.Fatal(err)
	}
	table := c.GetTable("dbo", "t")
	table.Triggers = []*metadata.TriggerMetadata{
		{Name: "trg", Timing: "AFTER", Event: "INSERT", Body: "CREATE TRIGGER trg ON t AFTER INSERT AS SET NOCOUNT ON"},
		{Name: "TRG", Timing: "AFTER", Event: "INSERT", Body: "CREATE TRIGGER TRG ON t AFTER INSERT AS SET NOCOUNT ON"},
	}
	results, err := c.Exec("ALTER TRIGGER TRG ON t AFTER UPDATE AS SET NOCOUNT ON", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatal(r.Error)
		}
	}
	if table.Triggers[0].Event != "INSERT" || table.Triggers[1].Event != "UPDATE" {
		t.Fatalf("got %s %s, %s %s", table.Triggers[0].Name, table.Triggers[0].Event, table.Triggers[1].Name, table.Triggers[1].Event)
	}
}
