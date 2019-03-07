package wingo

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/valyala/fasthttp"
	"hash"
	"strings"
)

type TokenPayload struct {
	Id     int64             `json:"id"`
	Name   string            `json:"name"`
	Expire int64             `json:"expire"`
	Extra  map[string]string `json:"extra"`
}

var (
	mac   hash.Hash
	codec Codec
)

func init() {
	codec = &JsonCodec{}
}

func SetAuthSecretKey(authSecretKey string) {
	mac = hmac.New(sha256.New, String2Bytes(authSecretKey))
}

func checkAuthorization(ctx *Context) {
	if mac == nil {
		panic("need set authorization secret key at first!")
	}
	token := ctx.Request.Header.Peek("Authorization")
	if token == nil {
		ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.Stop()
		return
	}
	if pass, payload := unmarshalToken(Bytes2String(token)); pass {
		ctx.TokenPayload = payload
		ctx.Next()
	} else {
		ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.Stop()
	}
}

func marshalToken(payload *TokenPayload) string {
	if mac == nil {
		panic("need set authorization secret key at first!")
	}
	// payload
	bytes, err := codec.Marshal(payload)
	CheckError(err)
	payloadBuf := make([]byte, base64.RawURLEncoding.EncodedLen(len(bytes)))
	base64.RawURLEncoding.Encode(payloadBuf, bytes)
	// sign
	mac.Reset()
	mac.Write(payloadBuf)
	sign := mac.Sum(nil)
	signBuf := make([]byte, base64.RawURLEncoding.EncodedLen(len(sign)))
	base64.RawURLEncoding.Encode(signBuf, sign)
	return Bytes2String(payloadBuf) + "/" + Bytes2String(signBuf)
}

func unmarshalToken(token string) (bool, *TokenPayload) {
	if mac == nil {
		panic("need set authorization secret key at first!")
	}
	arr := strings.Split(token, "/")
	if len(arr) != 2 {
		return false, nil
	}
	p := String2Bytes(arr[0])
	// sign
	sign := String2Bytes(arr[1])
	signBuf := make([]byte, base64.RawURLEncoding.DecodedLen(len(arr[1])))
	base64.RawURLEncoding.Decode(signBuf, sign)

	mac.Reset()
	mac.Write(p)
	if !hmac.Equal(signBuf, mac.Sum(nil)) {
		return false, nil
	}
	// payload
	payloadBuf := make([]byte, base64.RawURLEncoding.DecodedLen(len(p)))
	base64.RawURLEncoding.Decode(payloadBuf, p)

	tp := &TokenPayload{}
	err := codec.Unmarshal(payloadBuf, tp)
	CheckError(err)

	return true, tp
}
