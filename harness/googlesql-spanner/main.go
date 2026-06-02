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
// In short, classification is FAIL-CLOSED and kind-aware (see classify); parse
// rejects key on the error-message PREFIX (InvalidArgument is overloaded by
// semantic failures in both kinds):
//
//	query/DML: InvalidArgument + "Syntax error:" -> REJECT; other InvalidArgument
//	           (Table not found, Unrecognized name, "X is not supported") and a
//	           runtime OutOfRange -> ACCEPT.
//	DDL:       InvalidArgument + "Error parsing Spanner DDL statement:" -> REJECT;
//	           other InvalidArgument (bad index col, type change, generated/check/
//	           default expr), FailedPrecondition/NotFound, and the Internal
//	           GOOGLESQL_RET_CHECK quirk -> ACCEPT (semantic).
//	anything else (Unavailable, DeadlineExceeded, Canceled, Aborted, a resource-
//	level miss, a generic Internal, a non-gRPC error) -> verdict "error": the
//	oracle could not decide. It NEVER returns ACCEPT for an infra failure.
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
	"io"
	"os"
	"strings"
	"time"
	"unicode/utf8"

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
	emit := func(v verdict) {
		if err := enc.Encode(v); err != nil {
			// stdout is broken; nothing useful left to do — fail loudly rather
			// than silently produce a truncated verdict stream.
			fmt.Fprintf(os.Stderr, "googlesql-spanner-harness: stdout write failed: %v\n", err)
			os.Exit(1)
		}
	}

	if os.Getenv("GOOGLESQL_HARNESS_LINE") == "0" {
		// Single mode: all of stdin is one statement. Read raw (no line cap).
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "googlesql-spanner-harness: stdin read failed: %v\n", err)
			os.Exit(1)
		}
		emit(o.evaluate(ctx, strings.TrimSpace(string(data))))
		return
	}

	// Batch line mode: one base64 SQL per line, one verdict line out per input
	// line. The output-line-N <-> input-line-N contract is load-bearing for
	// callers that zip verdicts to fixtures, so EVERY input line (incl. blank
	// and undecodable) emits exactly one verdict, and a scanner error fails the
	// process rather than silently dropping the tail.
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 0, 1<<20), 64<<20) // 64 MiB cap (base64 inflates 4/3)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			emit(verdict{Verdict: "error", Reason: "blank line"})
			continue
		}
		raw, err := base64.StdEncoding.DecodeString(line)
		if err != nil {
			emit(verdict{Verdict: "error", Reason: "bad base64", Message: err.Error()})
			continue
		}
		emit(o.evaluate(ctx, string(raw)))
	}
	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "googlesql-spanner-harness: stdin read failed (line over 64MiB cap?): %v\n", err)
		os.Exit(1)
	}
}

type verdict struct {
	Verdict string `json:"verdict"`           // accept | reject | error
	Kind    string `json:"kind,omitempty"`    // query | dml | ddl
	Reason  string `json:"reason,omitempty"`  // accept/reject: ok|syntax|semantic; error: infra|empty|"blank line"|"bad base64"
	Code    string `json:"code,omitempty"`    // gRPC status code (or "non-status")
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
	// Blank, or comment/hint-only (stripLeading reduces it to nothing): there is
	// no statement to give a grammar verdict on. Return "error", not a verdict,
	// rather than shipping a bare comment to the emulator (which would reject it).
	if strings.TrimSpace(stripLeading(sql)) == "" {
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
	v := classify(kind, err)
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
// does mutate schema and state accumulates across a batch, but parse happens
// before schema validation, so the accept/reject syntax VERDICT is order-stable
// (only the semantic `reason` — ok vs Duplicate-name/NotFound — can vary with
// prior DDL in the same run). Verdict correctness, the thing callers diff on,
// is unaffected.
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

// classify maps a Spanner error to a grammar verdict for a statement of the
// given kind ("query" | "dml" | "ddl"). See oracle.md.
//
// FAIL-CLOSED: only outcomes that are genuinely a parser/semantic result of the
// emulator yield a grammar verdict (accept/reject). Anything that means "the
// oracle could not decide" — transport failures, timeouts, cancellation,
// transaction-abort retry exhaustion, a resource-level (database/instance/
// session) miss, a generic Internal, a non-gRPC error — returns verdict
// "error". An oracle must NEVER let an infra failure masquerade as accept: a
// silent false-accept would mask real parser bugs (or fabricate divergences)
// in every grammar node that trusts this verdict.
//
// The discriminator is kind-aware and keys on the gRPC CODE wherever possible,
// not on exact message text (the emulator's strings float with its image):
//   - DDL: the emulator returns a PARSE failure as InvalidArgument and a
//     SEMANTIC failure as FailedPrecondition / NotFound / Internal, so the code
//     alone decides — robust to message drift.
//   - query/DML: parse failures are InvalidArgument with a "Syntax error:"
//     message; every OTHER InvalidArgument is a semantic result (Table not
//     found, Unrecognized name, "X is not supported") => grammar parsed.
func classify(kind string, err error) verdict {
	if err == nil {
		return verdict{Verdict: "accept", Reason: "ok", Code: codes.OK.String()}
	}
	st, ok := status.FromError(err)
	if !ok {
		// Not a gRPC status (dial failure, context.DeadlineExceeded/Canceled,
		// wrapped I/O error, ...). The oracle cannot decide.
		return verdict{Verdict: "error", Reason: "infra", Code: "non-status", Message: oneLine(err.Error())}
	}
	code := st.Code()
	msg := st.Message()
	v := verdict{Code: code.String(), Message: oneLine(msg)}

	// Transport / availability / lifecycle codes are never a grammar verdict,
	// regardless of kind.
	switch code {
	case codes.Unavailable, codes.DeadlineExceeded, codes.Canceled,
		codes.Aborted, codes.ResourceExhausted, codes.Unknown,
		codes.Unauthenticated, codes.PermissionDenied:
		v.Verdict, v.Reason = "error", "infra"
		return v
	}
	// A resource-level NotFound/AlreadyExists (the scratch database/instance/
	// session is gone — e.g. after an emulator restart) is infra, NOT a
	// per-statement semantic result. Match the specific lifecycle phrases so a
	// table/column literally named "Session" (msg "Table not found: Session")
	// is NOT misrouted here.
	if (code == codes.NotFound || code == codes.AlreadyExists) && isResourceLevel(msg) {
		v.Verdict, v.Reason = "error", "infra"
		return v
	}

	if kind == "ddl" {
		// DDL PARSE failures (InvalidArgument) carry the exact prefix below;
		// DDL SEMANTIC failures are ALSO InvalidArgument but do NOT (verified:
		// bad index column, column type change, bad length, generated/CHECK/
		// DEFAULT expression errors — the latter start "Error parsing " but
		// diverge before "...statement:"). So the prefix, not the code, is the
		// discriminator. (Drift risk is mitigated by pinning the emulator
		// digest + a required-green live TestOracleLive; see oracle.md.)
		switch code {
		case codes.InvalidArgument:
			if strings.HasPrefix(msg, "Error parsing Spanner DDL statement:") {
				v.Verdict, v.Reason = "reject", "syntax"
			} else {
				v.Verdict, v.Reason = "accept", "semantic" // semantic InvalidArgument => grammar parsed
			}
		case codes.FailedPrecondition, codes.NotFound: // Duplicate name / missing INTERLEAVE parent
			v.Verdict, v.Reason = "accept", "semantic"
		case codes.Internal:
			if strings.Contains(msg, "GOOGLESQL_RET_CHECK") { // emulator quirk on unknown type (parsed)
				v.Verdict, v.Reason = "accept", "semantic"
			} else {
				v.Verdict, v.Reason = "error", "infra"
			}
		default:
			v.Verdict, v.Reason = "error", "infra"
		}
		return v
	}

	// query / dml
	switch code {
	case codes.InvalidArgument:
		if strings.HasPrefix(msg, "Syntax error:") {
			v.Verdict, v.Reason = "reject", "syntax"
		} else {
			v.Verdict, v.Reason = "accept", "semantic" // Table not found / Unrecognized name / "X is not supported"
		}
	case codes.OutOfRange:
		// A first-row runtime error (e.g. division by zero) means the statement
		// parsed and analyzed fine — grammar ACCEPT.
		v.Verdict, v.Reason = "accept", "semantic"
	default:
		// generic Internal and anything else: the oracle cannot decide.
		v.Verdict, v.Reason = "error", "infra"
	}
	return v
}

// isResourceLevel reports whether a NotFound/AlreadyExists message is about the
// scratch database/instance/session lifecycle (infra) rather than a per-
// statement object reference.
func isResourceLevel(msg string) bool {
	for _, p := range []string{
		"Database not found", "Instance not found", "Session not found",
		"Database already exists", "Instance already exists",
	} {
		if strings.Contains(msg, p) {
			return true
		}
	}
	return false
}

func oneLine(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	if len(s) > 200 {
		cut := 200
		for cut > 0 && !utf8.RuneStart(s[cut]) { // back up to a rune boundary
			cut--
		}
		s = s[:cut] + "…"
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
		//
		// Known limitation: a WITH-led DML (`WITH cte AS (...) INSERT/UPDATE/
		// DELETE ...`) routes here, not through the DML abort path. In practice
		// Spanner does not support WITH-led DML at all (it returns
		// "Syntax error: Unexpected keyword INSERT"), so this routes to a
		// deterministic `reject` — non-authoritative for the BigQuery+Spanner
		// union, so such forms should be triangulated, not trusted from here.
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
