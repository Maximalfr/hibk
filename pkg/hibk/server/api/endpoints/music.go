package endpoints

import (
	"github.com/Maximalfr/hibk/pkg/hibk/database"
    // "./errorcodes"
	"github.com/Maximalfr/hibk/pkg/hibk/models"
	"github.com/gin-gonic/gin"
)

// fullJson is the json structure that will be sent
type fullJson struct {
	Artists *[]models.Artist `json:"artists"`
	Albums *[]models.Album `json:albums`
	Tracks *[]models.Track `json:tracks`
}

// applyMusicRoutes applies routes for music endpoints
func applyMusicRoutes(r *gin.RouterGroup) {
	r.GET("/tracks", getTracks)
}

// getTracks sends the tracks list to the client
// The track contains also all albums and artists
func getTracks(c *gin.Context) {
	artists := database.GetArtists()
	albums := database.GetAlbums()
	tracks := database.GetTracks()
	full := fullJson{artists, albums, tracks}

    c.JSON(200, full)
}
