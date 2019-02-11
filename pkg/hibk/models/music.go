package models

type Artist struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Album struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ArtistId int    `json:"artist_id"`
	Year     int    `json:"year, omitempty"`
}

type Track struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	AlbumId     int    `json:"album_id"`
	ArtistId    int    `json:"artist_id"`
	TrackNumber int    `json:"track_number"`
	DiscNumber  int    `json:"disc_number"`
	Genre       string `json:"genre, omitempty"`
	Path        string `json:"-"` // Ignore the path when it is serialized in json
}
