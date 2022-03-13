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
	d.heapFile = &dataFile
	d.nextPageID = uint64(nextPageID)
	return d, nil
}

func (d *DiskManager) Open(dataFilePath string) (*DiskManager, error) {
	file, err := os.OpenFile(dataFilePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	return d.new(*file)
}

func (d *DiskManager) AllocatePage() PageID {
	return PageID(d.nextPageID + 1)
}

func (d *DiskManager) ReadPageData(pageID PageID, data []byte) ([]byte, error) {
	offset := os.Getpagesize() * int(pageID)
	_, err := d.heapFile.ReadAt(data, int64(offset))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (d *DiskManager) WritePageData(pageID PageID, data []byte) ([]byte, error) {
	offset := os.Getpagesize() * int(pageID)
	_, err := d.heapFile.WriteAt(data, int64(offset))
	if err != nil {
		return nil, err
	}
	return data, nil
}
