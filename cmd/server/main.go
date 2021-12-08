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
	server = NewServer()
	server.Open()
	server.Close()
}

var server *Server

type Server struct {
	innerServer http.Server
	activeFile  *storage.DBFile
}

func NewServer() *Server {
	return &Server{
		innerServer: http.Server{
			Addr: Address,
		},
	}
}

func (s *Server) Open() {
	server.activeFile = storage.GetActiveFile(ActiveFilePath)

	mux := http.NewServeMux()
	mux.HandleFunc("/get", logWrapper(get))
	mux.HandleFunc("/set", logWrapper(set))
	s.innerServer.Handler = mux

	if err := s.innerServer.ListenAndServe(); err != nil {
		fmt.Printf("ERR: %s\n", err)
	}
}

func (s *Server) Close() {
	_ = s.activeFile.Close()
	fmt.Println("Server closed")
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
	vb := server.activeFile.Scan(b)
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
	var sepIdx = -1
	if sepIdx = strings.IndexByte(string(b), '='); sepIdx < 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return errors.New("the correct format of request message is \"key=val\"")
	}
	ent := storage.NewEntry(b[:sepIdx], b[sepIdx+1:])
	err = server.activeFile.AppendEntry(ent)
	if err != nil {
		return err
	}
	writer.WriteHeader(http.StatusOK)
	return nil
}
