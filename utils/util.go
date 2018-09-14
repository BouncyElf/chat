package utils

import (
	"encoding/base64"
	"errors"
	"strconv"

	"github.com/aofei/air"
)

type M map[string]interface{}

func Error(code int, err error) error {
	return &air.Error{
		Code:    code,
		Message: err.Error(),
	}
}

func Success(res *air.Response, data interface{}) error {
	ret := M{}
	ret["code"] = 0
	ret["error"] = ""
	ret["data"] = data
	if data == nil {
		ret["data"] = ""
	}
	return res.JSON(ret)
}

func GetInt(req *air.Request, key string) (int, error) {
	v, ok := req.Params[key]
	if !ok {
		return 0, errors.New("no specific key")
	}
	return strconv.Atoi(v)
}

func GetInt64(req *air.Request, key string) (int64, error) {
	v, ok := req.Params[key]
	if !ok {
		return 0, errors.New("no specific key")
	}
	return strconv.ParseInt(v, 10, 64)
}

func Base64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
