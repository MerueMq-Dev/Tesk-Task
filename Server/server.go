package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	files, _ := ioutil.ReadDir("./music")
	for _, f := range files {
		http.HandleFunc(
			strings.ReplaceAll("/"+f.Name(), " ", "_"),
			func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, "music/"+f.Name())
			})

	}
	fmt.Println("Server is listening...")
	http.ListenAndServe("localhost:8181", nil)
	//	http.ListenAndServe("localhost:8181", http.FileServer(http.Dir("music")))
}
