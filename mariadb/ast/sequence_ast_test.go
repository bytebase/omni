package ast

import (
	"reflect"
	"testing"
)

func seqBoolPtr(b bool) *bool { return &b }

// TestSequenceNodesOutput is the outfuncs round-trip assertion for the 5 sequence
// nodes (BYT-9135). Each node is serialized via NodeToString and matched against
// its deterministic S-expression form, which exercises the writeNode dispatch +
// per-node writer and locks the field set.
func TestSequenceNodesOutput(t *testing.T) {
	tests := []struct {
		name string
		node Node
		want string
	}{
		{
			name: "create sequence with value options",
			node: &CreateSequenceStmt{
				Name:      &TableRef{Name: "s"},
				OrReplace: true,
				Start:     &IntLit{Value: 1},
			},
			want: "{CREATE_SEQUENCE :loc 0 :or_replace true :name {TABLEREF :loc 0 :name s} :start {INT_LIT :val 1 :loc 0}}",
		},
		{
			name: "create sequence with no-form flags and nocycle",
			node: &CreateSequenceStmt{
				Name:        &TableRef{Name: "s"},
				IfNotExists: true,
				NoMinValue:  true,
				NoMaxValue:  true,
				NoCache:     true,
				Cycle:       seqBoolPtr(false),
			},
			want: "{CREATE_SEQUENCE :loc 0 :if_not_exists true :name {TABLEREF :loc 0 :name s} :nominvalue true :nomaxvalue true :nocache true :cycle false}",
		},
		{
			name: "alter sequence restart with",
			node: &AlterSequenceStmt{
				Name:        &TableRef{Name: "s"},
				IfExists:    true,
				Restart:     true,
				RestartWith: &IntLit{Value: 500},
			},
			want: "{ALTER_SEQUENCE :loc 0 :if_exists true :name {TABLEREF :loc 0 :name s} :restart true :restart_with {INT_LIT :val 500 :loc 0}}",
		},
		{
			name: "drop sequence multiple",
			node: &DropSequenceStmt{
				IfExists:  true,
				Sequences: []*TableRef{{Name: "a"}, {Name: "b"}},
			},
			want: "{DROP_SEQUENCE :loc 0 :if_exists true :sequences {TABLEREF :loc 0 :name a} {TABLEREF :loc 0 :name b}}",
		},
		{
			name: "next value for",
			node: &NextValueForExpr{Sequence: &TableRef{Name: "s"}},
			want: "{NEXT_VALUE_FOR :loc 0 :sequence {TABLEREF :loc 0 :name s}}",
		},
		{
			name: "previous value for",
			node: &PreviousValueForExpr{Sequence: &TableRef{Name: "s"}},
			want: "{PREVIOUS_VALUE_FOR :loc 0 :sequence {TABLEREF :loc 0 :name s}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NodeToString(tt.node); got != tt.want {
				t.Errorf("NodeToString mismatch\n got: %s\nwant: %s", got, tt.want)
			}
		})
	}
}

// TestSequenceNodesWalkChildren verifies genwalker wired each new node's Node-typed
// children so Walk/Inspect recurses into them (the silent-wrong-lineage guard: a
// missing walk case drops children without error).
func TestSequenceNodesWalkChildren(t *testing.T) {
	tests := []struct {
		name string
		node Node
		want []string
	}{
		{
			name: "create sequence walks name + option exprs",
			node: &CreateSequenceStmt{
				Name:      &TableRef{Name: "s"},
				Start:     &IntLit{Value: 1},
				Increment: &IntLit{Value: 5},
			},
			want: []string{"CreateSequenceStmt", "TableRef", "IntLit"},
		},
		{
			name: "alter sequence walks name + restart_with",
			node: &AlterSequenceStmt{
				Name:        &TableRef{Name: "s"},
				RestartWith: &IntLit{Value: 500},
			},
			want: []string{"AlterSequenceStmt", "TableRef", "IntLit"},
		},
		{
			name: "drop sequence walks each name",
			node: &DropSequenceStmt{
				Sequences: []*TableRef{{Name: "a"}, {Name: "b"}},
			},
			want: []string{"DropSequenceStmt", "TableRef"},
		},
		{
			name: "next value for walks sequence",
			node: &NextValueForExpr{Sequence: &TableRef{Name: "s"}},
			want: []string{"NextValueForExpr", "TableRef"},
		},
		{
			name: "previous value for walks sequence",
			node: &PreviousValueForExpr{Sequence: &TableRef{Name: "s"}},
			want: []string{"PreviousValueForExpr", "TableRef"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			visited := map[string]bool{}
			Inspect(tt.node, func(n Node) bool {
				if n == nil {
					return false
				}
				visited[reflect.TypeOf(n).Elem().Name()] = true
				return true
			})
			for _, w := range tt.want {
				if !visited[w] {
					t.Errorf("expected Walk to visit %s; visited: %v", w, visited)
				}
			}
		})
	}
}
