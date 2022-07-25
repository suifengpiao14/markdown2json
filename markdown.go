package markdown2json

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func Markdown2XML(source []byte) (xhtml []byte, err error) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)
	var buf bytes.Buffer
	err = md.Convert(source, &buf)
	return buf.Bytes(), err
}

func Xhtml2json(xhtml []byte) (out []byte, err error) {
	doc, err := xmlquery.Parse(bytes.NewReader(xhtml))
	if err != nil {
		return nil, err
	}
	commentArr, err := xmlquery.QueryAll(doc, "//comment()")
	if err != nil {
		return nil, err
	}
	for _, comment := range commentArr {
		err = CommentParse(*comment)
		if err != nil {
			return nil, err
		}

	}
	return nil, nil
}

func CommentParse(commentNode xmlquery.Node) (err error) {
	if commentNode.Type != xmlquery.CommentNode {
		err = errors.Errorf("want node type %d ,got %d", xmlquery.CommentNode, commentNode.Type)
		return err
	}
	tag := ""
	firstCommaIndex := strings.Index(commentNode.Data, " ")
	if firstCommaIndex > -1 {
		tag = commentNode.Data[:firstCommaIndex]
	}
	fmt.Println(tag)

	return nil
}
