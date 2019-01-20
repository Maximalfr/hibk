package database

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"

	"../models"
	"../util"

	_ "github.com/go-sql-driver/mysql"
)

// Creates all tables for music data
func initMusic(db *sql.DB) {
	artists := `CREATE TABLE IF NOT EXISTS artists(
					id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
					name VARCHAR(60) NOT NULL UNIQUE
				)`

	albums := `CREATE TABLE IF NOT EXISTS albums(
					id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
					name VARCHAR(60) NOT NULL,
					album_artist_id INT UNSIGNED NOT NULL,
					year INT(4),
					UNIQUE INDEX (name, album_artist_id),
					FOREIGN KEY (album_artist_id) REFERENCES artists(id)
				)`

	tracks := `CREATE TABLE IF NOT EXISTS tracks(
					id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
					tracknumber INT(3),
					name VARCHAR(60) NOT NULL,
					album_id INT UNSIGNED NOT NULL,
					artist_id INT UNSIGNED NOT NULL,
					disc INT(2) UNSIGNED NOT NULL DEFAULT 1,
					genre VARCHAR(60),
					keyp VARCHAR(32) NOT NULL UNIQUE,
					path VARCHAR(300) NOT NULL,
					UNIQUE INDEX (name, album_id, artist_id, keyp),
					FOREIGN KEY (album_id) REFERENCES albums(id),
					FOREIGN KEY (artist_id) REFERENCES artists(id)
		)`

	inits := []string{artists, albums, tracks}

	for _, ex := range inits {
		_, err := db.Exec(ex)
		util.EH("initMusic", err, false)
	}
}

// Add an artist it the database
func AddArtist(name string) int {
	db, err := open()
	util.EH("AddArtist", err, true)
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO artists(name) VALUES(?)")
	util.EH("AddArtist", err, true)

	res, err := stmt.Exec(name)
	util.EH("AddArtist", err, true)

	lastID, err := res.LastInsertId()
	util.EH("AddArtist", err, true)

	return int(lastID)
}

// Add an album it the database
func AddAlbum(album *models.Album) int {
	db, err := open()
	util.EH("AddAlbum", err, true)
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO albums(name, album_artist_id, year) VALUES(?, ?, ?)")
	util.EH("AddAlbum", err, true)

	res, err := stmt.Exec(album.Name, album.ArtistId, album.Year)
	util.EH("AddAlbum", err, true)

	lastID, err := res.LastInsertId()
	util.EH("AddAlbum", err, true)

	return int(lastID)
}

// Add an album it the database
func AddTrack(track *models.Track) {
	db, err := open()
	util.EH("AddTrack", err, true)

	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO tracks values(NULL, ?, ?, ?, ?, ?, ?, ?, ?)")
	util.EH("AddTrack", err, true)

	h := md5.New()
	h.Write([]byte(track.Path))
	hash := hex.EncodeToString(h.Sum(nil))

	_, err = stmt.Exec(track.TrackNumber,
		track.Name,
		track.AlbumId,
		track.ArtistId,
		track.DiscNumber,
		track.Genre,
		hash,
		track.Path)

	util.EH("AddTrack", err, true)
}

// Get Artists list contained in the database
func GetArtists() *[]models.Artist {
	var artist models.Artist
	var artists []models.Artist

	db, err := open()
	defer db.Close()
	util.EH("GetArtists", err, true)

	rows, err := db.Query("SELECT id, name FROM artists")
	util.EH("GetArtists", err, true)

	for rows.Next() {
		err = rows.Scan(&artist.Id, &artist.Name)
		util.EH("GetArtists", err, true)
		artists = append(artists, artist)
	}
	err = rows.Err()
	util.EH("GetArtists", err, true)

	return &artists
}

// Get the albums list contained in the database
func GetAlbums() *[]models.Album{
	var albums []models.Album
	var album models.Album

	db, err := open()
	defer db.Close()
	util.EH("GetAlbums", err, true)

	rows, err := db.Query("SELECT * FROM albums")

	util.EH("GetAlbums", err, true)

	for rows.Next() {
		err = rows.Scan(&album.Id,
			&album.Name,
			&album.ArtistId,
			&album.Year)
		util.EH("GetAlbums", err, true)
		albums = append(albums, album)
	}

	err = rows.Err()
	util.EH("GetAlbums", err, true)
	return & albums
}

// Get all tracks from database
func GetTracks() *[]models.Track {
	var tracks []models.Track
	var track models.Track
	var skip string // Skip the keyp value

	db, err := open()
	defer db.Close()
	util.EH("GetTracks", err, true)

	rows, err := db.Query("SELECT * FROM tracks")

	for rows.Next() {
		err = rows.Scan(&track.Id,
			&track.TrackNumber,
			&track.Name,
			&track.AlbumId,
			&track.ArtistId,
			&track.DiscNumber,
			&track.Genre,
			&skip,
			&track.Path)
		util.EH("GetTracks", err, true)
		tracks = append(tracks, track)
	}

	return &tracks
}

// Returns all filepaths for tracks
func GetTracksPath() []string {
	db, err := open()
	defer db.Close()
	util.EH("GetTracksPath", err, true)

	// Counts the number of rows for the slice initialisation (paths)
	var count int
	err = db.QueryRow("SELECT count(*) FROM tracks").Scan(&count)
	util.EH("GetTracksPath", err, true)

	var path string
	paths := make([]string, count)

	// Selects path for each track and put them in the paths slice
	rows, err := db.Query("SELECT path FROM tracks")
	util.EH("GetTracksPath", err, true)
	defer rows.Close()

	i := 0
	for rows.Next() {
		err := rows.Scan(&path)
		util.EH("GetTracksPath", err, false)
		paths[i] = path
		i++
	}
	err = rows.Err()
	util.EH("GetTracksPath", err, true)

	return paths
}
