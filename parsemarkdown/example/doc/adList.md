# 广告列表
***说明***: 广告列表

***请求路径***: <Uri>/admin/v1/ad/add</Uri>

***请求方法***: POST
<Ref file="./example/sql/ad.sql#ad" name="db.ad"/>
<!--db.ad  _ref="file:///./example/sql/ad.sql#ad"  -->
<Ref file="./example/doc/common.md#doc.parameter" namespace="common.head"/>
<!--common.head _ref="file:///./example/doc/common.md#doc.parameter"-->

***请求头***:
|参数名|类型|必选|默认值|说明|
|:----    |:---|:----- |-----   |-----   |
|{{jsonGet . "common.head.content-type.name"}}| {{jsonGet . "common.head.content-type.type"}}|{{jsonGet . "common.head.content-type.required"}}|{{jsonGet . "common.head.content-type.default"}}|{{jsonGet . "common.head.content-type.description"}}|
|{{jsonGet . "common.head.appid.name"}}| {{jsonGet . "common.head.appid.type"}}|{{jsonGet . "common.head.appid.required"}}|{{jsonGet . "common.head.appid.default"}}|{{jsonGet . "common.head.appid.description"}}|
|{{jsonGet . "common.head.signature.name"}}| {{jsonGet . "common.head.signature.type"}}|{{jsonGet . "common.head.signature.required"}}|{{jsonGet . "common.head.signature.default"}}|{{jsonGet . "common.head.signature.description"}}|


***请求参数***:
<Parameter  namespace="doc.parameter.requestParamter" column="name,type,required,description,default,example" position="body" encoding="markdown/table" >
<!--doc.parameter.requestParamter _column="name,type,required,description,default,example"  position=body-->
|参数名|类型|必选|说明|默认值|示例|
|:----    |:---|:----- |-----   |-----   |-----   |
|title| string|是|广告标题||新年豪礼|
|advertiserId| string|是|广告主||123|
|beginAt| string|是|可以投放开始时间||2023-01-12 00:00:00|
|endAt| string|是|投放结束时间||2023-01-30 00:00:00|
|index| string|是|页索引,0开始|0||
|size| string|是|每页数量|10||
</Parameter>

**请求示例：**
```json
{{jsonExample . "doc.parameter.requestParamter" -}}
``` 
<Parameter namespace="doc.parameter.response" encoding="markdown/table" column="name,type,description,example" position="body"  httpSttus="200">

**返回参数：**
<!--doc.parameter.responseParameter position=body httpStatus="200" _column="name,type,description,example"-->
| 参数名                | 参数类型 | 描述             | 示例                      |
| --------------------- | -------- | ---------------- | ------------------------- |
|code                  | string   | 业务状态码         | 0                         |
| message   | string   | 业务提示           | ok                         |
| items               | array | 数组         | -                        |
|items[].{{dbRefCamel . "db.ad.advertise.id.name"}}|string |{{jsonGet . "db.ad.advertise.id.comment"}}|0| 
| <!--map _value="db.ad.advertise.title"-->items[].title|string |{{jsonGet . "db.ad.advertise.title.comment"}}|新年好礼| 
| <!--map _value="db.ad.advertise.advertiser_id"-->items[].advertiserId|string |{{jsonGet . "db.ad.advertise.advertiser_id.comment"}}|123| 
| <!--map _value="db.ad.advertise.summary"-->items[].summary|string |{{jsonGet . "db.ad.advertise.summary.comment"}}|下单有豪礼| 
| <!--map _value="db.ad.advertise.image"-->items[].image|string |{{jsonGet . "db.ad.advertise.image.comment"}}|http://image.service.cn/new_year.jpg"| 
| <!--map _value="db.ad.advertise.link"-->items[].link|string |{{jsonGet . "db.ad.advertise.link.comment"}}|http://gift.servcice.cn/new_year_git| 
| <!--map _value="db.ad.advertise.type"-->items[].type|string |{{jsonGet . "db.ad.advertise.type.comment"}}|image| 
| <!--map _value="db.ad.advertise.beginAt"-->items[].beginAt|string |{{jsonGet . "db.ad.advertise.begin_at.comment"}}|2023-01-12 00:00:00| 
| <!--map _value="db.ad.advertise.endAt"-->items[].endAt|string |{{jsonGet . "db.ad.advertise.end_at.comment"}}|2023-01-30 00:00:00| 
| <!--map _value="db.ad.advertise.remark"-->items[].remark|string |{{jsonGet . "db.ad.advertise.remark.comment"}}|营养早餐广告| 
| <!--map _value="db.ad.advertise.valueObj"-->items[].valueObj|string |{{jsonGet . "db.ad.advertise.value_obj.comment"}}|值对象| 
| pagination|object |对象|| 
| <!--map _value="doc.parameter.requestParamter.index"-->pagination.index|string |{{jsonGet . "doc.parameter.requestParamter.index.description"}}|{{jsonGet . "doc.parameter.requestParamter.index.default"}}| 
| <!--map _value="doc.parameter.requestParamter.size"-->pagination.size|string |{{jsonGet . "doc.parameter.requestParamter.size.description"}}|{{jsonGet . "doc.parameter.requestParamter.size.default"}}| 
| pagination.total|string |总数|60| 

</Parameter>

<Example namespace="api.example.response.200" encoding="markdown/code">

**返回示例(正常)：**
```json 
{{jsonExample . "doc.parameter.responseParameter" -}}
``` 
</Example>
<!--common.response.example.error _ref="file:///./example/doc/common.md#doc.response.example.error" /-->
**返回示例(错误)：**
{{jsonGet . "common.response.example.error" -}}