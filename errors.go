package wingo

type TestRouteError byte

func (e TestRouteError) Error() string {
	return "need ignore TestRouteError"
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
