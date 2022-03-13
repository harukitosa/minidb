package db

import (
	"os"
)

type DiskManager struct {
	heapFile   *os.File
	nextPageID uint64
}

type PageID uint64

func (d *DiskManager) new(dataFile os.File) (*DiskManager, error) {
	fileStat, err := dataFile.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fileStat.Size()
	nextPageID := fileSize / int64(os.Getpagesize())
	return &DiskManager{heapFile: &dataFile, nextPageID: uint64(nextPageID)}, nil
}

func (d *DiskManager) Open(dataFilePath string) (*DiskManager, error) {
	file, err := os.OpenFile(dataFilePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	return d.new(*file)
}

func (d *DiskManager) AllocatePage(diskManager DiskManager) int64 {
	return 0
}

func (d *DiskManager) ReadPageData(diskManager DiskManager, pageID PageID) ([]byte, error) {
	return nil, nil
}

func (d *DiskManager) WritePageData(diskManager DiskManager, pageID PageID) ([]byte, error) {
	return nil, nil
}
