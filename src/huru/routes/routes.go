package routes

import (
	"huru/routes/login"
	"huru/routes/nearby"
	"huru/routes/person"
	"huru/routes/register"
	"huru/routes/share"
)

// HuruRouteHandlers foobar
type HuruRouteHandlers struct {
	Register *register.RegisterHandler
	Login    *login.LoginHandler
	Nearby   *nearby.NearbyHandler
	Share    *share.ShareHandler
	Person   *person.PersonHandler
}

var People *person.PeopleInjection
var Share *share.ShareInjection
var Nearby *nearby.NearbyInjection
