package main

import (
	"github.com/carlos7ags/folio/document"
	"github.com/carlos7ags/folio/font"
	"github.com/carlos7ags/folio/layout"
)

func main() {
	doc := document.NewDocument(document.PageSizeA4)
	doc.Info.Title = "Hello World"
	doc.SetAutoBookmarks(true)

	doc.Add(layout.NewHeading("Hello, Folio!", layout.H1))
	doc.Add(layout.NewParagraph(
		"Generated with Folio — the modern PDF library for Go.",
		font.Helvetica, 12,
	))

	doc.Save("hello.pdf")
}
