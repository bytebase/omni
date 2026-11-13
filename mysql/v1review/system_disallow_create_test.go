package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

func TestSystemDisallowCreate(t *testing.T) {
	const sql = "CREATE PROCEDURE p() BEGIN CREATE VIEW w AS SELECT 1; END;\n" +
		"CREATE FUNCTION f() RETURNS INT RETURN 1;\n" +
		"CREATE AGGREGATE FUNCTION g RETURNS STRING SONAME 'g.so';\n" +
		"CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1;\n" +
		"CREATE VIEW v AS SELECT 1;\n" +
		"CREATE OR REPLACE VIEW v AS SELECT 2;\n" +
		"ALTER VIEW v AS SELECT 3;\n" +
		"ALTER PROCEDURE p COMMENT 'x';\n" +
		"ALTER EVENT e DISABLE;\n" +
		"CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;\n" +
		"CREATE TABLE t (a INT);"
	for _, tt := range []struct {
		name string
		rule func([]ast.Node) []v1review.Finding
		want []systemFinding
	}{
		{name: "procedure", rule: SystemProcedureDisallowCreate, want: []systemFinding{
			{Code: v1review.ProcedureCreated, Text: "CREATE PROCEDURE p() BEGIN CREATE VIEW w AS SELECT 1; END"},
		}},
		{name: "function", rule: SystemFunctionDisallowCreate, want: []systemFinding{
			{Code: v1review.FunctionCreated, Text: "CREATE FUNCTION f() RETURNS INT RETURN 1"},
			{Code: v1review.FunctionCreated, Text: "CREATE AGGREGATE FUNCTION g RETURNS STRING SONAME 'g.so'"},
		}},
		{name: "event", rule: SystemEventDisallowCreate, want: []systemFinding{
			{Code: v1review.EventCreated, Text: "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1"},
		}},
		{name: "view, not in a routine body and not altered", rule: SystemViewDisallowCreate, want: []systemFinding{
			{Code: v1review.ViewCreated, Text: "CREATE VIEW v AS SELECT 1"},
			{Code: v1review.ViewCreated, Text: "CREATE OR REPLACE VIEW v AS SELECT 2"},
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, judgeSystem(t, sql, tt.rule)); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
