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
func (p *Parser) parseAlterDatabaseStmt(start int) (nodes.StmtNode, error) {
	// DATABASE already consumed by caller
	stmt := &nodes.AdminDDLStmt{
		Action:     "ALTER",
		ObjectType: nodes.OBJECT_DATABASE,
		Loc:        nodes.Loc{Start: start},
	}

	// Optional database name — if next token is an identifier (not a keyword clause starter)
	if p.isIdentLike() || p.cur.Type == tokQIDENT {
		if !p.isDatabaseClauseKeyword() {
			var parseErr629 error
			stmt.Name, parseErr629 = p.parseObjectName()
			if parseErr629 != nil {
				return nil, parseErr629
			}
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// ---- recovery_clauses ----
		// RECOVER ...
		case p.isIdentLike() && p.cur.Str == "RECOVER":
			p.advance()
			parseErr630 := p.parseAlterDatabaseRecoverClause(opts, optStart)
			if parseErr630 !=

				// BEGIN BACKUP
				nil {
				return nil, parseErr630
			}

		case p.cur.Type == kwBEGIN:
			p.advance()
			if p.isIdentLike() && p.cur.Str == "BACKUP" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "BEGIN_BACKUP",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// END BACKUP
		case p.cur.Type == kwEND:
			p.advance()
			if p.isIdentLike() && p.cur.Str == "BACKUP" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "END_BACKUP",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// ---- database_file_clauses ----
		// RENAME FILE 'f1' [, ...] TO 'f2' [, ...]
		case p.cur.Type == kwRENAME:
			p.advance()
			if p.cur.Type == kwFILE {
				p.advance()
				fromFiles, parseErr631 := p.parseStringList()
				if parseErr631 != nil {
					return nil, parseErr631
				}
				if p.cur.Type == kwTO {
					p.advance()
				}
				toFiles, parseErr632 := p.parseStringList()
				if parseErr632 != nil {
					return nil, parseErr632
				}
				items := &nodes.List{}
				for _, f := range fromFiles {
					items.Items = append(items.Items, &nodes.DDLOption{Key: "FROM", Value: f})
				}
				for _, f := range toFiles {
					items.Items = append(items.Items, &nodes.DDLOption{Key: "TO", Value: f})
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "RENAME_FILE", Items: items,
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc:   nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
			parseErr633 := p.parseAlterDatafileTempfileAction(opts, optStart, "DATAFILE", file)
			if parseErr633 !=

				// TEMPFILE 'file' { ONLINE | OFFLINE | RESIZE | AUTOEXTEND | DROP | END BACKUP }
				nil {
				return nil, parseErr633
			}

		case p.isIdentLike() && p.cur.Str == "TEMPFILE":
			p.advance()
			file := ""
			if p.cur.Type == tokSCONST {
				file = p.cur.Str
				p.advance()
			}
			parseErr634 := p.parseAlterDatafileTempfileAction(opts, optStart, "TEMPFILE", file)
			if parseErr634 !=

				// MOVE DATAFILE 'old' TO 'new'
				nil {
				return nil, parseErr634
			}

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
				Loc:   nodes.Loc{Start: optStart, End: p.prev.End},
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
						Loc:   nodes.Loc{Start: optStart, End: p.prev.End},
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
							df, parseErr635 := p.parseDatafileClause()
							if parseErr635 != nil {
								return nil, parseErr635
							}
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
							Loc:   nodes.Loc{Start: lfStart, End: p.prev.End},
						})
						if p.cur.Type == ',' {
							p.advance()
							continue
						}
						break
					}
					opts.Items = append(opts.Items, &nodes.DDLOption{
						Key: prefix, Items: logItems,
						Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
						Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
						Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			} else if p.isIdentLike() && p.cur.Str == "MIRROR" {
				// DROP MIRROR COPY
				p.advance() // MIRROR
				if p.isIdentLike() && p.cur.Str == "COPY" {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DROP_MIRROR_COPY",
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			} else if p.cur.Type == tokSCONST {
				file := p.cur.Str
				p.advance()
				if p.isIdentLike() && p.cur.Str == "REUSE" {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "BACKUP_CONTROLFILE", Value: file,
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
							Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
						Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
							Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				df, parseErr636 := p.parseDatafileClause()
				if parseErr636 != nil {
					return nil, parseErr636
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			} else if p.isIdentLike() && p.cur.Str == "RESTRICTED" {
				p.advance() // RESTRICTED
				if p.cur.Type == kwSESSION {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "ENABLE_RESTRICTED_SESSION",
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			} else if p.isIdentLike() && p.cur.Str == "RESTRICTED" {
				p.advance() // RESTRICTED
				if p.cur.Type == kwSESSION {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DISABLE_RESTRICTED_SESSION",
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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

		// SUSPEND / RESUME
		case p.isIdentLike() && p.cur.Str == "SUSPEND":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "SUSPEND",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		case p.isIdentLike() && p.cur.Str == "RESUME":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "RESUME",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			} else if p.isIdentLike() && p.cur.Str == "REPLAY" {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "START_REPLAY",
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			} else if p.isIdentLike() && p.cur.Str == "REPLAY" {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "STOP_REPLAY",
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// ARCHIVELOG / NOARCHIVELOG / MANUAL (logfile_clauses)
		case p.isIdentLike() && p.cur.Str == "ARCHIVELOG":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "ARCHIVELOG",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		case p.isIdentLike() && p.cur.Str == "NOARCHIVELOG":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "NOARCHIVELOG",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		case p.isIdentLike() && p.cur.Str == "MANUAL":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "MANUAL",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

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

// parseAlterDatabaseRecoverClause parses the RECOVER sub-clauses of ALTER DATABASE.
//
//	RECOVER [ AUTOMATIC ] [ FROM 'location' ] DATABASE
//	  [ UNTIL { CANCEL | TIME date | CHANGE integer } ]
//	  [ USING BACKUP CONTROLFILE ]
//	RECOVER [ AUTOMATIC ] [ FROM 'location' ] [ STANDBY ] DATAFILE 'file' [, ...]
//	RECOVER [ AUTOMATIC ] [ FROM 'location' ] [ STANDBY ] TABLESPACE name [, ...]
//	RECOVER MANAGED STANDBY DATABASE { CANCEL | DISCONNECT [FROM SESSION] | FINISH | ... }
func (p *Parser) parseAlterDatabaseRecoverClause(opts *nodes.List, optStart int) error {
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
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
		return nil
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
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
		return nil
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})
			return nil
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
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
		return nil
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
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
		return nil
	}

	// Fallback: generic RECOVER
	opts.Items = append(opts.Items, &nodes.DDLOption{
		Key: "RECOVER", Value: strings.TrimSpace(val),
		Loc: nodes.Loc{Start: optStart, End: p.prev.End},
	})
	return nil
}

// parseAlterDatafileTempfileAction parses actions after DATAFILE/TEMPFILE 'file' in ALTER DATABASE.
//
//	{ ONLINE | OFFLINE [ FOR DROP ] | RESIZE size_clause | AUTOEXTEND ... | END BACKUP | DROP [INCLUDING DATAFILES] }
func (p *Parser) parseAlterDatafileTempfileAction(opts *nodes.List, optStart int, prefix string, file string) error {
	switch {
	case p.cur.Type == kwONLINE:
		p.advance()
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: prefix + "_ONLINE", Value: file,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
	case p.isIdentLike() && p.cur.Str == "RESIZE":
		p.advance()
		size, parseErr637 := p.parseSizeValue()
		if parseErr637 != nil {
			return parseErr637
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: prefix + "_RESIZE", Value: file,
			Items: &nodes.List{Items: []nodes.Node{&nodes.DDLOption{Key: "SIZE", Value: size}}},
			Loc:   nodes.Loc{Start: optStart, End: p.prev.End},
		})
	case p.isIdentLike() && p.cur.Str == "AUTOEXTEND":
		ac, parseErr638 := p.parseAutoextendClause()
		if parseErr638 != nil {
			return parseErr638
		}
		items := &nodes.List{}
		if ac != nil {
			items.Items = append(items.Items, ac)
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: prefix + "_AUTOEXTEND", Value: file,
			Items: items,
			Loc:   nodes.Loc{Start: optStart, End: p.prev.End},
		})
	case p.cur.Type == kwEND:
		p.advance()
		if p.isIdentLike() && p.cur.Str == "BACKUP" {
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: prefix + "_END_BACKUP", Value: file,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
	case p.isIdentLike() && p.cur.Str == "ENCRYPT":
		p.advance()
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: prefix + "_ENCRYPT", Value: file,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
	case p.isIdentLike() && p.cur.Str == "DECRYPT":
		p.advance()
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: prefix + "_DECRYPT", Value: file,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
	default:
		// unknown action, still record the file reference
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: prefix, Value: file,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
	}
	return nil
}

// parseStringList parses a comma-separated list of string constants.
func (p *Parser) parseStringList() ([]string, error) {
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
	return result, nil
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

// parseCreatePluggableDatabaseStmt parses a CREATE PLUGGABLE DATABASE statement.
// The CREATE keyword has been consumed. PLUGGABLE DATABASE has been consumed by caller.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-PLUGGABLE-DATABASE.html
//
//	CREATE PLUGGABLE DATABASE pdb_name
//	  { create_pdb_from_seed | create_pdb_clone | create_pdb_from_xml
//	  | create_pdb_from_mirror_copy | create_pdb_decrypt_from_xml }
func (p *Parser) parseCreatePluggableDatabaseStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_PLUGGABLE_DATABASE,
		Loc:        nodes.Loc{Start: start},
	}
	var parseErr639 error

	// Parse PDB name
	stmt.Name, parseErr639 = p.parseObjectName()
	if parseErr639 != nil {
		return nil,

			// Determine which variant: FROM, USING, AS, or from_seed (default)
			parseErr639
	}

	opts := &nodes.List{}

	switch {
	case p.cur.Type == kwFROM:
		// create_pdb_clone or create_pdb_from_mirror_copy
		p.advance() // consume FROM
		if p.isIdentLike() && p.cur.Str == "SEED" {
			p.advance() // consume SEED
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "FROM_SEED"})
		} else {
			// FROM src_pdb_name[@dblink]
			srcName, parseErr640 := p.parseObjectName()
			if parseErr640 != nil {
				return nil, parseErr640
			}
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
			srcName, parseErr641 := p.parseObjectName()
			if parseErr641 != nil {
				return nil, parseErr641
			}
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
	parseErr642 :=

		// Parse common PDB creation clauses
		p.parsePDBCreationClauses(opts)
	if parseErr642 != nil {
		return nil, parseErr642
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parsePDBCreationClauses parses common clauses for CREATE PLUGGABLE DATABASE variants.
func (p *Parser) parsePDBCreationClauses(opts *nodes.List) error {
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		optStart := p.pos()

		switch {
		// USING 'filename' (for XML variant after AS CLONE)
		case p.isIdentLike() && p.cur.Str == "USING":
			p.advance()
			if p.cur.Type == tokSCONST {
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "USING", Value: p.cur.Str,
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// ROLES = (role, ...)
		case p.cur.Type == kwROLE || (p.isIdentLike() && p.cur.Str == "ROLES"):
			p.advance() // ROLES
			if p.cur.Type == '=' {
				p.advance()
			}
			parseErr643 := p.parsePDBParenList(opts, "ROLES", optStart)
			if parseErr643 !=

				// PARALLEL [integer]
				nil {
				return parseErr643
			}

		case p.isIdentLike() && p.cur.Str == "PARALLEL":
			p.advance()
			val := ""
			if p.cur.Type == tokICONST {
				val = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "PARALLEL", Value: val,
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
			parseErr644 := p.parsePDBConvertClause(opts, "FILE_NAME_CONVERT", optStart)
			if parseErr644 !=

				// SERVICE_NAME_CONVERT = (...) | SERVICE_NAME_CONVERT = NONE
				nil {
				return parseErr644
			}

		case p.isIdentLike() && p.cur.Str == "SERVICE_NAME_CONVERT":
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			parseErr645 := p.parsePDBConvertClause(opts, "SERVICE_NAME_CONVERT", optStart)
			if parseErr645 !=

				// SOURCE_FILE_NAME_CONVERT = (...) | SOURCE_FILE_NAME_CONVERT = NONE
				nil {
				return parseErr645
			}

		case p.isIdentLike() && p.cur.Str == "SOURCE_FILE_NAME_CONVERT":
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			parseErr646 := p.parsePDBConvertClause(opts, "SOURCE_FILE_NAME_CONVERT", optStart)
			if parseErr646 !=

				// SOURCE_FILE_DIRECTORY = 'path' | SOURCE_FILE_DIRECTORY = NONE
				nil {
				return parseErr646
			}

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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// STORAGE (...) | STORAGE UNLIMITED
		case p.isIdentLike() && p.cur.Str == "STORAGE":
			p.advance()
			if p.isIdentLike() && p.cur.Str == "UNLIMITED" {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "STORAGE", Value: "UNLIMITED",
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			} else if p.cur.Type == '(' {
				p.advance()
				storageItems := &nodes.List{}
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					sStart := p.pos()
					if p.isIdentLike() && p.cur.Str == "MAXSIZE" {
						p.advance()
						val, parseErr647 := p.parsePDBSizeOrUnlimited()
						if parseErr647 != nil {
							return parseErr647
						}
						storageItems.Items = append(storageItems.Items, &nodes.DDLOption{
							Key: "MAXSIZE", Value: val,
							Loc: nodes.Loc{Start: sStart, End: p.prev.End},
						})
					} else if p.isIdentLike() && p.cur.Str == "MAX_AUDIT_SIZE" {
						p.advance()
						val, parseErr648 := p.parsePDBSizeOrUnlimited()
						if parseErr648 != nil {
							return parseErr648
						}
						storageItems.Items = append(storageItems.Items, &nodes.DDLOption{
							Key: "MAX_AUDIT_SIZE", Value: val,
							Loc: nodes.Loc{Start: sStart, End: p.prev.End},
						})
					} else if p.isIdentLike() && p.cur.Str == "MAX_DIAG_SIZE" {
						p.advance()
						val, parseErr649 := p.parsePDBSizeOrUnlimited()
						if parseErr649 != nil {
							return parseErr649
						}
						storageItems.Items = append(storageItems.Items, &nodes.DDLOption{
							Key: "MAX_DIAG_SIZE", Value: val,
							Loc: nodes.Loc{Start: sStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// TEMPFILE REUSE
		case p.isIdentLike() && p.cur.Str == "TEMPFILE":
			p.advance()
			if p.isIdentLike() && p.cur.Str == "REUSE" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "TEMPFILE_REUSE",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// USER_TABLESPACES = ALL [EXCEPT (...)] | (list) | NONE
		case p.isIdentLike() && p.cur.Str == "USER_TABLESPACES":
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			parseErr650 := p.parsePDBAllExceptNoneList(opts, "USER_TABLESPACES", optStart)
			if parseErr650 !=

				// STANDBYS = ALL [EXCEPT (...)] | (list) | NONE
				nil {
				return parseErr650
			}

		case p.isIdentLike() && p.cur.Str == "STANDBYS":
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			parseErr651 := p.parsePDBAllExceptNoneList(opts, "STANDBYS", optStart)
			if parseErr651 !=

				// LOGGING / NOLOGGING
				nil {
				return parseErr651
			}

		case p.isIdentLike() && p.cur.Str == "LOGGING":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "LOGGING",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})
		case p.isIdentLike() && p.cur.Str == "NOLOGGING":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "NOLOGGING",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			} else if p.isIdentLike() && p.cur.Str == "EXTERNAL" {
				p.advance() // EXTERNAL
				if p.isIdentLike() && p.cur.Str == "STORE" {
					p.advance() // STORE
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "KEYSTORE_EXTERNAL",
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// SNAPSHOT COPY [NO DATA] / SNAPSHOT = NONE | MANUAL | EVERY n HOURS/MINUTES
		case p.isIdentLike() && p.cur.Str == "SNAPSHOT":
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			parseErr652 := p.parsePDBSnapshotClause(opts, optStart)
			if parseErr652 !=

				// REFRESH MODE { MANUAL | EVERY n MINUTES/HOURS | NONE }
				nil {
				return parseErr652
			}

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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			}

		// NO DATA
		case p.isIdentLike() && p.cur.Str == "NO":
			p.advance() // NO
			if p.isIdentLikeStr("DATA") {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "NO_DATA",
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// COPY / MOVE / NOCOPY (for XML variant)
		case p.isIdentLike() && p.cur.Str == "COPY":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "COPY",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})
		case p.isIdentLike() && p.cur.Str == "MOVE":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "MOVE",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})
		case p.isIdentLike() && p.cur.Str == "NOCOPY":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "NOCOPY",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// AS CLONE (for XML variant)
		case p.cur.Type == kwAS:
			p.advance() // AS
			if p.isIdentLike() && p.cur.Str == "CLONE" {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "AS_CLONE",
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			} else if p.isIdentLike() && p.cur.Str == "APPLICATION" {
				p.advance() // APPLICATION
				if p.isIdentLike() && p.cur.Str == "CONTAINER" {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "AS_APPLICATION_CONTAINER",
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			} else if p.isIdentLikeStr("SEED") {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "AS_SEED",
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		default:
			// Unknown token — stop parsing
			return nil
		}
	}
	return nil
}

// parsePDBConvertClause parses FILE_NAME_CONVERT/SERVICE_NAME_CONVERT/SOURCE_FILE_NAME_CONVERT
// = (...) | = NONE
func (p *Parser) parsePDBConvertClause(opts *nodes.List, key string, optStart int) error {
	if p.isIdentLike() && p.cur.Str == "NONE" {
		p.advance()
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: key, Value: "NONE",
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
	}
	return nil
}

// parsePDBAllExceptNoneList parses ALL [EXCEPT (...)] | (list) | NONE for
// USER_TABLESPACES, STANDBYS, INSTANCES, SERVICES clauses.
func (p *Parser) parsePDBAllExceptNoneList(opts *nodes.List, key string, optStart int) error {
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})
		} else {
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: key, Value: "ALL",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})
		}
	} else if p.isIdentLike() && p.cur.Str == "NONE" {
		p.advance()
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: key, Value: "NONE",
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
	} else {
		parseErr653 :=
			// Parenthesized list
			p.parsePDBParenList(opts, key, optStart)
		if parseErr653 !=

			// parsePDBParenList parses a parenthesized list of identifiers.
			nil {
			return parseErr653
		}
	}
	return nil
}

func (p *Parser) parsePDBParenList(opts *nodes.List, key string, optStart int) error {
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
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
	}
	return nil
}

// parsePDBSizeOrUnlimited parses a size value or UNLIMITED keyword.
func (p *Parser) parsePDBSizeOrUnlimited() (string, error) {
	if p.isIdentLike() && p.cur.Str == "UNLIMITED" {
		p.advance()
		return "UNLIMITED", nil
	}
	return p.parseSizeValue()
}

// parsePDBSnapshotClause parses SNAPSHOT clause variants:
// COPY [NO DATA], NONE, MANUAL, EVERY n HOURS/MINUTES
func (p *Parser) parsePDBSnapshotClause(opts *nodes.List, optStart int) error {
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
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
	} else if p.isIdentLike() && p.cur.Str == "NONE" {
		p.advance()
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "SNAPSHOT", Value: "NONE",
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
	} else if p.isIdentLike() && p.cur.Str == "MANUAL" {
		p.advance()
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "SNAPSHOT", Value: "MANUAL",
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
	}
	return nil
}

// parseAlterPluggableDatabaseStmt parses an ALTER PLUGGABLE DATABASE statement.
// The ALTER keyword has been consumed. PLUGGABLE DATABASE has been consumed by caller.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-PLUGGABLE-DATABASE.html
func (p *Parser) parseAlterPluggableDatabaseStmt(start int) (nodes.StmtNode, error) {
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
		var parseErr654 error
		// PDB name
		stmt.Name, parseErr654 = p.parseObjectName()
		if parseErr654 !=

			// Parse the main clause
			nil {
			return nil, parseErr654
		}
	}
	parseErr655 := p.parseAlterPDBClauses(opts)
	if parseErr655 != nil {
		return nil, parseErr655
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
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
func (p *Parser) parseAlterPDBClauses(opts *nodes.List) error {
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})
			parseErr656 :=
				// instances_clause / services_clause
				p.parsePDBInstancesServices(opts)
			if parseErr656 !=

				// CLOSE [IMMEDIATE | ABORT] [instances_clause] [relocate_clause]
				nil {
				return parseErr656
			}

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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})
			parseErr657 :=
				// instances_clause
				p.parsePDBInstancesServices(opts)
			if parseErr657 !=
				// relocate_clause
				nil {
				return parseErr657
			}

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
					Loc: nodes.Loc{Start: relStart, End: p.prev.End},
				})
			} else if p.isIdentLike() && p.cur.Str == "NORELOCATE" {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "NORELOCATE",
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})
			parseErr658 :=
				// instances_clause
				p.parsePDBInstancesServices(opts)
			if parseErr658 !=

				// UNPLUG INTO 'filename' [ENCRYPT USING ...]
				nil {
				return parseErr658
			}

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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
						Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// STORAGE (...)  | STORAGE UNLIMITED
		case p.isIdentLike() && p.cur.Str == "STORAGE":
			p.advance()
			if p.isIdentLike() && p.cur.Str == "UNLIMITED" {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "STORAGE", Value: "UNLIMITED",
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			} else if p.cur.Type == '(' {
				p.advance()
				storageItems := &nodes.List{}
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					sStart := p.pos()
					if p.isIdentLike() && p.cur.Str == "MAXSIZE" {
						p.advance()
						val, parseErr659 := p.parsePDBSizeOrUnlimited()
						if parseErr659 != nil {
							return parseErr659
						}
						storageItems.Items = append(storageItems.Items, &nodes.DDLOption{
							Key: "MAXSIZE", Value: val,
							Loc: nodes.Loc{Start: sStart, End: p.prev.End},
						})
					} else if p.isIdentLike() && p.cur.Str == "MAX_AUDIT_SIZE" {
						p.advance()
						val, parseErr660 := p.parsePDBSizeOrUnlimited()
						if parseErr660 != nil {
							return parseErr660
						}
						storageItems.Items = append(storageItems.Items, &nodes.DDLOption{
							Key: "MAX_AUDIT_SIZE", Value: val,
							Loc: nodes.Loc{Start: sStart, End: p.prev.End},
						})
					} else if p.isIdentLike() && p.cur.Str == "MAX_DIAG_SIZE" {
						p.advance()
						val, parseErr661 := p.parsePDBSizeOrUnlimited()
						if parseErr661 != nil {
							return parseErr661
						}
						storageItems.Items = append(storageItems.Items, &nodes.DDLOption{
							Key: "MAX_DIAG_SIZE", Value: val,
							Loc: nodes.Loc{Start: sStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			}

		// LOGGING / NOLOGGING
		case p.isIdentLike() && p.cur.Str == "LOGGING":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "LOGGING",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})
		case p.isIdentLike() && p.cur.Str == "NOLOGGING":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "NOLOGGING",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			} else if p.isIdentLike() && p.cur.Str == "RECOVERY" {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "ENABLE_RECOVERY",
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			} else if p.isIdentLike() && p.cur.Str == "BACKUP" {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "ENABLE_BACKUP",
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			} else if p.isIdentLike() && p.cur.Str == "RECOVERY" {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DISABLE_RECOVERY",
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
				})
			} else if p.isIdentLike() && p.cur.Str == "BACKUP" {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{
					Key: "DISABLE_BACKUP",
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		// APPLICATION clauses
		case p.isIdentLike() && p.cur.Str == "APPLICATION":
			p.advance()
			parseErr662 := // APPLICATION
				p.parseAlterPDBApplicationClause(opts, optStart)
			if parseErr662 !=

				// SNAPSHOT clauses: SNAPSHOT NONE/MANUAL/EVERY
				nil {
				return parseErr662
			}

		case p.isIdentLike() && p.cur.Str == "SNAPSHOT":
			p.advance()
			parseErr663 := p.parsePDBSnapshotClause(opts, optStart)
			if parseErr663 !=

				// MATERIALIZE
				nil {
				return parseErr663
			}

		case p.isIdentLike() && p.cur.Str == "MATERIALIZE":
			p.advance()
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "MATERIALIZE",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
					Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})

		default:
			// Unknown clause — stop
			return nil
		}
	}
	return nil
}

// parsePDBInstancesServices parses optional INSTANCES = (...) and SERVICES = (...) clauses.
func (p *Parser) parsePDBInstancesServices(opts *nodes.List) error {
	for {
		optStart := p.pos()
		if p.isIdentLike() && p.cur.Str == "INSTANCES" {
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			parseErr664 := p.parsePDBAllExceptNoneList(opts, "INSTANCES", optStart)
			if parseErr664 != nil {
				return parseErr664
			}
		} else if p.isIdentLike() && p.cur.Str == "SERVICES" {
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			parseErr665 := p.parsePDBAllExceptNoneList(opts, "SERVICES", optStart)
			if parseErr665 != nil {
				return parseErr665

				// parseAlterPDBApplicationClause parses APPLICATION clauses in ALTER PLUGGABLE DATABASE.
			}
		} else {
			return nil
		}
	}
	return nil
}

func (p *Parser) parseAlterPDBApplicationClause(opts *nodes.List, optStart int) error {
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
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})
		} else {
			// ALL SYNC
			if p.isIdentLike() && p.cur.Str == "SYNC" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "APP_ALL_SYNC",
				Loc: nodes.Loc{Start: optStart, End: p.prev.End},
			})
		}
		return nil
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
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
		return nil
	}

	appName := ""
	if len(appNames) > 0 {
		appName = appNames[0]
	}

	// Single app operations
	switch {
	case p.cur.Type == kwBEGIN || (p.isIdentLike() && p.cur.Str == "BEGIN"):
		p.advance()
		parseErr666 := // BEGIN
			p.parseAppBeginClause(opts, appName, optStart)
		if parseErr666 != nil {
			return parseErr666
		}

	case p.cur.Type == kwEND || (p.isIdentLike() && p.cur.Str == "END"):
		p.advance()
		parseErr667 := // END
			p.parseAppEndClause(opts, appName, optStart)
		if parseErr667 != nil {
			return parseErr667
		}

	case p.cur.Type == kwSET:
		p.advance()
		parseErr668 := // SET
			p.parseAppSetClause(opts, appName, optStart)
		if parseErr668 != nil {
			return parseErr668
		}

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
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
	}
	return nil
}

// parseAppBeginClause parses APPLICATION app_name BEGIN { INSTALL | PATCH | UPGRADE | UNINSTALL }.
func (p *Parser) parseAppBeginClause(opts *nodes.List, appName string, optStart int) error {
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
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	case p.isIdentLike() && p.cur.Str == "UNINSTALL":
		p.advance()
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "APP_BEGIN_UNINSTALL", Value: appName,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
	}
	return nil
}

// parseAppEndClause parses APPLICATION app_name END { INSTALL | PATCH | UPGRADE | UNINSTALL }.
func (p *Parser) parseAppEndClause(opts *nodes.List, appName string, optStart int) error {
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
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})

	case p.isIdentLike() && p.cur.Str == "UNINSTALL":
		p.advance()
		opts.Items = append(opts.Items, &nodes.DDLOption{
			Key: "APP_END_UNINSTALL", Value: appName,
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
	}
	return nil
}

// parseAppSetClause parses APPLICATION app_name SET { PATCH | VERSION | COMPATIBILITY }.
func (p *Parser) parseAppSetClause(opts *nodes.List, appName string, optStart int) error {
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
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
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
			Loc: nodes.Loc{Start: optStart, End: p.prev.End},
		})
	}
	return nil
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
