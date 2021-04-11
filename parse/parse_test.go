package parse

import "testing"

func TestEnvFlagParser_Parse(t *testing.T) {
	cases := map[string]struct {
		input EnvFlagParser
	}{}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			t.Fail()
			_ = c
		})
	}
}
