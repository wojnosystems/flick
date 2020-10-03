package optional_set_value

import (
	"errors"
	"github.com/wojnosystems/flick/pkg/set_value"
	"github.com/wojnosystems/go-optional"
	"reflect"
	"strconv"
)

// Register the optional types
func Register(r *set_value.Registry) *set_value.Registry {
	r.Register(reflect.TypeOf((*optional.String)(nil)).Elem(), func(settableDst interface{}, src string) (err error) {
		s := settableDst.(*optional.String)
		s.Set(src)
		return
	})
	r.Register(reflect.TypeOf((*optional.Int)(nil)).Elem(), func(settableDst interface{}, src string) (err error) {
		s := settableDst.(*optional.Int)
		var v int64
		v, err = strconv.ParseInt(src, 10, 64)
		if err != nil {
			return
		}
		s.Set(int(v))
		return
	})
	r.Register(reflect.TypeOf((*optional.Int8)(nil)).Elem(), func(settableDst interface{}, src string) (err error) {
		s := settableDst.(*optional.Int8)
		var v int64
		v, err = strconv.ParseInt(src, 10, 8)
		if err != nil {
			return
		}
		s.Set(int8(v))
		return
	})
	r.Register(reflect.TypeOf((*optional.Int16)(nil)).Elem(), func(settableDst interface{}, src string) (err error) {
		s := settableDst.(*optional.Int16)
		var v int64
		v, err = strconv.ParseInt(src, 10, 16)
		if err != nil {
			return
		}
		s.Set(int16(v))
		return
	})
	r.Register(reflect.TypeOf((*optional.Int32)(nil)).Elem(), func(settableDst interface{}, src string) (err error) {
		s := settableDst.(*optional.Int32)
		var v int64
		v, err = strconv.ParseInt(src, 10, 32)
		if err != nil {
			return
		}
		s.Set(int32(v))
		return
	})
	r.Register(reflect.TypeOf((*optional.Int64)(nil)).Elem(), func(settableDst interface{}, src string) (err error) {
		s := settableDst.(*optional.Int64)
		var v int64
		v, err = strconv.ParseInt(src, 10, 64)
		if err != nil {
			return
		}
		s.Set(v)
		return
	})
	r.Register(reflect.TypeOf((*optional.Uint)(nil)).Elem(), func(settableDst interface{}, src string) (err error) {
		s := settableDst.(*optional.Uint)
		var v uint64
		v, err = strconv.ParseUint(src, 10, 64)
		if err != nil {
			return
		}
		s.Set(uint(v))
		return
	})
	r.Register(reflect.TypeOf((*optional.Uint8)(nil)).Elem(), func(settableDst interface{}, src string) (err error) {
		s := settableDst.(*optional.Uint8)
		var v uint64
		v, err = strconv.ParseUint(src, 10, 8)
		if err != nil {
			return
		}
		s.Set(uint8(v))
		return
	})
	r.Register(reflect.TypeOf((*optional.Uint16)(nil)).Elem(), func(settableDst interface{}, src string) (err error) {
		s := settableDst.(*optional.Uint16)
		var v uint64
		v, err = strconv.ParseUint(src, 10, 16)
		if err != nil {
			return
		}
		s.Set(uint16(v))
		return
	})
	r.Register(reflect.TypeOf((*optional.Uint32)(nil)).Elem(), func(settableDst interface{}, src string) (err error) {
		s := settableDst.(*optional.Uint32)
		var v uint64
		v, err = strconv.ParseUint(src, 10, 32)
		if err != nil {
			return
		}
		s.Set(uint32(v))
		return
	})
	r.Register(reflect.TypeOf((*optional.Uint64)(nil)).Elem(), func(settableDst interface{}, src string) (err error) {
		s := settableDst.(*optional.Uint64)
		var v uint64
		v, err = strconv.ParseUint(src, 10, 64)
		if err != nil {
			return
		}
		s.Set(uint64(v))
		return
	})
	r.Register(reflect.TypeOf((*optional.Float32)(nil)).Elem(), func(settableDst interface{}, src string) (err error) {
		s := settableDst.(*optional.Float32)
		var v float64
		v, err = strconv.ParseFloat(src, 32)
		if err != nil {
			return
		}
		s.Set(float32(v))
		return
	})
	r.Register(reflect.TypeOf((*optional.Float64)(nil)).Elem(), func(settableDst interface{}, src string) (err error) {
		s := settableDst.(*optional.Float64)
		var v float64
		v, err = strconv.ParseFloat(src, 64)
		if err != nil {
			return
		}
		s.Set(v)
		return
	})
	r.Register(reflect.TypeOf((*optional.Bool)(nil)).Elem(), func(settableDst interface{}, src string) (err error) {
		if src == "t" || src == "true" {
			settableDst.(*optional.Bool).Set(true)
			return
		} else if src == "" || src == "f" || src == "false" {
			settableDst.(*optional.Bool).Set(false)
			return
		}
		return errors.New("unable to convert string to boolean value")
	})
	return r
}
