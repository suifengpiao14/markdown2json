# 公共参数

### 公共HTTP请求头
> 注意: 列举的名称为_head 属性名称，具体请求格式参考案例
<!--doc.parameter  _column="name,type,required,example,description" position=head-->
| 名称|类型|必选|案例|说明|
|:--|:--|:--|:--|:--|
|Content-Type|string|是|application/json|请求格式,当前只支持<!--enum-->application/json|
|HSB-OPENAPI-CALLERSERVICEID|string|是|110001|发起请求方的服务ID|
|HSB-OPENAPI-SIGNATURE|string|是|request|签名值|

### 公共参数
> 注意: 列举的名称为_head 属性名称，具体请求格式参考案例
<!--doc.parameter  _column="name,type,required,example,description"  position=body-->
| 名称|类型|必选|案例|说明|
|:--|:--|:--|:--|:--|
|_head._version|string|是|0.01|协议版本号<br/>可选值:0.01|
|_head._msgType|string|是|request|报文类型<br/>可选值:request(请求)、response(响应)|
|_head._timestamps|string|是|1523330331|请求时间戳(单位毫秒)|
|_head._invokeId|string|是|book1523330331358|当前请求标识(每次请求要求唯一)|
|_head._callerServiceId|string|是|110001|发起http请求方的服务ID|
|_head._groupNo|string|是|1|请求分组号|
|_head._interface|string|是|templateList|请求接口标识|
|_head._remark|string|是|0.01|备注|



### 公共参数案例
<!--doc.example.body -->
```json
{
    "_head":{
        "_version":"0.01",
        "_msgType":"request",
        "_timestamps":"1523330331",
        "_invokeId":"book1523330331358",
        "_callerServiceId":"110001",
        "_groupNo":"1",
        "_interface":"templateList",
        "_remark":""
    }
}
```