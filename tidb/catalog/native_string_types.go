package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// normalizeTiDBStringType follows adjustBlobTypesFlen and
// checkTooBigFieldLengthAndTryAutoConvert, preserving the caller's AST.
// https://github.com/pingcap/tidb/blob/v8.5.5/pkg/ddl/add_column.go
func normalizeTiDBStringType(def *nodes.ColumnDef, tableCharset string, strict bool) (*nodes.ColumnDef, error) {
	if def == nil || def.TypeName == nil || def.TypeName.Length == 0 {
		return def, nil
	}
	dt := def.TypeName
	name := strings.ToLower(dt.Name)
	if name != "text" && name != "blob" && name != "varchar" && name != "varbinary" {
		return def, nil
	}
	charset := dt.Charset
	if charset == "" && dt.Collate == "" {
		charset = tableCharset
	}
	cs, _, err := resolveCharsetCollation(charset, dt.Collate)
	if err != nil {
		return nil, err
	}
	binary := name == "blob" || name == "varbinary" || strings.EqualFold(cs, "binary")
	width := 1
	if !binary {
		switch strings.ToLower(cs) {
		case "utf8mb4":
			width = 4
		case "utf8", "utf8mb3":
			width = 3
		case "gbk":
			width = 2
		}
	}
	if name == "varchar" || name == "varbinary" {
		if dt.Length <= 65535/width {
			return def, nil
		}
		if strict {
			return nil, &Error{Code: 1074, SQLState: "42000", Message: fmt.Sprintf("Column length too big for column '%s' (max = %d); use BLOB or TEXT instead", def.Name, 65535/width)}
		}
	}
	bytes := uint64(dt.Length) * uint64(width)
	prefix := ""
	switch {
	case bytes <= 255:
		prefix = "tiny"
	case bytes <= 65535:
	case bytes <= 16777215:
		prefix = "medium"
	default:
		prefix = "long"
	}
	copyDef, copyType := *def, *dt
	copyType.Name = prefix + "text"
	if binary {
		copyType.Name = prefix + "blob"
	}
	copyType.Length = 0
	copyDef.TypeName = &copyType
	return &copyDef, nil
}

// Only the strict-mode property is needed by string DDL. The accepted mode
// names come from TiDB's Str2SQLMode; global changes do not change this session.
// https://github.com/pingcap/tidb/blob/v8.5.5/pkg/parser/mysql/const.go
func (c *Catalog) setTiDBSQLMode(asgn *nodes.Assignment) error {
	strict := false
	switch value := asgn.Value.(type) {
	case *nodes.StringLit:
		for _, part := range strings.Split(strings.ToUpper(value.Value), ",") {
			mode := strings.TrimSpace(part)
			switch mode {
			case "STRICT_TRANS_TABLES", "STRICT_ALL_TABLES", "TRADITIONAL":
				strict = true
			case "", "REAL_AS_FLOAT", "PIPES_AS_CONCAT", "ANSI_QUOTES", "IGNORE_SPACE", "NOT_USED", "ONLY_FULL_GROUP_BY", "NO_UNSIGNED_SUBTRACTION", "NO_DIR_IN_CREATE", "POSTGRESQL", "ORACLE", "MSSQL", "DB2", "MAXDB", "NO_KEY_OPTIONS", "NO_TABLE_OPTIONS", "NO_FIELD_OPTIONS", "MYSQL323", "MYSQL40", "ANSI", "NO_AUTO_VALUE_ON_ZERO", "NO_BACKSLASH_ESCAPES", "NO_ZERO_IN_DATE", "NO_ZERO_DATE", "INVALID_DATES", "ERROR_FOR_DIVISION_BY_ZERO", "NO_AUTO_CREATE_USER", "HIGH_NOT_PRECEDENCE", "NO_ENGINE_SUBSTITUTION", "PAD_CHAR_TO_FULL_LENGTH", "ALLOW_INVALID_DATES":
			default:
				return &Error{Code: 1105, SQLState: "HY000", Message: fmt.Sprintf("ERROR 1231 (42000): Variable 'sql_mode' can't be set to the value of '%s'", mode)}
			}
		}
	case *nodes.IntLit:
		// ModeStrictTransTables and ModeStrictAllTables in TiDB's SQLMode.
		strict = value.Value&((1<<21)|(1<<22)) != 0
	case *nodes.ColumnRef:
		if !strings.EqualFold(value.Column, "default") {
			return nil
		}
		strict = true
	default:
		return nil // Nonconstant assignments require execution-time evaluation.
	}
	if asgn.Scope != "GLOBAL" {
		c.strictDDL = strict
	}
	return nil
}
