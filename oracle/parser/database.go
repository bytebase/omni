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
			// Unknown token, advance to avoid infinite loop
			p.advance()
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
		if p.isIdentLike() && p.cur.Str == "CREDENTIALS" {
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "ENCRYPT_CREDENTIALS",
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	// REKEY CREDENTIALS
	case p.isIdentLike() && p.cur.Str == "REKEY":
		p.advance()
		if p.isIdentLike() && p.cur.Str == "CREDENTIALS" {
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "REKEY_CREDENTIALS",
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
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

// parsePDBSizeOrUnlimited parses a size value or UNLIMITED keyword.
func (p *Parser) parsePDBSizeOrUnlimited() (string, error) {
	if p.isIdentLike() && p.cur.Str == "UNLIMITED" {
		p.advance()
		return "UNLIMITED", nil
	}
	return p.parseSizeValue()
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
	stmt.Name, parseErr669 = p.parseObjectName()
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
			if p.isIdentLike() && p.cur.Str == "DATAFILES" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "INCLUDING_DATAFILES",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})
		case p.isIdentLike() && p.cur.Str == "KEEP":
			p.advance() // KEEP
			if p.isIdentLike() && p.cur.Str == "DATAFILES" {
				p.advance()
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
			if p.isIdentLike() && p.cur.Str == "REDUNDANCY" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "REDUNDANCY", Value: red})
		}
	}

	// Parse FAILGROUP/DISK clauses
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
				fgName := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT {
					fgName = p.cur.Str
					p.advance()
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

		// DISK qualified_disk_clause [, qualified_disk_clause]...
		case p.isIdentLike() && p.cur.Str == "DISK":
			p.advance() // DISK
			for {
				dStart := p.pos()
				// path_name (string constant or identifier)
				path := ""
				if p.cur.Type == tokSCONST {
					path = p.cur.Str
					p.advance()
				} else if p.isIdentLike() || p.cur.Type == tokIDENT {
					path = p.cur.Str
					p.advance()
				}
				name := ""
				size := ""
				forceFlag := ""
				if p.isIdentLike() && p.cur.Str == "NAME" {
					p.advance()
					if p.isIdentLike() || p.cur.Type == tokIDENT {
						name = p.cur.Str
						p.advance()
					}
				}
				if p.cur.Type == kwSIZE {
					p.advance()
					var parseErr671 error
					size, parseErr671 = p.parseSizeValue()
					if parseErr671 != nil {
						return nil, parseErr671
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
				if p.cur.Type == tokSCONST {
					attrName := p.cur.Str
					p.advance()
					if p.cur.Type == '=' {
						p.advance()
					}
					attrVal := ""
					if p.cur.Type == tokSCONST {
						attrVal = p.cur.Str
						p.advance()
					}
					attrItems.Items = append(attrItems.Items, &nodes.DDLOption{
						Key: attrName, Value: attrVal,
					})
				}
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
	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterDiskgroupStmt parses an ALTER DISKGROUP statement.
// The ALTER keyword has been consumed. DISKGROUP has been consumed by caller.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-DISKGROUP.html
func (p *Parser) parseAlterDiskgroupStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "ALTER",
		ObjectType: nodes.OBJECT_DISKGROUP,
		Loc:        nodes.Loc{Start: start},
	}
	var parseErr672 error

	// Diskgroup name
	stmt.Name, parseErr672 = p.parseObjectName()
	if parseErr672 != nil {
		return nil,

			// Parse the main clause
			parseErr672
	}

	opts := &nodes.List{}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		optStart := p.pos()

		switch {
		// ADD DISK / ADD DIRECTORY / ADD ALIAS / ADD TEMPLATE / ADD VOLUME /
		// ADD USERGROUP / ADD USER / ADD QUOTAGROUP / ADD FILEGROUP
		case p.cur.Type == kwADD || (p.isIdentLike() && p.cur.Str == "ADD"):
			p.advance()
			parseErr673 := // ADD
				p.parseDiskgroupAddClause(opts, optStart)
			if parseErr673 !=

				// DROP DISK / DROP DISKS / DROP TEMPLATE / DROP DIRECTORY / DROP ALIAS /
				// DROP VOLUME / DROP USERGROUP / DROP USER / DROP QUOTAGROUP / DROP FILEGROUP / DROP FILE
				nil {
				return nil, parseErr673
			}

		case p.cur.Type == kwDROP:
			p.advance()
			parseErr674 := // DROP
				p.parseDiskgroupDropClause(opts, optStart)
			if parseErr674 !=

				// RESIZE ALL SIZE / RESIZE VOLUME / RESIZE DISK
				nil {
				return nil, parseErr674
			}

		case p.isIdentLike() && p.cur.Str == "RESIZE":
			p.advance() // RESIZE
			if p.cur.Type == kwALL || (p.isIdentLike() && p.cur.Str == "ALL") {
				p.advance() // ALL
				if p.cur.Type == kwSIZE {
					p.advance()
				}
				size, parseErr675 := p.parseSizeValue()
				if parseErr675 != nil {
					return nil, parseErr675
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "RESIZE_ALL", Value: size,
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			} else if p.isIdentLike() && p.cur.Str == "VOLUME" {
				p.advance() // VOLUME
				volName := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT {
					volName = p.cur.Str
					p.advance()
				}
				if p.cur.Type == kwSIZE {
					p.advance()
				}
				size, parseErr676 := p.parseSizeValue()
				if parseErr676 != nil {
					return nil, parseErr676
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "RESIZE_VOLUME", Value: volName + ":" + size,
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			} else if p.isIdentLike() && p.cur.Str == "DISK" {
				p.advance() // DISK
				diskName := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT {
					diskName = p.cur.Str
					p.advance()
				}
				if p.cur.Type == kwSIZE {
					p.advance()
				}
				size, parseErr677 := p.parseSizeValue()
				if parseErr677 != nil {
					return nil, parseErr677
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "RESIZE_DISK", Value: diskName + ":" + size,
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			}

		// REPLACE DISK / REPLACE USER
		case p.isIdentLike() && p.cur.Str == "REPLACE":
			p.advance() // REPLACE
			if p.isIdentLike() && p.cur.Str == "DISK" {
				p.advance() // DISK
				oldName := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT {
					oldName = p.cur.Str
					p.advance()
				}
				if p.cur.Type == kwWITH {
					p.advance() // WITH
				}
				newPath := ""
				if p.cur.Type == tokSCONST {
					newPath = p.cur.Str
					p.advance()
				} else if p.isIdentLike() || p.cur.Type == tokIDENT {
					newPath = p.cur.Str
					p.advance()
				}
				val := oldName + ":" + newPath
				forceWait, parseErr678 := p.parseDiskgroupForceWaitPower()
				if parseErr678 != nil {
					return nil, parseErr678
				}
				val += forceWait
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "REPLACE_DISK", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			} else if p.cur.Type == kwUSER || (p.isIdentLike() && p.cur.Str == "USER") {
				p.advance() // USER
				oldUser := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT {
					oldUser = p.cur.Str
					p.advance()
				}
				if p.cur.Type == kwWITH {
					p.advance()
				}
				newUser := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT {
					newUser = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "REPLACE_USER", Value: oldUser + ":" + newUser,
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			}

		// RENAME DISK / RENAME DISKS ALL / RENAME DIRECTORY / RENAME ALIAS
		case p.cur.Type == kwRENAME:
			p.advance() // RENAME
			if p.isIdentLike() && p.cur.Str == "DISK" {
				p.advance()
				var pairs []string
				for {
					oldName := ""
					if p.isIdentLike() || p.cur.Type == tokIDENT {
						oldName = p.cur.Str
						p.advance()
					}
					if p.cur.Type == kwTO {
						p.advance()
					}
					newName := ""
					if p.isIdentLike() || p.cur.Type == tokIDENT {
						newName = p.cur.Str
						p.advance()
					}
					pairs = append(pairs, oldName+":"+newName)
					if p.cur.Type == ',' {
						p.advance()
						continue
					}
					break
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "RENAME_DISK", Value: strings.Join(pairs, ","),
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			} else if p.isIdentLike() && p.cur.Str == "DISKS" {
				p.advance() // DISKS
				if p.cur.Type == kwALL || (p.isIdentLike() && p.cur.Str == "ALL") {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "RENAME_DISKS_ALL",
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			} else if p.isIdentLike() && p.cur.Str == "DIRECTORY" {
				p.advance()
				oldPath := ""
				if p.cur.Type == tokSCONST {
					oldPath = p.cur.Str
					p.advance()
				}
				if p.cur.Type == kwTO {
					p.advance()
				}
				newPath := ""
				if p.cur.Type == tokSCONST {
					newPath = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "RENAME_DIRECTORY", Value: oldPath + ":" + newPath,
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			} else if p.isIdentLike() && p.cur.Str == "ALIAS" {
				p.advance()
				oldAlias := ""
				if p.cur.Type == tokSCONST {
					oldAlias = p.cur.Str
					p.advance()
				}
				if p.cur.Type == kwTO {
					p.advance()
				}
				newAlias := ""
				if p.cur.Type == tokSCONST {
					newAlias = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "RENAME_ALIAS", Value: oldAlias + ":" + newAlias,
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			}

		// ONLINE DISK / ONLINE DISKS IN FAILGROUP / ONLINE ALL
		case p.isIdentLike() && p.cur.Str == "ONLINE":
			p.advance()
			val := ""
			if p.isIdentLike() && p.cur.Str == "DISK" {
				p.advance()
				parseValue67, parseErr68 := p.parseDiskNameList()
				if parseErr68 != nil {
					return nil, parseErr68
				}
				val = "DISK:" + parseValue67
			} else if p.isIdentLike() && p.cur.Str == "DISKS" {
				p.advance()
				if p.cur.Type == kwIN {
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "FAILGROUP" {
					p.advance()
				}
				fgName := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT {
					fgName = p.cur.Str
					p.advance()
				}
				val = "FAILGROUP:" + fgName
			} else if p.cur.Type == kwALL || (p.isIdentLike() && p.cur.Str == "ALL") {
				p.advance()
				val = "ALL"
			}
			forceWait, parseErr679 := p.parseDiskgroupForceWaitPower()
			if parseErr679 != nil {
				return nil, parseErr679
			}
			val += forceWait
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "ONLINE", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// OFFLINE DISK / OFFLINE DISKS IN FAILGROUP
		case p.isIdentLike() && p.cur.Str == "OFFLINE":
			p.advance()
			val := ""
			if p.isIdentLike() && p.cur.Str == "DISK" {
				p.advance()
				parseValue69, parseErr70 := p.parseDiskNameList()
				if parseErr70 != nil {
					return nil, parseErr70
				}
				val = "DISK:" + parseValue69
			} else if p.isIdentLike() && p.cur.Str == "DISKS" {
				p.advance()
				if p.cur.Type == kwIN {
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "FAILGROUP" {
					p.advance()
				}
				fgName := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT {
					fgName = p.cur.Str
					p.advance()
				}
				val = "FAILGROUP:" + fgName
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "OFFLINE", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// REBALANCE
		case p.isIdentLike() && p.cur.Str == "REBALANCE":
			p.advance()
			val := ""
			if p.cur.Type == kwWITH {
				p.advance()
				var modes []string
				for p.isIdentLike() && (p.cur.Str == "RESTORE" || p.cur.Str == "BALANCE" || p.cur.Str == "PREPARE" || p.cur.Str == "COMPACT") {
					modes = append(modes, p.cur.Str)
					p.advance()
					if p.cur.Type == ',' {
						p.advance()
					}
				}
				val = "WITH:" + strings.Join(modes, ",")
			} else if p.isIdentLike() && p.cur.Str == "WITHOUT" {
				p.advance()
				var modes []string
				for p.isIdentLike() && (p.cur.Str == "BALANCE" || p.cur.Str == "PREPARE" || p.cur.Str == "COMPACT") {
					modes = append(modes, p.cur.Str)
					p.advance()
					if p.cur.Type == ',' {
						p.advance()
					}
				}
				val = "WITHOUT:" + strings.Join(modes, ",")
			}
			forceWait, parseErr680 := p.parseDiskgroupForceWaitPower()
			if parseErr680 != nil {
				return nil, parseErr680
			}
			val += forceWait
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "REBALANCE", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// CHECK
		case p.isIdentLike() && p.cur.Str == "CHECK":
			p.advance()
			val := ""
			if p.cur.Type == kwALL || (p.isIdentLike() && p.cur.Str == "ALL") {
				p.advance()
				val = "ALL"
			} else if p.isIdentLike() && p.cur.Str == "DISK" {
				p.advance()
				parseValue71, parseErr72 := p.parseDiskNameList()
				if parseErr72 != nil {
					return nil, parseErr72
				}
				val = "DISK:" + parseValue71
			} else if p.isIdentLike() && p.cur.Str == "DISKS" {
				p.advance()
				if p.cur.Type == kwIN {
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "FAILGROUP" {
					p.advance()
				}
				fgName := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT {
					fgName = p.cur.Str
					p.advance()
				}
				val = "FAILGROUP:" + fgName
			} else if p.isIdentLike() && p.cur.Str == "FILE" {
				p.advance()
				val = "FILE"
				for p.cur.Type == tokSCONST || (p.isIdentLike() && p.cur.Str != "REPAIR" && p.cur.Str != "NOREPAIR") {
					p.advance()
					if p.cur.Type == ',' {
						p.advance()
					}
				}
			}
			if p.isIdentLike() && p.cur.Str == "REPAIR" {
				p.advance()
				val += ":REPAIR"
			} else if p.isIdentLike() && p.cur.Str == "NOREPAIR" {
				p.advance()
				val += ":NOREPAIR"
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "CHECK", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// MODIFY TEMPLATE / MODIFY VOLUME / MODIFY USERGROUP / MODIFY FILEGROUP /
		// MODIFY FILE / MODIFY POWER / MODIFY QUOTAGROUP
		case p.isIdentLike() && p.cur.Str == "MODIFY":
			p.advance()
			parseErr681 := p.parseDiskgroupModifyClause(opts, optStart)
			if parseErr681 !=

				// SET ATTRIBUTE / SET PERMISSION / SET OWNER / SET GROUP
				nil {
				return nil, parseErr681
			}

		case p.cur.Type == kwSET:
			p.advance()
			parseErr682 := p.parseDiskgroupSetClause(opts, optStart)
			if parseErr682 !=

				// CONVERT REDUNDANCY
				nil {
				return nil, parseErr682
			}

		case p.isIdentLike() && p.cur.Str == "CONVERT":
			p.advance()
			if p.isIdentLike() && p.cur.Str == "REDUNDANCY" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "CONVERT_REDUNDANCY",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// UNDROP DISKS
		case p.isIdentLike() && p.cur.Str == "UNDROP":
			p.advance()
			if p.isIdentLike() && p.cur.Str == "DISKS" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "UNDROP_DISKS",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// MOUNT
		case p.isIdentLike() && p.cur.Str == "MOUNT":
			p.advance()
			val := ""
			if p.cur.Type == kwALL || (p.isIdentLike() && p.cur.Str == "ALL") {
				val = "ALL"
				p.advance()
			}
			if p.isIdentLike() && (p.cur.Str == "RESTRICTED" || p.cur.Str == "NORMAL") {
				if val != "" {
					val += ","
				}
				val += p.cur.Str
				p.advance()
			}
			if p.cur.Type == kwFORCE {
				if val != "" {
					val += ","
				}
				val += "FORCE"
				p.advance()
			} else if p.isIdentLike() && p.cur.Str == "NOFORCE" {
				if val != "" {
					val += ","
				}
				val += "NOFORCE"
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "MOUNT", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// DISMOUNT
		case p.isIdentLike() && p.cur.Str == "DISMOUNT":
			p.advance()
			val := ""
			if p.cur.Type == kwALL || (p.isIdentLike() && p.cur.Str == "ALL") {
				val = "ALL"
				p.advance()
			}
			if p.cur.Type == kwFORCE {
				if val != "" {
					val += ","
				}
				val += "FORCE"
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "DISMOUNT", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// ENABLE VOLUME / DISABLE VOLUME
		case p.cur.Type == kwENABLE:
			p.advance()
			if p.isIdentLike() && p.cur.Str == "VOLUME" {
				p.advance()
				val, parseErr683 := p.parseDiskgroupVolumeList()
				if parseErr683 != nil {
					return nil, parseErr683
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "ENABLE_VOLUME", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			}
		case p.cur.Type == kwDISABLE:
			p.advance()
			if p.isIdentLike() && p.cur.Str == "VOLUME" {
				p.advance()
				val, parseErr684 := p.parseDiskgroupVolumeList()
				if parseErr684 != nil {
					return nil, parseErr684
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DISABLE_VOLUME", Value: val,
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			}

		// SCRUB
		case p.isIdentLike() && p.cur.Str == "SCRUB":
			p.advance()
			val := ""
			if p.isIdentLike() && p.cur.Str == "STOP" {
				p.advance()
				val = "STOP"
			} else {
				if p.isIdentLike() && p.cur.Str == "FILE" {
					p.advance()
					val = "FILE"
					if p.cur.Type == tokSCONST || p.isIdentLike() {
						val += ":" + p.cur.Str
						p.advance()
					}
				} else if p.isIdentLike() && p.cur.Str == "DISK" {
					p.advance()
					val = "DISK"
					if p.isIdentLike() || p.cur.Type == tokIDENT {
						val += ":" + p.cur.Str
						p.advance()
					}
				}
				if p.isIdentLike() && p.cur.Str == "REPAIR" {
					p.advance()
					val += ",REPAIR"
				} else if p.isIdentLike() && p.cur.Str == "NOREPAIR" {
					p.advance()
					val += ",NOREPAIR"
				}
				for {
					if p.isIdentLike() && p.cur.Str == "POWER" {
						p.advance()
						if p.isIdentLike() || p.cur.Type == tokICONST {
							val += ",POWER:" + p.cur.Str
							p.advance()
						}
					} else if p.isIdentLike() && p.cur.Str == "WAIT" {
						val += ",WAIT"
						p.advance()
					} else if p.isIdentLike() && p.cur.Str == "NOWAIT" {
						val += ",NOWAIT"
						p.advance()
					} else if p.cur.Type == kwFORCE {
						val += ",FORCE"
						p.advance()
					} else if p.isIdentLike() && p.cur.Str == "NOFORCE" {
						val += ",NOFORCE"
						p.advance()
					} else {
						break
					}
				}
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "SCRUB", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// MOVE FILE ... TO FILEGROUP / MOVE FILEGROUP ... TO QUOTAGROUP
		case p.isIdentLike() && p.cur.Str == "MOVE":
			p.advance()
			if p.isIdentLike() && p.cur.Str == "FILE" {
				p.advance()
				fileSpec := ""
				if p.cur.Type == tokSCONST {
					fileSpec = p.cur.Str
					p.advance()
				} else if p.isIdentLike() || p.cur.Type == tokIDENT {
					fileSpec = p.cur.Str
					p.advance()
				}
				if p.cur.Type == kwTO {
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "FILEGROUP" {
					p.advance()
				}
				fgName := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT {
					fgName = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "MOVE_FILE_TO_FILEGROUP", Value: fileSpec + ":" + fgName,
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			} else if p.isIdentLike() && p.cur.Str == "FILEGROUP" {
				p.advance()
				fgName := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT {
					fgName = p.cur.Str
					p.advance()
				}
				if p.cur.Type == kwTO {
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "QUOTAGROUP" {
					p.advance()
				}
				qgName := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT {
					qgName = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "MOVE_FILEGROUP_TO_QUOTAGROUP", Value: fgName + ":" + qgName,
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			}

		// ALTER TEMPLATE
		case p.cur.Type == kwALTER:
			p.advance()
			if p.isIdentLike() && p.cur.Str == "TEMPLATE" {
				p.advance()
				tmplName := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT {
					tmplName = p.cur.Str
					p.advance()
				}
				parseErr685 := p.parseDiskgroupSkipParens()
				if parseErr685 != nil {
					return nil, parseErr685
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "ALTER_TEMPLATE", Value: tmplName,
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			}

		default:
			goto alterDgDone
		}
	}
alterDgDone:
	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDiskgroupAddClause handles ADD sub-clauses for ALTER DISKGROUP.
func (p *Parser) parseDiskgroupAddClause(opts *nodes.List, optStart int) error {
	switch {
	case p.isIdentLike() && p.cur.Str == "DISK":
		p.advance()
		for {
			path := ""
			if p.cur.Type == tokSCONST {
				path = p.cur.Str
				p.advance()
			} else if p.isIdentLike() || p.cur.Type == tokIDENT {
				path = p.cur.Str
				p.advance()
			}
			val := path
			if p.isIdentLike() && p.cur.Str == "NAME" {
				p.advance()
				if p.isIdentLike() || p.cur.Type == tokIDENT {
					val += ":NAME:" + p.cur.Str
					p.advance()
				}
			}
			if p.cur.Type == kwSIZE {
				p.advance()
				parseValue73, parseErr74 := p.parseSizeValue()
				if parseErr74 != nil {
					return parseErr74
				}
				val += ":SIZE:" + parseValue73
			}
			if p.isIdentLike() && p.cur.Str == "FAILGROUP" {
				p.advance()
				if p.isIdentLike() || p.cur.Type == tokIDENT {
					val += ":FAILGROUP:" + p.cur.Str
					p.advance()
				}
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "ADD_DISK", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})
			if p.cur.Type == ',' {
				p.advance()
				continue
			}
			break
		}

	case p.isIdentLike() && p.cur.Str == "DIRECTORY":
		p.advance()
		path := ""
		if p.cur.Type == tokSCONST {
			path = p.cur.Str
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "ADD_DIRECTORY", Value: path,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	case p.isIdentLike() && p.cur.Str == "ALIAS":
		p.advance()
		alias := ""
		if p.cur.Type == tokSCONST {
			alias = p.cur.Str
			p.advance()
		}
		if p.cur.Type == kwFOR {
			p.advance()
		}
		target := ""
		if p.cur.Type == tokSCONST {
			target = p.cur.Str
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "ADD_ALIAS", Value: alias + ":" + target,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	case p.isIdentLike() && p.cur.Str == "TEMPLATE":
		p.advance()
		tmplName := ""
		if p.isIdentLike() || p.cur.Type == tokIDENT {
			tmplName = p.cur.Str
			p.advance()
		}
		parseErr686 := p.parseDiskgroupSkipParens()
		if parseErr686 != nil {
			return parseErr686
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "ADD_TEMPLATE", Value: tmplName,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	case p.isIdentLike() && p.cur.Str == "VOLUME":
		p.advance()
		volName := ""
		if p.isIdentLike() || p.cur.Type == tokIDENT {
			volName = p.cur.Str
			p.advance()
		}
		size := ""
		if p.cur.Type == kwSIZE {
			p.advance()
			var parseErr687 error
			size, parseErr687 = p.parseSizeValue()
			if parseErr687 != nil {
				return parseErr687
			}
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "ADD_VOLUME", Value: volName + ":" + size,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
		for p.isIdentLike() && (p.cur.Str == "MIRROR" || p.cur.Str == "HIGH" || p.cur.Str == "UNPROTECTED" ||
			p.cur.Str == "PARITY" || p.cur.Str == "DOUBLE" || p.cur.Str == "FINE" || p.cur.Str == "COARSE" ||
			p.cur.Str == "STRIPE_WIDTH" || p.cur.Str == "STRIPE_COLUMNS") {
			p.advance()
			if p.cur.Type == tokICONST || p.cur.Type == tokSCONST {
				p.advance()
			}
			if p.cur.Type == kwSIZE {
				p.advance()
				parseDiscard689, parseErr688 := p.parseSizeValue()
				_ = parseDiscard689
				if parseErr688 != nil {
					return parseErr688
				}
			}
		}

	case p.isIdentLike() && p.cur.Str == "USERGROUP":
		p.advance()
		name := ""
		if p.isIdentLike() || p.cur.Type == tokIDENT {
			name = p.cur.Str
			p.advance()
		}
		member := ""
		if p.isIdentLike() && p.cur.Str == "MEMBER" {
			p.advance()
			if p.isIdentLike() || p.cur.Type == tokIDENT {
				member = p.cur.Str
				p.advance()
			}
		}
		val := name
		if member != "" {
			val += ":" + member
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "ADD_USERGROUP", Value: val,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	case p.cur.Type == kwUSER || (p.isIdentLike() && p.cur.Str == "USER"):
		p.advance()
		var users []string
		for {
			if p.isIdentLike() || p.cur.Type == tokIDENT {
				users = append(users, p.cur.Str)
				p.advance()
			}
			if p.cur.Type == ',' {
				p.advance()
				continue
			}
			break
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "ADD_USER", Value: strings.Join(users, ","),
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	case p.isIdentLike() && p.cur.Str == "QUOTAGROUP":
		p.advance()
		name := ""
		if p.isIdentLike() || p.cur.Type == tokIDENT {
			name = p.cur.Str
			p.advance()
		}
		val := name
		if p.cur.Type == kwSET {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "QUOTA" {
				p.advance()
			}
			if p.cur.Type == '=' {
				p.advance()
			}
			parseValue75, parseErr76 := p.parsePDBSizeOrUnlimited()
			if parseErr76 != nil {
				return parseErr76
			}
			val += ":" + parseValue75
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "ADD_QUOTAGROUP", Value: val,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	case p.isIdentLike() && p.cur.Str == "FILEGROUP":
		p.advance()
		name := ""
		if p.isIdentLike() || p.cur.Type == tokIDENT {
			name = p.cur.Str
			p.advance()
		}
		for p.cur.Type != ';' && p.cur.Type != tokEOF {
			if p.cur.Type == kwDATABASE || p.cur.Type == kwCLUSTER ||
				(p.isIdentLike() && (p.cur.Str == "VOLUME" || p.cur.Str == "TEMPLATE")) {
				p.advance()
				if p.isIdentLike() || p.cur.Type == tokIDENT {
					p.advance()
				}
			} else if p.cur.Type == kwSET {
				p.advance()
				for p.cur.Type != ';' && p.cur.Type != tokEOF {
					p.advance()
					if p.cur.Type == '=' {
						p.advance()
						if p.cur.Type == tokSCONST || p.cur.Type == tokICONST || p.isIdentLike() {
							p.advance()
						}
					}
					if p.cur.Type != ',' {
						break
					}
					p.advance()
				}
			} else {
				break
			}
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "ADD_FILEGROUP", Value: name,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
	}
	return nil
}

// parseDiskgroupDropClause handles DROP sub-clauses for ALTER DISKGROUP.
func (p *Parser) parseDiskgroupDropClause(opts *nodes.List, optStart int) error {
	switch {
	case p.isIdentLike() && p.cur.Str == "DISK":
		p.advance()
		parseValue77, parseErr78 := p.parseDiskNameList()
		if parseErr78 != nil {
			return parseErr78
		}
		val := "DISK:" + parseValue77
		if p.cur.Type == kwFORCE {
			val += ":FORCE"
			p.advance()
		} else if p.isIdentLike() && p.cur.Str == "NOFORCE" {
			val += ":NOFORCE"
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "DROP_DISK", Value: val,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	case p.isIdentLike() && p.cur.Str == "DISKS":
		p.advance()
		if p.cur.Type == kwIN {
			p.advance()
		}
		if p.isIdentLike() && p.cur.Str == "FAILGROUP" {
			p.advance()
		}
		fgName := ""
		if p.isIdentLike() || p.cur.Type == tokIDENT {
			fgName = p.cur.Str
			p.advance()
		}
		val := fgName
		if p.cur.Type == kwFORCE {
			val += ":FORCE"
			p.advance()
		} else if p.isIdentLike() && p.cur.Str == "NOFORCE" {
			val += ":NOFORCE"
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "DROP_DISKS_FAILGROUP", Value: val,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	case p.isIdentLike() && p.cur.Str == "TEMPLATE":
		p.advance()
		name := ""
		if p.isIdentLike() || p.cur.Type == tokIDENT {
			name = p.cur.Str
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "DROP_TEMPLATE", Value: name,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	case p.isIdentLike() && p.cur.Str == "DIRECTORY":
		p.advance()
		path := ""
		if p.cur.Type == tokSCONST {
			path = p.cur.Str
			p.advance()
		}
		force := ""
		if p.cur.Type == kwFORCE {
			force = ":FORCE"
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "DROP_DIRECTORY", Value: path + force,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	case p.isIdentLike() && p.cur.Str == "ALIAS":
		p.advance()
		alias := ""
		if p.cur.Type == tokSCONST {
			alias = p.cur.Str
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "DROP_ALIAS", Value: alias,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	case p.isIdentLike() && p.cur.Str == "VOLUME":
		p.advance()
		name := ""
		if p.isIdentLike() || p.cur.Type == tokIDENT {
			name = p.cur.Str
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "DROP_VOLUME", Value: name,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	case p.isIdentLike() && p.cur.Str == "USERGROUP":
		p.advance()
		name := ""
		if p.isIdentLike() || p.cur.Type == tokIDENT {
			name = p.cur.Str
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "DROP_USERGROUP", Value: name,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	case p.cur.Type == kwUSER || (p.isIdentLike() && p.cur.Str == "USER"):
		p.advance()
		var users []string
		for {
			if p.isIdentLike() || p.cur.Type == tokIDENT {
				users = append(users, p.cur.Str)
				p.advance()
			}
			if p.cur.Type == kwCASCADE || (p.isIdentLike() && p.cur.Str == "CASCADE") {
				p.advance()
			}
			if p.cur.Type == ',' {
				p.advance()
				continue
			}
			break
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "DROP_USER", Value: strings.Join(users, ","),
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	case p.isIdentLike() && p.cur.Str == "QUOTAGROUP":
		p.advance()
		name := ""
		if p.isIdentLike() || p.cur.Type == tokIDENT {
			name = p.cur.Str
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "DROP_QUOTAGROUP", Value: name,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	case p.isIdentLike() && p.cur.Str == "FILEGROUP":
		p.advance()
		name := ""
		if p.isIdentLike() || p.cur.Type == tokIDENT {
			name = p.cur.Str
			p.advance()
		}
		val := name
		if p.cur.Type == kwCASCADE || (p.isIdentLike() && p.cur.Str == "CASCADE") {
			val += ":CASCADE"
			p.advance()
			for p.cur.Type == kwFOR {
				p.advance()
				if p.cur.Type == kwDATABASE || (p.isIdentLike() && p.cur.Str == "PLUGGABLE") {
					p.advance()
					if p.cur.Type == kwDATABASE {
						p.advance()
					}
					if p.isIdentLike() || p.cur.Type == tokIDENT {
						p.advance()
					}
				}
			}
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "DROP_FILEGROUP", Value: val,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	case p.isIdentLike() && p.cur.Str == "FILE":
		p.advance()
		fileSpec := ""
		if p.cur.Type == tokSCONST {
			fileSpec = p.cur.Str
			p.advance()
		} else if p.isIdentLike() || p.cur.Type == tokIDENT {
			fileSpec = p.cur.Str
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "DROP_FILE", Value: fileSpec,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
	}
	return nil
}

// parseDiskgroupModifyClause handles MODIFY sub-clauses for ALTER DISKGROUP.
func (p *Parser) parseDiskgroupModifyClause(opts *nodes.List, optStart int) error {
	switch {
	case p.isIdentLike() && p.cur.Str == "TEMPLATE":
		p.advance()
		name := ""
		if p.isIdentLike() || p.cur.Type == tokIDENT {
			name = p.cur.Str
			p.advance()
		}
		parseErr690 := p.parseDiskgroupSkipParens()
		if parseErr690 != nil {
			return parseErr690
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "MODIFY_TEMPLATE", Value: name,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	case p.isIdentLike() && p.cur.Str == "VOLUME":
		p.advance()
		name := ""
		if p.isIdentLike() || p.cur.Type == tokIDENT {
			name = p.cur.Str
			p.advance()
		}
		for p.isIdentLike() && (p.cur.Str == "MOUNTPATH" || p.cur.Str == "USAGE") {
			p.advance()
			if p.cur.Type == tokSCONST {
				p.advance()
			}
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "MODIFY_VOLUME", Value: name,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	case p.isIdentLike() && p.cur.Str == "USERGROUP":
		p.advance()
		name := ""
		if p.isIdentLike() || p.cur.Type == tokIDENT {
			name = p.cur.Str
			p.advance()
		}
		action := ""
		if p.cur.Type == kwADD || (p.isIdentLike() && p.cur.Str == "ADD") {
			action = "ADD"
			p.advance()
		} else if p.cur.Type == kwDROP {
			action = "DROP"
			p.advance()
		}
		if p.isIdentLike() && p.cur.Str == "MEMBER" {
			p.advance()
		}
		member := ""
		if p.isIdentLike() || p.cur.Type == tokIDENT {
			member = p.cur.Str
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "MODIFY_USERGROUP", Value: name + ":" + action + ":" + member,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	case p.isIdentLike() && p.cur.Str == "FILEGROUP":
		p.advance()
		name := ""
		if p.isIdentLike() || p.cur.Type == tokIDENT {
			name = p.cur.Str
			p.advance()
		}
		if p.cur.Type == kwSET {
			p.advance()
			for p.cur.Type != ';' && p.cur.Type != tokEOF {
				if p.cur.Type == tokSCONST || p.isIdentLike() {
					p.advance()
				}
				if p.cur.Type == '=' {
					p.advance()
					if p.cur.Type == tokSCONST || p.cur.Type == tokICONST || p.isIdentLike() {
						p.advance()
					}
				}
				if p.cur.Type == ',' {
					p.advance()
					continue
				}
				break
			}
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "MODIFY_FILEGROUP", Value: name,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	case p.isIdentLike() && p.cur.Str == "FILE":
		p.advance()
		fileSpec := ""
		if p.cur.Type == tokSCONST {
			fileSpec = p.cur.Str
			p.advance()
		}
		for p.isIdentLike() && (p.cur.Str == "MIRROR" || p.cur.Str == "HIGH" || p.cur.Str == "UNPROTECTED" || p.cur.Str == "PARITY" || p.cur.Str == "DOUBLE") {
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "MODIFY_FILE", Value: fileSpec,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	case p.isIdentLike() && p.cur.Str == "POWER":
		p.advance()
		val := ""
		if p.cur.Type == tokICONST {
			val = p.cur.Str
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "MODIFY_POWER", Value: val,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	case p.isIdentLike() && p.cur.Str == "QUOTAGROUP":
		p.advance()
		name := ""
		if p.isIdentLike() || p.cur.Type == tokIDENT {
			name = p.cur.Str
			p.advance()
		}
		if p.cur.Type == kwSET {
			p.advance()
		}
		if p.isIdentLike() && p.cur.Str == "QUOTA" {
			p.advance()
		}
		if p.cur.Type == '=' {
			p.advance()
		}
		parseValue79, parseErr80 := p.parsePDBSizeOrUnlimited()
		if parseErr80 != nil {
			return parseErr80
		}
		val := name + ":" + parseValue79
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "MODIFY_QUOTAGROUP", Value: val,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
	}
	return nil
}

// parseDiskgroupSetClause handles SET sub-clauses for ALTER DISKGROUP.
func (p *Parser) parseDiskgroupSetClause(opts *nodes.List, optStart int) error {
	if p.isIdentLike() && p.cur.Str == "ATTRIBUTE" {
		p.advance()
		attrItems := &nodes.List{}
		for {
			if p.cur.Type == tokSCONST {
				attrName := p.cur.Str
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
				}
				attrVal := ""
				if p.cur.Type == tokSCONST {
					attrVal = p.cur.Str
					p.advance()
				}
				attrItems.Items = append(attrItems.Items, &nodes.DDLOption{
					Key: attrName, Value: attrVal,
				})
			}
			if p.cur.Type == ',' {
				p.advance()
				continue
			}
			break
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "SET_ATTRIBUTE", Items: attrItems,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
	} else if p.isIdentLike() && p.cur.Str == "PERMISSION" {
		p.advance()
		for p.cur.Type != ';' && p.cur.Type != tokEOF {
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "SET_PERMISSION",
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
	} else if p.isIdentLike() && p.cur.Str == "OWNER" {
		p.advance()
		for p.cur.Type != ';' && p.cur.Type != tokEOF {
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "SET_OWNER",
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
	} else if p.cur.Type == kwGROUP || (p.isIdentLike() && p.cur.Str == "GROUP") {
		p.advance()
		for p.cur.Type != ';' && p.cur.Type != tokEOF {
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "SET_GROUP",
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
	}
	return nil
}

// parseDiskNameList parses a comma-separated list of disk names.
func (p *Parser) parseDiskNameList() (string, error) {
	var names []string
	for {
		if p.isIdentLike() || p.cur.Type == tokIDENT {
			names = append(names, p.cur.Str)
			p.advance()
		}
		if p.cur.Type == ',' {
			next := p.peekNext()
			if next.Type == tokIDENT || (next.Type >= 256 && next.Str != "") {
				p.advance()
				continue
			}
		}
		break
	}
	return strings.Join(names, ","), nil
}

// parseDiskgroupVolumeList parses ALL or comma-separated volume names.
func (p *Parser) parseDiskgroupVolumeList() (string, error) {
	if p.cur.Type == kwALL || (p.isIdentLike() && p.cur.Str == "ALL") {
		p.advance()
		return "ALL", nil
	}
	var names []string
	for {
		if p.isIdentLike() || p.cur.Type == tokIDENT {
			names = append(names, p.cur.Str)
			p.advance()
		}
		if p.cur.Type == ',' {
			p.advance()
			continue
		}
		break
	}
	return strings.Join(names, ","), nil
}

// parseDiskgroupForceWaitPower parses optional FORCE/NOFORCE, POWER n, WAIT/NOWAIT.
func (p *Parser) parseDiskgroupForceWaitPower() (string, error) {
	val := ""
	for {
		if p.cur.Type == kwFORCE {
			val += ":FORCE"
			p.advance()
		} else if p.isIdentLike() && p.cur.Str == "NOFORCE" {
			val += ":NOFORCE"
			p.advance()
		} else if p.isIdentLike() && p.cur.Str == "POWER" {
			p.advance()
			power := ""
			if p.cur.Type == tokICONST {
				power = p.cur.Str
				p.advance()
			} else if p.isIdentLike() && p.cur.Str == "LIMIT" {
				p.advance()
				if p.cur.Type == tokICONST {
					power = p.cur.Str
					p.advance()
				}
			}
			val += ":POWER:" + power
		} else if p.isIdentLike() && p.cur.Str == "WAIT" {
			val += ":WAIT"
			p.advance()
		} else if p.isIdentLike() && p.cur.Str == "NOWAIT" {
			val += ":NOWAIT"
			p.advance()
		} else {
			break
		}
	}
	return val, nil
}

// parseDiskgroupSkipParens skips an optional ATTRIBUTES (...) clause.
func (p *Parser) parseDiskgroupSkipParens() error {
	if p.isIdentLike() && p.cur.Str == "ATTRIBUTES" {
		p.advance()
	}
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
	return nil
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
			if p.isIdentLike() && p.cur.Str == "CONTENTS" {
				p.advance()
			}
			val := "INCLUDING_CONTENTS"
			if p.cur.Type == kwFORCE {
				val += ":FORCE"
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: val,
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})
		case p.isIdentLike() && p.cur.Str == "EXCLUDING":
			p.advance()
			if p.isIdentLike() && p.cur.Str == "CONTENTS" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "EXCLUDING_CONTENTS",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})
		case p.cur.Type == kwFORCE:
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "FORCE",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})
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
