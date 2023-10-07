package app

import (
	"bytes"
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
	ExternLink bool    `json:"extern,omitempty"`
	DataLink   *string `json:"data,omitempty"`
	Pages      []Page  `json:"pages,omitempty"`
}

type Navi struct {
	Title  *string `json:"title,omitempty"`
	Header []Page  `json:"header,omitempty"`
	Footer []Page  `json:"footer,omitempty"`
	Pages  []Page  `json:"pages,omitempty"`
	Root   *Page   `json:"root,omitempty"`
}

func getUrlId(url string) string {
	if val, ok := naviUrlIds[url]; ok {
		return val
	}
	naviUrlIds[url] = ksuid.New().String()
	return naviUrlIds[url]
}
func createLinkString(link *string) *string {
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
func createLink(page *Page, navi *Navi) {
	if page.Link != nil {
		page.DataLink = page.Link

		if page.DataLink != nil {

			matched, _ := regexp.MatchString(`^(?:[a-z+]+:)?//`, *page.DataLink)
			if matched {
				page.ExternLink = true
				page.Link = page.DataLink
				page.DataLink = nil
			} else {

				if strings.HasSuffix(*page.DataLink, ".md") {
					page.Id = getUrlId(*page.Link)
				}
				page.Link = createLinkString(page.DataLink)
				if navi.Root == nil && *page.DataLink == "README.md" {
					navi.Root = page
				}

			}
		}

	}
}

func listitemlink(page *Page, node ast.Node, source *[]byte, navi *Navi) {

	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkStatus(ast.WalkContinue)

		if entering {
			if n.Kind() == ast.KindLink {

				l := n.(*ast.Link)

				titleStr := string(n.Text([]byte(*source)))
				linkStr := string(l.Destination)
				page.Title = &titleStr

				page.Link = &linkStr
				createLink(page, navi)
				/*
					page.DataLink = &linkStr
					if page.DataLink != nil {

						matched, _ := regexp.MatchString(`^(?:[a-z+]+:)?//`, *page.DataLink)
						if matched {
							page.ExternLink = true
							page.Link = page.DataLink
							page.DataLink = nil
						} else {
							if strings.HasSuffix(*page.DataLink, ".md") {
								page.Link = createLinkString(&linkStr)
								page.Id = getUrlId(*page.Link)

							}

						}
					}
				*/

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
func list(node ast.Node, initLevel int, page *Page, source *[]byte, navi *Navi) {
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
					//pg.Id = ksuid.New().String()
					pg.Type = TypeLink
					pg.Level = level
					pg.Parent = page

					listitemlink(&pg, n.FirstChild(), source, navi)
					if page.Id != "" {
						pg.ParentId = &page.Id
						pg.ParentLink = createLinkString(page.Link)

					}

					list(n, level, &pg, source, navi)

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

func metalinks(pages []Page, navi *Navi, parent *Page) {
	for idx := range pages {
		s := &pages[idx]
		/*
			if parent != nil {
				s.Parent = s
				s.ParentId = &s.Id

			}
		*/
		createLink(s, navi)

		if s.Pages != nil {
			metalinks(s.Pages, navi, s)

		}

	}
}

func genNavi() (*Navi, error) {
	var navi Navi
	naviUrlIds = make(map[string]string)

	source, err := os.ReadFile("data/content/SUMMARY.md")
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

	doc := markdown.Parser().Parse(text.NewReader([]byte(source)))
	//doc := goldmark.DefaultParser().Parse(text.NewReader([]byte(source)))
	listLevel := 0
	inNavi := false
	var entry Page

	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkStatus(ast.WalkContinue)
		var err error

		if entering {

			if n.Kind() == ast.KindHeading {

				h := n.(*ast.Heading)

				if h.Level == 1 {
					t := string(n.Text([]byte(source)))
					if !inNavi {
						inNavi = true
						navi.Title = &t
					}

					/*
						if !titleSet && h.Level == 1 && navi.Title == "" {
							titleSet = true
							navi.Title = string(n.Text([]byte(source)))
					*/

				} else if h.Level == 2 {
					if inNavi {
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

				}

			} else if n.Kind() == ast.KindThematicBreak {
				if inNavi {

					if entry.Set {
						navi.Pages = append(navi.Pages, entry)
						entry = Page{}
					}
				}

			} else if n.Kind() == ast.KindList {
				listLevel = listLevel + 1

			} else if n.Kind() == ast.KindListItem {

				if listLevel == 1 {

					if inNavi {

						pg := Page{}
						pg.Set = true

						pg.Type = TypeMenuItem
						pg.Level = listLevel

						if entry.Id != "" {
							//pg.Parent = &entry
							id := entry.Id
							pg.ParentId = &id
						}

						listitemlink(&pg, n.FirstChild(), &source, &navi)

						/*
							if entry.Id != "" {
								pg.Parent = &entry
								pg.ParentId = &entry.Id
							}
						*/

						list(n, 1, &pg, &source, &navi)

						if entry.Type == 1 {
							entry.Pages = append(entry.Pages, pg)
						} else {
							navi.Pages = append(navi.Pages, pg)
						}
					}

				}
			}

		} else {
			if n.Kind() == ast.KindList {
				listLevel = listLevel - 1

			} else if n.Kind() == ast.KindDocument {
				if inNavi {
					if entry.Set {
						navi.Pages = append(navi.Pages, entry)

					}

				}
			}
		}

		return s, err
	})

	if navi.Root == nil {
		pg := Page{}
		//pg.Id = ksuid.New().String()
		pg.Type = TypeLink
		pg.Level = 0
		l := "README.md"
		pg.Link = &l
		createLink(&pg, &navi)
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	metaData := meta.Get(context)

	header := metaData["header"]
	if header != nil {
		jsonString, err := json.Marshal(header)
		if err == nil {
			var ml []Page
			json.Unmarshal(jsonString, &ml)
			navi.Header = ml
			metalinks(navi.Header, &navi, nil)
		}

	}
	footer := metaData["footer"]
	if footer != nil {
		jsonString, err := json.Marshal(footer)
		if err == nil {
			var ml []Page
			json.Unmarshal(jsonString, &ml)
			navi.Footer = ml
			metalinks(navi.Footer, &navi, nil)
		}

	}

	return &navi, nil

}
