package main

import (
	"io/ioutil"
	"net/http"
)

func main() {
	handler := http.HandlerFunc(handleRequest)
	http.Handle("/photo", handler)
	http.ListenAndServe(":30472", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	fileBytes, err := ioutil.ReadFile("test.png")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
	return
}
