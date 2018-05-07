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
 * import_wrappers.go
 *
 *  Created on May 26, 2017
 *      Author Massimiliano Ghilardi
 */

package base

import (
	"bytes"
	"fmt"
	"go/constant"
	"go/types"
	"math"
	"strconv"
	"strings"
)

type genimport struct {
	globals     *Globals
	mode        ImportMode
	gpkg        *types.Package
	scope       *types.Scope
	names       []string
	pkgrenames  map[string]string // map[path]name of packages to import, where name:s are guaranteed to be unique
	out         *bytes.Buffer
	path        string
	name, name_ string
	proxyprefix string
	reflect     string
}

func (g *Globals) writeImportFile(out *bytes.Buffer, path string, gpkg *types.Package, mode ImportMode) (isEmpty bool) {

	gen := g.newGenImport(out, path, gpkg, mode)
	if gen == nil {
		return true
	}
	gen.write()
	return false
}

func (g *Globals) newGenImport(out *bytes.Buffer, path string, gpkg *types.Package, mode ImportMode) *genimport {
	scope := gpkg.Scope()
	names := scope.Names()

	isEmpty := true
	for _, name := range names {
		if obj := scope.Lookup(name); obj.Exported() {
			switch obj.(type) {
			case *types.Const, *types.Var, *types.Func, *types.TypeName:
				isEmpty = false
				break
			}
		}
	}
	if isEmpty {
		return nil
	}

	gen := &genimport{globals: g, mode: mode, gpkg: gpkg, scope: scope, names: names, out: out, path: path}

	name := FileName(path)
	name = sanitizeIdentifier(name)
	gen.name = name

	if mode == ImInception {
		gen.reflect = "r."
	} else {
		gen.name_ = name + "."
	}
	if mode == ImPlugin {
		gen.proxyprefix = "P_"
	} else {
		gen.proxyprefix = fmt.Sprintf("P_%s_", sanitizeIdentifier(path))
	}
	return gen
}

func (gen *genimport) write() {

	gen.writePreamble()

	gen.writeBinds()
	gen.writeTypes()
	gen.writeProxies()
	gen.writeUntypeds()
	gen.writeWrappers()

	gen.out.WriteString("\n\t}\n}\n")
	gen.writeInterfaceProxies()
}

type mapdecl struct {
	out  *bytes.Buffer
	head string
	foot string
}

func (gen *genimport) mapdecl(head string) mapdecl {
	if strings.IndexByte(head, '%') >= 0 {
		head = fmt.Sprintf(head, gen.reflect)
	}
	return mapdecl{gen.out, head, ""}
}

func (d *mapdecl) header() {
	if len(d.head) != 0 {
		d.out.WriteString(d.head)
		d.out.WriteByte('{')
		d.head = ""
		d.foot = "\n\t}"
	}
}

func (d *mapdecl) footer() {
	if len(d.foot) != 0 {
		d.out.WriteString(d.foot)
		d.out.WriteString(", ")
	}
}

func (d *mapdecl) footer1(comma bool) {
	if len(d.foot) != 0 {
		d.out.WriteString(d.foot)
		if comma {
			d.out.WriteString(", ")
		}
	}
}

func (gen *genimport) collectPackageImportsWithRename(requireAllInterfaceMethodsExported bool) {
	gen.pkgrenames = gen.globals.CollectPackageImportsWithRename(gen.gpkg, requireAllInterfaceMethodsExported)
}

func (gen *genimport) writePreamble() {
	mode := gen.mode
	out := gen.out
	path := gen.path

	var alias, filepkg string
	switch mode {
	case ImBuiltin:
		alias = "_b "
		filepkg = "imports"
	case ImThirdParty:
		filepkg = "thirdparty"
	case ImPlugin:
		filepkg = "main"
	case ImInception:
		alias = "_i "
		filepkg = gen.name
	}

	fmt.Fprintf(gen.out, `// this file was generated by gomacro command: import %s%q
// DO NOT EDIT! Any change will be lost when the file is re-generated

package %s

import (`, alias, path, filepkg)

	var imports string
	if mode == ImInception {
		fmt.Fprintf(gen.out, "\n\tr \"reflect\"\n\t\"github.com/cosmos72/gomacro/imports\"")
		imports = "imports."
	} else {
		fmt.Fprintf(out, "\n\t. \"reflect\"")
	}
	gen.collectPackageImportsWithRename(true)
	for path, name := range gen.pkgrenames {
		if mode == ImInception && path == gen.path {
			continue // writing inside the package: it should not import itself
		} else if name == FileName(path) {
			fmt.Fprintf(out, "\n\t%q", path)
		} else {
			fmt.Fprintf(out, "\n\t%s %q", name, path)
		}
	}
	fmt.Fprintf(out, "\n)\n")

	if mode == ImPlugin {
		fmt.Fprint(out, `
type Package = struct {
	Binds    map[string]Value
	Types    map[string]Type
	Proxies  map[string]Type
	Untypeds map[string]string
	Wrappers map[string][]string
}

var Packages = make(map[string]Package)

func main() {
}

`)
	}

	fmt.Fprintf(out, `
// reflection: allow interpreted code to import %q
func init() {
	%sPackages[%q] = %sPackage{
	`, path, imports, path, imports)
}

func (gen *genimport) writeBinds() {
	d := gen.mapdecl("Binds: map[string]%sValue")

	for _, name := range gen.names {
		if obj := gen.scope.Lookup(name); obj.Exported() {
			switch obj := obj.(type) {
			case *types.Const:
				val := obj.Val()
				var conv1, conv2 string
				if t, ok := obj.Type().(*types.Basic); ok && t.Info()&types.IsUntyped != 0 {
					// untyped constants have arbitrary precision... they may overflow integers.
					// this is just an approximation, use Package.Untypeds for exact value
					if val.Kind() == constant.Int {
						str := val.ExactString()
						conv1, conv2 = gen.globals.detectIntKind(gen.path, name, str)
					}
				}
				d.header()
				fmt.Fprintf(gen.out, "\n\t\t%q:\t%sValueOf(%s%s%s%s),", name, gen.reflect, conv1, gen.name_, name, conv2)
			case *types.Var:
				d.header()
				fmt.Fprintf(gen.out, "\n\t\t%q:\t%sValueOf(&%s%s).Elem(),", name, gen.reflect, gen.name_, name)
			case *types.Func:
				d.header()
				fmt.Fprintf(gen.out, "\n\t\t%q:\t%sValueOf(%s%s),", name, gen.reflect, gen.name_, name)
			}
		}
	}
	d.footer()
}

func (gen *genimport) writeTypes() {
	d := gen.mapdecl("Types: map[string]%sType")

	for _, name := range gen.names {
		if obj := gen.scope.Lookup(name); obj.Exported() {
			switch obj.(type) {
			case *types.TypeName:
				d.header()
				fmt.Fprintf(gen.out, "\n\t\t%q:\t%sTypeOf((*%s%s)(nil)).Elem(),", name, gen.reflect, gen.name_, name)
			}
		}
	}
	d.footer()
}

func (gen *genimport) writeProxies() {
	d := gen.mapdecl("Proxies: map[string]%sType")

	for _, name := range gen.names {
		if obj := gen.scope.Lookup(name); obj.Exported() {
			if t := extractInterface(obj, true); t != nil {
				d.header()
				fmt.Fprintf(gen.out, "\n\t\t%q:\t%sTypeOf((*%s%s)(nil)).Elem(),", name, gen.reflect, gen.proxyprefix, name)
			}
		}
	}
	d.footer()
}

func (gen *genimport) writeUntypeds() {
	d := gen.mapdecl("Untypeds: map[string]string")

	for _, name := range gen.names {
		if obj := gen.scope.Lookup(name); obj.Exported() {
			switch obj := obj.(type) {
			case *types.Const:
				if t, ok := obj.Type().(*types.Basic); ok && t.Info()&types.IsUntyped != 0 {
					rkind := UntypedKindToReflectKind(t.Kind())
					str := MarshalUntyped(rkind, obj.Val())
					if len(str) != 0 {
						d.header()
						fmt.Fprintf(gen.out, "\n\t\t%q:\t%q,", name, str)
					}
				}
			}
		}
	}
	d.footer()
}

// find wrapper methods and write them. needed for accurate method selection.
func (gen *genimport) writeWrappers() {
	d := gen.mapdecl("Wrappers: map[string][]string")

	for _, name := range gen.names {
		if obj := gen.scope.Lookup(name); obj.Exported() {
			switch obj.(type) {
			case *types.TypeName:
				if t, ok := obj.Type().(*types.Named); ok {
					// only structs can have embedded fields, and thus wrapper methods for embedded fields
					if _, ok := t.Underlying().(*types.Struct); ok {
						wrappers := new(analyzer).Analyze(t)
						if len(wrappers) != 0 {
							d.header()
							fmt.Fprintf(gen.out, "\n\t\t%q:\t[]string{", obj.Name())
							for _, wrapper := range wrappers {
								fmt.Fprintf(gen.out, "%q,", wrapper)
							}
							fmt.Fprint(gen.out, "},")
						}
					}
				}
			}
		}
	}
	d.footer()
}

// write proxies that pre-implement package's interfaces
func (gen *genimport) writeInterfaceProxies() {
	path := gen.gpkg.Path()
	for _, name := range gen.names {
		obj := gen.scope.Lookup(name)
		if t := extractInterface(obj, true); t != nil {
			gen.writeInterfaceProxy(path, name, t)
		}
	}
}

func (g *Globals) detectIntKind(path, name, str string) (string, string) {
	i, err := strconv.ParseInt(str, 0, 64)
	if err == nil {
		if i == int64(int32(i)) {
			// constant fits int32. We can use the default (i.e. int)
			// on both 32-bit and 64-bit platforms
			return "", ""
		} else if i == int64(uint32(i)) {
			// constant fits uint32
			return "uint32(", ")"
		} else {
			return "int64(", ")"
		}
	}
	_, err = strconv.ParseUint(str, 0, 64)
	if err == nil {
		return "uint64(", ")"
	}
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		// nothing fits... leave the default
		return "", ""
	} else {
		prefix := "float64"
		f = math.Abs(f)
		if f == float64(float32(f)) && f <= math.MaxFloat32 && f >= math.SmallestNonzeroFloat32 {
			// float32 loses no precision vs. float64
			prefix = "float32"
		}
		g.Warnf("package %q: integer constant %s = %s overflows both int64 and uint64, converting to %s", path, name, str, prefix)
		return prefix + "(", ")"
	}
}
