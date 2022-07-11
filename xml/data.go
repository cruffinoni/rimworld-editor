package xml

import (
	"log"
	"reflect"
	"strconv"
)

type Data struct {
	data string
	t    reflect.Kind
}

func CreateDataType(data string) *Data {
	return &Data{data: data, t: reflect.TypeOf(data).Kind()}
}

func (d *Data) GetKind() reflect.Kind {
	return d.t
}

func (d *Data) GetData() any {
	return d.data
}

func (d *Data) GetInt64() int64 {
	i, err := strconv.ParseInt(d.data, 10, 64)
	if err != nil {
		log.Printf("can't convert to int64: %v > %v", d.data, err)
		return 0
	}
	return i
}

func (d *Data) GetInt() int {
	return int(d.GetInt64())
}

func (d *Data) GetInt8() int8 {
	return int8(d.GetInt64())
}

func (d *Data) GetInt16() int16 {
	return int16(d.GetInt64())
}

func (d *Data) GetInt32() int32 {
	return int32(d.GetInt64())
}

func (d *Data) GetString() string {
	return d.data
}

func (d *Data) GetFloat() float64 {
	f, err := strconv.ParseFloat(d.data, 64)
	if err != nil {
		log.Printf("can't convert to float: %v > %v", d.data, err)
		return 0.0
	}
	return f
}

func (d *Data) GetBool() bool {
	b, err := strconv.ParseBool(d.data)
	if err != nil {
		log.Printf("can't convert to bool: %v > %v", d.data, err)
		return false
	}
	return b
}
