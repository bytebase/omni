package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// COPY INTO <location>: an external location's parameters
// ---------------------------------------------------------------------------
//
// COPY INTO <location> writes an external location's access and encryption
// parameters right after the location, before FROM:
//
//	externalLocation (for Amazon S3) ::=
//	  '<protocol>://<bucket>[/<path>]'
//	  [ { STORAGE_INTEGRATION = <integration_name> }
//	  | { CREDENTIALS = ( { AWS_KEY_ID = '<string>' AWS_SECRET_KEY = '<string>' [ AWS_TOKEN = '<string>' ] } ) } ]
//	  [ ENCRYPTION = ( [ TYPE = 'AWS_CSE' ] [ MASTER_KEY = '<string>' ] |
//	                   [ TYPE = 'AWS_SSE_S3' ] |
//	                   [ TYPE = 'AWS_SSE_KMS' [ KMS_KEY_ID = '<string>' ] ] |
//	                   [ TYPE = 'NONE' ] ) ]
//
//	externalLocation (for Google Cloud Storage) ::=
//	  'gcs://<bucket>[/<path>]'
//	  [ STORAGE_INTEGRATION = <integration_name> ]
//	  [ ENCRYPTION = ( [ TYPE = 'GCS_SSE_KMS' ] [ KMS_KEY_ID = '<string>' ] | [ TYPE = 'NONE' ] ) ]
//
//	externalLocation (for Microsoft Azure) ::=
//	  'azure://<account>.blob.core.windows.net/<container>[/<path>]'
//	  [ { STORAGE_INTEGRATION = <integration_name> } | { CREDENTIALS = ( [ AZURE_SAS_TOKEN = '<string>' ] ) } ]
//	  [ ENCRYPTION = ( [ TYPE = { 'AZURE_CSE' | 'NONE' } ] [ MASTER_KEY = '<string>' ] ) ]
//
// https://docs.snowflake.com/en/sql-reference/sql/copy-into-location
//
// Each group is read as the page writes it: only the keys the location's
// cloud documents, each once, each a string, and a TYPE one the cloud lists.
// S3's CREDENTIALS needs both its keys.
// Every part of an ENCRYPTION (and of Azure's CREDENTIALS) is optional, so
// the page writes `( )` too. The same parameters written after FROM are
// copy options (copy.go), as the page's examples place them.

// locationCloud is the cloud an external location's protocol names.
type locationCloud int

const (
	cloudNone locationCloud = iota
	cloudS3
	cloudGCS
	cloudAzure
)

// externalLocationCloud reports the cloud of an external location's URI.
// '<protocol>' is s3, s3china or s3gov for Amazon S3 ("Amazon S3 compatible
// storage" is not a protocol the page lists).
func externalLocationCloud(uri string) locationCloud {
	scheme, _, ok := strings.Cut(uri, "://")
	if !ok {
		return cloudNone
	}
	switch strings.ToLower(scheme) {
	case "s3", "s3china", "s3gov":
		return cloudS3
	case "gcs":
		return cloudGCS
	case "azure":
		return cloudAzure
	}
	return cloudNone
}

// parseCopyLocationParams reads the parameters an external location of cloud
// takes before FROM: a STORAGE_INTEGRATION or a CREDENTIALS (GCS takes no
// CREDENTIALS), then an ENCRYPTION.
func (p *Parser) parseCopyLocationParams(cloud locationCloud) ([]*ast.CopyOption, error) {
	if cloud == cloudNone {
		return nil, nil
	}
	var opts []*ast.CopyOption
	switch {
	case p.atCopyLocationParam("STORAGE_INTEGRATION"):
		opt, err := p.parseObjectProperty()
		if err != nil {
			return nil, err
		}
		if !propName.matches(opt) {
			return nil, copyLocationParamError(opt, "STORAGE_INTEGRATION takes <integration_name>")
		}
		opts = append(opts, opt)
	case cloud != cloudGCS && p.atCopyLocationParam("CREDENTIALS"):
		opt, err := p.parseCopyLocationGroup()
		if err != nil {
			return nil, err
		}
		if err := checkCopyLocationGroup(opt, copyLocationCredentials(cloud)); err != nil {
			return nil, err
		}
		if cloud == cloudS3 && (findLocationOption(opt.Group, "AWS_KEY_ID") == nil || findLocationOption(opt.Group, "AWS_SECRET_KEY") == nil) {
			return nil, copyLocationParamError(opt, "CREDENTIALS takes AWS_KEY_ID and AWS_SECRET_KEY")
		}
		opts = append(opts, opt)
	}
	if p.atCopyLocationParam("ENCRYPTION") {
		opt, err := p.parseCopyLocationGroup()
		if err != nil {
			return nil, err
		}
		if err := checkCopyLocationGroup(opt, copyLocationEncryption(cloud)); err != nil {
			return nil, err
		}
		if err := checkCopyLocationEncryption(opt, cloud); err != nil {
			return nil, err
		}
		opts = append(opts, opt)
	}
	return opts, nil
}

// atCopyLocationParam reports whether cur begins `name =`, name unquoted.
func (p *Parser) atCopyLocationParam(name string) bool {
	return isWordToken(p.cur, name) && p.peekNext().Type == '='
}

// copyLocationCredentials are the keys a cloud's CREDENTIALS takes.
func copyLocationCredentials(cloud locationCloud) []string {
	if cloud == cloudS3 {
		return []string{"AWS_KEY_ID", "AWS_SECRET_KEY", "AWS_TOKEN"}
	}
	return []string{"AZURE_SAS_TOKEN"}
}

// copyLocationEncryption are the keys a cloud's ENCRYPTION takes.
func copyLocationEncryption(cloud locationCloud) []string {
	switch cloud {
	case cloudS3:
		return []string{"TYPE", "MASTER_KEY", "KMS_KEY_ID"}
	case cloudGCS:
		return []string{"TYPE", "KMS_KEY_ID"}
	}
	return []string{"TYPE", "MASTER_KEY"}
}

// parseCopyLocationGroup reads `<name> = ( [ <key> = '<string>' ... ] )`,
// cur being <name>: each key an unquoted word, each value a string, and the
// entries separated by blanks or, as COPY's option groups also take
// (copy.go), a comma.
func (p *Parser) parseCopyLocationGroup() (*ast.CopyOption, error) {
	key := p.advance()
	opt := &ast.CopyOption{Name: strings.ToUpper(key.Str), Loc: ast.Loc{Start: key.Loc.Start}}
	p.advance() // consume '='
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	opt.Group = []*ast.CopyOption{}
	for p.cur.Type != ')' {
		if len(opt.Group) > 0 && p.cur.Type == ',' && p.peekNext().Type != ')' {
			p.advance() // consume ','
		}
		if !(p.cur.Type == tokIdent || p.cur.Type >= 700) || p.peekNext().Type != '=' {
			return nil, p.syntaxErrorAtCur()
		}
		k := p.advance()
		p.advance() // consume '='
		v, err := p.expect(tokString)
		if err != nil {
			return nil, err
		}
		opt.Group = append(opt.Group, &ast.CopyOption{
			Name: strings.ToUpper(k.Str),
			Lit:  &ast.Literal{Kind: ast.LitString, Value: v.Str, Loc: v.Loc},
			Loc:  ast.Loc{Start: k.Loc.Start, End: v.Loc.End},
		})
	}
	p.advance() // consume ')'
	opt.Loc.End = p.prev.Loc.End
	return opt, nil
}

// checkCopyLocationGroup holds a group's keys to keys, each set once.
func checkCopyLocationGroup(opt *ast.CopyOption, keys []string) error {
	for i, e := range opt.Group {
		if !containsString(keys, e.Name) {
			return copyLocationParamError(e, opt.Name+" takes "+strings.Join(keys, ", "))
		}
		if findLocationOption(opt.Group[:i], e.Name) != nil {
			return &ParseError{Position: e.Loc.Start, End: e.Loc.End, Message: "syntax error: " + e.Name + " is set twice"}
		}
	}
	return nil
}

// checkCopyLocationEncryption holds an ENCRYPTION group to the TYPEs its
// cloud lists and to the key each TYPE takes: on S3, MASTER_KEY goes with
// AWS_CSE (the default) and KMS_KEY_ID with AWS_SSE_KMS; on GCS, KMS_KEY_ID
// goes with GCS_SSE_KMS (the default); Azure's two TYPEs both may carry a
// MASTER_KEY.
func checkCopyLocationEncryption(opt *ast.CopyOption, cloud locationCloud) error {
	typ := ""
	if t := findLocationOption(opt.Group, "TYPE"); t != nil {
		typ = strings.ToUpper(t.Lit.Value)
		var types []string
		switch cloud {
		case cloudS3:
			types = []string{"AWS_CSE", "AWS_SSE_S3", "AWS_SSE_KMS", "NONE"}
		case cloudGCS:
			types = []string{"GCS_SSE_KMS", "NONE"}
		default:
			types = []string{"AZURE_CSE", "NONE"}
		}
		if !containsString(types, typ) {
			return copyLocationParamError(t, "ENCRYPTION's TYPE is one of '"+strings.Join(types, "', '")+"'")
		}
	}
	for _, e := range opt.Group {
		ok := true
		switch {
		case cloud == cloudS3 && e.Name == "MASTER_KEY":
			ok = typ == "" || typ == "AWS_CSE"
		case cloud == cloudS3 && e.Name == "KMS_KEY_ID":
			ok = typ == "AWS_SSE_KMS"
		case cloud == cloudGCS && e.Name == "KMS_KEY_ID":
			ok = typ == "" || typ == "GCS_SSE_KMS"
		}
		if !ok {
			return copyLocationParamError(e, e.Name+" does not go with TYPE = '"+typ+"'")
		}
	}
	return nil
}

// findLocationOption returns the option of opts named name, or nil.
func findLocationOption(opts []*ast.CopyOption, name string) *ast.CopyOption {
	for _, o := range opts {
		if o.Name == name {
			return o
		}
	}
	return nil
}

func copyLocationParamError(o *ast.CopyOption, msg string) *ParseError {
	return &ParseError{Position: o.Loc.Start, End: o.Loc.End, Message: "syntax error: " + msg}
}
