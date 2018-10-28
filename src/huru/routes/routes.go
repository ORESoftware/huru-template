package routes

import (
	"huru/routes/login"
	"huru/routes/nearby"
	"huru/routes/person"
	"huru/routes/register"
	"huru/routes/share"
)

// Handlers
type RegisterHandler = register.Handler
type LoginHandler = login.Handler
type NearbyHandler = nearby.Handler
type ShareHandler = share.Handler
type PersonHandler = person.Handler

type HandlerCreator = func() struct{}

var Handlers = map[string]HandlerCreator{
	"Register": register.CreateHandler,
}

// Injection
type NearbyInjection = nearby.NearbyInjection
type ShareInjection = share.ShareInjection
type PersonInjection = person.PeopleInjection

// HuruRouteHandlers foobar
type HuruRouteHandlers struct {
	Register register.Handler
	Login    login.Handler
	Nearby   nearby.Handler
	Share    share.Handler
	Person   person.Handler
}

// GetHandlers money sauce
func (h HuruRouteHandlers) GetHandlers() HuruRouteHandlers {
	return HuruRouteHandlers{
		register.Handler{},
		login.Handler{},
		nearby.Handler{},
		share.Handler{},
		person.Handler{},
	}
}
