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

func DetermineDataType(data string) *Data {
	return &Data{data: data, t: reflect.TypeOf(data).Kind()}
}

func (d *Data) GetData() any {
	return d.data
}

func (d *Data) GetInt() int {
	i, err := strconv.Atoi(d.data)
	if err != nil {
		log.Printf("can't convert to int: %v > %v", d.data, err)
		return 0
	}
	return i
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
