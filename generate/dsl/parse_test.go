package dsl

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wojnosystems/go-optional/v2"
	"github.com/wojnosystems/okey-dokey/bad"
	"testing"
)

func TestParse(t *testing.T) {
	cases := map[string]struct {
		input    string
		expected Document
	}{
		"empty": {
			input:    "",
			expected: Document{},
		},
		"version": {
			input: `
optionapi:
  version: 4
`,
			expected: Document{
				OptionApi: OptionApi{Version: optional.StringFrom("4")},
			},
		},
		"command with option": {
			input: `
commands:
  server:
    options:
      - $ref: "#/components/options/ConnectTimeout"
    minArgs: 2
    maxArgs: 3
components:
  options:
    ConnectTimeout:
      type: duration
      description: "how long to wait when connecting to the server"
      usage: "Ns"
      env: 
        name: "CONNECT_TIMEOUT"
      flag:
        name: "connectTimeout"
        aliases: ["c"]
      default: 30s
`,
			expected: Document{
				Commands: NamedCommands{
					"server": Command{
						Options: []OptionOrReference{
							{
								Option: Option{
									Type:        "duration",
									Description: optional.StringFrom("how long to wait when connecting to the server"),
									Usage:       optional.StringFrom("Ns"),
									Env: EnvDef{
										Name: "CONNECT_TIMEOUT",
									},
									Flag: FlagDef{
										Name:    "connectTimeout",
										Aliases: []string{"c"},
									},
									Default: optional.StringFrom("30s"),
								},
							},
						},
						MinArgs: 2,
						MaxArgs: 3,
					},
				},
				Components: Components{
					Options: NamedOptions{
						"ConnectTimeout": {
							Type:        "duration",
							Description: optional.StringFrom("how long to wait when connecting to the server"),
							Usage:       optional.StringFrom("Ns"),
							Env: EnvDef{
								Name: "CONNECT_TIMEOUT",
							},
							Flag: FlagDef{
								Name:    "connectTimeout",
								Aliases: []string{"c"},
							},
							Default: optional.StringFrom("30s"),
						},
					},
				},
			},
		},
		"command with subcommands": {
			input: `
commands:
  server:
    commands:
      start:
        usage: "start"
      stop:
        usage: "stop"
      restart:
        usage: "restart"
`,
			expected: Document{
				Commands: NamedCommands{
					"server": Command{
						Commands: NamedCommands{
							"start": Command{
								Usage: optional.StringFrom("start"),
							},
							"stop": Command{
								Usage: optional.StringFrom("stop"),
							},
							"restart": Command{
								Usage: optional.StringFrom("restart"),
							},
						},
					},
				},
			},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			input := bytes.NewReader([]byte(c.input))
			actual, err := Parse(input, bad.NewCollection())
			require.NoError(t, err)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestParseWithValidationErrors(t *testing.T) {
	cases := map[string]struct {
		input       string
		expected    bad.ReceiveCollector
		expectedErr error
	}{
		"within alloptions it requires minArgs is less than maxArgs": {
			input: `---
minArgs: 3
maxArgs: 1
`,
			expected: func() (c bad.ReceiveCollector) {
				c = bad.NewCollection()
				c.Emit("minArgs must be less than maxArgs")
				return
			}(),
			expectedErr: ErrValidation,
		},
		"within sub-command it requires minArgs is less than maxArgs": {
			input: `---
commands:
  server:
    minArgs: 3
    maxArgs: 1
`,
			expected: func() (c bad.ReceiveCollector) {
				c = bad.NewCollection()
				c.Into("server").Emit("minArgs must be less than maxArgs")
				return
			}(),
			expectedErr: ErrValidation,
		},
		"with sub-commands maxArgs must be 0": {
			input: `---
maxArgs: 2
commands:
  server:
`,
			expected: func() (c bad.ReceiveCollector) {
				c = bad.NewCollection()
				c.Emit("when sub-commands are specified, maxArgs must be 0")
				return
			}(),
			expectedErr: ErrValidation,
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			input := bytes.NewReader([]byte(c.input))
			validationErrors := bad.NewCollection()
			_, err := Parse(input, validationErrors)
			assertEqualNilSafeError(t, err, c.expectedErr)
			assert.Equal(t, c.expected, validationErrors)
		})
	}
}

func assertEqualNilSafeError(t *testing.T, actual error, expectedOrNil error) {
	if expectedOrNil != nil {
		assert.EqualError(t, actual, expectedOrNil.Error())
	} else {
		require.NoError(t, actual)
	}
}
