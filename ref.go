package markdown2json

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
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
		path := refAttr.Value
		var source []byte
		u, err := url.ParseRequestURI(path)
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
				p = p[1:] // window 下删除开头的/
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
		subRecords, err := Parse(source)
		if err != nil {
			err = errors.WithMessagef(err, "ref:%s", refAttr.Value)
			return nil, err
		}
		indexAttr, ok := refRecord.GetKV(KEY_INER_INDEX)
		if ok {
			subRecords = subRecords.GetByIndexWithChildren(indexAttr.Value) //筛选 和ref record id 相同的元素及其子元素
			for _, subRecord := range subRecords {
				inerRefAttr, ok := subRecord.GetKV(KEY_INER_REF)
				if !ok {
					inerRefAttr = &KV{
						Key:   KEY_INER_REF,
						Value: refAttr.Value,
					}
				} else {
					inerRefAttr.Value = fmt.Sprintf("%s,%s", inerRefAttr.Value, refAttr.Value)
				}
				subRecord.AddKV(*inerRefAttr)
			}
		}
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
