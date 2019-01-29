package apiutils

import (
	"log"

	"github.com/Maximalfr/hibk/api/errorcodes"

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
