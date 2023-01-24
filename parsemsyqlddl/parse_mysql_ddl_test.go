package parsemsyqlddl

import (
	"fmt"
	"os"
	"testing"
)

func TestParseDDL(t *testing.T) {
	ddls := GetDDL()
	tables, err := ParseDDL(ddls, "ad")
	if err != nil {
		panic(err)
	}
	fmt.Println(tables)
}

func GetDDL() (ddl string) {
	filename := "./example/ad.sql"
	b, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	ddl = string(b)
	return ddl
}
