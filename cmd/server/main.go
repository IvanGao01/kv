package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/ivangao01/kv/storage"
)

const (
	ActiveFilePath = "D:\\Workspace\\github.comacaee\\kv"
	Address        = "127.0.0.1:3700"
)

func main() {
	server = &Server{}
	server.activeFile = storage.GetActiveFile(ActiveFilePath)
	http.HandleFunc("/get", logWrapper(get))
	http.HandleFunc("/set", logWrapper(set))

	if err := http.ListenAndServe(Address, nil); err != nil {
		fmt.Printf("ERR: %s\n", err)
	}
}

var server *Server

type Server struct {
	activeFile *storage.DBFile
}

type handlerWithError func(http.ResponseWriter, *http.Request) error

func logWrapper(fn handlerWithError) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		fmt.Printf("%s : %s\n", request.RemoteAddr, request.RequestURI)
		if err := fn(writer, request); err != nil {
			fmt.Printf("ERR: %s", err)
		}
	}
}

func get(writer http.ResponseWriter, request *http.Request) error {
	b, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println()
	}
	vb := server.activeFile.Read(b)
	if vb == nil {
		writer.WriteHeader(http.StatusOK)
		return nil
	}
	writer.WriteHeader(http.StatusOK)
	if _, err = writer.Write(vb); err != nil {
		return err
	}

	return nil
}

func set(writer http.ResponseWriter, request *http.Request) error {
	b, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return err
	}
	if strings.IndexByte(string(b), '=') < 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return errors.New("the correct format of request message is \"key=val\"")
	}
	err = server.activeFile.Write(b)
	if err != nil {
		return err
	}
	writer.WriteHeader(http.StatusOK)
	return nil
}
