package goland

import (
	"context"
	"errors"
	"fmt"
	"github.com/wojnosystems/flick/generate/dsl"
	"github.com/wojnosystems/flick/string_set"
	"github.com/wojnosystems/flick/string_writer"
	"io"
	"sort"
	"strings"
)

type GoLang struct {
	packageNameValue            string
	interfaceNameValue          string
	structNameValue             string
	globalOptionStructNameValue string
	optionTypes                 optionTypeRegistry
}

const (
	likelyCommandMaxNestingDepth = 10
	singleIndent                 = `  `

	defaultPackageName            = "flickstub"
	defaultInterfaceName          = "Interface"
	defaultStructName             = "Unimplemented"
	defaultGlobalOptionStructName = "AllCommand"

	goOptionalLibraryImportPath = "github.com/wojnosystems/go-optional/v2"
)

type golangSubStruct struct {
	name       string
	parentName string
	options    []dsl.Option
}

func (g *GoLang) Generate(ctx context.Context, document *dsl.Document, output io.Writer) (bytesWritten int, err error) {
	if g.optionTypes == nil {
		g.optionTypes = make(optionTypeRegistry)
		registerOptionalTypes(g.optionTypes)
	}

	out := string_writer.New(
		output)
	defer func() {
		bytesWritten = out.Counter.Count()
	}()
	err = out.Write2Ln("package " + g.packageName())
	if err != nil {
		return
	}

	var imports importRegistryType
	imports, err = collectImports(document, g.optionTypes)
	if err != nil {
		return
	}
	imports["context"] = goImport{
		Path:  "context",
		Alias: "",
	}

	err = out.WriteLn("import (")
	if err != nil {
		return
	}
	importKeysSorted := make([]string, 0, len(imports))
	for s := range imports {
		importKeysSorted = append(importKeysSorted, s)
	}
	sort.Strings(importKeysSorted)
	for _, key := range importKeysSorted {
		record := `"` + imports[key].Path + `"`
		if len(imports[key].Alias) != 0 {
			record = imports[key].Alias + " " + record
		}
		err = out.WriteLnF(`%s%s`, singleIndent, record)
		if err != nil {
			return
		}
	}
	err = out.Write2Ln(")")
	if err != nil {
		return
	}

	err = out.WriteLnF(`type %s interface {`, g.interfaceName())
	if err != nil {
		return
	}
	if len(document.Options) != 0 {
		err = out.WriteLnF(`%sHookBefore(ctx context.Context, opts *%sOptions) error`, singleIndent, g.globalOptionStructName())
		if err != nil {
			return
		}
		err = out.WriteLnF(`%sHookAfter(ctx context.Context, opts *%sOptions, err error) error`, singleIndent, g.globalOptionStructName())
		if err != nil {
			return
		}
	} else {
		err = out.WriteLnF(`%sHookBefore(ctx context.Context) error`, singleIndent)
		if err != nil {
			return
		}
		err = out.WriteLnF(`%sHookAfter(ctx context.Context, err error) error`, singleIndent)
		if err != nil {
			return
		}
	}

	registeredStructs := string_set.Collection{}

	optionStructsToGenerate := make([]golangSubStruct, 0, 10)
	var globalStruct *golangSubStruct
	if len(document.Options) != 0 {
		globalStruct = &golangSubStruct{
			name:       g.globalOptionStructName(),
			parentName: "",
			options:    make([]dsl.Option, len(document.Options)),
		}
		for i, reference := range document.Options {
			globalStruct.options[i] = reference.Option
		}
		optionStructsToGenerate = append(optionStructsToGenerate, *globalStruct)
		registeredStructs.Add(globalStruct.name)
	}

	err = walkCommands(document, func(prefix []string, cmd dsl.Command) (walkErr error) {
		methodName := joinPrefixesAsMethodName(prefix)

		var optionStructName string
		if len(cmd.Options) != 0 {
			optionStructName = methodName
		} else {
			optionStructName = getStructOrBlank(stringSlicePop(prefix), &registeredStructs, globalStruct)
		}

		optionFormal := ""
		if optionStructName != "" {
			optionFormal = fmt.Sprintf(", opts *%sOptions", optionStructName)
		}

		if cmd.Commands.HasAny() {
			err = out.WriteLnF(`%s%sHookBefore(ctx context.Context%s) error`, singleIndent, methodName, optionFormal)
			if err != nil {
				return
			}
			err = out.WriteLnF(`%s%sHookAfter(ctx context.Context%s, err error) error`, singleIndent, methodName, optionFormal)
			if err != nil {
				return
			}
		} else {
			walkErr = out.WriteLnF(`%s%s(ctx context.Context%s) error`, singleIndent, methodName, optionFormal)
			if walkErr != nil {
				return
			}

			if len(cmd.Options) != 0 {
				commandOption := golangSubStruct{
					name:       prefixToOptionStructName(prefix),
					parentName: getStructOrBlank(stringSlicePop(prefix), &registeredStructs, globalStruct),
					options:    make([]dsl.Option, len(cmd.Options)),
				}
				for i, reference := range cmd.Options {
					commandOption.options[i] = reference.Option
				}
				optionStructsToGenerate = append(optionStructsToGenerate, commandOption)
				registeredStructs.Add(commandOption.name)
			}
		}
		return
	})
	if err != nil {
		return
	}
	err = out.WriteLn(`}`)
	if err != nil {
		return
	}

	for _, subStruct := range optionStructsToGenerate {
		err = out.WriteLn("")
		if err != nil {
			return
		}
		err = out.WriteLnF(`type %sOptions struct {`, subStruct.name)
		if err != nil {
			return
		}
		if len(subStruct.parentName) != 0 {
			err = out.WriteLnF(`%s%s %sOptions`, singleIndent, subStruct.parentName, subStruct.parentName)
			if err != nil {
				return
			}
		}
		for _, optionDef := range subStruct.options {
			err = writeOptionStructField(out, optionDef.Name, optionDef, g.optionTypes)
			if err != nil {
				return
			}
		}
		err = out.WriteLn(`}`)
		if err != nil {
			return
		}
	}

	return
}

func (g GoLang) packageName() string {
	if len(g.packageNameValue) == 0 {
		return defaultPackageName
	}
	return g.packageNameValue
}

func (g GoLang) interfaceName() string {
	if len(g.interfaceNameValue) == 0 {
		return defaultInterfaceName
	}
	return g.interfaceNameValue
}

func (g GoLang) structName() string {
	if len(g.structNameValue) == 0 {
		return defaultStructName
	}
	return g.structNameValue
}

func (g GoLang) globalOptionStructName() string {
	if len(g.globalOptionStructNameValue) == 0 {
		return defaultGlobalOptionStructName
	}
	return g.globalOptionStructNameValue
}

func collectImports(document *dsl.Document, optionTypes optionTypeRegistry) (out importRegistryType, err error) {
	out = make(importRegistryType)

	for _, option := range document.Options {
		err = addOptionToImports(out, []string{}, &option, optionTypes)
		if err != nil {
			return
		}
	}

	err = walkCommands(document, func(prefix []string, cmd dsl.Command) (err error) {
		for _, option := range cmd.Options {
			err = addOptionToImports(out, []string{}, &option, optionTypes)
			if err != nil {
				return
			}
		}
		return
	})
	return
}

func addOptionToImports(out importRegistryType, prefix []string, option *dsl.OptionOrReference, optionTypes optionTypeRegistry) (err error) {
	if t, ok := optionTypes[option.Type]; !ok {
		err = fmt.Errorf(`unsupported option type: "%s"`, option.Type)
		return
	} else {
		var imp goImport
		var useOptional bool
		useOptional, err = shouldUseOptional(option.Option, prefix)
		if err != nil {
			return
		}
		if useOptional {
			imp = t.ImportOptional
		} else {
			imp = t.Import
		}
		if !imp.Empty() {
			out[imp.Path] = imp
		}
	}
	return
}

func shouldUseOptional(option dsl.Option, prefix []string) (useOptional bool, err error) {
	if option.Required {
		if option.Default.IsSet() {
			path := "\"" + strings.Join(prefix, "/") + "/" + option.Name + "\""
			err = errors.New("option at " + path + " cannot have a default value and also be required")
			return
		}
		// optional not required
	} else {
		if option.Default.IsSet() {
			// optional is not required as there will always be a value
		} else {
			// optional required as there may not be a value
			useOptional = true
		}
	}
	return
}

func walkCommands(doc *dsl.Document, callback func(prefix []string, cmd dsl.Command) error) error {
	prefixes := make([]string, 0, likelyCommandMaxNestingDepth)
	return walkCommandsRecursive(doc.Commands, &prefixes, callback)
}

func walkCommandsRecursive(cmd dsl.NamedCommands, recurPrefix *[]string, callback func(prefix []string, cmd dsl.Command) error) (err error) {
	*recurPrefix = append(*recurPrefix, "")
	keys := make([]string, 0, len(cmd))
	for key := range cmd {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, name := range keys {
		command := cmd[name]
		(*recurPrefix)[len(*recurPrefix)-1] = name
		err = callback(*recurPrefix, command)
		if err != nil {
			return
		}
		if len(command.Commands) != 0 {
			err = walkCommandsRecursive(command.Commands, recurPrefix, callback)
			if err != nil {
				return
			}
		}
	}
	*recurPrefix = (*recurPrefix)[0 : len(*recurPrefix)-1]
	return
}

func joinPrefixesAsMethodName(prefix []string) string {
	titlized := make([]string, len(prefix))
	for i, s := range prefix {
		titlized[i] = strings.Title(s)
	}
	return strings.Join(titlized, "")
}

func writeOptionStructField(out *string_writer.Type, optionKey string, optionDef dsl.Option, optionTypes optionTypeRegistry) (err error) {
	useOptional := false
	useOptional, _ = shouldUseOptional(optionDef, []string{})

	if t, ok := optionTypes[optionDef.Type]; !ok {
		err = errors.New(`unsupported option type: "` + optionDef.Type + `"`)
	} else {
		typeToUse := ""
		if !useOptional {
			typeToUse = t.Type
		} else {
			typeToUse = t.OptionalType
		}
		err = out.WriteLnF(`%s%s %s`, singleIndent, optionDef.Name, typeToUse)
	}
	return
}

func prefixToOptionStructName(prefix []string) string {
	return strings.Join(eachString(prefix, func(s string) string {
		return strings.Title(s)
	}), "")
}

func eachString(items []string, callback func(string) string) (out []string) {
	out = make([]string, len(items))
	for i, item := range items {
		out[i] = callback(item)
	}
	return
}

func getStructOrBlank(prefix []string, set *string_set.Collection, globalStruct *golangSubStruct) string {
	for i := len(prefix) - 1; i > 0; i-- {
		p := prefixToOptionStructName(prefix[0:i])
		if set.Exists(p) {
			return p
		}
	}
	if globalStruct != nil {
		return globalStruct.name
	}
	return ""
}

func stringSlicePop(s []string) []string {
	if len(s) > 1 {
		return s[0 : len(s)-2]
	}
	return s[0:0]
}
