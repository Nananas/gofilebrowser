package filebrowser

import (
	"bytes"
	"html/template"
	"log"
	"math"
	"strconv"
	"strings"
)

type TemplateData struct {
	Title  string
	Files  []*TemplateFile
	Static string
}

type TemplateFile struct {
	Type int
	Path string
	Name string
	Size string
}

const (
	DIR  = -1
	DOC  = 0
	TXT  = 1
	ARCH = 2
	IMG  = 3
	PDF  = 4
	AUD  = 5
	WORD = 6
	PPT  = 7
	XLS  = 8
	VID  = 9
	CODE = 10
)

const VIEWER_BLOCKS = `
<!DOCTYPE html>
<html>
	<head>
	    <meta charset="UTF-8">
		<title>{{.Title}}</title>

		<link rel="stylesheet" href="{{$.Static}}/css/fontello.css">
		<link rel="stylesheet" href="{{$.Static}}/css/style.css">
		<script>
			(function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
			(i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
			m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
			})(window,document,'script','//www.google-analytics.com/analytics.js','ga');

			ga('create', 'UA-45215610-2', 'auto');
			ga('send', 'pageview');
		</script>
	</head>
	<body>
		<div class="wrapper container">
			<h1> <a href=".." style="float:left; font-size:26px">..</a> {{.Title}} </h1>
			<div class="wrapper"> {{range .Files}}
				{{if not .IsIndex}}
				<div class="block">
					<a {{if .Is 0}} class="icon-doc icon"
					{{else if .Is -1}} class="icon-folder-empty icon"
					{{else if .Is 1}} class="icon-doc-text icon"
					{{else if .Is 2}} class="icon-file-archive icon"
					{{else if .Is 3}} class="icon-file-image icon"
					{{else if .Is 4}} class="icon-file-pdf icon"
					{{else if .Is 5}} class="icon-file-audio icon"
					{{else if .Is 6}} class="icon-file-word icon"
					{{else if .Is 7}} class="icon-file-powerpoint icon"
					{{else if .Is 8}} class="icon-file-excel icon" 
					{{else if .Is 9}} class="icon-file-video icon" 
					{{else if .Is 10}} class="icon-file-code icon" 
					{{end}}
					href="{{.Path}}"><p>{{.Name}}</p>
					<div class="stats">{{.Size}}</div>
					</a>
				</div>
				{{end}} {{end}}
			</div>
		</div>
	</body>
</html>
`

func CreateIndex(data *TemplateData) []byte {

	templ := template.New("index.html")
	t, err := templ.Parse(VIEWER_BLOCKS)
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBuffer(nil)
	// log.Println()

	err = t.Execute(buf, data)

	if err != nil {
		log.Fatal(err)
	}
	// log.Println(buf.String())

	return buf.Bytes()
}

func GetType(filename string) int {

	s := strings.Split(filename, ".")
	ext := s[len(s)-1]
	switch ext {
	case "zip", "rar", "gz", "7z":
		return ARCH
	case "jpg", "jpeg", "gif", "png", "bmp":
		return IMG
	case "pdf":
		return PDF
	case "flac", "mp3":
		return AUD
	case "doc", "docx":
		return WORD
	case "ppt", "pptx":
		return PPT
	case "xls", "xlsx":
		return XLS
	case "go", "conf", "css", "html":
		return CODE
	default:
	}

	if len(s) == 0 {
		return DOC
	}

	return TXT

}

func Human(size int64) string {
	suffixes := []string{
		" B", "KB", "MB", "GB",
	}

	if size == 0 {
		return "0  B"
	}

	b := math.Log(float64(size)) / math.Log(1024)
	ii := int(b)
	s := math.Pow(1000, float64(ii))

	return strconv.Itoa(int(float64(size)/s)) + " " + suffixes[ii]
}

// Helper function for template
//
func (f *TemplateFile) Is(i int) bool {
	return f.Type == i
}

// Helper function for template
//
func (f *TemplateFile) IsIndex() bool {
	return f.Name == "index.html"
}
