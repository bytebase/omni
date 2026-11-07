package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/deparse"
	tidbparser "github.com/bytebase/omni/tidb/parser"
)

// createTable installs the table stmt defines. declared is true when stmt lists
// its indexes in the order they were declared, as a CREATE TABLE statement does,
// and false for a loader's snapshot, which lists them in an order of its own.
func (c *Catalog) createTable(stmt *nodes.CreateTableStmt, declared bool) error {
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
		return errDupTable(db.Name + "." + tableName)
	}
	if db.Views[key] != nil {
		return errDupTable(db.Name + "." + tableName)
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
		Engine:      "InnoDB",
		Temporary:   stmt.Temporary,
		// A global temporary table is its own kind: an ordinary schema object
		// whose rows are per-session, so Temporary is not set for it and the
		// ON COMMIT clause it reports is kept beside it.
		GlobalTemporary: stmt.Global,
		OnCommit:        stmt.OnCommit,
	}

	// Apply table options.
	tblCharsetExplicit := false
	tblCollationExplicit := false
	ttlEnableSet := false
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
		case "compression":
			tbl.Compression = strings.Trim(opt.Value, "'\"")
		case "auto_increment":
			fmt.Sscanf(opt.Value, "%d", &tbl.AutoIncrement)
		case "row_format":
			tbl.RowFormat = opt.Value
		case "key_block_size":
			fmt.Sscanf(opt.Value, "%d", &tbl.KeyBlockSize)
		// TiDB-specific table options.
		case "shard_row_id_bits":
			fmt.Sscanf(opt.Value, "%d", &tbl.ShardRowIDBits)
		case "pre_split_regions":
			fmt.Sscanf(opt.Value, "%d", &tbl.PreSplitRegions)
		case "auto_id_cache":
			fmt.Sscanf(opt.Value, "%d", &tbl.AutoIDCache)
		case "auto_random_base":
			fmt.Sscanf(opt.Value, "%d", &tbl.AutoRandomBase)
		case "placement policy":
			if err := c.validatePolicyRef(opt.Value); err != nil {
				return err
			}
			tbl.PlacementPolicy = resolvePolicyRef(opt.Value)
		case "ttl":
			col, interval, err := extractTTLParts(opt.Value)
			if err != nil {
				return err
			}
			tbl.TTLColumn = col
			tbl.TTLInterval = interval
			if !ttlEnableSet {
				// TTL_ENABLE defaults to ON once a table has a TTL.
				tbl.TTLEnable = true
			}
		case "ttl_enable":
			tbl.TTLEnable = strings.EqualFold(opt.Value, "ON")
			ttlEnableSet = true
		case "ttl_job_interval":
			tbl.TTLJobInterval = opt.Value
		}
	}
	// A charset or collation the statement names is checked against TiDB's
	// list, spelled as TiDB spells it, and completes the other half.
	if tblCharsetExplicit || tblCollationExplicit {
		cs, coll := "", ""
		if tblCharsetExplicit {
			cs = tbl.Charset
		}
		if tblCollationExplicit {
			coll = tbl.Collation
		}
		cs, coll, err := resolveCharsetCollation(cs, coll)
		if err != nil {
			return err
		}
		tbl.Charset, tbl.Collation = cs, coll
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

	// foreignKeyID numbers the foreign keys of this CREATE TABLE in the order
	// they are declared, named or not, and an unnamed one takes fk_<its number>
	// (buildTableInfo in TiDB v8.5.5 pkg/ddl/create_table.go).
	var foreignKeyID int

	// TiDB reserves explicit names before naming unnamed CHECKs, skipping
	// occupied suffixes. See setEmptyCheckConstraintName in pkg/ddl/create_table.go.
	checkNames := make(map[string]bool)
	for _, con := range stmt.Constraints {
		switch con.Type {
		case nodes.ConstrUnique, nodes.ConstrIndex, nodes.ConstrFulltextIndex, nodes.ConstrSpatialIndex, nodes.ConstrVector:
			if strings.EqualFold(con.Name, "PRIMARY") {
				return errWrongNameForIndex(con.Name)
			}
		}
		if con.Type == nodes.ConstrCheck && con.Name != "" {
			checkNames[toLower(con.Name)] = true
		}
	}
	for _, col := range stmt.Columns {
		for _, con := range col.Constraints {
			if con.Type == nodes.ColConstrCheck && con.Name != "" {
				checkNames[toLower(con.Name)] = true
			}
		}
	}

	// columnChecks are the CHECK clauses written on columns. TiDB v8.5.5 adds
	// them after the table-level constraints and numbers their unnamed ones
	// after the table-level ones: `a INT CHECK (a > 0), CHECK (a <> b)` gives
	// t_chk_1 for the second and t_chk_2 for the first, listed in that order
	// (read on 2026-11-06, with tidb_enable_check_constraint ON).
	var columnChecks []*nodes.ColumnConstraint
	columnCheckColumns := make(map[*nodes.ColumnConstraint]string)
	primaryColumns := make(map[string]bool)
	for _, con := range stmt.Constraints {
		if con.Type == nodes.ConstrPrimaryKey {
			for _, name := range extractColumnNames(con) {
				primaryColumns[toLower(name)] = true
			}
		}
	}

	// Process columns.
	for i, colDef := range stmt.Columns {
		var err error
		colDef, err = normalizeTiDBStringType(colDef, tbl.Charset, c.strictDDL)
		if err != nil {
			return err
		}
		if err := validateTiDBColumn(colDef); err != nil {
			return err
		}
		if err := validateTiDBColumnNullability(colDef, primaryColumns[toLower(colDef.Name)]); err != nil {
			return err
		}
		if err := checkColumnQualifier(colDef.Schema, colDef.Table, db.Name, tableName); err != nil {
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
				col.DataType = tiDBDataType(colDef.TypeName)
				col.ColumnType = formatColumnType(colDef.TypeName)
			}
			cs, coll, err := resolveCharsetCollation(colDef.TypeName.Charset, colDef.TypeName.Collate)
			col.Zerofill = colDef.TypeName.Zerofill && col.DataType != "year"
			if err != nil {
				return err
			}
			if cs != "" {
				col.Charset, col.Collation = cs, coll
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

		// Standalone BINARY modifier (CHAR(10) BINARY) → the effective charset's
		// binary collation, matching TiDB. Resolved after the charset default
		// above so a bare CHAR BINARY uses the table/db charset. An explicit
		// COLLATE wins over BINARY, so only apply when no COLLATE was given.
		if colDef.TypeName != nil && colDef.TypeName.Binary && isStringType(col.DataType) &&
			colDef.TypeName.Collate == "" {
			col.Collation = binModifierCollation(col.Charset)
		}

		// Top-level column properties.
		if colDef.TypeName != nil && colDef.TypeName.SRID != 0 {
			col.SRID = colDef.TypeName.SRID
		}
		if colDef.AutoIncrement {
			col.AutoIncrement = true
			col.Nullable = false
		}
		if colDef.AutoRandom {
			col.AutoRandom = true
			col.AutoRandomShardBits = colDef.AutoRandomShardBits
			col.AutoRandomRangeBits = colDef.AutoRandomRangeBits
		}
		if colDef.Comment != "" {
			col.Comment = colDef.Comment
		}
		if colDef.DefaultValue != nil {
			s := c.defaultToSQL(colDef.DefaultValue)
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
					s := c.defaultToSQL(cc.Expr)
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
				// Added after the table's own constraints: see columnChecks.
				if c.checkConstraints {
					columnChecks = append(columnChecks, cc)
					columnCheckColumns[cc] = colDef.Name
				}
			case nodes.ColConstrReferences:
				// TiDB reads a REFERENCES clause on a column and keeps nothing:
				// v8.5.5 adds no foreign key and no index for
				// `a INT REFERENCES p(id)` (read on 2026-11-06), as MySQL does.
			case nodes.ColConstrCollate:
				// Collation specified via constraint.
				if cc.Expr != nil {
					if s, ok := cc.Expr.(*nodes.StringLit); ok {
						if err := applyColumnCollate(col, s.Value, colDef.TypeName != nil && colDef.TypeName.Charset != ""); err != nil {
							return err
						}
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

		if err := checkLiteralDefaultAllowed(col); err != nil {
			return err
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
					// TiDB: hoist CLUSTERED/NONCLUSTERED from inline column PK
					// (`id INT PRIMARY KEY CLUSTERED`) onto the synthesized
					// table-level constraint so that later lookups by
					// con.Type == ConPrimaryKey see the flag.
					Clustered: cc.Clustered,
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

	// The indexes so far come from column definitions. TiDB checks them after
	// the table-level ones, when it looks for an implicit primary key and when
	// it holds unique keys to the partitioning columns.
	columnIndexes := len(tbl.Indexes)

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
			// TiDB v8.5.5 refuses the table rather than create an invisible key,
			// as its ALTER TABLE ... ADD PRIMARY KEY path does.
			if !pkIdx.Visible {
				return errPrimaryKeyInvisible()
			}
			tbl.Indexes = append(tbl.Indexes, pkIdx)
			tbl.Constraints = append(tbl.Constraints, &Constraint{
				Name:      "PRIMARY",
				Type:      ConPrimaryKey,
				Table:     tbl,
				Columns:   cols,
				IndexName: "PRIMARY",
				Clustered: con.Clustered, // TiDB: CLUSTERED/NONCLUSTERED on table-level PK
			})

		case nodes.ConstrUnique:
			idxName := con.Name
			if idxName == "" {
				idxName = implicitTiDBIndexName(tbl, con)
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
			foreignKeyID++
			conName := fkName(con)
			if conName == "" {
				conName = fmt.Sprintf("fk_%d", foreignKeyID)
			}
			if findForeignKey(tbl, conName) != nil {
				return errCannotAddForeign()
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
			pendingFKs = append(pendingFKs, pendingFK{conName: conName, cols: cols, idxCols: buildIndexColumns(con)})

		case nodes.ConstrCheck:
			if !c.checkConstraints {
				break
			}
			conName := con.Name
			if conName == "" {
				conName = allocateTiDBCheckName(tableName, checkNames)
			}
			if err := validateTiDBCheck(tbl, conName, con.Expr, ""); err != nil {
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
			idxName := con.Name
			if idxName == "" {
				idxName = implicitTiDBIndexName(tbl, con)
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
			// v8.5.5 reads a FULLTEXT key in CREATE TABLE and keeps nothing:
			// SHOW CREATE TABLE lists no index and SHOW WARNINGS is empty (read
			// on 2026-11-06). CREATE FULLTEXT INDEX is refused instead.

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

		case nodes.ConstrVector:
			// TiDB's VECTOR INDEX, which CREATE TABLE declares as well as
			// ALTER TABLE adds. The index type is the USING name (HNSW).
			vecIdx := &Index{
				Name:      vectorIndexName(tbl, con, cols),
				Table:     tbl,
				Columns:   buildIndexColumns(con),
				Vector:    true,
				IndexType: vectorIndexHNSW,
				Visible:   true,
			}
			applyIndexOptions(vecIdx, con.IndexOptions)
			tbl.Indexes = append(tbl.Indexes, vecIdx)
		}
	}

	// Without a primary key TiDB v8.5.5 refuses a table whose implicit one is
	// invisible. A snapshot does not keep the declaration order the check needs,
	// and the table it describes passed the check on the server.
	if !hasPK && declared {
		if err := checkImplicitPrimaryKeyVisible(tbl, tbl.Indexes[columnIndexes:], tbl.Indexes[:columnIndexes]); err != nil {
			return err
		}
	}

	for _, cc := range columnChecks {
		conName := cc.Name
		if conName == "" {
			conName = allocateTiDBCheckName(tableName, checkNames)
		}
		if err := validateTiDBCheck(tbl, conName, cc.Expr, columnCheckColumns[cc]); err != nil {
			return err
		}
		tbl.Constraints = append(tbl.Constraints, &Constraint{
			Name:        conName,
			Type:        ConCheck,
			Table:       tbl,
			CheckExpr:   nodeToSQL(cc.Expr),
			NotEnforced: cc.NotEnforced,
		})
	}

	// TiDB keeps a table's keys in the order buildTableInfo meets them: the
	// table-level ones first, then those written on columns, in column order
	// (`b INT UNIQUE, KEY k (a)` lists k before b, read on 2026-11-06). The
	// catalog built the column-level ones first, for the implicit primary key
	// check above, and moves them after the rest here.
	if columnIndexes > 0 && columnIndexes < len(tbl.Indexes) {
		reordered := make([]*Index, 0, len(tbl.Indexes))
		reordered = append(reordered, tbl.Indexes[columnIndexes:]...)
		reordered = append(reordered, tbl.Indexes[:columnIndexes]...)
		tbl.Indexes = reordered
	}

	// TiDB numbers foreign keys by ID, not by name: a table created with one
	// named fk_7 still gives the next unnamed one fk_2 (read on 2026-11-06).
	// A table loaded from a snapshot comes through here too, with only its
	// foreign keys' count to go on.
	tbl.maxForeignKeyID = foreignKeyID

	// Process deferred FK backing indexes now that all explicit indexes are in place.
	for _, fk := range pendingFKs {
		if err := ensureFKBackingIndex(tbl, fk.conName, fk.cols, fk.idxCols); err != nil {
			return err
		}
	}

	// Structural checks still apply when foreign_key_checks permits a missing parent.
	for _, idx := range tbl.Indexes {
		if err := validateTiDBIndexColumns(tbl, idx.Columns); err != nil {
			return err
		}
		if err := c.prepareTiDBIndexExpressions(tbl, idx); err != nil {
			return err
		}
	}
	if err := c.validateForeignKeys(db, tbl); err != nil {
		return err
	}

	if err := validateTiDBPrimaryGenerated(tbl); err != nil {
		return err
	}

	// Process partition clause. v8.5.5 runs SYSTEM_TIME partitioning and leaves the table with
	// none, whatever its keys cover.
	if stmt.Partitions != nil && stmt.Partitions.Type != nodes.PartitionSystemTime {
		if err := validateTiDBPartition(tbl, stmt.Partitions); err != nil {
			return err
		}
		if err := validateCreateTableUniqueKeys(tbl, stmt.Partitions); err != nil {
			return err
		}
		tbl.Partitioning = buildPartitionInfo(stmt.Partitions)
		tbl.Partitioning.keyColumns = partitionKeyColumns(tbl, stmt.Partitions)
		if stmt.Partitions.Interval != nil {
			if err := applyPartitionInterval(tbl, stmt.Partitions, tbl.Partitioning); err != nil {
				return err
			}
		}
	}

	// Phase 3: analyze DEFAULT, GENERATED, and CHECK expressions now that all
	// columns are present in the table.
	c.analyzeTableExpressions(tbl, stmt)

	db.Tables[key] = tbl
	return nil
}

// analyzeTableExpressions performs best-effort semantic analysis on DEFAULT,
// GENERATED, and CHECK expressions after all columns have been added to the table.
func (c *Catalog) analyzeTableExpressions(tbl *Table, stmt *nodes.CreateTableStmt) {
	// Analyze DEFAULT and GENERATED expressions from column definitions.
	for i, colDef := range stmt.Columns {
		if i >= len(tbl.Columns) {
			break
		}
		col := tbl.Columns[i]

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

	// Analyze CHECK expressions on constraints, in the order they were added:
	// the table-level ones, then those written on columns.
	checkIdx := 0
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
}

// buildPartitionInfo converts an AST PartitionClause to a catalog PartitionInfo,
// keeping what TiDB keeps. v8.5.5 reads LINEAR, KEY ALGORITHM, a partition's
// ENGINE and the whole SUBPARTITION clause and stores none of them: SHOW CREATE
// TABLE writes `LINEAR HASH (a)` as `HASH (a)` and a subpartitioned table with no
// subpartitions (read on 2026-11-06).
func buildPartitionInfo(pc *nodes.PartitionClause) *PartitionInfo {
	pi := &PartitionInfo{
		NumParts: pc.NumParts,
	}

	switch pc.Type {
	case nodes.PartitionRange:
		if len(pc.Columns) > 0 {
			pi.Type = "RANGE COLUMNS"
			pi.Columns = pc.Columns
		} else {
			pi.Type = "RANGE"
			pi.Expr = partitionExpressionSQL(pc.Expr)
		}
	case nodes.PartitionList:
		if len(pc.Columns) > 0 {
			pi.Type = "LIST COLUMNS"
			pi.Columns = pc.Columns
		} else {
			pi.Type = "LIST"
			pi.Expr = partitionExpressionSQL(pc.Expr)
		}
	case nodes.PartitionHash:
		pi.Type = "HASH"
		pi.Expr = partitionExpressionSQL(pc.Expr)
	case nodes.PartitionKey:
		pi.Type = "KEY"
		pi.Columns = pc.Columns
	}

	// Partition definitions.
	for _, pd := range pc.Partitions {
		pdi := &PartitionDefInfo{
			Name: pd.Name,
		}
		if pd.Values != nil {
			pdi.ValueExpr = partitionValueToString(pd.Values, pc.Type)
		}
		for _, opt := range pd.Options {
			if toLower(opt.Name) == "comment" {
				pdi.Comment = opt.Value
			}
		}
		pi.Partitions = append(pi.Partitions, pdi)
	}

	// Auto-generate partition definitions for HASH/KEY when PARTITIONS N is
	// specified without explicit partition definitions: p0, p1, p2, ...
	if (pc.Type == nodes.PartitionHash || pc.Type == nodes.PartitionKey) && pi.NumParts == 0 {
		pi.NumParts = len(pi.Partitions)
		if pi.NumParts == 0 {
			pi.NumParts = 1
		}
	}
	if len(pi.Partitions) == 0 && pi.NumParts > 0 {
		for i := 0; i < pi.NumParts; i++ {
			pi.Partitions = append(pi.Partitions, &PartitionDefInfo{
				Name: fmt.Sprintf("p%d", i),
			})
		}
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
			if subList, ok := item.(*nodes.List); ok {
				// Tuple: (val1, val2) for multi-column LIST COLUMNS
				subParts := make([]string, len(subList.Items))
				for j, sub := range subList.Items {
					subParts[j] = partitionBoundSQL(sub)
				}
				parts[i] = "(" + strings.Join(subParts, ",") + ")"
			} else {
				parts[i] = partitionBoundSQL(item)
			}
		}
		return strings.Join(parts, ",")
	case nodes.ExprNode:
		return nodeToSQL(n)
	default:
		return ""
	}
}

// partitionBoundSQL renders one bound in a partition's value list. MAXVALUE
// reaches it as a String beside the expressions, as in RANGE COLUMNS
// VALUES LESS THAN (5, MAXVALUE); asserting every item to an expression
// panicked Exec on that statement.
func partitionBoundSQL(n nodes.Node) string {
	switch n := n.(type) {
	case *nodes.String:
		return n.Str
	case nodes.ExprNode:
		return nodeToSQL(n)
	}
	return ""
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
	// TiDB validates the child definition even with foreign_key_checks=0.
	// https://github.com/pingcap/tidb/blob/v8.5.5/pkg/ddl/create_table.go
	for _, name := range con.Columns {
		col := tbl.GetColumn(name)
		if col == nil {
			continue
		}
		if col.Generated != nil && !col.Generated.Stored {
			return tiDBVirtualFKColumn(con.Name, name)
		}
		if !col.Nullable && (con.OnDelete == "SET NULL" || con.OnUpdate == "SET NULL") {
			return &Error{Code: 1830, SQLState: "HY000", Message: fmt.Sprintf("Column '%s' cannot be NOT NULL: needed in a foreign key constraint '%s' SET NULL", name, con.Name)}
		}
	}
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
		if !c.foreignKeyChecks {
			return nil
		}
		return errFKNoRefTable(con.RefTable)
	}
	// Parent virtual columns are rejected before checking the backing index.
	// https://github.com/pingcap/tidb/blob/v8.5.5/pkg/ddl/foreign_key.go
	for _, name := range con.RefColumns {
		col := refTbl.GetColumn(name)
		if col == nil {
			return &Error{Code: 3734, SQLState: "HY000", Message: fmt.Sprintf("Failed to add the foreign key constraint. Missing column '%s' for constraint '%s' in the referenced table '%s'", name, con.Name, con.RefTable)}
		}
		if col.Generated != nil && !col.Generated.Stored {
			return tiDBVirtualFKColumn(con.Name, name)
		}
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
			// `((v))` is a functional key part over a column, which a vector index requires.
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
	for strings.EqualFold(candidate, "PRIMARY") || indexNameExists(tbl, candidate) {
		candidate = fmt.Sprintf("%s_%d", baseName, suffix)
		suffix++
	}
	return candidate
}

// TiDB calls unnamed expression indexes expression_index, using the first
// key part's column name only when that part is a plain column reference.
// https://github.com/pingcap/tidb/blob/v8.5.5/pkg/ddl/index.go
func implicitTiDBIndexName(tbl *Table, con *nodes.Constraint) string {
	if len(con.IndexColumns) > 0 {
		first := con.IndexColumns[0]
		if col, ok := first.Expr.(*nodes.ColumnRef); ok && !first.Functional {
			return allocIndexName(tbl, col.Column)
		}
		return allocIndexName(tbl, "expression_index")
	}
	if len(con.Columns) > 0 {
		return allocIndexName(tbl, con.Columns[0])
	}
	return allocIndexName(tbl, "expression_index")
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

// fkName returns the name a foreign key takes, and the index it adds when
// none serves it: the CONSTRAINT symbol, else the name after FOREIGN KEY,
// else "" for a generated one (verified on v8.5.5).
func fkName(con *nodes.Constraint) string {
	if con.Name != "" {
		return con.Name
	}
	return con.IndexName
}

// ensureFKBackingIndex creates an implicit backing index for FK columns
// if no existing index already covers them (MySQL 8.0 behavior).
// MySQL uses the constraint name as the index name when provided;
// otherwise falls back to the first column name via allocIndexName.
func ensureFKBackingIndex(tbl *Table, conName string, cols []string, idxCols []*IndexColumn) error {
	if hasIndexCoveringColumns(tbl, cols) {
		return nil
	}
	idxName := conName
	if idxName == "" {
		idxName = allocIndexName(tbl, cols[0])
	}
	// TiDB names the index after the foreign key and refuses the statement
	// when another index already has that name, rather than choosing another:
	// 1061 from addIndexForForeignKey in CREATE TABLE and from createIndex in
	// ALTER TABLE (pkg/ddl/create_table.go, executor.go).
	if indexNameExists(tbl, idxName) {
		return errDupKeyName(idxName)
	}
	tbl.Indexes = append(tbl.Indexes, &Index{
		Name:    idxName,
		Table:   tbl,
		Columns: idxCols,
		Visible: true,
	})
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
	if declared == "HASH" {
		return "BTREE"
	}
	return declared
}

// resolveConstraintIndexType returns the index type from a constraint, from either
// of the positions resolveIndexUsingType reads.
func resolveConstraintIndexType(con *nodes.Constraint) string {
	return resolveIndexUsingType(con.IndexType, con.IndexOptions)
}

// vectorIndexName returns the name a vector index goes by: its own, or for one
// left unnamed, the name of its one column or, for its usual functional key part,
// the name v8.5.5 gives any unnamed key over an expression: expression_index, then
// expression_index_2.
func vectorIndexName(tbl *Table, con *nodes.Constraint, cols []string) string {
	if con.Name != "" {
		return con.Name
	}
	if len(cols) > 0 {
		return allocIndexName(tbl, cols[0])
	}
	return allocIndexName(tbl, "expression_index")
}

// vectorIndexHNSW is a vector index's type. It is the only one it has, and the grammar's actions
// refuse any other, so every vector index that parses has it, whichever way it spelled that:
// CREATE VECTOR INDEX lets a type among the options win and a constraint lets the one before the
// key parts win, and with neither the action writes HNSW itself.
//
// Ref: CreateIndexStmt and ConstraintVectorIndex in TiDB v8.5.5 pkg/parser/parser.y.
const vectorIndexHNSW = "HNSW"

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
		// TiDB's GLOBAL and LOCAL. The option list may hold both, and the
		// last one decides, as the grammar's own action does. Dropping GLOBAL
		// here left the partition-coverage check applying the local rule to a
		// global index, which is the 8264 that option exists to avoid.
		case "GLOBAL":
			idx.Global = true
		case "LOCAL":
			idx.Global = false
		}
	}
}

func findForeignKey(tbl *Table, name string) *Constraint {
	for _, con := range tbl.Constraints {
		if con.Type == ConForeignKey && eqFoldStr(con.Name, name) {
			return con
		}
	}
	return nil
}

// errCannotAddForeign is 1215, which TiDB answers for a foreign key whose name
// another one on the table already has.
func errCannotAddForeign() error {
	return &Error{Code: 1215, SQLState: "HY000", Message: "Cannot add foreign key constraint"}
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
	return storedTiDBExpressionSQL(node, charset, false)
}

func storedTiDBExpressionSQL(node nodes.ExprNode, charset string, index bool) string {
	if node == nil {
		return ""
	}
	switch n := node.(type) {
	case *nodes.ColumnRef:
		if n.Table != "" {
			return quoteIdent(n.Table) + "." + quoteIdent(n.Column)
		}
		return quoteIdent(n.Column)
	case *nodes.IntLit:
		return fmt.Sprintf("%d", n.Value)
	case *nodes.StringLit:
		// MySQL adds charset introducer for string literals in generated columns.
		if n.Charset != "" {
			charset = n.Charset
		}
		if charset != "" {
			return "_" + charset + "'" + outputFormat(n.Value) + "'"
		}
		return "'" + outputFormat(n.Value) + "'"
	case *nodes.FuncCallExpr:
		funcName := strings.ToLower(n.Name)
		if n.Star {
			return funcName + "(*)"
		}
		var args []string
		for _, a := range n.Args {
			args = append(args, storedTiDBExpressionSQL(a, charset, index))
		}
		return funcName + "(" + strings.Join(args, ", ") + ")"
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
		return "(" + storedTiDBExpressionSQL(n.Expr, charset, index) + ")"
	case *nodes.KeywordArg:
		// Bare keyword argument (timestampdiff unit, get_format type) —
		// stored unquoted in generated-column bodies exactly like view bodies.
		return n.Keyword
	case *nodes.ExtractExpr:
		// Oracle TiDB v8.5.5 stored form: extract(day_hour from `a`).
		return "extract(" + strings.ToLower(n.Unit) + " from " + storedTiDBExpressionSQL(n.Expr, charset, index) + ")"
	case *nodes.PositionExpr:
		// Oracle TiDB v8.5.5 stored form: position(_utf8mb4'x' in `s`) — the
		// keyword, where MySQL 8.0 stores locate(_utf8mb4'x',`s`).
		return "position(" + storedTiDBExpressionSQL(n.Substr, charset, index) + " in " + storedTiDBExpressionSQL(n.Str, charset, index) + ")"
	case *nodes.BinaryExpr:
		left := storedTiDBExpressionSQL(n.Left, charset, index)
		right := storedTiDBExpressionSQL(n.Right, charset, index)
		// MySQL rewrites JSON operators to function calls in generated column expressions.
		separator := ","
		if index {
			separator = ", "
		}
		switch n.Op {
		case nodes.BinOpJsonExtract:
			return "json_extract(" + left + separator + right + ")"
		case nodes.BinOpJsonUnquote:
			return "json_unquote(json_extract(" + left + separator + right + "))"
		}
		op := binaryOpToString(n.Op)
		return "(" + left + " " + op + " " + right + ")"
	case *nodes.CastExpr:
		target := formatDataType(n.TypeName)
		if n.TypeName != nil && n.TypeName.Charset != "" {
			target += " charset " + n.TypeName.Charset
		}
		if n.Array {
			target += " array"
		}
		return "cast(" + storedTiDBExpressionSQL(n.Expr, charset, index) + " as " + target + ")"
	case *nodes.UnaryExpr:
		operand := storedTiDBExpressionSQL(n.Operand, charset, index)
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
		if index {
			return nodeToSQL(node)
		}
		return "(?)"
	}
}

func nodeToSQL(node nodes.ExprNode) string {
	return deparse.Deparse(node)
}

// defaultToSQL writes a column default the way the catalog keeps it: an
// expression default in its parentheses, so it stays apart from a literal, and
// TRUE and FALSE as the 1 and 0 TiDB stores for them.
func (c *Catalog) defaultToSQL(node nodes.ExprNode) string {
	switch n := node.(type) {
	case *nodes.ParenExpr:
		return "(" + storedTiDBExpressionSQL(n.Expr, c.defaultCharset, true) + ")"
	case *nodes.FuncCallExpr:
		return storedTiDBExpressionSQL(n, c.defaultCharset, true)
	case *nodes.BoolLit:
		if n.Value {
			return "1"
		}
		return "0"
	}
	return nodeToSQL(node)
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
		// BOOL takes the numeric modifiers (BooleanType FieldOpts, parser.y:12919,
		// v8.5.5) and keeps them: TiDB writes BOOL UNSIGNED as tinyint(1) unsigned.
		t := "tinyint(1)"
		if dt.Unsigned || dt.Zerofill {
			t += " unsigned"
		}
		return t
	case "numeric":
		name = "decimal"
	case "serial":
		return "bigint unsigned"
	case "geometrycollection":
		name = "geomcollection"
	case "float":
		// FLOAT(n) is a precision in bits, which picks the type and is not kept: TiDB v8.5.5
		// shows FLOAT(10) as float and FLOAT(30) as double. With a scale the first number is
		// no such thing, and FLOAT(30,0) stays float.
		if dt.Length > 24 && !dt.HasScale {
			name = "double"
		}
	}

	var buf strings.Builder
	buf.WriteString(name)

	// Integer display width, as TiDB's FieldType.CompactStr writes it: dropped,
	// except under ZEROFILL (with default widths per type: tinyint(3),
	// smallint(5), mediumint(8), int(10), bigint(20)) and for tinyint(1), which
	// keeps its width because connectors read it as a boolean.
	isIntType := isIntegerType(name)
	if isIntType {
		if name == "tinyint" && dt.Length == 1 && !dt.Zerofill {
			buf.WriteString("(1)")
		} else if dt.Zerofill {
			width := dt.Length
			if width == 0 {
				width = defaultIntDisplayWidth(name, true) // ZEROFILL implies UNSIGNED
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
		// TiDB writes every YEAR column as year(4).
		buf.WriteString("(4)")
	} else if name == "bit" && dt.Length == 0 {
		buf.WriteString("(1)")
	} else if (name == "char" || name == "binary") && dt.Length == 0 {
		// CHAR/BINARY with no length → MySQL shows char(1)/binary(1)
		buf.WriteString("(1)")
	} else if dt.Length > 0 && dt.Scale > 0 {
		fmt.Fprintf(&buf, "(%d,%d)", dt.Length, dt.Scale)
	} else if name == "float" || name == "double" {
		// Without a scale neither shows a precision: TiDB v8.5.5 shows DOUBLE(10,0) as double,
		// and DOUBLE(10) is a syntax error there, so printing it would not read back.
	} else if name == "decimal" && dt.Length > 0 {
		fmt.Fprintf(&buf, "(%d,0)", dt.Length)
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
	// YEAR takes UNSIGNED and ZEROFILL and keeps neither: TiDB's SHOW CREATE
	// TABLE writes `year(4)` for `year unsigned zerofill`, so a column that was
	// declared with them is not one that reports them.
	if name != "year" {
		// ZEROFILL implies UNSIGNED, which SHOW CREATE TABLE writes out.
		if dt.Unsigned || dt.Zerofill {
			buf.WriteString(" unsigned")
		}
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
// MySQL escapes a single quote in enum values by doubling it.
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
		Name:         tableName,
		Database:     db,
		Columns:      make([]*Column, 0, len(srcTbl.Columns)),
		colByName:    make(map[string]int),
		Indexes:      make([]*Index, 0, len(srcTbl.Indexes)),
		Constraints:  make([]*Constraint, 0, len(srcTbl.Constraints)),
		Engine:       srcTbl.Engine,
		Charset:      srcTbl.Charset,
		Collation:    srcTbl.Collation,
		Comment:      srcTbl.Comment,
		Compression:  srcTbl.Compression,
		RowFormat:    srcTbl.RowFormat,
		KeyBlockSize: srcTbl.KeyBlockSize,
		Temporary:    stmt.Temporary,
		// CREATE GLOBAL TEMPORARY TABLE ... LIKE t carries its own ON COMMIT,
		// so these come from the statement and not from the source table.
		GlobalTemporary: stmt.Global,
		OnCommit:        stmt.OnCommit,
	}

	// Copy columns.
	for i, srcCol := range srcTbl.Columns {
		col := &Column{
			Position:      srcCol.Position,
			Name:          srcCol.Name,
			DataType:      srcCol.DataType,
			ColumnType:    srcCol.ColumnType,
			Zerofill:      srcCol.Zerofill,
			Nullable:      srcCol.Nullable,
			AutoIncrement: srcCol.AutoIncrement,
			Charset:       srcCol.Charset,
			Collation:     srcCol.Collation,
			Comment:       srcCol.Comment,
			OnUpdate:      srcCol.OnUpdate,
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
			Vector:       srcIdx.Vector,
			IndexType:    srcIdx.IndexType,
			Visible:      srcIdx.Visible,
			Comment:      srcIdx.Comment,
			KeyBlockSize: srcIdx.KeyBlockSize,
			Global:       srcIdx.Global,
			// The copy keeps the source key's clustering: TiDB's LIKE of a
			// NONCLUSTERED primary key is NONCLUSTERED (read on 2026-11-06).
			addedByAlter: srcIdx.addedByAlter,
		}
		cols := make([]*IndexColumn, len(srcIdx.Columns))
		for i, sc := range srcIdx.Columns {
			cols[i] = &IndexColumn{
				Name:       sc.Name,
				Length:     sc.Length,
				Descending: sc.Descending,
				Expr:       sc.Expr,
				ExprNode:   sc.ExprNode,
				ExprType:   sc.ExprType,
			}
		}
		idx.Columns = cols
		tbl.Indexes = append(tbl.Indexes, idx)
	}

	// Copy constraints (skip FK — LIKE copies no foreign keys). TiDB renames
	// every CHECK it copies <new table>_chk_<n>, the named ones too: LIKE of a
	// table with chk_age gives dst_chk_1 (v8.5.5 with
	// tidb_enable_check_constraint ON, read on 2026-11-06).
	checks := 0
	for _, srcCon := range srcTbl.Constraints {
		if srcCon.Type == ConForeignKey {
			continue
		}
		name := srcCon.Name
		if srcCon.Type == ConCheck {
			checks++
			name = fmt.Sprintf("%s_chk_%d", tableName, checks)
		}
		con := &Constraint{
			Name:        name,
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
			Clustered:   srcCon.Clustered,
		}
		tbl.Constraints = append(tbl.Constraints, con)
	}

	db.Tables[key] = tbl
	return nil
}

// extractTTLParts parses a TiDB TTL expression string (stored verbatim on
// TableOption.Value by the parser) and splits it into (column, interval).
// Accepted shapes in TiDB v8.5:
//   - `<col> + INTERVAL <n> <unit>`
//   - `DATE_ADD(<col>, INTERVAL <n> <unit>)` (and sibling date functions)
//
// Any other shape returns a validation error — TiDB accepts other expressions
// but the catalog only tracks the column reference, so unknown shapes are not
// silently lost.
func extractTTLParts(val string) (string, string, error) {
	expr, err := tidbparser.ParseExpr(val)
	if err != nil {
		return "", "", fmt.Errorf("invalid TTL expression %q: %w", val, err)
	}
	switch e := expr.(type) {
	case *nodes.BinaryExpr:
		// Shape: <col> + INTERVAL <n> <unit>
		if e.Op == nodes.BinOpAdd {
			if col, ok := ttlColumnName(e.Left); ok {
				return col, deparse.Deparse(e.Right), nil
			}
		}
	case *nodes.FuncCallExpr:
		// Shape: DATE_ADD(<col>, INTERVAL <n> <unit>) / ADDDATE
		name := strings.ToUpper(e.Name)
		if (name == "DATE_ADD" || name == "ADDDATE") && len(e.Args) == 2 {
			if col, ok := ttlColumnName(e.Args[0]); ok {
				return col, deparse.Deparse(e.Args[1]), nil
			}
		}
	}
	return "", "", fmt.Errorf("TTL expression %q is not in a recognized form (expected `<col> + INTERVAL <n> <unit>` or `DATE_ADD(<col>, INTERVAL <n> <unit>)`)", val)
}

// ttlColumnName extracts a column name from an expression node if it is a
// plain column reference. Returns ("", false) otherwise.
func ttlColumnName(e nodes.ExprNode) (string, bool) {
	if c, ok := e.(*nodes.ColumnRef); ok {
		return c.Column, true
	}
	return "", false
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
		// No action given. TiDB writes nothing for it in SHOW CREATE TABLE,
		// where an explicit NO ACTION is written back.
		return ""
	}
}
