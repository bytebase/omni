package catalog

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

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
			Columns:     viewColumnsToProto(v),
			Algorithm:   v.Algorithm,
			SqlSecurity: v.SqlSecurity,
			CheckOption: v.CheckOption,
		})
	}

	for _, fName := range sortedMapKeys(db.Functions) {
		f := db.Functions[fName]
		fnMeta := &metadata.FunctionMetadata{
			Name:              f.Name,
			Definition:        c.ShowCreateFunction(dbName, f.Name),
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
		procMeta := &metadata.ProcedureMetadata{
			Name:              p.Name,
			Definition:        c.ShowCreateProcedure(dbName, p.Name),
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

func viewColumnsToProto(v *View) []*metadata.ColumnMetadata {
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
				col.Type = deparseResolvedTypeQ(cm.Type)
			}
			cols = append(cols, col)
		}
		return cols
	}
	if len(v.Columns) == 0 {
		return nil
	}
	cols := make([]*metadata.ColumnMetadata, 0, len(v.Columns))
	for i, name := range v.Columns {
		cols = append(cols, &metadata.ColumnMetadata{
			Name:     name,
			Position: int32(i + 1),
		})
	}
	return cols
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
			fk := &metadata.ForeignKeyMetadata{
				Name:              con.Name,
				Columns:           slices.Clone(con.Columns),
				ReferencedTable:   con.RefTable,
				ReferencedColumns: slices.Clone(con.RefColumns),
				OnUpdate:          con.OnUpdate,
				OnDelete:          con.OnDelete,
				MatchType:         con.MatchType,
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
			idxMeta.Expressions = append(idxMeta.Expressions, col.Expr)
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

func sortedMapKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
