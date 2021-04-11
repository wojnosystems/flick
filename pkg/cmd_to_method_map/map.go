package cmd_to_method_map

import (
	"github.com/wojnosystems/flick/pkg/cmd_definitions"
	"github.com/wojnosystems/go-nested-map/nested_string_map"
)

type Map struct {
	t nested_string_map.T
}

func (m *Map) Get(path ...string) (method cmd_definitions.MethodDesc, ok bool) {
	v, ok := m.t.Get(path...)
	if ok {
		method, ok = v.(cmd_definitions.MethodDesc)
	}
	return
}

func (m *Map) Put(method cmd_definitions.MethodDesc, path ...string) {
	m.t.Put(method, path...)
}
