package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	env, err := ReadDir("./testdata/env")
	require.NoError(t, err)

	retCode := RunCmd([]string{"bad_command"}, env)
	require.Equal(t, -1, retCode)
}
