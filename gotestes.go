package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"regexp"
)

var testFileReg = regexp.MustCompile("^.*_test\\.go$")

type visitor struct {
	from string
	to   string
}

func (self visitor) Visit(node ast.Node) (w ast.Visitor) {
	if fd, ok := node.(*ast.FuncDecl); ok {
		if fd.Name.Name >= self.from && fd.Name.Name <= self.to {
			cmd := exec.Command("go", "test", "-v", fmt.Sprintf("-test.run=^%v$", fd.Name.Name))
			cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
			if err := cmd.Run(); err != nil {
				panic(err)
			}
		}
	}
	return self
}

func main() {
	from := flag.String("from", "", "The first test to run")
	to := flag.String("to", "", "The last test to run")

	flag.Parse()

	dir, err := os.Open(".")
	if err != nil {
		panic(err)
	}
	children, err := dir.Readdirnames(-1)
	if err != nil {
		panic(err)
	}
	for _, child := range children {
		if testFileReg.MatchString(child) {
			f, err := parser.ParseFile(&token.FileSet{}, child, nil, 0)
			if err != nil {
				panic(err)
			}
			ast.Walk(visitor{
				from: *from,
				to:   *to,
			}, f)
		}
	}
}
