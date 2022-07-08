**简要描述：**
> 手动获取spu问卷
**协议：**
- 级别：二层
- 路径：<!--api.request.uri-->/spu/ajaxSpuUpdateQuestion

**环境：**
- 开发：193.112.197.63 http://opms.huishoubao.com.cn
- 测试：xx.xx.xx.xx  http://opms.huishoubao.com.cn
- 线上：<!--api.request.host--> http://opms.huishoubao.com.cn

**请求参数：**
<!--api.contentType=application/json-->
|参数名|类型|必选|默认值|说明|
|:----    |:---|:----- |-----   |-----   |
|scene| string|是|-|场景 |
|Fxy_spu_id| string|是|-|闲鱼SPU ID |
|Fxy_product_name| string|是|-|闲鱼SPU 名称|
|Fhsb_product_id| string|是|-|回收宝产品ID |



**请求示例：**
<!--api.request.example position=body -->
```form
scene=3c&Fxy_spu_id=1&Fhsb_product_id=100&Fxy_product_name=闲鱼产品名称
```

**返回结果：**
<!--api.response-->
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
