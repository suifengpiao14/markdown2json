package markdown2json

import (
	"crypto/md5"
	"fmt"
)

func ConvertMapI2MapS(i interface{}) interface{} {
	if i == nil {
		return i
	}
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = ConvertMapI2MapS(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = ConvertMapI2MapS(v)
		}
	}
	return i
}

func Md5(s string) string {
	data := []byte(s)
	has := md5.Sum(data)
	out := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return out
}
