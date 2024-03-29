package filltemplateapidoc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/suifengpiao14/jsonschemaline"
	"github.com/tidwall/gjson"
)

// View 解析模板
func View(tplContent string, jsonStr string) (out string, err error) {
	tpl := template.New("").Funcs(FuncMap)
	tpl, err = tpl.Parse(tplContent)
	if err != nil {
		return "", err
	}
	var wr bytes.Buffer
	err = tpl.Execute(&wr, jsonStr)
	if err != nil {
		return "", err
	}
	out = wr.String()
	return out, nil
}

var FuncMap = template.FuncMap{
	"jsonGet":     JsonGet,
	"jsonExample": JsonExample,
}

// JsonGet 在模板中使用gjson 路径获取值
func JsonGet(data interface{}, path string) (out string, err error) {
	tpl := fmt.Sprintf(`{{jsonGet . "%s"}}`, path)
	out = tpl //默认为无法解析,输出原模板
	s, err := interface2string(data)
	if err != nil {
		return "", err
	}

	result := gjson.Get(s, path)
	if result.Exists() {
		str := result.String()
		if !strings.Contains(str, "{{") { // 内部不再有模板的情况,可以执行模板函数,否则保持原来模板输出
			out = str
		}
	}
	return out, nil
}

func JsonExample(data interface{}, path string) (out string, err error) {
	tpl := fmt.Sprintf(`{{jsonExample . "%s"}}`, path)
	out = tpl //默认为无法解析,输出原模板
	dataStr, err := interface2string(data)
	if err != nil {
		return "", err
	}
	result := gjson.Parse(dataStr)
	if path != "" {
		result = result.Get(path)
	}
	str := result.String()
	if str == "" {
		err = errors.Errorf("empty str got by path: %s,in data: %s", path, dataStr)
		return "", err
	}
	if strings.Contains(str, "{{") { // 内部有模板的情况,保持原来模板输出
		return out, nil
	}
	lineSchema, err := jsonschemaline.Json2lineSchema(str)
	if err != nil {
		return "", err
	}

	uniqueueMap := map[string]bool{}
	pathArr := make([]string, 0) // 使用数组保持原有顺序
	for _, item := range lineSchema.Items {
		path := item.Fullname
		lastDot := strings.LastIndex(item.Fullname, ".")
		if lastDot > -1 {
			path = item.Fullname[:lastDot]
		}
		if _, ok := uniqueueMap[path]; ok {
			continue
		}
		attrResult := gjson.Get(str, path)
		typeResult := attrResult.Get("type")
		if !typeResult.Exists() { // 一定要存在类型属性
			continue
		}
		uniqueueMap[path] = true
		pathArr = append(pathArr, path)
	}

	var w bytes.Buffer
	w.WriteString("version=http://json-schema.org/draft-07/schema#,direction=out,id=example\n")
	for _, path := range pathArr {
		attrResult := gjson.Get(str, path)
		typeResult := attrResult.Get("type")
		w.WriteString(fmt.Sprintf("fullname=%s,", path))
		format := attrResult.Get("format").String()
		defaul := attrResult.Get("default").String()
		example := attrResult.Get("example").String()
		typ := typeResult.String()
		if typ == "" {
			typ = "string"
		}
		w.WriteString(fmt.Sprintf("type=%s", typ))
		if format != "" {
			w.WriteString(fmt.Sprintf(",format=%s", format))
		}
		if defaul != "" {
			w.WriteString(fmt.Sprintf(",default=%s", defaul))
		}
		if example != "" {
			w.WriteString(fmt.Sprintf(",example=%s", example))
		}
		w.WriteString(fmt.Sprintf(",dst=%s", path)) //为了符合 lineschema规则,填充作用,此处无特殊意义

		w.WriteByte('\n')
	}
	lineschemaStr := w.String()
	lineschema, err := jsonschemaline.ParseJsonschemaline(lineschemaStr)
	if err != nil {
		return "", err
	}

	out, err = lineschema.Jsonschemaline2json()
	if err != nil {
		return "", err
	}
	out = gjson.Get(out, "@this|@pretty").String()

	return out, nil
}

func CurlExample(data interface{}, path string) (out string, err error) {
	tpl := fmt.Sprintf(`{{curlExample . "%s"}}`, path)
	out = tpl //默认为无法解析,输出原模板
	dataStr, err := interface2string(data)
	if err != nil {
		return "", err
	}
	return dataStr, nil
}

func interface2string(data interface{}) (out string, err error) {
	out, ok := data.(string)
	if !ok {
		b, err := json.Marshal(data)
		if err != nil {
			return "", err
		}
		out = string(b)
	}
	if !gjson.Valid(out) {
		err = errors.Errorf("invalid json data")
		return "", err
	}
	return out, nil
}
