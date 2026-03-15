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
