package apiutils

import (
	"log"

	"../errorcodes"

	"github.com/gin-gonic/gin"
)

func GetWithJSON(c *gin.Context, i interface{}) bool {
	err := c.BindJSON(i)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(errorcodes.BadRequest())
		return false
	}
	return true
}
