package v1review

import (
	"testing"

	"github.com/bytebase/omni/v1review"
)

func TestStatementDisallowAddNotNull(t *testing.T) {
	testStatementRule(t, StatementDisallowAddNotNull, []statementCase{
		{name: "set not null", sql: "ALTER TABLE t ALTER COLUMN a SET NOT NULL;", want: []stmtFinding{{Code: v1review.NotNullSet, Name: "a"}}},
		{name: "folded and quoted", sql: "ALTER TABLE t ALTER COLUMN Abc SET NOT NULL, ALTER \"Def\" SET NOT NULL;", want: []stmtFinding{{Code: v1review.NotNullSet, Name: "abc"}, {Code: v1review.NotNullSet, Name: "Def"}}},
		{name: "drop not null not judged", sql: "ALTER TABLE t ALTER COLUMN a DROP NOT NULL;"},
		{name: "added not null column not judged", sql: "ALTER TABLE t ADD COLUMN a int NOT NULL;"},
	})
}
