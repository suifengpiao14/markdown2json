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
}

type Attr struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Attrs []*Attr

func (attrs Attrs) GetByName(name string) (attr *Attr, ok bool) {
	for _, at := range attrs {
		if at.Name == name {
			return at, true
		}
	}
	return nil, false
}

type ApiElement struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Attrs Attrs  `json:"attrs"`
}

type ApiElements []ApiElement

func parseApi(node ast.Node, source []byte) (apiElements []*ApiElement) {
	apiElements = make([]*ApiElement, 0)
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
		apiElements = append(apiElements, attribute)
	}
	if rawHTML, ok := node.(*ast.RawHTML); ok {
		txt := Node2RawText(rawHTML, source)
		attribute := parseRawHtml(string(txt))
		if attribute.Value == "" {
			nextNode := node.NextSibling()
			attribute.Value = Node2RawText(nextNode, source)
		}
		apiElements = append(apiElements, attribute)
	}
	if node.HasChildren() {
		firstChild := node.FirstChild()
		subAttr := parseApi(firstChild, source)
		apiElements = append(apiElements, subAttr...)
	}
	nextNode := node.NextSibling()
	if nextNode != nil {
		subAttr := parseApi(nextNode, source)
		apiElements = append(apiElements, subAttr...)
	}
	return apiElements
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

func parseRawHtml(rawText string) (attribute *ApiElement) {

	attribute = &ApiElement{
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

type Example struct {
	Name     string `json:"name"`
	Request  string `json:"request"`
	Response string `json:"response"`
}

func (example Example) FillAttr(attr Attr) {
	switch attr.Name {
	case "name":
		example.Name = attr.Value
	case "direction":

	}
}

type Examples []*Example

func (examples Examples) GetByName(name string) (example *Example, ok bool) {
	for _, ex := range examples {
		if ex.Name == name {
			return ex, true
		}
	}
	return nil, false
}

//Add 增加案例，存在替换，不存在新增
func (examples Examples) Add(example *Example) {
	for i, ex := range examples {
		if ex.Name == example.Name {
			examples[i] = example
			return
		}
	}
	examples = append(examples, example)
	return
}

type Api struct {
	URL        string   `json:"url"`
	Method     string   `json:"metod"`
	Header     []*KV    `json:"header"`
	Query      string   `json:"query"`
	Body       string   `json:"body"`
	PreRequest string   `json:"preRequest"`
	Examples   Examples `json:"examples"`
}

const (
	API_URI         = "api.uri"
	API_HOST        = "api.host"
	API_METHOD      = "api.method"
	API_HEADER      = "api.header"
	API_QUERY       = "api.query"
	API_BODY        = "api.body"
	API_PRE_REQUEST = "api.preRequest"
	API_EXAMPLE     = "api.example"
)

func FormatApi(apiAttributes []*ApiElement) (api *Api, err error) {
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
			api.PreRequest = apiAttr.Value
		case API_EXAMPLE:
			nameAttr, ok := apiAttr.Attrs.GetByName("name")
			if !ok {
				nameAttr = &Attr{}
			}
			name := nameAttr.Value
			example, ok := api.Examples.GetByName(name)
			if !ok {
				example = &Example{
					Name: name,
				}
			}

			api.Examples.Add(example)

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
