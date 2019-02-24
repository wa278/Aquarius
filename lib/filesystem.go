package lib

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"strings"
)
//判断文件存在
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil && !os.IsExist(err) {
		return false
	}
	return true
}

//读取一行
func ReadLines(filename string) ([]string, error) {
	var content = make([]string, 0)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var buf = bufio.NewReader(file)
	for {
		//推荐方式
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if line != "" {
			content = append(content, line)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	return content, nil
}
//读取全部内容
func ReadFileAll(filename string) ([]byte, error) {
	file, fileErr := os.Open(filename)
	if fileErr != nil {
		return nil, fileErr
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}
//以追加方式写入
func FileAppendContent(filename, data string) error {
	file, fErr := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if fErr != nil {
		return fErr
	}
	defer file.Close()
	if _, wErr := file.WriteString(data); wErr != nil {
		return wErr
	}
	return nil
}
//重写方式写入
func FilePutContent(filename, data string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.WriteString(data); err != nil {
		return err
	}
	return nil
}

