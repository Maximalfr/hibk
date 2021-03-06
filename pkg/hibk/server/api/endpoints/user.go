package endpoints

import (
	"log"
	"net/http"

	"github.com/Maximalfr/hibk/pkg/hibk/server/api/apiutils"
	"github.com/Maximalfr/hibk/pkg/hibk/server/api/errorcodes"
	"github.com/Maximalfr/hibk/pkg/hibk/database"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// getUsername returns the username in the header (jwt)
func getUsername(c *gin.Context) string {
	username, _ := c.Get("username") // The username is set in the context with the jwt middleware
	return  username.(string)		// Shouldn't be empty
}

// applyUserRoutes applies routes for user endpoints
func applyUserRoutes(r *gin.RouterGroup) {
	r.PUT("/password", changePwd)
	r.GET("/user", getUserInfo)
}

// getUserInfo sends user info to the client
func getUserInfo(c *gin.Context) {
	username := getUsername(c)
	c.JSON(200, struct {
		Username string `json:"username"`
	}{username})
}

// changePwd changes the password for the given user
func changePwd(c *gin.Context) {
	username := getUsername(c)
	cps := struct {
		Password    string
		NewPassword string
	}{}

	if !apiutils.GetWithJSON(c, &cps) {
		return // Not a json
	}

	if same, _ := checkPassword(username, cps.Password); same {
		newPassword, err := bcrypt.GenerateFromPassword([]byte(cps.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Println("changePwd: " + err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		err = database.ChangePassword(username, string(newPassword))
		if err != nil {
			log.Println(err)
			return
		}
		c.JSON(errorcodes.OK())
	} else {
		c.AbortWithStatusJSON(errorcodes.BadPassword())
	}
}
