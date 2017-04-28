package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
)

type moduleExtractor struct {
	fs      *token.FileSet
	modules []ast.Node
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
		fmt.Printf("Call: %T %+v %+v\n", n.Fun, n.Fun, n.Args)
		switch fun := n.Fun.(type) {
		case *ast.SelectorExpr:
			debug("fun", fun.X)
			if x, ok := fun.X.(*ast.Ident); ok && x.Name == "service" {
				if fun.Sel.Name == "WithModule" {
					// TODO
					ast.Fprint(os.Stderr, m.fs, n, nil)
					return m
				}
			}
		}
		return m

	default:
		return m
	}
}

func (m *moduleExtractor) extractMain(n *ast.FuncDecl) {
	for _, stmt := range n.Body.List {
		debug("statement", stmt)
		fmt.Println()

		switch stmt := stmt.(type) {
		case *ast.AssignStmt:
			fmt.Println("Inspecting assign")
			m.inspectAssign(stmt)
		case *ast.ExprStmt:
			if call, ok := stmt.X.(*ast.CallExpr); ok {
				m.inspectCall(call)
			}
		}
	}
}

func (m *moduleExtractor) inspectAssign(assn *ast.AssignStmt) {
	debug("assn", assn)
	for _, lh := range assn.Lhs {
		debug("lhs", lh)
	}
	for _, rh := range assn.Rhs {
		debug("rhs", rh)
	}
}

func (m *moduleExtractor) inspectCall(call *ast.CallExpr) {
	debug("callexpr", call)

	debug("call fun", call.Fun)

	for _, arg := range call.Args {
		debug("arg", arg)
	}
}

func debug(prefix string, node ast.Node) {
	fmt.Printf("%s: %T %+v\n", prefix, node, node)
}

func (m *moduleExtractor) summarize(out io.Writer) int {
	if len(m.modules) == 0 {
		fmt.Fprintf(out, "[ERROR] No UberFx modules detected")
		return 1
	}

	for _, mod := range m.modules {
		fmt.Fprintf(out, fmt.Sprintf("%+v\n", mod))
	}

	return 0
}
