package wingo

import (
	"github.com/gogo/protobuf/proto"
	"github.com/json-iterator/go"
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
	Name() string
}

func SetCodec(c Codec) {
	globalCodec = c
}

// JSON消息协议
type JsonCodec struct {
}

func (j *JsonCodec) Marshal(v interface{}) ([]byte, error) {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(v)
}

func (j *JsonCodec) Unmarshal(data []byte, v interface{}) error {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(data, v)
}

func (j *JsonCodec) Name() string {
	return "json"
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

func (p *ProtobufCodec) Name() string {
	return "protobuf"
}
