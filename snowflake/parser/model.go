package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// MODEL and MODEL MONITOR (Snowflake ML), read as an ast.ObjectStmt
// ---------------------------------------------------------------------------
//
// A model's FROM @<stage> form (CREATE MODEL, ALTER MODEL ... ADD VERSION) is
// not read yet: ObjectClause has no field for a stage location.

func init() {
	registerObject(&objectForm{
		words:     []string{"MODEL"},
		create:    parseCreateModel,
		alter:     parseAlterModel,
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-model writes
		// DROP MODEL <name> with no IF EXISTS, which the registry's DROP always
		// takes, so DROP MODEL is not registered.
		props: objectProps{"COMMENT": propString, "DEFAULT_VERSION": propString},
		unset: []string{},
	})
	registerObject(&objectForm{
		words:     []string{"MODEL", "MONITOR"},
		create:    parseCreateModelMonitor,
		alter:     parseAlterModelMonitor,
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-model-monitor
		drop: true,
		// ALTER ... SET's properties; CREATE's WITH run is held to
		// modelMonitorCreateProps.
		props: objectProps{
			"REFRESH_INTERVAL": propString,
			"WAREHOUSE":        propName,
			"COMMENT":          propString,
			"BASELINE":         propString,
		},
		unset: []string{"COMMENT", "BASELINE"},
	})
}

// The properties a model version's SET takes: ALIAS after VERSION, COMMENT
// and METADATA after MODIFY VERSION.
var (
	modelVersionAliasProps  = objectProps{"ALIAS": propName}
	modelVersionModifyProps = objectProps{"COMMENT": propString, "METADATA": propString}
)

// modelMonitorCreateProps are the properties of CREATE MODEL MONITOR's WITH.
var modelMonitorCreateProps = objectProps{
	"MODEL":                    propName,
	"VERSION":                  propString,
	"GATEWAY":                  propName,
	"FUNCTION":                 propString,
	"SOURCE":                   propName,
	"WAREHOUSE":                propName,
	"REFRESH_INTERVAL":         propString,
	"AGGREGATION_WINDOW":       propString,
	"TIMESTAMP_COLUMN":         propName,
	"BASELINE":                 propName,
	"GROUND_TRUTH":             propName,
	"ID_COLUMNS":               propStrings,
	"PREDICTION_CLASS_COLUMNS": propStrings,
	"PREDICTION_SCORE_COLUMNS": propStrings,
	"ACTUAL_CLASS_COLUMNS":     propStrings,
	"ACTUAL_SCORE_COLUMNS":     propStrings,
	"SEGMENT_COLUMNS":          propStrings,
	"CUSTOM_METRIC_COLUMNS":    propStrings,
	"COMMENT":                  propString,
}

// parseCreateModel reads what follows the name of
//
//	CREATE [ OR REPLACE ] MODEL [ IF NOT EXISTS ] <name> [ WITH VERSION <version_name> ]
//	  FROM MODEL <source_model_name> [ VERSION <version_name> ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-model
func parseCreateModel(p *Parser, s *ast.ObjectStmt) error {
	if start := p.cur.Loc.Start; p.acceptWords("WITH", "VERSION") {
		if err := p.parseAiNameClause(s, "WITH VERSION", start); err != nil {
			return err
		}
	}
	return p.parseModelFromModel(s)
}

// parseModelFromModel reads FROM MODEL <model_name> [ VERSION <version_name> ]
// into a FROM MODEL clause and, when given, a VERSION clause.
func (p *Parser) parseModelFromModel(s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	if _, err := p.expect(kwFROM); err != nil {
		return err
	}
	if err := p.expectWords("MODEL"); err != nil {
		return err
	}
	name, err := p.parseObjectName()
	if err != nil {
		return err
	}
	p.addObjectKeywordsClause(s, "FROM MODEL", start).Names = []*ast.ObjectName{name}
	if start := p.cur.Loc.Start; p.acceptWords("VERSION") {
		return p.parseAiNameClause(s, "VERSION", start)
	}
	return nil
}

// parseAlterModel reads what follows the name of
//
//	ALTER MODEL [ IF EXISTS ] <name> SET
//	  [ COMMENT = '<string_literal>' ] [ DEFAULT_VERSION = '<version_name>' ]
//	ALTER MODEL [ IF EXISTS ] <name> SET TAG <tag_name> = '<tag_value>' [ , ... ]
//	ALTER MODEL [ IF EXISTS ] <name> UNSET TAG <tag_name> [ , ... ]
//	ALTER MODEL [ IF EXISTS ] <name> VERSION <version_or_alias_name> SET ALIAS = <alias_name>
//	ALTER MODEL [ IF EXISTS ] <name> VERSION <version_or_alias_name> UNSET ALIAS
//	ALTER MODEL <name> RENAME TO <new_name>
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-model
//
//	ALTER MODEL [ IF EXISTS ] <name> ADD VERSION <version_name>
//	  FROM MODEL <source_model_name> [ VERSION <source_version_or_alias_name> ]
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-model-add-version
//
//	ALTER MODEL [ IF EXISTS ] <name> DROP VERSION <version_name>
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-model-drop-version
//
//	ALTER MODEL [ IF EXISTS ] <name> MODIFY VERSION <version_name> SET
//	  [ COMMENT = '<string_literal>' ] [ METADATA = '<json_string>' ]
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-model-modify-version
func parseAlterModel(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	switch {
	case p.acceptWords("VERSION"):
		// VERSION <v> { SET ALIAS = <alias> | UNSET ALIAS }, one clause.
		version, err := p.parseAiSingleName()
		if err != nil {
			return err
		}
		switch {
		case p.acceptWords("SET"):
			if !p.curIsWord("ALIAS") || !p.startsObjectProperty() {
				return p.syntaxErrorAtCur()
			}
			alias, err := p.parseObjectProperty()
			if err != nil {
				return err
			}
			if err := aiCheckOptions([]*ast.CopyOption{alias}, modelVersionAliasProps); err != nil {
				return err
			}
			c := p.addObjectKeywordsClause(s, "VERSION SET", start)
			c.Names, c.Options = []*ast.ObjectName{version}, []*ast.CopyOption{alias}
			return nil
		case p.acceptWords("UNSET", "ALIAS"):
			c := p.addObjectKeywordsClause(s, "VERSION UNSET", start)
			c.Names, c.Keys = []*ast.ObjectName{version}, []string{"ALIAS"}
			return nil
		}
	case p.acceptWords("ADD", "VERSION"):
		if err := p.parseAiNameClause(s, "ADD VERSION", start); err != nil {
			return err
		}
		return p.parseModelFromModel(s)
	case p.acceptWords("DROP", "VERSION"):
		return p.parseAiNameClause(s, "DROP VERSION", start)
	case p.acceptWords("MODIFY", "VERSION"):
		// MODIFY VERSION <v> SET <property> [ ... ], one clause.
		version, err := p.parseAiSingleName()
		if err != nil {
			return err
		}
		if _, err := p.expect(kwSET); err != nil {
			return err
		}
		var opts []*ast.CopyOption
		for p.startsObjectProperty() {
			opt, err := p.parseObjectProperty()
			if err != nil {
				return err
			}
			opts = append(opts, opt)
		}
		if len(opts) == 0 {
			return p.syntaxErrorAtCur()
		}
		if err := aiCheckOptions(opts, modelVersionModifyProps); err != nil {
			return err
		}
		c := p.addObjectKeywordsClause(s, "MODIFY VERSION SET", start)
		c.Names, c.Options = []*ast.ObjectName{version}, opts
		return nil
	default:
		allow := actSet | actTags
		if !s.IfExists {
			allow |= actRename
		}
		if ok, err := p.parseObjectAlterCommon(s, allow); ok || err != nil {
			return err
		}
	}
	return p.syntaxErrorAtCur()
}

// parseCreateModelMonitor reads what follows the name of
//
//	CREATE [ OR REPLACE ] MODEL MONITOR [ IF NOT EXISTS ] <monitor_name> WITH
//	  MODEL = <model_name> VERSION = '<version_name>' FUNCTION = '<function_name>'
//	  SOURCE = <source_name> WAREHOUSE = <warehouse_name>
//	  REFRESH_INTERVAL = '<refresh_interval>' AGGREGATION_WINDOW = '<aggregation_window>'
//	  TIMESTAMP_COLUMN = <timestamp_name>
//	  [ BASELINE = <baseline_name> ]
//	  [ ID_COLUMNS = <id_column_name_array> ]
//	  [ PREDICTION_CLASS_COLUMNS = <prediction_class_column_name_array> ]
//	  [ PREDICTION_SCORE_COLUMNS = <prediction_column_name_array> ]
//	  [ ACTUAL_CLASS_COLUMNS = <actual_class_column_name_array> ]
//	  [ ACTUAL_SCORE_COLUMNS = <actual_column_name_array> ]
//	  [ SEGMENT_COLUMNS = <segment_column_name_array> ]
//	  [ CUSTOM_METRIC_COLUMNS = <custom_metric_column_name_array> ]
//	  [ COMMENT = '<comment>' ]
//
//	CREATE [ OR REPLACE ] MODEL MONITOR [ IF NOT EXISTS ] <monitor_name> WITH
//	  MODEL = <model_name> GATEWAY = <gateway_name> FUNCTION = '<function_name>'
//	  WAREHOUSE = <warehouse_name>
//	  REFRESH_INTERVAL = '<refresh_interval>' AGGREGATION_WINDOW = '<aggregation_window>'
//	  [ GROUND_TRUTH = <ground_truth_table> ] [ ... the column arrays above ... ]
//	  [ COMMENT = '<comment>' ]
//
// The properties are one WITH clause. A monitor on a gateway needs no VERSION,
// SOURCE or TIMESTAMP_COLUMN.
//
// https://docs.snowflake.com/en/sql-reference/sql/create-model-monitor
func parseCreateModelMonitor(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	if _, err := p.expect(kwWITH); err != nil {
		return err
	}
	props, err := p.parseObjectProperties(s, "WITH", true)
	if err != nil {
		return err
	}
	props.Loc.Start = start
	if err := aiCheckOptions(props.Options, modelMonitorCreateProps); err != nil {
		return err
	}
	if err := p.requireObjectProperties(props, "MODEL", "FUNCTION", "WAREHOUSE", "REFRESH_INTERVAL", "AGGREGATION_WINDOW"); err != nil {
		return err
	}
	for _, o := range props.Options {
		if o.Name == "GATEWAY" {
			return nil
		}
	}
	return p.requireObjectProperties(props, "VERSION", "SOURCE", "TIMESTAMP_COLUMN")
}

// parseAlterModelMonitor reads what follows the name of
//
//	ALTER MODEL MONITOR [ IF EXISTS ] <monitor_name> { SUSPEND | RESUME }
//	ALTER MODEL MONITOR [ IF EXISTS ] <monitor_name> SET
//	  [ REFRESH_INTERVAL = '<refresh_interval>' ] [ WAREHOUSE = <warehouse_name> ]
//	  [ COMMENT = '<comment>' ]
//	ALTER MODEL MONITOR [ IF EXISTS ] <monitor_name> UNSET COMMENT
//	ALTER MODEL MONITOR [ IF EXISTS ] <monitor_name> SET BASELINE = '<baseline_table_name>'
//	ALTER MODEL MONITOR [ IF EXISTS ] <monitor_name> UNSET BASELINE
//	ALTER MODEL MONITOR [ IF EXISTS ] <monitor_name> ADD SEGMENT_COLUMN = '<segment_column>'
//	ALTER MODEL MONITOR [ IF EXISTS ] <monitor_name> DROP SEGMENT_COLUMN = '<segment_column>'
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-model-monitor
func parseAlterModelMonitor(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	switch {
	case p.acceptWords("SUSPEND"):
		p.addObjectKeywordsClause(s, "SUSPEND", start)
		return nil
	case p.acceptWords("RESUME"):
		p.addObjectKeywordsClause(s, "RESUME", start)
		return nil
	case p.atWords("ADD", "SEGMENT_COLUMN") || p.atWords("DROP", "SEGMENT_COLUMN"):
		verb := p.advance()
		if !p.startsObjectProperty() {
			return p.syntaxErrorAtCur()
		}
		col, err := p.parseObjectProperty()
		if err != nil {
			return err
		}
		if col.Lit == nil || col.Lit.Kind != ast.LitString {
			return &ParseError{Position: col.Loc.Start, End: col.Loc.End, Message: "syntax error: a segment column is a string"}
		}
		keywords := "ADD"
		if verb.Type == kwDROP {
			keywords = "DROP"
		}
		p.addObjectKeywordsClause(s, keywords, start).Options = []*ast.CopyOption{col}
		return nil
	}
	if ok, err := p.parseObjectAlterCommon(s, actSet|actUnset); ok || err != nil {
		return err
	}
	return p.syntaxErrorAtCur()
}

// startsShowModelScope reports the ML scope word at cur (after SHOW's IN),
// MODEL, DATASET, EXPERIMENT or GATEWAY, when a name follows it, or "". A
// scope word with no name after it is still a schema's name. It consumes
// nothing. The scopes come from
//
//	SHOW VERSIONS [ LIKE '<pattern>' ] IN MODEL <model_name>
//	SHOW FUNCTIONS [ LIKE '<pattern>' ] IN MODEL <model_name>
//	SHOW VERSIONS [ LIKE '<pattern>' ] IN DATASET <dataset_name> [ LIMIT <rows> ]
//	SHOW RUNS [ LIKE '<pattern>' ] IN EXPERIMENT <experiment_name>
//	SHOW RUN { METRICS | PARAMETERS } [ LIKE '<pattern>' ] IN EXPERIMENT <experiment_name>
//	SHOW MODEL MONITORS IN GATEWAY <gateway_name>
//
// https://docs.snowflake.com/en/sql-reference/sql/show-versions-in-model
// https://docs.snowflake.com/en/sql-reference/sql/show-functions-in-model
// https://docs.snowflake.com/en/sql-reference/sql/show-versions-in-dataset
// https://docs.snowflake.com/en/sql-reference/sql/show-runs-in-experiment
// https://docs.snowflake.com/en/sql-reference/sql/show-run-in-experiment
func (p *Parser) startsShowModelScope() string {
	for _, w := range []string{"MODEL", "DATASET", "EXPERIMENT", "GATEWAY"} {
		if !p.atWords(w) {
			continue
		}
		switch next := p.peekNext().Type; next {
		case kwSTARTS, kwLIMIT, kwWITH, kwLIKE, kwFOR, tokFlow, tokEOF, ';':
			return ""
		default:
			if p.isObjectTypeWord(next) {
				return w
			}
		}
		return ""
	}
	return ""
}

// parseShowModelScope reads <word> <name> as SHOW's scope. cur is word.
func (p *Parser) parseShowModelScope(stmt *ast.ShowStmt, word string) error {
	p.advance() // consume the scope word
	name, err := p.parseObjectName()
	if err != nil {
		return err
	}
	stmt.Scope = ast.ShowScopeOther
	stmt.ScopeText = word
	stmt.ScopeName = name
	return nil
}
