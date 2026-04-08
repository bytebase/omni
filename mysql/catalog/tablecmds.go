package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

func (c *Catalog) createTable(stmt *nodes.CreateTableStmt) error {
	// Resolve database.
	dbName := ""
	if stmt.Table != nil {
		dbName = stmt.Table.Schema
	}
	db, err := c.resolveDatabase(dbName)
	if err != nil {
		return err
	}

	tableName := stmt.Table.Name
	key := toLower(tableName)

	// Check for duplicate table or view with the same name.
	if db.Tables[key] != nil {
		if stmt.IfNotExists {
			return nil
		}
		return errDupTable(tableName)
	}
	if db.Views[key] != nil {
		return errDupTable(tableName)
	}

	// CREATE TABLE ... LIKE
	if stmt.Like != nil {
		return c.createTableLike(db, tableName, key, stmt)
	}

	// CREATE TABLE ... AS SELECT (CTAS) — not supported yet, skip silently
	if stmt.Select != nil && len(stmt.Columns) == 0 {
		return nil
	}

	tbl := &Table{
		Name:      tableName,
		Database:  db,
		Columns:   make([]*Column, 0, len(stmt.Columns)),
		colByName: make(map[string]int),
		Indexes:   make([]*Index, 0),
		Constraints: make([]*Constraint, 0),
		Charset:   db.Charset,
		Collation: db.Collation,
		Engine:    "InnoDB",
		Temporary: stmt.Temporary,
	}

	// Apply table options.
	tblCharsetExplicit := false
	tblCollationExplicit := false
	for _, opt := range stmt.Options {
		switch toLower(opt.Name) {
		case "engine":
			tbl.Engine = opt.Value
		case "charset", "character set", "default charset", "default character set":
			tbl.Charset = opt.Value
			tblCharsetExplicit = true
		case "collate", "default collate":
			tbl.Collation = opt.Value
			tblCollationExplicit = true
		case "comment":
			tbl.Comment = opt.Value
		case "auto_increment":
			fmt.Sscanf(opt.Value, "%d", &tbl.AutoIncrement)
		case "row_format":
			tbl.RowFormat = opt.Value
		case "key_block_size":
			fmt.Sscanf(opt.Value, "%d", &tbl.KeyBlockSize)
		}
	}
	// When charset is specified without explicit collation, derive the default collation.
	if tblCharsetExplicit && !tblCollationExplicit {
		if dc, ok := defaultCollationForCharset[toLower(tbl.Charset)]; ok {
			tbl.Collation = dc
		}
	}
	// When collation is specified without explicit charset, derive the charset from collation.
	if tblCollationExplicit && !tblCharsetExplicit {
		if cs := charsetForCollation(tbl.Collation); cs != "" {
			tbl.Charset = cs
		}
	}
	// Track whether we have a primary key (to detect multiple PKs).
	hasPK := false

	// Process columns.
	for i, colDef := range stmt.Columns {
		colKey := toLower(colDef.Name)
		if _, exists := tbl.colByName[colKey]; exists {
			return errDupColumn(colDef.Name)
		}

		col := &Column{
			Position: i + 1,
			Name:     colDef.Name,
			Nullable: true, // default nullable
		}

		// Type info.
		isSerial := false
		if colDef.TypeName != nil {
			typeName := toLower(colDef.TypeName.Name)
			// Handle SERIAL: expands to BIGINT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE
			if typeName == "serial" {
				isSerial = true
				col.DataType = "bigint"
				col.ColumnType = "bigint unsigned"
				col.AutoIncrement = true
				col.Nullable = false
			} else if typeName == "boolean" {
				col.DataType = "tinyint"
				col.ColumnType = formatColumnType(colDef.TypeName)
			} else if typeName == "numeric" {
				col.DataType = "decimal"
				col.ColumnType = formatColumnType(colDef.TypeName)
			} else {
				col.DataType = typeName
				col.ColumnType = formatColumnType(colDef.TypeName)
			}
			if colDef.TypeName.Charset != "" {
				col.Charset = colDef.TypeName.Charset
			}
			if colDef.TypeName.Collate != "" {
				col.Collation = colDef.TypeName.Collate
			}

			// MySQL converts string types with CHARACTER SET binary to binary types.
			// ENUM and SET are not converted — they keep CHARACTER SET binary annotation.
			if strings.EqualFold(col.Charset, "binary") && isStringType(col.DataType) && !isEnumSetType(col.DataType) {
				col = convertToBinaryType(col, colDef.TypeName)
			}
		}

		// Default charset/collation for string types.
		if isStringType(col.DataType) {
			if col.Charset == "" {
				col.Charset = tbl.Charset
			}
			if col.Collation == "" {
				// If column charset differs from table charset, use the default
				// collation for the column's charset, not the table's collation.
				if !strings.EqualFold(col.Charset, tbl.Charset) {
					if dc, ok := defaultCollationForCharset[toLower(col.Charset)]; ok {
						col.Collation = dc
					}
				} else {
					col.Collation = tbl.Collation
				}
			}
		}

		// Top-level column properties.
		if colDef.TypeName != nil && colDef.TypeName.SRID != 0 {
			col.SRID = colDef.TypeName.SRID
		}
		if colDef.AutoIncrement {
			col.AutoIncrement = true
			col.Nullable = false
		}
		if colDef.Comment != "" {
			col.Comment = colDef.Comment
		}
		if colDef.DefaultValue != nil {
			s := nodeToSQL(colDef.DefaultValue)
			col.Default = &s
		}
		if colDef.OnUpdate != nil {
			col.OnUpdate = nodeToSQL(colDef.OnUpdate)
		}
		if colDef.Generated != nil {
			col.Generated = &GeneratedColumnInfo{
				Expr:   nodeToSQLGenerated(colDef.Generated.Expr, tbl.Charset),
				Stored: colDef.Generated.Stored,
			}
		}

		// Process column-level constraints.
		for _, cc := range colDef.Constraints {
			switch cc.Type {
			case nodes.ColConstrNotNull:
				col.Nullable = false
			case nodes.ColConstrNull:
				col.Nullable = true
			case nodes.ColConstrDefault:
				if cc.Expr != nil {
					s := nodeToSQL(cc.Expr)
					col.Default = &s
				}
			case nodes.ColConstrPrimaryKey:
				if hasPK {
					return errMultiplePriKey()
				}
				hasPK = true
				col.Nullable = false
				// Add PK index and constraint after all columns are processed.
				// We'll defer this—record it for now.
			case nodes.ColConstrUnique:
				// Handled after columns are added.
			case nodes.ColConstrAutoIncrement:
				col.AutoIncrement = true
				col.Nullable = false
			case nodes.ColConstrCheck:
				// Add check constraint.
				conName := cc.Name
				if conName == "" {
					conName = fmt.Sprintf("%s_chk_%d", tableName, nextCheckNumber(tbl))
				}
				tbl.Constraints = append(tbl.Constraints, &Constraint{
					Name:        conName,
					Type:        ConCheck,
					Table:       tbl,
					CheckExpr:   nodeToSQL(cc.Expr),
					NotEnforced: cc.NotEnforced,
				})
			case nodes.ColConstrReferences:
				// Column-level FK.
				refDB := ""
				refTable := ""
				if cc.RefTable != nil {
					refDB = cc.RefTable.Schema
					refTable = cc.RefTable.Name
				}
				conName := cc.Name
				if conName == "" {
					conName = fmt.Sprintf("%s_ibfk_%d", tableName, countFKConstraints(tbl)+1)
				}
				tbl.Constraints = append(tbl.Constraints, &Constraint{
					Name:       conName,
					Type:       ConForeignKey,
					Table:      tbl,
					Columns:    []string{colDef.Name},
					RefDatabase: refDB,
					RefTable:   refTable,
					RefColumns: cc.RefColumns,
					OnDelete:   refActionToString(cc.OnDelete),
					OnUpdate:   refActionToString(cc.OnUpdate),
				})
				// Add implicit backing index for FK if needed.
				ensureFKBackingIndex(tbl, cc.Name, []string{colDef.Name}, []*IndexColumn{{Name: colDef.Name}})
			case nodes.ColConstrVisible:
				col.Invisible = false
			case nodes.ColConstrInvisible:
				col.Invisible = true
			case nodes.ColConstrCollate:
				// Collation specified via constraint.
				if cc.Expr != nil {
					if s, ok := cc.Expr.(*nodes.StringLit); ok {
						col.Collation = s.Value
					}
				}
			}
		}

		// SERIAL implies UNIQUE KEY — add after the column is fully configured.
		if isSerial {
			idxName := allocIndexName(tbl, colDef.Name)
			tbl.Indexes = append(tbl.Indexes, &Index{
				Name:      idxName,
				Table:     tbl,
				Columns:   []*IndexColumn{{Name: colDef.Name}},
				Unique:    true,
				IndexType: "",
				Visible:   true,
			})
			tbl.Constraints = append(tbl.Constraints, &Constraint{
				Name:      idxName,
				Type:      ConUniqueKey,
				Table:     tbl,
				Columns:   []string{colDef.Name},
				IndexName: idxName,
			})
		}

		tbl.Columns = append(tbl.Columns, col)
		tbl.colByName[colKey] = i
	}

	// Second pass: add column-level PK and UNIQUE indexes/constraints.
	for _, colDef := range stmt.Columns {
		for _, cc := range colDef.Constraints {
			switch cc.Type {
			case nodes.ColConstrPrimaryKey:
				tbl.Indexes = append(tbl.Indexes, &Index{
					Name:      "PRIMARY",
					Table:     tbl,
					Columns:   []*IndexColumn{{Name: colDef.Name}},
					Unique:    true,
					Primary:   true,
					IndexType: "",
					Visible:   true,
				})
				tbl.Constraints = append(tbl.Constraints, &Constraint{
					Name:      "PRIMARY",
					Type:      ConPrimaryKey,
					Table:     tbl,
					Columns:   []string{colDef.Name},
					IndexName: "PRIMARY",
				})
			case nodes.ColConstrUnique:
				idxName := allocIndexName(tbl, colDef.Name)
				tbl.Indexes = append(tbl.Indexes, &Index{
					Name:      idxName,
					Table:     tbl,
					Columns:   []*IndexColumn{{Name: colDef.Name}},
					Unique:    true,
					IndexType: "",
					Visible:   true,
				})
				tbl.Constraints = append(tbl.Constraints, &Constraint{
					Name:      idxName,
					Type:      ConUniqueKey,
					Table:     tbl,
					Columns:   []string{colDef.Name},
					IndexName: idxName,
				})
			}
		}
	}

	// Process table-level constraints.
	for _, con := range stmt.Constraints {
		cols := extractColumnNames(con)

		switch con.Type {
		case nodes.ConstrPrimaryKey:
			if hasPK {
				return errMultiplePriKey()
			}
			hasPK = true
			// Mark PK columns as NOT NULL.
			for _, colName := range cols {
				c := tbl.GetColumn(colName)
				if c != nil {
					c.Nullable = false
				}
			}
			idxCols := buildIndexColumns(con)
			pkIdx := &Index{
				Name:      "PRIMARY",
				Table:     tbl,
				Columns:   idxCols,
				Unique:    true,
				Primary:   true,
				IndexType: "",
				Visible:   true,
			}
			applyIndexOptions(pkIdx, con.IndexOptions)
			tbl.Indexes = append(tbl.Indexes, pkIdx)
			tbl.Constraints = append(tbl.Constraints, &Constraint{
				Name:      "PRIMARY",
				Type:      ConPrimaryKey,
				Table:     tbl,
				Columns:   cols,
				IndexName: "PRIMARY",
			})

		case nodes.ConstrUnique:
			idxName := con.Name
			if idxName == "" && len(cols) > 0 {
				idxName = allocIndexName(tbl, cols[0])
			}
			idxCols := buildIndexColumns(con)
			uqIdx := &Index{
				Name:      idxName,
				Table:     tbl,
				Columns:   idxCols,
				Unique:    true,
				IndexType: resolveConstraintIndexType(con),
				Visible:   true,
			}
			applyIndexOptions(uqIdx, con.IndexOptions)
			tbl.Indexes = append(tbl.Indexes, uqIdx)
			tbl.Constraints = append(tbl.Constraints, &Constraint{
				Name:      idxName,
				Type:      ConUniqueKey,
				Table:     tbl,
				Columns:   cols,
				IndexName: idxName,
			})

		case nodes.ConstrForeignKey:
			conName := con.Name
			if conName == "" {
				conName = fmt.Sprintf("%s_ibfk_%d", tableName, countFKConstraints(tbl)+1)
			}
			refDB := ""
			refTable := ""
			if con.RefTable != nil {
				refDB = con.RefTable.Schema
				refTable = con.RefTable.Name
			}
			tbl.Constraints = append(tbl.Constraints, &Constraint{
				Name:       conName,
				Type:       ConForeignKey,
				Table:      tbl,
				Columns:    cols,
				RefDatabase: refDB,
				RefTable:   refTable,
				RefColumns: con.RefColumns,
				OnDelete:   refActionToString(con.OnDelete),
				OnUpdate:   refActionToString(con.OnUpdate),
			})
			// Add implicit backing index for FK if needed.
			ensureFKBackingIndex(tbl, con.Name, cols, buildIndexColumns(con))

		case nodes.ConstrCheck:
			conName := con.Name
			if conName == "" {
				conName = fmt.Sprintf("%s_chk_%d", tableName, nextCheckNumber(tbl))
			}
			tbl.Constraints = append(tbl.Constraints, &Constraint{
				Name:        conName,
				Type:        ConCheck,
				Table:       tbl,
				CheckExpr:   nodeToSQL(con.Expr),
				NotEnforced: con.NotEnforced,
			})

		case nodes.ConstrIndex:
			idxName := con.Name
			if idxName == "" && len(cols) > 0 {
				idxName = allocIndexName(tbl, cols[0])
			}
			idxCols := buildIndexColumns(con)
			keyIdx := &Index{
				Name:      idxName,
				Table:     tbl,
				Columns:   idxCols,
				IndexType: resolveConstraintIndexType(con),
				Visible:   true,
			}
			applyIndexOptions(keyIdx, con.IndexOptions)
			tbl.Indexes = append(tbl.Indexes, keyIdx)

		case nodes.ConstrFulltextIndex:
			idxName := con.Name
			if idxName == "" && len(cols) > 0 {
				idxName = allocIndexName(tbl, cols[0])
			}
			idxCols := buildIndexColumns(con)
			ftIdx := &Index{
				Name:      idxName,
				Table:     tbl,
				Columns:   idxCols,
				Fulltext:  true,
				IndexType: "FULLTEXT",
				Visible:   true,
			}
			applyIndexOptions(ftIdx, con.IndexOptions)
			tbl.Indexes = append(tbl.Indexes, ftIdx)

		case nodes.ConstrSpatialIndex:
			idxName := con.Name
			if idxName == "" && len(cols) > 0 {
				idxName = allocIndexName(tbl, cols[0])
			}
			idxCols := buildIndexColumns(con)
			spIdx := &Index{
				Name:      idxName,
				Table:     tbl,
				Columns:   idxCols,
				Spatial:   true,
				IndexType: "SPATIAL",
				Visible:   true,
			}
			applyIndexOptions(spIdx, con.IndexOptions)
			tbl.Indexes = append(tbl.Indexes, spIdx)
		}
	}

	// Validate foreign key constraints (unless foreign_key_checks=0).
	if c.foreignKeyChecks {
		if err := c.validateForeignKeys(db, tbl); err != nil {
			return err
		}
	}

	// Process partition clause.
	if stmt.Partitions != nil {
		tbl.Partitioning = buildPartitionInfo(stmt.Partitions)
	}

	db.Tables[key] = tbl
	return nil
}

// buildPartitionInfo converts an AST PartitionClause to a catalog PartitionInfo.
func buildPartitionInfo(pc *nodes.PartitionClause) *PartitionInfo {
	pi := &PartitionInfo{
		Linear:   pc.Linear,
		NumParts: pc.NumParts,
	}

	switch pc.Type {
	case nodes.PartitionRange:
		if len(pc.Columns) > 0 {
			pi.Type = "RANGE COLUMNS"
			pi.Columns = pc.Columns
		} else {
			pi.Type = "RANGE"
			pi.Expr = nodeToSQL(pc.Expr)
		}
	case nodes.PartitionList:
		if len(pc.Columns) > 0 {
			pi.Type = "LIST COLUMNS"
			pi.Columns = pc.Columns
		} else {
			pi.Type = "LIST"
			pi.Expr = nodeToSQL(pc.Expr)
		}
	case nodes.PartitionHash:
		pi.Type = "HASH"
		pi.Expr = nodeToSQL(pc.Expr)
	case nodes.PartitionKey:
		pi.Type = "KEY"
		pi.Columns = pc.Columns
		pi.Algorithm = pc.Algorithm
	}

	// Subpartition info.
	if pc.SubPartType != 0 || pc.SubPartExpr != nil || len(pc.SubPartColumns) > 0 {
		switch pc.SubPartType {
		case nodes.PartitionHash:
			pi.SubType = "HASH"
			pi.SubExpr = nodeToSQL(pc.SubPartExpr)
		case nodes.PartitionKey:
			pi.SubType = "KEY"
			pi.SubColumns = pc.SubPartColumns
			pi.SubAlgo = pc.SubPartAlgo
		}
		pi.SubLinear = false // TODO: track linear for subpartitions if parser supports it
		pi.NumSubParts = pc.NumSubParts
	}

	// Partition definitions.
	for _, pd := range pc.Partitions {
		pdi := &PartitionDefInfo{
			Name: pd.Name,
		}
		// Values.
		if pd.Values != nil {
			pdi.ValueExpr = partitionValueToString(pd.Values, pc.Type)
		}
		// Options.
		for _, opt := range pd.Options {
			switch toLower(opt.Name) {
			case "engine":
				pdi.Engine = opt.Value
			case "comment":
				pdi.Comment = opt.Value
			}
		}
		// Subpartitions.
		for _, spd := range pd.SubPartitions {
			spdi := &SubPartitionDefInfo{
				Name: spd.Name,
			}
			for _, opt := range spd.Options {
				switch toLower(opt.Name) {
				case "engine":
					spdi.Engine = opt.Value
				case "comment":
					spdi.Comment = opt.Value
				}
			}
			pdi.SubPartitions = append(pdi.SubPartitions, spdi)
		}
		pi.Partitions = append(pi.Partitions, pdi)
	}

	return pi
}

// partitionValueToString converts a partition value node to SQL string.
func partitionValueToString(v nodes.Node, ptype nodes.PartitionType) string {
	switch n := v.(type) {
	case *nodes.String:
		if n.Str == "MAXVALUE" {
			return "MAXVALUE"
		}
		return n.Str
	case *nodes.List:
		parts := make([]string, len(n.Items))
		for i, item := range n.Items {
			parts[i] = nodeToSQL(item.(nodes.ExprNode))
		}
		return strings.Join(parts, ",")
	case nodes.ExprNode:
		return nodeToSQL(n)
	default:
		return ""
	}
}

// validateForeignKeys checks all FK constraints on a table against the referenced tables.
// It validates: (1) referenced table exists, (2) referenced columns have an index,
// (3) column types are compatible.
func (c *Catalog) validateForeignKeys(db *Database, tbl *Table) error {
	for _, con := range tbl.Constraints {
		if con.Type != ConForeignKey {
			continue
		}
		if err := c.validateSingleFK(db, tbl, con); err != nil {
			return err
		}
	}
	return nil
}

// validateSingleFK validates a single FK constraint against its referenced table.
func (c *Catalog) validateSingleFK(db *Database, tbl *Table, con *Constraint) error {
	// Resolve the referenced table.
	refDBName := con.RefDatabase
	var refDB *Database
	if refDBName != "" {
		refDB = c.GetDatabase(refDBName)
	} else {
		refDB = db
	}

	var refTbl *Table
	if refDB != nil {
		// Self-referencing FK: the table being created references itself.
		if toLower(con.RefTable) == toLower(tbl.Name) && refDB == db {
			refTbl = tbl
		} else {
			refTbl = refDB.GetTable(con.RefTable)
		}
	}

	if refTbl == nil {
		return errFKNoRefTable(con.RefTable)
	}

	// Check that referenced columns have an index (PK or UNIQUE or KEY)
	// that starts with the referenced columns in order.
	if !hasIndexOnColumns(refTbl, con.RefColumns) {
		return errFKMissingIndex(con.Name, con.RefTable)
	}

	// Check column type compatibility.
	for i, colName := range con.Columns {
		if i >= len(con.RefColumns) {
			break
		}
		col := tbl.GetColumn(colName)
		refCol := refTbl.GetColumn(con.RefColumns[i])
		if col == nil || refCol == nil {
			continue
		}
		if !fkTypesCompatible(col, refCol) {
			return errFKIncompatibleColumns(colName, con.RefColumns[i], con.Name)
		}
	}

	return nil
}

// hasIndexOnColumns checks whether a table has an index (PK, UNIQUE, or regular KEY)
// whose leading columns match the given columns.
func hasIndexOnColumns(tbl *Table, cols []string) bool {
	for _, idx := range tbl.Indexes {
		if len(idx.Columns) < len(cols) {
			continue
		}
		match := true
		for i, col := range cols {
			if toLower(idx.Columns[i].Name) != toLower(col) {
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

// fkTypesCompatible checks whether two columns have compatible types for FK relationships.
// MySQL requires that FK and referenced columns have the same storage type.
func fkTypesCompatible(col, refCol *Column) bool {
	// Compare base data types.
	if col.DataType != refCol.DataType {
		return false
	}

	// For integer types, check signedness (unsigned must match).
	colUnsigned := strings.Contains(strings.ToLower(col.ColumnType), "unsigned")
	refUnsigned := strings.Contains(strings.ToLower(refCol.ColumnType), "unsigned")
	if colUnsigned != refUnsigned {
		return false
	}

	// For string types, check charset compatibility.
	if isStringType(col.DataType) {
		colCharset := col.Charset
		refCharset := refCol.Charset
		if colCharset != "" && refCharset != "" && toLower(colCharset) != toLower(refCharset) {
			return false
		}
	}

	return true
}

// extractColumnNames returns column names from an AST constraint.
func extractColumnNames(con *nodes.Constraint) []string {
	if len(con.IndexColumns) > 0 {
		names := make([]string, 0, len(con.IndexColumns))
		for _, ic := range con.IndexColumns {
			if cr, ok := ic.Expr.(*nodes.ColumnRef); ok {
				names = append(names, cr.Column)
			}
		}
		return names
	}
	return con.Columns
}

// buildIndexColumns converts AST IndexColumn list to catalog IndexColumn list.
func buildIndexColumns(con *nodes.Constraint) []*IndexColumn {
	if len(con.IndexColumns) > 0 {
		result := make([]*IndexColumn, 0, len(con.IndexColumns))
		for _, ic := range con.IndexColumns {
			idxCol := &IndexColumn{
				Length:     ic.Length,
				Descending: ic.Desc,
			}
			if cr, ok := ic.Expr.(*nodes.ColumnRef); ok {
				idxCol.Name = cr.Column
			} else {
				idxCol.Expr = nodeToSQL(ic.Expr)
			}
			result = append(result, idxCol)
		}
		return result
	}
	// Fallback to simple column names.
	result := make([]*IndexColumn, 0, len(con.Columns))
	for _, name := range con.Columns {
		result = append(result, &IndexColumn{Name: name})
	}
	return result
}

// allocIndexName generates a unique index name based on the first column,
// appending _2, _3, etc. on collision.
func allocIndexName(tbl *Table, baseName string) string {
	candidate := baseName
	suffix := 2
	for indexNameExists(tbl, candidate) {
		candidate = fmt.Sprintf("%s_%d", baseName, suffix)
		suffix++
	}
	return candidate
}

// hasIndexCoveringColumns returns true if the table already has an index whose
// leading columns match the given FK columns (left-prefix match). MySQL 8.0
// reuses such an index instead of creating an implicit backing index for the FK.
func hasIndexCoveringColumns(tbl *Table, fkCols []string) bool {
	for _, idx := range tbl.Indexes {
		if len(idx.Columns) < len(fkCols) {
			continue
		}
		match := true
		for i, col := range fkCols {
			if !strings.EqualFold(idx.Columns[i].Name, col) {
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

// ensureFKBackingIndex creates an implicit backing index for FK columns
// if no existing index already covers them (MySQL 8.0 behavior).
// MySQL uses the constraint name as the index name when provided;
// otherwise falls back to the first column name via allocIndexName.
func ensureFKBackingIndex(tbl *Table, conName string, cols []string, idxCols []*IndexColumn) {
	if hasIndexCoveringColumns(tbl, cols) {
		return
	}
	idxName := conName
	if idxName == "" {
		idxName = allocIndexName(tbl, cols[0])
	}
	tbl.Indexes = append(tbl.Indexes, &Index{
		Name:    idxName,
		Table:   tbl,
		Columns: idxCols,
		Visible: true,
	})
}

func indexNameExists(tbl *Table, name string) bool {
	key := toLower(name)
	for _, idx := range tbl.Indexes {
		if toLower(idx.Name) == key {
			return true
		}
	}
	return false
}

func indexTypeOrDefault(indexType, defaultType string) string {
	if indexType != "" {
		return indexType
	}
	return defaultType
}

// resolveConstraintIndexType returns the index type from a constraint,
// checking both IndexType (USING before key parts) and IndexOptions (USING after key parts).
func resolveConstraintIndexType(con *nodes.Constraint) string {
	if con.IndexType != "" {
		return strings.ToUpper(con.IndexType)
	}
	for _, opt := range con.IndexOptions {
		if strings.EqualFold(opt.Name, "USING") {
			if s, ok := opt.Value.(*nodes.StringLit); ok {
				return strings.ToUpper(s.Value)
			}
		}
	}
	return ""
}

// applyIndexOptions extracts COMMENT, VISIBLE/INVISIBLE, and KEY_BLOCK_SIZE
// from AST IndexOptions and applies them to the given Index.
func applyIndexOptions(idx *Index, opts []*nodes.IndexOption) {
	for _, opt := range opts {
		switch strings.ToUpper(opt.Name) {
		case "COMMENT":
			if s, ok := opt.Value.(*nodes.StringLit); ok {
				idx.Comment = s.Value
			}
		case "VISIBLE":
			idx.Visible = true
		case "INVISIBLE":
			idx.Visible = false
		case "KEY_BLOCK_SIZE":
			if n, ok := opt.Value.(*nodes.Integer); ok {
				idx.KeyBlockSize = int(n.Ival)
			}
		}
	}
}

// countFKConstraints counts the number of foreign key constraints on a table.
func countFKConstraints(tbl *Table) int {
	count := 0
	for _, c := range tbl.Constraints {
		if c.Type == ConForeignKey {
			count++
		}
	}
	return count
}

// nextCheckNumber returns the next available check constraint number for auto-naming.
// MySQL uses tableName_chk_N where N starts at 1 and increments, skipping existing names.
func nextCheckNumber(tbl *Table) int {
	n := 1
	for {
		name := fmt.Sprintf("%s_chk_%d", tbl.Name, n)
		exists := false
		for _, c := range tbl.Constraints {
			if toLower(c.Name) == toLower(name) {
				exists = true
				break
			}
		}
		if !exists {
			return n
		}
		n++
	}
}

func isStringType(dt string) bool {
	switch dt {
	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext",
		"enum", "set":
		return true
	}
	return false
}

// convertToBinaryType converts a string-type column with CHARACTER SET binary
// to the equivalent binary type (char->binary, varchar->varbinary, text->blob, etc.).
func convertToBinaryType(col *Column, dt *nodes.DataType) *Column {
	switch col.DataType {
	case "char":
		col.DataType = "binary"
		length := dt.Length
		if length == 0 {
			length = 1
		}
		col.ColumnType = fmt.Sprintf("binary(%d)", length)
	case "varchar":
		col.DataType = "varbinary"
		col.ColumnType = fmt.Sprintf("varbinary(%d)", dt.Length)
	case "tinytext":
		col.DataType = "tinyblob"
		col.ColumnType = "tinyblob"
	case "text":
		col.DataType = "blob"
		col.ColumnType = "blob"
	case "mediumtext":
		col.DataType = "mediumblob"
		col.ColumnType = "mediumblob"
	case "longtext":
		col.DataType = "longblob"
		col.ColumnType = "longblob"
	}
	// Binary types don't have charset/collation in SHOW CREATE TABLE.
	col.Charset = ""
	col.Collation = ""
	return col
}

// nodeToSQLGenerated converts an AST expression to SQL for use in a generated
// column definition. MySQL prefixes string literals with a charset introducer
// (e.g., _utf8mb4'value') in generated column expressions.
func nodeToSQLGenerated(node nodes.ExprNode, charset string) string {
	if node == nil {
		return ""
	}
	switch n := node.(type) {
	case *nodes.ColumnRef:
		if n.Table != "" {
			return "`" + n.Table + "`.`" + n.Column + "`"
		}
		return "`" + n.Column + "`"
	case *nodes.IntLit:
		return fmt.Sprintf("%d", n.Value)
	case *nodes.StringLit:
		// MySQL adds charset introducer for string literals in generated columns.
		if charset != "" {
			return "_" + charset + "'" + n.Value + "'"
		}
		return "'" + n.Value + "'"
	case *nodes.FuncCallExpr:
		funcName := strings.ToLower(n.Name)
		if n.Star {
			return funcName + "(*)"
		}
		var args []string
		for _, a := range n.Args {
			args = append(args, nodeToSQLGenerated(a, charset))
		}
		return funcName + "(" + strings.Join(args, ",") + ")"
	case *nodes.NullLit:
		return "NULL"
	case *nodes.BoolLit:
		if n.Value {
			return "1"
		}
		return "0"
	case *nodes.FloatLit:
		return n.Value
	case *nodes.BitLit:
		val := strings.TrimLeft(n.Value, "0")
		if val == "" {
			val = "0"
		}
		return "b'" + val + "'"
	case *nodes.ParenExpr:
		return "(" + nodeToSQLGenerated(n.Expr, charset) + ")"
	case *nodes.BinaryExpr:
		left := nodeToSQLGenerated(n.Left, charset)
		right := nodeToSQLGenerated(n.Right, charset)
		// MySQL rewrites JSON operators to function calls in generated column expressions.
		switch n.Op {
		case nodes.BinOpJsonExtract:
			return "json_extract(" + left + "," + right + ")"
		case nodes.BinOpJsonUnquote:
			return "json_unquote(json_extract(" + left + "," + right + "))"
		}
		op := binaryOpToString(n.Op)
		return "(" + left + " " + op + " " + right + ")"
	case *nodes.UnaryExpr:
		operand := nodeToSQLGenerated(n.Operand, charset)
		switch n.Op {
		case nodes.UnaryMinus:
			return "-" + operand
		case nodes.UnaryNot:
			return "NOT " + operand
		case nodes.UnaryBitNot:
			return "~" + operand
		default:
			return operand
		}
	default:
		return "(?)"
	}
}

func nodeToSQL(node nodes.ExprNode) string {
	if node == nil {
		return ""
	}
	switch n := node.(type) {
	case *nodes.ColumnRef:
		if n.Table != "" {
			return "`" + n.Table + "`.`" + n.Column + "`"
		}
		return "`" + n.Column + "`"
	case *nodes.IntLit:
		return fmt.Sprintf("%d", n.Value)
	case *nodes.StringLit:
		return "'" + n.Value + "'"
	case *nodes.FuncCallExpr:
		funcName := strings.ToLower(n.Name)
		if n.Star {
			return funcName + "(*)"
		}
		var args []string
		for _, a := range n.Args {
			args = append(args, nodeToSQL(a))
		}
		return funcName + "(" + strings.Join(args, ",") + ")"
	case *nodes.NullLit:
		return "NULL"
	case *nodes.BoolLit:
		if n.Value {
			return "1"
		}
		return "0"
	case *nodes.FloatLit:
		return n.Value
	case *nodes.BitLit:
		// MySQL strips leading zeros from bit literals in SHOW CREATE TABLE.
		val := strings.TrimLeft(n.Value, "0")
		if val == "" {
			val = "0"
		}
		return "b'" + val + "'"
	case *nodes.ParenExpr:
		return "(" + nodeToSQL(n.Expr) + ")"
	case *nodes.BinaryExpr:
		left := nodeToSQL(n.Left)
		right := nodeToSQL(n.Right)
		op := binaryOpToString(n.Op)
		// MySQL wraps binary expressions in parentheses in SHOW CREATE TABLE.
		return "(" + left + " " + op + " " + right + ")"
	case *nodes.UnaryExpr:
		operand := nodeToSQL(n.Operand)
		switch n.Op {
		case nodes.UnaryMinus:
			return "-" + operand
		case nodes.UnaryNot:
			return "NOT " + operand
		case nodes.UnaryBitNot:
			return "~" + operand
		default:
			return operand
		}
	default:
		return "(?)"
	}
}

func binaryOpToString(op nodes.BinaryOp) string {
	switch op {
	case nodes.BinOpAdd:
		return "+"
	case nodes.BinOpSub:
		return "-"
	case nodes.BinOpMul:
		return "*"
	case nodes.BinOpDiv:
		return "/"
	case nodes.BinOpMod:
		return "%"
	case nodes.BinOpEq:
		return "="
	case nodes.BinOpNe:
		return "!="
	case nodes.BinOpLt:
		return "<"
	case nodes.BinOpGt:
		return ">"
	case nodes.BinOpLe:
		return "<="
	case nodes.BinOpGe:
		return ">="
	case nodes.BinOpAnd:
		return "and"
	case nodes.BinOpOr:
		return "or"
	case nodes.BinOpBitAnd:
		return "&"
	case nodes.BinOpBitOr:
		return "|"
	case nodes.BinOpBitXor:
		return "^"
	case nodes.BinOpShiftLeft:
		return "<<"
	case nodes.BinOpShiftRight:
		return ">>"
	case nodes.BinOpDivInt:
		return "DIV"
	default:
		return "?"
	}
}

func formatColumnType(dt *nodes.DataType) string {
	name := strings.ToLower(dt.Name)

	// MySQL type aliases: BOOLEAN/BOOL → tinyint(1), NUMERIC → decimal, SERIAL → bigint unsigned
	// GEOMETRYCOLLECTION → geomcollection (MySQL 8.0 normalized form)
	switch name {
	case "boolean":
		return "tinyint(1)"
	case "numeric":
		name = "decimal"
	case "serial":
		return "bigint unsigned"
	case "geometrycollection":
		name = "geomcollection"
	}

	var buf strings.Builder
	buf.WriteString(name)

	// Integer display width handling for MySQL 8.0:
	// - Display width is deprecated and NOT shown by default
	// - EXCEPTION: When ZEROFILL is used, MySQL 8.0 still shows the display width
	//   with default widths per type: tinyint(3), smallint(5), mediumint(8), int(10), bigint(20)
	isIntType := isIntegerType(name)
	if isIntType {
		if dt.Zerofill {
			width := dt.Length
			if width == 0 {
				width = defaultIntDisplayWidth(name, dt.Unsigned)
			}
			fmt.Fprintf(&buf, "(%d)", width)
		}
		// Non-zerofill integer types: strip display width (MySQL 8.0 deprecated)
	} else if name == "decimal" && dt.Length == 0 && dt.Scale == 0 {
		// DECIMAL with no precision → MySQL shows decimal(10,0)
		buf.WriteString("(10,0)")
	} else if isTextBlobLengthStripped(name) {
		// TEXT(n) and BLOB(n) — MySQL stores the length internally but
		// SHOW CREATE TABLE displays just TEXT / BLOB without the length.
		// Do not emit length.
	} else if name == "year" {
		// YEAR(4) is deprecated in MySQL 8.0 — SHOW CREATE TABLE shows just `year`.
	} else if (name == "char" || name == "binary") && dt.Length == 0 {
		// CHAR/BINARY with no length → MySQL shows char(1)/binary(1)
		buf.WriteString("(1)")
	} else if dt.Length > 0 && dt.Scale > 0 {
		fmt.Fprintf(&buf, "(%d,%d)", dt.Length, dt.Scale)
	} else if dt.Length > 0 {
		fmt.Fprintf(&buf, "(%d)", dt.Length)
	}

	if len(dt.EnumValues) > 0 {
		buf.WriteString("(")
		for i, v := range dt.EnumValues {
			if i > 0 {
				buf.WriteString(",")
			}
			buf.WriteString("'" + escapeEnumValue(v) + "'")
		}
		buf.WriteString(")")
	}
	if dt.Unsigned {
		buf.WriteString(" unsigned")
	}
	if dt.Zerofill {
		buf.WriteString(" zerofill")
	}
	return buf.String()
}

// isIntegerType returns true for MySQL integer types.
func isIntegerType(dt string) bool {
	switch dt {
	case "tinyint", "smallint", "mediumint", "int", "integer", "bigint":
		return true
	}
	return false
}

// defaultIntDisplayWidth returns the default display width for integer types
// when ZEROFILL is used. These are the MySQL defaults.
func defaultIntDisplayWidth(typeName string, unsigned bool) int {
	switch typeName {
	case "tinyint":
		if unsigned {
			return 3
		}
		return 4
	case "smallint":
		if unsigned {
			return 5
		}
		return 6
	case "mediumint":
		if unsigned {
			return 8
		}
		return 9
	case "int", "integer":
		if unsigned {
			return 10
		}
		return 11
	case "bigint":
		if unsigned {
			return 20
		}
		return 20
	}
	return 11
}

// isTextBlobLengthStripped returns true for types where MySQL strips the length
// in SHOW CREATE TABLE output (TEXT(n) → text, BLOB(n) → blob).
func isTextBlobLengthStripped(dt string) bool {
	switch dt {
	case "text", "blob":
		return true
	}
	return false
}

// escapeEnumValue escapes single quotes in ENUM/SET values for SHOW CREATE TABLE.
// MySQL uses '' (two single quotes) to escape a single quote in enum values.
func escapeEnumValue(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// createTableLike implements CREATE TABLE t2 LIKE t1.
// It copies the structure (columns, indexes, constraints) from the source table.
func (c *Catalog) createTableLike(db *Database, tableName, key string, stmt *nodes.CreateTableStmt) error {
	// Resolve source table.
	srcDBName := stmt.Like.Schema
	srcDB, err := c.resolveDatabase(srcDBName)
	if err != nil {
		return err
	}
	srcTbl := srcDB.GetTable(stmt.Like.Name)
	if srcTbl == nil {
		return errNoSuchTable(srcDB.Name, stmt.Like.Name)
	}

	tbl := &Table{
		Name:      tableName,
		Database:  db,
		Columns:   make([]*Column, 0, len(srcTbl.Columns)),
		colByName: make(map[string]int),
		Indexes:   make([]*Index, 0, len(srcTbl.Indexes)),
		Constraints: make([]*Constraint, 0, len(srcTbl.Constraints)),
		Engine:    srcTbl.Engine,
		Charset:   srcTbl.Charset,
		Collation: srcTbl.Collation,
		Comment:   srcTbl.Comment,
		RowFormat: srcTbl.RowFormat,
		KeyBlockSize: srcTbl.KeyBlockSize,
		Temporary: stmt.Temporary,
	}

	// Copy columns.
	for i, srcCol := range srcTbl.Columns {
		col := &Column{
			Position:      srcCol.Position,
			Name:          srcCol.Name,
			DataType:      srcCol.DataType,
			ColumnType:    srcCol.ColumnType,
			Nullable:      srcCol.Nullable,
			AutoIncrement: srcCol.AutoIncrement,
			Charset:       srcCol.Charset,
			Collation:     srcCol.Collation,
			Comment:       srcCol.Comment,
			OnUpdate:      srcCol.OnUpdate,
			Invisible:     srcCol.Invisible,
		}
		if srcCol.Default != nil {
			def := *srcCol.Default
			col.Default = &def
		}
		if srcCol.Generated != nil {
			col.Generated = &GeneratedColumnInfo{
				Expr:   srcCol.Generated.Expr,
				Stored: srcCol.Generated.Stored,
			}
		}
		tbl.Columns = append(tbl.Columns, col)
		tbl.colByName[toLower(col.Name)] = i
	}

	// Copy indexes.
	for _, srcIdx := range srcTbl.Indexes {
		idx := &Index{
			Name:         srcIdx.Name,
			Table:        tbl,
			Unique:       srcIdx.Unique,
			Primary:      srcIdx.Primary,
			Fulltext:     srcIdx.Fulltext,
			Spatial:      srcIdx.Spatial,
			IndexType:    srcIdx.IndexType,
			Visible:      srcIdx.Visible,
			Comment:      srcIdx.Comment,
			KeyBlockSize: srcIdx.KeyBlockSize,
		}
		cols := make([]*IndexColumn, len(srcIdx.Columns))
		for i, sc := range srcIdx.Columns {
			cols[i] = &IndexColumn{
				Name:       sc.Name,
				Length:     sc.Length,
				Descending: sc.Descending,
				Expr:       sc.Expr,
			}
		}
		idx.Columns = cols
		tbl.Indexes = append(tbl.Indexes, idx)
	}

	// Copy constraints (skip FK — MySQL 8.0 does not copy FKs with LIKE).
	for _, srcCon := range srcTbl.Constraints {
		if srcCon.Type == ConForeignKey {
			continue
		}
		con := &Constraint{
			Name:        srcCon.Name,
			Type:        srcCon.Type,
			Table:       tbl,
			Columns:     append([]string{}, srcCon.Columns...),
			IndexName:   srcCon.IndexName,
			CheckExpr:   srcCon.CheckExpr,
			NotEnforced: srcCon.NotEnforced,
			RefDatabase: srcCon.RefDatabase,
			RefTable:    srcCon.RefTable,
			RefColumns:  append([]string{}, srcCon.RefColumns...),
			OnDelete:    srcCon.OnDelete,
			OnUpdate:    srcCon.OnUpdate,
		}
		tbl.Constraints = append(tbl.Constraints, con)
	}

	db.Tables[key] = tbl
	return nil
}

func refActionToString(action nodes.ReferenceAction) string {
	switch action {
	case nodes.RefActRestrict:
		return "RESTRICT"
	case nodes.RefActCascade:
		return "CASCADE"
	case nodes.RefActSetNull:
		return "SET NULL"
	case nodes.RefActSetDefault:
		return "SET DEFAULT"
	case nodes.RefActNoAction:
		return "NO ACTION"
	default:
		return "NO ACTION"
	}
}
