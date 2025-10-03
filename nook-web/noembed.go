//go:build !embedweb
// +build !embedweb

package orbitweb

import (
	"fmt"
	"net/http"
)

func NewHandler() (http.Handler, error) {
	return nil, fmt.Errorf("not implemented")
}
