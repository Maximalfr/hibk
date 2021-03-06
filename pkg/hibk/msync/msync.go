package msync

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Maximalfr/hibk/pkg/hibk/models"
	"github.com/dhowden/tag"

	"github.com/Maximalfr/hibk/pkg/hibk/database"
	"github.com/Maximalfr/hibk/pkg/hibk/util"
)

// Currents contains artists and albums contained in the database
type Current struct {
	Artists map[string]models.Artist
	Albums  map[string]models.Album
}

type completeTrack struct {
	Title            string
	Album            string
	Artist           string
	AlbumArtist      string
	TrackNumber      int
	TotalTrackNumber int
	DiscNumber       int
	TotalDiscNumber  int
	Year             int
	Genre            string
	Path             string
}

// isValid checks file extension.
// Return true is the file extension is contained in the function
func isValid(name string) bool {
	// Split to get the extension
	sName := strings.Split(name, ".")
	if len(sName) < 2 { // If it's true, the file doesn't have explicit extension
		return false
	}
	// Look the extension
	switch sName[1] {
	case "mp3", "flac":
		return true
	}

	return false

}

// walk finds all files in the directory and adds in files just the valids files.
// This function is needed to work with filepath.Walk. So it returns a WalkFunc
func walk(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isValid(info.Name()) {
			*files = append(*files, path)
		}
		return nil
	}
}

// getFiles calls walk function
func getFiles(path string) (files []string, err error) {
	files = make([]string, 0)
	err = filepath.Walk(path, walk(&files))
	return
}

// artistMapCreation creates a map with the artist name as key and artist id as value
func artistMapCreation(artistsSlice *[]models.Artist) *map[string]int {
	artists := make(map[string]int)

	for _, artist := range *artistsSlice {
		artists[artist.Name] = artist.Id
	}

	return &artists
}

// artistMapCreation creates a map with a special key and album struct as value
// The special key is the concatenation between the album name and the album artist id
func albumMapCreation(albumsSlice *[]models.Album) *map[string]models.Album {
	albums := make(map[string]models.Album)

	for _, album := range *albumsSlice {
		key := album.Name + strconv.Itoa(album.ArtistId)
		albums[key] = album
	}
	return &albums
}

// Sync is the main function of music syncing. Give it a filepath and it will fill the 
// database with all music contained in.
func Sync(filepath string) {
	// loading existing artists and albums
	artists := artistMapCreation(database.GetArtists())
	albums := albumMapCreation(database.GetAlbums())

	files, err := getFiles(filepath)
	db_files := database.GetTracksPath()
	if util.CheckErr("Sync", err, false) {
		return
	}
	for _, file := range files {
		// Processing for the files in the directoy but not in the database
		if in := util.StringInSlice(file, db_files); !in {
			track := tagReading(file)
			if track.Title != "" {
				processing(&track, artists, albums)
			}
		}
	}
}

// tagReading reads tags from an audio file. Returns a completeTrack type
// If an error occured, the function can return a default completeTrack
func tagReading(filepath string) completeTrack {
	file, err := os.Open(filepath)
	util.CheckErr("tagReading", err, true)

	t, err := tag.ReadFrom(file)
	if err != nil {
		if err.Error() == "EOF" { // Skip if the file cannot be read
			log.Printf("%v skipped\n", filepath)
			return completeTrack{}
		} else {
			panic(err)
		}
	}

	trackNumber, totalTrackNumber := t.Track()
	discNumber, totalDiscNumber := t.Disc()

	return completeTrack{
		Title:            t.Title(),
		Album:            t.Album(),
		Artist:           t.Artist(),
		AlbumArtist:      t.AlbumArtist(),
		TrackNumber:      trackNumber,
		TotalTrackNumber: totalTrackNumber,
		DiscNumber:       discNumber,
		TotalDiscNumber:  totalDiscNumber,
		Year:             t.Year(),
		Genre:            t.Genre(),
		Path:             filepath}
}

// processing adds the track to the database. But before that, it searchs the
// artist id and album id for the track. If the album or artist doesn't exists, this function
// handles this issue and fills the database with missing values.
// (see getArtistId)
// (see getAlbumId)
func processing(track *completeTrack, artists *map[string]int, albums *map[string]models.Album) {
	// Artist
	var artistId int
	if track.Artist == "" {
		artistId = getArtistId("Unknow Artist", artists)
	} else {
		artistId = getArtistId(track.Artist, artists)
	}

	// Album
	var albumArtistId int
	if track.AlbumArtist == "" {
		albumArtistId = artistId
	} else {
		albumArtistId = getArtistId(track.AlbumArtist, artists)
	}

	album := models.Album{0, track.Album, albumArtistId, track.Year}

	if album.Name == "" {
		album.Name = "UnknowAlbum"
	}

	albumId := getAlbumId(&album, albums)

	//Track
	// track name isn't empty, it's verified in Sync function
	// Normally, this track ins't in the database, it's also checked in Sync function
	// but if can have the same track with a different path.
	newTrack := models.Track{
		0, // The id is useless for addTrack
		track.Title,
		albumId,
		artistId,
		track.TrackNumber,
		track.DiscNumber,
		track.Genre,
		track.Path}
	database.AddTrack(&newTrack)
}

// getArtistId returns the id for an artist in database. If this artist doesn't exists, the function
// creates it and returns the new id
func getArtistId(name string, artists *map[string]int) int {
	artistId := (*artists)[name]
	if artistId == 0 {
		artistId = database.AddArtist(name)
		(*artists)[name] = artistId
	}
	return artistId
}

// getAlbumId returns the id for an album in database. If this album doesn't exists, the function
// creates it and returns the new id
func getAlbumId(albumModel *models.Album, albums *map[string]models.Album) int {
	album := (*albums)[albumModel.Name+strconv.Itoa(albumModel.ArtistId)]
	if (album != models.Album{}) {
		return album.Id
	}
	id := database.AddAlbum(albumModel)
	albumModel.Id = id
	(*albums)[albumModel.Name+strconv.Itoa(albumModel.ArtistId)] = *albumModel
	return id
}
