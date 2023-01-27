package xml

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/suifengpiao14/markdown2json/parsemarkdown"
)

func TestXml2Data(t *testing.T) {
	filename := "../example/doc/adList.md"
	b, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	s := string(b)
	records, err := Xml2Data(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(records.Json())
}

func TestParseSubEncoding(t *testing.T) {
	recordStr := `[{"key":"_tag","value":"Parameter"},{"key":"ns","value":"doc.parameter.requestParamter"},{"key":"column","value":"name,type,required,description,default,example"},{"key":"position","value":"body"},{"key":"encoding","value":"markdown/table"},{"key":"_text","value":"\n\u003c!--doc.parameter.requestParamter column=\"name,type,required,description,default,example\"  position=body--\u003e\n|参数名|类型|必选|说明|默认值|示例|\n|:----    |:---|:----- |-----   |-----   |-----   |\n|title| string|是|广告标题||新年豪礼|\n|advertiserId| string|是|广告主||123|\n|beginAt| string|是|可以投放开始时间||2023-01-12 00:00:00|\n|endAt| string|是|投放结束时间||2023-01-30 00:00:00|\n|index| string|是|页索引,0开始|0||\n|size| string|是|每页数量|10||\n"}]`
	record := &parsemarkdown.Record{}
	err := json.Unmarshal([]byte(recordStr), record)
	if err != nil {
		panic(err)
	}
	records, err := ParseSubEncoding(record)
	if err != nil {
		panic(err)
	}
	fmt.Println(records)
}
