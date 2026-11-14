package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// PROVISIONED THROUGHPUT (Cortex AI), read as an ast.ObjectStmt
// ---------------------------------------------------------------------------
//
// The command reference documents CREATE alone, with no IF NOT EXISTS, no
// tags, and every property required.

func init() {
	registerObject(&objectForm{
		words:  []string{"PROVISIONED", "THROUGHPUT"},
		create: parseCreateProvisionedThroughput,
		// https://docs.snowflake.com/en/sql-reference/sql/create-provisioned-throughput
		orReplace:     true,
		noIfNotExists: true,
		props: objectProps{
			"CLOUD_PROVIDER": propString,
			"MODEL":          propString,
			"PTUS":           propNumber,
			"TERM_START":     propString,
			"TERM_END":       propString,
		},
	})
}

// parseCreateProvisionedThroughput reads what follows the name of
//
//	CREATE [ OR REPLACE ] PROVISIONED THROUGHPUT <name>
//	  CLOUD_PROVIDER = '<cloud_provider>'
//	  MODEL = '<model_name>'
//	  PTUS = <num_ptus>
//	  TERM_START = '<start_date>'
//	  TERM_END = '<end_date>'
//
// https://docs.snowflake.com/en/sql-reference/sql/create-provisioned-throughput
func parseCreateProvisionedThroughput(p *Parser, s *ast.ObjectStmt) error {
	c, err := p.parseObjectProperties(s, "", true)
	if err != nil {
		return err
	}
	if len(c.Tags) > 0 {
		return &ParseError{
			Position: c.Tags[0].Name.Loc.Start, End: c.Tags[0].Name.Loc.End,
			Message: "syntax error: PROVISIONED THROUGHPUT does not take tags",
		}
	}
	return p.requireObjectProperties(c, "CLOUD_PROVIDER", "MODEL", "PTUS", "TERM_START", "TERM_END")
}
