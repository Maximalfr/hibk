package api

import (
	"github.com/gin-gonic/gin"
)

// main function to boot up everything
func Run() {
	r := gin.Default()
	r.Use(CORSMiddleware())
	api := r.Group("/api")
	api.Any("/ping", ping)
	applyAuthRoutes(api)

	// Authentification required for this group
	authRequired := api.Group("/a")
	authRequired.Use(jwtMiddleware())
	authRequired.Any("/ping", ping)

	applyUserRoutes(authRequired)

	r.Run("localhost:8081") // listen and serve on 0.0.0.0:8080
}

type loginStruct struct {
	Username string
	Password string
}

func ping(c *gin.Context) {
	c.Writer.WriteString("ping")
}

// Just a cors middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
