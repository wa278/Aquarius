package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)
//转化int64
func GetInt64(v interface{}, errmsg string, defaultValue int64) (int64, error) {
	switch val := v.(type) {
	case float64:
		return int64(val), nil
	case float32:
		return int64(val), nil
	default:
		strv := fmt.Sprintf("%v", v)
		if i, err := strconv.ParseInt(strv, 10, 64); err != nil {
			if len(errmsg) > 0 {
				return defaultValue, errors.New(errmsg)
			} else {
				return defaultValue, nil
			}
		} else {
			return i, nil
		}
	}

}

func GetSafeInt64(v interface{}, defaultValue int64) int64 {
	val, _ := GetInt64(v, "", defaultValue)
	return val
}

//转化float64
func GetFloat64(v interface{}, errmsg string, defaultValue float64) (float64, error) {
	strv := fmt.Sprintf("%v", v)
	if fv, err := strconv.ParseFloat(strv, 64); err != nil {
		if len(errmsg) > 0 {
			return defaultValue, errors.New(errmsg)
		} else {
			return defaultValue, nil
		}
	} else {
		return fv, nil
	}
}

func GetSafeFloat64(v interface{}, defaultValue float64) float64 {
	val, _ := GetFloat64(v, "", defaultValue)
	return val
}
//转化字符串
func GetString(v interface{}, errmsg string, defaultValue string) (string, error) {
	if v == nil {
		v = ""
	}
	switch reflect.ValueOf(v).Kind() {
	case reflect.String, reflect.Int64, reflect.Uint64, reflect.Int32, reflect.Uint32, reflect.Int, reflect.Uint, reflect.Int16, reflect.Uint16, reflect.Int8, reflect.Uint8:
		val := fmt.Sprintf("%v", v)
		if len(val) > 0 {
			return val, nil
		}
		if len(errmsg) > 0 {
			return defaultValue, errors.New(errmsg)
		} else {
			return defaultValue, nil
		}
	case reflect.Float64, reflect.Float32:
		if f64, ok := v.(float64); ok {
			return strconv.FormatFloat(f64, 'f', -1, 64), nil
		}
		if f32, ok := v.(float32); ok {
			return strconv.FormatFloat(float64(f32), 'f', -1, 64), nil
		}
		return defaultValue, errors.New("Float can not convert.")
	case reflect.Bool:
		if b, ok := v.(bool); ok {
			if b {
				return "1", nil
			} else {
				return "0", nil
			}
		} else {
			return defaultValue, errors.New("Bool can not convert.")
		}
	default:
		if len(errmsg) > 0 {
			return defaultValue, errors.New(errmsg)
		} else {
			return defaultValue, nil
		}
	}
}

func GetSafeString(v interface{}, defaultValue string) string {
	val, _ := GetString(v, "", defaultValue)
	return val
}


func GetMustString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch reflect.ValueOf(v).Kind() {
	case reflect.Map, reflect.Slice, reflect.Array:
		jsn, _ := json.Marshal(v)
		return string(jsn)
	default:
		return fmt.Sprintf("%v", v)
	}
}
