package app

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/blubooks/blubooks-cli/pkg/tools"
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
	DataLink   *string `json:"data"`
	Navi       *string `json:"navi"`
	Pages      []Page  `json:"pages,omitempty"`
}

/*

type MetaLinks struct {
	Name       string      `json:"name"`
	Link       string      `json:"link,omitempty"`
	DataLink   string      `json:"data"`
	Id         string      `json:"id,omitempty"`
	ExternLink bool        `json:"extern,omitempty"`
	Navi       string      `json:"navi,omitempty"`
	Links      []MetaLinks `json:"links,omitempty"`
}
*/

type Navi struct {
	Title    string `json:"title,omitempty"`
	Id       string `json:"id"`
	Path     *string
	Pages    []Page `json:"pages,omitempty"`
	Header   []Page `json:"header,omitempty"`
	Footer   []Page `json:"footer,omitempty"`
	Navis    []Navi
	FileName string
}

func createLink(link *string, nav *string) *string {
	if link != nil {

		if filepath.Base(*link) == "SUMMARY.md" {
			l := filepath.Dir(*link)
			l = strings.Replace(l, "/", "-", 0)
			l = "n/" + l
			if nav != nil {
				l = *nav + l
			}
			l = tools.SetFirstLash(l)
			return &l
		} else if *link == "README.md" {
			l := "/"
			return &l
		} else if filepath.Base(*link) == "README.md" {
			l := filepath.Dir(*link)
			if nav != nil {
				l = *nav + l
			}
			l = tools.SetFirstLash(l)
			return &l
		}
		l := strings.TrimSuffix(*link, filepath.Ext(*link))
		l = tools.SetLastLash(l)
		if nav != nil {
			l = *nav + l
		}
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
					pg.Id = ksuid.New().String()
					pg.Type = TypeLink
					pg.Level = level
					pg.Parent = page
					pg.ParentId = &page.Id
					pg.ParentLink = createLink(page.Link, nil)

					if filepath.Base(*page.DataLink) == "SUMMARY.md" {
						l := filepath.Dir(*page.DataLink)
						l = strings.Replace(l, "/", "-", 0)
						n := &Navi{
							Path:     &l,
							FileName: *page.DataLink,
						}
						err := genNavi(n, false)
						if err != nil {
							log.Println("Fehler - genNavi", err)
						} else {
							navi.Navis = append(navi.Navis, *n)
						}
					}

					listitemlink(&pg, n.FirstChild(), source, navi)
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

func listitemlink(page *Page, node ast.Node, source *[]byte, navi *Navi) {

	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkStatus(ast.WalkContinue)

		if entering {
			if n.Kind() == ast.KindLink {

				l := n.(*ast.Link)

				titleStr := string(n.Text([]byte(*source)))
				linkStr := string(l.Destination)
				page.Title = &titleStr
				page.DataLink = &linkStr
				page.Link = createLink(&linkStr, nil)

				if filepath.Base(*page.DataLink) == "SUMMARY.md" {
					l := filepath.Dir(*page.DataLink)
					l = strings.Replace(l, "/", "-", 0)
					n := &Navi{
						Path:     &l,
						FileName: *page.DataLink,
					}
					err := genNavi(n, false)
					if err != nil {
						log.Println("Fehler - genNavi", err)
					} else {
						navi.Navis = append(navi.Navis, *n)
					}
				}

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

/*
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
*/
func genNavi(navi *Navi, isRoot bool) error {
	log.Println("testbools", isRoot)
	/*


		if (navi == nil) {
			dir := filepath.Dir(filename)
			base := filepath.Base(filename)

			if base != "SUMMARY.md" {
				return nil, errors.New("No SUMMARY.md")
			}

			var navi Navi

		}
		if dir == "." {
			isRoot = true
		} else {
			s := strings.Replace(dir, "/", "-", 0)
			navi.Path = &s
		}

		if !strings.HasSuffix(filename, "SUMMARY.md") {
			return nil, errors.New("no SUMMERY.md found")
		}

	*/
	//var navi Navi

	if isRoot {
		navi.FileName = "SUMMARY.md"
	}

	source, err := os.ReadFile("data/content/" + navi.FileName)
	if err != nil {
		return err
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

	navi.Id = getUrlId("README.md")

	doc := markdown.Parser().Parse(text.NewReader([]byte(source)))
	//doc := goldmark.DefaultParser().Parse(text.NewReader([]byte(source)))
	listLevel := 0

	var entry Page
	naviIsSet := false
	loopType := 0 // 1: header, 2: navi, 3: footer
	headerIsSet := false
	footerIsSet := false
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkStatus(ast.WalkContinue)
		var err error

		if entering {

			if n.Kind() == ast.KindHeading {

				h := n.(*ast.Heading)
				if h.Level == 1 {
					t := string(n.Text([]byte(source)))
					if t == "NAV" && !naviIsSet {
						naviIsSet = true
						loopType = 2
					}
					if isRoot && t == "HEADER" && !headerIsSet {
						headerIsSet = true
						loopType = 1
					}
					if isRoot && t == "FOOTER" && !footerIsSet {
						footerIsSet = true
						loopType = 3
					}
					/*
						if !titleSet && h.Level == 1 && navi.Title == "" {
							titleSet = true
							navi.Title = string(n.Text([]byte(source)))
					*/

				} else if h.Level == 2 {
					if loopType == 2 {
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

				if entry.Set {
					if loopType == 1 {
						navi.Header = append(navi.Header, entry)
					}
					if loopType == 2 {
						navi.Pages = append(navi.Pages, entry)
					}
					if loopType == 3 {
						navi.Footer = append(navi.Footer, entry)
					}
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
					listitemlink(&pg, n.FirstChild(), &source, navi)

					list(n, 1, &pg, &source, navi)
					if entry.Type == 1 {
						entry.Pages = append(entry.Pages, pg)
					} else {
						if loopType == 1 {
							navi.Header = append(navi.Header, pg)
						}
						if loopType == 2 {
							navi.Pages = append(navi.Pages, pg)
						}
						if loopType == 3 {
							navi.Footer = append(navi.Footer, pg)
						}
					}

				}
			}

		} else {
			if n.Kind() == ast.KindList {
				listLevel = listLevel - 1

			} else if n.Kind() == ast.KindDocument {
				if entry.Set {
					if loopType == 1 {
						navi.Header = append(navi.Header, entry)
					}
					if loopType == 2 {
						navi.Pages = append(navi.Pages, entry)

					}
					if loopType == 3 {
						navi.Footer = append(navi.Footer, entry)
					}

				}
			}
		}

		return s, err
	})

	return nil
	//return &navi, nil
}
