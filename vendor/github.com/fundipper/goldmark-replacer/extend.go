// Package replacer is a extension for the goldmark
// (http://github.com/yuin/goldmark).
//
// This extension adds support for authomaticaly replacing text in markdowns.
package replacer

import (
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

var extender *Extender

type Extender struct {
	*strings.Replacer
}

// New return initialized image render with source url replacing support.
func NewExtender(oldnew ...string) goldmark.Extender {
	extender = &Extender{
		Replacer: strings.NewReplacer(oldnew...),
	}
	return extender
}

func (e *Extender) Extend(m goldmark.Markdown) {
	// m.Parser().AddOptions(
	// 	parser.WithASTTransformers(
	// util.Prioritized(NewTransformer(), 500),
	// ),
	// )
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(NewRenderer(), 500),
		),
	)
}
