package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/go-sql-driver/mysql"
)

func TestClassifyTiDBExecError(t *testing.T) {
	cases := []struct {
		name     string
		err      error
		want     Verdict
		wantCode int
		wantMsg  string
	}{
		{"nil is accept", nil, VerdictAccept, 0, ""},
		{"1064 is parse reject", &mysql.MySQLError{Number: 1064, Message: "syntax error"}, VerdictReject, 1064, "syntax error"},
		{"1149 ErrSyntax is parse reject", &mysql.MySQLError{Number: 1149, Message: "You have an error in your SQL syntax"}, VerdictReject, 1149, "You have an error in your SQL syntax"},
		{"1221 grammar-action abort is parse reject", &mysql.MySQLError{Number: 1221, Message: "Incorrect usage of ALL and DISTINCT"}, VerdictReject, 1221, "Incorrect usage of ALL and DISTINCT"},
		{"1492 ast-validator abort is parse reject", &mysql.MySQLError{Number: 1492, Message: "For LIST partitions each partition must be defined"}, VerdictReject, 1492, "For LIST partitions each partition must be defined"},
		{"1102 wrong-db-name is parse reject (runtime collisions fail closed via the label)", &mysql.MySQLError{Number: 1102, Message: "Incorrect database name ''"}, VerdictReject, 1102, "Incorrect database name ''"},
		{"8108 unsupported-type is parsed", &mysql.MySQLError{Number: 8108, Message: "Unsupported type"}, VerdictAccept, 8108, "Unsupported type"},
		{"1146 no-such-table is parsed", &mysql.MySQLError{Number: 1146, Message: "no such table"}, VerdictAccept, 1146, "no such table"},
		{"1105 runtime unknown-error is parsed", &mysql.MySQLError{Number: 1105, Message: "unknown error"}, VerdictAccept, 1105, "unknown error"},
		{"wrapped 1064 unwraps", fmt.Errorf("exec: %w", &mysql.MySQLError{Number: 1064, Message: "syntax error"}), VerdictReject, 1064, "syntax error"},
		{"infra error is none", errors.New("driver: bad connection"), VerdictNone, 0, "driver: bad connection"},

		// Connection-scope: as root, 1045 can only be a failed handshake after a
		// probed batch mutated the credentials — infra, never a parse verdict.
		{"1045 access-denied is connection-scope infra", &mysql.MySQLError{Number: 1045, Message: "Access denied for user 'root'@'127.0.0.1'"}, VerdictNone, 0, "Error 1045: Access denied for user 'root'@'127.0.0.1'"},
		{"wrapped 1045 unwraps to connection-scope infra", fmt.Errorf("exec: %w", &mysql.MySQLError{Number: 1045, Message: "Access denied"}), VerdictNone, 0, "exec: Error 1045: Access denied"},
		// Statement-level schema/privilege codes prove the statement parsed:
		// `USE nonexistent` → 1049, unqualified name with no default schema →
		// 1046, denied database → 1044. None of them are connection-scope here
		// (the DSN carries no default schema and we connect as root).
		{"1049 unknown-database is parsed", &mysql.MySQLError{Number: 1049, Message: "Unknown database 'nonexistent_db_xyz'"}, VerdictAccept, 1049, "Unknown database 'nonexistent_db_xyz'"},
		{"1046 no-database-selected is parsed", &mysql.MySQLError{Number: 1046, Message: "No database selected"}, VerdictAccept, 1046, "No database selected"},
		{"1044 database-access-denied is parsed", &mysql.MySQLError{Number: 1044, Message: "Access denied for user 'root'@'%' to database 'x'"}, VerdictAccept, 1044, "Access denied for user 'root'@'%' to database 'x'"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v, code, msg := classifyTiDBExecError(c.err)
			if v != c.want {
				t.Fatalf("verdict = %q, want %q", v, c.want)
			}
			if code != c.wantCode {
				t.Errorf("code = %d, want %d", code, c.wantCode)
			}
			if msg != c.wantMsg {
				t.Errorf("msg = %q, want %q", msg, c.wantMsg)
			}
		})
	}
}

func TestUnsafeToAdjudicate(t *testing.T) {
	cases := []struct {
		sql  string
		want bool
	}{
		{"shutdown", true},
		{"SHUTDOWN;", true},
		{"ShUtDoWn", true},
		{"kill 5", true},
		{"kill tidb 23123", true},
		{"KILL CONNECTION_ID()", true},
		{"restart", true},
		{"  /* c */ shutdown", true},
		{"-- c\nshutdown", true},
		{"shutdown;select 1", true},
		// Every statement in the batch is scanned: a mid-batch SHUTDOWN kills
		// the oracle just as dead, and a mid-batch SET PASSWORD doesn't kill it
		// at all — it poisons the handshake of every later fresh connection,
		// so the ping-abort backstop never fires.
		{"select 1; shutdown", true},
		{"SELECT 1; SET PASSWORD = 'x'", true},
		// Documented false-positive: the `;` split is naive, so an unsafe word
		// inside a string literal over-matches. Conservative direction only —
		// the row lands in INDETERMINATE unsafe_to_adjudicate, never executes.
		{"SELECT 'x; shutdown'", true},
		{"SELECT 1; SELECT 2", false},

		// First-keyword only: unsafe words elsewhere are fine.
		{"select shutdown_col from t", false},
		{"select * from kill", false},
		{"SELECT 1", false},
		// Identifier characters extend the token: KILLER is not KILL.
		{"killer stmt", false},
		{"shutdown_proc()", false},
		{"", false},

		// SET PASSWORD would change the oracle's credentials mid-sweep.
		{"SET PASSWORD = 'x'", true},
		{"set password for u = 'x'", true},
		// GLOBAL/PERSIST mutations outlive the probe's session: fresh-conn-
		// per-row guards session state only, and e.g. a global sql_mode of
		// ANSI_QUOTES changes how every later probe parses.
		{"SET GLOBAL sql_mode = 'ANSI_QUOTES'", true},
		{"set global max_connections = 100", true},
		{"SET PERSIST sql_mode = ''", true},
		{"SET @@GLOBAL.sql_mode = ''", true},
		{"set @@global.sql_mode=''", true},
		{"SET @@PERSIST.max_connections = 100", true},
		// Session-scoped SET forms stay adjudicable.
		{"SET NAMES utf8", false},
		{"SET sql_mode=''", false},
		{"SET @@session.sql_mode=''", false},
		{"SET @v = 1", false},
		{"select password from t", false},

		// Executable comments: /*! ... */ (with optional version digits) is
		// EXECUTED by TiDB/MySQL, so its content must be scanned, not stripped.
		{"/*! SET PASSWORD = 'x' */", true},
		{"/*!40101 SET GLOBAL sql_mode='ANSI_QUOTES' */", true},
		{"/*! KILL 5 */", true},
		// Safe executable-comment content stays adjudicable.
		{"/*!40101 SET NAMES utf8 */", false},
		// TiDB-specific executable comments: omni's own splitter treats /*T!
		// as executable SQL too (tidb/parser/split.go, Segment.Empty), and its
		// prefix match covers digit forms like /*T!50000 the same way.
		{"/*T! SET GLOBAL sql_mode='ANSI_QUOTES' */", true},
		{"/*T!50000 KILL 5 */", true},
		// Safe TiDB executable-comment content stays adjudicable.
		{"/*T!40101 SET NAMES utf8 */", false},
		// Ordinary comments still strip — their content is never executed.
		{"/* comment */ SELECT 1", false},
		{"/* SET PASSWORD = 'x' */ SELECT 1", false},
		// Mixed: a leading ordinary comment must not mask an executable one.
		{"/* c */ /*!40101 SET GLOBAL sql_mode='' */", true},
		// The server treats ordinary comments as whitespace ANYWHERE, not just
		// statement-leading: a comment between keywords must not blind the scan.
		{"SET /*c*/ PASSWORD = 'x'", true},
		{"SET /*a*/ /*b*/ GLOBAL sql_mode=''", true},
		{"CREATE /*c*/ USER 'root'@'127.0.0.1' IDENTIFIED BY 'x'", true},
		{"ALTER -- c\nUSER u IDENTIFIED BY 'x'", true},
		// An unsafe statement hiding inside an ORDINARY comment is never
		// executed by the server, and comments are stripped before the `;`
		// split, so the commented-out `;` cannot fabricate a phantom segment:
		// false is genuinely correct here, not an accepted miss.
		{"SELECT 1 /* ; SET PASSWORD='x' */", false},
		// Comment openers inside string literals are literal text, not
		// comments: stripping them would swallow the real statements that
		// follow — the one direction the deny-list must never err.
		{"SELECT 'a -- b'; KILL 5", true},
		{"SELECT '/*'; KILL 5", true},

		// Account-mutation DCL poisons later handshakes (identity channel).
		{"CREATE USER 'root'@'127.0.0.1' IDENTIFIED BY 'x'", true},
		{"ALTER USER root IDENTIFIED BY 'x'", true},
		{"DROP USER u", true},
		{"RENAME USER a TO b", true},
		// Second keyword must be USER: ordinary DDL stays adjudicable.
		{"CREATE TABLE user (a int)", false},
		{"DROP TABLE users", false},
		// GRANT/REVOKE cannot change our connection identity — not blocked.
		{"GRANT ALL ON *.* TO u", false},
		{"REVOKE ALL ON *.* FROM u", false},
	}
	for _, c := range cases {
		if got := unsafeToAdjudicate(c.sql); got != c.want {
			t.Errorf("unsafeToAdjudicate(%q) = %v, want %v", c.sql, got, c.want)
		}
	}
}

func TestNormalizeTiDBDSN(t *testing.T) {
	// H1: the sweep is incorrect without multiStatements (multi-statement
	// corpus rows would false-reject with 1064) and unbounded without
	// timeouts, so both are forced onto any DSN.
	got, err := normalizeTiDBDSN("root@tcp(127.0.0.1:14001)/test")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"multiStatements=true", "timeout=5s", "readTimeout=10s", "writeTimeout=10s"} {
		if !strings.Contains(got, want) {
			t.Errorf("normalized DSN %q missing %q", got, want)
		}
	}

	// A default schema is droppable by adjudicated DDL, after which every
	// later fresh connection fails its handshake with 1049 — so the schema is
	// forced empty even when a user-supplied DSN names one.
	cfg, err := mysql.ParseDSN(got)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DBName != "" {
		t.Errorf("normalized DSN %q kept default schema %q, want none", got, cfg.DBName)
	}

	// Explicit timeouts are respected, not clobbered.
	got, err = normalizeTiDBDSN("root@tcp(127.0.0.1:14001)/test?readTimeout=3s")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "readTimeout=3s") {
		t.Errorf("normalized DSN %q clobbered explicit readTimeout", got)
	}

	if _, err := normalizeTiDBDSN("not a dsn ("); err == nil {
		t.Error("invalid DSN did not error")
	}
}

func TestPrepareAdjudication(t *testing.T) {
	t.Run("unsafe row is marked indeterminate without touching the container", func(t *testing.T) {
		r := Row{
			SQL: "kill tidb 23123", Expected: VerdictAccept, OmniVerdict: VerdictReject,
			OmniError: "unexpected token", Class: ClassGap, DivergenceKey: "unexpected token",
		}
		if prepareAdjudication(&r) {
			t.Fatal("unsafe row was cleared for execution")
		}
		if r.Class != ClassIndeterminate {
			t.Errorf("class = %q, want %q", r.Class, ClassIndeterminate)
		}
		if r.ClassifierReason != "unsafe_to_adjudicate" {
			t.Errorf("reason = %q, want unsafe_to_adjudicate", r.ClassifierReason)
		}
		if r.EngineVerdict != VerdictNone {
			t.Errorf("engine verdict = %q, want empty", r.EngineVerdict)
		}
		if r.DivergenceKey != "" {
			t.Errorf("divergence key = %q, want empty (INDETERMINATE rows are not clustered)", r.DivergenceKey)
		}
	})

	t.Run("H3: duplicate_label_conflict clears the arbitrary label", func(t *testing.T) {
		r := Row{
			SQL: "SELECT 1", Expected: VerdictReject, OmniVerdict: VerdictAccept,
			Class: ClassIndeterminate, ClassifierReason: "duplicate_label_conflict",
		}
		if !prepareAdjudication(&r) {
			t.Fatal("dup-label row was not cleared for execution")
		}
		if r.Expected != VerdictNone {
			t.Errorf("expected label = %q, want cleared", r.Expected)
		}
	})

	t.Run("ordinary row passes through untouched", func(t *testing.T) {
		r := Row{SQL: "SELECT 1", Expected: VerdictAccept, OmniVerdict: VerdictReject, Class: ClassGap}
		if !prepareAdjudication(&r) {
			t.Fatal("ordinary row was not cleared for execution")
		}
		if r.Expected != VerdictAccept || r.Class != ClassGap {
			t.Errorf("row mutated: expected=%q class=%q", r.Expected, r.Class)
		}
	})
}

func TestApplyContainerVerdict(t *testing.T) {
	t.Run("H3 end to end: container is sole truth after label clear", func(t *testing.T) {
		// A dup-label row kept the arbitrary first-seen label (reject). If the
		// label survived, container-accept would flip a coin into
		// label_container_disagree; H3 clears it so the container alone rules.
		r := Row{
			SQL: "SELECT 1", Expected: VerdictReject, OmniVerdict: VerdictAccept,
			Class: ClassIndeterminate, ClassifierReason: "duplicate_label_conflict",
		}
		if !prepareAdjudication(&r) {
			t.Fatal("dup-label row was not cleared for execution")
		}
		if infra := applyContainerVerdict(&r, nil); infra {
			t.Fatal("nil exec error reported as infra")
		}
		if r.Class != ClassAgreeAccept {
			t.Fatalf("class = %q, want %q (container verdict alone)", r.Class, ClassAgreeAccept)
		}
		if r.ClassifierReason != "" {
			t.Errorf("reason = %q, want empty", r.ClassifierReason)
		}
	})

	t.Run("engine reject re-keys an OVER row on the engine message", func(t *testing.T) {
		r := Row{
			SQL: "ALTER TABLE t BOGUS", Expected: VerdictReject, OmniVerdict: VerdictAccept,
			Class: ClassOver, DivergenceKey: "ALTER TABLE T BOGUS",
		}
		err := &mysql.MySQLError{Number: 1064, Message: "You have an error in your SQL syntax near 'BOGUS' at line 1"}
		if infra := applyContainerVerdict(&r, err); infra {
			t.Fatal("MySQL error reported as infra")
		}
		if r.Class != ClassOver {
			t.Fatalf("class = %q, want %q", r.Class, ClassOver)
		}
		if r.EngineVerdict != VerdictReject || r.RawErrorCode != 1064 {
			t.Errorf("engine verdict/code = %q/%d, want reject/1064", r.EngineVerdict, r.RawErrorCode)
		}
		if want := "You have an error in your SQL syntax near ? at line N"; r.DivergenceKey != want {
			t.Errorf("divergence key = %q, want %q", r.DivergenceKey, want)
		}
	})

	t.Run("label vs container disagreement stays fail-closed", func(t *testing.T) {
		r := Row{
			SQL: "SELECT 1", Expected: VerdictReject, OmniVerdict: VerdictReject,
			Class: ClassAgreeReject,
		}
		if infra := applyContainerVerdict(&r, nil); infra {
			t.Fatal("nil exec error reported as infra")
		}
		if r.Class != ClassIndeterminate || r.ClassifierReason != "label_container_disagree" {
			t.Fatalf("class/reason = %q/%q, want INDETERMINATE/label_container_disagree", r.Class, r.ClassifierReason)
		}
	})

	t.Run("infra error is fail-closed indeterminate, never a verdict", func(t *testing.T) {
		r := Row{
			SQL: "SELECT SLEEP(100)", Expected: VerdictAccept, OmniVerdict: VerdictReject,
			Class: ClassGap, DivergenceKey: "unexpected token",
		}
		infra := applyContainerVerdict(&r, errors.New("read tcp 127.0.0.1: i/o timeout"))
		if !infra {
			t.Fatal("infra error not reported as infra")
		}
		if r.Class != ClassIndeterminate || r.ClassifierReason != "infra_error" {
			t.Fatalf("class/reason = %q/%q, want INDETERMINATE/infra_error", r.Class, r.ClassifierReason)
		}
		if r.EngineVerdict != VerdictNone {
			t.Errorf("engine verdict = %q, want empty", r.EngineVerdict)
		}
		if r.RawErrorMessage != "read tcp 127.0.0.1: i/o timeout" {
			t.Errorf("raw error message = %q, want the driver error text", r.RawErrorMessage)
		}
		if r.DivergenceKey != "" {
			t.Errorf("divergence key = %q, want empty", r.DivergenceKey)
		}
	})
}
