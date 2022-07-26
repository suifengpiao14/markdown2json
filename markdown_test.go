package markdown2json_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	parsemarkdown "github.com/suifengpiao14/markdown2json"
)

func TestXhtml2json(t *testing.T) {
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
	xmlByte, err := parsemarkdown.Markdown2XML(source)
	if err != nil {
		panic(err)
	}
	str, err := parsemarkdown.Xhtml2json(xmlByte)
	fmt.Println(str)
}
