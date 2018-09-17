package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
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

func ParseInt(req *air.Request, key string) (int, error) {
	v, ok := req.Params[key]
	if !ok {
		return 0, fmt.Errorf("no specific key: %s.", key)
	}
	return strconv.Atoi(v)
}

func ParseInt64(req *air.Request, key string) (int64, error) {
	v, ok := req.Params[key]
	if !ok {
		return 0, fmt.Errorf("no specific key: %s.", key)
	}
	return strconv.ParseInt(v, 10, 64)
}

func MD5(v string) string {
	h := md5.New()
	h.Write([]byte(v))
	return hex.EncodeToString(h.Sum(nil))
}
