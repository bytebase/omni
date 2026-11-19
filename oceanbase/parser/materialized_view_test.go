package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// TestMaterializedView pins where the materialized view statements' clauses
// land. Whether each form parses is held to OceanBase by the create_mview_stmt,
// create_mlog_stmt, alter_mview_action, alter_mlog_action and drop_mlog_stmt
// rows of testdata/compat.tsv.
func TestMaterializedView(t *testing.T) {
	one := func(sql string) nodes.Node {
		t.Helper()
		list, err := Parse(sql)
		if err != nil || len(list.Items) != 1 {
			t.Fatalf("Parse(%q) = %v, %v", sql, list, err)
		}
		return list.Items[0]
	}

	log := one("CREATE OR REPLACE MATERIALIZED VIEW LOG ON t PARALLEL 2 WITH ROWID, PRIMARY KEY (a, b) INCLUDING NEW VALUES PURGE START WITH 1 NEXT 2").(*nodes.CreateMaterializedViewLogStmt)
	if !log.OrReplace || log.Table.Name != "t" || len(log.Options) != 1 || log.Options[0].Degree != 2 {
		t.Fatalf("CREATE MATERIALIZED VIEW LOG = %+v", log)
	}
	if log.With == nil || len(log.With.SpecialColumns) != 2 || log.With.SpecialColumns[0] != "ROWID" || log.With.SpecialColumns[1] != "PRIMARY KEY" ||
		!log.With.HasReferenceList || len(log.With.ReferenceColumns) != 2 {
		t.Errorf("WITH = %+v", log.With)
	}
	if log.NewValues != "INCLUDING" || log.Purge == nil || log.Purge.StartWith == nil || log.Purge.Next == nil {
		t.Errorf("new values %q, purge %+v", log.NewValues, log.Purge)
	}
	empty := one("CREATE MATERIALIZED VIEW LOG ON t WITH").(*nodes.CreateMaterializedViewLogStmt)
	if empty.With == nil || len(empty.With.SpecialColumns) != 0 || empty.With.HasReferenceList {
		t.Errorf("bare WITH = %+v", empty.With)
	}

	mv := one("CREATE MATERIALIZED VIEW mv (a, b, PRIMARY KEY USING BTREE (a) COMMENT 'p') ENGINE = InnoDB WITH COLUMN GROUP(ALL COLUMNS) REFRESH FAST INDIVIDUAL ON COMMIT START WITH 1 NEXT 2 DISABLE ON QUERY COMPUTATION ENABLE QUERY REWRITE AS SELECT a FROM t").(*nodes.CreateMaterializedViewStmt)
	if mv.Name.Name != "mv" || len(mv.Columns) != 2 || mv.PrimaryKey == nil || mv.PrimaryKey.IndexType != "BTREE" || mv.PrimaryKey.Comment != "p" {
		t.Fatalf("CREATE MATERIALIZED VIEW = %+v", mv)
	}
	if len(mv.Options) != 1 || mv.Options[0].Name != "ENGINE" || mv.ColumnGroup != "ALL COLUMNS" {
		t.Errorf("options %+v, column group %q", mv.Options, mv.ColumnGroup)
	}
	if o := mv.MviewOptions; o == nil || o.Refresh == nil || o.Refresh.Method != "FAST" || o.Refresh.NestedRefresh != "INDIVIDUAL" ||
		o.Refresh.On != "COMMIT" || o.Refresh.StartWith == nil || o.Refresh.Next == nil ||
		o.OnQueryComputation != "DISABLE" || o.QueryRewrite != "ENABLE" {
		t.Errorf("mview options = %+v", o)
	}
	if mv.Select == nil || mv.SelectText != "SELECT a FROM t" {
		t.Errorf("select %+v, text %q", mv.Select, mv.SelectText)
	}

	alter := one("ALTER MATERIALIZED VIEW mv PARALLEL 4, NOPARALLEL, ENABLE QUERY REWRITE, REFRESH FAST, NEVER REFRESH, REFRESH START WITH 1 NEXT 2, REFRESH CONSISTENT").(*nodes.AlterMaterializedViewStmt)
	wantTypes := []nodes.AlterMviewActionType{
		nodes.AMViewParallel, nodes.AMViewParallel, nodes.AMViewEnableQueryRewrite,
		nodes.AMViewRefreshMethod, nodes.AMViewNeverRefresh, nodes.AMViewRefreshInterval, nodes.AMViewNestedRefresh,
	}
	if len(alter.Actions) != len(wantTypes) {
		t.Fatalf("ALTER MATERIALIZED VIEW has %d actions", len(alter.Actions))
	}
	for i, want := range wantTypes {
		if alter.Actions[i].Type != want {
			t.Errorf("action %d type = %d, want %d", i, alter.Actions[i].Type, want)
		}
	}
	if alter.Actions[0].Parallel.NoParallel || alter.Actions[0].Parallel.Degree != 4 || !alter.Actions[1].Parallel.NoParallel {
		t.Errorf("parallel actions = %+v, %+v", alter.Actions[0].Parallel, alter.Actions[1].Parallel)
	}
	if !alter.Actions[2].Enable || alter.Actions[3].Refresh.Method != "FAST" || !alter.Actions[4].Refresh.Never ||
		alter.Actions[5].Refresh.StartWith == nil || alter.Actions[5].Refresh.Next == nil || alter.Actions[6].Refresh.NestedRefresh != "CONSISTENT" {
		t.Errorf("ALTER MATERIALIZED VIEW actions = %+v", alter.Actions)
	}

	alterLog := one("ALTER MATERIALIZED VIEW LOG ON t NOPARALLEL, PURGE START WITH 1 NEXT 2, LOB_INROW_THRESHOLD = 3").(*nodes.AlterMaterializedViewLogStmt)
	if len(alterLog.Actions) != 3 || alterLog.Actions[0].Type != nodes.AMLogParallel || !alterLog.Actions[0].Parallel.NoParallel ||
		alterLog.Actions[1].Type != nodes.AMLogPurge || alterLog.Actions[1].StartWith == nil || alterLog.Actions[1].Next == nil ||
		alterLog.Actions[2].Type != nodes.AMLogLobInRowThreshold || alterLog.Actions[2].Threshold != 3 {
		t.Errorf("ALTER MATERIALIZED VIEW LOG = %+v", alterLog)
	}

	drop := one("DROP MATERIALIZED VIEW IF EXISTS mv").(*nodes.DropViewStmt)
	if !drop.Materialized || !drop.IfExists || len(drop.Views) != 1 {
		t.Errorf("DROP MATERIALIZED VIEW = %+v", drop)
	}
	dropLog := one("DROP MATERIALIZED VIEW LOG ON t").(*nodes.DropMaterializedViewLogStmt)
	if dropLog.Table.Name != "t" {
		t.Errorf("DROP MATERIALIZED VIEW LOG = %+v", dropLog)
	}

	for _, sql := range []string{
		// opt_check_option requires CHECK and OPTION after WITH.
		"CREATE MATERIALIZED VIEW a AS SELECT a WITH CHECK",
		"CREATE MATERIALIZED VIEW a AS SELECT a WITH LOCAL",
		"CREATE MATERIALIZED VIEW a AS SELECT a WITH CHECK a",
		// Neither materialized view drop production takes TEMPORARY.
		"DROP TEMPORARY MATERIALIZED VIEW a",
		"DROP TEMPORARY MATERIALIZED VIEW LOG ON a",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) was accepted", sql)
		}
	}
}

// TestMaterializedViewCompletion pins that the new statements are reachable
// from completion: MATERIALIZED heads a CREATE, ALTER and DROP, and VIEW
// follows it.
func TestMaterializedViewCompletion(t *testing.T) {
	for _, sql := range []string{"CREATE ", "ALTER ", "DROP "} {
		if cs := Collect(sql, len(sql)); !cs.HasToken(kwMATERIALIZED) {
			t.Errorf("Collect(%q) does not offer MATERIALIZED", sql)
		}
	}
	for _, sql := range []string{"CREATE MATERIALIZED ", "ALTER MATERIALIZED ", "DROP MATERIALIZED "} {
		if cs := Collect(sql, len(sql)); !cs.HasToken(kwVIEW) {
			t.Errorf("Collect(%q) does not offer VIEW", sql)
		}
	}
	for _, sql := range []string{"CREATE MATERIALIZED VIEW ", "ALTER MATERIALIZED VIEW ", "DROP MATERIALIZED VIEW "} {
		if cs := Collect(sql, len(sql)); !cs.HasToken(kwLOG) {
			t.Errorf("Collect(%q) does not offer LOG", sql)
		}
	}
	if cs := Collect("DROP MATERIALIZED VIEW ", len("DROP MATERIALIZED VIEW ")); !cs.HasRule("view_ref") {
		t.Error("Collect(\"DROP MATERIALIZED VIEW \") does not offer view_ref")
	}
	if cs := Collect("CREATE OR REPLACE ", len("CREATE OR REPLACE ")); !cs.HasToken(kwMATERIALIZED) || !cs.HasToken(kwVIEW) {
		t.Error("Collect(\"CREATE OR REPLACE \") does not offer MATERIALIZED and VIEW")
	}
	for _, sql := range []string{
		"CREATE MATERIALIZED VIEW LOG ",
		"ALTER MATERIALIZED VIEW LOG ",
		"DROP MATERIALIZED VIEW LOG ",
	} {
		if cs := Collect(sql, len(sql)); !cs.HasToken(kwON) {
			t.Errorf("Collect(%q) does not offer ON", sql)
		}
	}
	if cs := Collect("CREATE MATERIALIZED VIEW mv AS ", len("CREATE MATERIALIZED VIEW mv AS ")); !cs.HasToken(kwSELECT) {
		t.Error("Collect(\"CREATE MATERIALIZED VIEW mv AS \") does not offer SELECT")
	}
	if cs := Collect("CREATE MATERIALIZED VIEW mv ", len("CREATE MATERIALIZED VIEW mv ")); !cs.HasToken(kwAS) {
		t.Error("Collect(\"CREATE MATERIALIZED VIEW mv \") does not offer AS")
	}
	for _, sql := range []string{
		"CREATE MATERIALIZED VIEW mv REFRESH FAST ",
		"CREATE MATERIALIZED VIEW mv WITH COLUMN GROUP(ALL COLUMNS) ",
	} {
		if cs := Collect(sql, len(sql)); !cs.HasToken(kwAS) {
			t.Errorf("Collect(%q) does not offer AS", sql)
		}
	}
	if cs := Collect("ALTER MATERIALIZED VIEW ", len("ALTER MATERIALIZED VIEW ")); !cs.HasRule("view_ref") {
		t.Error("Collect(\"ALTER MATERIALIZED VIEW \") does not offer view_ref")
	}
	for _, sql := range []string{
		"CREATE MATERIALIZED VIEW LOG ON ",
		"ALTER MATERIALIZED VIEW LOG ON ",
		"DROP MATERIALIZED VIEW LOG ON ",
	} {
		if cs := Collect(sql, len(sql)); !cs.HasRule("table_ref") {
			t.Errorf("Collect(%q) does not offer table_ref", sql)
		}
	}
	if cs := Collect("ALTER MATERIALIZED VIEW mv ", len("ALTER MATERIALIZED VIEW mv ")); !cs.HasToken(kwENABLE) || !cs.HasToken(kwPARALLEL) {
		t.Error("ALTER MATERIALIZED VIEW offers no action candidates")
	}
	if cs := Collect("ALTER MATERIALIZED VIEW LOG ON t ", len("ALTER MATERIALIZED VIEW LOG ON t ")); !cs.HasToken(kwPURGE) || !cs.HasToken(kwPARALLEL) {
		t.Error("ALTER MATERIALIZED VIEW LOG offers no action candidates")
	}
	if cs := Collect("CREATE MATERIALIZED VIEW mv ", len("CREATE MATERIALIZED VIEW mv ")); !cs.HasToken(kwWITH) || !cs.HasToken(kwENABLE) {
		t.Error("CREATE MATERIALIZED VIEW offers no option candidates")
	}
	if cs := Collect("CREATE MATERIALIZED VIEW LOG ON t ", len("CREATE MATERIALIZED VIEW LOG ON t ")); !cs.HasToken(kwPARALLEL) || !cs.HasToken(kwWITH) || !cs.HasToken(kwPURGE) {
		t.Error("CREATE MATERIALIZED VIEW LOG offers no clause candidates")
	}
	if cs := Collect("CREATE MATERIALIZED VIEW ", len("CREATE MATERIALIZED VIEW ")); !cs.HasRule("view_ref") {
		t.Error("Collect(\"CREATE MATERIALIZED VIEW \") does not offer view_ref")
	}
	if cs := Collect("ALTER MATERIALIZED VIEW mv NOPARALLEL, ", len("ALTER MATERIALIZED VIEW mv NOPARALLEL, ")); !cs.HasToken(kwPARALLEL) {
		t.Error("ALTER MATERIALIZED VIEW after a comma offers no action candidates")
	}
	if cs := Collect("ALTER MATERIALIZED VIEW LOG ON t NOPARALLEL, ", len("ALTER MATERIALIZED VIEW LOG ON t NOPARALLEL, ")); !cs.HasToken(kwPURGE) {
		t.Error("ALTER MATERIALIZED VIEW LOG after a comma offers no action candidates")
	}
}
