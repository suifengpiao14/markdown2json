# 广告列表
***说明***: <Line id="title">广告列表</Line>

***请求路径***: <Line id="path">/admin/v1/ad/add</Line>

***请求方法***: <Line id="method">POST</Line>
<Ref file="./example/doc/common.md#doc.request.header" id="common.header"/>

***请求头***:
{{jsonGet . "common.header.text"}}

***请求参数***:
<Ref file="./example/sql/ad.sql#ad" id="db.ad"/>
<Set  id="doc.parameter.request" column="name,type,required,description,default,example" position="body" encoding="markdown/table" >

|参数名|类型|必选|说明|默认值|示例|
|:----    |:---|:----- |-----   |-----   |-----   |
|title| string|是|{{jsonGet . "db.ad.advertise.title.comment"}}||新年豪礼|
|advertiserId| string|是|广告主||123|
|beginAt| string|是|可以投放开始时间||2023-01-12 00:00:00|
|endAt| string|是|投放结束时间||2023-01-30 00:00:00|
|index| string|是|页索引,0开始|0||
|size| string|是|每页数量|10||
</Set>

**请求示例：**
```json
{{jsonExample . "doc.parameter.request" -}}
``` 
<Set id="doc.parameter.response" encoding="markdown/table" column="name,type,description,example" >

**返回参数：**
| 参数名                | 参数类型 | 描述             | 示例                      |
| --------------------- | -------- | ---------------- | ------------------------- |
|code                  | string   | 业务状态码         | 0                         |
| message   | string   | 业务提示           | ok                         |
| items               | array | 数组         | -                        |
|items[].id|string |{{jsonGet . "db.ad.advertise.id.comment"}}|0| 
|items[].title|string |{{jsonGet . "db.ad.advertise.title.comment"}}|{{jsonGet . "doc.parameter.request.title.example"}}| 
|items[].advertiserId|string |{{jsonGet . "db.ad.advertise.advertiser_id.comment"}}|{{jsonGet . "doc.parameter.request.advertiserId.example"}}| 
|items[].summary|string |{{jsonGet . "db.ad.advertise.summary.comment"}}|下单有豪礼| 
|items[].image|string |{{jsonGet . "db.ad.advertise.image.comment"}}|http://image.service.cn/new_year.jpg"| 
|items[].link|string |{{jsonGet . "db.ad.advertise.link.comment"}}|http://gift.servcice.cn/new_year_git| 
|items[].type|string |{{jsonGet . "db.ad.advertise.type.comment"}}|image| 
|items[].beginAt|string |{{jsonGet . "db.ad.advertise.begin_at.comment"}}|2023-01-12 00:00:00| 
|items[].endAt|string |{{jsonGet . "db.ad.advertise.end_at.comment"}}|2023-01-30 00:00:00| 
|items[].remark|string |{{jsonGet . "db.ad.advertise.remark.comment"}}|营养早餐广告| 
|items[].valueObj|string |{{jsonGet . "db.ad.advertise.value_obj.comment"}}|值对象| 
| pagination|object |对象|| 
|pagination.index|string |{{jsonGet . "doc.parameter.requestParamter.index.description"}}|{{jsonGet . "doc.parameter.requestParamter.index.default"}}| 
| pagination.size|string |{{jsonGet . "doc.parameter.requestParamter.size.description"}}|{{jsonGet . "doc.parameter.requestParamter.size.default"}}| 
| pagination.total|string |总数|60| 

</Set>

<Obj id="doc.example.response.200" encoding="markdown/code">

**返回示例(正常)：**
```json 
{{jsonExample . "doc.parameter.response" -}}
``` 
</Obj>

***返回示例(错误)***:
<Ref file="./example/doc/common.md#doc.example.response.error" id="common.response.error"/>
{{Ref .common.response.error }}
