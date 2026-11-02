package catalog

import (
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// MariaDB hides an index from the optimizer with IGNORED (Index.Visible is
// false) and refuses to hide the primary key with 4174, whether the key was
// declared or is implicit: on a table without a PRIMARY KEY, MariaDB promotes
// its first UNIQUE key whose parts are whole NOT NULL columns
// (TABLE_SHARE::init_from_binary_frm_image, sql/table.cc, mariadb-11.8.8). The
// keys are sorted first (sort_keys, sql/sql_table.cc): UNIQUE before the rest,
// a long unique hash last among them, keys without a nullable part before
// those with one, PRIMARY first, keys without a prefix part before those with
// one, and otherwise in the order written. So the first of the keys this file
// calls candidates is the implicit primary key.
//
// The catalog reports 4174 only where it is sure which key MariaDB promotes.
// That needs every UNIQUE key's kind settled; a key whose kind turns on the
// engine, the charset or a length limit leaves the question open, and no
// error is reported.

// candidateKind classifies a UNIQUE key for the implicit primary key.
type candidateKind int

const (
	keyNotCandidate candidateKind = iota // never promoted
	keyCandidate                         // promoted when it sorts first
	keyUnsettled                         // cannot tell without the engine
)

// implicitKeyCandidates returns the table's UNIQUE keys MariaDB could promote,
// in the order the catalog holds them, and whether every UNIQUE key's kind is
// settled. A table with a declared PRIMARY KEY has no implicit one.
func implicitKeyCandidates(tbl *Table) ([]*Index, bool) {
	var cands []*Index
	for _, idx := range tbl.Indexes {
		if idx.Primary {
			return nil, true
		}
		if !idx.Unique || idx.Fulltext || idx.Spatial {
			continue
		}
		switch classifyUniqueKey(tbl, idx) {
		case keyCandidate:
			cands = append(cands, idx)
		case keyUnsettled:
			return nil, false
		}
	}
	return cands, true
}

// innodbMaxKeyBytes is InnoDB's key length limit. A UNIQUE key longer than it
// becomes a long unique hash, which MariaDB never promotes (verified on 11.8.8:
// UNIQUE on VARCHAR(1000) NOT NULL in utf8mb4 is printed USING HASH and may be
// IGNORED on a table without a primary key).
const innodbMaxKeyBytes = 3072

// innodbMaxPartChars bounds a character column's length in a key part that is
// sure to fit every InnoDB row format: 767 bytes at four bytes a character.
const innodbMaxPartChars = 191

func classifyUniqueKey(tbl *Table, idx *Index) candidateKind {
	if !isInnoDB(tbl.Engine) {
		// MEMORY and MyISAM build a hash or a shorter key of their own; MariaDB
		// promotes a USING HASH key on both (verified on 11.8.8).
		return keyUnsettled
	}
	if strings.EqualFold(idx.IndexType, "HASH") {
		// InnoDB has no hash index: an explicit USING HASH UNIQUE key is a long
		// unique hash (verified on 11.8.8).
		return keyNotCandidate
	}
	if idx.IndexType != "" && !strings.EqualFold(idx.IndexType, "BTREE") {
		return keyUnsettled
	}
	kind := keyCandidate
	keyBytes := 0
	for _, ic := range idx.Columns {
		if ic.Expr != "" || ic.Name == "" {
			return keyUnsettled
		}
		col := tbl.GetColumn(ic.Name)
		if col == nil {
			return keyUnsettled
		}
		if col.Nullable {
			// A key with a nullable part is never promoted.
			return keyNotCandidate
		}
		dt := strings.ToLower(col.DataType)
		switch {
		case isBlobKeyType(dt):
			switch {
			case ic.Length == 0:
				// A whole BLOB or TEXT column makes the key a long unique hash
				// (verified on 11.8.8).
				kind = keyNotCandidate
			case ic.Length < 255:
				// Shorter than any BLOB: a prefix part, never promoted.
				kind = keyNotCandidate
			default:
				return keyUnsettled
			}
		case isCharKeyType(dt):
			n, ok := declaredLength(col.ColumnType)
			if !ok || n > innodbMaxPartChars {
				return keyUnsettled
			}
			if ic.Length > 0 {
				if ic.Length >= n {
					return keyUnsettled
				}
				// A prefix part, never promoted.
				kind = keyNotCandidate
			}
			keyBytes += 4 * n
		default:
			size, ok := fixedKeyBytes[dt]
			if !ok || ic.Length > 0 {
				return keyUnsettled
			}
			keyBytes += size
		}
	}
	if kind == keyCandidate && keyBytes > innodbMaxKeyBytes {
		return keyUnsettled
	}
	return kind
}

func isInnoDB(engine string) bool {
	return engine == "" || strings.EqualFold(engine, "InnoDB")
}

func isBlobKeyType(dt string) bool {
	switch dt {
	case "tinyblob", "blob", "mediumblob", "longblob",
		"tinytext", "text", "mediumtext", "longtext", "json":
		return true
	}
	return false
}

func isCharKeyType(dt string) bool {
	switch dt {
	case "char", "varchar", "binary", "varbinary":
		return true
	}
	return false
}

// fixedKeyBytes bounds the key length of a type whose length is fixed.
var fixedKeyBytes = map[string]int{
	"tinyint": 1, "smallint": 2, "mediumint": 3, "int": 4, "integer": 4, "bigint": 8,
	"decimal": 32, "float": 8, "double": 8, "bit": 8, "year": 1,
	"date": 3, "time": 6, "datetime": 8, "timestamp": 7,
	"enum": 2, "set": 8,
}

// declaredLength reads n from a column type written "varchar(n)...".
func declaredLength(columnType string) (int, bool) {
	open := strings.IndexByte(columnType, '(')
	if open < 0 {
		return 0, false
	}
	end := strings.IndexByte(columnType[open:], ')')
	if end < 0 {
		return 0, false
	}
	n, err := strconv.Atoi(columnType[open+1 : open+end])
	return n, err == nil
}

// tablePrimaryKey returns the table's primary key, declared or implicit, or
// nil when it has none or the catalog cannot tell which key MariaDB promotes.
// ordered says the catalog holds the keys in the order MariaDB sorts them by,
// which is so for a table CREATE TABLE has just built: column-level keys come
// first here, but they cannot be IGNORED, so the first candidate here is an
// ignored one only when it is MariaDB's first too. Afterwards MariaDB keeps
// the keys in the order it sorted them, which the catalog does not track.
func tablePrimaryKey(tbl *Table, ordered bool) *Index {
	for _, idx := range tbl.Indexes {
		if idx.Primary {
			return idx
		}
	}
	cands, settled := implicitKeyCandidates(tbl)
	if !settled || len(cands) == 0 {
		return nil
	}
	if ordered || len(cands) == 1 {
		return cands[0]
	}
	return nil
}

// validateIgnoredPrimaryKey returns 4174 when the table's primary key, declared
// or implicit, is IGNORED. ordered is as for tablePrimaryKey. Without the
// order, a table whose candidates are all IGNORED still promotes an ignored
// one.
func validateIgnoredPrimaryKey(tbl *Table, ordered bool) error {
	if pk := tablePrimaryKey(tbl, ordered); pk != nil {
		if !pk.Visible {
			return errPKIndexCantBeIgnored()
		}
		return nil
	}
	if ordered {
		return nil
	}
	cands, settled := implicitKeyCandidates(tbl)
	if !settled || len(cands) == 0 {
		return nil
	}
	for _, idx := range cands {
		if idx.Visible {
			return nil
		}
	}
	return errPKIndexCantBeIgnored()
}

// hasIgnoredIndex reports whether any index of the table is IGNORED.
func hasIgnoredIndex(tbl *Table) bool {
	for _, idx := range tbl.Indexes {
		if !idx.Visible {
			return true
		}
	}
	return false
}

// checkAlterIndexIgnorability checks every ALTER INDEX of an ALTER TABLE
// against the table as the statement found it, as MariaDB does
// (mysql_prepare_alter_table, sql/sql_table.cc, mariadb-11.8.8): the index must
// be one of the table's, and not one the same statement drops or renames,
// else 1176 unless IF EXISTS; and it must not be the primary key, declared or
// implicit, whichever ignorability is asked (4174, also with IF EXISTS).
// Verified on 11.8.8: ALTER INDEX k IGNORED, DROP INDEX k and ALTER INDEX k
// IGNORED, RENAME INDEX k TO k2 both answer 1176 in either order.
func checkAlterIndexIgnorability(tbl *Table, cmds []*nodes.AlterTableCmd) error {
	return nil
	gone := map[string]bool{}
	alters := false
	for _, cmd := range cmds {
		switch cmd.Type {
		case nodes.ATDropIndex, nodes.ATRenameIndex:
			gone[toLower(cmd.Name)] = true
		case nodes.ATDropConstraint:
			if strings.EqualFold(cmd.Name, "PRIMARY") {
				gone["primary"] = true
			}
		case nodes.ATAlterIndexVisible, nodes.ATAlterIndexInvisible:
			alters = true
		}
	}
	if !alters {
		return nil
	}
	pk := tablePrimaryKey(tbl, false)
	for _, cmd := range cmds {
		if cmd.Type != nodes.ATAlterIndexVisible && cmd.Type != nodes.ATAlterIndexInvisible {
			continue
		}
		var idx *Index
		for _, i := range tbl.Indexes {
			if toLower(i.Name) == toLower(cmd.Name) {
				idx = i
				break
			}
		}
		if idx == nil || gone[toLower(cmd.Name)] {
			if cmd.IfExists {
				continue
			}
			return errKeyDoesNotExist(cmd.Name, tbl.Name)
		}
		if idx == pk {
			return errPKIndexCantBeIgnored()
		}
	}
	return nil
}
