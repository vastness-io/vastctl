package render

import "encoding/json"

type Render func(interface{}) (interface{}, error)

func PrettyPrintJSON(x interface{}) (interface{}, error) {
	b, err := json.MarshalIndent(x, "", "    ") //pretty print

	if err != nil {
		return nil, err
	}

	return string(b), err
}
