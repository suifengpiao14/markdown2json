package xml

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"

	"github.com/antchfx/xmlquery"
	"github.com/pkg/errors"
	"github.com/suifengpiao14/markdown2json/parsemarkdown"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

func Xml2Data(xml string) (records parsemarkdown.Records, err error) {
	doc, err := xmlquery.Parse(strings.NewReader(xml))
	if err != nil {
		return nil, err
	}

	allElement := xmlquery.Find(doc, "//*")
	records = make(parsemarkdown.Records, 0)
	for _, element := range allElement {
		record := parsemarkdown.Record{}
		tag := element.Data
		record.SetKV(parsemarkdown.KV{
			Key:   parsemarkdown.KEY_TAG,
			Value: tag,
		})
		for _, attr := range element.Attr {
			record.SetKV(parsemarkdown.KV{
				Key:   attr.Name.Local,
				Value: attr.Value,
			})
		}
		element.SetAttr("xml:space", "preserve") // 保留文本中的空格——维持原样,方便后续再解析
		record.SetKV(parsemarkdown.KV{
			Key:   parsemarkdown.KEY_TEXT,
			Value: element.OutputXML(false),
		})

		records = append(records, record)
		subRecords, err := ParseSubEncoding(&record)
		if err != nil {
			return nil, err
		}
		records = append(records, subRecords...)
	}
	return records, nil
}

func ParseSubEncoding(record *parsemarkdown.Record) (subRecords parsemarkdown.Records, err error) {
	encoding, ok := record.GetKV(parsemarkdown.KEY_ENCODING)
	if !ok {
		return nil, nil
	}
	markdownTxtKV, ok := record.GetKV(parsemarkdown.KEY_TEXT)
	if !ok || markdownTxtKV.Value == "" {
		err = errors.Errorf("markdown text required")
		return nil, err
	}
	source := []byte(markdownTxtKV.Value)
	node, err := getMarkdownDocument(source)
	if err != nil {
		return nil, err
	}
	switch encoding.Value {
	case "markdown/table":
		{
			subRecords, err := ParseMarkdownTable(source, node, *record)
			if err != nil {
				return nil, err
			}
			return subRecords, nil
		}
	case "markdown/code":
		{
			subRecord := parseMarkdownCode(source, node)
			*record = append(*record, subRecord...)
			return nil, nil
		}
	}

	return subRecords, nil

}

func ParseMarkdownTable(source []byte, node ast.Node, record parsemarkdown.Record) (records parsemarkdown.Records, err error) {
	columnKV, ok := record.GetKV(parsemarkdown.KEY_COLUMN)
	if !ok {
		err = errors.Errorf(`encodiig="markdown/table" require attribute %s`, parsemarkdown.KEY_COLUMN)
		return nil, err
	}
	columnArr := strings.Split(columnKV.Value, ",")
	// 处理表格元素
	if tableHTML, ok := node.(*extast.Table); ok {
		firstNode := tableHTML.FirstChild()
		headNode, ok := firstNode.(*extast.TableHeader)
		if !ok {
			err = errors.Errorf("first children is not header")
			return nil, err
		}
		columnLen := len(columnArr)
		if columnLen != headNode.ChildCount() {
			err = errors.Errorf("column filed not match table head field._column:%s,ref:", strings.Join(columnArr, ","))
			return nil, err
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
				err = errors.Errorf("subNode must be tableRow")
				return nil, err
			}
			cellIndex := 0
			cellNode := tableRow.FirstChild()
			newRecord := cloneTabHeader(record)
			var kv = parsemarkdown.KV{}
			for {
				if cellNode == nil {
					break
				}
				name := columnArr[cellIndex]
				value := string(cellNode.Text(source))

				if cellIndex == 0 {
					parentNamespace := record.GetNamespace()
					ns := fmt.Sprintf("%s.%s", parentNamespace, value)
					newRecord.SetKV(parsemarkdown.KV{
						Key:   parsemarkdown.KEY_NAMESPACE,
						Value: ns,
					})
				}
				kv = parsemarkdown.KV{
					Key:   name,
					Value: value,
				}
				newRecord.SetKV(kv)
				cellNode = cellNode.NextSibling()
				cellIndex++
			}
			records = append(records, newRecord)
			subNode = subNode.NextSibling()
			i++
		}
		return records, nil
	}

	if node.HasChildren() {
		firstChild := node.FirstChild()
		subRecords, err := ParseMarkdownTable(source, firstChild, record)
		if err != nil {
			return nil, err
		}
		records = append(records, subRecords...)
	}
	nextNode := node.NextSibling()
	if nextNode != nil {
		subRecords, err := ParseMarkdownTable(source, nextNode, record)
		if err != nil {
			return nil, err
		}
		records = append(records, subRecords...)
	}
	return records, nil
}

func getMarkdownDocument(source []byte) (document ast.Node, err error) {
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
	document = md.Parser().Parse(reader)
	return
}

func cloneTabHeader(record parsemarkdown.Record) parsemarkdown.Record { // 表格元素需要把db、table 等基本属性复制到子元素
	newRecord := parsemarkdown.Record{}
	for _, kv := range record {
		switch kv.Key {
		case parsemarkdown.KEY_COLUMN, parsemarkdown.KEY_ID, parsemarkdown.KEY_REF: // 这些属性不复制
			continue
		default:
			newRecord.SetKV(*kv)
		}
	}
	return newRecord
}

func parseMarkdownCode(source []byte, node ast.Node) (record parsemarkdown.Record) {
	record = parsemarkdown.Record{}
	if fencedCodeNode, ok := node.(*ast.FencedCodeBlock); ok {
		attr := &parsemarkdown.KV{
			Key:   "language",
			Value: string(fencedCodeNode.Language(source)),
		}
		record.SetKV(*attr)
		value := Node2RawText(node, source)
		record.SetKV(parsemarkdown.KV{
			Key:   parsemarkdown.KEY_TEXT,
			Value: value, // 修改标签名称的值
		})
		return record
	}
	if node.HasChildren() {
		firstChild := node.FirstChild()
		subRecord := parseMarkdownCode(source, firstChild)
		record = append(record, subRecord...)
	}
	nextNode := node.NextSibling()
	if nextNode != nil {
		subRecord := parseMarkdownCode(source, nextNode)
		record = append(record, subRecord...)
	}
	return record
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
