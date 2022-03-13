package main

import (
	"log"
	"minidb/db"
)

func main() {
	// dm := db.DiskManager{}
	// d, err := dm.Open("healFile.db")
	// if err != nil {
	// 	log.Println(err)
	// }
	// buf := []byte("HelloWorld!")
	// _, err = d.WritePageData(db.PageID(0), buf)
	// if err != nil {
	// 	log.Println(err)
	// }
	// bufr := make([]byte, 64)
	// log.Println(bufr)
	// _, err = d.ReadPageData(db.PageID(1), bufr)
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(bufr)

	// NewBufferPoolManager()
	dm := db.DiskManager{}
	dpm := db.NewBufferPoolManager()
	d, err := dm.Open("healFile.db")
	if err != nil {
		log.Println(err)
		return
	}
	dpm.Disk = d
	res, err := dpm.FetchPage(db.PageID(0))
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(res.Page))
	res3, err := dpm.FetchPage(db.PageID(0))
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(res3.Page))
	re, err := dpm.FetchPage(db.PageID(0))
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(re.Page))
	re2, err := dpm.FetchPage(db.PageID(1))
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(re2.Page))
}
