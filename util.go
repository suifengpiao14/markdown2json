package markdown2json

func ConvertMapI2MapS(i interface{}) interface{} {
	if i == nil {
		return i
	}
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = ConvertMapI2MapS(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = ConvertMapI2MapS(v)
		}
	}
	return i
}
