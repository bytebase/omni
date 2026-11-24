package review

import (
	"encoding/json"
	"math/big"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/parser"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// obKeyKind distinguishes the index column type constraints OceanBase enforces.
type obKeyKind uint8

const (
	obKeyPlain obKeyKind = iota
	obKeyFulltext
	obKeySpatial
	obKeyVector
)

// onlyCommentAlter reports whether an ALTER TABLE consists solely of COMMENT
// table options, the only ALTER TABLE form OceanBase accepts on a view.
func onlyCommentAlter(cmds []*ast.AlterTableCmd) bool {
	if len(cmds) == 0 {
		return false
	}
	for _, c := range cmds {
		if c == nil || c.Type != ast.ATTableOption || c.Option == nil || !strings.EqualFold(c.Option.Name, "COMMENT") {
			return false
		}
	}
	return true
}

// alterTableOptionsFail checks OceanBase ALTER TABLE options, character set /
// collation conversions, column groups, and TTL clauses that the engine
// unconditionally refuses.
func (b *builder) alterTableOptionsFail(r tableRef, cmds []*ast.AlterTableCmd) bool {
	for _, c := range cmds {
		if c == nil {
			continue
		}
		rng := rangeOf(c.Loc)
		switch c.Type {
		case ast.ATTableOption:
			if c.Option != nil && b.alterTableOptionFailsOB(r, c.Option, rng) {
				return true
			}
		case ast.ATConvertCharset:
			if b.checkCharsetNameOB(c.Name, rng) {
				return true
			}
			if c.NewName != "" {
				if b.checkCollationNameOB(c.NewName, rng) {
					return true
				}
				if b.checkCharsetCollationPairOB(c.Name, c.NewName, rng) {
					return true
				}
			}
		case ast.ATAddColumnGroup:
			if b.checkColumnGroupCmdOB(c, true, rng) {
				return true
			}
		case ast.ATDropColumnGroup:
			if b.checkColumnGroupCmdOB(c, false, rng) {
				return true
			}
		}
	}
	return false
}

func (b *builder) alterTableOptionFailsOB(r tableRef, opt *ast.TableOption, rng review.Range) bool {
	name := strings.ToUpper(opt.Name)
	value := unquoteOBOption(opt.Value)
	switch {
	case name == "FORMAT" || name == "LOCATION" || name == "PATTERN" || name == "PROPERTIES" || name == "PARTITION_TYPE":
		b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses external table option "+review.Code(name)+" on a base table.")
		return true
	case name == "TABLE_ID":
		b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses altering "+review.Code("TABLE_ID")+".")
		return true
	case strings.HasPrefix(name, "ORGANIZATION"):
		b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses "+review.Code("ORGANIZATION")+" on ALTER TABLE.")
		return true
	case name == "MERGE_ENGINE":
		b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses altering "+review.Code("MERGE_ENGINE")+".")
		return true
	case name == "AUTO_REFRESH":
		b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses "+review.Code("AUTO_REFRESH")+" on ALTER TABLE.")
		return true
	case name == "STORAGE_CACHE_POLICY":
		b.fail(rng, nil, review.UnsupportedForm, b.engine()+" has no storage cache policy for a table.")
		return true
	case name == "JSON":
		b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses LOB storage parameters on ALTER TABLE.")
		return true
	case name == "TABLE_MODE":
		switch strings.ToLower(value) {
		case "normal", "queuing", "moderate", "super", "extreme":
		default:
			b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses "+review.Code("TABLE_MODE")+" value "+review.Code(value)+".")
			return true
		}
	case name == "DUPLICATE_SCOPE":
		switch strings.ToLower(value) {
		case "none", "cluster":
		default:
			b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses "+review.Code("DUPLICATE_SCOPE")+" value "+review.Code(value)+".")
			return true
		}
	case name == "DUPLICATE_READ_CONSISTENCY":
		switch strings.ToLower(value) {
		case "strong", "weak":
		default:
			b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses "+review.Code("DUPLICATE_READ_CONSISTENCY")+" value "+review.Code(value)+".")
			return true
		}
	case name == "SEMISTRUCT_ENCODING_TYPE":
		switch strings.ToLower(value) {
		case "", "none", "encoding":
		default:
			b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses "+review.Code("SEMISTRUCT_ENCODING_TYPE")+" value "+review.Code(value)+".")
			return true
		}
	case name == "AUTO_INCREMENT_MODE":
		switch strings.ToLower(value) {
		case "order", "noorder":
		default:
			b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses "+review.Code("AUTO_INCREMENT_MODE")+" value "+review.Code(value)+".")
			return true
		}
	case name == "KV_ATTRIBUTES":
		if value != "" && !json.Valid([]byte(value)) {
			b.fail(rng, nil, review.UnsupportedForm, b.engine()+" requires valid JSON for "+review.Code("KV_ATTRIBUTES")+".")
			return true
		}
	case name == "BLOCK_SIZE":
		if n, err := strconv.ParseInt(value, 10, 64); err != nil || n < 1024 || n > 1048576 {
			b.fail(rng, nil, review.UnsupportedForm, b.engine()+" requires "+review.Code("BLOCK_SIZE")+" between 1024 and 1048576.")
			return true
		}
	case name == "TABLET_SIZE":
		if n, err := strconv.ParseInt(value, 10, 64); err != nil || n%(2<<20) != 0 {
			b.fail(rng, nil, review.UnsupportedForm, b.engine()+" requires "+review.Code("TABLET_SIZE")+" to be a multiple of 2 MiB.")
			return true
		}
	case name == "COMPRESSION":
		comp := strings.Trim(strings.TrimSpace(opt.Value), "'\"")
		switch strings.ToLower(comp) {
		case "none", "lz4_1.0", "lz4_1.9.1", "snappy_1.0", "zstd_1.0", "zstd_1.3.8":
		default:
			b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses "+review.Code("COMPRESSION")+" value "+review.Code(comp)+".")
			return true
		}
	case name == "CHARSET" || name == "CHARACTER SET" || name == "DEFAULT CHARSET" || name == "DEFAULT CHARACTER SET":
		return b.checkCharsetNameOB(opt.Value, rng)
	case name == "COLLATE" || name == "DEFAULT COLLATE":
		return b.checkCollationNameOB(opt.Value, rng)
	case name == "DYNAMIC_PARTITION_POLICY":
		return b.checkDynamicPartitionPolicyOB(opt.Value, rng)
	}
	return false
}

var (
	obDynPartTimeRe = regexp.MustCompile(`(?i)^\s*(-1|0|\+?\d+\s*(hour|day|week|month|year))\s*$`)
	obDynPartZoneRe = regexp.MustCompile(`^([+-])(\d{1,2}):(\d{2})$`)
)

func (b *builder) checkDynamicPartitionPolicyOB(text string, rng review.Range) bool {
	entries, ok := obOptionEntries(text)
	if !ok {
		b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses invalid "+review.Code("DYNAMIC_PARTITION_POLICY")+" syntax.")
		return true
	}
	for _, entry := range entries {
		key, value := strings.ToUpper(entry[0]), entry[1]
		switch key {
		case "TIME_UNIT":
			switch strings.ToLower(value) {
			case "hour", "day", "week", "month", "year":
			default:
				b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses "+review.Code("DYNAMIC_PARTITION_POLICY")+" option "+review.Code(key)+".")
				return true
			}
		case "PRECREATE_TIME", "EXPIRE_TIME":
			if !obDynPartTimeRe.MatchString(value) {
				b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses "+review.Code("DYNAMIC_PARTITION_POLICY")+" option "+review.Code(key)+".")
				return true
			}
		case "TIME_ZONE":
			if !strings.EqualFold(value, "default") && !obValidOffsetZone(value) {
				b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses "+review.Code("TIME_ZONE")+" value "+review.Code(value)+" in "+review.Code("DYNAMIC_PARTITION_POLICY")+".")
				return true
			}
		case "BIGINT_PRECISION":
			switch strings.ToLower(value) {
			case "none", "us", "ms", "s":
			default:
				b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses "+review.Code("DYNAMIC_PARTITION_POLICY")+" option "+review.Code(key)+".")
				return true
			}
		}
	}
	return false
}

func obValidOffsetZone(text string) bool {
	m := obDynPartZoneRe.FindStringSubmatch(text)
	if m == nil {
		return false
	}
	h, _ := strconv.Atoi(m[2])
	mm, _ := strconv.Atoi(m[3])
	if mm > 59 {
		return false
	}
	minutes := h*60 + mm
	if m[1] == "-" {
		return minutes <= 12*60+59
	}
	return minutes <= 13*60
}

func (b *builder) checkTTLOptionOB(r tableRef, opt *ast.TableOption, rng review.Range) bool {
	for _, colName := range opt.TTLColumns {
		ck, ok := b.key(spaceColumn, colName)
		if !ok {
			continue
		}
		colK := colKey(r.db, r.name, ck)
		b.need(colK, walk.Present, nil, rng, notFound("column", ident(colName), r.written), "")
		if b.touchedCols[ck] {
			continue
		}
		if f, _ := b.l.Lookup(colK); f == walk.Present {
			if typ := b.w.colType(r, ck); typ != "" && typ != "datetime" && typ != "timestamp" {
				b.fail(rng, []walk.Key{colK}, review.UnsupportedForm, b.engine()+" requires a DATETIME or TIMESTAMP column for TTL, not "+review.Code(ident(colName))+".")
				return true
			}
		}
	}
	return false
}

func (b *builder) checkColumnGroupCmdOB(c *ast.AlterTableCmd, isAdd bool, rng review.Range) bool {
	delayed := c.Option != nil && strings.EqualFold(c.Option.Name, "DELAYED")
	if !isAdd && delayed {
		b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses "+review.Code("DELAYED")+" on DROP COLUMN GROUP.")
		return true
	}
	allCols, eachCol, named, dup := parseOBColumnGroups(c.Name)
	if named {
		b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses named column groups.")
		return true
	}
	if dup != "" {
		b.fail(rng, nil, review.DuplicateName, statesTwice("declares", "column group", dup))
		return true
	}
	if isAdd && delayed && !eachCol && allCols {
		b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses "+review.Code("DELAYED")+" on ADD COLUMN GROUP without EACH COLUMN.")
		return true
	}
	return false
}

func parseOBColumnGroups(text string) (allCols, eachCol, named bool, dup string) {
	var elem []parser.Token
	depth := 0
	flush := func() bool {
		if len(elem) == 0 {
			return false
		}
		for _, tok := range elem {
			if tok.Type == '(' {
				named = true
				return true
			}
		}
		if strings.EqualFold(elem[0].Str, "ALL") {
			if allCols {
				dup = "ALL COLUMNS"
				return true
			}
			allCols = true
		} else {
			if eachCol {
				dup = "EACH COLUMN"
				return true
			}
			eachCol = true
		}
		elem = elem[:0]
		return false
	}
	for _, tok := range parser.Tokenize(text) {
		switch {
		case tok.Type == '(':
			depth++
		case tok.Type == ')':
			depth--
		case tok.Type == ',' && depth == 0:
			if flush() {
				return
			}
			continue
		}
		elem = append(elem, tok)
	}
	flush()
	return
}

// indexOptionsFailOB checks index options that OceanBase refuses on non-vector
// keys (WITH, DATA_TABLE_ID, INDEX_TABLE_ID, VIRTUAL_COLUMN_ID,
// STORAGE_CACHE_POLICY) and invalid BLOCK_SIZE on any key.
func (b *builder) indexOptionsFailOB(opts []*ast.IndexOption, kind obKeyKind, rng review.Range) bool {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		name := strings.ToUpper(opt.Name)
		if name == "BLOCK_SIZE" {
			if lit, ok := opt.Value.(*ast.IntLit); ok && (lit.Value < 1024 || lit.Value > 1048576) {
				b.fail(rng, nil, review.UnsupportedForm, b.engine()+" requires "+review.Code("BLOCK_SIZE")+" between 1024 and 1048576.")
				return true
			}
		}
		if kind == obKeyVector {
			continue
		}
		switch name {
		case "WITH":
			b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses "+review.Code("WITH")+" index parameters on a non-vector index.")
			return true
		case "DATA_TABLE_ID", "INDEX_TABLE_ID", "VIRTUAL_COLUMN_ID":
			b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses "+review.Code(name)+" on an index.")
			return true
		case "STORAGE_CACHE_POLICY":
			b.fail(rng, nil, review.UnsupportedForm, b.engine()+" has no storage cache policy for an index.")
			return true
		}
	}
	return false
}

func constraintKeyKindOB(c *ast.Constraint) (obKeyKind, bool) {
	switch c.Type {
	case ast.ConstrPrimaryKey, ast.ConstrUnique, ast.ConstrIndex:
		return obKeyPlain, true
	case ast.ConstrFulltextIndex:
		return obKeyFulltext, true
	case ast.ConstrSpatialIndex:
		return obKeySpatial, true
	case ast.ConstrVector:
		return obKeyVector, true
	}
	return obKeyPlain, false
}

// checkKeyPartsOB validates an index/constraint's options, duplicate columns,
// column existence, functional bare-column expressions, prefix lengths, and
// FULLTEXT/SPATIAL/VECTOR column types.
func (b *builder) checkKeyPartsOB(r tableRef, simpleCols []string, idxCols []*ast.IndexColumn, kind obKeyKind, opts []*ast.IndexOption, rng review.Range) bool {
	if b.indexOptionsFailOB(opts, kind, rng) {
		return true
	}
	seen := make(map[string]bool)
	checkDup := func(name string) bool {
		k, ok := b.key(spaceColumn, name)
		if !ok {
			return false
		}
		if seen[k] {
			b.fail(rng, nil, review.DuplicateName, statesTwice("names", "column", ident(name)))
			return true
		}
		seen[k] = true
		return false
	}
	if len(idxCols) > 0 {
		for _, ic := range idxCols {
			if ic == nil || ic.Functional {
				continue
			}
			if cr, ok := ic.Expr.(*ast.ColumnRef); ok && cr.Column != "" {
				if checkDup(cr.Column) {
					return true
				}
			}
		}
		for _, ic := range idxCols {
			if ic == nil {
				continue
			}
			if !ic.Functional {
				if cr, ok := ic.Expr.(*ast.ColumnRef); ok && cr.Column != "" {
					if b.checkKeyColumnOB(r, cr.Column, ic.Length, kind, rng) {
						return true
					}
					continue
				}
			}
			// Functional key part: check referenced columns exist, and reject bare-column ((a)).
			if col, ok := obBareColumnExpr(ic.Expr); ok {
				ck, okk := b.key(spaceColumn, col)
				if okk {
					colK := colKey(r.db, r.name, ck)
					if f, _ := b.l.Lookup(colK); f == walk.Absent {
						b.need(colK, walk.Present, nil, rng, notFound("column", ident(col), r.written), "")
						return true
					}
					b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses a functional index on bare column "+review.Code(ident(col))+".")
					return true
				}
			} else {
				for _, col := range obFreeColumnNames(ic.Expr) {
					if ck, ok := b.key(spaceColumn, col); ok {
						colK := colKey(r.db, r.name, ck)
						if f, _ := b.l.Lookup(colK); f == walk.Absent {
							b.need(colK, walk.Present, nil, rng, notFound("column", ident(col), r.written), "")
							return true
						}
					}
				}
			}
		}
	} else {
		for _, col := range simpleCols {
			if checkDup(col) {
				return true
			}
		}
		for _, col := range simpleCols {
			if b.checkKeyColumnOB(r, col, 0, kind, rng) {
				return true
			}
		}
	}
	for _, opt := range opts {
		if opt != nil && strings.EqualFold(opt.Name, "STORING") {
			for _, col := range opt.Columns {
				if ck, ok := b.key(spaceColumn, col); ok {
					b.need(colKey(r.db, r.name, ck), walk.Present, nil, rng, notFound("column", ident(col), r.written), "")
				}
			}
		}
	}
	return false
}

func (b *builder) checkKeyColumnOB(r tableRef, colName string, prefixLen int, kind obKeyKind, rng review.Range) bool {
	ck, ok := b.key(spaceColumn, colName)
	if !ok {
		return false
	}
	colK := colKey(r.db, r.name, ck)
	b.need(colK, walk.Present, nil, rng, notFound("column", ident(colName), r.written), "")
	if f, _ := b.l.Lookup(colK); f != walk.Present || b.touchedCols[ck] {
		return false
	}
	typ := b.w.colType(r, ck)
	if typ == "" {
		return false
	}
	if prefixLen > 0 && !obIsStringKeyType(typ) {
		b.fail(rng, []walk.Key{colK}, review.UnsupportedForm, b.engine()+" refuses a prefix key on non-string column "+review.Code(ident(colName))+".")
		return true
	}
	switch kind {
	case obKeyFulltext:
		if !obIsFulltextKeyType(typ) {
			b.fail(rng, []walk.Key{colK}, review.UnsupportedForm, b.engine()+" refuses a FULLTEXT index on non-string column "+review.Code(ident(colName))+".")
			return true
		}
	case obKeySpatial:
		if !obIsSpatialType(typ) {
			b.fail(rng, []walk.Key{colK}, review.UnsupportedForm, b.engine()+" refuses a SPATIAL index on non-spatial column "+review.Code(ident(colName))+".")
			return true
		}
	case obKeyVector:
		if typ != "vector" {
			b.fail(rng, []walk.Key{colK}, review.UnsupportedForm, b.engine()+" refuses a VECTOR index on non-vector column "+review.Code(ident(colName))+".")
			return true
		}
	}
	return false
}

// columnDefFailsOB extends columnFormFails with OceanBase-verified column
// definition checks: collection element types, generated column restrictions,
// invalid defaults and ON UPDATE expressions, SRID on non-spatial types, and
// non-boolean inline CHECK expressions.
func (b *builder) columnDefFailsOB(col *ast.ColumnDef, rng review.Range, inAlter bool) bool {
	if !b.w.oceanBase || col == nil || col.TypeName == nil {
		return false
	}
	if obBadCollectionElement(col.TypeName) {
		b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses this collection element type on column "+review.Code(ident(col.Name))+".")
		return true
	}
	if col.Generated != nil {
		for _, cc := range col.Constraints {
			if inAlter && cc.Type == ast.ColConstrCheck {
				b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses a column-level CHECK constraint on a generated column in ALTER TABLE.")
				return true
			}
			if !col.Generated.Stored && cc.Type == ast.ColConstrPrimaryKey {
				b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses a PRIMARY KEY on virtual generated column "+review.Code(ident(col.Name))+".")
				return true
			}
		}
	}
	baseType := obNormalizedDataType(col.TypeName)
	if col.OnUpdate != nil && obIsTemporalFunc(col.OnUpdate) && baseType != "datetime" && baseType != "timestamp" {
		b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses "+review.Code("ON UPDATE CURRENT_TIMESTAMP")+" on non-temporal column "+review.Code(ident(col.Name))+".")
		return true
	}
	defaults := []ast.ExprNode{col.DefaultValue}
	references := false
	for _, cc := range col.Constraints {
		switch cc.Type {
		case ast.ColConstrDefault:
			defaults = append(defaults, cc.Expr)
		case ast.ColConstrReferences:
			references = true
		case ast.ColConstrCheck:
			if !references && obNotCheckCondition(cc.Expr) {
				b.fail(rng, nil, review.UnsupportedForm, b.engine()+" requires a boolean expression for the CHECK constraint on column "+review.Code(ident(col.Name))+".")
				return true
			}
		}
	}
	for _, d := range defaults {
		if b.defaultExprFailsOB(col.Name, col.TypeName, baseType, d, rng) {
			return true
		}
	}
	if col.TypeName.HasSRID && !obIsSpatialType(baseType) {
		b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses "+review.Code("SRID")+" on non-spatial column "+review.Code(ident(col.Name))+".")
		return true
	}
	return false
}

func (b *builder) defaultExprFailsOB(colName string, dt *ast.DataType, baseType string, expr ast.ExprNode, rng review.Range) bool {
	if expr == nil {
		return false
	}
	if paren, ok := expr.(*ast.ParenExpr); ok {
		if len(obFreeColumnNames(paren)) > 0 {
			b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses column references in a default value expression.")
			return true
		}
		return false
	}
	if obIsTemporalFunc(expr) {
		if baseType != "" && baseType != "datetime" && baseType != "timestamp" {
			b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses a temporal default on column "+review.Code(ident(colName))+".")
			return true
		}
		return false
	}
	if baseType == "" {
		return false
	}
	if !obLiteralDefaultFits(dt, baseType, expr) {
		b.fail(rng, nil, review.UnsupportedForm, "Invalid default value for column "+review.Code(ident(colName))+".")
		return true
	}
	return false
}

func (b *builder) checkGeneratedColumnRefsOB(r tableRef, col *ast.ColumnDef, addedHere map[string]bool, rng review.Range) bool {
	if !b.w.oceanBase || col == nil || col.Generated == nil {
		return false
	}
	selfKey, _ := b.key(spaceColumn, col.Name)
	for _, ref := range obFreeColumnNames(col.Generated.Expr) {
		ck, ok := b.key(spaceColumn, ref)
		if !ok {
			continue
		}
		if ck == selfKey {
			b.fail(rng, nil, review.ObjectNotFound, notFound("column", ident(ref), r.written))
			return true
		}
		if addedHere[ck] {
			continue
		}
		colK := colKey(r.db, r.name, ck)
		b.need(colK, walk.Present, nil, rng, notFound("column", ident(ref), r.written), "")
	}
	return false
}

func (b *builder) checkAlterColumnDefaultOB(r tableRef, c *ast.AlterTableCmd, rng review.Range) bool {
	if !b.w.oceanBase || c.DefaultExpr == nil {
		return false
	}
	if len(obFreeColumnNames(c.DefaultExpr)) > 0 {
		b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses column references in a default value expression.")
		return true
	}
	ck, ok := b.key(spaceColumn, c.Name)
	if !ok {
		return false
	}
	colK := colKey(r.db, r.name, ck)
	if f, _ := b.l.Lookup(colK); f != walk.Present || b.touchedCols[ck] {
		return false
	}
	typ := b.w.colType(r, ck)
	if typ == "" {
		return false
	}
	if obIsTemporalFunc(c.DefaultExpr) && typ != "datetime" && typ != "timestamp" {
		b.fail(rng, []walk.Key{colK}, review.UnsupportedForm, b.engine()+" refuses a temporal default on column "+review.Code(ident(c.Name))+".")
		return true
	}
	// Only check simple numeric types where we do not need length/precision/enum members.
	switch typ {
	case "tinyint", "smallint", "mediumint", "int", "integer", "bigint", "float", "double", "real", "decimal", "numeric", "dec", "fixed":
		if text, _, ok := obLiteralText(c.DefaultExpr); ok {
			if _, ok := obLiteralNumber(text); !ok {
				b.fail(rng, []walk.Key{colK}, review.UnsupportedForm, "Invalid default value for column "+review.Code(ident(c.Name))+".")
				return true
			}
		}
	}
	return false
}

func (b *builder) checkTableCheckConstraintOB(r tableRef, c *ast.Constraint, rng review.Range) bool {
	if !b.w.oceanBase || c == nil || c.Expr == nil {
		return false
	}
	for _, col := range obFreeColumnNames(c.Expr) {
		ck, ok := b.key(spaceColumn, col)
		if !ok {
			continue
		}
		colK := colKey(r.db, r.name, ck)
		if f, _ := b.l.Lookup(colK); f == walk.Absent {
			b.need(colK, walk.Present, nil, rng, notFound("column", ident(col), r.written), "")
			return true
		}
	}
	if obNotCheckCondition(c.Expr) {
		b.fail(rng, nil, review.UnsupportedForm, b.engine()+" requires a boolean expression for a CHECK constraint.")
		return true
	}
	return false
}

// checkPartitionByOB checks OceanBase repartitioning rules on ALTER TABLE ...
// PARTITION BY that depend only on the clause itself, the table's columns, and
// its known primary key columns.
func (b *builder) checkPartitionByOB(r tableRef, pc *ast.PartitionClause, rng review.Range) bool {
	if pc == nil {
		return false
	}
	if pc.SplitSize != "" {
		b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses partition SIZE on ALTER TABLE.")
		return true
	}
	for _, names := range [][]string{pc.Columns, obFreeColumnNames(pc.Expr), pc.SubPartColumns, obFreeColumnNames(pc.SubPartExpr)} {
		for _, col := range names {
			ck, ok := b.key(spaceColumn, col)
			if !ok {
				continue
			}
			colK := colKey(r.db, r.name, ck)
			if f, _ := b.l.Lookup(colK); f == walk.Absent {
				b.need(colK, walk.Present, nil, rng, notFound("column", ident(col), r.written), "")
				return true
			}
		}
	}
	subpartitioned := pc.SubPartType != ast.PartitionRange || pc.SubPartExpr != nil || len(pc.SubPartColumns) > 0
	if !subpartitioned && (pc.Type == ast.PartitionHash || pc.Type == ast.PartitionKey) {
		for _, pd := range pc.Partitions {
			if pd != nil && len(pd.SubPartitions) > 0 {
				b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses subpartitions without SUBPARTITION BY.")
				return true
			}
		}
	}
	seen := make(map[string]bool, len(pc.Partitions))
	for _, pd := range pc.Partitions {
		if pd == nil || pd.Name == "" {
			continue
		}
		k := strings.ToLower(pd.Name)
		if seen[k] {
			b.fail(rng, nil, review.DuplicateName, statesTwice("declares", "partition", ident(pd.Name)))
			return true
		}
		seen[k] = true
	}
	var bounds []ast.Node
	for _, pd := range pc.Partitions {
		if pd == nil {
			continue
		}
		bounds = append(bounds, pd.Values)
		for _, sp := range pd.SubPartitions {
			if sp != nil {
				bounds = append(bounds, sp.Values)
			}
		}
	}
	for _, sp := range pc.SubPartTemplate {
		if sp != nil {
			bounds = append(bounds, sp.Values)
		}
	}
	for _, bound := range bounds {
		if col := obBoundColumn(bound); col != "" {
			b.fail(rng, nil, review.ObjectNotFound, notFound("column", ident(col), r.written))
			return true
		}
	}
	// Primary key prefix check for RANGE/LIST auto-repartition (no partition list),
	// and primary key coverage check (1503) for all partitioning types.
	pkCols := b.w.pkCols[relKey(r.db, r.name)]
	if len(pkCols) > 0 && !b.changesPK {
		if (pc.Type == ast.PartitionRange || pc.Type == ast.PartitionList) && len(pc.Partitions) == 0 {
			var key []string
			if len(pc.Columns) > 0 {
				key = pc.Columns
			} else if col, ok := obBareColumnExpr(pc.Expr); ok {
				key = []string{col}
			}
			if len(key) > 0 {
				matches := len(key) <= len(pkCols)
				for i := 0; matches && i < len(key); i++ {
					matches = strings.EqualFold(pkCols[i], key[i])
				}
				if !matches {
					b.fail(rng, nil, review.UnsupportedForm, b.engine()+" requires the auto-partitioning key to be a prefix of the primary key on "+review.Code(r.written)+".")
					return true
				}
			}
		}
		partCols := obPartitionKeyColumns(pc, pkCols)
		for _, pcol := range partCols {
			covered := false
			for _, pkc := range pkCols {
				if strings.EqualFold(pkc, pcol) {
					covered = true
					break
				}
			}
			if !covered {
				b.fail(rng, nil, review.UnsupportedForm, b.engine()+" requires the primary key of "+review.Code(r.written)+" to include partitioning column "+review.Code(ident(pcol))+".")
				return true
			}
		}
	}
	return false
}

func obPartitionKeyColumns(pc *ast.PartitionClause, pkCols []string) []string {
	var cols []string
	add := func(names []string) {
		for _, name := range names {
			found := false
			for _, existing := range cols {
				if strings.EqualFold(existing, name) {
					found = true
					break
				}
			}
			if !found {
				cols = append(cols, name)
			}
		}
	}
	switch {
	case len(pc.Columns) > 0:
		add(pc.Columns)
	case pc.Type == ast.PartitionKey:
		add(pkCols)
	default:
		add(obFreeColumnNames(pc.Expr))
	}
	switch {
	case len(pc.SubPartColumns) > 0:
		add(pc.SubPartColumns)
	case pc.SubPartType == ast.PartitionKey:
		add(pkCols)
	case pc.SubPartExpr != nil:
		add(obFreeColumnNames(pc.SubPartExpr))
	}
	return cols
}

// checkAlterPartitionCmdOB validates partition management subcommands whose
// statement text alone establishes refusal (duplicate partition names, column
// references in partition bounds, exchanging a partition with the same table).
func (b *builder) checkAlterPartitionCmdOB(r tableRef, c *ast.AlterTableCmd, rng review.Range) bool {
	switch c.Type {
	case ast.ATAddPartition:
		seen := make(map[string]bool, len(c.PartitionDefs))
		for _, pd := range c.PartitionDefs {
			if pd == nil {
				continue
			}
			if pd.Name != "" {
				k := strings.ToLower(pd.Name)
				if seen[k] {
					b.fail(rng, nil, review.DuplicateName, statesTwice("adds", "partition", ident(pd.Name)))
					return true
				}
				seen[k] = true
			}
			if col := obBoundColumn(pd.Values); col != "" {
				b.fail(rng, nil, review.ObjectNotFound, notFound("column", ident(col), r.written))
				return true
			}
		}
	case ast.ATAddSubpartition:
		for _, pd := range c.PartitionDefs {
			if pd == nil {
				continue
			}
			seen := make(map[string]bool, len(pd.SubPartitions))
			for _, sp := range pd.SubPartitions {
				if sp == nil {
					continue
				}
				if sp.Name != "" {
					k := strings.ToLower(sp.Name)
					if seen[k] {
						b.fail(rng, nil, review.DuplicateName, statesTwice("adds", "subpartition", ident(sp.Name)))
						return true
					}
					seen[k] = true
				}
				if col := obBoundColumn(sp.Values); col != "" {
					b.fail(rng, nil, review.ObjectNotFound, notFound("column", ident(col), r.written))
					return true
				}
			}
		}
	case ast.ATDropPartition:
		seen := make(map[string]bool, len(c.PartitionNames))
		for _, name := range c.PartitionNames {
			k := strings.ToLower(name)
			if seen[k] {
				b.fail(rng, nil, review.DuplicateName, statesTwice("drops", "partition", ident(name)))
				return true
			}
			seen[k] = true
		}
	case ast.ATExchangePartition:
		if c.ExchangeTable != nil {
			ex, st := b.ref(c.ExchangeTable)
			if st == refOK && ex.db == r.db && ex.name == r.name {
				b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses exchanging a partition of "+review.Code(r.written)+" with itself.")
				return true
			}
		}
	}
	return false
}

// checkForeignKeyColumnsOB validates foreign key local and referenced columns
// and checks that referenced columns lead an index on the parent table when
// that parent table's indexes are completely known.
func (b *builder) checkForeignKeyColumnsOB(r tableRef, c *ast.Constraint, rng review.Range) bool {
	if !b.w.oceanBase || !b.fks || c == nil || c.Type != ast.ConstrForeignKey {
		return false
	}
	for _, col := range c.Columns {
		if ck, ok := b.key(spaceColumn, col); ok {
			colK := colKey(r.db, r.name, ck)
			if f, _ := b.l.Lookup(colK); f == walk.Absent {
				b.need(colK, walk.Present, nil, rng, notFound("column", ident(col), r.written), "")
				return true
			}
		}
	}
	if b.w.fkChecks != review.SwitchOn || c.RefTable == nil || b.renames {
		return false
	}
	ref := *c.RefTable
	if ref.Schema == "" {
		if b.w.lost || b.w.noDatabase {
			return false
		}
		db := b.w.current
		if strings.HasPrefix(db, nonASCII) {
			return false
		}
		ref.Schema = db
	}
	t, st := b.ref(&ref)
	if st != refOK {
		return false
	}
	parentRel := relKey(t.db, t.name)
	if !(t.db == r.db && t.name == r.name) {
		f, kind := b.l.Lookup(parentRel)
		if f != walk.Present || kind != relTable {
			return false
		}
	}
	for _, col := range c.RefColumns {
		if ck, ok := b.key(spaceColumn, col); ok {
			colK := colKey(t.db, t.name, ck)
			if f, _ := b.l.Lookup(colK); f == walk.Absent {
				b.need(colK, walk.Present, nil, rng, notFound("column", ident(col), table(c.RefTable)), "")
				return true
			}
		}
	}
	if len(c.RefColumns) > 0 && b.w.indexLeadsKnown[parentRel] && !(t.db == r.db && t.name == r.name && b.addsIndex) {
		if !obHasLeadingIndex(b.w.indexLeads[parentRel], c.RefColumns) {
			b.fail(rng, nil, review.UnsupportedForm, b.engine()+" requires an index on referenced columns of "+review.Code(table(c.RefTable))+" for a foreign key.")
			return true
		}
	}
	return false
}

func obHasLeadingIndex(indexes [][]string, cols []string) bool {
	for _, idxCols := range indexes {
		if len(idxCols) < len(cols) {
			continue
		}
		match := true
		for i, c := range cols {
			if !strings.EqualFold(idxCols[i], c) {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// checkColumnUsedInFKOB checks OceanBase 5320: dropping a column used in a
// foreign key, or changing its data type via MODIFY/CHANGE COLUMN.
func (b *builder) checkColumnUsedInFKOB(r tableRef, colName string, newType *ast.DataType, isDrop bool, rng review.Range) bool {
	if !b.w.oceanBase || b.dropsFK {
		return false
	}
	ck, ok := b.key(spaceColumn, colName)
	if !ok {
		return false
	}
	rel := relKey(r.db, r.name)
	if !b.w.fkCols[rel][ck] {
		return false
	}
	if isDrop {
		b.fail(rng, []walk.Key{colKey(r.db, r.name, ck)}, review.UnsupportedForm, "Column "+review.Code(ident(colName))+" is used in a foreign key constraint on "+review.Code(r.written)+".")
		return true
	}
	if newType == nil {
		return false
	}
	oldBase := b.w.colType(r, ck)
	newBase := obNormalizedDataType(newType)
	if oldBase != "" && newBase != "" && oldBase != newBase {
		b.fail(rng, []walk.Key{colKey(r.db, r.name, ck)}, review.UnsupportedForm, "Column "+review.Code(ident(colName))+" is used in a foreign key constraint on "+review.Code(r.written)+".")
		return true
	}
	return false
}

// checkCreateTriggerOB validates OceanBase trigger rules beyond table/name
// checks: trigger schema matching its table's schema (1049/1435), FOLLOWS /
// PRECEDES referencing an existing trigger (3011), OLD in INSERT / NEW in
// DELETE (1363), and SET NEW in AFTER triggers (1362).
func (b *builder) checkCreateTriggerOB(r tableRef, v *ast.CreateTriggerStmt, rng review.Range) bool {
	if !b.w.oceanBase {
		return false
	}
	if v.Schema != "" {
		if trigDB, ok := b.key(spaceDatabase, v.Schema); ok && trigDB != r.db {
			b.fail(rng, nil, review.TriggerSchemaMismatch, "Trigger "+review.Code(ident(v.Schema)+"."+ident(v.Name))+" and table "+review.Code(r.written)+" are in different databases.")
			return true
		}
	}
	if v.Order != nil && v.Order.TriggerName != "" {
		if v.Order.Schema != "" {
			if orderDB, ok := b.key(spaceDatabase, v.Order.Schema); ok && orderDB != r.db {
				b.fail(rng, nil, review.ObjectNotFound, notFound("trigger", ident(v.Order.Schema)+"."+ident(v.Order.TriggerName), ""))
				return true
			}
		}
		if okKey, ok := b.key(spaceTrigger, v.Order.TriggerName); ok {
			b.need(objKey(spaceTrigger, r.db, okKey), walk.Present, nil, rng, notFound("trigger", ident(v.Order.TriggerName), ""), "")
		}
	}
	return b.checkTriggerBodyOB(v.Timing, v.Event, v.BodyText, rng)
}

func (b *builder) checkTriggerBodyOB(timing, event, body string, rng review.Range) bool {
	if strings.TrimSpace(body) == "" {
		return false
	}
	toks := parser.Tokenize(body)
	// Skip OLD./NEW. checks if the body binds old/new as a table name, implicit
	// or explicit alias, or local variable (any standalone identifier not
	// preceded or followed by '.').
	var oldBound, newBound bool
	for i := 0; i < len(toks); i++ {
		if toks[i].Type == obStringTokenType {
			continue
		}
		if i > 0 && toks[i-1].Type == '.' {
			continue
		}
		if i+1 < len(toks) && toks[i+1].Type == '.' {
			continue
		}
		switch strings.ToLower(toks[i].Str) {
		case "old":
			oldBound = true
		case "new":
			newBound = true
		}
	}
	lowerEvent := strings.ToLower(event)
	lowerTiming := strings.ToLower(timing)
	for i := 0; i+1 < len(toks); i++ {
		if toks[i].Type == obStringTokenType || toks[i+1].Type != '.' {
			continue
		}
		if i > 0 && toks[i-1].Type == '.' {
			continue
		}
		word := strings.ToLower(toks[i].Str)
		if word == "old" && !oldBound && lowerEvent == "insert" {
			b.fail(rng, nil, review.UnsupportedForm, b.engine()+" has no OLD row in an INSERT trigger.")
			return true
		}
		if word == "new" && !newBound && lowerEvent == "delete" {
			b.fail(rng, nil, review.UnsupportedForm, b.engine()+" has no NEW row in a DELETE trigger.")
			return true
		}
		if word == "new" && !newBound && lowerTiming == "after" && i > 0 && strings.EqualFold(toks[i-1].Str, "set") {
			b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses updating NEW row in an AFTER trigger.")
			return true
		}
	}
	return false
}

// checkEventScheduleOB checks that CREATE/ALTER EVENT schedule expressions
// (AT, STARTS, ENDS) do not reference bare columns (1525).
func (b *builder) checkEventScheduleOB(schedule *ast.EventSchedule, rng review.Range) bool {
	if !b.w.oceanBase || schedule == nil {
		return false
	}
	for _, expr := range []ast.ExprNode{schedule.At, schedule.Starts, schedule.Ends} {
		if col, ok := obFirstScheduleColRef(expr); ok {
			b.fail(rng, nil, review.UnsupportedForm, b.engine()+" refuses column reference "+review.Code(ident(col))+" in an event schedule.")
			return true
		}
	}
	return false
}

func obFirstScheduleColRef(expr ast.ExprNode) (string, bool) {
	switch e := expr.(type) {
	case nil:
		return "", false
	case *ast.ColumnRef:
		if e.Column != "" {
			return e.Column, true
		}
	case *ast.ParenExpr:
		return obFirstScheduleColRef(e.Expr)
	case *ast.UnaryExpr:
		return obFirstScheduleColRef(e.Operand)
	case *ast.BinaryExpr:
		if col, ok := obFirstScheduleColRef(e.Left); ok {
			return col, true
		}
		return obFirstScheduleColRef(e.Right)
	}
	return "", false
}

// --- Lexing, expression, type, charset, and collation helpers ---

var obStringTokenType = func() int {
	if toks := parser.Tokenize("''"); len(toks) == 1 {
		return toks[0].Type
	}
	return 0
}()

type obOptionToken struct {
	text  string
	str   bool
	punct byte
}

func lexOBOption(text string) []obOptionToken {
	var toks []obOptionToken
	for _, tok := range parser.Tokenize(text) {
		switch {
		case tok.Type == obStringTokenType:
			if n := len(toks) - 1; n >= 0 && toks[n].str {
				toks[n].text += tok.Str
				continue
			}
			toks = append(toks, obOptionToken{text: tok.Str, str: true})
		case tok.Type == '(' || tok.Type == ')' || tok.Type == ',' || tok.Type == '=':
			toks = append(toks, obOptionToken{text: tok.Str, punct: byte(tok.Type)})
		default:
			toks = append(toks, obOptionToken{text: tok.Str})
		}
	}
	return toks
}

func unquoteOBOption(value string) string {
	toks := lexOBOption(value)
	if len(toks) > 0 && toks[0].punct == '=' {
		toks = toks[1:]
	}
	if len(toks) == 1 {
		return toks[0].text
	}
	words := make([]string, len(toks))
	for i, t := range toks {
		words[i] = t.text
	}
	return strings.Join(words, " ")
}

func obOptionEntries(text string) ([][2]string, bool) {
	toks := lexOBOption(text)
	if len(toks) >= 2 && toks[0].punct == '(' && toks[len(toks)-1].punct == ')' {
		toks = toks[1 : len(toks)-1]
	}
	var entries [][2]string
	var entry []obOptionToken
	flush := func() bool {
		if len(entry) != 3 || entry[1].punct != '=' || entry[0].punct != 0 || entry[2].punct != 0 {
			return false
		}
		entries = append(entries, [2]string{entry[0].text, entry[2].text})
		entry = nil
		return true
	}
	for _, t := range toks {
		if t.punct == ',' {
			if !flush() {
				return nil, false
			}
			continue
		}
		entry = append(entry, t)
	}
	return entries, flush()
}

func obFreeColumnNames(root ast.Node) []string {
	if root == nil {
		return nil
	}
	var names []string
	var walkNode func(n ast.Node, masked map[string]bool)
	walkNode = func(n ast.Node, masked map[string]bool) {
		if n == nil {
			return
		}
		if lam, ok := n.(*ast.LambdaExpr); ok {
			next := make(map[string]bool, len(masked)+len(lam.Params))
			for k, v := range masked {
				next[k] = v
			}
			for _, p := range lam.Params {
				next[strings.ToLower(p)] = true
			}
			walkNode(lam.Body, next)
			return
		}
		if ref, ok := n.(*ast.ColumnRef); ok && ref.Column != "" {
			if ref.Table != "" || !masked[strings.ToLower(ref.Column)] {
				names = append(names, ref.Column)
			}
			return
		}
		ast.Inspect(n, func(child ast.Node) bool {
			if child == nil || child == n {
				return true
			}
			walkNode(child, masked)
			return false
		})
	}
	walkNode(root, nil)
	return names
}

func obBareColumnExpr(expr ast.ExprNode) (string, bool) {
	for {
		switch e := expr.(type) {
		case *ast.ParenExpr:
			expr = e.Expr
		case *ast.ColumnRef:
			if e.Column != "" {
				return e.Column, true
			}
			return "", false
		default:
			return "", false
		}
	}
}

func obBoundColumn(n ast.Node) string {
	switch v := n.(type) {
	case nil:
		return ""
	case *ast.List:
		for _, item := range v.Items {
			if name := obBoundColumn(item); name != "" {
				return name
			}
		}
		return ""
	case ast.ExprNode:
		if names := obFreeColumnNames(v); len(names) > 0 {
			return names[0]
		}
	}
	return ""
}

func obNotCheckCondition(expr ast.ExprNode) bool {
	switch e := expr.(type) {
	case *ast.ParenExpr:
		return obNotCheckCondition(e.Expr)
	case *ast.ColumnRef, *ast.StringLit:
		return true
	case *ast.BinaryExpr:
		switch e.Op {
		case ast.BinOpAdd, ast.BinOpSub, ast.BinOpMul, ast.BinOpDiv, ast.BinOpMod, ast.BinOpDivInt,
			ast.BinOpBitAnd, ast.BinOpBitOr, ast.BinOpBitXor, ast.BinOpShiftLeft, ast.BinOpShiftRight:
			return true
		}
	}
	return false
}

func obBaseTypeFromMetadata(colType string) string {
	s := strings.TrimSpace(colType)
	if idx := strings.IndexAny(s, " ("); idx >= 0 {
		s = s[:idx]
	}
	return obNormalizeBaseTypeName(s)
}

func obNormalizedDataType(dt *ast.DataType) string {
	if dt == nil {
		return ""
	}
	return obNormalizeBaseTypeName(dt.Name)
}

func obNormalizeBaseTypeName(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "integer", "int4", "int":
		return "int"
	case "int1", "tinyint":
		return "tinyint"
	case "int2", "smallint":
		return "smallint"
	case "int3", "middleint", "mediumint":
		return "mediumint"
	case "int8", "bigint":
		return "bigint"
	case "bool", "boolean":
		return "tinyint"
	case "character", "char", "nchar", "national char":
		return "char"
	case "character varying", "varchar", "nvarchar", "national varchar", "string":
		return "varchar"
	case "dec", "numeric", "fixed", "number", "decimal":
		return "decimal"
	case "float4", "float":
		return "float"
	case "float8", "double precision", "real", "double":
		return "double"
	case "long varbinary", "mediumblob":
		return "mediumblob"
	case "long varchar", "mediumtext":
		return "mediumtext"
	case "geomcollection":
		return "geometrycollection"
	case "serial":
		return "bigint"
	default:
		return strings.ToLower(strings.TrimSpace(name))
	}
}

func obIsStringKeyType(typ string) bool {
	switch typ {
	case "char", "varchar", "binary", "varbinary", "tinytext", "text", "mediumtext", "longtext",
		"tinyblob", "blob", "mediumblob", "longblob":
		return true
	}
	return false
}

func obIsFulltextKeyType(typ string) bool {
	switch typ {
	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext":
		return true
	}
	return false
}

func obIsSpatialType(typ string) bool {
	switch typ {
	case "geometry", "point", "linestring", "polygon", "multipoint", "multilinestring",
		"multipolygon", "geometrycollection", "geomcollection":
		return true
	}
	return false
}

func obBadCollectionElement(dt *ast.DataType) bool {
	for dt != nil {
		switch strings.ToLower(dt.Name) {
		case "map":
			if obBadCollectionElement(dt.MapKey) {
				return true
			}
			dt = dt.ArrayElem
			continue
		case "array":
			elem := dt.ArrayElem
			if elem == nil {
				return false
			}
			name := obNormalizeBaseTypeName(elem.Name)
			switch name {
			case "array", "map", "tinyint", "smallint", "mediumint", "int", "bigint", "boolean", "float", "double", "varchar", "vector":
			default:
				return true
			}
			dt = elem
			continue
		}
		return false
	}
	return false
}

var (
	obNumericTextRe = regexp.MustCompile(`^[+-]?(\d+\.?\d*|\.\d+)([eE][+-]?\d+)?$`)
	obDateTextRe    = regexp.MustCompile(`^(\d{2,4})[^\d](\d{1,2})[^\d](\d{1,2})(?:[ T](\d{1,2})[^\d](\d{1,2})(?:[^\d](\d{1,2})(?:\.\d+)?)?)?$`)
	obDateDigitsRe  = regexp.MustCompile(`^(\d{4})(\d{2})(\d{2})(?:(\d{2})(\d{2})(\d{2}))?$`)
	obTimeTextRe    = regexp.MustCompile(`^-?(?:\d+ )?\d{1,3}(?::(\d{1,2})(?::(\d{1,2})(?:\.\d+)?)?)?$`)
	obTimeDigitsRe  = regexp.MustCompile(`^-?\d+(?:\.\d+)?$`)
)

func obIsTemporalFunc(expr ast.ExprNode) bool {
	fn, ok := expr.(*ast.FuncCallExpr)
	if !ok {
		return false
	}
	switch strings.ToLower(fn.Name) {
	case "current_timestamp", "now", "localtime", "localtimestamp":
		return true
	}
	return false
}

func obLiteralText(expr ast.ExprNode) (string, bool, bool) {
	switch e := expr.(type) {
	case *ast.StringLit:
		return e.Value, true, true
	case *ast.IntLit:
		if e.Text != "" {
			return e.Text, false, obNumericTextRe.MatchString(e.Text)
		}
		return big.NewInt(e.Value).String(), false, true
	case *ast.FloatLit:
		return e.Value, false, true
	case *ast.UnaryExpr:
		if e.Op == ast.UnaryMinus {
			if text, isString, ok := obLiteralText(e.Operand); ok && !isString {
				return "-" + text, false, true
			}
		}
	}
	return "", false, false
}

func obLiteralNumber(text string) (*big.Float, bool) {
	text = strings.TrimSpace(text)
	if !obNumericTextRe.MatchString(text) {
		return nil, false
	}
	n, _, err := big.ParseFloat(text, 10, 256, big.ToNearestEven)
	return n, err == nil
}

func obLiteralDefaultFits(dt *ast.DataType, typeName string, expr ast.ExprNode) bool {
	text, isString, ok := obLiteralText(expr)
	if !ok {
		return true
	}
	switch typeName {
	case "tinyint", "smallint", "mediumint", "int", "bigint":
		n, ok := obLiteralNumber(text)
		if !ok {
			return false
		}
		if dt == nil {
			return true
		}
		return obIntInRange(n, typeName, dt.Unsigned || dt.Zerofill)
	case "decimal":
		n, ok := obLiteralNumber(text)
		if !ok {
			return false
		}
		if dt == nil {
			return true
		}
		precision, scale := dt.Length, dt.Scale
		if precision == 0 {
			precision = 10
		}
		if precision < 0 || precision > 65 || scale < 0 || scale > 30 || scale > precision {
			return true
		}
		limit := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(precision-scale)), nil))
		half := new(big.Float).Quo(big.NewFloat(0.5), new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(scale)), nil)))
		bound := new(big.Float).SetPrec(256).Sub(limit, half)
		return new(big.Float).Abs(n).Cmp(bound) < 0
	case "float", "double":
		_, ok := obLiteralNumber(text)
		return ok
	case "char", "binary":
		if !isString || dt == nil {
			return true
		}
		n := dt.Length
		if !dt.HasLength {
			n = 1
		}
		if typeName == "char" {
			return utf8.RuneCountInString(text) <= n
		}
		return len(text) <= n
	case "varchar":
		if !isString || dt == nil || dt.Length == 0 {
			return true
		}
		return utf8.RuneCountInString(text) <= dt.Length
	case "varbinary":
		if !isString || dt == nil || dt.Length == 0 {
			return true
		}
		return len(text) <= dt.Length
	case "date", "datetime", "timestamp":
		if !isString {
			return true
		}
		return obValidDateText(text)
	case "time":
		if !isString {
			return true
		}
		return obValidTimeText(text)
	}
	return true
}

func obIntInRange(n *big.Float, typeName string, unsigned bool) bool {
	bits := map[string]uint{"tinyint": 8, "smallint": 16, "mediumint": 24, "int": 32, "bigint": 64}[typeName]
	var lo, hi *big.Float
	if unsigned {
		lo = new(big.Float)
		hi = new(big.Float).SetInt(new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), bits), big.NewInt(1)))
	} else {
		hi = new(big.Float).SetInt(new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), bits-1), big.NewInt(1)))
		lo = new(big.Float).Neg(new(big.Float).SetInt(new(big.Int).Lsh(big.NewInt(1), bits-1)))
	}
	half := big.NewFloat(0.5)
	upper := new(big.Float).SetPrec(256).Add(hi, half)
	lower := new(big.Float).SetPrec(256).Sub(lo, half)
	return n.Cmp(upper) < 0 && n.Cmp(lower) > 0
}

func obValidDateText(text string) bool {
	var parts []string
	if m := obDateTextRe.FindStringSubmatch(text); m != nil {
		parts = m[1:]
	} else if m := obDateDigitsRe.FindStringSubmatch(text); m != nil {
		parts = m[1:]
	} else {
		return false
	}
	limits := []int{9999, 12, 31, 23, 59, 59}
	for i, p := range parts {
		if p == "" {
			continue
		}
		v := 0
		for _, c := range p {
			v = v*10 + int(c-'0')
		}
		if v > limits[i] {
			return false
		}
	}
	return true
}

func obValidTimeText(text string) bool {
	if obTimeDigitsRe.MatchString(text) {
		return true
	}
	m := obTimeTextRe.FindStringSubmatch(text)
	if m == nil {
		return false
	}
	for _, p := range m[1:] {
		if len(p) == 2 && p[0] >= '6' {
			return false
		}
	}
	return true
}

func (b *builder) checkCharsetNameOB(raw string, rng review.Range) bool {
	name := strings.ToLower(strings.Trim(strings.TrimSpace(raw), "'\"`"))
	if name == "" || name == "default" {
		return false
	}
	if !obCharsets[name] {
		b.fail(rng, nil, review.UnsupportedForm, b.engine()+" has no character set "+review.Code(ident(raw))+".")
		return true
	}
	return false
}

func (b *builder) checkCollationNameOB(raw string, rng review.Range) bool {
	name := obNormalizeCollationName(raw)
	if name == "" || name == "default" {
		return false
	}
	if _, ok := obCollations[name]; !ok {
		b.fail(rng, nil, review.UnsupportedForm, b.engine()+" has no collation "+review.Code(ident(raw))+".")
		return true
	}
	return false
}

func (b *builder) checkCharsetCollationPairOB(charset, collation string, rng review.Range) bool {
	cs := strings.ToLower(strings.Trim(strings.TrimSpace(charset), "'\"`"))
	if cs == "utf8" || cs == "utf8mb3" {
		cs = "utf8mb4"
	}
	col := obNormalizeCollationName(collation)
	if want, ok := obCollations[col]; ok && cs != "" && cs != "default" && cs != want {
		b.fail(rng, nil, review.UnsupportedForm, "Collation "+review.Code(ident(collation))+" is not valid for character set "+review.Code(ident(charset))+".")
		return true
	}
	return false
}

func obNormalizeCollationName(raw string) string {
	s := strings.ToLower(strings.Trim(strings.TrimSpace(raw), "'\"`"))
	if strings.HasPrefix(s, "utf8mb3_") {
		return "utf8mb4_" + s[len("utf8mb3_"):]
	}
	if strings.HasPrefix(s, "utf8_") {
		return "utf8mb4_" + s[len("utf8_"):]
	}
	return s
}

var obCharsets = map[string]bool{
	"binary": true, "utf8mb4": true, "utf8mb3": true, "utf8": true,
	"gbk": true, "utf16": true, "gb18030": true, "latin1": true,
	"gb2312": true, "gb18030_2022": true, "ascii": true, "tis620": true,
	"ujis": true, "euckr": true, "eucjpms": true, "cp932": true,
	"utf16le": true, "sjis": true, "big5": true, "hkscs": true,
	"hkscs31": true, "dec8": true, "cp850": true, "hp8": true,
	"macroman": true, "swe7": true,
}

var obCollations = map[string]string{
	"ascii_bin": "ascii", "ascii_general_ci": "ascii",
	"big5_bin": "big5", "big5_chinese_ci": "big5",
	"binary":    "binary",
	"cp850_bin": "cp850", "cp850_general_ci": "cp850",
	"cp932_bin": "cp932", "cp932_japanese_ci": "cp932",
	"dec8_bin": "dec8", "dec8_swedish_ci": "dec8",
	"eucjpms_bin": "eucjpms", "eucjpms_japanese_ci": "eucjpms",
	"euckr_bin": "euckr", "euckr_korean_ci": "euckr",
	"gb18030_2022_bin": "gb18030_2022", "gb18030_2022_chinese_ci": "gb18030_2022",
	"gb18030_2022_chinese_cs": "gb18030_2022", "gb18030_2022_radical_ci": "gb18030_2022",
	"gb18030_2022_radical_cs": "gb18030_2022", "gb18030_2022_stroke_ci": "gb18030_2022",
	"gb18030_2022_stroke_cs": "gb18030_2022",
	"gb18030_bin":            "gb18030", "gb18030_chinese_ci": "gb18030",
	"gb2312_bin": "gb2312", "gb2312_chinese_ci": "gb2312",
	"gbk_bin": "gbk", "gbk_chinese_ci": "gbk",
	"hkscs31_bin": "hkscs31", "hkscs_bin": "hkscs",
	"hp8_bin": "hp8", "hp8_english_ci": "hp8",
	"latin1_bin": "latin1", "latin1_danish_ci": "latin1",
	"latin1_general_ci": "latin1", "latin1_general_cs": "latin1",
	"latin1_german1_ci": "latin1", "latin1_german2_ci": "latin1",
	"latin1_spanish_ci": "latin1", "latin1_swedish_ci": "latin1",
	"macroman_bin": "macroman", "macroman_general_ci": "macroman",
	"sjis_bin": "sjis", "sjis_japanese_ci": "sjis",
	"swe7_bin": "swe7", "swe7_swedish_ci": "swe7",
	"tis620_bin": "tis620", "tis620_thai_ci": "tis620",
	"ujis_bin": "ujis", "ujis_japanese_ci": "ujis",
	"utf16_bin": "utf16", "utf16_croatian_ci": "utf16", "utf16_czech_ci": "utf16",
	"utf16_danish_ci": "utf16", "utf16_esperanto_ci": "utf16", "utf16_estonian_ci": "utf16",
	"utf16_general_ci": "utf16", "utf16_german2_ci": "utf16", "utf16_hungarian_ci": "utf16",
	"utf16_icelandic_ci": "utf16", "utf16_latvian_ci": "utf16", "utf16_lithuanian_ci": "utf16",
	"utf16_persian_ci": "utf16", "utf16_polish_ci": "utf16", "utf16_roman_ci": "utf16",
	"utf16_romanian_ci": "utf16", "utf16_sinhala_ci": "utf16", "utf16_slovak_ci": "utf16",
	"utf16_slovenian_ci": "utf16", "utf16_spanish2_ci": "utf16", "utf16_spanish_ci": "utf16",
	"utf16_swedish_ci": "utf16", "utf16_turkish_ci": "utf16", "utf16_unicode_520_ci": "utf16",
	"utf16_unicode_ci": "utf16", "utf16_vietnamese_ci": "utf16",
	"utf16le_bin": "utf16le", "utf16le_general_ci": "utf16le",
	"utf8mb4_0900_ai_ci": "utf8mb4", "utf8mb4_0900_as_ci": "utf8mb4",
	"utf8mb4_0900_as_cs": "utf8mb4", "utf8mb4_0900_bin": "utf8mb4",
	"utf8mb4_bg_0900_ai_ci": "utf8mb4", "utf8mb4_bg_0900_as_cs": "utf8mb4",
	"utf8mb4_bin": "utf8mb4", "utf8mb4_bs_0900_ai_ci": "utf8mb4",
	"utf8mb4_bs_0900_as_cs": "utf8mb4", "utf8mb4_croatian_ci": "utf8mb4",
	"utf8mb4_cs_0900_ai_ci": "utf8mb4", "utf8mb4_cs_0900_as_cs": "utf8mb4",
	"utf8mb4_czech_ci": "utf8mb4", "utf8mb4_da_0900_ai_ci": "utf8mb4",
	"utf8mb4_da_0900_as_cs": "utf8mb4", "utf8mb4_danish_ci": "utf8mb4",
	"utf8mb4_de_pb_0900_ai_ci": "utf8mb4", "utf8mb4_de_pb_0900_as_cs": "utf8mb4",
	"utf8mb4_eo_0900_ai_ci": "utf8mb4", "utf8mb4_eo_0900_as_cs": "utf8mb4",
	"utf8mb4_es_0900_ai_ci": "utf8mb4", "utf8mb4_es_0900_as_cs": "utf8mb4",
	"utf8mb4_es_trad_0900_ai_ci": "utf8mb4", "utf8mb4_es_trad_0900_as_cs": "utf8mb4",
	"utf8mb4_esperanto_ci": "utf8mb4", "utf8mb4_estonian_ci": "utf8mb4",
	"utf8mb4_et_0900_ai_ci": "utf8mb4", "utf8mb4_et_0900_as_cs": "utf8mb4",
	"utf8mb4_general_ci": "utf8mb4", "utf8mb4_german2_ci": "utf8mb4",
	"utf8mb4_gl_0900_ai_ci": "utf8mb4", "utf8mb4_gl_0900_as_cs": "utf8mb4",
	"utf8mb4_hr_0900_ai_ci": "utf8mb4", "utf8mb4_hr_0900_as_cs": "utf8mb4",
	"utf8mb4_hu_0900_ai_ci": "utf8mb4", "utf8mb4_hu_0900_as_cs": "utf8mb4",
	"utf8mb4_hungarian_ci": "utf8mb4", "utf8mb4_icelandic_ci": "utf8mb4",
	"utf8mb4_is_0900_ai_ci": "utf8mb4", "utf8mb4_is_0900_as_cs": "utf8mb4",
	"utf8mb4_ja_0900_as_cs": "utf8mb4", "utf8mb4_ja_0900_as_cs_ks": "utf8mb4",
	"utf8mb4_la_0900_ai_ci": "utf8mb4", "utf8mb4_la_0900_as_cs": "utf8mb4",
	"utf8mb4_latvian_ci": "utf8mb4", "utf8mb4_lithuanian_ci": "utf8mb4",
	"utf8mb4_lt_0900_ai_ci": "utf8mb4", "utf8mb4_lt_0900_as_cs": "utf8mb4",
	"utf8mb4_lv_0900_ai_ci": "utf8mb4", "utf8mb4_lv_0900_as_cs": "utf8mb4",
	"utf8mb4_mn_cyrl_0900_ai_ci": "utf8mb4", "utf8mb4_mn_cyrl_0900_as_cs": "utf8mb4",
	"utf8mb4_nb_0900_ai_ci": "utf8mb4", "utf8mb4_nb_0900_as_cs": "utf8mb4",
	"utf8mb4_nn_0900_ai_ci": "utf8mb4", "utf8mb4_nn_0900_as_cs": "utf8mb4",
	"utf8mb4_persian_ci": "utf8mb4", "utf8mb4_pl_0900_ai_ci": "utf8mb4",
	"utf8mb4_pl_0900_as_cs": "utf8mb4", "utf8mb4_polish_ci": "utf8mb4",
	"utf8mb4_ro_0900_ai_ci": "utf8mb4", "utf8mb4_ro_0900_as_cs": "utf8mb4",
	"utf8mb4_roman_ci": "utf8mb4", "utf8mb4_romanian_ci": "utf8mb4",
	"utf8mb4_ru_0900_ai_ci": "utf8mb4", "utf8mb4_ru_0900_as_cs": "utf8mb4",
	"utf8mb4_sinhala_ci": "utf8mb4", "utf8mb4_sk_0900_ai_ci": "utf8mb4",
	"utf8mb4_sk_0900_as_cs": "utf8mb4", "utf8mb4_sl_0900_ai_ci": "utf8mb4",
	"utf8mb4_sl_0900_as_cs": "utf8mb4", "utf8mb4_slovak_ci": "utf8mb4",
	"utf8mb4_slovenian_ci": "utf8mb4", "utf8mb4_spanish2_ci": "utf8mb4",
	"utf8mb4_spanish_ci": "utf8mb4", "utf8mb4_sr_latn_0900_ai_ci": "utf8mb4",
	"utf8mb4_sr_latn_0900_as_cs": "utf8mb4", "utf8mb4_sv_0900_ai_ci": "utf8mb4",
	"utf8mb4_sv_0900_as_cs": "utf8mb4", "utf8mb4_swedish_ci": "utf8mb4",
	"utf8mb4_tr_0900_ai_ci": "utf8mb4", "utf8mb4_tr_0900_as_cs": "utf8mb4",
	"utf8mb4_turkish_ci": "utf8mb4", "utf8mb4_unicode_520_ci": "utf8mb4",
	"utf8mb4_unicode_ci": "utf8mb4", "utf8mb4_vi_0900_ai_ci": "utf8mb4",
	"utf8mb4_vi_0900_as_cs": "utf8mb4", "utf8mb4_vietnamese_ci": "utf8mb4",
	"utf8mb4_zh_0900_as_cs": "utf8mb4",
}

// Whether a table is partitioned is one key of its path on OceanBase,
// read only when Present: partitionedNo or partitionedYes, as a CREATE
// TABLE of the change writes a PARTITION BY or not, and partitionedYes for
// a table the sync lists partitions of. A table the sync lists no
// partitions of is Unknown: oceanbase/driver's sync lists none. A synced
// table no statement of the change has repartitioned has a kind that says
// more (partitionFacts): its partitions' names, whether it has
// subpartitions, and the columns its partitioning reads.
const (
	spacePartitioned = "partitioned"
	spacePartition   = "partition"
	partitionedName  = "is"
	partitionedYes   = "yes"
	partitionedNo    = "no"
)

func partitionedKey(db, table string) walk.Key {
	return walk.Key{Space: spacePartitioned, Parent: tablePath(db, table), Name: partitionedName}
}

// partitionFacts is what the sync says of a partitioned table.
type partitionFacts struct {
	names []string // keys of the partitions' names
	sub   bool     // some partition has subpartitions
	cols  []string // keys of the columns the partitioning reads, nil when not known
}

// syncedPartitioned is the kind of the partitioned key of a table the sync
// lists partitions of.
func syncedPartitioned(parts []*metadata.TablePartitionMetadata) string {
	var f partitionFacts
	for _, p := range parts {
		k, ok := nameKey(spacePartition, p.GetName(), true)
		if !ok || strings.ContainsAny(k, "\x1e\x1f") {
			return partitionedYes
		}
		f.names = append(f.names, k)
		f.sub = f.sub || len(p.GetSubpartitions()) > 0
	}
	// The expression of a COLUMNS or KEY partitioning, or of a lone
	// column, is the quoted columns, comma-separated.
	for _, c := range strings.Split(parts[0].GetExpression(), ",") {
		c = strings.TrimSpace(c)
		if len(c) < 3 || c[0] != '`' || c[len(c)-1] != '`' || strings.ContainsAny(c[1:len(c)-1], "`\x1e\x1f") {
			f.cols = nil
			break
		}
		k, _ := nameKey(spaceColumn, c[1:len(c)-1], true)
		f.cols = append(f.cols, k)
	}
	sub := "0"
	if f.sub {
		sub = "1"
	}
	return strings.Join([]string{partitionedYes, sub, strings.Join(f.cols, "\x1e"), strings.Join(f.names, "\x1e")}, "\x1f")
}

// partitionFactsOf decodes a partitioned key's kind, false when it says no
// more than whether the table is partitioned.
func partitionFactsOf(kind string) (partitionFacts, bool) {
	fields := strings.Split(kind, "\x1f")
	if len(fields) != 4 {
		return partitionFacts{}, false
	}
	f := partitionFacts{sub: fields[1] == "1", names: strings.Split(fields[3], "\x1e")}
	if fields[2] != "" {
		f.cols = strings.Split(fields[2], "\x1e")
	}
	return f, true
}

// isPartitioned reports whether a partitioned key's kind says the table is
// partitioned.
func isPartitioned(kind string) bool {
	return kind == partitionedYes || strings.HasPrefix(kind, partitionedYes+"\x1f")
}

// partitionedFails fails an ALTER TABLE of one subcommand by what the
// table's partitioning refutes, each verified on 4.3.5 (the
// objects@partitioned and one_table@ fixtures of oceanbase/catalog's
// simulate.tsv), and before anything else the subcommand says:
//
//   - managing the partitions of a table that has none (1505): ADD, DROP,
//     TRUNCATE, REORGANIZE and EXCHANGE PARTITION, DROP, TRUNCATE and ADD
//     SUBPARTITION, and a partition's storage cache policy;
//   - managing the subpartitions of a table that has none (5875);
//   - DROP PARTITION of a partition the table doesn't have (1507), and ADD
//     PARTITION of one it has (1517);
//   - dropping or renaming a column the partitioning reads (3855).
//
// It records what the statement leaves: PARTITION BY a partitioned table,
// REMOVE PARTITIONING one that is not, and any other partition command
// partitions the walk no longer knows.
func (b *builder) partitionedFails(r tableRef, v *ast.AlterTableStmt) bool {
	k := partitionedKey(r.db, r.name)
	f, kind := b.l.Lookup(k)
	for _, c := range v.Commands {
		switch {
		case c == nil:
		case c.Type == ast.ATPartitionBy:
			b.Set(k, walk.Present, partitionedYes)
		case c.Type == ast.ATRemovePartitioning:
			b.Set(k, walk.Present, partitionedNo)
		case f == walk.Present && isPartitioned(kind) && !alterNoPartitionFacts(c):
			b.Set(k, walk.Present, partitionedYes)
		}
	}
	if len(v.Commands) != 1 || v.Commands[0] == nil || f != walk.Present {
		return false
	}
	c := v.Commands[0]
	rng := rangeOf(c.Loc)
	switch c.Type {
	case ast.ATAddPartition, ast.ATDropPartition, ast.ATTruncatePartition, ast.ATReorganizePartition,
		ast.ATExchangePartition, ast.ATDropSubpartition, ast.ATTruncateSubpartition, ast.ATAddSubpartition,
		ast.ATAlterPartitionStorageCachePolicy:
		if kind == partitionedNo {
			b.fail(rng, nil, review.UnsupportedForm, "Table "+review.Code(r.written)+" is not partitioned.")
			return true
		}
	}
	facts, ok := partitionFactsOf(kind)
	if !ok {
		return false
	}
	has := func(name string) bool {
		pk, ok := nameKey(spacePartition, name, true)
		return ok && slices.Contains(facts.names, pk)
	}
	switch c.Type {
	case ast.ATDropSubpartition, ast.ATTruncateSubpartition, ast.ATAddSubpartition:
		if !facts.sub {
			b.fail(rng, nil, review.UnsupportedForm, "Table "+review.Code(r.written)+" has no subpartitions.")
			return true
		}
	case ast.ATDropPartition:
		for _, n := range c.PartitionNames {
			if !has(n) {
				b.fail(rng, nil, review.ObjectNotFound, notFound("partition", ident(n), r.written))
				return true
			}
		}
	case ast.ATAddPartition:
		for _, pd := range c.PartitionDefs {
			if pd != nil && has(pd.Name) {
				b.fail(rng, nil, review.NameTaken, alreadyExists("partition", ident(pd.Name), r.written))
				return true
			}
		}
	case ast.ATDropColumn, ast.ATRenameColumn:
		if ck, ok := b.key(spaceColumn, c.Name); ok && slices.Contains(facts.cols, ck) {
			b.fail(rng, nil, review.UnsupportedForm, "Column "+review.Code(ident(c.Name))+" is read by the partitioning of "+review.Code(r.written)+", and can't be dropped or renamed.")
			return true
		}
	}
	return false
}

// alterNoPartitionFacts reports whether a subcommand leaves a table's
// partitions and partitioning columns as they were.
func alterNoPartitionFacts(c *ast.AlterTableCmd) bool {
	switch c.Type {
	case ast.ATAddPartition, ast.ATDropPartition, ast.ATReorganizePartition, ast.ATCoalescePartition,
		ast.ATExchangePartition, ast.ATDropSubpartition, ast.ATAddSubpartition, ast.ATChangeColumn,
		ast.ATRenameColumn, ast.ATDropColumn, ast.ATModifyColumn:
		return false
	}
	return true
}
