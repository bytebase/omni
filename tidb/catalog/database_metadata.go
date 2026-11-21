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
// DatabaseSchemaMetadata snapshot matching the TiDB sync and LoadMetadata
// round-trip contract.
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
	dbMeta.PlacementPolicy = db.PlacementPolicy
	dbMeta.TiflashReplica = int32(db.TiFlashReplica)
	dbMeta.TiflashLocationLabels = slices.Clone(db.TiFlashLocationLabels)

	schemaMeta := &metadata.SchemaMetadata{
		Name: "",
	}

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
		trMeta := &metadata.TriggerMetadata{
			Name:     tr.Name,
			Event:    tr.Event,
			Timing:   tr.Timing,
			Body:     tr.Body,
			SkipDump: tr.SkipDump,
		}
		if tr.HasSessionContext {
			trMeta.SqlMode = tr.SQLMode
			trMeta.CharacterSetClient = tr.CharacterSetClient
			trMeta.CollationConnection = tr.CollationConnection
		}
		tableProto.Triggers = append(tableProto.Triggers, trMeta)
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

	for _, sName := range sortedMapKeys(db.Sequences) {
		seq := db.Sequences[sName]
		schemaMeta.Sequences = append(schemaMeta.Sequences, &metadata.SequenceMetadata{
			Name:      seq.Name,
			DataType:  "BIGINT",
			Start:     strconv.FormatInt(seq.Start, 10),
			MinValue:  strconv.FormatInt(seq.MinValue, 10),
			MaxValue:  strconv.FormatInt(seq.MaxValue, 10),
			Increment: strconv.FormatInt(seq.Increment, 10),
			CacheSize: strconv.FormatInt(seq.Cache, 10),
			Cycle:     seq.Cycle,
			Comment:   seq.Comment,
			SkipDump:  seq.SkipDump,
		})
	}

	for _, fName := range sortedMapKeys(db.Functions) {
		f := db.Functions[fName]
		fnMeta := &metadata.FunctionMetadata{
			Name:              f.Name,
			Definition:        c.ShowCreateFunction(dbName, f.Name),
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
			Definition:        c.ShowCreateProcedure(dbName, p.Name),
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
	var nonJunkTargets []*TargetEntryQ
	if leftQ := leftmostAnalyzedQuery(v.AnalyzedQuery); leftQ != nil {
		for _, te := range leftQ.TargetList {
			if te != nil && !te.ResJunk {
				nonJunkTargets = append(nonJunkTargets, te)
			}
		}
	}
	n := len(v.Columns)
	if n == 0 {
		n = len(v.ColumnMetadata)
	}
	if len(nonJunkTargets) > n {
		n = len(nonJunkTargets)
	}
	if n == 0 {
		return nil
	}
	cols := make([]*metadata.ColumnMetadata, 0, n)
	for i := 0; i < n; i++ {
		name := ""
		if i < len(v.Columns) {
			name = v.Columns[i]
		} else if i < len(v.ColumnMetadata) {
			name = v.ColumnMetadata[i].Name
		}
		var te *TargetEntryQ
		if i < len(nonJunkTargets) {
			te = nonJunkTargets[i]
		}
		if (name == "" || name == "*" || strings.HasSuffix(name, ".*")) && te != nil && te.ResName != "" {
			name = te.ResName
		}
		col := &metadata.ColumnMetadata{
			Name:     name,
			Position: int32(i + 1),
			// TiDB information_schema.COLUMNS reports IS_NULLABLE=YES and NULL
			// charset/collation for every view column.
			Nullable: true,
		}
		if v.AnalyzedQuery != nil && v.AnalyzedQuery.SetOp != SetOpNone {
			col.Type = resolveSetOpColumnType(c, db, v.AnalyzedQuery, i)
		} else if te != nil {
			if te.ResOrigTable != "" && te.ResOrigCol != "" {
				origDB := db
				if te.ResOrigDB != "" {
					origDB = c.GetDatabase(te.ResOrigDB)
				}
				if origDB != nil {
					if tbl := origDB.GetTable(te.ResOrigTable); tbl != nil {
						if srcCol := tbl.GetColumn(te.ResOrigCol); srcCol != nil {
							col.Type = srcCol.ColumnType
						}
					}
				}
			}
			if col.Type == "" && te.Expr != nil {
				col.Type = resolveAnalyzedExprColumnType(c, db, v.AnalyzedQuery, te.Expr)
			}
		}
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

func resolveSetOpColumnType(c *Catalog, db *Database, q *Query, idx int) string {
	t, _ := resolveSetOpColumnTypeExpr(c, db, q, idx)
	return t
}

func resolveSetOpColumnTypeExpr(c *Catalog, db *Database, q *Query, idx int) (string, AnalyzedExpr) {
	if q == nil {
		return "", nil
	}
	if q.SetOp == SetOpNone {
		te := nthNonJunkTarget(q, idx)
		if te == nil {
			return "", nil
		}
		if te.ResOrigTable != "" && te.ResOrigCol != "" {
			origDB := db
			if te.ResOrigDB != "" {
				origDB = c.GetDatabase(te.ResOrigDB)
			}
			if origDB != nil {
				if tbl := origDB.GetTable(te.ResOrigTable); tbl != nil {
					if srcCol := tbl.GetColumn(te.ResOrigCol); srcCol != nil {
						return srcCol.ColumnType, te.Expr
					}
				}
			}
		}
		if te.Expr != nil {
			return resolveAnalyzedExprColumnType(c, db, q, te.Expr), te.Expr
		}
		return "", nil
	}
	lType, lExpr := resolveSetOpColumnTypeExpr(c, db, q.LArg, idx)
	// TiDB v8.5.5 keeps the left branch's column type for INTERSECT/EXCEPT views
	// and coerces across arms for UNION.
	if q.SetOp != SetOpUnion {
		return lType, lExpr
	}
	rType, rExpr := resolveSetOpColumnTypeExpr(c, db, q.RArg, idx)
	if lType == "" {
		return rType, rExpr
	}
	if rType == "" {
		return lType, lExpr
	}
	return combineSetOpColumnTypes(lType, lExpr, rType, rExpr), lExpr
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
		w := max(lp.width, rp.width)
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
		w := max(lp.width, rp.width)
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
			if lp.cat == "int" && lPrec < 19 {
				lPrec = 19
			}
			if rp.cat == "int" && rPrec < 19 {
				rPrec = 19
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
	return fmt.Sprintf("varchar(%d)", max(lp.width, rp.width))
}

func resolveAnalyzedExprColumnType(c *Catalog, db *Database, q *Query, expr AnalyzedExpr) string {
	if expr == nil {
		return ""
	}
	if cn, ok := expr.(*ConstExprQ); ok {
		if cn.IsNull {
			return "binary(0)"
		}
		if cn.IsString {
			return fmt.Sprintf("var_string(%d)", utf8.RuneCountInString(cn.Value))
		}
	}
	if rt := resolveExprResolvedType(c, db, q, expr); rt != nil {
		if rt.BaseType == BaseTypeTinyIntBool {
			// TiDB information_schema.COLUMNS reports comparisons and boolean
			// expressions in views as bigint rather than tinyint(1).
			return "bigint"
		}
		return resolvedTypeToColumnType(rt)
	}
	return ""
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
						if sv.AnalyzedQuery.SetOp == SetOpUnion {
							if rt := parseColumnResolvedType(resolveSetOpColumnType(c, origDB, sv.AnalyzedQuery, e.AttNum-1)); rt != nil {
								return rt
							}
						}
						if subTE := nthNonJunkTarget(sv.AnalyzedQuery, e.AttNum-1); subTE != nil {
							return resolveExprResolvedType(c, origDB, sv.AnalyzedQuery, subTE.Expr)
						}
					}
				}
			}
			if rte.Kind == RTESubquery && rte.Subquery != nil {
				if rte.Subquery.SetOp == SetOpUnion {
					if rt := parseColumnResolvedType(resolveSetOpColumnType(c, db, rte.Subquery, e.AttNum-1)); rt != nil {
						return rt
					}
				}
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
		switch opUpper {
		case "=", "!=", "<>", "<", "<=", ">", ">=", "<=>",
			"LIKE", "NOT LIKE", "REGEXP", "RLIKE", "NOT REGEXP", "NOT RLIKE",
			"IS TRUE", "IS NOT TRUE", "IS FALSE", "IS NOT FALSE", "IS UNKNOWN", "IS NOT UNKNOWN":
			return BoolType
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
	default:
		return nil
	}
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
		Name:                  t.Name,
		Engine:                t.Engine,
		Charset:               t.Charset,
		Collation:             t.Collation,
		Comment:               formatSyncedComment(t.Comment),
		IsGlobalTemporary:     t.GlobalTemporary,
		OnCommit:              t.OnCommit,
		PlacementPolicy:       t.PlacementPolicy,
		TtlColumn:             t.TTLColumn,
		TtlInterval:           t.TTLInterval,
		TtlEnable:             t.TTLEnable,
		TtlJobInterval:        t.TTLJobInterval,
		AutoIdCache:           t.AutoIDCache,
		PreSplitRegions:       int32(t.PreSplitRegions),
		TiflashReplica:        int32(t.TiFlashReplica),
		TiflashLocationLabels: slices.Clone(t.TiFlashLocationLabels),
		SkipDump:              t.SkipDump,
	}
	var opts []string
	if t.Compression != "" {
		opts = append(opts, fmt.Sprintf("compression='%s'", escapeComment(strings.Trim(t.Compression, "'\"`"))))
	}
	if t.Cached {
		opts = append(opts, "cached=on")
	}
	if len(opts) > 0 {
		table.CreateOptions = strings.Join(opts, " ")
	}

	var autoRandCol *Column
	for _, col := range t.Columns {
		colMeta := &metadata.ColumnMetadata{
			Name:         col.Name,
			Position:     int32(col.Position),
			Type:         col.ColumnType,
			Nullable:     col.Nullable,
			CharacterSet: col.Charset,
			Collation:    col.Collation,
			Comment:      formatSyncedComment(col.Comment),
		}
		switch {
		case col.AutoRandom:
			autoRandCol = col
			switch {
			case col.AutoRandomShardBits == 0 && col.AutoRandomRangeBits == 0:
				colMeta.Default = autoRandomSymbol
			case col.AutoRandomRangeBits != 0:
				colMeta.Default = fmt.Sprintf("%s(%d, %d)", autoRandomSymbol, col.AutoRandomShardBits, col.AutoRandomRangeBits)
			default:
				colMeta.Default = fmt.Sprintf("%s(%d)", autoRandomSymbol, col.AutoRandomShardBits)
			}
		case col.AutoIncrement:
			colMeta.Default = autoIncrementSymbol
		case col.Default != nil:
			if strings.EqualFold(*col.Default, "NULL") {
				colMeta.Default = "NULL"
			} else {
				colMeta.Default = formatDefault(*col.Default, col)
			}
		default:
		}
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
		colMeta.IsInvisible = col.Invisible
		colMeta.DefaultDropped = col.DefaultDropped
		table.Columns = append(table.Columns, colMeta)
	}

	switch {
	case t.ShardRowIDBits > 0:
		table.ShardingInfo = fmt.Sprintf("SHARD_BITS=%d", t.ShardRowIDBits)
	case autoRandCol != nil && autoRandCol.AutoRandomShardBits > 0:
		if autoRandCol.AutoRandomRangeBits != 0 {
			table.ShardingInfo = fmt.Sprintf("PK_AUTO_RANDOM_BITS=%d, RANGE BITS=%d", autoRandCol.AutoRandomShardBits, autoRandCol.AutoRandomRangeBits)
		} else {
			table.ShardingInfo = fmt.Sprintf("PK_AUTO_RANDOM_BITS=%d", autoRandCol.AutoRandomShardBits)
		}
	default:
	}

	sortedIndexes := slices.Clone(t.Indexes)
	slices.SortStableFunc(sortedIndexes, func(a, b *Index) int {
		return cmp.Compare(a.Name, b.Name)
	})
	for _, idx := range sortedIndexes {
		indexType := strings.ToUpper(idx.IndexType)
		if indexType == "" {
			indexType = "BTREE"
		}
		idxMeta := &metadata.IndexMetadata{
			Name:         idx.Name,
			Type:         indexType,
			Unique:       idx.Unique,
			Primary:      idx.Primary,
			Visible:      idx.Visible,
			Comment:      idx.Comment,
			KeyBlockSize: int64(idx.KeyBlockSize),
			IsVector:     idx.Vector,
			IsGlobal:     idx.Global,
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
			idxMeta.KeyLength = append(idxMeta.KeyLength, keyLength)
			idxMeta.Descending = append(idxMeta.Descending, col.Descending)
		}
		table.Indexes = append(table.Indexes, idxMeta)
	}

	for _, con := range t.Constraints {
		switch con.Type {
		case ConPrimaryKey:
			if con.Clustered != nil && table.PrimaryKeyType == "" {
				if *con.Clustered {
					table.PrimaryKeyType = "CLUSTERED"
				} else {
					table.PrimaryKeyType = "NONCLUSTERED"
				}
			}
		case ConForeignKey:
			// Omitted actions default to NO ACTION in TiDB REFERENTIAL_CONSTRAINTS (scenarios_c21_test.go).
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
	if p.SubType != "" && p.NumSubParts > 0 {
		subUseDefault = strconv.Itoa(p.NumSubParts)
	}

	if len(p.Partitions) == 0 && p.IntervalText != "" {
		return []*metadata.TablePartitionMetadata{{
			Type:           partType,
			Expression:     expr,
			IntervalClause: p.IntervalText,
			Algorithm:      int32(p.Algorithm),
		}}
	}

	var result []*metadata.TablePartitionMetadata
	for _, pd := range p.Partitions {
		part := &metadata.TablePartitionMetadata{
			Name:           pd.Name,
			Type:           partType,
			Expression:     expr,
			UseDefault:     partUseDefault,
			Value:          pd.ValueExpr,
			Comment:        pd.Comment,
			Algorithm:      int32(p.Algorithm),
			IntervalClause: p.IntervalText,
		}
		if p.SubType != "" {
			if len(pd.SubPartitions) > 0 {
				for _, spd := range pd.SubPartitions {
					part.Subpartitions = append(part.Subpartitions, &metadata.TablePartitionMetadata{
						Name:       spd.Name,
						Type:       subType,
						Expression: subExpr,
						UseDefault: subUseDefault,
						Comment:    spd.Comment,
						Algorithm:  int32(p.SubAlgo),
					})
				}
			} else {
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
