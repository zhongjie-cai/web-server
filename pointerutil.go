package webserver

import "reflect"

func isInterfaceValueNil(i interface{}) bool {
	if i == nil {
		return true
	}
	var v = reflectValueOf(i)
	if v.Kind() == reflect.Ptr {
		return v.IsNil()
	}
	return !v.IsValid()
}
