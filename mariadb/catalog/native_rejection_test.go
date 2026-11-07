package catalog

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/mariadb/parser"
	gomysql "github.com/go-sql-driver/mysql"
)

// These inherited MySQL fixtures are negative dialect probes on MariaDB.
// They must remain covered, but cannot serve as valid SHOW CREATE inputs.
func TestMariaDBRejectsMySQLDDL(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()
	if err := ctr.execSQL("CREATE TABLE rejected_check (a INT, CONSTRAINT ck CHECK(a > 0))"); err != nil {
		t.Fatal(err)
	}
	for _, sql := range []string{
		"CREATE TABLE rejected (a INT, CHECK(a > 0) NOT ENFORCED)",
		"CREATE TABLE rejected (a POINT NOT NULL SRID 4326)",
		"CREATE TABLE rejected (a LINESTRING SRID 4326)",
		"CREATE TABLE rejected (a POLYGON NOT NULL SRID 4326)",
		"ALTER TABLE rejected_check ALTER CHECK ck NOT ENFORCED",
		"ALTER TABLE rejected_check ALTER CHECK ck ENFORCED",
	} {
		t.Run(sql, func(t *testing.T) {
			err := ctr.execSQL(sql)
			var native *gomysql.MySQLError
			if !errors.As(err, &native) || native.Number != 1064 {
				t.Fatalf("native syntax rejection: %v", err)
			}
			if _, err := parser.Parse(sql); err == nil {
				t.Fatal("parser accepted native syntax rejection")
			}
		})
	}
}

// Dropping a member of a composite primary key is a semantic DDL error,
// distinct from the MySQL-only syntax rejections above.
func TestMariaDBCompositePrimaryKeyDrop(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()
	c := New()
	for _, sql := range []string{"CREATE DATABASE test", "USE test", "CREATE TABLE rejected_pk (a INT, b INT, PRIMARY KEY (a,b))"} {
		if sql != "CREATE DATABASE test" && sql != "USE test" {
			if err := ctr.execSQL(sql); err != nil {
				t.Fatal(err)
			}
		}
		results, err := c.Exec(sql, nil)
		if err != nil || results[0].Error != nil {
			t.Fatalf("catalog setup: %v %v", err, results)
		}
	}
	sql := "ALTER TABLE rejected_pk DROP COLUMN a"
	var native *gomysql.MySQLError
	if err := ctr.execSQL(sql); !errors.As(err, &native) || native.Number != 1072 {
		t.Fatalf("native error: %v", err)
	}
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatal(err)
	}
	var modeled *Error
	if !errors.As(results[0].Error, &modeled) || modeled.Code != 1072 {
		t.Fatalf("catalog error: %v", results[0].Error)
	}
	if len(c.GetDatabase("test").GetTable("rejected_pk").Columns) != 2 {
		t.Fatal("failed DDL changed the catalog")
	}
}

func mariaOracleErrorCode(err error) uint16 {
	var native *gomysql.MySQLError
	if err == nil {
		return 0
	}
	if errors.As(err, &native) {
		return native.Number
	}
	return 65535
}
