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

//Md5ID 在md5数据后面附带长度，降低md5撞key的可能性
func Md5ID(s string) string {
	data := []byte(s)
	has := md5.Sum(data)
	out := fmt.Sprintf("%x%d", has, len(s)) //将[]byte转成16进制

	return out
}
