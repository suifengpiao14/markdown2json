package parseapi_test

import (
	"io"
	"os"

	"github.com/suifengpiao14/markdown2json/parsemarkdown"
)

func GetRecords() parsemarkdown.Records {
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
	records, err := parsemarkdown.ParseWithRef(source)
	if err != nil {
		panic(err)
	}
	return records
}
