package v1review

import (
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// SystemCharsetAllowlist judges the encoding each top-level CREATE DATABASE
// names against the allowed list, compared without regard to letter case.
// Only the first ENCODING option written as a string is read; a numeric or
// DEFAULT encoding is skipped, as is a statement under EXPLAIN or in a
// routine body. Every other statement, ALTER DATABASE included, is not
// judged.
//
// Name is the encoding lower-cased. Every finding anchors the whole
// statement. Ported from Bytebase's V1 rule SYSTEM_CHARSET_ALLOWLIST.
func SystemCharsetAllowlist(stmts []ast.Node, allowed []string) []v1review.Finding {
	allowedSet := make(map[string]bool, len(allowed))
	for _, encoding := range allowed {
		allowedSet[strings.ToLower(encoding)] = true
	}
	var findings []v1review.Finding
	for _, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		n, ok := stmt.(*ast.CreatedbStmt)
		if !ok {
			continue
		}
		encoding := strings.ToLower(createdbEncoding(n))
		if encoding == "" || allowedSet[encoding] {
			continue
		}
		start, end := span(n.Loc)
		findings = append(findings, v1review.Finding{Code: v1review.CharsetNotAllowed, Name: encoding, Start: start, End: end})
	}
	return findings
}

// createdbEncoding returns the first ENCODING option of n written as a
// string, "" when there is none.
func createdbEncoding(n *ast.CreatedbStmt) string {
	if n.Options == nil {
		return ""
	}
	for _, item := range n.Options.Items {
		def, ok := item.(*ast.DefElem)
		if !ok || !strings.EqualFold(def.Defname, "encoding") {
			continue
		}
		if s, ok := def.Arg.(*ast.String); ok {
			return s.Str
		}
	}
	return ""
}
