package endpoints

import "github.com/gin-gonic/gin"

func ApplyAnonRoutes(r *gin.RouterGroup) {
    applyAuthRoutes(r)
}

func ApplyAuthRoutes(r *gin.RouterGroup) {
    applyUserRoutes(r)
    applyMusicRoutes(r)
}
