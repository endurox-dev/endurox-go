/**
 * @brief Unified Buffer Format (UBF) marshal/unmarshal to/from structures
 *
 * @file typed_ubf_tag.go
 */
/* -----------------------------------------------------------------------------
 * Enduro/X Middleware Platform for Distributed Transaction Processing
 * Copyright (C) 2009-2016, ATR Baltic, Ltd. All Rights Reserved.
 * Copyright (C) 2017-2019, Mavimax, Ltd. All Rights Reserved.
 * This software is released under one of the following licenses:
 * LGPL or Mavimax's license for commercial use.
 * See LICENSE file for full text.
 *
 * C (as designed by Dennis Ritchie and later authors) language code is licensed
 * under Enduro/X Modified GNU Affero General Public License, version 3.
 * See LICENSE_C file for full text.
 * -----------------------------------------------------------------------------
 * LGPL license:
 *
 * This program is free software; you can redistribute it and/or modify it under
 * the terms of the GNU Lesser General Public License, version 3 as published
 * by the Free Software Foundation;
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT ANY
 * WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
 * PARTICULAR PURPOSE. See the GNU Lesser General Public License, version 3
 * for more details.
 *
 * You should have received a copy of the Lesser General Public License along
 * with this program; if not, write to the Free Software Foundation, Inc.,
 * 59 Temple Place, Suite 330, Boston, MA 02111-1307 USA
 *
 * -----------------------------------------------------------------------------
 * A commercial use license is available from Mavimax, Ltd
 * contact@mavimax.com
 * -----------------------------------------------------------------------------
 */
package atmi

import "reflect"
import "fmt"

// Unmarshal the value (UBF -> struct)
//@param occOne single field occurrence, if not FAIL, for unmarshalling
func (u *TypedUBF) unmarshalValue(p *reflect.StructField,
	rvv *reflect.Value, fldid int, occOne int) UBFError {

	//fmt.Printf("Field: [%s] Type: [%s]/%s Tag: [%s] %s -> %T\n",
	//		p.Name, p.Type, p.Type.Name(), p.Tag.Get("ubf"), p.Type.Kind(), p)

	switch p.Type.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		//fmt.Printf("This is int...\n")

		if FAIL == occOne {
			occOne = 0
		}

		if v, err := u.BGetInt64(fldid, occOne); err == nil {
			rvv.FieldByName(p.Name).SetInt(v)
		}

	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		//fmt.Printf("This is uint...\n")

		if FAIL == occOne {
			occOne = 0
		}

		if v, err := u.BGetInt64(fldid, occOne); err == nil {
			rvv.FieldByName(p.Name).SetUint(uint64(v))
		}
	case reflect.Float32,
		reflect.Float64:
		//fmt.Printf("This is float32/64...\n")

		if FAIL == occOne {
			occOne = 0
		}

		if v, err := u.BGetFloat64(fldid, occOne); err == nil {
			rvv.FieldByName(p.Name).SetFloat(v)
		}
	case reflect.String:
		//fmt.Printf("This is string...\n")

		if FAIL == occOne {
			occOne = 0
		}

		if v, err := u.BGetString(fldid, occOne); err == nil {
			rvv.FieldByName(p.Name).SetString(v)
		}
	case reflect.Slice:
		occ, _ := u.BOccur(fldid)
		if occ > 0 {

			capac := occ

			if FAIL != occOne {
				capac = 1
			}

			x := reflect.MakeSlice(p.Type, capac, capac)

			switch p.Type.Elem().Kind() {
			case reflect.Int,
				reflect.Int8,
				reflect.Int16,
				reflect.Int32,
				reflect.Int64:
				//fmt.Printf("Slice array... of int ...\n")

				occStart := 0
				occStop := occ

				if FAIL != occOne {
					occStart = occOne
					occStop = occOne + 1
				}

				for i := occStart; i < occStop; i++ {
					if v, err := u.BGetInt64(fldid, i); err == nil {
						if FAIL == occOne {
							x.Index(i).SetInt(v)
						} else {
							x.Index(0).SetInt(v)
						}

					}
				}
				rvv.FieldByName(p.Name).Set(x)
			case reflect.Uint,
				reflect.Uint8,
				reflect.Uint16,
				reflect.Uint32,
				reflect.Uint64:
				//fmt.Printf("Slice array... of uint ...\n")

				occStart := 0
				occStop := occ

				if FAIL != occOne {
					occStart = occOne
					occStop = occOne + 1
				}

				for i := occStart; i < occStop; i++ {
					if v, err := u.BGetInt64(fldid, i); err == nil {
						if FAIL == occOne {
							x.Index(i).SetUint(uint64(v))
						} else {
							x.Index(0).SetUint(uint64(v))
						}
					}
				}
				rvv.FieldByName(p.Name).Set(x)
			case reflect.Float32,
				reflect.Float64:
				//fmt.Printf("Slice array... of float64 ...\n")

				occStart := 0
				occStop := occ

				if FAIL != occOne {
					occStart = occOne
					occStop = occOne + 1
				}

				for i := occStart; i < occStop; i++ {
					if v, err := u.BGetFloat64(fldid, i); err == nil {
						if FAIL == occOne {
							x.Index(i).SetFloat(v)
						} else {
							x.Index(0).SetFloat(v)
						}

					}
				}
				rvv.FieldByName(p.Name).Set(x)
			case reflect.String:
				//fmt.Printf("Slice array... of string ...\n")

				occStart := 0
				occStop := occ

				if FAIL != occOne {
					occStart = occOne
					occStop = occOne + 1
				}

				for i := occStart; i < occStop; i++ {
					if v, err := u.BGetString(fldid, i); err == nil {
						if FAIL == occOne {
							x.Index(i).SetString(v)
						} else {
							x.Index(0).SetString(v)
						}
					}
				}
				rvv.FieldByName(p.Name).Set(x)
			case reflect.Slice:
				if reflect.Uint8 == p.Type.Elem().Elem().Kind() {
					//fmt.Printf("C_ARRAY support...\n")

					occStart := 0
					occStop := occ

					if FAIL != occOne {
						occStart = occOne
						occStop = occOne + 1
					}

					for i := occStart; i < occStop; i++ {
						if v, err := u.BGetByteArr(fldid, i); err == nil {
							//x.Index(i).Set(v)

							//Convert each
							val_len := len(v)
							y := reflect.MakeSlice(p.Type.Elem(), val_len, val_len)

							for j := 0; j < val_len; j++ {
								y.Index(j).SetUint(uint64(v[j]))
							}

							if FAIL == occOne {
								x.Index(i).Set(y)
							} else {
								x.Index(0).Set(y)
							}
						}
					}
					rvv.FieldByName(p.Name).Set(x)
				} else {
					return NewCustomUBFError(BEINVAL,
						fmt.Sprintf("%s - Not a C array!", p.Name))
				}

			default:
				return NewCustomUBFError(BEINVAL,
					fmt.Sprintf("%s - Unsupported slice!", p.Name))
			}
		}
	default:
		return NewCustomUBFError(BEINVAL,
			fmt.Sprintf("%s - Unsupported field type!", p.Name))
	}

	return nil
}

// Marshal the value (Struct -> UBF)
//@param occOne optional occurrence (if not FAIL) used for single field occurrences
// applies to arrays only.
func (u *TypedUBF) marshalValue(p *reflect.StructField,
	rvv *reflect.Value, fldid int, occOne int) UBFError {

	//fmt.Printf("Field: [%s] Type: [%s]/%s Tag: [%s] %s -> %T\n",
	//		p.Name, p.Type, p.Type.Name(), p.Tag.Get("ubf"), p.Type.Kind(), p)

	switch p.Type.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		//fmt.Printf("This is int...\n")

		if err := u.BChg(fldid, 0, rvv.FieldByName(p.Name).Int()); err != nil {
			return err
		}

	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		//fmt.Printf("This is uint...\n")

		if err := u.BChg(fldid, 0, rvv.FieldByName(p.Name).Uint()); err != nil {
			return err
		}
	case reflect.Float32,
		reflect.Float64:
		//fmt.Printf("This is float32/64...\n")
		if err := u.BChg(fldid, 0, rvv.FieldByName(p.Name).Float()); err != nil {
			return err
		}
	case reflect.String:
		//fmt.Printf("This is string...\n")

		if err := u.BChg(fldid, 0, rvv.FieldByName(p.Name).String()); err != nil {
			return err
		}
	case reflect.Slice:
		occ := rvv.FieldByName(p.Name).Len()

		if occOne != FAIL {
			if occOne >= occ {
				return NewCustomUBFError(BEINVAL,
					fmt.Sprintf("Invalid occurrence requested: %d - out of "+
						"bounds, tot occs: %d (valid range %d..%d",
						occOne, occ, 0, occ-1))
			}
		}

		if occ > 0 {
			x := reflect.MakeSlice(p.Type, occ, occ)

			switch p.Type.Elem().Kind() {
			case reflect.Int,
				reflect.Int8,
				reflect.Int16,
				reflect.Int32,
				reflect.Int64:
				//fmt.Printf("Slice array... of int ...\n")

				if occOne == FAIL {
					for i := 0; i < occ; i++ {
						if err := u.BChg(fldid, i,
							rvv.FieldByName(p.Name).Index(i).Int()); err != nil {
							return err
						}
					}
				} else {
					if err := u.BChg(fldid, 0,
						rvv.FieldByName(p.Name).Index(occOne).Int()); err != nil {
						return err
					}
				}

				rvv.FieldByName(p.Name).Set(x)
			case reflect.Uint,
				reflect.Uint8,
				reflect.Uint16,
				reflect.Uint32,
				reflect.Uint64:
				//fmt.Printf("Slice array... of uint ...\n")

				if occOne == FAIL {
					for i := 0; i < occ; i++ {
						if err := u.BChg(fldid, i,
							rvv.FieldByName(p.Name).Index(i).Uint()); err != nil {
							return err
						}
					}
				} else {
					if err := u.BChg(fldid, 0,
						rvv.FieldByName(p.Name).Index(occOne).Uint()); err != nil {
						return err
					}
				}
				rvv.FieldByName(p.Name).Set(x)
			case reflect.Float32,
				reflect.Float64:
				//fmt.Printf("Slice array... of float64 ...\n")

				if occOne == FAIL {
					for i := 0; i < occ; i++ {
						if err := u.BChg(fldid, i,
							rvv.FieldByName(p.Name).Index(i).Float()); err != nil {
							return err
						}
					}
				} else {
					if err := u.BChg(fldid, 0,
						rvv.FieldByName(p.Name).Index(occOne).Float()); err != nil {
						return err
					}
				}
				rvv.FieldByName(p.Name).Set(x)
			case reflect.String:
				//fmt.Printf("Slice array... of string ...\n")

				if occOne == FAIL {
					for i := 0; i < occ; i++ {
						if err := u.BChg(fldid, i,
							rvv.FieldByName(p.Name).Index(i).String()); err != nil {
							return err
						}
					}
				} else {
					if err := u.BChg(fldid, 0,
						rvv.FieldByName(p.Name).Index(occOne).String()); err != nil {
						return err
					}
				}
				rvv.FieldByName(p.Name).Set(x)
			case reflect.Slice:
				if reflect.Uint8 == p.Type.Elem().Elem().Kind() {
					//fmt.Printf("C_ARRAY support...\n")

					if occOne == FAIL {

						for i := 0; i < occ; i++ {
							if err := u.BChg(fldid, i,
								rvv.FieldByName(p.Name).Index(i).Bytes()); err != nil {
								return err
							}
						}
					} else {
						if err := u.BChg(fldid, 0,
							rvv.FieldByName(p.Name).Index(occOne).Bytes()); err != nil {
							return err
						}
					}

				} else {
					return NewCustomUBFError(BEINVAL,
						fmt.Sprintf("%s - Not a C array!", p.Name))
				}

			default:
				return NewCustomUBFError(BEINVAL,
					fmt.Sprintf("%s - Unsupported slice!", p.Name))
			}
		}
	default:
		return NewCustomUBFError(BEINVAL,
			fmt.Sprintf("%s - Unsupported field type!", p.Name))
	}

	return nil
}

//TODO: Add versions like UnmarshalOcc and MarshalOcc.

//Copy the specified fields to the local structure
//according to the `ubf' (i.e. take fields from UBF and copy to v structure).
//@param v  local struct
//@return UBF error
func (u *TypedUBF) Unmarshal(v interface{}) UBFError {
	return u._marshal(false, v, FAIL)
}

//Copy the structur in v struct to UBF
//@param v  local struct
//@return UBF error
func (u *TypedUBF) Marshal(v interface{}) UBFError {
	return u._marshal(true, v, FAIL)
}

//Copy the specified fields to the local structure
//according to the `ubf' (i.e. take fields from UBF and copy to v structure).
//@param v  local struct
//@param occ single occurrence in UBF to copy to either simple structure elements
//	or array elements. If copied to array, then it goes to first array element.
//@return UBF error
func (u *TypedUBF) UnmarshalSingle(v interface{}, occ int) UBFError {
	return u._marshal(false, v, occ)
}

//Copy the structur in v struct to UBF
//@param v  local struct
//@param occ single occurrence to marshal to UBF, i.e. single struct arrays occurrence
//	to copy to UBF. Applies only to structure array elements. In case of occurrence
//  is set which is out of the bounds of the array, the error will be generated as
//	BEINVAL
//@return UBF error
func (u *TypedUBF) MarshalSingle(v interface{}, occ int) UBFError {
	return u._marshal(true, v, occ)
}

//Copy the specified fields to the local structure
//or copy the local struct to UBF
//according to the `ubf'
//@param is_marshal true -> marshal mode, false -> unmarshal mode
//@param v  local struct
//@param occ occurrence to marshal/unmarshal
//@return UBF error
func (u *TypedUBF) _marshal(is_marshal bool, v interface{}, occ int) UBFError {

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return NewCustomUBFError(BEINVAL, "Struct is not ptr or nil")
	}

	rvv := rv.Elem()

	//fmt.Printf("rvv = %T\n", rvv)

	typ := reflect.TypeOf(v)
	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return NewCustomUBFError(BEINVAL, "Not a struct passed in!")
	}

	for i := 0; i < typ.NumField(); i++ {
		p := typ.Field(i)
		if !p.Anonymous {
			if p.Tag.Get("ubf") != "" {
				if fldid, _ := u.Buf.Ctx.BFldId(p.Tag.Get("ubf")); fldid != BBADFLDID {
					if is_marshal {
						if err := u.marshalValue(&p, &rvv, fldid, occ); nil != err {
							return err
						}
					} else {
						if err := u.unmarshalValue(&p, &rvv, fldid, occ); nil != err {
							return err
						}
					}
				} else {
					return NewCustomUBFError(BEINVAL,
						fmt.Sprintf("Field Name [%s] not resolved!",
							p.Tag.Get("ubf")))
				}
			}
		}
	}

	return nil
}
/* vim: set ts=4 sw=4 et smartindent: */
