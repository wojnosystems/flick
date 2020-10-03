package source

import (
	"github.com/stretchr/testify/assert"
	"github.com/wojnosystems/flick/pkg/optional_set_value"
	"github.com/wojnosystems/flick/pkg/set_value"
	"github.com/wojnosystems/go-optional"
	"testing"
)

func TestEnv_Unmarshall(t *testing.T) {
	cases := map[string]struct {
		env      *envMock
		expected appConfigMock
	}{
		"nothing": {
			env: &envMock{},
		},
		"name": {
			env: &envMock{
				mock: map[string]string{
					"Name": "SuperServer",
				},
			},
			expected: appConfigMock{
				Name: optional.StringFrom("SuperServer"),
			},
		},
		"db[1].Host": {
			env: &envMock{
				mock: map[string]string{
					"Databases_1_Host": "example.com",
				},
			},
			expected: appConfigMock{
				Databases: []dbConfigMock{
					{},
					{
						Host: optional.StringFrom("example.com"),
					},
				},
			},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := &appConfigMock{}
			e := Env{
				envs:          c.env,
				parseRegistry: optional_set_value.Register(set_value.RegisterGoPrimitives(&set_value.Registry{})),
			}
			err := e.Unmarshall(actual)
			assert.NoError(t, err)
			assert.True(t, c.expected.IsEqual(actual))
		})
	}
}
