package api

import (
	_ "context"
	_ "fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Api struct {
	Version string
}

func (A *Api) Getid(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello q1mi!",
	})
}
