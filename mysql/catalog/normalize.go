package catalog

// Package-level diff-time canonicalization for MySQL SDL (declarative) migration.
//
// MySQL implicitly rewrites DDL on storage: it drops integer display widths (8.0),
// injects them (5.7), folds BOOLEAN to tinyint(1), echoes or strips column charset
// clauses depending on version, single-quotes numeric defaults, and so on. A
// declarative diff compares the user's TARGET schema against the SYNCED (stored)
// schema read back via SHOW CREATE TABLE. Unless both sides are reduced to MySQL's
// canonical stored form first, every implicit rewrite produces a phantom diff forever.
//
// This file is the single owner of that canonical form. diff-core calls the
// Canonical* / Resolve* methods below wherever it compares columns, types, defaults,
// charset/collation, generated expressions, or table options. The defining property,
// proven against the live oracle, is:
//
//	Canonical(userTargetForm) == Canonical(syncedStoredForm)   (for a given engine version)
//
// Every entry in docs/migration/mysql/sdl/normalization.md is a readback that this
// layer must collapse onto its target form. Behavior that diverges across MySQL 5.7
// and 8.0 is branched on Version; behavior that depends on a session variable
// (explicit_defaults_for_timestamp) is taken from the Normalizer context, never
// hardcoded to a version (see the flag explicit-defaults-for-timestamp-is-session-var).

import (
	"fmt"
	"strconv"
	"strings"
)

// Version identifies the target MySQL major version whose stored form the
// canonicalizer must reproduce. Charset spelling, integer display width, YEAR width,
// descending indexes, the utf8mb4 default collation, CHECK persistence, functional
// defaults, and the expression charset introducer all diverge across these.
type Version int

const (
	// MySQL57 reproduces MySQL 5.7's stored form (display widths injected,
	// utf8 spelling, utf8mb4_general_ci default, no descending indexes, CHECK dropped).
	MySQL57 Version = iota
	// MySQL80 reproduces MySQL 8.0's stored form (display widths dropped except
	// tinyint(1)/zerofill, utf8mb3 spelling, utf8mb4_0900_ai_ci default, descending
	// indexes, CHECK kept, functional defaults, _latin1 expression introducer).
	MySQL80
)

// Normalizer carries the engine version plus the session/server variables that the
// canonical stored form depends on. diff-core constructs one per side from the
// synced schema's captured version and variables, then resolves both the target and
// the current catalog through it so the two are compared in the same normal form.
type Normalizer struct {
	Version Version
	// ExplicitDefaultsForTimestamp is the effective explicit_defaults_for_timestamp.
	// It governs the implicit NOT NULL / DEFAULT CURRENT_TIMESTAMP magic on bare
	// TIMESTAMP columns. It is a session/server variable, NOT a version trait: a 5.7
	// server can run with it ON and an 8.0 server with it OFF. diff-core must capture
	// it per source database and pass it here. NormalizerFor seeds it to each
	// version's box default for convenience, but callers should override it from the
	// synced schema when known.
	ExplicitDefaultsForTimestamp bool
}

// NormalizerFor builds a Normalizer for a version, seeding session variables to that
// version's oracle-box default (explicit_defaults_for_timestamp=0 on 5.7, 1 on 8.0).
// Callers with a captured value should set ExplicitDefaultsForTimestamp explicitly.
func NormalizerFor(v Version) *Normalizer {
	return &Normalizer{
		Version:                      v,
		ExplicitDefaultsForTimestamp: v == MySQL80,
	}
}

// is80 reports whether the target reproduces 8.0 stored form.
func (n *Normalizer) is80() bool { return n.Version == MySQL80 }

// ---------------------------------------------------------------------------
// Charset / collation — the dominant phantom-diff surface (HIGHEST RISK).
// ---------------------------------------------------------------------------

// foldCharset maps a charset name to a version-independent identity so that the two
// spellings of the same charset compare equal. MySQL stores utf8 (5.7) vs utf8mb3
// (8.0) for the same charset; folding both to utf8mb3 lets the resolved pair compare
// equal regardless of which version's surface spelling each side carried.
// (entry utf8-utf8mb3-alias)
func foldCharset(cs string) string {
	cs = strings.ToLower(strings.TrimSpace(cs))
	if cs == "utf8" {
		return "utf8mb3"
	}
	return cs
}

// foldCollation maps a collation name to a version-independent identity. It rewrites a
// utf8_* collation prefix to utf8mb3_* to match foldCharset, so utf8_general_ci (5.7)
// and utf8mb3_general_ci (8.0) — the same collation — compare equal.
func foldCollation(coll string) string {
	coll = strings.ToLower(strings.TrimSpace(coll))
	if strings.HasPrefix(coll, "utf8_") {
		return "utf8mb3_" + coll[len("utf8_"):]
	}
	return coll
}

// defaultCollationFor returns the version-correct default collation for a charset,
// already folded. The utf8mb4 default is the single most dangerous case: 8.0 uses
// utf8mb4_0900_ai_ci (which does not even exist on 5.7) while 5.7 uses
// utf8mb4_general_ci. (entry utf8mb4-default-collation)
func (n *Normalizer) defaultCollationFor(charset string) string {
	cs := foldCharset(charset)
	if cs == "utf8mb4" {
		if n.is80() {
			return "utf8mb4_0900_ai_ci"
		}
		return "utf8mb4_general_ci"
	}
	if dc, ok := defaultCollationForCharset[cs]; ok {
		return foldCollation(dc)
	}
	// Last resort: derive <charset>_general_ci.
	if cs != "" {
		return foldCollation(cs + "_general_ci")
	}
	return ""
}

// ResolveColumnCharsetCollation returns the effective (charset, collation) pair of a
// column, folded to a version-independent identity, so that diff-core can compare the
// resolved pair instead of the surface CHARACTER SET/COLLATE tokens. This is the
// canonical strategy for the whole charset family (flag
// column-charset-collation-echo-asymmetry): surface-token comparison phantom-diffs
// across versions because 5.7 STRIPS a column's charset when it equals the table's
// while 8.0 KEEPS it as a full pair; resolving the effective pair on BOTH sides and
// comparing that is version-robust.
//
// Resolution rules (entries column-charset-echo-57/-80, -only-collation-resolution),
// driven by the loader's explicit-clause flags so the three cases MySQL distinguishes
// are kept apart even though the loader fills Charset/Collation from inheritance:
//   - charset: the column's explicit CHARACTER SET if it wrote one, else the table's.
//   - collation: the column's explicit COLLATE if it wrote one; else, if it wrote only a
//     CHARACTER SET, that charset's default collation (NOT the table COLLATE — the subtle
//     -only-collation-resolution rule); else (a bare column) the table's effective
//     collation, re-resolved to the target version.
//
// For non-string columns (no charset concept) it returns ("", "").
func (n *Normalizer) ResolveColumnCharsetCollation(table *Table, c *Column) (charset, collation string) {
	if !isStringType(c.DataType) && !isEnumSetType(c.DataType) {
		return "", ""
	}

	tableCharset := foldCharset(table.Charset)
	tableCollation := n.effectiveTableCollation(table)

	// Effective charset:
	//   - an explicit column CHARACTER SET wins;
	//   - an explicit column COLLATE (without CHARACTER SET) implies the collation's
	//     charset — the loader derives c.Charset from the collation, so honor it;
	//   - otherwise the column inherits the table charset.
	switch {
	case c.CharsetExplicit && c.Charset != "":
		charset = foldCharset(c.Charset)
	case c.CollationExplicit && c.Charset != "":
		charset = foldCharset(c.Charset)
	default:
		charset = tableCharset
	}

	// Effective collation, resolved with the version-correct rules. The explicit flags
	// let us tell apart the three otherwise-indistinguishable cases the loader collapses:
	switch {
	case c.CollationExplicit && c.Collation != "":
		// The user wrote COLLATE: honor it verbatim (folded).
		collation = foldCollation(c.Collation)
	case c.CharsetExplicit && c.Charset != "":
		// The user wrote only CHARACTER SET: the collation is THAT charset's default,
		// never the table COLLATE. (entry column-charset-only-collation-resolution)
		collation = n.defaultCollationFor(c.Charset)
	default:
		// Bare column: inherits the table's effective collation (re-resolved to the
		// target version). MySQL stores a bare column with the table's collation, not the
		// charset default — so a table-level COLLATE change must propagate to bare columns.
		collation = tableCollation
	}
	return charset, collation
}

// effectiveTableCollation returns the table's collation, folded and re-resolved to the
// target version. Two distinct "this is really the charset default" cases must both
// collapse onto the version-correct default so a bare CHARSET table never phantom-diffs:
//
//   - The omni loader bakes the STATIC (8.0) charset-default collation into
//     table.Collation when the table has no explicit COLLATE (it cannot know the target
//     version at load time) — so a user's `... DEFAULT CHARSET=utf8mb4` arrives carrying
//     utf8mb4_0900_ai_ci even when the target is 5.7.
//   - The SYNCED side is dumped by MetadataToSDL with the server's STORED default spelled
//     out as an explicit COLLATE (a 5.7 utf8mb4 table reads back
//     `DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci`) — so it arrives carrying the
//     version's stored default as an "explicit" clause.
//
// Both are the charset default for SOME version and are indistinguishable from a bare
// charset on the engine, so we fold any collation equal to either the loader-baked
// static default OR the version-appropriate default down to defaultCollationFor. A
// genuinely non-default table COLLATE (e.g. utf8mb4_unicode_ci) is honored as-is.
func (n *Normalizer) effectiveTableCollation(table *Table) string {
	tableCharset := foldCharset(table.Charset)
	tableCollation := foldCollation(table.Collation)
	if tableCollation == "" {
		return n.defaultCollationFor(tableCharset)
	}
	if n.isCharsetDefaultCollation(tableCharset, tableCollation) {
		return n.defaultCollationFor(tableCharset)
	}
	return tableCollation
}

// isCharsetDefaultCollation reports whether folded collation `coll` is "the default" for
// folded charset `cs` for the purpose of effective-collation resolution: it matches
// either the version-appropriate default (defaultCollationFor) or the static collation
// the omni loader bakes for a bare charset (the 8.0 default, from
// defaultCollationForCharset). The utf8mb4 case is the only one whose default diverges by
// version (5.7 utf8mb4_general_ci vs 8.0 utf8mb4_0900_ai_ci); folding both spellings keeps
// a 5.7-synced `COLLATE=utf8mb4_general_ci` and a loader-baked `utf8mb4_0900_ai_ci` equal
// to each other and to a bare charset under either version's normalizer.
func (n *Normalizer) isCharsetDefaultCollation(cs, coll string) bool {
	if coll == "" {
		return true
	}
	if coll == n.defaultCollationFor(cs) {
		return true
	}
	if baked, ok := defaultCollationForCharset[cs]; ok && coll == foldCollation(baked) {
		return true
	}
	return false
}

// ---------------------------------------------------------------------------
// Column type — integer display width, BOOLEAN, decimal/float, year, zerofill.
// ---------------------------------------------------------------------------

// parsedType is a decomposed column-type string. The loader stores the 8.0 form (e.g.
// "int"); raw readback strings may carry widths (e.g. "int(11)"). CanonicalColumnType
// parses either and re-emits the version-correct canonical form, so both collapse onto
// one value and never phantom-diff.
type parsedType struct {
	base     string // lowercased base type name (already alias-folded: numeric→decimal, etc.)
	length   int    // first paren arg (display width, char length, precision); 0 = absent
	scale    int    // second paren arg (decimal/float scale); -1 = absent
	hasLen   bool
	hasScale bool
	unsigned bool
	zerofill bool
	suffix   string // verbatim tail for types we do not reformat (enum/set member list)
}

// parseColumnType decomposes a column-type string into base/length/scale/modifiers.
func parseColumnType(columnType string) parsedType {
	s := strings.TrimSpace(columnType)
	low := strings.ToLower(s)

	p := parsedType{scale: -1}

	// Trailing modifiers.
	if strings.Contains(low, "zerofill") {
		p.zerofill = true
		p.unsigned = true // zerofill implies unsigned
	}
	if strings.Contains(low, "unsigned") {
		p.unsigned = true
	}
	// Strip the modifier words to isolate "<base>[(args)]".
	low = strings.TrimSpace(strings.NewReplacer("unsigned", "", "zerofill", "").Replace(low))

	// enum/set: keep the member list verbatim (order is significant), normalize quoting.
	if strings.HasPrefix(low, "enum(") || strings.HasPrefix(low, "set(") {
		open := strings.IndexByte(low, '(')
		p.base = normalizedTypeName(low[:open])
		p.suffix = s[open:] // original-cased member list
		return p
	}

	if open := strings.IndexByte(low, '('); open >= 0 {
		closeIdx := strings.LastIndexByte(low, ')')
		// Guard against a malformed type with no closing paren, or ')' before '('
		// (e.g. a truncated readback "int(11"): treat the paren as absent.
		if closeIdx <= open {
			p.base = normalizedTypeName(strings.TrimSpace(low[:open]))
			return p
		}
		argStr := low[open+1 : closeIdx]
		name := strings.TrimSpace(low[:open])
		p.base = normalizedTypeName(name)
		args := strings.SplitN(argStr, ",", 2)
		if v, err := strconv.Atoi(strings.TrimSpace(args[0])); err == nil {
			p.length = v
			p.hasLen = true
		}
		if len(args) == 2 {
			if v, err := strconv.Atoi(strings.TrimSpace(args[1])); err == nil {
				p.scale = v
				p.hasScale = true
			}
		}
	} else {
		p.base = normalizedTypeName(low)
	}
	return p
}

// CanonicalColumnType returns the canonical stored column type for the target version.
// It folds aliases (BOOLEAN/BOOL→tinyint(1), NUMERIC/DEC/FIXED→decimal, REAL→double,
// integer→int), applies the version's integer-display-width rule (8.0 drops widths
// except tinyint(1) and zerofill; 5.7 injects the type's default width), pads DECIMAL
// precision/scale, drops single-arg FLOAT precision, injects length-1 for bare
// CHAR/BINARY/BIT, and reproduces the version's YEAR width. The result is a stable
// comparison key, not necessarily a SHOW CREATE rendering.
func (n *Normalizer) CanonicalColumnType(c *Column) string {
	p := parseColumnType(c.ColumnType)
	return n.canonicalType(p)
}

func (n *Normalizer) canonicalType(p parsedType) string {
	switch p.base {
	case "boolean":
		// BOOLEAN/BOOL → tinyint(1) on both versions. (parseColumnType already folds
		// "bool" → "boolean" via normalizedTypeName, so only "boolean" reaches here.)
		return "tinyint(1)"
	case "serial":
		return "bigint unsigned"
	}

	if isIntegerType(p.base) {
		return n.canonicalIntType(p)
	}
	if p.base == "decimal" {
		length, scale := p.length, p.scale
		if !p.hasLen {
			length = 10
		}
		if !p.hasScale {
			scale = 0
		}
		return n.withUnsigned(fmt.Sprintf("decimal(%d,%d)", length, scale), p)
	}
	if p.base == "float" || p.base == "double" {
		base := p.base
		// FLOAT(M,D)/DOUBLE(M,D) retained; single-arg FLOAT(p) precision dropped;
		// FLOAT with precision > 24 was already folded to double upstream — fold here too.
		if p.base == "float" && p.hasLen && !p.hasScale && p.length > 24 {
			base = "double"
		}
		if p.hasLen && p.hasScale {
			return n.withUnsigned(fmt.Sprintf("%s(%d,%d)", base, p.length, p.scale), p)
		}
		return n.withUnsigned(base, p)
	}

	switch p.base {
	case "year":
		if n.is80() {
			return "year"
		}
		return "year(4)"
	case "bit":
		if !p.hasLen {
			return "bit(1)"
		}
		return fmt.Sprintf("bit(%d)", p.length)
	case "char", "binary":
		if !p.hasLen {
			return p.base + "(1)"
		}
		return fmt.Sprintf("%s(%d)", p.base, p.length)
	case "enum", "set":
		return p.base + canonicalEnumSuffix(p.suffix)
	}

	// Text/blob: never carry a width.
	if isTextBlobLengthStripped(p.base) {
		return p.base
	}

	// Default: keep length/scale verbatim (varchar(N), varbinary(N), time(N), etc.).
	if p.hasLen && p.hasScale {
		return fmt.Sprintf("%s(%d,%d)", p.base, p.length, p.scale)
	}
	if p.hasLen {
		return fmt.Sprintf("%s(%d)", p.base, p.length)
	}
	return p.base
}

// canonicalIntType applies the version-specific integer display-width rule.
// p.base is already alias-folded by normalizedTypeName (integer→int), so it is one of
// tinyint/smallint/mediumint/int/bigint here.
func (n *Normalizer) canonicalIntType(p parsedType) string {
	base := p.base

	// tinyint(1) is the boolean marker: width survives on BOTH versions (only when
	// signed and not zerofill — tinyint(1) unsigned/zerofill is a plain narrow int).
	if base == "tinyint" && p.hasLen && p.length == 1 && !p.unsigned && !p.zerofill {
		return "tinyint(1)"
	}

	// ZEROFILL retains the display width on both versions (8.0 does NOT drop it).
	if p.zerofill {
		width := p.length
		if !p.hasLen {
			width = defaultIntDisplayWidth(base, true)
		}
		return fmt.Sprintf("%s(%d) unsigned zerofill", base, width)
	}

	if n.is80() {
		// 8.0 drops the display width from all (non-zerofill) integer types.
		return n.withUnsigned(base, p)
	}
	// 5.7 keeps/injects the default width.
	width := p.length
	if !p.hasLen {
		width = defaultIntDisplayWidth(base, p.unsigned)
	}
	return n.withUnsigned(fmt.Sprintf("%s(%d)", base, width), p)
}

// withUnsigned appends " unsigned" to a rendered type when the column is unsigned
// (but not zerofill, which renders its own "unsigned zerofill").
func (n *Normalizer) withUnsigned(typeStr string, p parsedType) string {
	if p.unsigned && !p.zerofill {
		return typeStr + " unsigned"
	}
	return typeStr
}

// canonicalEnumSuffix normalizes an enum/set member list "('a','b',...)" to a canonical
// form: each member single-quoted, order preserved (order is semantically significant
// for ENUM/SET and is never reordered). Double-quoted members become single-quoted.
// (entry enum-set-quoting)
func canonicalEnumSuffix(suffix string) string {
	s := strings.TrimSpace(suffix)
	if len(s) < 2 || s[0] != '(' || s[len(s)-1] != ')' {
		return suffix
	}
	inner := s[1 : len(s)-1]
	members := splitEnumMembers(inner)
	parts := make([]string, len(members))
	for i, m := range members {
		parts[i] = "'" + strings.ReplaceAll(m, "'", "''") + "'"
	}
	return "(" + strings.Join(parts, ",") + ")"
}

// splitEnumMembers splits a comma-separated quoted member list into decoded member
// contents, honoring single- or double-quoting and doubled-quote escapes.
func splitEnumMembers(inner string) []string {
	var members []string
	i := 0
	for i < len(inner) {
		// Skip separators/whitespace.
		for i < len(inner) && (inner[i] == ',' || inner[i] == ' ' || inner[i] == '\t') {
			i++
		}
		if i >= len(inner) {
			break
		}
		q := inner[i]
		if q != '\'' && q != '"' {
			// Unexpected; bail with whatever remains as one member.
			members = append(members, inner[i:])
			break
		}
		i++
		var b strings.Builder
		for i < len(inner) {
			if inner[i] == q {
				if i+1 < len(inner) && inner[i+1] == q {
					b.WriteByte(q)
					i += 2
					continue
				}
				i++
				break
			}
			b.WriteByte(inner[i])
			i++
		}
		members = append(members, b.String())
	}
	return members
}

// ---------------------------------------------------------------------------
// Column default — value-based canonicalization.
// ---------------------------------------------------------------------------

// CanonicalDefault returns a comparison key for a column's DEFAULT. It is value-based,
// not string-based (flag numeric-default-quote-version): DEFAULT 0 and DEFAULT '0'
// yield the same key, because MySQL quotes numeric defaults in storage but the exact
// quoting drifts across 8.0 point releases. Rules:
//   - absent default OR explicit NULL → a single "no value" key (they are equivalent
//     for a nullable column; the differ compares nullability separately).
//   - CURRENT_TIMESTAMP family (incl. NOW/LOCALTIME/LOCALTIMESTAMP and (N) variants)
//     → "CURRENT_TIMESTAMP[(N)]" unquoted, by DefaultKind.
//   - functional DEFAULT (expression) → the expression run through the same generated-
//     expression normalizer (8.0-only construct).
//   - numeric literal on a scaled DECIMAL → padded to scale ('0' → '0.00') then keyed
//     as the numeric value.
//   - numeric literal otherwise → keyed by parsed numeric value (quote-insensitive).
//   - string literal → keyed by its decoded content (quote style/charset-independent).
func (n *Normalizer) CanonicalDefault(c *Column) string {
	// Absent / dropped default.
	if c.Default == nil {
		return defaultAbsentKey
	}
	raw := strings.TrimSpace(*c.Default)
	if strings.EqualFold(raw, "NULL") {
		return defaultAbsentKey
	}

	switch c.DefaultKind {
	case ColumnDefaultCurrentTimestamp:
		if ts, ok := normalizeCurrentTimestamp(raw); ok {
			return "ts:" + ts
		}
		return "ts:CURRENT_TIMESTAMP"
	case ColumnDefaultExpression:
		return "expr:" + n.CanonicalGeneratedExpr(stripOuterParens(raw))
	}

	// Constant. Decide numeric vs string by the column's data type and the literal shape.
	unq, wasQuoted := unquoteDefault(raw)

	if isNumericType(c.DataType) {
		// Numeric column: compare by numeric value, padded to the column's scale.
		if key, ok := numericDefaultKey(c.ColumnType, unq); ok {
			return key
		}
	}
	// Bit literal default (b'...' / 0b...).
	if looksLikeBitLiteral(raw) {
		return "bit:" + canonicalBitLiteral(raw)
	}
	// If it parses as a number and the column is not obviously a string type, key by value.
	if !wasQuoted && !isStringType(c.DataType) && !isEnumSetType(c.DataType) {
		if key, ok := numericDefaultKey(c.ColumnType, unq); ok {
			return key
		}
	}
	// String/other literal: key by decoded content.
	return "str:" + unq
}

const defaultAbsentKey = "<none>"

// normalizeCurrentTimestamp folds CURRENT_TIMESTAMP synonyms (NOW, LOCALTIME,
// LOCALTIMESTAMP) and (N) variants to "CURRENT_TIMESTAMP[(N)]".
func normalizeCurrentTimestamp(val string) (string, bool) {
	if ts, ok := currentTimestampSQL(val); ok {
		return ts, true
	}
	upper := strings.ToUpper(strings.TrimSpace(val))
	switch upper {
	case "LOCALTIME", "LOCALTIME()", "LOCALTIMESTAMP", "LOCALTIMESTAMP()":
		return "CURRENT_TIMESTAMP", true
	}
	for _, prefix := range []string{"LOCALTIME(", "LOCALTIMESTAMP("} {
		if strings.HasPrefix(upper, prefix) && strings.HasSuffix(upper, ")") {
			return "CURRENT_TIMESTAMP" + upper[len(prefix)-1:], true
		}
	}
	return "", false
}

// unquoteDefault strips one layer of surrounding single or double quotes from a default
// literal and un-doubles embedded quotes, returning the decoded content and whether it
// was quoted.
func unquoteDefault(s string) (string, bool) {
	if len(s) >= 2 {
		if (s[0] == '\'' && s[len(s)-1] == '\'') || (s[0] == '"' && s[len(s)-1] == '"') {
			q := s[0]
			inner := s[1 : len(s)-1]
			inner = strings.ReplaceAll(inner, string([]byte{q, q}), string(q))
			return inner, true
		}
	}
	return s, false
}

// numericDefaultKey renders a canonical value key for a numeric default by normalizing
// the literal AS A DECIMAL STRING — never via float64, which would lose precision for
// BIGINT/DECIMAL values beyond 2^53 and collide distinct defaults. It normalizes sign and
// leading zeros, then applies the column-kind-specific storage rule: scaled types
// (decimal(M,D), float(M,D)) round to scale with carry; integer types and bare DECIMAL
// round to integer; bare FLOAT/DOUBLE preserve their fraction (trailing zeros stripped).
// So '0' and '0.00' compare equal while 9223372036854775806 and ...807 stay distinct.
func numericDefaultKey(columnType, literal string) (string, bool) {
	sign, intPart, fracPart, ok := parseDecimalLiteral(literal)
	if !ok {
		return "", false
	}
	p := parseColumnType(columnType)
	switch {
	case p.hasScale:
		// A scaled type — decimal(M,D) or float(M,D)/double(M,D) — ROUNDS the default
		// (half-up) to the column scale on storage, with carry into the integer part
		// (MySQL: decimal(10,2) 0.999 → '1.00', 9.999 → '10.00'; float(10,2) 1.5 → '1.50').
		intPart, fracPart = roundToScale(intPart, fracPart, p.scale)
	case isIntegerType(p.base) || p.base == "decimal":
		// Integer types, and bare DECIMAL (which has implicit scale 0): round to zero
		// fractional digits, with carry (INT DEFAULT 1.9 → '2', decimal(10,0) 0.9 → '1').
		intPart, fracPart = roundToScale(intPart, fracPart, 0)
	default:
		// Bare FLOAT/DOUBLE/REAL (no explicit scale): the fraction is preserved verbatim
		// with only trailing zeros stripped (MySQL: float 1.9 → '1.9', 2.0 → '2',
		// double 3.14159 → '3.14159'). Rounding here would collapse distinct defaults.
		fracPart = strings.TrimRight(fracPart, "0")
	}
	// Re-normalize a negative zero produced by rounding (e.g. -0.4 at scale 0 → 0).
	if sign == "-" && intPart == "0" && strings.Trim(fracPart, "0") == "" {
		sign = ""
	}
	key := sign + intPart
	if fracPart != "" {
		key += "." + fracPart
	}
	return "num:" + key, true
}

// parseDecimalLiteral splits a numeric literal into (sign, integerDigits, fractionDigits)
// with leading integer zeros stripped (but a single "0" kept). Returns ok=false for
// non-numeric input or a literal with no digits at all. Exponent notation is not used by
// MySQL stored numeric defaults.
func parseDecimalLiteral(literal string) (sign, intPart, fracPart string, ok bool) {
	s := strings.TrimSpace(literal)
	if s == "" {
		return "", "", "", false
	}
	if s[0] == '+' || s[0] == '-' {
		if s[0] == '-' {
			sign = "-"
		}
		s = s[1:]
	}
	if s == "" {
		return "", "", "", false
	}
	intDigits, fracDigits := s, ""
	if dot := strings.IndexByte(s, '.'); dot >= 0 {
		intDigits = s[:dot]
		fracDigits = s[dot+1:]
	}
	if !allDigits(intDigits) || !allDigits(fracDigits) {
		return "", "", "", false
	}
	// Require at least one digit somewhere (reject "." / "+" / "-").
	if intDigits == "" && fracDigits == "" {
		return "", "", "", false
	}
	intDigits = strings.TrimLeft(intDigits, "0")
	if intDigits == "" {
		intDigits = "0"
	}
	// Normalize negative zero to positive.
	if intDigits == "0" && strings.Trim(fracDigits, "0") == "" {
		sign = ""
	}
	return sign, intDigits, fracDigits, true
}

// roundToScale rounds (intPart, fracPart) half-up to exactly scale fractional digits,
// propagating carry into the integer digits. Both inputs are bare digit strings; it
// returns the rounded (integerDigits, fractionDigits). Operating on digit strings keeps
// full precision for values beyond float64 range.
func roundToScale(intPart, fracPart string, scale int) (string, string) {
	if len(fracPart) <= scale {
		// No rounding needed; pad to scale.
		return intPart, fracPart + strings.Repeat("0", scale-len(fracPart))
	}
	kept := fracPart[:scale]
	roundUp := fracPart[scale] >= '5'
	if !roundUp {
		return intPart, kept
	}
	// Increment the digit string formed by intPart+kept by one, then split back.
	combined := incrementDigits(intPart + kept)
	// combined may have grown by one digit (carry out of the most-significant place).
	if scale == 0 {
		return combined, ""
	}
	return combined[:len(combined)-scale], combined[len(combined)-scale:]
}

// incrementDigits returns the decimal digit string plus one (e.g. "999" → "1000").
func incrementDigits(digits string) string {
	b := []byte(digits)
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] < '9' {
			b[i]++
			return string(b)
		}
		b[i] = '0'
	}
	return "1" + string(b)
}

func allDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// isNumericType reports whether a base data type carries a numeric default.
func isNumericType(dt string) bool {
	switch dt {
	case "tinyint", "smallint", "mediumint", "int", "integer", "bigint",
		"decimal", "numeric", "dec", "fixed", "float", "double", "real", "bit":
		return true
	}
	return false
}

func looksLikeBitLiteral(s string) bool {
	low := strings.ToLower(strings.TrimSpace(s))
	return strings.HasPrefix(low, "b'") || strings.HasPrefix(low, "0b")
}

// stripOuterParens removes one balanced layer of outer parentheses, used to normalize
// functional-default expressions whose stored form is double-parenthesized.
func stripOuterParens(s string) string {
	s = strings.TrimSpace(s)
	for len(s) >= 2 && s[0] == '(' && s[len(s)-1] == ')' && balancedFirstParenWrapsAll(s) {
		s = strings.TrimSpace(s[1 : len(s)-1])
	}
	return s
}

// balancedFirstParenWrapsAll reports whether the opening paren at index 0 matches the
// closing paren at the end (i.e. the whole string is one parenthesized group).
func balancedFirstParenWrapsAll(s string) bool {
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 && i != len(s)-1 {
				return false
			}
		}
	}
	return depth == 0
}

// ---------------------------------------------------------------------------
// Generated / functional-default expression canonicalization.
// ---------------------------------------------------------------------------

// CanonicalGeneratedExpr returns a comparison key for a generated-column expression,
// functional DEFAULT (expression), or CHECK expression. MySQL rewrites these on
// storage: identifiers are backticked, binary operators single-spaced, an outer paren
// pair is added, function names lowercased, and — on 8.0 only — string literals gain a
// charset introducer (_latin1'x') reflecting the connection charset at create time.
// The introducer is connection-dependent noise that makes 8.0 readbacks non-portable,
// so this canonical key strips it (entry generated-expr-string-introducer), collapses
// whitespace, removes outer parentheses, and lowercases identifiers/keywords outside
// string literals. The result lets a user's expression, a 5.7 readback, and an 8.0
// readback of the same logical expression all compare equal.
func (n *Normalizer) CanonicalGeneratedExpr(expr string) string {
	s := stripCharsetIntroducers(expr)
	s = stripOuterParens(s)
	s = collapseExprWhitespace(s)
	return s
}

// scanStringLiteral returns the index just past a single-quoted string literal that
// starts at i (expr[i]=='\”), honoring doubled-quote escapes (”) and backslash
// escapes (\'). If the literal is unterminated it returns len(expr). The returned index
// is always a valid slice bound (≤ len(expr)), so callers never read past the end.
func scanStringLiteral(expr string, i int) int {
	j := i + 1
	for j < len(expr) {
		switch expr[j] {
		case '\\':
			j += 2 // skip the escaped char (\' \\ etc.)
		case '\'':
			if j+1 < len(expr) && expr[j+1] == '\'' {
				j += 2 // doubled-quote escape
				continue
			}
			return j + 1 // past the closing quote
		default:
			j++
		}
	}
	return len(expr) // unterminated: consume to EOF
}

// stripCharsetIntroducers removes a charset introducer (e.g. _latin1, _utf8mb4) that
// immediately precedes a single-quoted string literal, leaving the bare literal.
func stripCharsetIntroducers(expr string) string {
	var b strings.Builder
	i := 0
	for i < len(expr) {
		c := expr[i]
		if c == '\'' {
			end := scanStringLiteral(expr, i)
			b.WriteString(expr[i:end])
			i = end
			continue
		}
		if c == '_' {
			// Possible charset introducer: _<ident>' — drop the "_<ident>" if a quote
			// follows the identifier run.
			j := i + 1
			for j < len(expr) && isIdentByte(expr[j]) {
				j++
			}
			if j < len(expr) && expr[j] == '\'' && j > i+1 {
				// Skip the introducer; the literal is copied on the next loop turn.
				i = j
				continue
			}
		}
		b.WriteByte(c)
		i++
	}
	return b.String()
}

func isIdentByte(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

// collapseExprWhitespace lowercases the expression outside string literals, strips
// backtick identifier quoting (lowercasing the inner name), and collapses each run of
// whitespace to a SINGLE space, so cosmetic spacing and quoting differences do not diff
// while token boundaries are preserved (e.g. `a` and `b` must not merge into `aandb`).
// String literals are copied verbatim (case and spacing inside them are significant) and
// unterminated literals/identifiers are consumed safely to EOF.
func collapseExprWhitespace(expr string) string {
	var b strings.Builder
	i := 0
	for i < len(expr) {
		c := expr[i]
		switch c {
		case '\'':
			end := scanStringLiteral(expr, i)
			b.WriteString(expr[i:end])
			i = end
		case '`':
			// Backticks are identifier quoting; drop them and lowercase the inner name
			// so a backticked identifier and the bare identifier compare equal. An
			// unterminated backtick is consumed to EOF.
			j := i + 1
			for j < len(expr) && expr[j] != '`' {
				j++
			}
			inner := expr[i+1 : j]
			for k := 0; k < len(inner); k++ {
				ch := inner[k]
				if ch >= 'A' && ch <= 'Z' {
					b.WriteByte(ch + ('a' - 'A'))
				} else {
					b.WriteByte(ch)
				}
			}
			if j < len(expr) {
				j++ // past the closing backtick
			}
			i = j
		case ' ', '\t', '\n', '\r':
			// Collapse a run of whitespace to a single space, but only between tokens —
			// drop it entirely at the start/end and adjacent to punctuation where it is
			// never significant. Keeping one space preserves keyword/identifier
			// boundaries so `a and b` does not become `aandb`.
			j := i
			for j < len(expr) && (expr[j] == ' ' || expr[j] == '\t' || expr[j] == '\n' || expr[j] == '\r') {
				j++
			}
			if exprSpaceIsSignificant(b.String(), expr, j) {
				b.WriteByte(' ')
			}
			i = j
		default:
			if c >= 'A' && c <= 'Z' {
				b.WriteByte(c + ('a' - 'A'))
			} else {
				b.WriteByte(c)
			}
			i++
		}
	}
	return b.String()
}

// exprSpaceIsSignificant reports whether a whitespace run between the already-emitted
// prefix and the next source byte (expr[next]) must be preserved as a single space.
// Whitespace is significant only when it separates two "word" tokens (identifier/number
// runs) — e.g. between a keyword and an identifier — and is dropped at the boundaries or
// next to punctuation/operators where MySQL never needs it. A backtick at expr[next]
// starts an identifier token, so it counts as a word boundary even though the backtick
// byte itself is not an identifier char (this is the fix for `a` and `b` → must stay
// `a and b`, not `a andb`).
func exprSpaceIsSignificant(emitted, expr string, next int) bool {
	if emitted == "" || next >= len(expr) {
		return false
	}
	prev := emitted[len(emitted)-1]
	nextStartsWord := isIdentByte(expr[next]) || expr[next] == '`'
	return isIdentByte(prev) && nextStartsWord
}

// ---------------------------------------------------------------------------
// Nullability + PK + the TIMESTAMP magic (explicit_defaults_for_timestamp).
// ---------------------------------------------------------------------------

// CanonicalNotNull returns the effective NOT NULL state of a column after MySQL's
// implicit rewrites:
//   - a PRIMARY KEY member is forced NOT NULL regardless of its declaration
//     (entry pk-implies-not-null);
//   - when explicit_defaults_for_timestamp is OFF, a TIMESTAMP column with NO explicit
//     nullability clause becomes NOT NULL (entry timestamp-magic-first-57). This is
//     governed by the captured session variable, NOT the version (flag #3): a 5.7
//     server with the variable ON leaves TIMESTAMP nullable, an 8.0 server with it OFF
//     forces NOT NULL. The Column.NullExplicit flag (set by the loader) tells a bare
//     `TIMESTAMP` apart from an explicit `TIMESTAMP NULL`, so a genuinely-nullable
//     `TIMESTAMP NULL` — even one with a constant default — is left nullable.
//
// For all other columns it returns the declared !Nullable.
func (n *Normalizer) CanonicalNotNull(table *Table, c *Column) bool {
	if !c.Nullable {
		return true
	}
	if n.columnInPrimaryKey(table, c) {
		return true
	}
	if isTimestampType(c.DataType) && !n.ExplicitDefaultsForTimestamp && !c.NullExplicit {
		return true
	}
	return false
}

// columnInPrimaryKey reports whether the column participates in the table's PK.
func (n *Normalizer) columnInPrimaryKey(table *Table, c *Column) bool {
	for _, idx := range table.Indexes {
		if !idx.Primary {
			continue
		}
		for _, ic := range idx.Columns {
			if strings.EqualFold(ic.Name, c.Name) {
				return true
			}
		}
	}
	// Constraints may also express PK.
	for _, con := range table.Constraints {
		if con.Type != ConPrimaryKey {
			continue
		}
		for _, col := range con.Columns {
			if strings.EqualFold(col, c.Name) {
				return true
			}
		}
	}
	return false
}

// CanonicalTimestampDefaults returns the canonical (default, onUpdate) keys for a
// TIMESTAMP/DATETIME column, applying the explicit_defaults_for_timestamp=OFF magic:
// the FIRST bare TIMESTAMP column in the table gets an implicit
// DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP. Both keys use the same
// value-based form as CanonicalDefault ("ts:CURRENT_TIMESTAMP", "<none>", etc.).
// When the user wrote an explicit default/on-update, that wins.
func (n *Normalizer) CanonicalTimestampDefaults(table *Table, c *Column) (def, onUpdate string) {
	def = n.CanonicalDefault(c)
	onUpdate = n.canonicalOnUpdate(c)

	if !isTimestampType(c.DataType) {
		return def, onUpdate
	}
	if n.ExplicitDefaultsForTimestamp {
		return def, onUpdate
	}
	// EDFT OFF: the FIRST TIMESTAMP column gets an implicit
	// DEFAULT CURRENT_TIMESTAMP(N) ON UPDATE CURRENT_TIMESTAMP(N) (matching its
	// fractional-seconds precision) UNLESS it has an explicit default or is explicitly
	// NULL. An explicit `TIMESTAMP NULL` suppresses the magic and reads back
	// `NULL DEFAULT NULL`; an explicit `TIMESTAMP NOT NULL` still receives the magic. We
	// key the condition off the canonical NOT NULL state (which already encodes the
	// explicit-NULL / EDFT rules), so both bare and `NOT NULL` first columns get it.
	if def == defaultAbsentKey && onUpdate == defaultAbsentKey &&
		n.isFirstTimestampColumn(table, c) && n.CanonicalNotNull(table, c) {
		ts := "ts:CURRENT_TIMESTAMP" + fractionalSecondsSuffix(c.ColumnType)
		return ts, ts
	}
	return def, onUpdate
}

// fractionalSecondsSuffix returns "(N)" for a TIMESTAMP(N)/DATETIME(N) column type, or
// "" when no fractional precision is declared, so the injected CURRENT_TIMESTAMP default
// carries the same precision MySQL stores (entry timestamp-magic-first-57: TIMESTAMP(3)
// → CURRENT_TIMESTAMP(3)).
func fractionalSecondsSuffix(columnType string) string {
	p := parseColumnType(columnType)
	if p.hasLen && p.length > 0 {
		return fmt.Sprintf("(%d)", p.length)
	}
	return ""
}

// canonicalOnUpdate returns the canonical ON UPDATE key for a column.
func (n *Normalizer) canonicalOnUpdate(c *Column) string {
	if c.OnUpdate == "" {
		return defaultAbsentKey
	}
	if ts, ok := normalizeCurrentTimestamp(c.OnUpdate); ok {
		return "ts:" + ts
	}
	return "str:" + c.OnUpdate
}

// isFirstTimestampColumn reports whether c is the first TIMESTAMP column (by position)
// in the table — the one that receives the EDFT-off implicit default.
func (n *Normalizer) isFirstTimestampColumn(table *Table, c *Column) bool {
	for _, col := range table.Columns {
		if isTimestampType(col.DataType) {
			return strings.EqualFold(col.Name, c.Name)
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// Table options — engine default + ignore-in-diff (AUTO_INCREMENT, ROW_FORMAT).
// ---------------------------------------------------------------------------

// defaultEngine is the server default storage engine that MySQL always emits in
// SHOW CREATE even when the user omits ENGINE.
const defaultEngine = "innodb"

// CanonicalEngine returns the table's effective storage engine, lowercased, resolving
// an empty value to the server default (InnoDB) so an omitted ENGINE matches the
// always-emitted ENGINE=InnoDB readback (entry engine-always-emitted).
func (n *Normalizer) CanonicalEngine(table *Table) string {
	e := strings.ToLower(strings.TrimSpace(table.Engine))
	if e == "" {
		return defaultEngine
	}
	return e
}

// IgnoreTableAutoIncrement always reports true: the table-level AUTO_INCREMENT=N option
// is the live next-value counter, not schema, and changes with every insert. diff-core
// must never compare it (entry auto-increment-counter). The column-level AUTO_INCREMENT
// attribute is real schema and is compared separately.
func (n *Normalizer) IgnoreTableAutoIncrement() bool { return true }

// IgnoreRowFormat reports whether the table's ROW_FORMAT should be excluded from the
// diff. The default ROW_FORMAT is environment-dependent (innodb_default_row_format) and
// not reliably reconstructible from DDL, so it is ignored unless the user explicitly set
// a non-default value (entry row-format-default-omitted; flag row-format-default-source).
func (n *Normalizer) IgnoreRowFormat(table *Table) bool {
	rf := strings.ToUpper(strings.TrimSpace(table.RowFormat))
	return rf == "" || rf == "DEFAULT"
}

// ---------------------------------------------------------------------------
// Comments — content comparison (decoded).
// ---------------------------------------------------------------------------

// CanonicalComment returns a comment's content for comparison. The omni loader already
// stores the DECODED comment content (embedded single quotes un-doubled), so the stored
// value IS the canonical content and is returned as-is. (entry comment-escaping: the diff
// compares decoded content, not the quoted surface form.) Re-running an un-doubling pass
// here would be wrong — it would collapse a genuine `a”b` content down to `a'b`.
func (n *Normalizer) CanonicalComment(comment string) string {
	return comment
}

// ---------------------------------------------------------------------------
// Index columns — descending direction is version-flagged.
// ---------------------------------------------------------------------------

// CanonicalIndexColumn returns a comparison key for one index column part, including
// its prefix length (entry prefix-index) and — only on 8.0 — its DESC direction. 5.7
// has no descending indexes and silently stores ascending, so direction is dropped from
// the 5.7 key (entry index-desc-asc). Explicit ASC is the default and never contributes.
func (n *Normalizer) CanonicalIndexColumn(ic *IndexColumn) string {
	var b strings.Builder
	if ic.Expr != "" {
		b.WriteString("expr:")
		b.WriteString(n.CanonicalGeneratedExpr(ic.Expr))
	} else {
		b.WriteString("col:")
		b.WriteString(strings.ToLower(ic.Name))
	}
	if ic.Length > 0 {
		fmt.Fprintf(&b, "(%d)", ic.Length)
	}
	// DESC is meaningful only on 8.0; 5.7 stores ascending regardless.
	if n.is80() && ic.Descending {
		b.WriteString(" desc")
	}
	return b.String()
}

// ---------------------------------------------------------------------------
// CanonicalColumn — the aggregate comparison key (the phantom-diff firewall).
// ---------------------------------------------------------------------------

// CanonicalColumn returns a single stable comparison key for a column, folding every
// per-aspect canonicalizer (type, resolved charset/collation pair, default, on-update,
// nullability, generated expression, auto-increment, invisibility) into one string.
// diff-core's column comparison (the MySQL analog of PG's columnsChanged at
// pg/catalog/diff_column.go:66) compares two columns by equality of this key: equal
// keys mean no change, so a column whose surface form differs from its synced readback
// but whose canonical form is identical produces no phantom diff.
//
// It is the primary entry point; the individual Canonical*/Resolve* methods are exposed
// for differs that need a single aspect (e.g. an ALTER that only touches the default).
func (n *Normalizer) CanonicalColumn(table *Table, c *Column) string {
	charset, collation := n.ResolveColumnCharsetCollation(table, c)
	def, onUpdate := n.CanonicalTimestampDefaults(table, c)

	gen, stored := "", ""
	if c.Generated != nil {
		gen = n.CanonicalGeneratedExpr(c.Generated.Expr)
		stored = strconv.FormatBool(c.Generated.Stored)
	}
	// Length-delimited field encoding: each value is prefixed with its byte length, so a
	// value containing the field separator (e.g. a comment of ";gen=...") cannot be
	// mistaken for another field and collide with a different column's key.
	return encodeKeyFields(
		"name", strings.ToLower(c.Name),
		"type", n.CanonicalColumnType(c),
		"cs", charset,
		"coll", collation,
		"notnull", strconv.FormatBool(n.CanonicalNotNull(table, c)),
		"default", def,
		"onupdate", onUpdate,
		"autoinc", strconv.FormatBool(c.AutoIncrement),
		"comment", n.CanonicalComment(c.Comment),
		"gen", gen,
		"stored", stored,
		"invisible", strconv.FormatBool(c.Invisible),
		"srid", strconv.Itoa(c.SRID),
	)
}

// encodeKeyFields joins (label, value) pairs into an unambiguous comparison key. Each
// value is written as "<label>:<len>:<value>" so no value byte sequence can be confused
// with a delimiter or another field — making the aggregate key collision-free regardless
// of value content (comments, defaults, expressions).
func encodeKeyFields(pairs ...string) string {
	var b strings.Builder
	for i := 0; i+1 < len(pairs); i += 2 {
		b.WriteString(pairs[i])
		b.WriteByte(':')
		b.WriteString(strconv.Itoa(len(pairs[i+1])))
		b.WriteByte(':')
		b.WriteString(pairs[i+1])
		b.WriteByte('|')
	}
	return b.String()
}

// ---------------------------------------------------------------------------
// CHECK constraints — version-flagged (8.0 normalized + auto-named; 5.7 dropped).
// ---------------------------------------------------------------------------

// CheckSupported reports whether CHECK constraints are represented in the stored form
// for the target version. 8.0 keeps and normalizes CHECK; 5.7 parses and silently drops
// it, so any target CHECK is unrepresentable and cannot round-trip (entry
// check-constraint). diff-core must skip CHECK on 5.7 (treat as unsupported) and compare
// the normalized expression on 8.0.
func (n *Normalizer) CheckSupported() bool { return n.is80() }

// CanonicalCheckExpr returns the canonical key for a CHECK constraint expression on 8.0,
// using the same expression normalizer as generated columns (entry check-constraint:
// "(a > 0)" -> "((`a` > 0))" normalized). On 5.7 callers should consult CheckSupported
// first; this still returns a normalized key for completeness.
func (n *Normalizer) CanonicalCheckExpr(expr string) string {
	return n.CanonicalGeneratedExpr(expr)
}
