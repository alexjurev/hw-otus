package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading directory %s : %w", dir, err)
	}

	env := newEnvironment()

	for _, file := range files {
		if !file.Type().IsRegular() || file.IsDir() {
			continue
		}

		fileName := file.Name()

		if strings.Contains(fileName, "=") {
			continue
		}

		info, err := file.Info()
		if err != nil {
			return nil, err
		}

		if info.Size() == 0 {
			env[fileName] = EnvValue{
				Value:      "",
				NeedRemove: true,
			}
			continue
		}

		val, err := getFirstLine(path.Join(dir, fileName))
		if err != nil {
			return nil, err
		}

		val = sanitizeVal(val)

		env[fileName] = EnvValue{Value: val}
	}

	return env, nil
}

func newEnvironment() Environment {
	e := make(Environment)
	return e
}

func getFirstLine(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("read file %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	return scanner.Text(), scanner.Err()
}

func sanitizeVal(s string) string {
	s = strings.TrimSuffix(s, "\t ")
	s = string(bytes.ReplaceAll([]byte(s), []byte{0x00}, []byte("\n")))

	if strings.TrimSpace(s) == "" {
		return ""
	}

	return s
}
