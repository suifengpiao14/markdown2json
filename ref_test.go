package markdown2json_test

import (
	"fmt"
	"testing"

	parsemarkdown "github.com/suifengpiao14/markdown2json"
)

func TestResolveRef(t *testing.T) {
	records := GetRecords()
	docRecords := records.FilterByKV(parsemarkdown.KV{Key: parsemarkdown.KEY_DB, Value: "doc"})
	fmt.Println(docRecords.String())

}
