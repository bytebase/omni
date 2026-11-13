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
// This package is entirely separate from package review, the SQL Review V2
// contract: neither imports the other, the two share no types, and nothing
// here joins the V2 standard rule set.
package v1review

import "regexp"

// Code names the condition a finding reports; each rule function documents
// the codes it uses.
type Code string

const (
	// NameMismatch: a name does not match the configured pattern.
	NameMismatch Code = "name-mismatch"
	// NameTooLong: a name exceeds the configured length.
	NameTooLong Code = "name-too-long"
)

// Finding is one violation a rule function reports, in statement order.
type Finding struct {
	// Code is the condition the finding reports.
	Code Code
	// Name is the judged identifier, as the engine's parser stores it,
	// without its qualifiers.
	Name string
	// Start and End are the byte range that anchors the finding in the
	// text the statements were parsed from: the judged name's range where
	// the AST keeps one, else its clause's. Both 0 means the AST kept no
	// range; the caller then anchors the statement.
	Start, End int
}

// NamingConfig configures a naming rule; Bytebase builds it from the V1
// NamingRulePayload. The zero value requires nothing.
type NamingConfig struct {
	// Pattern is an RE2 regular expression (Go regexp syntax) a name must
	// match, matched anywhere in the name unless anchored with ^ and $.
	// Empty requires nothing.
	Pattern string
	// MaxLength is the longest allowed name, counted in Unicode code
	// points. 0 means no limit.
	MaxLength int
}

// Compile returns the pattern compiled, nil when the config has none.
func (c NamingConfig) Compile() (*regexp.Regexp, error) {
	if c.Pattern == "" {
		return nil, nil
	}
	return regexp.Compile(c.Pattern)
}
