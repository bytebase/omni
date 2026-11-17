package analysis

import (
	"strings"

	"github.com/bytebase/omni/trino/parser"
)

// Diagnose returns statement's syntax diagnostics, leaving out the parser's
// "not yet supported" stubs for forms that are valid Trino.
//
// The omni Trino parser emits a "<NAME> statement parsing is not yet supported"
// stub from a handful of dispatch fallthroughs. Most of those fire only on
// genuinely malformed input — every valid CREATE/DROP/ALTER object keyword is
// implemented, so e.g. "CREATE TABLEE ..." (an unrecognized object) is a true
// syntax error (confirmed against the Trino 481 oracle) and must be kept. The
// only valid-Trino forms that still hit the stub are DESCRIBE INPUT / DESCRIBE
// OUTPUT (oracle: NOT_FOUND, i.e. they parse). Suppress only those so the editor
// does not show a false positive on valid SQL; everything else passes through.
// When omni implements DESCRIBE INPUT/OUTPUT the stub stops firing and this
// filter becomes a no-op.
func Diagnose(statement string) []parser.Diagnostic {
	diags := parser.Diagnose(statement)
	out := make([]parser.Diagnostic, 0, len(diags))
	for _, d := range diags {
		if isValidButUnimplementedStub(d.Msg) {
			continue
		}
		out = append(out, d)
	}
	return out
}

// isValidButUnimplementedStub reports whether msg is an omni "not yet supported"
// stub error for a statement form that is valid Trino (so it must not be flagged
// as a syntax error). Today that is only DESCRIBE INPUT / DESCRIBE OUTPUT; the
// CREATE/DROP/ALTER stubs fire exclusively on unrecognized object keywords,
// which are genuine syntax errors, so they are intentionally NOT matched here.
func isValidButUnimplementedStub(msg string) bool {
	return strings.HasPrefix(msg, "DESCRIBE ") &&
		strings.HasSuffix(msg, "statement parsing is not yet supported")
}
