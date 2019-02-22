package model

import (
	"context"
	"reflect"
)

type Setter interface {
	Set(string, interface{})
}

type Base interface {
	ToContext(c context.Context) context.Context
	SetContext(s Setter)
}

func MapToStruct(mapVal map[string]interface{}, val interface{}) (ok bool) {
	structVal := reflect.Indirect(reflect.ValueOf(val))
	for name, elem := range mapVal {
		structVal.FieldByName(name).Set(reflect.ValueOf(elem))
	}

	return
}
