package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateDatabaseStmt parses a CREATE DATABASE statement.
// The CREATE keyword has already been consumed, and DATABASE is the current token.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-DATABASE.html
//
//	CREATE DATABASE [ database_name ]
//	    [ USER SYS IDENTIFIED BY password ]
//	    [ USER SYSTEM IDENTIFIED BY password ]
//	    [ CONTROLFILE REUSE ]
//	    [ LOGFILE logfile_clause [, ...] ]
//	    [ MAXLOGFILES integer ]
//	    [ MAXLOGMEMBERS integer ]
//	    [ MAXLOGHISTORY integer ]
//	    [ MAXDATAFILES integer ]
//	    [ MAXINSTANCES integer ]
//	    [ { ARCHIVELOG | NOARCHIVELOG } ]
//	    [ FORCE LOGGING ]
//	    [ SET STANDBY NOLOGGING FOR { DATA AVAILABILITY | LOAD PERFORMANCE } ]
//	    [ CHARACTER SET charset ]
//	    [ NATIONAL CHARACTER SET charset ]
//	    [ SET DEFAULT { BIGFILE | SMALLFILE } TABLESPACE ]
//	    [ database_logging_clauses ]
//	    [ tablespace_clauses ]
//	    [ set_time_zone_clause ]
//	    [ [ BIGFILE | SMALLFILE ]
//	        DEFAULT TABLESPACE tablespace
//	        DATAFILE datafile_tempfile_spec
//	      | DEFAULT TEMPORARY TABLESPACE tablespace
//	        TEMPFILE datafile_tempfile_spec
//	      | UNDO TABLESPACE tablespace
//	        DATAFILE datafile_tempfile_spec ]
//	    [ enable_pluggable_database ]
func (p *Parser) parseCreateDatabaseStmt(start int) nodes.StmtNode {
	// DATABASE keyword already checked by caller but not consumed
	p.advance() // consume DATABASE

	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_DATABASE,
		Loc:        nodes.Loc{Start: start},
	}

	// Optional database name — peek at the next token to see if it's an identifier
	// (not a keyword that starts a clause like USER, LOGFILE, etc.)
	if (p.isIdentLike() || p.cur.Type == tokQIDENT) && !p.isCreateDatabaseClauseKeyword() {
		stmt.Name = p.parseObjectName()
	}

	// Parse CREATE DATABASE clauses
	opts := &nodes.List{}
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		optStart := p.pos()

		switch {
		// USER SYS IDENTIFIED BY password / USER SYSTEM IDENTIFIED BY password
		case p.cur.Type == kwUSER:
			p.advance() // consume USER
			userType := ""
			if p.isIdentLike() {
				userType = p.cur.Str // SYS or SYSTEM
				p.advance()
			}
			if p.cur.Type == kwIDENTIFIED {
				p.advance() // IDENTIFIED
				if p.isIdentLike() && p.cur.Str == "BY" {
					p.advance() // BY
				}
			}
			password := ""
			if p.isIdentLike() || p.cur.Type == tokSCONST {
				password = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "USER_" + userType, Value: password,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// CONTROLFILE REUSE
		case p.isIdentLike() && p.cur.Str == "CONTROLFILE":
			p.advance() // CONTROLFILE
			if p.isIdentLike() && p.cur.Str == "REUSE" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "CONTROLFILE_REUSE",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// LOGFILE GROUP n 'file' SIZE ... [, GROUP n 'file' SIZE ...]
		case p.isIdentLike() && p.cur.Str == "LOGFILE":
			p.advance() // LOGFILE
			logItems := &nodes.List{}
			for {
				lfStart := p.pos()
				groupNum := ""
				if p.cur.Type == kwGROUP {
					p.advance() // GROUP
					if p.cur.Type == tokICONST {
						groupNum = p.cur.Str
						p.advance()
					}
				}
				// Parse file specs
				var files []*nodes.DatafileClause
				for p.cur.Type == tokSCONST {
					df := p.parseDatafileClause()
					if df != nil {
						files = append(files, df)
					}
					if p.cur.Type == ',' && (p.isIdentLike() || p.cur.Type == tokSCONST) {
						// Could be next file or next GROUP — peek
						break
					}
				}
				fileList := &nodes.List{}
				for _, f := range files {
					fileList.Items = append(fileList.Items, f)
				}
				logItems.Items = append(logItems.Items, &nodes.DDLOption{
					Key: "LOGFILE_GROUP", Value: groupNum,
					Items: fileList,
					Loc:   nodes.Loc{Start: lfStart, End: p.pos()},
				})
				// Check for comma-separated groups
				if p.cur.Type == ',' {
					p.advance()
					continue
				}
				break
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "LOGFILE", Items: logItems,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// MAXLOGFILES integer
		case p.isIdentLike() && p.cur.Str == "MAXLOGFILES":
			p.advance()
			val := ""
			if p.cur.Type == tokICONST {
				val = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "MAXLOGFILES", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// MAXLOGMEMBERS integer
		case p.isIdentLike() && p.cur.Str == "MAXLOGMEMBERS":
			p.advance()
			val := ""
			if p.cur.Type == tokICONST {
				val = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "MAXLOGMEMBERS", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// MAXLOGHISTORY integer
		case p.isIdentLike() && p.cur.Str == "MAXLOGHISTORY":
			p.advance()
			val := ""
			if p.cur.Type == tokICONST {
				val = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "MAXLOGHISTORY", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// MAXDATAFILES integer
		case p.isIdentLike() && p.cur.Str == "MAXDATAFILES":
			p.advance()
			val := ""
			if p.cur.Type == tokICONST {
				val = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "MAXDATAFILES", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// MAXINSTANCES integer
		case p.isIdentLike() && p.cur.Str == "MAXINSTANCES":
			p.advance()
			val := ""
			if p.cur.Type == tokICONST {
				val = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "MAXINSTANCES", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// ARCHIVELOG
		case p.isIdentLike() && p.cur.Str == "ARCHIVELOG":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "ARCHIVELOG",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// NOARCHIVELOG
		case p.isIdentLike() && p.cur.Str == "NOARCHIVELOG":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "NOARCHIVELOG",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// FORCE LOGGING
		case p.cur.Type == kwFORCE:
			p.advance() // FORCE
			if p.cur.Type == kwLOGGING {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "FORCE_LOGGING",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// SET STANDBY NOLOGGING FOR ... / SET DEFAULT ... TABLESPACE / SET TIME_ZONE ...
		case p.cur.Type == kwSET:
			p.advance() // SET
			if p.isIdentLike() && p.cur.Str == "STANDBY" {
				// SET STANDBY NOLOGGING FOR { DATA AVAILABILITY | LOAD PERFORMANCE }
				p.advance() // STANDBY
				if p.cur.Type == kwNOLOGGING {
					p.advance() // NOLOGGING
				}
				if p.cur.Type == kwFOR {
					p.advance() // FOR
				}
				val := ""
				if p.isIdentLike() && p.cur.Str == "DATA" {
					p.advance()
					if p.isIdentLike() && p.cur.Str == "AVAILABILITY" {
						p.advance()
					}
					val = "DATA AVAILABILITY"
				} else if p.isIdentLike() && p.cur.Str == "LOAD" {
					p.advance()
					if p.isIdentLike() && p.cur.Str == "PERFORMANCE" {
						p.advance()
					}
					val = "LOAD PERFORMANCE"
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "SET_STANDBY_NOLOGGING", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.cur.Type == kwDEFAULT {
				// SET DEFAULT { BIGFILE | SMALLFILE } TABLESPACE
				p.advance() // DEFAULT
				val := ""
				if p.isIdentLike() && (p.cur.Str == "BIGFILE" || p.cur.Str == "SMALLFILE") {
					val = p.cur.Str
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "TABLESPACE" {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "SET_DEFAULT_TABLESPACE_TYPE", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "TIME_ZONE" {
				// SET TIME_ZONE = 'value'
				p.advance() // TIME_ZONE
				if p.cur.Type == '=' {
					p.advance()
				}
				val := ""
				if p.cur.Type == tokSCONST {
					val = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "SET_TIME_ZONE", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else {
				// Unknown SET clause, skip token
				p.advance()
			}

		// CHARACTER SET charset
		case p.isIdentLike() && p.cur.Str == "CHARACTER":
			p.advance() // CHARACTER
			if p.cur.Type == kwSET {
				p.advance() // SET
			}
			val := ""
			if p.isIdentLike() || p.cur.Type == tokQIDENT {
				val = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "CHARACTER_SET", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// NATIONAL CHARACTER SET charset
		case p.isIdentLike() && p.cur.Str == "NATIONAL":
			p.advance() // NATIONAL
			if p.isIdentLike() && p.cur.Str == "CHARACTER" {
				p.advance()
			}
			if p.cur.Type == kwSET {
				p.advance()
			}
			val := ""
			if p.isIdentLike() || p.cur.Type == tokQIDENT {
				val = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "NATIONAL_CHARACTER_SET", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// SYSAUX DATAFILE file_specification [, ...]
		case p.isIdentLike() && p.cur.Str == "SYSAUX":
			p.advance() // SYSAUX
			if p.isIdentLike() && (p.cur.Str == "DATAFILE" || p.cur.Str == "DATAFILES") {
				p.advance()
			}
			fileList := &nodes.List{}
			for p.cur.Type == tokSCONST {
				df := p.parseDatafileClause()
				if df != nil {
					fileList.Items = append(fileList.Items, df)
				}
				if p.cur.Type == ',' {
					p.advance()
					continue
				}
				break
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "SYSAUX_DATAFILE", Items: fileList,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// EXTENT MANAGEMENT LOCAL [ { AUTOALLOCATE | UNIFORM [ SIZE size_clause ] } ]
		case p.isIdentLike() && p.cur.Str == "EXTENT":
			p.advance() // EXTENT
			if p.isIdentLike() && p.cur.Str == "MANAGEMENT" {
				p.advance()
			}
			if p.cur.Type == kwLOCAL {
				p.advance()
			}
			val := ""
			if p.isIdentLike() && p.cur.Str == "AUTOALLOCATE" {
				val = "AUTOALLOCATE"
				p.advance()
			} else if p.isIdentLike() && p.cur.Str == "UNIFORM" {
				val = "UNIFORM"
				p.advance()
				if p.cur.Type == kwSIZE {
					p.advance()
					val += " " + p.parseSizeValue()
				}
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "EXTENT_MANAGEMENT_LOCAL", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// USING MIRROR COPY mirror_name
		case p.cur.Type == kwUSING:
			p.advance() // USING
			if p.isIdentLike() && p.cur.Str == "MIRROR" {
				p.advance() // MIRROR
				if p.isIdentLike() && p.cur.Str == "COPY" {
					p.advance()
				}
			}
			val := ""
			if p.isIdentLike() || p.cur.Type == tokQIDENT {
				val = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "USING_MIRROR_COPY", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// DATAFILE 'file' SIZE ... (standalone datafile spec, not part of tablespace)
		case p.isIdentLike() && p.cur.Str == "DATAFILE":
			p.advance() // DATAFILE
			fileList := &nodes.List{}
			for p.cur.Type == tokSCONST {
				df := p.parseDatafileClause()
				if df != nil {
					fileList.Items = append(fileList.Items, df)
				}
				if p.cur.Type == ',' {
					p.advance()
					continue
				}
				break
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "DATAFILE", Items: fileList,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// BIGFILE / SMALLFILE DEFAULT TABLESPACE
		case p.isIdentLike() && (p.cur.Str == "BIGFILE" || p.cur.Str == "SMALLFILE"):
			fileType := p.cur.Str
			p.advance()
			if p.cur.Type == kwDEFAULT {
				// BIGFILE/SMALLFILE DEFAULT TABLESPACE ...
				p.advance() // DEFAULT
				if p.isIdentLike() && p.cur.Str == "TABLESPACE" {
					p.advance() // TABLESPACE
				}
				tsName := ""
				if p.isIdentLike() || p.cur.Type == tokQIDENT {
					tsName = p.cur.Str
					p.advance()
				}
				// Parse DATAFILE spec
				fileList := &nodes.List{}
				if p.isIdentLike() && p.cur.Str == "DATAFILE" {
					p.advance()
					for p.cur.Type == tokSCONST {
						df := p.parseDatafileClause()
						if df != nil {
							fileList.Items = append(fileList.Items, df)
						}
						if p.cur.Type == ',' {
							p.advance()
							continue
						}
						break
					}
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DEFAULT_TABLESPACE", Value: fileType + " " + tsName,
					Items: fileList,
					Loc:   nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else {
				// Just BIGFILE/SMALLFILE without DEFAULT — skip
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "FILE_TYPE", Value: fileType,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// DEFAULT { TABLESPACE | [LOCAL] TEMPORARY TABLESPACE }
		case p.cur.Type == kwDEFAULT:
			p.advance() // DEFAULT
			if p.cur.Type == kwLOCAL {
				// DEFAULT LOCAL TEMPORARY TABLESPACE name [FOR {ALL|LEAF}] TEMPFILE ...
				p.advance() // LOCAL
				if p.cur.Type == kwTEMPORARY {
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "TABLESPACE" {
					p.advance()
				}
				tsName := ""
				if p.isIdentLike() || p.cur.Type == tokQIDENT {
					tsName = p.cur.Str
					p.advance()
				}
				// Optional FOR { ALL | LEAF }
				if p.cur.Type == kwFOR {
					p.advance()
					if p.cur.Type == kwALL || (p.isIdentLike() && p.cur.Str == "LEAF") {
						p.advance()
					}
				}
				fileList := &nodes.List{}
				if p.isIdentLike() && p.cur.Str == "TEMPFILE" {
					p.advance()
					for p.cur.Type == tokSCONST {
						df := p.parseDatafileClause()
						if df != nil {
							fileList.Items = append(fileList.Items, df)
						}
						if p.cur.Type == ',' {
							p.advance()
							continue
						}
						break
					}
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DEFAULT_LOCAL_TEMPORARY_TABLESPACE", Value: tsName,
					Items: fileList,
					Loc:   nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.cur.Type == kwTEMPORARY {
				// DEFAULT TEMPORARY TABLESPACE name TEMPFILE ...
				p.advance() // TEMPORARY
				if p.isIdentLike() && p.cur.Str == "TABLESPACE" {
					p.advance() // TABLESPACE
				}
				tsName := ""
				if p.isIdentLike() || p.cur.Type == tokQIDENT {
					tsName = p.cur.Str
					p.advance()
				}
				fileList := &nodes.List{}
				if p.isIdentLike() && p.cur.Str == "TEMPFILE" {
					p.advance()
					for p.cur.Type == tokSCONST {
						df := p.parseDatafileClause()
						if df != nil {
							fileList.Items = append(fileList.Items, df)
						}
						if p.cur.Type == ',' {
							p.advance()
							continue
						}
						break
					}
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DEFAULT_TEMPORARY_TABLESPACE", Value: tsName,
					Items: fileList,
					Loc:   nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "TABLESPACE" {
				// DEFAULT TABLESPACE name DATAFILE ...
				p.advance() // TABLESPACE
				tsName := ""
				if p.isIdentLike() || p.cur.Type == tokQIDENT {
					tsName = p.cur.Str
					p.advance()
				}
				fileList := &nodes.List{}
				if p.isIdentLike() && p.cur.Str == "DATAFILE" {
					p.advance()
					for p.cur.Type == tokSCONST {
						df := p.parseDatafileClause()
						if df != nil {
							fileList.Items = append(fileList.Items, df)
						}
						if p.cur.Type == ',' {
							p.advance()
							continue
						}
						break
					}
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DEFAULT_TABLESPACE", Value: tsName,
					Items: fileList,
					Loc:   nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else {
				// Unknown DEFAULT clause
				p.advance()
			}

		// UNDO TABLESPACE name DATAFILE ...
		case p.isIdentLike() && p.cur.Str == "UNDO":
			p.advance() // UNDO
			if p.isIdentLike() && p.cur.Str == "TABLESPACE" {
				p.advance() // TABLESPACE
			}
			tsName := ""
			if p.isIdentLike() || p.cur.Type == tokQIDENT {
				tsName = p.cur.Str
				p.advance()
			}
			fileList := &nodes.List{}
			if p.isIdentLike() && p.cur.Str == "DATAFILE" {
				p.advance()
				for p.cur.Type == tokSCONST {
					df := p.parseDatafileClause()
					if df != nil {
						fileList.Items = append(fileList.Items, df)
					}
					if p.cur.Type == ',' {
						p.advance()
						continue
					}
					break
				}
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "UNDO_TABLESPACE", Value: tsName,
				Items: fileList,
				Loc:   nodes.Loc{Start: optStart, End: p.pos()},
			})

		// ENABLE PLUGGABLE DATABASE [SEED [...]]
		case p.cur.Type == kwENABLE:
			p.advance() // ENABLE
			if p.isIdentLike() && p.cur.Str == "PLUGGABLE" {
				p.advance() // PLUGGABLE
				if p.cur.Type == kwDATABASE {
					p.advance() // DATABASE
				}
			}
			seedItems := &nodes.List{}
			// Optional SEED clause
			if p.isIdentLike() && p.cur.Str == "SEED" {
				p.advance()
				// Parse SEED sub-clauses: FILE_NAME_CONVERT, SYSTEM DATAFILES, SYSAUX DATAFILES, LOCAL UNDO
				for p.cur.Type != ';' && p.cur.Type != tokEOF {
					seedStart := p.pos()
					if p.isIdentLike() && p.cur.Str == "FILE_NAME_CONVERT" {
						p.advance()
						if p.cur.Type == '=' {
							p.advance()
						}
						if p.isIdentLike() && p.cur.Str == "NONE" {
							p.advance()
							seedItems.Items = append(seedItems.Items, &nodes.DDLOption{
								Key: "FILE_NAME_CONVERT", Value: "NONE",
								Loc: nodes.Loc{Start: seedStart, End: p.pos()},
							})
						} else if p.cur.Type == '(' {
							p.advance()
							var pairs []string
							for p.cur.Type != ')' && p.cur.Type != tokEOF {
								if p.cur.Type == tokSCONST {
									pairs = append(pairs, p.cur.Str)
								}
								p.advance()
							}
							if p.cur.Type == ')' {
								p.advance()
							}
							seedItems.Items = append(seedItems.Items, &nodes.DDLOption{
								Key: "FILE_NAME_CONVERT", Value: strings.Join(pairs, ","),
								Loc: nodes.Loc{Start: seedStart, End: p.pos()},
							})
						}
					} else if p.isIdentLike() && p.cur.Str == "SYSTEM" {
						p.advance() // SYSTEM
						if p.isIdentLike() && (p.cur.Str == "DATAFILES" || p.cur.Str == "DATAFILE") {
							p.advance()
						}
						// SIZE size_clause [AUTOEXTEND ...]
						size := ""
						if p.cur.Type == kwSIZE {
							p.advance()
							size = p.parseSizeValue()
						}
						ae := p.parseOptionalAutoextend()
						items := &nodes.List{}
						if ae != nil {
							items.Items = append(items.Items, ae)
						}
						seedItems.Items = append(seedItems.Items, &nodes.DDLOption{
							Key: "SYSTEM_DATAFILES", Value: size,
							Items: items,
							Loc:   nodes.Loc{Start: seedStart, End: p.pos()},
						})
					} else if p.isIdentLike() && p.cur.Str == "SYSAUX" {
						p.advance() // SYSAUX
						if p.isIdentLike() && (p.cur.Str == "DATAFILES" || p.cur.Str == "DATAFILE") {
							p.advance()
						}
						size := ""
						if p.cur.Type == kwSIZE {
							p.advance()
							size = p.parseSizeValue()
						}
						ae := p.parseOptionalAutoextend()
						items := &nodes.List{}
						if ae != nil {
							items.Items = append(items.Items, ae)
						}
						seedItems.Items = append(seedItems.Items, &nodes.DDLOption{
							Key: "SYSAUX_DATAFILES", Value: size,
							Items: items,
							Loc:   nodes.Loc{Start: seedStart, End: p.pos()},
						})
					} else if p.cur.Type == kwLOCAL {
						p.advance() // LOCAL
						if p.isIdentLike() && p.cur.Str == "UNDO" {
							p.advance()
						}
						val := ""
						if p.isIdentLike() && (p.cur.Str == "ON" || p.cur.Str == "OFF") {
							val = p.cur.Str
							p.advance()
						} else if p.cur.Type == kwON {
							val = "ON"
							p.advance()
						}
						seedItems.Items = append(seedItems.Items, &nodes.DDLOption{
							Key: "LOCAL_UNDO", Value: val,
							Loc: nodes.Loc{Start: seedStart, End: p.pos()},
						})
					} else {
						break
					}
				}
			}
			opt := &nodes.DDLOption{
				Key: "ENABLE_PLUGGABLE_DATABASE",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			}
			if len(seedItems.Items) > 0 {
				opt.Items = seedItems
			}
			opts.Items = append(opts.Items, opt)

		default:
			// Unknown token, advance to avoid infinite loop
			p.advance()
		}
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// isCreateDatabaseClauseKeyword returns true if the current token starts
// a CREATE DATABASE clause (and is not a database name).
func (p *Parser) isCreateDatabaseClauseKeyword() bool {
	switch p.cur.Type {
	case kwUSER, kwFORCE, kwSET, kwDEFAULT, kwENABLE:
		return true
	}
	if p.isIdentLike() {
		switch p.cur.Str {
		case "CONTROLFILE", "LOGFILE", "DATAFILE", "TEMPFILE",
			"MAXLOGFILES", "MAXLOGMEMBERS", "MAXLOGHISTORY",
			"MAXDATAFILES", "MAXINSTANCES",
			"ARCHIVELOG", "NOARCHIVELOG",
			"CHARACTER", "NATIONAL",
			"BIGFILE", "SMALLFILE",
			"UNDO", "PLUGGABLE":
			return true
		}
	}
	return false
}

// parseCreateControlfileStmt parses a CREATE CONTROLFILE statement.
// The CREATE keyword has already been consumed, and "CONTROLFILE" has been consumed by caller.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-CONTROLFILE.html
//
//	CREATE CONTROLFILE [ REUSE ]
//	    [ SET ] DATABASE database_name
//	    [ LOGFILE logfile_clause [, ...] ]
//	    { RESETLOGS | NORESETLOGS }
//	    [ DATAFILE file_specification [, ...] ]
//	    [ MAXLOGFILES integer ]
//	    [ MAXLOGMEMBERS integer ]
//	    [ MAXLOGHISTORY integer ]
//	    [ MAXDATAFILES integer ]
//	    [ MAXINSTANCES integer ]
//	    [ { ARCHIVELOG | NOARCHIVELOG } ]
//	    [ FORCE LOGGING ]
//	    [ SET STANDBY NOLOGGING FOR { DATA AVAILABILITY | LOAD PERFORMANCE } ]
//	    [ character_set_clause ]
func (p *Parser) parseCreateControlfileStmt(start int) nodes.StmtNode {
	// "CONTROLFILE" identifier already consumed by caller

	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_CONTROLFILE,
		Loc:        nodes.Loc{Start: start},
	}

	opts := &nodes.List{}

	// Optional REUSE
	if p.isIdentLike() && p.cur.Str == "REUSE" {
		optStart := p.pos()
		p.advance()
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "REUSE",
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
	}

	// [ SET ] DATABASE database_name
	isSet := false
	if p.cur.Type == kwSET {
		isSet = true
		p.advance()
	}
	if p.cur.Type == kwDATABASE {
		p.advance() // consume DATABASE
		dbName := ""
		if p.isIdentLike() || p.cur.Type == tokQIDENT {
			dbName = p.cur.Str
			p.advance()
		}
		key := "DATABASE"
		if isSet {
			key = "SET_DATABASE"
		}
		stmt.Name = &nodes.ObjectName{Name: dbName, Loc: nodes.Loc{Start: p.pos()}}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: key, Value: dbName,
		})
	}

	// Parse remaining clauses in a loop
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		optStart := p.pos()

		switch {
		// LOGFILE GROUP n 'file' SIZE ... [, GROUP n ...]
		case p.isIdentLike() && p.cur.Str == "LOGFILE":
			p.advance()
			logItems := &nodes.List{}
			for {
				lfStart := p.pos()
				groupNum := ""
				if p.cur.Type == kwGROUP {
					p.advance()
					if p.cur.Type == tokICONST {
						groupNum = p.cur.Str
						p.advance()
					}
				}
				var files []*nodes.DatafileClause
				for p.cur.Type == tokSCONST {
					df := p.parseDatafileClause()
					if df != nil {
						files = append(files, df)
					}
					if p.cur.Type == ',' {
						// peek: could be next file or next GROUP
						break
					}
				}
				fileList := &nodes.List{}
				for _, f := range files {
					fileList.Items = append(fileList.Items, f)
				}
				logItems.Items = append(logItems.Items, &nodes.DDLOption{
					Key: "LOGFILE_GROUP", Value: groupNum,
					Items: fileList,
					Loc:   nodes.Loc{Start: lfStart, End: p.pos()},
				})
				if p.cur.Type == ',' {
					p.advance()
					continue
				}
				break
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "LOGFILE", Items: logItems,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// RESETLOGS
		case p.isIdentLike() && p.cur.Str == "RESETLOGS":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "RESETLOGS",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// NORESETLOGS
		case p.isIdentLike() && p.cur.Str == "NORESETLOGS":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "NORESETLOGS",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// DATAFILE 'file' [, ...]
		case p.isIdentLike() && p.cur.Str == "DATAFILE":
			p.advance()
			fileList := &nodes.List{}
			for p.cur.Type == tokSCONST {
				df := p.parseDatafileClause()
				if df != nil {
					fileList.Items = append(fileList.Items, df)
				}
				if p.cur.Type == ',' {
					p.advance()
					continue
				}
				break
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "DATAFILE", Items: fileList,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// MAXLOGFILES, MAXLOGMEMBERS, MAXLOGHISTORY, MAXDATAFILES, MAXINSTANCES
		case p.isIdentLike() && (p.cur.Str == "MAXLOGFILES" || p.cur.Str == "MAXLOGMEMBERS" ||
			p.cur.Str == "MAXLOGHISTORY" || p.cur.Str == "MAXDATAFILES" || p.cur.Str == "MAXINSTANCES"):
			key := p.cur.Str
			p.advance()
			val := ""
			if p.cur.Type == tokICONST {
				val = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: key, Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// ARCHIVELOG
		case p.isIdentLike() && p.cur.Str == "ARCHIVELOG":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "ARCHIVELOG",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// NOARCHIVELOG
		case p.isIdentLike() && p.cur.Str == "NOARCHIVELOG":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "NOARCHIVELOG",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// FORCE LOGGING
		case p.cur.Type == kwFORCE:
			p.advance()
			if p.cur.Type == kwLOGGING {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "FORCE_LOGGING",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// SET STANDBY NOLOGGING FOR ...
		case p.cur.Type == kwSET:
			p.advance()
			if p.isIdentLike() && p.cur.Str == "STANDBY" {
				p.advance()
				if p.cur.Type == kwNOLOGGING {
					p.advance()
				}
				if p.cur.Type == kwFOR {
					p.advance()
				}
				val := ""
				if p.isIdentLike() && p.cur.Str == "DATA" {
					p.advance()
					if p.isIdentLike() && p.cur.Str == "AVAILABILITY" {
						p.advance()
					}
					val = "DATA AVAILABILITY"
				} else if p.isIdentLike() && p.cur.Str == "LOAD" {
					p.advance()
					if p.isIdentLike() && p.cur.Str == "PERFORMANCE" {
						p.advance()
					}
					val = "LOAD PERFORMANCE"
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "SET_STANDBY_NOLOGGING", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// CHARACTER SET charset
		case p.isIdentLike() && p.cur.Str == "CHARACTER":
			p.advance()
			if p.cur.Type == kwSET {
				p.advance()
			}
			val := ""
			if p.isIdentLike() || p.cur.Type == tokQIDENT {
				val = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "CHARACTER_SET", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		default:
			p.advance()
		}
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterDatabaseStmt parses an ALTER DATABASE statement.
// The ALTER keyword has already been consumed, and DATABASE is the current token.
//
// BNF: oracle/parser/bnf/ALTER-DATABASE.bnf
//
//	alter_database ::=
//	  ALTER DATABASE [ db_name ]
//	    { startup_clauses
//	    | recovery_clauses
//	    | database_file_clauses
//	    | logfile_clauses
//	    | controlfile_clauses
//	    | standby_database_clauses
//	    | default_settings_clauses
//	    | instance_clauses
//	    | security_clause
//	    | RENAME GLOBAL_NAME TO db_name
//	    | { ENABLE | DISABLE } BLOCK CHANGE TRACKING [ USING FILE 'filename' | USING FILE '+asm_diskgroup' ]
//	    | cdb_fleet_clauses
//	    | replay_upgrade_clause
//	    }
//
//	startup_clauses ::=
//	  { MOUNT [ STANDBY DATABASE | CLONE DATABASE ]
//	  | OPEN [ READ WRITE | READ ONLY ] [ RESETLOGS | NORESETLOGS ] [ UPGRADE | DOWNGRADE ]
//	  }
//
//	recovery_clauses ::= { BEGIN BACKUP | END BACKUP | general_recovery | managed_standby_recovery }
//	general_recovery ::= RECOVER [ AUTOMATIC | FROM 'location' ] { full_database_recovery | partial_database_recovery }
//	    [ LOGFILE 'filename' ] [ TEST ] [ ALLOW integer CORRUPTION ] [ CONTINUE [ DEFAULT ] | CANCEL ] [ parallel_clause ]
//	full_database_recovery ::= { DATABASE | STANDBY DATABASE [ UNTIL { CANCEL | TIME | CHANGE | CONSISTENT } ] [ USING BACKUP CONTROLFILE ] }
//	partial_database_recovery ::= { TABLESPACE ts [, ...] | DATAFILE { 'fn' | int } [, ...] }
//	managed_standby_recovery ::= RECOVER MANAGED STANDBY DATABASE [ USING ARCHIVED LOGFILE | USING CURRENT LOGFILE ]
//	    [ DISCONNECT [ FROM SESSION ] ] [ NODELAY ] [ UNTIL CHANGE int | UNTIL CONSISTENT ]
//	    [ USING INSTANCES { ALL | int } ] [ FINISH [ FORCE | WAIT | NOWAIT ] ] [ CANCEL [ IMMEDIATE | WAIT | NOWAIT ] ]
//	    [ TO LOGICAL STANDBY db_name [ KEEP IDENTITY ] ] [ parallel_clause ]
//
//	database_file_clauses ::=
//	  { RENAME FILE 'old' TO 'new' [, ...]
//	  | CREATE DATAFILE { 'fn' | int } AS { NEW | file_specification }
//	  | DATAFILE { 'fn' | int } { ONLINE | OFFLINE [ FOR DROP ] | RESIZE sz | END BACKUP | ENCRYPT | DECRYPT | autoextend }
//	  | TEMPFILE { 'fn' | int } { RESIZE sz | autoextend | DROP [ INCLUDING DATAFILES ] }
//	  | MOVE DATAFILE { 'fn' | ASM_fn | int } TO { 'fn' | ASM_fn } [ REUSE ] [ KEEP ]
//	  }
//
//	logfile_clauses ::=
//	  { ARCHIVELOG | MANUAL | NOARCHIVELOG | [ NO ] FORCE LOGGING
//	  | SET STANDBY NOLOGGING FOR { LOAD PERFORMANCE | DATA AVAILABILITY }
//	  | RENAME FILE ... TO ... | CLEAR LOGFILE logfile_descriptor [ UNARCHIVED | UNRECOVERABLE DATAFILE ]
//	  | add_logfile_clauses | drop_logfile_clauses | switch_logfile_clause | supplemental_db_logging }
//
//	controlfile_clauses ::=
//	  { CREATE [ PHYSICAL | LOGICAL ] STANDBY CONTROLFILE AS 'fn' [ REUSE ]
//	  | CREATE FAR SYNC INSTANCE CONTROLFILE AS 'fn' [ REUSE ]
//	  | BACKUP CONTROLFILE TO 'fn' [ REUSE ]
//	  | BACKUP CONTROLFILE TO TRACE [ AS 'fn' ] [ REUSE ] [ { RESETLOGS | NORESETLOGS } ]
//	  }
//
//	standby_database_clauses ::=
//	  { activate_standby_db_clause | maximize_standby_db_clause | register_logfile_clause
//	  | commit_switchover_clause | start_standby_clause | stop_standby_clause
//	  | convert_database_clause | switchover_clause | failover_clause }
//
//	default_settings_clauses ::=
//	  { DEFAULT EDITION = ed | DEFAULT TABLESPACE ts | DEFAULT [ LOCAL ] TEMPORARY TABLESPACE ts
//	  | SET DEFAULT { BIGFILE | SMALLFILE } TABLESPACE | flashback_mode_clause | undo_mode_clause | set_time_zone_clause }
//	instance_clauses ::= { ENABLE | DISABLE } RESTRICTED SESSION
//	security_clause ::= { prepare_clause | drop_mirror_copy | lost_write_protection }
//	cdb_fleet_clauses ::= { SET LEAD_CDB [=] { name | NONE } | SET LEAD_CDB_URI [=] { 'uri' | NONE } | SET PROPERTY { name = value } }
//	replay_upgrade_clause ::= { START REPLAY | STOP REPLAY }
func (p *Parser) parseAlterDatabaseStmt(start int) nodes.StmtNode {
	// DATABASE already consumed by caller
	stmt := &nodes.AdminDDLStmt{
		Action:     "ALTER",
		ObjectType: nodes.OBJECT_DATABASE,
		Loc:        nodes.Loc{Start: start},
	}

	// Optional database name — if next token is an identifier (not a keyword clause starter)
	if p.isIdentLike() || p.cur.Type == tokQIDENT {
		if !p.isDatabaseClauseKeyword() {
			stmt.Name = p.parseObjectName()
		}
	}

	opts := &nodes.List{}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		optStart := p.pos()

		switch {
		// ---- startup_clauses ----
		// MOUNT [ STANDBY | CLONE DATABASE ]
		case p.isIdentLike() && p.cur.Str == "MOUNT":
			p.advance()
			val := ""
			if p.isIdentLike() && p.cur.Str == "STANDBY" {
				p.advance()
				val = "STANDBY"
				if p.cur.Type == kwDATABASE {
					p.advance()
				}
			} else if p.isIdentLike() && p.cur.Str == "CLONE" {
				p.advance()
				val = "CLONE"
				if p.cur.Type == kwDATABASE {
					p.advance()
				}
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "MOUNT", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// OPEN { [ READ WRITE ] [ RESETLOGS | NORESETLOGS ] [ UPGRADE | DOWNGRADE ] | READ ONLY }
		case p.cur.Type == kwOPEN:
			p.advance()
			val := ""
			if p.cur.Type == kwREAD {
				p.advance()
				if p.cur.Type == kwONLY {
					p.advance()
					val = "READ ONLY"
				} else if p.cur.Type == kwWRITE {
					p.advance()
					val = "READ WRITE"
				}
			}
			// Optional RESETLOGS/NORESETLOGS
			if p.isIdentLike() && p.cur.Str == "RESETLOGS" {
				val += " RESETLOGS"
				p.advance()
			} else if p.isIdentLike() && p.cur.Str == "NORESETLOGS" {
				val += " NORESETLOGS"
				p.advance()
			}
			// Optional UPGRADE/DOWNGRADE
			if p.isIdentLike() && p.cur.Str == "UPGRADE" {
				val += " UPGRADE"
				p.advance()
			} else if p.isIdentLike() && p.cur.Str == "DOWNGRADE" {
				val += " DOWNGRADE"
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "OPEN", Value: strings.TrimSpace(val),
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// ---- recovery_clauses ----
		// RECOVER ...
		case p.isIdentLike() && p.cur.Str == "RECOVER":
			p.advance()
			p.parseAlterDatabaseRecoverClause(opts, optStart)

		// BEGIN BACKUP
		case p.cur.Type == kwBEGIN:
			p.advance()
			if p.isIdentLike() && p.cur.Str == "BACKUP" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "BEGIN_BACKUP",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// END BACKUP
		case p.cur.Type == kwEND:
			p.advance()
			if p.isIdentLike() && p.cur.Str == "BACKUP" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "END_BACKUP",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// ---- database_file_clauses ----
		// RENAME FILE 'f1' [, ...] TO 'f2' [, ...]
		case p.cur.Type == kwRENAME:
			p.advance()
			if p.cur.Type == kwFILE {
				p.advance()
				fromFiles := p.parseStringList()
				if p.cur.Type == kwTO {
					p.advance()
				}
				toFiles := p.parseStringList()
				items := &nodes.List{}
				for _, f := range fromFiles {
					items.Items = append(items.Items, &nodes.DDLOption{Key: "FROM", Value: f})
				}
				for _, f := range toFiles {
					items.Items = append(items.Items, &nodes.DDLOption{Key: "TO", Value: f})
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "RENAME_FILE", Items: items,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "GLOBAL_NAME" {
				// RENAME GLOBAL_NAME TO database.domain
				p.advance() // GLOBAL_NAME
				if p.cur.Type == kwTO {
					p.advance()
				}
				val := ""
				for p.isIdentLike() || p.cur.Type == tokQIDENT || p.cur.Type == '.' {
					val += p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "RENAME_GLOBAL_NAME", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// CREATE DATAFILE / CREATE [PHYSICAL|LOGICAL] STANDBY CONTROLFILE / CREATE FAR SYNC INSTANCE CONTROLFILE
		case p.cur.Type == kwCREATE:
			p.advance()
			if p.isIdentLike() && p.cur.Str == "DATAFILE" {
				p.advance()
				val := ""
				if p.cur.Type == tokSCONST || p.cur.Type == tokICONST {
					val = p.cur.Str
					p.advance()
				}
				asVal := ""
				if p.cur.Type == kwAS {
					p.advance()
					if p.cur.Type == tokSCONST {
						asVal = p.cur.Str
						p.advance()
					} else if p.isIdentLike() && p.cur.Str == "NEW" {
						asVal = "NEW"
						p.advance()
					}
				}
				items := &nodes.List{}
				if asVal != "" {
					items.Items = append(items.Items, &nodes.DDLOption{Key: "AS", Value: asVal})
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "CREATE_DATAFILE", Value: val,
					Items: items,
					Loc:   nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && (p.cur.Str == "PHYSICAL" || p.cur.Str == "LOGICAL") {
				// CREATE [PHYSICAL|LOGICAL] STANDBY CONTROLFILE AS 'filename' [REUSE]
				modifier := p.cur.Str
				p.advance()
				if p.isIdentLike() && p.cur.Str == "STANDBY" {
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "CONTROLFILE" {
					p.advance()
				}
				if p.cur.Type == kwAS {
					p.advance()
				}
				file := ""
				if p.cur.Type == tokSCONST {
					file = p.cur.Str
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "REUSE" {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "CREATE_STANDBY_CONTROLFILE", Value: modifier + " " + file,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "STANDBY" {
				// CREATE STANDBY CONTROLFILE AS 'filename' [REUSE]
				p.advance()
				if p.isIdentLike() && p.cur.Str == "CONTROLFILE" {
					p.advance()
				}
				if p.cur.Type == kwAS {
					p.advance()
				}
				file := ""
				if p.cur.Type == tokSCONST {
					file = p.cur.Str
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "REUSE" {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "CREATE_STANDBY_CONTROLFILE", Value: file,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "FAR" {
				// CREATE FAR SYNC INSTANCE CONTROLFILE AS 'filename' [REUSE]
				p.advance() // FAR
				if p.isIdentLike() && p.cur.Str == "SYNC" {
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "INSTANCE" {
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "CONTROLFILE" {
					p.advance()
				}
				if p.cur.Type == kwAS {
					p.advance()
				}
				file := ""
				if p.cur.Type == tokSCONST {
					file = p.cur.Str
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "REUSE" {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "CREATE_FAR_SYNC_CONTROLFILE", Value: file,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// DATAFILE 'file' { ONLINE | OFFLINE | RESIZE | AUTOEXTEND | END BACKUP }
		case p.isIdentLike() && p.cur.Str == "DATAFILE":
			p.advance()
			file := ""
			if p.cur.Type == tokSCONST {
				file = p.cur.Str
				p.advance()
			}
			p.parseAlterDatafileTempfileAction(opts, optStart, "DATAFILE", file)

		// TEMPFILE 'file' { ONLINE | OFFLINE | RESIZE | AUTOEXTEND | DROP | END BACKUP }
		case p.isIdentLike() && p.cur.Str == "TEMPFILE":
			p.advance()
			file := ""
			if p.cur.Type == tokSCONST {
				file = p.cur.Str
				p.advance()
			}
			p.parseAlterDatafileTempfileAction(opts, optStart, "TEMPFILE", file)

		// MOVE DATAFILE 'old' TO 'new'
		case p.isIdentLike() && p.cur.Str == "MOVE":
			p.advance()
			if p.isIdentLike() && p.cur.Str == "DATAFILE" {
				p.advance()
			}
			oldFile := ""
			if p.cur.Type == tokSCONST {
				oldFile = p.cur.Str
				p.advance()
			}
			if p.cur.Type == kwTO {
				p.advance()
			}
			newFile := ""
			if p.cur.Type == tokSCONST {
				newFile = p.cur.Str
				p.advance()
			}
			items := &nodes.List{}
			items.Items = append(items.Items, &nodes.DDLOption{Key: "TO", Value: newFile})
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "MOVE_DATAFILE", Value: oldFile,
				Items: items,
				Loc:   nodes.Loc{Start: optStart, End: p.pos()},
			})

		// ---- logfile_clauses ----
		// ADD [STANDBY] LOGFILE ...
		case p.cur.Type == kwADD:
			p.advance()
			standby := false
			if p.isIdentLike() && p.cur.Str == "STANDBY" {
				standby = true
				p.advance()
			}
			if p.isIdentLike() && p.cur.Str == "LOGFILE" {
				p.advance()
				prefix := "ADD_LOGFILE"
				if standby {
					prefix = "ADD_STANDBY_LOGFILE"
				}
				// ADD LOGFILE MEMBER 'file' TO GROUP n
				if p.isIdentLike() && p.cur.Str == "MEMBER" {
					p.advance()
					memberFile := ""
					if p.cur.Type == tokSCONST {
						memberFile = p.cur.Str
						p.advance()
					}
					if p.cur.Type == kwTO {
						p.advance()
					}
					groupNum := ""
					if p.cur.Type == kwGROUP {
						p.advance()
						if p.cur.Type == tokICONST {
							groupNum = p.cur.Str
							p.advance()
						}
					}
					opts.Items = append(opts.Items, &nodes.DDLOption{
						Key: prefix + "_MEMBER", Value: memberFile,
						Items: &nodes.List{Items: []nodes.Node{&nodes.DDLOption{Key: "GROUP", Value: groupNum}}},
						Loc:   nodes.Loc{Start: optStart, End: p.pos()},
					})
				} else {
					// ADD LOGFILE [INSTANCE 'inst'] [THREAD int] [GROUP n] 'file' SIZE ...
					logItems := &nodes.List{}
					// Optional INSTANCE 'name'
					if p.isIdentLike() && p.cur.Str == "INSTANCE" {
						p.advance()
						instVal := ""
						if p.cur.Type == tokSCONST {
							instVal = p.cur.Str
							p.advance()
						}
						logItems.Items = append(logItems.Items, &nodes.DDLOption{Key: "INSTANCE", Value: instVal})
					}
					// Optional THREAD int
					if p.isIdentLike() && p.cur.Str == "THREAD" {
						p.advance()
						threadVal := ""
						if p.cur.Type == tokICONST {
							threadVal = p.cur.Str
							p.advance()
						}
						logItems.Items = append(logItems.Items, &nodes.DDLOption{Key: "THREAD", Value: threadVal})
					}
					for {
						lfStart := p.pos()
						groupNum := ""
						if p.cur.Type == kwGROUP {
							p.advance()
							if p.cur.Type == tokICONST {
								groupNum = p.cur.Str
								p.advance()
							}
						}
						var files []*nodes.DatafileClause
						for p.cur.Type == tokSCONST {
							df := p.parseDatafileClause()
							if df != nil {
								files = append(files, df)
							}
							if p.cur.Type == ',' {
								break
							}
						}
						fileList := &nodes.List{}
						for _, f := range files {
							fileList.Items = append(fileList.Items, f)
						}
						logItems.Items = append(logItems.Items, &nodes.DDLOption{
							Key: "GROUP", Value: groupNum,
							Items: fileList,
							Loc:   nodes.Loc{Start: lfStart, End: p.pos()},
						})
						if p.cur.Type == ',' {
							p.advance()
							continue
						}
						break
					}
					opts.Items = append(opts.Items, &nodes.DDLOption{
						Key: prefix, Items: logItems,
						Loc: nodes.Loc{Start: optStart, End: p.pos()},
					})
				}
			} else if p.isIdentLike() && p.cur.Str == "SUPPLEMENTAL" {
				// ADD SUPPLEMENTAL LOG { DATA | supplemental_id_key_clause | supplemental_plsql_clause | DATA SUBSET DATABASE REPLICATION }
				p.advance() // SUPPLEMENTAL
				if p.cur.Type == kwLOG {
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "DATA" {
					p.advance()
				}
				val := ""
				// Check for DATA (column_key_clause) or DATA SUBSET DATABASE REPLICATION
				if p.cur.Type == '(' {
					p.advance() // (
					// Read column key type: PRIMARY KEY | UNIQUE | FOREIGN KEY | ALL | PL/SQL CALL
					var parts []string
					for p.cur.Type != ')' && p.cur.Type != tokEOF {
						if p.isIdentLike() || p.cur.Type == tokSCONST {
							parts = append(parts, p.cur.Str)
						} else if p.cur.Type == '/' {
							parts = append(parts, "/")
						}
						p.advance()
					}
					if p.cur.Type == ')' {
						p.advance()
					}
					val = strings.Join(parts, " ")
				} else if p.isIdentLike() && p.cur.Str == "SUBSET" {
					p.advance() // SUBSET
					if p.cur.Type == kwDATABASE {
						p.advance()
					}
					if p.isIdentLike() && p.cur.Str == "REPLICATION" {
						p.advance()
					}
					val = "SUBSET DATABASE REPLICATION"
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "ADD_SUPPLEMENTAL_LOG_DATA", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// DROP [STANDBY] LOGFILE ...
		case p.cur.Type == kwDROP:
			p.advance()
			standby := false
			if p.isIdentLike() && p.cur.Str == "STANDBY" {
				standby = true
				p.advance()
			}
			if p.isIdentLike() && p.cur.Str == "LOGFILE" {
				p.advance()
				prefix := "DROP_LOGFILE"
				if standby {
					prefix = "DROP_STANDBY_LOGFILE"
				}
				if p.isIdentLike() && p.cur.Str == "MEMBER" {
					p.advance()
					memberFile := ""
					if p.cur.Type == tokSCONST {
						memberFile = p.cur.Str
						p.advance()
					}
					opts.Items = append(opts.Items, &nodes.DDLOption{
						Key: prefix + "_MEMBER", Value: memberFile,
						Loc: nodes.Loc{Start: optStart, End: p.pos()},
					})
				} else {
					// DROP LOGFILE GROUP n
					groupNum := ""
					if p.cur.Type == kwGROUP {
						p.advance()
						if p.cur.Type == tokICONST {
							groupNum = p.cur.Str
							p.advance()
						}
					}
					opts.Items = append(opts.Items, &nodes.DDLOption{
						Key: prefix, Value: groupNum,
						Loc: nodes.Loc{Start: optStart, End: p.pos()},
					})
				}
			} else if p.isIdentLike() && p.cur.Str == "SUPPLEMENTAL" {
				// DROP SUPPLEMENTAL LOG DATA [sub-clause]
				p.advance()
				if p.cur.Type == kwLOG {
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "DATA" {
					p.advance()
				}
				val := ""
				if p.cur.Type == '(' {
					p.advance()
					var parts []string
					for p.cur.Type != ')' && p.cur.Type != tokEOF {
						if p.isIdentLike() || p.cur.Type == tokSCONST {
							parts = append(parts, p.cur.Str)
						} else if p.cur.Type == '/' {
							parts = append(parts, "/")
						}
						p.advance()
					}
					if p.cur.Type == ')' {
						p.advance()
					}
					val = strings.Join(parts, " ")
				} else if p.isIdentLike() && p.cur.Str == "SUBSET" {
					p.advance()
					if p.cur.Type == kwDATABASE {
						p.advance()
					}
					if p.isIdentLike() && p.cur.Str == "REPLICATION" {
						p.advance()
					}
					val = "SUBSET DATABASE REPLICATION"
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DROP_SUPPLEMENTAL_LOG_DATA", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "MIRROR" {
				// DROP MIRROR COPY
				p.advance() // MIRROR
				if p.isIdentLike() && p.cur.Str == "COPY" {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DROP_MIRROR_COPY",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// CLEAR [UNARCHIVED] LOGFILE logfile_descriptor [ UNARCHIVED | UNRECOVERABLE DATAFILE ]
		case p.isIdentLike() && p.cur.Str == "CLEAR":
			p.advance()
			unarchived := false
			if p.isIdentLike() && p.cur.Str == "UNARCHIVED" {
				unarchived = true
				p.advance()
			}
			if p.isIdentLike() && p.cur.Str == "LOGFILE" {
				p.advance()
			}
			groupNum := ""
			if p.cur.Type == kwGROUP {
				p.advance()
				if p.cur.Type == tokICONST {
					groupNum = p.cur.Str
					p.advance()
				}
			}
			key := "CLEAR_LOGFILE"
			if unarchived {
				key = "CLEAR_UNARCHIVED_LOGFILE"
			}
			// Optional trailing UNARCHIVED or UNRECOVERABLE DATAFILE
			if p.isIdentLike() && p.cur.Str == "UNARCHIVED" {
				p.advance()
				key = "CLEAR_UNARCHIVED_LOGFILE"
			} else if p.isIdentLike() && p.cur.Str == "UNRECOVERABLE" {
				p.advance()
				if p.isIdentLike() && p.cur.Str == "DATAFILE" {
					p.advance()
				}
				key = "CLEAR_LOGFILE_UNRECOVERABLE"
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: key, Value: groupNum,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// SWITCH { ALL LOGFILE | LOGFILE TO BLOCK SIZE integer }
		case p.isIdentLike() && p.cur.Str == "SWITCH":
			p.advance()
			if p.cur.Type == kwALL {
				p.advance()
				if p.isIdentLike() && p.cur.Str == "LOGFILE" {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "SWITCH_LOGFILE",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "LOGFILE" {
				p.advance() // LOGFILE
				val := ""
				if p.cur.Type == kwTO {
					p.advance() // TO
					if p.cur.Type == kwBLOCK {
						p.advance() // BLOCK
					}
					if p.cur.Type == kwSIZE {
						p.advance() // SIZE
					}
					if p.cur.Type == tokICONST {
						val = p.cur.Str
						p.advance()
					}
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "SWITCH_LOGFILE_BLOCK_SIZE", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// ---- controlfile_clauses ----
		// BACKUP CONTROLFILE TO { 'filename' [REUSE] | TRACE [AS 'filename' [REUSE]] }
		case p.isIdentLike() && p.cur.Str == "BACKUP":
			p.advance()
			if p.isIdentLike() && p.cur.Str == "CONTROLFILE" {
				p.advance()
			}
			if p.cur.Type == kwTO {
				p.advance()
			}
			if p.isIdentLike() && p.cur.Str == "TRACE" {
				p.advance()
				val := ""
				if p.cur.Type == kwAS {
					p.advance()
					if p.cur.Type == tokSCONST {
						val = p.cur.Str
						p.advance()
					}
				}
				if p.isIdentLike() && p.cur.Str == "REUSE" {
					p.advance()
				}
				// Optional RESETLOGS/NORESETLOGS
				if p.isIdentLike() && (p.cur.Str == "RESETLOGS" || p.cur.Str == "NORESETLOGS") {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "BACKUP_CONTROLFILE_TRACE", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.cur.Type == tokSCONST {
				file := p.cur.Str
				p.advance()
				if p.isIdentLike() && p.cur.Str == "REUSE" {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "BACKUP_CONTROLFILE", Value: file,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// ---- standby_database_clauses ----
		// ACTIVATE [PHYSICAL | LOGICAL] STANDBY DATABASE [FINISH APPLY]
		case p.isIdentLike() && p.cur.Str == "ACTIVATE":
			p.advance()
			val := ""
			if p.isIdentLike() && (p.cur.Str == "PHYSICAL" || p.cur.Str == "LOGICAL") {
				val = p.cur.Str
				p.advance()
			}
			if p.isIdentLike() && p.cur.Str == "STANDBY" {
				p.advance()
			}
			if p.cur.Type == kwDATABASE {
				p.advance()
			}
			if p.isIdentLike() && p.cur.Str == "FINISH" {
				p.advance()
				if p.isIdentLike() && p.cur.Str == "APPLY" {
					p.advance()
				}
				val += " FINISH APPLY"
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "ACTIVATE_STANDBY", Value: strings.TrimSpace(val),
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// SET { STANDBY ... | DEFAULT ... | TIME_ZONE | LEAD_CDB | LEAD_CDB_URI | PROPERTY | UNDO TABLESPACE }
		case p.cur.Type == kwSET:
			p.advance()
			if p.isIdentLike() && p.cur.Str == "STANDBY" {
				p.advance()
				if p.cur.Type == kwDATABASE {
					// SET STANDBY DATABASE TO MAXIMIZE ...
					p.advance()
					if p.cur.Type == kwTO {
						p.advance()
					}
					if p.isIdentLike() && p.cur.Str == "MAXIMIZE" {
						p.advance()
						val := ""
						if p.isIdentLike() {
							val = p.cur.Str
							p.advance()
						}
						opts.Items = append(opts.Items, &nodes.DDLOption{
							Key: "SET_STANDBY_MAXIMIZE", Value: val,
							Loc: nodes.Loc{Start: optStart, End: p.pos()},
						})
					}
				} else if p.cur.Type == kwNOLOGGING {
					// SET STANDBY NOLOGGING FOR { DATA AVAILABILITY | LOAD PERFORMANCE }
					p.advance() // NOLOGGING
					if p.cur.Type == kwFOR {
						p.advance()
					}
					val := ""
					if p.isIdentLike() && p.cur.Str == "DATA" {
						p.advance()
						if p.isIdentLike() && p.cur.Str == "AVAILABILITY" {
							p.advance()
						}
						val = "DATA AVAILABILITY"
					} else if p.isIdentLike() && p.cur.Str == "LOAD" {
						p.advance()
						if p.isIdentLike() && p.cur.Str == "PERFORMANCE" {
							p.advance()
						}
						val = "LOAD PERFORMANCE"
					}
					opts.Items = append(opts.Items, &nodes.DDLOption{
						Key: "SET_STANDBY_NOLOGGING", Value: val,
						Loc: nodes.Loc{Start: optStart, End: p.pos()},
					})
				} else {
					// SET STANDBY DATABASE TO MAXIMIZE ... (without DATABASE keyword)
					if p.cur.Type == kwTO {
						p.advance()
					}
					if p.isIdentLike() && p.cur.Str == "MAXIMIZE" {
						p.advance()
						val := ""
						if p.isIdentLike() {
							val = p.cur.Str
							p.advance()
						}
						opts.Items = append(opts.Items, &nodes.DDLOption{
							Key: "SET_STANDBY_MAXIMIZE", Value: val,
							Loc: nodes.Loc{Start: optStart, End: p.pos()},
						})
					}
				}
			} else if p.cur.Type == kwDEFAULT {
				p.advance()
				val := ""
				if p.isIdentLike() && (p.cur.Str == "BIGFILE" || p.cur.Str == "SMALLFILE") {
					val = p.cur.Str
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "TABLESPACE" {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "SET_DEFAULT_TABLESPACE_TYPE", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "TIME_ZONE" {
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
				}
				val := ""
				if p.cur.Type == tokSCONST || p.isIdentLike() {
					val = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "SET_TIME_ZONE", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "LEAD_CDB" {
				p.advance() // LEAD_CDB
				if p.cur.Type == '=' {
					p.advance()
				}
				val := ""
				if p.isIdentLike() || p.cur.Type == tokQIDENT {
					val = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "SET_LEAD_CDB", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "LEAD_CDB_URI" {
				p.advance() // LEAD_CDB_URI
				if p.cur.Type == '=' {
					p.advance()
				}
				val := ""
				if p.cur.Type == tokSCONST || p.isIdentLike() {
					val = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "SET_LEAD_CDB_URI", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "PROPERTY" {
				p.advance() // PROPERTY
				// property_name = property_value
				propName := ""
				if p.isIdentLike() || p.cur.Type == tokQIDENT {
					propName = p.cur.Str
					p.advance()
				}
				if p.cur.Type == '=' {
					p.advance()
				}
				propVal := ""
				if p.isIdentLike() || p.cur.Type == tokSCONST || p.cur.Type == tokICONST || p.cur.Type == tokQIDENT {
					propVal = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "SET_PROPERTY", Value: propName + "=" + propVal,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "UNDO" {
				// SET UNDO TABLESPACE = name
				p.advance() // UNDO
				if p.isIdentLike() && p.cur.Str == "TABLESPACE" {
					p.advance()
				}
				if p.cur.Type == '=' {
					p.advance()
				}
				val := ""
				if p.isIdentLike() || p.cur.Type == tokQIDENT {
					val = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "SET_UNDO_TABLESPACE", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// REGISTER [OR REPLACE] [PHYSICAL | LOGICAL] LOGFILE 'file' [, ...]
		case p.isIdentLike() && p.cur.Str == "REGISTER":
			p.advance()
			if p.isIdentLike() && (p.cur.Str == "PHYSICAL" || p.cur.Str == "LOGICAL") {
				p.advance()
			}
			if p.isIdentLike() && p.cur.Str == "LOGFILE" {
				p.advance()
			}
			fileList := &nodes.List{}
			for p.cur.Type == tokSCONST {
				df := p.parseDatafileClause()
				if df != nil {
					fileList.Items = append(fileList.Items, df)
				}
				if p.cur.Type == ',' {
					p.advance()
					continue
				}
				break
			}
			// OR REPLACE comes after file list
			if p.cur.Type == kwOR {
				p.advance()
				if p.cur.Type == kwREPLACE {
					p.advance()
				}
			}
			// Optional FOR logminer_session_name
			if p.cur.Type == kwFOR {
				p.advance()
				if p.isIdentLike() {
					p.advance()
				}
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "REGISTER_LOGFILE", Items: fileList,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// CONVERT TO [PHYSICAL | SNAPSHOT] STANDBY
		case p.isIdentLike() && p.cur.Str == "CONVERT":
			p.advance()
			if p.cur.Type == kwTO {
				p.advance()
			}
			val := ""
			if p.isIdentLike() && (p.cur.Str == "PHYSICAL" || p.cur.Str == "SNAPSHOT") {
				val = p.cur.Str
				p.advance()
			}
			if p.isIdentLike() && p.cur.Str == "STANDBY" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "CONVERT_TO_STANDBY", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// ---- default_settings_clauses ----
		// DEFAULT TABLESPACE | DEFAULT [LOCAL] TEMPORARY TABLESPACE | DEFAULT EDITION
		case p.cur.Type == kwDEFAULT:
			p.advance()
			if p.cur.Type == kwLOCAL {
				// DEFAULT LOCAL TEMPORARY TABLESPACE name
				p.advance() // LOCAL
				if p.cur.Type == kwTEMPORARY {
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "TABLESPACE" {
					p.advance()
				}
				tsName := ""
				if p.isIdentLike() || p.cur.Type == tokQIDENT {
					tsName = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DEFAULT_LOCAL_TEMPORARY_TABLESPACE", Value: tsName,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.cur.Type == kwTEMPORARY {
				p.advance()
				if p.isIdentLike() && p.cur.Str == "TABLESPACE" {
					p.advance()
				}
				tsName := ""
				if p.isIdentLike() || p.cur.Type == tokQIDENT {
					tsName = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DEFAULT_TEMPORARY_TABLESPACE", Value: tsName,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "TABLESPACE" {
				p.advance()
				tsName := ""
				if p.isIdentLike() || p.cur.Type == tokQIDENT {
					tsName = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DEFAULT_TABLESPACE", Value: tsName,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "EDITION" {
				p.advance()
				edName := ""
				if p.isIdentLike() || p.cur.Type == tokQIDENT {
					edName = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DEFAULT_EDITION", Value: edName,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else {
				p.advance()
			}

		// ENABLE { BLOCK CHANGE TRACKING | INSTANCE | RESTRICTED SESSION | LOST WRITE PROTECTION }
		case p.cur.Type == kwENABLE:
			p.advance()
			if p.isIdentLike() && p.cur.Str == "INSTANCE" {
				p.advance()
				val := ""
				if p.cur.Type == tokSCONST {
					val = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "ENABLE_INSTANCE", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.cur.Type == kwBLOCK {
				p.advance() // BLOCK
				if p.isIdentLike() && p.cur.Str == "CHANGE" {
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "TRACKING" {
					p.advance()
				}
				// Optional USING FILE 'filename' [REUSE]
				if p.cur.Type == kwUSING {
					p.advance()
					if p.cur.Type == kwFILE {
						p.advance()
					}
					if p.cur.Type == tokSCONST {
						p.advance()
					}
					if p.isIdentLike() && p.cur.Str == "REUSE" {
						p.advance()
					}
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "ENABLE_BLOCK_CHANGE_TRACKING",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "RESTRICTED" {
				p.advance() // RESTRICTED
				if p.cur.Type == kwSESSION {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "ENABLE_RESTRICTED_SESSION",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "LOST" {
				p.advance() // LOST
				if p.cur.Type == kwWRITE {
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "PROTECTION" {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "ENABLE_LOST_WRITE_PROTECTION",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// DISABLE { BLOCK CHANGE TRACKING | INSTANCE | RESTRICTED SESSION | LOST WRITE PROTECTION }
		case p.cur.Type == kwDISABLE:
			p.advance()
			if p.isIdentLike() && p.cur.Str == "INSTANCE" {
				p.advance()
				val := ""
				if p.cur.Type == tokSCONST {
					val = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DISABLE_INSTANCE", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.cur.Type == kwBLOCK {
				p.advance()
				if p.isIdentLike() && p.cur.Str == "CHANGE" {
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "TRACKING" {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DISABLE_BLOCK_CHANGE_TRACKING",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "RESTRICTED" {
				p.advance() // RESTRICTED
				if p.cur.Type == kwSESSION {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DISABLE_RESTRICTED_SESSION",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "LOST" {
				p.advance() // LOST
				if p.cur.Type == kwWRITE {
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "PROTECTION" {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DISABLE_LOST_WRITE_PROTECTION",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// FLASHBACK { ON | OFF }
		case p.cur.Type == kwFLASHBACK:
			p.advance()
			val := ""
			if p.cur.Type == kwON || (p.isIdentLike() && p.cur.Str == "ON") {
				val = "ON"
				p.advance()
			} else if p.isIdentLike() && p.cur.Str == "OFF" {
				val = "OFF"
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "FLASHBACK", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// GUARD { ALL | STANDBY | NONE }
		case p.isIdentLike() && p.cur.Str == "GUARD":
			p.advance()
			val := ""
			if p.cur.Type == kwALL {
				val = "ALL"
				p.advance()
			} else if p.isIdentLike() && (p.cur.Str == "STANDBY" || p.cur.Str == "NONE") {
				val = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "GUARD", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// FORCE LOGGING
		case p.cur.Type == kwFORCE:
			p.advance()
			if p.cur.Type == kwLOGGING {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "FORCE_LOGGING",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// SUSPEND / RESUME
		case p.isIdentLike() && p.cur.Str == "SUSPEND":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "SUSPEND",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		case p.isIdentLike() && p.cur.Str == "RESUME":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "RESUME",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// PREPARE { MIRROR COPY | TO SWITCHOVER }
		case p.isIdentLike() && p.cur.Str == "PREPARE":
			p.advance()
			if p.cur.Type == kwTO {
				// PREPARE TO SWITCHOVER TO { LOGICAL STANDBY | PRIMARY DATABASE }
				p.advance() // TO
				if p.isIdentLike() && p.cur.Str == "SWITCHOVER" {
					p.advance()
				}
				if p.cur.Type == kwTO {
					p.advance()
				}
				val := ""
				for p.isIdentLike() || p.cur.Type == kwDATABASE {
					if val != "" {
						val += " "
					}
					val += p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "PREPARE_SWITCHOVER", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else {
				// PREPARE MIRROR COPY name ...
				for p.cur.Type != ';' && p.cur.Type != tokEOF {
					if p.isIdentLike() || p.cur.Type == tokSCONST || p.cur.Type == tokQIDENT {
						p.advance()
					} else {
						break
					}
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "PREPARE",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// COMMIT TO SWITCHOVER TO { [PHYSICAL] STANDBY | LOGICAL STANDBY | PRIMARY DATABASE }
		case p.cur.Type == kwCOMMIT:
			p.advance() // COMMIT
			if p.cur.Type == kwTO {
				p.advance() // TO
			}
			if p.isIdentLike() && p.cur.Str == "SWITCHOVER" {
				p.advance()
			}
			if p.cur.Type == kwTO {
				p.advance()
			}
			val := ""
			for p.isIdentLike() || p.cur.Type == kwDATABASE {
				if val != "" {
					val += " "
				}
				val += p.cur.Str
				p.advance()
			}
			// Optional WITH/WITHOUT SESSION SHUTDOWN, WAIT/NOWAIT, CANCEL
			for p.cur.Type != ';' && p.cur.Type != tokEOF {
				if p.isIdentLike() || p.cur.Type == kwSESSION || p.cur.Type == kwWAIT || p.cur.Type == kwNOWAIT {
					p.advance()
				} else {
					break
				}
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "COMMIT_SWITCHOVER", Value: strings.TrimSpace(val),
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// START { LOGICAL STANDBY APPLY | REPLAY }
		case p.cur.Type == kwSTART:
			p.advance()
			if p.isIdentLike() && p.cur.Str == "LOGICAL" {
				p.advance() // LOGICAL
				if p.isIdentLike() && p.cur.Str == "STANDBY" {
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "APPLY" {
					p.advance()
				}
				// Optional modifiers: IMMEDIATE, NODELAY, INITIAL, NEW PRIMARY, SKIP FAILED, FINISH
				val := ""
				for p.cur.Type != ';' && p.cur.Type != tokEOF {
					if p.isIdentLike() {
						if val != "" {
							val += " "
						}
						val += p.cur.Str
						p.advance()
					} else {
						break
					}
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "START_LOGICAL_STANDBY_APPLY", Value: strings.TrimSpace(val),
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "REPLAY" {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "START_REPLAY",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// STOP { LOGICAL STANDBY APPLY | REPLAY }
		case p.isIdentLike() && p.cur.Str == "STOP":
			p.advance()
			if p.isIdentLike() && p.cur.Str == "LOGICAL" {
				p.advance()
				if p.isIdentLike() && p.cur.Str == "STANDBY" {
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "APPLY" {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "STOP_LOGICAL_STANDBY_APPLY",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "REPLAY" {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "STOP_REPLAY",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// SWITCHOVER TO { PRIMARY | PHYSICAL STANDBY | LOGICAL STANDBY } [VERIFY] [FORCE] target_db_name
		case p.isIdentLike() && p.cur.Str == "SWITCHOVER":
			p.advance()
			if p.cur.Type == kwTO {
				p.advance()
			}
			val := ""
			// Parse target type + modifiers + db name
			for p.cur.Type != ';' && p.cur.Type != tokEOF {
				if p.isIdentLike() || p.cur.Type == tokQIDENT {
					if val != "" {
						val += " "
					}
					val += p.cur.Str
					p.advance()
				} else {
					break
				}
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "SWITCHOVER_TO", Value: strings.TrimSpace(val),
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// FAILOVER TO { PRIMARY | PHYSICAL STANDBY | LOGICAL STANDBY } [FORCE] target_db_name
		case p.isIdentLike() && p.cur.Str == "FAILOVER":
			p.advance()
			if p.cur.Type == kwTO {
				p.advance()
			}
			val := ""
			for p.cur.Type != ';' && p.cur.Type != tokEOF {
				if p.isIdentLike() || p.cur.Type == tokQIDENT {
					if val != "" {
						val += " "
					}
					val += p.cur.Str
					p.advance()
				} else {
					break
				}
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "FAILOVER_TO", Value: strings.TrimSpace(val),
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// ARCHIVELOG / NOARCHIVELOG / MANUAL (logfile_clauses)
		case p.isIdentLike() && p.cur.Str == "ARCHIVELOG":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "ARCHIVELOG",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		case p.isIdentLike() && p.cur.Str == "NOARCHIVELOG":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "NOARCHIVELOG",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		case p.isIdentLike() && p.cur.Str == "MANUAL":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "MANUAL",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		default:
			// Unknown token, advance to avoid infinite loop
			p.advance()
		}
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterDatabaseRecoverClause parses the RECOVER sub-clauses of ALTER DATABASE.
//
//	RECOVER [ AUTOMATIC ] [ FROM 'location' ] DATABASE
//	  [ UNTIL { CANCEL | TIME date | CHANGE integer } ]
//	  [ USING BACKUP CONTROLFILE ]
//	RECOVER [ AUTOMATIC ] [ FROM 'location' ] [ STANDBY ] DATAFILE 'file' [, ...]
//	RECOVER [ AUTOMATIC ] [ FROM 'location' ] [ STANDBY ] TABLESPACE name [, ...]
//	RECOVER MANAGED STANDBY DATABASE { CANCEL | DISCONNECT [FROM SESSION] | FINISH | ... }
func (p *Parser) parseAlterDatabaseRecoverClause(opts *nodes.List, optStart int) {
	// RECOVER already consumed
	val := ""

	// AUTOMATIC
	if p.isIdentLike() && p.cur.Str == "AUTOMATIC" {
		val = "AUTOMATIC "
		p.advance()
	}

	// MANAGED STANDBY DATABASE ...
	if p.isIdentLike() && p.cur.Str == "MANAGED" {
		p.advance()
		if p.isIdentLike() && p.cur.Str == "STANDBY" {
			p.advance()
		}
		if p.cur.Type == kwDATABASE {
			p.advance()
		}
		// Parse managed standby options
		var parts []string
		for p.cur.Type != ';' && p.cur.Type != tokEOF {
			switch {
			case p.isIdentLike() && p.cur.Str == "USING":
				p.advance()
				if p.isIdentLike() && (p.cur.Str == "ARCHIVED" || p.cur.Str == "CURRENT") {
					parts = append(parts, "USING "+p.cur.Str)
					p.advance()
					if p.isIdentLike() && p.cur.Str == "LOGFILE" {
						parts[len(parts)-1] += " LOGFILE"
						p.advance()
					}
				} else if p.isIdentLike() && p.cur.Str == "INSTANCES" {
					p.advance()
					val := "USING INSTANCES"
					if p.cur.Type == kwALL {
						val += " ALL"
						p.advance()
					} else if p.cur.Type == tokICONST {
						val += " " + p.cur.Str
						p.advance()
					}
					parts = append(parts, val)
				}
			case p.isIdentLike() && p.cur.Str == "DISCONNECT":
				p.advance()
				val := "DISCONNECT"
				if p.cur.Type == kwFROM {
					p.advance()
					if p.cur.Type == kwSESSION {
						p.advance()
					}
					val = "DISCONNECT FROM SESSION"
				}
				parts = append(parts, val)
			case p.isIdentLike() && p.cur.Str == "NODELAY":
				parts = append(parts, "NODELAY")
				p.advance()
			case p.isIdentLike() && p.cur.Str == "UNTIL":
				p.advance()
				val := "UNTIL"
				if p.isIdentLike() && p.cur.Str == "CHANGE" {
					p.advance()
					val += " CHANGE"
					if p.cur.Type == tokICONST {
						val += " " + p.cur.Str
						p.advance()
					}
				} else if p.isIdentLike() && p.cur.Str == "CONSISTENT" {
					p.advance()
					val += " CONSISTENT"
				}
				parts = append(parts, val)
			case p.isIdentLike() && p.cur.Str == "FINISH":
				p.advance()
				val := "FINISH"
				if p.cur.Type == kwFORCE {
					val += " FORCE"
					p.advance()
				} else if p.cur.Type == kwWAIT {
					val += " WAIT"
					p.advance()
				} else if p.cur.Type == kwNOWAIT {
					val += " NOWAIT"
					p.advance()
				}
				parts = append(parts, val)
			case p.isIdentLike() && p.cur.Str == "CANCEL":
				p.advance()
				val := "CANCEL"
				if p.isIdentLike() && p.cur.Str == "IMMEDIATE" {
					val += " IMMEDIATE"
					p.advance()
				} else if p.cur.Type == kwWAIT {
					val += " WAIT"
					p.advance()
				} else if p.cur.Type == kwNOWAIT {
					val += " NOWAIT"
					p.advance()
				}
				parts = append(parts, val)
			case p.isIdentLike() && p.cur.Str == "REGISTER":
				// REGISTER LOGFILE in managed standby context
				p.advance()
				if p.isIdentLike() && p.cur.Str == "LOGFILE" {
					p.advance()
				}
				parts = append(parts, "REGISTER LOGFILE")
				// Parse files
				for p.cur.Type == tokSCONST {
					p.advance()
					if p.cur.Type == ',' {
						p.advance()
					}
				}
			case p.cur.Type == kwTO:
				p.advance()
				if p.isIdentLike() && p.cur.Str == "LOGICAL" {
					p.advance()
					if p.isIdentLike() && p.cur.Str == "STANDBY" {
						p.advance()
					}
					parts = append(parts, "TO LOGICAL STANDBY")
				}
			default:
				goto doneManaged
			}
		}
	doneManaged:
		action := strings.Join(parts, " ")
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "RECOVER_MANAGED_STANDBY", Value: action,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
		return
	}

	// FROM 'location'
	if p.cur.Type == kwFROM {
		p.advance()
		if p.cur.Type == tokSCONST {
			p.advance()
		}
	}

	// DATABASE
	if p.cur.Type == kwDATABASE {
		p.advance()
		// Optional UNTIL { CANCEL | TIME date | CHANGE integer | CONSISTENT }
		if p.isIdentLike() && p.cur.Str == "UNTIL" {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "CANCEL" {
				p.advance()
			} else if p.isIdentLike() && p.cur.Str == "TIME" {
				p.advance()
				if p.cur.Type == tokSCONST {
					p.advance()
				}
			} else if p.isIdentLike() && p.cur.Str == "CHANGE" {
				p.advance()
				if p.cur.Type == tokICONST {
					p.advance()
				}
			} else if p.isIdentLike() && p.cur.Str == "CONSISTENT" {
				p.advance()
			}
		}
		// Optional USING BACKUP CONTROLFILE
		if p.cur.Type == kwUSING {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "BACKUP" {
				p.advance()
			}
			if p.isIdentLike() && p.cur.Str == "CONTROLFILE" {
				p.advance()
			}
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "RECOVER_DATABASE", Value: strings.TrimSpace(val),
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
		return
	}

	// STANDBY [ DATABASE [ UNTIL ... ] [ USING BACKUP CONTROLFILE ] ]
	if p.isIdentLike() && p.cur.Str == "STANDBY" {
		p.advance()
		if p.cur.Type == kwDATABASE {
			p.advance()
			// Optional UNTIL { CANCEL | TIME | CHANGE | CONSISTENT }
			if p.isIdentLike() && p.cur.Str == "UNTIL" {
				p.advance()
				if p.isIdentLike() && p.cur.Str == "CANCEL" {
					p.advance()
				} else if p.isIdentLike() && p.cur.Str == "TIME" {
					p.advance()
					if p.cur.Type == tokSCONST {
						p.advance()
					}
				} else if p.isIdentLike() && p.cur.Str == "CHANGE" {
					p.advance()
					if p.cur.Type == tokICONST {
						p.advance()
					}
				} else if p.isIdentLike() && p.cur.Str == "CONSISTENT" {
					p.advance()
				}
			}
			if p.cur.Type == kwUSING {
				p.advance()
				if p.isIdentLike() && p.cur.Str == "BACKUP" {
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "CONTROLFILE" {
					p.advance()
				}
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "RECOVER_STANDBY_DATABASE", Value: strings.TrimSpace(val),
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})
			return
		}
	}

	// DATAFILE 'file' [, ...]
	if p.isIdentLike() && p.cur.Str == "DATAFILE" {
		p.advance()
		file := ""
		if p.cur.Type == tokSCONST {
			file = p.cur.Str
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "RECOVER_DATAFILE", Value: file,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
		return
	}

	// TABLESPACE name [, ...]
	if p.isIdentLike() && p.cur.Str == "TABLESPACE" {
		p.advance()
		tsName := ""
		if p.isIdentLike() || p.cur.Type == tokQIDENT {
			tsName = p.cur.Str
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "RECOVER_TABLESPACE", Value: tsName,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
		return
	}

	// Fallback: generic RECOVER
	opts.Items = append(opts.Items, &nodes.DDLOption{
		Key: "RECOVER", Value: strings.TrimSpace(val),
		Loc: nodes.Loc{Start: optStart, End: p.pos()},
	})
}

// parseAlterDatafileTempfileAction parses actions after DATAFILE/TEMPFILE 'file' in ALTER DATABASE.
//
//	{ ONLINE | OFFLINE [ FOR DROP ] | RESIZE size_clause | AUTOEXTEND ... | END BACKUP | DROP [INCLUDING DATAFILES] }
func (p *Parser) parseAlterDatafileTempfileAction(opts *nodes.List, optStart int, prefix string, file string) {
	switch {
	case p.cur.Type == kwONLINE:
		p.advance()
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: prefix + "_ONLINE", Value: file,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
	case p.cur.Type == kwOFFLINE:
		p.advance()
		// Optional FOR DROP
		if p.cur.Type == kwFOR {
			p.advance()
			if p.cur.Type == kwDROP {
				p.advance()
			}
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: prefix + "_OFFLINE", Value: file,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
	case p.isIdentLike() && p.cur.Str == "RESIZE":
		p.advance()
		size := p.parseSizeValue()
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: prefix + "_RESIZE", Value: file,
			Items: &nodes.List{Items: []nodes.Node{&nodes.DDLOption{Key: "SIZE", Value: size}}},
			Loc:   nodes.Loc{Start: optStart, End: p.pos()},
		})
	case p.isIdentLike() && p.cur.Str == "AUTOEXTEND":
		ac := p.parseAutoextendClause()
		items := &nodes.List{}
		if ac != nil {
			items.Items = append(items.Items, ac)
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: prefix + "_AUTOEXTEND", Value: file,
			Items: items,
			Loc:   nodes.Loc{Start: optStart, End: p.pos()},
		})
	case p.cur.Type == kwEND:
		p.advance()
		if p.isIdentLike() && p.cur.Str == "BACKUP" {
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: prefix + "_END_BACKUP", Value: file,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
	case p.cur.Type == kwDROP:
		p.advance()
		if p.isIdentLike() && p.cur.Str == "INCLUDING" {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "DATAFILES" {
				p.advance()
			}
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: prefix + "_DROP", Value: file,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
	case p.isIdentLike() && p.cur.Str == "ENCRYPT":
		p.advance()
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: prefix + "_ENCRYPT", Value: file,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
	case p.isIdentLike() && p.cur.Str == "DECRYPT":
		p.advance()
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: prefix + "_DECRYPT", Value: file,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
	default:
		// unknown action, still record the file reference
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: prefix, Value: file,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
	}
}

// parseStringList parses a comma-separated list of string constants.
func (p *Parser) parseStringList() []string {
	var result []string
	for p.cur.Type == tokSCONST {
		result = append(result, p.cur.Str)
		p.advance()
		if p.cur.Type == ',' {
			p.advance()
			continue
		}
		break
	}
	return result
}

// parseAlterDatabaseDictionaryStmt parses an ALTER DATABASE DICTIONARY statement.
// The ALTER keyword has already been consumed, DATABASE and DICTIONARY are consumed by caller.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ALTER-DATABASE-DICTIONARY.html
//
//	ALTER DATABASE DICTIONARY
//	    { ENCRYPT CREDENTIALS
//	    | REKEY CREDENTIALS
//	    | DELETE CREDENTIALS KEY }
func (p *Parser) parseAlterDatabaseDictionaryStmt(start int) nodes.StmtNode {
	// DATABASE and DICTIONARY already consumed by caller
	stmt := &nodes.AdminDDLStmt{
		Action:     "ALTER",
		ObjectType: nodes.OBJECT_DATABASE_DICTIONARY,
		Loc:        nodes.Loc{Start: start},
	}

	opts := &nodes.List{}
	optStart := p.pos()

	switch {
	// ENCRYPT CREDENTIALS
	case p.isIdentLike() && p.cur.Str == "ENCRYPT":
		p.advance()
		if p.isIdentLike() && p.cur.Str == "CREDENTIALS" {
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "ENCRYPT_CREDENTIALS",
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})

	// REKEY CREDENTIALS
	case p.isIdentLike() && p.cur.Str == "REKEY":
		p.advance()
		if p.isIdentLike() && p.cur.Str == "CREDENTIALS" {
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "REKEY_CREDENTIALS",
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})

	// DELETE CREDENTIALS KEY
	case p.cur.Type == kwDELETE:
		p.advance()
		if p.isIdentLike() && p.cur.Str == "CREDENTIALS" {
			p.advance()
		}
		if p.cur.Type == kwKEY {
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "DELETE_CREDENTIALS_KEY",
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// isDatabaseClauseKeyword returns true if the current token is a keyword
// that starts an ALTER DATABASE clause (not a database name).
func (p *Parser) isDatabaseClauseKeyword() bool {
	switch p.cur.Type {
	case kwOPEN, kwFORCE, kwSET, kwDEFAULT, kwADD, kwDROP,
		kwENABLE, kwDISABLE, kwRENAME, kwCREATE,
		kwFLASHBACK, kwNOLOGGING, kwLOGGING,
		kwBEGIN, kwEND, kwONLINE, kwOFFLINE,
		kwSTART, kwCOMMIT:
		return true
	}
	if p.isIdentLike() {
		switch p.cur.Str {
		case "MOUNT", "RECOVER", "RESETLOGS", "NORESETLOGS",
			"ARCHIVELOG", "NOARCHIVELOG", "ACTIVATE", "GUARD",
			"STANDBY", "CLEAR", "CONVERT", "DISMOUNT", "PREPARE",
			"BACKUP", "DATAFILE", "TEMPFILE", "MOVE", "SWITCH",
			"REGISTER", "SUSPEND", "RESUME", "STOP", "SWITCHOVER",
			"FAILOVER", "MANUAL", "NO":
			return true
		}
	}
	return false
}

// parseOptionalAutoextend parses an optional AUTOEXTEND clause if present, returning nil otherwise.
func (p *Parser) parseOptionalAutoextend() *nodes.AutoextendClause {
	if p.isIdentLike() && p.cur.Str == "AUTOEXTEND" {
		return p.parseAutoextendClause()
	}
	return nil
}

// objectNameStr returns a string representation of an ObjectName.
func objectNameStr(n *nodes.ObjectName) string {
	if n == nil {
		return ""
	}
	s := ""
	if n.Schema != "" {
		s = n.Schema + "."
	}
	s += n.Name
	if n.DBLink != "" {
		s += "@" + n.DBLink
	}
	return s
}

// parseCreatePluggableDatabaseStmt parses a CREATE PLUGGABLE DATABASE statement.
// The CREATE keyword has been consumed. PLUGGABLE DATABASE has been consumed by caller.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-PLUGGABLE-DATABASE.html
//
//	CREATE PLUGGABLE DATABASE pdb_name
//	  { create_pdb_from_seed | create_pdb_clone | create_pdb_from_xml
//	  | create_pdb_from_mirror_copy | create_pdb_decrypt_from_xml }
func (p *Parser) parseCreatePluggableDatabaseStmt(start int) nodes.StmtNode {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_PLUGGABLE_DATABASE,
		Loc:        nodes.Loc{Start: start},
	}

	// Parse PDB name
	stmt.Name = p.parseObjectName()

	opts := &nodes.List{}

	// Determine which variant: FROM, USING, AS, or from_seed (default)
	switch {
	case p.cur.Type == kwFROM:
		// create_pdb_clone or create_pdb_from_mirror_copy
		p.advance() // consume FROM
		if p.isIdentLike() && p.cur.Str == "SEED" {
			p.advance() // consume SEED
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "FROM_SEED"})
		} else {
			// FROM src_pdb_name[@dblink]
			srcName := p.parseObjectName()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key:   "FROM",
				Value: objectNameStr(srcName),
			})
		}
	case p.cur.Type == kwAS:
		p.advance() // consume AS
		if p.isIdentLike() && p.cur.Str == "PROXY" {
			// AS PROXY FROM src_pdb_name@dblink
			p.advance() // consume PROXY
			if p.cur.Type == kwFROM {
				p.advance() // consume FROM
			}
			srcName := p.parseObjectName()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key:   "AS_PROXY_FROM",
				Value: objectNameStr(srcName),
			})
		} else if p.isIdentLike() && p.cur.Str == "APPLICATION" {
			// AS APPLICATION CONTAINER
			p.advance() // consume APPLICATION
			if p.isIdentLike() && p.cur.Str == "CONTAINER" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "AS_APPLICATION_CONTAINER"})
		} else if p.isIdentLikeStr("SEED") {
			// AS SEED
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "AS_SEED"})
		} else if p.isIdentLike() && p.cur.Str == "CLONE" {
			// AS CLONE (from XML variant) — but USING should follow
			p.advance() // consume CLONE
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "AS_CLONE"})
		}
	case p.isIdentLike() && p.cur.Str == "USING":
		// create_pdb_from_xml / create_pdb_decrypt_from_xml
		p.advance() // consume USING
		if p.cur.Type == tokSCONST {
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key:   "USING",
				Value: p.cur.Str,
			})
			p.advance()
		}
	}

	// Parse common PDB creation clauses
	p.parsePDBCreationClauses(opts)

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.pos()
	return stmt
}

// parsePDBCreationClauses parses common clauses for CREATE PLUGGABLE DATABASE variants.
func (p *Parser) parsePDBCreationClauses(opts *nodes.List) {
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		optStart := p.pos()

		switch {
		// USING 'filename' (for XML variant after AS CLONE)
		case p.isIdentLike() && p.cur.Str == "USING":
			p.advance()
			if p.cur.Type == tokSCONST {
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "USING", Value: p.cur.Str,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
				p.advance()
			}

		// ADMIN USER ... IDENTIFIED BY ...
		case p.isIdentLike() && p.cur.Str == "ADMIN":
			p.advance() // ADMIN
			if p.cur.Type == kwUSER {
				p.advance() // USER
			}
			userName := ""
			if p.isIdentLike() || p.cur.Type == tokIDENT {
				userName = p.cur.Str
				p.advance()
			}
			if p.cur.Type == kwIDENTIFIED {
				p.advance() // IDENTIFIED
				if p.isIdentLike() && p.cur.Str == "BY" {
					p.advance() // BY
				}
			}
			password := ""
			if p.isIdentLike() || p.cur.Type == tokSCONST || p.cur.Type == tokIDENT {
				password = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "ADMIN_USER", Value: userName + ":" + password,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// ROLES = (role, ...)
		case p.cur.Type == kwROLE || (p.isIdentLike() && p.cur.Str == "ROLES"):
			p.advance() // ROLES
			if p.cur.Type == '=' {
				p.advance()
			}
			p.parsePDBParenList(opts, "ROLES", optStart)

		// PARALLEL [integer]
		case p.isIdentLike() && p.cur.Str == "PARALLEL":
			p.advance()
			val := ""
			if p.cur.Type == tokICONST {
				val = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "PARALLEL", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// DEFAULT TABLESPACE / DEFAULT TEMPORARY TABLESPACE / DEFAULT EDITION
		case p.cur.Type == kwDEFAULT:
			p.advance() // DEFAULT
			if p.isIdentLike() && p.cur.Str == "EDITION" {
				p.advance() // EDITION
				edName := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT {
					edName = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DEFAULT_EDITION", Value: edName,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "TEMPORARY" {
				p.advance() // TEMPORARY
				if p.cur.Type == kwTABLESPACE {
					p.advance() // TABLESPACE
				}
				tsName := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
					tsName = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DEFAULT_TEMPORARY_TABLESPACE", Value: tsName,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.cur.Type == kwTABLESPACE {
				p.advance() // TABLESPACE
				tsName := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
					tsName = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DEFAULT_TABLESPACE", Value: tsName,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else {
				p.advance() // skip unknown
			}

		// FILE_NAME_CONVERT = (...)  | FILE_NAME_CONVERT = NONE
		case p.isIdentLike() && p.cur.Str == "FILE_NAME_CONVERT":
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			p.parsePDBConvertClause(opts, "FILE_NAME_CONVERT", optStart)

		// SERVICE_NAME_CONVERT = (...) | SERVICE_NAME_CONVERT = NONE
		case p.isIdentLike() && p.cur.Str == "SERVICE_NAME_CONVERT":
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			p.parsePDBConvertClause(opts, "SERVICE_NAME_CONVERT", optStart)

		// SOURCE_FILE_NAME_CONVERT = (...) | SOURCE_FILE_NAME_CONVERT = NONE
		case p.isIdentLike() && p.cur.Str == "SOURCE_FILE_NAME_CONVERT":
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			p.parsePDBConvertClause(opts, "SOURCE_FILE_NAME_CONVERT", optStart)

		// SOURCE_FILE_DIRECTORY = 'path' | SOURCE_FILE_DIRECTORY = NONE
		case p.isIdentLike() && p.cur.Str == "SOURCE_FILE_DIRECTORY":
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			val := ""
			if p.isIdentLike() && p.cur.Str == "NONE" {
				val = "NONE"
				p.advance()
			} else if p.cur.Type == tokSCONST {
				val = p.cur.Str
				p.advance()
			} else if p.isIdentLike() {
				val = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "SOURCE_FILE_DIRECTORY", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// STORAGE (...) | STORAGE UNLIMITED
		case p.isIdentLike() && p.cur.Str == "STORAGE":
			p.advance()
			if p.isIdentLike() && p.cur.Str == "UNLIMITED" {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "STORAGE", Value: "UNLIMITED",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.cur.Type == '(' {
				p.advance()
				storageItems := &nodes.List{}
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					sStart := p.pos()
					if p.isIdentLike() && p.cur.Str == "MAXSIZE" {
						p.advance()
						val := p.parsePDBSizeOrUnlimited()
						storageItems.Items = append(storageItems.Items, &nodes.DDLOption{
							Key: "MAXSIZE", Value: val,
							Loc: nodes.Loc{Start: sStart, End: p.pos()},
						})
					} else if p.isIdentLike() && p.cur.Str == "MAX_AUDIT_SIZE" {
						p.advance()
						val := p.parsePDBSizeOrUnlimited()
						storageItems.Items = append(storageItems.Items, &nodes.DDLOption{
							Key: "MAX_AUDIT_SIZE", Value: val,
							Loc: nodes.Loc{Start: sStart, End: p.pos()},
						})
					} else if p.isIdentLike() && p.cur.Str == "MAX_DIAG_SIZE" {
						p.advance()
						val := p.parsePDBSizeOrUnlimited()
						storageItems.Items = append(storageItems.Items, &nodes.DDLOption{
							Key: "MAX_DIAG_SIZE", Value: val,
							Loc: nodes.Loc{Start: sStart, End: p.pos()},
						})
					} else {
						p.advance()
					}
				}
				if p.cur.Type == ')' {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "STORAGE", Items: storageItems,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// PATH_PREFIX = 'path' | PATH_PREFIX = NONE
		case p.isIdentLike() && p.cur.Str == "PATH_PREFIX":
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			val := ""
			if p.isIdentLike() && p.cur.Str == "NONE" {
				val = "NONE"
				p.advance()
			} else if p.cur.Type == tokSCONST {
				val = p.cur.Str
				p.advance()
			} else if p.isIdentLike() {
				val = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "PATH_PREFIX", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// TEMPFILE REUSE
		case p.isIdentLike() && p.cur.Str == "TEMPFILE":
			p.advance()
			if p.isIdentLike() && p.cur.Str == "REUSE" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "TEMPFILE_REUSE",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// USER_TABLESPACES = ALL [EXCEPT (...)] | (list) | NONE
		case p.isIdentLike() && p.cur.Str == "USER_TABLESPACES":
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			p.parsePDBAllExceptNoneList(opts, "USER_TABLESPACES", optStart)

		// STANDBYS = ALL [EXCEPT (...)] | (list) | NONE
		case p.isIdentLike() && p.cur.Str == "STANDBYS":
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			p.parsePDBAllExceptNoneList(opts, "STANDBYS", optStart)

		// LOGGING / NOLOGGING
		case p.isIdentLike() && p.cur.Str == "LOGGING":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "LOGGING",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})
		case p.isIdentLike() && p.cur.Str == "NOLOGGING":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "NOLOGGING",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// CREATE_FILE_DEST = 'path' | CREATE_FILE_DEST = NONE
		case p.isIdentLike() && p.cur.Str == "CREATE_FILE_DEST":
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			val := ""
			if p.isIdentLike() && p.cur.Str == "NONE" {
				val = "NONE"
				p.advance()
			} else if p.cur.Type == tokSCONST {
				val = p.cur.Str
				p.advance()
			} else if p.isIdentLike() {
				val = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "CREATE_FILE_DEST", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// KEYSTORE IDENTIFIED BY password [DECRYPT USING 'secret']
		// KEYSTORE EXTERNAL STORE [DECRYPT USING 'secret']
		case p.isIdentLike() && p.cur.Str == "KEYSTORE":
			p.advance()
			if p.cur.Type == kwIDENTIFIED {
				p.advance() // IDENTIFIED
				if p.isIdentLike() && p.cur.Str == "BY" {
					p.advance() // BY
				}
				pass := ""
				if p.isIdentLike() || p.cur.Type == tokSCONST {
					pass = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "KEYSTORE_PASSWORD", Value: pass,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "EXTERNAL" {
				p.advance() // EXTERNAL
				if p.isIdentLike() && p.cur.Str == "STORE" {
					p.advance() // STORE
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "KEYSTORE_EXTERNAL",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}
			// Check for DECRYPT USING after KEYSTORE
			if p.isIdentLike() && p.cur.Str == "DECRYPT" {
				p.advance() // DECRYPT
				if p.isIdentLike() && p.cur.Str == "USING" {
					p.advance() // USING
				}
				secret := ""
				if p.cur.Type == tokSCONST || p.isIdentLike() {
					secret = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DECRYPT_USING", Value: secret,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// DECRYPT USING 'secret' (standalone, for XML variant)
		case p.isIdentLike() && p.cur.Str == "DECRYPT":
			p.advance() // DECRYPT
			if p.isIdentLike() && p.cur.Str == "USING" {
				p.advance() // USING
			}
			secret := ""
			if p.cur.Type == tokSCONST || p.isIdentLike() {
				secret = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "DECRYPT_USING", Value: secret,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// SNAPSHOT COPY [NO DATA] / SNAPSHOT = NONE | MANUAL | EVERY n HOURS/MINUTES
		case p.isIdentLike() && p.cur.Str == "SNAPSHOT":
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			p.parsePDBSnapshotClause(opts, optStart)

		// REFRESH MODE { MANUAL | EVERY n MINUTES/HOURS | NONE }
		case p.isIdentLike() && p.cur.Str == "REFRESH":
			p.advance() // REFRESH
			if p.isIdentLike() && p.cur.Str == "MODE" {
				p.advance() // MODE
				val := ""
				if p.isIdentLike() && p.cur.Str == "MANUAL" {
					val = "MANUAL"
					p.advance()
				} else if p.isIdentLike() && p.cur.Str == "NONE" {
					val = "NONE"
					p.advance()
				} else if p.isIdentLike() && p.cur.Str == "EVERY" {
					p.advance() // EVERY
					interval := ""
					if p.cur.Type == tokICONST {
						interval = p.cur.Str
						p.advance()
					}
					unit := ""
					if p.isIdentLike() && (p.cur.Str == "MINUTES" || p.cur.Str == "HOURS") {
						unit = p.cur.Str
						p.advance()
					}
					val = "EVERY " + interval + " " + unit
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "REFRESH_MODE", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "SWITCHOVER" {
				// REFRESH SWITCHOVER TO PRIMARY dblink
				p.advance() // SWITCHOVER
				if p.cur.Type == kwTO {
					p.advance() // TO
				}
				if p.isIdentLike() && p.cur.Str == "PRIMARY" {
					p.advance() // PRIMARY
				}
				dblink := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT {
					dblink = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "REFRESH_SWITCHOVER", Value: dblink,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// NO DATA
		case p.isIdentLike() && p.cur.Str == "NO":
			p.advance() // NO
			if p.isIdentLikeStr("DATA") {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "NO_DATA",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// RELOCATE [KEEP SOURCE] [REFRESH MODE ...] [AVAILABILITY ...]
		case p.isIdentLike() && p.cur.Str == "RELOCATE":
			p.advance()
			val := ""
			if p.isIdentLike() && p.cur.Str == "KEEP" {
				p.advance() // KEEP
				if p.isIdentLike() && p.cur.Str == "SOURCE" {
					p.advance() // SOURCE
				}
				val = "KEEP_SOURCE"
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "RELOCATE", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// COPY / MOVE / NOCOPY (for XML variant)
		case p.isIdentLike() && p.cur.Str == "COPY":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "COPY",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})
		case p.isIdentLike() && p.cur.Str == "MOVE":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "MOVE",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})
		case p.isIdentLike() && p.cur.Str == "NOCOPY":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "NOCOPY",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// AS CLONE (for XML variant)
		case p.cur.Type == kwAS:
			p.advance() // AS
			if p.isIdentLike() && p.cur.Str == "CLONE" {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "AS_CLONE",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "APPLICATION" {
				p.advance() // APPLICATION
				if p.isIdentLike() && p.cur.Str == "CONTAINER" {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "AS_APPLICATION_CONTAINER",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLikeStr("SEED") {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "AS_SEED",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else {
				p.advance() // skip unknown after AS
			}

		// HOST 'hostname'
		case p.isIdentLike() && p.cur.Str == "HOST":
			p.advance()
			val := ""
			if p.cur.Type == tokSCONST {
				val = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "HOST", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// PORT number
		case p.isIdentLike() && p.cur.Str == "PORT":
			p.advance()
			val := ""
			if p.cur.Type == tokICONST {
				val = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "PORT", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// CONTAINER_MAP UPDATE (...)
		case p.isIdentLike() && p.cur.Str == "CONTAINER_MAP":
			p.advance() // CONTAINER_MAP
			if p.cur.Type == kwUPDATE || p.isIdentLikeStr("UPDATE") {
				p.advance()
			}
			// Skip the parenthesized content
			if p.cur.Type == '(' {
				depth := 1
				p.advance()
				for depth > 0 && p.cur.Type != tokEOF {
					if p.cur.Type == '(' {
						depth++
					} else if p.cur.Type == ')' {
						depth--
					}
					if depth > 0 {
						p.advance()
					}
				}
				if p.cur.Type == ')' {
					p.advance()
				}
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "CONTAINER_MAP",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// AVAILABILITY { NORMAL | MAX }
		case p.isIdentLike() && p.cur.Str == "AVAILABILITY":
			p.advance()
			val := ""
			if p.isIdentLike() {
				val = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "AVAILABILITY", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		default:
			// Unknown token — stop parsing
			return
		}
	}
}

// parsePDBConvertClause parses FILE_NAME_CONVERT/SERVICE_NAME_CONVERT/SOURCE_FILE_NAME_CONVERT
// = (...) | = NONE
func (p *Parser) parsePDBConvertClause(opts *nodes.List, key string, optStart int) {
	if p.isIdentLike() && p.cur.Str == "NONE" {
		p.advance()
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: key, Value: "NONE",
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
	} else if p.cur.Type == '(' {
		p.advance()
		var pairs []string
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.cur.Type == tokSCONST {
				pairs = append(pairs, p.cur.Str)
			}
			p.advance()
		}
		if p.cur.Type == ')' {
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: key, Value: strings.Join(pairs, ","),
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
	}
}

// parsePDBAllExceptNoneList parses ALL [EXCEPT (...)] | (list) | NONE for
// USER_TABLESPACES, STANDBYS, INSTANCES, SERVICES clauses.
func (p *Parser) parsePDBAllExceptNoneList(opts *nodes.List, key string, optStart int) {
	if p.cur.Type == kwALL || (p.isIdentLike() && p.cur.Str == "ALL") {
		p.advance()
		if p.cur.Type == kwEXCEPT || (p.isIdentLike() && p.cur.Str == "EXCEPT") {
			p.advance()
			// Parse exception list
			exceptItems := &nodes.List{}
			if p.cur.Type == '(' {
				p.advance()
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT || p.cur.Type == tokSCONST {
						exceptItems.Items = append(exceptItems.Items, &nodes.DDLOption{Value: p.cur.Str})
						p.advance()
					}
					if p.cur.Type == ',' {
						p.advance()
					}
				}
				if p.cur.Type == ')' {
					p.advance()
				}
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: key, Value: "ALL_EXCEPT", Items: exceptItems,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})
		} else {
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: key, Value: "ALL",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})
		}
	} else if p.isIdentLike() && p.cur.Str == "NONE" {
		p.advance()
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: key, Value: "NONE",
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
	} else {
		// Parenthesized list
		p.parsePDBParenList(opts, key, optStart)
	}
}

// parsePDBParenList parses a parenthesized list of identifiers.
func (p *Parser) parsePDBParenList(opts *nodes.List, key string, optStart int) {
	if p.cur.Type == '(' {
		p.advance()
		items := &nodes.List{}
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT || p.cur.Type == tokSCONST {
				items.Items = append(items.Items, &nodes.DDLOption{Value: p.cur.Str})
				p.advance()
			}
			if p.cur.Type == ',' {
				p.advance()
			}
		}
		if p.cur.Type == ')' {
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: key, Items: items,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
	}
}

// parsePDBSizeOrUnlimited parses a size value or UNLIMITED keyword.
func (p *Parser) parsePDBSizeOrUnlimited() string {
	if p.isIdentLike() && p.cur.Str == "UNLIMITED" {
		p.advance()
		return "UNLIMITED"
	}
	return p.parseSizeValue()
}

// parsePDBSnapshotClause parses SNAPSHOT clause variants:
// COPY [NO DATA], NONE, MANUAL, EVERY n HOURS/MINUTES
func (p *Parser) parsePDBSnapshotClause(opts *nodes.List, optStart int) {
	if p.isIdentLike() && p.cur.Str == "COPY" {
		p.advance()
		val := "COPY"
		if p.isIdentLike() && p.cur.Str == "NO" {
			p.advance() // NO
			if p.isIdentLikeStr("DATA") {
				p.advance()
			}
			val = "COPY_NO_DATA"
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "SNAPSHOT", Value: val,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
	} else if p.isIdentLike() && p.cur.Str == "NONE" {
		p.advance()
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "SNAPSHOT", Value: "NONE",
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
	} else if p.isIdentLike() && p.cur.Str == "MANUAL" {
		p.advance()
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "SNAPSHOT", Value: "MANUAL",
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
	} else if p.isIdentLike() && p.cur.Str == "EVERY" {
		p.advance()
		interval := ""
		if p.cur.Type == tokICONST {
			interval = p.cur.Str
			p.advance()
		}
		unit := ""
		if p.isIdentLike() && (p.cur.Str == "MINUTES" || p.cur.Str == "HOURS") {
			unit = p.cur.Str
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "SNAPSHOT", Value: "EVERY " + interval + " " + unit,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
	}
}

// parseAlterPluggableDatabaseStmt parses an ALTER PLUGGABLE DATABASE statement.
// The ALTER keyword has been consumed. PLUGGABLE DATABASE has been consumed by caller.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-PLUGGABLE-DATABASE.html
func (p *Parser) parseAlterPluggableDatabaseStmt(start int) nodes.StmtNode {
	stmt := &nodes.AdminDDLStmt{
		Action:     "ALTER",
		ObjectType: nodes.OBJECT_PLUGGABLE_DATABASE,
		Loc:        nodes.Loc{Start: start},
	}

	opts := &nodes.List{}

	// Parse optional PDB name, ALL, or ALL EXCEPT (...)
	// Special cases: APPLICATION, CONTAINERS, DEFAULT, SET, OPEN, CLOSE, SAVE, DISCARD, etc.
	// can appear without a PDB name.
	if p.cur.Type == kwALL || (p.isIdentLike() && p.cur.Str == "ALL") {
		p.advance() // ALL
		if p.cur.Type == kwEXCEPT || (p.isIdentLike() && p.cur.Str == "EXCEPT") {
			// ALL EXCEPT (pdb1, pdb2, ...)
			p.advance()
			exceptItems := &nodes.List{}
			if p.cur.Type == '(' {
				p.advance()
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
						exceptItems.Items = append(exceptItems.Items, &nodes.DDLOption{Value: p.cur.Str})
						p.advance()
					}
					if p.cur.Type == ',' {
						p.advance()
					}
				}
				if p.cur.Type == ')' {
					p.advance()
				}
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "TARGET", Value: "ALL_EXCEPT", Items: exceptItems,
			})
		} else {
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "TARGET", Value: "ALL"})
		}
	} else if p.isAlterPDBClauseKeyword() {
		// No PDB name — clause keyword directly
	} else if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
		// PDB name
		stmt.Name = p.parseObjectName()
	}

	// Parse the main clause
	p.parseAlterPDBClauses(opts)

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.pos()
	return stmt
}

// isAlterPDBClauseKeyword checks if current token starts an ALTER PDB clause (no PDB name).
func (p *Parser) isAlterPDBClauseKeyword() bool {
	switch p.cur.Type {
	case kwOPEN, kwCLOSE, kwSET, kwDEFAULT, kwRENAME,
		kwENABLE, kwDISABLE, kwCREATE, kwDROP,
		kwGRANT, kwREVOKE:
		return true
	}
	if p.isIdentLike() {
		switch p.cur.Str {
		case "APPLICATION", "CONTAINERS", "UNPLUG", "RECOVER",
			"BACKUP", "SAVE", "DISCARD", "SNAPSHOT", "MATERIALIZE",
			"PREPARE", "DATAFILE", "STORAGE", "LOGGING", "NOLOGGING",
			"REFRESH", "PRIORITY":
			return true
		}
	}
	return false
}

// parseAlterPDBClauses parses the body of ALTER PLUGGABLE DATABASE.
func (p *Parser) parseAlterPDBClauses(opts *nodes.List) {
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		optStart := p.pos()

		switch {
		// OPEN [READ WRITE | READ ONLY | HYBRID READ ONLY | UPGRADE]
		// [RESTRICTED] [FORCE] [RESETLOGS] [instances_clause] [services_clause]
		case p.cur.Type == kwOPEN:
			p.advance() // OPEN
			mode := ""
			switch {
			case p.cur.Type == kwREAD:
				p.advance() // READ
				if p.isIdentLike() && p.cur.Str == "WRITE" {
					p.advance()
					mode = "READ_WRITE"
				} else if p.isIdentLike() && p.cur.Str == "ONLY" {
					p.advance()
					mode = "READ_ONLY"
				}
			case p.isIdentLike() && p.cur.Str == "HYBRID":
				p.advance() // HYBRID
				if p.cur.Type == kwREAD {
					p.advance() // READ
				}
				if p.isIdentLike() && p.cur.Str == "ONLY" {
					p.advance()
				}
				mode = "HYBRID_READ_ONLY"
			case p.isIdentLike() && p.cur.Str == "UPGRADE":
				p.advance()
				mode = "UPGRADE"
			}
			// Optional: RESTRICTED, FORCE, RESETLOGS
			flags := ""
			for {
				if p.isIdentLike() && p.cur.Str == "RESTRICTED" {
					p.advance()
					flags += ",RESTRICTED"
				} else if p.cur.Type == kwFORCE {
					p.advance()
					flags += ",FORCE"
				} else if p.isIdentLike() && p.cur.Str == "RESETLOGS" {
					p.advance()
					flags += ",RESETLOGS"
				} else {
					break
				}
			}
			if flags != "" {
				mode += flags
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "OPEN", Value: mode,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})
			// instances_clause / services_clause
			p.parsePDBInstancesServices(opts)

		// CLOSE [IMMEDIATE | ABORT] [instances_clause] [relocate_clause]
		case p.cur.Type == kwCLOSE:
			p.advance()
			mode := ""
			if p.cur.Type == kwIMMEDIATE {
				p.advance()
				mode = "IMMEDIATE"
			} else if p.isIdentLike() && p.cur.Str == "ABORT" {
				p.advance()
				mode = "ABORT"
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "CLOSE", Value: mode,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})
			// instances_clause
			p.parsePDBInstancesServices(opts)
			// relocate_clause
			if p.isIdentLike() && p.cur.Str == "RELOCATE" {
				relStart := p.pos()
				p.advance()
				val := ""
				if p.cur.Type == kwTO {
					p.advance()
					if p.cur.Type == tokSCONST || p.isIdentLike() {
						val = p.cur.Str
						p.advance()
					}
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "RELOCATE", Value: val,
					Loc: nodes.Loc{Start: relStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "NORELOCATE" {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "NORELOCATE",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// SAVE STATE / DISCARD STATE [instances_clause]
		case p.isIdentLike() && (p.cur.Str == "SAVE" || p.cur.Str == "DISCARD"):
			action := p.cur.Str
			p.advance()
			if p.isIdentLike() && p.cur.Str == "STATE" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: action + "_STATE",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})
			// instances_clause
			p.parsePDBInstancesServices(opts)

		// UNPLUG INTO 'filename' [ENCRYPT USING ...]
		case p.isIdentLike() && p.cur.Str == "UNPLUG":
			p.advance() // UNPLUG
			if p.cur.Type == kwINTO {
				p.advance() // INTO
			}
			filename := ""
			if p.cur.Type == tokSCONST {
				filename = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "UNPLUG", Value: filename,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})
			// ENCRYPT USING
			if p.isIdentLike() && p.cur.Str == "ENCRYPT" {
				p.advance()
				if p.isIdentLike() && p.cur.Str == "USING" {
					p.advance()
				}
				secret := ""
				if p.cur.Type == tokSCONST || p.isIdentLike() {
					secret = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "ENCRYPT_USING", Value: secret,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// SET DEFAULT TABLESPACE | SET TIME_ZONE | SET STANDBY NOLOGGING | SET MAX_PDB_SNAPSHOTS
		case p.cur.Type == kwSET:
			p.advance() // SET
			if p.cur.Type == kwDEFAULT {
				p.advance() // DEFAULT
				if p.cur.Type == kwTABLESPACE {
					p.advance() // TABLESPACE
					val := ""
					if p.isIdentLike() && (p.cur.Str == "BIGFILE" || p.cur.Str == "SMALLFILE") {
						val = p.cur.Str
						p.advance()
					}
					opts.Items = append(opts.Items, &nodes.DDLOption{
						Key: "SET_DEFAULT_TABLESPACE", Value: val,
						Loc: nodes.Loc{Start: optStart, End: p.pos()},
					})
				}
			} else if p.isIdentLike() && p.cur.Str == "TIME_ZONE" {
				p.advance() // TIME_ZONE
				if p.cur.Type == '=' {
					p.advance()
				}
				val := ""
				if p.cur.Type == tokSCONST {
					val = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "SET_TIME_ZONE", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "STANDBY" {
				p.advance() // STANDBY
				if p.isIdentLike() && p.cur.Str == "NOLOGGING" {
					p.advance() // NOLOGGING
				}
				// FOR LOAD PERFORMANCE | FOR DATA AVAILABILITY
				val := "STANDBY_NOLOGGING"
				if p.cur.Type == kwFOR {
					p.advance() // FOR
					if p.isIdentLike() && p.cur.Str == "LOAD" {
						p.advance()
						val += "_LOAD"
					} else if p.isIdentLikeStr("DATA") {
						p.advance()
						val += "_DATA"
					}
					if p.isIdentLike() && (p.cur.Str == "PERFORMANCE" || p.cur.Str == "AVAILABILITY") {
						p.advance()
					}
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "SET_STANDBY_NOLOGGING", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "MAX_PDB_SNAPSHOTS" {
				p.advance() // MAX_PDB_SNAPSHOTS
				if p.cur.Type == '=' {
					p.advance()
				}
				val := ""
				if p.cur.Type == tokICONST {
					val = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "SET_MAX_PDB_SNAPSHOTS", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else {
				// Unknown SET clause — skip token
				p.advance()
			}

		// DEFAULT EDITION | DEFAULT TABLESPACE | DEFAULT TEMPORARY TABLESPACE
		case p.cur.Type == kwDEFAULT:
			p.advance() // DEFAULT
			if p.isIdentLike() && p.cur.Str == "EDITION" {
				p.advance()
				val := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT {
					val = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DEFAULT_EDITION", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "TEMPORARY" {
				p.advance() // TEMPORARY
				if p.cur.Type == kwTABLESPACE {
					p.advance()
				}
				val := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
					val = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DEFAULT_TEMPORARY_TABLESPACE", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.cur.Type == kwTABLESPACE {
				p.advance()
				val := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
					val = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DEFAULT_TABLESPACE", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else {
				p.advance()
			}

		// RENAME GLOBAL_NAME TO ...
		case p.cur.Type == kwRENAME:
			p.advance() // RENAME
			if p.isIdentLike() && p.cur.Str == "GLOBAL_NAME" {
				p.advance()
			}
			if p.cur.Type == kwTO {
				p.advance()
			}
			// Parse global database name (may contain dots)
			var parts []string
			for {
				if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
					parts = append(parts, p.cur.Str)
					p.advance()
				}
				if p.cur.Type == '.' {
					parts = append(parts, ".")
					p.advance()
					continue
				}
				break
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "RENAME_GLOBAL_NAME", Value: strings.Join(parts, ""),
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// STORAGE (...)  | STORAGE UNLIMITED
		case p.isIdentLike() && p.cur.Str == "STORAGE":
			p.advance()
			if p.isIdentLike() && p.cur.Str == "UNLIMITED" {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "STORAGE", Value: "UNLIMITED",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.cur.Type == '(' {
				p.advance()
				storageItems := &nodes.List{}
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					sStart := p.pos()
					if p.isIdentLike() && p.cur.Str == "MAXSIZE" {
						p.advance()
						val := p.parsePDBSizeOrUnlimited()
						storageItems.Items = append(storageItems.Items, &nodes.DDLOption{
							Key: "MAXSIZE", Value: val,
							Loc: nodes.Loc{Start: sStart, End: p.pos()},
						})
					} else if p.isIdentLike() && p.cur.Str == "MAX_AUDIT_SIZE" {
						p.advance()
						val := p.parsePDBSizeOrUnlimited()
						storageItems.Items = append(storageItems.Items, &nodes.DDLOption{
							Key: "MAX_AUDIT_SIZE", Value: val,
							Loc: nodes.Loc{Start: sStart, End: p.pos()},
						})
					} else if p.isIdentLike() && p.cur.Str == "MAX_DIAG_SIZE" {
						p.advance()
						val := p.parsePDBSizeOrUnlimited()
						storageItems.Items = append(storageItems.Items, &nodes.DDLOption{
							Key: "MAX_DIAG_SIZE", Value: val,
							Loc: nodes.Loc{Start: sStart, End: p.pos()},
						})
					} else {
						p.advance()
					}
				}
				if p.cur.Type == ')' {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "STORAGE", Items: storageItems,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// LOGGING / NOLOGGING
		case p.isIdentLike() && p.cur.Str == "LOGGING":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "LOGGING",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})
		case p.isIdentLike() && p.cur.Str == "NOLOGGING":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "NOLOGGING",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// ENABLE FORCE LOGGING | ENABLE RECOVERY | ENABLE LOST WRITE PROTECTION | ENABLE BACKUP
		case p.cur.Type == kwENABLE:
			p.advance() // ENABLE
			if p.cur.Type == kwFORCE {
				p.advance() // FORCE
				val := "ENABLE_FORCE_LOGGING"
				if p.isIdentLike() && p.cur.Str == "NOLOGGING" {
					val = "ENABLE_FORCE_NOLOGGING"
					p.advance()
				} else if p.isIdentLike() && p.cur.Str == "LOGGING" {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "RECOVERY" {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "ENABLE_RECOVERY",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "LOST" {
				p.advance() // LOST
				if p.isIdentLike() && p.cur.Str == "WRITE" {
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "PROTECTION" {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "ENABLE_LOST_WRITE_PROTECTION",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "BACKUP" {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "ENABLE_BACKUP",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else {
				p.advance()
			}

		// DISABLE FORCE LOGGING | DISABLE RECOVERY | DISABLE LOST WRITE PROTECTION | DISABLE BACKUP
		case p.cur.Type == kwDISABLE:
			p.advance() // DISABLE
			if p.cur.Type == kwFORCE {
				p.advance() // FORCE
				val := "DISABLE_FORCE_LOGGING"
				if p.isIdentLike() && p.cur.Str == "NOLOGGING" {
					val = "DISABLE_FORCE_NOLOGGING"
					p.advance()
				} else if p.isIdentLike() && p.cur.Str == "LOGGING" {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "RECOVERY" {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DISABLE_RECOVERY",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "LOST" {
				p.advance() // LOST
				if p.isIdentLike() && p.cur.Str == "WRITE" {
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "PROTECTION" {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DISABLE_LOST_WRITE_PROTECTION",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "BACKUP" {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DISABLE_BACKUP",
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else {
				p.advance()
			}

		// REFRESH MODE ... | REFRESH SWITCHOVER ...
		case p.isIdentLike() && p.cur.Str == "REFRESH":
			p.advance() // REFRESH
			if p.isIdentLike() && p.cur.Str == "MODE" {
				p.advance()
				val := ""
				if p.isIdentLike() && p.cur.Str == "MANUAL" {
					val = "MANUAL"
					p.advance()
				} else if p.isIdentLike() && p.cur.Str == "NONE" {
					val = "NONE"
					p.advance()
				} else if p.isIdentLike() && p.cur.Str == "EVERY" {
					p.advance()
					interval := ""
					if p.cur.Type == tokICONST {
						interval = p.cur.Str
						p.advance()
					}
					unit := ""
					if p.isIdentLike() && (p.cur.Str == "MINUTES" || p.cur.Str == "HOURS") {
						unit = p.cur.Str
						p.advance()
					}
					val = "EVERY " + interval + " " + unit
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "REFRESH_MODE", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "SWITCHOVER" {
				p.advance() // SWITCHOVER
				if p.cur.Type == kwTO {
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "PRIMARY" {
					p.advance()
				}
				dblink := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT {
					dblink = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "REFRESH_SWITCHOVER", Value: dblink,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// CONTAINERS DEFAULT TARGET | CONTAINERS HOST | CONTAINERS PORT
		case p.isIdentLike() && p.cur.Str == "CONTAINERS":
			p.advance() // CONTAINERS
			if p.cur.Type == kwDEFAULT {
				p.advance() // DEFAULT
				if p.isIdentLikeStr("TARGET") {
					p.advance()
				}
				val := ""
				if p.isIdentLike() && p.cur.Str == "NONE" {
					val = "NONE"
					p.advance()
				} else if p.isIdentLike() || p.cur.Type == tokIDENT {
					val = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "CONTAINERS_DEFAULT_TARGET", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "HOST" {
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
				}
				val := ""
				if p.cur.Type == tokSCONST {
					val = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "CONTAINERS_HOST", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "PORT" {
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
				}
				val := ""
				if p.cur.Type == tokICONST {
					val = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "CONTAINERS_PORT", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// PRIORITY { integer | NONE }
		case p.isIdentLike() && p.cur.Str == "PRIORITY":
			p.advance()
			val := ""
			if p.isIdentLike() && p.cur.Str == "NONE" {
				val = "NONE"
				p.advance()
			} else if p.cur.Type == tokICONST {
				val = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "PRIORITY", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// DATAFILE clause
		case p.isIdentLike() && p.cur.Str == "DATAFILE":
			p.advance()
			val := ""
			if p.cur.Type == kwALL || (p.isIdentLike() && p.cur.Str == "ALL") {
				val = "ALL"
				p.advance()
			} else if p.cur.Type == tokSCONST {
				val = p.cur.Str
				p.advance()
			} else if p.cur.Type == tokICONST {
				val = p.cur.Str
				p.advance()
			}
			action := ""
			if p.isIdentLike() && p.cur.Str == "ONLINE" {
				action = "ONLINE"
				p.advance()
			} else if p.isIdentLike() && p.cur.Str == "OFFLINE" {
				action = "OFFLINE"
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "DATAFILE", Value: val + ":" + action,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// RECOVER [STANDBY] [UNTIL ...] [USING BACKUP CONTROLFILE]
		case p.isIdentLike() && p.cur.Str == "RECOVER":
			p.advance() // RECOVER
			val := ""
			if p.isIdentLike() && p.cur.Str == "STANDBY" {
				p.advance()
				val = "STANDBY"
			}
			if p.cur.Type == kwUNTIL || (p.isIdentLike() && p.cur.Str == "UNTIL") {
				p.advance() // UNTIL
				if p.isIdentLike() && p.cur.Str == "CANCEL" {
					p.advance()
					val += ",UNTIL_CANCEL"
				} else if p.isIdentLike() && p.cur.Str == "TIME" {
					p.advance()
					if p.cur.Type == tokSCONST {
						val += ",UNTIL_TIME:" + p.cur.Str
						p.advance()
					}
				} else if p.isIdentLike() && p.cur.Str == "CHANGE" {
					p.advance()
					if p.cur.Type == tokICONST {
						val += ",UNTIL_CHANGE:" + p.cur.Str
						p.advance()
					}
				} else if p.isIdentLike() && p.cur.Str == "SEQUENCE" {
					p.advance()
					if p.cur.Type == tokICONST {
						val += ",UNTIL_SEQUENCE:" + p.cur.Str
						p.advance()
					}
				}
			}
			if p.isIdentLike() && p.cur.Str == "USING" {
				p.advance() // USING
				if p.isIdentLike() && p.cur.Str == "BACKUP" {
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "CONTROLFILE" {
					p.advance()
				}
				val += ",BACKUP_CONTROLFILE"
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "RECOVER", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// BACKUP BEGIN | BACKUP END
		case p.isIdentLike() && p.cur.Str == "BACKUP":
			p.advance()
			val := ""
			if p.cur.Type == kwBEGIN || (p.isIdentLike() && p.cur.Str == "BEGIN") {
				val = "BEGIN"
				p.advance()
			} else if p.cur.Type == kwEND || (p.isIdentLike() && p.cur.Str == "END") {
				val = "END"
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "BACKUP", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// APPLICATION clauses
		case p.isIdentLike() && p.cur.Str == "APPLICATION":
			p.advance() // APPLICATION
			p.parseAlterPDBApplicationClause(opts, optStart)

		// SNAPSHOT clauses: SNAPSHOT NONE/MANUAL/EVERY
		case p.isIdentLike() && p.cur.Str == "SNAPSHOT":
			p.advance()
			p.parsePDBSnapshotClause(opts, optStart)

		// MATERIALIZE
		case p.isIdentLike() && p.cur.Str == "MATERIALIZE":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "MATERIALIZE",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		// CREATE SNAPSHOT snap_name
		case p.cur.Type == kwCREATE:
			p.advance() // CREATE
			if p.isIdentLike() && p.cur.Str == "SNAPSHOT" {
				p.advance()
				name := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT {
					name = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "CREATE_SNAPSHOT", Value: name,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// DROP SNAPSHOT snap_name | DROP MIRROR COPY mirror_name
		case p.cur.Type == kwDROP:
			p.advance() // DROP
			if p.isIdentLike() && p.cur.Str == "SNAPSHOT" {
				p.advance()
				name := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT {
					name = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DROP_SNAPSHOT", Value: name,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			} else if p.isIdentLike() && p.cur.Str == "MIRROR" {
				p.advance() // MIRROR
				if p.isIdentLike() && p.cur.Str == "COPY" {
					p.advance()
				}
				name := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT {
					name = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DROP_MIRROR_COPY", Value: name,
					Loc: nodes.Loc{Start: optStart, End: p.pos()},
				})
			}

		// PREPARE MIRROR COPY mirror_name [WITH ... REDUNDANCY] [FOR DATABASE ...]
		case p.isIdentLike() && p.cur.Str == "PREPARE":
			p.advance() // PREPARE
			if p.isIdentLike() && p.cur.Str == "MIRROR" {
				p.advance() // MIRROR
			}
			if p.isIdentLike() && p.cur.Str == "COPY" {
				p.advance() // COPY
			}
			name := ""
			if p.isIdentLike() || p.cur.Type == tokIDENT {
				name = p.cur.Str
				p.advance()
			}
			val := name
			if p.cur.Type == kwWITH {
				p.advance() // WITH
				redundancy := ""
				if p.isIdentLike() && (p.cur.Str == "EXTERNAL" || p.cur.Str == "NORMAL" || p.cur.Str == "HIGH") {
					redundancy = p.cur.Str
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "REDUNDANCY" {
					p.advance()
				}
				val += ":" + redundancy
			}
			if p.cur.Type == kwFOR {
				p.advance() // FOR
				if p.cur.Type == kwDATABASE {
					p.advance() // DATABASE
				}
				dbName := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT {
					dbName = p.cur.Str
					p.advance()
				}
				val += ":FOR:" + dbName
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "PREPARE_MIRROR_COPY", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})

		default:
			// Unknown clause — stop
			return
		}
	}
}

// parsePDBInstancesServices parses optional INSTANCES = (...) and SERVICES = (...) clauses.
func (p *Parser) parsePDBInstancesServices(opts *nodes.List) {
	for {
		optStart := p.pos()
		if p.isIdentLike() && p.cur.Str == "INSTANCES" {
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			p.parsePDBAllExceptNoneList(opts, "INSTANCES", optStart)
		} else if p.isIdentLike() && p.cur.Str == "SERVICES" {
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			p.parsePDBAllExceptNoneList(opts, "SERVICES", optStart)
		} else {
			return
		}
	}
}

// parseAlterPDBApplicationClause parses APPLICATION clauses in ALTER PLUGGABLE DATABASE.
func (p *Parser) parseAlterPDBApplicationClause(opts *nodes.List, optStart int) {
	// Check for ALL [EXCEPT (...)] SYNC or multi-app SYNC
	if p.cur.Type == kwALL || (p.isIdentLike() && p.cur.Str == "ALL") {
		p.advance() // ALL
		if p.cur.Type == kwEXCEPT || (p.isIdentLike() && p.cur.Str == "EXCEPT") {
			p.advance()
			exceptItems := &nodes.List{}
			if p.cur.Type == '(' {
				p.advance()
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					if p.isIdentLike() || p.cur.Type == tokIDENT {
						exceptItems.Items = append(exceptItems.Items, &nodes.DDLOption{Value: p.cur.Str})
						p.advance()
					}
					if p.cur.Type == ',' {
						p.advance()
					}
				}
				if p.cur.Type == ')' {
					p.advance()
				}
			}
			if p.isIdentLike() && p.cur.Str == "SYNC" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "APP_ALL_EXCEPT_SYNC", Items: exceptItems,
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})
		} else {
			// ALL SYNC
			if p.isIdentLike() && p.cur.Str == "SYNC" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "APP_ALL_SYNC",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})
		}
		return
	}

	// Parse app name(s) — could be single or comma-separated
	var appNames []string
	if p.isIdentLike() || p.cur.Type == tokIDENT {
		appNames = append(appNames, p.cur.Str)
		p.advance()
	}
	// Check for comma-separated list (multi-app SYNC)
	for p.cur.Type == ',' {
		p.advance()
		if p.isIdentLike() || p.cur.Type == tokIDENT {
			appNames = append(appNames, p.cur.Str)
			p.advance()
		}
	}

	// If we have multiple apps, it must be SYNC
	if len(appNames) > 1 {
		if p.isIdentLike() && p.cur.Str == "SYNC" {
			p.advance()
		}
		items := &nodes.List{}
		for _, n := range appNames {
			items.Items = append(items.Items, &nodes.DDLOption{Value: n})
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "APP_MULTI_SYNC", Items: items,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
		return
	}

	appName := ""
	if len(appNames) > 0 {
		appName = appNames[0]
	}

	// Single app operations
	switch {
	case p.cur.Type == kwBEGIN || (p.isIdentLike() && p.cur.Str == "BEGIN"):
		p.advance() // BEGIN
		p.parseAppBeginClause(opts, appName, optStart)

	case p.cur.Type == kwEND || (p.isIdentLike() && p.cur.Str == "END"):
		p.advance() // END
		p.parseAppEndClause(opts, appName, optStart)

	case p.cur.Type == kwSET:
		p.advance() // SET
		p.parseAppSetClause(opts, appName, optStart)

	case p.isIdentLike() && p.cur.Str == "SYNC":
		p.advance() // SYNC
		val := ""
		if p.cur.Type == kwTO {
			p.advance() // TO
			if p.isIdentLike() && p.cur.Str == "PATCH" {
				p.advance()
				if p.cur.Type == tokICONST {
					val = "PATCH:" + p.cur.Str
					p.advance()
				}
			} else if p.cur.Type == tokSCONST || p.isIdentLike() {
				val = p.cur.Str
				p.advance()
			}
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "APP_SYNC", Value: appName + ":" + val,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
	}
}

// parseAppBeginClause parses APPLICATION app_name BEGIN { INSTALL | PATCH | UPGRADE | UNINSTALL }.
func (p *Parser) parseAppBeginClause(opts *nodes.List, appName string, optStart int) {
	switch {
	case p.isIdentLike() && p.cur.Str == "INSTALL":
		p.advance() // INSTALL
		version := ""
		if p.isIdentLike() && p.cur.Str == "VERSION" {
			p.advance()
			if p.cur.Type == tokSCONST || p.isIdentLike() {
				version = p.cur.Str
				p.advance()
			}
		}
		comment := ""
		if p.cur.Type == kwCOMMENT {
			p.advance()
			if p.cur.Type == tokSCONST {
				comment = p.cur.Str
				p.advance()
			}
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "APP_BEGIN_INSTALL", Value: appName + ":" + version + ":" + comment,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})

	case p.isIdentLike() && p.cur.Str == "PATCH":
		p.advance() // PATCH
		patchNum := ""
		if p.cur.Type == tokICONST {
			patchNum = p.cur.Str
			p.advance()
		}
		minVersion := ""
		if p.isIdentLike() && p.cur.Str == "MINIMUM" {
			p.advance() // MINIMUM
			if p.isIdentLike() && p.cur.Str == "VERSION" {
				p.advance()
			}
			if p.cur.Type == tokSCONST || p.isIdentLike() {
				minVersion = p.cur.Str
				p.advance()
			}
		}
		comment := ""
		if p.cur.Type == kwCOMMENT {
			p.advance()
			if p.cur.Type == tokSCONST {
				comment = p.cur.Str
				p.advance()
			}
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "APP_BEGIN_PATCH", Value: appName + ":" + patchNum + ":" + minVersion + ":" + comment,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})

	case p.isIdentLike() && p.cur.Str == "UPGRADE":
		p.advance() // UPGRADE
		fromVersion := ""
		toVersion := ""
		if p.cur.Type == kwFROM {
			p.advance()
			if p.cur.Type == tokSCONST || p.isIdentLike() {
				fromVersion = p.cur.Str
				p.advance()
			}
		}
		if p.cur.Type == kwTO {
			p.advance()
			if p.cur.Type == tokSCONST || p.isIdentLike() {
				toVersion = p.cur.Str
				p.advance()
			}
		}
		comment := ""
		if p.cur.Type == kwCOMMENT {
			p.advance()
			if p.cur.Type == tokSCONST {
				comment = p.cur.Str
				p.advance()
			}
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "APP_BEGIN_UPGRADE", Value: appName + ":" + fromVersion + ":" + toVersion + ":" + comment,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})

	case p.isIdentLike() && p.cur.Str == "UNINSTALL":
		p.advance()
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "APP_BEGIN_UNINSTALL", Value: appName,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
	}
}

// parseAppEndClause parses APPLICATION app_name END { INSTALL | PATCH | UPGRADE | UNINSTALL }.
func (p *Parser) parseAppEndClause(opts *nodes.List, appName string, optStart int) {
	switch {
	case p.isIdentLike() && p.cur.Str == "INSTALL":
		p.advance()
		version := ""
		if p.isIdentLike() && p.cur.Str == "VERSION" {
			p.advance()
			if p.cur.Type == tokSCONST || p.isIdentLike() {
				version = p.cur.Str
				p.advance()
			}
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "APP_END_INSTALL", Value: appName + ":" + version,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})

	case p.isIdentLike() && p.cur.Str == "PATCH":
		p.advance()
		patchNum := ""
		if p.cur.Type == tokICONST {
			patchNum = p.cur.Str
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "APP_END_PATCH", Value: appName + ":" + patchNum,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})

	case p.isIdentLike() && p.cur.Str == "UPGRADE":
		p.advance()
		toVersion := ""
		if p.cur.Type == kwTO {
			p.advance()
			if p.cur.Type == tokSCONST || p.isIdentLike() {
				toVersion = p.cur.Str
				p.advance()
			}
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "APP_END_UPGRADE", Value: appName + ":" + toVersion,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})

	case p.isIdentLike() && p.cur.Str == "UNINSTALL":
		p.advance()
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "APP_END_UNINSTALL", Value: appName,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
	}
}

// parseAppSetClause parses APPLICATION app_name SET { PATCH | VERSION | COMPATIBILITY }.
func (p *Parser) parseAppSetClause(opts *nodes.List, appName string, optStart int) {
	switch {
	case p.isIdentLike() && p.cur.Str == "PATCH":
		p.advance()
		val := ""
		if p.cur.Type == tokICONST {
			val = p.cur.Str
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "APP_SET_PATCH", Value: appName + ":" + val,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})

	case p.isIdentLike() && p.cur.Str == "VERSION":
		p.advance()
		val := ""
		if p.cur.Type == tokSCONST || p.isIdentLike() {
			val = p.cur.Str
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "APP_SET_VERSION", Value: appName + ":" + val,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})

	case p.isIdentLike() && p.cur.Str == "COMPATIBILITY":
		p.advance() // COMPATIBILITY
		if p.isIdentLike() && p.cur.Str == "VERSION" {
			p.advance()
		}
		val := ""
		if p.isIdentLike() && p.cur.Str == "CURRENT" {
			val = "CURRENT"
			p.advance()
		} else if p.cur.Type == tokSCONST || p.isIdentLike() {
			val = p.cur.Str
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "APP_SET_COMPATIBILITY", Value: appName + ":" + val,
			Loc: nodes.Loc{Start: optStart, End: p.pos()},
		})
	}
}

// parseDropPluggableDatabaseStmt parses a DROP PLUGGABLE DATABASE statement.
// The DROP keyword has been consumed. PLUGGABLE DATABASE has been consumed by caller.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-PLUGGABLE-DATABASE.html
//
//	DROP PLUGGABLE DATABASE [ IF EXISTS ] pdb_name [ FORCE ]
//	  { KEEP DATAFILES | INCLUDING DATAFILES }
func (p *Parser) parseDropPluggableDatabaseStmt(start int) nodes.StmtNode {
	stmt := &nodes.AdminDDLStmt{
		Action:     "DROP",
		ObjectType: nodes.OBJECT_PLUGGABLE_DATABASE,
		Loc:        nodes.Loc{Start: start},
	}

	// Optional IF EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwEXISTS {
			p.advance() // IF
			p.advance() // EXISTS
			stmt.IfExists = true
		}
	}

	// PDB name
	stmt.Name = p.parseObjectName()

	opts := &nodes.List{}

	// Optional trailing: FORCE, KEEP DATAFILES, INCLUDING DATAFILES
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		optStart := p.pos()
		switch {
		case p.cur.Type == kwFORCE:
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "FORCE",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})
		case p.isIdentLike() && p.cur.Str == "INCLUDING":
			p.advance() // INCLUDING
			if p.isIdentLike() && p.cur.Str == "DATAFILES" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "INCLUDING_DATAFILES",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})
		case p.isIdentLike() && p.cur.Str == "KEEP":
			p.advance() // KEEP
			if p.isIdentLike() && p.cur.Str == "DATAFILES" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "KEEP_DATAFILES",
				Loc: nodes.Loc{Start: optStart, End: p.pos()},
			})
		default:
			goto done
		}
	}
done:
	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.pos()
	return stmt
}
