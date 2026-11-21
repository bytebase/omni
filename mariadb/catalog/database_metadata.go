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

// DatabaseMetadata converts the named database in the catalog to a
// DatabaseSchemaMetadata snapshot matching the MariaDB driver sync and
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

	// MariaDB uses a single empty-name schema.
	schemaMeta := &metadata.SchemaMetadata{
		Name: "",
	}

	for _, sName := range sortedMapKeys(db.Sequences) {
		seq := db.Sequences[sName]
		if seq == nil || seq.Temporary {
			continue
		}
		schemaMeta.Sequences = append(schemaMeta.Sequences, &metadata.SequenceMetadata{
			Name:      seq.Name,
			DataType:  seq.DataType,
			Start:     seq.Start,
			MinValue:  seq.MinValue,
			MaxValue:  seq.MaxValue,
			Increment: seq.Increment,
			CacheSize: seq.Cache,
			Cycle:     seq.Cycle,
			Comment:   formatSyncedComment(seq.Comment),
		})
	}

	// tableProtoByLowerName lets the trigger pass below attach each trigger to its
	// owning table's metadata. The catalog keys tables case-insensitively, and a
	// trigger records its table via Trigger.Table, so the lookup is lower-cased.
	tableProtoByLowerName := make(map[string]*metadata.TableMetadata, len(db.Tables))
	for _, tName := range sortedMapKeys(db.Tables) {
		t := db.Tables[tName]
		if t == nil || t.Temporary {
			continue
		}
		tableProto := tableToProto(db.Name, t)
		schemaMeta.Tables = append(schemaMeta.Tables, tableProto)
		tableProtoByLowerName[strings.ToLower(t.Name)] = tableProto
	}

	// Triggers are database-level in MariaDB (db.Triggers), each bound to a table
	// via Trigger.Table; the proto nests them under TableMetadata.Triggers in
	// execution order (respecting FOLLOWS/PRECEDES), matching live MariaDB sync.
	triggers := make([]*Trigger, 0, len(db.Triggers))
	for _, tr := range db.Triggers {
		if tr != nil {
			triggers = append(triggers, tr)
		}
	}
	type triggerGroupKey struct{ table, event, timing string }
	groupOrders := make(map[triggerGroupKey]map[string]int)
	triggerOrderPos := func(tr *Trigger) int {
		gk := triggerGroupKey{
			table:  strings.ToUpper(tr.Table),
			event:  strings.ToUpper(tr.Event),
			timing: strings.ToUpper(tr.Timing),
		}
		order, ok := groupOrders[gk]
		if !ok {
			order = effectiveTriggerOrder(db, tr)
			groupOrders[gk] = order
		}
		return order[toLower(tr.Name)]
	}
	slices.SortFunc(triggers, func(a, b *Trigger) int {
		return cmp.Or(
			cmp.Compare(strings.ToUpper(a.Table), strings.ToUpper(b.Table)),
			cmp.Compare(strings.ToUpper(a.Event), strings.ToUpper(b.Event)),
			cmp.Compare(strings.ToUpper(a.Timing), strings.ToUpper(b.Timing)),
			cmp.Compare(triggerOrderPos(a), triggerOrderPos(b)),
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
		if v == nil {
			continue
		}
		def := v.Definition
		if opt := strings.ToUpper(strings.TrimSpace(v.CheckOption)); opt == "CASCADED" || opt == "LOCAL" {
			// Match live MariaDB sync (and writeView/writeViewSDL, which emit Definition verbatim).
			suffix := " WITH " + opt + " CHECK OPTION"
			if !strings.HasSuffix(strings.ToUpper(strings.TrimSpace(def)), suffix) {
				def += suffix
			}
		}
		schemaMeta.Views = append(schemaMeta.Views, &metadata.ViewMetadata{
			Name:        v.Name,
			Definition:  def,
			Columns:     c.viewColumnsToProto(db, v),
			Algorithm:   v.Algorithm,
			SqlSecurity: v.SqlSecurity,
			CheckOption: v.CheckOption,
		})
	}

	for _, fName := range sortedMapKeys(db.Functions) {
		f := db.Functions[fName]
		if f == nil {
			continue
		}
		fnMeta := &metadata.FunctionMetadata{
			Name:              f.Name,
			Definition:        normalizeRoutineDefinition(c.ShowCreateFunction(dbName, f.Name), false),
			DatabaseCollation: db.Collation,
			Comment:           f.Characteristics["COMMENT"],
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
		if p == nil {
			continue
		}
		procMeta := &metadata.ProcedureMetadata{
			Name:              p.Name,
			Definition:        normalizeRoutineDefinition(c.ShowCreateProcedure(dbName, p.Name), true),
			DatabaseCollation: db.Collation,
			Comment:           p.Characteristics["COMMENT"],
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
		if e == nil {
			continue
		}
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
		} else if te != nil && te.Expr != nil {
			if nullable, known := analyzedExprNullable(c, db, leftQ, optionalRT, te.Expr, true); known {
				col.Nullable = nullable
			}
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
			if isStringColumnType(meta.colType) && db != nil && db.Charset != "" {
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

	combinedType := combineSetOpColumnTypes(leftMeta.colType, leftMeta.expr, rightMeta.colType, rightMeta.expr)
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
	nullable, _ := analyzedExprNullable(c, db, q, opt, expr, widenOptionalJoin)
	return nullable
}

func analyzedExprNullable(c *Catalog, db *Database, q *Query, opt map[int]bool, expr AnalyzedExpr, widenOptionalJoin bool) (bool, bool) {
	if expr == nil {
		return false, false
	}
	switch e := expr.(type) {
	case *ConstExprQ:
		return e.IsNull, true
	case *VarExprQ:
		if widenOptionalJoin && opt[e.RangeIdx] {
			return true, true
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
							return col.Nullable, true
						}
					}
					if sv := origDB.Views[strings.ToLower(rte.TableName)]; sv != nil && sv.AnalyzedQuery != nil {
						if meta, ok := resolveQueryTargetMeta(c, origDB, sv.AnalyzedQuery, e.AttNum-1, widenOptionalJoin); ok {
							return meta.nullable, true
						}
					}
				}
			}
			if (rte.Kind == RTESubquery || rte.Kind == RTECTE) && rte.Subquery != nil && rte.Subquery != q {
				if meta, ok := resolveQueryTargetMeta(c, db, rte.Subquery, e.AttNum-1, widenOptionalJoin); ok {
					return meta.nullable, true
				}
			}
		}
		return false, false
	case *CastExprQ:
		return analyzedExprNullable(c, db, q, opt, e.Arg, widenOptionalJoin)
	case *OpExprQ:
		if e.Op == "/" || e.Op == "%" || strings.EqualFold(e.Op, "DIV") || strings.EqualFold(e.Op, "MOD") {
			return true, true
		}
		switch strings.ToUpper(e.Op) {
		case "<=>", "IS TRUE", "IS NOT TRUE", "IS FALSE", "IS NOT FALSE", "IS UNKNOWN", "IS NOT UNKNOWN":
			return false, true
		case "BINARY":
			return true, true
		}
		return isAnalyzedExprNullable(c, db, q, opt, e.Left, widenOptionalJoin) || isAnalyzedExprNullable(c, db, q, opt, e.Right, widenOptionalJoin), true
	case *NullTestExprQ:
		return false, true
	case *BetweenExprQ:
		return isAnalyzedExprNullable(c, db, q, opt, e.Arg, widenOptionalJoin) ||
			isAnalyzedExprNullable(c, db, q, opt, e.Lower, widenOptionalJoin) ||
			isAnalyzedExprNullable(c, db, q, opt, e.Upper, widenOptionalJoin), true
	case *InListExprQ:
		if isAnalyzedExprNullable(c, db, q, opt, e.Arg, widenOptionalJoin) {
			return true, true
		}
		for _, item := range e.List {
			if isAnalyzedExprNullable(c, db, q, opt, item, widenOptionalJoin) {
				return true, true
			}
		}
		return false, true
	case *BoolExprQ:
		for _, arg := range e.Args {
			if isAnalyzedExprNullable(c, db, q, opt, arg, widenOptionalJoin) {
				return true, true
			}
		}
		return false, true
	case *SubLinkExprQ:
		switch e.Kind {
		case SubLinkExists:
			return false, true
		case SubLinkScalar, SubLinkIn, SubLinkAny, SubLinkAll:
			return true, true
		}
		return false, false
	case *FuncCallExprQ:
		if strings.EqualFold(e.Name, "count") {
			return false, true
		}
		switch strings.ToLower(e.Name) {
		case "nullif", "mod",
			"date_add", "date_sub", "adddate", "subdate",
			"timestampadd", "timestampdiff", "datediff", "timediff", "addtime", "subtime":
			return true, true
		}
		if strings.EqualFold(e.Name, "if") && len(e.Args) == 3 {
			return isAnalyzedExprNullable(c, db, q, opt, e.Args[1], widenOptionalJoin) ||
				isAnalyzedExprNullable(c, db, q, opt, e.Args[2], widenOptionalJoin), true
		}
		if isViewStringMetadataNullableFunc(e.Name) || isViewNullableAggregateFunc(e.Name) {
			return true, true
		}
		for _, arg := range e.Args {
			if isAnalyzedExprNullable(c, db, q, opt, arg, widenOptionalJoin) {
				return true, true
			}
		}
		return false, true
	case *CoalesceExprQ:
		if len(e.Args) == 0 {
			return true, true
		}
		for _, arg := range e.Args {
			if !isAnalyzedExprNullable(c, db, q, opt, arg, widenOptionalJoin) {
				return false, true
			}
		}
		return true, true
	case *CaseExprQ:
		if e.Default == nil || isAnalyzedExprNullable(c, db, q, opt, e.Default, widenOptionalJoin) {
			return true, true
		}
		for _, w := range e.Args {
			if w != nil && isAnalyzedExprNullable(c, db, q, opt, w.Then, widenOptionalJoin) {
				return true, true
			}
		}
		return false, true
	default:
		return false, false
	}
}

func isViewNullableAggregateFunc(name string) bool {
	switch strings.ToLower(name) {
	case "sum", "avg", "min", "max", "group_concat",
		"bit_and", "bit_or", "bit_xor",
		"std", "stddev", "stddev_pop", "stddev_samp",
		"var_pop", "var_samp", "variance",
		"json_arrayagg", "json_objectagg", "any_value":
		return true
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
	cat      string
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

func isConstString(cn *ConstExprQ) bool {
	if cn == nil || cn.IsNull || cn.Keyword {
		return false
	}
	if cn.IsString {
		return true
	}
	if cn.Type != nil {
		return isResolvedStringType(cn.Type)
	}
	if cn.Value == "TRUE" || cn.Value == "FALSE" || isNumberLiteral(cn.Value) {
		return false
	}
	return true
}

func setOpStringWidth(pt parsedSetOpType, expr AnalyzedExpr) int {
	if pt.cat == "int" {
		if cn, ok := expr.(*ConstExprQ); ok && !cn.IsNull && !isConstString(cn) {
			digits := len(strings.TrimLeft(cn.Value, "+-"))
			if digits == 0 {
				digits = 1
			}
			return digits + 1
		}
	}
	return pt.width
}

func combineSetOpColumnTypes(lt string, le AnalyzedExpr, rt string, re AnalyzedExpr) string {
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
		w := max(setOpStringWidth(lp, le), setOpStringWidth(rp, re))
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
		w := max(setOpStringWidth(lp, le), setOpStringWidth(rp, re))
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
	return fmt.Sprintf("varchar(%d)", max(setOpStringWidth(lp, le), setOpStringWidth(rp, re)))
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
	case RTESubquery, RTECTE:
		if rte.Subquery != nil && rte.Subquery != q {
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
	if (rte.Kind == RTESubquery || rte.Kind == RTECTE) && rte.Subquery != nil && rte.Subquery != q {
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
		return combineControlFlowColumnTypes(c, db, q, branches)
	case *CoalesceExprQ:
		return combineControlFlowColumnTypes(c, db, q, e.Args)
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
	case *FuncCallExprQ:
		switch strings.ToLower(e.Name) {
		case "min", "max":
			if len(e.Args) == 1 {
				return resolveAnalyzedExprColumnType(c, db, q, e.Args[0])
			}
		case "if":
			if len(e.Args) == 3 {
				return combineControlFlowColumnTypes(c, db, q, []AnalyzedExpr{e.Args[1], e.Args[2]})
			}
		case "nullif":
			if len(e.Args) >= 1 {
				return combineControlFlowColumnTypes(c, db, q, []AnalyzedExpr{e.Args[0]})
			}
		default:
			if ct := resolveScalarFuncColumnType(c, db, q, e); ct != "" {
				return ct
			}
		}
	}
	if cn, ok := expr.(*ConstExprQ); ok {
		if cn.IsNull {
			return "binary(0)"
		}
		if isConstString(cn) {
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

func combineControlFlowColumnTypes(c *Catalog, db *Database, q *Query, branches []AnalyzedExpr) string {
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
		combinedType = combineSetOpColumnTypes(combinedType, combinedExpr, bt, b)
	}
	return normalizeControlFlowColumnType(combinedType)
}

func normalizeControlFlowColumnType(colType string) string {
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

func resolveExprResolvedType(c *Catalog, db *Database, q *Query, expr AnalyzedExpr) *ResolvedType {
	if expr == nil {
		return nil
	}
	switch e := expr.(type) {
	case *ConstExprQ:
		if e.Type != nil {
			return e.Type
		}
		if e.IsNull {
			return nil
		}
		if e.IsString {
			return &ResolvedType{BaseType: BaseTypeVarchar, Length: utf8.RuneCountInString(e.Value)}
		}
		if e.Value == "TRUE" || e.Value == "FALSE" {
			return BoolType
		}
		if isNumberLiteral(e.Value) {
			if strings.ContainsAny(e.Value, "eE") {
				return &ResolvedType{BaseType: BaseTypeDouble}
			}
			if strings.Contains(e.Value, ".") {
				whole, frac, _ := strings.Cut(strings.TrimLeft(e.Value, "+-"), ".")
				return &ResolvedType{BaseType: BaseTypeDecimal, Precision: len(whole) + len(frac), Scale: len(frac)}
			}
			return &ResolvedType{BaseType: BaseTypeInt}
		}
		return nil
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
			if (rte.Kind == RTESubquery || rte.Kind == RTECTE) && rte.Subquery != nil && rte.Subquery != q {
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
		if e.Left == nil && opUpper == "+" {
			return resolveExprResolvedType(c, db, q, e.Right)
		}
		if e.Left == nil && opUpper == "-" {
			rt := resolveExprResolvedType(c, db, q, e.Right)
			if rt == nil {
				return nil
			}
			if isArithmeticStringOperand(rt) {
				return &ResolvedType{BaseType: BaseTypeDouble}
			}
			switch rt.BaseType {
			case BaseTypeTinyInt, BaseTypeSmallInt, BaseTypeMediumInt:
				return &ResolvedType{BaseType: BaseTypeInt}
			case BaseTypeInt, BaseTypeBigInt:
				return &ResolvedType{BaseType: BaseTypeBigInt}
			}
			cp := *rt
			cp.Unsigned = false
			return &cp
		}
		switch opUpper {
		case "=", "!=", "<>", "<", "<=", ">", ">=", "<=>",
			"LIKE", "NOT LIKE", "REGEXP", "RLIKE", "NOT REGEXP", "NOT RLIKE",
			"IS TRUE", "IS NOT TRUE", "IS FALSE", "IS NOT FALSE", "IS UNKNOWN", "IS NOT UNKNOWN",
			"XOR":
			return BoolType
		case "&", "|", "^", "<<", ">>", "~":
			if bt := bitwiseBinaryStringType(c, db, q, e); bt != nil {
				return bt
			}
			return &ResolvedType{BaseType: BaseTypeBigInt, Unsigned: true}
		case "+", "-", "*", "/", "%", "MOD":
			lt := resolveExprResolvedType(c, db, q, e.Left)
			rt := resolveExprResolvedType(c, db, q, e.Right)
			if (lt != nil && (lt.BaseType == BaseTypeDouble || lt.BaseType == BaseTypeFloat)) ||
				(rt != nil && (rt.BaseType == BaseTypeDouble || rt.BaseType == BaseTypeFloat)) ||
				isArithmeticStringOperand(lt) || isArithmeticStringOperand(rt) {
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
			return &ResolvedType{BaseType: BaseTypeBigInt, Unsigned: (lt != nil && lt.Unsigned) || (rt != nil && rt.Unsigned)}
		case "DIV":
			lt := resolveExprResolvedType(c, db, q, e.Left)
			rt := resolveExprResolvedType(c, db, q, e.Right)
			if isArithmeticStringOperand(lt) || isArithmeticStringOperand(rt) {
				return &ResolvedType{BaseType: BaseTypeBigInt}
			}
			return &ResolvedType{BaseType: BaseTypeBigInt, Unsigned: (lt != nil && lt.Unsigned) || (rt != nil && rt.Unsigned)}
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
		case "timediff", "addtime", "subtime":
			if ct := resolveScalarFuncColumnType(c, db, q, e); ct != "" {
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
	if c, okConst := expr.(*ConstExprQ); okConst && !c.IsNull && !isConstString(c) && rt != nil &&
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
		return &ResolvedType{BaseType: BaseTypeDateTime, FSP: p}
	case "timestamp":
		return &ResolvedType{BaseType: BaseTypeTimestamp, FSP: p}
	case "time":
		return &ResolvedType{BaseType: BaseTypeTime, FSP: p}
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
		if e.IsNull || isConstString(e) {
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
	if cn, ok := expr.(*ConstExprQ); ok && isConstString(cn) {
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

func resolveBinaryUnaryColumnType(c *Catalog, db *Database, q *Query, operand AnalyzedExpr) string {
	if operand == nil {
		return ""
	}
	if cn, ok := operand.(*ConstExprQ); ok {
		if cn.IsNull {
			return "binary(0)"
		}
		if isConstString(cn) {
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
		if w := setOpStringWidth(p, operand); w > 0 {
			return fmt.Sprintf("varbinary(%d)", w)
		}
	}
	return ""
}

func charsetMaxBytes(charset string) int {
	switch normalizeCharsetName(charset) {
	case "utf8mb4", "utf32", "utf16", "utf16le":
		return 4
	case "utf8mb3", "utf8":
		return 3
	case "ucs2", "gbk", "big5", "sjis", "cp932", "euckr", "gb2312":
		return 2
	default:
		return 1
	}
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
		return &ResolvedType{BaseType: BaseTypeDateTime, FSP: argFSP()}
	case "curtime", "current_time", "utc_time":
		return &ResolvedType{BaseType: BaseTypeTime, FSP: argFSP()}
	case "time":
		fsp := 0
		if len(e.Args) >= 1 {
			fsp = exprFSP(e.Args[0])
		}
		return &ResolvedType{BaseType: BaseTypeTime, FSP: fsp}
	case "timestamp":
		fsp := 0
		for _, arg := range e.Args {
			fsp = max(fsp, exprFSP(arg))
		}
		return &ResolvedType{BaseType: BaseTypeDateTime, FSP: fsp}
	}
	return nil
}

func resolveScalarFuncColumnType(c *Catalog, db *Database, q *Query, e *FuncCallExprQ) string {
	if e == nil {
		return ""
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
		return "varchar(17)"
	case "uuid":
		return "varchar(36)"
	case "version", "database", "schema", "charset", "collation":
		return "varchar(64)"
	case "user", "current_user", "session_user", "system_user":
		return "varchar(288)"
	case "lower", "lcase", "upper", "ucase", "trim", "btrim", "ltrim", "rtrim", "reverse":
		if len(e.Args) >= 1 {
			if p, w, ok := argWidth(e.Args[0]); ok {
				if p.cat == "text" {
					if p.rank >= 3 {
						return p.raw
					}
					return "mediumtext"
				}
				return fmt.Sprintf("varchar(%d)", w)
			}
		}
	case "replace":
		if len(e.Args) == 3 {
			if p, w, ok := argWidth(e.Args[0]); ok {
				if p.cat == "text" {
					if p.rank >= 3 {
						return p.raw
					}
					return "mediumtext"
				}
				return fmt.Sprintf("varchar(%d)", w)
			}
		}
	case "left", "right":
		if len(e.Args) >= 2 {
			if p, w, ok := argWidth(e.Args[0]); ok {
				if n, okN := constIntScalarArg(e.Args[1]); okN && n >= 0 {
					w = min(w, n)
				}
				if p.cat == "text" {
					if w <= 65535 {
						return fmt.Sprintf("varchar(%d)", w)
					}
					return p.raw
				}
				return fmt.Sprintf("varchar(%d)", w)
			}
		}
	case "substring", "substr":
		if len(e.Args) >= 2 {
			if p, w, ok := argWidth(e.Args[0]); ok {
				if len(e.Args) >= 3 {
					if n, okN := constIntScalarArg(e.Args[2]); okN && n >= 0 {
						w = min(w, n)
					}
				} else if pos, okP := constIntScalarArg(e.Args[1]); okP && pos > 0 {
					w = max(0, w-(pos-1))
				}
				if p.cat == "text" {
					if w <= 65535 {
						return fmt.Sprintf("varchar(%d)", w)
					}
					return p.raw
				}
				return fmt.Sprintf("varchar(%d)", w)
			}
		}
	case "lpad", "rpad":
		if len(e.Args) >= 2 {
			if n, okN := constIntScalarArg(e.Args[1]); okN && n >= 0 {
				return fmt.Sprintf("varchar(%d)", n)
			}
		}
	case "repeat":
		if len(e.Args) == 2 {
			if _, w, ok := argWidth(e.Args[0]); ok {
				if n, okN := constIntScalarArg(e.Args[1]); okN && n >= 0 {
					return fmt.Sprintf("varchar(%d)", w*n)
				}
			}
		}
	case "space":
		if len(e.Args) == 1 {
			if n, okN := constIntScalarArg(e.Args[0]); okN && n >= 0 {
				return fmt.Sprintf("varchar(%d)", n)
			}
		}
	case "format":
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
			if maxRank <= 3 {
				return "mediumtext"
			}
			return "longtext"
		}
		return fmt.Sprintf("varchar(%d)", total)
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
		if rt.FSP > 0 {
			return fmt.Sprintf("datetime(%d)", rt.FSP)
		}
		return "datetime"
	case BaseTypeTimestamp:
		if rt.FSP > 0 {
			return fmt.Sprintf("timestamp(%d)", rt.FSP)
		}
		return "timestamp"
	case BaseTypeTime:
		if rt.FSP > 0 {
			return fmt.Sprintf("time(%d)", rt.FSP)
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
	engine := t.Engine
	if engine == "" {
		engine = "InnoDB"
	}
	table := &metadata.TableMetadata{
		Name:       t.Name,
		Engine:     engine,
		Charset:    t.Charset,
		Collation:  t.Collation,
		Comment:    formatSyncedComment(t.Comment),
		Tablespace: t.Tablespace,
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
			Name:           col.Name,
			Position:       pos,
			Type:           col.ColumnType,
			Nullable:       col.Nullable,
			CharacterSet:   col.Charset,
			Collation:      col.Collation,
			Comment:        formatSyncedComment(col.Comment),
			DefaultDropped: col.DefaultDropped,
			IsInvisible:    col.Invisible,
		}
		switch {
		case col.AutoIncrement:
			colMeta.Default = autoIncrementDefault
		case col.Generated != nil:
			// Generated columns carry no default.
		case col.Default != nil:
			colMeta.Default = formatSyncedDefault(col)
		case col.Nullable && !col.DefaultDropped:
			colMeta.Default = "NULL"
		default:
		}
		if col.OnUpdate != "" {
			colMeta.OnUpdate = col.OnUpdate
		}
		if col.Generated != nil {
			if col.Generated.RowBound != "" {
				// System-versioning period columns (GENERATED ALWAYS AS ROW START/END)
				// are always NOT NULL in MariaDB (information_schema reports IS_NULLABLE = 'NO')
				// and are reported as stored generated columns with the bound keyword.
				colMeta.Nullable = false
				colMeta.Generation = &metadata.GenerationMetadata{
					Type:       metadata.GenerationMetadata_TYPE_STORED,
					Expression: col.Generated.RowBound,
				}
			} else {
				genType := metadata.GenerationMetadata_TYPE_VIRTUAL
				if col.Generated.Stored {
					genType = metadata.GenerationMetadata_TYPE_STORED
				}
				colMeta.Generation = &metadata.GenerationMetadata{
					Type:       genType,
					Expression: col.Generated.Expr,
				}
			}
		}
		if col.SRID != 0 {
			srid := uint32(col.SRID)
			colMeta.Srid = &srid
		}
		table.Columns = append(table.Columns, colMeta)
	}

	sortedIndexes := slices.Clone(t.Indexes)
	slices.SortStableFunc(sortedIndexes, func(a, b *Index) int {
		return cmp.Compare(a.Name, b.Name)
	})
	for _, idx := range sortedIndexes {
		var indexType string
		switch {
		case idx.Fulltext:
			indexType = mysqlIndexFullText
		case idx.Spatial:
			indexType = mysqlIndexSpatial
		case idx.Vector:
			indexType = "VECTOR"
		default:
			indexType = strings.ToUpper(idx.IndexType)
			if indexType == "" {
				indexType = defaultIndexType(engine)
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
				idxMeta.Expressions = append(idxMeta.Expressions, escapeFunctionalIndexQuotes(formatIndexExpr(col.Expr)))
			} else {
				idxMeta.Expressions = append(idxMeta.Expressions, col.Name)
			}
			keyLength := int64(col.Length)
			if keyLength == 0 {
				keyLength = -1
			}
			if keyLength == -1 && indexType == mysqlIndexSpatial {
				// Spatial index key parts report sub_part 32 in information_schema.STATISTICS.
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
			// Omitted referential actions default to RESTRICT in MariaDB.
			onUpdate := con.OnUpdate
			if onUpdate == "" {
				onUpdate = "RESTRICT"
			}
			onDelete := con.OnDelete
			if onDelete == "" {
				onDelete = "RESTRICT"
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

// formatSyncedDefault renders a column's default in the form MariaDB's sync
// records: CURRENT_TIMESTAMP forms and NULL stay bare, while all other defaults
// are QUOTE(COLUMN_DEFAULT) — single-quoted with backslashes and quotes escaped
// around MariaDB's COLUMN_DEFAULT text (e.g. "'\\'abc\\”", "'0'", "'1 + 2'",
// "'curdate()'"). This round-trips losslessly through LoadMetadata and
// DatabaseDefinition.
func formatSyncedDefault(col *Column) string {
	val := *col.Default
	if strings.EqualFold(val, "NULL") {
		return "NULL"
	}
	if col.DefaultKind == ColumnDefaultCurrentTimestamp {
		return val
	}
	if ts, ok := currentTimestampSQL(val); ok {
		return ts
	}
	rendered := formatDefault(val, col)
	if strings.EqualFold(rendered, "NULL") {
		return "NULL"
	}
	if col.DefaultKind == ColumnDefaultExpression {
		// information_schema.COLUMNS.COLUMN_DEFAULT omits the outer parentheses
		// that SHOW CREATE TABLE adds around operator expressions; stripping one
		// pair here matches sync and keeps renderColumnDefault from double-wrapping.
		rendered = stripOuterParens(rendered)
	}
	return "'" + formatSyncedComment(rendered) + "'"
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
	if t.PageChecksum != "" {
		opts = append(opts, "page_checksum="+t.PageChecksum)
	}
	if t.DelayKeyWrite != "" && t.DelayKeyWrite != "0" {
		opts = append(opts, "delay_key_write="+t.DelayKeyWrite)
	}
	if t.Transactional != "" {
		opts = append(opts, "transactional="+t.Transactional)
	}
	for _, o := range t.EngineOptions {
		value := o.Value
		if o.Quoted {
			value = "'" + escapeComment(o.Value) + "'"
		}
		opts = append(opts, fmt.Sprintf("`%s`=%s", o.Name, value))
	}
	if t.Partitioning != nil {
		opts = append(opts, "partitioned")
	}
	if t.SystemVersioned {
		opts = append(opts, systemVersioningOption)
	}
	return strings.Join(opts, " ")
}

func normalizeRoutineDefinition(def string, isProcedure bool) string {
	if def == "" {
		return ""
	}
	def = stripLeadingDefiner(def)
	if !isProcedure {
		def = stripReturnsCharset(def)
	}
	return def
}

// stripReturnsCharset drops the CHARSET clause SHOW CREATE FUNCTION writes
// after a string return type, as live sync does. Every search skips quoted
// text: RETURNS ENUM('a CHARSET b','c') holds " CHARSET " in a member.
func stripReturnsCharset(def string) string {
	returnsIdx := indexOutsideQuotes(def, " RETURNS ", 0)
	if returnsIdx == -1 {
		return def
	}
	lineEnd := indexOutsideQuotes(def, "\n", returnsIdx)
	if lineEnd == -1 {
		return def
	}
	charsetIdx := indexOutsideQuotes(def[:lineEnd], " CHARSET ", returnsIdx)
	if charsetIdx == -1 {
		return def
	}
	return def[:charsetIdx] + def[lineEnd:]
}

// indexOutsideQuotes is the offset of the first sub in s at or after start
// that is not inside a '...', "..." or `...` quoted run, or -1.
func indexOutsideQuotes(s, sub string, start int) int {
	var quote byte
	for i := start; i < len(s); i++ {
		c := s[i]
		if quote != 0 {
			switch {
			case c == '\\' && quote != '`':
				i++
			case c == quote:
				if i+1 < len(s) && s[i+1] == quote {
					i++
				} else {
					quote = 0
				}
			}
			continue
		}
		if c == '\'' || c == '"' || c == '`' {
			quote = c
			continue
		}
		if strings.HasPrefix(s[i:], sub) {
			return i
		}
	}
	return -1
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

func partitionsToProto(p *PartitionInfo) []*metadata.TablePartitionMetadata {
	partType := partitionTypeToProto(p.Type, p.Linear)
	if partType == metadata.TablePartitionMetadata_TYPE_UNSPECIFIED {
		// SYSTEM_TIME partitioning has no TablePartitionMetadata_Type enum value;
		// live MariaDB sync omits it as well.
		return nil
	}
	expr := p.Expr
	if expr == "" && len(p.Columns) > 0 {
		expr = joinPartitionColumns(p.Columns)
	}
	subExpr := p.SubExpr
	if subExpr == "" && len(p.SubColumns) > 0 {
		subExpr = joinPartitionColumns(p.SubColumns)
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
		if p.SubType != "" {
			switch {
			case p.UseDefaultSubpartitions && p.UseDefaultNumSubpartitions:
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
		Name:   tr.Name,
		Event:  tr.Event,
		Timing: tr.Timing,
		Body:   tr.Body,
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

func joinPartitionColumns(cols []string) string {
	quoted := make([]string, len(cols))
	for i, c := range cols {
		quoted[i] = "`" + strings.ReplaceAll(c, "`", "``") + "`"
	}
	return strings.Join(quoted, ",")
}

func escapeFunctionalIndexQuotes(expr string) string {
	if !strings.ContainsAny(expr, `'\`) {
		return expr
	}
	var b strings.Builder
	b.Grow(len(expr) + 8)
	for i := 0; i < len(expr); i++ {
		if expr[i] == '\'' || expr[i] == '\\' {
			b.WriteByte('\\')
		}
		b.WriteByte(expr[i])
	}
	return b.String()
}

func isArithmeticStringOperand(rt *ResolvedType) bool {
	if rt == nil {
		return false
	}
	switch rt.BaseType {
	case BaseTypeChar, BaseTypeVarchar, BaseTypeTinyText, BaseTypeText, BaseTypeMediumText, BaseTypeLongText,
		BaseTypeBinary, BaseTypeVarBinary, BaseTypeTinyBlob, BaseTypeBlob, BaseTypeMediumBlob, BaseTypeLongBlob,
		BaseTypeEnum, BaseTypeSet:
		return true
	}
	return false
}

func bitwiseBinaryStringType(c *Catalog, db *Database, q *Query, e *OpExprQ) *ResolvedType {
	binLen := func(x AnalyzedExpr) (int, bool) {
		if x == nil {
			return 0, false
		}
		ct := strings.ToLower(strings.TrimSpace(resolveAnalyzedExprColumnType(c, db, q, x)))
		for _, prefix := range []string{"varbinary(", "binary("} {
			if strings.HasPrefix(ct, prefix) && strings.HasSuffix(ct, ")") {
				if n, err := strconv.Atoi(ct[len(prefix) : len(ct)-1]); err == nil && n > 0 {
					return n, true
				}
			}
		}
		return 0, false
	}
	switch e.Op {
	case "~":
		if n, ok := binLen(e.Right); ok {
			return &ResolvedType{BaseType: BaseTypeVarBinary, Length: n}
		}
	case "<<", ">>":
		if n, ok := binLen(e.Left); ok {
			return &ResolvedType{BaseType: BaseTypeVarBinary, Length: n}
		}
	default:
		ln, lok := binLen(e.Left)
		rn, rok := binLen(e.Right)
		if lok && rok {
			return &ResolvedType{BaseType: BaseTypeVarBinary, Length: max(ln, rn)}
		}
	}
	return nil
}
