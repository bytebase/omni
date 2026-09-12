package parser

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/testcontainers/testcontainers-go"
	tcmssql "github.com/testcontainers/testcontainers-go/modules/mssql"
)

type parserOracle struct {
	db  *sql.DB
	ctx context.Context
}

// sharedMSSQL is the package-wide SQL Server 2022 oracle. SQL Server takes
// 15-20s to boot, so one container serves every test in the package instead
// of one per test. The tests only ever run SET PARSEONLY ON checks, which
// execute nothing, so there is no state to isolate between them.
var sharedMSSQL struct {
	once      sync.Once
	container *tcmssql.MSSQLServerContainer
	db        *sql.DB
	err       error
}

// startParserOracle returns the shared SQL Server oracle, starting it on first
// use. It fails the test in CI when the container cannot start and skips
// locally, so a laptop without Docker still runs the rest of the package.
func startParserOracle(t *testing.T) *parserOracle {
	t.Helper()
	sharedMSSQL.once.Do(func() {
		ctx := context.Background()
		container, err := tcmssql.Run(ctx, "mcr.microsoft.com/mssql/server:2022-latest",
			tcmssql.WithAcceptEULA(),
			tcmssql.WithPassword("Str0ngPa$$w0rd!"),
		)
		if err != nil {
			sharedMSSQL.err = fmt.Errorf("start SQL Server container: %w", err)
			return
		}
		fail := func(err error) {
			_ = testcontainers.TerminateContainer(container)
			sharedMSSQL.err = err
		}
		connStr, err := container.ConnectionString(ctx)
		if err != nil {
			fail(fmt.Errorf("connection string: %w", err))
			return
		}
		db, err := sql.Open("sqlserver", connStr)
		if err != nil {
			fail(fmt.Errorf("open database: %w", err))
			return
		}
		// SET PARSEONLY and USE are session state: pin the pool to one
		// connection so every statement sees them.
		db.SetMaxOpenConns(1)
		db.SetMaxIdleConns(1)
		if err := db.PingContext(ctx); err != nil {
			_ = db.Close()
			fail(fmt.Errorf("ping SQL Server: %w", err))
			return
		}
		// Create a test database and table for queries that reference objects.
		for _, s := range []string{
			"CREATE DATABASE testdb",
			"USE testdb",
			"CREATE TABLE dbo.t (a INT, col INT, partition INT, encryption INT, window INT, bucket INT)",
		} {
			if _, err := db.ExecContext(ctx, s); err != nil {
				_ = db.Close()
				fail(fmt.Errorf("setup SQL failed (%s): %w", s, err))
				return
			}
		}
		sharedMSSQL.container, sharedMSSQL.db = container, db
	})
	if sharedMSSQL.err != nil {
		if os.Getenv("CI") != "" {
			t.Fatalf("SQL Server oracle required in CI but unavailable: %v", sharedMSSQL.err)
		}
		t.Skipf("SQL Server oracle unavailable: %v", sharedMSSQL.err)
	}
	return &parserOracle{db: sharedMSSQL.db, ctx: context.Background()}
}

func TestMain(m *testing.M) {
	code := m.Run()
	if sharedMSSQL.db != nil {
		_ = sharedMSSQL.db.Close()
	}
	if sharedMSSQL.container != nil {
		_ = testcontainers.TerminateContainer(sharedMSSQL.container)
	}
	os.Exit(code)
}

// canParse tests whether SQL Server accepts the given SQL without execution errors.
// It uses SET PARSEONLY ON to check syntax without executing.
func (o *parserOracle) canParse(sql string) (bool, error) {
	_, err := o.db.ExecContext(o.ctx, "SET PARSEONLY ON")
	if err != nil {
		return false, fmt.Errorf("SET PARSEONLY ON: %w", err)
	}
	defer o.db.ExecContext(o.ctx, "SET PARSEONLY OFF") //nolint:errcheck

	_, err = o.db.ExecContext(o.ctx, sql)
	if err != nil {
		return false, nil // Parse error — SQL Server rejects this syntax
	}
	return true, nil // SQL Server accepts this syntax
}

// TestKeywordOracleOptionPositions verifies whether omni's option parsing
// is too permissive compared to SQL Server across all major option positions.
func TestKeywordOracleOptionPositions(t *testing.T) {
	oracle := startParserOracle(t)

	tests := []struct {
		name string
		sql  string
	}{
		// ================================================================
		// 1. Index options: CREATE INDEX ix ON t(a) WITH (...)
		// ================================================================
		{"index/FILLFACTOR_valid", "CREATE INDEX ix ON dbo.t(a) WITH (FILLFACTOR = 80)"},
		{"index/SELECT_invalid", "CREATE INDEX ix ON dbo.t(a) WITH (SELECT = 1)"},
		{"index/FROM_invalid", "CREATE INDEX ix ON dbo.t(a) WITH (FROM = 1)"},
		{"index/DROP_invalid", "CREATE INDEX ix ON dbo.t(a) WITH (DROP = ON)"},

		// ================================================================
		// 2. SET predicate options: SET <option> ON/OFF
		// ================================================================
		{"set_predicate/ANSI_NULLS_valid", "SET ANSI_NULLS ON"},
		{"set_predicate/QUOTED_IDENTIFIER_valid", "SET QUOTED_IDENTIFIER OFF"},
		{"set_predicate/ARITHABORT_valid", "SET ARITHABORT ON"},
		{"set_predicate/SELECT_invalid", "SET SELECT ON"},
		{"set_predicate/FROM_invalid", "SET FROM OFF"},

		// ================================================================
		// 3. Database SET options: ALTER DATABASE db SET <option>
		// ================================================================
		{"db_set/RECOVERY_SIMPLE_valid", "ALTER DATABASE testdb SET RECOVERY SIMPLE"},
		{"db_set/READ_ONLY_valid", "ALTER DATABASE testdb SET READ_ONLY"},
		{"db_set/ANSI_NULLS_valid", "ALTER DATABASE testdb SET ANSI_NULLS ON"},
		{"db_set/SELECT_invalid", "ALTER DATABASE testdb SET SELECT ON"},
		{"db_set/FROM_invalid", "ALTER DATABASE testdb SET FROM OFF"},

		// ================================================================
		// 4. Table hints: SELECT * FROM t WITH (<hint>)
		// ================================================================
		{"table_hint/NOLOCK_valid", "SELECT * FROM dbo.t WITH (NOLOCK)"},
		{"table_hint/ROWLOCK_valid", "SELECT * FROM dbo.t WITH (ROWLOCK)"},
		{"table_hint/HOLDLOCK_valid", "SELECT * FROM dbo.t WITH (HOLDLOCK)"},
		{"table_hint/SELECT_invalid", "SELECT * FROM dbo.t WITH (SELECT)"},
		{"table_hint/FROM_invalid", "SELECT * FROM dbo.t WITH (FROM)"},

		// ================================================================
		// 5. Query hints: SELECT ... OPTION (<hint>)
		// ================================================================
		{"query_hint/RECOMPILE_valid", "SELECT * FROM dbo.t OPTION (RECOMPILE)"},
		{"query_hint/MAXDOP_valid", "SELECT * FROM dbo.t OPTION (MAXDOP 1)"},
		{"query_hint/SELECT_invalid", "SELECT * FROM dbo.t OPTION (SELECT)"},
		{"query_hint/FROM_invalid", "SELECT * FROM dbo.t OPTION (FROM)"},

		// ================================================================
		// 6. Cursor options: DECLARE c CURSOR <options> FOR SELECT 1
		// ================================================================
		{"cursor/FAST_FORWARD_valid", "DECLARE c CURSOR FAST_FORWARD FOR SELECT 1"},
		{"cursor/SCROLL_valid", "DECLARE c CURSOR SCROLL FOR SELECT 1"},
		{"cursor/STATIC_valid", "DECLARE c CURSOR STATIC FOR SELECT 1"},
		{"cursor/SELECT_invalid", "DECLARE c CURSOR SELECT FOR SELECT 1"},
		{"cursor/FROM_invalid", "DECLARE c CURSOR FROM FOR SELECT 1"},

		// ================================================================
		// 7. DBCC options: DBCC CHECKDB WITH (<option>)
		// ================================================================
		{"dbcc/NO_INFOMSGS_valid", "DBCC CHECKDB WITH (NO_INFOMSGS)"},
		{"dbcc/ALL_ERRORMSGS_valid", "DBCC CHECKDB WITH (ALL_ERRORMSGS)"},
		{"dbcc/SELECT_invalid", "DBCC CHECKDB WITH (SELECT)"},
		{"dbcc/FROM_invalid", "DBCC CHECKDB WITH (FROM)"},

		// ================================================================
		// 8. Backup options: BACKUP DATABASE ... WITH (<option>)
		// ================================================================
		{"backup/COMPRESSION_valid", "BACKUP DATABASE testdb TO DISK = '/tmp/test.bak' WITH COMPRESSION"},
		{"backup/INIT_valid", "BACKUP DATABASE testdb TO DISK = '/tmp/test.bak' WITH INIT"},
		{"backup/SELECT_invalid", "BACKUP DATABASE testdb TO DISK = '/tmp/test.bak' WITH SELECT"},
		{"backup/FROM_invalid", "BACKUP DATABASE testdb TO DISK = '/tmp/test.bak' WITH FROM"},

		// ================================================================
		// 9. Restore options: RESTORE DATABASE ... WITH (<option>)
		// ================================================================
		{"restore/NORECOVERY_valid", "RESTORE DATABASE testdb FROM DISK = '/tmp/test.bak' WITH NORECOVERY"},
		{"restore/REPLACE_valid", "RESTORE DATABASE testdb FROM DISK = '/tmp/test.bak' WITH REPLACE"},
		{"restore/SELECT_invalid", "RESTORE DATABASE testdb FROM DISK = '/tmp/test.bak' WITH SELECT"},
		{"restore/FROM_invalid", "RESTORE DATABASE testdb FROM DISK = '/tmp/test.bak' WITH FROM"},

		// ================================================================
		// 10. FOR XML options: SELECT 1 FOR XML RAW, <option>
		// ================================================================
		{"for_xml/ELEMENTS_valid", "SELECT 1 AS val FOR XML RAW, ELEMENTS"},
		{"for_xml/ROOT_valid", "SELECT 1 AS val FOR XML RAW, ROOT('r')"},
		{"for_xml/SELECT_invalid", "SELECT 1 AS val FOR XML RAW, SELECT"},
		{"for_xml/FROM_invalid", "SELECT 1 AS val FOR XML RAW, FROM"},

		// ================================================================
		// 11. FOR JSON options: SELECT 1 FOR JSON PATH, <option>
		// ================================================================
		{"for_json/ROOT_valid", "SELECT 1 AS val FOR JSON PATH, ROOT('r')"},
		{"for_json/WITHOUT_ARRAY_WRAPPER_valid", "SELECT 1 AS val FOR JSON PATH, WITHOUT_ARRAY_WRAPPER"},
		{"for_json/SELECT_invalid", "SELECT 1 AS val FOR JSON PATH, SELECT"},
		{"for_json/FROM_invalid", "SELECT 1 AS val FOR JSON PATH, FROM"},

		// ================================================================
		// 12. Bulk insert options: BULK INSERT t FROM 'x' WITH (<option>)
		// NOTE: SQL Server's SET PARSEONLY ON rejects all BULK INSERT statements
		// regardless of option validity, so "valid" cases cannot be oracle-tested.
		// Only invalid-option cases are tested (both omni and SQL Server reject).
		// ================================================================
		{"bulk_insert/SELECT_invalid", "BULK INSERT dbo.t FROM '/tmp/data.csv' WITH (SELECT = 1)"},
		{"bulk_insert/FROM_invalid", "BULK INSERT dbo.t FROM '/tmp/data.csv' WITH (FROM = 1)"},

		// ================================================================
		// 13. CREATE TABLE options: WITH (<option>)
		// Note: MEMORY_OPTIMIZED/DURABILITY are syntactically valid but PARSEONLY
		// rejects them (requires In-Memory OLTP engine), so only invalid cases tested.
		// ================================================================
		{"create_table/SELECT_invalid", "CREATE TABLE dbo.t_bad1 (a INT) WITH (SELECT = ON)"},
		{"create_table/FROM_invalid", "CREATE TABLE dbo.t_bad2 (a INT) WITH (FROM = ON)"},

		// ================================================================
		// 14. Procedure options: CREATE PROC p WITH <option> AS SELECT 1
		// ================================================================
		{"proc/RECOMPILE_valid", "CREATE PROCEDURE dbo.p_recomp WITH RECOMPILE AS SELECT 1"},
		{"proc/ENCRYPTION_valid", "CREATE PROCEDURE dbo.p_enc WITH ENCRYPTION AS SELECT 1"},
		{"proc/SELECT_invalid", "CREATE PROCEDURE dbo.p_bad1 WITH SELECT AS SELECT 1"},
		{"proc/FROM_invalid", "CREATE PROCEDURE dbo.p_bad2 WITH FROM AS SELECT 1"},

		// ================================================================
		// 15. Trigger options: CREATE TRIGGER tr ON t WITH <option> FOR INSERT AS ...
		// ================================================================
		{"trigger/ENCRYPTION_valid", "CREATE TRIGGER dbo.tr_enc ON dbo.t WITH ENCRYPTION FOR INSERT AS SELECT 1"},
		{"trigger/SELECT_invalid", "CREATE TRIGGER dbo.tr_bad1 ON dbo.t WITH SELECT FOR INSERT AS SELECT 1"},
		{"trigger/FROM_invalid", "CREATE TRIGGER dbo.tr_bad2 ON dbo.t WITH FROM FOR INSERT AS SELECT 1"},

		// ================================================================
		// 16. View options: CREATE VIEW v WITH <option> AS SELECT 1
		// ================================================================
		{"view/SCHEMABINDING_valid", "CREATE VIEW dbo.v_sb WITH SCHEMABINDING AS SELECT 1 AS val"},
		{"view/ENCRYPTION_valid", "CREATE VIEW dbo.v_enc WITH ENCRYPTION AS SELECT 1 AS val"},
		{"view/SELECT_invalid", "CREATE VIEW dbo.v_bad1 WITH SELECT AS SELECT 1 AS val"},
		{"view/FROM_invalid", "CREATE VIEW dbo.v_bad2 WITH FROM AS SELECT 1 AS val"},

		// ================================================================
		// 17. Fulltext index options
		// Note: CHANGE_TRACKING is syntactically valid but PARSEONLY rejects all
		// fulltext syntax (requires fulltext catalog), so only invalid cases tested.
		// ================================================================
		{"fulltext/SELECT_invalid", "CREATE FULLTEXT INDEX ON dbo.t (a) KEY INDEX ix WITH (SELECT = ON)"},

		// ================================================================
		// 18. Service broker options
		// Note: PARSEONLY rejects all broker syntax (requires Service Broker enabled),
		// so only invalid cases tested here.
		// ================================================================
		{"broker/SELECT_invalid", "CREATE SERVICE svc ON QUEUE dbo.q (SELECT)"},

		// ================================================================
		// 19. Availability group options
		// Note: PARSEONLY rejects all AG syntax (requires WSFC/HA), so only
		// invalid cases tested here.
		// ================================================================
		{"ag/SELECT_invalid", "CREATE AVAILABILITY GROUP ag WITH (SELECT = ON)"},

		// ================================================================
		// 20. Endpoint options
		// ================================================================
		{"endpoint/STATE_valid", "CREATE ENDPOINT ep STATE = STARTED AS TCP (LISTENER_PORT = 5022) FOR DATABASE_MIRRORING (ROLE = PARTNER)"},
		{"endpoint/SELECT_invalid", "CREATE ENDPOINT ep SELECT = STARTED AS TCP (LISTENER_PORT = 5022) FOR DATABASE_MIRRORING (ROLE = PARTNER)"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ssAccepts, err := oracle.canParse(tc.sql)
			if err != nil {
				t.Fatalf("oracle error: %v", err)
			}

			_, omniErr := Parse(tc.sql)
			omniAccepts := omniErr == nil

			if ssAccepts != omniAccepts {
				t.Errorf("MISMATCH: SQL Server %s, omni %s\n  SQL: %s",
					boolToAcceptReject(ssAccepts), boolToAcceptReject(omniAccepts), tc.sql)
			}
		})
	}
}

// TestKeywordOracleContextDisambiguation verifies that context keywords
// work correctly as identifiers in SQL Server.
func TestKeywordOracleContextDisambiguation(t *testing.T) {
	oracle := startParserOracle(t)

	tests := []struct {
		name string
		sql  string
	}{
		// Context keywords as table aliases
		{"window as table alias", "SELECT * FROM dbo.t window"},
		{"window as explicit alias", "SELECT * FROM dbo.t AS window"},
		{"encryption as table alias", "SELECT * FROM dbo.t encryption"},
		// Context keywords as column names
		{"partition as column name", "SELECT partition FROM dbo.t"},
		{"partition as qualified column", "SELECT t.partition FROM dbo.t"},
		{"encryption as qualified column", "SELECT t.encryption FROM dbo.t"},
		{"window as column name", "SELECT window FROM dbo.t"},
		// Context keywords as bare aliases
		{"encryption as bare alias", "SELECT 1 encryption FROM dbo.t"},
		{"partition as bare alias", "SELECT 1 partition FROM dbo.t"},
		{"window as bare alias", "SELECT 1 window FROM dbo.t"},
		// Context keywords as explicit aliases
		{"encryption as AS alias", "SELECT 1 AS encryption FROM dbo.t"},
		{"partition as AS alias", "SELECT 1 AS partition FROM dbo.t"},
		// Context keywords as table/column names in DDL
		{"CREATE TABLE encryption", "CREATE TABLE dbo.encryption_test (a INT)"},
		{"column named window", "CREATE TABLE dbo.t_window (window INT)"},
		{"column named partition", "CREATE TABLE dbo.t_partition (partition INT)"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ssAccepts, err := oracle.canParse(tc.sql)
			if err != nil {
				t.Fatalf("oracle error: %v", err)
			}

			_, omniErr := Parse(tc.sql)
			omniAccepts := omniErr == nil

			if ssAccepts != omniAccepts {
				t.Errorf("MISMATCH: SQL Server %s, omni %s\n  SQL: %s",
					boolToAcceptReject(ssAccepts), boolToAcceptReject(omniAccepts), tc.sql)
			}
		})
	}
}

// TestKeywordOracleCoreAsIdentifier verifies that core keywords are correctly
// rejected as unquoted identifiers and accepted when bracket-quoted.
func TestKeywordOracleCoreAsIdentifier(t *testing.T) {
	oracle := startParserOracle(t)

	coreKeywords := []string{
		"select", "from", "where", "insert", "update", "delete",
		"create", "drop", "alter", "order", "group", "having",
		"join", "on", "set", "into", "values", "table",
	}

	for _, kw := range coreKeywords {
		// Unquoted column name — should fail
		t.Run(fmt.Sprintf("column_%s_unquoted", kw), func(t *testing.T) {
			sql := fmt.Sprintf("CREATE TABLE dbo.t_core_%s (%s INT)", kw, kw)
			ssAccepts, err := oracle.canParse(sql)
			if err != nil {
				t.Fatalf("oracle error: %v", err)
			}

			_, omniErr := Parse(sql)
			omniAccepts := omniErr == nil

			if ssAccepts != omniAccepts {
				t.Errorf("MISMATCH: SQL Server %s, omni %s\n  SQL: %s",
					boolToAcceptReject(ssAccepts), boolToAcceptReject(omniAccepts), sql)
			}
		})

		// Bracket-quoted column name — should succeed
		t.Run(fmt.Sprintf("column_%s_bracketed", kw), func(t *testing.T) {
			sql := fmt.Sprintf("CREATE TABLE dbo.t_core_%s_q ([%s] INT)", kw, kw)
			ssAccepts, err := oracle.canParse(sql)
			if err != nil {
				t.Fatalf("oracle error: %v", err)
			}

			_, omniErr := Parse(sql)
			omniAccepts := omniErr == nil

			if ssAccepts != omniAccepts {
				t.Errorf("MISMATCH: SQL Server %s, omni %s\n  SQL: %s",
					boolToAcceptReject(ssAccepts), boolToAcceptReject(omniAccepts), sql)
			}
		})

		// Bare alias — should fail
		t.Run(fmt.Sprintf("alias_%s_bare", kw), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT 1 %s FROM dbo.t", kw)
			ssAccepts, err := oracle.canParse(sql)
			if err != nil {
				t.Fatalf("oracle error: %v", err)
			}

			_, omniErr := Parse(sql)
			omniAccepts := omniErr == nil

			if ssAccepts != omniAccepts {
				t.Errorf("MISMATCH: SQL Server %s, omni %s\n  SQL: %s",
					boolToAcceptReject(ssAccepts), boolToAcceptReject(omniAccepts), sql)
			}
		})
	}
}

func boolToAcceptReject(b bool) string {
	if b {
		return "ACCEPTS"
	}
	return "REJECTS"
}
