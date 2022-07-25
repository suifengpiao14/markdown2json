package markdown2json_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	parsemarkdown "github.com/suifengpiao_14/markdown2json"
)

func TestParseMarkdown(t *testing.T) {
	file := "./example/spuUpdateQuestion.md"
	fd, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	source, err := io.ReadAll(fd)
	if err != nil {
		panic(err)
	}
	html := parsemarkdown.ParseMarkdown(source)
	fmt.Println(html)
}

func TestQueryXML(t *testing.T) {
	file := "./example/spuUpdateQuestion.md"
	fd, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	source, err := io.ReadAll(fd)
	if err != nil {
		panic(err)
	}
	html := parsemarkdown.ParseMarkdown(source)
	err = parsemarkdown.QueryXML(html)
	if err != nil {
		panic(err)
	}
}

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
