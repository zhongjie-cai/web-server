package webserver

import "reflect"

func isInterfaceValueNil(i any) bool {
	if i == nil {
		return true
	}
	var v = reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		return v.IsNil()
	}
	return !v.IsValid()
}
