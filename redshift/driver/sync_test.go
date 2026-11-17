package driver

import "testing"

// TestGetPgVersionAndRedshiftVersion reads what SELECT version() reports, which
// is how a sync learns both the PostgreSQL release Redshift speaks and the
// Redshift release itself.
func TestGetPgVersionAndRedshiftVersion(t *testing.T) {
	for _, tc := range []struct {
		name         string
		version      string
		wantPg       string
		wantRedshift string
		wantErr      bool
	}{
		{
			name:         "the documented form",
			version:      "PostgreSQL 8.0.2 on i686-pc-linux-gnu, compiled by GCC gcc (GCC) 3.4.2 20041017 (Red Hat 3.4.2-6.fc3), Redshift 1.0.48042",
			wantPg:       "8.0.2",
			wantRedshift: "1.0.48042",
		},
		{
			name:         "another build",
			version:      "PostgreSQL 8.0.2 on x86_64-pc-linux-gnu, compiled by gcc (GCC) 3.4.2 20041017 (Red Hat 3.4.2-6.fc3), Redshift 1.0.12345",
			wantPg:       "8.0.2",
			wantRedshift: "1.0.12345",
		},
		{
			name:    "a string that is not a Redshift version",
			version: "Redshift 1.0",
			wantErr: true,
		},
		{
			name:    "nothing",
			wantErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pgVersion, redshiftVersion, err := getPgVersionAndRedshiftVersion(tc.version)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("parsed %q into %q and %q", tc.version, pgVersion, redshiftVersion)
				}
				return
			}
			if err != nil {
				t.Fatalf("getPgVersionAndRedshiftVersion: %v", err)
			}
			if pgVersion != tc.wantPg {
				t.Errorf("postgres version = %q, want %q", pgVersion, tc.wantPg)
			}
			if redshiftVersion != tc.wantRedshift {
				t.Errorf("redshift version = %q, want %q", redshiftVersion, tc.wantRedshift)
			}
		})
	}
}

// TestBuildRedshiftVersionString pins the version the instance list shows.
func TestBuildRedshiftVersionString(t *testing.T) {
	got := buildRedshiftVersionString("1.0.48042", "8.0.2")
	want := "Redshift 1.0.48042 based on PostgreSQL 8.0.2"
	if got != want {
		t.Errorf("buildRedshiftVersionString = %q, want %q", got, want)
	}
}

// TestExcludedDatabaseList pins the databases a sync leaves out: a cluster
// answers them in pg_database but they are not the user's.
func TestExcludedDatabaseList(t *testing.T) {
	for _, name := range []string{"padb_harvest", "awsdatacatalog", "sys:internal", "template0", "template1"} {
		if !excludedDatabaseList[name] {
			t.Errorf("%q is not excluded", name)
		}
	}
	for _, name := range []string{"dev", "analytics", "prod"} {
		if excludedDatabaseList[name] {
			t.Errorf("%q is excluded", name)
		}
	}
}

// TestForeignKeyCodesAreSpelledOut pins the one-character codes pg_constraint
// writes into the words a synced foreign key carries.
func TestForeignKeyCodesAreSpelledOut(t *testing.T) {
	for code, want := range map[string]string{
		"a": "NO ACTION",
		"r": "RESTRICT",
		"c": "CASCADE",
		"n": "SET NULL",
		"d": "SET DEFAULT",
		"x": "x",
	} {
		if got := convertForeignKeyActionCode(code); got != want {
			t.Errorf("convertForeignKeyActionCode(%q) = %q, want %q", code, got, want)
		}
	}
	for code, want := range map[string]string{
		"f": "FULL",
		"p": "PARTIAL",
		"s": "SIMPLE",
		"x": "x",
	} {
		if got := convertForeignKeyMatchType(code); got != want {
			t.Errorf("convertForeignKeyMatchType(%q) = %q, want %q", code, got, want)
		}
	}
}

// TestForeignKeyColumnsFromDefinition reads the column lists out of the
// constraint text the sync gets back.
func TestForeignKeyColumnsFromDefinition(t *testing.T) {
	columns, referenced, err := getForeignKeyColumnsAndReferencedColumns(`FOREIGN KEY (a, "B b") REFERENCES other.t(x, y)`)
	if err != nil {
		t.Fatalf("getForeignKeyColumnsAndReferencedColumns: %v", err)
	}
	for i, want := range []string{"a", "B b"} {
		if columns[i] != want {
			t.Errorf("column %d = %q, want %q", i, columns[i], want)
		}
	}
	for i, want := range []string{"x", "y"} {
		if referenced[i] != want {
			t.Errorf("referenced column %d = %q, want %q", i, referenced[i], want)
		}
	}

	if _, _, err := getForeignKeyColumnsAndReferencedColumns("FOREIGN KEY (a) REFERENCES t"); err == nil {
		t.Error("a definition without the referenced columns was accepted")
	}
}

// TestFormatTableNameFromRegclass reduces what
// quote_ident(schema) || '.' || quote_ident(table) writes.
func TestFormatTableNameFromRegclass(t *testing.T) {
	for in, want := range map[string]string{
		`public.users`:     "users",
		`"Weird Schema".t`: "t",
		`users`:            "users",
	} {
		if got := formatTableNameFromRegclass(in); got != want {
			t.Errorf("formatTableNameFromRegclass(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestGetIndexMethodType reads the access method an index definition names, and
// Redshift's default when it names none.
func TestGetIndexMethodType(t *testing.T) {
	for definition, want := range map[string]string{
		"CREATE INDEX i ON t USING btree (a)": "btree",
		"CREATE INDEX i ON t USING gin (a)":   "gin",
		"CREATE INDEX i ON t (a)":             "btree",
	} {
		if got := getIndexMethodType(definition); got != want {
			t.Errorf("getIndexMethodType(%q) = %q, want %q", definition, got, want)
		}
	}
}

// TestConvertYesNo pins the nullable answers a sync reads.
func TestConvertYesNo(t *testing.T) {
	for in, want := range map[string]bool{"YES": true, "Y": true, "1": true, "NO": false, "N": false, "0": false} {
		got, err := convertYesNo(in)
		if err != nil {
			t.Errorf("convertYesNo(%q): %v", in, err)
			continue
		}
		if got != want {
			t.Errorf("convertYesNo(%q) = %v, want %v", in, got, want)
		}
	}
	if _, err := convertYesNo("maybe"); err == nil {
		t.Error("convertYesNo accepted an answer that is neither")
	}
}

// TestGetDatabaseInCreateDatabaseStatement reads the database name a
// CREATE DATABASE step creates, in either case, and the forms that name none.
func TestGetDatabaseInCreateDatabaseStatement(t *testing.T) {
	for statement, want := range map[string]string{
		"CREATE DATABASE analytics":        "analytics",
		"create database analytics;":       "analytics",
		`CREATE DATABASE "Analytics Dev";`: "Analytics Dev",
	} {
		got, err := getDatabaseInCreateDatabaseStatement(statement)
		if err != nil {
			t.Errorf("getDatabaseInCreateDatabaseStatement(%q): %v", statement, err)
			continue
		}
		if got != want {
			t.Errorf("getDatabaseInCreateDatabaseStatement(%q) = %q, want %q", statement, got, want)
		}
	}

	for _, statement := range []string{"", "SELECT 1", "CREATE TABLE t (a int)"} {
		if _, err := getDatabaseInCreateDatabaseStatement(statement); err == nil {
			t.Errorf("getDatabaseInCreateDatabaseStatement(%q) found a name", statement)
		}
	}
}
