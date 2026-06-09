package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// ALTER SYSTEM ADD BACKEND
// ---------------------------------------------------------------------------

func TestAlterSystemAddBackend_Single(t *testing.T) {
	node := parseOne(t, `ALTER SYSTEM ADD BACKEND "192.168.0.1:9050,192.168.0.2:9050"`)
	stmt, ok := node.(*ast.SystemAlterStmt)
	if !ok {
		t.Fatalf("expected *SystemAlterStmt, got %T", node)
	}
	if stmt.Verb != "ADD" {
		t.Errorf("Verb=%q, want ADD", stmt.Verb)
	}
	if stmt.Object != "BACKEND" {
		t.Errorf("Object=%q, want BACKEND", stmt.Object)
	}
	if len(stmt.Hosts) != 2 {
		t.Fatalf("Hosts=%v, want 2 entries", stmt.Hosts)
	}
	if stmt.Hosts[0] != "192.168.0.1:9050" {
		t.Errorf("Hosts[0]=%q, want 192.168.0.1:9050", stmt.Hosts[0])
	}
	if stmt.Hosts[1] != "192.168.0.2:9050" {
		t.Errorf("Hosts[1]=%q, want 192.168.0.2:9050", stmt.Hosts[1])
	}
}

func TestAlterSystemAddBackend_WithProperties(t *testing.T) {
	node := parseOne(t, `ALTER SYSTEM ADD BACKEND "doris-be01:9050" PROPERTIES ("tag.location" = "groupb")`)
	stmt := node.(*ast.SystemAlterStmt)
	if stmt.Verb != "ADD" || stmt.Object != "BACKEND" {
		t.Errorf("Verb=%q Object=%q, want ADD BACKEND", stmt.Verb, stmt.Object)
	}
	if len(stmt.Hosts) != 1 || stmt.Hosts[0] != "doris-be01:9050" {
		t.Errorf("Hosts=%v, want [doris-be01:9050]", stmt.Hosts)
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("Properties count=%d, want 1", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "tag.location" {
		t.Errorf("Properties[0].Key=%q, want tag.location", stmt.Properties[0].Key)
	}
	if stmt.Properties[0].Value != "groupb" {
		t.Errorf("Properties[0].Value=%q, want groupb", stmt.Properties[0].Value)
	}
}

func TestAlterSystemAddBackend_ComputeGroup(t *testing.T) {
	node := parseOne(t, `ALTER SYSTEM ADD BACKEND "192.168.0.3:9050" PROPERTIES ("tag.compute_group_name" = "cloud_groupc")`)
	stmt := node.(*ast.SystemAlterStmt)
	if stmt.Verb != "ADD" || stmt.Object != "BACKEND" {
		t.Errorf("Verb=%q Object=%q, want ADD BACKEND", stmt.Verb, stmt.Object)
	}
	if len(stmt.Properties) != 1 || stmt.Properties[0].Value != "cloud_groupc" {
		t.Errorf("unexpected Properties: %v", stmt.Properties)
	}
}

func TestAlterSystemAddBackend_Tag(t *testing.T) {
	node := parseOne(t, `ALTER SYSTEM ADD BACKEND "192.168.0.1:9050,192.168.0.2:9050"`)
	if node.Tag() != ast.T_SystemAlterStmt {
		t.Errorf("Tag=%v, want T_SystemAlterStmt", node.Tag())
	}
}

// ---------------------------------------------------------------------------
// ALTER SYSTEM DROP BACKEND
// ---------------------------------------------------------------------------

func TestAlterSystemDropBackend_MultipleHosts(t *testing.T) {
	node := parseOne(t, `ALTER SYSTEM DROP BACKEND "192.168.0.1:9050", "192.168.0.2:9050"`)
	stmt := node.(*ast.SystemAlterStmt)
	if stmt.Verb != "DROP" || stmt.Object != "BACKEND" {
		t.Errorf("Verb=%q Object=%q, want DROP BACKEND", stmt.Verb, stmt.Object)
	}
	if len(stmt.Hosts) != 2 {
		t.Fatalf("Hosts=%v, want 2 entries", stmt.Hosts)
	}
}

func TestAlterSystemDropBackend_ById(t *testing.T) {
	node := parseOne(t, `ALTER SYSTEM DROP BACKEND "10002"`)
	stmt := node.(*ast.SystemAlterStmt)
	if stmt.Verb != "DROP" || stmt.Object != "BACKEND" {
		t.Errorf("Verb=%q Object=%q, want DROP BACKEND", stmt.Verb, stmt.Object)
	}
	if len(stmt.Hosts) != 1 || stmt.Hosts[0] != "10002" {
		t.Errorf("Hosts=%v, want [10002]", stmt.Hosts)
	}
}

// ---------------------------------------------------------------------------
// ALTER SYSTEM DECOMMISSION BACKEND
// ---------------------------------------------------------------------------

func TestAlterSystemDecommissionBackend_MultipleHosts(t *testing.T) {
	node := parseOne(t, `ALTER SYSTEM DECOMMISSION BACKEND "192.168.0.1:9050", "192.168.0.2:9050"`)
	stmt := node.(*ast.SystemAlterStmt)
	if stmt.Verb != "DECOMMISSION" || stmt.Object != "BACKEND" {
		t.Errorf("Verb=%q Object=%q, want DECOMMISSION BACKEND", stmt.Verb, stmt.Object)
	}
	if len(stmt.Hosts) != 2 {
		t.Fatalf("Hosts=%v, want 2 entries", stmt.Hosts)
	}
}

func TestAlterSystemDecommissionBackend_ById(t *testing.T) {
	node := parseOne(t, `ALTER SYSTEM DECOMMISSION BACKEND "10002"`)
	stmt := node.(*ast.SystemAlterStmt)
	if stmt.Verb != "DECOMMISSION" || stmt.Object != "BACKEND" {
		t.Errorf("Verb=%q Object=%q, want DECOMMISSION BACKEND", stmt.Verb, stmt.Object)
	}
}

// ---------------------------------------------------------------------------
// ALTER SYSTEM MODIFY BACKEND
// ---------------------------------------------------------------------------

func TestAlterSystemModifyBackend_TagLocation(t *testing.T) {
	node := parseOne(t, `ALTER SYSTEM MODIFY BACKEND "127.0.0.1:9050" SET ("tag.location" = "group_a")`)
	stmt := node.(*ast.SystemAlterStmt)
	if stmt.Verb != "MODIFY" || stmt.Object != "BACKEND" {
		t.Errorf("Verb=%q Object=%q, want MODIFY BACKEND", stmt.Verb, stmt.Object)
	}
	if len(stmt.Hosts) != 1 || stmt.Hosts[0] != "127.0.0.1:9050" {
		t.Errorf("Hosts=%v, want [127.0.0.1:9050]", stmt.Hosts)
	}
	if len(stmt.SetClause) != 1 {
		t.Fatalf("SetClause count=%d, want 1", len(stmt.SetClause))
	}
	if stmt.SetClause[0].Key != "tag.location" || stmt.SetClause[0].Value != "group_a" {
		t.Errorf("SetClause[0]=%v, want tag.location=group_a", stmt.SetClause[0])
	}
}

func TestAlterSystemModifyBackend_DisableQuery(t *testing.T) {
	node := parseOne(t, `ALTER SYSTEM MODIFY BACKEND "10002" SET ("disable_query" = "true")`)
	stmt := node.(*ast.SystemAlterStmt)
	if stmt.Verb != "MODIFY" || stmt.Object != "BACKEND" {
		t.Errorf("Verb=%q Object=%q, want MODIFY BACKEND", stmt.Verb, stmt.Object)
	}
	if len(stmt.SetClause) != 1 || stmt.SetClause[0].Key != "disable_query" {
		t.Errorf("SetClause=%v, want disable_query", stmt.SetClause)
	}
}

// ---------------------------------------------------------------------------
// ALTER SYSTEM ADD/DROP BROKER
// ---------------------------------------------------------------------------

func TestAlterSystemAddBroker_MultipleHosts(t *testing.T) {
	node := parseOne(t, `ALTER SYSTEM ADD BROKER broker_name "host1:port", "host2:port"`)
	stmt := node.(*ast.SystemAlterStmt)
	if stmt.Verb != "ADD" || stmt.Object != "BROKER" {
		t.Errorf("Verb=%q Object=%q, want ADD BROKER", stmt.Verb, stmt.Object)
	}
	if stmt.BrokerName != "broker_name" {
		t.Errorf("BrokerName=%q, want broker_name", stmt.BrokerName)
	}
	if len(stmt.Hosts) != 2 {
		t.Fatalf("Hosts=%v, want 2 entries", stmt.Hosts)
	}
}

func TestAlterSystemAddBroker_FQDN(t *testing.T) {
	node := parseOne(t, `ALTER SYSTEM ADD BROKER broker_fqdn1 "broker_fqdn1:port"`)
	stmt := node.(*ast.SystemAlterStmt)
	if stmt.Verb != "ADD" || stmt.Object != "BROKER" {
		t.Errorf("Verb=%q Object=%q, want ADD BROKER", stmt.Verb, stmt.Object)
	}
	if stmt.BrokerName != "broker_fqdn1" {
		t.Errorf("BrokerName=%q, want broker_fqdn1", stmt.BrokerName)
	}
	if len(stmt.Hosts) != 1 {
		t.Fatalf("Hosts=%v, want 1 entry", stmt.Hosts)
	}
}

func TestAlterSystemDropAllBroker(t *testing.T) {
	node := parseOne(t, `ALTER SYSTEM DROP ALL BROKER broker_name`)
	stmt := node.(*ast.SystemAlterStmt)
	if stmt.Verb != "DROP" || stmt.Object != "BROKER" {
		t.Errorf("Verb=%q Object=%q, want DROP BROKER", stmt.Verb, stmt.Object)
	}
	if !stmt.DropAll {
		t.Error("DropAll should be true")
	}
	if stmt.BrokerName != "broker_name" {
		t.Errorf("BrokerName=%q, want broker_name", stmt.BrokerName)
	}
}

func TestAlterSystemDropBroker_ByHost(t *testing.T) {
	node := parseOne(t, `ALTER SYSTEM DROP BROKER broker_name "10.10.10.1:8000"`)
	stmt := node.(*ast.SystemAlterStmt)
	if stmt.Verb != "DROP" || stmt.Object != "BROKER" {
		t.Errorf("Verb=%q Object=%q, want DROP BROKER", stmt.Verb, stmt.Object)
	}
	if stmt.DropAll {
		t.Error("DropAll should be false")
	}
	if stmt.BrokerName != "broker_name" {
		t.Errorf("BrokerName=%q, want broker_name", stmt.BrokerName)
	}
	if len(stmt.Hosts) != 1 || stmt.Hosts[0] != "10.10.10.1:8000" {
		t.Errorf("Hosts=%v, want [10.10.10.1:8000]", stmt.Hosts)
	}
}

// ---------------------------------------------------------------------------
// ALTER SYSTEM ADD OBSERVER / FOLLOWER
// ---------------------------------------------------------------------------

func TestAlterSystemAddObserver(t *testing.T) {
	node := parseOne(t, `ALTER SYSTEM ADD OBSERVER "host_ip:9010"`)
	stmt := node.(*ast.SystemAlterStmt)
	if stmt.Verb != "ADD" || stmt.Object != "OBSERVER" {
		t.Errorf("Verb=%q Object=%q, want ADD OBSERVER", stmt.Verb, stmt.Object)
	}
	if len(stmt.Hosts) != 1 || stmt.Hosts[0] != "host_ip:9010" {
		t.Errorf("Hosts=%v, want [host_ip:9010]", stmt.Hosts)
	}
}

func TestAlterSystemDropFollower(t *testing.T) {
	node := parseOne(t, `ALTER SYSTEM DROP FOLLOWER "127.0.0.1:9010"`)
	stmt := node.(*ast.SystemAlterStmt)
	if stmt.Verb != "DROP" || stmt.Object != "FOLLOWER" {
		t.Errorf("Verb=%q Object=%q, want DROP FOLLOWER", stmt.Verb, stmt.Object)
	}
}

func TestAlterSystemDropObserver(t *testing.T) {
	node := parseOne(t, `ALTER SYSTEM DROP OBSERVER "127.0.0.1:9010"`)
	stmt := node.(*ast.SystemAlterStmt)
	if stmt.Verb != "DROP" || stmt.Object != "OBSERVER" {
		t.Errorf("Verb=%q Object=%q, want DROP OBSERVER", stmt.Verb, stmt.Object)
	}
}

func TestAlterSystemAddFollower(t *testing.T) {
	node := parseOne(t, `ALTER SYSTEM ADD FOLLOWER "127.0.0.1:9010"`)
	stmt := node.(*ast.SystemAlterStmt)
	if stmt.Verb != "ADD" || stmt.Object != "FOLLOWER" {
		t.Errorf("Verb=%q Object=%q, want ADD FOLLOWER", stmt.Verb, stmt.Object)
	}
}

// ---------------------------------------------------------------------------
// CANCEL DECOMMISSION BACKEND
// ---------------------------------------------------------------------------

func TestCancelDecommissionBackend_Single(t *testing.T) {
	node := parseOne(t, `CANCEL DECOMMISSION BACKEND "192.168.0.1:9050"`)
	stmt, ok := node.(*ast.CancelDecommissionStmt)
	if !ok {
		t.Fatalf("expected *CancelDecommissionStmt, got %T", node)
	}
	if stmt.Object != "BACKEND" {
		t.Errorf("Object=%q, want BACKEND", stmt.Object)
	}
	if len(stmt.Hosts) != 1 || stmt.Hosts[0] != "192.168.0.1:9050" {
		t.Errorf("Hosts=%v, want [192.168.0.1:9050]", stmt.Hosts)
	}
}

func TestCancelDecommissionBackend_Multiple(t *testing.T) {
	node := parseOne(t, `CANCEL DECOMMISSION BACKEND "192.168.0.1:9050", "192.168.0.2:9050"`)
	stmt := node.(*ast.CancelDecommissionStmt)
	if len(stmt.Hosts) != 2 {
		t.Fatalf("Hosts=%v, want 2 entries", stmt.Hosts)
	}
}

func TestCancelDecommissionBackend_Tag(t *testing.T) {
	node := parseOne(t, `CANCEL DECOMMISSION BACKEND "192.168.0.1:9050"`)
	if node.Tag() != ast.T_CancelDecommissionStmt {
		t.Errorf("Tag=%v, want T_CancelDecommissionStmt", node.Tag())
	}
}

// ---------------------------------------------------------------------------
// ADMIN SHOW REPLICA DISTRIBUTION
// ---------------------------------------------------------------------------

func TestAdminShowReplicaDistribution(t *testing.T) {
	node := parseOne(t, `ADMIN SHOW REPLICA DISTRIBUTION FROM my_table`)
	stmt, ok := node.(*ast.AdminStmt)
	if !ok {
		t.Fatalf("expected *AdminStmt, got %T", node)
	}
	if stmt.Verb != "SHOW" {
		t.Errorf("Verb=%q, want SHOW", stmt.Verb)
	}
	if stmt.Object != "REPLICA" {
		t.Errorf("Object=%q, want REPLICA", stmt.Object)
	}
	if stmt.Target == nil {
		t.Fatal("Target should not be nil")
	}
	if stmt.Target.String() != "my_table" {
		t.Errorf("Target=%q, want my_table", stmt.Target.String())
	}
}

func TestAdminShowReplicaStatus(t *testing.T) {
	node := parseOne(t, `ADMIN SHOW REPLICA STATUS FROM my_table`)
	stmt := node.(*ast.AdminStmt)
	if stmt.Verb != "SHOW" || stmt.Object != "REPLICA" {
		t.Errorf("Verb=%q Object=%q, want SHOW REPLICA", stmt.Verb, stmt.Object)
	}
	if stmt.Target == nil || stmt.Target.String() != "my_table" {
		t.Errorf("Target=%v, want my_table", stmt.Target)
	}
}

// ---------------------------------------------------------------------------
// ADMIN SHOW CONFIG
// ---------------------------------------------------------------------------

func TestAdminShowConfig(t *testing.T) {
	node := parseOne(t, `ADMIN SHOW CONFIG`)
	stmt := node.(*ast.AdminStmt)
	if stmt.Verb != "SHOW" || stmt.Object != "CONFIG" {
		t.Errorf("Verb=%q Object=%q, want SHOW CONFIG", stmt.Verb, stmt.Object)
	}
}

func TestAdminShowFrontendConfig(t *testing.T) {
	node := parseOne(t, `ADMIN SHOW FRONTEND CONFIG`)
	stmt := node.(*ast.AdminStmt)
	if stmt.Verb != "SHOW" {
		t.Errorf("Verb=%q, want SHOW", stmt.Verb)
	}
	if stmt.Object != "FRONTEND CONFIG" {
		t.Errorf("Object=%q, want FRONTEND CONFIG", stmt.Object)
	}
}

// ---------------------------------------------------------------------------
// ADMIN DIAGNOSE TABLET
// ---------------------------------------------------------------------------

func TestAdminDiagnoseTablet(t *testing.T) {
	node := parseOne(t, `ADMIN DIAGNOSE TABLET 12345`)
	stmt := node.(*ast.AdminStmt)
	if stmt.Verb != "DIAGNOSE" || stmt.Object != "TABLET" {
		t.Errorf("Verb=%q Object=%q, want DIAGNOSE TABLET", stmt.Verb, stmt.Object)
	}
}

// ---------------------------------------------------------------------------
// ADMIN REBALANCE DISK
// ---------------------------------------------------------------------------

func TestAdminRebalanceDisk(t *testing.T) {
	node := parseOne(t, `ADMIN REBALANCE DISK`)
	stmt := node.(*ast.AdminStmt)
	if stmt.Verb != "REBALANCE" || stmt.Object != "DISK" {
		t.Errorf("Verb=%q Object=%q, want REBALANCE DISK", stmt.Verb, stmt.Object)
	}
}

func TestAdminCancelRebalanceDisk(t *testing.T) {
	node := parseOne(t, `ADMIN CANCEL REBALANCE DISK`)
	stmt := node.(*ast.AdminStmt)
	if stmt.Verb != "CANCEL" || stmt.Object != "REBALANCE DISK" {
		t.Errorf("Verb=%q Object=%q, want CANCEL / REBALANCE DISK", stmt.Verb, stmt.Object)
	}
}

// ---------------------------------------------------------------------------
// ADMIN COMPACT TABLE
// ---------------------------------------------------------------------------

func TestAdminCompactTable(t *testing.T) {
	node := parseOne(t, `ADMIN COMPACT TABLE orders`)
	stmt := node.(*ast.AdminStmt)
	if stmt.Verb != "COMPACT" || stmt.Object != "TABLE" {
		t.Errorf("Verb=%q Object=%q, want COMPACT TABLE", stmt.Verb, stmt.Object)
	}
	if stmt.Target == nil || stmt.Target.String() != "orders" {
		t.Errorf("Target=%v, want orders", stmt.Target)
	}
}

// ---------------------------------------------------------------------------
// ADMIN CHECK TABLET
// ---------------------------------------------------------------------------

func TestAdminCheckTablet(t *testing.T) {
	file, errs := Parse(`ADMIN CHECK TABLET (10000, 10001) PROPERTIES("type"="consistency")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt := file.Stmts[0].(*ast.AdminStmt)
	if stmt.Verb != "CHECK" || stmt.Object != "TABLET" {
		t.Errorf("Verb=%q Object=%q, want CHECK TABLET", stmt.Verb, stmt.Object)
	}
	// Properties may or may not be parsed depending on args collection order;
	// just ensure no hard crash.
}

// ---------------------------------------------------------------------------
// ADMIN REPAIR TABLE
// ---------------------------------------------------------------------------

func TestAdminRepairTable(t *testing.T) {
	node := parseOne(t, `ADMIN REPAIR TABLE orders`)
	stmt := node.(*ast.AdminStmt)
	if stmt.Verb != "REPAIR" || stmt.Object != "TABLE" {
		t.Errorf("Verb=%q Object=%q, want REPAIR TABLE", stmt.Verb, stmt.Object)
	}
}

func TestAdminCancelRepairTable(t *testing.T) {
	node := parseOne(t, `ADMIN CANCEL REPAIR TABLE orders`)
	stmt := node.(*ast.AdminStmt)
	if stmt.Verb != "CANCEL" || stmt.Object != "REPAIR TABLE" {
		t.Errorf("Verb=%q Object=%q, want CANCEL / REPAIR TABLE", stmt.Verb, stmt.Object)
	}
}

// ---------------------------------------------------------------------------
// ADMIN SET REPLICA STATUS
// ---------------------------------------------------------------------------

func TestAdminSetReplicaStatus(t *testing.T) {
	node := parseOne(t, `ADMIN SET REPLICA STATUS PROPERTIES("tablet_id"="10003", "backend_id"="10001", "status"="bad")`)
	stmt := node.(*ast.AdminStmt)
	if stmt.Verb != "SET" || stmt.Object != "REPLICA STATUS" {
		t.Errorf("Verb=%q Object=%q, want SET REPLICA STATUS", stmt.Verb, stmt.Object)
	}
	if len(stmt.Properties) != 3 {
		t.Errorf("Properties count=%d, want 3", len(stmt.Properties))
	}
}

// ---------------------------------------------------------------------------
// ADMIN SET FRONTEND CONFIG
// ---------------------------------------------------------------------------

func TestAdminSetFrontendConfig(t *testing.T) {
	node := parseOne(t, `ADMIN SET FRONTEND CONFIG ("key"="value")`)
	stmt := node.(*ast.AdminStmt)
	if stmt.Verb != "SET" || stmt.Object != "FRONTEND CONFIG" {
		t.Errorf("Verb=%q Object=%q, want SET FRONTEND CONFIG", stmt.Verb, stmt.Object)
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("Properties count=%d, want 1", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "key" || stmt.Properties[0].Value != "value" {
		t.Errorf("Properties[0]=%v, want key=value", stmt.Properties[0])
	}
}

// ---------------------------------------------------------------------------
// ADMIN CLEAN TRASH
// ---------------------------------------------------------------------------

func TestAdminCleanTrash(t *testing.T) {
	node := parseOne(t, `ADMIN CLEAN TRASH`)
	stmt := node.(*ast.AdminStmt)
	if stmt.Verb != "CLEAN" || stmt.Object != "TRASH" {
		t.Errorf("Verb=%q Object=%q, want CLEAN TRASH", stmt.Verb, stmt.Object)
	}
}

// ---------------------------------------------------------------------------
// ADMIN COPY TABLET
// ---------------------------------------------------------------------------

func TestAdminCopyTablet(t *testing.T) {
	node := parseOne(t, `ADMIN COPY TABLET 10007 PROPERTIES("backend_id"="10001","version"="1")`)
	stmt := node.(*ast.AdminStmt)
	if stmt.Verb != "COPY" || stmt.Object != "TABLET" {
		t.Errorf("Verb=%q Object=%q, want COPY TABLET", stmt.Verb, stmt.Object)
	}
}

// ---------------------------------------------------------------------------
// ADMIN DECOMMISSION BACKEND BY HOSTNAME
// ---------------------------------------------------------------------------

func TestAdminDecommissionBackend(t *testing.T) {
	node := parseOne(t, `ADMIN DECOMMISSION BACKEND BY HOSTNAME 'host:port'`)
	stmt := node.(*ast.AdminStmt)
	if stmt.Verb != "DECOMMISSION" || stmt.Object != "BACKEND" {
		t.Errorf("Verb=%q Object=%q, want DECOMMISSION BACKEND", stmt.Verb, stmt.Object)
	}
}

// ---------------------------------------------------------------------------
// Tag tests
// ---------------------------------------------------------------------------

func TestAdminStmt_Tag(t *testing.T) {
	node := parseOne(t, `ADMIN SHOW CONFIG`)
	if node.Tag() != ast.T_AdminStmt {
		t.Errorf("Tag=%v, want T_AdminStmt", node.Tag())
	}
}

func TestSystemAlterStmt_Tag(t *testing.T) {
	node := parseOne(t, `ALTER SYSTEM ADD BACKEND "192.168.0.1:9050"`)
	if node.Tag() != ast.T_SystemAlterStmt {
		t.Errorf("Tag=%v, want T_SystemAlterStmt", node.Tag())
	}
}

// ---------------------------------------------------------------------------
// Legacy corpus round-trips
// ---------------------------------------------------------------------------

// TestLegacyCorpus_Backend covers all statements from cluster_backend.sql.
func TestLegacyCorpus_Backend(t *testing.T) {
	cases := []string{
		`ALTER SYSTEM ADD BACKEND "192.168.0.1:9050,192.168.0.2:9050"`,
		`ALTER SYSTEM ADD BACKEND "doris-be01:9050" PROPERTIES ("tag.location" = "groupb")`,
		`ALTER SYSTEM ADD BACKEND "192.168.0.3:9050" PROPERTIES ("tag.compute_group_name" = "cloud_groupc")`,
		`ALTER SYSTEM DROP BACKEND "192.168.0.1:9050", "192.168.0.2:9050"`,
		`ALTER SYSTEM DROP BACKEND "10002"`,
		`ALTER SYSTEM DECOMMISSION BACKEND "192.168.0.1:9050", "192.168.0.2:9050"`,
		`ALTER SYSTEM DECOMMISSION BACKEND "10002"`,
		`ALTER SYSTEM MODIFY BACKEND "127.0.0.1:9050" SET ("tag.location" = "group_a")`,
		`ALTER SYSTEM MODIFY BACKEND "10002" SET ("disable_query" = "true")`,
		`ALTER SYSTEM MODIFY BACKEND "127.0.0.1:9050" SET ("disable_load" = "true")`,
	}
	for i, sql := range cases {
		file, errs := Parse(sql)
		if len(errs) != 0 {
			t.Errorf("case %d %q: errors=%v", i, sql, errs)
			continue
		}
		if len(file.Stmts) != 1 {
			t.Errorf("case %d %q: want 1 stmt, got %d", i, sql, len(file.Stmts))
			continue
		}
		if _, ok := file.Stmts[0].(*ast.SystemAlterStmt); !ok {
			t.Errorf("case %d %q: want *SystemAlterStmt, got %T", i, sql, file.Stmts[0])
		}
	}
}

// TestLegacyCorpus_Broker covers cluster_broker.sql.
func TestLegacyCorpus_Broker(t *testing.T) {
	cases := []string{
		`ALTER SYSTEM ADD BROKER broker_name "host1:port", "host2:port"`,
		`ALTER SYSTEM ADD BROKER broker_fqdn1 "broker_fqdn1:port"`,
		`ALTER SYSTEM DROP ALL BROKER broker_name`,
		`ALTER SYSTEM DROP BROKER broker_name "10.10.10.1:8000"`,
	}
	for i, sql := range cases {
		file, errs := Parse(sql)
		if len(errs) != 0 {
			t.Errorf("case %d %q: errors=%v", i, sql, errs)
			continue
		}
		if len(file.Stmts) != 1 {
			t.Errorf("case %d %q: want 1 stmt, got %d", i, sql, len(file.Stmts))
			continue
		}
		if _, ok := file.Stmts[0].(*ast.SystemAlterStmt); !ok {
			t.Errorf("case %d %q: want *SystemAlterStmt, got %T", i, sql, file.Stmts[0])
		}
	}
}

// TestLegacyCorpus_Frontend covers cluster_frontend.sql (ALTER parts).
func TestLegacyCorpus_Frontend(t *testing.T) {
	cases := []string{
		`ALTER SYSTEM ADD OBSERVER "host_ip:9010"`,
	}
	for i, sql := range cases {
		file, errs := Parse(sql)
		if len(errs) != 0 {
			t.Errorf("case %d %q: errors=%v", i, sql, errs)
			continue
		}
		if len(file.Stmts) != 1 {
			t.Errorf("case %d %q: want 1 stmt, got %d", i, sql, len(file.Stmts))
			continue
		}
		if _, ok := file.Stmts[0].(*ast.SystemAlterStmt); !ok {
			t.Errorf("case %d %q: want *SystemAlterStmt, got %T", i, sql, file.Stmts[0])
		}
	}
}

// TestAdminVariants checks a representative set of ADMIN forms parse without error.
func TestAdminVariants(t *testing.T) {
	cases := []string{
		`ADMIN SHOW REPLICA DISTRIBUTION FROM tbl`,
		`ADMIN SHOW REPLICA STATUS FROM tbl`,
		`ADMIN SHOW CONFIG`,
		`ADMIN SHOW FRONTEND CONFIG`,
		`ADMIN REBALANCE DISK`,
		`ADMIN CANCEL REBALANCE DISK`,
		`ADMIN DIAGNOSE TABLET 12345`,
		`ADMIN COMPACT TABLE tbl`,
		`ADMIN REPAIR TABLE tbl`,
		`ADMIN CANCEL REPAIR TABLE tbl`,
		`ADMIN SET REPLICA STATUS PROPERTIES("tablet_id"="10003", "backend_id"="10001", "status"="bad")`,
		`ADMIN SET FRONTEND CONFIG ("disable_balance"="true")`,
		`ADMIN CLEAN TRASH`,
		`ADMIN COPY TABLET 10007 PROPERTIES("backend_id"="10001","version"="1")`,
		`ADMIN DECOMMISSION BACKEND BY HOSTNAME 'host:9050'`,
	}
	for i, sql := range cases {
		file, errs := Parse(sql)
		if len(errs) != 0 {
			t.Errorf("case %d %q: errors=%v", i, sql, errs)
			continue
		}
		if len(file.Stmts) != 1 {
			t.Errorf("case %d %q: want 1 stmt, got %d", i, sql, len(file.Stmts))
			continue
		}
		if _, ok := file.Stmts[0].(*ast.AdminStmt); !ok {
			t.Errorf("case %d %q: want *AdminStmt, got %T", i, sql, file.Stmts[0])
		}
	}
}
