package testcontainer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"cloud.google.com/go/spanner"
	spannerdb "cloud.google.com/go/spanner/admin/database/apiv1"
	"cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bytebase/omni/enginetest"
)

const pgJudgeInstance = "pg-compat"

var pgJudgeSeq atomic.Uint64

// errPGAbort rolls back a DML validation transaction once the statement has
// parsed, so the PostgreSQL-dialect judge never mutates its scratch database.
var errPGAbort = errors.New("pg-judge: abort after parse")

// PostgresOracle evaluates SQL statements directly against a
// DatabaseDialect_POSTGRESQL database on the pinned Cloud Spanner emulator.
//
// Queries run in single-use read-only snapshots, DML runs inside a read-write
// transaction that always aborts before commit, and DDL runs against a
// dedicated empty PostgreSQL-dialect database that is dropped and recreated
// whenever a DDL statement succeeds, so no statement's side effects can alter a
// later statement's verdict.
type PostgresOracle struct {
	cli      *spanner.Client
	dbAdmin  *spannerdb.DatabaseAdminClient
	queryDB  string
	ddlDB    string
	ddlDirty bool
}

// JudgePostgres returns a PostgresOracle connected to a fresh empty
// DatabaseDialect_POSTGRESQL database on the pinned Spanner emulator, outside
// any tape. It runs only in a recording (OMNI_RECORD=1 or OMNI_RECORD=verify).
func JudgePostgres(t testing.TB) *PostgresOracle {
	t.Helper()
	if !enginetest.Recording(t) {
		t.Fatalf("testcontainer.JudgePostgres outside a recording: a replay reaches no server")
	}
	host := Host(t)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if err := ensureInstance(ctx, host, pgJudgeInstance); err != nil {
		t.Fatalf("create Spanner emulator instance for PostgreSQL judge: %v", err)
	}
	dbAdmin, err := databaseAdmin(ctx, host)
	if err != nil {
		t.Fatalf("create Spanner database admin client: %v", err)
	}
	queryDB, err := createEmptyPGDatabase(ctx, dbAdmin, "pgq")
	if err != nil {
		_ = dbAdmin.Close()
		t.Fatalf("create empty PostgreSQL-dialect query database: %v", err)
	}
	cli, err := spanner.NewClient(ctx, queryDB, emulatorOptions(host)...)
	if err != nil {
		_ = dbAdmin.DropDatabase(ctx, &databasepb.DropDatabaseRequest{Database: queryDB})
		_ = dbAdmin.Close()
		t.Fatalf("create Spanner data client for PostgreSQL judge: %v", err)
	}
	o := &PostgresOracle{
		cli:     cli,
		dbAdmin: dbAdmin,
		queryDB: queryDB,
	}
	t.Cleanup(o.close)
	return o
}

func (o *PostgresOracle) close() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if o.cli != nil {
		o.cli.Close()
	}
	if o.dbAdmin != nil {
		if o.queryDB != "" {
			_ = o.dbAdmin.DropDatabase(ctx, &databasepb.DropDatabaseRequest{Database: o.queryDB})
		}
		if o.ddlDB != "" {
			_ = o.dbAdmin.DropDatabase(ctx, &databasepb.DropDatabaseRequest{Database: o.ddlDB})
		}
		_ = o.dbAdmin.Close()
	}
}

func createEmptyPGDatabase(ctx context.Context, dbAdmin *spannerdb.DatabaseAdminClient, prefix string) (string, error) {
	id := fmt.Sprintf("%s_%d_%d", prefix, os.Getpid()%100000, pgJudgeSeq.Add(1))
	op, err := dbAdmin.CreateDatabase(ctx, &databasepb.CreateDatabaseRequest{
		Parent:          InstancePath(pgJudgeInstance),
		CreateStatement: fmt.Sprintf(`CREATE DATABASE "%s"`, id),
		DatabaseDialect: databasepb.DatabaseDialect_POSTGRESQL,
	})
	if err != nil {
		return "", err
	}
	if _, err := op.Wait(ctx); err != nil {
		return "", err
	}
	return DatabasePath(pgJudgeInstance, id), nil
}

// Eval sends sql to the appropriate Spanner API (query, dml, or ddl) on the
// PostgreSQL-dialect emulator database and returns the routed kind and the raw
// emulator error (or nil on success).
func (o *PostgresOracle) Eval(ctx context.Context, sql string) (kind string, err error) {
	kind = pgKindOf(sql)
	switch kind {
	case "ddl":
		return "ddl", o.validateDDL(ctx, sql)
	case "dml":
		return "dml", o.validateDML(ctx, sql)
	default:
		err = o.validateQuery(ctx, sql)
		if isSingleUseDMLError(err) {
			return "dml", o.validateDML(ctx, sql)
		}
		return "query", err
	}
}

func isSingleUseDMLError(err error) bool {
	if err == nil {
		return false
	}
	st, ok := status.FromError(err)
	return ok && st.Code() == codes.InvalidArgument &&
		strings.Contains(st.Message(), "DML statements may not be performed in single-use transactions")
}

func (o *PostgresOracle) validateQuery(ctx context.Context, sql string) error {
	iter := o.cli.Single().Query(ctx, spanner.Statement{SQL: sql})
	defer iter.Stop()
	_, err := iter.Next()
	if errors.Is(err, iterator.Done) {
		return nil
	}
	return err
}

func (o *PostgresOracle) validateDML(ctx context.Context, sql string) error {
	_, err := o.cli.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		if _, e := txn.Update(ctx, spanner.Statement{SQL: sql}); e != nil {
			return e
		}
		return errPGAbort
	})
	if errors.Is(err, errPGAbort) {
		return nil
	}
	return err
}

func (o *PostgresOracle) validateDDL(ctx context.Context, sql string) error {
	if o.ddlDB == "" || o.ddlDirty {
		if err := o.resetDDLDatabase(ctx); err != nil {
			return fmt.Errorf("reset PostgreSQL DDL database: %w", err)
		}
	}
	op, err := o.dbAdmin.UpdateDatabaseDdl(ctx, &databasepb.UpdateDatabaseDdlRequest{
		Database:   o.ddlDB,
		Statements: []string{sql},
	})
	if err != nil {
		return err
	}
	err = op.Wait(ctx)
	if err == nil {
		o.ddlDirty = true
	}
	return err
}

func (o *PostgresOracle) resetDDLDatabase(ctx context.Context) error {
	if o.ddlDB != "" {
		if err := o.dbAdmin.DropDatabase(ctx, &databasepb.DropDatabaseRequest{Database: o.ddlDB}); err != nil {
			return err
		}
		o.ddlDB = ""
	}
	dbPath, err := createEmptyPGDatabase(ctx, o.dbAdmin, "pgd")
	if err != nil {
		return err
	}
	o.ddlDB = dbPath
	o.ddlDirty = false
	return nil
}

// pgKindOf routes a PostgreSQL statement to the Spanner API that handles it:
// DML, DDL, or query. CTE-led DML (WITH ... INSERT/UPDATE/DELETE) starts in
// the query path and falls back to DML automatically when the emulator reports
// that DML cannot run in a single-use transaction.
func pgKindOf(sql string) string {
	s := stripPGLeading(sql)
	switch pgLeadingWord(s) {
	case "INSERT", "UPDATE", "DELETE", "MERGE":
		return "dml"
	case "CREATE", "ALTER", "DROP", "GRANT", "REVOKE", "RENAME", "TRUNCATE", "ANALYZE":
		return "ddl"
	default:
		return "query"
	}
}

// stripPGLeading skips leading whitespace, -- line comments, and (possibly
// nested) /* ... */ block comments.
func stripPGLeading(s string) string {
	for {
		s = strings.TrimLeft(s, " \t\r\n")
		switch {
		case strings.HasPrefix(s, "--"):
			if i := strings.IndexByte(s, '\n'); i >= 0 {
				s = s[i+1:]
				continue
			}
			return ""
		case strings.HasPrefix(s, "/*"):
			depth := 1
			i := 2
			for i < len(s) && depth > 0 {
				if i+1 < len(s) && s[i] == '/' && s[i+1] == '*' {
					depth++
					i += 2
				} else if i+1 < len(s) && s[i] == '*' && s[i+1] == '/' {
					depth--
					i += 2
				} else {
					i++
				}
			}
			if depth > 0 {
				return s // unclosed comment: let the emulator reject it
			}
			s = s[i:]
		default:
			return s
		}
	}
}

func pgLeadingWord(s string) string {
	i := 0
	for i < len(s) && ((s[i] >= 'a' && s[i] <= 'z') || (s[i] >= 'A' && s[i] <= 'Z') || s[i] == '_') {
		i++
	}
	return strings.ToUpper(s[:i])
}
