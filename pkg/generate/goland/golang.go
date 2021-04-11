package goland

import (
	"context"
	"errors"
	"fmt"
	"github.com/wojnosystems/flick/pkg/generate/dsl"
	"github.com/wojnosystems/flick/pkg/string_writer"
	"github.com/wojnosystems/go-string-set/string_set"
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
	goFlickLibraryImportPath    = "github.com/wojnosystems/flick/cli"
)

type optionStruct struct {
	name       string
	parentName string
	options    []dsl.Option
}

type structMethodDefinition struct {
	declaration string
	body        string
}

type collected struct {
	interfaceDeclarations []string
	baseStructMethodDefs  []structMethodDefinition
	globalStruct          *optionStruct
	optionStructs         []optionStruct
	optionStructRegistry  string_set.Interface
}

func (c *collected) addGlobalStruct(o optionStruct) {
	c.globalStruct = &o
	c.addOptionStruct(o)
}

func (c *collected) addOptionStruct(o optionStruct) {
	c.optionStructs = append(c.optionStructs, o)
	c.optionStructRegistry.Add(o.name)
}

func (c collected) getParentStructName(prefix []string) string {
	return getStructOrBlank(stringSlicePop(prefix), c.optionStructRegistry, c.globalStruct)
}

func (g *GoLang) Generate(_ context.Context, document *dsl.Document, output io.Writer) (bytesWritten int, err error) {
	if g.optionTypes == nil {
		g.optionTypes = make(optionTypeRegistry)
		registerOptionalTypes(g.optionTypes)
	}

	out := string_writer.New(
		output,
		singleIndent)
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
	err = writeImports(out, imports)
	if err != nil {
		return
	}

	generatedComponents := collected{
		interfaceDeclarations: make([]string, 0, 10),
		baseStructMethodDefs:  make([]structMethodDefinition, 0, 10),
		optionStructs:         make([]optionStruct, 0, 10),
		optionStructRegistry:  string_set.New(),
	}

	err = g.collectComponents(document, &generatedComponents)
	if err != nil {
		return
	}

	err = g.writeInterface(out, generatedComponents.interfaceDeclarations)
	if err != nil {
		return
	}

	err = g.writeOptionStructs(out, &generatedComponents)
	if err != nil {
		return
	}

	err = g.writeUnimplementedStruct(out, generatedComponents.baseStructMethodDefs)

	// TODO: write registration function that maps the parsed command with the options group to execute what is required.

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
	if err != nil {
		return
	}

	out["context"] = goImport{
		Path:  "context",
		Alias: "",
	}
	out[goFlickLibraryImportPath] = goImport{
		Path:  goFlickLibraryImportPath,
		Alias: "",
	}
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

func writeImports(out *string_writer.Type, imports importRegistryType) (err error) {
	importKeysSorted := make([]string, 0, len(imports))
	for s := range imports {
		importKeysSorted = append(importKeysSorted, s)
	}
	sort.Strings(importKeysSorted)
	err = out.WriteLn("import (")
	if err != nil {
		return
	}
	err = out.In(func(out *string_writer.Type) (err error) {
		for _, key := range importKeysSorted {
			record := `"` + imports[key].Path + `"`
			if len(imports[key].Alias) != 0 {
				record = imports[key].Alias + " " + record
			}
			err = out.WriteLn(record)
			if err != nil {
				return
			}
		}
		return
	})
	if err != nil {
		return
	}
	return out.Write2Ln(")")
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

func (g *GoLang) collectComponents(document *dsl.Document, c *collected) (err error) {
	if len(document.Options) != 0 {
		// BEFORE HOOK
		c.interfaceDeclarations = append(c.interfaceDeclarations,
			fmt.Sprintf(`HookBefore(ctx context.Context, opts *%sOptions) error`,
				g.globalOptionStructName()))
		c.baseStructMethodDefs = append(c.baseStructMethodDefs, structMethodDefinition{
			declaration: fmt.Sprintf(`HookBefore(_ context.Context, _ *%sOptions) error`, g.globalOptionStructName()),
			body:        "return nil",
		})

		// AFTER HOOK
		c.interfaceDeclarations = append(c.interfaceDeclarations,
			fmt.Sprintf(`HookAfter(ctx context.Context, opts *%sOptions, err error) error`,
				g.globalOptionStructName()))
		c.baseStructMethodDefs = append(c.baseStructMethodDefs, structMethodDefinition{
			declaration: fmt.Sprintf(`HookAfter(_ context.Context, _ *%sOptions, _ error) error`, g.globalOptionStructName()),
			body:        "return nil",
		})
	} else {
		c.interfaceDeclarations = append(c.interfaceDeclarations,
			`HookBefore(ctx context.Context) error`)
		c.baseStructMethodDefs = append(c.baseStructMethodDefs, structMethodDefinition{
			declaration: `HookBefore(_ context.Context) error`,
			body:        "return nil",
		})

		c.interfaceDeclarations = append(c.interfaceDeclarations,
			`HookAfter(ctx context.Context, err error) error`)
		c.baseStructMethodDefs = append(c.baseStructMethodDefs, structMethodDefinition{
			declaration: `HookAfter(_ context.Context, _ error) error`,
			body:        "return nil",
		})
	}

	if len(document.Options) != 0 {
		globalStruct := optionStruct{
			name:    g.globalOptionStructName(),
			options: make([]dsl.Option, len(document.Options)),
		}
		for i, reference := range document.Options {
			globalStruct.options[i] = reference.Option
		}
		c.addGlobalStruct(globalStruct)
	}

	err = walkCommands(document, func(prefix []string, cmd dsl.Command) (walkErr error) {
		methodName := joinPrefixesAsMethodName(prefix)

		var optionStructName string
		if len(cmd.Options) != 0 {
			optionStructName = methodName
		} else {
			optionStructName = getStructOrBlank(stringSlicePop(prefix), c.optionStructRegistry, c.globalStruct)
		}

		optionFormal := ""
		optionFormalWithoutNamedParam := ""
		if optionStructName != "" {
			optionFormal = fmt.Sprintf(", opts *%sOptions", optionStructName)
			optionFormalWithoutNamedParam = fmt.Sprintf(", _ *%sOptions", optionStructName)
		}

		if cmd.Commands.HasAny() {
			c.interfaceDeclarations = append(c.interfaceDeclarations, fmt.Sprintf(`%sHookBefore(ctx context.Context%s) error`, methodName, optionFormal))
			c.baseStructMethodDefs = append(c.baseStructMethodDefs, structMethodDefinition{
				declaration: fmt.Sprintf(`%sHookBefore(_ context.Context%s) error`, methodName, optionFormalWithoutNamedParam),
				body:        `return nil`,
			})

			c.interfaceDeclarations = append(c.interfaceDeclarations, fmt.Sprintf(`%sHookAfter(ctx context.Context%s, err error) error`, methodName, optionFormal))
			c.baseStructMethodDefs = append(c.baseStructMethodDefs, structMethodDefinition{
				declaration: fmt.Sprintf(`%sHookAfter(_ context.Context%s, _ error) error`, methodName, optionFormalWithoutNamedParam),
				body:        `return nil`,
			})
		} else {
			c.interfaceDeclarations = append(c.interfaceDeclarations, fmt.Sprintf(
				`%s(ctx context.Context%s) error`, methodName, optionFormal))
			c.baseStructMethodDefs = append(c.baseStructMethodDefs, structMethodDefinition{
				declaration: fmt.Sprintf(`%s(_ context.Context%s) error`, methodName, optionFormalWithoutNamedParam),
				body:        `return cli.ErrCommandUnimplemented`,
			})

			if len(cmd.Options) != 0 {
				commandOption := optionStruct{
					name:       prefixToOptionStructName(prefix),
					parentName: c.getParentStructName(prefix),
					options:    make([]dsl.Option, len(cmd.Options)),
				}
				for i, reference := range cmd.Options {
					commandOption.options[i] = reference.Option
				}
				c.addOptionStruct(commandOption)
			}
		}
		return
	})
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

func (g *GoLang) writeInterface(out *string_writer.Type, declarations []string) (err error) {
	err = out.WriteLnF(`type %s interface {`, g.interfaceName())
	if err != nil {
		return
	}
	err = out.In(func(out *string_writer.Type) (err error) {
		for _, declaration := range declarations {
			err = out.WriteLn(declaration)
		}
		return
	})
	if err != nil {
		return
	}
	err = out.WriteLn(`}`)
	return
}

func (g *GoLang) writeOptionStructs(out *string_writer.Type, declarations *collected) (err error) {
	for _, subStruct := range declarations.optionStructs {
		err = out.WriteLn("")
		if err != nil {
			return
		}
		err = out.WriteLnF(`type %sOptions struct {`, subStruct.name)
		if err != nil {
			return
		}
		err = out.In(func(out *string_writer.Type) (err error) {
			if len(subStruct.parentName) != 0 {
				err = out.WriteLnF(`%s %sOptions`, subStruct.parentName, subStruct.parentName)
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
			return
		})
		if err != nil {
			return
		}
		err = out.WriteLn(`}`)
	}
	return
}

func (g *GoLang) writeUnimplementedStruct(out *string_writer.Type, declarations []structMethodDefinition) (err error) {
	err = out.WriteLn("")
	if err != nil {
		return
	}
	err = out.WriteLnF("type %s struct {", g.structName())
	if err != nil {
		return
	}
	err = out.In(func(out *string_writer.Type) (err error) {
		for _, s := range declarations {
			err = out.WriteLnF(`%s {`, s.declaration)
			if err != nil {
				return
			}
			err = out.In(func(out *string_writer.Type) error {
				return out.WriteLn(s.body)
			})
			if err != nil {
				return
			}
			err = out.WriteLn(`}`)
			if err != nil {
				return
			}
		}
		return
	})
	if err != nil {
		return
	}
	err = out.WriteLn("}")
	return nil
}

func joinPrefixesAsMethodName(prefix []string) string {
	title := make([]string, len(prefix))
	for i, s := range prefix {
		title[i] = strings.Title(s)
	}
	return strings.Join(title, "")
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
		err = out.WriteLnF(`%s %s`, optionDef.Name, typeToUse)
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

func getStructOrBlank(prefix []string, set string_set.Tester, globalStruct *optionStruct) string {
	for i := len(prefix) - 1; i > 0; i-- {
		p := prefixToOptionStructName(prefix[0:i])
		if set.Includes(p) {
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
