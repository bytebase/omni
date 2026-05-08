package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

var p2ModernFeaturePositiveSQL = []string{
	"CREATE ANALYTIC VIEW sales_av USING sales_fact DIMENSION BY (time_dim) MEASURES (amount)",
	"ALTER ANALYTIC VIEW sales_av RENAME TO sales_av2",
	"DROP ANALYTIC VIEW sales_av",
	"CREATE ATTRIBUTE DIMENSION time_attr_dim USING time_dim ATTRIBUTES (yr) LEVEL detail KEY yr",
	"ALTER ATTRIBUTE DIMENSION time_attr_dim RENAME TO time_attr_dim2",
	"DROP ATTRIBUTE DIMENSION time_attr_dim",
	"CREATE HIERARCHY time_hier USING time_attr_dim (month CHILD OF year)",
	"ALTER HIERARCHY time_hier RENAME TO time_hier2",
	"DROP HIERARCHY time_hier",
	"CREATE DOMAIN email_domain AS VARCHAR2(255) NOT NULL",
	"ALTER DOMAIN email_domain ADD DISPLAY email_col",
	"DROP DOMAIN email_domain",
	"CREATE JSON RELATIONAL DUALITY VIEW emp_dv AS SELECT employee_id FROM employees",
	"ALTER JSON RELATIONAL DUALITY VIEW emp_dv ENABLE LOGICAL REPLICATION",
	"DROP JSON RELATIONAL DUALITY VIEW emp_dv",
	"CREATE DIMENSION time_dim LEVEL day IS (t.day_id)",
	"ALTER DIMENSION time_dim COMPILE",
	"DROP DIMENSION time_dim",
	"CREATE MATERIALIZED ZONEMAP sales_zmap ON sales (region_id)",
	"ALTER MATERIALIZED ZONEMAP sales_zmap ENABLE PRUNING",
	"DROP MATERIALIZED ZONEMAP sales_zmap",
	"CREATE INMEMORY JOIN GROUP my_jg (sales(prod_id), products(prod_id))",
	"ALTER INMEMORY JOIN GROUP my_jg ADD (orders(prod_id))",
	"DROP INMEMORY JOIN GROUP my_jg",
	"CREATE INDEXTYPE my_itype FOR my_op(NUMBER) USING my_type",
	"ALTER INDEXTYPE my_itype ADD my_op2(VARCHAR2)",
	"DROP INDEXTYPE my_itype",
	"CREATE OPERATOR my_eq BINDING (NUMBER, NUMBER) RETURN NUMBER USING my_eq_func",
	"ALTER OPERATOR my_eq ADD BINDING (VARCHAR2, VARCHAR2) RETURN NUMBER USING my_eq_str",
	"DROP OPERATOR my_eq",
	"CREATE LOCKDOWN PROFILE my_profile",
	"ALTER LOCKDOWN PROFILE my_profile DISABLE FEATURE = (NETWORK_ACCESS)",
	"DROP LOCKDOWN PROFILE my_profile",
	"CREATE OUTLINE my_outline ON SELECT * FROM employees",
	"ALTER OUTLINE my_outline RENAME TO new_outline",
	"DROP OUTLINE my_outline",
	"CREATE PROPERTY GRAPH my_graph VERTEX TABLES (persons KEY (person_id)) EDGE TABLES (knows SOURCE persons DESTINATION persons)",
	"CREATE VECTOR INDEX vec_idx ON docs (embedding) ORGANIZATION INMEMORY NEIGHBOR GRAPH DISTANCE COSINE",
	"CREATE JAVA SOURCE NAMED my_java AS 'public class Foo {}'",
	"ALTER JAVA SOURCE my_java COMPILE",
	"DROP JAVA SOURCE my_java",
	"CREATE LIBRARY my_lib AS '/usr/lib/mylib.so'",
	"ALTER LIBRARY my_lib COMPILE",
	"DROP LIBRARY my_lib",
	"CREATE MLE ENV my_env IMPORTS ('module1')",
	"DROP MLE ENV my_env",
	"CREATE MLE MODULE my_module LANGUAGE JAVASCRIPT AS 'export function hello() { return 1; }'",
	"DROP MLE MODULE my_module",
	"CREATE LOGICAL PARTITION TRACKING ON my_table PARTITION BY RANGE (id) (PARTITION p1 VALUES LESS THAN (1000))",
	"CREATE PMEM FILESTORE my_fs MOUNTPOINT '/pmem0' SIZE 100G",
}

func TestP2ModernFeaturePositive(t *testing.T) {
	for _, sql := range p2ModernFeaturePositiveSQL {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() != 1 {
				t.Fatalf("expected 1 statement, got %d", result.Len())
			}
		})
	}
}

func TestP2ModernFeatureNegative(t *testing.T) {
	tests := []string{
		"CREATE ANALYTIC VIEW USING sales_fact",
		"ALTER ATTRIBUTE DIMENSION RENAME TO new_dim",
		"CREATE DOMAIN AS NUMBER",
		"CREATE PROPERTY GRAPH VERTEX TABLES (persons)",
		"CREATE VECTOR INDEX ON docs (embedding)",
		"CREATE LOCKDOWN PROFILE",
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseShouldFail(t, sql)
		})
	}
}

func TestP2ModernFeatureStructuredOptions(t *testing.T) {
	result := ParseAndCheck(t, "CREATE ANALYTIC VIEW sales_av USING sales_fact DIMENSION BY (time_dim) MEASURES (amount) CACHE")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.CreateAnalyticViewStmt)
	if stmt.Options == nil || stmt.Options.Len() != 1 {
		t.Fatalf("expected analytic view tail option, got %+v", stmt.Options)
	}
}

func TestP2ModernFeatureLoc(t *testing.T) {
	for _, sql := range p2ModernFeaturePositiveSQL {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			raw := result.Items[0].(*ast.RawStmt)
			loc := ast.NodeLoc(raw.Stmt)
			if loc.Start != 0 || loc.End <= loc.Start {
				t.Fatalf("invalid statement Loc=%+v", loc)
			}
		})
	}
}
