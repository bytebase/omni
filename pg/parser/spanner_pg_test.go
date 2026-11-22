package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

func TestDialectSpannerPG_Isolation(t *testing.T) {
	// Spanner-PG-only DDL extensions must be accepted in DialectSpannerPG
	// and rejected in DialectDefault (strict PostgreSQL 17).
	spannerOnlyDDL := []string{
		`CREATE TABLE child (id bigint, cid bigint, PRIMARY KEY (id, cid)) INTERLEAVE IN PARENT parent`,
		`CREATE TABLE child (id bigint, cid bigint, PRIMARY KEY (id, cid)) INTERLEAVE IN PARENT parent ON DELETE CASCADE`,
		`CREATE TABLE child (id bigint, cid bigint, PRIMARY KEY (id, cid)) INTERLEAVE IN parent`,
		`CREATE TABLE t (id bigint PRIMARY KEY, ts timestamptz) TTL INTERVAL '30 days' ON ts`,
		`CREATE TABLE t (id bigint PRIMARY KEY, c text HIDDEN)`,
		`CREATE TABLE t (id bigint PRIMARY KEY, c text, g text GENERATED ALWAYS AS (upper(c)) VIRTUAL)`,
		`ALTER TABLE t ADD TTL INTERVAL '30 days' ON ts`,
		`ALTER TABLE t ALTER TTL INTERVAL '60 days' ON ts`,
		`ALTER TABLE t SET INTERLEAVE IN PARENT parent ON DELETE CASCADE`,
		`ALTER TABLE t SET ON DELETE CASCADE`,
		`CREATE CHANGE STREAM cs`,
		`CREATE CHANGE STREAM cs FOR ALL`,
		`CREATE CHANGE STREAM cs FOR t1, t2(c1, c2) WITH (retention_period = '7d')`,
		`ALTER CHANGE STREAM cs SET FOR ALL`,
		`ALTER CHANGE STREAM cs SET (retention_period = '7d')`,
		`ALTER CHANGE STREAM cs RESET (retention_period)`,
		`ALTER CHANGE STREAM cs DROP FOR ALL`,
		`DROP CHANGE STREAM IF EXISTS cs`,
		`CREATE VIEW v SQL SECURITY INVOKER AS SELECT 1`,
		`CREATE OR REPLACE VIEW v SQL SECURITY DEFINER AS SELECT 1`,
		`CREATE SEQUENCE s BIT_REVERSED_POSITIVE`,
		`CREATE SEQUENCE s BIT_REVERSED_POSITIVE SKIP RANGE 1000 5000 START COUNTER WITH 100`,
		`ALTER SEQUENCE s SKIP RANGE 1000 5000`,
		`ALTER SEQUENCE s NO SKIP RANGE`,
		`ALTER SEQUENCE s RESTART COUNTER WITH 10`,
		`CREATE INDEX idx ON t (c1) INCLUDE (c2) INTERLEAVE IN parent`,
		`CREATE SEARCH INDEX idx ON t (tok_col)`,
		`CREATE SEARCH INDEX idx ON t (tok_col) PARTITION BY p1 ORDER BY o1 DESC INTERLEAVE IN parent WHERE c1 IS NOT NULL WITH (sort_order_sharding = true)`,
		`ALTER SEARCH INDEX idx ADD COLUMN c2`,
		`ALTER SEARCH INDEX idx DROP COLUMN c2`,
		`DROP SEARCH INDEX IF EXISTS idx`,
		`CREATE LOCALITY GROUP IF NOT EXISTS lg`,
		`DROP LOCALITY GROUP IF EXISTS lg`,
	}

	for _, sql := range spannerOnlyDDL {
		if _, err := ParseWithDialect(sql, DialectSpannerPG); err != nil {
			t.Errorf("ParseWithDialect(%q, DialectSpannerPG) unexpected error: %v", sql, err)
		}
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) in DialectDefault unexpectedly succeeded; want error", sql)
		}
	}

	// Standard PG 17 features rejected in Spanner PG must remain accepted in DialectDefault.
	pg17Only := []string{
		``,
		`;`,
		`SELECT a IS DOCUMENT`,
		`SELECT XMLCONCAT(a)`,
		`SELECT XMLELEMENT(NAME a, XMLATTRIBUTES(a), a)`,
		`SELECT XMLFOREST(a AS b)`,
		`SELECT XMLPARSE(DOCUMENT a)`,
		`SELECT XMLPI(NAME a)`,
		`SELECT XMLROOT(a, VERSION a, STANDALONE YES)`,
		`SELECT JSON_QUERY(a, a PASSING a AS a)`,
		`SELECT JSON_QUERY(a, a WITH WRAPPER ERROR ON EMPTY)`,
		`SELECT JSON_EXISTS(a, a ERROR ON ERROR)`,
		`SELECT a AT LOCAL`,
		`ALTER DOMAIN a ADD CONSTRAINT c NOT NULL`,
		`COPY a FROM 's' FORCE NOT NULL *`,
		`COPY a FROM 's' (a DEFAULT)`,
		`MERGE INTO a USING b ON a.id = b.id WHEN NOT MATCHED BY SOURCE THEN DELETE`,
		`ALTER TABLE a ALTER c SET EXPRESSION AS (1)`,
		`ALTER TABLE a SET ACCESS METHOD DEFAULT`,
		`ALTER STATISTICS a SET STATISTICS DEFAULT`,
		`ALTER TYPE a SET (a)`,
		`ALTER INDEX a DROP a`,
		`CREATE RULE r AS ON UPDATE TO a DO NOTHING`,
		`CLUSTER (VERBOSE)`,
	}

	for _, sql := range pg17Only {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q) in DialectDefault unexpected error: %v", sql, err)
		}
		if _, err := ParseWithDialect(sql, DialectSpannerPG); err == nil {
			t.Errorf("ParseWithDialect(%q, DialectSpannerPG) unexpectedly succeeded; want error", sql)
		}
	}
}

func TestDialectSpannerPG_AST(t *testing.T) {
	t.Run("CreateTableInterleaveAndTTL", func(t *testing.T) {
		sql := `CREATE TABLE child (id bigint, cid bigint, ts timestamptz, c text HIDDEN, g text GENERATED ALWAYS AS (upper(c)) VIRTUAL, PRIMARY KEY (id, cid)) INTERLEAVE IN PARENT parent ON DELETE CASCADE TTL INTERVAL '30 days' ON ts`
		list, err := ParseWithDialect(sql, DialectSpannerPG)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		stmt := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.CreateStmt)
		if stmt.Interleave == nil || !stmt.Interleave.InParent || stmt.Interleave.Parent.Relname != "parent" || stmt.Interleave.OnDelete != 'c' {
			t.Fatalf("unexpected Interleave: %+v", stmt.Interleave)
		}
		if stmt.TTL == nil || stmt.TTL.Column != "ts" || stmt.TTL.Interval == nil {
			t.Fatalf("unexpected TTL: %+v", stmt.TTL)
		}
		cCol := stmt.TableElts.Items[3].(*nodes.ColumnDef)
		if !cCol.IsHidden || cCol.IsVirtual {
			t.Errorf("column c: IsHidden=%v IsVirtual=%v, want true/false", cCol.IsHidden, cCol.IsVirtual)
		}
		gCol := stmt.TableElts.Items[4].(*nodes.ColumnDef)
		if !gCol.IsVirtual || gCol.IsHidden {
			t.Errorf("column g: IsHidden=%v IsVirtual=%v, want false/true", gCol.IsHidden, gCol.IsVirtual)
		}
	})

	t.Run("CreateViewSQLSecurity", func(t *testing.T) {
		list, err := ParseWithDialect(`CREATE OR REPLACE VIEW v SQL SECURITY INVOKER AS SELECT 1`, DialectSpannerPG)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		stmt := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.ViewStmt)
		if !stmt.Replace || stmt.SqlSecurity != "invoker" {
			t.Fatalf("unexpected ViewStmt: Replace=%v SqlSecurity=%q", stmt.Replace, stmt.SqlSecurity)
		}
	})

	t.Run("CreateChangeStream", func(t *testing.T) {
		list, err := ParseWithDialect(`CREATE CHANGE STREAM IF NOT EXISTS cs FOR t1(), t2(c1, c2) WITH (retention_period = '7d')`, DialectSpannerPG)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		stmt := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.CreateChangeStreamStmt)
		if !stmt.IfNotExists || stmt.Name.Relname != "cs" || stmt.ForAll {
			t.Fatalf("unexpected CreateChangeStreamStmt: %+v", stmt)
		}
		if stmt.WatchedObjects == nil || len(stmt.WatchedObjects.Items) != 2 {
			t.Fatalf("expected 2 watched items, got %v", stmt.WatchedObjects)
		}
		item0 := stmt.WatchedObjects.Items[0].(*nodes.ChangeStreamForItem)
		if item0.AllColumns || len(item0.Columns.Items) != 0 {
			t.Errorf("item0: AllColumns=%v len(Columns)=%d, want false/0", item0.AllColumns, len(item0.Columns.Items))
		}
		item1 := stmt.WatchedObjects.Items[1].(*nodes.ChangeStreamForItem)
		if item1.AllColumns || len(item1.Columns.Items) != 2 {
			t.Errorf("item1: AllColumns=%v len(Columns)=%d, want false/2", item1.AllColumns, len(item1.Columns.Items))
		}
		if stmt.Options == nil || len(stmt.Options.Items) != 1 {
			t.Errorf("expected 1 option, got %v", stmt.Options)
		}
	})

	t.Run("CreateSearchIndex", func(t *testing.T) {
		sql := `CREATE SEARCH INDEX idx ON t (tok_col) INCLUDE (c2) PARTITION BY p1, p2 ORDER BY o1 DESC INTERLEAVE IN parent WHERE c1 IS NOT NULL WITH (sort_order_sharding = true)`
		list, err := ParseWithDialect(sql, DialectSpannerPG)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		stmt := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.IndexStmt)
		if !stmt.IsSearch || stmt.Idxname != "idx" {
			t.Fatalf("unexpected IndexStmt: IsSearch=%v Idxname=%q", stmt.IsSearch, stmt.Idxname)
		}
		if stmt.PartitionBy == nil || len(stmt.PartitionBy.Items) != 2 {
			t.Errorf("expected 2 PARTITION BY cols, got %v", stmt.PartitionBy)
		}
		if stmt.OrderBy == nil || len(stmt.OrderBy.Items) != 1 {
			t.Errorf("expected 1 ORDER BY elem, got %v", stmt.OrderBy)
		}
		if stmt.InterleaveParent == nil || stmt.InterleaveParent.Relname != "parent" {
			t.Errorf("unexpected InterleaveParent: %+v", stmt.InterleaveParent)
		}
		if stmt.WhereClause == nil || stmt.Options == nil {
			t.Errorf("expected WhereClause and Options to be set")
		}
	})

	t.Run("CreateSequenceSpannerOptions", func(t *testing.T) {
		sql := `CREATE SEQUENCE s BIT_REVERSED_POSITIVE SKIP RANGE 1000 5000 START COUNTER WITH 100`
		list, err := ParseWithDialect(sql, DialectSpannerPG)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		stmt := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.CreateSeqStmt)
		if stmt.Options == nil || len(stmt.Options.Items) != 3 {
			t.Fatalf("expected 3 sequence options, got %v", stmt.Options)
		}
		if stmt.Options.Items[0].(*nodes.DefElem).Defname != "bit_reversed_positive" ||
			stmt.Options.Items[1].(*nodes.DefElem).Defname != "skip_range" ||
			stmt.Options.Items[2].(*nodes.DefElem).Defname != "start_counter" {
			t.Errorf("unexpected sequence options: %+v", stmt.Options.Items)
		}
	})

	t.Run("WalkAllNewNodes", func(t *testing.T) {
		stmts := []string{
			`CREATE TABLE child (id bigint PRIMARY KEY, ts timestamptz) INTERLEAVE IN PARENT parent ON DELETE CASCADE TTL INTERVAL '30 days' ON ts`,
			`CREATE CHANGE STREAM cs FOR t1(c1) WITH (retention_period = '7d')`,
			`ALTER CHANGE STREAM cs SET FOR t1(c1)`,
			`CREATE SEARCH INDEX idx ON t (tok) PARTITION BY p1 ORDER BY o1 DESC INTERLEAVE IN parent WHERE p1 IS NOT NULL WITH (sort_order_sharding = true)`,
			`CREATE LOCALITY GROUP lg`,
		}
		for _, sql := range stmts {
			list, err := ParseWithDialect(sql, DialectSpannerPG)
			if err != nil {
				t.Fatalf("parse %q: %v", sql, err)
			}
			visited := 0
			nodes.Inspect(list, func(n nodes.Node) bool {
				visited++
				return true
			})
			if visited < 3 {
				t.Errorf("walk(%q) visited only %d nodes", sql, visited)
			}
		}
	})

	t.Run("QuotedIdentifiersNotTreatedAsSpannerKeywords", func(t *testing.T) {
		// ALTER TABLE t DROP "ttl" and DROP "interleave" must remain AT_DropColumn,
		// not AT_SpannerDropTTL / AT_SpannerDropInterleave.
		for _, tc := range []struct {
			sql     string
			colName string
		}{
			{`ALTER TABLE t DROP "ttl"`, "ttl"},
			{`ALTER TABLE t DROP U&"ttl"`, "ttl"},
			{`ALTER TABLE t DROP "interleave"`, "interleave"},
			{`ALTER TABLE t DROP U&"interleave"`, "interleave"},
		} {
			list, err := ParseWithDialect(tc.sql, DialectSpannerPG)
			if err != nil {
				t.Fatalf("parse %q: %v", tc.sql, err)
			}
			stmt := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.AlterTableStmt)
			cmd := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
			if cmd.Subtype != int(nodes.AT_DropColumn) || cmd.Name != tc.colName {
				t.Errorf("%q: got Subtype=%d Name=%q, want AT_DropColumn (%d) Name=%q",
					tc.sql, cmd.Subtype, cmd.Name, nodes.AT_DropColumn, tc.colName)
			}
		}

		// Quoted "parent" is a table name, not the optional PARENT keyword.
		list, err := ParseWithDialect(`CREATE TABLE child (id bigint PRIMARY KEY) INTERLEAVE IN "parent"`, DialectSpannerPG)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		cs := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.CreateStmt)
		if cs.Interleave == nil || cs.Interleave.InParent || cs.Interleave.Parent.Relname != "parent" {
			t.Errorf("unexpected Interleave for quoted parent: %+v", cs.Interleave)
		}

		// Quoted Spanner extension keywords must not be accepted as DDL keywords.
		for _, bad := range []string{
			`CREATE "change" STREAM cs`,
			`CREATE CHANGE "stream" cs`,
			`CREATE "locality" GROUP lg`,
			`CREATE SEQUENCE s "bit_reversed_positive"`,
			`CREATE TABLE t (id bigint PRIMARY KEY, c text "hidden")`,
			`CREATE TABLE t (id bigint PRIMARY KEY, c text GENERATED ALWAYS AS (id + 1) "virtual")`,
		} {
			if _, err := ParseWithDialect(bad, DialectSpannerPG); err == nil {
				t.Errorf("ParseWithDialect(%q, DialectSpannerPG) unexpectedly succeeded; want error", bad)
			}
		}
	})

	t.Run("InterleaveParentNamedTTL", func(t *testing.T) {
		// Parent table explicitly named ttl after INTERLEAVE IN PARENT.
		for _, tc := range []struct {
			sql        string
			wantInPar  bool
			wantSchema string
			wantRel    string
			wantTTLCol string
		}{
			{
				sql:       `CREATE TABLE child (id bigint PRIMARY KEY) INTERLEAVE IN PARENT ttl`,
				wantInPar: true,
				wantRel:   "ttl",
			},
			{
				sql:        `CREATE TABLE child (id bigint PRIMARY KEY, ts timestamptz) INTERLEAVE IN PARENT ttl TTL INTERVAL '30 days' ON ts`,
				wantInPar:  true,
				wantRel:    "ttl",
				wantTTLCol: "ts",
			},
			{
				sql:        `CREATE TABLE child (id bigint PRIMARY KEY, ts timestamptz) INTERLEAVE IN parent TTL INTERVAL '30 days' ON ts`,
				wantInPar:  false,
				wantRel:    "parent",
				wantTTLCol: "ts",
			},
			{
				sql:        `CREATE TABLE child (id bigint PRIMARY KEY) INTERLEAVE IN PARENT ttl.p`,
				wantInPar:  true,
				wantSchema: "ttl",
				wantRel:    "p",
			},
		} {
			list, err := ParseWithDialect(tc.sql, DialectSpannerPG)
			if err != nil {
				t.Fatalf("parse %q: %v", tc.sql, err)
			}
			stmt := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.CreateStmt)
			if stmt.Interleave == nil || stmt.Interleave.InParent != tc.wantInPar ||
				stmt.Interleave.Parent.Schemaname != tc.wantSchema || stmt.Interleave.Parent.Relname != tc.wantRel {
				t.Errorf("%q: got Interleave=%+v (Parent=%+v), want InParent=%v Schema=%q Rel=%q",
					tc.sql, stmt.Interleave, stmt.Interleave.Parent, tc.wantInPar, tc.wantSchema, tc.wantRel)
			}
			if tc.wantTTLCol == "" {
				if stmt.TTL != nil {
					t.Errorf("%q: expected nil TTL, got %+v", tc.sql, stmt.TTL)
				}
			} else if stmt.TTL == nil || stmt.TTL.Column != tc.wantTTLCol {
				t.Errorf("%q: got TTL=%+v, want Column=%q", tc.sql, stmt.TTL, tc.wantTTLCol)
			}
		}

		// ALTER TABLE SET INTERLEAVE IN PARENT ttl (with and without ON DELETE).
		for _, sql := range []string{
			`ALTER TABLE child SET INTERLEAVE IN PARENT ttl`,
			`ALTER TABLE child SET INTERLEAVE IN PARENT ttl ON DELETE CASCADE`,
		} {
			list, err := ParseWithDialect(sql, DialectSpannerPG)
			if err != nil {
				t.Fatalf("parse %q: %v", sql, err)
			}
			stmt := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.AlterTableStmt)
			cmd := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
			if cmd.Subtype != int(nodes.AT_SpannerSetInterleave) {
				t.Fatalf("%q: got Subtype=%d, want AT_SpannerSetInterleave", sql, cmd.Subtype)
			}
			il := cmd.Def.(*nodes.SpannerInterleaveClause)
			if !il.InParent || il.Parent.Relname != "ttl" {
				t.Errorf("%q: got Interleave=%+v (Parent=%+v), want InParent=true Rel=ttl", sql, il, il.Parent)
			}
		}
	})

	t.Run("FuncTableCommonSubexprInSpannerPG", func(t *testing.T) {
		for _, sql := range []string{
			`SELECT * FROM COALESCE(1, 2)`,
			`SELECT * FROM CAST(1 AS bigint)`,
			`SELECT * FROM GREATEST(1, 2)`,
			`SELECT * FROM NULLIF(1, 2)`,
			`SELECT * FROM JSON_EXISTS(a)`,
		} {
			if _, err := ParseWithDialect(sql, DialectSpannerPG); err != nil {
				t.Errorf("ParseWithDialect(%q, DialectSpannerPG) unexpected error: %v", sql, err)
			}
		}
		if _, err := ParseWithDialect(`SELECT * FROM JSON_QUERY(a, a PASSING a AS a)`, DialectSpannerPG); err == nil {
			t.Errorf("expected PG17 JSON_QUERY PASSING in FROM to fail in DialectSpannerPG")
		}
	})
}
