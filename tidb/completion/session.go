package completion

import (
	"context"

	"github.com/bytebase/omni/internal/mysqlcomplete"
	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/catalog"
	"github.com/bytebase/omni/tidb/parser"
)

// SessionConfig configures multi-database, session-aware SQL completion.
type SessionConfig = mysqlcomplete.SessionConfig

// CompleteSession runs session-aware auto-completion over statement at the
// 1-based line and 0-based UTF-16 code-unit caretOffset.
func CompleteSession(ctx context.Context, cfg SessionConfig, statement string, caretLine, caretOffset int) []Candidate {
	raw := mysqlcomplete.CompleteSession(ctx, cfg, mysqlcomplete.Dialect{
		CurrentStatement:  currentStatement,
		IsReservedKeyword: parser.IsReservedKeyword,
		HasDatabaseRule:   hasDatabaseRule,
		NewCatalog: func(ctx context.Context) mysqlcomplete.CatalogHooks {
			cat := catalog.New()
			return mysqlcomplete.CatalogHooks{
				RegisterDB: func(name string) {
					_ = cat.DefineDatabase(&ast.CreateDatabaseStmt{Name: name, IfNotExists: true})
				},
				LoadDB: func(name string, meta *metadatapb.DatabaseSchemaMetadata) {
					_, _ = cat.LoadMetadataInto(ctx, meta, name)
				},
				SetCurrentDB: cat.SetCurrentDatabase,
				Complete: func(stmt string, pos int) []mysqlcomplete.Candidate {
					cands := Complete(stmt, pos, cat)
					out := make([]mysqlcomplete.Candidate, len(cands))
					for i, c := range cands {
						out[i] = mysqlcomplete.Candidate{
							Text:       c.Text,
							Type:       mysqlcomplete.CandidateType(c.Type),
							Definition: c.Definition,
							Comment:    c.Comment,
						}
					}
					return out
				},
			}
		},
	}, statement, caretLine, caretOffset)

	out := make([]Candidate, len(raw))
	for i, c := range raw {
		out[i] = Candidate{
			Text:       c.Text,
			Type:       CandidateType(c.Type),
			Definition: c.Definition,
			Comment:    c.Comment,
		}
	}
	return out
}

// QuoteIdentifierIfNeeded backtick-quotes name when it would otherwise be
// invalid as a bare identifier or is a reserved keyword. When caretInBacktick
// is true, name is returned unquoted.
func QuoteIdentifierIfNeeded(name string, caretInBacktick bool) string {
	return mysqlcomplete.QuoteIdentifierIfNeeded(name, caretInBacktick, parser.IsReservedKeyword)
}

// CaretInsideBacktickIdentifier reports whether byte offset pos in statement
// sits inside an open backtick-quoted identifier.
func CaretInsideBacktickIdentifier(statement string, pos int) bool {
	return mysqlcomplete.CaretInsideBacktickIdentifier(statement, pos)
}

// StatementReferencesDatabase reports whether dbName appears as a database
// qualifier (`dbName.`) in statement, at an identifier boundary.
func StatementReferencesDatabase(statement, dbName string) bool {
	return mysqlcomplete.StatementReferencesDatabase(statement, dbName)
}

// LineOffsetToBytePos converts a 1-based line and 0-based UTF-16 code-unit
// column offset into a byte offset into s.
func LineOffsetToBytePos(s string, line, offset int) int {
	return mysqlcomplete.LineOffsetToBytePos(s, line, offset)
}

func hasDatabaseRule(sql string, offset int) bool {
	cs := parser.Collect(sql, offset)
	return cs != nil && cs.HasRule("database_ref")
}

func currentStatement(statement string, pos int) (string, int) {
	for _, seg := range parser.Split(statement) {
		if pos >= seg.ByteStart && pos <= seg.ByteEnd {
			return seg.Text, pos - seg.ByteStart
		}
	}
	return "", 0
}
