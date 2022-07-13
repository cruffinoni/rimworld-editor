package saver

import "reflect"

type Transformer interface {
	TransformToXML() ([]byte, error)
}

func TransformPrimaryType(value reflect.Value) ([]byte, error) {
	switch value.Kind() {
	case reflect.String:
		return []byte(value.String()), nil
	case reflect.Bool:
		return []byte(value.String()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(value.String()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []byte(value.String()), nil
	case reflect.Float32, reflect.Float64:
		return []byte(value.String()), nil
	}
	return nil, nil
}
