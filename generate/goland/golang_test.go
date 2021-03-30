package goland

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
  "github.com/wojnosystems/flick/cli"
`

	cases := map[string]struct {
		input    dsl.Document
		expected string
	}{
		"empty": {
			input: dsl.Document{},
			expected: globalHeader + `)

type Interface interface {
  HookBefore(ctx context.Context) error
  HookAfter(ctx context.Context, err error) error
}

type Unimplemented struct {
  HookBefore(_ context.Context) error {
    return nil
  }
  HookAfter(_ context.Context, _ error) error {
    return nil
  }
}
`,
		},
		"with global options": {
			input: dsl.Document{
				Options: []dsl.OptionOrReference{
					dsl.OptionOrReference{
						Option: dsl.Option{
							Name: "key1",
							Type: "int",
						},
					},
				},
			},
			expected: globalHeader + `  "github.com/wojnosystems/go-optional/v2"
)

type Interface interface {
  HookBefore(ctx context.Context, opts *AllCommandOptions) error
  HookAfter(ctx context.Context, opts *AllCommandOptions, err error) error
}

type AllCommandOptions struct {
  key1 optional.Int
}

type Unimplemented struct {
  HookBefore(_ context.Context, _ *AllCommandOptions) error {
    return nil
  }
  HookAfter(_ context.Context, _ *AllCommandOptions, _ error) error {
    return nil
  }
}
`,
		},
		"two commands without options": {
			input: dsl.Document{
				Commands: dsl.NamedCommands{
					"bar": dsl.Command{},
					"foo": dsl.Command{},
				},
			},
			expected: globalHeader + `)

type Interface interface {
  HookBefore(ctx context.Context) error
  HookAfter(ctx context.Context, err error) error
  Bar(ctx context.Context) error
  Foo(ctx context.Context) error
}

type Unimplemented struct {
  HookBefore(_ context.Context) error {
    return nil
  }
  HookAfter(_ context.Context, _ error) error {
    return nil
  }
  Bar(_ context.Context) error {
    return cli.ErrUnimplemented
  }
  Foo(_ context.Context) error {
    return cli.ErrUnimplemented
  }
}
`,
		},
		"two commands with nested options": {
			input: dsl.Document{
				Commands: dsl.NamedCommands{
					"bar": dsl.Command{
						Options: []dsl.OptionOrReference{
							{
								Option: dsl.Option{
									Name: "puppy",
									Type: "int",
								},
							},
						},
					},
					"foo": dsl.Command{
						Options: []dsl.OptionOrReference{
							{
								Option: dsl.Option{
									Name: "cat",
									Type: "int",
								},
							},
						},
					},
				},
			},
			expected: globalHeader + `  "github.com/wojnosystems/go-optional/v2"
)

type Interface interface {
  HookBefore(ctx context.Context) error
  HookAfter(ctx context.Context, err error) error
  Bar(ctx context.Context, opts *BarOptions) error
  Foo(ctx context.Context, opts *FooOptions) error
}

type BarOptions struct {
  puppy optional.Int
}

type FooOptions struct {
  cat optional.Int
}

type Unimplemented struct {
  HookBefore(_ context.Context) error {
    return nil
  }
  HookAfter(_ context.Context, _ error) error {
    return nil
  }
  Bar(_ context.Context, _ *BarOptions) error {
    return cli.ErrUnimplemented
  }
  Foo(_ context.Context, _ *FooOptions) error {
    return cli.ErrUnimplemented
  }
}
`,
		},
		"sub-command inherit global option": {
			input: dsl.Document{
				Options: []dsl.OptionOrReference{
					{
						Option: dsl.Option{
							Name: "puppy",
							Type: "int",
						},
					},
				},
				Commands: dsl.NamedCommands{
					"bar": dsl.Command{},
				},
			},
			expected: globalHeader + `  "github.com/wojnosystems/go-optional/v2"
)

type Interface interface {
  HookBefore(ctx context.Context, opts *AllCommandOptions) error
  HookAfter(ctx context.Context, opts *AllCommandOptions, err error) error
  Bar(ctx context.Context, opts *AllCommandOptions) error
}

type AllCommandOptions struct {
  puppy optional.Int
}

type Unimplemented struct {
  HookBefore(_ context.Context, _ *AllCommandOptions) error {
    return nil
  }
  HookAfter(_ context.Context, _ *AllCommandOptions, _ error) error {
    return nil
  }
  Bar(_ context.Context, _ *AllCommandOptions) error {
    return cli.ErrUnimplemented
  }
}
`,
		},
		"sub-command with nested global option": {
			input: dsl.Document{
				Options: []dsl.OptionOrReference{
					{
						Option: dsl.Option{
							Name: "puppy",
							Type: "int",
						},
					},
				},
				Commands: dsl.NamedCommands{
					"bar": dsl.Command{
						Options: []dsl.OptionOrReference{
							{
								Option: dsl.Option{
									Name: "barOption",
									Type: "duration",
								},
							},
						},
					},
				},
			},
			expected: globalHeader + `  "github.com/wojnosystems/go-optional/v2"
)

type Interface interface {
  HookBefore(ctx context.Context, opts *AllCommandOptions) error
  HookAfter(ctx context.Context, opts *AllCommandOptions, err error) error
  Bar(ctx context.Context, opts *BarOptions) error
}

type AllCommandOptions struct {
  puppy optional.Int
}

type BarOptions struct {
  AllCommand AllCommandOptions
  barOption optional.Duration
}

type Unimplemented struct {
  HookBefore(_ context.Context, _ *AllCommandOptions) error {
    return nil
  }
  HookAfter(_ context.Context, _ *AllCommandOptions, _ error) error {
    return nil
  }
  Bar(_ context.Context, _ *BarOptions) error {
    return cli.ErrUnimplemented
  }
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
