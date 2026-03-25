package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// pgOracle wraps a real PostgreSQL container connection for oracle testing.
type pgOracle struct {
	db  *sql.DB
	ctx context.Context
}

var (
	pgOracleOnce    sync.Once
	pgOracleInst    *pgOracle
	pgOracleCleanup func()
	pgSchemaCounter atomic.Int64
)

// startPGOracle starts a shared PostgreSQL 16 container. The container is reused
// across all tests via sync.Once; individual tests get schema-level isolation.
func startPGOracle(t *testing.T) *pgOracle {
	t.Helper()
	pgOracleOnce.Do(func() {
		ctx := context.Background()
		container, err := tcpg.Run(ctx, "postgres:16-alpine",
			tcpg.WithDatabase("omni_test"),
			tcpg.WithUsername("postgres"),
			tcpg.WithPassword("test"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2)),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to start PG container: %v", err))
		}

		connStr, err := container.ConnectionString(ctx, "sslmode=disable")
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			panic(fmt.Sprintf("failed to get connection string: %v", err))
		}

		db, err := sql.Open("pgx", connStr)
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			panic(fmt.Sprintf("failed to open database: %v", err))
		}

		if err := db.PingContext(ctx); err != nil {
			db.Close()
			_ = testcontainers.TerminateContainer(container)
			panic(fmt.Sprintf("failed to ping: %v", err))
		}

		pgOracleInst = &pgOracle{db: db, ctx: ctx}
		pgOracleCleanup = func() {
			db.Close()
			_ = testcontainers.TerminateContainer(container)
		}
	})
	t.Cleanup(func() {
		// Don't cleanup here — container shared across tests.
		// Will be cleaned up when process exits.
	})
	return pgOracleInst
}

// execSQL executes a SQL statement on the PG container.
func (o *pgOracle) execSQL(t *testing.T, sqlStr string) {
	t.Helper()
	_, err := o.db.ExecContext(o.ctx, sqlStr)
	if err != nil {
		t.Fatalf("PG exec failed: %v\nSQL: %s", err, truncateSQL(sqlStr, 500))
	}
}

func truncateSQL(sqlStr string, maxLen int) string {
	if len(sqlStr) <= maxLen {
		return sqlStr
	}
	return sqlStr[:maxLen] + "..."
}

// ---------------------------------------------------------------------------
// Schema isolation helpers
// ---------------------------------------------------------------------------

// freshSchema creates a new uniquely-named schema and registers cleanup.
func (o *pgOracle) freshSchema(t *testing.T) string {
	t.Helper()
	n := pgSchemaCounter.Add(1)
	name := fmt.Sprintf("test_%d", n)
	o.execSQL(t, fmt.Sprintf("DROP SCHEMA IF EXISTS %q CASCADE", name))
	o.execSQL(t, fmt.Sprintf("CREATE SCHEMA %q", name))
	t.Cleanup(func() {
		o.execSQL(t, fmt.Sprintf("DROP SCHEMA IF EXISTS %q CASCADE", name))
	})
	return name
}

// execInSchema sets the search_path and executes DDL.
func (o *pgOracle) execInSchema(t *testing.T, schema, ddl string) {
	t.Helper()
	sqlStr := fmt.Sprintf("SET search_path TO %q, public;\n%s", schema, ddl)
	o.execSQL(t, sqlStr)
}

// ---------------------------------------------------------------------------
// Schema comparison row types
// ---------------------------------------------------------------------------

type tableRow struct {
	name     string
	relkind  string
}

type columnRow struct {
	table     string
	name      string
	dataType  string
	nullable  string
	defVal    sql.NullString
	identity  sql.NullString
	generated sql.NullString
	position  int
}

type indexRow struct {
	name       string
	definition string
}

type constraintRow struct {
	name       string
	conType    string
	table      string
	definition string
}

type functionRow struct {
	name       string
	resultType string
	argTypes   string
	body       sql.NullString
	language   string
	volatility string
	strict     bool
	security   string
	leakproof  bool
	parallel   string
}

type triggerRow struct {
	name       string
	table      string
	definition string
}

type sequenceRow struct {
	name      string
	dataType  string
	start     int64
	increment int64
	min       int64
	max       int64
	cycle     bool
}

type enumRow struct {
	name   string
	values string
}

type viewRow struct {
	name        string
	definition  string
	checkOption sql.NullString
}

type commentRow struct {
	objectType string
	objectName string
	comment    string
}

type policyRow struct {
	name       string
	table      string
	cmd        string
	permissive string
	roles      string
	usingExpr  sql.NullString
	withCheck  sql.NullString
}

// ---------------------------------------------------------------------------
// Schema query functions
// ---------------------------------------------------------------------------

func (o *pgOracle) queryTables(t *testing.T, schema string) []tableRow {
	t.Helper()
	rows, err := o.db.QueryContext(o.ctx, `
		SELECT c.relname, c.relkind::text
		FROM pg_class c
		JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE n.nspname = $1
		  AND c.relkind IN ('r', 'v', 'm', 'f', 'p')
		ORDER BY c.relname`, schema)
	if err != nil {
		t.Fatalf("queryTables: %v", err)
	}
	defer rows.Close()
	var result []tableRow
	for rows.Next() {
		var r tableRow
		if err := rows.Scan(&r.name, &r.relkind); err != nil {
			t.Fatalf("queryTables scan: %v", err)
		}
		result = append(result, r)
	}
	return result
}

func (o *pgOracle) queryColumns(t *testing.T, schema string) []columnRow {
	t.Helper()
	rows, err := o.db.QueryContext(o.ctx, `
		SELECT table_name, column_name, data_type, is_nullable,
		       column_default, identity_generation, generation_expression,
		       ordinal_position
		FROM information_schema.columns
		WHERE table_schema = $1
		ORDER BY table_name, ordinal_position`, schema)
	if err != nil {
		t.Fatalf("queryColumns: %v", err)
	}
	defer rows.Close()
	var result []columnRow
	for rows.Next() {
		var r columnRow
		if err := rows.Scan(&r.table, &r.name, &r.dataType, &r.nullable,
			&r.defVal, &r.identity, &r.generated, &r.position); err != nil {
			t.Fatalf("queryColumns scan: %v", err)
		}
		result = append(result, r)
	}
	return result
}

func (o *pgOracle) queryIndexes(t *testing.T, schema string) []indexRow {
	t.Helper()
	rows, err := o.db.QueryContext(o.ctx, `
		SELECT indexname, indexdef
		FROM pg_indexes
		WHERE schemaname = $1
		ORDER BY indexname`, schema)
	if err != nil {
		t.Fatalf("queryIndexes: %v", err)
	}
	defer rows.Close()
	var result []indexRow
	for rows.Next() {
		var r indexRow
		if err := rows.Scan(&r.name, &r.definition); err != nil {
			t.Fatalf("queryIndexes scan: %v", err)
		}
		result = append(result, r)
	}
	return result
}

func (o *pgOracle) queryConstraints(t *testing.T, schema string) []constraintRow {
	t.Helper()
	rows, err := o.db.QueryContext(o.ctx, `
		SELECT con.conname, con.contype::text,
		       c.relname,
		       pg_get_constraintdef(con.oid)
		FROM pg_constraint con
		JOIN pg_class c ON c.oid = con.conrelid
		JOIN pg_namespace n ON n.oid = con.connamespace
		WHERE n.nspname = $1
		ORDER BY c.relname, con.conname`, schema)
	if err != nil {
		t.Fatalf("queryConstraints: %v", err)
	}
	defer rows.Close()
	var result []constraintRow
	for rows.Next() {
		var r constraintRow
		if err := rows.Scan(&r.name, &r.conType, &r.table, &r.definition); err != nil {
			t.Fatalf("queryConstraints scan: %v", err)
		}
		result = append(result, r)
	}
	return result
}

func (o *pgOracle) queryFunctions(t *testing.T, schema string) []functionRow {
	t.Helper()
	rows, err := o.db.QueryContext(o.ctx, `
		SELECT p.proname,
		       pg_get_function_result(p.oid),
		       pg_get_function_arguments(p.oid),
		       p.prosrc,
		       l.lanname,
		       p.provolatile::text,
		       p.proisstrict,
		       CASE WHEN p.prosecdef THEN 'definer' ELSE 'invoker' END,
		       p.proleakproof,
		       p.proparallel::text
		FROM pg_proc p
		JOIN pg_namespace n ON n.oid = p.pronamespace
		JOIN pg_language l ON l.oid = p.prolang
		WHERE n.nspname = $1
		ORDER BY p.proname, pg_get_function_arguments(p.oid)`, schema)
	if err != nil {
		t.Fatalf("queryFunctions: %v", err)
	}
	defer rows.Close()
	var result []functionRow
	for rows.Next() {
		var r functionRow
		if err := rows.Scan(&r.name, &r.resultType, &r.argTypes,
			&r.body, &r.language, &r.volatility, &r.strict,
			&r.security, &r.leakproof, &r.parallel); err != nil {
			t.Fatalf("queryFunctions scan: %v", err)
		}
		result = append(result, r)
	}
	return result
}

func (o *pgOracle) queryTriggers(t *testing.T, schema string) []triggerRow {
	t.Helper()
	rows, err := o.db.QueryContext(o.ctx, `
		SELECT tg.tgname, c.relname,
		       pg_get_triggerdef(tg.oid)
		FROM pg_trigger tg
		JOIN pg_class c ON c.oid = tg.tgrelid
		JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE n.nspname = $1
		  AND NOT tg.tgisinternal
		ORDER BY c.relname, tg.tgname`, schema)
	if err != nil {
		t.Fatalf("queryTriggers: %v", err)
	}
	defer rows.Close()
	var result []triggerRow
	for rows.Next() {
		var r triggerRow
		if err := rows.Scan(&r.name, &r.table, &r.definition); err != nil {
			t.Fatalf("queryTriggers scan: %v", err)
		}
		result = append(result, r)
	}
	return result
}

func (o *pgOracle) querySequences(t *testing.T, schema string) []sequenceRow {
	t.Helper()
	rows, err := o.db.QueryContext(o.ctx, `
		SELECT s.sequencename, s.data_type,
		       s.start_value, s.increment_by,
		       s.min_value, s.max_value, s.cycle
		FROM pg_sequences s
		WHERE s.schemaname = $1
		ORDER BY s.sequencename`, schema)
	if err != nil {
		t.Fatalf("querySequences: %v", err)
	}
	defer rows.Close()
	var result []sequenceRow
	for rows.Next() {
		var r sequenceRow
		if err := rows.Scan(&r.name, &r.dataType, &r.start, &r.increment,
			&r.min, &r.max, &r.cycle); err != nil {
			t.Fatalf("querySequences scan: %v", err)
		}
		result = append(result, r)
	}
	return result
}

func (o *pgOracle) queryEnumTypes(t *testing.T, schema string) []enumRow {
	t.Helper()
	rows, err := o.db.QueryContext(o.ctx, `
		SELECT t.typname,
		       string_agg(e.enumlabel, ',' ORDER BY e.enumsortorder)
		FROM pg_type t
		JOIN pg_namespace n ON n.oid = t.typnamespace
		JOIN pg_enum e ON e.enumtypid = t.oid
		WHERE n.nspname = $1
		GROUP BY t.typname
		ORDER BY t.typname`, schema)
	if err != nil {
		t.Fatalf("queryEnumTypes: %v", err)
	}
	defer rows.Close()
	var result []enumRow
	for rows.Next() {
		var r enumRow
		if err := rows.Scan(&r.name, &r.values); err != nil {
			t.Fatalf("queryEnumTypes scan: %v", err)
		}
		result = append(result, r)
	}
	return result
}

func (o *pgOracle) queryViews(t *testing.T, schema string) []viewRow {
	t.Helper()
	rows, err := o.db.QueryContext(o.ctx, `
		SELECT viewname, definition, NULL::text
		FROM pg_views
		WHERE schemaname = $1
		ORDER BY viewname`, schema)
	if err != nil {
		t.Fatalf("queryViews: %v", err)
	}
	defer rows.Close()
	var result []viewRow
	for rows.Next() {
		var r viewRow
		if err := rows.Scan(&r.name, &r.definition, &r.checkOption); err != nil {
			t.Fatalf("queryViews scan: %v", err)
		}
		result = append(result, r)
	}
	return result
}

func (o *pgOracle) queryComments(t *testing.T, schema string) []commentRow {
	t.Helper()
	rows, err := o.db.QueryContext(o.ctx, `
		SELECT
			CASE c.classoid
				WHEN 'pg_class'::regclass THEN
					CASE rel.relkind
						WHEN 'r' THEN 'TABLE'
						WHEN 'v' THEN 'VIEW'
						WHEN 'm' THEN 'MATERIALIZED VIEW'
						WHEN 'S' THEN 'SEQUENCE'
						WHEN 'i' THEN 'INDEX'
						ELSE 'RELATION'
					END
				WHEN 'pg_type'::regclass THEN 'TYPE'
				WHEN 'pg_proc'::regclass THEN 'FUNCTION'
				WHEN 'pg_namespace'::regclass THEN 'SCHEMA'
				ELSE c.classoid::regclass::text
			END AS object_type,
			CASE c.classoid
				WHEN 'pg_class'::regclass THEN rel.relname
				WHEN 'pg_type'::regclass THEN typ.typname
				WHEN 'pg_proc'::regclass THEN pro.proname
				WHEN 'pg_namespace'::regclass THEN nsp.nspname
				ELSE c.objoid::text
			END AS object_name,
			d.description
		FROM pg_description d
		JOIN (
			SELECT objoid, classoid,
			       CASE WHEN classoid = 'pg_class'::regclass THEN (SELECT relnamespace FROM pg_class WHERE oid = objoid)
			            WHEN classoid = 'pg_type'::regclass THEN (SELECT typnamespace FROM pg_type WHERE oid = objoid)
			            WHEN classoid = 'pg_proc'::regclass THEN (SELECT pronamespace FROM pg_proc WHERE oid = objoid)
			            WHEN classoid = 'pg_namespace'::regclass THEN objoid
			            ELSE 0
			       END AS nsoid
			FROM pg_description
			WHERE objsubid = 0
		) c ON c.objoid = d.objoid AND c.classoid = d.classoid
		LEFT JOIN pg_class rel ON d.classoid = 'pg_class'::regclass AND d.objoid = rel.oid
		LEFT JOIN pg_type typ ON d.classoid = 'pg_type'::regclass AND d.objoid = typ.oid
		LEFT JOIN pg_proc pro ON d.classoid = 'pg_proc'::regclass AND d.objoid = pro.oid
		LEFT JOIN pg_namespace nsp ON d.classoid = 'pg_namespace'::regclass AND d.objoid = nsp.oid
		JOIN pg_namespace fns ON fns.oid = c.nsoid
		WHERE fns.nspname = $1
		  AND d.objsubid = 0
		ORDER BY object_type, object_name`, schema)
	if err != nil {
		t.Fatalf("queryComments: %v", err)
	}
	defer rows.Close()
	var result []commentRow
	for rows.Next() {
		var r commentRow
		if err := rows.Scan(&r.objectType, &r.objectName, &r.comment); err != nil {
			t.Fatalf("queryComments scan: %v", err)
		}
		result = append(result, r)
	}
	return result
}

func (o *pgOracle) queryPolicies(t *testing.T, schema string) []policyRow {
	t.Helper()
	rows, err := o.db.QueryContext(o.ctx, `
		SELECT pol.polname, c.relname,
		       pol.polcmd::text,
		       CASE WHEN pol.polpermissive THEN 'PERMISSIVE' ELSE 'RESTRICTIVE' END,
		       array_to_string(
		           ARRAY(SELECT rolname FROM pg_roles WHERE oid = ANY(pol.polroles)), ','
		       ),
		       pg_get_expr(pol.polqual, pol.polrelid),
		       pg_get_expr(pol.polwithcheck, pol.polrelid)
		FROM pg_policy pol
		JOIN pg_class c ON c.oid = pol.polrelid
		JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE n.nspname = $1
		ORDER BY c.relname, pol.polname`, schema)
	if err != nil {
		t.Fatalf("queryPolicies: %v", err)
	}
	defer rows.Close()
	var result []policyRow
	for rows.Next() {
		var r policyRow
		if err := rows.Scan(&r.name, &r.table, &r.cmd, &r.permissive,
			&r.roles, &r.usingExpr, &r.withCheck); err != nil {
			t.Fatalf("queryPolicies scan: %v", err)
		}
		result = append(result, r)
	}
	return result
}

// ---------------------------------------------------------------------------
// Schema comparison
// ---------------------------------------------------------------------------

func (o *pgOracle) compareSchemas(t *testing.T, schemaA, schemaB string) []string {
	t.Helper()
	var diffs []string

	// Compare tables
	tablesA := o.queryTables(t, schemaA)
	tablesB := o.queryTables(t, schemaB)
	diffs = append(diffs, compareSlices("tables", tablesA, tablesB,
		func(r tableRow) string { return r.name },
		func(r tableRow) string { return fmt.Sprintf("%s (kind=%s)", r.name, r.relkind) })...)

	// Compare columns
	colsA := o.queryColumns(t, schemaA)
	colsB := o.queryColumns(t, schemaB)
	diffs = append(diffs, compareSlices("columns", colsA, colsB,
		func(r columnRow) string { return fmt.Sprintf("%s.%s", r.table, r.name) },
		func(r columnRow) string {
			return fmt.Sprintf("%s.%s type=%s nullable=%s default=%v identity=%v generated=%v pos=%d",
				r.table, r.name, r.dataType, r.nullable, r.defVal, r.identity, r.generated, r.position)
		})...)

	// Compare indexes
	idxA := o.queryIndexes(t, schemaA)
	idxB := o.queryIndexes(t, schemaB)
	diffs = append(diffs, compareIndexes(schemaA, schemaB, idxA, idxB)...)

	// Compare constraints
	consA := o.queryConstraints(t, schemaA)
	consB := o.queryConstraints(t, schemaB)
	diffs = append(diffs, compareSlices("constraints", consA, consB,
		func(r constraintRow) string { return fmt.Sprintf("%s.%s", r.table, r.name) },
		func(r constraintRow) string {
			return fmt.Sprintf("%s.%s type=%s def=%s", r.table, r.name, r.conType, r.definition)
		})...)

	// Compare functions
	funcsA := o.queryFunctions(t, schemaA)
	funcsB := o.queryFunctions(t, schemaB)
	diffs = append(diffs, compareSlices("functions", funcsA, funcsB,
		func(r functionRow) string { return fmt.Sprintf("%s(%s)", r.name, r.argTypes) },
		func(r functionRow) string {
			return fmt.Sprintf("%s(%s) returns=%s lang=%s vol=%s strict=%v sec=%s leak=%v par=%s body=%v",
				r.name, r.argTypes, r.resultType, r.language, r.volatility,
				r.strict, r.security, r.leakproof, r.parallel, r.body)
		})...)

	// Compare triggers
	trigsA := o.queryTriggers(t, schemaA)
	trigsB := o.queryTriggers(t, schemaB)
	diffs = append(diffs, compareTriggers(schemaA, schemaB, trigsA, trigsB)...)

	// Compare sequences
	seqsA := o.querySequences(t, schemaA)
	seqsB := o.querySequences(t, schemaB)
	diffs = append(diffs, compareSlices("sequences", seqsA, seqsB,
		func(r sequenceRow) string { return r.name },
		func(r sequenceRow) string {
			return fmt.Sprintf("%s type=%s start=%d inc=%d min=%d max=%d cycle=%v",
				r.name, r.dataType, r.start, r.increment, r.min, r.max, r.cycle)
		})...)

	// Compare enum types
	enumsA := o.queryEnumTypes(t, schemaA)
	enumsB := o.queryEnumTypes(t, schemaB)
	diffs = append(diffs, compareSlices("enums", enumsA, enumsB,
		func(r enumRow) string { return r.name },
		func(r enumRow) string { return fmt.Sprintf("%s values=[%s]", r.name, r.values) })...)

	// Compare views
	viewsA := o.queryViews(t, schemaA)
	viewsB := o.queryViews(t, schemaB)
	diffs = append(diffs, compareSlices("views", viewsA, viewsB,
		func(r viewRow) string { return r.name },
		func(r viewRow) string {
			return fmt.Sprintf("%s def=%s check=%v", r.name, r.definition, r.checkOption)
		})...)

	// Compare comments
	cmtsA := o.queryComments(t, schemaA)
	cmtsB := o.queryComments(t, schemaB)
	diffs = append(diffs, compareSlices("comments", cmtsA, cmtsB,
		func(r commentRow) string { return fmt.Sprintf("%s.%s", r.objectType, r.objectName) },
		func(r commentRow) string {
			return fmt.Sprintf("%s %s: %s", r.objectType, r.objectName, r.comment)
		})...)

	// Compare policies
	polsA := o.queryPolicies(t, schemaA)
	polsB := o.queryPolicies(t, schemaB)
	diffs = append(diffs, compareSlices("policies", polsA, polsB,
		func(r policyRow) string { return fmt.Sprintf("%s.%s", r.table, r.name) },
		func(r policyRow) string {
			return fmt.Sprintf("%s.%s cmd=%s perm=%s roles=%s using=%v check=%v",
				r.table, r.name, r.cmd, r.permissive, r.roles, r.usingExpr, r.withCheck)
		})...)

	return diffs
}

// compareSlices is a generic helper that compares two sorted slices by key and detail string.
func compareSlices[T any](category string, a, b []T, keyFn func(T) string, detailFn func(T) string) []string {
	var diffs []string
	mapA := make(map[string]string, len(a))
	mapB := make(map[string]string, len(b))
	var keysA, keysB []string
	for _, r := range a {
		k := keyFn(r)
		mapA[k] = detailFn(r)
		keysA = append(keysA, k)
	}
	for _, r := range b {
		k := keyFn(r)
		mapB[k] = detailFn(r)
		keysB = append(keysB, k)
	}
	sort.Strings(keysA)
	sort.Strings(keysB)

	allKeys := make(map[string]bool)
	for _, k := range keysA {
		allKeys[k] = true
	}
	for _, k := range keysB {
		allKeys[k] = true
	}
	sorted := make([]string, 0, len(allKeys))
	for k := range allKeys {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)

	for _, k := range sorted {
		dA, okA := mapA[k]
		dB, okB := mapB[k]
		if okA && !okB {
			diffs = append(diffs, fmt.Sprintf("%s: only in schemaA: %s", category, dA))
		} else if !okA && okB {
			diffs = append(diffs, fmt.Sprintf("%s: only in schemaB: %s", category, dB))
		} else if dA != dB {
			diffs = append(diffs, fmt.Sprintf("%s: differ:\n  A: %s\n  B: %s", category, dA, dB))
		}
	}
	return diffs
}

// compareIndexes compares indexes, normalizing schema names in index definitions.
func compareIndexes(schemaA, schemaB string, a, b []indexRow) []string {
	// Normalize index definitions by replacing schema-specific references.
	// PG may output schema names quoted or unquoted depending on the name.
	normalize := func(def, schema string) string {
		def = strings.ReplaceAll(def, fmt.Sprintf("%q.", schema), "")
		def = strings.ReplaceAll(def, schema+".", "")
		return def
	}
	var diffs []string
	mapA := make(map[string]string, len(a))
	mapB := make(map[string]string, len(b))
	for _, r := range a {
		mapA[r.name] = normalize(r.definition, schemaA)
	}
	for _, r := range b {
		mapB[r.name] = normalize(r.definition, schemaB)
	}

	allKeys := make(map[string]bool)
	for k := range mapA {
		allKeys[k] = true
	}
	for k := range mapB {
		allKeys[k] = true
	}
	sorted := make([]string, 0, len(allKeys))
	for k := range allKeys {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)

	for _, k := range sorted {
		dA, okA := mapA[k]
		dB, okB := mapB[k]
		if okA && !okB {
			diffs = append(diffs, fmt.Sprintf("indexes: only in schemaA: %s = %s", k, dA))
		} else if !okA && okB {
			diffs = append(diffs, fmt.Sprintf("indexes: only in schemaB: %s = %s", k, dB))
		} else if dA != dB {
			diffs = append(diffs, fmt.Sprintf("indexes: differ:\n  A: %s = %s\n  B: %s = %s", k, dA, k, dB))
		}
	}
	return diffs
}

// compareTriggers compares triggers, normalizing schema names in definitions.
func compareTriggers(schemaA, schemaB string, a, b []triggerRow) []string {
	normalize := func(def, schema string) string {
		def = strings.ReplaceAll(def, fmt.Sprintf("%q.", schema), "")
		def = strings.ReplaceAll(def, schema+".", "")
		return def
	}
	var diffs []string
	mapA := make(map[string]string, len(a))
	mapB := make(map[string]string, len(b))
	for _, r := range a {
		k := fmt.Sprintf("%s.%s", r.table, r.name)
		mapA[k] = normalize(r.definition, schemaA)
	}
	for _, r := range b {
		k := fmt.Sprintf("%s.%s", r.table, r.name)
		mapB[k] = normalize(r.definition, schemaB)
	}

	allKeys := make(map[string]bool)
	for k := range mapA {
		allKeys[k] = true
	}
	for k := range mapB {
		allKeys[k] = true
	}
	sorted := make([]string, 0, len(allKeys))
	for k := range allKeys {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)

	for _, k := range sorted {
		dA, okA := mapA[k]
		dB, okB := mapB[k]
		if okA && !okB {
			diffs = append(diffs, fmt.Sprintf("triggers: only in schemaA: %s = %s", k, dA))
		} else if !okA && okB {
			diffs = append(diffs, fmt.Sprintf("triggers: only in schemaB: %s = %s", k, dB))
		} else if dA != dB {
			diffs = append(diffs, fmt.Sprintf("triggers: differ:\n  A: %s = %s\n  B: %s = %s", k, dA, k, dB))
		}
	}
	return diffs
}

func (o *pgOracle) assertSchemasEqual(t *testing.T, schemaA, schemaB string) {
	t.Helper()
	diffs := o.compareSchemas(t, schemaA, schemaB)
	if len(diffs) > 0 {
		t.Errorf("schemas %s and %s differ:\n%s", schemaA, schemaB, strings.Join(diffs, "\n"))
	}
}

// ---------------------------------------------------------------------------
// Core roundtrip helper
// ---------------------------------------------------------------------------

// assertOracleRoundtrip tests that omni's generated migration produces the same
// schema state as applying the "after" DDL directly.
func assertOracleRoundtrip(t *testing.T, o *pgOracle, beforeDDL, afterDDL string) {
	t.Helper()

	migrated := o.freshSchema(t)
	expected := o.freshSchema(t)

	// 1. Apply "before" to the migrated schema.
	if beforeDDL != "" {
		o.execInSchema(t, migrated, beforeDDL)
	}

	// 2. Generate migration via omni.
	fromCat := New()
	if beforeDDL != "" {
		var err error
		fromCat, err = LoadSQL(beforeDDL)
		if err != nil {
			t.Fatalf("LoadSQL(before) failed: %v", err)
		}
	}
	toCat := New()
	if afterDDL != "" {
		var err error
		toCat, err = LoadSQL(afterDDL)
		if err != nil {
			t.Fatalf("LoadSQL(after) failed: %v", err)
		}
	}
	diff := Diff(fromCat, toCat)
	plan := GenerateMigration(fromCat, toCat, diff)

	// 3. Apply migration to PG.
	// Strip explicit "public." schema qualifiers since we're running in a test schema.
	migrationSQL := plan.SQL()
	migrationSQL = strings.ReplaceAll(migrationSQL, "public.", "")
	migrationSQL = strings.ReplaceAll(migrationSQL, `"public".`, "")
	if migrationSQL != "" {
		o.execInSchema(t, migrated, migrationSQL)
	}

	// 4. Apply "after" directly to expected schema.
	if afterDDL != "" {
		o.execInSchema(t, expected, afterDDL)
	}

	// 5. Compare schemas.
	o.assertSchemasEqual(t, migrated, expected)
}
