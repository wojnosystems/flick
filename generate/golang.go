package generate

import (
	"context"
	"errors"
	"github.com/wojnosystems/flick/generate/dsl"
	sorted_set "github.com/wojnosystems/go-sorted-set"
	"io"
	"strings"
)

type GoLang struct {
	packageNameValue            string
	interfaceNameValue          string
	structNameValue             string
	globalOptionStructNameValue string
}

const (
	likelyCommandMaxNestingDepth = 10
	singleIndent                 = `  `

	defaultPackageName            = "flickstub"
	defaultInterfaceName          = "Interface"
	defaultStructName             = "Unimplemented"
	defaultGlobalOptionStructName = "AllCommandOptions"
)

func (g *GoLang) Generate(ctx context.Context, document *dsl.Document, output io.Writer) (bytesWritten int, err error) {
	out := stringWriter{
		stream: output,
	}
	defer func() {
		bytesWritten = out.counter.Count()
	}()
	err = out.write2Ln("package " + g.packageName())
	if err != nil {
		return
	}

	err = out.writeLn("import (")
	if err != nil {
		return
	}
	imports := sorted_set.NewString(
		"context",
		"errors",
		// TODO: how do we know we'll need go-optional?
		"github.com/wojnosystems/go-optional/v2").
		Sort()
	for _, value := range imports {
		err = out.writeLnF(`%s"%s"`, singleIndent, value)
		if err != nil {
			return
		}
	}
	err = out.write2Ln(")")
	if err != nil {
		return
	}

	err = out.writeLnF(`type %s interface {`, g.interfaceName())
	if err != nil {
		return
	}
	if document.Options.HasAny() {
		err = out.writeLnF(`%sHookBefore(ctx context.Context, opts *%s) error`, singleIndent, g.globalOptionStructName())
		if err != nil {
			return
		}
		err = out.writeLnF(`%sHookAfter(ctx context.Context, opts *%s, err error) error`, singleIndent, g.globalOptionStructName())
		if err != nil {
			return
		}
	} else {
		err = out.writeLnF(`%sHookBefore(ctx context.Context) error`, singleIndent)
		if err != nil {
			return
		}
		err = out.writeLnF(`%sHookAfter(ctx context.Context, err error) error`, singleIndent)
		if err != nil {
			return
		}
	}
	err = walkCommands(document, func(prefix []string, cmd dsl.Command) (walkErr error) {
		methodName := joinPrefixesAsMethodName(prefix)
		if cmd.Commands.HasAny() {
			err = out.writeLnF(`%s%sHookBefore(ctx context.Context, opts *%sOptions) error`, singleIndent, methodName, methodName)
			if err != nil {
				return
			}
			err = out.writeLnF(`%s%sHookAfter(ctx context.Context, opts *%sOptions, err error) error`, singleIndent, methodName, methodName)
			if err != nil {
				return
			}
		} else {
			walkErr = out.writeLnF(`%s%s(ctx context.Context, opts *%sOptions) error`, singleIndent, methodName, methodName)
			if walkErr != nil {
				return
			}
		}
		return
	})
	if err != nil {
		return
	}
	err = out.writeLn(`}`)
	if err != nil {
		return
	}

	if document.Options.HasAny() {
		err = out.writeLn("")
		if err != nil {
			return
		}
		err = out.writeLnF(`type %s struct {`, g.globalOptionStructName())
		if err != nil {
			return
		}
		for optionKey, optionDef := range document.Options {
			err = writeOptionStructField(out, []string{}, optionKey, optionDef.Option)
			if err != nil {
				return
			}
		}
		err = out.writeLn(`}`)
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

func walkCommands(doc *dsl.Document, callback func(prefix []string, cmd dsl.Command) error) error {
	prefixes := make([]string, 0, likelyCommandMaxNestingDepth)
	return walkCommandsRecursive(doc.Commands, &prefixes, callback)
}

func walkCommandsRecursive(cmd dsl.NamedCommands, recurPrefix *[]string, callback func(prefix []string, cmd dsl.Command) error) (err error) {
	*recurPrefix = append(*recurPrefix, "")
	for name, command := range cmd {
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

func writeOptionStructField(out stringWriter, prefixes []string, optionKey string, optionDef dsl.Option) (err error) {
	useOptional := false
	if optionDef.Required {
		if optionDef.Default.IsSet() {
			path := "\"" + strings.Join(prefixes, "/") + "/" + optionKey + "\""
			err = errors.New("option at " + path + " cannot have a default value and also be required")
			return
		}
		// optional not required
	} else {
		if optionDef.Default.IsSet() {
			// optional is not required as there will always be a value
		} else {
			// optional required as there may not be a value
			useOptional = true
		}
	}

	// TODO: replace this with a registry
	typeDefinition := ""
	switch optionDef.Type {
	// types supported by optional:
	case "int":
		if useOptional {
			typeDefinition = "optional.Int"
		} else {
			typeDefinition = "int"
		}
	default:
		err = errors.New(`unsupported option type: "` + optionDef.Type + `"`)
		return
	}

	err = out.writeLnF(`%s%s %s`, singleIndent, optionKey, typeDefinition)
	return
}
