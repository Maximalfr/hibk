package jwt

import (
	"crypto/rand"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/Maximalfr/hibk/api/errorcodes"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var signingKey = []byte("keymaker")
//var signingKey = keyGenerator(64)
var ErrTokenExpired = errors.New("token")

func JwtMiddleware() gin.HandlerFunc {
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
			if ve, ok := err.(*jwtgo.ValidationError); ok {
				if ve.Errors&(jwtgo.ValidationErrorExpired|jwtgo.ValidationErrorNotValidYet) != 0 {
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
		name := claims.(jwtgo.MapClaims)["username"].(string)

		c.Header("username", name)
		c.Next()
	}
}

func GetToken(username string) (string, error) {
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, jwtgo.MapClaims{
		"username": username,
		"exp":      time.Now().UTC().Add(48 * time.Hour).Unix(), // the session is valid for 48h
	})
	tokenString, err := token.SignedString(signingKey)
	return tokenString, err
}

func verifyToken(tokenString string) (jwtgo.Claims, error) {
	token, err := jwtgo.Parse(tokenString, func(token *jwtgo.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}
	return token.Claims, err
}

func keyGenerator(size int) []byte {
	key := make([]byte, size)
	_, err := rand.Read(key)
	if err != nil {
		log.Println(err)
	}
	return key
}
