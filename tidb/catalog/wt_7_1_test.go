package catalog

import (
	"strings"
	"testing"
)

// --- Section 3.1 (Phase 3): RANGE Partitioning (8 scenarios) ---
// File target: wt_7_1_test.go
// Proof: go test ./mysql/catalog/ -short -count=1 -run "TestWalkThrough_7_1"

func TestWalkThrough_7_1_RangePartitioning(t *testing.T) {
	t.Run("range_expr_3_partitions_maxvalue", func(t *testing.T) {
		// Scenario 1: CREATE TABLE PARTITION BY RANGE (expr) with 3 partitions + MAXVALUE
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t1 (
			id INT NOT NULL,
			store_id INT NOT NULL
		)
		PARTITION BY RANGE (store_id) (
			PARTITION p0 VALUES LESS THAN (10),
			PARTITION p1 VALUES LESS THAN (20),
			PARTITION p2 VALUES LESS THAN MAXVALUE
		)`)

		tbl := c.GetDatabase("testdb").GetTable("t1")
		if tbl == nil {
			t.Fatal("table not found")
		}
		if tbl.Partitioning == nil {
			t.Fatal("expected partitioning info, got nil")
		}
		if tbl.Partitioning.Type != "RANGE" {
			t.Errorf("expected type RANGE, got %q", tbl.Partitioning.Type)
		}
		if len(tbl.Partitioning.Partitions) != 3 {
			t.Fatalf("expected 3 partitions, got %d", len(tbl.Partitioning.Partitions))
		}
		if tbl.Partitioning.Partitions[0].Name != "p0" {
			t.Errorf("expected partition name p0, got %q", tbl.Partitioning.Partitions[0].Name)
		}

		ddl := c.ShowCreateTable("testdb", "t1")
		// Verify version comment /*!50100
		if !strings.Contains(ddl, "/*!50100") {
			t.Errorf("expected /*!50100 version comment in DDL:\n%s", ddl)
		}
		// Verify MAXVALUE without parens for plain RANGE
		if !strings.Contains(ddl, "VALUES LESS THAN MAXVALUE") {
			t.Errorf("expected 'VALUES LESS THAN MAXVALUE' (no parens) in DDL:\n%s", ddl)
		}
		// Verify VALUES LESS THAN (10)
		if !strings.Contains(ddl, "VALUES LESS THAN (10)") {
			t.Errorf("expected 'VALUES LESS THAN (10)' in DDL:\n%s", ddl)
		}
	})

	t.Run("range_columns_single", func(t *testing.T) {
		// Scenario 2: CREATE TABLE PARTITION BY RANGE COLUMNS (col) — single column
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t2 (
			id INT NOT NULL,
			name VARCHAR(50) NOT NULL
		)
		PARTITION BY RANGE COLUMNS (name) (
			PARTITION p0 VALUES LESS THAN ('g'),
			PARTITION p1 VALUES LESS THAN ('n'),
			PARTITION p2 VALUES LESS THAN (MAXVALUE)
		)`)

		tbl := c.GetDatabase("testdb").GetTable("t2")
		if tbl == nil {
			t.Fatal("table not found")
		}
		if tbl.Partitioning == nil {
			t.Fatal("expected partitioning info, got nil")
		}
		if tbl.Partitioning.Type != "RANGE COLUMNS" {
			t.Errorf("expected type RANGE COLUMNS, got %q", tbl.Partitioning.Type)
		}
		if len(tbl.Partitioning.Columns) != 1 {
			t.Fatalf("expected 1 column, got %d", len(tbl.Partitioning.Columns))
		}

		ddl := c.ShowCreateTable("testdb", "t2")
		// RANGE COLUMNS uses /*!50500
		if !strings.Contains(ddl, "/*!50500") {
			t.Errorf("expected /*!50500 version comment in DDL:\n%s", ddl)
		}
		// RANGE COLUMNS MAXVALUE is parenthesized
		if !strings.Contains(ddl, "VALUES LESS THAN (MAXVALUE)") {
			t.Errorf("expected 'VALUES LESS THAN (MAXVALUE)' in DDL:\n%s", ddl)
		}
		// MySQL uses double space before COLUMNS
		if !strings.Contains(ddl, "RANGE  COLUMNS") {
			t.Errorf("expected 'RANGE  COLUMNS' (double space) in DDL:\n%s", ddl)
		}
	})

	t.Run("range_columns_multi", func(t *testing.T) {
		// Scenario 3: CREATE TABLE PARTITION BY RANGE COLUMNS (col1, col2) — multi-column
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t3 (
			a INT NOT NULL,
			b INT NOT NULL,
			c INT NOT NULL
		)
		PARTITION BY RANGE COLUMNS (a, b) (
			PARTITION p0 VALUES LESS THAN (10, 20),
			PARTITION p1 VALUES LESS THAN (20, 30),
			PARTITION p2 VALUES LESS THAN (MAXVALUE, MAXVALUE)
		)`)

		tbl := c.GetDatabase("testdb").GetTable("t3")
		if tbl == nil {
			t.Fatal("table not found")
		}
		if tbl.Partitioning == nil {
			t.Fatal("expected partitioning info, got nil")
		}
		if tbl.Partitioning.Type != "RANGE COLUMNS" {
			t.Errorf("expected type RANGE COLUMNS, got %q", tbl.Partitioning.Type)
		}
		if len(tbl.Partitioning.Columns) != 2 {
			t.Fatalf("expected 2 columns, got %d", len(tbl.Partitioning.Columns))
		}
		if len(tbl.Partitioning.Partitions) != 3 {
			t.Fatalf("expected 3 partitions, got %d", len(tbl.Partitioning.Partitions))
		}

		ddl := c.ShowCreateTable("testdb", "t3")
		if !strings.Contains(ddl, "RANGE  COLUMNS(a,b)") {
			t.Errorf("expected 'RANGE  COLUMNS(a,b)' in DDL:\n%s", ddl)
		}
		if !strings.Contains(ddl, "VALUES LESS THAN (10,20)") {
			t.Errorf("expected 'VALUES LESS THAN (10,20)' in DDL:\n%s", ddl)
		}
	})

	t.Run("alter_add_partition", func(t *testing.T) {
		// Scenario 4: ALTER TABLE ADD PARTITION to RANGE table
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t4 (
			id INT NOT NULL,
			val INT NOT NULL
		)
		PARTITION BY RANGE (val) (
			PARTITION p0 VALUES LESS THAN (10),
			PARTITION p1 VALUES LESS THAN (20)
		)`)

		tbl := c.GetDatabase("testdb").GetTable("t4")
		if len(tbl.Partitioning.Partitions) != 2 {
			t.Fatalf("expected 2 partitions before ADD, got %d", len(tbl.Partitioning.Partitions))
		}

		wtExec(t, c, `ALTER TABLE t4 ADD PARTITION (PARTITION p2 VALUES LESS THAN (30))`)

		tbl = c.GetDatabase("testdb").GetTable("t4")
		if len(tbl.Partitioning.Partitions) != 3 {
			t.Fatalf("expected 3 partitions after ADD, got %d", len(tbl.Partitioning.Partitions))
		}
		if tbl.Partitioning.Partitions[2].Name != "p2" {
			t.Errorf("expected new partition name p2, got %q", tbl.Partitioning.Partitions[2].Name)
		}

		ddl := c.ShowCreateTable("testdb", "t4")
		if !strings.Contains(ddl, "PARTITION p2 VALUES LESS THAN (30)") {
			t.Errorf("expected p2 partition in DDL:\n%s", ddl)
		}
	})

	t.Run("alter_drop_partition", func(t *testing.T) {
		// Scenario 5: ALTER TABLE DROP PARTITION from RANGE table
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t5 (
			id INT NOT NULL,
			val INT NOT NULL
		)
		PARTITION BY RANGE (val) (
			PARTITION p0 VALUES LESS THAN (10),
			PARTITION p1 VALUES LESS THAN (20),
			PARTITION p2 VALUES LESS THAN (30)
		)`)

		tbl := c.GetDatabase("testdb").GetTable("t5")
		if len(tbl.Partitioning.Partitions) != 3 {
			t.Fatalf("expected 3 partitions before DROP, got %d", len(tbl.Partitioning.Partitions))
		}

		wtExec(t, c, `ALTER TABLE t5 DROP PARTITION p1`)

		tbl = c.GetDatabase("testdb").GetTable("t5")
		if len(tbl.Partitioning.Partitions) != 2 {
			t.Fatalf("expected 2 partitions after DROP, got %d", len(tbl.Partitioning.Partitions))
		}
		// Remaining should be p0 and p2
		if tbl.Partitioning.Partitions[0].Name != "p0" {
			t.Errorf("expected first partition p0, got %q", tbl.Partitioning.Partitions[0].Name)
		}
		if tbl.Partitioning.Partitions[1].Name != "p2" {
			t.Errorf("expected second partition p2, got %q", tbl.Partitioning.Partitions[1].Name)
		}

		ddl := c.ShowCreateTable("testdb", "t5")
		if strings.Contains(ddl, "p1") {
			t.Errorf("dropped partition p1 should not appear in DDL:\n%s", ddl)
		}
	})

	t.Run("alter_reorganize_split", func(t *testing.T) {
		// Scenario 6: ALTER TABLE REORGANIZE PARTITION split (1->2)
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t6 (
			id INT NOT NULL,
			val INT NOT NULL
		)
		PARTITION BY RANGE (val) (
			PARTITION p0 VALUES LESS THAN (10),
			PARTITION p1 VALUES LESS THAN (30),
			PARTITION p2 VALUES LESS THAN MAXVALUE
		)`)

		wtExec(t, c, `ALTER TABLE t6 REORGANIZE PARTITION p1 INTO (
			PARTITION p1a VALUES LESS THAN (20),
			PARTITION p1b VALUES LESS THAN (30)
		)`)

		tbl := c.GetDatabase("testdb").GetTable("t6")
		if len(tbl.Partitioning.Partitions) != 4 {
			t.Fatalf("expected 4 partitions after split, got %d", len(tbl.Partitioning.Partitions))
		}
		names := make([]string, len(tbl.Partitioning.Partitions))
		for i, p := range tbl.Partitioning.Partitions {
			names[i] = p.Name
		}
		expected := []string{"p0", "p1a", "p1b", "p2"}
		for i, exp := range expected {
			if names[i] != exp {
				t.Errorf("partition %d: expected %q, got %q", i, exp, names[i])
			}
		}

		ddl := c.ShowCreateTable("testdb", "t6")
		if !strings.Contains(ddl, "PARTITION p1a VALUES LESS THAN (20)") {
			t.Errorf("expected p1a partition in DDL:\n%s", ddl)
		}
		if !strings.Contains(ddl, "PARTITION p1b VALUES LESS THAN (30)") {
			t.Errorf("expected p1b partition in DDL:\n%s", ddl)
		}
	})

	t.Run("alter_reorganize_merge", func(t *testing.T) {
		// Scenario 7: ALTER TABLE REORGANIZE PARTITION merge (2->1)
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t7 (
			id INT NOT NULL,
			val INT NOT NULL
		)
		PARTITION BY RANGE (val) (
			PARTITION p0 VALUES LESS THAN (10),
			PARTITION p1 VALUES LESS THAN (20),
			PARTITION p2 VALUES LESS THAN (30),
			PARTITION p3 VALUES LESS THAN MAXVALUE
		)`)

		wtExec(t, c, `ALTER TABLE t7 REORGANIZE PARTITION p1, p2 INTO (
			PARTITION p_merged VALUES LESS THAN (30)
		)`)

		tbl := c.GetDatabase("testdb").GetTable("t7")
		if len(tbl.Partitioning.Partitions) != 3 {
			t.Fatalf("expected 3 partitions after merge, got %d", len(tbl.Partitioning.Partitions))
		}
		names := make([]string, len(tbl.Partitioning.Partitions))
		for i, p := range tbl.Partitioning.Partitions {
			names[i] = p.Name
		}
		expected := []string{"p0", "p_merged", "p3"}
		for i, exp := range expected {
			if names[i] != exp {
				t.Errorf("partition %d: expected %q, got %q", i, exp, names[i])
			}
		}

		ddl := c.ShowCreateTable("testdb", "t7")
		if !strings.Contains(ddl, "PARTITION p_merged VALUES LESS THAN (30)") {
			t.Errorf("expected p_merged partition in DDL:\n%s", ddl)
		}
		if strings.Contains(ddl, "PARTITION p1 ") || strings.Contains(ddl, "PARTITION p2 ") {
			t.Errorf("old partitions p1/p2 should not appear in DDL:\n%s", ddl)
		}
	})

	t.Run("range_date_expression", func(t *testing.T) {
		// Scenario 8: RANGE partition with date expression (YEAR(col))
		c := wtSetup(t)
		wtExec(t, c, `CREATE TABLE t8 (
			id INT NOT NULL,
			created_at DATE NOT NULL
		)
		PARTITION BY RANGE (YEAR(created_at)) (
			PARTITION p2020 VALUES LESS THAN (2021),
			PARTITION p2021 VALUES LESS THAN (2022),
			PARTITION p2022 VALUES LESS THAN (2023),
			PARTITION pfuture VALUES LESS THAN MAXVALUE
		)`)

		tbl := c.GetDatabase("testdb").GetTable("t8")
		if tbl == nil {
			t.Fatal("table not found")
		}
		if tbl.Partitioning == nil {
			t.Fatal("expected partitioning info, got nil")
		}
		if tbl.Partitioning.Type != "RANGE" {
			t.Errorf("expected type RANGE, got %q", tbl.Partitioning.Type)
		}
		if len(tbl.Partitioning.Partitions) != 4 {
			t.Fatalf("expected 4 partitions, got %d", len(tbl.Partitioning.Partitions))
		}

		ddl := c.ShowCreateTable("testdb", "t8")
		// Verify the expression rendering includes YEAR(...) (MySQL renders as lowercase year())
		upperDDL := strings.ToUpper(ddl)
		if !strings.Contains(upperDDL, "YEAR(") {
			t.Errorf("expected YEAR() expression in DDL:\n%s", ddl)
		}
		if !strings.Contains(ddl, "VALUES LESS THAN (2021)") {
			t.Errorf("expected 'VALUES LESS THAN (2021)' in DDL:\n%s", ddl)
		}
		if !strings.Contains(ddl, "VALUES LESS THAN MAXVALUE") {
			t.Errorf("expected 'VALUES LESS THAN MAXVALUE' (no parens) in DDL:\n%s", ddl)
		}
	})
}
