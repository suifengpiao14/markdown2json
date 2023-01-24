package filltemplateapidoc

import (
	"fmt"
	"os"
	"testing"
)

func TestView(t *testing.T) {
	tplContent, err := GetTplContent()
	if err != nil {
		panic(err)
	}
	//|{{.common.head["content-type"].name}}| {{.common.head["content-type"].type}}|{{.common.head["content-type"].required}}|{{.common.head["content-type"].default}}|{{.common.head["content-type"].desc}}|
	jsonStr := `
	{
		"common":{
			"head":{
				"content-type":{
					"name":"content-type",
					"type":"string",
					"required":"是",
					"default":"application/json",
					"desc":"文件格式"
				},
				"appid":{
					"name":"appid",
					"type":"string",
					"required":"是",
					"default":"",
					"desc":"访问服务的备案id"
				},
				"signature":{
					"name":"appid",
					"type":"string",
					"required":"是",
					"default":"",
					"desc":"签名,外网访问需开启签名"
				}
			}
		}
}
	`
	//common.head["content-type"].name
	out, err := View(tplContent, jsonStr)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}

func GetTplContent() (tplContent string, err error) {
	filename := "./example/doc/adList.md"
	b, err := os.ReadFile(filename)
	tplContent = string(b)
	if err != nil {
		return "", err
	}
	return tplContent, nil
}
