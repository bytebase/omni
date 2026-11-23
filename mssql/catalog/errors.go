package catalog

import "fmt"

// Error is a SQL Server catalog execution error carrying the SQL Server error
// number (matching mssql.Error.Number).
type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("Msg %d: %s", e.Code, e.Message)
}

const (
	ErrSyntax                     = 102
	ErrInvalidTableDef            = 173
	ErrInvalidColumnName          = 207
	ErrInvalidObjectName          = 208
	ErrUnknownType                = 243
	ErrSelectStarWithoutFrom      = 263
	ErrObjectNotFound             = 1088
	ErrDropIndexFilestreamOn      = 1716
	ErrAddDefaultColumnNotFound   = 1752
	ErrFKColumnNotFound           = 1769
	ErrFKColumnCountMismatch      = 1774
	ErrMultiplePrimaryKeys        = 1779
	ErrColumnAlreadyHasDefault    = 1781
	ErrMultipleClusteredIndex     = 1902
	ErrDupIndexColumn             = 1909
	ErrIndexColumnNotFound        = 1911
	ErrDupIndex                   = 1913
	ErrIndexInvalidObject         = 1914
	ErrInvalidFilegroup           = 1921
	ErrViewNotSchemaBound         = 1939
	ErrIndexOnViewNotAllowed      = 1940
	ErrFileTableNotSupported      = 1969
	ErrIncompatibleAlterType      = 2010
	ErrDupColumn                  = 2705
	ErrDupObject                  = 2714
	ErrUnknownColumnType          = 2715
	ErrIndexNotFound              = 2727
	ErrMultipleIdentityColumns    = 2744
	ErrRowguidcolType             = 2761
	ErrCannotDropNotFound         = 3701
	ErrCannotDropWrongKind        = 3705
	ErrDropIndexReferencedByFK    = 3723
	ErrDropConstraintReferencedByFK = 3725
	ErrDropTableReferencedByFK    = 3726
	ErrConstraintNotFound         = 3728
	ErrDropSchemaNotEmpty         = 3729
	ErrOnlineDropNonClustered     = 3745
	ErrDropIndexClusteredOpt      = 3748
	ErrTruncateNotFound           = 4701
	ErrTruncateNotTable           = 4708
	ErrTruncateReferencedByFK     = 4712
	ErrAlterTableNotFound         = 4902
	ErrAlterTableNotTable         = 4909
	ErrCheckConstraintNotFound    = 4917
	ErrAlterColPersistedNonComp   = 4919
	ErrDropLastColumn             = 4923
	ErrAlterColNotFound           = 4924
	ErrDropRowguidcolMissing      = 4926
	ErrSwitchSameTable            = 4955
	ErrAlterColNFRNotIdentity     = 4991
	ErrColumnDependentObject      = 5074
	ErrXmlIndexNonXmlCol          = 6330
	ErrXmlIndexMissingPrimary     = 6333
	ErrXmlIndexNotTable           = 6334
	ErrAssemblyNotFound           = 6528
	ErrRebuildPartitionInvalid    = 7729
	ErrAlterTableAddIndexSyntax   = 8135
	ErrDupTargetInAlter           = 8168
	ErrTriggerTargetInvalid       = 8197
	ErrFilteredIndexOnView        = 10610
	ErrNoResumableIndexOp         = 10638
	ErrAlterTableIndexDiskTable   = 10785
	ErrCheckConstraintWrongKind   = 11415
	ErrSequenceZeroIncrement      = 11717
	ErrSpatialIndexNonSpatialCol  = 12002
	ErrSystemTimeInvalidCols      = 13501
	ErrSystemVersioningNoPeriod   = 13510
	ErrDropPeriodSystemVersioned  = 13554
	ErrSystemVersioningNotEnabled = 13591
	ErrDropPeriodNotDefined       = 13593
	ErrAlterColHiddenNotPeriod    = 13735
	ErrEdgeConstraintNonEdge      = 13930
	ErrSecurableNotFound          = 15151
	ErrInvalidMaskingFunction     = 16006
	ErrFileTableDirectoryNonFT    = 33411
	ErrClusteredCSOnView          = 35305
	ErrClusteredCSExists          = 35372
)

func mssqlErr(code int, format string, args ...any) *Error {
	return &Error{Code: code, Message: fmt.Sprintf(format, args...)}
}
