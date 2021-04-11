package cmd_definitions

import "github.com/wojnosystems/flick/pkg/cmd_to_method_map"

type ServiceDesc struct {
	Root    MethodDesc
	Methods cmd_to_method_map.Map
}
