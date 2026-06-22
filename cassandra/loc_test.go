package cassandra

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bytebase/omni/cassandra/ast"
)

type LocViolation struct {
	Path    string
	Start   int
	End     int
	Reason  string
}

func (v LocViolation) String() string {
	return fmt.Sprintf("%s: %s (Start=%d, End=%d)", v.Path, v.Reason, v.Start, v.End)
}

func CheckLocations(t *testing.T, sql string) []LocViolation {
	t.Helper()
	stmts, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	var violations []LocViolation
	for i, s := range stmts {
		path := fmt.Sprintf("stmts[%d]", i)
		walkNodeLocs(reflect.ValueOf(s.AST), path, &violations)
	}
	return violations
}

var locType = reflect.TypeOf(ast.Loc{})

func walkNodeLocs(v reflect.Value, path string, violations *[]LocViolation) {
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return
		}
		walkNodeLocs(v.Elem(), path, violations)
	case reflect.Interface:
		if v.IsNil() {
			return
		}
		elem := v.Elem()
		typeName := elem.Type().Name()
		if elem.Kind() == reflect.Ptr {
			typeName = elem.Type().Elem().Name()
		}
		walkNodeLocs(elem, path+"("+typeName+")", violations)
	case reflect.Struct:
		t := v.Type()
		locField := v.FieldByName("Loc")
		if locField.IsValid() && locField.Type() == locType {
			loc := locField.Interface().(ast.Loc)
			if loc.Start >= 0 && loc.End >= 0 {
				if loc.End <= loc.Start {
					*violations = append(*violations, LocViolation{
						Path:   path,
						Start:  loc.Start,
						End:    loc.End,
						Reason: "End <= Start",
					})
				}
			} else if (loc.Start < 0) != (loc.End < 0) {
				*violations = append(*violations, LocViolation{
					Path:   path,
					Start:  loc.Start,
					End:    loc.End,
					Reason: "mixed unknown sentinel",
				})
			}
		}
		for i := range t.NumField() {
			field := t.Field(i)
			if !field.IsExported() || field.Name == "Loc" {
				continue
			}
			walkNodeLocs(v.Field(i), path+"."+field.Name, violations)
		}
	case reflect.Slice:
		for i := range v.Len() {
			elem := v.Index(i)
			walkNodeLocs(elem, fmt.Sprintf("%s[%d]", path, i), violations)
		}
	}
}

func TestCheckLocations(t *testing.T) {
	tests := []string{
		// DML
		"SELECT * FROM users",
		"SELECT DISTINCT name, age FROM ks.users WHERE id = 1 LIMIT 10 ALLOW FILTERING",
		"SELECT JSON name AS n FROM users",
		"INSERT INTO users (id, name) VALUES (1, 'Alice')",
		"INSERT INTO ks.users (id, name) VALUES (1, 'Bob') IF NOT EXISTS",
		"INSERT INTO users (id, name) VALUES (1, 'Charlie') USING TTL 86400",
		"INSERT INTO users JSON '{\"id\": 1, \"name\": \"Dave\"}'",
		"UPDATE users SET name = 'Alice' WHERE id = 1",
		"UPDATE ks.users USING TTL 3600 SET name = 'Bob' WHERE id = 2",
		"UPDATE users SET name = 'Charlie' WHERE id = 3 IF EXISTS",
		"UPDATE users SET name = 'Dave' WHERE id = 4 IF name = 'old'",
		"DELETE FROM users WHERE id = 1",
		"DELETE name FROM ks.users WHERE id = 2",
		`BEGIN BATCH
			INSERT INTO users (id, name) VALUES (1, 'Alice');
			UPDATE users SET name = 'Bob' WHERE id = 2;
			DELETE FROM users WHERE id = 3;
		APPLY BATCH`,
		"TRUNCATE users",
		"USE my_keyspace",

		// DDL
		"CREATE KEYSPACE ks WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': '1'}",
		"ALTER KEYSPACE ks WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': '3'}",
		"DROP KEYSPACE IF EXISTS ks",
		"CREATE TABLE users (id int, name text, PRIMARY KEY (id))",
		"CREATE TABLE t (id int, name text, age int, PRIMARY KEY ((id, name), age)) WITH CLUSTERING ORDER BY (age DESC)",
		"ALTER TABLE users ADD email text",
		"DROP TABLE IF EXISTS users",
		"CREATE INDEX ON users (name)",
		"DROP INDEX IF EXISTS users_name_idx",
		"CREATE TYPE address (street text, city text)",
		"ALTER TYPE address ADD zip text",
		"DROP TYPE IF EXISTS address",
		"CREATE MATERIALIZED VIEW mv AS SELECT * FROM users WHERE id IS NOT NULL PRIMARY KEY (id)",
		"CREATE MATERIALIZED VIEW mv AS SELECT col1, col2 FROM users WHERE col1 IS NOT NULL AND col2 IS NOT NULL PRIMARY KEY (col1, col2)",

		// Auth
		"GRANT SELECT ON TABLE users TO reader",
		"REVOKE ALL ON ALL KEYSPACES FROM admin",
		"LIST ALL PERMISSIONS OF admin",
		"LIST ROLES",
		"LIST ROLES OF admin NORECURSIVE",
		"CREATE ROLE myrole WITH PASSWORD = 'secret' AND LOGIN = true",
		"ALTER ROLE myrole WITH PASSWORD = 'newsecret'",
		"DROP ROLE IF EXISTS myrole",

		// Multi-statement
		"SELECT * FROM users; INSERT INTO users (id) VALUES (1); USE ks",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			violations := CheckLocations(t, sql)
			for _, v := range violations {
				t.Errorf("Loc violation: %s", v)
			}
		})
	}
}
