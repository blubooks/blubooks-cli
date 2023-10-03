package app

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/blubooks/blubooks-cli/pkg/tools"
	jsoniter "github.com/json-iterator/go"
	"github.com/segmentio/ksuid"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

const (
	TypeGroup    int = 1
	TypeMenuItem int = 2
	TypeLink     int = 3
)

var naviUrlIds map[string]string

type Page struct {
	Id         string  `json:"id,omitempty"`
	Set        bool    `json:"-"`
	Parent     *Page   `json:"-"`
	ParentLink *string `json:"parent,omitempty"`
	ParentId   *string `json:"parentId,omitempty"`
	Level      int     `json:"level,omitempty"`
	Type       int     `json:"type,omitempty"`
	Title      *string `json:"title,omitempty"`
	Link       *string `json:"link,omitempty"`
	DataLink   *string `json:"data"`
	Pages      []Page  `json:"pages,omitempty"`
}

type MetaLinks struct {
	Name       string      `json:"name"`
	Link       string      `json:"link,omitempty"`
	DataLink   string      `json:"data"`
	Id         string      `json:"id,omitempty"`
	ExternLink bool        `json:"extern,omitempty"`
	Navi       string      `json:"navi,omitempty"`
	Links      []MetaLinks `json:"links,omitempty"`
}

type Navi struct {
	Title  string      `json:"title,omitempty"`
	Id     string      `json:"id"`
	Pages  []Page      `json:"pages,omitempty"`
	Header []MetaLinks `json:"header,omitempty"`
	Footer []MetaLinks `json:"footer,omitempty"`
}

func createLink(link *string) *string {
	if link != nil {

		if *link == "README.md" {
			l := "/"
			return &l
		} else if filepath.Base(*link) == "README.md" {
			l := filepath.Dir(*link)
			l = tools.SetFirstLash(l)
			return &l
		}
		l := strings.TrimSuffix(*link, filepath.Ext(*link))
		l = tools.SetLastLash(l)
		l = tools.SetFirstLash(l)

		return &l

	}
	return nil
}

func getUrlId(url string) string {
	if val, ok := naviUrlIds[url]; ok {
		return val
	}
	naviUrlIds[url] = ksuid.New().String()
	return naviUrlIds[url]
}

func list(node ast.Node, initLevel int, page *Page, source *[]byte) {
	level := initLevel
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkStatus(ast.WalkContinue)

		if entering {
			if n.Kind() == ast.KindList {
				level = level + 1

			}
			if n.Kind() == ast.KindListItem {
				if level == initLevel+1 {
					pg := Page{}
					pg.Type = TypeLink
					pg.Level = level
					pg.Parent = page
					pg.ParentId = &page.Id
					pg.ParentLink = createLink(page.Link)
					pg.Id = ksuid.New().String()

					listitemlink(&pg, n.FirstChild(), source)

					list(n, level, &pg, source)

					page.Pages = append(page.Pages, pg)

				}
			}
		} else {
			if n.Kind() == ast.KindList {
				level = level - 1
			}
		}
		var err error
		return s, err
	})
}

func listitemlink(page *Page, node ast.Node, source *[]byte) {

	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkStatus(ast.WalkContinue)

		if entering {
			if n.Kind() == ast.KindLink {
				l := n.(*ast.Link)
				titleStr := string(n.Text([]byte(*source)))
				linkStr := string(l.Destination)
				page.Title = &titleStr
				page.DataLink = &linkStr
				page.Link = createLink(&linkStr)
				page.Id = getUrlId(*page.Link)

			}
		}
		var err error
		return s, err
	})

	if page.Title == nil {
		titleStr := string(node.FirstChild().Text([]byte(*source)))
		page.Title = &titleStr
	}

}

func metalinks(links []MetaLinks) {
	for idx := range links {
		s := &links[idx]

		matched, _ := regexp.MatchString(`^(?:[a-z+]+:)?//`, s.Link)
		if matched {
			s.ExternLink = true
		} else {
			if strings.HasSuffix(s.Link, ".md") {
				s.DataLink = s.Link
				l := createLink(&s.DataLink)
				s.Link = *l
				s.Id = getUrlId(*l)
			}

			if s.Links != nil {
				metalinks(s.Links)
			}

		}

	}
}

func genNavi(filename string) (*Navi, error) {
	naviUrlIds = make(map[string]string)

	isRoot := false

	if filename == "SUMMARY.md" {
		isRoot = true
	}

	if !strings.HasSuffix(filename, "SUMMARY.md") {
		return nil, errors.New("no SUMMERY.md found")
	}

	source, err := os.ReadFile("data/content/" + filename)
	if err != nil {
		return nil, err
	}

	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	var buf bytes.Buffer
	context := parser.NewContext()

	if err := markdown.Convert(source, &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}

	var navi Navi
	navi.Id = getUrlId("README.md")

	if isRoot {
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		metaData := meta.Get(context)

		header := metaData["header"]
		jsonString, err := json.Marshal(header)
		if err == nil {
			var ml []MetaLinks
			json.Unmarshal(jsonString, &ml)
			metalinks(ml)
			navi.Header = ml

		}
		footer := metaData["footer"]
		jsonString, err = json.Marshal(footer)
		if err == nil {
			var ml []MetaLinks
			metalinks(ml)
			json.Unmarshal(jsonString, &ml)
			navi.Footer = ml

		}
	}

	doc := markdown.Parser().Parse(text.NewReader([]byte(source)))
	//doc := goldmark.DefaultParser().Parse(text.NewReader([]byte(source)))
	listLevel := 0

	var entry Page
	titleSet := false
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkStatus(ast.WalkContinue)
		var err error

		if entering {

			if n.Kind() == ast.KindHeading {

				h := n.(*ast.Heading)
				if !titleSet && h.Level == 1 && navi.Title == "" {
					titleSet = true
					navi.Title = string(n.Text([]byte(source)))

				} else if h.Level == 2 {

					if entry.Set {

						navi.Pages = append(navi.Pages, entry)
					}

					titleStr := string(n.Text([]byte(source)))

					entry = Page{
						Id:    ksuid.New().String(),
						Set:   true,
						Level: listLevel,
						Type:  TypeGroup,
						Title: &titleStr,
					}
				}

			} else if n.Kind() == ast.KindThematicBreak {

				if entry.Set {

					navi.Pages = append(navi.Pages, entry)
					//navi.Pages = append(entry.Pages, entry)
					entry = Page{}
				}

			} else if n.Kind() == ast.KindList {
				listLevel = listLevel + 1

			} else if n.Kind() == ast.KindListItem {

				if listLevel == 1 {

					pg := Page{}
					pg.Set = true

					pg.Type = TypeMenuItem
					pg.Level = listLevel
					pg.Parent = &entry
					pg.ParentId = &entry.Id
					listitemlink(&pg, n.FirstChild(), &source)
					list(n, 1, &pg, &source)
					if entry.Type == 1 {
						entry.Pages = append(entry.Pages, pg)
					} else {
						navi.Pages = append(navi.Pages, pg)
					}

				}
			}

		} else {
			if n.Kind() == ast.KindList {
				listLevel = listLevel - 1

			} else if n.Kind() == ast.KindDocument {
				if entry.Set {

					navi.Pages = append(navi.Pages, entry)

				}
			}
		}

		return s, err
	})

	return &navi, nil
}
