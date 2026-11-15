package driver

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/pg/internal/testcontainer"
)

func TestGetDatabaseInCreateDatabaseStatement(t *testing.T) {
	cases := []struct {
		stmt    string
		want    string
		wantErr bool
	}{
		{`CREATE DATABASE "hello" ENCODING "UTF8";`, "hello", false},
		{`CREATE DATABASE "hello";`, "hello", false},
		{`CREATE DATABASE hello;`, "hello", false},
		{`CREATE DATABASE hello ENCODING "UTF8";`, "hello", false},
		{`CREATE DATABASE;`, "", true},
	}
	for _, tc := range cases {
		got, err := getDatabaseInCreateDatabaseStatement(tc.stmt)
		if (err != nil) != tc.wantErr {
			t.Errorf("%q: err = %v, wantErr %v", tc.stmt, err, tc.wantErr)
		}
		if got != tc.want {
			t.Errorf("%q: got %q, want %q", tc.stmt, got, tc.want)
		}
	}
}

type stepLogger struct {
	steps []string
	cmds  []driver.Command
}

func (l *stepLogger) CommandExecute(cmd driver.Command) {
	l.cmds = append(l.cmds, cmd)
	l.steps = append(l.steps, fmt.Sprintf("exec[%d:%d]", cmd.ByteStart, cmd.ByteEnd))
}

func (l *stepLogger) CommandResponse(affected int64, perResult []int64, err error) {
	if err != nil {
		l.steps = append(l.steps, fmt.Sprintf("resp:err:%v", perResult == nil))
		return
	}
	l.steps = append(l.steps, fmt.Sprintf("resp:%d:%v", affected, perResult))
}

func (l *stepLogger) TransactionControl(op driver.TransactionOp, err error) {
	if err != nil {
		l.steps = append(l.steps, fmt.Sprintf("tx:%s:err", op))
		return
	}
	l.steps = append(l.steps, fmt.Sprintf("tx:%s", op))
}

func (l *stepLogger) Retry(_ error, attempt, max int) {
	l.steps = append(l.steps, fmt.Sprintf("retry:%d/%d", attempt, max))
}

// TestExecute covers transactional execution, non-transactional statement
// deferral (CREATE/DROP INDEX CONCURRENTLY, VACUUM), rollback on failure,
// auto-commit mode, PL/pgSQL DO blocks, and tenant mode.
func TestExecute(t *testing.T) {
	ctx := context.Background()
	for _, role := range []string{"pg12", testcontainer.Default, testcontainer.PG18} {
		t.Run(role, func(t *testing.T) {
			d := openQueryDriver(t, role)

			// 1. Transactional script with non-transactional statements deferred after COMMIT.
			const txScript = "CREATE TABLE exec_t (id int PRIMARY KEY, v text);\nINSERT INTO exec_t VALUES (1, 'a'), (2, 'b'), (3, 'c');\nCREATE INDEX CONCURRENTLY idx_exec_v ON exec_t (v);\nUPDATE exec_t SET v = 'x' WHERE id <= 2;\nVACUUM exec_t;\nDROP INDEX CONCURRENTLY idx_exec_v;"
			log := &stepLogger{}
			affected, err := d.Execute(ctx, txScript, driver.ExecuteOptions{
				Transaction: driver.Transactional,
				Log:         log,
			})
			if err != nil {
				t.Fatal(err)
			}
			if affected != 5 {
				t.Fatalf("affected = %d, want 5", affected)
			}
			for _, cmd := range log.cmds {
				if txScript[cmd.ByteStart:cmd.ByteEnd] != cmd.Text {
					t.Errorf("command span [%d:%d] = %q, cmd.Text = %q", cmd.ByteStart, cmd.ByteEnd, txScript[cmd.ByteStart:cmd.ByteEnd], cmd.Text)
				}
			}
			wantSteps := []string{
				"tx:BEGIN",
				"exec[0:49]", "resp:0:[0]",
				"exec[49:105]", "resp:3:[3]",
				"exec[157:198]", "resp:2:[2]",
				"tx:COMMIT",
				"exec[105:157]", "resp:0:[0]",
				"exec[198:213]", "resp:0:[0]",
				"exec[213:249]", "resp:0:[0]",
			}
			if !slices.Equal(log.steps, wantSteps) {
				t.Fatalf("tx steps =\n%v\nwant\n%v", log.steps, wantSteps)
			}

			// 2. Transactional rollback on error leaves prior rows untouched.
			rbLog := &stepLogger{}
			if _, err := d.Execute(ctx, "INSERT INTO exec_t VALUES (10, 'rolled_back');\nSELECT 1 / 0;", driver.ExecuteOptions{
				Transaction: driver.Transactional,
				Log:         rbLog,
			}); err == nil {
				t.Fatal("expected division by zero error")
			}
			if !slices.Contains(rbLog.steps, "tx:ROLLBACK") {
				t.Fatalf("rollback not logged: %v", rbLog.steps)
			}
			var count int
			if err := d.db.QueryRowContext(ctx, "SELECT count(*) FROM exec_t WHERE id = 10").Scan(&count); err != nil || count != 0 {
				t.Fatalf("rolled-back row persisted: count=%d err=%v", count, err)
			}

			// 3. Auto-commit mode commits earlier statements before a failure.
			acLog := &stepLogger{}
			acAffected, err := d.Execute(ctx, "INSERT INTO exec_t VALUES (20, 'kept');\nSELECT 1 / 0;", driver.ExecuteOptions{
				Transaction: driver.AutoCommit,
				Log:         acLog,
			})
			if err == nil {
				t.Fatal("expected division by zero error in auto-commit")
			}
			if acAffected != 1 {
				t.Fatalf("auto-commit partial affected = %d, want 1", acAffected)
			}
			for _, st := range acLog.steps {
				if strings.HasPrefix(st, "tx:") {
					t.Errorf("auto-commit logged transaction step %q", st)
				}
			}
			if err := d.db.QueryRowContext(ctx, "SELECT count(*) FROM exec_t WHERE id = 20").Scan(&count); err != nil || count != 1 {
				t.Fatalf("auto-commit row missing: count=%d err=%v", count, err)
			}

			// 4. Single PL/pgSQL DO block runs outside an outer BEGIN/COMMIT.
			doLog := &stepLogger{}
			const doScript = "DO $$ BEGIN INSERT INTO exec_t VALUES (30, 'do'); COMMIT; INSERT INTO exec_t VALUES (31, 'do2'); END $$;"
			if _, err := d.Execute(ctx, doScript, driver.ExecuteOptions{
				Transaction: driver.Transactional,
				Log:         doLog,
			}); err != nil {
				t.Fatalf("DO block failed: %v", err)
			}
			if !slices.Equal(doLog.steps, []string{fmt.Sprintf("exec[0:%d]", len(doScript)), "resp:0:[0]"}) {
				t.Fatalf("DO steps = %v", doLog.steps)
			}

			// 5. Tenant mode creates objects as the database owner and elevates for superuser statements.
			mustExec(t, d.db, `
				CREATE ROLE exec_tenant_owner NOLOGIN;
				CREATE ROLE exec_reader NOLOGIN;
				GRANT CREATE ON SCHEMA public TO exec_tenant_owner;
			`)
			mustExec(t, d.db, fmt.Sprintf("ALTER DATABASE %s OWNER TO exec_tenant_owner;", quoteIdent(d.databaseName)))
			tenantDriver := &Driver{db: d.db, databaseName: d.databaseName, tenantMode: true}
			if _, err := tenantDriver.Execute(ctx, "CREATE TABLE tenant_owned (id int);\nGRANT SELECT ON tenant_owned TO exec_reader;", driver.ExecuteOptions{
				Transaction: driver.Transactional,
			}); err != nil {
				t.Fatalf("tenant Execute failed: %v", err)
			}
			var tableOwner string
			if err := d.db.QueryRowContext(ctx, "SELECT tableowner FROM pg_tables WHERE schemaname = 'public' AND tablename = 'tenant_owned'").Scan(&tableOwner); err != nil || tableOwner != "exec_tenant_owner" {
				t.Fatalf("tenant tableowner = %q, err = %v", tableOwner, err)
			}
		})
	}
}

func TestExecuteCreateDatabase(t *testing.T) {
	ctx := context.Background()
	for _, role := range []string{"pg12", testcontainer.Default} {
		t.Run(role, func(t *testing.T) {
			d := openQueryDriver(t, role)
			mustExec(t, d.db, "CREATE ROLE created_db_owner NOLOGIN;")

			stmt := "CREATE DATABASE \"omni_new_db\" ENCODING \"UTF8\";\nALTER DATABASE \"omni_new_db\" OWNER TO \"created_db_owner\";"
			if _, err := d.Execute(ctx, stmt, driver.ExecuteOptions{CreateDatabase: true}); err != nil {
				t.Fatal(err)
			}
			var owner string
			if err := d.db.QueryRowContext(ctx, "SELECT pg_catalog.pg_get_userbyid(datdba) FROM pg_database WHERE datname = 'omni_new_db'").Scan(&owner); err != nil || owner != "created_db_owner" {
				t.Fatalf("created database owner = %q, err = %v", owner, err)
			}

			// Re-running CreateDatabase for an existing database is a no-op.
			if _, err := d.Execute(ctx, stmt, driver.ExecuteOptions{CreateDatabase: true}); err != nil {
				t.Fatalf("idempotent CreateDatabase failed: %v", err)
			}
		})
	}
}
