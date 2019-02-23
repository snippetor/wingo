package wingo

import (
	"encoding/json"
	"github.com/gogo/protobuf/proto"
)

var (
	globalCodec Codec
)

func init() {
	globalCodec = &JsonCodec{}
}

// 数据传输协议，即包体格式
type Codec interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

func SetCodec(c Codec) {
	globalCodec = c
}

// JSON消息协议
type JsonCodec struct {
}

func (j *JsonCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j *JsonCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// Protobuf消息协议
type ProtobufCodec struct {
}

func (p *ProtobufCodec) Marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (p *ProtobufCodec) Unmarshal(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}
