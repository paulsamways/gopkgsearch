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
	ekInterface
	ekType
	ekConst
	ekVariable
)

func (ek ElementKind) String() string {
	switch ek {
	case ekFunction:
		return "Function"
	case ekStruct:
		return "Struct"
	case ekInterface:
		return "Interface"
	case ekType:
		return "Type"
	case ekConst:
		return "Const"
	case ekVariable:
		return "Variable"
	}

	return ""
}

type Element struct {
	Package    string
	FilePath   string
	Name, name string
	Kind       string
	Source     string
	Doc        string
	recv       string
}

func printNode(fset *token.FileSet, node interface{}) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, node)

	return buf.String()
}

func printFieldList(fset *token.FileSet, fields *ast.FieldList) string {
	r := "("

	for _, f := range fields.List {
		if len(r) > 1 {
			r += ", "
		}

		r += printNode(fset, f.Type)
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

type makeElement func() *Element

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

			mE := func() *Element {
				return &Element{
					Package:  pkg.Name,
					FilePath: p,
				}
			}

			for _, decl := range file.Decls {
				switch t := decl.(type) {
				case *ast.FuncDecl:
					e := mE()
					indexFunc(e, fset, t)
					elements = append(elements, e)
				case *ast.GenDecl:
					es := indexGen(mE, fset, t)
					elements = append(elements, es...)
				}

			}
		}
	}

	return elements, nil
}

func indexFunc(e *Element, fset *token.FileSet, t *ast.FuncDecl) {
	e.Kind = ekFunction.String()
	e.Source = printNode(fset, t.Type)

	if t.Body != nil {
		e.Source += " " + printNode(fset, t.Body)
	}

	if t.Doc != nil {
		e.Doc = t.Doc.Text()
	}

	e.Name = t.Name.Name
	e.name = strings.ToLower(e.Name)

	if t.Recv != nil {
		e.Name += " " + printFieldList(fset, t.Recv)

		expr := t.Recv.List[0].Type

		switch exprT := expr.(type) {
		case *ast.Ident:
			e.recv = strings.ToLower(exprT.Name)
		case *ast.StarExpr:
			ident := exprT.X.(*ast.Ident)
			e.recv = strings.ToLower(ident.Name)
		}
	}
}

func indexGen(mE makeElement, fset *token.FileSet, t *ast.GenDecl) []*Element {
	es := make([]*Element, 0)

	switch t.Tok {
	case token.TYPE:
		for _, spec := range t.Specs {
			tSpec := spec.(*ast.TypeSpec)
			e := mE()
			indexType(e, fset, t, tSpec)
			es = append(es, e)
		}
	case token.CONST, token.VAR:
		for _, spec := range t.Specs {
			vSpec := spec.(*ast.ValueSpec)
			e := mE()
			indexConstOrVar(e, fset, t, vSpec)
			es = append(es, e)
		}
	}

	return es
}

func indexType(e *Element, fset *token.FileSet, d *ast.GenDecl, t *ast.TypeSpec) {
	if t.Doc != nil {
		e.Doc = t.Doc.Text()
	} else if d.Doc != nil {
		e.Doc = d.Doc.Text()
	}

	// need to inspect t.Type to determine whether its a struct,interface,etc...

	e.Source = printNode(fset, t)
	e.Name = t.Name.Name
	e.name = strings.ToLower(e.Name)

	switch t.Type.(type) {
	case *ast.InterfaceType:
		e.Kind = ekInterface.String()
	case *ast.StructType:
		e.Kind = ekStruct.String()
	default:
		e.Kind = ekType.String()
	}
}

func indexConstOrVar(e *Element, fset *token.FileSet, d *ast.GenDecl, v *ast.ValueSpec) {
	if v.Doc != nil {
		e.Doc = v.Doc.Text()
	} else if d.Doc != nil {
		e.Doc = d.Doc.Text()
	}

	e.Source = printNode(fset, v)
	e.Name = v.Names[0].Name // AFAIK this is always len == 1 for consts
	e.name = strings.ToLower(e.Name)

	switch d.Tok {
	case token.CONST:
		e.Kind = ekConst.String()
	case token.VAR:
		e.Kind = ekVariable.String()
	}
}
