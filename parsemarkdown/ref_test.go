package parsemarkdown_test

import (
	"fmt"
	"testing"

	"github.com/suifengpiao14/markdown2json/parsemarkdown"
)

func TestResolveRef(t *testing.T) {
	records := GetRecords()
	newRecords, err := parsemarkdown.ResolveRef(records)
	if err != nil {
		panic(err)
	}
	fmt.Println(newRecords.Json())

}
