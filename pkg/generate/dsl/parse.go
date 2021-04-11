package dsl

import (
	"errors"
	"github.com/wojnosystems/flick/parse"
	optional_parse_registry "github.com/wojnosystems/go-optional-parse-registry/v2"
	"github.com/wojnosystems/okey-dokey/bad"
	"io"
)

var ErrValidation = errors.New("failed to validate optionapi spec")

func Parse(r io.Reader, emitter bad.MemberEmitter) (out Document, err error) {
	registry := optional_parse_registry.RegisterFluent(optional_parse_registry.NewWithGoPrimitives())

	yamlParser := parse.YamlWithParseRegister(registry)
	err = yamlParser.UnmarshalFile(r, &out)
	if err != nil {
		return
	}

	tracked := newTrackedEmitter(emitter)
	DocumentValidations.Validate(&out, tracked)
	if tracked.isInvalid() {
		err = ErrValidation
		return
	}

	// populate Document Refs, walk the entire tree
	err = replaceDocumentReferences(&out)
	return
}

func replaceDocumentReferences(doc *Document) (err error) {
	refLookup := make(map[string]*Option)
	for optionName, opt := range doc.Components.Options {
		refLookup["#/components/options/"+optionName] = &opt
	}
	return replaceDocumentReferencesRecursive(doc.Commands, refLookup)
}

// replaceDocumentReferencesRecursive is inefficient, replace with stack-based one later
func replaceDocumentReferencesRecursive(namedCommand NamedCommands, lookup map[string]*Option) (err error) {
	for cmdName, cmd := range namedCommand {
		for dex, opt := range cmd.Options {
			if isBlank(opt.Reference) {
				continue
			}
			if ref, ok := lookup[opt.Reference]; !ok {
				return newErrReference(opt.Reference)
			} else {
				namedCommand[cmdName].Options[dex].Reference = ""
				namedCommand[cmdName].Options[dex].Option = *ref
			}
		}
		err = replaceDocumentReferencesRecursive(cmd.Commands, lookup)
		if err != nil {
			return
		}
	}
	return
}
