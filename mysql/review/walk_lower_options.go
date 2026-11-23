package review

import (
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// mysqlCharsets are the character sets MySQL 8.0 and 8.4 have (SHOW
// CHARACTER SET on mysql:8.0.46 and mysql:8.4.13), with utf8 for utf8mb3. A
// server has no others: character sets are compiled in.
var mysqlCharsets = map[string]bool{
	"armscii8": true, "ascii": true, "big5": true, "binary": true, "cp1250": true, "cp1251": true,
	"cp1256": true, "cp1257": true, "cp850": true, "cp852": true, "cp866": true, "cp932": true,
	"dec8": true, "eucjpms": true, "euckr": true, "gb18030": true, "gb2312": true, "gbk": true,
	"geostd8": true, "greek": true, "hebrew": true, "hp8": true, "keybcs2": true, "koi8r": true,
	"koi8u": true, "latin1": true, "latin2": true, "latin5": true, "latin7": true, "macce": true,
	"macroman": true, "sjis": true, "swe7": true, "tis620": true, "ucs2": true, "ujis": true,
	"utf16": true, "utf16le": true, "utf32": true, "utf8": true, "utf8mb3": true, "utf8mb4": true,
}

// normalCharset returns a character set's name as MySQL resolves it, utf8
// as utf8mb3, "" for a name that is no character set.
func normalCharset(name string) string {
	n := strings.ToLower(strings.Trim(strings.TrimSpace(name), "'\"`"))
	if !mysqlCharsets[n] {
		return ""
	}
	if n == "utf8" {
		return "utf8mb3"
	}
	return n
}

// collationCharset returns the character set a collation belongs to, which
// its name begins with (binary is binary's own), "" when no character set
// begins it. The longest name wins, so utf8mb4_bin is utf8mb4's and not
// utf8's.
func collationCharset(collation string) string {
	c := strings.ToLower(strings.Trim(strings.TrimSpace(collation), "'\"`"))
	if c == "binary" {
		return "binary"
	}
	best := ""
	for cs := range mysqlCharsets {
		if strings.HasPrefix(c, cs+"_") && len(cs) > len(best) {
			best = cs
		}
	}
	return normalCharset(best)
}

// mysqlEngines are the storage engines a MySQL 8.0 or 8.4 server may have,
// by every name it takes for them: the ones the pinned servers build in
// (SHOW ENGINES on mysql:8.0.46), the ones a plugin of MySQL, of Percona
// Server, or of a third party may add, and the aliases MySQL resolves. Any other name is 1286
// under the default sql_mode, which has NO_ENGINE_SUBSTITUTION.
var mysqlEngines = map[string]bool{
	"innodb": true, "innobase": true, "myisam": true, "memory": true, "heap": true, "csv": true,
	"archive": true, "blackhole": true, "mrg_myisam": true, "merge": true, "federated": true,
	"performance_schema": true, "ndbcluster": true, "ndb": true, "ndbinfo": true, "example": true,
	"temptable": true, "rocksdb": true, "tokudb": true, "mroonga": true, "sphinx": true, "spider": true,
	"connect": true, "columnstore": true, "aria": true, "default": true,
}

// optionName returns a table option's name without a DEFAULT before it and
// with CHARSET for CHARACTER SET.
func optionName(o *ast.TableOption) string {
	n := strings.ToUpper(strings.Join(strings.Fields(o.Name), " "))
	n = strings.TrimPrefix(n, "DEFAULT ")
	switch n {
	case "CHARACTER SET":
		return "CHARSET"
	case "STORAGE ENGINE":
		return "ENGINE"
	}
	return n
}

func optionValue(o *ast.TableOption) string {
	return strings.Trim(strings.TrimSpace(o.Value), "'\"`")
}

// unknownCharset is the text for a character set MySQL doesn't have
// (1115).
func unknownCharset(name string) string {
	return "Character set " + review.Code(name) + " doesn't exist."
}

// collationMismatch is the text for a collation of another character set
// than the one beside it (1253).
func collationMismatch(collation, charset string) string {
	return "Collation " + review.Code(collation) + " isn't one of character set " + review.Code(charset) + "."
}

// optionTextFails fails a CREATE TABLE or ALTER TABLE whose table options,
// CONVERT TO, or column types MySQL refuses while it parses them, before it
// looks the table up (each verified on mysql:8.0.46 against a table that
// doesn't exist):
//
//   - a character set MySQL doesn't have (1115), in a table option, CONVERT
//     TO, or a column type;
//   - a storage engine no server has (1286);
//   - a collation of another character set than the one written beside it
//     (1253), in a column type, CONVERT TO, or a table's CHARSET and then
//     COLLATE; or COLLATE and then a CHARSET of another set (1302).
func (b *builder) optionTextFails(opts []*ast.TableOption, convert *ast.AlterTableCmd, cols []*ast.ColumnDef) bool {
	rng := b.anchor()
	for _, o := range opts {
		if o != nil && optionName(o) == "CHARSET" && !strings.EqualFold(optionValue(o), "DEFAULT") && normalCharset(o.Value) == "" {
			b.fail(rng, nil, review.UnsupportedForm, unknownCharset(optionValue(o)))
			return true
		}
	}
	if convert != nil && !strings.EqualFold(convert.Name, "DEFAULT") && normalCharset(convert.Name) == "" {
		b.fail(rng, nil, review.UnsupportedForm, unknownCharset(convert.Name))
		return true
	}
	for _, c := range cols {
		if c != nil && c.TypeName != nil && c.TypeName.Charset != "" && normalCharset(c.TypeName.Charset) == "" {
			b.fail(rangeOf(c.Loc), nil, review.UnsupportedForm, unknownCharset(c.TypeName.Charset))
			return true
		}
	}
	if !b.w.sqlModeSet {
		for _, o := range opts {
			if o != nil && optionName(o) == "ENGINE" && !mysqlEngines[strings.ToLower(optionValue(o))] {
				b.fail(rng, nil, review.UnsupportedForm, engineName+" has no storage engine "+review.Code(optionValue(o))+".")
				return true
			}
		}
	}
	for _, c := range cols {
		if c == nil || c.TypeName == nil || c.TypeName.Charset == "" || c.TypeName.Collate == "" {
			continue
		}
		cs, ccs := normalCharset(c.TypeName.Charset), collationCharset(c.TypeName.Collate)
		if cs != "" && ccs != "" && cs != ccs {
			b.fail(rangeOf(c.Loc), nil, review.UnsupportedForm, collationMismatch(c.TypeName.Collate, c.TypeName.Charset))
			return true
		}
	}
	if convert != nil && convert.NewName != "" {
		cs, ccs := normalCharset(convert.Name), collationCharset(convert.NewName)
		if cs != "" && ccs != "" && cs != ccs {
			b.fail(rng, nil, review.UnsupportedForm, collationMismatch(convert.NewName, convert.Name))
			return true
		}
	}
	var charset, collation *ast.TableOption
	for _, o := range opts {
		if o == nil {
			continue
		}
		switch optionName(o) {
		case "CHARSET":
			if charset != nil || strings.EqualFold(optionValue(o), "DEFAULT") {
				return false
			}
			charset = o
		case "COLLATE":
			if collation != nil || strings.EqualFold(optionValue(o), "DEFAULT") {
				return false
			}
			collation = o
		}
	}
	if charset != nil && collation != nil {
		cs, ccs := normalCharset(charset.Value), collationCharset(collation.Value)
		if cs != "" && ccs != "" && cs != ccs {
			b.fail(rng, nil, review.UnsupportedForm, collationMismatch(optionValue(collation), optionValue(charset)))
			return true
		}
	}
	return false
}

// alterOptionTextFails is optionTextFails for an ALTER TABLE.
func (b *builder) alterOptionTextFails(v *ast.AlterTableStmt) bool {
	var opts []*ast.TableOption
	var convert *ast.AlterTableCmd
	var cols []*ast.ColumnDef
	for _, c := range v.Commands {
		if c == nil {
			continue
		}
		switch c.Type {
		case ast.ATTableOption:
			if c.Option != nil {
				opts = append(opts, c.Option)
			}
		case ast.ATConvertCharset:
			convert = c
		case ast.ATAddColumn, ast.ATModifyColumn, ast.ATChangeColumn:
			if c.Column != nil {
				cols = append(cols, c.Column)
			}
			cols = append(cols, c.Columns...)
		}
	}
	return b.optionTextFails(opts, convert, cols)
}

// The secondary engine space holds one key per table, secondaryName,
// Present when the table has a secondary engine (SECONDARY_ENGINE) and
// Absent when it has none, as the sync's CREATE_OPTIONS tells.
const (
	spaceSecondaryEngine = "secondary engine"
	secondaryName        = "is"
)

func secondaryKey(db, table string) walk.Key {
	return walk.Key{Space: spaceSecondaryEngine, Parent: tablePath(db, table), Name: secondaryName}
}

// secondaryEngineOption returns the SECONDARY_ENGINE a list of table
// options sets: set is false when none does, and has is false for NULL.
func secondaryEngineOption(opts []*ast.TableOption) (has, set bool) {
	for _, o := range opts {
		if o != nil && optionName(o) == "SECONDARY_ENGINE" {
			set = true
			has = !strings.EqualFold(optionValue(o), "NULL")
		}
	}
	return has, set
}
