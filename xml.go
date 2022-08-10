package markdown2json

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// 为啥不用xml，1.转换为xml 后标记增多，不方便获取相邻元素文本。2. 引号" 等特殊字符被转义，需要反转义恢复原样。相比处理起来比较麻烦，所以保留，后续可以再权衡考虑
func Markdown2XML(source []byte) (out string) {
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
	if err := md.Convert(source, &buf); err != nil {
		panic(err)
	}
	return buf.String()
}

func QueryXML(xhtml string) (err error) {
	doc, err := xmlquery.Parse(strings.NewReader(xhtml))
	if err != nil {
		return err
	}
	list := xmlquery.Find(doc, ".//comment()")
	for _, comment := range list {

		fmt.Println(comment.Data)
	}

	return nil
}
