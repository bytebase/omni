package catalog

import (
	"strings"
	"testing"
)

// The normalize-core module is the diff-time canonicalization layer. Its defining
// property is phantom-diff elimination: for a given target engine version,
//
//	Canonical(userTargetForm) == Canonical(syncedStoredForm)
//
// where both forms describe the same logical schema. The synced stored form is what
// MySQL writes back via SHOW CREATE TABLE (captured in normalization.md as oracle
// readbacks); the user target form is what the user declares. These tests assert the
// canonicalizer collapses both onto the same value, version by version.

func col(name, columnType string) *Column {
	return &Column{Name: name, DataType: baseTypeOf(columnType), ColumnType: columnType, Nullable: true}
}

// baseTypeOf extracts the lowercased base type name from a column-type string for
// test fixtures (the real loader sets DataType directly; this mirrors it).
func baseTypeOf(columnType string) string {
	for i := 0; i < len(columnType); i++ {
		c := columnType[i]
		if c == '(' || c == ' ' {
			return columnType[:i]
		}
	}
	return columnType
}

// ---------------------------------------------------------------------------
// Highest-risk rule: ResolveColumnCharsetCollation — resolved-pair comparison.
// (flag column-charset-collation-echo-asymmetry; entries column-charset-echo-57/-80,
//  column-charset-only-collation-resolution, utf8mb4-default-collation)
// ---------------------------------------------------------------------------

func TestResolveColumnCharsetCollation_57_StripsWhenEqualsTable(t *testing.T) {
	// 5.7 readback: a column whose charset == table charset stores NO charset/collation
	// clause (all four columns read back as bare `varchar(10)`). The user's target may
	// write the clause explicitly. Both must resolve to the same effective pair.
	tbl := &Table{Name: "c1", Charset: "utf8mb4", Collation: "utf8mb4_general_ci"}

	// User target: explicit charset+collation equal to the table default.
	userCol := &Column{Name: "a", DataType: "varchar", ColumnType: "varchar(10)",
		Charset: "utf8mb4", CharsetExplicit: true, Collation: "utf8mb4_general_ci", CollationExplicit: true, Nullable: true}
	// Synced stored form: 5.7 stripped both clauses, leaving them empty (inherited).
	storedCol := &Column{Name: "a", DataType: "varchar", ColumnType: "varchar(10)",
		Charset: "", Collation: "", Nullable: true}

	n := NormalizerFor(MySQL57)
	uc, ucol := n.ResolveColumnCharsetCollation(tbl, userCol)
	sc, scol := n.ResolveColumnCharsetCollation(tbl, storedCol)
	if uc != sc || ucol != scol {
		t.Fatalf("5.7 equal-to-table charset must resolve identically: user=(%q,%q) stored=(%q,%q)", uc, ucol, sc, scol)
	}
	if uc != "utf8mb4" || ucol != "utf8mb4_general_ci" {
		t.Fatalf("expected resolved pair (utf8mb4, utf8mb4_general_ci), got (%q,%q)", uc, ucol)
	}
}

func TestResolveColumnCharsetCollation_57_KeepsGenuineOverride(t *testing.T) {
	// 5.7 readback: a column whose charset differs from the table keeps it
	// (`CHARACTER SET utf8` under a utf8mb4 table).
	tbl := &Table{Name: "cu", Charset: "utf8mb4", Collation: "utf8mb4_general_ci"}
	userCol := &Column{Name: "a", DataType: "varchar", ColumnType: "varchar(10)",
		Charset: "utf8", CharsetExplicit: true, Nullable: true} // user wrote CHARACTER SET utf8, no collate
	storedCol := &Column{Name: "a", DataType: "varchar", ColumnType: "varchar(10)",
		Charset: "utf8", CharsetExplicit: true, Nullable: true} // 5.7 keeps `CHARACTER SET utf8`

	n := NormalizerFor(MySQL57)
	uc, ucol := n.ResolveColumnCharsetCollation(tbl, userCol)
	sc, scol := n.ResolveColumnCharsetCollation(tbl, storedCol)
	if uc != sc || ucol != scol {
		t.Fatalf("override must resolve identically: user=(%q,%q) stored=(%q,%q)", uc, ucol, sc, scol)
	}
	// utf8 ≡ utf8mb3; collation resolves to the charset default, not the table's.
	if uc != "utf8mb3" {
		t.Fatalf("expected resolved charset utf8mb3 (utf8 alias), got %q", uc)
	}
	if ucol != "utf8_general_ci" && ucol != "utf8mb3_general_ci" {
		t.Fatalf("expected utf8 charset-default collation, got %q", ucol)
	}
}

func TestResolveColumnCharsetCollation_80_KeepsRedundantPair(t *testing.T) {
	// 8.0 readback: when the user specifies charset OR collation, 8.0 echoes BOTH as a
	// pair even if it equals the table default. A bare column collapses to varchar(10).
	tbl := &Table{Name: "c1", Charset: "utf8mb4", Collation: "utf8mb4_0900_ai_ci"}

	// `c` -> CHARACTER SET utf8mb4 (no collate). Stored as the resolved pair.
	userCol := &Column{Name: "c", DataType: "varchar", ColumnType: "varchar(10)",
		Charset: "utf8mb4", CharsetExplicit: true, Nullable: true}
	storedCol := &Column{Name: "c", DataType: "varchar", ColumnType: "varchar(10)",
		Charset: "utf8mb4", CharsetExplicit: true, Collation: "utf8mb4_0900_ai_ci", CollationExplicit: true, Nullable: true}

	n := NormalizerFor(MySQL80)
	uc, ucol := n.ResolveColumnCharsetCollation(tbl, userCol)
	sc, scol := n.ResolveColumnCharsetCollation(tbl, storedCol)
	if uc != sc || ucol != scol {
		t.Fatalf("8.0 pair must resolve identically: user=(%q,%q) stored=(%q,%q)", uc, ucol, sc, scol)
	}
	if ucol != "utf8mb4_0900_ai_ci" {
		t.Fatalf("8.0 must resolve missing collation to charset default, got %q", ucol)
	}
}

func TestResolveColumnCharsetCollation_80_BareInheritsTable(t *testing.T) {
	// 8.0 readback: a fully-bare column (`d`) collapses to varchar(10) — it inherits
	// the table's charset+collation. The resolved pair equals the table's.
	tbl := &Table{Name: "c1", Charset: "utf8mb4", Collation: "utf8mb4_0900_ai_ci"}
	bare := &Column{Name: "d", DataType: "varchar", ColumnType: "varchar(10)", Nullable: true}

	n := NormalizerFor(MySQL80)
	cs, coll := n.ResolveColumnCharsetCollation(tbl, bare)
	if cs != "utf8mb4" || coll != "utf8mb4_0900_ai_ci" {
		t.Fatalf("bare column must inherit table pair, got (%q,%q)", cs, coll)
	}
}

func TestResolveColumnCharsetCollation_80_CharsetOnlyResolvesFromCharsetDefault(t *testing.T) {
	// entry column-charset-only-collation-resolution: under a table collated
	// utf8mb4_unicode_ci, a column with only CHARACTER SET utf8mb4 resolves its
	// collation from the CHARSET default (utf8mb4_0900_ai_ci on 8.0), NOT the table COLLATE.
	tbl := &Table{Name: "t_d1", Charset: "utf8mb4", Collation: "utf8mb4_unicode_ci"}
	userCol := &Column{Name: "a", DataType: "varchar", ColumnType: "varchar(10)",
		Charset: "utf8mb4", CharsetExplicit: true, Nullable: true}

	n := NormalizerFor(MySQL80)
	_, coll := n.ResolveColumnCharsetCollation(tbl, userCol)
	if coll != "utf8mb4_0900_ai_ci" {
		t.Fatalf("charset-only column must resolve to charset default utf8mb4_0900_ai_ci, not table COLLATE; got %q", coll)
	}
}

// ---------------------------------------------------------------------------
// CanonicalColumnType — integer display width, BOOLEAN, decimal, year, unsigned,
// zerofill. (entries int-display-width, tinyint1-boolean, int-unsigned-width,
//  int-zerofill, decimal-precision-scale, float-double-aliasing, year-width,
//  bit/char/binary length defaults)
// ---------------------------------------------------------------------------

// canonTypeEq asserts that every input spelling maps to the same canonical type for
// the given version. The loader stores the 8.0 form (e.g. "int"); a raw 5.7 readback
// string would be "int(11)". Both, plus any user spelling, must collapse to one value.
func canonTypeEq(t *testing.T, v Version, want string, inputs ...string) {
	t.Helper()
	n := NormalizerFor(v)
	for _, in := range inputs {
		got := n.CanonicalColumnType(col("x", in))
		if got != want {
			t.Errorf("[%v] CanonicalColumnType(%q) = %q, want %q", v, in, got, want)
		}
	}
}

func TestCanonicalColumnType_IntWidth_80DropsWidth(t *testing.T) {
	canonTypeEq(t, MySQL80, "int", "int", "int(11)", "INT(11)", "integer", "int(5)")
	canonTypeEq(t, MySQL80, "bigint", "bigint", "bigint(20)")
	canonTypeEq(t, MySQL80, "smallint", "smallint", "smallint(6)")
	canonTypeEq(t, MySQL80, "mediumint", "mediumint", "mediumint(9)")
	canonTypeEq(t, MySQL80, "tinyint", "tinyint", "tinyint", "tinyint(4)")
}

func TestCanonicalColumnType_IntWidth_57InjectsDefault(t *testing.T) {
	// 5.7 injects the default width when omitted. Loader-form "int" and readback
	// "int(11)" must both canonicalize to "int(11)".
	canonTypeEq(t, MySQL57, "int(11)", "int", "int(11)", "integer")
	canonTypeEq(t, MySQL57, "bigint(20)", "bigint", "bigint(20)")
	canonTypeEq(t, MySQL57, "smallint(6)", "smallint", "smallint(6)")
	canonTypeEq(t, MySQL57, "mediumint(9)", "mediumint", "mediumint(9)")
	canonTypeEq(t, MySQL57, "tinyint(4)", "tinyint", "tinyint(4)")
}

func TestCanonicalColumnType_Tinyint1Boolean(t *testing.T) {
	// tinyint(1) survives on BOTH versions (boolean marker); BOOLEAN/BOOL → tinyint(1).
	canonTypeEq(t, MySQL80, "tinyint(1)", "tinyint(1)", "boolean", "bool")
	canonTypeEq(t, MySQL57, "tinyint(1)", "tinyint(1)", "boolean", "bool")
	// bare tinyint differs: 8.0 drops to tinyint, 5.7 injects tinyint(4).
	canonTypeEq(t, MySQL80, "tinyint", "tinyint(4)")
	canonTypeEq(t, MySQL57, "tinyint(4)", "tinyint")
}

func TestCanonicalColumnType_UnsignedWidth(t *testing.T) {
	canonTypeEq(t, MySQL80, "int unsigned", "int unsigned", "int(11) unsigned", "int(10) unsigned")
	// 5.7 bare INT UNSIGNED → int(10) unsigned (unsigned default width is 10, not 11).
	canonTypeEq(t, MySQL57, "int(10) unsigned", "int unsigned", "int(10) unsigned")
	canonTypeEq(t, MySQL57, "bigint(20) unsigned", "bigint unsigned", "bigint(20) unsigned")
	canonTypeEq(t, MySQL57, "tinyint(3) unsigned", "tinyint unsigned")
}

func TestCanonicalColumnType_Zerofill_WidthKeptBothVersions(t *testing.T) {
	// ZEROFILL retains display width on BOTH versions and implies unsigned.
	want := "int(5) unsigned zerofill"
	canonTypeEq(t, MySQL80, want, "int(5) zerofill", "int(5) unsigned zerofill")
	canonTypeEq(t, MySQL57, want, "int(5) zerofill", "int(5) unsigned zerofill")
	// bare INT ZEROFILL → width 10 injected, both versions.
	canonTypeEq(t, MySQL80, "int(10) unsigned zerofill", "int zerofill")
	canonTypeEq(t, MySQL57, "int(10) unsigned zerofill", "int zerofill")
}

func TestCanonicalColumnType_DecimalPrecisionScale(t *testing.T) {
	// NUMERIC/DEC/FIXED → decimal; bare → decimal(10,0); DECIMAL(10) → decimal(10,0).
	for _, v := range []Version{MySQL57, MySQL80} {
		canonTypeEq(t, v, "decimal(10,0)", "decimal", "decimal(10)", "numeric", "dec", "fixed", "decimal(10,0)")
		canonTypeEq(t, v, "decimal(8,3)", "decimal(8,3)", "numeric(8,3)")
	}
}

func TestCanonicalColumnType_FloatDoubleAliasing(t *testing.T) {
	for _, v := range []Version{MySQL57, MySQL80} {
		canonTypeEq(t, v, "float", "float", "float(5)") // single-arg precision dropped
		canonTypeEq(t, v, "double", "double", "real")
		canonTypeEq(t, v, "float(10,2)", "float(10,2)")
		canonTypeEq(t, v, "double(15,4)", "double(15,4)")
	}
}

func TestCanonicalColumnType_YearWidth(t *testing.T) {
	canonTypeEq(t, MySQL80, "year", "year", "year(4)")
	canonTypeEq(t, MySQL57, "year(4)", "year", "year(4)")
}

func TestCanonicalColumnType_BitCharBinaryLengthDefault(t *testing.T) {
	for _, v := range []Version{MySQL57, MySQL80} {
		canonTypeEq(t, v, "bit(1)", "bit", "bit(1)")
		canonTypeEq(t, v, "char(1)", "char", "char(1)")
		canonTypeEq(t, v, "binary(1)", "binary", "binary(1)")
		canonTypeEq(t, v, "varchar(10)", "varchar(10)")
		canonTypeEq(t, v, "varbinary(32)", "varbinary(32)")
	}
}

func TestCanonicalColumnType_TextBlobNoWidth(t *testing.T) {
	for _, v := range []Version{MySQL57, MySQL80} {
		canonTypeEq(t, v, "text", "text")
		canonTypeEq(t, v, "blob", "blob")
		canonTypeEq(t, v, "longtext", "longtext")
	}
}

// ---------------------------------------------------------------------------
// CanonicalDefault — value-based comparison. (entries default-literal-quoting,
//  decimal-default-padding, boolean-default, functional-default,
//  nullable-default-null, timestamp-datetime-default-expr; flag
//  numeric-default-quote-version)
// ---------------------------------------------------------------------------

func defCol(columnType string, def *string, kind ColumnDefaultKind) *Column {
	c := col("x", columnType)
	c.Default = def
	c.DefaultKind = kind
	return c
}
func sp(s string) *string { return &s }

func defaultKeyEq(t *testing.T, v Version, a, b *Column) {
	t.Helper()
	n := NormalizerFor(v)
	ka, kb := n.CanonicalDefault(a), n.CanonicalDefault(b)
	if ka != kb {
		t.Fatalf("[%v] defaults must canonicalize equal: %q vs %q", v, ka, kb)
	}
}

func TestCanonicalDefault_NumericValueBased(t *testing.T) {
	// flag numeric-default-quote-version: compare numeric defaults by VALUE, not string.
	// Loader stores "0" for DEFAULT 0; readback is "'0'". Both must be equal.
	for _, v := range []Version{MySQL57, MySQL80} {
		defaultKeyEq(t, v,
			defCol("int", sp("0"), ColumnDefaultConstant),
			defCol("int", sp("'0'"), ColumnDefaultConstant))
		defaultKeyEq(t, v,
			defCol("int", sp("5"), ColumnDefaultConstant),
			defCol("int", sp("'5'"), ColumnDefaultConstant))
	}
}

func TestCanonicalDefault_DecimalScalePadding(t *testing.T) {
	// A numeric default on a scaled DECIMAL is padded to the column scale: 0 → '0.00'.
	for _, v := range []Version{MySQL57, MySQL80} {
		defaultKeyEq(t, v,
			defCol("decimal(10,2)", sp("0"), ColumnDefaultConstant),
			defCol("decimal(10,2)", sp("'0.00'"), ColumnDefaultConstant))
		defaultKeyEq(t, v,
			defCol("decimal(10,2)", sp("0.00"), ColumnDefaultConstant),
			defCol("decimal(10,2)", sp("'0.00'"), ColumnDefaultConstant))
	}
}

func TestCanonicalDefault_StringQuoting(t *testing.T) {
	for _, v := range []Version{MySQL57, MySQL80} {
		// Double-quoted vs single-quoted string default are equal.
		defaultKeyEq(t, v,
			defCol("varchar(10)", sp(`"double"`), ColumnDefaultConstant),
			defCol("varchar(10)", sp(`'double'`), ColumnDefaultConstant))
		// Empty-string default.
		defaultKeyEq(t, v,
			defCol("varchar(10)", sp(`''`), ColumnDefaultConstant),
			defCol("varchar(10)", sp(`''`), ColumnDefaultConstant))
	}
}

func TestCanonicalDefault_NullHandling(t *testing.T) {
	// A nullable column with explicit NULL / DEFAULT NULL / no default are equivalent.
	for _, v := range []Version{MySQL57, MySQL80} {
		nullDef := defCol("int", sp("NULL"), ColumnDefaultConstant)
		noDef := col("x", "int") // Default == nil
		defaultKeyEq(t, v, nullDef, noDef)
	}
}

func TestCanonicalDefault_BooleanLiteral(t *testing.T) {
	// BOOLEAN DEFAULT TRUE/FALSE → '1'/'0'. Loader already folds TRUE→"1".
	for _, v := range []Version{MySQL57, MySQL80} {
		defaultKeyEq(t, v,
			defCol("tinyint(1)", sp("1"), ColumnDefaultConstant),
			defCol("tinyint(1)", sp("'1'"), ColumnDefaultConstant))
	}
}

func TestCanonicalDefault_CurrentTimestampSynonyms(t *testing.T) {
	// NOW()/LOCALTIME/LOCALTIMESTAMP → CURRENT_TIMESTAMP, unquoted, case-insensitive.
	for _, v := range []Version{MySQL57, MySQL80} {
		defaultKeyEq(t, v,
			defCol("datetime", sp("CURRENT_TIMESTAMP"), ColumnDefaultCurrentTimestamp),
			defCol("datetime", sp("now()"), ColumnDefaultCurrentTimestamp))
		defaultKeyEq(t, v,
			defCol("timestamp(3)", sp("CURRENT_TIMESTAMP(3)"), ColumnDefaultCurrentTimestamp),
			defCol("timestamp(3)", sp("current_timestamp(3)"), ColumnDefaultCurrentTimestamp))
	}
}

func TestCanonicalDefault_DistinctValuesDiffer(t *testing.T) {
	// Guard against over-collapsing: different defaults must NOT compare equal.
	n := NormalizerFor(MySQL80)
	if n.CanonicalDefault(defCol("int", sp("0"), ColumnDefaultConstant)) ==
		n.CanonicalDefault(defCol("int", sp("1"), ColumnDefaultConstant)) {
		t.Fatal("DEFAULT 0 and DEFAULT 1 must not be equal")
	}
	if n.CanonicalDefault(defCol("varchar(10)", sp("'a'"), ColumnDefaultConstant)) ==
		n.CanonicalDefault(defCol("varchar(10)", sp("'b'"), ColumnDefaultConstant)) {
		t.Fatal("DEFAULT 'a' and DEFAULT 'b' must not be equal")
	}
	// A NULL default and a non-null default must differ.
	if n.CanonicalDefault(col("x", "int")) ==
		n.CanonicalDefault(defCol("int", sp("0"), ColumnDefaultConstant)) {
		t.Fatal("no default and DEFAULT 0 must not be equal")
	}
}

// ---------------------------------------------------------------------------
// CanonicalGeneratedExpr — expression canonicalization. (entries
//  generated-expr-normalization, generated-expr-string-introducer, functional-default,
//  check-constraint)
// ---------------------------------------------------------------------------

func genExprEq(t *testing.T, v Version, inputs ...string) {
	t.Helper()
	n := NormalizerFor(v)
	want := n.CanonicalGeneratedExpr(inputs[0])
	for _, in := range inputs[1:] {
		if got := n.CanonicalGeneratedExpr(in); got != want {
			t.Errorf("[%v] CanonicalGeneratedExpr(%q)=%q != CanonicalGeneratedExpr(%q)=%q", v, in, got, inputs[0], want)
		}
	}
}

func TestCanonicalGeneratedExpr_WhitespaceAndParens(t *testing.T) {
	// (a+1), ( a + 1 ), `a` + 1, ((`a` + 1)) all collapse to one canonical form.
	for _, v := range []Version{MySQL57, MySQL80} {
		genExprEq(t, v,
			"`a` + 1",
			"(`a` + 1)",
			"((`a` + 1))",
			"a+1",
			"( a   +   1 )",
		)
	}
}

func TestCanonicalGeneratedExpr_LatinIntroducerStripped(t *testing.T) {
	// 8.0 readback injects _latin1 before string literals; 5.7 and user input do not.
	// Both must compare equal (the introducer is connection-dependent noise).
	for _, v := range []Version{MySQL57, MySQL80} {
		genExprEq(t, v,
			"concat(`a`,'x')",
			"concat(`a`,_latin1'x')",
			"concat(`a`, 'x')",
		)
	}
}

func TestCanonicalGeneratedExpr_FunctionLowercasing(t *testing.T) {
	for _, v := range []Version{MySQL57, MySQL80} {
		genExprEq(t, v,
			"concat(`a`,`b`)",
			"CONCAT(`a`,`b`)",
			"Concat(`a`, `b`)",
		)
	}
}

func TestCanonicalGeneratedExpr_DistinctExprsDiffer(t *testing.T) {
	n := NormalizerFor(MySQL80)
	if n.CanonicalGeneratedExpr("`a` + 1") == n.CanonicalGeneratedExpr("`a` + 2") {
		t.Fatal("a+1 and a+2 must not compare equal")
	}
	if n.CanonicalGeneratedExpr("`a` * 2") == n.CanonicalGeneratedExpr("`a` + 2") {
		t.Fatal("a*2 and a+2 must not compare equal")
	}
}

// ---------------------------------------------------------------------------
// Nullability + PK. (entries nullable-default-null, pk-implies-not-null)
// ---------------------------------------------------------------------------

func TestCanonicalNullability_PKForcesNotNull(t *testing.T) {
	// A PK-member column is forced NOT NULL even if declared nullable.
	tbl := &Table{Name: "t_pk"}
	idCol := &Column{Name: "id", DataType: "int", ColumnType: "int", Nullable: true}
	tbl.Columns = []*Column{idCol}
	tbl.Indexes = []*Index{{Primary: true, Columns: []*IndexColumn{{Name: "id"}}}}

	n := NormalizerFor(MySQL80)
	if n.CanonicalNotNull(tbl, idCol) != true {
		t.Fatal("PK column must canonicalize to NOT NULL")
	}
}

func TestCanonicalNullability_NonPKNullablePreserved(t *testing.T) {
	tbl := &Table{Name: "t"}
	c := &Column{Name: "a", DataType: "int", ColumnType: "int", Nullable: true}
	tbl.Columns = []*Column{c}
	n := NormalizerFor(MySQL80)
	if n.CanonicalNotNull(tbl, c) != false {
		t.Fatal("non-PK nullable column must stay nullable")
	}
}

// ---------------------------------------------------------------------------
// TIMESTAMP magic — governed by explicit_defaults_for_timestamp (a session var, not a
// version trait). (entry timestamp-magic-first-57; flag explicit-defaults-for-timestamp)
// ---------------------------------------------------------------------------

func tsTable(cols ...*Column) *Table {
	t := &Table{Name: "t_ts"}
	t.Columns = cols
	return t
}
func tsCol(name string) *Column {
	return &Column{Name: name, DataType: "timestamp", ColumnType: "timestamp", Nullable: true}
}

func TestCanonicalTimestamp_EDFTOff_FirstBareGetsMagic(t *testing.T) {
	// explicit_defaults_for_timestamp=0: the FIRST bare TIMESTAMP gets
	// NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP.
	n := &Normalizer{Version: MySQL57, ExplicitDefaultsForTimestamp: false}
	a := tsCol("a")
	tbl := tsTable(a)

	if !n.CanonicalNotNull(tbl, a) {
		t.Fatal("EDFT=0: first bare TIMESTAMP must be NOT NULL")
	}
	def, onUpdate := n.CanonicalTimestampDefaults(tbl, a)
	if def != "ts:CURRENT_TIMESTAMP" {
		t.Fatalf("EDFT=0: first bare TIMESTAMP default = %q, want ts:CURRENT_TIMESTAMP", def)
	}
	if onUpdate != "ts:CURRENT_TIMESTAMP" {
		t.Fatalf("EDFT=0: first bare TIMESTAMP on-update = %q, want ts:CURRENT_TIMESTAMP", onUpdate)
	}
}

func TestCanonicalTimestamp_EDFTOff_NullExplicitDistinguishesBareFromNull(t *testing.T) {
	// EDFT=0: a TIMESTAMP with NO explicit nullability clause is forced NOT NULL — even a
	// non-first one — while an explicit `TIMESTAMP NULL` stays nullable. The loader's
	// NullExplicit flag tells them apart (a bare TIMESTAMP and `TIMESTAMP NULL` are
	// otherwise the same catalog state).
	n := &Normalizer{Version: MySQL57, ExplicitDefaultsForTimestamp: false}
	a := tsCol("a")            // first, bare
	bare := tsCol("b")         // non-first, bare → NOT NULL
	explicitNull := tsCol("c") // non-first, explicit NULL → nullable
	explicitNull.NullExplicit = true
	tbl := tsTable(a, bare, explicitNull)

	if !n.CanonicalNotNull(tbl, bare) {
		t.Fatal("EDFT=0: a non-first bare TIMESTAMP must be NOT NULL")
	}
	if n.CanonicalNotNull(tbl, explicitNull) {
		t.Fatal("EDFT=0: an explicit `TIMESTAMP NULL` must stay nullable")
	}
}

func TestCanonicalTimestamp_EDFTOff_FirstWithPrecisionGetsPreciseDefault(t *testing.T) {
	// First TIMESTAMP(3) under EDFT=0 → DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE
	// CURRENT_TIMESTAMP(3), matching the column's fractional precision.
	n := &Normalizer{Version: MySQL57, ExplicitDefaultsForTimestamp: false}
	a := &Column{Name: "a", DataType: "timestamp", ColumnType: "timestamp(3)", Nullable: true}
	tbl := tsTable(a)
	def, onUpdate := n.CanonicalTimestampDefaults(tbl, a)
	if def != "ts:CURRENT_TIMESTAMP(3)" || onUpdate != "ts:CURRENT_TIMESTAMP(3)" {
		t.Fatalf("TIMESTAMP(3) first column must inject CURRENT_TIMESTAMP(3); got def=%q onUpdate=%q", def, onUpdate)
	}
}

func TestCanonicalTimestamp_EDFTOn_BareIsNullable(t *testing.T) {
	// explicit_defaults_for_timestamp=1 (8.0 box default): bare TIMESTAMP → NULL DEFAULT NULL.
	n := &Normalizer{Version: MySQL80, ExplicitDefaultsForTimestamp: true}
	a := tsCol("a")
	tbl := tsTable(a)
	if n.CanonicalNotNull(tbl, a) {
		t.Fatal("EDFT=1: bare TIMESTAMP must be NULLABLE")
	}
}

func TestCanonicalTimestamp_SameVersionDifferentEDFT_Diverges(t *testing.T) {
	// The whole point of flag #3: a 5.7 server with EDFT=1 behaves like 8.0, NOT like
	// the 5.7 box. Nullability must follow the captured variable, not the version.
	a1 := tsCol("a")
	a2 := tsCol("a")
	on := &Normalizer{Version: MySQL57, ExplicitDefaultsForTimestamp: true}   // 5.7 + EDFT ON
	off := &Normalizer{Version: MySQL57, ExplicitDefaultsForTimestamp: false} // 5.7 + EDFT OFF
	if on.CanonicalNotNull(tsTable(a1), a1) == off.CanonicalNotNull(tsTable(a2), a2) {
		t.Fatal("TIMESTAMP nullability must depend on explicit_defaults_for_timestamp, not version")
	}
}

// ---------------------------------------------------------------------------
// Table options — engine default, ignore-in-diff (AUTO_INCREMENT, ROW_FORMAT).
// (entries engine-always-emitted, row-format-default-omitted, auto-increment-counter;
//  flag row-format-default-source)
// ---------------------------------------------------------------------------

func TestCanonicalEngine_DefaultResolved(t *testing.T) {
	n := NormalizerFor(MySQL80)
	// Empty engine resolves to the server default InnoDB; explicit InnoDB matches.
	if n.CanonicalEngine(&Table{Engine: ""}) != n.CanonicalEngine(&Table{Engine: "InnoDB"}) {
		t.Fatal("empty engine must resolve to default InnoDB")
	}
	if n.CanonicalEngine(&Table{Engine: "innodb"}) != n.CanonicalEngine(&Table{Engine: "INNODB"}) {
		t.Fatal("engine comparison must be case-insensitive")
	}
}

func TestIgnoreInDiff_AutoIncrementCounter(t *testing.T) {
	// The table-level AUTO_INCREMENT=N counter is a runtime value, ignored in diff.
	n := NormalizerFor(MySQL80)
	if !n.IgnoreTableAutoIncrement() {
		t.Fatal("table AUTO_INCREMENT counter must be ignored in diff")
	}
}

func TestIgnoreInDiff_RowFormatDefault(t *testing.T) {
	// ROW_FORMAT is ignored when it is empty/the environment default; a user-set value
	// is not ignored. (flag row-format-default-source: exclude unless explicitly set)
	n := NormalizerFor(MySQL80)
	if !n.IgnoreRowFormat(&Table{RowFormat: ""}) {
		t.Fatal("default (empty) ROW_FORMAT must be ignored")
	}
	if !n.IgnoreRowFormat(&Table{RowFormat: "DEFAULT"}) {
		t.Fatal("ROW_FORMAT=DEFAULT must be ignored")
	}
	if n.IgnoreRowFormat(&Table{RowFormat: "DYNAMIC"}) {
		t.Fatal("explicit ROW_FORMAT=DYNAMIC must NOT be ignored")
	}
}

// ---------------------------------------------------------------------------
// Comments — content comparison, decoded. (entry comment-escaping)
// ---------------------------------------------------------------------------

func TestCanonicalComment_ContentIsAlreadyDecoded(t *testing.T) {
	n := NormalizerFor(MySQL80)
	// The loader stores DECODED comment content. The canonicalizer must NOT re-decode:
	// a genuine content of `a''b` (from DDL COMMENT 'a''''b') must stay distinct from a
	// content of `a'b` (from DDL COMMENT 'a''b').
	if n.CanonicalComment("a''b") == n.CanonicalComment("a'b") {
		t.Fatal("distinct decoded comments `a''b` and `a'b` must not collapse")
	}
	if n.CanonicalComment("") != "" {
		t.Fatal("empty comment must yield empty key")
	}
	if n.CanonicalComment("a") == n.CanonicalComment("b") {
		t.Fatal("distinct comments must differ")
	}
}

// ---------------------------------------------------------------------------
// Index column direction — descending is version-flagged. (entry index-desc-asc)
// ---------------------------------------------------------------------------

func TestCanonicalIndexColumn_DescendingVersionFlagged(t *testing.T) {
	descCol := &IndexColumn{Name: "a", Descending: true}
	ascCol := &IndexColumn{Name: "a", Descending: false}

	// 8.0 supports descending: DESC and ASC produce different keys.
	n80 := NormalizerFor(MySQL80)
	if n80.CanonicalIndexColumn(descCol) == n80.CanonicalIndexColumn(ascCol) {
		t.Fatal("8.0: DESC and ASC index columns must differ")
	}
	// 5.7 ignores direction: DESC and ASC produce the SAME key (stored ascending).
	n57 := NormalizerFor(MySQL57)
	if n57.CanonicalIndexColumn(descCol) != n57.CanonicalIndexColumn(ascCol) {
		t.Fatal("5.7: direction must be ignored (no descending indexes)")
	}
}

func TestCanonicalIndexColumn_PrefixLengthPreserved(t *testing.T) {
	for _, v := range []Version{MySQL57, MySQL80} {
		n := NormalizerFor(v)
		a10 := &IndexColumn{Name: "a", Length: 10}
		a20 := &IndexColumn{Name: "a", Length: 20}
		if n.CanonicalIndexColumn(a10) == n.CanonicalIndexColumn(a20) {
			t.Fatalf("[%v] prefix lengths 10 and 20 must differ", v)
		}
	}
}

// ---------------------------------------------------------------------------
// ENUM / SET quoting. (entry enum-set-quoting)
// ---------------------------------------------------------------------------

func TestCanonicalColumnType_EnumSetQuoting(t *testing.T) {
	for _, v := range []Version{MySQL57, MySQL80} {
		// Double-quoted members canonicalize to single-quoted; order preserved.
		canonTypeEq(t, v, "enum('x','y','z')", "enum('x','y','z')", `enum("x","y","z")`)
		canonTypeEq(t, v, "set('a','b','c')", "set('a','b','c')")
	}
	// Different member ORDER must NOT be equal (order is significant for ENUM).
	n := NormalizerFor(MySQL80)
	if n.CanonicalColumnType(col("x", "enum('a','b')")) == n.CanonicalColumnType(col("x", "enum('b','a')")) {
		t.Fatal("ENUM member order is significant; reordered members must differ")
	}
}

// ---------------------------------------------------------------------------
// CanonicalColumn — the aggregate key diff-core's column comparison calls. Ties
// type + resolved charset/collation + default + nullability + generated expr into one
// comparison key per column, so equal logical columns produce equal keys regardless of
// surface form. This is the phantom-diff firewall.
// ---------------------------------------------------------------------------

func TestCanonicalColumn_8080PhantomDiffEliminated(t *testing.T) {
	// 8.0: user writes `INT(11)` with redundant explicit charset; synced readback is
	// `int` bare. The aggregate keys must be equal — no phantom diff.
	tbl := &Table{Name: "t", Charset: "utf8mb4", Collation: "utf8mb4_0900_ai_ci"}
	n := NormalizerFor(MySQL80)

	user := &Column{Name: "a", DataType: "int", ColumnType: "int(11)", Nullable: true}
	stored := &Column{Name: "a", DataType: "int", ColumnType: "int", Nullable: true}
	if n.CanonicalColumn(tbl, user) != n.CanonicalColumn(tbl, stored) {
		t.Fatalf("8.0 int(11) vs int must not phantom-diff:\n user=%s\n stor=%s",
			n.CanonicalColumn(tbl, user), n.CanonicalColumn(tbl, stored))
	}
}

func TestCanonicalColumn_57CharsetEchoPhantomEliminated(t *testing.T) {
	// 5.7: user writes explicit `CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci`
	// equal to the table default; synced readback dropped both. Keys must match.
	tbl := &Table{Name: "t", Charset: "utf8mb4", Collation: "utf8mb4_general_ci"}
	n := NormalizerFor(MySQL57)

	user := &Column{Name: "a", DataType: "varchar", ColumnType: "varchar(10)",
		Charset: "utf8mb4", Collation: "utf8mb4_general_ci", Nullable: true}
	stored := &Column{Name: "a", DataType: "varchar", ColumnType: "varchar(10)", Nullable: true}
	if n.CanonicalColumn(tbl, user) != n.CanonicalColumn(tbl, stored) {
		t.Fatalf("5.7 charset-echo must not phantom-diff:\n user=%s\n stor=%s",
			n.CanonicalColumn(tbl, user), n.CanonicalColumn(tbl, stored))
	}
}

func TestCanonicalColumn_GenuineDifferenceDetected(t *testing.T) {
	// The firewall must still detect a REAL change (int → bigint).
	tbl := &Table{Name: "t", Charset: "utf8mb4", Collation: "utf8mb4_0900_ai_ci"}
	n := NormalizerFor(MySQL80)
	a := &Column{Name: "a", DataType: "int", ColumnType: "int", Nullable: true}
	b := &Column{Name: "a", DataType: "bigint", ColumnType: "bigint", Nullable: true}
	if n.CanonicalColumn(tbl, a) == n.CanonicalColumn(tbl, b) {
		t.Fatal("int vs bigint must produce different keys")
	}
}

// ---------------------------------------------------------------------------
// Idempotence: Canonical(Canonical(x)) == Canonical(x) for every canonicalizer.
// ---------------------------------------------------------------------------

func TestIdempotence_ColumnType(t *testing.T) {
	for _, v := range []Version{MySQL57, MySQL80} {
		n := NormalizerFor(v)
		for _, in := range []string{"int(11)", "int", "bigint", "tinyint(1)", "boolean",
			"decimal(10,2)", "decimal", "float(10,2)", "year(4)", "varchar(10)",
			"int(5) zerofill", "int unsigned", "enum('a','b')", "char", "bit"} {
			once := n.CanonicalColumnType(col("x", in))
			twice := n.CanonicalColumnType(col("x", once))
			if once != twice {
				t.Errorf("[%v] CanonicalColumnType not idempotent for %q: %q -> %q", v, in, once, twice)
			}
		}
	}
}

func TestIdempotence_GeneratedExpr(t *testing.T) {
	n := NormalizerFor(MySQL80)
	for _, in := range []string{"`a` + 1", "((`a` + 1))", "concat(`a`,_latin1'x')", "CONCAT(`a`, `b`)"} {
		once := n.CanonicalGeneratedExpr(in)
		twice := n.CanonicalGeneratedExpr(once)
		if once != twice {
			t.Errorf("CanonicalGeneratedExpr not idempotent for %q: %q -> %q", in, once, twice)
		}
	}
}

func TestIdempotence_CharsetCollation(t *testing.T) {
	tbl := &Table{Charset: "utf8mb4", Collation: "utf8mb4_0900_ai_ci"}
	for _, v := range []Version{MySQL57, MySQL80} {
		n := NormalizerFor(v)
		in := &Column{Name: "a", DataType: "varchar", ColumnType: "varchar(10)", Charset: "utf8", Nullable: true}
		cs, coll := n.ResolveColumnCharsetCollation(tbl, in)
		// Feed the resolved pair back as an explicit column; must be stable.
		in2 := &Column{Name: "a", DataType: "varchar", ColumnType: "varchar(10)", Charset: cs, Collation: coll, Nullable: true}
		cs2, coll2 := n.ResolveColumnCharsetCollation(tbl, in2)
		if cs != cs2 || coll != coll2 {
			t.Errorf("[%v] charset/collation resolution not idempotent: (%q,%q) -> (%q,%q)", v, cs, coll, cs2, coll2)
		}
	}
}

// ---------------------------------------------------------------------------
// Regression tests for review findings (correctness fixes).
// ---------------------------------------------------------------------------

func TestCanonicalDefault_LargeBigintPrecisionPreserved(t *testing.T) {
	// Numeric defaults must compare by full-precision decimal string, never float64,
	// or large BIGINT values beyond 2^53 collide (a missed diff).
	n := NormalizerFor(MySQL80)
	a := defCol("bigint", sp("9223372036854775806"), ColumnDefaultConstant)
	b := defCol("bigint", sp("9223372036854775807"), ColumnDefaultConstant)
	if n.CanonicalDefault(a) == n.CanonicalDefault(b) {
		t.Fatalf("distinct large bigint defaults must differ: %q vs %q", n.CanonicalDefault(a), n.CanonicalDefault(b))
	}
	// Quote-insensitivity and value-equality must still hold at full precision.
	c := defCol("bigint", sp("'9223372036854775807'"), ColumnDefaultConstant)
	if n.CanonicalDefault(b) != n.CanonicalDefault(c) {
		t.Fatal("quoted and unquoted large bigint default must be equal")
	}
	// Negative and negative-zero handling.
	if n.CanonicalDefault(defCol("int", sp("-0"), ColumnDefaultConstant)) !=
		n.CanonicalDefault(defCol("int", sp("0"), ColumnDefaultConstant)) {
		t.Fatal("-0 and 0 must be equal")
	}
	if n.CanonicalDefault(defCol("int", sp("-5"), ColumnDefaultConstant)) ==
		n.CanonicalDefault(defCol("int", sp("5"), ColumnDefaultConstant)) {
		t.Fatal("-5 and 5 must differ")
	}
}

func TestCanonicalGeneratedExpr_TokenBoundaryPreserved(t *testing.T) {
	// Whitespace must collapse to a single space between word tokens, not vanish, so
	// `a` and `b` (an AND expression) does not merge into the identifier `aandb`.
	n := NormalizerFor(MySQL80)
	if n.CanonicalGeneratedExpr("`a` and `b`") == n.CanonicalGeneratedExpr("`aandb`") {
		t.Fatal("`a` and `b` must not collapse to aandb")
	}
	// But spacing variants of the same expression are still equal.
	if n.CanonicalGeneratedExpr("`a`  and  `b`") != n.CanonicalGeneratedExpr("`a` and `b`") {
		t.Fatal("spacing variants of the same AND expression must be equal")
	}
}

func TestExpressionTokenizers_NoPanicOnMalformedInput(t *testing.T) {
	// The hand-written tokenizers must never panic on malformed/truncated input
	// (unterminated quotes, backslash escapes, missing parens, unterminated backticks).
	n := NormalizerFor(MySQL80)
	inputs := []string{
		"concat(`a`,'x",  // unterminated string
		`concat(a,'\'')`, // backslash-escaped quote
		"_latin1'x",      // introducer + unterminated literal
		"`unterminated",  // unterminated backtick
		"'",              // lone quote
		"",               // empty
		"((((",           // unbalanced parens
	}
	for _, in := range inputs {
		_ = n.CanonicalGeneratedExpr(in) // must not panic
	}
	// parseColumnType on malformed types must not panic.
	for _, in := range []string{"int(", "int(11", "a)b(c", "varchar(", "(", ")", "decimal(10,"} {
		_ = n.CanonicalColumnType(col("x", in)) // must not panic
	}
}

func TestCanonicalColumn_NoKeyInjectionViaComment(t *testing.T) {
	// A column whose comment contains the field-delimiter syntax must not collide with a
	// genuinely different column (length-delimited field encoding prevents injection).
	tbl := &Table{Name: "t", Charset: "utf8mb4", Collation: "utf8mb4_0900_ai_ci"}
	n := NormalizerFor(MySQL80)

	withComment := &Column{Name: "a", DataType: "int", ColumnType: "int", Nullable: true,
		Comment: "x|gen:3:a+1|stored:5:false|"}
	gen := &Column{Name: "a", DataType: "int", ColumnType: "int", Nullable: true,
		Comment: "x", Generated: &GeneratedColumnInfo{Expr: "a+1", Stored: false}}
	if n.CanonicalColumn(tbl, withComment) == n.CanonicalColumn(tbl, gen) {
		t.Fatal("a comment mimicking field syntax must not collide with a generated column")
	}
}

// ---------------------------------------------------------------------------
// Regression tests for review round 2.
// ---------------------------------------------------------------------------

func TestCanonicalDefault_DecimalRoundsNotTruncates(t *testing.T) {
	// MySQL ROUNDS a numeric default to the column scale on storage, with carry:
	// 0.999 → '1.00', 9.999 → '10.00', INT DEFAULT 1.9 → '2'.
	n := NormalizerFor(MySQL80)
	eq := func(colType, lit, storedReadback string) {
		a := n.CanonicalDefault(defCol(colType, sp(lit), ColumnDefaultConstant))
		b := n.CanonicalDefault(defCol(colType, sp(storedReadback), ColumnDefaultConstant))
		if a != b {
			t.Errorf("%s DEFAULT %s must equal stored %s: %q vs %q", colType, lit, storedReadback, a, b)
		}
	}
	eq("decimal(10,2)", "0.999", "'1.00'")
	eq("decimal(10,2)", "9.999", "'10.00'")
	eq("decimal(10,2)", "0.005", "'0.01'")
	eq("int", "1.9", "'2'")
	eq("decimal(10,2)", "0.994", "'0.99'") // rounds down
	// Negative rounding to zero must normalize sign.
	if n.CanonicalDefault(defCol("int", sp("-0.4"), ColumnDefaultConstant)) !=
		n.CanonicalDefault(defCol("int", sp("0"), ColumnDefaultConstant)) {
		t.Error("-0.4 at scale 0 must round to 0")
	}
}

func TestCanonicalGeneratedExpr_SpaceBeforeBacktickPreserved(t *testing.T) {
	// `a` and `b` must NOT become `a andb` — a space before a backticked identifier is a
	// token boundary even though the backtick byte is not an identifier char.
	n := NormalizerFor(MySQL80)
	got := n.CanonicalGeneratedExpr("`a` and `b`")
	if strings.Contains(got, "andb") {
		t.Fatalf("space before backtick must be preserved, got %q", got)
	}
	// The bare user form and the 8.0 backticked readback must still compare equal.
	if n.CanonicalGeneratedExpr("a and b") != n.CanonicalGeneratedExpr("`a` and `b`") {
		t.Fatalf("user `a and b` must equal readback `\\`a\\` and \\`b\\``: %q vs %q",
			n.CanonicalGeneratedExpr("a and b"), n.CanonicalGeneratedExpr("`a` and `b`"))
	}
}

func TestResolveColumnCharsetCollation_CollateOnlyDerivesCharset(t *testing.T) {
	// A column with only COLLATE (no CHARACTER SET) under a different-charset table: the
	// effective charset is the collation's charset, not the table's.
	tbl := &Table{Name: "t", Charset: "utf8mb4", Collation: "utf8mb4_0900_ai_ci"}
	// Loader derives col.Charset=latin1 from the collation and sets CollationExplicit.
	c := &Column{Name: "a", DataType: "varchar", ColumnType: "varchar(10)",
		Charset: "latin1", Collation: "latin1_german1_ci", CollationExplicit: true, Nullable: true}
	n := NormalizerFor(MySQL80)
	cs, coll := n.ResolveColumnCharsetCollation(tbl, c)
	if cs != "latin1" || coll != "latin1_german1_ci" {
		t.Fatalf("COLLATE-only column must resolve to (latin1, latin1_german1_ci), got (%q,%q)", cs, coll)
	}
}

func TestCanonicalTimestamp_NotNullFirstColumnGetsMagic(t *testing.T) {
	// EDFT=0: the first TIMESTAMP column gets the implicit CURRENT_TIMESTAMP default even
	// when explicitly NOT NULL (only an explicit NULL suppresses it).
	n := &Normalizer{Version: MySQL57, ExplicitDefaultsForTimestamp: false}
	a := &Column{Name: "a", DataType: "timestamp", ColumnType: "timestamp", Nullable: false, NullExplicit: true}
	tbl := tsTable(a)
	def, onUpdate := n.CanonicalTimestampDefaults(tbl, a)
	if def != "ts:CURRENT_TIMESTAMP" || onUpdate != "ts:CURRENT_TIMESTAMP" {
		t.Fatalf("first TIMESTAMP NOT NULL must get CURRENT_TIMESTAMP magic; got def=%q onUpdate=%q", def, onUpdate)
	}
	// Explicit NULL first column suppresses the magic.
	b := &Column{Name: "a", DataType: "timestamp", ColumnType: "timestamp", Nullable: true, NullExplicit: true}
	defB, _ := n.CanonicalTimestampDefaults(tsTable(b), b)
	if defB != defaultAbsentKey {
		t.Fatalf("first TIMESTAMP NULL must NOT get magic; got %q", defB)
	}
}

func TestCanonicalDefault_FloatFractionPreserved(t *testing.T) {
	// Bare FLOAT/DOUBLE defaults preserve their fraction (NOT rounded to int): MySQL
	// stores float 1.9 → '1.9', double 3.14159 → '3.14159', float 2.0 → '2'.
	n := NormalizerFor(MySQL80)
	// Distinct float defaults must NOT collapse.
	if n.CanonicalDefault(defCol("float", sp("1.9"), ColumnDefaultConstant)) ==
		n.CanonicalDefault(defCol("float", sp("2.1"), ColumnDefaultConstant)) {
		t.Fatal("FLOAT DEFAULT 1.9 and 2.1 must not collapse")
	}
	// Value-equality with the quoted readback and trailing-zero stripping still hold.
	eq := func(colType, lit, readback string) {
		a := n.CanonicalDefault(defCol(colType, sp(lit), ColumnDefaultConstant))
		b := n.CanonicalDefault(defCol(colType, sp(readback), ColumnDefaultConstant))
		if a != b {
			t.Errorf("%s DEFAULT %s must equal stored %s: %q vs %q", colType, lit, readback, a, b)
		}
	}
	eq("float", "1.9", "'1.9'")
	eq("double", "3.14159", "'3.14159'")
	eq("float", "2.0", "'2'")
	// A scaled float(10,2) DOES pad to scale.
	eq("float(10,2)", "1.5", "'1.50'")
}

// TestPartitionConstantFold_DateHelpers locks the date helpers used by the partition bound
// constant-folder (entry partition-constant-folding) without an engine: calcDayNumber reproduces
// MySQL TO_DAYS exactly, and parseDateTimeLiteral is strict.
func TestPartitionConstantFold_DateHelpers(t *testing.T) {
	// calcDayNumber == MySQL TO_DAYS (verified values against the live engine).
	days := []struct {
		y, m, d int
		want    int64
	}{
		{2020, 1, 1, 737790},
		{2021, 1, 1, 738156},
		{2010, 1, 1, 734138},
		{0, 0, 0, 0},
	}
	for _, c := range days {
		if got := calcDayNumber(c.y, c.m, c.d); got != c.want {
			t.Errorf("calcDayNumber(%d,%d,%d) = %d, want %d", c.y, c.m, c.d, got, c.want)
		}
	}

	// parseDateTimeLiteral accepts canonical date/datetime forms and rejects malformed ones.
	if y, mo, d, _, _, _, ok := parseDateTimeLiteral("2020-06-15"); !ok || y != 2020 || mo != 6 || d != 15 {
		t.Errorf("parseDateTimeLiteral(date) = %d-%d-%d ok=%v", y, mo, d, ok)
	}
	if _, _, _, hh, mi, ss, ok := parseDateTimeLiteral("2020-06-15 13:45:30"); !ok || hh != 13 || mi != 45 || ss != 30 {
		t.Errorf("parseDateTimeLiteral(datetime) time = %d:%d:%d ok=%v", hh, mi, ss, ok)
	}
	for _, bad := range []string{"2020/06/15", "2020-13-01", "2020-02-30", "not-a-date", "2020-06-15 25:00:00", "20200615"} {
		if _, _, _, _, _, _, ok := parseDateTimeLiteral(bad); ok {
			t.Errorf("parseDateTimeLiteral(%q) unexpectedly ok", bad)
		}
	}
}

// TestPartitionConstantFold_Bounds locks the end-to-end load-time fold through the public LoadSQL
// path: a non-literal bound folds to the engine's stored integer literal, an unsupported expression
// stays verbatim (the flagged residual), and a literal bound is unchanged.
func TestPartitionConstantFold_Bounds(t *testing.T) {
	boundOf := func(sql, part string) string {
		t.Helper()
		cat, err := LoadSQL(sql)
		if err != nil {
			t.Fatalf("load %q: %v", sql, err)
		}
		for _, db := range cat.Databases() {
			if tbl := db.GetTable("t"); tbl != nil && tbl.Partitioning != nil {
				for _, pd := range tbl.Partitioning.Partitions {
					if pd.Name == part {
						return pd.ValueExpr
					}
				}
			}
		}
		t.Fatalf("partition %q not found", part)
		return ""
	}
	const hdr = "CREATE DATABASE d; USE d; "
	cases := []struct {
		name, sql, part, want string
	}{
		{"add", hdr + "CREATE TABLE t (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (5+5));", "p0", "10"},
		{"mul", hdr + "CREATE TABLE t (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10*2));", "p0", "20"},
		{"div-int", hdr + "CREATE TABLE t (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (100 DIV 3));", "p0", "33"},
		{"neg-literal", hdr + "CREATE TABLE t (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (-10));", "p0", "-10"},
		{"plain-literal", hdr + "CREATE TABLE t (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (100));", "p0", "100"},
		{"todays", hdr + "CREATE TABLE t (id INT, dt DATE) PARTITION BY RANGE (TO_DAYS(dt)) (PARTITION p0 VALUES LESS THAN (TO_DAYS('2020-01-01')));", "p0", "737790"},
		// A DATE/TIMESTAMP literal arg folds symmetrically with the bare-string form above (stringLitValue
		// reads TemporalLit.Value), so DATE '2020-01-01' and '2020-01-01' canonicalize to the same bound.
		{"todays-date-literal", hdr + "CREATE TABLE t (id INT, dt DATE) PARTITION BY RANGE (TO_DAYS(dt)) (PARTITION p0 VALUES LESS THAN (TO_DAYS(DATE '2020-01-01')));", "p0", "737790"},
		{"todays-timestamp-literal", hdr + "CREATE TABLE t (id INT, dt DATETIME) PARTITION BY RANGE (TO_DAYS(dt)) (PARTITION p0 VALUES LESS THAN (TO_DAYS(TIMESTAMP '2020-01-01 00:00:00')));", "p0", "737790"},
		// A TIME literal is NOT folded (even one spelled to look like a date): its value is a clock
		// string, and TIME '2020-01-01' is rejected by MySQL itself — so the bound stays verbatim
		// rather than fold to a wrong day number.
		{"todays-time-literal-verbatim", hdr + "CREATE TABLE t (id INT, dt DATE) PARTITION BY RANGE (TO_DAYS(dt)) (PARTITION p0 VALUES LESS THAN (TO_DAYS(TIME '2020-01-01')));", "p0", "to_days(TIME'2020-01-01')"},
		{"year-fn", hdr + "CREATE TABLE t (dt DATE) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (YEAR('2020-06-15')));", "p0", "2020"},
		{"list-arith", hdr + "CREATE TABLE t (id INT) PARTITION BY LIST (id) (PARTITION p0 VALUES IN (1+1, 2+2));", "p0", "2,4"},
		// Flagged residual: an unsupported function stays verbatim (does not get a wrong fold).
		{"unsupported-verbatim", hdr + "CREATE TABLE t (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (ABS(-5)));", "p0", "abs(-5)"},
		// Deliberately-declined operators stay verbatim rather than fold to a wrong literal. MySQL
		// evaluates MOD and the bitwise/shift operators with UNSIGNED 64-bit semantics that signed
		// int64 does not match, so foldConstIntExpr declines them (the value the engine actually
		// stores is the verbatim expression on both versions; folding would phantom-diff).
		{"mod-verbatim", hdr + "CREATE TABLE t (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10 MOD 3));", "p0", "(10 % 3)"},
		{"bitand-verbatim", hdr + "CREATE TABLE t (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (6 & 3));", "p0", "(6 & 3)"},
		{"shift-verbatim", hdr + "CREATE TABLE t (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (1 << 4));", "p0", "(1 << 4)"},
		// Overflow is declined (not wrapped): the multiply would overflow int64, so the bound stays
		// verbatim. A non-evenly-dividing `/` likewise declines (the engine stores a decimal).
		{"overflow-verbatim", hdr + "CREATE TABLE t (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (9223372036854775807 * 2));", "p0", "(9223372036854775807 * 2)"},
		{"uneven-div-verbatim", hdr + "CREATE TABLE t (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10 / 3));", "p0", "(10 / 3)"},
		// A two-digit-year date literal is declined (MySQL's 70→1970 expansion is not reproduced),
		// so YEAR on it stays verbatim instead of folding to a wrong year.
		{"two-digit-year-verbatim", hdr + "CREATE TABLE t (dt DATE) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (YEAR('20-06-15')));", "p0", "year('20-06-15')"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := boundOf(c.sql, c.part); got != c.want {
				t.Errorf("ValueExpr = %q, want %q", got, c.want)
			}
		})
	}
}

// TestCanonicalPartitionText locks the diff-time whitespace canonicalization (entry
// partition-constant-folding) used by CanonicalPartitionValue/Expr. Incidental whitespace OUTSIDE a
// quoted string literal is collapsed so two specs that fold to the same values compare equal, while
// the content INSIDE a single-quoted string bound (a LIST/RANGE COLUMNS string value) is copied
// byte-for-byte — MySQL echoes string bounds verbatim, so a space or comma inside one is part of the
// value and must NOT be collapsed (doing so would mask a real bound change).
func TestCanonicalPartitionText(t *testing.T) {
	cases := []struct {
		name, in, want string
	}{
		// Whitespace outside strings collapses; spaces around commas/parens are dropped.
		{"collapse-runs", "10  +   5", "10 + 5"},
		{"comma-spacing", "1 ,  2 , 3", "1,2,3"},
		{"paren-spacing", "( 1 , 2 )", "(1,2)"},
		{"tab-newline", "1\t,\n2", "1,2"},
		// A quoted string bound is preserved exactly — interior spaces and commas are part of the
		// value, not separators.
		{"string-with-space", "'a b'", "'a b'"},
		{"string-with-comma", "'a, b'", "'a, b'"},
		{"tuple-mixed", "( 10 , 'x, y' )", "(10,'x, y')"},
		// Escapes inside the string must not end it early: a backslash-escaped quote and a doubled
		// '' quote both stay inside the literal.
		{"backslash-escaped-quote", "'a\\'b, c'", "'a\\'b, c'"},
		{"doubled-quote", "'a''b, c'", "'a''b, c'"},
		// Two adjacent string values separated by a comma: each is preserved, the separator space is
		// dropped.
		{"two-strings", "'a, b' , 'c d'", "'a, b','c d'"},
	}
	n := NormalizerFor(MySQL80)
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := n.CanonicalPartitionExpr(c.in); got != c.want {
				t.Errorf("CanonicalPartitionExpr(%q) = %q, want %q", c.in, got, c.want)
			}
			// CanonicalPartitionValue shares the same text canonicalizer.
			if got := n.CanonicalPartitionValue(c.in); got != c.want {
				t.Errorf("CanonicalPartitionValue(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}
