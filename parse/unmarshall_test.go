package parse

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/wojnosystems/go-optional/v2"
	"testing"
	"time"
)

func TestUnmarshall(t *testing.T) {
	cases := map[string]struct {
		input    func(configuration interface{}) (err error)
		expected appConfig
	}{
		"empty": {
			input: func(configuration interface{}) (err error) {
				return Unmarshall(configuration)
			},
		},
		"optional file missing": {
			input: func(configuration interface{}) (err error) {
				return Unmarshall(configuration,
					FileIsOptional(optional.StringFrom("/tmp/non-existant"), Yaml()))
			},
		},
		"yaml parsing": {
			input: func(configuration interface{}) (err error) {
				return Unmarshall(configuration,
					newFileAsBytes([]byte(`---
hostname: test.example.com
delay: 30s
`), Yaml()))
			},
			expected: appConfig{
				Hostname: optional.StringFrom("test.example.com"),
				Delay:    optional.DurationFrom(30 * time.Second),
			},
		},
		"yaml, env, flag parsing": {
			input: func(configuration interface{}) (err error) {
				return Unmarshall(configuration,
					newFileAsBytes([]byte(`---
hostname: test.example.com
delay: 30s
`), Yaml()),
				)
			},
			expected: appConfig{
				Hostname: optional.StringFrom("test.example.com"),
				Delay:    optional.DurationFrom(30 * time.Second),
			},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			var actual appConfig
			err := c.input(&actual)
			assert.NoError(t, err)
			assert.Equal(t, c.expected, actual)
		})
	}
}

type fileAsBytes struct {
	content      []byte
	unmarshaller FileUnmarshaler
}

func newFileAsBytes(content []byte, unmarshaller FileUnmarshaler) Unmarshaler {
	return &fileAsBytes{
		content:      content,
		unmarshaller: unmarshaller,
	}
}

func (f *fileAsBytes) Unmarshal(config interface{}) (err error) {
	rdr := bytes.NewReader(f.content)
	return f.unmarshaller.UnmarshalFile(rdr, config)
}
