package markdown2json

/*
* 根据传入的数据meta格式，整理[]*Record 数据，整理规则：
* 1. 将DB、Table 压入Uniqueue队头，
* 2. 按顺序从Uniqueue中取key，[]*Record 中取出值，拼接字符串后MD5,指纹相同的合并剩余属性
 */
type Table struct {
	Data []*Record
	Meta KVType
}

func (table *Table) RegisterMeta(meta KVType) {
	table.Meta = meta
}

type KVType struct {
	DB       string   `json:"db"`
	Table    string   `json:"table"`
	Uniqueue []string `json:"uniqueue"`
	Many     []string `json:"many"`
}

func FormatRecord(records []*Record) []*Record {
	return nil
}
