package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type ElementKind int

const (
	ekFunction ElementKind = iota
	ekStruct
	ekConst
)

func (ek ElementKind) String() string {
	switch ek {
	case ekFunction:
		return "Function"
	case ekStruct:
		return "Struct"
	case ekConst:
		return "Const"
	}

	return ""
}

type Element struct {
	Package, lowerPkg string
	FilePath          string
	Name, lowerName   string
	Kind              string
	Source            string
	Doc               string
}

func printNode(fset *token.FileSet, node interface{}) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, node)

	return buf.String()
}

func printFieldList(fset *token.FileSet, fields *ast.FieldList) string {
	r := "("

	for _, f := range fields.List {
		if len(f.Names) > 0 {
			r += f.Names[0].Name + " "
		}

		r += fmt.Sprintf("%T", f.Type)
	}

	return r + ")"
}

func getFilePath(fset *token.FileSet, path string, pos token.Pos) string {
	onDiskFile := fset.File(pos)

	if onDiskFile != nil {
		if rel, err := filepath.Rel(path, onDiskFile.Name()); err == nil {
			return rel
		} else {
			return onDiskFile.Name()
		}
	}
	return ""
}

func index(base string) ([]*Element, error) {
	elements := make([]*Element, 0)

	err := filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if strings.ToLower(info.Name()) == "testdata" {
				return filepath.SkipDir
			}

			e, err := parse(path, base)

			if err != nil {
				fmt.Printf("Couldn't parse Go files in directory %s. %s\n\n", path, err)
				return nil
			}

			elements = append(elements, e...)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return elements, nil
}

func parse(path, basePath string) ([]*Element, error) {
	elements := make([]*Element, 0)
	fset := token.NewFileSet()

	filter := func(f os.FileInfo) bool {
		return !f.IsDir() &&
			strings.HasSuffix(f.Name(), ".go") &&
			!strings.HasSuffix(f.Name(), "_test.go")
	}

	pkgs, err := parser.ParseDir(fset, path, filter, parser.ParseComments)

	if err != nil {
		return nil, err
	}

	for _, pkg := range pkgs {
		ast.PackageExports(pkg)

		for _, file := range pkg.Files {
			p := getFilePath(fset, basePath, file.Package)

			for _, decl := range file.Decls {
				switch t := decl.(type) {
				case *ast.FuncDecl:

					s := printNode(fset, t.Type)
					if t.Body != nil {
						s += " " + printNode(fset, t.Body)
					}

					d := ""
					if t.Doc != nil {
						d = t.Doc.Text()
					}

					n := t.Name.Name
					if t.Recv != nil {
						n += " " + printFieldList(fset, t.Recv)
					}

					e := Element{
						Package:   pkg.Name,
						lowerPkg:  strings.ToLower(pkg.Name),
						FilePath:  p,
						Name:      n,
						lowerName: strings.ToLower(n),
						Kind:      ekFunction.String(),
						Source:    s,
						Doc:       d,
					}

					elements = append(elements, &e)

					break
				case *ast.GenDecl:
					switch t.Tok {
					case token.TYPE:
						d := ""
						if t.Doc != nil {
							d = t.Doc.Text()
						}

						tSpec := t.Specs[0].(*ast.TypeSpec)

						s := printNode(fset, tSpec)

						e := Element{
							Package:   pkg.Name,
							lowerPkg:  strings.ToLower(pkg.Name),
							FilePath:  p,
							Name:      tSpec.Name.Name,
							lowerName: strings.ToLower(tSpec.Name.Name),
							Kind:      ekStruct.String(),
							Source:    s,
							Doc:       d,
						}

						elements = append(elements, &e)
						break
					case token.CONST:
						d := ""
						if t.Doc != nil {
							d = t.Doc.Text()
						}

						s := printNode(fset, t)

						for _, spec := range t.Specs {
							vSpec := spec.(*ast.ValueSpec)

							for _, name := range vSpec.Names {
								e := Element{
									Package:   pkg.Name,
									lowerPkg:  strings.ToLower(pkg.Name),
									FilePath:  p,
									Name:      name.Name,
									lowerName: strings.ToLower(name.Name),
									Kind:      ekConst.String(),
									Source:    s,
									Doc:       d,
								}

								elements = append(elements, &e)
							}
						}

						break
					}
					break
				}
			}
		}
	}

	return elements, nil
}
