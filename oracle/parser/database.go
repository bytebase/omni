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
func (p *Parser) parseCreateDatabaseStmt(start int) (nodes.StmtNode, error) {
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
		var parseErr614 error
		stmt.Name, parseErr614 = p.parseObjectName()
		if parseErr614 !=

			// Parse CREATE DATABASE clauses
			nil {
			return nil, parseErr614
		}
	}

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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// CONTROLFILE REUSE
		case p.isIdentLike() && p.cur.Str == "CONTROLFILE":
			p.advance() // CONTROLFILE
			if p.isIdentLike() && p.cur.Str == "REUSE" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "CONTROLFILE_REUSE",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					df, parseErr615 := p.parseDatafileClause()
					if parseErr615 != nil {
						return nil, parseErr615
					}
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
					Loc:   nodes.Loc{Start: lfStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// ARCHIVELOG
		case p.isIdentLike() && p.cur.Str == "ARCHIVELOG":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "ARCHIVELOG",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// NOARCHIVELOG
		case p.isIdentLike() && p.cur.Str == "NOARCHIVELOG":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "NOARCHIVELOG",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// FORCE LOGGING
		case p.cur.Type == kwFORCE:
			p.advance() // FORCE
			if p.cur.Type == kwLOGGING {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "FORCE_LOGGING",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// SYSAUX DATAFILE file_specification [, ...]
		case p.isIdentLike() && p.cur.Str == "SYSAUX":
			p.advance() // SYSAUX
			if p.isIdentLike() && (p.cur.Str == "DATAFILE" || p.cur.Str == "DATAFILES") {
				p.advance()
			}
			fileList := &nodes.List{}
			for p.cur.Type == tokSCONST {
				df, parseErr616 := p.parseDatafileClause()
				if parseErr616 != nil {
					return nil, parseErr616
				}
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					parseValue65, parseErr66 := p.parseSizeValue()
					if parseErr66 != nil {
						return nil, parseErr66
					}
					val += " " + parseValue65
				}
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "EXTENT_MANAGEMENT_LOCAL", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// DATAFILE 'file' SIZE ... (standalone datafile spec, not part of tablespace)
		case p.isIdentLike() && p.cur.Str == "DATAFILE":
			p.advance() // DATAFILE
			fileList := &nodes.List{}
			for p.cur.Type == tokSCONST {
				df, parseErr617 := p.parseDatafileClause()
				if parseErr617 != nil {
					return nil, parseErr617
				}
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
						df, parseErr618 := p.parseDatafileClause()
						if parseErr618 != nil {
							return nil, parseErr618
						}
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
					Loc:   nodes.Loc{Start: optStart, End: p.prev.End},
				})
			} else {
				// Just BIGFILE/SMALLFILE without DEFAULT — skip
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "FILE_TYPE", Value: fileType,
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
						df, parseErr619 := p.parseDatafileClause()
						if parseErr619 != nil {
							return nil, parseErr619
						}
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
					Loc:   nodes.Loc{Start: optStart, End: p.prev.End},
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
						df, parseErr620 := p.parseDatafileClause()
						if parseErr620 != nil {
							return nil, parseErr620
						}
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
					Loc:   nodes.Loc{Start: optStart, End: p.prev.End},
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
						df, parseErr621 := p.parseDatafileClause()
						if parseErr621 != nil {
							return nil, parseErr621
						}
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
					Loc:   nodes.Loc{Start: optStart, End: p.prev.End},
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
					df, parseErr622 := p.parseDatafileClause()
					if parseErr622 != nil {
						return nil, parseErr622
					}
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
				Loc:   nodes.Loc{Start: optStart, End: p.prev.End},
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
								Loc: nodes.Loc{Start: seedStart, End: p.prev.End},
							})
						} else if p.cur.Type == '(' {
							p.advance()
							var pairs []string
							for progress := -1; p.cur.Type != ')' && p.cur.Type != tokEOF; {
								// An item that consumes nothing would come round again.
								if p.pos() == progress {
									return nil, p.syntaxErrorAtCur()
								}
								progress = p.pos()
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
								Loc: nodes.Loc{Start: seedStart, End: p.prev.End},
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
							var parseErr623 error
							size, parseErr623 = p.parseSizeValue()
							if parseErr623 != nil {
								return nil, parseErr623
							}
						}
						ae, parseErr624 := p.parseOptionalAutoextend()
						if parseErr624 != nil {
							return nil, parseErr624
						}
						items := &nodes.List{}
						if ae != nil {
							items.Items = append(items.Items, ae)
						}
						seedItems.Items = append(seedItems.Items, &nodes.DDLOption{
							Key: "SYSTEM_DATAFILES", Value: size,
							Items: items,
							Loc:   nodes.Loc{Start: seedStart, End: p.prev.End},
						})
					} else if p.isIdentLike() && p.cur.Str == "SYSAUX" {
						p.advance() // SYSAUX
						if p.isIdentLike() && (p.cur.Str == "DATAFILES" || p.cur.Str == "DATAFILE") {
							p.advance()
						}
						size := ""
						if p.cur.Type == kwSIZE {
							p.advance()
							var parseErr625 error
							size, parseErr625 = p.parseSizeValue()
							if parseErr625 != nil {
								return nil, parseErr625
							}
						}
						ae, parseErr626 := p.parseOptionalAutoextend()
						if parseErr626 != nil {
							return nil, parseErr626
						}
						items := &nodes.List{}
						if ae != nil {
							items.Items = append(items.Items, ae)
						}
						seedItems.Items = append(seedItems.Items, &nodes.DDLOption{
							Key: "SYSAUX_DATAFILES", Value: size,
							Items: items,
							Loc:   nodes.Loc{Start: seedStart, End: p.prev.End},
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
							Loc: nodes.Loc{Start: seedStart, End: p.prev.End},
						})
					} else {
						break
					}
				}
			}
			opt := &nodes.DDLOption{
				Key: "ENABLE_PLUGGABLE_DATABASE",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			}
			if len(seedItems.Items) > 0 {
				opt.Items = seedItems
			}
			opts.Items = append(opts.Items, opt)

		default:
			// No other clause (ORA-02000 at USING).
			return nil, p.syntaxErrorAtCur()
		}
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
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
func (p *Parser) parseCreateControlfileStmt(start int) (nodes.StmtNode, error) {
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
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					df, parseErr627 := p.parseDatafileClause()
					if parseErr627 != nil {
						return nil, parseErr627
					}
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
					Loc:   nodes.Loc{Start: lfStart, End: p.prev.End},
				})
				if p.cur.Type == ',' {
					p.advance()
					continue
				}
				break
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "LOGFILE", Items: logItems,
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// RESETLOGS
		case p.isIdentLike() && p.cur.Str == "RESETLOGS":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "RESETLOGS",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// NORESETLOGS
		case p.isIdentLike() && p.cur.Str == "NORESETLOGS":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "NORESETLOGS",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// DATAFILE 'file' [, ...]
		case p.isIdentLike() && p.cur.Str == "DATAFILE":
			p.advance()
			fileList := &nodes.List{}
			for p.cur.Type == tokSCONST {
				df, parseErr628 := p.parseDatafileClause()
				if parseErr628 != nil {
					return nil, parseErr628
				}
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// ARCHIVELOG
		case p.isIdentLike() && p.cur.Str == "ARCHIVELOG":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "ARCHIVELOG",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// NOARCHIVELOG
		case p.isIdentLike() && p.cur.Str == "NOARCHIVELOG":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "NOARCHIVELOG",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// FORCE LOGGING
		case p.cur.Type == kwFORCE:
			p.advance()
			if p.cur.Type == kwLOGGING {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "FORCE_LOGGING",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		default:
			p.advance()
		}
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
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
func (p *Parser) parseAlterDatabaseDictionaryStmt(start int) (nodes.StmtNode, error) {
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
		if err := p.expectWords("CREDENTIALS"); err != nil {
			return nil, err
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "ENCRYPT_CREDENTIALS",
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	// REKEY CREDENTIALS
	case p.isIdentLike() && p.cur.Str == "REKEY":
		p.advance()
		if err := p.expectWords("CREDENTIALS"); err != nil {
			return nil, err
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "REKEY_CREDENTIALS",
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	// DELETE CREDENTIALS KEY
	case p.cur.Type == kwDELETE:
		p.advance()
		if err := p.expectWords("CREDENTIALS"); err != nil {
			return nil, err
		}
		if p.cur.Type != kwKEY {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "DELETE_CREDENTIALS_KEY",
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
	default:
		return nil, p.syntaxErrorAtCur()
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseOptionalAutoextend parses an optional AUTOEXTEND clause if present, returning nil otherwise.
func (p *Parser) parseOptionalAutoextend() (*nodes.AutoextendClause, error) {
	if p.isIdentLike() && p.cur.Str == "AUTOEXTEND" {
		return p.parseAutoextendClause()
	}
	return nil, nil
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

// parseDropPluggableDatabaseStmt parses a DROP PLUGGABLE DATABASE statement.
// The DROP keyword has been consumed. PLUGGABLE DATABASE has been consumed by caller.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-PLUGGABLE-DATABASE.html
//
//	DROP PLUGGABLE DATABASE [ IF EXISTS ] pdb_name [ FORCE ]
//	  { KEEP DATAFILES | INCLUDING DATAFILES }
func (p *Parser) parseDropPluggableDatabaseStmt(start int) (nodes.StmtNode, error) {
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
	var parseErr669 error

	// PDB name
	stmt.Name, parseErr669 = p.parseReservedCheckedObjectName()
	if parseErr669 != nil {
		return nil,

			// Optional trailing: FORCE, KEEP DATAFILES, INCLUDING DATAFILES
			parseErr669
	}

	opts := &nodes.List{}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		optStart := p.pos()
		switch {
		case p.cur.Type == kwFORCE:
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "FORCE",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})
		case p.isIdentLike() && p.cur.Str == "INCLUDING":
			p.advance() // INCLUDING
			if err := p.expectWords("DATAFILES"); err != nil {
				return nil, err
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "INCLUDING_DATAFILES",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})
		case p.isIdentLike() && p.cur.Str == "KEEP":
			p.advance() // KEEP
			if err := p.expectWords("DATAFILES"); err != nil {
				return nil, err
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "KEEP_DATAFILES",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})
		default:
			goto done
		}
	}
done:
	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCreateDiskgroupStmt parses a CREATE DISKGROUP statement.
// The CREATE keyword has been consumed. DISKGROUP has been consumed by caller.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-DISKGROUP.html
func (p *Parser) parseCreateDiskgroupStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_DISKGROUP,
		Loc:        nodes.Loc{Start: start},
	}
	var parseErr670 error

	// Diskgroup name
	stmt.Name, parseErr670 = p.parseObjectName()
	if parseErr670 != nil {
		return nil,

			// Optional redundancy: { NORMAL | HIGH | FLEX | EXTENDED | EXTERNAL } REDUNDANCY
			parseErr670
	}

	opts := &nodes.List{}

	if p.isIdentLike() {
		switch p.cur.Str {
		case "NORMAL", "HIGH", "FLEX", "EXTENDED", "EXTERNAL":
			red := p.cur.Str
			p.advance()
			if err := p.expectWords("REDUNDANCY"); err != nil {
				return nil, err
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "REDUNDANCY", Value: red})
		}
	}

	// Parse FAILGROUP/DISK clauses; one DISK at least (ORA-02000).
	hasDisk := false
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		optStart := p.pos()

		switch {
		// { QUORUM | REGULAR } FAILGROUP failgroup_name DISK ...
		case p.isIdentLike() && (p.cur.Str == "QUORUM" || p.cur.Str == "REGULAR" || p.cur.Str == "FAILGROUP"):
			qualifier := ""
			if p.cur.Str == "QUORUM" || p.cur.Str == "REGULAR" {
				qualifier = p.cur.Str
				p.advance()
			}
			if p.isIdentLike() && p.cur.Str == "FAILGROUP" {
				p.advance() // FAILGROUP
				fgName, err := p.parseRequiredIdentifier()
				if err != nil {
					return nil, err
				}
				val := fgName
				if qualifier != "" {
					val = qualifier + ":" + fgName
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "FAILGROUP", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			}
			// The disks of the failure group follow.
			if !p.isKeywordStr("DISK") {
				return nil, p.syntaxErrorAtCur()
			}

		// DISK qualified_disk_clause [, qualified_disk_clause]...
		case p.isIdentLike() && p.cur.Str == "DISK":
			p.advance() // DISK
			hasDisk = true
			for {
				dStart := p.pos()
				// search_string
				if p.cur.Type != tokSCONST {
					return nil, p.syntaxErrorAtCur()
				}
				path := p.cur.Str
				p.advance()
				name := ""
				size := ""
				forceFlag := ""
				if p.isIdentLike() && p.cur.Str == "NAME" {
					p.advance()
					var err error
					if name, err = p.parseRequiredIdentifier(); err != nil {
						return nil, err
					}
				}
				if p.cur.Type == kwSIZE {
					p.advance()
					var parseErr671 error
					size, parseErr671 = p.parseSizeValue()
					if parseErr671 != nil {
						return nil, parseErr671
					}
					if size == "" {
						return nil, p.syntaxErrorAtCur()
					}
				}
				if p.cur.Type == kwFORCE {
					forceFlag = "FORCE"
					p.advance()
				} else if p.isIdentLike() && p.cur.Str == "NOFORCE" {
					forceFlag = "NOFORCE"
					p.advance()
				}
				val := path
				if name != "" {
					val += ":NAME:" + name
				}
				if size != "" {
					val += ":SIZE:" + size
				}
				if forceFlag != "" {
					val += ":" + forceFlag
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DISK", Value: val,
					Loc: nodes.Loc{Start: dStart, End: p.prev.End},
				})
				if p.cur.Type == ',' {
					p.advance()
					continue
				}
				break
			}

		// ATTRIBUTE 'name' = 'value' [, ...]
		case p.isIdentLike() && p.cur.Str == "ATTRIBUTE":
			p.advance()
			attrItems := &nodes.List{}
			for {
				// 'attribute_name' = 'attribute_value'
				if p.cur.Type != tokSCONST {
					return nil, p.syntaxErrorAtCur()
				}
				attrName := p.cur.Str
				p.advance()
				if p.cur.Type != '=' {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance()
				if p.cur.Type != tokSCONST {
					return nil, p.syntaxErrorAtCur()
				}
				attrVal := p.cur.Str
				p.advance()
				attrItems.Items = append(attrItems.Items, &nodes.DDLOption{
					Key: attrName, Value: attrVal,
				})
				if p.cur.Type == ',' {
					p.advance()
					continue
				}
				break
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "ATTRIBUTE", Items: attrItems,
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		default:
			goto createDgDone
		}
	}
createDgDone:
	if !hasDisk {
		return nil, p.syntaxErrorAtCur()
	}
	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropDiskgroupStmt parses a DROP DISKGROUP statement.
// The DROP keyword has been consumed. DISKGROUP has been consumed by caller.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-DISKGROUP.html
func (p *Parser) parseDropDiskgroupStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "DROP",
		ObjectType: nodes.OBJECT_DISKGROUP,
		Loc:        nodes.Loc{Start: start},
	}
	var parseErr691 error

	if p.cur.Type == ';' || p.cur.Type == tokEOF {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Name, parseErr691 = p.parseObjectName()
	if parseErr691 != nil {
		return nil, parseErr691
	}

	opts := &nodes.List{}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		optStart := p.pos()
		switch {
		case p.isIdentLike() && p.cur.Str == "INCLUDING":
			p.advance()
			if err := p.expectWords("CONTENTS"); err != nil {
				return nil, err
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "INCLUDING_CONTENTS",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})
			goto dropDgDone
		case p.isIdentLike() && p.cur.Str == "EXCLUDING" && len(opts.Items) == 0:
			p.advance()
			if err := p.expectWords("CONTENTS"); err != nil {
				return nil, err
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "EXCLUDING_CONTENTS",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})
			goto dropDgDone
		case p.cur.Type == kwFORCE && len(opts.Items) == 0:
			// FORCE INCLUDING CONTENTS
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "FORCE",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})
			if !p.isKeywordStr("INCLUDING") {
				return nil, p.syntaxErrorAtCur()
			}
		default:
			goto dropDgDone
		}
	}
dropDgDone:
	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}
