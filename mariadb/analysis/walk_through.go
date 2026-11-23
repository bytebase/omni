package analysis

import (
	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/mariadb/parser"
)

// TableRef is an engine-neutral schema-qualified table or object reference.
type TableRef = analysis.TableRef

// RenamePair captures the source and destination tables of a rename statement.
type RenamePair = analysis.RenamePair

// WalkThroughStatementKind classifies a statement for DDL walk-through precheck.
type WalkThroughStatementKind = analysis.WalkThroughStatementKind

const (
	WalkThroughStatementOther          = analysis.WalkThroughStatementOther
	WalkThroughStatementCreateDatabase = analysis.WalkThroughStatementCreateDatabase
	WalkThroughStatementAlterDatabase  = analysis.WalkThroughStatementAlterDatabase
	WalkThroughStatementDropDatabase   = analysis.WalkThroughStatementDropDatabase
	WalkThroughStatementCreateTable    = analysis.WalkThroughStatementCreateTable
	WalkThroughStatementDropTable      = analysis.WalkThroughStatementDropTable
	WalkThroughStatementTableRefCheck  = analysis.WalkThroughStatementTableRefCheck
	WalkThroughStatementRenameTable    = analysis.WalkThroughStatementRenameTable
)

// WalkThroughStatement holds the fields extracted from a MariaDB AST statement
// for DDL walk-through precheck and cross-database rename target discovery.
type WalkThroughStatement = analysis.WalkThroughStatement

// ClassifyWalkThroughStatement inspects a parsed MariaDB AST node and extracts
// its engine-neutral walk-through classification and referenced objects.
func ClassifyWalkThroughStatement(node ast.Node) WalkThroughStatement {
	toRef := func(t *ast.TableRef) TableRef {
		if t == nil {
			return TableRef{}
		}
		return TableRef{Schema: t.Schema, Name: t.Name}
	}
	toRefs := func(tables []*ast.TableRef) []TableRef {
		refs := make([]TableRef, 0, len(tables))
		for _, t := range tables {
			if t != nil {
				refs = append(refs, toRef(t))
			}
		}
		return refs
	}

	var info WalkThroughStatement
	switch n := node.(type) {
	case *ast.CreateDatabaseStmt:
		if n == nil {
			return info
		}
		info.Kind = WalkThroughStatementCreateDatabase
		info.Database = n.Name
	case *ast.AlterDatabaseStmt:
		if n == nil {
			return info
		}
		info.Kind = WalkThroughStatementAlterDatabase
		info.Database = n.Name
	case *ast.DropDatabaseStmt:
		if n == nil {
			return info
		}
		info.Kind = WalkThroughStatementDropDatabase
		info.Database = n.Name
	case *ast.CreateTableStmt:
		if n == nil {
			return info
		}
		info.Kind = WalkThroughStatementCreateTable
		info.IfNotExists = n.IfNotExists
		info.IsCTAS = n.Select != nil
		if n.Table != nil {
			info.TableRefs = []TableRef{toRef(n.Table)}
		}
	case *ast.AlterTableStmt:
		if n == nil {
			return info
		}
		info.Kind = WalkThroughStatementTableRefCheck
		if n.Table != nil {
			info.TableRefs = []TableRef{toRef(n.Table)}
		}
		for _, cmd := range n.Commands {
			if cmd != nil && cmd.Type == ast.ATRenameTable && n.Table != nil && cmd.NewTable != nil {
				info.RenamePairs = append(info.RenamePairs, RenamePair{
					Old: toRef(n.Table),
					New: toRef(cmd.NewTable),
				})
			}
		}
	case *ast.DropTableStmt:
		if n == nil {
			return info
		}
		info.Kind = WalkThroughStatementDropTable
		info.TableRefs = toRefs(n.Tables)
	case *ast.CreateIndexStmt:
		if n == nil {
			return info
		}
		info.Kind = WalkThroughStatementTableRefCheck
		if n.Table != nil {
			info.TableRefs = []TableRef{toRef(n.Table)}
		}
	case *ast.DropIndexStmt:
		if n == nil {
			return info
		}
		info.Kind = WalkThroughStatementTableRefCheck
		if n.Table != nil {
			info.TableRefs = []TableRef{toRef(n.Table)}
		}
	case *ast.TruncateStmt:
		if n == nil {
			return info
		}
		info.Kind = WalkThroughStatementTableRefCheck
		info.TableRefs = toRefs(n.Tables)
	case *ast.RenameTableStmt:
		if n == nil {
			return info
		}
		info.Kind = WalkThroughStatementRenameTable
		for _, pair := range n.Pairs {
			if pair == nil {
				continue
			}
			oldRef := toRef(pair.Old)
			if pair.Old != nil {
				info.TableRefs = append(info.TableRefs, oldRef)
			}
			if pair.Old != nil && pair.New != nil {
				info.RenamePairs = append(info.RenamePairs, RenamePair{
					Old: oldRef,
					New: toRef(pair.New),
				})
			}
		}
	case *ast.CreateViewStmt:
		if n == nil {
			return info
		}
		info.Kind = WalkThroughStatementTableRefCheck
		if n.Name != nil {
			info.TableRefs = []TableRef{toRef(n.Name)}
		}
	case *ast.CreateSequenceStmt:
		if n == nil {
			return info
		}
		info.Kind = WalkThroughStatementTableRefCheck
		if n.Name != nil {
			info.TableRefs = []TableRef{toRef(n.Name)}
		}
	case *ast.AlterSequenceStmt:
		if n == nil {
			return info
		}
		info.Kind = WalkThroughStatementTableRefCheck
		if n.Name != nil {
			info.TableRefs = []TableRef{toRef(n.Name)}
		}
	case *ast.DropSequenceStmt:
		if n == nil {
			return info
		}
		info.Kind = WalkThroughStatementTableRefCheck
		info.TableRefs = toRefs(n.Sequences)
	default:
	}
	return info
}

var (
	ctasTokCreate    = parser.NewLexer("CREATE").NextToken().Type
	ctasTokOr        = parser.NewLexer("OR").NextToken().Type
	ctasTokReplace   = parser.NewLexer("REPLACE").NextToken().Type
	ctasTokTemporary = parser.NewLexer("TEMPORARY").NextToken().Type
	ctasTokTable     = parser.NewLexer("TABLE").NextToken().Type
	ctasTokSelect    = parser.NewLexer("SELECT").NextToken().Type
	ctasTokValues    = parser.NewLexer("VALUES").NextToken().Type
	ctasTokWith      = parser.NewLexer("WITH").NextToken().Type
)

// IsCreateTableAs reports whether statement is a CREATE [OR REPLACE] [TEMPORARY]
// TABLE ... query_expression (SELECT, VALUES, or parenthesized query) statement
// by scanning its MariaDB token stream. It provides a lexer-backed fallback for
// CTAS forms that do not produce a CreateTableStmt AST.
func IsCreateTableAs(statement string) bool {
	lex := parser.NewLexer(statement)
	if lex.NextToken().Type != ctasTokCreate {
		return false
	}
	tok := lex.NextToken()
	if tok.Type == ctasTokOr {
		if lex.NextToken().Type != ctasTokReplace {
			return false
		}
		tok = lex.NextToken()
	}
	if tok.Type == ctasTokTemporary {
		tok = lex.NextToken()
	}
	if tok.Type != ctasTokTable {
		return false
	}
	parenDepth := 0
	atParenStart := false
	for {
		tok = lex.NextToken()
		switch tok.Type {
		case 0, ';':
			return false
		case '(':
			if parenDepth == 0 {
				atParenStart = true
			}
			parenDepth++
		case ')':
			if parenDepth > 0 {
				parenDepth--
			}
			atParenStart = false
		case ctasTokSelect:
			return true
		case ctasTokValues:
			if parenDepth == 0 || atParenStart {
				return true
			}
			atParenStart = false
		case ctasTokWith:
			if atParenStart {
				return true
			}
			atParenStart = false
		default:
			atParenStart = false
		}
	}
}
