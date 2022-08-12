package markdown2json_test

import (
	"fmt"
	"testing"

	"github.com/suifengpiao14/markdown2json"
)

func TestMd5(t *testing.T) {
	a := "N9KU3"
	b := "QlYRH"
	ma := markdown2json.Md5(a)
	mb := markdown2json.Md5(b)
	fmt.Println(ma)
	fmt.Println(mb)
}
