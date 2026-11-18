package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftCreateDatabaseFromIntegration(t *testing.T) {
	sql := "CREATE DATABASE a FROM INTEGRATION 'id' DATABASE 'db' SET ACCEPTINVCHARS = TRUE QUERY_ALL_STATES FALSE REFRESH_INTERVAL 30 TRUNCATECOLUMNS TRUE HISTORY_MODE FALSE AUTO_REMEDIATION = TRUE"
	stmt := singleStmt(t, sql)
	create, ok := stmt.(*nodes.CreatedbStmt)
	if !ok {
		t.Fatalf("expected CreatedbStmt, got %T", stmt)
	}
	if create.Dbname != "a" {
		t.Fatalf("expected Dbname=a, got %q", create.Dbname)
	}
	if create.Loc.Start != 0 || create.Loc.End != len(sql) {
		t.Fatalf("expected Loc=[0,%d), got %v", len(sql), create.Loc)
	}
	assertDefElemString(t, create.Options, "from_integration", "id")
	assertDefElemString(t, create.Options, "integration_database", "db")
	assertDefElemBool(t, create.Options, "set", true)
	assertDefElemBool(t, create.Options, "acceptinvchars", true)
	assertDefElemBool(t, create.Options, "query_all_states", false)
	assertDefElemInt(t, create.Options, "refresh_interval", 30)
	assertDefElemBool(t, create.Options, "truncatecolumns", true)
	assertDefElemBool(t, create.Options, "history_mode", false)
	assertDefElemBool(t, create.Options, "auto_remediation", true)
}

func TestRedshiftCreateDatabaseFromARN(t *testing.T) {
	t.Run("schema", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE DATABASE a FROM ARN 'arn' WITH DATA CATALOG SCHEMA 'cat'")
		create := stmt.(*nodes.CreatedbStmt)
		assertDefElemString(t, create.Options, "from_arn", "arn")
		assertDefElemString(t, create.Options, "data_catalog_schema", "cat")
	})
	t.Run("no schema", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE DATABASE a FROM ARN 'arn' WITH NO DATA CATALOG SCHEMA")
		create := stmt.(*nodes.CreatedbStmt)
		assertDefElemString(t, create.Options, "from_arn", "arn")
		assertDefElemBool(t, create.Options, "no_data_catalog_schema", true)
	})
}

func TestRedshiftCreateDatabaseFromDatashare(t *testing.T) {
	t.Run("namespace", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE DATABASE a FROM DATASHARE share OF NAMESPACE ns")
		create := stmt.(*nodes.CreatedbStmt)
		assertDefElemString(t, create.Options, "from_datashare", "share")
		assertDefElemString(t, create.Options, "datashare_namespace", "ns")
		if findDefElem(create.Options, "datashare_account") != nil {
			t.Fatal("expected no datashare_account")
		}
	})
	t.Run("with permissions and account", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE DATABASE a WITH PERMISSIONS FROM DATASHARE share OF ACCOUNT 123 NAMESPACE 'ns'")
		create := stmt.(*nodes.CreatedbStmt)
		assertDefElemBool(t, create.Options, "with_permissions", true)
		assertDefElemString(t, create.Options, "from_datashare", "share")
		assertDefElemString(t, create.Options, "datashare_account", "123")
		assertDefElemString(t, create.Options, "datashare_namespace", "ns")
	})
}

func TestRedshiftAlterDatabaseIntegrationRefresh(t *testing.T) {
	t.Run("all tables in schema", func(t *testing.T) {
		stmt := singleStmt(t, "ALTER DATABASE a INTEGRATION REFRESH ALL TABLES IN SCHEMA s1, s2")
		alter, ok := stmt.(*nodes.AlterDatabaseStmt)
		if !ok {
			t.Fatalf("expected AlterDatabaseStmt, got %T", stmt)
		}
		assertDefElemString(t, alter.Options, "refresh", "all")
		names := defElemNames(t, alter.Options, "refresh_schemas")
		if len(names) != 2 || names[0] != "s1" || names[1] != "s2" {
			t.Fatalf("expected refresh_schemas [s1 s2], got %v", names)
		}
	})
	t.Run("tables", func(t *testing.T) {
		stmt := singleStmt(t, "ALTER DATABASE a INTEGRATION REFRESH TABLE t1, t2")
		alter := stmt.(*nodes.AlterDatabaseStmt)
		assertDefElemString(t, alter.Options, "refresh", "table")
		names := defElemNames(t, alter.Options, "refresh_tables")
		if len(names) != 2 || names[0] != "t1" || names[1] != "t2" {
			t.Fatalf("expected refresh_tables [t1 t2], got %v", names)
		}
	})
}

func TestRedshiftAlterDatabaseIntegrationSetHistoryMode(t *testing.T) {
	t.Run("all tables in schema", func(t *testing.T) {
		stmt := singleStmt(t, "ALTER DATABASE a INTEGRATION SET HISTORY_MODE TRUE FOR ALL TABLES IN SCHEMA s")
		alter := stmt.(*nodes.AlterDatabaseStmt)
		assertDefElemBool(t, alter.Options, "history_mode", true)
		assertDefElemBool(t, alter.Options, "history_mode_all_tables", true)
		names := defElemNames(t, alter.Options, "history_mode_schemas")
		if len(names) != 1 || names[0] != "s" {
			t.Fatalf("expected history_mode_schemas [s], got %v", names)
		}
	})
	t.Run("tables", func(t *testing.T) {
		stmt := singleStmt(t, "ALTER DATABASE a INTEGRATION SET QUERY_ALL_STATES = FALSE HISTORY_MODE FALSE FOR TABLE t1, t2 AUTO_REMEDIATION TRUE")
		alter := stmt.(*nodes.AlterDatabaseStmt)
		assertDefElemBool(t, alter.Options, "query_all_states", false)
		assertDefElemBool(t, alter.Options, "history_mode", false)
		assertDefElemBool(t, alter.Options, "auto_remediation", true)
		names := defElemNames(t, alter.Options, "history_mode_tables")
		if len(names) != 2 || names[0] != "t1" || names[1] != "t2" {
			t.Fatalf("expected history_mode_tables [t1 t2], got %v", names)
		}
	})
}

func TestRedshiftCopyDocOptionTails(t *testing.T) {
	t.Run("column placeholder", func(t *testing.T) {
		stmt := firstCopyStmt(t, "COPY a zzz FROM 's' IAM_ROLE DEFAULT")
		if stmt.Attlist == nil || len(stmt.Attlist.Items) != 1 {
			t.Fatalf("expected one column placeholder, got %#v", stmt.Attlist)
		}
		if got := stmt.Attlist.Items[0].(*nodes.String).Str; got != "zzz" {
			t.Fatalf("expected placeholder zzz, got %q", got)
		}
	})
	t.Run("bare identifier options", func(t *testing.T) {
		stmt := firstCopyStmt(t, "COPY a FROM 's' IAM_ROLE DEFAULT a, zzz a")
		if got := countDefElems(stmt.Options, "a"); got != 1 {
			t.Fatalf("expected one a option, got %d", got)
		}
		assertDefElemString(t, stmt.Options, "zzz", "a")
	})
	t.Run("csv format", func(t *testing.T) {
		stmt := firstCopyStmt(t, "COPY a FROM 's' IAM_ROLE DEFAULT CSV")
		assertDefElemString(t, stmt.Options, "format", "csv")
	})
}

func TestRedshiftUnloadDocOptionTails(t *testing.T) {
	t.Run("comma separated formats", func(t *testing.T) {
		stmt := firstUnloadStmt(t, "UNLOAD('s') TO 's' IAM_ROLE DEFAULT CSV, JSON")
		if got := countDefElems(stmt.Options, "format"); got != 2 {
			t.Fatalf("expected two format options, got %d", got)
		}
	})
	t.Run("parallel on", func(t *testing.T) {
		stmt := firstUnloadStmt(t, "UNLOAD('s') TO 's' IAM_ROLE DEFAULT PARALLEL")
		assertDefElemString(t, stmt.Options, "parallel", "on")
	})
}

func defElemNames(t *testing.T, list *nodes.List, name string) []string {
	t.Helper()
	elem := findDefElem(list, name)
	if elem == nil {
		t.Fatalf("expected option %q in %#v", name, list)
	}
	names, ok := elem.Arg.(*nodes.List)
	if !ok {
		t.Fatalf("expected option %q list arg, got %T", name, elem.Arg)
	}
	result := make([]string, 0, len(names.Items))
	for _, item := range names.Items {
		result = append(result, item.(*nodes.String).Str)
	}
	return result
}

func countDefElems(list *nodes.List, name string) int {
	if list == nil {
		return 0
	}
	count := 0
	for _, item := range list.Items {
		if elem, ok := item.(*nodes.DefElem); ok && elem.Defname == name {
			count++
		}
	}
	return count
}
