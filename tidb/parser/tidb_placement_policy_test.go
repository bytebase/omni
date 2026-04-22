package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// TestParsePlacementPolicyDDL covers the three new statements and the
// surface of the shared option list. Each case asserts on AST shape,
// not on round-trip string equality — outfuncs stability is covered by
// writeCreatePlacementPolicyStmt / writeAlterPlacementPolicyStmt tests
// downstream in the ast package.
func TestParsePlacementPolicyDDL(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
		check   func(t *testing.T, list *nodes.List)
	}{
		{
			name: "create_minimal",
			sql:  "CREATE PLACEMENT POLICY p1 PRIMARY_REGION = 'us-east-1'",
			check: func(t *testing.T, list *nodes.List) {
				s, ok := list.Items[0].(*nodes.CreatePlacementPolicyStmt)
				if !ok {
					t.Fatalf("want CreatePlacementPolicyStmt, got %T", list.Items[0])
				}
				if s.Name != "p1" {
					t.Errorf("Name = %q, want p1", s.Name)
				}
				if len(s.Options) != 1 || s.Options[0].Name != "PRIMARY_REGION" || s.Options[0].Value != "us-east-1" {
					t.Errorf("options: %+v", s.Options)
				}
			},
		},
		{
			name: "create_or_replace_if_not_exists",
			sql:  "CREATE OR REPLACE PLACEMENT POLICY IF NOT EXISTS p1 FOLLOWERS = 3",
			check: func(t *testing.T, list *nodes.List) {
				s := list.Items[0].(*nodes.CreatePlacementPolicyStmt)
				if !s.OrReplace {
					t.Error("OrReplace should be true")
				}
				if !s.IfNotExists {
					t.Error("IfNotExists should be true")
				}
				if len(s.Options) != 1 || !s.Options[0].IsInt || s.Options[0].IntValue != 3 {
					t.Errorf("FOLLOWERS option: %+v", s.Options[0])
				}
			},
		},
		{
			name: "create_multi_option_comma",
			sql:  "CREATE PLACEMENT POLICY p PRIMARY_REGION='us', REGIONS='us,eu', FOLLOWERS=2",
			check: func(t *testing.T, list *nodes.List) {
				s := list.Items[0].(*nodes.CreatePlacementPolicyStmt)
				if len(s.Options) != 3 {
					t.Fatalf("want 3 options, got %d", len(s.Options))
				}
			},
		},
		{
			name: "create_multi_option_whitespace",
			// TiDB's Restore output omits commas between options.
			sql: "CREATE PLACEMENT POLICY p PRIMARY_REGION='us' REGIONS='us,eu' LEARNERS=1",
			check: func(t *testing.T, list *nodes.List) {
				s := list.Items[0].(*nodes.CreatePlacementPolicyStmt)
				if len(s.Options) != 3 {
					t.Fatalf("want 3 options, got %d", len(s.Options))
				}
			},
		},
		{
			// The upstream grammar accepts mixing comma and whitespace
			// separators in the same option list (parser.y:2000/2004/2008).
			// Neither comma-only nor whitespace-only covers this.
			name: "create_multi_option_mixed_separators",
			sql:  "CREATE PLACEMENT POLICY p PRIMARY_REGION='us', REGIONS='us,eu' FOLLOWERS=2",
			check: func(t *testing.T, list *nodes.List) {
				s := list.Items[0].(*nodes.CreatePlacementPolicyStmt)
				if len(s.Options) != 3 {
					t.Fatalf("want 3 options, got %d: %+v", len(s.Options), s.Options)
				}
			},
		},
		{
			name: "create_constraints_json_dict",
			sql:  `CREATE PLACEMENT POLICY p CONSTRAINTS='{"+region=us-east-1":2,"+region=us-west-1":1}'`,
			check: func(t *testing.T, list *nodes.List) {
				s := list.Items[0].(*nodes.CreatePlacementPolicyStmt)
				if len(s.Options) != 1 || s.Options[0].Name != "CONSTRAINTS" {
					t.Fatalf("options: %+v", s.Options)
				}
				// Value is the raw string literal content (no quotes).
				if !strings.Contains(s.Options[0].Value, "us-east-1") {
					t.Errorf("CONSTRAINTS value lost contents: %q", s.Options[0].Value)
				}
			},
		},
		{
			name: "create_constraints_json_list",
			sql:  `CREATE PLACEMENT POLICY p CONSTRAINTS='[+region=us-east-1,-region=us-west-1]'`,
			check: func(t *testing.T, list *nodes.List) {
				s := list.Items[0].(*nodes.CreatePlacementPolicyStmt)
				if !strings.Contains(s.Options[0].Value, "+region=us-east-1") {
					t.Errorf("CONSTRAINTS list form lost contents: %q", s.Options[0].Value)
				}
			},
		},
		{
			name: "create_no_equals",
			// The `=` is optional for every option per EqOpt rule.
			sql: "CREATE PLACEMENT POLICY p PRIMARY_REGION 'us' FOLLOWERS 3",
			check: func(t *testing.T, list *nodes.List) {
				s := list.Items[0].(*nodes.CreatePlacementPolicyStmt)
				if len(s.Options) != 2 {
					t.Fatalf("want 2 options, got %d", len(s.Options))
				}
			},
		},
		{
			name: "create_all_constraint_variants",
			sql: `CREATE PLACEMENT POLICY p ` +
				`LEADER_CONSTRAINTS='[+region=us-east]' ` +
				`FOLLOWER_CONSTRAINTS='{+region=us-west: 2}' ` +
				`VOTER_CONSTRAINTS='[+zone=a]' ` +
				`LEARNER_CONSTRAINTS='[+zone=b]' ` +
				`SURVIVAL_PREFERENCES='[region, zone]'`,
			check: func(t *testing.T, list *nodes.List) {
				s := list.Items[0].(*nodes.CreatePlacementPolicyStmt)
				if len(s.Options) != 5 {
					t.Fatalf("want 5 options, got %d", len(s.Options))
				}
			},
		},
		{
			// Upstream parser.y:2025-2028 rejects FOLLOWERS=0 at parse
			// time. Our parser matches to avoid accepting SQL that real
			// TiDB would reject.
			name:    "error_followers_zero",
			sql:     "CREATE PLACEMENT POLICY p FOLLOWERS = 0",
			wantErr: true,
		},
		{
			// VOTERS = 0 / LEARNERS = 0 are NOT parse errors upstream —
			// asymmetric with FOLLOWERS, but upstream-consistent.
			name: "voters_zero_ok",
			sql:  "CREATE PLACEMENT POLICY p VOTERS = 0",
			check: func(t *testing.T, list *nodes.List) {
				s := list.Items[0].(*nodes.CreatePlacementPolicyStmt)
				if len(s.Options) != 1 || s.Options[0].IntValue != 0 {
					t.Errorf("VOTERS=0 option: %+v", s.Options[0])
				}
			},
		},
		{
			name: "alter_with_if_exists",
			sql:  "ALTER PLACEMENT POLICY IF EXISTS p1 PRIMARY_REGION = 'us-west'",
			check: func(t *testing.T, list *nodes.List) {
				s, ok := list.Items[0].(*nodes.AlterPlacementPolicyStmt)
				if !ok {
					t.Fatalf("want AlterPlacementPolicyStmt, got %T", list.Items[0])
				}
				if !s.IfExists {
					t.Error("IfExists should be true")
				}
				if s.Name != "p1" {
					t.Errorf("Name = %q", s.Name)
				}
			},
		},
		{
			name: "drop_basic",
			sql:  "DROP PLACEMENT POLICY p1",
			check: func(t *testing.T, list *nodes.List) {
				s, ok := list.Items[0].(*nodes.DropPlacementPolicyStmt)
				if !ok {
					t.Fatalf("want DropPlacementPolicyStmt, got %T", list.Items[0])
				}
				if s.Name != "p1" {
					t.Errorf("Name = %q", s.Name)
				}
			},
		},
		{
			name: "drop_if_exists",
			sql:  "DROP PLACEMENT POLICY IF EXISTS p1",
			check: func(t *testing.T, list *nodes.List) {
				s := list.Items[0].(*nodes.DropPlacementPolicyStmt)
				if !s.IfExists {
					t.Error("IfExists should be true")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := Parse(tt.sql)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			if list == nil || len(list.Items) == 0 {
				t.Fatal("empty parse result")
			}
			if tt.check != nil {
				tt.check(t, list)
			}
		})
	}
}
