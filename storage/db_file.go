package storage

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const KEY_LENGTH int64 = 2

type DBFile struct {
	Path  string
	File  *os.File
	mutex sync.RWMutex
}

func (f *DBFile) Write(entry []byte) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	var offset int64
	if info, err := f.File.Stat(); err != nil {
		return err
	} else {
		offset = info.Size()
	}
	var size = uint16(len(entry))
	_, err := f.File.WriteAt(entry, offset)
	if err != nil {
		return err
	}
	offset += int64(size)
	_, err = f.File.WriteAt(Uint16ToBytes(size), offset)
	if err != nil {
		return err
	}
	offset += KEY_LENGTH

	f.File.Sync()
	return nil
}

func (f *DBFile) Read(key []byte) []byte {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	var offset int64
	if info, err := f.File.Stat(); err != nil {
		panic(err)
	} else {
		offset = info.Size()
	}
	for {
		if offset <= KEY_LENGTH {
			return nil
		}
		offset -= KEY_LENGTH
		size := int64(readInt16(f.File, offset))
		offset -= size
		data := readBytes(f.File, offset, size)
		sepIdx := bytes.IndexByte(data, '=')
		if sepIdx > 0 {
			if bytes.Compare(data[:sepIdx], key) == 0 {
				return data[sepIdx+1:]
			}
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

func readInt16(reader io.ReaderAt, offset int64) int16 {
	b := make([]byte, 2)
	if _, err := reader.ReadAt(b, offset); err != nil {
		fmt.Printf("ERR while reading int16: %s\n", err)
		return -1
	}
	return int16(binary.BigEndian.Uint16(b))
}

func readBytes(reader io.ReaderAt, offset int64, size int64) []byte {
	b := make([]byte, size)
	if _, err := reader.ReadAt(b, offset); err != nil {
		fmt.Printf("ERR while reading bytes: %s\n", err)
		return b
	}
	return b
}

func BytesToUint32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

func BytesToUint16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

func Uint16ToBytes(n uint16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, n)
	return buf
}
func Uint32ToBytes(n uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, n)
	return buf
}
