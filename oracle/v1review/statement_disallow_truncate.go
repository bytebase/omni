package v1review

import (
	"strings"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementDisallowTruncate reports what top-level statements truncate: the
// table a TRUNCATE TABLE names (TRUNCATE CLUSTER is not judged), and each
// ALTER TABLE ... TRUNCATE PARTITION or SUBPARTITION action. text is the
// statement text the nodes were parsed from; names are printed from it as
// written.
//
// Table is the table's source text, or its unqualified name when its range
// does not lie within text. For a partition action, Value is "PARTITION" or
// "SUBPARTITION" (the action's subtype compared case-insensitively) and Name
// is the action's source text after the first occurrence of that keyword,
// trimmed, with a leading "FOR " written "FOR"; when that is empty, the
// parsed partition name. A TRUNCATE TABLE finding anchors the statement; a
// partition finding anchors its action. Codes: Truncate, TruncatePartition.
// Ported from Bytebase's V1 rule STATEMENT_DISALLOW_TRUNCATE.
func StatementDisallowTruncate(stmts []ast.Node, text string) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.TruncateStmt:
			if n.Cluster {
				continue
			}
			name := statementObjectName(n.Table)
			if n.Table != nil {
				if raw := statementSourceText(text, n.Table.Loc); raw != "" {
					name = raw
				}
			}
			start, end := anchor(n.Loc)
			findings = append(findings, v1review.Finding{Code: v1review.Truncate, Table: name, Start: start, End: end})
		case *ast.AlterTableStmt:
			table := statementObjectName(n.Name)
			if n.Name != nil {
				if raw := statementSourceText(text, n.Name.Loc); raw != "" {
					table = raw
				}
			}
			if n.Actions == nil {
				continue
			}
			for _, item := range n.Actions.Items {
				cmd, ok := item.(*ast.AlterTableCmd)
				if !ok || cmd.Action != ast.AT_TRUNCATE_PARTITION {
					continue
				}
				keyword := "PARTITION"
				if strings.EqualFold(cmd.Subtype, "SUBPARTITION") {
					keyword = "SUBPARTITION"
				}
				target := cmd.ColumnName
				if target == "" {
					target = cmd.NewName
				}
				if rawAction := statementSourceText(text, cmd.Loc); rawAction != "" {
					if idx := strings.Index(strings.ToUpper(rawAction), keyword); idx >= 0 {
						if rawTarget := strings.TrimSpace(rawAction[idx+len(keyword):]); rawTarget != "" {
							if strings.HasPrefix(strings.ToUpper(rawTarget), "FOR ") {
								rawTarget = "FOR" + strings.TrimSpace(rawTarget[len("FOR "):])
							}
							target = rawTarget
						}
					}
				}
				start, end := anchor(cmd.Loc)
				findings = append(findings, v1review.Finding{Code: v1review.TruncatePartition, Name: target, Table: table, Value: keyword, Start: start, End: end})
			}
		}
	}
	return findings
}

// statementObjectName returns an object's unqualified name, empty for none.
func statementObjectName(name *ast.ObjectName) string {
	if name == nil {
		return ""
	}
	return name.Name
}

// statementSourceText returns the text loc spans, empty when loc does not lie within
// text.
func statementSourceText(text string, loc ast.Loc) string {
	if loc.Start < 0 || loc.End < loc.Start || text == "" || loc.End > len(text) {
		return ""
	}
	return text[loc.Start:loc.End]
}
