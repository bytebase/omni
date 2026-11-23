package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/mssql/internal/testcontainer"
	mssqlparser "github.com/bytebase/omni/mssql/parser"
)

// TestSimulation holds the catalog to testdata/simulate.tsv: every core DDL
// statement SQL Server accepts in the parser's compat.tsv, run on each
// fixture, with what SQL Server 2022 did and what Exec does. An outcome is the
// primary SQL Server error number the statement is refused with, or ok: and a
// hash of the DDL that rebuilds what the database holds afterwards
// (simulationForms).
//
// SQL Server's side is on the test's tape, so the test starts no server.
// OMNI_RECORD=1 runs every row on the pinned SQL Server again; -update-compat
// commits the rows a catalog fix moved; and -simulate-show=<text> prints both
// schemas of the rows that hold the text.
func TestSimulation(t *testing.T) {
	rows, err := simulate.FromCompat(filepath.Join("..", "parser", "testdata", "compat.tsv"), simulationFixtures, simulated)
	if err != nil {
		t.Fatal(err)
	}
	simulate.Check(t, simulate.Suite{
		Path:     filepath.Join("testdata", "simulate.tsv"),
		Pins:     testcontainer.Pins,
		Fixtures: filepath.Join("testdata", "fixtures"),
		Rows:     rows,
		MadeUp:   testcontainer.Codec{}.MadeUp(),
		Omni:     simulationOmni,
		Engine:   simulationEngine,
	})
}

// simulationFixtures are the starting schemas, in testdata/fixtures. A
// compat.tsv statement writes a for every identifier it holds, so what a
// statement finds is set by what the fixture names a.
var simulationFixtures = []string{
	"one_table@missing", "one_table@table", "one_table@column", "one_table@index",
	"objects@view", "objects@fk", "objects@schema", "objects@sequence",
	"objects@synonym", "objects@procedure", "objects@function", "objects@trigger",
}

// simulatedStatement selects the core catalog DDL families whose effects stay
// inside the current database and match the objects mssql/catalog models:
// schemas, tables, indexes, views, sequences, synonyms, triggers, procedures,
// and functions. Server-, database-, security-, partition-, fulltext-, and
// external-object statements are left out.
var simulatedStatement = regexp.MustCompile(`(?i)^(?:` +
	`CREATE TABLE |ALTER TABLE |DROP TABLE |TRUNCATE TABLE ` +
	`|CREATE (?:(?:UNIQUE )?(?:CLUSTERED |NONCLUSTERED )?|(?:CLUSTERED |NONCLUSTERED )?COLUMNSTORE |PRIMARY XML |XML |SELECTIVE XML |SPATIAL )?INDEX |ALTER INDEX |DROP INDEX ` +
	`|CREATE (?:OR ALTER )?VIEW |ALTER VIEW |DROP VIEW ` +
	`|CREATE SCHEMA |ALTER SCHEMA |DROP SCHEMA ` +
	`|CREATE SEQUENCE |ALTER SEQUENCE |DROP SEQUENCE ` +
	`|CREATE SYNONYM |DROP SYNONYM ` +
	`|CREATE (?:OR ALTER )?(?:PROC|PROCEDURE) |ALTER (?:PROC|PROCEDURE) |DROP (?:PROC|PROCEDURE) ` +
	`|CREATE (?:OR ALTER )?FUNCTION |ALTER FUNCTION |DROP FUNCTION ` +
	`|CREATE (?:OR ALTER )?TRIGGER |ALTER TRIGGER |DROP TRIGGER ` +
	`)`)

func simulated(statement string) bool { return simulatedStatement.MatchString(statement) }

// simulationResetProcSQL installs [dbo].[_sim_reset] once per recording so
// each row resets the database with a 23-byte EXEC call instead of repeating
// the full reset script in every tape entry.
const simulationResetProcSQL = `CREATE OR ALTER PROCEDURE [dbo].[_sim_reset] AS
BEGIN
	SET NOCOUNT ON;
	DECLARE @sql NVARCHAR(MAX) = N'';
	SELECT @sql += N'ALTER TABLE ' + QUOTENAME(SCHEMA_NAME(schema_id)) + N'.' + QUOTENAME(name) + N' SET (SYSTEM_VERSIONING = OFF);' + CHAR(10)
	FROM sys.tables WHERE temporal_type = 2;
	IF @sql <> N'' EXEC sp_executesql @sql;
	SET @sql = N'';
	SELECT @sql += N'DROP TRIGGER ' + QUOTENAME(name) + N' ON DATABASE;' + CHAR(10)
	FROM sys.triggers WHERE parent_class = 0;
	SELECT @sql += N'DROP VIEW ' + QUOTENAME(SCHEMA_NAME(schema_id)) + N'.' + QUOTENAME(name) + N';' + CHAR(10)
	FROM sys.views WHERE is_ms_shipped = 0;
	SELECT @sql += N'ALTER TABLE ' + QUOTENAME(SCHEMA_NAME(parent.schema_id)) + N'.' + QUOTENAME(parent.name) + N' DROP CONSTRAINT ' + QUOTENAME(fk.name) + N';' + CHAR(10)
	FROM sys.foreign_keys fk JOIN sys.objects parent ON parent.object_id = fk.parent_object_id;
	IF @sql <> N'' EXEC sp_executesql @sql;
	SET @sql = N'';
	SELECT @sql += N'DROP TABLE ' + QUOTENAME(SCHEMA_NAME(schema_id)) + N'.' + QUOTENAME(name) + N';' + CHAR(10)
	FROM sys.tables WHERE is_ms_shipped = 0;
	IF @sql <> N'' EXEC sp_executesql @sql;
	SET @sql = N'';
	SELECT @sql += N'DROP PROCEDURE ' + QUOTENAME(SCHEMA_NAME(schema_id)) + N'.' + QUOTENAME(name) + N';' + CHAR(10)
	FROM sys.procedures WHERE is_ms_shipped = 0 AND NOT (schema_id = SCHEMA_ID(N'dbo') AND name IN (N'_sim_reset', N'_sim_forms'));
	SELECT @sql += N'DROP FUNCTION ' + QUOTENAME(SCHEMA_NAME(schema_id)) + N'.' + QUOTENAME(name) + N';' + CHAR(10)
	FROM sys.objects WHERE type IN ('FN', 'IF', 'TF') AND is_ms_shipped = 0;
	SELECT @sql += N'DROP SEQUENCE ' + QUOTENAME(SCHEMA_NAME(schema_id)) + N'.' + QUOTENAME(name) + N';' + CHAR(10)
	FROM sys.sequences;
	SELECT @sql += N'DROP SYNONYM ' + QUOTENAME(SCHEMA_NAME(schema_id)) + N'.' + QUOTENAME(name) + N';' + CHAR(10)
	FROM sys.synonyms;
	IF @sql <> N'' EXEC sp_executesql @sql;
	SET @sql = N'';
	SELECT @sql += N'DROP SCHEMA ' + QUOTENAME(name) + N';' + CHAR(10)
	FROM sys.schemas WHERE schema_id > 4 AND schema_id < 16384;
	IF @sql <> N'' EXEC sp_executesql @sql;
END;`

// simulationEngine asks SQL Server 2022, or the tape. Each row resets the
// database, builds the fixture batches, executes the statement, and after an
// accepted statement reads the stored objects back as DDL.
func simulationEngine(t testing.TB) func(simulate.Fixture, string) (simulate.Answer, error) {
	db := testcontainer.Open(t)
	ctx := context.Background()
	if _, err := db.ExecContext(ctx, simulationResetProcSQL); err != nil {
		t.Fatalf("install _sim_reset: %v", err)
	}
	if _, err := db.ExecContext(ctx, "CREATE OR ALTER PROCEDURE [dbo].[_sim_forms] AS\nBEGIN\n\tSET NOCOUNT ON;\n"+simulationFormsQuery+";\nEND;"); err != nil {
		t.Fatalf("install _sim_forms: %v", err)
	}
	return func(fx simulate.Fixture, statement string) (simulate.Answer, error) {
		if _, err := db.ExecContext(ctx, "EXEC [dbo].[_sim_reset]"); err != nil {
			return simulate.Answer{}, fmt.Errorf("reset database for fixture %s: %w", fx.Name, err)
		}
		for _, batch := range splitBatches(fx.SQL) {
			if _, err := db.ExecContext(ctx, batch); err != nil {
				return simulate.Answer{}, fmt.Errorf("build fixture %s: %w\n%s", fx.Name, err, batch)
			}
		}
		if _, err := db.ExecContext(ctx, statement+"\n"); err != nil {
			errs, ok := testcontainer.ServerErrors(err)
			if !ok || len(errs) == 0 {
				return simulate.Answer{}, err
			}
			return simulate.Answer{Code: strconv.Itoa(int(errs[0].Number))}, nil
		}
		return simulationForms(ctx, db)
	}
}

func splitBatches(script string) []string {
	var batches []string
	var cur strings.Builder
	flush := func() {
		if s := strings.TrimSpace(cur.String()); s != "" {
			batches = append(batches, s)
		}
		cur.Reset()
	}
	for _, line := range strings.Split(script, "\n") {
		if strings.EqualFold(strings.TrimSpace(line), "GO") {
			flush()
			continue
		}
		cur.WriteString(line)
		cur.WriteByte('\n')
	}
	flush()
	return batches
}

// simulationFormsQuery reads the database's user objects as T-SQL DDL in
// deterministic name order. What Diff compares is written, so the catalog is
// held to all of it. System-generated constraint names (is_system_named = 1)
// are omitted so unnamed constraints stay deterministic across recordings.
const simulationFormsQuery = `WITH key_cols AS (
	SELECT i.object_id, i.index_id,
		STRING_AGG(CAST(QUOTENAME(c2.name) + CASE WHEN ic2.is_descending_key = 1 THEN N' DESC' ELSE N' ASC' END AS NVARCHAR(MAX)), N', ')
			WITHIN GROUP (ORDER BY ic2.key_ordinal) AS cols
	FROM sys.indexes i
	JOIN sys.index_columns ic2 ON ic2.object_id = i.object_id AND ic2.index_id = i.index_id AND ic2.is_included_column = 0
	JOIN sys.columns c2 ON c2.object_id = ic2.object_id AND c2.column_id = ic2.column_id
	WHERE i.type > 0
	GROUP BY i.object_id, i.index_id
),
inc_cols AS (
	SELECT i.object_id, i.index_id,
		STRING_AGG(CAST(QUOTENAME(c2.name) AS NVARCHAR(MAX)), N', ')
			WITHIN GROUP (ORDER BY ic2.index_column_id) AS cols
	FROM sys.indexes i
	JOIN sys.index_columns ic2 ON ic2.object_id = i.object_id AND ic2.index_id = i.index_id AND ic2.is_included_column = 1
	JOIN sys.columns c2 ON c2.object_id = ic2.object_id AND c2.column_id = ic2.column_id
	WHERE i.type > 0
	GROUP BY i.object_id, i.index_id
),
all_idx_cols AS (
	SELECT i.object_id, i.index_id,
		STRING_AGG(CAST(QUOTENAME(c2.name) AS NVARCHAR(MAX)), N', ')
			WITHIN GROUP (ORDER BY ic2.index_column_id) AS cols
	FROM sys.indexes i
	JOIN sys.index_columns ic2 ON ic2.object_id = i.object_id AND ic2.index_id = i.index_id
	JOIN sys.columns c2 ON c2.object_id = ic2.object_id AND c2.column_id = ic2.column_id
	WHERE i.type = 6
	GROUP BY i.object_id, i.index_id
),
fk_parent_cols AS (
	SELECT fkc.constraint_object_id,
		STRING_AGG(CAST(QUOTENAME(pc.name) AS NVARCHAR(MAX)), N', ')
			WITHIN GROUP (ORDER BY fkc.constraint_column_id) AS cols
	FROM sys.foreign_key_columns fkc
	JOIN sys.columns pc ON pc.object_id = fkc.parent_object_id AND pc.column_id = fkc.parent_column_id
	GROUP BY fkc.constraint_object_id
),
fk_ref_cols AS (
	SELECT fkc.constraint_object_id,
		STRING_AGG(CAST(QUOTENAME(rc.name) AS NVARCHAR(MAX)), N', ')
			WITHIN GROUP (ORDER BY fkc.constraint_column_id) AS cols
	FROM sys.foreign_key_columns fkc
	JOIN sys.columns rc ON rc.object_id = fkc.referenced_object_id AND rc.column_id = fkc.referenced_column_id
	GROUP BY fkc.constraint_object_id
),
tbl_cols AS (
	SELECT c.object_id,
		STRING_AGG(
			CAST(
				QUOTENAME(c.name) + N' ' +
				CASE
					WHEN cc.object_id IS NOT NULL THEN
						N'AS ' + cc.definition
						+ CASE WHEN cc.is_persisted = 1 THEN N' PERSISTED' ELSE N'' END
						+ CASE WHEN cc.is_persisted = 1 AND c.is_nullable = 0 THEN N' NOT NULL' ELSE N'' END
					ELSE
						CASE
							WHEN ty.name IN (N'decimal', N'numeric') THEN
								ty.name + N'(' + CAST(c.precision AS NVARCHAR(10)) + N',' + CAST(c.scale AS NVARCHAR(10)) + N')'
							WHEN ty.name IN (N'datetime2', N'datetimeoffset', N'time') THEN
								ty.name + N'(' + CAST(c.scale AS NVARCHAR(10)) + N')'
							WHEN ty.name IN (N'nvarchar', N'nchar') THEN
								ty.name + N'(' + CASE WHEN c.max_length = -1 THEN N'max' ELSE CAST(c.max_length / 2 AS NVARCHAR(10)) END + N')'
							WHEN ty.name IN (N'varchar', N'char', N'varbinary', N'binary') THEN
								ty.name + N'(' + CASE WHEN c.max_length = -1 THEN N'max' ELSE CAST(c.max_length AS NVARCHAR(10)) END + N')'
							ELSE ty.name
						END
						+ CASE WHEN ic.object_id IS NOT NULL THEN
							N' IDENTITY(' + CAST(ic.seed_value AS NVARCHAR(40)) + N',' + CAST(ic.increment_value AS NVARCHAR(40)) + N')'
							ELSE N'' END
						+ CASE WHEN c.is_nullable = 0 THEN N' NOT NULL' ELSE N' NULL' END
						+ CASE WHEN dc.object_id IS NOT NULL THEN
							CASE WHEN dc.is_system_named = 0 THEN N' CONSTRAINT ' + QUOTENAME(dc.name) ELSE N'' END
							+ N' DEFAULT ' + dc.definition
							ELSE N'' END
				END
			AS NVARCHAR(MAX)),
			N', '
		) WITHIN GROUP (ORDER BY c.column_id) AS text
	FROM sys.columns c
	JOIN sys.tables t ON t.object_id = c.object_id AND t.is_ms_shipped = 0
	JOIN sys.types ty ON ty.user_type_id = c.user_type_id
	LEFT JOIN sys.computed_columns cc ON cc.object_id = c.object_id AND cc.column_id = c.column_id
	LEFT JOIN sys.identity_columns ic ON ic.object_id = c.object_id AND ic.column_id = c.column_id
	LEFT JOIN sys.default_constraints dc ON dc.parent_object_id = c.object_id AND dc.parent_column_id = c.column_id
	WHERE c.is_hidden = 0 AND c.graph_type IS NULL
	GROUP BY c.object_id
),
tbl_kcs AS (
	SELECT kc.parent_object_id AS object_id,
		STRING_AGG(
			CAST(
				N'ALTER TABLE ' + QUOTENAME(SCHEMA_NAME(t.schema_id)) + N'.' + QUOTENAME(t.name) + N' ADD '
				+ CASE WHEN kc.is_system_named = 0 THEN N'CONSTRAINT ' + QUOTENAME(kc.name) + N' ' ELSE N'' END
				+ CASE kc.type WHEN 'PK' THEN N'PRIMARY KEY ' ELSE N'UNIQUE ' END
				+ CASE i.type WHEN 1 THEN N'CLUSTERED' ELSE N'NONCLUSTERED' END
				+ N' (' + ISNULL(kcols.cols, N'') + N');' + CHAR(10)
			AS NVARCHAR(MAX)),
			N''
		) WITHIN GROUP (ORDER BY kc.is_system_named, CASE WHEN kc.is_system_named = 0 THEN kc.name ELSE N'' END, kc.type, i.index_id) AS text
	FROM sys.key_constraints kc
	JOIN sys.tables t ON t.object_id = kc.parent_object_id AND t.is_ms_shipped = 0
	JOIN sys.indexes i ON i.object_id = kc.parent_object_id AND i.index_id = kc.unique_index_id
	LEFT JOIN key_cols kcols ON kcols.object_id = i.object_id AND kcols.index_id = i.index_id
	GROUP BY kc.parent_object_id
),
tbl_cks AS (
	SELECT ck.parent_object_id AS object_id,
		STRING_AGG(
			CAST(
				N'ALTER TABLE ' + QUOTENAME(SCHEMA_NAME(t.schema_id)) + N'.' + QUOTENAME(t.name) + N' ADD '
				+ CASE WHEN ck.is_system_named = 0 THEN N'CONSTRAINT ' + QUOTENAME(ck.name) + N' ' ELSE N'' END
				+ N'CHECK ' + ck.definition + N';' + CHAR(10)
			AS NVARCHAR(MAX)),
			N''
		) WITHIN GROUP (ORDER BY ck.is_system_named, CASE WHEN ck.is_system_named = 0 THEN ck.name ELSE N'' END, ck.definition) AS text
	FROM sys.check_constraints ck
	JOIN sys.tables t ON t.object_id = ck.parent_object_id AND t.is_ms_shipped = 0
	GROUP BY ck.parent_object_id
),
tbl_fks AS (
	SELECT fk.parent_object_id AS object_id,
		STRING_AGG(
			CAST(
				N'ALTER TABLE ' + QUOTENAME(SCHEMA_NAME(t.schema_id)) + N'.' + QUOTENAME(t.name) + N' ADD '
				+ CASE WHEN fk.is_system_named = 0 THEN N'CONSTRAINT ' + QUOTENAME(fk.name) + N' ' ELSE N'' END
				+ N'FOREIGN KEY (' + ISNULL(fpc.cols, N'') + N') REFERENCES '
				+ QUOTENAME(OBJECT_SCHEMA_NAME(fk.referenced_object_id)) + N'.' + QUOTENAME(OBJECT_NAME(fk.referenced_object_id))
				+ N' (' + ISNULL(frc.cols, N'') + N')'
				+ CASE fk.delete_referential_action WHEN 1 THEN N' ON DELETE CASCADE' WHEN 2 THEN N' ON DELETE SET NULL' WHEN 3 THEN N' ON DELETE SET DEFAULT' ELSE N'' END
				+ CASE fk.update_referential_action WHEN 1 THEN N' ON UPDATE CASCADE' WHEN 2 THEN N' ON UPDATE SET NULL' WHEN 3 THEN N' ON UPDATE SET DEFAULT' ELSE N'' END
				+ N';' + CHAR(10)
			AS NVARCHAR(MAX)),
			N''
		) WITHIN GROUP (ORDER BY fk.is_system_named, CASE WHEN fk.is_system_named = 0 THEN fk.name ELSE N'' END, fk.object_id) AS text
	FROM sys.foreign_keys fk
	JOIN sys.tables t ON t.object_id = fk.parent_object_id AND t.is_ms_shipped = 0
	LEFT JOIN fk_parent_cols fpc ON fpc.constraint_object_id = fk.object_id
	LEFT JOIN fk_ref_cols frc ON frc.constraint_object_id = fk.object_id
	GROUP BY fk.parent_object_id
),
tbl_idxs AS (
	SELECT i.object_id,
		STRING_AGG(
			CAST(
				N'CREATE '
				+ CASE WHEN i.is_unique = 1 THEN N'UNIQUE ' ELSE N'' END
				+ CASE i.type
					WHEN 1 THEN N'CLUSTERED '
					WHEN 2 THEN N'NONCLUSTERED '
					WHEN 5 THEN N'CLUSTERED COLUMNSTORE '
					WHEN 6 THEN N'NONCLUSTERED COLUMNSTORE '
					ELSE N''
				  END
				+ N'INDEX ' + QUOTENAME(i.name) + N' ON ' + QUOTENAME(SCHEMA_NAME(t.schema_id)) + N'.' + QUOTENAME(t.name)
				+ CASE
					WHEN i.type = 5 THEN N''
					WHEN i.type = 6 THEN N' (' + ISNULL(acols.cols, N'') + N')'
					ELSE N' (' + ISNULL(kcols.cols, N'') + N')' + ISNULL(N' INCLUDE (' + icols.cols + N')', N'')
				  END
				+ CASE WHEN i.has_filter = 1 AND i.filter_definition IS NOT NULL THEN N' WHERE ' + i.filter_definition ELSE N'' END
				+ N';' + CHAR(10)
			AS NVARCHAR(MAX)),
			N''
		) WITHIN GROUP (ORDER BY i.name) AS text
	FROM sys.indexes i
	JOIN sys.tables t ON t.object_id = i.object_id AND t.is_ms_shipped = 0
	LEFT JOIN key_cols kcols ON kcols.object_id = i.object_id AND kcols.index_id = i.index_id
	LEFT JOIN inc_cols icols ON icols.object_id = i.object_id AND icols.index_id = i.index_id
	LEFT JOIN all_idx_cols acols ON acols.object_id = i.object_id AND acols.index_id = i.index_id
	WHERE i.type > 0
	  AND i.is_primary_key = 0
	  AND i.is_unique_constraint = 0
	  AND NOT EXISTS (
		  SELECT 1 FROM sys.index_columns ic3
		  JOIN sys.columns c3 ON c3.object_id = ic3.object_id AND c3.column_id = ic3.column_id
		  WHERE ic3.object_id = i.object_id AND ic3.index_id = i.index_id AND (c3.is_hidden = 1 OR c3.graph_type IS NOT NULL)
	  )
	GROUP BY i.object_id
),
obj_trgs AS (
	SELECT tr.parent_id AS object_id,
		STRING_AGG(CAST(N'GO' + CHAR(10) + LTRIM(RTRIM(OBJECT_DEFINITION(tr.object_id))) + CHAR(10) AS NVARCHAR(MAX)), N'')
			WITHIN GROUP (ORDER BY tr.name) AS text
	FROM sys.triggers tr
	WHERE tr.parent_class = 1
	GROUP BY tr.parent_id
)
SELECT kind, name, text FROM (
	SELECT CAST(N'SCHEMA' AS NVARCHAR(32)) AS kind,
		CAST(s.name AS NVARCHAR(260)) AS name,
		CAST(N'CREATE SCHEMA ' + QUOTENAME(s.name) + N';' + CHAR(10) AS NVARCHAR(MAX)) AS text
	FROM sys.schemas s
	WHERE s.schema_id > 4 AND s.schema_id < 16384
	UNION ALL
	SELECT CAST(N'TABLE' AS NVARCHAR(32)),
		CAST(s.name + N'.' + t.name AS NVARCHAR(260)),
		CAST(
			CASE
				WHEN t.is_edge = 1 AND tc.text IS NULL
					THEN N'CREATE TABLE ' + QUOTENAME(s.name) + N'.' + QUOTENAME(t.name) + N' AS EDGE;' + CHAR(10)
				ELSE
					N'CREATE TABLE ' + QUOTENAME(s.name) + N'.' + QUOTENAME(t.name) + N' (' + ISNULL(tc.text, N'') + N')'
					+ CASE WHEN t.is_edge = 1 THEN N' AS EDGE' ELSE N'' END + N';' + CHAR(10)
			END
			+ ISNULL(tk.text, N'')
			+ ISNULL(tck.text, N'')
			+ ISNULL(tfk.text, N'')
			+ ISNULL(ti.text, N'')
			+ ISNULL(ttr.text, N'')
		AS NVARCHAR(MAX))
	FROM sys.tables t
	JOIN sys.schemas s ON s.schema_id = t.schema_id
	LEFT JOIN tbl_cols tc ON tc.object_id = t.object_id
	LEFT JOIN tbl_kcs tk ON tk.object_id = t.object_id
	LEFT JOIN tbl_cks tck ON tck.object_id = t.object_id
	LEFT JOIN tbl_fks tfk ON tfk.object_id = t.object_id
	LEFT JOIN tbl_idxs ti ON ti.object_id = t.object_id
	LEFT JOIN obj_trgs ttr ON ttr.object_id = t.object_id
	WHERE t.is_ms_shipped = 0
	UNION ALL
	SELECT CAST(N'VIEW' AS NVARCHAR(32)),
		CAST(s.name + N'.' + v.name AS NVARCHAR(260)),
		CAST(LTRIM(RTRIM(OBJECT_DEFINITION(v.object_id))) + CHAR(10) + ISNULL(vtr.text, N'') AS NVARCHAR(MAX))
	FROM sys.views v
	JOIN sys.schemas s ON s.schema_id = v.schema_id
	LEFT JOIN obj_trgs vtr ON vtr.object_id = v.object_id
	WHERE v.is_ms_shipped = 0
	UNION ALL
	SELECT CAST(N'SEQUENCE' AS NVARCHAR(32)),
		CAST(s.name + N'.' + sq.name AS NVARCHAR(260)),
		CAST(
			N'CREATE SEQUENCE ' + QUOTENAME(s.name) + N'.' + QUOTENAME(sq.name)
			+ N' AS ' + ty.name
			+ N' START WITH ' + CAST(sq.start_value AS NVARCHAR(40))
			+ N' INCREMENT BY ' + CAST(sq.increment AS NVARCHAR(40))
			+ N' MINVALUE ' + CAST(sq.minimum_value AS NVARCHAR(40))
			+ N' MAXVALUE ' + CAST(sq.maximum_value AS NVARCHAR(40))
			+ CASE WHEN sq.is_cycling = 1 THEN N' CYCLE' ELSE N' NO CYCLE' END
			+ CASE WHEN sq.is_cached = 0 THEN N' NO CACHE' WHEN sq.cache_size IS NOT NULL THEN N' CACHE ' + CAST(sq.cache_size AS NVARCHAR(40)) ELSE N'' END
			+ N';' + CHAR(10)
		AS NVARCHAR(MAX))
	FROM sys.sequences sq
	JOIN sys.schemas s ON s.schema_id = sq.schema_id
	JOIN sys.types ty ON ty.user_type_id = sq.user_type_id
	UNION ALL
	SELECT CAST(N'SYNONYM' AS NVARCHAR(32)),
		CAST(s.name + N'.' + sy.name AS NVARCHAR(260)),
		CAST(N'CREATE SYNONYM ' + QUOTENAME(s.name) + N'.' + QUOTENAME(sy.name) + N' FOR ' + sy.base_object_name + N';' + CHAR(10) AS NVARCHAR(MAX))
	FROM sys.synonyms sy
	JOIN sys.schemas s ON s.schema_id = sy.schema_id
	UNION ALL
	SELECT CAST(N'PROCEDURE' AS NVARCHAR(32)),
		CAST(s.name + N'.' + p.name AS NVARCHAR(260)),
		CAST(LTRIM(RTRIM(OBJECT_DEFINITION(p.object_id))) + CHAR(10) AS NVARCHAR(MAX))
	FROM sys.procedures p
	JOIN sys.schemas s ON s.schema_id = p.schema_id
	WHERE p.is_ms_shipped = 0 AND NOT (s.name = N'dbo' AND p.name IN (N'_sim_reset', N'_sim_forms'))
	UNION ALL
	SELECT CAST(N'FUNCTION' AS NVARCHAR(32)),
		CAST(s.name + N'.' + o.name AS NVARCHAR(260)),
		CAST(LTRIM(RTRIM(OBJECT_DEFINITION(o.object_id))) + CHAR(10) AS NVARCHAR(MAX))
	FROM sys.objects o
	JOIN sys.schemas s ON s.schema_id = o.schema_id
	WHERE o.type IN ('FN', 'IF', 'TF') AND o.is_ms_shipped = 0
) o ORDER BY name COLLATE Latin1_General_BIN2, kind`

func simulationForms(ctx context.Context, db *sql.DB) (simulate.Answer, error) {
	rows, err := db.QueryContext(ctx, "EXEC [dbo].[_sim_forms]")
	if err != nil {
		return simulate.Answer{}, fmt.Errorf("read the stored schema: %w", err)
	}
	defer rows.Close()
	var forms []simulate.Form
	for rows.Next() {
		var f simulate.Form
		if err := rows.Scan(&f.Kind, &f.Name, &f.Text); err != nil {
			return simulate.Answer{}, err
		}
		forms = append(forms, f)
	}
	return simulate.Answer{Forms: forms}, rows.Err()
}

// simulationOmni runs the row in the catalog. The catalog stores what SQL
// Server stores when Diff finds nothing between it and a catalog loaded from
// SQL Server's forms in either direction.
func simulationOmni(fx simulate.Fixture, statement string, engine simulate.Answer) simulate.Omni {
	c, err := LoadSQL(fx.SQL)
	if err != nil {
		panic(fmt.Sprintf("the catalog does not load fixture %s: %v", fx.Name, err))
	}
	results, err := c.Exec(statement, nil)
	var syntax *mssqlparser.ParseError
	switch {
	case errors.As(err, &syntax) && syntax.Code != 0:
		return simulate.Omni{Code: strconv.Itoa(syntax.Code), Why: err.Error()}
	case errors.As(err, &syntax):
		return simulate.Omni{Code: "102", Why: err.Error()}
	case err != nil:
		return simulate.Omni{Code: simulate.NoCode, Why: err.Error()}
	}
	for _, r := range results {
		var refusal *Error
		switch {
		case r.Error == nil:
		case errors.As(r.Error, &refusal) && refusal.Code > 0:
			return simulate.Omni{Code: strconv.Itoa(refusal.Code), Why: r.Error.Error()}
		default:
			return simulate.Omni{Code: simulate.NoCode, Why: r.Error.Error()}
		}
	}
	if engine.Refused() {
		return simulate.Omni{}
	}

	want, err := loadSimulationForms(engine.Forms)
	if err != nil {
		return simulate.Omni{Why: fmt.Sprintf("the catalog does not load what SQL Server stores: %v", err)}
	}
	there := Diff(c, want)
	back := Diff(want, c)
	if there.IsEmpty() && back.IsEmpty() {
		return simulate.Omni{Same: true}
	}
	plan := GenerateMigration(c, want, there).SQL()
	if there.IsEmpty() {
		plan = "nothing, though from there it would run: " + GenerateMigration(want, c, back).SQL()
	}
	return simulate.Omni{Why: "to store what SQL Server stores, the catalog would run: " + strings.Join(strings.Fields(plan), " ")}
}

// loadSimulationForms builds a catalog from SQL Server's forms, ordering
// schemas and tables before objects that reference them (constraints, views,
// routines, triggers).
func loadSimulationForms(forms []simulate.Form) (*Catalog, error) {
	var script strings.Builder
	write := func(text string) {
		text = strings.TrimSpace(text)
		if text == "" {
			return
		}
		if script.Len() > 0 {
			script.WriteString("\nGO\n")
		}
		script.WriteString(text)
		script.WriteByte('\n')
	}
	var deferredFKs []string
	for _, kind := range []string{"SCHEMA", "SEQUENCE", "SYNONYM", "TABLE", "VIEW", "FUNCTION", "PROCEDURE"} {
		for _, f := range forms {
			if f.Kind != kind {
				continue
			}
			if kind == "TABLE" {
				var keep []string
				for _, line := range strings.Split(f.Text, "\n") {
					if strings.HasPrefix(line, "ALTER TABLE ") && strings.Contains(line, " FOREIGN KEY ") {
						deferredFKs = append(deferredFKs, line)
						continue
					}
					keep = append(keep, line)
				}
				write(strings.Join(keep, "\n"))
				continue
			}
			write(f.Text)
		}
		if kind == "TABLE" && len(deferredFKs) > 0 {
			write(strings.Join(deferredFKs, "\n"))
		}
	}
	c, err := LoadSQL(script.String())
	if err != nil {
		return nil, err
	}
	// Every SQL Server database always contains the built-in dbo schema even
	// when all user objects in dbo have been dropped or transferred away.
	c.getOrCreateSchema(defaultSchema)
	c.sortSchemas()
	return c, nil
}
