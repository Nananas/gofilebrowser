package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

func main() {
	log.SetFlags(log.Lshortfile)
	http.HandleFunc("/", handleFiles)
	http.HandleFunc("/static/", handleStatic)
	fmt.Println("Listening from port 9876")

	http.ListenAndServe(":9876", nil)
}

func handleStatic(rw http.ResponseWriter, req *http.Request) {
	log.Println("STATIC")
	log.Println(req.RequestURI)

	// mimetype := http.DetectContentType(data)
	// if strings.Contains(mimetype, "text/html") {
	// mimetype = strings.Replace(mimetype, "text/html", "text/plain", 1)
	// }

	s := strings.Split(req.RequestURI, ".")
	log.Println(s)
	if s[len(s)-1] == "ttf" {
		rw.Header().Add("content-type", "application/x-font-ttf")
	}

	http.ServeFile(rw, req, filepath.Join("../", req.URL.EscapedPath()))
}

func handleFiles(rw http.ResponseWriter, req *http.Request) {
	log.Println("INDEX")
	log.Println(req.RequestURI)
	http.ServeFile(rw, req, filepath.Join(".", req.URL.EscapedPath()))
}
