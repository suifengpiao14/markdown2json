package markdown2json

import (
	"crypto/md5"
	"fmt"
)

func Md5(s string) string {
	data := []byte(s)
	has := md5.Sum(data)
	out := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return out
}
