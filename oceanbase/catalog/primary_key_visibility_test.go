package catalog

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// addPrimaryKeyIgnoredOptionAlters each add a primary key with an option OceanBase
// takes and ignores: INVISIBLE, after which the key reads back visible, and COMMENT,
// after which it reads back with none.
var addPrimaryKeyIgnoredOptionAlters = []string{
	"ALTER TABLE t ADD PRIMARY KEY (id) INVISIBLE",
	"ALTER TABLE t ADD CONSTRAINT pkc PRIMARY KEY (id) INVISIBLE",
	"ALTER TABLE t ADD PRIMARY KEY (id) INVISIBLE VISIBLE",
	"ALTER TABLE t ADD PRIMARY KEY (id) VISIBLE INVISIBLE",
	"ALTER TABLE t ADD PRIMARY KEY (id) GLOBAL INVISIBLE",
	"ALTER TABLE t ADD PRIMARY KEY (id) COMMENT 'alter pk comment'",
	"ALTER TABLE t ADD CONSTRAINT PRIMARY KEY (id) COMMENT 'x' INVISIBLE",
	"ALTER TABLE t ADD PRIMARY KEY (id) KEY_BLOCK_SIZE 8 COMMENT 'pkc' USING HASH",
}

// createPrimaryKeyComment keeps its COMMENT on OceanBase, unlike the same key added
// by ALTER TABLE.
const createPrimaryKeyComment = "CREATE TABLE t (id INT, PRIMARY KEY (id) COMMENT 'pk comment')"

// alterIndexPrimaryAlters name the primary key in ALTER INDEX, which OceanBase does
// not reach: it answers 1176, Key 'PRIMARY' doesn't exist.
var alterIndexPrimaryAlters = []string{
	"ALTER TABLE t ALTER INDEX `PRIMARY` INVISIBLE",
	"ALTER TABLE t ALTER INDEX `PRIMARY` VISIBLE",
}

// primaryIndex returns tbl's primary key, or nil.
func primaryIndex(tbl *Table) *Index {
	for _, idx := range tbl.Indexes {
		if idx.Primary {
			return idx
		}
	}
	return nil
}

// TestAlterAddPrimaryKeyIgnoredOptions pins that an ALTER TABLE primary key is added
// visible and without a comment, whatever the DDL asks, as OceanBase adds it. MySQL
// refuses INVISIBLE there with 3522 and keeps the comment.
// TestOracle_OceanBaseAlterAddPrimaryKey asks the engine.
func TestAlterAddPrimaryKeyIgnoredOptions(t *testing.T) {
	for _, alter := range addPrimaryKeyIgnoredOptionAlters {
		t.Run(alter, func(t *testing.T) {
			c := New()
			sql := "CREATE DATABASE d; USE d; CREATE TABLE t (id INT NOT NULL, a INT); " + alter + ";"
			results, err := c.Exec(sql, nil)
			if err != nil {
				t.Fatal(err)
			}
			for _, r := range results {
				if r.Error != nil {
					t.Fatalf("Exec(%q): %v", alter, r.Error)
				}
			}
			pk := primaryIndex(c.GetDatabase("d").GetTable("t"))
			if pk == nil {
				t.Fatal("no primary key after the ALTER")
			}
			if !pk.Visible {
				t.Error("the primary key is invisible; OceanBase keeps it visible")
			}
			if pk.Comment != "" {
				t.Errorf("the primary key has comment %q; OceanBase drops it", pk.Comment)
			}
			if ddl := c.ShowCreateTable("d", "t"); strings.Contains(ddl, "INVISIBLE") || strings.Contains(ddl, "COMMENT") {
				t.Errorf("SHOW CREATE TABLE keeps an option OceanBase drops:\n%s", ddl)
			}
		})
	}
}

// TestCreatePrimaryKeyKeepsComment pins the other side: a primary key's COMMENT in
// CREATE TABLE stays, as OceanBase keeps it.
func TestCreatePrimaryKeyKeepsComment(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE d; USE d; "+createPrimaryKeyComment+";", nil); err != nil {
		t.Fatal(err)
	}
	if pk := primaryIndex(c.GetDatabase("d").GetTable("t")); pk == nil || pk.Comment != "pk comment" {
		t.Errorf("primary key = %+v, want the comment 'pk comment'", pk)
	}
}

// TestAlterIndexPrimaryRefused pins that ALTER INDEX does not reach the primary key,
// as on OceanBase, and leaves it visible.
func TestAlterIndexPrimaryRefused(t *testing.T) {
	for _, alter := range alterIndexPrimaryAlters {
		t.Run(alter, func(t *testing.T) {
			c := New()
			if _, err := c.Exec("CREATE DATABASE d; USE d; CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT);", nil); err != nil {
				t.Fatal(err)
			}
			results, err := c.Exec(alter+";", nil)
			if err != nil {
				t.Fatal(err)
			}
			if len(results) != 1 || results[0].Error == nil {
				t.Fatalf("Exec(%q) succeeded; OceanBase answers 1176", alter)
			}
			if msg := results[0].Error.Error(); !strings.Contains(msg, "Key 'PRIMARY' doesn't exist") {
				t.Errorf("Exec(%q) error = %q, want OceanBase's Key 'PRIMARY' doesn't exist", alter, msg)
			}
			if pk := primaryIndex(c.GetDatabase("d").GetTable("t")); pk == nil || !pk.Visible {
				t.Error("the refused ALTER changed the primary key")
			}
		})
	}
}

// TestOracle_OceanBaseAlterAddPrimaryKey applies each of
// addPrimaryKeyIgnoredOptionAlters on OceanBase, reads the table back, and requires the
// engine to report the key visible and without a comment, and the catalog's result to
// diff empty against the readback on its indexes. It requires the same of
// createPrimaryKeyComment, whose comment stays, and the engine to refuse
// alterIndexPrimaryAlters with 1176.
//
//	scripts/oceanbase-oracle.sh go test ./oceanbase/catalog/ -run OceanBase -count=1
//
// requirePrimaryKeyMatchesReadback loads userDDL and the engine's readback and requires
// the readback's primary key to be visible with comment, and the two to diff empty on
// their indexes.
func requirePrimaryKeyMatchesReadback(t *testing.T, o *oracleConn, userDDL, readback, comment string) {
	t.Helper()
	sc := serverCharsetFor(o.version)
	userCat, _ := loadOneTable(t, sc, userDDL, "t")
	storedCat, stored := loadOneTable(t, sc, readback, "t")
	if pk := primaryIndex(stored); pk == nil || !pk.Visible || pk.Comment != comment {
		t.Errorf("the readback's primary key is %+v, want it visible with comment %q:\n%s", pk, comment, readback)
	}
	for _, te := range DiffWithNormalizer(userCat, storedCat, NormalizerFor(o.version)).Tables {
		for _, ie := range te.Indexes {
			t.Errorf("index %s reported %s against the readback\n%s", ie.Name, ie.Action, readback)
		}
	}
}

func TestOracle_OceanBaseAlterAddPrimaryKey(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx := context.Background()
	const create = "CREATE TABLE t (id INT NOT NULL, a INT)"
	for i, alter := range addPrimaryKeyIgnoredOptionAlters {
		t.Run(alter, func(t *testing.T) {
			dbName := fmt.Sprintf("obpkvis_%d", i)
			if _, err := storeAndRead(ctx, o, dbName, create, "t"); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _, _ = o.db.ExecContext(ctx, "DROP DATABASE IF EXISTS "+dbName) })
			if _, err := o.db.ExecContext(ctx, alter); err != nil {
				t.Fatalf("OceanBase refused %q: %v", alter, err)
			}
			var visible string
			if err := o.db.QueryRowContext(ctx, "SELECT IS_VISIBLE FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' AND INDEX_NAME = 'PRIMARY'", dbName).Scan(&visible); err != nil {
				t.Fatal(err)
			}
			if visible != "YES" {
				t.Errorf("OceanBase reports the key IS_VISIBLE = %q, the catalog assumes YES", visible)
			}
			var readback string
			if err := o.db.QueryRowContext(ctx, "SHOW CREATE TABLE "+dbName+".t").Scan(new(string), &readback); err != nil {
				t.Fatal(err)
			}
			requirePrimaryKeyMatchesReadback(t, o, create+";\n"+alter+";", readback, "")
		})
	}
	t.Run("create table comment", func(t *testing.T) {
		readback, err := storeAndRead(ctx, o, "obpkvis_create", createPrimaryKeyComment, "t")
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _, _ = o.db.ExecContext(ctx, "DROP DATABASE IF EXISTS obpkvis_create") })
		requirePrimaryKeyMatchesReadback(t, o, createPrimaryKeyComment+";", readback, "pk comment")
	})
	t.Run("alter index primary", func(t *testing.T) {
		if _, err := storeAndRead(ctx, o, "obpkvis_alter", "CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)", "t"); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _, _ = o.db.ExecContext(ctx, "DROP DATABASE IF EXISTS obpkvis_alter") })
		for _, alter := range alterIndexPrimaryAlters {
			_, err := o.db.ExecContext(ctx, alter)
			var myErr *gomysql.MySQLError
			if !errors.As(err, &myErr) || myErr.Number != 1176 {
				t.Errorf("OceanBase answered %v for %q, want 1176", err, alter)
			}
		}
	})
}
