package markdown2json_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	parsemarkdown "github.com/suifengpiao_14/markdown2json"
)

func TestParse(t *testing.T) {
	//file := "./example/first-doc.mdx"
	file := "./example/spuUpdateQuestion.md"
	fd, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	source, err := io.ReadAll(fd)
	if err != nil {
		panic(err)
	}
	apiElements, err := parsemarkdown.Parse(source)
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
