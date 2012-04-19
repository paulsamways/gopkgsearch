package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var gopaths = strings.Split(os.Getenv("GOPATH"), string(os.PathListSeparator))
var goroot = os.Getenv("GOROOT")

var addr = flag.String("addr", ":8000", "address to listen for web requests.")
var useGoPaths = flag.Bool("usegopath", false, "include GOPATH when searching for packages.")

var elements []*Element

func main() {
	flag.Parse()

	fmt.Println("Indexing...")

	var err error

	elements, err = index(filepath.Join(goroot, "src/pkg"))

	if err != nil {
		fmt.Printf("Couldn't index the Go packages. %s\n", err)
		elements = make([]*Element, 0)
	}

	if *useGoPaths {
		for _, gp := range gopaths {
			if e, err := index(filepath.Join(gp, "src")); err != nil {
				fmt.Printf("Couldn't index the packages in %s. %s\n", filepath.Join(gp, "src"), err)
			} else {
				elements = append(elements, e...)
			}
		}
	}

	fmt.Printf("Indexed %d elements\n", len(elements))
	fmt.Printf("Starting web server on %s\n", *addr)

	listen()
}
