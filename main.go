package main

import (
	"flag"
	"fmt"
	pUtils "github.com/PaulSamways/utils/path"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

var gopaths = strings.Split(os.Getenv("GOPATH"), string(os.PathListSeparator))
var goroot = os.Getenv("GOROOT")

var addr = flag.String("addr", ":8000", "address to listen for web requests.")
var useGoPaths = flag.Bool("usegopath", false, "include GOPATH when searching for packages.")

type Element struct {
	Name, ImportPath, Comment, Kind string
}

func (e Element) String() string {
	return fmt.Sprintf("%v %s.%s", e.Kind, e.ImportPath, e.Name)
}

var elements []Element

func main() {
	flag.Parse()

	fmt.Println("Finding packages...")
	elements = find(filepath.Join(goroot, "src"))

	if *useGoPaths {
		for _, gp := range gopaths {
			v := find(filepath.Join(gp, "src"))
			elements = append(elements, v...)
		}
	}

	fmt.Printf("Starting web server on %s...\n", *addr)

	listen()
}

func find(path string) []Element {
	elements := make([]Element, 0)

	pUtils.WalkDirs(path, func(path string, files []string, err error) {
		hasGoFiles := false

		for _, v := range files {
			if strings.HasSuffix(v, ".go") && !strings.HasSuffix(v, "_test.go") {
				hasGoFiles = true
				break
			}
		}

		if hasGoFiles {
			e, err := parse(path)

			if err != nil {
				return
			}

			elements = append(elements, e...)
		}
	})

	return elements
}

func parse(path string) ([]Element, error) {
	fs := token.NewFileSet()
	pkgs, err := parser.ParseDir(fs, path, nil, 0)

	if err != nil {
		return nil, err
	}

	result := make([]Element, 0)

	for _, pkg := range pkgs {
		for _, f := range pkg.Files {

			if !strings.HasSuffix(f.Name.Name, "_test") && f.Scope != nil {
				for _, v := range f.Scope.Objects {
					if ast.IsExported(v.Name) {

						result = append(result, Element{
							v.Name,
							pkg.Name,
							"Comment",
							v.Kind.String(),
						})

					}
				}
			}

		}
	}

	return result, nil
}
