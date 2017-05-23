package atmi
/* 
** Unified Buffer Format (UBF) marshal/unmarshal to/from structures
**
** @file typed_ubf_tag.go
** 
** -----------------------------------------------------------------------------
** Enduro/X Middleware Platform for Distributed Transaction Processing
** Copyright (C) 2015, Mavimax, Ltd. All Rights Reserved.
** This software is released under one of the following licenses:
** GPL or Mavimax's license for commercial use.
** -----------------------------------------------------------------------------
** GPL license:
** 
** This program is free software; you can redistribute it and/or modify it under
** the terms of the GNU General Public License as published by the Free Software
** Foundation; either version 2 of the License, or (at your option) any later
** version.
**
** This program is distributed in the hope that it will be useful, but WITHOUT ANY
** WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
** PARTICULAR PURPOSE. See the GNU General Public License for more details.
**
** You should have received a copy of the GNU General Public License along with
** this program; if not, write to the Free Software Foundation, Inc., 59 Temple
** Place, Suite 330, Boston, MA 02111-1307 USA
**
** -----------------------------------------------------------------------------
** A commercial use license is available from Mavimax, Ltd
** contact@mavimax.com
** -----------------------------------------------------------------------------
*/


import "reflect"
import "fmt"

// Unmarshal the value
func (u *TypedUBF) unmarshalValue(p *reflect.StructField,
	rvv *reflect.Value, fldid int) UBFError {

	//fmt.Printf("Field: [%s] Type: [%s]/%s Tag: [%s] %s -> %T\n",
	//		p.Name, p.Type, p.Type.Name(), p.Tag.Get("ubf"), p.Type.Kind(), p)

	switch p.Type.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		//fmt.Printf("This is int...\n")
		if v, err := u.BGetInt64(fldid, 0); err == nil {
			rvv.FieldByName(p.Name).SetInt(v)
		}
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		//fmt.Printf("This is uint...\n")
		if v, err := u.BGetInt64(fldid, 0); err == nil {
			rvv.FieldByName(p.Name).SetUint(uint64(v))
		}
	case reflect.Float32,
		reflect.Float64:
		//fmt.Printf("This is float32/64...\n")
		if v, err := u.BGetFloat64(fldid, 0); err == nil {
			rvv.FieldByName(p.Name).SetFloat(v)
		}
	case reflect.String:
		//fmt.Printf("This is string...\n")
		if v, err := u.BGetString(fldid, 0); err == nil {
			rvv.FieldByName(p.Name).SetString(v)
		}
	case reflect.Slice:
		occ, _ := u.BOccur(fldid)
		if occ > 0 {
			x := reflect.MakeSlice(p.Type, occ, occ)

			switch p.Type.Elem().Kind() {
			case reflect.Int,
				reflect.Int8,
				reflect.Int16,
				reflect.Int32,
				reflect.Int64:
				//fmt.Printf("Slice array... of int ...\n")
				for i := 0; i < occ; i++ {
					if v, err := u.BGetInt64(fldid, i); err == nil {
						x.Index(i).SetInt(v)
					}
				}
				rvv.FieldByName(p.Name).Set(x)
			case reflect.Uint,
				reflect.Uint8,
				reflect.Uint16,
				reflect.Uint32,
				reflect.Uint64:
				//fmt.Printf("Slice array... of uint ...\n")
				for i := 0; i < occ; i++ {
					if v, err := u.BGetInt64(fldid, i); err == nil {
						x.Index(i).SetUint(uint64(v))
					}
				}
				rvv.FieldByName(p.Name).Set(x)
			case reflect.Float32,
				reflect.Float64:
				//fmt.Printf("Slice array... of float64 ...\n")
				for i := 0; i < occ; i++ {
					if v, err := u.BGetFloat64(fldid, i); err == nil {
						x.Index(i).SetFloat(v)
					}
				}
				rvv.FieldByName(p.Name).Set(x)
			case reflect.String:
				//fmt.Printf("Slice array... of string ...\n")
				for i := 0; i < occ; i++ {
					if v, err := u.BGetString(fldid, i); err == nil {
						x.Index(i).SetString(v)
					}
				}
				rvv.FieldByName(p.Name).Set(x)
			case reflect.Slice:
				if reflect.Uint8 == p.Type.Elem().Elem().Kind() {
					//fmt.Printf("C_ARRAY support...\n")
					for i := 0; i < occ; i++ {
						if v, err := u.BGetByteArr(fldid, i); err == nil {
							//x.Index(i).Set(v)

							//Convert each
							val_len := len(v)
							y := reflect.MakeSlice(p.Type.Elem(), val_len, val_len)

							for j := 0; j < val_len; j++ {
								y.Index(j).SetUint(uint64(v[j]))
							}
							x.Index(i).Set(y)
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

// Unmarshal the value
func (u *TypedUBF) marshalValue(p *reflect.StructField,
	rvv *reflect.Value, fldid int) UBFError {

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
		if occ > 0 {
			x := reflect.MakeSlice(p.Type, occ, occ)

			switch p.Type.Elem().Kind() {
			case reflect.Int,
				reflect.Int8,
				reflect.Int16,
				reflect.Int32,
				reflect.Int64:
				//fmt.Printf("Slice array... of int ...\n")
				for i := 0; i < occ; i++ {
					if err := u.BChg(fldid, i,
						rvv.FieldByName(p.Name).Index(i).Int()); err != nil {
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
				for i := 0; i < occ; i++ {
					if err := u.BChg(fldid, i,
						rvv.FieldByName(p.Name).Index(i).Uint()); err != nil {
						return err
					}
				}
				rvv.FieldByName(p.Name).Set(x)
			case reflect.Float32,
				reflect.Float64:
				//fmt.Printf("Slice array... of float64 ...\n")
				for i := 0; i < occ; i++ {
					if err := u.BChg(fldid, i,
						rvv.FieldByName(p.Name).Index(i).Float()); err != nil {
						return err
					}
				}
				rvv.FieldByName(p.Name).Set(x)
			case reflect.String:
				//fmt.Printf("Slice array... of string ...\n")
				for i := 0; i < occ; i++ {
					if err := u.BChg(fldid, i,
						rvv.FieldByName(p.Name).Index(i).String()); err != nil {
						return err
					}
				}
				rvv.FieldByName(p.Name).Set(x)
			case reflect.Slice:
				if reflect.Uint8 == p.Type.Elem().Elem().Kind() {
					//fmt.Printf("C_ARRAY support...\n")

					for i := 0; i < occ; i++ {
						if err := u.BChg(fldid, i,
							rvv.FieldByName(p.Name).Index(i).Bytes()); err != nil {
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

//Copy the specified fields to the local structure
//according to the `ubf'
//@param v  local struct
//@return UBF error
func (u *TypedUBF) Unmarshal(v interface{}) UBFError {
	return u._marshal(false, v)
}

//Copy the specified fields to the local structure
//Copy the local struct to UBF
//@param v  local struct
//@return UBF error
func (u *TypedUBF) Marshal(v interface{}) UBFError {
	return u._marshal(true, v)
}

//Copy the specified fields to the local structure
//or copy the local struct to UBF
//according to the `ubf'
//@param is_marshal true -> marshal mode, false -> unmarshal mode
//@param v  local struct
//@return UBF error
func (u *TypedUBF) _marshal(is_marshal bool, v interface{}) UBFError {

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
						if err := u.marshalValue(&p, &rvv, fldid); nil != err {
							return err
						}
					} else {
						if err := u.unmarshalValue(&p, &rvv, fldid); nil != err {
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
