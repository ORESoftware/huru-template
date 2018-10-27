package routes

import (
	// nearbym "huru/models/nearby"
	// personm "huru/models/person"
	// sharem "huru/models/share"

	"huru/routes/login"
	"huru/routes/nearby"
	"huru/routes/person"
	"huru/routes/register"
	"huru/routes/share"
)

// HuruRouteHandlers foobar
type HuruRouteHandlers struct {
	Register register.RegisterHandler
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
		register.RegisterHandler{},
		login.LoginHandler{},
		nearby.NearbyHandler{},
		share.ShareHandler{},
		person.PersonHandler{},
	}
}

// GetInjections money sauce
func (h HuruInjection) GetInjections(p personm.Map, n nearbym.Map, s sharem.Map) HuruInjection {
	return HuruInjection{
		person.PeopleInjection{People: p},
		share.ShareInjection{Share: s},
		nearby.NearbyInjection{Nearby: n},
	}
}
