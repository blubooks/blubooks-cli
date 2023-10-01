package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/blubooks/blubooks-cli/pkg/goldmark/baseurl"
	"github.com/blubooks/blubooks-cli/pkg/tools"
	replacer "github.com/fundipper/goldmark-replacer"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
)

type App struct {
}

func New() *App {
	return &App{}
}

type PageContent struct {
	Title string `json:"title,omitempty"`
	Html  string `json:"html"`
	Id    string `json:"id"`
}

func Build(dev bool) error {
	navi, err := genNavi()
	if err != nil {
		return err
	}

	naviBytes, err := json.Marshal(navi)
	if err != nil {
		return err
	}

	// CopyAll
	_ = os.RemoveAll("public")
	_ = os.MkdirAll("public", os.ModePerm)
	_ = os.MkdirAll("public/files", os.ModePerm)
	tools.CopyDir("data/content/.data/assets", "public/files")

	if !dev {

		tools.CopyDir("client/default", "public")

		// add BaseUL
		data, err := ioutil.ReadFile("client/default/index.html")
		if err != nil {
			fmt.Println(err)
		}
		str := string(data)
		str = strings.Replace(str, "<!-- ##BASE## -->", "<base href=\"http://localhost:4080/public/\">", 1)
		err = ioutil.WriteFile("public/index.html", []byte(str), 0777)
		if err != nil {
			return err
		}

	}
	// ApiFiles
	_ = os.MkdirAll("public/api/", os.ModePerm)
	err = os.WriteFile("public/api/navi.json", naviBytes, os.ModePerm)
	if err != nil {
		return err
	}
	writeJson("README.md", navi.Id)
	page(navi.Pages)

	return nil

}

func writeJson(filename string, id string) {
	var err error
	content, err := loadMarkdown(filename)
	if err != nil {
		log.Printf("Error in err, page() -> loadMarkdown(): %v", err)
	}
	var c PageContent
	c.Html = content.String()
	c.Id = id
	cJson, err := json.Marshal(c)
	if err != nil {
		log.Printf("Error in err, page() -> loadMarkdown(): %v", err)
	}

	err = os.WriteFile("public/api/"+id+".json", cJson, os.ModePerm)
	if err != nil {
		log.Printf("Error in err, page() -> loadMarkdown(): %v", err)
	}
}

func page(pages []Page) {

	for _, s := range pages {

		if s.Link != nil && s.DataLink != nil {
			writeJson(*s.DataLink, s.Id)

		}
		if len(s.Pages) > 0 {
			page(s.Pages)
		}
	}
}

func loadMarkdown(filename string) (bytes.Buffer, error) {
	var buf bytes.Buffer

	source, err := os.ReadFile("data/content/" + filename)
	if err != nil {
		return buf, err
	}

	/*
		str := string(source)
		re := regexp.MustCompile("page3")
		newStr := re.ReplaceAllString(str, "PAGE3")
	*/

	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
			replacer.NewExtender(
				"(c)", "&copy;",
				"(r)", "&reg;",
				"...", "&hellip;",
				"(tm)", "&trade;",
				"<-", "&larr;",
				"->", "&rarr;",
				"<->", "&harr;",
				"--", "&mdash;",
			),
			baseurl.NewExtender("", "files/")),
	)

	if err := markdown.Convert(source, &buf); err != nil {
		return buf, err
	}
	return buf, nil
}
