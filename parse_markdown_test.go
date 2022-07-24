package markdown2json_test

import (
	"encoding/json"
	"fmt"
	"testing"

	parsemarkdown "github.com/suifengpiao_14/markdown2json"
)

func TestParse(t *testing.T) {
	//file := "./example/first-doc.mdx"
	file := "./example/spuUpdateQuestion.md"
	apiElements, err := parsemarkdown.Parse(file)
	if err != nil {
		panic(err)
	}
	b, err := json.Marshal(apiElements)
	if err != nil {
		panic(err)
	}
	str := string(b)
	fmt.Println(str)
}
