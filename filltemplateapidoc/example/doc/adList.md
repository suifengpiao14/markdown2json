# 广告列表
***说明***: 广告列表

***请求路径***: /admin/v1/ad/add

***请求方法***: POST

***请求头***:
|参数名|类型|必选|默认值|说明|
|:----    |:---|:----- |-----   |-----   |
|{{jsonGet . "common.head.content-type.name"}}| {{jsonGet . "common.head.content-type.type"}}|{{jsonGet . "common.head.content-type.required"}}|{{jsonGet . "common.head.content-type.default"}}|{{jsonGet . "common.head.content-type.desc"}}|
|{{jsonGet . "common.head.appid.name"}}| {{jsonGet . "common.head.appid.type"}}|{{jsonGet . "common.head.appid.required"}}|{{jsonGet . "common.head.appid.default"}}|{{jsonGet . "common.head.appid.desc"}}|
|{{jsonGet . "common.head.signature.name"}}| {{jsonGet . "common.head.signature.type"}}|{{jsonGet . "common.head.signature.required"}}|{{jsonGet . "common.head.signature.default"}}|{{jsonGet . "common.head.signature.desc"}}|


***请求参数***:
<!--doc.parameter.requestParamter _column="name,type,required,description,default"  position=body-->
|参数名|类型|必选|说明|默认值|
|:----    |:---|:----- |-----   |-----   |
|title| string|是|广告标题||
|advertiserId| string|是|广告主||
|beginAt| string|是|可以投放开始时间||
|endAt| string|是|投放结束时间||
|index| string|是|页索引,0开始||
|size| string|是|每页数量||


**请求示例：**
```json
{
    "title" : "新年豪礼",
    "advertiserId":"123",
    "beginAt":"2023-01-12 00:00:00",
    "endAt":"2023-01-30 00:00:00"
}
``` 
**返回参数：**
| 参数名                | 参数类型 | 描述             | 示例                      |
| --------------------- | -------- | ---------------- | ------------------------- |
<!--doc.parameter.responseParameter position=body httpStatus="200" _column="name,type,description,example"-->
|code                  | string   | 业务状态码         | -                         |
| message   | string   | 业务提示           | -                         |
| items               | array | 数组         | -                        |
| <!--map _rev="db.ad.advertise.id"-->items[].id|string |{{jsonGet . "db.ad.advertise.id.comment"}}|0| 
| <!--map _rev="db.ad.advertise.title"-->items[].title|string |{{jsonGet . "db.ad.advertise.title.comment"}}|新年好礼| 
| <!--map _rev="db.ad.advertise.advertiserId"-->items[].advertiserId|string |{{jsonGet . "db.ad.advertise.advertiserId.comment"}}|123| 
| <!--map _rev="db.ad.advertise.summary"-->items[].summary|string |{{jsonGet . "db.ad.advertise.summary.comment"}}|下单有豪礼| 
| <!--map _rev="db.ad.advertise.image"-->items[].image|string |{{jsonGet . "db.ad.advertise.image.comment"}}|http://image.service.cn/new_year.jpg"| 
| <!--map _rev="db.ad.advertise.link"-->items[].link|string |{{jsonGet . "db.ad.advertise.link.comment"}}|http://gift.servcice.cn/new_year_git| 
| <!--map _rev="db.ad.advertise.type"-->items[].type|string |{{jsonGet . "db.ad.advertise.type.comment"}}|image| 
| <!--map _rev="db.ad.advertise.beginAt"-->items[].beginAt|string |{{jsonGet . "db.ad.advertise.beginAt.comment"}}|2023-01-12 00:00:00| 
| <!--map _rev="db.ad.advertise.endAt"-->items[].endAt|string |{{jsonGet . "db.ad.advertise.endAt.comment"}}|2023-01-30 00:00:00| 
| <!--map _rev="db.ad.advertise.remark"-->items[].remark|string |{{jsonGet . "db.ad.advertise.remark.comment"}}|营养早餐广告| 
| <!--map _rev="db.ad.advertise.valueObj"-->items[].valueObj|string |{{jsonGet . "db.ad.advertise.valueObj.comment"}}|值对象| 
| pagination|object |对象|| 
| <!--map _rev="doc.parameter.requestParamter.index"-->pagination.index|string |{{jsonGet . "doc.parameter.requestParamter.index.description"}}|| 
| <!--map _rev="doc.parameter.requestParamter.size"-->pagination.size|string |{{jsonGet . "doc.parameter.requestParamter.size.description"}}|10| 
| pagination.total|string |总数|60| 


**返回示例：**
```json 
{
"code":"0",
"message":"ok",
"items":[
    {{jsonGet . "doc.parameter.responseParameter.items.0"}}
],
"pagination":{{jsonGet . "doc.parameter.responseParameter.pagination"}}
}
``` 