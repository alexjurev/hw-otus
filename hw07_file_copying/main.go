package main

import (
	"flag"
	"fmt"

	"github.com/spf13/afero"
)

var (
	from, to                 string
	limit, offset, chunkSize int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
	flag.Int64Var(&chunkSize, "chunk-size", 1024, "size of part of file which is copyed simultaneously")
}

func main() {
	flag.Parse()
	fs := afero.NewOsFs()

	fmt.Printf("Copying file %s to %s\n", from, to)
	err := CopyFile(fs, from, to, limit, offset, chunkSize)
	if err != nil {
		fmt.Printf("Failed to copy files: %s\n", err)
		return
	}
	fmt.Println("Copying successfully done")
}
