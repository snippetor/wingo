package wingo

import (
	"github.com/gogo/protobuf/proto"
	"github.com/json-iterator/go"
	"github.com/shamaton/msgpack"
)

var (
	globalCodec Codec
)

func init() {
	globalCodec = &JsonCodec{}
}

type Codec interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	ContentType() string
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

func (j *JsonCodec) ContentType() string {
	return "application/json"
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

func (j *ProtobufCodec) ContentType() string {
	return "application/protobuf"
}

// Msgpack消息协议
type MsgPackCodec struct {
}

func (p *MsgPackCodec) Marshal(v interface{}) ([]byte, error) {
	return msgpack.Encode(v)
}

func (p *MsgPackCodec) Unmarshal(data []byte, v interface{}) error {
	return msgpack.Decode(data, v)
}

func (p *MsgPackCodec) Name() string {
	return "msgpack"
}

func (j *MsgPackCodec) ContentType() string {
	return "application/msgpack"
}
