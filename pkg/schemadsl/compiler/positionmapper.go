package compiler

import (
	"fmt"
	"strings"

	"github.com/authzed/spicedb/pkg/schemadsl/input"
)

type positionMapper struct {
	schemas         []InputSchema
	mappersBySource map[input.InputSource]input.SourcePositionMapper
}

func newPositionMapper(schemas []InputSchema) input.PositionMapper {
	mappersBySource := map[input.InputSource]input.SourcePositionMapper{}
	for _, schema := range schemas {
		mappersBySource[schema.Source] = input.CreateSourcePositionMapper([]byte(schema.SchemaString))
	}

	return &positionMapper{
		schemas:         schemas,
		mappersBySource: mappersBySource,
	}
}

func (pm *positionMapper) RunePositionToLineAndCol(runePosition int, source input.InputSource) (int, int, error) {
	sourceMapper := pm.mappersBySource[source]
	return sourceMapper.RunePositionToLineAndCol(runePosition)
}

func (pm *positionMapper) LineAndColToRunePosition(lineNumber int, colPosition int, source input.InputSource) (int, error) {
	sourceMapper := pm.mappersBySource[source]
	return sourceMapper.LineAndColToRunePosition(lineNumber, colPosition)
}

func (pm *positionMapper) TextForLine(lineNumber int, source input.InputSource) (string, error) {
	for _, schema := range pm.schemas {
		if schema.Source == source {
			lines := strings.Split(schema.SchemaString, "\n")
			return lines[lineNumber], nil
		}
	}
	return "", fmt.Errorf("unknown schema source")
}
