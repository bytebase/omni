package catalog

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/bytebase/omni/metadata"
)

func boolPtr(v bool) *bool { return &v }

func stringPtr(v string) *string { return &v }

// DatabaseMetadata converts the named database in the catalog to a
// DatabaseSchemaMetadata snapshot matching the MySQL driver sync and
// LoadMetadata round-trip contract.
func (c *Catalog) DatabaseMetadata(dbName string) *metadata.DatabaseSchemaMetadata {
	dbMeta := &metadata.DatabaseSchemaMetadata{
		Name: dbName,
	}
	if c == nil {
		return dbMeta
	}
	db := c.GetDatabase(dbName)
	if db == nil {
		return dbMeta
	}

	dbMeta.CharacterSet = db.Charset
	dbMeta.Collation = db.Collation

	// MySQL uses a single empty-name schema.
	schemaMeta := &metadata.SchemaMetadata{
		Name: "",
	}

	// tableProtoByLowerName lets the trigger pass below attach each trigger to its
	// owning table's metadata. The catalog keys tables case-insensitively, and a
	// trigger records its table via Trigger.Table, so the lookup is lower-cased.
	tableProtoByLowerName := make(map[string]*metadata.TableMetadata, len(db.Tables))
	for _, tName := range sortedMapKeys(db.Tables) {
		t := db.Tables[tName]
		if t.Temporary {
			continue
		}
		tableProto := tableToProto(db.Name, t)
		schemaMeta.Tables = append(schemaMeta.Tables, tableProto)
		tableProtoByLowerName[strings.ToLower(t.Name)] = tableProto
	}

	// Triggers are database-level in MySQL (db.Triggers), each bound to a table via
	// Trigger.Table; the proto nests them under their owning TableMetadata.Triggers
	// in ACTION_ORDER, matching live MySQL sync.
	triggers := make([]*Trigger, 0, len(db.Triggers))
	for _, tr := range db.Triggers {
		triggers = append(triggers, tr)
	}
	slices.SortFunc(triggers, func(a, b *Trigger) int {
		return cmp.Or(
			cmp.Compare(strings.ToUpper(a.Table), strings.ToUpper(b.Table)),
			cmp.Compare(strings.ToUpper(a.Event), strings.ToUpper(b.Event)),
			cmp.Compare(strings.ToUpper(a.Timing), strings.ToUpper(b.Timing)),
			cmp.Compare(a.ActionOrder, b.ActionOrder),
			cmp.Compare(a.Name, b.Name),
		)
	})
	for _, tr := range triggers {
		tableProto := tableProtoByLowerName[strings.ToLower(tr.Table)]
		if tableProto == nil {
			continue
		}
		tableProto.Triggers = append(tableProto.Triggers, triggerToProto(tr))
	}

	for _, vName := range sortedMapKeys(db.Views) {
		v := db.Views[vName]
		schemaMeta.Views = append(schemaMeta.Views, &metadata.ViewMetadata{
			Name:        v.Name,
			Definition:  v.Definition,
			Columns:     c.viewColumnsToProto(db, v),
			Algorithm:   v.Algorithm,
			SqlSecurity: v.SqlSecurity,
			CheckOption: v.CheckOption,
			SkipDump:    v.SkipDump,
		})
	}

	for _, fName := range sortedMapKeys(db.Functions) {
		f := db.Functions[fName]
		fnMeta := &metadata.FunctionMetadata{
			Name:              f.Name,
			Definition:        normalizeRoutineDefinition(c.ShowCreateFunction(dbName, f.Name), false),
			DatabaseCollation: db.Collation,
			Comment:           f.Characteristics["COMMENT"],
			SkipDump:          f.SkipDump,
		}
		if f.HasSessionContext {
			fnMeta.SqlMode = f.SQLMode
			fnMeta.CharacterSetClient = f.CharacterSetClient
			fnMeta.CollationConnection = f.CollationConnection
		}
		schemaMeta.Functions = append(schemaMeta.Functions, fnMeta)
	}

	for _, pName := range sortedMapKeys(db.Procedures) {
		p := db.Procedures[pName]
		procMeta := &metadata.ProcedureMetadata{
			Name:              p.Name,
			Definition:        normalizeRoutineDefinition(c.ShowCreateProcedure(dbName, p.Name), true),
			DatabaseCollation: db.Collation,
			Comment:           p.Characteristics["COMMENT"],
			SkipDump:          p.SkipDump,
		}
		if p.HasSessionContext {
			procMeta.SqlMode = p.SQLMode
			procMeta.CharacterSetClient = p.CharacterSetClient
			procMeta.CollationConnection = p.CollationConnection
		}
		schemaMeta.Procedures = append(schemaMeta.Procedures, procMeta)
	}

	for _, eName := range sortedMapKeys(db.Events) {
		e := db.Events[eName]
		evMeta := &metadata.EventMetadata{
			Name:       e.Name,
			Definition: c.ShowCreateEvent(dbName, e.Name),
			Comment:    e.Comment,
		}
		if e.HasSessionContext {
			evMeta.SqlMode = e.SQLMode
			evMeta.CharacterSetClient = e.CharacterSetClient
			evMeta.CollationConnection = e.CollationConnection
			evMeta.TimeZone = e.TimeZone
		}
		schemaMeta.Events = append(schemaMeta.Events, evMeta)
	}

	dbMeta.Schemas = append(dbMeta.Schemas, schemaMeta)
	return dbMeta
}

func (c *Catalog) viewColumnsToProto(db *Database, v *View) []*metadata.ColumnMetadata {
	leftQ := leftmostAnalyzedQuery(v.AnalyzedQuery)
	var nonJunkTargets []*TargetEntryQ
	if leftQ != nil {
		for _, te := range leftQ.TargetList {
			if te != nil && !te.ResJunk {
				nonJunkTargets = append(nonJunkTargets, te)
			}
		}
	}
	optionalRT := optionalJoinRangeIndices(leftQ)
	resolveTargetType := func(i int) string {
		if i < 0 || i >= len(nonJunkTargets) {
			return ""
		}
		te := nonJunkTargets[i]
		if te.ResOrigTable != "" && te.ResOrigCol != "" {
			origDB := db
			if te.ResOrigDB != "" {
				origDB = c.GetDatabase(te.ResOrigDB)
			}
			if origDB != nil {
				if tbl := origDB.GetTable(te.ResOrigTable); tbl != nil {
					if srcCol := tbl.GetColumn(te.ResOrigCol); srcCol != nil {
						return srcCol.ColumnType
					}
				}
			}
		}
		if te.Expr != nil {
			return resolveAnalyzedExprColumnType(c, db, leftQ, te.Expr)
		}
		return ""
	}
	resolveSourceAttrs := func(i int) (colType string, nullable bool, charset, collation string, ok bool) {
		if i < 0 || i >= len(nonJunkTargets) {
			return "", false, "", "", false
		}
		return resolveTargetColumnMetadata(c, db, leftQ, optionalRT, nonJunkTargets[i], true)
	}
	fillTargetCol := func(col *metadata.ColumnMetadata, i int, te *TargetEntryQ) {
		if v.AnalyzedQuery != nil && v.AnalyzedQuery.SetOp != SetOpNone {
			if meta, ok := resolveQueryTargetMeta(c, db, v.AnalyzedQuery, i, true); ok {
				if meta.colType != "" {
					col.Type = meta.colType
				}
				col.Nullable = meta.nullable
				col.CharacterSet = meta.charset
				col.Collation = meta.collation
				if te != nil && (col.Name == "" || col.Name == "*" || strings.HasSuffix(col.Name, ".*")) && te.ResName != "" {
					col.Name = te.ResName
				}
				return
			}
		}
		if colType, nullable, charset, collation, ok := resolveSourceAttrs(i); ok {
			if col.Type == "" {
				col.Type = colType
			}
			col.Nullable = nullable
			if charset != "" {
				col.CharacterSet = charset
			}
			if collation != "" {
				col.Collation = collation
			}
		} else if !col.Nullable && te != nil && te.Expr != nil && isAnalyzedExprNullable(c, db, leftQ, optionalRT, te.Expr, true) {
			col.Nullable = true
		}
		if col.Type == "" {
			col.Type = resolveTargetType(i)
		}
		if te != nil && (col.Name == "" || col.Name == "*" || strings.HasSuffix(col.Name, ".*")) && te.ResName != "" {
			col.Name = te.ResName
		}
	}
	if len(nonJunkTargets) > len(v.ColumnMetadata) && (!v.ExplicitColumns || len(v.Columns) == len(nonJunkTargets)) {
		cols := make([]*metadata.ColumnMetadata, 0, len(nonJunkTargets))
		for i, te := range nonJunkTargets {
			name := te.ResName
			if i < len(v.Columns) && v.Columns[i] != "" && v.Columns[i] != "*" && !strings.HasSuffix(v.Columns[i], ".*") {
				name = v.Columns[i]
			}
			col := &metadata.ColumnMetadata{
				Name:     name,
				Position: int32(i + 1),
			}
			if i < len(v.ColumnMetadata) && v.ColumnMetadata[i].Name == name {
				cm := v.ColumnMetadata[i]
				col.Nullable = cm.Nullable
				col.CharacterSet = cm.Charset
				col.Collation = cm.Collation
				if cm.Type != nil {
					col.Type = resolvedTypeToColumnType(cm.Type)
				} else if cm.ColumnType != "" {
					col.Type = cm.ColumnType
				}
			}
			fillTargetCol(col, i, te)
			cols = append(cols, col)
		}
		return cols
	}
	if len(v.ColumnMetadata) > 0 {
		cols := make([]*metadata.ColumnMetadata, 0, len(v.ColumnMetadata))
		for i, cm := range v.ColumnMetadata {
			col := &metadata.ColumnMetadata{
				Name:         cm.Name,
				Position:     int32(i + 1),
				Nullable:     cm.Nullable,
				CharacterSet: cm.Charset,
				Collation:    cm.Collation,
			}
			if cm.Type != nil {
				col.Type = resolvedTypeToColumnType(cm.Type)
			} else if cm.ColumnType != "" {
				col.Type = cm.ColumnType
			}
			var te *TargetEntryQ
			if i < len(nonJunkTargets) {
				te = nonJunkTargets[i]
			}
			fillTargetCol(col, i, te)
			cols = append(cols, col)
		}
		return cols
	}
	if len(v.Columns) == 0 {
		return nil
	}
	cols := make([]*metadata.ColumnMetadata, 0, len(v.Columns))
	for i, name := range v.Columns {
		col := &metadata.ColumnMetadata{
			Name:     name,
			Position: int32(i + 1),
		}
		var te *TargetEntryQ
		if i < len(nonJunkTargets) {
			te = nonJunkTargets[i]
		}
		fillTargetCol(col, i, te)
		cols = append(cols, col)
	}
	return cols
}

func leftmostAnalyzedQuery(q *Query) *Query {
	for q != nil && q.SetOp != SetOpNone && q.LArg != nil {
		q = q.LArg
	}
	return q
}

type viewTargetMeta struct {
	colType   string
	nullable  bool
	charset   string
	collation string
	expr      AnalyzedExpr
}

func resolveQueryTargetMeta(c *Catalog, db *Database, q *Query, idx int, widenOptionalJoin bool) (viewTargetMeta, bool) {
	if q == nil {
		return viewTargetMeta{}, false
	}
	if q.SetOp == SetOpNone {
		te := nthNonJunkTarget(q, idx)
		if te == nil {
			return viewTargetMeta{}, false
		}
		var opt map[int]bool
		if widenOptionalJoin {
			opt = optionalJoinRangeIndices(q)
		}
		meta := viewTargetMeta{expr: te.Expr}
		if ct, nl, cs, cl, ok := resolveTargetColumnMetadata(c, db, q, opt, te, widenOptionalJoin); ok {
			meta.colType = ct
			meta.nullable = nl
			meta.charset = cs
			meta.collation = cl
		} else {
			meta.nullable = isAnalyzedExprNullable(c, db, q, opt, te.Expr, widenOptionalJoin)
		}
		if meta.colType == "" {
			if te.ResOrigTable != "" && te.ResOrigCol != "" {
				origDB := db
				if te.ResOrigDB != "" && c != nil {
					origDB = c.GetDatabase(te.ResOrigDB)
				}
				if origDB != nil {
					if tbl := origDB.GetTable(te.ResOrigTable); tbl != nil {
						if srcCol := tbl.GetColumn(te.ResOrigCol); srcCol != nil {
							meta.colType = srcCol.ColumnType
						}
					}
				}
			}
			if meta.colType == "" && te.Expr != nil {
				meta.colType = resolveAnalyzedExprColumnType(c, db, q, te.Expr)
			}
		}
		if meta.charset == "" {
			if op, ok := te.Expr.(*OpExprQ); ok && op.Op == "->>" {
				meta.charset = "utf8mb4"
				meta.collation = "utf8mb4_bin"
			} else if isStringColumnType(meta.colType) && db != nil && db.Charset != "" {
				meta.charset = db.Charset
				meta.collation = db.Collation
			}
		}
		return meta, true
	}

	leftMeta, lok := resolveQueryTargetMeta(c, db, q.LArg, idx, widenOptionalJoin)
	rightMeta, rok := resolveQueryTargetMeta(c, db, q.RArg, idx, widenOptionalJoin)
	if !lok && !rok {
		return viewTargetMeta{}, false
	}
	if !lok {
		return rightMeta, true
	}
	if !rok {
		return leftMeta, true
	}

	combinedType := combineSetOpColumnTypes(leftMeta.colType, leftMeta.expr, rightMeta.colType, rightMeta.expr, false)
	var nullable bool
	switch q.SetOp {
	case SetOpUnion:
		nullable = leftMeta.nullable || rightMeta.nullable
	case SetOpExcept:
		nullable = leftMeta.nullable
	case SetOpIntersect:
		nullable = leftMeta.nullable && rightMeta.nullable
	}
	var cs, cl string
	if isStringColumnType(combinedType) {
		cs, cl = combineSetOpCollation(db, leftMeta.charset, leftMeta.collation, rightMeta.charset, rightMeta.collation)
	}
	return viewTargetMeta{
		colType:   combinedType,
		nullable:  nullable,
		charset:   cs,
		collation: cl,
	}, true
}

func isAnalyzedExprNullable(c *Catalog, db *Database, q *Query, opt map[int]bool, expr AnalyzedExpr, widenOptionalJoin bool) bool {
	if expr == nil {
		return false
	}
	switch e := expr.(type) {
	case *ConstExprQ:
		return e.IsNull
	case *VarExprQ:
		if widenOptionalJoin && opt[e.RangeIdx] {
			return true
		}
		if q != nil && e.LevelsUp == 0 && e.RangeIdx >= 0 && e.RangeIdx < len(q.RangeTable) {
			rte := q.RangeTable[e.RangeIdx]
			if rte.Kind == RTERelation && c != nil && e.AttNum >= 1 && e.AttNum <= len(rte.ColNames) {
				origDB := db
				if rte.DBName != "" {
					origDB = c.GetDatabase(rte.DBName)
				}
				if origDB != nil {
					if tbl := origDB.GetTable(rte.TableName); tbl != nil {
						if col := tbl.GetColumn(rte.ColNames[e.AttNum-1]); col != nil {
							return col.Nullable
						}
					}
					if sv := origDB.Views[strings.ToLower(rte.TableName)]; sv != nil && sv.AnalyzedQuery != nil {
						if meta, ok := resolveQueryTargetMeta(c, origDB, sv.AnalyzedQuery, e.AttNum-1, widenOptionalJoin); ok {
							return meta.nullable
						}
					}
				}
			}
		}
		return false
	case *CastExprQ:
		return isAnalyzedExprNullable(c, db, q, opt, e.Arg, widenOptionalJoin)
	case *OpExprQ:
		if _, lok := e.Left.(*IntervalExprQ); lok {
			return true
		}
		if _, rok := e.Right.(*IntervalExprQ); rok {
			return true
		}
		if e.Op == "/" || strings.EqualFold(e.Op, "DIV") {
			return true
		}
		switch strings.ToUpper(e.Op) {
		case "<=>", "IS TRUE", "IS NOT TRUE", "IS FALSE", "IS NOT FALSE", "IS UNKNOWN", "IS NOT UNKNOWN":
			return false
		case "->", "->>", "BINARY":
			return true
		}
		return isAnalyzedExprNullable(c, db, q, opt, e.Left, widenOptionalJoin) || isAnalyzedExprNullable(c, db, q, opt, e.Right, widenOptionalJoin)
	case *NullTestExprQ:
		return false
	case *BetweenExprQ:
		return isAnalyzedExprNullable(c, db, q, opt, e.Arg, widenOptionalJoin) ||
			isAnalyzedExprNullable(c, db, q, opt, e.Lower, widenOptionalJoin) ||
			isAnalyzedExprNullable(c, db, q, opt, e.Upper, widenOptionalJoin)
	case *LikeExprQ:
		return isAnalyzedExprNullable(c, db, q, opt, e.Arg, widenOptionalJoin) ||
			isAnalyzedExprNullable(c, db, q, opt, e.Pattern, widenOptionalJoin) ||
			(e.Escape != nil && isAnalyzedExprNullable(c, db, q, opt, e.Escape, widenOptionalJoin))
	case *InListExprQ:
		if isAnalyzedExprNullable(c, db, q, opt, e.Arg, widenOptionalJoin) {
			return true
		}
		for _, item := range e.List {
			if isAnalyzedExprNullable(c, db, q, opt, item, widenOptionalJoin) {
				return true
			}
		}
		return false
	case *BoolExprQ:
		for _, arg := range e.Args {
			if isAnalyzedExprNullable(c, db, q, opt, arg, widenOptionalJoin) {
				return true
			}
		}
		return false
	case *IntervalExprQ:
		return isAnalyzedExprNullable(c, db, q, opt, e.Value, widenOptionalJoin)
	case *SubLinkExprQ:
		switch e.Kind {
		case SubLinkExists:
			return false
		case SubLinkScalar, SubLinkIn, SubLinkAny, SubLinkAll:
			return true
		}
		return false
	case *FuncCallExprQ:
		if strings.EqualFold(e.Name, "count") {
			return false
		}
		switch strings.ToLower(e.Name) {
		case "nullif", "mod",
			"date_add", "date_sub", "adddate", "subdate",
			"timestampadd", "timestampdiff", "datediff", "timediff", "addtime", "subtime":
			return true
		}
		if strings.EqualFold(e.Name, "if") && len(e.Args) == 3 {
			return isAnalyzedExprNullable(c, db, q, opt, e.Args[1], widenOptionalJoin) ||
				isAnalyzedExprNullable(c, db, q, opt, e.Args[2], widenOptionalJoin)
		}
		if isViewStringMetadataNullableFunc(e.Name) || isViewNullableAggregateFunc(e.Name) {
			return true
		}
		for _, arg := range e.Args {
			if isAnalyzedExprNullable(c, db, q, opt, arg, widenOptionalJoin) {
				return true
			}
		}
		return false
	case *CoalesceExprQ:
		if len(e.Args) == 0 {
			return true
		}
		for _, arg := range e.Args {
			if !isAnalyzedExprNullable(c, db, q, opt, arg, widenOptionalJoin) {
				return false
			}
		}
		return true
	case *CaseExprQ:
		if e.Default == nil || isAnalyzedExprNullable(c, db, q, opt, e.Default, widenOptionalJoin) {
			return true
		}
		for _, w := range e.Args {
			if w != nil && isAnalyzedExprNullable(c, db, q, opt, w.Then, widenOptionalJoin) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

func combineSetOpCollation(db *Database, lcs, lcl, rcs, rcl string) (string, string) {
	if lcs == "" && rcs == "" {
		if db != nil && db.Charset != "" {
			return db.Charset, db.Collation
		}
		return "", ""
	}
	if lcs == "" {
		return rcs, rcl
	}
	if rcs == "" {
		return lcs, lcl
	}
	if strings.EqualFold(lcl, rcl) {
		return lcs, lcl
	}
	if strings.HasSuffix(strings.ToLower(lcl), "_bin") {
		return lcs, lcl
	}
	if strings.HasSuffix(strings.ToLower(rcl), "_bin") {
		return rcs, rcl
	}
	return lcs, lcl
}

type parsedSetOpType struct {
	cat      string // null, int, float, double, decimal, string, text, binary, blob, json, date, time, datetime, timestamp, year, bit, other
	raw      string
	rank     int
	width    int
	prec     int
	scale    int
	unsigned bool
	isChar   bool
}

func isStringColumnType(colType string) bool {
	p := parseSetOpColType(colType)
	return p.cat == "string" || p.cat == "text"
}

func parseSetOpColType(colType string) parsedSetOpType {
	raw := strings.TrimSpace(colType)
	s := strings.ToLower(raw)
	if s == "" || s == "binary(0)" {
		return parsedSetOpType{cat: "null", raw: raw}
	}
	unsigned := strings.Contains(s, "unsigned")
	s = strings.ReplaceAll(s, " unsigned", "")
	s = strings.ReplaceAll(s, " zerofill", "")
	s = strings.TrimSpace(s)

	base := s
	var p, sc int
	if open := strings.IndexByte(s, '('); open >= 0 {
		base = strings.TrimSpace(s[:open])
		if closeIdx := strings.IndexByte(s[open+1:], ')'); closeIdx >= 0 {
			args := strings.Split(s[open+1:open+1+closeIdx], ",")
			if len(args) >= 1 {
				p, _ = strconv.Atoi(strings.TrimSpace(args[0]))
			}
			if len(args) >= 2 {
				sc, _ = strconv.Atoi(strings.TrimSpace(args[1]))
			}
		}
	}

	switch base {
	case "tinyint", "bool", "boolean":
		return parsedSetOpType{cat: "int", raw: raw, rank: 1, width: 4, unsigned: unsigned}
	case "smallint":
		return parsedSetOpType{cat: "int", raw: raw, rank: 2, width: 6, unsigned: unsigned}
	case "mediumint":
		return parsedSetOpType{cat: "int", raw: raw, rank: 3, width: 9, unsigned: unsigned}
	case "int", "integer":
		return parsedSetOpType{cat: "int", raw: raw, rank: 4, width: 11, unsigned: unsigned}
	case "bigint":
		return parsedSetOpType{cat: "int", raw: raw, rank: 5, width: 20, unsigned: unsigned}
	case "float":
		return parsedSetOpType{cat: "float", raw: raw, width: 12, unsigned: unsigned}
	case "double", "real":
		return parsedSetOpType{cat: "double", raw: raw, width: 22, unsigned: unsigned}
	case "decimal", "numeric", "dec", "fixed":
		if p == 0 {
			p = 10
		}
		w := p + 1
		if sc > 0 {
			w = p + 2
		}
		return parsedSetOpType{cat: "decimal", raw: raw, prec: p, scale: sc, width: w, unsigned: unsigned}
	case "char":
		return parsedSetOpType{cat: "string", raw: raw, width: p, isChar: true}
	case "varchar", "var_string":
		return parsedSetOpType{cat: "string", raw: raw, width: p}
	case "tinytext":
		return parsedSetOpType{cat: "text", raw: raw, rank: 1, width: 255}
	case "text":
		return parsedSetOpType{cat: "text", raw: raw, rank: 2, width: 65535}
	case "mediumtext":
		return parsedSetOpType{cat: "text", raw: raw, rank: 3, width: 16777215}
	case "longtext":
		return parsedSetOpType{cat: "text", raw: raw, rank: 4, width: 4294967295}
	case "binary":
		return parsedSetOpType{cat: "binary", raw: raw, width: p, isChar: true}
	case "varbinary":
		return parsedSetOpType{cat: "binary", raw: raw, width: p}
	case "tinyblob":
		return parsedSetOpType{cat: "blob", raw: raw, rank: 1, width: 255}
	case "blob":
		return parsedSetOpType{cat: "blob", raw: raw, rank: 2, width: 65535}
	case "mediumblob":
		return parsedSetOpType{cat: "blob", raw: raw, rank: 3, width: 16777215}
	case "longblob":
		return parsedSetOpType{cat: "blob", raw: raw, rank: 4, width: 4294967295}
	case "json":
		return parsedSetOpType{cat: "json", raw: raw, width: 4294967295}
	case "date":
		return parsedSetOpType{cat: "date", raw: raw, width: 10}
	case "time":
		w := 8
		if p > 0 {
			w += p + 1
		}
		return parsedSetOpType{cat: "time", raw: raw, prec: p, width: w}
	case "datetime", "timestamp":
		w := 19
		if p > 0 {
			w += p + 1
		}
		return parsedSetOpType{cat: base, raw: raw, prec: p, width: w}
	case "year":
		return parsedSetOpType{cat: "year", raw: raw, width: 4}
	case "bit":
		w := p
		if w == 0 {
			w = 1
		}
		return parsedSetOpType{cat: "bit", raw: raw, prec: w, width: w}
	default:
		return parsedSetOpType{cat: "other", raw: raw}
	}
}

func setOpStringWidth(pt parsedSetOpType, expr AnalyzedExpr, isTiDB bool) int {
	if !isTiDB && pt.cat == "int" {
		if cn, ok := expr.(*ConstExprQ); ok && !cn.IsNull && !cn.IsString {
			digits := len(strings.TrimLeft(cn.Value, "+-"))
			if digits == 0 {
				digits = 1
			}
			return digits + 1
		}
	}
	return pt.width
}

func combineSetOpColumnTypes(lt string, le AnalyzedExpr, rt string, re AnalyzedExpr, isTiDB bool) string {
	lp := parseSetOpColType(lt)
	rp := parseSetOpColType(rt)
	if lp.cat == "null" {
		return rt
	}
	if rp.cat == "null" {
		return lt
	}
	if lp.cat == "json" || rp.cat == "json" {
		if lp.cat == "json" && rp.cat == "json" {
			return "json"
		}
		return "longtext"
	}
	isStrLike := func(p parsedSetOpType) bool {
		return p.cat == "string" || p.cat == "text"
	}
	isBinLike := func(p parsedSetOpType) bool {
		return p.cat == "binary" || p.cat == "blob"
	}
	if isBinLike(lp) || isBinLike(rp) {
		if lp.cat == "blob" || rp.cat == "blob" {
			switch max(lp.rank, rp.rank) {
			case 1:
				return "tinyblob"
			case 2:
				return "blob"
			case 3:
				return "mediumblob"
			default:
				return "longblob"
			}
		}
		w := max(setOpStringWidth(lp, le, isTiDB), setOpStringWidth(rp, re, isTiDB))
		if lp.cat == "binary" && rp.cat == "binary" && lp.isChar && rp.isChar {
			return fmt.Sprintf("binary(%d)", w)
		}
		return fmt.Sprintf("varbinary(%d)", w)
	}
	if isStrLike(lp) || isStrLike(rp) {
		if lp.cat == "text" || rp.cat == "text" {
			switch max(lp.rank, rp.rank) {
			case 1:
				return "tinytext"
			case 2:
				return "text"
			case 3:
				return "mediumtext"
			default:
				return "longtext"
			}
		}
		w := max(setOpStringWidth(lp, le, isTiDB), setOpStringWidth(rp, re, isTiDB))
		if lp.cat == "string" && rp.cat == "string" && lp.isChar && rp.isChar {
			return fmt.Sprintf("char(%d)", w)
		}
		return fmt.Sprintf("varchar(%d)", w)
	}
	isNumeric := func(p parsedSetOpType) bool {
		return p.cat == "int" || p.cat == "float" || p.cat == "double" || p.cat == "decimal" || p.cat == "year" || p.cat == "bit"
	}
	if isNumeric(lp) && isNumeric(rp) {
		if lp.cat == "double" || rp.cat == "double" {
			return "double"
		}
		if lp.cat == "float" || rp.cat == "float" {
			return "float"
		}
		if lp.cat == "decimal" || rp.cat == "decimal" {
			lrt := parseColumnResolvedType(lt)
			rrt := parseColumnResolvedType(rt)
			if lrt == nil {
				lrt = &ResolvedType{BaseType: BaseTypeBigInt}
			}
			if rrt == nil {
				rrt = &ResolvedType{BaseType: BaseTypeBigInt}
			}
			lPrec, lScale, lok := decimalOperandPrecisionScale(lrt, le)
			rPrec, rScale, rok := decimalOperandPrecisionScale(rrt, re)
			if !lok {
				lPrec, lScale = lp.prec, lp.scale
			}
			if !rok {
				rPrec, rScale = rp.prec, rp.scale
			}
			if isTiDB {
				// TiDB widens integer operands to DECIMAL(20,...) when coercing with DECIMAL.
				if lp.cat == "int" && lPrec < 20 {
					lPrec = 20
				}
				if rp.cat == "int" && rPrec < 20 {
					rPrec = 20
				}
			}
			scale := max(lScale, rScale)
			prec := min(65, max(lPrec-lScale, rPrec-rScale)+scale)
			return fmt.Sprintf("decimal(%d,%d)", prec, scale)
		}
		if lp.cat == "bit" && rp.cat == "bit" {
			return fmt.Sprintf("bit(%d)", max(lp.prec, rp.prec))
		}
		if lp.cat == "year" && rp.cat == "year" {
			return "year"
		}
		rank := max(lp.rank, rp.rank)
		if rank == 0 {
			rank = 4
		}
		unsigned := false
		if lp.unsigned == rp.unsigned {
			unsigned = lp.unsigned
		} else {
			var sRank, uRank int
			if lp.unsigned {
				uRank, sRank = lp.rank, rp.rank
			} else {
				sRank, uRank = lp.rank, rp.rank
			}
			if sRank > uRank {
				rank = sRank
			} else if uRank >= 5 {
				rank = 5
				unsigned = true
			} else {
				rank = uRank + 1
			}
		}
		suffix := ""
		if unsigned {
			suffix = " unsigned"
		}
		switch rank {
		case 1:
			return "tinyint" + suffix
		case 2:
			return "smallint" + suffix
		case 3:
			return "mediumint" + suffix
		case 4:
			return "int" + suffix
		default:
			return "bigint" + suffix
		}
	}
	if lp.cat == "date" && rp.cat == "date" {
		return "date"
	}
	if (lp.cat == "datetime" || lp.cat == "timestamp" || lp.cat == "date") &&
		(rp.cat == "datetime" || rp.cat == "timestamp" || rp.cat == "date") {
		if lp.cat == "timestamp" && rp.cat == "timestamp" {
			fsp := max(lp.prec, rp.prec)
			if fsp > 0 {
				return fmt.Sprintf("timestamp(%d)", fsp)
			}
			return "timestamp"
		}
		fsp := max(lp.prec, rp.prec)
		if fsp > 0 {
			return fmt.Sprintf("datetime(%d)", fsp)
		}
		return "datetime"
	}
	if lp.cat == "time" && rp.cat == "time" {
		fsp := max(lp.prec, rp.prec)
		if fsp > 0 {
			return fmt.Sprintf("time(%d)", fsp)
		}
		return "time"
	}
	if lp.raw == rp.raw {
		return lp.raw
	}
	return fmt.Sprintf("varchar(%d)", max(setOpStringWidth(lp, le, isTiDB), setOpStringWidth(rp, re, isTiDB)))
}

func optionalJoinRangeIndices(q *Query) map[int]bool {
	if q == nil || q.JoinTree == nil {
		return nil
	}
	out := make(map[int]bool)
	for _, node := range q.JoinTree.FromList {
		markOptionalJoinNodes(node, false, out)
	}
	return out
}

func markOptionalJoinNodes(node JoinNode, optional bool, out map[int]bool) {
	switch n := node.(type) {
	case *RangeTableRefQ:
		if optional {
			out[n.RTIndex] = true
		}
	case *JoinExprNodeQ:
		if optional {
			out[n.RTIndex] = true
		}
		leftOpt := optional || n.JoinType == JoinRight
		rightOpt := optional || n.JoinType == JoinLeft
		markOptionalJoinNodes(n.Left, leftOpt, out)
		markOptionalJoinNodes(n.Right, rightOpt, out)
	}
}

func exprUsesOptionalJoin(expr AnalyzedExpr, optionalRT map[int]bool) bool {
	if expr == nil || len(optionalRT) == 0 {
		return false
	}
	if ve, ok := expr.(*VarExprQ); ok && ve.LevelsUp == 0 {
		return optionalRT[ve.RangeIdx]
	}
	return false
}

func resolveTargetColumnMetadata(c *Catalog, db *Database, q *Query, optionalRT map[int]bool, te *TargetEntryQ, widenOptionalJoin bool) (colType string, nullable bool, charset, collation string, ok bool) {
	if te == nil {
		return "", false, "", "", false
	}
	opt := widenOptionalJoin && exprUsesOptionalJoin(te.Expr, optionalRT)
	if te.ResOrigTable != "" && te.ResOrigCol != "" {
		origDB := db
		if te.ResOrigDB != "" && c != nil {
			origDB = c.GetDatabase(te.ResOrigDB)
		}
		if origDB != nil {
			if tbl := origDB.GetTable(te.ResOrigTable); tbl != nil {
				if srcCol := tbl.GetColumn(te.ResOrigCol); srcCol != nil {
					ci := viewColumnCollation(srcCol)
					return srcCol.ColumnType, srcCol.Nullable || opt, ci.charset, ci.collation, true
				}
			}
		}
	}
	ve, isVar := te.Expr.(*VarExprQ)
	if !isVar || q == nil || ve.LevelsUp != 0 || ve.RangeIdx < 0 || ve.RangeIdx >= len(q.RangeTable) {
		return "", false, "", "", false
	}
	rte := q.RangeTable[ve.RangeIdx]
	switch rte.Kind {
	case RTERelation:
		if ve.AttNum < 1 || ve.AttNum > len(rte.ColNames) {
			return "", false, "", "", false
		}
		origDB := db
		if rte.DBName != "" && c != nil {
			origDB = c.GetDatabase(rte.DBName)
		}
		if origDB == nil {
			return "", false, "", "", false
		}
		if tbl := origDB.GetTable(rte.TableName); tbl != nil {
			if srcCol := tbl.GetColumn(rte.ColNames[ve.AttNum-1]); srcCol != nil {
				ci := viewColumnCollation(srcCol)
				return srcCol.ColumnType, srcCol.Nullable || opt, ci.charset, ci.collation, true
			}
		}
		if sv := origDB.Views[strings.ToLower(rte.TableName)]; sv != nil {
			if sv.AnalyzedQuery != nil {
				if meta, ok := resolveQueryTargetMeta(c, origDB, sv.AnalyzedQuery, ve.AttNum-1, widenOptionalJoin); ok {
					return meta.colType, meta.nullable || opt, meta.charset, meta.collation, true
				}
			}
			if ve.AttNum-1 < len(sv.ColumnMetadata) {
				cm := sv.ColumnMetadata[ve.AttNum-1]
				return resolvedTypeToColumnType(cm.Type), cm.Nullable || opt, cm.Charset, cm.Collation, true
			}
		}
	case RTESubquery:
		if rte.Subquery != nil {
			if meta, ok := resolveQueryTargetMeta(c, db, rte.Subquery, ve.AttNum-1, widenOptionalJoin); ok {
				return meta.colType, meta.nullable || opt, meta.charset, meta.collation, true
			}
		}
	}
	return "", false, "", "", false
}

func resolveVarExprColumnType(c *Catalog, db *Database, q *Query, e *VarExprQ) string {
	if e == nil || q == nil || e.LevelsUp != 0 || e.RangeIdx < 0 || e.RangeIdx >= len(q.RangeTable) {
		return ""
	}
	rte := q.RangeTable[e.RangeIdx]
	if rte.Kind == RTERelation && c != nil && e.AttNum >= 1 && e.AttNum <= len(rte.ColNames) {
		origDB := db
		if rte.DBName != "" {
			origDB = c.GetDatabase(rte.DBName)
		}
		if origDB != nil {
			if tbl := origDB.GetTable(rte.TableName); tbl != nil {
				if col := tbl.GetColumn(rte.ColNames[e.AttNum-1]); col != nil {
					return col.ColumnType
				}
			}
			if sv := origDB.Views[strings.ToLower(rte.TableName)]; sv != nil && sv.AnalyzedQuery != nil {
				if meta, ok := resolveQueryTargetMeta(c, origDB, sv.AnalyzedQuery, e.AttNum-1, true); ok && meta.colType != "" {
					return meta.colType
				}
			}
		}
	}
	if rte.Kind == RTESubquery && rte.Subquery != nil {
		if meta, ok := resolveQueryTargetMeta(c, db, rte.Subquery, e.AttNum-1, true); ok && meta.colType != "" {
			return meta.colType
		}
	}
	return ""
}

func resolveAnalyzedExprColumnType(c *Catalog, db *Database, q *Query, expr AnalyzedExpr) string {
	if expr == nil {
		return ""
	}
	if vr, isVar := expr.(*VarExprQ); isVar {
		if colType := resolveVarExprColumnType(c, db, q, vr); colType != "" {
			return colType
		}
	}
	switch e := expr.(type) {
	case *CaseExprQ:
		var branches []AnalyzedExpr
		for _, w := range e.Args {
			if w != nil && w.Then != nil {
				branches = append(branches, w.Then)
			}
		}
		if e.Default != nil {
			branches = append(branches, e.Default)
		}
		return combineControlFlowColumnTypes(c, db, q, branches, false)
	case *CoalesceExprQ:
		return combineControlFlowColumnTypes(c, db, q, e.Args, e.IsIfNull)
	case *SubLinkExprQ:
		if ct := resolveSubLinkColumnType(c, db, e); ct != "" {
			return ct
		}
	case *OpExprQ:
		if strings.EqualFold(e.Op, "BINARY") {
			if ct := resolveBinaryUnaryColumnType(c, db, q, e.Right); ct != "" {
				return ct
			}
		}
		if ct := resolveIntervalOpColumnType(c, db, q, e, false); ct != "" {
			return ct
		}
	case *FuncCallExprQ:
		switch strings.ToLower(e.Name) {
		case "min", "max":
			if len(e.Args) == 1 {
				return resolveAnalyzedExprColumnType(c, db, q, e.Args[0])
			}
		case "if":
			if len(e.Args) == 3 {
				return combineControlFlowColumnTypes(c, db, q, []AnalyzedExpr{e.Args[1], e.Args[2]}, false)
			}
		case "nullif":
			if len(e.Args) >= 1 {
				return combineControlFlowColumnTypes(c, db, q, []AnalyzedExpr{e.Args[0]}, false)
			}
		default:
			if ct := resolveScalarFuncColumnType(c, db, q, e, false); ct != "" {
				return ct
			}
		}
	}
	if cn, ok := expr.(*ConstExprQ); ok {
		if cn.IsNull {
			return "binary(0)"
		}
		if cn.IsString {
			n := utf8.RuneCountInString(cn.Value)
			if n == 0 {
				return "char(0)"
			}
			return fmt.Sprintf("varchar(%d)", n)
		}
	}
	if rt := resolveExprResolvedType(c, db, q, expr); rt != nil {
		if rt.BaseType == BaseTypeTinyIntBool {
			return "int"
		}
		return resolvedTypeToColumnType(rt)
	}
	return ""
}

func combineControlFlowColumnTypes(c *Catalog, db *Database, q *Query, branches []AnalyzedExpr, preserveSmallInt bool) string {
	var combinedType string
	var combinedExpr AnalyzedExpr
	for _, b := range branches {
		bt := resolveAnalyzedExprColumnType(c, db, q, b)
		if bt == "" {
			continue
		}
		if combinedType == "" {
			combinedType = bt
			combinedExpr = b
			continue
		}
		lp := parseSetOpColType(combinedType)
		rp := parseSetOpColType(bt)
		if (lp.cat == "float" && rp.cat == "decimal") || (lp.cat == "decimal" && rp.cat == "float") {
			combinedType = "double"
			combinedExpr = b
			continue
		}
		combinedType = combineSetOpColumnTypes(combinedType, combinedExpr, bt, b, false)
	}
	return normalizeControlFlowColumnType(combinedType, preserveSmallInt)
}

func normalizeControlFlowColumnType(colType string, preserveSmallInt bool) string {
	p := parseSetOpColType(colType)
	if p.cat == "null" || p.raw == "" {
		return colType
	}
	if p.cat == "string" && p.isChar {
		return fmt.Sprintf("varchar(%d)", p.width)
	}
	if p.cat == "int" {
		suffix := ""
		if p.unsigned {
			suffix = " unsigned"
		}
		if p.rank >= 4 {
			return "bigint" + suffix
		}
		if !preserveSmallInt {
			return "int" + suffix
		}
	}
	return colType
}

func resolveExprResolvedType(c *Catalog, db *Database, q *Query, expr AnalyzedExpr) *ResolvedType {
	if expr == nil {
		return nil
	}
	switch e := expr.(type) {
	case *VarExprQ:
		if e.Type != nil {
			return e.Type
		}
		if q != nil && e.LevelsUp == 0 && e.RangeIdx >= 0 && e.RangeIdx < len(q.RangeTable) {
			rte := q.RangeTable[e.RangeIdx]
			if e.AttNum >= 1 && e.AttNum <= len(rte.ColTypes) && rte.ColTypes[e.AttNum-1] != nil {
				return rte.ColTypes[e.AttNum-1]
			}
			if rte.Kind == RTERelation && c != nil && e.AttNum >= 1 && e.AttNum <= len(rte.ColNames) {
				origDB := db
				if rte.DBName != "" {
					origDB = c.GetDatabase(rte.DBName)
				}
				if origDB != nil {
					if tbl := origDB.GetTable(rte.TableName); tbl != nil {
						if col := tbl.GetColumn(rte.ColNames[e.AttNum-1]); col != nil {
							return parseColumnResolvedType(col.ColumnType)
						}
					}
					if sv := origDB.Views[strings.ToLower(rte.TableName)]; sv != nil && sv.AnalyzedQuery != nil {
						if subTE := nthNonJunkTarget(sv.AnalyzedQuery, e.AttNum-1); subTE != nil {
							return resolveExprResolvedType(c, origDB, sv.AnalyzedQuery, subTE.Expr)
						}
					}
				}
			}
			if rte.Kind == RTESubquery && rte.Subquery != nil {
				if subTE := nthNonJunkTarget(rte.Subquery, e.AttNum-1); subTE != nil {
					return resolveExprResolvedType(c, db, rte.Subquery, subTE.Expr)
				}
			}
		}
		return nil
	case *CastExprQ:
		if e.TargetType != nil && e.TargetType.BaseType == BaseTypeDecimal && e.TargetType.Precision == 0 {
			return &ResolvedType{BaseType: BaseTypeDecimal, Precision: 10, Scale: 0}
		}
		return e.TargetType
	case *OpExprQ:
		opUpper := strings.ToUpper(e.Op)
		if e.Left == nil && (opUpper == "+" || opUpper == "-") {
			return resolveExprResolvedType(c, db, q, e.Right)
		}
		if ct := resolveIntervalOpColumnType(c, db, q, e, false); ct != "" {
			return parseColumnResolvedType(ct)
		}
		switch opUpper {
		case "=", "!=", "<>", "<", "<=", ">", ">=", "<=>",
			"LIKE", "NOT LIKE", "REGEXP", "RLIKE", "NOT REGEXP", "NOT RLIKE",
			"IS TRUE", "IS NOT TRUE", "IS FALSE", "IS NOT FALSE", "IS UNKNOWN", "IS NOT UNKNOWN",
			"XOR":
			return BoolType
		case "&", "|", "^", "<<", ">>", "~":
			return &ResolvedType{BaseType: BaseTypeBigInt, Unsigned: true}
		case "->":
			return &ResolvedType{BaseType: BaseTypeJSON}
		case "->>":
			return &ResolvedType{BaseType: BaseTypeLongText}
		case "+", "-", "*", "/", "%", "MOD":
			lt := resolveExprResolvedType(c, db, q, e.Left)
			rt := resolveExprResolvedType(c, db, q, e.Right)
			if (lt != nil && (lt.BaseType == BaseTypeDouble || lt.BaseType == BaseTypeFloat)) ||
				(rt != nil && (rt.BaseType == BaseTypeDouble || rt.BaseType == BaseTypeFloat)) {
				return &ResolvedType{BaseType: BaseTypeDouble}
			}
			if opUpper == "/" ||
				(lt != nil && lt.BaseType == BaseTypeDecimal) ||
				(rt != nil && rt.BaseType == BaseTypeDecimal) {
				if dt := combineDecimalArithmetic(opUpper, lt, e.Left, rt, e.Right); dt != nil {
					return dt
				}
				return &ResolvedType{BaseType: BaseTypeDecimal}
			}
			return &ResolvedType{BaseType: BaseTypeBigInt}
		case "DIV":
			return &ResolvedType{BaseType: BaseTypeBigInt}
		}
		return e.ResultType
	case *CaseExprQ, *CoalesceExprQ:
		if ct := resolveAnalyzedExprColumnType(c, db, q, expr); ct != "" {
			return parseColumnResolvedType(ct)
		}
		return expr.exprType()
	case *SubLinkExprQ:
		if ct := resolveSubLinkColumnType(c, db, e); ct != "" {
			return parseColumnResolvedType(ct)
		}
		return e.ResultType
	case *FuncCallExprQ:
		switch strings.ToLower(e.Name) {
		case "if", "nullif":
			if ct := resolveAnalyzedExprColumnType(c, db, q, expr); ct != "" {
				return parseColumnResolvedType(ct)
			}
		case "date_add", "date_sub", "adddate", "subdate",
			"timestampadd", "timediff", "addtime", "subtime":
			if ct := resolveScalarFuncColumnType(c, db, q, e, false); ct != "" {
				return parseColumnResolvedType(ct)
			}
		case "abs", "ceil", "ceiling", "floor", "round", "truncate", "mod":
			if rt := resolveNumericScalarFuncType(c, db, q, e); rt != nil {
				return rt
			}
		case "now", "current_timestamp", "sysdate", "localtime", "localtimestamp", "utc_timestamp",
			"curtime", "current_time", "utc_time", "time", "timestamp":
			if rt := resolveTemporalScalarFuncType(c, db, q, e); rt != nil {
				return rt
			}
		case "min", "max":
			if len(e.Args) == 1 {
				return resolveExprResolvedType(c, db, q, e.Args[0])
			}
		case "sum":
			if len(e.Args) == 1 {
				if rt := sumAvgResolvedType(c, db, q, e.Args[0], false); rt != nil {
					return rt
				}
			}
		case "avg":
			if len(e.Args) == 1 {
				if rt := sumAvgResolvedType(c, db, q, e.Args[0], true); rt != nil {
					return rt
				}
			}
		}
		return e.ResultType
	default:
		return expr.exprType()
	}
}

func nthNonJunkTarget(q *Query, idx int) *TargetEntryQ {
	if q == nil || idx < 0 {
		return nil
	}
	pos := 0
	for _, te := range q.TargetList {
		if te == nil || te.ResJunk {
			continue
		}
		if pos == idx {
			return te
		}
		pos++
	}
	return nil
}

func combineDecimalArithmetic(op string, lt *ResolvedType, le AnalyzedExpr, rt *ResolvedType, re AnalyzedExpr) *ResolvedType {
	lp, ls, lok := decimalOperandPrecisionScale(lt, le)
	rp, rs, rok := decimalOperandPrecisionScale(rt, re)
	if !lok || !rok {
		return nil
	}
	var prec, scale int
	switch op {
	case "+", "-":
		scale = max(ls, rs)
		prec = min(65, max(lp-ls, rp-rs)+1+scale)
	case "*":
		scale = min(30, ls+rs)
		prec = min(65, lp+rp)
	case "%", "MOD":
		scale = max(ls, rs)
		prec = min(65, max(lp-ls, rp-rs)+scale)
	case "/":
		scale = min(30, ls+4)
		prec = min(65, max(lp-ls, 0)+rs+scale)
	default:
		return nil
	}
	return &ResolvedType{
		BaseType:  BaseTypeDecimal,
		Precision: prec,
		Scale:     scale,
	}
}

func decimalOperandPrecisionScale(rt *ResolvedType, expr AnalyzedExpr) (precision, scale int, ok bool) {
	if c, okConst := expr.(*ConstExprQ); okConst && !c.IsNull && !c.IsString && rt != nil &&
		(rt.BaseType == BaseTypeBigInt || rt.BaseType == BaseTypeInt || rt.BaseType == BaseTypeTinyIntBool) {
		if c.Value == "TRUE" || c.Value == "FALSE" {
			return 1, 0, true
		}
		digits := len(strings.TrimLeft(c.Value, "+-"))
		if digits == 0 {
			digits = 1
		}
		return digits, 0, true
	}
	if rt == nil {
		return 0, 0, false
	}
	switch rt.BaseType {
	case BaseTypeDecimal:
		if rt.Precision > 0 {
			return rt.Precision, rt.Scale, true
		}
		return 10, 0, true
	case BaseTypeTinyInt:
		return 3, 0, true
	case BaseTypeTinyIntBool:
		return 1, 0, true
	case BaseTypeSmallInt:
		return 5, 0, true
	case BaseTypeMediumInt:
		if rt.Unsigned {
			return 8, 0, true
		}
		return 7, 0, true
	case BaseTypeInt:
		return 10, 0, true
	case BaseTypeBigInt:
		if rt.Unsigned {
			return 20, 0, true
		}
		return 19, 0, true
	default:
		return 0, 0, false
	}
}

func parseColumnResolvedType(colType string) *ResolvedType {
	s := strings.ToLower(strings.TrimSpace(colType))
	if s == "" {
		return nil
	}
	unsigned := strings.Contains(s, "unsigned")
	s = strings.ReplaceAll(s, " unsigned", "")
	s = strings.ReplaceAll(s, " zerofill", "")
	s = strings.TrimSpace(s)

	base := s
	var p, sc int
	if open := strings.IndexByte(s, '('); open >= 0 {
		base = strings.TrimSpace(s[:open])
		if closeIdx := strings.IndexByte(s[open+1:], ')'); closeIdx >= 0 {
			args := strings.Split(s[open+1:open+1+closeIdx], ",")
			if len(args) >= 1 {
				p, _ = strconv.Atoi(strings.TrimSpace(args[0]))
			}
			if len(args) >= 2 {
				sc, _ = strconv.Atoi(strings.TrimSpace(args[1]))
			}
		}
	}
	switch base {
	case "decimal", "numeric", "dec", "fixed":
		if p == 0 {
			p = 10
		}
		return &ResolvedType{BaseType: BaseTypeDecimal, Precision: p, Scale: sc, Unsigned: unsigned}
	case "tinyint":
		return &ResolvedType{BaseType: BaseTypeTinyInt, Unsigned: unsigned}
	case "smallint":
		return &ResolvedType{BaseType: BaseTypeSmallInt, Unsigned: unsigned}
	case "mediumint":
		return &ResolvedType{BaseType: BaseTypeMediumInt, Unsigned: unsigned}
	case "int", "integer":
		return &ResolvedType{BaseType: BaseTypeInt, Unsigned: unsigned}
	case "bigint":
		return &ResolvedType{BaseType: BaseTypeBigInt, Unsigned: unsigned}
	case "float":
		return &ResolvedType{BaseType: BaseTypeFloat, Unsigned: unsigned}
	case "double", "real":
		return &ResolvedType{BaseType: BaseTypeDouble, Unsigned: unsigned}
	case "date":
		return &ResolvedType{BaseType: BaseTypeDate}
	case "datetime":
		return &ResolvedType{BaseType: BaseTypeDateTime, Precision: p}
	case "timestamp":
		return &ResolvedType{BaseType: BaseTypeTimestamp, Precision: p}
	case "time":
		return &ResolvedType{BaseType: BaseTypeTime, Precision: p}
	case "year":
		return &ResolvedType{BaseType: BaseTypeYear}
	case "varchar", "var_string":
		return &ResolvedType{BaseType: BaseTypeVarchar, Length: p}
	case "char":
		return &ResolvedType{BaseType: BaseTypeChar, Length: p}
	default:
		return nil
	}
}

func constIntScalarArg(expr AnalyzedExpr) (int, bool) {
	switch e := expr.(type) {
	case *ConstExprQ:
		if e.IsNull || e.IsString {
			return 0, false
		}
		v, err := strconv.Atoi(strings.TrimSpace(e.Value))
		return v, err == nil
	case *OpExprQ:
		if e.Left == nil && (e.Op == "+" || e.Op == "-") {
			v, ok := constIntScalarArg(e.Right)
			if !ok {
				return 0, false
			}
			if e.Op == "-" {
				return -v, true
			}
			return v, true
		}
	}
	return 0, false
}

func resolveNumericScalarFuncType(c *Catalog, db *Database, q *Query, e *FuncCallExprQ) *ResolvedType {
	if e == nil || len(e.Args) == 0 {
		return nil
	}
	name := strings.ToLower(e.Name)
	arg0 := resolveExprResolvedType(c, db, q, e.Args[0])
	if arg0 == nil {
		return nil
	}
	isInt := func(bt MySQLBaseType) bool {
		switch bt {
		case BaseTypeTinyInt, BaseTypeTinyIntBool, BaseTypeSmallInt, BaseTypeMediumInt, BaseTypeInt, BaseTypeBigInt:
			return true
		}
		return false
	}
	switch name {
	case "abs":
		if arg0.BaseType == BaseTypeFloat || arg0.BaseType == BaseTypeDouble {
			return &ResolvedType{BaseType: BaseTypeDouble}
		}
		if arg0.BaseType == BaseTypeDecimal {
			return &ResolvedType{BaseType: BaseTypeDecimal, Precision: arg0.Precision, Scale: arg0.Scale}
		}
		if isInt(arg0.BaseType) {
			return &ResolvedType{BaseType: BaseTypeBigInt}
		}
	case "ceil", "ceiling", "floor":
		if arg0.BaseType == BaseTypeFloat || arg0.BaseType == BaseTypeDouble {
			return &ResolvedType{BaseType: BaseTypeDouble}
		}
		if arg0.BaseType == BaseTypeDecimal {
			if arg0.Precision <= 0 || arg0.Precision-arg0.Scale <= 18 {
				return &ResolvedType{BaseType: BaseTypeBigInt}
			}
			return &ResolvedType{BaseType: BaseTypeDecimal, Precision: min(65, arg0.Precision-arg0.Scale+1), Scale: 0}
		}
		if isInt(arg0.BaseType) {
			return &ResolvedType{BaseType: BaseTypeBigInt}
		}
	case "round", "truncate":
		if arg0.BaseType == BaseTypeFloat || arg0.BaseType == BaseTypeDouble {
			return &ResolvedType{BaseType: BaseTypeDouble}
		}
		if isInt(arg0.BaseType) {
			return &ResolvedType{BaseType: BaseTypeBigInt}
		}
		if arg0.BaseType == BaseTypeDecimal {
			d := 0
			if len(e.Args) >= 2 {
				if v, ok := constIntScalarArg(e.Args[1]); ok {
					d = v
				}
			}
			intDigits := max(0, arg0.Precision-arg0.Scale)
			if arg0.Precision == 0 {
				intDigits = 10
			}
			targetScale := max(0, min(arg0.Scale, d))
			if name == "truncate" {
				return &ResolvedType{BaseType: BaseTypeDecimal, Precision: min(65, max(1, intDigits+targetScale)), Scale: targetScale}
			}
			extra := 0
			if targetScale < arg0.Scale {
				extra = 1
			}
			return &ResolvedType{BaseType: BaseTypeDecimal, Precision: min(65, max(1, intDigits+extra+targetScale)), Scale: targetScale}
		}
	case "mod":
		if len(e.Args) == 2 {
			arg1 := resolveExprResolvedType(c, db, q, e.Args[1])
			if arg0.BaseType == BaseTypeFloat || arg0.BaseType == BaseTypeDouble ||
				(arg1 != nil && (arg1.BaseType == BaseTypeFloat || arg1.BaseType == BaseTypeDouble)) {
				return &ResolvedType{BaseType: BaseTypeDouble}
			}
			if arg0.BaseType == BaseTypeDecimal || (arg1 != nil && arg1.BaseType == BaseTypeDecimal) {
				if dt := combineDecimalArithmetic("MOD", arg0, e.Args[0], arg1, e.Args[1]); dt != nil {
					return dt
				}
				return &ResolvedType{BaseType: BaseTypeDecimal}
			}
			return &ResolvedType{BaseType: BaseTypeBigInt}
		}
	}
	return nil
}

func resolveSubLinkColumnType(c *Catalog, db *Database, s *SubLinkExprQ) string {
	if s == nil {
		return ""
	}
	switch s.Kind {
	case SubLinkScalar:
		if s.Subquery == nil {
			return ""
		}
		if meta, ok := resolveQueryTargetMeta(c, db, s.Subquery, 0, true); ok && meta.colType != "" {
			return normalizeScalarSubqueryColumnType(meta.colType)
		}
	case SubLinkExists, SubLinkIn, SubLinkAny, SubLinkAll:
		return "int"
	}
	return ""
}

func normalizeScalarSubqueryColumnType(colType string) string {
	p := parseSetOpColType(colType)
	if p.cat == "null" || p.raw == "" {
		return colType
	}
	if p.cat == "string" && p.isChar {
		return fmt.Sprintf("varchar(%d)", p.width)
	}
	if p.cat == "int" {
		suffix := ""
		if p.unsigned {
			suffix = " unsigned"
		}
		if p.rank >= 4 {
			return "bigint" + suffix
		}
		return "int" + suffix
	}
	return colType
}

func intervalUnitHasTimePart(unit string) bool {
	switch strings.ToUpper(unit) {
	case "MICROSECOND", "SECOND", "MINUTE", "HOUR",
		"SECOND_MICROSECOND", "MINUTE_MICROSECOND", "MINUTE_SECOND",
		"HOUR_MICROSECOND", "HOUR_SECOND", "HOUR_MINUTE",
		"DAY_MICROSECOND", "DAY_SECOND", "DAY_MINUTE", "DAY_HOUR":
		return true
	}
	return false
}

func intervalUnitHasMicrosecond(unit string) bool {
	return strings.Contains(strings.ToUpper(unit), "MICROSECOND")
}

func formatTemporalTypeWithFSP(base string, fsp int) string {
	if fsp > 0 {
		return fmt.Sprintf("%s(%d)", base, min(fsp, 6))
	}
	return base
}

func exprTemporalFSP(c *Catalog, db *Database, q *Query, expr AnalyzedExpr) int {
	if ct := resolveAnalyzedExprColumnType(c, db, q, expr); ct != "" {
		p := parseSetOpColType(ct)
		if p.cat == "time" || p.cat == "datetime" || p.cat == "timestamp" {
			return p.prec
		}
	}
	if cn, ok := expr.(*ConstExprQ); ok && cn.IsString {
		if _, frac, found := strings.Cut(cn.Value, "."); found {
			digits := 0
			for i := 0; i < len(frac) && frac[i] >= '0' && frac[i] <= '9'; i++ {
				digits++
			}
			return min(digits, 6)
		}
	}
	return 0
}

func resolveIntervalArithmeticColumnType(c *Catalog, db *Database, q *Query, base AnalyzedExpr, unit string, isTiDB bool) string {
	if base == nil {
		return ""
	}
	ct := resolveAnalyzedExprColumnType(c, db, q, base)
	if ct == "" {
		return ""
	}
	p := parseSetOpColType(ct)
	hasMicro := intervalUnitHasMicrosecond(unit)
	switch p.cat {
	case "date":
		if intervalUnitHasTimePart(unit) {
			if hasMicro {
				return "datetime(6)"
			}
			return "datetime"
		}
		return "date"
	case "datetime", "timestamp":
		fsp := p.prec
		if hasMicro {
			fsp = 6
		}
		return formatTemporalTypeWithFSP("datetime", fsp)
	case "time":
		fsp := p.prec
		if hasMicro {
			fsp = 6
		}
		return formatTemporalTypeWithFSP("time", fsp)
	case "string", "text", "null":
		if isTiDB {
			return "var_string(29)"
		}
		return "varchar(29)"
	}
	return ""
}

func resolveBinaryUnaryColumnType(c *Catalog, db *Database, q *Query, operand AnalyzedExpr) string {
	if operand == nil {
		return ""
	}
	if cn, ok := operand.(*ConstExprQ); ok {
		if cn.IsNull {
			return "binary(0)"
		}
		if cn.IsString {
			cs := ""
			if db != nil {
				cs = db.Charset
			}
			if cs == "" {
				cs = "utf8mb4"
			}
			return fmt.Sprintf("varbinary(%d)", utf8.RuneCountInString(cn.Value)*charsetMaxBytes(cs))
		}
	}
	ct := resolveAnalyzedExprColumnType(c, db, q, operand)
	if ct == "" {
		return ""
	}
	p := parseSetOpColType(ct)
	switch {
	case p.cat == "blob", p.cat == "binary":
		return p.raw
	case p.cat == "string", p.cat == "text":
		bytes := p.width * exprCharsetMaxBytes(c, db, q, operand)
		if bytes <= 65535 {
			return fmt.Sprintf("varbinary(%d)", bytes)
		}
		if bytes <= 16777215 {
			return "mediumblob"
		}
		return "longblob"
	case p.cat == "json":
		return "longblob"
	default:
		if w := setOpStringWidth(p, operand, false); w > 0 {
			return fmt.Sprintf("varbinary(%d)", w)
		}
	}
	return ""
}

func exprCharsetMaxBytes(c *Catalog, db *Database, q *Query, expr AnalyzedExpr) int {
	cs := ""
	if ve, ok := expr.(*VarExprQ); ok && q != nil && ve.LevelsUp == 0 && ve.RangeIdx >= 0 && ve.RangeIdx < len(q.RangeTable) {
		rte := q.RangeTable[ve.RangeIdx]
		if rte.Kind == RTERelation && c != nil && ve.AttNum >= 1 && ve.AttNum <= len(rte.ColNames) {
			origDB := db
			if rte.DBName != "" {
				origDB = c.GetDatabase(rte.DBName)
			}
			if origDB != nil {
				if tbl := origDB.GetTable(rte.TableName); tbl != nil {
					if col := tbl.GetColumn(rte.ColNames[ve.AttNum-1]); col != nil {
						cs = col.Charset
					}
				}
			}
		}
	}
	if cs == "" && db != nil {
		cs = db.Charset
	}
	if cs == "" {
		cs = "utf8mb4"
	}
	return charsetMaxBytes(cs)
}

func resolveIntervalOpColumnType(c *Catalog, db *Database, q *Query, e *OpExprQ, isTiDB bool) string {
	if e == nil {
		return ""
	}
	switch e.Op {
	case "+":
		if iv, ok := e.Right.(*IntervalExprQ); ok {
			return resolveIntervalArithmeticColumnType(c, db, q, e.Left, iv.Unit, isTiDB)
		}
		if iv, ok := e.Left.(*IntervalExprQ); ok {
			return resolveIntervalArithmeticColumnType(c, db, q, e.Right, iv.Unit, isTiDB)
		}
	case "-":
		if iv, ok := e.Right.(*IntervalExprQ); ok {
			return resolveIntervalArithmeticColumnType(c, db, q, e.Left, iv.Unit, isTiDB)
		}
	}
	return ""
}

func resolveTemporalScalarFuncType(c *Catalog, db *Database, q *Query, e *FuncCallExprQ) *ResolvedType {
	if e == nil {
		return nil
	}
	name := strings.ToLower(e.Name)
	argFSP := func() int {
		if len(e.Args) >= 1 {
			if fsp, ok := constIntScalarArg(e.Args[0]); ok && fsp > 0 {
				return min(fsp, 6)
			}
		}
		return 0
	}
	exprFSP := func(arg AnalyzedExpr) int {
		ct := resolveAnalyzedExprColumnType(c, db, q, arg)
		if ct == "" {
			return 0
		}
		p := parseSetOpColType(ct)
		if p.cat == "time" || p.cat == "datetime" || p.cat == "timestamp" {
			return p.prec
		}
		return 0
	}
	switch name {
	case "now", "current_timestamp", "sysdate", "localtime", "localtimestamp", "utc_timestamp":
		return &ResolvedType{BaseType: BaseTypeDateTime, Precision: argFSP()}
	case "curtime", "current_time", "utc_time":
		return &ResolvedType{BaseType: BaseTypeTime, Precision: argFSP()}
	case "time":
		fsp := 0
		if len(e.Args) >= 1 {
			fsp = exprFSP(e.Args[0])
		}
		return &ResolvedType{BaseType: BaseTypeTime, Precision: fsp}
	case "timestamp":
		fsp := 0
		for _, arg := range e.Args {
			fsp = max(fsp, exprFSP(arg))
		}
		return &ResolvedType{BaseType: BaseTypeDateTime, Precision: fsp}
	}
	return nil
}

func resolveScalarFuncColumnType(c *Catalog, db *Database, q *Query, e *FuncCallExprQ, isTiDB bool) string {
	if e == nil {
		return ""
	}
	strType := func(w int) string {
		if isTiDB {
			return fmt.Sprintf("var_string(%d)", w)
		}
		return fmt.Sprintf("varchar(%d)", w)
	}
	argWidth := func(arg AnalyzedExpr) (parsedSetOpType, int, bool) {
		ct := resolveAnalyzedExprColumnType(c, db, q, arg)
		if ct == "" {
			return parsedSetOpType{}, 0, false
		}
		p := parseSetOpColType(ct)
		return p, p.width, true
	}
	name := strings.ToLower(e.Name)
	switch name {
	case "date_add", "date_sub", "adddate", "subdate":
		if len(e.Args) == 2 {
			unit := "DAY"
			if iv, ok := e.Args[1].(*IntervalExprQ); ok {
				unit = iv.Unit
			}
			return resolveIntervalArithmeticColumnType(c, db, q, e.Args[0], unit, isTiDB)
		}
	case "timestampadd":
		if len(e.Args) == 3 {
			unit := ""
			if kw, ok := e.Args[0].(*ConstExprQ); ok {
				unit = kw.Value
			}
			if isTiDB {
				if intervalUnitHasMicrosecond(unit) {
					return "var_string(26)"
				}
				return "var_string(19)"
			}
			return resolveIntervalArithmeticColumnType(c, db, q, e.Args[2], unit, false)
		}
	case "timediff":
		if len(e.Args) == 2 {
			fsp := max(exprTemporalFSP(c, db, q, e.Args[0]), exprTemporalFSP(c, db, q, e.Args[1]))
			return formatTemporalTypeWithFSP("time", fsp)
		}
	case "addtime", "subtime":
		if len(e.Args) == 2 {
			if ct := resolveAnalyzedExprColumnType(c, db, q, e.Args[0]); ct != "" {
				p := parseSetOpColType(ct)
				fsp := max(p.prec, exprTemporalFSP(c, db, q, e.Args[1]))
				switch p.cat {
				case "time":
					return formatTemporalTypeWithFSP("time", fsp)
				case "datetime", "timestamp":
					return formatTemporalTypeWithFSP("datetime", fsp)
				}
			}
		}
	case "json_type":
		if isTiDB {
			return "var_string(51)"
		}
		return "varchar(17)"
	case "uuid":
		return strType(36)
	case "version", "database", "schema", "charset", "collation":
		return strType(64)
	case "user", "current_user", "session_user", "system_user":
		if isTiDB {
			return "var_string(64)"
		}
		return "varchar(288)"
	case "lower", "lcase", "upper", "ucase", "trim", "btrim", "ltrim", "rtrim", "reverse":
		if len(e.Args) >= 1 {
			if p, w, ok := argWidth(e.Args[0]); ok {
				if p.cat == "text" {
					if isTiDB {
						return fmt.Sprintf("var_string(%d)", w)
					}
					if p.rank >= 3 {
						return p.raw
					}
					return "mediumtext"
				}
				return strType(w)
			}
		}
	case "replace":
		if len(e.Args) == 3 {
			if p, w, ok := argWidth(e.Args[0]); ok {
				if isTiDB {
					if fromC, ok1 := e.Args[1].(*ConstExprQ); ok1 && fromC.IsString {
						if toC, ok2 := e.Args[2].(*ConstExprQ); ok2 && toC.IsString {
							fl := max(1, utf8.RuneCountInString(fromC.Value))
							tl := utf8.RuneCountInString(toC.Value)
							if tl > fl {
								w = w * tl / fl
							}
						}
					}
				}
				if p.cat == "text" {
					if p.rank >= 3 {
						return p.raw
					}
					return "mediumtext"
				}
				return strType(w)
			}
		}
	case "left", "right":
		if len(e.Args) >= 2 {
			if p, w, ok := argWidth(e.Args[0]); ok {
				if !isTiDB {
					if n, okN := constIntScalarArg(e.Args[1]); okN && n >= 0 {
						w = min(w, n)
					}
				}
				if p.cat == "text" && !isTiDB {
					if w <= 65535 {
						return fmt.Sprintf("varchar(%d)", w)
					}
					return p.raw
				}
				return strType(w)
			}
		}
	case "substring", "substr":
		if len(e.Args) >= 2 {
			if p, w, ok := argWidth(e.Args[0]); ok {
				if !isTiDB {
					if len(e.Args) >= 3 {
						if n, okN := constIntScalarArg(e.Args[2]); okN && n >= 0 {
							w = min(w, n)
						}
					} else if pos, okP := constIntScalarArg(e.Args[1]); okP && pos > 0 {
						w = max(0, w-(pos-1))
					}
				}
				if p.cat == "text" && !isTiDB {
					if w <= 65535 {
						return fmt.Sprintf("varchar(%d)", w)
					}
					return p.raw
				}
				return strType(w)
			}
		}
	case "lpad", "rpad":
		if len(e.Args) >= 2 {
			if n, okN := constIntScalarArg(e.Args[1]); okN && n >= 0 {
				if isTiDB {
					return fmt.Sprintf("var_string(%d)", n*4)
				}
				return fmt.Sprintf("varchar(%d)", n)
			}
		}
	case "repeat":
		if len(e.Args) == 2 {
			if isTiDB {
				return "longtext"
			}
			if _, w, ok := argWidth(e.Args[0]); ok {
				if n, okN := constIntScalarArg(e.Args[1]); okN && n >= 0 {
					return fmt.Sprintf("varchar(%d)", w*n)
				}
			}
		}
	case "space":
		if len(e.Args) == 1 {
			if isTiDB {
				return "longtext"
			}
			if n, okN := constIntScalarArg(e.Args[0]); okN && n >= 0 {
				return fmt.Sprintf("varchar(%d)", n)
			}
		}
	case "format":
		if isTiDB {
			return "longtext"
		}
		return "varchar(48)"
	case "concat", "concat_ws":
		if len(e.Args) == 0 {
			return ""
		}
		total := 0
		hasText := false
		maxRank := 0
		if name == "concat_ws" {
			if len(e.Args) < 2 {
				return ""
			}
			sp, sw, sok := argWidth(e.Args[0])
			if !sok {
				return ""
			}
			if sp.cat == "text" {
				hasText = true
				maxRank = max(maxRank, sp.rank)
			}
			nonNull := 0
			for _, arg := range e.Args[1:] {
				ap, aw, aok := argWidth(arg)
				if !aok {
					return ""
				}
				if ap.cat == "null" {
					continue
				}
				if ap.cat == "text" {
					hasText = true
					maxRank = max(maxRank, ap.rank)
				}
				total += aw
				nonNull++
			}
			if nonNull > 1 {
				total += sw * (nonNull - 1)
			}
		} else {
			for _, arg := range e.Args {
				ap, aw, aok := argWidth(arg)
				if !aok {
					return ""
				}
				if ap.cat == "null" {
					return "binary(0)"
				}
				if ap.cat == "text" {
					hasText = true
					maxRank = max(maxRank, ap.rank)
				}
				total += aw
			}
		}
		if hasText || total > 65535 {
			if maxRank <= 2 {
				return "mediumtext"
			}
			if maxRank == 3 {
				return "mediumtext"
			}
			return "longtext"
		}
		return strType(total)
	}
	return ""
}

func resolvedTypeToColumnType(rt *ResolvedType) string {
	if rt == nil {
		return ""
	}
	unsigned := ""
	if rt.Unsigned {
		unsigned = " unsigned"
	}
	switch rt.BaseType {
	case BaseTypeTinyInt:
		return "tinyint" + unsigned
	case BaseTypeTinyIntBool:
		return "tinyint(1)"
	case BaseTypeSmallInt:
		return "smallint" + unsigned
	case BaseTypeMediumInt:
		return "mediumint" + unsigned
	case BaseTypeInt:
		return "int" + unsigned
	case BaseTypeBigInt:
		return "bigint" + unsigned
	case BaseTypeDecimal:
		if rt.Precision > 0 {
			return fmt.Sprintf("decimal(%d,%d)%s", rt.Precision, rt.Scale, unsigned)
		}
		return "decimal" + unsigned
	case BaseTypeFloat:
		return "float" + unsigned
	case BaseTypeDouble:
		return "double" + unsigned
	case BaseTypeBit:
		if rt.Length > 0 {
			return fmt.Sprintf("bit(%d)", rt.Length)
		}
		return "bit(1)"
	case BaseTypeDate:
		return "date"
	case BaseTypeDateTime:
		if rt.Precision > 0 {
			return fmt.Sprintf("datetime(%d)", rt.Precision)
		}
		return "datetime"
	case BaseTypeTimestamp:
		if rt.Precision > 0 {
			return fmt.Sprintf("timestamp(%d)", rt.Precision)
		}
		return "timestamp"
	case BaseTypeTime:
		if rt.Precision > 0 {
			return fmt.Sprintf("time(%d)", rt.Precision)
		}
		return "time"
	case BaseTypeYear:
		return "year"
	case BaseTypeChar:
		if rt.Length > 0 {
			return fmt.Sprintf("char(%d)", rt.Length)
		}
		return "char"
	case BaseTypeVarchar:
		if rt.Length > 0 {
			return fmt.Sprintf("varchar(%d)", rt.Length)
		}
		return "varchar"
	case BaseTypeBinary:
		if rt.Length > 0 {
			return fmt.Sprintf("binary(%d)", rt.Length)
		}
		return "binary"
	case BaseTypeVarBinary:
		if rt.Length > 0 {
			return fmt.Sprintf("varbinary(%d)", rt.Length)
		}
		return "varbinary"
	case BaseTypeTinyText:
		return "tinytext"
	case BaseTypeText:
		return "text"
	case BaseTypeMediumText:
		return "mediumtext"
	case BaseTypeLongText:
		return "longtext"
	case BaseTypeTinyBlob:
		return "tinyblob"
	case BaseTypeBlob:
		return "blob"
	case BaseTypeMediumBlob:
		return "mediumblob"
	case BaseTypeLongBlob:
		return "longblob"
	case BaseTypeJSON:
		return "json"
	default:
		return ""
	}
}

func tableToProto(dbName string, t *Table) *metadata.TableMetadata {
	table := &metadata.TableMetadata{
		Name:       t.Name,
		Engine:     t.Engine,
		Charset:    t.Charset,
		Collation:  t.Collation,
		Comment:    formatSyncedComment(t.Comment),
		Tablespace: t.Tablespace,
		SkipDump:   t.SkipDump,
	}
	if opts := tableCreateOptions(t); opts != "" {
		table.CreateOptions = opts
	}

	pos := int32(0)
	for _, col := range t.Columns {
		if col.Hidden == ColumnHiddenSystem {
			// System-hidden columns (e.g. backing columns of functional
			// indexes) are not reported by information_schema.
			continue
		}
		pos++
		colMeta := &metadata.ColumnMetadata{
			Name:                           col.Name,
			Position:                       pos,
			Type:                           col.ColumnType,
			Nullable:                       col.Nullable,
			CharacterSet:                   col.Charset,
			Collation:                      col.Collation,
			Comment:                        formatSyncedComment(col.Comment),
			NotSecondary:                   col.NotSecondary,
			IsGeneratedInvisiblePrimaryKey: col.GeneratedInvisiblePrimaryKey,
		}
		if col.Charset != "" {
			colMeta.CharsetInherited = boolPtr(!col.CharsetExplicit)
		}
		if col.Collation != "" {
			colMeta.CollationInherited = boolPtr(!col.CollationExplicit)
		}
		switch {
		case col.AutoIncrement:
			colMeta.Default = autoIncrementSymbol
		case col.Default != nil:
			colMeta.Default = normalizeDefault(*col.Default, col.DefaultKind)
		case col.Nullable && !col.DefaultDropped:
			colMeta.Default = "NULL"
		default:
		}
		colMeta.DefaultDropped = col.DefaultDropped
		if col.OnUpdate != "" {
			colMeta.OnUpdate = col.OnUpdate
		}
		if col.Generated != nil {
			genType := metadata.GenerationMetadata_TYPE_VIRTUAL
			if col.Generated.Stored {
				genType = metadata.GenerationMetadata_TYPE_STORED
			}
			colMeta.Generation = &metadata.GenerationMetadata{
				Type:       genType,
				Expression: col.Generated.Expr,
			}
		}
		// Spatial SRID (MySQL 8.0): presence-carrying — an explicit SRID (including 0)
		// is distinct from no SRID.
		if col.HasSRID {
			srid := uint32(col.SRID)
			colMeta.Srid = &srid
		}
		colMeta.IsInvisible = col.Invisible
		table.Columns = append(table.Columns, colMeta)
	}

	sortedIndexes := slices.Clone(t.Indexes)
	slices.SortStableFunc(sortedIndexes, func(a, b *Index) int {
		return cmp.Compare(a.Name, b.Name)
	})
	for _, idx := range sortedIndexes {
		indexType := strings.ToUpper(idx.IndexType)
		if indexType == "" {
			indexType = "BTREE"
			if strings.EqualFold(t.Engine, "MEMORY") || strings.EqualFold(t.Engine, "HEAP") {
				indexType = "HASH"
			}
		}
		idxMeta := &metadata.IndexMetadata{
			Name:         idx.Name,
			Type:         indexType,
			Unique:       idx.Unique,
			Primary:      idx.Primary,
			Visible:      idx.Visible,
			Comment:      idx.Comment,
			KeyBlockSize: int64(idx.KeyBlockSize),
		}
		for _, col := range idx.Columns {
			if col.Expr != "" {
				idxMeta.Expressions = append(idxMeta.Expressions, formatIndexExpr(col.Expr))
			} else {
				idxMeta.Expressions = append(idxMeta.Expressions, col.Name)
			}
			keyLength := int64(col.Length)
			if keyLength == 0 {
				keyLength = -1
			}
			if keyLength == -1 && indexType == mysqlIndexSpatial {
				// MySQL reports a key length of 32 for spatial index columns.
				keyLength = 32
			}
			idxMeta.KeyLength = append(idxMeta.KeyLength, keyLength)
			idxMeta.Descending = append(idxMeta.Descending, col.Descending)
		}
		table.Indexes = append(table.Indexes, idxMeta)
	}

	for _, con := range t.Constraints {
		switch con.Type {
		case ConForeignKey:
			// Omitted actions default to NO ACTION in MySQL 8.0 REFERENTIAL_CONSTRAINTS.
			onUpdate := con.OnUpdate
			if onUpdate == "" {
				onUpdate = "NO ACTION"
			}
			onDelete := con.OnDelete
			if onDelete == "" {
				onDelete = "NO ACTION"
			}
			matchType := con.MatchType
			if matchType == "" {
				matchType = "NONE"
			}
			fk := &metadata.ForeignKeyMetadata{
				Name:              con.Name,
				Columns:           slices.Clone(con.Columns),
				ReferencedTable:   con.RefTable,
				ReferencedColumns: slices.Clone(con.RefColumns),
				OnUpdate:          onUpdate,
				OnDelete:          onDelete,
				MatchType:         matchType,
			}
			if con.RefDatabase != "" && !eqFoldStr(con.RefDatabase, dbName) {
				fk.ReferencedSchema = con.RefDatabase
			}
			table.ForeignKeys = append(table.ForeignKeys, fk)
		case ConCheck:
			table.CheckConstraints = append(table.CheckConstraints, &metadata.CheckConstraintMetadata{
				Name:        con.Name,
				Expression:  con.CheckExpr,
				NotEnforced: con.NotEnforced,
			})
		default:
		}
	}

	slices.SortStableFunc(table.ForeignKeys, func(a, b *metadata.ForeignKeyMetadata) int {
		return cmp.Compare(a.Name, b.Name)
	})
	slices.SortStableFunc(table.CheckConstraints, func(a, b *metadata.CheckConstraintMetadata) int {
		return cmp.Compare(a.Name, b.Name)
	})

	if t.Partitioning != nil {
		table.Partitions = partitionsToProto(t.Partitioning)
	}

	return table
}

func tableCreateOptions(t *Table) string {
	var opts []string
	if t.KeyBlockSize > 0 {
		opts = append(opts, fmt.Sprintf("key_block_size=%d", t.KeyBlockSize))
	}
	if t.RowFormat != "" {
		opts = append(opts, "row_format="+t.RowFormat)
	}
	if t.Compression != "" {
		opts = append(opts, fmt.Sprintf("compression='%s'", escapeComment(unquoteTableOption(t.Compression))))
	}
	if t.Encryption != "" {
		opts = append(opts, fmt.Sprintf("encryption='%s'", escapeComment(unquoteTableOption(t.Encryption))))
	}
	if t.StatsPersistent != "" && !eqFoldStr(t.StatsPersistent, "DEFAULT") {
		opts = append(opts, "stats_persistent="+t.StatsPersistent)
	}
	if t.StatsAutoRecalc != "" && !eqFoldStr(t.StatsAutoRecalc, "DEFAULT") {
		opts = append(opts, "stats_auto_recalc="+t.StatsAutoRecalc)
	}
	if t.StatsSamplePages != "" && !eqFoldStr(t.StatsSamplePages, "DEFAULT") {
		opts = append(opts, "stats_sample_pages="+t.StatsSamplePages)
	}
	if t.MinRows != "" && t.MinRows != "0" {
		opts = append(opts, "min_rows="+t.MinRows)
	}
	if t.MaxRows != "" && t.MaxRows != "0" {
		opts = append(opts, "max_rows="+t.MaxRows)
	}
	if t.AvgRowLength != "" && t.AvgRowLength != "0" {
		opts = append(opts, "avg_row_length="+t.AvgRowLength)
	}
	if t.PackKeys != "" && !eqFoldStr(t.PackKeys, "DEFAULT") {
		opts = append(opts, "pack_keys="+t.PackKeys)
	}
	if t.Checksum != "" && t.Checksum != "0" {
		opts = append(opts, "checksum="+t.Checksum)
	}
	if t.DelayKeyWrite != "" && t.DelayKeyWrite != "0" {
		opts = append(opts, "delay_key_write="+t.DelayKeyWrite)
	}
	if t.SecondaryEngine != "" {
		opts = append(opts, fmt.Sprintf("secondary_engine='%s'", escapeComment(unquoteTableOption(t.SecondaryEngine))))
	}
	if t.AutoExtendSize > 0 {
		opts = append(opts, fmt.Sprintf("autoextend_size=%d", t.AutoExtendSize))
	}
	if t.Partitioning != nil {
		opts = append(opts, "partitioned")
	}
	return strings.Join(opts, " ")
}

func normalizeRoutineDefinition(def string, isProcedure bool) string {
	if def == "" {
		return ""
	}
	keyword := " FUNCTION "
	if isProcedure {
		keyword = " PROCEDURE "
	}
	if idx := strings.Index(def, keyword); idx >= 0 {
		def = "CREATE" + def[idx:]
	}
	if !isProcedure {
		def = stripReturnsCharset(def)
	}
	return def
}

func stripReturnsCharset(def string) string {
	returnsIdx := strings.Index(def, " RETURNS ")
	if returnsIdx == -1 {
		return def
	}
	lineEnd := strings.Index(def[returnsIdx:], "\n")
	if lineEnd == -1 {
		return def
	}
	charsetIdx := strings.Index(def[returnsIdx:returnsIdx+lineEnd], " CHARSET ")
	if charsetIdx == -1 {
		return def
	}
	return def[:returnsIdx+charsetIdx] + def[returnsIdx+lineEnd:]
}

func sumAvgResolvedType(c *Catalog, db *Database, q *Query, arg AnalyzedExpr, isAvg bool) *ResolvedType {
	rt := resolveExprResolvedType(c, db, q, arg)
	if rt == nil {
		return nil
	}
	switch rt.BaseType {
	case BaseTypeFloat, BaseTypeDouble:
		return &ResolvedType{BaseType: BaseTypeDouble}
	case BaseTypeTinyInt, BaseTypeTinyIntBool, BaseTypeSmallInt, BaseTypeMediumInt, BaseTypeInt, BaseTypeBigInt, BaseTypeDecimal:
		p, s, ok := decimalOperandPrecisionScale(rt, arg)
		if !ok {
			return nil
		}
		if isAvg {
			return &ResolvedType{BaseType: BaseTypeDecimal, Precision: min(65, p+4), Scale: min(30, s+4)}
		}
		return &ResolvedType{BaseType: BaseTypeDecimal, Precision: min(65, p+22), Scale: s}
	default:
		return nil
	}
}

func typedTemporalLiteralValue(value string) (string, bool) {
	trimmed := strings.TrimSpace(value)
	upper := strings.ToUpper(trimmed)
	for _, kw := range []string{"TIMESTAMP", "DATE", "TIME"} {
		if strings.HasPrefix(upper, kw) {
			rest := strings.TrimSpace(trimmed[len(kw):])
			if len(rest) >= 2 && rest[0] == '\'' && rest[len(rest)-1] == '\'' {
				return rest, true
			}
		}
	}
	return "", false
}

// normalizeDefault converts a catalog default value to the driver representation:
// static defaults are single-quoted for mysqldump compatibility (except bit/hex
// literals) and expression defaults are wrapped in parentheses.
func normalizeDefault(value string, kind ColumnDefaultKind) string {
	switch kind {
	case ColumnDefaultConstant:
		if strings.EqualFold(value, "NULL") {
			return "NULL"
		}
		if scalar, ok := typedTemporalLiteralValue(value); ok {
			return scalar
		}
		if strings.HasPrefix(value, "b'") || strings.HasPrefix(value, "B'") ||
			strings.HasPrefix(value, "0x") || strings.HasPrefix(value, "0X") ||
			isBitOrHexLiteral(value) {
			return value
		}
		if len(value) >= 2 && strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
			return value
		}
		return fmt.Sprintf("'%s'", value)
	case ColumnDefaultExpression:
		return fmt.Sprintf("(%s)", value)
	default:
		return value
	}
}

func partitionsToProto(p *PartitionInfo) []*metadata.TablePartitionMetadata {
	expr := p.Expr
	if expr == "" && len(p.Columns) > 0 {
		expr = strings.Join(p.Columns, ",")
	}
	partType := partitionTypeToProto(p.Type, p.Linear)
	subExpr := p.SubExpr
	if subExpr == "" && len(p.SubColumns) > 0 {
		subExpr = strings.Join(p.SubColumns, ",")
	}
	subType := partitionTypeToProto(p.SubType, p.SubLinear)

	var partUseDefault string
	if p.NumParts > 0 && (p.Type == "HASH" || p.Type == "KEY") {
		partUseDefault = strconv.Itoa(p.NumParts)
	}
	var subUseDefault string
	if p.SubType != "" && p.UseDefaultSubpartitions && !p.UseDefaultNumSubpartitions && p.NumSubParts > 0 {
		subUseDefault = strconv.Itoa(p.NumSubParts)
	}

	var result []*metadata.TablePartitionMetadata
	for _, pd := range p.Partitions {
		part := &metadata.TablePartitionMetadata{
			Name:       pd.Name,
			Type:       partType,
			Expression: expr,
			UseDefault: partUseDefault,
			Value:      pd.ValueExpr,
			Comment:    pd.Comment,
			Algorithm:  int32(p.Algorithm),
		}
		if pd.Tablespace != "" {
			part.Tablespace = stringPtr(pd.Tablespace)
		}
		if p.SubType != "" {
			switch {
			case p.UseDefaultSubpartitions && p.UseDefaultNumSubpartitions:
				// Preserve implicit subpartition generation rather than exporting auto-populated entries.
				part.Subpartitions = []*metadata.TablePartitionMetadata{{
					Type:       subType,
					Expression: subExpr,
					UseDefault: "default",
					Algorithm:  int32(p.SubAlgo),
				}}
			case len(pd.SubPartitions) > 0:
				for _, spd := range pd.SubPartitions {
					spMeta := &metadata.TablePartitionMetadata{
						Name:       spd.Name,
						Type:       subType,
						Expression: subExpr,
						UseDefault: subUseDefault,
						Comment:    spd.Comment,
						Algorithm:  int32(p.SubAlgo),
					}
					if spd.Tablespace != "" {
						spMeta.Tablespace = stringPtr(spd.Tablespace)
					}
					part.Subpartitions = append(part.Subpartitions, spMeta)
				}
			case subUseDefault != "":
				part.Subpartitions = []*metadata.TablePartitionMetadata{{
					Type:       subType,
					Expression: subExpr,
					UseDefault: subUseDefault,
					Algorithm:  int32(p.SubAlgo),
				}}
			}
		}
		result = append(result, part)
	}
	return result
}

func partitionTypeToProto(t string, linear bool) metadata.TablePartitionMetadata_Type {
	switch strings.ToUpper(t) {
	case "RANGE":
		return metadata.TablePartitionMetadata_RANGE
	case "RANGE COLUMNS":
		return metadata.TablePartitionMetadata_RANGE_COLUMNS
	case "LIST":
		return metadata.TablePartitionMetadata_LIST
	case "LIST COLUMNS":
		return metadata.TablePartitionMetadata_LIST_COLUMNS
	case "HASH":
		if linear {
			return metadata.TablePartitionMetadata_LINEAR_HASH
		}
		return metadata.TablePartitionMetadata_HASH
	case "LINEAR HASH":
		return metadata.TablePartitionMetadata_LINEAR_HASH
	case "KEY":
		if linear {
			return metadata.TablePartitionMetadata_LINEAR_KEY
		}
		return metadata.TablePartitionMetadata_KEY
	case "LINEAR KEY":
		return metadata.TablePartitionMetadata_LINEAR_KEY
	default:
		return metadata.TablePartitionMetadata_TYPE_UNSPECIFIED
	}
}

func triggerToProto(tr *Trigger) *metadata.TriggerMetadata {
	tm := &metadata.TriggerMetadata{
		Name:     tr.Name,
		Event:    tr.Event,
		Timing:   tr.Timing,
		Body:     tr.Body,
		SkipDump: tr.SkipDump,
	}
	if tr.HasSessionContext {
		tm.SqlMode = tr.SQLMode
		tm.CharacterSetClient = tr.CharacterSetClient
		tm.CollationConnection = tr.CollationConnection
	}
	return tm
}

func formatSyncedComment(s string) string {
	if !strings.ContainsAny(s, "\\'\x00\x1a") {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) + 4)
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\\':
			b.WriteString(`\\`)
		case '\'':
			b.WriteString(`\'`)
		case 0:
			b.WriteString(`\0`)
		case 26:
			b.WriteString(`\Z`)
		default:
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

func formatIndexExpr(expr string) string {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return expr
	}
	if expr[0] == '(' && expr[len(expr)-1] == ')' {
		depth := 0
		wrapped := true
		inSingle, inDouble, inBacktick := false, false, false
		for i := 0; i < len(expr); i++ {
			ch := expr[i]
			if inSingle {
				if ch == '\\' {
					i++
					continue
				}
				if ch == '\'' {
					if i+1 < len(expr) && expr[i+1] == '\'' {
						i++
						continue
					}
					inSingle = false
				}
				continue
			}
			if inDouble {
				if ch == '\\' {
					i++
					continue
				}
				if ch == '"' {
					inDouble = false
				}
				continue
			}
			if inBacktick {
				if ch == '`' {
					if i+1 < len(expr) && expr[i+1] == '`' {
						i++
						continue
					}
					inBacktick = false
				}
				continue
			}
			switch ch {
			case '\'':
				inSingle = true
			case '"':
				inDouble = true
			case '`':
				inBacktick = true
			case '(':
				depth++
			case ')':
				depth--
				if depth == 0 && i < len(expr)-1 {
					wrapped = false
				}
			}
		}
		if wrapped && depth == 0 {
			return expr
		}
	}
	return "(" + expr + ")"
}

func sortedMapKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
