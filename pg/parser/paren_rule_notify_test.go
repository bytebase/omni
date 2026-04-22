package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// TestParenRuleNotifyAction covers KB-2a: parseRuleActionStmt must handle
// NotifyStmt as a valid RuleActionStmt alternative.
//
// PG grammar:
//
//	RuleActionStmt:
//	    SelectStmt | InsertStmt | UpdateStmt | DeleteStmt | NotifyStmt
//
// Before the fix, `CREATE RULE ... DO [INSTEAD|ALSO] NOTIFY foo` silently
// produced a RuleStmt with nil/empty Actions plus a trailing NotifyStmt
// as a separate top-level RawStmt (Parse()'s outer loop reprocessed the
// unconsumed tail). This test asserts that (a) we get exactly one
// statement, (b) the RuleStmt carries a NotifyStmt action, and (c) the
// Instead/Replace flags are honored.
func TestParenRuleNotifyAction(t *testing.T) {
	cases := []struct {
		name    string
		sql     string
		replace bool
		instead bool
		cond    string
		event   nodes.CmdType
	}{
		{
			// copydml.sql:69
			name:    "do_instead_notify_insert",
			sql:     "create rule qqq as on insert to copydml_test do instead notify copydml_test",
			instead: true,
			cond:    "copydml_test",
			event:   nodes.CMD_INSERT,
		},
		{
			// rules.sql:1017 — bare DO (no INSTEAD / no ALSO)
			name:  "do_notify_delete",
			sql:   "create rule r4 as on delete to rules_src do notify rules_src_deletion",
			cond:  "rules_src_deletion",
			event: nodes.CMD_DELETE,
		},
		{
			// with.sql:1721
			name:    "or_replace_do_instead_notify",
			sql:     "CREATE OR REPLACE RULE y_rule AS ON INSERT TO y DO INSTEAD NOTIFY foo",
			replace: true,
			instead: true,
			cond:    "foo",
			event:   nodes.CMD_INSERT,
		},
		{
			// with.sql:1726 — DO ALSO keeps Instead=false
			name:    "or_replace_do_also_notify",
			sql:     "CREATE OR REPLACE RULE y_rule AS ON INSERT TO y DO ALSO NOTIFY foo",
			replace: true,
			instead: false,
			cond:    "foo",
			event:   nodes.CMD_INSERT,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			if result == nil || len(result.Items) != 1 {
				n := 0
				if result != nil {
					n = len(result.Items)
				}
				t.Fatalf("expected exactly 1 statement, got %d (KB-2a regression: NOTIFY body leaked to a second RawStmt)", n)
			}
			raw, ok := result.Items[0].(*nodes.RawStmt)
			if !ok {
				t.Fatalf("expected *nodes.RawStmt, got %T", result.Items[0])
			}
			stmt, ok := raw.Stmt.(*nodes.RuleStmt)
			if !ok {
				t.Fatalf("expected *nodes.RuleStmt, got %T", raw.Stmt)
			}
			if stmt.Replace != tc.replace {
				t.Errorf("Replace: got %v, want %v", stmt.Replace, tc.replace)
			}
			if stmt.Instead != tc.instead {
				t.Errorf("Instead: got %v, want %v", stmt.Instead, tc.instead)
			}
			if stmt.Event != tc.event {
				t.Errorf("Event: got %d, want %d", stmt.Event, tc.event)
			}
			if stmt.Actions == nil || len(stmt.Actions.Items) != 1 {
				got := 0
				if stmt.Actions != nil {
					got = len(stmt.Actions.Items)
				}
				t.Fatalf("expected 1 action in RuleStmt.Actions, got %d", got)
			}
			notify, ok := stmt.Actions.Items[0].(*nodes.NotifyStmt)
			if !ok {
				t.Fatalf("expected *nodes.NotifyStmt action, got %T", stmt.Actions.Items[0])
			}
			if notify.Conditionname != tc.cond {
				t.Errorf("Conditionname: got %q, want %q", notify.Conditionname, tc.cond)
			}
			if notify.Payload != "" {
				t.Errorf("Payload: got %q, want empty", notify.Payload)
			}
		})
	}
}
