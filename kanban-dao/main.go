package main

import (
	"log"

	db "irpl.com/kanban-dao/db"
	srv "irpl.com/kanban-dao/server"
)

func main() {

	// Initialize database connection
	database := db.GetDB()
	if database == nil {
		log.Fatal("Failed to establish database connection.")
	}

	// Start the web server
	srv.Web()
}
