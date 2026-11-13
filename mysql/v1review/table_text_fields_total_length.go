package v1review

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// TableTextFieldsTotalLength judges each top-level CREATE TABLE and ALTER
// TABLE by the table it names as it stands in final, the schema after every
// statement of the change: a table whose columns' text lengths add up to more
// than maximum is TableTextTooLong, Name the table name as written without
// its database and Value the total in decimal. Statements nested in a
// routine body are not judged, and a table final lacks is not judged, so a
// table the change creates and then drops is never reported, and every
// statement naming a table reports the same final total. The table is looked
// up in final's unnamed schema, its name compared as final says; a
// partition's name finds its table.
//
// A column's text length comes from its type, lower-cased: 1 for char and
// binary, 255 for tinytext and tinyblob, 65,535 for text and blob, 16,777,215
// for mediumtext and mediumblob, 4,294,967,295 for longtext and longblob;
// otherwise the number in the first "name(digits)" in the type, so int(11)
// counts 11 and varchar(20) 20; otherwise 0. Each finding anchors the
// statement. Ported from Bytebase's V1 rule TABLE_TEXT_FIELDS_TOTAL_LENGTH.
//
// A nil final panics at the first judged statement, as the V1 advisor did
// with no final metadata.
func TableTextFieldsTotalLength(stmts []ast.Node, final *v1review.Schema, maximum int) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		var ref *ast.TableRef
		var loc ast.Loc
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			ref, loc = n.Table, n.Loc
		case *ast.AlterTableStmt:
			ref, loc = n.Table, n.Loc
		default:
			continue
		}
		if ref == nil || ref.Name == "" {
			continue
		}
		columns, found := final.TableColumns("", ref.Name)
		if !found {
			continue
		}
		var total int64
		for _, column := range columns {
			total += textLength(column.GetType())
		}
		if total > int64(maximum) {
			start, end := span(loc)
			findings = append(findings, v1review.Finding{Code: v1review.TableTextTooLong, Name: ref.Name, Value: strconv.FormatInt(total, 10), Start: start, End: end})
		}
	}
	return findings
}

var textLengthPattern = regexp.MustCompile(`[a-z]+\((\d+)\)`)

func textLength(columnType string) int64 {
	switch columnType = strings.ToLower(columnType); columnType {
	case "char", "binary":
		return 1
	case "tinyblob", "tinytext":
		return 255
	case "blob", "text":
		return 65_535
	case "mediumblob", "mediumtext":
		return 16_777_215
	case "longblob", "longtext":
		return 4_294_967_295
	default:
		if match := textLengthPattern.FindStringSubmatch(columnType); len(match) >= 2 {
			if n, err := strconv.ParseInt(match[1], 10, 64); err == nil {
				return n
			}
		}
	}
	return 0
}
