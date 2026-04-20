package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

func (c *Catalog) createRoutine(stmt *nodes.CreateFunctionStmt) error {
	db, err := c.resolveDatabase(stmt.Name.Schema)
	if err != nil {
		return err
	}

	name := stmt.Name.Name
	key := toLower(name)

	routineMap := db.Functions
	if stmt.IsProcedure {
		routineMap = db.Procedures
	}

	if _, exists := routineMap[key]; exists {
		if !stmt.IfNotExists {
			if stmt.IsProcedure {
				return errDupProcedure(name)
			}
			return errDupFunction(name)
		}
		return nil
	}

	// Build params.
	params := make([]*RoutineParam, 0, len(stmt.Params))
	for _, p := range stmt.Params {
		params = append(params, &RoutineParam{
			Direction: p.Direction,
			Name:      p.Name,
			TypeName:  formatParamType(p.TypeName),
		})
	}

	// Build return type -- MySQL shows lowercase with CHARSET for string types.
	var returns string
	if stmt.Returns != nil {
		returns = formatReturnType(stmt.Returns, db.Charset)
	}

	// MySQL always sets a definer. Default to `root`@`%` when not specified.
	definer := stmt.Definer
	if definer == "" {
		definer = "`root`@`%`"
	}

	// Build characteristics.
	chars := make(map[string]string)
	for _, ch := range stmt.Characteristics {
		chars[ch.Name] = ch.Value
	}

	routine := &Routine{
		Name:            name,
		Database:        db,
		IsProcedure:     stmt.IsProcedure,
		Definer:         definer,
		Params:          params,
		Returns:         returns,
		Body:            strings.TrimSpace(stmt.Body),
		Characteristics: chars,
	}

	routineMap[key] = routine
	return nil
}

func (c *Catalog) dropRoutine(stmt *nodes.DropRoutineStmt) error {
	db, err := c.resolveDatabase(stmt.Name.Schema)
	if err != nil {
		if stmt.IfExists {
			return nil
		}
		return err
	}

	name := stmt.Name.Name
	key := toLower(name)

	routineMap := db.Functions
	if stmt.IsProcedure {
		routineMap = db.Procedures
	}

	if _, exists := routineMap[key]; !exists {
		if stmt.IfExists {
			return nil
		}
		if stmt.IsProcedure {
			return errNoSuchProcedure(db.Name, name)
		}
		return errNoSuchFunction(name)
	}

	delete(routineMap, key)
	return nil
}

func (c *Catalog) alterRoutine(stmt *nodes.AlterRoutineStmt) error {
	db, err := c.resolveDatabase(stmt.Name.Schema)
	if err != nil {
		return err
	}

	name := stmt.Name.Name
	key := toLower(name)

	routineMap := db.Functions
	if stmt.IsProcedure {
		routineMap = db.Procedures
	}

	routine, exists := routineMap[key]
	if !exists {
		if stmt.IsProcedure {
			return errNoSuchProcedure(db.Name, name)
		}
		return errNoSuchFunction(name)
	}

	// Update characteristics.
	for _, ch := range stmt.Characteristics {
		routine.Characteristics[ch.Name] = ch.Value
	}

	return nil
}

// formatDataType formats a DataType node into a display string for routine parameters/returns.
func formatDataType(dt *nodes.DataType) string {
	if dt == nil {
		return ""
	}

	name := strings.ToLower(dt.Name)

	switch name {
	case "int", "integer":
		name = "int"
	case "tinyint":
		name = "tinyint"
	case "smallint":
		name = "smallint"
	case "mediumint":
		name = "mediumint"
	case "bigint":
		name = "bigint"
	case "float":
		name = "float"
	case "double", "real":
		name = "double"
	case "decimal", "numeric", "dec", "fixed":
		name = "decimal"
	case "varchar":
		name = "varchar"
	case "char":
		name = "char"
	case "text", "tinytext", "mediumtext", "longtext":
		// keep as-is
	case "blob", "tinyblob", "mediumblob", "longblob":
		// keep as-is
	case "date", "time", "datetime", "timestamp", "year":
		// keep as-is
	case "json":
		name = "json"
	case "bool", "boolean":
		name = "tinyint"
	}

	var b strings.Builder
	b.WriteString(name)

	// Length/precision.
	if dt.Length > 0 {
		if dt.Scale > 0 {
			b.WriteString(fmt.Sprintf("(%d,%d)", dt.Length, dt.Scale))
		} else {
			b.WriteString(fmt.Sprintf("(%d)", dt.Length))
		}
	}

	if dt.Unsigned {
		b.WriteString(" unsigned")
	}

	return b.String()
}

// formatParamType formats a DataType for display in a routine parameter list.
// MySQL 8.0 shows parameter types in UPPERCASE (INT, VARCHAR(100), etc.)
func formatParamType(dt *nodes.DataType) string {
	raw := formatDataType(dt)
	return strings.ToUpper(raw)
}

// formatReturnType formats a DataType for display in the RETURNS clause.
// MySQL 8.0 shows return types in lowercase but adds CHARSET for string types.
func formatReturnType(dt *nodes.DataType, dbCharset string) string {
	raw := formatDataType(dt)
	// MySQL 8.0 appends CHARSET for string return types.
	name := strings.ToLower(dt.Name)
	if isStringRoutineType(name) {
		charset := dt.Charset
		if charset == "" {
			charset = dbCharset
		}
		if charset == "" {
			charset = "utf8mb4"
		}
		raw += " CHARSET " + charset
	}
	return raw
}

// isStringRoutineType returns true for types where MySQL shows CHARSET in routine RETURNS.
func isStringRoutineType(dt string) bool {
	switch dt {
	case "varchar", "char", "text", "tinytext", "mediumtext", "longtext",
		"enum", "set":
		return true
	}
	return false
}

// ShowCreateFunction produces MySQL 8.0-compatible SHOW CREATE FUNCTION output.
func (c *Catalog) ShowCreateFunction(database, name string) string {
	db := c.GetDatabase(database)
	if db == nil {
		return ""
	}
	routine := db.Functions[toLower(name)]
	if routine == nil {
		return ""
	}
	return showCreateRoutine(routine)
}

// ShowCreateProcedure produces MySQL 8.0-compatible SHOW CREATE PROCEDURE output.
func (c *Catalog) ShowCreateProcedure(database, name string) string {
	db := c.GetDatabase(database)
	if db == nil {
		return ""
	}
	routine := db.Procedures[toLower(name)]
	if routine == nil {
		return ""
	}
	return showCreateRoutine(routine)
}

// showCreateRoutine produces the SHOW CREATE output for a stored routine.
// MySQL 8.0 format:
//
//	CREATE DEFINER=`root`@`%` FUNCTION `name`(a INT, b INT) RETURNS int
//	    DETERMINISTIC
//	RETURN a + b
func showCreateRoutine(r *Routine) string {
	var b strings.Builder

	b.WriteString("CREATE")

	// DEFINER -- MySQL 8.0 always shows DEFINER.
	if r.Definer != "" {
		b.WriteString(fmt.Sprintf(" DEFINER=%s", r.Definer))
	}

	if r.IsProcedure {
		b.WriteString(fmt.Sprintf(" PROCEDURE `%s`(", r.Name))
	} else {
		b.WriteString(fmt.Sprintf(" FUNCTION `%s`(", r.Name))
	}

	// Parameters -- MySQL 8.0 separates with ", " (comma space).
	for i, p := range r.Params {
		if i > 0 {
			b.WriteString(", ")
		}
		if r.IsProcedure && p.Direction != "" {
			b.WriteString(fmt.Sprintf("%s %s %s", p.Direction, p.Name, p.TypeName))
		} else {
			b.WriteString(fmt.Sprintf("%s %s", p.Name, p.TypeName))
		}
	}
	b.WriteString(")")

	// RETURNS (functions only)
	if !r.IsProcedure && r.Returns != "" {
		b.WriteString(fmt.Sprintf(" RETURNS %s", r.Returns))
	}

	// Characteristics -- MySQL 8.0 outputs each on its own line with 4-space indent.
	// Order: DETERMINISTIC, DATA ACCESS, SQL SECURITY, COMMENT
	if v, ok := r.Characteristics["DETERMINISTIC"]; ok {
		if v == "YES" {
			b.WriteString("\n    DETERMINISTIC")
		} else {
			b.WriteString("\n    NOT DETERMINISTIC")
		}
	}

	if v, ok := r.Characteristics["DATA ACCESS"]; ok {
		b.WriteString(fmt.Sprintf("\n    %s", v))
	}

	if v, ok := r.Characteristics["SQL SECURITY"]; ok {
		b.WriteString(fmt.Sprintf("\n    SQL SECURITY %s", strings.ToUpper(v)))
	}

	if v, ok := r.Characteristics["COMMENT"]; ok {
		b.WriteString(fmt.Sprintf("\n    COMMENT '%s'", escapeComment(v)))
	}

	// Body -- MySQL 8.0 starts body on its own line with no indent.
	if r.Body != "" {
		b.WriteString(fmt.Sprintf("\n%s", r.Body))
	}

	return b.String()
}
