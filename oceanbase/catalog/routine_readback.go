package catalog

import "strings"

// RoutineReadbackWithComment returns def, the text OceanBase's SHOW CREATE FUNCTION or
// SHOW CREATE PROCEDURE printed, with its COMMENT written as a literal that reads back
// as comment, the routine's own from information_schema.ROUTINES.ROUTINE_COMMENT.
//
// OceanBase prints the comment between quotes as it is, with no escapes: the comment
// it's here comes out as COMMENT 'it's here', which does not parse, and a\b as
// COMMENT 'a\b', which reads back as another comment. Its stored-program parser reads a backslash in
// a routine's COMMENT as an escape under every sql_mode, NO_BACKSLASH_ESCAPES included
// (measured on 4.3.5.6: COMMENT 'C:\path' stores C:path, COMMENT 'C:\\path' stores
// C:\path), so a backslash is always written doubled.
//
// OceanBase prints each characteristic on a line of its own, indented four spaces,
// before the body, so the first such line holding exactly comment is the clause. def
// comes back unchanged when it has none.
func RoutineReadbackWithComment(def, comment string) string {
	if comment == "" {
		return def
	}
	printed := "\n    COMMENT '" + comment + "'\n"
	i := strings.Index(def, printed)
	if i < 0 {
		return def
	}
	return def[:i] + "\n    COMMENT " + quoteStringLiteral(comment) + "\n" + def[i+len(printed):]
}
