package endpoints

import (
	"errors"
	"log"
	"strings"

	"github.com/Maximalfr/hibk/pkg/hibk/server/api/apiutils"
	"github.com/Maximalfr/hibk/pkg/hibk/server/api/errorcodes"
	"github.com/Maximalfr/hibk/pkg/hibk/server/api/jwt"
	"github.com/Maximalfr/hibk/pkg/hibk/database"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserNotExist = errors.New("user doesn't exist")


func applyAuthRoutes(r *gin.RouterGroup) {
	r.POST("/login", authenticate)
	r.POST("/register", register)
}

// authenticate authenticates an user and if the credentials are goods, generates a
// a jwt and sends it to the client
func authenticate(c *gin.Context) {
	var ar = struct {
		Username string
		Password string
	}{}

	if !apiutils.GetWithJSON(c, &ar) {
		return // Not a json
	}

	samePwd, err := checkPassword(ar.Username, ar.Password)
	if err != nil {
		if err == database.ErrDatabaseNotResponding {
			c.AbortWithStatusJSON(errorcodes.InternalError(err.Error()))
		} else {
			log.Println(err)
			c.AbortWithStatusJSON(errorcodes.BadCredentials())
		}
		return
	}

	if samePwd { // If it's the same password
		token, err := jwt.GetToken(ar.Username) // Generate the jwt
		if err != nil {
			c.AbortWithStatusJSON(errorcodes.InternalError(err.Error()))
			log.Println("Error generating JWT token: " + err.Error())

		} else { // Send the confirmation
			c.Header("Authorization", "Bearer "+token)
			c.JSON(errorcodes.OK())
			// Needed because can't read the auth header with cors
			// var jwt = struct {
			// 	Jwt_token string `json:"jwt"`
			// }{token}
			// c.JSON(200, jwt)
		}

	} else { // Wrong password
		c.AbortWithStatusJSON(errorcodes.BadPassword())
	}
}

// register registers a new user for hibk
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

// checkPassword retrieves the user's password from the database and checks if
// it matches the password sent by the client.
func checkPassword(username string, password string) (bool, error) {

	user, err := database.GetUser(username)
	if err != nil {
		return false, err
	}
	if len(user.Username) == 0 || len(user.Password) == 0 {
		return false, ErrUserNotExist
	}

	res := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if res == nil {
		return true, nil
	} else if res == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	return false, res //became an error
}
