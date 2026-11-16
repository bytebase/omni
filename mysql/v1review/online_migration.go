package v1review

import (
	"fmt"
	"strconv"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// OnlineMigration reports, in statement order, each top-level ALTER TABLE that
// changes a table whose estimated row count reaches minRows: the changes a
// blocking ALTER may copy or lock the table for. An ALTER with no command is
// not reported, and neither is one whose table the snapshot does not hold.
//
// The database the advice prints is currentDatabase for an unqualified table
// and the ALTER's own qualifier otherwise; the table is found in the
// snapshot's first schema by name alone, as the V1 advisor looked it up, so an
// ALTER of another database's table of that name is judged against this
// database's row count. A finding carries the ALTER's range, its row count in
// Value and its index in the statements in Statement; the V1 advisor's end
// position was the statement's last line of content, which the caller reads
// from the statement text. Two of V1's decisions stay in Bytebase with the
// advice: the database the online migration tool writes to, and the advice code
// V1 chose by how many statements the review holds. Ported from Bytebase's V1
// rule ADVICE_ONLINE_MIGRATION.
func OnlineMigration(stmts []ast.Node, schema v1review.Schema, currentDatabase string, minRows int64) []v1review.Finding {
	schemas := schema.Metadata.GetSchemas()
	if len(schemas) == 0 {
		return nil
	}
	var findings []v1review.Finding
	for i, stmt := range stmts {
		alter, ok := stmt.(*ast.AlterTableStmt)
		if !ok || alter.Table == nil || len(alter.Commands) == 0 {
			continue
		}
		database := currentDatabase
		if alter.Table.Schema != "" {
			database = alter.Table.Schema
		}
		for _, table := range schemas[0].GetTables() {
			if table.GetName() != alter.Table.Name {
				continue
			}
			if table.GetRowCount() >= minRows {
				start, end := span(alter.Loc)
				findings = append(findings, v1review.Finding{
					Code:      v1review.OnlineMigrationNeeded,
					Table:     fmt.Sprintf("%s.%s", database, alter.Table.Name),
					Value:     strconv.FormatInt(table.GetRowCount(), 10),
					Statement: i,
					Start:     start,
					End:       end,
				})
			}
			break
		}
	}
	return findings
}
