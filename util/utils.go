package util

import (
	"unsafe"
	"reflect"
)

func B2S(b []byte) string {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&b))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*string)(unsafe.Pointer(&bh))
}

func S2B(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
	Data: sh.Data,
	Len:  sh.Len,
	Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}