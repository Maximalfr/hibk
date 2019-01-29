package main

import (
	"github.com/Maximalfr/hibk/api"
	"github.com/Maximalfr/hibk/database"
	//"./msync"
)

func main() {
	database.Init()
	//msync.Sync(".")
	api.Run()
}
