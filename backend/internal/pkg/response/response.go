package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PageResult struct {
	Items    interface{} `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"size"`
}

const (
	OK              = 0
	ErrBadRequest   = 400
	ErrUnauthorized = 401
	ErrForbidden    = 403
	ErrNotFound     = 404
	ErrInternal     = 500
)

var msgMap = map[int]string{
	OK:              "success",
	ErrBadRequest:   "bad request",
	ErrUnauthorized: "unauthorized",
	ErrForbidden:    "forbidden",
	ErrNotFound:     "not found",
	ErrInternal:     "internal server error",
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    OK,
		Message: msgMap[OK],
		Data:    data,
	})
}

func Error(c *gin.Context, code int, msg ...string) {
	message := msgMap[code]
	if len(msg) > 0 {
		message = msg[0]
	}
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func ErrorWithStatus(c *gin.Context, httpStatus int, code int, msg ...string) {
	message := msgMap[code]
	if len(msg) > 0 {
		message = msg[0]
	}
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func Page(c *gin.Context, list interface{}, total int64, page int, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    OK,
		Message: msgMap[OK],
		Data: PageResult{
			Items:    list,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}
