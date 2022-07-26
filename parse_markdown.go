package markdown2json

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func Parse(source []byte) (xmlTags XMLTags, err error) {

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
	xmlTags = parseNode(document, source)
	return xmlTags, nil
}

type Attr struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Attrs []*Attr

func (attrs *Attrs) GroupByName() (group map[string]Attrs) {
	group = make(map[string]Attrs)
	for _, attr := range *attrs {
		_, ok := group[attr.Name]
		if !ok {
			group[attr.Name] = Attrs{}
		}
		group[attr.Name] = append(group[attr.Name], attr)
	}
	return group
}

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

func (attrs *Attrs) Uniqueue() (uniqueueAttrs Attrs) {
	uniqueueAttrs = Attrs{}
	keyMap := make(map[string]bool)
	for _, at := range *attrs {
		valueKey := at.Value
		if len(at.Value) > 20 {
			valueKey = Md5(at.Value)
		}
		key := fmt.Sprintf("%s_%s", at.Name, valueKey)
		if _, ok := keyMap[key]; !ok {
			keyMap[key] = true
			uniqueueAttrs = append(uniqueueAttrs, at)

		}
	}
	return uniqueueAttrs
}

func (attrs *Attrs) ToMap() (out map[string]string) {
	out = make(map[string]string)
	for name, sameNameAttr := range attrs.GroupByName() {
		uniqueueAttrs := sameNameAttr.Uniqueue()
		valueArr := make([]string, 0)
		for _, attr := range uniqueueAttrs {
			valueArr = append(valueArr, attr.Value)
		}
		value := strings.Join(valueArr, ",")
		out[name] = value
	}
	return out
}

type XMLTag struct {
	Tag      string `json:"tag"`
	ID       string `json:"id"`
	RefLevel string `json:"refLevel"`
	Attrs    Attrs  `json:"attrs"`
}

func (xmlTag *XMLTag) Clone() (clone XMLTag) {
	clone = XMLTag{
		ID:    xmlTag.ID,
		Tag:   xmlTag.Tag,
		Attrs: Attrs{},
	}
	for _, attr := range xmlTag.Attrs {
		tmpAttr := *attr
		clone.Attrs = append(clone.Attrs, &tmpAttr)
	}
	return
}

func (xmlTag *XMLTag) AddAttr(attr *Attr) {
	if attr.Name == "id" {
		xmlTag.ID = attr.Value
		return
	}
	xmlTag.Attrs = append(xmlTag.Attrs, attr)
}

type XMLTags []XMLTag

func (xmlTags *XMLTags) GetByTag(tag string) (subXMLTags XMLTags) {
	subXMLTags = XMLTags{}
	for _, elem := range *xmlTags {
		if elem.Tag == tag {
			subXMLTags = append(subXMLTags, elem)
		}
	}
	return subXMLTags
}

func (xmlTags *XMLTags) GroupByTag() (group map[string]XMLTags) {
	group = make(map[string]XMLTags)
	for _, elem := range *xmlTags {
		_, ok := group[elem.Tag]
		if !ok {
			group[elem.Tag] = make(XMLTags, 0)
		}
		group[elem.Tag] = append(group[elem.Tag], elem)
	}
	return group
}

func (xmlTags *XMLTags) GetByID(id string) (subXMLTags XMLTags) {
	subXMLTags = XMLTags{}
	for _, elem := range *xmlTags {
		if elem.ID == id {
			subXMLTags = append(subXMLTags, elem)
		}
	}
	return subXMLTags
}

//Format 将ID为空的xmlTag 属性赋值到ID不为空的xmlTag 中，并删除，如果不存在ID不为空的xmlTag 则直接返回ID为空的
func (xmlTags *XMLTags) CopyEmptyID2NotEmpty() XMLTags {
	emptyIDXmls := xmlTags.GetByEmptyID()
	notEmptyIDXmls := xmlTags.GetByNotEmptyID()
	if len(notEmptyIDXmls) < 1 {
		return emptyIDXmls
	}
	for _, emptyIDXml := range emptyIDXmls {
		for i, notEmptyIDXml := range notEmptyIDXmls {
			if emptyIDXml.Tag == notEmptyIDXml.Tag { // 相同tag才复制
				notEmptyIDXmls[i].Attrs = append(notEmptyIDXml.Attrs, emptyIDXml.Attrs...) //此处非引用，所以指定index替换
			}
		}
	}
	return notEmptyIDXmls
}

//GetByEmptyID 获取ID为空的元素,ID为空的元素，会将其attrs 全部复制到id不为空的元素中，因此最后先按tag、id聚合后执行该操作）
func (xmlTags *XMLTags) GetByEmptyID() (emptyIDXMLTags XMLTags) {
	return xmlTags.GetByID("")
}

//GetByNotEmptyID 获取ID不为空的元素
func (xmlTags *XMLTags) GetByNotEmptyID() (subXMLTags XMLTags) {
	subXMLTags = XMLTags{}
	for _, elem := range *xmlTags {
		if elem.ID != "" {
			subXMLTags = append(subXMLTags, elem)
		}
	}
	return subXMLTags
}

//MergeSameIDXmlTags 合并相同Tag、相同ID的记录
func (xmlTags *XMLTags) MergeSameIDXmlTags() (newXMLTags XMLTags) {
	newXMLTags = XMLTags{}
	for tag, sameTagXmlTags := range xmlTags.GroupByTag() {
		for id, sameIDXmlTags := range sameTagXmlTags.GroupByID() {
			xmlTag := XMLTag{
				Tag:   tag,
				ID:    id,
				Attrs: Attrs{},
			}
			for _, xmlElement := range sameIDXmlTags {
				xmlTag.Attrs = append(xmlTag.Attrs, xmlElement.Attrs...)
			}
			newXMLTags = append(newXMLTags, xmlTag)
		}
	}
	return newXMLTags
}

func (xmlTags *XMLTags) GroupByID() (group map[string]XMLTags) {
	group = make(map[string]XMLTags)
	for _, elem := range *xmlTags {
		_, ok := group[elem.ID]
		if !ok {
			group[elem.ID] = make(XMLTags, 0)
		}
		group[elem.ID] = append(group[elem.ID], elem)
	}
	return group
}

func (xmlTags *XMLTags) ToMap() (out []map[string]string) {
	out = make([]map[string]string, 0)
	formatXMLTags := xmlTags.CopyEmptyID2NotEmpty()
	for _, xmlTag := range formatXMLTags.MergeSameIDXmlTags() {

		attrs := Attrs{}
		attrs = append(attrs, &Attr{
			Name:  "tag",
			Value: xmlTag.Tag,
		}, &Attr{
			Name:  "id",
			Value: xmlTag.ID,
		})
		attrs = append(attrs, xmlTag.Attrs...)
		oneRecord := attrs.ToMap()
		out = append(out, oneRecord)
	}
	return out
}

func (xmlTags *XMLTags) ToStruct(dst interface{}) (err error) {
	m := xmlTags.ToMap()
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, dst)
	if err != nil {
		return err
	}
	return
}

func parseNode(node ast.Node, source []byte) (xmlTags []XMLTag) {
	xmlTags = make([]XMLTag, 0)
	if htmlBlock, ok := node.(*ast.HTMLBlock); ok && htmlBlock.HTMLBlockType == ast.HTMLBlockType2 {
		htmlRaw := Node2RawText(htmlBlock, source)
		xmlTag, baseAttr, ok := parseRawHtml(htmlRaw)
		if ok {
			if xmlTag.Tag == "" {
				err := errors.Errorf("unexcept block tag xml:%s", htmlRaw)
				panic(err)
			}
			if baseAttr.Value == "" {
				nextNode := node.NextSibling()
				if fencedCodeNode, ok := nextNode.(*ast.FencedCodeBlock); ok {
					attr := &Attr{
						Name:  "language",
						Value: string(fencedCodeNode.Language(source)),
					}
					xmlTag.Attrs = append(xmlTag.Attrs, attr)
					value := Node2RawText(nextNode, source)
					xmlTag.Attrs = append(xmlTag.Attrs, &Attr{
						Name:  baseAttr.Name,
						Value: value,
					})
				} else if tableHTML, ok := nextNode.(*extast.Table); ok {
					if xmlTag.ID == "" {
						err := errors.Errorf("table element required id attribute")
						panic(err)
					}
					columnAttr, ok := xmlTag.Attrs.PopByName("column")
					if !ok {
						err := errors.Errorf("table element required coulmn attribute")
						panic(err)
					}
					columnArr := strings.Split(columnAttr.Value, ",")
					keyMap := map[string]string{}
					keyMapAttr, ok := xmlTag.Attrs.PopByName("keymap")
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
					i := 0
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
						rowXmlTag := xmlTag.Clone()
						rowXmlTag.ID = fmt.Sprintf("%s_%d", rowXmlTag.ID, i) // 表格行，自动增加id
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
							rowXmlTag.AddAttr(lastAttr)
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
									rowXmlTag.AddAttr(&Attr{
										Name:  key,
										Value: value,
									})
								} else if txt == "" {
									lastAttr.Value = txt
								}
							}
						}

						xmlTags = append(xmlTags, rowXmlTag)
						subNode = subNode.NextSibling()
						i++
					}
				}

			}
			xmlTags = append(xmlTags, *xmlTag)
		}

	} else if rawHTML, ok := node.(*ast.RawHTML); ok {
		txt := Node2RawText(rawHTML, source)
		xmlTag, baseAttr, ok := parseRawHtml(string(txt))
		if ok {
			if baseAttr.Value == "" {
				nextNode := node.NextSibling()
				value := Node2RawText(nextNode, source)
				xmlTag.AddAttr(&Attr{
					Name:  baseAttr.Name,
					Value: value,
				})
			}
			xmlTags = append(xmlTags, *xmlTag)
		}

	}
	if node.HasChildren() {
		firstChild := node.FirstChild()
		subAttr := parseNode(firstChild, source)
		xmlTags = append(xmlTags, subAttr...)
	}
	nextNode := node.NextSibling()
	if nextNode != nil {
		subAttr := parseNode(nextNode, source)
		xmlTags = append(xmlTags, subAttr...)
	}
	return xmlTags
}

func parseRawHtml(rawText string) (xmlTag *XMLTag, baseAttr *Attr, ok bool) {

	xmlTag = &XMLTag{
		Attrs: make([]*Attr, 0),
	}
	baseAttr = &Attr{}
	rawText = strings.TrimSpace(rawText)
	tag := strings.Trim(rawText, "<!->")
	attributeStr := ""
	index := strings.Index(tag, " ")
	if index > -1 {
		attributeStr = fmt.Sprintf("{%s}", tag[index+1:])
		tag = tag[:index]
	}
	if strings.Contains(tag, "=") {
		arr := strings.SplitN(tag, "=", 2)
		tag = arr[0]
		baseAttr.Value = arr[1]
	}
	if !strings.HasPrefix(tag, "doc.") {
		return nil, nil, false
	}
	tag = strings.TrimSpace(strings.TrimPrefix(tag, "doc."))
	if strings.Contains(tag, ".") {
		arr := strings.SplitN(tag, ".", 2)
		tag = arr[0]
		baseAttr.Name = strings.TrimSpace(arr[1])
	}

	if attributeStr != "" {
		txtReader := text.NewReader([]byte(attributeStr))
		attrs, ok := parser.ParseAttributes(txtReader)
		if !ok {
			err := errors.Errorf("convert to attribute err:%s", attributeStr)
			panic(err)
		}
		for _, parseAttr := range attrs {
			attr := Attr{
				Name:  string(parseAttr.Name),
				Value: fmt.Sprintf("%s", parseAttr.Value),
			}
			xmlTag.AddAttr(&attr)
		}
	}
	if baseAttr.Value != "" {
		xmlTag.Attrs = append(xmlTag.Attrs, baseAttr)
	}
	xmlTag.Tag = tag
	return xmlTag, baseAttr, true
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
