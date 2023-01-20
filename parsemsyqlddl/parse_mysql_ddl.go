package parsemsyqlddl

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	executor "github.com/bytewatch/ddl-executor"
	"github.com/pkg/errors"
	"goa.design/goa/codegen"
)

const (
	DEFAULT_VALUE_CURRENT_TIMESTAMP = "current_timestamp"
)

type Column struct {
	Name          string
	GoType        string
	DBType        string
	Comment       string
	Nullable      bool
	Enums         []string
	AutoIncrement bool
	DefaultValue  string
	OnUpdate      bool // 根据数据表ddl 配置
	Unsigned      bool
	Size          int
}

// IsDefaultValueCurrentTimestamp 判断默认值是否为自动填充时间
func (c *Column) IsDefaultValueCurrentTimestamp() bool {
	return strings.Contains(strings.ToLower(c.DefaultValue), DEFAULT_VALUE_CURRENT_TIMESTAMP) // 测试发现有 current_timestamp() 情况
}

type Enum struct {
	ConstKey        string // 枚举类型定义 常量 名称
	ConstValue      string // 枚举类型定义值
	Title           string // 枚举类型 标题（中文）
	ColumnNameCamel string //枚举类型分组（字段名，每个枚举字段有多个值，全表通用，需要分组）
	Type            string // 类型 int-整型，string-字符串，默认string
}

type Enums []*Enum

func (e Enums) Len() int { // 重写 Len() 方法
	return len(e)
}
func (e Enums) Swap(i, j int) { // 重写 Swap() 方法
	e[i], e[j] = e[j], e[i]
}
func (e Enums) Less(i, j int) bool { // 重写 Less() 方法， 从小到大排序
	return e[i].ConstKey < e[j].ConstKey
}

// UniqueItems 去重
func (e Enums) UniqueItems() (uniq Enums) {
	emap := make(map[string]*Enum)
	for _, enum := range e {
		emap[enum.ConstKey] = enum
	}
	uniq = Enums{}
	for _, enum := range emap {
		uniq = append(uniq, enum)
	}
	return
}

// ColumnNameCamels 获取所有分组
func (e Enums) ColumnNameCamels() (output []string) {
	columnNameCamelMap := make(map[string]string)
	for _, enum := range e {
		columnNameCamelMap[enum.ColumnNameCamel] = enum.ColumnNameCamel
	}
	output = make([]string, 0)
	for _, columnNameCamel := range columnNameCamelMap {
		output = append(output, columnNameCamel)
	}
	return
}

// GetByGroup 通过分组名称获取enum
func (e Enums) GetByColumnNameCamel(ColumnNameCamel string) (enums Enums) {
	enums = Enums{}
	for _, enum := range e {
		if enum.ColumnNameCamel == ColumnNameCamel {
			enums = append(enums, enum)
		}
	}
	return
}

type Table struct {
	DatabaseName string
	TableName    string
	PrimaryKey   string
	Columns      []*Column
	EnumsConst   Enums
	Comment      string
	TableDef     *executor.TableDef
}

const (
	ERROR_UNKNOW_DATABASE_SCAN_FORMAT = "Error 1049: Unknown database %s"
)

func GetDatabaseNameFromError(executorErr executor.Error) (dbName string, err error) {
	_, err = fmt.Sscanf(executorErr.Error(), ERROR_UNKNOW_DATABASE_SCAN_FORMAT, &dbName)
	if err != nil {
		return "", err
	}
	dbName = strings.Trim(dbName, "'")
	return dbName, nil
}

// TryExecDDLs 尝试解析ddls,其中,包含数据库不存在情况,自动创建
func TryExecDDLs(ddls string) (db *executor.Executor, err error) {

	var do = func() (err error) {
		conf := executor.NewDefaultConfig()
		db = executor.NewExecutor(conf)
		err = db.Exec(ddls)
		return err
	}
	for {
		err = do()
		executorErr, ok := err.(*executor.Error)
		if ok && executorErr.Code() == executor.ErrBadDB.Code() {
			dbName, err := GetDatabaseNameFromError(*executorErr)
			if err != nil {
				return nil, err
			}
			if dbName != "" {
				ddls = fmt.Sprintf("create database `%s`; %s", dbName, ddls)
			}
		} else {
			break
		}
	}
	if err != nil {
		return
	}
	return
}

func ParseDDL(ddls string, dbName string) (tables []*Table, err error) {
	tables = make([]*Table, 0)
	ddlDB := fmt.Sprintf("create database `%s`;use `%s`;", dbName, dbName)
	ddls = fmt.Sprintf("%s%s", ddlDB, ddls)
	db, err := TryExecDDLs(ddls)
	if err != nil {
		return
	}
	dbNameList := db.GetDatabases()
	for _, dbName := range dbNameList {
		tableNameList, err := db.GetTables(dbName)
		if err != nil {
			return nil, err
		}
		for _, tableName := range tableNameList {
			tableDef, err := db.GetTableDef(dbName, tableName)
			if err != nil {
				return nil, err
			}

			table, err := ConvertTabDef2Table(*tableDef)
			if err != nil {
				return nil, err
			}
			tables = append(tables, table)
		}
	}
	return tables, nil

}

func ConvertTabDef2Table(tableDef executor.TableDef) (table *Table, err error) {
	table = &Table{
		DatabaseName: tableDef.Database,
		TableName:    tableDef.Name,
		Columns:      make([]*Column, 0),
		EnumsConst:   Enums{},
		Comment:      tableDef.Comment,
	}
	for _, indice := range tableDef.Indices {
		if indice.Name == "PRIMARY" {
			table.PrimaryKey = indice.Columns[0] // 暂时取第一个为主键，不支持多字段主键
		}
	}
	for _, columnDef := range tableDef.Columns {

		goType, size, err := mysql2GoType(columnDef.Type, true)

		if err != nil {
			return nil, err
		}
		columnPt := &Column{
			Name:          columnDef.Name,
			GoType:        goType,
			Size:          size,
			DBType:        columnDef.Type,
			Unsigned:      columnDef.Unsigned,
			Comment:       columnDef.Comment,
			Nullable:      columnDef.Nullable,
			Enums:         columnDef.Elems,
			AutoIncrement: columnDef.AutoIncrement,
			DefaultValue:  columnDef.DefaultValue,
			OnUpdate:      columnDef.OnUpdate,
		}
		if len(columnPt.Enums) > 0 {
			subEnumConst := enumsConst("", columnPt)
			table.EnumsConst = append(table.EnumsConst, subEnumConst...)
		}
		table.Columns = append(table.Columns, columnPt)
	}
	return
}

// map for converting mysql type to golang types
var typeForMysqlToGo = map[string]string{
	"int":                "int",
	"integer":            "int",
	"tinyint":            "int",
	"smallint":           "int",
	"mediumint":          "int",
	"bigint":             "int",
	"int unsigned":       "int",
	"integer unsigned":   "int",
	"tinyint unsigned":   "int",
	"smallint unsigned":  "int",
	"mediumint unsigned": "int",
	"bigint unsigned":    "int",
	"bit":                "int",
	"bool":               "bool",
	"enum":               "string",
	"set":                "string",
	"varchar":            "string",
	"char":               "string",
	"tinytext":           "string",
	"mediumtext":         "string",
	"text":               "string",
	"longtext":           "string",
	"blob":               "string",
	"tinyblob":           "string",
	"mediumblob":         "string",
	"longblob":           "string",
	"date":               "time.Time", // time.Time or string
	"datetime":           "time.Time", // time.Time or string
	"timestamp":          "time.Time", // time.Time or string
	"time":               "time.Time", // time.Time or string
	"float":              "float64",
	"double":             "float64",
	"decimal":            "float64",
	"binary":             "string",
	"varbinary":          "string",
}

func mysql2GoType(mysqlType string, time2str bool) (goType string, size int, err error) {
	if time2str {
		typeForMysqlToGo["date"] = "string"
		typeForMysqlToGo["datetime"] = "string"
		typeForMysqlToGo["timestamp"] = "string"
		typeForMysqlToGo["time"] = "string"
	}
	subType := mysqlType
	index := strings.Index(mysqlType, "(")
	if index > -1 {
		endIndex := strings.Index(mysqlType, ")")
		if endIndex > -1 { //获取大小
			number := mysqlType[index+1 : endIndex]
			size, _ = strconv.Atoi(number)
		}
		subType = mysqlType[:index]

	}
	goType, ok := typeForMysqlToGo[subType]
	if !ok {
		err = errors.Errorf("mysql2GoType: not found mysql type %s to go type", mysqlType)
	}
	return

}

// 封装 goa.design/goa/v3/codegen 方便后续可定制
func ToCamel(name string) string {
	return codegen.CamelCase(name, true, true)
}

func ToLowerCamel(name string) string {
	return codegen.CamelCase(name, false, true)
}

func SnakeCase(name string) string {
	return codegen.SnakeCase(name)
}

func enumsConst(tablePrefix string, columnPt *Column) (enumsConsts Enums) {
	prefix := fmt.Sprintf("%s_%s", tablePrefix, columnPt.Name)
	enumsConsts = Enums{}
	comment := strings.ReplaceAll(columnPt.Comment, " ", ",") // 替换中文逗号(兼容空格和逗号另种分割符号)
	reg, err := regexp.Compile(`\W`)
	if err != nil {
		panic(err)
	}
	for _, constValue := range columnPt.Enums {
		constKey := fmt.Sprintf("%s_%s", prefix, constValue)
		valueFormat := fmt.Sprintf("%s-", constValue) // 枚举类型的comment 格式 value1-title1,value2-title2
		index := strings.Index(comment, valueFormat)
		if index < 0 {
			err := errors.Errorf("column %s(enum) comment except contains %s-xxx,got:%s", columnPt.Name, constValue, comment)
			panic(err)
		}
		title := comment[index+len(valueFormat):]
		comIndex := strings.Index(title, ",")
		if comIndex > -1 {
			title = title[:comIndex]
		} else {
			title = strings.TrimRight(title, " )")
		}
		constKey = reg.ReplaceAllString(constKey, "_") //替换非字母字符
		constKey = strings.ToUpper(constKey)
		enumsConst := &Enum{
			ConstKey:   constKey,
			ConstValue: constValue,
			Title:      title,
			Type:       "string",
		}
		enumsConsts = append(enumsConsts, enumsConst)
	}
	return
}
