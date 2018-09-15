package handlers

import (
	"github.com/BouncyElf/chat/utils"

	"github.com/aofei/air"
)

func init() {
	air.ErrorHandler = errorHandler
	air.GET("/", indexHandler)
}

func errorHandler(err error, req *air.Request, res *air.Response) {
	e, ok := err.(*air.Error)
	if !ok {
		e = &air.Error{
			Code:    500,
			Message: "Server Internal Error",
		}
		air.ERROR("error", utils.M{
			"err": err.Error(),
		})
	}
	if !res.Written {
		if req.Method == "GET" || req.Method == "HEAD" {
			delete(res.Headers, "ETag")
			delete(res.Headers, "Last-Modified")
		}
		ret := utils.M{}
		ret["data"] = ""
		ret["code"] = e.Code
		ret["error"] = e.Message
		res.StatusCode = e.Code
		res.JSON(ret)
		return
	}

}

func indexHandler(req *air.Request, res *air.Response) error {
	return res.String("hello chat")
}
