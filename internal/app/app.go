package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/blubooks/blubooks-cli/internal/app/models"
	"github.com/blubooks/blubooks-cli/pkg/goldmark/baseurl"
	"github.com/segmentio/ksuid"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/toc"
	"google.golang.org/protobuf/proto"

	embed "github.com/13rac1/goldmark-embed"
	d2 "github.com/FurqanSoftware/goldmark-d2"
	katex "github.com/FurqanSoftware/goldmark-katex"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/blubooks/blubooks-cli/pkg/tools"
	replacer "github.com/fundipper/goldmark-replacer"
	pikchr "github.com/jchenry/goldmark-pikchr"
	attributes "github.com/mdigger/goldmark-attributes"
	img64 "github.com/tenkoh/goldmark-img64"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"go.abhg.dev/goldmark/frontmatter"
	"go.abhg.dev/goldmark/mermaid"
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
type PagePDF struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
	Id    string `json:"id"`
	Path  string `json:"path"`
}
type PDF struct {
	Title    string `json:"title,omitempty"`
	SubTtile string `json:"subTitle,omitempty"`
	Text     string `json:"text,omitempty"`
	Pages    string `json:"page"`
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
	_ = os.MkdirAll("public/files/print", os.ModePerm)
	tools.CopyDir("data/content/.data/assets", "public/files")
	tools.CopyDir("data/content/.data/print/default", "public/files/print")

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

	elliot := &models.Person{
		Name: "Elliot",
		Age:  24,
	}
	data, err := proto.Marshal(elliot)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	fmt.Println("dddddddddddddddddd", data)

	newElliot := &models.Person{}
	err = proto.Unmarshal(data, newElliot)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
	fmt.Println(newElliot.GetAge())
	fmt.Println(newElliot.GetName())
	/*
		b := (*[2]uint8)(data)

		outf, _ := os.Create("public/api/data.bin")
		// ApiFiles
		err = binary.Write(outf, binary.LittleEndian, b)

		if err != nil {
			return err
		}
		outf.Close()

	*/
	//7s := string(data[:])
	//log.Println("DDDDDDDDDDDDDDDDDDDD", s)
	err = os.WriteFile("public/api/data.bin", data, os.ModePerm)
	if err != nil {
		return err
	}

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

	markdown := goldmark.New(attributes.Enable,
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithExtensions(
			extension.Table,
			extension.Strikethrough,
			extension.TaskList,
			extension.GFM,
			extension.DefinitionList,
			extension.Typographer,
			extension.Footnote,
			//meta.Meta,
			emoji.Emoji,
			embed.New(),
			&frontmatter.Extender{},
			&mermaid.Extender{},
			&pikchr.Extender{},
			//latex.NewLatex(),
			&d2.Extender{},
			&katex.Extender{},
			img64.NewImg64(),
			//&mathjax.MathJax.Extend(),

			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
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
