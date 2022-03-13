package main

import (
	"log"
	"minidb/db"
)

func main() {
	dm := db.DiskManager{}
	d, err := dm.Open("healFile.db")
	if err != nil {
		log.Println(err)
	}
	buf := make([]byte, 64)
	buf[0] = 115
	buf[1] = 111
	buf[2] = 109
	buf[3] = 101
	buf[4] = 10
	log.Println(buf)
	_, err = d.WritePageData(db.PageID(3), buf)
	if err != nil {
		log.Println(err)
	}
	bufr := make([]byte, 64)
	log.Println(bufr)
	_, err = d.ReadPageData(db.PageID(1), bufr)
	if err != nil {
		log.Println(err)
	}
	log.Println(bufr)
}
