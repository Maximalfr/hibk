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
	Id			int
	Name        string
	AlbumId     int
	ArtistId    int
	TrackNumber int
	DiscNumber  int
	Genre       string
	Path        string
}
