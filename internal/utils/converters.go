package utils

import (
	"bytes"
	"encoding/json"
)

/*

	This file contains helper to convert input to various output(type)

*/

// ToStruct convert the in interface to the out interface
func ToStruct(inputMap, out interface{}) error {

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(inputMap)

	if err != nil {
		return err
	}

	json.NewDecoder(buf).Decode(out)

	if err != nil {
		return err
	}

	return nil
}

// convert err to map
func ErrToMap(err error) map[string]interface{} {
	if err != nil {
		var serialized_err map[string]interface{}
		json_err, _ := json.Marshal(err)
		if err = json.Unmarshal(json_err, &serialized_err); err != nil {
			return map[string]interface{}{"all_fields": "parsing error"}
		}
		return serialized_err
	}

	return nil
}
