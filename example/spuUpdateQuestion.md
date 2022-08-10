# 手动获取spu问卷
**简要描述：**
> 手动获取spu问卷
**协议：**
- 级别：二层
- 路径：<!--doc.api.uri--> /spu/ajaxSpuUpdateQuestion
- 请求方法 :<!--doc.api.method --> POST

**环境：**
- 开发： <!--doc.server.proxy id=dev description="开发环境"-->193.112.197.63 http://opms.huishoubao.com.cn
- 测试：<!--doc.server.proxy id=test description="测试环境"-->xx.xx.xx.xx  http://opms.huishoubao.com.cn
- 线上：<!--doc.server.url-->http://opms.huishoubao.com.cn
<!--doc.parameter._ref="file:///D:/go/markdown2json/example/commonArgs.md#requestParamter" position=body id=requestParamter-->
**请求参数：**
<!--doc.parameter id=requestParamter prefix="_param" position=body _column="name,type,required,default,description"-->
|参数名|类型|必选|默认值|说明|
|:----    |:---|:----- |-----   |-----   |
|scene| string|是|-|场景<br/>枚举值:3C(3C),3C_NEW(3C_NEW)|
|Fxy_spuid| string|是|-|闲鱼SPU ID <br/>格式: number(数字类型)|
|Fxy_product_name| string|是|-|闲鱼SPU 名称|
|Fhsb_productid| string|是|-|回收宝产品ID<br/>格式: number(数字类型) |

<Attr name="hello world" />

<!--doc.server.variables=signature:joenebfhefeh -->
<!--doc.service._ref="file:///D:/go/markdown2json/example/readme.md"-->

<!--doc.service.preRequest-->
```javascript
var serviceId=options.variables.filter(function(row){return row.name=="serviceId"})[0]?.value;
options.headers['HSB-OPENAPI-CALLERSERVICEID']=String(serviceId);
var callersStr=options.variables.filter(function(row){return row.name=="signature"})[0]?.value;
var callers=JSON.parse(callersStr);
var caller = options.headers['HSB-OPENAPI-CALLERSERVICEID'];
var secret = callers[caller];
var paramStr  = JSON.stringify(options.requestBody);
var singnature = MD5(paramStr.concat("_",secret));
options.headers['HSB-OPENAPI-SIGNATURE']= String(singnature);
```



**请求示例：**
<!--doc.example.body-->
```json
{
   "scene":"3c",
   "Fxy_spuid":"1",
   "Fhsb_productid":"100",
   "Fxy_product_name":"闲鱼产品名称"
}
```

**返回结果：**
<!--doc.parameter id=responseParameter position=body httpStatus="200" _column="name,type,description"-->
|参数名|类型|说明|
|:-----  |:-----|----- |
|_ret |string   |0成功 1失败  |
|_errCode |string   |状态码  |
|_errStr |string   |返回信息  |

**返回示例：**

成功示例
<!--doc.example.response id=ok -->
```json
{
   "errcode":"0",
   "errmsg":"ok"
}

```
**返回示例：**

失败示例
<!--doc.example.response id=err -->
```json
{
   "errcode":"15785",
   "errmsg":"错误原因"
}
```

**说明：**

其它补充说明。。。
