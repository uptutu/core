package util

import "fmt"

const (
	// core.mapper.{type}.{entityID}.{name}.
	fmtMapperString = "core.mapper.%s.%s.%s"
)

func FormatMapper(typ, id, name string) string {
	return fmt.Sprintf(fmtMapperString, typ, id, name)
}
