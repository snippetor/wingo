package wingo

import (
	"fmt"
	"testing"
)

func TestJsonCodec(t *testing.T) {
	c := JsonCodec{}
	bytes, err := c.Marshal(map[string]string{"name": "carl", "id": "12345"})
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	m := make(map[string]string)
	c.Unmarshal(bytes, &m)
	fmt.Println(m)

}

func TestMsgPackCodec(t *testing.T) {
	c := MsgPackCodec{}
	bytes, err := c.Marshal(map[string]string{"name": "carl", "id": "12345"})
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	m := make(map[string]string)
	c.Unmarshal(bytes, &m)
	fmt.Println(m)

}

func BenchmarkJsonCodec(b *testing.B) {
	c := JsonCodec{}
	for i := 0; i < b.N; i++ {
		_, _ = c.Marshal(map[string]string{"name": "carl", "id": "12345"})
	}
}

func BenchmarkMsgPackCodec(b *testing.B) {
	c := MsgPackCodec{}
	for i := 0; i < b.N; i++ {
		_, _ = c.Marshal(map[string]string{"name": "carl", "id": "12345"})
	}
}