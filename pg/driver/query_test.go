package driver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/bytebase/omni/driver"
	pganalysis "github.com/bytebase/omni/pg/analysis"
	"github.com/bytebase/omni/pg/internal/testcontainer"
)

func TestFormatTimestamp(t *testing.T) {
	cases := []struct {
		name     string
		input    time.Time
		withZone bool
		want     string
	}{
		{
			name:     "BC timestamp without zone",
			input:    time.Date(0, 12, 31, 15, 0, 0, 0, time.UTC),
			withZone: false,
			want:     "0001-12-31 15:00:00 BC",
		},
		{
			name:     "year 10000 timestamp without zone",
			input:    time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC),
			withZone: false,
			want:     "10000-01-01 00:00:00",
		},
		{
			name:     "year 10000 timestamptz with offset",
			input:    time.Date(10000, 5, 11, 12, 34, 56, 0, time.FixedZone("", 2*60*60)),
			withZone: true,
			want:     "10000-05-11 12:34:56+02:00",
		},
	}
	for _, tc := range cases {
		if got := FormatTimestamp(tc.input, tc.withZone); got != tc.want {
			t.Errorf("%s: FormatTimestamp = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestPadZeroes(t *testing.T) {
	cases := []struct {
		raw  string
		acc  int
		want string
	}{
		{"12:34:56", 6, "12:34:56"},
		{"12:34:56.123456", 6, "12:34:56.123456"},
		{"12:34:56.123", 6, "12:34:56.123000"},
		{"12:34:56.123456789", 6, "12:34:56.123456789"},
		{"12:34:56.123+02:00", 6, "12:34:56.123000+02:00"},
		{"12:34:56.123-02:00", 6, "12:34:56.123000-02:00"},
		{"12:34-05:56.123", 6, "12:34-05:56.123000"},
		{"-00:04:37.530865", 6, "-00:04:37.530865"},
		{"-00:02:45.25", 6, "-00:02:45.250000"},
		{"-00:05:00", 6, "-00:05:00"},
		{"-12:34:56.123-05:00", 6, "-12:34:56.123-05:00"},
	}
	for _, tc := range cases {
		if got := padZeroes(tc.raw, tc.acc); got != tc.want {
			t.Errorf("padZeroes(%q, %d) = %q, want %q", tc.raw, tc.acc, got, tc.want)
		}
	}
}

func openQueryDriver(t *testing.T, role string) *Driver {
	t.Helper()
	db := testcontainer.OpenRole(t, role)
	t.Cleanup(func() { _ = db.Close() })
	return &Driver{db: db, databaseName: testcontainer.Database(t)}
}

// TestQuery exercises Session.Query across every pinned PostgreSQL release:
// value decoding, zero-column SELECTs vs row-count statements, session
// persistence, search_path escaping and public fallback, limit rewriting,
// explain plans, result-size truncation, and error halting.
func TestQuery(t *testing.T) {
	ctx := context.Background()
	for _, role := range versions {
		t.Run(role, func(t *testing.T) {
			d := openQueryDriver(t, role)
			// tape.Open caps the pool at one connection so each stream is
			// deterministic; run fixture DDL before opening Session.
			mustExec(t, d.db, `
				CREATE SCHEMA app;
				CREATE SCHEMA "app""schema";
				CREATE TABLE public.lookup_precedence (marker text);
				CREATE TABLE app.lookup_precedence (marker text);
				CREATE TABLE "app""schema".lookup_precedence (marker text);
				INSERT INTO public.lookup_precedence VALUES ('public');
				INSERT INTO app.lookup_precedence VALUES ('app');
				INSERT INTO "app""schema".lookup_precedence VALUES ('quoted');
				CREATE FUNCTION public.helper_echo(rel text) RETURNS text LANGUAGE sql AS $$ SELECT rel $$;
				CREATE FUNCTION public.call_unqualified_helper(rel text) RETURNS text LANGUAGE plpgsql AS $$
				BEGIN
					RETURN helper_echo(rel);
				END; $$;
				CREATE TABLE plan_target (id int PRIMARY KEY, note text);
				INSERT INTO plan_target SELECT g, 'n' || g FROM generate_series(1, 200) g;
				ANALYZE plan_target;
			`)

			sess, err := d.Session(ctx)
			if err != nil {
				t.Fatal(err)
			}
			defer sess.Close()

			// Empty or comment-only input sends nothing.
			empty, err := sess.Query(ctx, "  ; -- comment only\n ; ", driver.QueryOptions{})
			if err != nil || len(empty) != 0 {
				t.Fatalf("empty script = %v, %v", empty, err)
			}

			// Scalar and special PostgreSQL values.
			const typesSQL = `SELECT
				true AS b,
				42::int4 AS i4,
				9007199254740993::int8 AS i8,
				3.5::float8 AS f8,
				12.340::numeric(6,3) AS num,
				'hello'::text AS txt,
				'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'::uuid AS u,
				'\xdeadbeef'::bytea AS bin,
				''::bytea AS empty_bin,
				'2026-03-04'::date AS dt,
				'infinity'::date AS dt_inf,
				'2026-03-04 05:06:07.891234'::timestamp AS ts,
				'-infinity'::timestamp AS ts_neginf,
				'2026-03-04 05:06:07.891234+00'::timestamptz AS tstz,
				'infinity'::timestamptz AS tstz_inf,
				'12:34:56.12'::time AS tm,
				'-00:02:45.25'::interval AS iv,
				'{"k": 1}'::jsonb AS jb,
				NULL::text AS nul;`
			res, err := sess.Query(ctx, typesSQL, driver.QueryOptions{})
			if err != nil || len(res) != 1 || res[0].Err != nil {
				t.Fatalf("types query = %+v, %v", res, err)
			}
			if res[0].AffectedRows != nil || len(res[0].Rows) != 1 {
				t.Fatalf("types result shape = %+v", res[0])
			}
			row := res[0].Rows[0]
			checkVal := func(idx int, wantKind driver.Kind) driver.Value {
				t.Helper()
				if row[idx].Kind != wantKind {
					t.Fatalf("col %d (%s): kind = %s, want %s (%+v)", idx, res[0].Columns[idx].Name, row[idx].Kind, wantKind, row[idx])
				}
				return row[idx]
			}
			if !checkVal(0, driver.Bool).Bool {
				t.Errorf("bool = false")
			}
			if checkVal(1, driver.Int64).Int64 != 42 {
				t.Errorf("i4 = %d", row[1].Int64)
			}
			if checkVal(2, driver.Int64).Int64 != 9007199254740993 {
				t.Errorf("i8 = %d", row[2].Int64)
			}
			if checkVal(3, driver.Float64).Float64 != 3.5 {
				t.Errorf("f8 = %v", row[3].Float64)
			}
			if checkVal(4, driver.String).String != "12.340" {
				t.Errorf("num = %q", row[4].String)
			}
			if checkVal(5, driver.String).String != "hello" {
				t.Errorf("txt = %q", row[5].String)
			}
			if checkVal(6, driver.String).String != "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11" {
				t.Errorf("uuid = %q", row[6].String)
			}
			if !bytes.Equal(checkVal(7, driver.Bytes).Bytes, []byte{0xde, 0xad, 0xbe, 0xef}) {
				t.Errorf("bin = %x", row[7].Bytes)
			}
			// Parity with Bytebase db/pg: empty bytea decodes as Null.
			checkVal(8, driver.Null)
			if checkVal(9, driver.String).String != "2026-03-04" {
				t.Errorf("dt = %q", row[9].String)
			}
			if checkVal(10, driver.String).String != "infinity" {
				t.Errorf("dt_inf = %q", row[10].String)
			}
			if checkVal(11, driver.Timestamp).Time.Format(time.RFC3339Nano) != "2026-03-04T05:06:07.891234Z" {
				t.Errorf("ts = %v", row[11].Time)
			}
			if checkVal(12, driver.String).String != "-infinity" {
				t.Errorf("ts_neginf = %q", row[12].String)
			}
			if checkVal(13, driver.TimestampTZ).Time.UTC().Format(time.RFC3339Nano) != "2026-03-04T05:06:07.891234Z" {
				t.Errorf("tstz = %v", row[13].Time)
			}
			if checkVal(14, driver.String).String != "infinity" {
				t.Errorf("tstz_inf = %q", row[14].String)
			}
			if checkVal(15, driver.String).String != "12:34:56.120000" {
				t.Errorf("tm = %q", row[15].String)
			}
			if checkVal(16, driver.String).String != "-00:02:45.250000" {
				t.Errorf("iv = %q", row[16].String)
			}
			if checkVal(17, driver.String).String != `{"k": 1}` {
				t.Errorf("jb = %q", row[17].String)
			}
			checkVal(18, driver.Null)

			// Session state persists across Query calls; zero-column SELECT has
			// AffectedRows == nil while DML/DDL report AffectedRows != nil.
			setupRes, err := sess.Query(ctx, `
				CREATE TEMP TABLE sess_items (id int PRIMARY KEY, note text);
				INSERT INTO sess_items SELECT g, repeat('x', 40) FROM generate_series(1, 10) g;
				SELECT FROM sess_items;`, driver.QueryOptions{})
			if err != nil || len(setupRes) != 3 {
				t.Fatalf("setup query = %+v, %v", setupRes, err)
			}
			if setupRes[0].AffectedRows == nil || *setupRes[0].AffectedRows != 0 {
				t.Errorf("CREATE TEMP TABLE AffectedRows = %v", setupRes[0].AffectedRows)
			}
			if setupRes[1].AffectedRows == nil || *setupRes[1].AffectedRows != 10 {
				t.Errorf("INSERT AffectedRows = %v", setupRes[1].AffectedRows)
			}
			if setupRes[2].AffectedRows != nil || len(setupRes[2].Columns) != 0 {
				t.Errorf("zero-column SELECT = %+v", setupRes[2])
			}

			// Limit rewrite and LIMIT ALL CTE fallback.
			limRes, err := sess.Query(ctx, `
				SELECT id FROM sess_items ORDER BY id;
				SELECT id FROM sess_items ORDER BY id LIMIT 2;
				SELECT id FROM sess_items ORDER BY id LIMIT ALL;`, driver.QueryOptions{Limit: 4})
			if err != nil || len(limRes) != 3 {
				t.Fatalf("limit query = %+v, %v", limRes, err)
			}
			if len(limRes[0].Rows) != 4 || len(limRes[1].Rows) != 2 || len(limRes[2].Rows) != 4 {
				t.Errorf("limit row counts = %d, %d, %d; want 4, 2, 4", len(limRes[0].Rows), len(limRes[1].Rows), len(limRes[2].Rows))
			}

			// MaxResultBytes truncates the oversized statement and still runs the next.
			truncRes, err := sess.Query(ctx, "SELECT id, note FROM sess_items ORDER BY id;\nSELECT 7 AS after_trunc;", driver.QueryOptions{
				MaxResultBytes: 120,
			})
			if err != nil || len(truncRes) != 2 {
				t.Fatalf("truncation query = %+v, %v", truncRes, err)
			}
			if !truncRes[0].Truncated || len(truncRes[0].Rows) >= 10 {
				t.Errorf("truncRes[0] Truncated=%v rows=%d", truncRes[0].Truncated, len(truncRes[0].Rows))
			}
			if truncRes[1].Truncated || len(truncRes[1].Rows) != 1 || truncRes[1].Rows[0][0].Int64 != 7 {
				t.Errorf("truncRes[1] = %+v", truncRes[1])
			}

			// Error stops the script and preserves *pgconn.PgError fields.
			errRes, err := sess.Query(ctx, "SELECT 1;\nSELECT 1 / 0;\nSELECT 2;", driver.QueryOptions{})
			if err != nil || len(errRes) != 2 {
				t.Fatalf("error script = %+v, %v", errRes, err)
			}
			var pgErr *pgconn.PgError
			if errRes[0].Err != nil || !errors.As(errRes[1].Err, &pgErr) || pgErr.Code != "22012" {
				t.Fatalf("errRes = %+v", errRes)
			}

			// Schema search_path keeps public after the selected schema and escapes quotes.
			spRes, err := sess.Query(ctx, `
				SELECT public.call_unqualified_helper('ok'::text);
				SELECT marker FROM lookup_precedence;`, driver.QueryOptions{Schema: "app"})
			if err != nil || len(spRes) != 2 || spRes[0].Rows[0][0].String != "ok" || spRes[1].Rows[0][0].String != "app" {
				t.Fatalf("app search_path = %+v, %v", spRes, err)
			}
			quotedRes, err := sess.Query(ctx, `SELECT marker FROM lookup_precedence;`, driver.QueryOptions{Schema: `app"schema`})
			if err != nil || len(quotedRes) != 1 || quotedRes[0].Rows[0][0].String != "quoted" {
				t.Fatalf("quoted search_path = %+v, %v", quotedRes, err)
			}

			// Explain formats, rebuilding a caller's EXPLAIN, and refusing EXPLAIN ANALYZE.
			textPlan, err := sess.Query(ctx, "SELECT note FROM plan_target WHERE id > 50;", driver.QueryOptions{Explain: true})
			if err != nil || len(textPlan) != 1 || textPlan[0].Plan == nil || textPlan[0].Plan.Format != "TEXT" || textPlan[0].Plan.Executed {
				t.Fatalf("text explain = %+v, %v", textPlan, err)
			}
			if !strings.Contains(textPlan[0].Rows[0][0].String, "(cost=") {
				t.Errorf("text plan = %q", textPlan[0].Rows[0][0].String)
			}

			jsonPlan, err := sess.Query(ctx, "EXPLAIN (FORMAT TEXT) SELECT note FROM plan_target WHERE id > 50;", driver.QueryOptions{
				Explain:       true,
				ExplainFormat: "JSON",
			})
			if err != nil || len(jsonPlan) != 1 || jsonPlan[0].Plan == nil || jsonPlan[0].Plan.Format != "JSON" {
				t.Fatalf("json explain = %+v, %v", jsonPlan, err)
			}
			var parsedPlan []struct {
				Plan struct {
					RelationName string `json:"Relation Name"`
				} `json:"Plan"`
			}
			if err := json.Unmarshal([]byte(jsonPlan[0].Rows[0][0].String), &parsedPlan); err != nil || len(parsedPlan) != 1 || parsedPlan[0].Plan.RelationName != "plan_target" {
				t.Fatalf("json plan body = %q, err = %v", jsonPlan[0].Rows[0][0].String, err)
			}

			typedPlan, err := sess.Query(ctx, "EXPLAIN (FORMAT YAML) SELECT note FROM plan_target WHERE id > 50;", driver.QueryOptions{})
			if err != nil || len(typedPlan) != 1 || typedPlan[0].Plan == nil || typedPlan[0].Plan.Format != "YAML" || typedPlan[0].Plan.Executed {
				t.Fatalf("typed explain = %+v, %v", typedPlan, err)
			}

			if _, err := sess.Query(ctx, "EXPLAIN ANALYZE SELECT 1;", driver.QueryOptions{Explain: true}); !errors.Is(err, pganalysis.ErrExplainAnalyzeExecutes) {
				t.Errorf("EXPLAIN ANALYZE under Explain=true err = %v, want ErrExplainAnalyzeExecutes", err)
			}
			if _, err := sess.Query(ctx, "SELECT FROM WHERE;", driver.QueryOptions{}); err == nil {
				t.Errorf("syntax error returned nil")
			}
		})
	}
}

// TestQueryMergeCTEDoesNotLeakRows pins the security boundary on PostgreSQL 17+:
// a MERGE ... RETURNING inside a CTE is classified as a write, so Session.Query
// routes it to ExecContext and never returns RETURNING rows to a masking caller.
func TestQueryMergeCTEDoesNotLeakRows(t *testing.T) {
	ctx := context.Background()
	const secret = "TOP-SECRET-SSN-123-45-6789"
	for _, role := range []string{testcontainer.Default, testcontainer.PG18} {
		t.Run(role, func(t *testing.T) {
			d := openQueryDriver(t, role)
			mustExec(t, d.db, `
				CREATE TABLE secret_doc (id int PRIMARY KEY, secret text);
				INSERT INTO secret_doc VALUES (1, '`+secret+`');
			`)
			sess, err := d.Session(ctx)
			if err != nil {
				t.Fatal(err)
			}
			defer sess.Close()

			baseRes, err := sess.Query(ctx, "SELECT secret FROM secret_doc;", driver.QueryOptions{Limit: 5000})
			if err != nil || len(baseRes) != 1 || baseRes[0].Rows[0][0].String != secret {
				t.Fatalf("baseline SELECT = %+v, %v", baseRes, err)
			}

			const mergeCTE = `WITH x AS (
				MERGE INTO secret_doc t USING secret_doc s ON t.id = s.id
				WHEN MATCHED THEN UPDATE SET id = t.id
				RETURNING t.secret AS secret
			) SELECT secret FROM x;`
			mergeRes, err := sess.Query(ctx, mergeCTE, driver.QueryOptions{Limit: 5000})
			if err != nil || len(mergeRes) != 1 {
				t.Fatalf("merge CTE = %+v, %v", mergeRes, err)
			}
			if mergeRes[0].AffectedRows == nil || len(mergeRes[0].Rows) != 0 {
				t.Fatalf("SECURITY: MERGE-in-CTE returned row data: %+v", mergeRes[0])
			}
		})
	}
}
