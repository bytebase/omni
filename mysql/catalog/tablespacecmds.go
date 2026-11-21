package catalog

import (
	"encoding/json"
	"fmt"
	"math"
	"path"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

const (
	minAutoextendSizeBytes = 4 * 1024 * 1024    // 4M
	maxAutoextendSizeBytes = 4096 * 1024 * 1024 // 4096M (4G)
	innodbDefaultPageSize  = 16384              // 16K
)

func parseMySQLByteSize(raw string) (uint64, bool) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return 0, false
	}
	mult := uint64(1)
	switch s[len(s)-1] {
	case 'k', 'K':
		mult = 1024
		s = s[:len(s)-1]
	case 'm', 'M':
		mult = 1024 * 1024
		s = s[:len(s)-1]
	case 'g', 'G':
		mult = 1024 * 1024 * 1024
		s = s[:len(s)-1]
	}
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, false
	}
	if val == 0 {
		return 0, true
	}
	if val > math.MaxUint64/mult {
		return 0, false
	}
	return val * mult, true
}

// minAutoextendSizeBytes returns the minimum non-zero AUTOEXTEND_SIZE step
// for the catalog's configured innodb_page_size: 4M for 4K/8K/16K, 8M for 32K,
// and 16M for 64K (4 InnoDB extents).
func (c *Catalog) minAutoextendSizeBytes() uint64 {
	if c != nil {
		if step := c.InnoDBPageSize() * 256; step > minAutoextendSizeBytes {
			return step
		}
	}
	return minAutoextendSizeBytes
}

// parseAutoextendSize parses and validates an AUTOEXTEND_SIZE option value as
// InnoDB enforces it on mysql:8.0.46: 0 disables custom autoextend, and any
// non-zero value must lie in [minStep, 4096M] (4025) and be an exact multiple
// of minStep (4023).
func (c *Catalog) parseAutoextendSize(raw string) (uint64, error) {
	if strings.TrimSpace(raw) == "" {
		return 0, nil
	}
	minBytes := c.minAutoextendSizeBytes()
	bytes, ok := parseMySQLByteSize(raw)
	if !ok {
		return 0, errAutoextendSizeOutOfRange(minBytes)
	}
	if bytes == 0 {
		return 0, nil
	}
	if bytes < minBytes || bytes > maxAutoextendSizeBytes {
		return 0, errAutoextendSizeOutOfRange(minBytes)
	}
	if bytes%minBytes != 0 {
		return 0, errAutoextendSizeNotMultiple(minBytes)
	}
	return bytes, nil
}

func (c *Catalog) validateAutoextendSize(raw string) error {
	_, err := c.parseAutoextendSize(raw)
	return err
}

// parseTableAutoextendSize validates and resolves an InnoDB table's effective
// AUTOEXTEND_SIZE option against temporary/system/general tablespace rules
// (measured on mysql:8.0.46: 0 is accepted everywhere; non-zero answers 4026 on
// temporary tables, 1478 on innodb_system, and 4024 on general tablespaces
// ahead of range/alignment checks).
func (c *Catalog) parseTableAutoextendSize(raw string, temporary bool, rawTS string) (uint64, error) {
	if strings.TrimSpace(raw) == "" {
		return 0, nil
	}
	if bytes, ok := parseMySQLByteSize(raw); ok && bytes == 0 {
		return 0, nil
	}
	if temporary {
		return 0, errAutoextendSizeTemporaryTablespace()
	}
	ts := unquoteTableOption(strings.TrimSpace(rawTS))
	if strings.EqualFold(ts, "innodb_system") {
		return 0, errAutoextendSizeSystemTablespace()
	}
	if ts != "" && !strings.EqualFold(ts, "innodb_file_per_table") {
		return 0, errAutoextendSizeGeneralTablespace()
	}
	return c.parseAutoextendSize(raw)
}

// SetInnoDBPageSize sets the server's configured innodb_page_size in bytes
// (allowed values: 4096, 8192, 16384, 32768, 65536; 0 resets to 16384).
func (c *Catalog) SetInnoDBPageSize(bytes uint64) {
	switch bytes {
	case 0:
		c.innodbPageSize = 0
	case 4096, 8192, 16384, 32768, 65536:
		c.innodbPageSize = bytes
	}
}

// InnoDBPageSize returns the server's configured innodb_page_size in bytes,
// defaulting to 16384 (16 KiB).
func (c *Catalog) InnoDBPageSize() uint64 {
	if c != nil && c.innodbPageSize != 0 {
		return c.innodbPageSize
	}
	return innodbDefaultPageSize
}

func validateTablespaceEngineAttribute(attr string) error {
	if attr == "" || json.Valid([]byte(attr)) {
		return nil
	}
	return errInvalidJSONAttribute(attr)
}

func (c *Catalog) resolveTablespaceFileBlockSize(raw string) (uint64, error) {
	pageSize := c.InnoDBPageSize()
	if strings.TrimSpace(raw) == "" {
		return pageSize, nil
	}
	bytes, ok := parseMySQLByteSize(raw)
	if !ok {
		return 0, &Error{Code: 1478, SQLState: "HY000",
			Message: fmt.Sprintf("InnoDB does not support FILE_BLOCK_SIZE=%s", strings.TrimSpace(raw))}
	}
	if bytes == 0 || bytes == pageSize {
		if bytes == 0 {
			return pageSize, nil
		}
		return bytes, nil
	}
	switch bytes {
	case 1024, 2048, 4096, 8192, 16384, 32768, 65536:
		if bytes > pageSize {
			return 0, &Error{Code: 1478, SQLState: "HY000",
				Message: fmt.Sprintf("InnoDB: Cannot create a tablespace with FILE_BLOCK_SIZE=%d because INNODB_PAGE_SIZE=%d.", bytes, pageSize)}
		}
		if pageSize <= innodbDefaultPageSize {
			return bytes, nil
		}
		return 0, &Error{Code: 1478, SQLState: "HY000",
			Message: fmt.Sprintf("InnoDB does not support FILE_BLOCK_SIZE=%d", bytes)}
	default:
		return 0, &Error{Code: 1478, SQLState: "HY000",
			Message: fmt.Sprintf("InnoDB does not support FILE_BLOCK_SIZE=%d", bytes)}
	}
}

func (c *Catalog) validateTablespaceDatafile(tsName, df string, hasDF, undo bool) (string, error) {
	if !undo && !hasDF {
		return "", nil
	}
	ext := ".ibd"
	if undo {
		ext = ".ibu"
	}
	if df == "" {
		return "", errTablespaceDatafileProperFilename()
	}
	sepIdx := strings.LastIndexAny(df, `/\`)
	base := df
	if sepIdx >= 0 {
		base = df[sepIdx+1:]
	}
	if base == "" || base == ext {
		return "", errTablespaceDatafileProperFilename()
	}
	// On Linux mysql:8.0.46, InnoDB checks the lowercase '.ibd'/'.ibu' suffix case-sensitively.
	if !strings.HasSuffix(base, ext) {
		return "", errTablespaceDatafileExt(ext)
	}
	baseDir := c.DataDir()
	if sepIdx >= 0 {
		resolvedDir, ok := c.resolveTablespaceDatafileDir(df)
		if !ok {
			return "", errTablespaceDatafileKnownDir()
		}
		fullPath := path.Clean(resolvedDir + "/" + base)
		for _, used := range c.tablespaceFiles {
			if used == fullPath {
				return "", errDupTablespaceDatafile(tsName)
			}
		}
		return fullPath, nil
	}
	fullPath := path.Clean(baseDir + "/" + base)
	for _, used := range c.tablespaceFiles {
		if used == fullPath {
			return "", errDupTablespaceDatafile(tsName)
		}
	}
	return fullPath, nil
}

// resolveTablespaceDatafileDir resolves and validates an InnoDB tablespace
// datafile's directory against configured InnoDB directories (defaulting to the
// server datadir `/var/lib/mysql`, measured on mysql:8.0.46). Circular `..`
// segments following a non-`..` path component are rejected per Fil_path::is_circular.
func (c *Catalog) resolveTablespaceDatafileDir(df string) (string, bool) {
	baseDir := c.DataDir()
	norm := strings.ReplaceAll(df, `\`, "/")
	lastSlash := strings.LastIndexByte(norm, '/')
	if lastSlash < 0 {
		return baseDir, true
	}
	rawDir := norm[:lastSlash]
	seenNamedSeg := false
	for _, seg := range strings.Split(rawDir, "/") {
		if seg == "" || seg == "." {
			continue
		}
		if seg == ".." {
			if seenNamedSeg {
				return "", false
			}
			continue
		}
		seenNamedSeg = true
	}
	var cleanFull string
	if strings.HasPrefix(norm, "/") {
		cleanFull = path.Clean(norm)
	} else {
		cleanFull = path.Clean(baseDir + "/" + norm)
	}
	cleanDir := path.Dir(cleanFull)
	if cleanDir == baseDir || c.tablespaceDirs[cleanDir] {
		return cleanDir, true
	}
	return "", false
}

// validateInnoDBTablespace validates a TABLESPACE option on an InnoDB table or
// partition (measured on mysql:8.0.46). Non-InnoDB engines ignore TABLESPACE.
// General tablespace names are instance-wide and case-sensitive.
func (c *Catalog) validateInnoDBTablespace(tableName string, temporary bool, rawTS string) error {
	ts := unquoteTableOption(strings.TrimSpace(rawTS))
	if ts == "" {
		return nil
	}
	switch {
	case ts == "innodb_file_per_table" || ts == "innodb_system":
		if temporary {
			return &Error{Code: 1478, SQLState: "HY000",
				Message: fmt.Sprintf("InnoDB: Tablespace `%s` cannot contain TEMPORARY tables.", ts)}
		}
		return nil
	case ts == "innodb_temporary":
		if temporary {
			return nil
		}
		return errTempTablespaceNonTempTable()
	case ts == "mysql":
		return errReservedTablespaceName(tableName, "mysql")
	case strings.HasPrefix(ts, "innodb_"):
		// The prefix is case-sensitive: INNODB_X is an ordinary name.
		return errReservedInnoDBTablespacePrefix()
	case c.tablespaces[ts]:
		if temporary {
			return &Error{Code: 1478, SQLState: "HY000",
				Message: fmt.Sprintf("InnoDB: Tablespace `%s` cannot contain TEMPORARY tables.", ts)}
		}
		return nil
	case c.loadingSDL:
		if temporary {
			return &Error{Code: 1478, SQLState: "HY000",
				Message: fmt.Sprintf("InnoDB: Tablespace `%s` cannot contain TEMPORARY tables.", ts)}
		}
		c.tablespaces[ts] = true
		return nil
	case c.tablespacesPartial && !c.droppedTablespaces[ts]:
		if temporary {
			return &Error{Code: 1478, SQLState: "HY000",
				Message: fmt.Sprintf("InnoDB: Tablespace `%s` cannot contain TEMPORARY tables.", ts)}
		}
		return nil
	default:
		return errTablespaceDoesNotExist(ts)
	}
}

func (c *Catalog) validateInnoDBCompressedTableOption(tableName, rawTS, rowFormat string, keyBlockSize int, forAlter bool) error {
	rfCompressed := strings.EqualFold(strings.TrimSpace(rowFormat), "COMPRESSED")
	if !rfCompressed && keyBlockSize <= 0 {
		return nil
	}
	ts := unquoteTableOption(strings.TrimSpace(rawTS))
	if strings.EqualFold(ts, "innodb_system") {
		return &Error{Code: 1478, SQLState: "HY000",
			Message: fmt.Sprintf("InnoDB: Tablespace `%s` cannot contain a COMPRESSED table", ts)}
	}
	if c.InnoDBPageSize() > 16384 && (ts == "" || strings.EqualFold(ts, "innodb_file_per_table")) {
		if forAlter {
			return errUnsupportedRowOption("ROW_TYPE")
		}
		return errTableStorageEngineOption(tableName)
	}
	return nil
}

func (c *Catalog) validateInnoDBTableBlockSize(rawTS, rowFormat string, keyBlockSize int) error {
	ts := unquoteTableOption(strings.TrimSpace(rawTS))
	if ts == "" || ts == "innodb_file_per_table" || ts == "innodb_system" || ts == "innodb_temporary" {
		return nil
	}
	if c.tablespaceBlockSizes == nil {
		return nil
	}
	tsBlockSize, ok := c.tablespaceBlockSizes[ts]
	if !ok || tsBlockSize == 0 {
		return nil
	}
	pageSize := c.InnoDBPageSize()
	rfCompressed := strings.EqualFold(strings.TrimSpace(rowFormat), "COMPRESSED")
	if tsBlockSize == pageSize {
		if keyBlockSize > 0 || rfCompressed {
			return &Error{Code: 1478, SQLState: "HY000",
				Message: fmt.Sprintf("InnoDB: Tablespace `%s` cannot contain a COMPRESSED table", ts)}
		}
		return nil
	}
	tblPageSize := pageSize
	if keyBlockSize > 0 {
		tblPageSize = uint64(keyBlockSize) * 1024
	} else if rfCompressed {
		tblPageSize = pageSize / 2
	}
	if tblPageSize != tsBlockSize {
		return &Error{Code: 1478, SQLState: "HY000",
			Message: fmt.Sprintf("InnoDB: Tablespace `%s` uses block size %d and cannot contain a table with physical page size %d", ts, tsBlockSize, tblPageSize)}
	}
	return nil
}

// validateInnoDBTablespaceEncryption checks that an InnoDB table's effective
// ENCRYPTION setting matches the general tablespace's ENCRYPTION state (3825).
func (c *Catalog) validateInnoDBTablespaceEncryption(rawTS, tblEncryption string) error {
	ts := unquoteTableOption(strings.TrimSpace(rawTS))
	if ts == "" || !c.tablespaces[ts] {
		return nil
	}
	enc := unquoteTableOption(strings.TrimSpace(tblEncryption))
	if enc != "" && !strings.EqualFold(enc, "Y") && !strings.EqualFold(enc, "N") {
		return nil
	}
	tblEnc := strings.EqualFold(enc, "Y")
	tsEnc, known := c.tablespaceEncrypted[ts]
	if !known {
		if c.loadingSDL {
			if c.tablespaceEncrypted == nil {
				c.tablespaceEncrypted = make(map[string]bool)
			}
			c.tablespaceEncrypted[ts] = tblEnc
		}
		return nil
	}
	if tsEnc != tblEnc {
		return errTablespaceEncryptionMismatch(tblEnc)
	}
	return nil
}

func (c *Catalog) updateTablespaceTablesEncryption(tsName string, encrypted bool) {
	for _, db := range c.databases {
		for _, tbl := range db.Tables {
			if tbl == nil || !isInnoDBEngine(tbl.Engine) || unquoteTableOption(tbl.Tablespace) != tsName {
				continue
			}
			if encrypted {
				tbl.Encryption = "Y"
			} else {
				tbl.Encryption = ""
			}
		}
	}
}

func validateInnoDBPartitionedSharedTablespace(rawTS string) error {
	ts := unquoteTableOption(strings.TrimSpace(rawTS))
	if ts != "" && ts != "innodb_file_per_table" {
		return errPartitionedTableSharedTablespace()
	}
	return nil
}

func effectivePartEngine(tblEngine, partEngine, subEngine string) string {
	if subEngine != "" {
		return subEngine
	}
	if partEngine != "" {
		return partEngine
	}
	return tblEngine
}

// tableUsesTablespace reports whether tbl occupies general tablespace tsName in
// InnoDB. On mysql:8.0.46, non-InnoDB engines such as MyISAM retain a syntactic
// TABLESPACE clause in SHOW CREATE TABLE without occupying the InnoDB tablespace,
// so DROP TABLESPACE succeeds while only a non-InnoDB table names it.
func tableUsesTablespace(tbl *Table, tsName string) bool {
	if tbl == nil {
		return false
	}
	if isInnoDBEngine(tbl.Engine) && unquoteTableOption(tbl.Tablespace) == tsName {
		return true
	}
	if tbl.Partitioning != nil {
		for _, pd := range tbl.Partitioning.Partitions {
			if pd == nil {
				continue
			}
			if isInnoDBEngine(effectivePartEngine(tbl.Engine, pd.Engine, "")) && unquoteTableOption(pd.Tablespace) == tsName {
				return true
			}
			for _, spd := range pd.SubPartitions {
				if spd != nil && isInnoDBEngine(effectivePartEngine(tbl.Engine, pd.Engine, spd.Engine)) && unquoteTableOption(spd.Tablespace) == tsName {
					return true
				}
			}
		}
	}
	return false
}

func renameTableTablespaceRefs(tbl *Table, oldName, newName string) {
	if tbl == nil {
		return
	}
	if isInnoDBEngine(tbl.Engine) && unquoteTableOption(tbl.Tablespace) == oldName {
		tbl.Tablespace = newName
	}
	if tbl.Partitioning != nil {
		for _, pd := range tbl.Partitioning.Partitions {
			if pd == nil {
				continue
			}
			if isInnoDBEngine(effectivePartEngine(tbl.Engine, pd.Engine, "")) && unquoteTableOption(pd.Tablespace) == oldName {
				pd.Tablespace = newName
			}
			for _, spd := range pd.SubPartitions {
				if spd != nil && isInnoDBEngine(effectivePartEngine(tbl.Engine, pd.Engine, spd.Engine)) && unquoteTableOption(spd.Tablespace) == oldName {
					spd.Tablespace = newName
				}
			}
		}
	}
}

func validateTablespaceCreateName(name string) error {
	if strings.Contains(name, "/") {
		return errTablespaceSlashInName()
	}
	// The reserved prefix is lowercase innodb_ exactly (WL#6205):
	// mysql:8.0.46 takes CREATE TABLESPACE INNODB_X and Innodb_z.
	if strings.HasPrefix(name, "innodb_") {
		return errReservedInnoDBTablespacePrefix()
	}
	if name == "mysql" {
		return errTablespaceExists(name)
	}
	return nil
}

func (c *Catalog) validateTablespaceCreateEngine(engine, op string) error {
	if engine == "" {
		return nil
	}
	if err := validateStorageEngine(engine, c.noEngineSubstitution()); err != nil {
		return err
	}
	canon := storedEngineName(engine)
	if !isInnoDBEngine(canon) {
		if strings.EqualFold(canon, "MRG_MyISAM") {
			canon = "MRG_MYISAM"
		}
		return errTablespaceUnsupportedEngine(canon, op)
	}
	return nil
}

// isGeneralTablespaceName reports a general tablespace or a built-in one, by
// its exact name: tablespace names are case-sensitive, so mysql:8.0.46 answers
// 3510 for DROP UNDO TABLESPACE MYSQL and INNODB_SYSTEM.
func (c *Catalog) isGeneralTablespaceName(name string) bool {
	if c.tablespaces[name] {
		return true
	}
	return name == "innodb_system" || name == "innodb_temporary" || name == "mysql"
}

func (c *Catalog) createTablespace(stmt *nodes.CreateTablespaceStmt) error {
	op := "CREATE TABLESPACE"
	if stmt.Undo {
		op = "CREATE UNDO TABLESPACE"
	}
	if err := c.validateTablespaceCreateEngine(stmt.Engine, op); err != nil {
		return err
	}
	if err := validateTablespaceCreateName(stmt.Name); err != nil {
		return err
	}
	if _, isUndo := c.undoTablespaces[stmt.Name]; isUndo || c.tablespaces[stmt.Name] {
		return errTablespaceExists(stmt.Name)
	}
	blockSize := c.InnoDBPageSize()
	if !stmt.Undo {
		var bsErr error
		blockSize, bsErr = c.resolveTablespaceFileBlockSize(stmt.FileBlockSize)
		if bsErr != nil {
			return bsErr
		}
	}
	if stmt.AutoextendSize != "" {
		if err := c.validateAutoextendSize(stmt.AutoextendSize); err != nil {
			return err
		}
	}
	if err := validateTablespaceEngineAttribute(stmt.EngineAttribute); err != nil {
		return err
	}
	dfBase, err := c.validateTablespaceDatafile(stmt.Name, stmt.DataFile, stmt.HasDataFile || stmt.Undo, stmt.Undo)
	if err != nil {
		return err
	}
	if (stmt.HasEncryption || stmt.Encryption != "") && !strings.EqualFold(stmt.Encryption, "Y") && !strings.EqualFold(stmt.Encryption, "N") {
		return errInvalidEncryptionOption()
	}
	if stmt.Undo {
		c.undoTablespaces[stmt.Name] = true
		if dfBase != "" {
			c.tablespaceFiles[stmt.Name] = dfBase
		}
		return nil
	}
	c.tablespaces[stmt.Name] = true
	delete(c.droppedTablespaces, stmt.Name)
	if c.tablespaceBlockSizes == nil {
		c.tablespaceBlockSizes = make(map[string]uint64)
	}
	c.tablespaceBlockSizes[stmt.Name] = blockSize
	if c.tablespaceEncrypted == nil {
		c.tablespaceEncrypted = make(map[string]bool)
	}
	tsEnc := strings.EqualFold(stmt.Encryption, "Y")
	if !stmt.HasEncryption && stmt.Encryption == "" {
		tsEnc = c.session.DefaultTableEncryption
	}
	c.tablespaceEncrypted[stmt.Name] = tsEnc
	if dfBase != "" {
		c.tablespaceFiles[stmt.Name] = dfBase
	}
	return nil
}

func (c *Catalog) alterTablespace(stmt *nodes.AlterTablespaceStmt) error {
	if stmt.Undo {
		if err := c.validateTablespaceCreateEngine(stmt.Engine, "ALTER UNDO TABLESPACE"); err != nil {
			return err
		}
		if stmt.AutoextendSize != "" {
			if err := c.validateAutoextendSize(stmt.AutoextendSize); err != nil {
				return err
			}
		}
		if (stmt.HasEncryption || stmt.Encryption != "") && !strings.EqualFold(stmt.Encryption, "Y") && !strings.EqualFold(stmt.Encryption, "N") {
			return errInvalidEncryptionOption()
		}
		if err := validateTablespaceEngineAttribute(stmt.EngineAttribute); err != nil {
			return err
		}
		if c.isGeneralTablespaceName(stmt.Name) {
			return errWrongTablespaceKind("ALTER UNDO TABLESPACE", stmt.Name, "a general", "ALTER TABLESPACE")
		}
		if _, exists := c.undoTablespaces[stmt.Name]; !exists {
			return errTablespaceDoesNotExist(stmt.Name)
		}
		if stmt.SetInactive {
			c.undoTablespaces[stmt.Name] = false
		} else if stmt.SetActive {
			c.undoTablespaces[stmt.Name] = true
		}
		return nil
	}
	if _, isUndo := c.undoTablespaces[stmt.Name]; isUndo {
		return errWrongTablespaceKind("ALTER TABLESPACE", stmt.Name, "an undo", "ALTER UNDO TABLESPACE")
	}
	if stmt.Name == "mysql" {
		return c.alterMySQLTablespace(stmt)
	}
	// Only the exact built-in names exist beside the general tablespaces:
	// ALTER TABLESPACE innodb_nosuch and INNODB_SYSTEM are 3510 on
	// mysql:8.0.46.
	exists := c.tablespaces[stmt.Name] || stmt.Name == "innodb_system" || stmt.Name == "innodb_temporary"
	if !exists {
		return errTablespaceDoesNotExist(stmt.Name)
	}
	if exists && stmt.Engine != "" && !strings.EqualFold(stmt.Engine, "InnoDB") {
		return errTablespaceEngineMismatch(stmt.Engine, stmt.Name)
	}
	if stmt.AutoextendSize != "" {
		if err := c.validateAutoextendSize(stmt.AutoextendSize); err != nil {
			return err
		}
	}
	if err := validateTablespaceEngineAttribute(stmt.EngineAttribute); err != nil {
		return err
	}
	if (stmt.HasEncryption || stmt.Encryption != "") && !strings.EqualFold(stmt.Encryption, "Y") && !strings.EqualFold(stmt.Encryption, "N") {
		return errInvalidEncryptionOption()
	}
	if err := validateTablespaceCreateName(stmt.Name); err != nil {
		return err
	}
	if !c.tablespaces[stmt.Name] {
		return errTablespaceDoesNotExist(stmt.Name)
	}
	if stmt.HasAddDataFile || stmt.AddDataFile != "" {
		return errAlterTablespaceFailed(stmt.Name)
	}
	if stmt.HasDropDataFile || stmt.DropDataFile != "" {
		if stmt.DropDataFile == "" {
			return errTablespaceNoSuchFile(stmt.Name, "")
		}
		return errAlterTablespaceFailed(stmt.Name)
	}
	targetName := stmt.Name
	if stmt.RenameTo != "" {
		if err := validateTablespaceCreateName(stmt.RenameTo); err != nil {
			return err
		}
		// The target may not exist, the tablespace itself included:
		// mysql:8.0.46 answers 1813 for ALTER TABLESPACE ts RENAME TO ts and
		// takes RENAME TO TS, a name of its own.
		if _, isUndo := c.undoTablespaces[stmt.RenameTo]; c.tablespaces[stmt.RenameTo] || isUndo {
			return errTablespaceExists(stmt.RenameTo)
		}
		delete(c.tablespaces, stmt.Name)
		c.tablespaces[stmt.RenameTo] = true
		if stmt.RenameTo != stmt.Name {
			c.markTablespaceDropped(stmt.Name)
		}
		delete(c.droppedTablespaces, stmt.RenameTo)
		if bs, ok := c.tablespaceBlockSizes[stmt.Name]; ok && stmt.RenameTo != stmt.Name {
			delete(c.tablespaceBlockSizes, stmt.Name)
			c.tablespaceBlockSizes[stmt.RenameTo] = bs
		}
		if enc, ok := c.tablespaceEncrypted[stmt.Name]; ok && stmt.RenameTo != stmt.Name {
			delete(c.tablespaceEncrypted, stmt.Name)
			c.tablespaceEncrypted[stmt.RenameTo] = enc
		}
		if df, ok := c.tablespaceFiles[stmt.Name]; ok && stmt.RenameTo != stmt.Name {
			delete(c.tablespaceFiles, stmt.Name)
			c.tablespaceFiles[stmt.RenameTo] = df
		}
		for _, db := range c.databases {
			for _, tbl := range db.Tables {
				renameTableTablespaceRefs(tbl, stmt.Name, stmt.RenameTo)
			}
		}
		targetName = stmt.RenameTo
	}
	if stmt.HasEncryption || stmt.Encryption != "" {
		enc := strings.EqualFold(stmt.Encryption, "Y")
		if c.tablespaceEncrypted == nil {
			c.tablespaceEncrypted = make(map[string]bool)
		}
		c.tablespaceEncrypted[targetName] = enc
		c.updateTablespaceTablesEncryption(targetName, enc)
	}
	return nil
}

// alterMySQLTablespace applies ALTER TABLESPACE to the data dictionary
// tablespace `mysql` (the name is case-sensitive; MYSQL is no tablespace). It
// exists on every server, is InnoDB, keeps its one file mysql.ibd and its name,
// and takes AUTOEXTEND_SIZE, ENGINE_ATTRIBUTE and ENCRYPTION. ENCRYPTION is the
// documented way to encrypt it; a server with no keyring answers 3185, which
// the catalog, like CREATE TABLESPACE ... ENCRYPTION='Y', does not model.
func (c *Catalog) alterMySQLTablespace(stmt *nodes.AlterTablespaceStmt) error {
	if stmt.Engine != "" && !strings.EqualFold(stmt.Engine, "InnoDB") {
		return errTablespaceEngineMismatch(stmt.Engine, stmt.Name)
	}
	if stmt.AutoextendSize != "" {
		if err := c.validateAutoextendSize(stmt.AutoextendSize); err != nil {
			return err
		}
	}
	if err := validateTablespaceEngineAttribute(stmt.EngineAttribute); err != nil {
		return err
	}
	if (stmt.HasEncryption || stmt.Encryption != "") && !strings.EqualFold(stmt.Encryption, "Y") && !strings.EqualFold(stmt.Encryption, "N") {
		return errInvalidEncryptionOption()
	}
	if stmt.HasAddDataFile || stmt.AddDataFile != "" {
		return errAlterTablespaceFailed(stmt.Name)
	}
	if stmt.HasDropDataFile || stmt.DropDataFile != "" {
		if stmt.DropDataFile != "mysql.ibd" {
			return errTablespaceNoSuchFile(stmt.Name, stmt.DropDataFile)
		}
		return errAlterTablespaceFailed(stmt.Name)
	}
	if stmt.RenameTo != "" {
		// The target is looked up by its exact name: RENAME TO MYSQL or
		// INNODB_SYSTEM names no tablespace and answers 3119.
		_, taken := c.undoTablespaces[stmt.RenameTo]
		switch stmt.RenameTo {
		case "mysql", "innodb_system", "innodb_temporary":
			taken = true
		}
		if taken || c.tablespaces[stmt.RenameTo] {
			return errTablespaceExists(stmt.RenameTo)
		}
		return errWrongTablespaceNameReserved(stmt.Name)
	}
	if stmt.HasEncryption || stmt.Encryption != "" {
		if c.tablespaceEncrypted == nil {
			c.tablespaceEncrypted = make(map[string]bool)
		}
		c.tablespaceEncrypted[stmt.Name] = strings.EqualFold(stmt.Encryption, "Y")
	}
	return nil
}

func (c *Catalog) dropTablespace(stmt *nodes.DropTablespaceStmt) error {
	if stmt.Undo {
		if err := c.validateTablespaceCreateEngine(stmt.Engine, "DROP UNDO TABLESPACE"); err != nil {
			return err
		}
		if c.isGeneralTablespaceName(stmt.Name) {
			return errWrongTablespaceKind("DROP UNDO TABLESPACE", stmt.Name, "a general", "DROP TABLESPACE")
		}
		// An innodb_ prefix is reserved only for CREATE: DROP UNDO TABLESPACE
		// innodb_x names no tablespace, 3510.
		active, exists := c.undoTablespaces[stmt.Name]
		if !exists {
			return errTablespaceDoesNotExist(stmt.Name)
		}
		if active {
			return errDropUndoTablespaceFailed(stmt.Name)
		}
		delete(c.undoTablespaces, stmt.Name)
		delete(c.tablespaceFiles, stmt.Name)
		return nil
	}
	if _, isUndo := c.undoTablespaces[stmt.Name]; isUndo && !strings.HasPrefix(stmt.Name, "innodb_") {
		return errWrongTablespaceKind("DROP TABLESPACE", stmt.Name, "an undo", "DROP UNDO TABLESPACE")
	}
	// Only the exact built-in names exist beside the general tablespaces, and
	// they are reserved: DROP TABLESPACE innodb_system is 3119, while
	// innodb_nosuch and INNODB_SYSTEM are 3510 on mysql:8.0.46.
	builtin := stmt.Name == "innodb_system" || stmt.Name == "innodb_temporary"
	if !c.tablespaces[stmt.Name] && !builtin {
		return errTablespaceDoesNotExist(stmt.Name)
	}
	if stmt.Engine != "" && !strings.EqualFold(stmt.Engine, "InnoDB") {
		return errTablespaceEngineMismatch(stmt.Engine, stmt.Name)
	}
	if builtin {
		return errWrongTablespaceNameReserved(stmt.Name)
	}
	if !c.tablespaces[stmt.Name] {
		return errTablespaceDoesNotExist(stmt.Name)
	}
	for _, db := range c.databases {
		for _, tbl := range db.Tables {
			if tableUsesTablespace(tbl, stmt.Name) {
				return errTablespaceIsNotEmpty(stmt.Name)
			}
		}
	}
	delete(c.tablespaces, stmt.Name)
	delete(c.tablespaceBlockSizes, stmt.Name)
	delete(c.tablespaceEncrypted, stmt.Name)
	delete(c.tablespaceFiles, stmt.Name)
	c.markTablespaceDropped(stmt.Name)
	return nil
}

func (c *Catalog) markTablespaceDropped(name string) {
	if !c.tablespacesPartial || name == "" {
		return
	}
	if c.droppedTablespaces == nil {
		c.droppedTablespaces = make(map[string]bool)
	}
	c.droppedTablespaces[name] = true
}

// isBuiltinTablespaceName reports the built-in names by their exact spelling:
// INNODB_SYSTEM is an ordinary general tablespace on mysql:8.0.46.
func isBuiltinTablespaceName(ts string) bool {
	switch strings.TrimSpace(ts) {
	case "", "innodb_file_per_table", "innodb_system", "innodb_temporary":
		return true
	}
	return false
}

func (c *Catalog) registerMetadataTableTablespaces(tbl *Table) {
	if c == nil || tbl == nil || !isInnoDBEngine(tbl.Engine) {
		return
	}
	registerTS := func(raw string, encrypted bool) {
		ts := unquoteTableOption(strings.TrimSpace(raw))
		if isBuiltinTablespaceName(ts) {
			return
		}
		if c.tablespaces == nil {
			c.tablespaces = make(map[string]bool)
		}
		c.tablespaces[ts] = true
		delete(c.droppedTablespaces, ts)
		if c.tablespaceEncrypted == nil {
			c.tablespaceEncrypted = make(map[string]bool)
		}
		if _, known := c.tablespaceEncrypted[ts]; !known {
			c.tablespaceEncrypted[ts] = encrypted
		}
	}
	tblEnc := strings.EqualFold(unquoteTableOption(strings.TrimSpace(tbl.Encryption)), "Y")
	registerTS(tbl.Tablespace, tblEnc)
	if tbl.Partitioning != nil {
		for _, pd := range tbl.Partitioning.Partitions {
			if pd == nil {
				continue
			}
			registerTS(pd.Tablespace, tblEnc)
			for _, spd := range pd.SubPartitions {
				if spd != nil {
					registerTS(spd.Tablespace, tblEnc)
				}
			}
		}
	}
}

func tableHasDiscardedTablespace(tbl *Table) bool {
	if tbl == nil {
		return false
	}
	if tbl.TablespaceDiscarded {
		return true
	}
	if tbl.Partitioning != nil {
		for _, pd := range tbl.Partitioning.Partitions {
			if pd != nil && pd.Discarded {
				return true
			}
		}
	}
	return false
}

// alterDiscardImportTablespace handles ALTER TABLE ... DISCARD/IMPORT TABLESPACE
// as InnoDB enforces it on mysql:8.0.46.
func (c *Catalog) alterDiscardImportTablespace(tbl *Table, partitionNames []string, isImport bool) error {
	if !isInnoDBEngine(tbl.Engine) {
		return errTableStorageEngineOption(tbl.Name)
	}
	if tbl.Temporary {
		return errCannotDiscardImportTempTable()
	}
	if tbl.Tablespace != "" && !strings.EqualFold(unquoteTableOption(tbl.Tablespace), "innodb_file_per_table") {
		return errNotAllowedCommand()
	}
	dbName := ""
	if tbl.Database != nil {
		dbName = tbl.Database.Name
	}
	if len(partitionNames) > 0 {
		if tbl.Partitioning == nil {
			return errPartitionMgmtOnNonPartitioned()
		}
		targetAll := len(partitionNames) == 1 && strings.EqualFold(partitionNames[0], "ALL")
		var targets []*PartitionDefInfo
		if targetAll {
			targets = tbl.Partitioning.Partitions
		} else {
			for _, name := range partitionNames {
				var found *PartitionDefInfo
				for _, pd := range tbl.Partitioning.Partitions {
					if pd != nil && strings.EqualFold(pd.Name, name) {
						found = pd
						break
					}
				}
				if found == nil {
					return &Error{Code: 1735, SQLState: "HY000",
						Message: fmt.Sprintf("Unknown partition '%s' in table '%s'", name, tbl.Name)}
				}
				targets = append(targets, found)
			}
		}
		if !isImport {
			if c.foreignKeyChecks && dbName != "" {
				if _, _, referenced := c.findReferencingFK(dbName, tbl.Name); referenced {
					return errRowIsReferenced2()
				}
			}
			for _, pd := range targets {
				pd.Discarded = true
			}
			return nil
		}
		for _, pd := range targets {
			if !pd.Discarded && !tbl.TablespaceDiscarded {
				return errTablespaceExists(fmt.Sprintf("%s/%s#p#%s", dbName, tbl.Name, pd.Name))
			}
		}
		return errTablespaceMissing(dbName, tbl.Name)
	}
	if !isImport {
		if c.foreignKeyChecks && dbName != "" {
			if _, _, referenced := c.findReferencingFK(dbName, tbl.Name); referenced {
				return errRowIsReferenced2()
			}
		}
		// On mysql:8.0.46, repeating DISCARD TABLESPACE after a prior discard
		// succeeds (0) with Warning 1812 rather than failing with 1814.
		tbl.TablespaceDiscarded = true
		if tbl.Partitioning != nil {
			for _, pd := range tbl.Partitioning.Partitions {
				if pd != nil {
					pd.Discarded = true
				}
			}
		}
		return nil
	}
	if !tbl.TablespaceDiscarded {
		if tbl.Partitioning != nil && len(tbl.Partitioning.Partitions) > 0 {
			for _, pd := range tbl.Partitioning.Partitions {
				if pd != nil && !pd.Discarded {
					return errTablespaceExists(fmt.Sprintf("%s/%s#p#%s", dbName, tbl.Name, pd.Name))
				}
			}
			return errTablespaceMissing(dbName, tbl.Name)
		}
		return errTablespaceExists(fmt.Sprintf("%s/%s", dbName, tbl.Name))
	}
	return errTablespaceMissing(dbName, tbl.Name)
}

// validatePartitionTablespaces checks TABLESPACE options on partition and
// subpartition definitions when their effective storage engine is InnoDB.
func (c *Catalog) validatePartitionTablespaces(tableName string, temporary bool, tableEngine string, defs []*nodes.PartitionDef) error {
	for _, pd := range defs {
		if pd == nil {
			continue
		}
		partEngine := tableEngine
		partTS := ""
		for _, opt := range pd.Options {
			if opt == nil {
				continue
			}
			if engineOptionName(opt.Name) {
				partEngine = storedEngineName(opt.Value)
			} else if strings.EqualFold(opt.Name, "tablespace") {
				partTS = opt.Value
			}
		}
		if partTS != "" && isInnoDBEngine(partEngine) {
			if err := c.validateInnoDBTablespace(tableName, temporary, partTS); err != nil {
				return err
			}
			if !c.loadingSDL {
				if err := validateInnoDBPartitionedSharedTablespace(partTS); err != nil {
					return err
				}
			}
		}
		for _, sub := range pd.SubPartitions {
			if sub == nil {
				continue
			}
			subEngine := partEngine
			subTS := ""
			for _, opt := range sub.Options {
				if opt == nil {
					continue
				}
				if engineOptionName(opt.Name) {
					subEngine = storedEngineName(opt.Value)
				} else if strings.EqualFold(opt.Name, "tablespace") {
					subTS = opt.Value
				}
			}
			if subTS != "" && isInnoDBEngine(subEngine) {
				if err := c.validateInnoDBTablespace(tableName, temporary, subTS); err != nil {
					return err
				}
				if !c.loadingSDL {
					if err := validateInnoDBPartitionedSharedTablespace(subTS); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
