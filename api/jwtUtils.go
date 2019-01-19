package api

import (
	"errors"
	"strings"
	"time"

	"./errorcodes"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var signingKey = []byte("keymaker")
var ErrTokenExpired = errors.New("token")

func jwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		// The token isn't in the header
		if len(tokenString) == 0 {
			c.AbortWithStatusJSON(errorcodes.MissingAuthorizationHeader())
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1) // Remove the bearer text
		claims, err := verifyToken(tokenString)
		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					// Token is either expired or not active yet
					c.AbortWithStatusJSON(errorcodes.ErrorTokenExpired())
				} else {
					c.AbortWithStatusJSON(errorcodes.ErrorToken(err.Error()))
				}
			} else {
				c.AbortWithStatusJSON(errorcodes.InternalError(err.Error()))
			}
			return
		}
		name := claims.(jwt.MapClaims)["username"].(string)

		c.Header("username", name)
		c.Next()
	}
}

func getToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().UTC().Add(48 * time.Hour).Unix(), // the session is valid for 48h
	})
	tokenString, err := token.SignedString(signingKey)
	return tokenString, err
}

func verifyToken(tokenString string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}
	return token.Claims, err
}
