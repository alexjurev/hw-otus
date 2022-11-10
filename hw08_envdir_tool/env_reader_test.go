package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	env, err := ReadDir("./testdata/env")

	require.NoError(t, err)
	expected := Environment(map[string]EnvValue{
		"BAR":   {Value: "bar"},
		"EMPTY": {Value: "", NeedRemove: false},
		"FOO":   {Value: "   foo\nwith new line"},
		"HELLO": {Value: `"hello"`},
		"UNSET": {Value: "", NeedRemove: true},
	})
	require.Equal(t, expected, env)
}
