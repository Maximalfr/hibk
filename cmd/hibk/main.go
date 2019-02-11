package main

import (
	"github.com/Maximalfr/hibk/pkg/hibk/server"
	"github.com/Maximalfr/hibk/pkg/hibk/database"
	//"./msync"
)

func main() {
	database.Init()
	//msync.Sync(".")
	server.Run()
}
