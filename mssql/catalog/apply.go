package catalog

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql/ast"
)

// parseSpatialOptions reads the WITH options of CREATE SPATIAL INDEX into the
// metadata Bytebase syncs and renders.
func parseSpatialOptions(options *ast.List, index *metadata.IndexMetadata) {
	if options == nil {
		return
	}
	index.SpatialConfig.Storage = &metadata.StorageConfig{
		AllowRowLocks:  true,
		AllowPageLocks: true,
	}
	for _, item := range options.Items {
		option, ok := item.(*ast.String)
		if !ok {
			continue
		}
		name, value, ok := strings.Cut(option.Str, "=")
		if !ok {
			continue
		}
		name = strings.ToUpper(strings.TrimSpace(name))
		value = strings.TrimSpace(value)
		switch name {
		case "BOUNDING_BOX":
			parseBoundingBox(value, index)
		case "GRIDS":
			parseGridLevels(value, index)
		case "CELLS_PER_OBJECT":
			if cellsPerObject, err := strconv.ParseInt(value, 10, 32); err == nil {
				index.SpatialConfig.Tessellation.CellsPerObject = int32(cellsPerObject)
			}
		case "PAD_INDEX", "FILLFACTOR", "SORT_IN_TEMPDB", "ONLINE", "ALLOW_ROW_LOCKS", "ALLOW_PAGE_LOCKS", "MAXDOP", "DATA_COMPRESSION":
			parseSpatialStorageOption(name, value, index)
		}
	}
}

func parseBoundingBox(value string, index *metadata.IndexMetadata) {
	value = strings.TrimSpace(strings.Trim(value, "()"))
	parts := strings.Split(value, ",")
	if len(parts) != 4 {
		return
	}
	var nums [4]float64
	for i, part := range parts {
		num, err := strconv.ParseFloat(strings.TrimSpace(part), 64)
		if err != nil {
			return
		}
		nums[i] = num
	}
	index.SpatialConfig.Tessellation.BoundingBox = &metadata.BoundingBox{
		Xmin: nums[0],
		Ymin: nums[1],
		Xmax: nums[2],
		Ymax: nums[3],
	}
}

func parseGridLevels(value string, index *metadata.IndexMetadata) {
	value = strings.TrimSpace(strings.Trim(value, "()"))
	var gridLevels []*metadata.GridLevel
	for _, part := range strings.Split(value, ",") {
		levelPart, densityPart, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		levelPart = strings.TrimSpace(levelPart)
		densityPart = strings.TrimSpace(densityPart)
		var levelNum int32
		if strings.HasPrefix(strings.ToUpper(levelPart), "LEVEL_") {
			if num, err := strconv.ParseInt(levelPart[6:], 10, 32); err == nil {
				levelNum = int32(num)
			}
		}
		gridLevels = append(gridLevels, &metadata.GridLevel{
			Level:   levelNum,
			Density: strings.ToUpper(densityPart),
		})
	}
	index.SpatialConfig.Tessellation.GridLevels = gridLevels
}

func parseSpatialStorageOption(name, value string, index *metadata.IndexMetadata) {
	if index.SpatialConfig.Storage == nil {
		index.SpatialConfig.Storage = &metadata.StorageConfig{}
	}
	storage := index.SpatialConfig.Storage
	switch name {
	case "PAD_INDEX":
		storage.PadIndex = strings.EqualFold(value, "ON")
	case "FILLFACTOR":
		if fillFactor, err := strconv.ParseInt(value, 10, 32); err == nil {
			storage.Fillfactor = int32(fillFactor)
		}
	case "SORT_IN_TEMPDB":
		storage.SortInTempdb = strings.ToUpper(value)
	case "ONLINE":
		storage.Online = strings.EqualFold(value, "ON")
	case "ALLOW_ROW_LOCKS":
		storage.AllowRowLocks = strings.EqualFold(value, "ON")
	case "ALLOW_PAGE_LOCKS":
		storage.AllowPageLocks = strings.EqualFold(value, "ON")
	case "MAXDOP":
		if maxdop, err := strconv.ParseInt(value, 10, 32); err == nil {
			storage.Maxdop = int32(maxdop)
		}
	case "DATA_COMPRESSION":
		storage.DataCompression = strings.ToUpper(value)
	}
}
