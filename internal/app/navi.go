package app

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/blubooks/blubooks-cli/pkg/tools"
	"github.com/segmentio/ksuid"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/frontmatter"
)

const (
	TypeBook       string = "book"
	TypeGroup      string = "group"
	TypeLink       string = "link"
	TypeExternLink string = "link-extern"
	TypeOtherLink  string = "link-other"
)

var naviUrlIds map[string]string

type NaviMetaData struct {
	Header *struct {
		Pages      []Page `json:"pages,omitempty"`
		Hide       bool   `json:"hide,omitempty"`
		HideNavi   bool   `json:"showNavi,omitempty"`
		BeforeNavi bool   `json:"beforeNavi,omitempty"`
	} `yaml:"header"`
	Footer *struct {
		Hide  bool   `json:"hide,omitempty"`
		Pages []Page `json:"pages,omitempty"`
	} `yaml:"footer"`
	Options *Options `yaml:"options" json:"options,omitempty"`
}

type Options struct {
	Accordion bool `json:"accordion,omitempty"`
}

type Page struct {
	Set        bool    `json:"-"`
	Parent     *Page   `json:"-"`
	Id         string  `json:"id,omitempty"`
	ParentLink *string `json:"parent,omitempty"`
	ParentId   *string `json:"parentId,omitempty"`
	Level      int     `json:"level,omitempty"`
	Type       string  `json:"type,omitempty"`
	Title      *string `json:"title,omitempty"`
	Link       *string `json:"link,omitempty"`
	ExternLink bool    `json:"extern,omitempty"`
	DataLink   *string `json:"data,omitempty"`
	Pages      []Page  `json:"pages,omitempty"`
}

/*
type Navi struct {
	Title       *string  `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Header      []Page   `json:"header,omitempty"`
	Footer      []Page   `json:"footer,omitempty"`
	Pages       []Page   `json:"pages,omitempty"`
	Root        *Page    `json:"root,omitempty"`
	Options     *Options `json:"options,omitempty"`
	SearchId    *string  `json:"searchId,omitempty"`
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
				page.Type = TypeExternLink

			} else {
				if strings.HasSuffix(*page.DataLink, ".md") {
					page.Id = getUrlId(*page.Link)
					page.Type = TypeLink
				} else {
					page.Type = TypeOtherLink
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
	var metaData NaviMetaData
	//var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var subEntry Page
	var book Page
	var buf bytes.Buffer

	context := parser.NewContext()
	listLevel := 0
	inNavi := false

	naviUrlIds = make(map[string]string)

	source, err := os.ReadFile("data/content/SUMMARY.md")
	if err != nil {
		return nil, err
	}

	markdown := goldmark.New(
		goldmark.WithExtensions(
			&frontmatter.Extender{
				//Mode: frontmatter.SetMetadata,
			},
		),
	)

	if err := markdown.Convert(source, &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}

	doc := markdown.Parser().Parse(text.NewReader([]byte(source)))
	//doc := goldmark.DefaultParser().Parse(text.NewReader([]byte(source)))

	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkStatus(ast.WalkContinue)
		var err error

		if entering {

			if n.Kind() == ast.KindHeading {

				h := n.(*ast.Heading)
				t := string(n.Text([]byte(source)))
				if h.Level == 1 {
					if !book.Set && !inNavi {
						navi.Title = &t
						inNavi = true
					}

				}
				if h.Level == 2 {
					inNavi = false
					if book.Set {
						navi.Pages = append(navi.Pages, book)
					}
					book = Page{
						Id:    ksuid.New().String(),
						Set:   true,
						Level: listLevel,
						Type:  TypeBook,
						Title: &t,
					}
					p := book

					navi.Header = append(navi.Header, p)

				} else if h.Level == 3 {
					if book.Set {
						if subEntry.Set {
							book.Pages = append(book.Pages, subEntry)
						}

						subEntry = Page{
							Id:    ksuid.New().String(),
							Set:   true,
							Level: listLevel,
							Type:  TypeGroup,
							Title: &t,
						}

					}

				}

			} else if n.Kind() == ast.KindThematicBreak {
				if book.Set {
					if subEntry.Set {
						book.Pages = append(book.Pages, subEntry)
						subEntry = Page{}
					}
				}

			} else if n.Kind() == ast.KindList {
				listLevel = listLevel + 1

			} else if n.Kind() == ast.KindListItem {

				if listLevel == 1 {

					if book.Set {

						pg := Page{}
						pg.Set = true
						pg.Level = listLevel
						if subEntry.Id != "" {
							id := subEntry.Id
							pg.ParentId = &id
						}
						listitemlink(&pg, n.FirstChild(), &source, &navi)
						list(n, 1, &pg, &source, &navi)

						if subEntry.Type == TypeGroup {
							subEntry.Pages = append(subEntry.Pages, pg)
						} else {
							book.Pages = append(book.Pages, pg)
						}
					}

				}
			} else if n.Kind() == ast.KindParagraph {
				if inNavi && navi.Title != nil {
					if navi.Description == "" {
						navi.Description = string(n.Text(source))
					} else {
						navi.Description = navi.Description + "\n" + string(n.Text(source))
					}
				}
			}

		} else {
			if n.Kind() == ast.KindList {
				listLevel = listLevel - 1

			} else if n.Kind() == ast.KindDocument {
				if book.Set {
					navi.Pages = append(navi.Pages, book)
				}

			}
		}

		return s, err
	})

	if navi.Root == nil {
		pg := Page{}
		pg.Type = TypeLink
		pg.Level = 0
		l := "README.md"
		pg.Link = &l
		createLink(&pg, &navi)
	}

	d := frontmatter.Get(context)
	if err := d.Decode(&metaData); err != nil {
		log.Println(err)
	}

	if metaData.Header != nil {
		if metaData.Header.Hide {
			navi.Header = nil
		} else {
			if metaData.Header.HideNavi {
				navi.Header = nil
			}

			if len(metaData.Header.Pages) > 0 {
				metalinks(metaData.Header.Pages, &navi, nil)
				if metaData.Header.BeforeNavi {
					navi.Header = append(metaData.Header.Pages, navi.Header...)
				} else {
					navi.Header = append(navi.Header, metaData.Header.Pages...)

				}
			}

		}
	}

	if metaData.Footer != nil {

		if metaData.Footer.Hide {
			navi.Footer = nil

		} else {

			if len(metaData.Footer.Pages) > 0 {
				metalinks(metaData.Footer.Pages, &navi, nil)

				navi.Footer = append(navi.Footer, metaData.Footer.Pages...)

			}

		}
	}

	if metaData.Options != nil {
		navi.Options = metaData.Options
	}

	return &navi, nil

}
