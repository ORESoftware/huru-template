package models

import (
	"huru/models/nearby"
	"huru/models/person"
	"huru/models/share"
	"reflect"
)

type Container struct {
	Person person.Model
}

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
