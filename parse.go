package markdown2json

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/suifengpiao_14/markdown2json/convmap"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
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
			parser.WithAttribute(),
			//parser.WithASTTransformers(util.PrioritizedValue{}),
			//parser.WithInlineParsers(util.PrioritizedValue{}),
			//parser.WithBlockParsers(util.PrioritizedValue{}),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)
	reader := text.NewReader(source)
	document := md.Parser().Parse(reader)
	metaData := document.OwnerDocument().Meta()
	//document.Dump(source, 10)
	parseApi(document, source)

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

func parseApi(node ast.Node, source []byte) string {

	if node.Type() == ast.TypeBlock {
		htmlBlock, ok := node.(*ast.HTMLBlock)
		if ok && htmlBlock.HTMLBlockType == ast.HTMLBlockType2 {
			for i := 0; i < htmlBlock.Lines().Len(); i++ {
				line := htmlBlock.Lines().At(i)
				fmt.Printf("%s", line.Value(source))
			}
			fmt.Printf("\"\n")
		}
	}
	if node.HasChildren() {
		firstChild := node.FirstChild()
		parseApi(firstChild, source)
	}

	nextNode := node.NextSibling()
	if nextNode != nil {
		parseApi(nextNode, source)
	}

	return ""
}
