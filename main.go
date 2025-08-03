package main

import (
	"dokja/db"
	"dokja/routes"
)

func main() {
	db.ConnectToMongo()
	routes.Listen()
}
