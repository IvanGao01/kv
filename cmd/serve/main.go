package main

import (
	"github.com/ivangao01/kv/storage"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	ActiveFilePath = "D:\\Workspace\\github.comacaee\\kv"
	Address        = "127.0.0.1:3700"
)

func main() {

	Open()

}

var serve *Serve

type Serve struct {
	address    string // default address is 127.0.0.1
	handler    *http.HandlerFunc
	activeFile *storage.DBFile
}

func Open() {
	serve = &Serve{}
	serve.activeFile = storage.GetActiveFile(ActiveFilePath)
	serve.address = Address
	http.HandleFunc("/get", Get)
	http.HandleFunc("/set", Set)
	http.ListenAndServe(serve.address, nil)

}

func Get(writer http.ResponseWriter, request *http.Request) {
	b, _ := ioutil.ReadAll(request.Body)
	kvb := serve.activeFile.Read(b)
	writer.Write(kvb)


}

func Set(writer http.ResponseWriter, request *http.Request) {
	b, _ := ioutil.ReadAll(request.Body)
	keyval := strings.Split(string(b), "=")
	if len(keyval) == 2 { //符合格式，可以写入
		serve.activeFile.Write(b)
	}

}
