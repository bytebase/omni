package catalog

import (
	"slices"
	"sort"
	"strings"
	"testing"
)

// TestScenario_PS covers section PS of mysql/catalog/SCENARIOS-mysql-implicit-behavior.md
// — "Path-split behaviors (CREATE vs ALTER)". Each subtest runs the scenario's
// DDL against both OceanBase and the omni catalog and asserts
// both match the expected value. Existing TestBugFix_* tests remain unchanged;
// these TestScenario_PS tests are the durable dual-assertion versions.
func TestScenario_PS(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// PS.1 CHECK constraint counter — CREATE path (fresh counter).
	// User-named t_chk_5 is NOT seeded into the generated counter; the
	// single unnamed CHECK receives t_chk_1.
	t.Run("PS_1_Check_counter_CREATE_fresh", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (
            a INT,
            CONSTRAINT t_chk_5 CHECK (a > 0),
            b INT,
            CHECK (b < 100)
        )`)

		// OceanBase names the unnamed CHECK after the table and the time,
		// compared here without the count a tape holds for the time.
		oracleWant := []string{"t_chk_5", "t_OBCHECK_<n>"}
		// omni gap: the catalog numbers it t_chk_1, as MySQL does.
		want := []string{"t_chk_1", "t_chk_5"}

		rows := oracleRows(t, mc, `
            SELECT CONSTRAINT_NAME FROM information_schema.CHECK_CONSTRAINTS
            WHERE CONSTRAINT_SCHEMA='testdb'
            ORDER BY CONSTRAINT_NAME`)
		var oracleNames []string
		for _, r := range rows {
			oracleNames = append(oracleNames, withoutCount(asString(r[0])))
		}
		if !slices.Equal(oracleNames, oracleWant) {
			t.Errorf("PS.1 oracle CHECK names: got %v, want %v", oracleNames, oracleWant)
		}

		omniNames := psCheckNames(c, "t")
		if !slices.Equal(omniNames, want) {
			t.Errorf("PS.1 omni CHECK names: got %v, want %v", omniNames, want)
		}
	})

	// PS.2 CHECK constraint counter — ALTER path (max+1).
	// The user-named t_chk_20 IS seeded on the ALTER path; new unnamed
	// CHECK gets t_chk_21.
	t.Run("PS_2_Check_counter_ALTER_maxplus1", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			`CREATE TABLE t (a INT, b INT, CONSTRAINT t_chk_20 CHECK (a>0))`)
		runOnBoth(t, mc, c, `ALTER TABLE t ADD CHECK (b>0)`)

		// OceanBase names the unnamed CHECK after the table and the time.
		oracleWant := []string{"t_chk_20", "t_OBCHECK_<n>"}
		// omni gap: the catalog continues MySQL's counter to t_chk_21.
		want := []string{"t_chk_20", "t_chk_21"}

		rows := oracleRows(t, mc, `
            SELECT CONSTRAINT_NAME FROM information_schema.CHECK_CONSTRAINTS
            WHERE CONSTRAINT_SCHEMA='testdb'
            ORDER BY CONSTRAINT_NAME`)
		var oracleNames []string
		for _, r := range rows {
			oracleNames = append(oracleNames, withoutCount(asString(r[0])))
		}
		if !slices.Equal(oracleNames, oracleWant) {
			t.Errorf("PS.2 oracle CHECK names: got %v, want %v", oracleNames, oracleWant)
		}

		omniNames := psCheckNames(c, "t")
		if !slices.Equal(omniNames, want) {
			t.Errorf("PS.2 omni CHECK names: got %v, want %v", omniNames, want)
		}
	})

	// PS.3 FK counter — CREATE path (fresh, user-named NOT seeded).
	t.Run("PS_3_FK_counter_CREATE_fresh", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE parent (id INT PRIMARY KEY)`)
		runOnBoth(t, mc, c, `CREATE TABLE child (
            a INT,
            CONSTRAINT child_ibfk_5 FOREIGN KEY (a) REFERENCES parent(id),
            b INT,
            FOREIGN KEY (b) REFERENCES parent(id)
        )`)

		// OceanBase names the unnamed FK after the table and the time.
		oracleWant := []string{"child_ibfk_5", "child_OBFK_<n>"}
		// omni gap: the catalog numbers it child_ibfk_1, as MySQL does.
		want := []string{"child_ibfk_1", "child_ibfk_5"}

		rows := oracleRows(t, mc, `
            SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='child'
              AND CONSTRAINT_TYPE='FOREIGN KEY'
            ORDER BY CONSTRAINT_NAME`)
		var oracleNames []string
		for _, r := range rows {
			oracleNames = append(oracleNames, withoutCount(asString(r[0])))
		}
		if !slices.Equal(oracleNames, oracleWant) {
			t.Errorf("PS.3 oracle FK names: got %v, want %v", oracleNames, oracleWant)
		}

		omniNames := psFKNames(c, "child")
		if !slices.Equal(omniNames, want) {
			t.Errorf("PS.3 omni FK names: got %v, want %v", omniNames, want)
		}
	})

	// PS.4 FK counter — ALTER path (max+1 over existing generated numbers).
	t.Run("PS_4_FK_counter_ALTER_maxplus1", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE parent (id INT PRIMARY KEY)`)
		runOnBoth(t, mc, c, `CREATE TABLE child (
            a INT,
            b INT,
            CONSTRAINT child_ibfk_20 FOREIGN KEY (a) REFERENCES parent(id)
        )`)
		runOnBoth(t, mc, c,
			`ALTER TABLE child ADD FOREIGN KEY (b) REFERENCES parent(id)`)

		// OceanBase names the unnamed FK after the table and the time.
		oracleWant := []string{"child_ibfk_20", "child_OBFK_<n>"}
		// omni gap: the catalog continues MySQL's counter to child_ibfk_21.
		want := []string{"child_ibfk_20", "child_ibfk_21"}

		rows := oracleRows(t, mc, `
            SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='child'
              AND CONSTRAINT_TYPE='FOREIGN KEY'
            ORDER BY CONSTRAINT_NAME`)
		var oracleNames []string
		for _, r := range rows {
			oracleNames = append(oracleNames, withoutCount(asString(r[0])))
		}
		if !slices.Equal(oracleNames, oracleWant) {
			t.Errorf("PS.4 oracle FK names: got %v, want %v", oracleNames, oracleWant)
		}

		omniNames := psFKNames(c, "child")
		if !slices.Equal(omniNames, want) {
			t.Errorf("PS.4 omni FK names: got %v, want %v", omniNames, want)
		}
	})

	// PS.5 DEFAULT NOW() / fsp precision mismatch must error.
	// MySQL rejects DATETIME(6) DEFAULT NOW() with ER_INVALID_DEFAULT (1067).
	// omni should reject as well.
	t.Run("PS_5_Datetime_fsp_mismatch_errors", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE t (a DATETIME(6) DEFAULT NOW())`

		_, oracleErr := mc.db.ExecContext(mc.ctx, ddl)
		if oracleErr == nil {
			t.Errorf("PS.5 oracle: expected MySQL to reject DATETIME(6) DEFAULT NOW() with ER_INVALID_DEFAULT, got success")
		} else if !strings.Contains(oracleErr.Error(), "1067") &&
			!strings.Contains(strings.ToLower(oracleErr.Error()), "invalid default") {
			t.Errorf("PS.5 oracle: unexpected error text: %v", oracleErr)
		}

		var omniErr error
		results, parseErr := c.Exec(ddl, nil)
		if parseErr != nil {
			omniErr = parseErr
		} else {
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
					break
				}
			}
		}
		if omniErr == nil {
			t.Errorf("PS.5 omni: expected ER_INVALID_DEFAULT-style error, got success")
		}
	})

	// PS.6 HASH partition ADD PARTITIONS n. MySQL seeds the new partitions from
	// the count, and this scenario held that when the package was forked from
	// mysql. OceanBase's ADD PARTITION takes only a list of RANGE or LIST
	// partitions and answers 1064 for the count (parser partitionOptionCases), so
	// omni must refuse it and keep the three partitions.
	t.Run("PS_6_Hash_partition_ADD_count_refused", func(t *testing.T) {
		c := scenarioNewCatalog(t)
		if _, err := c.Exec(`CREATE TABLE t (id INT) PARTITION BY HASH(id) PARTITIONS 3`, nil); err != nil {
			t.Fatal(err)
		}
		if _, err := c.Exec(`ALTER TABLE t ADD PARTITION PARTITIONS 2`, nil); err == nil {
			t.Errorf("PS.6 omni: ADD PARTITION PARTITIONS 2 parsed; OceanBase answers 1064")
		}
		var omniNames []string
		if tbl := c.GetDatabase("testdb").GetTable("t"); tbl != nil && tbl.Partitioning != nil {
			for _, p := range tbl.Partitioning.Partitions {
				omniNames = append(omniNames, p.Name)
			}
		}
		if want := []string{"p0", "p1", "p2"}; !slices.Equal(omniNames, want) {
			t.Errorf("PS.6 omni: partition names %v after the refused ADD, want %v", omniNames, want)
		}
	})

	// PS.7 FK name collision — user-named t_ibfk_1 collides with the
	// generator's implicit first unnamed FK name. MySQL errors with
	// ER_FK_DUP_NAME (1826). omni should reject as well.
	t.Run("PS_7_FK_name_collision_errors", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// Parent table first (both sides).
		runOnBoth(t, mc, c, `CREATE TABLE p (id INT PRIMARY KEY)`)

		ddl := `CREATE TABLE c (
            a INT,
            CONSTRAINT c_ibfk_1 FOREIGN KEY (a) REFERENCES p(id),
            b INT,
            FOREIGN KEY (b) REFERENCES p(id)
        )`

		// MySQL names the unnamed key c_ibfk_1, which the user took, and
		// answers ER_FK_DUP_NAME (1826). OceanBase's generated name cannot
		// collide, so the table is created.
		if _, err := mc.db.ExecContext(mc.ctx, ddl); err != nil {
			t.Errorf("PS.7 oracle: %v", err)
		}

		// omni gap: the catalog generates c_ibfk_1 and refuses the table, as
		// MySQL does.

		var omniErr error
		results, parseErr := c.Exec(ddl, nil)
		if parseErr != nil {
			omniErr = parseErr
		} else {
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
					break
				}
			}
		}
		if omniErr == nil {
			t.Errorf("PS.7 omni: expected ER_FK_DUP_NAME, got success")
		}
	})

	// PS.8 CHECK constraint duplicate name in schema — must error.
	// CHECK constraint names are schema-scoped in MySQL; the second
	// CREATE TABLE with the same check name fails with
	// ER_CHECK_CONSTRAINT_DUP_NAME (3822). omni should reject as well.
	t.Run("PS_8_Check_dup_name_schema_scope", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		// First create must succeed on both.
		runOnBoth(t, mc, c,
			`CREATE TABLE t1 (a INT, CONSTRAINT my_rule CHECK (a > 0))`)

		ddl2 := `CREATE TABLE t2 (b INT, CONSTRAINT my_rule CHECK (b > 0))`

		_, oracleErr := mc.db.ExecContext(mc.ctx, ddl2)
		if oracleErr == nil {
			t.Errorf("PS.8 oracle: expected ER_CHECK_CONSTRAINT_DUP_NAME, got success")
		} else if !strings.Contains(oracleErr.Error(), "3822") &&
			!strings.Contains(strings.ToLower(oracleErr.Error()), "duplicate") {
			t.Errorf("PS.8 oracle: unexpected error text: %v", oracleErr)
		}

		var omniErr error
		results, parseErr := c.Exec(ddl2, nil)
		if parseErr != nil {
			omniErr = parseErr
		} else {
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
					break
				}
			}
		}
		if omniErr == nil {
			t.Errorf("PS.8 omni: expected ER_CHECK_CONSTRAINT_DUP_NAME, got success")
		}
	})
}

// omniCheckNames returns the CHECK constraint names from the omni catalog
// for the given table (in testdb), sorted.
func psCheckNames(c *Catalog, table string) []string {
	db := c.GetDatabase("testdb")
	if db == nil {
		return nil
	}
	tbl := db.GetTable(table)
	if tbl == nil {
		return nil
	}
	var names []string
	for _, con := range tbl.Constraints {
		if con.Type == ConCheck {
			names = append(names, con.Name)
		}
	}
	sort.Strings(names)
	return names
}

// omniFKNames returns the FOREIGN KEY constraint names from the omni catalog
// for the given table (in testdb), sorted.
func psFKNames(c *Catalog, table string) []string {
	db := c.GetDatabase("testdb")
	if db == nil {
		return nil
	}
	tbl := db.GetTable(table)
	if tbl == nil {
		return nil
	}
	var names []string
	for _, con := range tbl.Constraints {
		if con.Type == ConForeignKey {
			names = append(names, con.Name)
		}
	}
	sort.Strings(names)
	return names
}
