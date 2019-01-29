package api

import (
	"errors"
	"log"
	"strings"

	"github.com/Maximalfr/hibk/database"
	"github.com/Maximalfr/hibk/api/apiutils"
	"github.com/Maximalfr/hibk/api/errorcodes"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func applyAuthRoutes(r *gin.RouterGroup) {
	r.POST("/auth", authenticate)
	r.POST("/register", register)
}

func authenticate(c *gin.Context) {
	var ar = struct {
		Username string
		Password string
	}{}

	if !apiutils.GetWithJSON(c, &ar) {
		return // Not a json
	}

	samePwd, err := checkPassword(ar.Username, ar.Password)
	if err != nil { // TODO: Handle the future custom error when the user doesn't exists
		if err == database.ErrDatabaseNotResponding {
			c.AbortWithStatusJSON(errorcodes.InternalError(err.Error()))
		} else {
			log.Println(err)
			c.AbortWithStatusJSON(errorcodes.BadCredentials())
		}
		return
	}

	if samePwd { // If it's the same password
		token, err := getToken(ar.Username) // Generate the jwt
		if err != nil {
			c.AbortWithStatusJSON(errorcodes.InternalError(err.Error()))
			log.Println("Error generating JWT token: " + err.Error())

		} else { // Send the confirmation
			c.Header("Authorization", "Bearer "+token)
			//c.JSON(errorcodes.OK())
			// Needed because can't read the auth header with cors
			var jwt = struct{ Jwt_token string `json:"jwt"` }{token}
			c.JSON(200, jwt)
		}

	} else { // Wrong password
		c.AbortWithStatusJSON(errorcodes.BadPassword())
	}
}

func register(c *gin.Context) {
	var newUser = struct {
		Username string
		Password string
	}{}

	if !apiutils.GetWithJSON(c, &newUser) {
		return // Not a json
	}

	password, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("register: " + err.Error())
		c.AbortWithStatusJSON(errorcodes.InternalError(err.Error()))
		return
	}
	err = database.RegisterUser(newUser.Username, string(password))
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1062") {
			c.AbortWithStatusJSON(errorcodes.UsernameAlreadyExists())
			return
		}
		c.AbortWithStatusJSON(errorcodes.InternalError(err.Error()))
		return
	}
	c.JSON(errorcodes.OK())
}

func checkPassword(username string, password string) (bool, error) {

	user, err := database.GetUser(username)
	if err != nil {
		return false, err
	}
	if len(user.Username) == 0 || len(user.Password) == 0 {
		return false, errors.New("This user doesn't exists") // TODO: Make a custom error
	}

	res := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if res == nil {
		return true, nil
	} else if res == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	return false, res //became an error
}
