// Package communicator is a system for reading and writing to http responses
package communicator

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func New(w http.ResponseWriter) *Communicator {
	return NewCommunicatorWithCtx(w, nil)
}

func NewCommunicatorWithCtx(w http.ResponseWriter, ctx interface{}) *Communicator {
	return &Communicator{
		e:       json.NewEncoder(w),
		w:       w,
		context: ctx,
	}
}

type Communicator struct {
	e       *json.Encoder
	w       http.ResponseWriter
	context interface{}
}

func (c *Communicator) With(d interface{}) *Communicator {
	return NewCommunicatorWithCtx(c.w, d)
}

func (c *Communicator) OKf(form string, s ...interface{}) {
	c.OK(fmt.Sprintf(form, s...))
}

func (c *Communicator) OK(s ...interface{}) {
	c.write(response{
		Message: format(s...),
		Error:   false,
		Code:    http.StatusOK,
		Data:    c.context,
	})
}

func (c *Communicator) Failf(form string, s ...interface{}) {
	c.Fail(fmt.Sprintf(form, s...))
}

func (c *Communicator) Fail(s ...interface{}) {
	c.write(response{
		Message: format(s...),
		Error:   true,
		Code:    http.StatusConflict,
		Data:    c.context,
	})
}

func (c *Communicator) Errorf(form string, s ...interface{}) {
	c.Error(fmt.Sprintf(form, s...))

}

func (c *Communicator) Error(s ...interface{}) {
	c.write(response{
		Message: format(s...),
		Error:   true,
		Code:    http.StatusInternalServerError,
		Data:    c.context,
	})
}

func (c *Communicator) write(r response) {
	c.w.WriteHeader(r.Code)

	if err := c.e.Encode(r); err != nil {
		c.e.Encode("Error encoding data")
	}
}

type response struct {
	Message string      `json:"message"`
	Error   bool        `json:"error"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

func format(s ...interface{}) string {
	var message string
	if len(s) != 0 {
		message = fmt.Sprint(s...)
	}

	return message
}
