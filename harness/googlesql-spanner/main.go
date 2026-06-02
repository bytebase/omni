// Command googlesql-spanner-harness is the differential ORACLE side for the
// omni `googlesql` parser migration. It runs SQL through a live Cloud Spanner
// emulator (gRPC) and emits a compact accept/reject verdict per statement, so
// that grammar-node tests can diff omni's parser accept/reject against a real
// GoogleSQL implementation.
//
// It is the GoogleSQL analogue of harness/mssql-scriptdom: an external
// reference parser/engine that omni diffs against. Unlike ScriptDOM (which
// exposes an AST), the Spanner emulator only answers "did it parse?" — so this
// harness emits a verdict + error classification, NOT an AST shape.
//
// WHY the emulator and how errors are classified: see
//
//	docs/migration/googlesql/oracle.md
//
// In short: the gRPC error CODE is not a discriminator (syntax and semantic
// failures are both InvalidArgument); the message PREFIX is:
//
//	"Syntax error:"                        -> grammar REJECT (queries / DML)
//	"Error parsing Spanner DDL statement:" -> grammar REJECT (DDL)
//	anything else (Table not found, Unrecognized name, "X is not supported",
//	FailedPrecondition, Internal)          -> grammar ACCEPT (semantic)
//
// CAVEAT: the emulator speaks Spanner's GoogleSQL, a SUBSET of the BigQuery +
// Spanner union the omni parser must accept. A reject of a BigQuery-only form
// (scripting, EXPORT DATA, CREATE MODEL, PIVOT/UNPIVOT, GQL, ...) is NOT
// authoritative — callers must tag such forms and triangulate against the
// docs corpus + the legacy .g4. This harness reports only the Spanner verdict;
// authoritativeness is the caller's call.
//
// Protocol (matches mssql-scriptdom):
//   - Batch line mode (default, or GOOGLESQL_HARNESS_LINE=1): read one
//     base64-encoded SQL per line on stdin; emit one JSON verdict per line.
//     Base64 avoids newline-in-SQL ambiguity and lets callers reuse one
//     process across many fixtures. Output line N corresponds to input line N.
//   - Single mode (GOOGLESQL_HARNESS_LINE=0): read all of stdin as one
//     statement; emit one JSON verdict.
//
// Verdict JSON:
//
//	{"verdict":"accept|reject","kind":"query|dml|ddl","reason":"ok|syntax|semantic",
//	 "code":"InvalidArgument","message":"..."}
//
// Run with the emulator up (see oracle.md):
//
//	SPANNER_EMULATOR_HOST=localhost:9010 go run .
package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	project    = "test-project"
	instanceID = "test-instance"
	dbID       = "googlesql_harness"
	parent     = "projects/" + project + "/instances/" + instanceID
	dbPath     = parent + "/databases/" + dbID
	configName = "projects/" + project + "/instanceConfigs/emulator-config"
)

// errAbort rolls back a DML validation transaction once the statement has
// parsed, so the harness never mutates the scratch database.
var errAbort = errors.New("harness: abort after parse")

func main() {
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
		fmt.Fprintln(os.Stderr, "googlesql-spanner-harness: SPANNER_EMULATOR_HOST is unset; point it at the emulator gRPC port (e.g. localhost:9010). See docs/migration/googlesql/oracle.md")
		os.Exit(2)
	}

	ctx := context.Background()
	o, err := newOracle(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "googlesql-spanner-harness: bootstrap failed: %v\n", err)
		os.Exit(1)
	}
	defer o.close()

	enc := json.NewEncoder(os.Stdout)

	if os.Getenv("GOOGLESQL_HARNESS_LINE") == "0" {
		// Single mode: all of stdin is one statement.
		var sb strings.Builder
		sc := bufio.NewScanner(os.Stdin)
		sc.Buffer(make([]byte, 0, 1<<20), 1<<24)
		for sc.Scan() {
			sb.WriteString(sc.Text())
			sb.WriteByte('\n')
		}
		_ = enc.Encode(o.evaluate(ctx, strings.TrimSpace(sb.String())))
		return
	}

	// Batch line mode: one base64 SQL per line.
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 0, 1<<20), 1<<24)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		raw, err := base64.StdEncoding.DecodeString(line)
		if err != nil {
			_ = enc.Encode(verdict{Verdict: "error", Reason: "bad base64", Message: err.Error()})
			continue
		}
		_ = enc.Encode(o.evaluate(ctx, string(raw)))
	}
}

type verdict struct {
	Verdict string `json:"verdict"`           // accept | reject | error
	Kind    string `json:"kind,omitempty"`    // query | dml | ddl
	Reason  string `json:"reason,omitempty"`  // ok | syntax | semantic
	Code    string `json:"code,omitempty"`    // gRPC status code
	Message string `json:"message,omitempty"` // truncated, single-line
}

type oracle struct {
	cli     *spanner.Client
	dbAdmin *database.DatabaseAdminClient
	instCli *instance.InstanceAdminClient
}

func newOracle(ctx context.Context) (*oracle, error) {
	instCli, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("instance admin client: %w", err)
	}
	// Create the instance; ignore AlreadyExists.
	if op, err := instCli.CreateInstance(ctx, &instancepb.CreateInstanceRequest{
		Parent:     "projects/" + project,
		InstanceId: instanceID,
		Instance:   &instancepb.Instance{Config: configName, DisplayName: "omni-harness", NodeCount: 1},
	}); err == nil {
		if _, err := op.Wait(ctx); err != nil && status.Code(err) != codes.AlreadyExists {
			return nil, fmt.Errorf("create instance: %w", err)
		}
	} else if status.Code(err) != codes.AlreadyExists {
		return nil, fmt.Errorf("create instance: %w", err)
	}

	dbAdmin, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("database admin client: %w", err)
	}
	// Fresh, empty scratch database each run for deterministic verdicts.
	_ = dbAdmin.DropDatabase(ctx, &databasepb.DropDatabaseRequest{Database: dbPath})
	op, err := dbAdmin.CreateDatabase(ctx, &databasepb.CreateDatabaseRequest{
		Parent:          parent,
		CreateStatement: "CREATE DATABASE `" + dbID + "`",
	})
	if err != nil {
		return nil, fmt.Errorf("create database: %w", err)
	}
	if _, err := op.Wait(ctx); err != nil {
		return nil, fmt.Errorf("create database wait: %w", err)
	}

	cli, err := spanner.NewClient(ctx, dbPath)
	if err != nil {
		return nil, fmt.Errorf("data client: %w", err)
	}
	return &oracle{cli: cli, dbAdmin: dbAdmin, instCli: instCli}, nil
}

func (o *oracle) close() {
	if o.cli != nil {
		o.cli.Close()
	}
	if o.dbAdmin != nil {
		_ = o.dbAdmin.Close()
	}
	if o.instCli != nil {
		_ = o.instCli.Close()
	}
}

// evaluate routes a statement to the right Spanner API by its leading keyword
// and classifies the result.
func (o *oracle) evaluate(ctx context.Context, sql string) verdict {
	if strings.TrimSpace(sql) == "" {
		return verdict{Verdict: "error", Reason: "empty"}
	}
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	kind := kindOf(sql)
	var err error
	switch kind {
	case "ddl":
		err = o.validateDDL(ctx, sql)
	case "dml":
		err = o.validateDML(ctx, sql)
	default: // query
		err = o.validateQuery(ctx, sql)
	}
	v := classify(err)
	v.Kind = kind
	return v
}

// validateQuery executes the query and reads only the first row. SELECT/WITH/
// VALUES have no side effects, so execution is safe and — unlike the
// emulator's PLAN-mode AnalyzeQuery, which returns "query plan unavailable" —
// it yields a clean accept for valid queries. A syntax error surfaces before
// any row; iterator.Done (zero rows) is an accept.
func (o *oracle) validateQuery(ctx context.Context, sql string) error {
	iter := o.cli.Single().Query(ctx, spanner.Statement{SQL: sql})
	defer iter.Stop()
	_, err := iter.Next()
	if errors.Is(err, iterator.Done) {
		return nil
	}
	return err
}

// validateDML runs the statement inside a read-write transaction and aborts
// after it parses, so no rows are mutated. Update() surfaces parse/semantic
// errors before the abort.
func (o *oracle) validateDML(ctx context.Context, sql string) error {
	_, err := o.cli.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		if _, e := txn.Update(ctx, spanner.Statement{SQL: sql}); e != nil {
			return e
		}
		return errAbort
	})
	if errors.Is(err, errAbort) {
		return nil
	}
	return err
}

// validateDDL submits to UpdateDatabaseDdl against the scratch database. DDL
// does mutate schema, but only the syntax verdict (message prefix) matters;
// semantic outcomes (Duplicate name, missing parent, ...) classify as accept.
func (o *oracle) validateDDL(ctx context.Context, sql string) error {
	op, err := o.dbAdmin.UpdateDatabaseDdl(ctx, &databasepb.UpdateDatabaseDdlRequest{
		Database:   dbPath,
		Statements: []string{sql},
	})
	if err != nil {
		return err
	}
	return op.Wait(ctx)
}

// classify maps a Spanner error to a grammar verdict. See oracle.md.
func classify(err error) verdict {
	if err == nil {
		return verdict{Verdict: "accept", Reason: "ok", Code: codes.OK.String()}
	}
	st, _ := status.FromError(err)
	msg := st.Message()
	v := verdict{Code: st.Code().String(), Message: oneLine(msg)}
	switch {
	case strings.HasPrefix(msg, "Syntax error:"),
		strings.HasPrefix(msg, "Error parsing Spanner DDL statement"):
		v.Verdict, v.Reason = "reject", "syntax"
	default:
		// Table not found / Unrecognized name / "X is not supported" /
		// FailedPrecondition / Internal => the grammar accepted the input.
		v.Verdict, v.Reason = "accept", "semantic"
	}
	return v
}

func oneLine(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	if len(s) > 200 {
		s = s[:200] + "…"
	}
	return s
}

// kindOf routes by the first significant keyword, skipping leading whitespace,
// line (-- and #) and block (/* */) comments, and statement hints (@{...}).
func kindOf(sql string) string {
	s := stripLeading(sql)
	word := leadingWord(s)
	switch word {
	case "INSERT", "UPDATE", "DELETE", "MERGE":
		return "dml"
	case "CREATE", "ALTER", "DROP", "GRANT", "REVOKE", "RENAME", "TRUNCATE", "ANALYZE":
		return "ddl"
	default:
		// SELECT, WITH, VALUES, FROM, TABLE, "(", and everything else
		// (incl. BigQuery scripting/other) route through the query path,
		// which still yields a Spanner verdict.
		return "query"
	}
}

func stripLeading(s string) string {
	for {
		s = strings.TrimLeft(s, " \t\r\n")
		switch {
		case strings.HasPrefix(s, "--"), strings.HasPrefix(s, "#"):
			if i := strings.IndexByte(s, '\n'); i >= 0 {
				s = s[i+1:]
				continue
			}
			return ""
		case strings.HasPrefix(s, "/*"):
			if i := strings.Index(s, "*/"); i >= 0 {
				s = s[i+2:]
				continue
			}
			return ""
		case strings.HasPrefix(s, "@{"):
			if i := strings.IndexByte(s, '}'); i >= 0 {
				s = s[i+1:]
				continue
			}
			return ""
		default:
			return s
		}
	}
}

func leadingWord(s string) string {
	i := 0
	for i < len(s) && (isAlpha(s[i]) || s[i] == '_') {
		i++
	}
	return strings.ToUpper(s[:i])
}

func isAlpha(b byte) bool { return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') }
