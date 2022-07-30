package markdown2json_test

import (
	"fmt"
	"testing"

	parsemarkdown "github.com/suifengpiao14/markdown2json"
)

func TestResolveRef(t *testing.T) {
	records := GetRecords()
	newRecords, err := parsemarkdown.ResolveRef(records)
	if err != nil {
		panic(err)
	}
	fmt.Println(newRecords.Filter(parsemarkdown.KV{Key: parsemarkdown.KEY_DB, Value: "doc"}).String())

}
