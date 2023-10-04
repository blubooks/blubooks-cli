# goldmark-replacer

[Goldmark](https://github.com/yuin/goldmark) text replacer extension.

## code

```go
func Example() {
	md := goldmark.New(
		goldmark.WithExtensions(
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
		),
	)
	var source = []byte("(c)Dmitry Sedykh")
	err := md.Convert(source, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
```

## view

```html
<p>©Dmitry Sedykh</p>
```

## thanks

[Goldmark](https://github.com/yuin/goldmark)

[goldmark-text-replacer](https://github.com/mdigger/goldmark-text-replacer)
