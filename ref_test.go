package markdown2json_test

import (
	"encoding/json"
	"fmt"
	"testing"

	parsemarkdown "github.com/suifengpiao14/markdown2json"
)

func TestResolveRef(t *testing.T) {
	records := GetRecords()

	rb, err := json.Marshal(records)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(rb))

	newRecords, err := parsemarkdown.ResolveRef(records)
	if err != nil {
		panic(err)
	}
	docRecords := newRecords.Filter(parsemarkdown.KV{Key: parsemarkdown.KEY_DB, Value: "doc"})
	b, err := json.Marshal(docRecords)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	fmt.Println(docRecords.String())

}
