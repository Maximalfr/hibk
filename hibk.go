package main

import (
	// "./api"
	"./database"
	"./msync"
)

func main() {
	database.Init()
	msync.Sync(".")
	//api.Run()
}
