package markdown2json

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestParseDoc(t *testing.T) {
	s := GetConent()
	out, err := ParseDoc(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}

func GetConent() string {
	file := "./example/doc/adList.md"
	fd, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	source, err := io.ReadAll(fd)
	if err != nil {
		panic(err)
	}
	s := string(source)
	return s

}
