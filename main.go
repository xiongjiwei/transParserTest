package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/astrewrite"
	tidbparser "github.com/pingcap/parser"
	"github.com/pingcap/parser/format"
	_ "github.com/pingcap/parser/test_driver"
)

func restore(sql string) string {
	p := tidbparser.New()
	p.EnableWindowFunc(true)
	var sb strings.Builder
	sb.Reset()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		fmt.Printf("%+v", err)
		panic("parse error")
	}

	restoreSQLs := ""
	for _, stmt := range stmtNodes {
		sb.Reset()
		err = stmt.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, &sb))
		if restoreSQLs != "" {
			restoreSQLs += "; "
		}
		restoreSQLs += sb.String()
	}

	return restoreSQLs
}

func main() {
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", os.Stdin, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	rewritten := astrewrite.Walk(f, rewriteFunc)
	_ = printer.Fprint(os.Stdout, fset, rewritten)
}

func rewriteFunc(n ast.Node) (ast.Node, bool) {
	if ret, ok := n.(*ast.FuncDecl); ok {
		if ret.Recv != nil && ret.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name == "testParserSuite" {
			for _, stmt := range ret.Body.List {
				if assign, ok := stmt.(*ast.AssignStmt); ok {
					for _, rh := range assign.Rhs {
						if lst, ok1 := rh.(*ast.CompositeLit); ok1 {
							if array, ok2 := lst.Type.(*ast.ArrayType); ok2 {
								if t, ok3 := array.Elt.(*ast.Ident); ok3 {
									if t.Name == "testCase" {
										for _, elt := range lst.Elts {
											elts := elt.(*ast.CompositeLit).Elts
											if elts[1].(*ast.Ident).Name == "true" {
												src := unquoteString(elts[0].(*ast.BasicLit).Value)
												re := restore(src)
												except := unquoteString(elts[2].(*ast.BasicLit).Value)
												if re != except && strings.Contains(re, "_UTF8MB4") && !strings.Contains(strings.ToLower(except), "_utf8mb4") {
													elts[2].(*ast.BasicLit).Value = strconv.Quote(re)
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return n, true
}

func unquoteString(s string) string {
	unquote, err := strconv.Unquote(s)
	if err != nil {
		panic(err)
	}
	return unquote
}
