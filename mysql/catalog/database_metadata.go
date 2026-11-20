package catalog

import (
	"fmt"
	"slices"
	"strings"

	"github.com/bytebase/omni/metadata"
)

// DatabaseMetadata converts the named database in the catalog to a
// DatabaseSchemaMetadata snapshot matching the MySQL driver sync format.
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
		tableProto := tableToProto(t)
		schemaMeta.Tables = append(schemaMeta.Tables, tableProto)
		tableProtoByLowerName[strings.ToLower(t.Name)] = tableProto
	}

	// Triggers are database-level in MySQL (db.Triggers), each bound to a table via
	// Trigger.Table; the proto nests them under their owning TableMetadata.Triggers,
	// matching the live MySQL sync.
	for _, trName := range sortedMapKeys(db.Triggers) {
		tr := db.Triggers[trName]
		tableProto := tableProtoByLowerName[strings.ToLower(tr.Table)]
		if tableProto == nil {
			continue
		}
		tableProto.Triggers = append(tableProto.Triggers, triggerToProto(tr))
	}

	for _, vName := range sortedMapKeys(db.Views) {
		v := db.Views[vName]
		schemaMeta.Views = append(schemaMeta.Views, &metadata.ViewMetadata{
			Name:       v.Name,
			Definition: v.Definition,
		})
	}

	for _, fName := range sortedMapKeys(db.Functions) {
		f := db.Functions[fName]
		schemaMeta.Functions = append(schemaMeta.Functions, &metadata.FunctionMetadata{
			Name:       f.Name,
			Definition: f.Body,
		})
	}

	for _, pName := range sortedMapKeys(db.Procedures) {
		p := db.Procedures[pName]
		schemaMeta.Procedures = append(schemaMeta.Procedures, &metadata.ProcedureMetadata{
			Name:       p.Name,
			Definition: p.Body,
		})
	}

	for _, eName := range sortedMapKeys(db.Events) {
		e := db.Events[eName]
		schemaMeta.Events = append(schemaMeta.Events, &metadata.EventMetadata{
			Name:       e.Name,
			Definition: c.ShowCreateEvent(dbName, e.Name),
		})
	}

	dbMeta.Schemas = append(dbMeta.Schemas, schemaMeta)
	return dbMeta
}

func tableToProto(t *Table) *metadata.TableMetadata {
	table := &metadata.TableMetadata{
		Name:      t.Name,
		Engine:    t.Engine,
		Charset:   t.Charset,
		Collation: t.Collation,
		Comment:   t.Comment,
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
			Comment:      col.Comment,
		}
		switch {
		case col.AutoIncrement:
			colMeta.Default = autoIncrementSymbol
		case col.Default != nil:
			colMeta.Default = normalizeDefault(*col.Default, col.DefaultKind)
		case col.Nullable:
			colMeta.Default = "NULL"
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
		// Spatial SRID (MySQL 8.0): presence-carrying — an explicit SRID (including 0)
		// is distinct from no SRID.
		if col.HasSRID {
			srid := uint32(col.SRID)
			colMeta.Srid = &srid
		}
		colMeta.IsInvisible = col.Invisible
		table.Columns = append(table.Columns, colMeta)
	}

	for _, idx := range t.Indexes {
		indexType := strings.ToUpper(idx.IndexType)
		if indexType == "" {
			indexType = "BTREE"
			if strings.EqualFold(t.Engine, "MEMORY") || strings.EqualFold(t.Engine, "HEAP") {
				indexType = "HASH"
			}
		}
		idxMeta := &metadata.IndexMetadata{
			Name:    idx.Name,
			Type:    indexType,
			Unique:  idx.Unique,
			Primary: idx.Primary,
			Visible: idx.Visible,
			Comment: idx.Comment,
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
			fk := &metadata.ForeignKeyMetadata{
				Name:              con.Name,
				Columns:           con.Columns,
				ReferencedTable:   con.RefTable,
				ReferencedColumns: con.RefColumns,
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
				Name:       con.Name,
				Expression: con.CheckExpr,
			})
		default:
		}
	}

	if t.Partitioning != nil {
		table.Partitions = partitionsToProto(t.Partitioning)
	}

	return table
}

// normalizeDefault converts a catalog default value to the driver representation:
// static defaults are single-quoted for mysqldump compatibility and expression
// defaults are wrapped in parentheses.
func normalizeDefault(value string, kind ColumnDefaultKind) string {
	switch kind {
	case ColumnDefaultConstant:
		if strings.EqualFold(value, "NULL") {
			return "NULL"
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

	var result []*metadata.TablePartitionMetadata
	for _, pd := range p.Partitions {
		part := &metadata.TablePartitionMetadata{
			Name:       pd.Name,
			Type:       partType,
			Expression: expr,
			Value:      pd.ValueExpr,
			Comment:    pd.Comment,
		}
		if p.SubType != "" {
			for _, spd := range pd.SubPartitions {
				part.Subpartitions = append(part.Subpartitions, &metadata.TablePartitionMetadata{
					Name:       spd.Name,
					Type:       subType,
					Expression: subExpr,
					Comment:    spd.Comment,
				})
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
	return &metadata.TriggerMetadata{
		Name:   tr.Name,
		Event:  tr.Event,
		Timing: tr.Timing,
		Body:   tr.Body,
	}
}

func sortedMapKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
