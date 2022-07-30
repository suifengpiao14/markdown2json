package markdown2json

// 只负责解析markdown 到[]*Record 格式，不负责数据整合及有效性验证
import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
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
	KEY_INER_NEXT_SIBLING_COLUMN = "_nextsiblingName_"
	KEY_INER_INDEX               = "_index_"
	KEY_ID                       = "id"
	KEY_COLUMN                   = "column"
	KEY_DB                       = "db"
	KEY_TABLE                    = "table"
	KEY_OFFSET                   = "_offset" //内联元素指定截取字符串位置
	KEY_LENGTH                   = "_length" //内联元素指定截取字符串长度
	ID_SEPARATOR                 = "-"
	KEY_REF                      = "ref"
	KEY_INER_REF                 = "_ref_"    // 内部记录来源,方便出错时,提示信息更有正对性
	KEY_INHERIT                  = "_inherit" // 是否基础其它相同id的属性(公共参数有时需明确指出不继承其它优先级标签的更多属性)
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
	return records, nil
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Record []*KV

type Records []*Record

func (records Records) MoveInternalKey() (newRecords Records) {
	newRecords = make(Records, 0)
	for _, record := range records {
		newRecord := record.MoveInternalKey()
		newRecords = append(newRecords, &newRecord)
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
	group := map[string]Records{} // 按照index 分组
	for _, record := range records {
		index := record.GetIndex()
		_, ok := group[index]
		if !ok {
			group[index] = make(Records, 0)
		}
		group[index] = append(group[index], record)
	}
	// 相同index 合并
	tmpNewRecords := Records{}
	for _, sameIndexRecords := range group {
		newRecord, err := MergeRecords(sameIndexRecords...)
		if err != nil {
			return nil, err
		}
		tmpNewRecords = append(tmpNewRecords, &newRecord)
	}
	// 合并父类
	for _, record := range tmpNewRecords {
		mergeRecord := Records{}
		mergeRecord = append(mergeRecord, record)
		index := record.GetIndex()
		parentIndex := GetParentIndex(index)
		for {
			if parentIndex == "" {
				break
			}
			sameIndexParents := tmpNewRecords.GetByIndex(parentIndex)
			mergeRecord = append(mergeRecord, sameIndexParents...)
			parentIndex = GetParentIndex(parentIndex)
		}
		newRecord, err := MergeRecords(mergeRecord...)
		if err != nil {
			return nil, err
		}
		newRecords = append(newRecords, &newRecord)
	}
	return newRecords, nil
}

func (records Records) Filter(kv KV) (subRecords Records) {
	subRecords = make(Records, 0)
	for _, record := range records {
		eKv, ok := record.GetKV(kv.Key)
		if ok && eKv.Value == kv.Value {
			subRecords = append(subRecords, record)
		}
	}
	return subRecords
}

func (records Records) String() (out string) {
	newRecords, err := records.Format()
	if err != nil {
		panic(err)
	}
	newRecords = newRecords.MoveInternalKey()
	arr := make([]map[string]string, 0)
	for _, record := range newRecords {
		mp := make(map[string]string)
		for _, kv := range *record {
			mp[kv.Key] = kv.Value
		}
		arr = append(arr, mp)
	}
	b, err := json.Marshal(arr)
	if err != nil {
		panic(err)
	}
	out = string(b)
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

func (records Records) GetByIndexWithChildren(index string) (newRecords Records) {
	newRecords = make(Records, 0)
	for _, record := range records {
		if strings.HasPrefix(record.GetIndex(), index) {
			newRecords = append(newRecords, record)
		}
	}
	return newRecords
}

func RecordError(record Record, err error) error {
	idAttr, ok := record.GetKV(KEY_DB)
	if ok {
		err = errors.WithMessagef(err, "id: %s", idAttr.Value)
	}
	innerRefAttr, ok := record.GetKV(KEY_INER_REF)
	if ok {
		err = errors.WithMessagef(err, "ref: %s", innerRefAttr.Value)
	}
	return err
}

//MergeRecords 将多条记录中的kv，按保留最早出现的原则，合并成一条
func MergeRecords(records ...*Record) (newRecord Record, err error) {
	kvMap := map[string]*KV{}
	for _, record := range records {
		inheritAttr, ok := record.GetKV(KEY_INHERIT)
		if ok {
			bol, err := strconv.ParseBool(inheritAttr.Value)
			if err != nil {
				err = RecordError(*record, err)
				return nil, err
			}
			if !bol {

				return *record, nil
			}
		}
		for _, kv := range *record {
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
	}
	newRecord = Record{}
	for _, kv := range kvMap {
		newRecord = append(newRecord, kv)
	}
	return newRecord, nil
}

func (record *Record) AddKV(kv KV) {
	*record = append(*record, &kv) // 首先添加
	//针对db、table、id 属性特殊处理,3个全部设置好后生成 _index_属性
	if !(kv.Key == KEY_DB || kv.Key == KEY_TABLE || kv.Key == KEY_ID) {
		return
	}
	// add index
	dbValue := ""
	tableValue := ""
	idValue := ""
	if db, ok := record.GetKV(KEY_DB); ok {
		dbValue = db.Value
	}
	if table, ok := record.GetKV(KEY_TABLE); ok {
		tableValue = table.Value
	}
	if id, ok := record.GetKV(KEY_ID); ok {
		idValue = id.Value
	}
	index := fmt.Sprintf("%s%s%s%s%s", dbValue, ID_SEPARATOR, tableValue, ID_SEPARATOR, idValue)
	index = strings.Trim(index, ID_SEPARATOR)
	record.ResetKV(KV{
		Key:   KEY_INER_INDEX,
		Value: index,
	})

}

//MoveInternalKey 删除内部使用的KV
func (record *Record) MoveInternalKey() (new Record) {
	newRecord := Record{}
	for _, kv := range *record {
		switch kv.Key {
		case KEY_INER_NEXT_SIBLING_COLUMN, KEY_INER_INDEX, KEY_LENGTH, KEY_OFFSET: // 删除内部使用的KV
		case KEY_COLUMN, KEY_REF, KEY_INHERIT: // 删除内部使用的KV
		default:
			newRecord = append(newRecord, kv)
		}
	}
	return newRecord
}

//GetIndex  获取记录的key
func (record *Record) GetIndex() (index string) {
	if indexAttr, ok := record.GetKV(KEY_INER_INDEX); ok {
		return indexAttr.Value
	}
	return ""
}

// 克隆基本信息
func (record *Record) Clone() (newRecord Record) {
	newRecord = Record{}
	dbAttr, ok := record.GetKV(KEY_DB)
	if ok {
		newRecord = append(newRecord, dbAttr)
	}
	tableAttr, ok := record.GetKV(KEY_TABLE)
	if ok {
		newRecord = append(newRecord, tableAttr)
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
func (record *Record) PopKV(key string) (popKV *KV, ok bool) {
	newRecord := Record{}
	for _, kv := range *record {
		if kv.Key == key {
			popKV = kv
		} else {
			newRecord = append(newRecord, newRecord...)
		}
	}
	return popKV, false
}

func (record *Record) ResetKV(kv KV) {
	exists := false
	for _, okv := range *record {
		if okv.Key == kv.Key {
			exists = true
			okv.Value = kv.Value
		}
	}
	if !exists {
		*record = append(*record, &kv)
	}
}

//GetIndex  获取记录的 父类key
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
		case KEY_COLUMN, KEY_ID, KEY_REF: // 这些属性不复制
			continue
		default:
			newRecord.AddKV(*kv)
		}
	}
	return newRecord
}
func SetNextSiblingValue(nextNode ast.Node, record *Record, records *Records, source []byte) (err error) {
	nextsiblingAttr, ok := record.GetKV(KEY_INER_NEXT_SIBLING_COLUMN)
	if !ok {
		err = errors.Errorf("record attr %s required", KEY_INER_NEXT_SIBLING_COLUMN)
		return err
	}
	// 处理代码块元素
	if fencedCodeNode, ok := nextNode.(*ast.FencedCodeBlock); ok {
		attr := &KV{
			Key:   "language",
			Value: string(fencedCodeNode.Language(source)),
		}
		record.AddKV(*attr)
		value := Node2RawText(nextNode, source)
		record.AddKV(KV{
			Key:   nextsiblingAttr.Value,
			Value: value,
		})
		return nil
	}
	// 处理表格元素
	if tableHTML, ok := nextNode.(*extast.Table); ok {
		idAttr, ok := record.GetKV(KEY_ID)
		if !ok {
			err = errors.Errorf("table element required %s attribute", KEY_ID)
			return err
		}
		columnAttr, ok := record.GetKV(KEY_COLUMN)
		if !ok {
			err = errors.Errorf("table element required %s attribute", KEY_COLUMN)
			return err
		}
		columnArr := strings.Split(columnAttr.Value, ",")
		firstColumnName := columnArr[0]
		firstNode := tableHTML.FirstChild()
		headNode, ok := firstNode.(*extast.TableHeader)
		if !ok {
			err = errors.Errorf("first children is not header")
			return err
		}
		columnLen := len(columnArr)
		if columnLen != headNode.ChildCount() {
			err = errors.Errorf("column filed not match table head field.column:%s,ref:", strings.Join(columnArr, ","))
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
				value := cellNode.Text(source)
				kv = KV{
					Key:   name,
					Value: string(value),
				}
				newRecord.AddKV(kv)
				cellNode = cellNode.NextSibling()
				cellIndex++
			}
			idValue := fmt.Sprintf("%s%s%d", idAttr.Value, ID_SEPARATOR, i)
			firstAttr, ok := newRecord.GetKV(firstColumnName)
			if ok {
				idValue = fmt.Sprintf("%s%s%s", idAttr.Value, ID_SEPARATOR, firstAttr.Value)
			}
			idKV := KV{
				Key:   KEY_ID,
				Value: idValue,
			}
			newRecord.AddKV(idKV)
			*records = append(*records, &newRecord)
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
		Key:   nextsiblingAttr.Value,
		Value: value,
	})

	return nil
}

func parseNode(node ast.Node, source []byte) (records Records, err error) {
	records = Records{}
	if htmlBlock, ok := node.(*ast.HTMLBlock); ok && htmlBlock.HTMLBlockType == ast.HTMLBlockType2 {
		htmlRaw := Node2RawText(htmlBlock, source)
		record, err := rawHtml2Record(htmlRaw)
		if err != nil {
			return nil, err
		}
		records = append(records, record) // 先保存记录,后续可以通过引用修改当前record 部分值
		_, exists := record.GetKV(KEY_INER_NEXT_SIBLING_COLUMN)
		if exists {
			nextNode := node.NextSibling()
			SetNextSiblingValue(nextNode, record, &records, source)
		}

	} else if rawHTML, ok := node.(*ast.RawHTML); ok { // 内联元素
		txt := Node2RawText(rawHTML, source)
		record, err := rawHtml2Record(txt)
		if err != nil {
			return nil, err
		}

		_, exists := record.GetKV(KEY_INER_NEXT_SIBLING_COLUMN)
		if exists {
			nextNode := node.NextSibling()
			SetNextSiblingValue(nextNode, record, &records, source)
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

var (
	regArrWithValue  = regexp.MustCompile(`(\[\]) *=`)
	replArrWithValue = fmt.Sprintf(`%s=`, ARRAY_SUFFIX)
	regArrWithName   = regexp.MustCompile(`([\w\.]+)\[\] `)
	replArrWithName  = fmt.Sprintf(`$1%s`, ARRAY_SUFFIX)

	regWithValue  = regexp.MustCompile(`(\w+)\.(\w+)\.(\w+)=(.*)`)
	replWithValue = fmt.Sprintf(`%s=$1 %s=$2 $3=$4`, KEY_DB, KEY_TABLE)

	regWithName  = regexp.MustCompile(`(\w+)\.(\w+)\.(\w+)(.*)`)
	replWithName = fmt.Sprintf(`%s=$1 %s=$2 %s=$3 $4`, KEY_DB, KEY_TABLE, KEY_INER_NEXT_SIBLING_COLUMN)

	regTableName  = regexp.MustCompile(`(\w+)\.(\w+)(.*)`)
	replTableName = fmt.Sprintf(`%s=$1 %s=$2 %s="" $3`, KEY_DB, KEY_TABLE, KEY_INER_NEXT_SIBLING_COLUMN)

	regDBName  = regexp.MustCompile(`(\w+)(.*)`)
	replDBName = fmt.Sprintf(`%s=$1 %s="" $2`, KEY_DB, KEY_INER_NEXT_SIBLING_COLUMN)
)

func FormatRawText(s string) string {
	s = regArrWithValue.ReplaceAllString(s, replArrWithValue)
	s = regArrWithName.ReplaceAllString(s, replArrWithName)
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

func rawHtml2Record(rawText string) (record *Record, err error) {
	record = &Record{}
	rawText = strings.Trim(rawText, "<!-/>")
	rawText = strings.TrimSpace(rawText)
	formatText := FormatRawText(rawText)

	attrStr := fmt.Sprintf("{%s}", formatText)
	txtReader := text.NewReader([]byte(attrStr))
	attrs, ok := parser.ParseAttributes(txtReader)
	if !ok {
		err := errors.Errorf("convert to attribute err:%s,rawTxt:%s", attrStr, rawText)
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
		record.AddKV(attr)
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
