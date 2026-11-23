package review

import (
	"strings"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/review"
)

// mariadbCharsets are the character sets MariaDB 11.8 has (SHOW
// CHARACTER SET on mariadb:11.8.8), with utf8 for utf8mb3 as the default
// old_mode reads it. A server has no others: character sets are compiled
// in.
var mariadbCharsets = map[string]bool{
	"armscii8": true, "ascii": true, "big5": true, "binary": true, "cp1250": true, "cp1251": true,
	"cp1256": true, "cp1257": true, "cp850": true, "cp852": true, "cp866": true, "cp932": true,
	"dec8": true, "eucjpms": true, "euckr": true, "gb2312": true, "gbk": true, "geostd8": true,
	"greek": true, "hebrew": true, "hp8": true, "keybcs2": true, "koi8r": true, "koi8u": true,
	"latin1": true, "latin2": true, "latin5": true, "latin7": true, "macce": true, "macroman": true,
	"sjis": true, "swe7": true, "tis620": true, "ucs2": true, "ujis": true, "utf16": true,
	"utf16le": true, "utf32": true, "utf8": true, "utf8mb3": true, "utf8mb4": true,
}

// normalCharset returns a character set's name as MariaDB resolves it,
// utf8 as utf8mb3, "" for a name that is no character set.
func normalCharset(name string) string {
	n := strings.ToLower(strings.Trim(strings.TrimSpace(name), "'\"`"))
	if !mariadbCharsets[n] {
		return ""
	}
	if n == "utf8" {
		return "utf8mb3"
	}
	return n
}

// collationCharset returns the character set a collation belongs to, the
// name before its first underscore (binary is binary's own), "" when no
// character set begins it: MariaDB's uca1400_ai_ci and its kin serve
// every Unicode character set.
func collationCharset(collation string) string {
	c := strings.ToLower(strings.Trim(strings.TrimSpace(collation), "'\"`"))
	cs, _, _ := strings.Cut(c, "_")
	return normalCharset(cs)
}

// mariadbEngines are the storage engines a MariaDB server may have, by
// every name it takes for them: the ones mariadb:11.8.8 builds in (SHOW
// ENGINES), the ones a plugin MariaDB ships or a third party writes may
// add, and the aliases it resolves. Any other name is 1286 under the
// default sql_mode, which has NO_ENGINE_SUBSTITUTION (verified on
// mariadb:11.8.8).
var mariadbEngines = map[string]bool{
	"innodb": true, "innobase": true, "xtradb": true, "myisam": true, "aria": true, "maria": true,
	"memory": true, "heap": true, "csv": true, "mrg_myisam": true, "mrg_maria": true, "merge": true,
	"archive": true, "blackhole": true, "federated": true, "federatedx": true, "sequence": true,
	"sql_sequence": true, "performance_schema": true, "s3": true, "rocksdb": true, "tokudb": true,
	"mroonga": true, "sphinx": true, "spider": true, "connect": true, "columnstore": true,
	"oqgraph": true, "cassandra": true, "example": true, "ndbcluster": true, "ndb": true,
	"partition": true, "wsrep": true, "test_sql_discovery": true, "pbxt": true, "ibmdb2i": true,
}

// optionName returns a table option's name with CHARSET for CHARACTER SET.
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

// unknownCharset is the text for a character set MariaDB doesn't have
// (1115).
func unknownCharset(name string) string {
	return "Character set " + review.Code(name) + " doesn't exist."
}

// collationMismatch is the text for a collation of another character set
// than the one beside it (1253).
func collationMismatch(collation, charset string) string {
	return "Collation " + review.Code(collation) + " isn't one of character set " + review.Code(charset) + "."
}

// unknownEngine is the text for a storage engine no MariaDB server has
// (1286).
func unknownEngine(name string) string {
	return engineName + " has no storage engine " + review.Code(name) + "."
}

// optionTextFails fails a CREATE TABLE or ALTER TABLE whose table
// options, CONVERT TO, or column types MariaDB refuses while it parses
// them, before it looks the table up (each verified on mariadb:11.8.8,
// on a table that doesn't exist too):
//
//   - a character set MariaDB doesn't have (1115), in a table option,
//     CONVERT TO, or a column type;
//   - a storage engine no server has (1286), unless the session set
//     sql_mode, which may let MariaDB substitute the default one;
//   - a collation of another character set than the one written beside it
//     (1253), in a column type (ASCII is latin1), CONVERT TO, or a table's
//     CHARSET and then COLLATE.
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
			if o != nil && optionName(o) == "ENGINE" && !mariadbEngines[strings.ToLower(optionValue(o))] {
				b.fail(rng, nil, review.UnsupportedForm, unknownEngine(optionValue(o)))
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
			b.fail(rangeOf(c.Loc), nil, review.UnsupportedForm, collationMismatch(strings.ToLower(c.TypeName.Collate), cs))
			return true
		}
	}
	if convert != nil && convert.NewName != "" {
		cs, ccs := normalCharset(convert.Name), collationCharset(convert.NewName)
		if cs != "" && ccs != "" && cs != ccs {
			b.fail(rng, nil, review.UnsupportedForm, collationMismatch(strings.ToLower(convert.NewName), cs))
			return true
		}
	}
	// A table's CHARSET and then its COLLATE; two of either, or COLLATE
	// first, the parser judges or the walk leaves alone.
	var charset, collation *ast.TableOption
	for _, o := range opts {
		if o == nil {
			continue
		}
		switch optionName(o) {
		case "CHARSET":
			if charset != nil || collation != nil || strings.EqualFold(optionValue(o), "DEFAULT") {
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
			b.fail(rng, nil, review.UnsupportedForm, collationMismatch(strings.ToLower(optionValue(collation)), cs))
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
