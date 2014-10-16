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
	"strings"
)

var testFileReg = regexp.MustCompile("^.*_test\\.go$")
var testFuncReg = regexp.MustCompile("^Test")

func isTestFunc(decl *ast.FuncDecl) bool {
	if !testFuncReg.MatchString(decl.Name.Name) {
		return false
	}
	if len(decl.Type.Params.List) != 1 {
		return false
	}
	starExpr, ok := decl.Type.Params.List[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	selExpr, ok := starExpr.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := selExpr.X.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == "testing" && selExpr.Sel.Name == "T"
}

const (
	searching = iota
	running
	last
	done
)

type visitor struct {
	state int
	from  string
	to    string
	tests []string
}

func (self *visitor) Visit(node ast.Node) (w ast.Visitor) {
	if fd, ok := node.(*ast.FuncDecl); ok {
		if isTestFunc(fd) {
			switch self.state {
			case searching:
				if fd.Name.Name == self.from {
					self.state = running
				}
			case running:
				if fd.Name.Name == self.to {
					self.state = last
				}
			case last:
				self.state = done
			}
			if self.state == running || self.state == last {
				self.tests = append(self.tests, fmt.Sprintf("(%v)", fd.Name.Name))
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
	v := &visitor{
		from: *from,
		to:   *to,
	}
	for _, child := range children {
		if testFileReg.MatchString(child) {
			f, err := parser.ParseFile(&token.FileSet{}, child, nil, 0)
			if err != nil {
				panic(err)
			}
			ast.Walk(v, f)
		}
	}

	cmd := exec.Command("go", "test", "-v", fmt.Sprintf("-run=^%v$", strings.Join(v.tests, "|")))
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
