package catalog

import (
	"context"
	"sort"
	"strings"
	"testing"

	"github.com/bytebase/omni/mariadb/internal/testcontainer"
)

// TestOracle_DatabaseMetadataRestoresOnServer executes a MariaDB schema on both
// the pinned MariaDB 11.8.8 server and omni's catalog, exports the catalog via
// DatabaseMetadata, renders both dump and SDL forms with DatabaseDefinition,
// restores each onto a fresh database on the server, and requires every stored
// object (tables, views with CHECK OPTION, sequences, routines, triggers,
// events) to match what the server stored from the original DDL.
func TestOracle_DatabaseMetadataRestoresOnServer(t *testing.T) {
	conn, err := testcontainer.OpenDatabase(t).Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	database := testcontainer.DatabaseName(t)

	schema := `
CREATE SEQUENCE seq_order START WITH 100 INCREMENT BY 5 MINVALUE 1 MAXVALUE 9999 CACHE 20 CYCLE COMMENT='order seq';
CREATE SEQUENCE seq_u32 AS INT UNSIGNED START WITH 1 NOCACHE NOCYCLE;

CREATE TABLE parent (
	id INT AUTO_INCREMENT PRIMARY KEY,
	code VARCHAR(20) NOT NULL DEFAULT 'abc',
	expr_str VARCHAR(20) NOT NULL DEFAULT '1 + 2',
	expr_col INT DEFAULT (1 + 2),
	fn_col DATE DEFAULT curdate(),
	bit_col BIT(3) DEFAULT b'101',
	nullable_col INT,
	ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	gen_stored VARCHAR(30) GENERATED ALWAYS AS (CONCAT(code, '_s')) STORED,
	gen_virt VARCHAR(30) GENERATED ALWAYS AS (CONCAT(code, '_v')) VIRTUAL,
	g POINT NOT NULL,
	secret INT INVISIBLE NOT NULL DEFAULT 0,
	SPATIAL INDEX sp_idx (g),
	INDEX idx_ignored (code) IGNORED COMMENT 'not used by optimizer',
	INDEX idx_prefix (code(10) DESC)
) ENGINE=InnoDB COMMENT='parent''s table';

CREATE TABLE child (
	id INT PRIMARY KEY,
	parent_id INT,
	score INT,
	CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id) ON DELETE CASCADE,
	CONSTRAINT chk_score CHECK (score >= 0)
) ENGINE=InnoDB;

CREATE TABLE part_tbl (
	id INT NOT NULL,
	yr INT NOT NULL,
	PRIMARY KEY (id, yr)
) ENGINE=InnoDB
PARTITION BY RANGE (yr)
SUBPARTITION BY HASH (id) SUBPARTITIONS 2 (
	PARTITION p0 VALUES LESS THAN (2025),
	PARTITION p1 VALUES LESS THAN MAXVALUE
);

CREATE TABLE versioned_explicit (
	id INT PRIMARY KEY,
	val VARCHAR(50),
	row_start TIMESTAMP(6) GENERATED ALWAYS AS ROW START,
	row_end TIMESTAMP(6) GENERATED ALWAYS AS ROW END,
	PERIOD FOR SYSTEM_TIME (row_start, row_end)
) ENGINE=InnoDB WITH SYSTEM VERSIONING;

CREATE VIEW v_cascaded (vid, vcode) AS SELECT id, code FROM parent WITH CASCADED CHECK OPTION;
CREATE VIEW v_local AS SELECT id, code FROM parent WITH LOCAL CHECK OPTION;

CREATE FUNCTION fn_add(a INT, b INT) RETURNS INT DETERMINISTIC COMMENT 'adds two ints'
BEGIN
	RETURN a + b;
END;

CREATE PROCEDURE proc_noop() COMMENT 'noop proc'
BEGIN
	SELECT 1;
END;

CREATE TRIGGER trg_parent AFTER INSERT ON parent FOR EACH ROW
BEGIN
	INSERT INTO child(id) VALUES (NEW.id);
END;

CREATE TRIGGER trg_z_first AFTER INSERT ON parent FOR EACH ROW PRECEDES trg_parent
BEGIN
	INSERT INTO child(id) VALUES (NEW.id + 1);
END;

CREATE EVENT ev_daily ON SCHEDULE EVERY 1 DAY COMMENT 'daily job' DO SELECT 1;
`

	// Apply original DDL directly on the pinned MariaDB server.
	resetDatabase(t, conn, database)
	applySchema(t, conn, schema)
	want := storedObjects(t, conn, false)

	// Simulate the same DDL in omni's catalog and export via DatabaseMetadata.
	c := newCatalogAt(MariaDB118)
	setup := "CREATE DATABASE " + quoteIdent(database) + ";\nUSE " + quoteIdent(database) + ";\n" + schema
	results, err := c.Exec(setup, nil)
	if err != nil {
		t.Fatalf("catalog Exec: %v", err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("catalog statement %d failed: %v", r.Index, r.Error)
		}
	}
	meta := c.DatabaseMetadata(database)

	compare := func(t *testing.T, label string, got map[string]string, ignoreTriggerOrder bool) {
		t.Helper()
		var names []string
		for name := range want {
			names = append(names, name)
		}
		for name := range got {
			if _, ok := want[name]; !ok {
				names = append(names, name)
			}
		}
		sort.Strings(names)
		for _, name := range names {
			w, g := want[name], got[name]
			if ignoreTriggerOrder && strings.HasPrefix(name, "trigger ") {
				w = stripTriggerActionOrder(w)
				g = stripTriggerActionOrder(g)
			}
			if g != w {
				t.Errorf("%s: %s differs\n--- restored\n%s\n--- server\n%s", label, name, g, w)
			}
		}
	}

	// Restore via dump format on the pinned MariaDB server (must preserve PRECEDES order too).
	dump, err := DatabaseDefinition(meta, DefinitionOptions{PrintHeader: true})
	if err != nil {
		t.Fatalf("DatabaseDefinition(dump): %v", err)
	}
	resetDatabase(t, conn, database)
	for _, stmt := range splitLikeMySQLClient(dump) {
		if _, err := conn.ExecContext(context.Background(), stmt); err != nil {
			t.Fatalf("restore dump %q: %v\n%s", stmt, err, dump)
		}
	}
	compare(t, "dump", storedObjects(t, conn, false), false)

	// Restore via SDL format on the pinned MariaDB server (SDL emits triggers sorted by name).
	sdl, err := DatabaseDefinition(meta, DefinitionOptions{SDLFormat: true})
	if err != nil {
		t.Fatalf("DatabaseDefinition(sdl): %v", err)
	}
	resetDatabase(t, conn, database)
	applySchema(t, conn, sdl)
	compare(t, "sdl", storedObjects(t, conn, false), true)
}

func stripTriggerActionOrder(s string) string {
	lines := strings.Split(s, "\n")
	if len(lines) >= 4 {
		lines[3] = "*"
	}
	return strings.Join(lines, "\n")
}
