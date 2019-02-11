package server

import (
	"github.com/Maximalfr/hibk/pkg/hibk/server/api/endpoints"
	"github.com/Maximalfr/hibk/pkg/hibk/server/api/jwt"
    "github.com/Maximalfr/hibk/pkg/hibk/server/middlewares"
	"github.com/gin-gonic/gin"
)

const baseRoot = "/api/v1"


// Run is the main function to boot up everything
func Run() {
	r := gin.Default()
	r.Use(middlewares.CORS())

    // Public routes
    pcRoute := r.Group(baseRoot)
	endpoints.ApplyAnonRoutes(pcRoute)

    // Private routes
	// Authentification required for this group
	peRoute := r.Group(baseRoot)
	peRoute.Use(jwt.JwtMiddleware())
	endpoints.ApplyAuthRoutes(peRoute)

	r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080
}
