package jsonnumber

import (
	"bytes"
	"encoding/json"
	"reflect"
)

// decode object using `json.Number` and convert to actual number type.
func decodeUseNumber(content []byte, in interface{}) error {
	// decode using json number
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.UseNumber()
	err := decoder.Decode(in)
	if err != nil {
		return err
	}
	var numInt int64

	v := reflect.ValueOf(in).Elem()
	return iterateMapFields(v, func(mField map[string]interface{}) error {
		for k, v := range mField {
			if jsonNum, ok := v.(json.Number); ok {
				// try with int first
				if numInt, err = jsonNum.Int64(); err == nil {
					mField[k] = int(numInt)
					continue
				}
				mField[k], err = jsonNum.Float64()
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

}

// fetch all `map[string]interface` fields and execute function on that field value
func iterateMapFields(v reflect.Value, fn func(map[string]interface{}) error) (err error) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem() // dereference pointer
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch field.Kind() {
		case reflect.Map:
			// possibly a map[string]interface field
			m, ok := field.Interface().(map[string]interface{})
			if ok {
				if err = fn(m); err != nil {
					return
				}
			}
		case reflect.Struct:
			// recusive loop map fields
			if err = iterateMapFields(field, fn); err != nil {
				return
			}
		}
	}
	return
}
