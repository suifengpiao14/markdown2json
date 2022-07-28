package markdown2json

// 只负责解析markdown 到[]*Record 格式，不负责数据整合及有效性验证
import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
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

const (
	NEXT_SIBLING_COLUMN_KEY = "_nextsiblingName_"
	KEY_ID                  = "id"
	KEY_COLUMN              = "column"
	KEY_DB                  = "db"
	KEY_TABLE               = "table"
	KEY_OFFSET              = "_offset" //内联元素指定截取字符串位置
	KEY_LENGTH              = "_length" //内联元素指定截取字符串长度
)

func Parse(source []byte) (records []*Record, err error) {
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
	records = parseNode(document, source)
	return records, nil
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Record []*KV

func (record *Record) AdddKV(kv KV) {
	*record = append(*record, &kv)
}

// 克隆基本信息
func (record *Record) Clone() (newRecord Record) {
	newRecord = Record{}
	dbAttr, ok := record.GetKVFirst(KEY_DB)
	if ok {
		newRecord = append(newRecord, dbAttr)
	}
	tableAttr, ok := record.GetKVFirst(KEY_TABLE)
	if ok {
		newRecord = append(newRecord, tableAttr)
	}
	return newRecord
}

func (record *Record) GetKVFirst(key string) (kv *KV, ok bool) {
	for _, kv := range *record {
		if kv.Key == key {
			return kv, true
		}
	}
	return nil, false
}

func CloneTabHeader(record Record) Record {
	newRecord := Record{}
	for _, kv := range record {
		if kv.Key == KEY_COLUMN || kv.Key == KEY_ID {
			continue
		}
		newKv := *kv
		newRecord = append(newRecord, &newKv)
	}
	return newRecord
}

func parseNode(node ast.Node, source []byte) (records []*Record) {
	records = make([]*Record, 0)
	if htmlBlock, ok := node.(*ast.HTMLBlock); ok && htmlBlock.HTMLBlockType == ast.HTMLBlockType2 {
		htmlRaw := Node2RawText(htmlBlock, source)
		record, ok := rawHtml2Record(htmlRaw)
		if ok {
			nextsiblingAttr, exists := record.GetKVFirst(NEXT_SIBLING_COLUMN_KEY)
			if exists {
				nextNode := node.NextSibling()
				if fencedCodeNode, ok := nextNode.(*ast.FencedCodeBlock); ok {
					attr := &KV{
						Key:   "language",
						Value: string(fencedCodeNode.Language(source)),
					}
					record.AdddKV(*attr)
					value := Node2RawText(nextNode, source)
					record.AdddKV(KV{
						Key:   nextsiblingAttr.Value,
						Value: value,
					})
				} else if tableHTML, ok := nextNode.(*extast.Table); ok {
					idAttr, ok := record.GetKVFirst(KEY_ID)
					if !ok {
						err := errors.Errorf("table element required %s attribute", KEY_ID)
						panic(err)
					}
					columnAttr, ok := record.GetKVFirst(KEY_COLUMN)
					if !ok {
						err := errors.Errorf("table element required %s attribute", KEY_COLUMN)
						panic(err)
					}
					columnArr := strings.Split(columnAttr.Value, ",")
					firstColumnName := columnArr[0]
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
						newRecord := CloneTabHeader(*record)
						var kv = KV{}
						for {
							if cellNode == nil {
								break
							}
							name := columnArr[cellIndex]
							value := cellNode.Text(source)
							kv = KV{
								Key:   name,
								Value: string(value),
							}
							newRecord.AdddKV(kv)
							cellNode = cellNode.NextSibling()
							cellIndex++
						}
						idValue := fmt.Sprintf("%s_%d", idAttr.Value, i)
						firstAttr, ok := newRecord.GetKVFirst(firstColumnName)
						if ok {
							idValue = fmt.Sprintf("%s_%s", idAttr.Value, firstAttr.Value)
						}
						idKV := KV{
							Key:   KEY_ID,
							Value: idValue,
						}
						newRecord.AdddKV(idKV)
						records = append(records, &newRecord)
						subNode = subNode.NextSibling()
						i++
					}
				}
			} else {
				records = append(records, record)
			}

		}

	} else if rawHTML, ok := node.(*ast.RawHTML); ok { // 内联元素
		txt := Node2RawText(rawHTML, source)
		record, ok := rawHtml2Record(txt)
		if ok {
			nextsiblingAttr, exists := record.GetKVFirst(NEXT_SIBLING_COLUMN_KEY)
			if exists {
				nextNode := node.NextSibling()
				value := Node2RawText(nextNode, source)
				lenAttr, ok := record.GetKVFirst(KEY_LENGTH)
				if ok {
					l, err := strconv.Atoi(lenAttr.Value)
					if err != nil {
						err = errors.WithMessagef(err, "convert %s attr err, rawHTML:%s", KEY_LENGTH, txt)
						panic(err)
					}
					if l > len(value) {
						err = errors.WithMessagef(err, "value length less then  %d, rawHTML:%s", l, txt)
						panic(err)
					}
					value = value[:l+1]
				}
				record.AdddKV(KV{
					Key:   nextsiblingAttr.Value,
					Value: value,
				})
			}
			records = append(records, record)
		}
	}
	if node.HasChildren() {
		firstChild := node.FirstChild()
		subRecords := parseNode(firstChild, source)
		records = append(records, subRecords...)
	}
	nextNode := node.NextSibling()
	if nextNode != nil {
		subRecords := parseNode(nextNode, source)
		records = append(records, subRecords...)
	}

	for _, record := range records {
		newRecord := Record{}
		for _, kv := range *record {
			switch kv.Key {
			case NEXT_SIBLING_COLUMN_KEY: // 删除内部使用的KV
			default:
				newRecord = append(newRecord, kv)
			}
		}
		*record = newRecord
	}

	return records
}

var (
	regWithValue  = regexp.MustCompile(`(\w+)\.(\w+)\.(\w+)=(.*)`)
	replWithValue = fmt.Sprintf(`%s=$1 %s=$2 $3=$4`, KEY_DB, KEY_TABLE)

	regWithName  = regexp.MustCompile(`(\w+)\.(\w+)\.(\w+)(.*)`)
	replWithName = fmt.Sprintf(`%s=$1 %s=$2 %s=$3 $4`, KEY_DB, KEY_TABLE, NEXT_SIBLING_COLUMN_KEY)

	regTableName  = regexp.MustCompile(`(\w+)\.(\w+)(.*)`)
	replTableName = fmt.Sprintf(`%s=$1 %s=$2 %s="" $3`, KEY_DB, KEY_TABLE, NEXT_SIBLING_COLUMN_KEY)

	regDBName  = regexp.MustCompile(`(\w+)(.*)`)
	replDBName = fmt.Sprintf(`%s=$1 %s="" $2`, KEY_DB, NEXT_SIBLING_COLUMN_KEY)
)

func FormatRawText(s string) string {
	if regWithValue.MatchString(s) {
		return regWithValue.ReplaceAllString(s, replWithValue)
	}
	if regWithName.MatchString(s) {
		return regWithName.ReplaceAllString(s, replWithName)
	}
	if regTableName.MatchString(s) {
		return regTableName.ReplaceAllString(s, replTableName)
	}
	if regDBName.MatchString(s) {
		return regDBName.ReplaceAllString(s, replDBName)
	}
	return s
}

func rawHtml2Record(rawText string) (record *Record, ok bool) {
	record = &Record{}
	rawText = strings.Trim(rawText, "<!-/>")
	rawText = strings.TrimSpace(rawText)
	formatText := FormatRawText(rawText)

	attrStr := fmt.Sprintf("{%s}", formatText)
	txtReader := text.NewReader([]byte(attrStr))
	attrs, ok := parser.ParseAttributes(txtReader)
	if !ok {
		err := errors.Errorf("convert to attribute err:%s,rawTxt:%s", attrStr, rawText)
		panic(err)
	}
	for _, parseAttr := range attrs {
		attr := KV{
			Key:   string(parseAttr.Name),
			Value: fmt.Sprintf("%s", parseAttr.Value),
		}
		record.AdddKV(attr)
	}
	return record, true
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
