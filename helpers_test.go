package goseeder

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFindString(t *testing.T) {
	input := []string{
		"dandy",
		"trout",
		"fish",
		"more",
		"fish",
		"ok",
	}

	pos, found := findString(input, "fish")
	require.Equal(t, 2, pos)
	require.Equal(t, true, found)
}

func TestPrepareStatement(t *testing.T) {
	table := "categories"
	data := map[string]interface{}{
		"id":   "100",
		"name": "common",
	}
	sb, args := prepareStatement(table, data)

	require.Equal(
		t,
		"insert into categories (id, name) values (?, ?)",
		sb.String(),
	)

	require.Equal(
		t,
		[]interface{}{
			"100",
			"common"},
		args,
	)
}

var testCases = []struct {
	name     string
	value    interface{}
	expected interface{}
}{
	{"bool_string", "true", "true"},
	{"bool_true", true, true},
	{"bool_string", "false", "false"},
	{"bool_false", false, false},
	{"int_string", "12", "12"},
	{"int", 12, int(12)},
	{"float_string", "12.77", "12.77"},
	{"float", 12.77, 12.77},
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
