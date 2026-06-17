package catalog

import (
	"slices"
	"sort"
	"strings"
	"testing"
)

// TestScenario_PS covers section PS of SCENARIOS-mysql-implicit-behavior.md
// — "Path-split behaviors (CREATE vs ALTER)". Each subtest runs the scenario's
// DDL against both a real MySQL 8.0 container and the omni catalog and asserts
// both match the expected value. Existing TestBugFix_* tests remain unchanged;
// these TestScenario_PS tests are the durable dual-assertion versions.
func TestScenario_PS(t *testing.T) {
	scenariosSkipIfShort(t)
	scenariosSkipIfNoDocker(t)

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

		want := []string{"t_chk_1", "t_chk_5"}

		rows := oracleRows(t, mc, `
            SELECT CONSTRAINT_NAME FROM information_schema.CHECK_CONSTRAINTS
            WHERE CONSTRAINT_SCHEMA='testdb'
            ORDER BY CONSTRAINT_NAME`)
		var oracleNames []string
		for _, r := range rows {
			oracleNames = append(oracleNames, asString(r[0]))
		}
		if !slices.Equal(oracleNames, want) {
			t.Errorf("PS.1 oracle CHECK names: got %v, want %v", oracleNames, want)
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

		want := []string{"t_chk_20", "t_chk_21"}

		rows := oracleRows(t, mc, `
            SELECT CONSTRAINT_NAME FROM information_schema.CHECK_CONSTRAINTS
            WHERE CONSTRAINT_SCHEMA='testdb'
            ORDER BY CONSTRAINT_NAME`)
		var oracleNames []string
		for _, r := range rows {
			oracleNames = append(oracleNames, asString(r[0]))
		}
		if !slices.Equal(oracleNames, want) {
			t.Errorf("PS.2 oracle CHECK names: got %v, want %v", oracleNames, want)
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

		want := []string{"child_ibfk_1", "child_ibfk_5"}

		rows := oracleRows(t, mc, `
            SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='child'
              AND CONSTRAINT_TYPE='FOREIGN KEY'
            ORDER BY CONSTRAINT_NAME`)
		var oracleNames []string
		for _, r := range rows {
			oracleNames = append(oracleNames, asString(r[0]))
		}
		if !slices.Equal(oracleNames, want) {
			t.Errorf("PS.3 oracle FK names: got %v, want %v", oracleNames, want)
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

		want := []string{"child_ibfk_20", "child_ibfk_21"}

		rows := oracleRows(t, mc, `
            SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='child'
              AND CONSTRAINT_TYPE='FOREIGN KEY'
            ORDER BY CONSTRAINT_NAME`)
		var oracleNames []string
		for _, r := range rows {
			oracleNames = append(oracleNames, asString(r[0]))
		}
		if !slices.Equal(oracleNames, want) {
			t.Errorf("PS.4 oracle FK names: got %v, want %v", oracleNames, want)
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

	// PS.6 HASH partition ADD — seeded from count.
	// Oracle and omni should both produce the expanded partition list.
	t.Run("PS_6_Hash_partition_ADD_seeded", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c,
			`CREATE TABLE t (id INT) PARTITION BY HASH(id) PARTITIONS 3`)
		// Run the ALTER on both sides and compare the resulting partition list.
		alter := `ALTER TABLE t ADD PARTITION PARTITIONS 2`
		if _, err := mc.db.ExecContext(mc.ctx, alter); err != nil {
			t.Errorf("PS.6 oracle ALTER failed: %v", err)
		}

		// Oracle partition names.
		rows := oracleRows(t, mc, `
            SELECT PARTITION_NAME FROM information_schema.PARTITIONS
            WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
            ORDER BY PARTITION_ORDINAL_POSITION`)
		var oracleNames []string
		for _, r := range rows {
			oracleNames = append(oracleNames, asString(r[0]))
		}
		wantOracle := []string{"p0", "p1", "p2", "p3", "p4"}
		if !slices.Equal(oracleNames, wantOracle) {
			t.Errorf("PS.6 oracle partition names: got %v, want %v", oracleNames, wantOracle)
		}

		// omni side: verify the expected 5-partition layout.
		results, parseErr := c.Exec(alter, nil)
		var omniAlterErr error
		if parseErr != nil {
			omniAlterErr = parseErr
		} else {
			for _, r := range results {
				if r.Error != nil {
					omniAlterErr = r.Error
					break
				}
			}
		}

		var omniNames []string
		if tbl := c.GetDatabase("testdb").GetTable("t"); tbl != nil && tbl.Partitioning != nil {
			for _, p := range tbl.Partitioning.Partitions {
				omniNames = append(omniNames, p.Name)
			}
		}
		if omniAlterErr != nil || !slices.Equal(omniNames, wantOracle) {
			t.Errorf("PS.6 omni: expected partition names %v; got %v (alter err: %v)", wantOracle, omniNames, omniAlterErr)
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

		_, oracleErr := mc.db.ExecContext(mc.ctx, ddl)
		if oracleErr == nil {
			t.Errorf("PS.7 oracle: expected ER_FK_DUP_NAME (1826), got success")
		} else if !strings.Contains(oracleErr.Error(), "1826") &&
			!strings.Contains(strings.ToLower(oracleErr.Error()), "duplicate") {
			t.Errorf("PS.7 oracle: unexpected error text: %v", oracleErr)
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
