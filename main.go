package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

func main() {
	run := extract(dirOrHere())
	os.Exit(run.summarize(os.Stdout))
}

func extract(dir string) *moduleExtractor {
	fs := token.NewFileSet()
	pkgs, err := parser.ParseDir(fs, dir, nil, parser.ParseComments)
	if err != nil {
		log.Fatal("Unable to parse dir:", err)
	}

	linter := &moduleExtractor{
		fs: fs,
	}

	for _, pkg := range pkgs {
		ast.Walk(linter, pkg)
	}

	return linter
}

func dirOrHere() string {
	if len(os.Args) > 1 {
		here := os.Args[1]
		stat, err := os.Stat(here)
		if err == nil && stat.IsDir() {
			return here
		}

		log.Fatalf("Not a valid directory: %q, error: %v", here, err)
	}
	if dir, err := os.Getwd(); err != nil {
		panic(err)
	} else {
		return dir
	}
}
