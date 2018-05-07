// -------------------------------------------------------------
// DO NOT EDIT! this file was generated automatically by gomacro
// Any change will be lost when the file is re-generated
// -------------------------------------------------------------

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
 *     along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *
 * func1ret0.go
 *
 *  Created on Apr 16, 2017
 *      Author Massimiliano Ghilardi
 */

package fast

import (
	r "reflect"
	"unsafe"

	. "github.com/cosmos72/gomacro/base"
	xr "github.com/cosmos72/gomacro/xreflect"
)

func (c *Comp) func1ret0(t xr.Type, m *funcMaker) func(*Env) r.Value {

	nbinds := m.nbinds
	nintbinds := m.nintbinds
	funcbody := m.funcbody
	param0index := m.parambinds[0].Desc.Index()

	var debugC *Comp
	if c.Globals.Options&OptDebugger != 0 {
		debugC = c
	}

	targ0 := t.In(0)
	karg0 := targ0.Kind()
	switch karg0 {
	case r.Bool:
		{
			if funcbody == nil {
				return func(env *Env) r.Value {
					return r.ValueOf(func(bool) {
					})
				}
			}
			return func(env *Env) r.Value {

				env.MarkUsedByClosure()
				return r.ValueOf(func(arg0 bool) {
					env := newEnv4Func(env, nbinds, nintbinds, debugC)

					*(*bool)(unsafe.Pointer(&env.Ints[param0index])) = arg0

					funcbody(env)

					env.FreeEnv()
				})
			}
		}
	case r.Int:
		{
			if funcbody == nil {
				return func(env *Env) r.Value {
					return r.ValueOf(func(int) {
					})
				}
			}
			return func(env *Env) r.Value {

				env.MarkUsedByClosure()
				return r.ValueOf(func(arg0 int) {
					env := newEnv4Func(env, nbinds, nintbinds, debugC)

					*(*int)(unsafe.Pointer(&env.Ints[param0index])) = arg0

					funcbody(env)

					env.FreeEnv()
				})
			}
		}
	case r.Int8:
		{
			if funcbody == nil {
				return func(env *Env) r.Value {
					return r.ValueOf(func(int8) {
					})
				}
			}
			return func(env *Env) r.Value {

				env.MarkUsedByClosure()
				return r.ValueOf(func(arg0 int8) {
					env := newEnv4Func(env, nbinds, nintbinds, debugC)

					*(*int8)(unsafe.Pointer(&env.Ints[param0index])) = arg0

					funcbody(env)

					env.FreeEnv()
				})
			}
		}
	case r.Int16:
		{
			if funcbody == nil {
				return func(env *Env) r.Value {
					return r.ValueOf(func(int16) {
					})
				}
			}
			return func(env *Env) r.Value {

				env.MarkUsedByClosure()
				return r.ValueOf(func(arg0 int16) {
					env := newEnv4Func(env, nbinds, nintbinds, debugC)

					*(*int16)(unsafe.Pointer(&env.Ints[param0index])) = arg0

					funcbody(env)

					env.FreeEnv()
				})
			}
		}
	case r.Int32:
		{
			if funcbody == nil {
				return func(env *Env) r.Value {
					return r.ValueOf(func(int32) {
					})
				}
			}
			return func(env *Env) r.Value {

				env.MarkUsedByClosure()
				return r.ValueOf(func(arg0 int32) {
					env := newEnv4Func(env, nbinds, nintbinds, debugC)

					*(*int32)(unsafe.Pointer(&env.Ints[param0index])) = arg0

					funcbody(env)

					env.FreeEnv()
				})
			}
		}
	case r.Int64:
		{
			if funcbody == nil {
				return func(env *Env) r.Value {
					return r.ValueOf(func(int64) {
					})
				}
			}
			return func(env *Env) r.Value {

				env.MarkUsedByClosure()
				return r.ValueOf(func(arg0 int64) {
					env := newEnv4Func(env, nbinds, nintbinds, debugC)

					*(*int64)(unsafe.Pointer(&env.Ints[param0index])) = arg0

					funcbody(env)

					env.FreeEnv()
				})
			}
		}
	case r.Uint:
		{
			if funcbody == nil {
				return func(env *Env) r.Value {
					return r.ValueOf(func(uint) {
					})
				}
			}
			return func(env *Env) r.Value {

				env.MarkUsedByClosure()
				return r.ValueOf(func(arg0 uint) {
					env := newEnv4Func(env, nbinds, nintbinds, debugC)

					*(*uint)(unsafe.Pointer(&env.Ints[param0index])) = arg0

					funcbody(env)

					env.FreeEnv()
				})
			}
		}
	case r.Uint8:
		{
			if funcbody == nil {
				return func(env *Env) r.Value {
					return r.ValueOf(func(uint8) {
					})
				}
			}
			return func(env *Env) r.Value {

				env.MarkUsedByClosure()
				return r.ValueOf(func(arg0 uint8) {
					env := newEnv4Func(env, nbinds, nintbinds, debugC)

					*(*uint8)(unsafe.Pointer(&env.Ints[param0index])) = arg0

					funcbody(env)

					env.FreeEnv()
				})
			}
		}
	case r.Uint16:
		{
			if funcbody == nil {
				return func(env *Env) r.Value {
					return r.ValueOf(func(uint16) {
					})
				}
			}
			return func(env *Env) r.Value {

				env.MarkUsedByClosure()
				return r.ValueOf(func(arg0 uint16) {
					env := newEnv4Func(env, nbinds, nintbinds, debugC)

					*(*uint16)(unsafe.Pointer(&env.Ints[param0index])) = arg0

					funcbody(env)

					env.FreeEnv()
				})
			}
		}
	case r.Uint32:
		{
			if funcbody == nil {
				return func(env *Env) r.Value {
					return r.ValueOf(func(uint32) {
					})
				}
			}
			return func(env *Env) r.Value {

				env.MarkUsedByClosure()
				return r.ValueOf(func(arg0 uint32) {
					env := newEnv4Func(env, nbinds, nintbinds, debugC)

					*(*uint32)(unsafe.Pointer(&env.Ints[param0index])) = arg0

					funcbody(env)

					env.FreeEnv()
				})
			}
		}
	case r.Uint64:
		{
			if funcbody == nil {
				return func(env *Env) r.Value {
					return r.ValueOf(func(uint64) {
					})
				}
			}
			return func(env *Env) r.Value {

				env.MarkUsedByClosure()
				return r.ValueOf(func(arg0 uint64) {
					env := newEnv4Func(env, nbinds, nintbinds, debugC)

					env.Ints[param0index] = arg0

					funcbody(env)

					env.FreeEnv()
				})
			}
		}

	case r.Uintptr:
		{
			if funcbody == nil {
				return func(env *Env) r.Value {
					return r.ValueOf(func(uintptr) {
					})
				}
			}
			return func(env *Env) r.Value {

				env.MarkUsedByClosure()
				return r.ValueOf(func(arg0 uintptr) {
					env := newEnv4Func(env, nbinds, nintbinds, debugC)

					*(*uintptr)(unsafe.Pointer(&env.Ints[param0index])) = arg0

					funcbody(env)

					env.FreeEnv()
				})
			}
		}

	case r.Float32:
		{
			if funcbody == nil {
				return func(env *Env) r.Value {
					return r.ValueOf(func(float32) {
					})
				}
			}
			return func(env *Env) r.Value {

				env.MarkUsedByClosure()
				return r.ValueOf(func(arg0 float32) {
					env := newEnv4Func(env, nbinds, nintbinds, debugC)

					*(*float32)(unsafe.Pointer(&env.Ints[param0index])) = arg0

					funcbody(env)

					env.FreeEnv()
				})
			}
		}

	case r.Float64:
		{
			if funcbody == nil {
				return func(env *Env) r.Value {
					return r.ValueOf(func(float64) {
					})
				}
			}
			return func(env *Env) r.Value {

				env.MarkUsedByClosure()
				return r.ValueOf(func(arg0 float64) {
					env := newEnv4Func(env, nbinds, nintbinds, debugC)

					*(*float64)(unsafe.Pointer(&env.Ints[param0index])) = arg0

					funcbody(env)

					env.FreeEnv()
				})
			}
		}

	case r.Complex64:
		{
			if funcbody == nil {
				return func(env *Env) r.Value {
					return r.ValueOf(func(complex64) {
					})
				}
			}
			return func(env *Env) r.Value {

				env.MarkUsedByClosure()
				return r.ValueOf(func(arg0 complex64) {
					env := newEnv4Func(env, nbinds, nintbinds, debugC)

					*(*complex64)(unsafe.Pointer(&env.Ints[param0index])) = arg0

					funcbody(env)

					env.FreeEnv()
				})
			}
		}

	case r.Complex128:
		{
			if funcbody == nil {
				return func(env *Env) r.Value {
					return r.ValueOf(func(complex128) {
					})
				}
			}
			return func(env *Env) r.Value {

				env.MarkUsedByClosure()
				return r.ValueOf(func(arg0 complex128) {
					env := newEnv4Func(env, nbinds, nintbinds, debugC)
					{
						place := r.New(TypeOfComplex128).Elem()
						place.SetComplex(arg0,
						)
						env.Vals[param0index] = place
					}

					funcbody(env)

					env.FreeEnv()
				})
			}
		}

	case r.String:
		{
			if funcbody == nil {
				return func(env *Env) r.Value {
					return r.ValueOf(func(string) {
					})
				}
			}
			return func(env *Env) r.Value {

				env.MarkUsedByClosure()
				return r.ValueOf(func(arg0 string) {
					env := newEnv4Func(env, nbinds, nintbinds, debugC)

					{
						place := r.New(TypeOfString).Elem()
						place.SetString(arg0,
						)
						env.Vals[param0index] = place
					}
					funcbody(env)

					env.FreeEnv()
				})
			}
		}

	default:
		{
			rtype := t.ReflectType()
			if funcbody == nil {
				return func(env *Env) r.Value {
					return r.MakeFunc(rtype, func([]r.Value) []r.Value { return ZeroValues },
					)
				}
			} else {
				return func(env *Env) r.Value {

					env.MarkUsedByClosure()
					rtarg0 := targ0.ReflectType()
					return r.MakeFunc(rtype, func(args []r.Value) []r.Value {
						env := newEnv4Func(env, nbinds, nintbinds, debugC)

						if param0index != NoIndex {
							place := r.New(rtarg0).Elem()
							if arg0 := args[0]; arg0 != Nil && arg0 != None {
								place.Set(arg0.Convert(rtarg0))
							}

							env.Vals[param0index] = place
						}

						funcbody(env)
						return ZeroValues
					})
				}
			}

		}
	}
}
