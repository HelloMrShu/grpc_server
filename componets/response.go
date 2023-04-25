package componets

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	code    int
	message string
	content gin.H
	attach  gin.H
}

func NewResponse() (r *Response) {
	r = &Response{
		code:    0,
		message: "",
		content: make(map[string]interface{}),
		attach:  make(map[string]interface{}),
	}

	return
}

func (r *Response) Code(code int) *Response {
	r.code = code

	return r
}

func (r *Response) Message(msg string) *Response {
	r.message = msg

	return r
}

func (r *Response) Append(key string, value interface{}) *Response {
	r.content[key] = value

	return r
}

func (r *Response) Error(msg string) *Response {
	r.Code(-1).
		Message(msg)

	return r
}

func (r *Response) JSON(c *gin.Context) {
	resp := gin.H{
		"code":    r.code,
		"message": r.message,
		"content": r.content,
	}

	for k, v := range r.attach {
		resp[k] = v
	}

	c.JSON(http.StatusOK, resp)
}
