package catalog

import (
	"fmt"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// maxRecordLength is the longest row MySQL's server stores, BLOB and TEXT
// values counted by their pointers alone: HA_MAX_REC_LENGTH.
const maxRecordLength = 65535

func errTooBigFieldLength(col string, max int) error {
	return &Error{Code: 1074, SQLState: "42000",
		Message: fmt.Sprintf("Column length too big for column '%s' (max = %d); use BLOB or TEXT instead", col, max)}
}

func errTooBigRowSize() error {
	return &Error{Code: 1118, SQLState: "42000",
		Message: "Row size too large. The maximum row size for the used table type, not counting BLOBs, is 65535. This includes storage overhead, check the manual. You have to change some columns to TEXT or BLOBs"}
}

// checkRowLengths holds the table a statement leaves to the
// lengths the server stores. Under a strict sql_mode a VARCHAR or VARBINARY
// whose length takes more than 65535 bytes is 1074 (VARCHAR(16384) in
// utf8mb4, max = 16383); outside one MySQL makes it a TEXT, which the catalog
// does not model, so nothing is judged. Then the row is 1118 when its columns,
// virtual ones included, take more than 65535 bytes with one byte for each
// eight nullable columns: latin1 VARCHAR(65533) NOT NULL is the longest row a
// lone column makes, and as a nullable column VARCHAR(65532) (mysql:8.0.46,
// CREATE TABLE, ADD COLUMN and CONVERT TO alike). A column whose length the
// catalog cannot read counts as none, so a row is refused only when it
// certainly is too long.
func (c *Catalog) checkRowLengths(tbl *Table) error {
	if !c.isStrictSQLMode() {
		return nil
	}
	record, nullable := 0, 0
	for _, col := range tbl.Columns {
		if col == nil {
			continue
		}
		if col.Nullable {
			nullable++
		}
		n, err := columnRecordBytes(tbl, col)
		if err != nil {
			return err
		}
		record += n
	}
	if record+(nullable+7)/8 > maxRecordLength {
		return errTooBigRowSize()
	}
	return nil
}

// checkOversizeVarchars refuses, under a strict sql_mode, a VARCHAR a
// statement defines past 65535 bytes (1074), which the catalog otherwise
// stores as the MEDIUMTEXT MySQL makes of it outside one: VARCHAR(70000) in
// latin1 is "max = 65535" and VARCHAR(20000) in utf8mb4 "max = 16383"
// (mysql:8.0.46). A length past the bytes the column's character set allows
// below 65536 characters is the row-size check's (checkRowLengths).
func (c *Catalog) checkOversizeVarchars(defs []*nodes.ColumnDef, tableCharset string) error {
	if !c.isStrictSQLMode() {
		return nil
	}
	for _, def := range defs {
		if def == nil || def.TypeName == nil || normalizedTypeName(def.TypeName.Name) != "varchar" || def.TypeName.Length <= maxRecordLength {
			continue
		}
		charset := def.TypeName.Charset
		if charset == "" && def.TypeName.Collate != "" {
			charset = charsetForCollation(def.TypeName.Collate)
		}
		if charset == "" {
			charset = tableCharset
		}
		mb := charsetMaxBytes(normalizeCharsetName(charset))
		if def.TypeName.Binary || strings.EqualFold(charset, "binary") {
			mb = 1
		}
		return errTooBigFieldLength(def.Name, maxRecordLength/mb)
	}
	return nil
}

// columnRecordBytes is what a column takes in the server's row, MySQL's
// Field::pack_length, or 0 when the catalog cannot tell.
func columnRecordBytes(tbl *Table, col *Column) (int, error) {
	dt := strings.ToLower(col.DataType)
	length := columnTypeLength(col.ColumnType)
	switch dt {
	case "varchar", "varbinary":
		mb := 1
		if dt == "varchar" {
			charset := col.Charset
			if charset == "" {
				charset = tbl.Charset
			}
			mb = charsetMaxBytes(charset)
		}
		if length*mb > maxRecordLength {
			return 0, errTooBigFieldLength(col.Name, maxRecordLength/mb)
		}
		bytes := length * mb
		if bytes > 255 {
			return bytes + 2, nil
		}
		return bytes + 1, nil
	case "char":
		charset := col.Charset
		if charset == "" {
			charset = tbl.Charset
		}
		return length * charsetMaxBytes(charset), nil
	case "binary":
		return length, nil
	case "tinyint", "year":
		return 1, nil
	case "smallint":
		return 2, nil
	case "mediumint", "date":
		return 3, nil
	case "int", "integer", "float":
		return 4, nil
	case "bigint", "double", "real":
		return 8, nil
	case "tinyblob", "tinytext":
		return 9, nil
	case "blob", "text":
		return 10, nil
	case "mediumblob", "mediumtext":
		return 11, nil
	case "longblob", "longtext", "json":
		return 12, nil
	case "decimal", "numeric":
		return decimalRecordBytes(col.ColumnType), nil
	case "datetime":
		return 5 + (temporalFSP(col.ColumnType)+1)/2, nil
	case "timestamp":
		return 4 + (temporalFSP(col.ColumnType)+1)/2, nil
	case "time":
		return 3 + (temporalFSP(col.ColumnType)+1)/2, nil
	}
	return 0, nil
}

// decimalRecordBytes is the size of DECIMAL(p, s) in its binary form: four
// bytes for each nine digits on either side of the point, and fewer for the
// rest.
func decimalRecordBytes(columnType string) int {
	precision, scale := 10, 0
	if open := strings.IndexByte(columnType, '('); open >= 0 {
		if end := strings.IndexByte(columnType[open:], ')'); end > 0 {
			parts := strings.Split(columnType[open+1:open+end], ",")
			if p, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil {
				precision = p
			}
			if len(parts) > 1 {
				if s, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
					scale = s
				}
			}
		}
	}
	digits := [10]int{0, 1, 1, 2, 2, 3, 3, 4, 4, 4}
	intg := precision - scale
	if intg < 0 || scale < 0 {
		return 0
	}
	return intg/9*4 + digits[intg%9] + scale/9*4 + digits[scale%9]
}

// temporalFSP is the fractional seconds precision a temporal type declares.
func temporalFSP(columnType string) int {
	n := columnTypeLength(columnType)
	if n < 0 || n > 6 {
		return 0
	}
	return n
}
