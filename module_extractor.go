package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
)

const (
	// TODO: will be replaced with `fx` once di-refactor is done
	_svcPackage = "service"
)

type moduleExtractor struct {
	fs      *token.FileSet
	modules []moduleCreator
}

func (m *moduleExtractor) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.Package:
		if n.Name != "main" {
			return nil
		}
		debug("entering main package")
		return m
	case *ast.FuncDecl:
		// We only care about the main func
		if n.Name.Name != "main" {
			return nil
		}
		debug("entering main func")
		return m
	case *ast.CallExpr:
		switch fun := n.Fun.(type) {
		case *ast.SelectorExpr:
			if x, ok := fun.X.(*ast.Ident); ok {
				if x.Name == _svcPackage && fun.Sel.Name == "WithModule" {
					m.extractWithModule(n.Args)
					return m
				} else {
					debug("Ignoring SelectorExpr.X of non-WithModule %T %+v", x, x.Name)
				}
			}
		}
		return m
	}
	return m
}

func (m *moduleExtractor) extractWithModule(args []ast.Expr) {
	for _, arg := range args {
		switch arg := arg.(type) {
		case *ast.CallExpr:
			m.addModuleCall(arg)
		case *ast.Ident:
			m.addModuleIdent(arg)
		case *ast.BasicLit:
			m.addSimple(arg)
		default:
			debug("unknown arg to serice.WithModule, expected *ast.CallExpr, got %T %+v", arg, arg)
		}
	}
}

func (m *moduleExtractor) addSimple(arg *ast.BasicLit) {
	m.modules = append(m.modules, moduleCreator{
		lit: arg,
	})
}

func (m *moduleExtractor) addModuleIdent(arg *ast.Ident) {
	m.modules = append(m.modules, moduleCreator{
		pkgSel: arg,
	})
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
			debug("Unknown selector epr (non-Ident): %T %+v", x, x)
			return
		}
	default:
		debug("Unknown call.Fun (non-SelectorExpr): %T %+v", fun, fun)
		return
	}
	m.modules = append(m.modules, mc)
}

func (m *moduleExtractor) summarize(out io.Writer) int {
	if len(m.modules) == 0 {
		fmt.Fprintf(out, "[ERROR] No UberFx modules detected\n")
		return 1
	}

	fmt.Fprintf(out, "Input modules: \n\n")

	for _, mod := range m.modules {
		fmt.Fprintf(out, fmt.Sprintf("\t%v\n", mod))
	}

	fmt.Fprintf(out, "\n\nCatalyst func:\n\nfunc init() {\n")

	for _, mod := range m.modules {
		fmt.Fprintf(out, fmt.Sprintf("\t%s\n", mod.AsCatalyst()))
	}

	fmt.Fprintf(out, "}\n")

	return 0
}
