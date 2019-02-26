package wingo

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"strings"
	"time"
)

type TokenPayload struct {
	Id     uint32            `json:"id"`
	Name   string            `json:"name"`
	Expire time.Duration     `json:"expire"`
	Extra  map[string]string `json:"extra"`
}

func checkAuthorization(ctx *Context) {
	if authSecretKey == "" {
		Log.E("need set authorization secret key at first!")
		return
	}
	token := ctx.Request.Header.Peek("Authorization")
	if token == nil {
		ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.Stop()
		return
	}
	if pass, payload := unmarshalToken(Bytes2String(token), authSecretKey); pass {
		ctx.TokenPayload = payload
		ctx.Next()
	} else {
		ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.Stop()
	}
}

func marshalToken(payload *TokenPayload, privateSecret string) string {
	// payload
	bytes, err := json.Marshal(payload)
	CheckError(err)
	payloadBuf := make([]byte, base64.RawURLEncoding.EncodedLen(len(bytes)))
	base64.RawURLEncoding.Encode(payloadBuf, bytes)
	// sign
	mac := hmac.New(sha256.New, String2Bytes(privateSecret))
	mac.Write(payloadBuf)
	sign := mac.Sum(nil)
	signBuf := make([]byte, base64.RawURLEncoding.EncodedLen(len(sign)))
	base64.RawURLEncoding.Encode(signBuf, sign)
	return Bytes2String(payloadBuf) + "/" + Bytes2String(signBuf)
}

func unmarshalToken(token, privateSecret string) (bool, *TokenPayload) {
	arr := strings.Split(token, "/")
	if len(arr) != 2 {
		return false, nil
	}
	p := String2Bytes(arr[0])
	// sign
	sign := String2Bytes(arr[1])
	signBuf := make([]byte, base64.RawURLEncoding.DecodedLen(len(arr[1])))
	base64.RawURLEncoding.Decode(signBuf, sign)

	mac := hmac.New(sha256.New, String2Bytes(privateSecret))
	mac.Write(p)
	if !hmac.Equal(signBuf, mac.Sum(nil)) {
		return false, nil
	}
	// payload
	payloadBuf := make([]byte, base64.RawURLEncoding.DecodedLen(len(p)))
	base64.RawURLEncoding.Decode(payloadBuf, p)

	tp := &TokenPayload{}
	err := json.Unmarshal(payloadBuf, tp)
	CheckError(err)

	return true, tp
}
