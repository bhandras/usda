package main

import (
	"github.com/bhandras/usda"
	"log"
)

func main() {
	var db usda.DB
	if err := db.Read("./data"); err != nil {
		log.Fatalf("Error loading database: %s", err)
	}
}
