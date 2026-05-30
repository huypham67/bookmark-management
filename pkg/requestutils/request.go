package requestutils

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Bind[T any](c *gin.Context) (*T, error) {
	req := new(T)

	if c.Request.Method != http.MethodGet {
		if err := c.ShouldBindJSON(req); err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}
	}

	if err := c.ShouldBindUri(req); err != nil {
		return nil, err
	}

	if err := c.ShouldBindQuery(req); err != nil {
		return nil, err
	}

	if err := c.ShouldBindHeader(req); err != nil {
		return nil, err
	}

	return req, nil
}
