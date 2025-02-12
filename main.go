package main

import (
	"url-shortener/config"
	"url-shortener/routes"
)

func main() {
	config.ConnectDatabase()

	r := routes.SetupRouter()
	r.Run(":9080")
}
