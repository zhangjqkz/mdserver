package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/russross/blackfriday.v2"
)

const (
	name    = "MarkdownServer"
	version = "0.1.1"
)

var config = struct {
	Addr         string
	RootPath     string
	IndexPaths   string
	TemplatePath string
}{}

func main() {
	flag.StringVar(&config.Addr, "a", ":8080", "http addr")
	flag.StringVar(&config.RootPath, "r", "./", "root path")
	flag.StringVar(&config.IndexPaths, "i", "index.md,README.md,readme.md", "index paths")
	flag.StringVar(&config.TemplatePath, "m", "default", "markdown template path")
	versionFlag := flag.Bool("v", false, "show version")
	flag.Parse()
	rootPath := flag.Arg(0)
	if rootPath != "" {
		config.RootPath = rootPath
	}
	absPath, err := filepath.Abs(config.RootPath)
	if err != nil {
		log.Fatalln(err)
	}
	absPath = strings.ReplaceAll(absPath, "\\", "/")
	config.RootPath = absPath
	if *versionFlag {
		fmt.Printf("%s version %s\n", name, version)
		fmt.Println("happy new year!")
		return
	}
	log.Printf("[info]: root path = %q", config.RootPath)
	log.Printf("[info]: %s listen and serve on %q ...", name, config.Addr)
	if err := http.ListenAndServe(config.Addr, http.HandlerFunc(render)); err != nil {
		log.Fatalln(err)
	}
}

func render(w http.ResponseWriter, r *http.Request) {
	file, err := loadFileByURLPath(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("[info]: url path map: %q => %q\n", r.URL.String(), file.Name())
	buffer, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	ext := filepath.Ext(file.Name())
	switch ext {
	case "", ".md":
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		// CRLF to LF, fix: https://github.com/russross/blackfriday/issues/551
		buffer := bytes.ReplaceAll(buffer, []byte("\r"), []byte(""))
		content := string(blackfriday.Run(buffer, blackfriday.WithExtensions(blackfriday.CommonExtensions)))
		context := struct {
			FileName string
			Content  template.HTML
		}{
			FileName: filepath.Base(file.Name()),
			Content:  template.HTML(content),
		}
		var tpl *template.Template
		if config.TemplatePath == "" {
			err = fmt.Errorf("require markdown template path")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if config.TemplatePath == "default" {
			tpl, err = template.New("default").Parse(defaultTemplate)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			tpl, err = template.ParseFiles(config.TemplatePath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		if err := tpl.Execute(w, context); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case ".htm", ".html":
		w.Header().Set("Content-Type", "text/html")
		w.Write(buffer)
		return
	case ".css":
		w.Header().Set("Content-Type", "text/css")
		w.Write(buffer)
		return
	default:
		w.Write(buffer)
		return
	}
}

func loadFileByURLPath(urlPath string) (*os.File, error) {
	filePath := joinFilePath(config.RootPath, urlPath)
	if strings.HasSuffix(urlPath, "/") {
		findIndexPath := ""
		items := strings.Split(config.IndexPaths, ",")
		for _, item := range items {
			item = strings.TrimSpace(item)
			_, err := os.Stat(joinFilePath(filePath, item))
			if err == nil {
				findIndexPath = item
				break
			}
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		if findIndexPath == "" {
			return nil, fmt.Errorf("index file not found on dir %q", filePath)
		}
		filePath = joinFilePath(filePath, findIndexPath)
	}
	if filepath.Ext(filePath) == "" {
		filePath += ".md"
	}
	return os.Open(filePath)
}

func joinFilePath(a, b string) string {
	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)
	if a == "" && b == "" {
		return ""
	} else if a != "" && b == "" {
		return a
	} else if a == "" && b != "" {
		return b
	}
	return strings.TrimRight(a, "/") + "/" + strings.TrimLeft(b, "/")
}
