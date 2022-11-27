package main

import (
	"fmt"
	"os"
	"regexp"
)

func main() {
	fmt.Println(regexp.MustCompile(`([^/\s]+)$`).FindStringSubmatch(os.Args[0])[1])
}
