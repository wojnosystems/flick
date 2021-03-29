package goland

import (
	"strings"
)

type optionType struct {
	Import         goImport
	ImportOptional goImport
	Type           string
	OptionalType   string
}

// map[value_type like int] = option Type
type optionTypeRegistry map[string]optionType

func registerOptionalTypes(registry optionTypeRegistry) optionTypeRegistry {
	for _, s := range []string{"int", "uint", "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "string", "float32", "float64", "bool", "byte", "rune"} {
		registry[s] = optionType{
			Import: goImport{},
			ImportOptional: goImport{
				Path: goOptionalLibraryImportPath,
			},
			Type:         s,
			OptionalType: "optional." + strings.Title(s),
		}
	}
	for _, s := range []string{"duration", "time"} {
		registry[s] = optionType{
			Import: goImport{
				Path: "time",
			},
			ImportOptional: goImport{
				Path: goOptionalLibraryImportPath,
			},
			Type:         "time." + strings.Title(s),
			OptionalType: "optional." + strings.Title(s),
		}
	}
	return registry
}
