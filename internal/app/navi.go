package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/blubooks/blubooks-cli/pkg/tools"
	"github.com/segmentio/ksuid"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

const (
	TypeGroup    int = 1
	TypeMenuItem int = 2
	TypeLink     int = 3
)

type Page struct {
	Id         string  `json:"id"`
	Set        bool    `json:"-"`
	Parent     *Page   `json:"-"`
	ParentLink *string `json:"parent,omitempty"`
	Level      int     `json:"level,omitempty"`
	Type       int     `json:"type,omitempty"`
	Title      *string `json:"title,omitempty"`
	Link       *string `json:"link,omitempty"`
	DataLink   *string `json:"-"`
	Pages      []Page  `json:"pages,omitempty"`
}

type Navi struct {
	Title string `json:"title,omitempty"`
	Id    string `json:"id"`
	Pages []Page `json:"pages,omitempty"`
}

func createLink(link *string) *string {
	if link != nil {

		if filepath.Base(*link) == "README.md" {
			l := "/" + filepath.Dir(*link)
			return &l
		}
		l := strings.TrimSuffix(*link, filepath.Ext(*link))
		l = "/" + tools.SetLastLash(l)

		return &l

	}
	return nil
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
				page.Id = ksuid.New().String()

				page.Link = createLink(&linkStr)

			}
		}
		var err error
		return s, err
	})

	if page.Title == nil {
		titleStr := string(node.FirstChild().Text([]byte(*source)))
		page.Title = &titleStr
	}

	return

}

func genNavi() (*Navi, error) {

	source, err := os.ReadFile("data/content/SUMMARY.md")
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := goldmark.Convert(source, &buf); err != nil {
		panic(err)
	}

	doc := goldmark.DefaultParser().Parse(text.NewReader([]byte(source)))
	listLevel := 0

	var navi Navi

	var entry Page
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkStatus(ast.WalkContinue)
		var err error

		if entering {

			if n.Kind() == ast.KindHeading {

				h := n.(*ast.Heading)
				if h.Level == 1 && navi.Title == "" {
					navi.Title = string(n.Text([]byte(source)))
					navi.Id = ksuid.New().String()
				} else if h.Level == 2 {

					if entry.Set {
						navi.Pages = append(navi.Pages, entry)
					}

					titleStr := string(n.Text([]byte(source)))

					entry = Page{
						Set:   true,
						Level: listLevel,
						Type:  TypeGroup,
						Title: &titleStr,
					}
				}

			} else if n.Kind() == ast.KindThematicBreak {
				if entry.Set {
					navi.Pages = append(navi.Pages, entry)
					entry = Page{}
				}
			} else if n.Kind() == ast.KindList {
				listLevel = listLevel + 1

			} else if n.Kind() == ast.KindListItem {

				if listLevel == 1 {

					pg := Page{}
					pg.Set = true
					pg.Id = ksuid.New().String()
					pg.Type = TypeMenuItem
					pg.Level = listLevel
					listitemlink(&pg, n.FirstChild(), &source)

					if entry.Type == 1 {
						list(n, 1, &pg, &source)
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
