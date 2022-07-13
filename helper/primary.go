package helper

import "reflect"

func IsReflectPrimaryType(value reflect.Kind) bool {
	switch value {
	case reflect.Ptr, reflect.Interface, reflect.Struct, reflect.Map, reflect.Slice, reflect.Array, reflect.Chan, reflect.Func:
		return false
	}
	return true
}
