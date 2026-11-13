package catalog

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/deparse"
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
	if err := checkSomeColumnVisible(stmt); err != nil {
		return err
	}

	tableName := stmt.Table.Name
	key := toLower(tableName)

	// The name is checked last: OceanBase reads the whole definition first, so
	// CREATE TABLE t (x INT, x INT) is 1060 when t exists, even with IF NOT EXISTS
	// (TestOracle_OceanBaseNameOrder). nameTaken answers for it; done is set when
	// IF NOT EXISTS makes the statement do nothing.
	nameTaken := func() (done bool, err error) {
		if db.Tables[key] != nil {
			if stmt.IfNotExists {
				return true, nil
			}
			return false, errDupTable(tableName)
		}
		if db.Views[key] != nil {
			return false, errDupTable(tableName)
		}
		return false, nil
	}

	// CREATE TABLE ... LIKE
	if stmt.Like != nil {
		return c.createTableLike(db, tableName, key, stmt, nameTaken)
	}

	// CREATE TABLE ... AS SELECT (CTAS) — not supported yet, skip silently
	if stmt.Select != nil && len(stmt.Columns) == 0 {
		_, err := nameTaken()
		return err
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
			tbl.Collation = normalizeCollationName(opt.Value)
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
			if err := validateCompressor(opt.Value); err != nil {
				return err
			}
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
	autoIncrementColumn := ""

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
				col.ColumnType = "bigint(20) unsigned"
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
				// OceanBase's national character set is utf8mb4, where MySQL's is utf8mb3.
				// A national type names its character set as CHARACTER SET does, so its
				// collation is that character set's and not its table's
				// (TestOracle_OceanBaseNationalCollation).
				col.Charset = "utf8mb4"
				col.CharsetExplicit = true
			}
			if colDef.TypeName.Charset != "" {
				col.Charset = normalizeCharsetName(colDef.TypeName.Charset)
				col.CharsetExplicit = true
			}
			if colDef.TypeName.Collate != "" {
				col.Collation = normalizeCollationName(colDef.TypeName.Collate)
				col.CollationExplicit = true
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
				// A column that names its character set, or has another than its
				// table's, takes that character set's default collation, not the
				// table's collation.
				if col.CharsetExplicit || !strings.EqualFold(col.Charset, tbl.Charset) {
					if dc, ok := defaultCollationForCharset[toLower(col.Charset)]; ok {
						col.Collation = dc
					}
				} else {
					col.Collation = tbl.Collation
				}
			}
			applyBinaryModifierCollation(col, colDef.TypeName)
		}

		// Top-level column properties. Copy SRID by PRESENCE (HasSRID), not by a nonzero
		// value: MySQL keeps `SRID 0` as a distinct, present clause, so a `SRID 0` column
		// (SRID==0, HasSRID==true) must not be mistaken for a no-SRID column.
		if colDef.TypeName != nil && colDef.TypeName.HasSRID {
			col.SRID = colDef.TypeName.SRID
			col.HasSRID = true
		}
		if colDef.AutoIncrement {
			col.AutoIncrement = true
			col.Nullable = false
		}
		if colDef.Comment != "" {
			col.Comment = colDef.Comment
		}
		if colDef.DefaultValue != nil {
			setColumnDefaultFromExpr(col, colDef.DefaultValue)
		}
		if colDef.OnUpdate != nil {
			setColumnOnUpdateFromExpr(col, colDef.OnUpdate)
		}
		if colDef.Generated != nil {
			// Generated expressions are stored by the engine in resolver-rewritten
			// form — TIMESTAMPADD/DATE_ADD fold to interval arithmetic, NOT folds
			// to inverted comparisons — identical to view bodies (oracle 8.0.32:
			// AS (TIMESTAMPADD(HOUR,1,a)) stores as ((`a` + interval 1 hour)),
			// AS (NOT a) as ((0 = `a`))). Apply the same rewrite pass before
			// rendering so user-declared SDL canonicalizes to the stored form.
			col.Generated = &GeneratedColumnInfo{
				Expr:   nodeToSQLGenerated(deparse.RewriteExpr(colDef.Generated.Expr), tbl.Charset),
				Stored: colDef.Generated.Stored,
			}
		}
		if col.AutoIncrement {
			if autoIncrementColumn != "" {
				return errIncorrectAutoIncrementDefinition()
			}
			autoIncrementColumn = colDef.Name
		}

		// Process column-level constraints.
		for _, cc := range colDef.Constraints {
			switch cc.Type {
			case nodes.ColConstrNotNull:
				col.Nullable = false
				col.NullExplicit = true
			case nodes.ColConstrNull:
				col.Nullable = true
				col.NullExplicit = true
			case nodes.ColConstrDefault:
				if cc.Expr != nil {
					setColumnDefaultFromExpr(col, cc.Expr)
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
					conName = c.madeUpConstraintName(tableName, madeUpCheckInfix, func(name string) bool {
						return checkConstraintNameExists(db, tbl, name)
					})
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
						col.Collation = normalizeCollationName(s.Value)
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

	c.applyTimestampSessionDefaults(tbl.Columns, stmt.Columns)

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
		if err := validateStoringColumns(tbl, con.IndexOptions); err != nil {
			return err
		}
		cols := extractColumnNames(con)

		switch con.Type {
		case nodes.ConstrPrimaryKey:
			if hasPK {
				return errDuplicatedPriKey()
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
				conName = c.madeUpConstraintName(tableName, madeUpForeignKeyInfix, func(name string) bool {
					return foreignKeyConstraintNameExists(db, tbl, name)
				})
			}
			// A name twice in one table is 1215 and a name another table's key has is
			// 1022, where MySQL answers ER_FK_DUP_NAME (1826) for both.
			if constraintNameExistsInTable(tbl, ConForeignKey, conName) {
				return errCannotAddForeign()
			}
			if foreignKeyConstraintNameExists(db, nil, conName) {
				return errFKDupName(tableName)
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
			// No index is made for the foreign key: OceanBase makes none, where MySQL
			// adds one when no index serves it. TestOracle_OceanBaseForeignKeyIndex
			// holds this to the engine.

		case nodes.ConstrCheck:
			conName := con.Name
			if conName == "" {
				conName = c.madeUpConstraintName(tableName, madeUpCheckInfix, func(name string) bool {
					return checkConstraintNameExists(db, tbl, name)
				})
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

		case nodes.ConstrVector:
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
			vecIdx := &Index{
				Name:      idxName,
				Table:     tbl,
				Columns:   buildIndexColumns(con),
				Vector:    true,
				IndexType: resolveConstraintIndexType(con),
				Visible:   true,
			}
			applyIndexOptions(vecIdx, con.IndexOptions)
			tbl.Indexes = append(tbl.Indexes, vecIdx)

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
		if err := validatePartitionClause(tbl, stmt.Partitions); err != nil {
			return err
		}
		tbl.Partitioning = buildPartitionInfo(tbl, stmt.Partitions)
	}

	// Phase 3: analyze DEFAULT, GENERATED, and CHECK expressions now that all
	// columns are present in the table.
	c.analyzeTableExpressions(tbl, stmt)

	if done, err := nameTaken(); done || err != nil {
		return err
	}
	db.Tables[key] = tbl
	return nil
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
		Linear:     pc.Linear,
		NumParts:   pc.NumParts,
		keyColumns: partitionKeyColumns(tbl, pc),
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
		pi.UseDefaultSubpartitions = true
		pi.UseDefaultNumSubpartitions = true
		switch pc.SubPartType {
		case nodes.PartitionHash:
			pi.SubType = "HASH"
			pi.SubExpr = nodeToSQL(pc.SubPartExpr)
		case nodes.PartitionKey:
			pi.SubType = "KEY"
			pi.SubColumns = pc.SubPartColumns
			pi.SubAlgo = pc.SubPartAlgo
		// OceanBase also subpartitions by RANGE and LIST, with or without COLUMNS.
		case nodes.PartitionRange, nodes.PartitionList:
			pi.SubType = "RANGE"
			if pc.SubPartType == nodes.PartitionList {
				pi.SubType = "LIST"
			}
			if len(pc.SubPartColumns) > 0 {
				pi.SubType += " COLUMNS"
				pi.SubColumns = pc.SubPartColumns
			} else {
				pi.SubExpr = nodeToSQL(pc.SubPartExpr)
			}
		}
		pi.SubLinear = false // TODO: track linear for subpartitions if parser supports it
		if pc.NumSubParts > 0 {
			pi.NumSubParts = pc.NumSubParts
			pi.UseDefaultNumSubpartitions = false
		} else {
			pi.NumSubParts = 1
		}
	}

	// Partition definitions.
	explicitSubpartitions := false
	explicitSubpartitionCount := 0
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
		if len(pd.SubPartitions) > 0 {
			if !explicitSubpartitions {
				explicitSubpartitions = true
				explicitSubpartitionCount = len(pd.SubPartitions)
			}
		}
		for _, spd := range pd.SubPartitions {
			spdi := &SubPartitionDefInfo{
				Name:    spd.Name,
				Engine:  pdi.Engine,
				Comment: pdi.Comment,
			}
			if spd.Values != nil {
				spdi.ValueExpr = partitionValueToString(spd.Values, pc.SubPartType)
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
	if explicitSubpartitions {
		pi.UseDefaultSubpartitions = false
		pi.UseDefaultNumSubpartitions = false
		// OceanBase lets a partition list no subpartitions beside ones that do, and
		// gives it SUBPARTITIONS N of them, or one when N is unstated, named
		// <partition>sp<i> (4.3.5.6). The table then has no template, so the count
		// stops at CREATE: a partition added later gets one (addedPartitionSubpartitions).
		fill := pc.NumSubParts
		if fill == 0 {
			fill = 1
		}
		generated := 0
		for _, pdi := range pi.Partitions {
			if len(pdi.SubPartitions) > 0 {
				continue
			}
			for j := 0; j < fill; j++ {
				name := fmt.Sprintf("%ssp%d", pdi.Name, j)
				if defaultSubpartitionValue(pi.SubType) != "" {
					name = generatedSubpartitionName(&generated)
				}
				pdi.SubPartitions = append(pdi.SubPartitions, &SubPartitionDefInfo{
					Name:      name,
					ValueExpr: defaultSubpartitionValue(pi.SubType),
					Engine:    pdi.Engine,
					Comment:   pdi.Comment,
				})
			}
		}
		pi.NumSubParts = explicitSubpartitionCount
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

	// OceanBase's SUBPARTITION TEMPLATE gives every partition the same subpartitions,
	// where MySQL repeats them inside each one. This sits after the outer partitions are
	// built, generated ones included, because HASH and KEY with PARTITIONS N define none
	// of their own and a template has to reach those too. Clearing
	// UseDefaultSubpartitions is also what stops the block below inventing a second set.
	// The template is kept, because a partition added later takes it too.
	if len(pc.SubPartTemplate) > 0 {
		for _, spd := range pc.SubPartTemplate {
			st := &SubPartitionDefInfo{Name: spd.Name}
			if spd.Values != nil {
				st.ValueExpr = partitionValueToString(spd.Values, pc.SubPartType)
			}
			for _, opt := range spd.Options {
				switch toLower(opt.Name) {
				case "engine":
					st.Engine = opt.Value
				case "comment":
					st.Comment = opt.Value
				}
			}
			pi.SubPartTemplate = append(pi.SubPartTemplate, st)
		}
		for _, pdi := range pi.Partitions {
			if len(pdi.SubPartitions) == 0 {
				pdi.SubPartitions = templateSubpartitions(pi.SubPartTemplate, pdi)
			}
		}
		pi.NumSubParts = len(pc.SubPartTemplate)
		pi.UseDefaultSubpartitions = false
		pi.UseDefaultNumSubpartitions = false
	}

	// RANGE and LIST subpartitioning with no subpartitions listed at all. OceanBase
	// gives each partition one, to MAXVALUE or DEFAULT, named as generatedSubpartitionName
	// names them, prints them per partition, and keeps no template: a partition added
	// later gets <partition>sp0 (rangeListSubpartitionCases).
	if defaultSubpartitionValue(pi.SubType) != "" && pi.UseDefaultSubpartitions {
		generated := 0
		for _, part := range pi.Partitions {
			if len(part.SubPartitions) == 0 {
				part.SubPartitions = append(part.SubPartitions, &SubPartitionDefInfo{
					Name:      generatedSubpartitionName(&generated),
					ValueExpr: defaultSubpartitionValue(pi.SubType),
					Engine:    part.Engine,
					Comment:   part.Comment,
				})
			}
		}
		pi.UseDefaultSubpartitions = false
		pi.UseDefaultNumSubpartitions = false
	}

	// Auto-generate subpartition definitions when SUBPARTITIONS N is specified
	// without explicit subpartition definitions.
	// MySQL naming convention: <partition_name>sp0, <partition_name>sp1, ...
	if pi.SubType != "" && pi.UseDefaultSubpartitions && pi.NumSubParts > 0 {
		for _, part := range pi.Partitions {
			if len(part.SubPartitions) == 0 {
				for j := 0; j < pi.NumSubParts; j++ {
					part.SubPartitions = append(part.SubPartitions, &SubPartitionDefInfo{
						Name:      fmt.Sprintf("%ssp%d", part.Name, j),
						ValueExpr: defaultSubpartitionValue(pi.SubType),
					})
				}
			}
		}
		// OceanBase keeps SUBPARTITIONS N as a template of p0 ... p<N-1>: SHOW CREATE
		// TABLE prints `subpartition template (subpartition p0, subpartition p1)`, and a
		// partition added later takes it, so templateSubpartitions names its
		// subpartitions <partition>sp<i> as above.
		for j := 0; j < pi.NumSubParts; j++ {
			pi.SubPartTemplate = append(pi.SubPartTemplate, &SubPartitionDefInfo{
				Name:      fmt.Sprintf("p%d", j),
				ValueExpr: defaultSubpartitionValue(pi.SubType),
			})
		}
	}

	return pi
}

// defaultSubpartitionValue is the bound OceanBase gives a subpartition it makes up for a
// RANGE or LIST subpartitioned table, one the DDL did not list: MAXVALUE for RANGE and
// DEFAULT for LIST, so it holds every row the listed ones do not. A partition added to
// a RANGE subpartitioned table without subpartitions reads back as
// (subpartition `p1sp0` values less than (MAXVALUE)) on 4.3.5.6.
func defaultSubpartitionValue(subType string) string {
	switch {
	case strings.HasPrefix(subType, "RANGE"):
		return "MAXVALUE"
	case strings.HasPrefix(subType, "LIST"):
		return "DEFAULT"
	}
	return ""
}

// generatedSubpartitionName names the next subpartition OceanBase makes up at CREATE
// TABLE for a RANGE or LIST subpartitioned partition that listed none: P8192, P8193,
// ..., counted across the table in partition order (4.3.5.6 reads back
// PARTITION BY HASH(a) PARTITIONS 2 SUBPARTITION BY LIST(b) as p0 (P8192) and
// p1 (P8193)). count is the number named so far.
func generatedSubpartitionName(count *int) string {
	name := fmt.Sprintf("P%d", 8192+*count)
	*count++
	return name
}

// templateSubpartitions gives partition pd the subpartitions of a SUBPARTITION
// TEMPLATE. OceanBase names each one after the partition and the template entry,
// partition name, "s", entry name: a template of s0 and s1 gives partition p0 the
// subpartitions p0ss0 and p0ss1 (information_schema.PARTITIONS on 4.3.5.6), which keeps
// the names unique across the table, as a subpartition name has to be.
func templateSubpartitions(template []*SubPartitionDefInfo, pd *PartitionDefInfo) []*SubPartitionDefInfo {
	out := make([]*SubPartitionDefInfo, 0, len(template))
	for _, st := range template {
		sp := &SubPartitionDefInfo{
			Name:      pd.Name + "s" + st.Name,
			ValueExpr: st.ValueExpr,
			Engine:    pd.Engine,
			Comment:   pd.Comment,
		}
		if st.Engine != "" {
			sp.Engine = st.Engine
		}
		if st.Comment != "" {
			sp.Comment = st.Comment
		}
		out = append(out, sp)
	}
	return out
}

func (c *Catalog) applyTimestampSessionDefaults(cols []*Column, defs []*nodes.ColumnDef) {
	if c.session.ExplicitDefaultsForTimestamp {
		return
	}
	promotedFirst := false
	for i, col := range cols {
		if col == nil || !strings.EqualFold(col.DataType, "timestamp") || i >= len(defs) {
			continue
		}
		def := defs[i]
		if hasExplicitNull(def) {
			continue
		}
		col.Nullable = false

		hasDefault := hasExplicitDefault(def)
		hasOnUpdate := hasExplicitOnUpdate(def)
		if !promotedFirst && !hasDefault && !hasOnUpdate {
			current := timestampCurrentExpr(timestampFsp(def))
			col.Default = &current
			col.DefaultKind = ColumnDefaultCurrentTimestamp
			col.OnUpdate = current
			col.OnUpdateKind = ColumnDefaultCurrentTimestamp
			promotedFirst = true
			continue
		}
		if !hasDefault && col.Default == nil {
			zero := timestampZeroDefault(timestampFsp(def))
			col.Default = &zero
			col.DefaultKind = ColumnDefaultConstant
		}
	}
}

func hasExplicitNull(def *nodes.ColumnDef) bool {
	if def == nil {
		return false
	}
	for _, cc := range def.Constraints {
		if cc.Type == nodes.ColConstrNull {
			return true
		}
	}
	return false
}

func hasExplicitDefault(def *nodes.ColumnDef) bool {
	if def == nil {
		return false
	}
	if def.DefaultValue != nil {
		return true
	}
	for _, cc := range def.Constraints {
		if cc.Type == nodes.ColConstrDefault {
			return true
		}
	}
	return false
}

func hasExplicitOnUpdate(def *nodes.ColumnDef) bool {
	return def != nil && def.OnUpdate != nil
}

func timestampFsp(def *nodes.ColumnDef) int {
	if def == nil || def.TypeName == nil {
		return 0
	}
	return def.TypeName.Length
}

func timestampCurrentExpr(fsp int) string {
	if fsp > 0 {
		return fmt.Sprintf("CURRENT_TIMESTAMP(%d)", fsp)
	}
	return "CURRENT_TIMESTAMP"
}

func timestampZeroDefault(fsp int) string {
	if fsp > 0 {
		return "0000-00-00 00:00:00." + strings.Repeat("0", fsp)
	}
	return "0000-00-00 00:00:00"
}

func setColumnDefaultFromExpr(col *Column, expr nodes.ExprNode) {
	val, kind := columnDefaultValue(expr)
	col.Default = &val
	col.DefaultKind = kind
	col.DefaultDropped = false
}

func setColumnOnUpdateFromExpr(col *Column, expr nodes.ExprNode) {
	val, kind := columnDefaultValue(expr)
	col.OnUpdate = val
	col.OnUpdateKind = kind
}

func columnDefaultValue(expr nodes.ExprNode) (string, ColumnDefaultKind) {
	switch e := expr.(type) {
	case *nodes.ParenExpr:
		if ts, ok := currentTimestampExpr(e.Expr); ok {
			return ts, ColumnDefaultCurrentTimestamp
		}
		return nodeToSQL(e.Expr), ColumnDefaultExpression
	case *nodes.StringLit:
		// MySQL stores a bare string default by VALUE: a charset introducer is
		// not part of the stored default — SHOW CREATE TABLE echoes
		// DEFAULT _utf8mb4'x' as DEFAULT 'x' (oracle 8.0.32 + 5.7.25; same for
		// the folded _utf8mb4'x' 'y' form). Deparsing the introducer into the
		// stored value would render DEFAULT '_utf8mb4'x'' — invalid SQL that
		// also phantom-diffs against the engine readback. Expression defaults
		// (the ParenExpr case above) legitimately keep the introducer: the
		// engine itself stores DEFAULT (_utf8mb4'xy').
		if e.Charset != "" {
			stripped := *e
			stripped.Charset = ""
			return nodeToSQL(&stripped), ColumnDefaultConstant
		}
		return nodeToSQL(expr), ColumnDefaultConstant
	case *nodes.BoolLit:
		if e.Value {
			return "1", ColumnDefaultConstant
		}
		return "0", ColumnDefaultConstant
	case *nodes.BitLit:
		return canonicalBitLiteral(e.Value), ColumnDefaultConstant
	case *nodes.IntLit:
		// Use the exact source digits when present: deparse renders IntLit from its int64
		// Value, which the lexer clamps to math.MaxInt64 for an unsigned literal beyond
		// int64 range (BIGINT UNSIGNED DEFAULT 18446744073709551615 would otherwise be
		// stored as 9223372036854775807 and silently corrupt the column default).
		if s, ok := exactIntLiteral(e, ""); ok {
			return s, ColumnDefaultConstant
		}
		return nodeToSQL(expr), ColumnDefaultConstant
	case *nodes.UnaryExpr:
		// A signed integer default (DEFAULT -9223372036854775808) is a unary minus over an
		// IntLit. The magnitude 9223372036854775808 overflows int64 on the way in (clamped
		// to MaxInt64, then negated to -MaxInt64), so render the sign + exact digits to keep
		// the int64-MIN boundary exact rather than off-by-one.
		if e.Op == nodes.UnaryMinus || e.Op == nodes.UnaryPlus {
			if lit, ok := e.Operand.(*nodes.IntLit); ok {
				sign := ""
				if e.Op == nodes.UnaryMinus {
					sign = "-"
				}
				if s, ok := exactIntLiteral(lit, sign); ok {
					return s, ColumnDefaultConstant
				}
			}
		}
		return nodeToSQL(expr), ColumnDefaultConstant
	default:
		if ts, ok := currentTimestampExpr(expr); ok {
			return ts, ColumnDefaultCurrentTimestamp
		}
		return nodeToSQL(expr), ColumnDefaultConstant
	}
}

func currentTimestampExpr(expr nodes.ExprNode) (string, bool) {
	fn, ok := expr.(*nodes.FuncCallExpr)
	if !ok {
		return "", false
	}
	name := strings.ToLower(fn.Name)
	if name != "current_timestamp" && name != "now" {
		return "", false
	}
	if len(fn.Args) == 0 {
		return "CURRENT_TIMESTAMP", true
	}
	if len(fn.Args) == 1 {
		if lit, ok := fn.Args[0].(*nodes.IntLit); ok {
			return fmt.Sprintf("CURRENT_TIMESTAMP(%d)", lit.Value), true
		}
	}
	if ts, ok := currentTimestampSQL(nodeToSQL(expr)); ok {
		return ts, true
	}
	return "CURRENT_TIMESTAMP", true
}

// exactIntLiteral renders an integer literal default from its exact source digits
// (lit.Text) with an optional leading sign, so a literal whose magnitude exceeds int64
// range round-trips exactly instead of being clamped to math.MaxInt64. It returns ok=false
// when Text is absent (a synthesized node) so the caller falls back to int64-based deparse.
// A negative all-zero magnitude drops the sign (e.g. "-0"/"-00" -> the digits as-is, never
// a "-0"-prefixed string) so a signed-zero default does not render with a leading minus.
func exactIntLiteral(lit *nodes.IntLit, sign string) (string, bool) {
	if lit.Text == "" {
		return "", false
	}
	if sign == "-" && strings.Trim(lit.Text, "0") == "" {
		return lit.Text, true
	}
	return sign + lit.Text, true
}

func canonicalBitLiteral(bits string) string {
	bits = strings.TrimSpace(bits)
	bits = strings.TrimPrefix(bits, "0b")
	bits = strings.TrimPrefix(bits, "0B")
	if len(bits) >= 3 && (strings.HasPrefix(bits, "b'") || strings.HasPrefix(bits, "B'")) && strings.HasSuffix(bits, "'") {
		bits = bits[2 : len(bits)-1]
	}
	bits = strings.TrimLeft(bits, "0")
	if bits == "" {
		bits = "0"
	}
	return "b'" + bits + "'"
}

// errSubpartitionBoundKind is OceanBase's answer to a subpartition whose bound does not
// fit SUBPARTITION BY: a RANGE subpartition without VALUES LESS THAN, a LIST one
// without VALUES IN, or a HASH or KEY one with either. The grammar takes them, and
// 4.3.5.6 refuses the table with 1210.
func errSubpartitionBoundKind() error {
	return &Error{Code: 1210, SQLState: "HY000", Message: "Invalid argument"}
}

// checkSubpartitionBoundKinds checks the subpartitions a partition lists against the
// table's SUBPARTITION BY type.
func checkSubpartitionBoundKinds(subType nodes.PartitionType, subs []*nodes.SubPartitionDef) error {
	for _, spd := range subs {
		var ok bool
		switch subType {
		case nodes.PartitionRange:
			ok = spd.Values != nil && !spd.ValuesIn
		case nodes.PartitionList:
			ok = spd.Values != nil && spd.ValuesIn
		default:
			ok = spd.Values == nil
		}
		if !ok {
			return errSubpartitionBoundKind()
		}
	}
	return nil
}

// checkSomeColumnVisible refuses a table element list with no visible column, as
// (KEY (b)) or (b INT INVISIBLE): OceanBase answers 4028 before it looks at the name
// (1050), the query of a CREATE TABLE ... SELECT, or anything the list names
// (TestOracle_OceanBaseTableWithoutColumns).
func checkSomeColumnVisible(stmt *nodes.CreateTableStmt) error {
	if len(stmt.Columns) == 0 && len(stmt.Constraints) == 0 {
		return nil // no element list: CREATE TABLE ... LIKE or ... SELECT alone
	}
	for _, col := range stmt.Columns {
		if !slices.ContainsFunc(col.Constraints, func(cc *nodes.ColumnConstraint) bool { return cc.Type == nodes.ColConstrInvisible }) {
			return nil
		}
	}
	return &Error{Code: 4028, SQLState: "HY000", Message: "table must have at least one column that is not invisible"}
}

func validatePartitionClause(tbl *Table, pc *nodes.PartitionClause) error {
	if err := checkPartitionNames(tbl, pc); err != nil {
		return err
	}
	if pc.SubPartExpr != nil || len(pc.SubPartColumns) > 0 {
		for _, pd := range pc.Partitions {
			if err := checkSubpartitionBoundKinds(pc.SubPartType, pd.SubPartitions); err != nil {
				return err
			}
		}
	}
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
	if pc.Type == nodes.PartitionList {
		for _, pd := range pc.Partitions {
			if containsPartitionDefaultExpr(pd.Values) {
				return &Error{Code: 1064, SQLState: "42000", Message: "DEFAULT is not valid in LIST partition values"}
			}
		}
	}
	// MySQL requires every partition to list the same number of subpartitions, or
	// none to. OceanBase does not: 4.3.5.6 takes partitions with different numbers of
	// subpartitions, partitions that list none beside ones that do, and lists that
	// disagree with SUBPARTITIONS N, and prints the first of these back in SHOW CREATE
	// TABLE once a partition is added with subpartitions of its own. Holding it to the
	// MySQL rule refused the engine's own readbacks, so there is no such check here.
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
	return validateUniqueKeysCoverPartitionColumns(tbl, partitionKeyColumns(tbl, pc))
}

func containsPartitionDefaultExpr(n nodes.Node) bool {
	switch v := n.(type) {
	case nil:
		return false
	case *nodes.DefaultExpr:
		return true
	case *nodes.List:
		for _, item := range v.Items {
			if containsPartitionDefaultExpr(item) {
				return true
			}
		}
	}
	return false
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

// partitionValueToString converts a partition value node to SQL string.
//
// Each scalar bound expression is constant-folded (partitionValueItemSQL): MySQL evaluates a
// non-literal RANGE/LIST bound at storage time (`5+5`→`10`, `100 DIV 3`→`33`, `TO_DAYS('2020-01-01')`
// →`737790`) and SHOW CREATE echoes only the folded literal, so the user form must fold to the same
// literal or it phantom-diffs against the readback (verified on live 5.7.25 + 8.0.32). Folding is
// best-effort: anything the folder cannot evaluate (a non-deterministic function, a column
// reference, an unsupported function) is rendered verbatim via nodeToSQL and round-trips only when
// the user already wrote the engine's stored form.
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
					subParts[j] = partitionValueItemSQL(sub.(nodes.ExprNode))
				}
				parts[i] = "(" + strings.Join(subParts, ",") + ")"
			} else if s, ok := item.(*nodes.String); ok {
				// MAXVALUE beside other bounds, as in RANGE COLUMNS
				// VALUES LESS THAN (5, MAXVALUE). Asserting it to an expression
				// panicked.
				parts[i] = s.Str
			} else {
				parts[i] = partitionValueItemSQL(item.(nodes.ExprNode))
			}
		}
		return strings.Join(parts, ",")
	case nodes.ExprNode:
		return partitionValueItemSQL(n)
	default:
		return ""
	}
}

// partitionValueItemSQL renders one scalar partition bound, constant-folding it to an integer
// literal when the expression is a constant integer expression (integer arithmetic, or a supported
// deterministic function on a literal — see foldConstIntExpr). On any non-foldable expression it
// falls back to the verbatim deparse, so a literal bound (the common case) is rendered unchanged and
// an unsupported expression is preserved (and flagged as a documented round-trip limitation).
func partitionValueItemSQL(n nodes.ExprNode) string {
	// An OceanBase LIST subpartition's VALUES IN (DEFAULT), which 4.3.5.6 prints as
	// written. A partition's is refused before it gets here (validatePartitionClause).
	if _, ok := n.(*nodes.DefaultExpr); ok {
		return "DEFAULT"
	}
	if v, ok := foldConstIntExpr(n); ok {
		return strconv.FormatInt(v, 10)
	}
	return nodeToSQL(n)
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

// validateSingleFK validates a single FK constraint against its referenced table, as
// OceanBase does (TestOracle_OceanBaseKeyErrors): the table and the columns on both sides
// are there, the action is not SET DEFAULT, the referenced columns lead an index, and the
// types match. A generated column, VIRTUAL or STORED, may stand on either side.
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
		return errFKNoRefTable()
	}
	for _, colName := range con.Columns {
		if tbl.GetColumn(colName) == nil {
			return errFKNoSuchColumn()
		}
	}
	for _, colName := range con.RefColumns {
		if refTbl.GetColumn(colName) == nil {
			return errFKNoSuchColumn()
		}
	}
	if strings.EqualFold(con.OnDelete, "SET DEFAULT") || strings.EqualFold(con.OnUpdate, "SET DEFAULT") {
		return errCannotAddForeign()
	}

	// The referenced columns are the leading columns of an index, in any order.
	if !hasIndexOnColumns(refTbl, con.RefColumns) {
		return errCannotAddForeign()
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
		if !fkTypesCompatible(col, refCol) {
			return errCannotAddForeign()
		}
	}

	return nil
}

// hasIndexOnColumns checks whether a table has an index (PK, UNIQUE, or regular KEY)
// whose leading columns match the given columns.
func hasIndexOnColumns(tbl *Table, cols []string) bool {
	want := make(map[string]bool, len(cols))
	for _, col := range cols {
		want[toLower(col)] = true
	}
	for _, idx := range tbl.Indexes {
		if len(idx.Columns) < len(want) {
			continue
		}
		match := true
		for _, ic := range idx.Columns[:len(want)] {
			if !want[toLower(ic.Name)] {
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

// validateCompressor refuses a COMPRESSION that names none of OceanBase's compressors
// (ER_WRONG_ARGUMENTS, 1210), MySQL's 'ZLIB' and 'LZ4' among them, and zlib_1.0, which
// 4.3.5 lists and does not support (1235). TestOracle_OceanBaseRefusals holds each name
// to the engine.
func validateCompressor(value string) error {
	name := strings.Trim(value, "'\"")
	switch toLower(name) {
	case "none", "lz4_1.0", "lz4_1.9.1", "snappy_1.0", "zstd_1.0", "zstd_1.3.8":
		return nil
	case "zlib_1.0":
		return &Error{Code: 1235, SQLState: "0A000", Message: "zlib compressor not supported"}
	}
	return &Error{Code: ErrWrongArguments, SQLState: sqlState(ErrWrongArguments),
		Message: fmt.Sprintf("Incorrect arguments to %s", name)}
}

// allocIndexName generates a unique index name based on the first column,
// appending _2, _3, etc. on collision. A key over a column named PRIMARY is named
// PRIMARY on a table without a primary key: OceanBase refuses the name only when it is
// written (TestScenario_C1, 1.10), where MySQL falls back to PRIMARY_2.
func allocIndexName(tbl *Table, baseName string) string {
	candidate := baseName
	suffix := 2
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

// resolveIndexUsingType returns the index type a statement declares, from the two
// positions the grammar lets an index_type clause stand in: before the key parts
// (the dedicated indexType slot) and among the index_option list after them.
//
// A statement may fill both, and the one written last stands, which is why the
// options are scanned after the leading slot rather than before it. Read against
// mysql:8.0 on 2026-11-04 through SHOW CREATE TABLE, which echoes the type it kept:
// `CREATE INDEX i ON t (a) TYPE BTREE` comes back as `USING BTREE`, while
// `CREATE INDEX i TYPE BTREE ON t (a) TYPE HASH` and
// `KEY k USING BTREE (a) USING HASH` come back with no index type at all, the HASH
// each of them kept having become a BTREE index by itself.
//
// Either spelling of the clause, USING or TYPE, reaches here as an option named
// USING. "" means the statement named no index type.
func resolveIndexUsingType(indexType string, opts []*nodes.IndexOption) string {
	declared := strings.ToUpper(indexType)
	for _, opt := range opts {
		if strings.EqualFold(opt.Name, "USING") {
			if s, ok := opt.Value.(*nodes.StringLit); ok {
				declared = strings.ToUpper(s.Value)
			}
		}
	}
	return declared
}

// resolveConstraintIndexType returns the index type from a constraint,
// checking both IndexType (USING before key parts) and IndexOptions (USING after key parts).
func resolveConstraintIndexType(con *nodes.Constraint) string {
	return resolveIndexUsingType(con.IndexType, con.IndexOptions)
}

// applyIndexOptions extracts COMMENT, VISIBLE/INVISIBLE, and KEY_BLOCK_SIZE
// from AST IndexOptions and applies them to the given Index.
// vectorParamDefaults are the parameters 4.3.5.6 fills in when a WITH (...) leaves them
// out, in the order it writes them after the ones that were written. CREATE TABLE ... WITH
// (type=hnsw, lib=vsag, m=8, distance=l2) reads back WITH (TYPE=HNSW, LIB=VSAG, M=8,
// DISTANCE=L2, EF_CONSTRUCTION=200, EF_SEARCH=64).
var vectorParamDefaults = []struct{ name, value string }{
	{"LIB", "VSAG"}, {"M", "16"}, {"EF_CONSTRUCTION", "200"}, {"EF_SEARCH", "64"},
}

// vectorParamsFromOption reads a WITH (...) as SHOW CREATE TABLE writes it: upper-cased,
// in the order written, with the defaults of those left out after them. A name given
// twice keeps its first place and its last value.
func vectorParamsFromOption(opt *nodes.IndexOption) []string {
	values := make(map[string]string, len(opt.Properties))
	var written []string
	for _, prop := range opt.Properties {
		value := ""
		switch v := prop.Value.(type) {
		case *nodes.IntLit:
			value = strconv.FormatInt(v.Value, 10)
		case *nodes.StringLit:
			value = strings.ToUpper(v.Value)
		}
		name := strings.ToUpper(prop.Name)
		if _, seen := values[name]; !seen {
			written = append(written, name)
		}
		values[name] = value
	}
	params := make([]string, 0, len(written)+len(vectorParamDefaults))
	for _, name := range written {
		params = append(params, name+"="+values[name])
	}
	for _, def := range vectorParamDefaults {
		if _, ok := values[def.name]; !ok {
			params = append(params, def.name+"="+def.value)
		}
	}
	return params
}

// validateStoringColumns answers 1054 for a STORING (...) column the table lacks, as
// 4.3.5.6 does.
func validateStoringColumns(tbl *Table, opts []*nodes.IndexOption) error {
	for _, opt := range opts {
		if !strings.EqualFold(opt.Name, "STORING") {
			continue
		}
		for _, name := range opt.Columns {
			if tbl.GetColumn(name) == nil {
				return errNoSuchColumn(name, "'"+tbl.Name+"'")
			}
		}
	}
	return nil
}

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
		case "PARSER":
			// A repeated WITH PARSER takes the last one, as 4.3.5.6 does, and the
			// engine prints the name lower-cased.
			if s, ok := opt.Value.(*nodes.StringLit); ok {
				idx.Parser = strings.ToLower(s.Value)
			}
		case "STORING":
			// A second STORING adds to the first: STORING (b) STORING (c) reads back
			// STORING (`b`, `c`) on 4.3.5.6.
			idx.Storing = append(idx.Storing, opt.Columns...)
		case "WITH":
			// A later WITH replaces an earlier one.
			idx.VectorParams = vectorParamsFromOption(opt)
		case "PARSER_PROPERTIES":
			// A later list replaces an earlier one rather than merging with it:
			// PARSER_PROPERTIES=(min_token_size=5) PARSER_PROPERTIES=(max_token_size=40)
			// reads back min_token_size=3,max_token_size=40 on 4.3.5.6.
			idx.ParserProperties = parserPropertiesFromOption(opt)
		case "LOCAL", "GLOBAL":
			// A repeated scope takes the last one, as 4.3.5.6 does: KEY k (a) LOCAL
			// GLOBAL reads back GLOBAL. The primary key keeps none: ALTER TABLE ...
			// ADD PRIMARY KEY (id) GLOBAL is accepted and reads back as PRIMARY KEY
			// (`id`), and CREATE TABLE refuses a scope on it.
			if !idx.Primary {
				idx.Scope = strings.ToUpper(opt.Name)
			}
		}
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

// applyBinaryModifierCollation resolves the legacy BINARY column modifier
// (e.g. `VARCHAR(n) BINARY`, `CHAR(n) BINARY`, `TEXT BINARY`) to the binary
// collation of the column's effective charset. MySQL treats the trailing
// BINARY attribute as shorthand for `COLLATE <charset>_bin` and stores the
// column with an explicit `CHARACTER SET <cs> COLLATE <cs>_bin` pair (verified
// on 5.7 + 8.0; see normalization.md entry char-binary-attribute). The charset
// is resolved by the caller (column charset if specified, else table/db
// default) before this runs, so col.Charset already holds the effective
// charset. Both the charset and the derived collation must be flagged explicit
// so the normalizer's ResolveColumnCharsetCollation honors the _bin collation
// instead of falling back to the table-default collation — otherwise the
// user's `BINARY` form and the stored `<cs>_bin` form canonicalize differently
// and phantom-diff forever.
func applyBinaryModifierCollation(col *Column, dt *nodes.DataType) {
	if col == nil || dt == nil || !dt.Binary || col.Charset == "" {
		return
	}
	col.Collation = fmt.Sprintf("%s_bin", normalizeCharsetName(col.Charset))
	col.CharsetExplicit = true
	col.CollationExplicit = true
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
		// MySQL adds charset introducer for string literals in generated columns. The value is
		// decoded, so a quote or a backslash in it is written escaped, or the text never closes:
		// OceanBase 4.3.5.6 stores 'it''s' as 'it\'s' and _latin1 X'27' as '\''. It stores no
		// introducer there either, which this package's MySQL-shaped tests still expect.
		escaped := strings.ReplaceAll(strings.ReplaceAll(n.Value, `\`, `\\`), "'", `\'`)
		if charset != "" {
			return "_" + charset + "'" + escaped + "'"
		}
		return "'" + escaped + "'"
	case *nodes.FuncCallExpr:
		funcName := strings.ToLower(n.Name)
		if n.Star {
			return funcName + "(*)"
		}
		if n.JSONValueArm || n.JSONQueryArm {
			// Its clauses change what it answers. OceanBase 4.3.5.6 stores the strings of this
			// call without an introducer, which writes it the way a view does.
			return deparse.Deparse(n)
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
		// A binary operation is written in parentheses of its own, so the pair written
		// around it adds nothing: ((`a` >= 1)) is (`a` >= 1), as in a CHECK and a view.
		if _, binary := n.Expr.(*nodes.BinaryExpr); binary {
			return nodeToSQLGenerated(n.Expr, charset)
		}
		return "(" + nodeToSQLGenerated(n.Expr, charset) + ")"
	case *nodes.KeywordArg:
		// Bare keyword argument (timestampdiff unit, get_format type) —
		// stored unquoted in generated-column bodies exactly like view
		// bodies (oracle 8.0.32: AS (timestampdiff(SECOND,`a`,`b`))).
		return n.Keyword
	case *nodes.ExtractExpr:
		// Oracle 8.0.32 stored form: extract(day_hour from `a`).
		return "extract(" + strings.ToLower(n.Unit) + " from " + nodeToSQLGenerated(n.Expr, charset) + ")"
	case *nodes.IntervalExpr:
		// Oracle 8.0.32 stored form: ((`a` + interval 1 day)).
		return "interval " + nodeToSQLGenerated(n.Value, charset) + " " + strings.ToLower(n.Unit)
	case *nodes.CastExpr:
		// Oracle 8.0.32 stored forms: cast(`a` as char(4) charset latin1),
		// cast(`s` as signed), and the WEIGHT_STRING(... AS BINARY(n))
		// desugar weight_string(cast(`s` as char(4) charset binary)).
		result := "cast(" + nodeToSQLGenerated(n.Expr, charset) + " as " + castTypeToSQLGenerated(n.TypeName)
		if n.Array {
			// Multi-valued index form: stored with a lowercase trailing array.
			result += " array"
		}
		return result + ")"
	case *nodes.PositionExpr:
		// Oracle mysql:8.0 stored form: locate(_utf8mb4'x',`s`) — the engine
		// keeps the function under its other name.
		return "locate(" + nodeToSQLGenerated(n.Substr, charset) + "," + nodeToSQLGenerated(n.Str, charset) + ")"
	case *nodes.WeightStringExpr:
		// Oracle 8.0.32 stored form: weight_string(`s` as char(4)).
		var b strings.Builder
		b.WriteString("weight_string(")
		b.WriteString(nodeToSQLGenerated(n.Expr, charset))
		if n.AsChar != nil {
			fmt.Fprintf(&b, " as char(%d)", n.AsChar.Length)
		}
		for i, lv := range n.Levels {
			if i == 0 {
				b.WriteString(" level ")
			} else {
				b.WriteString(",")
			}
			fmt.Fprintf(&b, "%d", lv.Level)
			if lv.Desc {
				b.WriteString(" desc")
			}
			if lv.Reverse {
				b.WriteString(" reverse")
			}
		}
		b.WriteString(")")
		return b.String()
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
	case *nodes.CaseExpr:
		// CASE, which is also what deparse.RewriteExpr makes of IF(c, a, b).
		var b strings.Builder
		b.WriteString("(case")
		if n.Operand != nil {
			b.WriteString(" " + nodeToSQLGenerated(n.Operand, charset))
		}
		for _, w := range n.Whens {
			b.WriteString(" when " + nodeToSQLGenerated(w.Cond, charset) + " then " + nodeToSQLGenerated(w.Result, charset))
		}
		if n.Default != nil {
			b.WriteString(" else " + nodeToSQLGenerated(n.Default, charset))
		}
		b.WriteString(" end)")
		return b.String()
	default:
		// Any other expression is written as the deparser writes it everywhere else. A
		// placeholder here would be stored as the column's expression and sent to the
		// engine in a plan.
		return deparse.Deparse(n)
	}
}

// castTypeToSQLGenerated renders a CAST target type in the form SHOW CREATE
// TABLE stores inside generated-column expressions (oracle 8.0.32 + 5.7.25):
// lowercase names, length/precision kept, and the charset attribute rendered
// only when the parsed type carries one — readbacks always do for CHAR
// (cast(`a` as char(4) charset latin1)); a user form without a charset is
// rendered as written.
func castTypeToSQLGenerated(dt *nodes.DataType) string {
	if dt == nil {
		return ""
	}
	name := strings.ToLower(dt.Name)
	if name == "array" {
		return deparse.CastArrayType(dt)
	}
	if name == "binary" {
		// CAST(x AS BINARY[(n)]) is stored as cast(x as char[(n)] charset
		// binary) — same rewrite the view deparser applies (oracle 8.0.32 +
		// 5.7.25, view/generated/CHECK positions alike).
		result := "char"
		if dt.Length > 0 {
			result += fmt.Sprintf("(%d)", dt.Length)
		}
		return result + " charset binary"
	}
	result := name
	switch name {
	case "decimal":
		if dt.Length > 0 || dt.Scale > 0 {
			result += fmt.Sprintf("(%d,%d)", dt.Length, dt.Scale)
		}
	case "datetime", "time":
		if dt.Length > 0 {
			result += fmt.Sprintf("(%d)", dt.Length)
		}
	default:
		if dt.Length > 0 {
			result += fmt.Sprintf("(%d)", dt.Length)
		}
	}
	if dt.Charset != "" {
		result += " charset " + strings.ToLower(dt.Charset)
	}
	return result
}

// nodeToSQL writes an expression as the catalog stores it: a CHECK, a generated column, a
// default, an index's key part. BETWEEN over constant bounds is two comparisons there,
// as OceanBase stores it (deparse.FoldBetween).
func nodeToSQL(node nodes.ExprNode) string {
	if node == nil {
		return deparse.Deparse(node)
	}
	return deparse.Deparse(deparse.FoldBetween(node))
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
	// TEXT(n) and BLOB(n) are the smallest type that holds n, written without it.
	if name == "text" && dt.Length > 0 {
		return promoteTextType(dt.Length)
	}
	if name == "blob" && dt.Length > 0 {
		return promoteBlobType(dt.Length)
	}

	// Type aliases: BOOLEAN/BOOL → tinyint(1), NUMERIC → decimal, SERIAL → bigint(20) unsigned
	switch name {
	case "boolean", "bool":
		return "tinyint(1)"
	case "serial":
		return "bigint(20) unsigned"
	}

	// OceanBase writes the vector type in upper case, the one type it does: `v` VECTOR(3).
	if name == "vector" {
		return fmt.Sprintf("VECTOR(%d)", dt.Length)
	}

	var buf strings.Builder
	buf.WriteString(name)

	// OceanBase keeps an integer's display width, as MySQL 5.7 does and 8.0 does not: the
	// one written, or the type's own, int(11) and int(10) unsigned. ZEROFILL makes the
	// column UNSIGNED, so a bare INT ZEROFILL is int(10) unsigned zerofill.
	// TestOracle_OceanBaseColumnTypes holds every width to the engine.
	isIntType := isIntegerType(name)
	if isIntType {
		width := dt.Length
		if width == 0 {
			width = defaultIntDisplayWidth(name, dt.Unsigned || dt.Zerofill)
		}
		fmt.Fprintf(&buf, "(%d)", width)
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
		// A YEAR of any width, or of none, is year(4).
		buf.WriteString("(4)")
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
	if name == "text" && dt.Length > 0 {
		return promoteTextType(dt.Length)
	}
	if name == "blob" && dt.Length > 0 {
		return promoteBlobType(dt.Length)
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

// normalizeCharsetName lower-cases a character set's name. utf8 and utf8mb3 are utf8mb4
// in OceanBase, which has no three-byte UTF-8 (TestOracle_OceanBaseColumnTypes).
func normalizeCharsetName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "utf8" || name == "utf8mb3" {
		return "utf8mb4"
	}
	return name
}

// normalizeCollationName writes a collation of utf8 or utf8mb3 as utf8mb4's, which is
// what OceanBase stores: utf8_general_ci is utf8mb4_general_ci.
func normalizeCollationName(name string) string {
	for _, prefix := range []string{"utf8_", "utf8mb3_"} {
		if len(name) > len(prefix) && strings.EqualFold(name[:len(prefix)], prefix) {
			return "utf8mb4_" + name[len(prefix):]
		}
	}
	return name
}

func charsetMatchesCollation(charset, collation string) bool {
	derived := normalizeCharsetName(charsetForCollation(collation))
	return derived == "" || strings.EqualFold(normalizeCharsetName(charset), derived)
}

// promoteBlobType is the smallest BLOB type that holds n bytes.
func promoteBlobType(n int) string {
	switch {
	case n <= 255:
		return "tinyblob"
	case n <= 65535:
		return "blob"
	case n <= 16777215:
		return "mediumblob"
	default:
		return "longblob"
	}
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
// MySQL escapes a single quote in enum values by doubling it.
func escapeEnumValue(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// createTableLike implements CREATE TABLE t2 LIKE t1.
// It copies the structure (columns, indexes, constraints) from the source table.
func (c *Catalog) createTableLike(db *Database, tableName, key string, stmt *nodes.CreateTableStmt, nameTaken func() (bool, error)) error {
	// Resolve source table.
	srcDBName := stmt.Like.Schema
	srcDB, err := c.resolveDatabase(srcDBName)
	if err != nil {
		return err
	}
	// A table copied onto its own name is 1066, before the source is looked up (1146)
	// and before the name (1050).
	if srcDB == db && toLower(stmt.Like.Name) == key {
		return errNotUniqueTable(tableName)
	}
	srcTbl := srcDB.GetTable(stmt.Like.Name)
	if srcTbl == nil {
		return errNoSuchTable(srcDB.Name, stmt.Like.Name)
	}
	if done, err := nameTaken(); done || err != nil {
		return err
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
			OnUpdateKind:  srcCol.OnUpdateKind,
			Invisible:     srcCol.Invisible,
			Hidden:        srcCol.Hidden,
			// SRID copied by presence: MySQL's CREATE TABLE ... LIKE preserves the source
			// column's SRID clause (incl. an explicit `SRID 0`), so the cloned column must
			// carry both the value and its presence flag, not silently drop to no-SRID.
			SRID:    srcCol.SRID,
			HasSRID: srcCol.HasSRID,
		}
		if srcCol.Default != nil {
			def := *srcCol.Default
			col.Default = &def
			col.DefaultKind = srcCol.DefaultKind
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
			Vector:       srcIdx.Vector,
			VectorParams: srcIdx.VectorParams,
			Storing:      slices.Clone(srcIdx.Storing),
			IndexType:    srcIdx.IndexType,
			Visible:      srcIdx.Visible,
			Comment:      srcIdx.Comment,
			KeyBlockSize: srcIdx.KeyBlockSize,
			Scope:        srcIdx.Scope,
			// OceanBase copies a fulltext index's parser with the table
			// (TestOracle_OceanBaseFulltextParserEcho).
			Parser:           srcIdx.Parser,
			ParserProperties: copyParserProperties(srcIdx.ParserProperties),
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
		// An omitted action reads RESTRICT in OceanBase, in information_schema and in
		// SHOW CREATE TABLE, where MySQL reads NO ACTION (TestScenario_C21, 21.5).
		return "RESTRICT"
	}
}
