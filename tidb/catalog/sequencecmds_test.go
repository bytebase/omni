package catalog

import (
	"context"
	"database/sql"
	"testing"

	"github.com/bytebase/omni/metadata"
)

// sequenceScript is run statement by statement on TiDB and through the
// catalog. Each statement must succeed or fail with the same code on both,
// and after each the catalog must write every sequence as TiDB does.
var sequenceScript = []string{
	"CREATE SEQUENCE s_default",
	"CREATE SEQUENCE s_opts START WITH 10 INCREMENT BY 2 MINVALUE 5 MAXVALUE 1000 CACHE 20 CYCLE COMMENT 'it''s'",
	"CREATE SEQUENCE s_desc INCREMENT -1",
	"CREATE SEQUENCE s_desc_start INCREMENT -2 START -10",
	"CREATE SEQUENCE s_desc_max INCREMENT -1 MAXVALUE -5",
	"CREATE SEQUENCE s_min MINVALUE 5",
	"CREATE SEQUENCE s_nocache NOCACHE",
	"CREATE SEQUENCE s_no START WITH 3 INCREMENT BY 4 NO MINVALUE NO MAXVALUE NO CYCLE",
	"CREATE SEQUENCE s_engine ENGINE=MyISAM",
	"CREATE SEQUENCE s_all_three START 1 MINVALUE -10 MAXVALUE 10 INCREMENT -1",
	"CREATE SEQUENCE s_conflict MINVALUE 10 MAXVALUE 5",
	"CREATE SEQUENCE s_start0 START 0",
	"CREATE SEQUENCE s_incr0 INCREMENT 0",
	"CREATE SEQUENCE s_cache0 CACHE 0",
	"CREATE SEQUENCE s_max MAXVALUE 9223372036854775807",
	"CREATE SEQUENCE s_desc_start_above INCREMENT -1 START 5",
	"CREATE SEQUENCE s_default",
	"CREATE SEQUENCE IF NOT EXISTS s_default",
	"CREATE TABLE s_default (a int)",
	"CREATE VIEW s_default AS SELECT 1",
	"CREATE TABLE t (a int)",
	"CREATE SEQUENCE t",
	"CREATE SEQUENCE IF NOT EXISTS t",
	"ALTER SEQUENCE s_desc MINVALUE -100",
	"ALTER SEQUENCE s_desc CACHE 5 INCREMENT 1",
	"ALTER SEQUENCE s_desc MAXVALUE 1000",
	"ALTER SEQUENCE s_desc START 5",
	"ALTER SEQUENCE s_desc NOCACHE",
	"ALTER SEQUENCE s_desc CACHE 0",
	"ALTER SEQUENCE s_min MINVALUE 10",
	"ALTER SEQUENCE s_opts RESTART",
	"ALTER SEQUENCE s_opts RESTART WITH 50",
	"ALTER SEQUENCE s_opts NOCYCLE",
	"ALTER SEQUENCE s_missing INCREMENT 2",
	"ALTER SEQUENCE IF EXISTS s_missing INCREMENT 2",
	"DROP SEQUENCE s_missing",
	"DROP SEQUENCE IF EXISTS s_missing",
	"DROP SEQUENCE s_no, s_missing",
	"DROP SEQUENCE s_engine",
	"DROP TABLE s_nocache",
}

func TestContainerSequencesMatchTiDB(t *testing.T) {
	ctr, done := startContainer(t)
	defer done()
	db := "omni_sequences"
	for _, stmt := range []string{"DROP DATABASE IF EXISTS " + db, "CREATE DATABASE " + db, "USE " + db} {
		if _, err := ctr.db.ExecContext(ctr.ctx, stmt); err != nil {
			t.Fatal(err)
		}
	}
	defer func() { _, _ = ctr.db.ExecContext(ctr.ctx, "USE test") }()
	c := New()
	if _, err := c.Exec("CREATE DATABASE "+db+"; USE "+db, nil); err != nil {
		t.Fatal(err)
	}
	for _, stmt := range sequenceScript {
		_, serverErr := ctr.db.ExecContext(ctr.ctx, stmt)
		results, err := c.Exec(stmt, nil)
		if err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
		if got, want := mysqlErrorCode(results[0].Error), mysqlErrorCode(serverErr); got != want {
			t.Errorf("%s: catalog %d (%v), TiDB %d (%v)", stmt, got, results[0].Error, want, serverErr)
		}
		for _, name := range sequencesOnTiDB(t, ctr, db) {
			if got, want := c.ShowCreateSequence(db, name), showCreateSequenceOnTiDB(t, ctr, name); got != want {
				t.Errorf("after %s: %s:\ncatalog %s\nTiDB    %s", stmt, name, got, want)
			}
		}
		if got, want := len(c.GetDatabase(db).Sequences), len(sequencesOnTiDB(t, ctr, db)); got != want {
			t.Errorf("after %s: catalog has %d sequences, TiDB %d", stmt, got, want)
		}
	}
	// Bytebase's sync reads a sequence from information_schema.SEQUENCES, and
	// the snapshot loads back to what TiDB writes.
	loaded := New()
	if _, err := loaded.LoadMetadataInto(context.Background(), &metadata.DatabaseSchemaMetadata{
		Name: db, Schemas: []*metadata.SchemaMetadata{{Sequences: syncSequences(t, ctr, db)}},
	}, db); err != nil {
		t.Fatal(err)
	}
	for _, name := range sequencesOnTiDB(t, ctr, db) {
		if got, want := loaded.ShowCreateSequence(db, name), showCreateSequenceOnTiDB(t, ctr, name); got != want {
			t.Errorf("loaded %s:\ncatalog %s\nTiDB    %s", name, got, want)
		}
	}
}

func sequencesOnTiDB(t *testing.T, ctr *mysqlContainer, db string) []string {
	t.Helper()
	rows, err := ctr.db.QueryContext(ctr.ctx, "SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME", db)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatal(err)
		}
		names = append(names, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return names
}

func showCreateSequenceOnTiDB(t *testing.T, ctr *mysqlContainer, name string) string {
	t.Helper()
	var n, text string
	if err := ctr.db.QueryRowContext(ctr.ctx, "SHOW CREATE SEQUENCE "+quoteIdent(name)).Scan(&n, &text); err != nil {
		t.Fatal(err)
	}
	return text
}

// syncSequences reads a database's sequences as Bytebase's TiDB sync does
// (getSequenceList in backend/plugin/db/tidb/sync.go).
func syncSequences(t *testing.T, ctr *mysqlContainer, db string) []*metadata.SequenceMetadata {
	t.Helper()
	rows, err := ctr.db.QueryContext(ctr.ctx, `SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME`, db)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var out []*metadata.SequenceMetadata
	for rows.Next() {
		s := &metadata.SequenceMetadata{DataType: "BIGINT"}
		var cycle sql.NullInt64
		if err := rows.Scan(&s.Name, &s.Start, &s.MinValue, &s.MaxValue, &s.Increment, &cycle, &s.CacheSize, &s.Comment); err != nil {
			t.Fatal(err)
		}
		s.Cycle = cycle.Int64 != 0
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return out
}
