package wingo

type TestRouteError byte

func (e TestRouteError) Error() string {
	return "need ignore TestRouteError"
}

type AuthorizationBreakError byte

func (e AuthorizationBreakError) Error() string {
	return "need ignore AuthorizationBreakError"
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
