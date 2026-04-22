// Package parser - resource_governor.go implements T-SQL Resource Governor statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateWorkloadGroupStmt parses CREATE WORKLOAD GROUP.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-workload-group-transact-sql
//
//	CREATE WORKLOAD GROUP group_name
//	[ WITH
//	    ( [ IMPORTANCE = { LOW | MEDIUM | HIGH } ]
//	      [ [ , ] REQUEST_MAX_MEMORY_GRANT_PERCENT = value ]
//	      [ [ , ] REQUEST_MAX_CPU_TIME_SEC = value ]
//	      [ [ , ] REQUEST_MEMORY_GRANT_TIMEOUT_SEC = value ]
//	      [ [ , ] MAX_DOP = value ]
//	      [ [ , ] GROUP_MAX_REQUESTS = value ]
//	      [ [ , ] GROUP_MAX_TEMPDB_DATA_MB = value ]
//	      [ [ , ] GROUP_MAX_TEMPDB_DATA_PERCENT = value ] )
//	]
//	[ USING {
//	    [ pool_name | [default] ]
//	    [ [ , ] EXTERNAL external_pool_name | [ default ] ]
//	    } ]
func (p *Parser) parseCreateWorkloadGroupStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "WORKLOAD GROUP",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// group_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options = p.parseResourceGovernorOptions()
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterWorkloadGroupStmt parses ALTER WORKLOAD GROUP.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-workload-group-transact-sql
//
//	ALTER WORKLOAD GROUP { group_name | [ default ] }
//	[ WITH
//	    ( [ IMPORTANCE = { LOW | MEDIUM | HIGH } ]
//	      [ [ , ] REQUEST_MAX_MEMORY_GRANT_PERCENT = value ]
//	      [ [ , ] REQUEST_MAX_CPU_TIME_SEC = value ]
//	      [ [ , ] REQUEST_MEMORY_GRANT_TIMEOUT_SEC = value ]
//	      [ [ , ] MAX_DOP = value ]
//	      [ [ , ] GROUP_MAX_REQUESTS = value ]
//	      [ [ , ] GROUP_MAX_TEMPDB_DATA_MB = value ]
//	      [ [ , ] GROUP_MAX_TEMPDB_DATA_PERCENT = value ] )
//	]
//	[ USING { pool_name | [default] } ]
func (p *Parser) parseAlterWorkloadGroupStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "WORKLOAD GROUP",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// group_name | [default]
	if p.isIdentLike() || p.cur.Type == tokSCONST || p.cur.Type == kwDEFAULT {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options = p.parseResourceGovernorOptions()
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropWorkloadGroupStmt parses DROP WORKLOAD GROUP.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-workload-group-transact-sql
//
//	DROP WORKLOAD GROUP group_name
func (p *Parser) parseDropWorkloadGroupStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: "WORKLOAD GROUP",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// group_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateResourcePoolStmt parses CREATE RESOURCE POOL.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-resource-pool-transact-sql
//
//	CREATE RESOURCE POOL pool_name
//	[ WITH
//	    (
//	        [ MIN_CPU_PERCENT = value ]
//	        [ [ , ] MAX_CPU_PERCENT = value ]
//	        [ [ , ] CAP_CPU_PERCENT = value ]
//	        [ [ , ] AFFINITY {SCHEDULER =
//	                  AUTO
//	                | ( <scheduler_range_spec> )
//	                | NUMANODE = ( <NUMA_node_range_spec> )
//	                } ]
//	        [ [ , ] MIN_MEMORY_PERCENT = value ]
//	        [ [ , ] MAX_MEMORY_PERCENT = value ]
//	        [ [ , ] MIN_IOPS_PER_VOLUME = value ]
//	        [ [ , ] MAX_IOPS_PER_VOLUME = value ]
//	    )
//	]
//
//	<scheduler_range_spec> ::=
//	{ SCHED_ID | SCHED_ID TO SCHED_ID }[,...n]
//
//	<NUMA_node_range_spec> ::=
//	{ NUMA_node_ID | NUMA_node_ID TO NUMA_node_ID }[,...n]
func (p *Parser) parseCreateResourcePoolStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "RESOURCE POOL",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// pool_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options = p.parseResourceGovernorOptions()
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterResourcePoolStmt parses ALTER RESOURCE POOL.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-resource-pool-transact-sql
//
//	ALTER RESOURCE POOL { pool_name | [default] }
//	[WITH
//	    ( [ MIN_CPU_PERCENT = value ]
//	    [ [ , ] MAX_CPU_PERCENT = value ]
//	    [ [ , ] CAP_CPU_PERCENT = value ]
//	    [ [ , ] AFFINITY {
//	                        SCHEDULER = AUTO
//	                      | ( <scheduler_range_spec> )
//	                      | NUMANODE = ( <NUMA_node_range_spec> )
//	                      }]
//	    [ [ , ] MIN_MEMORY_PERCENT = value ]
//	    [ [ , ] MAX_MEMORY_PERCENT = value ]
//	    [ [ , ] MIN_IOPS_PER_VOLUME = value ]
//	    [ [ , ] MAX_IOPS_PER_VOLUME = value ]
//	)]
func (p *Parser) parseAlterResourcePoolStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "RESOURCE POOL",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// pool_name | [default]
	if p.isIdentLike() || p.cur.Type == tokSCONST || p.cur.Type == kwDEFAULT {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options = p.parseResourceGovernorOptions()
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropResourcePoolStmt parses DROP RESOURCE POOL.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-resource-pool-transact-sql
//
//	DROP RESOURCE POOL pool_name
func (p *Parser) parseDropResourcePoolStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: "RESOURCE POOL",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// pool_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateExternalResourcePoolStmt parses CREATE EXTERNAL RESOURCE POOL.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-external-resource-pool-transact-sql
//
//	CREATE EXTERNAL RESOURCE POOL pool_name
//	[ WITH (
//	    [ MAX_CPU_PERCENT = value ]
//	    [ [ , ] MAX_MEMORY_PERCENT = value ]
//	    [ [ , ] MAX_PROCESSES = value ]
//	    )
//	]
func (p *Parser) parseCreateExternalResourcePoolStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "EXTERNAL RESOURCE POOL",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// pool_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options = p.parseResourceGovernorOptions()
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterExternalResourcePoolStmt parses ALTER EXTERNAL RESOURCE POOL.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-external-resource-pool-transact-sql
//
//	ALTER EXTERNAL RESOURCE POOL { pool_name | "default" }
//	[ WITH (
//	    [ MAX_CPU_PERCENT = value ]
//	    [ [ , ] MAX_MEMORY_PERCENT = value ]
//	    [ [ , ] MAX_PROCESSES = value ]
//	    )
//	]
func (p *Parser) parseAlterExternalResourcePoolStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "EXTERNAL RESOURCE POOL",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// pool_name | "default"
	if p.isIdentLike() || p.cur.Type == tokSCONST || p.cur.Type == kwDEFAULT {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options = p.parseResourceGovernorOptions()
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropExternalResourcePoolStmt parses DROP EXTERNAL RESOURCE POOL.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-external-resource-pool-transact-sql
//
//	DROP EXTERNAL RESOURCE POOL pool_name
func (p *Parser) parseDropExternalResourcePoolStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: "EXTERNAL RESOURCE POOL",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// pool_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterResourceGovernorStmt parses ALTER RESOURCE GOVERNOR.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-resource-governor-transact-sql
//
//	ALTER RESOURCE GOVERNOR
//	    { RECONFIGURE
//	      | DISABLE
//	      | RESET STATISTICS
//	      | WITH
//	              ( [ CLASSIFIER_FUNCTION = { schema_name.function_name | NULL } ]
//	                [ [ , ] MAX_OUTSTANDING_IO_PER_VOLUME = value ]
//	              )
//	    }
func (p *Parser) parseAlterResourceGovernorStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "RESOURCE GOVERNOR",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	stmt.Options = p.parseResourceGovernorOptions()
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseResourceGovernorOptions consumes the remaining tokens of a Resource Governor
// statement (WITH (...), USING ..., RECONFIGURE, DISABLE, RESET STATISTICS, etc.)
// and returns them as a list of String nodes.
//
// Used by: CREATE/ALTER WORKLOAD GROUP, CREATE/ALTER RESOURCE POOL,
//
//	CREATE/ALTER EXTERNAL RESOURCE POOL, ALTER RESOURCE GOVERNOR,
//	CREATE/ALTER WORKLOAD CLASSIFIER.
//
// WITH clause options:
//
//	WORKLOAD GROUP:     IMPORTANCE, REQUEST_MAX_MEMORY_GRANT_PERCENT, REQUEST_MAX_CPU_TIME_SEC,
//	                    REQUEST_MEMORY_GRANT_TIMEOUT_SEC, MAX_DOP, GROUP_MAX_REQUESTS,
//	                    GROUP_MAX_TEMPDB_DATA_MB, GROUP_MAX_TEMPDB_DATA_PERCENT
//	RESOURCE POOL:      MIN_CPU_PERCENT, MAX_CPU_PERCENT, CAP_CPU_PERCENT, MIN_MEMORY_PERCENT,
//	                    MAX_MEMORY_PERCENT, MIN_IOPS_PER_VOLUME, MAX_IOPS_PER_VOLUME,
//	                    AFFINITY { SCHEDULER = AUTO | (range) | NUMANODE = (range) }
//	EXTERNAL POOL:      MAX_CPU_PERCENT, MAX_MEMORY_PERCENT, MAX_PROCESSES
//	RESOURCE GOVERNOR:  CLASSIFIER_FUNCTION = schema.func | NULL, MAX_OUTSTANDING_IO_PER_VOLUME
//	WORKLOAD CLASSIFIER: WORKLOAD_GROUP, MEMBERNAME, WLM_LABEL, WLM_CONTEXT,
//	                     START_TIME, END_TIME, IMPORTANCE
//
// Non-WITH options:
//
//	USING pool_name [, EXTERNAL ext_pool] (for WORKLOAD GROUP)
//	RECONFIGURE | DISABLE | RESET STATISTICS (for ALTER RESOURCE GOVERNOR)
func (p *Parser) parseResourceGovernorOptions() *nodes.List {
	var opts []nodes.Node

	for p.cur.Type != ';' && p.cur.Type != tokEOF && p.cur.Type != kwGO {
		if p.cur.Type == kwWITH {
			p.advance()
			if p.cur.Type == '(' {
				p.advance()
				opts = append(opts, p.parseResourceGovernorWithOptions()...)
				p.match(')')
			}
		} else if p.isAnyKeywordIdent() || p.cur.Type == kwAS || p.cur.Type == kwFOR ||
			p.cur.Type == kwON || p.cur.Type == kwOFF || p.cur.Type == kwDEFAULT ||
			p.cur.Type == kwNULL {
			optLoc := p.pos()
			name := strings.ToUpper(p.cur.Str)
			p.advance()
			val := ""
			if p.cur.Type == '=' {
				p.advance()
				val = p.parseResourceGovernorValue()
			} else if (name == "USING" || name == "EXTERNAL") &&
				(p.isAnyKeywordIdent() || p.cur.Type == kwDEFAULT) {
				// USING pool_name or EXTERNAL ext_pool_name (no = sign)
				val = strings.ToUpper(p.cur.Str)
				p.advance()
			} else if name == "RESET" && p.cur.Type == kwSTATISTICS {
				// RESET STATISTICS
				val = strings.ToUpper(p.cur.Str)
				p.advance()
			}
			// Handle qualified names like schema.function
			for p.cur.Type == '.' {
				p.advance()
				if p.isIdentLike() {
					val += "." + strings.ToUpper(p.cur.Str)
					p.advance()
				}
			}
			opts = append(opts, &nodes.ResourceGovernorOption{
				Name:  name,
				Value: val,
				Loc:   nodes.Loc{Start: optLoc, End: p.prevEnd()},
			})
			if p.cur.Type == ',' {
				p.advance()
			}
		} else {
			break
		}
	}

	if len(opts) == 0 {
		return nil
	}
	return &nodes.List{Items: opts}
}

// parseResourceGovernorWithOptions parses key=value pairs inside a WITH (...) clause.
func (p *Parser) parseResourceGovernorWithOptions() []nodes.Node {
	var opts []nodes.Node

	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		if p.isAnyKeywordIdent() || p.cur.Type == kwNULL || p.cur.Type == kwON || p.cur.Type == kwOFF || p.cur.Type == kwDEFAULT {
			optName := strings.ToUpper(p.cur.Str)
			p.advance()

			// Special handling for AFFINITY (no = sign, sub-syntax follows directly)
			if optName == "AFFINITY" {
				opts = append(opts, p.parseResourceGovernorAffinityValue()...)
			} else if p.cur.Type == '=' {
				p.advance()

				if optName == "CLASSIFIER_FUNCTION" {
					// CLASSIFIER_FUNCTION = schema.func | NULL
					val := p.parseResourceGovernorQualifiedValue()
					opts = append(opts, &nodes.String{Str: optName + "=" + val})
				} else {
					val := p.parseResourceGovernorValue()
					opts = append(opts, &nodes.String{Str: optName + "=" + val})
				}
			} else {
				opts = append(opts, &nodes.String{Str: optName})
			}
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}

	return opts
}

// parseResourceGovernorValue consumes a single value: number, string, identifier, ON, OFF, NULL, DEFAULT.
func (p *Parser) parseResourceGovernorValue() string {
	switch {
	case p.cur.Type == tokICONST || p.cur.Type == tokFCONST:
		val := strings.ToUpper(p.cur.Str)
		p.advance()
		return val
	case p.cur.Type == tokSCONST:
		val := p.cur.Str
		p.advance()
		return val
	case p.cur.Type == kwON:
		p.advance()
		return "ON"
	case p.cur.Type == kwOFF:
		p.advance()
		return "OFF"
	case p.cur.Type == kwNULL:
		p.advance()
		return "NULL"
	case p.cur.Type == kwDEFAULT:
		val := strings.ToUpper(p.cur.Str)
		p.advance()
		return val
	case p.isAnyKeywordIdent():
		val := strings.ToUpper(p.cur.Str)
		p.advance()
		return val
	default:
		return ""
	}
}

// parseResourceGovernorQualifiedValue consumes a possibly dot-qualified value (schema.func or NULL).
func (p *Parser) parseResourceGovernorQualifiedValue() string {
	val := p.parseResourceGovernorValue()
	for p.cur.Type == '.' {
		p.advance()
		if p.isIdentLike() {
			val += "." + strings.ToUpper(p.cur.Str)
			p.advance()
		}
	}
	return val
}

// parseResourceGovernorAffinityValue parses the AFFINITY sub-syntax:
//
//	AFFINITY { SCHEDULER = AUTO | ( range_spec ) | NUMANODE = ( range_spec ) }
//	range_spec: { id | id TO id } [,...n]
func (p *Parser) parseResourceGovernorAffinityValue() []nodes.Node {
	var opts []nodes.Node

	if p.cur.Type == kwSCHEDULER {
		p.advance() // consume SCHEDULER
		if p.cur.Type == '=' {
			p.advance()
		}
		if p.cur.Type == kwAUTO {
			opts = append(opts, &nodes.String{Str: "AFFINITY=SCHEDULER=AUTO"})
			p.advance()
		} else if p.cur.Type == '(' {
			p.advance()
			rangeStr := p.parseResourceGovernorRangeSpec()
			p.match(')')
			opts = append(opts, &nodes.String{Str: "AFFINITY=SCHEDULER=(" + rangeStr + ")"})
		}
	} else if p.cur.Type == kwNUMANODE {
		p.advance() // consume NUMANODE
		if p.cur.Type == '=' {
			p.advance()
		}
		if p.cur.Type == '(' {
			p.advance()
			rangeStr := p.parseResourceGovernorRangeSpec()
			p.match(')')
			opts = append(opts, &nodes.String{Str: "AFFINITY=NUMANODE=(" + rangeStr + ")"})
		}
	}

	return opts
}

// parseResourceGovernorRangeSpec parses: { id | id TO id } [,...n]
func (p *Parser) parseResourceGovernorRangeSpec() string {
	var parts []string
	for {
		if p.cur.Type != tokICONST && !p.isAnyKeywordIdent() {
			break
		}
		val := p.cur.Str
		p.advance()
		if p.cur.Type == kwTO {
			p.advance()
			if p.cur.Type == tokICONST || p.isAnyKeywordIdent() {
				val += " TO " + p.cur.Str
				p.advance()
			}
		}
		parts = append(parts, val)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	return strings.Join(parts, ", ")
}

// parseCreateWorkloadClassifierStmt parses CREATE WORKLOAD CLASSIFIER.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-workload-classifier-transact-sql
//
//	CREATE WORKLOAD CLASSIFIER classifier_name
//	WITH
//	    ( WORKLOAD_GROUP = 'name'
//	    , MEMBERNAME = 'security_account'
//	    [ [ , ] WLM_LABEL = 'label' ]
//	    [ [ , ] WLM_CONTEXT = 'context' ]
//	    [ [ , ] START_TIME = 'HH:MM' ]
//	    [ [ , ] END_TIME = 'HH:MM' ]
//	    [ [ , ] IMPORTANCE = { LOW | BELOW_NORMAL | NORMAL | ABOVE_NORMAL | HIGH } ] )
func (p *Parser) parseCreateWorkloadClassifierStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "WORKLOAD CLASSIFIER",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// classifier_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options = p.parseResourceGovernorOptions()
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterWorkloadClassifierStmt parses ALTER WORKLOAD CLASSIFIER.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-workload-classifier-transact-sql
//
//	ALTER WORKLOAD CLASSIFIER classifier_name
//	WITH
//	    ( WORKLOAD_GROUP = 'name'
//	    , MEMBERNAME = 'security_account'
//	    [ [ , ] WLM_LABEL = 'label' ]
//	    [ [ , ] WLM_CONTEXT = 'context' ]
//	    [ [ , ] START_TIME = 'HH:MM' ]
//	    [ [ , ] END_TIME = 'HH:MM' ]
//	    [ [ , ] IMPORTANCE = { LOW | BELOW_NORMAL | NORMAL | ABOVE_NORMAL | HIGH } ] )
func (p *Parser) parseAlterWorkloadClassifierStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "WORKLOAD CLASSIFIER",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// classifier_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options = p.parseResourceGovernorOptions()
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropWorkloadClassifierStmt parses DROP WORKLOAD CLASSIFIER.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-workload-classifier-transact-sql
//
//	DROP WORKLOAD CLASSIFIER classifier_name
func (p *Parser) parseDropWorkloadClassifierStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: "WORKLOAD CLASSIFIER",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// classifier_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}
