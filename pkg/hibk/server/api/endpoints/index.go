package endpoints
// index.go applies all routes for endpoints. It's easier to do this than calling all other
// apply functions.


import "github.com/gin-gonic/gin"


// ApplyAnonRoutes applies routes for endpoints for which authentication is not required.
func ApplyAnonRoutes(r *gin.RouterGroup) {
    applyAuthRoutes(r)
}

// ApplyAuthRoutes applies routes for endpoints for which authentication is required.
func ApplyAuthRoutes(r *gin.RouterGroup) {
    applyUserRoutes(r)
    applyMusicRoutes(r)
}
