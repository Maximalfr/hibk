package endpoints

import (
	"log"
	"net/http"

	"github.com/Maximalfr/hibk/api/apiutils"
	"github.com/Maximalfr/hibk/api/errorcodes"
	"github.com/Maximalfr/hibk/database"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// getUsername returns the username in the header (jwt)
func getUsername(c *gin.Context) string {
	return c.Writer.Header().Get("username") // With the jwt, the username is in the header
}

func applyUserRoutes(r *gin.RouterGroup) {
	r.POST("/changepwd", changePwd)
	r.GET("/getuserinfo", getUserInfo)
}

func getUserInfo(c *gin.Context) {
	username := getUsername(c)
	c.JSON(200, struct {
		Username string `json:"username"`
	}{username})
}

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
