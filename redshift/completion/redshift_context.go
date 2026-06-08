package completion

import "strings"

var redshiftTopLevelKeywords = []string{
	"SELECT",
	"INSERT",
	"UPDATE",
	"DELETE",
	"MERGE",
	"CREATE",
	"ALTER",
	"DROP",
	"COPY",
	"UNLOAD",
	"SHOW",
}

var redshiftCreateTableOptions = []string{
	"DISTSTYLE",
	"DISTKEY",
	"SORTKEY",
	"INTERLEAVED",
	"ENCODE",
	"BACKUP",
}

var redshiftCopyOptions = []string{
	"IAM_ROLE",
	"CREDENTIALS",
	"FORMAT",
	"DELIMITER",
	"HEADER",
	"MANIFEST",
	"ENCRYPTED",
	"GZIP",
	"BZIP2",
}

var redshiftUnloadOptions = []string{
	"IAM_ROLE",
	"FORMAT",
	"DELIMITER",
	"HEADER",
	"MANIFEST",
	"ENCRYPTED",
	"GZIP",
	"BZIP2",
}

var redshiftShowSubcommands = []string{
	"DATABASES",
	"SCHEMAS",
	"TABLES",
	"COLUMNS",
	"GRANTS",
	"DATASHARES",
}

func redshiftContextCandidates(sql string, offset int) []Candidate {
	stmt := statementBeforeCursor(sql, offset)
	stmtUpper := strings.ToUpper(stmt)

	var result []Candidate
	if stmt == "" {
		result = appendKeywords(result, redshiftTopLevelKeywords...)
	}
	if strings.HasPrefix(stmtUpper, "CREATE TABLE") {
		result = appendKeywords(result, redshiftCreateTableOptions...)
	}
	if strings.HasPrefix(stmtUpper, "COPY ") {
		result = appendKeywords(result, redshiftCopyOptions...)
	}
	if strings.HasPrefix(stmtUpper, "UNLOAD") {
		result = appendKeywords(result, redshiftUnloadOptions...)
	}
	if stmtUpper == "SHOW" {
		result = appendKeywords(result, redshiftShowSubcommands...)
	}
	return result
}

func statementBeforeCursor(sql string, offset int) string {
	if offset < 0 {
		offset = 0
	}
	if offset > len(sql) {
		offset = len(sql)
	}
	before := sql[:offset]
	if i := strings.LastIndex(before, ";"); i >= 0 {
		before = before[i+1:]
	}
	return strings.TrimSpace(before)
}

func appendKeywords(result []Candidate, keywords ...string) []Candidate {
	for _, keyword := range keywords {
		result = append(result, Candidate{Text: keyword, Type: CandidateKeyword})
	}
	return result
}
