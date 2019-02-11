package apiutils

import (
	"log"

	"github.com/Maximalfr/hibk/pkg/hibk/server/api/errorcodes"

	"github.com/gin-gonic/gin"
)

// GetWithJSON parses the json in the body of the context request. if the json is
// correct for the struct, the struct is filled and the function return false.
// If the json doesn't corresponding to the struct (i), return true.
func GetWithJSON(c *gin.Context, i interface{}) bool {
	err := c.BindJSON(i)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(errorcodes.BadRequest())
		return false
	}
	return true
}
