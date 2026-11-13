// Package v1review is the contract between Bytebase SQL Review V1 and the
// engine v1review packages (mysql/v1review, pg/v1review, ...). V1 stays
// Bytebase's product surface — the rule registry, proto payloads, severity
// levels, and advice rendering — and the engine packages hold only the
// judgment ported from Bytebase's advisor: which parts of the SQL violate a
// rule, decided from the engine's AST. Bytebase converts a rule's payload
// into a config type here, calls the engine's rule function with the
// statements it already parsed, and renders the findings as advice.
//
// Each rule is one function per engine, such as NamingTable, taking the
// engine's statement nodes and the rule's config and returning findings. A
// config the engine cannot use, such as a Pattern that does not compile, is
// an error: the rule was misconfigured, not evaluated.
//
// A ported rule keeps exactly the judgment the V1 advisor made: the same
// statement forms, the same name text, the same unit of length, and the
// same anchor, defects included. A port is a move; a change to what V1
// reports is a separate commit, named as one.
//
// This package is entirely separate from package review, the SQL Review V2
// contract: neither imports the other, the two share no types, and nothing
// here joins the V2 standard rule set.
package v1review

import (
	"regexp"
	"strings"
)

// Code names the condition a finding reports; each rule function documents
// the codes it uses.
type Code string

const (
	// NameMismatch: a name does not match the configured pattern.
	NameMismatch Code = "name-mismatch"
	// NameTooLong: a name exceeds the configured length.
	NameTooLong Code = "name-too-long"
	// PatternInvalid: a template rule's pattern, expanded for one object,
	// does not compile, so that object's name could not be judged. Detail
	// holds the compile error.
	PatternInvalid Code = "pattern-invalid"
	// NameIsKeyword: a name is a keyword in the table the engine's V1
	// advisor checked names against.
	NameIsKeyword Code = "name-is-keyword"
	// NameCaseMismatch: a name is not in the configured letter case.
	NameCaseMismatch Code = "name-case-mismatch"
	// NameUnqualified: an object is named without its schema.
	NameUnqualified Code = "name-unqualified"
)

// Finding is one violation a rule function reports, in statement order.
type Finding struct {
	// Code is the condition the finding reports.
	Code Code
	// Name is the judged identifier, the text the V1 advisor matched and
	// printed.
	Name string
	// Table is the table the judged object belongs to, as the V1 advisor
	// printed it. Rules whose advice does not name a table leave it empty.
	Table string
	// Pattern is the pattern the name was held to. Template rules set it,
	// since theirs differs per object; other rules leave it empty.
	Pattern string
	// Detail is the compile error's text, set with PatternInvalid.
	Detail string
	// Search is the text the V1 advisor searched the statement for to find
	// the finding's line, which can differ from Name: the source spelling
	// where Name is lower-cased, say. It is set only on a finding the V1
	// advisor placed that way; the caller repeats the search, and uses
	// Start and End for every other finding.
	Search string
	// Start and End are the byte range that anchors the finding in the
	// text the statements were parsed from: the node the V1 advisor took
	// its line from, which each engine's rule function names. Where the V1
	// advisor took its line from the statement text instead, by searching
	// it for the name or reading where the statement starts, the range is
	// the statement's and the caller keeps that computation. Both 0 means
	// the AST kept no range.
	Start, End int
}

// NamingConfig configures a naming rule; Bytebase builds it from the V1
// NamingRulePayload. The zero value requires nothing.
type NamingConfig struct {
	// Pattern is an RE2 regular expression (Go regexp syntax) a name must
	// match, matched anywhere in the name unless anchored with ^ and $.
	// Empty requires nothing.
	Pattern string
	// MaxLength is the longest allowed name, counted in bytes of its UTF-8
	// text, as every V1 advisor counts it. 0 means no limit.
	MaxLength int
}

// Compile returns the pattern compiled, nil when the config has none.
func (c NamingConfig) Compile() (*regexp.Regexp, error) {
	if c.Pattern == "" {
		return nil, nil
	}
	return regexp.Compile(c.Pattern)
}

// CaseConfig configures an identifier case rule; Bytebase builds it from the
// V1 NamingCaseRulePayload.
type CaseConfig struct {
	// Upper requires upper case; false requires lower case.
	Upper bool
}

// The tokens a template rule's Format may hold. Which of them a rule gives
// values to is the rule's own: an index rule the table and its column list,
// a foreign key rule the four referencing and referenced names.
const (
	TableToken             = "{{table}}"
	ColumnListToken        = "{{column_list}}"
	ReferencingTableToken  = "{{referencing_table}}"
	ReferencingColumnToken = "{{referencing_column}}"
	ReferencedTableToken   = "{{referenced_table}}"
	ReferencedColumnToken  = "{{referenced_column}}"
)

// TemplateConfig configures a naming rule whose pattern depends on the object
// it names; Bytebase builds it from the V1 NamingRulePayload after checking
// that Format holds only tokens its rule allows.
type TemplateConfig struct {
	// Format is an RE2 regular expression holding tokens such as {{table}},
	// each replaced by the object's own value before the pattern is compiled.
	Format string
	// MaxLength is the longest allowed name, counted in bytes. 0 means no
	// limit.
	MaxLength int
}

var templateToken = regexp.MustCompile(`{{[^{}]+}}`)

// Expand returns Format compiled for one object: each token Format holds, in
// the order it writes them, is replaced everywhere by its value, and a token
// with no value is left as written. A value goes in as it is, not quoted, so
// one holding a regular expression metacharacter changes the pattern and may
// keep it from compiling; the V1 advisors do the same.
func (c TemplateConfig) Expand(values map[string]string) (*regexp.Regexp, error) {
	pattern := c.Format
	for _, token := range templateToken.FindAllString(c.Format, -1) {
		if value, ok := values[token]; ok {
			pattern = strings.ReplaceAll(pattern, token, value)
		}
	}
	return regexp.Compile(pattern)
}
