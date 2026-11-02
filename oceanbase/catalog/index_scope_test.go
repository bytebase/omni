package catalog

import (
	"errors"
	"strings"
	"testing"
)

// TestIndexScopeLastWins pins how a repeated scope reads. OceanBase's index options
// repeat freely and the last scope stands: on 5.7.25-OceanBase_CE-v4.3.5.6, KEY k (a)
// LOCAL GLOBAL reads back GLOBAL and GLOBAL LOCAL reads back LOCAL, and a UNIQUE KEY
// that leaves out the partitioning column is accepted as LOCAL GLOBAL and refused with
// 1503 as GLOBAL LOCAL.
func TestIndexScopeLastWins(t *testing.T) {
	for _, tc := range []struct {
		ddl, scope string
	}{
		{"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) LOCAL GLOBAL)", "GLOBAL"},
		{"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) GLOBAL LOCAL)", "LOCAL"},
		{"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) LOCAL COMMENT 'x' GLOBAL)", "GLOBAL"},
		{"CREATE TABLE t (id INT, a INT, PRIMARY KEY (id), UNIQUE KEY k (a) LOCAL GLOBAL) PARTITION BY HASH(id) PARTITIONS 2", "GLOBAL"},
	} {
		c := setupWithDB(t)
		mustExec(t, c, tc.ddl)
		var got string
		for _, idx := range c.GetDatabase("testdb").GetTable("t").Indexes {
			if idx.Name == "k" {
				got = idx.Scope
			}
		}
		if got != tc.scope {
			t.Errorf("%s: scope = %q, want %q", tc.ddl, got, tc.scope)
		}
	}

	c := setupWithDB(t)
	results, err := c.Exec("CREATE TABLE t (id INT, a INT, PRIMARY KEY (id), UNIQUE KEY k (a) GLOBAL LOCAL) PARTITION BY HASH(id) PARTITIONS 2", nil)
	for _, r := range results {
		if r.Error != nil {
			err = r.Error
		}
	}
	var ce *Error
	if !errors.As(err, &ce) || ce.Code != 1503 {
		t.Errorf("UNIQUE KEY ... GLOBAL LOCAL leaving out the partitioning column: err = %v, want 1503", err)
	}
}

// TestAlterAddPrimaryKeyKeepsNoScope pins that a primary key carries no scope.
// ALTER TABLE ... ADD PRIMARY KEY takes LOCAL and GLOBAL, and 4.3.5.6 accepts them and
// reads the key back as PRIMARY KEY (`id`) with neither. Keeping the scope would also
// print it into SHOW CREATE TABLE, where CREATE TABLE refuses it with 1064.
func TestAlterAddPrimaryKeyKeepsNoScope(t *testing.T) {
	for _, alter := range []string{
		"ALTER TABLE t ADD PRIMARY KEY (id) GLOBAL",
		"ALTER TABLE t ADD PRIMARY KEY (id) LOCAL",
		"ALTER TABLE t ADD CONSTRAINT pk PRIMARY KEY (id) GLOBAL",
	} {
		c := setupWithDB(t)
		mustExec(t, c, "CREATE TABLE t (id INT NOT NULL, a INT)")
		mustExec(t, c, alter)
		for _, idx := range c.GetDatabase("testdb").GetTable("t").Indexes {
			if idx.Primary && idx.Scope != "" {
				t.Errorf("%s: the primary key kept scope %q", alter, idx.Scope)
			}
		}
		if show := c.ShowCreateTable("testdb", "t"); strings.Contains(show, "GLOBAL") || strings.Contains(show, "LOCAL") {
			t.Errorf("%s: SHOW CREATE TABLE prints a scope on the primary key:\n%s", alter, show)
		}
	}
}

// TestPrimaryKeyCoversPartitionColumns pins that only a secondary GLOBAL unique index
// skips the partition-key coverage rule. A primary key that leaves out the
// partitioning column is a 1503 whatever its scope: 4.3.5.6 answers 1503 for ALTER
// TABLE ... ADD PRIMARY KEY (a) GLOBAL on a table partitioned by HASH(id).
func TestPrimaryKeyCoversPartitionColumns(t *testing.T) {
	for _, tc := range []struct {
		name    string
		idx     *Index
		wantErr bool
	}{
		{"global primary key", &Index{Name: "PRIMARY", Primary: true, Unique: true, Scope: "GLOBAL"}, true},
		{"global unique key", &Index{Name: "ua", Unique: true, Scope: "GLOBAL"}, false},
		{"local unique key", &Index{Name: "ua", Unique: true, Scope: "LOCAL"}, true},
	} {
		tc.idx.Columns = []*IndexColumn{{Name: "a"}}
		err := validateUniqueKeysCoverPartitionColumns(&Table{Name: "t", Indexes: []*Index{tc.idx}}, []string{"id"})
		if (err != nil) != tc.wantErr {
			t.Errorf("%s: err = %v, want error %v", tc.name, err, tc.wantErr)
		}
	}
}
