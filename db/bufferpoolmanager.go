package db

import (
	"log"
	"os"
)

type Page []byte
type BufferID uint64

type Buffer struct {
	PageID  PageID
	Page    Page
	isDirty bool
}

type Frame struct {
	usageCount uint64
	buffer     Buffer
}

type BufferPool struct {
	buffers      []Frame
	nextVictimID BufferID
}

type BufferPoolManager struct {
	Disk      *DiskManager
	Pool      BufferPool
	pageTable map[PageID]BufferID
}

func NewBufferPoolManager() BufferPoolManager {
	return BufferPoolManager{
		Disk:      &DiskManager{},
		Pool:      BufferPool{},
		pageTable: map[PageID]BufferID{},
	}
}

func (b *BufferPool) size() int {
	return len(b.buffers)
}

func (b *BufferPool) evict() (BufferID, error) {
	victimID := 0
	consecutivePinned := 0
	log.Println(b)
	for {
		nextVictimID := b.nextVictimID
		if len(b.buffers) <= int(nextVictimID) {
			break
		}
		frame := b.buffers[nextVictimID]

		if frame.usageCount == 0 {
			victimID = int(b.nextVictimID)
			break
		}
		// 参照カウンタ作成する
		frame.usageCount -= 1
		consecutivePinned += 1

		b.nextVictimID = b.incrementID(b.nextVictimID)
	}
	return BufferID(victimID), nil
}

func (b *BufferPool) incrementID(bufferID BufferID) BufferID {
	return (bufferID + 1) % BufferID(b.size())
}

func (b *BufferPoolManager) FetchPage(pageID PageID) (Buffer, error) {
	bufferID, _ := b.pageTable[pageID]
	// if ok {
	// 	var frame Frame
	// 	if len(b.Pool.buffers) <= int(bufferID) {
	// 		frame = Frame{}
	// 	} else {
	// 		frame = b.Pool.buffers[bufferID]
	// 	}
	// 	frame.usageCount += 1
	// 	return frame.buffer, nil
	// } else {
	bufferID, err := b.Pool.evict()
	if err != nil {
		log.Println(err)
	}

	var frame Frame
	if len(b.Pool.buffers) <= int(bufferID) {
		frame = Frame{}
	} else {
		frame = b.Pool.buffers[bufferID]
	}
	evictPageID := frame.buffer.PageID
	buffer := frame.buffer
	if buffer.isDirty {
		b.Disk.WritePageData(evictPageID, buffer.Page)
	}
	buffer.PageID = pageID
	buffer.isDirty = false
	buffer.Page = make([]byte, os.Getpagesize())
	buffer.Page, err = b.Disk.ReadPageData(pageID, buffer.Page)
	if err != nil {
		return buffer, err
	}
	frame.usageCount = 1
	delete(b.pageTable, evictPageID)
	b.pageTable[pageID] = bufferID
	return buffer, nil
	// }
}

func (b *BufferPoolManager) CreatePage() (*Buffer, error) {
	bufferID, err := b.Pool.evict()
	if err != nil {
		return nil, err
	}
	frame := b.Pool.buffers[bufferID]
	evictPageID := frame.buffer.PageID

	buffer := frame.buffer
	if buffer.isDirty {
		b.Disk.WritePageData(evictPageID, buffer.Page)
	}
	pageID := b.Disk.AllocatePage()

	buffer.PageID = pageID
	buffer.isDirty = true
	frame.usageCount = 1

	page := frame.buffer
	delete(b.pageTable, evictPageID)
	b.pageTable[pageID] = bufferID
	return &page, nil
}

func (b *BufferPoolManager) flush() {
	for pageID, bufferID := range b.pageTable {
		frame := b.Pool.buffers[bufferID]
		page := frame.buffer.Page
		b.Disk.WritePageData(pageID, page)
		frame.buffer.isDirty = false
	}
}
