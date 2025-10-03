package api

import (
	"net/http"
)

func TypedWithBody[REQ any, RESP any](next func(c *Context, r REQ) (*RESP, error)) HandlerFunc {
	return func(c *Context) error {
		var v REQ
		if err := c.ParseBody(&v); err != nil {
			return err
		}

		if sanitizable, ok := interface{}(&v).(BodySanitize); ok {
			sanitizable.Sanitize()
		}

		if validatable, ok := interface{}(&v).(BodyValidate); ok {
			if err := validatable.Validate(); err != nil {
				return err
			}
		}

		resp, err := next(c, v)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, APIResponse[*RESP]{
			Success: true,
			Data:    resp,
		})
	}
}

func Typed[RESP any](next func(c *Context) (*RESP, error)) HandlerFunc {
	return func(c *Context) error {
		resp, err := next(c)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, APIResponse[*RESP]{
			Success: true,
			Data:    resp,
		})
	}
}

type BodyValidate interface {
	Validate() error
}

type BodySanitize interface {
	Sanitize()
}
