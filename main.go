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

	elements = Search(filepath.Join(goroot, "src/pkg"))

	if *useGoPaths {
		for _, gp := range gopaths {
			e := Search(filepath.Join(gp, "src"))
			elements = append(elements, e...)
		}
	}

	fmt.Printf("Indexed %d elements\n", len(elements))
	fmt.Printf("Starting web server on %s\n", *addr)

	listen()
}
