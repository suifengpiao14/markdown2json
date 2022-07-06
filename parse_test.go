package markdown2json_test

import (
	"fmt"
	"testing"

	parsemarkdown "github.com/suifengpiao_14/markdown2json"
)

func TestParse(t *testing.T) {
	file := "./example/first-doc.mdx"
	str, err := parsemarkdown.Parse(file)
	if err != nil {
		panic(err)
	}
	fmt.Println(str)
}
