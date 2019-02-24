package lib

import (
	"fmt"
	"strings"
	"time"
)

func GetConfigFloat(key string, defaultValue float64) (float64, error) {
	result, rErr := getConfigs()
	if rErr != nil {
		return 0, rErr
	}
	val, ok := result[key]
	if ok {
		return GetFloat64(val, fmt.Sprintf("Item `%s` is invalid in config file.", key), defaultValue)
	}
	return defaultValue, nil
}

func GetConfigInt64(key string, defaultValue int64) (int64, error) {
	result, rErr := getConfigs()
	if rErr != nil {
		return 0, rErr
	}
	val, ok := result[key]
	if ok {
		return GetInt64(val, fmt.Sprintf("Item `%s` is invalid in config file.", key), defaultValue)
	}
	return defaultValue, nil
}

func GetConfigString(key string, defaultValue string) (string, error) {
	result, rErr := getConfigs()
	if rErr != nil {
		return defaultValue, rErr
	}
	val, ok := result[key]
	if ok {
		return val, nil
	}
	return defaultValue, nil
}

func GetConfigStringList(key string, defaultValue []string) ([]string, error) {
	result, rErr := getConfigs()
	if rErr != nil {
		return defaultValue, rErr
	}
	val, ok := result[key]
	if ok {
		return StringToList(val,";"), nil
	}
	return defaultValue, nil
}
//获取config
func getConfigs() (map[string]string, error) {
	cache := NewCache()
	key := "SystemConfigFilesCache"

	if data := cache.GetMap(key); data == nil {
		confItems := make(map[string]string)

		pathConst := NewPathConst()
		confFiles := []string{
			pathConst.ConfigDir + pathConst.DS + "com-config.conf",
		}
		for _, fname := range confFiles {
			lines, linesErr := ReadLines(fname)
			if linesErr != nil {
				return nil, linesErr
			}
			if len(lines) > 0 {
				for _, line := range lines {
					if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
						continue
					}
					segs := strings.SplitN(line, "=", 2)
					if len(segs) != 2 {
						continue
					}
					itemKey := strings.TrimSpace(segs[0])
					itemValue := strings.TrimSpace(segs[1])
					confItems[itemKey] = itemValue
				}
			}
		}
		//将config以map的形式存在cache中
		cache.Set(key, confItems, time.Second*10)
		return confItems, nil
	} else {
		return data, nil
	}
}
