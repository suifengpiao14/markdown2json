package markdown2json

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

func Parse(source []byte) (apiElements ApiElements, err error) {

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
	apiElements = parseApi(document, source, nil)
	return apiElements, nil
}

type Attr struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Attrs []*Attr

func (attrs *Attrs) GetByName(name string) (attr *Attr, ok bool) {
	for _, at := range *attrs {
		if at.Name == name {
			return at, true
		}
	}
	return nil, false
}
func (attrs *Attrs) PopByName(name string) (attr *Attr, ok bool) {
	for i, at := range *attrs {
		if at.Name == name {
			*attrs = append((*attrs)[:i], (*attrs)[i+1:]...)
			return at, true
		}
	}
	return nil, false
}

type ApiElement struct {
	Tag      string `json:"tag"`
	Name     string `json:"name"`
	RefLevel string `json:"refLevel"`
	Value    string `json:"value"`
	Attrs    Attrs  `json:"attrs"`
}

func (apiElement *ApiElement) Clone() (clone ApiElement) {
	clone = ApiElement{
		Tag:   apiElement.Tag,
		Name:  apiElement.Name,
		Value: apiElement.Value,
	}
	for _, attr := range apiElement.Attrs {
		tmpAttr := *attr
		clone.Attrs = append(clone.Attrs, &tmpAttr)
	}
	return
}

func (apiElement *ApiElement) AddAttr(attr *Attr) {
	if attr.Name == "name" && apiElement.Name == "" {
		apiElement.Name = attr.Name
		apiElement.Value = attr.Value
		return
	}
	apiElement.Attrs = append(apiElement.Attrs, attr)
}

type ApiElements []ApiElement

func (apiElements *ApiElements) GetByName(name string) (subApiElements []ApiElement, ok bool) {
	subApiElements = ApiElements{}
	for _, elem := range *apiElements {
		if elem.Name == name {
			subApiElements = append(subApiElements, elem)
		}
	}
	return subApiElements, len(subApiElements) > 0
}

//弹出找到的第一个
func (apiElements *ApiElements) PopByName(name string) (apiElement *ApiElement, ok bool) {
	for i, elem := range *apiElements {
		if elem.Name == name {
			*apiElements = append((*apiElements)[:i], (*apiElements)[i+1:]...)
			return &elem, true
		}
	}
	return nil, false
}

func parseApi(node ast.Node, source []byte, parent *ApiElement) (apiElements []ApiElement) {
	apiElements = make([]ApiElement, 0)
	var attribute *ApiElement
	if htmlBlock, ok := node.(*ast.HTMLBlock); ok && htmlBlock.HTMLBlockType == ast.HTMLBlockType2 {
		htmlRaw := Node2RawText(htmlBlock, source)
		attribute, ok = parseRawHtml(htmlRaw)
		if ok {
			if attribute.Tag == "" {
				attribute.Tag = "unexcept block tag"
			}
			if attribute.Value == "" {
				nextNode := node.NextSibling()
				if fencedCodeNode, ok := nextNode.(*ast.FencedCodeBlock); ok {
					attr := &Attr{
						Name:  "language",
						Value: string(fencedCodeNode.Language(source)),
					}
					attribute.Attrs = append(attribute.Attrs, attr)
					attribute.Value = Node2RawText(nextNode, source)
				} else if tableHTML, ok := nextNode.(*extast.Table); ok {
					columnAttr, ok := attribute.Attrs.PopByName("column")
					if !ok {
						err := errors.Errorf("table element required coulmn attribute")
						panic(err)
					}
					columnArr := strings.Split(columnAttr.Value, ",")
					keyMap := map[string]string{}
					keyMapAttr, ok := attribute.Attrs.PopByName("keymap")
					if ok {
						formatArr := strings.Split(keyMapAttr.Value, ",")
						for _, kvStr := range formatArr {
							kvStr = strings.TrimSpace(kvStr)
							key := ""
							value := ""
							if strings.Contains(kvStr, ":") {
								arr := strings.SplitN(kvStr, ":", 2)
								key = arr[0]
								value = arr[1]
							} else {
								key = kvStr
								value = kvStr
							}
							keyMap[key] = value
						}
					}

					firstNode := tableHTML.FirstChild()
					headNode, ok := firstNode.(*extast.TableHeader)
					if !ok {
						err := errors.Errorf("first children is not header")
						panic(err)
					}
					columnLen := len(columnArr)
					if columnLen != headNode.ChildCount() {
						err := errors.Errorf("coulmn filed not match table head field")
						panic(err)
					}

					var subNode ast.Node
					subNode = headNode.NextSibling()
					for {
						if subNode == nil {
							break
						}
						tableRow, ok := subNode.(*extast.TableRow)
						if !ok {
							err := errors.Errorf("subNode must be tableRow")
							panic(err)
						}
						cellIndex := 0
						cellNode := tableRow.FirstChild()
						rowApiElement := attribute.Clone()
						var lastAttr = &Attr{}
						for {
							if cellNode == nil {
								break
							}
							name := columnArr[cellIndex]
							value := cellNode.Text(source)
							lastAttr = &Attr{
								Name:  name,
								Value: string(value),
							}
							rowApiElement.AddAttr(lastAttr)
							cellNode = cellNode.NextSibling()
							cellIndex++
						}
						if keyMap != nil { // 获取最后一列更多内容
							lineNode := tableRow.LastChild().FirstChild()
							txtArr := make([]string, 0)
							txt := ""
							for {
								if lineNode == nil {
									if txt != "" {
										txt = strings.TrimSpace(txt)
										txtArr = append(txtArr, txt)
									}
									break
								}
								textNode, ok := lineNode.(*ast.Text)
								if ok {
									txt = fmt.Sprintf("%s%s", txt, string(textNode.Text(source)))
									lineNode = lineNode.NextSibling()
									continue
								}
								if txt != "" {
									txt = strings.TrimSpace(txt)
									txtArr = append(txtArr, txt)
								}
								txt = ""
								lineNode = lineNode.NextSibling()
							}

							for _, txt := range txtArr {
								if strings.Contains(txt, ":") {
									arr := strings.SplitN(txt, ":", 2)
									key := arr[0]
									value := arr[1]
									key, ok := keyMap[key]
									if !ok {
										continue
									}
									rowApiElement.AddAttr(&Attr{
										Name:  key,
										Value: value,
									})
								} else if attribute.Value == "" {
									lastAttr.Value = txt
								}
							}
						}

						apiElements = append(apiElements, rowApiElement)
						subNode = subNode.NextSibling()
					}
				}

			}
			apiElements = append(apiElements, *attribute)
		}

	} else if rawHTML, ok := node.(*ast.RawHTML); ok {
		txt := Node2RawText(rawHTML, source)
		attribute, ok := parseRawHtml(string(txt))
		if ok {
			if attribute.Value == "" {
				nextNode := node.NextSibling()
				attribute.Value = Node2RawText(nextNode, source)
			}
			apiElements = append(apiElements, *attribute)
		}

	}
	if node.HasChildren() {
		firstChild := node.FirstChild()
		subParent := parent
		if len(apiElements) > 0 {
			subParent = &(apiElements[len(apiElements)-1])
		}
		subAttr := parseApi(firstChild, source, subParent)
		apiElements = append(apiElements, subAttr...)
	}
	nextNode := node.NextSibling()
	if nextNode != nil {
		subAttr := parseApi(nextNode, source, parent)
		apiElements = append(apiElements, subAttr...)
	}
	return apiElements
}

func parseRawHtml(rawText string) (attribute *ApiElement, ok bool) {

	attribute = &ApiElement{
		Attrs: make([]*Attr, 0),
	}
	rawText = strings.TrimSpace(rawText)
	tag := strings.Trim(rawText, "<!->")
	value := ""
	attributeStr := ""
	index := strings.Index(tag, " ")
	if index > -1 {
		attributeStr = fmt.Sprintf("{%s}", tag[index+1:])
		tag = tag[:index]
	}
	if strings.Contains(tag, "=") {
		arr := strings.SplitN(tag, "=", 2)
		tag = arr[0]
		value = arr[1]
	}
	if !strings.HasPrefix(tag, "doc.") {
		return nil, false
	}
	tag = strings.TrimSpace(strings.TrimPrefix(tag, "doc."))
	if strings.Contains(tag, ".") {
		arr := strings.SplitN(tag, ".", 2)
		tag = arr[0]
		attribute.Name = strings.TrimSpace(arr[1])
	}

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
	attribute.Tag = tag
	attribute.Value = value
	return attribute, true
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
func (examples Examples) Add(example *Example) (newExamples Examples) {
	for i, ex := range examples {
		if ex.Name == example.Name {
			examples[i] = example
			return examples
		}
	}
	examples = append(examples, example)
	return examples
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

func FormatApi(apiElement []*ApiElement) (api *Api, err error) {
	api = &Api{
		Header: make([]*KV, 0),
	}
	u := &url.URL{}
	uv := url.Values{}
	uri := ""
	for _, apiAttr := range apiElement {
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
