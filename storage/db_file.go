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

const FOOTER_SIZE int64 = 4

type DBFile struct {
	Path  string
	File  *os.File
	mutex sync.RWMutex
}

func (f *DBFile) Close() error {
	return f.File.Close()
}

func (f *DBFile) AppendEntry(entry *Entry) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	var offset int64
	if info, err := f.File.Stat(); err != nil {
		return err
	} else {
		offset = info.Size()
	}

	_, err := f.File.WriteAt(entry.Marshal(), offset)
	if err != nil {
		return err
	}

	return f.File.Sync()
}

func (f *DBFile) Scan(key []byte) []byte {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	var offset int64
	if info, err := f.File.Stat(); err != nil {
		panic(err)
	} else {
		offset = info.Size()
	}
	for {
		if offset <= FOOTER_SIZE {
			return nil
		}

		// |  n    |  n      |  2 bytes   | 2 bytes     |
		// |  key  |  value  |  key size  | value size  |
		offset -= 2
		valueSize := readInt16(f.File, offset)
		offset -= 2
		keySize := readInt16(f.File, offset)

		offset -= int64(valueSize + keySize)
		curKey := readBytes(f.File, offset, int64(keySize))
		if bytes.Compare(curKey, key) == 0 {
			return readBytes(f.File, offset+int64(keySize), int64(valueSize))
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
	key   []byte
	value []byte
}

func NewEntry(key, value []byte) *Entry {
	return &Entry{
		key:   key,
		value: value,
	}
}

func (e *Entry) getKeyLength() uint16 {
	return uint16(len(e.key))
}

func (e *Entry) getValueLength() uint16 {
	return uint16(len(e.value))
}

func (e Entry) Marshal() []byte {
	keySize := e.getKeyLength()
	valueSize := e.getValueLength()
	result := make([]byte, int64(keySize+valueSize)+FOOTER_SIZE)

	ret := result

	copy(ret, e.key)
	ret = ret[keySize:]

	copy(ret, e.value)
	ret = ret[valueSize:]

	copy(ret, Uint16ToBytes(keySize))
	ret = ret[2:]

	copy(ret, Uint16ToBytes(valueSize))

	return result
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
