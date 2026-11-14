package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// COMPUTE POOL (Snowpark Container Services), read as an ast.ObjectStmt
// ---------------------------------------------------------------------------

func init() {
	registerObject(&objectForm{
		words:  []string{"COMPUTE", "POOL"},
		create: parseCreateComputePool,
		alter:  parseAlterComputePool,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-compute-pool
		drop:       true,
		nameParts:  1, // an account-level object
		props:      computePoolProps,
		alterProps: computePoolAlterProps,
		unset:      []string{"AUTO_SUSPEND_SECS", "AUTO_RESUME", "PLACEMENT_GROUP", "BACKUP_INSTANCE_FAMILIES", "COMMENT"},
	})
}

// computePoolProps are the properties CREATE COMPUTE POOL and ALTER COMPUTE
// POOL ... SET document. INITIALLY_SUSPENDED is CREATE's alone.
var computePoolProps = objectProps{
	"MIN_NODES":                propNumber,
	"MAX_NODES":                propNumber,
	"INSTANCE_FAMILY":          propName,
	"AUTO_RESUME":              propBool,
	"INITIALLY_SUSPENDED":      propBool,
	"AUTO_SUSPEND_SECS":        propNumber,
	"COMMENT":                  propString,
	"PLACEMENT_GROUP":          propString,
	"BACKUP_INSTANCE_FAMILIES": propStrings,
}

// computePoolAlterProps are ALTER COMPUTE POOL ... SET's: CREATE's but
// INITIALLY_SUSPENDED.
var computePoolAlterProps = func() objectProps {
	m := objectProps{}
	for k, v := range computePoolProps {
		if k != "INITIALLY_SUSPENDED" {
			m[k] = v
		}
	}
	return m
}()

// parseCreateComputePool reads what follows the name of
//
//	CREATE COMPUTE POOL [ IF NOT EXISTS ] <name>
//	  [ FOR APPLICATION <app_name> ]
//	  MIN_NODES = <num> MAX_NODES = <num> INSTANCE_FAMILY = <instance_family_name>
//	  [ AUTO_RESUME = { TRUE | FALSE } ] [ INITIALLY_SUSPENDED = { TRUE | FALSE } ]
//	  [ AUTO_SUSPEND_SECS = <num> ] [ [ WITH ] TAG ( ... ) ] [ COMMENT = '<string>' ]
//	  [ PLACEMENT_GROUP = '<name>' ] [ BACKUP_INSTANCE_FAMILIES = ( '<name>' [ , ... ] ) ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-compute-pool
func parseCreateComputePool(p *Parser, s *ast.ObjectStmt) error {
	if start := p.cur.Loc.Start; p.acceptWords("FOR", "APPLICATION") {
		name, err := p.parseObjectName()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "FOR APPLICATION", start).Names = []*ast.ObjectName{name}
	}
	props, err := p.parseObjectProperties(s, "", true)
	if err != nil {
		return err
	}
	return p.requireObjectProperties(props, "MIN_NODES", "MAX_NODES", "INSTANCE_FAMILY")
}

// parseAlterComputePool reads what follows the name of
//
//	ALTER COMPUTE POOL [ IF EXISTS ] <name> { SUSPEND | RESUME }
//	ALTER COMPUTE POOL [ IF EXISTS ] <name> STOP ALL [ OF TYPE <type> [ , ... ] ]
//	ALTER COMPUTE POOL [ IF EXISTS ] <name> SET <property> [ ... ]
//	ALTER COMPUTE POOL [ IF EXISTS ] <name> UNSET <property> [ , ... ]
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-compute-pool
func parseAlterComputePool(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	switch {
	case p.acceptWords("SUSPEND"):
		p.addObjectKeywordsClause(s, "SUSPEND", start)
		return nil
	case p.acceptWords("RESUME"):
		p.addObjectKeywordsClause(s, "RESUME", start)
		return nil
	case p.acceptWords("STOP", "ALL"):
		c := p.addObjectKeywordsClause(s, "STOP ALL", start)
		if p.acceptWords("OF", "TYPE") {
			names, err := p.parseObjectNameList()
			if err != nil {
				return err
			}
			c.Keywords = "STOP ALL OF TYPE"
			c.Names = names
			c.Loc.End = p.prev.Loc.End
		}
		return nil
	}
	ok, err := p.parseObjectAlterCommon(s, actSet|actUnset)
	if err != nil || ok {
		return err
	}
	return p.syntaxErrorAtCur()
}
