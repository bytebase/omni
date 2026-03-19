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

	// Check for duplicate table.
	if db.Tables[key] != nil {
		if stmt.IfNotExists {
			return nil
		}
		return errDupTable(tableName)
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
	for _, opt := range stmt.Options {
		switch toLower(opt.Name) {
		case "engine":
			tbl.Engine = opt.Value
		case "charset", "character set", "default charset", "default character set":
			tbl.Charset = opt.Value
		case "collate", "default collate":
			tbl.Collation = opt.Value
		case "comment":
			tbl.Comment = opt.Value
		case "auto_increment":
			fmt.Sscanf(opt.Value, "%d", &tbl.AutoIncrement)
		case "row_format":
			tbl.RowFormat = opt.Value
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
		}

		// Default charset/collation for string types.
		if isStringType(col.DataType) {
			if col.Charset == "" {
				col.Charset = tbl.Charset
			}
			if col.Collation == "" {
				col.Collation = tbl.Collation
			}
		}

		// Top-level column properties.
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
				Expr:   nodeToSQL(colDef.Generated.Expr),
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
					conName = fmt.Sprintf("%s_chk_%d", tableName, len(tbl.Constraints)+1)
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
					conName = fmt.Sprintf("%s_ibfk_%d", tableName, len(tbl.Constraints)+1)
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
				// Add implicit backing index for FK.
				idxName := allocIndexName(tbl, colDef.Name)
				tbl.Indexes = append(tbl.Indexes, &Index{
					Name:      idxName,
					Table:     tbl,
					Columns:   []*IndexColumn{{Name: colDef.Name}},
					Unique:    false,
					IndexType: "BTREE",
					Visible:   true,
				})
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
				IndexType: "BTREE",
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
					IndexType: "BTREE",
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
					IndexType: "BTREE",
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
			tbl.Indexes = append(tbl.Indexes, &Index{
				Name:      "PRIMARY",
				Table:     tbl,
				Columns:   idxCols,
				Unique:    true,
				Primary:   true,
				IndexType: "BTREE",
				Visible:   true,
			})
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
			tbl.Indexes = append(tbl.Indexes, &Index{
				Name:      idxName,
				Table:     tbl,
				Columns:   idxCols,
				Unique:    true,
				IndexType: indexTypeOrDefault(con.IndexType, "BTREE"),
				Visible:   true,
			})
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
				conName = fmt.Sprintf("%s_ibfk_%d", tableName, len(tbl.Constraints)+1)
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
			// Add implicit backing index for FK.
			idxName := con.Name
			if idxName == "" && len(cols) > 0 {
				idxName = allocIndexName(tbl, cols[0])
			}
			idxCols := buildIndexColumns(con)
			tbl.Indexes = append(tbl.Indexes, &Index{
				Name:      idxName,
				Table:     tbl,
				Columns:   idxCols,
				IndexType: "BTREE",
				Visible:   true,
			})

		case nodes.ConstrCheck:
			conName := con.Name
			if conName == "" {
				conName = fmt.Sprintf("%s_chk_%d", tableName, len(tbl.Constraints)+1)
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
			tbl.Indexes = append(tbl.Indexes, &Index{
				Name:      idxName,
				Table:     tbl,
				Columns:   idxCols,
				IndexType: indexTypeOrDefault(con.IndexType, "BTREE"),
				Visible:   true,
			})

		case nodes.ConstrFulltextIndex:
			idxName := con.Name
			if idxName == "" && len(cols) > 0 {
				idxName = allocIndexName(tbl, cols[0])
			}
			idxCols := buildIndexColumns(con)
			tbl.Indexes = append(tbl.Indexes, &Index{
				Name:      idxName,
				Table:     tbl,
				Columns:   idxCols,
				Fulltext:  true,
				IndexType: "FULLTEXT",
				Visible:   true,
			})

		case nodes.ConstrSpatialIndex:
			idxName := con.Name
			if idxName == "" && len(cols) > 0 {
				idxName = allocIndexName(tbl, cols[0])
			}
			idxCols := buildIndexColumns(con)
			tbl.Indexes = append(tbl.Indexes, &Index{
				Name:      idxName,
				Table:     tbl,
				Columns:   idxCols,
				Spatial:   true,
				IndexType: "SPATIAL",
				Visible:   true,
			})
		}
	}

	db.Tables[key] = tbl
	return nil
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

func isStringType(dt string) bool {
	switch dt {
	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext",
		"enum", "set":
		return true
	}
	return false
}

func nodeToSQL(node nodes.ExprNode) string {
	if node == nil {
		return ""
	}
	switch n := node.(type) {
	case *nodes.ColumnRef:
		if n.Table != "" {
			return n.Table + "." + n.Column
		}
		return n.Column
	case *nodes.IntLit:
		return fmt.Sprintf("%d", n.Value)
	case *nodes.StringLit:
		return "'" + n.Value + "'"
	case *nodes.FuncCallExpr:
		if n.Star {
			return n.Name + "(*)"
		}
		var args []string
		for _, a := range n.Args {
			args = append(args, nodeToSQL(a))
		}
		return n.Name + "(" + strings.Join(args, ", ") + ")"
	case *nodes.NullLit:
		return "NULL"
	case *nodes.BoolLit:
		if n.Value {
			return "TRUE"
		}
		return "FALSE"
	default:
		return "(?)"
	}
}

func formatColumnType(dt *nodes.DataType) string {
	name := strings.ToLower(dt.Name)

	// MySQL type aliases: BOOLEAN/BOOL → tinyint(1), NUMERIC → decimal, SERIAL → bigint unsigned
	switch name {
	case "boolean":
		return "tinyint(1)"
	case "numeric":
		name = "decimal"
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
	} else if name == "decimal" && dt.Length == 0 && dt.Scale == 0 {
		// DECIMAL with no precision → MySQL shows decimal(10,0)
		buf.WriteString("(10,0)")
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
			buf.WriteString("'" + v + "'")
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
