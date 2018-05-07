/*
 * gomacro - A Go interpreter with Lisp-like macros
 *
 * Copyright (C) 2017-2018 Massimiliano Ghilardi
 *
 *     This program is free software: you can redistribute it and/or modify
 *     it under the terms of the GNU Lesser General Public License as published
 *     by the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 *
 *     This program is distributed in the hope that it will be useful,
 *     but WITHOUT ANY WARRANTY; without even the implied warranty of
 *     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *     GNU Lesser General Public License for more details.
 *
 *     You should have received a copy of the GNU Lesser General Public License
 *     along with this program.  If not, see <https://www.gnu.org/licenses/lgpl>.
 *
 *
 * compile.go
 *
 *  Created on Apr 01, 2017
 *      Author Massimiliano Ghilardi
 */

package fast

import (
	"go/ast"
	"go/token"
	"go/types"
	r "reflect"

	. "github.com/cosmos72/gomacro/ast2"
	. "github.com/cosmos72/gomacro/base"
	xr "github.com/cosmos72/gomacro/xreflect"
)

func New() *Interp {
	top := newTopInterp("builtin")
	top.env.UsedByClosure = true // do not free this *Env
	file := NewInnerInterp(top, "main", "main")
	file.env.UsedByClosure = true // do not free this *Env
	return file
}

func newTopInterp(path string) *Interp {
	name := FileName(path)

	globals := NewGlobals()
	universe := xr.NewUniverse()

	compGlobals := &CompGlobals{
		Globals:      globals,
		Universe:     universe,
		KnownImports: make(map[string]*Import),
		interf2proxy: make(map[r.Type]r.Type),
		proxy2interf: make(map[r.Type]xr.Type),
		Prompt:       "gomacro> ",
	}
	envGlobals := &ThreadGlobals{Globals: globals}
	ce := &Interp{
		Comp: &Comp{
			CompGlobals: compGlobals,
			CompBinds: CompBinds{
				Name: name,
				Path: path,
			},
			UpCost: 1,
			Depth:  0,
			Outer:  nil,
		},
		env: &Env{
			Outer:         nil,
			ThreadGlobals: envGlobals,
		},
	}
	// tell xreflect about our packages "fast" and "main"
	universe.CachePackage(types.NewPackage("fast", "fast"))
	universe.CachePackage(types.NewPackage("main", "main"))

	// no need to scavenge for Builtin, Function,  Macro and UntypedLit fields and methods.
	// actually, making them opaque helps securing against malicious interpreted code.
	for _, rtype := range []r.Type{rtypeOfBuiltin, rtypeOfFunction, rtypeOfPtrImport, rtypeOfMacro} {
		compGlobals.opaqueType(rtype, "fast")
	}
	compGlobals.opaqueType(rtypeOfUntypedLit, "untyped")

	envGlobals.TopEnv = ce.env
	ce.addBuiltins()
	return ce
}

func NewInnerInterp(outer *Interp, name string, path string) *Interp {
	if len(name) == 0 {
		name = FileName(path)
	}

	outerComp := outer.Comp
	outerEnv := outer.env
	ir := &Interp{
		Comp: &Comp{
			CompGlobals: outerComp.CompGlobals,
			CompBinds: CompBinds{
				Name: name,
				Path: path,
			},
			UpCost: 1,
			Depth:  outerComp.Depth + 1,
			Outer:  outerComp,
		},
		env: &Env{
			Outer:         outerEnv,
			ThreadGlobals: outerEnv.ThreadGlobals,
		},
	}
	if outerEnv.Outer == nil {
		outerEnv.ThreadGlobals.FileEnv = ir.env
	}
	return ir
}

func NewComp(outer *Comp, code *Code) *Comp {
	if outer == nil {
		return &Comp{UpCost: 1}
	}
	c := Comp{
		UpCost:      1,
		Depth:       outer.Depth + 1,
		Outer:       outer,
		CompGlobals: outer.CompGlobals,
	}
	// Debugf("NewComp(%p->%p) %s", outer, &c, debug.Stack())
	if code != nil {
		c.Code = *code
	}
	return &c
}

func (c *Comp) TopComp() *Comp {
	for ; c != nil; c = c.Outer {
		if c.Outer == nil {
			break
		}
	}
	return c
}

func (c *Comp) FileComp() *Comp {
	for ; c != nil; c = c.Outer {
		outer := c.Outer
		if outer == nil || outer.Outer == nil {
			break
		}
	}
	return c
}

// if a function Env only declares ignored binds, it gets this scratch buffers
var ignoredBinds = []r.Value{Nil}
var ignoredIntBinds = []uint64{0}

func NewEnv(outer *Env, nbinds int, nintbinds int) *Env {
	g := outer.ThreadGlobals
	pool := &g.Pool // pool is an array, do NOT copy it!
	index := g.PoolSize - 1
	var env *Env
	if index >= 0 {
		g.PoolSize = index
		env = pool[index]
		pool[index] = nil
	} else {
		env = &Env{}
	}
	if nbinds <= 1 {
		env.Vals = ignoredBinds
	} else if cap(env.Vals) < nbinds {
		env.Vals = make([]r.Value, nbinds)
	} else {
		env.Vals = env.Vals[0:nbinds]
	}
	if nintbinds <= 1 {
		env.Ints = ignoredIntBinds
	} else if cap(env.Ints) < nintbinds {
		env.Ints = make([]uint64, nintbinds)
	} else {
		env.Ints = env.Ints[0:nintbinds]
	}
	env.Outer = outer
	env.IP = outer.IP
	env.Code = outer.Code
	env.ThreadGlobals = g
	env.DebugPos = outer.DebugPos
	env.CallDepth = g.CallDepth
	return env
}

func newEnv4Func(outer *Env, nbinds int, nintbinds int, debugComp *Comp) *Env {
	g := outer.ThreadGlobals
	pool := &g.Pool // pool is an array, do NOT copy it!
	index := g.PoolSize - 1
	var env *Env
	if index >= 0 {
		g.PoolSize = index
		env = pool[index]
		pool[index] = nil
	} else {
		env = &Env{}
	}
	if nbinds <= 1 {
		env.Vals = ignoredBinds
	} else if cap(env.Vals) < nbinds {
		env.Vals = make([]r.Value, nbinds)
	} else {
		env.Vals = env.Vals[0:nbinds]
	}
	if nintbinds <= 1 {
		env.Ints = ignoredIntBinds
	} else if cap(env.Ints) < nintbinds {
		env.Ints = make([]uint64, nintbinds)
	} else {
		env.Ints = env.Ints[0:nintbinds]
	}
	env.Outer = outer
	env.ThreadGlobals = g
	env.DebugComp = debugComp
	env.CallDepth = g.CallDepth + 1
	// Debugf("newEnv4Func(%p->%p) binds=%d intbinds=%d", outer, env, nbinds, nintbinds)
	return env
}

func (env *Env) MarkUsedByClosure() {
	for ; env != nil && !env.UsedByClosure; env = env.Outer {
		env.UsedByClosure = true
	}
}

// FreeEnv tells the interpreter that given Env is no longer needed.
func (env *Env) FreeEnv() {
	// Debugf("FreeEnv(%p->%p), IP = %d of %d", env, env.Outer, env.Outer.IP, len(env.Outer.Code))
	if env.UsedByClosure {
		// in use, cannot recycle
		return
	}
	common := env.ThreadGlobals
	n := common.PoolSize
	if n >= PoolCapacity {
		return
	}
	if env.IntAddressTaken {
		env.Ints = nil
		env.IntAddressTaken = false
	}
	env.Outer = nil
	env.Code = nil
	env.DebugPos = nil
	env.DebugComp = nil
	env.ThreadGlobals = nil
	common.Pool[n] = env // pool is an array, be careful NOT to copy it!
	common.PoolSize = n + 1
}

func (env *Env) Top() *Env {
	for ; env != nil; env = env.Outer {
		if env.Outer == nil {
			break
		}
	}
	return env
}

func (env *Env) File() *Env {
	for ; env != nil; env = env.Outer {
		outer := env.Outer
		if outer == nil || outer.Outer == nil {
			break
		}
	}
	return env
}

// combined Parse + MacroExpandCodeWalk
func (c *Comp) Parse(src string) Ast {
	c.Line = 0
	nodes := c.ParseBytes([]byte(src))
	forms := anyToAst(nodes, "Parse")

	forms, _ = c.MacroExpandCodewalk(forms)
	if c.Options&OptShowMacroExpand != 0 {
		c.Debugf("after macroexpansion: %v", forms.Interface())
	}
	return forms
}

func (c *Comp) Compile(in Ast) *Expr {
	switch form := in.(type) {
	case nil:
		return nil
	case AstWithNode:
		return c.CompileNode(form.Node())
	case AstWithSlice:
		n := form.Size()
		var list []*Expr
		for i := 0; i < n; i++ {
			e := c.Compile(form.Get(i))
			if e != nil {
				list = append(list, e)
			}
		}
		return exprList(list, c.CompileOptions())
	}
	c.Errorf("unsupported Ast node, expecting <AstWithNode> or <AstWithSlice>, found %v <%v>", in, r.TypeOf(in))
	return nil
}

// compileExpr is a wrapper for Compile
// that guarantees Code does not get clobbered/cleared.
// Used by Comp.Quasiquote
func (c *Comp) compileExpr(in Ast) *Expr {
	cf := NewComp(c, nil)
	cf.UpCost = 0
	cf.Depth--
	return cf.Compile(in)
}

func (c *Comp) CompileNode(node ast.Node) *Expr {
	if n := c.Code.Len(); n != 0 {
		c.Warnf("Compile: discarding %d previously compiled statements from code buffer", n)
	}
	c.Code.Clear()
	if node == nil {
		return nil
	}
	c.Pos = node.Pos()
	switch node := node.(type) {
	case ast.Decl:
		c.Decl(node)
	case ast.Expr:
		return c.Expr(node, nil)
	case *ast.ExprStmt:
		// special case of statement
		return c.Expr(node.X, nil)
	case ast.Stmt:
		c.Stmt(node)
	case *ast.File:
		c.File(node)
	default:
		c.Errorf("unsupported node type, expecting <ast.Decl>, <ast.Expr>, <ast.Stmt> or <*ast.File>, found %v <%v>", node, r.TypeOf(node))
		return nil
	}
	return c.Code.AsExpr()
}

func (c *Comp) File(node *ast.File) {
	c.Name = node.Name.Name
	for _, decl := range node.Decls {
		c.Decl(decl)
	}
}

func (c *Comp) Append(stmt Stmt, pos token.Pos) {
	c.Code.Append(stmt, pos)
}

func (c *Comp) append(stmt Stmt) {
	c.Code.Append(stmt, c.Pos)
}
