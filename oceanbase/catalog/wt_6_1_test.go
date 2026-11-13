package catalog

import (
	"slices"
	"testing"
)

// TestWalkThrough_6_1_ForeignKeyMakesNoIndex pins that a table has the indexes it declares
// and no other, whatever foreign keys it has and however they got there. OceanBase makes no
// index for a foreign key, where MySQL adds one when no index serves the key
// (TestOracle_OceanBaseForeignKeyIndex asks the engine).
func TestWalkThrough_6_1_ForeignKeyMakesNoIndex(t *testing.T) {
	const parent = "CREATE TABLE parent (id INT NOT NULL, PRIMARY KEY (id))"
	for _, tc := range []struct {
		name string
		ddl  []string
		want []string
	}{
		{"named", []string{parent,
			"CREATE TABLE child (id INT NOT NULL, parent_id INT NOT NULL, PRIMARY KEY (id), CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id))",
		}, []string{"PRIMARY"}},
		{"unnamed", []string{parent,
			"CREATE TABLE child (id INT NOT NULL, parent_id INT NOT NULL, PRIMARY KEY (id), FOREIGN KEY (parent_id) REFERENCES parent(id))",
		}, []string{"PRIMARY"}},
		{"beside an index over its column", []string{parent,
			"CREATE TABLE child (id INT NOT NULL, parent_id INT NOT NULL, PRIMARY KEY (id), INDEX idx_parent (parent_id), CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id))",
		}, []string{"PRIMARY", "idx_parent"}},
		{"beside a unique key over its column", []string{parent,
			"CREATE TABLE child (id INT NOT NULL, parent_id INT NOT NULL, PRIMARY KEY (id), UNIQUE KEY uk_parent (parent_id), CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id))",
		}, []string{"PRIMARY", "uk_parent"}},
		{"over the primary key", []string{parent,
			"CREATE TABLE child (id INT NOT NULL, PRIMARY KEY (id), CONSTRAINT fk_parent FOREIGN KEY (id) REFERENCES parent(id))",
		}, []string{"PRIMARY"}},
		{"two columns, an index over one", []string{
			"CREATE TABLE parent (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a, b))",
			"CREATE TABLE child (id INT NOT NULL, a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (id), INDEX idx_a (a), CONSTRAINT fk_ab FOREIGN KEY (a, b) REFERENCES parent(a, b))",
		}, []string{"PRIMARY", "idx_a"}},
		{"added by ALTER TABLE", []string{parent,
			"CREATE TABLE child (id INT NOT NULL, parent_id INT NOT NULL, PRIMARY KEY (id))",
			"ALTER TABLE child ADD CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id)",
		}, []string{"PRIMARY"}},
		{"added by ALTER TABLE beside an index", []string{parent,
			"CREATE TABLE child (id INT NOT NULL, parent_id INT NOT NULL, PRIMARY KEY (id), INDEX idx_parent (parent_id))",
			"ALTER TABLE child ADD CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id)",
		}, []string{"PRIMARY", "idx_parent"}},
		{"dropped", []string{parent,
			"CREATE TABLE child (id INT NOT NULL, parent_id INT NOT NULL, PRIMARY KEY (id), INDEX idx_parent (parent_id), CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id))",
			"ALTER TABLE child DROP FOREIGN KEY fk_parent",
		}, []string{"PRIMARY", "idx_parent"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := setupWithDB(t)
			for _, stmt := range tc.ddl {
				mustExec(t, c, stmt)
			}
			var got []string
			for _, idx := range c.GetDatabase("testdb").GetTable("child").Indexes {
				got = append(got, idx.Name)
			}
			if !slices.Equal(got, tc.want) {
				t.Errorf("indexes %v, want %v", got, tc.want)
			}
		})
	}

	// An index named as the foreign key was never made, so there is none to drop: OceanBase
	// answers 1091 (measured on 4.3.5.6).
	t.Run("DROP INDEX of the foreign key's name", func(t *testing.T) {
		c := setupWithDB(t)
		mustExec(t, c, parent)
		mustExec(t, c, "CREATE TABLE child (id INT NOT NULL, parent_id INT NOT NULL, PRIMARY KEY (id), CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id))")
		assertError(t, execError(t, c, "ALTER TABLE child DROP INDEX fk_parent"), 1091)
		assertError(t, execError(t, c, "ALTER TABLE child DROP FOREIGN KEY fk_parent, DROP INDEX fk_parent"), 1091)
	})
}
