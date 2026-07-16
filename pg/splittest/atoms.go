// Package splittest is the PG splitter test framework: enumeration
// fences, a constructive generator whose cases carry their own expected
// boundaries, and invariant checks that need no ground truth at all.
//
// Design: docs 2026-07-16 "PG Splitter 测试框架设计" (three-layer tower).
// The generator's key property is constructive truth: every atom
// guarantees that any semicolon inside its text is NOT a statement
// boundary, so when the composer joins statements with top-level
// semicolons, the expected segmentation is derived during construction
// and holds at any scale without consulting an engine.
package splittest

import (
	"fmt"
	"math/rand"
	"strings"
)

// Atom is a fragment of statement text whose internal semicolons (if
// any) are guaranteed by construction to be inside a quoting/comment
// construct and therefore must never split. Atoms never end in a state
// that leaks into following text (all constructs are closed), and are
// composed with single-space separators so adjacent atoms cannot fuse
// into different tokens (e.g. an identifier tail absorbing an E prefix).
type Atom struct {
	// Class is the construct class from the framework design (C1..C19).
	Class string
	// Gen produces the fragment text. It must return text that is
	// self-contained: every quote/comment/paren opened is closed.
	Gen func(r *rand.Rand) string
	// RequiresFix names an unfixed audit defect (e.g. "D2", "D5") that
	// must land before this atom can be enabled in the composer. Empty
	// means the construct is safe on current main.
	RequiresFix string
	// WholeStatement marks constructs that are only meaningful as a
	// complete statement (e.g. COPY FROM stdin, whose data mode
	// requires COPY as the statement's first word). The composer emits
	// them verbatim as full segments and never embeds them mid-statement.
	// Their Gen output must carry its own statement terminator.
	WholeStatement bool
}

// identChars are safe identifier-continuation characters for generated names.
const identChars = "abcdefghijklmnopqrstuvwxyz_"

func randIdent(r *rand.Rand, n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteByte(identChars[r.Intn(len(identChars))])
	}
	return b.String()
}

// payload returns dangerous inner content for quoted constructs:
// semicolons, quote-lookalikes, comment openers, dollar signs.
func payload(r *rand.Rand) string {
	pieces := []string{
		";", "; ", "a;b", "--;", "/*;", "$x$;", "$;$", "end;", "END ;",
		"select 1;", ");(", ";;",
	}
	var b strings.Builder
	n := 1 + r.Intn(3)
	for i := 0; i < n; i++ {
		b.WriteString(pieces[r.Intn(len(pieces))])
	}
	return b.String()
}

// Atoms is the construct-atom library. Each entry corresponds to a
// class from the enumeration layer; the composer draws from the subset
// whose RequiresFix is empty (or explicitly enabled after a fix lands).
var Atoms = []Atom{
	{
		Class: "C1-plain-string",
		Gen: func(r *rand.Rand) string {
			// '' doubling belongs inside; backslash is literal
			// (standard_conforming_strings=on).
			inner := strings.ReplaceAll(payload(r), "'", "''")
			if r.Intn(3) == 0 {
				inner += `\`
			}
			return "'" + inner + "'"
		},
	},
	{
		Class: "C2-escape-string",
		Gen: func(r *rand.Rand) string {
			// E-string: backslash escapes; embed \' and '' forms
			// around semicolons (audit D1, fixed in #374).
			forms := []string{`\';`, `'';`, `\\`, `\n;`}
			var b strings.Builder
			prefix := "E"
			if r.Intn(2) == 0 {
				prefix = "e"
			}
			b.WriteString(prefix + "'")
			n := 1 + r.Intn(3)
			for i := 0; i < n; i++ {
				b.WriteString(forms[r.Intn(len(forms))])
			}
			b.WriteString("'")
			return b.String()
		},
	},
	{
		Class: "C5-dollar-quote",
		Gen: func(r *rand.Rand) string {
			tag := ""
			if r.Intn(2) == 0 {
				tag = randIdent(r, 1+r.Intn(4))
			}
			delim := "$" + tag + "$"
			inner := payload(r)
			// Similar-but-different tag inside is fair game.
			if tag != "" && r.Intn(2) == 0 {
				inner += "$" + tag[:len(tag)-1] + "x$"
			}
			// Constructive guarantee: the payload must not contain the
			// closing delimiter itself, or the atom is no longer sealed.
			inner = strings.ReplaceAll(inner, delim, "@")
			return delim + inner + delim
		},
	},
	{
		Class: "C3-unicode-string",
		Gen: func(r *rand.Rand) string {
			// U&'...' lexes with plain-string quote rules: '' doubling,
			// backslash NOT special at the string boundary (UESCAPE
			// processing happens after lexing). Audit-verified.
			inner := strings.ReplaceAll(payload(r), "'", "''")
			return "U&'" + inner + "'"
		},
	},
	{
		Class: "C4-bit-string",
		Gen: func(r *rand.Rand) string {
			// B''/X'' consume to the closing quote like plain strings;
			// the splitter does not validate content, so a semicolon
			// inside must stay sealed even though the value is invalid.
			prefix := []string{"B", "b", "X", "x"}[r.Intn(4)]
			return prefix + "'01;10'"
		},
	},
	{
		Class: "C7-line-comment",
		Gen: func(r *rand.Rand) string {
			// Line comment swallows a semicolon; newline ends it so
			// the atom stays self-contained.
			return "--" + strings.ReplaceAll(payload(r), "\n", " ") + "\n"
		},
	},
	{
		Class: "C8-block-comment",
		Gen: func(r *rand.Rand) string {
			// Neutralize both delimiters: a stray closer would end the
			// comment early, a stray opener would deepen nesting beyond
			// the closers we emit — either way the atom leaks.
			inner := strings.ReplaceAll(payload(r), "*/", "* /")
			inner = strings.ReplaceAll(inner, "/*", "/ *")
			depth := 1 + r.Intn(3) // PG block comments nest
			return strings.Repeat("/*", depth) + inner + strings.Repeat("*/", depth)
		},
	},
	{
		Class: "C9-quoted-ident",
		Gen: func(r *rand.Rand) string {
			inner := strings.ReplaceAll(payload(r), `"`, `""`)
			return `"` + inner + `"`
		},
	},
	{
		Class: "C11-begin-atomic",
		Gen: func(r *rand.Rand) string {
			// SQL-standard function body: internal semicolons must not
			// split. Keep body constructs simple (quoted payloads) —
			// keyword-boundary pathology (D4) is enumeration-layer
			// territory until the hardening lands.
			body := "SELECT " + "'" + strings.ReplaceAll(payload(r), "'", "''") + "'" + ";"
			if r.Intn(2) == 0 {
				body += " SELECT 1;"
			}
			return "CREATE FUNCTION " + randIdent(r, 4) + "() RETURNS int LANGUAGE sql BEGIN ATOMIC " + body + " END"
		},
	},
	{
		Class:       "C10-paren-semicolon",
		RequiresFix: "D2",
		Gen: func(r *rand.Rand) string {
			return "(SELECT 1; SELECT 2)"
		},
	},
	{
		Class:       "C6-ident-dollar",
		RequiresFix: "D3",
		Gen: func(r *rand.Rand) string {
			return randIdent(r, 2) + "$" + randIdent(r, 1) + "$" + randIdent(r, 1)
		},
	},
	{
		Class:          "C12-copy-stdin",
		RequiresFix:    "D5",
		WholeStatement: true,
		Gen: func(r *rand.Rand) string {
			// Complete segment: statement, its semicolon, data lines
			// (with sealed semicolons), terminator line and its newline
			// (contract: the terminator's newline belongs to this segment).
			rows := []string{"a;b", "1\t\\N", "c;d;e"}
			data := rows[r.Intn(len(rows))] + "\n" + rows[r.Intn(len(rows))]
			return "COPY " + randIdent(r, 3) + " FROM stdin;\n" + data + "\n\\.\n"
		},
	},
}

// EnabledAtoms returns the atoms usable by the composer: those with no
// pending fix, plus any whose defect id is in the enabled set (used to
// switch atoms on as D2/D3/D5 fixes merge).
func EnabledAtoms(enabledFixes ...string) []Atom {
	on := map[string]bool{}
	for _, f := range enabledFixes {
		on[f] = true
	}
	var out []Atom
	for _, a := range Atoms {
		if a.RequiresFix == "" || on[a.RequiresFix] {
			out = append(out, a)
		}
	}
	return out
}

// glue is plain SQL text between atoms within one statement. Single
// spaces prevent token fusion between adjacent atoms.
func glue(r *rand.Rand) string {
	pieces := []string{"SELECT", "1", "+", "col", ",", "FROM t WHERE a =", "AS x"}
	return pieces[r.Intn(len(pieces))]
}

// Statement composes 1..k atoms with glue into a statement containing
// no top-level semicolon by construction.
func Statement(r *rand.Rand, atoms []Atom) string {
	var b strings.Builder
	n := 1 + r.Intn(3)
	b.WriteString(glue(r))
	for i := 0; i < n; i++ {
		a := atoms[r.Intn(len(atoms))]
		for a.WholeStatement {
			a = atoms[r.Intn(len(atoms))]
		}
		b.WriteString(" ")
		b.WriteString(a.Gen(r))
		b.WriteString(" ")
		b.WriteString(glue(r))
	}
	return b.String()
}

// Script is a generated multi-statement script together with its
// expected segmentation, derived during construction.
type Script struct {
	SQL string
	// Want holds the expected segment texts, in order. A segment is a
	// statement plus its terminating semicolon; inter-statement trivia
	// belongs to the FOLLOWING segment (current splitter contract).
	Want []string
}

// trivia is inter-statement noise attached to the start of the next
// segment: whitespace and full-line comments.
func trivia(r *rand.Rand) string {
	switch r.Intn(4) {
	case 0:
		return " "
	case 1:
		return "\n"
	case 2:
		return " -- trailing note;\n"
	default:
		return ""
	}
}

// Compose builds a script of n statements. Every semicolon that acts
// as a boundary is placed by this function; every other semicolon is
// sealed inside an atom. That is the constructive-truth guarantee.
func Compose(r *rand.Rand, atoms []Atom, n int) Script {
	var whole []Atom
	for _, a := range atoms {
		if a.WholeStatement {
			whole = append(whole, a)
		}
	}
	var sql strings.Builder
	var want []string
	pending := "" // trivia belonging to the upcoming segment
	for i := 0; i < n; i++ {
		var stmt string
		if len(whole) > 0 && r.Intn(4) == 0 {
			// Whole-statement construct: emitted verbatim, carries its
			// own terminator (e.g. COPY block through the \. line).
			stmt = pending + whole[r.Intn(len(whole))].Gen(r)
		} else {
			stmt = pending + Statement(r, atoms) + ";"
		}
		sql.WriteString(stmt)
		want = append(want, stmt)
		pending = trivia(r)
	}
	if pending != "" && strings.TrimSpace(pending) == "" {
		// Trailing pure-whitespace trivia forms no segment in the
		// current contract only if empty; keep scripts closed instead.
		pending = ""
	}
	if pending != "" {
		sql.WriteString(pending)
		want = append(want, pending)
	}
	return Script{SQL: sql.String(), Want: want}
}

// String implements fmt.Stringer for failure dumps.
func (s Script) String() string {
	return fmt.Sprintf("script(%d segs): %q", len(s.Want), s.SQL)
}
