package main

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	const (
		testPath    = "test/"
		simpleFile  = "simple_file"
		emptyString = ""
		allPerm     = 0o775
		filePerm    = 0o660
	)

	file := "test/file"
	emptyFile := "test/emptyFile"
	targetFile1 := "test/targetFile1"
	targetFile2 := "test/targetFile2"
	targetFile3 := "test/targetFile3"
	targetFile4 := "test/targetFile4"
	targetFile5 := "test/targetFile5"
	targetFile6 := "test/targetFile6"

	fakeFs := afero.NewMemMapFs()

	fakeFs.MkdirAll(testPath, allPerm)

	afero.WriteFile(fakeFs, file, []byte(simpleFile), filePerm)

	afero.WriteFile(fakeFs, emptyFile, []byte(emptyString), filePerm)
	afero.WriteFile(fakeFs, targetFile1, []byte{}, filePerm)
	afero.WriteFile(fakeFs, targetFile2, []byte{}, filePerm)
	afero.WriteFile(fakeFs, targetFile3, []byte{}, filePerm)
	afero.WriteFile(fakeFs, targetFile4, []byte{}, filePerm)
	afero.WriteFile(fakeFs, targetFile5, []byte{}, filePerm)
	afero.WriteFile(fakeFs, targetFile6, []byte{}, filePerm)

	tests := []struct {
		name      string
		from      string
		to        string
		limit     int64
		offset    int64
		chunkSize int64
	}{
		{name: "simple test", from: file, to: targetFile1, limit: 1024, offset: 0, chunkSize: 1024},
		{name: "copying empty file", from: emptyFile, to: targetFile2, limit: 1024, offset: 0, chunkSize: 1024},
		{name: "testing limit (limit > chunk size)", from: file, to: targetFile3, limit: 3, offset: 0, chunkSize: 1024},
		{name: "testing limit (limit < chunk size)", from: file, to: targetFile4, limit: 3, offset: 0, chunkSize: 2},
		{name: "testing offset (offset > chunk size)", from: file, to: targetFile5, limit: 1024, offset: 3, chunkSize: 1024},
		{name: "testing offset (offset < chunk size)", from: file, to: targetFile6, limit: 1024, offset: 3, chunkSize: 2},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := CopyFile(fakeFs, tc.from, tc.to, tc.limit, tc.offset, tc.chunkSize)
			require.NoError(t, err)
			source, err := afero.ReadFile(fakeFs, tc.from)
			require.NoError(t, err)
			dest, err := afero.ReadFile(fakeFs, tc.to)
			require.NoError(t, err)
			var lim int
			switch {
			case len(source) == 0:
				lim = 0
			case int(tc.limit) < len(source):
				lim = int(tc.limit)
			default:
				lim = len(source)
			}
			require.Equal(t, source[tc.offset:lim], dest)
		})
	}
}
