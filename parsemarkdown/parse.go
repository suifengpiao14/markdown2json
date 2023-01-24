package parsemarkdown

// 只负责解析markdown 到[]*Record 格式，不负责数据整合及有效性验证
import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
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
	KEY_PRIVATE_PREFIX           = "_"
	KEY_INER_NEXT_SIBLING_COLUMN = "_nextsiblingName_"
	KEY_INER_INDEX               = "_index_"
	KEY_ID                       = "id"
	KEY_TAG                      = "_tag"
	KEY_IS_END_BY_BACKSLASH      = "_isEndByBackslash"
	KEY_TEXT                     = "_text"
	KEY_COLUMN                   = "_column"
	KEY_DB                       = "db"
	KEY_TABLE                    = "table"
	KEY_OFFSET                   = "_offset" //内联元素指定截取字符串位置
	KEY_LENGTH                   = "_length" //内联元素指定截取字符串长度
	ID_SEPARATOR                 = "-"
	KEY_REF                      = "_ref"
	KEY_INER_REF                 = "_ref_" // 内部记录来源,方便出错时,提示信息更有正对性
)

const (
	STRING_TRUE  = "true"
	STRING_FALSE = "false"
)

func Parse(source []byte) (records Records, err error) {
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
	records, err = parseNode(document, source)
	if err != nil {
		return nil, err
	}
	records, err = ResolveRef(records)
	if err != nil {
		return nil, err
	}
	return records, nil
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Record []*KV

type Records []Record

func (records Records) MoveInternalKey() (newRecords Records) {
	newRecords = make(Records, 0)
	for _, record := range records {
		newRecord := record.MoveInternalKey()
		newRecords = append(newRecords, newRecord)
	}
	return newRecords
}
func (records Records) GetRefs() (refRecords Records) {
	refRecords = make(Records, 0)
	for _, record := range records {
		_, ok := record.GetKV(KEY_REF)
		if ok {
			refRecords = append(refRecords, record)
		}

	}
	return refRecords
}

func (records Records) Format() (newRecords Records, err error) {
	newRecords = make(Records, 0)
	group := map[string][]Record{} // 按照tag 分组
	for _, record := range records {
		tag := record.GetTag()
		_, ok := group[tag]
		if !ok {
			group[tag] = make([]Record, 0)
		}
		group[tag] = append(group[tag], record)
	}
	// 相同tag 合并
	for _, sameTagRecords := range group {
		l := len(sameTagRecords)
		record := sameTagRecords[l-1]
		for i := l - 2; i > -1; i-- {
			prevRecord := sameTagRecords[i]
			for _, kv := range prevRecord {
				record.AddKV(*kv)
			}
		}
		newRecords = append(newRecords, record)
	}
	return newRecords, nil
}

func (records Records) FilterByKV(kv KV) (subRecords Records) {
	subRecords = make(Records, 0)
	for _, record := range records {
		ekv, ok := record.GetKV(kv.Key)
		if ok && ekv.Value == kv.Value {
			subRecords = append(subRecords, record)
		}
	}
	return subRecords
}

func (records Records) Filter(fn func(record Record) bool) (subRecords Records) {
	subRecords = make(Records, 0)
	for _, record := range records {
		if fn(record) {
			subRecords = append(subRecords, record)
		}
	}
	return subRecords
}

func (records Records) Walk(fn func(record Record) Record) (subRecords Records) {
	subRecords = make(Records, 0)
	for _, record := range records {
		newRecord := fn(record)
		subRecords = append(subRecords, newRecord)
	}
	return subRecords
}

func (records Records) First() (record Record) {
	if len(records) > 0 {
		record = records[0]
	}
	return record
}

func (records Records) Json() (out string) {
	// 所有属性,转换为绝对名称
	l := len(records)
	for i := l - 1; i > -1; i-- { // 倒序,相同key时,最早出现的覆盖后面出现的
		record := records[i]
		tagName := record.GetTag()
		for _, kv := range record {
			if kv.Key == KEY_IS_END_BY_BACKSLASH || kv.Key == KEY_REF || (kv.Key == KEY_TEXT && kv.Value == "") {
				continue
			}
			if kv.Key != KEY_TAG {
				kv.Key = fmt.Sprintf("%s.%s", tagName, kv.Key)
			}
			var err error
			var key = kv.Key
			result := gjson.Get(out, key)
			if result.IsArray() {
				key = fmt.Sprintf("%s.-1", kv.Key)
			}
			out, err = sjson.Set(out, key, kv.Value)
			if err != nil {
				fmt.Println(out)
				panic(err)
			}
		}
	}
	return out
}

func (records Records) GetByIndex(index string) (newRecords Records) {
	newRecords = make(Records, 0)
	for _, record := range records {
		if record.GetIndex() == index {
			newRecords = append(newRecords, record)
		}
	}
	return newRecords
}

func (records Records) GetByTag(tag string) (newRecords Records) {
	newRecords = make(Records, 0)
	for _, record := range records {
		if record.GetTag() == tag {
			newRecords = append(newRecords, record)
		}
	}
	return newRecords
}

func (records Records) GetByTagWithChildren(tag string) (newRecords Records) {
	newRecords = make(Records, 0)
	for _, record := range records {
		if strings.HasPrefix(record.GetTag(), tag) {
			newRecords = append(newRecords, record)
		}
	}
	return newRecords
}

func RecordError(record Record, err error) error {
	idAttr, ok := record.GetKV(KEY_ID)
	if ok {
		err = errors.WithMessagef(err, "id: %s", idAttr.Value)
	}
	innerRefAttr, ok := record.GetKV(KEY_INER_REF)
	if ok {
		err = errors.WithMessagef(err, "ref: %s", innerRefAttr.Value)
	}
	return err
}

// MergeRecords 将多条记录中的kv，按保留最早出现的原则，合并成一条
func MergeRecords(records ...Record) (newRecord Record, err error) {
	kvMap := map[string]*KV{}
	breakInherit := false
	for _, record := range records {
		for _, kv := range record {
			if strings.HasPrefix(kv.Key, KEY_PRIVATE_PREFIX) { // 过滤私有属性，私有属性不继承
				continue
			}
			okv, ok := kvMap[kv.Key]
			if !ok { // 不存在，直接填充后跳过
				kvMap[kv.Key] = kv
				continue
			}
			isArr := strings.HasSuffix(kv.Key, "[]")
			if !isArr {
				continue // 非数组，跳过
			}
			okv.Value = fmt.Sprintf("%s,%s", okv.Value, kv.Value)
			kvMap[kv.Key] = okv
		}
		if breakInherit {
			break // 后续属性，不再继承
		}
	}
	newRecord = Record{}
	for _, kv := range kvMap {
		newRecord.AddKV(*kv)
	}
	return newRecord, nil
}

func (record *Record) AddKV(kv KV) {
	*record = append(*record, &kv) // 首先添加
}

func (record *Record) IsEmptyKey(key string) (empty bool) {
	kv, exists := record.GetKV(key)
	if !exists {
		empty = true
		return
	}
	empty = kv.Value == ""
	return empty
}

func (record *Record) IsExists(key string) (exists bool) {
	_, exists = record.GetKV(key)
	return exists
}

func (record *Record) SetNotExistsKV(kv KV) {
	_, exists := record.GetKV(kv.Key)
	if exists {
		return
	}
	*record = append(*record, &kv) // 首先添加
}
func (record *Record) ResetKV(kv KV) {
	exists := false
	okv, exists := record.GetKV(kv.Key)
	if exists {
		okv.Value = kv.Value
		return
	}
	*record = append(*record, &kv)
}

func (record Record) String() (out string) {
	newRecord := record.MoveInternalKey()
	mp := make(map[string]string)
	for _, kv := range newRecord {
		mp[kv.Key] = kv.Value
	}
	b, err := json.Marshal(mp)
	if err != nil {
		panic(err)
	}
	out = string(b)
	return out
}

// MoveInternalKey 删除内部使用的KV
func (record *Record) MoveInternalKey() (new Record) {
	newRecord := Record{}
	for _, kv := range *record {
		switch kv.Key {
		case KEY_INER_NEXT_SIBLING_COLUMN, KEY_INER_INDEX, KEY_LENGTH, KEY_OFFSET, KEY_INER_REF: // 删除内部使用的KV
		case KEY_COLUMN, KEY_REF: // 删除内部使用的KV
		default:
			newRecord = append(newRecord, kv)
		}
	}
	return newRecord
}

// GetIndex  获取记录的_index_
func (record *Record) GetIndex() (index string) {
	if indexAttr, ok := record.GetKV(KEY_INER_INDEX); ok {
		return indexAttr.Value
	}
	return ""
}

// GetID 获取记录的ID
func (record *Record) GetID() (index string) {
	if idAttr, ok := record.GetKV(KEY_ID); ok {
		return idAttr.Value
	}
	return ""
}

// GetID 获取记录的ID
func (record *Record) GetTag() (index string) {
	if tagAttr, ok := record.GetKV(KEY_TAG); ok {
		return tagAttr.Value
	}
	return ""
}

// 克隆记录
func (record *Record) Clone() (newRecord Record) {
	newRecord = Record{}
	for _, kv := range *record {
		newRecord.AddKV(*kv)
	}
	return newRecord
}

func (record *Record) GetKV(key string) (kv *KV, ok bool) {
	for _, kv := range *record {
		if kv.Key == key {
			return kv, true
		}
	}
	return nil, false
}

func (record *Record) GetValue(key string) (value string) {
	for _, kv := range *record {
		if kv.Key == key {
			return kv.Value
		}
	}
	return ""
}

func (record *Record) PopKV(key string) (popKV *KV, ok bool) {
	newRecord := Record{}
	for _, kv := range *record {
		if kv.Key == key {
			popKV = kv
		} else {
			newRecord = append(newRecord, kv)
		}
	}
	*record = newRecord // 替换原有的
	if popKV == nil {
		return nil, false
	}
	return popKV, true
}

func (record *Record) IsNotEndBlacslash() (yes bool) {
	kv, ok := record.GetKV(KEY_IS_END_BY_BACKSLASH)
	if !ok {
		return false
	}
	yes = kv.Value != STRING_TRUE
	return yes
}

// GetIndex  获取记录的 父类key
func GetParentIndex(index string) (parentIndex string) {
	p := strings.LastIndex(index, ID_SEPARATOR)
	if p > -1 {
		return index[:p]
	}
	return ""
}
func CloneTabHeader(record Record) Record { // 表格元素需要把db、table 等基本属性复制到子元素
	newRecord := Record{}
	for _, kv := range record {
		switch kv.Key {
		case KEY_COLUMN, KEY_ID, KEY_REF, KEY_INER_NEXT_SIBLING_COLUMN: // 这些属性不复制
			continue
		default:
			newRecord.AddKV(*kv)
		}
	}
	return newRecord
}
func SetNextSiblingValue(nextNode ast.Node, record *Record, records *Records, source []byte) (err error) {
	// 处理代码块元素
	if fencedCodeNode, ok := nextNode.(*ast.FencedCodeBlock); ok {
		attr := &KV{
			Key:   "language",
			Value: string(fencedCodeNode.Language(source)),
		}
		record.ResetKV(*attr)
		value := Node2RawText(nextNode, source)

		record.AddKV(KV{
			Key:   KEY_TEXT,
			Value: value, // 修改标签名称的值
		})
		return nil
	}
	// 处理表格元素
	if tableHTML, ok := nextNode.(*extast.Table); ok {
		columnKey := KEY_COLUMN
		columnAttr, ok := record.GetKV(columnKey)
		if !ok {
			err = errors.Errorf("table element required %s attribute", KEY_COLUMN)
			return err
		}
		columnArr := strings.Split(columnAttr.Value, ",")
		firstNode := tableHTML.FirstChild()
		headNode, ok := firstNode.(*extast.TableHeader)
		if !ok {
			err = errors.Errorf("first children is not header")
			return err
		}
		columnLen := len(columnArr)
		if columnLen != headNode.ChildCount() {
			err = errors.Errorf("column filed not match table head field._column:%s,ref:", strings.Join(columnArr, ","))
			return err
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
				return err
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
				value := string(cellNode.Text(source))

				if cellIndex == 0 {
					tableTagName := newRecord.GetTag()
					columnTagName := fmt.Sprintf("%s.%s", tableTagName, value)
					newRecord.ResetKV(KV{
						Key:   KEY_TAG,
						Value: columnTagName,
					})
				}
				kv = KV{
					Key:   name,
					Value: value,
				}
				newRecord.ResetKV(kv)
				cellNode = cellNode.NextSibling()
				cellIndex++
			}
			*records = append(*records, newRecord)
			subNode = subNode.NextSibling()
			i++
		}
		return
	}
	// 其它元素
	value := Node2RawText(nextNode, source)
	lenAttr, ok := record.GetKV(KEY_LENGTH)
	if ok {
		l, err := strconv.Atoi(lenAttr.Value)
		if err != nil {
			err = errors.WithMessagef(err, "convert %s attr err, rawHTML:%s", KEY_LENGTH, value)
			return err
		}
		if l > len(value) {
			err = errors.WithMessagef(err, "value length less then  %d, rawHTML:%s", l, value)
			return err
		}
		value = value[:l+1]
	}
	record.AddKV(KV{
		Key:   KEY_TEXT,
		Value: value,
	})
	return nil
}

func parseNode(node ast.Node, source []byte) (records Records, err error) {
	records = Records{}
	if htmlBlock, ok := node.(*ast.HTMLBlock); ok {
		switch htmlBlock.HTMLBlockType {
		case ast.HTMLBlockType2:
			htmlRaw := Node2RawText(htmlBlock, source)
			record, err := rawHtml2Record(htmlRaw)
			if err != nil {
				return nil, err
			}
			if record.IsNotEndBlacslash() && record.IsEmptyKey(KEY_TEXT) { // 非闭合标签,名称的值为空,则取相邻下一个元素为值
				nextNode := node.NextSibling()
				err = SetNextSiblingValue(nextNode, &record, &records, source)
				if err != nil {
					return nil, err
				}
			}
			records = append(records, record)
		case ast.HTMLBlockType7: // 自关闭自定义标签，后面需要加空行，如果存在这种情况，抛出错误，提示增加空行
			segments := htmlBlock.Lines()
			for i := 0; i < segments.Len(); i++ {
				segment := segments.At(i)
				lineByte := segment.Value(source)
				lineByte = bytes.TrimSpace(lineByte)
				if bytes.HasPrefix(lineByte, []byte("<!--")) {
					err := errors.Errorf("rawHtml:(%s) is  ast.HTMLBlockType7 type ,and must end with blank line .see https://spec.commonmark.org/0.30/#html-blocks point 7", Node2RawText(htmlBlock, source))
					return nil, err
				}

			}

		}

	} else if rawHTML, ok := node.(*ast.RawHTML); ok { // 内联元素
		txt := Node2RawText(rawHTML, source)
		record, err := rawHtml2Record(txt)
		if err != nil {
			return nil, err
		}
		if record.IsNotEndBlacslash() { // 标签名称的值为空,则取相邻下一个元素为值
			nextNode := node.NextSibling()
			SetNextSiblingValue(nextNode, &record, &records, source)
		}
		records = append(records, record)

	}
	if node.HasChildren() {
		firstChild := node.FirstChild()
		subRecords, err := parseNode(firstChild, source)
		if err != nil {
			return nil, err
		}
		records = append(records, subRecords...)
	}
	nextNode := node.NextSibling()
	if nextNode != nil {
		subRecords, err := parseNode(nextNode, source)
		if err != nil {
			return nil, err
		}
		records = append(records, subRecords...)
	}

	return records, nil
}

const (
	ARRAY_SUFFIX = "__array__"
)

func rawHtml2Record(rawText string) (record Record, err error) {
	record = Record{}
	rawText = strings.TrimSpace(rawText)
	//解决以注释开头的行,后续带有值的情况 "<!--doc.service.description-->记录公司spuID和阿里spuID关联关系，以及状态同步"
	endPos := strings.Index(rawText, ">")
	if endPos != len(rawText)-1 {
		siblingValue := rawText[endPos+1:]
		rawText = rawText[:endPos+1]
		record.AddKV(KV{
			Key:   KEY_TEXT,
			Value: siblingValue,
		})
	}
	rawText = strings.Trim(rawText, "<!-> ")
	tagName := strings.Trim(rawText, "/")
	lastIndex := len(rawText) - 1
	endBackslash := STRING_FALSE
	if rawText[lastIndex] == '/' {
		rawText = rawText[:lastIndex]
		endBackslash = STRING_TRUE
		// /结束不会再设置 text,此处设置默认kv
		record.ResetKV(KV{
			Key:   KEY_TEXT,
			Value: "",
		})
	}
	record.ResetKV(KV{
		Key:   KEY_IS_END_BY_BACKSLASH,
		Value: endBackslash,
	})

	index := strings.Index(rawText, " ")
	otherAttrStr := ""
	if index > -1 {
		tagName = rawText[:index]
		otherAttrStr = strings.TrimSpace(rawText[index:])
	}
	record.ResetKV(KV{
		Key:   KEY_TAG,
		Value: tagName,
	})
	if otherAttrStr == "" {
		return record, nil
	}
	// 标签名称之外还有其它属性
	otherAttrStr = fmt.Sprintf("{%s}", otherAttrStr)
	txtReader := text.NewReader([]byte(otherAttrStr))
	attrs, ok := parser.ParseAttributes(txtReader)
	if !ok {
		err := errors.Errorf("convert to attribute err:%s,rawTxt:%s", otherAttrStr, rawText)
		return nil, err
	}
	for _, parseAttr := range attrs {
		value := ""
		rv := reflect.Indirect(reflect.ValueOf(parseAttr.Value))
		switch rv.Kind() {
		case reflect.String:
			value = rv.String()
		case reflect.Bool:
			value = strconv.FormatBool(rv.Bool())
		case reflect.Float64:
			value = strconv.FormatFloat(rv.Float(), 'f', 0, 64)
		case reflect.Float32:
			value = strconv.FormatFloat(rv.Float(), 'f', 0, 32)
		case reflect.Int, reflect.Int64:
			value = strconv.FormatInt(rv.Int(), 10)
		default:
			value = fmt.Sprintf("%s", parseAttr.Value)
		}
		attr := KV{
			Key:   string(parseAttr.Name),
			Value: value,
		}
		if strings.HasSuffix(attr.Key, ARRAY_SUFFIX) { // 数组元素替换为原样
			attr.Key = fmt.Sprintf("%s[]", strings.TrimSuffix(attr.Key, ARRAY_SUFFIX))
		}
		record.ResetKV(attr)
	}

	return record, nil
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
