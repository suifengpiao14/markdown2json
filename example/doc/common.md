# 公共参数


***请求头***:
<Set ns="doc.request.header" encoding="markdown/table" column="name,type,required,default,description">

|参数名|类型|必选|默认值|说明|
|:----    |:---|:----- |-----   |-----   |
|content-type| string|是|application/json|文件格式|
|appid|string|是||访问服务的备案id|
|signature|string|是||签名,外网访问需开启签名|
</Set>


***返回示例(错误)***:
<Obj ns="doc.example.response.error" encoding="markdown/code">
```json
{
    "code":"xxx",
    "message":"xxx提示"
}
```
</Obj>

***域名***:
<Line id="host">http://localhost:80 </Line>

***服务器地址***:
<Obj sn="server">

</Obj>

***联系人***:<Line id="name">彭政</Line>(手机号:<Line id="phone">15999646794<Line>)

***签名算法***:
<Obj>
</Obj>