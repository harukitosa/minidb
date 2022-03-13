package main

import (
	"log"
	"minidb/db"
)

func main() {
	f := db.DiskManager{}
	file, err := f.Open("healFile.db")
	if err != nil {
		log.Println(err)
	}
	log.Println(file)
}
