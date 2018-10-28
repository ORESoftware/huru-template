package models

import (
	"huru/models/nearby"
	"huru/models/person"
	"huru/models/share"
	"reflect"
)

// Container model container
// type Container struct {
// 	Person person.Model
// 	Nearby nearby.Model
// 	Share  share.Model
// }

type Person = person.Model
type Nearby = nearby.Model
type Share = share.Model

type PersonMap = person.Map
type NearbyMap = nearby.Map
type ShareMap = share.Map

var PersonInit = person.Init
var NearbyInit = nearby.Init
var ShareInit = share.Init

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
