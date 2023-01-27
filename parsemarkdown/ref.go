package parsemarkdown

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/suifengpiao14/markdown2json/parsemsyqlddl"
)

/**
* ref 处理流程
* 1. markdown/html 解析成 records 不合并(ref 属性不扩散,即不复制到子元素上)
* 2. 获取所有含有ref属性的record,逐个解析
* 3. curl/read file 获取 markdown/html 数据
* 4. 解析 markdown/html
* 5. 筛选 和ref record id 相同的元素及其子元素 增加到 records 末尾,降低优先级
* 6. 当前元素删除ref属性
* 8. 重复3-7
 */
//ResolveRef 解决markdown/html 中的 ref 标记
func ResolveRef(records Records) (newRecords Records, err error) {
	newRecords = Records{}
	newRecords = append(newRecords, records...)
	refRecords := newRecords.GetRefs()
	i := 0
	for len(refRecords) > i { // refRecords 有元素追加
		refRecord := refRecords[i]
		i++
		refAttr, _ := refRecord.GetKV(KEY_REF)
		refRecord.PopKV(refAttr.Key) // 删除 ref key
		uri := refAttr.Value
		var source []byte
		u, err := url.Parse(uri)
		if err != nil {
			return nil, err
		}
		switch u.Scheme {
		case "http", "https":
			source, err = LoadFromURL(u.String())
			if err != nil {
				return nil, err
			}
		case "file":
			p := u.Path
			if len(p) > 3 && p[0] == '/' && p[2] == ':' {
				p = p[1:] // window 下删除开头的/ ,第一个/ 代表host
			}
			if len(p) > 2 && p[0] == '/' && p[1] == '.' {
				p = p[1:] // 相对路径时,删除开头的/,第一个/ 代表host
			}
			source, err = LoadFromFile(p)
			if err != nil {
				return nil, err
			}
		default:
			err := errors.Errorf("unsuport scheme:%s", u.String())
			return nil, err

		}

		if source == nil { // 内容为空,跳过
			continue
		}
		ext := path.Ext(u.Path)
		subRecords := Records{}
		switch ext {
		case ".md", ".markdown":
			subRecords, err = ParseWithRef(source)

		case ".sql":
			ddls := string(source)
			baseName := path.Base(u.Path)
			dbName := strings.TrimSuffix(baseName, ext)
			subRecords, err = ParsSQLDDL(ddls, dbName)

		}
		if err != nil {
			err = errors.WithMessagef(err, "ref:%s", refAttr.Value)
			return nil, err
		}

		if u.Fragment != "" {
			fragmentRecords := subRecords.GetByTagWithChildren(u.Fragment)
			subRecords = Records{}                   //重置 引用集合
			for _, record := range fragmentRecords { // 重置标签kv 值
				oldTagName := record.GetTag()
				// 替换tag值
				newTagName := strings.Replace(oldTagName, u.Fragment, refRecord.GetTag(), 1)
				record.SetKV(KV{
					Key:   KEY_TAG,
					Value: newTagName,
				})
				subRecords = append(subRecords, record)
			}
		}

		// 统一重置除了 _tag、_text、_ref 之外的所有其它属性
		tmpRcords := Records{}

		for _, record := range subRecords {
			for _, kv := range refRecord {
				if kv.Key == KEY_TAG || kv.Key == KEY_TEXT || kv.Key == KEY_REF || kv.Key == KEY_IS_END_BY_BACKSLASH {
					continue
				}
				record.SetKV(*kv)
			}
			tmpRcords = append(tmpRcords, record)
		}

		subRecords = tmpRcords
		newRecords = append(newRecords, subRecords...)
		subRefRecords := subRecords.GetRefs()
		refRecords = append(refRecords, subRefRecords...)
	}

	return newRecords, nil
}

func LoadFromURL(u string) (md []byte, err error) {
	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		Get(u)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		err := errors.Errorf("http code %d,body:%s", resp.StatusCode(), resp.Body())
		return nil, err
	}
	return resp.Body(), nil
}

func LoadFromFile(file string) (md []byte, err error) {
	fd, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	source, err := io.ReadAll(fd)
	if err != nil {
		return nil, err
	}
	return source, nil
}

func ParsSQLDDL(ddls string, dbName string) (records Records, err error) {
	records = Records{}
	tables, err := parsemsyqlddl.ParseDDL(ddls, dbName)
	if err != nil {
		return nil, err
	}
	for _, table := range tables {
		for _, colum := range table.Columns {
			record := Record{}
			record.SetKV(KV{
				Key:   KEY_TAG,
				Value: fmt.Sprintf("%s.%s.%s", table.DatabaseName, table.TableName, colum.Name),
			})
			record.SetKV(KV{
				Key:   "goType",
				Value: colum.GoType,
			})
			record.SetKV(KV{
				Key:   "comment",
				Value: colum.Comment,
			})

			record.SetKV(KV{
				Key:   "nullable",
				Value: strconv.FormatBool(colum.Nullable),
			})
			record.SetKV(KV{
				Key:   "enums",
				Value: strings.Join(colum.Enums, ","),
			})
			record.SetKV(KV{
				Key:   "default",
				Value: colum.DefaultValue,
			})
			record.SetKV(KV{
				Key:   "unsigned",
				Value: strconv.FormatBool(colum.Unsigned),
			})
			record.SetKV(KV{
				Key:   "size",
				Value: strconv.Itoa(colum.Size),
			})
			records = append(records, record)

		}

	}
	return records, err
}
