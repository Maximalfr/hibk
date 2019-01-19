package models

type Artist struct {
	Id   int
	Name string
}

type Album struct {
	Id       int
	Name     string
	ArtistId int
	Year     int
}

type Track struct {
	Name        string
	AlbumId     int
	ArtistId    int
	TrackNumber int
	DiscNumber  int
	Genre       string
	Path        string
}

type CompleteTrack struct {
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
