package parsemark_test

import (
	"fmt"
	"testing"

	parsemark "github.com/suifengpiao14/markdown2json/parsemark_back"
)

func TestResolveRef(t *testing.T) {
	records := GetRecords()
	docRecords := records.FilterByKV(parsemark.KV{Key: parsemark.KEY_DB, Value: "doc"})
	fmt.Println(docRecords.String())

}
