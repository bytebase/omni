package catalog

import (
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// checkIndexOptions refuses the index options OceanBase refuses on a key that is not
// a vector index: WITH (...), the parameters only a vector index takes (1235);
// DATA_TABLE_ID and INDEX_TABLE_ID, which only a restore may give (1235, and 4179 in
// CREATE INDEX); and STORAGE_CACHE_POLICY, which the shared-nothing server does not
// support (1235). TestOracle_OceanBaseIndexOptionRefusals holds them to the engine.
func checkIndexOptions(opts []*nodes.IndexOption, createIndex bool) error {
	for _, opt := range opts {
		switch strings.ToUpper(opt.Name) {
		case "WITH":
			return errIndexOptionNotSupported("set index params in not vector index is")
		case "DATA_TABLE_ID", "INDEX_TABLE_ID":
			if createIndex {
				return &Error{Code: 4179, SQLState: "HY000", Message: "create index with index_table_id/data_table_id not allowed"}
			}
			return errIndexOptionNotSupported("specify " + strings.ToLower(strings.TrimSuffix(opt.Name, "_ID")) + " id not for restore purpose in create index")
		case "STORAGE_CACHE_POLICY":
			return errIndexOptionNotSupported("storage cache policy is not supported in shared storage mode")
		}
	}
	return nil
}

// checkKeyOptions runs checkIndexOptions over the keys a CREATE TABLE or an ALTER
// TABLE defines, before anything else of the statement reads them: OceanBase answers
// ADD KEY (missing) WITH (a = 1) with 1235, not 1072, and CREATE TABLE t (...,
// KEY (b) WITH (a = 1)) with 1235 when t exists. A FULLTEXT key is checked after
// its columns (1283 comes first) and a vector index takes its parameters.
func checkKeyOptions(cons []*nodes.Constraint) error {
	for _, con := range cons {
		if con == nil {
			continue
		}
		switch con.Type {
		case nodes.ConstrIndex, nodes.ConstrUnique, nodes.ConstrSpatialIndex:
			if err := checkIndexOptions(con.IndexOptions, false); err != nil {
				return err
			}
		}
	}
	return nil
}

// alterKeys is the keys an ALTER TABLE adds.
func alterKeys(cmds []*nodes.AlterTableCmd) []*nodes.Constraint {
	var cons []*nodes.Constraint
	for _, cmd := range cmds {
		if cmd.Type == nodes.ATAddIndex || cmd.Type == nodes.ATAddConstraint {
			cons = append(cons, cmd.Constraint)
		}
	}
	return cons
}

// errIndexOptionNotSupported is OceanBase's 1235, whose text is what it names and "not supported".
func errIndexOptionNotSupported(what string) error {
	return &Error{Code: 1235, SQLState: "0A000", Message: what + " not supported"}
}
