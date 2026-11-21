package catalog

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/bytebase/omni/metadata"
	"google.golang.org/protobuf/proto"
)

// DatabaseMetadata converts the named database in the catalog to a
// DatabaseSchemaMetadata snapshot matching the MySQL-mode OceanBase sync and
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

	// OceanBase MySQL mode uses a single empty-name schema.
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
		tableProto := tableToProto(t)
		tableProto.HasMaterializedViewLog = db.MaterializedViewLogs[strings.ToLower(t.Name)]
		schemaMeta.Tables = append(schemaMeta.Tables, tableProto)
		tableProtoByLowerName[strings.ToLower(t.Name)] = tableProto
	}

	// Triggers are database-level in OceanBase MySQL mode (db.Triggers), each bound
	// to a table via Trigger.Table; the proto nests them under TableMetadata.Triggers
	// in ACTION_ORDER, matching live sync.
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
		if v.Materialized {
			mvMeta := &metadata.MaterializedViewMetadata{
				Name:       v.Name,
				Definition: v.Definition,
				Comment:    formatSyncedComment(v.Comment),
				SkipDump:   v.SkipDump,
			}
			sortedMVIndexes := slices.Clone(v.Indexes)
			slices.SortStableFunc(sortedMVIndexes, func(a, b *Index) int {
				return cmp.Compare(a.Name, b.Name)
			})
			for _, idx := range sortedMVIndexes {
				mvMeta.Indexes = append(mvMeta.Indexes, indexToProto(idx, ""))
			}
			schemaMeta.MaterializedViews = append(schemaMeta.MaterializedViews, mvMeta)
			continue
		}
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
	leftQ := leftmostAnalyzedQuery(v.AnalyzedQuery)
	var nonJunkTargets []*TargetEntryQ
	if leftQ != nil {
		for _, te := range leftQ.TargetList {
			if te != nil && !te.ResJunk {
				nonJunkTargets = append(nonJunkTargets, te)
			}
		}
	}
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
		return resolveTargetColumnMetadata(c, db, leftQ, nonJunkTargets[i])
	}
	fillTargetCol := func(col *metadata.ColumnMetadata, i int, te *TargetEntryQ) {
		if colType, nullable, charset, collation, ok := resolveSourceAttrs(i); ok {
			if col.Type == "" {
				col.Type = colType
			}
			// OceanBase keeps base-column nullability across outer joins (viewcmds.go).
			col.Nullable = nullable
			if charset != "" {
				col.CharacterSet = charset
			}
			if collation != "" {
				col.Collation = collation
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

func resolveTargetColumnMetadata(c *Catalog, db *Database, q *Query, te *TargetEntryQ) (colType string, nullable bool, charset, collation string, ok bool) {
	if te == nil {
		return "", false, "", "", false
	}
	if te.ResOrigTable != "" && te.ResOrigCol != "" {
		origDB := db
		if te.ResOrigDB != "" && c != nil {
			origDB = c.GetDatabase(te.ResOrigDB)
		}
		if origDB != nil {
			if tbl := origDB.GetTable(te.ResOrigTable); tbl != nil {
				if srcCol := tbl.GetColumn(te.ResOrigCol); srcCol != nil {
					ci := viewColumnCollation(srcCol)
					return srcCol.ColumnType, srcCol.Nullable, ci.charset, ci.collation, true
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
				return srcCol.ColumnType, srcCol.Nullable, ci.charset, ci.collation, true
			}
		}
		if sv := origDB.Views[strings.ToLower(rte.TableName)]; sv != nil {
			if innerQ := leftmostAnalyzedQuery(sv.AnalyzedQuery); innerQ != nil {
				if subTE := nthNonJunkTarget(innerQ, ve.AttNum-1); subTE != nil {
					if ct, nl, cs, cl, found := resolveTargetColumnMetadata(c, origDB, innerQ, subTE); found {
						return ct, nl, cs, cl, true
					}
				}
			}
			if ve.AttNum-1 < len(sv.ColumnMetadata) {
				cm := sv.ColumnMetadata[ve.AttNum-1]
				return resolvedTypeToColumnType(cm.Type), cm.Nullable, cm.Charset, cm.Collation, true
			}
		}
	case RTESubquery:
		if innerQ := leftmostAnalyzedQuery(rte.Subquery); innerQ != nil {
			if subTE := nthNonJunkTarget(innerQ, ve.AttNum-1); subTE != nil {
				if ct, nl, cs, cl, found := resolveTargetColumnMetadata(c, db, innerQ, subTE); found {
					return ct, nl, cs, cl, true
				}
			}
		}
	}
	return "", false, "", "", false
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

func tableToProto(t *Table) *metadata.TableMetadata {
	table := &metadata.TableMetadata{
		Name:         t.Name,
		Engine:       t.Engine,
		Charset:      t.Charset,
		Collation:    t.Collation,
		Comment:      formatSyncedComment(t.Comment),
		ColumnGroups: t.ColumnGroups.String(),
		Tablespace:   t.Tablespace,
		SkipDump:     t.SkipDump,
	}
	if opts := tableCreateOptions(t); opts != "" {
		table.CreateOptions = opts
	}

	for _, col := range t.Columns {
		if col.Hidden == ColumnHiddenSystem {
			// System-hidden columns (e.g. backing columns of functional
			// indexes) are not reported by information_schema.
			continue
		}
		colMeta := &metadata.ColumnMetadata{
			Name:         col.Name,
			Position:     int32(col.Position),
			Type:         col.ColumnType,
			Nullable:     col.Nullable,
			CharacterSet: col.Charset,
			Collation:    col.Collation,
			Comment:      formatSyncedComment(col.Comment),
		}
		if col.Charset != "" {
			colMeta.CharsetInherited = proto.Bool(!col.CharsetExplicit)
		}
		if col.Collation != "" {
			colMeta.CollationInherited = proto.Bool(!col.CollationExplicit)
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
		table.Indexes = append(table.Indexes, indexToProto(idx, t.Engine))
	}

	for _, con := range t.Constraints {
		switch con.Type {
		case ConForeignKey:
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
			if con.RefDatabase != "" {
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
	if t.Partitioning != nil {
		opts = append(opts, "partitioned")
	}
	return strings.Join(opts, " ")
}

func normalizeDefault(value string, kind ColumnDefaultKind) string {
	switch kind {
	case ColumnDefaultConstant:
		if strings.EqualFold(value, "NULL") {
			return "NULL"
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
	if p.NumParts > 0 && (p.Type == "HASH" || p.Type == "KEY") && (p.SubType == "" || p.UseDefaultSubpartitions) {
		partUseDefault = strconv.Itoa(p.NumParts)
	}
	var subUseDefault string
	if p.SubType != "" && p.UseDefaultSubpartitions && !p.UseDefaultNumSubpartitions && p.NumSubParts > 0 {
		subUseDefault = strconv.Itoa(p.NumSubParts)
	}

	if len(p.Partitions) == 0 && (p.AutoSize != "" || p.SplitSize != "") {
		part := &metadata.TablePartitionMetadata{
			Type:       partType,
			Expression: expr,
			UseDefault: partUseDefault,
			SplitSize:  p.SplitSize,
			AutoSize:   p.AutoSize,
			Algorithm:  int32(p.Algorithm),
		}
		if p.SubType != "" {
			part.Subpartitions = []*metadata.TablePartitionMetadata{{
				Type:       subType,
				Expression: subExpr,
				UseDefault: subUseDefault,
				Algorithm:  int32(p.SubAlgo),
			}}
		}
		if len(p.SubPartTemplate) > 0 && !p.UseDefaultSubpartitions {
			part.SubpartitionTemplate = subpartitionTemplateToProto(p.SubPartTemplate, subType, subExpr, p.SubAlgo)
		}
		return []*metadata.TablePartitionMetadata{part}
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
			SplitSize:  p.SplitSize,
			AutoSize:   p.AutoSize,
			Algorithm:  int32(p.Algorithm),
		}
		if p.SubType != "" {
			if len(pd.SubPartitions) > 0 {
				for _, spd := range pd.SubPartitions {
					part.Subpartitions = append(part.Subpartitions, &metadata.TablePartitionMetadata{
						Name:       spd.Name,
						Type:       subType,
						Expression: subExpr,
						UseDefault: subUseDefault,
						Value:      spd.ValueExpr,
						Comment:    spd.Comment,
						Algorithm:  int32(p.SubAlgo),
					})
				}
			} else if subUseDefault != "" {
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
	if len(result) > 0 && len(p.SubPartTemplate) > 0 && !p.UseDefaultSubpartitions {
		result[0].SubpartitionTemplate = subpartitionTemplateToProto(p.SubPartTemplate, subType, subExpr, p.SubAlgo)
	}
	return result
}

func subpartitionTemplateToProto(tpl []*SubPartitionDefInfo, subType metadata.TablePartitionMetadata_Type, subExpr string, subAlgo int) []*metadata.TablePartitionMetadata {
	out := make([]*metadata.TablePartitionMetadata, 0, len(tpl))
	for _, st := range tpl {
		out = append(out, &metadata.TablePartitionMetadata{
			Name:       st.Name,
			Type:       subType,
			Expression: subExpr,
			Value:      st.ValueExpr,
			Comment:    st.Comment,
			Algorithm:  int32(subAlgo),
		})
	}
	return out
}

func indexToProto(idx *Index, engine string) *metadata.IndexMetadata {
	indexType := strings.ToUpper(idx.IndexType)
	if indexType == "" && !idx.Vector && !idx.Fulltext && !idx.Spatial {
		indexType = "BTREE"
		if strings.EqualFold(engine, "MEMORY") || strings.EqualFold(engine, "HEAP") {
			indexType = "HASH"
		}
	}
	idxMeta := &metadata.IndexMetadata{
		Name:             idx.Name,
		Type:             indexType,
		Unique:           idx.Unique,
		Primary:          idx.Primary,
		Visible:          idx.Visible,
		Comment:          idx.Comment,
		KeyBlockSize:     int64(idx.KeyBlockSize),
		IsVector:         idx.Vector,
		IsGlobal:         strings.EqualFold(idx.Scope, "GLOBAL"),
		Scope:            idx.Scope,
		StoringColumns:   slices.Clone(idx.Storing),
		VectorParams:     slices.Clone(idx.VectorParams),
		Parser:           idx.Parser,
		ParserProperties: maps.Clone(idx.ParserProperties),
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
			keyLength = 32
		}
		idxMeta.KeyLength = append(idxMeta.KeyLength, keyLength)
		idxMeta.Descending = append(idxMeta.Descending, col.Descending)
	}
	return idxMeta
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
