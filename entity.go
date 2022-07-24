package markdown2json

const (
	API_METHOD_GET    = "get"
	API_METHOD_POST   = "post"
	API_METHOD_PUT    = "put"
	API_METHOD_DELETE = "delete"
	API_METHOD_HEAD   = "head"
)

type APIModel struct {
	// api 文档标识
	APIID string `json:"apiID"`
	// 服务标识
	ServiceID string `json:"serviceID"`
	// 路由名称(英文)
	Name string `json:"name"`
	// 标题
	Title string `json:"title"`
	// 标签ID集合
	Tags string `json:"tags"`
	// 路径
	URI string `json:"uri"`
	// 请求方法(get-GET,post-POST,put-PUT,delete-DELETE,head-HEAD)
	Method string `json:"method"`
	// 摘要
	Summary string `json:"summary"`
	// 介绍
	Description string `json:"description"`
	// 创建时间
	CreatedAt string `json:"createdAt"`
	// 修改时间
	UpdatedAt string `json:"updatedAt"`
	// 删除时间
	DeletedAt string `json:"deletedAt"`
}

func (t *APIModel) TableName() string {
	return "api"
}
func (t *APIModel) PrimaryKey() string {
	return "api_id"
}
func (t *APIModel) PrimaryKeyCamel() string {
	return "APIID"
}

func (t *APIModel) MethodTitleMap() map[string]string {
	enumMap := make(map[string]string)
	enumMap[API_METHOD_GET] = "GET"
	enumMap[API_METHOD_POST] = "POST"
	enumMap[API_METHOD_PUT] = "PUT"
	enumMap[API_METHOD_DELETE] = "DELETE"
	enumMap[API_METHOD_HEAD] = "HEAD"
	return enumMap
}
func (t *APIModel) MethodTitle() string {
	enumMap := t.MethodTitleMap()
	title, ok := enumMap[t.Method]
	if !ok {
		msg := "func MethodTitle not found title by key " + t.Method
		panic(msg)
	}
	return title
}

const (
	EXAMPLE_METHOD_GET                    = "get"
	EXAMPLE_METHOD_POST                   = "post"
	EXAMPLE_METHOD_PUT                    = "put"
	EXAMPLE_METHOD_DELETE                 = "delete"
	EXAMPLE_METHOD_HEAD                   = "head"
	EXAMPLE_CONTENT_TYPE_APPLICATION_JSON = "application/json"
	EXAMPLE_CONTENT_TYPE_PLAIN_TEXT       = "plain/text"
)

type ExampleModel struct {
	// 测试用例标识
	ExampleID string `json:"exampleID"`
	// 服务标识
	ServiceID string `json:"serviceID"`
	// api 文档标识
	APIID string `json:"apiID"`
	// 标签,mock数据时不同接口案例优先返回相同tag案例
	Tag string `json:"tag"`
	// 案例名称
	Title string `json:"title"`
	// 简介
	Summary string `json:"summary"`
	// URL
	URL string `json:"url"`
	// 请求方法(get-GET,post-POST,put-PUT,delete-DELETE,head-HEAD)
	Method string `json:"method"`
	// 前置请求脚本
	PreRequestScript string `json:"preRequestScript"`
	// query 鉴权
	Auth string `json:"auth"`
	// query 请求头
	Headers string `json:"headers"`
	// query 请求参数(json序列化)
	Parameters string `json:"parameters"`
	// 请求格式(application/json-json,plain/text-文本)
	ContentType string `json:"contentType"`
	// query 请求体
	Body string `json:"body"`
	// query 返回体测试脚本
	TestScript string `json:"testScript"`
	// 请求体
	Response string `json:"response"`
	// 创建时间
	CreatedAt string `json:"createdAt"`
	// 修改时间
	UpdatedAt string `json:"updatedAt"`
	// 删除时间
	DeletedAt string `json:"deletedAt"`
}

func (t *ExampleModel) TableName() string {
	return "example"
}
func (t *ExampleModel) PrimaryKey() string {
	return "example_id"
}
func (t *ExampleModel) PrimaryKeyCamel() string {
	return "ExampleID"
}
func (t *ExampleModel) MethodTitleMap() map[string]string {
	enumMap := make(map[string]string)
	enumMap[EXAMPLE_METHOD_GET] = "GET"
	enumMap[EXAMPLE_METHOD_POST] = "POST"
	enumMap[EXAMPLE_METHOD_PUT] = "PUT"
	enumMap[EXAMPLE_METHOD_DELETE] = "DELETE"
	enumMap[EXAMPLE_METHOD_HEAD] = "HEAD"
	return enumMap
}
func (t *ExampleModel) MethodTitle() string {
	enumMap := t.MethodTitleMap()
	title, ok := enumMap[t.Method]
	if !ok {
		msg := "func MethodTitle not found title by key " + t.Method
		panic(msg)
	}
	return title
}
func (t *ExampleModel) ContentTypeTitleMap() map[string]string {
	enumMap := make(map[string]string)
	enumMap[EXAMPLE_CONTENT_TYPE_APPLICATION_JSON] = "json"
	enumMap[EXAMPLE_CONTENT_TYPE_PLAIN_TEXT] = "文本"
	return enumMap
}
func (t *ExampleModel) ContentTypeTitle() string {
	enumMap := t.ContentTypeTitleMap()
	title, ok := enumMap[t.ContentType]
	if !ok {
		msg := "func ContentTypeTitle not found title by key " + t.ContentType
		panic(msg)
	}
	return title
}

type MarkdownModel struct {
	// markdown 文档标识
	MarkdownID string `json:"markdownID"`
	// 服务标识
	ServiceID string `json:"serviceID"`
	// api 文档标识
	APIID string `json:"apiID"`
	// 唯一名称
	Name string `json:"name"`
	// 文章标题
	Title string `json:"title"`
	//
	Markdown string `json:"markdown"`
	//
	Content string `json:"content"`
	// 作者ID
	OwnerID int `json:"ownerID"`
	// 作者名称
	OwnerName string `json:"ownerName"`
	// 创建时间
	CreatedAt string `json:"createdAt"`
	// 修改时间
	UpdatedAt string `json:"updatedAt"`
	// 删除时间
	DeletedAt string `json:"deletedAt"`
}

func (t *MarkdownModel) TableName() string {
	return "markdown"
}
func (t *MarkdownModel) PrimaryKey() string {
	return "markdown_id"
}
func (t *MarkdownModel) PrimaryKeyCamel() string {
	return "MarkdownID"
}

const (
	PARAMETER_TYPE_STRING             = "string"
	PARAMETER_TYPE_INT                = "int"
	PARAMETER_TYPE_NUMBER             = "number"
	PARAMETER_TYPE_ARRAY              = "array"
	PARAMETER_TYPE_OBJECT             = "object"
	PARAMETER_POSITION_BODY           = "body"
	PARAMETER_POSITION_HEAD           = "head"
	PARAMETER_POSITION_PATH           = "path"
	PARAMETER_POSITION_QUERY          = "query"
	PARAMETER_POSITION_COOKIE         = "cookie"
	PARAMETER_DEPRECATED_TRUE         = "true"
	PARAMETER_DEPRECATED_FALSE        = "false"
	PARAMETER_REQUIRED_TRUE           = "true"
	PARAMETER_REQUIRED_FALSE          = "false"
	PARAMETER_EXPLODE_TRUE            = "true"
	PARAMETER_EXPLODE_FALSE           = "false"
	PARAMETER_ALLOW_EMPTY_VALUE_TRUE  = "true"
	PARAMETER_ALLOW_EMPTY_VALUE_FALSE = "false"
	PARAMETER_ALLOW_RESERVED_TRUE     = "true"
	PARAMETER_ALLOW_RESERVED_FALSE    = "false"
)

type ParameterModel struct {
	// 参数标识
	ParameterID string `json:"parameterID"`
	// 服务标识
	ServiceID string `json:"serviceID"`
	// api 文档标识
	APIID string `json:"apiID"`
	// 验证规则标识
	SchemaID string `json:"schemaID"`
	// 全称
	FullName string `json:"fullName"`
	// 名称(冗余local.en)
	Name string `json:"name"`
	// 名称(冗余local.zh)
	Title string `json:"title"`
	// 参数类型(string-字符,int-整型,number-数字,array-数组,object-对象)
	Type string `json:"type"`
	// 所属标签
	Tag string `json:"tag"`

	// http状态码
	HTTPStatus string `json:"httpStatus"`
	// 参数所在的位置(body-BODY,head-HEAD,path-PATH,query-QUERY,cookie-COOKIE)
	Position string `json:"position"`
	// 案例
	Example string `json:"example"`
	// 是否弃用(true-是,false-否)
	Deprecated string `json:"deprecated"`
	// 是否必须(true-是,false-否)
	Required string `json:"required"`
	// 对数组、对象序列化方法,参照openapi parameters.style
	Serialize string `json:"serialize"`
	// 对象的key,是否单独成参数方式,参照openapi parameters.explode(true-是,false-否)
	Explode string `json:"explode"`
	// 是否容许空值(true-是,false-否)
	AllowEmptyValue string `json:"allowEmptyValue"`
	// 特殊字符是否容许出现在uri参数中(true-是,false-否)
	AllowReserved string `json:"allowReserved"`
	// 简介
	Description string `json:"description"`
	// 创建时间
	CreatedAt string `json:"createdAt"`
	// 修改时间
	UpdatedAt string `json:"updatedAt"`
	// 删除时间
	DeletedAt string `json:"deletedAt"`
}

func (t *ParameterModel) TableName() string {
	return "parameter"
}
func (t *ParameterModel) PrimaryKey() string {
	return "parameter_id"
}
func (t *ParameterModel) PrimaryKeyCamel() string {
	return "ParameterID"
}
func (t *ParameterModel) AllowEmptyValueTitleMap() map[string]string {
	enumMap := make(map[string]string)
	enumMap[PARAMETER_ALLOW_EMPTY_VALUE_TRUE] = "是"
	enumMap[PARAMETER_ALLOW_EMPTY_VALUE_FALSE] = "否"
	return enumMap
}
func (t *ParameterModel) AllowEmptyValueTitle() string {
	enumMap := t.AllowEmptyValueTitleMap()
	title, ok := enumMap[t.AllowEmptyValue]
	if !ok {
		msg := "func AllowEmptyValueTitle not found title by key " + t.AllowEmptyValue
		panic(msg)
	}
	return title
}
func (t *ParameterModel) AllowReservedTitleMap() map[string]string {
	enumMap := make(map[string]string)
	enumMap[PARAMETER_ALLOW_RESERVED_TRUE] = "是"
	enumMap[PARAMETER_ALLOW_RESERVED_FALSE] = "否"
	return enumMap
}
func (t *ParameterModel) AllowReservedTitle() string {
	enumMap := t.AllowReservedTitleMap()
	title, ok := enumMap[t.AllowReserved]
	if !ok {
		msg := "func AllowReservedTitle not found title by key " + t.AllowReserved
		panic(msg)
	}
	return title
}
func (t *ParameterModel) TypeTitleMap() map[string]string {
	enumMap := make(map[string]string)
	enumMap[PARAMETER_TYPE_STRING] = "字符"
	enumMap[PARAMETER_TYPE_INT] = "整型"
	enumMap[PARAMETER_TYPE_NUMBER] = "数字"
	enumMap[PARAMETER_TYPE_ARRAY] = "数组"
	enumMap[PARAMETER_TYPE_OBJECT] = "对象"
	return enumMap
}
func (t *ParameterModel) TypeTitle() string {
	enumMap := t.TypeTitleMap()
	title, ok := enumMap[t.Type]
	if !ok {
		msg := "func TypeTitle not found title by key " + t.Type
		panic(msg)
	}
	return title
}

func (t *ParameterModel) PositionTitleMap() map[string]string {
	enumMap := make(map[string]string)
	enumMap[PARAMETER_POSITION_BODY] = "BODY"
	enumMap[PARAMETER_POSITION_HEAD] = "HEAD"
	enumMap[PARAMETER_POSITION_PATH] = "PATH"
	enumMap[PARAMETER_POSITION_QUERY] = "QUERY"
	enumMap[PARAMETER_POSITION_COOKIE] = "COOKIE"
	return enumMap
}
func (t *ParameterModel) PositionTitle() string {
	enumMap := t.PositionTitleMap()
	title, ok := enumMap[t.Position]
	if !ok {
		msg := "func PositionTitle not found title by key " + t.Position
		panic(msg)
	}
	return title
}
func (t *ParameterModel) DeprecatedTitleMap() map[string]string {
	enumMap := make(map[string]string)
	enumMap[PARAMETER_DEPRECATED_TRUE] = "是"
	enumMap[PARAMETER_DEPRECATED_FALSE] = "否"
	return enumMap
}
func (t *ParameterModel) DeprecatedTitle() string {
	enumMap := t.DeprecatedTitleMap()
	title, ok := enumMap[t.Deprecated]
	if !ok {
		msg := "func DeprecatedTitle not found title by key " + t.Deprecated
		panic(msg)
	}
	return title
}
func (t *ParameterModel) RequiredTitleMap() map[string]string {
	enumMap := make(map[string]string)
	enumMap[PARAMETER_REQUIRED_TRUE] = "是"
	enumMap[PARAMETER_REQUIRED_FALSE] = "否"
	return enumMap
}
func (t *ParameterModel) RequiredTitle() string {
	enumMap := t.RequiredTitleMap()
	title, ok := enumMap[t.Required]
	if !ok {
		msg := "func RequiredTitle not found title by key " + t.Required
		panic(msg)
	}
	return title
}
func (t *ParameterModel) ExplodeTitleMap() map[string]string {
	enumMap := make(map[string]string)
	enumMap[PARAMETER_EXPLODE_TRUE] = "是"
	enumMap[PARAMETER_EXPLODE_FALSE] = "否"
	return enumMap
}
func (t *ParameterModel) ExplodeTitle() string {
	enumMap := t.ExplodeTitleMap()
	title, ok := enumMap[t.Explode]
	if !ok {
		msg := "func ExplodeTitle not found title by key " + t.Explode
		panic(msg)
	}
	return title
}

const (
	SCHEMA_TYPE_INTEGER            = "integer"
	SCHEMA_TYPE_ARRAY              = "array"
	SCHEMA_TYPE_STRING             = "string"
	SCHEMA_TYPE_OBJECT             = "object"
	SCHEMA_DEPRECATED_TRUE         = "true"
	SCHEMA_DEPRECATED_FALSE        = "false"
	SCHEMA_REQUIRED_TRUE           = "true"
	SCHEMA_REQUIRED_FALSE          = "false"
	SCHEMA_NULLABLE_TRUE           = "true"
	SCHEMA_NULLABLE_FALSE          = "false"
	SCHEMA_EXCLUSIVE_MAXIMUM_TRUE  = "true"
	SCHEMA_EXCLUSIVE_MAXIMUM_FALSE = "false"
	SCHEMA_EXCLUSIVE_MINIMUM_TRUE  = "true"
	SCHEMA_EXCLUSIVE_MINIMUM_FALSE = "false"
	SCHEMA_UNIQUE_ITEMS_TRUE       = "true"
	SCHEMA_UNIQUE_ITEMS_FALSE      = "false"
	SCHEMA_ALLOW_EMPTY_VALUE_TRUE  = "true"
	SCHEMA_ALLOW_EMPTY_VALUE_FALSE = "false"
	SCHEMA_ALLOW_RESERVED_TRUE     = "true"
	SCHEMA_ALLOW_RESERVED_FALSE    = "false"
	SCHEMA_READ_ONLY_TRUE          = "true"
	SCHEMA_READ_ONLY_FALSE         = "false"
	SCHEMA_WRITE_ONLY_TRUE         = "true"
	SCHEMA_WRITE_ONLY_FALSE        = "false"
)

type SchemaModel struct {
	// api schema 标识
	SchemaID string `json:"schemaID"`
	// 所属服务标识
	ServiceID string `json:"serviceID"`
	// 描述
	Description string `json:"description"`
	// 备注
	Remark string `json:"remark"`
	// 类型(integer-整数,array-数组,string-字符串,object-对象)
	Type string `json:"type"`
	// 案例
	Example string `json:"example"`
	// 是否弃用(true-是,false-否)
	Deprecated string `json:"deprecated"`
	// 是否必须(true-是,false-否)
	Required string `json:"required"`
	// 枚举值
	Enum string `json:"enum"`
	// 枚举名称
	EnumNames string `json:"enumNames"`
	// 枚举标题
	EnumTitles string `json:"enumTitles"`
	// 格式
	Format string `json:"format"`
	// 默认值
	Default string `json:"default"`
	// 是否可以为空(true-是,false-否)
	Nullable string `json:"nullable"`
	// 倍数
	MultipleOf int `json:"multipleOf"`
	// 最大值
	Maxnum int `json:"maxnum"`
	// 是否不包含最大项(true-是,false-否)
	ExclusiveMaximum string `json:"exclusiveMaximum"`
	// 最小值
	Minimum int `json:"minimum"`
	// 是否不包含最小项(true-是,false-否)
	ExclusiveMinimum string `json:"exclusiveMinimum"`
	// 最大长度
	MaxLength int `json:"maxLength"`
	// 最小长度
	MinLength int `json:"minLength"`
	// 正则表达式
	Pattern string `json:"pattern"`
	// 最大项数
	MaxItems int `json:"maxItems"`
	// 最小项数
	MinItems int `json:"minItems"`
	// 所有项是否需要唯一(true-是,false-否)
	UniqueItems string `json:"uniqueItems"`
	// 最多属性项
	MaxProperties int `json:"maxProperties"`
	// 最少属性项
	MinProperties int `json:"minProperties"`
	// 所有
	AllOf string `json:"allOf"`
	// 只满足一个
	OneOf string `json:"oneOf"`
	// 任何一个SchemaID
	AnyOf string `json:"anyOf"`
	// 是否容许空值(true-是,false-否)
	AllowEmptyValue string `json:"allowEmptyValue"`
	// 特殊字符是否容许出现在uri参数中(true-是,false-否)
	AllowReserved string `json:"allowReserved"`
	// 不包含的schemaID
	Not string `json:"not"`
	// boolean
	AdditionalProperties string `json:"additionalProperties"`
	// schema鉴别
	Discriminator string `json:"discriminator"`
	// 是否只读(true-是,false-否)
	ReadOnly string `json:"readOnly"`
	// 是否只写(true-是,false-否)
	WriteOnly string `json:"writeOnly"`
	// xml对象
	XML string `json:"xml"`
	// 附加文档
	ExternalDocs string `json:"externalDocs"`
	// 附加文档
	ExternalPros string `json:"externalPros"`
	// 扩展字段
	Extensions string `json:"extensions"`
	// 创建时间
	CreatedAt string `json:"createdAt"`
	// 修改时间
	UpdatedAt string `json:"updatedAt"`
	// 删除时间
	DeletedAt string `json:"deletedAt"`
	// 简介
	Summary string `json:"summary"`
}

func (t *SchemaModel) TableName() string {
	return "schema"
}
func (t *SchemaModel) PrimaryKey() string {
	return "schema_id"
}
func (t *SchemaModel) PrimaryKeyCamel() string {
	return "SchemaID"
}
func (t *SchemaModel) TypeTitleMap() map[string]string {
	enumMap := make(map[string]string)
	enumMap[SCHEMA_TYPE_INTEGER] = "整数"
	enumMap[SCHEMA_TYPE_ARRAY] = "数组"
	enumMap[SCHEMA_TYPE_STRING] = "字符串"
	enumMap[SCHEMA_TYPE_OBJECT] = "对象"
	return enumMap
}
func (t *SchemaModel) TypeTitle() string {
	enumMap := t.TypeTitleMap()
	title, ok := enumMap[t.Type]
	if !ok {
		msg := "func TypeTitle not found title by key " + t.Type
		panic(msg)
	}
	return title
}
func (t *SchemaModel) DeprecatedTitleMap() map[string]string {
	enumMap := make(map[string]string)
	enumMap[SCHEMA_DEPRECATED_TRUE] = "是"
	enumMap[SCHEMA_DEPRECATED_FALSE] = "否"
	return enumMap
}
func (t *SchemaModel) DeprecatedTitle() string {
	enumMap := t.DeprecatedTitleMap()
	title, ok := enumMap[t.Deprecated]
	if !ok {
		msg := "func DeprecatedTitle not found title by key " + t.Deprecated
		panic(msg)
	}
	return title
}
func (t *SchemaModel) UniqueItemsTitleMap() map[string]string {
	enumMap := make(map[string]string)
	enumMap[SCHEMA_UNIQUE_ITEMS_TRUE] = "是"
	enumMap[SCHEMA_UNIQUE_ITEMS_FALSE] = "否"
	return enumMap
}
func (t *SchemaModel) UniqueItemsTitle() string {
	enumMap := t.UniqueItemsTitleMap()
	title, ok := enumMap[t.UniqueItems]
	if !ok {
		msg := "func UniqueItemsTitle not found title by key " + t.UniqueItems
		panic(msg)
	}
	return title
}
func (t *SchemaModel) RequiredTitleMap() map[string]string {
	enumMap := make(map[string]string)
	enumMap[SCHEMA_REQUIRED_TRUE] = "是"
	enumMap[SCHEMA_REQUIRED_FALSE] = "否"
	return enumMap
}
func (t *SchemaModel) RequiredTitle() string {
	enumMap := t.RequiredTitleMap()
	title, ok := enumMap[t.Required]
	if !ok {
		msg := "func RequiredTitle not found title by key " + t.Required
		panic(msg)
	}
	return title
}
func (t *SchemaModel) NullableTitleMap() map[string]string {
	enumMap := make(map[string]string)
	enumMap[SCHEMA_NULLABLE_TRUE] = "是"
	enumMap[SCHEMA_NULLABLE_FALSE] = "否"
	return enumMap
}
func (t *SchemaModel) NullableTitle() string {
	enumMap := t.NullableTitleMap()
	title, ok := enumMap[t.Nullable]
	if !ok {
		msg := "func NullableTitle not found title by key " + t.Nullable
		panic(msg)
	}
	return title
}
func (t *SchemaModel) ExclusiveMaximumTitleMap() map[string]string {
	enumMap := make(map[string]string)
	enumMap[SCHEMA_EXCLUSIVE_MAXIMUM_TRUE] = "是"
	enumMap[SCHEMA_EXCLUSIVE_MAXIMUM_FALSE] = "否"
	return enumMap
}
func (t *SchemaModel) ExclusiveMaximumTitle() string {
	enumMap := t.ExclusiveMaximumTitleMap()
	title, ok := enumMap[t.ExclusiveMaximum]
	if !ok {
		msg := "func ExclusiveMaximumTitle not found title by key " + t.ExclusiveMaximum
		panic(msg)
	}
	return title
}
func (t *SchemaModel) ExclusiveMinimumTitleMap() map[string]string {
	enumMap := make(map[string]string)
	enumMap[SCHEMA_EXCLUSIVE_MINIMUM_TRUE] = "是"
	enumMap[SCHEMA_EXCLUSIVE_MINIMUM_FALSE] = "否"
	return enumMap
}
func (t *SchemaModel) ExclusiveMinimumTitle() string {
	enumMap := t.ExclusiveMinimumTitleMap()
	title, ok := enumMap[t.ExclusiveMinimum]
	if !ok {
		msg := "func ExclusiveMinimumTitle not found title by key " + t.ExclusiveMinimum
		panic(msg)
	}
	return title
}
func (t *SchemaModel) AllowEmptyValueTitleMap() map[string]string {
	enumMap := make(map[string]string)
	enumMap[SCHEMA_ALLOW_EMPTY_VALUE_TRUE] = "是"
	enumMap[SCHEMA_ALLOW_EMPTY_VALUE_FALSE] = "否"
	return enumMap
}
func (t *SchemaModel) AllowEmptyValueTitle() string {
	enumMap := t.AllowEmptyValueTitleMap()
	title, ok := enumMap[t.AllowEmptyValue]
	if !ok {
		msg := "func AllowEmptyValueTitle not found title by key " + t.AllowEmptyValue
		panic(msg)
	}
	return title
}
func (t *SchemaModel) AllowReservedTitleMap() map[string]string {
	enumMap := make(map[string]string)
	enumMap[SCHEMA_ALLOW_RESERVED_TRUE] = "是"
	enumMap[SCHEMA_ALLOW_RESERVED_FALSE] = "否"
	return enumMap
}
func (t *SchemaModel) AllowReservedTitle() string {
	enumMap := t.AllowReservedTitleMap()
	title, ok := enumMap[t.AllowReserved]
	if !ok {
		msg := "func AllowReservedTitle not found title by key " + t.AllowReserved
		panic(msg)
	}
	return title
}
func (t *SchemaModel) ReadOnlyTitleMap() map[string]string {
	enumMap := make(map[string]string)
	enumMap[SCHEMA_READ_ONLY_TRUE] = "是"
	enumMap[SCHEMA_READ_ONLY_FALSE] = "否"
	return enumMap
}
func (t *SchemaModel) ReadOnlyTitle() string {
	enumMap := t.ReadOnlyTitleMap()
	title, ok := enumMap[t.ReadOnly]
	if !ok {
		msg := "func ReadOnlyTitle not found title by key " + t.ReadOnly
		panic(msg)
	}
	return title
}
func (t *SchemaModel) WriteOnlyTitleMap() map[string]string {
	enumMap := make(map[string]string)
	enumMap[SCHEMA_WRITE_ONLY_TRUE] = "是"
	enumMap[SCHEMA_WRITE_ONLY_FALSE] = "否"
	return enumMap
}
func (t *SchemaModel) WriteOnlyTitle() string {
	enumMap := t.WriteOnlyTitleMap()
	title, ok := enumMap[t.WriteOnly]
	if !ok {
		msg := "func WriteOnlyTitle not found title by key " + t.WriteOnly
		panic(msg)
	}
	return title
}

type ServerModel struct {
	// 服务标识
	ServerID string `json:"serverID"`
	// 服务标识
	ServiceID string `json:"serviceID"`
	// 服务器地址
	URL string `json:"url"`
	// 介绍
	Description string `json:"description"`
	// 代理地址
	Proxy string `json:"proxy"`
	// 扩展字段
	ExtensionIds string `json:"extensionIds"`
	// 创建时间
	CreatedAt string `json:"createdAt"`
	// 修改时间
	UpdatedAt string `json:"updatedAt"`
	// 删除时间
	DeletedAt string `json:"deletedAt"`
}

func (t *ServerModel) TableName() string {
	return "server"
}
func (t *ServerModel) PrimaryKey() string {
	return "server_id"
}
func (t *ServerModel) PrimaryKeyCamel() string {
	return "ServerID"
}

type ServiceModel struct {
	// 服务标识
	ServiceID string `json:"serviceID"`
	// 标题
	Title string `json:"title"`
	// 介绍
	Description string `json:"description"`
	// 版本
	Version string `json:"version"`
	// 联系人
	ContactIds string `json:"contactIds"`
	// 协议
	License string `json:"license"`
	// 鉴权
	Security string `json:"security"`
	// 代理地址
	Proxy string `json:"proxy"`
	// json字符串
	Variables string `json:"variables"`
	// 创建时间
	CreatedAt string `json:"createdAt"`
	// 修改时间
	UpdatedAt string `json:"updatedAt"`
	// 删除时间
	DeletedAt string `json:"deletedAt"`
}

func (t *ServiceModel) TableName() string {
	return "service"
}
func (t *ServiceModel) PrimaryKey() string {
	return "service_id"
}
func (t *ServiceModel) PrimaryKeyCamel() string {
	return "ServiceID"
}
