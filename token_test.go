package wingo

import (
	"encoding/json"
	"fmt"
	"github.com/json-iterator/go"
	"testing"
	"time"
)

var TestSecretKey = "7npc8w4g7BZsQQTVafadd2zhjRyBhPnv"

func init() {
	SetAuthSecretKey(TestSecretKey)
}

func TestToken(a *testing.T) {
	token := marshalToken(&TokenPayload{Id: 20001, Expire: time.Now().Unix()}, )
	fmt.Println(len(token))

	if ok, t := unmarshalToken(token); !ok {
		a.Fail()
	} else {
		fmt.Println(t)
		if t.Id != 20001 {
			a.Fail()
		}

		if t.Name != "carl" {
			a.Fail()
		}
	}
}

func BenchmarkMarshalToken(b *testing.B) {
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = marshalToken(&TokenPayload{Id: 20001, Name: "carl", Expire: time.Now().Unix()})
	}
	b.StopTimer()
}

func BenchmarkUnmarshalToken(b *testing.B) {
	b.ReportAllocs()
	token := marshalToken(&TokenPayload{Id: 20001, Name: "carl", Expire: time.Now().Unix()})
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = unmarshalToken(token)
	}
	b.StopTimer()
}

func BenchmarkMarshalJSON(b *testing.B) {
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(&TokenPayload{Id: 20001, Name: "carl", Extra: map[string]string{"test": "test"}})
	}
	b.StopTimer()
}

func BenchmarkMarshalJSONGoIter(b *testing.B) {
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = jsoniter.Marshal(&TokenPayload{Id: 20001, Name: "carl", Extra: map[string]string{"test": "test"}})
	}
	b.StopTimer()
}
