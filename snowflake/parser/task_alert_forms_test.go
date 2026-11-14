package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// The CREATE / ALTER TASK and ALERT forms the docs give beyond the original
// T4.3 shapes.

func TestParseTaskContactsAndSessionParams(t *testing.T) {
	c := mustCreateTask(t, "CREATE TASK t WITH TAG (x = 'v') WITH CONTACT (STEWARD = a, SUPPORT = b) "+
		"WAREHOUSE = wh TIMEZONE = 'UTC', QUERY_TAG = 'q' AS SELECT 1")
	if len(c.Contacts) != 2 || c.Contacts[0].Name != "STEWARD" || c.Contacts[1].Words != "B" {
		t.Errorf("contacts = %+v", c.Contacts)
	}
	// A contact is a schema object: db.schema.name may qualify it.
	mustCreateTask(t, "CREATE TASK t WITH CONTACT (STEWARD = governance.contacts.owner) AS SELECT 1")
	mustAlterTask(t, "ALTER TASK t SET CONTACT STEWARD = governance.contacts.owner")
	if len(c.Options) != 3 {
		t.Errorf("options = %+v", c.Options)
	}

	a := mustAlterTask(t, "ALTER TASK t SET LOG_LEVEL = 'INFO' TIMEZONE = 'UTC', QUERY_TAG = 'q' CONTACT STEWARD = a, SUPPORT = b EXECUTE AS USER u")
	if a.Action != ast.AlterTaskSet || len(a.Options) != 5 {
		t.Fatalf("got %+v", a)
	}
	if o := a.Options[3]; o.Name != "CONTACT" || len(o.Group) != 2 {
		t.Errorf("contact = %+v", o)
	}
	if o := a.Options[4]; o.Name != "EXECUTE AS USER" || o.Words != "u" {
		t.Errorf("execute as = %+v", o)
	}

	u := mustAlterTask(t, "ALTER TASK t UNSET CONTACT STEWARD, EXECUTE AS USER, COMMENT")
	if u.Action != ast.AlterTaskUnset || len(u.UnsetProps) != 3 || u.UnsetProps[0] != "CONTACT STEWARD" || u.UnsetProps[1] != "EXECUTE AS USER" {
		t.Errorf("got %+v", u)
	}

	e := mustCreateTask(t, "CREATE TASK t WAREHOUSE = wh EXECUTE AS USER u WHEN x > 1 AS SELECT 1")
	if len(e.Options) != 2 || e.Options[1].Name != "EXECUTE AS USER" || e.Options[1].Words != "u" || e.When == nil || e.BodyRaw != "SELECT 1" {
		t.Errorf("got %+v", e)
	}

	r := mustAlterTask(t, "ALTER TASK IF EXISTS t REMOVE WHEN")
	if r.Action != ast.AlterTaskRemoveWhen || !r.IfExists {
		t.Errorf("got %+v", r)
	}

	for _, in := range []string{
		"CREATE TASK t WITH CONTACT () AS SELECT 1",
		"CREATE TASK t WITH CONTACT (STEWARD = 'a') AS SELECT 1",
		"CREATE TASK t TIMEZONE = 'UTC', AS SELECT 1",
		"ALTER TASK t SET TIMEZONE = 'UTC',",
		"ALTER TASK t SET CONTACT STEWARD",
		"ALTER TASK t SET EXECUTE AS USER",
		"ALTER TASK t UNSET EXECUTE AS",
		"ALTER TASK t REMOVE WHEN x",
		"CREATE TASK t zzz AS SELECT 1", // a task's options all take a value
		"CREATE TASK t WAREHOUSE = wh zzz AS SELECT 1",
		"ALTER TASK t SET zzz",
		"ALTER TASK t SET WAREHOUSE wh",
	} {
		mustRejectObject(t, in)
	}
}

func TestParseAlertTemplateCloneAndTags(t *testing.T) {
	c := mustCreateAlert(t, "CREATE OR REPLACE ALERT a CLONE db.s.b")
	if c.Clone == nil || c.Clone.String() != "db.s.b" || !c.OrReplace {
		t.Errorf("got %+v", c)
	}
	f := mustCreateAlert(t, "CREATE ALERT IF NOT EXISTS a FROM TEMPLATE tpl WITH TAG (t = 'v') "+
		"WAREHOUSE = wh SCHEDULE = '1 MINUTE' COMMENT = 'c' RUNBOOK = 'r' SUSPEND_ALERT_AFTER_NUM_FAILURES = 3 TEMPLATE_PARAMS = '{}'")
	if f.Template == nil || f.Template.String() != "tpl" || len(f.Tags) != 1 || len(f.Options) != 6 || f.ConditionRaw != "" {
		t.Errorf("got %+v", f)
	}
	mustCreateAlert(t, "CREATE ALERT a FROM TEMPLATE tpl")

	s := mustAlterAlert(t, "ALTER ALERT IF EXISTS a FROM TEMPLATE tpl SET WAREHOUSE = wh TEMPLATE_PARAMS = '{}'")
	if s.Action != ast.AlterAlertSet || s.Template == nil || len(s.Options) != 2 {
		t.Errorf("got %+v", s)
	}
	st := mustAlterAlert(t, "ALTER ALERT a SET TAG t = 'v', u = 'w'")
	if st.Action != ast.AlterAlertSetTag || len(st.Tags) != 2 {
		t.Errorf("got %+v", st)
	}
	ut := mustAlterAlert(t, "ALTER ALERT a UNSET TAG t, u")
	if ut.Action != ast.AlterAlertUnsetTag || len(ut.UnsetTags) != 2 {
		t.Errorf("got %+v", ut)
	}

	e := mustParseOne(t, "EXECUTE ALERT db.s.a")
	if x, ok := e.(*ast.ExecuteAlertStmt); !ok || x.Name.String() != "db.s.a" || x.Loc.End != len("EXECUTE ALERT db.s.a") {
		t.Errorf("got %#v", e)
	}

	for _, in := range []string{
		"CREATE ALERT IF NOT EXISTS a CLONE b",
		"CREATE OR ALTER ALERT a CLONE b",
		"CREATE OR ALTER ALERT a FROM TEMPLATE tpl",
		"CREATE ALERT a FROM TEMPLATE db.tpl",
		"CREATE ALERT a FROM TEMPLATE tpl CONFIG = 'x'",
		"CREATE ALERT a FROM TEMPLATE tpl WAREHOUSE = 'wh'",
		"CREATE ALERT a FROM TEMPLATE tpl COMMENT = 'a' COMMENT = 'b'",
		"CREATE ALERT a FROM TEMPLATE tpl IF (EXISTS (SELECT 1)) THEN SELECT 1",
		"ALTER ALERT a FROM TEMPLATE tpl SET",
		"ALTER ALERT a FROM TEMPLATE tpl SET NOT_A_PROPERTY = 1",
		"ALTER ALERT a FROM TEMPLATE tpl RESUME",
		"ALTER ALERT a UNSET TAG",
		"EXECUTE ALERT",
		"EXECUTE ALERT a b",
	} {
		mustRejectObject(t, in)
	}
}

func TestParseShowCommittedTasksInGraph(t *testing.T) {
	s := mustParseOne(t, "SHOW COMMITTED TASKS LIKE 'x%' IN GRAPH g LIMIT 5").(*ast.ShowStmt)
	if s.Scope != ast.ShowScopeOther || s.ScopeText != "GRAPH" || s.ScopeName.String() != "g" {
		t.Errorf("got %+v", s)
	}
	// IN GRAPH alone names the schema GRAPH.
	s = mustParseOne(t, "SHOW COMMITTED TASKS IN GRAPH").(*ast.ShowStmt)
	if s.Scope != ast.ShowScopeSchema || s.ScopeName.String() != "GRAPH" {
		t.Errorf("got %+v", s)
	}
	mustRejectObject(t, "SHOW COMMITTED TASKS IN GRAPH db.g")
}
