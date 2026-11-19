package review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
)

// dropKind is what an ALTER TABLE ... DROP subcommand drops.
type dropKind int

const (
	dropNone dropKind = iota
	dropColumn
	dropIndex
	dropPrimaryKey
	dropForeignKey
	dropCheck
	// dropConstraint is DROP CONSTRAINT, which drops the foreign key,
	// unique key, or check constraint of the name.
	dropConstraint
	dropPartition
)

// dropKindOf returns what an ALTER TABLE subcommand drops, dropNone for
// one that drops nothing. The parser gives DROP PRIMARY KEY, DROP FOREIGN
// KEY, DROP CHECK, and DROP CONSTRAINT one command type and keeps the keyword
// each of them read on the command, which tells them apart.
func dropKindOf(cmd *ast.AlterTableCmd) dropKind {
	switch cmd.Type {
	case ast.ATDropColumn:
		return dropColumn
	case ast.ATDropIndex:
		// The primary key is the index named PRIMARY.
		if sameName(cmd.Name, "PRIMARY") {
			return dropPrimaryKey
		}
		return dropIndex
	case ast.ATDropPartition:
		return dropPartition
	case ast.ATDropConstraint:
		switch cmd.ConstraintKeyword {
		case ast.ConstraintKeywordPrimaryKey:
			return dropPrimaryKey
		case ast.ConstraintKeywordForeignKey:
			return dropForeignKey
		case ast.ConstraintKeywordCheck:
			return dropCheck
		case ast.ConstraintKeywordConstraint:
			return dropConstraint
		}
	}
	return dropNone
}

func tokenText(text string, tok parser.Token) string {
	if tok.Loc < 0 || tok.End > len(text) || tok.End <= tok.Loc {
		return ""
	}
	return text[tok.Loc:tok.End]
}
