# 新增广告
***说明***: <Attr id="title">新增广告</Attr>

***请求路径***: <Attr id="path">/admin/v1/ad/add</Attr>

***请求方法***: <Attr id="method">POST</Attr>
<Ref file="./example/doc/common.md#doc.request.header" id="common.header"/>

***请求头***:
{{Ref .common.header}}

***请求参数***:

<Ref file= "./adList.md#doc.parameter.response.items[]" id="obj" />
<Set  id="doc.parameter.request" column="name,type,required,description,default,example" position="body" encoding="markdown/table">

|参数名|类型|必选|说明|默认值|示例|
|:----    |:---|:----- |-----   |-----   |----   |
|{{Strtr .obj.title ".name|string|是|.description|.default|.example"}}|
|{{Strtr .obj.title ".advertiserId|string|是|.description|.default|.example"}}|
|{{Strtr .obj.title ".summary|string|是|.description|.default|.example"}}|
|{{Strtr .obj.title ".image|string|是|.description|.default|.example"}}|
|{{Strtr .obj.title ".link|string|是|.description|.default|.example"}}|
|{{Strtr .obj.title ".type|string|是|.description|.default|.example"}}|
|{{Strtr .obj.title ".beginAt|string|是|.description|.default|.example"}}|
|{{Strtr .obj.title ".endAt|string|是|.description|.default|.example"}}|
|{{Strtr .obj.title ".remark|string|是|.description|.default|.example"}}|
|{{Strtr .obj.title ".valueObj|string|是|.description|.default|.example"}}|

**请求示例：**
```json
{{jsonExample .doc.parameter.request -}}
``` 
**返回参数：**
<Set id="doc.parameter.response" encoding="markdown/table" column="name,type,description,example" >
| 参数名                | 参数类型 | 描述             | 示例                      |
| --------------------- | -------- | ---------------- | ------------------------- |
|code                  | string   | 业务状态码         | -                         |
| message   | string   | 业务提示           | -                         |
| data               | object | 对象         | -                        |
| data.id|string |新增广告ID标识|0| 
</Set>

**返回示例(正常)：**
```json 
{{jsonExample . "doc.parameter.response" -}}
``` 
***返回示例(错误)***:
<Ref file="./example/doc/common.md#doc.example.response.error" id="common.response.error"/>
{{jsonGet . "common.response.error.text"}}