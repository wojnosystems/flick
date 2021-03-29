package string_set

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestCollection_Add(t *testing.T) {
	cases := map[string]struct {
		input    []string
		expected []string
	}{
		"empty": {
			input:    []string{},
			expected: []string{},
		},
		"one": {
			input:    []string{"a"},
			expected: []string{"a"},
		},
		"two": {
			input:    []string{"a", "b"},
			expected: []string{"a", "b"},
		},
		"two with duplicate": {
			input:    []string{"a", "b", "a"},
			expected: []string{"a", "b"},
		},
		"three with duplicates": {
			input:    []string{"a", "a", "a"},
			expected: []string{"a"},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			working := Collection{}
			for _, s := range c.input {
				working.Add(s)
			}
			actual := working.ToSlice()
			sort.Strings(actual)
			assert.Equal(t, c.expected, actual)

			for _, s := range c.expected {
				assert.True(t, working.Exists(s))
			}

			assert.False(t, working.Exists("notInSet"))
		})
	}
}
