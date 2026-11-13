package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestIndexTypeAllowList(t *testing.T) {
	notAllowed := func(indexType, text string) indexFinding {
		return indexFinding{Code: v1review.IndexTypeNotAllowed, Type: indexType, Text: text}
	}
	for _, tt := range []struct {
		name  string
		sql   string
		allow []string
		want  []indexFinding
	}{
		{name: "BTREE by default", sql: "CREATE TABLE t (a INT, PRIMARY KEY (a), KEY idx (a))", allow: []string{"BTREE"}},
		{
			name:  "each kind's default type",
			sql:   "CREATE TABLE t (a TEXT, b GEOMETRY NOT NULL,\n  KEY idx (a),\n  FULLTEXT KEY ft (a),\n  SPATIAL KEY sp (b))",
			allow: []string{"HASH"},
			want: []indexFinding{
				notAllowed("BTREE", "KEY idx (a)"),
				notAllowed("FULLTEXT", "FULLTEXT KEY ft (a)"),
				notAllowed("SPATIAL", "SPATIAL KEY sp (b)"),
			},
		},
		{
			name:  "USING names the type, compared as written",
			sql:   "CREATE TABLE t (a INT, KEY idx USING HASH (a))",
			allow: []string{"hash"},
			want:  []indexFinding{notAllowed("HASH", "KEY idx USING HASH (a)")},
		},
		{name: "foreign keys and column keys are not judged", sql: "CREATE TABLE t (a INT PRIMARY KEY, b INT UNIQUE, FOREIGN KEY (b) REFERENCES p (b))"},
		{
			name:  "alter table anchors the statement",
			sql:   "ALTER TABLE t ADD FULLTEXT INDEX ft (a)",
			allow: []string{"BTREE"},
			want:  []indexFinding{notAllowed("FULLTEXT", "ALTER TABLE t ADD FULLTEXT INDEX ft (a)")},
		},
		{
			name:  "create index of every kind",
			sql:   "CREATE SPATIAL INDEX sp ON t (a)",
			allow: []string{"BTREE"},
			want:  []indexFinding{notAllowed("SPATIAL", "CREATE SPATIAL INDEX sp ON t (a)")},
		},
		{name: "allowed", sql: "CREATE INDEX idx USING HASH ON t (a)", allow: []string{"HASH"}},
		{
			name:  "USING after the key parts is not read",
			sql:   "CREATE INDEX idx ON t (a) USING HASH",
			allow: []string{"HASH"},
			want:  []indexFinding{notAllowed("BTREE", "CREATE INDEX idx ON t (a) USING HASH")},
		},
		{name: "a routine body is not entered", sql: "CREATE PROCEDURE p() BEGIN CREATE INDEX idx ON t (a); END"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			sqls := []string{tt.sql}
			got := indexFindings(sqls, IndexTypeAllowList(parseEach(t, sqls), tt.allow))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
