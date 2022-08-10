package markdown2json_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/suifengpiao14/markdown2json"
)

func TestMarkdown2XML(t *testing.T) {
	file := "./example/spuUpdateQuestion.md"
	fd, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	source, err := io.ReadAll(fd)
	if err != nil {
		panic(err)
	}
	out := markdown2json.Markdown2XML(source)
	fmt.Println(out)

}
