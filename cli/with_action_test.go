package cli

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli/v2"
	"github.com/wojnosystems/flick/builtin_actions"
	flag_unmarshaler "github.com/wojnosystems/go-flag-unmarshaler"
	"github.com/wojnosystems/okey-dokey/bad"
	"testing"
)

type withActionTestMock struct {
	mock.Mock
	Name string
}

func (m *withActionTestMock) Validate(emitter bad.MemberEmitter) (err error) {
	a := m.Called(emitter)
	return a.Error(0)
}

func TestWithAction_Run(t *testing.T) {
	cases := map[string]struct {
		app         WithAction
		flags       []flag_unmarshaler.Group
		expected    interface{}
		expectedErr string
	}{
		"no parsing": {
			app: WithAction{
				Action: builtin_actions.CommandPrintVersion("v1.0.0"),
			},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			m := &withActionTestMock{}
			m.On("Validate", mock.Anything).
				Once().
				Return(nil)
			c.app.Options = m
			err := c.app.Run(c.flags)
			assert.EqualError(t, err, c.expectedErr)
			assert.Equal(t, c.expected, c.app.Options)
			m.AssertExpectations(t)
		})
	}
}
