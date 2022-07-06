package markdown2json

import (
	"encoding/json"
	"io"
	"os"

	"github.com/suifengpiao_14/markdown2json/convmap"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

func Parse(mdxFile string) (out string, err error) {
	fd, err := os.OpenFile(mdxFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return "", err
	}
	source, err := io.ReadAll(fd)
	if err != nil {
		return "", err
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			meta.New(
				meta.WithStoresInDocument(),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)
	reader := text.NewReader(source)
	document := md.Parser().Parse(reader)
	metaData := document.OwnerDocument().Meta()
	newMetaData, err := convmap.Convert(metaData, convmap.ConvertMapKeyStrict)
	if err != nil {
		return "", err
	}
	byt, err := json.Marshal(newMetaData)
	if err != nil {
		return "", err
	}
	out = string(byt)
	return out, nil
}
