package catalog

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/deparse"
)

const generatedInvisiblePrimaryKeyColumnName = "my_row_id"

// createTable installs the table stmt defines. declared is true when stmt lists
// its keys in the order they were declared, as a CREATE TABLE statement does,
// and false for a loader's snapshot, which lists them in an order of its own.
// source is the text the statement was parsed from, which a CREATE TABLE ...
// SELECT body is judged against (isDualRef), or "" when there is none.
func (c *Catalog) createTable(stmt *nodes.CreateTableStmt, declared bool, source string) error {
	// MySQL checks a FULLTEXT key's parser, a column's display width and an
	// ENGINE while preparing the statement, so each answers ahead of the
	// database and the table it names: 1128 for a FULLTEXT key whose parser it
	// does not have, 1439 for `CREATE TABLE nosuchdb.a (a INT(256))`, and 1286
	// for an engine it does not ship (mysql:8.0.46).
	if declared {
		// A loader's table may carry a parser plugin another server installs.
		for _, con := range stmt.Constraints {
			if err := validateFulltextParser(con); err != nil {
				return err
			}
		}
	}
	for _, colDef := range stmt.Columns {
		if err := validateColumnDisplayWidth(colDef); err != nil {
			return err
		}
		// A bare NOW()-family DEFAULT on a column that is neither DATETIME nor
		// TIMESTAMP is 1067 while the statement is prepared, ahead of the name
		// collision and of the database the name sits in (mysql:8.0.46).
		if err := validateBareNowDefault(colDef); err != nil {
			return err
		}
		if err := validateColumnOnUpdate(colDef); err != nil {
			return err
		}
	}
	// A loader's statement (DefineTable) is a snapshot of a server that may
	// ship another engine, so only a CREATE TABLE statement is judged, and its
	// partitions are judged as ADD PARTITION's are.
	if declared {
		for _, opt := range stmt.Options {
			if opt == nil || !engineOptionName(opt.Name) {
				continue
			}
			if err := validateStorageEngine(opt.Value, c.noEngineSubstitution()); err != nil {
				return err
			}
		}
		if stmt.Partitions != nil {
			if err := validatePartitionDefinitions(stmt.Partitions.Partitions, c.noEngineSubstitution()); err != nil {
				return err
			}
		}
	}
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

	// A CREATE TABLE ... LIKE that names its own source is 1066, "Not unique
	// table/alias", before the table it would create is looked up: mysql:8.0.46
	// answers it whether the target exists (where 1050 would follow) or not
	// (where the source's 1146 would). Only the self case is decided here; every
	// other source is resolved after the target's collision checks, as before.
	if stmt.Like != nil {
		srcDB := c.GetDatabase(c.currentDB)
		if stmt.Like.Schema != "" {
			srcDB = c.GetDatabase(stmt.Like.Schema)
		}
		if srcDB == db && toLower(stmt.Like.Name) == key {
			return errNonUniqueTable(stmt.Like.Name)
		}
	}

	if declared {
		effEngine := c.sessionDefaultEngine(stmt.Temporary)
		effTS := ""
		effEncryption := ""
		effRowFormat := ""
		effKeyBlockSize := 0
		for _, opt := range stmt.Options {
			if opt == nil {
				continue
			}
			if engineOptionName(opt.Name) {
				if engineSubstituted(opt.Value) {
					effEngine = c.sessionDefaultEngine(stmt.Temporary)
				} else {
					effEngine = storedEngineName(opt.Value)
				}
			} else if strings.EqualFold(opt.Name, "tablespace") {
				effTS = opt.Value
			} else if strings.EqualFold(opt.Name, "encryption") {
				effEncryption = opt.Value
			} else if strings.EqualFold(opt.Name, "row_format") {
				effRowFormat = opt.Value
			} else if strings.EqualFold(opt.Name, "key_block_size") {
				if n, err := strconv.Atoi(strings.TrimSpace(opt.Value)); err == nil {
					effKeyBlockSize = n
				}
			}
		}
		if isInnoDBEngine(effEngine) {
			if err := c.validateInnoDBCompressedTableOption(tableName, effTS, effRowFormat, effKeyBlockSize, false); err != nil {
				return err
			}
			if effTS != "" {
				if err := c.validateInnoDBTablespace(tableName, stmt.Temporary, effTS); err != nil {
					return err
				}
				if err := c.validateInnoDBTableBlockSize(effTS, effRowFormat, effKeyBlockSize); err != nil {
					return err
				}
				if err := c.validateInnoDBTablespaceEncryption(effTS, effEncryption); err != nil {
					return err
				}
			}
		}
		if stmt.Partitions != nil {
			if err := c.validatePartitionTablespaces(tableName, stmt.Temporary, effEngine, stmt.Partitions.Partitions); err != nil {
				return err
			}
		}
	}

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

	// CREATE TABLE ... AS SELECT (CTAS). MySQL reads the body after the table's
	// own name, which the checks above answer 1050 for: over a table that is not
	// there, `CREATE TABLE z SELECT nosuch` is 1054 `Unknown column 'nosuch' in
	// 'field list'`, and a star with no table to expand is 1096 (`SELECT *` and
	// `SELECT * FROM DUAL` alike). Both are a view body's refusals, so the body
	// is read as one (analyzeViewBody), whether or not the statement declares its
	// columns: `CREATE TABLE z (a INT) AS SELECT nosuch` is the same 1054. The
	// catalog does not model the table a body makes without declared columns, so
	// one that reads installs nothing.
	if stmt.Select != nil {
		if declared {
			if _, err := c.analyzeViewBody(stmt.Select, source); err != nil {
				return err
			}
		}
		if len(stmt.Columns) == 0 {
			return nil
		}
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
	// A table without ENGINE takes the session's default engine,
	// default_tmp_storage_engine for a TEMPORARY one: under
	// default_storage_engine = MyISAM, mysql:8.0.46 stores ENGINE=MyISAM.
	if declared {
		tbl.Engine = c.sessionDefaultEngine(stmt.Temporary)
	}

	// Apply table options.
	tblCharsetExplicit := false
	tblCollationExplicit := false
	for _, opt := range stmt.Options {
		switch toLower(opt.Name) {
		case "engine":
			// A loader's snapshot keeps the engine its server reported, which the
			// pinned one may not ship; an authored statement keeps what the
			// pinned server would store.
			// An engine MySQL substitutes is the session's default.
			if declared && engineSubstituted(opt.Value) {
				tbl.Engine = c.sessionDefaultEngine(stmt.Temporary)
			} else if declared {
				tbl.Engine = storedEngineName(opt.Value)
			} else {
				tbl.Engine = opt.Value
			}
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
		case "checksum", "table_checksum":
			tbl.Checksum = normalizedChecksum(opt.Value)
		case "delay_key_write":
			tbl.DelayKeyWrite = opt.Value
		case "secondary_engine":
			if strings.EqualFold(opt.Value, "NULL") {
				tbl.SecondaryEngine = ""
			} else {
				tbl.SecondaryEngine = opt.Value
			}
		}
	}
	hasAutoExtend := false
	rawAutoExtend := ""
	for _, opt := range stmt.Options {
		if opt != nil && strings.EqualFold(opt.Name, "autoextend_size") {
			hasAutoExtend = true
			rawAutoExtend = opt.Value
		}
	}
	if hasAutoExtend && isInnoDBEngine(tbl.Engine) {
		if declared && !c.loadingSDL {
			aeBytes, err := c.parseTableAutoextendSize(rawAutoExtend, stmt.Temporary, tbl.Tablespace)
			if err != nil {
				return err
			}
			tbl.AutoExtendSize = aeBytes
		} else if aeBytes, ok := parseMySQLByteSize(rawAutoExtend); ok {
			tbl.AutoExtendSize = aeBytes
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
	// The UNIQUE keys in the order the statement declares them, for the
	// implicit primary key.
	var uniqueKeys []declaredUniqueKey

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
		if err := validateColumnDefSemantics(colDef, c.isIgnoreSpaceSQLMode()); err != nil {
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
				col.CharsetExplicit = true
			}
			if colDef.TypeName.Collate != "" {
				col.Collation = colDef.TypeName.Collate
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
			c.setColumnDefaultFromExpr(col, colDef.DefaultValue)
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
					c.setColumnDefaultFromExpr(col, cc.Expr)
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
				if autoIncrementColumn != "" && !strings.EqualFold(autoIncrementColumn, colDef.Name) {
					return errIncorrectAutoIncrementDefinition()
				}
				autoIncrementColumn = colDef.Name
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
			case nodes.ColConstrNotSecondary:
				col.NotSecondary = true
			case nodes.ColConstrCollate:
				// Collation specified via constraint.
				collName := cc.Name
				if collName == "" && cc.Expr != nil {
					if s, ok := cc.Expr.(*nodes.StringLit); ok {
						collName = s.Value
					}
				}
				if collName != "" {
					col.Collation = collName
					col.CollationExplicit = true
					if !col.CharsetExplicit {
						if cs := charsetForCollation(collName); cs != "" {
							col.Charset = normalizeCharsetName(cs)
						}
					}
					if strings.EqualFold(col.Charset, "binary") && isStringType(col.DataType) && !isEnumSetType(col.DataType) && colDef.TypeName != nil {
						col = convertToBinaryType(col, colDef.TypeName)
					}
				}
			}
		}

		if col.CharsetExplicit && col.CollationExplicit && !charsetMatchesCollation(col.Charset, col.Collation) {
			return &Error{Code: 1253, SQLState: "42000", Message: "COLLATION is not valid for CHARACTER SET"}
		}
		if declared {
			if colDef.TypeName != nil {
				if err := validateEnumSetDuplicates(col.Name, col.DataType, col.Collation, col.Charset, colDef.TypeName.EnumValues, c.isStrictSQLMode()); err != nil {
					return err
				}
			}
			if err := c.validateColumnLiteralDefault(col, columnDefEffectiveDefaultExpr(colDef)); err != nil {
				return err
			}
		} else if effDef := columnDefEffectiveDefaultExpr(colDef); effDef != nil {
			c.normalizeColumnStoredDefault(col, effDef, true, false)
		}

		// SERIAL implies UNIQUE KEY — add after the column is fully configured.
		if isSerial {
			idxName := allocIndexName(tbl, colDef.Name)
			serialCols := []*IndexColumn{{Name: colDef.Name}}
			if declared {
				if err := c.validateIndexColumns(tbl, serialCols, false, false, true); err != nil {
					return err
				}
			}
			serialIdx := &Index{
				Name:      idxName,
				Table:     tbl,
				Columns:   serialCols,
				Unique:    true,
				IndexType: "",
				Visible:   true,
			}
			tbl.Indexes = append(tbl.Indexes, serialIdx)
			uniqueKeys = append(uniqueKeys, declaredUniqueKey{serialIdx, colDef.Loc.Start})
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
				pkCols := []*IndexColumn{{Name: colDef.Name}}
				if declared {
					if err := c.validateIndexColumns(tbl, pkCols, false, false, true); err != nil {
						return err
					}
				}
				tbl.Indexes = append(tbl.Indexes, &Index{
					Name:      "PRIMARY",
					Table:     tbl,
					Columns:   pkCols,
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
				uqCols := []*IndexColumn{{Name: colDef.Name}}
				if declared {
					if err := c.validateIndexColumns(tbl, uqCols, false, false, true); err != nil {
						return err
					}
				}
				uqIdx := &Index{
					Name:      idxName,
					Table:     tbl,
					Columns:   uqCols,
					Unique:    true,
					IndexType: "",
					Visible:   true,
				}
				tbl.Indexes = append(tbl.Indexes, uqIdx)
				uniqueKeys = append(uniqueKeys, declaredUniqueKey{uqIdx, cc.Loc.Start})
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
	//
	// A loader's statement (DefineTable) keeps a key part it could not read as a
	// column (metadataIndexConstraint), so which columns a key names is judged
	// only for a CREATE TABLE statement.
	checkKeyColumns := func(con *nodes.Constraint) error {
		if !declared {
			return nil
		}
		return validateKeyColumnsExist(tbl, nil, con)
	}
	for _, con := range stmt.Constraints {
		cols := extractColumnNames(con)

		switch con.Type {
		case nodes.ConstrPrimaryKey:
			if declared && constraintHasFunctionalKeyPart(con) {
				return errFunctionalIndexPrimaryKey()
			}
			if hasPK {
				return errMultiplePriKey()
			}
			if err := checkKeyColumns(con); err != nil {
				return err
			}
			if err := validatePrimaryKeyColumns(tbl, cols); err != nil {
				return err
			}
			idxCols := buildIndexColumns(con)
			if declared {
				if err := validateUniqueOrPrimaryNotSpatial(tbl, idxCols, resolveConstraintIndexType(con)); err != nil {
					return err
				}
				if err := c.validateIndexColumns(tbl, idxCols, false, false, true); err != nil {
					return err
				}
				if name := duplicateKeyColumnName(con); name != "" {
					return errDupColumn(name)
				}
			}
			hasPK = true
			// Mark PK columns as NOT NULL.
			for _, colName := range cols {
				c := tbl.GetColumn(colName)
				if c != nil {
					c.Nullable = false
				}
			}
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
			if declared {
				if err := validateNoColumnFunctionalParts(con.IndexColumns); err != nil {
					return err
				}
			}
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
			if err := checkKeyColumns(con); err != nil {
				return err
			}
			if declared {
				if err := validateUniqueOrPrimaryNotSpatial(tbl, idxCols, resolveConstraintIndexType(con)); err != nil {
					return err
				}
			}
			if err := c.validateIndexColumns(tbl, idxCols, false, false, true); err != nil {
				return err
			}
			if declared {
				if name := duplicateKeyColumnName(con); name != "" {
					return errDupColumn(name)
				}
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
			coerceInnoDBHashIndex(tbl, uqIdx)
			if err := synthesizeFunctionalIndexColumnsForStmt(tbl, uqIdx, nil, declared, c.isIgnoreSpaceSQLMode()); err != nil {
				return err
			}
			if hasFunctionalIndexColumns(uqIdx.Columns) {
				if err := c.validateIndexTotalKeyWidth(tbl, uqIdx.Columns); err != nil {
					return err
				}
			}
			tbl.Indexes = append(tbl.Indexes, uqIdx)
			uniqueKeys = append(uniqueKeys, declaredUniqueKey{uqIdx, con.Loc.Start})
			tbl.Constraints = append(tbl.Constraints, &Constraint{
				Name:      idxName,
				Type:      ConUniqueKey,
				Table:     tbl,
				Columns:   cols,
				IndexName: idxName,
			})

		case nodes.ConstrForeignKey:
			if err := checkKeyColumns(con); err != nil {
				return err
			}
			fkIdxCols := buildIndexColumns(con)
			if declared {
				if err := c.validateIndexColumns(tbl, fkIdxCols, false, false, false); err != nil {
					return err
				}
				if name := duplicateKeyColumnName(con); name != "" {
					return errDupColumn(name)
				}
				if len(con.RefColumns) == 0 || len(cols) != len(con.RefColumns) {
					return errWrongFkDef(fkErrorDisplayName(con), "Key reference and table reference don't match")
				}
			}
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
			pendingFKs = append(pendingFKs, pendingFK{conName: fkIndexName(con), cols: cols, idxCols: fkIdxCols})

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
			if declared {
				if err := validateNoColumnFunctionalParts(con.IndexColumns); err != nil {
					return err
				}
			}
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
			if err := checkKeyColumns(con); err != nil {
				return err
			}
			spatial := isSpatialIndexType(resolveConstraintIndexType(con))
			if err := c.validateIndexColumns(tbl, idxCols, false, spatial, false); err != nil {
				return err
			}
			if declared {
				if name := duplicateKeyColumnName(con); name != "" {
					return errDupColumn(name)
				}
			}
			keyIdx := &Index{
				Name:      idxName,
				Table:     tbl,
				Columns:   idxCols,
				Spatial:   spatial,
				IndexType: storedIndexType(resolveConstraintIndexType(con)),
				Visible:   true,
			}
			applyIndexOptions(keyIdx, con.IndexOptions)
			coerceInnoDBHashIndex(tbl, keyIdx)
			if err := synthesizeFunctionalIndexColumnsForStmt(tbl, keyIdx, nil, declared, c.isIgnoreSpaceSQLMode()); err != nil {
				return err
			}
			if !spatial && hasFunctionalIndexColumns(keyIdx.Columns) {
				if err := c.validateIndexTotalKeyWidth(tbl, keyIdx.Columns); err != nil {
					return err
				}
			}
			tbl.Indexes = append(tbl.Indexes, keyIdx)

		case nodes.ConstrFulltextIndex:
			if declared {
				if err := validateNoColumnFunctionalParts(con.IndexColumns); err != nil {
					return err
				}
			}
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
			if err := checkKeyColumns(con); err != nil {
				return err
			}
			idxCols := buildIndexColumns(con)
			if err := c.validateIndexColumns(tbl, idxCols, true, false, false); err != nil {
				return err
			}
			if err := validateFulltextSpatialFunctionalParts(tbl, nil, declared, idxName, idxCols, false, c.isIgnoreSpaceSQLMode()); err != nil {
				return err
			}
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
			if declared {
				if err := validateNoColumnFunctionalParts(con.IndexColumns); err != nil {
					return err
				}
			}
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
			if err := checkKeyColumns(con); err != nil {
				return err
			}
			idxCols := buildIndexColumns(con)
			if err := c.validateIndexColumns(tbl, idxCols, false, true, false); err != nil {
				return err
			}
			if err := validateFulltextSpatialFunctionalParts(tbl, nil, declared, idxName, idxCols, true, c.isIgnoreSpaceSQLMode()); err != nil {
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

	// PERFORMANCE_SCHEMA holds the performance schema's own tables, which MySQL
	// refuses with 1683 once the table's own definition is accepted
	// (mysql:8.0.46).
	if declared {
		for _, opt := range stmt.Options {
			if opt == nil || !engineOptionName(opt.Name) {
				continue
			}
			if strings.EqualFold(strings.TrimSpace(opt.Value), performanceSchemaEngine) {
				return errInvalidPerformanceSchema()
			}
		}
	}

	if c.generateGIPK {
		if tbl.GetColumn(generatedInvisiblePrimaryKeyColumnName) != nil {
			return errGIPKColumnNameReserved()
		}
		if !hasPK {
			// MySQL 8.0 generates no key for a partitioned table and refuses it, so the
			// partition-key check below never sees a my_row_id key the engine would not add.
			if stmt.Partitions != nil {
				return errGIPKPartitioned()
			}
			addGeneratedInvisiblePrimaryKey(tbl)
			hasPK = true
		}
	}

	// Without a primary key MySQL 8.0 refuses a table whose implicit one is
	// invisible. A generated invisible primary key is a primary key, so the
	// check comes after it. A snapshot does not keep the declaration order the
	// check needs, and the table it describes passed the check on the server.
	if !hasPK && declared {
		if err := checkImplicitPrimaryKeyVisible(tbl, uniqueKeys); err != nil {
			return err
		}
	}

	// Process deferred FK backing indexes now that all explicit indexes are in place.
	for _, fk := range pendingFKs {
		ensureFKBackingIndex(tbl, fk.conName, fk.cols, fk.idxCols)
	}

	// A table with no visible column is 4028, not only one with no column at
	// all: `CREATE TABLE q (a INT INVISIBLE)` is refused, and so it is with
	// sql_generate_invisible_primary_key on, where the only column is the
	// invisible my_row_id. MySQL answers it after the name collision (1050) and
	// after a key's own columns are resolved (1072), and before the partition
	// clause, which alone cannot make a column: `CREATE TABLE q (KEY (x))` is
	// 1072 and `CREATE TABLE q PARTITION BY KEY()` 4028 (mysql:8.0.46).
	// A CTAS is never this error's: MySQL appends the query's own columns to the
	// declared ones as visible columns, so `CREATE TABLE t2 (col2 INT INVISIBLE)
	// AS SELECT col1, col2 FROM t1` is a table it takes (mysql:8.0.46).
	if stmt.Select == nil && !tableHasVisibleColumn(tbl) {
		return errTableMustHaveColumns()
	}
	if declared {
		if !c.loadingSDL {
			if err := validateAutoIncrementKeyed(tbl); err != nil {
				return err
			}
		}
		if err := c.validateTableDefaultExpressions(tbl); err != nil {
			return err
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
		if declared && !c.loadingSDL && isInnoDBEngine(tbl.Engine) {
			if err := validateInnoDBPartitionedSharedTablespace(tbl.Tablespace); err != nil {
				return err
			}
		}
		tbl.Partitioning = buildPartitionInfo(tbl, stmt.Partitions, declared)
	}

	// Phase 3: analyze DEFAULT, GENERATED, and CHECK expressions now that all
	// columns are present in the table.
	c.analyzeTableExpressions(tbl, stmt)

	if !declared {
		c.registerMetadataTableTablespaces(tbl)
	}

	db.Tables[key] = tbl
	return nil
}

// tableHasVisibleColumn reports whether the table keeps a column a SELECT * and
// a reader may see. my_row_id is invisible, so a table the generated invisible
// primary key alone fills has none, and neither has one whose only other column
// is a functional index's hidden system column.
func tableHasVisibleColumn(tbl *Table) bool {
	for _, col := range tbl.Columns {
		if col != nil && !col.Invisible && col.Hidden == ColumnHiddenNone {
			return true
		}
	}
	return false
}

func errGIPKPartitioned() error {
	return &Error{
		Code:     1235,
		SQLState: "42000",
		Message:  "This version of MySQL doesn't yet support 'generating invisible primary key for the partitioned tables'",
	}
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
func buildPartitionInfo(tbl *Table, pc *nodes.PartitionClause, declared bool) *PartitionInfo {
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
		// Options. A loader's snapshot keeps the engine its server reported,
		// which the pinned one may not ship.
		for _, opt := range pd.Options {
			switch toLower(opt.Name) {
			case "engine":
				pdi.Engine = partitionEngineName(opt.Value, declared)
			case "tablespace":
				pdi.Tablespace = unquoteTableOption(opt.Value)
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
				Name:       spd.Name,
				Engine:     pdi.Engine,
				Tablespace: pdi.Tablespace,
				Comment:    pdi.Comment,
			}
			for _, opt := range spd.Options {
				switch toLower(opt.Name) {
				case "engine":
					spdi.Engine = partitionEngineName(opt.Value, declared)
				case "tablespace":
					spdi.Tablespace = unquoteTableOption(opt.Value)
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

	// Auto-generate subpartition definitions when SUBPARTITIONS N is specified
	// without explicit subpartition definitions.
	// MySQL naming convention: <partition_name>sp0, <partition_name>sp1, ...
	if pi.SubType != "" && pi.UseDefaultSubpartitions && pi.NumSubParts > 0 {
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

func (c *Catalog) resolveNormalizedStringDefaultLiteral(col *Column, sl *nodes.StringLit) (string, bool, bool) {
	litVal := sl.Value
	decoded := false
	clientCS := "utf8mb4"
	connCS := "utf8mb4"
	if c != nil {
		clientCS = c.effectiveClientCharset()
		connCS = c.effectiveConnCharset()
	}
	hasSessionCS := clientCS != "utf8mb4" || connCS != "utf8mb4"
	if sl.Charset != "" || sl.National || hasSessionCS {
		if dv, decOK := c.decodeStringLiteralDefault(sl, col.Charset); decOK {
			litVal = dv
			decoded = true
		}
	}
	colIsBinary := strings.EqualFold(col.DataType, "binary") ||
		strings.EqualFold(col.DataType, "varbinary") ||
		strings.EqualFold(col.Charset, "binary")
	if colIsBinary && sl.Charset == "" && !sl.National && clientCS != connCS {
		connBytes := encodeStringToCharsetBytes(litVal, connCS)
		hasHighOrCtrl := false
		for _, b := range connBytes {
			if b >= 0x80 || b < 0x20 {
				hasHighOrCtrl = true
				break
			}
		}
		if hasHighOrCtrl {
			return fmt.Sprintf("0x%X", connBytes), true, true
		}
		return string(connBytes), false, true
	}
	litIsBinary := strings.EqualFold(strings.TrimSpace(sl.Charset), "binary") ||
		(colIsBinary && (sl.Charset != "" || sl.National || clientCS == connCS))
	normVal, changed := normalizeCharDefaultStringLiteral(col, litVal)
	return normVal, false, changed || (decoded && !litIsBinary)
}

func (c *Catalog) setColumnDefaultFromExpr(col *Column, expr nodes.ExprNode) {
	if pe, ok := expr.(*nodes.ParenExpr); ok {
		if stripped := stripColumnQualifiers(pe.Expr); stripped != pe.Expr {
			cp := *pe
			cp.Expr = stripped
			expr = &cp
		}
	}
	if sl, ok := expr.(*nodes.StringLit); ok {
		if normVal, isHex, shouldRewrite := c.resolveNormalizedStringDefaultLiteral(col, sl); shouldRewrite {
			if isHex {
				v := normVal
				col.Default = &v
				col.DefaultKind = ColumnDefaultConstant
				col.DefaultDropped = false
				col.DefaultIsExpr = false
				col.DefaultExprNode = nil
				return
			}
			cp := *sl
			cp.Charset = ""
			cp.National = false
			cp.Value = normVal
			expr = &cp
		}
	} else if coerced, ok := coerceNumericDefaultForCharColumn(col, expr); ok {
		expr = &nodes.StringLit{Value: coerced}
	}
	val, kind := columnDefaultValue(expr)
	if _, isParen := expr.(*nodes.ParenExpr); !isParen {
		if coerced, ok := coerceNumericColumnDefault(col, expr, true, false, utcResolvedTimeZone()); ok {
			val = coerced
		}
	}
	col.Default = &val
	col.DefaultKind = kind
	col.DefaultDropped = false
	if pe, ok := expr.(*nodes.ParenExpr); ok {
		col.DefaultIsExpr = true
		col.DefaultExprNode = pe.Expr
	} else {
		col.DefaultIsExpr = false
		col.DefaultExprNode = nil
	}
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

// validatePartitionClause judges a partitioning a CREATE TABLE writes: the clause
// itself, and then every unique key of the table it leaves covering it.
func validatePartitionClause(tbl *Table, pc *nodes.PartitionClause) error {
	if err := validatePartitionClauseShape(tbl, pc); err != nil {
		return err
	}
	return validateUniqueKeysCoverPartitionColumns(tbl, partitionKeyColumns(tbl, pc))
}

// validatePartitionClauseShape runs the checks the clause makes on its own, without
// the unique keys that must cover it: what it defines (1492), what it mixes (1500),
// the columns it names (1652, 1488, 1054) and its bounds (1481, 1493, 1491).
func validatePartitionClauseShape(tbl *Table, pc *nodes.PartitionClause) error {
	if (pc.Type == nodes.PartitionRange || pc.Type == nodes.PartitionList) && len(pc.Partitions) == 0 {
		return errPartitionsMustBeDefined()
	}
	// A top-level HASH or KEY partitioning takes no SUBPARTITION BY: MySQL mixes
	// subpartitioning with RANGE and LIST only, and answers 1500 before it reads a
	// column either function names (mysql:8.0.46).
	if (pc.Type == nodes.PartitionHash || pc.Type == nodes.PartitionKey) && partitionClauseHasSubpartitions(pc) {
		return errPartitionSubpartMix()
	}
	// Every column the clause names must be a column of the table. KEY's name list
	// answers 1488, a partitioning expression 1054, and a name the list repeats 1652,
	// which the server answers first: `PARTITION BY KEY(zzz, zzz)` is 1652 and
	// `PARTITION BY KEY(zzz)` on a table keyed by another column is 1488 rather than
	// the 1503 its unique keys would raise (mysql:8.0.46).
	if err := validatePartitionClauseColumns(tbl, pc); err != nil {
		return err
	}
	if pc.Type == nodes.PartitionRange && len(pc.Columns) > 0 {
		// RANGE COLUMNS takes MAXVALUE in any partition: what it requires is
		// bounds that strictly increase, compared column by column, and it
		// answers 1493 otherwise. Holding it to RANGE's 1481 refused
		// "(5, MAXVALUE), (6, 1)", which the engine runs.
		kinds := rangeColumnsBoundKinds(tbl, pc.Columns)
		for i := 1; i < len(pc.Partitions); i++ {
			prev, cur := pc.Partitions[i-1].Values, pc.Partitions[i].Values
			if prev != nil && cur != nil && rangeColumnsNotIncreasing(prev, cur, kinds) {
				return &Error{Code: 1493, SQLState: "HY000", Message: "VALUES LESS THAN value must be strictly increasing for each partition"}
			}
		}
	} else if pc.Type == nodes.PartitionRange {
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
	if err := validateSubpartitionDefinitions(pc); err != nil {
		return err
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
	return nil
}

// partitionClauseHasSubpartitions reports whether the clause writes a SUBPARTITION
// BY. PartitionType's zero value is RANGE, which no SUBPARTITION BY takes, so a
// subpartitioning type is one of the other two.
func partitionClauseHasSubpartitions(pc *nodes.PartitionClause) bool {
	return pc.SubPartType == nodes.PartitionHash || pc.SubPartType == nodes.PartitionKey ||
		pc.SubPartExpr != nil || len(pc.SubPartColumns) > 0
}

// validatePartitionClauseColumns refuses a partition function that names a column
// the table does not hold, and a name one list repeats.
func validatePartitionClauseColumns(tbl *Table, pc *nodes.PartitionClause) error {
	if err := validatePartitionFields(tbl, pc.Type, pc.Columns, pc.Expr); err != nil {
		return err
	}
	return validatePartitionFields(tbl, pc.SubPartType, pc.SubPartColumns, pc.SubPartExpr)
}

// validatePartitionFields checks one partitioning function's columns against the
// table: a KEY list's names, which MySQL answers 1652 and then 1488 for, and the
// columns an expression reads, which it answers 1054 for. A COLUMNS list is read
// as the names it is compared by.
func validatePartitionFields(tbl *Table, typ nodes.PartitionType, columns []string, expr nodes.ExprNode) error {
	seen := map[string]bool{}
	for _, name := range columns {
		if seen[toLower(name)] {
			return errDupPartitionField(name)
		}
		seen[toLower(name)] = true
	}
	for _, name := range columns {
		if tbl.GetColumn(name) == nil {
			if typ == nodes.PartitionKey {
				return errPartitionFieldNotFound()
			}
			return errNoSuchColumn(name, "partition function")
		}
	}
	for _, name := range exprColumnNames(expr) {
		if tbl.GetColumn(name) == nil {
			return errNoSuchColumn(name, "partition function")
		}
	}
	return nil
}

// rangeColumnsNotIncreasing reports whether the RANGE COLUMNS bound cur is
// certainly not greater than prev, which the engine refuses with 1493. It
// compares column by column, as MySQL and MariaDB do: MAXVALUE is greater than
// any value, and a MAXVALUE in both ends the comparison as equal, so
// "(MAXVALUE, 5), (MAXVALUE, 10)" is refused. Below MAXVALUE, two bounds are
// compared only where their column's type settles the order (kinds); any other
// pair leaves the answer open, and nothing is reported.
func rangeColumnsNotIncreasing(prev, cur nodes.Node, kinds []rangeBoundKind) bool {
	ps, cs := partitionBoundItems(prev), partitionBoundItems(cur)
	if len(ps) != len(cs) || len(ps) != len(kinds) {
		return false
	}
	for i := range ps {
		pm, cm := isPartitionMaxvalue(ps[i]), isPartitionMaxvalue(cs[i])
		switch {
		case pm:
			return true
		case cm:
			return false
		}
		order, known := comparePartitionBounds(ps[i], cs[i], kinds[i])
		if !known {
			return false
		}
		if order != 0 {
			return order > 0
		}
	}
	return true
}

// rangeBoundKind is how a RANGE COLUMNS column's literal bounds compare where
// the order is certain.
type rangeBoundKind int

const (
	rangeBoundOpen     rangeBoundKind = iota // nothing is decided
	rangeBoundInteger                        // integer literals, by value
	rangeBoundString                         // string literals, by byte where the collation keeps it
	rangeBoundTemporal                       // date and time literals of one fixed, zero-padded shape
)

// rangeColumnsBoundKinds reads each partitioning column's kind off the table.
// A string column compares only under a collation with no tailoring of ASCII
// letters and digits: a language collation such as utf8mb4_spanish2_ci sorts
// "ch" after "cz", which byte order cannot see.
func rangeColumnsBoundKinds(tbl *Table, cols []string) []rangeBoundKind {
	kinds := make([]rangeBoundKind, len(cols))
	for i, name := range cols {
		col := tbl.GetColumn(name)
		if col == nil {
			continue
		}
		switch strings.ToLower(col.DataType) {
		case "tinyint", "smallint", "mediumint", "int", "integer", "bigint":
			kinds[i] = rangeBoundInteger
		case "binary", "varbinary":
			kinds[i] = rangeBoundString
		case "char", "varchar":
			if plainCollation(col.Collation) {
				kinds[i] = rangeBoundString
			}
		case "date", "datetime":
			kinds[i] = rangeBoundTemporal
		}
	}
	return kinds
}

// plainCollation reports whether a collation orders ASCII digits and lowercase
// letters as their bytes do: binary ones and the untailored defaults.
func plainCollation(c string) bool {
	c = strings.ToLower(c)
	if c == "" || c == "binary" || strings.HasSuffix(c, "_bin") {
		return true
	}
	switch c {
	case "utf8mb4_0900_ai_ci", "utf8mb4_0900_as_ci", "utf8mb4_0900_as_cs", "utf8mb4_general_ci",
		"utf8mb4_unicode_ci", "utf8mb4_unicode_520_ci", "utf8mb3_general_ci", "utf8mb3_unicode_ci",
		"utf8_general_ci", "utf8_unicode_ci", "latin1_swedish_ci", "latin1_general_ci", "latin1_general_cs",
		"ascii_general_ci":
		return true
	}
	return false
}

func partitionBoundItems(n nodes.Node) []nodes.Node {
	if l, ok := n.(*nodes.List); ok {
		return l.Items
	}
	return []nodes.Node{n}
}

func isPartitionMaxvalue(n nodes.Node) bool {
	s, ok := n.(*nodes.String)
	return ok && strings.EqualFold(s.Str, "MAXVALUE")
}

// comparePartitionBounds orders two bounds of a column of the given kind, a
// against b, when the order is certain: -1, 0 or 1, and whether it is known.
func comparePartitionBounds(a, b nodes.Node, kind rangeBoundKind) (int, bool) {
	switch kind {
	case rangeBoundInteger:
		x, ok1 := partitionBoundInt(a)
		y, ok2 := partitionBoundInt(b)
		if !ok1 || !ok2 {
			return 0, false
		}
		switch {
		case x < y:
			return -1, true
		case x > y:
			return 1, true
		}
		return 0, true
	case rangeBoundString, rangeBoundTemporal:
		x, ok1 := a.(*nodes.StringLit)
		y, ok2 := b.(*nodes.StringLit)
		if !ok1 || !ok2 {
			return 0, false
		}
		if x.Value == y.Value {
			return 0, true
		}
		if kind == rangeBoundTemporal {
			return compareTemporalLiterals(x.Value, y.Value)
		}
		for i := 0; i < len(x.Value) && i < len(y.Value); i++ {
			c, d := x.Value[i], y.Value[i]
			if c == d {
				continue
			}
			digits := c >= '0' && c <= '9' && d >= '0' && d <= '9'
			lower := c >= 'a' && c <= 'z' && d >= 'a' && d <= 'z'
			if !digits && !lower {
				return 0, false
			}
			if c < d {
				return -1, true
			}
			return 1, true
		}
	}
	return 0, false
}

// compareTemporalLiterals orders two date or time literals written in one
// shape: equal length, digits in the same places and the same separators
// between them, as '2020-01-31' and '2021-12-01' are, and a four-digit year
// first. Zero-padded fields of one width then sort as their bytes do;
// '2020-9-1' against '2020-10-01' does not, and neither does a two-digit year,
// which reads 70-99 as 1970-1999 and 00-69 as 2000-2069, so both are left
// open.
func compareTemporalLiterals(x, y string) (int, bool) {
	if len(x) != len(y) || !fourDigitYear(x) {
		return 0, false
	}
	for i := 0; i < len(x); i++ {
		xd, yd := x[i] >= '0' && x[i] <= '9', y[i] >= '0' && y[i] <= '9'
		if xd != yd || (!xd && x[i] != y[i]) {
			return 0, false
		}
	}
	if x < y {
		return -1, true
	}
	return 1, true
}

// fourDigitYear reports whether a date literal leads with a four-digit year
// and a separator after it.
func fourDigitYear(s string) bool {
	if len(s) < 5 {
		return false
	}
	for i := 0; i < 4; i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return s[4] < '0' || s[4] > '9'
}

func partitionBoundInt(n nodes.Node) (int64, bool) {
	if e, ok := n.(nodes.ExprNode); ok {
		if v, ok := foldConstIntExpr(e); ok {
			// The lexer clamps a literal past int64 to MaxInt64, and a minus
			// before it folds to -MaxInt64, so a value at either end may not
			// be the one written: -9223372036854775808 reads as
			// -9223372036854775807.
			if v >= math.MaxInt64 || v <= -math.MaxInt64 {
				return 0, false
			}
			return v, true
		}
	}
	return 0, false
}

func validateSubpartitionDefinitions(pc *nodes.PartitionClause) error {
	if pc.SubPartType == 0 && pc.SubPartExpr == nil && len(pc.SubPartColumns) == 0 {
		return nil
	}
	explicitCount := 0
	seenExplicit := false
	for i, pd := range pc.Partitions {
		count := len(pd.SubPartitions)
		if count == 0 {
			if seenExplicit {
				return errWrongNumberOfSubpartitions()
			}
			continue
		}
		if !seenExplicit {
			if i > 0 {
				return errWrongNumberOfSubpartitions()
			}
			seenExplicit = true
			explicitCount = count
		}
		if count != explicitCount {
			return errWrongNumberOfSubpartitions()
		}
		if pc.NumSubParts > 0 && count != pc.NumSubParts {
			return errWrongNumberOfSubpartitions()
		}
	}
	return nil
}

func errWrongNumberOfSubpartitions() error {
	return &Error{
		Code:     1064,
		SQLState: "42000",
		Message:  "Wrong number of subpartitions defined, mismatch with previous setting",
	}
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
					subParts[j] = partitionBoundSQL(sub)
				}
				parts[i] = "(" + strings.Join(subParts, ",") + ")"
			} else {
				parts[i] = partitionBoundSQL(item)
			}
		}
		return strings.Join(parts, ",")
	case nodes.ExprNode:
		return partitionValueItemSQL(n)
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
		return partitionValueItemSQL(n)
	}
	return ""
}

// partitionValueItemSQL renders one scalar partition bound, constant-folding it to an integer
// literal when the expression is a constant integer expression (integer arithmetic, or a supported
// deterministic function on a literal — see foldConstIntExpr). On any non-foldable expression it
// falls back to the verbatim deparse, so a literal bound (the common case) is rendered unchanged and
// an unsupported expression is preserved (and flagged as a documented round-trip limitation).
func partitionValueItemSQL(n nodes.ExprNode) string {
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
				expr := stripColumnQualifiers(ic.Expr)
				idxCol.Expr = nodeToSQL(expr)
				idxCol.ExprNode = expr
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

// fkIndexName returns the name of the index a foreign key adds when none
// serves it: the CONSTRAINT symbol, else the name after FOREIGN KEY, else
// "" for one named after its first column. The foreign key itself takes
// the symbol, or a generated name without one, never the name after
// FOREIGN KEY (verified on 8.0.46).
func fkIndexName(con *nodes.Constraint) string {
	if con.Name != "" {
		return con.Name
	}
	return con.IndexName
}

func fkErrorDisplayName(con *nodes.Constraint) string {
	if con.Name != "" {
		return con.Name
	}
	if con.IndexName != "" {
		return con.IndexName
	}
	return "foreign key without name"
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

func (c *Catalog) validateIndexColumns(tbl *Table, idxCols []*IndexColumn, fulltext, spatial, unique bool) error {
	return c.validateIndexColumnsOpts(tbl, idxCols, fulltext, spatial, unique, true)
}

func (c *Catalog) validateIndexColumnsOpts(tbl *Table, idxCols []*IndexColumn, fulltext, spatial, unique, checkWidths bool) error {
	for _, ic := range idxCols {
		if ic.Name == "" {
			continue
		}
		col := tbl.GetColumn(ic.Name)
		if col == nil {
			continue
		}
		// A FULLTEXT key takes the character types only: BLOB, VARBINARY, ENUM,
		// SET and every non-string type are 1283 (measured on mysql:8.0.46).
		if fulltext {
			if !isFulltextColumnType(col.DataType) {
				return errCannotPartOfFulltext(ic.Name)
			}
			// MySQL ignores any specified prefix length on a FULLTEXT key part
			// without raising 1089 and omits it from SHOW CREATE TABLE
			// (measured on mysql:8.0.46).
			ic.Length = 0
			continue
		}
		if spatial {
			// The column's type is judged before its nullability: a SPATIAL key over an
			// INT column is 1687 whether or not the column is NOT NULL, and 1252 is
			// what a nullable geometry column answers (mysql:8.0.46).
			if !isGeometryColumnType(col.DataType) {
				return errSpatialNotGeometrical()
			}
			if col.Nullable {
				return &Error{Code: 1252, SQLState: "42000", Message: "All parts of a SPATIAL index must be NOT NULL"}
			}
			if ic.Length > 0 {
				return errIncorrectPrefixKey()
			}
			continue
		}
		if strings.EqualFold(col.DataType, "json") {
			return errJSONUsedAsKey(ic.Name)
		}
		if isTextBlobType(col.DataType) && ic.Length == 0 {
			return &Error{Code: 1170, SQLState: "42000", Message: "BLOB/TEXT column used in key specification without a key length"}
		}
		if ic.Length > 0 {
			strict := c == nil || c.isStrictSQLMode()
			if err := c.validateKeyPrefixLength(tbl, col, ic, strict || unique || !checkWidths); err != nil {
				return err
			}
		}
	}
	if checkWidths && !fulltext && !spatial {
		if err := c.validateIndexWidths(tbl, idxCols, unique); err != nil {
			return err
		}
	}
	return nil
}

func isAutoTruncatableKeyType(dt string) bool {
	switch strings.ToLower(dt) {
	case "char", "varchar", "binary", "varbinary",
		"tinytext", "text", "mediumtext", "longtext",
		"tinyblob", "blob", "mediumblob", "longblob":
		return true
	default:
		return false
	}
}

func indexKeyPartCharsetBytes(col *Column) int {
	if col == nil {
		return 1
	}
	switch strings.ToLower(col.DataType) {
	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext":
		return charsetMaxBytes(col.Charset)
	default:
		return 1
	}
}

func (c *Catalog) validateIndexWidths(tbl *Table, idxCols []*IndexColumn, unique bool) error {
	partMax, hasPartMax := c.maxIndexColumnPartBytes(tbl)
	totalMax, hasTotalMax := c.maxIndexTotalKeyBytes(tbl)

	if !unique && c != nil && !c.isStrictSQLMode() && len(idxCols) == 1 {
		if ic := idxCols[0]; ic != nil && ic.Expr == "" && ic.Name != "" {
			if col := tbl.GetColumn(ic.Name); col != nil && isAutoTruncatableKeyType(col.DataType) {
				limit := 0
				hasLimit := false
				if hasPartMax {
					limit = partMax
					hasLimit = true
				}
				if hasTotalMax && (!hasLimit || totalMax < limit) {
					limit = totalMax
					hasLimit = true
				}
				if hasLimit && indexKeyPartBytes(col, ic) > limit {
					mb := indexKeyPartCharsetBytes(col)
					if mb <= 0 {
						mb = 1
					}
					if truncLen := limit / mb; truncLen > 0 {
						ic.Length = truncLen
					}
				}
			}
		}
	}

	if hasPartMax {
		for _, ic := range idxCols {
			if ic == nil || ic.Name == "" {
				continue
			}
			col := tbl.GetColumn(ic.Name)
			if col == nil {
				continue
			}
			if indexKeyPartBytes(col, ic) > partMax {
				return &Error{Code: 1071, SQLState: "42000",
					Message: fmt.Sprintf("Specified key was too long; max key length is %d bytes", partMax)}
			}
		}
	}

	return c.validateIndexTotalKeyWidth(tbl, idxCols)
}

func (c *Catalog) validateIndexTotalKeyWidth(tbl *Table, idxCols []*IndexColumn) error {
	maxBytes, ok := c.maxIndexTotalKeyBytes(tbl)
	if !ok {
		return nil
	}
	totalBytes := 0
	for _, ic := range idxCols {
		if ic == nil || ic.Name == "" {
			continue
		}
		col := tbl.GetColumn(ic.Name)
		if col == nil {
			continue
		}
		// A functional part whose width the catalog could not infer counts
		// toward neither limit: its placeholder VARCHAR(255) is no measure.
		if ic.Expr != "" && !functionalTypeWidthKnown(inferFunctionalIndexExprType(tbl, ic.ExprNode)) {
			continue
		}
		partBytes := indexKeyPartBytes(col, ic)
		// A functional key part's hidden column meets the per-part limit a
		// named column does, in strict mode: mysql:8.0.46 answers 1071 (767
		// bytes) for KEY ((LOWER(varchar_192_utf8mb4))) on a COMPACT table and
		// takes it with sql_mode ''.
		if ic.Expr != "" && isInnoDBEngine(tbl.Engine) && c.isStrictSQLMode() {
			if limit := c.maxInnoDBColumnPartBytes(tbl); partBytes > limit {
				return &Error{Code: 1071, SQLState: "42000",
					Message: fmt.Sprintf("Specified key was too long; max key length is %d bytes", limit)}
			}
		}
		totalBytes += partBytes
	}
	if totalBytes > maxBytes {
		return &Error{Code: 1071, SQLState: "42000",
			Message: fmt.Sprintf("Specified key was too long; max key length is %d bytes", maxBytes)}
	}
	return nil
}

func charsetMaxBytes(charset string) int {
	switch strings.ToLower(strings.TrimSpace(charset)) {
	case "utf8mb4", "utf16", "utf16le", "utf32", "gb18030":
		return 4
	case "utf8", "utf8mb3", "ujis", "eucjpms":
		return 3
	case "ucs2", "sjis", "cp932", "gbk", "gb2312", "big5", "euckr":
		return 2
	default:
		return 1
	}
}

func (c *Catalog) maxInnoDBColumnPartBytes(tbl *Table) int {
	if tbl != nil {
		rf := strings.ToUpper(strings.TrimSpace(tbl.RowFormat))
		if rf == "COMPACT" || rf == "REDUNDANT" {
			return 767
		}
	}
	return c.maxInnoDBIndexKeyBytes(tbl)
}

func (c *Catalog) maxInnoDBIndexKeyBytes(tbl *Table) int {
	switch c.InnoDBPageSize() {
	case 4096:
		return 768
	case 8192:
		return 1536
	default:
		return 3072
	}
}

func (c *Catalog) maxIndexColumnPartBytes(tbl *Table) (int, bool) {
	if tbl == nil {
		return 0, false
	}
	if strings.EqualFold(tbl.Engine, "MyISAM") {
		return 1000, true
	}
	if isInnoDBEngine(tbl.Engine) {
		return c.maxInnoDBColumnPartBytes(tbl), true
	}
	return 0, false
}

func (c *Catalog) maxIndexTotalKeyBytes(tbl *Table) (int, bool) {
	if tbl == nil {
		return 0, false
	}
	if strings.EqualFold(tbl.Engine, "MyISAM") {
		return 1000, true
	}
	if isInnoDBEngine(tbl.Engine) {
		return c.maxInnoDBIndexKeyBytes(tbl), true
	}
	return 0, false
}

func (c *Catalog) maxBlobTextIndexPrefixBytes(tbl *Table, dt string) int {
	if dt == "tinytext" || dt == "tinyblob" {
		return 255
	}
	if tbl != nil && strings.EqualFold(tbl.Engine, "MyISAM") {
		return 1000
	}
	return c.maxInnoDBColumnPartBytes(tbl)
}

func indexKeyPartBytes(col *Column, ic *IndexColumn) int {
	if col == nil || ic == nil {
		return 0
	}
	dt := strings.ToLower(col.DataType)
	mb := 1
	switch dt {
	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext":
		mb = charsetMaxBytes(col.Charset)
	}
	if ic.Length > 0 {
		return ic.Length * mb
	}
	if maxLen, ok := fixedLengthStringColumnLimit(col); ok {
		return maxLen * mb
	}
	return fixedWidthScalarColumnBytes(col)
}

func fixedWidthScalarColumnBytes(col *Column) int {
	dt := strings.ToLower(col.DataType)
	switch dt {
	case "tinyint", "bool", "boolean", "year":
		return 1
	case "smallint":
		return 2
	case "mediumint", "date":
		return 3
	case "int", "integer":
		return 4
	case "bigint", "double", "real", "double precision":
		return 8
	case "float":
		if p, _, ok := parseNumericTypeArgs(col.ColumnType); ok && p > 24 {
			return 8
		}
		return 4
	case "time":
		return 3 + fspStorageBytes(temporalColumnFSP(col.ColumnType))
	case "timestamp":
		return 4 + fspStorageBytes(temporalColumnFSP(col.ColumnType))
	case "datetime":
		return 5 + fspStorageBytes(temporalColumnFSP(col.ColumnType))
	case "decimal", "numeric":
		prec, scale, ok := parseNumericTypeArgs(col.ColumnType)
		if !ok {
			prec, scale = 10, 0
		}
		if prec < scale {
			prec = scale
		}
		return decimalPackedBytes(prec-scale) + decimalPackedBytes(scale)
	case "bit":
		bits, _, ok := parseNumericTypeArgs(col.ColumnType)
		if !ok || bits <= 0 {
			bits = 1
		}
		return (bits + 7) / 8
	case "enum":
		if countQuotedListArgs(col.ColumnType) > 255 {
			return 2
		}
		return 1
	case "set":
		n := countQuotedListArgs(col.ColumnType)
		if n <= 0 {
			return 1
		}
		b := (n + 7) / 8
		if b > 4 {
			return 8
		}
		return b
	default:
		return 0
	}
}

func fspStorageBytes(fsp int) int {
	switch {
	case fsp <= 0:
		return 0
	case fsp <= 2:
		return 1
	case fsp <= 4:
		return 2
	default:
		return 3
	}
}

func decimalPackedBytes(digits int) int {
	if digits <= 0 {
		return 0
	}
	dig2bytes := [10]int{0, 1, 1, 2, 2, 3, 3, 4, 4, 4}
	return (digits/9)*4 + dig2bytes[digits%9]
}

func parseNumericTypeArgs(colType string) (int, int, bool) {
	open := strings.IndexByte(colType, '(')
	close := strings.LastIndexByte(colType, ')')
	if open < 0 || close <= open+1 {
		return 0, 0, false
	}
	inner := colType[open+1 : close]
	parts := strings.Split(inner, ",")
	a, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, false
	}
	b := 0
	if len(parts) > 1 {
		if bv, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
			b = bv
		}
	}
	return a, b, true
}

func countQuotedListArgs(colType string) int {
	open := strings.IndexByte(colType, '(')
	close := strings.LastIndexByte(colType, ')')
	if open < 0 || close <= open+1 {
		return 0
	}
	inner := colType[open+1 : close]
	count := 1
	inQuote := false
	for i := 0; i < len(inner); i++ {
		switch inner[i] {
		case '\'':
			if inQuote && i+1 < len(inner) && inner[i+1] == '\'' {
				i++
			} else {
				inQuote = !inQuote
			}
		case '\\':
			if inQuote && i+1 < len(inner) {
				i++
			}
		case ',':
			if !inQuote {
				count++
			}
		}
	}
	return count
}

// validateKeyPrefixLength validates a non-zero key part prefix length against
// the column's data type and declared width (1089 for fixed-width types, 1071
// or non-strict clamping for BLOB/TEXT types, measured on mysql:8.0.46).
func (c *Catalog) validateKeyPrefixLength(tbl *Table, col *Column, ic *IndexColumn, strict bool) error {
	prefixLen := ic.Length
	dt := strings.ToLower(col.DataType)
	// MySQL holds a key part's byte length in 32 bits, so a prefix whose
	// bytes pass 2^32-1 wraps: on a utf8mb4 TEXT, mysql:8.0.46 stores
	// s(1073741825) as s(1) and answers 1170 for s(1073741824), whose bytes
	// wrap to 0, as for no prefix; VARCHAR(10) takes s(1073741824) whole.
	// The parser admits no prefix past 2^31-1, so 64 bits never overflow.
	mb := uint64(1)
	switch dt {
	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext":
		mb = uint64(max(charsetMaxBytes(col.Charset), 1))
	}
	if bytes := uint64(prefixLen) * mb; bytes > math.MaxUint32 {
		prefixLen = int(uint64(uint32(bytes)) / mb)
		ic.Length = prefixLen
		if prefixLen == 0 {
			if isTextBlobType(col.DataType) {
				return &Error{Code: 1170, SQLState: "42000", Message: "BLOB/TEXT column used in key specification without a key length"}
			}
			return nil
		}
	}
	switch dt {
	case "char", "varchar", "binary", "varbinary":
		if maxLen, ok := fixedLengthStringColumnLimit(col); ok && prefixLen > maxLen {
			return errIncorrectPrefixKey()
		}
		return nil
	case "tinytext", "text", "mediumtext", "longtext",
		"tinyblob", "blob", "mediumblob", "longblob":
		maxBytes := c.maxBlobTextIndexPrefixBytes(tbl, dt)
		mb := 1
		if strings.HasSuffix(dt, "text") {
			mb = charsetMaxBytes(col.Charset)
		}
		if int64(prefixLen)*int64(mb) > int64(maxBytes) {
			if strict {
				return &Error{Code: 1071, SQLState: "42000",
					Message: fmt.Sprintf("Specified key was too long; max key length is %d bytes", maxBytes)}
			}
			ic.Length = maxBytes / mb
		}
		return nil
	default:
		return errIncorrectPrefixKey()
	}
}

// isFulltextColumnType reports whether a column's type may be part of a
// FULLTEXT index: the character types, beside which MySQL refuses BLOB, VARBINARY,
// ENUM and SET with 1283.
func isFulltextColumnType(dt string) bool {
	switch strings.ToLower(dt) {
	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext":
		return true
	}
	return false
}

// isGeometryColumnType reports whether a column's type is one a SPATIAL key takes.
func isGeometryColumnType(dt string) bool {
	switch strings.ToLower(dt) {
	case "geometry", "geomcollection", "geometrycollection", "point", "linestring", "polygon",
		"multipoint", "multilinestring", "multipolygon":
		return true
	}
	return false
}

// isSpatialIndexType reports whether an index type takes the rules of a SPATIAL
// key. USING RTREE makes the key a spatial one, whatever its own type says.
func isSpatialIndexType(indexType string) bool {
	switch strings.ToUpper(strings.TrimSpace(indexType)) {
	case "SPATIAL", "RTREE":
		return true
	}
	return false
}

// validateUniqueOrPrimaryNotSpatial rejects PRIMARY KEY / UNIQUE KEY / CREATE
// UNIQUE INDEX with USING RTREE: mysql:8.0.46 checks that every key part is
// geometrical (1687) and then refuses unique/primary spatial indexes with 3728.
func validateUniqueOrPrimaryNotSpatial(tbl *Table, idxCols []*IndexColumn, indexType string) error {
	if !isSpatialIndexType(indexType) {
		return nil
	}
	for _, ic := range idxCols {
		if ic == nil || ic.Name == "" {
			continue
		}
		if col := tbl.GetColumn(ic.Name); col != nil && !isGeometryColumnType(col.DataType) {
			return errSpatialNotGeometrical()
		}
	}
	return errSpatialUniqueIndex()
}

// storedIndexType is the type an index whose USING names a spatial access method is
// stored as: MySQL keeps SPATIAL KEY and drops the RTREE spelling, so a table built
// with `KEY USING RTREE (g)` reads back as `SPATIAL KEY g (g)` (mysql:8.0.46).
func storedIndexType(indexType string) string {
	if strings.EqualFold(strings.TrimSpace(indexType), "RTREE") {
		return "SPATIAL"
	}
	return indexType
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

// coerceInnoDBHashIndex folds a meaningless `USING HASH` on an InnoDB index to the no-explicit-type
// form (IndexType ""), so a HASH index on InnoDB round-trips empty on both 5.7 and 8.0.
//
// It must be called from EVERY path that builds a *secondary* index from a user `USING HASH`, or
// that path phantom-diffs. Call sites: CREATE TABLE UNIQUE + plain KEY (tablecmds.go), ALTER ADD
// UNIQUE + plain KEY (altercmds.go), and standalone CREATE [UNIQUE] INDEX (indexcmds.go). The
// PRIMARY KEY paths are intentionally excluded: they hard-code IndexType "" (MySQL never accepts or
// echoes a USING clause on a PK, so there is nothing to fold).
//
// InnoDB does not implement hash indexes: it silently treats `USING HASH` as BTREE. The SHOW CREATE
// echo of that meaningless clause is version-, kind-, and engine-divergent (verified on live 5.7.25
// + 8.0.32, across all of the above authoring paths):
//   - 8.0 InnoDB DROPS the clause entirely;
//   - 5.7 InnoDB KEEPS `USING HASH` in the echo;
//   - `USING BTREE` on InnoDB is echoed verbatim on BOTH versions (BTREE is the real storage), so it
//     is NOT folded — only HASH is meaningless;
//   - MEMORY (and other engines where HASH is a real index type) KEEP `USING HASH` on both versions,
//     so the engine guard below leaves them untouched.
//
// Folding HASH → "" (rather than → "BTREE") is what makes both echoes converge: a user's `USING
// HASH`, 5.7's kept `USING HASH`, and 8.0's dropped clause all load to IndexType "" on InnoDB, which
// the index differ (canonicalIndexType) and the index generator (formatIndexDefinition) both treat
// as "no USING clause" — matching the plain no-USING readback on either version. (The earlier
// HASH → "BTREE" fold matched 5.7's echo but still phantom-diffed against 8.0's dropped clause.)
func coerceInnoDBHashIndex(tbl *Table, idx *Index) {
	if isInnoDBEngine(tbl.Engine) && strings.EqualFold(idx.IndexType, "HASH") {
		idx.IndexType = ""
	}
}

// isInnoDBEngine reports whether a table's declared engine is InnoDB, treating an empty (omitted)
// engine as InnoDB — the MySQL default storage engine. The user form often omits ENGINE while the
// SHOW CREATE readback always echoes ENGINE=InnoDB, so the HASH fold must apply to both or a `USING
// HASH` written on an engine-less CREATE TABLE would phantom-diff against its InnoDB readback.
func isInnoDBEngine(engine string) bool {
	e := strings.TrimSpace(engine)
	return e == "" || strings.EqualFold(e, "InnoDB")
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

// columnTypeLength is the length a column type such as varchar(20) declares,
// 0 when it declares none.
func columnTypeLength(columnType string) int {
	open := strings.IndexByte(columnType, '(')
	if open < 0 {
		return 0
	}
	end := strings.IndexByte(columnType[open:], ')')
	if end < 0 {
		return 0
	}
	n, err := strconv.Atoi(strings.TrimSpace(columnType[open+1 : open+end]))
	if err != nil {
		return 0
	}
	return n
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
		// A generated column stores a string escaped, with the introducer it was written with or
		// else the column's charset, and N'x' as _utf8mb3'x': mysql:8.0.46 stores 'it''s' as
		// _utf8mb4'it\'s' and _latin1 X'27' as _latin1'\''. JSON_VALUE's strings follow the
		// same rules, so this is its writer.
		return deparse.JSONValueString(n, charset)
	case *nodes.FuncCallExpr:
		funcName := strings.ToLower(n.Name)
		if n.Star {
			return funcName + "(*)"
		}
		if n.JSONValueArm {
			// Its clauses change what it answers, and the server stores them with a space
			// after the document's comma (mysql:8.0.46).
			return jsonValueToSQLGenerated(n, charset)
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
	default:
		return "(?)"
	}
}

func jsonValueToSQLGenerated(n *nodes.FuncCallExpr, charset string) string {
	doc := nodeToSQLGenerated(n.Args[0], charset)
	path := n.Args[1].(*nodes.StringLit) // the parser holds JSON_VALUE to a literal path
	return "json_value(" + doc + ", " + deparse.JSONValueString(path, charset) +
		deparse.JSONValueClauses(n.JSONValue, charset) + ")"
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
		if name == "tinyint" && dt.Length == 1 && !dt.Unsigned && !dt.Zerofill {
			buf.WriteString("(1)")
		} else if dt.Zerofill {
			width := dt.Length
			if width == 0 {
				// ZEROFILL forces the column UNSIGNED, so a bare INT ZEROFILL takes the
				// UNSIGNED default width (int → 10), not the signed default (11). MySQL
				// stores `int(10) unsigned zerofill` for `INT ZEROFILL`. (Passing
				// dt.Unsigned here was a bug: the zerofill-implies-unsigned promotion has
				// not yet flipped dt.Unsigned at this point.)
				width = defaultIntDisplayWidth(name, true)
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
	// YEAR takes UNSIGNED and ZEROFILL and keeps neither: every engine's SHOW
	// CREATE TABLE writes `year` for `year unsigned zerofill`, so a column that
	// was declared with them is not one that reports them.
	if name != "year" {
		if dt.Unsigned || dt.Zerofill {
			buf.WriteString(" unsigned")
		}
		if dt.Zerofill {
			buf.WriteString(" zerofill")
		}
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
// MySQL escapes a single quote in enum values by doubling it.
func escapeEnumValue(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// createTableLike implements CREATE TABLE t2 LIKE t1.
// It copies the structure (columns, indexes, constraints) from the source table.
func likeAutoExtendSize(srcTbl *Table, temporary bool) uint64 {
	if temporary || srcTbl == nil {
		return 0
	}
	return srcTbl.AutoExtendSize
}

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
	if isInnoDBEngine(srcTbl.Engine) && srcTbl.Tablespace != "" {
		if err := c.validateInnoDBTablespace(tableName, stmt.Temporary, srcTbl.Tablespace); err != nil {
			return err
		}
	}

	tbl := &Table{
		Name:             tableName,
		Database:         db,
		Columns:          make([]*Column, 0, len(srcTbl.Columns)),
		colByName:        make(map[string]int),
		Indexes:          make([]*Index, 0, len(srcTbl.Indexes)),
		Constraints:      make([]*Constraint, 0, len(srcTbl.Constraints)),
		Engine:           srcTbl.Engine,
		SecondaryEngine:  srcTbl.SecondaryEngine,
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
		AutoExtendSize:   likeAutoExtendSize(srcTbl, stmt.Temporary),
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
			OnUpdateKind:                 srcCol.OnUpdateKind,
			Invisible:                    srcCol.Invisible,
			NotSecondary:                 srcCol.NotSecondary,
			GeneratedInvisiblePrimaryKey: srcCol.GeneratedInvisiblePrimaryKey,
			Hidden:                       srcCol.Hidden,
			// SRID copied by presence: MySQL's CREATE TABLE ... LIKE preserves the source
			// column's SRID clause (incl. an explicit `SRID 0`), so the cloned column must
			// carry both the value and its presence flag, not silently drop to no-SRID.
			SRID:    srcCol.SRID,
			HasSRID: srcCol.HasSRID,
		}
		if srcCol.Default != nil {
			def := *srcCol.Default
			col.Default = &def
			if srcCol.DefaultTimestampUTC != nil {
				tsUTC := *srcCol.DefaultTimestampUTC
				col.DefaultTimestampUTC = &tsUTC
			}
			col.DefaultKind = srcCol.DefaultKind
			col.DefaultIsExpr = srcCol.DefaultIsExpr
			col.DefaultExprNode = srcCol.DefaultExprNode
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

	likeSrcDB := srcDB.Name
	if srcTbl.Database != nil {
		likeSrcDB = srcTbl.Database.Name
	}
	stripTableQualifiersFromTableExprs(tbl, likeSrcDB, srcTbl.Name, true, nil)

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

// normalizedChecksum is the value CHECKSUM reports, which is not the value it
// was given: the option is a flag, so the server records any non-zero number as
// 1 and a zero as no option at all. `CHECKSUM=2`, `TABLE_CHECKSUM=1` and
// `CHECKSUM=1.5` all read back as `CHECKSUM=1`, and `CHECKSUM=0` reads back as
// nothing (MySQL 8.0 and MariaDB 11.8.8, information_schema.tables.create_options).
//
// TABLE_CHECKSUM is CHECKSUM's other name and reports under CHECKSUM, so it
// comes here too.
func normalizedChecksum(value string) string {
	for _, c := range value {
		if c != '0' && c != '.' && c != '+' && c != '-' {
			return "1"
		}
	}
	return "0"
}
