package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// The ALGORITHM and LOCK clauses of ALTER TABLE, CREATE INDEX and DROP INDEX
// on InnoDB. The server refuses a clause the statement's operations cannot run
// under after it has judged everything else about the statement: a missing
// column is 1054 whatever ALGORITHM says. Every case below was read on
// mysql:8.0.46 and mysql:8.4.13, which answer them alike. An operation the
// catalog does not classify is taken to run under any clause, so the catalog
// refuses only what MySQL certainly refuses.

// ddlLimit is what one operation cannot run under. reason is the one InnoDB
// gives (1846), "" for an operation it refuses without one (1845); rank
// orders the reasons, since MySQL answers with the strongest whatever the
// order the operations were written in: a type change over a lock an index
// needs, over a refusal with no reason.
type ddlLimit struct {
	notInstant    bool
	instantReason string
	copyOnly      bool   // not INPLACE, and so not LOCK=NONE either
	needsLock     bool   // INPLACE, but not LOCK=NONE
	reason        string // for INPLACE and LOCK=NONE
	rank          int
}

const (
	reasonTypeChange   = "Cannot change column type INPLACE"
	reasonRebuildType  = "Need to rebuild the table to change column type"
	reasonDropPrimary  = "Dropping a primary key is not allowed without also adding a new primary key"
	reasonAutoInc      = "Adding an auto-increment column requires a lock"
	reasonFulltextLock = "Fulltext index creation requires a lock"
	reasonCopyLock     = "COPY algorithm requires a lock"
)

var (
	limitNotInstant = ddlLimit{notInstant: true}
	limitTypeChange = ddlLimit{notInstant: true, instantReason: reasonRebuildType, copyOnly: true, reason: reasonTypeChange, rank: 3}
	limitDropPK     = ddlLimit{notInstant: true, instantReason: reasonDropPrimary, copyOnly: true, reason: reasonDropPrimary, rank: 2}
	limitAutoInc    = ddlLimit{notInstant: true, instantReason: reasonAutoInc, needsLock: true, reason: reasonAutoInc, rank: 1}
	limitFulltext   = ddlLimit{notInstant: true, instantReason: reasonFulltextLock, needsLock: true, reason: reasonFulltextLock, rank: 1}
)

// alterDDLLimits classifies the operations of an ALTER TABLE: old is the
// table before it and tbl the table it leaves. serverCopy reports an
// operation the server itself runs by copying the table alone: adding a CHECK
// constraint.
func alterDDLLimits(old, tbl *Table, commands []*nodes.AlterTableCmd) (limits []ddlLimit, serverCopy bool) {
	addsPrimary := false
	for _, cmd := range commands {
		if cmd != nil && cmd.Type == nodes.ATAddConstraint && cmd.Constraint != nil && cmd.Constraint.Type == nodes.ConstrPrimaryKey {
			addsPrimary = true
		}
	}
	for _, cmd := range commands {
		if cmd == nil {
			continue
		}
		for _, def := range alterColumnDefs(cmd) {
			for _, cc := range def.Constraints {
				if cc == nil {
					continue
				}
				switch cc.Type {
				case nodes.ColConstrCheck:
					serverCopy = true
				case nodes.ColConstrPrimaryKey, nodes.ColConstrUnique:
					limits = append(limits, limitNotInstant)
				}
			}
		}
		switch cmd.Type {
		case nodes.ATAddColumn:
			for _, def := range alterColumnDefs(cmd) {
				col := tbl.GetColumn(def.Name)
				switch {
				case col == nil:
				case col.AutoIncrement:
					limits = append(limits, limitAutoInc)
				case col.Generated != nil && col.Generated.Stored:
					limits = append(limits, limitNotInstant)
				}
			}
		case nodes.ATDropColumn:
			if columnIndexed(old, cmd.Name) {
				limits = append(limits, limitNotInstant)
			}
		case nodes.ATModifyColumn, nodes.ATChangeColumn:
			if cmd.Column == nil {
				continue
			}
			if limit, ok := redefinedColumnLimit(old, tbl, old.GetColumn(cmd.Name), tbl.GetColumn(cmd.Column.Name)); ok {
				limits = append(limits, limit)
			}
		case nodes.ATConvertCharset:
			if convertChangesColumns(old, tbl) {
				limits = append(limits, limitTypeChange)
			}
		case nodes.ATAddConstraint, nodes.ATAddIndex:
			con := cmd.Constraint
			if con == nil {
				continue
			}
			switch con.Type {
			case nodes.ConstrCheck:
				serverCopy = true
			case nodes.ConstrFulltextIndex:
				limits = append(limits, limitFulltext)
			case nodes.ConstrPrimaryKey, nodes.ConstrUnique, nodes.ConstrIndex:
				limits = append(limits, limitNotInstant)
			}
		case nodes.ATDropIndex, nodes.ATRenameIndex, nodes.ATAlterIndexVisible, nodes.ATAlterIndexInvisible:
			limits = append(limits, limitNotInstant)
		case nodes.ATDropConstraint:
			if cmd.ConstraintKeyword == nodes.ConstraintKeywordPrimaryKey {
				if addsPrimary {
					limits = append(limits, limitNotInstant)
				} else {
					limits = append(limits, limitDropPK)
				}
			}
		case nodes.ATTableOption:
			if cmd.Option != nil {
				switch toLower(cmd.Option.Name) {
				case "engine", "storage engine", "comment", "auto_increment":
					limits = append(limits, limitNotInstant)
				}
			}
		}
	}
	return limits, serverCopy
}

// redefinedColumnLimit classifies a MODIFY or CHANGE that turns column was
// into is. An ENUM or SET is left alone: appending a member runs INSTANT.
func redefinedColumnLimit(old, tbl *Table, was, is *Column) (ddlLimit, bool) {
	if was == nil || is == nil || isEnumSetType(was.DataType) || isEnumSetType(is.DataType) {
		return ddlLimit{}, false
	}
	// A character set is a string column's alone: CONVERT TO beside MODIFY
	// id INT changes nothing of id.
	sameCharset := !isStringType(was.DataType) && !isStringType(is.DataType) ||
		columnCharset(old, was) == columnCharset(tbl, is) && columnCollation(old, was) == columnCollation(tbl, is)
	if strings.EqualFold(was.ColumnType, is.ColumnType) && sameCharset {
		if was.Nullable != is.Nullable {
			return limitNotInstant, true
		}
		return ddlLimit{}, false
	}
	// A VARCHAR grows in place while its length prefix keeps its size.
	if was.DataType == "varchar" && is.DataType == "varchar" && sameCharset {
		mb := charsetMaxBytes(columnCharset(tbl, is))
		from, to := columnTypeLength(was.ColumnType), columnTypeLength(is.ColumnType)
		if to >= from && (from*mb <= 255) == (to*mb <= 255) {
			return limitNotInstant, true
		}
	}
	return limitTypeChange, true
}

func columnCharset(tbl *Table, col *Column) string {
	if col.Charset != "" {
		return toLower(col.Charset)
	}
	return toLower(tbl.Charset)
}

func columnCollation(tbl *Table, col *Column) string {
	if col.Collation != "" {
		return toLower(col.Collation)
	}
	return toLower(tbl.Collation)
}

// convertChangesColumns reports whether CONVERT TO CHARACTER SET changed a
// character column's character set.
func convertChangesColumns(old, tbl *Table) bool {
	for _, was := range old.Columns {
		is := tbl.GetColumn(was.Name)
		if is == nil || !isStringType(was.DataType) {
			continue
		}
		if columnCharset(old, was) != columnCharset(tbl, is) {
			return true
		}
	}
	return false
}

func columnIndexed(tbl *Table, name string) bool {
	for _, idx := range tbl.Indexes {
		for _, ic := range idx.Columns {
			if strings.EqualFold(ic.Name, name) {
				return true
			}
		}
	}
	return false
}

// checkDDLClauses holds a statement's operations to its ALGORITHM and LOCK.
func checkDDLClauses(algorithm, lock string, limits []ddlLimit, serverCopy bool) error {
	algorithm, lock = strings.ToUpper(strings.TrimSpace(algorithm)), strings.ToUpper(strings.TrimSpace(lock))
	if lock == "DEFAULT" {
		lock = ""
	}
	if algorithm == "DEFAULT" {
		algorithm = ""
	}
	if algorithm == "INSTANT" && lock != "" {
		return &Error{Code: 1221, SQLState: "HY000", Message: "Incorrect usage of ALGORITHM=INSTANT and LOCK=NONE/SHARED/EXCLUSIVE"}
	}
	if serverCopy {
		switch {
		case algorithm == "INPLACE":
			return errAlterNotSupported("ALGORITHM=INPLACE", "ALGORITHM=COPY")
		case algorithm == "INSTANT":
			return errAlterNotSupported("ALGORITHM=INSTANT", "ALGORITHM=COPY/INPLACE")
		case lock == "NONE":
			return errAlterNotSupportedReason("LOCK=NONE", reasonCopyLock, "LOCK=SHARED")
		}
	}
	if algorithm == "COPY" && lock == "NONE" {
		return errAlterNotSupportedReason("LOCK=NONE", reasonCopyLock, "LOCK=SHARED")
	}
	strongest := func(applies func(ddlLimit) bool) (ddlLimit, bool) {
		var best ddlLimit
		found := false
		for _, l := range limits {
			if applies(l) && (!found || l.rank > best.rank) {
				best, found = l, true
			}
		}
		return best, found
	}
	switch algorithm {
	case "INSTANT":
		if l, ok := strongest(func(l ddlLimit) bool { return l.notInstant }); ok {
			if l.instantReason == "" {
				return errAlterNotSupported("ALGORITHM=INSTANT", "ALGORITHM=COPY/INPLACE")
			}
			return errAlterNotSupportedReason("ALGORITHM=INSTANT", l.instantReason, "ALGORITHM=COPY/INPLACE")
		}
	case "INPLACE":
		if l, ok := strongest(func(l ddlLimit) bool { return l.copyOnly }); ok {
			return errAlterNotSupportedReason("ALGORITHM=INPLACE", l.reason, "ALGORITHM=COPY")
		}
	}
	if lock == "NONE" && algorithm != "INSTANT" {
		if l, ok := strongest(func(l ddlLimit) bool { return l.copyOnly || l.needsLock }); ok {
			return errAlterNotSupportedReason("LOCK=NONE", l.reason, "LOCK=SHARED")
		}
	}
	return nil
}

func errAlterNotSupported(clause, try string) error {
	return &Error{Code: 1845, SQLState: "0A000", Message: fmt.Sprintf("%s is not supported for this operation. Try %s.", clause, try)}
}

func errAlterNotSupportedReason(clause, reason, try string) error {
	return &Error{Code: 1846, SQLState: "0A000", Message: fmt.Sprintf("%s is not supported. Reason: %s. Try %s.", clause, reason, try)}
}

// alterClauses returns the ALGORITHM and LOCK an ALTER TABLE writes, the last
// of each.
func alterClauses(commands []*nodes.AlterTableCmd) (algorithm, lock string) {
	for _, cmd := range commands {
		if cmd == nil {
			continue
		}
		switch cmd.Type {
		case nodes.ATAlgorithm:
			algorithm = cmd.Name
		case nodes.ATLock:
			lock = cmd.Name
		}
	}
	return algorithm, lock
}
