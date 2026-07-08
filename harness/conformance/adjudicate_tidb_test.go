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
		// First-statement-only is deliberate: a mid-batch unsafe statement is
		// caught by the container-death abort backstop, not this predicate.
		{"select 1; shutdown", false},

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
		// Other SET forms stay adjudicable.
		{"SET NAMES utf8", false},
		{"SET sql_mode=''", false},
		{"select password from t", false},
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
