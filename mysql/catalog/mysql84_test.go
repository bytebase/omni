package catalog

import (
	"errors"
	"testing"
)

func execExpectCode(t *testing.T, c *Catalog, sql string, wantCode int) {
	t.Helper()
	results, err := c.Exec(sql, nil)
	if err != nil {
		if wantCode == 0 {
			t.Fatalf("Exec(%q) unexpected parse/exec error: %v", sql, err)
		}
		var ce *Error
		if errors.As(err, &ce) && ce.Code == wantCode {
			return
		}
		t.Fatalf("Exec(%q) error = %v, want code %d", sql, err, wantCode)
	}
	var gotErr error
	for _, r := range results {
		if r.Error != nil {
			gotErr = r.Error
			break
		}
	}
	if wantCode == 0 {
		if gotErr != nil {
			t.Fatalf("Exec(%q) unexpected error: %v", sql, gotErr)
		}
		return
	}
	var ce *Error
	if !errors.As(gotErr, &ce) || ce.Code != wantCode {
		t.Fatalf("Exec(%q) error = %v, want code %d", sql, gotErr, wantCode)
	}
}

func TestMySQL84_RestrictFKOnNonStandardKey(t *testing.T) {
	setup := `
		CREATE DATABASE d;
		USE d;
		CREATE TABLE p (
			b INT PRIMARY KEY,
			c INT,
			d INT,
			e INT,
			f INT,
			s VARCHAR(50),
			KEY idx_d (d),
			UNIQUE KEY uk_ce (c, e),
			UNIQUE KEY uk_s (s(10))
		);
	`

	t.Run("MySQL84 rejects non-unique, unindexed, partial, reordered, and prefix unique parent keys with 6125", func(t *testing.T) {
		cases := []struct {
			name string
			stmt string
			code int
		}{
			{"non-unique parent key", "CREATE TABLE c1 (id INT PRIMARY KEY, d INT, FOREIGN KEY (d) REFERENCES p (d))", ErrFKNoUniqueIndexParent},
			{"unindexed parent column", "CREATE TABLE c2 (id INT PRIMARY KEY, f INT, FOREIGN KEY (f) REFERENCES p (f))", ErrFKNoUniqueIndexParent},
			{"partial prefix of composite unique key", "CREATE TABLE c3 (id INT PRIMARY KEY, c INT, FOREIGN KEY (c) REFERENCES p (c))", ErrFKNoUniqueIndexParent},
			{"reordered composite unique key", "CREATE TABLE c4 (id INT PRIMARY KEY, e INT, c INT, FOREIGN KEY (e, c) REFERENCES p (e, c))", ErrFKNoUniqueIndexParent},
			{"prefix-length unique key", "CREATE TABLE c5 (id INT PRIMARY KEY, s VARCHAR(50), FOREIGN KEY (s) REFERENCES p (s))", ErrFKNoUniqueIndexParent},
			{"alter add fk non-unique", "CREATE TABLE c6 (id INT PRIMARY KEY, d INT); ALTER TABLE c6 ADD FOREIGN KEY (d) REFERENCES p (d)", ErrFKNoUniqueIndexParent},
			{"exact primary key accepted", "CREATE TABLE c7 (id INT PRIMARY KEY, b INT, FOREIGN KEY (b) REFERENCES p (b))", 0},
			{"exact composite unique key accepted", "CREATE TABLE c8 (id INT PRIMARY KEY, c INT, e INT, FOREIGN KEY (c, e) REFERENCES p (c, e))", 0},
			{"foreign_key_checks=0 still rejects existing non-unique parent", "SET foreign_key_checks = 0; CREATE TABLE c9 (id INT PRIMARY KEY, d INT, FOREIGN KEY (d) REFERENCES p (d))", ErrFKNoUniqueIndexParent},
			{"foreign_key_checks=0 allows forward ref to missing table", "SET foreign_key_checks = 0; CREATE TABLE c10 (id INT PRIMARY KEY, x INT, FOREIGN KEY (x) REFERENCES nosuch (x))", 0},
			{"restrict_fk_on_non_standard_key=OFF allows non-unique parent", "SET SESSION restrict_fk_on_non_standard_key = OFF; CREATE TABLE c11 (id INT PRIMARY KEY, d INT, FOREIGN KEY (d) REFERENCES p (d))", 0},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				c, err := LoadSQLWithVersion(setup, MySQL84)
				if err != nil {
					t.Fatalf("setup: %v", err)
				}
				execExpectCode(t, c, tc.stmt, tc.code)
			})
		}
	})

	t.Run("MySQL80 keeps lenient non-unique and partial parent keys", func(t *testing.T) {
		c, err := LoadSQLWithVersion(setup, MySQL80)
		if err != nil {
			t.Fatalf("setup: %v", err)
		}
		execExpectCode(t, c, "CREATE TABLE c1 (id INT PRIMARY KEY, d INT, FOREIGN KEY (d) REFERENCES p (d))", 0)
		execExpectCode(t, c, "CREATE TABLE c2 (id INT PRIMARY KEY, c INT, FOREIGN KEY (c) REFERENCES p (c))", 0)
		execExpectCode(t, c, "CREATE TABLE c3 (id INT PRIMARY KEY, f INT, FOREIGN KEY (f) REFERENCES p (f))", ErrFKMissingIndex)
	})
}

func TestMySQL84_DropUniqueIndexNeededByFK(t *testing.T) {
	setup := `
		CREATE DATABASE d;
		USE d;
		CREATE TABLE p (id INT PRIMARY KEY, u INT, UNIQUE KEY uk_u (u), KEY k_u (u));
		CREATE TABLE c (id INT PRIMARY KEY, u INT, FOREIGN KEY (u) REFERENCES p (u));
	`

	// On MySQL 8.4, dropping uk_u leaves only the non-unique k_u, so InnoDB rejects it with 1553.
	c84, err := LoadSQLWithVersion(setup, MySQL84)
	if err != nil {
		t.Fatalf("setup 8.4: %v", err)
	}
	execExpectCode(t, c84, "ALTER TABLE p DROP INDEX uk_u", ErrDropIndexNeededInFK)

	// Adding a second unique key on (u) allows dropping uk_u on 8.4.
	execExpectCode(t, c84, "ALTER TABLE p ADD UNIQUE KEY uk_u2 (u)", 0)
	execExpectCode(t, c84, "ALTER TABLE p DROP INDEX uk_u", 0)

	// On MySQL 8.0, non-unique k_u still covers the FK, so dropping uk_u succeeds.
	c80, err := LoadSQLWithVersion(setup, MySQL80)
	if err != nil {
		t.Fatalf("setup 8.0: %v", err)
	}
	execExpectCode(t, c80, "ALTER TABLE p DROP INDEX uk_u", 0)
}

func TestMySQL84_FloatDoubleAutoIncrementRemoved(t *testing.T) {
	stmts := []struct {
		name string
		sql  string
	}{
		{"create float auto_increment", "CREATE TABLE f1 (id FLOAT AUTO_INCREMENT PRIMARY KEY)"},
		{"create double auto_increment", "CREATE TABLE f2 (id DOUBLE AUTO_INCREMENT PRIMARY KEY)"},
		{"create real auto_increment", "CREATE TABLE f3 (id REAL AUTO_INCREMENT PRIMARY KEY)"},
		{"alter modify float auto_increment", "CREATE TABLE f4 (id INT PRIMARY KEY); ALTER TABLE f4 MODIFY id FLOAT AUTO_INCREMENT"},
		{"alter modify double auto_increment", "CREATE TABLE f5 (id INT PRIMARY KEY); ALTER TABLE f5 MODIFY id DOUBLE AUTO_INCREMENT"},
	}

	for _, tc := range stmts {
		t.Run("8.4/"+tc.name, func(t *testing.T) {
			c, err := LoadSQLWithVersion("CREATE DATABASE d; USE d;", MySQL84)
			if err != nil {
				t.Fatal(err)
			}
			execExpectCode(t, c, tc.sql, ErrWrongFieldSpec)
		})
		t.Run("8.0/"+tc.name, func(t *testing.T) {
			c, err := LoadSQLWithVersion("CREATE DATABASE d; USE d;", MySQL80)
			if err != nil {
				t.Fatal(err)
			}
			execExpectCode(t, c, tc.sql, 0)
		})
	}
}

func TestMySQL84_NormalizerAndVersion(t *testing.T) {
	c := New()
	if c.Version() != MySQL80 {
		t.Fatalf("default Version = %v, want MySQL80", c.Version())
	}
	c.SetVersion(MySQL84)
	if c.Version() != MySQL84 {
		t.Fatalf("Version after SetVersion(MySQL84) = %v, want MySQL84", c.Version())
	}
	if !c.session.RestrictFKOnNonStandardKey {
		t.Fatalf("expected RestrictFKOnNonStandardKey=true on MySQL84")
	}
	n := NormalizerFor(MySQL84)
	if n.Version != MySQL84 || !n.ExplicitDefaultsForTimestamp {
		t.Fatalf("NormalizerFor(MySQL84) = %+v", n)
	}
	if versionName(MySQL84) != "8.4" {
		t.Fatalf("versionName(MySQL84) = %q, want 8.4", versionName(MySQL84))
	}
}
