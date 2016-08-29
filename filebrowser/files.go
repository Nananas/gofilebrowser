package filebrowser

import (
	"io/ioutil"
	"path/filepath"
)

func (l *YLocation) ReadToData(config *YConfig) (*TemplateData, error) {

	infos, err := ioutil.ReadDir(l.Watch)
	if err != nil {
		return nil, err
	}

	files := make([]*TemplateFile, 0)

	for _, i := range infos {
		f := &TemplateFile{
			Name: i.Name(),
			Path: filepath.Join(".", i.Name()),
			Size: Human(i.Size()),
			Type: GetType(i.Name()),
		}

		if i.IsDir() {
			f.Type = DIR
		}

		files = append(files, f)
	}

	return &TemplateData{
		Title:  l.Title,
		Files:  files,
		Static: config.Static.Serve,
	}, nil

}
