package parsemark_test

import (
	"fmt"
	"testing"

	"github.com/suifengpiao14/markdown2json/parsemark"
)

func TestResolveRef(t *testing.T) {
	records := GetRecords()
	docRecords := records.FilterByKV(parsemark.KV{Key: parsemark.KEY_DB, Value: "doc"})
	fmt.Println(docRecords.String())

}
