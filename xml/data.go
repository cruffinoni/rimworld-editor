package xml

import (
	"log"
	"reflect"
	"regexp"
	"strconv"
)

type Data struct {
	data      string
	t         reflect.Kind
	checkDone bool
}

const (
	Matrix = reflect.UnsafePointer + iota
	MatrixWithMapIndex
	RGBA
)

type associatedRegex struct {
	pattern *regexp.Regexp
	kind    reflect.Kind
}

var (
	HexRegex     = regexp.MustCompile(`(0x|#)[0-9a-fA-F]{2,}`)
	IntegerRegex = regexp.MustCompile(`^-?[0-9]+$`)
	FloatRegex   = regexp.MustCompile(`^-?[0-9]*\.?[0-9]+$`)
	BoolRegex    = regexp.MustCompile(`^(TRUE|FALSE|true|false)$`)

	AllPatterns = []associatedRegex{
		{pattern: IntegerRegex, kind: reflect.Int64},
		{pattern: FloatRegex, kind: reflect.Float64},
		{pattern: BoolRegex, kind: reflect.Bool},
		{pattern: HexRegex, kind: reflect.Uint64},
	}
)

func CreateDataType(data string) *Data {
	if data == "" {
		return &Data{data: data, t: reflect.Invalid}
	}
	for _, p := range AllPatterns {
		if p.pattern.MatchString(data) {
			return &Data{data: data, t: p.kind}
		}
	}
	return &Data{data: data, t: reflect.String}
}

func (d *Data) lazyCheck(destKind reflect.Kind) {
	if d.checkDone {
		return
	}
	defer func() {
		d.checkDone = true
	}()
	if d.t == reflect.Invalid {
		log.Panicf("lazyCheck: invalid data type")
	}
	// If the dest kind is the same as the source kind, we accept it anyway
	if destKind == reflect.String ||
		// Sometimes, float numbers doesn't have a decimal point so let's accept it as a float
		destKind == reflect.Float64 && d.t == reflect.Int64 {
		return
	}
	if destKind != d.t {
		log.Panicf("lazyCheck: type mismatch: %v != %v", d.t, destKind)
	}
}

func (d *Data) Kind() reflect.Kind {
	return d.t
}

func (d *Data) GetData() any {
	return d.data
}

func (d *Data) GetInt64() int64 {
	d.lazyCheck(reflect.Int64)
	i, err := strconv.ParseInt(d.data, 10, 64)
	if err != nil {
		log.Panicf("can't convert to int64: %v > %v", d.data, err)
		return 0
	}
	return i
}

func (d *Data) GetUint64() uint64 {
	d.lazyCheck(reflect.Uint64)
	i, err := strconv.ParseUint(d.data, 10, 64)
	if err != nil {
		log.Panicf("can't convert to uint64: %v > %v", d.data, err)
		return 0
	}
	return i
}

func (d *Data) GetString() string {
	d.lazyCheck(reflect.String)
	return d.data
}

func (d *Data) GetFloat64() float64 {
	d.lazyCheck(reflect.Float64)
	f, err := strconv.ParseFloat(d.data, 64)
	if err != nil {
		log.Panicf("can't convert to float: %v > %v", d.data, err)
		return 0.0
	}
	return f
}

func (d *Data) GetBool() bool {
	d.lazyCheck(reflect.Bool)
	b, err := strconv.ParseBool(d.data)
	if err != nil {
		log.Panicf("can't convert to bool: %v > %v", d.data, err)
		return false
	}
	return b
}
