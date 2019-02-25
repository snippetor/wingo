package wingo

type TestRouteError byte

func (e TestRouteError) Error() string {
	return "need ignore TestRouteError"
}
