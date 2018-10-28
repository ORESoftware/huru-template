package routes

import (
	"huru/models"
	"huru/routes/login"
	"huru/routes/nearby"
	"huru/routes/person"
	"huru/routes/register"
	"huru/routes/share"
)

// Handlers
type RegisterHandler = register.Handler
type LoginHandler = login.LoginHandler
type NearbyHandler = nearby.NearbyHandler
type ShareHandler = share.ShareHandler
type PersonHandler = person.PersonHandler

// Injection
type NearbyInjection = nearby.NearbyInjection
type ShareInjection = share.ShareInjection
type PersonInjection = person.PeopleInjection

// HuruRouteHandlers foobar
type HuruRouteHandlers struct {
	Register register.Handler
	Login    login.LoginHandler
	Nearby   nearby.NearbyHandler
	Share    share.ShareHandler
	Person   person.PersonHandler
}

// HuruInjection foobar
type HuruInjection struct {
	People person.PeopleInjection
	Share  share.ShareInjection
	Nearby nearby.NearbyInjection
}

// GetHandlers money sauce
func (h HuruRouteHandlers) GetHandlers() HuruRouteHandlers {
	return HuruRouteHandlers{
		register.Handler{},
		login.LoginHandler{},
		nearby.NearbyHandler{},
		share.ShareHandler{},
		person.PersonHandler{},
	}
}

// GetInjections money sauce
func (h HuruInjection) GetInjections(p models.PersonMap, n models.NearbyMap, s models.ShareMap) HuruInjection {
	return HuruInjection{
		person.PeopleInjection{People: p},
		share.ShareInjection{Share: s},
		nearby.NearbyInjection{Nearby: n},
	}
}
