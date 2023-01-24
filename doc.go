package markdown2json

import (
	"strings"

	"github.com/suifengpiao14/markdown2json/filltemplateapidoc"
	"github.com/suifengpiao14/markdown2json/parsemarkdown"
)

func ParseDoc(str string) (s string, err error) {
	s = str
	preTplVariableCount := -1
	for {
		tplVariableCount := strings.Count(s, "{{")
		if preTplVariableCount == tplVariableCount {
			return s, nil
		}
		records, err := parsemarkdown.Parse([]byte(s))
		if err != nil {
			panic(err)
		}

		preTplVariableCount = tplVariableCount
		s, err = filltemplateapidoc.View(s, records.Json())
		if err != nil {
			return "", err
		}
	}
}
