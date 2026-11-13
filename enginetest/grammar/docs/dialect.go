// Package docs converts the Syntax blocks of a vendor's SQL command
// reference to an ANTLR 4 grammar, the .g4 a docs-judged compat.tsv is
// generated from and judged by (§2.2 of
// docs/specs/2026-11-11-engine-compatibility-without-servers-design.md).
//
// The blocks are near-BNF:
//
//	ALTER WAREHOUSE [ IF EXISTS ] <name> SET TAG <tag_name> = '<tag_value>' [ , ... ]
//
// Uppercase words are keywords, <placeholders> are values the reader
// supplies, [ ] is optional, { a | b } is a choice, ... repeats, and a
// lowercase word names a fragment the page defines below with ::=. Other
// punctuation, ( ) included, is literal. An engine's converter is a main
// that calls Run with its Dialect: where its docs are, how a page's Syntax
// blocks are found, and what its placeholders stand for. A Dialect whose
// docs write the notation another way rewrites its blocks into this one in
// Blocks.
package docs

import "path"

// A Dialect is what one vendor's docs need beyond the shared notation.
type Dialect struct {
	// Grammar names the grammar and its file, <Grammar>.g4.
	Grammar string
	// Command is the converter's package, as the files it writes cite it.
	Command string
	// Testdata is the default directory of the grammar, docs-fixes.tsv and
	// pins.tsv, relative to the repo root.
	Testdata string
	// Header is the .g4's opening comment, each line starting with //.
	Header string
	// Site+Index is the command index; Site+path is a page's URL.
	Site, Index string
	// Pages returns the command pages the index lists, in its order, each
	// once, as paths relative to Site.
	Pages func(index string) []string
	// Slug names a page in docs-fixes.tsv and its rule. Nil is the path's
	// last element.
	Slug func(path string) string
	// RuleName turns a slug into the page's rule name. Nil is RuleName.
	RuleName func(slug string) string
	// Blocks returns the text of a page's Syntax blocks, in the notation.
	Blocks func(page string) []string
	// Prepare corrects a page's blocks in place after its docs-fixes.tsv
	// rows, for typos the docs repeat on many pages, and returns a note
	// for each kind it corrected. It may be nil.
	Prepare func(slug string, blocks []string) []string
	// Classify returns the token class a <placeholder> stands for.
	Classify func(placeholder string) string
	// StatementWords returns the words a statement form on a page can
	// start with.
	StatementWords func(slug string) map[string]bool
}

func (d *Dialect) slug(p string) string {
	if d.Slug != nil {
		return d.Slug(p)
	}
	return path.Base(p)
}

func (d *Dialect) ruleName(slug string) string {
	if d.RuleName != nil {
		return d.RuleName(slug)
	}
	return RuleName(slug)
}
