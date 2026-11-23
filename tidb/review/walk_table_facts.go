package review

import (
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/parser"
)

// A table's character set is one key of its own under the table's path,
// read only when Present, with the character set as its kind. The walk
// knows it from the sync (Bytebase reads it from the table's collation) or
// from a statement that names it; a table created without one takes the
// database's, which the walk leaves Unknown.
const (
	spaceTableCharset = "table charset"
	tableCharsetName  = "is"
)

func tableCharsetKey(db, table string) walk.Key {
	return walk.Key{Space: spaceTableCharset, Parent: tablePath(db, table), Name: tableCharsetName}
}

// normalCharset folds a character set name as TiDB reads it: utf8mb3 is
// utf8. "" for DEFAULT or none.
func normalCharset(name string) string {
	cs := strings.ToLower(strings.TrimSpace(name))
	switch cs {
	case "default":
		return ""
	case "utf8mb3":
		return "utf8"
	}
	return cs
}

// collationCharset is the character set a collation belongs to: binary's
// is binary, every other's the prefix before its first underscore, as
// Bytebase's TiDB sync reads it.
func collationCharset(collation string) string {
	c := strings.ToLower(strings.TrimSpace(collation))
	if c == "" || c == "default" {
		return ""
	}
	if c == "binary" {
		return "binary"
	}
	cs, _, _ := strings.Cut(c, "_")
	return normalCharset(cs)
}

// charsetOption returns what a table option says of the table's character
// set: the set it names, and whether it names one at all.
func charsetOption(opt *ast.TableOption) (cs string, names bool) {
	if opt == nil {
		return "", false
	}
	switch strings.ToLower(opt.Name) {
	case "charset", "character set", "default charset", "default character set":
		return normalCharset(opt.Value), true
	case "collate", "default collate":
		return collationCharset(opt.Value), true
	}
	return "", false
}

// alterCharset returns the character set an ALTER TABLE gives the table:
// cs is "" when the walk cannot tell it, and names is false when the
// statement leaves it alone. A CHARSET option wins over a COLLATE one, as
// TiDB resolves the pair.
func alterCharset(v *ast.AlterTableStmt) (cs string, names bool) {
	var fromCharset, fromCollate string
	for _, c := range v.Commands {
		switch c.Type {
		case ast.ATConvertCharset:
			names = true
			if s := normalCharset(c.Name); s != "" {
				fromCharset = s
			} else if s := collationCharset(c.NewName); s != "" {
				fromCollate = s
			}
		case ast.ATTableOption:
			s, ok := charsetOption(c.Option)
			if !ok {
				continue
			}
			names = true
			if n := strings.ToLower(c.Option.Name); strings.Contains(n, "collate") {
				fromCollate = s
			} else {
				fromCharset = s
			}
		}
	}
	if fromCharset != "" {
		return fromCharset, names
	}
	return fromCollate, names
}

// charsetChangeFails fails an ALTER TABLE that moves a table to another
// character set TiDB will not convert it to: "Unsupported modify charset"
// (8200, verified on v8.5.5 for CHARSET, COLLATE and CONVERT TO binary on a
// utf8mb4 table). TiDB changes a table's character set only to utf8mb4, from
// utf8 or latin1 (ddl's modifiableCharsetAndCollation); the walk claims
// only a move to a set other than utf8mb4, and leaves utf8 from latin1 or
// ascii alone.
func (b *builder) charsetChangeFails(v *ast.AlterTableStmt, r tableRef) bool {
	to, names := alterCharset(v)
	if !names || to == "" || to == "utf8mb4" {
		return false
	}
	f, from := b.l.Lookup(tableCharsetKey(r.db, r.name))
	if f != walk.Present || from == "" || from == to {
		return false
	}
	if to == "utf8" && (from == "latin1" || from == "ascii") {
		return false
	}
	b.fail(b.anchor(), nil, review.UnsupportedForm, engineName+" can't change the character set of "+review.Code(r.written)+" from "+review.Code(from)+" to "+review.Code(to)+".")
	return true
}

// charsetEffects records the character set an ALTER TABLE leaves the
// table with.
func (b *builder) charsetEffects(v *ast.AlterTableStmt, r tableRef) {
	to, names := alterCharset(v)
	switch {
	case !names:
	case to == "":
		b.Set(tableCharsetKey(r.db, r.name), walk.Unknown, "")
	default:
		b.Set(tableCharsetKey(r.db, r.name), walk.Present, to)
	}
}

// createdCharset records the character set a CREATE TABLE names, when it
// names one.
func (b *builder) createdCharset(v *ast.CreateTableStmt, r tableRef) {
	var fromCharset, fromCollate string
	for _, opt := range v.Options {
		s, ok := charsetOption(opt)
		if !ok {
			continue
		}
		if strings.Contains(strings.ToLower(opt.Name), "collate") {
			fromCollate = s
		} else {
			fromCharset = s
		}
	}
	cs := fromCharset
	if cs == "" {
		cs = fromCollate
	}
	if cs != "" {
		b.Set(tableCharsetKey(r.db, r.name), walk.Present, cs)
	}
}

// ttlColumn returns the column a TTL option's expression leads with, ""
// when it is not one the walk can read.
func ttlColumn(opt *ast.TableOption) string {
	if opt == nil || !strings.EqualFold(opt.Name, "TTL") {
		return ""
	}
	toks := parser.Tokenize(opt.Value)
	if len(toks) == 0 {
		return ""
	}
	t := tokenText(opt.Value, toks[0])
	if strings.HasPrefix(t, "`") && strings.HasSuffix(t, "`") && len(t) >= 2 {
		t = strings.ReplaceAll(t[1:len(t)-1], "``", "`")
	}
	if t == "" || strings.ContainsAny(t, ".'\"") {
		return ""
	}
	return t
}

// columnRefDefault fails ALTER COLUMN SET DEFAULT (col): TiDB evaluates a
// default expression with no columns in scope, so a bare column reference
// is unknown (1054, verified on v8.5.5 with the column present).
func (b *builder) columnRefDefault(c *ast.AlterTableCmd, rng review.Range) bool {
	if c.DefaultExpr == nil {
		return false
	}
	ref, ok := unwrapParens(c.DefaultExpr).(*ast.ColumnRef)
	if !ok || ref.Column == "" || ref.Star {
		return false
	}
	b.fail(rng, nil, review.ObjectNotFound, "Column "+review.Code(ident(ref.Column))+" is unknown in a "+review.Code("DEFAULT")+" expression, which reads no columns.")
	return true
}

// Whether a table has an AUTO_RANDOM column is one key under its path, read
// only when Present, with kind autoRandomYes or autoRandomNo. Bytebase's
// sync shows it in the table's sharding info (PK_AUTO_RANDOM_BITS) and on
// the column's default; a CREATE TABLE shows it in its column list.
const (
	spaceAutoRandom = "auto random"
	autoRandomName  = "has"
	autoRandomYes   = "yes"
	autoRandomNo    = "no"
)

func autoRandomKey(db, table string) walk.Key {
	return walk.Key{Space: spaceAutoRandom, Parent: tablePath(db, table), Name: autoRandomName}
}

// syncedAutoRandom reads whether a synced table has an AUTO_RANDOM column.
func syncedAutoRandom(t *metadata.TableMetadata) string {
	if strings.Contains(strings.ToUpper(t.GetShardingInfo()), "PK_AUTO_RANDOM_BITS") {
		return autoRandomYes
	}
	for _, c := range t.GetColumns() {
		if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(c.GetDefault())), "AUTO_RANDOM") {
			return autoRandomYes
		}
	}
	return autoRandomNo
}

// createdAutoRandom records whether a CREATE TABLE's columns hold an
// AUTO_RANDOM one.
func (b *builder) createdAutoRandom(v *ast.CreateTableStmt, r tableRef) {
	has := autoRandomNo
	for _, c := range v.Columns {
		if c != nil && c.AutoRandom {
			has = autoRandomYes
		}
	}
	b.Set(autoRandomKey(r.db, r.name), walk.Present, has)
}

// autoRandomFails fails AUTO_RANDOM_BASE on a table with no AUTO_RANDOM
// column (8216, verified on v8.5.5), as the table was before the statement.
func (b *builder) autoRandomFails(v *ast.AlterTableStmt, r tableRef) bool {
	for _, c := range v.Commands {
		if c.Type != ast.ATTableOption || c.Option == nil || !strings.EqualFold(c.Option.Name, "auto_random_base") {
			continue
		}
		if f, has := b.l.Lookup(autoRandomKey(r.db, r.name)); f == walk.Present && has == autoRandomNo {
			b.fail(rangeOf(c.Loc), nil, review.UnsupportedForm, "Table "+review.Code(r.written)+" has no "+review.Code("AUTO_RANDOM")+" column to rebase.")
			return true
		}
	}
	return false
}

// autoRandomEffects records what an ALTER TABLE does to whether the table
// has an AUTO_RANDOM column: a column it adds with the attribute gives it
// one, and a column it modifies, changes, or drops, or a primary key it
// adds or drops, leaves the answer Unknown.
func (b *builder) autoRandomEffects(v *ast.AlterTableStmt, r tableRef) {
	k := autoRandomKey(r.db, r.name)
	for _, c := range v.Commands {
		cols := c.Columns
		if c.Column != nil {
			cols = append([]*ast.ColumnDef{c.Column}, cols...)
		}
		for _, col := range cols {
			if col != nil && col.AutoRandom {
				b.Set(k, walk.Unknown, "")
				return
			}
		}
		switch c.Type {
		case ast.ATModifyColumn, ast.ATChangeColumn, ast.ATDropColumn, ast.ATDropConstraint, ast.ATAddConstraint, ast.ATAddIndex:
			b.Set(k, walk.Unknown, "")
			return
		}
	}
}
