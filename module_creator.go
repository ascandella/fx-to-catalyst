package main

import (
	"bytes"
	"go/ast"
	"io"
)

var pkgMap = map[string]string{
	"uhttp": "catalysthttp",
}

type moduleCreator struct {
	pkgSel *ast.Ident
	fnSel  *ast.SelectorExpr
	lit    *ast.BasicLit
	args   []ast.Expr
}

func (m moduleCreator) String() string {
	out := &bytes.Buffer{}
	out.WriteString("module:\t")
	out.WriteString(m.pkgSel.Name)
	out.WriteString("\tfn:\t")
	out.WriteString(m.fnSel.Sel.Name)
	out.WriteString("\t args: [")
	for _, arg := range m.args {
		m.writeArg(out, arg)
	}
	out.WriteString("]")
	return out.String()
}

func (m moduleCreator) writeArg(out io.Writer, arg ast.Expr) {
	switch arg := arg.(type) {
	case *ast.Ident:
		io.WriteString(out, arg.Name)
		// TODO handle other cases
		// case *ast.CallExpr:
		// 	fun := arg.Fun
		// 	io.WriteString(out, fmt.Sprintf("(fun: %T %+v)", fun, fun))
		// default:
		// 	fmt.Printf("%T %+v\n", arg, arg)
	}
}

func (m moduleCreator) AsCatalyst() string {
	out := &bytes.Buffer{}
	out.WriteString("catalyst.Register(")

	if m.pkgSel != nil {
		modName, ok := pkgMap[m.pkgSel.Name]
		if !ok {
			modName = m.pkgSel.Name
		}

		out.WriteString(modName)
		// package-local methods will just have an ident
		if m.fnSel != nil {
			out.WriteString(".")
			out.WriteString(m.fnSel.Sel.Name)
		}
	}
	// TODO handle case of invoked ctor with zero args.
	if len(m.args) > 0 {
		out.WriteString("(")

		for _, arg := range m.args {
			m.writeArg(out, arg)
		}

		out.WriteString(")")
	}
	out.WriteString(")")

	return out.String()
}
