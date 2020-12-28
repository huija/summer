package utils

import "encoding/json"

// MergeStructByMarshal merge b to a by json.Marshal
// 0. a & b should be pointer
// 1. json tag: omitempty is necessary for every `no-zero` attr
// 2. The performance of this method is not very high for lagre struct
func MergeStructByMarshal(a, b interface{}) (err error) {
	marshal, err := json.Marshal(a)
	if err != nil {
		return
	}
	marshal2, err := json.Marshal(b)
	if err != nil {
		return
	}

	res := make(map[string]interface{})
	err = json.Unmarshal(marshal2, &res)
	if err != nil {
		return
	}
	err = json.Unmarshal(marshal, &res)
	if err != nil {
		return
	}

	bytes, err := json.Marshal(res)
	if err != nil {
		return
	}
	return json.Unmarshal(bytes, &a)
}
