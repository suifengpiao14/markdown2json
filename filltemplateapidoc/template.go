package filltemplateapidoc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

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
	"jsonGet": JsonGet,
}

// JsonGet 在模板中使用gjson 路径获取值
func JsonGet(data interface{}, path string) (out string, err error) {
	tpl := fmt.Sprintf(`{{jsonGet . "%s"}}`, path)
	out = tpl //默认为无法解析,输出原模板
	s, ok := data.(string)
	if !ok {
		b, err := json.Marshal(data)
		if err != nil {
			return "", err
		}
		s = string(b)
	}
	result := gjson.Get(s, path)
	if result.Exists() {
		out = result.String()
	}
	return out, nil
}
