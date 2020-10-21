package goseeder

import (
	"github.com/stretchr/testify/require"
	"testing"
)

var testCases = []struct {
	name     string
	value    string
	expected interface{}
}{
	{"bool", "true", true},
	{"bool", "false", false},
	{"int", "12", int64(12)},
	{"float", "12.77", 12.770000457763672},
	{"string", "justastring", "justastring"},
	{"string_with_nr", "12justastring", "12justastring"},
}

func TestParseValue(t *testing.T) {
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, parseValue(tt.value))
		})
	}
}
