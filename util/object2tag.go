package util

import (
	"encoding/json"
	"reflect"
	"strconv"
)

func Object2Tag(obj any, tag string) (map[string]string, error) {
	params := make(map[string]string, 0)
	if obj != nil {
		valueOf := reflect.ValueOf(obj)
		typeOf := reflect.TypeOf(obj)
		if reflect.TypeOf(obj).Kind() == reflect.Ptr {
			valueOf = reflect.ValueOf(obj).Elem()
			typeOf = reflect.TypeOf(obj).Elem()
		}
		numField := valueOf.NumField()
		for i := 0; i < numField; i++ {
			tag := typeOf.Field(i).Tag.Get(tag)
			if len(tag) > 0 && tag != "-" {
				switch valueOf.Field(i).Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16,
					reflect.Int32, reflect.Int64:
					params[tag] = strconv.FormatInt(valueOf.Field(i).Int(), 10)
				case reflect.Uint, reflect.Uint8, reflect.Uint16,
					reflect.Uint32, reflect.Uint64:
					params[tag] = strconv.FormatUint(valueOf.Field(i).Uint(), 10)
				case reflect.Float32, reflect.Float64:
					params[tag] = strconv.FormatFloat(valueOf.Field(i).Float(), 'f', -1, 64)
				case reflect.Bool:
					params[tag] = strconv.FormatBool(valueOf.Field(i).Bool())
				case reflect.String:
					if len(valueOf.Field(i).String()) > 0 {
						params[tag] = valueOf.Field(i).String()
					}
				case reflect.Map:
					if !valueOf.Field(i).IsNil() {
						bytes, err := json.Marshal(valueOf.Field(i).Interface())
						if err != nil {
							return nil, err
						} else {
							params[tag] = string(bytes)
						}
					}
				case reflect.Slice:
					if ss, ok := valueOf.Field(i).Interface().([]string); ok {
						var pv string
						for _, sv := range ss {
							pv += sv + ","
						}
						if len(pv) >= len(",") && pv[len(pv)-len(","):] == "," {
							pv = pv[:len(pv)-1]
						}
						if len(pv) > 0 {
							params[tag] = pv
						}
					}
				}
			}
		}
	}
	return params, nil
}
