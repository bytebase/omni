package parser

import (
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func (p *Parser) parseRedshiftDataMovementOptions() (*nodes.List, error) {
	var items []nodes.Node
	for {
		item, err := p.parseRedshiftDataMovementOption()
		if err != nil {
			return nil, err
		}
		if item == nil {
			break
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: items}, nil
}

func (p *Parser) parseRedshiftDataMovementOption() (*nodes.DefElem, error) {
	switch {
	case p.consumeRedshiftWord("iam_role"):
		value, err := p.parseRedshiftStringOrKeyword()
		if err != nil {
			return nil, err
		}
		return makeDefElem("iam_role", &nodes.String{Str: value}), nil
	case p.consumeRedshiftWord("credentials"):
		value, err := p.parseRedshiftStringOrKeyword()
		if err != nil {
			return nil, err
		}
		return makeDefElem("credentials", &nodes.String{Str: value}), nil
	case p.consumeRedshiftWord("access_key_id"):
		value, err := p.parseRedshiftCopyStringOption()
		if err != nil {
			return nil, err
		}
		return makeDefElem("access_key_id", &nodes.String{Str: value}), nil
	case p.consumeRedshiftWord("secret_access_key"):
		value, err := p.parseRedshiftCopyStringOption()
		if err != nil {
			return nil, err
		}
		return makeDefElem("secret_access_key", &nodes.String{Str: value}), nil
	case p.consumeRedshiftWord("session_token"):
		value, err := p.parseRedshiftCopyStringOption()
		if err != nil {
			return nil, err
		}
		return makeDefElem("session_token", &nodes.String{Str: value}), nil
	case p.cur.Type == FORMAT:
		p.advance()
		p.parseOptAs()
		value, err := p.parseRedshiftOptionWord()
		if err != nil {
			return nil, err
		}
		elem := makeDefElem("format", &nodes.String{Str: strings.ToLower(value)})
		if strings.EqualFold(value, "avro") && p.cur.Type == SCONST {
			elem.Defname = "avro"
			elem.Arg = &nodes.String{Str: p.cur.Str}
			p.advance()
		}
		return elem, nil
	case p.cur.Type == DELIMITER:
		p.advance()
		p.parseOptAs()
		value, err := p.parseRedshiftStringOrKeyword()
		if err != nil {
			return nil, err
		}
		return makeDefElem("delimiter", &nodes.String{Str: value}), nil
	case p.cur.Type == HEADER_P:
		p.advance()
		return makeDefElem("header", &nodes.Boolean{Boolval: true}), nil
	case p.cur.Type == QUOTE:
		p.advance()
		value, err := p.parseRedshiftCopyStringOption()
		if err != nil {
			return nil, err
		}
		return makeDefElem("quote", &nodes.String{Str: value}), nil
	case p.cur.Type == NULL_P:
		p.advance()
		value, err := p.parseRedshiftCopyStringOption()
		if err != nil {
			return nil, err
		}
		return makeDefElem("null", &nodes.String{Str: value}), nil
	case p.cur.Type == ESCAPE:
		p.advance()
		return makeDefElem("escape", &nodes.Boolean{Boolval: true}), nil
	case p.consumeRedshiftWord("manifest"):
		if p.consumeRedshiftWord("verbose") {
			return makeDefElem("manifest", &nodes.String{Str: "verbose"}), nil
		}
		return makeDefElem("manifest", &nodes.Boolean{Boolval: true}), nil
	case p.consumeRedshiftWord("encrypted"):
		if p.consumeRedshiftWord("auto") {
			return makeDefElem("encrypted", &nodes.String{Str: "auto"}), nil
		}
		return makeDefElem("encrypted", &nodes.Boolean{Boolval: true}), nil
	case p.consumeRedshiftWord("kms_key_id"):
		value, err := p.parseRedshiftCopyStringOption()
		if err != nil {
			return nil, err
		}
		return makeDefElem("kms_key_id", &nodes.String{Str: value}), nil
	case p.consumeRedshiftWord("region"):
		value, err := p.parseRedshiftCopyStringOption()
		if err != nil {
			return nil, err
		}
		return makeDefElem("region", &nodes.String{Str: value}), nil
	case p.consumeRedshiftWord("json"):
		value, err := p.parseRedshiftCopyStringOption()
		if err != nil {
			return nil, err
		}
		return makeDefElem("json", &nodes.String{Str: value}), nil
	case p.consumeRedshiftWord("avro"):
		value, err := p.parseRedshiftCopyStringOption()
		if err != nil {
			return nil, err
		}
		return makeDefElem("avro", &nodes.String{Str: value}), nil
	case p.consumeRedshiftWord("fixedwidth"):
		value, err := p.parseRedshiftCopyStringOption()
		if err != nil {
			return nil, err
		}
		return makeDefElem("fixedwidth", &nodes.String{Str: value}), nil
	case p.consumeRedshiftWord("dateformat"):
		value, err := p.parseRedshiftCopyStringOption()
		if err != nil {
			return nil, err
		}
		return makeDefElem("dateformat", &nodes.String{Str: value}), nil
	case p.consumeRedshiftWord("timeformat"):
		value, err := p.parseRedshiftCopyStringOption()
		if err != nil {
			return nil, err
		}
		return makeDefElem("timeformat", &nodes.String{Str: value}), nil
	case p.consumeRedshiftWord("encoding"):
		value, err := p.parseRedshiftCopyStringOption()
		if err != nil {
			return nil, err
		}
		return makeDefElem("encoding", &nodes.String{Str: strings.ToLower(value)}), nil
	case p.consumeRedshiftWord("acceptinvchars"):
		if p.cur.Type == SCONST || p.cur.Type == AS {
			value, err := p.parseRedshiftCopyStringOption()
			if err != nil {
				return nil, err
			}
			return makeDefElem("acceptinvchars", &nodes.String{Str: value}), nil
		}
		return makeDefElem("acceptinvchars", &nodes.Boolean{Boolval: true}), nil
	case p.consumeRedshiftWord("ignoreheader"):
		return makeDefElem("ignoreheader", &nodes.Integer{Ival: p.parseRedshiftCopyIntOption()}), nil
	case p.consumeRedshiftWord("readratio"):
		return makeDefElem("readratio", &nodes.Integer{Ival: p.parseRedshiftCopyIntOption()}), nil
	case p.consumeRedshiftWord("maxerror"):
		return makeDefElem("maxerror", &nodes.Integer{Ival: p.parseRedshiftCopyIntOption()}), nil
	case p.consumeRedshiftWord("comprows"):
		return makeDefElem("comprows", &nodes.Integer{Ival: p.parseRedshiftCopyIntOption()}), nil
	case p.consumeRedshiftWord("maxfilesize"):
		value, err := p.parseRedshiftUnloadSizeOption()
		if err != nil {
			return nil, err
		}
		return makeDefElem("maxfilesize", &nodes.String{Str: value}), nil
	case p.consumeRedshiftWord("rowgroupsize"):
		value, err := p.parseRedshiftUnloadSizeOption()
		if err != nil {
			return nil, err
		}
		return makeDefElem("rowgroupsize", &nodes.String{Str: value}), nil
	case p.consumeRedshiftWord("compupdate"):
		value, err := p.parseRedshiftOptionWord()
		if err != nil {
			return nil, err
		}
		return makeDefElem("compupdate", &nodes.String{Str: strings.ToLower(value)}), nil
	case p.consumeRedshiftWord("statupdate"):
		value, err := p.parseRedshiftOptionWord()
		if err != nil {
			return nil, err
		}
		return makeDefElem("statupdate", &nodes.String{Str: strings.ToLower(value)}), nil
	case p.consumeRedshiftWord("gzip"):
		return makeDefElem("compression", &nodes.String{Str: "gzip"}), nil
	case p.consumeRedshiftWord("bzip2"):
		return makeDefElem("compression", &nodes.String{Str: "bzip2"}), nil
	case p.consumeRedshiftWord("lzop"):
		return makeDefElem("compression", &nodes.String{Str: "lzop"}), nil
	case p.consumeRedshiftWord("zstd"):
		return makeDefElem("compression", &nodes.String{Str: "zstd"}), nil
	case p.consumeRedshiftWord("parquet"):
		return makeDefElem("format", &nodes.String{Str: "parquet"}), nil
	case p.consumeRedshiftWord("orc"):
		return makeDefElem("format", &nodes.String{Str: "orc"}), nil
	case p.consumeRedshiftFlagOption("removequotes"):
		return makeDefElem("removequotes", &nodes.Boolean{Boolval: true}), nil
	case p.consumeRedshiftFlagOption("blanksasnull"):
		return makeDefElem("blanksasnull", &nodes.Boolean{Boolval: true}), nil
	case p.consumeRedshiftFlagOption("emptyasnull"):
		return makeDefElem("emptyasnull", &nodes.Boolean{Boolval: true}), nil
	case p.consumeRedshiftFlagOption("explicit_ids"):
		return makeDefElem("explicit_ids", &nodes.Boolean{Boolval: true}), nil
	case p.consumeRedshiftFlagOption("truncatecolumns"):
		return makeDefElem("truncatecolumns", &nodes.Boolean{Boolval: true}), nil
	case p.consumeRedshiftFlagOption("acceptanydate"):
		return makeDefElem("acceptanydate", &nodes.Boolean{Boolval: true}), nil
	case p.consumeRedshiftFlagOption("roundec"):
		return makeDefElem("roundec", &nodes.Boolean{Boolval: true}), nil
	case p.consumeRedshiftFlagOption("trimblanks"):
		return makeDefElem("trimblanks", &nodes.Boolean{Boolval: true}), nil
	case p.consumeRedshiftFlagOption("ignoreblanklines"):
		return makeDefElem("ignoreblanklines", &nodes.Boolean{Boolval: true}), nil
	case p.consumeRedshiftFlagOption("fillrecord"):
		return makeDefElem("fillrecord", &nodes.Boolean{Boolval: true}), nil
	case p.consumeRedshiftFlagOption("noload"):
		return makeDefElem("noload", &nodes.Boolean{Boolval: true}), nil
	case p.consumeRedshiftFlagOption("ignoreallerrors"):
		return makeDefElem("ignoreallerrors", &nodes.Boolean{Boolval: true}), nil
	case p.consumeRedshiftFlagOption("addquotes"):
		return makeDefElem("addquotes", &nodes.Boolean{Boolval: true}), nil
	case p.consumeRedshiftFlagOption("allowoverwrite"):
		return makeDefElem("allowoverwrite", &nodes.Boolean{Boolval: true}), nil
	case p.consumeRedshiftFlagOption("cleanpath"):
		return makeDefElem("cleanpath", &nodes.Boolean{Boolval: true}), nil
	case p.consumeRedshiftWord("parallel"):
		value, err := p.parseRedshiftOptionWord()
		if err != nil {
			return nil, err
		}
		return makeDefElem("parallel", &nodes.String{Str: strings.ToLower(value)}), nil
	case p.consumeRedshiftWord("extension"):
		value, err := p.parseRedshiftCopyStringOption()
		if err != nil {
			return nil, err
		}
		return makeDefElem("extension", &nodes.String{Str: value}), nil
	case p.consumeRedshiftWord("partition"):
		if _, err := p.expect(BY); err != nil {
			return nil, err
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		columns := p.parseColumnList()
		if columns == nil {
			return nil, p.syntaxErrorAtCur()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		name := "partition_by"
		if p.consumeRedshiftWord("include") {
			name = "partition_by_include"
		}
		return makeDefElem(name, columns), nil
	// optional-probe: Redshift COPY/UNLOAD options are optional and terminate on the first non-option token.
	default:
		return nil, nil
	}
}

func (p *Parser) consumeRedshiftFlagOption(name string) bool {
	return p.consumeRedshiftWord(name)
}

func (p *Parser) parseRedshiftCopyStringOption() (string, error) {
	p.parseOptAs()
	return p.parseRedshiftStringOrKeyword()
}

func (p *Parser) parseRedshiftCopyIntOption() int64 {
	p.parseOptAs()
	return p.parseSignedIconst()
}

func (p *Parser) parseRedshiftUnloadSizeOption() (string, error) {
	p.parseOptAs()
	size := p.parseSignedIconst()
	unit := ""
	if p.redshiftWordEqual("mb") || p.redshiftWordEqual("gb") {
		unit = strings.ToLower(p.advance().Str)
	}
	if unit == "" {
		return strconv.FormatInt(size, 10), nil
	}
	return strconv.FormatInt(size, 10) + " " + unit, nil
}

func (p *Parser) parseRedshiftStringOrKeyword() (string, error) {
	if p.cur.Type == SCONST {
		value := p.cur.Str
		p.advance()
		return value, nil
	}
	return p.parseRedshiftOptionWord()
}

func (p *Parser) parseUnloadStmt(loc int) (nodes.Node, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	query := p.cur.Str
	if _, err := p.expect(SCONST); err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	if _, err := p.expect(TO); err != nil {
		return nil, err
	}
	target := p.cur.Str
	if _, err := p.expect(SCONST); err != nil {
		return nil, err
	}
	options, err := p.parseRedshiftDataMovementOptions()
	if err != nil {
		return nil, err
	}
	return &nodes.UnloadStmt{
		Query:   query,
		Target:  target,
		Options: options,
		Loc:     nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}
