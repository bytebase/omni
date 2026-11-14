// Package parser - resource_governor.go implements T-SQL Resource Governor statement parsing.
package parser

import (
	"fmt"
	"math"
	"strconv"
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

	if err := p.rgWorkloadGroupBody(stmt, false); err != nil {
		return nil, err
	}
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

	// group_name | [default]: SQL Server takes the default group's name only
	// bracketed, and requires a WITH or a USING clause (Msg 10917).
	if err := p.rgWorkloadGroupBody(stmt, true); err != nil {
		return nil, err
	}
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

	name, err := p.rgName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

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

	if err := p.rgPoolBody(stmt, false); err != nil {
		return nil, err
	}
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

	// pool_name | [default]: SQL Server takes the default pool's name only
	// bracketed, and requires the WITH clause.
	if err := p.rgPoolBody(stmt, true); err != nil {
		return nil, err
	}
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

	name, err := p.rgName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

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

	if err := p.parseExternalResourcePoolBody(stmt, false); err != nil {
		return nil, err
	}
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

	// pool_name | [default]: SQL Server takes the default pool's name only
	// bracketed, and requires the WITH clause.
	if err := p.parseExternalResourcePoolBody(stmt, true); err != nil {
		return nil, err
	}
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

	name, err := p.extName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

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

	opt, err := p.rgGovernorAction()
	if err != nil {
		return nil, err
	}
	stmt.Options = &nodes.List{Items: []nodes.Node{opt}}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseResourceGovernorOptions consumes the remaining tokens of a Resource Governor
// statement (WITH (...), USING ..., RECONFIGURE, DISABLE, RESET STATISTICS, etc.)
// and returns them as a list of String nodes.
//
// Used by: CREATE/ALTER EXTERNAL RESOURCE POOL. The other Resource Governor
// statements read their clauses by the grammar (rgPoolBody,
// rgWorkloadGroupBody, rgGovernorAction).
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
		} else if p.isKeywordOrIdent() || p.cur.Type == kwAS || p.cur.Type == kwFOR ||
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
				(p.isIdentLike() || p.cur.Type == kwDEFAULT) {
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
		if p.isIdentLike() || p.cur.Type == kwNULL || p.cur.Type == kwON || p.cur.Type == kwOFF || p.cur.Type == kwDEFAULT {
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
	case p.isIdentLike():
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
		if p.cur.Type != tokICONST && !p.isIdentLike() {
			break
		}
		val := p.cur.Str
		p.advance()
		if p.cur.Type == kwTO {
			p.advance()
			if p.cur.Type == tokICONST || p.isIdentLike() {
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

// parseCreateWorkloadClassifierStmt parses CREATE WORKLOAD CLASSIFIER, which
// SQL Server 2022 refuses once it has read the name: "CREATE WORKLOAD
// CLASSIFIER is allowed only when connected to Synapse frontend" (Msg 22487).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-workload-classifier-transact-sql
// (Azure Synapse Analytics only); TSql160.g:22634.
func (p *Parser) parseCreateWorkloadClassifierStmt() (*nodes.SecurityStmt, error) {
	return nil, p.rgSynapseOnlyClassifier("CREATE")
}

// parseAlterWorkloadClassifierStmt parses ALTER WORKLOAD CLASSIFIER, which
// SQL Server 2022 refuses once it has read the name (Msg 22487), as it does
// CREATE WORKLOAD CLASSIFIER. ScriptDOM has no ALTER form.
func (p *Parser) parseAlterWorkloadClassifierStmt() (*nodes.SecurityStmt, error) {
	return nil, p.rgSynapseOnlyClassifier("ALTER")
}

// parseDropWorkloadClassifierStmt parses DROP WORKLOAD CLASSIFIER, which SQL
// Server 2022 refuses once it has read the name (Msg 22487).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-workload-classifier-transact-sql
// (Azure Synapse Analytics only); TSql160.g:15730.
func (p *Parser) parseDropWorkloadClassifierStmt() (*nodes.SecurityStmt, error) {
	return nil, p.rgSynapseOnlyClassifier("DROP")
}

// rgSynapseOnlyClassifier reads the classifier's name and refuses the
// statement as SQL Server 2022 does: a syntax error before the name, and
// Msg 22487 after it.
func (p *Parser) rgSynapseOnlyClassifier(action string) error {
	loc := p.pos()
	if _, err := p.rgName(); err != nil {
		return err
	}
	return p.newParseError(loc, action+" WORKLOAD CLASSIFIER is allowed only in Azure Synapse Analytics")
}

// rgIsWord reports whether the current token is word, written without
// brackets or quotes: SQL Server refuses [USING], [AFFINITY], [DISABLE] and
// the bracketed option names of these statements.
func (p *Parser) rgIsWord(word string) bool {
	if kw := lookupKeyword(word); kw != tokIDENT {
		return p.cur.Type == kw
	}
	return p.cur.Type == tokIDENT && !p.isDelimited(p.cur) && strings.ToUpper(p.cur.Str) == word
}

// rgName reads the name of a pool, a workload group or a classifier: an
// identifier, a reserved word only when delimited (TSql160.g identifier), so
// [default] but not default or 'a'.
func (p *Parser) rgName() (string, error) {
	if !p.isIdentLike() {
		return "", p.syntaxErrorAtCur()
	}
	name := p.cur.Str
	p.advance()
	return name, nil
}

// rgInt32 reads an unsigned decimal integer that fits an Int32, as
// TSql160.g's integer rule and the Int32.TryParse in its actions do: SQL
// Server refuses a sign, a hex literal, a decimal and 2147483648.
func (p *Parser) rgInt32() (int64, error) {
	if p.cur.Type != tokICONST || strings.HasPrefix(strings.ToLower(p.cur.Str), "0x") || p.cur.Ival > math.MaxInt32 {
		return 0, p.syntaxErrorAtCur()
	}
	v := p.cur.Ival
	p.advance()
	return v, nil
}

// rgPoolParam is a resource pool parameter and the range SQL Server takes for
// it (TSql160.g:22150-22200, Msg 1080 outside it).
type rgPoolParam struct {
	word     string
	min, max int64
}

// rgPoolParams are the resource pool parameters SQL Server 2022 takes besides
// AFFINITY: ScriptDOM's ResourcePoolParameterType, each probed on the pinned
// server.
var rgPoolParams = []rgPoolParam{
	{"MIN_CPU_PERCENT", 0, 100},
	{"MAX_CPU_PERCENT", 1, 100},
	{"CAP_CPU_PERCENT", 1, 100},
	{"MIN_MEMORY_PERCENT", 0, 100},
	{"MAX_MEMORY_PERCENT", 1, 100},
	{"TARGET_MEMORY_PERCENT", 1, 100},
	{"MIN_IO_PERCENT", 0, 100},
	{"MAX_IO_PERCENT", 1, 100},
	{"CAP_IO_PERCENT", 1, 100},
	{"MIN_IOPS_PER_VOLUME", 0, math.MaxInt32},
	{"MAX_IOPS_PER_VOLUME", 0, math.MaxInt32},
}

// rgPoolBody reads the name and the options of CREATE or ALTER RESOURCE POOL
// (TSql160.g:22112, 22127-22250):
//
//	pool_name [ WITH ( <parameter> [ ,...n ] ) ]
//
//	<parameter> ::= name = integer
//	  | AFFINITY { SCHEDULER = { AUTO | ( <range> [ ,...n ] ) } | NUMANODE = ( <range> [ ,...n ] ) }
//	<range> ::= integer [ TO integer ]
//
// SQL Server 2022 requires the WITH clause in ALTER and takes each parameter,
// AFFINITY included, once (Msg 1039).
func (p *Parser) rgPoolBody(stmt *nodes.SecurityStmt, requireWith bool) error {
	name, err := p.rgName()
	if err != nil {
		return err
	}
	stmt.Name = name
	if p.cur.Type != kwWITH {
		if requireWith {
			return p.syntaxErrorAtCur()
		}
		return nil
	}
	p.advance()
	if _, err := p.expect('('); err != nil {
		return err
	}
	var opts []nodes.Node
	seen := map[string]bool{}
	for {
		optLoc := p.pos()
		var key, str string
		if p.rgIsWord("AFFINITY") {
			p.advance()
			key = "AFFINITY"
			if str, err = p.rgPoolAffinity(); err != nil {
				return err
			}
		} else {
			var param *rgPoolParam
			for i := range rgPoolParams {
				if p.rgIsWord(rgPoolParams[i].word) {
					param = &rgPoolParams[i]
					break
				}
			}
			if param == nil {
				return p.syntaxErrorAtCur()
			}
			key = param.word
			p.advance()
			if _, err := p.expect('='); err != nil {
				return err
			}
			valLoc := p.pos()
			v, err := p.rgInt32()
			if err != nil {
				return err
			}
			if v < param.min || v > param.max {
				return p.newParseError(valLoc, fmt.Sprintf("the integer value %s is out of range", key))
			}
			str = key + "=" + strconv.FormatInt(v, 10)
		}
		if seen[key] {
			return p.newParseError(optLoc, fmt.Sprintf("option '%s' is specified more than once", key))
		}
		seen[key] = true
		opts = append(opts, &nodes.String{Str: str})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return err
	}
	stmt.Options = &nodes.List{Items: opts}
	return nil
}

// rgPoolAffinity reads what follows AFFINITY in a resource pool
// (TSql160.g:22202): SCHEDULER = { AUTO | ( ranges ) } or NUMANODE =
// ( ranges ). It returns the option as AFFINITY=SCHEDULER=AUTO or
// AFFINITY=NUMANODE=(0 TO 1, 3).
func (p *Parser) rgPoolAffinity() (string, error) {
	var kind string
	switch {
	case p.rgIsWord("SCHEDULER"):
		kind = "SCHEDULER"
	case p.rgIsWord("NUMANODE"):
		kind = "NUMANODE"
	default:
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	if _, err := p.expect('='); err != nil {
		return "", err
	}
	if kind == "SCHEDULER" && p.rgIsWord("AUTO") {
		p.advance()
		return "AFFINITY=SCHEDULER=AUTO", nil
	}
	if _, err := p.expect('('); err != nil {
		return "", err
	}
	var ranges []string
	for {
		lo, err := p.rgInt32()
		if err != nil {
			return "", err
		}
		r := strconv.FormatInt(lo, 10)
		if p.cur.Type == kwTO {
			p.advance()
			hi, err := p.rgInt32()
			if err != nil {
				return "", err
			}
			r += " TO " + strconv.FormatInt(hi, 10)
		}
		ranges = append(ranges, r)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return "", err
	}
	return "AFFINITY=" + kind + "=(" + strings.Join(ranges, ", ") + ")", nil
}

// rgWorkloadGroupBody reads the name, the WITH options and the USING clause of
// CREATE or ALTER WORKLOAD GROUP (TSql160.g:22467-22600):
//
//	group_name [ WITH ( <parameter> [ ,...n ] ) ]
//	[ USING { pool_name | EXTERNAL external_pool_name } [ ,...n ] ]
//
//	<parameter> ::= IMPORTANCE = { LOW | MEDIUM | HIGH }
//	  | REQUEST_MAX_MEMORY_GRANT_PERCENT = 0..100, an integer or a decimal
//	  | { REQUEST_MAX_CPU_TIME_SEC | REQUEST_MEMORY_GRANT_TIMEOUT_SEC
//	    | MAX_DOP | GROUP_MAX_REQUESTS } = integer
//
// SQL Server 2022 takes each parameter, the pool and the external pool once
// (Msg 1039), refuses GROUP_MAX_TEMPDB_DATA_MB, GROUP_MAX_TEMPDB_DATA_PERCENT
// and the Synapse parameters, and requires WITH or USING in ALTER (Msg 10917).
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-workload-group-transact-sql
func (p *Parser) rgWorkloadGroupBody(stmt *nodes.SecurityStmt, isAlter bool) error {
	name, err := p.rgName()
	if err != nil {
		return err
	}
	stmt.Name = name
	var opts []nodes.Node
	hasClause := false
	if p.cur.Type == kwWITH {
		hasClause = true
		p.advance()
		if _, err := p.expect('('); err != nil {
			return err
		}
		seen := map[string]bool{}
		for {
			optLoc := p.pos()
			key, opt, err := p.rgWorkloadGroupParam()
			if err != nil {
				return err
			}
			if seen[key] {
				return p.newParseError(optLoc, fmt.Sprintf("option '%s' is specified more than once", key))
			}
			seen[key] = true
			opts = append(opts, &nodes.String{Str: opt})
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if _, err := p.expect(')'); err != nil {
			return err
		}
	}
	if p.rgIsWord("USING") {
		hasClause = true
		p.advance()
		seen := map[string]bool{}
		for {
			optLoc := p.pos()
			key := "USING"
			if p.rgIsWord("EXTERNAL") {
				key = "EXTERNAL"
				p.advance()
			}
			if seen[key] {
				return p.newParseError(optLoc, "a pool is specified more than once")
			}
			seen[key] = true
			pool, err := p.rgName()
			if err != nil {
				return err
			}
			opts = append(opts, &nodes.ResourceGovernorOption{
				Name:  key,
				Value: strings.ToUpper(pool),
				Loc:   nodes.Loc{Start: optLoc, End: p.prevEnd()},
			})
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}
	if isAlter && !hasClause {
		return p.newParseError(p.pos(), "ALTER WORKLOAD GROUP requires a WITH or a USING clause")
	}
	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}
	return nil
}

// rgWorkloadGroupParam reads one workload group parameter and returns its name
// and NAME=VALUE. The value forms and ranges are what SQL Server 2022 takes
// (TSql160.g:22544): IMPORTANCE takes a word, REQUEST_MAX_MEMORY_GRANT_PERCENT
// a number from 0 to 100 with no exponent (Msg 1031 outside it), the rest an
// Int32.
func (p *Parser) rgWorkloadGroupParam() (string, string, error) {
	switch {
	case p.rgIsWord("IMPORTANCE"):
		p.advance()
		if _, err := p.expect('='); err != nil {
			return "", "", err
		}
		for _, w := range []string{"LOW", "MEDIUM", "HIGH"} {
			if p.rgIsWord(w) {
				p.advance()
				return "IMPORTANCE", "IMPORTANCE=" + w, nil
			}
		}
		return "", "", p.syntaxErrorAtCur()
	case p.rgIsWord("REQUEST_MAX_MEMORY_GRANT_PERCENT"):
		const key = "REQUEST_MAX_MEMORY_GRANT_PERCENT"
		p.advance()
		if _, err := p.expect('='); err != nil {
			return "", "", err
		}
		text := p.cur.Str
		switch {
		case p.cur.Type == tokICONST && !strings.HasPrefix(strings.ToLower(text), "0x"):
		case p.cur.Type == tokFCONST && !strings.ContainsAny(text, "eExX"):
		default:
			return "", "", p.syntaxErrorAtCur()
		}
		if v, err := strconv.ParseFloat(text, 64); err != nil || v < 0 || v > 100 {
			return "", "", p.newParseError(p.pos(), "percent values must be between 0 and 100")
		}
		p.advance()
		return key, key + "=" + text, nil
	}
	for _, w := range []string{"REQUEST_MAX_CPU_TIME_SEC", "REQUEST_MEMORY_GRANT_TIMEOUT_SEC", "MAX_DOP", "GROUP_MAX_REQUESTS"} {
		if !p.rgIsWord(w) {
			continue
		}
		p.advance()
		if _, err := p.expect('='); err != nil {
			return "", "", err
		}
		v, err := p.rgInt32()
		if err != nil {
			return "", "", err
		}
		return w, w + "=" + strconv.FormatInt(v, 10), nil
	}
	return "", "", p.syntaxErrorAtCur()
}

// rgGovernorAction reads what follows ALTER RESOURCE GOVERNOR
// (TSql160.g:22060):
//
//	DISABLE | RECONFIGURE | RESET STATISTICS
//	| WITH ( CLASSIFIER_FUNCTION = { schema_name.function_name | NULL } )
//	| WITH ( MAX_OUTSTANDING_IO_PER_VOLUME = 1..100 )
//
// SQL Server 2022 takes one option in WITH, and MAX_OUTSTANDING_IO_PER_VOLUME
// (which ScriptDOM lacks) from 1 to 100 (Msg 1040).
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-resource-governor-transact-sql
func (p *Parser) rgGovernorAction() (nodes.Node, error) {
	loc := p.pos()
	switch {
	case p.rgIsWord("DISABLE"), p.rgIsWord("RECONFIGURE"):
		name := strings.ToUpper(p.cur.Str)
		p.advance()
		return &nodes.ResourceGovernorOption{Name: name, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
	case p.rgIsWord("RESET"):
		p.advance()
		if !p.rgIsWord("STATISTICS") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		return &nodes.ResourceGovernorOption{Name: "RESET", Value: "STATISTICS", Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
	case p.cur.Type != kwWITH:
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var str string
	switch {
	case p.rgIsWord("CLASSIFIER_FUNCTION"):
		p.advance()
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		if p.cur.Type == kwNULL {
			p.advance()
			str = "CLASSIFIER_FUNCTION=NULL"
			break
		}
		// schemaObjectNonEmptyTwoPartName: schema_name.function_name.
		schema, err := p.rgName()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect('.'); err != nil {
			return nil, err
		}
		fn, err := p.rgName()
		if err != nil {
			return nil, err
		}
		str = "CLASSIFIER_FUNCTION=" + schema + "." + fn
	case p.rgIsWord("MAX_OUTSTANDING_IO_PER_VOLUME"):
		p.advance()
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		valLoc := p.pos()
		v, err := p.rgInt32()
		if err != nil {
			return nil, err
		}
		if v < 1 || v > 100 {
			return nil, p.newParseError(valLoc, "MAX_OUTSTANDING_IO_PER_VOLUME must be greater than 0 and at most 100")
		}
		str = "MAX_OUTSTANDING_IO_PER_VOLUME=" + strconv.FormatInt(v, 10)
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.String{Str: str}, nil
}
