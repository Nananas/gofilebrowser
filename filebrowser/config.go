package filebrowser

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type YConfig struct {
	Static    *YStatic
	Locations map[string]*YLocation
}

type YLocation struct {
	Watch     string
	Recursive bool   "yaml:,omitempty"
	Title     string "yaml:,omitempty"

	Children    map[string]chan bool "yaml:-"
	Stopchannel chan bool            "yaml:-"
}

type YStatic struct {
	Serve string
	Path  string
}

func LoadConfig() *YConfig {
	home_dir := os.Getenv("HOME")
	configfile, err := ioutil.ReadFile(filepath.Join(home_dir, ".config/gofilebrowser.conf"))
	if err != nil {
		log.Fatal(err)
	}

	var yamlconfig YConfig

	err = yaml.Unmarshal(configfile, &yamlconfig)
	if err != nil {
		log.Fatal(err)
	}

	return &yamlconfig

}
