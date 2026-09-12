package catalog

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
)

// TestContainer_LoadMetadataIndexTypesMatchEngine loads a snapshot introspected
// from real tables and checks every index against the type MySQL itself
// resolved. Sync records INDEX_TYPE verbatim; the catalog then stores the
// engine's default blank so SHOW CREATE TABLE renders no redundant USING, and
// defaultIndexType fills it back in. Only a live engine can say whether it
// fills in the right one — HASH on MEMORY, BTREE elsewhere — and whether a key
// declared USING BTREE on a MEMORY table survives that round trip.
func TestContainer_LoadMetadataIndexTypesMatchEngine(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	ddl := []string{
		`CREATE TABLE idx_innodb (
			id INT NOT NULL AUTO_INCREMENT,
			email VARCHAR(255) NOT NULL,
			body TEXT,
			PRIMARY KEY (id),
			UNIQUE KEY uq_email (email),
			KEY k_email_prefix (email(32)),
			FULLTEXT KEY ft_body (body)
		) ENGINE=InnoDB`,
		`CREATE TABLE idx_memory (
			id INT NOT NULL,
			label VARCHAR(64) NOT NULL,
			PRIMARY KEY (id),
			KEY k_label (label),
			KEY k_label_btree (label) USING BTREE
		) ENGINE=MEMORY`,
	}
	for _, stmt := range ddl {
		if err := ctr.execSQL(stmt); err != nil {
			t.Fatalf("container exec: %v", err)
		}
	}

	meta := mysqlLiveSnapshot(t, ctr, "test")
	c := New()
	report, err := c.LoadMetadata(context.Background(), meta)
	if err != nil {
		t.Fatalf("LoadMetadata: %v", err)
	}
	if len(report.Degraded) != 0 || len(report.Missing) != 0 {
		t.Fatalf("a snapshot taken from live tables must install whole; degraded=%v missing=%v", report.Degraded, report.Missing)
	}

	db := c.GetDatabase("test")
	if db == nil {
		t.Fatal("database test did not install")
	}
	for _, want := range mysqlLiveIndexTypes(t, ctr, "test") {
		tbl := db.Tables[strings.ToLower(want.table)]
		if tbl == nil {
			t.Errorf("table %s did not install", want.table)
			continue
		}
		found := false
		for _, idx := range tbl.Indexes {
			if !strings.EqualFold(idx.Name, want.index) {
				continue
			}
			found = true
			// An index whose type is the engine's default is stored blank, so
			// that SHOW CREATE TABLE does not render a redundant USING. Resolve
			// it the way a consumer must before comparing with the engine.
			got := idx.IndexType
			if got == "" {
				got = defaultIndexType(tbl.Engine)
			}
			if !strings.EqualFold(got, want.indexType) {
				t.Errorf("%s.%s index type = %q, want the engine's %q", want.table, want.index, got, want.indexType)
			}
		}
		if !found {
			t.Errorf("%s.%s did not install", want.table, want.index)
		}
	}
}

type liveIndexType struct {
	table     string
	index     string
	indexType string
}

// mysqlLiveIndexTypes reports the index type MySQL resolved for every index,
// which is what the loaded catalog has to agree with.
func mysqlLiveIndexTypes(t *testing.T, ctr *mysqlContainer, database string) []liveIndexType {
	t.Helper()
	rows, err := ctr.db.QueryContext(ctr.ctx, `
		SELECT TABLE_NAME, INDEX_NAME, INDEX_TYPE
		FROM information_schema.STATISTICS
		WHERE TABLE_SCHEMA = ? AND SEQ_IN_INDEX = 1
		ORDER BY TABLE_NAME, INDEX_NAME`, database)
	if err != nil {
		t.Fatalf("introspect index types: %v", err)
	}
	defer rows.Close()
	var out []liveIndexType
	for rows.Next() {
		var r liveIndexType
		if err := rows.Scan(&r.table, &r.index, &r.indexType); err != nil {
			t.Fatalf("scan index type: %v", err)
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("introspect index types: %v", err)
	}
	return out
}

// mysqlLiveSnapshot introspects one database the way sync does: COLUMN_TYPE for
// the column's full declared type, and INDEX_TYPE verbatim for every index.
func mysqlLiveSnapshot(t *testing.T, ctr *mysqlContainer, database string) *metadata.DatabaseSchemaMetadata {
	t.Helper()
	sm := &metadata.SchemaMetadata{}
	tables := map[string]*metadata.TableMetadata{}

	rows, err := ctr.db.QueryContext(ctr.ctx, `
		SELECT t.TABLE_NAME, t.ENGINE, t.TABLE_COLLATION,
		       c.COLUMN_NAME, c.COLUMN_TYPE, c.IS_NULLABLE, c.COLUMN_DEFAULT,
		       c.EXTRA, c.ORDINAL_POSITION, c.CHARACTER_SET_NAME, c.COLLATION_NAME
		FROM information_schema.TABLES t
		JOIN information_schema.COLUMNS c
			ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME
		WHERE t.TABLE_SCHEMA = ? AND t.TABLE_TYPE = 'BASE TABLE'
		ORDER BY t.TABLE_NAME, c.ORDINAL_POSITION`, database)
	if err != nil {
		t.Fatalf("introspect columns: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var table, columnName, columnType, nullable, extra string
		var engine, tableCollation, def, charset, collation sql.NullString
		var position int32
		if err := rows.Scan(&table, &engine, &tableCollation, &columnName, &columnType,
			&nullable, &def, &extra, &position, &charset, &collation); err != nil {
			t.Fatalf("scan column: %v", err)
		}
		tbl, ok := tables[table]
		if !ok {
			tbl = &metadata.TableMetadata{Name: table, Engine: engine.String, Collation: tableCollation.String}
			tables[table] = tbl
			sm.Tables = append(sm.Tables, tbl)
		}
		col := &metadata.ColumnMetadata{
			Name: columnName, Type: columnType, Nullable: nullable == "YES", Position: position,
			CharacterSet: charset.String, Collation: collation.String,
		}
		switch {
		case strings.Contains(strings.ToLower(extra), "auto_increment"):
			col.Default = "AUTO_INCREMENT"
		case def.Valid:
			col.Default = def.String
		}
		tbl.Columns = append(tbl.Columns, col)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("introspect columns: %v", err)
	}

	for _, tbl := range sm.Tables {
		tbl.Indexes = mysqlLiveIndexes(t, ctr, database, tbl.GetName())
	}
	return &metadata.DatabaseSchemaMetadata{Name: database, Schemas: []*metadata.SchemaMetadata{sm}}
}

func mysqlLiveIndexes(t *testing.T, ctr *mysqlContainer, database, table string) []*metadata.IndexMetadata {
	t.Helper()
	rows, err := ctr.db.QueryContext(ctr.ctx, `
		SELECT INDEX_NAME, NON_UNIQUE, INDEX_TYPE, COLUMN_NAME, SUB_PART, SEQ_IN_INDEX
		FROM information_schema.STATISTICS
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
		ORDER BY INDEX_NAME, SEQ_IN_INDEX`, database, table)
	if err != nil {
		t.Fatalf("introspect indexes: %v", err)
	}
	defer rows.Close()
	byName := map[string]*metadata.IndexMetadata{}
	var out []*metadata.IndexMetadata
	for rows.Next() {
		var name, indexType, column string
		var nonUnique, seq int
		var subPart sql.NullInt64
		if err := rows.Scan(&name, &nonUnique, &indexType, &column, &subPart, &seq); err != nil {
			t.Fatalf("scan index: %v", err)
		}
		idx, ok := byName[name]
		if !ok {
			idx = &metadata.IndexMetadata{
				Name:    name,
				Unique:  nonUnique == 0,
				Primary: name == "PRIMARY",
				Visible: true,
			}
			idx.Type = indexType
			byName[name] = idx
			out = append(out, idx)
		}
		expr := column
		if subPart.Valid {
			idx.KeyLength = append(idx.KeyLength, subPart.Int64)
		} else {
			idx.KeyLength = append(idx.KeyLength, -1)
		}
		idx.Expressions = append(idx.Expressions, expr)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("introspect indexes: %v", err)
	}
	return out
}
