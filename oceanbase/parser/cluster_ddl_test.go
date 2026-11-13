package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// TestClusterDDL pins the action and the name an object statement parses
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
		{"CREATE TABLEGROUP IF NOT EXISTS tg1 SHARDING = 'ADAPTIVE' PARTITION BY RANGE COLUMNS 1 PARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN MAXVALUE)", "CREATE TABLEGROUP", "tg1"},
		{"ALTER TABLEGROUP tg1 ADD TABLE t1, d.t2", "ALTER TABLEGROUP", "tg1"},
		{"ALTER TABLEGROUP tg1 SET SHARDING = 'NONE'", "ALTER TABLEGROUP", "tg1"},
		{"ALTER TABLEGROUP tg1 DROP PARTITION (p0, p1)", "ALTER TABLEGROUP", "tg1"},
		{"DROP TABLEGROUP IF EXISTS tg1", "DROP TABLEGROUP", "tg1"},
		{"CREATE SEQUENCE IF NOT EXISTS d.s1 START WITH 1 INCREMENT BY -1 NOMAXVALUE CACHE 20 NOCYCLE", "CREATE SEQUENCE", "s1"},
		{"ALTER SEQUENCE s1 RESTART MAXVALUE 100", "ALTER SEQUENCE", "s1"},
		{"DROP SEQUENCE IF EXISTS d.s1", "DROP SEQUENCE", "s1"},
		// A number too long for an INTNUM is a DECIMAL_VAL, which simple_num takes.
		{"CREATE SEQUENCE s1 MINVALUE -9223372036854775808 MAXVALUE 9223372036854775807", "CREATE SEQUENCE", "s1"},
		{"CREATE SEQUENCE s1 MAXVALUE 18446744073709551616", "CREATE SEQUENCE", "s1"},
		// name_list opens with a keyword or a NAME_OB, and goes on with NAME_OBs.
		{"ALTER TABLEGROUP tg1 DROP PARTITION ACCESS_INFO, p1", "ALTER TABLEGROUP", "tg1"},
		{"ALTER TABLEGROUP tg1 TRUNCATE PARTITION ACCESS_INFO p1, p2", "ALTER TABLEGROUP", "tg1"},
		{"ALTER TABLEGROUP tg1 DROP PARTITION p0 p1", "ALTER TABLEGROUP", "tg1"},
		// tablegroup_option_list: the options with a comma, then a run without one.
		{"CREATE TABLEGROUP tg1 BINDING TRUE, SHARDING 'NONE' MAX_USED_PART_ID 1", "CREATE TABLEGROUP", "tg1"},
		{"CREATE TABLEGROUP tg1 BINDING TRUE, SHARDING 'NONE', MAX_USED_PART_ID 1", "CREATE TABLEGROUP", "tg1"},
		{"CREATE TABLEGROUP tg1 BINDING TRUE SHARDING 'NONE' MAX_USED_PART_ID 1", "CREATE TABLEGROUP", "tg1"},
		// A variable's value is an expression before it is a bare BINARY.
		{"CREATE TENANT t1 SET x = BINARY 'a', y = BINARY, z = ON", "CREATE TENANT", "t1"},
		{"ALTER TENANT t1 VARIABLES x = BINARY 'a'", "ALTER TENANT", "t1"},
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
		"CREATE TENANT t2 FROM t1 WITH RESOURCE_POOL p1", // a clone takes two options
		"CREATE TENANT t2 FROM t1 WITH RESOURCE_POOL p1, UNIT u1, UNIT u2", // and no more
		"DROP RESOURCE GROUP g1", // MySQL's, which OceanBase has not
		"CREATE TABLEGROUP tg1 PARTITION BY RANGE (c) (PARTITION p0 VALUES LESS THAN (10))", // a tablegroup names no column
		"CREATE SEQUENCE s1 START 1",
		"ALTER TABLEGROUP tg1",
		"CREATE TABLEGROUP tg1 BINDING TRUE SHARDING 'NONE', MAX_USED_PART_ID 1", // no comma after a run without one
		"CREATE TABLEGROUP tg1 BINDING TRUE,",
		"ALTER TABLEGROUP tg1 DROP PARTITION p0, ACCESS_INFO", // only the first name may be a keyword
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) was accepted", sql)
		}
	}
}
