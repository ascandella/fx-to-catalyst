package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
)

type moduleExtractor struct {
	fs      *token.FileSet
	modules []moduleCreator
}

func (m *moduleExtractor) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.FuncDecl:
		// We only care about the main func
		if n.Name.Name != "main" {
			return nil
		}
		return m
	case *ast.CallExpr:
		switch fun := n.Fun.(type) {
		case *ast.SelectorExpr:
			if x, ok := fun.X.(*ast.Ident); ok && x.Name == "service" {
				if fun.Sel.Name == "WithModule" {
					m.extractWithModule(n.Args)
					return m
				}
			}
		}
		return m

	default:
		return m
	}
}

func (m *moduleExtractor) extractWithModule(args []ast.Expr) {
	for _, arg := range args {
		switch arg := arg.(type) {
		case *ast.CallExpr:
			m.addModuleCall(arg)
		default:
			// TODO info warn logging
		}
	}
}

func (m *moduleExtractor) addModuleCall(call *ast.CallExpr) {
	mc := moduleCreator{
		args: call.Args,
	}
	switch fun := call.Fun.(type) {
	case *ast.SelectorExpr:
		mc.fnSel = fun
		if x, ok := fun.X.(*ast.Ident); ok {
			mc.pkgSel = x
		} else {
			// TODO error logging
		}
	default:
		// TODO error logging
		return
	}
	m.modules = append(m.modules, mc)
}

func debug(prefix string, node ast.Node) {
	fmt.Printf("%s: %T %+v\n", prefix, node, node)
}

func (m *moduleExtractor) summarize(out io.Writer) int {
	if len(m.modules) == 0 {
		fmt.Fprintf(out, "[ERROR] No UberFx modules detected")
		return 1
	}

	fmt.Fprintf(out, "Modules: \n\n")

	for _, mod := range m.modules {
		fmt.Fprintf(out, fmt.Sprintf("\t%v\n", mod))
	}

	return 0
}
