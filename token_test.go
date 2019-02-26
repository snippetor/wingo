package wingo

import (
	"fmt"
	"testing"
)

var TestSecretKey = "7npc8w4g7BZsQQTVafadd2zhjRyBhPnv"

func TestToken(a *testing.T) {
	token := marshalToken(&TokenPayload{Id: 20001, Name: "carl", Admin: true}, TestSecretKey)
	fmt.Println(token)
	t := &TokenPayload{}

	if !unmarshalToken(token, "7npc8w4g7BZsQQTVafadd2zhjRyBhPnv", t) {
		a.Fail()
	}

	fmt.Println(t)

	if t.Id != 20001 {
		a.Fail()
	}

	if t.Name != "carl" {
		a.Fail()
	}

	if t.Admin != true {
		a.Fail()
	}
}
