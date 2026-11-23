package catalog

import (
	"fmt"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// maxRecordLength is the longest row MariaDB's server stores, BLOB and TEXT
// values counted by their pointers alone: HA_MAX_REC_LENGTH.
const maxRecordLength = 65535

// maxVarcharBytes is the most bytes one VARCHAR holds, MAX_FIELD_VARCHARLENGTH:
// mariadb:11.8.8 answers 1074 with "max = 65532" for latin1 VARCHAR(65533) and
// "max = 16383" for utf8mb4 VARCHAR(16384), NOT NULL or not, where MySQL
// takes up to 65535 bytes.
const maxVarcharBytes = 65532

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
// whose length takes more than 65532 bytes is 1074 (VARCHAR(16384) in
// utf8mb4, max = 16383); outside one MariaDB makes it a TEXT, which the catalog
// does not model, so nothing is judged. Then the row is 1118 when its columns,
// virtual ones included, take more than 65535 bytes with one byte for each
// eight nullable columns: latin1 INT, VARCHAR(65530) is 1118 and INT,
// VARCHAR(65526) runs (mariadb:11.8.8). A column whose length the
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
// stores as the MEDIUMTEXT MariaDB makes of it outside one: VARCHAR(70000) in
// latin1 is "max = 65532" and VARCHAR(20000) in utf8mb4 "max = 16383"
// (mariadb:11.8.8). A length past the bytes the column's character set allows
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
		return errTooBigFieldLength(def.Name, maxVarcharBytes/mb)
	}
	return nil
}

// columnRecordBytes is what a column takes in the server's row, the server's
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
		if length*mb > maxVarcharBytes {
			return 0, errTooBigFieldLength(col.Name, maxVarcharBytes/mb)
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

// columnTypeLength reads the length a column type declares, VARCHAR(20) is 20,
// or 0 when it declares none.
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

// autoPrefixTextKey gives a plain key's BLOB, TEXT or geometry part written
// without a length the longest prefix the engine keys, as MariaDB does with a
// 1071 note: measured on mariadb:11.8.8, KEY (t) over a utf8mb4 TEXT is stored
// as KEY (`t`(768)) on InnoDB, (`t`(250)) on MyISAM and (`t`(575)) on Aria,
// over a BLOB or a GEOMETRY as (`t`(3072)), over a TINYTEXT as (`t`(255)), the
// type's own length, and over a POINT as (`g`(25)). The key is then held to
// the engine's length whole: KEY (t, a) over a TEXT and an INT is 1071, and
// KEY (g, a) over a POINT is (`g`(25),`a`). A key that is certainly longer is
// refused, one that may fit is taken. An engine the catalog does not know the
// key length of is left to the 1170 the key-part checks answer, as are
// PRIMARY, UNIQUE, FULLTEXT and SPATIAL keys, which do not come here.
func autoPrefixTextKey(tbl *Table, idxCols []*IndexColumn) error {
	maxKey := 0
	switch toLower(tbl.Engine) {
	case "", "innodb", "default":
		maxKey = 3072
	case "myisam":
		maxKey = 1000
	case "aria":
		maxKey = 2300
	default:
		return nil
	}
	type prefixed struct {
		ic     *IndexColumn
		length int
	}
	var parts []prefixed
	bytes := 0
	for _, ic := range idxCols {
		if ic == nil || ic.Name == "" || ic.Expr != "" {
			continue
		}
		col := tbl.GetColumn(ic.Name)
		if col == nil || (!isTextBlobType(col.DataType) && !isSpatialColumnType(col.DataType)) {
			continue
		}
		// A geometry part's own prefix is kept up to the engine's longest
		// key, and a POINT's is its 25 bytes whatever it says: KEY (g(10))
		// is (`g`(10)) over a GEOMETRY and (`g`(25)) over a POINT, and KEY
		// (g(5000)) over a POLYGON is (`g`(3072)).
		if ic.Length > 0 {
			if isSpatialColumnType(col.DataType) {
				length := min(ic.Length, maxKey)
				if strings.EqualFold(col.DataType, "point") {
					length = 25
				}
				parts = append(parts, prefixed{ic, length})
				bytes += length
			}
			continue
		}
		mb := 1
		if isStringType(col.DataType) {
			charset := col.Charset
			if charset == "" {
				charset = tbl.Charset
			}
			mb = charsetMaxBytes(charset)
		}
		chars := 0
		switch toLower(col.DataType) {
		case "point":
			chars = 25
		case "tinytext", "tinyblob":
			chars = 255
		case "text", "blob":
			chars = 65535
		default:
			chars = 16777215
		}
		length := min(maxKey, chars*mb) / mb
		parts = append(parts, prefixed{ic, length})
		bytes += length * mb
	}
	// Every other part takes a byte at least.
	if len(parts) > 0 && bytes+len(idxCols)-len(parts) > maxKey {
		return &Error{Code: 1071, SQLState: "42000", Message: fmt.Sprintf("Specified key was too long; max key length is %d bytes", maxKey)}
	}
	for _, p := range parts {
		p.ic.Length = p.length
	}
	return nil
}
