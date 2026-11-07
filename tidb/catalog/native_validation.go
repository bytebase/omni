package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

func tiDBVirtualFKColumn(constraint, column string) error {
	return &Error{Code: 3733, SQLState: "HY000", Message: fmt.Sprintf("Foreign key '%s' uses virtual column '%s' which is not supported.", constraint, column)}
}

// checkModifyCharsetAndCollation permits utf8/latin1 -> utf8mb4, while
// DEFAULT CHARACTER SET keeps existing columns unchanged.
// https://github.com/pingcap/tidb/blob/v8.5.5/pkg/ddl/executor.go
func tiDBTableCharsetChangeAllowed(from, to string) bool {
	from, to = strings.ToLower(from), strings.ToLower(to)
	return from == "" || from == to || to == "utf8mb4" && (from == "utf8" || from == "utf8mb3" || from == "latin1")
}

// FLOAT(p) chooses a storage type by binary precision; the same choice must
// reach DATA_TYPE as well as the full type used by SHOW CREATE TABLE.
// https://github.com/pingcap/tidb/blob/v8.5.5/pkg/parser/parser.y
func tiDBDataType(dt *nodes.DataType) string {
	switch name := strings.ToLower(dt.Name); name {
	case "boolean":
		return "tinyint"
	case "numeric":
		return "decimal"
	case "serial":
		return "bigint"
	case "float":
		if dt.Length > 24 && !dt.HasScale {
			return "double"
		}
		return name
	default:
		return name
	}
}

// The whitelist is TiDB v8.5.5's getFuncCallDefaultValue, not MySQL's
// general expression-default support.
// https://github.com/pingcap/tidb/blob/v8.5.5/pkg/ddl/add_column.go
func validateTiDBColumn(def *nodes.ColumnDef) error {
	if def == nil {
		return nil
	}
	if err := validateTiDBColumnNullability(def, false); err != nil {
		return err
	}
	if def.TypeName != nil {
		switch strings.ToLower(def.TypeName.Name) {
		case "decimal", "numeric":
			dt := def.TypeName
			if dt.Length > 65 {
				return &Error{Code: 1426, SQLState: "42000", Message: fmt.Sprintf("Too-big precision %d specified for '%s'. Maximum is 65.", dt.Length, def.Name)}
			}
			if dt.Scale > 30 {
				return &Error{Code: 1425, SQLState: "42000", Message: fmt.Sprintf("Too big scale %d specified for column '%s'. Maximum is 30.", dt.Scale, def.Name)}
			}
			if dt.Scale > dt.Length {
				return &Error{Code: 1427, SQLState: "42000", Message: fmt.Sprintf("For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column '%s').", def.Name)}
			}
		case "datetime", "timestamp", "time":
			if def.TypeName.Length > 6 {
				return &Error{Code: 1426, SQLState: "42000", Message: fmt.Sprintf("Too big precision %d specified for '%s'. Maximum is 6.", def.TypeName.Length, def.Name)}
			}
		}
	}
	if def.OnUpdate != nil {
		fn, ok := def.OnUpdate.(*nodes.FuncCallExpr)
		if !ok || !tiDBCurrentTimestamp(fn.Name) || def.TypeName == nil || (!strings.EqualFold(def.TypeName.Name, "datetime") && !strings.EqualFold(def.TypeName.Name, "timestamp")) || tiDBFunctionFSP(fn) != def.TypeName.Length {
			return &Error{Code: 1294, SQLState: "HY000", Message: fmt.Sprintf("Invalid ON UPDATE clause for '%s' column", def.Name)}
		}
	}
	defaults := []nodes.ExprNode{def.DefaultValue}
	for _, con := range def.Constraints {
		if con.Type == nodes.ColConstrDefault {
			defaults = append(defaults, con.Expr)
		}
	}
	for _, expr := range defaults {
		for {
			if p, ok := expr.(*nodes.ParenExpr); ok {
				expr = p.Expr
			} else {
				break
			}
		}
		if fn, ok := expr.(*nodes.FuncCallExpr); ok {
			if def.TypeName != nil && tiDBCurrentTimestamp(fn.Name) && (strings.EqualFold(def.TypeName.Name, "datetime") || strings.EqualFold(def.TypeName.Name, "timestamp")) && tiDBFunctionFSP(fn) != def.TypeName.Length {
				return &Error{Code: 1067, SQLState: "42000", Message: fmt.Sprintf("Invalid default value for '%s'", def.Name)}
			}
			allowed := false
			switch strings.ToLower(fn.Name) {
			case "now", "current_timestamp", "curdate", "current_date", "localtime", "localtimestamp", "nextval", "rand", "uuid", "uuid_to_bin", "json_object", "json_array", "json_quote", "vec_from_text":
				allowed = true
			case "date_format":
				if len(fn.Args) == 2 {
					if now, ok := fn.Args[0].(*nodes.FuncCallExpr); ok && strings.EqualFold(now.Name, "now") {
						if format, ok := fn.Args[1].(*nodes.StringLit); ok {
							switch format.Value {
							case "%Y-%m", "%Y-%m-%d", "%Y-%m-%d %H.%i.%s", "%Y-%m-%d %H:%i:%s":
								allowed = true
							}
						}
					}
				}
			case "upper":
				if len(fn.Args) == 1 {
					if sub, ok := fn.Args[0].(*nodes.FuncCallExpr); ok && strings.EqualFold(sub.Name, "substring_index") && len(sub.Args) == 3 {
						if user, ok := sub.Args[0].(*nodes.FuncCallExpr); ok && strings.EqualFold(user.Name, "user") {
							if delimiter, ok := sub.Args[1].(*nodes.StringLit); ok && delimiter.Value == "@" {
								allowed = true
							}
						}
					}
				}
			case "replace":
				if len(fn.Args) == 3 {
					first := fn.Args[0]
					if conv, ok := first.(*nodes.ConvertExpr); ok {
						first = conv.Expr
					}
					if up, ok := first.(*nodes.FuncCallExpr); ok && strings.EqualFold(up.Name, "upper") && len(up.Args) == 1 {
						if uuid, ok := up.Args[0].(*nodes.FuncCallExpr); ok && strings.EqualFold(uuid.Name, "uuid") && len(uuid.Args) == 0 {
							allowed = true
						}
					}
				}
			case "str_to_date":
				if len(fn.Args) == 2 {
					_, first := fn.Args[0].(*nodes.StringLit)
					_, second := fn.Args[1].(*nodes.StringLit)
					allowed = first && second
				}
			}
			if !allowed {
				return &Error{Code: 3770, SQLState: "HY000", Message: fmt.Sprintf("Default value expression of column '%s' contains a disallowed function: `%s`.", def.Name, fn.Name)}
			}
		}
	}
	if err := validateTiDBFunctionPrecision(def); err != nil {
		return err
	}
	if def.Generated != nil {
		invalid := tiDBIllegalGeneratedExpr(def.Generated.Expr)
		if invalid {
			return &Error{Code: 3102, SQLState: "HY000", Message: fmt.Sprintf("Expression of generated column '%s' contains a disallowed function.", def.Name)}
		}
	}
	return nil
}

// checkPriKeyConstraint in TiDB's add_column.go rejects explicit NULL on a
// primary key, and checks a NULL default before that explicit-NULL flag.
func validateTiDBColumnNullability(def *nodes.ColumnDef, primary bool) error {
	notNull, explicitNull := false, false
	defaultExpr := def.DefaultValue
	for _, con := range def.Constraints {
		switch con.Type {
		case nodes.ColConstrPrimaryKey:
			primary = true
		case nodes.ColConstrNotNull:
			notNull = true
		case nodes.ColConstrNull:
			notNull, explicitNull = false, true
		case nodes.ColConstrDefault:
			defaultExpr = con.Expr
		}
	}
	for {
		if paren, ok := defaultExpr.(*nodes.ParenExpr); ok {
			defaultExpr = paren.Expr
		} else {
			break
		}
	}
	if _, isNull := defaultExpr.(*nodes.NullLit); isNull && (notNull || primary) {
		return &Error{Code: 1067, SQLState: "42000", Message: fmt.Sprintf("Invalid default value for '%s'", def.Name)}
	}
	if primary && explicitNull {
		return &Error{Code: 1171, SQLState: "42000", Message: "All parts of a PRIMARY KEY must be NOT NULL; if you need NULL in a key, use UNIQUE instead"}
	}
	return nil
}

// TiDB resolves added keys and multi-schema-change positions against the
// original table, before any subcommand has changed it.
// https://github.com/pingcap/tidb/blob/v8.5.5/pkg/ddl/multi_schema_change.go
func validateTiDBAlterColumns(tbl *Table, commands []*nodes.AlterTableCmd) error {
	for _, cmd := range commands {
		if len(commands) > 1 && cmd.After != "" && tbl.GetColumn(cmd.After) == nil {
			return errNoSuchColumn(cmd.After, tbl.Name)
		}
		cons := cmd.Constraints
		if cmd.Constraint != nil {
			cons = append(append([]*nodes.Constraint(nil), cons...), cmd.Constraint)
		}
		for _, con := range cons {
			switch con.Type {
			case nodes.ConstrPrimaryKey, nodes.ConstrUnique, nodes.ConstrIndex, nodes.ConstrForeignKey:
				if err := validateTiDBIndexColumns(tbl, buildIndexColumns(con)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func validateTiDBIndexColumns(tbl *Table, columns []*IndexColumn) error {
	for _, col := range columns {
		if col.Name != "" && tbl.GetColumn(col.Name) == nil {
			return &Error{Code: 1072, SQLState: "42000", Message: "column does not exist: " + col.Name}
		}
	}
	return nil
}

// Native CHECK validation is conditional on tidb_enable_check_constraint.
// Keep expression rejection ahead of column lookup, as buildTableInfo does.
// https://github.com/pingcap/tidb/blob/v8.5.5/pkg/table/constraint.go
// https://github.com/pingcap/tidb/blob/v8.5.5/pkg/ddl/create_table.go
func validateTiDBCheck(tbl *Table, name string, expr nodes.ExprNode, ownColumn string) error {
	tables := []*Table{tbl}
	if tbl.Database != nil {
		for _, other := range tbl.Database.Tables {
			if other != tbl {
				tables = append(tables, other)
			}
		}
	}
	for _, table := range tables {
		for _, con := range table.Constraints {
			if con.Type == ConCheck && strings.EqualFold(con.Name, name) {
				return &Error{Code: 3822, SQLState: "HY000", Message: fmt.Sprintf("Duplicate check constraint name '%s'.", name)}
			}
		}
	}
	var exprErr error
	var columns []string
	nodes.Inspect(expr, func(node nodes.Node) bool {
		if exprErr != nil {
			return false
		}
		switch n := node.(type) {
		case *nodes.VariableRef:
			exprErr = &Error{Code: 3816, SQLState: "HY000", Message: fmt.Sprintf("An expression of a check constraint '%s' cannot refer to a user or system variable.", name)}
		case *nodes.SubqueryExpr, *nodes.ExistsExpr:
			exprErr = &Error{Code: 3815, SQLState: "HY000", Message: fmt.Sprintf("An expression of a check constraint '%s' contains disallowed function.", name)}
		case *nodes.DefaultExpr:
			exprErr = tiDBCheckFunctionError(name, "default")
		case *nodes.FuncCallExpr:
			switch strings.ToLower(n.Name) {
			case "now", "current_timestamp", "curdate", "current_date", "curtime", "current_time", "localtime", "localtimestamp", "unix_timestamp", "utc_date", "utc_timestamp", "utc_time", "connection_id", "current_user", "session_user", "version", "found_rows", "last_insert_id", "system_user", "user", "rand", "row_count", "get_lock", "is_free_lock", "is_used_lock", "release_lock", "release_all_locks", "load_file", "uuid", "uuid_short", "sleep":
				exprErr = tiDBCheckFunctionError(name, strings.ToLower(n.Name))
			}
		case *nodes.ColumnRef:
			columns = append(columns, n.Column)
		}
		return exprErr == nil
	})
	if exprErr != nil {
		return exprErr
	}
	for _, column := range columns {
		if ownColumn != "" && !strings.EqualFold(ownColumn, column) {
			return &Error{Code: 3813, SQLState: "HY000", Message: fmt.Sprintf("Column check constraint '%s' references other column.", name)}
		}
		col := tbl.GetColumn(column)
		if col == nil {
			return &Error{Code: 3820, SQLState: "HY000", Message: fmt.Sprintf("Check constraint '%s' refers to non-existing column '%s'.", name, column)}
		}
		if col.AutoIncrement {
			return &Error{Code: 3818, SQLState: "HY000", Message: fmt.Sprintf("Check constraint '%s' cannot refer to an auto-increment column.", name)}
		}
	}
	return nil
}

func allocateTiDBCheckName(table string, reserved map[string]bool) string {
	for n := 1; ; n++ {
		name := fmt.Sprintf("%s_chk_%d", strings.ToLower(table), n)
		if !reserved[name] {
			reserved[name] = true
			return name
		}
	}
}

func tiDBCheckFunctionError(name, function string) error {
	return &Error{Code: 3814, SQLState: "HY000", Message: fmt.Sprintf("An expression of a check constraint '%s' contains disallowed function: %s.", name, function)}
}

// TiDB's partition expression uses its AST Restore spelling: function names are
// uppercase, while column identifiers remain quoted. Generated columns use a
// different stored form, so this does not change the general expression printer.
// https://github.com/pingcap/tidb/blob/v8.5.5/pkg/ddl/partition.go
func partitionExpressionSQL(expr nodes.ExprNode) string {
	var b strings.Builder
	var write func(nodes.ExprNode)
	write = func(expr nodes.ExprNode) {
		switch n := expr.(type) {
		case *nodes.FuncCallExpr:
			b.WriteString(strings.ToUpper(n.Name))
			b.WriteByte('(')
			for i, arg := range n.Args {
				if i > 0 {
					b.WriteString(", ")
				}
				write(arg)
			}
			b.WriteByte(')')
		case *nodes.ParenExpr:
			b.WriteByte('(')
			write(n.Expr)
			b.WriteByte(')')
		case *nodes.BinaryExpr:
			b.WriteByte('(')
			write(n.Left)
			b.WriteByte(' ')
			b.WriteString(binaryOpToString(n.Op))
			b.WriteByte(' ')
			write(n.Right)
			b.WriteByte(')')
		default:
			b.WriteString(nodeToSQL(expr))
		}
	}
	write(expr)
	return b.String()
}

func tiDBCurrentTimestamp(name string) bool {
	switch strings.ToLower(name) {
	case "now", "current_timestamp", "localtime", "localtimestamp":
		return true
	}
	return false
}

func tiDBFunctionFSP(fn *nodes.FuncCallExpr) int {
	if len(fn.Args) == 1 {
		if value, ok := fn.Args[0].(*nodes.IntLit); ok {
			return int(value.Value)
		}
	}
	return 0
}

func validateTiDBFunctionPrecision(node nodes.Node) error {
	var err error
	nodes.Inspect(node, func(node nodes.Node) bool {
		fn, ok := node.(*nodes.FuncCallExpr)
		if !ok {
			return true
		}
		switch strings.ToLower(fn.Name) {
		case "now", "current_timestamp", "localtime", "localtimestamp", "curtime", "current_time", "utc_time", "utc_timestamp", "sysdate":
			if fsp := tiDBFunctionFSP(fn); fsp > 6 {
				err = &Error{Code: 1426, SQLState: "42000", Message: fmt.Sprintf("Too big precision %d specified for '%s'. Maximum is 6.", fsp, fn.Name)}
			}
		}
		return err == nil
	})
	return err
}

func validateTiDBPrimaryGenerated(tbl *Table) error {
	for _, idx := range tbl.Indexes {
		if !idx.Primary {
			continue
		}
		for _, part := range idx.Columns {
			if col := tbl.GetColumn(part.Name); col != nil && col.Generated != nil && !col.Generated.Stored {
				return &Error{Code: 3106, SQLState: "HY000", Message: "'Defining a virtual generated column as primary key' is not supported for generated columns."}
			}
		}
	}
	return nil
}

// TiDB v8.5.5 checks the partition expression before storing the table.
// https://github.com/pingcap/tidb/blob/v8.5.5/pkg/ddl/partition.go
// https://github.com/pingcap/tidb/blob/v8.5.5/pkg/expression/function_traits.go
func validateTiDBPartition(tbl *Table, pc *nodes.PartitionClause) error {
	if pc == nil || pc.Type == nodes.PartitionSystemTime {
		return nil
	}
	if pc.Type == nodes.PartitionRange && len(pc.Columns) == 0 {
		for i, part := range pc.Partitions {
			if i < len(pc.Partitions)-1 && strings.EqualFold(partitionValueToString(part.Values, pc.Type), "MAXVALUE") {
				return &Error{Code: 1481, SQLState: "HY000", Message: "MAXVALUE can only be used in last partition definition"}
			}
		}
	}
	if pc.Expr == nil || len(pc.Columns) != 0 {
		return nil
	}
	var err error
	columns := 0
	checkOperands := func(forbidTimestamp bool, args ...nodes.ExprNode) error {
		for _, arg := range args {
			ref, ok := arg.(*nodes.ColumnRef)
			if !ok {
				continue
			}
			col := tbl.GetColumn(ref.Column)
			if col == nil {
				return errNoSuchColumn(ref.Column, "partition function")
			}
			if forbidTimestamp && strings.EqualFold(col.DataType, "timestamp") {
				return &Error{Code: 1486, SQLState: "HY000", Message: "Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed"}
			}
		}
		return nil
	}
	nodes.Inspect(pc.Expr, func(node nodes.Node) bool {
		switch n := node.(type) {
		case *nodes.FuncCallExpr:
			if err = checkOperands(false, n.Args...); err != nil {
				return false
			}
			if n.Schema != "" {
				err = &Error{Code: 1564, SQLState: "HY000", Message: "This partition function is not allowed"}
				return false
			}
			switch strings.ToLower(n.Name) {
			case "to_days", "to_seconds", "dayofmonth", "month", "dayofyear", "quarter", "yearweek", "year", "weekday", "dayofweek", "day", "hour", "minute", "second", "time_to_sec", "microsecond", "unix_timestamp", "from_days", "extract", "abs", "ceiling", "ceil", "datediff", "floor", "mod":
			default:
				err = &Error{Code: 1564, SQLState: "HY000", Message: "This partition function is not allowed"}
			}
		case *nodes.BinaryExpr:
			switch n.Op {
			case nodes.BinOpAdd, nodes.BinOpSub, nodes.BinOpMul, nodes.BinOpDivInt, nodes.BinOpMod:
				err = checkOperands(true, n.Left, n.Right)
			default:
				err = &Error{Code: 1564, SQLState: "HY000", Message: "This partition function is not allowed"}
			}
		case *nodes.UnaryExpr:
			if n.Op != nodes.UnaryPlus && n.Op != nodes.UnaryMinus {
				err = &Error{Code: 1564, SQLState: "HY000", Message: "This partition function is not allowed"}
			} else {
				err = checkOperands(true, n.Operand)
			}
		case *nodes.ColumnRef:
			columns++
			if tbl.GetColumn(n.Column) == nil {
				err = errNoSuchColumn(n.Column, "partition function")
			}
		case *nodes.VariableRef, *nodes.SubqueryExpr:
			err = &Error{Code: 1564, SQLState: "HY000", Message: "This partition function is not allowed"}
		}
		return err == nil
	})
	if err != nil {
		return err
	}
	if columns == 0 {
		return &Error{Code: 1486, SQLState: "HY000", Message: "Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed"}
	}
	expr := pc.Expr
	for {
		paren, ok := expr.(*nodes.ParenExpr)
		if !ok {
			break
		}
		expr = paren.Expr
	}
	if ref, ok := expr.(*nodes.ColumnRef); ok {
		col := tbl.GetColumn(ref.Column)
		switch strings.ToLower(col.DataType) {
		case "tinyint", "smallint", "mediumint", "int", "integer", "bigint", "year", "bit":
		default:
			return &Error{Code: 1659, SQLState: "HY000", Message: fmt.Sprintf("Field '%s' is of a not allowed type for this type of partitioning", ref.Column)}
		}
	}
	return nil
}

// Generated columns and expression indexes share TiDB's forbidden-function list.
// https://github.com/pingcap/tidb/blob/v8.5.5/pkg/expression/function_traits.go
func tiDBIllegalGeneratedExpr(expr nodes.ExprNode) bool {
	invalid := false
	nodes.Inspect(expr, func(node nodes.Node) bool {
		switch n := node.(type) {
		case *nodes.VariableRef, *nodes.SubqueryExpr:
			invalid = true
		case *nodes.FuncCallExpr:
			switch strings.ToLower(n.Name) {
			case "connection_id", "load_file", "last_insert_id", "rand", "uuid", "uuid_short", "curdate", "current_date", "curtime", "current_time", "current_timestamp", "localtime", "localtimestamp", "now", "unix_timestamp", "utc_date", "utc_time", "utc_timestamp", "benchmark", "current_user", "database", "found_rows", "get_lock", "is_free_lock", "is_used_lock", "master_pos_wait", "name_const", "release_lock", "row", "row_count", "schema", "session_user", "sleep", "sysdate", "system_user", "user", "values", "encrypt", "version", "json_merge", "release_all_locks":
				invalid = true
			}
		}
		return !invalid
	})
	return invalid
}
