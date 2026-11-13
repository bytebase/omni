package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// TestClusterDDL pins the action and the name a cluster statement parses
// into. Whether each form parses is held to OceanBase by its rows in
// testdata/compat.tsv.
func TestClusterDDL(t *testing.T) {
	for _, tc := range []struct{ sql, action, name string }{
		{"CREATE TENANT IF NOT EXISTS t1 RESOURCE_POOL_LIST = ('p1'), PRIMARY_ZONE = 'z1', LOCALITY = 'F@z1' SET ob_tcp_invited_nodes = '%'", "CREATE TENANT", "t1"},
		{"CREATE TENANT t2 FROM t1 USING SNAPSHOT s1 WITH RESOURCE_POOL = p2, UNIT = u1", "CREATE TENANT", "t2"},
		{"CREATE STANDBY TENANT t1 LOG_RESTORE_SOURCE = 'SERVICE=1.1.1.1:2881 USER=u@t PASSWORD=p' RESOURCE_POOL_LIST = ('p1')", "CREATE STANDBY TENANT", "t1"},
		{"ALTER TENANT t1 PRIMARY_ZONE = RANDOM, COMMENT 'c'", "ALTER TENANT", "t1"},
		{"ALTER TENANT ALL SET VARIABLES a = 1, b TO DEFAULT", "ALTER TENANT", ""},
		{"ALTER TENANT t1 RENAME GLOBAL_NAME TO t2", "ALTER TENANT", "t1"},
		{"ALTER TENANT t1 LOCK", "ALTER TENANT", "t1"},
		{"DROP TENANT IF EXISTS t1 FORCE", "DROP TENANT", "t1"},
		{"CREATE RESOURCE UNIT u1 MAX_CPU = 2, MEMORY_SIZE = '2G'", "CREATE RESOURCE UNIT", "u1"},
		{"CREATE RESOURCE POOL p1 UNIT = 'u1', UNIT_NUM = 1, ZONE_LIST = ('z1', 'z2')", "CREATE RESOURCE POOL", "p1"},
		{"ALTER RESOURCE POOL p1 UNIT_NUM = 2 DELETE UNIT = (1001, 1002)", "ALTER RESOURCE POOL", "p1"},
		{"ALTER RESOURCE POOL MERGE ('p1', 'p2') INTO ('p3')", "ALTER RESOURCE POOL", ""},
		{"ALTER RESOURCE TENANT t1 UNIT_NUM = 2", "ALTER RESOURCE TENANT", "t1"},
		{"DROP RESOURCE UNIT IF EXISTS u1", "DROP RESOURCE UNIT", "u1"},
		// A variable's value is an expression before it is a bare BINARY.
		{"CREATE TENANT t1 SET x = BINARY 'a', y = BINARY, z = ON", "CREATE TENANT", "t1"},
		{"ALTER TENANT t1 VARIABLES x = BINARY 'a'", "ALTER TENANT", "t1"},
		// A keyword var_name leaves out is a variable's name once it is quoted.
		{"CREATE TENANT t1 SET `range` = 1, `maxvalue` = 2", "CREATE TENANT", "t1"},
		// A user variable is one value, quoted or not.
		{"CREATE TENANT t1 PRIMARY_ZONE = @'zone', REPLICA_NUM = 1", "CREATE TENANT", "t1"},
		{"ALTER TENANT t1 PRIMARY_ZONE = @\"zone\"", "ALTER TENANT", "t1"},
		{"ALTER TENANT t1 PRIMARY_ZONE = @zone", "ALTER TENANT", "t1"},
	} {
		list, err := Parse(tc.sql)
		if err != nil {
			t.Errorf("Parse(%q): %v", tc.sql, err)
			continue
		}
		stmt, ok := list.Items[0].(*nodes.ObjectDDLStmt)
		if len(list.Items) != 1 || !ok {
			t.Errorf("Parse(%q) = %T, want one *ast.ObjectDDLStmt", tc.sql, list.Items[0])
			continue
		}
		if stmt.Action != tc.action || stmt.Name != tc.name {
			t.Errorf("Parse(%q): %s %q, want %s %q", tc.sql, stmt.Action, stmt.Name, tc.action, tc.name)
		}
	}
	for _, sql := range []string{
		"ALTER TENANT t1",     // an option or a variable
		"ALTER TENANT t1 SET", // the same
		"CREATE TENANT",
		"CREATE RESOURCE UNIT u1 MAX_CPU",
		"CREATE TENANT t1 SET range = 1", // var_name leaves out MAXVALUE, RANGE, PASSWORD and ACCESS
		"ALTER TENANT t1 VARIABLES maxvalue = 1",
		"ALTER TENANT t1 VARIABLES password = 1",
		"ALTER TENANT t1 VARIABLES access = 1",
		"CREATE TENANT t2 FROM t1 WITH RESOURCE_POOL p1",                   // a clone takes two options
		"CREATE TENANT t2 FROM t1 WITH RESOURCE_POOL p1, UNIT u1, UNIT u2", // and no more
		"DROP RESOURCE GROUP g1",                                           // MySQL's, which OceanBase has not
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) was accepted", sql)
		}
	}
}
