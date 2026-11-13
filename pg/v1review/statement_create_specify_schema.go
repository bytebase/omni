package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementCreateSpecifySchema reports each top-level CREATE TABLE whose
// table has no schema, one SchemaUnspecified finding with Name "Table".
// That covers partition and typed tables; CREATE TABLE AS, views,
// sequences, indexes and every other object are not judged. Findings
// anchor nothing: Bytebase's V1 advisor placed them on the statement's
// first line. Ported from Bytebase's V1 rule
// STATEMENT_CREATE_SPECIFY_SCHEMA.
func StatementCreateSpecifySchema(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		create, ok := unwrapRaw(stmt).(*ast.CreateStmt)
		if !ok {
			continue
		}
		if create.Relation != nil && create.Relation.Schemaname == "" {
			findings = append(findings, v1review.Finding{Code: v1review.SchemaUnspecified, Name: "Table"})
		}
	}
	return findings
}
