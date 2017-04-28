package main

import (
	"bytes"
	"fmt"
	"go/ast"
)

type moduleCreator struct {
	pkgSel *ast.Ident
	fnSel  *ast.SelectorExpr
	args   []ast.Expr
}

func (m moduleCreator) String() string {
	out := &bytes.Buffer{}
	out.WriteString(m.pkgSel.Name)
	out.WriteString(".")
	out.WriteString(m.fnSel.Sel.Name)
	out.WriteString("\t args: [")
	for _, arg := range m.args {
		out.WriteString(fmt.Sprintf("%T %+v", arg, arg))
	}
	out.WriteString("]")
	return out.String()
}
