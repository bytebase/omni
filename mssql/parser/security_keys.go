// Package parser - security_keys.go implements parsing for security key/cert/credential statements.
//
// Every statement here is read by the grammar: ScriptDOM's TSql160.g (the rule
// is cited beside each function) where SQL Server 2022 agrees with it, and the
// server's verdict where it does not (mssql/parser/testdata/compat.tsv). Option
// words are matched by their unquoted text, since SQL Server refuses a
// bracketed option word (`WITH [SUBJECT] = 'a'`), and each option of a list
// may appear once.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// keysWordIs reports whether the current token is the unquoted word w (an
// identifier or a keyword), compared case-insensitively.
func (p *Parser) keysWordIs(w string) bool {
	return p.keysTokenIsWord(p.cur, w)
}

// keysTokenIsWord reports whether tok is the unquoted word w (upper-case): its
// keyword token when w is a keyword, else an undelimited identifier.
func (p *Parser) keysTokenIsWord(tok Token, w string) bool {
	if kw := lookupKeyword(w); kw != tokIDENT {
		return tok.Type == kw
	}
	return tok.Type == tokIDENT && !p.isDelimited(tok) && strings.ToUpper(tok.Str) == w
}

// keysMatchWord consumes the current token if it is the unquoted word w.
func (p *Parser) keysMatchWord(w string) bool {
	if p.keysWordIs(w) {
		p.advance()
		return true
	}
	return false
}

// keysExpectWords consumes the unquoted words ws in order.
func (p *Parser) keysExpectWords(ws ...string) error {
	for _, w := range ws {
		if !p.keysMatchWord(w) {
			return p.syntaxErrorAtCur()
		}
	}
	return nil
}

// keysExpectTok consumes a token of type t.
func (p *Parser) keysExpectTok(t int) error {
	if p.cur.Type != t {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// keysName reads a one-part name.
func (p *Parser) keysName() (string, error) {
	name, ok := p.parseIdentifier()
	if !ok {
		return "", p.syntaxErrorAtCur()
	}
	return name, nil
}

// keysString reads a string literal ('...' or N'...').
func (p *Parser) keysString() (string, error) {
	if p.cur.Type != tokSCONST && p.cur.Type != tokNSCONST {
		return "", p.syntaxErrorAtCur()
	}
	s := p.cur.Str
	p.advance()
	return s, nil
}

// keysBinary reads a binary literal (0x...).
func (p *Parser) keysBinary() (string, error) {
	if (p.cur.Type != tokICONST && p.cur.Type != tokFCONST) ||
		len(p.cur.Str) < 2 || p.cur.Str[0] != '0' || (p.cur.Str[1] != 'x' && p.cur.Str[1] != 'X') {
		return "", p.syntaxErrorAtCur()
	}
	s := p.cur.Str
	p.advance()
	return s, nil
}

// keysEqString reads `= 'string'`.
func (p *Parser) keysEqString() (string, error) {
	if err := p.keysExpectTok('='); err != nil {
		return "", err
	}
	return p.keysString()
}

// keysEqWord reads `= WORD` where WORD is one of the unquoted words in set and
// returns it upper-cased.
func (p *Parser) keysEqWord(set []string) (string, error) {
	if err := p.keysExpectTok('='); err != nil {
		return "", err
	}
	for _, w := range set {
		if p.keysMatchWord(w) {
			return w, nil
		}
	}
	return "", p.syntaxErrorAtCur()
}

// keysPasswordEq reads `PASSWORD = 'string'` (passwordCrypto, TSql160.g:10634).
func (p *Parser) keysPasswordEq() (string, error) {
	if err := p.keysExpectWords("PASSWORD"); err != nil {
		return "", err
	}
	return p.keysEqString()
}

// keysByPassword reads `<word> BY PASSWORD = 'string'`, as in
// ENCRYPTION BY PASSWORD = '...' (encryptClause, TSql160.g:1504).
func (p *Parser) keysByPassword(word string) (string, error) {
	if err := p.keysExpectWords(word, "BY"); err != nil {
		return "", err
	}
	return p.keysPasswordEq()
}

// keysOnce reports an error at the current token when the option has already
// been read; SQL Server refuses an option given twice.
func (p *Parser) keysOnce(seen map[string]bool, name string) error {
	if seen[name] {
		return p.syntaxErrorAtCur()
	}
	seen[name] = true
	return nil
}

// keysAlgorithms is what ALGORITHM takes where SQL Server 2022 generates a key
// (CREATE SYMMETRIC KEY ... WITH, CREATE ASYMMETRIC KEY ... WITH). The server
// refuses DES, TRIPLE_DES, TRIPLE_DES_3KEY, RC2, RC4, RC4_128, DESX, RSA_512
// and RSA_1024 there, though the docs list them
// (https://learn.microsoft.com/en-us/sql/t-sql/statements/create-symmetric-key-transact-sql).
var keysAlgorithms = []string{"AES_128", "AES_192", "AES_256", "RSA_2048", "RSA_3072", "RSA_4096"}

// keysProviderAlgorithms is what ALGORITHM takes for a key from an EKM provider
// (providerKeySourceOption, TSql160.g:10672): every algorithm the docs list.
var keysProviderAlgorithms = []string{
	"DES", "TRIPLE_DES", "TRIPLE_DES_3KEY", "RC2", "RC4", "RC4_128", "DESX",
	"AES_128", "AES_192", "AES_256", "RSA_512", "RSA_1024", "RSA_2048", "RSA_3072", "RSA_4096",
}

// keysDEKAlgorithms is what ALGORITHM takes for a database encryption key.
// SQL Server 2022 refuses TRIPLE_DES_3KEY, which the docs list
// (https://learn.microsoft.com/en-us/sql/t-sql/statements/create-database-encryption-key-transact-sql).
var keysDEKAlgorithms = []string{"AES_128", "AES_192", "AES_256"}

// keysOpts collects the clauses of a statement as String nodes.
type keysOpts struct{ items []nodes.Node }

func (o *keysOpts) add(s string) { o.items = append(o.items, &nodes.String{Str: s}) }

// kv adds an option NAME = value as two items, name then value.
func (o *keysOpts) kv(name, value string) { o.add(name); o.add(value) }

func (o *keysOpts) list() *nodes.List {
	if len(o.items) == 0 {
		return nil
	}
	return &nodes.List{Items: o.items}
}

// parseSecurityKeyStmt parses CREATE/ALTER/DROP of MASTER KEY, SYMMETRIC KEY,
// ASYMMETRIC KEY, CERTIFICATE, CREDENTIAL and CRYPTOGRAPHIC PROVIDER. The
// caller has consumed CREATE, ALTER or DROP.
func (p *Parser) parseSecurityKeyStmt(action string) (*nodes.SecurityKeyStmt, error) {
	stmt := &nodes.SecurityKeyStmt{
		Action: action,
		Loc:    nodes.Loc{Start: p.pos(), End: -1},
	}
	var o keysOpts
	var err error

	switch {
	case p.keysWordIs("MASTER"):
		p.advance()
		if err := p.keysExpectWords("KEY"); err != nil {
			return nil, err
		}
		stmt.ObjectType = "MASTER KEY"
		err = p.keysMasterKeyBody(action, &o)
	case p.keysWordIs("SYMMETRIC"):
		p.advance()
		if err := p.keysExpectWords("KEY"); err != nil {
			return nil, err
		}
		stmt.ObjectType = "SYMMETRIC KEY"
		if stmt.Name, err = p.keysName(); err != nil {
			return nil, err
		}
		err = p.keysSymmetricKeyBody(action, &o)
	case p.keysWordIs("ASYMMETRIC"):
		p.advance()
		if err := p.keysExpectWords("KEY"); err != nil {
			return nil, err
		}
		stmt.ObjectType = "ASYMMETRIC KEY"
		if stmt.Name, err = p.keysName(); err != nil {
			return nil, err
		}
		err = p.keysAsymmetricKeyBody(action, &o)
	case p.keysWordIs("CERTIFICATE"):
		p.advance()
		stmt.ObjectType = "CERTIFICATE"
		if stmt.Name, err = p.keysName(); err != nil {
			return nil, err
		}
		err = p.keysCertificateBody(action, &o)
	case p.keysWordIs("CREDENTIAL"):
		p.advance()
		stmt.ObjectType = "CREDENTIAL"
		if stmt.Name, err = p.keysName(); err != nil {
			return nil, err
		}
		if action != "DROP" {
			err = p.keysCredentialBody(action, &o)
		}
	case p.keysWordIs("CRYPTOGRAPHIC"):
		p.advance()
		if err := p.keysExpectWords("PROVIDER"); err != nil {
			return nil, err
		}
		stmt.ObjectType = "CRYPTOGRAPHIC PROVIDER"
		if stmt.Name, err = p.keysName(); err != nil {
			return nil, err
		}
		err = p.keysCryptographicProviderBody(action, &o)
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if err != nil {
		return nil, err
	}

	stmt.Options = o.list()
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// keysMasterKeyBody reads what follows MASTER KEY.
//
//	CREATE MASTER KEY ENCRYPTION BY PASSWORD = 'password'
//	    (createMasterKeyStatement, TSql160.g:7546; SQL Server 2022 refuses it
//	    without the password: "Database master keys without password are not
//	    supported in this version of SQL Server")
//	ALTER MASTER KEY { [ FORCE ] REGENERATE WITH ENCRYPTION BY PASSWORD = 'password'
//	    | { ADD | DROP } ENCRYPTION BY { SERVICE MASTER KEY | PASSWORD = 'password' } }
//	    (alterMasterKeyStatement, TSql160.g:12079)
//	DROP MASTER KEY (dropMasterKeyStatement, TSql160.g:15346)
func (p *Parser) keysMasterKeyBody(action string, o *keysOpts) error {
	switch action {
	case "CREATE":
		pwd, err := p.keysByPassword("ENCRYPTION")
		if err != nil {
			return err
		}
		o.add("ENCRYPTION BY PASSWORD = " + pwd)
	case "ALTER":
		switch {
		case p.keysWordIs("FORCE") || p.keysWordIs("REGENERATE"):
			if p.keysMatchWord("FORCE") {
				o.add("FORCE")
			}
			if err := p.keysExpectWords("REGENERATE", "WITH"); err != nil {
				return err
			}
			pwd, err := p.keysByPassword("ENCRYPTION")
			if err != nil {
				return err
			}
			o.add("REGENERATE")
			o.add("ENCRYPTION BY PASSWORD = " + pwd)
		case p.cur.Type == kwADD || p.cur.Type == kwDROP:
			op := strings.ToUpper(p.cur.Str)
			p.advance()
			o.add(op)
			if err := p.keysExpectWords("ENCRYPTION", "BY"); err != nil {
				return err
			}
			if p.keysMatchWord("SERVICE") {
				if err := p.keysExpectWords("MASTER", "KEY"); err != nil {
					return err
				}
				o.add("ENCRYPTION BY SERVICE MASTER KEY")
				return nil
			}
			pwd, err := p.keysPasswordEq()
			if err != nil {
				return err
			}
			o.add("ENCRYPTION BY PASSWORD = " + pwd)
		default:
			return p.syntaxErrorAtCur()
		}
	}
	return nil
}

// keysSymmetricKeyBody reads what follows SYMMETRIC KEY name.
//
//	CREATE SYMMETRIC KEY name [ AUTHORIZATION owner ]
//	    { WITH key_spec [ , ...n ] ENCRYPTION BY crypto [ , ...n ]
//	    | FROM PROVIDER provider [ WITH provider_option [ , ...n ] ] [ ENCRYPTION BY crypto [ , ...n ] ] }
//	    (createSymmetricKeyStatement, TSql160.g:10438)
//	ALTER SYMMETRIC KEY name { ADD | DROP } ENCRYPTION BY crypto [ , ...n ]
//	    (alterSymmetricKeyStatement, TSql160.g:10684)
//	DROP SYMMETRIC KEY name [ REMOVE PROVIDER KEY ] (dropSymmetricKeyStatement, TSql160.g:15368)
func (p *Parser) keysSymmetricKeyBody(action string, o *keysOpts) error {
	switch action {
	case "CREATE":
		if err := p.keysAuthorization(o); err != nil {
			return err
		}
		if p.keysMatchWord("FROM") {
			if err := p.keysExpectWords("PROVIDER"); err != nil {
				return err
			}
			prov, err := p.keysName()
			if err != nil {
				return err
			}
			o.add("FROM PROVIDER " + prov)
			if p.cur.Type == kwWITH {
				p.advance()
				// SQL Server takes IDENTITY_VALUE here, which the grammar's
				// providerKeySourceOption (TSql160.g:10672) does not.
				if err := p.keysProviderOptions(o, true); err != nil {
					return err
				}
			}
			if p.keysWordIs("ENCRYPTION") {
				return p.keysEncryptorList(o)
			}
			return nil
		}
		if err := p.keysExpectWords("WITH"); err != nil {
			return err
		}
		if err := p.keysSymmetricKeySpecs(o); err != nil {
			return err
		}
		return p.keysEncryptorList(o)
	case "ALTER":
		if p.cur.Type != kwADD && p.cur.Type != kwDROP {
			return p.syntaxErrorAtCur()
		}
		o.add(strings.ToUpper(p.cur.Str))
		p.advance()
		return p.keysEncryptorList(o)
	case "DROP":
		return p.keysRemoveProviderKey(o)
	}
	return nil
}

// keysSymmetricKeySpecs reads key_spec [ , ...n ] (keySpecList, TSql160.g:10470):
//
//	KEY_SOURCE = 'pass_phrase' | ALGORITHM = algorithm | IDENTITY_VALUE = 'identity_phrase'
func (p *Parser) keysSymmetricKeySpecs(o *keysOpts) error {
	seen := map[string]bool{}
	for {
		switch {
		case p.keysWordIs("KEY_SOURCE"), p.keysWordIs("IDENTITY_VALUE"):
			name := strings.ToUpper(p.cur.Str)
			if err := p.keysOnce(seen, name); err != nil {
				return err
			}
			p.advance()
			v, err := p.keysEqString()
			if err != nil {
				return err
			}
			o.kv(name, v)
		case p.keysWordIs("ALGORITHM"):
			if err := p.keysOnce(seen, "ALGORITHM"); err != nil {
				return err
			}
			p.advance()
			alg, err := p.keysEqWord(keysAlgorithms)
			if err != nil {
				return err
			}
			o.kv("ALGORITHM", alg)
		default:
			return p.syntaxErrorAtCur()
		}
		if _, ok := p.match(','); !ok {
			return nil
		}
	}
}

// keysProviderOptions reads provider_option [ , ...n ] after WITH
// (providerKeySourceOptions, TSql160.g:10654):
//
//	ALGORITHM = algorithm | PROVIDER_KEY_NAME = 'name'
//	| CREATION_DISPOSITION = { CREATE_NEW | OPEN_EXISTING }
//
// withIdentityValue also allows IDENTITY_VALUE = 'phrase', which SQL Server
// takes for a symmetric key and refuses for an asymmetric one.
func (p *Parser) keysProviderOptions(o *keysOpts, withIdentityValue bool) error {
	seen := map[string]bool{}
	for {
		switch {
		case p.keysWordIs("ALGORITHM"):
			if err := p.keysOnce(seen, "ALGORITHM"); err != nil {
				return err
			}
			p.advance()
			alg, err := p.keysEqWord(keysProviderAlgorithms)
			if err != nil {
				return err
			}
			o.kv("ALGORITHM", alg)
		case p.keysWordIs("PROVIDER_KEY_NAME"), withIdentityValue && p.keysWordIs("IDENTITY_VALUE"):
			name := strings.ToUpper(p.cur.Str)
			if err := p.keysOnce(seen, name); err != nil {
				return err
			}
			p.advance()
			v, err := p.keysEqString()
			if err != nil {
				return err
			}
			o.kv(name, v)
		case p.keysWordIs("CREATION_DISPOSITION"):
			// creationDispositionSourceOption, TSql160.g:10562
			if err := p.keysOnce(seen, "CREATION_DISPOSITION"); err != nil {
				return err
			}
			p.advance()
			v, err := p.keysEqWord([]string{"CREATE_NEW", "OPEN_EXISTING"})
			if err != nil {
				return err
			}
			o.kv("CREATION_DISPOSITION", v)
		default:
			return p.syntaxErrorAtCur()
		}
		if _, ok := p.match(','); !ok {
			return nil
		}
	}
}

// keysEncryptorList reads ENCRYPTION BY crypto [ , ...n ]
// (encryptorListWithNoPassword, TSql160.g:10583):
//
//	crypto ::= CERTIFICATE name | PASSWORD = 'password'
//	    | SYMMETRIC KEY name | ASYMMETRIC KEY name
func (p *Parser) keysEncryptorList(o *keysOpts) error {
	if err := p.keysExpectWords("ENCRYPTION", "BY"); err != nil {
		return err
	}
	for {
		mech, err := p.keysCrypto()
		if err != nil {
			return err
		}
		o.add("ENCRYPTION BY " + mech)
		if _, ok := p.match(','); !ok {
			return nil
		}
	}
}

// keysCrypto reads one crypto of cryptoWithNoPasswordOrJustPassword
// (TSql160.g:10599).
func (p *Parser) keysCrypto() (string, error) {
	switch {
	case p.keysMatchWord("CERTIFICATE"):
		name, err := p.keysName()
		if err != nil {
			return "", err
		}
		return "CERTIFICATE " + name, nil
	case p.keysWordIs("SYMMETRIC"), p.keysWordIs("ASYMMETRIC"):
		kind := strings.ToUpper(p.cur.Str)
		p.advance()
		if err := p.keysExpectWords("KEY"); err != nil {
			return "", err
		}
		name, err := p.keysName()
		if err != nil {
			return "", err
		}
		return kind + " KEY " + name, nil
	case p.keysWordIs("PASSWORD"):
		pwd, err := p.keysPasswordEq()
		if err != nil {
			return "", err
		}
		return "PASSWORD = " + pwd, nil
	}
	return "", p.syntaxErrorAtCur()
}

// keysAuthorization reads an optional AUTHORIZATION owner (authorizationOpt,
// TSql160.g:6883).
func (p *Parser) keysAuthorization(o *keysOpts) error {
	if p.cur.Type != kwAUTHORIZATION {
		return nil
	}
	p.advance()
	owner, err := p.keysName()
	if err != nil {
		return err
	}
	o.add("AUTHORIZATION " + owner)
	return nil
}

// keysRemoveProviderKey reads an optional REMOVE PROVIDER KEY
// (removeProviderKeyOpt, TSql160.g:15400).
func (p *Parser) keysRemoveProviderKey(o *keysOpts) error {
	if !p.keysMatchWord("REMOVE") {
		return nil
	}
	if err := p.keysExpectWords("PROVIDER", "KEY"); err != nil {
		return err
	}
	o.add("REMOVE PROVIDER KEY")
	return nil
}

// keysAsymmetricKeyBody reads what follows ASYMMETRIC KEY name.
//
//	CREATE ASYMMETRIC KEY name [ AUTHORIZATION owner ]
//	    { FROM { [ EXECUTABLE ] FILE = 'path' | ASSEMBLY name
//	           | PROVIDER provider WITH provider_option [ , ...n ] }
//	    | WITH ALGORITHM = algorithm }
//	    [ ENCRYPTION BY PASSWORD = 'password' ]
//	    (createAsymmetricKeyStatement, TSql160.g:1324; createAsymmetricKeyParams, TSql160.g:1347)
//	ALTER ASYMMETRIC KEY name { REMOVE { PRIVATE KEY | ATTESTED OPTION }
//	    | WITH PRIVATE KEY ( password_option [ , ...n ] ) | ATTESTED BY 'string' }
//	    (alterAsymmetricKeyStatement, TSql160.g:11321)
//	DROP ASYMMETRIC KEY name [ REMOVE PROVIDER KEY ] (dropAsymmetricKeyStatement, TSql160.g:15384)
func (p *Parser) keysAsymmetricKeyBody(action string, o *keysOpts) error {
	switch action {
	case "CREATE":
		if err := p.keysAuthorization(o); err != nil {
			return err
		}
		switch {
		case p.keysMatchWord("FROM"):
			switch {
			case p.keysWordIs("EXECUTABLE"), p.keysWordIs("FILE"):
				// fileEncryptionSource, TSql160.g:1389
				prefix := "FROM "
				if p.keysMatchWord("EXECUTABLE") {
					prefix += "EXECUTABLE "
				}
				if err := p.keysExpectWords("FILE"); err != nil {
					return err
				}
				path, err := p.keysEqString()
				if err != nil {
					return err
				}
				o.add(prefix + "FILE = " + path)
			case p.keysMatchWord("ASSEMBLY"):
				name, err := p.keysName()
				if err != nil {
					return err
				}
				o.add("FROM ASSEMBLY " + name)
			case p.keysMatchWord("PROVIDER"):
				// providerEncryptionSource, TSql160.g:1377: the WITH list is
				// required here.
				prov, err := p.keysName()
				if err != nil {
					return err
				}
				o.add("FROM PROVIDER " + prov)
				if err := p.keysExpectWords("WITH"); err != nil {
					return err
				}
				if err := p.keysProviderOptions(o, false); err != nil {
					return err
				}
			default:
				return p.syntaxErrorAtCur()
			}
		case p.keysMatchWord("WITH"):
			// asymKeySpec, TSql160.g:1405
			if err := p.keysExpectWords("ALGORITHM"); err != nil {
				return err
			}
			alg, err := p.keysEqWord(keysAlgorithms)
			if err != nil {
				return err
			}
			o.kv("ALGORITHM", alg)
		default:
			return p.syntaxErrorAtCur()
		}
		if p.keysWordIs("ENCRYPTION") {
			pwd, err := p.keysByPassword("ENCRYPTION")
			if err != nil {
				return err
			}
			o.add("ENCRYPTION BY PASSWORD = " + pwd)
		}
	case "ALTER":
		switch {
		case p.keysWordIs("REMOVE"), p.keysWordIs("ATTESTED"):
			return p.keysRemoveOrAttested(o)
		case p.keysMatchWord("WITH"):
			if err := p.keysExpectWords("PRIVATE", "KEY"); err != nil {
				return err
			}
			if err := p.keysExpectTok('('); err != nil {
				return err
			}
			// passwordChangeOption, TSql160.g:1534
			seen := map[string]bool{}
			for {
				if !p.keysWordIs("ENCRYPTION") && !p.keysWordIs("DECRYPTION") {
					return p.syntaxErrorAtCur()
				}
				word := strings.ToUpper(p.cur.Str)
				if err := p.keysOnce(seen, word); err != nil {
					return err
				}
				pwd, err := p.keysByPassword(word)
				if err != nil {
					return err
				}
				o.add(word + " BY PASSWORD = " + pwd)
				if _, ok := p.match(','); !ok {
					break
				}
			}
			return p.keysExpectTok(')')
		default:
			return p.syntaxErrorAtCur()
		}
	case "DROP":
		return p.keysRemoveProviderKey(o)
	}
	return nil
}

// keysRemoveOrAttested reads REMOVE { PRIVATE KEY | ATTESTED OPTION } or
// ATTESTED BY 'string', shared by ALTER CERTIFICATE (TSql160.g:11415) and
// ALTER ASYMMETRIC KEY (TSql160.g:11321).
func (p *Parser) keysRemoveOrAttested(o *keysOpts) error {
	if p.keysMatchWord("ATTESTED") {
		if err := p.keysExpectWords("BY"); err != nil {
			return err
		}
		s, err := p.keysString()
		if err != nil {
			return err
		}
		o.add("ATTESTED BY " + s)
		return nil
	}
	if err := p.keysExpectWords("REMOVE"); err != nil {
		return err
	}
	if p.keysMatchWord("PRIVATE") {
		if err := p.keysExpectWords("KEY"); err != nil {
			return err
		}
		o.add("REMOVE PRIVATE KEY")
		return nil
	}
	if err := p.keysExpectWords("ATTESTED", "OPTION"); err != nil {
		return err
	}
	o.add("REMOVE ATTESTED OPTION")
	return nil
}

// keysCertificateBody reads what follows CERTIFICATE name.
//
//	CREATE CERTIFICATE name [ AUTHORIZATION owner ]
//	    { FROM source | [ ENCRYPTION BY PASSWORD = 'password' ] WITH option [ , ...n ] }
//	    [ ACTIVE FOR BEGIN_DIALOG = { ON | OFF } ]
//	    (createCertificateStatement, TSql160.g:1414; createCertificateParams, TSql160.g:1433)
//	source ::= ASSEMBLY name
//	    | { [ EXECUTABLE ] FILE = 'path' | BINARY = bits } [ WITH key_item [ , ...n ] ]
//	key_item ::= FORMAT = 'format' | PRIVATE KEY ( private_key_option [ , ...n ] )
//	option ::= SUBJECT = 'subject' | START_DATE = 'date' | EXPIRY_DATE = 'date'
//	ALTER CERTIFICATE name { REMOVE { PRIVATE KEY | ATTESTED OPTION }
//	    | WITH { PRIVATE KEY ( private_key_option [ , ...n ] ) | ACTIVE FOR BEGIN_DIALOG = { ON | OFF } }
//	    | ATTESTED BY 'string' }
//	    (alterCertificateStatement, TSql160.g:11415)
//	DROP CERTIFICATE name (dropCertificateStatement, TSql160.g:15414)
//
// BINARY and FORMAT are SQL Server's, not the grammar's
// (https://learn.microsoft.com/en-us/sql/t-sql/statements/create-certificate-transact-sql).
func (p *Parser) keysCertificateBody(action string, o *keysOpts) error {
	switch action {
	case "CREATE":
		if err := p.keysAuthorization(o); err != nil {
			return err
		}
		if p.keysMatchWord("FROM") {
			if p.keysMatchWord("ASSEMBLY") {
				// assemblyEncryptionSource, TSql160.g:1366
				name, err := p.keysName()
				if err != nil {
					return err
				}
				o.add("FROM ASSEMBLY " + name)
			} else {
				if err := p.keysCertificateFileSource(o); err != nil {
					return err
				}
				if p.cur.Type == kwWITH {
					p.advance()
					if err := p.keysCertificateKeyItems(o, keysPrivateKeyCreate); err != nil {
						return err
					}
				}
			}
		} else {
			if p.keysWordIs("ENCRYPTION") {
				pwd, err := p.keysByPassword("ENCRYPTION")
				if err != nil {
					return err
				}
				o.add("ENCRYPTION BY PASSWORD = " + pwd)
			}
			if err := p.keysExpectWords("WITH"); err != nil {
				return err
			}
			// certificateOption, TSql160.g:1472
			seen := map[string]bool{}
			for {
				if !p.keysWordIs("SUBJECT") && !p.keysWordIs("START_DATE") && !p.keysWordIs("EXPIRY_DATE") {
					return p.syntaxErrorAtCur()
				}
				name := strings.ToUpper(p.cur.Str)
				if err := p.keysOnce(seen, name); err != nil {
					return err
				}
				p.advance()
				v, err := p.keysEqString()
				if err != nil {
					return err
				}
				o.kv(name, v)
				if _, ok := p.match(','); !ok {
					break
				}
			}
		}
		if p.keysWordIs("ACTIVE") {
			return p.keysActiveForBeginDialog(o)
		}
	case "ALTER":
		switch {
		case p.keysWordIs("REMOVE"), p.keysWordIs("ATTESTED"):
			return p.keysRemoveOrAttested(o)
		case p.keysMatchWord("WITH"):
			if p.keysWordIs("ACTIVE") {
				return p.keysActiveForBeginDialog(o)
			}
			if err := p.keysExpectWords("PRIVATE", "KEY"); err != nil {
				return err
			}
			return p.keysPrivateKeySpec(o, keysPrivateKeyCreate)
		default:
			return p.syntaxErrorAtCur()
		}
	}
	return nil
}

// keysCertificateFileSource reads [ EXECUTABLE ] FILE = 'path' | BINARY = bits.
func (p *Parser) keysCertificateFileSource(o *keysOpts) error {
	if p.keysMatchWord("BINARY") {
		bits, err := p.keysBinaryEq()
		if err != nil {
			return err
		}
		o.add("FROM BINARY = " + bits)
		return nil
	}
	prefix := "FROM "
	if p.keysMatchWord("EXECUTABLE") {
		prefix += "EXECUTABLE "
	}
	if err := p.keysExpectWords("FILE"); err != nil {
		return err
	}
	path, err := p.keysEqString()
	if err != nil {
		return err
	}
	o.add(prefix + "FILE = " + path)
	return nil
}

// keysBinaryEq reads `= 0x...`.
func (p *Parser) keysBinaryEq() (string, error) {
	if err := p.keysExpectTok('='); err != nil {
		return "", err
	}
	return p.keysBinary()
}

// keysActiveForBeginDialog reads ACTIVE FOR BEGIN_DIALOG = { ON | OFF }
// (createCertificateActivityFlag, TSql160.g:1460).
func (p *Parser) keysActiveForBeginDialog(o *keysOpts) error {
	if err := p.keysExpectWords("ACTIVE", "FOR", "BEGIN_DIALOG"); err != nil {
		return err
	}
	v, err := p.keysEqWord([]string{"ON", "OFF"})
	if err != nil {
		return err
	}
	o.add("ACTIVE FOR BEGIN_DIALOG = " + v)
	return nil
}

// keysPrivateKeyMode says which private key options a statement takes.
type keysPrivateKeyMode int

const (
	// keysPrivateKeyCreate is CREATE and ALTER CERTIFICATE: FILE or BINARY,
	// not both, and the passwords; the list may hold only passwords.
	keysPrivateKeyCreate keysPrivateKeyMode = iota
	// keysPrivateKeyBackup is BACKUP CERTIFICATE: FILE and BINARY may both
	// appear and one must, and ALGORITHM = 'algorithm' is allowed.
	keysPrivateKeyBackup
)

// keysCertificateKeyItems reads key_item [ , ...n ] after WITH, where
// key_item ::= FORMAT = 'format' | PRIVATE KEY ( ... ), each once.
func (p *Parser) keysCertificateKeyItems(o *keysOpts, mode keysPrivateKeyMode) error {
	seen := map[string]bool{}
	for {
		switch {
		case p.keysWordIs("FORMAT"):
			if err := p.keysOnce(seen, "FORMAT"); err != nil {
				return err
			}
			p.advance()
			v, err := p.keysEqString()
			if err != nil {
				return err
			}
			o.kv("FORMAT", v)
		case p.keysWordIs("PRIVATE"):
			if err := p.keysOnce(seen, "PRIVATE"); err != nil {
				return err
			}
			p.advance()
			if err := p.keysExpectWords("KEY"); err != nil {
				return err
			}
			if err := p.keysPrivateKeySpec(o, mode); err != nil {
				return err
			}
		default:
			return p.syntaxErrorAtCur()
		}
		if _, ok := p.match(','); !ok {
			return nil
		}
	}
}

// keysPrivateKeySpec reads ( private_key_option [ , ...n ] ) after PRIVATE KEY
// (privateKeySpec, TSql160.g:1512; certificatePrivateKeySpec, TSql160.g:1520):
//
//	private_key_option ::= FILE = 'path' | BINARY = bits
//	    | DECRYPTION BY PASSWORD = 'password' | ENCRYPTION BY PASSWORD = 'password'
//
// Each option appears once.
func (p *Parser) keysPrivateKeySpec(o *keysOpts, mode keysPrivateKeyMode) error {
	start := p.cur
	if err := p.keysExpectTok('('); err != nil {
		return err
	}
	seen := map[string]bool{}
	for {
		switch {
		case p.keysWordIs("FILE"), p.keysWordIs("BINARY"):
			name := strings.ToUpper(p.cur.Str)
			if err := p.keysOnce(seen, name); err != nil {
				return err
			}
			if mode == keysPrivateKeyCreate && seen["FILE"] && seen["BINARY"] {
				return p.syntaxErrorAtCur()
			}
			p.advance()
			var v string
			var err error
			if name == "FILE" {
				v, err = p.keysEqString()
			} else {
				v, err = p.keysBinaryEq()
			}
			if err != nil {
				return err
			}
			o.add("PRIVATE KEY " + name + " = " + v)
		case p.keysWordIs("ENCRYPTION"), p.keysWordIs("DECRYPTION"):
			word := strings.ToUpper(p.cur.Str)
			if err := p.keysOnce(seen, word); err != nil {
				return err
			}
			pwd, err := p.keysByPassword(word)
			if err != nil {
				return err
			}
			o.add("PRIVATE KEY " + word + " BY PASSWORD = " + pwd)
		case mode == keysPrivateKeyBackup && p.keysWordIs("ALGORITHM"):
			if err := p.keysOnce(seen, "ALGORITHM"); err != nil {
				return err
			}
			p.advance()
			v, err := p.keysEqString()
			if err != nil {
				return err
			}
			o.add("PRIVATE KEY ALGORITHM = " + v)
		default:
			return p.syntaxErrorAtCur()
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if mode == keysPrivateKeyBackup && !seen["FILE"] && !seen["BINARY"] {
		// SQL Server: "Incorrect syntax near 'PRIVATE KEY'".
		return &ParseError{Message: "PRIVATE KEY needs FILE or BINARY", Position: start.Loc, Source: p.source}
	}
	return p.keysExpectTok(')')
}

// keysCredentialBody reads what follows CREDENTIAL name for CREATE and ALTER:
//
//	WITH IDENTITY = 'identity' [ , SECRET = 'secret' ]
//	    [ FOR CRYPTOGRAPHIC PROVIDER provider ]   -- CREATE only
//
// (credentialStatementBody, TSql160.g:1641; createCredentialStatement,
// TSql160.g:1625; alterCredentialStatement, TSql160.g:11478).
func (p *Parser) keysCredentialBody(action string, o *keysOpts) error {
	if err := p.keysExpectWords("WITH", "IDENTITY"); err != nil {
		return err
	}
	id, err := p.keysEqString()
	if err != nil {
		return err
	}
	o.kv("IDENTITY", id)
	if _, ok := p.match(','); ok {
		if err := p.keysExpectWords("SECRET"); err != nil {
			return err
		}
		secret, err := p.keysEqString()
		if err != nil {
			return err
		}
		o.kv("SECRET", secret)
	}
	if action == "CREATE" && p.cur.Type == kwFOR {
		p.advance()
		if err := p.keysExpectWords("CRYPTOGRAPHIC", "PROVIDER"); err != nil {
			return err
		}
		prov, err := p.keysName()
		if err != nil {
			return err
		}
		o.add("FOR CRYPTOGRAPHIC PROVIDER " + prov)
	}
	return nil
}

// keysCryptographicProviderBody reads what follows CRYPTOGRAPHIC PROVIDER name.
//
//	CREATE CRYPTOGRAPHIC PROVIDER name FROM FILE = 'path'
//	    (createCryptographicProviderStatement, TSql160.g:21980)
//	ALTER CRYPTOGRAPHIC PROVIDER name { FROM FILE = 'path' | ENABLE | DISABLE }
//	    (alterCryptographicProviderStatement, TSql160.g:22008; SQL Server
//	    refuses ENABLE or DISABLE after FROM FILE, which the docs show)
//	DROP CRYPTOGRAPHIC PROVIDER name
func (p *Parser) keysCryptographicProviderBody(action string, o *keysOpts) error {
	switch action {
	case "CREATE":
		return p.keysCryptographicProviderFile(o)
	case "ALTER":
		if p.keysWordIs("ENABLE") || p.keysWordIs("DISABLE") {
			o.add(strings.ToUpper(p.cur.Str))
			p.advance()
			return nil
		}
		return p.keysCryptographicProviderFile(o)
	}
	return nil
}

// keysCryptographicProviderFile reads FROM FILE = 'path' (cryptographicProviderFile, TSql160.g:21998).
func (p *Parser) keysCryptographicProviderFile(o *keysOpts) error {
	if err := p.keysExpectWords("FROM", "FILE"); err != nil {
		return err
	}
	path, err := p.keysEqString()
	if err != nil {
		return err
	}
	o.add("FROM FILE = " + path)
	return nil
}

// parseSecurityKeyStmtColumn parses CREATE/ALTER/DROP COLUMN { ENCRYPTION KEY | MASTER KEY }.
// Called after "COLUMN" has been matched by the dispatcher.
//
//	CREATE COLUMN MASTER KEY name WITH ( KEY_STORE_PROVIDER_NAME = 'name', KEY_PATH = 'path'
//	    [ , ENCLAVE_COMPUTATIONS ( SIGNATURE = bits ) ] )
//	    (createColumnMasterKeyStatement, TSql160.g:23148; SQL Server takes the
//	    first two in either order and ENCLAVE_COMPUTATIONS only after them)
//	CREATE COLUMN ENCRYPTION KEY name WITH VALUES value [ , value ]
//	    (createColumnEncryptionKeyStatement, TSql160.g:23235)
//	ALTER COLUMN ENCRYPTION KEY name { ADD VALUE value | DROP VALUE ( COLUMN_MASTER_KEY = name ) }
//	    (alterColumnEncryptionKeyStatement, TSql160.g:23331)
//	DROP COLUMN { MASTER | ENCRYPTION } KEY name
func (p *Parser) parseSecurityKeyStmtColumn(action string) (*nodes.SecurityKeyStmt, error) {
	stmt := &nodes.SecurityKeyStmt{
		Action: action,
		Loc:    nodes.Loc{Start: p.pos(), End: -1},
	}
	var o keysOpts

	switch {
	case p.keysWordIs("ENCRYPTION"):
		stmt.ObjectType = "COLUMN ENCRYPTION KEY"
	case p.keysWordIs("MASTER") && action != "ALTER":
		stmt.ObjectType = "COLUMN MASTER KEY"
	default:
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if err := p.keysExpectWords("KEY"); err != nil {
		return nil, err
	}
	name, err := p.keysName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	switch {
	case action == "CREATE" && stmt.ObjectType == "COLUMN MASTER KEY":
		err = p.keysColumnMasterKeyParams(&o)
	case action == "CREATE":
		if err := p.keysExpectWords("WITH", "VALUES"); err != nil {
			return nil, err
		}
		if err = p.keysColumnEncryptionKeyValue(&o); err == nil {
			if _, ok := p.match(','); ok {
				err = p.keysColumnEncryptionKeyValue(&o)
			}
		}
	case action == "ALTER":
		switch p.cur.Type {
		case kwADD:
			p.advance()
			if err := p.keysExpectWords("VALUE"); err != nil {
				return nil, err
			}
			o.add("ADD")
			o.add("VALUE")
			err = p.keysColumnEncryptionKeyValue(&o)
		case kwDROP:
			// alterColumnEncryptionKeyDropStatement, TSql160.g:23364
			p.advance()
			if err := p.keysExpectWords("VALUE"); err != nil {
				return nil, err
			}
			if err := p.keysExpectTok('('); err != nil {
				return nil, err
			}
			if err := p.keysExpectWords("COLUMN_MASTER_KEY"); err != nil {
				return nil, err
			}
			if err := p.keysExpectTok('='); err != nil {
				return nil, err
			}
			cmk, err := p.keysName()
			if err != nil {
				return nil, err
			}
			o.add("DROP")
			o.add("VALUE")
			o.kv("COLUMN_MASTER_KEY", cmk)
			if err := p.keysExpectTok(')'); err != nil {
				return nil, err
			}
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}
	if err != nil {
		return nil, err
	}

	stmt.Options = o.list()
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// keysColumnMasterKeyParams reads the WITH ( ... ) list of CREATE COLUMN
// MASTER KEY (columnMasterkeyParameter, TSql160.g:23178).
func (p *Parser) keysColumnMasterKeyParams(o *keysOpts) error {
	if err := p.keysExpectWords("WITH"); err != nil {
		return err
	}
	if err := p.keysExpectTok('('); err != nil {
		return err
	}
	seen := map[string]bool{}
	for i := 0; i < 2; i++ {
		if i > 0 {
			if err := p.keysExpectTok(','); err != nil {
				return err
			}
		}
		if !p.keysWordIs("KEY_STORE_PROVIDER_NAME") && !p.keysWordIs("KEY_PATH") {
			return p.syntaxErrorAtCur()
		}
		name := strings.ToUpper(p.cur.Str)
		if err := p.keysOnce(seen, name); err != nil {
			return err
		}
		p.advance()
		v, err := p.keysEqString()
		if err != nil {
			return err
		}
		o.kv(name, v)
	}
	if _, ok := p.match(','); ok {
		if err := p.keysExpectWords("ENCLAVE_COMPUTATIONS"); err != nil {
			return err
		}
		if err := p.keysExpectTok('('); err != nil {
			return err
		}
		if err := p.keysExpectWords("SIGNATURE"); err != nil {
			return err
		}
		sig, err := p.keysBinaryEq()
		if err != nil {
			return err
		}
		if err := p.keysExpectTok(')'); err != nil {
			return err
		}
		o.add("ENCLAVE_COMPUTATIONS (SIGNATURE = " + sig + ")")
	}
	return p.keysExpectTok(')')
}

// keysColumnEncryptionKeyValue reads one value of a column encryption key
// (createColumnEncryptionKeyValue, TSql160.g:23259):
//
//	( COLUMN_MASTER_KEY = name, ALGORITHM = 'algorithm', ENCRYPTED_VALUE = bits )
//
// in any order, each once.
func (p *Parser) keysColumnEncryptionKeyValue(o *keysOpts) error {
	if err := p.keysExpectTok('('); err != nil {
		return err
	}
	seen := map[string]bool{}
	for i := 0; i < 3; i++ {
		if i > 0 {
			if err := p.keysExpectTok(','); err != nil {
				return err
			}
		}
		var v string
		var err error
		name := strings.ToUpper(p.cur.Str)
		switch {
		case p.keysWordIs("COLUMN_MASTER_KEY"):
			if err := p.keysOnce(seen, name); err != nil {
				return err
			}
			p.advance()
			if err := p.keysExpectTok('='); err != nil {
				return err
			}
			v, err = p.keysName()
		case p.keysWordIs("ALGORITHM"):
			if err := p.keysOnce(seen, name); err != nil {
				return err
			}
			p.advance()
			v, err = p.keysEqString()
		case p.keysWordIs("ENCRYPTED_VALUE"):
			if err := p.keysOnce(seen, name); err != nil {
				return err
			}
			p.advance()
			v, err = p.keysBinaryEq()
		default:
			return p.syntaxErrorAtCur()
		}
		if err != nil {
			return err
		}
		o.kv(name, v)
	}
	return p.keysExpectTok(')')
}

// parseSecurityKeyStmtDatabaseEncryption parses CREATE/ALTER/DROP DATABASE ENCRYPTION KEY
// and CREATE/ALTER/DROP DATABASE SCOPED CREDENTIAL.
// Called after "DATABASE" has been consumed and the next token is "ENCRYPTION" or "SCOPED".
//
//	CREATE DATABASE ENCRYPTION KEY WITH ALGORITHM = algorithm
//	    ENCRYPTION BY SERVER { CERTIFICATE name | ASYMMETRIC KEY name }
//	    (createDatabaseEncryptionKey, TSql160.g:5506)
//	ALTER DATABASE ENCRYPTION KEY { REGENERATE WITH ALGORITHM = algorithm
//	    | ENCRYPTION BY SERVER { CERTIFICATE name | ASYMMETRIC KEY name } }
//	    (alterDatabaseEncryptionKey, TSql160.g:4235; SQL Server refuses an
//	    ENCRYPTION BY after REGENERATE, which the grammar allows)
//	DROP DATABASE ENCRYPTION KEY
//	{ CREATE | ALTER } DATABASE SCOPED CREDENTIAL name WITH IDENTITY = 'identity'
//	    [ , SECRET = 'secret' ] [ FOR CRYPTOGRAPHIC PROVIDER provider ]
//	    (createDatabaseScopedCredentialStatement, TSql160.g:1614; SQL Server
//	    takes the FOR clause on CREATE)
//	DROP DATABASE SCOPED CREDENTIAL name
func (p *Parser) parseSecurityKeyStmtDatabaseEncryption(action string) (*nodes.SecurityKeyStmt, error) {
	stmt := &nodes.SecurityKeyStmt{
		Action: action,
		Loc:    nodes.Loc{Start: p.pos(), End: -1},
	}
	var o keysOpts

	switch {
	case p.keysMatchWord("ENCRYPTION"):
		if err := p.keysExpectWords("KEY"); err != nil {
			return nil, err
		}
		stmt.ObjectType = "DATABASE ENCRYPTION KEY"
		var err error
		switch action {
		case "CREATE":
			if err := p.keysExpectWords("WITH", "ALGORITHM"); err != nil {
				return nil, err
			}
			alg, err := p.keysEqWord(keysDEKAlgorithms)
			if err != nil {
				return nil, err
			}
			o.kv("ALGORITHM", alg)
			err = p.keysDEKEncryptor(&o)
			if err != nil {
				return nil, err
			}
		case "ALTER":
			if p.keysMatchWord("REGENERATE") {
				if err := p.keysExpectWords("WITH", "ALGORITHM"); err != nil {
					return nil, err
				}
				alg, err := p.keysEqWord(keysDEKAlgorithms)
				if err != nil {
					return nil, err
				}
				o.add("REGENERATE WITH ALGORITHM = " + alg)
			} else {
				err = p.keysDEKEncryptor(&o)
			}
		}
		if err != nil {
			return nil, err
		}
	case p.keysMatchWord("SCOPED"):
		if err := p.keysExpectWords("CREDENTIAL"); err != nil {
			return nil, err
		}
		stmt.ObjectType = "DATABASE SCOPED CREDENTIAL"
		name, err := p.keysName()
		if err != nil {
			return nil, err
		}
		stmt.Name = name
		if action != "DROP" {
			if err := p.keysCredentialBody(action, &o); err != nil {
				return nil, err
			}
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Options = o.list()
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// keysDEKEncryptor reads ENCRYPTION BY SERVER { CERTIFICATE | ASYMMETRIC KEY } name
// (dekEncryptor, TSql160.g:5522; dekEncryptorType, TSql160.g:5532).
func (p *Parser) keysDEKEncryptor(o *keysOpts) error {
	if err := p.keysExpectWords("ENCRYPTION", "BY", "SERVER"); err != nil {
		return err
	}
	kind := "CERTIFICATE"
	if p.keysMatchWord("ASYMMETRIC") {
		if err := p.keysExpectWords("KEY"); err != nil {
			return err
		}
		kind = "ASYMMETRIC KEY"
	} else if err := p.keysExpectWords("CERTIFICATE"); err != nil {
		return err
	}
	name, err := p.keysName()
	if err != nil {
		return err
	}
	o.add("ENCRYPTION BY SERVER " + kind + " " + name)
	return nil
}

// parseOpenSymmetricKeyStmt parses OPEN SYMMETRIC KEY statements.
//
//	OPEN SYMMETRIC KEY name DECRYPTION BY
//	    { CERTIFICATE name [ WITH PASSWORD = 'password' ]
//	    | ASYMMETRIC KEY name [ WITH PASSWORD = 'password' ]
//	    | SYMMETRIC KEY name
//	    | PASSWORD = 'password' }
//
// (openSymmetricKeyStatement, TSql160.g:14463; decryptionMechanism, TSql160.g:14430)
func (p *Parser) parseOpenSymmetricKeyStmt() (*nodes.SecurityKeyStmt, error) {
	loc := p.pos()
	p.advance() // consume OPEN
	stmt := &nodes.SecurityKeyStmt{
		Action:     "OPEN",
		ObjectType: "SYMMETRIC KEY",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	var o keysOpts
	if err := p.keysExpectWords("SYMMETRIC", "KEY"); err != nil {
		return nil, err
	}
	name, err := p.keysName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if err := p.keysExpectWords("DECRYPTION", "BY"); err != nil {
		return nil, err
	}
	mech, err := p.keysCrypto()
	if err != nil {
		return nil, err
	}
	o.add("DECRYPTION BY " + mech)
	// The action at TSql160.g:14444 refuses WITH after SYMMETRIC KEY.
	if p.cur.Type == kwWITH && (strings.HasPrefix(mech, "CERTIFICATE ") || strings.HasPrefix(mech, "ASYMMETRIC ")) {
		p.advance()
		pwd, err := p.keysPasswordEq()
		if err != nil {
			return nil, err
		}
		o.add("WITH PASSWORD = " + pwd)
	}

	stmt.Options = o.list()
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCloseSymmetricKeyStmt parses CLOSE SYMMETRIC KEY or CLOSE ALL SYMMETRIC KEYS
// (closeSymmetricKeyStatement, TSql160.g:14527).
func (p *Parser) parseCloseSymmetricKeyStmt() (*nodes.SecurityKeyStmt, error) {
	loc := p.pos()
	p.advance() // consume CLOSE
	stmt := &nodes.SecurityKeyStmt{
		Action: "CLOSE",
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	if p.cur.Type == kwALL {
		p.advance() // consume ALL
		if err := p.keysExpectWords("SYMMETRIC", "KEYS"); err != nil {
			return nil, err
		}
		stmt.ObjectType = "ALL SYMMETRIC KEYS"
	} else {
		if err := p.keysExpectWords("SYMMETRIC", "KEY"); err != nil {
			return nil, err
		}
		stmt.ObjectType = "SYMMETRIC KEY"
		name, err := p.keysName()
		if err != nil {
			return nil, err
		}
		stmt.Name = name
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseBackupCertificateStmt parses BACKUP CERTIFICATE|MASTER KEY|SYMMETRIC KEY statements.
//
//	BACKUP CERTIFICATE name TO FILE = 'path' [ WITH key_item [ , ...n ] ]
//	    (backupCertificateStatement, TSql160.g:6005)
//	key_item ::= FORMAT = 'format' | PRIVATE KEY ( private_key_option [ , ...n ] )
//	private_key_option ::= FILE = 'path' | BINARY = bits | ALGORITHM = 'algorithm'
//	    | ENCRYPTION BY PASSWORD = 'password' | DECRYPTION BY PASSWORD = 'password'
//	BACKUP MASTER KEY TO { FILE | URL } = 'path' ENCRYPTION BY PASSWORD = 'password'
//	    (backupMasterKeyStatement, TSql160.g:6061)
//	BACKUP SYMMETRIC KEY name TO { FILE | URL } = 'path' ENCRYPTION BY PASSWORD = 'password'
//	    (https://learn.microsoft.com/en-us/sql/t-sql/statements/backup-symmetric-key-transact-sql)
//
// FORMAT, BINARY, ALGORITHM and URL are SQL Server's, not the grammar's
// (https://learn.microsoft.com/en-us/sql/t-sql/statements/backup-certificate-transact-sql).
func (p *Parser) parseBackupCertificateStmt() (*nodes.SecurityKeyStmt, error) {
	loc := p.pos()
	p.advance() // consume BACKUP
	stmt := &nodes.SecurityKeyStmt{
		Action: "BACKUP",
		Loc:    nodes.Loc{Start: loc, End: -1},
	}
	var o keysOpts

	switch {
	case p.keysMatchWord("CERTIFICATE"):
		stmt.ObjectType = "CERTIFICATE"
		name, err := p.keysName()
		if err != nil {
			return nil, err
		}
		stmt.Name = name
		if err := p.keysExpectWords("TO", "FILE"); err != nil {
			return nil, err
		}
		path, err := p.keysEqString()
		if err != nil {
			return nil, err
		}
		o.add("TO FILE = " + path)
		if p.cur.Type == kwWITH {
			p.advance()
			if err := p.keysCertificateKeyItems(&o, keysPrivateKeyBackup); err != nil {
				return nil, err
			}
		}
	case p.keysMatchWord("SYMMETRIC"):
		if err := p.keysExpectWords("KEY"); err != nil {
			return nil, err
		}
		stmt.ObjectType = "SYMMETRIC KEY"
		name, err := p.keysName()
		if err != nil {
			return nil, err
		}
		stmt.Name = name
		if err := p.keysKeyFileTransfer("TO", true, &o); err != nil {
			return nil, err
		}
	case p.keysMatchWord("MASTER"):
		if err := p.keysExpectWords("KEY"); err != nil {
			return nil, err
		}
		stmt.ObjectType = "MASTER KEY"
		if err := p.keysKeyFileTransfer("TO", true, &o); err != nil {
			return nil, err
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Options = o.list()
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// keysKeyFileTransfer reads `{ TO | FROM } { FILE | URL } = 'path'` (URL only
// when allowURL) and the password clause after it: ENCRYPTION BY PASSWORD for
// TO, DECRYPTION BY PASSWORD for FROM.
func (p *Parser) keysKeyFileTransfer(dir string, allowURL bool, o *keysOpts) error {
	if err := p.keysExpectWords(dir); err != nil {
		return err
	}
	kind := "FILE"
	if allowURL && p.keysMatchWord("URL") {
		kind = "URL"
	} else if err := p.keysExpectWords("FILE"); err != nil {
		return err
	}
	path, err := p.keysEqString()
	if err != nil {
		return err
	}
	o.add(dir + " " + kind + " = " + path)
	word := "ENCRYPTION"
	if dir == "FROM" {
		word = "DECRYPTION"
	}
	pwd, err := p.keysByPassword(word)
	if err != nil {
		return err
	}
	o.add(word + " BY PASSWORD = " + pwd)
	return nil
}

// parseOpenMasterKeyStmt parses OPEN MASTER KEY DECRYPTION BY PASSWORD = 'password'
// (openMasterKeyStatement, TSql160.g:14407).
func (p *Parser) parseOpenMasterKeyStmt() (*nodes.SecurityKeyStmt, error) {
	loc := p.pos()
	p.advance() // consume OPEN
	stmt := &nodes.SecurityKeyStmt{
		Action:     "OPEN",
		ObjectType: "MASTER KEY",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	var o keysOpts
	if err := p.keysExpectWords("MASTER", "KEY"); err != nil {
		return nil, err
	}
	pwd, err := p.keysByPassword("DECRYPTION")
	if err != nil {
		return nil, err
	}
	o.add("DECRYPTION BY PASSWORD = " + pwd)
	stmt.Options = o.list()
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCloseMasterKeyStmt parses CLOSE MASTER KEY (closeMasterKeyStatement, TSql160.g:14515).
func (p *Parser) parseCloseMasterKeyStmt() (*nodes.SecurityKeyStmt, error) {
	loc := p.pos()
	p.advance() // consume CLOSE
	stmt := &nodes.SecurityKeyStmt{
		Action:     "CLOSE",
		ObjectType: "MASTER KEY",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	if err := p.keysExpectWords("MASTER", "KEY"); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseRestoreMasterKeyStmt parses
//
//	RESTORE MASTER KEY FROM { FILE | URL } = 'path'
//	    DECRYPTION BY PASSWORD = 'password'
//	    ENCRYPTION BY PASSWORD = 'password' [ FORCE ]
//
// (restoreMasterKeyStatement, TSql160.g:6079; URL is SQL Server's:
// https://learn.microsoft.com/en-us/sql/t-sql/statements/restore-master-key-transact-sql).
func (p *Parser) parseRestoreMasterKeyStmt() (*nodes.SecurityKeyStmt, error) {
	loc := p.pos()
	p.advance() // consume RESTORE
	stmt := &nodes.SecurityKeyStmt{
		Action:     "RESTORE",
		ObjectType: "MASTER KEY",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	var o keysOpts
	if err := p.keysExpectWords("MASTER", "KEY"); err != nil {
		return nil, err
	}
	if err := p.keysKeyFileTransfer("FROM", true, &o); err != nil {
		return nil, err
	}
	pwd, err := p.keysByPassword("ENCRYPTION")
	if err != nil {
		return nil, err
	}
	o.add("ENCRYPTION BY PASSWORD = " + pwd)
	if p.keysMatchWord("FORCE") {
		o.add("FORCE")
	}
	stmt.Options = o.list()
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseRestoreSymmetricKeyStmt parses
//
//	RESTORE SYMMETRIC KEY name FROM { FILE | URL } = 'path'
//	    DECRYPTION BY PASSWORD = 'password'
//	    ENCRYPTION BY PASSWORD = 'password'
//
// (https://learn.microsoft.com/en-us/sql/t-sql/statements/restore-symmetric-key-transact-sql).
func (p *Parser) parseRestoreSymmetricKeyStmt() (*nodes.SecurityKeyStmt, error) {
	loc := p.pos()
	p.advance() // consume RESTORE
	stmt := &nodes.SecurityKeyStmt{
		Action:     "RESTORE",
		ObjectType: "SYMMETRIC KEY",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	var o keysOpts
	if err := p.keysExpectWords("SYMMETRIC", "KEY"); err != nil {
		return nil, err
	}
	name, err := p.keysName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if err := p.keysKeyFileTransfer("FROM", true, &o); err != nil {
		return nil, err
	}
	pwd, err := p.keysByPassword("ENCRYPTION")
	if err != nil {
		return nil, err
	}
	o.add("ENCRYPTION BY PASSWORD = " + pwd)
	stmt.Options = o.list()
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// keysServiceMasterKeyBody reads what follows { ALTER | BACKUP | RESTORE }
// SERVICE MASTER KEY.
//
//	ALTER SERVICE MASTER KEY { [ FORCE ] REGENERATE
//	    | WITH OLD_ACCOUNT = 'account', OLD_PASSWORD = 'password'
//	    | WITH NEW_ACCOUNT = 'account', NEW_PASSWORD = 'password' }
//	    (alterServiceMasterKeyStatement, TSql160.g:12732)
//	BACKUP SERVICE MASTER KEY TO FILE = 'path' ENCRYPTION BY PASSWORD = 'password'
//	    (backupServiceMasterKeyStatement, TSql160.g:6022)
//	RESTORE SERVICE MASTER KEY FROM FILE = 'path' DECRYPTION BY PASSWORD = 'password' [ FORCE ]
//	    (restoreServiceMasterKeyStatement, TSql160.g:6030)
func (p *Parser) keysServiceMasterKeyBody(action string, stmt *nodes.SecurityKeyStmt) error {
	var o keysOpts
	switch action {
	case "ALTER":
		switch {
		case p.keysWordIs("FORCE"), p.keysWordIs("REGENERATE"):
			if p.keysMatchWord("FORCE") {
				o.add("FORCE")
			}
			if err := p.keysExpectWords("REGENERATE"); err != nil {
				return err
			}
			o.add("REGENERATE")
		case p.keysMatchWord("WITH"):
			prefix := ""
			switch {
			case p.keysWordIs("OLD_ACCOUNT"):
				prefix = "OLD_"
			case p.keysWordIs("NEW_ACCOUNT"):
				prefix = "NEW_"
			default:
				return p.syntaxErrorAtCur()
			}
			p.advance()
			account, err := p.keysEqString()
			if err != nil {
				return err
			}
			if err := p.keysExpectTok(','); err != nil {
				return err
			}
			if err := p.keysExpectWords(prefix + "PASSWORD"); err != nil {
				return err
			}
			pwd, err := p.keysEqString()
			if err != nil {
				return err
			}
			o.kv(prefix+"ACCOUNT", account)
			o.kv(prefix+"PASSWORD", pwd)
		default:
			return p.syntaxErrorAtCur()
		}
	case "BACKUP":
		if err := p.keysKeyFileTransfer("TO", false, &o); err != nil {
			return err
		}
	case "RESTORE":
		if err := p.keysKeyFileTransfer("FROM", false, &o); err != nil {
			return err
		}
		if p.keysMatchWord("FORCE") {
			o.add("FORCE")
		}
	}
	stmt.Options = o.list()
	return nil
}

// keysAlterDatabaseScopedCredential parses the rest of ALTER DATABASE SCOPED
// CREDENTIAL name WITH IDENTITY = 'identity' [ , SECRET = 'secret' ] after
// CREDENTIAL (createDatabaseScopedCredentialStatement, TSql160.g:1614).
func (p *Parser) keysAlterDatabaseScopedCredential(loc int) (*nodes.SecurityKeyStmt, error) {
	stmt := &nodes.SecurityKeyStmt{
		Action:     "ALTER",
		ObjectType: "DATABASE SCOPED CREDENTIAL",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	name, err := p.keysName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	var o keysOpts
	if err := p.keysCredentialBody("ALTER", &o); err != nil {
		return nil, err
	}
	stmt.Options = o.list()
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseSignatureStmt parses ADD/DROP [COUNTER] SIGNATURE after SIGNATURE.
//
//	ADD [ COUNTER ] SIGNATURE TO signable BY crypto [ , ...n ]
//	DROP [ COUNTER ] SIGNATURE FROM signable BY crypto [ , ...n ]
//	signable ::= [ { OBJECT | ASSEMBLY | DATABASE } :: ] [ [ db . ] schema . ] name
//	crypto ::= { CERTIFICATE name | ASYMMETRIC KEY name }
//	        [ WITH { PASSWORD = 'password' | SIGNATURE = bits } ]
//	    | SYMMETRIC KEY name | PASSWORD = 'password'
//
// (addSignatureStatement, TSql160.g:15747; dropSignatureStatement,
// TSql160.g:15751; signableElement, TSql160.g:15768;
// cryptoWithOptionalPasswordSignature, TSql160.g:15809; SQL Server takes
// SYMMETRIC KEY without the WITH clause.)
func (p *Parser) parseSignatureStmt(action string) (*nodes.SignatureStmt, error) {
	stmt := &nodes.SignatureStmt{
		Action: action,
		Loc:    nodes.Loc{Start: p.pos(), End: -1},
	}
	dir := "TO"
	if action == "DROP" {
		dir = "FROM"
	}
	if err := p.keysExpectWords(dir); err != nil {
		return nil, err
	}

	if (p.keysWordIs("OBJECT") || p.keysWordIs("ASSEMBLY") || p.keysWordIs("DATABASE")) &&
		p.peekNext().Type == tokCOLONCOLON {
		stmt.ModuleClass = strings.ToUpper(p.cur.Str)
		p.advance() // class
		p.advance() // ::
	}
	ref, err := p.keysThreePartName()
	if err != nil {
		return nil, err
	}
	stmt.ModuleName = ref

	if err := p.keysExpectWords("BY"); err != nil {
		return nil, err
	}
	var cryptos []nodes.Node
	for {
		item, err := p.keysSignatureCrypto()
		if err != nil {
			return nil, err
		}
		cryptos = append(cryptos, item)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	stmt.CryptoList = &nodes.List{Items: cryptos}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// keysThreePartName reads [ [ db . ] schema . ] name (schemaObjectThreePartName).
func (p *Parser) keysThreePartName() (*nodes.TableRef, error) {
	loc := p.pos()
	parts := []string{}
	for {
		name, err := p.keysName()
		if err != nil {
			return nil, err
		}
		parts = append(parts, name)
		if len(parts) == 3 || p.cur.Type != '.' {
			break
		}
		p.advance()
	}
	ref := &nodes.TableRef{Object: parts[len(parts)-1], Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}
	if len(parts) > 1 {
		ref.Schema = parts[len(parts)-2]
	}
	if len(parts) > 2 {
		ref.Database = parts[0]
	}
	return ref, nil
}

// keysSignatureCrypto reads one crypto of a signature's BY list.
func (p *Parser) keysSignatureCrypto() (*nodes.CryptoItem, error) {
	item := &nodes.CryptoItem{Loc: nodes.Loc{Start: p.pos(), End: -1}}
	withAllowed := true
	switch {
	case p.keysMatchWord("CERTIFICATE"):
		item.Mechanism = "CERTIFICATE"
	case p.keysWordIs("ASYMMETRIC"), p.keysWordIs("SYMMETRIC"):
		item.Mechanism = strings.ToUpper(p.cur.Str) + " KEY"
		withAllowed = item.Mechanism == "ASYMMETRIC KEY"
		p.advance()
		if err := p.keysExpectWords("KEY"); err != nil {
			return nil, err
		}
	case p.keysWordIs("PASSWORD"):
		pwd, err := p.keysPasswordEq()
		if err != nil {
			return nil, err
		}
		item.Mechanism = "PASSWORD"
		item.WithType = "PASSWORD"
		item.WithValue = pwd
		item.Loc.End = p.prevEnd()
		return item, nil
	default:
		return nil, p.syntaxErrorAtCur()
	}
	name, err := p.keysName()
	if err != nil {
		return nil, err
	}
	item.Name = name
	if withAllowed && p.cur.Type == kwWITH {
		p.advance()
		switch {
		case p.keysWordIs("PASSWORD"):
			pwd, err := p.keysPasswordEq()
			if err != nil {
				return nil, err
			}
			item.WithType, item.WithValue = "PASSWORD", pwd
		case p.keysMatchWord("SIGNATURE"):
			sig, err := p.keysBinaryEq()
			if err != nil {
				return nil, err
			}
			item.WithType, item.WithValue = "SIGNATURE", sig
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}
	item.Loc.End = p.prevEnd()
	return item, nil
}

// isStmtStart returns true if the token type could start a new statement.
func isStmtStart(t int) bool {
	switch t {
	case kwSELECT, kwINSERT, kwUPDATE, kwDELETE, kwCREATE, kwALTER, kwDROP,
		kwGRANT, kwREVOKE, kwDENY, kwEXEC, kwEXECUTE, kwDECLARE, kwSET,
		kwIF, kwWHILE, kwBEGIN, kwRETURN, kwGOTO, kwWAITFOR,
		kwPRINT, kwRAISERROR, kwTHROW, kwUSE, kwGO,
		kwCOMMIT, kwROLLBACK, kwSAVE, kwTRUNCATE, kwMERGE,
		kwOPEN, kwFETCH, kwCLOSE, kwDEALLOCATE, kwBACKUP:
		return true
	}
	return false
}
