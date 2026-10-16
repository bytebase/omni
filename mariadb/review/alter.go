package review

import (
	"strings"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/mariadb/parser"
)

// dropKind is what an ALTER TABLE ... DROP subcommand drops.
type dropKind int

const (
	dropNone dropKind = iota
	dropColumn
	dropIndex
	dropPrimaryKey
	dropForeignKey
	// dropConstraint is DROP CONSTRAINT, which drops the foreign key,
	// unique key, or check constraint of the name.
	dropConstraint
	dropPartition
	// dropSystemVersioning is DROP SYSTEM VERSIONING, which drops the
	// table's history rows.
	dropSystemVersioning
	dropPeriod
)

// dropKindOf returns what an ALTER TABLE subcommand drops, dropNone for
// one that drops nothing. The parser gives DROP [CONSTRAINT] PRIMARY KEY,
// DROP FOREIGN KEY, and DROP CONSTRAINT one command type; the keywords
// after DROP in the subcommand's own text tell them apart. MariaDB has no
// DROP CHECK, which is MySQL's.
func dropKindOf(sql string, cmd *ast.AlterTableCmd) dropKind {
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
	case ast.ATDropSystemVersioning:
		return dropSystemVersioning
	case ast.ATDropPeriod:
		return dropPeriod
	case ast.ATDropConstraint:
		switch wordAfterDrop(sql, cmd.Loc, 0) {
		case "PRIMARY":
			return dropPrimaryKey
		case "FOREIGN":
			return dropForeignKey
		case "CONSTRAINT":
			// DROP CONSTRAINT PRIMARY KEY. A constraint named `PRIMARY`
			// keeps its quotes here, and is not the primary key.
			if wordAfterDrop(sql, cmd.Loc, 1) == "PRIMARY" {
				return dropPrimaryKey
			}
			return dropConstraint
		}
	}
	return dropNone
}

// wordAfterDrop returns word i, counted from 0, of the words that follow
// DROP in the text of a location, upper-cased, "" when the text does not
// start with DROP or has no such word.
func wordAfterDrop(sql string, loc ast.Loc, i int) string {
	if loc.Start < 0 || loc.End > len(sql) || loc.End <= loc.Start {
		return ""
	}
	text := sql[loc.Start:loc.End]
	tokens := parser.Tokenize(text)
	if len(tokens) < i+2 || !strings.EqualFold(tokenText(text, tokens[0]), "DROP") {
		return ""
	}
	return strings.ToUpper(tokenText(text, tokens[i+1]))
}

func tokenText(text string, tok parser.Token) string {
	if tok.Loc < 0 || tok.End > len(text) || tok.End <= tok.Loc {
		return ""
	}
	return text[tok.Loc:tok.End]
}
