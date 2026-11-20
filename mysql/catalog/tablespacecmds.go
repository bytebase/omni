package catalog

import (
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

// parseAutoextendSize parses and validates an AUTOEXTEND_SIZE option value as
// InnoDB enforces it on mysql:8.0.46: 0 disables custom autoextend, and any
// non-zero value must lie in [4M, 4096M] (4025) and be an exact multiple of 4M (4023).
func parseAutoextendSize(raw string) (uint64, error) {
	if strings.TrimSpace(raw) == "" {
		return 0, nil
	}
	bytes, ok := parseMySQLByteSize(raw)
	if !ok {
		return 0, errAutoextendSizeOutOfRange()
	}
	if bytes == 0 {
		return 0, nil
	}
	if bytes < minAutoextendSizeBytes || bytes > maxAutoextendSizeBytes {
		return 0, errAutoextendSizeOutOfRange()
	}
	if bytes%minAutoextendSizeBytes != 0 {
		return 0, errAutoextendSizeNotMultiple()
	}
	return bytes, nil
}

func validateAutoextendSize(raw string) error {
	_, err := parseAutoextendSize(raw)
	return err
}

// parseTableAutoextendSize validates and resolves an InnoDB table's effective
// AUTOEXTEND_SIZE option against temporary/system/general tablespace rules
// (measured on mysql:8.0.46: 0 is accepted everywhere; non-zero answers 4026 on
// temporary tables, 1478 on innodb_system, and 4024 on general tablespaces
// ahead of range/alignment checks).
func parseTableAutoextendSize(raw string, temporary bool, rawTS string) (uint64, error) {
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
	return parseAutoextendSize(raw)
}

func resolveTablespaceFileBlockSize(raw string) (uint64, error) {
	if strings.TrimSpace(raw) == "" {
		return innodbDefaultPageSize, nil
	}
	bytes, ok := parseMySQLByteSize(raw)
	if !ok {
		return 0, &Error{Code: 1478, SQLState: "HY000",
			Message: fmt.Sprintf("InnoDB does not support FILE_BLOCK_SIZE=%s", strings.TrimSpace(raw))}
	}
	if bytes == 0 {
		return innodbDefaultPageSize, nil
	}
	switch bytes {
	case 1024, 2048, 4096, 8192, innodbDefaultPageSize:
		return bytes, nil
	case 32768, 65536:
		return 0, &Error{Code: 1478, SQLState: "HY000",
			Message: fmt.Sprintf("InnoDB: Cannot create a tablespace with FILE_BLOCK_SIZE=%d because INNODB_PAGE_SIZE=%d.", bytes, innodbDefaultPageSize)}
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
		return nil
	case ts == "innodb_temporary":
		if temporary {
			return nil
		}
		return errTempTablespaceNonTempTable()
	case ts == "mysql":
		return errReservedTablespaceName(tableName, "mysql")
	case strings.HasPrefix(strings.ToLower(ts), "innodb_"):
		return errReservedInnoDBTablespacePrefix()
	case c.tablespaces[ts]:
		return nil
	case c.loadingSDL:
		// SDL captures schema-scoped DDL without instance-wide CREATE TABLESPACE
		// statements; accept and register general tablespaces referenced by SDL tables.
		c.tablespaces[ts] = true
		return nil
	default:
		return errTablespaceDoesNotExist(ts)
	}
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
	rfCompressed := strings.EqualFold(strings.TrimSpace(rowFormat), "COMPRESSED")
	if tsBlockSize == innodbDefaultPageSize {
		if keyBlockSize > 0 || rfCompressed {
			return &Error{Code: 1478, SQLState: "HY000",
				Message: fmt.Sprintf("InnoDB: Tablespace `%s` cannot contain a COMPRESSED table", ts)}
		}
		return nil
	}
	var tblPageSize uint64 = innodbDefaultPageSize
	if keyBlockSize > 0 {
		tblPageSize = uint64(keyBlockSize) * 1024
	} else if rfCompressed {
		tblPageSize = 8192
	}
	if tblPageSize != tsBlockSize {
		return &Error{Code: 1478, SQLState: "HY000",
			Message: fmt.Sprintf("InnoDB: Tablespace `%s` uses block size %d and cannot contain a table with physical page size %d", ts, tsBlockSize, tblPageSize)}
	}
	return nil
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
	if strings.HasPrefix(strings.ToLower(name), "innodb_") {
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

func (c *Catalog) isGeneralTablespaceName(name string) bool {
	if c.tablespaces[name] {
		return true
	}
	lower := strings.ToLower(name)
	return lower == "innodb_system" || lower == "innodb_temporary" || lower == "mysql"
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
	blockSize := uint64(innodbDefaultPageSize)
	if !stmt.Undo {
		var bsErr error
		blockSize, bsErr = resolveTablespaceFileBlockSize(stmt.FileBlockSize)
		if bsErr != nil {
			return bsErr
		}
	}
	if stmt.AutoextendSize != "" {
		if err := validateAutoextendSize(stmt.AutoextendSize); err != nil {
			return err
		}
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
	if c.tablespaceBlockSizes == nil {
		c.tablespaceBlockSizes = make(map[string]uint64)
	}
	c.tablespaceBlockSizes[stmt.Name] = blockSize
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
			if err := validateAutoextendSize(stmt.AutoextendSize); err != nil {
				return err
			}
		}
		if (stmt.HasEncryption || stmt.Encryption != "") && !strings.EqualFold(stmt.Encryption, "Y") && !strings.EqualFold(stmt.Encryption, "N") {
			return errInvalidEncryptionOption()
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
	lower := strings.ToLower(stmt.Name)
	exists := c.tablespaces[stmt.Name] || lower == "innodb_system" || lower == "innodb_temporary"
	if !exists && !strings.HasPrefix(lower, "innodb_") && lower != "mysql" {
		return errTablespaceDoesNotExist(stmt.Name)
	}
	if exists && stmt.Engine != "" && !strings.EqualFold(stmt.Engine, "InnoDB") {
		return errTablespaceEngineMismatch(stmt.Engine, stmt.Name)
	}
	if stmt.AutoextendSize != "" {
		if err := validateAutoextendSize(stmt.AutoextendSize); err != nil {
			return err
		}
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
	if stmt.RenameTo != "" {
		if err := validateTablespaceCreateName(stmt.RenameTo); err != nil {
			return err
		}
		if _, isUndo := c.undoTablespaces[stmt.RenameTo]; stmt.RenameTo != stmt.Name && (c.tablespaces[stmt.RenameTo] || isUndo) {
			return errTablespaceExists(stmt.RenameTo)
		}
		delete(c.tablespaces, stmt.Name)
		c.tablespaces[stmt.RenameTo] = true
		if bs, ok := c.tablespaceBlockSizes[stmt.Name]; ok && stmt.RenameTo != stmt.Name {
			delete(c.tablespaceBlockSizes, stmt.Name)
			c.tablespaceBlockSizes[stmt.RenameTo] = bs
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
		if strings.HasPrefix(strings.ToLower(stmt.Name), "innodb_") {
			return errReservedInnoDBTablespacePrefix()
		}
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
	if _, isUndo := c.undoTablespaces[stmt.Name]; isUndo && !strings.HasPrefix(strings.ToLower(stmt.Name), "innodb_") {
		return errWrongTablespaceKind("DROP TABLESPACE", stmt.Name, "an undo", "DROP UNDO TABLESPACE")
	}
	lower := strings.ToLower(stmt.Name)
	exists := c.tablespaces[stmt.Name] || lower == "innodb_system" || lower == "innodb_temporary"
	if !exists && !strings.HasPrefix(lower, "innodb_") {
		return errTablespaceDoesNotExist(stmt.Name)
	}
	if exists && stmt.Engine != "" && !strings.EqualFold(stmt.Engine, "InnoDB") {
		return errTablespaceEngineMismatch(stmt.Engine, stmt.Name)
	}
	if strings.HasPrefix(lower, "innodb_") {
		return errReservedInnoDBTablespacePrefix()
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
	delete(c.tablespaceFiles, stmt.Name)
	return nil
}

// alterDiscardImportTablespace handles ALTER TABLE ... DISCARD/IMPORT TABLESPACE
// as InnoDB enforces it on mysql:8.0.46.
func (c *Catalog) alterDiscardImportTablespace(tbl *Table, isImport bool) error {
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
	if !isImport {
		if c.foreignKeyChecks && dbName != "" {
			if _, _, referenced := c.findReferencingFK(dbName, tbl.Name); referenced {
				return errRowIsReferenced2()
			}
		}
		tbl.TablespaceDiscarded = true
		return nil
	}
	if !tbl.TablespaceDiscarded {
		if tbl.Partitioning != nil && len(tbl.Partitioning.Partitions) > 0 {
			return errTablespaceExists(fmt.Sprintf("%s/%s#p#%s", dbName, tbl.Name, tbl.Partitioning.Partitions[0].Name))
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
