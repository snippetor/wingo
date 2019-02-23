package wingo

import "reflect"

var controllerType = reflect.TypeOf((*Controller)(nil)).Elem()

func IsController(i interface{}) bool {
	return reflect.TypeOf(i).Implements(controllerType)
}

type Controller interface {
	Route(router *Router)
}
