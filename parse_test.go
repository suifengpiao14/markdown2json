package markdown2json_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	parsemarkdown "github.com/suifengpiao14/markdown2json"
)

func TestParse(t *testing.T) {
	records := GetRecords()
	b, err := json.Marshal(records)
	if err != nil {
		panic(err)
	}
	str := string(b)
	fmt.Println(str)
}

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
	records, err := parsemarkdown.Parse(source)
	if err != nil {
		panic(err)
	}
	return records
}

func TestGetRefs(t *testing.T) {
	records := GetRecords()
	refRecords := records.GetRefs()
	fmt.Println(refRecords.String())
}

func TestMerge(t *testing.T) {
	records := GetRecords()
	newRecords := records.Format()
	b, err := json.Marshal(newRecords)
	if err != nil {
		panic(err)
	}
	str := string(b)
	fmt.Println(str)
}
func TestRecordString(t *testing.T) {
	records := GetRecords()
	str := records.Filter(parsemarkdown.KV{Key: parsemarkdown.KEY_DB, Value: "doc"}).String()
	fmt.Println(str)
}
