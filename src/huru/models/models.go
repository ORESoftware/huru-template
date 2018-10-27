package models

import (
	"huru/models/nearby"
	"huru/models/person"
	"huru/models/share"
	"reflect"
)

// GetModels foo
func GetModels() interface{} {
	return struct {
		Near  interface{}
		Pers  interface{}
		Sharz interface{}
	}{
		reflect.TypeOf(nearby.Map{}),
		person.Map{},
		share.Map{},
	}
}
