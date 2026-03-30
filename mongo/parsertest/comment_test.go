package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/parser"
)

func TestComments(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int // expected number of statements
	}{
		{
			name:  "line comment before statement",
			input: "// this is a comment\ndb.users.find()",
			want:  1,
		},
		{
			name:  "line comment after statement",
			input: "db.users.find() // trailing comment",
			want:  1,
		},
		{
			name:  "block comment before statement",
			input: "/* block */ db.users.find()",
			want:  1,
		},
		{
			name:  "block comment after statement",
			input: "db.users.find() /* trailing */",
			want:  1,
		},
		{
			name:  "multiline block comment",
			input: "/*\n * multi\n * line\n */\ndb.users.find()",
			want:  1,
		},
		{
			name:  "comment between statements",
			input: "db.users.find()\n// separator\ndb.orders.find()",
			want:  2,
		},
		{
			name:  "block comment between statements",
			input: "db.users.find()\n/* gap */\ndb.orders.find()",
			want:  2,
		},
		{
			name:  "only comments no statements",
			input: "// just a comment\n/* another one */",
			want:  0,
		},
		{
			name: "comment inside document value (as string)",
			input: `db.users.insertOne({ note: "// not a comment" })`,
			want:  1,
		},
		{
			name: "comment inside string single quotes",
			input: `db.users.insertOne({ note: '/* not a comment */' })`,
			want:  1,
		},
		{
			name:  "inline comment in argument list",
			input: "db.users.find(\n  { name: \"alice\" } /* filter */,\n  { _id: 0 } // projection\n)",
			want:  1,
		},
		{
			name:  "comment inside nested document",
			input: "db.users.find({\n  // match active users\n  status: \"active\",\n  /* check age */\n  age: { $gt: 18 }\n})",
			want:  1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nodes, err := parser.Parse(tc.input)
			if tc.want == 0 {
				if err != nil {
					t.Fatalf("parse error: %v", err)
				}
				if len(nodes) != 0 {
					t.Errorf("expected 0 nodes, got %d", len(nodes))
				}
				return
			}
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			if len(nodes) != tc.want {
				t.Errorf("expected %d nodes, got %d", tc.want, len(nodes))
			}
		})
	}
}
