package catalog

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

// OceanBase names a CHECK or a foreign key written without a name after the table and
// the time, `t_OBCHECK_1794529978914345`, where the loader has a count for the time. These pin
// what the diff makes of two constraints that differ only in such names; no server is
// needed, and TestOracle_CheckDiffIdempotence and TestOracle_ForeignKeyDiffIdempotence
// hold the same to the engine's own readbacks.

// loaderCount is a name the loader made up, told from the engine's here by its zeros.
var loaderCount = regexp.MustCompile(`(_OB(?:FK|CHECK)_)0{8}\d{8}`)

func madeUpNameDiff(t testing.TB, from, to string) string {
	t.Helper()
	d := Diff(loadCat(t, dbDDL+from), loadCat(t, dbDDL+to))
	var parts []string
	for _, te := range d.Tables {
		for _, c := range te.Checks {
			parts = append(parts, fmt.Sprintf("check %s %s", c.Action, loaderCount.ReplaceAllString(c.Name, "${1}<n>")))
		}
		for _, f := range te.ForeignKeys {
			parts = append(parts, fmt.Sprintf("fk %s %s", f.Action, loaderCount.ReplaceAllString(f.Name, "${1}<n>")))
		}
		if len(te.Checks)+len(te.ForeignKeys) == 0 {
			parts = append(parts, fmt.Sprintf("table %s %s", te.Action, te.Name))
		}
	}
	return strings.Join(parts, "; ")
}

func TestDiff_ChecksUnderMadeUpNames(t *testing.T) {
	const stored = "CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT `t_OBCHECK_1794529978914345` CHECK ((`a` > 0)));"
	tests := []struct {
		name     string
		from, to string
		want     string
	}{
		{"an unnamed CHECK is its own readback",
			stored, "CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a > 0));", ""},
		{"and the other way round",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a > 0));", stored, ""},
		{"a name made up under the table's old name",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT `old_OBCHECK_1794529978914345` CHECK ((`a` > 0)));",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a > 0));", ""},
		{"another expression is another CHECK",
			stored, "CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a > 1));",
			"check ADD t_OBCHECK_<n>; check DROP t_OBCHECK_1794529978914345"},
		{"NOT ENFORCED is part of the definition",
			stored, "CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a > 0) NOT ENFORCED);",
			"check ADD t_OBCHECK_<n>; check DROP t_OBCHECK_1794529978914345"},
		{"naming an unnamed CHECK is a change",
			stored, "CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT positive CHECK (a > 0));",
			"check ADD positive; check DROP t_OBCHECK_1794529978914345"},
		{"a name somebody wrote is not paired by its definition",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT positive CHECK (a > 0));",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a > 0));",
			"check DROP positive; check ADD t_OBCHECK_<n>"},
		{"two copies against one leave one over",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT `t_OBCHECK_1794529978914345` CHECK ((`a` > 0)), CONSTRAINT `t_OBCHECK_1794529978914346` CHECK ((`a` > 0)));",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a > 0));",
			"check DROP t_OBCHECK_1794529978914346"},
		{"each unnamed CHECK finds its own",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT `t_OBCHECK_1794529978914345` CHECK ((`b` < 100)), CONSTRAINT `t_OBCHECK_1794529978914346` CHECK ((`a` > 0)));",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (a > 0), CHECK (b < 100));", ""},
		{"a made-up name both sides have, over two definitions, is two CHECKs",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT `t_OBCHECK_0000000000000001` CHECK ((`a` > 0)));",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a > 5));",
			"check MODIFY t_OBCHECK_<n>"},
		{"made-up names both sides have, the other way round, are the same CHECKs",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT `t_OBCHECK_0000000000000001` CHECK ((`b` < 100)), CONSTRAINT `t_OBCHECK_0000000000000002` CHECK ((`a` > 0)));",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (a > 0), CHECK (b < 100));", ""},
		{"a name of MySQL's shape is a name somebody wrote",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT t_chk_1 CHECK (a > 0));",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a > 0));",
			"check DROP t_chk_1; check ADD t_OBCHECK_<n>"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := madeUpNameDiff(t, tc.from, tc.to); got != tc.want {
				t.Errorf("diff = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestMigration_LeavesMadeUpNamesOut pins that a plan gives a constraint no name nobody
// wrote: OceanBase makes its own, and the loader's count never reaches the engine.
func TestMigration_LeavesMadeUpNamesOut(t *testing.T) {
	const from = dbDDL + "CREATE TABLE p (id INT PRIMARY KEY);\nCREATE TABLE c (id INT PRIMARY KEY, pid INT);"
	for name, tc := range map[string]struct{ to, want string }{
		"unnamed": {
			"CREATE TABLE p (id INT PRIMARY KEY);\nCREATE TABLE c (id INT PRIMARY KEY, pid INT, CHECK (pid > 0), FOREIGN KEY (pid) REFERENCES p (id));",
			"ALTER TABLE `app`.`c` ADD CHECK ((`pid` > 0));\nALTER TABLE `app`.`c` ADD FOREIGN KEY (`pid`) REFERENCES `p` (`id`)",
		},
		"named": {
			"CREATE TABLE p (id INT PRIMARY KEY);\nCREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT pos CHECK (pid > 0), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id));",
			"ALTER TABLE `app`.`c` ADD CONSTRAINT `pos` CHECK ((`pid` > 0));\nALTER TABLE `app`.`c` ADD CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `p` (`id`)",
		},
	} {
		t.Run(name, func(t *testing.T) {
			fromCat, toCat := loadCat(t, from), loadCat(t, dbDDL+tc.to)
			if got := GenerateMigration(fromCat, toCat, Diff(fromCat, toCat)).SQL(); got != tc.want {
				t.Errorf("plan:\n%s\nwant:\n%s", got, tc.want)
			}
		})
	}
}

func TestDiff_ForeignKeysUnderMadeUpNames(t *testing.T) {
	const parent = "CREATE TABLE p (id INT PRIMARY KEY);\n"
	const stored = parent + "CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT `c_OBFK_1794529609863994` FOREIGN KEY (`pid`) REFERENCES `app`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT);"
	tests := []struct {
		name     string
		from, to string
		want     string
	}{
		{"an unnamed foreign key is its own readback",
			stored, parent + "CREATE TABLE c (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES p (id));", ""},
		{"and the other way round",
			parent + "CREATE TABLE c (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES p (id));", stored, ""},
		{"another action is another key",
			stored, parent + "CREATE TABLE c (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES p (id) ON DELETE CASCADE);",
			"fk ADD c_OBFK_<n>; fk DROP c_OBFK_1794529609863994"},
		{"naming an unnamed foreign key is a change",
			stored, parent + "CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk_pid FOREIGN KEY (pid) REFERENCES p (id));",
			"fk DROP c_OBFK_1794529609863994; fk ADD fk_pid"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := madeUpNameDiff(t, tc.from, tc.to); got != tc.want {
				t.Errorf("diff = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestMadeUpConstraintNames(t *testing.T) {
	tbl := &Table{Name: "Orders"}
	for name, want := range map[string]bool{
		"orders_OBCHECK_0000000000000001": true, // the loader's, a count for the time
		"ORDERS_obcheck_12":               true,
		"orders_OBCHECK_1794529978914345": true,
		"old_OBCHECK_1794529978914345":    true,  // made under the table's old name
		"orders_chk_1":                    false, // MySQL's shape: somebody wrote it
		"orders_OBCHECK_1a":               false,
		"positive":                        false,
		"orders_OBCHECK_":                 false,
		"orders_ibfk_1":                   false, // a foreign key's name, on a CHECK
	} {
		if got := madeUpCheckName(&Constraint{Name: name, Table: tbl}); got != want {
			t.Errorf("madeUpCheckName(%q) = %v, want %v", name, got, want)
		}
	}
	for name, want := range map[string]bool{
		"orders_ibfk_1":                false, // MySQL's shape: somebody wrote it
		"orders_OBFK_1794529609863994": true,
		"x_obfk_7":                     true,
		"fk_orders":                    false,
		"orders_chk_1":                 false,
	} {
		if got := madeUpForeignKeyName(&Constraint{Name: name, Table: tbl}); got != want {
			t.Errorf("madeUpForeignKeyName(%q) = %v, want %v", name, got, want)
		}
	}
}

// A schema written by one catalog and loaded by another carries the names the first made
// up, as names somebody wrote. The second catalog's count starts over, so it has to pass
// over the names taken, or the first unnamed constraint it meets collides with one.
func TestMadeUpNamesPassOverNamesTaken(t *testing.T) {
	first := New()
	if code := omniCode(first, `CREATE DATABASE d; USE d;
		CREATE TABLE p (id INT PRIMARY KEY);
		CREATE TABLE t (a INT, b INT, CHECK (a > 0), FOREIGN KEY (b) REFERENCES p (id));`); code != 0 {
		t.Fatalf("setup: error %d", code)
	}
	second := New()
	if code := omniCode(second, "CREATE DATABASE d; USE d; "+first.ShowCreateTable("d", "p")+"; "+first.ShowCreateTable("d", "t")); code != 0 {
		t.Fatalf("reload: error %d\n%s", code, first.ShowCreateTable("d", "t"))
	}
	for _, stmt := range []string{
		"ALTER TABLE t ADD CHECK (a < 10)",
		"ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES p (id)",
		"CREATE TABLE u (a INT, CHECK (a > 0), FOREIGN KEY (a) REFERENCES p (id))",
	} {
		if code := omniCode(second, stmt); code != 0 {
			t.Errorf("%s: error %d", stmt, code)
		}
	}
	seen := map[string]bool{}
	for _, name := range []string{"t", "u"} {
		for _, con := range second.GetDatabase("d").GetTable(name).Constraints {
			if con.Type != ConCheck && con.Type != ConForeignKey {
				continue
			}
			if key := strings.ToLower(con.Name); seen[key] {
				t.Errorf("two constraints are named %s", con.Name)
			} else {
				seen[key] = true
			}
		}
	}
	if len(seen) != 6 {
		t.Errorf("want 6 constraints over t and u, got %d: %v", len(seen), seen)
	}
}

// A unique key over a column named PRIMARY is named PRIMARY by the engine and by the
// catalog, and that name cannot be written (1280), so a plan writes the key without one.
// The statement is the one TestScenario_C1 sends OceanBase in 1.10.
func TestMigration_LeavesTheNamePrimaryOut(t *testing.T) {
	const without = "CREATE TABLE t (`PRIMARY` INT);"
	const with = "CREATE TABLE t (`PRIMARY` INT, UNIQUE KEY (`PRIMARY`));"
	for name, from := range map[string]string{"a key added": dbDDL + without, "a table created": dbDDL} {
		t.Run(name, func(t *testing.T) {
			fromCat, toCat := loadCat(t, from), loadCat(t, dbDDL+with)
			plan := GenerateMigration(fromCat, toCat, Diff(fromCat, toCat)).SQL()
			if !strings.Contains(plan, "UNIQUE KEY (`PRIMARY`)") {
				t.Errorf("the plan does not write the key without its name:\n%s", plan)
			}
			// The plan runs on the catalog it was made for, and leaves the target.
			if code := omniCode(fromCat, "USE app; "+plan); code != 0 {
				t.Fatalf("the plan does not run: error %d\n%s", code, plan)
			}
			if d := Diff(fromCat, toCat); !d.IsEmpty() {
				t.Errorf("the plan leaves a difference: %+v", d.Tables)
			}
		})
	}
}
