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
	"github.com/segmentio/ksuid"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/toc"

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

type PageTocItem struct {
	Title string        `json:"title,omitempty"`
	ID    string        `json:"id,omitempty"`
	Items []PageTocItem `json:"items,omitempty"`
}

type PageContent struct {
	Title string        `json:"title,omitempty"`
	Html  string        `json:"html"`
	TOC   []PageTocItem `json:"toc"`
	Id    string        `json:"id"`
}

type SearchPage struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
	Id    string `json:"id"`
	Path  string `json:"path"`
}

var search map[string]SearchPage

func Build(dev bool) error {
	search = make(map[string]SearchPage)
	navi, err := genNavi()
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
	_ = os.MkdirAll("public/api/", os.ModePerm)

	//writeJson("README.md", navi.Id)
	genPages(navi.Pages)

	if navi.Root != nil {
		genPage(navi.Root)
	}

	var searchResults []SearchPage
	for _, v := range search {
		searchResults = append(searchResults, v)
	}

	searchBytes, err := json.Marshal(searchResults)
	if err != nil {
		return err
	}

	searchId := ksuid.New().String()
	err = os.WriteFile("public/api/"+searchId+".json", searchBytes, os.ModePerm)
	if err == nil {
		navi.SearchId = &searchId
	} else {
		log.Println(err)
	}

	naviBytes, err := json.Marshal(navi)
	if err != nil {
		return err
	}

	// ApiFiles
	err = os.WriteFile("public/api/navi.json", naviBytes, os.ModePerm)
	if err != nil {
		return err
	}

	/*
		for _, n := range navi.Header {
			if len(n.Pages) > 0 {
				page(n.Pages)
			}
		}
		for _, n := range navi.Footer {
			if len(n.Pages) > 0 {
				page(n.Pages)
			}
		}


		for _, n := range navi.Navis {
			writeJson(strings.Replace(n.Summary, "SUMMARY.md", "README.md", 1), n.Id)
			if len(n.Pages) > 0 {
				page(n.Pages)

			}
		}
	*/
	return nil

}

func writeJson(filename string, id string) {
	var err error
	var c PageContent

	err = loadMarkdown(filename, &c)
	if err != nil {
		log.Printf("Error in err, page() -> loadMarkdown(): %v", err)
	}

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

func genPage(s *Page) {

	if s.Link != nil && s.DataLink != nil {
		writeJson(*s.DataLink, s.Id)

		_, ok := search[s.Id]
		if !ok {
			var searchPage SearchPage

			searchPage.Id = s.Id
			searchPage.Title = *s.Title
			searchPage.Path = *s.Link

			source, err := os.ReadFile("data/content/" + *s.DataLink)
			if err == nil {
				markdown := goldmark.New(
					goldmark.WithExtensions(
						meta.Meta,
					),
				)
				doc := markdown.Parser().Parse(text.NewReader([]byte(source)))
				ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
					walk := ast.WalkStatus(ast.WalkContinue)
					var err error
					if entering {
						if n.FirstChild() != nil {
							t := string(n.FirstChild().Text(source))

							searchPage.Text = searchPage.Text + " " + t

						}
					}
					return walk, err
				})
				searchPage.Text = strings.Trim(searchPage.Text, " ")

			}
			search[s.Id] = searchPage

		}

	}

	if len(s.Pages) > 0 {
		genPages(s.Pages)
	}
}

func genPages(pages []Page) {

	for _, s := range pages {

		genPage(&s)
	}
}

func tocElements(items *toc.Items) []PageTocItem {
	var its []PageTocItem
	for _, item := range *items {
		t := PageTocItem{
			Title: string(item.Title),
			ID:    string(item.ID),
		}
		if item.Items != nil {
			t.Items = tocElements(&item.Items)
		}
		its = append(its, t)

	}

	return its

}

func loadMarkdown(filename string, content *PageContent) error {
	var buf bytes.Buffer

	source, err := os.ReadFile("data/content/" + filename)
	if err != nil {
		return err
	}

	markdown := goldmark.New(
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
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

	doc := markdown.Parser().Parse(text.NewReader(source))

	tree, err := toc.Inspect(doc, source, toc.MinDepth(2), toc.MaxDepth(4), toc.Compact(true))
	if err == nil {
		pageTocItems := tocElements(&tree.Items)
		/*
			for _, s := range tree.Items {
				log.Println("title" + string(s.Title))
				for _, s1 := range s.Items {
					log.Println("title1" + string(s1.Title))

				}
			}
		*/

		list := toc.RenderList(tree)

		if list != nil {
			var tocBuf bytes.Buffer
			markdown.Renderer().Render(&tocBuf, source, list)

			//content.TOC = tocBuf.String()
			content.TOC = pageTocItems

		}

	}

	if err := markdown.Convert(source, &buf); err != nil {
		return err
	}

	content.Html = buf.String()

	return nil
}
