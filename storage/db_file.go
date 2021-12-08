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


func (f *DBFile) Write(b []byte) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	var offset int64
	if info, err := f.File.Stat(); err != nil {
		panic(err)
	} else {
		offset = info.Size()
	}
	var size = uint16(len(b))
	f.File.WriteAt(b, offset)
	offset += int64(size)
	f.File.WriteAt(Uint16ToBytes(size), offset)
	offset += 2

	f.File.Sync()
}

func (f *DBFile) Read(b []byte) []byte {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	var offset int64
	if info, err := f.File.Stat(); err != nil {
		panic(err)
	}else {
		offset = info.Size()
	}
	for {
		if offset <= 0 {
			return nil
		}
		offset -= 2
		var sizeb  = make([]byte, 2)
		f.File.ReadAt(sizeb, offset)
		offset -= int64(BytesToUint16(sizeb)) //返回到数据的读取位置
		var datab  = make([]byte, BytesToUint16(sizeb))
		f.File.ReadAt(datab, offset)
		if bytes.Compare(bytes.Split(datab, []byte("="))[0], b) == 0 {
			return datab
		}

	}

}


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
