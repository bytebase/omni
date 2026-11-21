package catalog

import (
	"errors"
	"strings"
	"testing"
)

func execSingleError(t *testing.T, setup, stmt string) *Error {
	t.Helper()
	c, err := LoadSQL("CREATE DATABASE test;\nUSE test;\n" + setup)
	if err != nil {
		t.Fatalf("LoadSQL(%q): %v", setup, err)
	}
	results, err := c.Exec(stmt, nil)
	if err != nil {
		var ce *Error
		if errors.As(err, &ce) {
			return ce
		}
		t.Fatalf("Exec(%q): non-catalog error %v", stmt, err)
	}
	for _, r := range results {
		if r.Error != nil {
			var ce *Error
			if errors.As(r.Error, &ce) {
				return ce
			}
			t.Fatalf("Exec(%q): non-catalog result error %v", stmt, r.Error)
		}
	}
	return nil
}

func TestViewObjectKindDDL(t *testing.T) {
	setup := `
		CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
		CREATE VIEW ident AS SELECT b, ident, d FROM t;
		CREATE DATABASE other;
	`
	tests := []struct {
		name     string
		stmt     string
		wantCode int
	}{
		{"alter view add column", "ALTER TABLE ident ADD COLUMN e INT", ErrWrongObject},
		{"alter view comment", "ALTER TABLE ident COMMENT 'x'", ErrWrongObject},
		{"alter view compact succeeds", "ALTER TABLE ident COMPACT", 0},
		{"alter view compact tiflash succeeds", "ALTER TABLE ident COMPACT TIFLASH REPLICA", 0},
		{"alter view analyze partition", "ALTER TABLE ident ANALYZE PARTITION ident", ErrUnknownError},
		{"alter view analyze partition index", "ALTER TABLE ident ANALYZE PARTITION ident INDEX", ErrPartitionMgmtOnNonPartitioned},
		{"alter view analyze partition named index", "ALTER TABLE ident ANALYZE PARTITION ident INDEX zzz", ErrPartitionMgmtOnNonPartitioned},
		{"create table like view", "CREATE TABLE nt LIKE ident", ErrWrongObject},
		{"create table like view colliding target", "CREATE TABLE ident LIKE ident", ErrWrongObject},
		{"create table paren like view", "CREATE TABLE ident (LIKE ident)", ErrWrongObject},
		{"rename view onto self", "RENAME TABLE ident TO ident", ErrDupTable},
		{"rename view onto table", "RENAME TABLE ident TO t", ErrDupTable},
		{"rename table onto view", "RENAME TABLE t TO ident", ErrDupTable},
		{"alter table rename onto view", "ALTER TABLE t RENAME TO ident", ErrDupTable},
		{"rename view same db", "RENAME TABLE ident TO ident2", 0},
		{"rename view across db", "RENAME TABLE ident TO other.ident2", ErrRenameViewAcrossDB},
		{"create index on view", "CREATE INDEX i ON ident (ident)", ErrUnknownError},
		{"create unique index on view", "CREATE UNIQUE INDEX i ON ident (ident)", ErrUnknownError},
		{"create fulltext index on view", "CREATE FULLTEXT INDEX i ON ident (ident)", ErrUnsupportedDDL},
		{"create spatial index on view", "CREATE SPATIAL INDEX i ON ident (ident)", ErrUnsupportedDDL},
		{"drop index on view", "DROP INDEX ident ON ident", ErrCantDropKey},
		{"drop index if exists on view", "DROP INDEX IF EXISTS ident ON ident", 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ce := execSingleError(t, setup, tc.stmt)
			gotCode := 0
			if ce != nil {
				gotCode = ce.Code
			}
			if gotCode != tc.wantCode {
				t.Errorf("%q: got code %d (%v), want %d", tc.stmt, gotCode, ce, tc.wantCode)
			}
		})
	}

	// Verify RENAME TABLE actually moves a view and preserves SHOW CREATE VIEW.
	c, err := LoadSQL("CREATE DATABASE test; USE test; " + setup)
	if err != nil {
		t.Fatal(err)
	}
	if results, err := c.Exec("RENAME TABLE ident TO ident_renamed", nil); err != nil || results[0].Error != nil {
		t.Fatalf("RENAME TABLE view failed: err=%v results=%v", err, results)
	}
	if c.ShowCreateView("test", "ident") != "" {
		t.Errorf("old view name still exists after RENAME TABLE")
	}
	if got := c.ShowCreateView("test", "ident_renamed"); !strings.Contains(got, "VIEW `ident_renamed`") {
		t.Errorf("unexpected ShowCreateView after rename: %q", got)
	}
}

func TestPartitionedTableAlterValidation(t *testing.T) {
	setup := `
		CREATE TABLE ident (
			b INT,
			ident INT,
			d INT,
			PRIMARY KEY (b, ident),
			KEY idx_d (d)
		) PARTITION BY RANGE (ident) (
			PARTITION p0 VALUES LESS THAN (100),
			PARTITION ident VALUES LESS THAN (200)
		);
		CREATE TABLE swap_plain (
			b INT,
			ident INT,
			d INT,
			PRIMARY KEY (b, ident),
			KEY idx_d (d)
		);
	`
	tests := []struct {
		name     string
		stmt     string
		wantCode int
	}{
		{"add partition values in on range", "ALTER TABLE ident ADD PARTITION (PARTITION p2 VALUES IN (1))", ErrWrongPartitionTypeValues},
		{"add partition default on range", "ALTER TABLE ident ADD PARTITION (PARTITION p2 DEFAULT)", ErrWrongPartitionTypeValues},
		{"add partition without values on range", "ALTER TABLE ident ADD PARTITION (PARTITION p2)", ErrPartitionRequiresValues},
		{"reorganize without values on range", "ALTER TABLE ident REORGANIZE PARTITION ident INTO (PARTITION p2)", ErrPartitionRequiresValues},
		{"add partition bare on range", "ALTER TABLE ident ADD PARTITION", ErrPartitionsMustBeDefined},
		{"add partition count on range", "ALTER TABLE ident ADD PARTITION PARTITIONS 1", ErrPartitionsMustBeDefined},
		{"add partition duplicate name", "ALTER TABLE ident ADD PARTITION (PARTITION ident VALUES LESS THAN MAXVALUE)", ErrSameNamePartition},
		{"add partition if not exists duplicate", "ALTER TABLE ident ADD PARTITION IF NOT EXISTS (PARTITION ident VALUES LESS THAN MAXVALUE)", 0},
		{"add partition duplicate in same stmt", "ALTER TABLE ident ADD PARTITION (PARTITION p2 VALUES LESS THAN (300), PARTITION p2 VALUES LESS THAN (400))", ErrSameNamePartition},
		{"add partition too many range values", "ALTER TABLE ident ADD PARTITION (PARTITION p2 VALUES LESS THAN (MAXVALUE, MAXVALUE))", ErrTooManyValues},
		{"add partition column ref bound", "ALTER TABLE ident ADD PARTITION (PARTITION p2 VALUES LESS THAN (ident))", ErrNoSuchColumn},
		{"first partition column ref bound", "ALTER TABLE ident FIRST PARTITION LESS THAN (ident)", ErrNoSuchColumn},
		{"last partition column ref bound", "ALTER TABLE ident LAST PARTITION LESS THAN (ident)", ErrNoSuchColumn},
		{"coalesce on range", "ALTER TABLE ident COALESCE PARTITION 1", ErrCoalesceOnlyOnHashPartition},
		{"exchange with partitioned table", "ALTER TABLE ident EXCHANGE PARTITION ident WITH TABLE ident", ErrPartitionExchangePartTable},
		{"exchange with plain table", "ALTER TABLE ident EXCHANGE PARTITION ident WITH TABLE swap_plain", 0},
		{"analyze partition missing index", "ALTER TABLE ident ANALYZE PARTITION ident INDEX zzz", ErrAnalyzeIndexNotFound},
		{"analyze partition existing index", "ALTER TABLE ident ANALYZE PARTITION ident INDEX idx_d", 0},
		{"analyze partition primary index", "ALTER TABLE ident ANALYZE PARTITION ident INDEX PRIMARY", 0},
		{"add valid range partition", "ALTER TABLE ident ADD PARTITION (PARTITION p2 VALUES LESS THAN (300))", 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ce := execSingleError(t, setup, tc.stmt)
			gotCode := 0
			if ce != nil {
				gotCode = ce.Code
			}
			if gotCode != tc.wantCode {
				t.Errorf("%q: got code %d (%v), want %d", tc.stmt, gotCode, ce, tc.wantCode)
			}
		})
	}
}

func TestForeignKeyTableAlterValidation(t *testing.T) {
	setup := `
		CREATE TABLE parent (b INT PRIMARY KEY);
		CREATE TABLE ident (
			b INT PRIMARY KEY,
			ident INT,
			d INT,
			KEY idx_d (d),
			CONSTRAINT ident FOREIGN KEY (ident) REFERENCES parent (b)
		);
	`
	tests := []struct {
		name     string
		stmt     string
		wantCode int
	}{
		{"alter drop fk backing index", "ALTER TABLE ident DROP INDEX ident", ErrDropIndexNeededInForeignKey},
		{"alter drop fk backing key", "ALTER TABLE ident DROP KEY ident", ErrDropIndexNeededInForeignKey},
		{"standalone drop fk backing index", "DROP INDEX ident ON ident", ErrDropIndexNeededInForeignKey},
		{"modify fk child column incompatible type", "ALTER TABLE ident MODIFY ident TINYINT", ErrFKIncompatibleColumns},
		{"change fk child column incompatible type", "ALTER TABLE ident CHANGE ident ident TINYINT", ErrFKIncompatibleColumns},
		{"modify fk parent column incompatible type", "ALTER TABLE parent MODIFY b BIGINT", ErrFKIncompatibleColumns},
		{"modify fk child column compatible type", "ALTER TABLE ident MODIFY ident INT NOT NULL", 0},
		{"drop check does not error or drop fk", "ALTER TABLE ident DROP CHECK ident", 0},
		{"drop constraint does not error or drop fk", "ALTER TABLE ident DROP CONSTRAINT ident", 0},
		{"drop fk succeeds", "ALTER TABLE ident DROP FOREIGN KEY ident", 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ce := execSingleError(t, setup, tc.stmt)
			gotCode := 0
			if ce != nil {
				gotCode = ce.Code
			}
			if gotCode != tc.wantCode {
				t.Errorf("%q: got code %d (%v), want %d", tc.stmt, gotCode, ce, tc.wantCode)
			}
		})
	}

	// Verify DROP CHECK / DROP CONSTRAINT leaves the identically named FK intact.
	for _, stmt := range []string{"ALTER TABLE ident DROP CHECK ident", "ALTER TABLE ident DROP CONSTRAINT ident"} {
		c, err := LoadSQL("CREATE DATABASE test; USE test; " + setup)
		if err != nil {
			t.Fatal(err)
		}
		if results, err := c.Exec(stmt, nil); err != nil || results[0].Error != nil {
			t.Fatalf("%q failed: err=%v results=%v", stmt, err, results)
		}
		ddl := c.ShowCreateTable("test", "ident")
		if !strings.Contains(ddl, "CONSTRAINT `ident` FOREIGN KEY") {
			t.Errorf("%q wrongly removed FK:\n%s", stmt, ddl)
		}
	}

	// Verify dropping FK backing index succeeds when another covering index exists.
	c, err := LoadSQL("CREATE DATABASE test; USE test; " + setup)
	if err != nil {
		t.Fatal(err)
	}
	if results, err := c.Exec("ALTER TABLE ident ADD INDEX ident_cover (ident, d); ALTER TABLE ident DROP INDEX ident;", nil); err != nil {
		t.Fatalf("Exec failed: %v", err)
	} else {
		for _, r := range results {
			if r.Error != nil {
				t.Fatalf("unexpected error dropping redundant FK index: %v", r.Error)
			}
		}
	}
	ddl := c.ShowCreateTable("test", "ident")
	if !strings.Contains(ddl, "CONSTRAINT `ident` FOREIGN KEY") {
		t.Errorf("dropping redundant index removed FK:\n%s", ddl)
	}
}
