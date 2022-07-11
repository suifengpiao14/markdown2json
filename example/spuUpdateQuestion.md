## heading {#id .className attrName=attrValue class="class1 class2"}
**简要描述：**
> 手动获取spu问卷
**协议：**
- 级别：二层
- 路径：<!--api.uri--> /spu/ajaxSpuUpdateQuestion

**环境：**
- 开发：193.112.197.63 http://opms.huishoubao.com.cn
- 测试：xx.xx.xx.xx  http://opms.huishoubao.com.cn
- 线上：<!--api.host-->http://opms.huishoubao.com.cn

**请求参数：**
|参数名|类型|必选|默认值|说明|
|:----    |:---|:----- |-----   |-----   |
|scene| string|是|-|场景 |
|Fxy_spu_id| string|是|-|闲鱼SPU ID |
|Fxy_product_name| string|是|-|闲鱼SPU 名称|
|Fhsb_product_id| string|是|-|回收宝产品ID |

<!--api.header=Hsb_service_id:1001 -->
<!--api.header=Hsb-signature:123 -->

<!--api.variable=signature:joenebfhefeh -->

<!--api.preRequest args=options
   var serviceId=options.variables.filter(function(row){return row.name=="serviceId"})[0]?.value;
options.headers['HSB-OPENAPI-CALLERSERVICEID']=String(serviceId);


var callersStr=options.variables.filter(function(row){return row.name=="signature"})[0]?.value;
var callers=JSON.parse(callersStr);
var caller = options.headers['HSB-OPENAPI-CALLERSERVICEID'];
var secret = callers[caller];
var paramStr  = JSON.stringify(options.requestBody);
var singnature = MD5(paramStr.concat("_",secret));
options.headers['HSB-OPENAPI-SIGNATURE']= String(singnature);

-->


**请求示例：**
<!--api.body -->
```urlencoded
scene=3c&Fxy_spu_id=1&Fhsb_product_id=100&Fxy_product_name=闲鱼产品名称
```

**返回结果：**
|参数名|类型|说明|
|:-----  |:-----|----- |
|_ret |string   |0成功 1失败  |
|_errCode |string   |状态码  |
|_errStr |string   |返回信息  |

**返回示例：**

成功示例
<!--api.response.example type=ok-->
```json
{
   "errcode":"0",
   "errmsg":"ok"
}

```
**返回示例：**

失败示例
<!--api.response.examle type=error-->
```json
{
   "errcode":"15785",
   "errmsg":"错误原因"
}
```

**说明：**

其它补充说明。。。
