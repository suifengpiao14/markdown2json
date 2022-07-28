package markdown2json

/*
* 根据传入的数据meta格式，整理[]*Record 数据，整理规则：
* 1. 将DB、Table 压入Uniqueue队头，
* 2. 按顺序从Uniqueue中取key，[]*Record 中取出值，拼接字符串后MD5,指纹相同的合并剩余属性(剩余属性按照many元素拼接,其余覆盖方式合并)
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
	Uniqueue []string `json:"uniqueue"` //doc.server.proxy.url
	Many     []string `json:"many"`
}

const (
	INDEX_SERVER    = "proxy"
	INDEX_PARAMETER = "httpStatus,position,fullname"
)

func Merge(key string, records []*Record) (newRecords []*Record) {
	newRecords = make([]*Record, 0)
	mp := make(map[string]*Record)
	for _, record := range records {
		kv, ok := record.GetKVFirst(key)
		if !ok {
			kv = &KV{
				Key: key,
			}
		}
		existsRecord, ok := mp[kv.Value]
		if !ok {
			mp[kv.Value] = record
			continue
		}

	}
	return newRecords
}

func MergeSameIdentifyRecord(first, second Record) (newRecord Record) {
	first = append(first, second...)
}

func FormatRecord(records []*Record) []*Record {
	return nil
}
