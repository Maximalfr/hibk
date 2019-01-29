package api

import (
	"github.com/Maximalfr/hibk/database"
    // "./errorcodes"
	"github.com/Maximalfr/hibk/models"
	"github.com/gin-gonic/gin"
)

type fullJson struct {
	Artists *[]models.Artist `json:"artists"`
	Albums *[]models.Album `json:albums`
	Tracks *[]models.Track `json:tracks`
}

func applyMusicRoutes(r *gin.RouterGroup) {
	r.GET("/tracks", getTracks)
}

func getTracks(c *gin.Context) {
	artists := database.GetArtists()
	albums := database.GetAlbums()
	tracks := database.GetTracks()
	full := fullJson{artists, albums, tracks}

    c.JSON(200, full)
}
