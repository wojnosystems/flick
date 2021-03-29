package goland

type goImport struct {
	Path  string
	Alias string
}

func (i goImport) Empty() bool {
	return len(i.Path) == 0 && len(i.Alias) == 0
}

// map["import.Path"] = goImport
type importRegistryType map[string]goImport
