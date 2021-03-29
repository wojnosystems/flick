package generate

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wojnosystems/flick/generate/dsl"
	"testing"
)

func TestGoLang_Generate(t *testing.T) {
	globalHeader := `package flickstub

import (
  "context"
  "errors"
  "github.com/wojnosystems/go-optional/v2"
)
`

	cases := map[string]struct {
		input    dsl.Document
		expected string
	}{
		"empty": {
			input: dsl.Document{},
			expected: globalHeader + `
type Interface interface {
  HookBefore(ctx context.Context) error
  HookAfter(ctx context.Context, err error) error
}
`,
		},
		"with global options": {
			input: dsl.Document{
				Options: map[string]dsl.OptionOrReference{
					"key1": dsl.OptionOrReference{
						Option: dsl.Option{
							Type: "int",
						},
					},
				},
			},
			expected: globalHeader + `
type Interface interface {
  HookBefore(ctx context.Context, opts *AllCommandOptions) error
  HookAfter(ctx context.Context, opts *AllCommandOptions, err error) error
}

type AllCommandOptions struct {
  key1 optional.Int
}
`,
		},
		"two commands": {
			input: dsl.Document{
				Commands: dsl.NamedCommands{
					"bar": dsl.Command{},
					"foo": dsl.Command{},
				},
			},
			expected: globalHeader + `
type Interface interface {
  HookBefore(ctx context.Context) error
  HookAfter(ctx context.Context, err error) error
  Bar(ctx context.Context, opts *BarOptions) error
  Foo(ctx context.Context, opts *FooOptions) error
}
`,
		},
	}

	g := GoLang{}
	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := bytes.Buffer{}
			_, err := g.Generate(context.TODO(), &c.input, &actual)
			require.NoError(t, err)
			assert.Equal(t, c.expected, actual.String())
		})
	}
}
