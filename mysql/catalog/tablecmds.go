package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/deparse"
)

const generatedInvisiblePrimaryKeyColumnName = "my_row_id"

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
		Name:        tableName,
		Database:    db,
		Columns:     make([]*Column, 0, len(stmt.Columns)),
		colByName:   make(map[string]int),
		Indexes:     make([]*Index, 0),
		Constraints: make([]*Constraint, 0),
		Charset:     db.Charset,
		Collation:   db.Collation,
		Engine:      "",
		Temporary:   stmt.Temporary,
	}

	// Apply table options.
	tblCharsetExplicit := false
	tblCollationExplicit := false
	for _, opt := range stmt.Options {
		switch toLower(opt.Name) {
		case "engine":
			tbl.Engine = opt.Value
		case "charset", "character set", "default charset", "default character set":
			tbl.Charset = normalizeCharsetName(opt.Value)
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
		case "compression":
			tbl.Compression = opt.Value
		case "encryption":
			tbl.Encryption = opt.Value
		case "stats_persistent":
			tbl.StatsPersistent = opt.Value
		case "stats_auto_recalc":
			tbl.StatsAutoRecalc = opt.Value
		case "stats_sample_pages":
			tbl.StatsSamplePages = opt.Value
		case "min_rows":
			tbl.MinRows = opt.Value
		case "max_rows":
			tbl.MaxRows = opt.Value
		case "avg_row_length":
			tbl.AvgRowLength = opt.Value
		case "tablespace":
			tbl.Tablespace = opt.Value
		case "pack_keys":
			tbl.PackKeys = opt.Value
		case "checksum":
			tbl.Checksum = opt.Value
		case "delay_key_write":
			tbl.DelayKeyWrite = opt.Value
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
			tbl.Charset = normalizeCharsetName(cs)
		}
	}
	if tblCharsetExplicit && tblCollationExplicit && !charsetMatchesCollation(tbl.Charset, tbl.Collation) {
		return &Error{Code: 1253, SQLState: "42000", Message: "COLLATION is not valid for CHARACTER SET"}
	}
	// Track whether we have a primary key (to detect multiple PKs).
	hasPK := false

	// Defer FK backing index creation until after all explicit indexes are added,
	// so that explicit indexes can satisfy FK requirements without creating duplicates.
	type pendingFK struct {
		conName string
		cols    []string
		idxCols []*IndexColumn
	}
	var pendingFKs []pendingFK

	// unnamedFKCount counts FKs that received an auto-generated name in this
	// CREATE TABLE statement. Matches MySQL 8.0's generate_fk_name() counter,
	// which is initialized to 0 for CREATE TABLE and incremented per unnamed FK,
	// IGNORING user-named FKs. See sql/sql_table.cc:9252 (create_table_impl
	// is called with fk_max_generated_name_number = 0) and sql/sql_table.cc:5912
	// (generate_fk_name uses ++counter).
	//
	// Example: CREATE TABLE t (a INT, CONSTRAINT t_ibfk_5 FK, b INT, FK)
	//   → first unnamed FK gets t_ibfk_1 (not t_ibfk_2 or t_ibfk_6).
	// This differs from ALTER TABLE ADD FK, where the counter starts at
	// max(existing) — see altercmds.go.
	var unnamedFKCount int

	// unnamedCheckCount counts CHECK constraints that received an auto-generated
	// name in this CREATE TABLE statement. Matches MySQL 8.0's CHECK counter
	// at sql/sql_table.cc:19073 (cc_max_generated_number starts at 0, used
	// via ++cc_max_generated_number). Like the FK counter, this IGNORES
	// user-named CHECK constraints during CREATE TABLE.
	//
	// Example: CREATE TABLE t (a INT, CONSTRAINT t_chk_1 CHECK(a>0), b INT, CHECK(b<100))
	//   → unnamed CHECK gets t_chk_1, but t_chk_1 is already taken by user
	//   → real MySQL errors with ER_CHECK_CONSTRAINT_DUP_NAME
	//   (see sql/sql_table.cc:19595 check_constraint_dup_name check).
	// For ALTER TABLE ADD CHECK, the counter is loaded from existing max —
	// see altercmds.go which uses nextCheckNumber (gap-scan helper).
	var unnamedCheckCount int

	// Process columns.
	for i, colDef := range stmt.Columns {
		if err := validateColumnDefSemantics(colDef); err != nil {
			return err
		}
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
			typeName := normalizedColumnDataType(colDef.TypeName)
			// Handle SERIAL: expands to BIGINT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE
			if strings.EqualFold(colDef.TypeName.Name, "serial") {
				isSerial = true
				col.DataType = "bigint"
				col.ColumnType = "bigint unsigned"
				col.AutoIncrement = true
				col.Nullable = false
			} else if typeName == "boolean" || typeName == "bool" {
				col.DataType = "tinyint"
				col.ColumnType = formatColumnType(colDef.TypeName)
			} else {
				col.DataType = typeName
				col.ColumnType = formatColumnType(colDef.TypeName)
			}
			if isNationalStringType(colDef.TypeName.Name) {
				col.Charset = "utf8mb3"
			}
			if colDef.TypeName.Charset != "" {
				col.Charset = normalizeCharsetName(colDef.TypeName.Charset)
			}
			if colDef.TypeName.Collate != "" {
				col.Collation = colDef.TypeName.Collate
				if col.Charset == "" {
					col.Charset = normalizeCharsetName(charsetForCollation(col.Collation))
				}
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
			applyBinaryModifierCollation(col, colDef.TypeName)
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
					unnamedCheckCount++
					conName = fmt.Sprintf("%s_chk_%d", tableName, unnamedCheckCount)
				}
				if checkConstraintNameExists(db, tbl, conName) {
					return errCheckConstraintDupName(conName)
				}
				if err := validateColumnCheckReferences(colDef.Name, conName, cc.Expr); err != nil {
					return err
				}
				tbl.Constraints = append(tbl.Constraints, &Constraint{
					Name:        conName,
					Type:        ConCheck,
					Table:       tbl,
					CheckExpr:   nodeToSQL(cc.Expr),
					NotEnforced: cc.NotEnforced,
				})
			case nodes.ColConstrReferences:
				// MySQL parses but ignores column-level REFERENCES for non-NDB tables.
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
				if err := validatePrimaryKeyColumns(tbl, []string{colDef.Name}); err != nil {
					return err
				}
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
			if err := validatePrimaryKeyColumns(tbl, cols); err != nil {
				return err
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
			if !pkIdx.Visible {
				return &Error{Code: 3522, SQLState: "HY000", Message: "A primary key index cannot be invisible"}
			}
			tbl.Indexes = append(tbl.Indexes, pkIdx)
			tbl.Constraints = append(tbl.Constraints, &Constraint{
				Name:      "PRIMARY",
				Type:      ConPrimaryKey,
				Table:     tbl,
				Columns:   cols,
				IndexName: "PRIMARY",
			})

		case nodes.ConstrUnique:
			idxCols := buildIndexColumns(con)
			idxName := con.Name
			if idxName == "" {
				idxName = defaultIndexName(tbl, cols, idxCols)
			} else if idxName != "" {
				if err := validateNonPrimaryIndexName(idxName); err != nil {
					return err
				}
				if indexNameExists(tbl, idxName) {
					return errDupKeyName(idxName)
				}
			}
			if err := validateIndexColumns(tbl, idxCols, false, false); err != nil {
				return err
			}
			uqIdx := &Index{
				Name:      idxName,
				Table:     tbl,
				Columns:   idxCols,
				Unique:    true,
				IndexType: resolveConstraintIndexType(con),
				Visible:   true,
			}
			applyIndexOptions(uqIdx, con.IndexOptions)
			if err := synthesizeFunctionalIndexColumns(tbl, uqIdx); err != nil {
				return err
			}
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
				unnamedFKCount++
				conName = fmt.Sprintf("%s_ibfk_%d", tableName, unnamedFKCount)
			}
			if foreignKeyConstraintNameExists(db, tbl, conName) {
				return errFKDupName(conName)
			}
			refDB := ""
			refTable := ""
			if con.RefTable != nil {
				refDB = con.RefTable.Schema
				refTable = con.RefTable.Name
			}
			tbl.Constraints = append(tbl.Constraints, &Constraint{
				Name:        conName,
				Type:        ConForeignKey,
				Table:       tbl,
				Columns:     cols,
				RefDatabase: refDB,
				RefTable:    refTable,
				RefColumns:  con.RefColumns,
				OnDelete:    refActionToString(con.OnDelete),
				OnUpdate:    refActionToString(con.OnUpdate),
			})
			// Defer implicit backing index for FK until after all explicit indexes are added.
			pendingFKs = append(pendingFKs, pendingFK{conName: con.Name, cols: cols, idxCols: buildIndexColumns(con)})

		case nodes.ConstrCheck:
			conName := con.Name
			if conName == "" {
				unnamedCheckCount++
				conName = fmt.Sprintf("%s_chk_%d", tableName, unnamedCheckCount)
			}
			if checkConstraintNameExists(db, tbl, conName) {
				return errCheckConstraintDupName(conName)
			}
			if err := validateCheckExpr(conName, con.Expr); err != nil {
				return err
			}
			tbl.Constraints = append(tbl.Constraints, &Constraint{
				Name:        conName,
				Type:        ConCheck,
				Table:       tbl,
				CheckExpr:   nodeToSQL(con.Expr),
				NotEnforced: con.NotEnforced,
			})

		case nodes.ConstrIndex:
			idxCols := buildIndexColumns(con)
			idxName := con.Name
			if idxName == "" {
				idxName = defaultIndexName(tbl, cols, idxCols)
			} else if idxName != "" {
				if err := validateNonPrimaryIndexName(idxName); err != nil {
					return err
				}
				if indexNameExists(tbl, idxName) {
					return errDupKeyName(idxName)
				}
			}
			if err := validateIndexColumns(tbl, idxCols, false, false); err != nil {
				return err
			}
			keyIdx := &Index{
				Name:      idxName,
				Table:     tbl,
				Columns:   idxCols,
				IndexType: resolveConstraintIndexType(con),
				Visible:   true,
			}
			applyIndexOptions(keyIdx, con.IndexOptions)
			coerceInnoDBHashIndex(tbl, keyIdx)
			if err := synthesizeFunctionalIndexColumns(tbl, keyIdx); err != nil {
				return err
			}
			tbl.Indexes = append(tbl.Indexes, keyIdx)

		case nodes.ConstrFulltextIndex:
			idxName := con.Name
			if idxName == "" && len(cols) > 0 {
				idxName = allocIndexName(tbl, cols[0])
			} else if idxName != "" {
				if err := validateNonPrimaryIndexName(idxName); err != nil {
					return err
				}
				if indexNameExists(tbl, idxName) {
					return errDupKeyName(idxName)
				}
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
			} else if idxName != "" {
				if err := validateNonPrimaryIndexName(idxName); err != nil {
					return err
				}
				if indexNameExists(tbl, idxName) {
					return errDupKeyName(idxName)
				}
			}
			idxCols := buildIndexColumns(con)
			if err := validateIndexColumns(tbl, idxCols, false, true); err != nil {
				return err
			}
			if strings.EqualFold(resolveConstraintIndexType(con), "BTREE") {
				return &Error{Code: 3500, SQLState: "HY000", Message: "Index type not supported for spatial index"}
			}
			spIdx := &Index{
				Name:      idxName,
				Table:     tbl,
				Columns:   idxCols,
				Spatial:   true,
				IndexType: "SPATIAL",
				Visible:   true,
			}
			applyIndexOptions(spIdx, con.IndexOptions)
			if strings.EqualFold(spIdx.IndexType, "BTREE") {
				return &Error{Code: 3500, SQLState: "HY000", Message: "Index type not supported for spatial index"}
			}
			tbl.Indexes = append(tbl.Indexes, spIdx)
		}
	}

	if c.generateGIPK {
		if tbl.GetColumn(generatedInvisiblePrimaryKeyColumnName) != nil {
			return errGIPKColumnNameReserved()
		}
		if !hasPK {
			addGeneratedInvisiblePrimaryKey(tbl)
			hasPK = true
		}
	}

	// Process deferred FK backing indexes now that all explicit indexes are in place.
	for _, fk := range pendingFKs {
		ensureFKBackingIndex(tbl, fk.conName, fk.cols, fk.idxCols)
	}

	// Validate foreign key constraints (unless foreign_key_checks=0).
	if c.foreignKeyChecks {
		if err := c.validateForeignKeys(db, tbl); err != nil {
			return err
		}
	}

	// Process partition clause.
	if stmt.Partitions != nil {
		if err := validatePartitionClause(tbl, stmt.Partitions); err != nil {
			return err
		}
		tbl.Partitioning = buildPartitionInfo(tbl, stmt.Partitions)
	}

	// Phase 3: analyze DEFAULT, GENERATED, and CHECK expressions now that all
	// columns are present in the table.
	c.analyzeTableExpressions(tbl, stmt)

	db.Tables[key] = tbl
	return nil
}

func errGIPKColumnNameReserved() error {
	return &Error{
		Code:     4108,
		SQLState: "HY000",
		Message:  "Failed to generate invisible primary key. Column 'my_row_id' already exists.",
	}
}

func addGeneratedInvisiblePrimaryKey(tbl *Table) {
	col := &Column{
		Position:                     0,
		Name:                         generatedInvisiblePrimaryKeyColumnName,
		DataType:                     "bigint",
		ColumnType:                   "bigint unsigned",
		Nullable:                     false,
		AutoIncrement:                true,
		Invisible:                    true,
		GeneratedInvisiblePrimaryKey: true,
	}

	tbl.Columns = append([]*Column{col}, tbl.Columns...)
	rebuildColByNamePreservePositions(tbl)

	tbl.Indexes = append(tbl.Indexes, &Index{
		Name:      "PRIMARY",
		Table:     tbl,
		Columns:   []*IndexColumn{{Name: generatedInvisiblePrimaryKeyColumnName}},
		Unique:    true,
		Primary:   true,
		IndexType: "",
		Visible:   true,
	})
	tbl.Constraints = append(tbl.Constraints, &Constraint{
		Name:      "PRIMARY",
		Type:      ConPrimaryKey,
		Table:     tbl,
		Columns:   []string{generatedInvisiblePrimaryKeyColumnName},
		IndexName: "PRIMARY",
	})
}

func rebuildColByNamePreservePositions(tbl *Table) {
	tbl.colByName = make(map[string]int, len(tbl.Columns))
	for i, col := range tbl.Columns {
		tbl.colByName[toLower(col.Name)] = i
	}
}

// analyzeTableExpressions performs best-effort semantic analysis on DEFAULT,
// GENERATED, and CHECK expressions after all columns have been added to the table.
func (c *Catalog) analyzeTableExpressions(tbl *Table, stmt *nodes.CreateTableStmt) {
	// Analyze DEFAULT and GENERATED expressions from column definitions.
	for _, colDef := range stmt.Columns {
		col := tbl.GetColumn(colDef.Name)
		if col == nil {
			continue
		}

		// Top-level DEFAULT.
		if colDef.DefaultValue != nil {
			if analyzed, err := c.AnalyzeStandaloneExpr(colDef.DefaultValue, tbl); err == nil {
				col.DefaultAnalyzed = analyzed
			}
		}

		// Column-constraint DEFAULT (may override top-level).
		for _, cc := range colDef.Constraints {
			if cc.Type == nodes.ColConstrDefault && cc.Expr != nil {
				if analyzed, err := c.AnalyzeStandaloneExpr(cc.Expr, tbl); err == nil {
					col.DefaultAnalyzed = analyzed
				}
			}
		}

		// GENERATED ALWAYS AS.
		if colDef.Generated != nil {
			if analyzed, err := c.AnalyzeStandaloneExpr(colDef.Generated.Expr, tbl); err == nil {
				col.GeneratedAnalyzed = analyzed
			}
		}
	}

	// Analyze CHECK expressions on constraints.
	// We iterate constraints and match CHECK ones; the AST sources are both
	// column-level and table-level constraint nodes.
	checkIdx := 0
	for _, colDef := range stmt.Columns {
		for _, cc := range colDef.Constraints {
			if cc.Type == nodes.ColConstrCheck && cc.Expr != nil {
				// Find matching CHECK constraint by index.
				for checkIdx < len(tbl.Constraints) {
					if tbl.Constraints[checkIdx].Type == ConCheck {
						if analyzed, err := c.AnalyzeStandaloneExpr(cc.Expr, tbl); err == nil {
							tbl.Constraints[checkIdx].CheckAnalyzed = analyzed
						}
						checkIdx++
						break
					}
					checkIdx++
				}
			}
		}
	}
	for _, con := range stmt.Constraints {
		if con.Type == nodes.ConstrCheck && con.Expr != nil {
			for checkIdx < len(tbl.Constraints) {
				if tbl.Constraints[checkIdx].Type == ConCheck {
					if analyzed, err := c.AnalyzeStandaloneExpr(con.Expr, tbl); err == nil {
						tbl.Constraints[checkIdx].CheckAnalyzed = analyzed
					}
					checkIdx++
					break
				}
				checkIdx++
			}
		}
	}
}

// buildPartitionInfo converts an AST PartitionClause to a catalog PartitionInfo.
func buildPartitionInfo(tbl *Table, pc *nodes.PartitionClause) *PartitionInfo {
	pi := &PartitionInfo{
		Linear:   pc.Linear,
		NumParts: pc.NumParts,
	}
	if pi.NumParts == 0 && len(pc.Partitions) == 0 && (pc.Type == nodes.PartitionHash || pc.Type == nodes.PartitionKey) {
		pi.NumParts = 1
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
		if len(pi.Columns) == 0 {
			pi.Columns = primaryKeyColumns(tbl)
		}
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
		if pi.NumSubParts == 0 {
			pi.NumSubParts = 1
		}
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

	// Auto-generate partition definitions for HASH/KEY/LINEAR HASH/LINEAR KEY
	// when PARTITIONS N is specified without explicit partition definitions.
	// MySQL naming convention: p0, p1, p2, ...
	if len(pi.Partitions) == 0 && pi.NumParts > 0 {
		for i := 0; i < pi.NumParts; i++ {
			pi.Partitions = append(pi.Partitions, &PartitionDefInfo{
				Name: fmt.Sprintf("p%d", i),
			})
		}
	}

	// Auto-generate subpartition definitions when SUBPARTITIONS N is specified
	// without explicit subpartition definitions.
	// MySQL naming convention: <partition_name>sp0, <partition_name>sp1, ...
	if pi.NumSubParts > 0 {
		for _, part := range pi.Partitions {
			if len(part.SubPartitions) == 0 {
				for j := 0; j < pi.NumSubParts; j++ {
					part.SubPartitions = append(part.SubPartitions, &SubPartitionDefInfo{
						Name: fmt.Sprintf("%ssp%d", part.Name, j),
					})
				}
			}
		}
	}

	return pi
}

func validatePartitionClause(tbl *Table, pc *nodes.PartitionClause) error {
	if (pc.Type == nodes.PartitionRange || pc.Type == nodes.PartitionList) && len(pc.Partitions) == 0 {
		return &Error{Code: 1492, SQLState: "HY000", Message: "Partitions must be defined for RANGE/LIST partitioning"}
	}
	if pc.Type == nodes.PartitionRange {
		for i, pd := range pc.Partitions {
			if pd.Values != nil && strings.Contains(strings.ToUpper(partitionValueToString(pd.Values, pc.Type)), "MAXVALUE") && i != len(pc.Partitions)-1 {
				return &Error{Code: 1481, SQLState: "HY000", Message: "MAXVALUE can only be used in last partition definition"}
			}
		}
	}
	expr := strings.ToLower(nodeToSQL(pc.Expr))
	if strings.Contains(expr, "concat(") {
		return &Error{Code: 1491, SQLState: "HY000", Message: "The PARTITION function returns the wrong type"}
	}
	if pc.Type == nodes.PartitionRange {
		if cr, ok := pc.Expr.(*nodes.ColumnRef); ok {
			if col := tbl.GetColumn(cr.Column); col != nil && strings.EqualFold(col.DataType, "timestamp") {
				return &Error{Code: 1491, SQLState: "HY000", Message: "A PRIMARY KEY must include all columns in the table's partitioning function"}
			}
		}
	}
	if pc.Type == nodes.PartitionHash && pc.Expr != nil {
		if cr, ok := pc.Expr.(*nodes.ColumnRef); ok {
			if err := validateUniqueKeysCoverPartitionColumns(tbl, []string{cr.Column}); err != nil {
				return err
			}
		}
	}
	return nil
}

func primaryKeyColumns(tbl *Table) []string {
	for _, con := range tbl.Constraints {
		if con.Type == ConPrimaryKey {
			return append([]string{}, con.Columns...)
		}
	}
	for _, idx := range tbl.Indexes {
		if idx.Primary {
			cols := make([]string, 0, len(idx.Columns))
			for _, ic := range idx.Columns {
				cols = append(cols, ic.Name)
			}
			return cols
		}
	}
	return nil
}

func validateUniqueKeysCoverPartitionColumns(tbl *Table, partCols []string) error {
	if len(partCols) == 0 {
		return nil
	}
	for _, idx := range tbl.Indexes {
		if !idx.Unique && !idx.Primary {
			continue
		}
		idxCols := map[string]bool{}
		for _, ic := range idx.Columns {
			idxCols[toLower(ic.Name)] = true
		}
		for _, pc := range partCols {
			if !idxCols[toLower(pc)] {
				return &Error{Code: 1503, SQLState: "HY000", Message: "A UNIQUE INDEX must include all columns in the table's partitioning function"}
			}
		}
	}
	return nil
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
			if subList, ok := item.(*nodes.List); ok {
				// Tuple: (val1, val2) for multi-column LIST COLUMNS
				subParts := make([]string, len(subList.Items))
				for j, sub := range subList.Items {
					subParts[j] = nodeToSQL(sub.(nodes.ExprNode))
				}
				parts[i] = "(" + strings.Join(subParts, ",") + ")"
			} else {
				parts[i] = nodeToSQL(item.(nodes.ExprNode))
			}
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
		if (strings.EqualFold(con.OnDelete, "SET NULL") || strings.EqualFold(con.OnUpdate, "SET NULL")) && !col.Nullable {
			return &Error{Code: 1830, SQLState: "HY000", Message: "Column cannot be NOT NULL: needed in a foreign key constraint"}
		}
		if col.Generated != nil && !col.Generated.Stored {
			return errFKCannotUseVirtualColumn(colName)
		}
		if refCol.Generated != nil && !refCol.Generated.Stored {
			return errFKCannotUseVirtualColumn(con.RefColumns[i])
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
			if cr, ok := ic.Expr.(*nodes.ColumnRef); ok && !ic.Functional {
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
			if cr, ok := ic.Expr.(*nodes.ColumnRef); ok && !ic.Functional {
				idxCol.Name = cr.Column
			} else {
				idxCol.Expr = nodeToSQL(ic.Expr)
				idxCol.ExprNode = ic.Expr
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
	if strings.EqualFold(candidate, "PRIMARY") {
		candidate = fmt.Sprintf("%s_%d", baseName, suffix)
		suffix++
	}
	for indexNameExists(tbl, candidate) {
		candidate = fmt.Sprintf("%s_%d", baseName, suffix)
		suffix++
	}
	return candidate
}

func validateNonPrimaryIndexName(name string) error {
	if strings.EqualFold(name, "PRIMARY") {
		return errWrongNameForIndex(name)
	}
	return nil
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

func validateIndexColumns(tbl *Table, idxCols []*IndexColumn, fulltext, spatial bool) error {
	for _, ic := range idxCols {
		if ic.Name == "" {
			continue
		}
		col := tbl.GetColumn(ic.Name)
		if col == nil {
			continue
		}
		if spatial {
			if col.Nullable {
				return &Error{Code: 1252, SQLState: "42000", Message: "All parts of a SPATIAL index must be NOT NULL"}
			}
			continue
		}
		if !fulltext && isTextBlobType(col.DataType) && ic.Length == 0 {
			return &Error{Code: 1170, SQLState: "42000", Message: "BLOB/TEXT column used in key specification without a key length"}
		}
		if ic.Length > 0 {
			if maxLen, ok := fixedLengthStringColumnLimit(col); ok && ic.Length > maxLen {
				return &Error{Code: 1089, SQLState: "HY000", Message: "Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys"}
			}
		}
	}
	return nil
}

func fixedLengthStringColumnLimit(col *Column) (int, bool) {
	switch strings.ToLower(col.DataType) {
	case "char", "varchar", "binary", "varbinary":
	default:
		return 0, false
	}
	open := strings.IndexByte(col.ColumnType, '(')
	close := strings.IndexByte(col.ColumnType, ')')
	if open < 0 || close <= open+1 {
		return 0, false
	}
	n := 0
	for _, ch := range col.ColumnType[open+1 : close] {
		if ch < '0' || ch > '9' {
			return 0, false
		}
		n = n*10 + int(ch-'0')
	}
	return n, n > 0
}

func coerceInnoDBHashIndex(tbl *Table, idx *Index) {
	if strings.EqualFold(tbl.Engine, "InnoDB") && strings.EqualFold(idx.IndexType, "HASH") {
		idx.IndexType = "BTREE"
	}
}

func constraintNameExistsInTable(tbl *Table, typ ConstraintType, name string) bool {
	key := toLower(name)
	for _, con := range tbl.Constraints {
		if con.Type == typ && toLower(con.Name) == key {
			return true
		}
	}
	return false
}

func checkConstraintNameExists(db *Database, pending *Table, name string) bool {
	if pending != nil && constraintNameExistsInTable(pending, ConCheck, name) {
		return true
	}
	key := toLower(name)
	for _, tbl := range db.Tables {
		for _, con := range tbl.Constraints {
			if con.Type == ConCheck && toLower(con.Name) == key {
				return true
			}
		}
	}
	return false
}

func foreignKeyConstraintNameExists(db *Database, pending *Table, name string) bool {
	if pending != nil && constraintNameExistsInTable(pending, ConForeignKey, name) {
		return true
	}
	key := toLower(name)
	for _, tbl := range db.Tables {
		for _, con := range tbl.Constraints {
			if con.Type == ConForeignKey && toLower(con.Name) == key {
				return true
			}
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
			switch n := opt.Value.(type) {
			case *nodes.Integer:
				idx.KeyBlockSize = int(n.Ival)
			case *nodes.IntLit:
				idx.KeyBlockSize = int(n.Value)
			}
		}
	}
}

// nextFKGeneratedNumber returns the next available counter for an auto-generated
// InnoDB FK constraint name of the form "<tableName>_ibfk_<N>".
//
// This matches MySQL 8.0's behavior in sql/sql_table.cc:5843
// (get_fk_max_generated_name_number): it scans existing FK constraints on the
// table, parses any name that looks like "<tableName>_ibfk_<digits>" as a
// generated name, and returns max(N)+1 (or 1 if no such names exist).
//
// Bytebase omni catalog is case-insensitive on table names (we lowercase the
// prefix before comparison). MySQL's own comparison is case-sensitive on
// already-lowered names, so this is equivalent for typical use.
//
// As in MySQL, pre-4.0.18-style names ("<tableName>_ibfk_0<digits>") are
// ignored — we skip anything whose counter substring starts with '0'.
func nextFKGeneratedNumber(tbl *Table, tableName string) int {
	prefix := toLower(tableName) + "_ibfk_"
	max := 0
	for _, con := range tbl.Constraints {
		if con.Type != ConForeignKey {
			continue
		}
		name := toLower(con.Name)
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		rest := name[len(prefix):]
		if rest == "" || rest[0] == '0' {
			continue
		}
		n := 0
		ok := true
		for _, ch := range rest {
			if ch < '0' || ch > '9' {
				ok = false
				break
			}
			n = n*10 + int(ch-'0')
		}
		if !ok {
			continue
		}
		if n > max {
			max = n
		}
	}
	return max + 1
}

// nextCheckNumber returns max(existing generated tableName_chk_N) + 1.
// ALTER TABLE seeds the generated CHECK counter from existing names.
func nextCheckNumber(tbl *Table) int {
	prefix := toLower(tbl.Name) + "_chk_"
	max := 0
	for _, con := range tbl.Constraints {
		if con.Type != ConCheck {
			continue
		}
		name := toLower(con.Name)
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		rest := name[len(prefix):]
		if rest == "" || rest[0] == '0' {
			continue
		}
		n := 0
		ok := true
		for _, ch := range rest {
			if ch < '0' || ch > '9' {
				ok = false
				break
			}
			n = n*10 + int(ch-'0')
		}
		if ok && n > max {
			max = n
		}
	}
	return max + 1
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

func applyBinaryModifierCollation(col *Column, dt *nodes.DataType) {
	if col == nil || dt == nil || !dt.Binary || col.Charset == "" {
		return
	}
	col.Collation = fmt.Sprintf("%s_bin", normalizeCharsetName(col.Charset))
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
	return deparse.Deparse(node)
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
	case nodes.BinOpXor:
		return "XOR"
	case nodes.BinOpRegexp:
		return "REGEXP"
	case nodes.BinOpLikeEscape:
		return "LIKE"
	case nodes.BinOpNullSafeEq:
		return "<=>"
	case nodes.BinOpJsonExtract:
		return "->"
	case nodes.BinOpJsonUnquote:
		return "->>"
	case nodes.BinOpSoundsLike:
		return "SOUNDS LIKE"
	default:
		return "?"
	}
}

func formatColumnType(dt *nodes.DataType) string {
	name := normalizedTypeName(dt.Name)
	if name == "float" && dt.Length > 24 && dt.Scale == 0 {
		name = "double"
	}
	if name == "varchar" && dt.Length > 65535 {
		return "mediumtext"
	}
	if name == "text" && dt.Length > 0 {
		name = promoteTextType(dt.Length)
	}

	// MySQL type aliases: BOOLEAN/BOOL → tinyint(1), NUMERIC → decimal, SERIAL → bigint unsigned
	// GEOMETRYCOLLECTION → geomcollection (MySQL 8.0 normalized form)
	switch name {
	case "boolean", "bool":
		return "tinyint(1)"
	case "serial":
		return "bigint unsigned"
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
	} else if name == "decimal" {
		if dt.Length == 0 && dt.Scale == 0 {
			// DECIMAL with no precision → MySQL shows decimal(10,0)
			buf.WriteString("(10,0)")
		} else {
			fmt.Fprintf(&buf, "(%d,%d)", dt.Length, dt.Scale)
		}
	} else if isTextBlobLengthStripped(name) {
		// TEXT(n) and BLOB(n) — MySQL stores the length internally but
		// SHOW CREATE TABLE displays just TEXT / BLOB without the length.
		// Do not emit length.
	} else if name == "year" {
		// YEAR(4) is deprecated in MySQL 8.0 — SHOW CREATE TABLE shows just `year`.
	} else if name == "bit" && dt.Length == 0 {
		buf.WriteString("(1)")
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
	if dt.Unsigned || dt.Zerofill {
		buf.WriteString(" unsigned")
	}
	if dt.Zerofill {
		buf.WriteString(" zerofill")
	}
	return buf.String()
}

func normalizedColumnDataType(dt *nodes.DataType) string {
	name := normalizedTypeName(dt.Name)
	if name == "float" && dt.Length > 24 && dt.Scale == 0 {
		return "double"
	}
	if name == "varchar" && dt.Length > 65535 {
		return "mediumtext"
	}
	if name == "text" && dt.Length > 0 {
		return promoteTextType(dt.Length)
	}
	return name
}

func normalizedTypeName(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "bool", "boolean":
		return "boolean"
	case "integer", "int4":
		return "int"
	case "int1":
		return "tinyint"
	case "int2":
		return "smallint"
	case "int3", "middleint":
		return "mediumint"
	case "int8":
		return "bigint"
	case "real", "float8", "double precision":
		return "double"
	case "float4":
		return "float"
	case "numeric", "dec", "fixed":
		return "decimal"
	case "character":
		return "char"
	case "character varying", "varcharacter", "nvarchar", "national varchar",
		"nchar varchar", "national char varying", "nchar varying":
		return "varchar"
	case "national char", "nchar":
		return "char"
	case "geometrycollection":
		return "geomcollection"
	default:
		return strings.ToLower(strings.TrimSpace(name))
	}
}

func isNationalStringType(name string) bool {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "nvarchar", "national varchar", "nchar varchar", "national char varying",
		"nchar varying", "national char", "nchar":
		return true
	default:
		return false
	}
}

func normalizeCharsetName(name string) string {
	if strings.EqualFold(name, "utf8") {
		return "utf8mb3"
	}
	return strings.ToLower(strings.TrimSpace(name))
}

func charsetMatchesCollation(charset, collation string) bool {
	derived := normalizeCharsetName(charsetForCollation(collation))
	return derived == "" || strings.EqualFold(normalizeCharsetName(charset), derived)
}

func promoteTextType(chars int) string {
	bytes := chars * 4
	switch {
	case bytes <= 255:
		return "tinytext"
	case bytes <= 65535:
		return "text"
	case bytes <= 16777215:
		return "mediumtext"
	default:
		return "longtext"
	}
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
// MySQL uses ” (two single quotes) to escape a single quote in enum values.
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
		Name:             tableName,
		Database:         db,
		Columns:          make([]*Column, 0, len(srcTbl.Columns)),
		colByName:        make(map[string]int),
		Indexes:          make([]*Index, 0, len(srcTbl.Indexes)),
		Constraints:      make([]*Constraint, 0, len(srcTbl.Constraints)),
		Engine:           srcTbl.Engine,
		Charset:          srcTbl.Charset,
		Collation:        srcTbl.Collation,
		Comment:          srcTbl.Comment,
		RowFormat:        srcTbl.RowFormat,
		KeyBlockSize:     srcTbl.KeyBlockSize,
		Compression:      srcTbl.Compression,
		Encryption:       srcTbl.Encryption,
		StatsPersistent:  srcTbl.StatsPersistent,
		StatsAutoRecalc:  srcTbl.StatsAutoRecalc,
		StatsSamplePages: srcTbl.StatsSamplePages,
		MinRows:          srcTbl.MinRows,
		MaxRows:          srcTbl.MaxRows,
		AvgRowLength:     srcTbl.AvgRowLength,
		Tablespace:       srcTbl.Tablespace,
		PackKeys:         srcTbl.PackKeys,
		Checksum:         srcTbl.Checksum,
		DelayKeyWrite:    srcTbl.DelayKeyWrite,
		Temporary:        stmt.Temporary,
	}

	// Copy columns.
	for i, srcCol := range srcTbl.Columns {
		col := &Column{
			Position:                     srcCol.Position,
			Name:                         srcCol.Name,
			DataType:                     srcCol.DataType,
			ColumnType:                   srcCol.ColumnType,
			Nullable:                     srcCol.Nullable,
			AutoIncrement:                srcCol.AutoIncrement,
			Charset:                      srcCol.Charset,
			Collation:                    srcCol.Collation,
			Comment:                      srcCol.Comment,
			OnUpdate:                     srcCol.OnUpdate,
			Invisible:                    srcCol.Invisible,
			GeneratedInvisiblePrimaryKey: srcCol.GeneratedInvisiblePrimaryKey,
			Hidden:                       srcCol.Hidden,
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
				ExprNode:   sc.ExprNode,
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
