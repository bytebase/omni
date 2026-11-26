package review

import (
	"slices"
	"strconv"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/tidb/ast"
)

// A table's columns as the sync lists them are one key of its shape, read
// only when Present: each column in order, with its type, its character
// set, whether it is generated and which columns it then reads, and
// whether an index holds it. The sync gives it; any ALTER TABLE of the
// table, and CREATE INDEX or DROP INDEX on it, leaves it Unknown, and a
// cleared table path Absent, which says nothing. It settles what TiDB v8.5.5 judges against a column's
// character set (checkModifyCharsetAndCollation in pkg/ddl) and its
// generation (checkModifyGeneratedColumn and
// checkModifyColumnWithGeneratedColumnsConstraint).
const spaceColumnSet = "column set"

func columnSetKey(db, table string) walk.Key {
	return walk.Key{Space: spaceColumnSet, Parent: tablePath(db, table), Name: shapeName}
}

// Generation kinds of a column.
const (
	genVirtual = "v"
	genStored  = "s"
)

// setColumn is a column of a table's column set: its key, its type as
// information_schema writes it, its character set ("" for none), its
// generation kind ("" for a column that isn't generated) and the keys of
// the columns its expression reads, whether it is a key part of an index,
// and whether an expression index reads it.
type setColumn struct {
	name, typ, charset, gen string
	reads                   []string
	indexed, exprIndexed    bool
}

// syncedColumnSet reads a synced table's columns, false when the walk
// cannot read a name or an expression of them.
func syncedColumnSet(t *metadata.TableMetadata) (string, bool) {
	plain := make(map[string]bool)
	inExpr := make(map[string]bool)
	for _, i := range t.GetIndexes() {
		for _, e := range i.GetExpressions() {
			if k, ok := plainName(spaceColumn, e); ok {
				plain[k] = true
				continue
			}
			cols, ok := checkColumns(e)
			if !ok {
				return "", false
			}
			for _, c := range cols {
				inExpr[c] = true
			}
		}
	}
	var cols []setColumn
	for _, c := range t.GetColumns() {
		k, ok := plainName(spaceColumn, c.GetName())
		typ := strings.ToLower(strings.TrimSpace(c.GetType()))
		if !ok || strings.ContainsAny(typ, sepEntry+sepField+sepItem) {
			return "", false
		}
		sc := setColumn{name: k, typ: typ, charset: normalCharset(c.GetCharacterSet()), indexed: plain[k], exprIndexed: inExpr[k]}
		if g := c.GetGeneration(); g != nil && g.GetExpression() != "" {
			sc.gen = genVirtual
			if g.GetType() == metadata.GenerationMetadata_TYPE_STORED {
				sc.gen = genStored
			}
			if sc.reads, ok = checkColumns(g.GetExpression()); !ok {
				return "", false
			}
		}
		cols = append(cols, sc)
	}
	return encodeColumnSet(cols), true
}

func encodeColumnSet(cols []setColumn) string {
	parts := make([]string, len(cols))
	for n, c := range cols {
		flags := ""
		if c.indexed {
			flags += "i"
		}
		if c.exprIndexed {
			flags += "e"
		}
		parts[n] = strings.Join([]string{c.name, c.typ, c.charset, c.gen, strings.Join(c.reads, sepItem), flags}, sepField)
	}
	return strings.Join(parts, sepEntry)
}

func decodeColumnSet(kind string) ([]setColumn, bool) {
	if kind == "" {
		return nil, true
	}
	var out []setColumn
	for _, e := range strings.Split(kind, sepEntry) {
		f := strings.Split(e, sepField)
		if len(f) != 6 {
			return nil, false
		}
		c := setColumn{name: f[0], typ: f[1], charset: f[2], gen: f[3], indexed: strings.Contains(f[5], "i"), exprIndexed: strings.Contains(f[5], "e")}
		if f[4] != "" {
			c.reads = strings.Split(f[4], sepItem)
		}
		out = append(out, c)
	}
	return out, true
}

// columnSet returns a table's columns before the statement, false when the
// walk does not know them.
func (b *builder) columnSet(r tableRef) ([]setColumn, bool) {
	f, kind := b.l.Lookup(columnSetKey(r.db, r.name))
	if f != walk.Present {
		return nil, false
	}
	return decodeColumnSet(kind)
}

// findSetColumn returns the position of a column in a column set, -1 when
// it is not there.
func findSetColumn(cols []setColumn, key string) int {
	for n, c := range cols {
		if c.name == key {
			return n
		}
	}
	return -1
}

// charsetConvertible reports whether TiDB moves a column or a table from
// one character set to another: to utf8mb4 from utf8 or latin1, and to the
// same set (checkModifyCharsetAndCollation in v8.5.5's ddl, which refuses
// every other move with 8200, "Unsupported modify charset", verified for
// ascii and utf8mb4 to latin1, ascii to utf8 and utf8mb4, latin1 to utf8,
// and utf8mb4 to binary).
func charsetConvertible(from, to string) bool {
	return from == to || to == "utf8mb4" && (from == "utf8" || from == "latin1")
}

// charsetMaxLen is the most bytes a character of each set takes, which
// TiDB's VARCHAR length check reads (IsVarcharTooBigFieldLength).
var charsetMaxLen = map[string]int{"utf8mb4": 4, "utf8": 3, "gbk": 2, "gb18030": 4, "latin1": 1, "ascii": 1, "binary": 1}

// varcharLength returns the length of a VARCHAR or VARBINARY type as
// information_schema writes it, 0 for another type, and false when the
// walk cannot read it.
func varcharLength(typ string) (int, bool) {
	for _, p := range []string{"varchar(", "varbinary("} {
		if rest, ok := strings.CutPrefix(typ, p); ok {
			n, err := strconv.Atoi(strings.TrimSuffix(rest, ")"))
			return n, err == nil && strings.HasSuffix(rest, ")")
		}
	}
	return 0, true
}

// convertColumnsFail fails ALTER TABLE ... CONVERT TO CHARACTER SET, alone
// in its statement, when a column has a character set TiDB will not
// convert to the one named (8200). TiDB checks the table's set first, with
// the same error, then each column in order, a VARCHAR's length under the
// new set (1074) before the column's set: the walk stops at a VARCHAR that
// may be too long.
func (b *builder) convertColumnsFail(v *ast.AlterTableStmt, r tableRef) bool {
	if len(v.Commands) != 1 || v.Commands[0].Type != ast.ATConvertCharset {
		return false
	}
	c := v.Commands[0]
	to := normalCharset(c.Name)
	maxLen, ok := charsetMaxLen[to]
	if !ok || to == "binary" || c.NewName != "" && collationCharset(c.NewName) != to {
		return false
	}
	cols, ok := b.columnSet(r)
	if !ok {
		return false
	}
	for _, col := range cols {
		n, ok := varcharLength(col.typ)
		if !ok || n*maxLen > 65535 {
			return false
		}
		if col.charset == "" || col.charset == "binary" || charsetConvertible(col.charset, to) {
			continue
		}
		b.fail(rangeOf(c.Loc), nil, review.UnsupportedForm, engineName+" can't change the character set of column "+review.Code(ident(col.name))+" of "+review.Code(r.written)+" from "+review.Code(col.charset)+" to "+review.Code(to)+".")
		return true
	}
	return false
}

// inForeignKey reports whether a column is a column of one of the table's
// foreign keys or one another table's foreign key references, which TiDB
// judges a change of first (3780).
func (b *builder) inForeignKey(r tableRef, key string) bool {
	for _, o := range b.ownFKs(r) {
		if contains(o.fk.cols, key) {
			return true
		}
	}
	for _, u := range b.referencingFKs(r) {
		if contains(u.fk.refCols, key) {
			return true
		}
	}
	return false
}

// defCharset returns the character set a column definition gives a
// character string column of the table: the one it names, or else the
// table's, false when the walk does not know it or the definition names a
// collation, which may pick the set.
func (b *builder) defCharset(r tableRef, def *ast.ColumnDef) (string, bool) {
	dt := def.TypeName
	if dt == nil || dt.Binary {
		return "", false
	}
	for _, cc := range def.Constraints {
		if cc.Type == ast.ColConstrCollate {
			return "", false
		}
	}
	if cs := normalCharset(dt.Charset); cs != "" {
		return cs, dt.Collate == "" || collationCharset(dt.Collate) == cs
	}
	if dt.Collate != "" {
		return "", false
	}
	f, cs := b.l.Lookup(tableCharsetKey(r.db, r.name))
	return cs, f == walk.Present && cs != ""
}

// modifyCharsetFails fails MODIFY or CHANGE of a character string column
// to a character string type of a character set TiDB will not move the
// column to (8200, verified for a utf8mb4 column of a latin1 table that
// MODIFY gives no set, so it takes the table's).
func (b *builder) modifyCharsetFails(r tableRef, col string, def *ast.ColumnDef, rng review.Range) bool {
	k, ok := b.key(spaceColumn, col)
	if !ok || b.inForeignKey(r, k) {
		return false
	}
	cols, ok := b.columnSet(r)
	if !ok {
		return false
	}
	n := findSetColumn(cols, k)
	if n < 0 || cols[n].charset == "" || cols[n].charset == "binary" {
		return false
	}
	from, ok1 := parseColumnType(syncedColumnType(cols[n].typ))
	to, ok2 := parseColumnType(astColumnType(def.TypeName))
	if !ok1 || !ok2 || !textTypes[from.base] || !textTypes[to.base] {
		return false
	}
	cs, ok := b.defCharset(r, def)
	if !ok || cs == "binary" || charsetConvertible(cols[n].charset, cs) {
		return false
	}
	b.fail(rng, nil, review.UnsupportedForm, engineName+" can't change the character set of column "+review.Code(ident(col))+" of "+review.Code(r.written)+" from "+review.Code(cols[n].charset)+" to "+review.Code(cs)+".")
	return true
}

// keepsStorage reports whether a column definition keeps a column's type
// and character set, a CHAR or VARCHAR growing no longer than an index of
// it can hold, so TiDB judges nothing of the change before its generation
// checks. The definition sets nothing but the type, the generation, NULL
// or NOT NULL, and a comment.
func (b *builder) keepsStorage(r tableRef, old setColumn, def *ast.ColumnDef) bool {
	dt := def.TypeName
	if dt == nil || def.DefaultValue != nil || def.OnUpdate != nil || def.AutoIncrement || def.AutoRandom {
		return false
	}
	for _, cc := range def.Constraints {
		switch cc.Type {
		case ast.ColConstrNotNull, ast.ColConstrNull, ast.ColConstrComment:
		default:
			return false
		}
	}
	if old.charset != "" {
		if cs, ok := b.defCharset(r, def); !ok || cs != old.charset {
			return false
		}
	} else if dt.Charset != "" || dt.Collate != "" || dt.Binary {
		return false
	}
	name := strings.ToLower(dt.Name)
	if base, rest, ok := strings.Cut(old.typ, "("); ok && (base == "char" || base == "varchar") && name == base {
		n, err := strconv.Atoi(strings.TrimSuffix(rest, ")"))
		// 768 characters of four bytes fill an index key of 3072.
		return err == nil && (dt.Length == n || dt.Length > n && dt.Length <= 768)
	}
	written := name
	if dt.Length > 0 {
		written += "(" + strconv.Itoa(dt.Length) + ")"
	}
	if dt.Unsigned {
		written += " unsigned"
	}
	return written == old.typ
}

// generatedModifyFails fails MODIFY or CHANGE, with no position, that a
// table's generated columns refute, verified on v8.5.5:
//   - a column made VIRTUAL from stored or not generated, or the other
//     way (3106, "Changing the STORED status");
//   - an indexed VIRTUAL column given an expression reading other columns,
//     or the same ones in another order (3106, "modifying an indexed
//     column");
//   - a column renamed that a generated column reads (3108).
//
// The walk claims these only for a definition that keeps the column's
// storage (keepsStorage), of a column in no foreign key, whose new
// expression reads columns of the table before it.
func (b *builder) generatedModifyFails(r tableRef, c *ast.AlterTableCmd, col string, rng review.Range) bool {
	def := c.Column
	if def == nil || c.First || c.After != "" || c.IfExists {
		return false
	}
	k, ok1 := b.key(spaceColumn, col)
	nk, ok2 := b.key(spaceColumn, def.Name)
	if !ok1 || !ok2 || b.inForeignKey(r, k) {
		return false
	}
	cols, ok := b.columnSet(r)
	if !ok {
		return false
	}
	n := findSetColumn(cols, k)
	if n < 0 || !b.keepsStorage(r, cols[n], def) {
		return false
	}
	old := cols[n]
	var reads []string
	if g := def.Generated; g != nil {
		if g.Stored || g.Expr == nil || nondeterministicCall(b.w.sql, g.Expr) != "" {
			return false
		}
		reads, ok = exprColumns(g.Expr)
		if !ok {
			return false
		}
		for _, read := range reads {
			j := findSetColumn(cols, read)
			if j < 0 || j == n || cols[j].gen != "" && j > n {
				return false
			}
		}
	}
	gcol := review.Code(ident(def.Name))
	if (old.gen == genVirtual) != (def.Generated != nil) {
		b.fail(rng, nil, review.UnsupportedForm, engineName+" can't change whether generated column "+gcol+" is stored.")
		return true
	}
	if old.gen == genVirtual && k == nk && old.indexed && !slices.Equal(reads, old.reads) {
		b.fail(rng, nil, review.UnsupportedForm, engineName+" can't change the expression of generated column "+gcol+", which an index holds.")
		return true
	}
	if k == nk || old.exprIndexed {
		return false
	}
	for _, other := range cols {
		if other.gen != "" && contains(other.reads, k) {
			b.fail(rng, nil, review.UnsupportedForm, "Column "+review.Code(ident(col))+" is read by generated column "+review.Code(ident(other.name))+", and can't be renamed.")
			return true
		}
	}
	return false
}

// columnCharset returns a column's character set before the statement,
// false when the walk does not know it or the column has none.
func (b *builder) columnCharset(r tableRef, col string) (string, bool) {
	k, ok := b.key(spaceColumn, col)
	if !ok {
		return "", false
	}
	cols, ok := b.columnSet(r)
	if !ok {
		return "", false
	}
	n := findSetColumn(cols, k)
	if n < 0 || cols[n].charset == "" {
		return "", false
	}
	return cols[n].charset, true
}

// createdColumnCharset returns the character set a CREATE TABLE gives a
// character string column: the one its definition names, or else the
// table's, false when the statement names neither, so the column takes
// the database's.
func createdColumnCharset(v *ast.CreateTableStmt, col string) (string, bool) {
	for _, def := range v.Columns {
		if def == nil || !strings.EqualFold(def.Name, col) || def.TypeName == nil {
			continue
		}
		for _, cc := range def.Constraints {
			if cc.Type == ast.ColConstrCollate {
				return "", false
			}
		}
		if cs := normalCharset(def.TypeName.Charset); cs != "" {
			return cs, true
		}
		if cs := collationCharset(def.TypeName.Collate); cs != "" {
			return cs, true
		}
		var fromCharset, fromCollate string
		for _, opt := range v.Options {
			if s, ok := charsetOption(opt); ok {
				if strings.Contains(strings.ToLower(opt.Name), "collate") {
					fromCollate = s
				} else {
					fromCharset = s
				}
			}
		}
		if fromCharset != "" {
			return fromCharset, true
		}
		return fromCollate, fromCollate != ""
	}
	return "", false
}
