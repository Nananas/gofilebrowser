package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	. "github.com/nananas/gofilebrowser/filebrowser"

	"github.com/fsnotify/fsnotify"
)

var CONFIG *YConfig
var WG sync.WaitGroup

func main() {
	CONFIG = LoadConfig()

	f, err := os.OpenFile("logfile.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}

	defer f.Close()

	if len(os.Args) > 1 && os.Args[1] == "-d" {
		fmt.Println("Starting in debug mode")
		log.SetFlags(log.Lshortfile)
	} else {
		log.SetOutput(f)
	}

	for _, l := range CONFIG.Locations {
		startAtLocation(l)
	}

	WG.Wait()
}

func startAtLocation(l *YLocation) {
	l.Children = make(map[string]chan bool)

	data, err := l.ReadToData(CONFIG)
	if err != nil {
		// log.Println(err)
		fmt.Println("File " + l.Watch + " does not exist anymore")
		return
	}

	b := CreateIndex(data)
	ioutil.WriteFile(filepath.Join(l.Watch, "./index.html"), b, 0644)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	err = watcher.Add(l.Watch)
	if err != nil {
		fmt.Println("Could not start watcher for " + l.Watch)
		fmt.Println("Does the directory exist?")
		os.Exit(1)
	}

	if l.Recursive {
		infos, err := ioutil.ReadDir(l.Watch)
		if err != nil {
			log.Println(err)
			return
		}

		for _, i := range infos {
			if i.IsDir() {
				ch := make(chan bool)
				l.Children[i.Name()] = ch

				fmt.Println("Recursive: " + filepath.Join(l.Watch, i.Name()))

				sl := &YLocation{
					Recursive:   l.Recursive,
					Watch:       filepath.Join(l.Watch, i.Name()),
					Title:       l.Title,
					Stopchannel: ch,
				}

				startAtLocation(sl)

			}
		}
	}

	WG.Add(1)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// log.Println(event)

				if filepath.Base(event.Name) != "index.html" {
					if l.Recursive {
						CheckNewRecursion(l, filepath.Base(event.Name))
					}

					data, err := l.ReadToData(CONFIG)
					if err != nil {
						// log.Println(data)
						l.Stopchannel <- true
						continue
					}

					b := CreateIndex(data)
					ioutil.WriteFile(filepath.Join(l.Watch, "./index.html"), b, 0644)
				}

			case err := <-watcher.Errors:
				log.Fatal("error:", err)
			case <-l.Stopchannel:
				for _, c := range l.Children {
					c <- true
				}

				return
			}
		}

		watcher.Close()
		WG.Done()
	}()
}

func CheckNewRecursion(l *YLocation, newname string) {
	if _, ok := l.Children[newname]; !ok {

		newL := &YLocation{
			Recursive:   l.Recursive,
			Watch:       filepath.Join(l.Watch, newname),
			Title:       l.Title,
			Stopchannel: make(chan bool),
		}

		startAtLocation(newL)
	}

}

func contains(list []string, e string) bool {
	for _, l := range list {
		if l == e {
			return true
		}
	}

	return false
}
