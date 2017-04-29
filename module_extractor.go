package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"strings"
)

const (
	// TODO: will be replaced with `fx` once di-refactor is done
	_svcPackage = "service"
)

type moduleExtractor struct {
	fs       *token.FileSet
	modules  []moduleCreator
	out      io.Writer
	debugOut io.Writer
	debugme  *bool
}

func (m *moduleExtractor) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.Package:
		if n.Name != "main" {
			return nil
		}
		m.debug("entering main package")
		return m
	case *ast.FuncDecl:
		// We only care about the main func
		if n.Name.Name != "main" {
			return nil
		}
		m.debug("entering main func")
		return m
	case *ast.CallExpr:
		switch fun := n.Fun.(type) {
		case *ast.SelectorExpr:
			if x, ok := fun.X.(*ast.Ident); ok {
				if x.Name == _svcPackage && fun.Sel.Name == "WithModule" {
					m.extractWithModule(n.Args)
					return m
				} else {
					m.debug("Ignoring SelectorExpr.X of non-WithModule %T %+v", x, x.Name)
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
			m.debug("unknown arg to serice.WithModule, expected *ast.CallExpr, got %T %+v", arg, arg)
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
			m.debug("Unknown selector epr (non-Ident): %T %+v", x, x)
			return
		}
	default:
		m.debug("Unknown call.Fun (non-SelectorExpr): %T %+v", fun, fun)
		return
	}
	m.modules = append(m.modules, mc)
}

func (m *moduleExtractor) isDebug() bool {
	if *m.debugme {
		return true
	}
	return os.Getenv("DEBUG_EXTRACT") != ""
}

func (m *moduleExtractor) debug(msg string, args ...interface{}) {
	if m.debugme == nil || !(*m.debugme) {
		return
	}

	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Fprintf(m.debugOut, msg, args...)
}

type extractOption func(*moduleExtractor)

func withWriter(out io.Writer) extractOption {
	return func(m *moduleExtractor) {
		m.out = out
		if m.debugOut != nil {
			m.debugOut = out
		}
	}
}

func withDebug(out io.Writer) extractOption {
	return func(m *moduleExtractor) {
		m.debugOut = out
	}
}

func (m *moduleExtractor) summarize(opts ...extractOption) int {
	// set up default
	m.out = os.Stdout
	m.debugOut = os.Stderr
	for _, opt := range opts {
		opt(m)
	}

	if len(m.modules) == 0 {
		fmt.Fprintf(m.out, "[ERROR] No UberFx modules detected\n")
		return 1
	}

	fmt.Fprintf(m.out, "Input modules: \n\n")

	for _, mod := range m.modules {
		fmt.Fprintf(m.out, fmt.Sprintf("\t%v\n", mod))
	}

	fmt.Fprintf(m.out, "\n\nCatalyst func:\n\nfunc init() {\n")

	for _, mod := range m.modules {
		fmt.Fprintf(m.out, fmt.Sprintf("\t%s\n", mod.AsCatalyst()))
	}

	fmt.Fprintf(m.out, "}\n")

	return 0
}
