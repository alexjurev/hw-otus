package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	reversedPhrase := stringutil.Reverse("Hello, OTUS!")
	fmt.Println(reversedPhrase)
}
