package markdown2json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

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
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)
	reader := text.NewReader(source)
	document := md.Parser().Parse(reader)
	//metaData := document.OwnerDocument().Meta()
	//document.Dump(source, 10)
	apiElements := parseApi(document, source)
	api, err := FormatApi(apiElements)
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(api)
	if err != nil {
		return "", err
	}
	out = string(b)
	return out, nil

	// newMetaData, err := convmap.Convert(metaData, convmap.ConvertMapKeyStrict)
	// if err != nil {
	// 	return "", err
	// }
	// byt, err := json.Marshal(newMetaData)
	// if err != nil {
	// 	return "", err
	// }
	// out = string(byt)
	// return out, nil
}

type Attr struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type ApiAttribute struct {
	Name  string  `json:"name"`
	Value string  `json:"value"`
	Attrs []*Attr `json:"attrs"`
}

func parseApi(node ast.Node, source []byte) (apiAttributes []*ApiAttribute) {
	apiAttributes = make([]*ApiAttribute, 0)
	if htmlBlock, ok := node.(*ast.HTMLBlock); ok && htmlBlock.HTMLBlockType == ast.HTMLBlockType2 {
		htmlRaw := Node2RawText(htmlBlock, source)
		attribute := parseRawHtml(htmlRaw)
		if attribute.Value == "" {
			nextNode := node.NextSibling()
			if fencedCodeNode, ok := nextNode.(*ast.FencedCodeBlock); ok {
				attr := fencedCodeContentType(*fencedCodeNode, source)
				attribute.Attrs = append(attribute.Attrs, attr)
			}
			attribute.Value = Node2RawText(nextNode, source)
		}
		apiAttributes = append(apiAttributes, attribute)
	}
	if rawHTML, ok := node.(*ast.RawHTML); ok {
		txt := Node2RawText(rawHTML, source)
		attribute := parseRawHtml(string(txt))
		if attribute.Value == "" {
			nextNode := node.NextSibling()
			attribute.Value = Node2RawText(nextNode, source)
		}
		apiAttributes = append(apiAttributes, attribute)
	}
	if node.HasChildren() {
		firstChild := node.FirstChild()
		subAttr := parseApi(firstChild, source)
		apiAttributes = append(apiAttributes, subAttr...)
	}
	nextNode := node.NextSibling()
	if nextNode != nil {
		subAttr := parseApi(nextNode, source)
		apiAttributes = append(apiAttributes, subAttr...)
	}
	return apiAttributes
}

func fencedCodeContentType(fencedCodeNode ast.FencedCodeBlock, source []byte) (attr *Attr) {
	info := fencedCodeNode.Info
	val := string(info.Text(source))
	value := ""
	switch val {
	case "multipart":
		value = "multipart/form-data"
	case "urlencoded":
		value = "application/x-www-form-urlencoded"
	case "json":
		value = "application/json"
	default:
		value = ""
	}
	attr = &Attr{
		Name:  "contentType",
		Value: value,
	}
	return attr
}

func parseRawHtml(rawText string) (attribute *ApiAttribute) {

	attribute = &ApiAttribute{
		Attrs: make([]*Attr, 0),
	}
	rawText = strings.TrimSpace(rawText)
	name := strings.Trim(rawText, "<!->")
	value := ""
	attributeStr := ""
	index := strings.Index(name, " ")
	if index > -1 {
		attributeStr = fmt.Sprintf("{%s}", name[index+1:])
		name = name[:index]
	}
	if strings.Contains(name, "=") {
		arr := strings.SplitN(name, "=", 2)
		name = arr[0]
		value = arr[1]
	}
	attribute.Name = strings.TrimSpace(name)

	if attributeStr != "" {
		txtReader := text.NewReader([]byte(attributeStr))
		attrs, ok := parser.ParseAttributes(txtReader)
		if ok {
			for _, parseAttr := range attrs {
				attr := Attr{
					Name:  string(parseAttr.Name),
					Value: fmt.Sprintf("%s", parseAttr.Value),
				}
				attribute.Attrs = append(attribute.Attrs, &attr)
			}
		}
	}
	attribute.Name = name
	attribute.Value = value
	return attribute
}

func Node2RawText(node ast.Node, source []byte) (out string) {
	if node == nil {
		return ""
	}
	var w bytes.Buffer
	if node.Type() == ast.TypeBlock {
		for i := 0; i < node.Lines().Len(); i++ {
			line := node.Lines().At(i)
			w.Write(line.Value(source))
		}
		out = strings.TrimSpace(w.String())
		return out
	}

	if rawHTML, ok := node.(*ast.RawHTML); ok {
		t := []string{}
		for i := 0; i < rawHTML.Segments.Len(); i++ {
			segment := rawHTML.Segments.At(i)
			t = append(t, string(segment.Value(source)))
		}
		out = strings.Join(t, "")
		return strings.TrimSpace(out)
	}
	if autoLink, ok := node.(*ast.AutoLink); ok {
		b := autoLink.URL(source)
		out = string(b)
		return strings.TrimSpace(out)
	}
	b := node.Text(source)
	out = string(b)
	return strings.TrimSpace(out)
}

type KV struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type Api struct {
	URL        string `json:"url"`
	Method     string `json:"metod"`
	Header     []*KV  `json:"header"`
	Query      string `json:"query"`
	Body       string `json:"body"`
	PreRequest string `json:"preRequest"`
}

const (
	API_URI         = "api.uri"
	API_HOST        = "api.host"
	API_METHOD      = "api.method"
	API_HEADER      = "api.header"
	API_QUERY       = "api.query"
	API_BODY        = "api.body"
	API_PRE_REQUEST = "api.preRequest"
)

func FormatApi(apiAttributes []*ApiAttribute) (api *Api, err error) {
	api = &Api{
		Header: make([]*KV, 0),
	}
	u := &url.URL{}
	uv := url.Values{}
	uri := ""
	for _, apiAttr := range apiAttributes {
		switch apiAttr.Name {
		case API_URI:
			uri = apiAttr.Value
		case API_HOST:
			if apiAttr.Value[:4] == "http" {
				u, err = url.Parse(apiAttr.Value)
				if err != nil {
					return nil, err
				}
			}
		case API_QUERY:
			uv, err = url.ParseQuery(apiAttr.Value)
			if err != nil {
				return nil, err
			}
		case API_BODY:
			api.Body = apiAttr.Value
		case API_METHOD:
			api.Method = apiAttr.Value
		case API_HEADER:
			index := strings.Index(apiAttr.Value, ":")
			kv := &KV{}
			if index > -1 {
				kv.Name = apiAttr.Value[:index]
				kv.Value = apiAttr.Value[index+1:]
			}
			api.Header = append(api.Header, kv)
		case API_PRE_REQUEST:
			fmt.Println(apiAttr.Value)

		}
	}
	if uri != "" {
		u.Path = uri
	}
	if len(uv) > 0 {
		u.RawQuery = uv.Encode()
	}
	api.URL = u.String()
	return api, nil
}
