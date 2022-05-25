package api

import (
	"github.com/gin-gonic/gin"
)

func ErrResponse(c *gin.Context, code int, err error) {
	c.IndentedJSON(code, ErrorResponse{err.Error()})
}
