package markdown2json_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	parsemarkdown "github.com/suifengpiao14/markdown2json"
)

func TestParse(t *testing.T) {
	records := GetRecords()
	b, err := json.Marshal(records)
	if err != nil {
		panic(err)
	}
	str := string(b)
	fmt.Println(str)
}

func GetRecords() parsemarkdown.Records {
	//file := "./example/first-doc.mdx"
	file := "./example/spuUpdateQuestion.md"
	fd, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	source, err := io.ReadAll(fd)
	if err != nil {
		panic(err)
	}
	records, err := parsemarkdown.Parse(source)
	if err != nil {
		panic(err)
	}
	return records
}

func TestGetRefs(t *testing.T) {
	records := GetRecords()
	refRecords := records.GetRefs()
	parsemarkdown.ResolveRef(refRecords)
	fmt.Println(refRecords.String())
}

func TestMerge(t *testing.T) {
	records := GetRecords()
	newRecords, err := parsemarkdown.MergeRecords(records...)
	if err != nil {
		panic(err)
	}
	b, err := json.Marshal(newRecords)
	if err != nil {
		panic(err)
	}
	out := string(b)
	fmt.Println(out)
}
func TestRecordString(t *testing.T) {
	records := GetRecords()
	out := records.FilterByKV(parsemarkdown.KV{Key: parsemarkdown.KEY_DB, Value: "doc"}).String()
	fmt.Println(out)
}

func TestFormat(t *testing.T) {
	str := `[[{"key":"db","value":"doc"},{"key":"_index_","value":"doc-parameter-requestParamter"},{"key":"table","value":"parameter"},{"key":"_ref","value":"file:///D:\\\\go\\\\markdown2json\\\\example\\\\commonArgs.md#requestParamter"},{"key":"position","value":"body"},{"key":"id","value":"requestParamter"}],[{"key":"db","value":"doc"},{"key":"_index_","value":"doc-parameter-requestParamter"},{"key":"table","value":"parameter"},{"key":"_nextsiblingName_","value":""},{"key":"id","value":"requestParamter"},{"key":"prefix","value":"_param"},{"key":"position","value":"body"},{"key":"_column","value":"name,type,required,default,description"}],[{"key":"db","value":"doc"},{"key":"_index_","value":"doc-parameter-requestParamter-scene"},{"key":"_index_","value":"doc-parameter-requestParamter-scene"},{"key":"table","value":"parameter"},{"key":"prefix","value":"_param"},{"key":"position","value":"body"},{"key":"name","value":"scene"},{"key":"type","value":"string"},{"key":"required","value":"是"},{"key":"default","value":"-"},{"key":"description","value":"场景枚举值:3C(3C),3C_NEW(3C_NEW)"},{"key":"id","value":"requestParamter-scene"}],[{"key":"db","value":"doc"},{"key":"_index_","value":"doc-parameter-requestParamter-Fxy_spuid"},{"key":"_index_","value":"doc-parameter-requestParamter-Fxy_spuid"},{"key":"table","value":"parameter"},{"key":"prefix","value":"_param"},{"key":"position","value":"body"},{"key":"name","value":"Fxy_spuid"},{"key":"type","value":"string"},{"key":"required","value":"是"},{"key":"default","value":"-"},{"key":"description","value":"闲鱼SPU ID 格式: number(数字类型)"},{"key":"id","value":"requestParamter-Fxy_spuid"}],[{"key":"db","value":"doc"},{"key":"_index_","value":"doc-parameter-requestParamter-Fxy_product_name"},{"key":"_index_","value":"doc-parameter-requestParamter-Fxy_product_name"},{"key":"table","value":"parameter"},{"key":"prefix","value":"_param"},{"key":"position","value":"body"},{"key":"name","value":"Fxy_product_name"},{"key":"type","value":"string"},{"key":"required","value":"是"},{"key":"default","value":"-"},{"key":"description","value":"闲鱼SPU 名称"},{"key":"id","value":"requestParamter-Fxy_product_name"}],[{"key":"db","value":"doc"},{"key":"_index_","value":"doc-parameter-requestParamter-Fhsb_productid"},{"key":"_index_","value":"doc-parameter-requestParamter-Fhsb_productid"},{"key":"table","value":"parameter"},{"key":"prefix","value":"_param"},{"key":"position","value":"body"},{"key":"name","value":"Fhsb_productid"},{"key":"type","value":"string"},{"key":"required","value":"是"},{"key":"default","value":"-"},{"key":"description","value":"回收宝产品ID格式: number(数字类型)"},{"key":"id","value":"requestParamter-Fhsb_productid"}],[{"key":"db","value":"doc"},{"key":"_index_","value":"doc-parameter-requestParamter"},{"key":"table","value":"parameter"},{"key":"_nextsiblingName_","value":""},{"key":"id","value":"requestParamter"},{"key":"_column","value":"name,type,required,example,description"},{"key":"prefix","value":"_head"},{"key":"position","value":"body"}],[{"key":"db","value":"doc"},{"key":"_index_","value":"doc-parameter-requestParamter-_version"},{"key":"_index_","value":"doc-parameter-requestParamter-_version"},{"key":"table","value":"parameter"},{"key":"prefix","value":"_head"},{"key":"position","value":"body"},{"key":"name","value":"_version"},{"key":"type","value":"string"},{"key":"required","value":"是"},{"key":"example","value":"0.01"},{"key":"description","value":"协议版本号可选值:0.01"},{"key":"id","value":"requestParamter-_version"}],[{"key":"db","value":"doc"},{"key":"_index_","value":"doc-parameter-requestParamter-_msgType"},{"key":"_index_","value":"doc-parameter-requestParamter-_msgType"},{"key":"table","value":"parameter"},{"key":"prefix","value":"_head"},{"key":"position","value":"body"},{"key":"name","value":"_msgType"},{"key":"type","value":"string"},{"key":"required","value":"是"},{"key":"example","value":"request"},{"key":"description","value":"报文类型可选值:request(请求)、response(响应)"},{"key":"id","value":"requestParamter-_msgType"}],[{"key":"db","value":"doc"},{"key":"_index_","value":"doc-parameter-requestParamter-_timestamps"},{"key":"_index_","value":"doc-parameter-requestParamter-_timestamps"},{"key":"table","value":"parameter"},{"key":"prefix","value":"_head"},{"key":"position","value":"body"},{"key":"name","value":"_timestamps"},{"key":"type","value":"string"},{"key":"required","value":"是"},{"key":"example","value":"1523330331"},{"key":"description","value":"请求时间戳(单位毫秒)"},{"key":"id","value":"requestParamter-_timestamps"}],[{"key":"db","value":"doc"},{"key":"_index_","value":"doc-parameter-requestParamter-_invokeId"},{"key":"_index_","value":"doc-parameter-requestParamter-_invokeId"},{"key":"table","value":"parameter"},{"key":"prefix","value":"_head"},{"key":"position","value":"body"},{"key":"name","value":"_invokeId"},{"key":"type","value":"string"},{"key":"required","value":"是"},{"key":"example","value":"book1523330331358"},{"key":"description","value":"当前请求标识(每次请求要求唯一)"},{"key":"id","value":"requestParamter-_invokeId"}],[{"key":"db","value":"doc"},{"key":"_index_","value":"doc-parameter-requestParamter-_callerServiceId"},{"key":"_index_","value":"doc-parameter-requestParamter-_callerServiceId"},{"key":"table","value":"parameter"},{"key":"prefix","value":"_head"},{"key":"position","value":"body"},{"key":"name","value":"_callerServiceId"},{"key":"type","value":"string"},{"key":"required","value":"是"},{"key":"example","value":"110001"},{"key":"description","value":"发起http请求方的服务ID"},{"key":"id","value":"requestParamter-_callerServiceId"}],[{"key":"db","value":"doc"},{"key":"_index_","value":"doc-parameter-requestParamter-_groupNo"},{"key":"_index_","value":"doc-parameter-requestParamter-_groupNo"},{"key":"table","value":"parameter"},{"key":"prefix","value":"_head"},{"key":"position","value":"body"},{"key":"name","value":"_groupNo"},{"key":"type","value":"string"},{"key":"required","value":"是"},{"key":"example","value":"1"},{"key":"description","value":"请求分组号"},{"key":"id","value":"requestParamter-_groupNo"}],[{"key":"db","value":"doc"},{"key":"_index_","value":"doc-parameter-requestParamter-_interface"},{"key":"_index_","value":"doc-parameter-requestParamter-_interface"},{"key":"table","value":"parameter"},{"key":"prefix","value":"_head"},{"key":"position","value":"body"},{"key":"name","value":"_interface"},{"key":"type","value":"string"},{"key":"required","value":"是"},{"key":"example","value":"templateList"},{"key":"description","value":"请求接口标识"},{"key":"id","value":"requestParamter-_interface"}],[{"key":"db","value":"doc"},{"key":"_index_","value":"doc-parameter-requestParamter-_remark"},{"key":"_index_","value":"doc-parameter-requestParamter-_remark"},{"key":"table","value":"parameter"},{"key":"prefix","value":"_head"},{"key":"position","value":"body"},{"key":"name","value":"_remark"},{"key":"type","value":"string"},{"key":"required","value":"是"},{"key":"example","value":"0.01"},{"key":"description","value":"备注"},{"key":"id","value":"requestParamter-_remark"}]]`
	records := parsemarkdown.Records{}
	err := json.Unmarshal([]byte(str), &records)
	if err != nil {
		panic(err)
	}
	newRecords, err := records.Format()
	if err != nil {
		panic(err)
	}
	fmt.Println(newRecords.String())
}
