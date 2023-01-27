# 新增广告
***说明***: <Line id="title">新增广告</Line>

***请求路径***: <Line id="path">/admin/v1/ad/add</Line>

***请求方法***: <Line id="method">POST</Line>
<Ref file="./example/doc/common.md#doc.request.header" id="common.header"/>

***请求头***:
{{jsonGet . "common.header.text"}}

***请求参数***:

<Table  id="doc.parameter.request" column="name,type,required,description,default,example" position="body" encoding="markdown/table" ref.obj.file="./adList.md#doc.parameter.response.items[]" ref.obj.map="name:name,type:type,required:是,description:description,default:-,example:example" >

|参数名|类型|必选|说明|默认值|示例|
|:----    |:---|:----- |-----   |-----   |----   |
|{{getJsonObj "obj.title"}}|
|{{getJsonObj "obj.advertiserId"}}|
|{{getJsonObj "obj.summary"}}|
|{{getJsonObj "obj.image"}}|
|{{getJsonObj "obj.link"}}|
|{{getJsonObj "obj.type"}}|
|{{getJsonObj "obj.beginAt"}}|
|{{getJsonObj "obj.endAt"}}|
|{{getJsonObj "obj.remark"}}|
|{{getJsonObj "obj.valueObj"}}|

**请求示例：**
```json
{{jsonExample . "doc.parameter.request" -}}
``` 
**返回参数：**
<Table id="doc.parameter.response" encoding="markdown/table" column="name,type,description,example" >
| 参数名                | 参数类型 | 描述             | 示例                      |
| --------------------- | -------- | ---------------- | ------------------------- |
|code                  | string   | 业务状态码         | -                         |
| message   | string   | 业务提示           | -                         |
| data               | object | 对象         | -                        |
| data.id|string |新增广告ID标识|0| 
</Table>

**返回示例(正常)：**
```json 
{{jsonExample . "doc.parameter.response" -}}
``` 
***返回示例(错误)***:
<Ref file="./example/doc/common.md#doc.example.response.error" id="common.response.error"/>
{{jsonGet . "common.response.error.text"}}