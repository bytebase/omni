package catalog

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

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

	schemaMeta := &metadata.SchemaMetadata{
		Name: "",
	}

	for _, tName := range sortedMapKeys(db.Tables) {
		schemaMeta.Tables = append(schemaMeta.Tables, tableToProto(db.Tables[tName]))
	}

	for _, vName := range sortedMapKeys(db.Views) {
		v := db.Views[vName]
		schemaMeta.Views = append(schemaMeta.Views, &metadata.ViewMetadata{
			Name:       v.Name,
			Definition: v.Definition,
			Columns:    viewColumnsToProto(v),
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
		})
	}

	dbMeta.Schemas = append(dbMeta.Schemas, schemaMeta)
	return dbMeta
}

func viewColumnsToProto(v *View) []*metadata.ColumnMetadata {
	if len(v.ColumnMetadata) > 0 {
		cols := make([]*metadata.ColumnMetadata, 0, len(v.ColumnMetadata))
		for i, cm := range v.ColumnMetadata {
			cols = append(cols, &metadata.ColumnMetadata{
				Name:         cm.Name,
				Position:     int32(i + 1),
				CharacterSet: cm.Charset,
				Collation:    cm.Collation,
			})
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
		Name:      t.Name,
		Engine:    t.Engine,
		Collation: t.Collation,
		Comment:   t.Comment,
	}
	if t.Cached {
		table.CreateOptions = "cached=on"
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
			Comment:      col.Comment,
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
			colMeta.Default = *col.Default
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

	for _, idx := range t.Indexes {
		indexType := strings.ToUpper(idx.IndexType)
		if indexType == "" {
			indexType = "BTREE"
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

func partitionsToProto(p *PartitionInfo) []*metadata.TablePartitionMetadata {
	expr := p.Expr
	if expr == "" && len(p.Columns) > 0 {
		expr = strings.Join(p.Columns, ",")
	}
	partType := partitionTypeToProto(p.Type, p.Linear)
	var result []*metadata.TablePartitionMetadata
	for _, pd := range p.Partitions {
		result = append(result, &metadata.TablePartitionMetadata{
			Name:       pd.Name,
			Type:       partType,
			Expression: expr,
			Value:      pd.ValueExpr,
			Comment:    pd.Comment,
		})
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

func sortedMapKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
