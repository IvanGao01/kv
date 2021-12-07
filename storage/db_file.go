package storage

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type DBFile struct {
	Path  string
	File  *os.File
	mutex sync.RWMutex
}

// 每个文件最大8589934592
// | header | entries |
// | uint32 | Entry
func (f *DBFile) Write(b []byte) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	var offset int64
	if info, err := f.File.Stat(); err != nil {
		panic(err)
	} else {
		offset = info.Size()
	}
	var size uint16 = uint16(len(b))
	f.File.WriteAt(Uint16ToBytes(size), offset)
	offset += 2
	f.File.WriteAt(b, offset)
	offset += int64(size)
	f.File.Sync()
}

func (f *DBFile) Read(b []byte) []byte {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	var offset int64
	for {
		var sizeb []byte = make([]byte, 2)
		f.File.ReadAt(sizeb, offset)
		offset += 2
		var datab []byte = make([]byte, BytesToUint16(sizeb))

		f.File.ReadAt(datab, offset)
		offset += int64(BytesToUint16(sizeb))
		if bytes.Compare(bytes.Split(datab, []byte("="))[0], b) == 0 {
			return datab
		}

	}
	// 4 + 2
	// size

}

// 获取ActiveFIle的Offset
//func (f *DBFile) offset(){
//	if f.Offset != 0 {
//		return
//	}
//	var b []byte = make([]byte, 4, 4)
//	f.File.ReadAt(b, 0)
//	if BytesToUint32(b) == 0 {
//		_, err := f.File.WriteAt([]byte{0,0,0,4}, 0)
//		if err != nil {
//			fmt.Println(err)
//		}
//		f.File.Sync()
//	}
//	f.Offset = BytesToUint32(b)
//}

func GetActiveFile(path string) *DBFile {
	var dbFile DBFile
	dbFile.Path = path
	dbFile.File = openActiveFile(path)
	return &dbFile
}

func openActiveFile(path string) *os.File {
	file, err := getActiveFile(path)
	if err != nil {
		fmt.Println(err)
	}

	return file
}

// Active File by active.data
// if return nil
func getActiveFile(path string) (*os.File, error) {
	var fName string
	filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, "avtive.data") {
			fName = path
		}
		return nil
	})

	file, err := os.OpenFile(fName, os.O_RDWR, 0600)
	return file, err
}

// | size | key | value |
// | int16 |   by size    |
type Entry struct {
	size  uint16
	key   string
	value string
}

func BytesToUint32(b []byte) uint32 {
	return uint32(binary.BigEndian.Uint32(b))
}

func BytesToUint16(b []byte) uint16 {
	return uint16(binary.BigEndian.Uint16(b))
}

func Uint16ToBytes(n uint16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(n))
	return buf
}
func Uint32ToBytes(n uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(n))
	return buf
}