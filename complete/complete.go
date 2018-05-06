package complete

import (
    "go/ast"
    "go/token"

    "errors"
    "fmt"
	"log"
    "reflect"
    "strings"

	"github.com/cosmos72/gomacro/ast2"
	"github.com/cosmos72/gomacro/base"
	interp "github.com/cosmos72/gomacro/fast"
)

func extractLastExpr(stmt ast.Stmt) ast.Expr {

    switch s := stmt.(type) {
    case *ast.ForStmt:
        // TODO(jpark): check Init, Cond, Post if Body is null
        if s.Body != nil && len(s.Body.List) > 0 {
            stmtList := s.Body.List
            return extractLastExpr(stmtList[len(stmtList)-1])
        }
    case *ast.AssignStmt:
        if len(s.Rhs) > 0 {
            return s.Rhs[len(s.Rhs)-1];
        }
    case *ast.ExprStmt:
        return s.X
    default:
        log.Printf("Unhandled statement type: %T\n", s)
    }

    return nil
}

// parse incomplete code
func parsePartial(src []byte) ([]ast.Stmt, ast.Expr) {
    var psr parser

    fset := token.NewFileSet()
    psr.init(fset, "", src, 0)

    psr.openScope()
    psr.pkgScope = psr.topScope

	var decls []ast.Stmt
    for psr.tok == token.IMPORT {
        decl := ast.DeclStmt{
            Decl: psr.parseGenDecl(token.IMPORT, psr.parseImportSpec),
        }
        decls = append(decls, &decl)
    }

    stmts := psr.parseStmtList()
    stmts = append(decls, stmts...)

    if len(psr.errors) > 1 {
        for i := 1; i < len(psr.errors); i++ {
            err := psr.errors[i]
            log.Printf("Parse error at %d:%d: %s\n",
                err.Pos.Line, err.Pos.Column,
                err.Msg)
        }
    }

    if len(stmts) == 0 {
        return nil, nil
    }

    prevStmts := stmts[:len(stmts)-1]
    lastStmt := stmts[len(stmts)-1]

    // extract expression from the last statement
    expr := extractLastExpr(lastStmt)

    return prevStmts, expr
}

// Based on the expression at the cursor, find completion matches 
func findMatchFromExpr(ir *interp.Interp, expr ast.Expr) (matches []string, prefix string, err error) {
    log.Printf("Expression: %v(%T)\n", expr, expr)

    switch e := expr.(type) {
    case *ast.BinaryExpr:
        return findMatchFromExpr(ir, e.Y)
    case *ast.CallExpr:
        args := e.Args
        if len(args) > 0 {
            lastArg := args[len(args)-1]

            return findMatchFromExpr(ir, lastArg)
        }
    case *ast.Ident:
        prefix = e.Name
        compiler := ir.Comp

        for _, bind := range compiler.Binds {
            if strings.HasPrefix(bind.Name, prefix) {
                matches = append(matches, bind.Name)
            }
        }
    case *ast.IndexExpr:
        return findMatchFromExpr(ir, e.Index)
    case *ast.ParenExpr:
        return findMatchFromExpr(ir, e.X)
    case *ast.SelectorExpr:
        if (e.Sel.Name == "_") {
            prefix = ""
        } else {
            prefix = e.Sel.Name
        }

        srcAst := ast2.AnyToAst(e.X, "")
        compiled := ir.CompileAst(srcAst)
        v := ir.RunExpr1(compiled)
        t := v.Type()

        if (t.Kind() == reflect.Ptr) {
            t = t.Elem()
        }

        log.Printf("Value: %v(%T)\n", v, v)
        log.Printf("Value Interface: %v(%T)\n", v.Interface(), v.Interface())
        log.Printf("Type: %v(%T)\n", t, t)

        importExpr, ok := v.Interface().(interp.Import)
        log.Printf("ok: %s\n", ok)
        if ok {
            log.Printf("Import expr: %v\n", importExpr)
            // base expr is a package
            for name := range importExpr.Binds {
                if strings.HasPrefix(name, prefix) {
                    matches = append(matches, name)
                }
            }
        } else {
            for i := 0; i < t.NumField(); i++ {
                field := t.Field(i)
                name := field.Name
                if e.Sel.Name == "_" || strings.HasPrefix(name, prefix) {
                    matches = append(matches, name)
                }
            }
            for i := 0; i < t.NumMethod(); i++ {
                method := t.Method(i)
                name := method.Name
                if e.Sel.Name == "_" || strings.HasPrefix(name, prefix) {
                    matches = append(matches, name)
                }
            }
        }
    case *ast.UnaryExpr:
        return findMatchFromExpr(ir, e.X)
    default:
        log.Printf("Unhandled expression type: %T\n", e)
    }

    return matches, prefix, err
}

// find matches for the given code at the cursor position
// this function parses the code into statements, and finds the expression at the cursor
// and then based on the expression, find the possible matches
func FindMatch(ir *interp.Interp, code string, curPos int) (matches []string, prefix string, err error) {
	// Capture a panic from the evaluation if one occurs and store it in the `err` return parameter.
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(error); !ok {
				err = errors.New(fmt.Sprint(r))
			}
		}
	}()

	// Prepare and perform the multiline evaluation.
	compiler := ir.Comp

	// Don't show the gomacro prompt.
	compiler.Options &^= base.OptShowPrompt

	// Don't swallow panics as they are recovered above and handled with a Jupyter `error` message instead.
	compiler.Options &^= base.OptTrapPanic

	// Reset the error line so that error messages correspond to the lines from the cell.
	compiler.Line = 0

    // Parse the code, and find the last expression
    src := []byte(code)[:curPos]
    stmts, expr := parsePartial(src)

    // Eval statements before the cursor
    for _, stmt := range stmts {
        srcBlock := src[stmt.Pos()-1:stmt.End()-1]
        nodes := compiler.ParseBytes(srcBlock)
        srcAst := ast2.AnyToAst(nodes, "doEval")
        compiled := ir.CompileAst(srcAst)
        ir.RunExpr(compiled)
    }

    return findMatchFromExpr(ir, expr)
}

