package analysis

import "github.com/bytebase/omni/cosmosdb"

// ValidateQuery reports whether statement may run in the SQL editor's read-only
// mode, and whether it returns data: (readOnly, returnsData, error).
//
// The Cosmos DB query language has no DML or DDL; every statement the parser
// accepts is a SELECT-shaped query, so a successful parse is both verdicts and a
// syntax error is the only refusal.
func ValidateQuery(statement string) (bool, bool, error) {
	if _, err := cosmosdb.Parse(statement); err != nil {
		return false, false, err
	}
	return true, true, nil
}
