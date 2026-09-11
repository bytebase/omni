package catalog

import (
	"slices"
	"testing"

	"github.com/bytebase/omni/metadata"
	nodes "github.com/bytebase/omni/pg/ast"
)

func TestParseTypeName(t *testing.T) {
	for _, typ := range []string{
		"integer",
		"numeric(10,2)",
		"timestamp(3) with time zone",
		"character varying(255)",
		"text[]",
		"public.task_status",
		`"My Schema"."My Type"`,
	} {
		if tn, err := parseTypeName(typ); err != nil || tn == nil {
			t.Errorf("parseTypeName(%q) = %v, %v", typ, tn, err)
		}
	}
	for _, typ := range []string{"", "integer; DROP TABLE t"} {
		if _, err := parseTypeName(typ); err == nil {
			t.Errorf("parseTypeName(%q) succeeded, want an error", typ)
		}
	}
}

func TestUserTypeRef(t *testing.T) {
	tests := []struct {
		typ          string
		schema, name string
	}{
		{typ: "public.task_status", schema: "public", name: "task_status"},
		{typ: "public.task_status[]", schema: "public", name: "task_status"},
		{typ: "public.money_amount(10,2)", schema: "public", name: "money_amount"},
		{typ: `"My Schema"."My Type"`, schema: "My Schema", name: "My Type"},
		{typ: `"a""b".c`, schema: `a"b`, name: "c"},
		// Built-in and array types are unqualified; sync qualifies every user type.
		{typ: "integer"},
		{typ: "character varying(255)"},
		{typ: "timestamp(3) with time zone"},
		{typ: "_int4"},
		{typ: ""},
	}
	for _, tt := range tests {
		schema, name, ok := userTypeRef(tt.typ)
		if ok != (tt.schema != "") || schema != tt.schema || name != tt.name {
			t.Errorf("userTypeRef(%q) = %q, %q, %v; want %q, %q", tt.typ, schema, name, ok, tt.schema, tt.name)
		}
	}
}

func TestSplitQualifiedIdent(t *testing.T) {
	tests := []struct {
		in           string
		schema, name string
		ok           bool
	}{
		{in: "public.foo", schema: "public", name: "foo", ok: true},
		{in: `"weird name"."another"`, schema: "weird name", name: "another", ok: true},
		{in: `"has.dot".x`, schema: "has.dot", name: "x", ok: true},
		{in: "foo"},
		{in: "a.b.c"},
		{in: `"unterminated.x`},
		{in: `"a"b.c`},
		{in: ".x"},
	}
	for _, tt := range tests {
		schema, name, ok := splitQualifiedIdent(tt.in)
		if schema != tt.schema || name != tt.name || ok != tt.ok {
			t.Errorf("splitQualifiedIdent(%q) = %q, %q, %v; want %q, %q, %v", tt.in, schema, name, ok, tt.schema, tt.name, tt.ok)
		}
	}
}

func TestSignatureArgTypes(t *testing.T) {
	tests := []struct {
		signature string
		want      []string
		wantErr   bool
	}{
		{signature: "fn()"},
		{signature: "fn"},
		{signature: "fn(integer)", want: []string{"integer"}},
		{signature: "fn(  int4 , text )", want: []string{"int4", "text"}},
		{signature: "fn(numeric(10,2), public.status[])", want: []string{"numeric(10,2)", "public.status[]"}},
		{signature: "fn(integer", wantErr: true},
	}
	for _, tt := range tests {
		got, err := signatureArgTypes(tt.signature)
		if (err != nil) != tt.wantErr || !slices.Equal(got, tt.want) {
			t.Errorf("signatureArgTypes(%q) = %q, %v; want %q", tt.signature, got, err, tt.want)
		}
	}
}

func TestNextvalSequence(t *testing.T) {
	tests := []struct {
		def          string
		schema, name string
	}{
		{def: "nextval('public.users_id_seq'::regclass)", schema: "public", name: "users_id_seq"},
		{def: "nextval('users_id_seq'::regclass)", schema: "app", name: "users_id_seq"},
		{def: `nextval('"My Schema"."My Seq"'::regclass)`, schema: "My Schema", name: "My Seq"},
		{def: `nextval('"My Seq"'::regclass)`, schema: "app", name: "My Seq"},
		{def: "now()"},
	}
	for _, tt := range tests {
		schema, name, ok := nextvalSequence(tt.def, "app")
		if ok != (tt.name != "") || schema != tt.schema || name != tt.name {
			t.Errorf("nextvalSequence(%q) = %q, %q, %v; want %q, %q", tt.def, schema, name, ok, tt.schema, tt.name)
		}
	}
}

func TestCollateClause(t *testing.T) {
	names := func(c *nodes.CollateClause) []string {
		var out []string
		for _, item := range c.Collname.Items {
			out = append(out, item.(*nodes.String).Str)
		}
		return out
	}
	if collateClause("") != nil {
		t.Error("an empty collation must produce no clause")
	}
	for in, want := range map[string][]string{
		`"C"`:             {"C"},
		"und":             {"und"},
		`"en_US"`:         {"en_US"},
		"locale.en_us":    {"locale", "en_us"},
		`"My Schema".x`:   {"My Schema", "x"},
		`"say ""hi"""`:    {`say "hi"`},
		`"a.b"`:           {"a.b"},
		`public."C.utf8"`: {"public", "C.utf8"},
	} {
		if got := names(collateClause(in)); !slices.Equal(got, want) {
			t.Errorf("collateClause(%q) names %q, want %q", in, got, want)
		}
	}
}

func TestMetadataIndexStmt(t *testing.T) {
	stmt, err := metadataIndexStmt("public", "t", &metadata.IndexMetadata{
		Name:        "t_idx",
		Type:        "btree",
		Expressions: []string{"id", `"Name"`, "lower(email)", "(payload ->> 'k'::text)"},
		Descending:  []bool{false, true},
	})
	if err != nil {
		t.Fatalf("metadataIndexStmt: %v", err)
	}
	elems := stmt.IndexParams.Items
	if len(elems) != 4 {
		t.Fatalf("got %d index keys, want 4", len(elems))
	}
	for i, want := range []string{"id", "Name", "", ""} {
		elem := elems[i].(*nodes.IndexElem)
		if elem.Name != want || (want == "") != (elem.Expr != nil) {
			t.Errorf("key %d: name %q, expr %v; want name %q", i, elem.Name, elem.Expr, want)
		}
	}
	if elems[1].(*nodes.IndexElem).Ordering != nodes.SORTBY_DESC {
		t.Error("second key must be descending")
	}
	if _, err := metadataIndexStmt("public", "t", &metadata.IndexMetadata{Expressions: []string{"id"}}); err == nil {
		t.Error("an index without a name must fail")
	}
}

func TestMetadataFunctionStmtFallsBackToSignature(t *testing.T) {
	stmt, err := metadataFunctionStmt("public", &metadata.FunctionMetadata{
		Name:       "fn",
		Signature:  "fn(integer, text)",
		Definition: "CREATE FUNCTION with a body omni cannot parse",
	})
	if err != nil {
		t.Fatalf("metadataFunctionStmt: %v", err)
	}
	if stmt.Parameters == nil || len(stmt.Parameters.Items) != 2 {
		t.Fatalf("want the signature's two parameters, got %+v", stmt.Parameters)
	}
	c := New()
	if err := c.CreateFunctionStmt(stmt); err != nil {
		t.Fatalf("CreateFunctionStmt: %v", err)
	}
}

func TestStandInSelectQuotesColumnNames(t *testing.T) {
	c := New()
	stmt, err := standInViewStmt("public", "v", []string{"id", "Mixed Case", `quote"d`, "select"})
	if err != nil {
		t.Fatalf("standInViewStmt: %v", err)
	}
	if err := c.DefineView(stmt); err != nil {
		t.Fatalf("DefineView: %v", err)
	}
	var got []string
	for _, col := range c.GetRelation("public", "v").Columns {
		got = append(got, col.Name)
	}
	if want := []string{"id", "Mixed Case", `quote"d`, "select"}; !slices.Equal(got, want) {
		t.Errorf("stand-in view columns = %q, want %q", got, want)
	}
}
