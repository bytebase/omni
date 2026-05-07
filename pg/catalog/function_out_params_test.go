package catalog

import "testing"

func TestCreateFunctionInfersReturnTypeFromSingleOutParam(t *testing.T) {
	c := New()

	execSQL(t, c, `
CREATE FUNCTION f_out_scalar(IN input_value integer, OUT output_value text)
LANGUAGE plpgsql
AS 'BEGIN output_value := input_value::text; END';
`)

	procs := c.LookupProcByName("f_out_scalar")
	if len(procs) != 1 {
		t.Fatalf("procs: got %d, want 1", len(procs))
	}
	up := c.userProcs[procs[0].OID]
	if up == nil {
		t.Fatal("user proc not found")
	}
	if got, want := string(up.ArgModes), "io"; got != want {
		t.Fatalf("arg modes: got %q, want %q", got, want)
	}
	if got, want := len(up.AllArgTypes), 2; got != want {
		t.Fatalf("all arg types count: got %d, want %d", got, want)
	}
	if procs[0].RetType != TEXTOID {
		t.Fatalf("ret type: got %d, want %d", procs[0].RetType, TEXTOID)
	}
	if procs[0].NArgs != 1 {
		t.Fatalf("nargs: got %d, want 1", procs[0].NArgs)
	}
}

func TestCreateFunctionInfersRecordReturnTypeFromMultipleOutParams(t *testing.T) {
	c := New()

	execSQL(t, c, `
CREATE FUNCTION f_out_record(IN input_value integer, OUT doubled integer, OUT label text)
LANGUAGE plpgsql
AS 'BEGIN doubled := input_value * 2; label := input_value::text; END';
`)

	procs := c.LookupProcByName("f_out_record")
	if len(procs) != 1 {
		t.Fatalf("procs: got %d, want 1", len(procs))
	}
	if procs[0].RetType != RECORDOID {
		t.Fatalf("ret type: got %d, want %d", procs[0].RetType, RECORDOID)
	}
	if procs[0].NArgs != 1 {
		t.Fatalf("nargs: got %d, want 1", procs[0].NArgs)
	}
}

func TestLoadSQLViewFromRecordReturningFunctionWithOutParams(t *testing.T) {
	c, err := LoadSQL(`
CREATE FUNCTION public.record_pair(OUT id integer, OUT name text)
RETURNS record
LANGUAGE sql
AS $$
  SELECT 1::integer, 'alice'::text;
$$;

CREATE VIEW public.record_pair_view AS
SELECT *
FROM public.record_pair();
`)
	if err != nil {
		t.Fatalf("LoadSQL() error: %v", err)
	}

	rel := c.GetRelation("public", "record_pair_view")
	if rel == nil {
		t.Fatal("view record_pair_view not found")
	}
	if got, want := len(rel.Columns), 2; got != want {
		t.Fatalf("view columns: got %d, want %d", got, want)
	}
	if got, want := rel.Columns[0].Name, "id"; got != want {
		t.Fatalf("first column name: got %q, want %q", got, want)
	}
	if got, want := rel.Columns[0].TypeOID, INT4OID; got != want {
		t.Fatalf("first column type: got %d, want %d", got, want)
	}
	if got, want := rel.Columns[1].Name, "name"; got != want {
		t.Fatalf("second column name: got %q, want %q", got, want)
	}
	if got, want := rel.Columns[1].TypeOID, TEXTOID; got != want {
		t.Fatalf("second column type: got %d, want %d", got, want)
	}
}

func TestLoadSQLViewFromRecordReturningFunctionWithInOutParam(t *testing.T) {
	c, err := LoadSQL(`
CREATE FUNCTION public.record_inout(INOUT id integer, OUT name text)
RETURNS record
LANGUAGE sql
AS $$
  SELECT id, 'alice'::text;
$$;

CREATE VIEW public.record_inout_view AS
SELECT *
FROM public.record_inout(1);
`)
	if err != nil {
		t.Fatalf("LoadSQL() error: %v", err)
	}

	rel := c.GetRelation("public", "record_inout_view")
	if rel == nil {
		t.Fatal("view record_inout_view not found")
	}
	if got, want := len(rel.Columns), 2; got != want {
		t.Fatalf("view columns: got %d, want %d", got, want)
	}
	if got, want := rel.Columns[0].Name, "id"; got != want {
		t.Fatalf("first column name: got %q, want %q", got, want)
	}
	if got, want := rel.Columns[0].TypeOID, INT4OID; got != want {
		t.Fatalf("first column type: got %d, want %d", got, want)
	}
	if got, want := rel.Columns[1].Name, "name"; got != want {
		t.Fatalf("second column name: got %q, want %q", got, want)
	}
	if got, want := rel.Columns[1].TypeOID, TEXTOID; got != want {
		t.Fatalf("second column type: got %d, want %d", got, want)
	}
}

func TestLoadSQLViewFromReturnsTableFunction(t *testing.T) {
	c, err := LoadSQL(`
CREATE FUNCTION public.table_pair()
RETURNS TABLE(id integer, name text)
LANGUAGE sql
AS $$
  SELECT 1::integer, 'alice'::text;
$$;

CREATE VIEW public.table_pair_view AS
SELECT *
FROM public.table_pair();
`)
	if err != nil {
		t.Fatalf("LoadSQL() error: %v", err)
	}

	rel := c.GetRelation("public", "table_pair_view")
	if rel == nil {
		t.Fatal("view table_pair_view not found")
	}
	if got, want := len(rel.Columns), 2; got != want {
		t.Fatalf("view columns: got %d, want %d", got, want)
	}
	if got, want := rel.Columns[0].Name, "id"; got != want {
		t.Fatalf("first column name: got %q, want %q", got, want)
	}
	if got, want := rel.Columns[0].TypeOID, INT4OID; got != want {
		t.Fatalf("first column type: got %d, want %d", got, want)
	}
	if got, want := rel.Columns[1].Name, "name"; got != want {
		t.Fatalf("second column name: got %q, want %q", got, want)
	}
	if got, want := rel.Columns[1].TypeOID, TEXTOID; got != want {
		t.Fatalf("second column type: got %d, want %d", got, want)
	}
}

func TestLoadSQLViewFromBuiltinRecordReturningFunctionWithOutParams(t *testing.T) {
	c, err := LoadSQL(`
CREATE VIEW public.pg_keywords_view AS
SELECT *
FROM pg_get_keywords();
`)
	if err != nil {
		t.Fatalf("LoadSQL() error: %v", err)
	}

	rel := c.GetRelation("public", "pg_keywords_view")
	if rel == nil {
		t.Fatal("view pg_keywords_view not found")
	}
	wantNames := []string{"word", "catcode", "barelabel", "catdesc", "baredesc"}
	wantTypes := []uint32{TEXTOID, CHAROID, BOOLOID, TEXTOID, TEXTOID}
	if got, want := len(rel.Columns), len(wantNames); got != want {
		t.Fatalf("view columns: got %d, want %d", got, want)
	}
	for i := range wantNames {
		if got, want := rel.Columns[i].Name, wantNames[i]; got != want {
			t.Fatalf("column %d name: got %q, want %q", i+1, got, want)
		}
		if got, want := rel.Columns[i].TypeOID, wantTypes[i]; got != want {
			t.Fatalf("column %d type: got %d, want %d", i+1, got, want)
		}
	}
}

func TestLoadSQLViewFromRecordReturningFunctionWithColumnDefinitionList(t *testing.T) {
	c, err := LoadSQL(`
CREATE FUNCTION public.record_pair_untyped()
RETURNS record
LANGUAGE sql
AS $$
  SELECT 1::integer, 'alice'::text;
$$;

CREATE VIEW public.record_pair_untyped_view AS
SELECT *
FROM public.record_pair_untyped() AS record_pair(id integer, name text);
`)
	if err != nil {
		t.Fatalf("LoadSQL() error: %v", err)
	}

	rel := c.GetRelation("public", "record_pair_untyped_view")
	if rel == nil {
		t.Fatal("view record_pair_untyped_view not found")
	}
	if got, want := len(rel.Columns), 2; got != want {
		t.Fatalf("view columns: got %d, want %d", got, want)
	}
	if got, want := rel.Columns[0].Name, "id"; got != want {
		t.Fatalf("first column name: got %q, want %q", got, want)
	}
	if got, want := rel.Columns[0].TypeOID, INT4OID; got != want {
		t.Fatalf("first column type: got %d, want %d", got, want)
	}
	if got, want := rel.Columns[1].Name, "name"; got != want {
		t.Fatalf("second column name: got %q, want %q", got, want)
	}
	if got, want := rel.Columns[1].TypeOID, TEXTOID; got != want {
		t.Fatalf("second column type: got %d, want %d", got, want)
	}
}

func TestParseProcArrayElementsStrict(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want []string
		ok   bool
	}{
		{
			name: "simple",
			raw:  "{i,o,b}",
			want: []string{"i", "o", "b"},
			ok:   true,
		},
		{
			name: "quoted comma",
			raw:  `{"a,b","c\"d","e\\f"}`,
			want: []string{"a,b", `c"d`, `e\f`},
			ok:   true,
		},
		{
			name: "empty element",
			raw:  "{input,,output}",
			want: []string{"input", "", "output"},
			ok:   true,
		},
		{
			name: "missing braces",
			raw:  "i,o,b",
			ok:   false,
		},
		{
			name: "unterminated quote",
			raw:  `{"a,b}`,
			ok:   false,
		},
		{
			name: "dangling escape",
			raw:  `{"a\}`,
			ok:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseProcArrayElementsStrict(tt.raw)
			if ok != tt.ok {
				t.Fatalf("ok: got %v, want %v", ok, tt.ok)
			}
			if !ok {
				return
			}
			if len(got) != len(tt.want) {
				t.Fatalf("elements: got %v, want %v", got, tt.want)
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Fatalf("element %d: got %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestCreateFunctionWithoutReturnsAndWithoutOutParamsIsRejected(t *testing.T) {
	c := New()

	err := execSQLErr(c, `
CREATE FUNCTION f_missing_returns(input_value integer)
LANGUAGE sql
AS 'SELECT input_value';
`)
	assertCode(t, err, CodeInvalidFunctionDefinition)
}
