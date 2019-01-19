package msync

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"log"

	"../models"
	"github.com/dhowden/tag"

	"../database"
	"../util"
)

type Current struct {
	Artists map[string]models.Artist
	Albums  map[string]models.Album
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

func getFiles(path string) (files []string, err error) {
	files = make([]string, 0)
	err = filepath.Walk(path, walk(&files))
	return
}

func Sync(filepath string) {
	// loading existing artists and albums
	artists := database.GetArtists()
	albums := database.GetAlbums()

	files, err := getFiles(filepath)
	db_files := database.GetTracksPath()
	if util.EH("Sync", err, false) {
		return
	}
	for _, file := range files {
		// Processing for the files in the directoy but not in the database
		if in := util.StringInSlice(file, db_files); !in {
			track := tagReading(file)
			if track.Title != "" {
				processing(&track, &artists, &albums)
			}
		}
	}
}

// Reads tags from an audi file. Return a models.CompleteTrack
// If an error occured, the function can return a default models.CompleteTrack
func tagReading(filepath string) models.CompleteTrack {
	file, err := os.Open(filepath)
	util.EH("tagReading", err, true)

	t, err := tag.ReadFrom(file)
	if err != nil {
		if err.Error() == "EOF" { // Skip if the file cannot be read
			log.Printf("%v skipped\n", filepath)
			return models.CompleteTrack{}
		} else {
			panic(err)
		}
	}

	trackNumber, totalTrackNumber := t.Track()
	discNumber, totalDiscNumber := t.Disc()

	return models.CompleteTrack{
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

func processing(track *models.CompleteTrack, artists *map[string]int, albums *map[string]models.Album) {
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
	// track name isn't empty, it's verify in Sync function
	// Normally, this track ins't in the database, it's also checked in Sync function
	// but if can have the same track with a different path.
	newTrack := models.Track{
		track.Title,
		albumId,
		artistId,
		track.TrackNumber,
		track.DiscNumber,
		track.Genre,
		track.Path}
	database.AddTrack(&newTrack)
}

// Return the id for an artist in database. If this artist doesn't exists, the function
// creates it and return the new id
func getArtistId(name string, artists *map[string]int) int {
	artistId := (*artists)[name]
	if artistId == 0 {
		artistId = database.AddArtist(name)
		(*artists)[name] = artistId
	}
	return artistId
}

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
